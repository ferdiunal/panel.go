package middleware

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

/// # AuditEvent
///
/// Bu yapı, güvenlik denetim olaylarını temsil eder ve sistemdeki tüm önemli
/// güvenlik olaylarının kaydedilmesi için kullanılır.
///
/// ## Kullanım Senaryoları
///
/// - Kullanıcı kimlik doğrulama olaylarının kaydedilmesi (giriş, çıkış, kayıt)
/// - API isteklerinin ve yanıtlarının izlenmesi
/// - Kaynak erişim ve değişiklik olaylarının takibi
/// - İzin kontrolü sonuçlarının loglanması
/// - Güvenlik ihlali girişimlerinin tespiti
/// - Uyumluluk (compliance) raporlaması için veri toplama
///
/// ## Alanlar
///
/// - `Timestamp`: Olayın gerçekleştiği zaman damgası
/// - `EventType`: Olay tipi (login_success, login_failure, resource_create, vb.)
/// - `UserID`: İşlemi gerçekleştiren kullanıcının ID'si (opsiyonel)
/// - `Email`: Kullanıcının e-posta adresi (opsiyonel)
/// - `IP`: İsteğin geldiği IP adresi
/// - `UserAgent`: İstemci tarayıcı/uygulama bilgisi
/// - `Method`: HTTP metodu (GET, POST, PUT, DELETE, vb.)
/// - `Path`: İstek yapılan URL yolu
/// - `StatusCode`: HTTP yanıt durum kodu
/// - `Success`: İşlemin başarılı olup olmadığı (status code < 400)
/// - `ErrorMessage`: Hata mesajı (varsa)
/// - `Resource`: Erişilen kaynak adı (opsiyonel)
/// - `Action`: Gerçekleştirilen aksiyon (read, create, update, delete, vb.)
/// - `Metadata`: Ek bilgiler için anahtar-değer çiftleri
///
/// ## Örnek Kullanım
///
/// ```go
/// event := AuditEvent{
///     Timestamp:  time.Now(),
///     EventType:  "login_success",
///     Email:      "user@example.com",
///     IP:         "192.168.1.1",
///     UserAgent:  "Mozilla/5.0...",
///     Success:    true,
/// }
/// ```
///
/// ## Önemli Notlar
///
/// - JSON serileştirme için struct tag'leri kullanılır
/// - `omitempty` tag'i ile boş alanlar JSON'da gösterilmez
/// - Zaman damgası her zaman UTC formatında saklanmalıdır
/// - IP adresi gizlilik düzenlemelerine uygun şekilde işlenmelidir
type AuditEvent struct {
	Timestamp    time.Time              `json:"timestamp"`
	EventType    string                 `json:"event_type"`
	UserID       string                 `json:"user_id,omitempty"`
	Email        string                 `json:"email,omitempty"`
	IP           string                 `json:"ip"`
	UserAgent    string                 `json:"user_agent"`
	Method       string                 `json:"method"`
	Path         string                 `json:"path"`
	StatusCode   int                    `json:"status_code"`
	Success      bool                   `json:"success"`
	ErrorMessage string                 `json:"error_message,omitempty"`
	Resource     string                 `json:"resource,omitempty"`
	Action       string                 `json:"action,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
}

/// # AuditLogger
///
/// Bu interface, denetim loglama implementasyonları için standart bir arayüz sağlar.
/// Farklı loglama stratejileri (konsol, dosya, veritabanı, uzak servis) bu interface'i
/// implement ederek sistemde kullanılabilir.
///
/// ## Kullanım Senaryoları
///
/// - Geliştirme ortamında konsola loglama (ConsoleAuditLogger)
/// - Production ortamında dosyaya loglama (FileAuditLogger)
/// - Merkezi log toplama sistemlerine gönderme (Elasticsearch, Splunk, vb.)
/// - Veritabanına kaydetme (PostgreSQL, MongoDB, vb.)
/// - Güvenlik bilgi ve olay yönetimi (SIEM) sistemlerine entegrasyon
///
/// ## Metodlar
///
/// - `Log(event AuditEvent) error`: Denetim olayını loglar
///
/// ## Örnek Kullanım
///
/// ```go
/// // Özel bir logger implementasyonu
/// type DatabaseAuditLogger struct {
///     db *gorm.DB
/// }
///
/// func (l *DatabaseAuditLogger) Log(event AuditEvent) error {
///     return l.db.Create(&event).Error
/// }
///
/// // Middleware'de kullanım
/// logger := &DatabaseAuditLogger{db: db}
/// app.Use(AuditMiddleware(logger))
/// ```
///
/// ## Avantajlar
///
/// - **Esneklik**: Farklı loglama stratejileri kolayca değiştirilebilir
/// - **Test Edilebilirlik**: Mock logger'lar ile test yazmak kolay
/// - **Genişletilebilirlik**: Yeni loglama hedefleri eklemek basit
/// - **Bağımsızlık**: Loglama mantığı iş mantığından ayrılmış
///
/// ## Önemli Notlar
///
/// - Log metodu asenkron olarak çağrılabilir, thread-safe olmalıdır
/// - Hata durumunda panic atmak yerine error döndürülmelidir
/// - Performans için buffering ve batch processing kullanılabilir
/// - Hassas bilgiler (şifreler, tokenlar) loglanmamalıdır
type AuditLogger interface {
	Log(event AuditEvent) error
}

/// # ConsoleAuditLogger
///
/// Bu yapı, denetim olaylarını konsola yazdıran basit bir logger implementasyonudur.
/// Geliştirme ortamında hızlı debugging ve test için idealdir.
///
/// ## Kullanım Senaryoları
///
/// - Geliştirme ortamında gerçek zamanlı log izleme
/// - Hızlı debugging ve sorun giderme
/// - Test ortamında log çıktılarını görme
/// - CI/CD pipeline'larında log toplama
/// - Docker container loglarını stdout'a yönlendirme
///
/// ## Örnek Kullanım
///
/// ```go
/// // Basit kullanım
/// logger := &ConsoleAuditLogger{}
/// app.Use(AuditMiddleware(logger))
///
/// // Veya nil geçerek otomatik olarak ConsoleAuditLogger kullanımı
/// app.Use(AuditMiddleware(nil))
/// ```
///
/// ## Avantajlar
///
/// - **Basitlik**: Kurulum gerektirmez, hemen kullanılabilir
/// - **Hız**: Dosya I/O yok, çok hızlı
/// - **Debugging**: Gerçek zamanlı log görüntüleme
/// - **Container Uyumlu**: Docker/Kubernetes logları için ideal
///
/// ## Dezavantajlar
///
/// - **Kalıcılık Yok**: Loglar saklanmaz, kaybolur
/// - **Arama Zorluğu**: Geçmiş logları aramak imkansız
/// - **Performans**: Yüksek trafikte konsol çıktısı yavaşlatabilir
/// - **Production Uygunsuz**: Production ortamı için önerilmez
///
/// ## Önemli Notlar
///
/// - Production ortamında FileAuditLogger veya uzak log servisi kullanın
/// - JSON formatında çıktı verir, log parsing araçları ile uyumlu
/// - Thread-safe değildir, yüksek concurrency'de karışık çıktı olabilir
/// - Hassas bilgiler konsola yazılır, dikkatli kullanın
type ConsoleAuditLogger struct{}

/// # Log
///
/// Bu fonksiyon, denetim olayını JSON formatında konsola yazdırır.
///
/// ## Parametreler
///
/// - `event`: Loglanacak AuditEvent yapısı
///
/// ## Dönüş Değeri
///
/// - `error`: JSON serileştirme hatası (nadiren oluşur), başarılı ise nil
///
/// ## Davranış
///
/// 1. AuditEvent'i JSON formatına dönüştürür
/// 2. "[AUDIT]" prefix'i ile konsola yazdırır
/// 3. Hata durumunda error döndürür
///
/// ## Örnek Çıktı
///
/// ```json
/// [AUDIT] {"timestamp":"2024-01-15T10:30:00Z","event_type":"login_success","email":"user@example.com","ip":"192.168.1.1","success":true}
/// ```
///
/// ## Önemli Notlar
///
/// - Blocking operation, yüksek trafikte performans etkisi olabilir
/// - JSON serileştirme hatası çok nadir, genellikle nil döner
/// - Konsol çıktısı buffer'lanabilir, gerçek zamanlı olmayabilir
func (l *ConsoleAuditLogger) Log(event AuditEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	fmt.Printf("[AUDIT] %s\n", string(data))
	return nil
}

/// # FileAuditLogger
///
/// Bu yapı, denetim olaylarını dosyaya yazan production-ready bir logger implementasyonudur.
/// Log rotation, buffering ve performans optimizasyonları ile production ortamı için tasarlanmıştır.
///
/// ## Kullanım Senaryoları
///
/// - Production ortamında kalıcı log saklama
/// - Uyumluluk (compliance) gereksinimleri için log arşivleme
/// - Güvenlik denetimi ve forensic analiz
/// - Log aggregation sistemlerine veri besleme
/// - Uzun vadeli log saklama ve analiz
///
/// ## Planlanan Özellikler (TODO)
///
/// - **Log Rotation**: Dosya boyutu veya zaman bazlı otomatik rotation
/// - **Compression**: Eski logların otomatik sıkıştırılması
/// - **Buffering**: Performans için batch yazma
/// - **Async Writing**: Non-blocking asenkron yazma
/// - **Retention Policy**: Eski logların otomatik silinmesi
/// - **Multiple Files**: Olay tipine göre farklı dosyalara yazma
///
/// ## Örnek Kullanım (Gelecek)
///
/// ```go
/// logger := &FileAuditLogger{
///     FilePath:       "/var/log/audit/app.log",
///     MaxSize:        100, // MB
///     MaxBackups:     10,
///     MaxAge:         30, // days
///     Compress:       true,
/// }
/// app.Use(AuditMiddleware(logger))
/// ```
///
/// ## Avantajlar
///
/// - **Kalıcılık**: Loglar disk üzerinde saklanır
/// - **Arama**: Geçmiş logları arama ve analiz imkanı
/// - **Uyumluluk**: Yasal gereksinimleri karşılar
/// - **Forensics**: Güvenlik olaylarını araştırma
///
/// ## Dezavantajlar
///
/// - **Disk Kullanımı**: Disk alanı gerektirir
/// - **I/O Overhead**: Dosya yazma performans etkisi
/// - **Yönetim**: Log rotation ve temizleme gerektirir
///
/// ## Önemli Notlar
///
/// - **ŞU ANDA IMPLEMENT EDİLMEMİŞTİR** - TODO olarak işaretlenmiştir
/// - Production kullanımı için önce implement edilmesi gerekir
/// - Alternatif olarak üçüncü parti kütüphaneler kullanılabilir (lumberjack, zap, vb.)
/// - Thread-safe implementasyon gereklidir
type FileAuditLogger struct {
	// TODO: Implement file-based logging
}

/// # Log
///
/// Bu fonksiyon, denetim olayını dosyaya yazmak için kullanılacaktır.
/// **ŞU ANDA IMPLEMENT EDİLMEMİŞTİR** - Placeholder implementasyon.
///
/// ## Parametreler
///
/// - `event`: Loglanacak AuditEvent yapısı
///
/// ## Dönüş Değeri
///
/// - `error`: Şu anda her zaman nil döner (TODO implementasyonu)
///
/// ## Planlanan Davranış
///
/// 1. AuditEvent'i JSON formatına dönüştürme
/// 2. Dosyaya thread-safe şekilde yazma
/// 3. Buffer'ı periyodik olarak flush etme
/// 4. Gerekirse log rotation tetikleme
/// 5. Hata durumunda error döndürme
///
/// ## Önemli Notlar
///
/// - **KULLANMAYIN**: Bu implementasyon henüz tamamlanmamıştır
/// - Production'da kullanmadan önce implement edilmesi gerekir
/// - Geçici olarak ConsoleAuditLogger kullanın
///
/// ## TODO
///
/// - Dosya açma ve yazma mantığı
/// - Log rotation implementasyonu
/// - Buffer yönetimi
/// - Error handling
/// - Thread-safety (mutex kullanımı)
func (l *FileAuditLogger) Log(event AuditEvent) error {
	// TODO: Implement file writing with rotation
	return nil
}

/// # AuditMiddleware
///
/// Bu fonksiyon, tüm HTTP isteklerini otomatik olarak loglayan bir Fiber middleware'i oluşturur.
/// Her istek için detaylı denetim kaydı tutar ve güvenlik olaylarını izler.
///
/// ## Kullanım Senaryoları
///
/// - Tüm API isteklerinin merkezi loglanması
/// - Güvenlik denetimi ve uyumluluk gereksinimleri
/// - Kullanıcı aktivitelerinin izlenmesi
/// - Performans analizi ve debugging
/// - Güvenlik ihlali tespiti ve forensic analiz
/// - API kullanım istatistikleri toplama
///
/// ## Parametreler
///
/// - `logger`: AuditLogger interface'ini implement eden logger instance'ı
///   - `nil` geçilirse otomatik olarak ConsoleAuditLogger kullanılır
///
/// ## Dönüş Değeri
///
/// - `fiber.Handler`: Fiber framework'ü ile uyumlu middleware handler fonksiyonu
///
/// ## Davranış
///
/// 1. İstek başlangıç zamanını kaydeder
/// 2. İsteği bir sonraki handler'a iletir (c.Next())
/// 3. İstek tamamlandıktan sonra AuditEvent oluşturur
/// 4. HTTP metodu, path, IP, User-Agent gibi bilgileri toplar
/// 5. Context'ten kullanıcı bilgilerini çıkarır (varsa)
/// 6. İstek tipini otomatik olarak belirler (determineEventType)
/// 7. Olayı logger'a gönderir
/// 8. Loglama hatalarını konsola yazdırır (uygulamayı durdurmaz)
///
/// ## Örnek Kullanım
///
/// ```go
/// // Basit kullanım (ConsoleAuditLogger)
/// app := fiber.New()
/// app.Use(AuditMiddleware(nil))
///
/// // Özel logger ile kullanım
/// logger := &DatabaseAuditLogger{db: db}
/// app.Use(AuditMiddleware(logger))
///
/// // Sadece belirli route'lar için
/// api := app.Group("/api")
/// api.Use(AuditMiddleware(logger))
/// ```
///
/// ## Loglanan Bilgiler
///
/// - **Timestamp**: İstek başlangıç zamanı
/// - **Method**: HTTP metodu (GET, POST, PUT, DELETE, vb.)
/// - **Path**: İstek URL'i
/// - **IP**: İstemci IP adresi
/// - **UserAgent**: Tarayıcı/uygulama bilgisi
/// - **StatusCode**: HTTP yanıt kodu
/// - **Success**: İşlem başarılı mı (status < 400)
/// - **EventType**: Otomatik belirlenen olay tipi
/// - **User**: Context'ten çıkarılan kullanıcı bilgisi (varsa)
///
/// ## Olay Tipleri
///
/// Middleware otomatik olarak şu olay tiplerini belirler:
/// - `login_success` / `login_failure`: Giriş işlemleri
/// - `registration`: Kayıt işlemleri
/// - `logout`: Çıkış işlemleri
/// - `password_reset_request`: Şifre sıfırlama
/// - `resource_read` / `resource_create` / `resource_update` / `resource_delete`: CRUD işlemleri
/// - `settings_change`: Ayar değişiklikleri
/// - `api_request`: Genel API istekleri
///
/// ## Performans Notları
///
/// - Middleware non-blocking şekilde çalışır
/// - Loglama hatası uygulamayı durdurmaz
/// - Asenkron loglama için özel logger implementasyonu kullanılabilir
/// - Yüksek trafikte buffering önerilir
///
/// ## Güvenlik Notları
///
/// - IP adresleri GDPR/KVKK uyumlu şekilde işlenmelidir
/// - Hassas bilgiler (şifreler, tokenlar) loglanmamalıdır
/// - User-Agent bilgisi fingerprinting için kullanılabilir
/// - Log dosyaları güvenli bir şekilde saklanmalıdır
///
/// ## Önemli Notlar
///
/// - Middleware sırası önemlidir, genellikle en başta kullanılmalıdır
/// - Context'e "user" key'i ile kullanıcı bilgisi eklenmelidir
/// - Loglama hatası uygulamayı durdurmaz, sadece konsola yazdırılır
/// - Her istek için bir AuditEvent oluşturulur
func AuditMiddleware(logger AuditLogger) fiber.Handler {
	if logger == nil {
		logger = &ConsoleAuditLogger{}
	}

	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		// Log after request completes
		event := AuditEvent{
			Timestamp:  start,
			Method:     c.Method(),
			Path:       c.Path(),
			IP:         c.IP(),
			UserAgent:  c.Get("User-Agent"),
			StatusCode: c.Response().StatusCode(),
			Success:    c.Response().StatusCode() < 400,
		}

		// Extract user info from context if available
		if user := c.Locals("user"); user != nil {
			// Type assertion to extract user details
			// This depends on your user struct
			event.Metadata = map[string]interface{}{
				"user": user,
			}
		}

		// Determine event type based on path and method
		event.EventType = determineEventType(c.Method(), c.Path(), c.Response().StatusCode())

		// Log the event
		if logErr := logger.Log(event); logErr != nil {
			fmt.Printf("[AUDIT ERROR] Failed to log event: %v\n", logErr)
		}

		return err
	}
}

