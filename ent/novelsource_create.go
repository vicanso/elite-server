// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/facebook/ent/dialect/sql/sqlgraph"
	"github.com/facebook/ent/schema/field"
	"github.com/vicanso/elite/ent/novelsource"
)

// NovelSourceCreate is the builder for creating a NovelSource entity.
type NovelSourceCreate struct {
	config
	mutation *NovelSourceMutation
	hooks    []Hook
}

// SetCreatedAt sets the created_at field.
func (nsc *NovelSourceCreate) SetCreatedAt(t time.Time) *NovelSourceCreate {
	nsc.mutation.SetCreatedAt(t)
	return nsc
}

// SetNillableCreatedAt sets the created_at field if the given value is not nil.
func (nsc *NovelSourceCreate) SetNillableCreatedAt(t *time.Time) *NovelSourceCreate {
	if t != nil {
		nsc.SetCreatedAt(*t)
	}
	return nsc
}

// SetUpdatedAt sets the updated_at field.
func (nsc *NovelSourceCreate) SetUpdatedAt(t time.Time) *NovelSourceCreate {
	nsc.mutation.SetUpdatedAt(t)
	return nsc
}

// SetNillableUpdatedAt sets the updated_at field if the given value is not nil.
func (nsc *NovelSourceCreate) SetNillableUpdatedAt(t *time.Time) *NovelSourceCreate {
	if t != nil {
		nsc.SetUpdatedAt(*t)
	}
	return nsc
}

// SetName sets the name field.
func (nsc *NovelSourceCreate) SetName(s string) *NovelSourceCreate {
	nsc.mutation.SetName(s)
	return nsc
}

// SetAuthor sets the author field.
func (nsc *NovelSourceCreate) SetAuthor(s string) *NovelSourceCreate {
	nsc.mutation.SetAuthor(s)
	return nsc
}

// SetSource sets the source field.
func (nsc *NovelSourceCreate) SetSource(i int) *NovelSourceCreate {
	nsc.mutation.SetSource(i)
	return nsc
}

// SetSourceID sets the source_id field.
func (nsc *NovelSourceCreate) SetSourceID(i int) *NovelSourceCreate {
	nsc.mutation.SetSourceID(i)
	return nsc
}

// SetStatus sets the status field.
func (nsc *NovelSourceCreate) SetStatus(i int) *NovelSourceCreate {
	nsc.mutation.SetStatus(i)
	return nsc
}

// SetNillableStatus sets the status field if the given value is not nil.
func (nsc *NovelSourceCreate) SetNillableStatus(i *int) *NovelSourceCreate {
	if i != nil {
		nsc.SetStatus(*i)
	}
	return nsc
}

// Mutation returns the NovelSourceMutation object of the builder.
func (nsc *NovelSourceCreate) Mutation() *NovelSourceMutation {
	return nsc.mutation
}

