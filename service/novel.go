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

package service

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/jinzhu/gorm"
	"github.com/vicanso/elite/config"
	"github.com/vicanso/elite/helper"
	"github.com/vicanso/elite/log"
	"github.com/vicanso/go-axios"
	"github.com/vicanso/hes"
	"go.uber.org/zap"
)

const (
	// 未知
	StatusUnknown = iota
	// 未完结
	StatusUnfinished
	// 完结
	StatusFinished
	// 下架
	StatusDiscontinued
)

const (
	defaultCover = "404.png"
)

type (
	// Novel novel
	Novel struct {
		ID        uint       `gorm:"primary_key" json:"id,omitempty"`
		CreatedAt *time.Time `json:"createdAt,omitempty"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		DeletedAt *time.Time `sql:"index" json:"deletedAt,omitempty"`

		Name         string `json:"name,omitempty" gorm:"type:varchar(100);not null;unique_index:idx_novels_name_author"`
		Author       string `json:"author,omitempty" gorm:"type:varchar(50);not null;unique_index:idx_novels_name_author"`
		Summary      string `json:"summary,omitempty"`
		WordCount    int    `json:"wordCount,omitempty"`
		ChapterCount int    `json:"chapterCount,omitempty"`
		// Grading 自定义分级
		Grading int `json:"grading,omitempty"`
		// Score 动态生成的打分
		Score int `json:"score,omitempty"`
		// Status 状态
		Status int `json:"status,omitempty"`
	}
	// NovelChapter novel chapter
	NovelChapter struct {
		ID        uint       `gorm:"primary_key" json:"id,omitempty"`
		CreatedAt *time.Time `json:"createdAt,omitempty"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		DeletedAt *time.Time `sql:"index" json:"deletedAt,omitempty"`

		BookID    uint   `json:"bookID,omitempty" gorm:"not null;unique_index:idx_novel_chapters_book_id_no"`
		NO        int    `json:"no,omitempty" gorm:"not null;unique_index:idx_novel_chapters_book_id_no"`
		Title     string `json:"title,omitempty"`
		Content   string `json:"content,omitempty"`
		WordCount int    `json:"wordCount,omitempty"`
	}
	// NovelCover novel cover
	NovelCover struct {
		ID        uint       `gorm:"primary_key" json:"id,omitempty"`
		CreatedAt *time.Time `json:"createdAt,omitempty"`
		UpdatedAt *time.Time `json:"updatedAt,omitempty"`
		DeletedAt *time.Time `sql:"index" json:"deletedAt,omitempty"`

		BookID      uint   `json:"bookID,omitempty" gorm:"not null;unique_index:idx_novel_covers_book_id"`
		Data        []byte `json:"data,omitempty"`
		ContentType string `json:"contentType,omitempty"`
	}

	NovelSrv struct{}

	WslSrv struct{}

	// ImageOptimParams image optim params
	ImageOptimParams struct {
		Width  int
		Height int
		Output string
		Effect string
	}

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

	tinyImage struct {
		Data []byte `json:"data,omitempty"`
	}
)

var (
	wslIns  *axios.Instance
	tinyIns *axios.Instance
)

const (
	wslBooksURL        = "/books/v1"
	wslBookChaptersURL = "/books/v1/%d/chapters?limit=%d&offset=%d&fields=id,no,title,word_count,content"
)

func init() {
	pgGetClient().
		AutoMigrate(&Novel{}).
		AutoMigrate(&NovelCover{}).
		AutoMigrate(&NovelChapter{})

	wslConfig := config.GetWslConfig()
	wslIns = axios.NewInstance(&axios.InstanceConfig{
		BaseURL: wslConfig.BaseURL,
		Timeout: 10 * time.Second,
		ResponseInterceptors: []axios.ResponseInterceptor{
			func(resp *axios.Response) (err error) {
				if resp.Status < 400 {
					return nil
				}
				message := standardJSON.Get(resp.Data, "message").ToString()
				if message == "" {
					message = "Unknown Error"
				}
				return &hes.Error{
					Message:    message,
					StatusCode: resp.Status,
				}
			},
		},
	})

	tinyConfig := config.GetTinyConfig()
	tinyIns = axios.NewInstance(&axios.InstanceConfig{
		BaseURL: tinyConfig.BaseURL,
		Timeout: 5 * time.Second,
	})

}

// FindOne find one novel
func (srv *NovelSrv) FindOne(condition *Novel) (*Novel, error) {
	result := new(Novel)
	err := pgGetClient().First(result, condition).Error
	return result, err
}

// List list novel
func (srv *NovelSrv) List(params *helper.DbParams, where ...interface{}) (novels []*Novel, err error) {
	novels = make([]*Novel, 0)
	err = helper.PGGetDB(params).Find(&novels, where...).Error
	return
}

// Count count novel
func (srv *NovelSrv) Count(where ...interface{}) (count int, err error) {
	db := pgGetClient().Model(&Novel{})
	if len(where) != 0 {
		db = db.Where(where[0], where[1:]...)
	}
	err = db.Count(&count).Error
	return
}

