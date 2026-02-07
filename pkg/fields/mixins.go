// Package fields, alan davranışlarını birleştirilebilir hale getiren mixin'leri sağlar.
//
// Bu paket, alanlara modüler işlevsellik eklemek için mixin pattern (kompozisyon) uygular.
// Her mixin, alan türlerine gömülebilen belirli bir yetenek sağlar.
//
// # Mevcut Mixin'ler
//
// - Searchable: Alanların özel arama mantığı ile aranabilir olmasını sağlar
// - Sortable: Alanların özel sıralama mantığı ile sıralanabilir olmasını sağlar
// - Filterable: Alanların özel filtreleme mantığı ile filtrelenebilir olmasını sağlar
// - Validatable: Alanların doğrulama kuralları ve özel doğrulayıcılara sahip olmasını sağlar
// - Displayable: Alanların görüntüleme formatını özelleştirmesini sağlar
// - Hideable: Alanların farklı bağlamlarda görünürlüğünü kontrol etmesini sağlar
//
// # Kullanım Örneği
//
//	type CustomField struct {
//	    fields.Base
//	    fields.Searchable
//	    fields.Sortable
//	}
//
//	field := &CustomField{}
//	field.SetSearchableColumns([]string{"name", "email"})
//	field.SetSortable(true)
//
// # Mimari
//
// Mixin'ler, Go Panel API mimarisinde önerilen kompozisyon pattern'ini takip eder.
// Kalıtım olmadan yeniden kullanılabilir işlevsellik sağlarlar, bu da tek bir alan türünde
// birden fazla yeteneği birleştirmeyi kolaylaştırır.
//
// Daha fazla bilgi için docs/Fields.md ve .docs/ARCHITECTURE.md dosyalarına bakın.
package fields

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"gorm.io/gorm"
)

// Searchable, alanların aranabilir olmasını sağlayan bir mixin'dir.
//
// Aranabilir sütunlar ve özel arama callback'leri için yapılandırma sağlar.
// Bu callback'ler, özel arama mantığı uygulamak için veritabanı sorgusunu değiştirebilir.
//
// # Kullanım
//
//	type TextField struct {
//	    fields.Base
//	    fields.Searchable
//	}
//
//	field := &TextField{}
//	field.SetSearchableColumns([]string{"name", "email"})
//	field.SetSearchCallback(func(db *gorm.DB, term string) *gorm.DB {
//	    return db.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(term)+"%")
//	})
//
// Daha fazla örnek için docs/Fields.md dosyasına bakın.
type Searchable struct {
	searchableColumns []string
	searchCallback    func(*gorm.DB, string) *gorm.DB
}

// SetSearchableColumns, aranabilir sütunları ayarlar.
// Bu sütunlar global arama işlemlerinde kullanılır.
//
// Parametreler:
//   - columns: Aranabilir olması gereken sütun adlarının listesi
//
// Örnek:
//
//	field.SetSearchableColumns([]string{"name", "email", "username"})
func (s *Searchable) SetSearchableColumns(columns []string) {
	s.searchableColumns = columns
}

// GetSearchableColumns, aranabilir sütunları döndürür.
// Hiçbir sütun yapılandırılmamışsa boş bir dilim döndürür.
//
// Döndürür:
//   - Aranabilir sütun adlarının dilimi
func (s *Searchable) GetSearchableColumns() []string {
	if s.searchableColumns == nil {
		return []string{}
	}
	return s.searchableColumns
}

// SetSearchCallback, özel bir arama callback fonksiyonu ayarlar.
//
// Callback, GORM veritabanı örneğini ve arama terimini alır,
// ve arama sorgusu uygulanmış değiştirilmiş bir veritabanı örneği döndürmelidir.
//
// Parametreler:
//   - cb: Özel arama mantığı için veritabanı sorgusunu değiştiren fonksiyon
//
// Örnek:
//
//	field.SetSearchCallback(func(db *gorm.DB, term string) *gorm.DB {
//	    return db.Where(
//	        "LOWER(name) LIKE ? OR LOWER(email) LIKE ?",
//	        "%"+strings.ToLower(term)+"%",
//	        "%"+strings.ToLower(term)+"%",
//	    )
//	})
func (s *Searchable) SetSearchCallback(cb func(*gorm.DB, string) *gorm.DB) {
	s.searchCallback = cb
}

