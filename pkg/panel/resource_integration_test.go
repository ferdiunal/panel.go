package panel

import (
	"encoding/json"
	"testing"

	"github.com/ferdiunal/panel.go/pkg/resource"
	resourceAccount "github.com/ferdiunal/panel.go/pkg/resource/account"
	resourceSession "github.com/ferdiunal/panel.go/pkg/resource/session"
	resourceSetting "github.com/ferdiunal/panel.go/pkg/resource/setting"
	resourceUser "github.com/ferdiunal/panel.go/pkg/resource/user"
	resourceVerification "github.com/ferdiunal/panel.go/pkg/resource/verification"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestPanelResourceRegistration, resource kayıt işlemini test eder
func TestPanelResourceRegistration(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	cfg := Config{
		Database: DatabaseConfig{
			Instance: db,
		},
		Server: ServerConfig{
			Host: "localhost",
			Port: "8080",
		},
		Environment: "test",
		Resources: []resource.Resource{
			resourceAccount.NewAccountResource(),
			resourceSession.NewSessionResource(),
			resourceVerification.NewVerificationResource(),
			resourceSetting.NewSettingResource(),
		},
	}

	p := New(cfg)

	// Check if resources are registered
	if len(p.resources) == 0 {
		t.Error("Expected resources to be registered")
	}

	// Check specific resources
	expectedSlugs := []string{"accounts", "sessions", "verifications", "settings", "users"}
	for _, slug := range expectedSlugs {
		if _, ok := p.resources[slug]; !ok {
			t.Errorf("Expected resource with slug '%s' to be registered", slug)
		}
	}
}

// TestPanelResourceCount, kayıtlı resource sayısını test eder
func TestPanelResourceCount(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	cfg := Config{
		Database: DatabaseConfig{
			Instance: db,
		},
		Server: ServerConfig{
			Host: "localhost",
			Port: "8080",
		},
		Environment: "test",
		Resources: []resource.Resource{
			resourceAccount.NewAccountResource(),
			resourceSession.NewSessionResource(),
			resourceVerification.NewVerificationResource(),
			resourceSetting.NewSettingResource(),
		},
	}

	p := New(cfg)

	// 4 custom resources + 1 default user resource = 5
	if len(p.resources) != 5 {
		t.Errorf("Expected 5 resources, got %d", len(p.resources))
	}
}

// TestPanelResourceProperties, resource özelliklerini test eder
func TestPanelResourceProperties(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	cfg := Config{
		Database: DatabaseConfig{
			Instance: db,
		},
		Server: ServerConfig{
			Host: "localhost",
			Port: "8080",
		},
		Environment: "test",
		Resources: []resource.Resource{
			resourceAccount.NewAccountResource(),
		},
	}

	p := New(cfg)

	// Get account resource
	accountRes, ok := p.resources["accounts"]
	if !ok {
		t.Error("Expected account resource to be registered")
		return
	}

	// Check properties
	if accountRes.Slug() != "accounts" {
		t.Errorf("Expected slug 'accounts', got '%s'", accountRes.Slug())
	}

	if accountRes.Title() != "Accounts" {
		t.Errorf("Expected title 'Accounts', got '%s'", accountRes.Title())
	}

	if accountRes.Icon() != "key" {
		t.Errorf("Expected icon 'key', got '%s'", accountRes.Icon())
	}

	if accountRes.Group() != "System" {
		t.Errorf("Expected group 'System', got '%s'", accountRes.Group())
	}
}

// TestPanelResourceJSON, resource'ları JSON'a serialize etme
func TestPanelResourceJSON(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	cfg := Config{
		Database: DatabaseConfig{
			Instance: db,
		},
		Server: ServerConfig{
			Host: "localhost",
			Port: "8080",
		},
		Environment: "test",
		Resources: []resource.Resource{
			resourceAccount.NewAccountResource(),
			resourceSession.NewSessionResource(),
		},
	}

	p := New(cfg)

	// Get resources
	resources := make(map[string]interface{})
	for slug, res := range p.resources {
		resources[slug] = map[string]interface{}{
			"slug":  res.Slug(),
			"title": res.Title(),
			"icon":  res.Icon(),
			"group": res.Group(),
		}
	}

	// Try to marshal to JSON
	data, err := json.Marshal(resources)
	if err != nil {
		t.Errorf("Failed to marshal resources to JSON: %v", err)
	}

	if len(data) == 0 {
		t.Error("Expected non-empty JSON data")
	}
}

// TestPanelUserResourceDefault, varsayılan User resource'u test eder
func TestPanelUserResourceDefault(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	cfg := Config{
		Database: DatabaseConfig{
			Instance: db,
		},
		Server: ServerConfig{
			Host: "localhost",
			Port: "8080",
		},
		Environment: "test",
	}

	p := New(cfg)

	// Check if user resource is registered
	userRes, ok := p.resources["users"]
	if !ok {
		t.Error("Expected user resource to be registered by default")
		return
	}

	if userRes.Slug() != "users" {
		t.Errorf("Expected slug 'users', got '%s'", userRes.Slug())
	}
}

// TestPanelCustomUserResource, özel User resource'u test eder
func TestPanelCustomUserResource(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	customUserRes := resourceUser.NewUserResource()

	cfg := Config{
		Database: DatabaseConfig{
			Instance: db,
		},
		Server: ServerConfig{
			Host: "localhost",
			Port: "8080",
		},
		Environment:  "test",
		UserResource: customUserRes,
	}

	p := New(cfg)

	// Check if custom user resource is registered
	userRes, ok := p.resources["users"]
	if !ok {
		t.Error("Expected user resource to be registered")
		return
	}

	if userRes.Slug() != "users" {
		t.Errorf("Expected slug 'users', got '%s'", userRes.Slug())
	}
}
