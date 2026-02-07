# İlişkiler (Relationships)

İlişkiler, veritabanı tablolarındaki ilişkileri Go Panel API'de temsil etmenin yoludur. BelongsTo, HasMany, HasOne, BelongsToMany ve MorphTo gibi ilişki türlerini destekler.

## Genel Bakış

Relationship fields, ilişkili verileri yüklemek, göstermek, aramak, filtrelemek ve sıralamak için fluent API sağlar.

### Slug-Based Yaklaşım (Geleneksel)

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

### Resource-Based Yaklaşım (Önerilen)

Resource instance kullanarak tip güvenli ilişki tanımlama:

```go
// BelongsTo: Post -> Author
field := NewBelongsToResource("Author", "author_id", blog.NewAuthorResource()).
    DisplayUsing("name").
    WithSearchableColumns("name", "email").
    WithEagerLoad()

// HasMany: Author -> Posts
field := NewHasManyResource("Posts", "posts", blog.NewPostResource()).
    ForeignKey("author_id").
    WithEagerLoad()

// HasOne: User -> Profile
field := NewHasOneResource("Profile", "profile", blog.NewProfileResource()).
    ForeignKey("user_id")

// BelongsToMany: User -> Roles
field := NewBelongsToManyResource("Roles", "roles", blog.NewRoleResource()).
    PivotTable("role_user").
    ForeignKey("user_id").
    RelatedKey("role_id")
```

**Resource-Based Avantajları:**
- ✅ Tip güvenliği (derleme zamanı kontrolü)
- ✅ Refactoring desteği (resource adı değişirse otomatik güncellenir)
- ✅ IDE desteği (autocomplete, go-to-definition)
- ✅ Tablo adı otomatik alınır (`resource.Slug()`)
- ✅ Backward compatible (eski slug-based yöntem hala çalışır)

**Detaylı bilgi için:** [Resource-Based İlişkiler Dokümantasyonu](../.docs/RESOURCE_BASED_RELATIONSHIPS.md)

## İlişki Türleri

### BelongsTo

Inverse one-to-one veya one-to-many relationship. Bir model'in başka bir model'e ait olduğunu belirtir. Örneğin, bir Post bir Author'a aittir.

**Temel Kullanım:**

**Slug-Based (Geleneksel):**
```go
field := NewBelongsTo("Author", "author_id", "authors")
```

**Resource-Based (Önerilen):**
```go
field := NewBelongsToResource("Author", "author_id", blog.NewAuthorResource())
```

Resource-based yaklaşımda, tablo adı (`authors`) otomatik olarak resource'dan alınır (`blog.NewAuthorResource().Slug()`). Bu sayede:
- ✅ Tip güvenliği sağlanır
- ✅ Refactoring desteği artar
- ✅ IDE autocomplete çalışır

**Metodlar:**

#### DisplayUsing(key string)
İlişkili kayıtta hangi field'ın label olarak gösterileceğini belirler. Default olarak "name" field'ı kullanılır.

```go
field := NewBelongsTo("Author", "author_id", "authors").
    DisplayUsing("email")  // Author'un email'i gösterilir
```

#### WithSearchableColumns(columns ...string)
İlişkili kayıtlarda arama yapılabilecek sütunları belirler. Bu sütunlar combobox'ta arama yaparken kullanılır.

```go
field := NewBelongsTo("Author", "author_id", "authors").
    WithSearchableColumns("name", "email", "username")
```

#### AutoOptions(displayField string)
Form elemanları (Combobox/Select) için seçenekleri veritabanından otomatik olarak yükler. Manuel olarak `Options` callback'i tanımlamaya gerek kalmaz.

```go
// Tüm author'ları getirir ve 'name' field'ını gösterir
field := NewBelongsTo("Author", "author_id", "authors").
    AutoOptions("name")
```

**Önemli:** `AutoOptions` kullanıldığında, backend otomatik olarak ilgili tablodan tüm kayıtları çeker ve belirtilen field'ı label olarak kullanır.

#### Query(callback func(*Query) *Query)
İlişkili kayıtları filtrelemek, sıralamak veya sınırlamak için özel query tanımlar.

```go
field := NewBelongsTo("Author", "author_id", "authors").
    Query(func(q *Query) *Query {
        return q.
            Where("status", "=", "active").
            OrderBy("name", "ASC")
    })
```

#### WithEagerLoad() / WithLazyLoad()
Yükleme stratejisini belirler.

