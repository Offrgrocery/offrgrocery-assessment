package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"entgo.io/ent/schema/mixin"
)

// Item holds the schema definition for the Item entity.
type Item struct {
	ent.Schema
}

// Fields of the Item.
func (Item) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty(),
		field.String("brand").
			NotEmpty(),
		field.Float("price").
			Positive(),
	}
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

func (Item) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name", "brand").
			Annotations(
				entsql.IndexType("FULLTEXT"),
			),
	}
}

func (Item) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}
