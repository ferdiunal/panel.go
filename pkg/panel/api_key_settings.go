package panel

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ferdiunal/panel.go/pkg/domain/setting"
	"gorm.io/gorm"
)

const (
	apiKeySettingEnabled = "api_key_enabled"
	apiKeySettingHeader  = "api_key_header"
	apiKeySettingKeys    = "api_keys"
	defaultAPIKeyHeader  = "X-API-Key"
)

type apiKeyRuntimeConfig struct {
	Enabled bool
	Header  string
	Keys    []string
}

func baseAPIKeyRuntimeConfig(cfg APIKeyConfig) apiKeyRuntimeConfig {
	keys := normalizeAPIKeyList(cfg.Keys)
	if len(keys) == 0 {
		if envKey := strings.TrimSpace(os.Getenv("PANEL_API_KEY")); envKey != "" {
			keys = []string{envKey}
		}
	}

	header := strings.TrimSpace(cfg.Header)
	if header == "" {
		header = defaultAPIKeyHeader
	}

	enabled := cfg.Enabled
	if !enabled && len(keys) > 0 {
		enabled = true
	}

	return apiKeyRuntimeConfig{
		Enabled: enabled,
		Header:  header,
		Keys:    keys,
	}
}

func applyAPIKeySettings(base apiKeyRuntimeConfig, values map[string]string) apiKeyRuntimeConfig {
	cfg := base

	_, hasEnabled := values[apiKeySettingEnabled]
	if hasEnabled {
		cfg.Enabled = parseBoolValue(values[apiKeySettingEnabled])
	}

	if raw, ok := values[apiKeySettingHeader]; ok {
		header := strings.TrimSpace(raw)
		if header == "" {
			header = defaultAPIKeyHeader
		}
		cfg.Header = header
	}

	if raw, ok := values[apiKeySettingKeys]; ok {
		cfg.Keys = parseAPIKeyList(raw)
		if !hasEnabled && len(cfg.Keys) > 0 {
			cfg.Enabled = true
		}
	}

	return cfg
}

func readAPIKeySettings(db *gorm.DB) (map[string]string, error) {
	if db == nil {
		return nil, fmt.Errorf("database is nil")
	}

	if !db.Migrator().HasTable(&setting.Setting{}) {
		return map[string]string{}, nil
	}

	var settingsRows []setting.Setting
	if err := db.Where(
		"key IN ?",
		[]string{apiKeySettingEnabled, apiKeySettingHeader, apiKeySettingKeys},
	).Find(&settingsRows).Error; err != nil {
		return nil, err
	}

	values := make(map[string]string, len(settingsRows))
	for _, row := range settingsRows {
		values[row.Key] = row.Value
	}

	return values, nil
}

func (p *Panel) resolveAPIKeyRuntimeConfig() apiKeyRuntimeConfig {
	if p == nil {
		return baseAPIKeyRuntimeConfig(APIKeyConfig{})
	}

	base := baseAPIKeyRuntimeConfig(p.Config.APIKey)
	if p.Db == nil {
		return base
	}

	values, err := readAPIKeySettings(p.Db)
	if err != nil {
		return base
	}

	return applyAPIKeySettings(base, values)
}

func (p *Panel) refreshAPIKeyAuthFromSettings() {
	if p == nil || p.apiKeyAuth == nil {
		return
	}

	cfg := p.resolveAPIKeyRuntimeConfig()
	p.apiKeyAuth.SetConfig(cfg.Enabled, cfg.Header, cfg.Keys)
}

func parseAPIKeyList(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return []string{}
	}

	if strings.HasPrefix(raw, "[") {
		var asJSON []string
		if err := json.Unmarshal([]byte(raw), &asJSON); err == nil {
			return normalizeAPIKeyList(asJSON)
		}
	}

	parts := strings.FieldsFunc(raw, func(r rune) bool {
		return r == '\n' || r == '\r' || r == ',' || r == ';' || r == '\t'
	})
	return normalizeAPIKeyList(parts)
}

func parseAPIKeyInput(value any) []string {
	switch v := value.(type) {
	case nil:
		return []string{}
	case string:
		return parseAPIKeyList(v)
	case []string:
		return normalizeAPIKeyList(v)
	case []any:
		values := make([]string, 0, len(v))
		for _, entry := range v {
			values = append(values, fmt.Sprintf("%v", entry))
		}
		return normalizeAPIKeyList(values)
	default:
		return parseAPIKeyList(fmt.Sprintf("%v", value))
	}
}

func normalizeAPIKeyList(keys []string) []string {
	normalized := make([]string, 0, len(keys))
	seen := make(map[string]struct{}, len(keys))

	for _, key := range keys {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		normalized = append(normalized, key)
	}

	return normalized
}

func upsertSettingValue(db *gorm.DB, key, value, valueType, group, label, help string) error {
	if db == nil {
		return fmt.Errorf("database is nil")
	}

	var existing setting.Setting
	err := db.Where("key = ?", key).First(&existing).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		return db.Create(&setting.Setting{
			Key:   key,
			Value: value,
			Type:  valueType,
			Group: group,
			Label: label,
			Help:  help,
		}).Error
	}

	existing.Value = value
	existing.Type = valueType
	existing.Group = group
	existing.Label = label
	existing.Help = help
	return db.Save(&existing).Error
}

func isAPIKeySettingKey(key string) bool {
	switch key {
	case apiKeySettingEnabled, apiKeySettingHeader, apiKeySettingKeys:
		return true
	default:
		return false
	}
}
