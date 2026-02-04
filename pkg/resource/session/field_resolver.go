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
		(&fields.Schema{
			Key:   "id",
			Name:  "ID",
			View:  "text",
			Props: make(map[string]interface{}),
		}).ReadOnly().OnlyOnDetail(),

		(&fields.Schema{
			Key:   "token",
			Name:  "Token",
			View:  "text",
			Props: make(map[string]interface{}),
		}).ReadOnly().OnList().OnDetail(),

		(&fields.Schema{
			Key:   "user_id",
			Name:  "User",
			View:  "belongs-to",
			Props: make(map[string]interface{}),
		}).OnList().OnDetail(),

		(&fields.Schema{
			Key:   "ip_address",
			Name:  "IP Address",
			View:  "text",
			Props: make(map[string]interface{}),
		}).OnList().OnDetail(),

		(&fields.Schema{
			Key:   "user_agent",
			Name:  "User Agent",
			View:  "text",
			Props: make(map[string]interface{}),
		}).OnDetail(),

		(&fields.Schema{
			Key:   "expires_at",
			Name:  "Expires At",
			View:  "datetime",
			Props: make(map[string]interface{}),
		}).OnList().OnDetail(),

		(&fields.Schema{
			Key:   "created_at",
			Name:  "Created At",
			View:  "datetime",
			Props: make(map[string]interface{}),
		}).ReadOnly().OnList().OnDetail(),

		(&fields.Schema{
			Key:   "updated_at",
			Name:  "Updated At",
			View:  "datetime",
			Props: make(map[string]interface{}),
		}).ReadOnly().OnList().OnDetail(),
	}
}
