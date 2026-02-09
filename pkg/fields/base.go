package fields

import (
	"reflect"
	"strings"

	"github.com/ferdiunal/panel.go/pkg/core"
	"github.com/gofiber/fiber/v2"

	"github.com/iancoleman/strcase"
)

// Schema, bir alanın temel yapılandırmasını ve durumunu tutan yapıdır.
//
// Schema, Go Panel API'nin alan sisteminin temel yapı taşıdır. Her alan, bir Schema
// örneği olarak temsil edilir ve frontend ile backend arasında veri taşımak için kullanılır.
//
// # Temel Özellikler
//
// - **JSON Serileştirme**: Tüm alanlar JSON formatında serileştirilebilir
// - **Veri Çıkarma**: Reflection kullanarak model'lerden veri çıkarır
// - **Görünürlük Kontrolü**: Farklı bağlamlarda (list, detail, form) görünürlük kontrolü
// - **Callback Desteği**: Extract, Modify, Resolve, Storage callback'leri
// - **Fluent API**: Zincirleme metod çağrıları ile kolay yapılandırma
//
// # Kategoriler
//
// Schema, 10 ana kategoride özellik içerir:
//
// 1. **Temel Özellikler**: Name, Key, View, Data, Type, Context
// 2. **Durum Özellikleri**: IsReadOnly, IsDisabled, IsImmutable, IsRequired, IsNullable
// 3. **UI Özellikleri**: PlaceholderText, LabelText, HelpTextContent, IsStacked, TextAlign
// 4. **Arama/Filtreleme**: IsFilterable, IsSortable, GlobalSearch
// 5. **Callback'ler**: ExtractCallback, VisibilityCallback, StorageCallback, ModifyCallback
// 6. **Doğrulama**: ValidationRules, CustomValidators
// 7. **Görüntüleme**: DisplayCallback, DisplayedAs, DisplayUsingLabelsFlag
// 8. **Bağımlılıklar**: DependsOnFields, DependencyRules
// 9. **Öneriler**: SuggestionsCallback, AutoCompleteURL, MinCharsForSuggestionsVal
// 10. **Dosya Yükleme**: AcceptedMimeTypes, MaxFileSize, StorageDisk, StoragePath
//
// # Kullanım Örneği
//
//	schema := &fields.Schema{
//	    Key:   "email",
//	    Name:  "E-posta",
//	    View:  "email-field",
//	    Type:  core.TYPE_EMAIL,
//	    Props: make(map[string]interface{}),
//	}
//	schema.OnList().OnDetail().OnForm().Required().Searchable()
//
// # Mixin Desteği
//
// Schema, mixin pattern'ini desteklemez (çünkü kendisi bir veri yapısıdır),
// ancak mixin özelliklerini (Searchable, Sortable, vb.) field'lar aracılığıyla kullanabilir.
//
// Daha fazla bilgi için docs/Fields.md ve .docs/ARCHITECTURE.md dosyalarına bakın.
type Schema struct {
	Name               string                                                              `json:"name"`      // Görünen Ad
	Key                string                                                              `json:"key"`       // Veri Anahtarı
	View               string                                                              `json:"view"`      // Frontend Bileşeni
	Data               interface{}                                                         `json:"data"`      // Alan Değeri
	Type               ElementType                                                         `json:"type"`      // Veri Tipi
	Context            ElementContext                                                      `json:"context"`   // Görünüm Bağlamı (List, Detail, Form)
	IsReadOnly         bool                                                                `json:"read_only"` // Salt okunur mu?
	IsDisabled         bool                                                                `json:"disabled"`  // Devre dışı mı?
	IsImmutable        bool                                                                `json:"immutable"` // Değiştirilemez mi?
	Props              map[string]interface{}                                              `json:"props"`     // Ekstra özellikler
	IsRequired         bool                                                                `json:"required"`  // Zorunlu mu?
	IsNullable         bool                                                                `json:"nullable"`  // Boş bırakılabilir mi?
	PlaceholderText    string                                                              `json:"placeholder"`
	LabelText          string                                                              `json:"label"`
	HelpTextContent    string                                                              `json:"help_text"`
	IsFilterable       bool                                                                `json:"filterable"`
	IsSortable         bool                                                                `json:"sortable"`
	GlobalSearch       bool                                                                `json:"searchable"`
	IsStacked          bool                                                                `json:"stacked"`
	TextAlign          string                                                              `json:"text_align"`
	Suggestions        []interface{}                                                       `json:"suggestions"`
	ExtractCallback    func(value interface{}, item interface{}, c *fiber.Ctx) interface{} `json:"-"`
	VisibilityCallback VisibilityFunc                                                      `json:"-"`
	StorageCallback    StorageCallbackFunc                                                 `json:"-"`
	ModifyCallback     func(value interface{}, c *fiber.Ctx) interface{}                   `json:"-"`
	AutoOptionsConfig  core.AutoOptionsConfig                                              `json:"-"`

	// Validation (Kategori 1)
	ValidationRules  []ValidationRule `json:"validation_rules"`
	CustomValidators []ValidatorFunc  `json:"-"`

	// Display (Kategori 2)
	DisplayCallback        func(interface{}) string `json:"-"`
	DisplayedAs            string                   `json:"displayed_as"`
	DisplayUsingLabelsFlag bool                     `json:"display_using_labels"`
	ResolveHandleValue     string                   `json:"resolve_handle"`

	// Dependencies (Kategori 3)
	DependsOnFields            []string               `json:"depends_on"`
	DependencyRules            map[string]interface{} `json:"dependency_rules"`
	DependencyCallback         DependencyCallbackFunc `json:"-"`
	DependencyCallbackOnCreate DependencyCallbackFunc `json:"-"`
	DependencyCallbackOnUpdate DependencyCallbackFunc `json:"-"`

	// Suggestions (Kategori 4)
	SuggestionsCallback       func(string) []interface{} `json:"-"`
	AutoCompleteURL           string                     `json:"autocomplete_url"`
	MinCharsForSuggestionsVal int                        `json:"min_chars_for_suggestions"`

	// Attachments (Kategori 5)
	AcceptedMimeTypes  []string                             `json:"accepted_mime_types"`
	MaxFileSize        int64                                `json:"max_file_size"`
	StorageDisk        string                               `json:"storage_disk"`
	StoragePath        string                               `json:"storage_path"`
	UploadCallback     func(interface{}, interface{}) error `json:"-"`
	RemoveEXIFDataFlag bool                                 `json:"remove_exif_data"`

	// Repeater (Kategori 6)
	RepeaterFields  []core.Element `json:"-"`
	MinRepeatsCount int            `json:"min_repeats"`
	MaxRepeatsCount int            `json:"max_repeats"`

	// Rich Text (Kategori 7)
	EditorType     string `json:"editor_type"`
	EditorLanguage string `json:"editor_language"`
	EditorTheme    string `json:"editor_theme"`

	// Status (Kategori 8)
	StatusColors map[string]string `json:"status_colors"`
	BadgeVariant string            `json:"badge_variant"`

	// Pivot (Kategori 9)
	IsPivotField      bool   `json:"is_pivot_field"`
	PivotResourceName string `json:"pivot_resource_name"`

	// GORM Veritabanı Yapılandırması (Kategori 10)
	GormConfiguration *GormConfig `json:"-"`
}

// Derleme zamanı kontrolü: Schema'nın core.Element interface'ini implement ettiğinden emin olur
var _ core.Element = (*Schema)(nil)

// GetKey, alanın benzersiz tanımlayıcısını döndürür.
// Key, alanı resource model'indeki bir field'a eşlemek için kullanılır.
//
// Döndürür:
//   - Alanın benzersiz tanımlayıcısı (key)
func (s *Schema) GetKey() string {
	return s.Key
}

// GetView, alanın görünüm tipini döndürür.
// View tipi, alanın UI'da nasıl render edileceğini belirler.
//
// Döndürür:
//   - Görünüm tipi (örn. "text-field", "email-field", "select-field")
func (s *Schema) GetView() string {
	return s.View
}

// GetContext, alanın görüntülendiği bağlamı döndürür.
// Bağlam, alanın nerede gösterilmesi gerektiğini belirtir (form, list, detail).
//
// Döndürür:
//   - Görüntüleme bağlamı (ElementContext)
func (s *Schema) GetContext() ElementContext {
	return s.Context
}

// GetName, bu element'in görünen adını döndürür.
//
// Name, element'in kullanıcı arayüzünde gösterilecek insan okunabilir adıdır.
// Genellikle form label'ı veya tablo başlığı olarak kullanılır.
//
// Döndürür:
//   - Element'in görünen adı (örn: "Name", "Email", "Author")
//
// Örnek:
//
//	field := fields.Text("Name", "name")
//	name := field.GetName() // "Name"
func (s *Schema) GetName() string {
	return s.Name
}

// GetType, bu element'in veri tipini döndürür.
//
// Type, element'in veri tipini belirler ve validation, formatting ve
// OpenAPI schema oluşturma için kullanılır.
// Örneğin: TYPE_TEXT, TYPE_NUMBER, TYPE_BOOLEAN, TYPE_DATE
//
// Döndürür:
//   - Element'in veri tipi (core.ElementType)
//
// Örnek:
//
//	field := fields.Text("Name", "name")
//	fieldType := field.GetType() // core.TYPE_TEXT
func (s *Schema) GetType() core.ElementType {
	return s.Type
}

// Extract, verilen resource'dan veri çıkarır ve alanı doldurur.
//
// Bu metod, reflection kullanarak resource'dan alan değerini çıkarır.
// Resource bir struct veya map olabilir. Struct için, alan adı veya JSON tag'i
// kullanılarak field bulunur. Map için, key kullanılarak değer alınır.
//
// # Özel Durumlar
//
// - **ID Suffix Uyumsuzluğu**: "author_id" gibi bir key için, önce "AuthorId" aranır,
//   bulunamazsa "AuthorID" aranır (Go naming convention)
// - **JSON Tag Desteği**: Struct field'ları JSON tag'leri ile de eşleştirilebilir
// - **Nil Güvenliği**: Resource nil ise hiçbir işlem yapılmaz
//
// Parametreler:
//   - resource: Veri çıkarılacak kaynak (struct veya map)
//
// Örnek:
//
//	type User struct {
//	    Name  string `json:"name"`
//	    Email string `json:"email"`
//	}
//
//	user := &User{Name: "John", Email: "john@example.com"}
//	schema := &Schema{Key: "name"}
//	schema.Extract(user)
//	// schema.Data artık "John" değerini içerir
func (s *Schema) Extract(resource interface{}) {
	if resource == nil {
		return
	}

	// Use reflection to get the value
	v := reflect.ValueOf(resource)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	var value interface{}

	switch v.Kind() {
	case reflect.Struct:
		// Try to find the field by name or json tag
		fieldVal := v.FieldByName(strcase.ToCamel(s.Key))

		// Check for ID suffix mismatch (e.g. key "author_id" -> camel "AuthorId", but struct "AuthorID")
		if !fieldVal.IsValid() {
			camelKey := strcase.ToCamel(s.Key)
			if strings.HasSuffix(camelKey, "Id") {
				fixedName := strings.TrimSuffix(camelKey, "Id") + "ID"
				fieldVal = v.FieldByName(fixedName)
			}
		}

		if !fieldVal.IsValid() {
			// Iterate over fields to check json tags if name doesn't match directly
			for i := 0; i < v.NumField(); i++ {
				typeField := v.Type().Field(i)
				tag := typeField.Tag.Get("json")
				if tag == s.Key || strings.Split(tag, ",")[0] == s.Key {
					fieldVal = v.Field(i)
					break
				}
			}
		}

		if fieldVal.IsValid() && fieldVal.CanInterface() {
			value = fieldVal.Interface()
		}
	case reflect.Map:
		// Check if it's a map[string]interface{} or similar
		val := v.MapIndex(reflect.ValueOf(s.Key))
		if val.IsValid() && val.CanInterface() {
			value = val.Interface()
		}
	}

	s.Data = value
}

// JsonSerialize, alanı JSON uyumlu bir map'e serileştirir.
//
// Bu metod, alanın tüm özelliklerini JSON encoding için hazır bir map'e dönüştürür.
// Frontend'e gönderilmek üzere alanın durumunu temsil eder.
//
// # Serileştirilen Özellikler
//
// - **view**: Görünüm tipi (frontend bileşeni)
// - **type**: Veri tipi (ElementType)
// - **key**: Benzersiz tanımlayıcı
// - **name**: Görünen ad
// - **data**: Alan değeri
// - **props**: Ekstra özellikler
// - **context**: Görüntüleme bağlamı
// - **placeholder**: Yer tutucu metni
// - **label**: Etiket metni
// - **help_text**: Yardım metni
// - **read_only**: Salt okunur durumu
// - **disabled**: Devre dışı durumu
// - **required**: Zorunlu durumu
// - **nullable**: Boş bırakılabilir durumu
// - **sortable**: Sıralanabilir durumu
// - **filterable**: Filtrelenebilir durumu
// - **stacked**: Yığılmış (tam genişlik) durumu
// - **text_align**: Metin hizalama
//
// Döndürür:
//   - JSON encoding için hazır map[string]any
//
// Örnek:
//
//	schema := &Schema{Key: "email", Name: "E-posta", View: "email-field"}
//	json := schema.JsonSerialize()
//	// json["key"] == "email"
//	// json["name"] == "E-posta"
//	// json["view"] == "email-field"
func (s *Schema) JsonSerialize() map[string]interface{} {
	return map[string]interface{}{
		"view":        s.View,
		"type":        s.Type,
		"key":         s.Key,
		"name":        s.Name,
		"data":        s.Data,
		"props":       s.Props,
		"context":     s.Context,
		"placeholder": s.PlaceholderText,
		"label":       s.LabelText,
		"help_text":   s.HelpTextContent,
		"read_only":   s.IsReadOnly,
		"disabled":    s.IsDisabled,
		"required":    s.IsRequired,
		"nullable":    s.IsNullable,
		"sortable":    s.IsSortable,
		"filterable":  s.IsFilterable,
		"stacked":     s.IsStacked,
		"text_align":  s.TextAlign,
	}
}

