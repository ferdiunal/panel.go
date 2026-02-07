package permission

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

// Bu yapı, permissions.toml dosyasında tanımlanan bir rolün yapısını temsil eder.
//
// Kullanım Senaryosu:
// - Uygulamada kullanıcı rollerinin tanımlanması ve yönetilmesi
// - Her rolün etiketini ve sahip olduğu izinleri depolamak
//
// Alanlar:
// - Label: Rolün insan tarafından okunabilir adı (örn: "Yönetici", "Kullanıcı")
// - Permissions: Rolün sahip olduğu izinlerin listesi (örn: ["users.create", "users.delete"])
//
// Örnek TOML Yapısı:
// [admin]
// label = "Yönetici"
// permissions = ["*"]
//
// [user]
// label = "Standart Kullanıcı"
// permissions = ["posts.read", "posts.create"]
type Role struct {
	Label       string   `toml:"label"`
	Permissions []string `toml:"permissions"`
}

// Bu yapı, permissions.toml dosyasının tamamını temsil eder.
// Dosya, üst düzey anahtarlar olarak roller içerir (örn: [admin], [user]).
// Bu nedenle, yapı bir harita (map) olarak parse edilir.
//
// Kullanım Senaryosu:
// - Tüm rol tanımlarını bellekte depolamak
// - Hızlı rol ve izin araması için
//
// Örnek:
// config := Config{
//     "admin": Role{Label: "Yönetici", Permissions: []string{"*"}},
//     "user": Role{Label: "Kullanıcı", Permissions: []string{"posts.read"}},
// }
type Config map[string]Role

// Bu yapı, izin konfigürasyonlarını yönetir ve kontrol eder.
// Uygulamada izin sistemi için merkezi yönetim noktası olarak görev yapar.
//
// Kullanım Senaryosu:
// - Uygulamada izin kontrolü yapmak
// - Rol ve izin bilgilerini sorgulamak
// - Kullanıcıların belirli işlemleri yapma yetkisini kontrol etmek
//
// Alanlar:
// - config: Yüklenen tüm rol ve izin konfigürasyonları
//
// Örnek Kullanım:
// manager, err := Load("permissions.toml")
// if err != nil {
//     log.Fatal(err)
// }
// if manager.HasPermission("admin", "users.delete") {
//     // Yönetici kullanıcıları silebilir
// }
type Manager struct {
	config Config
}

// Bu değişken, uygulamada global olarak erişilebilen Manager instance'ını tutar.
// GetInstance() fonksiyonu aracılığıyla erişilir.
//
// Önemli Not:
// - Singleton pattern kullanılır
// - Thread-safe olmayabilir, eğer concurrent erişim varsa senkronizasyon gerekebilir
var currentManager *Manager

// Bu fonksiyon, verilen yolda bulunan TOML dosyasını okur ve parse eder.
// Dosya, rol ve izin tanımlarını içerir ve Manager instance'ı oluşturur.
//
// Parametreler:
// - path: permissions.toml dosyasının tam yolu (örn: "/etc/app/permissions.toml")
//
// Dönüş Değerleri:
// - *Manager: Yüklenen konfigürasyonla oluşturulan Manager pointer'ı
// - error: Dosya okuma veya parse hatası (nil ise başarılı)
//
// Hata Senaryoları:
// - Dosya bulunamadığında: "permissions file could not be read"
// - TOML parse hatası: "permissions file could not be parsed"
//
// Kullanım Örneği:
// manager, err := Load("/etc/app/permissions.toml")
// if err != nil {
//     log.Fatalf("İzin dosyası yüklenemedi: %v", err)
// }
// // Artık manager kullanılabilir
//
// Önemli Notlar:
// - Fonksiyon, yüklenen Manager'ı global currentManager değişkenine atar
// - Aynı uygulamada birden fazla Load çağrısı, önceki konfigürasyonu üzerine yazar
func Load(path string) (*Manager, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("permissions file could not be read: %w", err)
	}

	var config Config
	if err := toml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("permissions file could not be parsed: %w", err)
	}

	mgr := &Manager{
		config: config,
	}
	currentManager = mgr
	return mgr, nil
}

// Bu metod, Manager tarafından yüklenen tüm konfigürasyonu döndürür.
// Eğer konfigürasyon nil ise, boş bir Config haritası döndürür.
//
// Alıcı:
// - m: Manager pointer'ı
//
// Dönüş Değeri:
// - Config: Tüm rol ve izin tanımlarını içeren harita
//
// Kullanım Örneği:
// manager, _ := Load("permissions.toml")
// config := manager.GetConfig()
// for roleName, roleData := range config {
//     fmt.Printf("Rol: %s, Etiket: %s\n", roleName, roleData.Label)
// }
//
// Önemli Notlar:
// - Nil kontrol yapılır, nil ise boş harita döndürülür
// - Döndürülen harita, orijinal konfigürasyonun referansıdır (deep copy değil)
// - Harita üzerinde yapılan değişiklikler Manager'ı etkileyebilir
func (m *Manager) GetConfig() Config {
	if m.config == nil {
		return make(Config)
	}
	return m.config
}

// Bu metod, konfigürasyonda tanımlanan tüm rol anahtarlarının listesini döndürür.
// Rol adlarını (key'leri) bir string slice'ı olarak döndürür.
//
// Alıcı:
// - m: Manager pointer'ı
//
// Dönüş Değeri:
// - []string: Tüm rol adlarının listesi (örn: ["admin", "user", "guest"])
//
// Kullanım Örneği:
// manager, _ := Load("permissions.toml")
// roles := manager.GetRoles()
// for _, role := range roles {
//     fmt.Printf("Mevcut Rol: %s\n", role)
// }
//
// Kullanım Senaryoları:
// - Uygulamada mevcut tüm rolleri listelemek
// - Rol seçim dropdown'ları doldurmak
// - Rol validasyonu yapmak
//
// Önemli Notlar:
// - Döndürülen slice'ın sırası garantili değildir (harita iterasyonu)
// - Boş konfigürasyon durumunda boş slice döndürülür
// - Slice kapasitesi, config haritasının boyutuna göre önceden tahsis edilir
func (m *Manager) GetRoles() []string {
	roles := make([]string, 0, len(m.config))
	for r := range m.config {
		roles = append(roles, r)
	}
	return roles
}