// GetSearchCallback, arama callback fonksiyonunu döndürür.
// Hiçbir callback ayarlanmamışsa nil döndürür.
//
// Döndürür:
//   - Arama callback fonksiyonu veya ayarlanmamışsa nil
func (s *Searchable) GetSearchCallback() func(*gorm.DB, string) *gorm.DB {
	return s.searchCallback
}

// Sortable, alanların sıralanabilir olmasını sağlayan bir mixin'dir.
//
// Sıralanabilir davranış, varsayılan sıralama yönü ve özel sıralama callback'leri
// için yapılandırma sağlar. Bu callback'ler veritabanı sorgusunu değiştirebilir.
//
// # Kullanım
//
//	type TextField struct {
//	    fields.Base
//	    fields.Sortable
//	}
//
//	field := &TextField{}
//	field.SetSortable(true)
//	field.SetSortDirection("asc")
//	field.SetSortCallback(func(db *gorm.DB, direction string) *gorm.DB {
//	    return db.Order("LOWER(name) " + direction)
//	})
//
// Daha fazla örnek için docs/Fields.md dosyasına bakın.
type Sortable struct {
	sortable      bool
	sortDirection string
	sortCallback  func(*gorm.DB, string) *gorm.DB
}

// SetSortable, alanın sıralanabilir olup olmadığını ayarlar.
//
// Parametreler:
//   - sortable: Sıralamayı etkinleştirmek için true, devre dışı bırakmak için false
//
// Örnek:
//
//	field.SetSortable(true)
func (s *Sortable) SetSortable(sortable bool) {
	s.sortable = sortable
}

// IsSortable, alanın sıralanabilir olup olmadığını döndürür.
//
// Döndürür:
//   - Alan sıralanabilirse true, değilse false
func (s *Sortable) IsSortable() bool {
	return s.sortable
}

// SetSortDirection, varsayılan sıralama yönünü ayarlar.
//
// Parametreler:
//   - direction: Artan sıralama için "asc" veya azalan sıralama için "desc"
//
// Örnek:
//
//	field.SetSortDirection("desc")
func (s *Sortable) SetSortDirection(direction string) {
	s.sortDirection = direction
}

// GetSortDirection, varsayılan sıralama yönünü döndürür.
// Hiçbir yön ayarlanmamışsa "asc" döndürür.
//
// Döndürür:
//   - Sıralama yönü ("asc" veya "desc")
func (s *Sortable) GetSortDirection() string {
	if s.sortDirection == "" {
		return "asc"
	}
	return s.sortDirection
}

// SetSortCallback, özel bir sıralama callback fonksiyonu ayarlar.
//
// Callback, GORM veritabanı örneğini ve sıralama yönünü alır,
// ve sıralama sorgusu uygulanmış değiştirilmiş bir veritabanı örneği döndürmelidir.
//
// Parametreler:
//   - cb: Özel sıralama mantığı için veritabanı sorgusunu değiştiren fonksiyon
//
// Örnek:
//
//	field.SetSortCallback(func(db *gorm.DB, direction string) *gorm.DB {
//	    return db.Order("LOWER(name) " + direction)
//	})
func (s *Sortable) SetSortCallback(cb func(*gorm.DB, string) *gorm.DB) {
	s.sortCallback = cb
}

// GetSortCallback, sıralama callback fonksiyonunu döndürür.
// Hiçbir callback ayarlanmamışsa nil döndürür.
//
// Döndürür:
//   - Sıralama callback fonksiyonu veya ayarlanmamışsa nil
func (s *Sortable) GetSortCallback() func(*gorm.DB, string) *gorm.DB {
	return s.sortCallback
}

// Filterable, alanların filtrelenebilir olmasını sağlayan bir mixin'dir.
//
// Filtrelenebilir davranış, filtre seçenekleri ve özel filtre callback'leri
// için yapılandırma sağlar. Bu callback'ler veritabanı sorgusunu değiştirebilir.
//
// # Kullanım
//
//	type SelectField struct {
//	    fields.Base
//	    fields.Filterable
//	}
//
//	field := &SelectField{}
//	field.SetFilterable(true)
//	field.SetFilterOptions(map[string]string{
//	    "active": "Aktif",
//	    "inactive": "İnaktif",
//	})
//	field.SetFilterCallback(func(db *gorm.DB, value any) *gorm.DB {
//	    return db.Where("status = ?", value)
//	})
//
// Daha fazla örnek için docs/Fields.md dosyasına bakın.
type Filterable struct {
	filterable     bool
	filterCallback func(*gorm.DB, any) *gorm.DB
	filterOptions  map[string]string
}