// Fluent Setters - Zincirleme Metod Çağrıları
//
// Bu bölüm, Schema yapısının özelliklerini ayarlamak için fluent API pattern'ini kullanır.
// Her metod, Schema pointer'ını döndürerek zincirleme çağrılara olanak tanır.

// SetName, alanın görünen adını ayarlar.
//
// Bu metod, UI'da kullanıcıya gösterilecek alan adını belirler.
// Genellikle form etiketleri, tablo başlıkları ve detay görünümlerinde kullanılır.
//
// # Parametreler
//
//   - name: Alanın görünen adı (örn. "E-posta Adresi", "Kullanıcı Adı")
//
// # Döndürür
//
//   - Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Örnek
//
//	field := &Schema{Key: "email"}
//	field.SetName("E-posta Adresi")
func (s *Schema) SetName(name string) Element {
	s.Name = name
	return s
}

// SetKey, alanın benzersiz tanımlayıcısını ayarlar.
//
// Key, alanı resource model'indeki bir field'a eşlemek için kullanılır.
// Genellikle veritabanı sütun adı veya struct field adı ile eşleşir.
//
// # Parametreler
//
//   - key: Alanın benzersiz tanımlayıcısı (örn. "email", "user_name", "created_at")
//
// # Döndürür
//
//   - Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Önemli Notlar
//
//   - Key, snake_case formatında olmalıdır
//   - Key, resource model'deki field adı ile eşleşmelidir
//   - Key değiştirildiğinde, veri çıkarma işlemi etkilenir
//
// # Örnek
//
//	field := &Schema{}
//	field.SetKey("email_address")
func (s *Schema) SetKey(key string) Element {
	s.Key = key
	return s
}

// SetContext, alanın görüntüleneceği bağlamı ayarlar.
//
// Context, alanın hangi sayfalarda (list, detail, form) görüneceğini belirler.
// Bu metod, görünürlük kontrolü için temel mekanizmayı sağlar.
//
// # Parametreler
//
//   - context: Görüntüleme bağlamı (ElementContext enum değeri)
//
// # Döndürür
//
//   - Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Kullanılabilir Context Değerleri
//
//   - SHOW_ON_LIST: Liste görünümünde göster
//   - SHOW_ON_DETAIL: Detay görünümünde göster
//   - SHOW_ON_FORM: Form görünümünde göster
//   - HIDE_ON_LIST: Liste görünümünde gizle
//   - HIDE_ON_DETAIL: Detay görünümünde gizle
//   - HIDE_ON_CREATE: Oluşturma formunda gizle
//   - HIDE_ON_UPDATE: Güncelleme formunda gizle
//   - ONLY_ON_LIST: Sadece liste görünümünde göster
//   - ONLY_ON_DETAIL: Sadece detay görünümünde göster
//   - ONLY_ON_CREATE: Sadece oluşturma formunda göster
//   - ONLY_ON_UPDATE: Sadece güncelleme formunda göster
//   - ONLY_ON_FORM: Sadece formlarda göster
//
// # Örnek
//
//	field := &Schema{Key: "id"}
//	field.SetContext(HIDE_ON_FORM)
func (s *Schema) SetContext(context ElementContext) Element {
	// Eğer mevcut context boş değilse ve yeni context farklıysa, birleştir
	if s.Context != "" && s.Context != context {
		// Context'i space-separated string olarak birleştir
		s.Context = ElementContext(string(s.Context) + " " + string(context))
	} else {
		s.Context = context
	}
	return s
}

// OnList, alanı liste görünümünde gösterir.
//
// Bu metod, alanın tablo/liste görünümünde görünür olmasını sağlar.
// Genellikle özet bilgiler için kullanılır.
//
// # Döndürür
//
//   - Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Örnek
//
//	field := Text("name", "İsim").OnList()
func (s *Schema) OnList() Element {
	return s.SetContext(SHOW_ON_LIST)
}

// OnDetail, alanı detay görünümünde gösterir.
//
// Bu metod, alanın kayıt detay sayfasında görünür olmasını sağlar.
// Genellikle tüm alan bilgileri için kullanılır.
//
// # Döndürür
//
//   - Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Örnek
//
//	field := Text("description", "Açıklama").OnDetail()
func (s *Schema) OnDetail() Element {
	return s.SetContext(SHOW_ON_DETAIL)
}

// OnForm, alanı form görünümünde gösterir.
//
// Bu metod, alanın hem oluşturma hem de güncelleme formlarında görünür olmasını sağlar.
// Kullanıcının düzenleyebileceği alanlar için kullanılır.
//
// # Döndürür
//
//   - Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Örnek
//
//	field := Text("email", "E-posta").OnForm()
func (s *Schema) OnForm() Element {
	return s.SetContext(SHOW_ON_FORM)
}

// HideOnList, alanı liste görünümünde gizler.
//
// Bu metod, alanın tablo/liste görünümünde gizlenmesini sağlar.
// Genellikle uzun metinler veya hassas bilgiler için kullanılır.
//
// # Döndürür
//
//   - Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Örnek
//
//	field := Textarea("bio", "Biyografi").HideOnList()
func (s *Schema) HideOnList() Element {
	return s.SetContext(HIDE_ON_LIST)
}

// HideOnDetail, alanı detay görünümünde gizler.
//
// Bu metod, alanın kayıt detay sayfasında gizlenmesini sağlar.
// Genellikle sadece form için gerekli olan alanlar için kullanılır.
//
// # Döndürür
//
//   - Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Örnek
//
//	field := Password("password", "Şifre").HideOnDetail()
func (s *Schema) HideOnDetail() Element {
	return s.SetContext(HIDE_ON_DETAIL)
}

// HideOnCreate, alanı oluşturma formunda gizler.
//
// Bu metod, alanın yeni kayıt oluşturma formunda gizlenmesini sağlar.
// Genellikle otomatik oluşturulan alanlar (ID, timestamps) için kullanılır.
//
// # Döndürür
//
//   - Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Örnek
//
//	field := Text("id", "ID").HideOnCreate()
func (s *Schema) HideOnCreate() Element {
	return s.SetContext(HIDE_ON_CREATE)
}

// HideOnUpdate, alanı güncelleme formunda gizler.
//
// Bu metod, alanın kayıt güncelleme formunda gizlenmesini sağlar.
// Genellikle değiştirilmemesi gereken alanlar için kullanılır.
//
// # Döndürür
//
//   - Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Örnek
//
//	field := Text("username", "Kullanıcı Adı").HideOnUpdate()
func (s *Schema) HideOnUpdate() Element {
	return s.SetContext(HIDE_ON_UPDATE)
}

// OnlyOnList, alanı sadece liste görünümünde gösterir.
//
// Bu metod, alanın yalnızca tablo/liste görünümünde görünür olmasını sağlar.
// Diğer tüm görünümlerde (detail, form) gizlenir.
//
// # Döndürür
//
//   - Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Kullanım Senaryoları
//
//   - Özet bilgiler (örn. kayıt sayısı)
//   - Hesaplanmış değerler (örn. toplam tutar)
//   - Hızlı erişim linkleri
//
// # Örnek
//
//	field := Text("summary", "Özet").OnlyOnList()
func (s *Schema) OnlyOnList() Element {
	return s.SetContext(ONLY_ON_LIST)
}

// OnlyOnDetail, alanı sadece detay görünümünde gösterir.
//
// Bu metod, alanın yalnızca kayıt detay sayfasında görünür olmasını sağlar.
// Liste ve formlarda gizlenir.
//
// # Döndürür
//
//   - Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Kullanım Senaryoları
//
//   - Detaylı açıklamalar
//   - İlişkili kayıt bilgileri
//   - Metadata bilgileri
//
// # Örnek
//
//	field := Textarea("full_description", "Detaylı Açıklama").OnlyOnDetail()
func (s *Schema) OnlyOnDetail() Element {
	return s.SetContext(ONLY_ON_DETAIL)
}

// OnlyOnCreate, alanı sadece oluşturma formunda gösterir.
//
// Bu metod, alanın yalnızca yeni kayıt oluşturma formunda görünür olmasını sağlar.
// Güncelleme formunda ve diğer görünümlerde gizlenir.
//
// # Döndürür
//
//   - Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Kullanım Senaryoları
//
//   - Şifre alanları (ilk oluşturmada)
//   - Başlangıç yapılandırma alanları
//   - Bir kez ayarlanan değerler
//
// # Örnek
//
//	field := Password("password", "Şifre").OnlyOnCreate()
func (s *Schema) OnlyOnCreate() Element {
	return s.SetContext(ONLY_ON_CREATE)
}

// OnlyOnUpdate, alanı sadece güncelleme formunda gösterir.
//
// Bu metod, alanın yalnızca kayıt güncelleme formunda görünür olmasını sağlar.
// Oluşturma formunda ve diğer görünümlerde gizlenir.
//
// # Döndürür
//
//   - Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Kullanım Senaryoları
//
//   - Şifre değiştirme alanları
//   - Durum güncelleme alanları
//   - Sadece mevcut kayıtlar için geçerli alanlar
//
// # Örnek
//
//	field := Password("new_password", "Yeni Şifre").OnlyOnUpdate()
func (s *Schema) OnlyOnUpdate() Element {
	return s.SetContext(ONLY_ON_UPDATE)
}

// OnlyOnForm, alanı sadece formlarda gösterir.
//
// Bu metod, alanın hem oluşturma hem de güncelleme formlarında görünür olmasını,
// ancak liste ve detay görünümlerinde gizlenmesini sağlar.
//
// # Döndürür
//
//   - Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Kullanım Senaryoları
//
//   - Sadece düzenleme için gerekli alanlar
//   - Kullanıcı girdisi gerektiren alanlar
//   - Görüntülenmesi gerekmeyen form alanları
//
// # Örnek
//
//	field := Hidden("csrf_token", "CSRF Token").OnlyOnForm()
func (s *Schema) OnlyOnForm() Element {
	return s.SetContext(ONLY_ON_FORM)
}

// ReadOnly, alanı salt okunur olarak işaretler.
//
// Salt okunur alanlar, kullanıcı tarafından düzenlenemez ancak görüntülenebilir.
// Genellikle sistem tarafından oluşturulan veya hesaplanan değerler için kullanılır.
//
// # Döndürür
//
//   - Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Kullanım Senaryoları
//
//   - Otomatik oluşturulan ID'ler
//   - Timestamp alanları (created_at, updated_at)
//   - Hesaplanan değerler
//   - Sistem tarafından yönetilen alanlar
//
// # Örnek
//
//	field := Text("created_at", "Oluşturma Tarihi").ReadOnly()
func (s *Schema) ReadOnly() Element {
	s.IsReadOnly = true
	return s
}

// WithProps, alana özel bir özellik ekler.
//
// Props, alanın frontend bileşenine iletilecek ekstra özellikleri saklar.
// Bu metod, özel bileşen davranışları veya stil ayarları için kullanılır.
//
// # Parametreler
//
//   - key: Özellik anahtarı
//   - value: Özellik değeri (herhangi bir tip olabilir)
//
// # Döndürür
//
//   - Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Kullanım Senaryoları
//
//   - Özel bileşen yapılandırması
//   - Frontend'e özel veri iletimi
//   - UI davranış kontrolü
//   - Stil ve tema ayarları
//
// # Örnek
//
//	field := Text("name", "İsim").
//	    WithProps("maxLength", 100).
//	    WithProps("autoFocus", true).
//	    WithProps("className", "custom-input")
func (s *Schema) WithProps(key string, value interface{}) Element {
	if s.Props == nil {
		s.Props = make(map[string]interface{})
	}
	s.Props[key] = value
	return s
}

// Tooltip, alana bir tooltip (bilgi balonu) ekler.
//
// Tooltip, kullanıcıya alan hakkında ek bilgi sağlamak için kullanılır.
// Frontend'te label'ın yanında bir info ikonu olarak gösterilir ve üzerine
// gelindiğinde tooltip metni görüntülenir.
//
// # Döndürür
//
//   - Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Kullanım Örneği
//
//	field := Text("Kullanıcı Adı", "username").
//	    Tooltip("Kullanıcı adınız benzersiz olmalıdır ve en az 3 karakter içermelidir")
//
//	field := Email("E-posta", "email").
//	    Tooltip("Geçerli bir e-posta adresi giriniz").
//	    Required()
//
// # Not
//
// Bu metod, arka planda WithProps("tooltip", text) çağrısı yapar.
// Tooltip özelliği frontend bileşenleri tarafından desteklenmelidir.
func (s *Schema) Tooltip(text string) Element {
	return s.WithProps("tooltip", text)
}

// Disabled, alanı devre dışı bırakır.
//
// Devre dışı alanlar, kullanıcı tarafından ne düzenlenebilir ne de etkileşime girebilir.
// ReadOnly'den farklı olarak, görsel olarak da devre dışı görünür.
//
// # Döndürür
//
//   - Element: Zincirleme çağrılar için Schema pointer'ı
//
// # ReadOnly vs Disabled
//
//   - **ReadOnly**: Görüntülenebilir, kopyalanabilir, ancak düzenlenemez
//   - **Disabled**: Tamamen etkileşimsiz, genellikle gri/soluk görünür
//
// # Kullanım Senaryoları
//
//   - Koşullu olarak devre dışı bırakılan alanlar
//   - İzin kontrolü gerektiren alanlar
//   - Bağımlı alanlar (başka bir alan seçilene kadar)
//
// # Örnek
//
//	field := Select("city", "Şehir").Disabled()
func (s *Schema) Disabled() Element {
	s.IsDisabled = true
	return s
}

