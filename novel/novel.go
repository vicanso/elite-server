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
	"context"
	"time"

	"github.com/vicanso/elite/cache"
	"github.com/vicanso/elite/ent"
	"github.com/vicanso/elite/ent/chapter"
	"github.com/vicanso/elite/ent/novel"
	"github.com/vicanso/elite/ent/novelsource"
	entSchema "github.com/vicanso/elite/ent/schema"
	"github.com/vicanso/elite/helper"
	"github.com/vicanso/elite/log"
	"github.com/vicanso/hes"
	"go.uber.org/zap"
)

var (
	getEntClient = helper.EntGetClient

	logger = log.Default()
)

const (
	errNovelCategory = "novel"

	defaultQueryTimeout = 3 * time.Second
)

// 小说来源
const (
	// NovelSourceBiQuGe biquge source
	NovelSourceBiQuGe = iota + 1
	// NovelSourceQiDian qidian source
	NovelSourceQiDian
)

type (
	Srv struct{}
	// Fetcher 小说拉取的interface
	Fetcher interface {
		GetDetail() (novel Novel, err error)
		GetChapters() (chapters []*Chapter, err error)
		GetChapterContent(no int) (content string, err error)
	}
	// Novel 小说
	Novel struct {
		Name     string
		Author   string
		Summary  string
		SourceID int
		Source   int
		CoverURL string
	}
	// Chapter 小说章节
	Chapter struct {
		Title string
		NO    int
		URL   string
	}
	// QueryParams 查询参数
	QueryParams struct {
		Name   string
		Author string
		Source int
	}
)

// AddToSource 添加至小说源
func (novel *Novel) AddToSource() (source *ent.NovelSource, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultQueryTimeout)
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
	ctx, cancel := context.WithTimeout(context.Background(), defaultQueryTimeout)
	defer cancel()
	result, err = getEntClient().Novel.Create().
		SetName(novel.Name).
		SetAuthor(novel.Author).
		SetSource(novel.Source).
		SetSummary(novel.Summary).
		SetCover(novel.CoverURL).
		Save(ctx)
	return
}

// FirstNovel 查询第一条符合条件的小说
func (params *QueryParams) FirstNovel() (*ent.Novel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultQueryTimeout)
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
	ctx, cancel := context.WithTimeout(context.Background(), defaultQueryTimeout)
	defer cancel()
	query := getEntClient().NovelSource.Query()
	if params.Name != "" {
		query = query.Where(novelsource.NameEQ(params.Name))
	}
	if params.Author != "" {
		query = query.Where(novelsource.AuthorEQ(params.Author))
	}
	if params.Source != 0 {
		query = query.Where(novelsource.SourceEQ(params.Source))
	}
	query = query.Order(ent.Asc("source"))
	return query.First(ctx)
}

// SyncSource 同步小说
func (*Srv) SyncSource() (err error) {
	redisSrv := new(cache.Redis)
	// 确保只有一个实例在更新
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	ok, done, err := redisSrv.LockWithDone(ctx, "novel-sync-source", time.Hour)
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

// GetFetcherByID 根据小说id获取其fetcher
func (srv *Srv) GetFetcherByID(id int) (fetcher Fetcher, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultQueryTimeout)
	defer cancel()
	result, err := getEntClient().Novel.Query().
		Where(novel.ID(id)).
		First(ctx)
	if err != nil {
		return
	}
	return srv.GetFetcher(QueryParams{
		Author: result.Author,
		Name:   result.Name,
		Source: result.Source,
	})
}

// GetFetcher 获取fetcher
func (*Srv) GetFetcher(params QueryParams) (fetcher Fetcher, err error) {
	novelSource, err := params.FirstNovelSource()
	if ent.IsNotFound(err) {
		err = nil
	}
	if err != nil {
		return
	}
	if novelSource == nil {
		err = hes.New("无法找到该小说的源", errNovelCategory)
		return
	}
	// TODO 后续添更多的源
	fetcher = NewBiQuGe().NewFetcher(novelSource.SourceID)
	return
}

