package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// OrganizationMember holds the schema definition for the OrganizationMember entity.
type OrganizationMember struct {
	ent.Schema
}

// Fields of the OrganizationMember.
func (OrganizationMember) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Unique().
			Immutable(),
		field.UUID("organization_id", uuid.UUID{}).
			Comment("Organization ID"),
		field.UUID("user_id", uuid.UUID{}).
			Comment("User ID"),
		field.Enum("role").
			Values("admin", "member").
			Default("member").
			Comment("Member role in the organization"),
		field.Time("joined_at").
			Default(time.Now).
			Immutable().
			Comment("When the user joined the organization"),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the OrganizationMember.
func (OrganizationMember) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("organization", Organization.Type).
			Ref("members").
			Field("organization_id").
			Required().
			Unique(),
		edge.From("user", User.Type).
			Ref("memberships").
			Field("user_id").
			Required().
			Unique(),
	}
}
