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

package request

import (
	"github.com/vicanso/elite/config"
	"github.com/vicanso/go-axios"
)

var locationIns = newLocation()
var qiDianIns = newQiDian()
var biQuGeIns = newBiQuGe()
var locationService = "location"
var qiDianService = "qidian"
var biQuGeService = "biquge"
var insList = map[string]*axios.Instance{
	locationService: locationIns,
	qiDianService:   qiDianIns,
	biQuGeService:   biQuGeIns,
}

type InstanceStats struct {
	Name           string `json:"name,omitempty"`
	MaxConcurrency int    `json:"maxConcurrency,omitempty"`
	Concurrency    int    `json:"concurrency,omitempty"`
}

// newLocation 初始化location的实例
func newLocation() *axios.Instance {
	locationConfig := config.GetLocationConfig()
	return NewHTTP(locationService, locationConfig.BaseURL, locationConfig.Timeout)
}

func newQiDian() *axios.Instance {
	conf := config.GetNovelConfigs().Find(qiDianService)
	return NewHTTP(qiDianService, conf.BaseURL, conf.Timeout)
}
func newBiQuGe() *axios.Instance {
	conf := config.GetNovelConfigs().Find(biQuGeService)
	return NewHTTP(biQuGeService, conf.BaseURL, conf.Timeout)
}

// GetLocation get location instance
func GetLocation() *axios.Instance {
	return locationIns
}

// GetQiDian get qi dian instance
func GetQiDian() *axios.Instance {
	return qiDianIns
}

// GetBiQuGe get bi qu ge instance
func GetBiQuGe() *axios.Instance {
	return biQuGeIns
}

// GetHTTPStats get http instance stats
func GetHTTPStats() []*InstanceStats {
	statsList := make([]*InstanceStats, len(insList))
	index := 0
	for name, ins := range insList {
		stats := InstanceStats{
			Name:           name,
			MaxConcurrency: int(ins.Config.MaxConcurrency),
			Concurrency:    int(ins.GetConcurrency()),
		}
		statsList[index] = &stats
		index++
	}
	return statsList
}

// UpdateConcurrencyLimit update the concurrency limit for instance
func UpdateConcurrencyLimit(limits map[string]int) {
	for name, ins := range insList {
		v := limits[name]
		limit := int32(v)
		if ins.Config.MaxConcurrency != limit {
			ins.SetMaxConcurrency(limit)
		}
	}
}
