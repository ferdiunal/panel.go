package user

import (
	"testing"

	domainUser "github.com/ferdiunal/panel.go/pkg/domain/user"
)

// TestUserPolicyDeleteWithoutModel, model olmadan Delete permission'ını test eder
func TestUserPolicyDeleteWithoutModel(t *testing.T) {
	policy := UserPolicy{}

	// Model nil ise true döner
	if !policy.Delete(nil, nil) {
		t.Error("Expected Delete to return true when model is nil")
	}
}

// TestUserPolicyDeleteInvalidModel, geçersiz model ile Delete
func TestUserPolicyDeleteInvalidModel(t *testing.T) {
	policy := UserPolicy{}

	// Geçersiz model
	if policy.Delete(nil, "invalid") {
		t.Error("Expected Delete to return false with invalid model")
	}
}

// TestUserPolicyDeleteWithUser, user model ile Delete
func TestUserPolicyDeleteWithUser(t *testing.T) {
	policy := UserPolicy{}

	user := &domainUser.User{
		ID: 1,
	}

	// User model ile
	result := policy.Delete(nil, user)
	// nil context'te false döner
	if result {
		t.Error("Expected Delete to return false with nil context")
	}
}
