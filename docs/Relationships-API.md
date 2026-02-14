# İlişkiler API Referansı

Bu referans, doğrudan `pkg/fields` ve `pkg/resource/registry.go` içindeki güncel API imzalarına göre hazırlanmıştır.

## Ortak Kavramlar

### Loading Strategy

```go
type LoadingStrategy string

const (
    EAGER_LOADING LoadingStrategy = "eager"
    LAZY_LOADING  LoadingStrategy = "lazy"
)
```

### RelationshipField Interface

```go
type RelationshipField interface {
    Element

    GetRelationshipType() string
    GetRelatedResourceSlug() string
    GetRelationshipName() string

    ResolveRelationship(item interface{}) (interface{}, error)
    GetQueryCallback() func(interface{}) interface{}
    GetLoadingStrategy() LoadingStrategy

    ValidateRelationship(value interface{}) error

    GetDisplayKey() string
    GetSearchableColumns() []string

    IsRequired() bool
    GetTypes() map[string]string
}
```

Notlar:
- `Query` callback imzası: `func(interface{}) interface{}`
- Callback içinde çoğu senaryoda `*gorm.DB` type assertion yapılır.
- Tüm ilişki tipleri bu interface’i birebir aynı şekilde implemente etmek zorunda değildir; pratikte endpoint tarafında view/type bazlı işleme de gidilir.

## Constructor'lar

```go
func BelongsTo(name, key string, relatedResource interface{}) *BelongsToField
func HasMany(name, key string, relatedResource interface{}) *HasManyField
func HasOne(name, key string, relatedResource interface{}) *HasOneField
func BelongsToMany(name, key string, relatedResource interface{}) *BelongsToManyField

func NewMorphTo(name, key string) *MorphTo
func NewMorphToMany(name, key string) *MorphToMany
```

`relatedResource` parametresi:
- `string` slug (örn. `"users"`)
- `Slug() string` dönen resource instance

## BelongsToField

Temel metodlar:

```go
func (b *BelongsToField) AutoOptions(displayField string) *BelongsToField
func (b *BelongsToField) DisplayUsing(key string) *BelongsToField
func (b *BelongsToField) WithSearchableColumns(columns ...string) *BelongsToField
func (b *BelongsToField) Query(fn func(interface{}) interface{}) *BelongsToField
func (b *BelongsToField) WithEagerLoad() *BelongsToField
func (b *BelongsToField) WithLazyLoad() *BelongsToField

func (b *BelongsToField) GetRelationshipType() string
func (b *BelongsToField) GetRelatedResource() string
func (b *BelongsToField) GetRelatedResourceSlug() string
func (b *BelongsToField) GetRelationshipName() string
func (b *BelongsToField) GetDisplayKey() string
func (b *BelongsToField) GetSearchableColumns() []string
func (b *BelongsToField) GetQueryCallback() func(interface{}) interface{}
func (b *BelongsToField) GetLoadingStrategy() LoadingStrategy
func (b *BelongsToField) IsRequired() bool

func (b *BelongsToField) GetForeignKey() string
func (b *BelongsToField) GetRelatedTableName() string
func (b *BelongsToField) GetOwnerKeyColumn() string

func (b *BelongsToField) WithHoverCard(config HoverCardConfig) *BelongsToField
func (b *BelongsToField) HoverCard(hoverStruct interface{}) *BelongsToField
func (b *BelongsToField) ResolveHoverCard(resolver HoverCardResolver) *BelongsToField
func (b *BelongsToField) GetHoverCard() *HoverCardConfig
```

## HasManyField

```go
func (h *HasManyField) AutoOptions(displayField string) *HasManyField
func (h *HasManyField) ForeignKey(key string) *HasManyField
func (h *HasManyField) OwnerKey(key string) *HasManyField
func (h *HasManyField) Query(fn func(interface{}) interface{}) *HasManyField
func (h *HasManyField) WithEagerLoad() *HasManyField
func (h *HasManyField) WithLazyLoad() *HasManyField
func (h *HasManyField) WithFullData() *HasManyField

func (h *HasManyField) SetRelatedResource(res interface{}) *HasManyField
func (h *HasManyField) GetRelatedResource() interface{}
func (h *HasManyField) GetRelatedResourceSlug() string

func (h *HasManyField) GetRelationshipType() string
func (h *HasManyField) GetRelationshipName() string
func (h *HasManyField) GetQueryCallback() func(interface{}) interface{}
func (h *HasManyField) GetLoadingStrategy() LoadingStrategy
func (h *HasManyField) IsRequired() bool

func (h *HasManyField) GetRelatedTableName() string
func (h *HasManyField) GetForeignKeyColumn() string
func (h *HasManyField) GetOwnerKeyColumn() string
```

## HasOneField