```go
// Eager loading: N+1 query problemini önler (önerilen)
field := NewBelongsTo("Author", "author_id", "authors").
    WithEagerLoad()

// Lazy loading: İhtiyaç anında yükler
field := NewBelongsTo("Author", "author_id", "authors").
    WithLazyLoad()
```

**Database Yapısı:**

BelongsTo relationship için database'de foreign key sütunu gereklidir:

```go
type Post struct {
    ID       uint   `json:"id"`
    Title    string `json:"title"`
    AuthorID uint   `json:"authorId"`  // Foreign key
    Author   *Author `json:"author"`    // İlişki (opsiyonel, eager loading için)
}

type Author struct {
    ID    uint   `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

**Frontend Davranışı:**

BelongsTo field frontend'de bir combobox/select olarak gösterilir. Kullanıcı listeden bir kayıt seçer.

Backend'den gelen data formatı:
```json
{
  "author": {
    "data": 1,  // Author ID
    "props": {
      "options": {
        "1": "John Doe",
        "2": "Jane Smith"
      },
      "related_resource": "authors"
    }
  }
}
```

**Kullanım Senaryoları:**

1. **Post -> Author**: Bir yazı bir yazara aittir
2. **Comment -> User**: Bir yorum bir kullanıcıya aittir
3. **Order -> Customer**: Bir sipariş bir müşteriye aittir
4. **Product -> Category**: Bir ürün bir kategoriye aittir

**Örnek: Post ve Author İlişkisi**

**Slug-Based:**
```go
// Post Resource
type PostResource struct {
    resource.OptimizedBase
}

func (r *PostResource) Fields() []core.Element {
    return []core.Element{
        fields.ID("ID").Sortable(),
        fields.Text("Title", "title").Required(),
        fields.Textarea("Content", "content"),

        // BelongsTo: Post -> Author (slug-based)
        fields.NewBelongsTo("Author", "author_id", "authors").
            DisplayUsing("name").
            WithSearchableColumns("name", "email").
            AutoOptions("name").
            WithEagerLoad().
            Required(),

        fields.DateTime("Created At", "createdAt").ReadOnly(),
    }
}
```

**Resource-Based (Önerilen):**
```go
// Post Resource
type PostResource struct {
    resource.OptimizedBase
}

func (r *PostResource) Fields() []core.Element {
    return []core.Element{
        fields.ID("ID").Sortable(),
        fields.Text("Title", "title").Required(),
        fields.Textarea("Content", "content"),

        // BelongsTo: Post -> Author (resource-based)
        fields.NewBelongsToResource("Author", "author_id", blog.NewAuthorResource()).
            DisplayUsing("name").
            WithSearchableColumns("name", "email").
            AutoOptions("name").
            WithEagerLoad().
            Required(),

        fields.DateTime("Created At", "createdAt").ReadOnly(),
    }
}
```

**Best Practices:**

1. **Eager Loading**: N+1 query problemini önlemek için eager loading kullan
2. **Searchable Columns**: Indexed sütunları searchable olarak belirle
3. **AutoOptions**: Manuel options tanımlamak yerine AutoOptions kullan
4. **Display Field**: Kullanıcı dostu bir field seç (name, title, email vb.)
5. **Required**: Zorunlu ilişkiler için `Required()` kullan

**Sorun Giderme:**

**Problem:** Options boş geliyor
**Çözüm:** `AutoOptions` veya manuel `Options` callback'i tanımla

**Problem:** N+1 query problemi
**Çözüm:** `WithEagerLoad()` kullan

**Problem:** Arama çalışmıyor
**Çözüm:** `WithSearchableColumns` ile aranabilir sütunları belirle

### HasMany

One-to-many relationship. Örneğin, bir Author birçok Post'a sahiptir.

**Temel Kullanım:**

**Slug-Based:**
```go
field := NewHasMany("Posts", "posts", "posts")
```

**Resource-Based (Önerilen):**
```go
field := NewHasManyResource("Posts", "posts", blog.NewPostResource())
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

One-to-one relationship. Bir model'in başka bir model'e sahip olduğunu belirtir. Örneğin, bir User bir Profile'a sahiptir.

**Temel Kullanım:**

**Slug-Based:**
```go
field := NewHasOne("Profile", "profile", "profiles")
```

**Resource-Based (Önerilen):**
```go
field := NewHasOneResource("Profile", "profile", blog.NewProfileResource())
```

**Metodlar:**

#### ForeignKey(key string)
İlişkili tablodaki foreign key sütununu belirler. Default olarak `{parent_model}_id` kullanılır.

```go
field := NewHasOne("Profile", "profile", "profiles").
    ForeignKey("user_id")  // profiles.user_id
