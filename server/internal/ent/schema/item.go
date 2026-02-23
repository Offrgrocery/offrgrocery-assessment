package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
)

// Item holds the schema definition for the Item entity.
type Item struct {
	ent.Schema
}

// Fields of the Item.
func (Item) Fields() []ent.Field {
	return nil
}

// Edges of the Item.
func (Item) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("store", Store.Type).
			Ref("items").
			Unique().
			Required(),
		edge.From("lists", List.Type).
			Ref("items"),
	}
}
