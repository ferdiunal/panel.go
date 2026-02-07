package panel

import (
	"github.com/ferdiunal/panel.go/pkg/page"
	"github.com/ferdiunal/panel.go/pkg/resource"
	"gorm.io/gorm"
)

/// # FeatureConfig - Panel Özellik Yapılandırması
///
/// Panel uygulamasının opsiyonel özelliklerini açıp kapatmak için kullanılan yapıdır.
/// Bu yapı, panelin hangi özelliklerin etkin olacağını kontrol etmeyi sağlar.
///
/// ## Kullanım Senaryoları
/// - Kullanıcı kaydını devre dışı bırakmak (sadece admin tarafından oluşturulan hesaplar)
/// - Şifre sıfırlama özelliğini kapatmak (LDAP/SSO entegrasyonunda)
/// - Belirli ortamlarda (test, staging) özellikleri kontrol etmek
///
/// ## Örnek Kullanım
/// ```go
/// config := FeatureConfig{
///     Register:       true,        // Kullanıcı kaydını etkinleştir
///     ForgotPassword: true,        // Şifre sıfırlama özelliğini etkinleştir
/// }
/// ```
///
/// ## Avantajlar
/// - Basit boolean yapısı ile kolay kontrol
/// - Runtime'da özellikleri değiştirmek mümkün
/// - Farklı ortamlar için farklı konfigürasyonlar oluşturulabilir
///
/// ## Önemli Notlar
/// - Varsayılan değer false'tur (özellikleri açıkça etkinleştirmeniz gerekir)
/// - Veritabanı ayarlarından da kontrol edilebilir (SettingsConfig)
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

/// # OAuthConfig - OAuth Sağlayıcı Yapılandırması
///
/// Panelde kullanılacak OAuth sağlayıcılarının yapılandırmasını tutar.
/// Şu anda sadece Google OAuth desteklenmektedir.
///
/// ## Kullanım Senaryoları
/// - Google hesabı ile giriş sağlamak
/// - Sosyal medya entegrasyonu
/// - Kurumsal SSO çözümleri
///
/// ## Örnek Kullanım
/// ```go
/// oauthConfig := OAuthConfig{
///     Google: GoogleConfig{
///         ClientID:     "your-client-id.apps.googleusercontent.com",
///         ClientSecret: "your-client-secret",
///         RedirectURL:  "https://example.com/auth/google/callback",
///     },
/// }
/// ```
///
/// ## Gelecek Genişletmeler
/// - GitHub OAuth desteği
/// - Microsoft OAuth desteği
/// - OIDC (OpenID Connect) desteği
type OAuthConfig struct {
	/// Google, Google OAuth sağlayıcısının yapılandırmasını içerir
	Google GoogleConfig
}

/// # GoogleConfig - Google OAuth Yapılandırması
///
/// Google OAuth 2.0 entegrasyonu için gerekli yapılandırma bilgilerini tutar.
/// Google Cloud Console'dan alınan kimlik bilgilerini içerir.
///
/// ## Kullanım Senaryoları
/// - Google hesabı ile giriş
/// - Google Workspace entegrasyonu
/// - Kurumsal ortamlarda SSO
///
/// ## Örnek Kullanım
/// ```go
/// googleConfig := GoogleConfig{
///     ClientID:     "123456789-abcdefghijklmnopqrstuvwxyz.apps.googleusercontent.com",
///     ClientSecret: "GOCSPX-1234567890abcdefghijklmno",
///     RedirectURL:  "https://panel.example.com/auth/google/callback",
/// }
///
/// if googleConfig.Enabled() {
///     // Google OAuth'u etkinleştir
/// }
/// ```
///
/// ## Önemli Notlar
/// - ClientSecret asla client-side'da açığa çıkarılmamalıdır
/// - RedirectURL, Google Cloud Console'da kayıtlı olmalıdır
/// - Üretim ortamında HTTPS kullanılmalıdır
/// - ClientID ve ClientSecret çevre değişkenlerinden yüklenmelidir
///
/// ## Güvenlik Uyarıları
/// - Kimlik bilgilerini kaynak kodunda saklamayın
/// - .env dosyasını .gitignore'a ekleyin
/// - Düzenli olarak token'ları döndürün
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

