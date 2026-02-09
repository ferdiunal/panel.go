// Package verification, doğrulama (verification) işlemleri için resource tanımlarını içerir.
// Bu paket, panel uygulamasında doğrulama verilerinin yönetimi, görüntülenmesi ve işlenmesi
// için gerekli tüm resource yapılarını sağlar.
package verification

import (
	"github.com/ferdiunal/panel.go/pkg/data"
	domainVerification "github.com/ferdiunal/panel.go/pkg/domain/verification"
	"github.com/ferdiunal/panel.go/pkg/resource"
	"gorm.io/gorm"
)

// Bu yapı, Verification entity'si için resource tanımını temsil eder.
// VerificationResource, panel uygulamasında doğrulama verilerinin CRUD işlemleri,
// filtreleme, sıralama ve görüntülenmesini yönetir.
//
// Kullanım Senaryoları:
// - Doğrulama kayıtlarının listelenmesi ve yönetimi
// - Doğrulama durumunun takip edilmesi
// - Sistem tarafından gerçekleştirilen doğrulama işlemlerinin kaydedilmesi
// - Doğrulama geçmişinin görüntülenmesi
//
// Önemli Notlar:
// - OptimizedBase'i embed ederek temel resource işlevselliğini miras alır
// - Doğrulama verilerine erişim kontrol politikası tarafından korunur
// - Sistem grubu altında navigasyonda 52. sırada görüntülenir
type VerificationResource struct {
	resource.OptimizedBase
}

// Bu fonksiyon, yeni bir Verification resource'u oluşturur ve yapılandırır.
// Fonksiyon, resource'un tüm gerekli ayarlarını (model, slug, başlık, ikon, grup, vb.)
// yapılandırarak kullanıma hazır bir resource örneği döner.
//
// Parametreler: Yok
//
// Dönüş Değeri:
// - *VerificationResource: Yapılandırılmış Verification resource pointer'ı
//
// Yapılandırılan Özellikler:
// - Model: domainVerification.Verification (veritabanı modeli)
// - Slug: "verifications" (URL'de kullanılan tanımlayıcı)
// - Başlık: "Verifications" (UI'da gösterilen başlık)
// - İkon: "shield-check" (navigasyon menüsünde gösterilen ikon)
// - Grup: "System" (resource'un ait olduğu grup)
// - Navigasyon Sırası: 52 (menüde gösterilme sırası)
// - Görünürlük: true (resource menüde görünür)
// - Field Resolver: VerificationFieldResolver (alan çözümleyici)
// - Card Resolver: VerificationCardResolver (kart görünümü çözümleyici)
// - Policy: VerificationPolicy (erişim kontrol politikası)
//
// Kullanım Örneği:
//   resource := NewVerificationResource()
//   // resource artık panel uygulamasında kullanılmaya hazırdır
//
// Önemli Notlar:
// - Fonksiyon her çağrıldığında yeni bir resource örneği oluşturur
// - Tüm ayarlar method chaining kullanılarak yapılandırılır
// - Resource, VerificationFieldResolver ve VerificationCardResolver'ı kullanır
func NewVerificationResource() *VerificationResource {
	r := &VerificationResource{}

	// Veritabanı modelini ayarla - Verification entity'sini kullan
	r.SetModel(&domainVerification.Verification{})

	// URL'de kullanılacak slug'ı ayarla - "/verifications" endpoint'i oluşturur
	r.SetSlug("verifications")

	// UI'da gösterilecek başlığı ayarla
	r.SetTitle("Verifications")

	// Navigasyon menüsünde gösterilecek ikon'u ayarla (shield-check: güvenlik simgesi)
	r.SetIcon("shield-check")

	// Resource'un ait olduğu grubu ayarla - "System" grubu altında görüntülenir
	r.SetGroup("System")

	// Navigasyon menüsünde gösterilme sırasını ayarla (52. sırada)
	r.SetNavigationOrder(52)

	// Resource'u menüde görünür yap
	r.SetVisible(true)

	// Kayıt başlığı için "id" field'ını kullan
	// Verification kayıtları için varsayılan olarak ID gösterilir
	r.SetRecordTitleKey("id")

	// Alan çözümleyicisini ayarla - alanların nasıl gösterileceğini belirler
	r.SetFieldResolver(&VerificationFieldResolver{})

	// Kart görünümü çözümleyicisini ayarla - kart formatında gösterilişi belirler
	r.SetCardResolver(&VerificationCardResolver{})

	// Erişim kontrol politikasını ayarla - kim hangi işlemleri yapabilir belirler
	r.SetPolicy(&VerificationPolicy{})

	// Yapılandırılmış resource pointer'ını döner
	return r
}

