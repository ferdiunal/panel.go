package middleware

import (
	"crypto/subtle"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2"
)

const (
	// APIKeyAuthenticatedLocalKey indicates that the request was authenticated via API key.
	APIKeyAuthenticatedLocalKey = "api_key_authenticated"
)

// APIKeyAuth validates incoming API key headers.
type APIKeyAuth struct {
	mu               sync.RWMutex
	enabled          bool
	header           string
	keys             []string
	dynamicValidator APIKeyValidator
}

// APIKeyValidator validates managed/dynamic keys (e.g. DB-backed keys).
type APIKeyValidator func(c *fiber.Ctx, incoming string) bool

// NewAPIKeyAuth creates a new API key middleware module.
func NewAPIKeyAuth(enabled bool, header string, keys []string) *APIKeyAuth {
	auth := &APIKeyAuth{}
	auth.SetConfig(enabled, header, keys)
	return auth
}

// SetConfig updates the middleware config at runtime.
func (a *APIKeyAuth) SetConfig(enabled bool, header string, keys []string) {
	if a == nil {
		return
	}

	normalizedHeader := strings.TrimSpace(header)
	if normalizedHeader == "" {
		normalizedHeader = "X-API-Key"
	}

	normalizedKeys := make([]string, 0, len(keys))
	for _, key := range keys {
		key = strings.TrimSpace(key)
		if key != "" {
			normalizedKeys = append(normalizedKeys, key)
		}
	}

	a.mu.Lock()
	a.enabled = enabled
	a.header = normalizedHeader
	a.keys = normalizedKeys
	a.mu.Unlock()
}

// SetDynamicValidator sets dynamic API key validation callback.
func (a *APIKeyAuth) SetDynamicValidator(validator APIKeyValidator) {
	if a == nil {
		return
	}

	a.mu.Lock()
	a.dynamicValidator = validator
	a.mu.Unlock()
}

// Enabled reports whether API key authentication is active.
func (a *APIKeyAuth) Enabled() bool {
	if a == nil {
		return false
	}
	enabled, _, keys, validator := a.snapshot()
	return enabled && (len(keys) > 0 || validator != nil)
}

// Middleware authenticates requests when an API key header is present.
//
// Behavior:
// - Disabled or not configured: no-op
// - Missing header: no-op (session auth can still continue)
// - Invalid header: 401 Unauthorized
// - Valid header: marks request as API-key authenticated
func (a *APIKeyAuth) Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		enabled, header, keys, validator := a.snapshot()
		if !enabled || (len(keys) == 0 && validator == nil) {
			return c.Next()
		}

		incoming := strings.TrimSpace(c.Get(header))
		if incoming == "" {
			return c.Next()
		}

		valid := isValidAPIKey(incoming, keys)
		if !valid && validator != nil {
			valid = validator(c, incoming)
		}

		if !valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}

		c.Locals(APIKeyAuthenticatedLocalKey, true)
		return c.Next()
	}
}

func (a *APIKeyAuth) snapshot() (bool, string, []string, APIKeyValidator) {
	if a == nil {
		return false, "", nil, nil
	}

	a.mu.RLock()
	defer a.mu.RUnlock()

	keys := make([]string, len(a.keys))
	copy(keys, a.keys)

	return a.enabled, a.header, keys, a.dynamicValidator
}

func isValidAPIKey(incoming string, keys []string) bool {
	for _, key := range keys {
		if subtle.ConstantTimeCompare([]byte(incoming), []byte(key)) == 1 {
			return true
		}
	}
	return false
}
