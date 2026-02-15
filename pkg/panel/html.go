// Package panel, HTML injection ve rendering işlemlerini sağlar.
//
// Bu paket, index.html dosyasını runtime'da modify ederek RTL, dark tema
// ve diğer dinamik bilgileri inject eder.
package panel

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/ferdiunal/panel.go/pkg/i18n"
	"github.com/ferdiunal/panel.go/pkg/rtl"
	"github.com/gofiber/fiber/v2"
	"gopkg.in/yaml.v3"
)

// embeddedLocalesFS, SDK ile gelen varsayılan locale dosyalarını içerir.
// Kullanıcı locale dosyaları bunların üzerine override edilir.
//
//go:embed locales/*.yaml
var embeddedLocalesFS embed.FS

// HTMLPlaceholders, index.html'de kullanılan placeholder'lar.
const (
	PlaceholderLang    = "{{PANEL_LANG}}"
	PlaceholderDir     = "{{PANEL_DIR}}"
	PlaceholderTitle   = "{{PANEL_TITLE}}"
	PlaceholderInit    = "{{PANEL_INIT}}"
	PlaceholderFavicon = "{{PANEL_FAVICON}}"
)

// HTMLInjectionData, HTML'e inject edilecek veri.
type HTMLInjectionData struct {
	Lang    string // Dil kodu (örn: "tr", "en", "ar")
	Dir     string // Text direction ("ltr" veya "rtl")
	Title   string // Site başlığı
	Favicon string // Favicon URL (logo veya varsayılan)
}

// GetHTMLInjectionData, request'ten HTML injection data'sını oluşturur.
func GetHTMLInjectionData(c *fiber.Ctx, config Config) HTMLInjectionData {
	// Dil bilgisini al
	lang := "en"

	if config.I18n.Enabled {
		lang = i18n.GetLocale(c)
		if lang == "" || (lang == "en" && config.I18n.DefaultLanguage.String() != "" && config.I18n.DefaultLanguage.String() != "en") {
			lang = config.I18n.DefaultLanguage.String()
		}
	} else if config.I18n.DefaultLanguage.String() != "" {
		lang = config.I18n.DefaultLanguage.String()
	}

	// Direction bilgisini al
	dir := rtl.GetDirectionString(lang)

	// Site başlığını al
	title := config.SettingsValues.SiteName
	if title == "" {
		title = "Panel.go"
	}

	// Favicon — settings'den logo varsa onu kullan
	favicon := "/vite.svg"
	if config.SettingsValues.Values != nil {
		if logo, ok := config.SettingsValues.Values["logo"]; ok {
			if logoStr, ok := logo.(string); ok && logoStr != "" {
				favicon = logoStr
			}
		}
	}

	return HTMLInjectionData{
		Lang:    lang,
		Dir:     dir,
		Title:   title,
		Favicon: favicon,
	}
}

// loadTranslations, belirtilen dil için YAML locale dosyasını okur
// ve tüm key'leri "dot notation" (nested.key.path) olarak flat bir map'e çevirir.
func loadTranslations(config Config, lang string) map[string]interface{} {
	merged := make(map[string]interface{})
	mergeTranslationMaps(merged, loadEmbeddedTranslations(lang))
	mergeTranslationMaps(merged, loadUserTranslations(config, lang))
	return merged
}

func loadEmbeddedTranslations(lang string) map[string]interface{} {
	translations := make(map[string]interface{})

	for _, fileName := range localeFileCandidates(lang) {
		filePath := path.Join("locales", fileName)
		data, err := fs.ReadFile(embeddedLocalesFS, filePath)
		if err != nil {
			continue
		}

		flat, err := parseAndFlattenLocaleYAML(data)
		if err != nil {
			continue
		}

		mergeTranslationMaps(translations, flat)
	}

	return translations
}

func loadUserTranslations(config Config, lang string) map[string]interface{} {
	translations := make(map[string]interface{})

	rootPath := config.I18n.RootPath
	if rootPath == "" {
		rootPath = "./locales"
	}

	for _, fileName := range localeFileCandidates(lang) {
		filePath := filepath.Join(rootPath, fileName)
		data, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		flat, err := parseAndFlattenLocaleYAML(data)
		if err != nil {
			continue
		}

		mergeTranslationMaps(translations, flat)
	}

	return translations
}

func parseAndFlattenLocaleYAML(data []byte) (map[string]interface{}, error) {
	var raw map[string]interface{}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	flat := make(map[string]interface{})
	flattenMap("", raw, flat)
	return flat, nil
}

