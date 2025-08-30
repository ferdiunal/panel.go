package service

import (
	"context"

	"panel.go/internal/repository"
)

type AuthService struct {
	AccountRepository repository.AccountRepository
	UserRepository    repository.UserRepository
	SessionRepository repository.SessionRepository
}

type LoginPayload struct {
	Email    string `json:"email" validate:"required,email,max=150"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type RegisterPayload struct {
	Name     string `json:"name" validate:"required,min=3,max=150"`
	Email    string `json:"email" validate:"required,email,max=150"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

func NewAuthService(
	accountRepository repository.AccountRepository,
	userRepository repository.UserRepository,
	sessionRepository repository.SessionRepository,
) *AuthService {
	return &AuthService{AccountRepository: accountRepository, UserRepository: userRepository, SessionRepository: sessionRepository}
}

func (s *AuthService) Login(ctx context.Context, email string, password string) (*string, error) {
	return nil, nil
}
