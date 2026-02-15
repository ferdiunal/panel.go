// Package handler provides HTTP request handlers for the panel API.
// This package contains the core FieldHandler struct and shared helper methods
// used by all controller functions to handle resource operations.
package handler

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/ferdiunal/panel.go/pkg/auth"
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/core"
	"github.com/ferdiunal/panel.go/pkg/data"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/notification"
	"github.com/ferdiunal/panel.go/pkg/resource"
	"github.com/ferdiunal/panel.go/pkg/widget"
	"github.com/iancoleman/strcase"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
)

const formNullSentinel = "__PANEL_NULL__"

// / Bu yapı, panel API'sinde kaynak (resource) işlemlerini yönetmek için kullanılan
// / merkezi HTTP istek işleyicisidir. Tüm CRUD operasyonları, alan (field) çözümleme,
// / yetkilendirme ve bildirim işlemlerini koordine eder.
// /
// / # Temel Sorumluluklar
// /
// / - **Veri Yönetimi**: GORM veritabanı bağlantısı ve veri sağlayıcı (DataProvider) üzerinden
// /   kaynak verilerini okuma, yazma, güncelleme ve silme işlemlerini gerçekleştirir
// / - **Alan Çözümleme**: Resource'a ait alanları (fields) çözümler, görünürlük kontrolü yapar
// /   ve frontend için uygun formatta serileştirir
// / - **Yetkilendirme**: Policy arayüzü üzerinden kullanıcı yetkilerini kontrol eder
// / - **Dosya Yönetimi**: Dosya yükleme işlemlerinde storage path ve URL yönetimini sağlar
// / - **Widget Desteği**: Dashboard kartları (cards) ve widget'ları yönetir
// / - **Bildirim**: Kullanıcı bildirimlerini NotificationService üzerinden gönderir
// /
// / # Kullanım Senaryoları
// /
// / 1. **Standart Resource İşlemleri**: CRUD operasyonları için NewResourceHandler ile oluşturulur
// / 2. **Lens (Filtered View) İşlemleri**: Filtrelenmiş görünümler için NewLensHandler ile oluşturulur
// / 3. **Özel Veri Sağlayıcılar**: Custom repository implementasyonları için NewFieldHandler ile oluşturulur
// /
// / # Örnek Kullanım
// /
// / ```go
// / // Resource'tan handler oluşturma
// / handler := NewResourceHandler(db, userResource, "/storage", "https://example.com/storage")
// /
// / // Fiber route'a bağlama
// / app.Get("/api/users", func(c *fiber.Ctx) error {
// /     ctx := context.New(c)
// /     return handler.Index(ctx)
// / })
// / ```
// /
// / # Önemli Notlar
// /
// / - **Thread Safety**: Elements alanı concurrent isteklerde mutasyona uğrayabilir,
// /   production ortamında element cloning düşünülmelidir
// / - **Memory Management**: Provider ve DB bağlantıları düzgün kapatılmalıdır
// / - **Policy Kontrolü**: Her işlem öncesi Policy.Can() kontrolü yapılmalıdır
// /
// / # İlişkili Tipler
// /
// / - `data.DataProvider`: Veri erişim katmanı arayüzü
// / - `fields.Element`: Alan tanımı arayüzü
// / - `auth.Policy`: Yetkilendirme arayüzü
// / - `resource.Resource`: Kaynak tanımı arayüzü
type FieldHandler struct {
	Provider            data.DataProvider
	Elements            []fields.Element
	Policy              auth.Policy
	StoragePath         string
	StorageURL          string
	Resource            resource.Resource
	Lens                resource.Lens
	Cards               []widget.Card
	Title               string
	DialogType          resource.DialogType
	NotificationService *notification.Service
}

func collectSearchableColumns(elements []fields.Element) []string {
	columns := make([]string, 0)
	seen := make(map[string]struct{})

	for _, element := range elements {
		searchable, ok := element.(interface{ IsSearchable() bool })
		if !ok || !searchable.IsSearchable() {
			continue
		}

		key := element.GetKey()
		if key == "" {
			continue
		}
		if _, exists := seen[key]; exists {
			continue
		}

		seen[key] = struct{}{}
		columns = append(columns, key)
	}

	return columns
}

func collectRelationshipPreloads(
	db *gorm.DB,
	model interface{},
	explicit []string,
	elements []fields.Element,
) []string {
	candidates := make([]string, 0, len(explicit)+len(elements)*2)
	candidates = append(candidates, explicit...)

	for _, element := range elements {
		switch element.GetView() {
		case "has-many-field", "has-one-field", "belongs-to-many-field", "belongs-to-field":
			key := strings.TrimSpace(element.GetKey())
			if key == "" {
				continue
			}

			candidates = append(candidates, key)
			if strings.HasSuffix(key, "_id") {
				candidates = append(candidates, strings.TrimSuffix(key, "_id"))
			}
		}
	}

	if len(candidates) == 0 {
		return []string{}
	}

	// If DB/schema is unavailable, return deduplicated candidates as-is.
	if db == nil {
		seen := make(map[string]struct{}, len(candidates))
		preloads := make([]string, 0, len(candidates))
		for _, candidate := range candidates {
			candidate = strings.TrimSpace(candidate)
			if candidate == "" {
				continue
			}
			if _, ok := seen[candidate]; ok {
				continue
			}
			seen[candidate] = struct{}{}
			preloads = append(preloads, candidate)
		}
		return preloads
	}

	// Build relationship name index from model schema.
	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(model); err != nil || stmt.Schema == nil {
		seen := make(map[string]struct{}, len(candidates))
		preloads := make([]string, 0, len(candidates))
		for _, candidate := range candidates {
			candidate = strings.TrimSpace(candidate)
			if candidate == "" {
				continue
			}
			if _, ok := seen[candidate]; ok {
				continue
			}
			seen[candidate] = struct{}{}
			preloads = append(preloads, candidate)
		}
		return preloads
	}

	relationshipIndex := make(map[string]string)
	for relationName := range stmt.Schema.Relationships.Relations {
		relationshipIndex[relationName] = relationName
		relationshipIndex[strings.ToLower(relationName)] = relationName
		relationshipIndex[strings.ToLower(strcase.ToSnake(relationName))] = relationName
	}

	seen := make(map[string]struct{}, len(candidates))
	preloads := make([]string, 0, len(candidates))
	addPreload := func(value string) {
		if value == "" {
			return
		}
		if _, ok := seen[value]; ok {
			return
		}
		seen[value] = struct{}{}
		preloads = append(preloads, value)
	}

	for _, candidate := range candidates {
		candidate = strings.TrimSpace(candidate)
		if candidate == "" {
			continue
		}

		// Nested preload paths should remain intact.
		if strings.Contains(candidate, ".") {
			addPreload(candidate)
			continue
		}

		tryKeys := []string{
			candidate,
			strings.ToLower(candidate),
			strcase.ToCamel(candidate),
			strings.ToLower(strcase.ToSnake(candidate)),
		}

		if strings.HasSuffix(candidate, "_id") {
			base := strings.TrimSuffix(candidate, "_id")
			tryKeys = append(tryKeys, base, strings.ToLower(base), strcase.ToCamel(base))
		}

		matched := ""
		for _, key := range tryKeys {
			if normalized, ok := relationshipIndex[key]; ok {
				matched = normalized
				break
			}
		}

		if matched != "" {
			addPreload(matched)
			continue
		}

		// Ignore unresolved simple relations to avoid runtime GORM preload errors.
		// Nested paths are handled above and preserved.
	}

	return preloads
}

// / Bu fonksiyon, özel bir veri sağlayıcı (DataProvider) ile yeni bir FieldHandler oluşturur.
// / Minimal konfigürasyonla handler oluşturmak için kullanılır.
// /
// / # Parametreler
// /
// / - `provider`: Veri erişim işlemlerini gerçekleştirecek DataProvider implementasyonu
// /
// / # Döndürür
// /
// / - Yapılandırılmış FieldHandler pointer'ı
// /
// / # Kullanım Senaryoları
// /
// / 1. **Özel Repository**: Custom repository implementasyonu kullanılacaksa
// / 2. **Test Ortamı**: Mock provider ile test senaryoları için
// / 3. **Minimal Setup**: Sadece provider'a ihtiyaç duyulan basit durumlar için
// /
// / # Örnek Kullanım
// /
// / ```go
// / // Özel repository ile handler oluşturma
// / customRepo := NewCustomUserRepository(db)
// / handler := NewFieldHandler(customRepo)
// / handler.SetElements(userFields)
// / ```
// /
// / # Önemli Notlar
// /
// / - Bu fonksiyon minimal bir handler oluşturur, diğer alanlar manuel olarak ayarlanmalıdır
// / - Production kullanımı için genellikle NewResourceHandler tercih edilmelidir
// / - Provider nil olmamalıdır, aksi halde runtime panic oluşur
// /
// / # İlişkili Fonksiyonlar
// /
// / - `NewResourceHandler`: Resource'tan tam yapılandırılmış handler oluşturur
// / - `NewLensHandler`: Lens ile filtrelenmiş handler oluşturur
func NewFieldHandler(provider data.DataProvider) *FieldHandler {
	return &FieldHandler{
		Provider: provider,
	}
}

// / Bu metod, handler'a ait alan (field) listesini ayarlar.
// /
// / # Parametreler
// /
// / - `elements`: Kaynak için tanımlanmış alan listesi
// /
// / # Kullanım Senaryoları
// /
// / 1. **Dinamik Alan Yönetimi**: Runtime'da alanları değiştirmek için
// / 2. **Test Senaryoları**: Mock alanlar ile test yapmak için
// / 3. **Özel Konfigürasyon**: NewFieldHandler sonrası manuel alan ataması için
// /
// / # Örnek Kullanım
// /
// / ```go
// / handler := NewFieldHandler(provider)
// / handler.SetElements([]fields.Element{
// /     fields.NewID(),
// /     fields.NewText("name").SetLabel("İsim"),
// /     fields.NewEmail("email").SetLabel("E-posta"),
// / })
// / ```
// /
// / # Önemli Notlar
// /
// / - Bu metod mevcut Elements listesini tamamen değiştirir
// / - Concurrent erişimlerde thread-safe değildir
// / - Genellikle initialization sırasında kullanılmalıdır
func (h *FieldHandler) SetElements(elements []fields.Element) {
	h.Elements = elements
}

