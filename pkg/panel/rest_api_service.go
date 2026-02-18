package panel

import (
	"crypto/subtle"
	"os"
	"strings"

	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/gofiber/fiber/v2"
)

const (
	defaultInternalRESTAPIBasePath = "/api/internal/rest"
	defaultInternalRESTAPIHeader   = "X-Internal-API-Key"
)

type internalRESTAPIRuntimeConfig struct {
	Enabled  bool
	BasePath string
	Header   string
	Keys     []string
}

func (p *Panel) resolveInternalRESTAPIRuntimeConfig() internalRESTAPIRuntimeConfig {
	cfg := baseInternalRESTAPIRuntimeConfig(Config{})
	if p != nil {
		cfg = baseInternalRESTAPIRuntimeConfig(p.Config)
	}

	if p == nil {
		return cfg
	}

	cfg.Enabled = p.Config.Features.RestAPI
	if len(cfg.Keys) == 0 {
		cfg.Enabled = false
	}

	return cfg
}

func baseInternalRESTAPIRuntimeConfig(cfg Config) internalRESTAPIRuntimeConfig {
	keys := normalizeAPIKeyList(cfg.RESTAPI.Keys)
	if len(keys) == 0 {
		if raw := strings.TrimSpace(os.Getenv("INTERNAL_REST_API_KEYS")); raw != "" {
			keys = parseAPIKeyList(raw)
		}
	}
	if len(keys) == 0 {
		if raw := strings.TrimSpace(os.Getenv("INTERNAL_REST_API_KEY")); raw != "" {
			keys = parseAPIKeyList(raw)
		}
	}

	header := strings.TrimSpace(cfg.RESTAPI.Header)
	if header == "" {
		if raw := strings.TrimSpace(os.Getenv("INTERNAL_REST_API_HEADER")); raw != "" {
			header = raw
		}
	}
	if header == "" {
		header = defaultInternalRESTAPIHeader
	}

	basePath := strings.TrimSpace(cfg.RESTAPI.BasePath)
	if basePath == "" {
		basePath = strings.TrimSpace(os.Getenv("INTERNAL_REST_API_BASE_PATH"))
	}

	return internalRESTAPIRuntimeConfig{
		BasePath: normalizeInternalRESTAPIBasePath(basePath),
		Header:   header,
		Keys:     keys,
	}
}

func normalizeInternalRESTAPIBasePath(basePath string) string {
	basePath = strings.TrimSpace(basePath)
	if basePath == "" || basePath == "/" {
		return defaultInternalRESTAPIBasePath
	}

	if !strings.HasPrefix(basePath, "/") {
		basePath = "/" + basePath
	}

	basePath = strings.TrimRight(basePath, "/")
	if basePath == "" {
		return defaultInternalRESTAPIBasePath
	}

	return basePath
}

func (p *Panel) registerInternalRESTAPIRoutes(app *fiber.App) {
	if p == nil || app == nil {
		return
	}

	cfg := p.resolveInternalRESTAPIRuntimeConfig()
	if !cfg.Enabled {
		return
	}

	internalAPI := app.Group(cfg.BasePath)
	internalAPI.Use(func(c *fiber.Ctx) error {
		incoming := strings.TrimSpace(c.Get(cfg.Header))
		if incoming == "" || !containsAPIKey(cfg.Keys, incoming) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}
		return c.Next()
	})

	internalAPI.Get("/:resource", context.Wrap(p.handleResourceIndex))
	internalAPI.Get("/:resource/:id", context.Wrap(p.handleResourceDetail))
	internalAPI.Put("/:resource/:id", context.Wrap(p.handleResourceUpdate))
	internalAPI.Patch("/:resource/:id", context.Wrap(p.handleResourceUpdate))
	internalAPI.Delete("/:resource/:id", context.Wrap(p.handleResourceDestroy))
}

func containsAPIKey(keys []string, incoming string) bool {
	for _, key := range keys {
		if subtle.ConstantTimeCompare([]byte(incoming), []byte(key)) == 1 {
			return true
		}
	}
	return false
}

func isPanelAPIPath(path string, cfg Config) bool {
	if hasPathPrefix(path, "/api") {
		return true
	}

	basePath := baseInternalRESTAPIRuntimeConfig(cfg).BasePath
	if hasPathPrefix(path, basePath) {
		return true
	}

	externalBasePath := baseExternalAPIRuntimeConfig(cfg).BasePath
	return hasPathPrefix(path, externalBasePath)
}

func hasPathPrefix(path string, prefix string) bool {
	if prefix == "" {
		return false
	}
	if path == prefix {
		return true
	}
	return strings.HasPrefix(path, prefix+"/")
}
