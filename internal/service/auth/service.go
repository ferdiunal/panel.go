package auth

import (
	"context"
	"errors"
	"flag"
	"net/mail"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ferdiunal/panel.go/internal/data/orm"
	"github.com/ferdiunal/panel.go/internal/domain/account"
	"github.com/ferdiunal/panel.go/internal/domain/session"
	"github.com/ferdiunal/panel.go/internal/domain/user"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidEmail       = errors.New("invalid email")
	ErrWeakPassword       = errors.New("password does not meet policy requirements")
	ErrInvalidName        = errors.New("invalid name")
	ErrTooManyAttempts    = errors.New("too many failed login attempts")
)

const (
	maxFailedAttempts = 5
	attemptWindow     = 15 * time.Minute
	lockoutDuration   = 15 * time.Minute
	minPasswordLength = 8
)

type loginAttemptState struct {
	Count       int
	FirstFailed time.Time
	LockedUntil time.Time
}

type Service struct {
	db          *gorm.DB
	userRepo    user.Repository
	sessionRepo session.Repository
	accountRepo account.Repository
	hashCost    int
	attemptMu   sync.Mutex
	attempts    map[string]loginAttemptState
}

func NewService(db *gorm.DB, u user.Repository, s session.Repository, a account.Repository) *Service {
	return &Service{
		db:          db,
		userRepo:    u,
		sessionRepo: s,
		accountRepo: a,
		hashCost:    resolvePasswordHashCost(),
		attempts:    make(map[string]loginAttemptState),
	}
}

func (s *Service) RegisterEmail(ctx context.Context, name, email, password string) (*user.User, error) {
	name = strings.TrimSpace(name)
	email = normalizeEmail(email)

	if name == "" {
		return nil, ErrInvalidName
	}
	if !isValidEmail(email) {
		return nil, ErrInvalidEmail
	}
	if !isStrongPassword(password) {
		return nil, ErrWeakPassword
	}

	// Keep registration atomic.
	if s.db != nil {
		var createdUser *user.User
		err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			txUserRepo := orm.NewUserRepository(tx)
			txAccountRepo := orm.NewAccountRepository(tx)

			existing, err := txUserRepo.FindByEmail(ctx, email)
			if err == nil && existing != nil {
				return ErrEmailAlreadyExists
			}
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}

			role := user.RoleUser
			var totalUsers int64
			if err := tx.Model(&user.User{}).Count(&totalUsers).Error; err == nil && totalUsers == 0 {
				role = user.RoleAdmin
			}

			hashed, err := bcrypt.GenerateFromPassword([]byte(password), s.hashCost)
			if err != nil {
				return err
			}

			userID, _ := uuid.NewV7()
			u := &user.User{
				ID:            userID.String(),
				Name:          name,
				Email:         email,
				Role:          role,
				EmailVerified: false,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			}

			if err := txUserRepo.CreateUser(ctx, u); err != nil {
				return err
			}

			accountID, _ := uuid.NewV7()
			acc := &account.Account{
				ID:         accountID.String(),
				UserID:     u.ID,
				ProviderID: "credential",
				AccountID:  email,
				Password:   string(hashed),
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}

			if err := txAccountRepo.Create(ctx, acc); err != nil {
				return err
			}

			createdUser = u
			return nil
		})
		if err != nil {
			if errors.Is(err, ErrEmailAlreadyExists) {
				return nil, err
			}
			return nil, err
		}
		return createdUser, nil
	}

	// Fallback for non-GORM setups.
	existing, err := s.userRepo.FindByEmail(ctx, email)
	if err == nil && existing != nil {
		return nil, ErrEmailAlreadyExists
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), s.hashCost)
	if err != nil {
		return nil, err
	}

	userID, _ := uuid.NewV7()
	u := &user.User{
		ID:            userID.String(),
		Name:          name,
		Email:         email,
		Role:          user.RoleUser,
		EmailVerified: false,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.userRepo.CreateUser(ctx, u); err != nil {
		return nil, err
	}

	accountID, _ := uuid.NewV7()
	acc := &account.Account{
		ID:         accountID.String(),
		UserID:     u.ID,
		ProviderID: "credential",
		AccountID:  email,
		Password:   string(hashed),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.accountRepo.Create(ctx, acc); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *Service) LoginEmail(ctx context.Context, email, password string, ip, userAgent string) (*session.Session, error) {
	email = normalizeEmail(email)
	attemptKey := buildAttemptKey(email, ip)
	if s.isLocked(attemptKey) {
		return nil, ErrTooManyAttempts
	}

	// Find User
	u, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		s.recordFailedAttempt(attemptKey)
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
			s.recordFailedAttempt(attemptKey)
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
			s.recordFailedAttempt(attemptKey)
			return nil, ErrInvalidCredentials
		}
	}

	// Verify Password
	if err := bcrypt.CompareHashAndPassword([]byte(acc.Password), []byte(password)); err != nil {
		s.recordFailedAttempt(attemptKey)
		return nil, ErrInvalidCredentials
	}

	// Create Session
	sessionId, _ := uuid.NewV7()
	sessionToken, _ := uuid.NewV7()
	sess := &session.Session{
		ID:        sessionId.String(),
		UserID:    u.ID,
		Token:     sessionToken.String(),
		User:      u,
		ExpiresAt: time.Now().Add(24 * 7 * time.Hour), // 7 days
		IPAddress: ip,
		UserAgent: userAgent,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.sessionRepo.Create(ctx, sess); err != nil {
		return nil, err
	}

	s.clearFailedAttempts(attemptKey)
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

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func isValidEmail(email string) bool {
	if email == "" {
		return false
	}
	_, err := mail.ParseAddress(email)
	return err == nil
}

func isStrongPassword(password string) bool {
	if len(password) < minPasswordLength {
		return false
	}

	for _, r := range password {
		if r == ' ' || r == '\t' || r == '\n' {
			return false
		}
	}
	return true
}

func buildAttemptKey(email, ip string) string {
	return email + "|" + strings.TrimSpace(ip)
}

func (s *Service) isLocked(key string) bool {
	s.attemptMu.Lock()
	defer s.attemptMu.Unlock()

	state, ok := s.attempts[key]
	if !ok {
		return false
	}

	now := time.Now()
	if !state.LockedUntil.IsZero() && now.Before(state.LockedUntil) {
		return true
	}
	if !state.LockedUntil.IsZero() && !now.Before(state.LockedUntil) {
		delete(s.attempts, key)
	}

	return false
}

func (s *Service) recordFailedAttempt(key string) {
	s.attemptMu.Lock()
	defer s.attemptMu.Unlock()

	now := time.Now()
	state := s.attempts[key]

	if state.FirstFailed.IsZero() || now.Sub(state.FirstFailed) > attemptWindow {
		state = loginAttemptState{
			Count:       1,
			FirstFailed: now,
		}
		s.attempts[key] = state
		return
	}

	state.Count++
	if state.Count >= maxFailedAttempts {
		state.LockedUntil = now.Add(lockoutDuration)
	}
	s.attempts[key] = state
}

func (s *Service) clearFailedAttempts(key string) {
	s.attemptMu.Lock()
	defer s.attemptMu.Unlock()
	delete(s.attempts, key)
}

func resolvePasswordHashCost() int {
	// Keep tests and race runs fast and deterministic.
	if flag.Lookup("test.v") != nil || strings.HasSuffix(os.Args[0], ".test") {
		return bcrypt.MinCost
	}
	return bcrypt.DefaultCost
}