// Bu metod, Verification entity'si için veritabanı repository'sini oluşturur ve döner.
// Repository, veritabanı işlemleri (CRUD) için gerekli tüm fonksiyonları sağlar.
//
// Parametreler:
// - db *gorm.DB: GORM veritabanı bağlantısı
//
// Dönüş Değeri:
// - data.DataProvider: Veritabanı işlemleri için provider interface'i
//
// Kullanım Senaryoları:
// - Doğrulama kayıtlarını veritabanından almak
// - Yeni doğrulama kaydı oluşturmak
// - Mevcut doğrulama kaydını güncellemek
// - Doğrulama kaydını silmek
// - Doğrulama kayıtlarını filtrelemek ve sıralamak
//
// Kullanım Örneği:
//   db := gorm.Open(...)
//   resource := NewVerificationResource()
//   provider := resource.Repository(db)
//   // provider artık veritabanı işlemleri için kullanılabilir
//
// Önemli Notlar:
// - Her çağrıda yeni bir DataProvider örneği oluşturulur
// - GORM ORM'i kullanarak veritabanı işlemleri gerçekleştirilir
// - Verification modeli otomatik olarak tablo adı ve yapısını belirler
func (r *VerificationResource) Repository(db *gorm.DB) data.DataProvider {
	return data.NewGormDataProvider(db, &domainVerification.Verification{})
}

// Bu metod, Verification resource'u için özel görünümleri (lens'leri) döner.
// Lens'ler, aynı verinin farklı perspektiflerden görüntülenmesini sağlar.
//
// Parametreler: Yok
//
// Dönüş Değeri:
// - []resource.Lens: Tanımlı lens'lerin dilimi (şu anda boş)
//
// Kullanım Senaryoları:
// - Doğrulama verilerinin farklı görünümlerini oluşturmak
// - Örneğin: "Başarılı Doğrulamalar", "Başarısız Doğrulamalar", "Beklemede Olanlar" gibi filtreli görünümler
// - Her lens, belirli bir filtreleme ve sıralama kombinasyonunu temsil eder
//
// Mevcut Durum:
// - Şu anda hiçbir lens tanımlanmamıştır (boş dilim döner)
// - Gelecekte özel görünümler eklenebilir
//
// Kullanım Örneği:
//   resource := NewVerificationResource()
//   lenses := resource.Lenses()
//   // lenses şu anda boş bir dilim içerir
//
// Önemli Notlar:
// - Lens'ler, kullanıcıya hızlı erişim sağlayan önceden tanımlanmış filtrelerdir
// - Her lens, belirli bir iş mantığı senaryosunu temsil eder
// - Gelecekte genişletilebilir bir yapıdır
func (r *VerificationResource) Lenses() []resource.Lens {
	return []resource.Lens{}
}

