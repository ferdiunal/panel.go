# API ve UI Entegrasyonu

Panel SDK, frontend (UI) geliştiricileri için öngörülebilir ve zengin bir JSON yanıt yapısı sunar. Bu yapı, dinamik arayüzler oluşturmak için gerekli tüm metaları (meta verileri) içerir.

## Standart Yanıt Yapısı

Tüm kaynak (Resource) istekleri (`Index`, `Show`, `Store`, `Update`) standart bir format izler:

```json
{
    "data": [ ... ] veya { ... },
    "meta": {
        "title": "Sayfa Başlığı",
        "policy": { ... },
        "current_page": 1,
        "per_page": 10,
        "total": 50
    }
}
```

### 1. Meta Verisi (`meta`)

`meta` objesi, UI'ın durumunu yönetmek için kritik bilgiler içerir.

*   **title**: Sayfanın veya kaynağın görünen adı (Örn: "Kullanıcılar"). Header veya Breadcrumb için kullanılır.
*   **policy**: Kullanıcının bu kaynak üzerindeki yetkilerini belirtir. UI butonlarını (Örn: "Sil", "Düzenle", "Yeni Ekle") gizlemek/göstermek için kullanılır.

**Örnek Policy (Index - Liste Sayfası):**
```json
"policy": {
    "create": true,   // "Yeni Ekle" butonu gösterilmeli mi?
    "view_any": true  // Listeyi görme yetkisi var mı?
}
```

**Örnek Policy (Show/Detail - Detay Sayfası):**
```json
"policy": {
    "view": true,     // Detayı görme yetkisi
    "update": false,  // "Düzenle" butonu gizlenmeli
    "delete": true    // "Sil" butonu gösterilmeli
}
```

### 2. Veri (`data`)

#### Liste Görünümü (Index)
`data`, her biri bir kayıt temsil eden objelerden oluşan bir dizidir. Her kayıt, `key` - `value` (alan tanımı) çiftlerinden oluşur.

```json
"data": [
    {
        "id": { "view": "id-field", "value": 1, ... },
        "name": { "view": "text-field", "value": "Ferdi", ... },
        "email": { "view": "email-field", "value": "test@example.com", ... }
    }
]
```

#### Detay Görünümü (Show)
`data`, tek bir kaydı temsil eden objedir.

## Widget API

Widget verilerine erişmek için aşağıdaki yapı kullanılır:

### Liste (`GET /api/resource/:resource/widgets`)

Available widget'ların listesini döner.

```json
{
    "data": [
        {
            "index": 0,
            "name": "Toplam Kullanıcı",
            "type": "value"
        }
    ]
}
```

### Detay (`GET /api/resource/:resource/widgets/:index`)

Belirli bir widget'ın verisini döner.

```json
{
    "data": 150 // Value widget örneği
}
```

## Alan Yapısı (Field Structure)

Her bir alan (field), sadece değerini değil, nasıl görüntüleneceğine dair tüm yapılandırmayı içerir. Frontend, bu yapıya göre dinamik component render etmelidir.

```json
{
    "view": "text-field",       // Component türü (Vue/React tarafında eşleşecek)
    "key": "name",              // Modeldeki alan adı
    "name": "Ad Soyad",         // Label (Etiket)
    "data": "Ferdi Ünal",       // Değer
    "placeholder": "Adınızı girin",
    "required": true,
    "read_only": false,
    "sortable": true,
    "props": {                  // Ekstra özellikler (varsa)
        "size": "large"
    }
}
```

### UI Entegrasyon İpuçları

1.  **Dinamik Componentler**: `view` özelliğini kullanarak (örn: `text-field`, `select-field`) frontend tarafında dinamik component yükleyici kullanın.
2.  **Yetki Kontrolü**: `meta.policy` objesini kontrol ederek aksiyon butonlarını (Sil, Düzenle) engelleyin.
3.  **Başlıklar**: Sayfa başlığını `meta.title`'dan okuyun.