// Publish 发布小说
func (srv *Srv) Publish(params QueryParams) (novel *ent.Novel, err error) {

	novel, err = params.FirstNovel()
	if ent.IsNotFound(err) {
		err = nil
	}
	if err != nil || novel != nil {
		return
	}
	fetcher, err := srv.GetFetcher(params)
	if err != nil {
		return
	}
	result, err := fetcher.GetDetail()
	if err != nil {
		return
	}
	// 如果从qidian中能获取，则替换简介
	qiDianResult, _ := NewQiDian().Search(params.Name, params.Author)
	if qiDianResult.SourceID != 0 {
		result.Summary = qiDianResult.Summary
		if qiDianResult.CoverURL != "" {
			result.CoverURL = qiDianResult.CoverURL
		}
	}

	// 添加小说
	novel, err = result.Add()
	if err != nil {
		return
	}
	// 更新小说来源为已发布
	go func() {
		_, err := getEntClient().NovelSource.Update().
			Where(novelsource.NameEQ(params.Name)).
			Where(novelsource.AuthorEQ(params.Author)).
			SetStatus(entSchema.NovelSourceStatusPublished).
			Save(context.Background())
		if err != nil {
			logger.Error("update novel source status fail",
				zap.String("name", params.Name),
				zap.String("author", params.Author),
				zap.Error(err),
			)
		}
	}()
	return
}

// UpdateChapters 拉取小说章节
func (srv *Srv) UpdateChapters(id int) (err error) {
	fetcher, err := srv.GetFetcherByID(id)
	if err != nil {
		return
	}
	// TODO 支持更多的小说源
	chapters, err := fetcher.GetChapters()
	if err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), defaultQueryTimeout)
	defer cancel()
	currentCount, err := getEntClient().Chapter.Query().
		Where(chapter.NovelEQ(id)).
		Count(ctx)
	if err != nil {
		return
	}
	// 如果所有章节都已更新
	if len(chapters) <= currentCount {
		return
	}
	chapters = chapters[currentCount:]
	bulk := make([]*ent.ChapterCreate, len(chapters))
	for i, item := range chapters {
		bulk[i] = getEntClient().Chapter.Create().
			SetTitle(item.Title).
			SetNo(item.NO).
			SetNovel(id)
	}
	_, err = getEntClient().Chapter.CreateBulk(bulk...).Save(ctx)
	if err != nil {
		return
	}
	return
}

// UpdateAllChaptersByWeight 根据更新权重更新所有小说章节
func (srv *Srv) UpdateAllChaptersByWeight(minUpdatedWeight int) (err error) {
	maxID, err := srv.GetMaxID()
	if err != nil {
		return
	}
	var item *ent.Novel
	for i := 0; i < maxID; i++ {
		// 小说id从1开始
		id := i + 1
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
		defer cancel()
		item, err = getEntClient().Novel.Get(ctx, id)
		if ent.IsNotFound(err) {
			err = nil
		}
		if err != nil {
			return
		}
		if item == nil {
			continue
		}
		if item.UpdatedWeight >= minUpdatedWeight {
			err = srv.UpdateChapters(id)
			if ent.IsNotFound(err) {
				err = nil
			}
			if err != nil {
				return
			}
		}
	}
	return
}

// UpdateWordCount 更新总字数
func (srv *Srv) UpdateWordCount(id int) (err error) {
	chapters := make([]*ent.Chapter, 0)
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, defaultQueryTimeout)
	defer cancel()
	err = getEntClient().Chapter.Query().
		Where(chapter.Novel(id)).
		Where(chapter.WordCountNotNil()).
		Select(chapter.FieldWordCount).
		Scan(ctx, &chapters)
	if err != nil {
		return
	}
	wordCount := 0
	for _, item := range chapters {
		wordCount += item.WordCount
	}
	_, err = getEntClient().Novel.UpdateOneID(id).
		SetWordCount(wordCount).
		Save(context.Background())
	if err != nil {
		return
	}
	return
}

