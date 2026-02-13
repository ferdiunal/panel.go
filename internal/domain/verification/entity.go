package verification

import (
	"context"
	"time"
)

type Verification struct {
	ID         string    `json:"id" gorm:"primaryKey"`
	Identifier string    `json:"identifier" gorm:"index"` // Email or UserID
	Token      string    `json:"token" gorm:"uniqueIndex"`
	ExpiresAt  time.Time `json:"expiresAt" gorm:"index"`
	CreatedAt  time.Time `json:"createdAt" gorm:"index"`
	UpdatedAt  time.Time `json:"updatedAt" gorm:"index"`
}

type Repository interface {
	Create(ctx context.Context, v *Verification) error
	FindByToken(ctx context.Context, token string) (*Verification, error)
	Delete(ctx context.Context, id string) error
	DeleteByIdentifier(ctx context.Context, identifier string) error
}