// Immutable, alanı değiştirilemez olarak işaretler.
//
// Immutable alanlar, oluşturulduktan sonra asla değiştirilemez.
// Genellikle oluşturma formunda gösterilir, güncelleme formunda gizlenir veya salt okunur olur.
//
// # Döndürür
//
//   - Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Kullanım Senaryoları
//
//   - Kullanıcı adı (username)
//   - E-posta adresi (bazı sistemlerde)
//   - Benzersiz tanımlayıcılar
//   - Kayıt oluşturma zamanı
//
// # Örnek
//
//	field := Text("username", "Kullanıcı Adı").Immutable()
func (s *Schema) Immutable() Element {
	s.IsImmutable = true
	return s
}

// Required, alanı zorunlu olarak işaretler.
//
// Zorunlu alanlar, form gönderilmeden önce doldurulmalıdır.
// Frontend ve backend validasyonunda kullanılır.
//
// # Döndürür
//
//   - Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Önemli Notlar
//
//   - Frontend'de otomatik validasyon eklenir
//   - Genellikle kırmızı yıldız (*) ile gösterilir
//   - Backend validasyonu ile birlikte kullanılmalıdır
//
// # Örnek
//
//	field := Text("email", "E-posta").Required().Email()
func (s *Schema) Required() Element {
	s.IsRequired = true
	return s
}

// Nullable, alanın boş (null) değer alabileceğini belirtir.
//
// Nullable alanlar, veritabanında NULL değer saklayabilir.
// Varsayılan olarak alanlar nullable değildir.
//
// # Döndürür
//
//   - Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Kullanım Senaryoları
//
//   - Opsiyonel ilişkiler (foreign key'ler)
//   - Opsiyonel tarih alanları
//   - Opsiyonel sayısal değerler
//
// # Örnek
//
//	field := BelongsTo("author", "Yazar").Nullable()
func (s *Schema) Nullable() Element {
	s.IsNullable = true
	return s
}

// Placeholder, alan için yer tutucu metin ayarlar.
//
// Placeholder, alan boşken gösterilen açıklayıcı metindir.
// Kullanıcıya ne girilmesi gerektiği hakkında ipucu verir.
//
// # Parametreler
//
//   - placeholder: Yer tutucu metin
//
// # Döndürür
//
//   - Element: Zincirleme çağrılar için Schema pointer'ı
//
// # İyi Uygulamalar
//
//   - Kısa ve açıklayıcı olmalı
//   - Örnek değer gösterilebilir
//   - Zorunlu alan açıklaması yerine kullanılmamalı
//
// # Örnek
//
//	field := Text("email", "E-posta").
//	    Placeholder("ornek@email.com")
func (s *Schema) Placeholder(placeholder string) Element {
	s.PlaceholderText = placeholder
	return s
}

// Label, alan için etiket metni ayarlar.
//
// Label, alanın yanında veya üstünde gösterilen açıklayıcı metindir.
// Name'den farklı olarak, daha uzun ve açıklayıcı olabilir.
//
// # Parametreler
//
//   - label: Etiket metni
//
// # Döndürür
//
//   - Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Name vs Label
//
//   - **Name**: Kısa, genel başlık (örn. "E-posta")
//   - **Label**: Daha açıklayıcı (örn. "İş E-posta Adresiniz")
//
// # Örnek
//
//	field := Text("email", "E-posta").
//	    Label("Lütfen geçerli bir e-posta adresi girin")
func (s *Schema) Label(label string) Element {
	s.LabelText = label
	return s
}

// HelpText, alan için yardım metni ayarlar.
//
// HelpText, alanın altında gösterilen açıklayıcı veya yönlendirici metindir.
// Kullanıcıya alan hakkında ek bilgi verir.
//
// # Parametreler
//
//   - helpText: Yardım metni
//
// # Döndürür
//
//   - Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Kullanım Senaryoları
//
//   - Format açıklaması (örn. "DD/MM/YYYY formatında")
//   - Karakter limiti bilgisi
//   - Güvenlik uyarıları
//   - Kullanım önerileri
//
// # Örnek
//
//	field := Password("password", "Şifre").
//	    HelpText("En az 8 karakter, bir büyük harf ve bir rakam içermelidir")
func (s *Schema) HelpText(helpText string) Element {
	s.HelpTextContent = helpText
	return s
}

// Filterable, alanın filtrelenebilir olduğunu belirtir.
//
// Filterable alanlar, liste görünümünde filtre seçenekleri sunar.
// Kullanıcılar bu alana göre kayıtları filtreleyebilir.
//
// # Döndürür
//
//   - Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Kullanım Senaryoları
//
//   - Durum alanları (aktif/pasif)
//   - Kategori alanları
//   - Tarih aralıkları
//   - Sayısal aralıklar
//
// # Örnek
//
//	field := Select("status", "Durum").Filterable()
func (s *Schema) Filterable() Element {
	s.IsFilterable = true
	return s
}

// Sortable, alanın sıralanabilir olduğunu belirtir.
//
// Sortable alanlar, liste görünümünde sıralama seçenekleri sunar.
// Kullanıcılar bu alana göre kayıtları artan veya azalan şekilde sıralayabilir.
//
// # Döndürür
//
//   - Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Kullanım Senaryoları
//
//   - İsim alanları
//   - Tarih alanları
//   - Sayısal alanlar
//   - Durum alanları
//
// # Örnek
//
//	field := Text("name", "İsim").Sortable()
func (s *Schema) Sortable() Element {
	s.IsSortable = true
	return s
}

// Searchable, alanın global aramada kullanılacağını belirtir.
//
// Searchable alanlar, genel arama kutusunda aranabilir.
// Kullanıcı arama yaptığında, bu alanların içeriği taranır.
//
// # Döndürür
//
//   - Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Kullanım Senaryoları
//
//   - İsim alanları
//   - E-posta alanları
//   - Açıklama alanları
//   - Başlık alanları
//
// # Önemli Notlar
//
//   - Performans için dikkatli kullanılmalı
//   - Çok fazla alan searchable yapılmamalı
//   - Genellikle text alanları için uygundur
//
// # Örnek
//
//	field := Text("name", "İsim").Searchable()
func (s *Schema) Searchable() Element {
	s.GlobalSearch = true
	return s
}

// IsSearchable, alanın aranabilir olup olmadığını kontrol eder.
//
// # Döndürür
//
//   - bool: Alanın aranabilir olup olmadığı
//
// # Örnek
//
//	if field.IsSearchable() {
//	    // Arama sorgusu oluştur
//	}
func (s *Schema) IsSearchable() bool {
	return s.GlobalSearch
}

// Stacked, alanın tam genişlikte görüntüleneceğini belirtir.
//
// Stacked alanlar, formda kendi satırını kaplar (100% genişlik).
// Varsayılan olarak alanlar yan yana dizilir, stacked alanlar alt alta dizilir.
//
// # Döndürür
//
//   - Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Kullanım Senaryoları
//
//   - Textarea alanları
//   - Uzun text alanları
//   - Önemli alanlar
//   - Görsel vurgu gerektiren alanlar
//
// # Örnek
//
//	field := Textarea("description", "Açıklama").Stacked()
func (s *Schema) Stacked() Element {
	s.IsStacked = true
	return s
}

// SetTextAlign, alan içeriğinin hizalamasını ayarlar.
//
// Text align, alanın içeriğinin nasıl hizalanacağını belirler.
// Genellikle liste görünümünde kullanılır.
//
// # Parametreler
//
//   - align: Hizalama değeri ("left", "center", "right", "justify")
//
// # Döndürür
//
//   - Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Kullanım Senaryoları
//
//   - Sayısal alanlar için sağa hizalama
//   - Başlıklar için ortaya hizalama
//   - Metin alanları için sola hizalama (varsayılan)
//
// # Örnek
//
//	field := Number("price", "Fiyat").SetTextAlign("right")
func (s *Schema) SetTextAlign(align string) Element {
	s.TextAlign = align
	return s
}

// IsVisible, alanın belirli bir bağlamda görünür olup olmadığını kontrol eder.
//
// Bu metod, VisibilityCallback varsa onu çağırır, yoksa true döner.
// Dinamik görünürlük kontrolü için kullanılır.
//
// # Parametreler
//
//   - ctx: Resource bağlamı (kullanıcı, izinler, vb.)
//
// # Döndürür
//
//   - bool: Alanın görünür olup olmadığı
//
// # Örnek
//
//	if field.IsVisible(ctx) {
//	    // Alanı göster
//	}
func (s *Schema) IsVisible(ctx *core.ResourceContext) bool {
	// Önce VisibilityCallback kontrolü yap
	if s.VisibilityCallback != nil {
		return s.VisibilityCallback(ctx)
	}

	// Context bazlı görünürlük kontrolü yap
	if ctx != nil && ctx.VisibilityCtx != "" {
		return s.IsVisibleInContext(ctx.VisibilityCtx)
	}

	return true
}

// CanSee, alan için görünürlük callback'i ayarlar.
//
// Bu metod, alanın dinamik olarak gösterilip gizlenmesini sağlar.
// Callback, kullanıcı izinlerine, rol'lere veya diğer koşullara göre çalışabilir.
//
// # Parametreler
//
//   - fn: Görünürlük kontrolü yapan fonksiyon
//
// # Döndürür
//
//   - Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Kullanım Senaryoları
//
//   - İzin tabanlı görünürlük
//   - Rol tabanlı görünürlük
//   - Koşullu alan gösterimi
//   - Dinamik form yapıları
//
// # Örnek
//
//	field := Text("salary", "Maaş").CanSee(func(ctx *core.ResourceContext) bool {
//	    return ctx.User.IsAdmin()
//	})
func (s *Schema) CanSee(fn VisibilityFunc) Element {
	s.VisibilityCallback = fn
	return s
}

// StoreAs, alan için özel depolama callback'i ayarlar.
//
// Bu metod, alanın veritabanına nasıl kaydedileceğini özelleştirmeye olanak tanır.
// Veri dönüşümü, şifreleme, hash'leme gibi işlemler için kullanılır.
//
// # Parametreler
//
//   - fn: Depolama işlemini yapan fonksiyon
//
// # Döndürür
//
//   - Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Kullanım Senaryoları
//
//   - Şifre hash'leme
//   - Veri şifreleme
//   - Format dönüşümü
//   - Hesaplanan değerler
//
// # Örnek
//
//	field := Password("password", "Şifre").StoreAs(func(value interface{}, ctx *fiber.Ctx) interface{} {
//	    return bcrypt.HashPassword(value.(string))
//	})
func (s *Schema) StoreAs(fn StorageCallbackFunc) Element {
	s.StorageCallback = fn
	return s
}

// GetStorageCallback, alanın depolama callback'ini döner.
//
// # Döndürür
//
//   - StorageCallbackFunc: Depolama callback fonksiyonu (nil olabilir)
func (s *Schema) GetStorageCallback() StorageCallbackFunc {
	return s.StorageCallback
}

// Resolve, alan için özel veri çıkarma callback'i ayarlar.
//
// Bu metod, alanın resource'dan nasıl çıkarılacağını özelleştirmeye olanak tanır.
// İlişkili veriler, hesaplanan değerler veya özel dönüşümler için kullanılır.
//
// # Parametreler
//
//   - fn: Veri çıkarma işlemini yapan fonksiyon
//
// # Döndürür
//
//   - Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Kullanım Senaryoları
//
//   - İlişkili veri çıkarma
//   - Hesaplanan değerler
//   - Format dönüşümü
//   - Özel veri işleme
//
// # Örnek
//
//	field := Text("full_name", "Tam İsim").Resolve(func(value interface{}, item interface{}, c *fiber.Ctx) interface{} {
//	    user := item.(*User)
//	    return user.FirstName + " " + user.LastName
//	})
func (s *Schema) Resolve(fn func(value interface{}, item interface{}, c *fiber.Ctx) interface{}) Element {
	s.ExtractCallback = fn
	return s
}

// GetResolveCallback, alanın resolve callback'ini döner.
//
// # Döndürür
//
//   - func: Resolve callback fonksiyonu (nil olabilir)
func (s *Schema) GetResolveCallback() func(value interface{}, item interface{}, c *fiber.Ctx) interface{} {
	return s.ExtractCallback
}

// Modify, alan için değer modifikasyon callback'i ayarlar.
//
// Bu metod, alanın değerinin frontend'e gönderilmeden önce değiştirilmesini sağlar.
// Format dönüşümü, maskeleme, hesaplama gibi işlemler için kullanılır.
//
// # Parametreler
//
//   - fn: Değer modifikasyonu yapan fonksiyon
//
// # Döndürür
//
//   - Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Kullanım Senaryoları
//
//   - Tarih formatı dönüşümü
//   - Para birimi formatı
//   - Hassas veri maskeleme
//   - Değer hesaplama
//
// # Örnek
//
//	field := Text("phone", "Telefon").Modify(func(value interface{}, c *fiber.Ctx) interface{} {
//	    phone := value.(string)
//	    return maskPhone(phone) // (555) ***-**-12
//	})
func (s *Schema) Modify(fn func(value interface{}, c *fiber.Ctx) interface{}) Element {
	s.ModifyCallback = fn
	return s
}

// GetModifyCallback, alanın modify callback'ini döner.
//
// # Döndürür
//
//   - func: Modify callback fonksiyonu (nil olabilir)
func (s *Schema) GetModifyCallback() func(value interface{}, c *fiber.Ctx) interface{} {
	return s.ModifyCallback
}

