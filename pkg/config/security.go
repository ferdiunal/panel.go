package config

import "time"

// SecurityConfig holds all security-related configuration
type SecurityConfig struct {
	// CORS configuration
	CORS CORSConfig

	// Rate limiting configuration
	RateLimit RateLimitConfig

	// Account lockout configuration
	AccountLockout AccountLockoutConfig

	// Session configuration
	Session SessionConfig

	// Encryption configuration
	Encryption EncryptionConfig

	// Audit logging configuration
	Audit AuditConfig
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	// AllowedOrigins is a list of allowed origins (e.g., ["https://example.com", "https://app.example.com"])
	// NEVER use "*" in production
	AllowedOrigins []string

	// AllowCredentials indicates whether credentials are allowed
	AllowCredentials bool

	// AllowedMethods lists allowed HTTP methods
	AllowedMethods []string

	// AllowedHeaders lists allowed headers
	AllowedHeaders []string

	// ExposeHeaders lists headers that can be exposed to the client
	ExposeHeaders []string

	// MaxAge indicates how long preflight results can be cached
	MaxAge int
}

// RateLimitConfig holds rate limiting configuration
type RateLimitConfig struct {
	// Enabled indicates whether rate limiting is enabled
	Enabled bool

	// AuthMaxRequests is the maximum number of auth requests per window
	AuthMaxRequests int

	// AuthWindow is the time window for auth rate limiting
	AuthWindow time.Duration

	// APIMaxRequests is the maximum number of API requests per window
	APIMaxRequests int

	// APIWindow is the time window for API rate limiting
	APIWindow time.Duration
}

// AccountLockoutConfig holds account lockout configuration
type AccountLockoutConfig struct {
	// Enabled indicates whether account lockout is enabled
	Enabled bool

	// MaxAttempts is the maximum number of failed login attempts before lockout
	MaxAttempts int

	// LockoutDuration is how long an account is locked after max attempts
	LockoutDuration time.Duration
}

// SessionConfig holds session security configuration
type SessionConfig struct {
	// CookieName is the name of the session cookie
	CookieName string

	// Secure indicates whether the cookie should only be sent over HTTPS
	Secure bool

	// HTTPOnly indicates whether the cookie should be HTTP-only
	HTTPOnly bool

	// SameSite sets the SameSite attribute ("Strict", "Lax", or "None")
	SameSite string

	// MaxAge is the maximum age of the session in seconds
	MaxAge int

	// Domain is the domain for the cookie
	Domain string

	// Path is the path for the cookie
	Path string
}

// EncryptionConfig holds encryption configuration
type EncryptionConfig struct {
	// Algorithm specifies the encryption algorithm ("AES-GCM" or "AES-CBC")
	Algorithm string

	// KeyHex is the hex-encoded encryption key
	KeyHex string

	// RotationEnabled indicates whether key rotation is enabled
	RotationEnabled bool

	// RotationInterval is how often to rotate keys
	RotationInterval time.Duration
}

// AuditConfig holds audit logging configuration
type AuditConfig struct {
	// Enabled indicates whether audit logging is enabled
	Enabled bool

	// LogLevel specifies what to log ("all", "security", "errors")
	LogLevel string

	// Destination specifies where to log ("console", "file", "siem")
	Destination string

	// FilePath is the path to the audit log file (if destination is "file")
	FilePath string

	// SIEMEndpoint is the SIEM endpoint URL (if destination is "siem")
	SIEMEndpoint string
}

// DefaultSecurityConfig returns secure default configuration
func DefaultSecurityConfig() SecurityConfig {
	return SecurityConfig{
		CORS: CORSConfig{
			AllowedOrigins:   []string{"http://localhost:3000"}, // Must be configured per environment
			AllowCredentials: true,
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Content-Type", "Authorization", "X-CSRF-Token"},
			ExposeHeaders:    []string{"Content-Length"},
			MaxAge:           3600,
		},
		RateLimit: RateLimitConfig{
			Enabled:         true,
			AuthMaxRequests: 10,
			AuthWindow:      1 * time.Minute,
			APIMaxRequests:  100,
			APIWindow:       1 * time.Minute,
		},
		AccountLockout: AccountLockoutConfig{
			Enabled:         true,
			MaxAttempts:     5,
			LockoutDuration: 15 * time.Minute,
		},
		Session: SessionConfig{
			CookieName: "__Host-session_token",
			Secure:     true,
			HTTPOnly:   true,
			SameSite:   "Strict",
			MaxAge:     86400, // 24 hours
			Domain:     "",
			Path:       "/",
		},
		Encryption: EncryptionConfig{
			Algorithm:        "AES-GCM",
			RotationEnabled:  false,
			RotationInterval: 90 * 24 * time.Hour, // 90 days
		},
		Audit: AuditConfig{
			Enabled:     true,
			LogLevel:    "security",
			Destination: "console",
		},
	}
}

// ProductionSecurityConfig returns production-ready security configuration
func ProductionSecurityConfig() SecurityConfig {
	config := DefaultSecurityConfig()

	// Production-specific overrides
	config.CORS.AllowedOrigins = []string{} // Must be explicitly configured
	config.Session.Secure = true
	config.Session.SameSite = "Strict"
	config.Audit.LogLevel = "all"
	config.Audit.Destination = "file"
	config.Audit.FilePath = "/var/log/panel/audit.log"

	return config
}

// DevelopmentSecurityConfig returns development-friendly security configuration
func DevelopmentSecurityConfig() SecurityConfig {
	config := DefaultSecurityConfig()

	// Development-specific overrides (still secure!)
	config.CORS.AllowedOrigins = []string{"http://localhost:3000", "http://localhost:5173"}
	config.Session.Secure = false // Allow HTTP in development
	config.Session.SameSite = "Lax"
	config.RateLimit.AuthMaxRequests = 50 // More lenient for development
	config.Audit.Destination = "console"

	return config
}
