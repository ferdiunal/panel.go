package panel

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ferdiunal/panel.go/pkg/domain/user"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestAPIKeyPageSave_RefreshesMiddleware(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect db: %v", err)
	}

	p := New(Config{
		Database:    DatabaseConfig{Instance: db},
		Environment: "test",
	})

	sessionCookie := registerAndLoginTestUser(t, p, "apikey-admin@example.com")
	if sessionCookie == nil {
		t.Fatal("session cookie is nil")
	}

	saveBody, _ := json.Marshal(map[string]any{
		"api_key_enabled": true,
		"api_key_header":  "X-App-Key",
		"api_keys":        "client-secret-key",
	})

	saveReq := httptest.NewRequest("POST", "/api/pages/api-settings", bytes.NewReader(saveBody))
	saveReq.Header.Set("Content-Type", "application/json")
	saveReq.AddCookie(sessionCookie)

	saveResp, err := testFiberRequest(p.Fiber, saveReq)
	if err != nil {
		t.Fatalf("api settings save request failed: %v", err)
	}
	if saveResp.StatusCode != 200 {
		t.Fatalf("expected status 200 from api settings save, got %d", saveResp.StatusCode)
	}

	validKeyReq := httptest.NewRequest("GET", "/api/resource/users", nil)
	validKeyReq.Header.Set("X-App-Key", "client-secret-key")

	validKeyResp, err := testFiberRequest(p.Fiber, validKeyReq)
	if err != nil {
		t.Fatalf("resource request with valid api key failed: %v", err)
	}
	if validKeyResp.StatusCode != 404 {
		t.Fatalf("expected status 404 with valid api key on internal resource, got %d", validKeyResp.StatusCode)
	}

	invalidKeyReq := httptest.NewRequest("GET", "/api/resource/users", nil)
	invalidKeyReq.Header.Set("X-App-Key", "wrong-key")

	invalidKeyResp, err := testFiberRequest(p.Fiber, invalidKeyReq)
	if err != nil {
		t.Fatalf("resource request with invalid api key failed: %v", err)
	}
	if invalidKeyResp.StatusCode != 401 {
		t.Fatalf("expected status 401 with invalid api key, got %d", invalidKeyResp.StatusCode)
	}
}

func registerAndLoginTestUser(t *testing.T, p *Panel, email string) *http.Cookie {
	t.Helper()

	registerBody, _ := json.Marshal(map[string]string{
		"name":     "API Key Admin",
		"email":    email,
		"password": "password",
	})
	registerReq := httptest.NewRequest("POST", "/api/auth/sign-up/email", bytes.NewReader(registerBody))
	registerReq.Header.Set("Content-Type", "application/json")

	if _, err := testFiberRequest(p.Fiber, registerReq); err != nil {
		t.Fatalf("register request failed: %v", err)
	}

	if p != nil && p.Db != nil {
		_ = p.Db.Model(&user.User{}).Where("email = ?", email).Update("role", "admin").Error
	}

	loginBody, _ := json.Marshal(map[string]string{
		"email":    email,
		"password": "password",
	})
	loginReq := httptest.NewRequest("POST", "/api/auth/sign-in/email", bytes.NewReader(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")

	loginResp, err := testFiberRequest(p.Fiber, loginReq)
	if err != nil {
		t.Fatalf("login request failed: %v", err)
	}

	for _, cookie := range loginResp.Cookies() {
		if cookie.Name == "session_token" {
			return cookie
		}
	}

	return nil
}
