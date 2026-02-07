---
name: panel-go-resource
description: Generate complete admin panel resources with CRUD operations, field resolvers, card resolvers, and policies for Panel.go projects. Use when creating new resources or setting up complete admin functionality.
allowed-tools: Read, Write, Edit, Bash
---

# Panel.go Resource Generator

Expert in generating complete admin panel resources with CRUD operations, field resolvers, card resolvers, and policies for Panel.go projects.

## Expertise

You are a Panel.go resource generation expert. You understand:

- **Resource Architecture**: OptimizedBase embedding, Repository pattern, GORM integration
- **Field Resolvers**: Dynamic field configuration, validation rules, display options
- **Card Resolvers**: Dashboard widgets, metrics, statistics cards
- **Policy System**: RBAC authorization, ViewAny/View/Create/Update/Delete methods
- **Repository Pattern**: Data provider abstraction, GORM repositories
- **Relationships**: BelongsTo, HasMany, HasOne, BelongsToMany, MorphTo
- **CRUD Operations**: Index, Show, Create, Store, Edit, Update, Destroy
- **Field Types**: Text, Number, Email, Image, Select, Date, Relationship fields

## Patterns

### Resource Structure

```go
// pkg/resource/{name}/resource.go
type {Name}Resource struct {
    resource.OptimizedBase
}

func New{Name}Resource() *{Name}Resource {
    r := &{Name}Resource{}
    r.Model = &domain.{Name}{}
    r.Slug = "{slug}"
    r.Title = "{Title}"
    r.Icon = "icon-name"
    r.Group = "group-name"
    r.NavigationOrder = 1
    r.Visible = true
    r.SetPolicy(&{Name}Policy{})
    r.SetFieldResolver(&{Name}FieldResolver{})
    r.SetCardResolver(&{Name}CardResolver{})
    return r
}

func (r *{Name}Resource) Repository(db *gorm.DB) data.DataProvider {
    return orm.New{Name}Repository(db)
}
```

### Field Resolver

```go
// pkg/resource/{name}/field_resolver.go
type {Name}FieldResolver struct{}

func (r *{Name}FieldResolver) ResolveFields(ctx *context.Context) []core.Element {
    return []core.Element{
        fields.ID().ReadOnly().OnlyOnDetail(),
        fields.Text("name").Required().Searchable(),
        fields.Email("email").Required().Unique(),
        fields.Select("status").Options([]fields.Option{
            {Label: "Active", Value: "active"},
            {Label: "Inactive", Value: "inactive"},
        }),
        fields.Link("user", &user.UserResource{}).DisplayKey("name"),
    }
}
```

### Card Resolver

```go
// pkg/resource/{name}/card_resolver.go
type {Name}CardResolver struct{}

func (r *{Name}CardResolver) ResolveCards(ctx *context.Context) []widget.Card {
    return []widget.Card{
        // Add dashboard cards here
    }
}
```

### Policy

```go
// pkg/resource/{name}/policy.go
type {Name}Policy struct{}

func (p *{Name}Policy) ViewAny(ctx *context.Context) bool {
    return ctx.User().HasPermission("{resource}.view_any")
}

func (p *{Name}Policy) View(ctx *context.Context, model interface{}) bool {
    return ctx.User().HasPermission("{resource}.view")
}

func (p *{Name}Policy) Create(ctx *context.Context) bool {
    return ctx.User().HasPermission("{resource}.create")
}

func (p *{Name}Policy) Update(ctx *context.Context, model interface{}) bool {
    return ctx.User().HasPermission("{resource}.update")
}

func (p *{Name}Policy) Delete(ctx *context.Context, model interface{}) bool {
    return ctx.User().HasPermission("{resource}.delete")
}
```

### Repository

```go
// pkg/data/orm/{name}_repository.go
type {Name}Repository struct {
    *data.GormDataProvider
    db *gorm.DB
}

func New{Name}Repository(db *gorm.DB) *{Name}Repository {
    return &{Name}Repository{
        GormDataProvider: data.NewGormDataProvider(db, &domain.{Name}{}),
        db: db,
    }
}
```

### Domain Entity

```go
// pkg/domain/{name}/entity.go
type {Name} struct {
    ID        uint      `gorm:"primaryKey"`
    Name      string    `gorm:"size:255;not null;index"`
    Email     string    `gorm:"size:255;unique;not null"`
    Status    string    `gorm:"size:50;default:'active'"`
    CreatedAt time.Time
    UpdatedAt time.Time
}

func (u *{Name}) TableName() string {
    return "{table_name}"
}
```

## Anti-Patterns

❌ **Don't hardcode permissions** - Use permission manager
❌ **Don't skip validation** - Always validate user input
❌ **Don't expose sensitive fields** - Use ReadOnly() or hide fields
❌ **Don't forget relationships** - Use proper relationship fields
❌ **Don't skip eager loading** - Use With() for N+1 prevention
❌ **Don't ignore context** - Always check context.User()
❌ **Don't forget indexes** - Add database indexes for searchable fields

## Decisions

### Resource Naming
- **Singular names** for resources (User, Post, Product)
- **Plural slugs** for URLs (users, posts, products)
- **PascalCase** for Go types
- **kebab-case** for slugs

### Field Configuration
- **Required fields** should have validation
- **Unique fields** should have database constraints
- **Searchable fields** should have indexes
- **Relationship fields** should use DisplayKey

### Authorization
- **Policy-based** authorization for all resources
- **Permission checks** in all policy methods
- **Context-aware** decisions based on user and model

### Repository Pattern
- **GORM-based** repositories for database operations
- **Interface abstraction** for testability
- **Transaction support** for complex operations

## Sharp Edges

⚠️ **Password Fields**: Never return password in API responses
⚠️ **File Uploads**: Handle file storage and validation properly
⚠️ **Relationships**: Use eager loading to prevent N+1 queries
⚠️ **Soft Deletes**: Use GORM's soft delete for audit trails
⚠️ **Validation**: Validate on both client and server side
⚠️ **Permissions**: Check permissions in both policy and middleware
⚠️ **Context**: Always pass context through the call chain

## Usage

When user asks to create a new resource:

1. **Create domain entity** in `pkg/domain/{name}/entity.go`
2. **Create repository** in `pkg/data/orm/{name}_repository.go`
3. **Create resource** in `pkg/resource/{name}/resource.go`
4. **Create field resolver** in `pkg/resource/{name}/field_resolver.go`
5. **Create card resolver** in `pkg/resource/{name}/card_resolver.go`
6. **Create policy** in `pkg/resource/{name}/policy.go`
7. **Register resource** in main.go: `panel.Register(resource.New{Name}Resource())`
8. **Create migration** for database table
9. **Add permissions** to permissions.toml
