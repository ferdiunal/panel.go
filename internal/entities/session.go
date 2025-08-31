package entities

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	_uuid "panel.go/shared/uuid"
)

// Session holds the schema definition for the Session entity.
type Session struct {
	ent.Schema
}

// Fields of the Session.
func (Session) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(_uuid.NewUUID).
			Unique().
			Immutable(),
		field.Time("expires_at").Optional(),
		field.String("token").Unique().Optional(),
		field.UUID("user_id", uuid.UUID{}),
		field.String("ip_address").Optional(),
		field.String("user_agent").Optional(),
		field.UUID("impersonator_id", uuid.UUID{}).Optional(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the Session.
func (Session) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("sessions").
			Field("user_id").
			Unique().
			Required(),
	}
}

func (Session) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("token", "user_id").Unique(),
	}
}