// Options, alan için seçenek listesi ayarlar.
//
// Bu metod, select, radio, checkbox gibi seçim alanları için kullanılır.
// Statik seçenek listesi sağlar.
//
// # Parametreler
//
//   - options: Seçenek listesi (slice veya map olabilir)
//
// # Döndürür
//
//   - Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Kullanım Senaryoları
//
//   - Sabit seçenek listeleri
//   - Enum değerleri
//   - Durum seçenekleri
//
// # Örnek
//
//	field := Select("status", "Durum").Options([]map[string]interface{}{
//	    {"value": "active", "label": "Aktif"},
//	    {"value": "inactive", "label": "Pasif"},
//	})
func (s *Schema) Options(options interface{}) Element {
	s.Props["options"] = options
	return s
}

// AutoOptions, alan için otomatik seçenek yükleme ayarlar.
//
// Bu metod, ilişkili model'den otomatik olarak seçeneklerin yüklenmesini sağlar.
// BelongsTo, HasOne gibi ilişki alanları için kullanılır.
//
// # Parametreler
//
//   - displayField: Seçeneklerde gösterilecek alan adı
//
// # Döndürür
//
//   - Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Kullanım Senaryoları
//
//   - İlişkili model seçenekleri
//   - Dinamik seçenek listeleri
//   - Veritabanından seçenek yükleme
//
// # Örnek
//
//	field := BelongsTo("category_id", "Kategori").AutoOptions("name")
func (s *Schema) AutoOptions(displayField string) Element {
	s.AutoOptionsConfig.Enabled = true
	s.AutoOptionsConfig.DisplayField = displayField
	return s
}

// GetAutoOptionsConfig, alanın auto options yapılandırmasını döner.
//
// # Döndürür
//
//   - core.AutoOptionsConfig: Auto options yapılandırması
func (s *Schema) GetAutoOptionsConfig() core.AutoOptionsConfig {
	return s.AutoOptionsConfig
}

// Default, alan için varsayılan değer ayarlar.
//
// Bu metod, alanın başlangıç değerini belirler.
// Form ilk açıldığında bu değer gösterilir.
//
// # Parametreler
//
//   - value: Varsayılan değer
//
// # Döndürür
//
//   - Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Kullanım Senaryoları
//
//   - Form başlangıç değerleri
//   - Önerilen değerler
//   - Varsayılan seçimler
//
// # Örnek
//
//	field := Select("status", "Durum").Default("active")
func (s *Schema) Default(value interface{}) Element {
	s.Data = value
	return s
}

// Core Interface Methods - Element Interface Implementasyonu
//
// Bu bölüm, core.Element interface'ini implement eden metodları içerir.
// Bu metodlar, alanın görünürlük kontrolü, veri çıkarma ve metadata yönetimi için kullanılır.

// IsHidden, alanın belirli bir görünürlük bağlamında gizli olup olmadığını kontrol eder.
//
// Bu metod, IsVisibleInContext metodunun tersini döner.
// Görünürlük kontrolü için kullanılır.
//
// # Parametreler
//
//   - ctx: Görünürlük bağlamı (ContextIndex, ContextDetail, ContextCreate, ContextUpdate, ContextPreview)
//
// # Döndürür
//
//   - bool: Alanın gizli olup olmadığı
//
// # Kullanım Senaryoları
//
//   - Koşullu alan gösterimi
//   - Bağlam bazlı görünürlük kontrolü
//   - Form ve liste filtreleme
//
// # Örnek
//
//	if field.IsHidden(ContextCreate) {
//	    // Alanı oluşturma formunda gösterme
//	}
func (s *Schema) IsHidden(ctx VisibilityContext) bool {
	return !s.IsVisibleInContext(ctx)
}

// ResolveForDisplay, alanın değerini görüntüleme için çözümler.
//
// Bu metod, verilen kayıttan alanın değerini çıkarır ve görüntüleme için hazırlar.
// Reflection kullanarak struct field'larını veya map değerlerini okur.
//
// # Parametreler
//
//   - item: Değer çıkarılacak kayıt (struct veya map)
//
// # Döndürür
//
//   - interface{}: Çıkarılan değer
//   - error: Hata durumunda hata mesajı (şu an her zaman nil)
//
// # Çalışma Mantığı
//
// 1. Item nil ise, mevcut Data değerini döner
// 2. Item struct ise, field adına göre değer çıkarır
// 3. JSON tag'leri de kontrol edilir
// 4. Item map ise, key'e göre değer çıkarır
//
// # Örnek
//
//	type User struct {
//	    Name string `json:"name"`
//	}
//	user := &User{Name: "John"}
//	value, err := field.ResolveForDisplay(user)
//	// value == "John"
func (s *Schema) ResolveForDisplay(item interface{}) (interface{}, error) {
	if item == nil {
		return s.Data, nil
	}

	// Extract the value from the item
	v := reflect.ValueOf(item)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	var value interface{}

	switch v.Kind() {
	case reflect.Struct:
		fieldVal := v.FieldByName(strcase.ToCamel(s.Key))
		if !fieldVal.IsValid() {
			for i := 0; i < v.NumField(); i++ {
				typeField := v.Type().Field(i)
				tag := typeField.Tag.Get("json")
				if tag == s.Key || strings.Split(tag, ",")[0] == s.Key {
					fieldVal = v.Field(i)
					break
				}
			}
		}

		if fieldVal.IsValid() && fieldVal.CanInterface() {
			value = fieldVal.Interface()
		}
	case reflect.Map:
		val := v.MapIndex(reflect.ValueOf(s.Key))
		if val.IsValid() && val.CanInterface() {
			value = val.Interface()
		}
	}

	return value, nil
}

// GetDependencies, alanın bağımlı olduğu diğer alanların listesini döner.
//
// Bu metod, Props içindeki "dependencies" anahtarını kontrol eder.
// Bağımlılık yönetimi için kullanılır.
//
// # Döndürür
//
//   - []string: Bağımlı alan adları listesi (boş slice olabilir)
//
// # Kullanım Senaryoları
//
//   - Dinamik form alanları
//   - Koşullu alan gösterimi
//   - Cascade seçim alanları
//   - Bağımlılık grafiği oluşturma
//
// # Örnek
//
//	deps := field.GetDependencies()
//	for _, dep := range deps {
//	    // Bağımlı alanı kontrol et
//	}
func (s *Schema) GetDependencies() []string {
	deps, ok := s.Props["dependencies"].([]string)
	if ok {
		return deps
	}
	return []string{}
}

// IsConditionallyVisible, alanın kayıt değerlerine göre görünür olup olmadığını belirler.
//
// Bu metod, VisibilityCallback varsa onu çağırır, yoksa true döner.
// Dinamik görünürlük kontrolü için kullanılır.
//
// # Parametreler
//
//   - item: Kontrol edilecek kayıt
//
// # Döndürür
//
//   - bool: Alanın görünür olup olmadığı
//
// # Kullanım Senaryoları
//
//   - Kayıt durumuna göre alan gösterimi
//   - Kullanıcı rolüne göre alan gösterimi
//   - Değer bazlı koşullu görünürlük
//
// # Örnek
//
//	if field.IsConditionallyVisible(user) {
//	    // Alanı göster
//	}
func (s *Schema) IsConditionallyVisible(item interface{}) bool {
	// If there's a visibility callback, use it
	if s.VisibilityCallback != nil {
		// Create a minimal ResourceContext for the callback
		ctx := &core.ResourceContext{}
		return s.VisibilityCallback(ctx)
	}
	return true
}

// GetMetadata, alan hakkında metadata bilgilerini döner.
//
// Bu metod, alanın tüm özelliklerini içeren bir map döner.
// API yanıtları, debugging ve introspection için kullanılır.
//
// # Döndürür
//
//   - map[string]interface{}: Alan metadata'sı
//
// # İçerilen Bilgiler
//
//   - name: Görünen ad
//   - key: Benzersiz tanımlayıcı
//   - view: Frontend bileşeni
//   - type: Veri tipi
//   - context: Görüntüleme bağlamı
//   - read_only: Salt okunur durumu
//   - disabled: Devre dışı durumu
//   - immutable: Değiştirilemez durumu
//   - required: Zorunlu durumu
//   - nullable: Boş bırakılabilir durumu
//   - filterable: Filtrelenebilir durumu
//   - sortable: Sıralanabilir durumu
//   - searchable: Aranabilir durumu
//   - stacked: Yığılmış durumu
//   - text_align: Metin hizalama
//   - dependencies: Bağımlılıklar
//   - props: Ekstra özellikler
//
// # Örnek
//
//	metadata := field.GetMetadata()
//	fmt.Printf("Field: %s, Type: %s\n", metadata["name"], metadata["type"])
func (s *Schema) GetMetadata() map[string]interface{} {
	metadata := make(map[string]interface{})
	metadata["name"] = s.Name
	metadata["key"] = s.Key
	metadata["view"] = s.View
	metadata["type"] = s.Type
	metadata["context"] = s.Context
	metadata["read_only"] = s.IsReadOnly
	metadata["disabled"] = s.IsDisabled
	metadata["immutable"] = s.IsImmutable
	metadata["required"] = s.IsRequired
	metadata["nullable"] = s.IsNullable
	metadata["filterable"] = s.IsFilterable
	metadata["sortable"] = s.IsSortable
	metadata["searchable"] = s.GlobalSearch
	metadata["stacked"] = s.IsStacked
	metadata["text_align"] = s.TextAlign
	metadata["dependencies"] = s.GetDependencies()
	metadata["props"] = s.Props

	return metadata
}

// IsVisibleInContext, alanın belirli bir bağlamda görünür olup olmadığını kontrol eder.
//
// Bu metod, ElementContext değerlerine göre görünürlük kontrolü yapar.
// Her bağlam için farklı kurallar uygulanır.
//
// # Parametreler
//
//   - ctx: Görünürlük bağlamı
//
// # Döndürür
//
//   - bool: Alanın görünür olup olmadığı
//
// # Bağlam Kuralları
//
// **ContextIndex (Liste Görünümü)**:
//   - HIDE_ON_LIST ise gizli
//   - ONLY_ON_DETAIL ise gizli
//   - ONLY_ON_FORM ise gizli
//   - Diğer durumlarda görünür
//
// **ContextDetail (Detay Görünümü)**:
//   - HIDE_ON_DETAIL ise gizli
//   - ONLY_ON_LIST ise gizli
//   - ONLY_ON_FORM ise gizli
//   - Diğer durumlarda görünür
//
// **ContextCreate (Oluşturma Formu)**:
//   - HIDE_ON_CREATE ise gizli
//   - ONLY_ON_UPDATE ise gizli
//   - ONLY_ON_LIST ise gizli
//   - ONLY_ON_DETAIL ise gizli
//   - Diğer durumlarda görünür
//
// **ContextUpdate (Güncelleme Formu)**:
//   - HIDE_ON_UPDATE ise gizli
//   - ONLY_ON_CREATE ise gizli
//   - ONLY_ON_LIST ise gizli
//   - ONLY_ON_DETAIL ise gizli
//   - Diğer durumlarda görünür
//
// **ContextPreview (Önizleme)**:
//   - HIDE_ON_DETAIL ise gizli
//   - ONLY_ON_LIST ise gizli
//   - ONLY_ON_FORM ise gizli
//   - Diğer durumlarda görünür
//
// # Örnek
//
//	if field.IsVisibleInContext(ContextCreate) {
//	    // Alanı oluşturma formunda göster
//	}
func (s *Schema) IsVisibleInContext(ctx VisibilityContext) bool {
	// Map VisibilityContext to ElementContext for compatibility
	switch ctx {
	case ContextIndex:
		return s.Context != HIDE_ON_LIST && s.Context != ONLY_ON_DETAIL && s.Context != ONLY_ON_FORM
	case ContextDetail:
		return s.Context != HIDE_ON_DETAIL && s.Context != ONLY_ON_LIST && s.Context != ONLY_ON_FORM
	case ContextCreate:
		return s.Context != HIDE_ON_CREATE && s.Context != ONLY_ON_UPDATE && s.Context != ONLY_ON_LIST && s.Context != ONLY_ON_DETAIL
	case ContextUpdate:
		return s.Context != HIDE_ON_UPDATE && s.Context != ONLY_ON_CREATE && s.Context != ONLY_ON_LIST && s.Context != ONLY_ON_DETAIL
	case ContextPreview:
		return s.Context != HIDE_ON_DETAIL && s.Context != ONLY_ON_LIST && s.Context != ONLY_ON_FORM
	default:
		return true
	}
}

// Validation Fluent API Methods - Doğrulama Kuralları
//
// Bu bölüm, alana doğrulama kuralları eklemek için kullanılan metodları içerir.
// Her metod, ValidationRule oluşturur ve alanın ValidationRules listesine ekler.

// AddValidationRule, alana bir doğrulama kuralı ekler.
//
// Bu metod, özel veya hazır doğrulama kurallarını alana eklemek için kullanılır.
// ValidationRule interface'ini implement eden herhangi bir kural eklenebilir.
//
// # Parametreler
//
//   - rule: Eklenecek doğrulama kuralı (ValidationRule interface'ini implement etmeli)
//
// # Döndürür
//
//   - core.Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Kullanım Senaryoları
//
//   - Özel doğrulama kuralları ekleme
//   - Hazır kural fonksiyonlarını kullanma
//   - Çoklu doğrulama kuralları birleştirme
//
// # Örnek
//
//	field := Text("age", "Yaş").
//	    AddValidationRule(Min(18)).
//	    AddValidationRule(Max(100))
func (s *Schema) AddValidationRule(rule interface{}) core.Element {
	if vr, ok := rule.(ValidationRule); ok {
		s.ValidationRules = append(s.ValidationRules, vr)
	}
	return s
}

// Email, e-posta doğrulama kuralı ekler.
//
// Bu metod, alanın geçerli bir e-posta adresi içermesini zorunlu kılar.
// RFC 5322 standardına göre e-posta formatı kontrol edilir.
//
// # Döndürür
//
//   - core.Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Doğrulama Kuralları
//
//   - @ işareti içermeli
//   - Geçerli domain formatı
//   - Geçerli local part formatı
//
// # Kullanım Senaryoları
//
//   - Kullanıcı kayıt formları
//   - İletişim formları
//   - Profil güncelleme formları
//
// # Örnek
//
//	field := Text("email", "E-posta").Required().Email()
func (s *Schema) Email() core.Element {
	return s.AddValidationRule(EmailRule())
}

