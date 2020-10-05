// Code generated by entc, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"math"

	"github.com/facebook/ent/dialect/sql"
	"github.com/facebook/ent/dialect/sql/sqlgraph"
	"github.com/facebook/ent/schema/field"
	"github.com/vicanso/elite/ent/novel"
	"github.com/vicanso/elite/ent/predicate"
)

// NovelQuery is the builder for querying Novel entities.
type NovelQuery struct {
	config
	limit      *int
	offset     *int
	order      []OrderFunc
	unique     []string
	predicates []predicate.Novel
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the builder.
func (nq *NovelQuery) Where(ps ...predicate.Novel) *NovelQuery {
	nq.predicates = append(nq.predicates, ps...)
	return nq
}

// Limit adds a limit step to the query.
func (nq *NovelQuery) Limit(limit int) *NovelQuery {
	nq.limit = &limit
	return nq
}

// Offset adds an offset step to the query.
func (nq *NovelQuery) Offset(offset int) *NovelQuery {
	nq.offset = &offset
	return nq
}

// Order adds an order step to the query.
func (nq *NovelQuery) Order(o ...OrderFunc) *NovelQuery {
	nq.order = append(nq.order, o...)
	return nq
}

// First returns the first Novel entity in the query. Returns *NotFoundError when no novel was found.
func (nq *NovelQuery) First(ctx context.Context) (*Novel, error) {
	nodes, err := nq.Limit(1).All(ctx)
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{novel.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (nq *NovelQuery) FirstX(ctx context.Context) *Novel {
	node, err := nq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first Novel id in the query. Returns *NotFoundError when no id was found.
func (nq *NovelQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = nq.Limit(1).IDs(ctx); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{novel.Label}
		return
	}
	return ids[0], nil
}

// FirstXID is like FirstID, but panics if an error occurs.
func (nq *NovelQuery) FirstXID(ctx context.Context) int {
	id, err := nq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns the only Novel entity in the query, returns an error if not exactly one entity was returned.
func (nq *NovelQuery) Only(ctx context.Context) (*Novel, error) {
	nodes, err := nq.Limit(2).All(ctx)
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{novel.Label}
	default:
		return nil, &NotSingularError{novel.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (nq *NovelQuery) OnlyX(ctx context.Context) *Novel {
	node, err := nq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID returns the only Novel id in the query, returns an error if not exactly one id was returned.
func (nq *NovelQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = nq.Limit(2).IDs(ctx); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{novel.Label}
	default:
		err = &NotSingularError{novel.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (nq *NovelQuery) OnlyIDX(ctx context.Context) int {
	id, err := nq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of Novels.
func (nq *NovelQuery) All(ctx context.Context) ([]*Novel, error) {
	if err := nq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	return nq.sqlAll(ctx)
}

// AllX is like All, but panics if an error occurs.
func (nq *NovelQuery) AllX(ctx context.Context) []*Novel {
	nodes, err := nq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of Novel ids.
func (nq *NovelQuery) IDs(ctx context.Context) ([]int, error) {
	var ids []int
	if err := nq.Select(novel.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (nq *NovelQuery) IDsX(ctx context.Context) []int {
	ids, err := nq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (nq *NovelQuery) Count(ctx context.Context) (int, error) {
	if err := nq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return nq.sqlCount(ctx)
}

// CountX is like Count, but panics if an error occurs.
func (nq *NovelQuery) CountX(ctx context.Context) int {
	count, err := nq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (nq *NovelQuery) Exist(ctx context.Context) (bool, error) {
	if err := nq.prepareQuery(ctx); err != nil {
		return false, err
	}
	return nq.sqlExist(ctx)
}

// ExistX is like Exist, but panics if an error occurs.
func (nq *NovelQuery) ExistX(ctx context.Context) bool {
	exist, err := nq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the query builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (nq *NovelQuery) Clone() *NovelQuery {
	return &NovelQuery{
		config:     nq.config,
		limit:      nq.limit,
		offset:     nq.offset,
		order:      append([]OrderFunc{}, nq.order...),
		unique:     append([]string{}, nq.unique...),
		predicates: append([]predicate.Novel{}, nq.predicates...),
		// clone intermediate query.
		sql:  nq.sql.Clone(),
		path: nq.path,
	}
}

// GroupBy used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		CreatedAt time.Time `json:"createdAt,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.Novel.Query().
//		GroupBy(novel.FieldCreatedAt).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
//
func (nq *NovelQuery) GroupBy(field string, fields ...string) *NovelGroupBy {
	group := &NovelGroupBy{config: nq.config}
	group.fields = append([]string{field}, fields...)
	group.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := nq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return nq.sqlQuery(), nil
	}
	return group
}

// Select one or more fields from the given query.
//
// Example:
//
//	var v []struct {
//		CreatedAt time.Time `json:"createdAt,omitempty"`
//	}
//
//	client.Novel.Query().
//		Select(novel.FieldCreatedAt).
//		Scan(ctx, &v)
//
func (nq *NovelQuery) Select(field string, fields ...string) *NovelSelect {
	selector := &NovelSelect{config: nq.config}
	selector.fields = append([]string{field}, fields...)
	selector.path = func(ctx context.Context) (prev *sql.Selector, err error) {
		if err := nq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		return nq.sqlQuery(), nil
	}
	return selector
}

func (nq *NovelQuery) prepareQuery(ctx context.Context) error {
	if nq.path != nil {
		prev, err := nq.path(ctx)
		if err != nil {
			return err
		}
		nq.sql = prev
	}
	return nil
}

func (nq *NovelQuery) sqlAll(ctx context.Context) ([]*Novel, error) {
	var (
		nodes = []*Novel{}
		_spec = nq.querySpec()
	)
	_spec.ScanValues = func() []interface{} {
		node := &Novel{config: nq.config}
		nodes = append(nodes, node)
		values := node.scanValues()
		return values
	}
	_spec.Assign = func(values ...interface{}) error {
		if len(nodes) == 0 {
			return fmt.Errorf("ent: Assign called without calling ScanValues")
		}
		node := nodes[len(nodes)-1]
		return node.assignValues(values...)
	}
	if err := sqlgraph.QueryNodes(ctx, nq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	return nodes, nil
}

func (nq *NovelQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := nq.querySpec()
	return sqlgraph.CountNodes(ctx, nq.driver, _spec)
}

func (nq *NovelQuery) sqlExist(ctx context.Context) (bool, error) {
	n, err := nq.sqlCount(ctx)
	if err != nil {
		return false, fmt.Errorf("ent: check existence: %v", err)
	}
	return n > 0, nil
}

func (nq *NovelQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := &sqlgraph.QuerySpec{
		Node: &sqlgraph.NodeSpec{
			Table:   novel.Table,
			Columns: novel.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeInt,
				Column: novel.FieldID,
			},
		},
		From:   nq.sql,
		Unique: true,
	}
	if ps := nq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := nq.limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := nq.offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := nq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector, novel.ValidColumn)
			}
		}
	}
	return _spec
}

func (nq *NovelQuery) sqlQuery() *sql.Selector {
	builder := sql.Dialect(nq.driver.Dialect())
	t1 := builder.Table(novel.Table)
	selector := builder.Select(t1.Columns(novel.Columns...)...).From(t1)
	if nq.sql != nil {
		selector = nq.sql
		selector.Select(selector.Columns(novel.Columns...)...)
	}
	for _, p := range nq.predicates {
		p(selector)
	}
	for _, p := range nq.order {
		p(selector, novel.ValidColumn)
	}
	if offset := nq.offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := nq.limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// NovelGroupBy is the builder for group-by Novel entities.
type NovelGroupBy struct {
	config
	fields []string
	fns    []AggregateFunc
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Aggregate adds the given aggregation functions to the group-by query.
func (ngb *NovelGroupBy) Aggregate(fns ...AggregateFunc) *NovelGroupBy {
	ngb.fns = append(ngb.fns, fns...)
	return ngb
}

// Scan applies the group-by query and scan the result into the given value.
func (ngb *NovelGroupBy) Scan(ctx context.Context, v interface{}) error {
	query, err := ngb.path(ctx)
	if err != nil {
		return err
	}
	ngb.sql = query
	return ngb.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (ngb *NovelGroupBy) ScanX(ctx context.Context, v interface{}) {
	if err := ngb.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from group-by. It is only allowed when querying group-by with one field.
func (ngb *NovelGroupBy) Strings(ctx context.Context) ([]string, error) {
	if len(ngb.fields) > 1 {
		return nil, errors.New("ent: NovelGroupBy.Strings is not achievable when grouping more than 1 field")
	}
	var v []string
	if err := ngb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (ngb *NovelGroupBy) StringsX(ctx context.Context) []string {
	v, err := ngb.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// String returns a single string from group-by. It is only allowed when querying group-by with one field.
func (ngb *NovelGroupBy) String(ctx context.Context) (_ string, err error) {
	var v []string
	if v, err = ngb.Strings(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{novel.Label}
	default:
		err = fmt.Errorf("ent: NovelGroupBy.Strings returned %d results when one was expected", len(v))
	}
	return
}

// StringX is like String, but panics if an error occurs.
func (ngb *NovelGroupBy) StringX(ctx context.Context) string {
	v, err := ngb.String(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from group-by. It is only allowed when querying group-by with one field.
func (ngb *NovelGroupBy) Ints(ctx context.Context) ([]int, error) {
	if len(ngb.fields) > 1 {
		return nil, errors.New("ent: NovelGroupBy.Ints is not achievable when grouping more than 1 field")
	}
	var v []int
	if err := ngb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (ngb *NovelGroupBy) IntsX(ctx context.Context) []int {
	v, err := ngb.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Int returns a single int from group-by. It is only allowed when querying group-by with one field.
func (ngb *NovelGroupBy) Int(ctx context.Context) (_ int, err error) {
	var v []int
	if v, err = ngb.Ints(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{novel.Label}
	default:
		err = fmt.Errorf("ent: NovelGroupBy.Ints returned %d results when one was expected", len(v))
	}
	return
}

// IntX is like Int, but panics if an error occurs.
func (ngb *NovelGroupBy) IntX(ctx context.Context) int {
	v, err := ngb.Int(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from group-by. It is only allowed when querying group-by with one field.
func (ngb *NovelGroupBy) Float64s(ctx context.Context) ([]float64, error) {
	if len(ngb.fields) > 1 {
		return nil, errors.New("ent: NovelGroupBy.Float64s is not achievable when grouping more than 1 field")
	}
	var v []float64
	if err := ngb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (ngb *NovelGroupBy) Float64sX(ctx context.Context) []float64 {
	v, err := ngb.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64 returns a single float64 from group-by. It is only allowed when querying group-by with one field.
func (ngb *NovelGroupBy) Float64(ctx context.Context) (_ float64, err error) {
	var v []float64
	if v, err = ngb.Float64s(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{novel.Label}
	default:
		err = fmt.Errorf("ent: NovelGroupBy.Float64s returned %d results when one was expected", len(v))
	}
	return
}

// Float64X is like Float64, but panics if an error occurs.
func (ngb *NovelGroupBy) Float64X(ctx context.Context) float64 {
	v, err := ngb.Float64(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from group-by. It is only allowed when querying group-by with one field.
func (ngb *NovelGroupBy) Bools(ctx context.Context) ([]bool, error) {
	if len(ngb.fields) > 1 {
		return nil, errors.New("ent: NovelGroupBy.Bools is not achievable when grouping more than 1 field")
	}
	var v []bool
	if err := ngb.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (ngb *NovelGroupBy) BoolsX(ctx context.Context) []bool {
	v, err := ngb.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bool returns a single bool from group-by. It is only allowed when querying group-by with one field.
func (ngb *NovelGroupBy) Bool(ctx context.Context) (_ bool, err error) {
	var v []bool
	if v, err = ngb.Bools(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{novel.Label}
	default:
		err = fmt.Errorf("ent: NovelGroupBy.Bools returned %d results when one was expected", len(v))
	}
	return
}

// BoolX is like Bool, but panics if an error occurs.
func (ngb *NovelGroupBy) BoolX(ctx context.Context) bool {
	v, err := ngb.Bool(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (ngb *NovelGroupBy) sqlScan(ctx context.Context, v interface{}) error {
	for _, f := range ngb.fields {
		if !novel.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("invalid field %q for group-by", f)}
		}
	}
	selector := ngb.sqlQuery()
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := ngb.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (ngb *NovelGroupBy) sqlQuery() *sql.Selector {
	selector := ngb.sql
	columns := make([]string, 0, len(ngb.fields)+len(ngb.fns))
	columns = append(columns, ngb.fields...)
	for _, fn := range ngb.fns {
		columns = append(columns, fn(selector, novel.ValidColumn))
	}
	return selector.Select(columns...).GroupBy(ngb.fields...)
}

// NovelSelect is the builder for select fields of Novel entities.
type NovelSelect struct {
	config
	fields []string
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Scan applies the selector query and scan the result into the given value.
func (ns *NovelSelect) Scan(ctx context.Context, v interface{}) error {
	query, err := ns.path(ctx)
	if err != nil {
		return err
	}
	ns.sql = query
	return ns.sqlScan(ctx, v)
}

// ScanX is like Scan, but panics if an error occurs.
func (ns *NovelSelect) ScanX(ctx context.Context, v interface{}) {
	if err := ns.Scan(ctx, v); err != nil {
		panic(err)
	}
}

// Strings returns list of strings from selector. It is only allowed when selecting one field.
func (ns *NovelSelect) Strings(ctx context.Context) ([]string, error) {
	if len(ns.fields) > 1 {
		return nil, errors.New("ent: NovelSelect.Strings is not achievable when selecting more than 1 field")
	}
	var v []string
	if err := ns.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// StringsX is like Strings, but panics if an error occurs.
func (ns *NovelSelect) StringsX(ctx context.Context) []string {
	v, err := ns.Strings(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// String returns a single string from selector. It is only allowed when selecting one field.
func (ns *NovelSelect) String(ctx context.Context) (_ string, err error) {
	var v []string
	if v, err = ns.Strings(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{novel.Label}
	default:
		err = fmt.Errorf("ent: NovelSelect.Strings returned %d results when one was expected", len(v))
	}
	return
}

// StringX is like String, but panics if an error occurs.
func (ns *NovelSelect) StringX(ctx context.Context) string {
	v, err := ns.String(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Ints returns list of ints from selector. It is only allowed when selecting one field.
func (ns *NovelSelect) Ints(ctx context.Context) ([]int, error) {
	if len(ns.fields) > 1 {
		return nil, errors.New("ent: NovelSelect.Ints is not achievable when selecting more than 1 field")
	}
	var v []int
	if err := ns.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// IntsX is like Ints, but panics if an error occurs.
func (ns *NovelSelect) IntsX(ctx context.Context) []int {
	v, err := ns.Ints(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Int returns a single int from selector. It is only allowed when selecting one field.
func (ns *NovelSelect) Int(ctx context.Context) (_ int, err error) {
	var v []int
	if v, err = ns.Ints(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{novel.Label}
	default:
		err = fmt.Errorf("ent: NovelSelect.Ints returned %d results when one was expected", len(v))
	}
	return
}

// IntX is like Int, but panics if an error occurs.
func (ns *NovelSelect) IntX(ctx context.Context) int {
	v, err := ns.Int(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64s returns list of float64s from selector. It is only allowed when selecting one field.
func (ns *NovelSelect) Float64s(ctx context.Context) ([]float64, error) {
	if len(ns.fields) > 1 {
		return nil, errors.New("ent: NovelSelect.Float64s is not achievable when selecting more than 1 field")
	}
	var v []float64
	if err := ns.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// Float64sX is like Float64s, but panics if an error occurs.
func (ns *NovelSelect) Float64sX(ctx context.Context) []float64 {
	v, err := ns.Float64s(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Float64 returns a single float64 from selector. It is only allowed when selecting one field.
func (ns *NovelSelect) Float64(ctx context.Context) (_ float64, err error) {
	var v []float64
	if v, err = ns.Float64s(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{novel.Label}
	default:
		err = fmt.Errorf("ent: NovelSelect.Float64s returned %d results when one was expected", len(v))
	}
	return
}

// Float64X is like Float64, but panics if an error occurs.
func (ns *NovelSelect) Float64X(ctx context.Context) float64 {
	v, err := ns.Float64(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bools returns list of bools from selector. It is only allowed when selecting one field.
func (ns *NovelSelect) Bools(ctx context.Context) ([]bool, error) {
	if len(ns.fields) > 1 {
		return nil, errors.New("ent: NovelSelect.Bools is not achievable when selecting more than 1 field")
	}
	var v []bool
	if err := ns.Scan(ctx, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// BoolsX is like Bools, but panics if an error occurs.
func (ns *NovelSelect) BoolsX(ctx context.Context) []bool {
	v, err := ns.Bools(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Bool returns a single bool from selector. It is only allowed when selecting one field.
func (ns *NovelSelect) Bool(ctx context.Context) (_ bool, err error) {
	var v []bool
	if v, err = ns.Bools(ctx); err != nil {
		return
	}
	switch len(v) {
	case 1:
		return v[0], nil
	case 0:
		err = &NotFoundError{novel.Label}
	default:
		err = fmt.Errorf("ent: NovelSelect.Bools returned %d results when one was expected", len(v))
	}
	return
}

// BoolX is like Bool, but panics if an error occurs.
func (ns *NovelSelect) BoolX(ctx context.Context) bool {
	v, err := ns.Bool(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

func (ns *NovelSelect) sqlScan(ctx context.Context, v interface{}) error {
	for _, f := range ns.fields {
		if !novel.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("invalid field %q for selection", f)}
		}
	}
	rows := &sql.Rows{}
	query, args := ns.sqlQuery().Query()
	if err := ns.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

func (ns *NovelSelect) sqlQuery() sql.Querier {
	selector := ns.sql
	selector.Select(selector.Columns(ns.fields...)...)
	return selector
}