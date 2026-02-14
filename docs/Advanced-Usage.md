# Gelişmiş Kullanım

Bu bölüm Panel SDK'nın ileri seviye özelliklerini ve tekniklerini kapsar. Temel kavramları zaten bildiğinizi varsayıyoruz.

## Özel Alanlar Oluşturma

Yerleşik alanlar yeterli değilse, kendi özel alanlarınızı oluşturabilirsiniz.

### Basit Özel Alan

```go
type ColorField struct {
    base.Base
    defaultValue string
}

func Color(label, attribute string) *ColorField {
    return &ColorField{
        Base: base.Base{
            Label:     label,
            Attribute: attribute,
        },
    }
}

func (f *ColorField) Resolve(data interface{}) interface{} {
    // Veriyi çöz ve renk değerini döndür
    return "#000000"
}

func (f *ColorField) Serialize(data interface{}) map[string]interface{} {
    return map[string]interface{}{
        "type":  "color",
        "value": f.Resolve(data),
    }
}
```

### Fluent API ile Özel Alan

```go
type RatingField struct {
    base.Base
    maxRating int
}

func Rating(label, attribute string) *RatingField {
    return &RatingField{
        Base: base.Base{
            Label:     label,
            Attribute: attribute,
        },
        maxRating: 5,
    }
}

func (f *RatingField) MaxRating(max int) *RatingField {
    f.maxRating = max
    return f
}

func (f *RatingField) Resolve(data interface{}) interface{} {
    // Rating değerini döndür
    return 4
}

func (f *RatingField) Serialize(data interface{}) map[string]interface{} {
    return map[string]interface{}{
        "type":       "rating",
        "value":      f.Resolve(data),
        "maxRating":  f.maxRating,
    }
}
```

## Özel Kaynaklar

Temel Resource sınıfını genişleterek özel kaynaklar oluşturun.

### Temel Özel Kaynak

```go
type BaseResource struct {
    resource.OptimizedBase
}

func (r *BaseResource) Group() string {
    return "Yönetim"
}

func (r *BaseResource) Icon() string {
    return "cog"
}

func (r *BaseResource) SortableBy() []string {
    return []string{"created_at", "updated_at"}
}
```

### Özel Kaynak Kullanımı

```go
type ProductResource struct {
    BaseResource
}

func (r *ProductResource) Model() interface{} {
    return &Product{}
}

func (r *ProductResource) Slug() string {
    return "products"
}

func (r *ProductResource) Title() string {
    return "Ürünler"
}

func (r *ProductResource) Fields() []fields.Element {
    return []fields.Element{
        fields.ID(),
        fields.Text("Adı", "name").Sortable().Required(),
        fields.Number("Fiyat", "price").Required(),
    }
}
```

## Middleware ve Hooks

Kaynakların yaşam döngüsüne müdahale edin.

### Ön İşleme (Before Hooks)

```go
type UserResource struct {
    resource.OptimizedBase
}

func (r *UserResource) BeforeCreate(ctx context.Context, data map[string]interface{}) error {
    // Şifreyi hash'le
    if password, ok := data["password"].(string); ok {
        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
        if err != nil {
            return err
        }
        data["password"] = string(hashedPassword)
    }
    return nil
}

func (r *UserResource) BeforeUpdate(ctx context.Context, id interface{}, data map[string]interface{}) error {
    // Güncelleme öncesi doğrulama
    if email, ok := data["email"].(string); ok {
        if !isValidEmail(email) {
            return fmt.Errorf("geçersiz e-posta: %s", email)
        }
    }
    return nil
}

func (r *UserResource) BeforeDelete(ctx context.Context, id interface{}) error {
    // Silme öncesi kontrol
    return nil
}
```

### Sonrası İşleme (After Hooks)

```go
func (r *UserResource) AfterCreate(ctx context.Context, model interface{}) error {
    // Yeni kullanıcı oluşturulduktan sonra
    user := model.(*User)
    
    // E-posta gönder
    sendWelcomeEmail(user.Email)
    
    // Log kaydı oluştur
    logAction("user_created", user.ID)
    
    return nil
}

func (r *UserResource) AfterUpdate(ctx context.Context, model interface{}) error {
    // Güncelleme sonrası
    user := model.(*User)
    logAction("user_updated", user.ID)
    return nil
}

func (r *UserResource) AfterDelete(ctx context.Context, id interface{}) error {
    // Silme sonrası
    logAction("user_deleted", id)
    return nil
}
```

## Özel Doğrulama

Alanlar için özel doğrulama kuralları tanımlayın.

### Doğrulama Kuralları

```go
type ProductResource struct {
    resource.OptimizedBase
}

func (r *ProductResource) Fields() []fields.Element {
    return []fields.Element{
        fields.ID(),
        fields.Text("Adı", "name").
            Required().
            Rules(
                validate.MinLength(3),
                validate.MaxLength(100),
                validate.Unique("products", "name"),
            ),
        fields.Number("Fiyat", "price").
            Required().
            Rules(
                validate.Min(0),
                validate.Max(1000000),
            ),
        fields.Number("Stok", "stock").
            Required().
            Rules(
                validate.Min(0),
                validate.Integer(),
            ),
    }
}
```

