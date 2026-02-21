# Alanlar (Fields) Rehberi (Legacy Teknik Akış)

Alanlar, veritabanı sütunlarını yönetim panelinde nasıl göstereceğinizi ve işleneceğini tanımlar.

Bu doküman, `fields.Schema` tabanlı düşük seviye/legacy akışı referans alır ve `Resource + FieldResolver` yaklaşımıyla uyumludur.

## Bu Doküman Ne Zaman Okunmalı?

Önerilen sıra:
1. [Başlarken](Getting-Started)
2. [Kaynaklar (Resource)](Resources)
3. Bu doküman (`Fields`)
4. [İlişkiler (Relationships)](Relationships)

## Hızlı Field Karar Akışı

- Basit metin/veri: `text`, `textarea`, `number`, `switch`
- Sabit seçenek: `select`, `radio`, `checkbox`
- Tarih/saat: `date`, `time`, `datetime` (dialog/native kararına göre)
- İlişki verisi: relationship field'ları (bu dokümanın ilişki bölümüne bakın)
- Dosya/medya: `image`, `video`, `audio`, `file`
- Gelişmiş içerik: `richtext`, `keyvalue`, `combobox`

## Alan Türleri

### Metin Alanı (Text)

Tek satırlık metin girişi için kullanılır.

```go
(&fields.Schema{
	Key:   "email",
	Name:  "E-posta",
	View:  "text",
	Props: make(map[string]interface{}),
}).OnList().OnDetail().OnForm()
```

**Seçenekler:**
- `Placeholder("Örnek@email.com")` - Yer tutucu metni
- `Required()` - Zorunlu alan
- `ReadOnly()` - Salt okunur

### Metin Alanı (Textarea)

Çok satırlık metin girişi için kullanılır.

```go
(&fields.Schema{
	Key:   "description",
	Name:  "Açıklama",
	View:  "textarea",
	Props: make(map[string]interface{}),
}).OnForm()
```

### Seçim Alanı (Select)

Önceden tanımlanmış seçeneklerden birini seçmek için kullanılır.

```go
(&fields.Schema{
	Key:   "status",
	Name:  "Durum",
	View:  "select",
	Props: map[string]interface{}{
		"options": map[string]string{
			"draft":     "Taslak",
			"published": "Yayınlandı",
			"archived":  "Arşivlendi",
		},
	},
}).OnList().OnDetail().OnForm()
```

### Tarih Alanı (Date)

Tarih seçimi için kullanılır. İki mod destekler:
- **Dialog Modu (Varsayılan)**: Popover içinde takvim arayüzü, tarih seçildiğinde otomatik kapanma
- **Native Modu**: HTML5 date input, mobil uyumlu

```go
// Dialog modu (varsayılan)
(&fields.Schema{
	Key:   "published_at",
	Name:  "Yayınlanma Tarihi",
	View:  "date",
	Props: make(map[string]interface{}),
}).OnDetail().OnForm()

// Native modu
(&fields.Schema{
	Key:   "birth_date",
	Name:  "Doğum Tarihi",
	View:  "date",
	Props: map[string]interface{}{
		"native": true, // Native HTML date input kullan
	},
}).OnForm().Required()
```

**Frontend Kullanımı:**
```tsx
// Dialog modu (varsayılan)
<DateField
  name="published_at"
  label="Yayınlanma Tarihi"
  value={publishedAt}
  onChange={setPublishedAt}
/>

// Native modu
<DateField
  name="birth_date"
  label="Doğum Tarihi"
  value={birthDate}
  onChange={setBirthDate}
  useNative
  required
/>
```

### Tarih-Saat Alanı (DateTime)

Tarih ve saat seçimi için kullanılır. İki mod destekler:
- **Dialog Modu (Varsayılan)**: Popover içinde takvim + saat girişi, "Tamam" butonu ile kapanma
- **Native Modu**: HTML5 datetime-local input, mobil uyumlu

```go
// Dialog modu (varsayılan)
(&fields.Schema{
	Key:   "created_at",
	Name:  "Oluşturulma Tarihi",
	View:  "datetime",
	Props: make(map[string]interface{}),
}).ReadOnly().OnList().OnDetail()

// Native modu
(&fields.Schema{
	Key:   "appointment_at",
	Name:  "Randevu Tarihi ve Saati",
	View:  "datetime",
	Props: map[string]interface{}{
		"native": true, // Native HTML datetime-local input kullan
	},
}).OnForm().Required()
```

**Frontend Kullanımı:**
```tsx
// Dialog modu (varsayılan)
<DateTimeField
  name="created_at"
  label="Oluşturulma Tarihi"
  value={createdAt}
  onChange={setCreatedAt}
/>

// Native modu
<DateTimeField
  name="appointment_at"
  label="Randevu Tarihi ve Saati"
  value={appointmentAt}
  onChange={setAppointmentAt}
  useNative
  required
/>
```

### Saat Alanı (Time)

Saat seçimi için kullanılır. İki mod destekler:
- **Dialog Modu (Varsayılan)**: Popover içinde saat girişi, "Tamam" butonu ile kapanma
- **Native Modu**: HTML5 time input, mobil uyumlu

```go
// Dialog modu (varsayılan)
(&fields.Schema{
	Key:   "start_time",
	Name:  "Başlangıç Saati",
	View:  "time",
	Props: make(map[string]interface{}),
}).OnForm()

// Native modu
(&fields.Schema{
	Key:   "work_hours",
	Name:  "Çalışma Saati",
	View:  "time",
	Props: map[string]interface{}{
		"native": true, // Native HTML time input kullan
	},
}).OnForm().Required()
```

**Frontend Kullanımı:**
```tsx
// Dialog modu (varsayılan)
<TimeField
  name="start_time"
  label="Başlangıç Saati"
  value={startTime}
  onChange={setStartTime}
/>

// Native modu
<TimeField
  name="work_hours"
  label="Çalışma Saati"
  value={workHours}
  onChange={setWorkHours}
  useNative
  required
/>
```

**Mod Seçimi Rehberi:**

| Özellik | Dialog Modu | Native Modu |
|---------|-------------|-------------|
| **Görsel** | Zengin takvim/saat arayüzü | Basit HTML input |
| **Mobil Uyumluluk** | İyi | Mükemmel (native picker) |
| **Bundle Boyutu** | Daha büyük (~20KB) | Minimal (~1KB) |
| **Özelleştirme** | Yüksek | Sınırlı |
| **Performans** | Orta | Hızlı |
| **Kullanım Senaryosu** | Desktop uygulamalar | Mobil uygulamalar, basit formlar |

### Evet/Hayır Alanı (Switch)

Boolean değerleri için kullanılır.

```go
(&fields.Schema{
	Key:   "is_active",
	Name:  "Aktif mi?",
	View:  "switch",
	Props: make(map[string]interface{}),
}).OnList().OnDetail().OnForm()
```

### Checkbox Alanı

Tek bir checkbox veya checkbox grubu için kullanılır.

**Tek Checkbox:**
```go
(&fields.Schema{
	Key:   "terms_accepted",
	Name:  "Kullanım koşullarını kabul ediyorum",
	View:  "checkbox",
	Props: make(map[string]interface{}),
}).OnForm().Required()
```

**Checkbox Grubu:**
```go
(&fields.Schema{
	Key:   "interests",
	Name:  "İlgi Alanları",
	View:  "checkbox",
	Props: map[string]interface{}{
		"options": []map[string]interface{}{
			{"value": "sports", "label": "Spor"},
			{"value": "music", "label": "Müzik"},
			{"value": "tech", "label": "Teknoloji"},
		},
	},
}).OnForm()
```

**Frontend Kullanımı:**
```tsx
// Tek checkbox
<CheckboxField
  name="terms"
  label="Kullanım koşullarını kabul ediyorum"
  checked={terms}
  onCheckedChange={setTerms}
  required
/>

// Checkbox grubu
<CheckboxField
  name="interests"
  label="İlgi Alanları"
  options={[
    { value: 'sports', label: 'Spor' },
    { value: 'music', label: 'Müzik' },
    { value: 'tech', label: 'Teknoloji' }
  ]}
  value={interests}
  onChange={setInterests}
/>
```

### Radio Group Alanı

Birden fazla seçenek arasından tek bir seçim için kullanılır.

```go
(&fields.Schema{
	Key:   "gender",
	Name:  "Cinsiyet",
	View:  "radio",
	Props: map[string]interface{}{
		"options": []map[string]interface{}{
			{"value": "male", "label": "Erkek"},
			{"value": "female", "label": "Kadın"},
			{"value": "other", "label": "Diğer"},
		},
	},
}).OnForm().Required()

// Açıklamalı seçenekler
(&fields.Schema{
	Key:   "plan",
	Name:  "Plan Seçimi",
	View:  "radio",
	Props: map[string]interface{}{
		"options": []map[string]interface{}{
			{
				"value":       "basic",
				"label":       "Temel",
				"description": "Temel özellikler",
			},
			{
				"value":       "pro",
				"label":       "Pro",
				"description": "Gelişmiş özellikler",
			},
		},
		"orientation": "horizontal", // veya "vertical" (varsayılan)
	},
}).OnForm().Required()
```

**Frontend Kullanımı:**
```tsx
// Dikey düzen (varsayılan)
<RadioGroupField
  name="gender"
  label="Cinsiyet"
  options={[
    { value: 'male', label: 'Erkek' },
    { value: 'female', label: 'Kadın' },
    { value: 'other', label: 'Diğer' }
  ]}
  value={gender}
  onChange={setGender}
  required
/>

// Yatay düzen
<RadioGroupField
  name="status"
  label="Durum"
  options={[
    { value: 'active', label: 'Aktif' },
    { value: 'inactive', label: 'Pasif' }
  ]}
  value={status}
  onChange={setStatus}
  orientation="horizontal"
/>

// Açıklamalı seçenekler
<RadioGroupField
  name="plan"
  label="Plan Seçimi"
  options={[
    { value: 'basic', label: 'Temel', description: 'Temel özellikler' },
    { value: 'pro', label: 'Pro', description: 'Gelişmiş özellikler' }
  ]}
  value={plan}
  onChange={setPlan}
/>
```

### Sayı Alanı (Number)

Sayısal değerler için kullanılır.

```go
(&fields.Schema{
	Key:   "price",
	Name:  "Fiyat",
	View:  "number",
	Props: make(map[string]interface{}),
}).OnList().OnDetail().OnForm()
```

### E-posta Alanı (Email)

E-posta doğrulaması ile metin alanı.

```go
(&fields.Schema{
	Key:   "email",
	Name:  "E-posta",
	View:  "email",
	Props: make(map[string]interface{}),
}).OnForm().Required()
```

### URL Alanı (URL)

URL doğrulaması ile metin alanı.

```go
(&fields.Schema{
	Key:   "website",
	Name:  "Web Sitesi",
	View:  "url",
	Props: make(map[string]interface{}),
}).OnForm()
```

### Şifre Alanı (Password)

Şifre girişi için maskelenmiş alan.

```go
fields.Password("Şifre").
	OnForm().
	Required().
	MinLength(8).
	HelpText("En az 8 karakter olmalıdır")
```

### Telefon Alanı (Tel)

Telefon numarası girişi için kullanılır.

```go
fields.Tel("Telefon", "phone").
	OnList().
	OnDetail().
	OnForm().
	Placeholder("+90 (555) 123-4567").
	Pattern(`^\+?[1-9]\d{1,14}$`)
```

### Görsel Alanı (Image)

Görsel dosyası yükleme için kullanılır.

```go
fields.Image("Profil Fotoğrafı", "avatar").
	OnDetail().
	OnForm().
	Accept("image/jpeg", "image/png", "image/webp").
	MaxSize(5 * 1024 * 1024). // 5MB
	Store("public", "avatars").
	RemoveEXIFData()
```

