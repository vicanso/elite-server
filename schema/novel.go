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

package schema

import (
	"errors"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

const (
	NovelStatusUnknown = iota
	// NovelStatusWritting 连载中
	NovelStatusWritting
	// NovelStatusDone 已完结
	NovelStatusDone
	// NovelStatusBan 禁止状态
	NovelStatusBan
	NovelStatusEnd
)

// Novel holds the schema definition for the Novel entity.
type Novel struct {
	ent.Schema
}

// Novel 小说的mixin
func (Novel) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields of the Novel.
func (Novel) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			Immutable().
			Comment("小说名称"),
		field.String("author").
			NotEmpty().
			Immutable().
			Comment("小说作者"),
		field.Int("source").
			Immutable().
			NonNegative().
			Comment("小说来源"),
		field.Int("status").
			Default(NovelStatusWritting).
			Validate(func(i int) error {
				if i <= NovelStatusUnknown || i >= NovelStatusEnd {
					return errors.New("status is invalid")
				}
				return nil
			}).
			Comment("小说状态"),
		field.Int("chapter_count").
			Default(0).
			// 如果需要select此字段，则需要设置sql
			StructTag(`json:"chapterCount" sql:"chapter_count"`).
			Comment("章节总数"),
		field.Int("word_count").
			Optional().
			Default(0).
			// 如果需要select此字段，则需要设置sql
			StructTag(`json:"wordCount" sql:"word_count"`).
			Comment("小说总字数"),
		field.Int("views").
			Default(0).
			Comment("小说阅读次数"),
		field.Int("downloads").
			Default(0).
			Comment("小说下载次数"),
		field.Int("favorites").
			Default(0).
			Comment("小说收藏次数"),
		// 根据更新频率计算权重
		field.Int("updated_weight").
			Optional().
			Default(0).
			StructTag(`json:"updatedWeight" sql:"updated_weight"`).
			Comment("小说更新权重"),
		field.String("cover").
			Optional().
			Comment("小说封面"),
		field.String("summary").
			Comment("小说简介"),
		// 小说分类
		field.Strings("categories").
			Optional().
			Comment("小说分类"),
	}
}

// Edges of the Novel.
func (Novel) Edges() []ent.Edge {
	return nil
}

func (Novel) Indexes() []ent.Index {
	return []ent.Index{
		// 唯一索引
		index.Fields("name", "author").Unique(),
		index.Fields("categories"),
		index.Fields("views"),
		index.Fields("favorites"),
	}
}
