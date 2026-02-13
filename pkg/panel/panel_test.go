package panel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ferdiunal/panel.go/pkg/auth"
	appContext "github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/data"
	"github.com/ferdiunal/panel.go/pkg/domain/user"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/resource"
	"github.com/ferdiunal/panel.go/pkg/widget"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Mock User Model
type MockUser struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "test_users"
}

// Mock Resource
type UserResource struct{}

func (u *UserResource) Model() interface{} {
	return &User{}
}

func (u *UserResource) Fields() []fields.Element {
	return []fields.Element{
		fields.ID(),
		fields.Text("Name", "name"),
	}
}

func (u *UserResource) With() []string {
	return nil
}

func (u *UserResource) Lenses() []resource.Lens {
	return nil
}

func (u *UserResource) Title() string                                         { return "Users" }
func (u *UserResource) Icon() string                                          { return "users" }
func (u *UserResource) Group() string                                         { return "System" }
func (u *UserResource) Policy() auth.Policy                                   { return nil }
func (r *UserResource) GetSortable() []resource.Sortable                      { return nil }
func (u *UserResource) Slug() string                                          { return "users" }
func (u *UserResource) GetDialogType() resource.DialogType                    { return resource.DialogTypeSheet }
func (u *UserResource) SetDialogType(t resource.DialogType) resource.Resource { return u }
func (u *UserResource) Repository(db *gorm.DB) data.DataProvider {
	return nil
}

func (u *UserResource) Cards() []widget.Card {
	return []widget.Card{}
}

func (u *UserResource) NavigationOrder() int {
	return 1
}

func (u *UserResource) Visible() bool {
	return true
}

func (u *UserResource) StoreHandler(c *appContext.Context, file *multipart.FileHeader, storagePath string, storageURL string) (string, error) {
	return "", nil
}

func (u *UserResource) GetFields(ctx *appContext.Context) []fields.Element {
	return u.Fields()
}

func (u *UserResource) GetCards(ctx *appContext.Context) []widget.Card {
	return u.Cards()
}

func (u *UserResource) GetLenses() []resource.Lens {
	return u.Lenses()
}

func (u *UserResource) GetPolicy() auth.Policy {
	return u.Policy()
}

func (u *UserResource) ResolveField(fieldName string, item interface{}) (interface{}, error) {
	return nil, nil
}

func (u *UserResource) GetActions() []resource.Action {
	return []resource.Action{}
}

func (u *UserResource) GetFilters() []resource.Filter {
	return []resource.Filter{}
}