/// # Enabled - Google OAuth Yapılandırmasının Geçerliliğini Kontrol Eder
///
/// Google OAuth'un etkinleştirilip etkinleştirilmediğini kontrol eder.
/// ClientID ve ClientSecret'ın boş olmadığını doğrular.
///
/// ## Dönüş Değeri
/// - true: Google OAuth yapılandırması geçerli ve etkindir
/// - false: Google OAuth yapılandırması eksik veya geçersizdir
///
/// ## Kullanım Senaryoları
/// - Giriş sayfasında Google butonu gösterilip gösterilmeyeceğini belirlemek
/// - OAuth middleware'inde yapılandırma kontrolü
/// - Başlangıç sırasında yapılandırma doğrulaması
///
/// ## Örnek Kullanım
/// ```go
/// config := GoogleConfig{
///     ClientID:     "123456789-abc.apps.googleusercontent.com",
///     ClientSecret: "GOCSPX-1234567890abc",
///     RedirectURL:  "https://example.com/callback",
/// }
///
/// if config.Enabled() {
///     // Google OAuth'u etkinleştir
///     setupGoogleOAuth(config)
/// } else {
///     // Google OAuth'u devre dışı bırak
///     log.Warn("Google OAuth yapılandırması eksik")
/// }
/// ```
///
/// ## Önemli Notlar
/// - RedirectURL boş olsa bile true döndürebilir (opsiyonel kontrol)
/// - ClientID ve ClientSecret zorunludur
/// - Üretim ortamında bu fonksiyon başlangıçta çağrılmalıdır
func (c GoogleConfig) Enabled() bool {
	return c.ClientID != "" && c.ClientSecret != ""
}

/// # CORSConfig - Cross-Origin Resource Sharing Yapılandırması
///
/// Panelin farklı origin'lerden gelen istekleri kabul etmesini kontrol eder.
/// Web uygulamaları arasında veri paylaşımını güvenli bir şekilde sağlar.
///
/// ## Kullanım Senaryoları
/// - Frontend ve backend farklı domain'lerde çalışıyor
/// - Mobil uygulamalar panel API'sini kullanıyor
/// - Microservices mimarisinde API erişimi
/// - Geliştirme ortamında localhost'tan erişim
///
/// ## Örnek Kullanım
/// ```go
/// // Üretim ortamı
/// corsConfig := CORSConfig{
///     AllowedOrigins: []string{
///         "https://app.example.com",
///         "https://admin.example.com",
///     },
/// }
///
/// // Geliştirme ortamı
/// devCorsConfig := CORSConfig{
///     AllowedOrigins: []string{
///         "http://localhost:3000",
///         "http://localhost:5173",
///         "http://127.0.0.1:8080",
///     },
/// }
/// ```
///
/// ## Avantajlar
/// - Güvenli cross-origin istekleri sağlar
/// - Farklı ortamlar için farklı konfigürasyonlar
/// - Dinamik olarak değiştirilebilir
///
/// ## Güvenlik Uyarıları
/// - AllowedOrigins'e "*" (wildcard) eklemeyin (güvenlik riski)
/// - Sadece güvenilir origin'leri ekleyin
/// - Üretim ortamında localhost'u eklemeyin
/// - HTTPS kullanın
///
/// ## Önemli Notlar
/// - Boş liste tüm origin'leri reddeder
/// - Origin kontrolü case-sensitive'dir
/// - Port numarası origin'in bir parçasıdır
type CORSConfig struct {
	/// AllowedOrigins, CORS isteklerine izin verilen origin'lerin listesidir.
	/// Örnek: []string{"https://example.com", "https://app.example.com"}
	/// Geliştirme için: []string{"http://localhost:3000", "http://localhost:5173"}
	///
	/// UYARI: Wildcard "*" kullanmayın, güvenlik riski oluşturur!
	AllowedOrigins []string
}