**Özellikler:**
- `Accept()` - Kabul edilen dosya tipleri
- `MaxSize()` - Maksimum dosya boyutu (byte)
- `Store()` - Depolama diski ve yolu
- `RemoveEXIFData()` - EXIF verilerini kaldır

### Video Alanı (Video)

Video dosyası yükleme için kullanılır.

```go
fields.Video("Tanıtım Videosu", "promo_video").
	OnDetail().
	OnForm().
	Accept("video/mp4", "video/webm").
	MaxSize(100 * 1024 * 1024). // 100MB
	Store("public", "videos").
	HelpText("MP4 veya WebM formatında, maksimum 100MB")
```

### Ses Alanı (Audio)

Ses dosyası yükleme için kullanılır.

```go
fields.Audio("Podcast", "audio_file").
	OnDetail().
	OnForm().
	Accept("audio/mpeg", "audio/wav", "audio/ogg").
	MaxSize(50 * 1024 * 1024). // 50MB
	Store("public", "audio")
```

### Dosya Alanı (File)

Genel dosya yükleme için kullanılır.

```go
fields.File("Döküman", "document").
	OnDetail().
	OnForm().
	Accept("application/pdf", "application/msword").
	MaxSize(10 * 1024 * 1024). // 10MB
	Store("private", "documents").
	Required()
```

### Zengin Metin Editörü (RichText)

WYSIWYG editör ile zengin metin girişi için kullanılır.

```go
fields.RichText("İçerik", "content").
	OnForm().
	OnDetail().
	WithEditor("tiptap"). // veya "quill", "tinymce"
	WithLanguage("tr").
	WithTheme("snow").
	Required()
```

**Özellikler:**
- `WithEditor()` - Editör tipi (tiptap, quill, tinymce)
- `WithLanguage()` - Editör dili
- `WithTheme()` - Editör teması

### Anahtar-Değer Alanı (KeyValue)

Dinamik anahtar-değer çiftleri için kullanılır.

```go
fields.KeyValue("Meta Veriler", "metadata").
	OnForm().
	OnDetail().
	HelpText("Özel meta veriler ekleyin")
```

**Kullanım Örneği:**
```json
{
  "og:title": "Sayfa Başlığı",
  "og:description": "Sayfa Açıklaması",
  "og:image": "https://example.com/image.jpg"
}
```

### Matrix Alanı (Matrix)

Satır-sütun bazlı dinamik veri girişi için kullanılır. Özellikle varyant/opsiyon gibi
çoklu satır senaryolarında uygundur.

```go
fields.Matrix("Variant Matrix", "variant_matrix").
	OnForm().
	WithProps("options", map[string]interface{}{
		"columns": []map[string]interface{}{
			{
				"key":   "product_option_id",
				"label": "Product Option",
				"type":  "select",
				"options": map[string]interface{}{
					"1": "Color",
					"2": "Size",
				},
			},
			{
				"key":       "product_option_value_id",
				"label":     "Product Option Value",
				"type":      "select",
				"dependsOn": "product_option_id",
				"optionsByDependency": map[string]interface{}{
					"1": map[string]interface{}{"10": "Red", "11": "Blue"},
					"2": map[string]interface{}{"20": "S", "21": "M"},
				},
			},
			{
				"key":   "is_variant",
				"label": "Variant Option Value",
				"type":  "radio",
			},
		},
		"allowAddingRows":   true,
		"allowDeletingRows": true,
		"addButtonText":     "Satır Ekle",
		"emptyMessage":      "Henüz satır yok",
	}),
	WithProps("keys", map[string]interface{}{
		"option":       "product_option_id",
		"option_value": "product_option_value_id",
		"variant":      "is_variant",
	})
```

**Desteklenen hücre tipleri:**
- `text` (varsayılan)
- `number`
- `textarea`
- `select`
- `checkbox`
- `radio`

**Önemli Notlar:**
- `columns[].type` verilmezse varsayılan olarak `text` kullanılır.
- `columns[].key` backend tarafındaki veri anahtarıdır; payload bu anahtarla gelir.
- Satır bazlı `+` ve silme aksiyonları ile alt kısımda ekle butonu desteklenir.
- `dependsOn + optionsByDependency` ile bağımlı select akışı kurulabilir.
- `radio` için `options` verilmezse satır-seçici gibi çalışır (tek satır true kalır).

**Gönderilen payload örneği:**
```json
{
  "variant_matrix": {
    "rows": [
      {
        "product_option_id": "1",
        "product_option_value_id": "10",
        "is_variant": true
      }
    ],
    "keys": {
      "option": "product_option_id",
      "option_value": "product_option_value_id",
      "variant": "is_variant"
    }
  }
}
```

### Combobox Alanı

Arama yapılabilir seçim listesi için kullanılır.

```go
fields.Combobox("Kategori", "category_id").
	OnList().
	OnDetail().
	OnForm().
	Options(map[string]string{
		"1": "Elektronik",
		"2": "Giyim",
		"3": "Kitap",
	}).
	Searchable().
	Required()
```

**AutoOptions ile:**
```go
fields.Combobox("Kategori", "category_id").
	OnForm().
	AutoOptions("name"). // Otomatik olarak ilişkili kayıtlardan seçenekler oluşturur
	Required()
```

## İlişki Alanları (Relationship Fields)

İlişki alanları, veritabanı tablolarınız arasındaki ilişkileri yönetmek için kullanılır.

### Link (BelongsTo İlişkisi)

Bir kaydın başka bir kayda ait olduğunu gösterir (N:1 ilişki).

```go
fields.Link("Yazar", "users", "author_id").
	OnList().
	OnDetail().
	OnForm().
	Required().
	HelpText("Bu yazının yazarını seçin")
```

**Örnek Kullanım:**
```go
// Post -> User ilişkisi
fields.Link("Yazar", "users", "author_id").
	OnList().
	OnDetail().
	OnForm()

// Order -> Customer ilişkisi
fields.Link("Müşteri", "customers", "customer_id").
	OnList().
	OnDetail().
	OnForm().
	Required()
```

### Detail (HasOne İlişkisi)

Bir kaydın tek bir ilişkili kaydı olduğunu gösterir (1:1 ilişki).

```go
fields.Detail("Profil", "profiles", "profile").
	OnDetail().
	HelpText("Kullanıcının profil bilgileri")
```

**Örnek Kullanım:**
```go
// User -> Profile ilişkisi
fields.Detail("Profil", "profiles", "profile").
	OnDetail()

// Order -> Invoice ilişkisi
fields.Detail("Fatura", "invoices", "invoice").
	OnDetail()
```

### Collection (HasMany İlişkisi)

Bir kaydın birden fazla ilişkili kaydı olduğunu gösterir (1:N ilişki).

```go
fields.Collection("Yorumlar", "comments", "comments").
	OnDetail().
	HelpText("Bu yazıya yapılan yorumlar")
```

**Örnek Kullanım:**
```go
// Post -> Comments ilişkisi
fields.Collection("Yorumlar", "comments", "comments").
	OnDetail()

// User -> Orders ilişkisi
fields.Collection("Siparişler", "orders", "orders").
	OnDetail()
```

### Connect (BelongsToMany İlişkisi)

Çoktan çoğa ilişki için kullanılır (N:M ilişki).

```go
fields.Connect("Etiketler", "tags", "tags").
	OnForm().
	OnDetail().
	HelpText("Bu yazıya etiket ekleyin")
```

**Örnek Kullanım:**
```go
// Post -> Tags ilişkisi
fields.Connect("Etiketler", "tags", "tags").
	OnForm().
	OnDetail()

// User -> Roles ilişkisi
fields.Connect("Roller", "roles", "roles").
	OnForm().
	OnDetail().
	Required()
```

### PolyLink (MorphTo İlişkisi)

Polimorfik ilişki için kullanılır (bir kayıt farklı tiplerdeki kayıtlara ait olabilir).

```go
fields.PolyLink("İlişkili Kayıt", "commentable").
	OnDetail().
	HelpText("Bu yorumun ait olduğu kayıt")
```

**Örnek Kullanım:**
```go
// Comment -> Post veya Video ilişkisi
fields.PolyLink("Yorumlanabilir", "commentable").
	OnDetail()

// Image -> Post, Product veya User ilişkisi
fields.PolyLink("Görselin Sahibi", "imageable").
	OnDetail()
```

### PolyDetail (MorphOne İlişkisi)

Polimorfik tekil ilişki için kullanılır.

```go
fields.PolyDetail("Görsel", "images", "image").
	OnDetail().
	HelpText("Bu kayda ait görsel")
```

### PolyCollection (MorphMany İlişkisi)

Polimorfik çoğul ilişki için kullanılır.

```go
fields.PolyCollection("Yorumlar", "comments", "comments").
	OnDetail().
	HelpText("Bu kayda yapılan yorumlar")
```

**Örnek Kullanım:**
```go
// Post -> Comments (polimorfik)
fields.PolyCollection("Yorumlar", "comments", "comments").
	OnDetail()

// Product -> Images (polimorfik)
fields.PolyCollection("Görseller", "images", "images").
	OnDetail()
```

### PolyConnect (MorphToMany İlişkisi)

Polimorfik çoktan çoğa ilişki için kullanılır.

```go
fields.PolyConnect("Etiketler", "tags", "tags").
	OnForm().
	OnDetail().
	HelpText("Bu kayda etiket ekleyin")
```

## Alan Seçenekleri

### Görünürlük Kontrolleri

```go
// Sadece listede göster
field.OnlyOnList()

// Sadece detayda göster
field.OnlyOnDetail()

// Sadece formda göster
field.OnlyOnForm()

// Listede gizle
field.HideOnList()

// Grid kart/listing görünümünde gizle
field.HideOnGrid()

// Grid'de zorunlu göster (HideOnList varsa grid'de override eder)
field.ShowOnGrid()

// Index kapsamlarında (table + grid) göster, form/detail'da gizle
field.ShowOnlyGrid()

// Detayda gizle
field.HideOnDetail()

// Oluşturmada gizle
field.HideOnCreate()

// Güncellemede gizle
field.HideOnUpdate()

// External API yanıtında gizle
field.HideOnApi()

// Tüm yerlerde göster
field.OnList().OnDetail().OnForm()
```

`HideOnApi()` sadece external API (`/api`) çıktısını etkiler; internal panel endpoint'leri (`/api/internal/*`) ve internal REST API (`/api/internal/rest/*`) davranışını değiştirmez.

#### Grid görünürlük kuralları (özet)

- `HideOnGrid`: Grid kart/listing görünümünde gizler; table/detail/form etkilenmez.
- `ShowOnGrid`: Grid'de görünür olmasını zorlar. Özellikle `HideOnList().ShowOnGrid()` kombinasyonunda table'da gizli, grid'de görünür olur.
- `ShowOnlyGrid`: Kaydı index kapsamlarında (table + grid) görünür tutar, create/update/detail'da gizler.

Öncelik:
- `HideOnGrid`, grid görünümünde baskındır.
- `HideOnList`, grid'de de gizler; `ShowOnGrid` ile override edilebilir.
- `OnlyOnDetail/OnlyOnForm/OnlyOnCreate/OnlyOnUpdate` kısıtları grid'de de korunur.

Not:
- `HideOnGrid`, kart/listing render'ını etkiler.
- Alan değeri row payload'da korunur.

### Validasyon

```go
// Zorunlu alan
field.Required()

// Boş bırakılabilir
field.Nullable()

// E-posta doğrulaması
field.Email()

// URL doğrulaması
field.URL()

// Minimum uzunluk
field.MinLength(5)

// Maksimum uzunluk
field.MaxLength(100)

// Minimum değer
field.Min(0)

// Maksimum değer
field.Max(1000)

// Regex deseni
field.Pattern("^[A-Z].*")

// Benzersiz değer
field.Unique("users", "email")

// Var olan değer
field.Exists("categories", "id")
```

