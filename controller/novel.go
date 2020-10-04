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
	"github.com/vicanso/elite/novel"
	"github.com/vicanso/elite/router"
	"github.com/vicanso/elite/validate"
	"github.com/vicanso/elton"
)

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
	result, err := novel.Publish(novel.QueryParams{
		Name:   params.Name,
		Author: params.Author,
	})
	if err != nil {
		return
	}
	c.Created(result)
	return
}
