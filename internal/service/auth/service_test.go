package auth

import (
	"context"
	"fmt"
	"testing"

	"github.com/ferdiunal/panel.go/internal/data/orm"
	"github.com/ferdiunal/panel.go/internal/domain/account"
	"github.com/ferdiunal/panel.go/internal/domain/session"
	"github.com/ferdiunal/panel.go/internal/domain/user"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newTestService(t *testing.T) *Service {
	t.Helper()

	// Her test için unique database oluştur
	db, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	if err := db.AutoMigrate(&user.User{}, &account.Account{}, &session.Session{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	return NewService(
		db,
		orm.NewUserRepository(db),
		orm.NewSessionRepository(db),
		orm.NewAccountRepository(db),
	)
}

func TestRegisterEmailValidation(t *testing.T) {
	service := newTestService(t)
	ctx := context.Background()

	if _, err := service.RegisterEmail(ctx, "Test", "invalid", "Password1"); err != ErrInvalidEmail {
		t.Fatalf("expected ErrInvalidEmail, got %v", err)
	}

	if _, err := service.RegisterEmail(ctx, "Test", "test@example.com", "weak"); err != ErrWeakPassword {
		t.Fatalf("expected ErrWeakPassword, got %v", err)
	}
}

func TestRegisterFirstUserBecomesAdmin(t *testing.T) {
	service := newTestService(t)
	ctx := context.Background()

	u, err := service.RegisterEmail(ctx, "Admin", "admin@example.com", "Password1")
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}
	if u.Role != user.RoleAdmin {
		t.Fatalf("expected first user role admin, got %s", u.Role)
	}
}

func TestLoginBruteForceLockout(t *testing.T) {
	service := newTestService(t)
	ctx := context.Background()

	_, err := service.RegisterEmail(ctx, "User", "user@example.com", "Password1")
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	for i := 0; i < maxFailedAttempts; i++ {
		_, err = service.LoginEmail(ctx, "user@example.com", "WrongPass9", "127.0.0.1", "test")
		if err == nil {
			t.Fatalf("expected login to fail on iteration %d", i)
		}
	}

	_, err = service.LoginEmail(ctx, "user@example.com", "Password1", "127.0.0.1", "test")
	if err != ErrTooManyAttempts {
		t.Fatalf("expected ErrTooManyAttempts, got %v", err)
	}
}

func TestLoginSuccess(t *testing.T) {
	service := newTestService(t)
	ctx := context.Background()

	_, err := service.RegisterEmail(ctx, "User", "user@example.com", "Password123")
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	sess, err := service.LoginEmail(ctx, "user@example.com", "Password123", "127.0.0.1", "test-agent")
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}

	if sess == nil {
		t.Fatal("expected session, got nil")
	}
	if sess.Token == "" {
		t.Fatal("expected session token")
	}
	if sess.User == nil {
		t.Fatal("expected user in session")
	}
	if sess.IPAddress != "127.0.0.1" {
		t.Errorf("expected IP 127.0.0.1, got %s", sess.IPAddress)
	}
	if sess.UserAgent != "test-agent" {
		t.Errorf("expected user agent test-agent, got %s", sess.UserAgent)
	}
}

func TestSessionValidation(t *testing.T) {
	service := newTestService(t)
	ctx := context.Background()

	u, err := service.RegisterEmail(ctx, "User", "user@example.com", "Password123")
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	sess, err := service.LoginEmail(ctx, "user@example.com", "Password123", "127.0.0.1", "test")
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}

	validSess, err := service.ValidateSession(ctx, sess.Token)
	if err != nil {
		t.Fatalf("validate session failed: %v", err)
	}
	if validSess.UserID != u.ID {
		t.Errorf("expected user ID %s, got %s", u.ID, validSess.UserID)
	}

	invalidSess, err := service.ValidateSession(ctx, "invalid-token")
	if err == nil {
		t.Fatal("expected error for invalid token")
	}
	if invalidSess != nil {
		t.Fatal("expected nil session for invalid token")
	}
}

func TestLogout(t *testing.T) {
	service := newTestService(t)
	ctx := context.Background()

	_, err := service.RegisterEmail(ctx, "User", "user@example.com", "Password123")
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	sess, err := service.LoginEmail(ctx, "user@example.com", "Password123", "127.0.0.1", "test")
	if err != nil {
		t.Fatalf("login failed: %v", err)
	}

	err = service.Logout(ctx, sess.Token)
	if err != nil {
		t.Fatalf("logout failed: %v", err)
	}

	_, err = service.ValidateSession(ctx, sess.Token)
	if err == nil {
		t.Fatal("expected error after logout")
	}
}

func TestEmailNormalization(t *testing.T) {
	service := newTestService(t)
	ctx := context.Background()

	_, err := service.RegisterEmail(ctx, "User", "  USER@EXAMPLE.COM  ", "Password123")
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	sess, err := service.LoginEmail(ctx, "user@example.com", "Password123", "127.0.0.1", "test")
	if err != nil {
		t.Fatalf("login with normalized email failed: %v", err)
	}
	if sess == nil {
		t.Fatal("expected session")
	}

	sess2, err := service.LoginEmail(ctx, "  USER@EXAMPLE.COM  ", "Password123", "127.0.0.1", "test")
	if err != nil {
		t.Fatalf("login with uppercase email failed: %v", err)
	}
	if sess2 == nil {
		t.Fatal("expected session")
	}
}

func TestPasswordStrength(t *testing.T) {
	service := newTestService(t)
	ctx := context.Background()

	tests := []struct {
		name     string
		password string
		wantErr  error
	}{
		{"too short", "Pass1", ErrWeakPassword},
		{"with space", "Pass word1", ErrWeakPassword},
		{"with tab", "Pass\tword1", ErrWeakPassword},
		{"with newline", "Pass\nword1", ErrWeakPassword},
		{"valid simple", "Password1", nil},
		{"valid long", "VeryLongPassword123", nil},
		{"valid special chars", "P@ssw0rd!", nil},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Unique ve geçerli email oluştur
			email := fmt.Sprintf("test%d@example.com", i)
			_, err := service.RegisterEmail(ctx, "User", email, tt.password)
			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("expected error %v, got %v", tt.wantErr, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			}
		})
	}
}

func TestDuplicateEmailRegistration(t *testing.T) {
	service := newTestService(t)
	ctx := context.Background()

	_, err := service.RegisterEmail(ctx, "User1", "duplicate@example.com", "Password123")
	if err != nil {
		t.Fatalf("first registration failed: %v", err)
	}

	_, err = service.RegisterEmail(ctx, "User2", "duplicate@example.com", "Password456")
	if err != ErrEmailAlreadyExists {
		t.Fatalf("expected ErrEmailAlreadyExists, got %v", err)
	}
}