/// # LogAuthEvent
///
/// Bu fonksiyon, kimlik doğrulama ile ilgili olayları loglamak için kullanılır.
/// Giriş, çıkış, kayıt ve şifre sıfırlama gibi authentication işlemlerini kaydeder.
///
/// ## Kullanım Senaryoları
///
/// - Kullanıcı giriş denemelerinin kaydedilmesi (başarılı/başarısız)
/// - Kullanıcı çıkış işlemlerinin loglanması
/// - Yeni kullanıcı kayıtlarının takibi
/// - Şifre sıfırlama isteklerinin kaydedilmesi
/// - Brute force saldırı tespiti
/// - Şüpheli giriş aktivitelerinin izlenmesi
/// - Uyumluluk raporlaması için authentication logları
///
/// ## Parametreler
///
/// - `logger`: AuditLogger interface'ini implement eden logger instance'ı
/// - `eventType`: Olay tipi (login_success, login_failure, registration, logout, vb.)
/// - `email`: Kullanıcının e-posta adresi
/// - `ip`: İsteğin geldiği IP adresi
/// - `userAgent`: İstemci tarayıcı/uygulama bilgisi
/// - `success`: İşlemin başarılı olup olmadığı
/// - `errorMsg`: Hata mesajı (başarısız işlemler için), boş string başarılı işlemler için
///
/// ## Dönüş Değeri
///
/// Fonksiyon değer döndürmez (void). Loglama hatası durumunda konsola yazdırır.
///
/// ## Örnek Kullanım
///
/// ```go
/// // Başarılı giriş
/// LogAuthEvent(
///     logger,
///     "login_success",
///     "user@example.com",
///     "192.168.1.1",
///     "Mozilla/5.0...",
///     true,
///     "",
/// )
///
/// // Başarısız giriş
/// LogAuthEvent(
///     logger,
///     "login_failure",
///     "user@example.com",
///     "192.168.1.1",
///     "Mozilla/5.0...",
///     false,
///     "Invalid credentials",
/// )
///
/// // Kullanıcı kaydı
/// LogAuthEvent(
///     logger,
///     "registration",
///     "newuser@example.com",
///     "192.168.1.1",
///     "Mozilla/5.0...",
///     true,
///     "",
/// )
/// ```
///
/// ## Olay Tipleri
///
/// - `login_success`: Başarılı giriş
/// - `login_failure`: Başarısız giriş
/// - `registration`: Yeni kullanıcı kaydı
/// - `logout`: Kullanıcı çıkışı
/// - `password_reset_request`: Şifre sıfırlama isteği
/// - `password_reset_complete`: Şifre sıfırlama tamamlandı
/// - `email_verification`: E-posta doğrulama
/// - `two_factor_auth`: İki faktörlü kimlik doğrulama
///
/// ## Güvenlik Notları
///
/// - Şifreler asla loglanmamalıdır
/// - IP adresleri GDPR/KVKK uyumlu şekilde işlenmelidir
/// - Başarısız giriş denemeleri rate limiting için kullanılabilir
/// - Aynı IP'den çok sayıda başarısız deneme brute force saldırısı olabilir
/// - E-posta adresleri hassas veri olarak kabul edilmelidir
///
/// ## Önemli Notlar
///
/// - Fonksiyon non-blocking şekilde çalışır
/// - Loglama hatası uygulamayı durdurmaz
/// - Zaman damgası otomatik olarak eklenir (time.Now())
/// - Error mesajı sadece başarısız işlemler için doldurulmalıdır
func LogAuthEvent(logger AuditLogger, eventType, email, ip, userAgent string, success bool, errorMsg string) {
	event := AuditEvent{
		Timestamp:    time.Now(),
		EventType:    eventType,
		Email:        email,
		IP:           ip,
		UserAgent:    userAgent,
		Success:      success,
		ErrorMessage: errorMsg,
	}

	if err := logger.Log(event); err != nil {
		fmt.Printf("[AUDIT ERROR] Failed to log auth event: %v\n", err)
	}
}

