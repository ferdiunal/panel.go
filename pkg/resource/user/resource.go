package user

import (
	"github.com/ferdiunal/panel.go/pkg/data"
	domainUser "github.com/ferdiunal/panel.go/pkg/domain/user"
	"github.com/ferdiunal/panel.go/pkg/resource"
	"gorm.io/gorm"
)

// init, user resource'unu global registry'ye register eder.
//
// Bu fonksiyon, package import edildiğinde otomatik olarak çalışır ve
// user resource'unu "users" slug'ı ile registry'ye ekler.
//
// # Kullanım
//
// Bu fonksiyon otomatik olarak çalışır, manuel çağrı gerekmez:
//
//	import _ "github.com/ferdiunal/panel.go/pkg/resource/user"
//
// # Önemli Notlar
//
// - init() fonksiyonu package import edildiğinde otomatik çalışır
// - Resource registry'ye "users" slug'ı ile eklenir
// - Circular dependency sorununu önlemek için kullanılır
func init() {
	resource.Register("users", NewUserResource())
}

// Bu yapı, panel yönetim sisteminde kullanıcı kaynağını (resource) temsil eder.
//
// UserResource, CRUD işlemleri, yetkilendirme, alan çözümleme ve veri sağlama
// gibi kullanıcı yönetimi ile ilgili tüm işlevleri kapsüller. OptimizedBase
// yapısından kalıtım alarak, panel sisteminin temel kaynak özelliklerini
// otomatik olarak miras alır.
//
// # Kullanım Senaryoları
//
// - Kullanıcı listesi görüntüleme ve yönetimi
// - Yeni kullanıcı oluşturma ve düzenleme
// - Kullanıcı silme ve geri yükleme
// - Kullanıcı izinleri ve rolleri yönetimi
// - Kullanıcı profil bilgilerinin görüntülenmesi
//
// # Önemli Notlar
//
// - UserResource, OptimizedBase'den kalıtım alarak performans optimizasyonlarından faydalanır
// - Tüm CRUD işlemleri UserPolicy tarafından yetkilendirilir
// - Alan çözümleme UserFieldResolver tarafından yönetilir
// - Kart görünümü UserCardResolver tarafından işlenir
//
// # Örnek Kullanım
//
//	userResource := NewUserResource()
//	// userResource artık panel sistemine kaydedilebilir ve kullanılabilir
type UserResource struct {
	// OptimizedBase, panel sisteminin temel kaynak özelliklerini sağlar.
	// Bu yapı, model tanımı, slug, başlık, ikon, grup, görünürlük,
	// navigasyon sırası, yetkilendirme politikası, alan çözümleyici,
	// kart çözümleyici ve veri sağlayıcı gibi özellikleri içerir.
	resource.OptimizedBase
}

