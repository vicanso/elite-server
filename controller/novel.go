// Copyright 2021 tree xie
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
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/vicanso/elite/cs"
	"github.com/vicanso/elite/ent"
	"github.com/vicanso/elite/ent/chapter"
	entNovel "github.com/vicanso/elite/ent/novel"
	"github.com/vicanso/elite/log"
	"github.com/vicanso/elite/novel"
	"github.com/vicanso/elite/router"
	"github.com/vicanso/elite/schema"
	"github.com/vicanso/elite/service"
	"github.com/vicanso/elite/util"
	"github.com/vicanso/elite/validate"
	"github.com/vicanso/elton"
	"github.com/vicanso/go-axios"
	"github.com/vicanso/hes"
)

type novelCtrl struct{}

const eliteCoverBucket = "elite-covers"
const errNovelCategory = "novel"

// 接口参数定义
type (
	// novelListParams 小说查询参数
	novelListParams struct {
		listParams

		Keyword       string `json:"keyword,omitempty" validate:"omitempty,xKeyword"`
		IDS           string `json:"ids,omitempty" validate:"omitempty,xNovelIDS"`
		AuthorKeyword string
		NameKeyword   string
	}
	// novelUpdateParams 更新小说参数
	novelUpdateParams struct {
		ID      int    `json:"id,omitempty"`
		Status  int    `json:"status,omitempty" validate:"omitempty,xNovelStatus"`
		Summary string `json:"summary,omitempty" validate:"omitempty,xNovelSummary"`
	}
	// novelChapterListParams 章节查询参数
	novelChapterListParams struct {
		listParams

		// ID 小说id，由route param中获取并设置，因此不设置validate
		ID int `json:"id,omitempty"`
		// ChapterID 章节id，由route param中获取并设置，因此不设置validate
		ChapterID int `json:"chapterID,omitempty"`
	}
	// novelChapterUpdateParams 章节更新参数
	novelChapterUpdateParams struct {
		Content string `json:"content,omitempty" validate:"required"`
	}
	// novelCoverParams 小说封面参数
	novelCoverParams struct {
		Type    string `json:"type,omitempty" validate:"required,xNovelCoverType"`
		Width   string `json:"width,omitempty" validate:"omitempty,xNovelCoverWidth"`
		Height  string `json:"height,omitempty" validate:"omitempty,xNovelCoverHeight"`
		Quality string `json:"quality,omitempty" validate:"required,xNovelCoverQuality"`
	}
)

// 接口响应定义
type (
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
	// novelHotKeywordListResp 热门搜索关键字列表响应
	novelHotKeywordListResp struct {
		Keywords []string `json:"keywords"`
	}
)

func init() {
	g := router.NewGroup("/novels")

	ctrl := novelCtrl{}

	// 小说查询
	g.GET(
		"/v1",
		ctrl.list,
	)
	// 单本小说查询
	g.GET(
		"/v1/{id}",
		ctrl.findByID,
	)
	// 单本小说更新
	g.PATCH(
		"/v1/{id}",
		newTrackerMiddleware(cs.ActionNovelUpdate),
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
		ctrl.getChapterDetail,
	)
	// 小说章节内容
	g.PATCH(
		"/v1/{id}/chapters/{no}",
		ctrl.updateChapterDetail,
	)
	// 小说封面
	g.GET(
		"/v1/{id}/cover",
		ctrl.getCover,
	)

	// 更新所有章节
	g.POST(
		"/v1/update-all-chapters",
		loadUserSession,
		shouldBeAdmin,
		ctrl.updateAllChapters,
	)

	// 发布所有的小说
	g.POST(
		"/v1/publish-all",
		loadUserSession,
		shouldBeAdmin,
		ctrl.publishAll,
	)

	g.GET(
		"/v1/hot-keywords",
		ctrl.listHotKeyword,
	)
}

