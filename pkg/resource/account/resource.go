// Bu paket, Account (Hesap) entity'si için admin panel resource tanımlarını içerir.
// Account resource'u, kullanıcı hesaplarının yönetimi, görüntülenmesi ve işlenmesi için
// gerekli tüm konfigürasyonları ve metodları sağlar.
package account

import (
	"github.com/ferdiunal/panel.go/pkg/data"
	domainAccount "github.com/ferdiunal/panel.go/pkg/domain/account"
	"github.com/ferdiunal/panel.go/pkg/resource"
	"gorm.io/gorm"
)

// Bu yapı, Account entity'si için admin panel resource tanımını temsil eder.
// AccountResource, OptimizedBase'i embed ederek resource sisteminin tüm özelliklerini
// miras alır ve Account-spesifik konfigürasyonları sağlar.
//
// Kullanım Senaryosu:
// - Admin panelinde Account (Hesap) yönetim sayfasının oluşturulması
// - Hesap listesi, detay, oluşturma ve düzenleme işlemlerinin yönetimi
// - Hesaplarla ilişkili User verilerinin eager loading'i
//
// Önemli Notlar:
// - Bu yapı, resource.OptimizedBase'i embed ederek performans optimizasyonlarından faydalanır
// - Tüm konfigürasyonlar NewAccountResource() fonksiyonunda yapılır
// - Field, Card ve Policy resolver'ları ayrı yapılar tarafından yönetilir
type AccountResource struct {
	resource.OptimizedBase
}

// Bu fonksiyon, yeni bir Account resource'u oluşturur ve tüm gerekli konfigürasyonları ayarlar.
// Fonksiyon, admin panelinde Account yönetim sayfasının temelini oluşturur.
//
// Parametreler: Yok
//
// Dönüş Değeri:
// - *AccountResource: Yapılandırılmış Account resource pointer'ı
//
// Yapılandırılan Özellikler:
// - Model: domainAccount.Account{} - Veritabanı modeli
// - Slug: "accounts" - URL'de kullanılan benzersiz tanımlayıcı
// - Title: "Accounts" - Admin panelinde görüntülenecek başlık
// - Icon: "key" - Navigasyon menüsünde gösterilecek ikon
// - Group: "System" - Menü grublandırması (Sistem kategorisi altında)
// - NavigationOrder: 50 - Menüde görüntülenme sırası
// - Visible: true - Admin panelinde görünür olması
// - FieldResolver: AccountFieldResolver - Alan tanımlarını yönetir
// - CardResolver: AccountCardResolver - Kart görünümlerini yönetir
// - Policy: AccountPolicy - Erişim kontrol politikasını yönetir
//
// Kullanım Örneği:
//   accountResource := NewAccountResource()
//   // Resource artık admin panelinde kullanılmaya hazır
//
// Önemli Notlar:
// - Fonksiyon, resource sisteminin başlatılması sırasında çağrılmalıdır
// - Tüm resolver'lar ve policy'ler bu fonksiyonda ayarlanmalıdır
// - Döndürülen pointer, resource registry'sine kaydedilmelidir
func NewAccountResource() *AccountResource {
	r := &AccountResource{}

	// Veritabanı modeli olarak Account entity'sini ayarla
	r.SetModel(&domainAccount.Account{})

	// URL'de kullanılacak slug'ı ayarla (örn: /admin/accounts)
	r.SetSlug("accounts")

	// Admin panelinde görüntülenecek başlığı ayarla
	r.SetTitle("Accounts")

	// Navigasyon menüsünde gösterilecek ikon'u ayarla
	r.SetIcon("key")

	// Menü grublandırması - "System" kategorisi altında görünecek
	r.SetGroup("System")

	// Menüde görüntülenme sırası (50. sırada)
	r.SetNavigationOrder(50)

	// Admin panelinde görünür olmasını sağla
	r.SetVisible(true)

	// Field resolver'ı ayarla - Account alanlarının tanımını yönetir
	// AccountFieldResolver, her alanın türünü, validasyonunu ve görünümünü belirler
	r.SetFieldResolver(&AccountFieldResolver{})

	// Card resolver'ı ayarla - Kart görünümlerini yönetir
	// AccountCardResolver, farklı görünüm modlarında (liste, detay vb.) verilerin sunumunu kontrol eder
	r.SetCardResolver(&AccountCardResolver{})

	// Policy'yi ayarla - Erişim kontrol ve yetkilendirme kurallarını yönetir
	// AccountPolicy, hangi kullanıcıların hangi işlemleri yapabileceğini belirler
	r.SetPolicy(&AccountPolicy{})

	// Yapılandırılmış resource pointer'ını döndür
	return r
}

