# Panel.go - Optimizasyon Rehberi

## Giriş

Panel.go, mixin pattern ve resolver pattern kullanarak kod tekrarını azaltmış ve bakım kolaylığını artırmıştır. Bu rehber, optimizasyon pattern'lerini nasıl kullanacağınızı gösterir.

## Temel Konseptler

### 1. Resolver Pattern

Resolver pattern, dinamik çözümleme için kullanılır. Her bir resolver, belirli bir sorumluluğu yerine getirir:

- **FieldResolver**: Alanları dinamik olarak çözer
- **CardResolver**: Card'ları dinamik olarak çözer
- **FilterResolver**: Filtreleri dinamik olarak çözer
- **LensResolver**: Lens'leri dinamik olarak çözer
- **ActionResolver**: İşlemleri dinamik olarak çözer

### 2. Mixin Pattern

Mixin'ler, ortak işlevselliği sağlayan embedded struct'lardır. Composition kullanarak kod tekrarını azaltırlar.

**Resource Mixin'leri:**
- `Authorizable` - Yetkilendirme
- `Resolvable` - Çözümleme
- `Navigable` - Navigasyon

**Field Mixin'leri:**
- `Searchable` - Arama
- `Sortable` - Sıralama
- `Filterable` - Filtreleme
- `Validatable` - Doğrulama
- `Displayable` - Görüntüleme
- `Hideable` - Gizlilik

## Resource Implementasyonu

### Eski Yol (Hala Destekleniyor)

```go
type UserResource struct {
    resource.Base
}

func (r UserResource) Model() interface{} {
    return &User{}
}

func (r UserResource) Fields() []fields.Element {
    return []fields.Element{
        fields.NewText("name"),
        fields.NewEmail("email"),
    }
}

// ... 20+ metod daha
```

### Yeni Yol (Önerilen)

```go
package user

import (
    "github.com/ferdiunal/panel.go/pkg/context"
    "github.com/ferdiunal/panel.go/pkg/fields"
    "github.com/ferdiunal/panel.go/pkg/resource"
    "github.com/ferdiunal/panel.go/pkg/widget"
)

// UserResource, kullanıcı kaynağını temsil eder
type UserResource struct {
    resource.OptimizedBase
}

// NewUserResource, yeni bir UserResource oluşturur
func NewUserResource() *UserResource {
    r := &UserResource{}
    
    // Core ayarları yap
    r.SetModel(&User{})
    r.SetSlug("users")
    r.SetTitle("Kullanıcılar")
    r.SetIcon("users")
    r.SetGroup("Yönetim")
    r.SetVisible(true)
    r.SetNavigationOrder(1)
    
    // Resolver'ları ayarla
    r.SetFieldResolver(&UserFieldResolver{})
    r.SetCardResolver(&UserCardResolver{})
    r.SetFilterResolver(&UserFilterResolver{})
    r.SetLensResolver(&UserLensResolver{})
    r.SetActionResolver(&UserActionResolver{})
    
    // Yetkilendirme politikasını ayarla
    r.SetPolicy(&UserPolicy{})
    
    return r
}

// UserFieldResolver, kullanıcı alanlarını çözer
type UserFieldResolver struct{}

func (r *UserFieldResolver) ResolveFields(ctx *context.Context) []fields.Element {
    return []fields.Element{
        fields.NewText("name").
            SetLabel("Ad Soyad").
            SetRequired(true),
        
        fields.NewEmail("email").
            SetLabel("E-posta").
            SetRequired(true),
        
        fields.NewText("phone").
            SetLabel("Telefon").
            SetRequired(false),
        
        fields.BelongsTo("Rol", "role_id", "roles").
            DisplayUsing("name"),
    }
}

// UserCardResolver, kullanıcı card'larını çözer
type UserCardResolver struct{}

func (r *UserCardResolver) ResolveCards(ctx *context.Context) []widget.Card {
    return []widget.Card{
        widget.NewValueCard("Toplam Kullanıcı", "1,234"),
        widget.NewTrendCard("Bu Ay", 12, 5),
    }
}

// UserFilterResolver, kullanıcı filtrelerini çözer
type UserFilterResolver struct{}

func (r *UserFilterResolver) ResolveFilters(ctx *context.Context) []resource.Filter {
    return []resource.Filter{
        // Filtreler
    }
}

// UserLensResolver, kullanıcı lens'lerini çözer
type UserLensResolver struct{}

func (r *UserLensResolver) ResolveLenses(ctx *context.Context) []resource.Lens {
    return []resource.Lens{
        // Lens'ler
    }
}

// UserActionResolver, kullanıcı işlemlerini çözer
type UserActionResolver struct{}

func (r *UserActionResolver) ResolveActions(ctx *context.Context) []resource.Action {
    return []resource.Action{
        // İşlemler
    }
}

// User, kullanıcı modeli
type User struct {
    ID    uint
    Name  string
    Email string
    Phone string
    RoleID uint
}
```

