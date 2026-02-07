package config

import "time"

/// # SecurityConfig
///
/// Bu yapı, uygulamanın tüm güvenlik ile ilgili yapılandırmalarını merkezi bir noktada toplar.
/// Web uygulamalarında güvenlik katmanlarını (CORS, rate limiting, session, encryption, audit)
/// tek bir yapıda yönetmeyi sağlar.
///
/// ## Kullanım Senaryoları
///
/// - **Çok Katmanlı Güvenlik**: Farklı güvenlik mekanizmalarını bir arada kullanma
/// - **Ortam Bazlı Yapılandırma**: Development, staging ve production için farklı güvenlik seviyeleri
/// - **Merkezi Yönetim**: Tüm güvenlik ayarlarını tek bir yerden kontrol etme
/// - **Compliance**: Güvenlik standartlarına (OWASP, PCI-DSS) uyum sağlama
///
/// ## Örnek Kullanım
///
/// ```go
/// // Varsayılan güvenlik yapılandırması
/// config := config.DefaultSecurityConfig()
///
/// // Production ortamı için
/// prodConfig := config.ProductionSecurityConfig()
///
/// // Özel yapılandırma
/// customConfig := config.SecurityConfig{
///     CORS: config.CORSConfig{
///         AllowedOrigins: []string{"https://example.com"},
///         AllowCredentials: true,
///     },
///     RateLimit: config.RateLimitConfig{
///         Enabled: true,
///         AuthMaxRequests: 5,
///         AuthWindow: 1 * time.Minute,
///     },
/// }
/// ```
///
/// ## Avantajlar
///
/// - **Modüler Yapı**: Her güvenlik bileşeni bağımsız olarak yapılandırılabilir
/// - **Tip Güvenliği**: Compile-time'da hataları yakalar
/// - **Varsayılan Değerler**: Güvenli varsayılan ayarlarla gelir
/// - **Esneklik**: Ortam ve ihtiyaca göre özelleştirilebilir
///
/// ## Önemli Notlar
///
/// ⚠️ **Güvenlik Uyarıları**:
/// - Production ortamında mutlaka `ProductionSecurityConfig()` kullanın
/// - CORS AllowedOrigins'i asla "*" olarak ayarlamayın
/// - Encryption key'lerini environment variable'lardan okuyun
/// - Audit logging'i production'da mutlaka aktif tutun
///
/// ⚠️ **Performans Notları**:
/// - Rate limiting çok sıkı ayarlanırsa meşru kullanıcıları etkileyebilir
/// - Audit logging disk I/O'sunu artırır, log rotation yapılandırın
///
/// ## İlgili Yapılar
///
/// - `CORSConfig`: Cross-Origin Resource Sharing yapılandırması
/// - `RateLimitConfig`: İstek hızı sınırlama yapılandırması
/// - `AccountLockoutConfig`: Hesap kilitleme yapılandırması
/// - `SessionConfig`: Oturum güvenliği yapılandırması
/// - `EncryptionConfig`: Veri şifreleme yapılandırması
/// - `AuditConfig`: Denetim günlüğü yapılandırması
type SecurityConfig struct {
	/// CORS (Cross-Origin Resource Sharing) yapılandırması.
	/// Tarayıcıların farklı origin'lerden gelen istekleri nasıl işleyeceğini belirler.
	CORS CORSConfig

	/// Rate limiting (istek hızı sınırlama) yapılandırması.
	/// Brute-force ve DDoS saldırılarına karşı koruma sağlar.
	RateLimit RateLimitConfig

	/// Hesap kilitleme yapılandırması.
	/// Başarısız giriş denemelerinden sonra hesapları otomatik olarak kilitler.
	AccountLockout AccountLockoutConfig

	/// Oturum (session) güvenlik yapılandırması.
	/// Cookie ayarları ve oturum yönetimi parametrelerini içerir.
	Session SessionConfig

	/// Veri şifreleme yapılandırması.
	/// Hassas verilerin şifrelenmesi için kullanılan algoritma ve key yönetimi.
	Encryption EncryptionConfig

	/// Denetim günlüğü (audit log) yapılandırması.
	/// Güvenlik olaylarının ve kullanıcı aktivitelerinin kaydedilmesi.
	Audit AuditConfig
}

/// # CORSConfig
///
/// Bu yapı, Cross-Origin Resource Sharing (CORS) yapılandırmasını yönetir.
/// CORS, tarayıcıların farklı origin'lerden (domain, protocol, port) gelen
/// HTTP isteklerini güvenli bir şekilde işlemesini sağlayan bir güvenlik mekanizmasıdır.
///
/// ## Kullanım Senaryoları
///
/// - **SPA (Single Page Application)**: Frontend ve backend farklı domain'lerde
/// - **Mikroservis Mimarisi**: Farklı servisler arası iletişim
/// - **API Gateway**: Çoklu client'lardan gelen istekleri yönetme
/// - **CDN Entegrasyonu**: Statik içerik farklı domain'den sunulduğunda
/// - **Mobil Uygulama API'leri**: Native app'lerden API çağrıları
///
/// ## Örnek Kullanım
///
/// ```go
/// // Production için güvenli CORS yapılandırması
/// corsConfig := config.CORSConfig{
///     AllowedOrigins: []string{
///         "https://app.example.com",
///         "https://admin.example.com",
///     },
///     AllowCredentials: true,
///     AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
///     AllowedHeaders: []string{"Content-Type", "Authorization"},
///     ExposeHeaders: []string{"X-Total-Count"},
///     MaxAge: 3600, // 1 saat
/// }
///
/// // Development için esnek yapılandırma
/// devCorsConfig := config.CORSConfig{
///     AllowedOrigins: []string{
///         "http://localhost:3000",
///         "http://localhost:5173",
///     },
///     AllowCredentials: true,
///     AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
///     AllowedHeaders: []string{"*"},
///     MaxAge: 86400,
/// }
/// ```
///
/// ## Avantajlar
///
/// - **Güvenlik**: Sadece izin verilen origin'lerden gelen istekleri kabul eder
/// - **Esneklik**: Method, header ve credential kontrolü sağlar
/// - **Performans**: Preflight cache ile gereksiz OPTIONS isteklerini azaltır
/// - **Standart Uyumluluk**: W3C CORS standardına tam uyumlu
///
/// ## Önemli Notlar
///
/// ⚠️ **KRİTİK GÜVENLİK UYARILARI**:
/// - **ASLA** production'da `AllowedOrigins: []string{"*"}` kullanmayın!
/// - `AllowCredentials: true` iken wildcard (*) origin kullanılamaz
/// - HTTPS kullanmayan origin'lere production'da izin vermeyin
/// - Subdomain wildcard'ları dikkatli kullanın (*.example.com)
///
/// ⚠️ **Performans Notları**:
/// - `MaxAge` değerini yüksek tutarak preflight isteklerini azaltın
/// - Gereksiz header'ları `AllowedHeaders`'a eklemeyin
/// - `ExposeHeaders` sadece gerekli header'ları içermeli
///
/// ⚠️ **Yaygın Hatalar**:
/// - Preflight (OPTIONS) isteklerini handle etmeyi unutmak
/// - `AllowCredentials` ile wildcard origin birlikte kullanmak
/// - Port numaralarını origin'de belirtmeyi unutmak
/// - Protocol (http/https) uyumsuzluğu
///
/// ## CORS İstek Akışı
///
/// 1. **Simple Request**: Doğrudan istek gönderilir
///    - GET, HEAD, POST (belirli content-type'lar)
///    - Basit header'lar
///
/// 2. **Preflight Request**: Önce OPTIONS isteği gönderilir
///    - PUT, DELETE, PATCH gibi methodlar
///    - Custom header'lar
///    - Özel content-type'lar
///
/// ## İlgili Kaynaklar
///
/// - MDN CORS: https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS
/// - W3C Spec: https://www.w3.org/TR/cors/
type CORSConfig struct {
	/// İzin verilen origin'lerin listesi.
	/// Her origin tam URL formatında olmalıdır (protocol + domain + port).
	///
	/// **Örnekler**:
	/// - `["https://example.com"]` - Tek origin
	/// - `["https://app.example.com", "https://admin.example.com"]` - Çoklu origin
	/// - `["http://localhost:3000"]` - Development için
	///
	/// **UYARI**: Production'da ASLA "*" kullanmayın!
	AllowedOrigins []string

	/// Credential'ların (cookie, authorization header, TLS certificate)
	/// cross-origin isteklerde gönderilmesine izin verilip verilmeyeceğini belirtir.
	///
	/// **true**: Cookie ve auth header'lar gönderilir (güvenli, kimlik doğrulama gerekli)
	/// **false**: Credential'lar gönderilmez (public API'ler için)
	///
	/// **UYARI**: `true` iken AllowedOrigins "*" olamaz!
	AllowCredentials bool

	/// İzin verilen HTTP method'larının listesi.
	///
	/// **Standart methodlar**:
	/// - GET, POST, PUT, DELETE, PATCH, OPTIONS, HEAD
	///
	/// **Not**: OPTIONS her zaman preflight için gereklidir.
	AllowedMethods []string

	/// Client'ın gönderebileceği HTTP header'larının listesi.
	///
	/// **Yaygın header'lar**:
	/// - "Content-Type" - Request body tipi
	/// - "Authorization" - Auth token
	/// - "X-CSRF-Token" - CSRF koruması
	/// - "Accept" - Response format
	///
	/// **Not**: Basit header'lar (Accept, Accept-Language, Content-Language)
	/// otomatik olarak izin verilir.
	AllowedHeaders []string

	/// Client'ın JavaScript'ten erişebileceği response header'larının listesi.
	///
	/// **Yaygın kullanımlar**:
	/// - "X-Total-Count" - Pagination için toplam kayıt sayısı
	/// - "X-RateLimit-Remaining" - Kalan rate limit
	/// - "Content-Length" - Response boyutu
	///
	/// **Not**: Basit response header'lar (Cache-Control, Content-Language,
	/// Content-Type, Expires, Last-Modified, Pragma) otomatik olarak expose edilir.
	ExposeHeaders []string

	/// Preflight (OPTIONS) isteğinin sonucunun tarayıcı tarafından
	/// ne kadar süre cache'lenebileceğini saniye cinsinden belirtir.
	///
	/// **Önerilen değerler**:
	/// - Development: 3600 (1 saat)
	/// - Production: 86400 (24 saat)
	/// - Yüksek trafik: 7200 (2 saat) - daha sık güncelleme için
	///
	/// **Avantaj**: Preflight isteklerini azaltarak performansı artırır.
	MaxAge int
}

