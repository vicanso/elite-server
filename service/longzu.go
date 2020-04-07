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

// 江南 龙族

package service

import (
	"bytes"
	"errors"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/vicanso/go-axios"
	"go.uber.org/zap"
)

type (
	LongzuSrv struct{}
)

// 初始化http请求实例
var longZuHeader = http.Header{}
var longZuIns = axios.NewInstance(&axios.InstanceConfig{
	BaseURL: "https://www.luoxia.com/longzu/",
	Timeout: 30 * time.Second,
	Headers: longZuHeader,
})

func init() {
	longZuHeader.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.130 Safari/537.36")
}

func (srv *LongzuSrv) sync(url string) (err error) {
	resp, err := longZuIns.Get(url)
	if err != nil {
		return
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Data))
	if err != nil {
		return
	}
	bookDescribe := doc.Find(".book-describe")
	novel := &Novel{
		Name:    bookDescribe.Find("h1").Text(),
		Author:  "江南",
		Summary: bookDescribe.Find(".describe-html").Text(),
	}
	novelSrv := &NovelSrv{}
	found, _ := novelSrv.FindOne(&Novel{
		Name:   novel.Name,
		Author: novel.Author,
	})
	var bookID uint
	if found != nil && found.ID != 0 {
		bookID = found.ID
	} else {
		novel, err = novelSrv.Add(*novel)
		if err != nil {
			return
		}
		bookID = novel.ID
	}
	if bookID == 0 {
		err = errors.New("get book id fail")
		return
	}
	doc.Find(".book-list li").Each(func(i int, s *goquery.Selection) {
		if err != nil {
			return
		}
		result, _ := novelSrv.FindOneChapter(&NovelChapter{
			BookID: bookID,
			NO:     i,
		})
		// 已存在
		if result != nil && result.ID != 0 {
			return
		}
		href, _ := s.Find("a").Attr("href")
		if href == "" {
			onClick, _ := s.Find("b").Attr("onclick")
			reg := regexp.MustCompile(`(https://www.luoxia.com/longzu/\d+.htm)`)
			href = reg.FindString(onClick)
		}
		if href == "" {
			err = errors.New("get page url fail")
			return
		}
		tmpResp, e := longZuIns.Get(href)
		if e != nil {
			err = e
			return
		}
		tmpDoc, e := goquery.NewDocumentFromReader(bytes.NewReader(tmpResp.Data))
		if e != nil {
			err = e
			return
		}
		contentArr := make([]string, 0)
		tmpDoc.Find("#nr1 p").Each(func(i int, s *goquery.Selection) {
			contentArr = append(contentArr, strings.TrimSpace(s.Text()))
		})
		content := strings.Join(contentArr, "\n")
		novelChapter := NovelChapter{
			BookID:    bookID,
			NO:        i,
			Title:     tmpDoc.Find("#nr_title").Text(),
			Content:   content,
			WordCount: len(content),
		}
		_, _ = novelSrv.AddChapter(novelChapter)
	})
	if err != nil {
		return
	}

	err = novelSrv.RefreshBasicInfo(&Novel{
		ID: bookID,
	})
	if err != nil {
		return
	}

	return
}

func (srv *LongzuSrv) Sync() (err error) {
	resp, err := longZuIns.Get("")
	if err != nil {
		return
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Data))
	if err != nil {
		return
	}
	doc.Find("#content-list .title").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Find("a").Attr("href")
		err := srv.sync(href)
		if err != nil {
			logger.Error("sync fail",
				zap.String("url", href),
				zap.Error(err),
			)
		}
	})

	return
}
