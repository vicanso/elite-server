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

package schedule

import (
	"os"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/vicanso/elite/cs"
	"github.com/vicanso/elite/helper"
	"github.com/vicanso/elite/log"
	"github.com/vicanso/elite/novel"
	"github.com/vicanso/elite/service"

	"go.uber.org/zap"
)

type (
	taskFn      func() error
	statsTaskFn func() map[string]interface{}
)

var (
	logger = log.Default()
)

func init() {
	c := cron.New()
	_, _ = c.AddFunc("@every 1m", redisPing)
	_, _ = c.AddFunc("@every 1m", entPing)
	_, _ = c.AddFunc("@every 1m", configRefresh)
	_, _ = c.AddFunc("@every 5m", redisStats)
	_, _ = c.AddFunc("@every 5m", entStats)
	_, _ = c.AddFunc("@every 30s", cpuUsageStats)
	_, _ = c.AddFunc("@every 1m", performanceStats)
	_, _ = c.AddFunc("@every 24h", updateAllNovelWordCount)
	_, _ = c.AddFunc("@every 24h", updateAllNovelUpdatedWeight)
	// 每小时更新权重>=50的小说
	_, _ = c.AddFunc("@every 1h", newUpdateNovelChapterByWeight(50))
	// 每12小时更新权重>=10的小说
	_, _ = c.AddFunc("@every 12h", newUpdateNovelChapterByWeight(10))
	// 每24小时更新权重>=1的小说
	_, _ = c.AddFunc("@every 24h", newUpdateNovelChapterByWeight(1))

	if os.Getenv("SYNC_SOURCE") != "" {
		// _, _ = c.AddFunc("@every 12h", syncNovelSource)
		go syncNovelSource()
	}
	c.Start()
}

func doTask(desc string, fn taskFn) {
	startedAt := time.Now()
	err := fn()
	if err != nil {
		logger.Error(desc+" fail",
			zap.String("category", "schedule"),
			zap.Duration("use", time.Since(startedAt)),
			zap.Error(err),
		)
		service.AlarmError(desc + " fail, " + err.Error())
	} else {
		logger.Info(desc+" success",
			zap.String("category", "schedule"),
			zap.Duration("use", time.Since(startedAt)),
		)
	}
}

func doStatsTask(desc string, fn statsTaskFn) {
	startedAt := time.Now()
	stats := fn()
	logger.Info(desc,
		zap.String("category", "schedule"),
		zap.Duration("use", time.Since(startedAt)),
		zap.Any("stats", stats),
	)
}

func redisPing() {
	doTask("redis ping", helper.RedisPing)
}

func configRefresh() {
	configSrv := new(service.ConfigurationSrv)
	doTask("config refresh", configSrv.Refresh)
}

func redisStats() {
	doStatsTask("redis stats", func() map[string]interface{} {
		// 统计中除了redis数据库的统计，还有当前实例的统计指标，因此所有实例都会写入统计
		stats := helper.RedisStats()
		helper.GetInfluxSrv().Write(cs.MeasurementRedisStats, stats, nil)
		return stats
	})
}

func entPing() {
	doTask("ent ping", helper.EntPing)
}

// entStats ent的性能统计
func entStats() {
	doStatsTask("ent stats", func() map[string]interface{} {
		stats := helper.EntGetStats()
		helper.GetInfluxSrv().Write(cs.MeasurementEntStats, stats, nil)
		return stats
	})
}

// syncNovelSource 同步小说源
func syncNovelSource() {
	srv := novel.Srv{}
	doTask("sync novel source", srv.SyncSource)
}

// updateAllNovelWordCount 更新小说总字数
func updateAllNovelWordCount() {
	srv := novel.Srv{}
	doTask("update all novel word count", srv.UpdateAllWordCount)
}

// newUpdateNovelChapterByWeight 创建更新小说章节任务
func newUpdateNovelChapterByWeight(updatedWeight int) func() {
	return func() {
		// 暂时不更新，等所有小说同步完成后再启用
		// srv := novel.Srv{}
		// doTask("update novel chapters by weight", func() error {
		// 	return srv.UpdateAllChaptersByWeight(updatedWeight)
		// })
	}
}

// updateAllNovelUpdatedWeight 更新小说更新权重
func updateAllNovelUpdatedWeight() {
	srv := novel.Srv{}
	doTask("update novel updated weight", srv.UpdateAllUpdatedWeight)
}

// cpuUsageStats cpu使用率
func cpuUsageStats() {
	doTask("update cpu usage", service.UpdateCPUUsage)
}

// prevMemFrees 上一次 memory objects 释放的数量
var prevMemFrees uint64

// prevNumGC 上一次 gc 的次数
var prevNumGC uint32

// prevPauseTotal 上一次 pause 的总时长
var prevPauseTotal time.Duration

// performanceStats 系统性能
func performanceStats() {
	doStatsTask("performance stats", func() map[string]interface{} {
		data := service.GetPerformance()
		fields := map[string]interface{}{
			"goMaxProcs":   data.GoMaxProcs,
			"concurrency":  data.Concurrency,
			"memSys":       data.MemSys,
			"memHeapSys":   data.MemHeapSys,
			"memHeapInuse": data.MemHeapInuse,
			"memFrees":     data.MemFrees - prevMemFrees,
			"routineCount": data.RoutineCount,
			"cpuUsage":     data.CPUUsage,
			"numGC":        data.NumGC - prevNumGC,
			"pause":        (data.PauseTotalNs - prevPauseTotal).Milliseconds(),
		}
		prevMemFrees = data.MemFrees
		prevNumGC = data.NumGC

		helper.GetInfluxSrv().Write(cs.MeasurementPerformance, fields, nil)
		return fields
	})
}