/// # RateLimitConfig
///
/// Bu yapı, istek hızı sınırlama (rate limiting) yapılandırmasını yönetir.
/// Rate limiting, belirli bir zaman diliminde yapılabilecek istek sayısını
/// sınırlayarak brute-force saldırıları, DDoS saldırıları ve API kötüye kullanımını önler.
///
/// ## Kullanım Senaryoları
///
/// - **Brute-Force Koruması**: Login endpoint'lerinde deneme yanılma saldırılarını engelleme
/// - **DDoS Önleme**: Aşırı istek yükünü sınırlayarak sunucu kaynaklarını koruma
/// - **API Kota Yönetimi**: Kullanıcı başına API kullanım limitlerini zorunlu kılma
/// - **Fair Usage**: Tüm kullanıcıların adil kaynak kullanımını sağlama
/// - **Cost Control**: Ücretli API'lerde maliyet kontrolü
///
/// ## Örnek Kullanım
///
/// ```go
/// // Sıkı güvenlik için auth rate limiting
/// authRateLimit := config.RateLimitConfig{
///     Enabled: true,
///     AuthMaxRequests: 5,           // 5 deneme
///     AuthWindow: 15 * time.Minute, // 15 dakikada
///     APIMaxRequests: 100,          // 100 istek
///     APIWindow: 1 * time.Minute,   // 1 dakikada
/// }
///
/// // Yüksek trafikli API için
/// apiRateLimit := config.RateLimitConfig{
///     Enabled: true,
///     AuthMaxRequests: 10,
///     AuthWindow: 5 * time.Minute,
///     APIMaxRequests: 1000,         // Daha yüksek limit
///     APIWindow: 1 * time.Minute,
/// }
///
/// // Development ortamı için esnek ayarlar
/// devRateLimit := config.RateLimitConfig{
///     Enabled: true,
///     AuthMaxRequests: 50,          // Test için yüksek
///     AuthWindow: 1 * time.Minute,
///     APIMaxRequests: 500,
///     APIWindow: 1 * time.Minute,
/// }
/// ```
///
/// ## Avantajlar
///
/// - **Güvenlik**: Otomatik saldırıları yavaşlatır veya durdurur
/// - **Kaynak Koruması**: Sunucu kaynaklarının tükenmesini önler
/// - **Adil Kullanım**: Tek bir kullanıcının tüm kaynakları tüketmesini engeller
/// - **Maliyet Kontrolü**: Aşırı kullanımdan kaynaklanan maliyetleri sınırlar
///
/// ## Dezavantajlar
///
/// - **False Positive**: Meşru kullanıcılar yanlışlıkla engellenebilir
/// - **Shared IP**: NAT arkasındaki kullanıcılar aynı limiti paylaşır
/// - **Performans**: Her istek için rate limit kontrolü yapılır
/// - **Karmaşıklık**: Distributed sistemlerde senkronizasyon gerekir
///
/// ## Önemli Notlar
///
/// ⚠️ **Yapılandırma Önerileri**:
/// - **Auth Endpoint'leri**: Çok sıkı limitler (5-10 istek/15 dakika)
/// - **Public API'ler**: Orta seviye limitler (100-1000 istek/dakika)
/// - **Internal API'ler**: Esnek limitler veya devre dışı
/// - **Admin Panel**: Daha yüksek limitler ama yine de korumalı
///
/// ⚠️ **Güvenlik Uyarıları**:
/// - Rate limiting tek başına yeterli değildir, diğer güvenlik önlemleriyle birlikte kullanın
/// - IP-based rate limiting NAT/proxy arkasında sorun yaratabilir
/// - User-based rate limiting için authentication gereklidir
/// - Distributed cache (Redis) kullanarak cluster'lar arası senkronizasyon sağlayın
///
/// ⚠️ **Performans Notları**:
/// - In-memory rate limiting tek sunucuda çalışır
/// - Redis/Memcached kullanarak distributed rate limiting yapın
/// - Rate limit kontrolü her istekte çalışır, optimize edin
/// - Sliding window algoritması daha adil ama daha maliyetlidir
///
/// ## Rate Limiting Stratejileri
///
/// 1. **Fixed Window**: Sabit zaman dilimlerinde sayaç (basit, hızlı)
/// 2. **Sliding Window**: Kayan pencere (daha adil, daha karmaşık)
/// 3. **Token Bucket**: Token bazlı (burst'lere izin verir)
/// 4. **Leaky Bucket**: Sabit hızda işleme (düzgün trafik)
///
/// ## HTTP Response Kodları
///
/// - **429 Too Many Requests**: Rate limit aşıldığında
/// - **Retry-After**: Ne zaman tekrar deneyebileceğini belirtir
/// - **X-RateLimit-Limit**: Toplam limit
/// - **X-RateLimit-Remaining**: Kalan istek sayısı
/// - **X-RateLimit-Reset**: Limitin sıfırlanacağı zaman
///
/// ## İlgili Kaynaklar
///
/// - OWASP Rate Limiting: https://owasp.org/www-community/controls/Blocking_Brute_Force_Attacks
/// - RFC 6585: https://tools.ietf.org/html/rfc6585
type RateLimitConfig struct {
	/// Rate limiting'in aktif olup olmadığını belirtir.
	///
	/// **true**: Rate limiting aktif (production için önerilir)
	/// **false**: Rate limiting devre dışı (sadece development için)
	///
	/// **UYARI**: Production ortamında mutlaka true olmalıdır!
	Enabled bool

	/// Authentication endpoint'leri için maksimum istek sayısı.
	/// Bu limit, login, register, password reset gibi kimlik doğrulama
	/// işlemleri için geçerlidir.
	///
	/// **Önerilen değerler**:
	/// - Production: 5-10 (çok sıkı)
	/// - Staging: 20-30 (orta)
	/// - Development: 50-100 (esnek)
	///
	/// **Örnek**: 5 başarısız login denemesinden sonra engelle
	AuthMaxRequests int

	/// Authentication rate limiting için zaman penceresi.
	/// AuthMaxRequests sayısı bu süre içinde aşılırsa istek reddedilir.
	///
	/// **Önerilen değerler**:
	/// - Yüksek güvenlik: 15 * time.Minute (15 dakika)
	/// - Orta güvenlik: 5 * time.Minute (5 dakika)
	/// - Düşük güvenlik: 1 * time.Minute (1 dakika)
	///
	/// **Örnek**: 15 dakika içinde 5 deneme
	AuthWindow time.Duration

	/// Genel API endpoint'leri için maksimum istek sayısı.
	/// Normal API çağrıları (CRUD işlemleri, veri sorgulama) için geçerlidir.
	///
	/// **Önerilen değerler**:
	/// - Düşük trafik: 100-500 istek/dakika
	/// - Orta trafik: 500-2000 istek/dakika
	/// - Yüksek trafik: 2000-10000 istek/dakika
	///
	/// **Not**: Kullanıcı sayısı ve uygulama tipine göre ayarlayın
	APIMaxRequests int

	/// API rate limiting için zaman penceresi.
	/// APIMaxRequests sayısı bu süre içinde aşılırsa istek reddedilir.
	///
	/// **Önerilen değerler**:
	/// - Standart: 1 * time.Minute (1 dakika)
	/// - Esnek: 5 * time.Minute (5 dakika)
	/// - Sıkı: 10 * time.Second (10 saniye)
	///
	/// **Örnek**: 1 dakika içinde 100 istek
	APIWindow time.Duration
}

/// # AccountLockoutConfig
///
/// Bu yapı, hesap kilitleme (account lockout) yapılandırmasını yönetir.
/// Başarısız giriş denemelerinden sonra kullanıcı hesaplarını otomatik olarak
/// kilitleyerek brute-force saldırılarına karşı ek bir koruma katmanı sağlar.
///
/// ## Kullanım Senaryoları
///
/// - **Brute-Force Koruması**: Şifre tahmin saldırılarını önleme
/// - **Credential Stuffing Engelleme**: Çalıntı kimlik bilgilerinin denenmesini durdurma
/// - **Hesap Güvenliği**: Şüpheli aktivitelerde hesabı koruma altına alma
/// - **Compliance**: PCI-DSS, HIPAA gibi standartların gereksinimlerini karşılama
/// - **Otomatik Tehdit Yanıtı**: Manuel müdahale gerektirmeden güvenlik önlemi alma
///
/// ## Örnek Kullanım
///
/// ```go
/// // Yüksek güvenlikli ortam için sıkı ayarlar
/// strictLockout := config.AccountLockoutConfig{
///     Enabled: true,
///     MaxAttempts: 3,                // 3 başarısız deneme
///     LockoutDuration: 30 * time.Minute, // 30 dakika kilitle
/// }
///
/// // Standart güvenlik ayarları
/// standardLockout := config.AccountLockoutConfig{
///     Enabled: true,
///     MaxAttempts: 5,                // 5 başarısız deneme
///     LockoutDuration: 15 * time.Minute, // 15 dakika kilitle
/// }
///
/// // Esnek ayarlar (kullanıcı dostu)
/// lenientLockout := config.AccountLockoutConfig{
///     Enabled: true,
///     MaxAttempts: 10,               // 10 başarısız deneme
///     LockoutDuration: 5 * time.Minute,  // 5 dakika kilitle
/// }
///
/// // Kalıcı kilitleme (manuel açılana kadar)
/// permanentLockout := config.AccountLockoutConfig{
///     Enabled: true,
///     MaxAttempts: 5,
///     LockoutDuration: 0, // 0 = kalıcı kilitleme
/// }
/// ```
///
/// ## Avantajlar
///
/// - **Otomatik Koruma**: Manuel müdahale gerektirmeden hesapları korur
/// - **Saldırı Yavaşlatma**: Brute-force saldırılarını ekonomik olmaktan çıkarır
/// - **Kullanıcı Farkındalığı**: Şüpheli aktivite konusunda kullanıcıyı uyarır
/// - **Audit Trail**: Başarısız denemeler loglanarak analiz edilebilir
/// - **Compliance**: Güvenlik standartlarının gereksinimlerini karşılar
///
/// ## Dezavantajlar
///
/// - **Denial of Service**: Saldırganlar meşru kullanıcıları kilitleyebilir
/// - **Kullanıcı Deneyimi**: Şifresini unutan kullanıcılar engellenebilir
/// - **Destek Yükü**: Kilitlenen hesaplar için destek talepleri artar
/// - **Shared Accounts**: Paylaşılan hesaplarda sorun yaratabilir
///
/// ## Önemli Notlar
///
/// ⚠️ **Güvenlik Önerileri**:
/// - **MaxAttempts**: 3-5 arası değer önerilir (çok düşük UX sorununa, çok yüksek güvenlik açığına yol açar)
/// - **LockoutDuration**: 15-30 dakika arası optimal (çok kısa etkisiz, çok uzun kullanıcı dostu değil)
/// - **Kalıcı Kilitleme**: Sadece yüksek güvenlikli sistemlerde kullanın (LockoutDuration: 0)
/// - **Bildirim**: Kullanıcıya email/SMS ile kilitleme bildirimi gönderin
/// - **Admin Override**: Admin'lerin hesapları manuel açabilmesi gerekir
///
/// ⚠️ **DoS Saldırısı Önleme**:
/// - IP-based rate limiting ile birlikte kullanın
/// - CAPTCHA ekleyerek otomatik denemeleri zorlaştırın
/// - Hesap kilitleme öncesi CAPTCHA gösterin
/// - Şüpheli IP'lerden gelen istekleri daha sıkı kontrol edin
///
/// ⚠️ **Kullanıcı Deneyimi**:
/// - Kalan deneme hakkını kullanıcıya gösterin
/// - "Şifremi Unuttum" linkini belirgin yapın
/// - Kilitleme süresi dolmadan "Şifremi Unuttum" ile açılabilsin
/// - Başarılı girişte sayacı sıfırlayın
/// - Kilitleme nedeni ve süresini açıkça belirtin
///
/// ⚠️ **Yaygın Hatalar**:
/// - Başarılı girişte sayacı sıfırlamamak
/// - Kilitleme durumunu cache'lememek (her istekte DB sorgusu)
/// - Admin hesaplarını da kilitlemek (bypass mekanizması olmalı)
/// - Kilitleme loglarını tutmamak
/// - Email/SMS bildirimi göndermemek
///
/// ## Uygulama Detayları
///
/// ### Sayaç Yönetimi
/// ```go
/// // Başarısız deneme
/// failedAttempts++
/// if failedAttempts >= MaxAttempts {
///     lockUntil = time.Now().Add(LockoutDuration)
/// }
///
/// // Başarılı giriş
/// failedAttempts = 0
/// lockUntil = nil
/// ```
///
/// ### Kilitleme Kontrolü
/// ```go
/// if lockUntil != nil && time.Now().Before(*lockUntil) {
///     return ErrAccountLocked
/// }
/// ```
///
/// ### Otomatik Açılma
/// ```go
/// if lockUntil != nil && time.Now().After(*lockUntil) {
///     failedAttempts = 0
///     lockUntil = nil
/// }
/// ```
///
/// ## Güvenlik Standartları
///
/// - **OWASP**: Authentication Cheat Sheet önerilerine uygun
/// - **PCI-DSS**: Requirement 8.1.6 - 6 başarısız denemeden sonra kilitleme
/// - **NIST 800-63B**: Account lockout politikaları
/// - **HIPAA**: Access control gereksinimleri
///
/// ## İlgili Kaynaklar
///
/// - OWASP Authentication: https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html
/// - NIST 800-63B: https://pages.nist.gov/800-63-3/sp800-63b.html
type AccountLockoutConfig struct {
	/// Hesap kilitleme özelliğinin aktif olup olmadığını belirtir.
	///
	/// **true**: Hesap kilitleme aktif (production için şiddetle önerilir)
	/// **false**: Hesap kilitleme devre dışı (sadece development/test için)
	///
	/// **UYARI**: Production ortamında mutlaka true olmalıdır!
	/// Devre dışı bırakmak brute-force saldırılarına karşı savunmasız bırakır.
	Enabled bool

	/// Hesap kilitlenmeden önce izin verilen maksimum başarısız giriş denemesi sayısı.
	///
	/// **Önerilen değerler**:
	/// - Yüksek güvenlik: 3 deneme (bankacılık, sağlık)
	/// - Standart güvenlik: 5 deneme (genel uygulamalar)
	/// - Düşük güvenlik: 10 deneme (internal tool'lar)
	///
	/// **PCI-DSS**: Maksimum 6 deneme gerektirir
	///
	/// **Örnek**: 5 başarısız deneme sonrası hesap kilitlenir
	///
	/// **Not**: Çok düşük değer kullanıcı deneyimini olumsuz etkiler,
	/// çok yüksek değer güvenlik açığı yaratır.
	MaxAttempts int

	/// Hesabın ne kadar süre kilitli kalacağını belirtir.
	///
	/// **Önerilen değerler**:
	/// - Standart: 15 * time.Minute (15 dakika)
	/// - Sıkı: 30 * time.Minute (30 dakika)
	/// - Çok sıkı: 1 * time.Hour (1 saat)
	/// - Kalıcı: 0 (manuel açılana kadar, sadece kritik sistemler için)
	///
	/// **Örnek**: 15 dakika sonra hesap otomatik olarak açılır
	///
	/// **UYARI**: 0 değeri kalıcı kilitleme anlamına gelir!
	/// Bu durumda hesap sadece admin tarafından veya "Şifremi Unuttum"
	/// akışı ile açılabilir.
	///
	/// **Not**: Çok kısa süre etkisiz, çok uzun süre kullanıcı dostu değil.
	/// Optimal değer 15-30 dakika arasıdır.
	LockoutDuration time.Duration
}