```

#### OwnerKey(key string)
Parent model'deki key sütununu belirler. Default olarak `id` kullanılır.

```go
field := NewHasOne("Profile", "profile", "profiles").
    OwnerKey("id")  // users.id
```

#### AutoOptions(displayField string)
Form elemanları için seçenekleri veritabanından otomatik olarak yükler. **HasOne için özel davranış:** Sadece "boşta olan" (henüz bir parent'a atanmamış) kayıtları getirir (`foreign_key IS NULL` filtresi uygular).

```go
// Sadece user_id'si boş olan profilleri getirir
field := NewHasOne("Profile", "profile", "profiles").
    ForeignKey("user_id").
    AutoOptions("bio")
```

**Önemli:** `AutoOptions` kullanıldığında, edit modunda mevcut ilişkili kayıt da listeye dahil edilir. Böylece kullanıcı mevcut ilişkiyi koruyabilir veya değiştirebilir.

#### Query(callback func(*Query) *Query)
İlişkili kayıtları filtrelemek veya sıralamak için özel query tanımlar.

```go
field := NewHasOne("Profile", "profile", "profiles").
    Query(func(q *Query) *Query {
        return q.Where("status", "=", "active")
    })
```

**Database Yapısı:**

HasOne relationship için ilişkili tabloda foreign key sütunu gereklidir:

```go
type User struct {
    ID      uint     `json:"id"`
    Name    string   `json:"name"`
    Profile *Profile `json:"profile"`  // İlişki (opsiyonel, eager loading için)
}

type Profile struct {
    ID     uint   `json:"id"`
    UserID *uint  `json:"userId"`  // Foreign key (nullable)
    Bio    string `json:"bio"`
    User   *User  `json:"user"`    // Reverse relationship (opsiyonel)
}
```

**Frontend Davranışı:**

HasOne field frontend'de bir combobox/select olarak gösterilir. Kullanıcı listeden bir kayıt seçer veya boş bırakır.

Backend'den gelen data formatı:
```json
{
  "profile": {
    "data": 1,  // Profile ID (veya null)
    "props": {
      "options": {
        "1": "Software Engineer",
        "2": "Product Manager"
      },
      "related_resource": "profiles",
      "foreign_key": "user_id"
    }
  }
}
```

**Kullanım Senaryoları:**

1. **User -> Profile**: Bir kullanıcının bir profili vardır
2. **Author -> Bio**: Bir yazarın bir biyografisi vardır
3. **Product -> Detail**: Bir ürünün detay bilgisi vardır
4. **Order -> Invoice**: Bir siparişin bir faturası vardır

**Örnek: User ve Profile İlişkisi**

**Slug-Based:**
```go
// User Resource
type UserResource struct {
    resource.OptimizedBase
}

func (r *UserResource) Fields() []core.Element {
    return []core.Element{
        fields.ID("ID").Sortable(),
        fields.Text("Name", "name").Required(),
        fields.Email("Email", "email").Required(),

        // HasOne: User -> Profile (slug-based)
        fields.NewHasOne("Profile", "profile", "profiles").
            ForeignKey("user_id").
            AutoOptions("bio"),

        fields.DateTime("Created At", "createdAt").ReadOnly(),
    }
}
```

**Resource-Based (Önerilen):**
```go
// User Resource
type UserResource struct {
    resource.OptimizedBase
}

