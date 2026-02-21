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

type mockResourceGridDisabled struct {
	MockResource
}

func (m *mockResourceGridDisabled) IsGridEnabled() bool {
	return false
}

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
	if _, ok := meta["description"]; !ok {
		t.Errorf("Expected meta.description to be present")
	}
	if meta["record_title_key"] != "full_name" {
		t.Errorf("Expected record_title_key=full_name, got %v", meta["record_title_key"])
	}
	if meta["grid_enabled"] != true {
		t.Errorf("Expected grid_enabled=true, got %v", meta["grid_enabled"])
	}
	headers, ok := meta["headers"].([]interface{})
	if !ok {
		t.Fatalf("Expected meta.headers to be an array")
	}
	for _, headerRaw := range headers {
		header, ok := headerRaw.(map[string]interface{})
		if !ok {
			continue
		}
		if header["data"] != nil {
			t.Errorf("Expected header data to be nil for key %v, got %v", header["key"], header["data"])
		}
	}
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

func TestHandleResourceIndex_GridViewVisibility(t *testing.T) {
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
		fields.Text("Email", "email").HideOnList().ShowOnGrid(),
		fields.Text("Avatar", "avatar").HideOnGrid(),
		fields.Text("Grid Label", "grid_label").ShowOnlyGrid(),
	}

	h := NewFieldHandler(mockProvider)
	h.Resource = &MockResource{}
	h.Elements = fieldDefs

	app.Get("/users", FieldContextMiddleware(nil, nil, core.ContextIndex, fieldDefs), appContext.Wrap(func(c *appContext.Context) error {
		return HandleResourceIndex(h, c)
	}))

	req := httptest.NewRequest("GET", "/users?users[view]=grid", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	meta := response["meta"].(map[string]interface{})
	headers := meta["headers"].([]interface{})
	headerKeys := map[string]bool{}
	for _, raw := range headers {
		if item, ok := raw.(map[string]interface{}); ok {
			if key, ok := item["key"].(string); ok {
				headerKeys[key] = true
			}
		}
	}

	if !headerKeys["email"] {
		t.Fatal("expected email to be visible on grid when HideOnList+ShowOnGrid is set")
	}
	if headerKeys["avatar"] {
		t.Fatal("expected avatar to be hidden on grid when HideOnGrid is set")
	}
	if !headerKeys["grid_label"] {
		t.Fatal("expected grid_label to be visible on grid when ShowOnlyGrid is set")
	}

	dataList := response["data"].([]interface{})
	if len(dataList) != 1 {
		t.Fatalf("expected 1 data item, got %d", len(dataList))
	}
	row := dataList[0].(map[string]interface{})
	if _, ok := row["email"]; !ok {
		t.Fatal("expected email field in grid row payload")
	}
	if _, ok := row["avatar"]; !ok {
		t.Fatal("expected avatar field in grid row payload (HideOnGrid should hide card listing only)")
	}
}

func TestHandleResourceIndex_IncludesStackFieldWithChildrenOnListAndGrid(t *testing.T) {
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
		fields.Stack([]core.Element{
			fields.Text("Full Name", "full_name"),
			fields.Text("Email", "email"),
		}),
		fields.Text("Full Name", "full_name"),
	}

	h := NewFieldHandler(mockProvider)
	h.Resource = &MockResource{}
	h.Elements = fieldDefs

	app.Get("/users", FieldContextMiddleware(nil, nil, core.ContextIndex, fieldDefs), appContext.Wrap(func(c *appContext.Context) error {
		return HandleResourceIndex(h, c)
	}))

	for _, rawURL := range []string{
		"/users?page=1&per_page=10",
		"/users?users[view]=grid&page=1&per_page=10",
	} {
		req := httptest.NewRequest("GET", rawURL, nil)
		resp, err := app.Test(req)
		if err != nil {
			t.Fatalf("Failed to perform request %s: %v", rawURL, err)
		}
		if resp.StatusCode != 200 {
			t.Fatalf("Expected status 200 for %s, got %d", rawURL, resp.StatusCode)
		}

		body, _ := io.ReadAll(resp.Body)
		var response map[string]interface{}
		if err := json.Unmarshal(body, &response); err != nil {
			t.Fatalf("Failed to unmarshal response for %s: %v", rawURL, err)
		}

			meta := response["meta"].(map[string]interface{})
			headers := meta["headers"].([]interface{})
			stackHeaderFound := false
			for _, raw := range headers {
				header, ok := raw.(map[string]interface{})
				if !ok {
					continue
				}
				if key, _ := header["key"].(string); key == "stack" {
					stackHeaderFound = true
					break
				}
			}
			if !stackHeaderFound {
				t.Fatalf("expected stack header to be present for %s", rawURL)
			}

		dataList := response["data"].([]interface{})
		if len(dataList) != 1 {
			t.Fatalf("expected 1 data item for %s, got %d", rawURL, len(dataList))
		}

			row := dataList[0].(map[string]interface{})
			stackRaw, ok := row["stack"]
			if !ok {
				t.Fatalf("expected stack row field to be present for %s", rawURL)
			}
			stackField, ok := stackRaw.(map[string]interface{})
			if !ok {
				t.Fatalf("expected stack row field to be a map for %s, got %T", rawURL, stackRaw)
			}
			props, ok := stackField["props"].(map[string]interface{})
			if !ok {
				t.Fatalf("expected stack props to be map for %s", rawURL)
			}
			children, ok := props["fields"].([]interface{})
			if !ok || len(children) < 2 {
				t.Fatalf("expected stack children to be present for %s, got %T len=%d", rawURL, props["fields"], len(children))
			}

			childByKey := map[string]map[string]interface{}{}
			for _, rawChild := range children {
				child, ok := rawChild.(map[string]interface{})
				if !ok {
					continue
				}
				if key, _ := child["key"].(string); key != "" {
					childByKey[key] = child
				}
			}

			if got := childByKey["full_name"]["data"]; got != "John Doe" {
				t.Fatalf("expected stack full_name child data to be John Doe for %s, got %v", rawURL, got)
			}
			if got := childByKey["email"]["data"]; got != "john@example.com" {
				t.Fatalf("expected stack email child data to be john@example.com for %s, got %v", rawURL, got)
			}

			if _, ok := row["full_name"]; !ok {
				t.Fatalf("expected regular fields to remain in row payload for %s", rawURL)
			}
		}
	}

