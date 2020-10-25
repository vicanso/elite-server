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
		field.Int("word_count").
			Optional().
			Default(0).
			// 如果需要select此字段，则需要设置sql
			StructTag(`json:"wordCount,omitempty" sql:"word_count"`).
			Comment("小说总字数"),
		field.Int("views").
			Optional().
			Default(0).
			Comment("小说阅读次数"),
		field.Int("downloads").
			Optional().
			Default(0).
			Comment("小说下载次数"),
		field.Int("favorites").
			Optional().
			Default(0).
			Comment("小说收藏次数"),
		// 根据更新频率计算权重
		field.Int("updated_weight").
			Optional().
			Default(0).
			StructTag(`json:"updatedWeight,omitempty" sql:"updated_weight"`).
			Comment("小说更新权重"),
		field.String("cover").
			Optional().
			Comment("小说封面"),
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
