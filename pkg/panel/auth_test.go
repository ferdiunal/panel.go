package panel

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ferdiunal/panel.go/pkg/domain/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestAuthFlow(t *testing.T) {
	// Setup InMemory DB
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)

	// App Config
	cfg := Config{
		Server: ServerConfig{Host: "localhost", Port: "3004"},
		Database: DatabaseConfig{
			Instance: db,
		},
		Environment: "test",
		Features: FeatureConfig{
			Register: true, // Enable registration for tests
		},
	}

	app := New(cfg)
	doReq := func(t *testing.T, req *http.Request) *http.Response {
		t.Helper()
		resp, err := app.Fiber.Test(req, 10000)
		require.NoError(t, err)
		return resp
	}

	// User Data
	email := "auth_test@example.com"
	password := "password123"
	name := "Auth Test User"

	// 1. Register
	t.Run("Register", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{
			"name":     name,
			"email":    email,
			"password": password,
		})
		req, _ := http.NewRequest("POST", "/api/internal/auth/sign-up/email", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp := doReq(t, req)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var res map[string]*user.User
		json.NewDecoder(resp.Body).Decode(&res)
		assert.Equal(t, email, res["user"].Email)
	})

	// 2. Login
	var sessionToken string
	t.Run("Login", func(t *testing.T) {
		body, _ := json.Marshal(map[string]string{
			"email":    email,
			"password": password,
		})
		req, _ := http.NewRequest("POST", "/api/internal/auth/sign-in/email", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp := doReq(t, req)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Check Cookie
		foundCookie := false
		for _, cookie := range resp.Cookies() {
			if cookie.Name == "session_token" {
				sessionToken = cookie.Value
				foundCookie = true
				break
			}
		}
		assert.True(t, foundCookie, "Session token cookie should be set")
	})

	// 3. Get Session
	t.Run("GetSession", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/internal/auth/session", nil)
		cookie := &http.Cookie{Name: "session_token", Value: sessionToken}
		req.AddCookie(cookie)
		resp := doReq(t, req)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		sessionData := res["session"].(map[string]interface{})
		assert.Equal(t, sessionToken, sessionData["token"])
	})

	// 4. Sign Out
	t.Run("SignOut", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/internal/auth/sign-out", nil)
		cookie := &http.Cookie{Name: "session_token", Value: sessionToken}
		req.AddCookie(cookie)
		resp := doReq(t, req)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify Session is gone
		req2, _ := http.NewRequest("GET", "/api/internal/auth/session", nil)
		req2.AddCookie(cookie) // Send old cookie
		resp2 := doReq(t, req2)

		var res map[string]interface{}
		json.NewDecoder(resp2.Body).Decode(&res)
		assert.Nil(t, res["session"])
	})
}
