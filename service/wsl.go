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

// 卫斯理科幻小说系列

package service

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/vicanso/elite/log"
	"go.uber.org/zap"
)

type (
	WslSrv struct{}

	wslBook struct {
		Name    string `json:"name,omitempty"`
		Author  string `json:"author,omitempty"`
		Summary string `json:"summary,omitempty"`
		Cover   string `json:"cover,omitempty"`
	}
	wslBooks struct {
		Books []*wslBook `json:"books,omitempty"`
	}
	wslChapter struct {
		Title   string
		Content string
		NO      int
	}
	wslChapters struct {
		Chapters []*wslChapter
	}
)

// listAll list all
func (srv *WslSrv) listAll() (books []*wslBook, err error) {
	limit := 10
	offset := 0
	books = make([]*wslBook, 0)
	for {
		url := fmt.Sprintf(wslBooksURL+"?limit=%d&offset=%d", limit, offset)
		resp, e := wslIns.Get(url)
		if e != nil {
			err = e
			return
		}
		result := new(wslBooks)
		err = resp.JSON(result)
		if err != nil {
			return
		}
		books = append(books, result.Books...)
		// 已至最后一页
		if len(result.Books) < limit {
			break
		}
		offset += limit
	}
	return
}

// listAllChapters list all chapters
func (srv *WslSrv) listAllChapters(bookID uint) (chapters []*wslChapter, err error) {
	offset := 0
	limit := 10
	chapters = make([]*wslChapter, 0)
	for {
		url := fmt.Sprintf(wslBookChaptersURL, bookID, limit, offset)
		resp, e := wslIns.Get(url)
		if e != nil {
			err = e
			return
		}
		result := new(wslChapters)
		err = resp.JSON(result)
		if err != nil {
			return
		}
		chapters = append(chapters, result.Chapters...)
		if len(result.Chapters) < limit {
			break
		}
		offset += limit
	}
	return
}

func (srv *WslSrv) syncChapters(bookID uint) (err error) {
	chapters, err := srv.listAllChapters(bookID)
	if err != nil {
		return
	}
	for _, item := range chapters {
		chapter := &NovelChapter{}
		result := pgGetClient().First(chapter, &NovelChapter{
			BookID: bookID,
			NO:     item.NO,
		})
		if result.Error != nil && !result.RecordNotFound() {
			err = result.Error
			return
		}
		if chapter.ID == 0 {
			chapter.BookID = bookID
			chapter.NO = item.NO
			chapter.Content = item.Content
			chapter.WordCount = len(item.Content)
			chapter.Title = item.Title
			err = pgGetClient().Save(chapter).Error
			if err != nil {
				return
			}
		}
	}
	return
}

// Sync sync novel from wsl
func (srv *WslSrv) Sync() (err error) {
	books, err := srv.listAll()
	if err != nil {
		return
	}
	for _, book := range books {
		result, err := new(NovelSrv).FindOne(&Novel{
			Name:   book.Name,
			Author: book.Author,
		})
		// 如果该记录不存在，则插入
		if err == gorm.ErrRecordNotFound {
			result = &Novel{
				Name:    book.Name,
				Author:  book.Author,
				Summary: book.Summary,
				Status:  StatusUnfinished,
			}
			pgGetClient().Save(result)
		}
		// 如果cover更新失败，忽略错误
		_ = new(NovelSrv).UpdateCover(result.ID, book.Cover, false)

		_ = srv.syncChapters(result.ID)
		log.Default().Info("wsl sync",
			zap.Uint("id", result.ID),
			zap.String("name", result.Name),
		)
	}
	return
}
