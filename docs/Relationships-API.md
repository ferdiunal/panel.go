# İlişkiler API Referansı

Relationship fields için kapsamlı API referansı.

## RelationshipField Interface

```go
type RelationshipField interface {
    Element
    GetRelationshipType() string
    GetRelatedResource() string
    GetRelationshipName() string
    ResolveRelationship(ctx context.Context, data interface{}) (interface{}, error)
    ValidateRelationship(ctx context.Context, value interface{}) error
}
```

## BelongsTo

Inverse one-to-one relationship.

### Constructor

```go
func NewBelongsTo(name, attribute, relatedResource string) *BelongsTo
```

**Parametreler:**
- `name` (string): Field adı
- `attribute` (string): Attribute key (foreign key column)
- `relatedResource` (string): İlişkili resource slug

**Örnek:**
```go
field := NewBelongsTo("Author", "user_id", "users")
```

### Metodlar

#### DisplayUsing

```go
func (b *BelongsTo) DisplayUsing(key string) *BelongsTo
```

Gösterilecek alanı belirle.

**Parametreler:**
- `key` (string): Gösterilecek alan adı (default: "name")

**Örnek:**
```go
field.DisplayUsing("email")
```

#### WithSearchableColumns

```go
func (b *BelongsTo) WithSearchableColumns(columns ...string) *BelongsTo
```

Aranabilir alanları belirle.

**Parametreler:**
- `columns` ([]string): Aranabilir alan adları

**Örnek:**
```go
field.WithSearchableColumns("name", "email")
```

#### Query

```go
func (b *BelongsTo) Query(callback func(*Query) *Query) *BelongsTo
```

Query'yi özelleştir.

**Parametreler:**
- `callback` (func): Query özelleştirme callback'i

**Örnek:**
```go
field.Query(func(q *Query) *Query {
    return q.Where("status", "=", "active")
})
```

#### WithEagerLoad

```go
func (b *BelongsTo) WithEagerLoad() *BelongsTo
```

Eager loading stratejisini kullan.

**Örnek:**
```go
field.WithEagerLoad()
```

#### WithLazyLoad

```go
func (b *BelongsTo) WithLazyLoad() *BelongsTo
```

Lazy loading stratejisini kullan.

**Örnek:**
```go
field.WithLazyLoad()
```

#### Required

```go
func (b *BelongsTo) Required() *BelongsTo
```

İlişkiyi zorunlu yap.

**Örnek:**
```go
field.Required()
```

#### GetRelationshipType

```go
func (b *BelongsTo) GetRelationshipType() string
```

İlişki türünü döndür. Döner: `"belongsTo"`

#### GetRelatedResource

```go
func (b *BelongsTo) GetRelatedResource() string
```

İlişkili resource slug'ını döndür.

#### GetDisplayKey

```go
func (b *BelongsTo) GetDisplayKey() string
```

Gösterilecek alanı döndür.

#### GetSearchableColumns

```go
func (b *BelongsTo) GetSearchableColumns() []string
```

Aranabilir alanları döndür.

#### GetLoadingStrategy

```go
func (b *BelongsTo) GetLoadingStrategy() LoadingStrategy
```

Yükleme stratejisini döndür.

---

## HasMany

One-to-many relationship.

### Constructor

```go
func NewHasMany(name, attribute, relatedResource string) *HasMany
```

**Parametreler:**
- `name` (string): Field adı
- `attribute` (string): Attribute key
- `relatedResource` (string): İlişkili resource slug

**Örnek:**
```go
field := NewHasMany("Posts", "posts", "posts")
```

### Metodlar

#### ForeignKey

```go
func (h *HasMany) ForeignKey(key string) *HasMany
```

Foreign key alanını belirle.

**Parametreler:**
- `key` (string): Foreign key column adı

**Örnek:**
```go
field.ForeignKey("author_id")
```

#### OwnerKey

```go
func (h *HasMany) OwnerKey(key string) *HasMany
```

Owner key alanını belirle.

**Parametreler:**
- `key` (string): Owner key column adı

**Örnek:**
```go
field.OwnerKey("id")
```

#### Query

```go
func (h *HasMany) Query(callback func(*Query) *Query) *HasMany
```

Query'yi özelleştir.

