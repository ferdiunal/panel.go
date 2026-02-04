package handler

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	appContext "github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/core"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/gofiber/fiber/v2"
)

func TestHandleResourceCreate_Success(t *testing.T) {
	app := fiber.New()

	mockProvider := &MockDataProvider{}

	fieldDefs := []fields.Element{
		fields.ID(),
		fields.Text("Full Name", "full_name"),
		fields.Email("Email", "email"),
	}

	h := NewFieldHandler(mockProvider)
	h.Resource = &MockResource{}
	h.Elements = fieldDefs

	app.Get("/users/create", FieldContextMiddleware(nil, nil, core.ContextCreate, fieldDefs), appContext.Wrap(func(c *appContext.Context) error {
		return HandleResourceCreate(h, c)
	}))

	req := httptest.NewRequest("GET", "/users/create", nil)
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

	fieldsData := response["fields"].([]interface{})
	if len(fieldsData) == 0 {
		t.Error("Expected fields in response")
	}
}

func TestHandleResourceCreate_Unauthorized(t *testing.T) {
	app := fiber.New()

	mockProvider := &MockDataProvider{}

	fieldDefs := []fields.Element{
		fields.ID(),
		fields.Text("Full Name", "full_name"),
	}

	h := NewFieldHandler(mockProvider)
	h.Resource = &MockResource{}
	h.Elements = fieldDefs
	h.Policy = &MockPolicy{AllowCreate: false}

	app.Get("/users/create", FieldContextMiddleware(nil, nil, core.ContextCreate, fieldDefs), appContext.Wrap(func(c *appContext.Context) error {
		return HandleResourceCreate(h, c)
	}))

	req := httptest.NewRequest("GET", "/users/create", nil)
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