// / Bu fonksiyon, bir Resource tanımından tam yapılandırılmış bir FieldHandler oluşturur.
// / Production ortamında kullanılması önerilen ana handler oluşturma fonksiyonudur.
// /
// / # Temel İşlevler
// /
// / 1. **Veri Sağlayıcı Kurulumu**: Resource'un custom repository'si varsa onu kullanır,
// /    yoksa otomatik olarak GormDataProvider oluşturur
// / 2. **Arama Kolonları**: Searchable olarak işaretlenmiş tüm alanları otomatik tespit eder
// / 3. **İlişki Yükleme**: EAGER_LOADING stratejisine sahip relationship alanlarını
// /    otomatik olarak GORM Preload listesine ekler
// / 4. **Widget Yönetimi**: Resource'a ait dashboard kartlarını yükler
// / 5. **Bildirim Servisi**: Kullanıcı bildirimlerini göndermek için servis başlatır
// /
// / # Parametreler
// /
// / - `db`: GORM veritabanı bağlantısı
// / - `res`: İşlenecek kaynak tanımı (Resource interface implementasyonu)
// / - `storagePath`: Dosya yüklemeleri için fiziksel depolama yolu (örn: "/var/www/storage")
// / - `storageURL`: Dosyalara erişim için public URL (örn: "https://example.com/storage")
// /
// / # Döndürür
// /
// / - Tam yapılandırılmış FieldHandler pointer'ı
// /
// / # Kullanım Senaryoları
// /
// / 1. **Standart CRUD API**: Resource tanımından otomatik API endpoint'leri oluşturma
// / 2. **Custom Repository**: Özel veri erişim mantığı ile çalışma
// / 3. **İlişkisel Veriler**: Eager loading ile performanslı ilişki yükleme
// / 4. **Dosya Yönetimi**: Dosya upload/download işlemleri için storage yapılandırması
// /
// / # Örnek Kullanım
// /
// / ```go
// / // Resource tanımı
// / type UserResource struct {
// /     resource.BaseResource
// / }
// /
// / func (r *UserResource) Model() interface{} {
// /     return &models.User{}
// / }
// /
// / func (r *UserResource) Fields() []fields.Element {
// /     return []fields.Element{
// /         fields.NewID(),
// /         fields.NewText("name").SetSearchable(true),
// /         fields.NewEmail("email").SetSearchable(true),
// /         fields.NewBelongsTo("role").SetLoadingStrategy(fields.EAGER_LOADING),
// /     }
// / }
// /
// / // Handler oluşturma
// / userResource := &UserResource{}
// / handler := NewResourceHandler(
// /     db,
// /     userResource,
// /     "/var/www/storage",
// /     "https://example.com/storage",
// / )
// /
// / // Fiber route'lara bağlama
// / app.Get("/api/users", handler.Index)
// / app.Post("/api/users", handler.Store)
// / app.Get("/api/users/:id", handler.Show)
// / app.Put("/api/users/:id", handler.Update)
// / app.Delete("/api/users/:id", handler.Destroy)
// / ```
// /
// / # İlişki Yükleme Stratejisi
// /
// / Fonksiyon, relationship alanlarını otomatik olarak analiz eder:
// / - **EAGER_LOADING**: Alan key'i CamelCase'e çevrilerek GORM Preload listesine eklenir
// / - **LAZY_LOADING**: İlişki sadece talep edildiğinde yüklenir (varsayılan)
// /
// / # Arama Kolonları
// /
// / `IsSearchable()` true dönen tüm alanlar otomatik olarak arama kolonları listesine eklenir.
// / Bu, Index endpoint'inde `?search=query` parametresi ile arama yapılmasını sağlar.
// /
// / # Önemli Notlar
// /
// / - **Repository Önceliği**: Resource.Repository() nil dönmezse, custom repository kullanılır
// / - **CamelCase Dönüşümü**: İlişki isimleri GORM struct field isimleriyle eşleşmeli
// / - **Duplicate Kontrolü**: Aynı ilişki birden fazla kez eklenmez
// / - **Notification Service**: Her handler için yeni bir notification service instance'ı oluşturulur
// /
// / # Performans Notları
// /
// / - Eager loading çok fazla ilişki için N+1 problemini önler ancak memory kullanımını artırır
// / - Büyük veri setlerinde lazy loading tercih edilebilir
// / - Search kolonları index'lenmeli performans için
// /
// / # İlişkili Fonksiyonlar
// /
// / - `NewFieldHandler`: Minimal handler oluşturma
// / - `NewLensHandler`: Filtrelenmiş görünüm için handler oluşturma
func NewResourceHandler(client interface{}, res resource.Resource, storagePath, storageURL string) *FieldHandler {
	var provider data.DataProvider

	// Type assertion for GORM DB
	db, ok := client.(*gorm.DB)
	if !ok {
		panic("client must be *gorm.DB")
	}

	if repo := res.Repository(db); repo != nil {
		provider = repo
	} else {
		provider = data.NewGormDataProvider(db, res.Model())
	}

	preloads := collectRelationshipPreloads(db, res.Model(), res.With(), res.Fields())
	provider.SetWith(preloads)
	provider.SetSearchColumns(collectSearchableColumns(res.Fields()))

	// Initialize notification service with provider
	notificationService := notification.NewService(provider)

	return &FieldHandler{
		Provider:            provider,
		Elements:            nil, // Lazy load edilecek
		Policy:              res.Policy(),
		Resource:            res,
		StoragePath:         storagePath,
		StorageURL:          storageURL,
		Cards:               res.Cards(),
		Title:               res.Title(),
		DialogType:          res.GetDialogType(),
		NotificationService: notificationService,
	}
}

// / Bu fonksiyon, bir Resource ve Lens tanımından filtrelenmiş görünüm için FieldHandler oluşturur.
// / Lens, belirli bir sorgu veya filtreleme mantığı ile kaynak verilerinin alt kümesini gösterir.
// /
// / # Temel İşlevler
// /
// / 1. **Lens Query Kullanımı**: Lens.Query() metodundan gelen özel sorguyu base query olarak kullanır
// / 2. **Özel Alan Listesi**: Lens'e özgü alan listesini (Lens.Fields()) kullanır
// / 3. **Arama Kolonları**: Lens alanlarından searchable olanları otomatik tespit eder
// / 4. **Başlık Yönetimi**: Lens.Name() ile özel başlık atar
// /
// / # Parametreler
// /
// / - `db`: GORM veritabanı bağlantısı
// / - `res`: Ana kaynak tanımı (Resource interface implementasyonu)
// / - `lens`: Filtreleme mantığını içeren Lens tanımı
// /
// / # Döndürür
// /
// / - Lens için yapılandırılmış FieldHandler pointer'ı
// /
// / # Kullanım Senaryoları
// /
// / 1. **Filtrelenmiş Listeler**: Belirli kriterlere uyan kayıtları gösterme
// /    - Örnek: "Aktif Kullanıcılar", "Bu Ayki Siparişler", "Onay Bekleyenler"
// / 2. **Özel Görünümler**: Farklı alan kombinasyonları ile aynı kaynağı gösterme
// / 3. **Dashboard Widgets**: Özet veya istatistiksel görünümler için
// / 4. **Raporlama**: Özel sorgu mantığı ile rapor oluşturma
// /
// / # Örnek Kullanım
// /
// / ```go
// / // Lens tanımı
// / type ActiveUsersLens struct{}
// /
// / func (l *ActiveUsersLens) Name() string {
// /     return "Aktif Kullanıcılar"
// / }
// /
// / func (l *ActiveUsersLens) Query(db *gorm.DB) *gorm.DB {
// /     return db.Where("status = ?", "active").Where("last_login > ?", time.Now().AddDate(0, -1, 0))
// / }
// /
// / func (l *ActiveUsersLens) Fields() []fields.Element {
// /     return []fields.Element{
// /         fields.NewID(),
// /         fields.NewText("name").SetSearchable(true),
// /         fields.NewText("email").SetSearchable(true),
// /         fields.NewDateTime("last_login"),
// /     }
// / }
// /
// / // Handler oluşturma
// / userResource := &UserResource{}
// / activeLens := &ActiveUsersLens{}
// / handler := NewLensHandler(db, userResource, activeLens)
// /
// / // Fiber route'a bağlama
// / app.Get("/api/users/active", handler.Index)
// / ```
// /
// / # Lens vs Resource Farkları
// /
// / | Özellik | Resource | Lens |
// / |---------|----------|------|
// / | Query | Tüm kayıtlar | Filtrelenmiş kayıtlar |
// / | Fields | Resource.Fields() | Lens.Fields() |
// / | Title | Resource.Title() | Lens.Name() |
// / | Cards | Resource.Cards() | Yok (Lens'te desteklenmez) |
// / | Policy | Resource.Policy() | Yok (Lens'te desteklenmez) |
// /
// / # Önemli Notlar
// /
// / - **Query Encapsulation**: Lens.Query() tüm filtreleme mantığını içermeli
// / - **Preload Yönetimi**: Lens query'si gerekli preload'ları içermeli
// / - **Sınırlı Özellikler**: Lens handler'lar Cards ve Policy desteği içermez
// / - **Resource Bağlantısı**: Lens her zaman bir Resource'a bağlıdır
// /
// / # Performans Notları
// /
// / - Lens query'leri index'lenmeli kolonlar kullanmalı
// / - Karmaşık JOIN'ler için view kullanımı düşünülebilir
// / - Büyük veri setlerinde pagination önemlidir
// /
// / # İlişkili Fonksiyonlar
// /
// / - `NewResourceHandler`: Tam özellikli resource handler oluşturma
// / - `NewFieldHandler`: Minimal handler oluşturma
func NewLensHandler(client interface{}, res resource.Resource, lens resource.Lens) *FieldHandler {
	// Type assertion for GORM DB
	db, ok := client.(*gorm.DB)
	if !ok {
		panic("client must be *gorm.DB for lens handler")
	}
	provider := data.NewGormDataProvider(db, res.Model())
	preloadElements := append([]fields.Element{}, res.Fields()...)
	preloadElements = append(preloadElements, lens.Fields()...)
	preloads := collectRelationshipPreloads(db, res.Model(), res.With(), preloadElements)
	provider.SetWith(preloads)
	provider.SetBaseQuery(lens.GetQuery())
	provider.SetSearchColumns(collectSearchableColumns(lens.Fields()))

	return &FieldHandler{
		Provider: provider,
		Elements: nil, // Lazy load edilecek
		Policy:   res.Policy(),
		Title:    lens.Name(),
		Resource: res,
		Lens:     lens,
	}
}

