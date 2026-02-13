---
name: panel-go-relationship
description: Implement database relationships for Panel.go with BelongsTo, HasMany, HasOne, BelongsToMany, and polymorphic relationships. Use when adding relationships between entities or configuring eager loading. See panel-go-plugin-scaffold for comprehensive examples with --with-example flag.
allowed-tools: Read, Write, Edit
---

# Panel.go Relationship Expert

Expert in implementing database relationships for Panel.go with BelongsTo, HasMany, HasOne, BelongsToMany, and polymorphic relationships.

## Expertise

- **BelongsTo**: One-to-one, child to parent (Post belongs to User)
- **HasOne**: One-to-one, parent to child (User has one Profile)
- **HasMany**: One-to-many (User has many Posts)
- **BelongsToMany**: Many-to-many (Post belongs to many Tags)
- **MorphTo**: Polymorphic relationships (Comment belongs to Commentable)
- **MorphToMany**: Polymorphic many-to-many (Tag belongs to many Taggables)
- **Eager Loading**: Preload, Joins to prevent N+1
- **Entity.go Pattern**: All entities in one file to prevent circular dependency
- **Registry Pattern**: Resource access without direct imports

## Quick Patterns

### BelongsTo
```go
type Post struct {
    ID     uint
    UserID uint
    User   *User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

fields.BelongsTo("User", "user_id", resources.GetOrPanic("users")).
    DisplayUsing("name").
    WithSearchableColumns("name", "email").
    WithEagerLoad()
```

### HasMany
```go
type User struct {
    ID    uint
    Posts []Post `gorm:"foreignKey:UserID"`
}

fields.HasMany("Posts", "posts", resources.GetOrPanic("posts")).
    ForeignKey("user_id").
    OwnerKey("id").
    WithEagerLoad()
```

### HasOne
```go
type Organization struct {
    ID          uint
    BillingInfo *BillingInfo `gorm:"foreignKey:OrganizationID"`
}

type BillingInfo struct {
    ID             uint
    OrganizationID uint `gorm:"not null;unique;index"`
    Organization   *Organization
}

fields.HasOne("BillingInfo", "billing_info", resources.GetOrPanic("billing-info")).
    ForeignKey("organization_id").
    OwnerKey("id").
    WithEagerLoad()
```

### BelongsToMany
```go
type Post struct {
    ID   uint
    Tags []Tag `gorm:"many2many:post_tags;"`
}

type Tag struct {
    ID    uint
    Posts []Post `gorm:"many2many:post_tags;"`
}

fields.BelongsToMany("Tags", "tags", resources.GetOrPanic("tags")).
    PivotTable("post_tags").
    WithEagerLoad()
```

### MorphTo (Polymorphic)
```go
type Comment struct {
    ID              uint
    CommentableID   uint   `gorm:"index:idx_commentable"`
    CommentableType string `gorm:"index:idx_commentable"`
    Content         string
}

// Entity'de tanımlanır, ancak field henüz implement edilmemiş
```

### MorphToMany (Polymorphic Many-to-Many)
```go
type Tag struct {
    ID   uint
    Name string
}

// Taggable pivot table ile Product ve Shipment'a bağlanır
// Entity'de tanımlanır, ancak field henüz implement edilmemiş
```

### Eager Loading
```go
func (r *PostResource) With() []string {
    return []string{"User", "Category", "Tags"}
}
```

## Entity.go Pattern

**Tüm entity'leri tek dosyada tut** (circular dependency önlenir):

```go
// entity/entity.go
package entity

type Organization struct {
    ID          uint
    Addresses   []Address    // HasMany
    BillingInfo *BillingInfo // HasOne
}

type Address struct {
    ID             uint
    OrganizationID uint
    Organization   *Organization // BelongsTo
}

type BillingInfo struct {
    ID             uint
    OrganizationID uint `gorm:"unique"`
    Organization   *Organization // BelongsTo (HasOne)
}
```

## Registry Pattern

**Resource'lar birbirini import etmeden erişir**:

```go
// resources/registry.go
var registry = make(map[string]func() resource.Resource)

func Register(slug string, factory func() resource.Resource) {
    registry[slug] = factory
    resource.Register(slug, factory())
}

func GetOrPanic(slug string) resource.Resource {
    if r := Get(slug); r != nil {
        return r
    }
    panic("resource not found: " + slug)
}

// Kullanım
fields.HasMany("Addresses", "addresses", resources.GetOrPanic("addresses"))
```

## Key Rules

- **Always index foreign keys**
- **Use Preload** to prevent N+1 queries
- **Define cascade rules** (CASCADE, RESTRICT, SET NULL)
- **Set DisplayKey** for relationships
- **Add With()** method for eager loading
- **Use entity.go pattern** for all entities
- **Use registry pattern** for resource access

## Usage

When user asks to add relationships:

1. Add foreign key to entity in `entity/entity.go`
2. Add relationship field with GORM tags
3. Add relationship field to field resolver using `resources.GetOrPanic()`
4. Set DisplayUsing, WithSearchableColumns, WithEagerLoad
5. Add eager loading in With() method

## Comprehensive Examples

For comprehensive examples with all relationship types:

```bash
panel plugin create example --with-example
```

This creates 9 entities and 7 resources demonstrating:
- BelongsTo, HasMany, HasOne, BelongsToMany
- MorphTo, MorphToMany (entity level)
- Entity.go pattern
- Registry pattern

See `panel-go-plugin-scaffold` skill for details.

## Anti-Patterns

❌ **Don't forget foreign key indexes** - Always index foreign keys
❌ **Don't skip eager loading** - Use With() to prevent N+1 queries
❌ **Don't forget DisplayUsing** - Set DisplayUsing for relationship fields
❌ **Don't skip cascade rules** - Define OnUpdate and OnDelete behavior
❌ **Don't use wrong field types** - BelongsTo for belongs_to, HasMany for has_many, etc.
❌ **Don't separate entities** - Use entity.go pattern to prevent circular dependency
❌ **Don't import resources directly** - Use registry pattern

## Sharp Edges

⚠️ **N+1 Queries**: Always use WithEagerLoad() for relationships
⚠️ **Cascade Deletes**: Be careful with CASCADE on delete
⚠️ **Many-to-Many**: Requires junction table (many2many tag)
⚠️ **DisplayUsing**: Must match actual field name in related entity
⚠️ **Polymorphic**: MorphTo/MorphToMany fields not yet implemented in panel.go
⚠️ **Circular Dependency**: Use entity.go pattern to prevent import cycles
⚠️ **Registry Timing**: init() functions run on import, order matters
