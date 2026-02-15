package handler

import (
	"encoding/json"
	"fmt"
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

type HandlerTestCategory struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
}

type HandlerTestVariant struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	ProductID uint   `json:"product_id"`
	Name      string `json:"name"`
}

type HandlerTestProduct struct {
	ID              uint                 `gorm:"primaryKey" json:"id"`
	CategoryID      uint                 `json:"category_id"`
	Category        HandlerTestCategory  `json:"category" gorm:"foreignKey:CategoryID"`
	ProductVariants []HandlerTestVariant `json:"product_variants" gorm:"foreignKey:ProductID"`
}

type MockDataProvider struct {
	Items []interface{}
	Total int64
}

type TrackingDataProvider struct {
	MockDataProvider
	setWithCalled bool
	with          []string
	searchColumns []string
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

func (m *MockDataProvider) SetSearchColumns(cols []string)                          {}
func (m *MockDataProvider) SetWith(rels []string)                                   {}
func (m *MockDataProvider) SetRelationshipFields(fields []fields.RelationshipField) {}
func (m *TrackingDataProvider) SetSearchColumns(cols []string) {
	m.searchColumns = append([]string{}, cols...)
}
func (m *TrackingDataProvider) SetWith(rels []string) {
	m.setWithCalled = true
	m.with = append([]string{}, rels...)
}
func (m *MockDataProvider) QueryTable(ctx *appContext.Context, table string, conditions map[string]interface{}) ([]map[string]interface{}, error) {
	return []map[string]interface{}{}, nil
}
func (m *MockDataProvider) QueryRelationship(ctx *appContext.Context, relationshipType string, foreignKey string, foreignValue interface{}, displayField string) (interface{}, error) {
	return nil, nil
}
func (m *MockDataProvider) BeginTx(ctx *appContext.Context) (data.DataProvider, error) { return m, nil }
func (m *MockDataProvider) Commit() error                                              { return nil }
func (m *MockDataProvider) Rollback() error                                            { return nil }
func (m *MockDataProvider) Raw(ctx *appContext.Context, sql string, args ...interface{}) ([]map[string]interface{}, error) {
	return []map[string]interface{}{}, nil
}
func (m *MockDataProvider) Exec(ctx *appContext.Context, sql string, args ...interface{}) error {
	return nil
}
func (m *MockDataProvider) GetClient() interface{} { return nil }

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

func TestFieldHandler_Index_NormalizesRelationshipHeaderData(t *testing.T) {
	app := fiber.New()

	items := []interface{}{
		HandlerTestProduct{
			ID:              1,
			ProductVariants: nil,
		},
	}

	mockProvider := &MockDataProvider{
		Items: items,
		Total: 1,
	}

	relationshipField := fields.NewField("Variants", "product_variants")
	relationshipField.View = "has-many-field"
	relationshipField.Type = fields.TYPE_RELATIONSHIP

	fieldDefs := []fields.Element{
		relationshipField,
	}

	h := NewFieldHandler(mockProvider)
	h.Resource = &MockResource{}

	app.Get("/products", FieldContextMiddleware(nil, nil, core.ContextIndex, fieldDefs), appContext.Wrap(h.Index))

	req := httptest.NewRequest("GET", "/products?page=1&per_page=10", nil)
	resp, _ := app.Test(req)

	if resp.StatusCode != 200 {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	meta, ok := response["meta"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected meta to be a map, got %T", response["meta"])
	}

	headers, ok := meta["headers"].([]interface{})
	if !ok || len(headers) == 0 {
		t.Fatalf("expected headers to be a non-empty slice, got %T", meta["headers"])
	}

	header, ok := headers[0].(map[string]interface{})
	if !ok {
		t.Fatalf("expected first header to be map, got %T", headers[0])
	}

	dataValue, exists := header["data"]
	if !exists {
		t.Fatalf("expected header to contain data key")
	}

	dataSlice, ok := dataValue.([]interface{})
	if !ok {
		t.Fatalf("expected relationship header data to be []interface{}, got %T", dataValue)
	}
	if len(dataSlice) != 0 {
		t.Fatalf("expected relationship header data to be empty slice, got %v", dataSlice)
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

type MockResourceWithCustomRepository struct {
	*MockResource
	provider data.DataProvider
}

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
func (m *MockResource) TitleWithContext(ctx *fiber.Ctx) string                { return m.Title() }
func (m *MockResource) Icon() string                                          { return "user" }
func (m *MockResource) Group() string                                         { return "Management" }
func (m *MockResource) GroupWithContext(ctx *fiber.Ctx) string                { return m.Group() }
func (m *MockResource) Policy() auth.Policy                                   { return nil }
func (m *MockResource) GetSortable() []resource.Sortable                      { return nil }
func (m *MockResource) Slug() string                                          { return "users" }
func (m *MockResource) GetDialogType() resource.DialogType                    { return resource.DialogTypeSheet }
func (m *MockResource) SetDialogType(t resource.DialogType) resource.Resource { return m }
func (m *MockResource) Repository(db *gorm.DB) data.DataProvider              { return nil }
func (m *MockResource) GetRecordTitleKey() string                             { return "full_name" }
func (m *MockResource) SetRecordTitleKey(key string) resource.Resource        { return m }
func (m *MockResource) OpenAPIEnabled() bool                                  { return false }

func (m *MockResourceWithCustomRepository) Repository(db *gorm.DB) data.DataProvider {
	return m.provider
}

func (m *MockResourceWithCustomRepository) With() []string {
	return []string{"profile"}
}

func (m *MockResourceWithCustomRepository) Fields() []fields.Element {
	return []fields.Element{
		fields.Text("Full Name", "full_name").Searchable(),
	}
}

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

func TestNewResourceHandler_ConfiguresCustomRepositoryProvider(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to db: %v", err)
	}

	provider := &TrackingDataProvider{}
	res := &MockResourceWithCustomRepository{
		MockResource: &MockResource{},
		provider:     provider,
	}

	h := NewResourceHandler(db, res, "", "")
	if h.Provider != provider {
		t.Fatalf("expected handler provider to be custom provider")
	}

	if !provider.setWithCalled {
		t.Fatalf("expected SetWith to be called")
	}

	if len(provider.with) != 0 {
		t.Fatalf("expected unresolved explicit relations to be filtered out, got %v", provider.with)
	}

	if len(provider.searchColumns) != 1 || provider.searchColumns[0] != "full_name" {
		t.Fatalf("expected SetSearchColumns to include searchable fields, got %v", provider.searchColumns)
	}
}

func TestCollectRelationshipPreloads_NormalizesHasManySnakeCase(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to db: %v", err)
	}

	if err := db.AutoMigrate(&HandlerTestCategory{}, &HandlerTestProduct{}, &HandlerTestVariant{}); err != nil {
		t.Fatalf("Failed to migrate test models: %v", err)
	}

	elements := []fields.Element{
		fields.BelongsTo("Kategori", "category_id", "categories"),
		fields.HasMany("Ürün Varyantları", "product_variants", "product-variants"),
	}

	preloads := collectRelationshipPreloads(db, &HandlerTestProduct{}, []string{}, elements)

	hasCategory := false
	hasProductVariants := false
	for _, preload := range preloads {
		if preload == "Category" {
			hasCategory = true
		}
		if preload == "ProductVariants" {
			hasProductVariants = true
		}
	}

	if !hasCategory {
		t.Fatalf("expected normalized preload to include Category, got %v", preloads)
	}
	if !hasProductVariants {
		t.Fatalf("expected normalized preload to include ProductVariants, got %v", preloads)
	}
}

func TestCollectRelationshipPreloads_NormalizesExplicitSnakeCase(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to db: %v", err)
	}

	if err := db.AutoMigrate(&HandlerTestProduct{}, &HandlerTestVariant{}); err != nil {
		t.Fatalf("Failed to migrate test models: %v", err)
	}

	preloads := collectRelationshipPreloads(
		db,
		&HandlerTestProduct{},
		[]string{"product_variants"},
		nil,
	)

	if len(preloads) != 1 || preloads[0] != "ProductVariants" {
		t.Fatalf("expected explicit snake_case preload to normalize to ProductVariants, got %v", preloads)
	}
}

func TestCollectRelationshipPreloads_DropsUnresolvedExplicitRelations(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to db: %v", err)
	}

	if err := db.AutoMigrate(&HandlerTestProduct{}, &HandlerTestVariant{}); err != nil {
		t.Fatalf("Failed to migrate test models: %v", err)
	}

	preloads := collectRelationshipPreloads(
		db,
		&HandlerTestProduct{},
		[]string{"ProductVariant"},
		nil,
	)

	if len(preloads) != 0 {
		t.Fatalf("expected unresolved explicit relations to be dropped, got %v", preloads)
	}
}