**Örnek:**
```go
field.Query(func(q *Query) *Query {
    return q.OrderBy("created_at", "DESC")
})
```

#### WithEagerLoad

```go
func (h *HasMany) WithEagerLoad() *HasMany
```

Eager loading stratejisini kullan.

#### WithLazyLoad

```go
func (h *HasMany) WithLazyLoad() *HasMany
```

Lazy loading stratejisini kullan.

#### GetRelationshipType

```go
func (h *HasMany) GetRelationshipType() string
```

İlişki türünü döndür. Döner: `"hasMany"`

---

## HasOne

One-to-one relationship.

### Constructor

```go
func NewHasOne(name, attribute, relatedResource string) *HasOne
```

**Parametreler:**
- `name` (string): Field adı
- `attribute` (string): Attribute key
- `relatedResource` (string): İlişkili resource slug

**Örnek:**
```go
field := NewHasOne("Profile", "profile", "profiles")
```

### Metodlar

#### ForeignKey

```go
func (h *HasOne) ForeignKey(key string) *HasOne
```

Foreign key alanını belirle.

#### OwnerKey

```go
func (h *HasOne) OwnerKey(key string) *HasOne
```

Owner key alanını belirle.

#### Query

```go
func (h *HasOne) Query(callback func(*Query) *Query) *HasOne
```

Query'yi özelleştir.

#### GetRelationshipType

```go
func (h *HasOne) GetRelationshipType() string
```

İlişki türünü döndür. Döner: `"hasOne"`

---

## BelongsToMany

Many-to-many relationship.

### Constructor

```go
func NewBelongsToMany(name, attribute, relatedResource string) *BelongsToMany
```

**Parametreler:**
- `name` (string): Field adı
- `attribute` (string): Attribute key
- `relatedResource` (string): İlişkili resource slug

**Örnek:**
```go
field := NewBelongsToMany("Roles", "role_user", "roles")
```

### Metodlar

#### PivotTable

```go
func (b *BelongsToMany) PivotTable(name string) *BelongsToMany
```

Pivot table adını belirle.

**Parametreler:**
- `name` (string): Pivot table adı

**Örnek:**
```go
field.PivotTable("role_user")
```

#### ForeignKey

```go
func (b *BelongsToMany) ForeignKey(key string) *BelongsToMany
```

Foreign key alanını belirle.

**Örnek:**
```go
field.ForeignKey("user_id")
```

#### RelatedKey

```go
func (b *BelongsToMany) RelatedKey(key string) *BelongsToMany
```

Related key alanını belirle.

**Örnek:**
```go
field.RelatedKey("role_id")
```

#### Query

```go
func (b *BelongsToMany) Query(callback func(*Query) *Query) *BelongsToMany
```

Query'yi özelleştir.

#### GetRelationshipType

```go
func (b *BelongsToMany) GetRelationshipType() string
```

İlişki türünü döndür. Döner: `"belongsToMany"`

---

## MorphTo

Polymorphic relationship.

### Constructor

```go
func NewMorphTo(name, attribute string) *MorphTo
```

**Parametreler:**
- `name` (string): Field adı
- `attribute` (string): Attribute key

**Örnek:**
```go
field := NewMorphTo("Commentable", "commentable")
```

### Metodlar

#### Types

```go
func (m *MorphTo) Types(types map[string]string) *MorphTo
```

Type mapping'i belirle.

**Parametreler:**
- `types` (map[string]string): Type name -> resource slug mapping

**Örnek:**
```go
field.Types(map[string]string{
    "post":  "posts",
    "video": "videos",
})
```

#### GetTypes

```go
func (m *MorphTo) GetTypes() map[string]string
```

Type mapping'i döndür.

#### GetRelationshipType

```go
func (m *MorphTo) GetRelationshipType() string
```

İlişki türünü döndür. Döner: `"morphTo"`

---

## RelationshipQuery

Query özelleştirmesi için kullanılan struct.

### Metodlar

#### Where

```go
func (q *RelationshipQuery) Where(column, operator string, value interface{}) *RelationshipQuery
```

WHERE clause ekle.

**Parametreler:**
- `column` (string): Column adı
- `operator` (string): Karşılaştırma operatörü (=, !=, >, <, >=, <=, LIKE, IN, vb.)
- `value` (interface{}): Karşılaştırma değeri