### Görünüm Seçenekleri

```go
// Salt okunur
field.ReadOnly()

// Devre dışı
field.Disabled()

// Değiştirilemez
field.Immutable()

// Varsayılan değer
field.Default("Varsayılan")

// Yer tutucu
field.Placeholder("Buraya yazın...")

// Yardım metni
field.HelpText("Bu alan hakkında bilgi")

// Etiket
field.Label("Özel Etiket")
```

### Arama ve Sıralama

```go
// Aranabilir
field.Searchable()

// Sıralanabilir
field.Sortable()

// Filtrelenebilir
field.Filterable()

// Yığılı gösterim
field.Stacked()

// Grid kolon genişliği (Form + Detail)
// 1-12 arası desteklenir
field.Span(6)
```

### Form/Detail Grid Yerleşimi (`Span`)

Alanları form ve detail görünümünde 12 kolonlu grid üzerinde konumlandırabilirsiniz.

```go
fields.Text("Ad", "first_name").Span(6)
fields.Text("Soyad", "last_name").Span(6)
fields.Email("E-posta", "email").Span(12)
```

**Notlar:**
- `Span(1)` ile `Span(12)` arası desteklenir.
- `Span` verilmezse varsayılan değer `12`'dir (tam genişlik).
- Geçersiz değerler otomatik düzeltilir: `<1 => 1`, `>12 => 12`.
- Bu özellik **form** ve **detail** görünümü için geçerlidir.
- **Index/list** görünümünde `Span` kullanılmaz.
- `has-many`, `belongs-to-many`, `morph-to-many` gibi ilişki tablo alanları detail'da tam genişlikte kalır.

## Gelişmiş Özellikler

### Display Callbacks ve Formatlama

Alan değerlerinin nasıl görüntüleneceğini özelleştirin.

```go
// Display callback (önerilen imza: value + item)
fields.Text("Fiyat", "price").
	OnList().
	Display(func(value interface{}, item interface{}) interface{} {
		if price, ok := value.(float64); ok {
			return fmt.Sprintf("₺%.2f", price)
		}
		return ""
	})

// Computed alan örneği
// "sizes" veritabanında yoksa value nil gelir, model item'dan hesaplanır
fields.Text("Bedenler", "sizes").
	OnList().
	OnDetail().
	Display(func(value interface{}, item interface{}) interface{} {
		product, ok := item.(*product.Product)
		if !ok {
			return value
		}

		if value == nil {
			return strings.Join(product.AvailableSizes, ", ")
		}
		return value
	})

// Component döndürme (Badge)
fields.Text("Bedenler", "sizes").
	OnList().
	Display(func(value interface{}, item interface{}) interface{} {
		product, ok := item.(*product.Product)
		if !ok {
			return value
		}

		return fields.Badge("Bedenler").
			Default(strings.Join(product.AvailableSizes, ", ")).
			WithProps("variant", "secondary")
	})

// Component döndürme (Stack + Badge)
fields.Text("Sizes", "sizes").
	Display(func(value interface{}, item interface{}) core.Element {
		pv, ok := item.(*entity.Productvariant)
		if !ok || pv == nil {
			return fields.Stack([]core.Element{})
		}

		return fields.Stack([]core.Element{
			fields.Badge(fmt.Sprintf("%d", pv.Weight)).WithProps("variant", "secondary"),
			fields.Badge(fmt.Sprintf("%d", pv.Height)).WithProps("variant", "secondary"),
			fields.Badge(fmt.Sprintf("%d", pv.Width)).WithProps("variant", "secondary"),
			fields.Badge(fmt.Sprintf("%d", pv.Volume)).WithProps("variant", "secondary"),
		})
	})

// DisplayAs ile format string
fields.Text("Durum", "status").
	OnList().
	DisplayAs("Durum: %s")

// DisplayUsingLabels ile etiket kullanımı
fields.Select("Kategori", "category").
	OnList().
	Options(map[string]string{
		"1": "Elektronik",
		"2": "Giyim",
	}).
	DisplayUsingLabels() // "1" yerine "Elektronik" gösterir
```

`Display` callback içinde `fields.Badge(...)` gibi `core.Element` döndürürseniz
backend bunu otomatik algılar ve frontend tarafı render eder.

`Display` callback içinde `fields.Stack(...)` döndürürseniz, içindeki `fields` listesi
de recursive serialize edilir ve çocuk bileşenler birlikte render edilir.

### Grid kartı + Stack kullanımı (önerilen akış)

Grid kartında içerik sırası backend/frontend tarafından otomatik uygulanır:

1. Varsa ilk görünür `image-field` kart başında gösterilir
2. Altına resource `record_title_key` değeri başlık olarak yazılır
3. Kalan görünür alanlar field sırasına göre kart gövdesinde gösterilir
4. `Display(...)->fields.Stack(...)` dönen computed içerik doğrudan kart gövdesinde render edilir

```go
func NewProductResource() resource.Resource {
	r := resource.New("products")

	// Resource seviyesinde grid görünümünü aç/kapat
	r.SetGridEnabled(true)

	// Kart başlığında kullanılacak alan
	r.SetRecordTitleKey("name")

	r.SetFields([]fields.Element{
		fields.ID(),

		// Grid kartı üst görseli
		fields.Image("Görsel", "image").
			OnList().
			OnDetail(),

		fields.Text("Ad", "name").
			OnList().
			OnDetail().
			OnForm(),

		// Table'da gizli, grid kartında görünür computed blok
		fields.Text("Özet", "summary").
			HideOnList().
			ShowOnGrid().
			Display(func(value interface{}, item interface{}) core.Element {
				p, ok := item.(*Product)
				if !ok || p == nil {
					return fields.Stack([]core.Element{})
				}

				return fields.Stack([]core.Element{
					fields.Badge(fmt.Sprintf("Stok: %d", p.Stock)).WithProps("variant", "secondary"),
					fields.Badge(fmt.Sprintf("Fiyat: ₺%.2f", p.Price)).WithProps("variant", "outline"),
				})
			}),
	})

	return r
}
```

Notlar:
- Bu akışta ek bir `Line` API yoktur; computed kart içeriği için mevcut `Display + Stack` kullanılır.
- Boş değerler kartta `—` olarak gösterilir.

`Display` için desteklenen callback imzaları:

- `func(value interface{}) string`
- `func(value interface{}) core.Element`
- `func(value interface{}) interface{}`
- `func(value interface{}, item interface{}) string`
- `func(value interface{}, item interface{}) core.Element`
- `func(value interface{}, item interface{}) interface{}`

Notlar:

- Sayısal değer formatlarken `int64` için `%d`, `float64` için `%f` kullanın.
- `fields.Badge("10")` kullanımında değer `name` alanından da okunur (`data` fallback).

### Bağımlılıklar (Dependencies)

Alanların birbirine bağımlı olmasını sağlayın.

```go
// Basit bağımlılık
fields.Text("Şehir", "city").
	OnForm().
	DependsOn("country")

// Koşullu bağımlılık
fields.Text("Vergi Numarası", "tax_number").
	OnForm().
	DependsOn("is_company").
	When("is_company", "=", true)

// Çoklu bağımlılık
fields.Select("İlçe", "district").
	OnForm().
	DependsOn("country", "city").
	When("country", "=", "TR")
```

### Öneriler ve Autocomplete

Kullanıcıya dinamik öneriler sunun.

```go
// Statik öneriler
fields.Text("E-posta", "email").
	OnForm().
	WithSuggestions(func(query string) []interface{} {
		return []interface{}{
			query + "@gmail.com",
			query + "@hotmail.com",
			query + "@outlook.com",
		}
	}).
	MinCharsForSuggestions(3)

// API'den autocomplete
fields.Text("Şehir", "city").
	OnForm().
	WithAutoComplete("/api/cities/search").
	MinCharsForSuggestions(2)
```

### Dosya Yükleme Yapılandırması

Dosya yükleme alanlarını detaylı yapılandırın.

```go
// Görsel yükleme - tam yapılandırma
fields.Image("Profil Fotoğrafı", "avatar").
	OnForm().
	OnDetail().
	Accept("image/jpeg", "image/png", "image/webp").
	MaxSize(5 * 1024 * 1024). // 5MB
	Store("public", "avatars").
	RemoveEXIFData().
	WithUpload(func(file interface{}, item interface{}) error {
		// Özel yükleme işlemi
		return nil
	})

// Çoklu dosya yükleme
fields.File("Dökümanlar", "documents").
	OnForm().
	Accept("application/pdf", "application/msword").
	MaxSize(10 * 1024 * 1024).
	Store("private", "documents").
	HelpText("PDF veya Word dosyaları yükleyebilirsiniz")
```

**Image Preview Notu**
- `fields.Image(...)` alanları formda mevcut resmi otomatik önizler.
- Kayıt değeri tam URL (`https://...`) veya root-relative (`/storage/...`) ise doğrudan gösterilir.
- Sadece dosya adı tutuyorsanız varsayılan olarak `/storage/<dosyaAdı>` üzerinden preview yapılır.
- Farklı storage prefix için `WithProps("storageUrl", "/uploads")` veya `WithProps("storageURL", "...")` verebilirsiniz.

### Repeater Fields

Dinamik olarak tekrarlanan alan grupları oluşturun.

```go
// Repeater field tanımı
fields.NewField("Telefon Numaraları", "phone_numbers").
	OnForm().
	Fields(
		fields.Text("Tip", "type").Required(),
		fields.Tel("Numara", "number").Required(),
	).
	MinRepeats(1).
	MaxRepeats(5).
	HelpText("En az 1, en fazla 5 telefon numarası ekleyebilirsiniz")

// Adres repeater örneği
fields.NewField("Adresler", "addresses").
	OnForm().
	Fields(
		fields.Text("Başlık", "title").Required(),
		fields.Text("Adres", "address").Required(),
		fields.Text("Şehir", "city").Required(),
		fields.Text("Posta Kodu", "postal_code"),
	).
	MinRepeats(1).
	MaxRepeats(3)
```

### Zengin Metin Editörü Yapılandırması

Rich text editörünü detaylı yapılandırın.

```go
// TipTap editör
fields.RichText("İçerik", "content").
	OnForm().
	WithEditor("tiptap").
	WithLanguage("tr").
	WithTheme("default").
	Required()

// Quill editör
fields.RichText("Açıklama", "description").
	OnForm().
	WithEditor("quill").
	WithTheme("snow").
	HelpText("Zengin metin formatı kullanabilirsiniz")

// TinyMCE editör
fields.RichText("Makale", "article").
	OnForm().
	WithEditor("tinymce").
	WithLanguage("tr_TR")
```

### Durum Renkleri ve Badge'ler

Durum alanlarına renkli badge'ler ekleyin.

```go
// Durum renkleri
fields.Select("Durum", "status").
	OnList().
	OnDetail().
	OnForm().
	Options(map[string]string{
		"draft":     "Taslak",
		"published": "Yayınlandı",
		"archived":  "Arşivlendi",
	}).
	WithStatusColors(map[string]string{
		"draft":     "gray",
		"published": "green",
		"archived":  "red",
	}).
	WithBadgeVariant("solid")

// Boolean için renkli gösterim
fields.Switch("Aktif", "is_active").
	OnList().
	WithStatusColors(map[string]string{
		"true":  "green",
		"false": "red",
	})
```

### Money Field (Currency + Intl + Mask)

Money alanı, para tutarlarını currency metadata ile saklamak ve frontend'de
tarayıcı locale'ine göre `Intl.NumberFormat` ile göstermek için kullanılır.

