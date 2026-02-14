package panel

import (
	"time"

	"github.com/ferdiunal/panel.go/pkg/page"
	"github.com/ferdiunal/panel.go/pkg/resource"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

// / # FeatureConfig - Panel Özellik Yapılandırması
// /
// / Panel uygulamasının opsiyonel özelliklerini açıp kapatmak için kullanılan yapıdır.
// / Bu yapı, panelin hangi özelliklerin etkin olacağını kontrol etmeyi sağlar.
// /
// / ## Kullanım Senaryoları
// / - Kullanıcı kaydını devre dışı bırakmak (sadece admin tarafından oluşturulan hesaplar)
// / - Şifre sıfırlama özelliğini kapatmak (LDAP/SSO entegrasyonunda)
// / - Belirli ortamlarda (test, staging) özellikleri kontrol etmek
// /
// / ## Örnek Kullanım
// / ```go
// / config := FeatureConfig{
// /     Register:       true,        // Kullanıcı kaydını etkinleştir
// /     ForgotPassword: true,        // Şifre sıfırlama özelliğini etkinleştir
// / }
// / ```
// /
// / ## Avantajlar
// / - Basit boolean yapısı ile kolay kontrol
// / - Runtime'da özellikleri değiştirmek mümkün
// / - Farklı ortamlar için farklı konfigürasyonlar oluşturulabilir
// /
// / ## Önemli Notlar
// / - Varsayılan değer false'tur (özellikleri açıkça etkinleştirmeniz gerekir)
// / - Veritabanı ayarlarından da kontrol edilebilir (SettingsConfig)
type FeatureConfig struct {
	/// Register, kullanıcı kayıt özelliğini aktif eder.
	/// true: Yeni kullanıcılar kendilerini kaydedebilir
	/// false: Sadece admin tarafından oluşturulan hesaplar kullanılabilir
	Register bool

	/// ForgotPassword, şifre sıfırlama özelliğini aktif eder.
	/// true: Kullanıcılar şifrelerini sıfırlayabilir
	/// false: Şifre sıfırlama devre dışı (admin tarafından sıfırlanması gerekir)
	ForgotPassword bool
}

// / # OAuthConfig - OAuth Sağlayıcı Yapılandırması
// /
// / Panelde kullanılacak OAuth sağlayıcılarının yapılandırmasını tutar.
// / Şu anda sadece Google OAuth desteklenmektedir.
// /
// / ## Kullanım Senaryoları
// / - Google hesabı ile giriş sağlamak
// / - Sosyal medya entegrasyonu
// / - Kurumsal SSO çözümleri
// /
// / ## Örnek Kullanım
// / ```go
// / oauthConfig := OAuthConfig{
// /     Google: GoogleConfig{
// /         ClientID:     "your-client-id.apps.googleusercontent.com",
// /         ClientSecret: "your-client-secret",
// /         RedirectURL:  "https://example.com/auth/google/callback",
// /     },
// / }
// / ```
// /
// / ## Gelecek Genişletmeler
// / - GitHub OAuth desteği
// / - Microsoft OAuth desteği
// / - OIDC (OpenID Connect) desteği
type OAuthConfig struct {
	/// Google, Google OAuth sağlayıcısının yapılandırmasını içerir
	Google GoogleConfig
}

// / # GoogleConfig - Google OAuth Yapılandırması
// /
// / Google OAuth 2.0 entegrasyonu için gerekli yapılandırma bilgilerini tutar.
// / Google Cloud Console'dan alınan kimlik bilgilerini içerir.
// /
// / ## Kullanım Senaryoları
// / - Google hesabı ile giriş
// / - Google Workspace entegrasyonu
// / - Kurumsal ortamlarda SSO
// /
// / ## Örnek Kullanım
// / ```go
// / googleConfig := GoogleConfig{
// /     ClientID:     "123456789-abcdefghijklmnopqrstuvwxyz.apps.googleusercontent.com",
// /     ClientSecret: "GOCSPX-1234567890abcdefghijklmno",
// /     RedirectURL:  "https://panel.example.com/auth/google/callback",
// / }
// /
// / if googleConfig.Enabled() {
// /     // Google OAuth'u etkinleştir
// / }
// / ```
// /
// / ## Önemli Notlar
// / - ClientSecret asla client-side'da açığa çıkarılmamalıdır
// / - RedirectURL, Google Cloud Console'da kayıtlı olmalıdır
// / - Üretim ortamında HTTPS kullanılmalıdır
// / - ClientID ve ClientSecret çevre değişkenlerinden yüklenmelidir
// /
// / ## Güvenlik Uyarıları
// / - Kimlik bilgilerini kaynak kodunda saklamayın
// / - .env dosyasını .gitignore'a ekleyin
// / - Düzenli olarak token'ları döndürün
type GoogleConfig struct {
	/// ClientID, Google Cloud Console'dan alınan OAuth 2.0 Client ID'sidir.
	/// Format: "123456789-abcdefghijklmnopqrstuvwxyz.apps.googleusercontent.com"
	ClientID string

	/// ClientSecret, Google Cloud Console'dan alınan OAuth 2.0 Client Secret'ıdır.
	/// Format: "GOCSPX-1234567890abcdefghijklmno"
	/// UYARI: Bu değer asla client-side'da açığa çıkarılmamalıdır!
	ClientSecret string

	/// RedirectURL, Google OAuth akışından sonra kullanıcının yönlendirileceği URL'dir.
	/// Örnek: "https://panel.example.com/auth/google/callback"
	/// UYARI: Bu URL, Google Cloud Console'da kayıtlı olmalıdır
	RedirectURL string
}

// / # Enabled - Google OAuth Yapılandırmasının Geçerliliğini Kontrol Eder
// /
// / Google OAuth'un etkinleştirilip etkinleştirilmediğini kontrol eder.
// / ClientID ve ClientSecret'ın boş olmadığını doğrular.
// /
// / ## Dönüş Değeri
// / - true: Google OAuth yapılandırması geçerli ve etkindir
// / - false: Google OAuth yapılandırması eksik veya geçersizdir
// /
// / ## Kullanım Senaryoları
// / - Giriş sayfasında Google butonu gösterilip gösterilmeyeceğini belirlemek
// / - OAuth middleware'inde yapılandırma kontrolü
// / - Başlangıç sırasında yapılandırma doğrulaması
// /
// / ## Örnek Kullanım
// / ```go
// / config := GoogleConfig{
// /     ClientID:     "123456789-abc.apps.googleusercontent.com",
// /     ClientSecret: "GOCSPX-1234567890abc",
// /     RedirectURL:  "https://example.com/callback",
// / }
// /
// / if config.Enabled() {
// /     // Google OAuth'u etkinleştir
// /     setupGoogleOAuth(config)
// / } else {
// /     // Google OAuth'u devre dışı bırak
// /     log.Warn("Google OAuth yapılandırması eksik")
// / }
// / ```
// /
// / ## Önemli Notlar
// / - RedirectURL boş olsa bile true döndürebilir (opsiyonel kontrol)
// / - ClientID ve ClientSecret zorunludur
// / - Üretim ortamında bu fonksiyon başlangıçta çağrılmalıdır
func (c GoogleConfig) Enabled() bool {
	return c.ClientID != "" && c.ClientSecret != ""
}

// / # CORSConfig - Cross-Origin Resource Sharing Yapılandırması
// /
// / Panelin farklı origin'lerden gelen istekleri kabul etmesini kontrol eder.
// / Web uygulamaları arasında veri paylaşımını güvenli bir şekilde sağlar.
// /
// / ## Kullanım Senaryoları
// / - Frontend ve backend farklı domain'lerde çalışıyor
// / - Mobil uygulamalar panel API'sini kullanıyor
// / - Microservices mimarisinde API erişimi
// / - Geliştirme ortamında localhost'tan erişim
// /
// / ## Örnek Kullanım
// / ```go
// / // Üretim ortamı
// / corsConfig := CORSConfig{
// /     AllowedOrigins: []string{
// /         "https://app.example.com",
// /         "https://admin.example.com",
// /     },
// / }
// /
// / // Geliştirme ortamı
// / devCorsConfig := CORSConfig{
// /     AllowedOrigins: []string{
// /         "http://localhost:3000",
// /         "http://localhost:5173",
// /         "http://127.0.0.1:8080",
// /     },
// / }
// / ```
// /
// / ## Avantajlar
// / - Güvenli cross-origin istekleri sağlar
// / - Farklı ortamlar için farklı konfigürasyonlar
// / - Dinamik olarak değiştirilebilir
// /
// / ## Güvenlik Uyarıları
// / - AllowedOrigins'e "*" (wildcard) eklemeyin (güvenlik riski)
// / - Sadece güvenilir origin'leri ekleyin
// / - Üretim ortamında localhost'u eklemeyin
// / - HTTPS kullanın
// /
// / ## Önemli Notlar
// / - Boş liste tüm origin'leri reddeder
// / - Origin kontrolü case-sensitive'dir
// / - Port numarası origin'in bir parçasıdır
type CORSConfig struct {
	/// AllowedOrigins, CORS isteklerine izin verilen origin'lerin listesidir.
	/// Örnek: []string{"https://example.com", "https://app.example.com"}
	/// Geliştirme için: []string{"http://localhost:3000", "http://localhost:5173"}
	///
	/// UYARI: Wildcard "*" kullanmayın, güvenlik riski oluşturur!
	AllowedOrigins []string
}

// / # Config - Panel Genel Yapılandırması
// /
// / Panelin tüm yapılandırma ayarlarını merkezi olarak tutar.
// / Sunucu, veritabanı, OAuth, depolama ve diğer tüm ayarları içerir.
// /
// / ## Yapı Bileşenleri
// / - Server: HTTP sunucu ayarları (port, host)
// / - Database: Veritabanı bağlantı bilgileri
// / - Environment: Çalışma ortamı (production, development, test)
// / - Features: Etkinleştirilmiş özellikler
// / - OAuth: OAuth sağlayıcı yapılandırması
// / - Storage: Dosya depolama ayarları
// / - Permissions: İzin yönetimi yapılandırması
// / - SettingsValues: Veritabanından gelen dinamik ayarlar
// / - Pages: Özel sayfalar (Dashboard, Settings, Account vb.)
// / - UserResource: Kullanıcı resource'u
// / - Resources: Ek resource'lar
// / - CORS: CORS yapılandırması
// /
// / ## Kullanım Senaryoları
// / - Panel başlangıcında yapılandırma yükleme
// / - Farklı ortamlar için farklı konfigürasyonlar
// / - Runtime'da ayarları değiştirme
// / - Yapılandırma doğrulaması
// /
// / ## Örnek Kullanım
// / ```go
// / config := Config{
// /     Server: ServerConfig{
// /         Port: "8080",
// /         Host: "0.0.0.0",
// /     },
// /     Database: DatabaseConfig{
// /         DSN:    "user:password@tcp(localhost:3306)/panel",
// /         Driver: "mysql",
// /     },
// /     Environment: "production",
// /     Features: FeatureConfig{
// /         Register:       true,
// /         ForgotPassword: true,
// /     },
// /     CORS: CORSConfig{
// /         AllowedOrigins: []string{"https://example.com"},
// /     },
// / }
// / ```
// /
// / ## Önemli Notlar
// / - Config yapısı thread-safe değildir, başlangıçta ayarlanmalıdır
// / - Veritabanı bağlantısı Instance alanında saklanır
// / - Resources dinamik olarak eklenebilir
// / - Pages nil olabilir (varsayılan sayfalar otomatik oluşturulur)
// /
// / ## Güvenlik Uyarıları
// / - Kimlik bilgilerini kaynak kodunda saklamayın
// / - Çevre değişkenlerinden yükleyin
// / - Üretim ortamında debug modu kapatın
type Config struct {
	/// Server, HTTP sunucu ayarlarını tutar (port, host)
	Server ServerConfig

	/// Database, veritabanı bağlantı bilgilerini tutar
	Database DatabaseConfig

	/// Environment, panelin çalışma ortamını belirtir
	/// Değerler: "production", "development", "test"
	Environment string

	/// Features, panelin etkinleştirilmiş özelliklerini tutar
	Features FeatureConfig

	/// OAuth, OAuth sağlayıcılarının yapılandırmasını tutar
	OAuth OAuthConfig

	/// Storage, dosya yükleme ve depolama ayarlarını tutar
	Storage StorageConfig

	/// Permissions, izin yönetimi yapılandırmasını tutar
	Permissions PermissionConfig

	/// SettingsValues, veritabanından gelen dinamik ayarları tutar
	SettingsValues SettingsConfig

	/// Pages, panele eklenen özel sayfaları tutar (Dashboard, Settings, Account vb.)
	/// Eğer boş ise varsayılan Dashboard, Settings ve Account sayfaları otomatik oluşturulur.
	/// Kullanıcı kendi sayfalarını ekleyerek varsayılan sayfaları geçersiz kılabilir.
	/// Aynı slug'a sahip kullanıcı sayfası, varsayılan sayfanın yerine geçer.
	///
	/// Örnek:
	///   Pages: []page.Page{
	///       pages.NewDashboard(),    // Özel dashboard
	///       pages.NewSettings(),     // Özel settings
	///       pages.NewReportsPage(),  // Yeni özel sayfa
	///   }
	Pages []page.Page

	/// UserResource, kullanıcı resource'unu tutar
	UserResource resource.Resource

	/// Resources, panele eklenen ek resource'ları tutar
	/// Örnek: Ürünler, Kategoriler, Siparişler vb.
	Resources []resource.Resource

	/// CORS, Cross-Origin Resource Sharing yapılandırmasını tutar
	CORS CORSConfig

	/// EncryptionKey, cookie şifreleme için kullanılan AES anahtarıdır.
	/// SECURITY: 32-byte (AES-256) base64-encoded key kullanın
	/// Örnek: openssl rand -base64 32
	/// CRITICAL: Key'i her startup'ta değiştirmeyin (mevcut cookie'ler okunamaz hale gelir)
	/// Production'da environment variable'dan alın: os.Getenv("COOKIE_ENCRYPTION_KEY")
	EncryptionCookie encryptcookie.Config

	/// EnableHTTP2Push, HTTP/2 Server Push özelliğini aktif eder.
	/// PERFORMANCE: Kritik kaynakları (JS, CSS) proaktif olarak tarayıcıya gönderir
	/// IMPORTANT: Yanlış kullanımda performans kötüleşebilir, mutlaka ölçüm yapın
	/// REQUIRES: HTTPS ve HTTP/2 aktif olmalı
	/// Default: false (disabled)
	EnableHTTP2Push bool

	/// HTTP2PushResources, HTTP/2 Server Push ile gönderilecek kaynakların listesidir.
	/// Sadece kritik kaynakları (main JS/CSS bundle) ekleyin
	/// Örnek: []string{"/assets/app.js", "/assets/main.css"}
	/// IMPORTANT: Asset isimleri build hash'li ise her build'de güncellemelisiniz
	HTTP2PushResources []string

	/// CircuitBreaker, Circuit Breaker yapılandırmasını tutar
	/// Servis hatalarını yönetmek ve sistem çökmelerini önlemek için kullanılır
	CircuitBreaker CircuitBreakerConfig

	/// I18n, çoklu dil desteği yapılandırmasını tutar
	/// Uygulamanın farklı dillerde gösterilmesini sağlar
	I18n I18nConfig

	/// Plugins, plugin sistemi yapılandırmasını tutar
	/// Plugin'lerin otomatik keşfi ve yüklenmesi için kullanılır
	Plugins PluginConfig
}

// / # SettingsConfig - Dinamik Ayarlar Yapılandırması
// /
// / Veritabanından gelen ve runtime'da değiştirilebilen dinamik ayarları tutar.
// / Panel ayarlarının merkezi yönetim noktasıdır.
// /
// / ## Kullanım Senaryoları
// / - Site adını değiştirme
// / - Kayıt özelliğini açıp kapatma
// / - Şifre sıfırlama özelliğini kontrol etme
// / - Özel ayarları depolama
// /
// / ## Örnek Kullanım
// / ```go
// / settings := SettingsConfig{
// /     SiteName:       "Yönetim Paneli",
// /     Register:       true,
// /     ForgotPassword: true,
// /     Values: map[string]interface{}{
// /         "theme":           "dark",
// /         "items_per_page":  25,
// /         "session_timeout": 3600,
// /     },
// / }
// / ```
// /
// / ## Avantajlar
// / - Dinamik ayarlar veritabanında saklanır
// / - Runtime'da değişiklikler yapılabilir
// / - Özel ayarlar için Values map'i kullanılabilir
// / - Farklı kiracılar için farklı ayarlar
// /
// / ## Önemli Notlar
// / - Values map'i type assertion gerektirir
// / - Veritabanı bağlantısı gereklidir
// / - Ayarlar cache'lenebilir (performans için)
type SettingsConfig struct {
	/// SiteName, panelin adını tutar
	/// Örnek: "Yönetim Paneli", "Admin Dashboard"
	SiteName string

	/// Register, kullanıcı kaydı özelliğinin etkin olup olmadığını tutar
	/// Veritabanından gelen değer, FeatureConfig'i geçersiz kılar
	Register bool

	/// ForgotPassword, şifre sıfırlama özelliğinin etkin olup olmadığını tutar
	/// Veritabanından gelen değer, FeatureConfig'i geçersiz kılar
	ForgotPassword bool

	/// Values, özel ayarları tutar
	/// Örnek: map[string]interface{}{
	///     "theme": "dark",
	///     "language": "tr",
	///     "items_per_page": 25,
	/// }
	Values map[string]interface{}
}

// / # StorageConfig - Dosya Depolama Yapılandırması
// /
// / Dosya yükleme ve depolama ayarlarını tutar.
// / Kullanıcılar tarafından yüklenen dosyaların nereye kaydedileceğini belirtir.
// /
// / ## Kullanım Senaryoları
// / - Profil fotoğrafları yükleme
// / - Belge ve raporlar depolama
// / - Ürün görselleri yükleme
// / - Medya dosyaları yönetimi
// /
// / ## Örnek Kullanım
// / ```go
// / storageConfig := StorageConfig{
// /     Path: "/var/www/panel/storage/uploads",
// /     URL:  "/uploads",
// / }
// /
// / // Dosya yolu: /var/www/panel/storage/uploads/profile.jpg
// / // Erişim URL: https://example.com/uploads/profile.jpg
// / ```
// /
// / ## Avantajlar
// / - Merkezi dosya yönetimi
// / - Fiziksel yol ve URL'yi ayrı tutma
// / - Farklı depolama sistemlerine geçiş kolaylığı
// /
// / ## Önemli Notlar
// / - Path mutlak yol olmalıdır
// / - URL, web sunucusu tarafından erişilebilir olmalıdır
// / - Dosya izinleri doğru ayarlanmalıdır (755 veya 775)
// / - Düzenli olarak eski dosyaları temizleyin
// /
// / ## Güvenlik Uyarıları
// / - Yüklenen dosyaları doğrulayın (tip, boyut)
// / - Yürütülebilir dosyaları yüklemeyi engelle
// / - Dosya adlarını sanitize edin
// / - Antivirus taraması yapın
type StorageConfig struct {
	/// Path, dosyaların fiziksel olarak saklanacağı sunucu yoludur.
	/// Örnek: "/var/www/panel/storage/uploads"
	/// UYARI: Dizin yazılabilir olmalıdır (chmod 755 veya 775)
	Path string

	/// URL, dosyalara erişim için kullanılacak temel URL yoludur.
	/// Örnek: "/uploads" veya "https://cdn.example.com/uploads"
	/// Tarayıcıda erişim: https://example.com/uploads/filename.jpg
	URL string
}

// / # PermissionConfig - İzin Yönetimi Yapılandırması
// /
// / Rol tabanlı erişim kontrolü (RBAC) ayarlarını tutar.
// / Kullanıcıların hangi işlemleri yapabileceğini belirtir.
// /
// / ## Kullanım Senaryoları
// / - Rol tanımlama (Admin, Editor, Viewer)
// / - Resource'lara erişim kontrolü
// / - İşlemlere izin verme (Create, Read, Update, Delete)
// / - Dinamik izin yönetimi
// /
// / ## Örnek Kullanım
// / ```go
// / permConfig := PermissionConfig{
// /     Path: "/etc/panel/permissions.toml",
// / }
// /
// / // permissions.toml içeriği:
// / // [admin]
// / // resources = ["*"]
// / // actions = ["create", "read", "update", "delete"]
// / //
// / // [editor]
// / // resources = ["posts", "pages"]
// / // actions = ["create", "read", "update"]
// / ```
// /
// / ## Avantajlar
// / - Merkezi izin yönetimi
// / - Dinamik rol tanımlama
// / - Granüler kontrol
// /
// / ## Önemli Notlar
// / - permissions.toml dosyası başlangıçta yüklenmelidir
// / - Dosya değişiklikleri runtime'da yeniden yüklenebilir
// / - Varsayılan izinler tanımlanmalıdır
type PermissionConfig struct {
	/// Path, permissions.toml dosyasının yoludur.
	/// Örnek: "/etc/panel/permissions.toml"
	/// UYARI: Dosya okunabilir olmalıdır
	Path string
}

// / # ServerConfig - HTTP Sunucu Yapılandırması
// /
// / HTTP sunucusunun temel ayarlarını tutar.
// / Sunucunun hangi port ve host'ta çalışacağını belirtir.
// /
// / ## Kullanım Senaryoları
// / - Sunucuyu belirli bir port'ta başlatma
// / - Belirli bir IP adresine bağlama
// / - Proxy arkasında çalıştırma
// / - Docker container'ında çalıştırma
// /
// / ## Örnek Kullanım
// / ```go
// / serverConfig := ServerConfig{
// /     Port: "8080",
// /     Host: "0.0.0.0",  // Tüm interface'lerde dinle
// / }
// /
// / // Erişim: http://localhost:8080
// / // veya: http://192.168.1.100:8080
// / ```
// /
// / ## Avantajlar
// / - Basit ve anlaşılır yapı
// / - Farklı ortamlar için kolay konfigürasyon
// / - Dinamik port atama
// /
// / ## Önemli Notlar
// / - Port 1-65535 arasında olmalıdır
// / - 1-1024 arası portlar root yetkisi gerektirir
// / - Host "0.0.0.0" tüm interface'lerde dinler
// / - Host "127.0.0.1" sadece localhost'ta dinler
// /
// / ## Güvenlik Uyarıları
// / - Üretim ortamında "0.0.0.0" kullanmayın (firewall kuralları ekleyin)
// / - HTTPS kullanın (TLS/SSL)
// / - Proxy arkasında çalıştırırken X-Forwarded-* header'larını kontrol edin
type ServerConfig struct {
	/// Port, HTTP sunucusunun dinleyeceği port numarasıdır.
	/// Örnek: "8080", "3000", "5000"
	/// UYARI: 1-1024 arası portlar root yetkisi gerektirir
	Port string

	/// Host, HTTP sunucusunun bağlanacağı host adresini belirtir.
	/// Örnek: "0.0.0.0" (tüm interface'ler), "127.0.0.1" (localhost)
	/// UYARI: Üretim ortamında "0.0.0.0" kullanırken firewall kuralları ekleyin
	Host string
}

// / # ShardingConfig - Veritabanı Sharding Yapılandırması
// /
// / GORM Sharding plugin için yapılandırma ayarlarını tutar.
// / Büyük tablolar için horizontal partitioning (sharding) desteği sağlar.
// /
// / ## Kullanım Senaryoları
// / - Büyük veri setlerini birden fazla tabloya bölme
// / - Performans optimizasyonu için veri dağıtımı
// / - Yüksek trafikli uygulamalarda ölçeklenebilirlik
// /
// / ## Örnek Kullanım
// / ```go
// / shardingConfig := ShardingConfig{
// /     Enabled:             true,
// /     ShardingKey:         "user_id",
// /     NumberOfShards:      4,
// /     PrimaryKeyGenerator: "snowflake",
// / }
// / ```
// /
// / ## Önemli Notlar
// / - ShardingKey, veri dağıtımı için kullanılan kolon adıdır
// / - NumberOfShards, oluşturulacak shard sayısını belirtir
// / - PrimaryKeyGenerator: "snowflake", "postgresql", "mysql" veya "custom"
// / - Sharding etkinleştirildiğinde migration stratejisi değişir
// /
// / ## Güvenlik ve Performans
// / - Sharding key dikkatli seçilmelidir (değişmez olmalı)
// / - Shard sayısı başlangıçta iyi planlanmalıdır (sonradan değiştirmek zor)
// / - Cross-shard query'ler performansı etkileyebilir
type ShardingConfig struct {
	/// Enabled, sharding özelliğinin aktif olup olmadığını belirtir
	Enabled bool

	/// ShardingKey, veri dağıtımı için kullanılan kolon adı
	/// Örnek: "user_id", "tenant_id", "organization_id"
	ShardingKey string

	/// NumberOfShards, oluşturulacak shard (tablo bölümü) sayısı
	/// Örnek: 4, 8, 16 (genellikle 2'nin katları kullanılır)
	NumberOfShards uint

	/// PrimaryKeyGenerator, primary key üretim stratejisi
	/// Değerler: "snowflake", "postgresql", "mysql", "custom"
	PrimaryKeyGenerator string
}

// / # DatabaseConfig - Veritabanı Bağlantı Yapılandırması
// /
// / Veritabanı bağlantı bilgilerini tutar.
// / GORM ORM kütüphanesi tarafından kullanılır.
// /
// / ## Desteklenen Veritabanları
// / - PostgreSQL: "postgres"
// / - MySQL: "mysql"
// / - SQLite: "sqlite"
// /
// / ## Kullanım Senaryoları
// / - Veritabanı bağlantısı kurma
// / - Farklı ortamlar için farklı veritabanları
// / - Mevcut GORM bağlantısını kullanma
// / - Sharding ile horizontal partitioning
// /
// / ## Örnek Kullanım
// / ```go
// / // MySQL
// / dbConfig := DatabaseConfig{
// /     DSN:    "user:password@tcp(localhost:3306)/panel?charset=utf8mb4&parseTime=True",
// /     Driver: "mysql",
// / }
// /
// / // PostgreSQL
// / dbConfig := DatabaseConfig{
// /     DSN:    "host=localhost user=postgres password=secret dbname=panel port=5432",
// /     Driver: "postgres",
// / }
// /
// / // SQLite
// / dbConfig := DatabaseConfig{
// /     DSN:    "panel.db",
// /     Driver: "sqlite",
// / }
// /
// / // Mevcut GORM bağlantısını kullanma
// / dbConfig := DatabaseConfig{
// /     Instance: existingDB,
// / }
// /
// / // Sharding ile kullanım
// / dbConfig := DatabaseConfig{
// /     DSN:    "host=localhost user=postgres password=secret dbname=panel port=5432",
// /     Driver: "postgres",
// /     Sharding: ShardingConfig{
// /         Enabled:             true,
// /         ShardingKey:         "user_id",
// /         NumberOfShards:      4,
// /         PrimaryKeyGenerator: "snowflake",
// /     },
// / }
// / ```
// /
// / ## Avantajlar
// / - Çoklu veritabanı desteği
// / - Mevcut bağlantıyı yeniden kullanma
// / - GORM ile entegrasyon
// / - Sharding desteği ile ölçeklenebilirlik
// /
// / ## Önemli Notlar
// / - DSN formatı veritabanı türüne göre değişir
// / - Instance nil ise DSN ve Driver kullanılır
// / - Bağlantı havuzu ayarları DSN'de yapılabilir
// / - Zaman dilimi ayarlarını DSN'ye ekleyin
// / - Sharding etkinleştirildiğinde performans ve migration stratejisi değişir
// /
// / ## Güvenlik Uyarıları
// / - Veritabanı şifrelerini kaynak kodunda saklamayın
// / - Çevre değişkenlerinden yükleyin
// / - SSL/TLS bağlantısı kullanın
// / - Veritabanı kullanıcısı için minimum izinler verin
// / - Düzenli olarak yedek alın
type DatabaseConfig struct {
	/// DSN (Data Source Name), veritabanı bağlantı dizesidir.
	/// Format veritabanı türüne göre değişir:
	/// - MySQL: "user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True"
	/// - PostgreSQL: "host=localhost user=postgres password=secret dbname=panel port=5432"
	/// - SQLite: "panel.db" veya "/path/to/panel.db"
	DSN string

	/// Driver, kullanılacak veritabanı sürücüsünü belirtir.
	/// Değerler: "postgres", "mysql", "sqlite"
	Driver string

	/// Instance, mevcut bir GORM bağlantısıdır (Opsiyonel).
	/// Eğer Instance nil değilse, DSN ve Driver yoksayılır.
	/// Kullanım: Mevcut bir veritabanı bağlantısını yeniden kullanmak
	Instance *gorm.DB

	/// Sharding, veritabanı sharding yapılandırmasıdır (Opsiyonel).
	/// Büyük tablolar için horizontal partitioning desteği sağlar.
	/// Kullanım: Yüksek trafikli uygulamalarda performans optimizasyonu
	Sharding ShardingConfig
}

// / # CircuitBreakerConfig - Circuit Breaker Yapılandırması
// /
// / Circuit Breaker, servis hatalarını yönetmek ve sistem çökmelerini önlemek için
// / kullanılan bir dayanıklılık (resilience) desenidir. Üç durum arasında geçiş yapar:
// / Closed (Normal), Open (Devre Dışı), Half-Open (Test).
// /
// / ## Kullanım Senaryoları
// / - Dış API çağrılarını koruma
// / - Veritabanı bağlantı hatalarını yönetme
// / - Yavaş yanıt veren servisleri izole etme
// / - Cascade failure'ları önleme
// /
// / ## Örnek Kullanım
// / ```go
// / cbConfig := CircuitBreakerConfig{
// /     Enabled:                true,
// /     FailureThreshold:       5,      // 5 ardışık hata sonrası devre aç
// /     Timeout:                10 * time.Second,  // 10 saniye bekle
// /     SuccessThreshold:       5,      // 5 başarılı istek sonrası devre kapat
// /     HalfOpenMaxConcurrent:  1,      // Half-open'da 1 eşzamanlı istek
// / }
// / ```
// /
// / ## Circuit Breaker Durumları
// / - **Closed (Kapalı)**: Normal çalışma, istekler geçer, hatalar sayılır
// / - **Open (Açık)**: Devre açık, istekler hemen reddedilir (503 Service Unavailable)
// / - **Half-Open (Yarı Açık)**: Test modu, sınırlı sayıda istek geçer
// /
// / ## Avantajlar
// / - Sistem çökmelerini önler
// / - Hızlı hata yanıtı (fail-fast)
// / - Otomatik kurtarma (self-healing)
// / - Cascade failure'ları engeller
// /
// / ## Önemli Notlar
// / - FailureThreshold: Kaç ardışık hata sonrası devre açılır
// / - Timeout: Devre açıldıktan sonra ne kadar beklenecek
// / - SuccessThreshold: Kaç başarılı istek sonrası devre kapanır
// / - HalfOpenMaxConcurrent: Half-open durumunda kaç eşzamanlı istek
// /
// / ## Best Practices
// / - Kritik servislere uygulayın (dış API'ler, veritabanı)
// / - Timeout değerini servis yanıt süresine göre ayarlayın
// / - Monitoring ve alerting ekleyin
// / - Fallback mekanizmaları tanımlayın
type CircuitBreakerConfig struct {
	/// Enabled, Circuit Breaker'ın aktif olup olmadığını belirtir
	/// true: Circuit Breaker etkin, false: Devre dışı
	Enabled bool

	/// FailureThreshold, devre açılmadan önce kaç ardışık hata olması gerektiğini belirtir
	/// Varsayılan: 5
	/// Örnek: 5 ardışık hata sonrası devre açılır
	FailureThreshold int

	/// Timeout, devre açıldıktan sonra ne kadar süre bekleneceğini belirtir
	/// Varsayılan: 10 * time.Second
	/// Bu süre sonunda Half-Open durumuna geçilir
	Timeout time.Duration

	/// SuccessThreshold, devre kapanmadan önce kaç başarılı istek olması gerektiğini belirtir
	/// Varsayılan: 5
	/// Half-Open durumunda bu kadar başarılı istek sonrası devre kapanır
	SuccessThreshold int

	/// HalfOpenMaxConcurrent, Half-Open durumunda kaç eşzamanlı istek izin verileceğini belirtir
	/// Varsayılan: 1
	/// Genellikle 1 olarak ayarlanır (tek bir test isteği)
	HalfOpenMaxConcurrent int
}

// / # I18nConfig - Çoklu Dil Desteği Yapılandırması
// /
// / Panel uygulamasında çoklu dil desteği sağlar. go-i18n kütüphanesi kullanılarak
// / mesajların farklı dillerde gösterilmesini sağlar.
// /
// / ## Kullanım Senaryoları
// / - Çok dilli kullanıcı arayüzü
// / - Uluslararası uygulamalar
// / - Bölgesel içerik sunumu
// / - Dinamik dil değiştirme
// /
// / ## Örnek Kullanım
// / ```go
// / i18nConfig := I18nConfig{
// /     Enabled:          true,
// /     RootPath:         "./locales",
// /     AcceptLanguages:  []language.Tag{language.Turkish, language.English},
// /     DefaultLanguage:  language.Turkish,
// /     FormatBundleFile: "yaml",
// / }
// / ```
// /
// / ## Dil Dosyası Yapısı
// / ```yaml
// / # locales/tr.yaml
// / welcome:
// /   other: "Hoş geldiniz"
// / welcomeWithName:
// /   other: "Hoş geldiniz, {{.Name}}"
// / ```
// /
// / ## Dil Seçimi
// / Dil, şu sırayla belirlenir:
// / 1. Query parametresi: ?lang=tr
// / 2. Accept-Language header
// / 3. DefaultLanguage (fallback)
// /
// / ## Avantajlar
// / - Kolay çoklu dil desteği
// / - Dinamik dil değiştirme
// / - Fallback dil desteği
// / - Template değişkenleri desteği
// /
// / ## Önemli Notlar
// / - RootPath: Dil dosyalarının bulunduğu dizin
// / - AcceptLanguages: Desteklenen diller listesi
// / - DefaultLanguage: Varsayılan dil (fallback)
// / - FormatBundleFile: Dil dosyası formatı (yaml, json, toml)
// /
// / ## Best Practices
// / - Dil dosyalarını organize edin (locales/tr.yaml, locales/en.yaml)
// / - Fallback dil her zaman tanımlayın
// / - Template değişkenlerini kullanın ({{.Name}})
// / - Çeviri anahtarlarını anlamlı isimlendirin
type I18nConfig struct {
	/// Enabled, i18n'in aktif olup olmadığını belirtir
	/// true: Çoklu dil desteği etkin, false: Devre dışı
	Enabled bool

	/// RootPath, dil dosyalarının bulunduğu dizini belirtir
	/// Varsayılan: "./locales"
	/// Örnek: "./locales" -> locales/tr.yaml, locales/en.yaml
	RootPath string

	/// AcceptLanguages, desteklenen dillerin listesidir
	/// Varsayılan: []language.Tag{language.Turkish, language.English}
	/// Örnek: []language.Tag{language.Turkish, language.English, language.German}
	AcceptLanguages []language.Tag

	/// DefaultLanguage, varsayılan dili belirtir (fallback)
	/// Varsayılan: language.Turkish
	/// Eğer kullanıcının dili desteklenmiyorsa bu dil kullanılır
	DefaultLanguage language.Tag

	/// FormatBundleFile, dil dosyası formatını belirtir
	/// Varsayılan: "yaml"
	/// Değerler: "yaml", "json", "toml"
	FormatBundleFile string

	/// UseURLPrefix, URL'lerde dil prefix'i kullanılıp kullanılmayacağını belirtir
	/// true: URL'ler dil prefix'i içerir (örn: /api/en/resource/users)
	/// false: URL'ler dil prefix'i içermez (örn: /api/resource/users)
	/// Varsayılan: false
	UseURLPrefix bool

	/// URLPrefixOptional, varsayılan dil için URL prefix'inin opsiyonel olup olmadığını belirtir
	/// true: Varsayılan dil için prefix yok (örn: /api/resource/users), diğer diller için var (örn: /api/en/resource/users)
	/// false: Tüm diller için prefix var (örn: /api/tr/resource/users, /api/en/resource/users)
	/// Varsayılan: true
	/// Not: Bu ayar sadece UseURLPrefix=true olduğunda etkilidir
	URLPrefixOptional bool
}

// / # PluginConfig - Plugin Sistemi Yapılandırması
// /
// / Plugin sisteminin yapılandırma ayarlarını tutar.
// / Plugin'lerin otomatik keşfi ve yüklenmesi için kullanılır.
// /
// / ## Kullanım Senaryoları
// / - Plugin'leri otomatik keşfetme
// / - Plugin klasörünü belirleme
// / - Plugin yükleme stratejisi
// /
// / ## Örnek Kullanım
// / ```go
// / pluginConfig := PluginConfig{
// /     AutoDiscover: true,
// /     Path:         "./plugins",
// / }
// / ```
// /
// / ## Plugin Yükleme Stratejileri
// / 1. **Manuel Import (Önerilen)**: Plugin'ler compile-time'da import edilir
// /    ```go
// /    import _ "github.com/user/my-plugin"
// /    ```
// / 2. **Auto-discovery (Sınırlı/Opsiyonel)**: Sadece plugin descriptor'ları (plugin.yaml) keşfedilir
// /    ```go
// /    config.Plugins.AutoDiscover = true
// /    config.Plugins.Path = "./plugins"
// /    ```
// /
// / ## Avantajlar
// / - Manuel import: Type-safe, compile-time kontrol
// / - Auto-discovery: Descriptor görünürlüğü ve keşif
// /
// / ## Önemli Notlar
// / - AutoDiscover varsayılan olarak false'tur
// / - Manuel import tercih edilir (type-safe)
// / - Auto-discovery opsiyoneldir, runtime plugin package import/boot yapmaz
type PluginConfig struct {
	/// AutoDiscover, plugin'lerin otomatik keşfedilip keşfedilmeyeceğini belirtir
	/// true: Plugin'ler otomatik keşfedilir, false: Manuel import gereklidir
	/// Varsayılan: false (manuel import önerilir)
	AutoDiscover bool

	/// Path, plugin'lerin bulunduğu klasör yoludur
	/// Varsayılan: "./plugins"
	/// Örnek: "./plugins", "/etc/panel/plugins"
	/// UYARI: AutoDiscover true olmalıdır
	Path string
}
