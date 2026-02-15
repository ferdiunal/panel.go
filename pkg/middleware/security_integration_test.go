package middleware_test

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ferdiunal/panel.go/pkg/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRateLimitingIntegration tests rate limiting middleware
func TestRateLimitingIntegration(t *testing.T) {
	app := fiber.New()

	// Apply rate limiting: 5 requests per minute
	app.Use(middleware.RateLimiter(middleware.RateLimitConfig{
		Max:        5,
		Expiration: 1 * time.Minute,
	}))

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// First 5 requests should succeed
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	}

	// 6th request should be rate limited
	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 429, resp.StatusCode)
}

// TestAuthRateLimiting tests strict rate limiting for auth endpoints
func TestAuthRateLimiting(t *testing.T) {
	app := fiber.New()

	authRoutes := app.Group("/auth")
	authRoutes.Use(middleware.AuthRateLimiter())

	authRoutes.Post("/login", func(c *fiber.Ctx) error {
		return c.SendString("Login attempt")
	})

	// First 10 requests should succeed
	for i := 0; i < 10; i++ {
		req := httptest.NewRequest("POST", "/auth/login", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	}

	// 11th request should be rate limited
	req := httptest.NewRequest("POST", "/auth/login", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 429, resp.StatusCode)

	// Verify error message
	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Too many authentication attempts")
}

// TestAccountLockout tests account lockout mechanism
func TestAccountLockout(t *testing.T) {
	lockout := middleware.NewAccountLockout(3, 5*time.Minute)
	t.Cleanup(lockout.Close)

	email := "test@example.com"

	// Account should not be locked initially
	assert.False(t, lockout.IsLocked(email))
	assert.Equal(t, 3, lockout.GetRemainingAttempts(email))

	// Record 3 failed attempts
	for i := 0; i < 3; i++ {
		lockout.RecordFailedAttempt(email)
	}

	// Account should now be locked
	assert.True(t, lockout.IsLocked(email))
	assert.Equal(t, 0, lockout.GetRemainingAttempts(email))

	// Reset should unlock the account
	lockout.ResetAttempts(email)
	assert.False(t, lockout.IsLocked(email))
	assert.Equal(t, 3, lockout.GetRemainingAttempts(email))
}

// TestAccountLockoutExpiration tests that lockout expires after duration
func TestAccountLockoutExpiration(t *testing.T) {
	lockout := middleware.NewAccountLockout(2, 100*time.Millisecond)
	t.Cleanup(lockout.Close)

	email := "test@example.com"

	// Lock the account
	lockout.RecordFailedAttempt(email)
	lockout.RecordFailedAttempt(email)
	assert.True(t, lockout.IsLocked(email))

	// Wait for lockout to expire
	time.Sleep(150 * time.Millisecond)

	// Account should be unlocked
	assert.False(t, lockout.IsLocked(email))
}

// TestSecurityHeaders tests security headers middleware
func TestSecurityHeaders(t *testing.T) {
	app := fiber.New()

	config := middleware.DefaultSecurityHeaders()
	app.Use(middleware.SecurityHeaders(config))

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)

	// Verify security headers are set
	assert.Equal(t, "DENY", resp.Header.Get("X-Frame-Options"))
	assert.Equal(t, "nosniff", resp.Header.Get("X-Content-Type-Options"))
	assert.Equal(t, "no-referrer", resp.Header.Get("Referrer-Policy"))
	assert.Contains(t, resp.Header.Get("Content-Security-Policy"), "default-src 'self'")
	assert.Contains(t, resp.Header.Get("Permissions-Policy"), "geolocation=()")
}

// TestRequestSizeLimit tests request size limiting
func TestRequestSizeLimit(t *testing.T) {
	app := fiber.New()

	// Set 1KB limit
	app.Use(middleware.RequestSizeLimit(1024))

	app.Post("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Small request should succeed
	smallBody := strings.NewReader(strings.Repeat("a", 500))
	req := httptest.NewRequest("POST", "/test", smallBody)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Large request should be rejected
	largeBody := strings.NewReader(strings.Repeat("a", 2000))
	req = httptest.NewRequest("POST", "/test", largeBody)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 413, resp.StatusCode)

	// Verify error message
	body, _ := io.ReadAll(resp.Body)
	assert.Contains(t, string(body), "Request body too large")
}

