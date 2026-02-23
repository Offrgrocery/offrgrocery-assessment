package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// List holds the schema definition for the List entity.
type List struct {
	ent.Schema
}

// Fields of the List.
func (List) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Default("my new list").
			NotEmpty(),
	}
}

// Edges of the List.
func (List) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("lists").
			Unique().
			Required(),
		edge.To("items", Item.Type),
	}
}

func (List) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}
