// Package plugin, Panel.go plugin sistemi için scaffold işlemlerini sağlar.
//
// Bu paket, plugin backend ve frontend dosyalarını oluşturur, workspace
// config'i günceller ve plugin metadata'sını okur.
package plugin

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"
	"unicode/utf8"

	"gopkg.in/yaml.v3"
)

// PluginMetadata, plugin.yaml metadata yapısı.
type PluginMetadata struct {
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Author      string `yaml:"author"`
	Description string `yaml:"description"`
}

// generateBackendFiles, backend dosyalarını oluşturur.
//
// Bu fonksiyon, plugin.go ve plugin.yaml dosyalarını oluşturur.
//
// ## Parametreler
//   - pluginDir: Plugin dizini
//   - name: Plugin adı
//   - withExample: Tüm relationship türlerini içeren örnek entity'ler ekle
//
// ## Dönüş Değeri
//   - error: Oluşturma hatası varsa hata, aksi takdirde nil
func generateBackendFiles(pluginDir, name string, withExample bool) error {
	if withExample {
		// Example template kullan - tüm relationship türlerini içeren comprehensive örnek
		if err := generateEntityFile(pluginDir, name); err != nil {
			return err
		}
		if err := generateResourceFiles(pluginDir, name); err != nil {
			return err
		}
		if err := generateRegistryFile(pluginDir, name); err != nil {
			return err
		}
		if err := generateExampleMainFile(pluginDir, name); err != nil {
			return err
		}

		// plugin.yaml oluştur
		pluginYAMLPath := filepath.Join(pluginDir, "plugin.yaml")
		pluginYAMLContent := fmt.Sprintf(pluginYAMLTemplate, name, name)
		if err := os.WriteFile(pluginYAMLPath, []byte(pluginYAMLContent), 0644); err != nil {
			return fmt.Errorf("plugin.yaml oluşturulamadı: %w", err)
		}

		return nil
	}

	// Basit template kullan (mevcut kod)
	// plugin.go oluştur
	packageName := strings.ReplaceAll(name, "-", "_")
	pluginGoPath := filepath.Join(pluginDir, "plugin.go")
	pluginGoContent := fmt.Sprintf(`package %s

import (
	"github.com/ferdiunal/panel.go/pkg/plugin"
)

// Plugin, %s plugin'i.
type Plugin struct {
	plugin.BasePlugin
}

// init, plugin'i global registry'ye kaydeder.
func init() {
	plugin.Register(&Plugin{})
}

// Name, plugin adını döndürür.
func (p *Plugin) Name() string {
	return "%s"
}

// Version, plugin versiyonunu döndürür.
func (p *Plugin) Version() string {
	return "1.0.0"
}

// Register, plugin'i Panel'e kaydeder.
func (p *Plugin) Register(panel interface{}) error {
	// Plugin registration logic
	// Örnek: Resource, Page, Middleware ekle
	return nil
}

// Boot, plugin'i boot eder.
func (p *Plugin) Boot(panel interface{}) error {
	// Plugin boot logic
	// Örnek: Database migration, event listener ekle
	return nil
}
	`, packageName, name, name)

	if err := os.WriteFile(pluginGoPath, []byte(pluginGoContent), 0644); err != nil {
		return fmt.Errorf("plugin.go oluşturulamadı: %w", err)
	}

	// plugin.yaml oluştur
	pluginYAMLPath := filepath.Join(pluginDir, "plugin.yaml")
	pluginYAMLContent := fmt.Sprintf(`name: %s
version: 1.0.0
author: Panel.go Team
description: %s plugin for Panel.go
`, name, name)

	if err := os.WriteFile(pluginYAMLPath, []byte(pluginYAMLContent), 0644); err != nil {
		return fmt.Errorf("plugin.yaml oluşturulamadı: %w", err)
	}

	return nil
}

// generateFrontendFiles, frontend dosyalarını oluşturur.
//
// Bu fonksiyon, index.ts, package.json ve tsconfig.json dosyalarını oluşturur.
//
// ## Parametreler
//   - pluginDir: Plugin dizini
//   - name: Plugin adı
//
// ## Dönüş Değeri
//   - error: Oluşturma hatası varsa hata, aksi takdirde nil
func generateFrontendFiles(pluginDir, name string) error {
	// frontend dizini oluştur
	frontendDir := filepath.Join(pluginDir, "frontend")
	if err := os.MkdirAll(frontendDir, 0755); err != nil {
		return fmt.Errorf("frontend dizini oluşturulamadı: %w", err)
	}

	// index.ts oluştur
	indexTSPath := filepath.Join(frontendDir, "index.ts")
	indexTSContent := fmt.Sprintf(`import { Plugin } from '@/plugins/types';

export const %sPlugin: Plugin = {
  name: '%s',
  version: '1.0.0',
  description: '%s plugin',
  author: 'Panel.go Team',

  fields: [
    // Custom field'lar buraya eklenecek
    // Örnek:
    // {
    //   type: 'custom-field',
    //   component: CustomField,
    // },
  ],
};

export default %sPlugin;
`, toPascalCase(name), name, name, toPascalCase(name))

	if err := os.WriteFile(indexTSPath, []byte(indexTSContent), 0644); err != nil {
		return fmt.Errorf("index.ts oluşturulamadı: %w", err)
	}

	// package.json oluştur
	packageJSONPath := filepath.Join(frontendDir, "package.json")
	packageJSONContent := fmt.Sprintf(`{
  "name": "@panel-plugins/%s",
  "version": "1.0.0",
  "private": true,
  "type": "module",
  "main": "index.ts",
  "dependencies": {
    "react": "workspace:*",
    "react-dom": "workspace:*"
  }
}
`, name)

	if err := os.WriteFile(packageJSONPath, []byte(packageJSONContent), 0644); err != nil {
		return fmt.Errorf("package.json oluşturulamadı: %w", err)
	}

	// tsconfig.json oluştur
	tsconfigJSONPath := filepath.Join(frontendDir, "tsconfig.json")
	tsconfigJSONContent := `{
  "extends": "../../tsconfig.json",
  "compilerOptions": {
    "composite": true,
    "baseUrl": ".",
    "paths": {
      "@/*": ["../../src/*"]
    }
  },
  "include": ["**/*.ts", "**/*.tsx"],
  "exclude": ["node_modules"]
}
`

	if err := os.WriteFile(tsconfigJSONPath, []byte(tsconfigJSONContent), 0644); err != nil {
		return fmt.Errorf("tsconfig.json oluşturulamadı: %w", err)
	}

	// fields dizini oluştur
	fieldsDir := filepath.Join(frontendDir, "fields")
	if err := os.MkdirAll(fieldsDir, 0755); err != nil {
		return fmt.Errorf("fields dizini oluşturulamadı: %w", err)
	}

	// .gitkeep oluştur
	gitkeepPath := filepath.Join(fieldsDir, ".gitkeep")
	if err := os.WriteFile(gitkeepPath, []byte(""), 0644); err != nil {
		return fmt.Errorf(".gitkeep oluşturulamadı: %w", err)
	}

	return nil
}

