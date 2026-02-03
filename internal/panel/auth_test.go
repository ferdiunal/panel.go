package panel

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ferdiunal/panel.go/internal/domain/user"
	"github.com/stretchr/testify/assert"
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
	}

	app := New(cfg)

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
		req, _ := http.NewRequest("POST", "/api/auth/sign-up/email", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Fiber.Test(req)
		assert.NoError(t, err)
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
		req, _ := http.NewRequest("POST", "/api/auth/sign-in/email", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Fiber.Test(req)
		assert.NoError(t, err)
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
		req, _ := http.NewRequest("GET", "/api/auth/session", nil)
		cookie := &http.Cookie{Name: "session_token", Value: sessionToken}
		req.AddCookie(cookie)
		resp, err := app.Fiber.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		sessionData := res["session"].(map[string]interface{})
		assert.Equal(t, sessionToken, sessionData["token"])
	})

	// 4. Sign Out
	t.Run("SignOut", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/auth/sign-out", nil)
		cookie := &http.Cookie{Name: "session_token", Value: sessionToken}
		req.AddCookie(cookie)
		resp, err := app.Fiber.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verify Session is gone
		req2, _ := http.NewRequest("GET", "/api/auth/session", nil)
		req2.AddCookie(cookie) // Send old cookie
		resp2, err := app.Fiber.Test(req2)
		assert.NoError(t, err)

		var res map[string]interface{}
		json.NewDecoder(resp2.Body).Decode(&res)
		assert.Nil(t, res["session"])
	})
}
