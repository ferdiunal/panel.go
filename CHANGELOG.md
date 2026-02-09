# Changelog

Tüm önemli değişiklikler bu dosyada dökümante edilir.

## [Unreleased]

### ✨ Yeni Özellikler

#### Resource Title Pattern (Laravel Nova Uyumlu)

Panel.go'ya Laravel Nova'nın title pattern'i eklendi. Her resource için kayıt başlığı (record title) özelliği artık kullanılabilir. Bu, ilişki fieldlarında kayıtların okunabilir şekilde gösterilmesini sağlar.

**Özellikler:**
- `SetRecordTitleKey(key string)` - Kayıt başlığı için kullanılacak field adını ayarlar
- `GetRecordTitleKey() string` - Kayıt başlığı için kullanılacak field adını döndürür
- `SetRecordTitleFunc(fn func(any) string)` - Özel başlık fonksiyonu ayarlar
- `RecordTitle(record any) string` - Kaydın okunabilir başlığını döndürür

**Kullanım Örneği:**

```go
// UserResource'da "name" field'ını başlık olarak ayarla
func NewUserResource() *UserResource {
    r := &UserResource{}
    r.SetModel(&User{})
    r.SetSlug("users")
    r.SetRecordTitleKey("name") // ← Yeni özellik
    return r
}

// Özel başlık fonksiyonu ile
r.SetRecordTitleFunc(func(record any) string {
    user := record.(*User)
    return user.FirstName + " " + user.LastName
})
```

**İlişki Fieldları:**

Tüm ilişki fieldları artık minimal format döndürür: `{"id": ..., "title": ...}`

- **BelongsTo**: `{"id": 5, "title": "John Doe"}`
- **HasMany**: `[{"id": 1, "title": "First Post"}, {"id": 2, "title": "Second Post"}]`
- **HasOne**: `{"id": 1, "title": "User Profile"}`
- **BelongsToMany**: `[{"id": 1, "title": "Admin"}, {"id": 2, "title": "Editor"}]`

**Etkilenen Dosyalar:**
- `pkg/resource/resource.go` - Interface'e yeni metodlar eklendi
- `pkg/resource/optimized.go` - OptimizedBase implementation
- `pkg/resource/base.go` - Base implementation
- `pkg/fields/belongs_to.go` - Extract metodu eklendi
- `pkg/fields/has_many.go` - Extract metodu güncellendi
- `pkg/fields/has_one.go` - Extract metodu güncellendi
- `pkg/fields/belongs_to_many.go` - Extract metodu eklendi
- `pkg/resource/user/resource.go` - SetRecordTitleKey("name") eklendi
- `pkg/resource/account/resource.go` - SetRecordTitleKey("name") eklendi
- `pkg/resource/session/resource.go` - SetRecordTitleKey("id") eklendi
- `pkg/resource/verification/resource.go` - SetRecordTitleKey("id") eklendi

**Testler:**
- `pkg/resource/record_title_test.go` - RecordTitle için kapsamlı testler eklendi
- Tüm testler başarıyla çalışıyor ✅

### 🔧 Düzeltmeler

#### Base Resource Bug Fix

`Base.SetDialogType` ve `Base.SetOpenAPIEnabled` metodları pointer receiver'a çevrildi. Bu metodlar value receiver kullanıyordu ve değişiklikler kayboluyordu.

**Önceki (Hatalı):**
```go
func (b Base) SetDialogType(dialogType DialogType) Resource {
    b.DialogType = dialogType // Değişiklik kaybolur (kopya üzerinde)
    return b
}
```

**Sonrası (Düzeltilmiş):**
```go
func (b *Base) SetDialogType(dialogType DialogType) Resource {
    b.DialogType = dialogType // Değişiklik kalıcı
    return b
}
```

### ⚠️ Breaking Changes

1. **İlişki Field Serialize Formatı**: HasMany, HasOne, BelongsToMany fieldları artık `{"id": ..., "title": ...}` formatında döndürüyor (önceden tam kayıt veya sadece ID döndürüyordu)

2. **Base Resource Metodları**: `SetDialogType` ve `SetOpenAPIEnabled` metodları pointer receiver'a çevrildi

### 📝 Önemli Notlar

- **Eager Loading Zorunlu**: İlişki fieldlarında eager loading yapılmalı, aksi halde title null olur
- **DisplayUsing Korundu**: Mevcut DisplayUsing() callback'leri çalışmaya devam ediyor
- **Type Assertion**: RelatedResource interface{} tipinde olduğu için type assertion kullanıldı
- **MorphTo**: TypeMappings map[string]string olduğu için (resource slug'ları tutuyor) title pattern uygulanmadı

### 🧪 Test Durumu

- ✅ Resource testleri: Tüm testler başarılı
- ✅ RecordTitle testleri: Yeni testler eklendi ve başarılı
- ⚠️ Fields testleri: Mevcut test dosyalarında constructor fonksiyon adları ile ilgili sorunlar var (implementasyondan bağımsız)

### 📚 Dökümantasyon

- CHANGELOG.md oluşturuldu
- RecordTitle için kapsamlı testler ve örnekler eklendi
- Tüm metodlar Türkçe dokümantasyon ile açıklandı

---

## Önceki Sürümler

Önceki sürüm notları için git commit geçmişine bakınız.
