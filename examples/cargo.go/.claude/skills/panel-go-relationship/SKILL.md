---
name: panel-go-relationship
description: Implement database relationships for Panel.go with BelongsTo, HasMany, HasOne, BelongsToMany, and polymorphic relationships. Use when adding relationships between entities or configuring eager loading.
allowed-tools: Read, Write, Edit
---

# Panel.go Relationship Expert

Expert in implementing database relationships for Panel.go with BelongsTo, HasMany, HasOne, BelongsToMany, and polymorphic relationships.

## Expertise

- **BelongsTo**: One-to-one, child to parent (Post belongs to User)
- **HasOne**: One-to-one, parent to child (User has one Profile)
- **HasMany**: One-to-many (User has many Posts)
- **BelongsToMany**: Many-to-many (Post belongs to many Tags)
- **MorphTo**: Polymorphic relationships
- **Eager Loading**: Preload, Joins to prevent N+1

## Quick Patterns

### BelongsTo
```go
type Post struct {
    ID     uint
    UserID uint
    User   *User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

fields.Link("user", &user.UserResource{}).DisplayKey("name")
```

### HasMany
```go
type User struct {
    ID    uint
    Posts []Post `gorm:"foreignKey:UserID"`
}

fields.Collection("posts", &post.PostResource{}).DisplayKey("title")
```

### BelongsToMany
```go
type Post struct {
    ID   uint
    Tags []Tag `gorm:"many2many:post_tags;"`
}

fields.Connect("tags", &tag.TagResource{}).DisplayKey("name")
```

### Eager Loading
```go
func (r *PostResource) With() []string {
    return []string{"User", "Category", "Tags"}
}
```

## Key Rules

- **Always index foreign keys**
- **Use Preload** to prevent N+1 queries
- **Define cascade rules** (CASCADE, RESTRICT, SET NULL)
- **Set DisplayKey** for relationships
- **Add With()** method for eager loading

## Usage

When user asks to add relationships:

1. Add foreign key to entity
2. Add relationship field with GORM tags
3. Add relationship field to field resolver
4. Set DisplayKey
5. Add eager loading in With() method

## Anti-Patterns

❌ **Don't forget foreign key indexes** - Always index foreign keys
❌ **Don't skip eager loading** - Use With() to prevent N+1 queries
❌ **Don't forget DisplayKey** - Set DisplayKey for relationship fields
❌ **Don't skip cascade rules** - Define OnUpdate and OnDelete behavior
❌ **Don't use wrong field types** - Link for BelongsTo, Collection for HasMany, Connect for BelongsToMany

## Sharp Edges

⚠️ **N+1 Queries**: Always use Preload or Joins for relationships
⚠️ **Cascade Deletes**: Be careful with CASCADE on delete
⚠️ **Many-to-Many**: Requires junction table (many2many tag)
⚠️ **DisplayKey**: Must match actual field name in related entity
⚠️ **Polymorphic**: Use MorphTo for polymorphic relationships