// TestAuditLogging tests audit logging middleware
func TestAuditLogging(t *testing.T) {
	app := fiber.New()

	// Create a test logger that captures events
	var capturedEvents []middleware.AuditEvent
	testLogger := &testAuditLogger{events: &capturedEvents}

	app.Use(middleware.AuditMiddleware(testLogger))

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("User-Agent", "TestAgent")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Verify audit event was logged
	require.Len(t, capturedEvents, 1)
	event := capturedEvents[0]
	assert.Equal(t, "GET", event.Method)
	assert.Equal(t, "/test", event.Path)
	assert.Equal(t, "TestAgent", event.UserAgent)
	assert.Equal(t, 200, event.StatusCode)
	assert.True(t, event.Success)
}

// TestAuditLoggingAuthEvents tests audit logging for authentication events
func TestAuditLoggingAuthEvents(t *testing.T) {
	app := fiber.New()

	var capturedEvents []middleware.AuditEvent
	testLogger := &testAuditLogger{events: &capturedEvents}

	app.Use(middleware.AuditMiddleware(testLogger))

	// Simulate successful login
	app.Post("/auth/sign-in", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"success": true})
	})

	req := httptest.NewRequest("POST", "/auth/sign-in", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Verify login_success event was logged
	require.Len(t, capturedEvents, 1)
	event := capturedEvents[0]
	assert.Equal(t, "login_success", event.EventType)
	assert.True(t, event.Success)

	// Simulate failed login
	capturedEvents = nil
	app.Post("/auth/sign-in-fail", func(c *fiber.Ctx) error {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	})

	req = httptest.NewRequest("POST", "/auth/sign-in-fail", nil)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 401, resp.StatusCode)

	// Verify login_failure event was logged
	require.Len(t, capturedEvents, 1)
	event = capturedEvents[0]
	assert.Equal(t, "login_failure", event.EventType) // Path contains /auth/sign-in so it's detected as login_failure
	assert.False(t, event.Success)
}

// TestCORSValidation tests CORS origin validation
func TestCORSValidation(t *testing.T) {
	allowedOrigins := []string{"https://example.com", "https://app.example.com"}
	validator := middleware.ValidateCORSOrigin(allowedOrigins)

	// Allowed origins should pass
	assert.True(t, validator("https://example.com"))
	assert.True(t, validator("https://app.example.com"))

	// Disallowed origins should fail
	assert.False(t, validator("https://evil.com"))
	assert.False(t, validator("https://example.org"))

	// Wildcard should allow all
	wildcardValidator := middleware.ValidateCORSOrigin([]string{"*"})
	assert.True(t, wildcardValidator("https://anything.com"))

	// Wildcard subdomain support
	subdomainValidator := middleware.ValidateCORSOrigin([]string{"*.example.com"})
	assert.True(t, subdomainValidator("https://api.example.com"))
	assert.True(t, subdomainValidator("https://app.example.com"))
	assert.False(t, subdomainValidator("https://example.com"))
}

