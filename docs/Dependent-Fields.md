# Bağımlı Alanlar (Dependent Fields)

Bağımlı alanlar, bir form alanının değeri değiştiğinde diğer alanların dinamik olarak güncellenmesini sağlar. Bu özellik sayesinde formlarınızı daha akıllı ve kullanıcı dostu hale getirebilirsiniz.

## Temel Kullanım

Bir alanı başka bir alana bağımlı hale getirmek için `DependsOn()` ve `OnDependencyChange()` metodlarını kullanın:

```go
fields.Text("Vergi Numarası", "tax_number").
    OnForm().
    DependsOn("company_type").
    OnDependencyChange(func(field *fields.Schema, formData map[string]interface{}, ctx *fiber.Ctx) *fields.FieldUpdate {
        update := fields.NewFieldUpdate()

        companyType, _ := formData["company_type"].(string)

        if companyType == "individual" {
            update.Hide().MakeOptional()
        } else if companyType == "company" {
            update.Show().MakeRequired().
                SetHelpText("Lütfen şirket vergi numaranızı girin")
        }

        return update
    })
```

## Kullanım Senaryoları

### 1. Görünürlük Kontrolü (Visibility Toggle)

Bir alanın değerine göre diğer alanları göster/gizle:

```go
fields.Select("Şirket Tipi", "company_type").
    OnForm().
    Options(map[string]string{
        "individual": "Şahıs",
        "company":    "Şirket",
    })

fields.Text("Vergi Numarası", "tax_number").
    OnForm().
    DependsOn("company_type").
    OnDependencyChange(func(field *fields.Schema, formData map[string]interface{}, ctx *fiber.Ctx) *fields.FieldUpdate {
        update := fields.NewFieldUpdate()

        if formData["company_type"] == "company" {
            update.Show().MakeRequired()
        } else {
            update.Hide().MakeOptional()
        }

        return update
    })
```

### 2. Dinamik Seçenekler (Dynamic Options)

Bir alanın değerine göre başka bir alanın seçeneklerini güncelle:

```go
fields.Select("Ülke", "country").
    OnForm().
    Options(map[string]string{
        "TR": "Türkiye",
        "US": "Amerika Birleşik Devletleri",
    })

fields.Select("Şehir", "city").
    OnForm().
    DependsOn("country").
    OnDependencyChange(func(field *fields.Schema, formData map[string]interface{}, ctx *fiber.Ctx) *fields.FieldUpdate {
        update := fields.NewFieldUpdate()

        country, _ := formData["country"].(string)

        if country == "TR" {
            update.SetOptions(map[string]interface{}{
                "34": "İstanbul",
                "06": "Ankara",
                "35": "İzmir",
            })
        } else if country == "US" {
            update.SetOptions(map[string]interface{}{
                "CA": "California",
                "NY": "New York",
                "TX": "Texas",
            })
        }

        return update
    })
```

### 3. Dinamik Doğrulama (Dynamic Validation)

Bir alanın değerine göre doğrulama kurallarını değiştir:

```go
fields.Select("Ödeme Yöntemi", "payment_method").
    OnForm().
    Options(map[string]string{
        "credit_card": "Kredi Kartı",
        "bank_transfer": "Banka Transferi",
    })

fields.Text("Kart Numarası", "card_number").
    OnForm().
    DependsOn("payment_method").
    OnDependencyChange(func(field *fields.Schema, formData map[string]interface{}, ctx *fiber.Ctx) *fields.FieldUpdate {
        update := fields.NewFieldUpdate()

        if formData["payment_method"] == "credit_card" {
            update.Show().MakeRequired().
                SetHelpText("16 haneli kart numaranızı girin")
        } else {
            update.Hide().MakeOptional()
        }

        return update
    })
```

### 4. Koşullu Salt Okunur (Conditional Readonly)

Bir alanın değerine göre diğer alanları salt okunur yap:

```go
fields.Switch("Otomatik Hesapla", "auto_calculate").
    OnForm()

fields.Number("Toplam", "total").
    OnForm().
    DependsOn("auto_calculate").
    OnDependencyChange(func(field *fields.Schema, formData map[string]interface{}, ctx *fiber.Ctx) *fields.FieldUpdate {
        update := fields.NewFieldUpdate()

        if formData["auto_calculate"] == true {
            update.MakeReadOnly().
                SetHelpText("Otomatik hesaplanıyor")
        } else {
            update.MakeEditable()
        }

        return update
    })
```

### 5. Değer Hesaplama (Value Calculation)

Bir alanın değerine göre diğer alanın değerini hesapla:

```go
fields.Number("Fiyat", "price").OnForm()
fields.Number("Miktar", "quantity").OnForm()

fields.Number("Toplam", "total").
    OnForm().
    DependsOn("price", "quantity").
    OnDependencyChange(func(field *fields.Schema, formData map[string]interface{}, ctx *fiber.Ctx) *fields.FieldUpdate {
        update := fields.NewFieldUpdate()

        price, _ := formData["price"].(float64)
        quantity, _ := formData["quantity"].(float64)

        total := price * quantity
        update.SetValue(total).MakeReadOnly()

        return update
    })
```