// Update update novel
func (srv *NovelSrv) Update(data *Novel) (err error) {
	err = pgGetClient().Model(&Novel{}).Update(data).Error
	return
}

func (srv *NovelSrv) UpdateChapter(data *NovelChapter) (err error) {
	err = pgGetClient().Model(&NovelChapter{}).Update(data).Error
	return
}

// ListChapters list chapters
func (srv *NovelSrv) ListChapters(params *helper.DbParams, where ...interface{}) (chapters []*NovelChapter, err error) {
	chapters = make([]*NovelChapter, 0)
	err = helper.PGGetDB(params).Find(&chapters, where...).Error
	return
}

// UpdateCover update cover
func (srv *NovelSrv) UpdateCover(bookID uint, coverURL string, force bool) (err error) {
	if coverURL == "" {
		return
	}
	cover := &NovelCover{}
	result := pgGetClient().First(cover, &NovelCover{
		BookID: bookID,
	})

	// 如果无图片，则添加
	if result.Error != nil && !result.RecordNotFound() {
		err = result.Error
		return
	}
	// 如果已有封面，且非强制更新
	if cover.ID != 0 && !force {
		return
	}
	resp, err := axios.Get(coverURL)
	if err != nil {
		return
	}
	if resp.Status != http.StatusOK {
		err = errors.New("get cover fail")
		return
	}
	cover.BookID = bookID
	cover.ContentType = resp.Headers.Get("Content-Type")
	cover.Data = resp.Data
	err = pgGetClient().Save(cover).Error
	if err != nil {
		return
	}
	return
}

// GetCover get cover
func (srv *NovelSrv) GetCover(bookID uint, params *ImageOptimParams) (cover *NovelCover, err error) {
	cover = new(NovelCover)
	err = pgGetClient().First(cover, &NovelCover{
		BookID: bookID,
	}).Error
	// 如果是无封面，使用默认封面
	if err == gorm.ErrRecordNotFound {
		err = nil
		buf, _ := GetAssetFile(defaultCover)
		if len(buf) != 0 {
			cover.Data = buf
			cover.ContentType = mime.TypeByExtension(filepath.Ext(defaultCover))
		}
	}
	if cover == nil {
		return
	}
	outputType := params.Output
	if outputType == "" {
		outputType = "jpeg"
	}
	source := strings.Replace(cover.ContentType, "image/", "", 1)
	data := cover.Data
	// 只针对jpeg处理
	if params.Effect == "blur" && source == "jpeg" {
		img, _, _ := image.Decode(bytes.NewReader(data))
		if img != nil {
			rgba := imaging.Blur(img, 10)
			buf := new(bytes.Buffer)
			_ = jpeg.Encode(buf, rgba, nil)
			if buf.Len() != 0 {
				data = buf.Bytes()
			}
		}
	}

	resp, _ := tinyIns.Post("/images/optim", map[string]interface{}{
		"data":    base64.StdEncoding.EncodeToString(data),
		"source":  source,
		"output":  outputType,
		"width":   params.Width,
		"height":  params.Height,
		"quality": 80,
	})
	// 如果压缩失败，则直接使用原数据
	if resp != nil {
		img := new(tinyImage)
		_ = resp.JSON(img)
		if len(img.Data) != 0 {
			cover.Data = img.Data
			cover.ContentType = "image/" + outputType
		}
	}
	return
}

func (srv *NovelSrv) RefreshBasicInfo(novel *Novel) (err error) {
	chapters, err := srv.ListChapters(&helper.DbParams{
		Fields: "word_count",
	}, &NovelChapter{
		BookID: novel.ID,
	})
	if err != nil {
		return
	}
	wordCount := 0
	chapterCount := len(chapters)
	if chapterCount == 0 {
		return
	}
	for _, chapter := range chapters {
		wordCount += chapter.WordCount
	}
	// 如果章节或字数不一致，则更新
	if chapterCount != novel.ChapterCount ||
		wordCount != novel.WordCount {
		novel.ChapterCount = chapterCount
		novel.WordCount = wordCount
		_ = pgGetClient().Model(&Novel{}).Update(novel)
	}
	return
}

// RefreshAllBasicInfo refresh basic info
func (srv *NovelSrv) RefreshAllBasicInfo() (err error) {
	// 计算章节、字数等
	limit := 10
	offset := 0
	for {
		novels, e := srv.List(&helper.DbParams{
			Limit:  limit,
			Offset: offset,
			Fields: "id,word_count,chapter_count",
		}, &Novel{
			Status: StatusUnfinished,
		})
		if e != nil {
			err = e
			return
		}
		for _, novel := range novels {
			_ = srv.RefreshBasicInfo(novel)
		}
		if len(novels) < limit {
			break
		}
		offset += limit
	}
	return
}

// Add add novel
func (srv *NovelSrv) Add(data Novel) (novel *Novel, err error) {
	novel = &data
	err = pgGetClient().Save(novel).Error
	if err != nil {
		return
	}
	return
}

// AddChapter add chapter
func (srv *NovelSrv) AddChapter(data NovelChapter) (chapter *NovelChapter, err error) {
	chapter = &data
	err = pgGetClient().Save(chapter).Error
	if err != nil {
		return
	}
	return
}

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