// Bu metod, Account verilerini veritabanından almak için kullanılan repository'yi döner.
// Repository, GORM ORM kullanarak veritabanı işlemlerini yönetir.
//
// Parametreler:
// - db *gorm.DB: GORM veritabanı bağlantısı
//
// Dönüş Değeri:
// - data.DataProvider: Veritabanı işlemlerini gerçekleştiren provider interface'i
//
// Kullanım Senaryosu:
// - Admin panelinde Account listesi yüklenirken
// - Filtreleme, sıralama ve pagination işlemleri sırasında
// - Tekil Account kaydı alınırken
//
// Kullanım Örneği:
//   repo := accountResource.Repository(db)
//   accounts, err := repo.Get(ctx, nil)
//
// Önemli Notlar:
// - Döndürülen provider, GORM ORM'nin tüm özelliklerini kullanır
// - Veritabanı bağlantısı (db) geçerli ve açık olmalıdır
// - Provider, transaction desteği sağlar
func (r *AccountResource) Repository(db *gorm.DB) data.DataProvider {
	// GORM provider'ı oluştur ve Account modeli ile başlat
	return data.NewGormDataProvider(db, &domainAccount.Account{})
}

// Bu metod, Account verilerini yüklerken eager loading yapılacak ilişkileri belirtir.
// Eager loading, N+1 sorgu problemini önleyerek performansı artırır.
//
// Dönüş Değeri:
// - []string: Eager loading yapılacak ilişki adlarının listesi
//
// Eager Loading İlişkileri:
// - "User": Account'ın ilişkili olduğu User (Kullanıcı) kaydı
//
// Kullanım Senaryosu:
// - Account listesi yüklenirken, ilişkili User bilgileri de otomatik olarak yüklenir
// - Detay sayfasında Account ve User bilgileri birlikte gösterilir
// - N+1 sorgu problemini önler (her Account için ayrı User sorgusu yapılmaz)
//
// Kullanım Örneği:
//   with := accountResource.With()
//   // with = []string{"User"}
//   // Veritabanı sorgusu: SELECT * FROM accounts WITH User
//
// Önemli Notlar:
// - İlişki adları, Account model'inde tanımlanan GORM tag'larıyla eşleşmelidir
// - Eager loading, sorgu performansını artırır ancak bellek kullanımını da artırabilir
// - Gerekli olmayan ilişkileri eklemekten kaçının
// - Döndürülen liste boş ise, eager loading yapılmaz
func (r *AccountResource) With() []string {
	// User ilişkisini eager loading'e ekle
	return []string{"User"}
}

// Bu metod, Account resource'u için özel görünümleri (lens'leri) döner.
// Lens'ler, aynı verinin farklı perspektiflerden görüntülenmesini sağlar.
//
// Dönüş Değeri:
// - []resource.Lens: Tanımlanan lens'lerin listesi
//
// Kullanım Senaryosu:
// - Farklı kullanıcı rollerine göre farklı görünümler sunmak
// - Aynı Account verilerini farklı şekillerde göstermek
// - Örn: "Aktif Hesaplar", "Pasif Hesaplar", "Yönetici Hesapları" gibi filtreli görünümler
//
// Mevcut Durum:
// - Şu anda hiçbir lens tanımlanmamıştır (boş liste döner)
// - Gelecekte Account-spesifik lens'ler eklenebilir
//
// Kullanım Örneği:
//   lenses := accountResource.Lenses()
//   // Şu anda: lenses = []resource.Lens{}
//   // Gelecekte: lenses = []resource.Lens{
//   //   {Name: "Aktif", Filter: "status = 'active'"},
//   //   {Name: "Pasif", Filter: "status = 'inactive'"},
//   // }
//
// Önemli Notlar:
// - Lens'ler, admin panelinde hızlı filtre seçenekleri olarak görünür
// - Her lens, önceden tanımlanmış bir filtreleme kuralı içerir
// - Lens'ler, kullanıcı deneyimini iyileştirmek için kullanılır
func (r *AccountResource) Lenses() []resource.Lens {
	// Şu anda hiçbir özel görünüm tanımlanmamıştır
	return []resource.Lens{}
}

