// Package panel, HTML injection ve rendering işlemlerini sağlar.
//
// Bu paket, index.html dosyasını runtime'da modify ederek RTL, dark tema
// ve diğer dinamik bilgileri inject eder.
package panel

import (
	"fmt"
	"os"
	"strings"

	"github.com/ferdiunal/panel.go/pkg/rtl"
	"github.com/gofiber/fiber/v2"
)

// HTMLPlaceholders, index.html'de kullanılan placeholder'lar.
const (
	PlaceholderLang  = "{{PANEL_LANG}}"
	PlaceholderDir   = "{{PANEL_DIR}}"
	PlaceholderTheme = "{{PANEL_THEME}}"
	PlaceholderTitle = "{{PANEL_TITLE}}"
)

// HTMLInjectionData, HTML'e inject edilecek veri.
type HTMLInjectionData struct {
	Lang  string // Dil kodu (örn: "tr", "en", "ar")
	Dir   string // Text direction ("ltr" veya "rtl")
	Theme string // Tema ("light" veya "dark")
	Title string // Site başlığı
}

// GetHTMLInjectionData, request'ten HTML injection data'sını oluşturur.
//
// Bu fonksiyon, request context'inden dil, RTL ve tema bilgilerini alır
// ve HTMLInjectionData struct'ı oluşturur.
//
// ## Parametreler
//   - c: Fiber context
//   - config: Panel config
//
// ## Dönüş Değeri
//   - HTMLInjectionData: Injection data
func GetHTMLInjectionData(c *fiber.Ctx, config Config) HTMLInjectionData {
	// Dil bilgisini al
	lang := "en"

	// Config'deki default language'i kullan
	if config.I18n.Enabled && config.I18n.DefaultLanguage.String() != "" {
		lang = config.I18n.DefaultLanguage.String()
	}

	// Direction bilgisini al
	dir := rtl.GetDirectionString(lang)

	// Tema bilgisini al (cookie veya query'den)
	theme := c.Query("theme", c.Cookies("theme", "light"))
	if theme != "dark" && theme != "light" {
		theme = "light" // Varsayılan tema
	}

	// Site başlığını al
	title := config.SettingsValues.SiteName
	if title == "" {
		title = "Panel.go" // Varsayılan başlık
	}

	fmt.Println(lang, dir, theme, title)

	return HTMLInjectionData{
		Lang:  lang,
		Dir:   dir,
		Theme: theme,
		Title: title,
	}
}

// InjectHTML, HTML içeriğine placeholder'ları inject eder.
//
// Bu fonksiyon, HTML string'indeki placeholder'ları gerçek değerlerle
// replace eder.
//
// ## Parametreler
//   - html: HTML içeriği
//   - data: Injection data
//
// ## Dönüş Değeri
//   - string: Inject edilmiş HTML
func InjectHTML(html string, data HTMLInjectionData) string {
	// Placeholder'ları replace et
	html = strings.ReplaceAll(html, PlaceholderLang, data.Lang)
	html = strings.ReplaceAll(html, PlaceholderDir, data.Dir)
	html = strings.ReplaceAll(html, PlaceholderTheme, data.Theme)
	html = strings.ReplaceAll(html, PlaceholderTitle, data.Title)

	return html
}

// ServeHTML, HTML dosyasını inject ederek serve eder.
//
// Bu fonksiyon, HTML dosyasını okur, placeholder'ları inject eder ve
// client'a döndürür.
//
// ## Parametreler
//   - c: Fiber context
//   - htmlPath: HTML dosya yolu
//   - config: Panel config
//
// ## Dönüş Değeri
//   - error: Hata varsa error, yoksa nil
func ServeHTML(c *fiber.Ctx, htmlPath string, config Config) error {
	// HTML'i oku
	htmlBytes, err := os.ReadFile(htmlPath)
	if err != nil {
		return fmt.Errorf("HTML okunamadı: %w", err)
	}

	// Injection data'sını al
	data := GetHTMLInjectionData(c, config)

	// HTML'i inject et
	html := InjectHTML(string(htmlBytes), data)

	// HTML döndür
	c.Set("Content-Type", "text/html; charset=utf-8")
	c.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	return c.SendString(html)
}
