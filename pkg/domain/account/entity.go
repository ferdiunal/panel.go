package account

import (
	"context"
	"time"

	"github.com/ferdiunal/panel.go/pkg/domain/user"
)

type Account struct {
	ID                    uint       `json:"id" gorm:"primaryKey"`
	AccountID             *string    `json:"accountId" gorm:"index"`  // Provider's user ID (nullable for credentials)
	ProviderID            string     `json:"providerId" gorm:"index"` // e.g. "credential", "google"
	UserID                uint       `json:"userId" gorm:"index"`
	AccessToken           string     `json:"accessToken,omitempty"`
	RefreshToken          string     `json:"refreshToken,omitempty"`
	IDToken               string     `json:"idToken,omitempty"`
	AccessTokenExpiresAt  *time.Time `json:"accessTokenExpiresAt,omitempty"`
	RefreshTokenExpiresAt *time.Time `json:"refreshTokenExpiresAt,omitempty"`
	Password              string     `json:"-"` // Hashed password, never returned
	Scope                 string     `json:"scope,omitempty"`
	CreatedAt             time.Time  `json:"createdAt" gorm:"index"`
	UpdatedAt             time.Time  `json:"updatedAt" gorm:"index"`
	User                  *user.User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

type Repository interface {
	Create(ctx context.Context, account *Account) error
	FindByID(ctx context.Context, id uint) (*Account, error)
	FindByProvider(ctx context.Context, providerID, accountID string) (*Account, error)
	FindByUserID(ctx context.Context, userID uint) ([]Account, error)
	Update(ctx context.Context, account *Account) error
	Delete(ctx context.Context, id uint) error
}