// / getElements, context ile field'ları lazy load eder.
// /
// / Bu metod, handler initialization sırasında field'ları resolve etmek yerine,
// / ilk request geldiğinde context ile birlikte lazy load eder. Bu sayede
// / i18n translation'lar doğru context ile çalışır.
// /
// / # Parametreler
// /
// / - `ctx`: İstek context'i (nil olabilir)
// /
// / # Döndürür
// /
// / - Field listesi
// /
// / # İş Akışı
// /
// / 1. Context nil ise eski davranışa fallback (Resource.Fields() çağır)
// / 2. Context'ten cache'lenmiş field'ları kontrol et
// / 3. Cache'de varsa cache'den döndür
// / 4. Cache'de yoksa GetFieldsWithContext ile resolve et ve cache'e ekle
// / 5. GetFieldsWithContext yoksa fallback (Resource.Fields() çağır)
func (h *FieldHandler) getElements(ctx *context.Context) []fields.Element {
	// Lens-specific field resolution has priority.
	if h.Lens != nil {
		if ctx != nil {
			if lensFields := h.Lens.GetFields(ctx); len(lensFields) > 0 {
				return lensFields
			}
		}
		if lensFields := h.Lens.Fields(); len(lensFields) > 0 {
			return lensFields
		}
	}

	if ctx == nil {
		// Fallback: context yoksa eski davranış
		return h.Resource.Fields()
	}

	// Field'ları resolve et
	if optimized, ok := h.Resource.(interface {
		GetFieldsWithContext(*context.Context) []fields.Element
	}); ok {
		return optimized.GetFieldsWithContext(ctx)
	}

	// Fallback: eski davranış
	return h.Resource.Fields()
}

// / Bu metod, kaynak listesi (index) endpoint'ini işler.
// / İşlemi resource_index_controller.go dosyasındaki HandleResourceIndex fonksiyonuna yönlendirir.
// /
// / # Parametreler
// /
// / - `c`: Panel context nesnesi (HTTP request/response wrapper)
// /
// / # Döndürür
// /
// / - `error`: İşlem hatası veya nil
// /
// / # Kullanım Senaryoları
// /
// / - Kaynak listesini sayfalama ile getirme
// / - Arama ve filtreleme işlemleri
// / - Sıralama (sorting) işlemleri
// /
// / # Örnek Kullanım
// /
// / ```go
// / app.Get("/api/users", func(c *fiber.Ctx) error {
// /     ctx := context.New(c)
// /     return handler.Index(ctx)
// / })
// / ```
func (h *FieldHandler) Index(c *context.Context) error {
	return HandleResourceIndex(h, c)
}

// / Bu metod, tek bir kaynağı detaylı olarak gösterir.
// / İşlemi resource_show_controller.go dosyasındaki HandleResourceShow fonksiyonuna yönlendirir.
// /
// / # Parametreler
// /
// / - `c`: Panel context nesnesi (HTTP request/response wrapper)
// /
// / # Döndürür
// /
// / - `error`: İşlem hatası veya nil
// /
// / # Kullanım Senaryoları
// /
// / - Tek bir kaydın tüm detaylarını görüntüleme
// / - İlişkili verileri (relationships) yükleme
// / - Detay sayfası için veri sağlama
// /
// / # Örnek Kullanım
// /
// / ```go
// / app.Get("/api/users/:id", func(c *fiber.Ctx) error {
// /     ctx := context.New(c)
// /     return handler.Show(ctx)
// / })
// / ```
func (h *FieldHandler) Show(c *context.Context) error {
	return HandleResourceShow(h, c)
}

// / Bu metod, kaynak düzenleme formunu hazırlar.
// / İşlemi resource_edit_controller.go dosyasındaki HandleResourceEdit fonksiyonuna yönlendirir.
// /
// / # Parametreler
// /
// / - `c`: Panel context nesnesi (HTTP request/response wrapper)
// /
// / # Döndürür
// /
// / - `error`: İşlem hatası veya nil
// /
// / # Kullanım Senaryoları
// /
// / - Düzenleme formu için mevcut veriyi getirme
// / - Alan seçeneklerini (options) çözümleme
// / - Form validasyon kurallarını sağlama
// /
// / # Örnek Kullanım
// /
// / ```go
// / app.Get("/api/users/:id/edit", func(c *fiber.Ctx) error {
// /     ctx := context.New(c)
// /     return handler.Edit(ctx)
// / })
// / ```
func (h *FieldHandler) Edit(c *context.Context) error {
	return HandleResourceEdit(h, c)
}

// / Bu metod, kaynak detay sayfası için veri sağlar.
// / İşlemi resource_detail_controller.go dosyasındaki HandleResourceDetail fonksiyonuna yönlendirir.
// /
// / # Parametreler
// /
// / - `c`: Panel context nesnesi (HTTP request/response wrapper)
// /
// / # Döndürür
// /
// / - `error`: İşlem hatası veya nil
// /
// / # Kullanım Senaryoları
// /
// / - Detay görünümü için özelleştirilmiş alan listesi
// / - Read-only alanları gösterme
// / - İlişkili verileri detaylı gösterme
// /
// / # Örnek Kullanım
// /
// / ```go
// / app.Get("/api/users/:id/detail", func(c *fiber.Ctx) error {
// /     ctx := context.New(c)
// /     return handler.Detail(ctx)
// / })
// / ```
func (h *FieldHandler) Detail(c *context.Context) error {
	return HandleResourceDetail(h, c)
}