// updateWorkspaceConfig, workspace config'i günceller.
//
// Bu fonksiyon, pnpm-workspace.yaml dosyasını günceller ve plugin'i ekler.
//
// ## Parametreler
//   - webUIPath: web-ui dizini
//   - name: Plugin adı
//   - pluginPath: Plugin dizini
//
// ## Dönüş Değeri
//   - error: Güncelleme hatası varsa hata, aksi takdirde nil
func updateWorkspaceConfig(webUIPath, name, pluginPath string) error {
	workspaceYAMLPath := filepath.Join(webUIPath, "pnpm-workspace.yaml")

	// Workspace config var mı kontrol et
	var workspaceConfig map[string]interface{}
	if _, err := os.Stat(workspaceYAMLPath); err == nil {
		// Mevcut config'i oku
		data, err := os.ReadFile(workspaceYAMLPath)
		if err != nil {
			return fmt.Errorf("workspace config okunamadı: %w", err)
		}

		if err := yaml.Unmarshal(data, &workspaceConfig); err != nil {
			return fmt.Errorf("workspace config parse edilemedi: %w", err)
		}
	} else {
		// Yeni config oluştur
		workspaceConfig = map[string]interface{}{
			"packages": []string{},
		}
	}

	// packages listesini al
	packages, ok := workspaceConfig["packages"].([]interface{})
	if !ok {
		packages = []interface{}{}
	}

	// Plugin path'i ekle (eğer yoksa)
	pluginWorkspacePath := "../plugins/*/frontend"
	found := false
	for _, pkg := range packages {
		if pkgStr, ok := pkg.(string); ok && pkgStr == pluginWorkspacePath {
			found = true
			break
		}
	}

	if !found {
		packages = append(packages, pluginWorkspacePath)
		workspaceConfig["packages"] = packages
	}

	// Config'i yaz
	data, err := yaml.Marshal(workspaceConfig)
	if err != nil {
		return fmt.Errorf("workspace config marshal edilemedi: %w", err)
	}

	if err := os.WriteFile(workspaceYAMLPath, data, 0644); err != nil {
		return fmt.Errorf("workspace config yazılamadı: %w", err)
	}

	return nil
}

// createPluginSymlink, plugin workspace reference oluşturur.
//
// Bu fonksiyon, web-ui/plugins/<name> -> ../../plugins/<name>/frontend
// symlink'i oluşturur.
//
// ## Parametreler
//   - webUIPath: web-ui dizini
//   - name: Plugin adı
//   - pluginPath: Plugin dizini
//
// ## Dönüş Değeri
//   - error: Oluşturma hatası varsa hata, aksi takdirde nil
func createPluginSymlink(webUIPath, name, pluginPath string) error {
	// web-ui/plugins dizini oluştur
	pluginsDir := filepath.Join(webUIPath, "plugins")
	if err := os.MkdirAll(pluginsDir, 0755); err != nil {
		return fmt.Errorf("plugins dizini oluşturulamadı: %w", err)
	}

	// Symlink path
	symlinkPath := filepath.Join(pluginsDir, name)

	// Symlink zaten var mı kontrol et
	if _, err := os.Lstat(symlinkPath); err == nil {
		// Mevcut symlink'i sil
		if err := os.Remove(symlinkPath); err != nil {
			return fmt.Errorf("mevcut symlink silinemedi: %w", err)
		}
	}

	// Target path (relative)
	// web-ui/plugins/<name> -> ../../plugins/<name>/frontend
	targetPath := filepath.Join("..", "..", pluginPath, "frontend")

	// Symlink oluştur
	if err := os.Symlink(targetPath, symlinkPath); err != nil {
		return fmt.Errorf("symlink oluşturulamadı: %w", err)
	}

	return nil
}

// readPluginMetadata, plugin metadata'sını okur.
//
// Bu fonksiyon, plugin.yaml dosyasını okur ve metadata'yı döndürür.
//
// ## Parametreler
//   - pluginDir: Plugin dizini
//
// ## Dönüş Değeri
//   - *PluginMetadata: Plugin metadata
//   - error: Okuma hatası varsa hata, aksi takdirde nil
func readPluginMetadata(pluginDir string) (*PluginMetadata, error) {
	pluginYAMLPath := filepath.Join(pluginDir, "plugin.yaml")

	// plugin.yaml var mı kontrol et
	if _, err := os.Stat(pluginYAMLPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("plugin.yaml bulunamadı")
	}

	// plugin.yaml oku
	data, err := os.ReadFile(pluginYAMLPath)
	if err != nil {
		return nil, fmt.Errorf("plugin.yaml okunamadı: %w", err)
	}

	// Parse et
	var metadata PluginMetadata
	if err := yaml.Unmarshal(data, &metadata); err != nil {
		return nil, fmt.Errorf("plugin.yaml parse edilemedi: %w", err)
	}

	return &metadata, nil
}

// toPascalCase, string'i PascalCase'e çevirir.
//
// Bu fonksiyon, kebab-case veya snake_case string'i PascalCase'e çevirir.
//
// ## Parametreler
//   - s: String
//
// ## Dönüş Değeri
//   - string: PascalCase string
func toPascalCase(s string) string {
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == '-' || r == '_' || unicode.IsSpace(r)
	})

	result := strings.Builder{}
	for _, part := range parts {
		if len(part) > 0 {
			r, size := utf8.DecodeRuneInString(part)
			if r == utf8.RuneError && size == 0 {
				continue
			}
			result.WriteRune(unicode.ToUpper(r))
			result.WriteString(part[size:])
		}
	}

	return result.String()
}

