// Bu paket, kullanıcı kaynağı (User Resource) yönetimini sağlar.
// Kullanıcı verilerinin CRUD işlemleri, doğrulama ve veri sağlayıcı konfigürasyonunu içerir.
package user

import (
	"github.com/ferdiunal/panel.go/pkg/data"
	"github.com/ferdiunal/panel.go/pkg/domain/user"
	"github.com/ferdiunal/panel.go/pkg/resource"
	"gorm.io/gorm"
)

// Bu yapı, kullanıcı kaynağı (Resource) için özel Repository davranışını sağlar.
//
// UserResourceWrapper, resource.Base yapısını gömülerek (embed) GenericResource'ın
// varsayılan Repository metodunu geçersiz kılmak (override) için kullanılır.
// Bu sayede kullanıcı verilerine özel bir veri sağlayıcı (DataProvider) tanımlanabilir.
//
// **Kullanım Senaryosu:**
// - Kullanıcı verilerine özel sorgu ve filtreleme mantığı uygulamak
// - Kullanıcı tablosuna özel ilişkiler (relationships) tanımlamak
// - Kullanıcı verilerinin yüklenmesi sırasında özel işlemler yapmak
//
// **Önemli Not:**
// Bu yapı artık kullanılmamaktadır (deprecated). Yeni kodlarda NewUserResource()
// fonksiyonunu kullanınız. Bu yapı sadece geriye uyumluluk (backward compatibility)
// için korunmaktadır.
//
// **Örnek:**
//
//	wrapper := UserResourceWrapper{}
//	db := gorm.Open(...)
//	provider := wrapper.Repository(db)  // Kullanıcı veri sağlayıcısını alır
//
// Deprecated: NewUserResource() fonksiyonunu kullanınız.
type UserResourceWrapper struct {
	// Base, resource.Base yapısını gömülerek temel kaynak işlevselliğini sağlar.
	// Bu sayede UserResourceWrapper, resource.Resource arayüzünü otomatik olarak
	// uygular ve tüm temel metotları miras alır.
	resource.Base
}

// Bu metod, kullanıcı verilerine özel bir veri sağlayıcı (DataProvider) oluşturur.
//
// Repository metodu, GORM veritabanı bağlantısını alarak, kullanıcı tablosuna
// özel bir veri sağlayıcı (UserRepository) döndürür. Bu veri sağlayıcı, kullanıcı
// verilerinin tüm CRUD işlemlerini (Create, Read, Update, Delete) yönetir.
//
// **Parametreler:**
//   - db (*gorm.DB): GORM veritabanı bağlantısı. Bu bağlantı, kullanıcı verilerine
//     erişmek için kullanılır.
//
// **Dönüş Değeri:**
//   - data.DataProvider: Kullanıcı verilerine erişim sağlayan veri sağlayıcı arayüzü.
//     Bu arayüz, Create, Read, Update, Delete ve List gibi temel veri işlemlerini
//     tanımlar.
//
// **Kullanım Senaryosu:**
// - Veritabanı bağlantısı kurulduktan sonra kullanıcı verilerine erişmek
// - Kullanıcı verilerini sorgulamak, güncellemek veya silmek
// - Kullanıcı tablosuna özel filtreleme ve sıralama uygulamak
//
// **Örnek:**
//
//	db := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
//	wrapper := UserResourceWrapper{}
//	provider := wrapper.Repository(db)
//	// provider artık kullanıcı verilerine erişebilir
//
// **Önemli Not:**
// Döndürülen veri sağlayıcı, orm.NewUserRepository() tarafından oluşturulur.
// Bu fonksiyon, UserRepository yapısını başlatır ve GORM bağlantısını ayarlar.
//
// Döndürür: - Yapılandırılmış UserRepository veri sağlayıcısı pointer'ı
func (r UserResourceWrapper) Repository(db *gorm.DB) data.DataProvider {
	return data.NewGormDataProvider(db, &user.User{})
}

// Bu fonksiyon, kullanıcı kaynağının (Resource) tam konfigürasyonunu döner.
//
// GetUserResource, kullanıcı yönetimi için gerekli tüm ayarları içeren bir
// resource.Resource nesnesi oluşturur ve döndürür. Bu nesne, kullanıcı verilerinin
// yönetimi, doğrulama, filtreleme ve API uç noktaları (endpoints) için kullanılır.
//
// **Kullanım Senaryosu:**
// - Panel uygulamasında kullanıcı yönetim sayfasını oluşturmak
// - Kullanıcı verilerine erişim sağlayan API uç noktalarını tanımlamak
// - Kullanıcı verilerinin doğrulama kurallarını ayarlamak
// - Kullanıcı tablosunun sütunlarını ve ilişkilerini tanımlamak
//
// **Dönüş Değeri:**
//   - resource.Resource: Kullanıcı kaynağının tam konfigürasyonunu içeren nesne.
//     Bu nesne, aşağıdaki bilgileri içerir:
//   - Kaynak adı ve açıklaması
//   - Veritabanı tablosu bilgileri
//   - Sütun tanımları ve doğrulama kuralları
//   - İlişkiler (relationships)
//   - Filtreleme ve arama seçenekleri
//   - Özel eylemler (actions)
//
// **Örnek:**
//
//	userResource := GetUserResource()
//	// userResource artık panel uygulamasında kullanılabilir
//	// Örneğin: panel.RegisterResource(userResource)
//
// **Önemli Not:**
// Bu fonksiyon, NewUserResource() fonksiyonunu çağırarak yeni bir kullanıcı
// kaynağı oluşturur. GetUserResource() artık kullanılmamaktadır (deprecated).
// Yeni kodlarda doğrudan NewUserResource() fonksiyonunu kullanınız.
//
// **Geriye Uyumluluk:**
// Bu fonksiyon, eski kodların çalışmaya devam etmesi için korunmaktadır.
// Ancak yeni projeler için NewUserResource() kullanılması önerilir.
//
// Deprecated: NewUserResource() fonksiyonunu kullanınız.
func GetUserResource() resource.Resource {
	return NewUserResource()
}