/// # Config - Panel Genel Yapılandırması
///
/// Panelin tüm yapılandırma ayarlarını merkezi olarak tutar.
/// Sunucu, veritabanı, OAuth, depolama ve diğer tüm ayarları içerir.
///
/// ## Yapı Bileşenleri
/// - Server: HTTP sunucu ayarları (port, host)
/// - Database: Veritabanı bağlantı bilgileri
/// - Environment: Çalışma ortamı (production, development, test)
/// - Features: Etkinleştirilmiş özellikler
/// - OAuth: OAuth sağlayıcı yapılandırması
/// - Storage: Dosya depolama ayarları
/// - Permissions: İzin yönetimi yapılandırması
/// - SettingsValues: Veritabanından gelen dinamik ayarlar
/// - SettingsPage: Ayarlar sayfası yapılandırması
/// - DashboardPage: Dashboard sayfası yapılandırması
/// - UserResource: Kullanıcı resource'u
/// - Resources: Ek resource'lar
/// - CORS: CORS yapılandırması
///
/// ## Kullanım Senaryoları
/// - Panel başlangıcında yapılandırma yükleme
/// - Farklı ortamlar için farklı konfigürasyonlar
/// - Runtime'da ayarları değiştirme
/// - Yapılandırma doğrulaması
///
/// ## Örnek Kullanım
/// ```go
/// config := Config{
///     Server: ServerConfig{
///         Port: "8080",
///         Host: "0.0.0.0",
///     },
///     Database: DatabaseConfig{
///         DSN:    "user:password@tcp(localhost:3306)/panel",
///         Driver: "mysql",
///     },
///     Environment: "production",
///     Features: FeatureConfig{
///         Register:       true,
///         ForgotPassword: true,
///     },
///     CORS: CORSConfig{
///         AllowedOrigins: []string{"https://example.com"},
///     },
/// }
/// ```
///
/// ## Önemli Notlar
/// - Config yapısı thread-safe değildir, başlangıçta ayarlanmalıdır
/// - Veritabanı bağlantısı Instance alanında saklanır
/// - Resources dinamik olarak eklenebilir
/// - SettingsPage ve DashboardPage nil olabilir
///
/// ## Güvenlik Uyarıları
/// - Kimlik bilgilerini kaynak kodunda saklamayın
/// - Çevre değişkenlerinden yükleyin
/// - Üretim ortamında debug modu kapatın
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

	/// SettingsPage, ayarlar sayfasının yapılandırmasını tutar (opsiyonel)
	SettingsPage *page.Settings

	/// DashboardPage, dashboard sayfasının yapılandırmasını tutar (opsiyonel)
	DashboardPage *page.Dashboard

	/// UserResource, kullanıcı resource'unu tutar
	UserResource resource.Resource

	/// Resources, panele eklenen ek resource'ları tutar
	/// Örnek: Ürünler, Kategoriler, Siparişler vb.
	Resources []resource.Resource

	/// CORS, Cross-Origin Resource Sharing yapılandırmasını tutar
	CORS CORSConfig
}

/// # SettingsConfig - Dinamik Ayarlar Yapılandırması
///
/// Veritabanından gelen ve runtime'da değiştirilebilen dinamik ayarları tutar.
/// Panel ayarlarının merkezi yönetim noktasıdır.
///
/// ## Kullanım Senaryoları
/// - Site adını değiştirme
/// - Kayıt özelliğini açıp kapatma
/// - Şifre sıfırlama özelliğini kontrol etme
/// - Özel ayarları depolama
///
/// ## Örnek Kullanım
/// ```go
/// settings := SettingsConfig{
///     SiteName:       "Yönetim Paneli",
///     Register:       true,
///     ForgotPassword: true,
///     Values: map[string]interface{}{
///         "theme":           "dark",
///         "items_per_page":  25,
///         "session_timeout": 3600,
///     },
/// }
/// ```
///
/// ## Avantajlar
/// - Dinamik ayarlar veritabanında saklanır
/// - Runtime'da değişiklikler yapılabilir
/// - Özel ayarlar için Values map'i kullanılabilir
/// - Farklı kiracılar için farklı ayarlar
///
/// ## Önemli Notlar
/// - Values map'i type assertion gerektirir
/// - Veritabanı bağlantısı gereklidir
/// - Ayarlar cache'lenebilir (performans için)
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

