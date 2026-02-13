# Relationship Loading Implementation

Bu dokümantasyon, Panel.go'daki relationship loading sisteminin teknik implementasyonunu açıklar. Geliştiriciler için mimari kararlar, GORM entegrasyonu ve implementation detayları hakkında bilgi sağlar.

## İçindekiler

- [Genel Bakış](#genel-bakış)
- [Mimari](#mimari)
- [GORM Geçişi](#gorm-geçişi)
- [Lazy Loading](#lazy-loading)
- [Eager Loading](#eager-loading)
- [İlişki Türleri](#ilişki-türleri)
- [Performans Optimizasyonları](#performans-optimizasyonları)
- [Best Practices](#best-practices)
- [Sorun Giderme](#sorun-giderme)

## Genel Bakış

Relationship loading sistemi, veritabanı ilişkilerini yüklemek için iki strateji sunar:

1. **Eager Loading**: İlişkili verileri önceden yükler (N+1 query problemini önler)
2. **Lazy Loading**: İlişkili verileri ihtiyaç anında yükler

Her iki strateji de GORM'un native API'lerini kullanır:
- **Lazy Loading**: `Association().Find()` ile reflection-based dinamik yükleme
- **Eager Loading**: `Where()` ve `Find()` ile batch loading

## Mimari

### Dosya Yapısı

```
pkg/data/
├── relationship_loader.go           # Interface tanımları
├── relationship_loader_impl.go      # Ana loader implementasyonu
├── relationship_strategies.go       # Yükleme stratejileri (lazy/eager)
└── gorm_provider.go                 # GORM provider entegrasyonu
```

### Ana Bileşenler

#### 1. RelationshipLoader Interface

```go
type RelationshipLoader interface {
    LoadRelationships(ctx context.Context, items []interface{}, fields []fields.RelationshipField) error
    LoadRelationship(ctx context.Context, item interface{}, field fields.RelationshipField) (interface{}, error)
}
```

#### 2. GormRelationshipLoader

GORM kullanarak relationship loading implementasyonu:

```go
type GormRelationshipLoader struct {
    db *gorm.DB
}
```

**Metodlar:**
- `LoadRelationships`: Batch loading (eager loading)
- `LoadRelationship`: Single item loading (lazy loading)

#### 3. Yükleme Stratejileri

Her ilişki türü için iki metod:
- `eagerLoad{Type}`: Batch loading için
- `lazyLoad{Type}`: Single item loading için

**Desteklenen İlişki Türleri:**
- BelongsTo
- HasMany
- HasOne
- BelongsToMany

## GORM Geçişi

### Önceki Durum (Raw SQL)

Önceden raw SQL sorguları kullanılıyordu:

```go
// Eski yaklaşım
query := fmt.Sprintf("SELECT * FROM %s WHERE %s IN ?", table, column)
db.Raw(query, ids).Scan(&results)
```

**Problemler:**
- SQL injection riski
- GORM'un query builder özelliklerinden faydalanamama
- Manuel string manipülasyonu
- Hata ayıklama zorluğu

### Yeni Durum (GORM Native API)

Şimdi GORM'un native API'leri kullanılıyor:

```go
// Yeni yaklaşım - Lazy Loading
db.Model(item).Association(relationshipName).Find(result)

// Yeni yaklaşım - Eager Loading
db.Table(table).Where(column+" IN ?", ids).Find(&results)
```

**Avantajlar:**
- ✅ SQL injection koruması
- ✅ GORM'un query builder özellikleri
- ✅ Tip güvenliği
- ✅ Daha iyi hata mesajları
- ✅ GORM middleware desteği (hooks, callbacks)
- ✅ Daha kolay test edilebilirlik

### Migration Adımları

1. **Lazy Loading Metodları** → GORM Association API
2. **Eager Loading Metodları** → GORM Query Builder
3. **String Concatenation** → `fmt.Sprintf` yerine `+` operatörü
4. **Dokümantasyon** → Yorumları güncelleme

## Lazy Loading

Lazy loading, tek bir kayıt için ilişkiyi yükler. GORM'un `Association` API'sini kullanır.

### Implementation

```go
func (l *GormRelationshipLoader) lazyLoadBelongsTo(ctx context.Context, item interface{}, field fields.RelationshipField) (interface{}, error) {
    if item == nil {
        return nil, nil
    }

    // Reflection kullanarak relationship field tipini al
    itemValue := reflect.ValueOf(item)
    if itemValue.Kind() == reflect.Ptr {
        itemValue = itemValue.Elem()
    }

    relField := itemValue.FieldByName(field.GetRelationshipName())
    if !relField.IsValid() {
        return nil, fmt.Errorf("relationship field %s not found", field.GetRelationshipName())
    }

    relType := relField.Type()

    // Yeni instance oluştur
    var relValue reflect.Value
    if relType.Kind() == reflect.Ptr {
        relValue = reflect.New(relType.Elem())
    } else {
        relValue = reflect.New(relType)
    }

    // GORM Association API kullanarak ilişkiyi yükle
    err := l.db.WithContext(ctx).
        Model(item).
        Association(field.GetRelationshipName()).
        Find(relValue.Interface())

    if err != nil {
        return nil, fmt.Errorf("failed to load BelongsTo relationship: %w", err)
    }

    // İlişki verisini set et
    var actualValue interface{}
    if relType.Kind() == reflect.Ptr {
        actualValue = relValue.Interface()
    } else {
        actualValue = relValue.Elem().Interface()
    }

    if err := setRelationshipData(item, field.GetRelationshipName(), actualValue); err != nil {
        return nil, err
    }

    return actualValue, nil
}
```

### Özellikler

1. **Reflection-Based**: Dinamik olarak struct tipini belirler
2. **Type-Safe**: Pointer ve non-pointer tipleri destekler
3. **Context-Aware**: Context propagation için `WithContext` kullanır
4. **Error Handling**: Detaylı hata mesajları

### Kullanım Senaryoları

- Tek bir kaydın detay sayfası
- API endpoint'lerinde tek kayıt dönüşü
- İhtiyaç anında yükleme

## Eager Loading

Eager loading, birden fazla kayıt için ilişkileri batch olarak yükler. N+1 query problemini önler.

### Implementation

```go
func (l *GormRelationshipLoader) eagerLoadBelongsTo(ctx context.Context, items []interface{}, field fields.RelationshipField) error {
    // BelongsTo field'ından gerekli bilgileri al
    belongsToField, ok := field.(*fields.BelongsToField)
    if !ok {
        return fmt.Errorf("field is not a BelongsTo field")
    }

    foreignKey := belongsToField.GetForeignKey()
    ownerKey := belongsToField.GetOwnerKeyColumn()
    relatedTable := belongsToField.GetRelatedTableName()

    // 1. Tüm item'lardan foreign key değerlerini çıkar
    foreignKeyValues := []interface{}{}
    itemsByForeignKey := map[interface{}][]interface{}{}

    for _, item := range items {
        fkValue := extractFieldValue(item, foreignKey)
        if fkValue != nil && !isZeroValue(fkValue) {
            foreignKeyValues = append(foreignKeyValues, fkValue)
            itemsByForeignKey[fkValue] = append(itemsByForeignKey[fkValue], item)
        }
    }

    if len(foreignKeyValues) == 0 {
        return nil // Hiç foreign key yok
    }

    // 2. GORM query builder ile ilişkili kayıtları yükle
    safeOwnerKey := SanitizeColumnName(ownerKey)
    safeTable := SanitizeColumnName(relatedTable)

    var relatedRecords []map[string]interface{}
    err := l.db.WithContext(ctx).
        Table(safeTable).
        Where(safeOwnerKey+" IN ?", foreignKeyValues).
        Find(&relatedRecords).Error

    if err != nil {
        return fmt.Errorf("failed to load BelongsTo relationship: %w", err)
    }

    // 3. İlişkili kayıtları ID'ye göre map et
    relatedByID := map[interface{}]map[string]interface{}{}
    for _, record := range relatedRecords {
        id := record[ownerKey]
        relatedByID[id] = record
    }

    // 4. Her item'a ilişkili kaydı set et
    for fkValue, itemList := range itemsByForeignKey {
        relatedRecord := relatedByID[fkValue]
        if relatedRecord != nil {
            for _, item := range itemList {
                if err := setRelationshipData(item, field.GetRelationshipName(), relatedRecord); err != nil {
                    fmt.Printf("[WARN] Failed to set BelongsTo relationship: %v\n", err)
                }
            }
        }
    }

    return nil
}
```

### İşlem Sırası

1. **Extraction**: Tüm item'lardan foreign key değerlerini çıkar
2. **Batch Query**: Tek sorguda tüm ilişkili kayıtları yükle
3. **Mapping**: İlişkili kayıtları ID'ye göre map et
4. **Assignment**: Her item'a ilişkili kaydı set et

### Performans Optimizasyonları

1. **Batch Loading**: N+1 query yerine 2 query (1 ana + 1 ilişki)
2. **Memory Mapping**: Hash map kullanarak O(1) lookup
3. **Context Propagation**: Timeout ve cancellation desteği
4. **Sanitization**: SQL injection koruması

### Kullanım Senaryoları

- Liste sayfaları
- API endpoint'lerinde çoklu kayıt dönüşü
- Performans kritik senaryolar

## İlişki Türleri

### BelongsTo

**Lazy Loading:**
```go
// GORM Association API kullanır
db.Model(item).Association("Author").Find(&author)
```

**Eager Loading:**
```go
// Batch loading: Tüm foreign key'leri topla ve tek sorguda yükle
db.Table("authors").Where("id IN ?", authorIDs).Find(&authors)
```

**Dosya:** `pkg/data/relationship_strategies.go:115-166` (lazy), `pkg/data/relationship_strategies.go:32-98` (eager)

### HasMany

**Lazy Loading:**
```go
// GORM Association API kullanır
db.Model(item).Association("Posts").Find(&posts)
```

**Eager Loading:**
```go
// Batch loading: Tüm owner key'leri topla ve tek sorguda yükle
db.Table("posts").Where("author_id IN ?", authorIDs).Find(&posts)
```

**Dosya:** `pkg/data/relationship_strategies.go:271-312` (lazy), `pkg/data/relationship_strategies.go:189-254` (eager)

### HasOne

**Lazy Loading:**
```go
// GORM Association API kullanır
db.Model(item).Association("Profile").Find(&profile)
```

**Eager Loading:**
```go
// Batch loading: Tüm owner key'leri topla ve tek sorguda yükle
db.Table("profiles").Where("user_id IN ?", userIDs).Find(&profiles)
```

**Dosya:** `pkg/data/relationship_strategies.go:419-470` (lazy), `pkg/data/relationship_strategies.go:335-402` (eager)

### BelongsToMany

**Lazy Loading:**
```go
// GORM Association API kullanır
db.Model(item).Association("Roles").Find(&roles)
```

**Eager Loading:**
```go
// 3 adımlı batch loading:
// 1. Pivot tablodan ilişkileri çek
db.Table("role_user").Where("user_id IN ?", userIDs).Find(&pivots)

// 2. İlişkili kayıt ID'lerini çıkar
relatedIDs := extractRelatedIDs(pivots)

// 3. İlişkili kayıtları yükle
db.Table("roles").Where("id IN ?", relatedIDs).Find(&roles)
```

**Dosya:** `pkg/data/relationship_strategies.go:819-860` (lazy), `pkg/data/relationship_strategies.go:686-802` (eager)

## Performans Optimizasyonları

### 1. N+1 Query Problemi Çözümü

**Problem:**
```go
// Her kayıt için ayrı query (N+1 query)
for _, post := range posts {
    db.Where("id = ?", post.AuthorID).First(&post.Author)
}
// Toplam: 1 (posts) + N (authors) = N+1 query
```

**Çözüm:**
```go
// Eager loading ile 2 query
// 1. Ana kayıtları yükle
db.Find(&posts)

// 2. Tüm ilişkili kayıtları tek sorguda yükle
authorIDs := extractAuthorIDs(posts)
db.Where("id IN ?", authorIDs).Find(&authors)
// Toplam: 2 query
```

### 2. Memory Mapping

Hash map kullanarak O(1) lookup:

```go
// İlişkili kayıtları ID'ye göre map et
relatedByID := map[interface{}]map[string]interface{}{}
for _, record := range relatedRecords {
    id := record[ownerKey]
    relatedByID[id] = record
}

// O(1) lookup
for _, item := range items {
    related := relatedByID[item.ForeignKey]
}
```

### 3. Context Propagation

Timeout ve cancellation desteği:

```go
err := l.db.WithContext(ctx).
    Table(table).
    Where(column+" IN ?", ids).
    Find(&results).Error
```

### 4. Batch Size Optimization

Büyük veri setleri için batch processing:

```go
// Örnek: 1000'er kayıt olarak işle
batchSize := 1000
for i := 0; i < len(items); i += batchSize {
    end := i + batchSize
    if end > len(items) {
        end = len(items)
    }
    batch := items[i:end]
    // Batch'i işle
}
```

## Best Practices

### 1. Eager Loading Kullanımı

**✅ İyi:**
```go
field := BelongsTo("Author", "author_id", "authors").
    WithEagerLoad()
```

**❌ Kötü:**
```go
field := BelongsTo("Author", "author_id", "authors").
    WithLazyLoad()
```

**Neden:** Liste sayfalarında N+1 query problemini önler.

### 2. Context Kullanımı

**✅ İyi:**
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

loader.LoadRelationships(ctx, items, fields)
```

**❌ Kötü:**
```go
loader.LoadRelationships(context.Background(), items, fields)
```

**Neden:** Timeout ve cancellation desteği sağlar.

### 3. Error Handling

**✅ İyi:**
```go
if err := loader.LoadRelationships(ctx, items, fields); err != nil {
    return fmt.Errorf("failed to load relationships: %w", err)
}
```

**❌ Kötü:**
```go
loader.LoadRelationships(ctx, items, fields)
```

**Neden:** Hataları yakalamak ve loglamak önemlidir.

### 4. Null Check

**✅ İyi:**
```go
if item == nil {
    return nil, nil
}
```

**❌ Kötü:**
```go
// Null check yok
relValue := reflect.ValueOf(item)
```

**Neden:** Nil pointer dereference'ı önler.

### 5. Type Safety

**✅ İyi:**
```go
belongsToField, ok := field.(*fields.BelongsToField)
if !ok {
    return fmt.Errorf("field is not a BelongsTo field")
}
```

**❌ Kötü:**
```go
belongsToField := field.(*fields.BelongsToField)
```

**Neden:** Type assertion panic'ini önler.

## Sorun Giderme

### Problem 1: N+1 Query

**Belirti:**
```
[SQL] SELECT * FROM posts WHERE id = 1
[SQL] SELECT * FROM authors WHERE id = 1
[SQL] SELECT * FROM authors WHERE id = 2
[SQL] SELECT * FROM authors WHERE id = 3
...
```

**Çözüm:**
```go
// WithEagerLoad() kullan
field := BelongsTo("Author", "author_id", "authors").
    WithEagerLoad()
```

### Problem 2: İlişki Yüklenmiyor

**Belirti:**
```go
// İlişki field'ı nil veya boş
post.Author == nil
```

**Çözüm:**
1. Field adının doğru olduğunu kontrol et
2. Foreign key'in doğru olduğunu kontrol et
3. İlişkili kaydın veritabanında olduğunu kontrol et

```go
// Debug için log ekle
fmt.Printf("Loading relationship: %s\n", field.GetRelationshipName())
fmt.Printf("Foreign key: %s\n", field.GetForeignKey())
```

### Problem 3: Reflection Hatası

**Belirti:**
```
panic: reflect: call of reflect.Value.FieldByName on zero Value
```

**Çözüm:**
```go
// Nil check ekle
if item == nil {
    return nil, nil
}

// Pointer check ekle
itemValue := reflect.ValueOf(item)
if itemValue.Kind() == reflect.Ptr {
    if itemValue.IsNil() {
        return nil, nil
    }
    itemValue = itemValue.Elem()
}
```

### Problem 4: SQL Injection

**Belirti:**
```go
// Güvenli değil
query := fmt.Sprintf("SELECT * FROM %s WHERE %s IN (%s)", table, column, ids)
```

**Çözüm:**
```go
// GORM query builder kullan
db.Table(table).Where(column+" IN ?", ids).Find(&results)

// Veya sanitize et
safeTable := SanitizeColumnName(table)
safeColumn := SanitizeColumnName(column)
```

### Problem 5: Memory Leak

**Belirti:**
```
Memory usage increasing over time
```

**Çözüm:**
```go
// Map'leri temizle
defer func() {
    relatedByID = nil
    itemsByForeignKey = nil
}()

// Context timeout kullan
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
```

## Gelecek Geliştirmeler

### 1. Preload Desteği

GORM'un `Preload` metodunu kullanarak daha native bir yaklaşım:

```go
// Şu anki yaklaşım
db.Table("posts").Where("id IN ?", ids).Find(&posts)
db.Table("authors").Where("id IN ?", authorIDs).Find(&authors)

// Gelecek yaklaşım
db.Preload("Author").Find(&posts)
```

### 2. Nested Relationships

İç içe ilişkileri destekleme:

```go
// Post -> Author -> Profile
field := BelongsTo("Author", "author_id", "authors").
    WithEagerLoad().
    WithNested(
        HasOne("Profile", "profile", "profiles"),
    )
```

### 3. Conditional Loading

Koşullu yükleme desteği:

```go
field := BelongsTo("Author", "author_id", "authors").
    LoadIf(func(item interface{}) bool {
        return item.(*Post).Status == "published"
    })
```

### 4. Caching

İlişki verilerini cache'leme:

```go
field := BelongsTo("Author", "author_id", "authors").
    WithCache(5 * time.Minute)
```

### 5. Polymorphic Relationships

Polymorphic ilişkileri destekleme:

```go
field := MorphTo("Commentable", "commentable").
    Types(map[string]string{
        "posts":  "posts",
        "videos": "videos",
    })
```

## Referanslar

- [GORM Documentation](https://gorm.io/docs/)
- [GORM Association](https://gorm.io/docs/associations.html)
- [GORM Query Builder](https://gorm.io/docs/query.html)
- [Go Reflection](https://go.dev/blog/laws-of-reflection)
- [Panel.go Relationships](./Relationships.md)

## Değişiklik Geçmişi

### 2026-02-08: GORM Geçişi

- ✅ Lazy loading metodları GORM Association API kullanacak şekilde güncellendi
- ✅ Eager loading metodları GORM Query Builder kullanacak şekilde güncellendi
- ✅ String concatenation ile SQL injection koruması eklendi
- ✅ Dokümantasyon güncellendi

**Değiştirilen Dosyalar:**
- `pkg/data/relationship_strategies.go`
- `pkg/data/relationship_loader_impl.go`
- `pkg/data/gorm_provider.go`

**Değişiklik Detayları:**
- Raw SQL sorguları GORM native API'leri ile değiştirildi
- `fmt.Sprintf` yerine string concatenation kullanıldı
- Tüm yorumlar "Raw SQL" yerine "GORM query builder" olarak güncellendi
- Reflection-based dinamik tip belirleme eklendi
- Context propagation desteği eklendi
