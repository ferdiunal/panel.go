// Bu paket, Session entity'si için admin panel resource tanımlarını içerir.
// Session yönetimi, kullanıcı oturumlarının izlenmesi ve yönetilmesi için gerekli
// tüm resource konfigürasyonlarını sağlar.
package session

import (
	"github.com/ferdiunal/panel.go/pkg/data"
	domainSession "github.com/ferdiunal/panel.go/pkg/domain/session"
	"github.com/ferdiunal/panel.go/pkg/resource"
	"gorm.io/gorm"
)

// Bu yapı, Session entity'si için admin panel resource tanımını temsil eder.
//
// SessionResource, OptimizedBase'i embed ederek tüm temel resource işlevlerini
// miras alır. Session verilerinin admin panelde nasıl gösterileceğini, filtreleneceğini,
// sıralanacağını ve yönetileceğini tanımlar.
//
// Kullanım Senaryoları:
// - Admin panelinde aktif oturumları görüntüleme
// - Kullanıcı oturumlarını yönetme ve izleme
// - Oturum verilerini filtreleme ve sıralama
// - Oturum güvenliği ve yönetimi
//
// Önemli Notlar:
// - Session verisi hassas bilgiler içerebilir, erişim kontrolleri önemlidir
// - Eager loading ile User ilişkisi otomatik yüklenir
// - Varsayılan sıralama en yeni oturumlar önce gösterilir
type SessionResource struct {
	resource.OptimizedBase
}

// Bu fonksiyon, yeni bir Session resource'u oluşturur ve yapılandırır.
//
// Parametreler: Yok
//
// Dönüş Değeri:
// - *SessionResource: Yapılandırılmış SessionResource pointer'ı
//
// Fonksiyon Açıklaması:
// NewSessionResource, Session entity'si için gerekli tüm konfigürasyonları
// yaparak hazır bir resource instance'ı döner. Bu fonksiyon:
// - Model tanımını ayarlar (domainSession.Session)
// - URL slug'ını "sessions" olarak belirler
// - Başlık ve ikonunu ayarlar
// - Sistem grubu altında organize eder
// - Navigation sırasını 51 olarak belirler
// - Field, Card ve Policy resolver'larını bağlar
//
// Kullanım Örneği:
//   sessionResource := NewSessionResource()
//   // sessionResource artık admin panelde kullanılmaya hazırdır
//
// Önemli Notlar:
// - Bu fonksiyon singleton pattern ile kullanılmalıdır
// - Döndürülen resource hemen kullanıma hazırdır
// - Tüm resolver'lar otomatik olarak bağlanır
// - Method chaining desteklenmez, direkt konfigürasyon yapılır
func NewSessionResource() *SessionResource {
	r := &SessionResource{}

	// Model tanımını Session domain entity'sine ayarla
	r.SetModel(&domainSession.Session{})

	// URL slug'ını "sessions" olarak belirle
	// Bu, admin panelinde /sessions URL'sinde erişilebilir olmasını sağlar
	r.SetSlug("sessions")

	// Admin panelinde görüntülenecek başlığı ayarla
	r.SetTitle("Sessions")

	// Sidebar'da gösterilecek ikonu ayarla
	// "clock" ikonu oturumların zaman tabanlı doğasını temsil eder
	r.SetIcon("clock")

	// Hangi grup altında organize edileceğini belirle
	// "System" grubu sistem yönetimi ile ilgili kaynakları içerir
	r.SetGroup("System")

	// Navigation menüsündeki sıra numarasını ayarla
	// Daha yüksek sayılar menüde daha aşağıda görünür
	r.SetNavigationOrder(51)

	// Resource'un admin panelinde görünür olup olmayacağını belirle
	r.SetVisible(true)

	// Kayıt başlığı için "id" field'ını kullan
	// Session kayıtları için varsayılan olarak ID gösterilir
	r.SetRecordTitleKey("id")

	// Field resolver'ı ayarla - form alanlarının nasıl render edileceğini tanımlar
	r.SetFieldResolver(&SessionFieldResolver{})

	// Card resolver'ı ayarla - liste görünümündeki kartların nasıl gösterileceğini tanımlar
	r.SetCardResolver(&SessionCardResolver{})

	// Policy'yi ayarla - erişim kontrolleri ve yetkilendirme kurallarını tanımlar
	r.SetPolicy(&SessionPolicy{})

	// Yapılandırılmış SessionResource pointer'ı döndür
	return r
}

