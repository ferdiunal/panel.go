package panel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestManagedAPIKeyLifecycle(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect db: %v", err)
	}

	p := New(Config{
		Database:    DatabaseConfig{Instance: db},
		Environment: "test",
	})

	sessionCookie := registerAndLoginTestUser(t, p, "managed-api-admin@example.com")
	if sessionCookie == nil {
		t.Fatal("session cookie is nil")
	}

	// Enable API key authentication so managed keys can be used.
	enableBody, _ := json.Marshal(map[string]any{
		"api_key_enabled": true,
		"api_key_header":  "X-API-Key",
		"api_keys":        "",
	})
	enableReq := httptest.NewRequest("POST", "/api/internal/pages/api-settings", bytes.NewReader(enableBody))
	enableReq.Header.Set("Content-Type", "application/json")
	enableReq.AddCookie(sessionCookie)
	enableResp, err := testFiberRequest(p.Fiber, enableReq)
	if err != nil {
		t.Fatalf("failed to enable api key auth: %v", err)
	}
	if enableResp.StatusCode != 200 {
		t.Fatalf("expected status 200 from api settings save, got %d", enableResp.StatusCode)
	}

	createBody, _ := json.Marshal(map[string]any{
		"name": "CI Key",
	})
	createReq := httptest.NewRequest("POST", "/api/internal/api-keys", bytes.NewReader(createBody))
	createReq.Header.Set("Content-Type", "application/json")
	createReq.AddCookie(sessionCookie)

	createResp, err := testFiberRequest(p.Fiber, createReq)
	if err != nil {
		t.Fatalf("create api key request failed: %v", err)
	}
	if createResp.StatusCode != 201 {
		t.Fatalf("expected status 201 from api key create, got %d", createResp.StatusCode)
	}

	var created struct {
		Data struct {
			ID uint `json:"id"`
		} `json:"data"`
		Key string `json:"key"`
	}
	if err := json.NewDecoder(createResp.Body).Decode(&created); err != nil {
		t.Fatalf("failed to decode create response: %v", err)
	}
	if created.Data.ID == 0 {
		t.Fatal("created api key id is zero")
	}
	if strings.TrimSpace(created.Key) == "" {
		t.Fatal("created api key secret is empty")
	}

	listReq := httptest.NewRequest("GET", "/api/internal/api-keys", nil)
	listReq.AddCookie(sessionCookie)
	listResp, err := testFiberRequest(p.Fiber, listReq)
	if err != nil {
		t.Fatalf("list api keys request failed: %v", err)
	}
	if listResp.StatusCode != 200 {
		t.Fatalf("expected status 200 from api key list, got %d", listResp.StatusCode)
	}

	validResourceReq := httptest.NewRequest("GET", "/api/internal/resource/users", nil)
	validResourceReq.Header.Set("X-API-Key", created.Key)
	validResourceResp, err := testFiberRequest(p.Fiber, validResourceReq)
	if err != nil {
		t.Fatalf("resource request with managed key failed: %v", err)
	}
	if validResourceResp.StatusCode != 404 {
		t.Fatalf("expected status 404 for internal resource with valid api key, got %d", validResourceResp.StatusCode)
	}

	// API key-authenticated requests should not manage API keys.
	forbiddenMgmtReq := httptest.NewRequest("GET", "/api/internal/api-keys", nil)
	forbiddenMgmtReq.Header.Set("X-API-Key", created.Key)
	forbiddenMgmtResp, err := testFiberRequest(p.Fiber, forbiddenMgmtReq)
	if err != nil {
		t.Fatalf("managed route request with api key failed: %v", err)
	}
	if forbiddenMgmtResp.StatusCode != 403 {
		t.Fatalf("expected status 403 for api-key based management access, got %d", forbiddenMgmtResp.StatusCode)
	}

	revokeReq := httptest.NewRequest("DELETE", fmt.Sprintf("/api/internal/api-keys/%d", created.Data.ID), nil)
	revokeReq.AddCookie(sessionCookie)
	revokeResp, err := testFiberRequest(p.Fiber, revokeReq)
	if err != nil {
		t.Fatalf("revoke api key request failed: %v", err)
	}
	if revokeResp.StatusCode != 200 {
		t.Fatalf("expected status 200 from api key revoke, got %d", revokeResp.StatusCode)
	}

	revokedReq := httptest.NewRequest("GET", "/api/internal/resource/users", nil)
	revokedReq.Header.Set("X-API-Key", created.Key)
	revokedResp, err := testFiberRequest(p.Fiber, revokedReq)
	if err != nil {
		t.Fatalf("resource request with revoked api key failed: %v", err)
	}
	if revokedResp.StatusCode != 401 {
		t.Fatalf("expected status 401 with revoked api key, got %d", revokedResp.StatusCode)
	}
}

func TestOpenAPIDocsEndpointsPublic(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect db: %v", err)
	}

	p := New(Config{
		Database:    DatabaseConfig{Instance: db},
		Environment: "test",
	})

	docsReq := httptest.NewRequest("GET", "/api/docs", nil)
	docsResp, err := testFiberRequest(p.Fiber, docsReq)
	if err != nil {
		t.Fatalf("docs request failed: %v", err)
	}
	if docsResp.StatusCode != 200 {
		t.Fatalf("expected status 200 from /api/docs, got %d", docsResp.StatusCode)
	}
	if !strings.Contains(docsResp.Header.Get("Content-Security-Policy"), "unpkg.com") {
		t.Fatal("expected docs CSP to allow unpkg.com")
	}

	redocReq := httptest.NewRequest("GET", "/api/docs/redoc", nil)
	redocResp, err := testFiberRequest(p.Fiber, redocReq)
	if err != nil {
		t.Fatalf("redoc request failed: %v", err)
	}
	if redocResp.StatusCode == 200 {
		t.Fatalf("expected /api/docs/redoc to be disabled, got status 200")
	}
}

func TestOpenAPISpec_ExcludesSystemResources(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect db: %v", err)
	}

	p := New(Config{
		Database:    DatabaseConfig{Instance: db},
		Environment: "test",
	})

	specReq := httptest.NewRequest("GET", "/api/openapi.json", nil)
	specResp, err := testFiberRequest(p.Fiber, specReq)
	if err != nil {
		t.Fatalf("openapi spec request failed: %v", err)
	}
	if specResp.StatusCode != 200 {
		t.Fatalf("expected status 200 from /api/openapi.json, got %d", specResp.StatusCode)
	}

	body, err := io.ReadAll(specResp.Body)
	if err != nil {
		t.Fatalf("failed to read spec body: %v", err)
	}

	specBody := string(body)
	if !strings.Contains(specBody, "\"apiKeyAuth\"") {
		t.Fatal("openapi spec should include apiKeyAuth security scheme")
	}
	if !strings.Contains(specBody, "\"name\":\"X-External-API-Key\"") {
		t.Fatal("openapi spec should use X-External-API-Key header for apiKeyAuth")
	}
	if strings.Contains(specBody, "/api/users") {
		t.Fatal("openapi spec should not expose system users resource")
	}
	if strings.Contains(specBody, "/api/verifications") {
		t.Fatal("openapi spec should not expose system verifications resource")
	}
	if strings.Contains(specBody, "/api/internal/auth/sign-in/email") {
		t.Fatal("openapi spec should not expose internal auth endpoints")
	}
	if strings.Contains(specBody, "/api/internal/init") {
		t.Fatal("openapi spec should not expose internal init endpoint")
	}
	if strings.Contains(specBody, "/api/internal/navigation") {
		t.Fatal("openapi spec should not expose internal navigation endpoint")
	}
}
