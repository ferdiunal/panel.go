package middleware

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	// Max number of requests per window
	Max int
	// Time window duration
	Expiration time.Duration
	// Custom key generator (default: IP address)
	KeyGenerator func(*fiber.Ctx) string
	// Custom response when limit exceeded
	LimitReached fiber.Handler
}

// SecurityHeadersConfig holds security headers configuration
type SecurityHeadersConfig struct {
	// Content Security Policy
	ContentSecurityPolicy string
	// X-Frame-Options
	XFrameOptions string
	// X-Content-Type-Options
	XContentTypeOptions string
	// Referrer-Policy
	ReferrerPolicy string
	// Permissions-Policy
	PermissionsPolicy string
}

// DefaultSecurityHeaders returns secure default headers
func DefaultSecurityHeaders() SecurityHeadersConfig {
	return SecurityHeadersConfig{
		ContentSecurityPolicy: "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' data:; connect-src 'self'; frame-ancestors 'none';",
		XFrameOptions:         "DENY",
		XContentTypeOptions:   "nosniff",
		ReferrerPolicy:        "no-referrer",
		PermissionsPolicy:     "geolocation=(), microphone=(), camera=()",
	}
}

// SecurityHeaders middleware adds security headers to responses
func SecurityHeaders(config SecurityHeadersConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if config.ContentSecurityPolicy != "" {
			c.Set("Content-Security-Policy", config.ContentSecurityPolicy)
		}
		if config.XFrameOptions != "" {
			c.Set("X-Frame-Options", config.XFrameOptions)
		}
		if config.XContentTypeOptions != "" {
			c.Set("X-Content-Type-Options", config.XContentTypeOptions)
		}
		if config.ReferrerPolicy != "" {
			c.Set("Referrer-Policy", config.ReferrerPolicy)
		}
		if config.PermissionsPolicy != "" {
			c.Set("Permissions-Policy", config.PermissionsPolicy)
		}
		return c.Next()
	}
}

// RateLimiter creates a rate limiting middleware
func RateLimiter(config RateLimitConfig) fiber.Handler {
	if config.Max == 0 {
		config.Max = 100
	}
	if config.Expiration == 0 {
		config.Expiration = 1 * time.Minute
	}
	if config.KeyGenerator == nil {
		config.KeyGenerator = func(c *fiber.Ctx) string {
			return c.IP()
		}
	}

	return limiter.New(limiter.Config{
		Max:        config.Max,
		Expiration: config.Expiration,
		KeyGenerator: func(c *fiber.Ctx) string {
			return config.KeyGenerator(c)
		},
		LimitReached: config.LimitReached,
	})
}

// AuthRateLimiter creates a strict rate limiter for authentication endpoints
func AuthRateLimiter() fiber.Handler {
	return RateLimiter(RateLimitConfig{
		Max:        10, // 10 requests per minute
		Expiration: 1 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many authentication attempts. Please try again later.",
			})
		},
	})
}

// APIRateLimiter creates a rate limiter for general API endpoints
func APIRateLimiter() fiber.Handler {
	return RateLimiter(RateLimitConfig{
		Max:        100, // 100 requests per minute
		Expiration: 1 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Rate limit exceeded. Please slow down.",
			})
		},
	})
}

// AccountLockout implements account lockout after failed login attempts
type AccountLockout struct {
	mu       sync.RWMutex
	attempts map[string]*lockoutEntry
	maxAttempts int
	lockoutDuration time.Duration
}

type lockoutEntry struct {
	count      int
	lockedUntil time.Time
}

// NewAccountLockout creates a new account lockout manager
func NewAccountLockout(maxAttempts int, lockoutDuration time.Duration) *AccountLockout {
	al := &AccountLockout{
		attempts:        make(map[string]*lockoutEntry),
		maxAttempts:     maxAttempts,
		lockoutDuration: lockoutDuration,
	}

	// Cleanup goroutine
	go al.cleanup()

	return al
}

// IsLocked checks if an account is locked
func (al *AccountLockout) IsLocked(identifier string) bool {
	al.mu.RLock()
	defer al.mu.RUnlock()

	entry, exists := al.attempts[identifier]
	if !exists {
		return false
	}

	if time.Now().Before(entry.lockedUntil) {
		return true
	}

	return false
}

// RecordFailedAttempt records a failed login attempt
func (al *AccountLockout) RecordFailedAttempt(identifier string) {
	al.mu.Lock()
	defer al.mu.Unlock()

	entry, exists := al.attempts[identifier]
	if !exists {
		entry = &lockoutEntry{count: 0}
		al.attempts[identifier] = entry
	}

	// Reset if lockout expired (only check if lockedUntil was previously set)
	if !entry.lockedUntil.IsZero() && time.Now().After(entry.lockedUntil) {
		entry.count = 0
	}

	entry.count++

	if entry.count >= al.maxAttempts {
		entry.lockedUntil = time.Now().Add(al.lockoutDuration)
	}
}

// ResetAttempts resets failed attempts for an identifier (after successful login)
func (al *AccountLockout) ResetAttempts(identifier string) {
	al.mu.Lock()
	defer al.mu.Unlock()

	delete(al.attempts, identifier)
}

// GetRemainingAttempts returns remaining attempts before lockout
func (al *AccountLockout) GetRemainingAttempts(identifier string) int {
	al.mu.RLock()
	defer al.mu.RUnlock()

	entry, exists := al.attempts[identifier]
	if !exists {
		return al.maxAttempts
	}

	// Only check expiration if lockout was previously set
	if !entry.lockedUntil.IsZero() && time.Now().After(entry.lockedUntil) {
		return al.maxAttempts
	}

	remaining := al.maxAttempts - entry.count
	if remaining < 0 {
		return 0
	}
	return remaining
}

// cleanup removes expired entries periodically
func (al *AccountLockout) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		al.mu.Lock()
		now := time.Now()
		for key, entry := range al.attempts {
			if now.After(entry.lockedUntil) && entry.count < al.maxAttempts {
				delete(al.attempts, key)
			}
		}
		al.mu.Unlock()
	}
}

// RequestSizeLimit middleware limits request body size
func RequestSizeLimit(maxSize int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if len(c.Body()) > maxSize {
			return c.Status(fiber.StatusRequestEntityTooLarge).JSON(fiber.Map{
				"error": fmt.Sprintf("Request body too large. Maximum size: %d bytes", maxSize),
			})
		}
		return c.Next()
	}
}

// ValidateCORSOrigin validates CORS origin against whitelist
func ValidateCORSOrigin(allowedOrigins []string) func(string) bool {
	return func(origin string) bool {
		for _, allowed := range allowedOrigins {
			if allowed == "*" {
				return true
			}
			if strings.EqualFold(origin, allowed) {
				return true
			}
			// Support wildcard subdomains like *.example.com
			if strings.HasPrefix(allowed, "*.") {
				domain := strings.TrimPrefix(allowed, "*.")
				// Extract host from origin (remove protocol)
				if strings.Contains(origin, "://") {
					parts := strings.Split(origin, "://")
					if len(parts) == 2 {
						host := parts[1]
						// Match if host ends with domain and is not exactly the domain
						if strings.HasSuffix(host, domain) && host != domain {
							return true
						}
					}
				}
			}
		}
		return false
	}
}