/// # SessionConfig
///
/// Bu yapı, oturum (session) güvenlik yapılandırmasını yönetir.
/// HTTP cookie'leri üzerinden oturum yönetimi yapan uygulamalarda
/// güvenli cookie ayarlarını ve oturum parametrelerini tanımlar.
///
/// ## Kullanım Senaryoları
///
/// - **Kimlik Doğrulama**: Kullanıcı giriş durumunu takip etme
/// - **Session Management**: Güvenli oturum yaşam döngüsü yönetimi
/// - **CSRF Koruması**: Cross-Site Request Forgery saldırılarını önleme
/// - **XSS Koruması**: Cross-Site Scripting saldırılarını sınırlama
/// - **Cookie Security**: Güvenli cookie politikaları uygulama
///
/// ## Örnek Kullanım
///
/// ```go
/// // Production için maksimum güvenlik
/// prodSession := config.SessionConfig{
///     CookieName: "__Host-session_token",  // __Host- prefix güvenlik sağlar
///     Secure: true,                        // Sadece HTTPS
///     HTTPOnly: true,                      // JavaScript erişimi yok
///     SameSite: "Strict",                  // CSRF koruması
///     MaxAge: 86400,                       // 24 saat
///     Domain: "",                          // Mevcut domain
///     Path: "/",                           // Tüm path'ler
/// }
///
/// // Development için esnek ayarlar
/// devSession := config.SessionConfig{
///     CookieName: "session_token",
///     Secure: false,                       // HTTP'ye izin ver
///     HTTPOnly: true,
///     SameSite: "Lax",                     // Daha esnek
///     MaxAge: 86400,
///     Domain: "",
///     Path: "/",
/// }
///
/// // API için kısa ömürlü session
/// apiSession := config.SessionConfig{
///     CookieName: "__Secure-api_session",
///     Secure: true,
///     HTTPOnly: true,
///     SameSite: "Strict",
///     MaxAge: 3600,                        // 1 saat
///     Domain: "api.example.com",
///     Path: "/api",
/// }
///
/// // Remember Me özelliği için uzun ömürlü
/// rememberSession := config.SessionConfig{
///     CookieName: "__Host-remember_token",
///     Secure: true,
///     HTTPOnly: true,
///     SameSite: "Strict",
///     MaxAge: 2592000,                     // 30 gün
///     Domain: "",
///     Path: "/",
/// }
/// ```
///
/// ## Avantajlar
///
/// - **Güvenlik**: HTTPOnly ve Secure flag'leri ile saldırılara karşı koruma
/// - **CSRF Koruması**: SameSite attribute ile cross-site saldırıları önleme
/// - **Esneklik**: Domain ve Path ile cookie scope'unu kontrol etme
/// - **Standart Uyumluluk**: RFC 6265 cookie standardına tam uyum
/// - **Tarayıcı Desteği**: Modern tarayıcılarda tam destek
///
/// ## Dezavantajlar
///
/// - **HTTPS Gereksinimi**: Secure flag production'da HTTPS zorunlu kılar
/// - **SameSite Kısıtlamaları**: Strict mode bazı meşru cross-site senaryoları engelleyebilir
/// - **Cookie Boyutu**: Cookie'ler her istekte gönderilir, boyut önemli
/// - **Subdomain Sorunları**: Domain ayarları subdomain'lerde dikkat gerektirir
///
/// ## Önemli Notlar
///
/// ⚠️ **KRİTİK GÜVENLİK UYARILARI**:
/// - **ASLA** production'da `Secure: false` kullanmayın!
/// - **ASLA** `HTTPOnly: false` yapmayın (XSS saldırılarına açık olur)
/// - **Cookie Prefix'leri** kullanın: `__Host-` veya `__Secure-`
/// - **SameSite**: Production'da "Strict" veya "Lax" kullanın, "None" kullanmayın
/// - **MaxAge**: Çok uzun süre güvenlik riski, çok kısa süre kullanıcı deneyimi sorunu
///
/// ⚠️ **Cookie Prefix Kuralları**:
/// - `__Host-`: Secure=true, Domain boş, Path=/ gerektirir (en güvenli)
/// - `__Secure-`: Sadece Secure=true gerektirir
/// - Prefix kullanımı tarayıcı tarafından zorunlu kılınır
///
/// ⚠️ **SameSite Değerleri**:
/// - **Strict**: En güvenli, cross-site isteklerde cookie gönderilmez
///   - Kullanım: Hassas işlemler (banking, admin panel)
///   - Dezavantaj: External link'lerden gelenlerde session kaybolur
/// - **Lax**: Dengeli, top-level navigation'da cookie gönderilir
///   - Kullanım: Genel web uygulamaları
///   - Avantaj: Kullanıcı dostu, yine de CSRF koruması sağlar
/// - **None**: En esnek, tüm cross-site isteklerde cookie gönderilir
///   - Kullanım: Embedded widget'lar, iframe'ler
///   - UYARI: Secure=true zorunlu, CSRF riski yüksek
///
/// ⚠️ **Domain ve Path Ayarları**:
/// - **Domain boş**: Sadece mevcut domain (subdomain'ler hariç)
/// - **Domain=".example.com"**: Tüm subdomain'ler dahil
/// - **Path="/"**: Tüm path'ler
/// - **Path="/admin"**: Sadece /admin altındaki path'ler
///
/// ⚠️ **MaxAge Önerileri**:
/// - **Kısa oturum**: 3600 (1 saat) - yüksek güvenlik
/// - **Standart**: 86400 (24 saat) - dengeli
/// - **Remember Me**: 2592000 (30 gün) - kullanıcı dostu
/// - **Kalıcı**: 31536000 (1 yıl) - dikkatli kullanın
///
/// ⚠️ **Yaygın Hatalar**:
/// - HTTPS olmadan Secure=true kullanmak (cookie set edilmez)
/// - SameSite=None ile Secure=false kullanmak (geçersiz)
/// - __Host- prefix ile Domain set etmek (geçersiz)
/// - __Host- prefix ile Path!="/" kullanmak (geçersiz)
/// - Cookie name'de özel karakterler kullanmak
///
/// ## Session Güvenlik Kontrol Listesi
///
/// ✅ **Zorunlu Güvenlik Önlemleri**:
/// - [ ] Secure=true (production)
/// - [ ] HTTPOnly=true (her zaman)
/// - [ ] SameSite="Strict" veya "Lax"
/// - [ ] Cookie prefix kullanımı (__Host- veya __Secure-)
/// - [ ] Makul MaxAge değeri (24 saat önerilir)
/// - [ ] Session rotation (login sonrası yeni session ID)
/// - [ ] Session invalidation (logout'ta session silme)
/// - [ ] HTTPS zorunlu (production)
///
/// ✅ **Ek Güvenlik Önlemleri**:
/// - [ ] Session fixation koruması
/// - [ ] Concurrent session kontrolü
/// - [ ] Idle timeout mekanizması
/// - [ ] IP address validation
/// - [ ] User-Agent validation
/// - [ ] CSRF token ile birlikte kullanım
///
/// ## HTTP Header Örnekleri
///
/// ```http
/// // Production için ideal cookie
/// Set-Cookie: __Host-session_token=abc123; Secure; HttpOnly; SameSite=Strict; Path=/; Max-Age=86400
///
/// // Development için
/// Set-Cookie: session_token=abc123; HttpOnly; SameSite=Lax; Path=/; Max-Age=86400
///
/// // Subdomain'ler için
/// Set-Cookie: __Secure-session=abc123; Secure; HttpOnly; SameSite=Lax; Domain=.example.com; Path=/; Max-Age=86400
/// ```
///
/// ## Güvenlik Standartları
///
/// - **OWASP**: Session Management Cheat Sheet
/// - **RFC 6265**: HTTP State Management Mechanism
/// - **NIST 800-63B**: Digital Identity Guidelines
/// - **PCI-DSS**: Requirement 6.5.10 - Session management
///
/// ## İlgili Kaynaklar
///
/// - OWASP Session Management: https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet.html
/// - MDN Set-Cookie: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Set-Cookie
/// - RFC 6265: https://tools.ietf.org/html/rfc6265
type SessionConfig struct {
	/// Session cookie'sinin adı.
	///
	/// **Güvenlik Prefix'leri**:
	/// - `__Host-`: En güvenli, Secure=true + Domain="" + Path="/" gerektirir
	/// - `__Secure-`: Güvenli, sadece Secure=true gerektirir
	/// - Prefix yok: Standart cookie (daha az güvenli)
	///
	/// **Örnekler**:
	/// - `__Host-session_token` - Production için ideal
	/// - `__Secure-api_session` - API için
	/// - `session_token` - Development için
	///
	/// **UYARI**: Cookie name'de özel karakterler (;, =, space) kullanmayın!
	///
	/// **Not**: Prefix kullanımı tarayıcı tarafından zorunlu kılınır,
	/// kurallara uymazsanız cookie set edilmez.
	CookieName string

	/// Cookie'nin sadece HTTPS üzerinden gönderilip gönderilmeyeceğini belirtir.
	///
	/// **true**: Cookie sadece HTTPS isteklerde gönderilir (production için zorunlu)
	/// **false**: HTTP ve HTTPS'de gönderilir (sadece development için)
	///
	/// **UYARI**: Production ortamında MUTLAKA true olmalıdır!
	/// false değeri man-in-the-middle saldırılarına karşı savunmasız bırakır.
	///
	/// **Not**: __Host- veya __Secure- prefix kullanıyorsanız true zorunludur.
	Secure bool

	/// Cookie'ye JavaScript'ten erişilip erişilemeyeceğini belirtir.
	///
	/// **true**: JavaScript document.cookie ile erişemez (XSS koruması)
	/// **false**: JavaScript erişebilir (GÜVENLİK RİSKİ!)
	///
	/// **UYARI**: MUTLAKA true olmalıdır!
	/// false değeri XSS saldırılarında session çalınmasına yol açar.
	///
	/// **İstisna**: Client-side'da cookie'ye erişim gerekiyorsa
	/// (örn: mobile app token), ayrı bir cookie kullanın.
	HTTPOnly bool

	/// Cross-site isteklerde cookie'nin gönderilip gönderilmeyeceğini kontrol eder.
	/// CSRF (Cross-Site Request Forgery) saldırılarına karşı koruma sağlar.
	///
	/// **Değerler**:
	/// - `"Strict"`: En güvenli, hiçbir cross-site istekte gönderilmez
	///   - Kullanım: Banking, admin panel, hassas işlemler
	///   - Dezavantaj: External link'lerden gelenlerde session kaybolur
	///
	/// - `"Lax"`: Dengeli, sadece top-level navigation'da gönderilir (GET)
	///   - Kullanım: Genel web uygulamaları (önerilen)
	///   - Avantaj: Kullanıcı dostu + CSRF koruması
	///
	/// - `"None"`: Tüm cross-site isteklerde gönderilir
	///   - Kullanım: Embedded widget'lar, iframe'ler, third-party API'ler
	///   - UYARI: Secure=true zorunlu, CSRF riski yüksek
	///
	/// **Önerilen**: Production'da "Strict" veya "Lax"
	///
	/// **UYARI**: "None" kullanıyorsanız mutlaka Secure=true olmalı!
	SameSite string

	/// Cookie'nin maksimum yaşam süresi (saniye cinsinden).
	/// Bu süre sonunda cookie otomatik olarak silinir.
	///
	/// **Önerilen değerler**:
	/// - 3600 (1 saat): Yüksek güvenlik gerektiren uygulamalar
	/// - 86400 (24 saat): Standart web uygulamaları (önerilen)
	/// - 604800 (7 gün): Uzun oturum gerektiren uygulamalar
	/// - 2592000 (30 gün): "Remember Me" özelliği
	///
	/// **0 değeri**: Session cookie (tarayıcı kapanınca silinir)
	///
	/// **UYARI**: Çok uzun süre (>30 gün) güvenlik riski yaratır!
	/// Uzun oturumlar için ayrı "remember me" token kullanın.
	///
	/// **Not**: MaxAge ve Expires birlikte kullanılırsa MaxAge önceliklidir.
	MaxAge int

	/// Cookie'nin geçerli olduğu domain.
	///
	/// **Değerler**:
	/// - `""` (boş): Sadece mevcut domain (subdomain'ler hariç)
	///   - Örnek: example.com'da set edilirse sadece example.com'da geçerli
	///   - Kullanım: Tek domain uygulamaları (önerilen)
	///
	/// - `".example.com"`: Tüm subdomain'ler dahil
	///   - Örnek: example.com, api.example.com, admin.example.com
	///   - Kullanım: Multi-subdomain uygulamaları
	///   - UYARI: Güvenlik riski, dikkatli kullanın
	///
	/// - `"example.com"`: Sadece example.com (subdomain'siz)
	///
	/// **UYARI**: __Host- prefix kullanıyorsanız Domain boş olmalı!
	///
	/// **Güvenlik Notu**: Domain ne kadar geniş o kadar risk.
	/// Mümkün olduğunca dar tutun.
	Domain string

	/// Cookie'nin geçerli olduğu URL path.
	///
	/// **Değerler**:
	/// - `"/"`: Tüm path'ler (en yaygın)
	/// - `"/admin"`: Sadece /admin ve altındaki path'ler
	/// - `"/api/v1"`: Sadece /api/v1 ve altındaki path'ler
	///
	/// **UYARI**: __Host- prefix kullanıyorsanız Path="/" olmalı!
	///
	/// **Kullanım Senaryoları**:
	/// - Admin panel için ayrı cookie: Path="/admin"
	/// - API için ayrı cookie: Path="/api"
	/// - Genel uygulama: Path="/"
	///
	/// **Güvenlik Notu**: Path ne kadar dar o kadar güvenli.
	/// Ancak çoğu durumda "/" yeterlidir.
	Path string
}

