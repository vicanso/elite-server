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
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/vicanso/elite/ent"
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
)

type novelCtrl struct{}

const eliteCoverBucket = "elite-covers"

// 接口参数定义
type (
	// novelListParams 小说查询参数
	novelListParams struct {
		listParams

		Keyword       string `json:"keyword,omitempty" validate:"omitempty,xKeyword"`
		AuthorKeyword string
		NameKeyword   string
	}
)

// 接口响应定义
type (
	// novelListResp 小说列表响应
	novelListResp struct {
		Novels []*ent.Novel `json:"novels,omitempty"`
		Count  int          `json:"count,omitempty"`
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
