package fields

import (
	"github.com/gofiber/fiber/v2"
)

// DependencyCallbackFunc, form alanları arasındaki bağımlılık ilişkilerini yönetmek için kullanılan
// callback fonksiyon tipidir.
//
// # Genel Bakış
//
// Bu fonksiyon tipi, bir alanın değeri değiştiğinde diğer alanların davranışını dinamik olarak
// değiştirmek için kullanılır. Örneğin, bir dropdown'da seçilen değere göre başka bir alanı
// göstermek/gizlemek, zorunlu hale getirmek veya seçeneklerini değiştirmek gibi işlemler yapılabilir.
//
// # Kullanım Senaryoları
//
// - **Koşullu Alan Görünürlüğü**: Bir alanın değerine göre diğer alanları göster/gizle
// - **Dinamik Validasyon**: Bir alanın değerine göre diğer alanların validasyon kurallarını değiştir
// - **Cascade Seçenekler**: Bir dropdown'ın seçimine göre diğer dropdown'ların seçeneklerini filtrele
// - **Koşullu Zorunluluk**: Belirli koşullarda alanları zorunlu veya opsiyonel yap
// - **Dinamik Yardım Metinleri**: Bağlama göre yardım metinlerini güncelle
//
// # Parametreler
//
// - `field`: Bağımlılık kuralının uygulanacağı hedef alan şeması
// - `formData`: Form verilerini içeren map (tüm alanların güncel değerleri)
// - `ctx`: Fiber context nesnesi (HTTP request/response bilgilerine erişim için)
//
// # Dönüş Değeri
//
// `*FieldUpdate`: Alana uygulanacak güncellemeleri içeren nesne. nil döndürülürse hiçbir güncelleme yapılmaz.
//
// # Kullanım Örneği
//
// ```go
// // Ülke seçimine göre şehir alanını göster/gizle
// countryField.OnChange(func(field *Schema, formData map[string]interface{}, ctx *fiber.Ctx) *FieldUpdate {
//     country := formData["country"]
//     if country == "TR" {
//         return NewFieldUpdate().Show().MakeRequired()
//     }
//     return NewFieldUpdate().Hide().MakeOptional()
// })
//
// // Ürün tipine göre fiyat alanının validasyonunu değiştir
// productTypeField.OnChange(func(field *Schema, formData map[string]interface{}, ctx *fiber.Ctx) *FieldUpdate {
//     productType := formData["product_type"]
//     if productType == "premium" {
//         return NewFieldUpdate().
//             SetHelpText("Premium ürünler için minimum 1000 TL").
//             AddRule(ValidationRule{Type: "min", Value: 1000})
//     }
//     return NewFieldUpdate().SetHelpText("Standart fiyatlandırma")
// })
// ```
//
// # Önemli Notlar
//
// - Callback fonksiyonu her form değişikliğinde çağrılabilir, bu nedenle performans açısından hafif olmalıdır
// - Sonsuz döngülerden kaçınmak için dikkatli olunmalıdır (A alanı B'yi tetikler, B alanı A'yı tetikler)
// - formData map'i güvenli bir şekilde kontrol edilmeli, nil veya eksik değerler için kontrol yapılmalıdır
// - Context üzerinden veritabanı sorguları yapılabilir ancak performans etkileri göz önünde bulundurulmalıdır
//
// # Avantajlar
//
// - Dinamik ve esnek form davranışı sağlar
// - Kullanıcı deneyimini iyileştirir (gereksiz alanları gizler)
// - Karmaşık iş kurallarını kolayca uygular
// - Sunucu tarafında kontrol sağlar (güvenlik)
//
// # Dezavantajlar
//
// - Yanlış kullanımda performans sorunlarına yol açabilir
// - Karmaşık bağımlılık zincirleri debug edilmesi zor olabilir
// - Sonsuz döngü riski vardır
type DependencyCallbackFunc func(
	field *Schema,
	formData map[string]interface{},
	ctx *fiber.Ctx,
) *FieldUpdate

// FieldUpdate, bağımlılık değişikliklerine göre bir alana uygulanacak güncellemeleri temsil eder.
//
// # Genel Bakış
//
// Bu struct, form alanlarının dinamik davranışını kontrol etmek için kullanılır. Bir alanın
// görünürlüğünden validasyon kurallarına kadar tüm özelliklerini değiştirmeye olanak tanır.
// Pointer kullanımı sayesinde sadece değiştirilmek istenen özellikler güncellenir, diğerleri
// mevcut hallerini korur.
//
// # Struct Alanları
//
// - `Visible`: Alanın görünür olup olmadığını kontrol eder (nil = değişiklik yok)
// - `ReadOnly`: Alanın salt okunur olup olmadığını kontrol eder (nil = değişiklik yok)
// - `Required`: Alanın zorunlu olup olmadığını kontrol eder (nil = değişiklik yok)
// - `Disabled`: Alanın devre dışı olup olmadığını kontrol eder (nil = değişiklik yok)
// - `HelpText`: Alanın yardım metnini günceller (nil = değişiklik yok)
// - `Placeholder`: Alanın placeholder metnini günceller (nil = değişiklik yok)
// - `Options`: Alanın seçeneklerini günceller (select, radio vb. için)
// - `Value`: Alanın değerini programatik olarak ayarlar
// - `Rules`: Alanın validasyon kurallarını günceller
//
// # Kullanım Senaryoları
//
// - **Koşullu Görünürlük**: Bir alanı belirli koşullarda göster/gizle
// - **Dinamik Validasyon**: Koşullara göre validasyon kurallarını değiştir
// - **Cascade Seçenekler**: Üst seçime göre alt seçenekleri filtrele
// - **Koşullu Zorunluluk**: Belirli durumlarda alanı zorunlu yap
// - **Dinamik Değer Atama**: Hesaplanan değerleri alanlara otomatik ata
// - **Kullanıcı Rehberliği**: Bağlama göre yardım metinlerini güncelle
//
// # Kullanım Örneği
//
// ```go
// // Basit görünürlük kontrolü
// update := NewFieldUpdate().Show().MakeRequired()
//
// // Çoklu özellik güncelleme (method chaining)
// update := NewFieldUpdate().
//     Show().
//     MakeRequired().
//     SetHelpText("Bu alan zorunludur").
//     SetPlaceholder("Lütfen bir değer girin")
//
// // Dinamik seçenekler güncelleme
// cities := map[string]interface{}{
//     "istanbul": "İstanbul",
//     "ankara": "Ankara",
//     "izmir": "İzmir",
// }
// update := NewFieldUpdate().SetOptions(cities)
//
// // Validasyon kuralları ekleme
// update := NewFieldUpdate().
//     MakeRequired().
//     AddRule(ValidationRule{Type: "min", Value: 18}).
//     AddRule(ValidationRule{Type: "max", Value: 65})
//
// // Değer atama ve salt okunur yapma
// update := NewFieldUpdate().
//     SetValue("Otomatik hesaplanan değer").
//     MakeReadOnly()
// ```
//
// # Önemli Notlar
//
// - **Pointer Kullanımı**: Tüm boolean ve string alanlar pointer'dır. Bu sayede nil değer
//   "değişiklik yok" anlamına gelir, false/true veya boş string ile karıştırılmaz
// - **Method Chaining**: Tüm setter metodlar *FieldUpdate döndürür, bu sayede zincirleme
//   çağrılar yapılabilir
// - **Seçici Güncelleme**: Sadece set edilen alanlar güncellenir, diğerleri korunur
// - **JSON Serileştirme**: omitempty tag'leri sayesinde nil alanlar JSON'a dahil edilmez
// - **Thread Safety**: Bu struct thread-safe değildir, concurrent kullanımda dikkatli olunmalıdır
//
// # Avantajlar
//
// - Esnek ve güçlü alan kontrolü sağlar
// - Method chaining ile okunabilir kod yazılmasını sağlar
// - Seçici güncelleme ile performans optimizasyonu sağlar
// - Type-safe alan güncellemeleri yapar
// - JSON serileştirme desteği ile API entegrasyonu kolaydır
//
// # Dezavantajlar
//
// - Pointer kullanımı nil kontrollerini gerektirir
// - Karmaşık güncellemelerde kod okunabilirliği azalabilir
// - Thread-safe değildir
//
// # Best Practices
//
// - Her zaman NewFieldUpdate() ile yeni instance oluşturun
// - Method chaining kullanarak okunabilir kod yazın
// - Nil kontrollerini ihmal etmeyin
// - Gereksiz güncellemelerden kaçının (performans için)
type FieldUpdate struct {
	Visible     *bool                  `json:"visible,omitempty"`
	ReadOnly    *bool                  `json:"readonly,omitempty"`
	Required    *bool                  `json:"required,omitempty"`
	Disabled    *bool                  `json:"disabled,omitempty"`
	HelpText    *string                `json:"helpText,omitempty"`
	Placeholder *string                `json:"placeholder,omitempty"`
	Options     map[string]interface{} `json:"options,omitempty"`
	Value       interface{}            `json:"value,omitempty"`
	Rules       []ValidationRule       `json:"rules,omitempty"`
}

// NewFieldUpdate, yeni bir FieldUpdate instance'ı oluşturur.
//
// # Genel Bakış
//
// Bu fonksiyon, FieldUpdate struct'ının yeni bir instance'ını oluşturur ve döndürür.
// Tüm alanlar nil/boş değerlerle başlatılır, bu sayede sadece set edilen özellikler
// güncelleme sırasında uygulanır.
//
// # Kullanım Senaryoları
//
// - Bağımlılık callback fonksiyonlarında yeni güncelleme nesnesi oluşturmak
// - Method chaining ile zincirleme güncelleme işlemleri yapmak
// - Temiz ve okunabilir kod yazmak için başlangıç noktası
//
// # Dönüş Değeri
//
// `*FieldUpdate`: Yeni oluşturulmuş, boş bir FieldUpdate pointer'ı
//
// # Kullanım Örneği
//
// ```go
// // Basit kullanım
// update := NewFieldUpdate()
// update.Show()
// update.MakeRequired()
//
// // Method chaining ile kullanım (önerilen)
// update := NewFieldUpdate().
//     Show().
//     MakeRequired().
//     SetHelpText("Bu alan zorunludur")
//
// // Bağımlılık callback'inde kullanım
// field.OnChange(func(field *Schema, formData map[string]interface{}, ctx *fiber.Ctx) *FieldUpdate {
//     if formData["country"] == "TR" {
//         return NewFieldUpdate().Show().MakeRequired()
//     }
//     return NewFieldUpdate().Hide()
// })
// ```
//
// # Önemli Notlar
//
// - Her güncelleme işlemi için yeni bir instance oluşturun, mevcut instance'ları yeniden kullanmayın
// - Method chaining kullanarak daha okunabilir kod yazabilirsiniz
// - Nil döndürmek yerine boş FieldUpdate döndürmek güvenlidir (hiçbir güncelleme yapılmaz)
//
// # Best Practices
//
// - Her zaman bu fonksiyonu kullanarak yeni instance oluşturun
// - Struct literal kullanmak yerine bu factory fonksiyonunu tercih edin
// - Method chaining ile zincirleme çağrılar yapın
func NewFieldUpdate() *FieldUpdate {
	return &FieldUpdate{}
}

