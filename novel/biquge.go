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
	"github.com/vicanso/elite/cache"
	"github.com/vicanso/elite/config"
	"github.com/vicanso/elite/ent/novelsource"
	"github.com/vicanso/elite/log"
	"github.com/vicanso/elite/request"
	"github.com/vicanso/go-axios"
	lruttl "github.com/vicanso/lru-ttl"
)

var biQuGeIns *axios.Instance

func init() {
	service := "biquge"
	conf := config.GetNovelConfigs().Find(service)
	biQuGeIns = request.NewHTTP(service, conf.BaseURL, conf.Timeout)
	request.Register(service, biQuGeIns)
}

const (
	// 详情接口
	biQuGeDetailURL = "/book/:id/"
	// 封面接口
	biQuGeCoverURL = "/files/article/image/:prefix/:id/:ids.jpg"
)

type biQuGe struct {
	ins   *axios.Instance
	max   int
	cache *lruttl.L2Cache
}
type biQuGeNovel struct {
	biQuGe *biQuGe
	id     int
}

// NewBiQuGe 初始化biquge小说网站实例
func NewBiQuGe() *biQuGe {
	return &biQuGe{
		ins:   biQuGeIns,
		max:   50000,
		cache: cache.NewMultilevelCache(50, 5*time.Minute, "biquge:"),
	}
}

func (n *biQuGeNovel) GetDetail() (novel Novel, err error) {
	return n.biQuGe.GetDetail(n.id)
}

func (n *biQuGeNovel) GetChapters() (chpaters []*Chapter, err error) {
	return n.biQuGe.GetChapters(n.id)
}

func (n *biQuGeNovel) GetChapterContent(no int) (content string, err error) {
	return n.biQuGe.GetChapterContent(n.id, no)
}

// NewFetcher 新建fetcher
func (bgq *biQuGe) NewFetcher(id int) Fetcher {
	return &biQuGeNovel{
		id:     id,
		biQuGe: bgq,
	}
}

type biQuGeDetail struct {
	Data []byte `json:"data,omitempty"`
}

func (bqg *biQuGe) getDetail(id int) (data []byte, err error) {
	key := fmt.Sprintf("detail-%d", id)
	detail := biQuGeDetail{}
	// 忽略出错
	_ = bqg.cache.Get(key, &detail)
	if len(detail.Data) != 0 {
		data = detail.Data
		return
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

	_ = bqg.cache.Set(key, &biQuGeDetail{
		Data: data,
	})
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
	summary := strings.TrimSpace(doc.Find("#maininfo #intro").Text())

	novel = Novel{
		Name:     name,
		Author:   authorInfos[1],
		Summary:  summary,
		Source:   NovelSourceBiQuGe,
		SourceID: id,
		CoverURL: bqg.getCoverURL(id),
	}
	return
}

func (bqg *biQuGe) getCoverURL(id int) string {
	prefix := strconv.Itoa(id / 1000)
	url := strings.ReplaceAll(biQuGeCoverURL, ":id", strconv.Itoa(id))
	url = strings.ReplaceAll(url, ":prefix", prefix)
	return bqg.ins.Config.BaseURL + url
}

// GetCover 获取封面
func (bqg *biQuGe) GetCover(id int) (img image.Image, err error) {
	resp, err := bqg.ins.Get(bqg.getCoverURL(id))
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
func (bqg *biQuGe) GetChapters(id int) (chapters []*Chapter, err error) {
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
	chapters = make([]*Chapter, max)
	for i := 0; i < max; i++ {
		item := items.Eq(i)
		title := item.Text()
		href, _ := item.Find("a").Attr("href")
		chapters[i] = &Chapter{
			Title: title,
			NO:    i,
			URL:   href,
		}
	}
	return
}

// GetChapterContent 获取小说章节内容
func (bqg *biQuGe) GetChapterContent(id, no int) (content string, err error) {
	chapters, err := bqg.GetChapters(id)
	if err != nil {
		return
	}
	if no >= len(chapters) {
		// 正常一般不会出错超出范围，因此不使用hes error
		err = errors.New("该章节已超出最新章节")
		return
	}

	var resp *axios.Response
	var doc *goquery.Document
	var html string
	for i := 0; i < 3; i++ {
		resp, err = bqg.ins.Get(chapters[no].URL)
		if err != nil {
			continue
		}
		doc, err = goquery.NewDocumentFromReader(bytes.NewReader(resp.Data))
		if err != nil {
			continue
		}
		html, err = doc.Find("#content").Html()
		if err != nil {
			continue
		}
		// 如果成功，则break
		break
	}
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
		data = append(data, value)
	}
	content = strings.Join(data, "\n")
	return
}

// Sync 同步小说来源
func (bqg *biQuGe) Sync() (err error) {
	id := 0
	for i := id + 1; i < bqg.max; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		// 如果已存在，则忽略
		exists, _ := getEntClient().NovelSource.Query().
			Where(novelsource.SourceIDEQ(i)).
			Where(novelsource.SourceEQ(NovelSourceBiQuGe)).
			Exist(ctx)
		if exists {
			continue
		}
		novel, err := bqg.GetDetail(i)
		if err != nil {
			log.Default().Error().
				Int("id", i).
				Msg("sync novel fail")
			continue
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
