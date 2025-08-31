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

// User holds the schema definition for the User entity.
type Account struct {
	ent.Schema
}

// Fields of the User.
func (Account) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(_uuid.NewUUID).
			Unique().
			Immutable(),
		field.String("provider").Optional(),
		field.String("provider_id").Unique().Optional(),
		field.UUID("user_id", uuid.UUID{}),
		field.String("access_token").Optional(),
		field.String("refresh_token").Optional(),
		field.String("id_token").Optional(),
		field.Time("access_token_expires_at").Optional(),
		field.Time("refresh_token_expires_at").Optional(),
		field.JSON("scopes", []string{}).Optional(),
		field.String("password").Sensitive().Optional(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the User.
func (Account) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("accounts").
			Field("user_id").
			Unique().
			Required(),
	}
}

func (Account) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id"),
	}
}