// generateEntityFile, entity/entity.go dosyasını oluşturur.
func generateEntityFile(pluginDir, name string) error {
	entityDir := filepath.Join(pluginDir, "entity")
	if err := os.MkdirAll(entityDir, 0755); err != nil {
		return fmt.Errorf("entity dizini oluşturulamadı: %w", err)
	}

	entityPath := filepath.Join(entityDir, "entity.go")
	// entityContent := fmt.Sprintf(exampleEntityTemplate, name)
	return os.WriteFile(entityPath, []byte(exampleEntityTemplate), 0644)
}

// generateResourceFiles, resources/ altında her entity için resource dosyası oluşturur.
func generateResourceFiles(pluginDir, name string) error {
	resourcesDir := filepath.Join(pluginDir, "resources")
	if err := os.MkdirAll(resourcesDir, 0755); err != nil {
		return fmt.Errorf("resources dizini oluşturulamadı: %w", err)
	}

	// Her entity için resource dosyası oluştur
	resources := map[string]string{
		"organization": organizationResourceTemplate,
		"billing_info": billingInfoResourceTemplate,
		"address":      addressResourceTemplate,
		"product":      productResourceTemplate,
		"category":     categoryResourceTemplate,
		"shipment":     shipmentResourceTemplate,
		"shipment_row": shipmentRowResourceTemplate,
		// comment ve tag resource'ları kaldırıldı (MorphTo/MorphToMany field'ları henüz implement edilmemiş)
	}

	for resName, template := range resources {
		resDir := filepath.Join(resourcesDir, resName)
		if err := os.MkdirAll(resDir, 0755); err != nil {
			return fmt.Errorf("%s dizini oluşturulamadı: %w", resName, err)
		}

		resPath := filepath.Join(resDir, "main.go")
		resContent := fmt.Sprintf(template, name, name)
		if err := os.WriteFile(resPath, []byte(resContent), 0644); err != nil {
			return fmt.Errorf("%s/main.go oluşturulamadı: %w", resName, err)
		}
	}

	return nil
}

// generateRegistryFile, resources/registry.go dosyasını oluşturur.
func generateRegistryFile(pluginDir, name string) error {
	resourcesDir := filepath.Join(pluginDir, "resources")
	registryPath := filepath.Join(resourcesDir, "registry.go")
	// registryContent := fmt.Sprintf(registryTemplate, name)
	return os.WriteFile(registryPath, []byte(registryTemplate), 0644)
}

// generateExampleMainFile, plugin.go dosyasını oluşturur.
func generateExampleMainFile(pluginDir, name string) error {
	mainPath := filepath.Join(pluginDir, "plugin.go")
	// Package adını sanitize et (tire'leri underscore'a çevir)
	packageName := strings.ReplaceAll(name, "-", "_")
	importBase := name
	if modulePath, err := detectModulePath(); err == nil {
		if relPath, err := filepath.Rel(".", pluginDir); err == nil && !strings.HasPrefix(relPath, "..") {
			importBase = filepath.ToSlash(filepath.Join(modulePath, relPath))
		}
	}
	mainContent := fmt.Sprintf(exampleMainTemplate, packageName, importBase, name)
	return os.WriteFile(mainPath, []byte(mainContent), 0644)
}

func detectModulePath() (string, error) {
	data, err := os.ReadFile("go.mod")
	if err != nil {
		return "", fmt.Errorf("go.mod okunamadı: %w", err)
	}

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "module ") {
			modulePath := strings.TrimSpace(strings.TrimPrefix(line, "module"))
			if modulePath != "" {
				return modulePath, nil
			}
		}
	}

	return "", fmt.Errorf("go.mod içinde module satırı bulunamadı")
}

// Template constant'ları

const pluginYAMLTemplate = `name: %s
version: 1.0.0
author: Panel.go Team
description: %s plugin for Panel.go
`

