package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Organization holds the schema definition for the Organization entity.
type Organization struct {
	ent.Schema
}

// Fields of the Organization.
func (Organization) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Unique().
			Immutable(),
		field.String("name").
			NotEmpty().
			Comment("Organization name"),
		field.Text("description").
			Optional().
			Comment("Organization description"),
		field.String("logo_url").
			Optional().
			Comment("Organization logo URL"),
		field.UUID("owner_id", uuid.UUID{}).
			Comment("User ID of the organization owner"),
		field.Bool("is_active").
			Default(true).
			Comment("Whether the organization is active"),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Organization.
func (Organization) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("members", OrganizationMember.Type),
		edge.To("events", Event.Type),
		edge.From("owner", User.Type).
			Ref("owned_organizations").
			Field("owner_id").
			Required().
			Unique(),
	}
}
