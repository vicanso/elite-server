// Copyright 2019 tree xie
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

package controller

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/vicanso/elite/router"
	"github.com/vicanso/elite/service"
	"github.com/vicanso/elite/util"
	"github.com/vicanso/elite/validate"
	"github.com/vicanso/elton"
	"github.com/vicanso/hes"
)

type novelCtrl struct{}

type (
	syncParams struct {
		Token string `json:"token,omitempty" valid:"runelength(1|30)"`
	}
	updateCoverParams struct {
		Cover string `json:"cover,omitempty" valid:"url"`
	}
)

func init() {
	g := router.NewGroup("/novels")
	ctrl := novelCtrl{}
	// 获取书籍列表
	g.GET(
		"/v1",
		ctrl.list,
	)
	// 获取书籍详情
	g.GET(
		"/v1/:id",
		ctrl.detail,
	)
	// 更新书籍详情
	g.PATCH(
		"/v1/:id",
		loadUserSession,
		shouldBeAdmin,
		ctrl.update,
	)
	// 获取章节列表
	g.GET(
		"/v1/:id/chapters",
		ctrl.listChapters,
	)
	// 更新章节内容
	g.PATCH(
		"/v1/:id/chapters/:chapterId",
		loadUserSession,
		shouldBeAdmin,
		ctrl.updateChapter,
	)
	// 更新书籍封面
	g.PATCH(
		"/v1/:id/cover",
		loadUserSession,
		shouldBeAdmin,
		ctrl.updateCover,
	)
	// 获取书籍封面
	g.GET(
		"/v1/:id/cover",
		ctrl.getCover,
	)

	g.POST(
		"/v1/sync-wsl",
		ctrl.sync,
	)
}

// list 书籍列表查询
func (ctrl novelCtrl) list(c *elton.Context) (err error) {
	params, err := getDbQueryParams(c)
	if err != nil {
		return
	}
	where := make([]interface{}, 0)
	keyword := c.QueryParam("keyword")
	// 关键字搜索，暂仅支持对书名搜索
	if keyword != "" {
		where = append(where, "name LIKE ?", "%"+keyword+"%")
	}

	// 指定ID返回
	ids := c.QueryParam("ids")
	if ids != "" {
		where = append(where, "id IN (?)", strings.Split(ids, ","))
	}

	novels, err := novelSrv.List(params, where...)
	if err != nil {
		return
	}
	c.CacheMaxAge("5m")
	c.Body = map[string]interface{}{
		"novels": novels,
	}
	return
}

// detail 书籍详情
func (ctrl novelCtrl) detail(c *elton.Context) (err error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return
	}
	novel, err := novelSrv.FindOne(&service.Novel{
		ID: uint(id),
	})
	if err != nil {
		return
	}
	c.CacheMaxAge("5m")
	c.Body = novel
	return
}

// update 更新书籍信息
func (ctrl novelCtrl) update(c *elton.Context) (err error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return
	}
	novel := &service.Novel{}
	err = json.Unmarshal(c.RequestBody, novel)
	if err != nil {
		return
	}
	novel.ID = uint(id)
	err = novelSrv.Update(novel)
	if err != nil {
		return
	}
	c.NoContent()
	return
}

// updateChapter 更新章节信息
func (ctrl novelCtrl) updateChapter(c *elton.Context) (err error) {
	// 直接根据ID则可更新
	chapterId, err := strconv.Atoi(c.Param("chapterId"))
	if err != nil {
		return
	}
	novelChapter := &service.NovelChapter{}
	err = json.Unmarshal(c.RequestBody, novelChapter)
	if err != nil {
		return
	}

	novelChapter.ID = uint(chapterId)
	if novelChapter.Content != "" {
		novelChapter.WordCount = len(novelChapter.Content)
	}
	err = novelSrv.UpdateChapter(novelChapter)
	if err != nil {
		return
	}
	c.NoContent()
	return
}

// listChapters 章节列表查询
func (ctrl novelCtrl) listChapters(c *elton.Context) (err error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return
	}
	params, err := getDbQueryParams(c)
	if err != nil {
		return
	}
	chapters, err := novelSrv.ListChapters(params, &service.NovelChapter{
		BookID: uint(id),
	})
	if err != nil {
		return
	}
	c.CacheMaxAge("5m")
	c.Body = map[string]interface{}{
		"chapters": chapters,
	}
	return
}

// updateCover 更新书籍封面
func (ctrl novelCtrl) updateCover(c *elton.Context) (err error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return
	}
	params := &updateCoverParams{}
	err = validate.Do(params, c.RequestBody)
	if err != nil {
		return
	}
	err = novelSrv.UpdateCover(uint(id), params.Cover, true)
	if err != nil {
		return
	}
	c.NoContent()
	return
}

// getCover 获取封面
func (ctrl novelCtrl) getCover(c *elton.Context) (err error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return
	}
	width, _ := strconv.Atoi(c.QueryParam("width"))
	height, _ := strconv.Atoi(c.QueryParam("height"))
	cover, err := novelSrv.GetCover(uint(id), &service.ImageOptimParams{
		Width:  width,
		Height: height,
		Output: c.QueryParam("output"),
		Effect: c.QueryParam("effect"),
	})
	if err != nil {
		return
	}

	c.CacheMaxAge("1h")
	if cover == nil {
		c.NoContent()
	} else {
		c.SetHeader(elton.HeaderContentType, cover.ContentType)
		c.BodyBuffer = bytes.NewBuffer(cover.Data)
	}
	return
}

func (ctrl novelCtrl) sync(c *elton.Context) (err error) {
	params := &syncParams{}
	err = validate.Do(params, c.RequestBody)
	if err != nil {
		return
	}
	if util.Sha256(params.Token) != "ua/wtwoSlY1sq90dLi3dnpasYBqhycmQXXJE2iw7MzM=" {
		err = hes.New("token is invalid")
		return
	}
	go func() {
		_ = new(service.WslSrv).Sync()
	}()
	c.NoContent()
	return
}
