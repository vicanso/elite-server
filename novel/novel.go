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
	"time"

	"github.com/vicanso/elite/config"
	"github.com/vicanso/elite/ent"
	"github.com/vicanso/elite/helper"
	"golang.org/x/net/context"
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

type (
	// Novel 小说
	Novel struct {
		Name        string
		Author      string
		Description string
		SourceID    int
		Source      int
	}
	// Chapter 小说章节
	Chapter struct {
		Title string
		NO    int
		URL   string
	}
)

// AddToSource 添加至小说源
func (novel *Novel) AddToSource() (source *ent.NovelSource, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	source, err = getEntClient().NovelSource.Create().
		SetName(novel.Name).
		SetAuthor(novel.Author).
		SetSource(novel.Source).
		SetSourceID(novel.SourceID).
		SetDescription(novel.Description).
		Save(ctx)
	if err != nil {
		return
	}
	return
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
func SyncSource() (err error) {
	// NewBiQuGe().GetChapterContent("/book/15517/3997542.html")
	// fmt.Println(NewBiQuGe().GetDetail(8349))
	// return
	redisSrv := new(helper.Redis)
	// 确保只有一个实例在更新
	ok, done, err := redisSrv.LockWithDone("novel-sync-source", time.Hour)
	if err != nil || !ok {
		return
	}
	defer func() {
		_ = done()
	}()
	biQuGe := NewBiQuGe()
	err = biQuGe.Sync()
	if err != nil {
		return
	}
	return
}
