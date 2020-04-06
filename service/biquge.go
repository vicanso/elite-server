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

// 笔趣阁
package service

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"go.uber.org/zap"

	"github.com/PuerkitoBio/goquery"
	"github.com/vicanso/elite/helper"
	"github.com/vicanso/go-axios"
	lru "github.com/vicanso/lru-ttl"
	"golang.org/x/text/encoding/simplifiedchinese"
)

type (
	BiQuGeSrv struct{}

	BiQuGe struct {
		ID        uint       `gorm:"primary_key" json:"id,omitempty"`
		CreatedAt *time.Time `json:"createdAt,omitempty"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		DeletedAt *time.Time `sql:"index" json:"deletedAt,omitempty"`

		BookID int    `json:"bookId,omitempty" `
		Name   string `json:"name,omitempty" gorm:"type:varchar(100);not null;unique_index:idx_biquges_name_author"`
		Author string `json:"author,omitempty" gorm:"type:varchar(50);not null;unique_index:idx_biquges_name_author"`
	}
)

const (
	biQuGeDetailURL = "/%d_%d/"
)

const (
	biQuGeSyncTask = "bi-qu-ge-sync-task"
)

var biqugeHeader = http.Header{}
var biqugeIns = axios.NewInstance(&axios.InstanceConfig{
	BaseURL: "https://www.cnoz.org/",
	Timeout: 30 * time.Second,
	Headers: longZuHeader,
})
var biqugeCache = lru.New(20, 5*time.Minute)

func init() {
	biqugeHeader.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.130 Safari/537.36")
	transformRepsonse := make([]axios.TransformResponse, 0)
	transformRepsonse = append(transformRepsonse, axios.DefaultTransformResponse...)
	// GBK to utf8
	transformRepsonse = append(transformRepsonse, func(body []byte, headers http.Header) (data []byte, err error) {
		decoder := simplifiedchinese.GB18030.NewDecoder()
		return decoder.Bytes(body)
	})

	biqugeIns.Config.TransformResponse = transformRepsonse

	pgGetClient().
		AutoMigrate(&BiQuGe{})
	// srv := new(BiQuGeSrv)
	// fmt.Println(srv.Sync(10))
}

func (srv *BiQuGeSrv) getBasicInfoDocument(id int) (doc *goquery.Document, err error) {
	prefix := id / 1000
	url := fmt.Sprintf(biQuGeDetailURL, prefix, id)
	value, ok := biqugeCache.Get(url)
	if ok {
		doc = value.(*goquery.Document)
		return
	}

	resp, err := biqugeIns.Get(url)
	if err != nil {
		return
	}
	doc, err = goquery.NewDocumentFromReader(bytes.NewReader(resp.Data))
	if err != nil {
		return
	}
	biqugeCache.Add(url, doc)
	return
}

func (srv *BiQuGeSrv) GetBasicInfo(id int) (novel *Novel, err error) {
	doc, err := srv.getBasicInfoDocument(id)
	if err != nil {
		return
	}
	q := doc.Find("#maininfo #info")
	if q.Length() == 0 {
		err = errors.New("novel not found")
		return
	}
	novel = &Novel{
		Name:    strings.TrimSpace(q.Find("h1").Text()),
		Summary: strings.TrimSpace(doc.Find("#maininfo #intro").Text()),
	}
	authorText := q.Find("p").First().Text()

	arr := strings.Split(authorText, "：")
	if len(arr) == 2 {
		novel.Author = strings.TrimSpace(arr[1])
	}
	return
}

func (srv *BiQuGeSrv) GetChapters(id int) (chapters []*NovelChapter, err error) {
	doc, err := srv.getBasicInfoDocument(id)
	if err != nil {
		return
	}

	chapters = make([]*NovelChapter, 0, 100)
	doc.Find("#list dd a").Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Text())
		// url, _ := s.Attr("href")
		chapters = append(chapters, &NovelChapter{
			NO:    i,
			Title: title,
		})
	})
	return
}

func (srv *BiQuGeSrv) Add(data BiQuGe) (biquge *BiQuGe, err error) {
	biquge = &data
	err = pgGetClient().Save(biquge).Error
	if err != nil {
		return
	}
	return
}

func (srv *BiQuGeSrv) Sync(max int) (err error) {
	ok, err := redisSrv.Lock(biQuGeSyncTask, 30*time.Minute)
	if err != nil || !ok {
		return
	}
	biquge := new(BiQuGe)
	err = helper.PGGetDB(&helper.DbParams{
		Order: "-bookId",
	}).First(biquge).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}
	start := 1
	if biquge.BookID != 0 {
		start = biquge.BookID + 1
	}
	logger.Info("start to bi qu ge sync",
		zap.Int("start", start),
		zap.Int("end", max),
	)
	for i := start; i < max; i++ {
		novel, e := srv.GetBasicInfo(i)
		if e != nil {
			err = e
			return
		}
		_, err = srv.Add(BiQuGe{
			Name:   novel.Name,
			Author: novel.Author,
			BookID: i,
		})
		if err != nil {
			return
		}
	}
	logger.Info("bi qu ge sync done",
		zap.Int("start", start),
		zap.Int("end", max),
	)
	return
}
