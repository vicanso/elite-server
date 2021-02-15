// Code generated by entc, DO NOT EDIT.

package novelsource

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/vicanso/elite/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldID), id))
	})
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldID), id))
	})
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
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
func IDNotIn(ids ...int) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
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
func IDGT(id int) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldID), id))
	})
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldID), id))
	})
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldID), id))
	})
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldID), id))
	})
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreatedAt), v))
	})
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdatedAt), v))
	})
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldName), v))
	})
}

// Author applies equality check predicate on the "author" field. It's identical to AuthorEQ.
func Author(v string) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldAuthor), v))
	})
}

// Source applies equality check predicate on the "source" field. It's identical to SourceEQ.
func Source(v int) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldSource), v))
	})
}

// SourceID applies equality check predicate on the "source_id" field. It's identical to SourceIDEQ.
func SourceID(v int) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldSourceID), v))
	})
}

// Status applies equality check predicate on the "status" field. It's identical to StatusEQ.
func Status(v int) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldStatus), v))
	})
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.NovelSource {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.NovelSource(func(s *sql.Selector) {
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
func CreatedAtNotIn(vs ...time.Time) predicate.NovelSource {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.NovelSource(func(s *sql.Selector) {
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
func CreatedAtGT(v time.Time) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldCreatedAt), v))
	})
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldCreatedAt), v))
	})
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldUpdatedAt), v))
	})
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldUpdatedAt), v))
	})
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.NovelSource {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.NovelSource(func(s *sql.Selector) {
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
func UpdatedAtNotIn(vs ...time.Time) predicate.NovelSource {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.NovelSource(func(s *sql.Selector) {
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
func UpdatedAtGT(v time.Time) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldUpdatedAt), v))
	})
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldUpdatedAt), v))
	})
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldUpdatedAt), v))
	})
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldUpdatedAt), v))
	})
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldName), v))
	})
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldName), v))
	})
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.NovelSource {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.NovelSource(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldName), v...))
	})
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.NovelSource {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.NovelSource(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldName), v...))
	})
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldName), v))
	})
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldName), v))
	})
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldName), v))
	})
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldName), v))
	})
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldName), v))
	})
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldName), v))
	})
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldName), v))
	})
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldName), v))
	})
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldName), v))
	})
}

// AuthorEQ applies the EQ predicate on the "author" field.
func AuthorEQ(v string) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldAuthor), v))
	})
}

// AuthorNEQ applies the NEQ predicate on the "author" field.
func AuthorNEQ(v string) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldAuthor), v))
	})
}

// AuthorIn applies the In predicate on the "author" field.
func AuthorIn(vs ...string) predicate.NovelSource {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.NovelSource(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldAuthor), v...))
	})
}

// AuthorNotIn applies the NotIn predicate on the "author" field.
func AuthorNotIn(vs ...string) predicate.NovelSource {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.NovelSource(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldAuthor), v...))
	})
}

// AuthorGT applies the GT predicate on the "author" field.
func AuthorGT(v string) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldAuthor), v))
	})
}

// AuthorGTE applies the GTE predicate on the "author" field.
func AuthorGTE(v string) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldAuthor), v))
	})
}

// AuthorLT applies the LT predicate on the "author" field.
func AuthorLT(v string) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldAuthor), v))
	})
}

// AuthorLTE applies the LTE predicate on the "author" field.
func AuthorLTE(v string) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldAuthor), v))
	})
}

// AuthorContains applies the Contains predicate on the "author" field.
func AuthorContains(v string) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.Contains(s.C(FieldAuthor), v))
	})
}

// AuthorHasPrefix applies the HasPrefix predicate on the "author" field.
func AuthorHasPrefix(v string) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.HasPrefix(s.C(FieldAuthor), v))
	})
}

// AuthorHasSuffix applies the HasSuffix predicate on the "author" field.
func AuthorHasSuffix(v string) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.HasSuffix(s.C(FieldAuthor), v))
	})
}

// AuthorEqualFold applies the EqualFold predicate on the "author" field.
func AuthorEqualFold(v string) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.EqualFold(s.C(FieldAuthor), v))
	})
}