/// # LogPermissionCheck
///
/// Bu fonksiyon, izin kontrolü olaylarını loglamak için kullanılır.
/// Kullanıcıların kaynaklara erişim yetkilerinin kontrol edilmesi ve sonuçlarının kaydedilmesi için tasarlanmıştır.
///
/// ## Kullanım Senaryoları
///
/// - Rol tabanlı erişim kontrolü (RBAC) sonuçlarının loglanması
/// - Yetki ihlali girişimlerinin tespiti
/// - Uyumluluk raporlaması için erişim kontrol kayıtları
/// - Güvenlik denetimi ve forensic analiz
/// - Yetkisiz erişim denemelerinin izlenmesi
/// - İzin politikalarının etkinliğinin değerlendirilmesi
/// - Kullanıcı davranış analizi
///
/// ## Parametreler
///
/// - `logger`: AuditLogger interface'ini implement eden logger instance'ı
/// - `userID`: İzin kontrolü yapılan kullanıcının ID'si
/// - `resource`: Erişilmeye çalışılan kaynak adı (örn: "users", "posts", "settings")
/// - `action`: Gerçekleştirilmeye çalışılan aksiyon (örn: "read", "create", "update", "delete")
/// - `granted`: İznin verilip verilmediği (true: izin verildi, false: reddedildi)
///
/// ## Dönüş Değeri
///
/// Fonksiyon değer döndürmez (void). Loglama hatası durumunda konsola yazdırır.
///
/// ## Örnek Kullanım
///
/// ```go
/// // İzin verildi
/// LogPermissionCheck(
///     logger,
///     "user-123",
///     "posts",
///     "create",
///     true,
/// )
///
/// // İzin reddedildi
/// LogPermissionCheck(
///     logger,
///     "user-456",
///     "admin_settings",
///     "update",
///     false,
/// )
///
/// // Middleware içinde kullanım
/// if !hasPermission(userID, resource, action) {
///     LogPermissionCheck(logger, userID, resource, action, false)
///     return fiber.NewError(fiber.StatusForbidden, "Permission denied")
/// }
/// LogPermissionCheck(logger, userID, resource, action, true)
/// ```
///
/// ## Kaynak ve Aksiyon Örnekleri
///
/// **Kaynaklar:**
/// - `users`: Kullanıcı yönetimi
/// - `posts`: İçerik yönetimi
/// - `settings`: Sistem ayarları
/// - `reports`: Raporlar
/// - `admin_panel`: Yönetim paneli
///
/// **Aksiyonlar:**
/// - `read`: Okuma/görüntüleme
/// - `create`: Oluşturma
/// - `update`: Güncelleme
/// - `delete`: Silme
/// - `export`: Dışa aktarma
/// - `import`: İçe aktarma
///
/// ## Güvenlik Notları
///
/// - Reddedilen izinler özellikle dikkatle izlenmelidir
/// - Aynı kullanıcıdan çok sayıda reddedilen izin şüpheli aktivite olabilir
/// - İzin kontrolleri her kritik işlemden önce yapılmalıdır
/// - Bypass girişimleri tespit edilmeli ve loglanmalıdır
/// - RBAC politikaları düzenli olarak gözden geçirilmelidir
///
/// ## Analiz ve Raporlama
///
/// Bu loglar şunlar için kullanılabilir:
/// - En çok reddedilen kaynakların belirlenmesi
/// - Kullanıcı davranış paternlerinin analizi
/// - Güvenlik ihlali girişimlerinin tespiti
/// - İzin politikalarının optimizasyonu
/// - Uyumluluk raporlarının oluşturulması
///
/// ## Önemli Notlar
///
/// - EventType otomatik olarak "permission_check" olarak ayarlanır
/// - Zaman damgası otomatik olarak eklenir (time.Now())
/// - Success alanı granted parametresi ile doldurulur
/// - Fonksiyon non-blocking şekilde çalışır
/// - Loglama hatası uygulamayı durdurmaz
func LogPermissionCheck(logger AuditLogger, userID, resource, action string, granted bool) {
	event := AuditEvent{
		Timestamp: time.Now(),
		EventType: "permission_check",
		UserID:    userID,
		Resource:  resource,
		Action:    action,
		Success:   granted,
	}

	if err := logger.Log(event); err != nil {
		fmt.Printf("[AUDIT ERROR] Failed to log permission check: %v\n", err)
	}
}