// Show, alanı görünür hale getirir.
//
// # Genel Bakış
//
// Bu metod, bir form alanının kullanıcı arayüzünde görünür olmasını sağlar. Visible özelliğini
// true olarak ayarlar ve method chaining desteği için *FieldUpdate döndürür.
//
// # Kullanım Senaryoları
//
// - Koşullu alan görünürlüğü: Belirli bir değer seçildiğinde ilgili alanları göster
// - Dinamik form yapısı: Kullanıcı seçimlerine göre form alanlarını göster/gizle
// - Wizard formlar: Adım adım form akışında ilgili alanları göster
// - Koşullu validasyon: Görünür olan alanlar için validasyon uygula
//
// # Dönüş Değeri
//
// `*FieldUpdate`: Method chaining için güncellenmiş FieldUpdate pointer'ı
//
// # Kullanım Örneği
//
// ```go
// // Basit kullanım
// update := NewFieldUpdate().Show()
//
// // Koşullu görünürlük
// if formData["has_company"] == true {
//     return NewFieldUpdate().Show().MakeRequired()
// }
//
// // Çoklu alan kontrolü
// countryField.OnChange(func(field *Schema, formData map[string]interface{}, ctx *fiber.Ctx) *FieldUpdate {
//     country := formData["country"]
//     if country == "TR" {
//         // Türkiye seçildiğinde TC Kimlik No alanını göster
//         return NewFieldUpdate().Show().MakeRequired().SetHelpText("11 haneli TC Kimlik No")
//     }
//     return NewFieldUpdate().Hide()
// })
//
// // Method chaining ile birlikte kullanım
// update := NewFieldUpdate().
//     Show().
//     MakeRequired().
//     Enable().
//     SetHelpText("Bu alan artık görünür ve zorunlu")
// ```
//
// # Önemli Notlar
//
// - Show() çağrıldığında alan görünür olur ancak disabled durumu değişmez
// - Görünür olan alanlar için validasyon kuralları uygulanır
// - Hide() ile gizlenen alanlar için validasyon atlanır
// - Method chaining desteklenir, zincirleme çağrılar yapılabilir
//
// # Best Practices
//
// - Show() ile birlikte MakeRequired() veya MakeOptional() kullanarak validasyon durumunu belirtin
// - Görünür hale getirilen alanlara uygun yardım metni ekleyin
// - Koşullu görünürlükte her iki durumu da (Show/Hide) ele alın
func (u *FieldUpdate) Show() *FieldUpdate {
	visible := true
	u.Visible = &visible
	return u
}

// Hide, alanı gizler.
//
// # Genel Bakış
//
// Bu metod, bir form alanının kullanıcı arayüzünde gizlenmesini sağlar. Visible özelliğini
// false olarak ayarlar ve method chaining desteği için *FieldUpdate döndürür.
//
// # Kullanım Senaryoları
//
// - Koşullu alan gizleme: Belirli bir değer seçildiğinde ilgisiz alanları gizle
// - Dinamik form yapısı: Kullanıcı seçimlerine göre gereksiz alanları gizle
// - Basitleştirilmiş form: Kullanıcıya sadece ilgili alanları göster
// - Validasyon atlama: Gizli alanlar için validasyon uygulanmaz
//
// # Dönüş Değeri
//
// `*FieldUpdate`: Method chaining için güncellenmiş FieldUpdate pointer'ı
//
// # Kullanım Örneği
//
// ```go
// // Basit kullanım
// update := NewFieldUpdate().Hide()
//
// // Koşullu gizleme
// if formData["payment_method"] == "cash" {
//     // Nakit ödeme seçildiğinde kredi kartı alanlarını gizle
//     return NewFieldUpdate().Hide().MakeOptional()
// }
//
// // Karmaşık koşullu mantık
// productTypeField.OnChange(func(field *Schema, formData map[string]interface{}, ctx *fiber.Ctx) *FieldUpdate {
//     productType := formData["product_type"]
//     if productType == "digital" {
//         // Dijital ürünler için fiziksel ürün alanlarını gizle
//         return NewFieldUpdate().Hide().MakeOptional().SetValue(nil)
//     }
//     return NewFieldUpdate().Show().MakeRequired()
// })
//
// // Method chaining ile birlikte kullanım
// update := NewFieldUpdate().
//     Hide().
//     MakeOptional().
//     SetValue(nil) // Gizli alan için değeri temizle
// ```
//
// # Önemli Notlar
//
// - Hide() çağrıldığında alan gizlenir ve validasyon kuralları uygulanmaz
// - Gizli alanların değerleri form submit edildiğinde gönderilmeyebilir
// - Güvenlik açısından hassas alanlar için Hide() yeterli değildir, sunucu tarafı kontrolü gereklidir
// - Method chaining desteklenir, zincirleme çağrılar yapılabilir
// - Gizlenen alanları MakeOptional() ile birlikte kullanmak önerilir
//
// # Best Practices
//
// - Hide() ile birlikte MakeOptional() kullanarak validasyon hatalarını önleyin
// - Gizlenen alanların değerlerini SetValue(nil) ile temizlemeyi düşünün
// - Koşullu gizlemede her iki durumu da (Show/Hide) ele alın
// - Güvenlik kritik alanlar için sunucu tarafı kontrolü yapın
func (u *FieldUpdate) Hide() *FieldUpdate {
	visible := false
	u.Visible = &visible
	return u
}

// MakeReadOnly, alanı salt okunur (read-only) hale getirir.
//
// # Genel Bakış
//
// Bu metod, bir form alanının kullanıcı tarafından düzenlenemez ancak görüntülenebilir olmasını
// sağlar. ReadOnly özelliğini true olarak ayarlar ve method chaining desteği için *FieldUpdate döndürür.
//
// # Kullanım Senaryoları
//
// - **Hesaplanan Değerler**: Otomatik hesaplanan alanları göster ancak düzenlemeye izin verme
// - **Sistem Alanları**: Sistem tarafından yönetilen alanları koruma altına al
// - **Onay Sonrası Kilitleme**: Onaylanmış kayıtların belirli alanlarını kilitle
// - **Referans Bilgileri**: Referans amaçlı gösterilen ancak değiştirilmemesi gereken bilgiler
// - **Audit Trail**: Değişiklik geçmişi için orijinal değerleri koruma
//
// # Dönüş Değeri
//
// `*FieldUpdate`: Method chaining için güncellenmiş FieldUpdate pointer'ı
//
// # Kullanım Örneği
//
// ```go
// // Basit kullanım
// update := NewFieldUpdate().MakeReadOnly()
//
// // Hesaplanan toplam alanı
// quantityField.OnChange(func(field *Schema, formData map[string]interface{}, ctx *fiber.Ctx) *FieldUpdate {
//     quantity := formData["quantity"].(float64)
//     price := formData["price"].(float64)
//     total := quantity * price
//
//     // Toplam alanını hesapla ve salt okunur yap
//     return NewFieldUpdate().
//         SetValue(total).
//         MakeReadOnly().
//         SetHelpText(fmt.Sprintf("Otomatik hesaplanan: %.2f TL", total))
// })
//
// // Onaylanmış kayıtlar için
// statusField.OnChange(func(field *Schema, formData map[string]interface{}, ctx *fiber.Ctx) *FieldUpdate {
//     status := formData["status"]
//     if status == "approved" {
//         // Onaylandıktan sonra fiyat alanını kilitle
//         return NewFieldUpdate().
//             MakeReadOnly().
//             SetHelpText("Onaylanmış kayıtlarda fiyat değiştirilemez")
//     }
//     return NewFieldUpdate().MakeEditable()
// })
//
// // Sistem alanları
// createdAtField := NewFieldUpdate().
//     MakeReadOnly().
//     SetHelpText("Sistem tarafından otomatik oluşturuldu")
// ```
//
// # Önemli Notlar
//
// - **Görünürlük**: ReadOnly alanlar görünür kalır ancak düzenlenemez
// - **Validasyon**: ReadOnly alanlar için validasyon kuralları uygulanmaz (değer değişmediği için)
// - **Form Submit**: ReadOnly alanların değerleri form submit edildiğinde gönderilir
// - **Güvenlik**: Sunucu tarafında da kontrol yapılmalıdır, sadece client-side koruma yeterli değildir
// - **Disabled vs ReadOnly**: Disabled alanlar form submit'e dahil edilmez, ReadOnly alanlar dahil edilir
//
// # Best Practices
//
// - Hesaplanan değerler için SetValue() ile birlikte kullanın
// - Kullanıcıya neden düzenleyemediğini açıklayan yardım metni ekleyin
// - Sunucu tarafında da aynı kontrolü uygulayın (güvenlik)
// - Geçici kilitleme durumları için koşullu kullanın
func (u *FieldUpdate) MakeReadOnly() *FieldUpdate {
	readOnly := true
	u.ReadOnly = &readOnly
	return u
}

// MakeEditable, alanı düzenlenebilir hale getirir.
//
// # Genel Bakış
//
// Bu metod, daha önce salt okunur yapılmış bir alanı tekrar düzenlenebilir hale getirir.
// ReadOnly özelliğini false olarak ayarlar ve method chaining desteği için *FieldUpdate döndürür.
//
// # Kullanım Senaryoları
//
// - **Koşullu Düzenleme**: Belirli koşullarda alanları düzenlenebilir yap
// - **Rol Bazlı Erişim**: Yetki seviyesine göre alanları aç/kapat
// - **Durum Bazlı Kontrol**: Kayıt durumuna göre düzenleme izni ver
// - **Dinamik Form**: Kullanıcı seçimlerine göre alanları aktif et
//
// # Dönüş Değeri
//
// `*FieldUpdate`: Method chaining için güncellenmiş FieldUpdate pointer'ı
//
// # Kullanım Örneği
//
// ```go
// // Basit kullanım
// update := NewFieldUpdate().MakeEditable()
//
// // Durum bazlı düzenleme
// statusField.OnChange(func(field *Schema, formData map[string]interface{}, ctx *fiber.Ctx) *FieldUpdate {
//     status := formData["status"]
//     if status == "draft" {
//         // Taslak durumunda tüm alanlar düzenlenebilir
//         return NewFieldUpdate().
//             MakeEditable().
//             SetHelpText("Taslak durumunda düzenleyebilirsiniz")
//     }
//     return NewFieldUpdate().
//         MakeReadOnly().
//         SetHelpText("Sadece taslak durumunda düzenlenebilir")
// })
//
// // Rol bazlı erişim
// if ctx.Locals("user_role") == "admin" {
//     return NewFieldUpdate().
//         MakeEditable().
//         SetHelpText("Admin olarak düzenleme yetkisine sahipsiniz")
// }
//
// // Koşullu düzenleme
// editModeField.OnChange(func(field *Schema, formData map[string]interface{}, ctx *fiber.Ctx) *FieldUpdate {
//     editMode := formData["edit_mode"].(bool)
//     if editMode {
//         return NewFieldUpdate().MakeEditable().Enable()
//     }
//     return NewFieldUpdate().MakeReadOnly().Disable()
// })
// ```
//
// # Önemli Notlar
//
// - MakeEditable() çağrıldığında alan düzenlenebilir olur ancak disabled durumu değişmez
// - Düzenlenebilir alanlar için validasyon kuralları uygulanır
// - Sunucu tarafında da aynı yetki kontrolü yapılmalıdır
// - Method chaining desteklenir
//
// # Best Practices
//
// - Düzenleme izni verirken uygun yardım metni ekleyin
// - Sunucu tarafında yetki kontrolü yapın (güvenlik)
// - Enable() ile birlikte kullanarak tam erişim sağlayın
// - Koşullu düzenlemede her iki durumu da (Editable/ReadOnly) ele alın
func (u *FieldUpdate) MakeEditable() *FieldUpdate {
	readOnly := false
	u.ReadOnly = &readOnly
	return u
}