## Field Mixin'leri Kullanımı

### Searchable Field

```go
type SearchableTextField struct {
    fields.Base
    fields.Searchable
}

func NewSearchableTextField(key string) *SearchableTextField {
    f := &SearchableTextField{}
    f.SetKey(key)
    f.SetLabel("Aranabilir Alan")
    
    // Aranabilir sütunları ayarla
    f.SetSearchableColumns([]string{"name", "email", "phone"})
    
    // Özel arama callback'i (opsiyonel)
    f.SetSearchCallback(func(db *gorm.DB, term string) *gorm.DB {
        return db.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(term)+"%")
    })
    
    return f
}
```

### Sortable Field

```go
type SortableTextField struct {
    fields.Base
    fields.Sortable
}

func NewSortableTextField(key string) *SortableTextField {
    f := &SortableTextField{}
    f.SetKey(key)
    f.SetLabel("Sıralanabilir Alan")
    
    // Sıralanabilir yap
    f.SetSortable(true)
    f.SetSortDirection("asc")
    
    // Özel sıralama callback'i (opsiyonel)
    f.SetSortCallback(func(db *gorm.DB, direction string) *gorm.DB {
        return db.Order("name " + direction)
    })
    
    return f
}
```

### Filterable Field

```go
type FilterableSelectField struct {
    fields.Base
    fields.Filterable
}

func NewFilterableSelectField(key string) *FilterableSelectField {
    f := &FilterableSelectField{}
    f.SetKey(key)
    f.SetLabel("Filtrelenebilir Alan")
    
    // Filtrelenebilir yap
    f.SetFilterable(true)
    f.SetFilterOptions(map[string]string{
        "active": "Aktif",
        "inactive": "İnaktif",
        "pending": "Beklemede",
    })
    
    // Özel filtreleme callback'i (opsiyonel)
    f.SetFilterCallback(func(db *gorm.DB, value any) *gorm.DB {
        return db.Where("status = ?", value)
    })
    
    return f
}
```

### Validatable Field

```go
type ValidatableEmailField struct {
    fields.Base
    fields.Validatable
}

func NewValidatableEmailField(key string) *ValidatableEmailField {
    f := &ValidatableEmailField{}
    f.SetKey(key)
    f.SetLabel("E-posta")
    
    // Doğrulama kuralları
    f.SetRules([]string{"required", "email", "unique:users,email"})
    
    // Özel doğrulayıcılar
    f.AddValidator(func(value any) error {
        email, ok := value.(string)
        if !ok {
            return fmt.Errorf("email must be a string")
        }
        if len(email) < 5 {
            return fmt.Errorf("email must be at least 5 characters")
        }
        return nil
    })
    
    return f
}
```

### Hideable Field

```go
type HideablePasswordField struct {
    fields.Base
    fields.Hideable
}

func NewHideablePasswordField(key string) *HideablePasswordField {
    f := &HideablePasswordField{}
    f.SetKey(key)
    f.SetLabel("Şifre")
    
    // Gizlilik ayarları
    f.SetShowOnIndex(false)      // Liste görünümünde gösterme
    f.SetShowOnDetail(false)     // Detay görünümünde gösterme
    f.SetShowOnCreate(true)      // Oluşturma formunda göster
    f.SetShowOnUpdate(true)      // Güncelleme formunda göster
    
    return f
}
```

### Kombinasyon

```go
type AdvancedField struct {
    fields.Base
    fields.Searchable
    fields.Sortable
    fields.Filterable
    fields.Validatable
    fields.Displayable
    fields.Hideable
}

func NewAdvancedField(key string) *AdvancedField {
    f := &AdvancedField{}
    f.SetKey(key)
    f.SetLabel("Gelişmiş Alan")
    
    // Searchable
    f.SetSearchableColumns([]string{"name", "email"})
    
    // Sortable
    f.SetSortable(true)
    f.SetSortDirection("asc")
    
    // Filterable
    f.SetFilterable(true)
    f.SetFilterOptions(map[string]string{
        "yes": "Evet",
        "no": "Hayır",
    })
    
    // Validatable
    f.SetRules([]string{"required"})
    
    // Displayable
    f.SetDisplayFormat("uppercase")
    
    // Hideable
    f.SetShowOnIndex(true)
    f.SetShowOnDetail(true)
    f.SetShowOnCreate(true)
    f.SetShowOnUpdate(true)
    
    return f
}
```

## Page Implementasyonu

### Eski Yol

```go
type DashboardPage struct {
    page.Base
}

func (p DashboardPage) Slug() string {
    return "dashboard"
}

func (p DashboardPage) Title() string {
    return "Dashboard"
}

// ... 8+ metod daha
```

### Yeni Yol (Önerilen)

