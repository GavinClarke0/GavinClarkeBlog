// Code generated by entc, DO NOT EDIT.

package ent

import (
	"Blog/ent/viewevent"
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
)

// ViewEvent is the model entity for the ViewEvent schema.
type ViewEvent struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Page holds the value of the "page" field.
	Page string `json:"page,omitempty"`
	// IPAddress holds the value of the "ip_address" field.
	IPAddress string `json:"ip_address,omitempty"`
	// EventTime holds the value of the "event_time" field.
	EventTime time.Time `json:"event_time,omitempty"`
}

// scanValues returns the types for scanning values from sql.Rows.
func (*ViewEvent) scanValues(columns []string) ([]interface{}, error) {
	values := make([]interface{}, len(columns))
	for i := range columns {
		switch columns[i] {
		case viewevent.FieldID:
			values[i] = new(sql.NullInt64)
		case viewevent.FieldPage, viewevent.FieldIPAddress:
			values[i] = new(sql.NullString)
		case viewevent.FieldEventTime:
			values[i] = new(sql.NullTime)
		default:
			return nil, fmt.Errorf("unexpected column %q for type ViewEvent", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the ViewEvent fields.
func (ve *ViewEvent) assignValues(columns []string, values []interface{}) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case viewevent.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			ve.ID = int(value.Int64)
		case viewevent.FieldPage:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field page", values[i])
			} else if value.Valid {
				ve.Page = value.String
			}
		case viewevent.FieldIPAddress:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field ip_address", values[i])
			} else if value.Valid {
				ve.IPAddress = value.String
			}
		case viewevent.FieldEventTime:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field event_time", values[i])
			} else if value.Valid {
				ve.EventTime = value.Time
			}
		}
	}
	return nil
}

// Update returns a builder for updating this ViewEvent.
// Note that you need to call ViewEvent.Unwrap() before calling this method if this ViewEvent
// was returned from a transaction, and the transaction was committed or rolled back.
func (ve *ViewEvent) Update() *ViewEventUpdateOne {
	return (&ViewEventClient{config: ve.config}).UpdateOne(ve)
}

// Unwrap unwraps the ViewEvent entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (ve *ViewEvent) Unwrap() *ViewEvent {
	tx, ok := ve.config.driver.(*txDriver)
	if !ok {
		panic("ent: ViewEvent is not a transactional entity")
	}
	ve.config.driver = tx.drv
	return ve
}

// String implements the fmt.Stringer.
func (ve *ViewEvent) String() string {
	var builder strings.Builder
	builder.WriteString("ViewEvent(")
	builder.WriteString(fmt.Sprintf("id=%v", ve.ID))
	builder.WriteString(", page=")
	builder.WriteString(ve.Page)
	builder.WriteString(", ip_address=")
	builder.WriteString(ve.IPAddress)
	builder.WriteString(", event_time=")
	builder.WriteString(ve.EventTime.Format(time.ANSIC))
	builder.WriteByte(')')
	return builder.String()
}

// ViewEvents is a parsable slice of ViewEvent.
type ViewEvents []*ViewEvent

func (ve ViewEvents) config(cfg config) {
	for _i := range ve {
		ve[_i].config = cfg
	}
}
