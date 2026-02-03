package resource

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/ferdiunal/panel.go/internal/auth"
	appContext "github.com/ferdiunal/panel.go/internal/context"
	"github.com/ferdiunal/panel.go/internal/data"
	"github.com/ferdiunal/panel.go/internal/fields"
	"github.com/ferdiunal/panel.go/internal/widget"
	"gorm.io/gorm"
)

// Base, temel Resource yapısıdır.
// resource.Resource arayüzünü implemente eder.
// Resource implementasyonları için gömülü (embedding) olarak kullanılabilir.
type Base struct {
	DataModel     interface{}
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
}

// SettingsSeed is a helper struct for seeding or grouping settings.
type SettingsSeed struct {
	Key   string
	Value map[string]interface{}
}

// Model, veri modelini döner.
func (r Base) Model() interface{} {
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