```go
func (h *HasOneField) AutoOptions(displayField string) *HasOneField
func (h *HasOneField) ForeignKey(key string) *HasOneField
func (h *HasOneField) OwnerKey(key string) *HasOneField
func (h *HasOneField) Query(fn func(interface{}) interface{}) *HasOneField
func (h *HasOneField) WithEagerLoad() *HasOneField
func (h *HasOneField) WithLazyLoad() *HasOneField

func (h *HasOneField) GetRelatedResourceSlug() string
func (h *HasOneField) GetRelationshipType() string
func (h *HasOneField) GetRelationshipName() string
func (h *HasOneField) GetQueryCallback() func(interface{}) interface{}
func (h *HasOneField) GetLoadingStrategy() LoadingStrategy
func (h *HasOneField) IsRequired() bool

func (h *HasOneField) GetRelatedTableName() string
func (h *HasOneField) GetForeignKeyColumn() string
func (h *HasOneField) GetOwnerKeyColumn() string

func (h *HasOneField) WithHoverCard(config HoverCardConfig) *HasOneField
func (h *HasOneField) HoverCard(hoverStruct interface{}) *HasOneField
func (h *HasOneField) ResolveHoverCard(resolver HoverCardResolver) *HasOneField
func (h *HasOneField) GetHoverCard() *HoverCardConfig
```

## BelongsToManyField

```go
func (b *BelongsToManyField) AutoOptions(displayField string) *BelongsToManyField
func (b *BelongsToManyField) PivotTable(table string) *BelongsToManyField
func (b *BelongsToManyField) ForeignKey(key string) *BelongsToManyField
func (b *BelongsToManyField) RelatedKey(key string) *BelongsToManyField
func (b *BelongsToManyField) Query(fn func(interface{}) interface{}) *BelongsToManyField
func (b *BelongsToManyField) WithEagerLoad() *BelongsToManyField
func (b *BelongsToManyField) WithLazyLoad() *BelongsToManyField

func (b *BelongsToManyField) GetRelatedResourceSlug() string
func (b *BelongsToManyField) GetRelationshipType() string
func (b *BelongsToManyField) GetRelationshipName() string
func (b *BelongsToManyField) GetQueryCallback() func(interface{}) interface{}
func (b *BelongsToManyField) GetLoadingStrategy() LoadingStrategy
func (b *BelongsToManyField) IsRequired() bool
```

## MorphTo

```go
func (m *MorphTo) Types(types map[string]string) *MorphTo
func (m *MorphTo) Displays(displays map[string]string) *MorphTo
func (m *MorphTo) Query(fn func(interface{}) interface{}) *MorphTo
func (m *MorphTo) WithEagerLoad() *MorphTo
func (m *MorphTo) WithLazyLoad() *MorphTo

func (m *MorphTo) GetRelationshipType() string
func (m *MorphTo) GetRelatedResourceSlug() string
func (m *MorphTo) GetRelationshipName() string
func (m *MorphTo) GetQueryCallback() func(interface{}) interface{}
func (m *MorphTo) GetLoadingStrategy() LoadingStrategy
func (m *MorphTo) GetTypes() map[string]string
func (m *MorphTo) GetResourceForType(morphType string) (string, error)

func (m *MorphTo) IsRequired() bool

func (m *MorphTo) WithHoverCard(config HoverCardConfig) *MorphTo
func (m *MorphTo) HoverCard(hoverStruct interface{}) *MorphTo
func (m *MorphTo) ResolveHoverCard(resolver HoverCardResolver) *MorphTo
func (m *MorphTo) GetHoverCard() *HoverCardConfig
```

## MorphToMany

```go
func (m *MorphToMany) Types(types map[string]string) *MorphToMany
func (m *MorphToMany) Displays(displays map[string]string) *MorphToMany
func (m *MorphToMany) PivotTable(tableName string) *MorphToMany
func (m *MorphToMany) ForeignKey(key string) *MorphToMany
func (m *MorphToMany) RelatedKey(key string) *MorphToMany
func (m *MorphToMany) MorphType(column string) *MorphToMany
func (m *MorphToMany) AutoOptions(displayField string) *MorphToMany
func (m *MorphToMany) Query(fn func(interface{}) interface{}) *MorphToMany
func (m *MorphToMany) WithEagerLoad() *MorphToMany
func (m *MorphToMany) WithLazyLoad() *MorphToMany

func (m *MorphToMany) GetRelationshipType() string
func (m *MorphToMany) GetRelatedResource() string
func (m *MorphToMany) GetRelationshipName() string
func (m *MorphToMany) GetQueryCallback() func(interface{}) interface{}
func (m *MorphToMany) GetLoadingStrategy() LoadingStrategy
func (m *MorphToMany) GetTypes() map[string]string
func (m *MorphToMany) IsRequired() bool
```

## Query Callback Örneği

```go
field := fields.BelongsTo("Organization", "organization_id", "organizations").
    Query(func(q interface{}) interface{} {
        db, ok := q.(*gorm.DB)
        if !ok || db == nil {
            return q
        }
        return db.Where("deleted_at IS NULL")
    })
```

## Resource Registry API

`pkg/resource` içindeki global registry API:

```go
func Register(slug string, res Resource)
func Get(slug string) Resource
func List() map[string]Resource
func Clear()
```

Notlar:
- `resource.Register` factory değil, doğrudan `Resource` instance alır.
- `resource.GetOrPanic` fonksiyonu çekirdek (`pkg/resource`) API’sinde yoktur.
- `Clear()` test senaryoları içindir.