const exampleEntityTemplate = `package entity

import (
	"time"
)

// Organization - Organizasyon bilgilerini tutan model
// HasMany: Addresses, Products, Shipments
// HasOne: BillingInfo
type Organization struct {
	ID          uint64       ` + "`" + `gorm:"primaryKey;autoIncrement;column:id;bigint" json:"id"` + "`" + `
	Name        string       ` + "`" + `gorm:"not null;column:name;varchar(255)" json:"name"` + "`" + `
	Email       string       ` + "`" + `gorm:"not null;column:email;varchar(255)" json:"email"` + "`" + `
	Phone       string       ` + "`" + `gorm:"not null;column:phone;varchar(255)" json:"phone"` + "`" + `
	Addresses   []Address    ` + "`" + `gorm:"foreignKey:OrganizationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"addresses,omitempty"` + "`" + `
	BillingInfo *BillingInfo ` + "`" + `gorm:"foreignKey:OrganizationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"billingInfo,omitempty"` + "`" + `
	Products    []Product    ` + "`" + `gorm:"foreignKey:OrganizationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"products,omitempty"` + "`" + `
	Shipments   []Shipment   ` + "`" + `gorm:"foreignKey:OrganizationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"shipments,omitempty"` + "`" + `
	CreatedAt   time.Time    ` + "`" + `gorm:"autoCreateTime;column:created_at;timestamptz" json:"createdAt"` + "`" + `
	UpdatedAt   time.Time    ` + "`" + `gorm:"autoUpdateTime;column:updated_at;timestamptz" json:"updatedAt"` + "`" + `
}

// BillingInfo - Fatura bilgilerini tutan model (HasOne örneği)
// BelongsTo: Organization
type BillingInfo struct {
	ID             uint64        ` + "`" + `gorm:"primaryKey;autoIncrement;column:id;bigint" json:"id"` + "`" + `
	OrganizationID uint64        ` + "`" + `gorm:"not null;unique;column:organization_id;bigint;index" json:"organizationId"` + "`" + `
	Organization   *Organization ` + "`" + `gorm:"foreignKey:OrganizationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"organization,omitempty"` + "`" + `
	TaxNumber      string        ` + "`" + `gorm:"not null;column:tax_number;varchar(255)" json:"taxNumber"` + "`" + `
	TaxOffice      string        ` + "`" + `gorm:"not null;column:tax_office;varchar(255)" json:"taxOffice"` + "`" + `
	CreatedAt      time.Time     ` + "`" + `gorm:"autoCreateTime;column:created_at;timestamptz" json:"createdAt"` + "`" + `
	UpdatedAt      time.Time     ` + "`" + `gorm:"autoUpdateTime;column:updated_at;timestamptz" json:"updatedAt"` + "`" + `
}

// Address - Adres bilgilerini tutan model
// BelongsTo: Organization
type Address struct {
	ID             uint64        ` + "`" + `gorm:"primaryKey;autoIncrement;column:id;bigint" json:"id"` + "`" + `
	OrganizationID uint64        ` + "`" + `gorm:"not null;column:organization_id;bigint;index" json:"organizationId"` + "`" + `
	Organization   *Organization ` + "`" + `gorm:"foreignKey:OrganizationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"organization,omitempty"` + "`" + `
	Name           string        ` + "`" + `gorm:"not null;column:name;varchar(255)" json:"name"` + "`" + `
	Address        string        ` + "`" + `gorm:"not null;column:address;text" json:"address"` + "`" + `
	City           string        ` + "`" + `gorm:"not null;column:city;varchar(255);index" json:"city"` + "`" + `
	Country        string        ` + "`" + `gorm:"not null;column:country;varchar(255);index" json:"country"` + "`" + `
	CreatedAt      time.Time     ` + "`" + `gorm:"autoCreateTime;column:created_at;timestamptz" json:"createdAt"` + "`" + `
	UpdatedAt      time.Time     ` + "`" + `gorm:"autoUpdateTime;column:updated_at;timestamptz" json:"updatedAt"` + "`" + `
}

// Product - Ürün bilgilerini tutan model
// BelongsTo: Organization
// BelongsToMany: Categories
// HasMany: ShipmentRows
type Product struct {
	ID             uint64        ` + "`" + `gorm:"primaryKey;autoIncrement;column:id;bigint" json:"id"` + "`" + `
	OrganizationID uint64        ` + "`" + `gorm:"not null;column:organization_id;bigint;index" json:"organizationId"` + "`" + `
	Organization   *Organization ` + "`" + `gorm:"foreignKey:OrganizationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"organization,omitempty"` + "`" + `
	Name           string        ` + "`" + `gorm:"not null;column:name;varchar(255);index" json:"name"` + "`" + `
	Categories     []Category    ` + "`" + `gorm:"many2many:product_categories;" json:"categories,omitempty"` + "`" + `
	ShipmentRows   []ShipmentRow ` + "`" + `gorm:"foreignKey:ProductID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"shipmentRows,omitempty"` + "`" + `
	CreatedAt      time.Time     ` + "`" + `gorm:"autoCreateTime;column:created_at;timestamptz" json:"createdAt"` + "`" + `
	UpdatedAt      time.Time     ` + "`" + `gorm:"autoUpdateTime;column:updated_at;timestamptz" json:"updatedAt"` + "`" + `
}

// Category - Kategori bilgilerini tutan model (BelongsToMany örneği)
// BelongsToMany: Products
type Category struct {
	ID        uint64    ` + "`" + `gorm:"primaryKey;autoIncrement;column:id;bigint" json:"id"` + "`" + `
	Name      string    ` + "`" + `gorm:"not null;column:name;varchar(255);index" json:"name"` + "`" + `
	Products  []Product ` + "`" + `gorm:"many2many:product_categories;" json:"products,omitempty"` + "`" + `
	CreatedAt time.Time ` + "`" + `gorm:"autoCreateTime;column:created_at;timestamptz" json:"createdAt"` + "`" + `
	UpdatedAt time.Time ` + "`" + `gorm:"autoUpdateTime;column:updated_at;timestamptz" json:"updatedAt"` + "`" + `
}

// Shipment - Gönderi bilgilerini tutan model
// BelongsTo: Organization
// HasMany: ShipmentRows
type Shipment struct {
	ID             uint64        ` + "`" + `gorm:"primaryKey;autoIncrement;column:id;bigint" json:"id"` + "`" + `
	OrganizationID uint64        ` + "`" + `gorm:"not null;column:organization_id;bigint;index" json:"organizationId"` + "`" + `
	Organization   *Organization ` + "`" + `gorm:"foreignKey:OrganizationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"organization,omitempty"` + "`" + `
	Name           string        ` + "`" + `gorm:"not null;column:name;varchar(255)" json:"name"` + "`" + `
	ShipmentRows   []ShipmentRow ` + "`" + `gorm:"foreignKey:ShipmentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"shipmentRows,omitempty"` + "`" + `
	CreatedAt      time.Time     ` + "`" + `gorm:"autoCreateTime;column:created_at;timestamptz" json:"createdAt"` + "`" + `
	UpdatedAt      time.Time     ` + "`" + `gorm:"autoUpdateTime;column:updated_at;timestamptz" json:"updatedAt"` + "`" + `
}

// ShipmentRow - Gönderi satırı bilgilerini tutan model
// BelongsTo: Shipment, Product
type ShipmentRow struct {
	ID         uint64    ` + "`" + `gorm:"primaryKey;autoIncrement;column:id;bigint" json:"id"` + "`" + `
	ShipmentID uint64    ` + "`" + `gorm:"not null;column:shipment_id;bigint;index" json:"shipmentId"` + "`" + `
	Shipment   *Shipment ` + "`" + `gorm:"foreignKey:ShipmentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"shipment,omitempty"` + "`" + `
	ProductID  uint64    ` + "`" + `gorm:"not null;column:product_id;bigint;index" json:"productId"` + "`" + `
	Product    *Product  ` + "`" + `gorm:"foreignKey:ProductID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"product,omitempty"` + "`" + `
	Quantity   int       ` + "`" + `gorm:"not null;column:quantity;int" json:"quantity"` + "`" + `
	CreatedAt  time.Time ` + "`" + `gorm:"autoCreateTime;column:created_at;timestamptz" json:"createdAt"` + "`" + `
	UpdatedAt  time.Time ` + "`" + `gorm:"autoUpdateTime;column:updated_at;timestamptz" json:"updatedAt"` + "`" + `
}

// Comment - Yorum bilgilerini tutan model (MorphTo örneği)
// MorphTo: Commentable (Product/Shipment)
type Comment struct {
	ID              uint64    ` + "`" + `gorm:"primaryKey;autoIncrement;column:id;bigint" json:"id"` + "`" + `
	CommentableID   uint64    ` + "`" + `gorm:"not null;column:commentable_id;bigint;index:idx_commentable" json:"commentableId"` + "`" + `
	CommentableType string    ` + "`" + `gorm:"not null;column:commentable_type;varchar(255);index:idx_commentable" json:"commentableType"` + "`" + `
	Content         string    ` + "`" + `gorm:"not null;column:content;text" json:"content"` + "`" + `
	CreatedAt       time.Time ` + "`" + `gorm:"autoCreateTime;column:created_at;timestamptz" json:"createdAt"` + "`" + `
	UpdatedAt       time.Time ` + "`" + `gorm:"autoUpdateTime;column:updated_at;timestamptz" json:"updatedAt"` + "`" + `
}

// Tag - Etiket bilgilerini tutan model (MorphToMany örneği)
// MorphToMany: Taggable (Product/Shipment)
type Tag struct {
	ID        uint64    ` + "`" + `gorm:"primaryKey;autoIncrement;column:id;bigint" json:"id"` + "`" + `
	Name      string    ` + "`" + `gorm:"not null;unique;column:name;varchar(255);index" json:"name"` + "`" + `
	CreatedAt time.Time ` + "`" + `gorm:"autoCreateTime;column:created_at;timestamptz" json:"createdAt"` + "`" + `
	UpdatedAt time.Time ` + "`" + `gorm:"autoUpdateTime;column:updated_at;timestamptz" json:"updatedAt"` + "`" + `
}
`

