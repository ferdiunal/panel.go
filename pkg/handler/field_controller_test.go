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

func TestHandleFieldList_Success(t *testing.T) {
	app := fiber.New()

	fieldDefs := []fields.Element{
		fields.ID(),
		fields.Text("Full Name", "full_name"),
		fields.Text("Email", "email"),
	}

	h := &FieldHandler{
		Elements: fieldDefs,
	}

	app.Get("/fields", FieldContextMiddleware(nil, nil, core.ContextIndex, fieldDefs), appContext.Wrap(func(c *appContext.Context) error {
		return HandleFieldList(h, c)
	}))

	req := httptest.NewRequest("GET", "/fields", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var response []map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(response) != 3 {
		t.Errorf("Expected 3 fields, got %d", len(response))
	}
}

func TestHandleFieldList_WithResource(t *testing.T) {
	app := fiber.New()

	user := User{ID: 1, FullName: "John Doe", Email: "john@example.com"}

	fieldDefs := []fields.Element{
		fields.ID(),
		fields.Text("Full Name", "full_name"),
	}

	h := &FieldHandler{
		Elements: fieldDefs,
	}

	app.Get("/fields", FieldContextMiddleware(user, nil, core.ContextIndex, fieldDefs), appContext.Wrap(func(c *appContext.Context) error {
		return HandleFieldList(h, c)
	}))

	req := httptest.NewRequest("GET", "/fields", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var response []map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(response) != 2 {
		t.Errorf("Expected 2 fields, got %d", len(response))
	}
}

func TestHandleFieldList_NoContext(t *testing.T) {
	app := fiber.New()

	h := &FieldHandler{}

	// Don't use middleware, so context won't be set
	app.Get("/fields", appContext.Wrap(func(c *appContext.Context) error {
		return HandleFieldList(h, c)
	}))

	req := httptest.NewRequest("GET", "/fields", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	if resp.StatusCode != 500 {
		t.Errorf("Expected status 500, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["error"] != "Field context not found" {
		t.Errorf("Expected 'Field context not found' error, got %v", response["error"])
	}
}