/// # LogDataAccess
///
/// Bu fonksiyon, veri erişim olaylarını loglamak için kullanılır.
/// Kullanıcıların veritabanı kayıtlarına erişimlerini ve veri işlemlerini takip eder.
///
/// ## Kullanım Senaryoları
///
/// - Hassas verilere erişim kayıtlarının tutulması
/// - CRUD işlemlerinin detaylı loglanması
/// - Veri değişikliklerinin izlenmesi (audit trail)
/// - Uyumluluk gereksinimleri için veri erişim kayıtları (GDPR, KVKK, HIPAA)
/// - Veri sızıntısı araştırmaları için forensic analiz
/// - Kullanıcı aktivite raporlaması
/// - Veri erişim paternlerinin analizi
///
/// ## Parametreler
///
/// - `logger`: AuditLogger interface'ini implement eden logger instance'ı
/// - `userID`: Veri erişimi yapan kullanıcının ID'si
/// - `resource`: Erişilen kaynak/tablo adı (örn: "users", "orders", "payments")
/// - `action`: Gerçekleştirilen işlem (örn: "read", "create", "update", "delete", "export")
/// - `recordID`: Erişilen kaydın benzersiz ID'si
///
/// ## Dönüş Değeri
///
/// Fonksiyon değer döndürmez (void). Loglama hatası durumunda konsola yazdırır.
///
/// ## Örnek Kullanım
///
/// ```go
/// // Kullanıcı kaydı okuma
/// LogDataAccess(
///     logger,
///     "user-123",
///     "users",
///     "read",
///     "user-456",
/// )
///
/// // Sipariş güncelleme
/// LogDataAccess(
///     logger,
///     "admin-789",
///     "orders",
///     "update",
///     "order-12345",
/// )
///
/// // Ödeme bilgisi silme
/// LogDataAccess(
///     logger,
///     "user-123",
///     "payments",
///     "delete",
///     "payment-67890",
/// )
///
/// // CRUD handler içinde kullanım
/// func GetUser(c *fiber.Ctx) error {
///     userID := c.Params("id")
///     user, err := db.FindUser(userID)
///     if err != nil {
///         return err
///     }
///
///     // Veri erişimini logla
///     LogDataAccess(
///         logger,
///         c.Locals("currentUserID").(string),
///         "users",
///         "read",
///         userID,
///     )
///
///     return c.JSON(user)
/// }
/// ```
///
/// ## Kaynak ve Aksiyon Örnekleri
///
/// **Kaynaklar:**
/// - `users`: Kullanıcı verileri
/// - `orders`: Sipariş kayıtları
/// - `payments`: Ödeme bilgileri
/// - `medical_records`: Tıbbi kayıtlar (HIPAA)
/// - `personal_data`: Kişisel veriler (GDPR/KVKK)
/// - `financial_data`: Finansal veriler
///
/// **Aksiyonlar:**
/// - `read`: Kayıt okuma/görüntüleme
/// - `create`: Yeni kayıt oluşturma
/// - `update`: Kayıt güncelleme
/// - `delete`: Kayıt silme
/// - `export`: Veri dışa aktarma
/// - `import`: Veri içe aktarma
/// - `bulk_delete`: Toplu silme
/// - `anonymize`: Veri anonimleştirme
///
/// ## Metadata Alanı
///
/// RecordID, event.Metadata içinde saklanır:
/// ```go
/// Metadata: map[string]interface{}{
///     "record_id": recordID,
/// }
/// ```
///
/// Ek bilgiler de metadata'ya eklenebilir:
/// - Değişen alanlar (before/after)
/// - İşlem süresi
/// - Etkilenen kayıt sayısı
/// - İlişkili kayıtlar
///
/// ## Uyumluluk ve Yasal Gereksinimler
///
/// Bu loglar şu düzenlemelere uyum için kritiktir:
/// - **GDPR**: Kişisel verilere erişim kayıtları
/// - **KVKK**: Türkiye veri koruma kanunu
/// - **HIPAA**: Sağlık verilerine erişim kayıtları
/// - **PCI-DSS**: Ödeme kartı verilerine erişim
/// - **SOX**: Finansal veri erişim kayıtları
///
/// ## Güvenlik Notları
///
/// - Hassas verilere erişim özellikle dikkatle izlenmelidir
/// - Toplu veri indirme işlemleri şüpheli olabilir
/// - Çalışma saatleri dışı erişimler incelenmelidir
/// - Aynı kullanıcıdan çok sayıda erişim anormal olabilir
/// - Silinen kayıtlar geri alınamaz, dikkatli loglanmalıdır
///
/// ## Analiz ve Raporlama
///
/// Bu loglar şunlar için kullanılabilir:
/// - En çok erişilen kaynakların belirlenmesi
/// - Kullanıcı davranış paternlerinin analizi
/// - Veri sızıntısı tespiti
/// - Uyumluluk raporlarının oluşturulması
/// - Performans optimizasyonu (sık erişilen veriler)
///
/// ## Önemli Notlar
///
/// - EventType otomatik olarak "data_access" olarak ayarlanır
/// - Zaman damgası otomatik olarak eklenir (time.Now())
/// - RecordID metadata içinde saklanır
/// - Fonksiyon non-blocking şekilde çalışır
/// - Loglama hatası uygulamayı durdurmaz
/// - Hassas veri içeriği değil, sadece erişim bilgisi loglanmalıdır
func LogDataAccess(logger AuditLogger, userID, resource, action string, recordID string) {
	event := AuditEvent{
		Timestamp: time.Now(),
		EventType: "data_access",
		UserID:    userID,
		Resource:  resource,
		Action:    action,
		Metadata: map[string]interface{}{
			"record_id": recordID,
		},
	}

	if err := logger.Log(event); err != nil {
		fmt.Printf("[AUDIT ERROR] Failed to log data access: %v\n", err)
	}
}

