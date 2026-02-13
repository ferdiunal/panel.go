package auth

import (
	"context"
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

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
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