// URL, URL doğrulama kuralı ekler.
//
// Bu metod, alanın geçerli bir URL içermesini zorunlu kılar.
// HTTP, HTTPS ve diğer protokoller desteklenir.
//
// # Döndürür
//
//   - core.Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Doğrulama Kuralları
//
//   - Geçerli URL formatı
//   - Protokol içermeli (http://, https://, vb.)
//   - Geçerli domain formatı
//
// # Kullanım Senaryoları
//
//   - Website alanları
//   - Sosyal medya profil linkleri
//   - API endpoint alanları
//
// # Örnek
//
//	field := Text("website", "Website").URL()
func (s *Schema) URL() core.Element {
	return s.AddValidationRule(URL())
}

// Min, minimum değer doğrulama kuralı ekler.
//
// Bu metod, alanın belirtilen minimum değerden büyük veya eşit olmasını zorunlu kılar.
// Sayısal alanlar için minimum değer, string alanlar için minimum uzunluk kontrolü yapar.
//
// # Parametreler
//
//   - min: Minimum değer (int, float64, string uzunluğu)
//
// # Döndürür
//
//   - core.Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Kullanım Senaryoları
//
//   - Yaş kontrolü (min 18)
//   - Fiyat kontrolü (min 0)
//   - Miktar kontrolü (min 1)
//   - String uzunluk kontrolü
//
// # Örnek
//
//	field := Number("age", "Yaş").Min(18).Max(100)
func (s *Schema) Min(min interface{}) core.Element {
	return s.AddValidationRule(Min(min))
}

// Max, maksimum değer doğrulama kuralı ekler.
//
// Bu metod, alanın belirtilen maksimum değerden küçük veya eşit olmasını zorunlu kılar.
// Sayısal alanlar için maksimum değer, string alanlar için maksimum uzunluk kontrolü yapar.
//
// # Parametreler
//
//   - max: Maksimum değer (int, float64, string uzunluğu)
//
// # Döndürür
//
//   - core.Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Kullanım Senaryoları
//
//   - Yaş kontrolü (max 100)
//   - Fiyat kontrolü (max 999999)
//   - Miktar kontrolü (max 1000)
//   - String uzunluk kontrolü
//
// # Örnek
//
//	field := Number("quantity", "Miktar").Min(1).Max(1000)
func (s *Schema) Max(max interface{}) core.Element {
	return s.AddValidationRule(Max(max))
}

// MinLength, minimum uzunluk doğrulama kuralı ekler.
//
// Bu metod, string alanların belirtilen minimum karakter sayısına sahip olmasını zorunlu kılar.
// Boşluklar da karakter olarak sayılır.
//
// # Parametreler
//
//   - length: Minimum karakter sayısı
//
// # Döndürür
//
//   - core.Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Kullanım Senaryoları
//
//   - Şifre uzunluğu (min 8 karakter)
//   - Kullanıcı adı uzunluğu (min 3 karakter)
//   - Açıklama uzunluğu (min 10 karakter)
//
// # Örnek
//
//	field := Password("password", "Şifre").MinLength(8).MaxLength(128)
func (s *Schema) MinLength(length int) core.Element {
	return s.AddValidationRule(MinLength(length))
}

// MaxLength, maksimum uzunluk doğrulama kuralı ekler.
//
// Bu metod, string alanların belirtilen maksimum karakter sayısını aşmamasını zorunlu kılar.
// Boşluklar da karakter olarak sayılır.
//
// # Parametreler
//
//   - length: Maksimum karakter sayısı
//
// # Döndürür
//
//   - core.Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Kullanım Senaryoları
//
//   - Başlık uzunluğu (max 255 karakter)
//   - Açıklama uzunluğu (max 1000 karakter)
//   - Kullanıcı adı uzunluğu (max 50 karakter)
//
// # Örnek
//
//	field := Text("title", "Başlık").MaxLength(255)
func (s *Schema) MaxLength(length int) core.Element {
	return s.AddValidationRule(MaxLength(length))
}

// Pattern, regex pattern doğrulama kuralı ekler.
//
// Bu metod, alanın belirtilen regex pattern'ine uymasını zorunlu kılar.
// Özel format kontrolleri için kullanılır.
//
// # Parametreler
//
//   - pattern: Regex pattern string'i
//
// # Döndürür
//
//   - core.Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Kullanım Senaryoları
//
//   - Telefon numarası formatı
//   - Posta kodu formatı
//   - Özel kod formatları
//   - Kimlik numarası formatı
//
// # Örnek
//
//	field := Text("phone", "Telefon").
//	    Pattern(`^(\+90|0)?[0-9]{10}$`).
//	    Placeholder("5XX XXX XX XX")
func (s *Schema) Pattern(pattern string) core.Element {
	return s.AddValidationRule(Pattern(pattern))
}

// Unique, benzersizlik doğrulama kuralı ekler.
//
// Bu metod, alanın veritabanında benzersiz olmasını zorunlu kılar.
// Aynı değere sahip başka bir kayıt varsa hata verir.
//
// # Parametreler
//
//   - table: Kontrol edilecek tablo adı
//   - column: Kontrol edilecek sütun adı
//
// # Döndürür
//
//   - core.Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Kullanım Senaryoları
//
//   - E-posta benzersizliği
//   - Kullanıcı adı benzersizliği
//   - Slug benzersizliği
//   - Kod benzersizliği
//
// # Önemli Notlar
//
//   - Güncelleme işlemlerinde mevcut kaydın kendisi hariç tutulur
//   - Performans için veritabanı indeksi önerilir
//   - Büyük/küçük harf duyarlılığı veritabanına bağlıdır
//
// # Örnek
//
//	field := Text("email", "E-posta").
//	    Email().
//	    Unique("users", "email")
func (s *Schema) Unique(table, column string) core.Element {
	return s.AddValidationRule(Unique(table, column))
}

// Exists, varlık doğrulama kuralı ekler.
//
// Bu metod, alanın değerinin belirtilen tabloda mevcut olmasını zorunlu kılar.
// Foreign key ilişkileri için kullanılır.
//
// # Parametreler
//
//   - table: Kontrol edilecek tablo adı
//   - column: Kontrol edilecek sütun adı
//
// # Döndürür
//
//   - core.Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Kullanım Senaryoları
//
//   - Foreign key doğrulama
//   - İlişkili kayıt kontrolü
//   - Referans doğrulama
//
// # Önemli Notlar
//
//   - Silinen kayıtlar (soft delete) için dikkatli kullanılmalı
//   - Performans için veritabanı indeksi önerilir
//   - NULL değerler için ayrıca Nullable() kullanılmalı
//
// # Örnek
//
//	field := BelongsTo("category_id", "Kategori").
//	    Required().
//	    Exists("categories", "id")
func (s *Schema) Exists(table, column string) core.Element {
	return s.AddValidationRule(Exists(table, column))
}

// Display Fluent API Methods - Görüntüleme Özelleştirme
//
// Bu bölüm, alanın nasıl görüntüleneceğini özelleştirmek için kullanılan metodları içerir.
// Değer formatı, etiket gösterimi ve client-side etkileşim için kullanılır.

// Display, alan için özel görüntüleme callback'i ayarlar.
//
// Bu metod, alanın değerinin nasıl görüntüleneceğini özelleştirmeye olanak tanır.
// Callback, değeri alır ve görüntülenecek string'i döner.
//
// # Parametreler
//
//   - fn: Görüntüleme işlemini yapan fonksiyon (interface{} -> string)
//
// # Döndürür
//
//   - core.Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Kullanım Senaryoları
//
//   - Tarih formatı özelleştirme
//   - Para birimi formatı
//   - Boolean değerleri metin olarak gösterme
//   - Enum değerlerini açıklayıcı metne çevirme
//
// # Örnek
//
//	field := Boolean("is_active", "Aktif").Display(func(value interface{}) string {
//	    if value.(bool) {
//	        return "Aktif ✓"
//	    }
//	    return "Pasif ✗"
//	})
func (s *Schema) Display(fn func(interface{}) string) core.Element {
	s.DisplayCallback = fn
	return s
}

// DisplayAs, alan için görüntüleme format string'i ayarlar.
//
// Bu metod, alanın değerinin belirli bir formatta görüntülenmesini sağlar.
// Format string'i, frontend tarafından yorumlanır.
//
// # Parametreler
//
//   - format: Görüntüleme format string'i
//
// # Döndürür
//
//   - core.Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Kullanım Senaryoları
//
//   - Tarih formatı belirleme ("DD/MM/YYYY")
//   - Para birimi formatı ("$0,0.00")
//   - Yüzde formatı ("%0.2f")
//
// # Örnek
//
//	field := Number("price", "Fiyat").DisplayAs("$0,0.00")
func (s *Schema) DisplayAs(format string) core.Element {
	s.DisplayedAs = format
	return s
}

// DisplayUsingLabels, alanın etiketler kullanarak görüntüleneceğini belirtir.
//
// Bu metod, select, radio gibi seçim alanlarında değer yerine etiketin gösterilmesini sağlar.
// Örneğin, "1" yerine "Aktif" gösterilir.
//
// # Döndürür
//
//   - core.Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Kullanım Senaryoları
//
//   - Select alanlarında etiket gösterimi
//   - Enum değerlerinin açıklayıcı gösterimi
//   - İlişkili kayıtların isim gösterimi
//
// # Örnek
//
//	field := Select("status", "Durum").
//	    Options([]map[string]interface{}{
//	        {"value": 1, "label": "Aktif"},
//	        {"value": 0, "label": "Pasif"},
//	    }).
//	    DisplayUsingLabels()
func (s *Schema) DisplayUsingLabels() core.Element {
	s.DisplayUsingLabelsFlag = true
	return s
}

// ResolveHandle, client-side bileşen etkileşimi için handle ayarlar.
//
// Bu metod, frontend bileşeninin özel işlemler için kullanacağı bir handle belirler.
// Genellikle özel bileşenler veya dinamik davranışlar için kullanılır.
//
// # Parametreler
//
//   - handle: Resolve handle string'i
//
// # Döndürür
//
//   - core.Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Kullanım Senaryoları
//
//   - Özel bileşen etkileşimi
//   - Dinamik veri yükleme
//   - Client-side işlem tetikleme
//
// # Örnek
//
//	field := Custom("map", "Harita").ResolveHandle("map-component")
func (s *Schema) ResolveHandle(handle string) core.Element {
	s.ResolveHandleValue = handle
	return s
}

// Dependency Fluent API Methods - Bağımlılık Yönetimi
//
// Bu bölüm, alanlar arası bağımlılıkları yönetmek için kullanılan metodları içerir.
// Dinamik form davranışları ve koşullu alan gösterimi için kullanılır.

// DependsOn, alanın bağımlı olduğu diğer alanları belirtir.
//
// Bu metod, alanın görünürlüğünün veya davranışının diğer alanlara bağlı olduğunu belirtir.
// Bağımlı alanlar değiştiğinde, bu alan da güncellenir.
//
// # Parametreler
//
//   - fields: Bağımlı olunan alan adları (variadic)
//
// # Döndürür
//
//   - core.Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Kullanım Senaryoları
//
//   - Cascade seçim alanları (ülke -> şehir -> ilçe)
//   - Koşullu alan gösterimi
//   - Dinamik form yapıları
//   - İlişkili alan güncellemeleri
//
// # Örnek
//
//	field := Select("city", "Şehir").
//	    DependsOn("country").
//	    When("country", "=", "TR")
func (s *Schema) DependsOn(fields ...string) core.Element {
	s.DependsOnFields = fields
	return s
}

// When, bağımlılık kuralı ekler.
//
// Bu metod, belirli bir alanın belirli bir değere sahip olması durumunda
// bu alanın nasıl davranacağını belirler.
//
// # Parametreler
//
//   - field: Kontrol edilecek alan adı
//   - operator: Karşılaştırma operatörü ("=", "!=", ">", "<", ">=", "<=", "in", "not_in")
//   - value: Karşılaştırılacak değer
//
// # Döndürür
//
//   - core.Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Kullanım Senaryoları
//
//   - Koşullu alan gösterimi
//   - Değer bazlı validasyon
//   - Dinamik form davranışı
//
// # Örnek
//
//	field := Text("other_reason", "Diğer Sebep").
//	    DependsOn("reason").
//	    When("reason", "=", "other")
func (s *Schema) When(field string, operator string, value interface{}) core.Element {
	if s.DependencyRules == nil {
		s.DependencyRules = make(map[string]interface{})
	}
	s.DependencyRules[field] = map[string]interface{}{
		"operator": operator,
		"value":    value,
	}
	return s
}

// OnDependencyChange, bağımlı alanlar değiştiğinde çağrılacak callback ayarlar.
//
// Bu metod, bağımlı alanların değeri değiştiğinde özel işlemler yapılmasını sağlar.
// Hem oluşturma hem de güncelleme bağlamında çalışır.
//
// # Parametreler
//
//   - fn: Bağımlılık değişikliği callback fonksiyonu
//
// # Döndürür
//
//   - core.Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Kullanım Senaryoları
//
//   - Dinamik seçenek yükleme
//   - Hesaplanan değerler
//   - Koşullu validasyon
//   - Alan değeri senkronizasyonu
//
// # Örnek
//
//	field := Select("city", "Şehir").
//	    DependsOn("country").
//	    OnDependencyChange(func(deps map[string]interface{}, ctx *fiber.Ctx) interface{} {
//	        countryID := deps["country"]
//	        return loadCitiesByCountry(countryID)
//	    })
func (s *Schema) OnDependencyChange(fn DependencyCallbackFunc) core.Element {
	s.DependencyCallback = fn
	return s
}