// Save creates the NovelSource in the database.
func (nsc *NovelSourceCreate) Save(ctx context.Context) (*NovelSource, error) {
	var (
		err  error
		node *NovelSource
	)
	nsc.defaults()
	if len(nsc.hooks) == 0 {
		if err = nsc.check(); err != nil {
			return nil, err
		}
		node, err = nsc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*NovelSourceMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = nsc.check(); err != nil {
				return nil, err
			}
			nsc.mutation = mutation
			node, err = nsc.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(nsc.hooks) - 1; i >= 0; i-- {
			mut = nsc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, nsc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (nsc *NovelSourceCreate) SaveX(ctx context.Context) *NovelSource {
	v, err := nsc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// defaults sets the default values of the builder before save.
func (nsc *NovelSourceCreate) defaults() {
	if _, ok := nsc.mutation.CreatedAt(); !ok {
		v := novelsource.DefaultCreatedAt()
		nsc.mutation.SetCreatedAt(v)
	}
	if _, ok := nsc.mutation.UpdatedAt(); !ok {
		v := novelsource.DefaultUpdatedAt()
		nsc.mutation.SetUpdatedAt(v)
	}
	if _, ok := nsc.mutation.Status(); !ok {
		v := novelsource.DefaultStatus
		nsc.mutation.SetStatus(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (nsc *NovelSourceCreate) check() error {
	if _, ok := nsc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New("ent: missing required field \"created_at\"")}
	}
	if _, ok := nsc.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New("ent: missing required field \"updated_at\"")}
	}
	if _, ok := nsc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New("ent: missing required field \"name\"")}
	}
	if v, ok := nsc.mutation.Name(); ok {
		if err := novelsource.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf("ent: validator failed for field \"name\": %w", err)}
		}
	}
	if _, ok := nsc.mutation.Author(); !ok {
		return &ValidationError{Name: "author", err: errors.New("ent: missing required field \"author\"")}
	}
	if v, ok := nsc.mutation.Author(); ok {
		if err := novelsource.AuthorValidator(v); err != nil {
			return &ValidationError{Name: "author", err: fmt.Errorf("ent: validator failed for field \"author\": %w", err)}
		}
	}
	if _, ok := nsc.mutation.Source(); !ok {
		return &ValidationError{Name: "source", err: errors.New("ent: missing required field \"source\"")}
	}
	if v, ok := nsc.mutation.Source(); ok {
		if err := novelsource.SourceValidator(v); err != nil {
			return &ValidationError{Name: "source", err: fmt.Errorf("ent: validator failed for field \"source\": %w", err)}
		}
	}
	if _, ok := nsc.mutation.SourceID(); !ok {
		return &ValidationError{Name: "source_id", err: errors.New("ent: missing required field \"source_id\"")}
	}
	if v, ok := nsc.mutation.SourceID(); ok {
		if err := novelsource.SourceIDValidator(v); err != nil {
			return &ValidationError{Name: "source_id", err: fmt.Errorf("ent: validator failed for field \"source_id\": %w", err)}
		}
	}
	if _, ok := nsc.mutation.Status(); !ok {
		return &ValidationError{Name: "status", err: errors.New("ent: missing required field \"status\"")}
	}
	if v, ok := nsc.mutation.Status(); ok {
		if err := novelsource.StatusValidator(v); err != nil {
			return &ValidationError{Name: "status", err: fmt.Errorf("ent: validator failed for field \"status\": %w", err)}
		}
	}
	return nil
}

func (nsc *NovelSourceCreate) sqlSave(ctx context.Context) (*NovelSource, error) {
	_node, _spec := nsc.createSpec()
	if err := sqlgraph.CreateNode(ctx, nsc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	return _node, nil
}

func (nsc *NovelSourceCreate) createSpec() (*NovelSource, *sqlgraph.CreateSpec) {
	var (
		_node = &NovelSource{config: nsc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: novelsource.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: novelsource.FieldID,
			},
		}
	)
	if value, ok := nsc.mutation.CreatedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: novelsource.FieldCreatedAt,
		})
		_node.CreatedAt = value
	}
	if value, ok := nsc.mutation.UpdatedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: novelsource.FieldUpdatedAt,
		})
		_node.UpdatedAt = value
	}
	if value, ok := nsc.mutation.Name(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: novelsource.FieldName,
		})
		_node.Name = value
	}
	if value, ok := nsc.mutation.Author(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: novelsource.FieldAuthor,
		})
		_node.Author = value
	}
	if value, ok := nsc.mutation.Source(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: novelsource.FieldSource,
		})
		_node.Source = value
	}
	if value, ok := nsc.mutation.SourceID(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: novelsource.FieldSourceID,
		})
		_node.SourceID = value
	}
	if value, ok := nsc.mutation.Status(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: novelsource.FieldStatus,
		})
		_node.Status = value
	}
	return _node, _spec
}

// NovelSourceCreateBulk is the builder for creating a bulk of NovelSource entities.
type NovelSourceCreateBulk struct {
	config
	builders []*NovelSourceCreate
}

// Save creates the NovelSource entities in the database.
func (nscb *NovelSourceCreateBulk) Save(ctx context.Context) ([]*NovelSource, error) {
	specs := make([]*sqlgraph.CreateSpec, len(nscb.builders))
	nodes := make([]*NovelSource, len(nscb.builders))
	mutators := make([]Mutator, len(nscb.builders))
	for i := range nscb.builders {
		func(i int, root context.Context) {
			builder := nscb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*NovelSourceMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				nodes[i], specs[i] = builder.createSpec()
				var err error
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, nscb.builders[i+1].mutation)
				} else {
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, nscb.driver, &sqlgraph.BatchCreateSpec{Nodes: specs}); err != nil {
						if cerr, ok := isSQLConstraintError(err); ok {
							err = cerr
						}
					}
				}
				mutation.done = true
				if err != nil {
					return nil, err
				}
				id := specs[i].ID.Value.(int64)
				nodes[i].ID = int(id)
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, nscb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX calls Save and panics if Save returns an error.
func (nscb *NovelSourceCreateBulk) SaveX(ctx context.Context) []*NovelSource {
	v, err := nscb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}