func (r *UserResource) Fields() []core.Element {
    return []core.Element{
        fields.ID("ID").Sortable(),
        fields.Text("Name", "name").Required(),
        fields.Email("Email", "email").Required(),

        // HasOne: User -> Profile (resource-based)
        fields.NewHasOneResource("Profile", "profile", blog.NewProfileResource()).
            ForeignKey("user_id").
            AutoOptions("bio"),

        fields.DateTime("Created At", "createdAt").ReadOnly(),
    }
}
```

**Best Practices:**

1. **Nullable Foreign Key**: İlişkili tablodaki foreign key nullable olmalı
2. **AutoOptions**: Manuel options tanımlamak yerine AutoOptions kullan
3. **Unique Constraint**: Foreign key'e unique constraint ekle (bir profile sadece bir user'a ait olabilir)
4. **Display Field**: Anlamlı bir field seç (bio, title, description vb.)

**Sorun Giderme:**

**Problem:** Tüm kayıtlar listeleniyor (boşta olanlar değil)
**Çözüm:** `AutoOptions` kullan, otomatik olarak `foreign_key IS NULL` filtresi uygular

**Problem:** Edit modunda mevcut ilişki gösterilmiyor
**Çözüm:** `AutoOptions` kullanıldığında mevcut kayıt otomatik olarak listeye eklenir

**Problem:** Birden fazla kayıt aynı parent'a atanabiliyor
**Çözüm:** Foreign key'e unique constraint ekle

### BelongsToMany

Many-to-many relationship. İki model arasında çoktan çoğa ilişki kurar. Örneğin, bir User birçok Role'e sahiptir ve bir Role birçok User'a sahiptir.

**Temel Kullanım:**

**Slug-Based:**
```go
field := NewBelongsToMany("Roles", "role_user", "roles")
```

**Resource-Based (Önerilen):**
```go
field := NewBelongsToManyResource("Roles", "roles", blog.NewRoleResource())
```

Resource-based yaklaşımda, pivot tablo adı otomatik oluşturulur (alfabetik sıralama ile).

**Metodlar:**

#### PivotTable(name string)
Pivot (ara) table'ın adını belirler. Bu table iki model arasındaki ilişkiyi saklar.

```go
field := NewBelongsToMany("Roles", "role_user", "roles").
    PivotTable("role_user")
```

#### ForeignKey(key string)
Pivot table'daki parent model'in foreign key sütununu belirler. Default olarak `{parent_model}_id` kullanılır.

```go
field := NewBelongsToMany("Roles", "role_user", "roles").
    ForeignKey("user_id")  // role_user.user_id
```

#### RelatedKey(key string)
Pivot table'daki related model'in foreign key sütununu belirler. Default olarak `{related_model}_id` kullanılır.

```go
field := NewBelongsToMany("Roles", "role_user", "roles").
    RelatedKey("role_id")  // role_user.role_id
```

#### DisplayUsing(key string)
İlişkili kayıtlarda hangi field'ın label olarak gösterileceğini belirler. Default olarak "name" field'ı kullanılır.

```go
field := NewBelongsToMany("Roles", "role_user", "roles").
    DisplayUsing("title")  // Role'ün title'ı gösterilir
```

#### WithSearchableColumns(columns ...string)
İlişkili kayıtlarda arama yapılabilecek sütunları belirler.

```go
field := NewBelongsToMany("Roles", "role_user", "roles").
    WithSearchableColumns("name", "description")
```

#### AutoOptions(displayField string)
Form elemanları için seçenekleri veritabanından otomatik olarak yükler. Tüm ilişkili kayıtları getirir.

```go
// Tüm role'leri getirir ve 'name' field'ını gösterir
field := NewBelongsToMany("Roles", "role_user", "roles").
    AutoOptions("name")
```

#### Query(callback func(*Query) *Query)
İlişkili kayıtları filtrelemek veya sıralamak için özel query tanımlar.

```go
field := NewBelongsToMany("Roles", "role_user", "roles").
    Query(func(q *Query) *Query {
        return q.
            Where("status", "=", "active").
            OrderBy("name", "ASC")
    })
```

**Database Yapısı:**

BelongsToMany relationship için üç tablo gereklidir:

```go
// Parent model
type User struct {
    ID    uint   `json:"id"`
    Name  string `json:"name"`
    Roles []Role `json:"roles" gorm:"many2many:role_user;"`  // İlişki
}

// Related model
type Role struct {
    ID    uint   `json:"id"`
    Name  string `json:"name"`
    Users []User `json:"users" gorm:"many2many:role_user;"`  // Reverse relationship
}

// Pivot table (migration ile oluşturulur)
// CREATE TABLE role_user (
//     user_id INT NOT NULL,
//     role_id INT NOT NULL,
//     PRIMARY KEY (user_id, role_id),
//     FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
//     FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
// );
```

**Frontend Davranışı:**

BelongsToMany field frontend'de bir multi-select combobox olarak gösterilir. Kullanıcı birden fazla kayıt seçebilir.

Backend'den gelen data formatı:
```json
{
  "roles": {
    "data": [1, 2, 3],  // Seçili role ID'leri
    "props": {
      "options": {
        "1": "Admin",
        "2": "Editor",
        "3": "Viewer"
      },
      "related_resource": "roles",
      "pivot_table": "role_user"
    }
  }
}
```

**Kullanım Senaryoları:**

1. **User -> Roles**: Bir kullanıcının birden fazla rolü olabilir
2. **Post -> Tags**: Bir yazının birden fazla etiketi olabilir
3. **Product -> Categories**: Bir ürün birden fazla kategoride olabilir
4. **Student -> Courses**: Bir öğrenci birden fazla kursa kayıtlı olabilir

**Örnek: User ve Roles İlişkisi**

**Slug-Based:**
```go
// User Resource
type UserResource struct {
    resource.OptimizedBase
}

