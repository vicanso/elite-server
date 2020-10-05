// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/facebook/ent/dialect/sql/sqlgraph"
	"github.com/facebook/ent/schema/field"
	"github.com/vicanso/elite/ent/novel"
)

// NovelCreate is the builder for creating a Novel entity.
type NovelCreate struct {
	config
	mutation *NovelMutation
	hooks    []Hook
}

// SetCreatedAt sets the created_at field.
func (nc *NovelCreate) SetCreatedAt(t time.Time) *NovelCreate {
	nc.mutation.SetCreatedAt(t)
	return nc
}

// SetNillableCreatedAt sets the created_at field if the given value is not nil.
func (nc *NovelCreate) SetNillableCreatedAt(t *time.Time) *NovelCreate {
	if t != nil {
		nc.SetCreatedAt(*t)
	}
	return nc
}

// SetUpdatedAt sets the updated_at field.
func (nc *NovelCreate) SetUpdatedAt(t time.Time) *NovelCreate {
	nc.mutation.SetUpdatedAt(t)
	return nc
}

// SetNillableUpdatedAt sets the updated_at field if the given value is not nil.
func (nc *NovelCreate) SetNillableUpdatedAt(t *time.Time) *NovelCreate {
	if t != nil {
		nc.SetUpdatedAt(*t)
	}
	return nc
}

// SetName sets the name field.
func (nc *NovelCreate) SetName(s string) *NovelCreate {
	nc.mutation.SetName(s)
	return nc
}

// SetAuthor sets the author field.
func (nc *NovelCreate) SetAuthor(s string) *NovelCreate {
	nc.mutation.SetAuthor(s)
	return nc
}

// SetSource sets the source field.
func (nc *NovelCreate) SetSource(i int) *NovelCreate {
	nc.mutation.SetSource(i)
	return nc
}

// SetStatus sets the status field.
func (nc *NovelCreate) SetStatus(i int) *NovelCreate {
	nc.mutation.SetStatus(i)
	return nc
}

// SetNillableStatus sets the status field if the given value is not nil.
func (nc *NovelCreate) SetNillableStatus(i *int) *NovelCreate {
	if i != nil {
		nc.SetStatus(*i)
	}
	return nc
}

// SetCover sets the cover field.
func (nc *NovelCreate) SetCover(s string) *NovelCreate {
	nc.mutation.SetCover(s)
	return nc
}

// SetNillableCover sets the cover field if the given value is not nil.
func (nc *NovelCreate) SetNillableCover(s *string) *NovelCreate {
	if s != nil {
		nc.SetCover(*s)
	}
	return nc
}

// SetSummary sets the summary field.
func (nc *NovelCreate) SetSummary(s string) *NovelCreate {
	nc.mutation.SetSummary(s)
	return nc
}

// Mutation returns the NovelMutation object of the builder.
func (nc *NovelCreate) Mutation() *NovelMutation {
	return nc.mutation
}