```go
package dashboard

import (
    "github.com/ferdiunal/panel.go/pkg/context"
    "github.com/ferdiunal/panel.go/pkg/fields"
    "github.com/ferdiunal/panel.go/pkg/page"
    "github.com/ferdiunal/panel.go/pkg/widget"
)

// DashboardPage, dashboard sayfasını temsil eder
type DashboardPage struct {
    page.OptimizedBase
}

// NewDashboardPage, yeni bir DashboardPage oluşturur
func NewDashboardPage() *DashboardPage {
    p := &DashboardPage{}
    
    // Core ayarları yap
    p.SetSlug("dashboard")
    p.SetTitle("Dashboard")
    p.SetIcon("layout-dashboard")
    p.SetGroup("Ana")
    p.SetNavigationOrder(0)
    p.SetVisible(true)
    
    // Resolver'ları ayarla
    p.SetFieldResolver(&DashboardFieldResolver{})
    p.SetCardResolver(&DashboardCardResolver{})
    
    return p
}

// DashboardFieldResolver, dashboard alanlarını çözer
type DashboardFieldResolver struct{}

func (r *DashboardFieldResolver) ResolveFields(ctx *context.Context) []fields.Element {
    return []fields.Element{
        // Dashboard alanları
    }
}

// DashboardCardResolver, dashboard card'larını çözer
type DashboardCardResolver struct{}

func (r *DashboardCardResolver) ResolveCards(ctx *context.Context) []widget.Card {
    return []widget.Card{
        widget.NewValueCard("Toplam Kullanıcı", "1,234"),
        widget.NewValueCard("Toplam Siparişler", "567"),
        widget.NewTrendCard("Bu Ay", 12, 5),
        widget.NewTrendCard("Geçen Ay", 8, -2),
    }
}
```

## Best Practices

### 1. Resolver'ları Ayrı Dosyalarda Tutun

```
pkg/resource/user/
├── resource.go          # UserResource
├── field_resolver.go    # UserFieldResolver
├── card_resolver.go     # UserCardResolver
├── filter_resolver.go   # UserFilterResolver
├── lens_resolver.go     # UserLensResolver
└── action_resolver.go   # UserActionResolver
```

### 2. Resolver'ları Lazy Load Edin

```go
func (r *UserResource) ResolveFields(ctx *context.Context) []fields.Element {
    if r.fieldResolver == nil {
        r.fieldResolver = &UserFieldResolver{}
    }
    return r.fieldResolver.ResolveFields(ctx)
}
```

### 3. Context Kullanın

```go
func (r *UserFieldResolver) ResolveFields(ctx *context.Context) []fields.Element {
    fields := []fields.Element{
        fields.NewText("name"),
    }
    
    // Context'ten kullanıcı bilgisini al
    if user, ok := ctx.Get("user"); ok {
        // Kullanıcıya göre alanları özelleştir
        if user.IsAdmin {
            fields = append(fields, fields.NewText("admin_notes"))
        }
    }
    
    return fields
}
```

### 4. Caching Kullanın

```go
type UserFieldResolver struct {
    cache []fields.Element
}

func (r *UserFieldResolver) ResolveFields(ctx *context.Context) []fields.Element {
    if r.cache != nil {
        return r.cache
    }
    
    r.cache = []fields.Element{
        fields.NewText("name"),
        fields.NewEmail("email"),
    }
    
    return r.cache
}
```

## Migration Guide

### Adım 1: Resolver'ları Oluştur

```go
type UserFieldResolver struct{}

func (r *UserFieldResolver) ResolveFields(ctx *context.Context) []fields.Element {
    // Eski Fields() metodundan kodu buraya taşı
    return []fields.Element{
        fields.NewText("name"),
        fields.NewEmail("email"),
    }
}
```

### Adım 2: Resource'ı Güncelle

```go
type UserResource struct {
    resource.OptimizedBase
}

func NewUserResource() *UserResource {
    r := &UserResource{}
    r.SetModel(&User{})
    r.SetSlug("users")
    r.SetTitle("Kullanıcılar")
    r.SetFieldResolver(&UserFieldResolver{})
    return r
}
```

### Adım 3: Eski Metodları Kaldır

Eski `Fields()`, `Cards()` vb. metodları kaldır.

## Troubleshooting

### Problem: Resolver nil döndürüyor

**Çözüm**: Resolver'ı ayarladığından emin ol

```go
r.SetFieldResolver(&UserFieldResolver{})
```

### Problem: Alanlar görünmüyor

**Çözüm**: `ResolveFields()` metodunun çağrıldığından emin ol

```go
func (b *OptimizedBase) Fields() []fields.Element {
    return b.ResolveFields(nil)
}
```

### Problem: Context nil

**Çözüm**: Context'i kontrol et

```go
func (r *UserFieldResolver) ResolveFields(ctx *context.Context) []fields.Element {
    if ctx == nil {
        // Varsayılan alanları döndür
        return []fields.Element{}
    }
    // Context'i kullan
    return []fields.Element{}
}
```

## Kaynaklar

- [Resource Documentation](./Resources.md)
- [Fields Documentation](./Fields.md)
- [Pages Documentation](./Pages.md)
- [Relationships Documentation](./Relationships.md)
