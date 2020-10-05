// Code generated by entc, DO NOT EDIT.

package chapter

import (
	"time"

	"github.com/facebook/ent/dialect/sql"
	"github.com/vicanso/elite/ent/predicate"
)

// ID filters vertices based on their identifier.
func ID(id int) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(ids) == 0 {
			s.Where(sql.False())
			return
		}
		v := make([]interface{}, len(ids))
		for i := range v {
			v[i] = ids[i]
		}
		s.Where(sql.In(s.C(FieldID), v...))
	})
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(ids) == 0 {
			s.Where(sql.False())
			return
		}
		v := make([]interface{}, len(ids))
		for i := range v {
			v[i] = ids[i]
		}
		s.Where(sql.NotIn(s.C(FieldID), v...))
	})
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreatedAt), v))
	})
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdatedAt), v))
	})
}

// Novel applies equality check predicate on the "novel" field. It's identical to NovelEQ.
func Novel(v int) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldNovel), v))
	})
}

// No applies equality check predicate on the "no" field. It's identical to NoEQ.
func No(v int) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldNo), v))
	})
}

// Title applies equality check predicate on the "title" field. It's identical to TitleEQ.
func Title(v string) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldTitle), v))
	})
}

// Content applies equality check predicate on the "content" field. It's identical to ContentEQ.
func Content(v string) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldContent), v))
	})
}

// WordCount applies equality check predicate on the "word_count" field. It's identical to WordCountEQ.
func WordCount(v int) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldWordCount), v))
	})
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.Chapter {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Chapter(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldCreatedAt), v...))
	})
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.Chapter {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Chapter(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldCreatedAt), v...))
	})
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCreatedAt), v))
	})
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdatedAt), v))
	})
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldUpdatedAt), v))
	})
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.Chapter {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Chapter(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldUpdatedAt), v...))
	})
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.Chapter {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Chapter(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldUpdatedAt), v...))
	})
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v time.Time) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldUpdatedAt), v))
	})
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldUpdatedAt), v))
	})
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldUpdatedAt), v))
	})
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldUpdatedAt), v))
	})
}

// NovelEQ applies the EQ predicate on the "novel" field.
func NovelEQ(v int) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldNovel), v))
	})
}

// NovelNEQ applies the NEQ predicate on the "novel" field.
func NovelNEQ(v int) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldNovel), v))
	})
}

// NovelIn applies the In predicate on the "novel" field.
func NovelIn(vs ...int) predicate.Chapter {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Chapter(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldNovel), v...))
	})
}

// NovelNotIn applies the NotIn predicate on the "novel" field.
func NovelNotIn(vs ...int) predicate.Chapter {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Chapter(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldNovel), v...))
	})
}

// NovelGT applies the GT predicate on the "novel" field.
func NovelGT(v int) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldNovel), v))
	})
}

// NovelGTE applies the GTE predicate on the "novel" field.
func NovelGTE(v int) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldNovel), v))
	})
}

// NovelLT applies the LT predicate on the "novel" field.
func NovelLT(v int) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldNovel), v))
	})
}

// NovelLTE applies the LTE predicate on the "novel" field.
func NovelLTE(v int) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldNovel), v))
	})
}

// NoEQ applies the EQ predicate on the "no" field.
func NoEQ(v int) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldNo), v))
	})
}

// NoNEQ applies the NEQ predicate on the "no" field.
func NoNEQ(v int) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldNo), v))
	})
}

// NoIn applies the In predicate on the "no" field.
func NoIn(vs ...int) predicate.Chapter {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Chapter(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldNo), v...))
	})
}

// NoNotIn applies the NotIn predicate on the "no" field.
func NoNotIn(vs ...int) predicate.Chapter {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Chapter(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldNo), v...))
	})
}

// NoGT applies the GT predicate on the "no" field.
func NoGT(v int) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldNo), v))
	})
}

// NoGTE applies the GTE predicate on the "no" field.
func NoGTE(v int) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldNo), v))
	})
}

// NoLT applies the LT predicate on the "no" field.
func NoLT(v int) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldNo), v))
	})
}

// NoLTE applies the LTE predicate on the "no" field.
func NoLTE(v int) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldNo), v))
	})
}

// TitleEQ applies the EQ predicate on the "title" field.
func TitleEQ(v string) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldTitle), v))
	})
}

// TitleNEQ applies the NEQ predicate on the "title" field.
func TitleNEQ(v string) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldTitle), v))
	})
}

// TitleIn applies the In predicate on the "title" field.
func TitleIn(vs ...string) predicate.Chapter {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Chapter(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldTitle), v...))
	})
}

// TitleNotIn applies the NotIn predicate on the "title" field.
func TitleNotIn(vs ...string) predicate.Chapter {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Chapter(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldTitle), v...))
	})
}

// TitleGT applies the GT predicate on the "title" field.
func TitleGT(v string) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldTitle), v))
	})
}