// Bu fonksiyon, yeni bir UserResource örneği oluşturur ve tüm gerekli
// konfigürasyonları uygulayarak hazır hale getirir.
//
// # Fonksiyon Açıklaması
//
// NewUserResource, UserResource yapısının factory fonksiyonudur. Yeni bir
// UserResource örneği oluşturur, tüm temel ayarları yapılandırır ve panel
// sistemine entegre olmak için gerekli tüm bileşenleri (policy, resolver'lar)
// atanır.
//
// # Dönüş Değeri
//
// Döndürür: - Tam olarak yapılandırılmış UserResource pointer'ı
//
// # Yapılandırılan Özellikler
//
// - Model: domainUser.User (kullanıcı domain modeli)
// - Slug: "users" (URL'de kullanılan benzersiz tanımlayıcı)
// - Başlık: "Users" (UI'da gösterilen başlık - i18n destekli)
// - İkon: "users" (UI'da gösterilen ikon)
// - Grup: "System" (panel menüsünde gösterileceği grup - i18n destekli)
// - Görünürlük: true (panel menüsünde görünür)
// - Navigasyon Sırası: 1 (menüde gösterilme sırası)
// - Yetkilendirme Politikası: UserPolicy (erişim kontrolü)
// - Alan Çözümleyici: UserFieldResolver (alan işleme)
// - Kart Çözümleyici: UserCardResolver (kart görünümü işleme)
//
// # Kullanım Örneği
//
//	userResource := NewUserResource()
//	// userResource artık panel sistemine kaydedilebilir
//	panel.RegisterResource(userResource)
//
// # Önemli Notlar
//
// - Bu fonksiyon her çağrıldığında yeni bir UserResource örneği oluşturur
// - Tüm ayarlar sırasıyla uygulanır, bu nedenle sıra önemlidir
// - UserPolicy, UserFieldResolver ve UserCardResolver paketinde tanımlanmış olmalıdır
// - Oluşturulan kaynak, panel sistemine kaydedilmeden önce ek konfigürasyonlar yapılabilir
// - Başlık ve grup i18n desteği ile çoklu dilde gösterilebilir
func NewUserResource() *UserResource {
	r := &UserResource{}

	// Temel kaynak ayarlarını yapılandır
	// Model, slug, başlık, ikon ve grup gibi temel özellikleri tanımla
	r.SetModel(&domainUser.User{})
	r.SetSlug("users")
	r.SetTitle("Users")
	r.SetIcon("users")
	r.SetGroup("System")
	r.SetVisible(true)
	r.SetNavigationOrder(1)

	// Kayıt başlığı için "name" field'ını kullan
	// İlişki fieldlarında kullanıcılar "John Doe" gibi okunabilir şekilde gösterilir
	r.SetRecordTitleKey("name")

	// Yetkilendirme politikasını ayarla
	// UserPolicy, bu kaynağa erişim kontrolü sağlar
	r.SetPolicy(&UserPolicy{})

	// Alan ve kart çözümleyicilerini ayarla
	// UserFieldResolver, alan işleme ve görüntüleme mantığını yönetir
	// UserCardResolver, kart görünümü işleme mantığını yönetir
	r.SetFieldResolver(&UserFieldResolver{})
	r.SetCardResolver(&UserCardResolver{})

	return r
}

// Bu metod, UserResource için özel bir veri sağlayıcı (data provider) oluşturur
// ve döndürür.
//
// # Metod Açıklaması
//
// Repository metodu, verilen GORM veritabanı bağlantısını kullanarak
// UserResource için özel bir veri sağlayıcı oluşturur. Bu veri sağlayıcı,
// kullanıcı verilerine erişim, sorgu, filtreleme ve manipülasyon işlemlerini
// yönetir.
//
// # Parametreler
//
//   - db (*gorm.DB): GORM veritabanı bağlantısı. Bu bağlantı, veri sağlayıcı
//     tarafından tüm veritabanı işlemleri için kullanılır.
//
// # Dönüş Değeri
//
// Döndürür: - Yapılandırılmış UserDataProvider pointer'ı (data.DataProvider interface'ini uygular)
//
// # Kullanım Senaryoları
//
// - Kullanıcı verilerini veritabanından sorgulamak
// - Yeni kullanıcı kayıtları oluşturmak
// - Mevcut kullanıcı bilgilerini güncellemek
// - Kullanıcı kayıtlarını silmek
// - Kullanıcı verilerini filtrelemek ve sıralamak
//
// # Kullanım Örneği
//
//	userResource := NewUserResource()
//	db := gorm.Open(sqlite.Open("test.db"))
//	dataProvider := userResource.Repository(db)
//	// dataProvider artık kullanıcı verilerine erişmek için kullanılabilir
//
// # Önemli Notlar
//
// - Bu metod, her çağrıldığında yeni bir UserDataProvider örneği oluşturur
// - Verilen db bağlantısı, veri sağlayıcı tarafından tüm işlemler boyunca kullanılır
// - UserDataProvider, data.DataProvider interface'ini uygulamalıdır
// - Veritabanı bağlantısı geçerli ve açık olmalıdır
// - Bağlantı kapatıldıktan sonra veri sağlayıcı kullanılamaz
// - Hata yönetimi, veri sağlayıcı tarafından yapılır
func (r *UserResource) Repository(client interface{}) data.DataProvider {
	// Type assertion to get Ent client
	db, ok := client.(*gorm.DB)
	if !ok {
		// TODO: Add GORM support
		return nil
	}

	return NewUserDataProvider(db)
}