const registryTemplate = `package resources

import (
	"sync"

	"github.com/ferdiunal/panel.go/pkg/resource"
)

// registry - Resource factory'lerini saklayan merkezi kayıt
// Circular dependency problemini çözmek için kullanılır
var (
	registry = make(map[string]func() resource.Resource)
	mu       sync.RWMutex
)

// Register bir resource factory'sini kayıt eder
func Register(slug string, factory func() resource.Resource) {
	mu.Lock()
	defer mu.Unlock()
	registry[slug] = factory
	resource.Register(slug, factory())
}

// Get kayıtlı bir resource'u slug'ına göre alır
func Get(slug string) resource.Resource {
	mu.RLock()
	defer mu.RUnlock()
	if factory, ok := registry[slug]; ok {
		return factory()
	}
	return nil
}

// GetOrPanic kayıtlı bir resource'u alır, bulamazsa panic yapar
func GetOrPanic(slug string) resource.Resource {
	r := Get(slug)
	if r == nil {
		panic("resource not found: " + slug)
	}
	return r
}
`

const organizationResourceTemplate = `package organization

import (
	"%s/entity"
	"%s/resources"
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/resource"
)

func init() {
	resources.Register("organizations", func() resource.Resource {
		return NewOrganizationResource()
	})
}

type OrganizationResource struct {
	resource.OptimizedBase
}

type OrganizationPolicy struct{}

func (p *OrganizationPolicy) ViewAny(ctx *context.Context) bool   { return true }
func (p *OrganizationPolicy) View(ctx *context.Context, model interface{}) bool { return true }
func (p *OrganizationPolicy) Create(ctx *context.Context) bool    { return true }
func (p *OrganizationPolicy) Update(ctx *context.Context, model interface{}) bool { return true }
func (p *OrganizationPolicy) Delete(ctx *context.Context, model interface{}) bool { return true }

type OrganizationResolveFields struct{}

func NewOrganizationResource() *OrganizationResource {
	r := &OrganizationResource{}
	r.SetSlug("organizations")
	r.SetTitle("Organizations")
	r.SetIcon("building")
	r.SetGroup("Management")
	r.SetModel(&entity.Organization{})
	r.SetFieldResolver(&OrganizationResolveFields{})
	r.SetVisible(true)
	r.SetPolicy(&OrganizationPolicy{})
	return r
}

func (r *OrganizationResolveFields) ResolveFields(ctx *context.Context) []fields.Element {
	return []fields.Element{
		fields.ID("ID", "id"),
		fields.Text("Name", "name").Label("Organizasyon Adı").Placeholder("Organizasyon Adı").Required(),
		fields.Email("Email", "email").Label("E-posta").Placeholder("E-posta").Required(),
		fields.Tel("Phone", "phone").Label("Telefon").Placeholder("Telefon").Required(),
		fields.HasOne("BillingInfo", "billing_info", resources.GetOrPanic("billing-info")).
			ForeignKey("organization_id").OwnerKey("id").WithEagerLoad().Label("Fatura Bilgisi").HideOnList().HideOnCreate(),
		fields.HasMany("Addresses", "addresses", resources.GetOrPanic("addresses")).
			ForeignKey("organization_id").OwnerKey("id").WithEagerLoad().Label("Adresler").HideOnList().HideOnCreate(),
		fields.HasMany("Products", "products", resources.GetOrPanic("products")).
			ForeignKey("organization_id").OwnerKey("id").WithEagerLoad().Label("Ürünler").HideOnList().HideOnCreate(),
		fields.HasMany("Shipments", "shipments", resources.GetOrPanic("shipments")).
			ForeignKey("organization_id").OwnerKey("id").WithEagerLoad().Label("Gönderiler").HideOnList().HideOnCreate(),
		fields.Date("CreatedAt", "created_at").HideOnCreate().HideOnUpdate(),
		fields.Date("UpdatedAt", "updated_at").HideOnCreate().HideOnUpdate(),
	}
}

func (r *OrganizationResource) RecordTitle(record interface{}) string {
	if org, ok := record.(*entity.Organization); ok {
		return org.Name
	}
	return ""
}

func (r *OrganizationResource) With() []string {
	return []string{"BillingInfo", "Addresses", "Products", "Shipments"}
}
`

