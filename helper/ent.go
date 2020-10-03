// Copyright 2020 tree xie
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package helper

import (
	"context"
	"database/sql"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/facebook/ent/dialect"
	entsql "github.com/facebook/ent/dialect/sql"
	"github.com/iancoleman/strcase"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/vicanso/elite/config"
	"github.com/vicanso/elite/cs"
	"github.com/vicanso/elite/ent"
	"github.com/vicanso/elite/ent/hook"
	"github.com/vicanso/elite/ent/migrate"
	"go.uber.org/atomic"
	"go.uber.org/zap"
)

var (
	defaultEntDriver, defaultEntClient = initClientX()
)
var (
	initSchemaOnce sync.Once

	maskRegExp = regexp.MustCompile(`(?i)password`)
)

// processingKeyAll 记录所有表的正在处理请求
const processingKeyAll = "All"

// entProcessingStats ent的处理请求统计
type entProcessingStats struct {
	data map[string]*atomic.Uint32
}

var currentEntProcessingStats = new(entProcessingStats)

// initClientX 初始化客户端与driver
func initClientX() (*entsql.Driver, *ent.Client) {
	postgresConfig := config.GetPostgresConfig()

	maskURI := postgresConfig.URI
	urlInfo, _ := url.Parse(maskURI)
	if urlInfo != nil {
		pass, ok := urlInfo.User.Password()
		if ok {
			maskURI = strings.ReplaceAll(maskURI, pass, "***")
		}
	}
	logger.Info("connect postgres",
		zap.String("uri", maskURI),
	)
	db, err := sql.Open("pgx", postgresConfig.URI)
	if err != nil {
		panic(err)
	}

	// Create an ent.Driver from `db`.
	driver := entsql.OpenDB(dialect.Postgres, db)
	c := ent.NewClient(ent.Driver(driver))

	ctx := context.Background()
	if err := c.Schema.Create(ctx); err != nil {
		panic(err)
	}
	initSchemaHooks(c)
	return driver, c
}

// init 初始化统计
func (stats *entProcessingStats) init(schemas []string) {
	data := make(map[string]*atomic.Uint32)
	data[processingKeyAll] = atomic.NewUint32(0)
	for _, schema := range schemas {
		data[schema] = atomic.NewUint32(0)
	}
	stats.data = data
}

// inc 处理数+1
func (stats *entProcessingStats) inc(schema string) (uint32, uint32) {
	total := stats.data[processingKeyAll].Inc()
	p, ok := stats.data[schema]
	if !ok {
		return total, 0
	}
	return total, p.Inc()
}

// desc 处理数-1
func (stats *entProcessingStats) dec(schema string) (uint32, uint32) {
	total := stats.data[processingKeyAll].Dec()
	p, ok := stats.data[schema]
	if !ok {
		return total, 0
	}
	return total, p.Dec()
}

// initSchemaHooks 初始化相关的hooks
func initSchemaHooks(c *ent.Client) {
	schemas := make([]string, len(migrate.Tables))
	for index, table := range migrate.Tables {
		name := strcase.ToCamel(table.Name)
		// 去除最后的复数s
		schemas[index] = name[:len(name)-1]
	}
	currentEntProcessingStats.init(schemas)
	// 禁止删除数据
	c.Use(hook.Reject(ent.OpDelete | ent.OpDeleteOne))
	// 数据库操作统计
	c.Use(func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			schemaType := m.Type()
			totalProcessing, processing := currentEntProcessingStats.inc(schemaType)
			defer currentEntProcessingStats.dec(schemaType)
			op := m.Op().String()

			startedAt := time.Now()
			result := 0
			message := ""
			value, err := next.Mutate(ctx, m)
			// 如果失败，则记录出错信息
			if err != nil {
				result = 1
				message = err.Error()
			}
			data := make(map[string]interface{})
			for _, name := range m.Fields() {
				// 更新时间字段忽略
				if name == "updated_at" {
					continue
				}
				value, ok := m.Field(name)
				if !ok {
					continue
				}
				if maskRegExp.MatchString(name) {
					data[name] = "***"
				} else {
					data[name] = value
				}
			}

			d := time.Since(startedAt)
			logger.Info("ent stats",
				zap.String("schema", schemaType),
				zap.String("op", op),
				zap.Int("result", result),
				zap.Uint32("processing", processing),
				zap.Uint32("totalProcessing", totalProcessing),
				zap.String("use", d.String()),
				zap.Any("data", data),
				zap.String("message", message),
			)
			fields := map[string]interface{}{
				"processing":      processing,
				"totalProcessing": totalProcessing,
				"use":             int(d.Milliseconds()),
				"data":            data,
				"message":         message,
			}
			tags := map[string]string{
				"schema": schemaType,
				"op":     op,
				"result": strconv.Itoa(result),
			}
			GetInfluxSrv().Write(cs.MeasurementEntOP, fields, tags)
			return value, err
		})
	})
}

// EntGetStats get ent stats
func EntGetStats() map[string]interface{} {
	info := defaultEntDriver.DB().Stats()
	stats := map[string]interface{}{
		"maxOpenConnections": info.MaxOpenConnections,
		"openConnections":    info.OpenConnections,
		"inUse":              info.InUse,
		"idle":               info.Idle,
		"waitCount":          info.WaitCount,
		"waitDuration":       info.WaitDuration.Milliseconds(),
		"maxIdleClosed":      info.MaxIdleClosed,
		"maxIdleTimeClosed":  info.MaxIdleTimeClosed,
		"maxLifetimeClosed":  info.MaxLifetimeClosed,
	}
	for name, p := range currentEntProcessingStats.data {
		stats[strcase.ToLowerCamel(name)] = p.Load()
	}
	return stats
}

// EntGetClient get ent client
func EntGetClient() *ent.Client {
	return defaultEntClient
}

// EntPing ent driver ping
func EntPing() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return defaultEntDriver.DB().PingContext(ctx)
}

// EntInitSchema 初始化schema
func EntInitSchema() (err error) {
	// 只执行一次shcema初始化以及hook
	initSchemaOnce.Do(func() {
		err = defaultEntClient.Schema.Create(context.Background())
	})
	return
}
