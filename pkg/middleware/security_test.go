package middleware

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestRateLimiter(t *testing.T) {
	app := fiber.New()

	// Add rate limiter: 3 requests per 10 seconds
	app.Use(RateLimiter(RateLimitConfig{
		Max:        3,
		Expiration: 10 * time.Second,
	}))

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// First 3 requests should succeed
	for i := 0; i < 3; i++ {
		req := httptest.NewRequest("GET", "/test", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	}

	// 4th request should be rate limited
	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 429, resp.StatusCode)
}

func TestAccountLockout(t *testing.T) {
	lockout := NewAccountLockout(3, 5*time.Minute)
	t.Cleanup(lockout.Close)

	email := "test@example.com"

	// Should not be locked initially
	assert.False(t, lockout.IsLocked(email))
	assert.Equal(t, 3, lockout.GetRemainingAttempts(email))

	// Record 2 failed attempts
	lockout.RecordFailedAttempt(email)
	lockout.RecordFailedAttempt(email)

	// Should not be locked yet
	assert.False(t, lockout.IsLocked(email))
	assert.Equal(t, 1, lockout.GetRemainingAttempts(email))

	// Record 3rd failed attempt
	lockout.RecordFailedAttempt(email)

	// Should be locked now
	assert.True(t, lockout.IsLocked(email))
	assert.Equal(t, 0, lockout.GetRemainingAttempts(email))

	// Reset attempts (successful login)
	lockout.ResetAttempts(email)

	// Should not be locked anymore
	assert.False(t, lockout.IsLocked(email))
	assert.Equal(t, 3, lockout.GetRemainingAttempts(email))
}

func TestAccountLockoutClose_Idempotent(t *testing.T) {
	lockout := NewAccountLockout(2, time.Minute)
	lockout.Close()
	lockout.Close()
}

func TestSecurityHeaders(t *testing.T) {
	app := fiber.New()

	app.Use(SecurityHeaders(DefaultSecurityHeaders()))

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)

	// Check security headers are set
	assert.Contains(t, resp.Header.Get("Content-Security-Policy"), "default-src 'self'")
	assert.Equal(t, "DENY", resp.Header.Get("X-Frame-Options"))
	assert.Equal(t, "nosniff", resp.Header.Get("X-Content-Type-Options"))
	assert.Equal(t, "no-referrer", resp.Header.Get("Referrer-Policy"))
}

func TestValidateCORSOrigin(t *testing.T) {
	validator := ValidateCORSOrigin([]string{
		"https://example.com",
		"https://app.example.com",
		"*.subdomain.com",
	})

	// Should allow exact matches
	assert.True(t, validator("https://example.com"))
	assert.True(t, validator("https://app.example.com"))

	// Should allow wildcard subdomains
	assert.True(t, validator("https://test.subdomain.com"))
	assert.True(t, validator("https://api.subdomain.com"))

	// Should reject non-matching origins
	assert.False(t, validator("https://evil.com"))
	assert.False(t, validator("https://example.com.evil.com"))
}
