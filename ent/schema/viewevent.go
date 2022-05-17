package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// ViewEvent holds the schema definition for the ViewEvent entity.
type ViewEvent struct {
	ent.Schema
}

// Fields of the ViewEvent.
func (ViewEvent) Fields() []ent.Field {
	return []ent.Field{
		field.String("page"),
		field.String("ip_address"),
		field.Time("event_time"),
	}
}

// Edges of the ViewEvent.
func (ViewEvent) Edges() []ent.Edge {
	return nil
}
