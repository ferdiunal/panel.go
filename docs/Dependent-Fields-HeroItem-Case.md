# Bağımlı Alanlar Vaka Notu: Hero Item (`target_type`)

Bu doküman, gerçek bir implementasyon sırasında karşılaşılan sorunları ve önerilen çözüm pattern'ini özetler.

Kapsam:
- `DependsOn("target_type")` ile koşullu alan gösterme
- Edit formunda `target_type` seçili gelmeme problemi
- Değer normalize etme ve fallback stratejisi
- Image upload ile birlikte görülen sürüm/tutum sorunları

## Problem Özeti

`HeroItem` formunda `product_id`, `category_id`, `static_url` alanları `target_type` değerine göre açılıp kapanıyor.

Sahada görülen belirtiler:
1. `resolve-dependencies` callback'i çalışıyor ama edit formda `target_type` select boş görünüyor.
2. Bağımlı alan (ör. `static_url`) doğru görünür olsa bile select placeholder'da kalıyor.
3. Bazı payload'larda `target_type` düz string yerine envelope/map olarak geliyor.

## Kök Nedenler

1. `target_type` her zaman `string` gelmiyor.
2. Edit response başlangıç değerinde `target_type` boş/uyumsuz olabilir.
3. Callback içinde sadece `formData["target_type"].(string)` okumak yetersiz kalıyor.

## Önerilen Pattern

Dosya: `internal/resource/heroitem/heroitem_field_resolver.go`

### 1) Parent alanı normalize et

- `normalizeTargetType(value interface{}) string`
- Desteklenecek tipler:
  - `string`
  - `*string`
  - `map[string]interface{}` (`data`, `value`, `target_type` anahtarları)
  - JSON string (`{"data":"static_url"}`)

### 2) Form fallback kuralı ekle

`targetTypeFromForm(formData)` içinde:
1. Önce `target_type` normalize et.
2. Boşsa aşağıdaki infer fallback'lerini uygula:
   - `product_id` dolu -> `product`
   - `category_id` dolu -> `category`
   - `static_url` dolu -> `static_url`

### 3) Edit başlangıç değerini `Resolve(...)` ile garanti et

`targetTypeField.Resolve(...)` içinde:
1. `value` normalize et.
2. Boşsa item/model içinden infer et (`TargetType`, `ProductID`, `CategoryID`, `StaticURL`).

### 4) Submit tutarlılığı için `Modify(...)` kullan

`targetTypeField.Modify(...)` içinde:
1. Mevcut değeri normalize et.
2. Hala boşsa form field'larından fallback ile tekrar üret.

## Bağımlı Alan Callback Standardı

Her bağımlı alan için aynı yaklaşım:
- aktif durum: `Show().Enable().MakeRequired()`
- pasif durum: `Hide().Disable().MakeOptional().SetValue(nil)`

Bu pattern eski/yanlış değerin formda taşınmasını engeller.

## Tip Güvenliği Notları

Relation ID alanlarında tip karışıklığı çok yaygın:
- frontend: string
- backend/body parser: string/float64
- model: `*uint`

Bu yüzden helper kullanın:
- `toNullableUint(...)`
- `hasPositiveUint(...)`

Pointer (`*uint`, `*string`) tiplerini de normalize katmanında destekleyin.

## Debug Checklist

1. `GET /resource/{resource}/{id}/edit` response:
   - `target_type.data` gerçekten dolu mu?
2. Form field JSON:
   - bağımlı alanlarda `depends_on: ["target_type"]` var mı?
3. `POST /resource/{resource}/fields/resolve-dependencies` request:
   - `changedFields` içinde `target_type` var mı?
   - `formData.target_type` beklenen formatta mı?
4. Resolve response:
   - doğru alanda `visible/disabled/required/value` dönüyor mu?
5. Sunucu restart:
   - resolver değişikliği sonrası zorunlu.
6. Frontend cache:
   - hard refresh sonrası tekrar test.

## Image Upload ile Birlikte Görülen Durum

`HeroItemResource`, `resource.OptimizedBase` kullanıyorsa upload davranışı panel.go sürümüne bağlıdır.

Kontrol:
- `OptimizedBase.StoreHandler(...)` gerçek dosya kaydı yapıyor mu?
- Projede `replace github.com/ferdiunal/panel.go => /local/path` varsa, çalışan kodun gerçekten güncel olduğundan emin olun.

## Kullanım Kararı

Ne zaman bu yaklaşımı seçmeliyim?
- Parent alan value'su farklı payload şekillerinde gelebiliyorsa
- Edit başlangıç state'i create ile aynı formatta değilse
- Business rule "tek tip aktif alan" ise (`product/category/static_url`)

Ne zaman gerekmez?
- Parent alan her zaman düz string ve form/create/edit payloadları tamamen tutarlıysa

## Kısa Sonuç

Bağımlı alanların stabil çalışması için yalnızca `DependsOn + OnDependencyChange` yetmez.
Üretimde güvenli çözüm:
1. parent value normalize,
2. edit başlangıç infer,
3. submit fallback,
4. pasif alanı `SetValue(nil)` ile temizle.