// MakeRequired, alanı zorunlu hale getirir.
//
// # Genel Bakış
//
// Bu metod, bir form alanının doldurulmasını zorunlu kılar. Required özelliğini true olarak
// ayarlar ve method chaining desteği için *FieldUpdate döndürür. Zorunlu alanlar için
// validasyon kuralları otomatik olarak uygulanır.
//
// # Kullanım Senaryoları
//
// - **Koşullu Zorunluluk**: Belirli bir seçim yapıldığında ilgili alanları zorunlu yap
// - **Dinamik Validasyon**: Kullanıcı akışına göre zorunlu alanları değiştir
// - **İş Kuralları**: İş mantığına göre alanları zorunlu hale getir
// - **Adım Bazlı Formlar**: Her adımda farklı alanları zorunlu yap
// - **Rol Bazlı Zorunluluk**: Kullanıcı rolüne göre zorunlu alanları belirle
//
// # Dönüş Değeri
//
// `*FieldUpdate`: Method chaining için güncellenmiş FieldUpdate pointer'ı
//
// # Kullanım Örneği
//
// ```go
// // Basit kullanım
// update := NewFieldUpdate().MakeRequired()
//
// // Koşullu zorunluluk
// shippingTypeField.OnChange(func(field *Schema, formData map[string]interface{}, ctx *fiber.Ctx) *FieldUpdate {
//     shippingType := formData["shipping_type"]
//     if shippingType == "express" {
//         // Ekspres kargo için telefon numarası zorunlu
//         return NewFieldUpdate().
//             Show().
//             MakeRequired().
//             SetHelpText("Ekspres kargo için telefon numarası zorunludur")
//     }
//     return NewFieldUpdate().MakeOptional()
// })
//
// // İş kuralı bazlı zorunluluk
// companyField.OnChange(func(field *Schema, formData map[string]interface{}, ctx *fiber.Ctx) *FieldUpdate {
//     hasCompany := formData["has_company"].(bool)
//     if hasCompany {
//         // Şirket varsa vergi numarası zorunlu
//         return NewFieldUpdate().
//             Show().
//             MakeRequired().
//             SetHelpText("Kurumsal müşteriler için vergi numarası zorunludur").
//             SetPlaceholder("10 haneli vergi numarası")
//     }
//     return NewFieldUpdate().Hide().MakeOptional()
// })
//
// // Çoklu alan kontrolü
// paymentMethodField.OnChange(func(field *Schema, formData map[string]interface{}, ctx *fiber.Ctx) *FieldUpdate {
//     paymentMethod := formData["payment_method"]
//     if paymentMethod == "credit_card" {
//         // Kredi kartı seçildiğinde kart bilgileri zorunlu
//         return NewFieldUpdate().
//             Show().
//             MakeRequired().
//             AddRule(ValidationRule{Type: "credit_card"})
//     }
//     return NewFieldUpdate().Hide().MakeOptional()
// })
//
// // Rol bazlı zorunluluk
// userRole := ctx.Locals("user_role").(string)
// if userRole == "premium" {
//     return NewFieldUpdate().
//         MakeRequired().
//         SetHelpText("Premium üyeler için bu alan zorunludur")
// }
// ```
//
// # Önemli Notlar
//
// - **Validasyon**: Zorunlu alanlar boş bırakıldığında validasyon hatası verir
// - **Görünürlük**: Gizli alanlar zorunlu yapılsa bile validasyon uygulanmaz
// - **Kullanıcı Deneyimi**: Zorunlu alanlar için açıklayıcı yardım metni ekleyin
// - **Form Submit**: Zorunlu alanlar doldurulmadan form submit edilemez
// - **API Validasyonu**: Sunucu tarafında da aynı zorunluluk kontrolü yapılmalıdır
//
// # Best Practices
//
// - Zorunlu alanlar için SetHelpText() ile açıklama ekleyin
// - Show() ile birlikte kullanarak görünürlüğü sağlayın
// - Kullanıcıya neden zorunlu olduğunu açıklayın
// - Sunucu tarafında da validasyon yapın (güvenlik)
// - Gereksiz zorunluluktan kaçının (kullanıcı deneyimi için)
//
// # Avantajlar
//
// - Veri bütünlüğünü sağlar
// - İş kurallarını zorlar
// - Kullanıcı hatalarını önler
// - Dinamik form davranışı sağlar
//
// # Dezavantajlar
//
// - Aşırı kullanımda kullanıcı deneyimini olumsuz etkiler
// - Karmaşık koşullu mantıkta hata yapma riski artar
func (u *FieldUpdate) MakeRequired() *FieldUpdate {
	required := true
	u.Required = &required
	return u
}

// MakeOptional, alanı opsiyonel (zorunlu olmayan) hale getirir.
//
// # Genel Bakış
//
// Bu metod, daha önce zorunlu yapılmış bir alanı opsiyonel hale getirir. Required özelliğini
// false olarak ayarlar ve method chaining desteği için *FieldUpdate döndürür.
//
// # Kullanım Senaryoları
//
// - **Koşullu Opsiyonellik**: Belirli durumlarda alanları opsiyonel yap
// - **Basitleştirilmiş Form**: Kullanıcı deneyimini iyileştirmek için gereksiz zorunlulukları kaldır
// - **Dinamik İş Kuralları**: Koşullara göre zorunluluk durumunu değiştir
// - **Adım Bazlı Formlar**: Farklı adımlarda farklı zorunluluk seviyeleri
//
// # Dönüş Değeri
//
// `*FieldUpdate`: Method chaining için güncellenmiş FieldUpdate pointer'ı
//
// # Kullanım Örneği
//
// ```go
// // Basit kullanım
// update := NewFieldUpdate().MakeOptional()
//
// // Koşullu opsiyonellik
// accountTypeField.OnChange(func(field *Schema, formData map[string]interface{}, ctx *fiber.Ctx) *FieldUpdate {
//     accountType := formData["account_type"]
//     if accountType == "personal" {
//         // Bireysel hesaplar için şirket adı opsiyonel
//         return NewFieldUpdate().
//             MakeOptional().
//             SetHelpText("Bireysel hesaplar için opsiyonel")
//     }
//     return NewFieldUpdate().
//         MakeRequired().
//         SetHelpText("Kurumsal hesaplar için zorunlu")
// })
//
// // Gizli alanlar için
// visibilityField.OnChange(func(field *Schema, formData map[string]interface{}, ctx *fiber.Ctx) *FieldUpdate {
//     visible := formData["show_field"].(bool)
//     if !visible {
//         // Gizli alanlar opsiyonel olmalı
//         return NewFieldUpdate().Hide().MakeOptional()
//     }
//     return NewFieldUpdate().Show().MakeRequired()
// })
//
// // Basitleştirilmiş form akışı
// quickModeField.OnChange(func(field *Schema, formData map[string]interface{}, ctx *fiber.Ctx) *FieldUpdate {
//     quickMode := formData["quick_mode"].(bool)
//     if quickMode {
//         // Hızlı mod: detaylı alanlar opsiyonel
//         return NewFieldUpdate().
//             MakeOptional().
//             SetHelpText("Hızlı modda bu alan opsiyoneldir")
//     }
//     return NewFieldUpdate().
//         MakeRequired().
//         SetHelpText("Detaylı modda bu alan zorunludur")
// })
// ```
//
// # Önemli Notlar
//
// - MakeOptional() çağrıldığında alan için validasyon zorunluluğu kaldırılır
// - Opsiyonel alanlar boş bırakılabilir
// - Gizli alanlar her zaman opsiyonel yapılmalıdır
// - Method chaining desteklenir
//
// # Best Practices
//
// - Hide() ile birlikte kullanarak tutarlılık sağlayın
// - Opsiyonel alanlar için açıklayıcı yardım metni ekleyin
// - Koşullu zorunlulukta her iki durumu da (Required/Optional) ele alın
// - Kullanıcı deneyimini önceliklendirin
//
// # Avantajlar
//
// - Kullanıcı deneyimini iyileştirir
// - Form doldurma süresini kısaltır
// - Esnek form yapısı sağlar
// - Gereksiz zorunlulukları kaldırır
func (u *FieldUpdate) MakeOptional() *FieldUpdate {
	required := false
	u.Required = &required
	return u
}

