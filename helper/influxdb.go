// Copyright 2019 tree xie
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
	"time"

	influxdb "github.com/influxdata/influxdb-client-go"
	influxdbAPI "github.com/influxdata/influxdb-client-go/api"
	"github.com/vicanso/elite/config"
	"github.com/vicanso/elite/log"
	"go.uber.org/zap"
)

var (
	defaultInfluxSrv *InfluxSrv
)

type (
	InfluxSrv struct {
		client influxdb.Client
		writer influxdbAPI.WriteApi
	}
)

func init() {
	influxdbConfig := config.GetInfluxdbConfig()
	if influxdbConfig.Disabled {
		defaultInfluxSrv = new(InfluxSrv)
		return
	}
	opts := influxdb.DefaultOptions()
	opts.SetBatchSize(influxdbConfig.BatchSize)
	if influxdbConfig.FlushInterval > time.Millisecond {
		v := influxdbConfig.FlushInterval / time.Millisecond
		opts.SetFlushInterval(uint(v))
	}
	log.Default().Info("new influxdb client",
		zap.String("uri", influxdbConfig.URI),
		zap.String("org", influxdbConfig.Org),
		zap.String("bucket", influxdbConfig.Bucket),
		zap.Uint("batchSize", influxdbConfig.BatchSize),
		zap.String("token", influxdbConfig.Token[:5]+"..."),
		zap.Duration("interval", influxdbConfig.FlushInterval),
	)
	c := influxdb.NewClientWithOptions(influxdbConfig.URI, influxdbConfig.Token, opts)
	writer := c.WriteApi(influxdbConfig.Org, influxdbConfig.Bucket)
	defaultInfluxSrv = &InfluxSrv{
		client: c,
		writer: writer,
	}
}

// GetInfluxSrv get default influx service
func GetInfluxSrv() *InfluxSrv {
	return defaultInfluxSrv
}

// Write write metric to influxdb
func (srv *InfluxSrv) Write(measurement string, fields map[string]interface{}, tags map[string]string) {
	if srv.writer == nil {
		return
	}
	srv.writer.WritePoint(influxdb.NewPoint(measurement, tags, fields, time.Now()))
}

// Flush flush metric list
func (srv *InfluxSrv) Flush() {
	if srv.writer == nil {
		return
	}
	srv.writer.Flush()
}

// Close flush the point to influxdb and close client
func (srv *InfluxSrv) Close() {
	if srv.client == nil {
		return
	}
	srv.client.Close()
}