// Bu metod, Session verilerine erişmek için veri sağlayıcısını döner.
//
// Parametreler:
// - db *gorm.DB: GORM veritabanı bağlantısı
//
// Dönüş Değeri:
// - data.DataProvider: Session verilerine erişim sağlayan veri sağlayıcı
//
// Metod Açıklaması:
// Repository metodu, verilen GORM veritabanı bağlantısını kullanarak
// Session entity'si için bir veri sağlayıcı oluşturur. Bu sağlayıcı,
// CRUD işlemleri (Create, Read, Update, Delete) için kullanılır.
//
// Kullanım Örneği:
//   sessionResource := NewSessionResource()
//   provider := sessionResource.Repository(db)
//   sessions, err := provider.Get(ctx)
//
// Önemli Notlar:
// - Her çağrıda yeni bir DataProvider instance'ı oluşturulur
// - GORM bağlantısı nil olmamalıdır
// - DataProvider, veritabanı işlemlerini optimize eder
// - Lazy loading yerine eager loading tercih edilir (With() metodu ile)
func (r *SessionResource) Repository(db *gorm.DB) data.DataProvider {
	// GORM veri sağlayıcısını oluştur ve döndür
	// Bu sağlayıcı Session entity'si için tüm veritabanı işlemlerini yönetir
	return data.NewGormDataProvider(db, &domainSession.Session{})
}

// Bu metod, eager loading yapılacak ilişkileri belirtir.
//
// Parametreler: Yok
//
// Dönüş Değeri:
// - []string: Eager loading yapılacak ilişki adlarının slice'ı
//
// Metod Açıklaması:
// With metodu, Session verisi yüklenirken otomatik olarak yüklenmesi gereken
// ilişkili entity'leri belirtir. Bu, N+1 sorgu problemini önler ve performansı
// artırır. Döndürülen string'ler GORM'un Preload() metoduna geçirilir.
//
// Kullanım Örneği:
//   sessionResource := NewSessionResource()
//   relations := sessionResource.With()
//   // relations = []string{"User"}
//   // GORM: db.Preload("User").Find(&sessions)
//
// Önemli Notlar:
// - "User" ilişkisi her Session yüklemesinde otomatik olarak yüklenir
// - Bu, Session'ın hangi kullanıcıya ait olduğunu göstermek için gereklidir
// - Fazla ilişki yüklemek performansı olumsuz etkileyebilir
// - İlişki adları domain model'deki tag'larla eşleşmelidir
// - Döndürülen slice boş olabilir, bu durumda eager loading yapılmaz
func (r *SessionResource) With() []string {
	// User ilişkisini eager loading yapılacak ilişkiler listesine ekle
	// Bu, her Session kaydı yüklenirken ilişkili User'ı da yükler
	return []string{"User"}
}

// Bu metod, Session resource'u için özel görünümleri (lenses) döner.
//
// Parametreler: Yok
//
// Dönüş Değeri:
// - []resource.Lens: Özel görünümlerin slice'ı
//
// Metod Açıklaması:
// Lenses metodu, Session verilerinin farklı perspektiflerden görüntülenmesini
// sağlayan özel görünümleri tanımlar. Örneğin, "Aktif Oturumlar", "Süresi Dolan Oturumlar"
// gibi filtrelenmiş görünümler oluşturabilir.
//
// Kullanım Örneği:
//   sessionResource := NewSessionResource()
//   lenses := sessionResource.Lenses()
//   // Şu anda boş, gelecekte özel görünümler eklenebilir
//
// Önemli Notlar:
// - Şu anda hiçbir özel görünüm tanımlanmamıştır
// - Gelecekte "Aktif", "Süresi Dolan", "Güvenlik Uyarıları" gibi lensler eklenebilir
// - Her lens, belirli filtreleme ve sıralama kurallarını içerebilir
// - Boş slice döndürülmesi, varsayılan görünümün kullanılacağı anlamına gelir
func (r *SessionResource) Lenses() []resource.Lens {
	// Şu anda özel görünüm tanımlanmamıştır
	// Gelecekte Session'lar için özel filtrelenmiş görünümler eklenebilir
	return []resource.Lens{}
}

