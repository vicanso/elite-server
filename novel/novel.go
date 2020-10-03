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

// 主要包括各类小说的抓取功能

package novel

import (
	"fmt"

	"github.com/vicanso/elite/config"
	"github.com/vicanso/elite/helper"
)

var novelConfigs = config.GetNovelConfigs()

var (
	getEntClient = helper.EntGetClient
)

const (
	novelBiQuGeName = "biquge"
)

// 小说来源
const (
	// NovelSourceBiQuGe biquge source
	NovelSourceBiQuGe = iota + 1
)

type Novel struct {
	Name     string
	Author   string
	SourceID int
	Source   int
}

// getNovelConfig 获取对应的novel配置
func getNovelConfig(name string) (conf config.NovelConfig) {
	for _, item := range novelConfigs {
		if item.Name == name {
			conf = item
		}
	}
	return
}

// SyncSource 同步小说
func SyncSource() {
	biQuGe := NewBiQuGe()
	fmt.Println(biQuGe.GetDetail(1))
}
