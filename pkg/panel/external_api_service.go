package panel

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/gofiber/fiber/v2"
)

const (
	defaultExternalAPIBasePath = "/external-api"
	defaultExternalAPIHeader   = "X-External-API-Key"
)

type externalAPIRuntimeConfig struct {
	Enabled  bool
	BasePath string
	Header   string
	Keys     []string
}

func (p *Panel) resolveExternalAPIRuntimeConfig() externalAPIRuntimeConfig {
	cfg := baseExternalAPIRuntimeConfig(Config{})
	if p != nil {
		cfg = baseExternalAPIRuntimeConfig(p.Config)
	}

	if p == nil {
		return cfg
	}

	cfg.Enabled = p.Config.Features.ExternalAPI
	if len(cfg.Keys) == 0 {
		cfg.Enabled = false
	}

	return cfg
}

func baseExternalAPIRuntimeConfig(cfg Config) externalAPIRuntimeConfig {
	keys := normalizeAPIKeyList(cfg.ExternalAPI.Keys)
	if len(keys) == 0 {
		if raw := strings.TrimSpace(os.Getenv("EXTERNAL_API_KEYS")); raw != "" {
			keys = parseAPIKeyList(raw)
		}
	}
	if len(keys) == 0 {
		if raw := strings.TrimSpace(os.Getenv("EXTERNAL_API_KEY")); raw != "" {
			keys = parseAPIKeyList(raw)
		}
	}

	header := strings.TrimSpace(cfg.ExternalAPI.Header)
	if header == "" {
		if raw := strings.TrimSpace(os.Getenv("EXTERNAL_API_HEADER")); raw != "" {
			header = raw
		}
	}
	if header == "" {
		header = defaultExternalAPIHeader
	}

	basePath := strings.TrimSpace(cfg.ExternalAPI.BasePath)
	if basePath == "" {
		basePath = strings.TrimSpace(os.Getenv("EXTERNAL_API_BASE_PATH"))
	}

	return externalAPIRuntimeConfig{
		BasePath: normalizeExternalAPIBasePath(basePath),
		Header:   header,
		Keys:     keys,
	}
}

func normalizeExternalAPIBasePath(basePath string) string {
	basePath = strings.TrimSpace(basePath)
	if basePath == "" || basePath == "/" {
		return defaultExternalAPIBasePath
	}

	if !strings.HasPrefix(basePath, "/") {
		basePath = "/" + basePath
	}

	basePath = strings.TrimRight(basePath, "/")
	if basePath == "" {
		return defaultExternalAPIBasePath
	}

	return basePath
}

func (p *Panel) registerExternalAPIRoutes(app *fiber.App) {
	if p == nil || app == nil {
		return
	}

	cfg := p.resolveExternalAPIRuntimeConfig()
	if !cfg.Enabled {
		return
	}

	externalAPI := app.Group(cfg.BasePath)
	externalAPI.Use(func(c *fiber.Ctx) error {
		incoming := strings.TrimSpace(c.Get(cfg.Header))
		if incoming == "" || !containsAPIKey(cfg.Keys, incoming) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized",
			})
		}
		return c.Next()
	})

	externalAPI.Get("/:resource", context.Wrap(p.handleExternalResourceIndex))
	externalAPI.Get("/:resource/:id", context.Wrap(p.handleExternalResourceShow))
	externalAPI.Post("/:resource", context.Wrap(p.handleExternalResourceStore))
	externalAPI.Put("/:resource/:id", context.Wrap(p.handleExternalResourceUpdate))
	externalAPI.Patch("/:resource/:id", context.Wrap(p.handleExternalResourceUpdate))
	externalAPI.Delete("/:resource/:id", context.Wrap(p.handleExternalResourceDestroy))
}

func (p *Panel) handleExternalResourceIndex(c *context.Context) error {
	if err := p.handleResourceIndex(c); err != nil {
		return err
	}
	flattenExternalResourceResponse(c)
	return nil
}

func (p *Panel) handleExternalResourceShow(c *context.Context) error {
	if err := p.handleResourceShow(c); err != nil {
		return err
	}
	flattenExternalResourceResponse(c)
	return nil
}

