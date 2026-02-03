package user

import (
	"context"
	"time"
)

type User struct {
	ID            string    `json:"id" gorm:"primaryKey;type:uuid"`
	Name          string    `json:"name" gorm:"index"`
	Email         string    `json:"email" gorm:"uniqueIndex"`
	EmailVerified bool      `json:"emailVerified"`
	Image         string    `json:"image"`
	Role          string    `json:"role" gorm:"index"`
	CreatedAt     time.Time `json:"createdAt" gorm:"index"`
	UpdatedAt     time.Time `json:"updatedAt" gorm:"index"`
}

type Repository interface {
	CreateUser(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	Delete(ctx context.Context, id string) error
}
