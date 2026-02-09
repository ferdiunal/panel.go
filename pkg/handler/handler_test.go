package handler

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ferdiunal/panel.go/pkg/auth"
	appContext "github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/core"
	"github.com/ferdiunal/panel.go/pkg/data"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/resource"
	"github.com/ferdiunal/panel.go/pkg/widget"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	FullName  string    `json:"full_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MockDataProvider struct {
	Items []interface{}
	Total int64
}

func (m *MockDataProvider) Index(ctx *appContext.Context, req data.QueryRequest) (*data.QueryResponse, error) {
	return &data.QueryResponse{
		Items:   m.Items,
		Total:   m.Total,
		Page:    req.Page,
		PerPage: req.PerPage,
	}, nil
}

func (m *MockDataProvider) Show(ctx *appContext.Context, id string) (interface{}, error) {
	return nil, nil
}

func (m *MockDataProvider) Create(ctx *appContext.Context, data map[string]interface{}) (interface{}, error) {
	return data, nil
}

func (m *MockDataProvider) Update(ctx *appContext.Context, id string, data map[string]interface{}) (interface{}, error) {
	return data, nil
}

func (m *MockDataProvider) Delete(ctx *appContext.Context, id string) error {
	return nil
}

func (m *MockDataProvider) SetSearchColumns(cols []string)                    {}
func (m *MockDataProvider) SetWith(rels []string)                             {}
func (m *MockDataProvider) SetRelationshipFields(fields []fields.RelationshipField) {}

func TestFieldHandler_List(t *testing.T) {
	app := fiber.New()

	user := User{
		ID:       1,
		FullName: "John Doe",
		Email:    "john@example.com",
	}

	fieldDefs := []fields.Element{
		fields.ID(),
		fields.Text("Full Name", "full_name"),
		fields.Email("Email Address", "email"),
	}

	h := NewFieldHandler(&MockDataProvider{})

	app.Get("/fields", FieldContextMiddleware(&user, nil, core.ContextIndex, fieldDefs), appContext.Wrap(h.List))

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

	// Verify ID
	if response[0]["key"] != "id" || response[0]["data"] != float64(1) { // JSON numbers are floats
		t.Errorf("ID field incorrect: %v", response[0])
	}
}

func TestFieldHandler_Index(t *testing.T) {
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

	app.Get("/users", FieldContextMiddleware(nil, nil, core.ContextIndex, fieldDefs), appContext.Wrap(h.Index))

	req := httptest.NewRequest("GET", "/users?page=1&per_page=10", nil)
	resp, _ := app.Test(req)

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

	item1 := dataList[0].(map[string]interface{})
	idField := item1["id"].(map[string]interface{})
	if idField["data"] != float64(1) {
		t.Errorf("Item 1 ID mismatch: %v", idField["data"])
	}
}

func TestFieldHandler_MapSupport(t *testing.T) {
	app := fiber.New()

	data := map[string]interface{}{
		"id":    100,
		"title": "Hello World",
	}

	fieldDefs := []fields.Element{
		fields.ID(),
		fields.Text("Title"),
	}

	h := NewFieldHandler(&MockDataProvider{})

	app.Get("/fields", FieldContextMiddleware(data, nil, core.ContextIndex, fieldDefs), appContext.Wrap(h.List))

	req := httptest.NewRequest("GET", "/fields", nil)
	resp, _ := app.Test(req)

	body, _ := io.ReadAll(resp.Body)
	var response []map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(response) != 2 {
		t.Errorf("Expected 2 fields, got %d", len(response))
	}

	if response[1]["data"] != "Hello World" {
		t.Errorf("Map value resolution failed: %v", response[1])
	}
}

// Mock Resource Implementation
type MockResource struct{}

func (m *MockResource) Model() interface{} {
	return &User{}
}

func (m *MockResource) Fields() []fields.Element {
	return []fields.Element{
		fields.ID(),
		fields.Text("Full Name", "full_name"),
	}
}

func (m *MockResource) With() []string {
	return nil
}

func (m *MockResource) Lenses() []resource.Lens {
	return nil
}

func (m *MockResource) Title() string                                         { return "Mock User" }
func (m *MockResource) Icon() string                                          { return "user" }
func (m *MockResource) Group() string                                         { return "Management" }
func (m *MockResource) Policy() auth.Policy                                   { return nil }
func (m *MockResource) GetSortable() []resource.Sortable                      { return nil }
func (m *MockResource) Slug() string                                          { return "users" }
func (m *MockResource) GetDialogType() resource.DialogType                    { return resource.DialogTypeSheet }
func (m *MockResource) SetDialogType(t resource.DialogType) resource.Resource { return m }
func (m *MockResource) Repository(db *gorm.DB) data.DataProvider              { return nil }
func (m *MockResource) GetRecordTitleKey() string                             { return "full_name" }
func (m *MockResource) SetRecordTitleKey(key string) resource.Resource        { return m }
func (m *MockResource) OpenAPIEnabled() bool                                  { return false }

func (m *MockResource) Cards() []widget.Card {
	return []widget.Card{
		widget.NewCard("Test Card", "test-card"),
	}
}

func (m *MockResource) NavigationOrder() int { return 1 }
func (m *MockResource) Visible() bool        { return true }
func (m *MockResource) StoreHandler(c *appContext.Context, file *multipart.FileHeader, storagePath string, storageURL string) (string, error) {
	return "", nil
}

func (m *MockResource) GetFields(ctx *appContext.Context) []fields.Element {
	return m.Fields()
}

func (m *MockResource) GetCards(ctx *appContext.Context) []widget.Card {
	return m.Cards()
}

func (m *MockResource) GetLenses() []resource.Lens {
	return m.Lenses()
}

func (m *MockResource) GetPolicy() auth.Policy {
	return m.Policy()
}

func (m *MockResource) ResolveField(fieldName string, item interface{}) (interface{}, error) {
	return nil, nil
}

func (m *MockResource) GetActions() []resource.Action {
	return []resource.Action{}
}

func (m *MockResource) GetFilters() []resource.Filter {
	return []resource.Filter{}
}

func (m *MockResource) RecordTitle(record interface{}) string {
	if user, ok := record.(*User); ok {
		return user.FullName
	}
	if user, ok := record.(User); ok {
		return user.FullName
	}
	return ""
}

func TestResourceHandler(t *testing.T) {
	// Setup DB for NewResourceHandler (requires GORM DB)
	// We use the same setup as GormDataProvider test but mocked differently or just basic sqlite
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to db: %v", err)
	}
	db.AutoMigrate(&User{})
	db.Create(&User{ID: 1, FullName: "Resource User", Email: "res@example.com"})

	// Create Handler using Resource abstraction
	h := NewResourceHandler(db, &MockResource{}, "", "")

	app := fiber.New()
	app.Get("/resource-users", appContext.Wrap(h.Index)) // No middleware needed!

	req := httptest.NewRequest("GET", "/resource-users", nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	dataList := response["data"].([]interface{})
	if len(dataList) != 1 {
		t.Errorf("Expected 1 item, got %d", len(dataList))
	}

	item1 := dataList[0].(map[string]interface{})
	nameField := item1["full_name"].(map[string]interface{})
	if nameField["data"] != "Resource User" {
		t.Errorf("Mismatch in data: %v", nameField["data"])
	}
}
