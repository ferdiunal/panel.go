package account

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/core"
	"github.com/ferdiunal/panel.go/pkg/fields"
)

// AccountFieldResolver, Account alanlarını çözer
type AccountFieldResolver struct{}

// ResolveFields, Account alanlarını döner
func (r *AccountFieldResolver) ResolveFields(ctx *context.Context) []core.Element {
	return []core.Element{
		(&fields.Schema{
			Key:   "id",
			Name:  "ID",
			View:  "text",
			Props: make(map[string]interface{}),
		}).ReadOnly().OnlyOnDetail(),

		(&fields.Schema{
			Key:   "provider_id",
			Name:  "Provider",
			View:  "text",
			Props: make(map[string]interface{}),
		}).OnList().OnDetail().OnForm(),

		(&fields.Schema{
			Key:   "account_id",
			Name:  "Account ID",
			View:  "text",
			Props: make(map[string]interface{}),
		}).OnList().OnDetail().OnForm(),

		(&fields.Schema{
			Key:   "user_id",
			Name:  "User",
			View:  "belongs-to",
			Props: make(map[string]interface{}),
		}).OnList().OnDetail().OnForm(),

		(&fields.Schema{
			Key:   "access_token",
			Name:  "Access Token",
			View:  "text",
			Props: make(map[string]interface{}),
		}).ReadOnly().OnlyOnDetail(),

		(&fields.Schema{
			Key:   "refresh_token",
			Name:  "Refresh Token",
			View:  "text",
			Props: make(map[string]interface{}),
		}).ReadOnly().OnlyOnDetail(),

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
