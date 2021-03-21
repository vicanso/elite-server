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
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/vicanso/elite/request"
	"github.com/vicanso/go-axios"
)

const (
	// 查询接口
	qiDianSearchURL = "/search"
)

type qiDian struct {
	ins *axios.Instance
}

// NewQiDian 初始化qidian小说网站实例
func NewQiDian() *qiDian {
	return &qiDian{
		ins: request.GetQiDian(),
	}
}

// Search 查询小说
func (qd *qiDian) Search(name, author string) (novel Novel, err error) {
	query := make(url.Values)
	query.Add("kw", name)
	resp, err := qd.ins.Get(qiDianSearchURL, query)
	if err != nil {
		return
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Data))
	if err != nil {
		return
	}
	items := doc.Find("#result-list .res-book-item")
	count := items.Length()
	for i := 0; i < count; i++ {
		item := items.Eq(i)
		curName := item.Find("h4").Text()
		curAuthor := item.Find(".author .name").Text()
		if curName == name && curAuthor == author {
			summary := strings.TrimSpace(item.Find(".intro").Text())
			bid, _ := item.Attr("data-bid")
			id, _ := strconv.Atoi(bid)
			cover, ok := item.Find(".book-img-box img").Attr("src")
			if ok {
				cover = strings.Replace("https:"+cover, bid+"/150", bid+"/180", 1)
			}
			// TODO 添加分类
			novel = Novel{
				Name:     curName,
				Author:   curAuthor,
				Summary:  summary,
				Source:   NovelSourceQiDian,
				SourceID: id,
				CoverURL: cover,
			}
		}
	}
	return
}
