/// # Example Plugin
///
/// Panel.go için örnek plugin implementasyonu.
/// Plugin sisteminin tüm özelliklerini gösterir.
///
/// ## Özellikler
/// - Custom resource ekleme
/// - Custom field ekleme
/// - Custom middleware ekleme
/// - Custom route ekleme
/// - Migration ekleme
///
/// ## Kullanım Örneği
/// ```go
/// // main.go
/// import _ "github.com/ferdiunal/panel.go/plugins/example-plugin"
///
/// func main() {
///     config := panel.Config{...}
///     p := panel.New(config)
///     p.Start()
/// }
/// ```

package example_plugin

import (
	"fmt"

	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/plugin"
	"github.com/ferdiunal/panel.go/pkg/resource"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// / # ExamplePlugin Struct
// /
// / Örnek plugin implementasyonu.
// / plugin.BasePlugin'i embed ederek sadece gerekli metodları override eder.
type ExamplePlugin struct {
	plugin.BasePlugin
}

// / # Metadata Metodları
// /
// / Plugin'in temel bilgilerini döndürür.
func (p *ExamplePlugin) Name() string        { return "example-plugin" }
func (p *ExamplePlugin) Version() string     { return "1.0.0" }
func (p *ExamplePlugin) Author() string      { return "Panel.go Team" }
func (p *ExamplePlugin) Description() string { return "Example plugin demonstrating plugin system" }

// / # Register Metodu
// /
// / Plugin kaydedildiğinde çağrılır.
// / Temel yapılandırma işlemleri burada yapılır.
func (p *ExamplePlugin) Register(panel interface{}) error {
	fmt.Println("ExamplePlugin: Register called")
	return nil
}

// / # Boot Metodu
// /
// / Panel başlatıldığında çağrılır.
// / Resource, page, middleware vb. ekleme işlemleri burada yapılır.
func (p *ExamplePlugin) Boot(panel interface{}) error {
	fmt.Println("ExamplePlugin: Boot called")
	return nil
}

// / # Resources Metodu
// /
// / Plugin'in sağladığı resource'ları döndürür.
func (p *ExamplePlugin) Resources() []resource.Resource {
	return []resource.Resource{
		NewExampleResource(),
	}
}

// / # Middleware Metodu
// /
// / Plugin'in sağladığı middleware'leri döndürür.
func (p *ExamplePlugin) Middleware() []fiber.Handler {
	return []fiber.Handler{
		func(c *fiber.Ctx) error {
			// Example middleware: Log request
			fmt.Printf("ExamplePlugin Middleware: %s %s\n", c.Method(), c.Path())
			return c.Next()
		},
	}
}

// / # Routes Metodu
// /
// / Plugin'in sağladığı özel route'ları tanımlar.
func (p *ExamplePlugin) Routes(router fiber.Router) {
	router.Get("/api/example-plugin/hello", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Hello from ExamplePlugin!",
			"version": p.Version(),
		})
	})
}

// / # Migrations Metodu
// /
// / Plugin'in sağladığı migration'ları döndürür.
func (p *ExamplePlugin) Migrations() []plugin.Migration {
	return []plugin.Migration{
		&CreateExampleTable{},
	}
}

// / # ExampleModel Struct
// /
// / Örnek veritabanı modeli.
type ExampleModel struct {
	gorm.Model
	Name        string `json:"name" gorm:"size:255;not null"`
	Description string `json:"description" gorm:"type:text"`
	Active      bool   `json:"active" gorm:"default:true"`
}

// ExampleResource is a modern resource definition for the example plugin.
type ExampleResource struct {
	resource.Base
}

// / # CreateExampleTable Migration
// /
// / Example tablosunu oluşturan migration.
type CreateExampleTable struct{}

func (m *CreateExampleTable) Name() string {
	return "create_example_table"
}

func (m *CreateExampleTable) Up(db interface{}) error {
	gormDB, ok := db.(*gorm.DB)
	if !ok {
		return fmt.Errorf("invalid database instance")
	}
	return gormDB.AutoMigrate(&ExampleModel{})
}

func (m *CreateExampleTable) Down(db interface{}) error {
	gormDB, ok := db.(*gorm.DB)
	if !ok {
		return fmt.Errorf("invalid database instance")
	}
	return gormDB.Migrator().DropTable(&ExampleModel{})
}

// / # NewExampleResource Fonksiyonu
// /
// / Örnek resource oluşturur.
func NewExampleResource() resource.Resource {
	return &ExampleResource{
		Base: resource.Base{
			DataModel:  &ExampleModel{},
			Identifier: "example",
			Label:      "Example Resource",
			IconName:   "cube",
			GroupName:  "Plugin Examples",
			FieldsVal: []fields.Element{
				fields.ID(),
				fields.Text("Name", "name").
					Label("Name").
					Placeholder("Enter name").
					Required().
					Sortable().
					Searchable(),
				fields.Textarea("Description", "description").
					Label("Description").
					Placeholder("Enter description").
					Rows(3),
				fields.Switch("Active", "active").
					Label("Active").
					HelpText("Enable or disable this item").
					Default(true),
				fields.DateTime("Created At", "created_at").
					Label("Created At").
					Sortable().
					OnlyOnList(),
				fields.DateTime("Updated At", "updated_at").
					Label("Updated At").
					Sortable().
					OnlyOnList(),
			},
		},
	}
}

// / # init Fonksiyonu
// /
// / Plugin'i global registry'ye kaydeder.
// / Bu fonksiyon otomatik olarak çağrılır (Go'nun init mekanizması).
func init() {
	plugin.Register(&ExamplePlugin{})
	fmt.Println("ExamplePlugin registered via init()")
}
