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

## Gelişmiş Özellikler

### Display Callbacks ve Formatlama

Alan değerlerinin nasıl görüntüleneceğini özelleştirin.

```go
// Display callback ile özel formatlama
fields.Text("Fiyat", "price").
	OnList().
	Display(func(value interface{}) string {
		if price, ok := value.(float64); ok {
			return fmt.Sprintf("₺%.2f", price)
		}
		return ""
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
			Display(func(value interface{}) string {
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

## Sonraki Adımlar

- [İlişkiler Rehberi](./Relationships.md) - Tablo ilişkileri
- [Politika Rehberi](./Authorization.md) - Yetkilendirme
- [API Referansı](./API-Reference.md) - Tüm metodlar