const billingInfoResourceTemplate = `package billing_info

import (
	"%s/entity"
	"%s/resources"
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/resource"
)

func init() {
	resources.Register("billing-info", func() resource.Resource {
		return NewBillingInfoResource()
	})
}

type BillingInfoResource struct {
	resource.OptimizedBase
}

type BillingInfoPolicy struct{}

func (p *BillingInfoPolicy) ViewAny(ctx *context.Context) bool   { return true }
func (p *BillingInfoPolicy) View(ctx *context.Context, model interface{}) bool { return true }
func (p *BillingInfoPolicy) Create(ctx *context.Context) bool    { return true }
func (p *BillingInfoPolicy) Update(ctx *context.Context, model interface{}) bool { return true }
func (p *BillingInfoPolicy) Delete(ctx *context.Context, model interface{}) bool { return true }

type BillingInfoResolveFields struct{}

func NewBillingInfoResource() *BillingInfoResource {
	r := &BillingInfoResource{}
	r.SetSlug("billing-info")
	r.SetTitle("Billing Info")
	r.SetIcon("receipt")
	r.SetGroup("Management")
	r.SetModel(&entity.BillingInfo{})
	r.SetFieldResolver(&BillingInfoResolveFields{})
	r.SetVisible(true)
	r.SetPolicy(&BillingInfoPolicy{})
	return r
}

func (r *BillingInfoResolveFields) ResolveFields(ctx *context.Context) []fields.Element {
	return []fields.Element{
		fields.ID("ID", "id"),
		fields.BelongsTo("Organization", "organization_id", resources.GetOrPanic("organizations")).
			DisplayUsing("name").WithSearchableColumns("name", "email").WithEagerLoad().Label("Organizasyon").Required(),
		fields.Text("TaxNumber", "tax_number").Label("Vergi Numarası").Placeholder("Vergi Numarası").Required(),
		fields.Text("TaxOffice", "tax_office").Label("Vergi Dairesi").Placeholder("Vergi Dairesi").Required(),
		fields.Date("CreatedAt", "created_at").HideOnCreate().HideOnUpdate(),
		fields.Date("UpdatedAt", "updated_at").HideOnCreate().HideOnUpdate(),
	}
}

func (r *BillingInfoResource) RecordTitle(record interface{}) string {
	if info, ok := record.(*entity.BillingInfo); ok {
		return info.TaxNumber
	}
	return ""
}

func (r *BillingInfoResource) With() []string {
	return []string{"Organization"}
}
`

const addressResourceTemplate = `package address

import (
	"%s/entity"
	"%s/resources"
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/resource"
)

func init() {
	resources.Register("addresses", func() resource.Resource {
		return NewAddressResource()
	})
}

type AddressResource struct {
	resource.OptimizedBase
}

type AddressPolicy struct{}

func (p *AddressPolicy) ViewAny(ctx *context.Context) bool   { return true }
func (p *AddressPolicy) View(ctx *context.Context, model interface{}) bool { return true }
func (p *AddressPolicy) Create(ctx *context.Context) bool    { return true }
func (p *AddressPolicy) Update(ctx *context.Context, model interface{}) bool { return true }
func (p *AddressPolicy) Delete(ctx *context.Context, model interface{}) bool { return true }

type AddressResolveFields struct{}

func NewAddressResource() *AddressResource {
	r := &AddressResource{}
	r.SetSlug("addresses")
	r.SetTitle("Addresses")
	r.SetIcon("map-pin")
	r.SetGroup("Management")
	r.SetModel(&entity.Address{})
	r.SetFieldResolver(&AddressResolveFields{})
	r.SetVisible(true)
	r.SetPolicy(&AddressPolicy{})
	return r
}

func (r *AddressResolveFields) ResolveFields(ctx *context.Context) []fields.Element {
	return []fields.Element{
		fields.ID("ID", "id"),
		fields.BelongsTo("Organization", "organization_id", resources.GetOrPanic("organizations")).
			DisplayUsing("name").WithSearchableColumns("name", "email").WithEagerLoad().Label("Organizasyon").Required(),
		fields.Text("Name", "name").Label("Adres Adı").Placeholder("Adres Adı").Required(),
		fields.Textarea("Address", "address").Label("Adres").Placeholder("Adres").Required(),
		fields.Text("City", "city").Label("Şehir").Placeholder("Şehir").Required(),
		fields.Text("Country", "country").Label("Ülke").Placeholder("Ülke").Required(),
		fields.Date("CreatedAt", "created_at").HideOnCreate().HideOnUpdate(),
		fields.Date("UpdatedAt", "updated_at").HideOnCreate().HideOnUpdate(),
	}
}

func (r *AddressResource) RecordTitle(record interface{}) string {
	if addr, ok := record.(*entity.Address); ok {
		return addr.Name
	}
	return ""
}

func (r *AddressResource) With() []string {
	return []string{"Organization"}
}
`

const productResourceTemplate = `package product

import (
	"%s/entity"
	"%s/resources"
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/resource"
)

func init() {
	resources.Register("products", func() resource.Resource {
		return NewProductResource()
	})
}

type ProductResource struct {
	resource.OptimizedBase
}

type ProductPolicy struct{}

func (p *ProductPolicy) ViewAny(ctx *context.Context) bool   { return true }
func (p *ProductPolicy) View(ctx *context.Context, model interface{}) bool { return true }
func (p *ProductPolicy) Create(ctx *context.Context) bool    { return true }
func (p *ProductPolicy) Update(ctx *context.Context, model interface{}) bool { return true }
func (p *ProductPolicy) Delete(ctx *context.Context, model interface{}) bool { return true }

type ProductResolveFields struct{}

func NewProductResource() *ProductResource {
	r := &ProductResource{}
	r.SetSlug("products")
	r.SetTitle("Products")
	r.SetIcon("package")
	r.SetGroup("Management")
	r.SetModel(&entity.Product{})
	r.SetFieldResolver(&ProductResolveFields{})
	r.SetVisible(true)
	r.SetPolicy(&ProductPolicy{})
	return r
}

func (r *ProductResolveFields) ResolveFields(ctx *context.Context) []fields.Element {
	return []fields.Element{
		fields.ID("ID", "id"),
		fields.BelongsTo("Organization", "organization_id", resources.GetOrPanic("organizations")).
			DisplayUsing("name").WithSearchableColumns("name", "email").WithEagerLoad().Label("Organizasyon").Required(),
		fields.Text("Name", "name").Label("Ürün Adı").Placeholder("Ürün Adı").Required(),
		fields.BelongsToMany("Categories", "categories", resources.GetOrPanic("categories")).
			PivotTable("product_categories").WithEagerLoad().Label("Kategoriler").HideOnList(),
		fields.HasMany("ShipmentRows", "shipment_rows", resources.GetOrPanic("shipment-rows")).
			ForeignKey("product_id").OwnerKey("id").WithEagerLoad().Label("Gönderi Satırları").HideOnList().HideOnCreate(),
		fields.Date("CreatedAt", "created_at").HideOnCreate().HideOnUpdate(),
		fields.Date("UpdatedAt", "updated_at").HideOnCreate().HideOnUpdate(),
	}
}

func (r *ProductResource) RecordTitle(record interface{}) string {
	if product, ok := record.(*entity.Product); ok {
		return product.Name
	}
	return ""
}

func (r *ProductResource) With() []string {
	return []string{"Organization", "Categories", "ShipmentRows"}
}
`