// SetFilterable, alanın filtrelenebilir olup olmadığını ayarlar.
//
// Parametreler:
//   - filterable: Filtrelemeyi etkinleştirmek için true, devre dışı bırakmak için false
//
// Örnek:
//
//	field.SetFilterable(true)
func (f *Filterable) SetFilterable(filterable bool) {
	f.filterable = filterable
}

// IsFilterable, alanın filtrelenebilir olup olmadığını döndürür.
//
// Döndürür:
//   - Alan filtrelenebilirse true, değilse false
func (f *Filterable) IsFilterable() bool {
	return f.filterable
}

// SetFilterCallback, özel bir filtre callback fonksiyonu ayarlar.
//
// Callback, GORM veritabanı örneğini ve filtre değerini alır,
// ve filtre sorgusu uygulanmış değiştirilmiş bir veritabanı örneği döndürmelidir.
//
// Parametreler:
//   - cb: Özel filtre mantığı için veritabanı sorgusunu değiştiren fonksiyon
//
// Örnek:
//
//	field.SetFilterCallback(func(db *gorm.DB, value any) *gorm.DB {
//	    status, ok := value.(string)
//	    if !ok {
//	        return db
//	    }
//	    return db.Where("status = ?", status)
//	})
func (f *Filterable) SetFilterCallback(cb func(*gorm.DB, any) *gorm.DB) {
	f.filterCallback = cb
}

// GetFilterCallback, filtre callback fonksiyonunu döndürür.
// Hiçbir callback ayarlanmamışsa nil döndürür.
//
// Döndürür:
//   - Filtre callback fonksiyonu veya ayarlanmamışsa nil
func (f *Filterable) GetFilterCallback() func(*gorm.DB, any) *gorm.DB {
	return f.filterCallback
}

// SetFilterOptions, kullanılabilir filtre seçeneklerini ayarlar.
//
// Seçenekler, kullanıcıların önceden tanımlanmış değerlere göre filtreleme yapmasına
// izin veren bir açılır liste veya onay kutusu listesi olarak UI'da gösterilir.
//
// Parametreler:
//   - options: Filtre seçenekleri için değer-etiket eşlemesi
//
// Örnek:
//
//	field.SetFilterOptions(map[string]string{
//	    "draft": "Taslak",
//	    "published": "Yayınlandı",
//	    "archived": "Arşivlendi",
//	})
func (f *Filterable) SetFilterOptions(options map[string]string) {
	f.filterOptions = options
}

// GetFilterOptions, filtre seçeneklerini döndürür.
// Hiçbir seçenek yapılandırılmamışsa boş bir map döndürür.
//
// Döndürür:
//   - Filtre seçenekleri için değer-etiket eşlemesi
func (f *Filterable) GetFilterOptions() map[string]string {
	if f.filterOptions == nil {
		return make(map[string]string)
	}
	return f.filterOptions
}

// Validatable, alanların doğrulama kurallarına sahip olmasını sağlayan bir mixin'dir.
//
// String tabanlı doğrulama kuralları ve karmaşık doğrulama mantığı uygulayabilen
// özel doğrulayıcı fonksiyonları için yapılandırma sağlar.
//
// # Kullanım
//
//	type EmailField struct {
//	    fields.Base
//	    fields.Validatable
//	}
//
//	field := &EmailField{}
//	field.SetRules([]string{"required", "email"})
//	field.AddValidator(func(value any) error {
//	    email, ok := value.(string)
//	    if !ok {
//	        return fmt.Errorf("email bir string olmalıdır")
//	    }
//	    if !strings.Contains(email, "@") {
//	        return fmt.Errorf("geçersiz email formatı")
//	    }
//	    return nil
//	})
//
// Daha fazla örnek için docs/Fields.md dosyasına bakın.
type Validatable struct {
	rules      []string
	validators []func(any) error
}

