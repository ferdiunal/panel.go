package account

import (
	"context"
	"time"

	"github.com/ferdiunal/panel.go/pkg/domain/user"
)

type Account struct {
	ID                    string     `json:"id" gorm:"primaryKey"`
	AccountID             string     `json:"accountId" gorm:"index"`  // Provider's user ID
	ProviderID            string     `json:"providerId" gorm:"index"` // e.g. "credential", "google"
	UserID                string     `json:"userId" gorm:"index;type:uuid"`
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
	FindByID(ctx context.Context, id string) (*Account, error)
	FindByProvider(ctx context.Context, providerID, accountID string) (*Account, error)
	FindByUserID(ctx context.Context, userID string) ([]Account, error)
	Update(ctx context.Context, account *Account) error
	Delete(ctx context.Context, id string) error
}