// OnDependencyChangeCreating, oluşturma bağlamında bağımlılık callback'i ayarlar.
//
// Bu metod, sadece yeni kayıt oluşturulurken bağımlılık değişikliklerinde çalışır.
// Güncelleme işlemlerinde çalışmaz.
//
// # Parametreler
//
//   - fn: Bağımlılık değişikliği callback fonksiyonu
//
// # Döndürür
//
//   - core.Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Kullanım Senaryoları
//
//   - Oluşturma sırasında özel davranış
//   - İlk değer hesaplama
//   - Oluşturma bağlamına özel validasyon
//
// # Örnek
//
//	field := Text("slug", "Slug").
//	    DependsOn("title").
//	    OnDependencyChangeCreating(func(deps map[string]interface{}, ctx *fiber.Ctx) interface{} {
//	        return generateSlug(deps["title"].(string))
//	    })
func (s *Schema) OnDependencyChangeCreating(fn DependencyCallbackFunc) core.Element {
	s.DependencyCallbackOnCreate = fn
	return s
}

// OnDependencyChangeUpdating, güncelleme bağlamında bağımlılık callback'i ayarlar.
//
// Bu metod, sadece kayıt güncellenirken bağımlılık değişikliklerinde çalışır.
// Oluşturma işlemlerinde çalışmaz.
//
// # Parametreler
//
//   - fn: Bağımlılık değişikliği callback fonksiyonu
//
// # Döndürür
//
//   - core.Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Kullanım Senaryoları
//
//   - Güncelleme sırasında özel davranış
//   - Değer yeniden hesaplama
//   - Güncelleme bağlamına özel validasyon
//
// # Örnek
//
//	field := DateTime("updated_at", "Güncelleme Tarihi").
//	    OnDependencyChangeUpdating(func(deps map[string]interface{}, ctx *fiber.Ctx) interface{} {
//	        return time.Now()
//	    })
func (s *Schema) OnDependencyChangeUpdating(fn DependencyCallbackFunc) core.Element {
	s.DependencyCallbackOnUpdate = fn
	return s
}

// GetDependencyCallback, bağlama göre uygun callback'i döner.
//
// Bu metod, verilen bağlama (create/update) göre uygun callback fonksiyonunu döner.
// Bağlama özel callback yoksa, genel callback döner.
//
// # Parametreler
//
//   - context: Bağlam string'i ("create" veya "update")
//
// # Döndürür
//
//   - DependencyCallbackFunc: Uygun callback fonksiyonu (nil olabilir)
//
// # Örnek
//
//	callback := field.GetDependencyCallback("create")
//	if callback != nil {
//	    result := callback(dependencies, ctx)
//	}
func (s *Schema) GetDependencyCallback(context string) DependencyCallbackFunc {
	switch context {
	case "create":
		if s.DependencyCallbackOnCreate != nil {
			return s.DependencyCallbackOnCreate
		}
	case "update":
		if s.DependencyCallbackOnUpdate != nil {
			return s.DependencyCallbackOnUpdate
		}
	}
	return s.DependencyCallback
}

// Suggestion Fluent API Methods - Öneri Sistemi
//
// Bu bölüm, alan için otomatik tamamlama ve öneri özelliklerini yönetir.
// Kullanıcı deneyimini iyileştirmek için kullanılır.

// WithSuggestions, alan için öneri callback'i ayarlar.
//
// Bu metod, kullanıcı yazarken gösterilecek önerileri dinamik olarak oluşturur.
// Callback, arama sorgusunu alır ve öneri listesi döner.
//
// # Parametreler
//
//   - fn: Öneri oluşturma fonksiyonu (string -> []interface{})
//
// # Döndürür
//
//   - core.Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Kullanım Senaryoları
//
//   - Otomatik tamamlama
//   - Arama önerileri
//   - Hızlı seçim listeleri
//   - Dinamik filtreleme
//
// # Örnek
//
//	field := Text("email", "E-posta").
//	    WithSuggestions(func(query string) []interface{} {
//	        return searchEmails(query)
//	    }).
//	    MinCharsForSuggestions(3)
func (s *Schema) WithSuggestions(fn func(string) []interface{}) core.Element {
	s.SuggestionsCallback = fn
	return s
}

// WithAutoComplete, alan için otomatik tamamlama URL'i ayarlar.
//
// Bu metod, önerilerin bir API endpoint'inden yüklenmesini sağlar.
// Callback yerine URL kullanarak server-side öneri sistemi oluşturur.
//
// # Parametreler
//
//   - url: Otomatik tamamlama API endpoint'i
//
// # Döndürür
//
//   - core.Element: Zincirleme çağrılar için Schema pointer'ı
//
// # API Beklentileri
//
//   - GET request ile çağrılır
//   - Query parameter: ?q=arama_metni
//   - Response: JSON array [{value, label}, ...]
//
// # Kullanım Senaryoları
//
//   - Büyük veri setlerinde arama
//   - Uzak API entegrasyonu
//   - Performans optimizasyonu
//
// # Örnek
//
//	field := Text("city", "Şehir").
//	    WithAutoComplete("/api/cities/search").
//	    MinCharsForSuggestions(2)
func (s *Schema) WithAutoComplete(url string) core.Element {
	s.AutoCompleteURL = url
	return s
}

// MinCharsForSuggestions, öneri gösterimi için minimum karakter sayısı ayarlar.
//
// Bu metod, kullanıcının kaç karakter yazdıktan sonra önerilerin gösterileceğini belirler.
// Performans optimizasyonu için kullanılır.
//
// # Parametreler
//
//   - min: Minimum karakter sayısı
//
// # Döndürür
//
//   - core.Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Önerilen Değerler
//
//   - 1-2: Küçük veri setleri için
//   - 3-4: Orta veri setleri için (önerilen)
//   - 5+: Büyük veri setleri için
//
// # Örnek
//
//	field := Text("product", "Ürün").
//	    WithAutoComplete("/api/products/search").
//	    MinCharsForSuggestions(3)
func (s *Schema) MinCharsForSuggestions(min int) core.Element {
	s.MinCharsForSuggestionsVal = min
	return s
}

// Attachment Fluent API Methods - Dosya Yükleme Yönetimi
//
// Bu bölüm, dosya yükleme alanları için kullanılan metodları içerir.
// Dosya tipi, boyut, depolama ve işleme ayarları için kullanılır.

// Accept, kabul edilen MIME tiplerini ayarlar.
//
// Bu metod, hangi dosya tiplerinin yüklenebileceğini belirler.
// Frontend'de dosya seçici bu tiplere göre filtrelenir.
//
// # Parametreler
//
//   - mimeTypes: Kabul edilen MIME tipleri (variadic)
//
// # Döndürür
//
//   - core.Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Yaygın MIME Tipleri
//
//   - Resimler: "image/jpeg", "image/png", "image/gif", "image/webp"
//   - Dökümanlar: "application/pdf", "application/msword"
//   - Video: "video/mp4", "video/webm"
//   - Audio: "audio/mpeg", "audio/wav"
//
// # Kullanım Senaryoları
//
//   - Profil fotoğrafı yükleme
//   - Döküman yükleme
//   - Medya dosyası yükleme
//
// # Örnek
//
//	field := File("avatar", "Profil Fotoğrafı").
//	    Accept("image/jpeg", "image/png", "image/webp").
//	    MaxSize(5 * 1024 * 1024) // 5MB
func (s *Schema) Accept(mimeTypes ...string) core.Element {
	s.AcceptedMimeTypes = append(s.AcceptedMimeTypes, mimeTypes...)
	return s
}

// MaxSize, maksimum dosya boyutunu ayarlar.
//
// Bu metod, yüklenebilecek dosyanın maksimum boyutunu byte cinsinden belirler.
// Daha büyük dosyalar reddedilir.
//
// # Parametreler
//
//   - bytes: Maksimum dosya boyutu (byte cinsinden)
//
// # Döndürür
//
//   - core.Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Önerilen Boyutlar
//
//   - Profil fotoğrafı: 2-5 MB
//   - Döküman: 10-50 MB
//   - Video: 100-500 MB
//
// # Boyut Hesaplama
//
//   - 1 KB = 1024 bytes
//   - 1 MB = 1024 * 1024 bytes
//   - 1 GB = 1024 * 1024 * 1024 bytes
//
// # Örnek
//
//	field := File("document", "Döküman").
//	    Accept("application/pdf").
//	    MaxSize(10 * 1024 * 1024) // 10MB
func (s *Schema) MaxSize(bytes int64) core.Element {
	s.MaxFileSize = bytes
	return s
}

// Store, dosya depolama ayarlarını belirler.
//
// Bu metod, yüklenen dosyaların nerede saklanacağını belirler.
// Disk ve path bilgilerini ayarlar.
//
// # Parametreler
//
//   - disk: Depolama disk adı (örn. "local", "s3", "gcs")
//   - path: Depolama yolu (örn. "uploads/avatars", "documents")
//
// # Döndürür
//
//   - core.Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Disk Tipleri
//
//   - **local**: Yerel dosya sistemi
//   - **s3**: Amazon S3
//   - **gcs**: Google Cloud Storage
//   - **azure**: Azure Blob Storage
//
// # Kullanım Senaryoları
//
//   - Yerel depolama
//   - Cloud depolama
//   - CDN entegrasyonu
//
// # Örnek
//
//	field := File("avatar", "Avatar").
//	    Store("s3", "uploads/avatars").
//	    Accept("image/*").
//	    MaxSize(5 * 1024 * 1024)
func (s *Schema) Store(disk, path string) core.Element {
	s.StorageDisk = disk
	s.StoragePath = path
	return s
}

// WithUpload, özel dosya yükleme callback'i ayarlar.
//
// Bu metod, dosya yükleme işlemini özelleştirmeye olanak tanır.
// Dosya işleme, dönüşüm, validasyon gibi işlemler için kullanılır.
//
// # Parametreler
//
//   - fn: Dosya yükleme işlemini yapan fonksiyon
//
// # Döndürür
//
//   - core.Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Kullanım Senaryoları
//
//   - Resim boyutlandırma
//   - Dosya dönüşümü
//   - Özel validasyon
//   - Metadata çıkarma
//
// # Örnek
//
//	field := File("image", "Resim").
//	    WithUpload(func(file interface{}, item interface{}) error {
//	        // Resmi yeniden boyutlandır
//	        return resizeImage(file, 800, 600)
//	    })
func (s *Schema) WithUpload(fn func(interface{}, interface{}) error) core.Element {
	s.UploadCallback = fn
	return s
}

// MarkRemoveEXIFData, EXIF verilerinin kaldırılacağını belirtir.
//
// Bu metod, yüklenen resimlerdeki EXIF metadata'sının (konum, kamera bilgisi, vb.)
// otomatik olarak kaldırılmasını sağlar. Gizlilik için önemlidir.
//
// # Döndürür
//
//   - core.Element: Zincirleme çağrılar için Schema pointer'ı
//
// # EXIF Verileri
//
//   - GPS koordinatları
//   - Kamera modeli
//   - Çekim tarihi
//   - Kamera ayarları
//
// # Kullanım Senaryoları
//
//   - Gizlilik koruması
//   - Güvenlik
//   - Dosya boyutu azaltma
//
// # Örnek
//
//	field := File("photo", "Fotoğraf").
//	    Accept("image/jpeg", "image/png").
//	    MarkRemoveEXIFData()
func (s *Schema) MarkRemoveEXIFData() core.Element {
	s.RemoveEXIFDataFlag = true
	return s
}

// Repeater Fluent API Methods - Tekrarlayan Alan Yönetimi
//
// Bu bölüm, tekrarlayan alan grupları için kullanılan metodları içerir.
// Dinamik form alanları ve çoklu veri girişi için kullanılır.

// Fields, repeater için alt alanları ayarlar.
//
// Bu metod, tekrarlayan grup içindeki alanları belirler.
// Her tekrar, bu alanların bir kopyasını içerir.
//
// # Parametreler
//
//   - fields: Tekrarlayan alan listesi (variadic)
//
// # Döndürür
//
//   - core.Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Kullanım Senaryoları
//
//   - Telefon numaraları listesi
//   - E-posta adresleri listesi
//   - Ürün varyantları
//   - Adres bilgileri
//
// # Örnek
//
//	field := Repeater("phones", "Telefon Numaraları").
//	    Fields(
//	        Text("number", "Numara").Required(),
//	        Select("type", "Tip").Options([]string{"Ev", "İş", "Mobil"}),
//	    ).
//	    MinRepeats(1).
//	    MaxRepeats(5)
func (s *Schema) Fields(fields ...core.Element) core.Element {
	s.RepeaterFields = fields
	return s
}

// MinRepeats, minimum tekrar sayısını ayarlar.
//
// Bu metod, repeater'da en az kaç tane alan grubu olması gerektiğini belirler.
// Kullanıcı bu sayının altına inemez.
//
// # Parametreler
//
//   - min: Minimum tekrar sayısı
//
// # Döndürür
//
//   - core.Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Kullanım Senaryoları
//
//   - En az bir telefon numarası zorunlu
//   - En az bir adres zorunlu
//   - Minimum varyant sayısı
//
// # Örnek
//
//	field := Repeater("addresses", "Adresler").
//	    Fields(
//	        Text("street", "Sokak"),
//	        Text("city", "Şehir"),
//	    ).
//	    MinRepeats(1)
func (s *Schema) MinRepeats(min int) core.Element {
	s.MinRepeatsCount = min
	return s
}