func TestPanel_DynamicRouting(t *testing.T) {
	// Setup DB
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to db: %v", err)
	}
	if err := db.AutoMigrate(&User{}); err != nil {
		t.Fatalf("Failed to auto migrate: %v", err)
	}
	db.Create(&User{Name: "Panel User"})

	// Setup Config
	cfg := Config{
		Database: DatabaseConfig{
			Instance: db,
		},
		Environment: "test",
	}

	// Initialize Panel
	app := New(cfg)

	// Register Resource
	app.Register("users", &UserResource{})

	// Auth (Register & Login)
	// Ensure User table for Auth

	// Register
	regBody, _ := json.Marshal(map[string]string{"name": "Panel Tester", "email": "panel@example.com", "password": "password"})
	regReq := httptest.NewRequest("POST", "/api/auth/sign-up/email", bytes.NewReader(regBody))
	regReq.Header.Set("Content-Type", "application/json")
	app.Fiber.Test(regReq)

	// Login
	loginBody, _ := json.Marshal(map[string]string{"email": "panel@example.com", "password": "password"})
	loginReq := httptest.NewRequest("POST", "/api/auth/sign-in/email", bytes.NewReader(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginResp, _ := app.Fiber.Test(loginReq)

	var sessionCookie *http.Cookie
	for _, cookie := range loginResp.Cookies() {
		if cookie.Name == "session_token" {
			sessionCookie = cookie
			break
		}
	}

	// Test Request to Dynamic Route
	// Test Request to Dynamic Route
	req := httptest.NewRequest("GET", "/api/resource/users", nil)
	if sessionCookie != nil {
		req.AddCookie(sessionCookie)
	}
	resp, err := app.Fiber.Test(req)
	if err != nil {
		t.Fatalf("Request failed: %v", err)
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

	item := dataList[0].(map[string]interface{})
	nameField := item["name"].(map[string]interface{})
	if nameField["data"] != "Panel User" {
		t.Errorf("Expected Panel User, got %v", nameField["data"])
	}
}

func TestPanel_ResourceNotFound(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to db: %v", err)
	}
	app := New(Config{
		Database:    DatabaseConfig{Instance: db},
		Environment: "test",
	})

	// Auth (Register & Login)
	db.AutoMigrate(&user.User{})
	// Register
	regBody, _ := json.Marshal(map[string]string{"name": "RNF", "email": "rnf@example.com", "password": "password"})
	regReq := httptest.NewRequest("POST", "/api/auth/sign-up/email", bytes.NewReader(regBody))
	regReq.Header.Set("Content-Type", "application/json")
	app.Fiber.Test(regReq)

	// Login
	loginBody, _ := json.Marshal(map[string]string{"email": "rnf@example.com", "password": "password"})
	loginReq := httptest.NewRequest("POST", "/api/auth/sign-in/email", bytes.NewReader(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginResp, _ := app.Fiber.Test(loginReq)

	var sessionCookie *http.Cookie
	for _, cookie := range loginResp.Cookies() {
		if cookie.Name == "session_token" {
			sessionCookie = cookie
			break
		}
	}

	// Test Request to Non-existent Resource
	req := httptest.NewRequest("GET", "/api/resource/unknown", nil)
	if sessionCookie != nil {
		req.AddCookie(sessionCookie)
	}
	resp, _ := app.Fiber.Test(req)

	if resp.StatusCode != 404 {
		t.Errorf("Expected status 404, got %d", resp.StatusCode)
	}
}

func TestPanel_CRUD(t *testing.T) {
	// Setup DB
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to db: %v", err)
	}
	if err := db.AutoMigrate(&User{}); err != nil {
		t.Fatalf("Failed to auto migrate: %v", err)
	}

	app := New(Config{
		Database:    DatabaseConfig{Instance: db},
		Environment: "test",
	})
	app.Register("users", &UserResource{})

	// Auth (Register & Login)
	db.AutoMigrate(&user.User{}) // Ensure User table
	// Register
	regBody, _ := json.Marshal(map[string]string{"name": "CRUD", "email": "crud@example.com", "password": "password"})
	regReq := httptest.NewRequest("POST", "/api/auth/sign-up/email", bytes.NewReader(regBody))
	regReq.Header.Set("Content-Type", "application/json")
	app.Fiber.Test(regReq)

	// Login
	loginBody, _ := json.Marshal(map[string]string{"email": "crud@example.com", "password": "password"})
	loginReq := httptest.NewRequest("POST", "/api/auth/sign-in/email", bytes.NewReader(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginResp, _ := app.Fiber.Test(loginReq)

	var sessionCookie *http.Cookie
	for _, cookie := range loginResp.Cookies() {
		if cookie.Name == "session_token" {
			sessionCookie = cookie
			break
		}
	}

	// 1. CREATE
	createBody, _ := json.Marshal(map[string]interface{}{
		"name": "New User",
	})
	req := httptest.NewRequest("POST", "/api/resource/users", bytes.NewReader(createBody))
	req.Header.Set("Content-Type", "application/json")
	if sessionCookie != nil {
		req.AddCookie(sessionCookie)
	}
	resp, err := app.Fiber.Test(req)
	if err != nil {
		t.Fatalf("Create request failed: %v", err)
	}
	if resp.StatusCode != 201 {
		t.Errorf("Expected status 201 Created, got %d", resp.StatusCode)
	}

	respBody, _ := io.ReadAll(resp.Body)
	var createResp map[string]interface{}
	if err := json.Unmarshal(respBody, &createResp); err != nil {
		t.Fatalf("Failed to unmarshal create response: %v", err)
	}

	// Create returns rich object now
	itemData := createResp["data"].(map[string]interface{})
	idData := itemData["id"].(map[string]interface{})
	id := idData["data"].(float64) // JSON numbers are float64

	// 2. SHOW
	req = httptest.NewRequest("GET", fmt.Sprintf("/api/resource/users/%d", int(id)), nil)
	if sessionCookie != nil {
		req.AddCookie(sessionCookie)
	}
	resp, _ = app.Fiber.Test(req)
	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200 Show, got %d", resp.StatusCode)
	}

	// 3. UPDATE
	putBody := `{"name": "Updated User"}`
	req = httptest.NewRequest("PUT", fmt.Sprintf("/api/resource/users/%d", int(id)), strings.NewReader(putBody))
	req.Header.Set("Content-Type", "application/json")
	if sessionCookie != nil {
		req.AddCookie(sessionCookie)
	}
	resp, _ = app.Fiber.Test(req)
	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200 Update, got %d", resp.StatusCode)
	}

	// Verify Update
	req = httptest.NewRequest("GET", fmt.Sprintf("/api/resource/users/%d", int(id)), nil)
	if sessionCookie != nil {
		req.AddCookie(sessionCookie)
	}
	resp, _ = app.Fiber.Test(req)
	body, _ := io.ReadAll(resp.Body)
	var showResp map[string]interface{}
	if err := json.Unmarshal(body, &showResp); err != nil {
		t.Fatalf("Failed to unmarshal show response: %v", err)
	}

	// Data is now a rich object: nested field structures
	itemData = showResp["data"].(map[string]interface{})
	nameData := itemData["name"].(map[string]interface{})

	if nameData["data"] != "Updated User" {
		t.Errorf("Update failed, name matches %v", nameData["data"])
	}

	// 4. DELETE
	req = httptest.NewRequest("DELETE", fmt.Sprintf("/api/resource/users/%d", int(id)), nil)
	if sessionCookie != nil {
		req.AddCookie(sessionCookie)
	}
	resp, _ = app.Fiber.Test(req)
	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200 Delete, got %d", resp.StatusCode)
	}

	// Verify Delete (Expect 404)
	req = httptest.NewRequest("GET", fmt.Sprintf("/api/resource/users/%d", int(id)), nil)
	if sessionCookie != nil {
		req.AddCookie(sessionCookie)
	}
	resp, _ = app.Fiber.Test(req)
	if resp.StatusCode != 404 {
		t.Errorf("Expected status 404 after delete, got %d", resp.StatusCode)
	}
}