const categoryResourceTemplate = `package category

import (
	"%s/entity"
	"%s/resources"
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/resource"
)

func init() {
	resources.Register("categories", func() resource.Resource {
		return NewCategoryResource()
	})
}

type CategoryResource struct {
	resource.OptimizedBase
}

type CategoryPolicy struct{}

func (p *CategoryPolicy) ViewAny(ctx *context.Context) bool   { return true }
func (p *CategoryPolicy) View(ctx *context.Context, model interface{}) bool { return true }
func (p *CategoryPolicy) Create(ctx *context.Context) bool    { return true }
func (p *CategoryPolicy) Update(ctx *context.Context, model interface{}) bool { return true }
func (p *CategoryPolicy) Delete(ctx *context.Context, model interface{}) bool { return true }

type CategoryResolveFields struct{}

func NewCategoryResource() *CategoryResource {
	r := &CategoryResource{}
	r.SetSlug("categories")
	r.SetTitle("Categories")
	r.SetIcon("tag")
	r.SetGroup("Management")
	r.SetModel(&entity.Category{})
	r.SetFieldResolver(&CategoryResolveFields{})
	r.SetVisible(true)
	r.SetPolicy(&CategoryPolicy{})
	return r
}

func (r *CategoryResolveFields) ResolveFields(ctx *context.Context) []fields.Element {
	return []fields.Element{
		fields.ID("ID", "id"),
		fields.Text("Name", "name").Label("Kategori Adı").Placeholder("Kategori Adı").Required(),
		fields.BelongsToMany("Products", "products", resources.GetOrPanic("products")).
			PivotTable("product_categories").WithEagerLoad().Label("Ürünler").HideOnList(),
		fields.Date("CreatedAt", "created_at").HideOnCreate().HideOnUpdate(),
		fields.Date("UpdatedAt", "updated_at").HideOnCreate().HideOnUpdate(),
	}
}

func (r *CategoryResource) RecordTitle(record interface{}) string {
	if cat, ok := record.(*entity.Category); ok {
		return cat.Name
	}
	return ""
}

func (r *CategoryResource) With() []string {
	return []string{"Products"}
}
`

const shipmentResourceTemplate = `package shipment

import (
	"%s/entity"
	"%s/resources"
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/resource"
)

func init() {
	resources.Register("shipments", func() resource.Resource {
		return NewShipmentResource()
	})
}

type ShipmentResource struct {
	resource.OptimizedBase
}

type ShipmentPolicy struct{}

func (p *ShipmentPolicy) ViewAny(ctx *context.Context) bool   { return true }
func (p *ShipmentPolicy) View(ctx *context.Context, model interface{}) bool { return true }
func (p *ShipmentPolicy) Create(ctx *context.Context) bool    { return true }
func (p *ShipmentPolicy) Update(ctx *context.Context, model interface{}) bool { return true }
func (p *ShipmentPolicy) Delete(ctx *context.Context, model interface{}) bool { return true }

type ShipmentResolveFields struct{}

func NewShipmentResource() *ShipmentResource {
	r := &ShipmentResource{}
	r.SetSlug("shipments")
	r.SetTitle("Shipments")
	r.SetIcon("truck")
	r.SetGroup("Operations")
	r.SetModel(&entity.Shipment{})
	r.SetFieldResolver(&ShipmentResolveFields{})
	r.SetVisible(true)
	r.SetPolicy(&ShipmentPolicy{})
	return r
}

func (r *ShipmentResolveFields) ResolveFields(ctx *context.Context) []fields.Element {
	return []fields.Element{
		fields.ID("ID", "id"),
		fields.BelongsTo("Organization", "organization_id", resources.GetOrPanic("organizations")).
			DisplayUsing("name").WithSearchableColumns("name", "email").WithEagerLoad().Label("Organizasyon").Required(),
		fields.Text("Name", "name").Label("Gönderi Adı").Placeholder("Gönderi Adı").Required(),
		fields.HasMany("ShipmentRows", "shipment_rows", resources.GetOrPanic("shipment-rows")).
			ForeignKey("shipment_id").OwnerKey("id").WithEagerLoad().Label("Gönderi Satırları").HideOnList().HideOnCreate(),
		fields.Date("CreatedAt", "created_at").HideOnCreate().HideOnUpdate(),
		fields.Date("UpdatedAt", "updated_at").HideOnCreate().HideOnUpdate(),
	}
}

func (r *ShipmentResource) RecordTitle(record interface{}) string {
	if shipment, ok := record.(*entity.Shipment); ok {
		return shipment.Name
	}
	return ""
}

func (r *ShipmentResource) With() []string {
	return []string{"Organization", "ShipmentRows"}
}
`

const shipmentRowResourceTemplate = `package shipment_row

import (
	"fmt"

	"%s/entity"
	"%s/resources"
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/resource"
)

func init() {
	resources.Register("shipment-rows", func() resource.Resource {
		return NewShipmentRowResource()
	})
}

type ShipmentRowResource struct {
	resource.OptimizedBase
}

type ShipmentRowPolicy struct{}

func (p *ShipmentRowPolicy) ViewAny(ctx *context.Context) bool   { return true }
func (p *ShipmentRowPolicy) View(ctx *context.Context, model interface{}) bool { return true }
func (p *ShipmentRowPolicy) Create(ctx *context.Context) bool    { return true }
func (p *ShipmentRowPolicy) Update(ctx *context.Context, model interface{}) bool { return true }
func (p *ShipmentRowPolicy) Delete(ctx *context.Context, model interface{}) bool { return true }

type ShipmentRowResolveFields struct{}

func NewShipmentRowResource() *ShipmentRowResource {
	r := &ShipmentRowResource{}
	r.SetSlug("shipment-rows")
	r.SetTitle("Shipment Rows")
	r.SetIcon("list")
	r.SetGroup("Operations")
	r.SetModel(&entity.ShipmentRow{})
	r.SetFieldResolver(&ShipmentRowResolveFields{})
	r.SetVisible(true)
	r.SetPolicy(&ShipmentRowPolicy{})
	return r
}

func (r *ShipmentRowResolveFields) ResolveFields(ctx *context.Context) []fields.Element {
	return []fields.Element{
		fields.ID("ID", "id"),
		fields.BelongsTo("Shipment", "shipment_id", resources.GetOrPanic("shipments")).
			DisplayUsing("name").WithSearchableColumns("name").WithEagerLoad().Label("Gönderi").Required(),
		fields.BelongsTo("Product", "product_id", resources.GetOrPanic("products")).
			DisplayUsing("name").WithSearchableColumns("name").WithEagerLoad().Label("Ürün").Required(),
		fields.Number("Quantity", "quantity").Label("Miktar").Placeholder("Miktar").Required(),
		fields.Date("CreatedAt", "created_at").HideOnCreate().HideOnUpdate(),
		fields.Date("UpdatedAt", "updated_at").HideOnCreate().HideOnUpdate(),
	}
}

func (r *ShipmentRowResource) RecordTitle(record interface{}) string {
	if row, ok := record.(*entity.ShipmentRow); ok {
		return fmt.Sprintf("Row #%d", row.ID)
	}
	return ""
}

func (r *ShipmentRowResource) With() []string {
	return []string{"Shipment", "Product"}
}
`

