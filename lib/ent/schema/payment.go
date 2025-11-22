package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Payment holds the schema definition for the Payment entity.
type Payment struct {
	ent.Schema
}

// Fields of the Payment.
func (Payment) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Unique().
			Immutable(),
		field.UUID("event_id", uuid.UUID{}).
			Comment("Event ID for this payment"),
		field.UUID("user_id", uuid.UUID{}).
			Optional().
			Comment("User ID who made this payment (optional for guest checkout)"),
		field.String("event_title").
			NotEmpty().
			Comment("Event title at the time of purchase"),
		field.Int("ticket_quantity").
			Positive().
			Comment("Number of tickets purchased"),
		field.Float("total_price").
			Min(0).
			Comment("Total price paid"),
		field.String("currency").
			Default("KRW").
			Comment("Currency code"),
		field.String("buyer_name").
			NotEmpty().
			Comment("Buyer's name"),
		field.String("buyer_email").
			NotEmpty().
			Comment("Buyer's email"),
		field.String("buyer_phone").
			NotEmpty().
			Comment("Buyer's phone number"),
		field.String("payment_key").
			Optional().
			Comment("Payment gateway payment key"),
		field.String("order_id").
			Optional().
			Unique().
			Comment("Order ID for payment gateway"),
		field.Enum("status").
			Values("pending", "completed", "failed", "cancelled", "refunded").
			Default("pending").
			Comment("Payment status"),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Payment.
func (Payment) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("event", Event.Type).
			Ref("payments").
			Field("event_id").
			Required().
			Unique(),
		edge.From("user", User.Type).
			Ref("payments").
			Field("user_id").
			Unique(),
	}
}
