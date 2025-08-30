package user_resource

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"panel.go/internal/constants"
	"panel.go/internal/ent"
	"panel.go/internal/interfaces/resource"
)

type initResource struct {
}

type UserResource struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Avatar    string    `json:"avatar"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewResource() resource.ResourceInterface[ent.User, *UserResource] {
	return &initResource{}
}

func (initResource) Resource(user *ent.User) *UserResource {
	return &UserResource{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Avatar:    fmt.Sprintf("%s/avatar/%s?vibe=ocean", constants.APP_URL, user.ID),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func (r *initResource) Collection(users []*ent.User) []*UserResource {
	response := []*UserResource{}
	for _, user := range users {
		response = append(response, r.Resource(user))
	}
	return response
}
