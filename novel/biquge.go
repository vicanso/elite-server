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

package novel

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/vicanso/elite/ent"
	"github.com/vicanso/elite/ent/novelsource"
	"github.com/vicanso/elite/helper"
	"github.com/vicanso/go-axios"
	lruttl "github.com/vicanso/lru-ttl"
)

const (
	// 详情接口
	biQuGeDetailURL = "/book/:id/"
	// 封面接口
	biQuGeCoverURL = "/files/article/image/:prefix/:id/:ids.jpg"
)

type biQuGe struct {
	ins   *axios.Instance
	max   int
	cache *lruttl.Cache
}

// NewBiQuGe 初始化biquge小说网站实例
func NewBiQuGe() *biQuGe {
	conf := getNovelConfig(novelBiQuGeName)
	if conf.Name == "" {
		panic("get biquge's config fail")
	}
	ins := helper.NewInstance(conf.Name, conf.BaseURL, conf.Timeout)
	return &biQuGe{
		ins:   ins,
		max:   conf.Max,
		cache: lruttl.New(50, 5*time.Minute),
	}
}

func (bqg *biQuGe) getDetail(id int) (data []byte, err error) {
	key := fmt.Sprintf("detail-%d", id)
	if result, ok := bqg.cache.Get(key); ok {
		data, ok = result.([]byte)
		if ok {
			return
		}
	}
	// 如果出错则继续拉取，拉取两次
	for i := 0; i < 2; i++ {
		conf := &axios.Config{
			URL: biQuGeDetailURL,
			Params: map[string]string{
				"id": strconv.Itoa(id),
			},
		}
		resp, e := bqg.ins.Request(conf)
		err = e
		if e != nil {
			continue
		}
		data = resp.Data
		break
	}
	if err != nil {
		return
	}

	bqg.cache.Add(key, data)
	return
}

// GetDetail 根据ID获取小说详情
func (bqg *biQuGe) GetDetail(id int) (novel Novel, err error) {
	data, err := bqg.getDetail(id)
	if err != nil {
		return
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		return
	}
	info := doc.Find("#maininfo #info")

	name := info.Find("h1").Text()
	if name == "" {
		return
	}
	authorInfos := strings.Split(info.Find("p").First().Text(), "：")
	if len(authorInfos) != 2 {
		return
	}

	novel = Novel{
		Name:     name,
		Author:   authorInfos[1],
		Source:   NovelSourceBiQuGe,
		SourceID: id,
	}
	return
}

// GetCover 获取封面
func (bqg *biQuGe) GetCover(id int) (img image.Image, err error) {
	conf := &axios.Config{
		URL: biQuGeCoverURL,
		Params: map[string]string{
			"id":     strconv.Itoa(id),
			"prefix": strconv.Itoa(id / 1000),
		},
	}
	resp, err := bqg.ins.Request(conf)
	if err != nil {
		return
	}
	// 少于10KB的封面认为无封面
	if len(resp.Data) < 10*1024 {
		err = errors.New("cover not found")
		return
	}
	img, err = jpeg.Decode(bytes.NewReader(resp.Data))
	if err != nil {
		return
	}
	return
}

// GetChapters 获取小说章节列表
func (bqg *biQuGe) GetChapters(id int) (chpaters []*NovelChapter, err error) {
	data, err := bqg.getDetail(id)
	if err != nil {
		return
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		return
	}
	items := doc.Find("#list dd")
	max := items.Length()
	chpaters = make([]*NovelChapter, max)
	for i := 0; i < max; i++ {
		item := items.Eq(i)
		title := item.Text()
		href, _ := item.Find("a").Attr("href")
		chpaters[i] = &NovelChapter{
			Title: title,
			NO:    i,
			URL:   href,
		}
	}
	return
}

// GetChapterContent 获取小说章节内容
func (bqg *biQuGe) GetChapterContent(url string) (content string, err error) {
	resp, err := bqg.ins.Get(url)
	if err != nil {
		return
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Data))
	if err != nil {
		return
	}
	html, err := doc.Find("#content").Html()
	if err != nil {
		return
	}
	arr := strings.Split(html, "<br/>")
	data := make([]string, 0, len(arr))
	for _, item := range arr {
		value := strings.TrimSpace(item)
		if value == "" {
			continue
		}
		data = append(data, item)
	}
	content = strings.Join(data, "\n")
	return
}

// Sync 同步小说来源
func (bqg *biQuGe) Sync() (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// 获取当前源的最后更新记录
	source, err := getEntClient().NovelSource.Query().
		Where(novelsource.SourceEQ(NovelSourceBiQuGe)).
		Order(ent.Desc("source_id")).
		First(ctx)

	if err != nil && !ent.IsNotFound(err) {
		return
	}
	var id int
	if source != nil {
		id = source.SourceID
	}
	for i := id + 1; i < bqg.max; i++ {
		novel, err := bqg.GetDetail(i)
		if err != nil {
			return err
		}
		if novel.SourceID == 0 {
			continue
		}
		_, err = novel.AddToSource()
		if err != nil {
			return err
		}
	}
	return
}
