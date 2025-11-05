package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Event holds the schema definition for the Event entity.
type Event struct {
	ent.Schema
}

// Fields of the Event.
func (Event) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Unique().
			Immutable(),
		field.UUID("organization_id", uuid.UUID{}).
			Comment("Organization that hosts this event"),
		field.String("title").
			NotEmpty().
			Comment("Event title"),
		field.Text("description").
			Optional().
			Comment("Event description"),
		field.String("location").
			Optional().
			Comment("Event location"),
		field.String("venue").
			Optional().
			Comment("Event venue name"),
		field.Time("start_time").
			Comment("Event start time"),
		field.Time("end_time").
			Comment("Event end time"),
		field.Int("total_tickets").
			Default(0).
			NonNegative().
			Comment("Total number of tickets available"),
		field.Int("available_tickets").
			Default(0).
			NonNegative().
			Comment("Number of tickets still available"),
		field.Float("ticket_price").
			Default(0.0).
			Min(0).
			Comment("Price per ticket"),
		field.String("currency").
			Default("KRW").
			Comment("Currency code (e.g., KRW, USD)"),
		field.String("thumbnail_url").
			Optional().
			Comment("Event thumbnail image URL"),
		field.Enum("status").
			Values("draft", "published", "ongoing", "completed", "cancelled").
			Default("draft").
			Comment("Event status"),
		field.Bool("is_public").
			Default(true).
			Comment("Whether the event is publicly visible"),
		field.UUID("created_by", uuid.UUID{}).
			Comment("User ID who created this event"),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Event.
func (Event) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("organization", Organization.Type).
			Ref("events").
			Field("organization_id").
			Required().
			Unique(),
		edge.From("creator", User.Type).
			Ref("created_events").
			Field("created_by").
			Required().
			Unique(),
	}
}
