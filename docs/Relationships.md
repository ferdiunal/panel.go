# İlişkiler (Relationships)

İlişkiler, veritabanı tablolarındaki ilişkileri Go Panel API'de temsil etmenin yoludur. BelongsTo, HasMany, HasOne, BelongsToMany ve MorphTo gibi ilişki türlerini destekler.

## Genel Bakış

Relationship fields, ilişkili verileri yüklemek, göstermek, aramak, filtrelemek ve sıralamak için fluent API sağlar.

```go
// BelongsTo: Post -> Author
field := NewBelongsTo("Author", "user_id", "users").
    DisplayUsing("name").
    WithSearchableColumns("name", "email").
    WithEagerLoad()

// HasMany: Author -> Posts
field := NewHasMany("Posts", "posts", "posts").
    ForeignKey("author_id").
    WithEagerLoad()

// HasOne: User -> Profile
field := NewHasOne("Profile", "profile", "profiles").
    ForeignKey("user_id")

// BelongsToMany: User -> Roles
field := NewBelongsToMany("Roles", "role_user", "roles").
    PivotTable("role_user").
    ForeignKey("user_id").
    RelatedKey("role_id")

// MorphTo: Comment -> Commentable (Post, Video, vb.)
field := NewMorphTo("Commentable", "commentable").
    Types(map[string]string{
        "post":    "posts",
        "video":   "videos",
    })
```

## İlişki Türleri

### BelongsTo

Inverse one-to-one relationship. Örneğin, bir Post bir Author'a aittir.

```go
field := NewBelongsTo("Author", "user_id", "users")
```

**Metodlar:**
- `DisplayUsing(key string)` - Gösterilecek alanı belirle (default: "name")
- `WithSearchableColumns(columns ...string)` - Aranabilir alanları belirle
- `Query(callback func(*Query) *Query)` - Query'yi özelleştir
- `WithEagerLoad()` - Eager loading kullan
- `WithLazyLoad()` - Lazy loading kullan

**Örnek:**
```go
field := NewBelongsTo("Author", "user_id", "users").
    DisplayUsing("email").
    WithSearchableColumns("name", "email")
```

### HasMany

One-to-many relationship. Örneğin, bir Author birçok Post'a sahiptir.

```go
field := NewHasMany("Posts", "posts", "posts")
```

**Metodlar:**
- `ForeignKey(key string)` - Foreign key alanını belirle
- `OwnerKey(key string)` - Owner key alanını belirle
- `Query(callback func(*Query) *Query)` - Query'yi özelleştir
- `WithEagerLoad()` - Eager loading kullan
- `WithLazyLoad()` - Lazy loading kullan

**Örnek:**
```go
field := NewHasMany("Posts", "posts", "posts").
    ForeignKey("author_id").
    Query(func(q *Query) *Query {
        return q.OrderBy("created_at", "DESC")
    })
```

### HasOne

One-to-one relationship. Örneğin, bir User bir Profile'a sahiptir.

```go
field := NewHasOne("Profile", "profile", "profiles")
```

**Metodlar:**
- `ForeignKey(key string)` - Foreign key alanını belirle
- `OwnerKey(key string)` - Owner key alanını belirle
- `Query(callback func(*Query) *Query)` - Query'yi özelleştir

**Örnek:**
```go
field := NewHasOne("Profile", "profile", "profiles").
    ForeignKey("user_id")
```

### BelongsToMany

Many-to-many relationship. Örneğin, bir User birçok Role'a sahiptir.

```go
field := NewBelongsToMany("Roles", "role_user", "roles")
```

**Metodlar:**
- `PivotTable(name string)` - Pivot table adını belirle
- `ForeignKey(key string)` - Foreign key alanını belirle
- `RelatedKey(key string)` - Related key alanını belirle
- `Query(callback func(*Query) *Query)` - Query'yi özelleştir

**Örnek:**
```go
field := NewBelongsToMany("Roles", "role_user", "roles").
    PivotTable("role_user").
    ForeignKey("user_id").
    RelatedKey("role_id")
```

### MorphTo

Polymorphic relationship. Örneğin, bir Comment bir Post veya Video'ya ait olabilir.