**Örnek:**
```go
query.Where("status", "=", "published")
```

#### WhereIn

```go
func (q *RelationshipQuery) WhereIn(column string, values []interface{}) *RelationshipQuery
```

WHERE IN clause ekle.

**Parametreler:**
- `column` (string): Column adı
- `values` ([]interface{}): Değer listesi

**Örnek:**
```go
query.WhereIn("category_id", []interface{}{1, 2, 3})
```

#### OrderBy

```go
func (q *RelationshipQuery) OrderBy(column, direction string) *RelationshipQuery
```

ORDER BY ekle.

**Parametreler:**
- `column` (string): Column adı
- `direction` (string): Sıralama yönü (ASC, DESC)

**Örnek:**
```go
query.OrderBy("created_at", "DESC")
```

#### Limit

```go
func (q *RelationshipQuery) Limit(count int) *RelationshipQuery
```

LIMIT ekle.

**Parametreler:**
- `count` (int): Limit sayısı

**Örnek:**
```go
query.Limit(10)
```

#### Offset

```go
func (q *RelationshipQuery) Offset(count int) *RelationshipQuery
```

OFFSET ekle.

**Parametreler:**
- `count` (int): Offset sayısı

**Örnek:**
```go
query.Offset(20)
```

#### GetWhereConditions

```go
func (q *RelationshipQuery) GetWhereConditions() []WhereCondition
```

WHERE conditions'ı döndür.

#### GetOrderByColumns

```go
func (q *RelationshipQuery) GetOrderByColumns() []OrderByColumn
```

ORDER BY columns'ı döndür.

#### GetLimit

```go
func (q *RelationshipQuery) GetLimit() int
```

LIMIT value'sini döndür.

#### GetOffset

```go
func (q *RelationshipQuery) GetOffset() int
```

OFFSET value'sini döndür.

---

## RelationshipValidator

Relationship doğrulaması için kullanılan struct.

### Metodlar

#### ValidateBelongsTo

```go
func (v *RelationshipValidator) ValidateBelongsTo(ctx context.Context, value interface{}, field *BelongsTo) error
```

BelongsTo relationship'i doğrula.

#### ValidateHasMany

```go
func (v *RelationshipValidator) ValidateHasMany(ctx context.Context, value interface{}, field *HasMany) error
```

HasMany relationship'i doğrula.

#### ValidateHasOne

```go
func (v *RelationshipValidator) ValidateHasOne(ctx context.Context, value interface{}, field *HasOne) error
```

HasOne relationship'i doğrula.

#### ValidateBelongsToMany

```go
func (v *RelationshipValidator) ValidateBelongsToMany(ctx context.Context, value interface{}, field *BelongsToMany) error
```

BelongsToMany relationship'i doğrula.

#### ValidateMorphTo

```go
func (v *RelationshipValidator) ValidateMorphTo(ctx context.Context, value interface{}, field *MorphTo) error
```

MorphTo relationship'i doğrula.

---

## RelationshipDisplay

Relationship görüntüleme özelleştirmesi için kullanılan struct.

### Metodlar

#### GetDisplayValue

```go
func (d *RelationshipDisplay) GetDisplayValue(data interface{}) (string, error)
```

Tek bir ilişkili kayıt için display value'sini döndür.

#### GetDisplayValues

```go
func (d *RelationshipDisplay) GetDisplayValues(data []interface{}) ([]string, error)
```

Birden fazla ilişkili kayıt için display value'lerini döndür.

---

## RelationshipSearch

Relationship araması için kullanılan struct.

### Metodlar

#### Search

```go
func (s *RelationshipSearch) Search(ctx context.Context, term string) ([]interface{}, error)
```

İlişkili kayıtlarda arama yap.

**Parametreler:**
- `ctx` (context.Context): Context
- `term` (string): Arama terimi

**Döner:**
- `[]interface{}`: Eşleşen kayıtlar
- `error`: Hata varsa

---

## RelationshipFilter

Relationship filtrelemesi için kullanılan struct.

### Metodlar

#### ApplyFilter

```go
func (f *RelationshipFilter) ApplyFilter(ctx context.Context, column, operator string, value interface{}) ([]interface{}, error)
```