// Bu metod, Session resource'u için özel işlemleri (actions) döner.
//
// Parametreler: Yok
//
// Dönüş Değeri:
// - []resource.Action: Özel işlemlerin slice'ı
//
// Metod Açıklaması:
// GetActions metodu, Session kayıtları üzerinde gerçekleştirilebilecek özel
// işlemleri tanımlar. Örneğin, "Oturumu Sonlandır", "Oturumu Doğrula" gibi
// custom action'lar eklenebilir.
//
// Kullanım Örneği:
//   sessionResource := NewSessionResource()
//   actions := sessionResource.GetActions()
//   // Şu anda boş, gelecekte özel işlemler eklenebilir
//
// Önemli Notlar:
// - Şu anda hiçbir özel işlem tanımlanmamıştır
// - Gelecekte "Oturumu Sonlandır", "Oturumu Doğrula", "Güvenlik Taraması" gibi
//   action'lar eklenebilir
// - Her action, belirli bir işlemi gerçekleştiren bir handler içerir
// - Boş slice döndürülmesi, sadece standart CRUD işlemlerinin kullanılacağı
//   anlamına gelir
func (r *SessionResource) GetActions() []resource.Action {
	// Şu anda özel işlem tanımlanmamıştır
	// Gelecekte Session'lar üzerinde gerçekleştirilebilecek custom action'lar
	// eklenebilir (örn: oturumu sonlandırma, doğrulama vb.)
	return []resource.Action{}
}

// Bu metod, Session resource'u için filtreleri döner.
//
// Parametreler: Yok
//
// Dönüş Değeri:
// - []resource.Filter: Filtrelerin slice'ı
//
// Metod Açıklaması:
// GetFilters metodu, admin panelinde Session verilerini filtrelemek için
// kullanılabilecek filtreleri tanımlar. Örneğin, "Kullanıcıya Göre", "Duruma Göre",
// "Tarih Aralığına Göre" gibi filtreler eklenebilir.
//
// Kullanım Örneği:
//   sessionResource := NewSessionResource()
//   filters := sessionResource.GetFilters()
//   // Şu anda boş, gelecekte filtreler eklenebilir
//
// Önemli Notlar:
// - Şu anda hiçbir filtre tanımlanmamıştır
// - Gelecekte "Kullanıcı", "Durum", "Oluşturulma Tarihi", "Son Aktivite" gibi
//   filtreler eklenebilir
// - Her filtre, belirli bir alan üzerinde filtreleme sağlar
// - Boş slice döndürülmesi, filtreleme özelliğinin devre dışı olduğu anlamına gelir
func (r *SessionResource) GetFilters() []resource.Filter {
	// Şu anda filtre tanımlanmamıştır
	// Gelecekte Session'ları filtrelemek için filtreler eklenebilir
	// (örn: kullanıcıya göre, duruma göre, tarih aralığına göre vb.)
	return []resource.Filter{}
}

// Bu metod, Session resource'u için varsayılan sıralama ayarlarını döner.
//
// Parametreler: Yok
//
// Dönüş Değeri:
// - []resource.Sortable: Sıralama ayarlarının slice'ı
//
// Metod Açıklaması:
// GetSortable metodu, admin panelinde Session verilerinin varsayılan olarak
// nasıl sıralanacağını belirtir. Birden fazla sıralama kriteri tanımlanabilir.
// Döndürülen slice'daki sıra, sıralama önceliğini belirtir.
//
// Kullanım Örneği:
//   sessionResource := NewSessionResource()
//   sortables := sessionResource.GetSortable()
//   // sortables[0].Column = "created_at"
//   // sortables[0].Direction = "desc"
//   // GORM: db.Order("created_at DESC")
//
// Önemli Notlar:
// - Varsayılan sıralama "created_at" alanına göre azalan (desc) sırada yapılır
// - Bu, en yeni oturumların önce gösterilmesini sağlar
// - Sıralama yönü "asc" (artan) veya "desc" (azalan) olabilir
// - Birden fazla sıralama kriteri tanımlanabilir (örn: created_at DESC, user_id ASC)
// - Sıralama alanları veritabanında mevcut olmalıdır
func (r *SessionResource) GetSortable() []resource.Sortable {
	// Varsayılan sıralama ayarlarını tanımla
	return []resource.Sortable{
		{
			// Sıralama yapılacak kolon adı
			Column: "created_at",
			// Sıralama yönü: "desc" = azalan (en yeni önce)
			Direction: "desc",
		},
	}
}
