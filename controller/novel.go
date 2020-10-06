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

// 小说相关的一些路由处理

package controller

import (
	"bytes"
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/vicanso/elite/cs"
	"github.com/vicanso/elite/ent"
	"github.com/vicanso/elite/ent/chapter"
	entNovel "github.com/vicanso/elite/ent/novel"
	"github.com/vicanso/elite/ent/novelsource"
	"github.com/vicanso/elite/novel"
	"github.com/vicanso/elite/router"
	"github.com/vicanso/elite/service"
	"github.com/vicanso/elite/util"
	"github.com/vicanso/elite/validate"
	"github.com/vicanso/elton"
	"github.com/vicanso/go-axios"
	"go.uber.org/zap"
)

const eliteCoverBucket = "elite-covers"

type (
	novelCtrl struct{}

	// novelListResp 小说列表响应
	novelListResp struct {
		Novels []*ent.Novel `json:"novels,omitempty"`
		Count  int          `json:"count,omitempty"`
	}
	// novelChapterListResp 小说章节列表响应
	novelChapterListResp struct {
		Chapters []*ent.Chapter `json:"chapters,omitempty"`
		Count    int            `json:"count,omitempty"`
	}
	// novelSourceListResp 小说源列表响应
	novelSourceListResp struct {
		NovelSources []*ent.NovelSource `json:"novelSources,omitempty"`
		Count        int                `json:"count,omitempty"`
	}

	// novelAddParams 添加小说参数
	novelAddParams struct {
		Name   string `json:"name,omitempty" validate:"required,xNovelName"`
		Author string `json:"author,omitempty" validate:"required,xNovelAuthor"`
	}
	// novelUpdateParams 更新小说参数
	novelUpdateParams struct {
		ID      int    `json:"id,omitempty"`
		Status  int    `json:"status,omitempty" validate:"omitempty,xNovelStatus"`
		Summary string `json:"summary,omitempty" validate:"omitempty,xNovelSummary"`
	}
	// novelListParams 小说查询参数
	novelListParams struct {
		listParams

		Keyword string `json:"keyword,omitempty" validate:"omitempty,xKeyword"`
	}
	// novelChpaterListParams 章节查询参数
	novelChpaterListParams struct {
		listParams

		ID int `json:"id,omitempty"`
	}

	// novelSourceListParams 小说源查询参数
	novelSourceListParams struct {
		listParams

		Keyword string `json:"keyword,omitempty" validate:"omitempty,xKeyword"`
	}
)

func init() {
	ctrl := novelCtrl{}

	g := router.NewGroup("/novels")

	// 获取小说源列表
	g.GET(
		"/v1/sources",
		loadUserSession,
		shouldBeAdmin,
		ctrl.listSource,
	)

	// 添加小说
	g.POST(
		"/v1",
		loadUserSession,
		shouldBeAdmin,
		ctrl.add,
	)
	// 小说查询
	g.GET(
		"/v1",
		setNoCacheIfMatched,
		ctrl.list,
	)
	// 单本小说查询
	g.GET(
		"/v1/{id}",
		setNoCacheIfMatched,
		ctrl.findByID,
	)
	// 单本小说更新
	g.PATCH(
		"/v1/{id}",
		newTracker(cs.ActionNovelUpdate),
		loadUserSession,
		shouldBeAdmin,
		ctrl.updateByID,
	)
	// 小说章节查询
	g.GET(
		"/v1/{id}/chapters",
		ctrl.listChapter,
	)
	// 小说章节内容
	g.GET(
		"/v1/{id}/chapters/{no}",
		ctrl.getChapterContent,
	)
	// 小说封面
	g.GET(
		"/v1/{id}/cover",
		ctrl.getCover,
	)
}

// where 将查询条件中的参数转换为对应的where条件
func (params *novelListParams) where(query *ent.NovelQuery) *ent.NovelQuery {
	if params.Keyword != "" {
		query = query.Where(entNovel.Or(entNovel.NameContains(params.Keyword), entNovel.AuthorContains(params.Keyword)))
	}

	return query
}