```go
fields.Money("Price", "price").
	OnList().
	OnDetail().
	OnForm().
	Required().
	CurrencyEnum(fields.CurrencyTRY). // varsayılan para birimi
	Currencies(
		fields.CurrencyTRY,
		fields.CurrencyUSD,
		fields.CurrencyEUR,
	).
	AllowCustomCurrency(true). // custom currency kodları tanımlamaya izin
	CustomCurrencies("SAR", "AED"). // enum dışı ek kodlar
	Mask("999999999999999.99"). // form input mask
	MaskChar("_").
	ShowCurrency(true)
```

Notlar:

- `Money` field backend'de `TYPE_MONEY`, frontend'de `money-field` olarak çözülür.
- List/detail görünümlerinde para değeri `Intl.NumberFormat` ile browser locale'ine göre formatlanır.
- `mask` verildiğinde form input'u masked text olarak çalışır.
- Currency kodları için enumlar:
  - `fields.CurrencyUSD`
  - `fields.CurrencyEUR`
  - `fields.CurrencyTRY`
  - `fields.CurrencyGBP`
  - `fields.CurrencyJPY`
  - `fields.CurrencyCHF`
  - `fields.CurrencyCAD`
  - `fields.CurrencyAUD`
  - `fields.CurrencyCNY`

### Pivot Fields

Çoktan çoğa ilişkilerde pivot tablo alanları.

```go
// Pivot field tanımı
fields.Text("Rol", "role").
	OnForm().
	AsPivot().
	WithPivotResource("user_roles")

// Pivot field ile ek bilgi
fields.Date("Atanma Tarihi", "assigned_at").
	OnForm().
	AsPivot().
	WithPivotResource("user_roles").
	Default(time.Now())
```

### GORM Yapılandırması

Veritabanı migration ve model oluşturma için GORM yapılandırması.

```go
// Birincil anahtar
fields.ID().
	GormPrimaryKey().
	GormAutoIncrement()

// UUID birincil anahtar
fields.Text("ID", "id").
	OnlyOnDetail().
	GormPrimaryKey().
	GormUUID()

// İndeksli alan
fields.Text("E-posta", "email").
	OnForm().
	Required().
	GormIndex("idx_email").
	GormSize(255)

// Benzersiz indeks
fields.Text("Kullanıcı Adı", "username").
	OnForm().
	Required().
	GormUniqueIndex("idx_username").
	GormSize(100)

// Foreign key
fields.Link("Kategori", "categories", "category_id").
	OnForm().
	GormForeignKey("category_id", "categories.id").
	GormOnDelete("CASCADE").
	GormOnUpdate("CASCADE")

// Özel SQL tipi
fields.Text("Metadata", "metadata").
	OnForm().
	GormType("jsonb").
	GormComment("JSON metadata alanı")

// Fulltext indeks
fields.Text("İçerik", "content").
	OnForm().
	GormFullTextIndex("idx_content_fulltext")

// Soft delete
fields.DateTime("Silinme Tarihi", "deleted_at").
	OnlyOnDetail().
	GormSoftDelete()

// Varsayılan değer
fields.Switch("Aktif", "is_active").
	OnForm().
	GormDefault(true).
	GormComment("Kayıt aktif mi?")
```

### Callback'ler

Alan değerlerini işlemek için callback'ler kullanın.

```go
// Resolve callback - değeri okurken
fields.Text("Tam Ad", "full_name").
	OnList().
	Resolve(func(value interface{}, item interface{}, c *fiber.Ctx) interface{} {
		if user, ok := item.(*User); ok {
			return user.FirstName + " " + user.LastName
		}
		return value
	})

// Modify callback - değeri kaydetmeden önce
fields.Password("Şifre", "password").
	OnForm().
	Modify(func(value interface{}, c *fiber.Ctx) interface{} {
		if password, ok := value.(string); ok {
			// Şifreyi hashle
			hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
			return string(hashed)
		}
		return value
	})

// StoreAs callback - değeri veritabanına kaydetme şekli
fields.Text("Etiketler", "tags").
	OnForm().
	StoreAs(func(ctx *core.ResourceContext) (interface{}, error) {
		// Virgülle ayrılmış string'i array'e çevir
		tags := strings.Split(ctx.Value.(string), ",")
		return tags, nil
	})
```

### Görünürlük Kontrolü

Dinamik görünürlük kontrolü için callback kullanın.

```go
// Kullanıcı rolüne göre görünürlük
fields.Text("Admin Notu", "admin_note").
	OnForm().
	CanSee(func(ctx *core.ResourceContext) bool {
		user := ctx.User
		return user.Role == "admin"
	})

// Kayıt durumuna göre görünürlük
fields.Text("Yayın Tarihi", "published_at").
	OnForm().
	CanSee(func(ctx *core.ResourceContext) bool {
		if ctx.Item != nil {
			if post, ok := ctx.Item.(*Post); ok {
				return post.Status == "published"
			}
		}
		return false
	})

// İzin kontrolü ile görünürlük
fields.Number("Fiyat", "price").
	OnList().
	CanSee(func(ctx *core.ResourceContext) bool {
		return ctx.User.Can("view-prices")
	})
```

### Metin Hizalama

Liste görünümünde metin hizalaması.

```go
// Sağa hizalı (sayılar için)
fields.Number("Fiyat", "price").
	OnList().
	SetTextAlign("right")

// Ortaya hizalı
fields.Switch("Aktif", "is_active").
	OnList().
	SetTextAlign("center")

// Sola hizalı (varsayılan)
fields.Text("Ad", "name").
	OnList().
	SetTextAlign("left")
```

### Props ile Özel Özellikler

Frontend bileşenlerine özel özellikler gönderin.

```go
// Özel props
fields.Text("Renk", "color").
	OnForm().
	WithProps("type", "color").
	WithProps("showAlpha", true)

// Çoklu props
fields.Number("Miktar", "quantity").
	OnForm().
	WithProps("step", 1).
	WithProps("min", 0).
	WithProps("max", 100)

// Number input +/- kontrollerini gizle
fields.Number("Fiyat", "price").
	OnForm().
	HideNumberControls()

// Aynı davranışın açık hali
fields.Number("Stok", "stock").
	OnForm().
	ShowNumberControls(false)
```

## Örnek: Tam Alan Tanımı

```go
type Product struct {
	ID          string
	Name        string
	Description string
	Price       float64
	Category    string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ProductFieldResolver struct{}

func (r *ProductFieldResolver) ResolveFields(ctx *context.Context) []core.Element {
	return []core.Element{
		// ID - Salt okunur, sadece detayda
		(&fields.Schema{
			Key:   "id",
			Name:  "ID",
			View:  "text",
			Props: make(map[string]interface{}),
		}).ReadOnly().OnlyOnDetail(),

		// Ürün Adı - Zorunlu, aranabilir, sıralanabilir
		(&fields.Schema{
			Key:   "name",
			Name:  "Ürün Adı",
			View:  "text",
			Props: make(map[string]interface{}),
		}).OnList().OnDetail().OnForm().Required().Searchable().Sortable(),

		// Açıklama - Çok satırlı, formda
		(&fields.Schema{
			Key:   "description",
			Name:  "Açıklama",
			View:  "textarea",
			Props: make(map[string]interface{}),
		}).OnForm().HelpText("Ürün hakkında detaylı bilgi"),

		// Fiyat - Sayı, zorunlu, sıralanabilir
		(&fields.Schema{
			Key:   "price",
			Name:  "Fiyat",
			View:  "number",
			Props: make(map[string]interface{}),
		}).OnList().OnDetail().OnForm().Required().Min(0).Sortable(),

		// Kategori - Seçim, filtrelenebilir
		(&fields.Schema{
			Key:   "category",
			Name:  "Kategori",
			View:  "select",
			Props: map[string]interface{}{
				"options": map[string]string{
					"electronics": "Elektronik",
					"clothing":    "Giyim",
					"books":       "Kitaplar",
				},
			},
		}).OnList().OnDetail().OnForm().Filterable(),

		// Aktif - Evet/Hayır, listede göster
		(&fields.Schema{
			Key:   "is_active",
			Name:  "Aktif mi?",
			View:  "switch",
			Props: make(map[string]interface{}),
		}).OnList().OnDetail().OnForm(),

		// Oluşturulma Tarihi - Salt okunur, tarih-saat
		(&fields.Schema{
			Key:   "created_at",
			Name:  "Oluşturulma Tarihi",
			View:  "datetime",
			Props: make(map[string]interface{}),
		}).ReadOnly().OnList().OnDetail(),

		// Güncelleme Tarihi - Salt okunur, tarih-saat
		(&fields.Schema{
			Key:   "updated_at",
			Name:  "Güncelleme Tarihi",
			View:  "datetime",
			Props: make(map[string]interface{}),
		}).ReadOnly().OnDetail(),
	}
}
```

## Fluent API

Alanları tanımlarken fluent API kullanarak zincirleme metod çağrıları yapabilirsiniz:

```go
// Yeni syntax ile
fields.Email("E-posta", "email").
	OnList().
	OnDetail().
	OnForm().
	Required().
	Searchable().
	HelpText("Geçerli bir e-posta adresi girin")
```

## Gerçek Dünya Örnekleri

### Blog Sistemi

```go
func (r *PostFieldResolver) ResolveFields(ctx *context.Context) []core.Element {
	return []core.Element{
		fields.ID().OnlyOnDetail(),

		fields.Text("Başlık", "title").
			OnList().
			OnDetail().
			OnForm().
			Required().
			Searchable().
			Sortable().
			MaxLength(200).
			GormIndex("idx_title"),

		fields.Text("Slug", "slug").
			OnList().
			OnForm().
			Required().
			Unique("posts", "slug").
			Pattern("^[a-z0-9-]+$").
			HelpText("URL dostu başlık (örn: merhaba-dunya)"),

		fields.RichText("İçerik", "content").
			OnForm().
			OnDetail().
			WithEditor("tiptap").
			Required().
			GormType("text"),

		fields.Textarea("Özet", "excerpt").
			OnForm().
			MaxLength(500).
			HelpText("Kısa özet (maksimum 500 karakter)"),

		fields.Link("Yazar", "users", "author_id").
			OnList().
			OnDetail().
			OnForm().
			Required().
			GormForeignKey("author_id", "users.id").
			GormOnDelete("CASCADE"),

		fields.Connect("Etiketler", "tags", "tags").
			OnForm().
			OnDetail(),

		fields.Link("Kategori", "categories", "category_id").
			OnList().
			OnForm().
			Required().
			Filterable(),

		fields.Select("Durum", "status").
			OnList().
			OnForm().
			Options(map[string]string{
				"draft":     "Taslak",
				"published": "Yayınlandı",
				"archived":  "Arşivlendi",
			}).
			Default("draft").
			WithStatusColors(map[string]string{
				"draft":     "gray",
				"published": "green",
				"archived":  "red",
			}).
			Filterable(),

		fields.Image("Kapak Görseli", "cover_image").
			OnDetail().
			OnForm().
			Accept("image/jpeg", "image/png", "image/webp").
			MaxSize(5 * 1024 * 1024).
			Store("public", "posts/covers"),

		fields.Switch("Öne Çıkan", "is_featured").
			OnList().
			OnForm().
			Default(false),

		fields.DateTime("Yayın Tarihi", "published_at").
			OnList().
			OnDetail().
			OnForm().
			Sortable(),

		fields.DateTime("Oluşturulma", "created_at").
			OnList().
			OnDetail().
			ReadOnly().
			Sortable(),

		fields.Collection("Yorumlar", "comments", "comments").
			OnDetail(),
	}
}
```

### E-Ticaret Ürün Yönetimi

