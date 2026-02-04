package resource

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/ferdiunal/panel.go/pkg/auth"
	appContext "github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/data"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/widget"
	"gorm.io/gorm"
)

// Base, temel Resource yapısıdır.
// resource.Resource arayüzünü implemente eder.
// Resource implementasyonları için gömülü (embedding) olarak kullanılabilir.
type Base struct {
	DataModel     any
	Identifier    string // Slug, URL tanımlayıcısı
	Label         string // Title, Görünen Başlık
	IconName      string
	GroupName     string
	FieldsVal     []fields.Element
	WidgetsVal    []widget.Card
	Sortable      []Sortable
	PolicyVal     auth.Policy
	DialogType    DialogType
	UploadHandler func(c *appContext.Context, file *multipart.FileHeader) (string, error)
	Seed          SettingsSeed
	ActionsVal    []Action
	FiltersVal    []Filter
}

// SettingsSeed is a helper struct for seeding or grouping settings.
type SettingsSeed struct {
	Key   string
	Value map[string]any
}

// Model, veri modelini döner.
func (r Base) Model() any {
	return r.DataModel
}

// Slug, URL tanımlayıcısını döner.
func (r Base) Slug() string {
	return r.Identifier
}

// Title, insan tarafından okunabilir başlığı döner.
func (r Base) Title() string {
	return r.Label
}

// Icon, menü ikonunu döner.
func (r Base) Icon() string {
	return r.IconName
}

// Group, menü grubunu döner.
func (r Base) Group() string {
	return r.GroupName
}

// Fields, kaynak alanlarını döner.
func (r Base) Fields() []fields.Element {
	return r.FieldsVal
}

// With, eager loading yapılacak ilişkileri döner. (Varsayılan boş)
func (r Base) With() []string {
	return []string{}
}

// Lenses, tanımlı özel görünümleri döner. (Varsayılan boş)
func (r Base) Lenses() []Lens {
	return []Lens{}
}

// Policy, yetkilendirme politikasını döner.
func (r Base) Policy() auth.Policy {
	return r.PolicyVal
}

// GetSortable, varsayılan sıralama ayarlarını döner.
func (r Base) GetSortable() []Sortable {
	return r.Sortable
}

// GetDialogType, diyalog tipini döner.
func (r Base) GetDialogType() DialogType {
	return r.DialogType
}

func (r Base) SetDialogType(dialogType DialogType) Resource {
	r.DialogType = dialogType
	return r
}

func (r Base) Repository(db *gorm.DB) data.DataProvider {
	return nil
}

func (r Base) Cards() []widget.Card {
	if r.WidgetsVal != nil {
		// Convert generic WidgetsVal to []widget.Card if necessary,
		// but since we updated the interface, we should update WidgetsVal type too.
		// For now, let's assume WidgetsVal holds []widget.Card
		return r.WidgetsVal
	}
	return []widget.Card{}
}

func (r Base) StoreHandler(c *appContext.Context, file *multipart.FileHeader, storagePath string, storageURL string) (string, error) {
	if r.UploadHandler != nil {
		return r.UploadHandler(c, file)
	}

	// Default Storage Logic
	// Use configured storage path
	if storagePath == "" {
		storagePath = "./storage/public"
	}
	if storageURL == "" {
		storageURL = "/storage"
	}

	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	localPath := filepath.Join(storagePath, filename)

	// Ensure directory exists
	_ = os.MkdirAll(storagePath, 0755)

	if err := c.Ctx.SaveFile(file, localPath); err != nil {
		return "", err
	}
	// Public URL path
	// Ensure storageURL doesn't end with slash if filename starts with one, but filename is just name.
	// Handle slashes cleanly.
	return fmt.Sprintf("%s/%s", storageURL, filename), nil
}

func (r Base) NavigationOrder() int {
	return 99
}

func (r Base) Visible() bool {
	return true
}

// GetFields, belirli bir bağlamda gösterilecek alanları döner.
// Requirement 11.1: Resource arayüzünü, alanları almak için metodlar içerecek şekilde genişlet
func (r Base) GetFields(ctx *appContext.Context) []fields.Element {
	return r.FieldsVal
}

// GetCards, belirli bir bağlamda gösterilecek card'ları döner.
// Requirement 11.1: Resource arayüzünü, card'ları almak için metodlar içerecek şekilde genişlet
func (r Base) GetCards(ctx *appContext.Context) []widget.Card {
	if r.WidgetsVal != nil {
		return r.WidgetsVal
	}
	return []widget.Card{}
}

// GetLenses, kaynağın tüm lens'lerini döner.
// Requirement 11.1: Resource arayüzünü, lens'leri almak için metodlar içerecek şekilde genişlet
func (r Base) GetLenses() []Lens {
	return r.Lenses()
}

// GetPolicy, kaynağın yetkilendirme politikasını döner.
// Requirement 11.1: Resource arayüzünü, politikaları almak için metodlar içerecek şekilde genişlet
func (r Base) GetPolicy() auth.Policy {
	return r.PolicyVal
}

// ResolveField, bir alanın değerini dinamik olarak hesaplayan ve dönüştüren fonksiyon.
// Requirement 11.2: Resource'ların kendi alan çözümleme mantığını tanımlamasına izin ver
func (r Base) ResolveField(fieldName string, item any) (any, error) {
	// Varsayılan implementasyon: alanı bul ve değerini döndür
	for _, field := range r.FieldsVal {
		if field.GetKey() == fieldName {
			// Alan bulundu, değerini çöz
			// Extract the value from the item
			field.Extract(item)
			// Return the serialized value
			serialized := field.JsonSerialize()
			if val, ok := serialized["value"]; ok {
				return val, nil
			}
			return nil, nil
		}
	}
	// Alan bulunamadı
	return nil, fmt.Errorf("field %s not found", fieldName)
}

// GetActions, kaynağın özel işlemlerini döner.
// Requirement 11.4: Resource arayüzünü, işlemleri almak için metodlar içerecek şekilde genişlet
func (r Base) GetActions() []Action {
	if r.ActionsVal != nil {
		return r.ActionsVal
	}
	return []Action{}
}

// GetFilters, kaynağın filtreleme seçeneklerini döner.
// Requirement 11.4: Resource arayüzünü, filtreleri almak için metodlar içerecek şekilde genişlet
func (r Base) GetFilters() []Filter {
	if r.FiltersVal != nil {
		return r.FiltersVal
	}
	return []Filter{}
}
