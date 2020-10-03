// Code generated by entc, DO NOT EDIT.

package ent

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/facebook/ent/dialect/sql"
	"github.com/vicanso/elite/ent/schema"
	"github.com/vicanso/elite/ent/user"
)

// User is the model entity for the User schema.
type User struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"createdAt,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
	// Status holds the value of the "status" field.
	Status schema.Status `json:"status,omitempty"`
	// Account holds the value of the "account" field.
	Account string `json:"account,omitempty"`
	// Password holds the value of the "password" field.
	Password string `json:"-"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Roles holds the value of the "roles" field.
	Roles []string `json:"roles,omitempty"`
	// Groups holds the value of the "groups" field.
	Groups []string `json:"groups,omitempty"`
	// Email holds the value of the "email" field.
	Email string `json:"email,omitempty"`
}

// scanValues returns the types for scanning values from sql.Rows.
func (*User) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},  // id
		&sql.NullTime{},   // created_at
		&sql.NullTime{},   // updated_at
		&sql.NullInt64{},  // status
		&sql.NullString{}, // account
		&sql.NullString{}, // password
		&sql.NullString{}, // name
		&[]byte{},         // roles
		&[]byte{},         // groups
		&sql.NullString{}, // email
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the User fields.
func (u *User) assignValues(values ...interface{}) error {
	if m, n := len(values), len(user.Columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	u.ID = int(value.Int64)
	values = values[1:]
	if value, ok := values[0].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field created_at", values[0])
	} else if value.Valid {
		u.CreatedAt = value.Time
	}
	if value, ok := values[1].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field updated_at", values[1])
	} else if value.Valid {
		u.UpdatedAt = value.Time
	}
	if value, ok := values[2].(*sql.NullInt64); !ok {
		return fmt.Errorf("unexpected type %T for field status", values[2])
	} else if value.Valid {
		u.Status = schema.Status(value.Int64)
	}
	if value, ok := values[3].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field account", values[3])
	} else if value.Valid {
		u.Account = value.String
	}
	if value, ok := values[4].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field password", values[4])
	} else if value.Valid {
		u.Password = value.String
	}
	if value, ok := values[5].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field name", values[5])
	} else if value.Valid {
		u.Name = value.String
	}

	if value, ok := values[6].(*[]byte); !ok {
		return fmt.Errorf("unexpected type %T for field roles", values[6])
	} else if value != nil && len(*value) > 0 {
		if err := json.Unmarshal(*value, &u.Roles); err != nil {
			return fmt.Errorf("unmarshal field roles: %v", err)
		}
	}

	if value, ok := values[7].(*[]byte); !ok {
		return fmt.Errorf("unexpected type %T for field groups", values[7])
	} else if value != nil && len(*value) > 0 {
		if err := json.Unmarshal(*value, &u.Groups); err != nil {
			return fmt.Errorf("unmarshal field groups: %v", err)
		}
	}
	if value, ok := values[8].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field email", values[8])
	} else if value.Valid {
		u.Email = value.String
	}
	return nil
}

// Update returns a builder for updating this User.
// Note that, you need to call User.Unwrap() before calling this method, if this User
// was returned from a transaction, and the transaction was committed or rolled back.
func (u *User) Update() *UserUpdateOne {
	return (&UserClient{config: u.config}).UpdateOne(u)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (u *User) Unwrap() *User {
	tx, ok := u.config.driver.(*txDriver)
	if !ok {
		panic("ent: User is not a transactional entity")
	}
	u.config.driver = tx.drv
	return u
}

// String implements the fmt.Stringer.
func (u *User) String() string {
	var builder strings.Builder
	builder.WriteString("User(")
	builder.WriteString(fmt.Sprintf("id=%v", u.ID))
	builder.WriteString(", created_at=")
	builder.WriteString(u.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", updated_at=")
	builder.WriteString(u.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", status=")
	builder.WriteString(fmt.Sprintf("%v", u.Status))
	builder.WriteString(", account=")
	builder.WriteString(u.Account)
	builder.WriteString(", password=<sensitive>")
	builder.WriteString(", name=")
	builder.WriteString(u.Name)
	builder.WriteString(", roles=")
	builder.WriteString(fmt.Sprintf("%v", u.Roles))
	builder.WriteString(", groups=")
	builder.WriteString(fmt.Sprintf("%v", u.Groups))
	builder.WriteString(", email=")
	builder.WriteString(u.Email)
	builder.WriteByte(')')
	return builder.String()
}

// Users is a parsable slice of User.
type Users []*User

func (u Users) config(cfg config) {
	for _i := range u {
		u[_i].config = cfg
	}
}
