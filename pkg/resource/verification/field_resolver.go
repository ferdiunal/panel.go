package verification

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/core"
	"github.com/ferdiunal/panel.go/pkg/fields"
)

// VerificationFieldResolver, Verification alanlarını çözer
type VerificationFieldResolver struct{}

// ResolveFields, Verification alanlarını döner
func (r *VerificationFieldResolver) ResolveFields(ctx *context.Context) []core.Element {
	return []core.Element{
		(&fields.Schema{
			Key:   "id",
			Name:  "ID",
			View:  "text",
			Props: make(map[string]interface{}),
		}).ReadOnly().OnlyOnDetail(),

		(&fields.Schema{
			Key:   "identifier",
			Name:  "Identifier",
			View:  "text",
			Props: make(map[string]interface{}),
		}).OnList().OnDetail().OnForm(),

		(&fields.Schema{
			Key:   "token",
			Name:  "Token",
			View:  "text",
			Props: make(map[string]interface{}),
		}).ReadOnly().OnList().OnDetail(),

		(&fields.Schema{
			Key:   "expires_at",
			Name:  "Expires At",
			View:  "datetime",
			Props: make(map[string]interface{}),
		}).OnList().OnDetail().OnForm(),

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