// Enable, alanı etkinleştirir (aktif hale getirir).
//
// # Genel Bakış
//
// Bu metod, devre dışı bırakılmış bir alanı tekrar etkinleştirir. Disabled özelliğini false
// olarak ayarlar ve method chaining desteği için *FieldUpdate döndürür. Etkin alanlar
// kullanıcı tarafından düzenlenebilir ve form submit edildiğinde değerleri gönderilir.
//
// # Kullanım Senaryoları
//
// - **Koşullu Aktivasyon**: Belirli koşullarda alanları aktif hale getir
// - **Adım Bazlı Formlar**: Form adımlarına göre alanları etkinleştir
// - **Yetki Bazlı Erişim**: Kullanıcı yetkisine göre alanları aktif et
// - **Dinamik Form Akışı**: Kullanıcı seçimlerine göre alanları etkinleştir
// - **Bağımlılık Çözümü**: Bağımlı alanları koşul sağlandığında aktif et
//
// # Dönüş Değeri
//
// `*FieldUpdate`: Method chaining için güncellenmiş FieldUpdate pointer'ı
//
// # Kullanım Örneği
//
// ```go
// // Basit kullanım
// update := NewFieldUpdate().Enable()
//
// // Koşullu aktivasyon
// agreementField.OnChange(func(field *Schema, formData map[string]interface{}, ctx *fiber.Ctx) *FieldUpdate {
//     agreed := formData["terms_agreed"].(bool)
//     if agreed {
//         // Sözleşme kabul edildiğinde devam butonu aktif
//         return NewFieldUpdate().
//             Enable().
//             SetHelpText("Devam edebilirsiniz")
//     }
//     return NewFieldUpdate().
//         Disable().
//         SetHelpText("Devam etmek için sözleşmeyi kabul edin")
// })
//
// // Bağımlılık çözümü
// countryField.OnChange(func(field *Schema, formData map[string]interface{}, ctx *fiber.Ctx) *FieldUpdate {
//     country := formData["country"]
//     if country != "" {
//         // Ülke seçildiğinde şehir alanını aktif et
//         return NewFieldUpdate().
//             Enable().
//             Show().
//             SetHelpText("Lütfen şehir seçin")
//     }
//     return NewFieldUpdate().Disable().Hide()
// })
//
// // Adım bazlı form
// stepField.OnChange(func(field *Schema, formData map[string]interface{}, ctx *fiber.Ctx) *FieldUpdate {
//     step := formData["current_step"].(int)
//     if step >= 2 {
//         // 2. adımda ödeme alanlarını aktif et
//         return NewFieldUpdate().
//             Enable().
//             Show().
//             MakeRequired()
//     }
//     return NewFieldUpdate().Disable().Hide()
// })
//
// // Yetki bazlı erişim
// userRole := ctx.Locals("user_role").(string)
// if userRole == "admin" || userRole == "editor" {
//     return NewFieldUpdate().
//         Enable().
//         MakeEditable().
//         SetHelpText("Düzenleme yetkisine sahipsiniz")
// }
// ```
//
// # Önemli Notlar
//
// - **Disabled vs ReadOnly**: Disabled alanlar form submit'e dahil edilmez, ReadOnly alanlar dahil edilir
// - **Görünürlük**: Enable() sadece disabled durumunu değiştirir, görünürlüğü etkilemez
// - **Validasyon**: Etkin alanlar için validasyon kuralları uygulanır
// - **Form Submit**: Etkin alanların değerleri form submit edildiğinde gönderilir
// - **Method Chaining**: Diğer metodlarla zincirleme kullanılabilir
//
// # Best Practices
//
// - Enable() ile birlikte Show() kullanarak tam erişim sağlayın
// - Etkinleştirilen alanlar için uygun yardım metni ekleyin
// - MakeEditable() ile birlikte kullanarak düzenleme izni verin
// - Koşullu aktivasyonda her iki durumu da (Enable/Disable) ele alın
//
// # Avantajlar
//
// - Dinamik form davranışı sağlar
// - Kullanıcı deneyimini iyileştirir
// - Bağımlılık yönetimini kolaylaştırır
// - Koşullu erişim kontrolü sağlar
//
// # Disabled vs ReadOnly Karşılaştırması
//
// | Özellik | Disabled | ReadOnly |
// |---------|----------|----------|
// | Görünürlük | Görünür (genelde soluk) | Görünür (normal) |
// | Düzenlenebilirlik | Hayır | Hayır |
// | Form Submit | Dahil edilmez | Dahil edilir |
// | Validasyon | Uygulanmaz | Uygulanmaz |
// | Kullanım Amacı | Geçici devre dışı | Kalıcı koruma |
func (u *FieldUpdate) Enable() *FieldUpdate {
	disabled := false
	u.Disabled = &disabled
	return u
}

// Disable, alanı devre dışı bırakır (pasif hale getirir).
//
// # Genel Bakış
//
// Bu metod, bir form alanını devre dışı bırakır. Disabled özelliğini true olarak ayarlar
// ve method chaining desteği için *FieldUpdate döndürür. Devre dışı alanlar kullanıcı
// tarafından düzenlenemez ve form submit edildiğinde değerleri gönderilmez.
//
// # Kullanım Senaryoları
//
// - **Koşullu Devre Dışı Bırakma**: Belirli koşullarda alanları pasif hale getir
// - **Bağımlılık Kontrolü**: Bağımlı alanları koşul sağlanmadığında devre dışı bırak
// - **Yetki Kısıtlaması**: Yetkisiz kullanıcılar için alanları kapat
// - **İş Kuralları**: İş mantığına göre alanları devre dışı bırak
// - **Geçici Kilitleme**: Belirli durumlarda alanları geçici olarak kilitle
//
// # Dönüş Değeri
//
// `*FieldUpdate`: Method chaining için güncellenmiş FieldUpdate pointer'ı
//
// # Kullanım Örneği
//
// ```go
// // Basit kullanım
// update := NewFieldUpdate().Disable()
//
// // Bağımlılık kontrolü
// countryField.OnChange(func(field *Schema, formData map[string]interface{}, ctx *fiber.Ctx) *FieldUpdate {
//     country := formData["country"]
//     if country == "" {
//         // Ülke seçilmediğinde şehir alanını devre dışı bırak
//         return NewFieldUpdate().
//             Disable().
//             SetHelpText("Önce ülke seçin").
//             SetValue(nil)
//     }
//     return NewFieldUpdate().Enable()
// })
//
// // Yetki kısıtlaması
// userRole := ctx.Locals("user_role").(string)
// if userRole == "viewer" {
//     return NewFieldUpdate().
//         Disable().
//         SetHelpText("Bu alanı düzenleme yetkiniz yok")
// }
//
// // İş kuralı bazlı devre dışı bırakma
// statusField.OnChange(func(field *Schema, formData map[string]interface{}, ctx *fiber.Ctx) *FieldUpdate {
//     status := formData["status"]
//     if status == "completed" {
//         // Tamamlanmış kayıtlar düzenlenemez
//         return NewFieldUpdate().
//             Disable().
//             MakeReadOnly().
//             SetHelpText("Tamamlanmış kayıtlar düzenlenemez")
//     }
//     return NewFieldUpdate().Enable().MakeEditable()
// })
//
// // Koşullu form akışı
// paymentMethodField.OnChange(func(field *Schema, formData map[string]interface{}, ctx *fiber.Ctx) *FieldUpdate {
//     paymentMethod := formData["payment_method"]
//     if paymentMethod == "cash" {
//         // Nakit ödeme için kredi kartı alanlarını devre dışı bırak
//         return NewFieldUpdate().
//             Disable().
//             Hide().
//             SetValue(nil)
//     }
//     return NewFieldUpdate().Enable().Show()
// })
//
// // Geçici kilitleme
// processingField.OnChange(func(field *Schema, formData map[string]interface{}, ctx *fiber.Ctx) *FieldUpdate {
//     isProcessing := formData["is_processing"].(bool)
//     if isProcessing {
//         return NewFieldUpdate().
//             Disable().
//             SetHelpText("İşlem devam ediyor, lütfen bekleyin...")
//     }
//     return NewFieldUpdate().Enable()
// })
// ```
//
// # Önemli Notlar
//
// - **Form Submit**: Disabled alanların değerleri form submit edildiğinde gönderilmez
// - **Validasyon**: Disabled alanlar için validasyon kuralları uygulanmaz
// - **Görünürlük**: Disable() sadece disabled durumunu değiştirir, görünürlüğü etkilemez
// - **Değer Temizleme**: Devre dışı bırakılan alanların değerlerini SetValue(nil) ile temizlemeyi düşünün
// - **Kullanıcı Deneyimi**: Neden devre dışı olduğunu açıklayan yardım metni ekleyin
//
// # Best Practices
//
// - Disable() ile birlikte SetHelpText() kullanarak kullanıcıyı bilgilendirin
// - Devre dışı bırakılan alanların değerlerini temizlemeyi düşünün
// - Hide() ile birlikte kullanarak gereksiz alanları tamamen gizleyin
// - Koşullu devre dışı bırakmada her iki durumu da (Enable/Disable) ele alın
// - Sunucu tarafında da aynı kontrolü uygulayın (güvenlik)
//
// # Avantajlar
//
// - Kullanıcı hatalarını önler
// - Form akışını kontrol eder
// - Bağımlılık yönetimini kolaylaştırır
// - Kullanıcı deneyimini iyileştirir
//
// # Dezavantajlar
//
// - Disabled alanların değerleri form submit'e dahil edilmez (veri kaybı riski)
// - Aşırı kullanımda kullanıcı deneyimini olumsuz etkiler
func (u *FieldUpdate) Disable() *FieldUpdate {
	disabled := true
	u.Disabled = &disabled
	return u
}