```go
field := NewMorphTo("Commentable", "commentable").
    Types(map[string]string{
        "post":  "posts",
        "video": "videos",
    })
```

**Metodlar:**
- `Types(types map[string]string)` - Type mapping'i belirle

**Örnek:**
```go
field := NewMorphTo("Commentable", "commentable").
    Types(map[string]string{
        "post":    "posts",
        "comment": "comments",
    })
```

## Yükleme Stratejileri

### Eager Loading

İlişkili verileri önceden yükle. Default strateji.

```go
field := NewBelongsTo("Author", "user_id", "users").
    WithEagerLoad()
```

**Avantajlar:**
- N+1 query problemini çözer
- Tüm ilişkili verileri tek sorguda yükler
- Performans optimizasyonu

### Lazy Loading

İlişkili verileri ihtiyaç anında yükle.

```go
field := NewBelongsTo("Author", "user_id", "users").
    WithLazyLoad()
```

**Avantajlar:**
- Sadece gerekli verileri yükler
- Bellek kullanımını azaltır
- Dinamik yükleme

## Query Özelleştirmesi

Query callback'leri kullanarak ilişkili verileri özelleştir.

```go
field := NewHasMany("Posts", "posts", "posts").
    Query(func(q *Query) *Query {
        return q.
            Where("status", "=", "published").
            OrderBy("created_at", "DESC").
            Limit(10)
    })
```

**Mevcut Metodlar:**
- `Where(column, operator, value)` - WHERE clause ekle
- `WhereIn(column, values)` - WHERE IN clause ekle
- `OrderBy(column, direction)` - ORDER BY ekle
- `Limit(count)` - LIMIT ekle
- `Offset(count)` - OFFSET ekle

## Görüntüleme Özelleştirmesi

İlişkili verilerin nasıl gösterileceğini özelleştir.

```go
// BelongsTo: Hangi alanı göster
field := NewBelongsTo("Author", "user_id", "users").
    DisplayUsing("email")

// HasMany: Sayı veya liste
field := NewHasMany("Posts", "posts", "posts")

// HasOne: İlişkili kayıt veya boş durum
field := NewHasOne("Profile", "profile", "profiles")

// BelongsToMany: İlişkili kayıtların listesi
field := NewBelongsToMany("Roles", "role_user", "roles")

// MorphTo: Kayıt ve type göstergesi
field := NewMorphTo("Commentable", "commentable").
    Types(map[string]string{
        "post":  "posts",
        "video": "videos",
    })
```

## Otomatik Seçenekler (AutoOptions)

`HasOne` ve `BelongsTo` ilişkilerinde form elemanları (Combobox/Select) için seçenekleri veritabanından otomatik olarak yüklemek için `AutoOptions` metodunu kullanabilirsiniz. Bu özellik, geliştiricinin manuel olarak veritabanı sorgusu yazmasını ve `Options` callback'i tanımlamasını gereksiz kılar.

### HasOne AutoOptions

