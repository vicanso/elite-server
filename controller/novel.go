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
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/vicanso/elite/cs"
	"github.com/vicanso/elite/ent"
	"github.com/vicanso/elite/ent/chapter"
	entNovel "github.com/vicanso/elite/ent/novel"
	"github.com/vicanso/elite/ent/novelsource"
	"github.com/vicanso/elite/ent/schema"
	"github.com/vicanso/elite/novel"
	"github.com/vicanso/elite/router"
	"github.com/vicanso/elite/service"
	"github.com/vicanso/elite/util"
	"github.com/vicanso/elite/validate"
	"github.com/vicanso/elton"
	"github.com/vicanso/go-axios"
	"github.com/vicanso/hes"
	"go.uber.org/zap"
)

const eliteCoverBucket = "elite-covers"

const errNovelCategory = "novel"

var (
	novelNoMatchRecord = &hes.Error{
		Message:    "没有匹配的记录",
		StatusCode: http.StatusBadRequest,
		Category:   errNovelCategory,
	}
	novelIDInvalid = &hes.Error{
		Message:    "ID不能为空",
		StatusCode: http.StatusBadRequest,
		Category:   errNovelCategory,
	}
)

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
	// novelChapterListParams 章节查询参数
	novelChapterListParams struct {
		listParams

		// ID 小说id，由route param中获取并设置，因此不设置validate
		ID int `json:"id,omitempty"`
		// ChapterID 章节id，由route param中获取并设置，因此不设置validate
		ChapterID int `json:"chapterID,omitempty"`
	}

	// novelSourceListParams 小说源查询参数
	novelSourceListParams struct {
		listParams

		Status  string `json:"status,omitempty" validate:"omitempty,xNovelStatus"`
		Keyword string `json:"keyword,omitempty" validate:"omitempty,xKeyword"`
	}
	// novelUserBehaviorParams 用户行为参数
	novelUserBehaviorParams struct {
		Category string `json:"category,omitempty" validate:"required,xNovelBehaviorCategory"`

		// ID 小说id，由route param中获取并设置，因此设置omitempty
		ID int `json:"id,omitempty" validate:"omitempty,xNovelID"`
	}
	// novelCoverParams 小说封面参数
	novelCoverParams struct {
		Type    string `json:"type,omitempty" validate:"required,xNovelCoverType"`
		Width   string `json:"width,omitempty" validate:"omitempty,xNovelCoverWidth"`
		Height  string `json:"height,omitempty" validate:"omitempty,xNovelCoverHeight"`
		Quality string `json:"quality,omitempty" validate:"required,xNovelCoverQuality"`
	}
	// novelChapterUpdateParams 小说章节更新参数
	novelChapterUpdateParams struct {
		ID      int    `json:"id,omitempty"`
		Title   string `json:"title,omitempty" validate:"omitempty,xNovelChapterTitle"`
		Content string `json:"content,omitempty" validate:"omitempty,xNovelChapterContent"`
	}
)

