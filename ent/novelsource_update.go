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

// NovelSourceUpdate is the builder for updating NovelSource entities.
type NovelSourceUpdate struct {
	config
	hooks      []Hook
	mutation   *NovelSourceMutation
	predicates []predicate.NovelSource
}

// Where adds a new predicate for the builder.
func (nsu *NovelSourceUpdate) Where(ps ...predicate.NovelSource) *NovelSourceUpdate {
	nsu.predicates = append(nsu.predicates, ps...)
	return nsu
}

// SetStatus sets the status field.
func (nsu *NovelSourceUpdate) SetStatus(i int) *NovelSourceUpdate {
	nsu.mutation.ResetStatus()
	nsu.mutation.SetStatus(i)
	return nsu
}

// SetNillableStatus sets the status field if the given value is not nil.
func (nsu *NovelSourceUpdate) SetNillableStatus(i *int) *NovelSourceUpdate {
	if i != nil {
		nsu.SetStatus(*i)
	}
	return nsu
}

// AddStatus adds i to status.
func (nsu *NovelSourceUpdate) AddStatus(i int) *NovelSourceUpdate {
	nsu.mutation.AddStatus(i)
	return nsu
}

// Mutation returns the NovelSourceMutation object of the builder.
func (nsu *NovelSourceUpdate) Mutation() *NovelSourceMutation {
	return nsu.mutation
}

// Save executes the query and returns the number of rows/vertices matched by this operation.
func (nsu *NovelSourceUpdate) Save(ctx context.Context) (int, error) {
	var (
		err      error
		affected int
	)
	nsu.defaults()
	if len(nsu.hooks) == 0 {
		if err = nsu.check(); err != nil {
			return 0, err
		}
		affected, err = nsu.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*NovelSourceMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = nsu.check(); err != nil {
				return 0, err
			}
			nsu.mutation = mutation
			affected, err = nsu.sqlSave(ctx)
			mutation.done = true
			return affected, err
		})
		for i := len(nsu.hooks) - 1; i >= 0; i-- {
			mut = nsu.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, nsu.mutation); err != nil {
			return 0, err
		}
	}
	return affected, err
}

// SaveX is like Save, but panics if an error occurs.
func (nsu *NovelSourceUpdate) SaveX(ctx context.Context) int {
	affected, err := nsu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (nsu *NovelSourceUpdate) Exec(ctx context.Context) error {
	_, err := nsu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (nsu *NovelSourceUpdate) ExecX(ctx context.Context) {
	if err := nsu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (nsu *NovelSourceUpdate) defaults() {
	if _, ok := nsu.mutation.UpdatedAt(); !ok {
		v := novelsource.UpdateDefaultUpdatedAt()
		nsu.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (nsu *NovelSourceUpdate) check() error {
	if v, ok := nsu.mutation.Status(); ok {
		if err := novelsource.StatusValidator(v); err != nil {
			return &ValidationError{Name: "status", err: fmt.Errorf("ent: validator failed for field \"status\": %w", err)}
		}
	}
	return nil
}

func (nsu *NovelSourceUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   novelsource.Table,
			Columns: novelsource.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: novelsource.FieldID,
			},
		},
	}
	if ps := nsu.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := nsu.mutation.UpdatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: novelsource.FieldUpdatedAt,
		})
	}
	if value, ok := nsu.mutation.Status(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: novelsource.FieldStatus,
		})
	}
	if value, ok := nsu.mutation.AddedStatus(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: novelsource.FieldStatus,
		})
	}
	if n, err = sqlgraph.UpdateNodes(ctx, nsu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{novelsource.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return 0, err
	}
	return n, nil
}

// NovelSourceUpdateOne is the builder for updating a single NovelSource entity.
type NovelSourceUpdateOne struct {
	config
	hooks    []Hook
	mutation *NovelSourceMutation
}

// SetStatus sets the status field.
func (nsuo *NovelSourceUpdateOne) SetStatus(i int) *NovelSourceUpdateOne {
	nsuo.mutation.ResetStatus()
	nsuo.mutation.SetStatus(i)
	return nsuo
}

// SetNillableStatus sets the status field if the given value is not nil.
func (nsuo *NovelSourceUpdateOne) SetNillableStatus(i *int) *NovelSourceUpdateOne {
	if i != nil {
		nsuo.SetStatus(*i)
	}
	return nsuo
}

// AddStatus adds i to status.
func (nsuo *NovelSourceUpdateOne) AddStatus(i int) *NovelSourceUpdateOne {
	nsuo.mutation.AddStatus(i)
	return nsuo
}

// Mutation returns the NovelSourceMutation object of the builder.
func (nsuo *NovelSourceUpdateOne) Mutation() *NovelSourceMutation {
	return nsuo.mutation
}

// Save executes the query and returns the updated entity.
func (nsuo *NovelSourceUpdateOne) Save(ctx context.Context) (*NovelSource, error) {
	var (
		err  error
		node *NovelSource
	)
	nsuo.defaults()
	if len(nsuo.hooks) == 0 {
		if err = nsuo.check(); err != nil {
			return nil, err
		}
		node, err = nsuo.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*NovelSourceMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = nsuo.check(); err != nil {
				return nil, err
			}
			nsuo.mutation = mutation
			node, err = nsuo.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(nsuo.hooks) - 1; i >= 0; i-- {
			mut = nsuo.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, nsuo.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX is like Save, but panics if an error occurs.
func (nsuo *NovelSourceUpdateOne) SaveX(ctx context.Context) *NovelSource {
	node, err := nsuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (nsuo *NovelSourceUpdateOne) Exec(ctx context.Context) error {
	_, err := nsuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (nsuo *NovelSourceUpdateOne) ExecX(ctx context.Context) {
	if err := nsuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (nsuo *NovelSourceUpdateOne) defaults() {
	if _, ok := nsuo.mutation.UpdatedAt(); !ok {
		v := novelsource.UpdateDefaultUpdatedAt()
		nsuo.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (nsuo *NovelSourceUpdateOne) check() error {
	if v, ok := nsuo.mutation.Status(); ok {
		if err := novelsource.StatusValidator(v); err != nil {
			return &ValidationError{Name: "status", err: fmt.Errorf("ent: validator failed for field \"status\": %w", err)}
		}
	}
	return nil
}

func (nsuo *NovelSourceUpdateOne) sqlSave(ctx context.Context) (_node *NovelSource, err error) {
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   novelsource.Table,
			Columns: novelsource.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: novelsource.FieldID,
			},
		},
	}
	id, ok := nsuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "ID", err: fmt.Errorf("missing NovelSource.ID for update")}
	}
	_spec.Node.ID.Value = id
	if value, ok := nsuo.mutation.UpdatedAt(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: novelsource.FieldUpdatedAt,
		})
	}
	if value, ok := nsuo.mutation.Status(); ok {
		_spec.Fields.Set = append(_spec.Fields.Set, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: novelsource.FieldStatus,
		})
	}
	if value, ok := nsuo.mutation.AddedStatus(); ok {
		_spec.Fields.Add = append(_spec.Fields.Add, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: novelsource.FieldStatus,
		})
	}
	_node = &NovelSource{config: nsuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues()
	if err = sqlgraph.UpdateNode(ctx, nsuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{novelsource.Label}
		} else if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	return _node, nil
}
