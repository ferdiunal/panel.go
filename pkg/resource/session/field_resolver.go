package session

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/core"
	"github.com/ferdiunal/panel.go/pkg/fields"
)

// SessionFieldResolver, Session alanlarını çözer
type SessionFieldResolver struct{}

// ResolveFields, Session alanlarını döner
func (r *SessionFieldResolver) ResolveFields(ctx *context.Context) []core.Element {
	return []core.Element{
		fields.ID("ID").ReadOnly().OnlyOnDetail(),

		fields.Text("Token", "token").ReadOnly().OnList().OnDetail(),

		// User relationship
		fields.Link("User", "users", "user").OnList().OnDetail(),

		fields.Text("IP Address", "ipAddress").OnList().OnDetail(),

		fields.Text("User Agent", "userAgent").OnDetail(),

		fields.DateTime("Expires At", "expiresAt").OnList().OnDetail(),

		fields.DateTime("Created At", "createdAt").ReadOnly().OnList().OnDetail(),

		fields.DateTime("Updated At", "updatedAt").ReadOnly().OnList().OnDetail(),
	}
}
