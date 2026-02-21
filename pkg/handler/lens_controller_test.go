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

type mockLensResourceGridDisabled struct {
	MockResource
}

func (m *mockLensResourceGridDisabled) IsGridEnabled() bool {
	return false
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
	if response["record_title_key"] != "full_name" {
		t.Errorf("Expected record_title_key=full_name, got %v", response["record_title_key"])
	}
	if response["grid_enabled"] != true {
		t.Errorf("Expected grid_enabled=true, got %v", response["grid_enabled"])
	}

	dataList := response["data"].([]interface{})
	if len(dataList) != 1 {
		t.Errorf("Expected 1 item, got %d", len(dataList))
	}
}

func TestHandleLens_GridViewVisibility(t *testing.T) {
	app := fiber.New()

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
		fields.Text("Email", "email").HideOnList().ShowOnGrid(),
		fields.Text("Avatar", "avatar").HideOnGrid(),
	}

	h := NewFieldHandler(mockProvider)
	h.Resource = &MockResource{}
	h.Elements = fieldDefs

	app.Get("/lens", FieldContextMiddleware(nil, nil, core.ContextIndex, fieldDefs), appContext.Wrap(func(c *appContext.Context) error {
		return HandleLens(h, c)
	}))

	req := httptest.NewRequest("GET", "/lens?view=grid", nil)
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

	headers := response["headers"].([]interface{})
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

	dataList := response["data"].([]interface{})
	if len(dataList) != 1 {
		t.Fatalf("expected 1 data item, got %d", len(dataList))
	}
	row := dataList[0].(map[string]interface{})
	if _, ok := row["avatar"]; !ok {
		t.Fatal("expected avatar field in grid row payload (HideOnGrid should hide card listing only)")
	}
}

func TestHandleLens_IncludesStackFieldWithChildrenOnListAndGrid(t *testing.T) {
	app := fiber.New()

	users := []interface{}{
		User{ID: 1, FullName: "Active User", Email: "active@example.com"},
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

	app.Get("/lens", FieldContextMiddleware(nil, nil, core.ContextIndex, fieldDefs), appContext.Wrap(func(c *appContext.Context) error {
		return HandleLens(h, c)
	}))

	for _, rawURL := range []string{
		"/lens?page=1&per_page=10",
		"/lens?view=grid&page=1&per_page=10",
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

			headers := response["headers"].([]interface{})
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

			if got := childByKey["full_name"]["data"]; got != "Active User" {
				t.Fatalf("expected stack full_name child data to be Active User for %s, got %v", rawURL, got)
			}
			if got := childByKey["email"]["data"]; got != "active@example.com" {
				t.Fatalf("expected stack email child data to be active@example.com for %s, got %v", rawURL, got)
			}

			if _, ok := row["full_name"]; !ok {
				t.Fatalf("expected regular fields to remain in row payload for %s", rawURL)
			}
		}
	}

func TestHandleLens_GridViewVisibilityWithoutContextMiddleware(t *testing.T) {
	app := fiber.New()

	users := []interface{}{
		User{ID: 1, FullName: "Active User", Email: "active@example.com"},
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
	app.Get("/lens", appContext.Wrap(func(c *appContext.Context) error {
		return HandleLens(h, c)
	}))

	req := httptest.NewRequest("GET", "/lens?view=grid", nil)
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

	headers := response["headers"].([]interface{})
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

func TestHandleLens_GridViewDisabledByResource(t *testing.T) {
	app := fiber.New()

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
		fields.Text("Email", "email").HideOnList().ShowOnGrid(),
	}

	h := NewFieldHandler(mockProvider)
	h.Resource = &mockLensResourceGridDisabled{}
	h.Elements = fieldDefs

	app.Get("/lens", FieldContextMiddleware(nil, nil, core.ContextIndex, fieldDefs), appContext.Wrap(func(c *appContext.Context) error {
		return HandleLens(h, c)
	}))

	req := httptest.NewRequest("GET", "/lens?view=grid", nil)
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

	if response["grid_enabled"] != false {
		t.Fatalf("expected grid_enabled=false, got %v", response["grid_enabled"])
	}

	headers := response["headers"].([]interface{})
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