### Özel Doğrulama Fonksiyonu

```go
func ValidateProductCode(value interface{}) error {
    code, ok := value.(string)
    if !ok {
        return fmt.Errorf("ürün kodu string olmalı")
    }
    
    if len(code) < 5 {
        return fmt.Errorf("ürün kodu en az 5 karakter olmalı")
    }
    
    if !isValidProductCode(code) {
        return fmt.Errorf("geçersiz ürün kodu formatı")
    }
    
    return nil
}

// Kullanım
fields.Text("Ürün Kodu", "code").
    Required().
    Rules(ValidateProductCode)
```

## Özel Sorgu Optimizasyonu

Sorguları optimize etmek için özel yükleme stratejileri kullanın.

### Eager Loading Özelleştirmesi

```go
type PostResource struct {
    resource.OptimizedBase
}

func (r *PostResource) Fields() []fields.Element {
    return []fields.Element{
        fields.ID(),
        fields.Text("Başlık", "title"),
        fields.BelongsTo("Yazar", "user_id", "users").
            DisplayUsing("name").
            WithEagerLoad(),
        fields.HasMany("Yorumlar", "comments", "comments").
            ForeignKey("post_id").
            Query(func(q interface{}) interface{} {
                db, ok := q.(*gorm.DB)
                if !ok || db == nil {
                    return q
                }
                return db.
                    Where("approved = ?", true).
                    Order("created_at DESC").
                    Limit(5)
            }).
            WithEagerLoad(),
    }
}
```

### Lazy Loading Özelleştirmesi

```go
func (r *PostResource) Fields() []fields.Element {
    return []fields.Element{
        fields.ID(),
        fields.Text("Başlık", "title"),
        fields.HasMany("Tüm Yorumlar", "comments", "comments").
            ForeignKey("post_id").
            WithLazyLoad(), // İhtiyaç anında yükle
    }
}
```

## Özel Görünüm Mantığı

Alanların nasıl görüntüleneceğini özelleştirin.

### Koşullu Görünürlük

```go
type UserResource struct {
    resource.OptimizedBase
}

func (r *UserResource) Fields() []fields.Element {
    return []fields.Element{
        fields.ID(),
        fields.Text("İsim", "name"),
        fields.Email("E-posta", "email"),
        fields.Password("Şifre", "password").
            OnlyOnForms(), // Sadece form'da göster
        fields.Text("Rol", "role").
            OnlyOnDetail(), // Sadece detay sayfasında göster
        fields.Text("Oluşturma Tarihi", "created_at").
            OnlyOnIndex(), // Sadece liste sayfasında göster
    }
}
```

### Dinamik Görünürlük

```go
type PostResource struct {
    resource.OptimizedBase
}

func (r *PostResource) Fields() []fields.Element {
    return []fields.Element{
        fields.ID(),
        fields.Text("Başlık", "title"),
        fields.Textarea("İçerik", "content"),
        fields.Text("Yayın Durumu", "status").
            HiddenIf(func(data interface{}) bool {
                post := data.(*Post)
                return post.Status == "draft"
            }),
    }
}
```

## Özel Arama Mantığı

Arama davranışını özelleştirin.

### Gelişmiş Arama

```go
type ProductResource struct {
    resource.OptimizedBase
}

func (r *ProductResource) Fields() []fields.Element {
    return []fields.Element{
        fields.ID(),
        fields.Text("Adı", "name").
            Searchable().
            SearchableColumns("name", "description", "sku"),
        fields.Number("Fiyat", "price"),
        fields.BelongsTo("Kategori", "category_id", "categories").
            DisplayUsing("name").
            WithSearchableColumns("name"),
    }
}
```

### Tam Metin Araması

```go
func (r *ProductResource) Search(query string) interface{} {
    // Tam metin araması yapılandırması
    return map[string]interface{}{
        "type": "fulltext",
        "columns": []string{"name", "description"},
        "query": query,
    }
}
```

## Özel Sıralama

Sıralama davranışını özelleştirin.

### Çoklu Alan Sıralaması

```go
type OrderResource struct {
    resource.OptimizedBase
}

func (r *OrderResource) Fields() []fields.Element {
    return []fields.Element{
        fields.ID(),
        fields.Text("Sipariş Numarası", "order_number").Sortable(),
        fields.Date("Tarih", "created_at").Sortable(),
        fields.Number("Toplam", "total").Sortable(),
    }
}

func (r *OrderResource) SortableBy() []string {
    return []string{"created_at", "total", "order_number"}
}
```

## Özel Filtreleme

Filtreleme davranışını özelleştirin.

### Gelişmiş Filtreler

