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

	// novelAddParams 添加小说参数
	novelAddParams struct {
		Name   string `json:"name,omitempty" validate:"required,xNovelName"`
		Author string `json:"author,omitempty" validate:"required,xNovelAuthor"`
	}
)

func init() {
	ctrl := novelCtrl{}

	g := router.NewGroup("/novels")

	g.POST(
		"/v1",
		loadUserSession,
		shouldBeAdmin,
		ctrl.add,
	)
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
