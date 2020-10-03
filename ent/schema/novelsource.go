package schema

import (
	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/field"
	"github.com/facebook/ent/schema/index"
)

// NovelSource holds the schema definition for the NovelSource entity.
type NovelSource struct {
	ent.Schema
}

// NovelSource 小说来源的mixin
func (NovelSource) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields of the NovelSource.
func (NovelSource) Fields() []ent.Field {
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
		field.Int("source_id").
			StructTag(`json:"sourceID,omitempty"`).
			NonNegative().
			Immutable().
			Comment("小说来源ID"),
	}
}

// Edges of the NovelSource.
func (NovelSource) Edges() []ent.Edge {
	return nil
}

// Indexes 小说来源索引
func (NovelSource) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name", "author"),
		index.Fields("source", "source_id").Unique(),
	}
}
