// Code generated by entc, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// ViewEventsColumns holds the columns for the "view_events" table.
	ViewEventsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "page", Type: field.TypeString},
		{Name: "ip_address", Type: field.TypeString},
		{Name: "event_time", Type: field.TypeTime},
	}
	// ViewEventsTable holds the schema information for the "view_events" table.
	ViewEventsTable = &schema.Table{
		Name:       "view_events",
		Columns:    ViewEventsColumns,
		PrimaryKey: []*schema.Column{ViewEventsColumns[0]},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		ViewEventsTable,
	}
)

func init() {
}