// Save creates the Novel in the database.
func (nc *NovelCreate) Save(ctx context.Context) (*Novel, error) {
	var (
		err  error
		node *Novel
	)
	nc.defaults()
	if len(nc.hooks) == 0 {
		if err = nc.check(); err != nil {
			return nil, err
		}
		node, err = nc.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*NovelMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = nc.check(); err != nil {
				return nil, err
			}
			nc.mutation = mutation
			node, err = nc.sqlSave(ctx)
			mutation.done = true
			return node, err
		})
		for i := len(nc.hooks) - 1; i >= 0; i-- {
			mut = nc.hooks[i](mut)
		}
		if _, err := mut.Mutate(ctx, nc.mutation); err != nil {
			return nil, err
		}
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (nc *NovelCreate) SaveX(ctx context.Context) *Novel {
	v, err := nc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// defaults sets the default values of the builder before save.
func (nc *NovelCreate) defaults() {
	if _, ok := nc.mutation.CreatedAt(); !ok {
		v := novel.DefaultCreatedAt()
		nc.mutation.SetCreatedAt(v)
	}
	if _, ok := nc.mutation.UpdatedAt(); !ok {
		v := novel.DefaultUpdatedAt()
		nc.mutation.SetUpdatedAt(v)
	}
	if _, ok := nc.mutation.Status(); !ok {
		v := novel.DefaultStatus
		nc.mutation.SetStatus(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (nc *NovelCreate) check() error {
	if _, ok := nc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New("ent: missing required field \"created_at\"")}
	}
	if _, ok := nc.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New("ent: missing required field \"updated_at\"")}
	}
	if _, ok := nc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New("ent: missing required field \"name\"")}
	}
	if v, ok := nc.mutation.Name(); ok {
		if err := novel.NameValidator(v); err != nil {
			return &ValidationError{Name: "name", err: fmt.Errorf("ent: validator failed for field \"name\": %w", err)}
		}
	}
	if _, ok := nc.mutation.Author(); !ok {
		return &ValidationError{Name: "author", err: errors.New("ent: missing required field \"author\"")}
	}
	if v, ok := nc.mutation.Author(); ok {
		if err := novel.AuthorValidator(v); err != nil {
			return &ValidationError{Name: "author", err: fmt.Errorf("ent: validator failed for field \"author\": %w", err)}
		}
	}
	if _, ok := nc.mutation.Source(); !ok {
		return &ValidationError{Name: "source", err: errors.New("ent: missing required field \"source\"")}
	}
	if v, ok := nc.mutation.Source(); ok {
		if err := novel.SourceValidator(v); err != nil {
			return &ValidationError{Name: "source", err: fmt.Errorf("ent: validator failed for field \"source\": %w", err)}
		}
	}
	if _, ok := nc.mutation.Status(); !ok {
		return &ValidationError{Name: "status", err: errors.New("ent: missing required field \"status\"")}
	}
	if v, ok := nc.mutation.Status(); ok {
		if err := novel.StatusValidator(v); err != nil {
			return &ValidationError{Name: "status", err: fmt.Errorf("ent: validator failed for field \"status\": %w", err)}
		}
	}
	if _, ok := nc.mutation.Summary(); !ok {
		return &ValidationError{Name: "summary", err: errors.New("ent: missing required field \"summary\"")}
	}
	return nil
}

func (nc *NovelCreate) sqlSave(ctx context.Context) (*Novel, error) {
	_node, _spec := nc.createSpec()
	if err := sqlgraph.CreateNode(ctx, nc.driver, _spec); err != nil {
		if cerr, ok := isSQLConstraintError(err); ok {
			err = cerr
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	return _node, nil
}

func (nc *NovelCreate) createSpec() (*Novel, *sqlgraph.CreateSpec) {
	var (
		_node = &Novel{config: nc.config}
		_spec = &sqlgraph.CreateSpec{
			Table: novel.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: novel.FieldID,
			},
		}
	)
	if value, ok := nc.mutation.CreatedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: novel.FieldCreatedAt,
		})
		_node.CreatedAt = value
	}
	if value, ok := nc.mutation.UpdatedAt(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeTime,
			Value:  value,
			Column: novel.FieldUpdatedAt,
		})
		_node.UpdatedAt = value
	}
	if value, ok := nc.mutation.Name(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: novel.FieldName,
		})
		_node.Name = value
	}
	if value, ok := nc.mutation.Author(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: novel.FieldAuthor,
		})
		_node.Author = value
	}
	if value, ok := nc.mutation.Source(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: novel.FieldSource,
		})
		_node.Source = value
	}
	if value, ok := nc.mutation.Status(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeInt,
			Value:  value,
			Column: novel.FieldStatus,
		})
		_node.Status = value
	}
	if value, ok := nc.mutation.Cover(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: novel.FieldCover,
		})
		_node.Cover = value
	}
	if value, ok := nc.mutation.Summary(); ok {
		_spec.Fields = append(_spec.Fields, &sqlgraph.FieldSpec{
			Type:   field.TypeString,
			Value:  value,
			Column: novel.FieldSummary,
		})
		_node.Summary = value
	}
	return _node, _spec
}

// NovelCreateBulk is the builder for creating a bulk of Novel entities.
type NovelCreateBulk struct {
	config
	builders []*NovelCreate
}

// Save creates the Novel entities in the database.
func (ncb *NovelCreateBulk) Save(ctx context.Context) ([]*Novel, error) {
	specs := make([]*sqlgraph.CreateSpec, len(ncb.builders))
	nodes := make([]*Novel, len(ncb.builders))
	mutators := make([]Mutator, len(ncb.builders))
	for i := range ncb.builders {
		func(i int, root context.Context) {
			builder := ncb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*NovelMutation)
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
					_, err = mutators[i+1].Mutate(root, ncb.builders[i+1].mutation)
				} else {
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, ncb.driver, &sqlgraph.BatchCreateSpec{Nodes: specs}); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, ncb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX calls Save and panics if Save returns an error.
func (ncb *NovelCreateBulk) SaveX(ctx context.Context) []*Novel {
	v, err := ncb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}