/// # EncryptionConfig
///
/// Bu yapı, veri şifreleme yapılandırmasını yönetir.
/// Hassas verilerin (şifreler, kişisel bilgiler, ödeme bilgileri) güvenli bir şekilde
/// saklanması ve iletilmesi için kullanılan şifreleme algoritması ve key yönetimi ayarlarını içerir.
///
/// ## Kullanım Senaryoları
///
/// - **Veri Koruma**: Database'de hassas verilerin şifreli saklanması
/// - **PII Encryption**: Kişisel tanımlanabilir bilgilerin (PII) korunması
/// - **Compliance**: GDPR, HIPAA, PCI-DSS gereksinimlerini karşılama
/// - **Token Encryption**: API token'ları ve session ID'lerinin şifrelenmesi
/// - **File Encryption**: Yüklenen dosyaların şifreli saklanması
/// - **Backup Security**: Yedeklerin şifreli tutulması
///
/// ## Örnek Kullanım
///
/// ```go
/// // Production için AES-GCM (önerilen)
/// prodEncryption := config.EncryptionConfig{
///     Algorithm: "AES-GCM",
///     KeyHex: os.Getenv("ENCRYPTION_KEY"), // 32 byte (256-bit)
///     RotationEnabled: true,
///     RotationInterval: 90 * 24 * time.Hour, // 90 gün
/// }
///
/// // Yüksek güvenlik için key rotation
/// highSecEncryption := config.EncryptionConfig{
///     Algorithm: "AES-GCM",
///     KeyHex: loadKeyFromVault(), // Key vault'tan yükle
///     RotationEnabled: true,
///     RotationInterval: 30 * 24 * time.Hour, // 30 gün
/// }
///
/// // Legacy sistem için AES-CBC
/// legacyEncryption := config.EncryptionConfig{
///     Algorithm: "AES-CBC",
///     KeyHex: os.Getenv("ENCRYPTION_KEY"),
///     RotationEnabled: false,
/// }
///
/// // Development için (test key)
/// devEncryption := config.EncryptionConfig{
///     Algorithm: "AES-GCM",
///     KeyHex: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
///     RotationEnabled: false,
/// }
/// ```
///
/// ## Avantajlar
///
/// - **Veri Güvenliği**: Hassas verileri yetkisiz erişimden korur
/// - **Compliance**: Yasal gereksinimleri karşılar
/// - **Key Rotation**: Düzenli key değişimi ile güvenliği artırır
/// - **Authenticated Encryption**: AES-GCM ile veri bütünlüğü garantisi
/// - **Standart Algoritmalar**: Endüstri standardı şifreleme kullanır
///
/// ## Dezavantajlar
///
/// - **Performans**: Şifreleme/deşifreleme CPU kullanır
/// - **Key Yönetimi**: Key'lerin güvenli saklanması kritik
/// - **Karmaşıklık**: Key rotation ve migration karmaşık olabilir
/// - **Arama Zorluğu**: Şifreli verilerde full-text search yapılamaz
/// - **Backup**: Key kaybı veri kaybı anlamına gelir
///
/// ## Önemli Notlar
///
/// ⚠️ **KRİTİK GÜVENLİK UYARILARI**:
/// - **ASLA** encryption key'i kodda hardcode etmeyin!
/// - **ASLA** key'i git repository'ye commit etmeyin!
/// - **ASLA** key'i log'lara yazdırmayın!
/// - **MUTLAKA** environment variable veya key vault kullanın
/// - **MUTLAKA** key'leri güvenli bir yerde yedekleyin
/// - **MUTLAKA** production ve development için farklı key'ler kullanın
///
/// ⚠️ **Algoritma Seçimi**:
/// - **AES-GCM**: Modern, hızlı, authenticated encryption (ÖNERİLİR)
///   - Avantaj: Veri bütünlüğü kontrolü dahil
///   - Avantaj: Paralel işleme desteği
///   - Kullanım: Yeni projeler için ideal
///
/// - **AES-CBC**: Eski, yaygın, sadece confidentiality
///   - Dezavantaj: Ayrı MAC gerektirir (HMAC)
///   - Dezavantaj: Padding oracle saldırılarına açık olabilir
///   - Kullanım: Legacy sistemler için
///
/// ⚠️ **Key Boyutları**:
/// - **128-bit (16 byte)**: Minimum güvenlik, hızlı
/// - **192-bit (24 byte)**: Orta güvenlik
/// - **256-bit (32 byte)**: Maksimum güvenlik (ÖNERİLİR)
///
/// ⚠️ **Key Generation**:
/// ```go
/// // Güvenli key üretimi
/// key := make([]byte, 32) // 256-bit
/// if _, err := rand.Read(key); err != nil {
///     panic(err)
/// }
/// keyHex := hex.EncodeToString(key)
/// ```
///
/// ⚠️ **Key Storage Seçenekleri**:
/// 1. **Environment Variables**: Basit, yaygın
///    ```bash
///    export ENCRYPTION_KEY="abc123..."
///    ```
///
/// 2. **Key Vault**: En güvenli (AWS KMS, Azure Key Vault, HashiCorp Vault)
///    ```go
///    key, err := kmsClient.Decrypt(encryptedKey)
///    ```
///
/// 3. **Config File**: Şifreli config dosyası
///    ```go
///    key, err := loadEncryptedConfig("keys.enc")
///    ```
///
/// 4. **Hardware Security Module (HSM)**: En yüksek güvenlik
///
/// ⚠️ **Key Rotation Stratejisi**:
/// 1. **Yeni key oluştur**: Yeni şifreleme key'i üret
/// 2. **Dual encryption**: Geçiş süresince iki key'le çalış
/// 3. **Re-encryption**: Eski verileri yeni key ile şifrele
/// 4. **Old key retirement**: Eski key'i güvenli şekilde sil
/// 5. **Audit**: Tüm işlemleri logla
///
/// ⚠️ **Yaygın Hatalar**:
/// - Key'i kodda hardcode etmek
/// - Aynı key'i tüm ortamlarda kullanmak
/// - Key rotation yapmamak
/// - Key backup almamak
/// - IV (Initialization Vector) tekrar kullanmak
/// - Şifreleme hatalarını ignore etmek
/// - Key'i plain text olarak loglamak
///
/// ## Şifreleme Best Practices
///
/// ### 1. Key Management
/// ```go
/// // ✅ DOĞRU: Environment variable'dan oku
/// key := os.Getenv("ENCRYPTION_KEY")
/// if key == "" {
///     log.Fatal("ENCRYPTION_KEY not set")
/// }
///
/// // ❌ YANLIŞ: Hardcode
/// key := "my-secret-key-123"
/// ```
///
/// ### 2. Error Handling
/// ```go
/// // ✅ DOĞRU: Hataları handle et
/// encrypted, err := encrypt(data, key)
/// if err != nil {
///     log.Error("Encryption failed", "error", err)
///     return err
/// }
///
/// // ❌ YANLIŞ: Hataları ignore et
/// encrypted, _ := encrypt(data, key)
/// ```
///
/// ### 3. IV (Initialization Vector) Kullanımı
/// ```go
/// // ✅ DOĞRU: Her şifreleme için yeni IV
/// iv := make([]byte, 12) // GCM için 12 byte
/// rand.Read(iv)
///
/// // ❌ YANLIŞ: Sabit IV
/// iv := []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}
/// ```
///
/// ### 4. Veri Şifreleme
/// ```go
/// // Şifrelenecek veri tipleri
/// type SensitiveData struct {
///     SSN             string `encrypt:"true"`  // Sosyal güvenlik no
///     CreditCard      string `encrypt:"true"`  // Kredi kartı
///     Password        string `encrypt:"true"`  // Şifre
///     Email           string `encrypt:"false"` // Email (aranabilir olmalı)
///     Name            string `encrypt:"false"` // İsim (aranabilir olmalı)
/// }
/// ```
///
/// ## Performans Optimizasyonu
///
/// ### 1. Selective Encryption
/// - Sadece hassas alanları şifreleyin
/// - Aranması gereken alanları şifrelemeyin
/// - Metadata'yı plain text tutun
///
/// ### 2. Caching
/// - Deşifre edilmiş verileri cache'leyin (dikkatli!)
/// - Key'leri memory'de cache'leyin
/// - Connection pooling kullanın
///
/// ### 3. Batch Operations
/// - Toplu şifreleme/deşifreleme yapın
/// - Paralel işleme kullanın
/// - Async encryption düşünün
///
/// ## Compliance Gereksinimleri
///
/// ### GDPR (General Data Protection Regulation)
/// - Kişisel verilerin şifrelenmesi önerilir
/// - Data breach durumunda şifreli veri "güvenli" sayılır
/// - Right to erasure: Key silme = veri silme
///
/// ### PCI-DSS (Payment Card Industry)
/// - Kart verileri mutlaka şifrelenmeli
/// - Minimum AES-256 gerekli
/// - Key rotation zorunlu
/// - Key'ler HSM'de saklanmalı
///
/// ### HIPAA (Health Insurance Portability)
/// - Sağlık verileri şifrelenmeli
/// - Encryption at rest ve in transit
/// - Key management audit trail
///
/// ## Key Rotation Örneği
///
/// ```go
/// type KeyRotation struct {
///     CurrentKey  string
///     PreviousKey string
///     NextRotation time.Time
/// }
///
/// func (kr *KeyRotation) Rotate() error {
///     // 1. Yeni key oluştur
///     newKey := generateKey()
///
///     // 2. Eski key'i sakla
///     kr.PreviousKey = kr.CurrentKey
///
///     // 3. Yeni key'i aktif et
///     kr.CurrentKey = newKey
///
///     // 4. Sonraki rotation zamanını ayarla
///     kr.NextRotation = time.Now().Add(90 * 24 * time.Hour)
///
///     // 5. Verileri re-encrypt et (async)
///     go kr.reencryptData()
///
///     return nil
/// }
/// ```
///
/// ## Güvenlik Standartları
///
/// - **NIST FIPS 140-2**: Cryptographic module standardı
/// - **NIST SP 800-57**: Key management önerileri
/// - **OWASP**: Cryptographic Storage Cheat Sheet
/// - **ISO 27001**: Information security management
///
/// ## İlgili Kaynaklar
///
/// - OWASP Crypto: https://cheatsheetseries.owasp.org/cheatsheets/Cryptographic_Storage_Cheat_Sheet.html
/// - NIST Guidelines: https://csrc.nist.gov/publications/detail/sp/800-57-part-1/rev-5/final
/// - Go Crypto: https://pkg.go.dev/crypto
type EncryptionConfig struct {
	/// Kullanılacak şifreleme algoritması.
	///
	/// **Desteklenen Algoritmalar**:
	/// - `"AES-GCM"`: AES-256 Galois/Counter Mode (ÖNERİLİR)
	///   - Authenticated encryption (veri bütünlüğü dahil)
	///   - Paralel işleme desteği
	///   - Modern, hızlı, güvenli
	///   - NIST onaylı
	///
	/// - `"AES-CBC"`: AES-256 Cipher Block Chaining (Legacy)
	///   - Sadece confidentiality (ayrı MAC gerekir)
	///   - Sıralı işleme
	///   - Eski sistemler için
	///   - Padding oracle saldırılarına dikkat
	///
	/// **Önerilen**: "AES-GCM" (yeni projeler için)
	///
	/// **UYARI**: Algoritma değişikliği mevcut şifreli verileri okunamaz hale getirir!
	/// Migration planı yapın.
	Algorithm string

	/// Hex formatında şifreleme key'i.
	/// Key, hexadecimal string olarak saklanır (örn: "0123456789abcdef...").
	///
	/// **Key Boyutları**:
	/// - 32 byte (64 hex karakter): AES-256 (ÖNERİLİR)
	/// - 24 byte (48 hex karakter): AES-192
	/// - 16 byte (32 hex karakter): AES-128
	///
	/// **Örnek Key Generation**:
	/// ```bash
	/// # OpenSSL ile
	/// openssl rand -hex 32
	///
	/// # Go ile
	/// go run -c 'package main; import ("crypto/rand"; "encoding/hex"; "fmt"); func main() { b := make([]byte, 32); rand.Read(b); fmt.Println(hex.EncodeToString(b)) }'
	/// ```
	///
	/// **KRİTİK UYARILAR**:
	/// - ❌ ASLA kodda hardcode etmeyin!
	/// - ❌ ASLA git'e commit etmeyin!
	/// - ❌ ASLA log'lara yazdırmayın!
	/// - ✅ Environment variable kullanın: `os.Getenv("ENCRYPTION_KEY")`
	/// - ✅ Key vault kullanın: AWS KMS, Azure Key Vault, HashiCorp Vault
	/// - ✅ Güvenli bir yerde yedekleyin (key kaybı = veri kaybı)
	/// - ✅ Production ve development için farklı key'ler kullanın
	///
	/// **Örnek Kullanım**:
	/// ```go
	/// // Environment variable'dan oku
	/// config.KeyHex = os.Getenv("ENCRYPTION_KEY")
	/// if config.KeyHex == "" {
	///     log.Fatal("ENCRYPTION_KEY environment variable not set")
	/// }
	///
	/// // Key vault'tan oku
	/// config.KeyHex, err = kmsClient.GetKey("encryption-key")
	/// ```
	KeyHex string

	/// Key rotation (key döndürme) özelliğinin aktif olup olmadığını belirtir.
	///
	/// **true**: Otomatik key rotation aktif (production için önerilir)
	/// **false**: Key rotation devre dışı (development için)
	///
	/// **Key Rotation Nedir?**
	/// Belirli aralıklarla şifreleme key'ini değiştirme işlemidir.
	/// Güvenlik best practice'i olarak önerilir.
	///
	/// **Avantajlar**:
	/// - Key'in ele geçirilme riskini azaltır
	/// - Compliance gereksinimlerini karşılar
	/// - Güvenlik olaylarının etkisini sınırlar
	///
	/// **Dezavantajlar**:
	/// - Karmaşıklık ekler
	/// - Re-encryption gerektirir
	/// - Downtime riski
	///
	/// **UYARI**: Key rotation aktifse, eski key'leri de saklamalısınız!
	/// Eski veriler eski key'lerle deşifre edilmelidir.
	RotationEnabled bool

	/// Key rotation aralığı (ne sıklıkla key değiştirilecek).
	///
	/// **Önerilen Değerler**:
	/// - Yüksek güvenlik: 30 * 24 * time.Hour (30 gün)
	/// - Standart: 90 * 24 * time.Hour (90 gün) - PCI-DSS önerisi
	/// - Düşük risk: 180 * 24 * time.Hour (180 gün)
	/// - Yıllık: 365 * 24 * time.Hour (1 yıl)
	///
	/// **PCI-DSS**: Minimum 90 günde bir rotation gerektirir
	///
	/// **UYARI**: Çok sık rotation performans sorununa yol açabilir!
	/// Çok seyrek rotation güvenlik riski yaratır.
	///
	/// **Key Rotation Süreci**:
	/// 1. Yeni key oluştur
	/// 2. Yeni veriler yeni key ile şifrele
	/// 3. Eski veriler eski key ile okunabilir kal
	/// 4. Background'da eski verileri yeni key ile re-encrypt et
	/// 5. Re-encryption tamamlandığında eski key'i retire et
	///
	/// **Not**: RotationEnabled=false ise bu değer kullanılmaz.
	RotationInterval time.Duration
}