```go
type ProductResource struct {
    resource.OptimizedBase
}

func (r *ProductResource) Filters() []interface{} {
    return []interface{}{
        // Kategori filtresi
        &CategoryFilter{},
        
        // Fiyat aralığı filtresi
        &PriceRangeFilter{},
        
        // Stok durumu filtresi
        &StockStatusFilter{},
    }
}
```

## Özel Sayfalandırma

Sayfalandırma davranışını özelleştirin.

### Sayfa Boyutu Yapılandırması

```go
type ProductResource struct {
    resource.OptimizedBase
}

func (r *ProductResource) PerPage() int {
    return 25 // Varsayılan 15 yerine 25
}

func (r *ProductResource) PerPageOptions() []int {
    return []int{10, 25, 50, 100}
}
```

## Özel JSON Serileştirmesi

JSON çıktısını özelleştirin.

### Serileştirme Özelleştirmesi

```go
type UserResource struct {
    resource.OptimizedBase
}

func (r *UserResource) Serialize(model interface{}) map[string]interface{} {
    user := model.(*User)
    
    return map[string]interface{}{
        "id":    user.ID,
        "name":  user.Name,
        "email": user.Email,
        "role":  user.Role,
        "created_at": user.CreatedAt.Format("2006-01-02"),
        "updated_at": user.UpdatedAt.Format("2006-01-02"),
    }
}
```

## Özel Hata İşleme

Hata işlemeyi özelleştirin.

### Hata Yakalama ve İşleme

```go
type UserResource struct {
    resource.OptimizedBase
}

func (r *UserResource) HandleError(err error) error {
    if err == nil {
        return nil
    }
    
    // Veritabanı hatalarını özelleştir
    if strings.Contains(err.Error(), "duplicate key") {
        return fmt.Errorf("bu e-posta zaten kullanılıyor")
    }
    
    // Yetkilendirme hatalarını özelleştir
    if strings.Contains(err.Error(), "permission denied") {
        return fmt.Errorf("bu işlemi yapmak için yetkiniz yok")
    }
    
    return err
}
```

## Özel Bağlam (Context) Kullanımı

Bağlam bilgisini kullanarak özel mantık uygulayın.

### Kullanıcı Bilgisine Erişim

```go
type PostResource struct {
    resource.OptimizedBase
}

func (r *PostResource) BeforeCreate(ctx context.Context, data map[string]interface{}) error {
    // Bağlamdan kullanıcı bilgisini al
    user := ctx.Value("user").(*User)
    
    // Yazarı otomatik olarak ayarla
    data["author_id"] = user.ID
    
    return nil
}
```

### Kiracı Yalıtımı (Multi-Tenancy)

```go
func (r *PostResource) BeforeCreate(ctx context.Context, data map[string]interface{}) error {
    // Bağlamdan kiracı bilgisini al
    tenantID := ctx.Value("tenant_id").(string)
    
    // Kiracı ID'sini ayarla
    data["tenant_id"] = tenantID
    
    return nil
}
```

## Performans İpuçları

### 1. Eager Loading Kullan

```go
// ✓ İyi - N+1 problemini çözer
fields.BelongsTo("Author", "user_id", "users").WithEagerLoad()

// ✗ Kaçın - N+1 problemi
fields.BelongsTo("Author", "user_id", "users").WithLazyLoad()
```

### 2. Sayfalandırma Kullan

```go
// ✓ İyi - Bellek kullanımını azaltır
func (r *ProductResource) PerPage() int {
    return 25
}

// ✗ Kaçın - Tüm verileri yükler
func (r *ProductResource) PerPage() int {
    return 10000
}
```

### 3. Sorguları Optimize Et

```go
// ✓ İyi - Sadece gerekli alanları seç
fields.HasMany("Posts", "posts", "posts").
    Query(func(q interface{}) interface{} {
        db, ok := q.(*gorm.DB)
        if !ok || db == nil {
            return q
        }
        return db.Select("id", "title", "created_at")
    })

// ✗ Kaçın - Tüm alanları yükle
fields.HasMany("Posts", "posts", "posts")
```

### 4. İndeksler Kullan

```go
// Veritabanında sık sorgulanan alanlar için indeks oluştur
CREATE INDEX idx_posts_author_id ON posts(author_id);
CREATE INDEX idx_posts_created_at ON posts(created_at);
```

### 5. Caching Kullan

```go
type ProductResource struct {
    resource.OptimizedBase
}

func (r *ProductResource) Fields() []fields.Element {
    // Kategorileri cache'le
    categories := getCachedCategories()
    
    return []fields.Element{
        fields.ID(),
        fields.Text("Adı", "name"),
        fields.Select("Kategori", "category_id", categories),
    }
}
```

## Ayrıca Bkz.

- [Alanlar](./Fields.md) - Field system genel bakış
- [Kaynaklar](./Resources.md) - Resource tanımı
- [İlişkiler](./Relationships.md) - Relationship fields
- [API Referansı](./API-Reference.md) - Tam API referansı
