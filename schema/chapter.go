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
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Chapter holds the schema definition for the Chapter entity.
type Chapter struct {
	ent.Schema
}

// Mixin chapter的mixin
func (Chapter) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields of the Chapter.
func (Chapter) Fields() []ent.Field {
	return []ent.Field{
		field.Int("novel").
			Comment("小说id"),
		field.Int("no").
			Comment("章节序号"),
		field.String("title").
			Comment("章节名称"),
		field.String("content").
			Optional().
			Comment("章节内容"),
		field.Int("word_count").
			Optional().
			StructTag(`json:"wordCount,omitempty" sql:"word_count"`).
			Comment("章节字数"),
	}
}

// Edges of the Chapter.
func (Chapter) Edges() []ent.Edge {
	return nil
}

func (Chapter) Indexes() []ent.Index {
	return []ent.Index{
		// 唯一索引
		index.Fields("novel", "no").Unique(),
	}
}
