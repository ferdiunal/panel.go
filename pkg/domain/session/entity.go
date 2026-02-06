package session

import (
	"context"
	"time"

	"github.com/ferdiunal/panel.go/pkg/domain/user"
)

type Session struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	Token     string     `json:"token" gorm:"uniqueIndex"`
	UserID    uint       `json:"userId" gorm:"index"`
	ExpiresAt time.Time  `json:"expiresAt" gorm:"index"`
	IPAddress string     `json:"ipAddress"`
	UserAgent string     `json:"userAgent"`
	CreatedAt time.Time  `json:"createdAt" gorm:"index"`
	UpdatedAt time.Time  `json:"updatedAt" gorm:"index"`
	User      *user.User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

type Repository interface {
	Create(ctx context.Context, session *Session) error
	FindByID(ctx context.Context, id uint) (*Session, error)
	FindByToken(ctx context.Context, token string) (*Session, error)
	Delete(ctx context.Context, id uint) error
	DeleteByToken(ctx context.Context, token string) error
	DeleteByUserID(ctx context.Context, userID uint) error
}
