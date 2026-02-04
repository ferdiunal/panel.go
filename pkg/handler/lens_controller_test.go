package handler

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	appContext "github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/core"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/resource"
	"github.com/ferdiunal/panel.go/pkg/widget"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// MockLens implements resource.Lens interface for testing
type MockLens struct {
	name   string
	slug   string
	fields []fields.Element
}

func (m *MockLens) Name() string               { return m.name }
func (m *MockLens) Slug() string               { return m.slug }
func (m *MockLens) Fields() []fields.Element   { return m.fields }
func (m *MockLens) Query(db *gorm.DB) *gorm.DB { return db }
func (m *MockLens) GetName() string            { return m.name }
func (m *MockLens) GetSlug() string            { return m.slug }
func (m *MockLens) GetQuery() func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB { return db }
}
func (m *MockLens) GetFields(ctx *appContext.Context) []fields.Element { return m.fields }
func (m *MockLens) GetCards(ctx *appContext.Context) []widget.Card     { return nil }

// MockResourceWithLenses extends MockResource to include lenses
type MockResourceWithLenses struct {
	MockResource
	lenses []resource.Lens
}

func (m *MockResourceWithLenses) Lenses() []resource.Lens {
	return m.lenses
}

func (m *MockResourceWithLenses) GetLenses() []resource.Lens {
	return m.lenses
}

func TestHandleLensIndex_Success(t *testing.T) {
	app := fiber.New()

	lenses := []resource.Lens{
		&MockLens{name: "Active Users", slug: "active-users"},
		&MockLens{name: "Premium Users", slug: "premium-users"},
	}

	mockResource := &MockResourceWithLenses{
		lenses: lenses,
	}

	h := &FieldHandler{
		Resource: mockResource,
	}

	app.Get("/lenses", appContext.Wrap(func(c *appContext.Context) error {
		return HandleLensIndex(h, c)
	}))

	req := httptest.NewRequest("GET", "/lenses", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	dataList := response["data"].([]interface{})
	if len(dataList) != 2 {
		t.Errorf("Expected 2 lenses, got %d", len(dataList))
	}

	firstLens := dataList[0].(map[string]interface{})
	if firstLens["name"] != "Active Users" {
		t.Errorf("Expected lens name 'Active Users', got %v", firstLens["name"])
	}
	if firstLens["slug"] != "active-users" {
		t.Errorf("Expected lens slug 'active-users', got %v", firstLens["slug"])
	}
}

func TestHandleLensIndex_NoResource(t *testing.T) {
	app := fiber.New()

	h := &FieldHandler{
		Resource: nil,
	}

	app.Get("/lenses", appContext.Wrap(func(c *appContext.Context) error {
		return HandleLensIndex(h, c)
	}))

	req := httptest.NewRequest("GET", "/lenses", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	if resp.StatusCode != 404 {
		t.Errorf("Expected status 404, got %d", resp.StatusCode)
	}
}

func TestHandleLens_Success(t *testing.T) {
	app := fiber.New()

	// Filtered users for lens
	users := []interface{}{
		User{ID: 1, FullName: "Active User", Email: "active@example.com"},
	}

	mockProvider := &MockDataProvider{
		Items: users,
		Total: 1,
	}

	fieldDefs := []fields.Element{
		fields.ID(),
		fields.Text("Full Name", "full_name"),
	}

	h := NewFieldHandler(mockProvider)
	h.Resource = &MockResource{}
	h.Elements = fieldDefs
	h.Title = "Active Users Lens"

	app.Get("/lens", FieldContextMiddleware(nil, nil, core.ContextIndex, fieldDefs), appContext.Wrap(func(c *appContext.Context) error {
		return HandleLens(h, c)
	}))

	req := httptest.NewRequest("GET", "/lens?page=1&per_page=10", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	dataList := response["data"].([]interface{})
	if len(dataList) != 1 {
		t.Errorf("Expected 1 item, got %d", len(dataList))
	}
}

func TestHandleLens_WithFiltering(t *testing.T) {
	app := fiber.New()

	// Lens with specific filter
	users := []interface{}{
		User{ID: 2, FullName: "Premium User", Email: "premium@example.com"},
	}

	mockProvider := &MockDataProvider{
		Items: users,
		Total: 1,
	}

	fieldDefs := []fields.Element{
		fields.ID(),
		fields.Text("Full Name", "full_name"),
	}

	h := NewFieldHandler(mockProvider)
	h.Resource = &MockResource{}
	h.Elements = fieldDefs
	h.Title = "Premium Users Lens"

	app.Get("/lens", FieldContextMiddleware(nil, nil, core.ContextIndex, fieldDefs), appContext.Wrap(func(c *appContext.Context) error {
		return HandleLens(h, c)
	}))

	req := httptest.NewRequest("GET", "/lens", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	dataList := response["data"].([]interface{})
	if len(dataList) != 1 {
		t.Errorf("Expected 1 filtered item, got %d", len(dataList))
	}
}

func TestHandleLens_EmptyResults(t *testing.T) {
	app := fiber.New()

	mockProvider := &MockDataProvider{
		Items: []interface{}{},
		Total: 0,
	}

	fieldDefs := []fields.Element{
		fields.ID(),
		fields.Text("Full Name", "full_name"),
	}

	h := NewFieldHandler(mockProvider)
	h.Resource = &MockResource{}
	h.Elements = fieldDefs
	h.Title = "Empty Lens"

	app.Get("/lens", FieldContextMiddleware(nil, nil, core.ContextIndex, fieldDefs), appContext.Wrap(func(c *appContext.Context) error {
		return HandleLens(h, c)
	}))

	req := httptest.NewRequest("GET", "/lens", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	dataList := response["data"].([]interface{})
	if len(dataList) != 0 {
		t.Errorf("Expected 0 items, got %d", len(dataList))
	}
}
