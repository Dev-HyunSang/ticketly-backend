package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Annotations(),
		field.String("first_name"),
		field.String("last_name"),
		field.String("nick_name"),
		field.String("birthday"),
		field.String("email"),
		field.String("password"),
		field.String("phone_number"),
		field.Bool("is_valid").
			Default(true),
		field.Time("created_at").
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("owned_organizations", Organization.Type),
		edge.To("memberships", OrganizationMember.Type),
		edge.To("created_events", Event.Type),
		edge.To("payments", Payment.Type),
	}
}
