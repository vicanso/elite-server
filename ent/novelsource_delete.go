// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"github.com/facebook/ent/dialect/sql"
	"github.com/facebook/ent/dialect/sql/sqlgraph"
	"github.com/facebook/ent/schema/field"
	"github.com/vicanso/elite/ent/novelsource"
	"github.com/vicanso/elite/ent/predicate"
)

// NovelSourceDelete is the builder for deleting a NovelSource entity.
type NovelSourceDelete struct {
	config
	hooks    []Hook
	mutation *NovelSourceMutation
}

// Where adds a new predicate to the NovelSourceDelete builder.
func (nsd *NovelSourceDelete) Where(ps ...predicate.NovelSource) *NovelSourceDelete {
	nsd.mutation.predicates = append(nsd.mutation.predicates, ps...)
	return nsd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (nsd *NovelSourceDelete) Exec(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(nsd.hooks) == 0 {
		affected, err = nsd.sqlExec(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*NovelSourceMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			nsd.mutation = mutation
			affected, err = nsd.sqlExec(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(nsd.hooks) - 1; i >= 0; i-- {
			mut = nsd.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, nsd.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// ExecX is like Exec, but panics if an error occurs.
func (nsd *NovelSourceDelete) ExecX(ctx context.Context) int {
	n, err := nsd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (nsd *NovelSourceDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: novelsource.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: novelsource.FieldID,
			},
		},
	}
	if ps := nsd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, nsd.driver, _spec)
}

// NovelSourceDeleteOne is the builder for deleting a single NovelSource entity.
type NovelSourceDeleteOne struct {
	nsd *NovelSourceDelete
}

// Exec executes the deletion query.
func (nsdo *NovelSourceDeleteOne) Exec(ctx context.Context) error {
	n, err := nsdo.nsd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{novelsource.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (nsdo *NovelSourceDeleteOne) ExecX(ctx context.Context) {
	nsdo.nsd.ExecX(ctx)
}