/// # StorageConfig - Dosya Depolama Yapılandırması
///
/// Dosya yükleme ve depolama ayarlarını tutar.
/// Kullanıcılar tarafından yüklenen dosyaların nereye kaydedileceğini belirtir.
///
/// ## Kullanım Senaryoları
/// - Profil fotoğrafları yükleme
/// - Belge ve raporlar depolama
/// - Ürün görselleri yükleme
/// - Medya dosyaları yönetimi
///
/// ## Örnek Kullanım
/// ```go
/// storageConfig := StorageConfig{
///     Path: "/var/www/panel/storage/uploads",
///     URL:  "/uploads",
/// }
///
/// // Dosya yolu: /var/www/panel/storage/uploads/profile.jpg
/// // Erişim URL: https://example.com/uploads/profile.jpg
/// ```
///
/// ## Avantajlar
/// - Merkezi dosya yönetimi
/// - Fiziksel yol ve URL'yi ayrı tutma
/// - Farklı depolama sistemlerine geçiş kolaylığı
///
/// ## Önemli Notlar
/// - Path mutlak yol olmalıdır
/// - URL, web sunucusu tarafından erişilebilir olmalıdır
/// - Dosya izinleri doğru ayarlanmalıdır (755 veya 775)
/// - Düzenli olarak eski dosyaları temizleyin
///
/// ## Güvenlik Uyarıları
/// - Yüklenen dosyaları doğrulayın (tip, boyut)
/// - Yürütülebilir dosyaları yüklemeyi engelle
/// - Dosya adlarını sanitize edin
/// - Antivirus taraması yapın
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

/// # PermissionConfig - İzin Yönetimi Yapılandırması
///
/// Rol tabanlı erişim kontrolü (RBAC) ayarlarını tutar.
/// Kullanıcıların hangi işlemleri yapabileceğini belirtir.
///
/// ## Kullanım Senaryoları
/// - Rol tanımlama (Admin, Editor, Viewer)
/// - Resource'lara erişim kontrolü
/// - İşlemlere izin verme (Create, Read, Update, Delete)
/// - Dinamik izin yönetimi
///
/// ## Örnek Kullanım
/// ```go
/// permConfig := PermissionConfig{
///     Path: "/etc/panel/permissions.toml",
/// }
///
/// // permissions.toml içeriği:
/// // [admin]
/// // resources = ["*"]
/// // actions = ["create", "read", "update", "delete"]
/// //
/// // [editor]
/// // resources = ["posts", "pages"]
/// // actions = ["create", "read", "update"]
/// ```
///
/// ## Avantajlar
/// - Merkezi izin yönetimi
/// - Dinamik rol tanımlama
/// - Granüler kontrol
///
/// ## Önemli Notlar
/// - permissions.toml dosyası başlangıçta yüklenmelidir
/// - Dosya değişiklikleri runtime'da yeniden yüklenebilir
/// - Varsayılan izinler tanımlanmalıdır
type PermissionConfig struct {
	/// Path, permissions.toml dosyasının yoludur.
	/// Örnek: "/etc/panel/permissions.toml"
	/// UYARI: Dosya okunabilir olmalıdır
	Path string
}

/// # ServerConfig - HTTP Sunucu Yapılandırması
///
/// HTTP sunucusunun temel ayarlarını tutar.
/// Sunucunun hangi port ve host'ta çalışacağını belirtir.
///
/// ## Kullanım Senaryoları
/// - Sunucuyu belirli bir port'ta başlatma
/// - Belirli bir IP adresine bağlama
/// - Proxy arkasında çalıştırma
/// - Docker container'ında çalıştırma
///
/// ## Örnek Kullanım
/// ```go
/// serverConfig := ServerConfig{
///     Port: "8080",
///     Host: "0.0.0.0",  // Tüm interface'lerde dinle
/// }
///
/// // Erişim: http://localhost:8080
/// // veya: http://192.168.1.100:8080
/// ```
///
/// ## Avantajlar
/// - Basit ve anlaşılır yapı
/// - Farklı ortamlar için kolay konfigürasyon
/// - Dinamik port atama
///
/// ## Önemli Notlar
/// - Port 1-65535 arasında olmalıdır
/// - 1-1024 arası portlar root yetkisi gerektirir
/// - Host "0.0.0.0" tüm interface'lerde dinler
/// - Host "127.0.0.1" sadece localhost'ta dinler
///
/// ## Güvenlik Uyarıları
/// - Üretim ortamında "0.0.0.0" kullanmayın (firewall kuralları ekleyin)
/// - HTTPS kullanın (TLS/SSL)
/// - Proxy arkasında çalıştırırken X-Forwarded-* header'larını kontrol edin
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