// Bu fonksiyon, global olarak yüklenen Manager instance'ını döndürür.
// Load() fonksiyonu tarafından ayarlanan currentManager değişkenini erişir.
//
// Dönüş Değeri:
// - *Manager: Global Manager pointer'ı (nil olabilir, eğer Load() çağrılmadıysa)
//
// Kullanım Örneği:
// // Uygulamanın başında
// Load("permissions.toml")
//
// // Uygulamanın başka yerinde
// manager := GetInstance()
// if manager != nil && manager.HasPermission("user", "posts.read") {
//     // İzin var
// }
//
// Kullanım Senaryoları:
// - Singleton pattern ile global Manager erişimi
// - Middleware'de izin kontrolü
// - HTTP handler'larında yetkilendirme
//
// Önemli Notlar:
// - Nil kontrolü yapılmalıdır, eğer Load() çağrılmadıysa nil döner
// - Thread-safe değildir, concurrent erişimde senkronizasyon gerekebilir
// - Uygulamada genellikle başlangıçta bir kez Load() çağrılır
func GetInstance() *Manager {
	return currentManager
}

// Bu metod, belirli bir rol anahtarı için Role tanımını döndürür.
// Rol bulunursa Role struct'ı ve true, bulunamazsa boş Role ve false döner.
//
// Alıcı:
// - m: Manager pointer'ı
//
// Parametreler:
// - roleKey: Sorgulanacak rolün anahtarı (örn: "admin", "user")
//
// Dönüş Değerleri:
// - Role: Bulunan rol tanımı (boş ise rol bulunamadı)
// - bool: Rolün bulunup bulunmadığını gösteren boolean (true = bulundu)
//
// Kullanım Örneği:
// manager, _ := Load("permissions.toml")
// if role, ok := manager.GetRole("admin"); ok {
//     fmt.Printf("Rol Adı: %s\n", role.Label)
//     fmt.Printf("İzinler: %v\n", role.Permissions)
// } else {
//     fmt.Println("Admin rolü bulunamadı")
// }
//
// Kullanım Senaryoları:
// - Belirli bir rolün detaylarını almak
// - Rol validasyonu yapmak
// - Rol bilgilerini UI'da göstermek
//
// Önemli Notlar:
// - Go'nun map erişim pattern'ı kullanılır (comma ok idiom)
// - Rol bulunamazsa boş Role struct'ı döner (Label: "", Permissions: nil)
// - Döndürülen Role, orijinal konfigürasyonun referansıdır
func (m *Manager) GetRole(roleKey string) (Role, bool) {
	role, ok := m.config[roleKey]
	return role, ok
}

// Bu metod, belirli bir rolün verilen izne sahip olup olmadığını kontrol eder.
// Rol bulunmazsa false döner. Rol bulunursa, izinler listesinde kontrol yapılır.
// Özel "*" izni, tüm izinleri temsil eder (wildcard).
//
// Alıcı:
// - m: Manager pointer'ı
//
// Parametreler:
// - roleName: Kontrol edilecek rolün adı (örn: "admin", "user")
// - permission: Kontrol edilecek izin (örn: "users.delete", "posts.create")
//
// Dönüş Değeri:
// - bool: İzin varsa true, yoksa false
//
// Kullanım Örneği:
// manager, _ := Load("permissions.toml")
//
// // Basit kontrol
// if manager.HasPermission("admin", "users.delete") {
//     fmt.Println("Yönetici kullanıcıları silebilir")
// }
//
// // Middleware'de kullanım
// func AuthMiddleware(next http.Handler) http.Handler {
//     return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//         userRole := r.Header.Get("X-User-Role")
//         if !GetInstance().HasPermission(userRole, "posts.create") {
//             http.Error(w, "Yetkisiz", http.StatusForbidden)
//             return
//         }
//         next.ServeHTTP(w, r)
//     })
// }
//
// Kullanım Senaryoları:
// - HTTP endpoint'lerinde yetkilendirme kontrolü
// - İş mantığında izin doğrulaması
// - Middleware'de erişim kontrolü
// - API endpoint'lerinde role-based access control (RBAC)
//
// İzin Kontrol Mantığı:
// 1. Rol bulunamadığında: false döner
// 2. Rol bulunduğunda, izinler listesi kontrol edilir:
//    - "*" (wildcard) bulunursa: true döner (tüm izinler)
//    - Tam eşleşme bulunursa: true döner
//    - Hiçbiri bulunmazsa: false döner
//
// Önemli Notlar:
// - Wildcard "*" izni, tüm diğer izinleri geçersiz kılar
// - İzin kontrolü case-sensitive'dir
// - Rol bulunamazsa false döner (hata fırlatmaz)
// - Performans: O(n) karmaşıklığında, n = izin sayısı
//
// Uyarılar:
// - Büyük izin listeleri için performans düşebilir
// - İzin adlandırması tutarlı olmalıdır (örn: "users.delete" vs "user.delete")
func (m *Manager) HasPermission(roleName string, permission string) bool {
	role, ok := m.config[roleName]
	if !ok {
		return false
	}

	for _, p := range role.Permissions {
		if p == "*" || p == permission {
			return true
		}
	}

	return false
}