// TestIntegratedSecurityStack tests all security middleware together
func TestIntegratedSecurityStack(t *testing.T) {
	app := fiber.New()

	// Apply all security middleware
	app.Use(middleware.SecurityHeaders(middleware.DefaultSecurityHeaders()))
	app.Use(middleware.RequestSizeLimit(1024))

	var capturedEvents []middleware.AuditEvent
	testLogger := &testAuditLogger{events: &capturedEvents}
	app.Use(middleware.AuditMiddleware(testLogger))

	app.Use(middleware.RateLimiter(middleware.RateLimitConfig{
		Max:        3,
		Expiration: 1 * time.Minute,
	}))

	app.Post("/api/test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"success": true})
	})

	// Test 1: Normal request should succeed
	body := strings.NewReader(`{"test":"data"}`)
	req := httptest.NewRequest("POST", "/api/test", body)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	// Verify security headers
	assert.Equal(t, "DENY", resp.Header.Get("X-Frame-Options"))

	// Verify audit log
	require.Len(t, capturedEvents, 1)
	assert.True(t, capturedEvents[0].Success)

	// Test 2: Oversized request should be rejected
	largeBody := strings.NewReader(strings.Repeat("a", 2000))
	req = httptest.NewRequest("POST", "/api/test", largeBody)
	resp, err = app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, 413, resp.StatusCode)

	// Test 3: Rate limiting should kick in after 3 requests
	for i := 0; i < 3; i++ {
		body := strings.NewReader(`{"test":"data"}`)
		req := httptest.NewRequest("POST", "/api/test", body)
		resp, err := app.Test(req)
		require.NoError(t, err)
		if i < 2 {
			assert.Equal(t, 200, resp.StatusCode)
		} else {
			assert.Equal(t, 429, resp.StatusCode)
		}
	}
}

// testAuditLogger is a test implementation of AuditLogger
type testAuditLogger struct {
	events *[]middleware.AuditEvent
}

func (l *testAuditLogger) Log(event middleware.AuditEvent) error {
	*l.events = append(*l.events, event)
	return nil
}

// TestAccountLockoutConcurrency tests account lockout under concurrent access
func TestAccountLockoutConcurrency(t *testing.T) {
	lockout := middleware.NewAccountLockout(5, 1*time.Minute)
	t.Cleanup(lockout.Close)
	email := "concurrent@example.com"

	// Simulate concurrent failed login attempts
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			lockout.RecordFailedAttempt(email)
			done <- true
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Account should be locked
	assert.True(t, lockout.IsLocked(email))
	assert.Equal(t, 0, lockout.GetRemainingAttempts(email))
}

// TestSecurityHeadersCustomization tests custom security headers
func TestSecurityHeadersCustomization(t *testing.T) {
	app := fiber.New()

	customConfig := middleware.SecurityHeadersConfig{
		ContentSecurityPolicy: "default-src 'none'",
		XFrameOptions:         "SAMEORIGIN",
		XContentTypeOptions:   "nosniff",
		ReferrerPolicy:        "strict-origin",
		PermissionsPolicy:     "camera=(), microphone=()",
	}

	app.Use(middleware.SecurityHeaders(customConfig))

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)

	// Verify custom headers
	assert.Equal(t, "default-src 'none'", resp.Header.Get("Content-Security-Policy"))
	assert.Equal(t, "SAMEORIGIN", resp.Header.Get("X-Frame-Options"))
	assert.Equal(t, "strict-origin", resp.Header.Get("Referrer-Policy"))
	assert.Equal(t, "camera=(), microphone=()", resp.Header.Get("Permissions-Policy"))
}

// BenchmarkRateLimiter benchmarks rate limiting performance
func BenchmarkRateLimiter(b *testing.B) {
	app := fiber.New()
	app.Use(middleware.RateLimiter(middleware.RateLimitConfig{
		Max:        1000,
		Expiration: 1 * time.Minute,
	}))

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		_, _ = app.Test(req)
	}
}

// BenchmarkAccountLockout benchmarks account lockout performance
func BenchmarkAccountLockout(b *testing.B) {
	lockout := middleware.NewAccountLockout(5, 15*time.Minute)
	b.Cleanup(lockout.Close)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		email := "bench@example.com"
		lockout.IsLocked(email)
		lockout.RecordFailedAttempt(email)
		lockout.GetRemainingAttempts(email)
	}
}

// BenchmarkAuditLogging benchmarks audit logging performance
func BenchmarkAuditLogging(b *testing.B) {
	app := fiber.New()
	app.Use(middleware.AuditMiddleware(&middleware.ConsoleAuditLogger{}))

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		_, _ = app.Test(req)
	}
}
