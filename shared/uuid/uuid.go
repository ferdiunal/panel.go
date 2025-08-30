package uuid

import (
	"github.com/google/uuid"
)

// NewUUID generates a new UUID v7 if available, otherwise falls back to UUID v4
func NewUUID() uuid.UUID {
	id, err := uuid.NewV7()
	if err != nil {
		// Fallback to UUID v4 if v7 generation fails
		return uuid.New()
	}
	return id
}