/// # ShardingConfig - Veritabanı Sharding Yapılandırması
///
/// GORM Sharding plugin için yapılandırma ayarlarını tutar.
/// Büyük tablolar için horizontal partitioning (sharding) desteği sağlar.
///
/// ## Kullanım Senaryoları
/// - Büyük veri setlerini birden fazla tabloya bölme
/// - Performans optimizasyonu için veri dağıtımı
/// - Yüksek trafikli uygulamalarda ölçeklenebilirlik
///
/// ## Örnek Kullanım
/// ```go
/// shardingConfig := ShardingConfig{
///     Enabled:             true,
///     ShardingKey:         "user_id",
///     NumberOfShards:      4,
///     PrimaryKeyGenerator: "snowflake",
/// }
/// ```
///
/// ## Önemli Notlar
/// - ShardingKey, veri dağıtımı için kullanılan kolon adıdır
/// - NumberOfShards, oluşturulacak shard sayısını belirtir
/// - PrimaryKeyGenerator: "snowflake", "postgresql", "mysql" veya "custom"
/// - Sharding etkinleştirildiğinde migration stratejisi değişir
///
/// ## Güvenlik ve Performans
/// - Sharding key dikkatli seçilmelidir (değişmez olmalı)
/// - Shard sayısı başlangıçta iyi planlanmalıdır (sonradan değiştirmek zor)
/// - Cross-shard query'ler performansı etkileyebilir
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

/// # DatabaseConfig - Veritabanı Bağlantı Yapılandırması
///
/// Veritabanı bağlantı bilgilerini tutar.
/// GORM ORM kütüphanesi tarafından kullanılır.
///
/// ## Desteklenen Veritabanları
/// - PostgreSQL: "postgres"
/// - MySQL: "mysql"
/// - SQLite: "sqlite"
///
/// ## Kullanım Senaryoları
/// - Veritabanı bağlantısı kurma
/// - Farklı ortamlar için farklı veritabanları
/// - Mevcut GORM bağlantısını kullanma
/// - Sharding ile horizontal partitioning
///
/// ## Örnek Kullanım
/// ```go
/// // MySQL
/// dbConfig := DatabaseConfig{
///     DSN:    "user:password@tcp(localhost:3306)/panel?charset=utf8mb4&parseTime=True",
///     Driver: "mysql",
/// }
///
/// // PostgreSQL
/// dbConfig := DatabaseConfig{
///     DSN:    "host=localhost user=postgres password=secret dbname=panel port=5432",
///     Driver: "postgres",
/// }
///
/// // SQLite
/// dbConfig := DatabaseConfig{
///     DSN:    "panel.db",
///     Driver: "sqlite",
/// }
///
/// // Mevcut GORM bağlantısını kullanma
/// dbConfig := DatabaseConfig{
///     Instance: existingDB,
/// }
///
/// // Sharding ile kullanım
/// dbConfig := DatabaseConfig{
///     DSN:    "host=localhost user=postgres password=secret dbname=panel port=5432",
///     Driver: "postgres",
///     Sharding: ShardingConfig{
///         Enabled:             true,
///         ShardingKey:         "user_id",
///         NumberOfShards:      4,
///         PrimaryKeyGenerator: "snowflake",
///     },
/// }
/// ```
///
/// ## Avantajlar
/// - Çoklu veritabanı desteği
/// - Mevcut bağlantıyı yeniden kullanma
/// - GORM ile entegrasyon
/// - Sharding desteği ile ölçeklenebilirlik
///
/// ## Önemli Notlar
/// - DSN formatı veritabanı türüne göre değişir
/// - Instance nil ise DSN ve Driver kullanılır
/// - Bağlantı havuzu ayarları DSN'de yapılabilir
/// - Zaman dilimi ayarlarını DSN'ye ekleyin
/// - Sharding etkinleştirildiğinde performans ve migration stratejisi değişir
///
/// ## Güvenlik Uyarıları
/// - Veritabanı şifrelerini kaynak kodunda saklamayın
/// - Çevre değişkenlerinden yükleyin
/// - SSL/TLS bağlantısı kullanın
/// - Veritabanı kullanıcısı için minimum izinler verin
/// - Düzenli olarak yedek alın
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