func mergeTranslationMaps(dest map[string]interface{}, src map[string]interface{}) {
	for key, value := range src {
		dest[key] = value
	}
}

func localeFileCandidates(lang string) []string {
	normalized := strings.TrimSpace(strings.ReplaceAll(lang, "_", "-"))
	if normalized == "" {
		return nil
	}

	normalized = strings.ToLower(normalized)
	base := normalized
	if idx := strings.Index(normalized, "-"); idx > 0 {
		base = normalized[:idx]
	}

	// Base -> Exact sırası:
	// tr.yaml + tr-tr.yaml gibi dosyalar varsa exact locale base'i override eder.
	candidates := []string{base + ".yaml"}
	if normalized != base {
		candidates = append(candidates, normalized+".yaml")
	}

	return candidates
}

// flattenMap, nested map'i dot notation'a çevirir.
// Örn: { auth: { login: { title: "Giriş" } } } -> { "auth.login.title": "Giriş" }
func flattenMap(prefix string, src map[string]interface{}, dest map[string]interface{}) {
	for key, val := range src {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		switch v := val.(type) {
		case map[string]interface{}:
			flattenMap(fullKey, v, dest)
		default:
			dest[fullKey] = v
		}
	}
}

// GetInitData, frontend'e inject edilecek init verisini oluşturur.
// handleInit API endpoint'i ile aynı veriyi döner.
func GetInitData(c *fiber.Ctx, config Config, injectionData HTMLInjectionData) fiber.Map {
	// Features
	registerEnabled := config.Features.Register
	forgotPasswordEnabled := config.Features.ForgotPassword

	if settings := config.SettingsValues.Values; settings != nil {
		if registerVal, ok := settings["registration_enabled"]; ok {
			registerEnabled = parseBoolValue(registerVal)
		}
		if forgotVal, ok := settings["forgot_password_enabled"]; ok {
			forgotPasswordEnabled = parseBoolValue(forgotVal)
		}
	}

	// i18n
	i18nData := fiber.Map{
		"lang":      injectionData.Lang,
		"direction": injectionData.Dir,
	}

	if config.I18n.Enabled {
		supportedLangs := []fiber.Map{}
		for _, lang := range config.I18n.AcceptLanguages {
			supportedLangs = append(supportedLangs, fiber.Map{
				"code": lang.String(),
				"name": getLanguageName(lang),
			})
		}
		i18nData["supported_languages"] = supportedLangs
		i18nData["default_language"] = config.I18n.DefaultLanguage.String()
		i18nData["use_url_prefix"] = config.I18n.UseURLPrefix
		i18nData["url_prefix_optional"] = config.I18n.URLPrefixOptional
	}

	// Translations — aktif dil için locale dosyasını oku
	translations := loadTranslations(config, injectionData.Lang)

	return fiber.Map{
		"features": fiber.Map{
			"register":        registerEnabled,
			"forgot_password": forgotPasswordEnabled,
		},
		"oauth": fiber.Map{
			"google": config.OAuth.Google.Enabled(),
		},
		"i18n":         i18nData,
		"translations": translations,
		"version":      "1.0.0",
		"settings":     config.SettingsValues.Values,
	}
}

// InjectHTML, HTML içeriğine placeholder'ları inject eder.
func InjectHTML(html string, data HTMLInjectionData, initJSON string) string {
	html = strings.ReplaceAll(html, PlaceholderLang, data.Lang)
	html = strings.ReplaceAll(html, PlaceholderDir, data.Dir)
	html = strings.ReplaceAll(html, PlaceholderTitle, data.Title)
	html = strings.ReplaceAll(html, PlaceholderFavicon, data.Favicon)
	html = strings.ReplaceAll(html, PlaceholderInit, initJSON)

	return html
}

// ServeHTML, HTML dosyasını inject ederek serve eder.
func ServeHTML(c *fiber.Ctx, htmlPath string, config Config) error {
	// HTML'i oku
	htmlBytes, err := os.ReadFile(htmlPath)
	if err != nil {
		return fmt.Errorf("HTML okunamadı: %w", err)
	}

	// Injection data'sını al
	data := GetHTMLInjectionData(c, config)

	// Init verisini JSON olarak oluştur
	initData := GetInitData(c, config, data)
	initJSON, err := json.Marshal(initData)
	if err != nil {
		initJSON = []byte("{}")
	}

	// HTML'i inject et
	html := InjectHTML(string(htmlBytes), data, string(initJSON))

	// HTML döndür
	c.Set("Content-Type", "text/html; charset=utf-8")
	c.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	return c.SendString(html)
}