const commentResourceTemplate = `package comment

import (
	"%s/entity"
	"%s/resources"
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/resource"
)

func init() {
	resources.Register("comments", func() resource.Resource {
		return NewCommentResource()
	})
}

type CommentResource struct {
	resource.OptimizedBase
}

type CommentPolicy struct{}

func (p *CommentPolicy) ViewAny(ctx *context.Context) bool   { return true }
func (p *CommentPolicy) View(ctx *context.Context, model interface{}) bool { return true }
func (p *CommentPolicy) Create(ctx *context.Context) bool    { return true }
func (p *CommentPolicy) Update(ctx *context.Context, model interface{}) bool { return true }
func (p *CommentPolicy) Delete(ctx *context.Context, model interface{}) bool { return true }

type CommentResolveFields struct{}

func NewCommentResource() *CommentResource {
	r := &CommentResource{}
	r.SetSlug("comments")
	r.SetTitle("Comments")
	r.SetIcon("message-square")
	r.SetGroup("Content")
	r.SetModel(&entity.Comment{})
	r.SetFieldResolver(&CommentResolveFields{})
	r.SetVisible(true)
	r.SetPolicy(&CommentPolicy{})
	return r
}

func (r *CommentResolveFields) ResolveFields(ctx *context.Context) []fields.Element {
	return []fields.Element{
		fields.ID("ID", "id"),
		fields.MorphTo("Commentable", "commentable").
			Types(
				fields.MorphToType{Type: "products", Resource: resources.GetOrPanic("products")},
				fields.MorphToType{Type: "shipments", Resource: resources.GetOrPanic("shipments")},
			).Label("İlgili Kayıt").Required(),
		fields.Textarea("Content", "content").Label("Yorum").Placeholder("Yorum").Required(),
		fields.Date("CreatedAt", "created_at").HideOnCreate().HideOnUpdate(),
		fields.Date("UpdatedAt", "updated_at").HideOnCreate().HideOnUpdate(),
	}
}

func (r *CommentResource) RecordTitle(record interface{}) string {
	if comment, ok := record.(*entity.Comment); ok {
		if len(comment.Content) > 50 {
			return comment.Content[:50] + "..."
		}
		return comment.Content
	}
	return ""
}
`

const tagResourceTemplate = `package tag

import (
	"%s/entity"
	"%s/resources"
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/resource"
)

func init() {
	resources.Register("tags", func() resource.Resource {
		return NewTagResource()
	})
}

type TagResource struct {
	resource.OptimizedBase
}

type TagPolicy struct{}

func (p *TagPolicy) ViewAny(ctx *context.Context) bool   { return true }
func (p *TagPolicy) View(ctx *context.Context, model interface{}) bool { return true }
func (p *TagPolicy) Create(ctx *context.Context) bool    { return true }
func (p *TagPolicy) Update(ctx *context.Context, model interface{}) bool { return true }
func (p *TagPolicy) Delete(ctx *context.Context, model interface{}) bool { return true }

type TagResolveFields struct{}

func NewTagResource() *TagResource {
	r := &TagResource{}
	r.SetSlug("tags")
	r.SetTitle("Tags")
	r.SetIcon("tag")
	r.SetGroup("Content")
	r.SetModel(&entity.Tag{})
	r.SetFieldResolver(&TagResolveFields{})
	r.SetVisible(true)
	r.SetPolicy(&TagPolicy{})
	return r
}

func (r *TagResolveFields) ResolveFields(ctx *context.Context) []fields.Element {
	return []fields.Element{
		fields.ID("ID", "id"),
		fields.Text("Name", "name").Label("Etiket Adı").Placeholder("Etiket Adı").Required(),
		fields.MorphToMany("Taggables", "taggables").
			Types(
				fields.MorphToManyType{Type: "products", Resource: resources.GetOrPanic("products")},
				fields.MorphToManyType{Type: "shipments", Resource: resources.GetOrPanic("shipments")},
			).PivotTable("taggables").Label("İlgili Kayıtlar").HideOnList(),
		fields.Date("CreatedAt", "created_at").HideOnCreate().HideOnUpdate(),
		fields.Date("UpdatedAt", "updated_at").HideOnCreate().HideOnUpdate(),
	}
}

func (r *TagResource) RecordTitle(record interface{}) string {
	if tag, ok := record.(*entity.Tag); ok {
		return tag.Name
	}
	return ""
}
`

const exampleMainTemplate = `package %[1]s

import (
	_ "%[2]s/resources/address"
	_ "%[2]s/resources/billing_info"
	_ "%[2]s/resources/category"
	_ "%[2]s/resources/organization"
	_ "%[2]s/resources/product"
	_ "%[2]s/resources/shipment"
	_ "%[2]s/resources/shipment_row"

	"github.com/ferdiunal/panel.go/pkg/plugin"
)

// Plugin, %[3]s plugin'i.
type Plugin struct {
	plugin.BasePlugin
}

// init, plugin'i global registry'ye kaydeder.
func init() {
	plugin.Register(&Plugin{})
}

// Name, plugin adını döndürür.
func (p *Plugin) Name() string {
	return "%[3]s"
}

// Version, plugin versiyonunu döndürür.
func (p *Plugin) Version() string {
	return "1.0.0"
}

// Register, plugin'i Panel'e kaydeder.
func (p *Plugin) Register(panel interface{}) error {
	// Tüm resource'lar init() fonksiyonları ile otomatik kayıt edilir
	return nil
}

// Boot, plugin'i boot eder.
func (p *Plugin) Boot(panel interface{}) error {
	// Plugin boot logic
	return nil
}
`
