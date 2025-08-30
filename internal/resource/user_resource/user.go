package user_resource

import (
	"panel.go/internal/ent"
	"panel.go/internal/interfaces/resource"
)

type userResource struct {
}

func NewResource() resource.ResourceInterface[ent.User, resource.Response] {
	return &userResource{}
}

func (userResource) Resource(user *ent.User) resource.Response {
	return resource.Response{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	}
}

func (r *userResource) Collection(users []*ent.User) []resource.Response {
	response := []resource.Response{}
	for _, user := range users {
		response = append(response, r.Resource(user))
	}
	return response
}
