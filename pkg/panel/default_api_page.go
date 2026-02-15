package panel

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	appContext "github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/page"
	"gorm.io/gorm"
)

type defaultAPIPage struct {
	page.OptimizedBase
	defaults apiKeyRuntimeConfig
}

func newDefaultAPIPage(cfg Config) page.Page {
	p := &defaultAPIPage{
		defaults: baseAPIKeyRuntimeConfig(cfg.APIKey),
	}

	p.SetSlug("api-settings")
	p.SetTitle("API")
	p.SetDescription("API endpoints, documentation, and API key settings")
	p.SetIcon("code")
	p.SetGroup("System")
	p.SetNavigationOrder(2)
	p.SetVisible(true)

	return p
}

func (p *defaultAPIPage) Fields() []fields.Element {
	return p.buildFields(p.defaults)
}

func (p *defaultAPIPage) GetFieldsWithContext(ctx *appContext.Context) []fields.Element {
	cfg := p.defaults
	if ctx != nil && ctx.DB() != nil {
		if values, err := readAPIKeySettings(ctx.DB()); err == nil {
			cfg = applyAPIKeySettings(cfg, values)
		}
	}
	return p.buildFields(cfg)
}

func (p *defaultAPIPage) buildFields(cfg apiKeyRuntimeConfig) []fields.Element {
	return []fields.Element{
		fields.Text("OpenAPI Spec", "openapi_spec").Default("/api/openapi.json").ReadOnly(),
		fields.Text("Swagger UI", "swagger_ui").Default("/api/docs").ReadOnly(),
	}
}

func (p *defaultAPIPage) Save(c *appContext.Context, db *gorm.DB, data map[string]any) error {
	if c == nil || !p.CanAccess(c) {
		return errors.New("access denied")
	}
	if db == nil {
		return errors.New("database is nil")
	}

	cfg := p.defaults
	if values, err := readAPIKeySettings(db); err == nil {
		cfg = applyAPIKeySettings(cfg, values)
	}

	enabled := cfg.Enabled
	if raw, ok := data[apiKeySettingEnabled]; ok {
		enabled = parseBoolValue(raw)
	}

	header := cfg.Header
	if raw, ok := data[apiKeySettingHeader]; ok {
		header = strings.TrimSpace(fmt.Sprintf("%v", raw))
	}
	if header == "" {
		header = defaultAPIKeyHeader
	}

	keys := cfg.Keys
	if raw, ok := data[apiKeySettingKeys]; ok {
		keys = parseAPIKeyInput(raw)
	}
	keysJSON, err := json.Marshal(keys)
	if err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if err := upsertSettingValue(
			tx,
			apiKeySettingEnabled,
			fmt.Sprintf("%t", enabled),
			"boolean",
			"api",
			"API Key Authentication",
			"Enable API key authentication for API endpoints.",
		); err != nil {
			return err
		}

		if err := upsertSettingValue(
			tx,
			apiKeySettingHeader,
			header,
			"string",
			"api",
			"API Key Header",
			"Header name that clients must send with each request.",
		); err != nil {
			return err
		}

		if err := upsertSettingValue(
			tx,
			apiKeySettingKeys,
			string(keysJSON),
			"json",
			"api",
			"API Keys",
			"Allowed API keys in JSON array format.",
		); err != nil {
			return err
		}

		return nil
	})
}

func (p *defaultAPIPage) CanAccess(ctx *appContext.Context) bool {
	if ctx == nil {
		return false
	}

	currentUser := ctx.User()
	if currentUser == nil {
		return false
	}

	return strings.EqualFold(currentUser.Role, "admin")
}
