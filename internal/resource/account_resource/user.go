package account_resource

import (
	"panel.go/internal/ent"
	"panel.go/internal/interfaces/resource"
)

type accountResource struct {
}

func NewResource() resource.ResourceInterface[ent.Account, resource.Response] {
	return &accountResource{}
}

func (accountResource) Resource(user *ent.Account) resource.Response {
	return resource.Response{
		"id": user.ID,
	}
}

func (r *accountResource) Collection(users []*ent.Account) []resource.Response {
	response := []resource.Response{}
	for _, user := range users {
		response = append(response, r.Resource(user))
	}
	return response
}