/// # determineEventType
///
/// Bu fonksiyon, HTTP istek detaylarına göre olay tipini otomatik olarak belirler.
/// URL path'i, HTTP metodu ve yanıt durum kodunu analiz ederek uygun olay tipini döndürür.
///
/// ## Kullanım Senaryoları
///
/// - AuditMiddleware içinde otomatik olay tipi belirleme
/// - İstek loglarını kategorize etme
/// - Güvenlik olaylarını sınıflandırma
/// - Raporlama için olay tiplerini standardize etme
/// - İstatistiksel analiz için olay gruplandırma
///
/// ## Parametreler
///
/// - `method`: HTTP metodu (GET, POST, PUT, PATCH, DELETE, vb.)
/// - `path`: İstek URL path'i (örn: "/auth/sign-in", "/resource/users")
/// - `statusCode`: HTTP yanıt durum kodu (200, 401, 404, 500, vb.)
///
/// ## Dönüş Değeri
///
/// - `string`: Belirlenen olay tipi
///
/// ## Belirlenen Olay Tipleri
///
/// ### Kimlik Doğrulama Olayları
/// - `login_success`: Başarılı giriş (path: /auth/sign-in, status < 400)
/// - `login_failure`: Başarısız giriş (path: /auth/sign-in, status >= 400)
/// - `registration`: Kullanıcı kaydı (path: /auth/sign-up)
/// - `logout`: Kullanıcı çıkışı (path: /auth/sign-out)
/// - `password_reset_request`: Şifre sıfırlama isteği (path: /auth/forgot-password)
///
/// ### Kaynak İşlemleri (CRUD)
/// - `resource_read`: Kaynak okuma (path: /resource/*, method: GET)
/// - `resource_create`: Kaynak oluşturma (path: /resource/*, method: POST)
/// - `resource_update`: Kaynak güncelleme (path: /resource/*, method: PUT/PATCH)
/// - `resource_delete`: Kaynak silme (path: /resource/*, method: DELETE)
///
/// ### Diğer Olaylar
/// - `settings_change`: Ayar değişikliği (path: /pages/settings)
/// - `api_request`: Genel API isteği (diğer tüm durumlar)
///
/// ## Örnek Kullanım
///
/// ```go
/// // Başarılı giriş
/// eventType := determineEventType("POST", "/auth/sign-in", 200)
/// // Sonuç: "login_success"
///
/// // Başarısız giriş
/// eventType := determineEventType("POST", "/auth/sign-in", 401)
/// // Sonuç: "login_failure"
///
/// // Kaynak oluşturma
/// eventType := determineEventType("POST", "/resource/users", 201)
/// // Sonuç: "resource_create"
///
/// // Kaynak okuma
/// eventType := determineEventType("GET", "/resource/posts", 200)
/// // Sonuç: "resource_read"
///
/// // Genel API isteği
/// eventType := determineEventType("GET", "/api/dashboard", 200)
/// // Sonuç: "api_request"
/// ```
///
/// ## Algoritma
///
/// 1. Önce authentication endpoint'lerini kontrol eder
/// 2. Sonra resource endpoint'lerini kontrol eder
/// 3. Settings endpoint'lerini kontrol eder
/// 4. Hiçbiri eşleşmezse "api_request" döndürür
///
/// ## Path Eşleştirme
///
/// Path eşleştirme için `contains` fonksiyonu kullanılır:
/// - Başlangıç eşleştirme: "/auth/sign-in" ile başlıyor mu?
/// - Bitiş eşleştirme: "/auth/sign-in" ile bitiyor mu?
/// - İçerik eşleştirme: "/resource/" içeriyor mu?
///
/// ## Genişletme
///
/// Yeni olay tipleri eklemek için:
/// ```go
/// // Yeni endpoint kontrolü ekle
/// if contains(path, "/api/export") {
///     return "data_export"
/// }
/// ```
///
/// ## Önemli Notlar
///
/// - Fonksiyon case-sensitive çalışır
/// - Path'ler tam eşleşme veya içerme ile kontrol edilir
/// - Status code sadece login için kontrol edilir
/// - Öncelik sırası önemlidir (auth > resource > settings > default)
/// - Özel endpoint'ler için fonksiyon genişletilebilir
func determineEventType(method, path string, statusCode int) string {
	// Authentication endpoints
	if contains(path, "/auth/sign-in") {
		if statusCode < 400 {
			return "login_success"
		}
		return "login_failure"
	}
	if contains(path, "/auth/sign-up") {
		return "registration"
	}
	if contains(path, "/auth/sign-out") {
		return "logout"
	}
	if contains(path, "/auth/forgot-password") {
		return "password_reset_request"
	}

	// Resource operations
	if contains(path, "/resource/") {
		switch method {
		case "GET":
			return "resource_read"
		case "POST":
			return "resource_create"
		case "PUT", "PATCH":
			return "resource_update"
		case "DELETE":
			return "resource_delete"
		}
	}

	// Settings operations
	if contains(path, "/pages/settings") {
		return "settings_change"
	}

	return "api_request"
}

