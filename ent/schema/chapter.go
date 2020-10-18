package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/field"
	"github.com/facebook/ent/schema/index"
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