// / Bu metod, bir kaynak öğesinden (item) alan verilerini çıkarır ve frontend için
// / uygun formatta serileştirir. Tüm alan çözümleme, görünürlük kontrolü ve özel
// / işlemleri (MorphTo, options, callbacks) gerçekleştirir.
// /
// / # Temel İşlevler
// /
// / 1. **Görünürlük Kontrolü**: Her alan için IsVisible() kontrolü yapar
// / 2. **Veri Çıkarma**: Element.Extract() ile item'dan veri çeker
// / 3. **Serileştirme**: Element.JsonSerialize() ile JSON formatına dönüştürür
// / 4. **Options Çözümleme**: ResolveFieldOptions ile dinamik seçenekleri yükler
// / 5. **MorphTo İşleme**: Polymorphic ilişkiler için display field'ları çözümler
// / 6. **Callback Uygulama**: ResolveCallback varsa veri üzerinde çalıştırır
// /
// / # Parametreler
// /
// / - `c`: Fiber context (HTTP request/response)
// / - `ctx`: Resource context (görünürlük ve metadata bilgisi)
// / - `item`: Veri çıkarılacak kaynak öğesi (model instance)
// / - `elements`: İşlenecek alan listesi
// /
// / # Döndürür
// /
// / - `map[string]interface{}`: Alan key'leri ile serileştirilmiş alan verilerinin map'i
// /
// / # Kullanım Senaryoları
// /
// / 1. **Index Endpoint**: Liste görünümünde her kayıt için alan verilerini hazırlama
// / 2. **Show Endpoint**: Detay görünümünde tek kayıt için alan verilerini hazırlama
// / 3. **Edit Endpoint**: Düzenleme formu için mevcut değerleri hazırlama
// / 4. **Custom Endpoints**: Özel endpoint'lerde alan verilerini serileştirme
// /
// / # MorphTo Alan İşleme
// /
// / MorphTo alanları (polymorphic relationships) için özel işlem yapılır:
// / - `type` ve `id` değerleri alınır
// / - `displays` prop'undan ilgili type için display field bulunur
// / - Veritabanından display field değeri sorgulanır
// / - Sonuç data objesine eklenir
// /
// / # Örnek Kullanım
// /
// / ```go
// / // Tek bir kayıt için alan verilerini çözümleme
// / user := &models.User{ID: 1, Name: "John", Email: "john@example.com"}
// / resourceData := handler.resolveResourceFields(c, ctx, user, handler.Elements)
// /
// / // Sonuç formatı:
// / // {
// / //   "id": {"key": "id", "data": 1, "type": "id", ...},
// / //   "name": {"key": "name", "data": "John", "type": "text", ...},
// / //   "email": {"key": "email", "data": "john@example.com", "type": "email", ...}
// / // }
// / ```
// /
// / # MorphTo Örneği
// /
// / ```go
// / // MorphTo alan tanımı
// / fields.NewMorphTo("commentable").
// /     SetDisplays(map[string]string{
// /         "posts": "title",
// /         "videos": "name",
// /     })
// /
// / // Kayıt verisi
// / comment := &models.Comment{
// /     ID: 1,
// /     CommentableType: "posts",
// /     CommentableID: 5,
// / }
// /
// / // resolveResourceFields çağrısı sonrası:
// / // {
// / //   "commentable": {
// / //     "key": "commentable",
// / //     "data": {
// / //       "type": "posts",
// / //       "id": 5,
// / //       "title": "My First Post"  // Veritabanından çekildi
// / //     },
// / //     ...
// / //   }
// / // }
// / ```
// /
// / # Önemli Notlar
// /
// / - **Mutasyon Uyarısı**: Element.Extract() metodu element state'ini değiştirir,
// /   concurrent isteklerde thread-safe değildir
// / - **Debug Çıktıları**: MorphTo işlemleri için fmt.Printf debug logları içerir
// / - **Null Kontrolü**: MorphTo için type ve id null kontrolü yapılır
// / - **Callback Önceliği**: ResolveCallback en son uygulanır, tüm işlemlerden sonra
// /
// / # Performans Notları
// /
// / - MorphTo alanları için her kayıt başına ek veritabanı sorgusu yapılır
// / - Büyük listelerde MorphTo alanları performans sorununa yol açabilir
// / - Eager loading ile ilişkileri önceden yüklemek daha performanslıdır
// / - Options çözümleme de ek sorgular yapabilir
// /
// / # İlişkili Metodlar
// /
// / - `ResolveFieldOptions`: Alan seçeneklerini çözümler
// / - `Element.Extract`: Item'dan veri çıkarır
// / - `Element.JsonSerialize`: Veriyi JSON formatına dönüştürür
// / - `Element.GetResolveCallback`: Özel callback fonksiyonunu alır
func (h *FieldHandler) resolveResourceFields(c *fiber.Ctx, ctx *core.ResourceContext, item interface{}, elements []fields.Element) map[string]interface{} {
	resourceData := make(map[string]interface{})
	for _, element := range elements {
		if !element.IsVisible(ctx) {
			continue
		}
		// Clone logic or direct access warning applies here too as noted in previous Index implementation.
		// For now, we proceed with direct extraction which mutates element state.
		// In a real high-concurrency scenario, elements should be cloned or Extract should return value.

		// Resolve RelatedResource for relationship fields (HasMany, BelongsTo, etc.)
		// Bu, circular dependency sorununu önlemek için string slug kullanıldığında gereklidir.
		// Resource registry'den resource instance'ı alınır ve field'a set edilir.
		if relField, ok := element.(*fields.HasManyField); ok {
			if relField.GetRelatedResource() == nil && relField.GetRelatedResourceSlug() != "" {
				relatedResource := resource.Get(relField.GetRelatedResourceSlug())
				if relatedResource != nil {
					relField.SetRelatedResource(relatedResource)
				}
			}
		}
		if relField, ok := element.(*fields.BelongsToManyField); ok {
			if relField.RelatedResource == nil && relField.GetRelatedResourceSlug() != "" {
				relatedResource := resource.Get(relField.GetRelatedResourceSlug())
				if relatedResource != nil {
					relField.RelatedResource = relatedResource
				}
			}
		}

		element.Extract(item)
		serialized := element.JsonSerialize()
		normalizeRelationshipCollectionData(element.GetView(), serialized)

		// Resolve options
		h.ResolveFieldOptions(element, serialized, item)

		// Resolve MorphTo display fields
		if element.GetView() == "morph-to-field" {
			fmt.Printf("[DEBUG] MorphTo field detected: %s\n", element.GetKey())
			if data, ok := serialized["data"].(map[string]interface{}); ok {
				morphType, _ := data["type"].(string)
				morphID := data["id"]
				fmt.Printf("[DEBUG] MorphTo data - type: %s, id: %v\n", morphType, morphID)

				if morphType != "" && morphID != nil {
					// Get display field from props
					if props, ok := serialized["props"].(map[string]interface{}); ok {
						if displaysRaw, ok := props["displays"]; ok {
							fmt.Printf("[DEBUG] Displays found: %+v\n", displaysRaw)
							// displays can be map[string]string or map[string]interface{}
							var displayField string
							switch displays := displaysRaw.(type) {
							case map[string]string:
								displayField = displays[morphType]
							case map[string]interface{}:
								if df, ok := displays[morphType].(string); ok {
									displayField = df
								}
							}
							fmt.Printf("[DEBUG] Display field for type '%s': %s\n", morphType, displayField)

							if displayField != "" {
								// TODO: Implement MorphTo display field resolution via Provider
								// This requires adding a QueryTable method to DataProvider interface
								// or using a specialized relationship resolver
								fmt.Printf("[DEBUG] TODO: Resolve MorphTo display for table '%s' field '%s' with id %v\n", strings.ToLower(morphType), displayField, morphID)
							}
						} else {
							fmt.Printf("[DEBUG] No displays found in props\n")
						}
					}
				}
			}
		}

		// Apply callback if exists
		if callback := element.GetResolveCallback(); callback != nil {
			if val, ok := serialized["data"]; ok {
				serialized["data"] = callback(val, item, c)
			}
		}
		applyDisplayCallback(element, serialized, item)

		resourceData[serialized["key"].(string)] = serialized
	}
	return resourceData
}

func normalizeRelationshipCollectionData(view string, serialized map[string]interface{}) {
	switch view {
	case "has-many-field", "belongs-to-many-field", "morph-to-many-field":
	default:
		return
	}

	data, ok := serialized["data"]
	if !ok || data == nil {
		serialized["data"] = []interface{}{}
		return
	}

	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Slice && v.IsNil() {
		serialized["data"] = []interface{}{}
	}
}

func applyDisplayCallback(element fields.Element, serialized map[string]interface{}, item interface{}) {
	callback := element.GetDisplayCallback()
	if callback == nil {
		return
	}

	value := serialized["data"]
	displayValue := callback(value, item)

	if displayElement, ok := displayValue.(core.Element); ok {
		component := serializeDisplayElement(displayElement)

		if componentData, exists := component["data"]; !exists || componentData == nil {
			component["data"] = value
		}

		if key, ok := component["key"].(string); !ok || key == "" {
			component["key"] = element.GetKey()
		}

		if name, ok := component["name"].(string); !ok || name == "" {
			component["name"] = element.GetName()
		}

		serialized["data"] = component
		return
	}

	serialized["data"] = displayValue
}

func serializeDisplayElement(element core.Element) map[string]interface{} {
	serialized := element.JsonSerialize()

	props, ok := serialized["props"].(map[string]interface{})
	if !ok || props == nil {
		return serialized
	}

	rawFields, exists := props["fields"]
	if !exists {
		return serialized
	}

	children, ok := rawFields.([]core.Element)
	if !ok {
		return serialized
	}

	serializedChildren := make([]map[string]interface{}, 0, len(children))
	for _, child := range children {
		serializedChildren = append(serializedChildren, serializeDisplayElement(child))
	}
	props["fields"] = serializedChildren
	serialized["props"] = props

	return serialized
}