// Bu metod, Verification resource'u için özel işlemleri (action'ları) döner.
// Action'lar, kayıtlar üzerinde gerçekleştirilebilecek özel işlemleri temsil eder.
//
// Parametreler: Yok
//
// Dönüş Değeri:
// - []resource.Action: Tanımlı action'ların dilimi (şu anda boş)
//
// Kullanım Senaryoları:
// - Doğrulama kayıtları üzerinde toplu işlemler yapmak
// - Örneğin: "Tümünü Onayla", "Tümünü Reddet", "E-posta Gönder" gibi işlemler
// - Her action, belirli bir iş mantığı işlemini temsil eder
//
// Mevcut Durum:
// - Şu anda hiçbir action tanımlanmamıştır (boş dilim döner)
// - Gelecekte özel işlemler eklenebilir
//
// Kullanım Örneği:
//   resource := NewVerificationResource()
//   actions := resource.GetActions()
//   // actions şu anda boş bir dilim içerir
//
// Önemli Notlar:
// - Action'lar, UI'da butonlar veya menü öğeleri olarak görüntülenir
// - Her action, belirli bir iş mantığı işlemini gerçekleştirir
// - Gelecekte genişletilebilir bir yapıdır
func (r *VerificationResource) GetActions() []resource.Action {
	return []resource.Action{}
}

// Bu metod, Verification resource'u için filtreleri döner.
// Filtreler, kullanıcıların doğrulama kayıtlarını belirli kriterlere göre filtrelemesini sağlar.
//
// Parametreler: Yok
//
// Dönüş Değeri:
// - []resource.Filter: Tanımlı filtrelerin dilimi (şu anda boş)
//
// Kullanım Senaryoları:
// - Doğrulama durumuna göre filtreleme (başarılı, başarısız, beklemede)
// - Doğrulama türüne göre filtreleme (e-posta, telefon, vb.)
// - Tarih aralığına göre filtreleme
// - Kullanıcıya göre filtreleme
//
// Mevcut Durum:
// - Şu anda hiçbir filtre tanımlanmamıştır (boş dilim döner)
// - Gelecekte özel filtreler eklenebilir
//
// Kullanım Örneği:
//   resource := NewVerificationResource()
//   filters := resource.GetFilters()
//   // filters şu anda boş bir dilim içerir
//
// Önemli Notlar:
// - Filtreler, UI'da dropdown veya checkbox'lar olarak görüntülenir
// - Her filtre, belirli bir sütun veya alan üzerinde çalışır
// - Gelecekte genişletilebilir bir yapıdır
func (r *VerificationResource) GetFilters() []resource.Filter {
	return []resource.Filter{}
}

// Bu metod, Verification resource'u için varsayılan sıralama ayarlarını döner.
// Sıralama, doğrulama kayıtlarının hangi sütuna göre ve hangi yönde sıralanacağını belirler.
//
// Parametreler: Yok
//
// Dönüş Değeri:
// - []resource.Sortable: Sıralama ayarlarının dilimi
//
// Varsayılan Sıralama:
// - Sütun: "created_at" (oluşturulma tarihi)
// - Yön: "desc" (azalan sıra - en yeni kayıtlar önce)
//
// Kullanım Senaryoları:
// - Doğrulama kayıtlarını en yeni olanlardan başlayarak göstermek
// - Kullanıcıların en son doğrulama işlemlerini hızlıca görmesini sağlamak
// - Varsayılan olarak en güncel verileri sunmak
//
// Kullanım Örneği:
//   resource := NewVerificationResource()
//   sortables := resource.GetSortable()
//   // sortables: [{Column: "created_at", Direction: "desc"}]
//   // Kayıtlar oluşturulma tarihine göre en yeniden en eskiye sıralanır
//
// Önemli Notlar:
// - Sıralama, liste görünümünde varsayılan olarak uygulanır
// - Kullanıcılar UI'da sıralamayı değiştirebilir
// - "desc" = azalan sıra (Z'den A'ya), "asc" = artan sıra (A'dan Z'ye)
// - created_at sütunu, veritabanında otomatik olarak ayarlanır
func (r *VerificationResource) GetSortable() []resource.Sortable {
	return []resource.Sortable{
		{
			Column:    "created_at",
			Direction: "desc",
		},
	}
}