```go
func (r *ProductFieldResolver) ResolveFields(ctx *context.Context) []core.Element {
	return []core.Element{
		fields.ID().OnlyOnDetail(),

		fields.Text("Ürün Adı", "name").
			OnList().
			OnDetail().
			OnForm().
			Required().
			Searchable().
			Sortable().
			MaxLength(255).
			GormIndex("idx_name"),

		fields.Text("SKU", "sku").
			OnList().
			OnForm().
			Required().
			Unique("products", "sku").
			Pattern("^[A-Z0-9-]+$").
			HelpText("Stok Kodu (örn: PROD-001)"),

		fields.RichText("Açıklama", "description").
			OnForm().
			OnDetail().
			WithEditor("tiptap"),

		fields.Number("Fiyat", "price").
			OnList().
			OnDetail().
			OnForm().
			Required().
			Min(0).
			Sortable().
			SetTextAlign("right").
			Display(func(value interface{}, item interface{}) interface{} {
				if price, ok := value.(float64); ok {
					return fmt.Sprintf("₺%.2f", price)
				}
				return ""
			}).
			GormType("decimal(10,2)"),

		fields.Number("İndirimli Fiyat", "sale_price").
			OnList().
			OnForm().
			Min(0).
			SetTextAlign("right").
			GormType("decimal(10,2)").
			CanSee(func(ctx *core.ResourceContext) bool {
				return ctx.User.Can("manage-prices")
			}),

		fields.Number("Stok", "stock").
			OnList().
			OnForm().
			Required().
			Min(0).
			SetTextAlign("right").
			WithStatusColors(map[string]string{
				"0":  "red",
				"1":  "orange",
				"10": "green",
			}),

		fields.Link("Kategori", "categories", "category_id").
			OnList().
			OnForm().
			Required().
			Filterable().
			GormForeignKey("category_id", "categories.id"),

		fields.Link("Marka", "brands", "brand_id").
			OnList().
			OnForm().
			Filterable(),

		fields.Connect("Etiketler", "tags", "tags").
			OnForm().
			OnDetail(),

		fields.Image("Ana Görsel", "main_image").
			OnList().
			OnDetail().
			OnForm().
			Accept("image/jpeg", "image/png", "image/webp").
			MaxSize(5 * 1024 * 1024).
			Store("public", "products/images").
			RemoveEXIFData(),

		fields.PolyCollection("Görseller", "images", "images").
			OnDetail(),

		fields.KeyValue("Özellikler", "attributes").
			OnForm().
			OnDetail().
			HelpText("Ürün özellikleri (örn: Renk: Kırmızı)"),

		fields.Select("Durum", "status").
			OnList().
			OnForm().
			Options(map[string]string{
				"active":      "Aktif",
				"inactive":    "Pasif",
				"out_of_stock": "Stokta Yok",
			}).
			Default("active").
			WithStatusColors(map[string]string{
				"active":       "green",
				"inactive":     "gray",
				"out_of_stock": "red",
			}).
			Filterable(),

		fields.Switch("Öne Çıkan", "is_featured").
			OnList().
			OnForm(),

		fields.DateTime("Oluşturulma", "created_at").
			OnList().
			OnDetail().
			ReadOnly().
			Sortable(),

		fields.Collection("Yorumlar", "reviews", "reviews").
			OnDetail(),

		fields.Collection("Varyantlar", "variants", "variants").
			OnDetail(),
	}
}
```

### Kullanıcı Yönetimi

```go
func (r *UserFieldResolver) ResolveFields(ctx *context.Context) []core.Element {
	return []core.Element{
		fields.ID().OnlyOnDetail(),

		fields.Text("Ad", "first_name").
			OnList().
			OnDetail().
			OnForm().
			Required().
			Searchable().
			MaxLength(100),

		fields.Text("Soyad", "last_name").
			OnList().
			OnDetail().
			OnForm().
			Required().
			Searchable().
			MaxLength(100),

		fields.Text("Tam Ad", "full_name").
			OnList().
			OnDetail().
			Resolve(func(value interface{}, item interface{}, c *fiber.Ctx) interface{} {
				if user, ok := item.(*User); ok {
					return user.FirstName + " " + user.LastName
				}
				return value
			}),

		fields.Email("E-posta", "email").
			OnList().
			OnDetail().
			OnForm().
			Required().
			Unique("users", "email").
			Searchable().
			GormIndex("idx_email"),

		fields.Password("Şifre", "password").
			OnForm().
			Required().
			MinLength(8).
			HideOnUpdate().
			Modify(func(value interface{}, c *fiber.Ctx) interface{} {
				if password, ok := value.(string); ok {
					hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
					return string(hashed)
				}
				return value
			}),

		fields.Tel("Telefon", "phone").
			OnDetail().
			OnForm().
			Pattern(`^\+?[1-9]\d{1,14}$`),

		fields.Image("Profil Fotoğrafı", "avatar").
			OnDetail().
			OnForm().
			Accept("image/jpeg", "image/png").
			MaxSize(2 * 1024 * 1024).
			Store("public", "avatars").
			RemoveEXIFData(),

		fields.Date("Doğum Tarihi", "birth_date").
			OnDetail().
			OnForm(),

		fields.Select("Cinsiyet", "gender").
			OnDetail().
			OnForm().
			Options(map[string]string{
				"male":   "Erkek",
				"female": "Kadın",
				"other":  "Diğer",
			}),

		fields.Connect("Roller", "roles", "roles").
			OnForm().
			OnDetail().
			Required().
			CanSee(func(ctx *core.ResourceContext) bool {
				return ctx.User.Can("manage-roles")
			}),

		fields.Switch("Aktif", "is_active").
			OnList().
			OnForm().
			Default(true).
			CanSee(func(ctx *core.ResourceContext) bool {
				return ctx.User.Can("manage-users")
			}),

		fields.Switch("E-posta Doğrulandı", "email_verified").
			OnList().
			OnDetail().
			ReadOnly(),

		fields.DateTime("Son Giriş", "last_login_at").
			OnList().
			OnDetail().
			ReadOnly().
			Sortable(),

		fields.DateTime("Kayıt Tarihi", "created_at").
			OnList().
			OnDetail().
			ReadOnly().
			Sortable(),

		fields.Collection("Siparişler", "orders", "orders").
			OnDetail(),

		fields.Collection("Adresler", "addresses", "addresses").
			OnDetail(),
	}
}
```

## Best Practices ve İpuçları

### Performans

1. **İndeksleme**: Aranabilir ve sıralanabilir alanlar için mutlaka veritabanı indeksi oluşturun
   ```go
   fields.Text("E-posta", "email").
       Searchable().
       Sortable().
       GormIndex("idx_email")
   ```

2. **Eager Loading**: İlişki alanlarında N+1 problemini önlemek için eager loading kullanın
   ```go
   fields.Link("Kategori", "categories", "category_id").
       OnList() // Otomatik eager loading
   ```

3. **Fulltext Search**: Büyük metin alanlarında fulltext indeks kullanın
   ```go
   fields.Text("İçerik", "content").
       Searchable().
       GormFullTextIndex("idx_content_fulltext")
   ```

### Güvenlik

1. **Hassas Alanlar**: Hassas bilgileri salt okunur yapın veya görünürlüğünü kısıtlayın
   ```go
   fields.Text("Kredi Kartı", "credit_card").
       OnForm().
       CanSee(func(ctx *core.ResourceContext) bool {
           return ctx.User.Can("view-payment-info")
       })
   ```

2. **Şifre Güvenliği**: Şifreleri mutlaka hashleyin
   ```go
   fields.Password("Şifre", "password").
       OnForm().
       Modify(func(value interface{}, c *fiber.Ctx) interface{} {
           hashed, _ := bcrypt.GenerateFromPassword([]byte(value.(string)), bcrypt.DefaultCost)
           return string(hashed)
       })
   ```

3. **XSS Koruması**: Kullanıcı girdilerini doğrulayın ve temizleyin
   ```go
   fields.Text("Kullanıcı Adı", "username").
       Pattern("^[a-zA-Z0-9_]+$").
       MaxLength(50)
   ```

### Validasyon

1. **Çift Taraflı Validasyon**: Hem frontend hem backend'de validasyon yapın
   ```go
   fields.Email("E-posta", "email").
       Required().
       Email().
       MaxLength(255)
   ```

2. **Özel Validasyon**: Karmaşık kurallar için özel validasyon kullanın
   ```go
   fields.Text("Vergi No", "tax_number").
       DependsOn("is_company").
       When("is_company", "=", true).
       Required()
   ```

3. **Benzersizlik Kontrolü**: Benzersiz olması gereken alanları işaretleyin
   ```go
   fields.Text("E-posta", "email").
       Unique("users", "email").
       GormUniqueIndex("idx_email")
   ```

### Kullanıcı Deneyimi

1. **Yardım Metinleri**: Kullanıcıları yönlendirmek için yardım metni ekleyin
   ```go
   fields.Password("Şifre", "password").
       HelpText("En az 8 karakter, bir büyük harf ve bir rakam içermelidir")
   ```

2. **Placeholder**: Örnek değerler gösterin
   ```go
   fields.Tel("Telefon", "phone").
       Placeholder("+90 (555) 123-4567")
   ```

3. **Varsayılan Değerler**: Mantıklı varsayılan değerler belirleyin
   ```go
   fields.Switch("Aktif", "is_active").
       Default(true)
   ```

4. **Görsel Geri Bildirim**: Durum renklerini kullanın
   ```go
   fields.Select("Durum", "status").
       WithStatusColors(map[string]string{
           "success": "green",
           "warning": "orange",
           "error":   "red",
       })
   ```

### Veritabanı Tasarımı

1. **Foreign Key Constraints**: İlişkilerde foreign key kullanın
   ```go
   fields.Link("Kategori", "categories", "category_id").
       GormForeignKey("category_id", "categories.id").
       GormOnDelete("CASCADE").
       GormOnUpdate("CASCADE")
   ```

2. **Soft Delete**: Silinen kayıtları korumak için soft delete kullanın
   ```go
   fields.DateTime("Silinme Tarihi", "deleted_at").
       GormSoftDelete()
   ```

3. **Uygun Veri Tipleri**: Doğru SQL tiplerini kullanın
   ```go
   fields.Number("Fiyat", "price").
       GormType("decimal(10,2)")

   fields.Text("Metadata", "metadata").
       GormType("jsonb")
   ```

### Organizasyon

1. **Mantıksal Gruplama**: İlgili alanları gruplandırın
   ```go
   // Temel bilgiler
   fields.Text("Ad", "name"),
   fields.Text("Soyad", "surname"),

   // İletişim bilgileri
   fields.Email("E-posta", "email"),
   fields.Tel("Telefon", "phone"),

   // Sistem alanları
   fields.DateTime("Oluşturulma", "created_at"),
   fields.DateTime("Güncellenme", "updated_at"),
   ```

2. **Tutarlı İsimlendirme**: Tutarlı alan isimleri kullanın
   ```go
   // İyi
   fields.DateTime("created_at")
   fields.DateTime("updated_at")
   fields.DateTime("deleted_at")

   // Kötü
   fields.DateTime("createdDate")
   fields.DateTime("update_time")
   fields.DateTime("deletedOn")
   ```

3. **Yorum Kullanımı**: Karmaşık alanları açıklayın
   ```go
   // Kullanıcının son 30 gün içindeki toplam harcaması
   fields.Number("Son 30 Gün Harcama", "spending_last_30_days").
       OnDetail().
       ReadOnly()
   ```

## Frontend Bileşenleri

Panel.go'nun React tabanlı frontend bileşenleri, backend field tanımlamalarını kullanıcı arayüzünde görselleştirir. Bu bölümde frontend bileşenlerinin özel özellikleri ve kullanımları açıklanmaktadır.