`HasOne` ilişkisinde, genellikle "boşta olan" (henüz bir parent'a atanmamış) kayıtların listelenmesi istenir. `AutoOptions` bunu otomatik halleder (`foreign_key IS NULL` filtresi uygular).

```go
// Author -> Profile (HasOne)
// 'profiles' tablosundan, 'author_id'si boş olan kayıtları getirir.
// Listede 'bio' alanını gösterir.
fields.NewHasOne("Profile", "profile", "profiles").
    AutoOptions("bio")
```

### BelongsTo AutoOptions

`BelongsTo` ilişkisinde, genellikle tüm olası parent kayıtların listelenmesi istenir. `AutoOptions` tüm kayıtları getirir.

```go
// Post -> Author (BelongsTo)
// 'authors' tablosundan tüm yazarları getirir.
// Listede 'name' alanını gösterir.
fields.NewBelongsTo("Author", "author_id", "authors").
    AutoOptions("name")
```

**Not:** `AutoOptions` kullanıldığında `RelatedResourceSlug` parametresinin (3. parametre) veritabanı tablosu adıyla eşleşmesi veya doğru yapılandırılması gerekir. Ayrıca `HasOne` için `ForeignKeyColumn` doğru ayarlanmalıdır.

## Arama

İlişkili verilerde arama yap.

```go
field := NewBelongsTo("Author", "user_id", "users").
    WithSearchableColumns("name", "email")
```

**Özellikler:**
- Case-insensitive arama
- Birden fazla alanda arama
- Tam metin araması

## Filtreleme

İlişkili verileri filtrele.

```go
field := NewHasMany("Posts", "posts", "posts").
    Query(func(q *Query) *Query {
        return q.Where("status", "=", "published")
    })
```

## Sıralama

İlişkili verileri sırala.

```go
field := NewHasMany("Posts", "posts", "posts").
    Query(func(q *Query) *Query {
        return q.OrderBy("created_at", "DESC")
    })
```

**Yönler:**
- `ASC` - Artan sıra
- `DESC` - Azalan sıra

## Sayfalandırma

İlişkili verileri sayfalandır.

```go
field := NewHasMany("Posts", "posts", "posts").
    Query(func(q *Query) *Query {
        return q.Limit(10).Offset(0)
    })
```

## Kısıtlamalar

İlişkili verilere kısıtlamalar uygula.

```go
field := NewHasMany("Posts", "posts", "posts").
    Query(func(q *Query) *Query {
        return q.
            Where("status", "=", "published").
            WhereIn("category_id", []int{1, 2, 3}).
            Limit(10).
            Offset(0)
    })
```

## Sayma

İlişkili verilerin sayısını al.

```go
// BelongsTo: 0 veya 1
count := field.Count(data)

// HasMany: İlişkili kayıt sayısı
count := field.Count(data)

// BelongsToMany: Pivot table girdileri
count := field.Count(data)
```

## Varlık Kontrolü

İlişkili verilerin varlığını kontrol et.

```go
// Exists: İlişkili veri var mı?
exists := field.Exists(data)

// DoesntExist: İlişkili veri yok mu?
doesntExist := field.DoesntExist(data)
```

## Doğrulama

İlişkili verileri doğrula.

```go
// Zorunlu ilişki
field := NewBelongsTo("Author", "user_id", "users").
    Required()

// İsteğe bağlı ilişki (default)
field := NewBelongsTo("Author", "user_id", "users")
```

**Doğrulama Kuralları:**
- BelongsTo: İlişkili kayıt var mı?
- HasMany: Foreign key referansları geçerli mi?
- HasOne: En fazla bir ilişkili kayıt var mı?
- BelongsToMany: Pivot table girdileri geçerli mi?
- MorphTo: Morph type kayıtlı mı?

## JSON Serileştirmesi

İlişkili verileri JSON'a dönüştür.

```go
// Relationship'i serialize et
jsonData := field.Serialize(data)

// JSON string'e dönüştür
jsonStr := field.ToJSON(data)
```

**Çıktı:**
```json
{
  "type": "belongsTo",
  "name": "author",
  "value": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com"
  }
}
```

## En İyi Uygulamalar

### 1. Eager Loading Kullan
N+1 query problemini önlemek için eager loading kullan.

```go
// ✓ İyi
field := NewBelongsTo("Author", "user_id", "users").
    WithEagerLoad()

// ✗ Kaçın
field := NewBelongsTo("Author", "user_id", "users").
    WithLazyLoad()
```

### 2. Aranabilir Alanları Belirle
Arama performansını artırmak için indexed alanları kullan.

```go
// ✓ İyi
field := NewBelongsTo("Author", "user_id", "users").
    WithSearchableColumns("name", "email")

// ✗ Kaçın
field := NewBelongsTo("Author", "user_id", "users").
    WithSearchableColumns("bio", "description")
```

### 3. Query Özelleştirmesi Kullan
Gereksiz verileri yüklemekten kaçın.

```go
// ✓ İyi
field := NewHasMany("Posts", "posts", "posts").
    Query(func(q *Query) *Query {
        return q.Where("status", "=", "published")
    })

// ✗ Kaçın
field := NewHasMany("Posts", "posts", "posts")
```

### 4. Sayfalandırma Kullan
Büyük koleksiyonlar için sayfalandırma kullan.

```go
// ✓ İyi
field := NewHasMany("Posts", "posts", "posts").
    Query(func(q *Query) *Query {
        return q.Limit(10).Offset(0)
    })

// ✗ Kaçın
field := NewHasMany("Posts", "posts", "posts")
```

### 5. Doğrulama Kullan
İlişkili verilerin varlığını doğrula.

```go
// ✓ İyi
field := NewBelongsTo("Author", "user_id", "users").
    Required()

// ✗ Kaçın
field := NewBelongsTo("Author", "user_id", "users")
```

## Örnekler

### Örnek 1: Blog Yazısı ve Yazar

```go
// Post resource
type PostResource struct {
    // ... diğer alanlar
}

func (r *PostResource) Fields() []Element {
    return []Element{
        Text("Başlık", "title"),
        Textarea("İçerik", "content"),
        NewBelongsTo("Yazar", "user_id", "users").
            DisplayUsing("name").
            WithSearchableColumns("name", "email").
            WithEagerLoad(),
    }
}
```

### Örnek 2: Yazar ve Yazıları

```go
// Author resource
type AuthorResource struct {
    // ... diğer alanlar
}

func (r *AuthorResource) Fields() []Element {
    return []Element{
        Text("İsim", "name"),
        Email("E-posta", "email"),
        NewHasMany("Yazılar", "posts", "posts").
            ForeignKey("author_id").
            Query(func(q *Query) *Query {
                return q.OrderBy("created_at", "DESC")
            }).
            WithEagerLoad(),
    }
}
```

### Örnek 3: Kullanıcı ve Profili

```go
// User resource
type UserResource struct {
    // ... diğer alanlar
}

func (r *UserResource) Fields() []Element {
    return []Element{
        Text("İsim", "name"),
        Email("E-posta", "email"),
        NewHasOne("Profil", "profile", "profiles").
            ForeignKey("user_id"),
    }
}
```

### Örnek 4: Kullanıcı ve Rolleri

```go
// User resource
type UserResource struct {
    // ... diğer alanlar
}

func (r *UserResource) Fields() []Element {
    return []Element{
        Text("İsim", "name"),
        Email("E-posta", "email"),
        NewBelongsToMany("Roller", "role_user", "roles").
            PivotTable("role_user").
            ForeignKey("user_id").
            RelatedKey("role_id"),
    }
}
```

### Örnek 5: Polymorphic Yorumlar

```go
// Comment resource
type CommentResource struct {
    // ... diğer alanlar
}

func (r *CommentResource) Fields() []Element {
    return []Element{
        Textarea("İçerik", "content"),
        NewMorphTo("Yorumlanabilir", "commentable").
            Types(map[string]string{
                "post":  "posts",
                "video": "videos",
            }),
    }
}
```

## Sorun Giderme

### N+1 Query Problemi

**Problem:** Her ilişkili kayıt için ayrı query çalışıyor.

**Çözüm:** Eager loading kullan.

```go
// ✓ İyi
field := NewBelongsTo("Author", "user_id", "users").
    WithEagerLoad()
```

### İlişkili Veri Eksik

**Problem:** İlişkili veri yüklenmemiş.

**Çözüm:** Yükleme stratejisini kontrol et.

```go
// ✓ İyi
field := NewBelongsTo("Author", "user_id", "users").
    WithEagerLoad()
```

### Doğrulama Hataları

**Problem:** İlişkili veri doğrulaması başarısız.

**Çözüm:** Doğrulama kurallarını kontrol et.

```go
// ✓ İyi
field := NewBelongsTo("Author", "user_id", "users").
    Required()
```

### Performans Sorunları

**Problem:** Relationship queries yavaş.

**Çözüm:** Query özelleştirmesi ve sayfalandırma kullan.

```go
// ✓ İyi
field := NewHasMany("Posts", "posts", "posts").
    Query(func(q *Query) *Query {
        return q.
            Where("status", "=", "published").
            Limit(10).
            Offset(0)
    })
```

## Ayrıca Bkz.

- [Alanlar](./Fields.md) - Field system genel bakış
- [Kaynaklar](./Resources.md) - Resource tanımı
- [API Referansı](./API-Reference.md) - Tam API referansı
