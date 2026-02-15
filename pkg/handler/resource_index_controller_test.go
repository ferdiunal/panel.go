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
	"github.com/gofiber/fiber/v2"
)

func TestHandleResourceIndex_Success(t *testing.T) {
	app := fiber.New()

	users := []interface{}{
		User{ID: 1, FullName: "John Doe", Email: "john@example.com"},
		User{ID: 2, FullName: "Jane Doe", Email: "jane@example.com"},
	}

	mockProvider := &MockDataProvider{
		Items: users,
		Total: 2,
	}

	fieldDefs := []fields.Element{
		fields.ID(),
		fields.Text("Full Name", "full_name"),
	}

	h := NewFieldHandler(mockProvider)
	h.Resource = &MockResource{}
	h.Elements = fieldDefs

	app.Get("/users", FieldContextMiddleware(nil, nil, core.ContextIndex, fieldDefs), appContext.Wrap(func(c *appContext.Context) error {
		return HandleResourceIndex(h, c)
	}))

	req := httptest.NewRequest("GET", "/users?page=1&per_page=10", nil)
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
		t.Errorf("Expected 2 items, got %d", len(dataList))
	}

	meta := response["meta"].(map[string]interface{})
	if meta["row_click_action"] != string(resource.IndexRowClickActionEdit) {
		t.Errorf("Expected default row_click_action=edit, got %v", meta["row_click_action"])
	}

	paginationMeta := meta["pagination"].(map[string]interface{})
	if paginationMeta["type"] != string(resource.IndexPaginationTypeLinks) {
		t.Errorf("Expected default pagination.type=links, got %v", paginationMeta["type"])
	}

	reorderMeta := meta["reorder"].(map[string]interface{})
	if reorderMeta["enabled"] != false {
		t.Errorf("Expected reorder.enabled=false, got %v", reorderMeta["enabled"])
	}
}

func TestHandleResourceIndex_CustomPaginationType(t *testing.T) {
	app := fiber.New()

	mockProvider := &MockDataProvider{
		Items: []interface{}{
			User{ID: 1, FullName: "John Doe", Email: "john@example.com"},
		},
		Total: 1,
	}

	fieldDefs := []fields.Element{
		fields.ID(),
		fields.Text("Full Name", "full_name"),
	}

	h := NewFieldHandler(mockProvider)
	h.Resource = &MockResource{}
	h.Elements = fieldDefs
	h.IndexPaginationType = resource.IndexPaginationTypeLoadMore

	app.Get("/users", FieldContextMiddleware(nil, nil, core.ContextIndex, fieldDefs), appContext.Wrap(func(c *appContext.Context) error {
		return HandleResourceIndex(h, c)
	}))

	req := httptest.NewRequest("GET", "/users", nil)
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

	meta := response["meta"].(map[string]interface{})
	paginationMeta := meta["pagination"].(map[string]interface{})
	if paginationMeta["type"] != string(resource.IndexPaginationTypeLoadMore) {
		t.Errorf("Expected pagination.type=load_more, got %v", paginationMeta["type"])
	}
}

func TestHandleResourceIndex_Unauthorized(t *testing.T) {
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
	h.Policy = &MockPolicy{AllowViewAny: false}

	app.Get("/users", FieldContextMiddleware(nil, nil, core.ContextIndex, fieldDefs), appContext.Wrap(func(c *appContext.Context) error {
		return HandleResourceIndex(h, c)
	}))

	req := httptest.NewRequest("GET", "/users", nil)
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

func TestHandleResourceIndex_Pagination(t *testing.T) {
	app := fiber.New()

	users := []interface{}{
		User{ID: 1, FullName: "User 1", Email: "user1@example.com"},
		User{ID: 2, FullName: "User 2", Email: "user2@example.com"},
		User{ID: 3, FullName: "User 3", Email: "user3@example.com"},
	}

	mockProvider := &MockDataProvider{
		Items: users,
		Total: 3,
	}

	fieldDefs := []fields.Element{
		fields.ID(),
		fields.Text("Full Name", "full_name"),
	}

	h := NewFieldHandler(mockProvider)
	h.Resource = &MockResource{}
	h.Elements = fieldDefs

	app.Get("/users", FieldContextMiddleware(nil, nil, core.ContextIndex, fieldDefs), appContext.Wrap(func(c *appContext.Context) error {
		return HandleResourceIndex(h, c)
	}))

	req := httptest.NewRequest("GET", "/users?page=2&per_page=2", nil)
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

	meta := response["meta"].(map[string]interface{})
	if meta["current_page"] != float64(2) {
		t.Errorf("Expected page 2, got %v", meta["current_page"])
	}
	if meta["per_page"] != float64(2) {
		t.Errorf("Expected per_page 2, got %v", meta["per_page"])
	}
}

func TestHandleResourceIndex_Sorting(t *testing.T) {
	app := fiber.New()

	users := []interface{}{
		User{ID: 1, FullName: "Alice", Email: "alice@example.com"},
		User{ID: 2, FullName: "Bob", Email: "bob@example.com"},
	}

	mockProvider := &MockDataProvider{
		Items: users,
		Total: 2,
	}

	fieldDefs := []fields.Element{
		fields.ID(),
		fields.Text("Full Name", "full_name"),
	}

	h := NewFieldHandler(mockProvider)
	h.Resource = &MockResource{}
	h.Elements = fieldDefs

	app.Get("/users", FieldContextMiddleware(nil, nil, core.ContextIndex, fieldDefs), appContext.Wrap(func(c *appContext.Context) error {
		return HandleResourceIndex(h, c)
	}))

	req := httptest.NewRequest("GET", "/users?sort_column=full_name&sort_direction=asc", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestHandleResourceIndex_Search(t *testing.T) {
	app := fiber.New()

	users := []interface{}{
		User{ID: 1, FullName: "John Doe", Email: "john@example.com"},
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

	app.Get("/users", FieldContextMiddleware(nil, nil, core.ContextIndex, fieldDefs), appContext.Wrap(func(c *appContext.Context) error {
		return HandleResourceIndex(h, c)
	}))

	req := httptest.NewRequest("GET", "/users?search=John", nil)
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

// MockPolicy for testing authorization
type MockPolicy struct {
	AllowViewAny bool
	AllowView    bool
	AllowCreate  bool
	AllowUpdate  bool
	AllowDelete  bool
}

func (m *MockPolicy) ViewAny(c *appContext.Context) bool {
	return m.AllowViewAny
}

func (m *MockPolicy) View(c *appContext.Context, model interface{}) bool {
	return m.AllowView
}

func (m *MockPolicy) Create(c *appContext.Context) bool {
	return m.AllowCreate
}

func (m *MockPolicy) Update(c *appContext.Context, model interface{}) bool {
	return m.AllowUpdate
}

func (m *MockPolicy) Delete(c *appContext.Context, model interface{}) bool {
	return m.AllowDelete
}
