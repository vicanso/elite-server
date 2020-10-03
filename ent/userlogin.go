// Code generated by entc, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"github.com/facebook/ent/dialect/sql"
	"github.com/vicanso/elite/ent/userlogin"
)

// UserLogin is the model entity for the UserLogin schema.
type UserLogin struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"createdAt,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
	// Account holds the value of the "account" field.
	Account string `json:"account,omitempty"`
	// UserAgent holds the value of the "user_agent" field.
	UserAgent string `json:"userAgent,omitempty"`
	// IP holds the value of the "ip" field.
	IP string `json:"ip,omitempty"`
	// TrackID holds the value of the "track_id" field.
	TrackID string `json:"trackID,omitempty"`
	// SessionID holds the value of the "session_id" field.
	SessionID string `json:"sessionID,omitempty"`
	// XForwardedFor holds the value of the "x_forwarded_for" field.
	XForwardedFor string `json:"xForwardedFor,omitempty"`
	// Country holds the value of the "country" field.
	Country string `json:"country,omitempty"`
	// Province holds the value of the "province" field.
	Province string `json:"province,omitempty"`
	// City holds the value of the "city" field.
	City string `json:"city,omitempty"`
	// Isp holds the value of the "isp" field.
	Isp string `json:"isp,omitempty"`
}

// scanValues returns the types for scanning values from sql.Rows.
func (*UserLogin) scanValues() []interface{} {
	return []interface{}{
		&sql.NullInt64{},  // id
		&sql.NullTime{},   // created_at
		&sql.NullTime{},   // updated_at
		&sql.NullString{}, // account
		&sql.NullString{}, // user_agent
		&sql.NullString{}, // ip
		&sql.NullString{}, // track_id
		&sql.NullString{}, // session_id
		&sql.NullString{}, // x_forwarded_for
		&sql.NullString{}, // country
		&sql.NullString{}, // province
		&sql.NullString{}, // city
		&sql.NullString{}, // isp
	}
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the UserLogin fields.
func (ul *UserLogin) assignValues(values ...interface{}) error {
	if m, n := len(values), len(userlogin.Columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	value, ok := values[0].(*sql.NullInt64)
	if !ok {
		return fmt.Errorf("unexpected type %T for field id", value)
	}
	ul.ID = int(value.Int64)
	values = values[1:]
	if value, ok := values[0].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field created_at", values[0])
	} else if value.Valid {
		ul.CreatedAt = value.Time
	}
	if value, ok := values[1].(*sql.NullTime); !ok {
		return fmt.Errorf("unexpected type %T for field updated_at", values[1])
	} else if value.Valid {
		ul.UpdatedAt = value.Time
	}
	if value, ok := values[2].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field account", values[2])
	} else if value.Valid {
		ul.Account = value.String
	}
	if value, ok := values[3].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field user_agent", values[3])
	} else if value.Valid {
		ul.UserAgent = value.String
	}
	if value, ok := values[4].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field ip", values[4])
	} else if value.Valid {
		ul.IP = value.String
	}
	if value, ok := values[5].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field track_id", values[5])
	} else if value.Valid {
		ul.TrackID = value.String
	}
	if value, ok := values[6].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field session_id", values[6])
	} else if value.Valid {
		ul.SessionID = value.String
	}
	if value, ok := values[7].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field x_forwarded_for", values[7])
	} else if value.Valid {
		ul.XForwardedFor = value.String
	}
	if value, ok := values[8].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field country", values[8])
	} else if value.Valid {
		ul.Country = value.String
	}
	if value, ok := values[9].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field province", values[9])
	} else if value.Valid {
		ul.Province = value.String
	}
	if value, ok := values[10].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field city", values[10])
	} else if value.Valid {
		ul.City = value.String
	}
	if value, ok := values[11].(*sql.NullString); !ok {
		return fmt.Errorf("unexpected type %T for field isp", values[11])
	} else if value.Valid {
		ul.Isp = value.String
	}
	return nil
}

// Update returns a builder for updating this UserLogin.
// Note that, you need to call UserLogin.Unwrap() before calling this method, if this UserLogin
// was returned from a transaction, and the transaction was committed or rolled back.
func (ul *UserLogin) Update() *UserLoginUpdateOne {
	return (&UserLoginClient{config: ul.config}).UpdateOne(ul)
}

// Unwrap unwraps the entity that was returned from a transaction after it was closed,
// so that all next queries will be executed through the driver which created the transaction.
func (ul *UserLogin) Unwrap() *UserLogin {
	tx, ok := ul.config.driver.(*txDriver)
	if !ok {
		panic("ent: UserLogin is not a transactional entity")
	}
	ul.config.driver = tx.drv
	return ul
}

// String implements the fmt.Stringer.
func (ul *UserLogin) String() string {
	var builder strings.Builder
	builder.WriteString("UserLogin(")
	builder.WriteString(fmt.Sprintf("id=%v", ul.ID))
	builder.WriteString(", created_at=")
	builder.WriteString(ul.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", updated_at=")
	builder.WriteString(ul.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", account=")
	builder.WriteString(ul.Account)
	builder.WriteString(", user_agent=")
	builder.WriteString(ul.UserAgent)
	builder.WriteString(", ip=")
	builder.WriteString(ul.IP)
	builder.WriteString(", track_id=")
	builder.WriteString(ul.TrackID)
	builder.WriteString(", session_id=")
	builder.WriteString(ul.SessionID)
	builder.WriteString(", x_forwarded_for=")
	builder.WriteString(ul.XForwardedFor)
	builder.WriteString(", country=")
	builder.WriteString(ul.Country)
	builder.WriteString(", province=")
	builder.WriteString(ul.Province)
	builder.WriteString(", city=")
	builder.WriteString(ul.City)
	builder.WriteString(", isp=")
	builder.WriteString(ul.Isp)
	builder.WriteByte(')')
	return builder.String()
}

// UserLogins is a parsable slice of UserLogin.
type UserLogins []*UserLogin

func (ul UserLogins) config(cfg config) {
	for _i := range ul {
		ul[_i].config = cfg
	}
}