// / Bu metod, bir alan için dinamik seçenekleri (options) çözümler ve serileştirilmiş
// / veri yapısına ekler. AutoOptions özelliği ve manuel option callback'lerini destekler.
// /
// / # Temel İşlevler
// /
// / 1. **AutoOptions İşleme**: İlişkisel alanlar için otomatik seçenek listesi oluşturur
// / 2. **HasOne Filtreleme**: HasOne alanları için kullanılabilir kayıtları filtreler
// / 3. **BelongsTo/BelongsToMany**: Tüm ilişkili kayıtları seçenek olarak getirir
// / 4. **Callback Çözümleme**: Manuel option callback fonksiyonlarını çalıştırır
// / 5. **Props Yönetimi**: Çözümlenen seçenekleri props objesine ekler
// /
// / # Parametreler
// /
// / - `element`: Seçenekleri çözümlenecek alan
// / - `serialized`: Alan için serileştirilmiş veri yapısı (JsonSerialize sonucu)
// / - `item`: Mevcut kaynak öğesi (HasOne filtreleme için kullanılır)
// /
// / # AutoOptions Yapılandırması
// /
// / AutoOptions, ilişkisel alanlar için otomatik seçenek listesi oluşturur:
// / - `Enabled`: AutoOptions özelliği aktif mi?
// / - `DisplayField`: Seçeneklerde gösterilecek alan (örn: "name", "title")
// / - `related_resource`: İlişkili tablonun adı (props'tan alınır)
// /
// / # Alan Tiplerine Göre İşlem
// /
// / ## HasOne Alanları
// /
// / HasOne ilişkileri için özel filtreleme yapılır:
// / - Foreign key NULL olan kayıtlar (henüz ilişkilendirilmemiş)
// / - VEYA mevcut item'a ait kayıtlar (düzenleme sırasında)
// /
// / ```sql
// / SELECT id, display_field FROM related_table
// / WHERE foreign_key IS NULL OR foreign_key = current_item_id
// / ```
// /
// / ## BelongsTo ve BelongsToMany Alanları
// /
// / Tüm ilişkili kayıtlar seçenek olarak getirilir:
// /
// / ```sql
// / SELECT id, display_field FROM related_table
// / ```
// /
// / # Kullanım Senaryoları
// /
// / 1. **Form Seçenekleri**: Düzenleme/oluşturma formlarında dropdown seçenekleri
// / 2. **Dinamik Listeler**: Runtime'da değişen seçenek listeleri
// / 3. **İlişki Yönetimi**: İlişkisel alanlar için kullanılabilir kayıtları gösterme
// / 4. **Özel Seçenekler**: Callback ile hesaplanan veya filtrelenmiş seçenekler
// /
// / # Örnek Kullanım
// /
// / ## AutoOptions ile HasOne
// /
// / ```go
// / // Alan tanımı
// / fields.NewHasOne("profile").
// /     SetRelatedResource("profiles").
// /     SetForeignKey("user_id").
// /     SetAutoOptions("full_name")
// /
// / // Kullanıcı düzenleme sırasında:
// / // - Henüz atanmamış profiller (user_id IS NULL)
// / // - Mevcut kullanıcının profili (user_id = 5)
// / // seçenek olarak gösterilir
// / ```
// /
// / ## AutoOptions ile BelongsTo
// /
// / ```go
// / // Alan tanımı
// / fields.NewBelongsTo("category").
// /     SetRelatedResource("categories").
// /     SetAutoOptions("name")
// /
// / // Tüm kategoriler seçenek olarak gösterilir:
// / // {
// / //   "1": "Electronics",
// / //   "2": "Books",
// / //   "3": "Clothing"
// / // }
// / ```
// /
// / ## Manuel Callback ile
// /
// / ```go
// / // Alan tanımı
// / fields.NewSelect("status").
// /     SetOptions(func() map[string]string {
// /         return map[string]string{
// /             "active": "Aktif",
// /             "inactive": "Pasif",
// /             "pending": "Beklemede",
// /         }
// /     })
// /
// / // Callback çalıştırılır ve sonuç props'a eklenir
// / ```
// /
// / # Seçenek Formatı
// /
// / Çözümlenen seçenekler `map[string]string` formatında props'a eklenir:
// /
// / ```go
// / {
// /   "props": {
// /     "options": {
// /       "1": "Option 1",
// /       "2": "Option 2",
// /       "3": "Option 3"
// /     }
// /   }
// / }
// / ```
// /
// / # Önemli Notlar
// /
// / - **Mevcut Seçenekler**: Props'ta zaten options varsa, AutoOptions çalışmaz
// / - **Debug Çıktıları**: fmt.Printf ile detaylı debug logları içerir
// / - **Null Kontrolü**: HasOne için item nil kontrolü yapılır
// / - **ID Çıkarma**: Item'dan ID çıkarılırken reflection kullanılır
// / - **Callback Tipleri**: İki tip callback desteklenir:
// /   - `func() map[string]string`
// /   - `func() map[string]interface{}`
// /
// / # Performans Notları
// /
// / - Her alan için ayrı veritabanı sorgusu yapılır
// / - Büyük ilişki tablolarında performans sorunu olabilir
// / - Cache mekanizması düşünülebilir
// / - Eager loading ile toplu yükleme tercih edilebilir
// /
// / # Hata Durumları
// /
// / - DB nil ise: Sorgu yapılmaz, boş seçenek listesi döner
// / - Table bulunamazsa: Boş sonuç döner
// / - Display field yoksa: Boş string değerler eklenir
// / - Item ID bulunamazsa: HasOne için tüm NULL kayıtlar gösterilir
// /
// / # İlişkili Metodlar
// /
// / - `resolveResourceFields`: Bu metodu çağırır
// / - `Element.GetAutoOptionsConfig`: AutoOptions yapılandırmasını alır
// / - `Element.JsonSerialize`: Serileştirilmiş veriyi sağlar
func (h *FieldHandler) ResolveFieldOptions(element fields.Element, serialized map[string]interface{}, item interface{}) {
	props, _ := serialized["props"].(map[string]interface{})
	if props == nil {
		props = make(map[string]interface{})
		serialized["props"] = props
	}

	// Handle AutoOptions via Config (works even if element is *Schema due to fluent API)
	config := element.GetAutoOptionsConfig()
	// fmt.Printf("[DEBUG] ResolveFieldOptions - Key: %s, View: %s, Enabled: %v\n", element.GetKey(), element.GetView(), config.Enabled)
	if config.Enabled {
		if _, hasOpts := props["options"]; !hasOpts {
			slug, _ := props["related_resource"].(string)
			display := config.DisplayField

			// Resolve table name from resource registry if possible
			// Slug != Table Name assumption
			tableName := slug
			if res := resource.Get(slug); res != nil {
				// Get GORM DB to parse model
				if db, ok := h.Provider.GetClient().(*gorm.DB); ok {
					stmt := &gorm.Statement{DB: db}
					if err := stmt.Parse(res.Model()); err == nil {
						tableName = stmt.Schema.Table
					}
				}
			}

			// fmt.Printf("[DEBUG] AutoOptions Slug: %s, Table: %s, Display: %s\n", slug, tableName, display)

			if tableName != "" && display != "" {
				var results []map[string]interface{}
				view := element.GetView()

				if view == "has-one-field" {
					fk, _ := props["foreign_key"].(string)
					// fmt.Printf("[DEBUG] HasOne Query - Table: %s, FK: %s\n", tableName, fk)

					// Get GORM DB from provider
					if db, ok := h.Provider.GetClient().(*gorm.DB); ok && fk != "" {
						query := db.Table(tableName).Select("id, " + display)

						var itemID interface{}
						if item != nil {
							val := reflect.ValueOf(item)
							if val.Kind() == reflect.Ptr {
								val = val.Elem()
							}
							if val.Kind() == reflect.Struct {
								// Try to find ID or Id field
								idField := val.FieldByName("ID")
								if !idField.IsValid() {
									idField = val.FieldByName("Id")
								}
								if idField.IsValid() {
									itemID = idField.Interface()
								}
							}
						}

						if itemID != nil {
							query = query.Where(fk+" IS NULL OR "+fk+" = ?", itemID)
						} else {
							query = query.Where(fk + " IS NULL OR " + fk + " = 0")
						}
						query.Find(&results)
					}
				} else if view == "belongs-to-field" || view == "belongs-to-many-field" || view == "has-many-field" || view == "morph-to-many-field" {
					// fmt.Printf("[DEBUG] BelongsTo/HasMany/MorphToMany Query - Table: %s\n", tableName)

					// Get GORM DB from provider
					if db, ok := h.Provider.GetClient().(*gorm.DB); ok {
						db.Table(tableName).Select("id, " + display).Find(&results)
					}
				}

				// fmt.Printf("[DEBUG] Query Result Count: %d\n", len(results))
				opts := make(map[string]string)
				for _, r := range results {
					if val, ok := r[display]; ok {
						opts[fmt.Sprint(r["id"])] = fmt.Sprint(val)
					}
				}
				props["options"] = opts
			}
		}
	}

	// Resolve dynamic options if present (callback)
	if optsFunc, ok := props["options"].(func() map[string]string); ok {
		props["options"] = optsFunc()
	} else if optsFunc, ok := props["options"].(func() map[string]interface{}); ok {
		props["options"] = optsFunc()
	}
}

