package panel

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	appContext "github.com/ferdiunal/panel.go/pkg/context" // Alias for our internal context
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/page"
	"github.com/ferdiunal/panel.go/pkg/widget"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Mock Page Implementation
type MockPage struct {
	page.Base
}

func (p *MockPage) Slug() string {
	return "mock-dashboard"
}

func (p *MockPage) Title() string {
	return "Mock Dashboard"
}

func (p *MockPage) Cards() []widget.Card {
	// Return a simple card mock if needed, or empty
	return []widget.Card{}
}

func (p *MockPage) Fields() []fields.Element {
	return []fields.Element{}
}

func (p *MockPage) Save(c *appContext.Context, db *gorm.DB, data map[string]interface{}) error {
	return nil
}

func TestPageRegistration(t *testing.T) {
	// Setup app
	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	app := New(Config{
		Database:    DatabaseConfig{Instance: db},
		Environment: "test",
	})

	// Register Page
	mockPage := &MockPage{}
	app.RegisterPage(mockPage)

	// Verify Registration
	if _, ok := app.pages["mock-dashboard"]; !ok {
		t.Error("Page not registered correctly")
	}
}

func TestPageAPI(t *testing.T) {
	// Setup app

	// Ensure Auth Domains Migration for Session Middleware
	// In integration tests we might need user/session tables
	// But here we might check if middleware blocks us or if we can bypass it/mock it.
	// Panel Auth Middleware is applied to /api group.

	// We need to bypass auth or login relative to how New() sets up routes.
	// New() adds SessionMiddleware to /api group.
	// So we need a valid session or mock context.

	// Ideally we can test handlers directly or use app.Test with valid cookie.
	// For simplicity, let's test the handlers directly via appContext wrapper if possible,
	// or register a route without middleware for testing (not easy with New()).

	// Strategy: Use a fresh Fiber app and register handlers directly to avoid auth middleware for unit testing logic.

	testApp := fiber.New()
	p := &Panel{
		pages: make(map[string]page.Page),
	}

	// Register Mock Page
	mockPage := &MockPage{}
	p.RegisterPage(mockPage)

	// Register Routes directly
	testApp.Get("/pages", appContext.Wrap(p.handlePages))
	testApp.Get("/pages/:slug", appContext.Wrap(p.handlePageDetail))

	// 1. Test List
	req := httptest.NewRequest("GET", "/pages", nil)
	resp, _ := testApp.Test(req)

	if resp.StatusCode != 200 {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var listResp map[string]interface{}
	json.Unmarshal(body, &listResp)

	data := listResp["data"].([]interface{})
	if len(data) != 1 {
		t.Errorf("Expected 1 page, got %d", len(data))
	}

	item := data[0].(map[string]interface{})
	if item["slug"] != "mock-dashboard" {
		t.Errorf("Expected slug mock-dashboard, got %v", item["slug"])
	}

	// 2. Test Detail
	req = httptest.NewRequest("GET", "/pages/mock-dashboard", nil)
	resp, _ = testApp.Test(req)

	if resp.StatusCode != 200 {
		t.Errorf("Expected 200 for detail, got %d", resp.StatusCode)
	}

	body, _ = io.ReadAll(resp.Body)
	var detailResp map[string]interface{}
	json.Unmarshal(body, &detailResp)

	if detailResp["slug"] != "mock-dashboard" {
		t.Errorf("Expected slug in detail, got %v", detailResp["slug"])
	}
}

func TestDefaultAPIPageRegistration(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	app := New(Config{
		Database:    DatabaseConfig{Instance: db},
		Environment: "test",
	})

	if _, ok := app.pages["api-settings"]; !ok {
		t.Error("default api page not registered")
	}
}
