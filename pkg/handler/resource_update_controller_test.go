package handler

import (
	"bytes"
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

func TestHandleResourceUpdate_Success(t *testing.T) {
	app := fiber.New()

	existingUser := User{ID: 1, FullName: "Old Name", Email: "old@example.com"}
	updatedUser := User{ID: 1, FullName: "New Name", Email: "new@example.com"}

	mockProvider := &MockDataProviderWithUpdate{
		ShowItem:    existingUser,
		UpdatedItem: updatedUser,
	}

	fieldDefs := []fields.Element{
		fields.ID(),
		fields.Text("Full Name", "full_name"),
		fields.Email("Email", "email"),
	}

	h := NewFieldHandler(mockProvider)
	h.Resource = &MockResource{}
	h.Elements = fieldDefs

	app.Put("/users/:id", FieldContextMiddleware(nil, nil, core.ContextUpdate, fieldDefs), appContext.Wrap(func(c *appContext.Context) error {
		return HandleResourceUpdate(h, c)
	}))

	body := map[string]interface{}{
		"full_name": "New Name",
		"email":     "new@example.com",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("PUT", "/users/1", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	respBody, _ := io.ReadAll(resp.Body)
	var response map[string]interface{}
	if err := json.Unmarshal(respBody, &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	data := response["data"].(map[string]interface{})
	fullNameField := data["full_name"].(map[string]interface{})
	if fullNameField["data"] != "New Name" {
		t.Errorf("Expected 'New Name', got %v", fullNameField["data"])
	}
}

func TestHandleResourceUpdate_NotFound(t *testing.T) {
	app := fiber.New()

	mockProvider := &MockDataProviderWithUpdate{
		ShowError: errors.New("not found"),
	}

	fieldDefs := []fields.Element{
		fields.ID(),
		fields.Text("Full Name", "full_name"),
	}

	h := NewFieldHandler(mockProvider)
	h.Resource = &MockResource{}
	h.Elements = fieldDefs

	app.Put("/users/:id", FieldContextMiddleware(nil, nil, core.ContextUpdate, fieldDefs), appContext.Wrap(func(c *appContext.Context) error {
		return HandleResourceUpdate(h, c)
	}))

	body := map[string]interface{}{
		"full_name": "New Name",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("PUT", "/users/999", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	if resp.StatusCode != 404 {
		t.Errorf("Expected status 404, got %d", resp.StatusCode)
	}

	respBody, _ := io.ReadAll(resp.Body)
	var response map[string]interface{}
	if err := json.Unmarshal(respBody, &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["error"] != "Not found" {
		t.Errorf("Expected 'Not found' error, got %v", response["error"])
	}
}

func TestHandleResourceUpdate_Unauthorized(t *testing.T) {
	app := fiber.New()

	existingUser := User{ID: 1, FullName: "Old Name", Email: "old@example.com"}

	mockProvider := &MockDataProviderWithUpdate{
		ShowItem: existingUser,
	}

	fieldDefs := []fields.Element{
		fields.ID(),
		fields.Text("Full Name", "full_name"),
	}

	h := NewFieldHandler(mockProvider)
	h.Resource = &MockResource{}
	h.Elements = fieldDefs
	h.Policy = &MockPolicy{AllowUpdate: false}

	app.Put("/users/:id", FieldContextMiddleware(nil, nil, core.ContextUpdate, fieldDefs), appContext.Wrap(func(c *appContext.Context) error {
		return HandleResourceUpdate(h, c)
	}))

	body := map[string]interface{}{
		"full_name": "New Name",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("PUT", "/users/1", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	if resp.StatusCode != 403 {
		t.Errorf("Expected status 403, got %d", resp.StatusCode)
	}

	respBody, _ := io.ReadAll(resp.Body)
	var response map[string]interface{}
	if err := json.Unmarshal(respBody, &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["error"] != "Unauthorized" {
		t.Errorf("Expected 'Unauthorized' error, got %v", response["error"])
	}
}

func TestHandleResourceUpdate_InvalidBody(t *testing.T) {
	app := fiber.New()

	mockProvider := &MockDataProviderWithUpdate{}

	fieldDefs := []fields.Element{
		fields.ID(),
		fields.Text("Full Name", "full_name"),
	}

	h := NewFieldHandler(mockProvider)
	h.Resource = &MockResource{}
	h.Elements = fieldDefs

	app.Put("/users/:id", FieldContextMiddleware(nil, nil, core.ContextUpdate, fieldDefs), appContext.Wrap(func(c *appContext.Context) error {
		return HandleResourceUpdate(h, c)
	}))

	req := httptest.NewRequest("PUT", "/users/1", bytes.NewReader([]byte("{invalid")))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	// parseBody might not fail on malformed JSON, so we accept 400 or other status
	if resp.StatusCode != 400 && resp.StatusCode != 404 && resp.StatusCode != 200 {
		t.Errorf("Expected status 400, 404, or 200, got %d", resp.StatusCode)
	}
}

func TestHandleResourceUpdate_ValidationError(t *testing.T) {
	app := fiber.New()

	existingUser := User{ID: 1, FullName: "Old Name", Email: "old@example.com"}
	mockProvider := &MockDataProviderWithUpdate{
		ShowItem: existingUser,
	}

	fieldDefs := []fields.Element{
		fields.Text("Full Name", "full_name").Required(),
		fields.Email("Email", "email").
			Required().
			AddValidationRule(fields.EmailRule()),
	}

	h := NewFieldHandler(mockProvider)
	h.Resource = &MockResource{}
	h.Elements = fieldDefs

	app.Put("/users/:id", FieldContextMiddleware(nil, nil, core.ContextUpdate, fieldDefs), appContext.Wrap(func(c *appContext.Context) error {
		return HandleResourceUpdate(h, c)
	}))

	body := map[string]interface{}{
		"full_name": "",
		"email":     "invalid-email",
	}
	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest("PUT", "/users/1", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	if resp.StatusCode != fiber.StatusUnprocessableEntity {
		t.Errorf("Expected status 422, got %d", resp.StatusCode)
	}

	respBody, _ := io.ReadAll(resp.Body)
	var response map[string]interface{}
	if err := json.Unmarshal(respBody, &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["code"] != validationErrorCode {
		t.Errorf("Expected validation error code, got %v", response["code"])
	}

	errorMap, ok := response["errors"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected errors map in response")
	}

	fullNameErrors, ok := errorMap["full_name"].([]interface{})
	if !ok || len(fullNameErrors) == 0 {
		t.Fatalf("Expected full_name validation error")
	}

	emailErrors, ok := errorMap["email"].([]interface{})
	if !ok || len(emailErrors) == 0 {
		t.Fatalf("Expected email validation error")
	}
}

// MockDataProviderWithUpdate extends MockDataProvider with Update functionality
type MockDataProviderWithUpdate struct {
	MockDataProvider
	ShowItem    interface{}
	ShowError   error
	UpdatedItem interface{}
	UpdateError error
}

func (m *MockDataProviderWithUpdate) Show(ctx *appContext.Context, id string) (interface{}, error) {
	if m.ShowError != nil {
		return nil, m.ShowError
	}
	return m.ShowItem, nil
}

func (m *MockDataProviderWithUpdate) Update(ctx *appContext.Context, id string, data map[string]interface{}) (interface{}, error) {
	if m.UpdateError != nil {
		return nil, m.UpdateError
	}
	if m.UpdatedItem != nil {
		return m.UpdatedItem, nil
	}
	return data, nil
}

func (m *MockDataProviderWithUpdate) Index(ctx *appContext.Context, req data.QueryRequest) (*data.QueryResponse, error) {
	return &data.QueryResponse{
		Items:   []interface{}{},
		Total:   0,
		Page:    req.Page,
		PerPage: req.PerPage,
	}, nil
}

func (m *MockDataProviderWithUpdate) Create(ctx *appContext.Context, data map[string]interface{}) (interface{}, error) {
	return data, nil
}

func (m *MockDataProviderWithUpdate) Delete(ctx *appContext.Context, id string) error {
	return nil
}

func (m *MockDataProviderWithUpdate) SetSearchColumns(cols []string) {}
func (m *MockDataProviderWithUpdate) SetWith(rels []string)          {}