func TestResolveResourceFields_NormalizesNilRelationshipCollectionData(t *testing.T) {
	h := &FieldHandler{}
	item := &HandlerTestProduct{
		ID:              1,
		ProductVariants: nil, // explicit nil slice to simulate non-loaded relationship
	}

	relationshipField := fields.NewField("Variants", "product_variants")
	relationshipField.View = "has-many-field"
	relationshipField.Type = fields.TYPE_RELATIONSHIP

	resolved := h.resolveResourceFields(
		nil,
		&core.ResourceContext{},
		item,
		[]fields.Element{relationshipField},
	)

	fieldData, ok := resolved["product_variants"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected resolved field to be a map, got %T", resolved["product_variants"])
	}

	dataValue, exists := fieldData["data"]
	if !exists {
		t.Fatalf("expected serialized field to contain data key")
	}

	dataSlice, ok := dataValue.([]interface{})
	if !ok {
		t.Fatalf("expected relationship data to be []interface{}, got %T", dataValue)
	}
	if len(dataSlice) != 0 {
		t.Fatalf("expected empty relationship collection, got %v", dataSlice)
	}
}

func TestResolveResourceFields_AppliesDisplayCallbackWithComputedValue(t *testing.T) {
	h := &FieldHandler{}
	item := &HandlerTestProduct{ID: 42}

	computedField := fields.NewField("Sizes", "sizes").Display(func(value interface{}, item interface{}) interface{} {
		if value != nil {
			t.Fatalf("expected computed field source value to be nil, got %v", value)
		}
		product, ok := item.(*HandlerTestProduct)
		if !ok {
			t.Fatalf("expected item type *HandlerTestProduct, got %T", item)
		}
		return fmt.Sprintf("size-for-%d", product.ID)
	})

	resolved := h.resolveResourceFields(
		nil,
		&core.ResourceContext{},
		item,
		[]fields.Element{computedField},
	)

	fieldData, ok := resolved["sizes"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected resolved field to be a map, got %T", resolved["sizes"])
	}

	if got, ok := fieldData["data"].(string); !ok || got != "size-for-42" {
		t.Fatalf("expected computed display value 'size-for-42', got %v (%T)", fieldData["data"], fieldData["data"])
	}
}