func (r *UserResource) Fields() []core.Element {
    return []core.Element{
        fields.ID("ID").Sortable(),
        fields.Text("Name", "name").Required(),
        fields.Email("Email", "email").Required(),

        // BelongsToMany: User -> Roles (slug-based)
        fields.NewBelongsToMany("Roles", "role_user", "roles").
            PivotTable("role_user").
            ForeignKey("user_id").
            RelatedKey("role_id").
            DisplayUsing("name").
            WithSearchableColumns("name", "description").
            AutoOptions("name"),

        fields.DateTime("Created At", "createdAt").ReadOnly(),
    }
}
```

**Resource-Based (Önerilen):**
```go
// User Resource
type UserResource struct {
    resource.OptimizedBase
}

func (r *UserResource) Fields() []core.Element {
    return []core.Element{
        fields.ID("ID").Sortable(),
        fields.Text("Name", "name").Required(),
        fields.Email("Email", "email").Required(),

        // BelongsToMany: User -> Roles (resource-based)
        fields.NewBelongsToManyResource("Roles", "roles", blog.NewRoleResource()).
            PivotTable("role_user").
            ForeignKey("user_id").
            RelatedKey("role_id").
            DisplayUsing("name").
            WithSearchableColumns("name", "description").
            AutoOptions("name"),

        fields.DateTime("Created At", "createdAt").ReadOnly(),
    }
}
```

**Pivot Table ile Ekstra Sütunlar:**

Pivot table'a ekstra sütunlar ekleyebilirsiniz (örn: created_at, expires_at):

```go
// Pivot table migration
// CREATE TABLE role_user (
//     user_id INT NOT NULL,
//     role_id INT NOT NULL,
//     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
//     expires_at TIMESTAMP NULL,
//     PRIMARY KEY (user_id, role_id)
// );

// GORM model
type User struct {
    Roles []Role `gorm:"many2many:role_user;"`
}
```

**Best Practices:**

1. **Pivot Table Naming**: `{model1}_{model2}` formatında adlandır (alfabetik sıra)
2. **Composite Primary Key**: Pivot table'da (foreign_key1, foreign_key2) composite primary key kullan
3. **Cascade Delete**: Foreign key'lere ON DELETE CASCADE ekle
4. **AutoOptions**: Manuel options tanımlamak yerine AutoOptions kullan
5. **Searchable Columns**: Arama için indexed sütunları kullan

**Sorun Giderme:**

**Problem:** İlişkiler kaydedilmiyor
**Çözüm:** Pivot table'ın doğru adlandırıldığından ve foreign key'lerin doğru olduğundan emin ol

**Problem:** Duplicate entry hatası
**Çözüm:** Pivot table'da composite primary key veya unique constraint olduğundan emin ol

**Problem:** Options boş geliyor
**Çözüm:** `AutoOptions` veya manuel `Options` callback'i tanımla

**Problem:** Cascade delete çalışmıyor
**Çözüm:** Foreign key constraint'lere ON DELETE CASCADE ekle

### MorphTo

Polymorphic (çok biçimli) relationship. Bir model'in birden fazla farklı model türüne ait olabilmesini sağlar. Örneğin, bir Comment hem Post'a hem de Video'ya ait olabilir.

**Temel Kullanım:**
```go
field := NewMorphTo("Commentable", "commentable").
    Types(map[string]string{
        "posts":  "posts",   // Database Type => Resource Slug
        "videos": "videos",
    })
```

**Metodlar:**

#### Types(types map[string]string)
Polymorphic relationship için type mapping'i belirler. Key olarak database'de saklanan type değeri, value olarak resource slug'ı kullanılır.

```go
field := NewMorphTo("Commentable", "commentable").
    Types(map[string]string{
        "posts":    "posts",     // commentable_type = "posts" -> posts resource
        "videos":   "videos",    // commentable_type = "videos" -> videos resource
        "comments": "comments",  // commentable_type = "comments" -> comments resource
    })
