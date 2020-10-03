package schema

import (
	"errors"

	"github.com/facebook/ent"
	"github.com/facebook/ent/schema/field"
	"github.com/facebook/ent/schema/index"
)

const (
	NovelStatusUnknown = iota
	// NovelStatusWritting 连载中
	NovelStatusWritting
	// NovelStatusDone 已完结
	NovelStatusDone
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
		field.String("summary").
			Comment("小说简介"),
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
	}
}
