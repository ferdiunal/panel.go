package core

// # Resolver Interface
//
// Bu interface, alan değerlerini dinamik olarak çözümlemek (resolve) için kullanılır.
// Alanların, kayıt verilerine ve parametrelere göre değerlerini hesaplamasına veya
// dönüştürmesine olanak tanır.
//
// ## Genel Bakış
//
// Resolver interface'i, field sisteminde değer çözümleme işlemlerini standartlaştırır.
// Bu sayede alanlar, veritabanından gelen ham verileri kullanıcı dostu formatlara
// dönüştürebilir, hesaplanmış değerler üretebilir veya ilişkili verileri yükleyebilir.
//
// ## Kullanım Senaryoları
//
// ### 1. Hesaplanmış Alanlar (Computed Fields)
//
// Birden fazla alandan yeni bir değer hesaplamak için kullanılır:
//
// ```go
// type FullNameResolver struct{}
//
// func (r *FullNameResolver) Resolve(item interface{}, params map[string]interface{}) (interface{}, error) {
//     if user, ok := item.(*User); ok {
//         return user.FirstName + " " + user.LastName, nil
//     }
//     return nil, fmt.Errorf("invalid item type")
// }
//
// // Field tanımında kullanım
// fields.Text("Tam Ad", "full_name").
//     OnList().
//     WithResolver(&FullNameResolver{})
// ```
//
// ### 2. Değer Formatlama
//
// Ham verileri kullanıcı dostu formatlara dönüştürmek için kullanılır:
//
// ```go
// type PriceResolver struct{}
//
// func (r *PriceResolver) Resolve(item interface{}, params map[string]interface{}) (interface{}, error) {
//     if product, ok := item.(*Product); ok {
//         return fmt.Sprintf("₺%.2f", product.Price), nil
//     }
//     return nil, fmt.Errorf("invalid item type")
// }
//
// // Field tanımında kullanım
// fields.Text("Fiyat", "price").
//     OnList().
//     WithResolver(&PriceResolver{})
// ```
//
// ### 3. İlişkili Veri Yükleme
//
// İlişkili kayıtlardan veri çekmek için kullanılır:
//
// ```go
// type AuthorNameResolver struct {
//     DB *gorm.DB
// }
//
// func (r *AuthorNameResolver) Resolve(item interface{}, params map[string]interface{}) (interface{}, error) {
//     if post, ok := item.(*Post); ok {
//         var author User
//         if err := r.DB.First(&author, post.AuthorID).Error; err != nil {
//             return nil, err
//         }
//         return author.Name, nil
//     }
//     return nil, fmt.Errorf("invalid item type")
// }
// ```
//
// ### 4. Koşullu Değer Dönüşümü
//
// Parametrelere veya kayıt durumuna göre farklı değerler döndürmek için kullanılır:
//
// ```go
// type StatusResolver struct{}
//
// func (r *StatusResolver) Resolve(item interface{}, params map[string]interface{}) (interface{}, error) {
//     if order, ok := item.(*Order); ok {
//         locale := params["locale"].(string)
//
//         statusMap := map[string]map[string]string{
//             "pending": {
//                 "tr": "Beklemede",
//                 "en": "Pending",
//             },
//             "completed": {
//                 "tr": "Tamamlandı",
//                 "en": "Completed",
//             },
//         }
//
//         if translations, ok := statusMap[order.Status]; ok {
//             if translation, ok := translations[locale]; ok {
//                 return translation, nil
//             }
//         }
//         return order.Status, nil
//     }
//     return nil, fmt.Errorf("invalid item type")
// }
// ```
//
// ### 5. Veri Maskeleme ve Güvenlik
//
// Hassas verileri maskelemek veya kullanıcı yetkilerine göre filtrelemek için kullanılır:
//
// ```go
// type EmailResolver struct{}
//
// func (r *EmailResolver) Resolve(item interface{}, params map[string]interface{}) (interface{}, error) {
//     if user, ok := item.(*User); ok {
//         // Kullanıcı yetkisini kontrol et
//         if canViewFullEmail, ok := params["can_view_email"].(bool); ok && canViewFullEmail {
//             return user.Email, nil
//         }
//
//         // E-postayı maskele
//         parts := strings.Split(user.Email, "@")
//         if len(parts) == 2 {
//             masked := parts[0][:2] + "***@" + parts[1]
//             return masked, nil
//         }
//         return "***", nil
//     }
//     return nil, fmt.Errorf("invalid item type")
// }
// ```
//
// ## Avantajlar
//
// - **Esneklik**: Değer çözümleme mantığını field tanımından ayırır
// - **Yeniden Kullanılabilirlik**: Aynı resolver'ı birden fazla field'da kullanabilirsiniz
// - **Test Edilebilirlik**: Resolver'ları bağımsız olarak test edebilirsiniz
// - **Bakım Kolaylığı**: Değer dönüşüm mantığı merkezi bir yerde toplanır
// - **Tip Güvenliği**: Interface sayesinde tüm resolver'lar aynı imzayı kullanır
//
// ## Dezavantajlar ve Dikkat Edilmesi Gerekenler
//
// - **Performans**: Her kayıt için resolver çağrılır, bu nedenle ağır işlemlerden kaçının
// - **N+1 Problem**: İlişkili veri yüklerken eager loading kullanmayı unutmayın
// - **Hata Yönetimi**: Resolver hataları kullanıcıya anlamlı mesajlar olarak iletilmelidir
// - **Null Güvenliği**: Item ve params değerlerinin nil olabileceğini unutmayın
//
// ## Best Practices
//
// ### 1. Hata Yönetimi
//
// ```go
// func (r *MyResolver) Resolve(item interface{}, params map[string]interface{}) (interface{}, error) {
//     // Tip kontrolü yap
//     user, ok := item.(*User)
//     if !ok {
//         return nil, fmt.Errorf("expected *User, got %T", item)
//     }
//
//     // Nil kontrolü yap
//     if user == nil {
//         return nil, fmt.Errorf("user is nil")
//     }
//
//     // İşlemi gerçekleştir
//     return user.FullName(), nil
// }
// ```
//
// ### 2. Performans Optimizasyonu
//
// ```go
// // Kötü: Her kayıt için veritabanı sorgusu
// func (r *BadResolver) Resolve(item interface{}, params map[string]interface{}) (interface{}, error) {
//     post := item.(*Post)
//     var author User
//     db.First(&author, post.AuthorID) // N+1 problem!
//     return author.Name, nil
// }
//
// // İyi: Eager loading kullan veya cache'den al
// func (r *GoodResolver) Resolve(item interface{}, params map[string]interface{}) (interface{}, error) {
//     post := item.(*Post)
//     if post.Author != nil { // Eager loaded
//         return post.Author.Name, nil
//     }
//     return "Unknown", nil
// }
// ```
//
// ### 3. Parametre Kullanımı
//
// ```go
// func (r *LocalizedResolver) Resolve(item interface{}, params map[string]interface{}) (interface{}, error) {
//     product := item.(*Product)
//
//     // Parametreyi güvenli şekilde al
//     locale, ok := params["locale"].(string)
//     if !ok {
//         locale = "en" // Varsayılan değer
//     }
//
//     // Locale'e göre değer döndür
//     if locale == "tr" {
//         return product.NameTR, nil
//     }
//     return product.NameEN, nil
// }
// ```
//
// ### 4. Bağımlılık Enjeksiyonu
//
// ```go
// type DatabaseResolver struct {
//     DB     *gorm.DB
//     Cache  cache.Cache
//     Logger logger.Logger
// }
//
// func (r *DatabaseResolver) Resolve(item interface{}, params map[string]interface{}) (interface{}, error) {
//     // Bağımlılıkları kullan
//     r.Logger.Debug("Resolving value")
//
//     // Cache'den kontrol et
//     if cached := r.Cache.Get("key"); cached != nil {
//         return cached, nil
//     }
//
//     // Veritabanından al
//     var result string
//     r.DB.Raw("SELECT value FROM table").Scan(&result)
//
//     // Cache'e kaydet
//     r.Cache.Set("key", result)
//
//     return result, nil
// }
// ```
//
// ## Field System ile Entegrasyon
//
// Resolver interface'i, field sisteminde `Resolve` callback'i ile kullanılır.
// Detaylı bilgi için [Fields Dokümantasyonu](../../docs/Fields.md#callback-ler) bölümüne bakınız.
//
// ```go
// // Inline resolver kullanımı
// fields.Text("Tam Ad", "full_name").
//     OnList().
//     Resolve(func(value interface{}, item interface{}, c *fiber.Ctx) interface{} {
//         if user, ok := item.(*User); ok {
//             return user.FirstName + " " + user.LastName
//         }
//         return value
//     })
//
// // Struct-based resolver kullanımı
// type FullNameResolver struct{}
//
// func (r *FullNameResolver) Resolve(item interface{}, params map[string]interface{}) (interface{}, error) {
//     if user, ok := item.(*User); ok {
//         return user.FirstName + " " + user.LastName, nil
//     }
//     return nil, fmt.Errorf("invalid item type")
// }
//
// fields.Text("Tam Ad", "full_name").
//     OnList().
//     WithResolver(&FullNameResolver{})
// ```
//
// ## İlişkiler ile Kullanım
//
// Resolver'lar, ilişkili verileri yüklemek ve formatlamak için de kullanılabilir.
// Detaylı bilgi için [Relationships Dokümantasyonu](../../docs/Relationships.md) bölümüne bakınız.
//
// ```go
// // BelongsTo ilişkisinde resolver kullanımı
// fields.BelongsTo("Yazar", "author_id", "authors").
//     DisplayUsing("name").
//     WithResolver(&AuthorDisplayResolver{})
//
// // HasMany ilişkisinde resolver kullanımı
// fields.HasMany("Yorumlar", "comments", "comments").
//     WithResolver(&CommentCountResolver{})
// ```
//
// ## Test Örneği
//
// ```go
// func TestFullNameResolver(t *testing.T) {
//     resolver := &FullNameResolver{}
//
//     user := &User{
//         FirstName: "John",
//         LastName:  "Doe",
//     }
//
//     result, err := resolver.Resolve(user, nil)
//
//     assert.NoError(t, err)
//     assert.Equal(t, "John Doe", result)
// }
//
// func TestFullNameResolver_InvalidType(t *testing.T) {
//     resolver := &FullNameResolver{}
//
//     result, err := resolver.Resolve("invalid", nil)
//
//     assert.Error(t, err)
//     assert.Nil(t, result)
// }
// ```
//
// ## Önemli Notlar
//
// - Resolver'lar thread-safe olmalıdır (birden fazla goroutine tarafından kullanılabilir)
// - Resolver içinde panic kullanmaktan kaçının, her zaman error döndürün
// - Resolver'lar mümkün olduğunca hafif ve hızlı olmalıdır
// - Ağır işlemler için cache kullanmayı düşünün
// - Resolver'ları test etmeyi unutmayın
//
// ## Uyarılar
//
// - **Sonsuz Döngü**: Resolver içinde başka bir resolver çağırmaktan kaçının
// - **Bellek Sızıntısı**: Resolver içinde büyük nesneleri cache'lemeyin
// - **Veritabanı Bağlantısı**: Her resolver çağrısında yeni bağlantı açmayın
// - **Güvenlik**: Resolver içinde kullanıcı girdilerini doğrulayın
//
// ## İlgili Kaynaklar
//
// - [Fields Dokümantasyonu](../../docs/Fields.md)
// - [Relationships Dokümantasyonu](../../docs/Relationships.md)
// - [Resource Dokümantasyonu](../../docs/Resources.md)
type Resolver interface {
	// # Resolve Metodu
	//
	// Bu metod, bir alan değerini kayıt verilerine ve parametrelere göre hesaplar veya dönüştürür.
	//
	// ## Parametreler
	//
	// - `item interface{}`: Çözümlenecek kayıt verisi. Genellikle bir struct veya map'tir.
	//   Örnek: `*User`, `*Post`, `map[string]interface{}`
	//
	// - `params map[string]interface{}`: Çözümleme için ek parametreler. İsteğe bağlıdır.
	//   Örnek parametreler:
	//   - `"locale"`: Dil kodu (string)
	//   - `"user"`: Mevcut kullanıcı (User)
	//   - `"context"`: HTTP context (fiber.Ctx)
	//   - `"can_view_email"`: Yetki kontrolü (bool)
	//
	// ## Dönüş Değerleri
	//
	// - `interface{}`: Çözümlenmiş değer. Herhangi bir tip olabilir:
	//   - `string`: Metin değerler için
	//   - `int`, `float64`: Sayısal değerler için
	//   - `bool`: Boolean değerler için
	//   - `time.Time`: Tarih/saat değerleri için
	//   - `[]interface{}`: Dizi değerler için
	//   - `map[string]interface{}`: Nesne değerler için
	//
	// - `error`: Çözümleme başarısız olursa hata döndürülür. Hata mesajları kullanıcıya gösterilir.
	//
	// ## Kullanım Örnekleri
	//
	// ### Örnek 1: Basit Değer Dönüşümü
	//
	// ```go
	// func (r *StatusResolver) Resolve(item interface{}, params map[string]interface{}) (interface{}, error) {
	//     order := item.(*Order)
	//
	//     statusMap := map[string]string{
	//         "pending":   "Beklemede",
	//         "completed": "Tamamlandı",
	//         "cancelled": "İptal Edildi",
	//     }
	//
	//     if translated, ok := statusMap[order.Status]; ok {
	//         return translated, nil
	//     }
	//
	//     return order.Status, nil
	// }
	// ```
	//
	// ### Örnek 2: Hesaplanmış Değer
	//
	// ```go
	// func (r *TotalPriceResolver) Resolve(item interface{}, params map[string]interface{}) (interface{}, error) {
	//     order := item.(*Order)
	//
	//     total := 0.0
	//     for _, item := range order.Items {
	//         total += item.Price * float64(item.Quantity)
	//     }
	//
	//     return fmt.Sprintf("₺%.2f", total), nil
	// }
	// ```
	//
	// ### Örnek 3: Parametreli Çözümleme
	//
	// ```go
	// func (r *LocalizedNameResolver) Resolve(item interface{}, params map[string]interface{}) (interface{}, error) {
	//     product := item.(*Product)
	//
	//     locale, ok := params["locale"].(string)
	//     if !ok {
	//         locale = "en"
	//     }
	//
	//     switch locale {
	//     case "tr":
	//         return product.NameTR, nil
	//     case "en":
	//         return product.NameEN, nil
	//     default:
	//         return product.NameEN, nil
	//     }
	// }
	// ```
	//
	// ### Örnek 4: Hata Yönetimi
	//
	// ```go
	// func (r *AgeResolver) Resolve(item interface{}, params map[string]interface{}) (interface{}, error) {
	//     user, ok := item.(*User)
	//     if !ok {
	//         return nil, fmt.Errorf("expected *User, got %T", item)
	//     }
	//
	//     if user.BirthDate == nil {
	//         return nil, fmt.Errorf("birth date is not set")
	//     }
	//
	//     age := time.Now().Year() - user.BirthDate.Year()
	//     return age, nil
	// }
	// ```
	//
	// ## Performans İpuçları
	//
	// - Ağır işlemler için cache kullanın
	// - Veritabanı sorgularından kaçının (eager loading tercih edin)
	// - Tip dönüşümlerini optimize edin
	// - Gereksiz string concatenation'dan kaçının
	//
	// ## Güvenlik İpuçları
	//
	// - Kullanıcı girdilerini doğrulayın
	// - Hassas verileri maskeleyerek döndürün
	// - SQL injection'a karşı korunun
	// - XSS saldırılarına karşı output'u escape edin
	//
	// ## Döndürür
	//
	// - Çözümlenmiş değer (interface{})
	// - Hata durumunda error
	Resolve(item interface{}, params map[string]interface{}) (interface{}, error)
}