### TextInput - Maskeli Metin Girişi

TextInput bileşeni, `react-input-mask` kütüphanesi kullanılarak geliştirilmiş maskeli metin girişi desteği sağlar. Bu özellik sayesinde telefon numarası, TC kimlik no, tarih, kredi kartı, IBAN gibi formatlanmış veri girişleri kolayca yapılabilir.

#### Mask Prop'u

TextInput bileşenine `mask` prop'u ekleyerek input maskesi tanımlayabilirsiniz:

```tsx
// Frontend (React/TypeScript)
<TextInput
  name="phone"
  label="Telefon Numarası"
  value={phone}
  onChange={setPhone}
  mask="(599) 999 99 99"
  placeholder="(5XX) XXX XX XX"
/>
```

#### Maske Karakterleri

Mask tanımlarken kullanabileceğiniz özel karakterler:

- **9**: Rakam (0-9)
- **a**: Harf (a-z, A-Z)
- **\***: Alfanumerik (harf veya rakam)

Diğer tüm karakterler (parantez, tire, boşluk vb.) olduğu gibi gösterilir.

#### Mask Özellikleri

TextInput bileşeni mask için üç prop kabul eder:

| Prop | Tip | Varsayılan | Açıklama |
|------|-----|-----------|----------|
| `mask` | `string` | `undefined` | Input maskesi formatı |
| `maskChar` | `string` | `"_"` | Boş karakterler için gösterilecek karakter |
| `alwaysShowMask` | `boolean` | `false` | Maskeyi her zaman göster (focus olmasa bile) |

#### Kullanım Örnekleri

##### Telefon Numarası (Türkiye)

```tsx
<TextInput
  name="phone"
  label="Telefon Numarası"
  value={phone}
  onChange={setPhone}
  mask="(599) 999 99 99"
  placeholder="(5XX) XXX XX XX"
  required
/>
```

**Backend Field Tanımı:**
```go
fields.Tel("Telefon", "phone").
	OnList().
	OnDetail().
	OnForm().
	Required().
	Pattern(`^\(5\d{2}\) \d{3} \d{2} \d{2}$`).
	HelpText("Türkiye cep telefonu formatında giriniz")
```

##### TC Kimlik Numarası

```tsx
<TextInput
  name="tc_no"
  label="TC Kimlik No"
  value={tcNo}
  onChange={setTcNo}
  mask="99999999999"
  placeholder="TC Kimlik No"
  required
/>
```

**Backend Field Tanımı:**
```go
fields.Text("TC Kimlik No", "tc_no").
	OnForm().
	OnDetail().
	Required().
	Pattern(`^\d{11}$`).
	MinLength(11).
	MaxLength(11).
	HelpText("11 haneli TC kimlik numaranızı giriniz")
```

##### Tarih Girişi

```tsx
<TextInput
  name="birth_date"
  label="Doğum Tarihi"
  value={birthDate}
  onChange={setBirthDate}
  mask="99/99/9999"
  placeholder="GG/AA/YYYY"
  maskChar="_"
  alwaysShowMask
/>
```

**Backend Field Tanımı:**
```go
fields.Date("Doğum Tarihi", "birth_date").
	OnForm().
	OnDetail().
	Required().
	HelpText("Doğum tarihinizi GG/AA/YYYY formatında giriniz")
```

##### Kredi Kartı Numarası

```tsx
<TextInput
  name="card_number"
  label="Kart Numarası"
  value={cardNumber}
  onChange={setCardNumber}
  mask="9999 9999 9999 9999"
  placeholder="Kart Numarası"
  required
/>
```

**Backend Field Tanımı:**
```go
fields.Text("Kart Numarası", "card_number").
	OnForm().
	Required().
	Pattern(`^\d{4} \d{4} \d{4} \d{4}$`).
	MinLength(19).
	MaxLength(19).
	HelpText("16 haneli kart numaranızı giriniz").
	Modify(func(value interface{}, c *fiber.Ctx) interface{} {
		// Boşlukları kaldır ve veritabanına kaydet
		if cardNo, ok := value.(string); ok {
			return strings.ReplaceAll(cardNo, " ", "")
		}
		return value
	})
```

##### IBAN Numarası (Türkiye)

```tsx
<TextInput
  name="iban"
  label="IBAN"
  value={iban}
  onChange={setIban}
  mask="TR99 9999 9999 9999 9999 9999 99"
  placeholder="TR00 0000 0000 0000 0000 0000 00"
  required
/>
```

**Backend Field Tanımı:**
```go
fields.Text("IBAN", "iban").
	OnForm().
	OnDetail().
	Required().
	Pattern(`^TR\d{2} \d{4} \d{4} \d{4} \d{4} \d{4} \d{2}$`).
	MinLength(32).
	MaxLength(32).
	HelpText("Türkiye IBAN numaranızı giriniz").
	Modify(func(value interface{}, c *fiber.Ctx) interface{} {
		// Boşlukları kaldır ve veritabanına kaydet
		if iban, ok := value.(string); ok {
			return strings.ReplaceAll(iban, " ", "")
		}
		return value
	})
```

##### Posta Kodu

```tsx
<TextInput
  name="postal_code"
  label="Posta Kodu"
  value={postalCode}
  onChange={setPostalCode}
  mask="99999"
  placeholder="Posta Kodu"
/>
```

**Backend Field Tanımı:**
```go
fields.Text("Posta Kodu", "postal_code").
	OnForm().
	Pattern(`^\d{5}$`).
	MinLength(5).
	MaxLength(5).
	HelpText("5 haneli posta kodunu giriniz")
```

##### Saat Girişi

```tsx
<TextInput
  name="time"
  label="Saat"
  value={time}
  onChange={setTime}
  mask="99:99"
  placeholder="SS:DD"
  maskChar="_"
/>
```

**Backend Field Tanımı:**
```go
fields.Text("Saat", "time").
	OnForm().
	Pattern(`^([0-1][0-9]|2[0-3]):[0-5][0-9]$`).
	HelpText("Saat formatında giriniz (örn: 14:30)")
```

#### Backend Entegrasyonu

Maskeli input'lardan gelen veriler genellikle formatlanmış şekilde gelir (örn: "(555) 123 45 67"). Backend'de bu verileri işlerken formatı temizlemek veya doğrulamak gerekebilir:

```go
// Telefon numarası için Modify callback
fields.Tel("Telefon", "phone").
	OnForm().
	Modify(func(value interface{}, c *fiber.Ctx) interface{} {
		if phone, ok := value.(string); ok {
			// Sadece rakamları al
			re := regexp.MustCompile(`\D`)
			cleaned := re.ReplaceAllString(phone, "")
			return cleaned
		}
		return value
	})

// IBAN için Modify callback
fields.Text("IBAN", "iban").
	OnForm().
	Modify(func(value interface{}, c *fiber.Ctx) interface{} {
		if iban, ok := value.(string); ok {
			// Boşlukları kaldır
			return strings.ReplaceAll(iban, " ", "")
		}
		return value
	})

// Kredi kartı için Modify callback
fields.Text("Kart Numarası", "card_number").
	OnForm().
	Modify(func(value interface{}, c *fiber.Ctx) interface{} {
		if cardNo, ok := value.(string); ok {
			// Boşlukları kaldır ve şifrele
			cleaned := strings.ReplaceAll(cardNo, " ", "")
			encrypted, _ := encrypt(cleaned)
			return encrypted
		}
		return value
	})
```

#### Resolve Callback ile Formatlama

Veritabanından gelen temiz veriyi frontend'de formatlanmış şekilde göstermek için Resolve callback kullanabilirsiniz:

```go
// Telefon numarasını formatla
fields.Tel("Telefon", "phone").
	OnList().
	OnDetail().
	Resolve(func(value interface{}, item interface{}, c *fiber.Ctx) interface{} {
		if phone, ok := value.(string); ok && len(phone) == 10 {
			// 5551234567 -> (555) 123 45 67
			return fmt.Sprintf("(%s) %s %s %s",
				phone[0:3],
				phone[3:6],
				phone[6:8],
				phone[8:10],
			)
		}
		return value
	})

// IBAN'ı formatla
fields.Text("IBAN", "iban").
	OnList().
	OnDetail().
	Resolve(func(value interface{}, item interface{}, c *fiber.Ctx) interface{} {
		if iban, ok := value.(string); ok && len(iban) == 26 {
			// TR123456789012345678901234 -> TR12 3456 7890 1234 5678 9012 34
			formatted := ""
			for i := 0; i < len(iban); i += 4 {
				end := i + 4
				if end > len(iban) {
					end = len(iban)
				}
				if i > 0 {
					formatted += " "
				}
				formatted += iban[i:end]
			}
			return formatted
		}
		return value
	})
```

#### Validasyon

Maskeli input'lar için backend validasyonu önemlidir:

```go
// Telefon numarası validasyonu
fields.Tel("Telefon", "phone").
	OnForm().
	Required().
	Pattern(`^5\d{9}$`). // Sadece rakamlar, 10 haneli, 5 ile başlayan
	MinLength(10).
	MaxLength(10).
	HelpText("Geçerli bir Türkiye cep telefonu numarası giriniz")

// TC Kimlik No validasyonu
fields.Text("TC Kimlik No", "tc_no").
	OnForm().
	Required().
	Pattern(`^\d{11}$`).
	MinLength(11).
	MaxLength(11).
	Modify(func(value interface{}, c *fiber.Ctx) interface{} {
		if tcNo, ok := value.(string); ok {
			// TC Kimlik No algoritması ile doğrula
			if !isValidTCNo(tcNo) {
				return errors.New("Geçersiz TC Kimlik No")
			}
			return tcNo
		}
		return value
	})

// IBAN validasyonu
fields.Text("IBAN", "iban").
	OnForm().
	Required().
	Pattern(`^TR\d{24}$`). // Boşluksuz format
	MinLength(26).
	MaxLength(26).
	Modify(func(value interface{}, c *fiber.Ctx) interface{} {
		if iban, ok := value.(string); ok {
			// IBAN checksum doğrulaması
			if !isValidIBAN(iban) {
				return errors.New("Geçersiz IBAN numarası")
			}
			return iban
		}
		return value
	})
```

#### Best Practices

1. **Tutarlı Format Kullanımı**: Aynı veri tipi için her zaman aynı mask formatını kullanın
   ```tsx
   // İyi ✓
   mask="(599) 999 99 99"  // Tüm telefon alanlarında

   // Kötü ✗
   mask="(599) 999 99 99"  // Bir yerde
   mask="599 999 99 99"    // Başka bir yerde
   ```

2. **Backend Temizleme**: Maskeli veriler backend'e gönderilmeden önce temizlenmelidir
   ```go
   // Modify callback ile temizleme
   Modify(func(value interface{}, c *fiber.Ctx) interface{} {
       if str, ok := value.(string); ok {
           return strings.ReplaceAll(str, " ", "")
       }
       return value
   })
   ```

3. **Placeholder Kullanımı**: Kullanıcıya beklenen formatı gösterin
   ```tsx
   <TextInput
     mask="(599) 999 99 99"
     placeholder="(5XX) XXX XX XX"  // Format örneği
   />
   ```

4. **Yardım Metni**: Karmaşık formatlar için açıklama ekleyin
   ```tsx
   <TextInput
     mask="TR99 9999 9999 9999 9999 9999 99"
     placeholder="TR00 0000 0000 0000 0000 0000 00"
     helpText="Türkiye IBAN numaranızı TR ile başlayarak giriniz"
   />
   ```

