package user

import (
	"context"
	"time"
)

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

type User struct {
	ID            string    `json:"id" gorm:"primaryKey;type:uuid"`
	Name          string    `json:"name" gorm:"index"`
	Email         string    `json:"email" gorm:"uniqueIndex"`
	Role          string    `json:"role" gorm:"index;default:user"`
	EmailVerified bool      `json:"emailVerified"`
	Image         string    `json:"image"`
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
