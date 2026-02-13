/// # Plugin Paketi
///
/// Bu paket, Panel.go için genişletilebilir plugin sistemi sağlar.
/// Plugin'ler compile-time'da yüklenir ve runtime'da boot edilir.
///
/// ## Temel Özellikler
/// - Compile-time plugin loading (type-safe)
/// - Resource, Page, Field, Middleware, Route ekleme
/// - Lifecycle management (Register, Boot)
/// - Thread-safe registry
///
/// ## Kullanım Örneği
/// ```go
/// type MyPlugin struct{}
///
/// func (p *MyPlugin) Name() string { return "my-plugin" }
/// func (p *MyPlugin) Version() string { return "1.0.0" }
/// func (p *MyPlugin) Author() string { return "Author Name" }
/// func (p *MyPlugin) Description() string { return "Plugin description" }
///
/// func (p *MyPlugin) Register(panel *panel.Panel) error {
///     // Plugin kayıt işlemleri
///     return nil
/// }
///
/// func (p *MyPlugin) Boot(panel *panel.Panel) error {
///     // Plugin başlatma işlemleri
///     return nil
/// }
///
/// func init() {
///     plugin.Register(&MyPlugin{})
/// }
/// ```

package plugin

import (
	"github.com/ferdiunal/panel.go/pkg/resource"
	"github.com/gofiber/fiber/v2"
)

/// # Plugin Interface
///
/// Plugin, Panel.go için genişletilebilir plugin sistemi interface'idir.
/// Her plugin bu interface'i implement etmelidir.
///
/// ## Lifecycle
/// 1. **Register**: Plugin registry'ye kaydedilir (init() fonksiyonunda)
/// 2. **Boot**: Panel başlatıldığında plugin boot edilir
///
/// ## Capabilities
/// Plugin'ler aşağıdaki özellikleri sağlayabilir:
/// - Resources: Yeni resource'lar ekler
/// - Pages: Yeni sayfalar ekler
/// - Middleware: HTTP middleware'ler ekler
/// - Routes: Özel API endpoint'leri ekler
/// - Migrations: Veritabanı migration'ları ekler
///
/// ## Önemli Notlar
/// - Plugin'ler thread-safe olmalıdır
/// - Register() metodu sadece bir kez çağrılır
/// - Boot() metodu Panel başlatıldığında çağrılır
/// - Hata durumunda error döndürülmelidir
type Plugin interface {
	// Metadata - Plugin bilgileri
	Name() string        // Plugin adı (örn: "analytics-plugin")
	Version() string     // Semantic versioning (örn: "1.0.0")
	Author() string      // Plugin yazarı
	Description() string // Kısa açıklama

	// Lifecycle - Plugin yaşam döngüsü
	// Register: Plugin registry'ye kaydedildiğinde çağrılır
	// Bu metod plugin'in temel yapılandırmasını yapar
	Register(panel interface{}) error

	// Boot: Panel başlatıldığında çağrılır
	// Bu metod plugin'in resource, page, middleware vb. eklemelerini yapar
	Boot(panel interface{}) error

	// Capabilities - Plugin yetenekleri (opsiyonel)
	// Bu metodlar nil dönebilir (plugin o özelliği sağlamıyorsa)

	// Resources: Plugin'in sağladığı resource'lar
	Resources() []resource.Resource

	// Pages: Plugin'in sağladığı sayfalar
	Pages() []Page

	// Middleware: Plugin'in sağladığı HTTP middleware'ler
	Middleware() []fiber.Handler

	// Routes: Plugin'in sağladığı özel route'lar
	// Router üzerinde özel endpoint'ler tanımlanabilir
	Routes(router fiber.Router)

	// Migrations: Plugin'in sağladığı veritabanı migration'ları
	Migrations() []Migration
}

/// # Page Interface
///
/// Plugin'lerin sağlayabileceği sayfa interface'i.
/// pkg/page.Page interface'i ile uyumlu olmalıdır.
type Page interface {
	Slug() string
	Title() string
	Icon() string
	Group() string
	Visible() bool
	NavigationOrder() int
}

/// # Migration Interface
///
/// Plugin'lerin sağlayabileceği migration interface'i.
/// Veritabanı şema değişiklikleri için kullanılır.
///
/// ## Kullanım Örneği
/// ```go
/// type CreateUsersTable struct{}
///
/// func (m *CreateUsersTable) Name() string {
///     return "create_users_table"
/// }
///
/// func (m *CreateUsersTable) Up(db *gorm.DB) error {
///     return db.AutoMigrate(&User{})
/// }
///
/// func (m *CreateUsersTable) Down(db *gorm.DB) error {
///     return db.Migrator().DropTable(&User{})
/// }
/// ```
type Migration interface {
	// Name: Migration adı (benzersiz olmalı)
	Name() string

	// Up: Migration'ı uygula
	Up(db interface{}) error

	// Down: Migration'ı geri al
	Down(db interface{}) error
}

/// # BasePlugin Struct
///
/// Plugin interface'ini implement etmek için temel struct.
/// Plugin geliştiriciler bu struct'ı embed ederek sadece ihtiyaç duydukları
/// metodları override edebilirler.
///
/// ## Kullanım Örneği
/// ```go
/// type MyPlugin struct {
///     plugin.BasePlugin
/// }
///
/// func (p *MyPlugin) Name() string { return "my-plugin" }
/// func (p *MyPlugin) Version() string { return "1.0.0" }
/// func (p *MyPlugin) Author() string { return "Author" }
/// func (p *MyPlugin) Description() string { return "Description" }
///
/// func (p *MyPlugin) Resources() []resource.Resource {
///     return []resource.Resource{&MyResource{}}
/// }
/// ```
type BasePlugin struct{}

// Metadata defaults
func (p *BasePlugin) Name() string        { return "unnamed-plugin" }
func (p *BasePlugin) Version() string     { return "0.0.0" }
func (p *BasePlugin) Author() string      { return "Unknown" }
func (p *BasePlugin) Description() string { return "" }

// Lifecycle defaults
func (p *BasePlugin) Register(panel interface{}) error { return nil }
func (p *BasePlugin) Boot(panel interface{}) error     { return nil }

// Capabilities defaults (nil = özellik yok)
func (p *BasePlugin) Resources() []resource.Resource { return nil }
func (p *BasePlugin) Pages() []Page                  { return nil }
func (p *BasePlugin) Middleware() []fiber.Handler    { return nil }
func (p *BasePlugin) Routes(router fiber.Router)     {}
func (p *BasePlugin) Migrations() []Migration        { return nil }
