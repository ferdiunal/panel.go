package panel

import (
	"github.com/ferdiunal/panel.go/internal/page"
	"github.com/ferdiunal/panel.go/internal/resource"
	"gorm.io/gorm"
)

// FeatureConfig, panelin opsiyonel özelliklerini açıp kapatmak için kullanılır.
type FeatureConfig struct {
	// Register, kullanıcı kayıt özelliğini aktif eder.
	Register bool
	// ForgotPassword, şifre sırlama özelliğini aktif eder.
	ForgotPassword bool
}

// OAuthConfig, OAuth sağlayıcılarının yapılandırmasını tutar.
type OAuthConfig struct {
	Google GoogleConfig
}

// GoogleConfig, Google OAuth yapılandırmasını içerir.
type GoogleConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

// Enabled, Google OAuth konfigürasyonunun geçerli olup olmadığını kontrol eder.
func (c GoogleConfig) Enabled() bool {
	return c.ClientID != "" && c.ClientSecret != ""
}

// Config, panelin genel yapılandırma yapısıdır.
type Config struct {
	Server         ServerConfig
	Database       DatabaseConfig
	Environment    string // "production", "development", "test"
	Features       FeatureConfig
	OAuth          OAuthConfig
	Storage        StorageConfig
	SettingsValues SettingsConfig
	SettingsPage   *page.Settings
	DashboardPage  *page.Dashboard
	UserResource   resource.Resource
}

// SettingsConfig, veritabanından gelen dinamik ayarları tutar.
type SettingsConfig struct {
	SiteName       string
	Register       bool
	ForgotPassword bool
	Values         map[string]interface{}
}

// StorageConfig, dosya yükleme ve depolama ayarlarını tutar.
type StorageConfig struct {
	// Path, dosyaların fiziksel olarak saklanacağı sunucu yolu.
	Path string
	// URL, dosyalara erişim için kullanılacak temel URL yolu.
	URL string
}

// ServerConfig, HTTP sunucu ayarlarını tutar.
type ServerConfig struct {
	Port string
	Host string
}

// DatabaseConfig, veritabanı bağlantı ayarlarını tutar.
type DatabaseConfig struct {
	DSN      string
	Driver   string   // "postgres", "mysql", "sqlite"
	Instance *gorm.DB // Mevcut bir GORM bağlantısı varsa (Opsiyonel)
}