## Bağlam Bazlı Callback'ler (Context-Aware Callbacks)

Oluşturma ve güncelleme formları için farklı davranışlar tanımlayabilirsiniz:

```go
fields.Text("Kullanıcı Adı", "username").
    OnForm().
    DependsOn("email").
    OnDependencyChangeCreating(func(field *fields.Schema, formData map[string]interface{}, ctx *fiber.Ctx) *fields.FieldUpdate {
        // Sadece oluşturma formunda çalışır
        update := fields.NewFieldUpdate()

        email, _ := formData["email"].(string)
        if email != "" {
            // Email'den kullanıcı adı öner
            username := strings.Split(email, "@")[0]
            update.SetValue(username)
        }

        return update
    }).
    OnDependencyChangeUpdating(func(field *fields.Schema, formData map[string]interface{}, ctx *fiber.Ctx) *fields.FieldUpdate {
        // Sadece güncelleme formunda çalışır
        update := fields.NewFieldUpdate()
        update.MakeReadOnly().
            SetHelpText("Kullanıcı adı değiştirilemez")

        return update
    })
```

## FieldUpdate API Referansı

`FieldUpdate` nesnesi, alanın özelliklerini güncellemek için kullanılır:

### Görünürlük Metodları

- `Show()` - Alanı görünür yap
- `Hide()` - Alanı gizle

### Doğrulama Metodları

- `MakeRequired()` - Alanı zorunlu yap
- `MakeOptional()` - Alanı isteğe bağlı yap

### Düzenleme Metodları

- `MakeReadOnly()` - Alanı salt okunur yap
- `MakeEditable()` - Alanı düzenlenebilir yap
- `Enable()` - Alanı etkinleştir
- `Disable()` - Alanı devre dışı bırak

### İçerik Metodları

- `SetHelpText(text string)` - Yardım metnini ayarla
- `SetPlaceholder(text string)` - Yer tutucu metnini ayarla
- `SetOptions(options map[string]interface{})` - Seçenekleri ayarla (Select alanları için)
- `SetValue(value interface{})` - Alanın değerini ayarla

### Zincirleme Kullanım

Tüm metodlar zincirleme kullanım için `*FieldUpdate` döndürür:

```go
update.Show().MakeRequired().SetHelpText("Bu alan zorunludur")
```

## Çoklu Bağımlılıklar

Bir alan birden fazla alana bağımlı olabilir:

```go
fields.Text("İndirimli Fiyat", "discounted_price").
    OnForm().
    DependsOn("price", "discount_rate", "tax_rate").
    OnDependencyChange(func(field *fields.Schema, formData map[string]interface{}, ctx *fiber.Ctx) *fields.FieldUpdate {
        update := fields.NewFieldUpdate()

        price, _ := formData["price"].(float64)
        discountRate, _ := formData["discount_rate"].(float64)
        taxRate, _ := formData["tax_rate"].(float64)

        discountedPrice := price * (1 - discountRate/100) * (1 + taxRate/100)
        update.SetValue(discountedPrice).MakeReadOnly()

        return update
    })
```

## Performans ve Best Practices

### 1. Debouncing

Frontend otomatik olarak 300ms debouncing uygular. Kullanıcı yazmayı bıraktıktan 300ms sonra API çağrısı yapılır.

### 2. Dairesel Bağımlılıkları Önleme

Sistem otomatik olarak dairesel bağımlılıkları tespit eder ve hata döndürür:

```go
// ❌ YANLIŞ - Dairesel bağımlılık
fields.Text("A", "field_a").DependsOn("field_b")
fields.Text("B", "field_b").DependsOn("field_a")
```

### 3. Hafif Callback'ler

Callback fonksiyonlarınızı mümkün olduğunca hafif tutun. Ağır işlemler için asenkron işleme kullanın:

```go
// ✅ DOĞRU - Hafif callback
OnDependencyChange(func(field *fields.Schema, formData map[string]interface{}, ctx *fiber.Ctx) *fields.FieldUpdate {
    update := fields.NewFieldUpdate()

    if formData["type"] == "premium" {
        update.Show()
    } else {
        update.Hide()
    }

    return update
})

// ❌ YANLIŞ - Ağır işlem
OnDependencyChange(func(field *fields.Schema, formData map[string]interface{}, ctx *fiber.Ctx) *fields.FieldUpdate {
    // Veritabanı sorgusu veya API çağrısı yapma
    // Bu işlemler callback'i yavaşlatır
    return update
})
```

### 4. Null Kontrolü

Form verilerini kullanırken her zaman null kontrolü yapın:

```go
OnDependencyChange(func(field *fields.Schema, formData map[string]interface{}, ctx *fiber.Ctx) *fields.FieldUpdate {
    update := fields.NewFieldUpdate()

    // ✅ DOĞRU - Null kontrolü
    value, ok := formData["field_name"].(string)
    if !ok || value == "" {
        return update
    }

    // Değeri kullan
    update.SetHelpText(value)

    return update
})
```

## Frontend Entegrasyonu

Frontend tarafında ResourceForm bileşeni otomatik olarak bağımlı alanları yönetir. Ek bir yapılandırma gerekmez:

```tsx
<ResourceForm
    fields={fields}
    initialData={initialData}
    onSubmit={handleSubmit}
    resource="users"
    context="create"
/>
```

### Gerekli Props

- `resource` - Kaynak adı (örn: "users", "posts")
- `context` - Form bağlamı ("create" veya "update")
- `resourceId` - Güncelleme formları için kaynak ID'si (isteğe bağlı)

## Mimari

### Backend

1. **Dependency Graph** - Bağımlılık grafiği BFS algoritması ile oluşturulur
2. **Affected Fields** - Değişen alandan etkilenen tüm alanlar bulunur
3. **Callback Execution** - Her etkilenen alan için callback çalıştırılır
4. **Response** - Güncellenmiş alan özellikleri döndürülür

### Frontend

1. **Field Change** - Kullanıcı bir alanı değiştirir
2. **Debounce** - 300ms beklenir
3. **API Call** - `/resource/:resource/fields/resolve-dependencies` endpoint'ine POST isteği
4. **Update UI** - Dönen güncellemeler UI'a uygulanır

## API Endpoint

```
POST /resource/:resource/fields/resolve-dependencies
```

**Request Body:**
```json
{
  "formData": {
    "company_type": "company",
    "name": "Acme Inc"
  },
  "context": "create",
  "changedFields": ["company_type"],
  "resourceId": null
}
```

**Response:**
```json
{
  "fields": {
    "tax_number": {
      "visible": true,
      "required": true,
      "helpText": "Lütfen şirket vergi numaranızı girin"
    }
  }
}
```

## Hata Yönetimi

Sistem otomatik olarak şu hataları yönetir:

1. **Dairesel Bağımlılık** - Tespit edilir ve hata döndürülür
2. **Geçersiz Bağlam** - "create" veya "update" dışında bir değer hata döndürür
3. **Callback Hataları** - Callback'te oluşan hatalar loglanır ve boş güncelleme döndürülür

## Örnekler

### Tam Örnek: E-Ticaret Formu

```go
func (r *ProductResource) Fields() []core.Element {
    return []core.Element{
        fields.Select("Ürün Tipi", "product_type").
            OnForm().
            Options(map[string]string{
                "physical": "Fiziksel Ürün",
                "digital":  "Dijital Ürün",
            }),

        fields.Number("Ağırlık (kg)", "weight").
            OnForm().
            DependsOn("product_type").
            OnDependencyChange(func(field *fields.Schema, formData map[string]interface{}, ctx *fiber.Ctx) *fields.FieldUpdate {
                update := fields.NewFieldUpdate()

                if formData["product_type"] == "physical" {
                    update.Show().MakeRequired().
                        SetHelpText("Kargo hesaplaması için gerekli")
                } else {
                    update.Hide().MakeOptional()
                }

                return update
            }),

        fields.Text("İndirme Linki", "download_url").
            OnForm().
            DependsOn("product_type").
            OnDependencyChange(func(field *fields.Schema, formData map[string]interface{}, ctx *fiber.Ctx) *fields.FieldUpdate {
                update := fields.NewFieldUpdate()

                if formData["product_type"] == "digital" {
                    update.Show().MakeRequired().
                        SetHelpText("Ürün indirme linki")
                } else {
                    update.Hide().MakeOptional()
                }

                return update
            }),

        fields.Number("Fiyat", "price").OnForm().Required(),

        fields.Number("KDV Oranı (%)", "tax_rate").
            OnForm().
            Default(18),

        fields.Number("KDV Dahil Fiyat", "price_with_tax").
            OnForm().
            DependsOn("price", "tax_rate").
            OnDependencyChange(func(field *fields.Schema, formData map[string]interface{}, ctx *fiber.Ctx) *fields.FieldUpdate {
                update := fields.NewFieldUpdate()

                price, _ := formData["price"].(float64)
                taxRate, _ := formData["tax_rate"].(float64)

                priceWithTax := price * (1 + taxRate/100)
                update.SetValue(priceWithTax).MakeReadOnly().
                    SetHelpText(fmt.Sprintf("%.2f TL", priceWithTax))

                return update
            }),
    }
}
```

## İlgili Dokümantasyon

- [Fields](Fields.md) - Alan türleri ve temel kullanım
- [Resources](Resources.md) - Kaynak tanımlama
- [API Reference](API-Reference.md) - API endpoint'leri