// GetMaxID 获取最大的小说ID
func (srv *Srv) GetMaxID() (id int, err error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	return getEntClient().Novel.Query().
		Order(ent.Desc("id")).
		FirstID(ctx)
}

func (srv *Srv) doAll(fn func(int) error) (err error) {
	maxID, err := srv.GetMaxID()
	if err != nil {
		return
	}
	for i := 0; i < maxID; i++ {
		// 小说id从1开始
		err = fn(i + 1)
		if ent.IsNotFound(err) {
			err = nil
		}
		if err != nil {
			return
		}
	}
	return
}

// UpdateAllWordCount 更新所有小说总字数
func (srv *Srv) UpdateAllWordCount() (err error) {
	return srv.doAll(srv.UpdateWordCount)
}

// UpdateUpdatedWeight 更新小说权重
func (srv *Srv) UpdateUpdatedWeight(id int) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*defaultQueryTimeout)
	defer cancel()

	chapters := make([]*ent.Chapter, 0)
	// 查询最新更新的10章
	err = getEntClient().Chapter.Query().
		Where(chapter.Novel(id)).
		Order(ent.Desc("no")).
		Limit(10).
		Select("updated_at").
		Scan(ctx, &chapters)
	if err != nil {
		return
	}
	current := time.Now()
	oneDay := 24 * time.Hour
	oneWeek := 7 * oneDay
	oneMonth := 30 * oneDay
	updatedWeight := 0
	for _, item := range chapters {
		subTime := current.Sub(item.UpdatedAt)
		if subTime < oneDay {
			updatedWeight += 10
		} else if subTime < oneWeek {
			updatedWeight += 2
		} else if subTime < oneMonth {
			updatedWeight++
		}
	}
	_, err = getEntClient().Novel.UpdateOneID(id).
		SetUpdatedWeight(updatedWeight).
		Save(ctx)
	if err != nil {
		return
	}
	return
}

// UpdateAllUpdatedWeight 更新所有小说更新权重
func (srv *Srv) UpdateAllUpdatedWeight() (err error) {
	return srv.doAll(srv.UpdateUpdatedWeight)
}

// FetchAllChapterContent 拉取所有章节内容
func (srv *Srv) FetchAllChapterContent(id int) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*defaultQueryTimeout)
	defer cancel()
	count, err := getEntClient().Chapter.Query().
		Where(chapter.NovelEQ(id)).
		Count(ctx)
	if err != nil {
		return
	}
	for i := 0; i < count; i++ {
		_, err = srv.GetChapterContent(id, i)
		if err != nil {
			return
		}
	}
	return
}

// GetChapterContent 获取小说章节内容
func (srv *Srv) GetChapterContent(id, no int) (result *ent.Chapter, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*defaultQueryTimeout)
	defer cancel()
	result, err = getEntClient().Chapter.Query().
		Where(chapter.NovelEQ(id)).
		Where(chapter.NoEQ(no)).
		First(ctx)
	if err != nil {
		return
	}
	if result.Content != "" {
		return
	}
	fetcher, err := srv.GetFetcherByID(id)
	if err != nil {
		return
	}
	content, err := fetcher.GetChapterContent(no)
	if err != nil {
		return
	}
	result, err = result.Update().
		SetContent(content).
		SetWordCount(len(content)).
		Save(ctx)
	if err != nil {
		return
	}
	return
}

// GetCover 获取小说封面
func (*Srv) GetCover(id int) (cover string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultQueryTimeout)
	defer cancel()
	result, err := getEntClient().Novel.Query().
		Where(novel.ID(id)).
		First(ctx)
	if err != nil {
		return
	}
	cover = result.Cover
	return
}