/// # AuditConfig
///
/// Bu yapı, denetim günlüğü (audit logging) yapılandırmasını yönetir.
/// Güvenlik olaylarını, kullanıcı aktivitelerini ve sistem değişikliklerini
/// kaydetmek için kullanılır. Compliance gereksinimleri ve güvenlik analizi için kritiktir.
///
/// ## Kullanım Senaryoları
///
/// - **Güvenlik İzleme**: Şüpheli aktiviteleri ve saldırı girişimlerini kaydetme
/// - **Compliance**: GDPR, HIPAA, PCI-DSS, SOX gibi standartların gereksinimlerini karşılama
/// - **Forensic Analysis**: Güvenlik olaylarının sonradan analiz edilmesi
/// - **User Activity Tracking**: Kullanıcı işlemlerinin izlenmesi
/// - **Change Management**: Sistem değişikliklerinin kaydedilmesi
/// - **Incident Response**: Güvenlik olaylarına müdahale için veri toplama
///
/// ## Örnek Kullanım
///
/// ```go
/// // Production için kapsamlı audit logging
/// prodAudit := config.AuditConfig{
///     Enabled: true,
///     LogLevel: "all",                    // Tüm olayları logla
///     Destination: "siem",                // SIEM sistemine gönder
///     SIEMEndpoint: "https://siem.example.com/api/logs",
/// }
///
/// // File-based audit logging
/// fileAudit := config.AuditConfig{
///     Enabled: true,
///     LogLevel: "security",               // Sadece güvenlik olayları
///     Destination: "file",
///     FilePath: "/var/log/panel/audit.log",
/// }
///
/// // Development için console logging
/// devAudit := config.AuditConfig{
///     Enabled: true,
///     LogLevel: "errors",                 // Sadece hatalar
///     Destination: "console",
/// }
///
/// // Minimal logging (sadece kritik olaylar)
/// minimalAudit := config.AuditConfig{
///     Enabled: true,
///     LogLevel: "security",
///     Destination: "file",
///     FilePath: "/var/log/audit.log",
/// }
///
/// // Audit logging devre dışı (sadece development)
/// noAudit := config.AuditConfig{
///     Enabled: false,
/// }
/// ```
///
/// ## Avantajlar
///
/// - **Güvenlik**: Saldırıları ve şüpheli aktiviteleri tespit etme
/// - **Compliance**: Yasal gereksinimleri karşılama ve denetim kolaylığı
/// - **Accountability**: Kullanıcı işlemlerinin izlenebilirliği
/// - **Forensics**: Olay sonrası analiz ve kanıt toplama
/// - **Monitoring**: Gerçek zamanlı güvenlik izleme
/// - **Alerting**: Anormal davranışlarda otomatik uyarı
///
/// ## Dezavantajlar
///
/// - **Performans**: Log yazma işlemi I/O overhead yaratır
/// - **Disk Kullanımı**: Log dosyaları hızla büyüyebilir
/// - **Privacy**: Hassas verilerin loglanması GDPR sorununa yol açabilir
/// - **Karmaşıklık**: Log yönetimi ve analizi kaynak gerektirir
/// - **Cost**: SIEM sistemleri maliyetli olabilir
///
/// ## Önemli Notlar
///
/// ⚠️ **KRİTİK GÜVENLİK UYARILARI**:
/// - **ASLA** şifreleri, token'ları veya API key'leri loglama!
/// - **ASLA** kredi kartı numaralarını tam olarak loglama!
/// - **MUTLAKA** hassas verileri maskeleyerek logla (örn: **** **** **** 1234)
/// - **MUTLAKA** production'da audit logging'i aktif tutun!
/// - **MUTLAKA** log dosyalarını düzenli olarak rotate edin!
/// - **MUTLAKA** log dosyalarına erişimi kısıtlayın (chmod 600)
///
/// ⚠️ **LogLevel Seçimi**:
/// - **"all"**: Tüm olaylar (en kapsamlı, en yüksek overhead)
///   - Kullanım: Production, compliance gerektiren sistemler
///   - Loglar: Login, logout, CRUD, API calls, errors, security events
///   - Disk kullanımı: Yüksek
///
/// - **"security"**: Sadece güvenlik olayları (dengeli)
///   - Kullanım: Çoğu production ortamı (önerilen)
///   - Loglar: Failed login, permission denied, suspicious activity
///   - Disk kullanımı: Orta
///
/// - **"errors"**: Sadece hatalar (minimal)
///   - Kullanım: Development, düşük riskli sistemler
///   - Loglar: Exceptions, errors, critical failures
///   - Disk kullanımı: Düşük
///
/// ⚠️ **Destination Seçimi**:
/// - **"console"**: Stdout'a yazdır
///   - Avantaj: Basit, hızlı, development için ideal
///   - Dezavantaj: Kalıcı değil, production için uygun değil
///   - Kullanım: Development, debugging
///
/// - **"file"**: Dosyaya yazdır
///   - Avantaj: Kalıcı, analiz edilebilir, basit
///   - Dezavantaj: Disk yönetimi gerekir, rotation şart
///   - Kullanım: Küçük-orta ölçekli production
///   - **UYARI**: Log rotation yapılandırın (logrotate)!
///
/// - **"siem"**: SIEM sistemine gönder
///   - Avantaj: Merkezi yönetim, gelişmiş analiz, alerting
///   - Dezavantaj: Maliyetli, karmaşık, network bağımlı
///   - Kullanım: Enterprise, yüksek güvenlik gerektiren sistemler
///   - Örnekler: Splunk, ELK Stack, Datadog, Sumo Logic
///
/// ⚠️ **Loglanması Gereken Olaylar**:
///
/// **Güvenlik Olayları (Security Events)**:
/// - ✅ Başarılı/başarısız login denemeleri
/// - ✅ Hesap kilitleme/açma işlemleri
/// - ✅ Şifre değişiklikleri
/// - ✅ Permission denied hataları
/// - ✅ Rate limit aşımları
/// - ✅ Şüpheli IP'lerden gelen istekler
/// - ✅ CSRF token hataları
/// - ✅ Session hijacking denemeleri
///
/// **Kullanıcı Aktiviteleri (User Activities)**:
/// - ✅ CRUD işlemleri (Create, Read, Update, Delete)
/// - ✅ Dosya yükleme/indirme
/// - ✅ Export işlemleri
/// - ✅ Admin panel erişimleri
/// - ✅ Rol/permission değişiklikleri
///
/// **Sistem Olayları (System Events)**:
/// - ✅ Uygulama başlatma/durdurma
/// - ✅ Configuration değişiklikleri
/// - ✅ Database migration'ları
/// - ✅ Kritik hatalar ve exception'lar
///
/// ⚠️ **Loglanmaması Gereken Veriler**:
/// - ❌ Şifreler (plain text veya hash)
/// - ❌ API key'ler ve secret'lar
/// - ❌ Session token'ları
/// - ❌ Kredi kartı numaraları (tam)
/// - ❌ Sosyal güvenlik numaraları
/// - ❌ Sağlık bilgileri (HIPAA)
/// - ❌ Request/response body'leri (hassas veri içerebilir)
///
/// ⚠️ **Veri Maskeleme Örnekleri**:
/// ```go
/// // Kredi kartı
/// "4532 **** **** 1234"  // ✅ Doğru
/// "4532 1234 5678 1234"  // ❌ Yanlış
///
/// // Email
/// "u***@example.com"     // ✅ Doğru
/// "user@example.com"     // ⚠️ Dikkatli (GDPR)
///
/// // IP Address
/// "192.168.1.***"        // ✅ Doğru (privacy)
/// "192.168.1.100"        // ⚠️ Dikkatli (GDPR)
///
/// // Şifre
/// "[REDACTED]"           // ✅ Doğru
/// "password123"          // ❌ ASLA!
/// ```
///
/// ## Log Format Önerileri
///
/// ### Structured Logging (JSON)
/// ```json
/// {
///   "timestamp": "2026-02-07T13:35:36Z",
///   "level": "security",
///   "event": "login_failed",
///   "user_id": "12345",
///   "username": "john.doe",
///   "ip_address": "192.168.1.100",
///   "user_agent": "Mozilla/5.0...",
///   "reason": "invalid_password",
///   "attempt_count": 3
/// }
/// ```
///
/// ### Plain Text Logging
/// ```
/// [2026-02-07 13:35:36] SECURITY: Login failed for user john.doe from 192.168.1.100 (attempt 3/5)
/// ```
///
/// ## Log Rotation Yapılandırması
///
/// ### logrotate (Linux)
/// ```
/// /var/log/panel/audit.log {
///     daily                    # Günlük rotation
///     rotate 90                # 90 gün sakla
///     compress                 # Sıkıştır
///     delaycompress            # Bir sonraki rotation'da sıkıştır
///     notifempty               # Boş dosyaları rotate etme
///     create 0600 app app      # Yeni dosya izinleri
///     postrotate
///         systemctl reload app # Uygulamayı reload et
///     endscript
/// }
/// ```
///
/// ## SIEM Entegrasyonu
///
/// ### HTTP POST Örneği
/// ```go
/// func sendToSIEM(event AuditEvent) error {
///     payload, _ := json.Marshal(event)
///     resp, err := http.Post(
///         config.SIEMEndpoint,
///         "application/json",
///         bytes.NewBuffer(payload),
///     )
///     if err != nil {
///         // Fallback: Local file'a yaz
///         logToFile(event)
///         return err
///     }
///     return nil
/// }
/// ```
///
/// ### Syslog Örneği
/// ```go
/// import "log/syslog"
///
/// func sendToSyslog(event AuditEvent) error {
///     logger, err := syslog.New(syslog.LOG_INFO|syslog.LOG_LOCAL0, "panel")
///     if err != nil {
///         return err
///     }
///     defer logger.Close()
///     return logger.Info(event.String())
/// }
/// ```
///
/// ## Performans Optimizasyonu
///
/// ### 1. Async Logging
/// ```go
/// // Log'ları async olarak yaz
/// logChan := make(chan AuditEvent, 1000)
/// go func() {
///     for event := range logChan {
///         writeLog(event)
///     }
/// }()
/// ```
///
/// ### 2. Buffered Writing
/// ```go
/// // Buffer kullanarak disk I/O'sunu azalt
/// writer := bufio.NewWriter(file)
/// defer writer.Flush()
/// ```
///
/// ### 3. Sampling
/// ```go
/// // Yüksek frekanslı olayları sample'la
/// if event.Type == "api_call" && rand.Float64() > 0.1 {
///     return // %10'unu logla
/// }
/// ```
///
/// ## Compliance Gereksinimleri
///
/// ### GDPR (General Data Protection Regulation)
/// - Kişisel verilerin işlenmesini logla
/// - Log retention policy belirle (genelde 1-2 yıl)
/// - Kullanıcı talebi üzerine logları sil
/// - Hassas verileri maskele
///
/// ### PCI-DSS (Payment Card Industry)
/// - Tüm kart işlemlerini logla
/// - Log'ları minimum 1 yıl sakla
/// - Log'lara erişimi kısıtla
/// - Log integrity kontrolü yap
///
/// ### HIPAA (Health Insurance Portability)
/// - Sağlık verilerine erişimi logla
/// - Log'ları minimum 6 yıl sakla
/// - Audit trail oluştur
/// - Log'ları şifrele
///
/// ### SOX (Sarbanes-Oxley)
/// - Finansal veri değişikliklerini logla
/// - Log'ları 7 yıl sakla
/// - Tamper-proof logging
/// - Regular audit
///
/// ## Monitoring ve Alerting
///
/// ### Kritik Olaylar için Alert
/// ```go
/// // 5 dakikada 10'dan fazla başarısız login
/// if failedLoginCount > 10 && timeWindow < 5*time.Minute {
///     sendAlert("Possible brute-force attack detected")
/// }
///
/// // Yeni IP'den admin erişimi
/// if isAdminAccess && isNewIP {
///     sendAlert("Admin access from new IP: " + ip)
/// }
///
/// // Toplu veri export
/// if exportedRecords > 1000 {
///     sendAlert("Large data export: " + exportedRecords + " records")
/// }
/// ```
///
/// ## Güvenlik Standartları
///
/// - **OWASP**: Logging Cheat Sheet
/// - **NIST 800-92**: Guide to Computer Security Log Management
/// - **ISO 27001**: Information security logging requirements
/// - **CIS Controls**: Audit log management
///
/// ## İlgili Kaynaklar
///
/// - OWASP Logging: https://cheatsheetseries.owasp.org/cheatsheets/Logging_Cheat_Sheet.html
/// - NIST 800-92: https://csrc.nist.gov/publications/detail/sp/800-92/final
type AuditConfig struct {
	/// Audit logging'in aktif olup olmadığını belirtir.
	///
	/// **true**: Audit logging aktif (production için şiddetle önerilir)
	/// **false**: Audit logging devre dışı (sadece development için)
	///
	/// **UYARI**: Production ortamında MUTLAKA true olmalıdır!
	/// Audit logging olmadan:
	/// - Güvenlik olaylarını tespit edemezsiniz
	/// - Compliance gereksinimlerini karşılayamazsınız
	/// - Forensic analysis yapılamaz
	/// - Kullanıcı aktivitelerini izleyemezsiniz
	///
	/// **Compliance**: GDPR, PCI-DSS, HIPAA, SOX audit logging gerektirir
	Enabled bool

	/// Hangi seviyedeki olayların loglanacağını belirtir.
	///
	/// **Değerler**:
	/// - `"all"`: Tüm olaylar loglanır (en kapsamlı)
	///   - Login/logout, CRUD işlemleri, API calls, errors, security events
	///   - Kullanım: Compliance gerektiren sistemler, yüksek güvenlik
	///   - Disk kullanımı: Yüksek (log rotation şart!)
	///   - Performans etkisi: Orta-yüksek
	///
	/// - `"security"`: Sadece güvenlik olayları (önerilen)
	///   - Failed login, permission denied, rate limit, suspicious activity
	///   - Kullanım: Çoğu production ortamı
	///   - Disk kullanımı: Orta
	///   - Performans etkisi: Düşük-orta
	///
	/// - `"errors"`: Sadece hatalar (minimal)
	///   - Exceptions, errors, critical failures
	///   - Kullanım: Development, düşük riskli sistemler
	///   - Disk kullanımı: Düşük
	///   - Performans etkisi: Düşük
	///
	/// **Önerilen**: Production'da "security" veya "all"
	///
	/// **Not**: LogLevel ne kadar yüksek o kadar fazla disk ve performans maliyeti.
	/// İhtiyacınıza göre dengeleyin.
	LogLevel string

	/// Log'ların nereye yazılacağını belirtir.
	///
	/// **Değerler**:
	/// - `"console"`: Stdout'a yazdır
	///   - Avantaj: Basit, hızlı, development için ideal
	///   - Dezavantaj: Kalıcı değil, production için uygun değil
	///   - Kullanım: Development, debugging, Docker container'lar (stdout logging)
	///   - **UYARI**: Production'da tek başına kullanmayın!
	///
	/// - `"file"`: Dosyaya yazdır
	///   - Avantaj: Kalıcı, analiz edilebilir, basit setup
	///   - Dezavantaj: Disk yönetimi gerekir, log rotation şart
	///   - Kullanım: Küçük-orta ölçekli production
	///   - **UYARI**: FilePath belirtilmeli, log rotation yapılandırın!
	///   - **Güvenlik**: Dosya izinlerini kısıtlayın (chmod 600)
	///
	/// - `"siem"`: SIEM sistemine gönder
	///   - Avantaj: Merkezi yönetim, gelişmiş analiz, real-time alerting
	///   - Dezavantaj: Maliyetli, karmaşık, network bağımlı
	///   - Kullanım: Enterprise, yüksek güvenlik gerektiren sistemler
	///   - **UYARI**: SIEMEndpoint belirtilmeli, fallback mekanizması olmalı
	///   - Örnekler: Splunk, ELK Stack, Datadog, Sumo Logic, Azure Sentinel
	///
	/// **Önerilen**: Production'da "file" veya "siem"
	///
	/// **Best Practice**: Hybrid yaklaşım kullanın:
	/// - Primary: SIEM (real-time monitoring)
	/// - Fallback: File (SIEM erişilemezse)
	/// - Development: Console
	Destination string

	/// Log dosyasının tam yolu (Destination="file" ise gerekli).
	///
	/// **Önerilen Konumlar**:
	/// - Linux: `/var/log/panel/audit.log`
	/// - Docker: `/app/logs/audit.log` (volume mount)
	/// - Windows: `C:\ProgramData\Panel\logs\audit.log`
	///
	/// **Dosya İzinleri**:
	/// ```bash
	/// # Sadece uygulama kullanıcısı okuyabilsin
	/// chmod 600 /var/log/panel/audit.log
	/// chown app:app /var/log/panel/audit.log
	/// ```
	///
	/// **Log Rotation**:
	/// ```bash
	/// # logrotate yapılandırması
	/// /etc/logrotate.d/panel-audit
	/// ```
	///
	/// **UYARI**:
	/// - Dizinin var olduğundan emin olun
	/// - Yazma izni olduğundan emin olun
	/// - Disk doluluk kontrolü yapın
	/// - Log rotation yapılandırın (aksi halde disk dolar!)
	///
	/// **Not**: Destination="file" değilse bu alan kullanılmaz.
	FilePath string

	/// SIEM endpoint URL'i (Destination="siem" ise gerekli).
	///
	/// **Format**: `https://siem.example.com/api/logs`
	///
	/// **Örnekler**:
	/// - Splunk: `https://splunk.example.com:8088/services/collector`
	/// - ELK: `https://elasticsearch.example.com:9200/audit-logs/_doc`
	/// - Datadog: `https://http-intake.logs.datadoghq.com/v1/input`
	/// - Sumo Logic: `https://collectors.sumologic.com/receiver/v1/http/[token]`
	///
	/// **Güvenlik**:
	/// - HTTPS kullanın (HTTP kullanmayın!)
	/// - Authentication token'ı environment variable'dan okuyun
	/// - Network timeout ayarlayın
	/// - Retry mekanizması ekleyin
	/// - Fallback olarak local file'a yazın
	///
	/// **Örnek Kullanım**:
	/// ```go
	/// config.SIEMEndpoint = os.Getenv("SIEM_ENDPOINT")
	/// if config.SIEMEndpoint == "" {
	///     log.Fatal("SIEM_ENDPOINT not set")
	/// }
	/// ```
	///
	/// **UYARI**:
	/// - Endpoint erişilebilir olmalı
	/// - Authentication gerekiyorsa header'larda gönder
	/// - Rate limiting'e dikkat et
	/// - Network hatalarını handle et
	///
	/// **Not**: Destination="siem" değilse bu alan kullanılmaz.
	SIEMEndpoint string
}

