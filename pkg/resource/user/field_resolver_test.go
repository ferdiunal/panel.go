package user

import (
	"testing"

	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/core"
)

// TestUserFieldResolverResolveFields, field resolver'ı test eder
func TestUserFieldResolverResolveFields(t *testing.T) {
	resolver := &UserFieldResolver{}

	// Nil context ile
	result := resolver.ResolveFields(nil)

	if result == nil {
		t.Error("Expected non-nil result")
	}

	if len(result) == 0 {
		t.Error("Expected at least one field")
	}

	// Field türlerini kontrol et
	fieldTypes := make(map[string]bool)
	for _, field := range result {
		fieldTypes[field.GetKey()] = true
	}

	// Beklenen alanlar
	expectedFields := []string{"id", "image", "name", "email", "role", "password"}
	for _, expected := range expectedFields {
		if !fieldTypes[expected] {
			t.Errorf("Expected field '%s' not found", expected)
		}
	}
}

// TestUserFieldResolverResolveFieldsWithContext, context ile field resolver'ı test eder
func TestUserFieldResolverResolveFieldsWithContext(t *testing.T) {
	resolver := &UserFieldResolver{}
	ctx := &context.Context{}

	result := resolver.ResolveFields(ctx)

	if result == nil {
		t.Error("Expected non-nil result")
	}

	if len(result) == 0 {
		t.Error("Expected at least one field")
	}
}

// TestUserFieldResolverFieldProperties, field özelliklerini test eder
func TestUserFieldResolverFieldProperties(t *testing.T) {
	resolver := &UserFieldResolver{}
	fields := resolver.ResolveFields(nil)

	// ID field'ı kontrol et
	var idField core.Element
	for _, f := range fields {
		if f.GetKey() == "id" {
			idField = f
			break
		}
	}

	if idField == nil {
		t.Error("Expected ID field")
	}

	// Email field'ı kontrol et
	var emailField core.Element
	for _, f := range fields {
		if f.GetKey() == "email" {
			emailField = f
			break
		}
	}

	if emailField == nil {
		t.Error("Expected Email field")
	}
}

// TestUserFieldResolverImplementsInterface, interface'i implement ettiğini test eder
func TestUserFieldResolverImplementsInterface(t *testing.T) {
	var _ interface {
		ResolveFields(*context.Context) []core.Element
	} = (*UserFieldResolver)(nil)
}
