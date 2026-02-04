package user

import (
	"testing"

	domainUser "github.com/ferdiunal/panel.go/pkg/domain/user"
	"github.com/ferdiunal/panel.go/pkg/resource"
)

// TestNewUserResource, UserResource oluşturulmasını test eder
func TestNewUserResource(t *testing.T) {
	r := NewUserResource()

	if r == nil {
		t.Fatal("Expected UserResource, got nil")
	}

	if r.Slug() != "users" {
		t.Errorf("Expected slug 'users', got '%s'", r.Slug())
	}

	if r.Title() != "Users" {
		t.Errorf("Expected title 'Users', got '%s'", r.Title())
	}

	if r.Icon() != "users" {
		t.Errorf("Expected icon 'users', got '%s'", r.Icon())
	}

	if r.Group() != "System" {
		t.Errorf("Expected group 'System', got '%s'", r.Group())
	}

	if !r.Visible() {
		t.Error("Expected resource to be visible")
	}

	if r.NavigationOrder() != 1 {
		t.Errorf("Expected navigation order 1, got %d", r.NavigationOrder())
	}
}

// TestUserResourceModel, model'in doğru ayarlandığını test eder
func TestUserResourceModel(t *testing.T) {
	r := NewUserResource()
	model := r.Model()

	if model == nil {
		t.Fatal("Expected model, got nil")
	}

	_, ok := model.(*domainUser.User)
	if !ok {
		t.Errorf("Expected *domainUser.User, got %T", model)
	}
}

// TestUserResourceFields, alanların çözümlendiğini test eder
func TestUserResourceFields(t *testing.T) {
	r := NewUserResource()
	fields := r.Fields()

	if fields == nil {
		t.Fatal("Expected fields, got nil")
	}

	if len(fields) == 0 {
		t.Fatal("Expected fields, got empty slice")
	}
}

// TestUserResourceImplementsResource, Resource interface'ini implement ettiğini test eder
func TestUserResourceImplementsResource(t *testing.T) {
	r := NewUserResource()

	var _ resource.Resource = r

	// Temel metodları test et
	_ = r.Model()
	_ = r.Fields()
	_ = r.Slug()
	_ = r.Title()
	_ = r.Cards()
	_ = r.Visible()
}

// TestUserResourceBackwardCompatibility, eski GetUserResource() fonksiyonunun çalıştığını test eder
func TestUserResourceBackwardCompatibility(t *testing.T) {
	r := GetUserResource()

	if r == nil {
		t.Fatal("Expected resource, got nil")
	}

	if r.Slug() != "users" {
		t.Errorf("Expected slug 'users', got '%s'", r.Slug())
	}
}

// TestUserResourceSetDialogType, dialog type'ı ayarlamayı test eder
func TestUserResourceSetDialogType(t *testing.T) {
	r := NewUserResource()
	r.SetDialogType(resource.DialogTypeModal)

	if r.GetDialogType() != resource.DialogTypeModal {
		t.Errorf("Expected DialogTypeModal, got %v", r.GetDialogType())
	}
}