/// # DefaultSecurityConfig
///
/// Bu fonksiyon, güvenli varsayılan değerlerle yapılandırılmış bir SecurityConfig döndürür.
/// Tüm ortamlar için temel güvenlik ayarlarını içerir ve özelleştirme için başlangıç noktası sağlar.
///
/// ## Kullanım Senaryoları
///
/// - **Hızlı Başlangıç**: Yeni projelerde güvenli varsayılan ayarlarla başlama
/// - **Temel Yapılandırma**: Ortam bazlı özelleştirmeler için temel oluşturma
/// - **Güvenlik Baseline**: Minimum güvenlik gereksinimlerini karşılama
/// - **Development**: Geliştirme ortamında güvenli ayarlarla çalışma
///
/// ## Döndürülen Yapılandırma
///
/// ### CORS Ayarları
/// - **AllowedOrigins**: `["http://localhost:3000"]` - Sadece localhost (değiştirilmeli!)
/// - **AllowCredentials**: `true` - Cookie ve auth header'lar aktif
/// - **AllowedMethods**: GET, POST, PUT, DELETE, OPTIONS
/// - **AllowedHeaders**: Content-Type, Authorization, X-CSRF-Token
/// - **MaxAge**: 3600 saniye (1 saat)
///
/// ### Rate Limiting
/// - **Enabled**: `true` - Aktif
/// - **AuthMaxRequests**: 10 istek / dakika
/// - **APIMaxRequests**: 100 istek / dakika
///
/// ### Account Lockout
/// - **Enabled**: `true` - Aktif
/// - **MaxAttempts**: 5 başarısız deneme
/// - **LockoutDuration**: 15 dakika
///
/// ### Session
/// - **CookieName**: `__Host-session_token` - Güvenli prefix
/// - **Secure**: `true` - Sadece HTTPS
/// - **HTTPOnly**: `true` - JavaScript erişimi yok
/// - **SameSite**: `Strict` - CSRF koruması
/// - **MaxAge**: 86400 saniye (24 saat)
///
/// ### Encryption
/// - **Algorithm**: `AES-GCM` - Modern, authenticated encryption
/// - **RotationEnabled**: `false` - Manuel key yönetimi
/// - **RotationInterval**: 90 gün
///
/// ### Audit
/// - **Enabled**: `true` - Aktif
/// - **LogLevel**: `security` - Sadece güvenlik olayları
/// - **Destination**: `console` - Stdout'a yazdır
///
/// ## Örnek Kullanım
///
/// ```go
/// // Varsayılan yapılandırma ile başla
/// config := config.DefaultSecurityConfig()
///
/// // Ortama göre özelleştir
/// if os.Getenv("ENV") == "production" {
///     config.CORS.AllowedOrigins = []string{"https://example.com"}
///     config.Session.Secure = true
///     config.Audit.Destination = "file"
///     config.Audit.FilePath = "/var/log/audit.log"
/// }
///
/// // Encryption key'i environment variable'dan oku
/// config.Encryption.KeyHex = os.Getenv("ENCRYPTION_KEY")
///
/// // Uygulamada kullan
/// app := panel.New(config)
/// ```
///
/// ## Özelleştirme Örnekleri
///
/// ### CORS Özelleştirme
/// ```go
/// config := config.DefaultSecurityConfig()
/// config.CORS.AllowedOrigins = []string{
///     "https://app.example.com",
///     "https://admin.example.com",
/// }
/// ```
///
/// ### Rate Limiting Özelleştirme
/// ```go
/// config := config.DefaultSecurityConfig()
/// config.RateLimit.AuthMaxRequests = 3  // Daha sıkı
/// config.RateLimit.APIMaxRequests = 1000 // Daha esnek
/// ```
///
/// ### Audit Logging Özelleştirme
/// ```go
/// config := config.DefaultSecurityConfig()
/// config.Audit.LogLevel = "all"
/// config.Audit.Destination = "siem"
/// config.Audit.SIEMEndpoint = "https://siem.example.com/api/logs"
/// ```
///
/// ## Önemli Notlar
///
/// ⚠️ **MUTLAKA ÖZELLEŞTİRİLMELİ**:
/// - **CORS AllowedOrigins**: Localhost yerine gerçek domain'leri ekleyin!
/// - **Encryption KeyHex**: Environment variable'dan okuyun!
/// - **Audit Destination**: Production'da "file" veya "siem" kullanın!
///
/// ⚠️ **Production İçin Yetersiz**:
/// Bu yapılandırma temel güvenlik sağlar ancak production için yeterli değildir.
/// `ProductionSecurityConfig()` kullanın veya manuel olarak sıkılaştırın.
///
/// ⚠️ **Development İçin Sıkı**:
/// Development ortamında bazı ayarlar (Secure=true, Strict SameSite) sorun yaratabilir.
/// `DevelopmentSecurityConfig()` kullanın.
///
/// ## Avantajlar
///
/// - **Güvenli Varsayılanlar**: Tüm güvenlik özellikleri aktif
/// - **Hızlı Başlangıç**: Minimal yapılandırma ile çalışır
/// - **Özelleştirilebilir**: Tüm ayarlar değiştirilebilir
/// - **Dokümante**: Her ayar açıklanmış
///
/// ## Dezavantajlar
///
/// - **Generic**: Özel ihtiyaçları karşılamayabilir
/// - **Localhost CORS**: Production için uygun değil
/// - **Console Logging**: Production için yetersiz
/// - **Manual Key**: Encryption key manuel set edilmeli
///
/// ## İlgili Fonksiyonlar
///
/// - `ProductionSecurityConfig()`: Production için optimize edilmiş ayarlar
/// - `DevelopmentSecurityConfig()`: Development için esnek ayarlar
///
/// ## Döndürür
///
/// - Güvenli varsayılan değerlerle yapılandırılmış `SecurityConfig` struct'ı
func DefaultSecurityConfig() SecurityConfig {
	return SecurityConfig{
		CORS: CORSConfig{
			AllowedOrigins:   []string{"http://localhost:3000"}, // Must be configured per environment
			AllowCredentials: true,
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Content-Type", "Authorization", "X-CSRF-Token"},
			ExposeHeaders:    []string{"Content-Length"},
			MaxAge:           3600,
		},
		RateLimit: RateLimitConfig{
			Enabled:         true,
			AuthMaxRequests: 10,
			AuthWindow:      1 * time.Minute,
			APIMaxRequests:  100,
			APIWindow:       1 * time.Minute,
		},
		AccountLockout: AccountLockoutConfig{
			Enabled:         true,
			MaxAttempts:     5,
			LockoutDuration: 15 * time.Minute,
		},
		Session: SessionConfig{
			CookieName: "__Host-session_token",
			Secure:     true,
			HTTPOnly:   true,
			SameSite:   "Strict",
			MaxAge:     86400, // 24 hours
			Domain:     "",
			Path:       "/",
		},
		Encryption: EncryptionConfig{
			Algorithm:        "AES-GCM",
			RotationEnabled:  false,
			RotationInterval: 90 * 24 * time.Hour, // 90 days
		},
		Audit: AuditConfig{
			Enabled:     true,
			LogLevel:    "security",
			Destination: "console",
		},
	}
}