İlişkili kayıtlara filter uygula.

---

## RelationshipSort

Relationship sıralaması için kullanılan struct.

### Metodlar

#### ApplySort

```go
func (s *RelationshipSort) ApplySort(ctx context.Context, column, direction string) ([]interface{}, error)
```

İlişkili kayıtlara sort uygula.

---

## RelationshipPagination

Relationship sayfalandırması için kullanılan struct.

### Metodlar

#### ApplyPagination

```go
func (p *RelationshipPagination) ApplyPagination(ctx context.Context, page, perPage int) ([]interface{}, error)
```

İlişkili kayıtlara pagination uygula.

---

## RelationshipConstraints

Relationship kısıtlamaları için kullanılan struct.

### Metodlar

#### ApplyLimit

```go
func (c *RelationshipConstraints) ApplyLimit(limit int) *RelationshipConstraints
```

LIMIT kısıtlaması uygula.

#### ApplyOffset

```go
func (c *RelationshipConstraints) ApplyOffset(offset int) *RelationshipConstraints
```

OFFSET kısıtlaması uygula.

#### ApplyWhere

```go
func (c *RelationshipConstraints) ApplyWhere(column, operator string, value interface{}) *RelationshipConstraints
```

WHERE kısıtlaması uygula.

#### ApplyWhereIn

```go
func (c *RelationshipConstraints) ApplyWhereIn(column string, values []interface{}) *RelationshipConstraints
```

WHERE IN kısıtlaması uygula.

---

## RelationshipCounting

Relationship sayması için kullanılan struct.

### Metodlar

#### Count

```go
func (c *RelationshipCounting) Count(ctx context.Context, data interface{}) (int, error)
```

İlişkili kayıtların sayısını döndür.

---

## RelationshipExistence

Relationship varlık kontrolü için kullanılan struct.

### Metodlar

#### Exists

```go
func (e *RelationshipExistence) Exists(ctx context.Context, data interface{}) (bool, error)
```

İlişkili kayıtların varlığını kontrol et.

#### DoesntExist

```go
func (e *RelationshipExistence) DoesntExist(ctx context.Context, data interface{}) (bool, error)
```

İlişkili kayıtların yokluğunu kontrol et.

---

## RelationshipSerialization

Relationship JSON serileştirmesi için kullanılan struct.

### Metodlar

#### SerializeRelationship

```go
func (s *RelationshipSerialization) SerializeRelationship(data interface{}) (map[string]interface{}, error)
```

İlişkili kayıtları JSON-compatible format'a dönüştür.

#### ToJSON

```go
func (s *RelationshipSerialization) ToJSON(data interface{}) (string, error)
```

İlişkili kayıtları JSON string'e dönüştür.

---

## RelationshipLoader

Relationship yüklemesi için kullanılan struct.

### Metodlar

#### EagerLoad

```go
func (l *RelationshipLoader) EagerLoad(ctx context.Context, items []interface{}, field RelationshipField) error
```

Eager loading stratejisini uygula.

#### LazyLoad

```go
func (l *RelationshipLoader) LazyLoad(ctx context.Context, item interface{}, field RelationshipField) (interface{}, error)
```

Lazy loading stratejisini uygula.

---

## Sabitler

### LoadingStrategy

```go
const (
    EAGER_LOADING LoadingStrategy = iota
    LAZY_LOADING
)
```

---

## Hata İşleme

### RelationshipError

```go
type RelationshipError struct {
    Type    string
    Field   string
    Message string
    Context map[string]interface{}
}
```

Relationship işlemleri sırasında oluşan hatalar `RelationshipError` türünde döndürülür.

**Örnek:**
```go
if err != nil {
    if relErr, ok := err.(*RelationshipError); ok {
        fmt.Printf("İlişki Hatası: %s\n", relErr.Message)
        fmt.Printf("Tür: %s\n", relErr.Type)
        fmt.Printf("Alan: %s\n", relErr.Field)
    }
}
```

---

## Ayrıca Bkz.

- [İlişkiler](./Relationships.md) - İlişkiler genel bakış
- [Alanlar](./Fields.md) - Alanlar genel bakış
- [API Referansı](./API-Reference.md) - Tam API referansı