/// # contains
///
/// Bu fonksiyon, bir string'in başında, sonunda veya içinde belirli bir substring'in
/// bulunup bulunmadığını kontrol eder. Üç farklı eşleştirme stratejisi kullanır.
///
/// ## Kullanım Senaryoları
///
/// - URL path eşleştirme (determineEventType içinde)
/// - Endpoint kontrolü ve sınıflandırma
/// - String içerik arama
/// - Path pattern matching
///
/// ## Parametreler
///
/// - `s`: Aranacak ana string
/// - `substr`: Aranacak substring
///
/// ## Dönüş Değeri
///
/// - `bool`: Substring bulunursa true, bulunamazsa false
///
/// ## Eşleştirme Stratejileri
///
/// 1. **Başlangıç Eşleştirme**: String, substring ile başlıyor mu?
///    - Örnek: "/auth/sign-in" içinde "/auth" başlangıçta var
///
/// 2. **Bitiş Eşleştirme**: String, substring ile bitiyor mu?
///    - Örnek: "/api/users/profile" içinde "/profile" sonda var
///
/// 3. **İçerik Eşleştirme**: String, substring'i içeriyor mu? (findSubstring ile)
///    - Örnek: "/api/resource/users" içinde "/resource/" var
///
/// ## Örnek Kullanım
///
/// ```go
/// // Başlangıç eşleştirme
/// contains("/auth/sign-in", "/auth") // true
///
/// // Bitiş eşleştirme
/// contains("/api/users", "/users") // true
///
/// // İçerik eşleştirme
/// contains("/api/resource/posts", "/resource/") // true
///
/// // Bulunamaz
/// contains("/api/users", "/admin") // false
/// ```
///
/// ## Algoritma Mantığı
///
/// ```
/// IF s uzunluğu >= substr uzunluğu AND s başlangıcı == substr THEN
///     return true
/// ELSE IF s uzunluğu > substr uzunluğu AND s sonu == substr THEN
///     return true
/// ELSE IF s uzunluğu > substr uzunluğu AND findSubstring(s, substr) THEN
///     return true
/// ELSE
///     return false
/// ```
///
/// ## Performans Notları
///
/// - Başlangıç ve bitiş kontrolü O(n) - n: substring uzunluğu
/// - İçerik kontrolü O(n*m) - n: string uzunluğu, m: substring uzunluğu
/// - Kısa string'ler için yeterince hızlı
/// - Çok uzun string'ler için strings.Contains() kullanılabilir
///
/// ## Önemli Notlar
///
/// - Case-sensitive çalışır (büyük/küçük harf duyarlı)
/// - Boş substring her zaman true döner (başlangıç eşleştirme)
/// - Substring, string'den uzunsa false döner
/// - Standart kütüphane strings.Contains() alternatif olabilir
func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr ||
	       len(s) > len(substr) && s[len(s)-len(substr):] == substr ||
	       len(s) > len(substr) && findSubstring(s, substr)
}