// queryAll 查询小说列表
func (params *novelListParams) queryAll(ctx context.Context) (novels []*ent.Novel, err error) {
	query := getEntClient().Novel.Query()

	query = query.Limit(params.GetLimit()).
		Offset(params.GetOffset()).
		Order(params.GetOrders()...)
	fields := params.GetFields()
	query = params.where(query)
	// 如果指定select的字段
	if len(fields) != 0 {
		err = query.Select(fields[0], fields[1:]...).Scan(ctx, &novels)
		if err != nil {
			return
		}
		return
	}
	return query.All(ctx)
}

// count 计算总数
func (params *novelListParams) count(ctx context.Context) (count int, err error) {
	query := getEntClient().Novel.Query()
	query = params.where(query)

	return query.Count(ctx)
}

// where 将查询条件转换为where
func (params *novelChpaterListParams) where(query *ent.ChapterQuery) *ent.ChapterQuery {
	query = query.Where(chapter.NovelEQ(params.ID))
	return query
}

// queryAll 查询小说章节
func (params *novelChpaterListParams) queryAll(ctx context.Context) (chapters []*ent.Chapter, err error) {
	query := getEntClient().Chapter.Query()
	query = query.Limit(params.GetLimit()).
		Offset(params.GetOffset()).
		Order(params.GetOrders()...)
	query = params.where(query)
	return query.All(ctx)
}

// count 计算章节总数
func (params *novelChpaterListParams) count(ctx context.Context) (count int, err error) {
	query := getEntClient().Chapter.Query()
	query = params.where(query)
	return query.Count(ctx)
}

// where 将查询参数转换为where条件
func (params *novelSourceListParams) where(query *ent.NovelSourceQuery) *ent.NovelSourceQuery {
	if params.Keyword != "" {
		query = query.Where(novelsource.Or(novelsource.NameContains(params.Keyword), novelsource.AuthorContains(params.Keyword)))
	}
	return query
}

// queryAll 查询符合条件的记录
func (params *novelSourceListParams) queryAll(ctx context.Context) (sources []*ent.NovelSource, err error) {
	query := getEntClient().NovelSource.Query()
	query = query.Limit(params.GetLimit()).
		Offset(params.GetOffset()).
		Order(params.GetOrders()...)
	query = params.where(query)

	return query.All(ctx)
}

// count 计算小说源总数
func (params *novelSourceListParams) count(ctx context.Context) (count int, err error) {
	query := getEntClient().NovelSource.Query()
	query = params.where(query)
	return query.Count(ctx)
}

// update 更新小说
func (params *novelUpdateParams) update(ctx context.Context) (count int, err error) {
	if params.ID == 0 {
		err = errors.New("id can't be nil")
		return
	}
	update := getEntClient().Novel.Update().
		Where(entNovel.IDEQ(params.ID))
	if params.Summary != "" {
		update = update.SetSummary(params.Summary)
	}
	if params.Status != 0 {
		update = update.SetStatus(params.Status)
	}
	count, err = update.Save(ctx)
	if err != nil {
		return
	}
	if count == 0 {
		err = errors.New("no record match")
		return
	}
	return
}

// list 查询小说列表
func (*novelCtrl) list(c *elton.Context) (err error) {
	params := novelListParams{}
	err = validate.Do(&params, c.Query())
	if err != nil {
		return
	}
	count := -1
	if params.GetOffset() == 0 {
		count, err = params.count(c.Context())
		if err != nil {
			return
		}
	}
	novels, err := params.queryAll(c.Context())
	if err != nil {
		return
	}
	c.CacheMaxAge("5m")
	c.Body = &novelListResp{
		Novels: novels,
		Count:  count,
	}
	return
}

// findByID 通过id查询书籍
func (*novelCtrl) findByID(c *elton.Context) (err error) {
	id, err := getIDFromParams(c)
	if err != nil {
		return
	}
	result, err := getEntClient().Novel.Query().
		Where(entNovel.IDEQ(id)).
		First(c.Context())
	if err != nil {
		return
	}
	c.CacheMaxAge("10m")
	c.Body = result
	return
}