```

#### Displays(displays map[string]string)
Her type için hangi field'ın label olarak gösterileceğini belirler. Bu sayede frontend'de ilişkili kaydın hangi field'ı gösterileceği kontrol edilir.

```go
field := NewMorphTo("Commentable", "commentable").
    Types(map[string]string{
        "posts":  "posts",
        "videos": "videos",
    }).
    Displays(map[string]string{
        "posts":  "title",      // Post için title field'ı gösterilir
        "videos": "name",       // Video için name field'ı gösterilir
    })
```

**Önemli:** `Displays` metodu kullanıldığında, backend otomatik olarak ilişkili kaydın display field'ını yükler ve frontend'e gönderir. Bu sayede frontend'de gereksiz API istekleri atılmaz.

#### WithEagerLoad() / WithLazyLoad()
Yükleme stratejisini belirler.

```go
// Eager loading: İlişkili veriyi önceden yükle
field := NewMorphTo("Commentable", "commentable").
    WithEagerLoad()

// Lazy loading: İlişkili veriyi ihtiyaç anında yükle (default)
field := NewMorphTo("Commentable", "commentable").
    WithLazyLoad()
```

**Database Yapısı:**

MorphTo relationship için database'de iki sütun gereklidir:
- `{field_name}_type`: İlişkili model'in type'ı (örn: "posts", "videos")
- `{field_name}_id`: İlişkili model'in ID'si

```go
type Comment struct {
    ID              uint   `json:"id"`
    Content         string `json:"content"`
    CommentableType string `json:"commentableType"` // "posts" veya "videos"
    CommentableID   uint   `json:"commentableId"`   // İlişkili kaydın ID'si
}
```

**Frontend Davranışı:**

MorphTo field frontend'de iki dropdown olarak gösterilir:
1. **Type Dropdown**: Hangi model türüne ait olacağını seçer (örn: Post, Video)
2. **Resource Dropdown**: Seçilen type'a göre ilgili kayıtları listeler

Backend'den gelen data formatı:
```json
{
  "commentable": {
    "data": {
      "type": "posts",
      "id": 1,
      "morphToType": "posts",
      "morphToId": 1,
      "title": "Post Title"  // Display field (Displays metodunda belirtilmişse)
    },
    "props": {
      "types": [
        {"label": "Posts", "slug": "posts", "value": "posts"},
        {"label": "Videos", "slug": "videos", "value": "videos"}
      ],
      "displays": {
        "posts": "title",
        "videos": "name"
      }
    }
  }
}
```

**Kullanım Senaryoları:**

1. **Yorumlar Sistemi**: Bir Comment hem Post'a hem de Video'ya ait olabilir
2. **Beğeni Sistemi**: Bir Like hem Post'a hem de Comment'e ait olabilir
3. **Etiketleme Sistemi**: Bir Tag hem Post'a hem de Video'ya ait olabilir
4. **Bildirim Sistemi**: Bir Notification farklı model türlerine referans verebilir

**Örnek: Yorum Sistemi**

```go
// Comment Resource
type CommentResource struct {
    resource.OptimizedBase
}

func (r *CommentResource) Fields() []core.Element {
    return []core.Element{
        fields.ID("ID").Sortable(),
        fields.Text("Content", "content").Required(),

        // MorphTo: Comment -> Commentable (Post, Video, vb.)
        fields.NewMorphTo("Commentable", "commentable").
            Types(map[string]string{
                "posts":  "posts",
                "videos": "videos",
            }).
            Displays(map[string]string{
                "posts":  "title",
                "videos": "name",
            }).
            WithLazyLoad(),

        fields.DateTime("Created At", "createdAt").ReadOnly(),
    }
}
```

**Best Practices:**

1. **Type Mapping**: Type key'leri database'de saklanan değerlerle eşleşmeli
2. **Display Field**: Her type için uygun bir display field belirle
3. **Lazy Loading**: Performans için lazy loading kullan (default)
4. **Validation**: Type'ların geçerli olduğundan emin ol

**Sorun Giderme:**

**Problem:** Display field gösterilmiyor
**Çözüm:** `Displays` metodunu kullan ve backend'den display field'ın geldiğinden emin ol

**Problem:** Type dropdown boş
**Çözüm:** `Types` metodunda doğru mapping'i kontrol et

**Problem:** Resource dropdown boş
**Çözüm:** İlgili resource'un index endpoint'inin çalıştığından emin ol

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