5. **Validasyon**: Hem frontend hem backend'de validasyon yapın
   ```tsx
   // Frontend
   <TextInput
     mask="99999999999"
     required
     error={tcNoError}
   />
   ```
   ```go
   // Backend
   fields.Text("TC Kimlik No", "tc_no").
       Required().
       Pattern(`^\d{11}$`).
       MinLength(11).
       MaxLength(11)
   ```

6. **Erişilebilirlik**: Mask kullanırken erişilebilirlik özelliklerini koruyun
   ```tsx
   <TextInput
     mask="(599) 999 99 99"
     aria-label="Telefon Numarası"
     aria-describedby="phone-help"
   />
   ```

#### Teknik Detaylar

- **Paket**: `react-input-mask` (v2.0.4+)
- **TypeScript Desteği**: `@types/react-input-mask`
- **Bileşen Konumu**: `web/src/components/fields/TextInput.tsx`
- **Shadcn/ui Entegrasyonu**: TextInput, shadcn/ui Input bileşenini kullanır
- **Performans**: Mask işlemleri client-side'da yapılır, backend'e minimal yük

#### Sınırlamalar

1. **Dinamik Maskeler**: Mask değeri runtime'da değiştirilemez (component re-render gerektirir)
2. **Karmaşık Formatlar**: Çok karmaşık formatlar için özel regex validasyonu gerekebilir
3. **Mobil Klavye**: Mobil cihazlarda sayısal klavye için `type="tel"` kullanımı önerilir
4. **Copy-Paste**: Kullanıcı formatlanmamış veri yapıştırırsa mask otomatik uygulanır

---

### TelInput - Telefon Numarası Girişi

TelInput bileşeni, telefon numarası girişi için özel olarak tasarlanmış esnek bir bileşendir. İki farklı mod destekler:

1. **PhoneInput Modu (Gelişmiş)**: Uluslararası telefon numarası girişi, ülke seçimi ve otomatik formatlama
2. **Native Modu (Basit)**: HTML tel input ile opsiyonel mask desteği

#### Mod Seçimi

TelInput bileşeni, `usePhoneInput` prop'una göre otomatik olarak uygun modu kullanır:

```tsx
// PhoneInput modu (gelişmiş) - Uluslararası telefon numaraları için
<TelInput
  name="phone"
  label="Telefon Numarası"
  value={phone}
  onChange={setPhone}
  usePhoneInput
  defaultCountry="TR"
/>

// Native modu (basit) - Yerel telefon numaraları için
<TelInput
  name="phone"
  label="Telefon Numarası"
  value={phone}
  onChange={setPhone}
  mask="(599) 999 99 99"
  placeholder="(5XX) XXX XX XX"
/>
```

#### PhoneInput Modu (Gelişmiş)

PhoneInput modu, `react-phone-number-input` kütüphanesi kullanarak uluslararası telefon numarası desteği sağlar.

**Özellikler:**
- Ülke bayrağı ve telefon kodu seçimi
- Otomatik telefon numarası formatlaması
- E.164 formatında değer döndürme (+905551234567)
- Arama yapılabilir ülke listesi
- 200+ ülke desteği

**Kullanım Örneği:**

```tsx
import { TelInput } from '@/components/fields/TelInput';

function ContactForm() {
  const [phone, setPhone] = useState('');

  return (
    <TelInput
      name="phone"
      label="Telefon Numarası"
      value={phone}
      onChange={setPhone}
      usePhoneInput
      defaultCountry="TR"
      placeholder="Telefon numaranızı girin"
      required
      helpText="Uluslararası format kullanılacaktır"
    />
  );
}
```

**Backend Field Tanımı:**

```go
fields.Tel("Telefon", "phone").
	OnList().
	OnDetail().
	OnForm().
	Required().
	Pattern(`^\+[1-9]\d{1,14}$`). // E.164 format
	HelpText("Uluslararası telefon numarası formatında giriniz").
	Resolve(func(value interface{}, item interface{}, c *fiber.Ctx) interface{} {
		// E.164 formatını görsel formata çevir
		if phone, ok := value.(string); ok && len(phone) > 0 {
			// +905551234567 -> +90 (555) 123 45 67
			return formatPhoneNumber(phone)
		}
		return value
	})
```

#### Native Modu (Basit)

Native modu, basit HTML tel input kullanır ve opsiyonel olarak mask desteği sağlar.

**Özellikler:**
- Hafif ve hızlı
- Opsiyonel input mask desteği
- Mobil cihazlarda sayısal klavye
- Basit validasyon

**Maskeli Kullanım:**

```tsx
// Türkiye telefon numarası formatı
<TelInput
  name="phone"
  label="Telefon Numarası"
  value={phone}
  onChange={setPhone}
  mask="(599) 999 99 99"
  placeholder="(5XX) XXX XX XX"
  required
/>
```

**Maskesiz Kullanım:**

```tsx
// Basit tel input
<TelInput
  name="phone"
  label="Telefon Numarası"
  value={phone}
  onChange={setPhone}
  placeholder="Telefon numaranızı girin"
  required
/>
```

**Backend Field Tanımı:**

```go
fields.Tel("Telefon", "phone").
	OnList().
	OnDetail().
	OnForm().
	Required().
	Pattern(`^\(5\d{2}\) \d{3} \d{2} \d{2}$`). // Maskeli format
	HelpText("Türkiye cep telefonu formatında giriniz").
	Modify(func(value interface{}, c *fiber.Ctx) interface{} {
		// Formatı temizle ve sadece rakamları al
		if phone, ok := value.(string); ok {
			re := regexp.MustCompile(`\D`)
			cleaned := re.ReplaceAllString(phone, "")
			return cleaned
		}
		return value
	})
```

#### Props Karşılaştırması

| Prop | PhoneInput Modu | Native Modu | Açıklama |
|------|----------------|-------------|----------|
| `usePhoneInput` | `true` | `false` / `undefined` | Mod seçimi |
| `defaultCountry` | ✅ | ❌ | Varsayılan ülke kodu (örn: "TR") |
| `mask` | ❌ | ✅ | Input maskesi formatı |
| `maskChar` | ❌ | ✅ | Mask için boş karakter |
| `alwaysShowMask` | ❌ | ✅ | Maskeyi her zaman göster |
| `placeholder` | ✅ | ✅ | Placeholder metni |
| `required` | ✅ | ✅ | Zorunlu alan |
| `disabled` | ✅ | ✅ | Devre dışı |
| `error` | ✅ | ✅ | Hata mesajı |
| `helpText` | ✅ | ✅ | Yardım metni |

#### Kullanım Senaryoları

##### Senaryo 1: Uluslararası İş Uygulaması

Farklı ülkelerden kullanıcıların telefon numarası girmesi gerekiyorsa PhoneInput modunu kullanın:

```tsx
<TelInput
  name="phone"
  label="Phone Number"
  value={phone}
  onChange={setPhone}
  usePhoneInput
  defaultCountry="US"
  placeholder="Enter your phone number"
  required
  helpText="We'll use this to contact you"
/>
```

```go
// Backend - E.164 format beklenir
fields.Tel("Phone", "phone").
	OnForm().
	Required().
	Pattern(`^\+[1-9]\d{1,14}$`).
	HelpText("International phone number format")
```

##### Senaryo 2: Yerel Türkiye Uygulaması

Sadece Türkiye telefon numaraları için Native mod ile mask kullanın:

```tsx
<TelInput
  name="phone"
  label="Telefon Numarası"
  value={phone}
  onChange={setPhone}
  mask="(599) 999 99 99"
  placeholder="(5XX) XXX XX XX"
  required
  helpText="Türkiye cep telefonu numaranızı girin"
/>
```

```go
// Backend - Maskeli format beklenir
fields.Tel("Telefon", "phone").
	OnForm().
	Required().
	Pattern(`^\(5\d{2}\) \d{3} \d{2} \d{2}$`).
	Modify(func(value interface{}, c *fiber.Ctx) interface{} {
		// Sadece rakamları al: (555) 123 45 67 -> 5551234567
		if phone, ok := value.(string); ok {
			re := regexp.MustCompile(`\D`)
			return re.ReplaceAllString(phone, "")
		}
		return value
	})
```

##### Senaryo 3: Basit Form

Mask olmadan basit telefon girişi:

```tsx
<TelInput
  name="phone"
  label="Telefon"
  value={phone}
  onChange={setPhone}
  placeholder="05551234567"
/>
```

```go
// Backend - Basit validasyon
fields.Tel("Telefon", "phone").
	OnForm().
	Pattern(`^0\d{10}$`).
	MinLength(11).
	MaxLength(11)
```

#### Backend Entegrasyonu

##### PhoneInput Modu için Backend

PhoneInput modu E.164 formatında değer döndürür (+905551234567). Backend'de bu formatı işleyin:

```go
// Validasyon
fields.Tel("Telefon", "phone").
	OnForm().
	Required().
	Pattern(`^\+[1-9]\d{1,14}$`). // E.164 format
	HelpText("Uluslararası telefon numarası").
	Modify(func(value interface{}, c *fiber.Ctx) interface{} {
		if phone, ok := value.(string); ok {
			// E.164 formatını doğrula
			if !isValidE164(phone) {
				return errors.New("Geçersiz telefon numarası formatı")
			}
			return phone
		}
		return value
	})

// Görüntüleme için formatlama
fields.Tel("Telefon", "phone").
	OnList().
	OnDetail().
	Resolve(func(value interface{}, item interface{}, c *fiber.Ctx) interface{} {
		if phone, ok := value.(string); ok {
			// +905551234567 -> +90 (555) 123 45 67
			return formatE164ToDisplay(phone)
		}
		return value
	})
```

##### Native Modu için Backend

Native mod maskeli veya maskesiz değer döndürür. Backend'de temizleme yapın:

```go
// Maskeli format için
fields.Tel("Telefon", "phone").
	OnForm().
	Required().
	Pattern(`^\(5\d{2}\) \d{3} \d{2} \d{2}$`).
	Modify(func(value interface{}, c *fiber.Ctx) interface{} {
		if phone, ok := value.(string); ok {
			// (555) 123 45 67 -> 5551234567
			re := regexp.MustCompile(`\D`)
			cleaned := re.ReplaceAllString(phone, "")

			// Türkiye telefon numarası validasyonu
			if !strings.HasPrefix(cleaned, "5") || len(cleaned) != 10 {
				return errors.New("Geçersiz Türkiye telefon numarası")
			}

			return cleaned
		}
		return value
	}).
	Resolve(func(value interface{}, item interface{}, c *fiber.Ctx) interface{} {
		if phone, ok := value.(string); ok && len(phone) == 10 {
			// 5551234567 -> (555) 123 45 67
			return fmt.Sprintf("(%s) %s %s %s",
				phone[0:3],
				phone[3:6],
				phone[6:8],
				phone[8:10],
			)
		}
		return value
	})
```

#### Validasyon Örnekleri

##### PhoneInput Modu Validasyonu

```go
// E.164 format validasyonu
func isValidE164(phone string) bool {
	// +[1-9][0-9]{1,14}
	re := regexp.MustCompile(`^\+[1-9]\d{1,14}$`)
	return re.MatchString(phone)
}

// Ülkeye özel validasyon
func validatePhoneByCountry(phone string, country string) error {
	switch country {
	case "TR":
		// Türkiye: +90 ile başlamalı, 13 karakter
		if !strings.HasPrefix(phone, "+90") || len(phone) != 13 {
			return errors.New("Geçersiz Türkiye telefon numarası")
		}
	case "US":
		// ABD: +1 ile başlamalı, 12 karakter
		if !strings.HasPrefix(phone, "+1") || len(phone) != 12 {
			return errors.New("Geçersiz ABD telefon numarası")
		}
	}
	return nil
}
```

##### Native Modu Validasyonu