// TitleGTE applies the GTE predicate on the "title" field.
func TitleGTE(v string) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldTitle), v))
	})
}

// TitleLT applies the LT predicate on the "title" field.
func TitleLT(v string) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldTitle), v))
	})
}

// TitleLTE applies the LTE predicate on the "title" field.
func TitleLTE(v string) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldTitle), v))
	})
}

// TitleContains applies the Contains predicate on the "title" field.
func TitleContains(v string) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldTitle), v))
	})
}

// TitleHasPrefix applies the HasPrefix predicate on the "title" field.
func TitleHasPrefix(v string) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldTitle), v))
	})
}

// TitleHasSuffix applies the HasSuffix predicate on the "title" field.
func TitleHasSuffix(v string) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldTitle), v))
	})
}

// TitleEqualFold applies the EqualFold predicate on the "title" field.
func TitleEqualFold(v string) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldTitle), v))
	})
}

// TitleContainsFold applies the ContainsFold predicate on the "title" field.
func TitleContainsFold(v string) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldTitle), v))
	})
}

// ContentEQ applies the EQ predicate on the "content" field.
func ContentEQ(v string) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldContent), v))
	})
}

// ContentNEQ applies the NEQ predicate on the "content" field.
func ContentNEQ(v string) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldContent), v))
	})
}

// ContentIn applies the In predicate on the "content" field.
func ContentIn(vs ...string) predicate.Chapter {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Chapter(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldContent), v...))
	})
}

// ContentNotIn applies the NotIn predicate on the "content" field.
func ContentNotIn(vs ...string) predicate.Chapter {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Chapter(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldContent), v...))
	})
}

// ContentGT applies the GT predicate on the "content" field.
func ContentGT(v string) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldContent), v))
	})
}

// ContentGTE applies the GTE predicate on the "content" field.
func ContentGTE(v string) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldContent), v))
	})
}

// ContentLT applies the LT predicate on the "content" field.
func ContentLT(v string) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldContent), v))
	})
}

// ContentLTE applies the LTE predicate on the "content" field.
func ContentLTE(v string) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldContent), v))
	})
}

// ContentContains applies the Contains predicate on the "content" field.
func ContentContains(v string) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldContent), v))
	})
}

// ContentHasPrefix applies the HasPrefix predicate on the "content" field.
func ContentHasPrefix(v string) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldContent), v))
	})
}

// ContentHasSuffix applies the HasSuffix predicate on the "content" field.
func ContentHasSuffix(v string) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldContent), v))
	})
}

// ContentIsNil applies the IsNil predicate on the "content" field.
func ContentIsNil() predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldContent)))
	})
}

// ContentNotNil applies the NotNil predicate on the "content" field.
func ContentNotNil() predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldContent)))
	})
}

// ContentEqualFold applies the EqualFold predicate on the "content" field.
func ContentEqualFold(v string) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldContent), v))
	})
}

// ContentContainsFold applies the ContainsFold predicate on the "content" field.
func ContentContainsFold(v string) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldContent), v))
	})
}

// WordCountEQ applies the EQ predicate on the "word_count" field.
func WordCountEQ(v int) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldWordCount), v))
	})
}

// WordCountNEQ applies the NEQ predicate on the "word_count" field.
func WordCountNEQ(v int) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldWordCount), v))
	})
}

// WordCountIn applies the In predicate on the "word_count" field.
func WordCountIn(vs ...int) predicate.Chapter {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Chapter(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldWordCount), v...))
	})
}

// WordCountNotIn applies the NotIn predicate on the "word_count" field.
func WordCountNotIn(vs ...int) predicate.Chapter {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.Chapter(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldWordCount), v...))
	})
}

// WordCountGT applies the GT predicate on the "word_count" field.
func WordCountGT(v int) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldWordCount), v))
	})
}

// WordCountGTE applies the GTE predicate on the "word_count" field.
func WordCountGTE(v int) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldWordCount), v))
	})
}

// WordCountLT applies the LT predicate on the "word_count" field.
func WordCountLT(v int) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldWordCount), v))
	})
}

// WordCountLTE applies the LTE predicate on the "word_count" field.
func WordCountLTE(v int) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldWordCount), v))
	})
}

// WordCountIsNil applies the IsNil predicate on the "word_count" field.
func WordCountIsNil() predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.IsNull(s.C(FieldWordCount)))
	})
}

// WordCountNotNil applies the NotNil predicate on the "word_count" field.
func WordCountNotNil() predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s.Where(sql.NotNull(s.C(FieldWordCount)))
	})
}

// And groups list of predicates with the AND operator between them.
func And(predicates ...predicate.Chapter) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups list of predicates with the OR operator between them.
func Or(predicates ...predicate.Chapter) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for i, p := range predicates {
			if i > 0 {
				s1.Or()
			}
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Chapter) predicate.Chapter {
	return predicate.Chapter(func(s *sql.Selector) {
		p(s.Not())
	})
}