// SetRules, alan için doğrulama kurallarını ayarlar.
//
// Kurallar, yerleşik doğrulayıcılara karşılık gelen string tabanlı tanımlayıcılardır
// (örn. "required", "email", "min:8", "max:255").
//
// Parametreler:
//   - rules: Doğrulama kuralı tanımlayıcılarının listesi
//
// Örnek:
//
//	field.SetRules([]string{"required", "email", "max:255"})
func (v *Validatable) SetRules(rules []string) {
	v.rules = rules
}

// GetRules, doğrulama kurallarını döndürür.
// Hiçbir kural yapılandırılmamışsa boş bir dilim döndürür.
//
// Döndürür:
//   - Doğrulama kuralı tanımlayıcılarının dilimi
func (v *Validatable) GetRules() []string {
	if v.rules == nil {
		return []string{}
	}
	return v.rules
}

// AddValidator, özel bir doğrulayıcı fonksiyonu ekler.
//
// Doğrulayıcı, alan değerini alır ve doğrulama başarısız olursa bir hata,
// doğrulama başarılı olursa nil döndürmelidir.
//
// Parametreler:
//   - validator: Alan değerini doğrulayan fonksiyon
//
// Örnek:
//
//	field.AddValidator(func(value any) error {
//	    str, ok := value.(string)
//	    if !ok {
//	        return fmt.Errorf("değer bir string olmalıdır")
//	    }
//	    if len(str) < 8 {
//	        return fmt.Errorf("değer en az 8 karakter olmalıdır")
//	    }
//	    return nil
//	})
func (v *Validatable) AddValidator(validator func(any) error) {
	v.validators = append(v.validators, validator)
}

// GetValidators, özel doğrulayıcı fonksiyonlarını döndürür.
// Hiçbir doğrulayıcı eklenmemişse boş bir dilim döndürür.
//
// Döndürür:
//   - Doğrulayıcı fonksiyonlarının dilimi
func (v *Validatable) GetValidators() []func(any) error {
	if v.validators == nil {
		return []func(any) error{}
	}
	return v.validators
}

// Displayable, alanların görüntüleme formatını özelleştirmesini sağlayan bir mixin'dir.
//
// Alan değerlerinin UI'da nasıl render edileceğini kontrol eden görüntüleme callback'leri
// ve format string'leri için yapılandırma sağlar.
//
// # Kullanım
//
//	type TextField struct {
//	    fields.Base
//	    fields.Displayable
//	}
//
//	field := &TextField{}
//	field.SetDisplayFormat("uppercase")
//	field.SetDisplayCallback(func(ctx *context.Context, value any) string {
//	    return strings.ToUpper(value.(string))
//	})
//
// Daha fazla örnek için docs/Fields.md dosyasına bakın.
type Displayable struct {
	displayCallback func(*context.Context, any) string
	displayFormat   string
}

// SetDisplayCallback, özel bir görüntüleme callback fonksiyonu ayarlar.
//
// Callback, context ve alan değerini alır, ve görüntüleme için değerin
// formatlanmış string temsilini döndürmelidir.
//
// Parametreler:
//   - cb: Alan değerini görüntüleme için formatlayan fonksiyon
//
// Örnek:
//
//	field.SetDisplayCallback(func(ctx *context.Context, value any) string {
//	    timestamp, ok := value.(time.Time)
//	    if !ok {
//	        return ""
//	    }
//	    return timestamp.Format("2006-01-02 15:04:05")
//	})
func (d *Displayable) SetDisplayCallback(cb func(*context.Context, any) string) {
	d.displayCallback = cb
}

// GetDisplayCallback, görüntüleme callback fonksiyonunu döndürür.
// Hiçbir callback ayarlanmamışsa nil döndürür.
//
// Döndürür:
//   - Görüntüleme callback fonksiyonu veya ayarlanmamışsa nil
func (d *Displayable) GetDisplayCallback() func(*context.Context, any) string {
	return d.displayCallback
}

// SetDisplayFormat, görüntüleme format string'ini ayarlar.
//
// Format string'i, frontend'e değerin nasıl görüntüleneceği hakkında bir ipucudur
// (örn. "uppercase", "lowercase", "currency", "date").
//
// Parametreler:
//   - format: Görüntüleme format tanımlayıcısı
//
// Örnek:
//
//	field.SetDisplayFormat("currency")
func (d *Displayable) SetDisplayFormat(format string) {
	d.displayFormat = format
}