func init() {
	ctrl := novelCtrl{}

	g := router.NewGroup("/novels", setNoCacheIfMatched)

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
	// 用户行为
	g.POST(
		"/v1/{id}/behaviors",
		loadUserSession,
		ctrl.addBehavior,
	)
	// 发布所有的小说
	g.POST(
		"/v1/publish-all",
		loadUserSession,
		shouldBeAdmin,
		ctrl.publishAll,
	)
	// 更新所有章节
	g.POST(
		"/v1/update-all-chapters",
		loadUserSession,
		shouldBeAdmin,
		ctrl.updateAllChapters,
	)
	// 指定小说更新所有章节
	g.POST(
		"/v1/{id}/update-chapters",
		newTracker(cs.ActionNovelChaptersUpdate),
		loadUserSession,
		shouldBeAdmin,
		ctrl.updateChaptersByID,
	)
	// 根据ID获取小说章节
	g.GET(
		"/v1/chapters/{id}",
		ctrl.getChapterByID,
	)
	// 根据ID更新小说章节
	g.PATCH(
		"/v1/chapters/{id}",
		newTracker(cs.ActionNovelChapterUpdate),
		loadUserSession,
		shouldBeAdmin,
		ctrl.updateChapterByID,
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
	query := getEntClient().Chapter.Query()
	query = query.Limit(params.GetLimit()).
		Offset(params.GetOffset()).
		Order(params.GetOrders()...)
	query = params.where(query)
	return query.All(ctx)
}

// count 计算章节总数
func (params *novelChapterListParams) count(ctx context.Context) (count int, err error) {
	query := getEntClient().Chapter.Query()
	query = params.where(query)
	return query.Count(ctx)
}

// where 将查询参数转换为where条件
func (params *novelSourceListParams) where(query *ent.NovelSourceQuery) *ent.NovelSourceQuery {
	if params.Keyword != "" {
		query = query.Where(novelsource.Or(novelsource.NameContains(params.Keyword), novelsource.AuthorContains(params.Keyword)))
	}
	if params.Status != "" {
		status, _ := strconv.Atoi(params.Status)
		query = query.Where(novelsource.StatusEQ(status))
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
func (params *novelUpdateParams) update(ctx context.Context) (err error) {
	if params.ID == 0 {
		err = novelIDInvalid
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
		err = novelNoMatchRecord
		return
	}
	return
}

func (params *novelChapterUpdateParams) update(ctx context.Context) (err error) {
	if params.ID == 0 {
		err = novelIDInvalid
		return
	}
	update := getEntClient().Chapter.UpdateOneID(params.ID)
	if params.Title != "" {
		update = update.SetTitle(params.Title)
	}
	if params.Content != "" {
		update = update.SetContent(params.Content).
			SetWordCount(len(params.Content))
	}
	result, err := update.Save(ctx)
	if err != nil {
		return
	}
	if result == nil {
		err = novelNoMatchRecord
		return
	}
	return
}

// do 用户行为记录
func (params *novelUserBehaviorParams) do(ctx context.Context) (err error) {
	update := getEntClient().Novel.UpdateOneID(params.ID)
	switch params.Category {
	case cs.ActionNovelUserView:
		update = update.AddViews(1)
	case cs.ActionNovelUserDownload:
		update = update.AddDownloads(1)
	case cs.ActionNovelUserFavorite:
		update = update.AddFavorites(1)
	default:
		err = errors.New(params.Category + " is not supported")
		return
	}
	_, err = update.Save(ctx)
	if err != nil {
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

func (*novelCtrl) publish(params novel.QueryParams) (result *ent.Novel, err error) {
	result, err = novelSrv.Publish(params)
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
					zap.Error(err),
				)
				return
			}
			contentType := resp.Headers.Get("Content-Type")
			fileType := strings.Split(contentType, "/")[1]
			if fileType == "html" {
				logger.Error("get cover fail",
					zap.String("name", params.Name),
					zap.Error(errors.New("type of cover is invalid")),
				)
				return
			}
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
	return
}

// add 添加小说
func (ctrl *novelCtrl) add(c *elton.Context) (err error) {
	params := novelAddParams{}
	err = validate.Do(&params, c.RequestBody)
	if err != nil {
		return
	}
	result, err := ctrl.publish(novel.QueryParams{
		Name:   params.Name,
		Author: params.Author,
	})
	if err != nil {
		return
	}

	c.Created(result)
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
				logger.Error("publish novel fail",
					zap.String("name", result.Name),
					zap.String("author", result.Author),
					zap.Error(err),
				)
			}
		}
		logger.Info("publish all done")
	}()
	c.NoContent()
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
	return novelSrv.UpdateWordCount(id)
}

// updateChaptersByID 更新单本小说章节
func (*novelCtrl) updateChaptersByID(c *elton.Context) (err error) {
	id, err := getIDFromParams(c)
	if err != nil {
		return
	}
	// 章节内容较多，开新的goroutine处理
	go func() {
		err := updateNovelChapters(id, true)
		if err != nil {
			logger.Error("fetch chapter content fail",
				zap.Int("id", id),
				zap.Error(err),
			)
		}
	}()
	c.NoContent()
	return
}

// updateAllChapters 更新所有小说章节
func (*novelCtrl) updateAllChapters(c *elton.Context) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	id, err := getEntClient().Novel.Query().
		Order(ent.Desc("id")).
		FirstID(ctx)
	if err != nil {
		return
	}
	fetchingContent := c.QueryParam("fetching") != ""
	go func() {
		for i := 1; i <= id; i++ {
			err := updateNovelChapters(i, fetchingContent)
			if err != nil {
				logger.Error("update chapters fail",
					zap.Int("id", i),
					zap.Error(err),
				)
				continue
			}
		}
	}()
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
	c.CacheMaxAge(5 * time.Minute)
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
	c.CacheMaxAge(10 * time.Minute)
	c.Body = result
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

// addBehavior 添加用户行为
func (*novelCtrl) addBehavior(c *elton.Context) (err error) {
	id, err := getIDFromParams(c)
	if err != nil {
		return
	}
	params := novelUserBehaviorParams{}
	err = validate.Do(&params, c.RequestBody)
	if err != nil {
		return
	}
	params.ID = id
	err = params.do(c.Context())
	if err != nil {
		return
	}
	c.NoContent()
	return
}

// getChapterByID 通过ID获取章节
func (*novelCtrl) getChapterByID(c *elton.Context) (err error) {
	id, err := getIDFromParams(c)
	if err != nil {
		return
	}
	params := novelChapterListParams{}
	params.ChapterID = id
	chapters, err := params.queryAll(c.Context())
	if err != nil {
		return
	}
	if len(chapters) == 0 {
		err = hes.New("无法获取该章节内容")
		return
	}
	c.Body = chapters[0]
	return
}

// updateChapterByID 根据章节ID更新章节
func (*novelCtrl) updateChapterByID(c *elton.Context) (err error) {
	id, err := getIDFromParams(c)
	if err != nil {
		return
	}
	params := novelChapterUpdateParams{}
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
