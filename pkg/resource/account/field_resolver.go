package account

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/fields"
)

// AccountFieldResolver, Account alanlarını çözer
type AccountFieldResolver struct{}

// ResolveFields, Account alanlarını döner
func (r *AccountFieldResolver) ResolveFields(ctx *context.Context) []fields.Element {
	return []fields.Element{
		fields.ID("ID").ReadOnly(),

		fields.Text("Provider", "providerId").ReadOnly(),

		// fields.Text("Account ID", "accountId").ReadOnly(),

		// User relationship
		// Using "user" key to map to Account.User struct
		fields.Link("User", "users", "user").
			ReadOnly(),

		fields.Text("Provider", "provider").
			ReadOnly(),

		// fields.Text("Access Token", "accessToken").ReadOnly(),

		// fields.Text("Refresh Token", "refreshToken").ReadOnly(),

		fields.DateTime("Created At", "createdAt").ReadOnly(),

		fields.DateTime("Updated At", "updatedAt").ReadOnly(),
	}
}
