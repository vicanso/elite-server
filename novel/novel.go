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
	"github.com/vicanso/elite/ent/novel"
	"github.com/vicanso/elite/ent/novelsource"
	"github.com/vicanso/elite/helper"
	"github.com/vicanso/hes"
	"golang.org/x/net/context"
)

var novelConfigs = config.GetNovelConfigs()

var (
	getEntClient = helper.EntGetClient
)

const (
	errNovelCategory = "novel"
)

var (
	errNovelSourceNotFound = &hes.Error{
		Message:    "无法找到该小说的源",
		StatusCode: 400,
		Category:   errNovelCategory,
	}
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
		Name     string
		Author   string
		Summary  string
		SourceID int
		Source   int
	}
	// Chapter 小说章节
	Chapter struct {
		Title string
		NO    int
		URL   string
	}
	QueryParams struct {
		Name   string
		Author string
	}
)

// getConfig 获取对应的novel配置
func getConfig(name string) (conf config.NovelConfig) {
	for _, item := range novelConfigs {
		if item.Name == name {
			conf = item
		}
	}
	return
}

// AddToSource 添加至小说源
func (novel *Novel) AddToSource() (source *ent.NovelSource, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	source, err = getEntClient().NovelSource.Create().
		SetName(novel.Name).
		SetAuthor(novel.Author).
		SetSource(novel.Source).
		SetSourceID(novel.SourceID).
		Save(ctx)
	if err != nil {
		return
	}
	return
}

// Add 添加小说
func (novel *Novel) Add() (result *ent.Novel, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err = getEntClient().Novel.Create().
		SetName(novel.Name).
		SetAuthor(novel.Author).
		SetSource(novel.Source).
		SetSummary(novel.Summary).
		Save(ctx)
	return
}

// FirstNovel 查询第一条符合条件的小说
func (params *QueryParams) FirstNovel() (*ent.Novel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := getEntClient().Novel.Query()
	if params.Name != "" {
		query = query.Where(novel.NameEQ(params.Name))
	}
	if params.Author != "" {
		query = query.Where(novel.AuthorEQ(params.Author))
	}
	return query.First(ctx)
}

// FirstNovelSOurce 获取第一个符合的小说源
func (params *QueryParams) FirstNovelSource() (*ent.NovelSource, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := getEntClient().NovelSource.Query()
	if params.Name != "" {
		query = query.Where(novelsource.NameEQ(params.Name))
	}
	if params.Author != "" {
		query = query.Where(novelsource.AuthorEQ(params.Author))
	}
	query = query.Order(ent.Asc("source"))
	return query.First(ctx)
}

// SyncSource 同步小说
func SyncSource() (err error) {
	redisSrv := new(helper.Redis)
	// 确保只有一个实例在更新
	ok, done, err := redisSrv.LockWithDone("novel-sync-source", time.Hour)
	if err != nil || !ok {
		return
	}
	defer func() {
		_ = done()
	}()
	err = NewBiQuGe().Sync()
	if err != nil {
		return
	}
	return
}

// Publish 发布小说
func Publish(params QueryParams) (novel *ent.Novel, err error) {
	novel, err = params.FirstNovel()
	if ent.IsNotFound(err) {
		err = nil
	}
	if err != nil || novel != nil {
		return
	}
	novelSource, err := params.FirstNovelSource()
	if ent.IsNotFound(err) {
		err = nil
	}
	if err != nil {
		return
	}
	if novelSource == nil {
		err = errNovelSourceNotFound
		return
	}
	// TODO 支持更多的小说源
	result, err := NewBiQuGe().GetDetail(novelSource.SourceID)
	if err != nil {
		return
	}

	// 添加小说
	novel, err = result.Add()
	if err != nil {
		return
	}

	return
}