func TestHandleResourceIndex_GridViewVisibilityWithoutContextMiddleware(t *testing.T) {
	app := fiber.New()

	users := []interface{}{
		User{ID: 1, FullName: "John Doe", Email: "john@example.com"},
	}

	mockProvider := &MockDataProvider{
		Items: users,
		Total: 1,
	}

	fieldDefs := []fields.Element{
		fields.ID().HideOnGrid(),
		fields.Text("Full Name", "full_name"),
		fields.Text("Email", "email").HideOnList().ShowOnGrid(),
		fields.Text("Created At", "created_at").OnlyOnDetail(),
	}

	h := NewFieldHandler(mockProvider)
	h.Resource = &MockResource{}
	h.Elements = fieldDefs

	// Intentionally no FieldContextMiddleware. Handler should bootstrap ResourceContext.
	app.Get("/users", appContext.Wrap(func(c *appContext.Context) error {
		return HandleResourceIndex(h, c)
	}))

	req := httptest.NewRequest("GET", "/users?users[view]=grid", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	meta := response["meta"].(map[string]interface{})
	headers := meta["headers"].([]interface{})
	headerKeys := map[string]bool{}
	for _, raw := range headers {
		if item, ok := raw.(map[string]interface{}); ok {
			if key, ok := item["key"].(string); ok {
				headerKeys[key] = true
			}
		}
	}

	if headerKeys["id"] {
		t.Fatal("expected id to be hidden in grid context")
	}
	if headerKeys["created_at"] {
		t.Fatal("expected only_on_detail field to stay hidden in grid context")
	}
	if !headerKeys["email"] {
		t.Fatal("expected hide_on_list + show_on_grid field to be visible in grid context")
	}

	dataList := response["data"].([]interface{})
	if len(dataList) != 1 {
		t.Fatalf("expected 1 data item, got %d", len(dataList))
	}
	row := dataList[0].(map[string]interface{})
	if _, ok := row["id"]; !ok {
		t.Fatal("expected id field in grid row payload (HideOnGrid should hide card listing only)")
	}
}

func TestHandleResourceIndex_GridViewDisabledByResource(t *testing.T) {
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
		fields.Text("Email", "email").HideOnList().ShowOnGrid(),
	}

	h := NewFieldHandler(mockProvider)
	h.Resource = &mockResourceGridDisabled{}
	h.Elements = fieldDefs

	app.Get("/users", FieldContextMiddleware(nil, nil, core.ContextIndex, fieldDefs), appContext.Wrap(func(c *appContext.Context) error {
		return HandleResourceIndex(h, c)
	}))

	req := httptest.NewRequest("GET", "/users?users[view]=grid", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	meta := response["meta"].(map[string]interface{})
	if meta["grid_enabled"] != false {
		t.Fatalf("expected meta.grid_enabled=false, got %v", meta["grid_enabled"])
	}

	headers := meta["headers"].([]interface{})
	headerKeys := map[string]bool{}
	for _, raw := range headers {
		if item, ok := raw.(map[string]interface{}); ok {
			if key, ok := item["key"].(string); ok {
				headerKeys[key] = true
			}
		}
	}

	if headerKeys["email"] {
		t.Fatal("expected email to remain hidden when grid is disabled at resource level")
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
