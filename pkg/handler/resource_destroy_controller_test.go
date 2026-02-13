package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	appContext "github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/core"
	"github.com/ferdiunal/panel.go/pkg/data"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/gofiber/fiber/v2"
)

func TestHandleResourceDestroy_Success(t *testing.T) {
	app := fiber.New()

	existingUser := User{ID: 1, FullName: "John Doe", Email: "john@example.com"}

	mockProvider := &MockDataProviderWithDelete{
		ShowItem: existingUser,
	}

	fieldDefs := []fields.Element{
		fields.ID(),
		fields.Text("Full Name", "full_name"),
	}

	h := NewFieldHandler(mockProvider)
	h.Resource = &MockResource{}
	h.Elements = fieldDefs

	app.Delete("/users/:id", FieldContextMiddleware(nil, nil, core.ContextDetail, fieldDefs), appContext.Wrap(func(c *appContext.Context) error {
		return HandleResourceDestroy(h, c)
	}))

	req := httptest.NewRequest("DELETE", "/users/1", nil)
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

	if response["message"] != "Deleted successfully" {
		t.Errorf("Expected 'Deleted successfully' message, got %v", response["message"])
	}
}

func TestHandleResourceDestroy_NotFound(t *testing.T) {
	app := fiber.New()

	mockProvider := &MockDataProviderWithDelete{
		ShowError: errors.New("not found"),
	}

	fieldDefs := []fields.Element{
		fields.ID(),
		fields.Text("Full Name", "full_name"),
	}

	h := NewFieldHandler(mockProvider)
	h.Resource = &MockResource{}
	h.Elements = fieldDefs

	app.Delete("/users/:id", FieldContextMiddleware(nil, nil, core.ContextDetail, fieldDefs), appContext.Wrap(func(c *appContext.Context) error {
		return HandleResourceDestroy(h, c)
	}))

	req := httptest.NewRequest("DELETE", "/users/999", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	if resp.StatusCode != 404 {
		t.Errorf("Expected status 404, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["error"] != "Not found" {
		t.Errorf("Expected 'Not found' error, got %v", response["error"])
	}
}

func TestHandleResourceDestroy_Unauthorized(t *testing.T) {
	app := fiber.New()

	existingUser := User{ID: 1, FullName: "John Doe", Email: "john@example.com"}

	mockProvider := &MockDataProviderWithDelete{
		ShowItem: existingUser,
	}

	fieldDefs := []fields.Element{
		fields.ID(),
		fields.Text("Full Name", "full_name"),
	}

	h := NewFieldHandler(mockProvider)
	h.Resource = &MockResource{}
	h.Elements = fieldDefs
	h.Policy = &MockPolicy{AllowDelete: false}

	app.Delete("/users/:id", FieldContextMiddleware(nil, nil, core.ContextDetail, fieldDefs), appContext.Wrap(func(c *appContext.Context) error {
		return HandleResourceDestroy(h, c)
	}))

	req := httptest.NewRequest("DELETE", "/users/1", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	if resp.StatusCode != 403 {
		t.Errorf("Expected status 403, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["error"] != "Unauthorized" {
		t.Errorf("Expected 'Unauthorized' error, got %v", response["error"])
	}
}

// MockDataProviderWithDelete extends MockDataProvider with Delete functionality
type MockDataProviderWithDelete struct {
	MockDataProvider
	ShowItem    interface{}
	ShowError   error
	DeleteError error
}

func (m *MockDataProviderWithDelete) Show(ctx *appContext.Context, id string) (interface{}, error) {
	if m.ShowError != nil {
		return nil, m.ShowError
	}
	return m.ShowItem, nil
}

func (m *MockDataProviderWithDelete) Delete(ctx *appContext.Context, id string) error {
	return m.DeleteError
}

func (m *MockDataProviderWithDelete) Index(ctx *appContext.Context, req data.QueryRequest) (*data.QueryResponse, error) {
	return &data.QueryResponse{
		Items:   []interface{}{},
		Total:   0,
		Page:    req.Page,
		PerPage: req.PerPage,
	}, nil
}

func (m *MockDataProviderWithDelete) Create(ctx *appContext.Context, data map[string]interface{}) (interface{}, error) {
	return data, nil
}

func (m *MockDataProviderWithDelete) Update(ctx *appContext.Context, id string, data map[string]interface{}) (interface{}, error) {
	return data, nil
}

func (m *MockDataProviderWithDelete) SetSearchColumns(cols []string) {}
func (m *MockDataProviderWithDelete) SetWith(rels []string)          {}