// SetHelpText, alanın yardım metnini ayarlar.
//
// # Genel Bakış
//
// Bu metod, bir form alanının altında veya yanında gösterilen yardım metnini dinamik olarak
// günceller. HelpText özelliğini verilen string ile ayarlar ve method chaining desteği için
// *FieldUpdate döndürür. Yardım metinleri kullanıcıya alan hakkında ek bilgi, format
// gereksinimleri veya örnekler sağlar.
//
// # Kullanım Senaryoları
//
// - **Bağlamsal Yardım**: Kullanıcı seçimlerine göre ilgili yardım metni göster
// - **Format Açıklaması**: Beklenen veri formatını açıkla (örn: "DD/MM/YYYY formatında")
// - **Validasyon Rehberi**: Validasyon kurallarını kullanıcıya açıkla
// - **Durum Bildirimi**: Alan durumu hakkında bilgi ver (örn: "Hesaplanıyor...")
// - **Hata Önleme**: Yaygın hataları önlemek için uyarılar ver
// - **Örnek Gösterme**: Geçerli değer örnekleri göster
//
// # Parametreler
//
// - `text`: Gösterilecek yardım metni (string)
//
// # Dönüş Değeri
//
// `*FieldUpdate`: Method chaining için güncellenmiş FieldUpdate pointer'ı
//
// # Kullanım Örneği
//
// ```go
// // Basit kullanım
// update := NewFieldUpdate().SetHelpText("Bu alan zorunludur")
//
// // Bağlamsal yardım
// countryField.OnChange(func(field *Schema, formData map[string]interface{}, ctx *fiber.Ctx) *FieldUpdate {
//     country := formData["country"]
//     if country == "TR" {
//         return NewFieldUpdate().
//             SetHelpText("11 haneli TC Kimlik Numaranızı girin (örn: 12345678901)")
//     } else if country == "US" {
//         return NewFieldUpdate().
//             SetHelpText("9 haneli SSN numaranızı girin (örn: 123-45-6789)")
//     }
//     return NewFieldUpdate().SetHelpText("Kimlik numaranızı girin")
// })
//
// // Format açıklaması
// dateField.OnChange(func(field *Schema, formData map[string]interface{}, ctx *fiber.Ctx) *FieldUpdate {
//     return NewFieldUpdate().
//         SetHelpText("Tarih formatı: GG/AA/YYYY (örn: 15/03/2024)")
// })
//
// // Validasyon rehberi
// passwordField.OnChange(func(field *Schema, formData map[string]interface{}, ctx *fiber.Ctx) *FieldUpdate {
//     return NewFieldUpdate().
//         SetHelpText("En az 8 karakter, 1 büyük harf, 1 küçük harf ve 1 rakam içermelidir")
// })
//
// // Durum bildirimi
// calculatingField.OnChange(func(field *Schema, formData map[string]interface{}, ctx *fiber.Ctx) *FieldUpdate {
//     isCalculating := formData["is_calculating"].(bool)
//     if isCalculating {
//         return NewFieldUpdate().
//             SetHelpText("⏳ Hesaplama yapılıyor, lütfen bekleyin...").
//             Disable()
//     }
//     return NewFieldUpdate().
//         SetHelpText("Toplam tutar otomatik hesaplanmıştır").
//         Enable()
// })
//
// // Koşullu uyarı
// amountField.OnChange(func(field *Schema, formData map[string]interface{}, ctx *fiber.Ctx) *FieldUpdate {
//     amount := formData["amount"].(float64)
//     if amount > 10000 {
//         return NewFieldUpdate().
//             SetHelpText("⚠️ Yüksek tutar: Ek onay gerekebilir")
//     }
//     return NewFieldUpdate().SetHelpText("Tutar giriniz")
// })
//
// // Çoklu dil desteği
// language := ctx.Locals("language").(string)
// if language == "tr" {
//     return NewFieldUpdate().SetHelpText("E-posta adresinizi girin")
// } else {
//     return NewFieldUpdate().SetHelpText("Enter your email address")
// }
// ```
//
// # Önemli Notlar
//
// - **Markdown Desteği**: Bazı UI framework'leri markdown formatını destekleyebilir
// - **Uzunluk**: Yardım metinleri kısa ve öz olmalıdır (ideal: 1-2 cümle)
// - **Emoji Kullanımı**: Dikkat çekmek için emoji kullanılabilir ancak aşırıya kaçılmamalı
// - **Dinamik Güncelleme**: Yardım metni form durumuna göre dinamik olarak güncellenebilir
// - **Erişilebilirlik**: Screen reader'lar için anlamlı metinler yazın
//
// # Best Practices
//
// - Kısa, açık ve anlaşılır metinler yazın
// - Örnekler vererek kullanıcıyı yönlendirin
// - Hata mesajları yerine önleyici bilgiler verin
// - Teknik jargondan kaçının, sade dil kullanın
// - Çoklu dil desteği için i18n kullanmayı düşünün
// - Kullanıcının mevcut bağlamına uygun metinler gösterin
//
// # Avantajlar
//
// - Kullanıcı deneyimini iyileştirir
// - Hata oranını azaltır
// - Form doldurma süresini kısaltır
// - Kullanıcı güvenini artırır
// - Destek taleplerini azaltır
func (u *FieldUpdate) SetHelpText(text string) *FieldUpdate {
	u.HelpText = &text
	return u
}

// SetPlaceholder, alanın placeholder (yer tutucu) metnini ayarlar.
//
// # Genel Bakış
//
// Bu metod, bir form alanının içinde gösterilen placeholder metnini dinamik olarak günceller.
// Placeholder özelliğini verilen string ile ayarlar ve method chaining desteği için
// *FieldUpdate döndürür. Placeholder metinleri, alan boşken kullanıcıya ne girmesi
// gerektiği hakkında ipucu verir.
//
// # Kullanım Senaryoları
//
// - **Format Örneği**: Beklenen format için örnek göster (örn: "örnek@email.com")
// - **Değer Önerisi**: Tipik değer örnekleri göster (örn: "Ahmet Yılmaz")
// - **Arama İpucu**: Arama alanlarında ne aranabileceğini göster
// - **Dinamik İpuçları**: Kullanıcı seçimlerine göre farklı örnekler göster
// - **Birim Belirtme**: Sayısal alanlar için birim göster (örn: "0.00 TL")
//
// # Parametreler
//
// - `text`: Gösterilecek placeholder metni (string)
//
// # Dönüş Değeri
//
// `*FieldUpdate`: Method chaining için güncellenmiş FieldUpdate pointer'ı
//
// # Kullanım Örneği
//
// ```go
// // Basit kullanım
// update := NewFieldUpdate().SetPlaceholder("Adınızı girin")
//
// // Format örneği
// emailField := NewFieldUpdate().
//     SetPlaceholder("ornek@email.com").
//     SetHelpText("Geçerli bir e-posta adresi girin")
//
// // Dinamik placeholder
// searchTypeField.OnChange(func(field *Schema, formData map[string]interface{}, ctx *fiber.Ctx) *FieldUpdate {
//     searchType := formData["search_type"]
//     switch searchType {
//     case "name":
//         return NewFieldUpdate().SetPlaceholder("Örn: Ahmet Yılmaz")
//     case "email":
//         return NewFieldUpdate().SetPlaceholder("Örn: ahmet@example.com")
//     case "phone":
//         return NewFieldUpdate().SetPlaceholder("Örn: 0532 123 45 67")
//     default:
//         return NewFieldUpdate().SetPlaceholder("Arama...")
//     }
// })
//
// // Birim belirtme
// priceField := NewFieldUpdate().
//     SetPlaceholder("0.00 TL").
//     SetHelpText("Ürün fiyatını girin")
//
// // Bağlamsal örnek
// countryField.OnChange(func(field *Schema, formData map[string]interface{}, ctx *fiber.Ctx) *FieldUpdate {
//     country := formData["country"]
//     if country == "TR" {
//         return NewFieldUpdate().
//             SetPlaceholder("İstanbul, Ankara, İzmir...").
//             SetHelpText("Türkiye'deki bir şehir seçin")
//     } else if country == "US" {
//         return NewFieldUpdate().
//             SetPlaceholder("New York, Los Angeles, Chicago...").
//             SetHelpText("Select a US city")
//     }
//     return NewFieldUpdate().SetPlaceholder("Şehir seçin...")
// })
//
// // Arama alanı
// searchField := NewFieldUpdate().
//     SetPlaceholder("Ürün adı, kategori veya marka ara...").
//     SetHelpText("En az 3 karakter girin")
//
// // Sayısal alan
// quantityField := NewFieldUpdate().
//     SetPlaceholder("1").
//     SetHelpText("Sipariş miktarını girin (minimum: 1)")
//
// // Tarih alanı
// dateField := NewFieldUpdate().
//     SetPlaceholder("GG/AA/YYYY").
//     SetHelpText("Doğum tarihinizi girin")
//
// // Çoklu dil desteği
// language := ctx.Locals("language").(string)
// if language == "tr" {
//     return NewFieldUpdate().SetPlaceholder("Mesajınızı yazın...")
// } else {
//     return NewFieldUpdate().SetPlaceholder("Type your message...")
// }
// ```
//
// # Önemli Notlar
//
// - **Görünürlük**: Placeholder sadece alan boşken görünür, değer girildiğinde kaybolur
// - **Erişilebilirlik**: Placeholder, label'ın yerini tutmamalıdır (her ikisi de olmalı)
// - **Uzunluk**: Kısa ve öz olmalıdır (ideal: 2-5 kelime)
// - **Renk**: Genelde soluk renkte gösterilir, asıl değerle karıştırılmamalı
// - **Validasyon**: Placeholder metni validasyon kuralı değildir, sadece ipucudur
//
// # Best Practices
//
// - Gerçekçi ve anlamlı örnekler kullanın
// - "Örn:" veya "..." gibi önekler ekleyerek placeholder olduğunu belirtin
// - Label ile tekrar etmeyin, ek bilgi verin
// - Kısa ve öz tutun (uzun metinler kesilir)
// - Birim veya format bilgisi içeren örnekler verin
// - Çoklu dil desteği için i18n kullanın
//
// # Placeholder vs HelpText Karşılaştırması
//
// | Özellik | Placeholder | HelpText |
// |---------|-------------|----------|
// | Konum | Alan içinde | Alan dışında (altında/yanında) |
// | Görünürlük | Sadece boşken | Her zaman |
// | Uzunluk | Kısa (2-5 kelime) | Orta (1-2 cümle) |
// | Amaç | Örnek göster | Açıklama yap |
// | Erişilebilirlik | Düşük | Yüksek |
//
// # Avantajlar
//
// - Kullanıcıya hızlı ipucu verir
// - Form alanını temiz tutar
// - Örnek göstererek anlaşılırlığı artırır
// - Kullanıcı hatalarını azaltır
func (u *FieldUpdate) SetPlaceholder(text string) *FieldUpdate {
	u.Placeholder = &text
	return u
}

