---
name: panel-go-migration
description: Generate GORM migrations for Panel.go database schema with proper indexes, constraints, and relationships. Use when creating database tables, adding columns, or defining entity structures.
allowed-tools: Read, Write, Edit, Bash
---

# Panel.go Migration Generator

Expert in generating GORM migrations for Panel.go database schema with proper indexes, constraints, and relationships.

## Expertise

- **GORM Migrations**: AutoMigrate, CreateTable, AddColumn, AddIndex
- **Data Types**: String, Int, Bool, Time, JSON, UUID
- **Constraints**: Primary Key, Foreign Key, Unique, Not Null, Default
- **Indexes**: Single column, composite, unique indexes
- **Relationships**: Foreign keys, cascade deletes, on update
- **Soft Deletes**: GORM's DeletedAt field
- **Timestamps**: CreatedAt, UpdatedAt auto-management

## Quick Patterns

### Basic Entity
```go
type Product struct {
    ID          uint      `gorm:"primaryKey"`
    Name        string    `gorm:"size:255;not null;index"`
    Description string    `gorm:"type:text"`
    Price       float64   `gorm:"type:decimal(10,2);not null"`
    Stock       int       `gorm:"default:0;not null"`
    CategoryID  uint      `gorm:"not null;index"`
    Category    *Category `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
    CreatedAt   time.Time `gorm:"index"`
    UpdatedAt   time.Time
    DeletedAt   gorm.DeletedAt `gorm:"index"`
}
```

### Run Migrations
```go
func RunMigrations(db *gorm.DB) error {
    return db.AutoMigrate(
        &domain.User{},
        &domain.Post{},
        &domain.Category{},
    )
}
```

## Key Rules

- **Always index foreign keys**
- **Use soft deletes** (gorm.DeletedAt)
- **Add timestamps** (CreatedAt, UpdatedAt)
- **Plural table names** (users, posts)
- **Define cascade rules** (OnUpdate, OnDelete)
- **Index searchable fields**

## Usage

When user asks to create a migration or entity:

1. Create entity in `pkg/domain/{name}/entity.go`
2. Add GORM tags for constraints and indexes
3. Define TableName() method
4. Add to AutoMigrate in migrations.go

## Anti-Patterns

❌ **Don't skip indexes on foreign keys** - Always index foreign keys
❌ **Don't forget soft deletes** - Use gorm.DeletedAt for audit trails
❌ **Don't skip timestamps** - Always add CreatedAt, UpdatedAt
❌ **Don't use singular table names** - Use plural (users, not user)
❌ **Don't forget cascade rules** - Define OnUpdate and OnDelete behavior

## Sharp Edges

⚠️ **Foreign Keys**: Always define constraint behavior (CASCADE, RESTRICT, SET NULL)
⚠️ **Indexes**: Index all foreign keys and searchable fields
⚠️ **Soft Deletes**: Use gorm.DeletedAt, not custom deleted_at
⚠️ **Table Names**: Override TableName() for custom names
⚠️ **Migration Order**: Migrate parent tables before child tables