func (p *Panel) handleExternalResourceStore(c *context.Context) error {
	if err := p.handleResourceStore(c); err != nil {
		return err
	}
	flattenExternalResourceResponse(c)
	return nil
}

func (p *Panel) handleExternalResourceUpdate(c *context.Context) error {
	if err := p.handleResourceUpdate(c); err != nil {
		return err
	}
	flattenExternalResourceResponse(c)
	return nil
}

func (p *Panel) handleExternalResourceDestroy(c *context.Context) error {
	return p.handleResourceDestroy(c)
}

func flattenExternalResourceResponse(c *context.Context) {
	if c == nil || c.Ctx == nil {
		return
	}

	resp := c.Ctx.Response()
	if resp == nil {
		return
	}
	if resp.StatusCode() >= fiber.StatusBadRequest {
		return
	}

	rawBody := resp.Body()
	if len(rawBody) == 0 {
		return
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(rawBody, &payload); err != nil {
		return
	}

	changed := false

	if rawData, ok := payload["data"]; ok {
		if flatData, ok := flattenExternalData(rawData); ok {
			payload["data"] = flatData
			changed = true
		}
	}

	if rawFields, ok := payload["fields"]; ok {
		if flatFields, ok := flattenExternalFieldList(rawFields); ok {
			payload["data"] = flatFields
			delete(payload, "fields")
			changed = true
		}
	}

	if !changed {
		return
	}

	encoded, err := json.Marshal(payload)
	if err != nil {
		return
	}

	resp.SetBody(encoded)
	resp.Header.SetContentType(fiber.MIMEApplicationJSONCharsetUTF8)
}

func flattenExternalData(raw interface{}) (interface{}, bool) {
	switch typed := raw.(type) {
	case map[string]interface{}:
		return flattenExternalRecord(typed), true
	case []interface{}:
		out := make([]interface{}, 0, len(typed))
		for _, entry := range typed {
			record, ok := entry.(map[string]interface{})
			if !ok {
				out = append(out, entry)
				continue
			}
			out = append(out, flattenExternalRecord(record))
		}
		return out, true
	default:
		return nil, false
	}
}

func flattenExternalFieldList(raw interface{}) (map[string]interface{}, bool) {
	items, ok := raw.([]interface{})
	if !ok {
		return nil, false
	}

	out := make(map[string]interface{}, len(items))
	for _, item := range items {
		fieldMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		if isHiddenExternalFieldPayload(fieldMap) {
			continue
		}

		key, _ := fieldMap["key"].(string)
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}

		if data, ok := fieldMap["data"]; ok {
			out[key] = data
			continue
		}

		out[key] = fieldMap
	}

	return out, true
}

func flattenExternalRecord(record map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(record))
	for key, value := range record {
		fieldPayload, ok := value.(map[string]interface{})
		if ok && isHiddenExternalFieldPayload(fieldPayload) {
			continue
		}
		out[key] = extractExternalFieldValue(value)
	}
	return out
}

func extractExternalFieldValue(value interface{}) interface{} {
	fieldPayload, ok := value.(map[string]interface{})
	if !ok {
		return value
	}
	if !isExternalFieldPayload(fieldPayload) {
		return value
	}

	if data, ok := fieldPayload["data"]; ok {
		return data
	}

	if data, ok := fieldPayload["value"]; ok {
		return data
	}

	return value
}

func isExternalFieldPayload(payload map[string]interface{}) bool {
	if payload == nil {
		return false
	}

	_, hasKey := payload["key"]
	_, hasView := payload["view"]
	_, hasType := payload["type"]
	return hasKey || hasView || hasType
}

func isHiddenExternalFieldPayload(payload map[string]interface{}) bool {
	if !isExternalFieldPayload(payload) {
		return false
	}

	rawContext, ok := payload["context"].(string)
	if !ok {
		return false
	}

	for _, ctx := range strings.Fields(rawContext) {
		if ctx == string(fields.HIDE_ON_API) {
			return true
		}
	}

	return false
}
