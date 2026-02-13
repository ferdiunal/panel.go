// Bu yapı, panel uygulamasında aksiyon (action) yürütülmesi sırasında gereken tüm bağlamsal
// bilgileri içerir. Aksiyon, bir kaynağa (resource) uygulanabilen özel bir işlemdir.
// Örneğin: "Yayınla", "Sil", "Gönder" gibi işlemler.
//
// # Kullanım Senaryoları
// - Toplu işlemler: Seçilen birden fazla kaynağa aynı işlemi uygulamak
// - Veri doğrulama: Kullanıcı tarafından gönderilen form verilerini işlemek
// - Yetkilendirme: Mevcut kullanıcının işlemi yapma yetkisini kontrol etmek
// - Veritabanı işlemleri: Seçilen kaynakları güncellemek veya silmek
//
// # Örnek Kullanım
// ```go
// ctx := &ActionContext{
//     Models: []interface{}{post1, post2, post3},
//     Fields: map[string]interface{}{
//         "status": "published",
//         "published_at": time.Now(),
//     },
//     User: currentUser,
//     Resource: "posts",
//     DB: db,
//     Ctx: fiberCtx,
// }
// ```
package action

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// Bu yapı, aksiyon yürütülmesi sırasında gereken tüm bağlamsal bilgileri tutar.
// Seçilen kaynaklar, alan değerleri, kullanıcı bilgileri ve veritabanı bağlantısını içerir.
//
// # Alanlar
// - Models: İşlem yapılacak seçilen kaynak örnekleri (örn: post, user vb.)
// - Fields: Kullanıcı tarafından gönderilen form alanlarının değerleri
// - User: İşlemi gerçekleştiren kimlik doğrulanmış kullanıcı
// - Resource: Kaynak türünün slug'ı (örn: "posts", "users", "comments")
// - DB: GORM veritabanı bağlantısı, veritabanı işlemleri için kullanılır
// - Ctx: Fiber HTTP bağlamı, HTTP isteği/yanıtı bilgilerine erişim sağlar
//
// # Önemli Notlar
// - Models slice'ı boş olabilir (hiçbir kaynak seçilmemişse)
// - Fields map'i nil olabilir (form verisi gönderilmemişse)
// - User ve Resource her zaman doldurulmalıdır
// - DB ve Ctx nil olmamalıdır, aksi takdirde işlem başarısız olur
type ActionContext struct {
	// Bu alan, işlem yapılacak seçilen kaynak örneklerini içerir.
	// Örneğin, "Yayınla" aksiyonu için seçilen tüm yazılar bu slice'da bulunur.
	//
	// # Kullanım
	// - Toplu işlemler için birden fazla kaynak içerebilir
	// - Boş olabilir (hiçbir kaynak seçilmemişse)
	// - Her eleman, ilgili kaynağın veritabanı modeli örneğidir
	//
	// # Örnek
	// ```go
	// for _, model := range ctx.Models {
	//     post := model.(*Post)
	//     // post üzerinde işlem yap
	// }
	// ```
	Models []interface{}

	// Bu alan, kullanıcı tarafından aksiyon formu aracılığıyla gönderilen
	// tüm form alanlarının değerlerini içerir.
	//
	// # Kullanım
	// - Aksiyon parametrelerini almak için kullanılır
	// - Örneğin: "Yayınla" aksiyonunda "yayın tarihi" bilgisi
	// - Doğrulama ve işleme tabi tutulmalıdır
	//
	// # Örnek
	// ```go
	// if status, ok := ctx.Fields["status"]; ok {
	//     // status değerini kullan
	// }
	// ```
	//
	// # Önemli Notlar
	// - Nil olabilir (form verisi gönderilmemişse)
	// - Tüm değerler interface{} türündedir, type assertion gereklidir
	// - Güvenlik açısından, bu veriler doğrulanmalı ve sanitize edilmelidir
	Fields map[string]interface{}

	// Bu alan, aksiyonu gerçekleştiren kimlik doğrulanmış kullanıcıyı temsil eder.
	// Yetkilendirme ve denetim (audit) işlemleri için kullanılır.
	//
	// # Kullanım
	// - Kullanıcının aksiyonu yapma yetkisini kontrol etmek
	// - İşlem günlüğüne (audit log) kullanıcı bilgisini kaydetmek
	// - Kullanıcıya özel işlemler yapmak
	//
	// # Örnek
	// ```go
	// user := ctx.User.(*User)
	// if !user.HasPermission("publish_posts") {
	//     return errors.New("yetkiniz yok")
	// }
	// ```
	//
	// # Önemli Notlar
	// - Nil olmamalıdır (kimlik doğrulama gereklidir)
	// - interface{} türündedir, type assertion ile kullanılmalıdır
	// - Genellikle *User türüne dönüştürülür
	User interface{}

	// Bu alan, işlem yapılacak kaynağın türünü belirten slug'ı içerir.
	// Örneğin: "posts", "users", "comments", "products" vb.
	//
	// # Kullanım
	// - Hangi kaynak türü üzerinde işlem yapıldığını belirlemek
	// - Yetkilendirme kurallarını uygulamak
	// - İşlem günlüğüne kaydetmek
	// - Doğru modeli yüklemek
	//
	// # Örnek
	// ```go
	// switch ctx.Resource {
	// case "posts":
	//     // yazılar için işlem
	// case "users":
	//     // kullanıcılar için işlem
	// }
	// ```
	//
	// # Önemli Notlar
	// - Boş olmamalıdır
	// - Küçük harfle yazılmalıdır (convention)
	// - Veritabanı tablosu adıyla eşleşmelidir
	Resource string

	// Bu alan, veritabanı işlemleri için GORM bağlantısını içerir.
	// Seçilen kaynakları güncellemek, silmek veya sorgulamak için kullanılır.
	//
	// # Kullanım
	// - Seçilen kaynakları veritabanından güncellemek
	// - Yeni kayıtlar oluşturmak
	// - İlişkili verileri yüklemek
	// - İşlem (transaction) başlatmak
	//
	// # Örnek
	// ```go
	// for _, model := range ctx.Models {
	//     if err := ctx.DB.Model(model).Update("status", "published").Error; err != nil {
	//         return err
	//     }
	// }
	// ```
	//
	// # Önemli Notlar
	// - Nil olmamalıdır
	// - Genellikle middleware tarafından doldurulur
	// - İşlem (transaction) içinde kullanılabilir
	// - Hata kontrolü yapılmalıdır
	DB *gorm.DB

	// Bu alan, HTTP isteği/yanıtı bilgilerine erişim sağlayan Fiber bağlamını içerir.
	// HTTP başlıkları, query parametreleri, body vb. bilgilere erişmek için kullanılır.
	//
	// # Kullanım
	// - HTTP başlıklarını okumak
	// - Query parametrelerini almak
	// - Yanıt göndermek
	// - Session/Cookie bilgilerine erişmek
	// - İstek IP adresini almak
	//
	// # Örnek
	// ```go
	// userAgent := ctx.Ctx.Get("User-Agent")
	// clientIP := ctx.Ctx.IP()
	// ctx.Ctx.JSON(200, result)
	// ```
	//
	// # Önemli Notlar
	// - Nil olmamalıdır
	// - Fiber framework tarafından sağlanır
	// - HTTP yanıtı göndermek için kullanılabilir
	// - İstek tamamlandıktan sonra kullanılamaz
	Ctx *fiber.Ctx
}