/// # findSubstring
///
/// Bu fonksiyon, bir string içinde substring'i brute-force yöntemiyle arar.
/// Her pozisyonda substring eşleşmesi kontrol eder.
///
/// ## Kullanım Senaryoları
///
/// - contains fonksiyonu tarafından içerik eşleştirme için kullanılır
/// - String içinde substring arama
/// - Pattern matching
///
/// ## Parametreler
///
/// - `s`: Aranacak ana string
/// - `substr`: Aranacak substring
///
/// ## Dönüş Değeri
///
/// - `bool`: Substring bulunursa true, bulunamazsa false
///
/// ## Algoritma
///
/// 1. String'in her pozisyonunu iterate eder
/// 2. Her pozisyonda substring uzunluğunda bir slice alır
/// 3. Slice, substring ile eşleşirse true döner
/// 4. Hiçbir pozisyonda eşleşme yoksa false döner
///
/// ## Örnek Kullanım
///
/// ```go
/// // Substring ortada
/// findSubstring("/api/resource/users", "/resource/") // true
///
/// // Substring başta
/// findSubstring("/auth/sign-in", "/auth") // true
///
/// // Substring sonda
/// findSubstring("/api/users", "/users") // true
///
/// // Bulunamaz
/// findSubstring("/api/users", "/admin") // false
/// ```
///
/// ## Algoritma Detayı
///
/// ```
/// FOR i = 0 TO len(s) - len(substr) DO
///     IF s[i:i+len(substr)] == substr THEN
///         return true
///     END IF
/// END FOR
/// return false
/// ```
///
/// ## Performans Analizi
///
/// - **Zaman Karmaşıklığı**: O(n*m)
///   - n: ana string uzunluğu
///   - m: substring uzunluğu
/// - **Alan Karmaşıklığı**: O(1) - ek bellek kullanmaz
/// - **En İyi Durum**: O(m) - substring başta bulunur
/// - **En Kötü Durum**: O(n*m) - substring bulunamaz veya sonda bulunur
///
/// ## Alternatifler
///
/// Daha performanslı alternatifler:
/// - `strings.Contains(s, substr)`: Standart kütüphane, optimize edilmiş
/// - `strings.Index(s, substr) != -1`: Index bulma
/// - Boyer-Moore algoritması: Çok uzun string'ler için
/// - KMP algoritması: Pattern matching için
///
/// ## Önemli Notlar
///
/// - Brute-force algoritma, basit ama yavaş
/// - Kısa string'ler için yeterli
/// - Production'da strings.Contains() tercih edilebilir
/// - Case-sensitive çalışır
/// - Boş substring için davranış tanımsız (contains tarafından kontrol edilir)
func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
