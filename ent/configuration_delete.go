// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"

	"github.com/facebook/ent/dialect/sql"
	"github.com/facebook/ent/dialect/sql/sqlgraph"
	"github.com/facebook/ent/schema/field"
	"github.com/vicanso/elite/ent/configuration"
	"github.com/vicanso/elite/ent/predicate"
)

// ConfigurationDelete is the builder for deleting a Configuration entity.
type ConfigurationDelete struct {
	config
	hooks    []Hook
	mutation *ConfigurationMutation
}

// Where adds a new predicate to the ConfigurationDelete builder.
func (cd *ConfigurationDelete) Where(ps ...predicate.Configuration) *ConfigurationDelete {
	cd.mutation.predicates = append(cd.mutation.predicates, ps...)
	return cd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (cd *ConfigurationDelete) Exec(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	if len(cd.hooks) == 0 {
		affected, err = cd.sqlExec(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*ConfigurationMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			cd.mutation = mutation
			affected, err = cd.sqlExec(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(cd.hooks) - 1; i >= 0; i-- {
			mut = cd.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, cd.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// ExecX is like Exec, but panics if an error occurs.
func (cd *ConfigurationDelete) ExecX(ctx context.Context) int {
	n, err := cd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (cd *ConfigurationDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := &sqlgraph.DeleteSpec{
		Node: &sqlgraph.NodeSpec{
			Table: configuration.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: configuration.FieldID,
			},
		},
	}
	if ps := cd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return sqlgraph.DeleteNodes(ctx, cd.driver, _spec)
}

// ConfigurationDeleteOne is the builder for deleting a single Configuration entity.
type ConfigurationDeleteOne struct {
	cd *ConfigurationDelete
}

// Exec executes the deletion query.
func (cdo *ConfigurationDeleteOne) Exec(ctx context.Context) error {
	n, err := cdo.cd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{configuration.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (cdo *ConfigurationDeleteOne) ExecX(ctx context.Context) {
	cdo.cd.ExecX(ctx)
}