// where 将查询条件中的参数转换为对应的where条件
func (params *novelListParams) where(query *ent.NovelQuery) *ent.NovelQuery {
	// 通过keyword转换而来
	if params.NameKeyword != "" {
		query = query.Where(entNovel.NameContains(params.NameKeyword))
	}
	// 通过keyword转换而来
	if params.AuthorKeyword != "" {
		query = query.Where(entNovel.AuthorContains(params.AuthorKeyword))
	}

	if params.IDS != "" {
		arr := strings.Split(params.IDS, ",")
		ids := make([]int, 0, len(arr))
		for _, v := range arr {
			id, _ := strconv.Atoi(v)
			if id != 0 {
				ids = append(ids, id)
			}
		}
		query = query.Where(entNovel.IDIn(ids...))
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
		novels = make([]*ent.Novel, 0)
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
func (params *novelChapterListParams) where(query *ent.ChapterQuery) *ent.ChapterQuery {
	if params.ID != 0 {
		query = query.Where(chapter.NovelEQ(params.ID))
	}
	if params.ChapterID != 0 {
		query = query.Where(chapter.ID(params.ChapterID))
	}
	return query
}

// queryAll 查询小说章节
func (params *novelChapterListParams) queryAll(ctx context.Context) (chapters []*ent.Chapter, err error) {
	if params.ID == 0 {
		return nil, hes.NewWithStatusCode("小说ID不能为空", 400)
	}
	query := getEntClient().Chapter.Query()
	query = query.Limit(params.GetLimit()).
		Offset(params.GetOffset()).
		Order(params.GetOrders()...)
	query = params.where(query)
	fields := params.GetFields()
	// 如果指定了select的字段
	if len(fields) != 0 {
		chapters = make([]*ent.Chapter, 0)
		err = query.Select(fields[0], fields[1:]...).Scan(ctx, &chapters)
		if err != nil {
			return
		}
		return
	}
	return query.All(ctx)
}

// count 计算章节总数
func (params *novelChapterListParams) count(ctx context.Context) (count int, err error) {
	if params.ID == 0 {
		return -1, hes.NewWithStatusCode("小说ID不能为空", 400)
	}
	query := getEntClient().Chapter.Query()
	query = params.where(query)
	return query.Count(ctx)
}

// update 更新小说
func (params *novelUpdateParams) update(ctx context.Context) (err error) {
	if params.ID == 0 {
		err = hes.New("小说ID不能为空", errNovelCategory)
		return
	}
	update := getEntClient().Novel.UpdateOneID(params.ID)
	if params.Summary != "" {
		update = update.SetSummary(params.Summary)
	}
	if params.Status != 0 {
		update = update.SetStatus(params.Status)
	}
	result, err := update.Save(ctx)
	if err != nil {
		return
	}
	if result == nil {
		err = hes.New("无匹配的小说记录", errNovelCategory)
		return
	}
	return
}

// getChapterDetail 获取章节内容
func (*novelCtrl) getChapterDetail(c *elton.Context) (err error) {
	id, err := getIDFromParams(c)
	if err != nil {
		return
	}
	no, err := strconv.Atoi(c.Param("no"))
	if err != nil {
		return
	}
	result, err := novelSrv.GetChapterDetail(id, no)
	if err != nil {
		return
	}
	c.CacheMaxAge(10 * time.Minute)
	c.Body = result
	return
}

// updateChapterContent 更新章节内容
func (*novelCtrl) updateChapterDetail(c *elton.Context) (err error) {
	id, err := getIDFromParams(c)
	if err != nil {
		return
	}
	no, err := strconv.Atoi(c.Param("no"))
	if err != nil {
		return
	}
	params := novelChapterUpdateParams{}
	err = validate.Do(&params, c.RequestBody)
	if err != nil {
		return
	}
	_, err = novelSrv.UpdateChapterContent(id, no, params.Content)
	if err != nil {
		return
	}
	c.NoContent()

	return
}

func updateCoverByURL(id int, coverURL string) (err error) {
	resp, err := axios.Get(coverURL)
	if err != nil {
		return
	}
	contentType := resp.Headers.Get("Content-Type")
	fileType := strings.Split(contentType, "/")[1]
	if fileType == "html" {
		err = errors.New("content type is invalid")
		return
	}
	name := util.GenXID() + "." + fileType
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
		return
	}
	_, err = getEntClient().Novel.UpdateOneID(id).
		SetCover(name).Save(context.Background())
	if err != nil {
		return
	}
	return
}

func (*novelCtrl) publish(params novel.QueryParams) (result *ent.Novel, err error) {
	result, err = novelSrv.Publish(params)
	if err != nil {
		return
	}
	// 更新封面
	go func() {
		// 如果是绝对地址（外网），则下载图片并保存
		if strings.HasPrefix(result.Cover, "http") {
			err := updateCoverByURL(result.ID, result.Cover)
			if err != nil {
				log.Default().Error().
					Err(err).
					Str("name", result.Name).
					Msg("update cover fail")
			}
		}
	}()
	return
}

// publishAll 将所有小说源的小说发布
func (ctrl *novelCtrl) publishAll(c *elton.Context) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	id, err := getEntClient().NovelSource.Query().
		Order(ent.Desc("id")).
		FirstID(ctx)
	if err != nil {
		return
	}
	go func() {
		startedAt := time.Now()
		for i := 1; i < id; i++ {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			result, _ := getEntClient().NovelSource.Get(ctx, i)
			if result == nil || result.Status == schema.NovelSourceStatusPublished {
				continue
			}

			_, err := ctrl.publish(novel.QueryParams{
				Name:   result.Name,
				Author: result.Author,
			})
			if err != nil {
				log.Default().Error().
					Str("name", result.Name).
					Str("author", result.Author).
					Err(err).
					Msg("publish novel fail")
			}
		}
		log.Default().Info().
			Dur("use", time.Since(startedAt)).
			Msg("publish all done")
	}()
	c.NoContent()
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
	var novels []*ent.Novel
	// 如果有关键字，则不计算总数
	if params.Keyword != "" {
		limit := params.GetLimit()
		// 优先查名字，再查作者
		keyword := params.Keyword
		params.NameKeyword = keyword
		novels, err = params.queryAll(c.Context())
		if err != nil {
			return
		}
		// 如果未达到limit，则查询名称
		if len(novels) < limit {
			nameMatchNovels := novels
			params.NameKeyword = ""
			params.AuthorKeyword = keyword
			novels, err = params.queryAll(c.Context())
			if err != nil {
				return
			}
			novels = append(nameMatchNovels, novels...)
		}
		if len(novels) > limit {
			novels = novels[:limit]
		}
		count = len(novels)
		// 有符合条件的搜索才记录关键字
		if count > 0 {
			// 如果添加不成功忽略
			_ = novelSrv.AddHotKeyword(params.Keyword)
		}
	} else {
		if params.ShouldCount() {
			count, err = params.count(c.Context())
			if err != nil {
				return
			}
		}
		novels, err = params.queryAll(c.Context())
		if err != nil {
			return
		}
	}
	c.CacheMaxAge(5 * time.Minute)
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
		Where(entNovel.ID(id)).
		First(c.Context())
	if err != nil {
		return
	}
	c.CacheMaxAge(10 * time.Minute)
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
	err = params.update(c.Context())
	if err != nil {
		return
	}
	c.NoContent()
	return
}

