package panel

import (
	"net/http/httptest"
	"strings"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestBuildContentSecurityPolicyDevelopmentDefaults(t *testing.T) {
	t.Setenv("POS_ECHO_HOST", "")
	t.Setenv("POS_ECHO_PORT", "")
	t.Setenv("SOKETI_HOST", "")
	t.Setenv("SOKETI_PORT", "")
	t.Setenv("VITE_POS_ECHO_HOST", "")
	t.Setenv("VITE_POS_ECHO_PORT", "")
	t.Setenv("VITE_SOKETI_HOST", "")
	t.Setenv("VITE_SOKETI_PORT", "")

	csp := buildContentSecurityPolicy("development")

	if !strings.Contains(csp, "ws://127.0.0.1:6001") {
		t.Fatalf("expected development CSP to allow ws localhost fallback, got %q", csp)
	}
	if !strings.Contains(csp, "wss://127.0.0.1:6001") {
		t.Fatalf("expected development CSP to allow wss localhost fallback, got %q", csp)
	}
	if !strings.Contains(csp, "ws://localhost:6001") {
		t.Fatalf("expected development CSP to allow ws localhost fallback, got %q", csp)
	}
	if !strings.Contains(csp, "wss://localhost:6001") {
		t.Fatalf("expected development CSP to allow wss localhost fallback, got %q", csp)
	}
	if strings.Contains(csp, "connect-src *") {
		t.Fatalf("expected CSP to avoid wildcard connect-src, got %q", csp)
	}
}

func TestBuildContentSecurityPolicyUsesEchoEnvOrigins(t *testing.T) {
	t.Setenv("POS_ECHO_HOST", "127.0.0.1")
	t.Setenv("POS_ECHO_PORT", "6002")
	t.Setenv("SOKETI_HOST", "")
	t.Setenv("SOKETI_PORT", "")
	t.Setenv("VITE_POS_ECHO_HOST", "")
	t.Setenv("VITE_POS_ECHO_PORT", "")
	t.Setenv("VITE_SOKETI_HOST", "")
	t.Setenv("VITE_SOKETI_PORT", "")

	csp := buildContentSecurityPolicy("production")

	if !strings.Contains(csp, "ws://127.0.0.1:6002") {
		t.Fatalf("expected production CSP to include configured ws origin, got %q", csp)
	}
	if !strings.Contains(csp, "wss://127.0.0.1:6002") {
		t.Fatalf("expected production CSP to include configured wss origin, got %q", csp)
	}
	if strings.Contains(csp, "ws://localhost:6001") {
		t.Fatalf("production CSP must not include development localhost defaults, got %q", csp)
	}
}

func TestBuildContentSecurityPolicyUsesViteFallbackEnv(t *testing.T) {
	t.Setenv("POS_ECHO_HOST", "")
	t.Setenv("POS_ECHO_PORT", "")
	t.Setenv("SOKETI_HOST", "")
	t.Setenv("SOKETI_PORT", "")
	t.Setenv("VITE_POS_ECHO_HOST", "http://echo.local")
	t.Setenv("VITE_POS_ECHO_PORT", "7001")
	t.Setenv("VITE_SOKETI_HOST", "")
	t.Setenv("VITE_SOKETI_PORT", "")

	csp := buildContentSecurityPolicy("production")

	if !strings.Contains(csp, "ws://echo.local:7001") {
		t.Fatalf("expected CSP to include VITE echo ws origin, got %q", csp)
	}
	if !strings.Contains(csp, "wss://echo.local:7001") {
		t.Fatalf("expected CSP to include VITE echo wss origin, got %q", csp)
	}
}

func TestSecurityHeadersIncludeDynamicRealtimeOrigins(t *testing.T) {
	t.Setenv("POS_ECHO_HOST", "127.0.0.1")
	t.Setenv("POS_ECHO_PORT", "6001")
	t.Setenv("SOKETI_HOST", "")
	t.Setenv("SOKETI_PORT", "")
	t.Setenv("VITE_POS_ECHO_HOST", "")
	t.Setenv("VITE_POS_ECHO_PORT", "")
	t.Setenv("VITE_SOKETI_HOST", "")
	t.Setenv("VITE_SOKETI_PORT", "")

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}

	p := New(Config{
		Database:    DatabaseConfig{Instance: db},
		Environment: "test",
	})

	req := httptest.NewRequest("GET", "/api/internal/init", nil)
	resp, err := testFiberRequest(p.Fiber, req)
	if err != nil {
		t.Fatalf("init request failed: %v", err)
	}
	if resp.StatusCode != 200 {
		t.Fatalf("expected status 200 from /api/internal/init, got %d", resp.StatusCode)
	}

	csp := resp.Header.Get("Content-Security-Policy")
	if !strings.Contains(csp, "ws://127.0.0.1:6001") {
		t.Fatalf("expected csp to allow ws realtime origin, got %q", csp)
	}
	if !strings.Contains(csp, "wss://127.0.0.1:6001") {
		t.Fatalf("expected csp to allow wss realtime origin, got %q", csp)
	}

	if got := resp.Header.Get("X-Frame-Options"); got != "DENY" {
		t.Fatalf("expected X-Frame-Options=DENY, got %q", got)
	}
	if got := resp.Header.Get("X-Content-Type-Options"); got != "nosniff" {
		t.Fatalf("expected X-Content-Type-Options=nosniff, got %q", got)
	}
	if got := resp.Header.Get("Referrer-Policy"); got != "no-referrer" {
		t.Fatalf("expected Referrer-Policy=no-referrer, got %q", got)
	}
	if got := resp.Header.Get("Permissions-Policy"); got != "geolocation=(), microphone=(), camera=()" {
		t.Fatalf("unexpected Permissions-Policy header: %q", got)
	}
}