// updateByID 根据id更新小说信息
func (*novelCtrl) updateByID(c *elton.Context) (err error) {
	id, err := getIDFromParams(c)
	if err != nil {
		return
	}
	params := novelUpdateParams{}
	err = validate.Do(&params, c.RequestBody)
	if err != nil {
		return
	}
	params.ID = id
	_, err = params.update(c.Context())
	if err != nil {
		return
	}
	c.NoContent()
	return
}

// add 添加小说
func (*novelCtrl) add(c *elton.Context) (err error) {
	params := novelAddParams{}
	err = validate.Do(&params, c.RequestBody)
	if err != nil {
		return
	}
	queryParmas := novel.QueryParams{
		Name:   params.Name,
		Author: params.Author,
	}
	result, err := novelSrv.Publish(queryParmas)
	if err != nil {
		return
	}
	// 更新封面
	go func() {
		// 如果是绝对地址（外网），则下载图片并保存
		if strings.HasPrefix(result.Cover, "http") {
			resp, err := axios.Get(result.Cover)
			if err != nil {
				logger.Error("get cover fail",
					zap.String("name", params.Name),
				)
				return
			}
			contentType := resp.Headers.Get("Content-Type")
			fileType := strings.Split(contentType, "/")[1]
			name := util.GenUlid() + "." + fileType
			_, err = fileSrv.Upload(context.Background(), service.UploadParams{
				Bucket: eliteCoverBucket,
				Name:   name,
				Reader: bytes.NewReader(resp.Data),
				Size:   int64(len(resp.Data)),
				Opts: minio.PutObjectOptions{
					ContentType: contentType,
				},
			})
			if err != nil {
				logger.Error("upload cover fail",
					zap.String("name", params.Name),
				)
				return
			}
			_, err = result.Update().
				SetCover(name).Save(context.Background())
			if err != nil {
				logger.Error("update cover fail",
					zap.String("name", params.Name),
				)
				return
			}

		}
	}()

	c.Created(result)
	return
}

// listChapter 获取小说章节
func (*novelCtrl) listChapter(c *elton.Context) (err error) {
	id, err := getIDFromParams(c)
	if err != nil {
		return
	}
	params := novelChpaterListParams{}
	err = validate.Do(&params, c.Query())
	if err != nil {
		return
	}
	params.ID = id
	count := -1
	if params.GetOffset() == 0 {
		count, err = params.count(c.Context())
		if err != nil {
			return
		}
	}
	// 如果章节总数为0，则fetch数据
	if count == 0 {
		err = novelSrv.UpdateChapters(id)
		if err != nil {
			return
		}
	}
	chapters, err := params.queryAll(c.Context())
	if err != nil {
		return
	}
	c.CacheMaxAge("5m")
	c.Body = &novelChapterListResp{
		Count:    count,
		Chapters: chapters,
	}
	return
}

// getChapterContent 获取章节内容
func (*novelCtrl) getChapterContent(c *elton.Context) (err error) {
	id, err := getIDFromParams(c)
	if err != nil {
		return
	}
	no, err := strconv.Atoi(c.Param("no"))
	if err != nil {
		return
	}
	result, err := novelSrv.GetChapterContent(id, no)
	if err != nil {
		return
	}
	c.CacheMaxAge("10m")
	c.Body = result
	return
}

// getCover 获取小说封面
func (*novelCtrl) getCover(c *elton.Context) (err error) {
	id, err := getIDFromParams(c)
	if err != nil {
		return
	}
	cover, err := novelSrv.GetCover(id)
	if err != nil {
		return
	}
	// TODO 后续优化cover压缩，支持webp
	data, header, err := fileSrv.GetData(c.Context(), eliteCoverBucket, cover)
	if err != nil {
		return
	}
	c.MergeHeader(header)
	c.Body = data
	return
}

// listSource 获取小说源列表
func (*novelCtrl) listSource(c *elton.Context) (err error) {
	params := novelSourceListParams{}
	err = validate.Do(&params, c.Query())
	if err != nil {
		return
	}
	count := -1
	if params.GetOffset() == 0 {
		count, err = params.count(c.Context())
		if err != nil {
			return
		}
	}
	novelSources, err := params.queryAll(c.Context())
	if err != nil {
		return
	}
	c.Body = &novelSourceListResp{
		NovelSources: novelSources,
		Count:        count,
	}
	return
}