func TestResolveResourceFields_AppliesDisplayCallbackElementResult(t *testing.T) {
	h := &FieldHandler{}
	item := &HandlerTestProduct{CategoryID: 7}

	field := fields.Number("Category", "category_id").Display(func(value interface{}, item interface{}) interface{} {
		return fields.Badge("Category")
	})

	resolved := h.resolveResourceFields(
		nil,
		&core.ResourceContext{},
		item,
		[]fields.Element{field},
	)

	fieldData, ok := resolved["category_id"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected resolved field to be a map, got %T", resolved["category_id"])
	}

	componentData, ok := fieldData["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected display data to be a serialized component map, got %T", fieldData["data"])
	}

	if view, _ := componentData["view"].(string); view != "badge-field" {
		t.Fatalf("expected display component view 'badge-field', got %v", componentData["view"])
	}

	if componentData["data"] != uint(7) {
		t.Fatalf("expected fallback badge data to be original value 7, got %v", componentData["data"])
	}
}

func TestResolveResourceFields_AppliesDisplayCallbackStackResult(t *testing.T) {
	h := &FieldHandler{}
	item := &HandlerTestProduct{ID: 9}

	field := fields.Text("Sizes", "sizes").Display(func(value interface{}, item interface{}) core.Element {
		return fields.Stack([]core.Element{
			fields.Badge("10").WithProps("variant", "secondary"),
			fields.Badge("20").WithProps("variant", "secondary"),
		})
	})

	resolved := h.resolveResourceFields(
		nil,
		&core.ResourceContext{},
		item,
		[]fields.Element{field},
	)

	fieldData, ok := resolved["sizes"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected resolved field to be a map, got %T", resolved["sizes"])
	}

	stackData, ok := fieldData["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected display data to be a stack component map, got %T", fieldData["data"])
	}

	if view, _ := stackData["view"].(string); view != "stack-field" {
		t.Fatalf("expected stack component view 'stack-field', got %v", stackData["view"])
	}

	props, ok := stackData["props"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected stack props to be map[string]interface{}, got %T", stackData["props"])
	}

	children, ok := props["fields"].([]map[string]interface{})
	if !ok {
		t.Fatalf("expected stack props.fields to be []map[string]interface{}, got %T", props["fields"])
	}
	if len(children) != 2 {
		t.Fatalf("expected stack to have 2 children, got %d", len(children))
	}
}