// / Bu metod, HTTP request'ten gelen veriyi parse eder ve kaynak işlemleri için hazırlar.
// / JSON, multipart/form-data ve standart form verilerini destekler. Dosya yükleme,
// / MorphTo alanları ve BelongsToMany ilişkileri için özel işlemler yapar.
// /
// / # Temel İşlevler
// /
// / 1. **Content-Type Tespiti**: Request'in content-type'ına göre uygun parse stratejisi seçer
// / 2. **JSON Parse**: application/json için body parser kullanır
// / 3. **Multipart Parse**: multipart/form-data için form ve dosya işleme yapar
// / 4. **Dosya Yükleme**: File, Video, Audio alanları için dosya upload işlemi
// / 5. **MorphTo İşleme**: Polymorphic ilişkiler için type ve id ayrıştırma
// / 6. **BelongsToMany Boş Değer**: Unchecked ilişkiler için boş slice ataması
// / 7. **Modify Callbacks**: Her alan için ModifyCallback çalıştırma
// /
// / # Parametreler
// /
// / - `c`: Panel context nesnesi (HTTP request/response wrapper)
// /
// / # Döndürür
// /
// / - `map[string]interface{}`: Parse edilmiş request verisi
// / - `error`: Parse hatası veya nil
// /
// / # Desteklenen Content-Type'lar
// /
// / - `application/json`: JSON body
// / - `multipart/form-data`: Form data + dosya yükleme
// / - `application/x-www-form-urlencoded`: Standart form data
// /
// / # Kullanım Senaryoları
// /
// / 1. **Create İşlemi**: Yeni kayıt oluşturma için form verisi parse etme
// / 2. **Update İşlemi**: Mevcut kayıt güncelleme için form verisi parse etme
// / 3. **Dosya Upload**: Resim, video, audio dosyalarını yükleme
// / 4. **İlişki Yönetimi**: BelongsToMany, MorphTo ilişkilerini işleme
// /
// / # JSON Parse Örneği
// /
// / ```go
// / // Request Body (JSON):
// / // {
// / //   "name": "John Doe",
// / //   "email": "john@example.com",
// / //   "role_id": 2
// / // }
// /
// / body, err := handler.parseBody(ctx)
// / // body = map[string]interface{}{
// / //   "name": "John Doe",
// / //   "email": "john@example.com",
// / //   "role_id": 2,
// / // }
// / ```
// /
// / # Multipart Form Parse Örneği
// /
// / ```go
// / // Request (multipart/form-data):
// / // name: John Doe
// / // email: john@example.com
// / // avatar: [file]
// / // tags[]: 1
// / // tags[]: 2
// / // tags[]: 3
// /
// / body, err := handler.parseBody(ctx)
// / // body = map[string]interface{}{
// / //   "name": "John Doe",
// / //   "email": "john@example.com",
// / //   "avatar": "/storage/avatars/abc123.jpg",
// / //   "tags": []interface{}{"1", "2", "3"},
// / // }
// / ```
// /
// / # Dosya Yükleme İşlemi
// /
// / Dosya alanları için iki yöntem desteklenir:
// /
// / ## 1. Custom Storage Callback
// /
// / ```go
// / fields.NewFile("avatar").
// /     SetStorageCallback(func(c *fiber.Ctx, file *multipart.FileHeader) (string, error) {
// /         // Özel storage mantığı
// /         return uploadToS3(file)
// /     })
// / ```
// /
// / ## 2. Resource StoreHandler (Varsayılan)
// /
// / ```go
// / func (r *UserResource) StoreHandler(ctx *context.Context, file *multipart.FileHeader,
// /     storagePath, storageURL string) (string, error) {
// /     // Varsayılan storage mantığı
// /     return saveToLocal(file, storagePath, storageURL)
// / }
// / ```
// /
// / # MorphTo Alan İşleme
// /
// / MorphTo alanları için özel parse mantığı:
// /
// / ```go
// / // Request Body:
// / // {
// / //   "commentable": "{\"type\":\"posts\",\"id\":\"5\"}"
// / // }
// /
// / // Parse sonrası:
// / // {
// / //   "commentable_type": "posts",
// / //   "commentable_id": "5"
// / // }
// / // (Orijinal "commentable" key'i silinir)
// / ```
// /
// / MorphTo değeri üç formatta gelebilir:
// / 1. **JSON String**: `"{\"type\":\"posts\",\"id\":\"5\"}"`
// / 2. **Map**: `map[string]interface{}{"type": "posts", "id": "5"}`
// / 3. **Ayrık Alanlar**: `commentable_type` ve `commentable_id` ayrı ayrı
// /
// / # BelongsToMany Boş Değer İşleme
// /
// / Multipart form-data'da unchecked checkbox'lar için:
// /
// / ```go
// / // Form'da "tags" alanı gönderilmemişse (tüm checkbox'lar unchecked):
// / // body["tags"] = []interface{}{}
// / // Bu, ilişkilerin tamamen temizlenmesini sağlar
// / ```
// /
// / # Array Alan İşleme
// /
// / Array alanları için özel işlem:
// /
// / ```go
// / // Form field: tags[]
// / // Normalize edilir: tags
// /
// / // Tek değer: tags[] = "1" -> tags = "1"
// / // Çoklu değer: tags[] = ["1", "2", "3"] -> tags = ["1", "2", "3"]
// / ```
// /
// / # Modify Callback İşleme
// /
// / Her alan için ModifyCallback varsa çalıştırılır:
// /
// / ```go
// / fields.NewText("slug").
// /     SetModifyCallback(func(value interface{}, c *fiber.Ctx) interface{} {
// /         // Slug'ı otomatik oluştur
// /         if value == "" || value == nil {
// /             name := c.FormValue("name")
// /             return slugify(name)
// /         }
// /         return value
// /     })
// / ```
// /
// / # Dosya Tipi Tespiti
// /
// / Dosya alanları için özel işlem yapılır:
// / - `TYPE_FILE`: Genel dosyalar
// / - `TYPE_VIDEO`: Video dosyaları
// / - `TYPE_AUDIO`: Audio dosyaları
// /
// / Bu alanlar için form value yerine dosya upload işlemi yapılır.
// /
// / # Önemli Notlar
// /
// / - **Content-Type Kontrolü**: İlk olarak content-type kontrol edilir
// / - **Fiber BodyParser**: JSON ve form data için Fiber'ın built-in parser'ı kullanılır
// / - **Dosya Önceliği**: Dosya alanları form value'larını override eder
// / - **MorphTo Temizleme**: MorphTo parse sonrası orijinal key silinir
// / - **Callback Sırası**: ModifyCallback en son uygulanır
// / - **Error Handling**: Dosya upload hatası tüm işlemi durdurur
// /
// / # Hata Durumları
// /
// / - **BodyParser Hatası**: JSON veya form parse edilemezse error döner
// / - **Dosya Upload Hatası**: Storage callback veya StoreHandler hata verirse error döner
// / - **MorphTo Parse Hatası**: JSON unmarshal hatası sessizce ignore edilir
// /
// / # Performans Notları
// /
// / - Büyük dosyalar için memory kullanımı artabilir
// / - Multipart parse tüm form'u memory'ye yükler
// / - Dosya upload senkron çalışır, async düşünülebilir
// / - ModifyCallback her alan için ayrı çalışır
// /
// / # Güvenlik Notları
// /
// / - Dosya tipi validasyonu yapılmalı
// / - Dosya boyutu limiti kontrol edilmeli
// / - Dosya adı sanitize edilmeli
// / - Path traversal saldırılarına karşı korunmalı
// /
// / # İlişkili Metodlar
// /
// / - `Resource.StoreHandler`: Varsayılan dosya upload mantığı
// / - `Element.GetStorageCallback`: Özel storage callback
// / - `Element.GetModifyCallback`: Veri modifikasyon callback
func (h *FieldHandler) parseBody(c *context.Context) (map[string]interface{}, error) {
	var body = make(map[string]interface{})
	elements := h.getElements(c)

	// Check content type
	ctype := c.Ctx.Get("Content-Type")
	if !strings.Contains(ctype, "multipart/form-data") && !strings.Contains(ctype, "application/json") {
		// Try to parse json body manually if not standard
		if err := c.Ctx.BodyParser(&body); err != nil {
			return nil, err
		}
		return body, nil
	}

	// Handle JSON Body
	if err := c.Ctx.BodyParser(&body); err == nil {
		// If success, we have data. But if it's multipart, we might have partial data in body?
		// Fiber BodyParser handles multipart form fields too.
	}

	// Handle Form Data (Multipart)
	if form, err := c.Ctx.MultipartForm(); err == nil {
		for key, values := range form.Value {
			if len(values) > 0 {
				normalizedKey := key
				if strings.HasSuffix(key, "[]") {
					normalizedKey = strings.TrimSuffix(key, "[]")
				}

				var isFileType bool
				for _, el := range elements {
					if el.GetKey() == normalizedKey {
						if el.JsonSerialize()["type"] == fields.TYPE_FILE ||
							el.JsonSerialize()["type"] == fields.TYPE_VIDEO ||
							el.JsonSerialize()["type"] == fields.TYPE_AUDIO {
							isFileType = true
						}
						break
					}
				}

				if !isFileType {
					if strings.HasSuffix(key, "[]") || len(values) > 1 {
						body[normalizedKey] = values
					} else {
						body[normalizedKey] = values[0]
					}
				}
			}
		}
		for key, files := range form.File {
			if len(files) > 0 {
				file := files[0]
				var path string

				// Check for matching element and callback
				var callback fields.StorageCallbackFunc
				for _, el := range elements {
					if el.GetKey() == key {
						callback = el.GetStorageCallback()
						break
					}
				}

				if callback != nil {
					// User defined storage logic (Assuming it takes *fiber.Ctx)
					path, err = callback(c.Ctx, file)
					if err != nil {
						return nil, err
					}
				} else {
					// Use Resource StoreHandler (Takes *context.Context)
					path, err = h.Resource.StoreHandler(c, file, h.StoragePath, h.StorageURL)
					if err != nil {
						return nil, err
					}
				}

				body[key] = path
			}
		}

		// Handle missing BelongsToMany fields in multipart/form-data
		// If a BelongsToMany field is missing from the request, it implies the user unchecked all options.
		// We set it to an empty slice so GormDataProvider clears the relationships.
		for _, el := range elements {
			if _, ok := fields.IsRelationshipField(el); ok {
				// We specifically check for BelongsToMany by view or type if interface allows
				// Using reflection or type check if possible, or View/Type string
				if strings.HasPrefix(el.GetView(), "belongs-to-many-field") {
					key := el.GetKey()
					if _, exists := body[key]; !exists {
						body[key] = []interface{}{}
					}
				}
			}
		}
	}

	// Convert explicit frontend null sentinel values to nil for relationship fields.
	for _, el := range elements {
		key := el.GetKey()
		rawVal, exists := body[key]
		if !exists {
			continue
		}

		strVal, ok := rawVal.(string)
		if !ok || strVal != formNullSentinel {
			continue
		}

		// MorphTo is handled below by splitting into *_type and *_id fields.
		if strings.HasPrefix(el.GetView(), "morph-to-field") {
			continue
		}

		if _, isRelationship := fields.IsRelationshipField(el); isRelationship {
			body[key] = nil
		}
	}

	// Apply ModifyCallback for all fields
	for _, el := range elements {
		if val, ok := body[el.GetKey()]; ok {
			if callback := el.GetModifyCallback(); callback != nil {
				body[el.GetKey()] = callback(val, c.Ctx)
			}
		}
	}

	// Handle MorphTo fields: parse JSON object {"type":"...", "id":"..."} into separate fields
	for _, el := range elements {
		if strings.HasPrefix(el.GetView(), "morph-to-field") {
			key := el.GetKey()
			typeKey := key + "_type"
			idKey := key + "_id"

			if val, ok := body[key]; ok && val != nil {
				// Parse MorphTo value - can be JSON string, map, or already separated
				switch v := val.(type) {
				case string:
					if v == formNullSentinel {
						body[typeKey] = nil
						body[idKey] = nil
						delete(body, key)
						continue
					}
					// JSON string from form-data: {"type":"posts","id":"1"}
					if strings.HasPrefix(v, "{") {
						var morphData map[string]interface{}
						if err := json.Unmarshal([]byte(v), &morphData); err == nil {
							if morphType, ok := morphData["type"].(string); ok && morphType != "" {
								body[typeKey] = morphType
							}
							if morphID, exists := morphData["id"]; exists {
								body[idKey] = morphID
							}
						}
					}
				case map[string]interface{}:
					// Already parsed JSON object
					if morphType, ok := v["type"].(string); ok && morphType != "" {
						body[typeKey] = morphType
					}
					if morphID, exists := v["id"]; exists {
						body[idKey] = morphID
					}
				}
				// Remove the original composite key
				delete(body, key)
			}
		}
	}

	return body, nil
}

// / Bu metod, yeni kaynak oluşturma (store) endpoint'ini işler.
// / İşlemi resource_store_controller.go dosyasındaki HandleResourceStore fonksiyonuna yönlendirir.
// /
// / # Parametreler
// /
// / - `c`: Panel context nesnesi (HTTP request/response wrapper)
// /
// / # Döndürür
// /
// / - `error`: İşlem hatası veya nil
// /
// / # Kullanım Senaryoları
// /
// / - Yeni kayıt oluşturma (POST request)
// / - Form verisi parse etme ve validasyon
// / - Dosya yükleme işlemleri
// / - İlişki (relationship) kurma
// /
// / # Örnek Kullanım
// /
// / ```go
// / app.Post("/api/users", func(c *fiber.Ctx) error {
// /     ctx := context.New(c)
// /     return handler.Store(ctx)
// / })
// / ```
func (h *FieldHandler) Store(c *context.Context) error {
	return HandleResourceStore(h, c)
}

// / Bu metod, mevcut kaynağı güncelleme (update) endpoint'ini işler.
// / İşlemi resource_update_controller.go dosyasındaki HandleResourceUpdate fonksiyonuna yönlendirir.
// /
// / # Parametreler
// /
// / - `c`: Panel context nesnesi (HTTP request/response wrapper)
// /
// / # Döndürür
// /
// / - `error`: İşlem hatası veya nil
// /
// / # Kullanım Senaryoları
// /
// / - Mevcut kaydı güncelleme (PUT/PATCH request)
// / - Kısmi güncelleme (partial update)
// / - Dosya değiştirme işlemleri
// / - İlişki güncelleme
// /
// / # Örnek Kullanım
// /
// / ```go
// / app.Put("/api/users/:id", func(c *fiber.Ctx) error {
// /     ctx := context.New(c)
// /     return handler.Update(ctx)
// / })
// / ```
func (h *FieldHandler) Update(c *context.Context) error {
	return HandleResourceUpdate(h, c)
}

