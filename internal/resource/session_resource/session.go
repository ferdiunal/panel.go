package session_resource

import (
	"time"

	"github.com/google/uuid"
	"panel.go/internal/ent"
	"panel.go/internal/interfaces/resource"
)

type SessionResource struct {
	ID        uuid.UUID `json:"id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	IPAddress string    `json:"ip_address"`
	UserAgent string    `json:"user_agent"`
}

type initResource struct {
}

func NewResource() resource.ResourceInterface[ent.Session, *SessionResource] {
	return &initResource{}
}

func (initResource) Resource(session *ent.Session) *SessionResource {
	return &SessionResource{
		ID:        session.ID,
		Token:     session.Token,
		ExpiresAt: session.ExpiresAt,
		IPAddress: session.IPAddress,
		UserAgent: session.UserAgent,
	}
}

func (r *initResource) Collection(sessions []*ent.Session) []*SessionResource {
	response := []*SessionResource{}
	for _, session := range sessions {
		response = append(response, r.Resource(session))
	}
	return response
}