// Bu metod, Account resource'u için özel işlemleri (action'ları) döner.
// Action'lar, toplu işlemler veya özel işlevler için kullanılır.
//
// Dönüş Değeri:
// - []resource.Action: Tanımlanan action'ların listesi
//
// Kullanım Senaryosu:
// - Toplu silme, toplu aktivasyon/deaktivasyonu
// - Özel raporlar oluşturma
// - Dış sistemlere veri gönderme
// - Örn: "Seçili Hesapları Sil", "Seçili Hesapları Aktifleştir" gibi işlemler
//
// Mevcut Durum:
// - Şu anda hiçbir action tanımlanmamıştır (boş liste döner)
// - Gelecekte Account-spesifik action'lar eklenebilir
//
// Kullanım Örneği:
//   actions := accountResource.GetActions()
//   // Şu anda: actions = []resource.Action{}
//   // Gelecekte: actions = []resource.Action{
//   //   {Name: "Sil", Handler: deleteAccounts},
//   //   {Name: "Aktifleştir", Handler: activateAccounts},
//   // }
//
// Önemli Notlar:
// - Action'lar, admin panelinde toplu işlem düğmeleri olarak görünür
// - Her action, seçili kayıtlar üzerinde çalışır
// - Action'lar, yetkilendirme kontrolleri ile korunmalıdır
func (r *AccountResource) GetActions() []resource.Action {
	// Şu anda hiçbir özel işlem tanımlanmamıştır
	return []resource.Action{}
}

// Bu metod, Account resource'u için filtreleri döner.
// Filtreler, kullanıcıların verileri belirli kriterlere göre filtrelemesini sağlar.
//
// Dönüş Değeri:
// - []resource.Filter: Tanımlanan filtrelerin listesi
//
// Kullanım Senaryosu:
// - Account durumuna göre filtreleme (aktif, pasif, askıya alınmış)
// - Oluşturma tarihine göre filtreleme
// - Kullanıcı türüne göre filtreleme
// - Örn: "Durum", "Oluşturma Tarihi", "Hesap Türü" gibi filtreler
//
// Mevcut Durum:
// - Şu anda hiçbir filtre tanımlanmamıştır (boş liste döner)
// - Gelecekte Account-spesifik filtreler eklenebilir
//
// Kullanım Örneği:
//   filters := accountResource.GetFilters()
//   // Şu anda: filters = []resource.Filter{}
//   // Gelecekte: filters = []resource.Filter{
//   //   {Name: "Durum", Field: "status", Type: "select"},
//   //   {Name: "Oluşturma Tarihi", Field: "created_at", Type: "date"},
//   // }
//
// Önemli Notlar:
// - Filtreler, admin panelinde filtreleme paneli olarak görünür
// - Her filtre, belirli bir veritabanı alanına karşılık gelir
// - Filtreler, sorgu performansını etkileyebilir (indekslenmiş alanlar tercih edilir)
func (r *AccountResource) GetFilters() []resource.Filter {
	// Şu anda hiçbir filtre tanımlanmamıştır
	return []resource.Filter{}
}

// Bu metod, Account resource'u için varsayılan sıralama ayarlarını döner.
// Sıralama, admin panelinde Account listesinin hangi sıraya göre gösterileceğini belirler.
//
// Dönüş Değeri:
// - []resource.Sortable: Sıralama kurallarının listesi
//
// Varsayılan Sıralama:
// - Column: "created_at" - Oluşturma tarihi alanına göre sırala
// - Direction: "desc" - Azalan sırada (en yeni hesaplar önce)
//
// Kullanım Senaryosu:
// - Admin panelinde Account listesi açıldığında, en yeni hesaplar önce gösterilir
// - Kullanıcı, admin panelinde sıralamayı değiştirebilir
// - Varsayılan sıralama, en sık kullanılan sıralamadır
//
// Kullanım Örneği:
//   sortable := accountResource.GetSortable()
//   // sortable = []resource.Sortable{
//   //   {Column: "created_at", Direction: "desc"},
//   // }
//   // Veritabanı sorgusu: SELECT * FROM accounts ORDER BY created_at DESC
//
// Önemli Notlar:
// - Sıralama alanı, Account model'inde mevcut olmalıdır
// - Direction değerleri: "asc" (artan) veya "desc" (azalan)
// - Sıralama, sorgu performansını etkileyebilir (indekslenmiş alanlar tercih edilir)
// - Birden fazla sıralama kuralı tanımlanabilir (çok seviyeli sıralama)
// - Kullanıcı, admin panelinde sıralamayı değiştirebilir
func (r *AccountResource) GetSortable() []resource.Sortable {
	// Varsayılan sıralama: oluşturma tarihine göre azalan sırada (en yeni önce)
	return []resource.Sortable{
		{
			Column:    "created_at",
			Direction: "desc",
		},
	}
}