// / Bu metod, kaynak silme (destroy) endpoint'ini işler.
// / İşlemi resource_destroy_controller.go dosyasındaki HandleResourceDestroy fonksiyonuna yönlendirir.
// /
// / # Parametreler
// /
// / - `c`: Panel context nesnesi (HTTP request/response wrapper)
// /
// / # Döndürür
// /
// / - `error`: İşlem hatası veya nil
// /
// / # Kullanım Senaryoları
// /
// / - Kayıt silme (DELETE request)
// / - Soft delete işlemleri
// / - İlişkili verileri temizleme
// / - Cascade delete işlemleri
// /
// / # Örnek Kullanım
// /
// / ```go
// / app.Delete("/api/users/:id", func(c *fiber.Ctx) error {
// /     ctx := context.New(c)
// /     return handler.Destroy(ctx)
// / })
// / ```
func (h *FieldHandler) Destroy(c *context.Context) error {
	return HandleResourceDestroy(h, c)
}

// / Bu metod, yeni kaynak oluşturma formu endpoint'ini işler.
// / İşlemi resource_create_controller.go dosyasındaki HandleResourceCreate fonksiyonuna yönlendirir.
// /
// / # Parametreler
// /
// / - `c`: Panel context nesnesi (HTTP request/response wrapper)
// /
// / # Döndürür
// /
// / - `error`: İşlem hatası veya nil
// /
// / # Kullanım Senaryoları
// /
// / - Oluşturma formu için alan listesi sağlama
// / - Varsayılan değerleri ayarlama
// / - Form validasyon kurallarını sağlama
// / - Seçenek listelerini (options) yükleme
// /
// / # Örnek Kullanım
// /
// / ```go
// / app.Get("/api/users/create", func(c *fiber.Ctx) error {
// /     ctx := context.New(c)
// /     return handler.Create(ctx)
// / })
// / ```
func (h *FieldHandler) Create(c *context.Context) error {
	return HandleResourceCreate(h, c)
}

// / Bu metod, dashboard kartları (cards) listesi endpoint'ini işler.
// / İşlemi card_controller.go dosyasındaki HandleCardList fonksiyonuna yönlendirir.
// /
// / # Parametreler
// /
// / - `c`: Panel context nesnesi (HTTP request/response wrapper)
// /
// / # Döndürür
// /
// / - `error`: İşlem hatası veya nil
// /
// / # Kullanım Senaryoları
// /
// / - Dashboard widget'larını listeleme
// / - İstatistik kartlarını gösterme
// / - Özet bilgileri sağlama
// /
// / # Örnek Kullanım
// /
// / ```go
// / app.Get("/api/users/cards", func(c *fiber.Ctx) error {
// /     ctx := context.New(c)
// /     return handler.ListCards(ctx)
// / })
// / ```
func (h *FieldHandler) ListCards(c *context.Context) error {
	return HandleCardList(h, c)
}

// / Bu metod, tek bir dashboard kartı (card) detayı endpoint'ini işler.
// / İşlemi card_detail_controller.go dosyasındaki HandleCardDetail fonksiyonuna yönlendirir.
// /
// / # Parametreler
// /
// / - `c`: Panel context nesnesi (HTTP request/response wrapper)
// /
// / # Döndürür
// /
// / - `error`: İşlem hatası veya nil
// /
// / # Kullanım Senaryoları
// /
// / - Belirli bir widget'ın detaylı verisini getirme
// / - Kart verilerini yenileme (refresh)
// / - Dinamik kart içeriği sağlama
// /
// / # Örnek Kullanım
// /
// / ```go
// / app.Get("/api/users/cards/:cardKey", func(c *fiber.Ctx) error {
// /     ctx := context.New(c)
// /     return handler.GetCard(ctx)
// / })
// / ```
func (h *FieldHandler) GetCard(c *context.Context) error {
	return HandleCardDetail(h, c)
}

// / Bu metod, alan listesi (field list) endpoint'ini işler.
// / İşlemi field_controller.go dosyasındaki HandleFieldList fonksiyonuna yönlendirir.
// /
// / # Parametreler
// /
// / - `c`: Panel context nesnesi (HTTP request/response wrapper)
// /
// / # Döndürür
// /
// / - `error`: İşlem hatası veya nil
// /
// / # Kullanım Senaryoları
// /
// / - Resource için tanımlı alanları listeleme
// / - Frontend için alan metadata'sı sağlama
// / - Dinamik form oluşturma için alan bilgisi
// /
// / # Örnek Kullanım
// /
// / ```go
// / app.Get("/api/users/fields", func(c *fiber.Ctx) error {
// /     ctx := context.New(c)
// /     return handler.List(ctx)
// / })
// / ```
func (h *FieldHandler) List(c *context.Context) error {
	return HandleFieldList(h, c)
}

// / Bu fonksiyon, ResourceContext oluşturan ve HTTP request'e enjekte eden bir Fiber middleware'i döndürür.
// / Middleware, kaynak metadata'sı, görünürlük context'i ve alan çözümleyicilerini başlatır.
// /
// / # Temel İşlevler
// /
// / 1. **Context Oluşturma**: core.NewResourceContextWithVisibility ile tam yapılandırılmış context oluşturur
// / 2. **Metadata Başlatma**: Resource ve Lens bilgilerini context'e ekler
// / 3. **Görünürlük Kontrolü**: VisibilityContext ile alan görünürlük kurallarını ayarlar
// / 4. **Context Enjeksiyonu**: Oluşturulan context'i fiber.Ctx.Locals'a ekler
// / 5. **Chain Devamı**: c.Next() ile middleware chain'ini devam ettirir
// /
// / # Parametreler
// /
// / - `resource`: İşlenecek kaynak (Resource interface implementasyonu)
// / - `lens`: Opsiyonel lens (filtrelenmiş görünüm), nil olabilir
// / - `visibilityCtx`: Alanların hangi context'te görünür olacağını belirler
// /   (örn: INDEX, SHOW, EDIT, CREATE)
// / - `elements`: Bu kaynak ile ilişkili alan listesi
// /
// / # Döndürür
// /
// / - `fiber.Handler`: ResourceContext'i enjekte eden middleware fonksiyonu
// /
// / # Kullanım Senaryoları
// /
// / 1. **Index Endpoint**: Liste görünümü için context hazırlama
// / 2. **Create Endpoint**: Yeni kayıt formu için context hazırlama
// / 3. **Lens Endpoint**: Filtrelenmiş görünüm için context hazırlama
// / 4. **Field List**: Alan listesi endpoint'i için context hazırlama
// /
// / # Örnek Kullanım
// /
// / ## Standart Resource Endpoint
// /
// / ```go
// / // Resource tanımı
// / userResource := &UserResource{}
// / elements := userResource.Fields()
// /
// / // Index endpoint için middleware
// / app.Get("/api/users",
// /     FieldContextMiddleware(
// /         userResource,
// /         nil, // lens yok
// /         core.INDEX, // liste görünümü
// /         elements,
// /     ),
// /     func(c *fiber.Ctx) error {
// /         ctx := context.New(c)
// /         return handler.Index(ctx)
// /     },
// / )
// / ```
// /
// / ## Lens ile Kullanım
// /
// / ```go
// / // Lens tanımı
// / activeLens := &ActiveUsersLens{}
// / lensElements := activeLens.Fields()
// /
// / // Lens endpoint için middleware
// / app.Get("/api/users/active",
// /     FieldContextMiddleware(
// /         userResource,
// /         activeLens, // lens var
// /         core.INDEX,
// /         lensElements,
// /     ),
// /     func(c *fiber.Ctx) error {
// /         ctx := context.New(c)
// /         return handler.Index(ctx)
// /     },
// / )
// / ```
// /
// / ## Create Form Endpoint
// /
// / ```go
// / // Create form için middleware
// / app.Get("/api/users/create",
// /     FieldContextMiddleware(
// /         userResource,
// /         nil,
// /         core.CREATE, // oluşturma formu
// /         elements,
// /     ),
// /     func(c *fiber.Ctx) error {
// /         ctx := context.New(c)
// /         return handler.Create(ctx)
// /     },
// / )
// / ```
// /
// / # Context Erişimi
// /
// / Downstream handler'lar context'e şu şekilde erişir:
// /
// / ```go
// / func MyHandler(c *fiber.Ctx) error {
// /     // Context'i al
// /     ctx := c.Locals(core.ResourceContextKey).(*core.ResourceContext)
// /
// /     // Context bilgilerini kullan
// /     resource := ctx.Resource()
// /     lens := ctx.Lens()
// /     visibilityCtx := ctx.VisibilityContext()
// /
// /     // Alanları filtrele
// /     for _, element := range ctx.Elements() {
// /         if element.IsVisible(ctx) {
// /             // Alan görünür, işle
// /         }
// /     }
// /
// /     return nil
// / }
// / ```
// /
// / # VisibilityContext Değerleri
// /
// / - `core.INDEX`: Liste görünümü (tablo)
// / - `core.SHOW`: Detay görünümü (read-only)
// / - `core.EDIT`: Düzenleme formu
// / - `core.CREATE`: Oluşturma formu
// / - `core.DETAIL`: Detay sayfası (özelleştirilmiş)
// /
// / # Önemli Notlar
// /
// / - **Item ve User Null**: Bu middleware item ve user parametrelerini nil olarak ayarlar,
// /   bunlar handler'lar tarafından daha sonra ayarlanmalıdır
// / - **Middleware Sırası**: Bu middleware, authentication ve authorization middleware'lerinden
// /   sonra çalıştırılmalıdır
// / - **Context Key**: Context, `core.ResourceContextKey` ile Locals'a kaydedilir
// / - **Chain Devamı**: Middleware her zaman c.Next() çağırır, hata durumunda bile
// /
// / # Requirement Karşılama
// /
// / - **Requirement 15.1**: ResourceContext'i oluşturmak için middleware güncellendi
// / - **Requirement 15.4**: Context oluşturulduğunda tüm gerekli kaynak bilgisi başlatılır
// /
// / # İlişkili Fonksiyonlar
// /
// / - `FieldContextMiddlewareWithItem`: Item ve user ile context oluşturur
// / - `core.NewResourceContextWithVisibility`: Context oluşturma fonksiyonu
// / - `core.ResourceContext`: Context yapısı
func FieldContextMiddleware(resource interface{}, lens interface{}, visibilityCtx core.VisibilityContext, elements []fields.Element) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Create ResourceContext with full visibility and metadata
		ctx := core.NewResourceContextWithVisibility(
			c,
			resource,
			lens,
			visibilityCtx,
			nil, // Item will be set by handlers if needed
			nil, // User will be set by handlers if needed
			elements,
		)
		c.Locals(core.ResourceContextKey, ctx)
		return c.Next()
	}
}