// AuthorContainsFold applies the ContainsFold predicate on the "author" field.
func AuthorContainsFold(v string) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.ContainsFold(s.C(FieldAuthor), v))
	})
}

// SourceEQ applies the EQ predicate on the "source" field.
func SourceEQ(v int) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldSource), v))
	})
}

// SourceNEQ applies the NEQ predicate on the "source" field.
func SourceNEQ(v int) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldSource), v))
	})
}

// SourceIn applies the In predicate on the "source" field.
func SourceIn(vs ...int) predicate.NovelSource {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.NovelSource(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldSource), v...))
	})
}

// SourceNotIn applies the NotIn predicate on the "source" field.
func SourceNotIn(vs ...int) predicate.NovelSource {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.NovelSource(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldSource), v...))
	})
}

// SourceGT applies the GT predicate on the "source" field.
func SourceGT(v int) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldSource), v))
	})
}

// SourceGTE applies the GTE predicate on the "source" field.
func SourceGTE(v int) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldSource), v))
	})
}

// SourceLT applies the LT predicate on the "source" field.
func SourceLT(v int) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldSource), v))
	})
}

// SourceLTE applies the LTE predicate on the "source" field.
func SourceLTE(v int) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldSource), v))
	})
}

// SourceIDEQ applies the EQ predicate on the "source_id" field.
func SourceIDEQ(v int) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldSourceID), v))
	})
}

// SourceIDNEQ applies the NEQ predicate on the "source_id" field.
func SourceIDNEQ(v int) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldSourceID), v))
	})
}

// SourceIDIn applies the In predicate on the "source_id" field.
func SourceIDIn(vs ...int) predicate.NovelSource {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.NovelSource(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldSourceID), v...))
	})
}

// SourceIDNotIn applies the NotIn predicate on the "source_id" field.
func SourceIDNotIn(vs ...int) predicate.NovelSource {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.NovelSource(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldSourceID), v...))
	})
}

// SourceIDGT applies the GT predicate on the "source_id" field.
func SourceIDGT(v int) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldSourceID), v))
	})
}

// SourceIDGTE applies the GTE predicate on the "source_id" field.
func SourceIDGTE(v int) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldSourceID), v))
	})
}

// SourceIDLT applies the LT predicate on the "source_id" field.
func SourceIDLT(v int) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldSourceID), v))
	})
}

// SourceIDLTE applies the LTE predicate on the "source_id" field.
func SourceIDLTE(v int) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldSourceID), v))
	})
}

// StatusEQ applies the EQ predicate on the "status" field.
func StatusEQ(v int) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.EQ(s.C(FieldStatus), v))
	})
}

// StatusNEQ applies the NEQ predicate on the "status" field.
func StatusNEQ(v int) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.NEQ(s.C(FieldStatus), v))
	})
}

// StatusIn applies the In predicate on the "status" field.
func StatusIn(vs ...int) predicate.NovelSource {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.NovelSource(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.In(s.C(FieldStatus), v...))
	})
}

// StatusNotIn applies the NotIn predicate on the "status" field.
func StatusNotIn(vs ...int) predicate.NovelSource {
	v := make([]interface{}, len(vs))
	for i := range v {
		v[i] = vs[i]
	}
	return predicate.NovelSource(func(s *sql.Selector) {
		// if not arguments were provided, append the FALSE constants,
		// since we can't apply "IN ()". This will make this predicate falsy.
		if len(v) == 0 {
			s.Where(sql.False())
			return
		}
		s.Where(sql.NotIn(s.C(FieldStatus), v...))
	})
}

// StatusGT applies the GT predicate on the "status" field.
func StatusGT(v int) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.GT(s.C(FieldStatus), v))
	})
}

// StatusGTE applies the GTE predicate on the "status" field.
func StatusGTE(v int) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.GTE(s.C(FieldStatus), v))
	})
}

// StatusLT applies the LT predicate on the "status" field.
func StatusLT(v int) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.LT(s.C(FieldStatus), v))
	})
}

// StatusLTE applies the LTE predicate on the "status" field.
func StatusLTE(v int) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s.Where(sql.LTE(s.C(FieldStatus), v))
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.NovelSource) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.NovelSource) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
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
func Not(p predicate.NovelSource) predicate.NovelSource {
	return predicate.NovelSource(func(s *sql.Selector) {
		p(s.Not())
	})
}
