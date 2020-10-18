// Code generated by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"github.com/facebook/ent/dialect/sql"
	"github.com/vicanso/elite/ent/novel"
)

// Novel is the model entity for the Novel schema.
type Novel struct {
	config `json:"-" sql:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"createdAt,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Author holds the value of the "author" field.
	Author string `json:"author,omitempty"`
	// Source holds the value of the "source" field.
	Source int `json:"source,omitempty"`
	// Status holds the value of the "status" field.
	Status int `json:"status,omitempty"`
	// WordCount holds the value of the "word_count" field.
	WordCount int `json:"wordCount,omitempty" sql:"word_count"`
	// Views holds the value of the "views" field.
	Views int `json:"views,omitempty"`
	// Downloads holds the value of the "downloads" field.
	Downloads int `json:"downloads,omitempty"`
	// Favorites holds the value of the "favorites" field.
	Favorites int `json:"favorites,omitempty"`
	// Cover holds the value of the "cover" field.
	Cover string `json:"cover,omitempty"`
	// Summary holds the value of the "summary" field.
	Summary string `json:"summary,omitempty"`
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Novel) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},  // id
		&sql.NullTime{},   // created_at
		&sql.NullTime{},   // updated_at
		&sql.NullString{}, // name
		&sql.NullString{}, // author
		&sql.NullInt64{},  // source
		&sql.NullInt64{},  // status
		&sql.NullInt64{},  // word_count
		&sql.NullInt64{},  // views
		&sql.NullInt64{},  // downloads
		&sql.NullInt64{},  // favorites
		&sql.NullString{}, // cover
		&sql.NullString{}, // summary
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Novel fields.
func (n *Novel) assignValues(values ...interface{}) error {
	if m, n := len(values), len(novel.Columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	n.ID = int(value.Int64)
	values = values[1:]
	if value, ok := values[0].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field created_at", values[0])
	} else if value.Valid {
		n.CreatedAt = value.Time
	}
	if value, ok := values[1].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field updated_at", values[1])
	} else if value.Valid {
		n.UpdatedAt = value.Time
	}
	if value, ok := values[2].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field name", values[2])
	} else if value.Valid {
		n.Name = value.String
	}
	if value, ok := values[3].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field author", values[3])
	} else if value.Valid {
		n.Author = value.String
	}
	if value, ok := values[4].(*sql.NullInt64); !ok {
		return fmt.Errorf("unexpected type %T for field source", values[4])
	} else if value.Valid {
		n.Source = int(value.Int64)
	}
	if value, ok := values[5].(*sql.NullInt64); !ok {
		return fmt.Errorf("unexpected type %T for field status", values[5])
	} else if value.Valid {
		n.Status = int(value.Int64)
	}
	if value, ok := values[6].(*sql.NullInt64); !ok {
		return fmt.Errorf("unexpected type %T for field word_count", values[6])
	} else if value.Valid {
		n.WordCount = int(value.Int64)
	}
	if value, ok := values[7].(*sql.NullInt64); !ok {
		return fmt.Errorf("unexpected type %T for field views", values[7])
	} else if value.Valid {
		n.Views = int(value.Int64)
	}
	if value, ok := values[8].(*sql.NullInt64); !ok {
		return fmt.Errorf("unexpected type %T for field downloads", values[8])
	} else if value.Valid {
		n.Downloads = int(value.Int64)
	}
	if value, ok := values[9].(*sql.NullInt64); !ok {
		return fmt.Errorf("unexpected type %T for field favorites", values[9])
	} else if value.Valid {
		n.Favorites = int(value.Int64)
	}
	if value, ok := values[10].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field cover", values[10])
	} else if value.Valid {
		n.Cover = value.String
	}
	if value, ok := values[11].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field summary", values[11])
	} else if value.Valid {
		n.Summary = value.String
	}
	return nil
}

// Update returns a builder for updating this Novel.
// Note that, you need to call Novel.Unwrap() before calling this method, if this Novel
// was returned from a transaction, and the transaction was committed or rolled back.
func (n *Novel) Update() *NovelUpdateOne {
	return (&NovelClient{config: n.config}).UpdateOne(n)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (n *Novel) Unwrap() *Novel {
	tx, ok := n.config.driver.(*txDriver)
	if !ok {
		panic("ent: Novel is not a transactional entity")
	}
	n.config.driver = tx.drv
	return n
}

// String implements the fmt.Stringer.
func (n *Novel) String() string {
	var builder strings.Builder
	builder.WriteString("Novel(")
	builder.WriteString(fmt.Sprintf("id=%v", n.ID))
	builder.WriteString(", created_at=")
	builder.WriteString(n.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", updated_at=")
	builder.WriteString(n.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", name=")
	builder.WriteString(n.Name)
	builder.WriteString(", author=")
	builder.WriteString(n.Author)
	builder.WriteString(", source=")
	builder.WriteString(fmt.Sprintf("%v", n.Source))
	builder.WriteString(", status=")
	builder.WriteString(fmt.Sprintf("%v", n.Status))
	builder.WriteString(", word_count=")
	builder.WriteString(fmt.Sprintf("%v", n.WordCount))
	builder.WriteString(", views=")
	builder.WriteString(fmt.Sprintf("%v", n.Views))
	builder.WriteString(", downloads=")
	builder.WriteString(fmt.Sprintf("%v", n.Downloads))
	builder.WriteString(", favorites=")
	builder.WriteString(fmt.Sprintf("%v", n.Favorites))
	builder.WriteString(", cover=")
	builder.WriteString(n.Cover)
	builder.WriteString(", summary=")
	builder.WriteString(n.Summary)
	builder.WriteByte(')')
	return builder.String()
}

// Novels is a parsable slice of Novel.
type Novels []*Novel

func (n Novels) config(cfg config) {
	for _i := range n {
		n[_i].config = cfg
	}
}