// / Bu fonksiyon, belirli bir kaynak öğesi (item) ve kullanıcı (user) ile ResourceContext oluşturan
// / ve HTTP request'e enjekte eden bir Fiber middleware'i döndürür. Öğe-spesifik işlemler için kullanılır.
// /
// / # Temel İşlevler
// /
// / 1. **Item Context**: Belirli bir kaynak öğesi ile context oluşturur
// / 2. **User Context**: İşlemi yapan kullanıcı bilgisini context'e ekler
// / 3. **Görünürlük Kontrolü**: Item ve user bazlı alan görünürlük kurallarını ayarlar
// / 4. **Metadata Başlatma**: Resource, Lens, Item ve User bilgilerini context'e ekler
// / 5. **Context Enjeksiyonu**: Oluşturulan context'i fiber.Ctx.Locals'a ekler
// / 6. **Chain Devamı**: c.Next() ile middleware chain'ini devam ettirir
// /
// / # Parametreler
// /
// / - `resource`: İşlenecek kaynak (Resource interface implementasyonu)
// / - `lens`: Opsiyonel lens (filtrelenmiş görünüm), nil olabilir
// / - `visibilityCtx`: Alanların hangi context'te görünür olacağını belirler
// /   (örn: INDEX, SHOW, EDIT, CREATE)
// / - `item`: İşlenecek spesifik kaynak öğesi (model instance)
// / - `user`: İşlemi yapan kullanıcı
// / - `elements`: Bu kaynak ile ilişkili alan listesi
// /
// / # Döndürür
// /
// / - `fiber.Handler`: Item ve user ile yapılandırılmış ResourceContext'i enjekte eden middleware fonksiyonu
// /
// / # Kullanım Senaryoları
// /
// / 1. **Show Endpoint**: Belirli bir kaydı görüntüleme
// / 2. **Edit Endpoint**: Belirli bir kaydı düzenleme formu
// / 3. **Detail Endpoint**: Belirli bir kaydın detay sayfası
// / 4. **Update Endpoint**: Belirli bir kaydı güncelleme
// / 5. **Delete Endpoint**: Belirli bir kaydı silme
// /
// / # Örnek Kullanım
// /
// / ## Show Endpoint
// /
// / ```go
// / // Resource tanımı
// / userResource := &UserResource{}
// / elements := userResource.Fields()
// /
// / // Show endpoint için middleware
// / app.Get("/api/users/:id", func(c *fiber.Ctx) error {
// /     // Önce item'ı veritabanından çek
// /     id := c.Params("id")
// /     var user models.User
// /     if err := db.First(&user, id).Error; err != nil {
// /         return err
// /     }
// /
// /     // Mevcut kullanıcıyı al (authentication middleware'den)
// /     currentUser := c.Locals("user").(*models.User)
// /
// /     // Middleware'i dinamik olarak uygula
// /     middleware := FieldContextMiddlewareWithItem(
// /         userResource,
// /         nil,
// /         core.SHOW,
// /         &user,
// /         currentUser,
// /         elements,
// /     )
// /
// /     // Middleware'i çalıştır
// /     if err := middleware(c); err != nil {
// /         return err
// /     }
// /
// /     // Handler'ı çağır
// /     ctx := context.New(c)
// /     return handler.Show(ctx)
// / })
// / ```
// /
// / ## Edit Endpoint
// /
// / ```go
// / // Edit form için middleware
// / app.Get("/api/users/:id/edit", func(c *fiber.Ctx) error {
// /     id := c.Params("id")
// /     var user models.User
// /     db.First(&user, id)
// /
// /     currentUser := c.Locals("user").(*models.User)
// /
// /     middleware := FieldContextMiddlewareWithItem(
// /         userResource,
// /         nil,
// /         core.EDIT, // düzenleme formu
// /         &user,
// /         currentUser,
// /         elements,
// /     )
// /
// /     middleware(c)
// /     ctx := context.New(c)
// /     return handler.Edit(ctx)
// / })
// / ```
// /
// / ## Koşullu Alan Görünürlüğü
// /
// / Item ve user context'i ile koşullu alan görünürlüğü:
// /
// / ```go
// / // Alan tanımı
// / fields.NewText("salary").
// /     SetVisibleCallback(func(ctx *core.ResourceContext) bool {
// /         // Sadece admin veya kendi kaydını görüntüleyen kullanıcı görebilir
// /         user := ctx.User().(*models.User)
// /         item := ctx.Item().(*models.User)
// /
// /         return user.IsAdmin() || user.ID == item.ID
// /     })
// / ```
// /
// / # Context Erişimi
// /
// / Downstream handler'lar item ve user'a şu şekilde erişir:
// /
// / ```go
// / func MyHandler(c *fiber.Ctx) error {
// /     // Context'i al
// /     ctx := c.Locals(core.ResourceContextKey).(*core.ResourceContext)
// /
// /     // Item'a eriş
// /     user := ctx.Item().(*models.User)
// /     fmt.Printf("Editing user: %s\n", user.Name)
// /
// /     // User'a eriş
// /     currentUser := ctx.User().(*models.User)
// /     fmt.Printf("Current user: %s\n", currentUser.Name)
// /
// /     // Yetki kontrolü
// /     if !currentUser.CanEdit(user) {
// /         return fiber.ErrForbidden
// /     }
// /
// /     // Alanları filtrele
// /     for _, element := range ctx.Elements() {
// /         if element.IsVisible(ctx) {
// /             // Alan görünür, işle
// /         }
// /     }
// /
// /     return nil
// / }
// / ```
// /
// / # FieldContextMiddleware ile Farklar
// /
// / | Özellik | FieldContextMiddleware | FieldContextMiddlewareWithItem |
// / |---------|------------------------|--------------------------------|
// / | Item | nil | Belirli bir kayıt |
// / | User | nil | Belirli bir kullanıcı |
// / | Kullanım | Liste, Create | Show, Edit, Update, Delete |
// / | Görünürlük | Genel kurallar | Item/User bazlı kurallar |
// / | Yetki Kontrolü | Genel | Öğe-spesifik |
// /
// / # Item ve User Bazlı Özellikler
// /
// / ## 1. Koşullu Alan Görünürlüğü
// /
// / ```go
// / fields.NewText("private_notes").
// /     SetVisibleCallback(func(ctx *core.ResourceContext) bool {
// /         user := ctx.User().(*models.User)
// /         return user.IsAdmin()
// /     })
// / ```
// /
// / ## 2. Dinamik Alan Değerleri
// /
// / ```go
// / fields.NewText("status").
// /     SetResolveCallback(func(value interface{}, item interface{}, c *fiber.Ctx) interface{} {
// /         user := c.Locals("user").(*models.User)
// /         if user.IsAdmin() {
// /             return value // Admin tüm durumları görebilir
// /         }
// /         return "hidden" // Normal kullanıcılar göremez
// /     })
// / ```
// /
// / ## 3. Öğe-Spesifik Validasyon
// /
// / ```go
// / fields.NewText("email").
// /     SetModifyCallback(func(value interface{}, c *fiber.Ctx) interface{} {
// /         ctx := c.Locals(core.ResourceContextKey).(*core.ResourceContext)
// /         item := ctx.Item().(*models.User)
// /         user := ctx.User().(*models.User)
// /
// /         // Sadece admin veya kendi email'ini değiştirebilir
// /         if !user.IsAdmin() && user.ID != item.ID {
// /             return item.Email // Değişikliği engelle
// /         }
// /         return value
// /     })
// / ```
// /
// / # Önemli Notlar
// /
// / - **Item Null Kontrolü**: Item nil olabilir, handler'larda kontrol edilmeli
// / - **User Null Kontrolü**: User nil olabilir, authentication middleware gerekli
// / - **Type Assertion**: Item ve User type assertion ile dönüştürülmeli
// / - **Middleware Sırası**: Authentication middleware'den sonra çalıştırılmalı
// / - **Context Key**: Context, `core.ResourceContextKey` ile Locals'a kaydedilir
// / - **Chain Devamı**: Middleware her zaman c.Next() çağırır
// /
// / # Güvenlik Notları
// /
// / - Item'a erişim yetkisi kontrol edilmeli (Policy.Can())
// / - User authentication doğrulanmalı
// / - Cross-user data access önlenmeli
// / - Sensitive alanlar için ekstra kontroller yapılmalı
// /
// / # Performans Notları
// /
// / - Item veritabanından önceden yüklenmiş olmalı
// / - Eager loading ile ilişkiler yüklenebilir
// / - Context oluşturma lightweight bir işlemdir
// / - Middleware overhead minimal düzeydedir
// /
// / # Requirement Karşılama
// /
// / - **Requirement 15.1**: ResourceContext'i oluşturmak için middleware güncellendi
// / - **Requirement 15.4**: Context oluşturulduğunda tüm gerekli kaynak bilgisi başlatılır
// /
// / # İlişkili Fonksiyonlar
// /
// / - `FieldContextMiddleware`: Item ve user olmadan context oluşturur
// / - `core.NewResourceContextWithVisibility`: Context oluşturma fonksiyonu
// / - `core.ResourceContext`: Context yapısı
func FieldContextMiddlewareWithItem(resource interface{}, lens interface{}, visibilityCtx core.VisibilityContext, item interface{}, user interface{}, elements []fields.Element) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Create ResourceContext with full visibility and metadata
		ctx := core.NewResourceContextWithVisibility(
			c,
			resource,
			lens,
			visibilityCtx,
			item,
			user,
			elements,
		)
		c.Locals(core.ResourceContextKey, ctx)
		return c.Next()
	}
}