// MaxRepeats, maksimum tekrar sayısını ayarlar.
//
// Bu metod, repeater'da en fazla kaç tane alan grubu olabileceğini belirler.
// Kullanıcı bu sayının üstüne çıkamaz.
//
// # Parametreler
//
//   - max: Maksimum tekrar sayısı
//
// # Döndürür
//
//   - core.Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Kullanım Senaryoları
//
//   - Maksimum telefon sayısı sınırı
//   - Maksimum adres sayısı sınırı
//   - Performans optimizasyonu
//
// # Örnek
//
//	field := Repeater("emails", "E-posta Adresleri").
//	    Fields(
//	        Text("email", "E-posta").Email(),
//	    ).
//	    MinRepeats(1).
//	    MaxRepeats(3)
func (s *Schema) MaxRepeats(max int) core.Element {
	s.MaxRepeatsCount = max
	return s
}

// RichText Fluent API Methods - Zengin Metin Editörü Ayarları
//
// Bu bölüm, zengin metin editörü alanları için kullanılan metodları içerir.
// Editör tipi, dil ve tema ayarları için kullanılır.

// WithEditor, editör tipini ayarlar.
//
// Bu metod, hangi zengin metin editörünün kullanılacağını belirler.
// Farklı editörler farklı özellikler sunar.
//
// # Parametreler
//
//   - editorType: Editör tipi ("tinymce", "quill", "ckeditor", "markdown")
//
// # Döndürür
//
//   - core.Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Editör Tipleri
//
//   - **tinymce**: Tam özellikli WYSIWYG editör
//   - **quill**: Modern, hafif editör
//   - **ckeditor**: Güçlü, özelleştirilebilir editör
//   - **markdown**: Markdown formatı editörü
//
// # Örnek
//
//	field := RichText("content", "İçerik").
//	    WithEditor("tinymce").
//	    WithLanguage("tr").
//	    WithTheme("modern")
func (s *Schema) WithEditor(editorType string) core.Element {
	s.EditorType = editorType
	return s
}

// WithLanguage, editör dilini ayarlar.
//
// Bu metod, editör arayüzünün hangi dilde gösterileceğini belirler.
// Menüler, butonlar ve mesajlar bu dilde gösterilir.
//
// # Parametreler
//
//   - language: Dil kodu ("tr", "en", "de", "fr", vb.)
//
// # Döndürür
//
//   - core.Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Desteklenen Diller
//
//   - tr: Türkçe
//   - en: İngilizce
//   - de: Almanca
//   - fr: Fransızca
//   - es: İspanyolca
//
// # Örnek
//
//	field := RichText("description", "Açıklama").
//	    WithEditor("quill").
//	    WithLanguage("tr")
func (s *Schema) WithLanguage(language string) core.Element {
	s.EditorLanguage = language
	return s
}

// WithTheme, editör temasını ayarlar.
//
// Bu metod, editörün görsel temasını belirler.
// Farklı temalar farklı görünümler sunar.
//
// # Parametreler
//
//   - theme: Tema adı ("modern", "classic", "dark", "light")
//
// # Döndürür
//
//   - core.Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Yaygın Temalar
//
//   - **modern**: Modern, minimal görünüm
//   - **classic**: Klasik, geleneksel görünüm
//   - **dark**: Koyu tema
//   - **light**: Açık tema
//
// # Örnek
//
//	field := RichText("article", "Makale").
//	    WithEditor("ckeditor").
//	    WithTheme("dark")
func (s *Schema) WithTheme(theme string) core.Element {
	s.EditorTheme = theme
	return s
}

// Status Fluent API Methods - Durum Badge Yönetimi
//
// Bu bölüm, durum alanları için badge renkleri ve varyantları yönetir.
// Görsel durum gösterimi için kullanılır.

// WithStatusColors, durum renk eşleştirmesi ayarlar.
//
// Bu metod, her durum değeri için bir renk tanımlar.
// Liste ve detay görünümlerinde renkli badge'ler gösterilir.
//
// # Parametreler
//
//   - colors: Durum-renk eşleştirme map'i
//
// # Döndürür
//
//   - core.Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Renk Değerleri
//
//   - Hex kodları: "#FF0000", "#00FF00"
//   - Renk isimleri: "red", "green", "blue"
//   - Tailwind sınıfları: "bg-red-500", "bg-green-500"
//
// # Kullanım Senaryoları
//
//   - Sipariş durumları
//   - Ödeme durumları
//   - Kullanıcı durumları
//   - Görev durumları
//
// # Örnek
//
//	field := Status("status", "Durum").
//	    WithStatusColors(map[string]string{
//	        "pending":   "yellow",
//	        "approved":  "green",
//	        "rejected":  "red",
//	        "cancelled": "gray",
//	    })
func (s *Schema) WithStatusColors(colors map[string]string) core.Element {
	s.StatusColors = colors
	return s
}

// WithBadgeVariant, badge varyantını ayarlar.
//
// Bu metod, badge'in görsel stilini belirler.
// Farklı varyantlar farklı görünümler sunar.
//
// # Parametreler
//
//   - variant: Badge varyantı ("solid", "outline", "subtle", "dot")
//
// # Döndürür
//
//   - core.Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Varyant Tipleri
//
//   - **solid**: Dolu, renkli arka plan
//   - **outline**: Sadece kenarlık
//   - **subtle**: Hafif renkli arka plan
//   - **dot**: Nokta ile gösterim
//
// # Örnek
//
//	field := Status("priority", "Öncelik").
//	    WithStatusColors(map[string]string{
//	        "high":   "red",
//	        "medium": "yellow",
//	        "low":    "green",
//	    }).
//	    WithBadgeVariant("solid")
func (s *Schema) WithBadgeVariant(variant string) core.Element {
	s.BadgeVariant = variant
	return s
}

// Pivot Fluent API Methods - Pivot Tablo Yönetimi
//
// Bu bölüm, many-to-many ilişkilerde pivot tablo alanlarını yönetir.
// İlişki tablosundaki ekstra alanlar için kullanılır.

// AsPivot, alanı pivot field olarak işaretler.
//
// Bu metod, alanın bir pivot tablo alanı olduğunu belirtir.
// Many-to-many ilişkilerde ara tablodaki ekstra alanlar için kullanılır.
//
// # Döndürür
//
//   - core.Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Kullanım Senaryoları
//
//   - Rol-izin ilişkisinde izin tarihi
//   - Ürün-kategori ilişkisinde sıralama
//   - Kullanıcı-grup ilişkisinde katılma tarihi
//
// # Örnek
//
//	field := DateTime("assigned_at", "Atanma Tarihi").
//	    AsPivot().
//	    WithPivotResource("user_roles")
func (s *Schema) AsPivot() core.Element {
	s.IsPivotField = true
	return s
}

// WithPivotResource, pivot resource adını ayarlar.
//
// Bu metod, pivot alanın hangi resource'a ait olduğunu belirtir.
// İlişki yönetimi için kullanılır.
//
// # Parametreler
//
//   - resourceName: Pivot resource adı
//
// # Döndürür
//
//   - core.Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Örnek
//
//	field := Number("quantity", "Miktar").
//	    AsPivot().
//	    WithPivotResource("order_products")
func (s *Schema) WithPivotResource(resourceName string) core.Element {
	s.PivotResourceName = resourceName
	return s
}

// Missing Attachment Methods - Dosya Yükleme Getter Metodları
//
// Bu bölüm, dosya yükleme ayarlarını okumak için kullanılan getter metodlarını içerir.

// GetAcceptedMimeTypes, kabul edilen MIME tiplerini döner.
//
// # Döndürür
//
//   - []string: Kabul edilen MIME tipleri listesi (boş slice olabilir)
//
// # Örnek
//
//	mimeTypes := field.GetAcceptedMimeTypes()
//	for _, mime := range mimeTypes {
//	    fmt.Println("Accepted:", mime)
//	}
func (s *Schema) GetAcceptedMimeTypes() []string {
	if s.AcceptedMimeTypes == nil {
		return []string{}
	}
	return s.AcceptedMimeTypes
}

// GetMaxFileSize, maksimum dosya boyutunu döner.
//
// # Döndürür
//
//   - int64: Maksimum dosya boyutu (byte cinsinden)
//
// # Örnek
//
//	maxSize := field.GetMaxFileSize()
//	fmt.Printf("Max size: %d MB\n", maxSize/(1024*1024))
func (s *Schema) GetMaxFileSize() int64 {
	return s.MaxFileSize
}

// GetStorageDisk, depolama disk adını döner.
//
// # Döndürür
//
//   - string: Depolama disk adı (örn. "local", "s3", "gcs")
//
// # Örnek
//
//	disk := field.GetStorageDisk()
func (s *Schema) GetStorageDisk() string {
	return s.StorageDisk
}

// GetStoragePath, depolama yolunu döner.
//
// # Döndürür
//
//   - string: Depolama yolu (örn. "uploads/avatars")
//
// # Örnek
//
//	path := field.GetStoragePath()
func (s *Schema) GetStoragePath() string {
	return s.StoragePath
}

// ValidateAttachment, dosya yükleme doğrulaması yapar.
//
// Bu metod, dosya adı ve boyutuna göre doğrulama yapar.
// Şu an için placeholder implementasyon içerir.
//
// # Parametreler
//
//   - filename: Dosya adı
//   - size: Dosya boyutu (byte)
//
// # Döndürür
//
//   - error: Doğrulama hatası (nil ise geçerli)
func (s *Schema) ValidateAttachment(filename string, size int64) error {
	return nil
}

// GetUploadCallback, dosya yükleme callback fonksiyonunu döner.
//
// # Döndürür
//
//   - func: Upload callback fonksiyonu (nil olabilir)
//
// # Örnek
//
//	callback := field.GetUploadCallback()
//	if callback != nil {
//	    err := callback(file, item)
//	}
func (s *Schema) GetUploadCallback() func(interface{}, interface{}) error {
	return s.UploadCallback
}

// ShouldRemoveEXIFData, EXIF verilerinin kaldırılıp kaldırılmayacağını döner.
//
// # Döndürür
//
//   - bool: EXIF verileri kaldırılacak mı?
//
// # Örnek
//
//	if field.ShouldRemoveEXIFData() {
//	    // EXIF verilerini kaldır
//	}
func (s *Schema) ShouldRemoveEXIFData() bool {
	return s.RemoveEXIFDataFlag
}

// RemoveEXIFData, dosyadan EXIF verilerini kaldırır.
//
// Bu metod, resim dosyalarındaki EXIF metadata'sını kaldırır.
// Şu an için placeholder implementasyon içerir.
//
// # Parametreler
//
//   - ctx: İşlem bağlamı
//   - file: Dosya nesnesi
//
// # Döndürür
//
//   - error: İşlem hatası (nil ise başarılı)
func (s *Schema) RemoveEXIFData(ctx interface{}, file interface{}) error {
	return nil
}

// Missing Repeater Methods - Tekrarlayan Alan Getter Metodları
//
// Bu bölüm, tekrarlayan alan ayarlarını okumak için kullanılan getter metodlarını içerir.

// IsRepeaterField, alanın bir repeater field olup olmadığını kontrol eder.
//
// # Döndürür
//
//   - bool: Repeater field ise true
//
// # Örnek
//
//	if field.IsRepeaterField() {
//	    fields := field.GetRepeaterFields()
//	}
func (s *Schema) IsRepeaterField() bool {
	return len(s.RepeaterFields) > 0
}

// GetRepeaterFields, repeater alt alanlarını döner.
//
// # Döndürür
//
//   - []Element: Alt alan listesi
//
// # Örnek
//
//	fields := field.GetRepeaterFields()
//	for _, f := range fields {
//	    fmt.Println(f.GetKey())
//	}
func (s *Schema) GetRepeaterFields() []Element {
	return s.RepeaterFields
}

// GetMinRepeats, minimum tekrar sayısını döner.
//
// # Döndürür
//
//   - int: Minimum tekrar sayısı
func (s *Schema) GetMinRepeats() int {
	return s.MinRepeatsCount
}

// GetMaxRepeats, maksimum tekrar sayısını döner.
//
// # Döndürür
//
//   - int: Maksimum tekrar sayısı
func (s *Schema) GetMaxRepeats() int {
	return s.MaxRepeatsCount
}

// ValidateRepeats, tekrar sayısını doğrular.
//
// Bu metod, verilen tekrar sayısının min/max sınırları içinde olup olmadığını kontrol eder.
// Şu an için placeholder implementasyon içerir.
//
// # Parametreler
//
//   - count: Tekrar sayısı
//
// # Döndürür
//
//   - error: Doğrulama hatası (nil ise geçerli)
func (s *Schema) ValidateRepeats(count int) error {
	return nil
}

// Missing Rich Text Methods - Zengin Metin Editörü Getter Metodları
//
// Bu bölüm, zengin metin editörü ayarlarını okumak için kullanılan getter metodlarını içerir.

// GetEditorType, editör tipini döner.
//
// # Döndürür
//
//   - string: Editör tipi (örn. "tinymce", "quill", "ckeditor")
func (s *Schema) GetEditorType() string {
	return s.EditorType
}

// GetEditorLanguage, editör dilini döner.
//
// # Döndürür
//
//   - string: Dil kodu (örn. "tr", "en")
func (s *Schema) GetEditorLanguage() string {
	return s.EditorLanguage
}

// GetEditorTheme, editör temasını döner.
//
// # Döndürür
//
//   - string: Tema adı (örn. "modern", "dark")
func (s *Schema) GetEditorTheme() string {
	return s.EditorTheme
}

// Missing Status Methods - Durum Badge Getter Metodları
//
// Bu bölüm, durum badge ayarlarını okumak için kullanılan getter metodlarını içerir.

// GetStatusColors, durum renk eşleştirmesini döner.
//
// # Döndürür
//
//   - map[string]string: Durum-renk eşleştirme map'i (boş map olabilir)
//
// # Örnek
//
//	colors := field.GetStatusColors()
//	if color, ok := colors["active"]; ok {
//	    fmt.Println("Active color:", color)
//	}
func (s *Schema) GetStatusColors() map[string]string {
	if s.StatusColors == nil {
		return make(map[string]string)
	}
	return s.StatusColors
}