// SetOptions, alanın seçeneklerini ayarlar (select, radio, checkbox gibi alanlar için).
//
// # Genel Bakış
//
// Bu metod, select, radio button, checkbox gibi seçim alanlarının seçeneklerini dinamik olarak
// günceller. Options özelliğini verilen map ile ayarlar ve method chaining desteği için
// *FieldUpdate döndürür. Bu özellik, cascade dropdown'lar ve dinamik seçim listeleri için
// kritik öneme sahiptir.
//
// # Kullanım Senaryoları
//
// - **Cascade Dropdown**: Üst seçime göre alt dropdown seçeneklerini filtrele
// - **Dinamik Seçenekler**: Kullanıcı seçimlerine göre seçenekleri değiştir
// - **Veritabanı Sorguları**: Context üzerinden veritabanından seçenekleri çek
// - **Koşullu Seçenekler**: İş kurallarına göre farklı seçenekler göster
// - **Filtrelenmiş Listeler**: Belirli kriterlere göre seçenekleri filtrele
// - **Rol Bazlı Seçenekler**: Kullanıcı rolüne göre farklı seçenekler sun
//
// # Parametreler
//
// - `options`: Seçenekleri içeren map[string]interface{}
//   - Key: Seçeneğin değeri (value)
//   - Value: Seçeneğin görünen metni (label) veya karmaşık nesne
//
// # Dönüş Değeri
//
// `*FieldUpdate`: Method chaining için güncellenmiş FieldUpdate pointer'ı
//
// # Kullanım Örneği
//
// ```go
// // Basit kullanım
// cities := map[string]interface{}{
//     "istanbul": "İstanbul",
//     "ankara": "Ankara",
//     "izmir": "İzmir",
// }
// update := NewFieldUpdate().SetOptions(cities)
//
// // Cascade dropdown - Ülkeye göre şehirler
// countryField.OnChange(func(field *Schema, formData map[string]interface{}, ctx *fiber.Ctx) *FieldUpdate {
//     country := formData["country"].(string)
//
//     var cities map[string]interface{}
//     switch country {
//     case "TR":
//         cities = map[string]interface{}{
//             "istanbul": "İstanbul",
//             "ankara": "Ankara",
//             "izmir": "İzmir",
//             "bursa": "Bursa",
//         }
//     case "US":
//         cities = map[string]interface{}{
//             "new_york": "New York",
//             "los_angeles": "Los Angeles",
//             "chicago": "Chicago",
//         }
//     default:
//         cities = map[string]interface{}{}
//     }
//
//     return NewFieldUpdate().
//         SetOptions(cities).
//         Show().
//         Enable().
//         SetValue(nil) // Önceki seçimi temizle
// })
//
// // Veritabanından dinamik seçenekler
// categoryField.OnChange(func(field *Schema, formData map[string]interface{}, ctx *fiber.Ctx) *FieldUpdate {
//     categoryID := formData["category_id"].(int)
//
//     // Veritabanından alt kategorileri çek
//     var subcategories []Subcategory
//     db := ctx.Locals("db").(*gorm.DB)
//     db.Where("category_id = ?", categoryID).Find(&subcategories)
//
//     // Map'e dönüştür
//     options := make(map[string]interface{})
//     for _, sub := range subcategories {
//         options[fmt.Sprintf("%d", sub.ID)] = sub.Name
//     }
//
//     return NewFieldUpdate().
//         SetOptions(options).
//         Show().
//         MakeRequired().
//         SetHelpText(fmt.Sprintf("%d alt kategori bulundu", len(options)))
// })
//
// // Karmaşık seçenek nesneleri
// productField.OnChange(func(field *Schema, formData map[string]interface{}, ctx *fiber.Ctx) *FieldUpdate {
//     options := map[string]interface{}{
//         "1": map[string]interface{}{
//             "label": "Premium Paket",
//             "price": 1000,
//             "description": "Tüm özellikler dahil",
//             "icon": "star",
//         },
//         "2": map[string]interface{}{
//             "label": "Standart Paket",
//             "price": 500,
//             "description": "Temel özellikler",
//             "icon": "check",
//         },
//     }
//
//     return NewFieldUpdate().SetOptions(options)
// })
//
// // Filtrelenmiş seçenekler
// priceRangeField.OnChange(func(field *Schema, formData map[string]interface{}, ctx *fiber.Ctx) *FieldUpdate {
//     priceRange := formData["price_range"].(string)
//
//     var products []Product
//     db := ctx.Locals("db").(*gorm.DB)
//
//     switch priceRange {
//     case "low":
//         db.Where("price < ?", 100).Find(&products)
//     case "medium":
//         db.Where("price BETWEEN ? AND ?", 100, 500).Find(&products)
//     case "high":
//         db.Where("price > ?", 500).Find(&products)
//     }
//
//     options := make(map[string]interface{})
//     for _, p := range products {
//         options[fmt.Sprintf("%d", p.ID)] = fmt.Sprintf("%s - %.2f TL", p.Name, p.Price)
//     }
//
//     return NewFieldUpdate().SetOptions(options)
// })
//
// // Rol bazlı seçenekler
// userRole := ctx.Locals("user_role").(string)
// var statusOptions map[string]interface{}
//
// if userRole == "admin" {
//     statusOptions = map[string]interface{}{
//         "draft": "Taslak",
//         "pending": "Beklemede",
//         "approved": "Onaylandı",
//         "rejected": "Reddedildi",
//         "archived": "Arşivlendi",
//     }
// } else {
//     statusOptions = map[string]interface{}{
//         "draft": "Taslak",
//         "pending": "Beklemede",
//     }
// }
//
// return NewFieldUpdate().SetOptions(statusOptions)
// ```
//
// # Önemli Notlar
//
// - **Map Formatı**: Key-value çiftleri şeklinde olmalıdır
// - **Değer Temizleme**: Seçenekler değiştiğinde mevcut değeri SetValue(nil) ile temizleyin
// - **Performans**: Büyük seçenek listeleri için pagination veya lazy loading düşünün
// - **Veritabanı Sorguları**: Context üzerinden DB erişimi yapılabilir ancak performans etkileri göz önünde bulundurulmalıdır
// - **Boş Seçenekler**: Boş map göndermek seçenekleri temizler
//
// # Best Practices
//
// - Cascade dropdown'larda üst seçim değiştiğinde alt seçimi temizleyin
// - Veritabanı sorgularını optimize edin (index, limit kullanın)
// - Büyük listelerde arama/filtreleme özelliği ekleyin
// - Seçenek sayısını kullanıcıya bildirin (yardım metni ile)
// - Hata durumlarını ele alın (veritabanı hatası, boş sonuç vb.)
// - Cache mekanizması kullanarak performansı artırın
//
// # Avantajlar
//
// - Dinamik ve esnek seçim listeleri sağlar
// - Cascade dropdown'ları kolayca uygular
// - Kullanıcı deneyimini iyileştirir
// - İş kurallarını kolayca uygular
// - Veritabanı entegrasyonu kolaydır
//
// # Dezavantajlar
//
// - Büyük seçenek listeleri performans sorunlarına yol açabilir
// - Veritabanı sorguları her değişiklikte çalışabilir (performans)
// - Karmaşık bağımlılık zincirleri debug edilmesi zor olabilir
//
// # Performans İpuçları
//
// - Veritabanı sorgularını cache'leyin
// - Lazy loading kullanın (büyük listeler için)
// - Pagination ekleyin (1000+ seçenek için)
// - Debounce mekanizması kullanın (arama için)
// - Index'leri optimize edin
func (u *FieldUpdate) SetOptions(options map[string]interface{}) *FieldUpdate {
	u.Options = options
	return u
}

// SetValue, alanın değerini programatik olarak ayarlar.
//
// # Genel Bakış
//
// Bu metod, bir form alanının değerini programatik olarak ayarlar. Value özelliğini verilen
// değer ile günceller ve method chaining desteği için *FieldUpdate döndürür. Bu özellik,
// hesaplanan değerler, otomatik doldurma ve değer senkronizasyonu için kullanılır.
//
// # Kullanım Senaryoları
//
// - **Hesaplanan Değerler**: Diğer alanların değerlerine göre otomatik hesaplama yap
// - **Otomatik Doldurma**: Kullanıcı seçimlerine göre alanları otomatik doldur
// - **Değer Senkronizasyonu**: Bir alanın değerini diğer alanlara kopyala
// - **Varsayılan Değerler**: Koşullara göre varsayılan değerler ata
// - **Değer Temizleme**: Belirli koşullarda alanları temizle (nil)
// - **Format Dönüşümü**: Bir formattan diğerine dönüştür
//
// # Parametreler
//
// - `value`: Atanacak değer (interface{} - herhangi bir tip olabilir)
//   - string, int, float64, bool, nil, map, slice vb.
//
// # Dönüş Değeri
//
// `*FieldUpdate`: Method chaining için güncellenmiş FieldUpdate pointer'ı
//
// # Kullanım Örneği
//
// ```go
// // Basit kullanım
// update := NewFieldUpdate().SetValue("Yeni değer")
//
// // Değer temizleme
// update := NewFieldUpdate().SetValue(nil)
//
// // Hesaplanan değer - Toplam hesaplama
// quantityField.OnChange(func(field *Schema, formData map[string]interface{}, ctx *fiber.Ctx) *FieldUpdate {
//     quantity := formData["quantity"].(float64)
//     price := formData["price"].(float64)
//     total := quantity * price
//
//     // Toplam alanını güncelle
//     return NewFieldUpdate().
//         SetValue(total).
//         MakeReadOnly().
//         SetHelpText(fmt.Sprintf("Toplam: %.2f TL", total))
// })
//
// // Otomatik doldurma - Şehir seçimine göre posta kodu
// cityField.OnChange(func(field *Schema, formData map[string]interface{}, ctx *fiber.Ctx) *FieldUpdate {
//     city := formData["city"].(string)
//
//     // Şehir kodlarını map'ten al
//     cityCodes := map[string]string{
//         "istanbul": "34",
//         "ankara": "06",
//         "izmir": "35",
//     }
//
//     postalCode := cityCodes[city]
//     return NewFieldUpdate().
//         SetValue(postalCode).
//         SetHelpText(fmt.Sprintf("%s için posta kodu: %s", city, postalCode))
// })
//
// // Değer senkronizasyonu - E-posta'yı kullanıcı adına kopyala
// emailField.OnChange(func(field *Schema, formData map[string]interface{}, ctx *fiber.Ctx) *FieldUpdate {
//     email := formData["email"].(string)
//     autoFillUsername := formData["auto_fill_username"].(bool)
//
//     if autoFillUsername {
//         // E-posta'nın @ öncesi kısmını kullanıcı adı olarak ata
//         username := strings.Split(email, "@")[0]
//         return NewFieldUpdate().
//             SetValue(username).
//             SetHelpText("E-posta adresinizden otomatik oluşturuldu")
//     }
//     return nil
// })
//
// // Cascade seçim temizleme
// countryField.OnChange(func(field *Schema, formData map[string]interface{}, ctx *fiber.Ctx) *FieldUpdate {
//     // Ülke değiştiğinde şehir seçimini temizle
//     return NewFieldUpdate().
//         SetValue(nil).
//         SetHelpText("Lütfen yeni ülkeye göre şehir seçin")
// })
//
// // Format dönüşümü - Telefon numarası formatlama
// phoneField.OnChange(func(field *Schema, formData map[string]interface{}, ctx *fiber.Ctx) *FieldUpdate {
//     phone := formData["phone"].(string)
//
//     // Sadece rakamları al
//     digits := regexp.MustCompile(`\D`).ReplaceAllString(phone, "")
//
//     // Formatla: (0532) 123 45 67
//     if len(digits) == 11 {
//         formatted := fmt.Sprintf("(%s) %s %s %s",
//             digits[0:4], digits[4:7], digits[7:9], digits[9:11])
//         return NewFieldUpdate().
//             SetValue(formatted).
//             SetHelpText("Telefon numarası formatlandı")
//     }
//     return nil
// })
//
// // Koşullu varsayılan değer
// paymentMethodField.OnChange(func(field *Schema, formData map[string]interface{}, ctx *fiber.Ctx) *FieldUpdate {
//     paymentMethod := formData["payment_method"].(string)
//
//     if paymentMethod == "installment" {
//         // Taksitli ödeme seçildiğinde varsayılan taksit sayısı
//         return NewFieldUpdate().
//             SetValue(3).
//             SetHelpText("Varsayılan: 3 taksit")
//     }
//     return NewFieldUpdate().SetValue(nil)
// })
//
// // Veritabanından değer çekme
// productField.OnChange(func(field *Schema, formData map[string]interface{}, ctx *fiber.Ctx) *FieldUpdate {
//     productID := formData["product_id"].(int)
//
//     // Ürün fiyatını veritabanından çek
//     var product Product
//     db := ctx.Locals("db").(*gorm.DB)
//     db.First(&product, productID)
//
//     // Fiyat alanını otomatik doldur
//     return NewFieldUpdate().
//         SetValue(product.Price).
//         MakeReadOnly().
//         SetHelpText(fmt.Sprintf("Ürün fiyatı: %.2f TL", product.Price))
// })
//
// // Karmaşık nesne değeri
// addressField.OnChange(func(field *Schema, formData map[string]interface{}, ctx *fiber.Ctx) *FieldUpdate {
//     addressData := map[string]interface{}{
//         "street": "Atatürk Caddesi",
//         "city": "İstanbul",
//         "postal_code": "34000",
//         "country": "TR",
//     }
//
//     return NewFieldUpdate().SetValue(addressData)
// })
// ```
//
// # Önemli Notlar
//
// - **Tip Güvenliği**: interface{} kullanıldığı için tip dönüşümlerine dikkat edin
// - **Nil Değer**: nil göndermek alanı temizler
// - **Validasyon**: SetValue ile atanan değerler de validasyon kurallarına tabidir
// - **Kullanıcı Girişi**: Programatik değer ataması kullanıcı girişini ezer
// - **ReadOnly ile Kullanım**: Hesaplanan değerler için MakeReadOnly() ile birlikte kullanın
//
// # Best Practices
//
// - Hesaplanan değerler için MakeReadOnly() kullanın
// - Tip dönüşümlerinde hata kontrolü yapın
// - Cascade seçimlerde önceki değeri temizleyin (nil)
// - Kullanıcıya değerin nereden geldiğini açıklayın (yardım metni)
// - Veritabanı sorgularında hata kontrolü yapın
// - Format dönüşümlerinde validasyon yapın
//
// # Avantajlar
//
// - Otomatik hesaplama sağlar
// - Kullanıcı deneyimini iyileştirir
// - Veri tutarlılığını sağlar
// - Form doldurma süresini kısaltır
// - Hata oranını azaltır
//
// # Dezavantajlar
//
// - Yanlış kullanımda kullanıcı girişini ezebilir
// - Tip güvenliği yoktur (interface{})
// - Karmaşık hesaplamalarda performans sorunları olabilir
//
// # Güvenlik Notları
//
// - Kullanıcı girişini doğrudan kullanmadan önce sanitize edin
// - Veritabanı sorgularında SQL injection'a dikkat edin
// - Hassas bilgileri (şifre, kredi kartı) SetValue ile atamayın
// - Sunucu tarafında da aynı hesaplamaları yapın (güvenlik)
func (u *FieldUpdate) SetValue(value interface{}) *FieldUpdate {
	u.Value = value
	return u
}

