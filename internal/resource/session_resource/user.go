package session_resource

import (
	"panel.go/internal/ent"
	"panel.go/internal/interfaces/resource"
)

type sessionResource struct {
}

func NewResource() resource.ResourceInterface[ent.Session, resource.Response] {
	return &sessionResource{}
}

func (sessionResource) Resource(session *ent.Session) resource.Response {
	return resource.Response{}
}

func (r *sessionResource) Collection(sessions []*ent.Session) []resource.Response {
	response := []resource.Response{}
	for _, session := range sessions {
		response = append(response, r.Resource(session))
	}
	return response
}