// GetBadgeVariant, badge varyantını döner.
//
// # Döndürür
//
//   - string: Badge varyantı (örn. "solid", "outline")
func (s *Schema) GetBadgeVariant() string {
	return s.BadgeVariant
}

// Missing Pivot Methods - Pivot Tablo Getter Metodları
//
// Bu bölüm, pivot tablo ayarlarını okumak için kullanılan getter metodlarını içerir.

// IsPivot, alanın pivot field olup olmadığını kontrol eder.
//
// # Döndürür
//
//   - bool: Pivot field ise true
//
// # Örnek
//
//	if field.IsPivot() {
//	    resource := field.GetPivotResourceName()
//	}
func (s *Schema) IsPivot() bool {
	return s.IsPivotField
}

// GetPivotResourceName, pivot resource adını döner.
//
// # Döndürür
//
//   - string: Pivot resource adı
func (s *Schema) GetPivotResourceName() string {
	return s.PivotResourceName
}

// Missing Display Methods - Görüntüleme Getter Metodları
//
// Bu bölüm, görüntüleme ayarlarını okumak için kullanılan getter metodlarını içerir.

// GetDisplayCallback, görüntüleme callback fonksiyonunu döner.
//
// # Döndürür
//
//   - func(interface{}) string: Display callback fonksiyonu (nil olabilir)
//
// # Örnek
//
//	callback := field.GetDisplayCallback()
//	if callback != nil {
//	    displayValue := callback(value)
//	}
func (s *Schema) GetDisplayCallback() func(interface{}) string {
	return s.DisplayCallback
}

// GetDisplayedAs, görüntüleme format string'ini döner.
//
// # Döndürür
//
//   - string: Format string'i
func (s *Schema) GetDisplayedAs() string {
	return s.DisplayedAs
}

// ShouldDisplayUsingLabels, etiket kullanarak görüntüleme yapılıp yapılmayacağını döner.
//
// # Döndürür
//
//   - bool: Etiket kullanılacak mı?
//
// # Örnek
//
//	if field.ShouldDisplayUsingLabels() {
//	    // Değer yerine etiketi göster
//	}
func (s *Schema) ShouldDisplayUsingLabels() bool {
	return s.DisplayUsingLabelsFlag
}

// Missing Dependency Methods - Bağımlılık Getter Metodları
//
// Bu bölüm, bağımlılık ayarlarını okumak ve yönetmek için kullanılan metodları içerir.

// SetDependencies, alan bağımlılıklarını ayarlar.
//
// Bu metod, alanın bağımlı olduğu diğer alanları belirler.
//
// # Parametreler
//
//   - deps: Bağımlı alan adları listesi
//
// # Döndürür
//
//   - Element: Zincirleme çağrılar için Schema pointer'ı
//
// # Örnek
//
//	field.SetDependencies([]string{"country", "state"})
func (s *Schema) SetDependencies(deps []string) Element {
	s.DependsOnFields = deps
	return s
}

// GetDependencyRules, bağımlılık kurallarını döner.
//
// # Döndürür
//
//   - map[string]interface{}: Bağımlılık kuralları map'i (boş map olabilir)
//
// # Örnek
//
//	rules := field.GetDependencyRules()
//	for field, rule := range rules {
//	    fmt.Printf("Field %s has rule: %v\n", field, rule)
//	}
func (s *Schema) GetDependencyRules() map[string]interface{} {
	if s.DependencyRules == nil {
		return make(map[string]interface{})
	}
	return s.DependencyRules
}

// ResolveDependencies, bağımlılık kurallarını değerlendirir.
//
// Bu metod, verilen bağlama göre bağımlılık kurallarının karşılanıp karşılanmadığını kontrol eder.
// Şu an için her zaman true döner (placeholder implementasyon).
//
// # Parametreler
//
//   - context: Değerlendirme bağlamı
//
// # Döndürür
//
//   - bool: Bağımlılıklar karşılanıyor mu?
func (s *Schema) ResolveDependencies(context interface{}) bool {
	return true
}

// Missing Suggestion Methods - Öneri Sistemi Getter Metodları
//
// Bu bölüm, öneri sistemi ayarlarını okumak için kullanılan getter metodlarını içerir.

// GetSuggestionsCallback, öneri callback fonksiyonunu döner.
//
// # Döndürür
//
//   - func(string) []interface{}: Suggestions callback fonksiyonu (nil olabilir)
//
// # Örnek
//
//	callback := field.GetSuggestionsCallback()
//	if callback != nil {
//	    suggestions := callback("query")
//	}
func (s *Schema) GetSuggestionsCallback() func(string) []interface{} {
	return s.SuggestionsCallback
}

// GetAutoCompleteURL, otomatik tamamlama URL'ini döner.
//
// # Döndürür
//
//   - string: Autocomplete API endpoint'i
func (s *Schema) GetAutoCompleteURL() string {
	return s.AutoCompleteURL
}

// GetMinCharsForSuggestions, öneri için minimum karakter sayısını döner.
//
// # Döndürür
//
//   - int: Minimum karakter sayısı
func (s *Schema) GetMinCharsForSuggestions() int {
	return s.MinCharsForSuggestionsVal
}

// GetSuggestions, verilen sorgu için önerileri döner.
//
// Bu metod, SuggestionsCallback varsa onu çağırır, yoksa mevcut Suggestions listesini döner.
//
// # Parametreler
//
//   - query: Arama sorgusu
//
// # Döndürür
//
//   - []interface{}: Öneri listesi
//
// # Örnek
//
//	suggestions := field.GetSuggestions("john")
//	for _, s := range suggestions {
//	    fmt.Println(s)
//	}
func (s *Schema) GetSuggestions(query string) []interface{} {
	if s.SuggestionsCallback != nil {
		return s.SuggestionsCallback(query)
	}
	return s.Suggestions
}

// Missing Validation Methods - Doğrulama Getter Metodları
//
// Bu bölüm, doğrulama ayarlarını okumak için kullanılan getter metodlarını içerir.

// GetValidationRules, doğrulama kurallarını döner.
//
// # Döndürür
//
//   - []interface{}: Doğrulama kuralları listesi
//
// # Örnek
//
//	rules := field.GetValidationRules()
//	for _, rule := range rules {
//	    // Kuralı uygula
//	}
func (s *Schema) GetValidationRules() []interface{} {
	rules := make([]interface{}, len(s.ValidationRules))
	for i, r := range s.ValidationRules {
		rules[i] = r
	}
	return rules
}

// ValidateValue, değeri doğrulama kurallarına göre kontrol eder.
//
// Bu metod, verilen değerin tüm doğrulama kurallarını karşılayıp karşılamadığını kontrol eder.
// Şu an için placeholder implementasyon içerir.
//
// # Parametreler
//
//   - value: Doğrulanacak değer
//
// # Döndürür
//
//   - error: Doğrulama hatası (nil ise geçerli)
func (s *Schema) ValidateValue(value interface{}) error {
	return nil
}

// GetCustomValidators, özel doğrulayıcı fonksiyonlarını döner.
//
// # Döndürür
//
//   - []interface{}: Özel doğrulayıcı fonksiyonları listesi
//
// # Örnek
//
//	validators := field.GetCustomValidators()
//	for _, validator := range validators {
//	    // Doğrulayıcıyı çalıştır
//	}
func (s *Schema) GetCustomValidators() []interface{} {
	validators := make([]interface{}, len(s.CustomValidators))
	for i, v := range s.CustomValidators {
		validators[i] = v
	}
	return validators
}

// Missing Extended Field System Methods - Genişletilmiş Alan Sistemi Getter Metodları
//
// Bu bölüm, genişletilmiş alan sistemi özelliklerini okumak için kullanılan getter metodlarını içerir.

// GetResolveHandle, resolve handle değerini döner.
//
// Bu metod, client-side bileşen etkileşimi için kullanılan handle'ı döner.
//
// # Döndürür
//
//   - string: Resolve handle değeri
//
// # Örnek
//
//	handle := field.GetResolveHandle()
//	if handle != "" {
//	    // Handle'ı kullan
//	}
func (s *Schema) GetResolveHandle() string {
	return s.ResolveHandleValue
}

// GORM Yapılandırma Metotları

// Gorm, alan için GORM veritabanı yapılandırmasını ayarlar.
// Bu metod, migration ve model oluşturma için kullanılır.
func (s *Schema) Gorm(config *GormConfig) Element {
	s.GormConfiguration = config
	return s
}

// GetGormConfig, alanın GORM yapılandırmasını döner.
func (s *Schema) GetGormConfig() *GormConfig {
	return s.GormConfiguration
}

// HasGormConfig, alanın GORM yapılandırması olup olmadığını kontrol eder.
func (s *Schema) HasGormConfig() bool {
	return s.GormConfiguration != nil
}

// GormPrimaryKey, alanı birincil anahtar olarak işaretler.
func (s *Schema) GormPrimaryKey() Element {
	if s.GormConfiguration == nil {
		s.GormConfiguration = NewGormConfig()
	}
	s.GormConfiguration.WithPrimaryKey()
	return s
}

// GormIndex, alan için indeks oluşturur.
func (s *Schema) GormIndex(name ...string) Element {
	if s.GormConfiguration == nil {
		s.GormConfiguration = NewGormConfig()
	}
	s.GormConfiguration.WithIndex(name...)
	return s
}

// GormUniqueIndex, alan için benzersiz indeks oluşturur.
func (s *Schema) GormUniqueIndex(name ...string) Element {
	if s.GormConfiguration == nil {
		s.GormConfiguration = NewGormConfig()
	}
	s.GormConfiguration.WithUniqueIndex(name...)
	return s
}

// GormType, alan için SQL tipini belirler.
func (s *Schema) GormType(sqlType string) Element {
	if s.GormConfiguration == nil {
		s.GormConfiguration = NewGormConfig()
	}
	s.GormConfiguration.WithType(sqlType)
	return s
}

// GormSize, alan için sütun boyutunu belirler.
func (s *Schema) GormSize(size int) Element {
	if s.GormConfiguration == nil {
		s.GormConfiguration = NewGormConfig()
	}
	s.GormConfiguration.WithSize(size)
	return s
}

// GormDefault, alan için varsayılan değer belirler.
func (s *Schema) GormDefault(value interface{}) Element {
	if s.GormConfiguration == nil {
		s.GormConfiguration = NewGormConfig()
	}
	s.GormConfiguration.WithDefault(value)
	return s
}

// GormComment, alan için veritabanı yorumu ekler.
func (s *Schema) GormComment(comment string) Element {
	if s.GormConfiguration == nil {
		s.GormConfiguration = NewGormConfig()
	}
	s.GormConfiguration.WithComment(comment)
	return s
}

// GormFullTextIndex, alan için fulltext indeks oluşturur.
func (s *Schema) GormFullTextIndex(name ...string) Element {
	if s.GormConfiguration == nil {
		s.GormConfiguration = NewGormConfig()
	}
	s.GormConfiguration.WithFullTextIndex(name...)
	return s
}

// GormSoftDelete, alan için soft delete desteği ekler.
func (s *Schema) GormSoftDelete() Element {
	if s.GormConfiguration == nil {
		s.GormConfiguration = NewGormConfig()
	}
	s.GormConfiguration.WithSoftDelete()
	return s
}

// GormAutoIncrement, alanı otomatik artış olarak işaretler.
func (s *Schema) GormAutoIncrement() Element {
	if s.GormConfiguration == nil {
		s.GormConfiguration = NewGormConfig()
	}
	s.GormConfiguration.WithAutoIncrement()
	return s
}

// GormForeignKey, alan için foreign key ilişkisi tanımlar.
func (s *Schema) GormForeignKey(fk, references string) Element {
	if s.GormConfiguration == nil {
		s.GormConfiguration = NewGormConfig()
	}
	s.GormConfiguration.WithForeignKey(fk, references)
	return s
}

// GormOnDelete, foreign key için ON DELETE davranışını belirler.
func (s *Schema) GormOnDelete(action string) Element {
	if s.GormConfiguration == nil {
		s.GormConfiguration = NewGormConfig()
	}
	s.GormConfiguration.WithOnDelete(action)
	return s
}

// GormOnUpdate, foreign key için ON UPDATE davranışını belirler.
func (s *Schema) GormOnUpdate(action string) Element {
	if s.GormConfiguration == nil {
		s.GormConfiguration = NewGormConfig()
	}
	s.GormConfiguration.WithOnUpdate(action)
	return s
}

// GormUUID, UUID tipinde ID kullanır.
func (s *Schema) GormUUID() Element {
	if s.GormConfiguration == nil {
		s.GormConfiguration = NewGormConfig()
	}
	s.GormConfiguration.WithUUID()
	return s
}

// GormSnowflake, Snowflake tipinde ID kullanır.
func (s *Schema) GormSnowflake() Element {
	if s.GormConfiguration == nil {
		s.GormConfiguration = NewGormConfig()
	}
	s.GormConfiguration.WithSnowflake()
	return s
}

// GormULID, ULID tipinde ID kullanır.
func (s *Schema) GormULID() Element {
	if s.GormConfiguration == nil {
		s.GormConfiguration = NewGormConfig()
	}
	s.GormConfiguration.WithULID()
	return s
}

// Rows, textarea field'ı için satır sayısını ayarlar.
//
// Bu metod, textarea field'larında görüntülenecek satır sayısını belirler.
// Satır sayısı bilgisi Props'a "rows" key'i ile kaydedilir.
//
// Parametreler:
//   - rows: Satır sayısı
//
// Döndürür:
//   - Element pointer'ı (method chaining için)
//
// Örnek:
//
//	field := fields.Textarea("Description", "description").Rows(5)
func (s *Schema) Rows(rows int) Element {
	return s.WithProps("rows", rows)
}
