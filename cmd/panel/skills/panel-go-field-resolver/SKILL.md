---
name: panel-go-field-resolver
description: Create dynamic field configurations for Panel.go admin resources with validation, display options, and relationship handling. Use when adding fields to resources, configuring forms, or setting up field validation.
allowed-tools: Read, Write, Edit
---

# Panel.go Field Resolver Expert

Expert in creating dynamic field configurations for Panel.go admin resources with validation, display options, and relationship handling.

## Expertise

You are a Panel.go field resolver expert. You understand:

- **Field Types**: Text, Number, Email, Image, Select, Date, Relationship fields (31 types)
- **Field Configuration**: Required, Unique, Searchable, ReadOnly, Nullable
- **Display Options**: OnList, OnDetail, OnForm, OnlyOnDetail, OnlyOnForm
- **Validation**: Built-in validators, custom validation rules
- **Relationships**: Link (BelongsTo), Detail (HasOne), Collection (HasMany), Connect (BelongsToMany)
- **Options**: Static options, dynamic options from database
- **File Uploads**: Image, Video, Audio, File fields with storage
- **Rich Content**: RichText, Code, Color, KeyValue fields

## Patterns

### Basic Field Resolver

```go
type {Name}FieldResolver struct{}

func (r *{Name}FieldResolver) ResolveFields(ctx *context.Context) []core.Element {
    return []core.Element{
        // Primary key
        fields.ID().ReadOnly().OnlyOnDetail(),

        // Text fields
        fields.Text("name").Required().Searchable(),
        fields.Email("email").Required().Unique(),
        fields.Textarea("description").Nullable(),

        // Number fields
        fields.Number("price").Required().Min(0),
        fields.Number("quantity").Default(0),

        // Select fields
        fields.Select("status").Options([]fields.Option{
            {Label: "Active", Value: "active"},
            {Label: "Inactive", Value: "inactive"},
        }).Default("active"),

        // Boolean
        fields.Switch("is_active").Default(true),

        // Date fields
        fields.Date("published_at").Nullable(),
        fields.DateTime("created_at").ReadOnly().OnlyOnDetail(),

        // Relationships
        fields.Link("user", &user.UserResource{}).DisplayKey("name"),
    }
}
```

### File Upload Fields

```go
fields.Image("avatar").
    Disk("public").
    Path("avatars").
    MaxSize(2048). // 2MB
    Accept("image/jpeg,image/png"),

fields.File("document").
    Disk("private").
    Path("documents").
    MaxSize(10240), // 10MB
```

### Relationship Fields

```go
// BelongsTo
fields.Link("category", &category.CategoryResource{}).
    DisplayKey("name").
    Searchable(),

// HasMany
fields.Collection("posts", &post.PostResource{}).
    DisplayKey("title"),

// BelongsToMany
fields.Connect("tags", &tag.TagResource{}).
    DisplayKey("name").
    Searchable(),
```

### Dynamic Options

```go
fields.Select("category_id").
    Options(func(ctx *context.Context) []fields.Option {
        var categories []domain.Category
        ctx.DB().Find(&categories)

        options := make([]fields.Option, len(categories))
        for i, cat := range categories {
            options[i] = fields.Option{
                Label: cat.Name,
                Value: cat.ID,
            }
        }
        return options
    }),
```

## Anti-Patterns

❌ **Don't expose sensitive fields** - Use ReadOnly() or hide completely
❌ **Don't skip validation** - Always validate required fields
❌ **Don't forget indexes** - Add Searchable() for indexed fields
❌ **Don't hardcode options** - Use dynamic options for database data
❌ **Don't ignore file size limits** - Set MaxSize() for uploads
❌ **Don't forget display keys** - Set DisplayKey() for relationships
❌ **Don't mix concerns** - Keep field logic in resolver, business logic in repository

## Decisions

### Field Naming
- **snake_case** for field names (matches database columns)
- **Descriptive names** that match domain model
- **Consistent naming** across resources

### Validation
- **Required** for non-nullable fields
- **Unique** for unique constraints
- **Min/Max** for number ranges
- **Custom validators** for complex rules

### Display
- **OnList** for table columns
- **OnDetail** for detail view
- **OnForm** for create/edit forms
- **OnlyOnDetail** for read-only info
- **OnlyOnForm** for input-only fields

## Sharp Edges

⚠️ **File Uploads**: Configure storage disk and path correctly
⚠️ **Relationships**: Use DisplayKey to show meaningful data
⚠️ **Options**: Cache dynamic options for performance
⚠️ **Validation**: Validate on both client and server
⚠️ **Searchable**: Only mark indexed fields as searchable
⚠️ **ReadOnly**: Use for computed or system fields
⚠️ **Nullable**: Match database schema nullability

## Usage

When user asks to add fields to a resource:

1. **Identify field types** based on data type
2. **Add validation** rules (Required, Unique, Min, Max)
3. **Configure display** options (OnList, OnDetail, OnForm)
4. **Add relationships** if needed
5. **Set defaults** for optional fields
6. **Configure file uploads** if needed
7. **Add searchable** for indexed fields
