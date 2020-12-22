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
	"github.com/vicanso/elite/config"
	"github.com/vicanso/go-axios"
)

const (
	novelBiQuGeName = "biquge"
	novelQiDianName = "qidian"
)

var locationIns = newLocationInstance()
var biqugeIns = newNovelInstance(novelBiQuGeName)
var qidianIns = newNovelInstance(novelQiDianName)

// newLocationInstance 初始化location的实例
func newLocationInstance() *axios.Instance {
	locationConfig := config.GetLocationConfig()
	return NewHTTPInstance(locationConfig.Name, locationConfig.BaseURL, locationConfig.Timeout)
}

func getNovelConfig(name string) config.NovelConfig {
	conf := config.NovelConfig{}
	novelConfigs := config.GetNovelConfigs()
	for _, item := range novelConfigs {
		if item.Name == name {
			conf = item
		}
	}
	return conf
}

func newNovelInstance(name string) *axios.Instance {
	conf := getNovelConfig(name)
	return NewHTTPInstance(conf.Name, conf.BaseURL, conf.Timeout)
}

// GetLocationInstance get location instance
func GetLocationInstance() *axios.Instance {
	return locationIns
}

// GetBiqugeInstance get biquge instance
func GetBiqugeInstance() *axios.Instance {
	return biqugeIns
}

// GetQidianInstance get qidian instance
func GetQidianInstance() *axios.Instance {
	return qidianIns
}

// GetHTTPInstanceStats get http instance stats
func GetHTTPInstanceStats() map[string]interface{} {
	return map[string]interface{}{
		"location": 0,
		"biquge":   0,
		"qidian":   0,
	}
}