// SetRules, alanın validasyon kurallarını ayarlar.
//
// # Genel Bakış
//
// Bu metod, bir form alanının tüm validasyon kurallarını tek seferde ayarlar. Mevcut kuralları
// tamamen değiştirir. Rules özelliğini verilen slice ile günceller ve method chaining desteği
// için *FieldUpdate döndürür. Dinamik validasyon senaryoları için kritik öneme sahiptir.
//
// # Kullanım Senaryoları
//
// - **Koşullu Validasyon**: Kullanıcı seçimlerine göre farklı validasyon kuralları uygula
// - **Dinamik Kural Setleri**: İş kurallarına göre validasyon kurallarını değiştir
// - **Rol Bazlı Validasyon**: Kullanıcı rolüne göre farklı validasyon seviyeleri
// - **Adım Bazlı Validasyon**: Form adımlarına göre farklı kurallar uygula
// - **Toplu Kural Güncelleme**: Birden fazla kuralı aynı anda güncelle
//
// # Parametreler
//
// - `rules`: ValidationRule slice'ı - Uygulanacak tüm validasyon kuralları
//
// # Dönüş Değeri
//
// `*FieldUpdate`: Method chaining için güncellenmiş FieldUpdate pointer'ı
//
// # Kullanım Örneği
//
// ```go
// // Basit kullanım
// rules := []ValidationRule{
//     {Type: "required", Message: "Bu alan zorunludur"},
//     {Type: "min", Value: 3, Message: "En az 3 karakter olmalıdır"},
//     {Type: "max", Value: 50, Message: "En fazla 50 karakter olabilir"},
// }
// update := NewFieldUpdate().SetRules(rules)
//
// // Koşullu validasyon - Kullanıcı tipine göre
// userTypeField.OnChange(func(field *Schema, formData map[string]interface{}, ctx *fiber.Ctx) *FieldUpdate {
//     userType := formData["user_type"].(string)
//
//     var rules []ValidationRule
//     if userType == "corporate" {
//         // Kurumsal kullanıcılar için sıkı kurallar
//         rules = []ValidationRule{
//             {Type: "required", Message: "Vergi numarası zorunludur"},
//             {Type: "length", Value: 10, Message: "Vergi numarası 10 haneli olmalıdır"},
//             {Type: "numeric", Message: "Sadece rakam içermelidir"},
//             {Type: "tax_number", Message: "Geçerli bir vergi numarası değil"},
//         }
//     } else {
//         // Bireysel kullanıcılar için esnek kurallar
//         rules = []ValidationRule{
//             {Type: "length", Value: 11, Message: "TC Kimlik No 11 haneli olmalıdır"},
//             {Type: "numeric", Message: "Sadece rakam içermelidir"},
//         }
//     }
//
//     return NewFieldUpdate().
//         SetRules(rules).
//         MakeRequired().
//         SetHelpText(fmt.Sprintf("%s için kimlik bilgisi", userType))
// })
//
// // Adım bazlı validasyon
// stepField.OnChange(func(field *Schema, formData map[string]interface{}, ctx *fiber.Ctx) *FieldUpdate {
//     step := formData["current_step"].(int)
//
//     var emailRules []ValidationRule
//     if step == 1 {
//         // İlk adımda basit validasyon
//         emailRules = []ValidationRule{
//             {Type: "required", Message: "E-posta zorunludur"},
//             {Type: "email", Message: "Geçerli bir e-posta adresi girin"},
//         }
//     } else if step == 2 {
//         // İkinci adımda detaylı validasyon
//         emailRules = []ValidationRule{
//             {Type: "required", Message: "E-posta zorunludur"},
//             {Type: "email", Message: "Geçerli bir e-posta adresi girin"},
//             {Type: "unique", Table: "users", Column: "email", Message: "Bu e-posta zaten kullanılıyor"},
//             {Type: "domain_whitelist", Value: []string{"company.com", "partner.com"}, Message: "Sadece şirket e-postaları kabul edilir"},
//         }
//     }
//
//     return NewFieldUpdate().SetRules(emailRules)
// })
//
// // Rol bazlı validasyon
// userRole := ctx.Locals("user_role").(string)
// var priceRules []ValidationRule
//
// if userRole == "admin" {
//     // Admin için sınırsız
//     priceRules = []ValidationRule{
//         {Type: "required", Message: "Fiyat zorunludur"},
//         {Type: "numeric", Message: "Geçerli bir sayı girin"},
//         {Type: "min", Value: 0, Message: "Fiyat negatif olamaz"},
//     }
// } else if userRole == "manager" {
//     // Manager için sınırlı
//     priceRules = []ValidationRule{
//         {Type: "required", Message: "Fiyat zorunludur"},
//         {Type: "numeric", Message: "Geçerli bir sayı girin"},
//         {Type: "min", Value: 0, Message: "Fiyat negatif olamaz"},
//         {Type: "max", Value: 10000, Message: "Maksimum 10.000 TL girebilirsiniz"},
//     }
// } else {
//     // Normal kullanıcı için çok sınırlı
//     priceRules = []ValidationRule{
//         {Type: "required", Message: "Fiyat zorunludur"},
//         {Type: "numeric", Message: "Geçerli bir sayı girin"},
//         {Type: "min", Value: 0, Message: "Fiyat negatif olamaz"},
//         {Type: "max", Value: 1000, Message: "Maksimum 1.000 TL girebilirsiniz"},
//     }
// }
//
// return NewFieldUpdate().
//     SetRules(priceRules).
//     SetHelpText(fmt.Sprintf("Rol: %s - Fiyat limitleri uygulanıyor", userRole))
//
// // Ürün tipine göre validasyon
// productTypeField.OnChange(func(field *Schema, formData map[string]interface{}, ctx *fiber.Ctx) *FieldUpdate {
//     productType := formData["product_type"].(string)
//
//     var descriptionRules []ValidationRule
//     if productType == "premium" {
//         // Premium ürünler için detaylı açıklama zorunlu
//         descriptionRules = []ValidationRule{
//             {Type: "required", Message: "Premium ürünler için açıklama zorunludur"},
//             {Type: "min", Value: 100, Message: "En az 100 karakter olmalıdır"},
//             {Type: "max", Value: 5000, Message: "En fazla 5000 karakter olabilir"},
//             {Type: "rich_text", Message: "Zengin metin formatı gereklidir"},
//         }
//     } else {
//         // Standart ürünler için basit açıklama
//         descriptionRules = []ValidationRule{
//             {Type: "min", Value: 20, Message: "En az 20 karakter olmalıdır"},
//             {Type: "max", Value: 500, Message: "En fazla 500 karakter olabilir"},
//         }
//     }
//
//     return NewFieldUpdate().SetRules(descriptionRules)
// })
//
// // Boş kural seti (tüm validasyonları kaldır)
// update := NewFieldUpdate().SetRules([]ValidationRule{})
// ```
//
// # Önemli Notlar
//
// - **Mevcut Kuralları Değiştirir**: SetRules() mevcut tüm kuralları siler ve yenileriyle değiştirir
// - **Boş Slice**: Boş slice göndermek tüm validasyon kurallarını kaldırır
// - **Kural Sırası**: Kurallar verilen sırada uygulanır
// - **Performans**: Çok fazla kural performansı etkileyebilir
// - **Sunucu Tarafı**: Client-side validasyon yeterli değildir, sunucu tarafında da validasyon yapın
//
// # ValidationRule Yapısı
//
// ```go
// type ValidationRule struct {
//     Type    string      // Kural tipi: "required", "min", "max", "email", "numeric" vb.
//     Value   interface{} // Kural değeri (min/max için sayı, pattern için regex vb.)
//     Message string      // Hata mesajı
//     Table   string      // Unique kontrolü için tablo adı
//     Column  string      // Unique kontrolü için kolon adı
// }
// ```
//
// # Yaygın Kural Tipleri
//
// - `required`: Alan zorunludur
// - `min`: Minimum değer/uzunluk
// - `max`: Maximum değer/uzunluk
// - `email`: E-posta formatı
// - `numeric`: Sayısal değer
// - `alpha`: Sadece harf
// - `alphanumeric`: Harf ve rakam
// - `url`: URL formatı
// - `regex`: Regex pattern
// - `unique`: Veritabanında benzersiz
// - `length`: Sabit uzunluk
// - `in`: Belirli değerlerden biri
// - `not_in`: Belirli değerlerden hiçbiri
//
// # Best Practices
//
// - Kullanıcı dostu hata mesajları yazın
// - Gereksiz kurallardan kaçının (performans)
// - Kuralları mantıksal sıraya göre düzenleyin (required önce, detaylı kurallar sonra)
// - Sunucu tarafında da aynı kuralları uygulayın
// - Karmaşık kurallar için custom validator yazın
// - Hata mesajlarında örnek verin
//
// # Avantajlar
//
// - Dinamik ve esnek validasyon sağlar
// - İş kurallarını kolayca uygular
// - Kullanıcı deneyimini iyileştirir
// - Veri bütünlüğünü sağlar
// - Toplu kural güncellemesi yapar
//
// # Dezavantajlar
//
// - Mevcut kuralları tamamen değiştirir (dikkatli kullanılmalı)
// - Çok fazla kural performansı etkileyebilir
// - Karmaşık kural setleri yönetimi zorlaştırabilir
//
// # SetRules vs AddRule Karşılaştırması
//
// | Özellik | SetRules | AddRule |
// |---------|----------|---------|
// | Mevcut Kurallar | Siler ve değiştirir | Korur ve ekler |
// | Kullanım | Toplu güncelleme | Tekli ekleme |
// | Performans | Tek seferde | Her çağrıda |
// | Esneklik | Düşük | Yüksek |
// | Kullanım Amacı | Tüm kuralları değiştir | Kural ekle |
func (u *FieldUpdate) SetRules(rules []ValidationRule) *FieldUpdate {
	u.Rules = rules
	return u
}

