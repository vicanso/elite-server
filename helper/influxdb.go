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
	"sync"
	"time"

	influxdb "github.com/influxdata/influxdb-client-go"
	"github.com/vicanso/elite/config"
)

var (
	influxdbClient   influxdb.InfluxDBClient
	defaultInfluxSrv *InfluxSrv

	initDefaultInfluxSrv sync.Once
)

type (
	InfluxSrv struct {
		writer influxdb.WriteApi
	}
)

func init() {
	influxbConfig := config.GetInfluxdbConfig()
	opts := influxdb.DefaultOptions()
	opts.SetBatchSize(influxbConfig.BatchSize)
	influxdbClient = influxdb.NewClientWithOptions(influxbConfig.URI, influxbConfig.Token, opts)
}

// GetInfluxSrv get default influx service
func GetInfluxSrv() *InfluxSrv {
	initDefaultInfluxSrv.Do(func() {
		influxbConfig := config.GetInfluxdbConfig()
		defaultInfluxSrv = &InfluxSrv{
			writer: influxdbClient.WriteApi(influxbConfig.Org, influxbConfig.Bucket),
		}
		// defaultInfluxSrv = &InfluxSrv{
		// 	BatchSize: influxbConfig.BatchSize,
		// 	Bucket:    influxbConfig.Bucket,
		// 	Org:       influxbConfig.Org,
		// }
	})
	return defaultInfluxSrv
}

// Write write metric to influxdb
func (srv *InfluxSrv) Write(measurement string, fields map[string]interface{}, tags map[string]string) {
	defaultInfluxSrv.writer.WritePoint(influxdb.NewPoint(measurement, tags, fields, time.Now()))
}