/// # ProductionSecurityConfig
///
/// Bu fonksiyon, production ortamı için optimize edilmiş ve sıkılaştırılmış
/// güvenlik yapılandırması döndürür. Tüm güvenlik özellikleri maksimum seviyede
/// aktif edilmiş ve production best practice'lerine göre ayarlanmıştır.
///
/// ## Kullanım Senaryoları
///
/// - **Production Deployment**: Canlı ortama çıkarken kullanılmalı
/// - **Staging Environment**: Production benzeri test ortamları için
/// - **High Security**: Yüksek güvenlik gerektiren uygulamalar için
/// - **Compliance**: GDPR, PCI-DSS, HIPAA gibi standartları karşılama
///
/// ## DefaultSecurityConfig'den Farkları
///
/// ### CORS Değişiklikleri
/// - **AllowedOrigins**: `[]` (boş) - MUTLAKA yapılandırılmalı!
///   - Varsayılan: `["http://localhost:3000"]`
///   - Production: Gerçek domain'ler eklenmelidir
///
/// ### Session Değişiklikleri
/// - **Secure**: `true` - HTTPS zorunlu (değişmedi)
/// - **SameSite**: `"Strict"` - Maksimum CSRF koruması (değişmedi)
///
/// ### Audit Değişiklikleri
/// - **LogLevel**: `"all"` - Tüm olayları logla
///   - Varsayılan: `"security"`
///   - Production: Kapsamlı loglama
/// - **Destination**: `"file"` - Kalıcı dosya logging
///   - Varsayılan: `"console"`
///   - Production: Dosyaya yazma
/// - **FilePath**: `"/var/log/panel/audit.log"` - Standart log konumu
///   - Varsayılan: Yok
///   - Production: Belirli dosya yolu
///
/// ## Döndürülen Yapılandırma
///
/// ### CORS Ayarları
/// - **AllowedOrigins**: `[]` (BOŞ) - ⚠️ MUTLAKA AYARLANMALI!
/// - **AllowCredentials**: `true`
/// - **AllowedMethods**: GET, POST, PUT, DELETE, OPTIONS
/// - **AllowedHeaders**: Content-Type, Authorization, X-CSRF-Token
/// - **MaxAge**: 3600 saniye
///
/// ### Rate Limiting
/// - **Enabled**: `true`
/// - **AuthMaxRequests**: 10 istek / dakika
/// - **APIMaxRequests**: 100 istek / dakika
///
/// ### Account Lockout
/// - **Enabled**: `true`
/// - **MaxAttempts**: 5 başarısız deneme
/// - **LockoutDuration**: 15 dakika
///
/// ### Session
/// - **CookieName**: `__Host-session_token`
/// - **Secure**: `true` - HTTPS zorunlu
/// - **HTTPOnly**: `true`
/// - **SameSite**: `Strict` - Maksimum CSRF koruması
/// - **MaxAge**: 86400 saniye (24 saat)
///
/// ### Encryption
/// - **Algorithm**: `AES-GCM`
/// - **RotationEnabled**: `false` - Manuel yönetim
/// - **RotationInterval**: 90 gün
///
/// ### Audit
/// - **Enabled**: `true`
/// - **LogLevel**: `all` - Tüm olaylar
/// - **Destination**: `file` - Dosyaya yazma
/// - **FilePath**: `/var/log/panel/audit.log`
///
/// ## Örnek Kullanım
///
/// ```go
/// // Production yapılandırması ile başla
/// config := config.ProductionSecurityConfig()
///
/// // CORS origin'lerini ayarla (ZORUNLU!)
/// config.CORS.AllowedOrigins = []string{
///     "https://app.example.com",
///     "https://admin.example.com",
/// }
///
/// // Encryption key'i environment variable'dan oku (ZORUNLU!)
/// config.Encryption.KeyHex = os.Getenv("ENCRYPTION_KEY")
/// if config.Encryption.KeyHex == "" {
///     log.Fatal("ENCRYPTION_KEY environment variable not set")
/// }
///
/// // SIEM kullanıyorsanız audit destination'ı değiştir
/// if siemEndpoint := os.Getenv("SIEM_ENDPOINT"); siemEndpoint != "" {
///     config.Audit.Destination = "siem"
///     config.Audit.SIEMEndpoint = siemEndpoint
/// }
///
/// // Uygulamayı başlat
/// app := panel.New(config)
/// app.Run(":443") // HTTPS port
/// ```
///
/// ## Production Checklist
///
/// ### ⚠️ ZORUNLU AYARLAMALAR
///
/// #### 1. CORS AllowedOrigins
/// ```go
/// // ❌ YANLIŞ: Boş bırakmak
/// config := config.ProductionSecurityConfig()
/// // AllowedOrigins boş, hiçbir origin'e izin verilmez!
///
/// // ✅ DOĞRU: Gerçek domain'leri ekle
/// config.CORS.AllowedOrigins = []string{
///     "https://app.example.com",
///     "https://admin.example.com",
/// }
/// ```
///
/// #### 2. Encryption Key
/// ```go
/// // ❌ YANLIŞ: Hardcode
/// config.Encryption.KeyHex = "abc123..."
///
/// // ✅ DOĞRU: Environment variable
/// config.Encryption.KeyHex = os.Getenv("ENCRYPTION_KEY")
/// if config.Encryption.KeyHex == "" {
///     log.Fatal("ENCRYPTION_KEY not set")
/// }
/// ```
///
/// #### 3. HTTPS Zorunlu
/// ```go
/// // ❌ YANLIŞ: HTTP port
/// app.Run(":80")
///
/// // ✅ DOĞRU: HTTPS port
/// app.Run(":443")
///
/// // ✅ DOĞRU: TLS yapılandırması
/// app.RunTLS(":443", "cert.pem", "key.pem")
/// ```
///
/// #### 4. Log Directory
/// ```bash
/// # Log dizinini oluştur
/// sudo mkdir -p /var/log/panel
/// sudo chown app:app /var/log/panel
/// sudo chmod 750 /var/log/panel
///
/// # Log rotation yapılandır
/// sudo nano /etc/logrotate.d/panel-audit
/// ```
///
/// ### 📋 İsteğe Bağlı Özelleştirmeler
///
/// #### 1. Rate Limiting Ayarları
/// ```go
/// // Yüksek trafikli API için
/// config.RateLimit.APIMaxRequests = 1000
/// config.RateLimit.APIWindow = 1 * time.Minute
///
/// // Daha sıkı auth koruması için
/// config.RateLimit.AuthMaxRequests = 3
/// config.RateLimit.AuthWindow = 15 * time.Minute
/// ```
///
/// #### 2. Session Ayarları
/// ```go
/// // Kısa ömürlü session (yüksek güvenlik)
/// config.Session.MaxAge = 3600 // 1 saat
///
/// // Subdomain'ler için
/// config.Session.Domain = ".example.com"
/// ```
///
/// #### 3. Encryption Key Rotation
/// ```go
/// // Otomatik key rotation aktif et
/// config.Encryption.RotationEnabled = true
/// config.Encryption.RotationInterval = 30 * 24 * time.Hour // 30 gün
/// ```
///
/// #### 4. SIEM Entegrasyonu
/// ```go
/// // SIEM kullan
/// config.Audit.Destination = "siem"
/// config.Audit.SIEMEndpoint = os.Getenv("SIEM_ENDPOINT")
/// config.Audit.LogLevel = "all"
/// ```
///
/// ## Önemli Notlar
///
/// ⚠️ **KRİTİK UYARILAR**:
/// - **CORS AllowedOrigins BOŞ**: Hiçbir origin'e izin verilmez!
///   - MUTLAKA gerçek domain'leri ekleyin
///   - ASLA "*" kullanmayın
/// - **HTTPS ZORUNLU**: Secure=true ayarı HTTPS gerektirir
///   - HTTP'de cookie set edilmez
///   - TLS sertifikası yapılandırın
/// - **Log Directory**: `/var/log/panel/` dizini var olmalı
///   - Yazma izni olmalı
///   - Log rotation yapılandırın
/// - **Encryption Key**: Environment variable'dan okunmalı
///   - ASLA kodda hardcode etmeyin
///   - Key vault kullanın (AWS KMS, Azure Key Vault)
///
/// ⚠️ **Production Hazırlık**:
/// 1. **Environment Variables**:
///    ```bash
///    export ENCRYPTION_KEY="..."
///    export SIEM_ENDPOINT="https://siem.example.com/api/logs"
///    ```
///
/// 2. **Log Rotation**:
///    ```bash
///    # /etc/logrotate.d/panel-audit
///    /var/log/panel/audit.log {
///        daily
///        rotate 90
///        compress
///        delaycompress
///        notifempty
///        create 0600 app app
///    }
///    ```
///
/// 3. **File Permissions**:
///    ```bash
///    chmod 600 /var/log/panel/audit.log
///    chown app:app /var/log/panel/audit.log
///    ```
///
/// 4. **Firewall Rules**:
///    ```bash
///    # Sadece HTTPS'e izin ver
///    sudo ufw allow 443/tcp
///    sudo ufw deny 80/tcp
///    ```
///
/// ⚠️ **Monitoring**:
/// - Audit log'larını düzenli kontrol edin
/// - Failed login denemelerini izleyin
/// - Rate limit aşımlarını takip edin
/// - Disk kullanımını monitör edin
/// - SIEM alert'lerini yapılandırın
///
/// ⚠️ **Backup**:
/// - Encryption key'lerini güvenli yedekleyin
/// - Audit log'larını yedekleyin
/// - Configuration'ı version control'de tutun
///
/// ## Avantajlar
///
/// - **Maksimum Güvenlik**: Tüm güvenlik özellikleri aktif
/// - **Production Ready**: Best practice'lere göre yapılandırılmış
/// - **Compliance**: Güvenlik standartlarına uygun
/// - **Audit Trail**: Kapsamlı loglama
/// - **HTTPS Zorunlu**: Secure communication
///
/// ## Dezavantajlar
///
/// - **Sıkı Ayarlar**: Bazı meşru kullanımları engelleyebilir
/// - **Yapılandırma Gerekli**: CORS ve encryption key ayarlanmalı
/// - **Disk Kullanımı**: "all" log level çok yer kaplar
/// - **HTTPS Gereksinimi**: TLS sertifikası gerekir
///
/// ## İlgili Fonksiyonlar
///
/// - `DefaultSecurityConfig()`: Temel güvenlik ayarları
/// - `DevelopmentSecurityConfig()`: Development için esnek ayarlar
///
/// ## Döndürür
///
/// - Production ortamı için optimize edilmiş `SecurityConfig` struct'ı
/// - CORS AllowedOrigins boş (manuel ayarlanmalı)
/// - Audit logging dosyaya yazacak şekilde yapılandırılmış
/// - Tüm güvenlik özellikleri maksimum seviyede aktif
func ProductionSecurityConfig() SecurityConfig {
	config := DefaultSecurityConfig()

	// Production-specific overrides
	config.CORS.AllowedOrigins = []string{} // Must be explicitly configured
	config.Session.Secure = true
	config.Session.SameSite = "Strict"
	config.Audit.LogLevel = "all"
	config.Audit.Destination = "file"
	config.Audit.FilePath = "/var/log/panel/audit.log"

	return config
}