```go
// Türkiye cep telefonu validasyonu
func isValidTurkishMobile(phone string) bool {
	// Sadece rakamlar, 10 haneli, 5 ile başlayan
	re := regexp.MustCompile(`^5\d{9}$`)
	return re.MatchString(phone)
}

// Maskeli format validasyonu
func isValidMaskedPhone(phone string) bool {
	// (5XX) XXX XX XX formatı
	re := regexp.MustCompile(`^\(5\d{2}\) \d{3} \d{2} \d{2}$`)
	return re.MatchString(phone)
}
```

#### Best Practices

1. **Doğru Modu Seçin**
   - Uluslararası kullanıcılar → PhoneInput modu
   - Tek ülke, yerel kullanıcılar → Native mod (maskeli)
   - Basit formlar → Native mod (maskesiz)

2. **Backend Temizleme**
   ```go
   // Her zaman backend'de formatı temizleyin
   Modify(func(value interface{}, c *fiber.Ctx) interface{} {
       if phone, ok := value.(string); ok {
           // Boşluk, parantez, tire vb. kaldır
           re := regexp.MustCompile(`\D`)
           return re.ReplaceAllString(phone, "")
       }
       return value
   })
   ```

3. **Tutarlı Format**
   ```tsx
   // İyi ✓ - Tüm projede aynı format
   <TelInput mask="(599) 999 99 99" />

   // Kötü ✗ - Farklı formatlar
   <TelInput mask="(599) 999 99 99" />  // Bir yerde
   <TelInput mask="599 999 99 99" />    // Başka yerde
   ```

4. **Placeholder Kullanımı**
   ```tsx
   // Format örneği gösterin
   <TelInput
     mask="(599) 999 99 99"
     placeholder="(5XX) XXX XX XX"
   />
   ```

5. **Yardım Metni**
   ```tsx
   // Kullanıcıyı yönlendirin
   <TelInput
     usePhoneInput
     helpText="Ülke kodunu seçip telefon numaranızı girin"
   />
   ```

6. **Hata Mesajları**
   ```tsx
   // Açıklayıcı hata mesajları
   <TelInput
     error="Geçerli bir Türkiye cep telefonu numarası giriniz (5XX ile başlamalı)"
   />
   ```

#### Teknik Detaylar

**PhoneInput Modu:**
- **Paket**: `react-phone-number-input` (v3.4.0+)
- **Bileşen**: `web/src/components/ui/phone-input.tsx`
- **Format**: E.164 (+905551234567)
- **Ülke Sayısı**: 200+
- **Bundle Boyutu**: ~50KB (gzipped)

**Native Modu:**
- **Paket**: `react-input-mask` (v2.0.4+)
- **Bileşen**: `web/src/components/fields/TelInput.tsx`
- **Format**: Özelleştirilebilir
- **Bundle Boyutu**: ~5KB (gzipped)

#### Performans Karşılaştırması

| Özellik | PhoneInput | Native (Maskeli) | Native (Maskesiz) |
|---------|-----------|------------------|-------------------|
| Bundle Boyutu | ~50KB | ~5KB | ~1KB |
| İlk Render | ~100ms | ~20ms | ~10ms |
| Ülke Desteği | 200+ | 1 | 1 |
| Otomatik Format | ✅ | ✅ | ❌ |
| Validasyon | ✅ | Kısmi | ❌ |

#### Sınırlamalar

**PhoneInput Modu:**
1. **Bundle Boyutu**: Daha büyük bundle boyutu (~50KB)
2. **Performans**: İlk render daha yavaş
3. **Özelleştirme**: Stil özelleştirmesi sınırlı
4. **Mobil**: Mobil cihazlarda bazen klavye sorunları

**Native Modu:**
1. **Ülke Desteği**: Tek ülke için optimize
2. **Validasyon**: Manuel validasyon gerekli
3. **Format**: Otomatik format düzeltme yok
4. **Uluslararası**: Uluslararası numaralar için uygun değil

---

## Tooltip Desteği

Tüm form field'larında tooltip desteği mevcuttur. Tooltip, label'ın yanında bir info ikonu olarak gösterilir ve kullanıcıya ek bilgi sağlar.

### Kullanım

**Backend Field Tanımı:**
```go
fields.Text("Kullanıcı Adı", "username").
	OnForm().
	Required().
	WithProps("tooltip", "Kullanıcı adınız benzersiz olmalıdır ve en az 3 karakter içermelidir").
	MinLength(3).
	Unique("users", "username")
```

**Frontend Kullanımı:**
```tsx
<TextInput
  name="username"
  label="Kullanıcı Adı"
  value={username}
  onChange={setUsername}
  tooltip="Kullanıcı adınız benzersiz olmalıdır ve en az 3 karakter içermelidir"
  required
/>
```

### Desteklenen Komponentler

Tooltip desteği olan tüm form field komponentleri:

| Komponent | Tooltip Desteği | Kullanım Yeri |
|-----------|----------------|---------------|
| **TextInput** | ✅ | Form, Index, Detail |
| **TelInput** | ✅ | Form, Index, Detail |
| **DateField** | ✅ | Form, Index, Detail |
| **DateTimeField** | ✅ | Form, Index, Detail |
| **TimeField** | ✅ | Form, Index, Detail |
| **CheckboxField** | ✅ | Form, Index, Detail |
| **RadioGroupField** | ✅ | Form, Index, Detail |
| **NumberInput** | ✅ | Form, Index, Detail |
| **EmailInput** | ✅ | Form, Index, Detail |
| **PasswordInput** | ✅ | Form, Index, Detail |
| **URLInput** | ✅ | Form, Index, Detail |
| **TextareaField** | ✅ | Form, Index, Detail |
| **SelectField** | ✅ | Form, Index, Detail |
| **SwitchField** | ✅ | Form, Index, Detail |

### Örnekler

**Telefon Numarası ile Tooltip:**
```go
fields.Tel("Telefon", "phone").
	OnForm().
	Required().
	WithProps("tooltip", "Türkiye cep telefonu formatında giriniz (5XX XXX XX XX)")
```

```tsx
<TelInput
  name="phone"
  label="Telefon Numarası"
  value={phone}
  onChange={setPhone}
  mask="(599) 999 99 99"
  tooltip="Türkiye cep telefonu formatında giriniz (5XX XXX XX XX)"
  required
/>
```

**Tarih ile Tooltip:**
```go
fields.Date("Doğum Tarihi", "birth_date").
	OnForm().
	Required().
	WithProps("tooltip", "18 yaşından büyük olmalısınız")
```

```tsx
<DateField
  name="birth_date"
  label="Doğum Tarihi"
  value={birthDate}
  onChange={setBirthDate}
  tooltip="18 yaşından büyük olmalısınız"
  useNative
  required
/>
```

**Checkbox ile Tooltip:**
```go
fields.Checkbox("Kullanım Koşulları", "terms_accepted").
	OnForm().
	Required().
	WithProps("tooltip", "Devam etmek için kullanım koşullarını okumalı ve kabul etmelisiniz")
```

```tsx
<CheckboxField
  name="terms"
  label="Kullanım koşullarını kabul ediyorum"
  checked={terms}
  onCheckedChange={setTerms}
  tooltip="Devam etmek için kullanım koşullarını okumalı ve kabul etmelisiniz"
  required
/>
```

### Görünüm

Tooltip, label'ın yanında bir info ikonu (ℹ️) olarak gösterilir:

```
[Label] ℹ️
```

Kullanıcı info ikonunun üzerine geldiğinde (hover), tooltip içeriği bir popover içinde gösterilir.

### Best Practices

1. **Kısa ve Öz**: Tooltip metni kısa ve anlaşılır olmalıdır (maksimum 1-2 cümle)
2. **Ek Bilgi**: Tooltip, label'da yer almayan ek bilgi sağlamalıdır
3. **Yönlendirici**: Kullanıcıya ne yapması gerektiğini açıkça belirtmelidir
4. **Tutarlı**: Benzer alanlar için benzer tooltip formatı kullanılmalıdır

**İyi Örnekler:**
```tsx
tooltip="Kullanıcı adınız benzersiz olmalıdır ve en az 3 karakter içermelidir"
tooltip="Türkiye cep telefonu formatında giriniz (5XX XXX XX XX)"
tooltip="18 yaşından büyük olmalısınız"
```

**Kötü Örnekler:**
```tsx
tooltip="Kullanıcı adı" // Çok kısa, ek bilgi yok
tooltip="Bu alan kullanıcı adınızı girmeniz için kullanılır. Kullanıcı adınız benzersiz olmalıdır..." // Çok uzun
tooltip="Girin" // Belirsiz
```

## InputGroup Addon Desteği

Form field'larında `shadcn/ui` `InputGroup` pattern'i ile alanın başına/sonuna bileşen veya metin ekleyebilirsiniz.

Referans: [shadcn Input Group](https://ui.shadcn.com/docs/components/base/input-group.md)

### Ne İşe Yarar?

- Para birimi, birim, protokol gibi sabit ön/son ekler (`₺`, `%`, `https://`)
- Inline aksiyonlar (örn. şifre göster/gizle) ile aynı hizada kullanım
- Formda tutarlı input-group görünümü

### Desteklenen Props (Backend `WithProps`)

Önerilen anahtarlar:
- Baş addon: `startAddon`
- Son addon: `endAddon`

Uyumluluk için desteklenen alias'lar:
- Baş addon alias: `start_component`, `prefix`, `prepend`
- Son addon alias: `end_component`, `suffix`, `append`

> Not: Alias'lar eski/karma kullanımlar için desteklenir. Yeni tanımlarda `startAddon` / `endAddon` kullanın.

### Backend Örnekleri

```go
// Basit para birimi + birim
fields.Number("Fiyat", "price").
	OnForm().
	WithProps("startAddon", "₺").
	WithProps("endAddon", "/ay")
```

```go
// URL alanında protokol sabitleme
fields.URL("Web Site", "website").
	OnForm().
	WithProps("prefix", "https://")
```

```go
// Relationship/combobox alanlarında da aynı props çalışır
fields.BelongsTo("Kategori", "category_id", "categories").
	OnForm().
	WithProps("startAddon", "🔎")
```

### Kapsam

Addon çözümleyici (`resolveFieldInputAddons`) tüm form field component'lerinde uygulanır.

- Input/textarea/select/combobox tabanlı field'larda doğrudan `InputGroup` render edilir.
- Relationship chips, file, tabs, dialog gibi kompleks field'larda container-level addon uygulanır.

Uygulama dosyaları:
- `/Users/ferdiunal/Web/panel.go/web/src/components/fields/form/input-group-addon.tsx`
- `/Users/ferdiunal/Web/panel.go/web/src/components/fields/form/input-group-addon-utils.ts`

## Sık Hata Kontrolü (Field Odaklı)

- Formda alan görünmüyor: `OnForm()` / `OnlyOnForm()` / `HideOnCreate()` / `HideOnUpdate()` kombinasyonlarını kontrol edin.
- Listede alan görünmüyor: `OnList()` ve `HideOnList()` çakışmalarını kontrol edin.
- Değer yazılıyor ama okunmuyor: field `Key` değeri ile model alanı uyuşmuyor olabilir.
- `select` boş geliyor: `options` map'i yanlış formatta olabilir veya value tipiyle uyuşmuyor olabilir.
- Dosya yükleme başarısız: storage path/yetki ve `StoreAs` callback dönüşünü kontrol edin.
- Validasyon beklenmiyor: `Required()` ve backend doğrulama katmanının birlikte çalıştığını doğrulayın.

## Sonraki Adımlar

- [İlişkiler (Relationships)](Relationships) - Tablo ilişkileri
- [Yetkilendirme](Authorization) - Policy ve erişim kontrolü
- [API Referansı](API-Reference) - Tüm metodlar