// listChapter 获取小说章节
func (*novelCtrl) listChapter(c *elton.Context) (err error) {
	id, err := getIDFromParams(c)
	if err != nil {
		return
	}
	params := novelChapterListParams{}
	err = validate.Do(&params, c.Query())
	if err != nil {
		return
	}
	params.ID = id
	count := -1
	if params.ShouldCount() {
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
	c.CacheMaxAge(5 * time.Minute)
	c.Body = &novelChapterListResp{
		Count:    count,
		Chapters: chapters,
	}
	return
}

// getCover 获取小说封面
func (*novelCtrl) getCover(c *elton.Context) (err error) {
	id, err := getIDFromParams(c)
	if err != nil {
		return
	}
	params := novelCoverParams{}
	err = validate.Do(&params, c.Query())
	if err != nil {
		return
	}
	cover, err := novelSrv.GetCover(id)
	if err != nil {
		return
	}
	width, _ := strconv.Atoi(params.Width)
	height, _ := strconv.Atoi(params.Height)
	quality, _ := strconv.Atoi(params.Quality)

	data, header, err := imageSrv.GetImageFromBucket(
		c.Context(),
		eliteCoverBucket,
		cover,
		service.ImageOptimizeParams{
			Type:    params.Type,
			Width:   width,
			Height:  height,
			Quality: quality,
		},
	)
	if err != nil {
		return
	}

	c.MergeHeader(header)
	c.CacheMaxAge(time.Hour)
	c.Body = data
	return
}

func updateNovelChapters(id int, fetchingContent bool) (err error) {
	err = novelSrv.UpdateChapters(id)
	if err != nil {
		return
	}
	if !fetchingContent {
		return
	}
	err = novelSrv.FetchAllChapterContent(id)
	if err != nil {
		return
	}
	// 如果更新章节数失败，则忽略
	_ = novelSrv.UpdateChapterCount(id)
	return novelSrv.UpdateWordCount(id, time.Unix(0, 0))
}

// updateAllChapters 更新所有小说章节
func (*novelCtrl) updateAllChapters(c *elton.Context) (err error) {
	id, err := novelSrv.GetMaxID()
	if err != nil {
		return
	}
	fetchingContent := c.QueryParam("fetching") != ""
	go func() {
		for i := 1; i <= id; i++ {
			err := updateNovelChapters(i, fetchingContent)
			if err != nil {
				log.Default().Error().
					Int("id", i).
					Err(err).
					Msg("update chapters fail")
				continue
			}
		}
	}()
	c.NoContent()
	return
}

// listHotKeyword 获取热门搜索关键字
func (*novelCtrl) listHotKeyword(c *elton.Context) (err error) {
	keywords, err := novelSrv.ListHotKeyword()
	if err != nil {
		return
	}
	c.CacheMaxAge(5 * time.Minute)
	c.Body = &novelHotKeywordListResp{
		Keywords: keywords,
	}
	return
}