// GetDisplayFormat, görüntüleme format string'ini döndürür.
// Hiçbir format ayarlanmamışsa boş string döndürür.
//
// Döndürür:
//   - Görüntüleme format tanımlayıcısı
func (d *Displayable) GetDisplayFormat() string {
	return d.displayFormat
}

// Hideable, alanların görünürlüğünü kontrol etmesini sağlayan bir mixin'dir.
//
// Farklı bağlamlarda (index, detail, create, update) görünürlük ve
// özel gizleme callback'leri için yapılandırma sağlar.
//
// # Kullanım
//
//	type PasswordField struct {
//	    fields.Base
//	    fields.Hideable
//	}
//
//	field := &PasswordField{}
//	field.SetShowOnIndex(false)
//	field.SetShowOnDetail(false)
//	field.SetShowOnCreate(true)
//	field.SetShowOnUpdate(true)
//
// Daha fazla örnek için docs/Fields.md dosyasına bakın.
type Hideable struct {
	hidden       bool
	hideCallback func(*context.Context) bool
	showOnIndex  bool
	showOnDetail bool
	showOnCreate bool
	showOnUpdate bool
}

// SetHidden, alanın gizli olup olmadığını ayarlar.
//
// Parametreler:
//   - hidden: Alanı gizlemek için true, göstermek için false
//
// Örnek:
//
//	field.SetHidden(true)
func (h *Hideable) SetHidden(hidden bool) {
	h.hidden = hidden
}

// IsHidden, alanın gizli olup olmadığını döndürür.
//
// Döndürür:
//   - Alan gizliyse true, değilse false
func (h *Hideable) IsHidden() bool {
	return h.hidden
}

// SetHideCallback, özel bir gizleme callback fonksiyonu ayarlar.
//
// Callback, context'i alır ve alan gizlenmeli ise true,
// gösterilmeli ise false döndürmelidir.
//
// Parametreler:
//   - cb: Alanın gizlenip gizlenmeyeceğini belirleyen fonksiyon
//
// Örnek:
//
//	field.SetHideCallback(func(ctx *context.Context) bool {
//	    user := ctx.User()
//	    return user == nil || !user.IsAdmin
//	})
func (h *Hideable) SetHideCallback(cb func(*context.Context) bool) {
	h.hideCallback = cb
}

// GetHideCallback, gizleme callback fonksiyonunu döndürür.
// Hiçbir callback ayarlanmamışsa nil döndürür.
//
// Döndürür:
//   - Gizleme callback fonksiyonu veya ayarlanmamışsa nil
func (h *Hideable) GetHideCallback() func(*context.Context) bool {
	return h.hideCallback
}

// SetShowOnIndex, alanın liste görünümünde gösterilip gösterilmeyeceğini ayarlar.
func (h *Hideable) SetShowOnIndex(show bool) {
	h.showOnIndex = show
}

// ShowOnIndex, alanın liste görünümünde gösterilip gösterilmeyeceğini döner.
func (h *Hideable) ShowOnIndex() bool {
	return h.showOnIndex
}

// SetShowOnDetail, alanın detay görünümünde gösterilip gösterilmeyeceğini ayarlar.
func (h *Hideable) SetShowOnDetail(show bool) {
	h.showOnDetail = show
}

// ShowOnDetail, alanın detay görünümünde gösterilip gösterilmeyeceğini döner.
func (h *Hideable) ShowOnDetail() bool {
	return h.showOnDetail
}

// SetShowOnCreate, alanın oluşturma formunda gösterilip gösterilmeyeceğini ayarlar.
func (h *Hideable) SetShowOnCreate(show bool) {
	h.showOnCreate = show
}

// ShowOnCreate, alanın oluşturma formunda gösterilip gösterilmeyeceğini döner.
func (h *Hideable) ShowOnCreate() bool {
	return h.showOnCreate
}

// SetShowOnUpdate, alanın güncelleme formunda gösterilip gösterilmeyeceğini ayarlar.
func (h *Hideable) SetShowOnUpdate(show bool) {
	h.showOnUpdate = show
}

// ShowOnUpdate, alanın güncelleme formunda gösterilip gösterilmeyeceğini döner.
func (h *Hideable) ShowOnUpdate() bool {
	return h.showOnUpdate
}
