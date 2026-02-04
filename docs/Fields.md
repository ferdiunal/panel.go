# Alanlar (Fields) Rehberi

Alanlar, veritabanı sütunlarını yönetim panelinde nasıl göstereceğinizi ve işleneceğini tanımlar.

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

Tarih seçimi için kullanılır.

```go
(&fields.Schema{
	Key:   "published_at",
	Name:  "Yayınlanma Tarihi",
	View:  "date",
	Props: make(map[string]interface{}),
}).OnDetail().OnForm()
```

### Tarih-Saat Alanı (DateTime)

Tarih ve saat seçimi için kullanılır.

```go
(&fields.Schema{
	Key:   "created_at",
	Name:  "Oluşturulma Tarihi",
	View:  "datetime",
	Props: make(map[string]interface{}),
}).ReadOnly().OnList().OnDetail()
```

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
(&fields.Schema{
	Key:   "password",
	Name:  "Şifre",
	View:  "password",
	Props: make(map[string]interface{}),
}).OnForm().Required()
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

// Detayda gizle
field.HideOnDetail()

// Oluşturmada gizle
field.HideOnCreate()

// Güncellemede gizle
field.HideOnUpdate()

// Tüm yerlerde göster
field.OnList().OnDetail().OnForm()
```

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
(&fields.Schema{
	Key:   "email",
	Name:  "E-posta",
	View:  "email",
	Props: make(map[string]interface{}),
}).
	OnList().
	OnDetail().
	OnForm().
	Required().
	Searchable().
	HelpText("Geçerli bir e-posta adresi girin")
```

## İpuçları

1. **Performans**: Aranabilir alanlar için veritabanında indeks oluşturun
2. **Validasyon**: Hem frontend hem backend'de validasyon yapın
3. **Güvenlik**: Hassas alanları `ReadOnly()` yapın
4. **UX**: Yardım metni ekleyerek kullanıcıları yönlendirin
5. **Sıralama**: Sık kullanılan alanları sıralanabilir yapın

## Sonraki Adımlar

- [İlişkiler Rehberi](./Relationships.md) - Tablo ilişkileri
- [Politika Rehberi](./Authorization.md) - Yetkilendirme
- [API Referansı](./API-Reference.md) - Tüm metodlar
