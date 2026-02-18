package panel

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ferdiunal/panel.go/pkg/resource"
	resourceAccount "github.com/ferdiunal/panel.go/pkg/resource/account"
	resourceSession "github.com/ferdiunal/panel.go/pkg/resource/session"
	resourceUser "github.com/ferdiunal/panel.go/pkg/resource/user"
	resourceVerification "github.com/ferdiunal/panel.go/pkg/resource/verification"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestE2EInitEndpoint, /api/internal/init endpoint'ini test eder
func TestE2EInitEndpoint(t *testing.T) {
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
		},
	}

	p := New(cfg)

	// Test /api/internal/init endpoint
	req := httptest.NewRequest("GET", "/api/internal/init", nil)
	resp, err := testFiberRequest(p.Fiber, req)

	if err != nil {
		t.Errorf("Failed to test /api/internal/init: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

// TestE2ENavigationEndpoint, /api/internal/navigation endpoint'ini test eder
func TestE2ENavigationEndpoint(t *testing.T) {
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
		},
	}

	p := New(cfg)

	// Test /api/internal/navigation endpoint
	req := httptest.NewRequest("GET", "/api/internal/navigation", nil)
	resp, err := testFiberRequest(p.Fiber, req)

	if err != nil {
		t.Errorf("Failed to test /api/internal/navigation: %v", err)
	}

	// Should be 401 (unauthorized) or 200 (if auth middleware is not enforced in test)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 200 or 401, got %d", resp.StatusCode)
	}
}

// TestE2EResourceListEndpoint, /api/internal/resource/:resource endpoint'ini test eder
func TestE2EResourceListEndpoint(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&resourceUser.UserResource{})

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

	// Test /api/internal/resource/users endpoint
	req := httptest.NewRequest("GET", "/api/internal/resource/users", nil)
	resp, err := testFiberRequest(p.Fiber, req)

	if err != nil {
		t.Errorf("Failed to test /api/internal/resource/users: %v", err)
	}

	// Should be 401 (unauthorized) or 200 (if auth middleware is not enforced in test)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 200 or 401, got %d", resp.StatusCode)
	}
}

// TestE2EResourceCreateEndpoint, /api/internal/resource/:resource/create endpoint'ini test eder
func TestE2EResourceCreateEndpoint(t *testing.T) {
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

	// Test /api/internal/resource/users/create endpoint
	req := httptest.NewRequest("GET", "/api/internal/resource/users/create", nil)
	resp, err := testFiberRequest(p.Fiber, req)

	if err != nil {
		t.Errorf("Failed to test /api/internal/resource/users/create: %v", err)
	}

	// Should be 401 (unauthorized) or 200 (if auth middleware is not enforced in test)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 200 or 401, got %d", resp.StatusCode)
	}
}

// TestE2EResourceCountEndpoint, resource sayısını test eder
func TestE2EResourceCountEndpoint(t *testing.T) {
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
		},
	}

	p := New(cfg)

	// Check resource count
	if len(p.resources) != 4 {
		t.Errorf("Expected 4 resources, got %d", len(p.resources))
	}
}

// TestE2EResourceSlugsEndpoint, resource slug'larını test eder
func TestE2EResourceSlugsEndpoint(t *testing.T) {
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
		},
	}

	p := New(cfg)

	// Check resource slugs
	expectedSlugs := []string{"accounts", "sessions", "verifications", "users"}
	for _, slug := range expectedSlugs {
		if _, ok := p.resources[slug]; !ok {
			t.Errorf("Expected resource with slug '%s'", slug)
		}
	}
}

// TestE2EResourcePropertiesEndpoint, resource özelliklerini test eder
func TestE2EResourcePropertiesEndpoint(t *testing.T) {
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
		t.Error("Expected account resource")
		return
	}

	// Check properties
	properties := map[string]interface{}{
		"slug":  accountRes.Slug(),
		"title": accountRes.Title(),
		"icon":  accountRes.Icon(),
		"group": accountRes.Group(),
	}

	// Try to marshal to JSON
	data, err := json.Marshal(properties)
	if err != nil {
		t.Errorf("Failed to marshal properties: %v", err)
	}

	if len(data) == 0 {
		t.Error("Expected non-empty JSON data")
	}
}

// TestE2EResourceFieldsEndpoint, resource alanlarını test eder
func TestE2EResourceFieldsEndpoint(t *testing.T) {
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
		},
	}

	p := New(cfg)

	// Check fields for each resource
	for slug, res := range p.resources {
		if slug == "users" {
			continue // Skip user resource
		}

		fields := res.Fields()
		if len(fields) == 0 {
			t.Errorf("Expected fields for resource '%s'", slug)
		}
	}
}

// TestE2EResourcePolicyEndpoint, resource policy'lerini test eder
func TestE2EResourcePolicyEndpoint(t *testing.T) {
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
		},
	}

	p := New(cfg)

	// Check policy for each resource
	for slug, res := range p.resources {
		if slug == "users" {
			continue // Skip user resource
		}

		policy := res.Policy()
		if policy == nil {
			t.Errorf("Expected policy for resource '%s'", slug)
		}
	}
}

// TestE2EResourceRepositoryEndpoint, resource repository'lerini test eder
func TestE2EResourceRepositoryEndpoint(t *testing.T) {
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
		},
	}

	p := New(cfg)

	// Check repository for each resource
	for slug, res := range p.resources {
		if slug == "users" {
			continue // Skip user resource
		}

		repo := res.Repository(db)
		if repo == nil {
			t.Errorf("Expected repository for resource '%s'", slug)
		}
	}
}
