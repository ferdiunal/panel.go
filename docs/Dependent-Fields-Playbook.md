# Bağımlı Alanlar Playbook (dependent_on)

Bu doküman, `DependsOn()` / `OnDependencyChange()` akışını hızlı ve sorunsuz kurmak için hazırlanmıştır.

## Ne zaman kullanılır?

- Bir alanın seçenekleri başka bir alanın değerine göre değişecekse
- Bir alan görünür/zorunlu/disabled durumu başka alana bağlıysa
- Cascade seçim (örn: ürün -> opsiyon -> opsiyon değeri) gerekiyorsa

## Kısa Mimari

1. Backend field tanımı `DependsOn("alan")` ile bağımlılığı tanımlar.
2. Frontend, initial field response içindeki `depends_on` bilgisini okuyup değişimi izler.
3. Değişiklikte `/api/resource/:resource/fields/resolve-dependencies` endpoint'ine istek atılır.
4. Backend callback (`OnDependencyChange`) `FieldUpdate` döner.
5. Frontend update'i ilgili alana uygular (`options`, `disabled`, `helpText`, `value` vb.).

## Minimum Doğru Kurulum

```go
optionField := fields.Select("Product Option", "variant_option_id")
optionField.OnlyOnForm()
optionField.DependsOn("product_id")
optionField.OnDependencyChange(func(
	field *fields.Schema,
	formData map[string]interface{},
	c *fiber.Ctx,
) *fields.FieldUpdate {
	productID := toUint(formData["product_id"])
	opts := loadOptionsByProduct(productID, c)

	update := fields.NewFieldUpdate().
		SetOptions(opts).
		SetValue(nil)

	if len(opts) == 0 {
		return update.Disable().SetHelpText("Bu seçim için seçenek bulunamadı.")
	}

	return update.Enable()
})
```

## Kontrol Listesi

1. Field response içinde bağımlı alanda `depends_on` var mı?
2. Parent alan değişince network'te `resolve-dependencies` isteği gidiyor mu?
3. Response `fields.<dependent_key>.options` dolu mu?
4. UI'da select seçenekleri güncelleniyor mu?
5. Boş durumda `SetValue(nil)` ile eski seçim temizleniyor mu?

## Beklenen Field JSON Örneği

```json
{
  "key": "variant_option_id",
  "type": "select",
  "view": "select-field-form",
  "depends_on": ["product_id"],
  "props": {}
}
```

## Sık Hata ve Çözüm

### 1) Select options hiç gelmiyor

- Belirti: `props.options` boş kalıyor, callback hiç çalışmıyor.
- Kontrol: initial field response içinde `depends_on` yoksa frontend tetikleme yapmaz.
- Çözüm: `DependsOn(...)` tanımlı olduğundan ve field serializer'ın `depends_on` döndürdüğünden emin ol.

### 2) Callback çalışıyor ama seçim eski değerde kalıyor

- Çözüm: Callback içinde `SetValue(nil)` kullan.

### 3) Her şey doğru ama endpoint'e istek gitmiyor

- Kontrol: form modu `create`/`edit` mi, field form context'te mi?
- Kontrol: parent alan gerçekten değişiyor mu (watch tetikleniyor mu)?

## Hızlı API Testi

`resolve-dependencies` endpoint'ini manuel test ederek callback'i doğrulayabilirsin:

```json
POST /api/resource/productvariants/fields/resolve-dependencies
{
  "formData": {
    "product_id": 1,
    "variant_option_id": null
  },
  "context": "create",
  "changedFields": ["product_id"],
  "resourceId": null
}
```

Beklenti:

- `fields.variant_option_id.options` dolu gelmeli
- Boş durumda `disabled: true` ve açıklayıcı `helpText` gelebilmeli

