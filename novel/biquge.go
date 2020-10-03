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
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/vicanso/go-axios"
)

const (
	// 详情接口
	biQuGeDetailURL = "/book/:id/"
)

type biQuGe struct {
	ins *axios.Instance
}

// NewBiQuGe 初始化biquge小说网站实例
func NewBiQuGe() *biQuGe {
	conf := getNovelConfig(novelBiQuGeName)
	if conf.Name == "" {
		panic("get biquge's config fail")
	}
	ins := axios.NewInstance(&axios.InstanceConfig{
		BaseURL: conf.BaseURL,
		Timeout: conf.Timeout,
	})
	return &biQuGe{
		ins: ins,
	}
}

// GetDetail 根据ID获取小说详情
func (bqg *biQuGe) GetDetail(id int) (novel Novel, err error) {
	conf := &axios.Config{
		URL: biQuGeDetailURL,
		Params: map[string]string{
			"id": strconv.Itoa(id),
		},
	}
	resp, err := bqg.ins.Request(conf)
	if err != nil {
		return
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Data))
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
