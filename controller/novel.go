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
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/vicanso/elite/ent"
	entNovel "github.com/vicanso/elite/ent/novel"
	"github.com/vicanso/elite/novel"
	"github.com/vicanso/elite/router"
	"github.com/vicanso/elite/service"
	"github.com/vicanso/elite/util"
	"github.com/vicanso/elite/validate"
	"github.com/vicanso/elton"
	"github.com/vicanso/go-axios"
	"go.uber.org/zap"
)

const eliteConverBucket = "elite-covers"

type (
	novelCtrl struct{}

	// novelListResp 小说列表响应
	novelListResp struct {
		Novels []*ent.Novel `json:"novels,omitempty"`
		Count  int          `json:"count,omitempty"`
	}

	// novelAddParams 添加小说参数
	novelAddParams struct {
		Name   string `json:"name,omitempty" validate:"required,xNovelName"`
		Author string `json:"author,omitempty" validate:"required,xNovelAuthor"`
	}
	// novelListParams 小说查询参数
	novelListParams struct {
		listParams

		Keyword string `json:"keyword,omitempty" validate:"omitempty,xKeyword"`
	}
)

func init() {
	ctrl := novelCtrl{}

	g := router.NewGroup("/novels")

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
	query = params.where(query)
	return query.All(ctx)
}

// count 计算总数
func (params *novelListParams) count(ctx context.Context) (count int, err error) {
	query := getEntClient().Novel.Query()
	query = params.where(query)

	return query.Count(ctx)
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
				Bucket: eliteConverBucket,
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