/// # DevelopmentSecurityConfig
///
/// Bu fonksiyon, development (geliştirme) ortamı için optimize edilmiş güvenlik
/// yapılandırması döndürür. Güvenlik özelliklerini korurken geliştirme sürecini
/// kolaylaştıracak şekilde esnek ayarlar içerir.
///
/// ## Kullanım Senaryoları
///
/// - **Local Development**: Yerel geliştirme ortamında çalışma
/// - **Testing**: Test senaryolarında esnek ayarlar
/// - **Debugging**: Hata ayıklama sürecinde kolaylık sağlama
/// - **Rapid Prototyping**: Hızlı prototip geliştirme
/// - **Learning**: Öğrenme ve deneme amaçlı kullanım
///
/// ## DefaultSecurityConfig'den Farkları
///
/// ### CORS Değişiklikleri
/// - **AllowedOrigins**: `["http://localhost:3000", "http://localhost:5173"]`
///   - Varsayılan: `["http://localhost:3000"]`
///   - Development: Yaygın frontend port'ları eklendi
///   - Port 3000: Create React App, Next.js
///   - Port 5173: Vite
///
/// ### Session Değişiklikleri
/// - **Secure**: `false` - HTTP'ye izin ver
///   - Varsayılan: `true` (HTTPS zorunlu)
///   - Development: HTTP üzerinde çalışabilir
///   - **UYARI**: Production'da ASLA false kullanmayın!
/// - **SameSite**: `"Lax"` - Daha esnek CSRF koruması
///   - Varsayılan: `"Strict"`
///   - Development: Cross-site navigation'da session korunur
///
/// ### Rate Limiting Değişiklikleri
/// - **AuthMaxRequests**: `50` - Daha yüksek limit
///   - Varsayılan: `10`
///   - Development: Test sırasında engellenmemek için
///
/// ### Audit Değişiklikleri
/// - **Destination**: `"console"` - Stdout'a yazdır
///   - Varsayılan: `"console"`
///   - Development: Terminal'de görüntüleme kolaylığı
///
/// ## Döndürülen Yapılandırma
///
/// ### CORS Ayarları
/// - **AllowedOrigins**: `["http://localhost:3000", "http://localhost:5173"]`
/// - **AllowCredentials**: `true`
/// - **AllowedMethods**: GET, POST, PUT, DELETE, OPTIONS
/// - **AllowedHeaders**: Content-Type, Authorization, X-CSRF-Token
/// - **MaxAge**: 3600 saniye
///
/// ### Rate Limiting
/// - **Enabled**: `true` - Yine de aktif (aşırı istekleri önler)
/// - **AuthMaxRequests**: `50` istek / dakika (esnek)
/// - **APIMaxRequests**: `100` istek / dakika
///
/// ### Account Lockout
/// - **Enabled**: `true` - Aktif
/// - **MaxAttempts**: `5` başarısız deneme
/// - **LockoutDuration**: `15` dakika
///
/// ### Session
/// - **CookieName**: `__Host-session_token`
/// - **Secure**: `false` - HTTP'de çalışır
/// - **HTTPOnly**: `true` - XSS koruması aktif
/// - **SameSite**: `Lax` - Esnek CSRF koruması
/// - **MaxAge**: `86400` saniye (24 saat)
///
/// ### Encryption
/// - **Algorithm**: `AES-GCM`
/// - **RotationEnabled**: `false`
/// - **RotationInterval**: `90` gün
///
/// ### Audit
/// - **Enabled**: `true`
/// - **LogLevel**: `errors` - Sadece hatalar
/// - **Destination**: `console` - Terminal'e yazdır
///
/// ## Örnek Kullanım
///
/// ### Temel Kullanım
/// ```go
/// // Development yapılandırması
/// config := config.DevelopmentSecurityConfig()
///
/// // Encryption key'i test key ile ayarla
/// config.Encryption.KeyHex = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
///
/// // Uygulamayı başlat
/// app := panel.New(config)
/// app.Run(":8080") // HTTP port
/// ```
///
/// ### Ortam Bazlı Yapılandırma
/// ```go
/// var config config.SecurityConfig
///
/// switch os.Getenv("ENV") {
/// case "production":
///     config = config.ProductionSecurityConfig()
///     config.CORS.AllowedOrigins = []string{"https://example.com"}
///     config.Encryption.KeyHex = os.Getenv("ENCRYPTION_KEY")
/// case "staging":
///     config = config.ProductionSecurityConfig()
///     config.CORS.AllowedOrigins = []string{"https://staging.example.com"}
///     config.Encryption.KeyHex = os.Getenv("ENCRYPTION_KEY")
/// default: // development
///     config = config.DevelopmentSecurityConfig()
///     config.Encryption.KeyHex = "test-key-for-development-only"
/// }
///
/// app := panel.New(config)
/// ```
///
/// ### Docker Compose ile Kullanım
/// ```yaml
/// # docker-compose.yml
/// version: '3.8'
/// services:
///   app:
///     build: .
///     ports:
///       - "8080:8080"
///     environment:
///       - ENV=development
///       - ENCRYPTION_KEY=test-key
///     volumes:
///       - .:/app
/// ```
///
/// ### Frontend Entegrasyonu
/// ```javascript
/// // React development server (port 3000)
/// const API_URL = 'http://localhost:8080';
///
/// fetch(`${API_URL}/api/users`, {
///   credentials: 'include', // Cookie gönder
///   headers: {
///     'Content-Type': 'application/json',
///   },
/// });
///
/// // Vite development server (port 5173)
/// // vite.config.js
/// export default {
///   server: {
///     port: 5173,
///     proxy: {
///       '/api': 'http://localhost:8080'
///     }
///   }
/// }
/// ```
///
/// ## Özelleştirme Örnekleri
///
/// ### Ek Frontend Port Ekleme
/// ```go
/// config := config.DevelopmentSecurityConfig()
/// config.CORS.AllowedOrigins = append(
///     config.CORS.AllowedOrigins,
///     "http://localhost:4200", // Angular
///     "http://localhost:8000", // Django
/// )
/// ```
///
/// ### Daha Detaylı Logging
/// ```go
/// config := config.DevelopmentSecurityConfig()
/// config.Audit.LogLevel = "all" // Tüm olayları logla
/// ```
///
/// ### Rate Limiting Devre Dışı (Dikkatli!)
/// ```go
/// config := config.DevelopmentSecurityConfig()
/// config.RateLimit.Enabled = false // Sadece development için!
/// ```
///
/// ### Test Veritabanı ile Kullanım
/// ```go
/// config := config.DevelopmentSecurityConfig()
///
/// // Test database
/// db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
///
/// // Test encryption key
/// config.Encryption.KeyHex = "test-key-32-bytes-hex-encoded-here"
///
/// app := panel.New(config)
/// ```
///
/// ## Önemli Notlar
///
/// ⚠️ **GÜVENLİK UYARILARI**:
/// - **ASLA** production'da bu yapılandırmayı kullanmayın!
/// - **Secure=false**: HTTP üzerinde cookie gönderilir (güvenli değil)
/// - **Esnek Rate Limiting**: Brute-force saldırılarına daha açık
/// - **Console Logging**: Log'lar kalıcı değil
/// - **Test Key'ler**: Production'da gerçek key'ler kullanın
///
/// ⚠️ **Development Best Practices**:
/// 1. **Environment Variable Kullanın**:
///    ```bash
///    export ENV=development
///    export ENCRYPTION_KEY=test-key
///    ```
///
/// 2. **Git Ignore**:
///    ```gitignore
///    # .gitignore
///    .env
///    .env.local
///    test.db
///    *.log
///    ```
///
/// 3. **Separate Config Files**:
///    ```
///    config/
///    ├── development.yaml
///    ├── staging.yaml
///    └── production.yaml
///    ```
///
/// 4. **Docker Development**:
///    ```dockerfile
///    # Dockerfile.dev
///    FROM golang:1.21
///    WORKDIR /app
///    COPY . .
///    RUN go mod download
///    CMD ["go", "run", "main.go"]
///    ```
///
/// ⚠️ **Yaygın Sorunlar ve Çözümler**:
///
/// **1. CORS Hatası**
/// ```
/// Error: CORS policy: No 'Access-Control-Allow-Origin' header
/// ```
/// Çözüm: Frontend port'unu AllowedOrigins'e ekleyin
/// ```go
/// config.CORS.AllowedOrigins = append(
///     config.CORS.AllowedOrigins,
///     "http://localhost:YOUR_PORT",
/// )
/// ```
///
/// **2. Cookie Set Edilmiyor**
/// ```
/// Warning: Cookie not set in browser
/// ```
/// Çözüm: SameSite="None" ve Secure=true gerekebilir (HTTPS ile)
/// ```go
/// config.Session.SameSite = "None"
/// config.Session.Secure = true // HTTPS gerektirir
/// ```
///
/// **3. Rate Limit Aşımı**
/// ```
/// Error: 429 Too Many Requests
/// ```
/// Çözüm: Limitleri artırın veya devre dışı bırakın
/// ```go
/// config.RateLimit.AuthMaxRequests = 100
/// // veya
/// config.RateLimit.Enabled = false
/// ```
///
/// **4. Session Kaybolması**
/// ```
/// Issue: Session lost after page refresh
/// ```
/// Çözüm: Cookie domain ve path ayarlarını kontrol edin
/// ```go
/// config.Session.Domain = "" // Boş bırakın
/// config.Session.Path = "/"
/// ```
///
/// ## Development Workflow
///
/// ### 1. İlk Kurulum
/// ```bash
/// # Projeyi klonla
/// git clone https://github.com/example/panel.go
/// cd panel.go
///
/// # Dependencies
/// go mod download
///
/// # Environment variables
/// cp .env.example .env
/// nano .env
/// ```
///
/// ### 2. Development Server Başlatma
/// ```bash
/// # Backend (Go)
/// ENV=development go run main.go
///
/// # Frontend (React)
/// cd web
/// npm install
/// npm run dev
/// ```
///
/// ### 3. Hot Reload ile Geliştirme
/// ```bash
/// # Air ile hot reload
/// air
///
/// # veya
/// go install github.com/cosmtrek/air@latest
/// air
/// ```
///
/// ### 4. Testing
/// ```bash
/// # Unit tests
/// go test ./...
///
/// # Integration tests
/// ENV=test go test -tags=integration ./...
///
/// # Coverage
/// go test -cover ./...
/// ```
///
/// ## Frontend Development Servers
///
/// ### Create React App (Port 3000)
/// ```json
/// // package.json
/// {
///   "proxy": "http://localhost:8080"
/// }
/// ```
///
/// ### Vite (Port 5173)
/// ```javascript
/// // vite.config.js
/// export default {
///   server: {
///     port: 5173,
///     proxy: {
///       '/api': {
///         target: 'http://localhost:8080',
///         changeOrigin: true,
///       }
///     }
///   }
/// }
/// ```
///
/// ### Next.js (Port 3000)
/// ```javascript
/// // next.config.js
/// module.exports = {
///   async rewrites() {
///     return [
///       {
///         source: '/api/:path*',
///         destination: 'http://localhost:8080/api/:path*',
///       },
///     ]
///   },
/// }
/// ```
///
/// ## Debugging Tips
///
/// ### 1. Verbose Logging
/// ```go
/// config := config.DevelopmentSecurityConfig()
/// config.Audit.LogLevel = "all"
/// ```
///
/// ### 2. CORS Debug
/// ```go
/// // Tüm origin'lere izin ver (sadece debug için!)
/// config.CORS.AllowedOrigins = []string{"*"}
/// config.CORS.AllowCredentials = false // * ile true kullanılamaz
/// ```
///
/// ### 3. Rate Limit Debug
/// ```go
/// // Rate limiting'i geçici olarak kapat
/// config.RateLimit.Enabled = false
/// ```
///
/// ### 4. Session Debug
/// ```go
/// // Session cookie'yi browser'da görmek için
/// config.Session.HTTPOnly = false // Sadece debug için!
/// ```
///
/// ## Avantajlar
///
/// - **Hızlı Geliştirme**: HTTP üzerinde çalışır, HTTPS gerekmez
/// - **Esnek CORS**: Yaygın frontend port'larına izin verir
/// - **Yüksek Rate Limit**: Test sırasında engellenmez
/// - **Console Logging**: Terminal'de anında görüntüleme
/// - **Kolay Debug**: Daha az kısıtlama, daha kolay hata ayıklama
///
/// ## Dezavantajlar
///
/// - **Düşük Güvenlik**: Production için uygun değil
/// - **HTTP**: Man-in-the-middle saldırılarına açık
/// - **Esnek Limitler**: Brute-force saldırılarına karşı zayıf
/// - **Geçici Log'lar**: Console log'ları kalıcı değil
/// - **Test Key'ler**: Gerçek veri şifrelemesi için uygun değil
///
/// ## Production'a Geçiş Checklist
///
/// Geliştirme tamamlandığında production'a geçmeden önce:
///
/// - [ ] `ProductionSecurityConfig()` kullan
/// - [ ] CORS AllowedOrigins'i gerçek domain'lerle değiştir
/// - [ ] Encryption key'i environment variable'dan oku
/// - [ ] HTTPS sertifikası yapılandır
/// - [ ] Audit logging'i dosyaya veya SIEM'e yönlendir
/// - [ ] Rate limiting ayarlarını gözden geçir
/// - [ ] Session Secure=true yap
/// - [ ] Test key'lerini production key'leriyle değiştir
/// - [ ] Environment variable'ları production'da ayarla
/// - [ ] Log rotation yapılandır
/// - [ ] Monitoring ve alerting kur
///
/// ## İlgili Fonksiyonlar
///
/// - `DefaultSecurityConfig()`: Temel güvenlik ayarları
/// - `ProductionSecurityConfig()`: Production için sıkı ayarlar
///
/// ## Döndürür
///
/// - Development ortamı için optimize edilmiş `SecurityConfig` struct'ı
/// - HTTP üzerinde çalışabilir (Secure=false)
/// - Esnek CORS ayarları (localhost port'ları)
/// - Yüksek rate limit (test için uygun)
/// - Console logging (terminal'de görüntüleme)
func DevelopmentSecurityConfig() SecurityConfig {
	config := DefaultSecurityConfig()

	// Development-specific overrides (still secure!)
	config.CORS.AllowedOrigins = []string{"http://localhost:3000", "http://localhost:5173"}
	config.Session.Secure = false // Allow HTTP in development
	config.Session.SameSite = "Lax"
	config.RateLimit.AuthMaxRequests = 50 // More lenient for development
	config.Audit.Destination = "console"

	return config
}
