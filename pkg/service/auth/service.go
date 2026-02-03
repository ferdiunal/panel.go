package auth

import (
	"context"
	"errors"
	"time"

	"github.com/ferdiunal/panel.go/pkg/domain/account"
	"github.com/ferdiunal/panel.go/pkg/domain/session"
	"github.com/ferdiunal/panel.go/pkg/domain/user"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailAlreadyExists = errors.New("email already exists")
)

type Service struct {
	userRepo    user.Repository
	sessionRepo session.Repository
	accountRepo account.Repository
}

func NewService(u user.Repository, s session.Repository, a account.Repository) *Service {
	return &Service{
		userRepo:    u,
		sessionRepo: s,
		accountRepo: a,
	}
}

func (s *Service) RegisterEmail(ctx context.Context, name, email, password string) (*user.User, error) {
	// Check if user exists
	existing, _ := s.userRepo.FindByEmail(ctx, email)
	if existing != nil {
		return nil, ErrEmailAlreadyExists
	}

	// Hash Password
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create User
	userId, _ := uuid.NewV7()
	u := &user.User{
		ID:            userId.String(),
		Name:          name,
		Email:         email,
		EmailVerified: false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.userRepo.CreateUser(ctx, u); err != nil {
		return nil, err
	}

	// Create Account
	accountId, _ := uuid.NewV7()
	acc := &account.Account{
		ID:         accountId.String(),
		UserID:     u.ID,
		ProviderID: "credential",
		AccountID:  email, // For credential provider, accountID is usually the email or a unique identifier
		Password:   string(hashed),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.accountRepo.Create(ctx, acc); err != nil {
		// Cleanup user if account creation fails?
		// For simplicity, we assume transaction or manual cleanup, but keeping it simple here.
		return nil, err
	}

	return u, nil
}

func (s *Service) LoginEmail(ctx context.Context, email, password string, ip, userAgent string) (*session.Session, error) {
	// Find User
	u, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Find Credential Account
	// We could query account repo directly or check user if we had loaded accounts,
	// but let's query account repo for specific provider
	acc, err := s.accountRepo.FindByProvider(ctx, "credential", email)
	if err != nil {
		// Fallback: try finding any credential account for this user if email is just login identifier
		accounts, err := s.accountRepo.FindByUserID(ctx, u.ID)
		if err != nil {
			return nil, ErrInvalidCredentials
		}
		found := false
		for _, a := range accounts {
			if a.ProviderID == "credential" {
				acc = &a
				found = true
				break
			}
		}
		if !found {
			return nil, ErrInvalidCredentials
		}
	}

	// Verify Password
	if err := bcrypt.CompareHashAndPassword([]byte(acc.Password), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Create Session
	sessionId, _ := uuid.NewV7()
	sessionToken, _ := uuid.NewV7()
	sess := &session.Session{
		ID:        sessionId.String(),
		UserID:    u.ID,
		Token:     sessionToken.String(),
		ExpiresAt: time.Now().Add(24 * 7 * time.Hour), // 7 days
		IPAddress: ip,
		UserAgent: userAgent,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.sessionRepo.Create(ctx, sess); err != nil {
		return nil, err
	}

	return sess, nil
}

func (s *Service) ValidateSession(ctx context.Context, token string) (*session.Session, error) {
	sess, err := s.sessionRepo.FindByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	if sess.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("session expired")
	}

	// Optionally extend session?
	return sess, nil
}

func (s *Service) Logout(ctx context.Context, token string) error {
	return s.sessionRepo.DeleteByToken(ctx, token)
}