// AddRule, mevcut validasyon kurallarına yeni bir kural ekler.
//
// # Genel Bakış
//
// Bu metod, bir form alanının mevcut validasyon kurallarına yeni bir kural ekler. Mevcut
// kuralları korur ve sonuna yeni kuralı ekler. Method chaining desteği için *FieldUpdate
// döndürür. Kademeli kural ekleme senaryoları için idealdir.
//
// # Kullanım Senaryoları
//
// - **Kademeli Kural Ekleme**: Koşullara göre adım adım kural ekle
// - **Ek Validasyon**: Mevcut kurallara ek kontroller ekle
// - **Koşullu Kural Ekleme**: Belirli durumlarda ekstra kurallar ekle
// - **Dinamik Kural Biriktirme**: Birden fazla koşula göre kuralları biriktir
// - **Modüler Validasyon**: Her koşul için ayrı kural ekle
//
// # Parametreler
//
// - `rule`: Eklenecek ValidationRule - Tek bir validasyon kuralı
//
// # Dönüş Değeri
//
// `*FieldUpdate`: Method chaining için güncellenmiş FieldUpdate pointer'ı
//
// # Kullanım Örneği
//
// ```go
// // Basit kullanım
// update := NewFieldUpdate().
//     AddRule(ValidationRule{Type: "required", Message: "Bu alan zorunludur"}).
//     AddRule(ValidationRule{Type: "min", Value: 3, Message: "En az 3 karakter"})
//
// // Koşullu kural ekleme
// passwordField.OnChange(func(field *Schema, formData map[string]interface{}, ctx *fiber.Ctx) *FieldUpdate {
//     update := NewFieldUpdate().
//         AddRule(ValidationRule{Type: "required", Message: "Şifre zorunludur"}).
//         AddRule(ValidationRule{Type: "min", Value: 8, Message: "En az 8 karakter"})
//
//     // Güvenlik seviyesine göre ek kurallar
//     securityLevel := formData["security_level"].(string)
//     if securityLevel == "high" {
//         update.
//             AddRule(ValidationRule{Type: "uppercase", Message: "En az 1 büyük harf"}).
//             AddRule(ValidationRule{Type: "lowercase", Message: "En az 1 küçük harf"}).
//             AddRule(ValidationRule{Type: "number", Message: "En az 1 rakam"}).
//             AddRule(ValidationRule{Type: "special", Message: "En az 1 özel karakter"})
//     } else if securityLevel == "medium" {
//         update.
//             AddRule(ValidationRule{Type: "uppercase", Message: "En az 1 büyük harf"}).
//             AddRule(ValidationRule{Type: "number", Message: "En az 1 rakam"})
//     }
//
//     return update
// })
//
// // Kademeli kural ekleme
// emailField.OnChange(func(field *Schema, formData map[string]interface{}, ctx *fiber.Ctx) *FieldUpdate {
//     update := NewFieldUpdate().
//         AddRule(ValidationRule{Type: "required", Message: "E-posta zorunludur"}).
//         AddRule(ValidationRule{Type: "email", Message: "Geçerli bir e-posta girin"})
//
//     // Kurumsal hesap kontrolü
//     accountType := formData["account_type"].(string)
//     if accountType == "corporate" {
//         update.AddRule(ValidationRule{
//             Type:    "domain_whitelist",
//             Value:   []string{"company.com", "partner.com"},
//             Message: "Sadece şirket e-postaları kabul edilir",
//         })
//     }
//
//     // Benzersizlik kontrolü
//     checkUnique := formData["check_unique"].(bool)
//     if checkUnique {
//         update.AddRule(ValidationRule{
//             Type:    "unique",
//             Table:   "users",
//             Column:  "email",
//             Message: "Bu e-posta zaten kullanılıyor",
//         })
//     }
//
//     return update
// })
//
// // Çoklu koşullu kural ekleme
// priceField.OnChange(func(field *Schema, formData map[string]interface{}, ctx *fiber.Ctx) *FieldUpdate {
//     update := NewFieldUpdate().
//         AddRule(ValidationRule{Type: "required", Message: "Fiyat zorunludur"}).
//         AddRule(ValidationRule{Type: "numeric", Message: "Geçerli bir sayı girin"}).
//         AddRule(ValidationRule{Type: "min", Value: 0, Message: "Fiyat negatif olamaz"})
//
//     // Ürün tipine göre
//     productType := formData["product_type"].(string)
//     if productType == "premium" {
//         update.AddRule(ValidationRule{
//             Type:    "min",
//             Value:   1000,
//             Message: "Premium ürünler minimum 1000 TL olmalıdır",
//         })
//     }
//
//     // Kullanıcı rolüne göre
//     userRole := ctx.Locals("user_role").(string)
//     if userRole != "admin" {
//         update.AddRule(ValidationRule{
//             Type:    "max",
//             Value:   10000,
//             Message: "Maksimum 10.000 TL girebilirsiniz",
//         })
//     }
//
//     // İndirim varsa
//     hasDiscount := formData["has_discount"].(bool)
//     if hasDiscount {
//         update.AddRule(ValidationRule{
//             Type:    "custom",
//             Value:   "validate_discount_price",
//             Message: "İndirimli fiyat normal fiyattan düşük olmalıdır",
//         })
//     }
//
//     return update
// })
//
// // Method chaining ile zincirleme ekleme
// update := NewFieldUpdate().
//     MakeRequired().
//     AddRule(ValidationRule{Type: "min", Value: 3, Message: "En az 3 karakter"}).
//     AddRule(ValidationRule{Type: "max", Value: 50, Message: "En fazla 50 karakter"}).
//     AddRule(ValidationRule{Type: "alphanumeric", Message: "Sadece harf ve rakam"}).
//     SetHelpText("Kullanıcı adı 3-50 karakter, harf ve rakam içermelidir")
//
// // Dinamik kural biriktirme
// ageField.OnChange(func(field *Schema, formData map[string]interface{}, ctx *fiber.Ctx) *FieldUpdate {
//     update := NewFieldUpdate().
//         AddRule(ValidationRule{Type: "required", Message: "Yaş zorunludur"}).
//         AddRule(ValidationRule{Type: "numeric", Message: "Geçerli bir sayı girin"})
//
//     // Yaş aralığı kontrolü
//     serviceType := formData["service_type"].(string)
//     switch serviceType {
//     case "child":
//         update.
//             AddRule(ValidationRule{Type: "min", Value: 0, Message: "Yaş 0'dan küçük olamaz"}).
//             AddRule(ValidationRule{Type: "max", Value: 12, Message: "Çocuk servisi 0-12 yaş arası"})
//     case "teen":
//         update.
//             AddRule(ValidationRule{Type: "min", Value: 13, Message: "Genç servisi 13 yaşından başlar"}).
//             AddRule(ValidationRule{Type: "max", Value: 17, Message: "Genç servisi 13-17 yaş arası"})
//     case "adult":
//         update.
//             AddRule(ValidationRule{Type: "min", Value: 18, Message: "Yetişkin servisi 18 yaşından başlar"}).
//             AddRule(ValidationRule{Type: "max", Value: 65, Message: "Yetişkin servisi 18-65 yaş arası"})
//     case "senior":
//         update.
//             AddRule(ValidationRule{Type: "min", Value: 65, Message: "Yaşlı servisi 65 yaşından başlar"}).
//             AddRule(ValidationRule{Type: "max", Value: 120, Message: "Geçerli bir yaş girin"})
//     }
//
//     return update
// })
// ```
//
// # Önemli Notlar
//
// - **Mevcut Kuralları Korur**: AddRule() mevcut kuralları korur ve sonuna ekler
// - **Sıralama**: Kurallar eklenme sırasına göre uygulanır
// - **Çoklu Çağrı**: Birden fazla kez çağrılabilir (method chaining)
// - **Performans**: Her çağrı slice'a append yapar
// - **Nil Kontrol**: Rules nil ise otomatik olarak initialize edilir
//
// # Best Practices
//
// - Temel kuralları önce ekleyin (required, type checks)
// - Detaylı kuralları sonra ekleyin (min, max, pattern)
// - Pahalı kuralları en sona ekleyin (unique, custom validators)
// - Method chaining kullanarak okunabilir kod yazın
// - Koşullu ekleme için if blokları kullanın
// - Kullanıcı dostu hata mesajları yazın
//
// # Avantajlar
//
// - Mevcut kuralları korur
// - Esnek ve modüler kural ekleme
// - Method chaining desteği
// - Koşullu kural ekleme için ideal
// - Okunabilir kod yazımı
//
// # Dezavantajlar
//
// - Çok fazla çağrı performansı etkileyebilir
// - Kural sırası önemlidir (dikkatli olunmalı)
// - Duplicate kural kontrolü yoktur
//
// # AddRule vs SetRules Karşılaştırması
//
// | Özellik | AddRule | SetRules |
// |---------|---------|----------|
// | Mevcut Kurallar | Korur ve ekler | Siler ve değiştirir |
// | Kullanım | Tekli ekleme | Toplu güncelleme |
// | Method Chaining | Evet | Evet |
// | Esneklik | Yüksek | Düşük |
// | Performans | Her çağrıda append | Tek seferde |
// | Kullanım Amacı | Kademeli ekleme | Tüm kuralları değiştir |
//
// # Performans İpuçları
//
// - Çok fazla AddRule çağrısı yerine SetRules kullanmayı düşünün
// - Pahalı validasyonları (DB sorguları) en sona ekleyin
// - Gereksiz kural eklemekten kaçının
// - Kural sayısını makul seviyede tutun (ideal: 3-7 kural)
func (u *FieldUpdate) AddRule(rule ValidationRule) *FieldUpdate {
	u.Rules = append(u.Rules, rule)
	return u
}
