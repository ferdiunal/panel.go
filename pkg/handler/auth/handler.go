// Package auth, kimlik doğrulama (authentication) işlemlerini yöneten HTTP handler'larını içerir.
// Bu paket, kullanıcı kaydı, giriş, çıkış, oturum yönetimi ve şifre sıfırlama gibi
// temel kimlik doğrulama işlemlerini gerçekleştirir.
//
// # Özellikler
//
// - Email/şifre tabanlı kayıt ve giriş
// - Güvenli oturum yönetimi (session management)
// - Hesap kilitleme koruması (account lockout protection)
// - Şifre sıfırlama işlemleri
// - CSRF koruması ile güvenli cookie yönetimi
// - IP ve User-Agent takibi
//
// # Güvenlik Özellikleri
//
// - HTTPOnly ve Secure cookie bayrakları
// - SameSite=Strict CSRF koruması
// - __Host- prefix ile gelişmiş cookie güvenliği
// - Başarısız giriş denemesi takibi
// - Hesap kilitleme mekanizması
// - Rate limiting desteği
//
// # Kullanım Senaryoları
//
// 1. **Kullanıcı Kaydı**: Yeni kullanıcıların sisteme kaydolması
// 2. **Kullanıcı Girişi**: Mevcut kullanıcıların kimlik doğrulaması
// 3. **Oturum Yönetimi**: Aktif oturumların doğrulanması ve yönetimi
// 4. **Güvenli Çıkış**: Kullanıcı oturumlarının güvenli sonlandırılması
// 5. **Şifre Kurtarma**: Unutulan şifrelerin güvenli sıfırlanması
package auth

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/domain/session"
	"github.com/ferdiunal/panel.go/pkg/domain/user"
	"github.com/ferdiunal/panel.go/pkg/middleware"
	"github.com/ferdiunal/panel.go/pkg/service/auth"
	"github.com/gofiber/fiber/v2"
)

// Handler, kimlik doğrulama işlemlerini yöneten HTTP handler yapısıdır.
//
// Bu yapı, kullanıcı kaydı, giriş, çıkış ve oturum yönetimi gibi tüm kimlik doğrulama
// işlemlerini koordine eder. Auth service ile iletişim kurarak iş mantığını yürütür
// ve HTTP isteklerini/yanıtlarını yönetir.
//
// # Alanlar
//
// - `service`: Kimlik doğrulama iş mantığını içeren servis katmanı
// - `accountLockout`: Başarısız giriş denemelerini takip eden ve hesap kilitleme uygulayan middleware
// - `environment`: Çalışma ortamı (production, test, vb.) - cookie güvenlik ayarlarını etkiler
//
// # Güvenlik Özellikleri
//
// - **Account Lockout**: Brute-force saldırılarına karşı koruma
// - **Environment-based Security**: Test ortamında HTTP, production'da HTTPS zorunluluğu
// - **Secure Cookie Management**: HTTPOnly, Secure, SameSite bayrakları
// - **IP Tracking**: Giriş yapan IP adreslerinin kaydedilmesi
// - **User-Agent Tracking**: Cihaz ve tarayıcı bilgilerinin takibi
//
// # Kullanım Örneği
//
// ```go
// authService := auth.NewService(db, emailService)
// lockout := middleware.NewAccountLockout(5, 15*time.Minute)
// handler := auth.NewHandler(authService, lockout, "production")
//
// // Route tanımlamaları
// app.Post("/auth/register", handler.RegisterEmail)
// app.Post("/auth/login", handler.LoginEmail)
// app.Post("/auth/logout", handler.SignOut)
// app.Get("/auth/session", handler.GetSession)
// ```
//
// # Önemli Notlar
//
// - Handler, HTTP katmanında çalışır ve iş mantığını service katmanına delege eder
// - Tüm hata durumları uygun HTTP status kodları ile döndürülür
// - Cookie isimlendirmesi environment'a göre değişir (__Host- prefix production'da)
// - Account lockout middleware opsiyoneldir (nil olabilir)
type Handler struct {
	service        *auth.Service
	accountLockout *middleware.AccountLockout
	environment    string
}

// NewHandler, yeni bir kimlik doğrulama handler'ı oluşturur ve yapılandırır.
//
// Bu fonksiyon, Handler yapısının constructor'ıdır ve tüm bağımlılıkları enjekte eder.
// Dependency injection pattern'i kullanarak gevşek bağlı (loosely coupled) bir mimari sağlar.
//
// # Parametreler
//
// - `service`: Kimlik doğrulama iş mantığını içeren servis katmanı (zorunlu)
// - `accountLockout`: Hesap kilitleme middleware'i (opsiyonel, nil olabilir)
// - `environment`: Çalışma ortamı ("production", "test", "development")
//
// # Döndürür
//
// - Yapılandırılmış Handler pointer'ı
//
// # Kullanım Örneği
//
// ```go
// // Tam özellikli production kurulumu
// authService := auth.NewService(db, emailService)
// lockout := middleware.NewAccountLockout(5, 15*time.Minute)
// handler := auth.NewHandler(authService, lockout, "production")
//
// // Test ortamı için basit kurulum (lockout olmadan)
// handler := auth.NewHandler(authService, nil, "test")
// ```
//
// # Önemli Notlar
//
// - Environment parametresi cookie güvenlik ayarlarını etkiler
// - "production" ortamında __Host- prefix ve Secure flag kullanılır
// - "test" ortamında HTTP üzerinden çalışabilmesi için güvenlik gevşetilir
// - accountLockout nil ise, hesap kilitleme koruması devre dışı kalır
//
// # Güvenlik Önerileri
//
// - Production ortamında mutlaka accountLockout middleware'i kullanın
// - Environment değerini doğru ayarladığınızdan emin olun
// - Test ortamını production'da asla kullanmayın
func NewHandler(service *auth.Service, accountLockout *middleware.AccountLockout, environment string) *Handler {
	return &Handler{
		service:        service,
		accountLockout: accountLockout,
		environment:    environment,
	}
}

// RegisterRequest, kullanıcı kayıt işlemi için gerekli bilgileri içeren istek yapısıdır.
//
// Bu yapı, yeni bir kullanıcının sisteme kaydolması için gerekli minimum bilgileri tanımlar.
// HTTP request body'sinden JSON formatında parse edilir.
//
// # Alanlar
//
// - `Name`: Kullanıcının tam adı (zorunlu)
// - `Email`: Kullanıcının email adresi (zorunlu, benzersiz olmalı)
// - `Password`: Kullanıcının şifresi (zorunlu, hash'lenecek)
//
// # JSON Örneği
//
// ```json
//
//	{
//	  "name": "Ahmet Yılmaz",
//	  "email": "ahmet@example.com",
//	  "password": "GüçlüŞifre123!"
//	}
//
// ```
//
// # Validasyon Kuralları
//
// - **Name**: Boş olmamalı, minimum 2 karakter
// - **Email**: Geçerli email formatında olmalı, sistemde benzersiz olmalı
// - **Password**: Minimum 8 karakter, güçlü şifre politikasına uymalı
//
// # Güvenlik Notları
//
// - Şifre asla düz metin olarak saklanmaz, bcrypt ile hash'lenir
// - Email adresi küçük harfe çevrilerek normalize edilir
// - Kayıt sonrası email doğrulama gerekebilir (implementasyona bağlı)
//
// # Kullanım Örneği
//
// ```go
// // Client tarafından gönderilen request
// POST /auth/register
// Content-Type: application/json
//
//	{
//	  "name": "Mehmet Demir",
//	  "email": "mehmet@example.com",
//	  "password": "Secure123!"
//	}
//
// ```
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterEmail, yeni bir kullanıcıyı email ve şifre ile sisteme kaydeder.
//
// Bu fonksiyon, kullanıcı kayıt işlemini yönetir ve başarılı kayıt sonrası kullanıcı bilgilerini döndürür.
// HTTP POST isteği ile çağrılır ve JSON formatında request body bekler.
//
// # HTTP Endpoint
//
// ```
// POST /auth/register
// Content-Type: application/json
// ```
//
// # Request Body
//
// ```json
//
//	{
//	  "name": "Kullanıcı Adı",
//	  "email": "kullanici@example.com",
//	  "password": "GüçlüŞifre123!"
//	}
//
// ```
//
// # Response
//
// **Başarılı (201 Created):**
// ```json
//
//	{
//	  "user": {
//	    "id": "uuid",
//	    "name": "Kullanıcı Adı",
//	    "email": "kullanici@example.com",
//	    "created_at": "2024-01-01T00:00:00Z"
//	  }
//	}
//
// ```
//
// **Hata Durumları:**
// - `400 Bad Request`: Geçersiz request body formatı
// - `409 Conflict`: Email adresi zaten kullanımda
// - `500 Internal Server Error`: Sunucu hatası
//
// # İş Akışı
//
// 1. Request body'den RegisterRequest parse edilir
// 2. Email benzersizliği kontrol edilir
// 3. Şifre bcrypt ile hash'lenir
// 4. Kullanıcı veritabanına kaydedilir
// 5. Kullanıcı bilgileri döndürülür (şifre hariç)
//
// # Güvenlik Özellikleri
//
// - Şifre asla düz metin olarak saklanmaz
// - Email adresi normalize edilir (küçük harf)
// - Benzersizlik kontrolü yapılır
// - Şifre hash'leme için bcrypt kullanılır (cost: 10)
//
// # Kullanım Örneği
//
// ```go
// // Route tanımı
// app.Post("/auth/register", handler.RegisterEmail)
//
// // Client tarafından kullanım
//
//	fetch('/auth/register', {
//	  method: 'POST',
//	  headers: { 'Content-Type': 'application/json' },
//	  body: JSON.stringify({
//	    name: 'Ahmet Yılmaz',
//	    email: 'ahmet@example.com',
//	    password: 'Secure123!'
//	  })
//	})
//
// ```
//
// # Önemli Notlar
//
// - Kayıt sonrası otomatik giriş yapılmaz (manuel login gerekir)
// - Email doğrulama implementasyona bağlıdır (şu an yok)
// - Rate limiting uygulanması önerilir (spam koruması)
// - CAPTCHA entegrasyonu düşünülebilir
//
// # Gelecek Geliştirmeler
//
// - Email doğrulama sistemi
// - Sosyal medya ile kayıt (OAuth)
// - İki faktörlü kimlik doğrulama (2FA)
// - Şifre güçlülük kontrolü
func (h *Handler) RegisterEmail(c *context.Context) error {
	var req RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	user, err := h.service.RegisterEmail(c.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		if err == auth.ErrEmailAlreadyExists {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Auto login after register? For now just return user.
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"user": user,
	})
}

// LoginRequest, kullanıcı giriş işlemi için gerekli kimlik bilgilerini içeren istek yapısıdır.
//
// Bu yapı, mevcut bir kullanıcının sisteme giriş yapması için gerekli bilgileri tanımlar.
// HTTP request body'sinden JSON formatında parse edilir.
//
// # Alanlar
//
// - `Email`: Kullanıcının kayıtlı email adresi (zorunlu)
// - `Password`: Kullanıcının şifresi (zorunlu, hash ile karşılaştırılacak)
//
// # JSON Örneği
//
// ```json
//
//	{
//	  "email": "kullanici@example.com",
//	  "password": "GüçlüŞifre123!"
//	}
//
// ```
//
// # Güvenlik Özellikleri
//
// - Şifre bcrypt hash ile karşılaştırılır
// - Başarısız giriş denemeleri takip edilir
// - Account lockout mekanizması ile brute-force koruması
// - IP adresi ve User-Agent bilgileri kaydedilir
// - Rate limiting uygulanabilir
//
// # Kullanım Örneği
//
// ```go
// // Client tarafından gönderilen request
// POST /auth/login
// Content-Type: application/json
//
//	{
//	  "email": "ahmet@example.com",
//	  "password": "Secure123!"
//	}
//
// ```
//
// # Önemli Notlar
//
// - Email adresi büyük/küçük harf duyarsızdır (normalize edilir)
// - Başarısız denemeler sonrası hesap geçici olarak kilitlenebilir
// - Başarılı giriş sonrası güvenli session cookie oluşturulur
// - Session token HTTPOnly ve Secure bayrakları ile korunur
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginEmail, kullanıcının email ve şifre ile sisteme giriş yapmasını sağlar.
//
// Bu fonksiyon, kimlik doğrulama sürecini yönetir ve başarılı giriş sonrası güvenli bir
// oturum (session) oluşturur. Brute-force saldırılarına karşı hesap kilitleme mekanizması,
// IP ve User-Agent takibi gibi gelişmiş güvenlik özellikleri içerir.
//
// # HTTP Endpoint
//
// ```
// POST /auth/login
// Content-Type: application/json
// ```
//
// # Request Body
//
// ```json
//
//	{
//	  "email": "kullanici@example.com",
//	  "password": "GüçlüŞifre123!"
//	}
//
// ```
//
// # Response
//
// **Başarılı (200 OK):**
// ```json
//
//	{
//	  "session": {
//	    "token": "secure-session-token-here",
//	    "expires": "2024-01-08T00:00:00Z"
//	  },
//	  "user": {
//	    "id": "uuid",
//	    "name": "Kullanıcı Adı",
//	    "email": "kullanici@example.com"
//	  }
//	}
//
// ```
//
// **Hata Durumları:**
// - `400 Bad Request`: Geçersiz request body formatı
// - `401 Unauthorized`: Yanlış email veya şifre
// - `429 Too Many Requests`: Hesap geçici olarak kilitli (çok fazla başarısız deneme)
//
// # İş Akışı
//
// 1. **Request Parsing**: JSON body'den LoginRequest parse edilir
// 2. **Account Lockout Check**: Hesabın kilitli olup olmadığı kontrol edilir
// 3. **IP Detection**: Client IP adresi tespit edilir (X-Forwarded-For desteği ile)
// 4. **Authentication**: Email ve şifre doğrulanır (bcrypt hash karşılaştırması)
// 5. **Failed Attempt Tracking**: Başarısız girişler kaydedilir
// 6. **Session Creation**: Başarılı giriş sonrası yeni session oluşturulur
// 7. **Cookie Setting**: Güvenli HTTPOnly cookie set edilir
// 8. **Response**: Session ve kullanıcı bilgileri döndürülür
//
// # Güvenlik Özellikleri
//
// ## 1. Account Lockout Protection
// - Başarısız giriş denemeleri takip edilir
// - Belirli sayıda başarısız denemeden sonra hesap geçici olarak kilitlenir
// - Kalan deneme hakkı kullanıcıya bildirilir
// - Başarılı giriş sonrası sayaç sıfırlanır
//
// ## 2. Secure Cookie Management
// - **HTTPOnly**: JavaScript erişimini engeller (XSS koruması)
// - **Secure**: HTTPS zorunluluğu (production'da)
// - **SameSite=Strict**: CSRF saldırılarını engeller
// - **__Host- Prefix**: Ek güvenlik katmanı (subdomain koruması)
// - **Path=/**: Tüm uygulama için geçerli
//
// ## 3. IP ve User-Agent Tracking
// - Her giriş denemesinin IP adresi kaydedilir
// - X-Forwarded-For header'ı desteklenir (proxy/load balancer arkasında)
// - User-Agent bilgisi kaydedilir (cihaz/tarayıcı takibi)
// - Şüpheli aktivite tespiti için kullanılabilir
//
// ## 4. Environment-based Security
// - **Production**: __Host- prefix, Secure=true, HTTPS zorunlu
// - **Test**: Normal cookie, Secure=false, HTTP izinli
//
// # Kullanım Örneği
//
// ```go
// // Route tanımı
// app.Post("/auth/login", handler.LoginEmail)
//
// // Client tarafından kullanım
//
//	const response = await fetch('/auth/login', {
//	  method: 'POST',
//	  headers: { 'Content-Type': 'application/json' },
//	  body: JSON.stringify({
//	    email: 'ahmet@example.com',
//	    password: 'Secure123!'
//	  }),
//	  credentials: 'include' // Cookie için gerekli
//	});
//
//	if (response.ok) {
//	  const { session, user } = await response.json();
//	  console.log('Giriş başarılı:', user.name);
//	} else if (response.status === 429) {
//
//	  console.error('Hesap kilitli, lütfen daha sonra tekrar deneyin');
//	} else if (response.status === 401) {
//
//	  const data = await response.json();
//	  if (data.remaining_attempts) {
//	    console.warn(`Yanlış şifre. Kalan deneme: ${data.remaining_attempts}`);
//	  }
//	}
//
// ```
//
// # Cookie Detayları
//
// **Production Cookie:**
// ```
// Set-Cookie: __Host-session_token=<token>;
//
//	Path=/;
//	Expires=<date>;
//	HttpOnly;
//	Secure;
//	SameSite=Strict
//
// ```
//
// **Test Cookie:**
// ```
// Set-Cookie: session_token=<token>;
//
//	Path=/;
//	Expires=<date>;
//	HttpOnly;
//	SameSite=Strict
//
// ```
//
// # Önemli Notlar
//
// - Session token otomatik olarak cookie'de saklanır
// - Client'ın cookie'leri kabul etmesi gerekir
// - CORS ayarları credentials: 'include' için yapılandırılmalı
// - Rate limiting ek bir katman olarak uygulanabilir
// - Başarısız giriş logları güvenlik analizi için saklanmalı
//
// # Güvenlik Uyarıları
//
// - ⚠️ Test environment'ı production'da asla kullanmayın
// - ⚠️ HTTPS olmadan production'a deploy etmeyin
// - ⚠️ Account lockout süresini çok kısa tutmayın (DDoS riski)
// - ⚠️ IP-based blocking dikkatli kullanın (shared IP'ler)
// - ⚠️ Session token'ları güvenli şekilde saklayın
//
// # Performans Notları
//
// - Bcrypt hash karşılaştırması CPU-intensive'dir (kasıtlı)
// - Account lockout kontrolü memory-based'dir (hızlı)
// - IP detection minimal overhead ekler
// - Session creation veritabanı yazma işlemi gerektirir
//
// # Gelecek Geliştirmeler
//
// - İki faktörlü kimlik doğrulama (2FA/TOTP)
// - Biometric authentication desteği
// - Sosyal medya ile giriş (OAuth)
// - Passwordless authentication (magic link)
// - Device fingerprinting
// - Anomaly detection (ML-based)
func (h *Handler) LoginEmail(c *context.Context) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// SECURITY: Check if account is locked
	if h.accountLockout != nil && h.accountLockout.IsLocked(req.Email) {
		return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
			"error": "Account temporarily locked due to too many failed login attempts. Please try again later.",
		})
	}

	// Get IP with fallback to X-Forwarded-For
	ip := c.IP()
	if forwarded := c.Get("X-Forwarded-For"); forwarded != "" {
		ip = forwarded
	}

	session, err := h.service.LoginEmail(c.Context(), req.Email, req.Password, ip, c.Get("User-Agent"))
	if err != nil {
		// SECURITY: Record failed login attempt
		if h.accountLockout != nil {
			h.accountLockout.RecordFailedAttempt(req.Email)
			remaining := h.accountLockout.GetRemainingAttempts(req.Email)
			if remaining > 0 {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error":              err.Error(),
					"remaining_attempts": remaining,
				})
			}
		}
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	// SECURITY: Reset failed attempts on successful login
	if h.accountLockout != nil {
		h.accountLockout.ResetAttempts(req.Email)
	}

	// SECURITY: Set secure session cookie
	// Use __Host- prefix for additional security (requires Secure=true, Path=/, no Domain)
	// In test environment, use regular cookie name to allow HTTP
	cookieName := "__Host-session_token"
	secure := true
	if h.environment == "test" {
		cookieName = "session_token"
		secure = false
	}

	c.Cookie(&fiber.Cookie{
		Name:     cookieName,
		Value:    session.Token,
		Expires:  session.ExpiresAt,
		HTTPOnly: true,
		Secure:   secure,
		SameSite: "Strict", // Strict for admin panels to prevent CSRF
		Path:     "/",
	})

	return c.JSON(fiber.Map{
		"session": fiber.Map{
			"token":   session.Token,
			"expires": session.ExpiresAt,
		},
		"user": session.User,
	})
}

// SignOut, kullanıcının aktif oturumunu sonlandırır ve güvenli çıkış yapar.
//
// Bu fonksiyon, kullanıcının mevcut oturumunu (session) geçersiz kılar ve session cookie'sini temizler.
// Hem sunucu tarafında session'ı veritabanından siler, hem de client tarafında cookie'yi temizler.
//
// # HTTP Endpoint
//
// ```
// POST /auth/logout
// Cookie: __Host-session_token=<token>
// ```
//
// # Response
//
// **Başarılı (200 OK):**
// ```json
//
//	{
//	  "message": "Signed out"
//	}
//
// ```
//
// # İş Akışı
//
// 1. **Cookie Reading**: Session token cookie'den okunur
// 2. **Session Invalidation**: Token varsa, veritabanından session silinir
// 3. **Cookie Clearing**: Client tarafında cookie temizlenir
// 4. **Success Response**: Başarı mesajı döndürülür
//
// # Güvenlik Özellikleri
//
// - Session token sunucu tarafında geçersiz kılınır (veritabanından silinir)
// - Cookie client tarafında temizlenir (Set-Cookie: Max-Age=0)
// - Token yoksa bile hata vermez (idempotent operation)
// - Logout işlemi her zaman başarılı döner (bilgi sızıntısı önlenir)
//
// # Kullanım Örneği
//
// ```go
// // Route tanımı
// app.Post("/auth/logout", handler.SignOut)
//
// // Client tarafından kullanım
//
//	const response = await fetch('/auth/logout', {
//	  method: 'POST',
//	  credentials: 'include' // Cookie gönderimi için gerekli
//	});
//
//	if (response.ok) {
//	  const data = await response.json();
//	  console.log(data.message); // "Signed out"
//	  // Kullanıcıyı login sayfasına yönlendir
//	  window.location.href = '/login';
//	}
//
// ```
//
// # Cookie Temizleme Detayları
//
// **Production:**
// ```
// Set-Cookie: __Host-session_token=;
//
//	Path=/;
//	Max-Age=0;
//	HttpOnly;
//	Secure;
//	SameSite=Strict
//
// ```
//
// **Test:**
// ```
// Set-Cookie: session_token=;
//
//	Path=/;
//	Max-Age=0;
//	HttpOnly;
//	SameSite=Strict
//
// ```
//
// # Önemli Notlar
//
// - Logout işlemi idempotent'tir (birden fazla kez çağrılabilir)
// - Token yoksa veya geçersizse bile başarılı döner
// - Session veritabanından fiziksel olarak silinir
// - Client tarafında cookie otomatik temizlenir
// - Logout sonrası korumalı endpoint'lere erişim engellenir
//
// # Best Practices
//
// - Logout sonrası kullanıcıyı login sayfasına yönlendirin
// - Client-side state'i temizleyin (localStorage, sessionStorage)
// - Logout işlemini POST method ile yapın (CSRF koruması)
// - Logout butonunu her sayfada erişilebilir yapın
//
// # Güvenlik Notları
//
// - Logout işlemi authentication gerektirmez (herkes çağırabilir)
// - Token yoksa bile hata vermez (timing attack önlenir)
// - Session veritabanından tamamen silinir (geri alınamaz)
// - Cookie temizleme işlemi tarayıcı tarafından garanti edilir
func (h *Handler) SignOut(c *context.Context) error {
	cookieName := "__Host-session_token"
	if h.environment == "test" {
		cookieName = "session_token"
	}

	token := c.Cookies(cookieName)
	if token != "" {
		h.service.Logout(c.Context(), token)
	}

	c.ClearCookie(cookieName)
	return c.JSON(fiber.Map{"message": "Signed out"})
}

// GetSession, kullanıcının mevcut oturum bilgilerini döndürür.
//
// Bu fonksiyon, client tarafındaki session cookie'sini doğrular ve geçerliyse
// oturum bilgileri ile kullanıcı bilgilerini döndürür. Geçersiz veya süresi dolmuş
// session'lar için null döner ve cookie temizlenir.
//
// # HTTP Endpoint
//
// ```
// GET /auth/session
// Cookie: __Host-session_token=<token>
// ```
//
// # Response
//
// **Geçerli Session (200 OK):**
// ```json
//
//	{
//	  "session": {
//	    "token": "secure-session-token-here",
//	    "expires": "2024-01-08T00:00:00Z"
//	  },
//	  "user": {
//	    "id": "uuid",
//	    "name": "Kullanıcı Adı",
//	    "email": "kullanici@example.com"
//	  }
//	}
//
// ```
//
// **Geçersiz/Yok Session (200 OK):**
// ```json
//
//	{
//	  "session": null
//	}
//
// ```
//
// # İş Akışı
//
// 1. **Cookie Reading**: Session token cookie'den okunur
// 2. **Token Validation**: Token yoksa null döndürülür
// 3. **Session Validation**: Token varsa veritabanından doğrulanır
// 4. **Expiry Check**: Session süresi kontrol edilir
// 5. **User Preloading**: Kullanıcı bilgileri eager loading ile yüklenir
// 6. **Cookie Cleanup**: Geçersiz session için cookie temizlenir
// 7. **Response**: Session ve kullanıcı bilgileri döndürülür
//
// # Kullanım Senaryoları
//
// 1. **Sayfa Yüklemesi**: Uygulama başlangıcında kullanıcı durumunu kontrol etme
// 2. **Session Refresh**: Periyodik olarak session geçerliliğini kontrol etme
// 3. **Protected Routes**: Korumalı sayfalara erişim öncesi doğrulama
// 4. **User Info Display**: Kullanıcı bilgilerini gösterme (navbar, profil)
//
// # Kullanım Örneği
//
// ```go
// // Route tanımı
// app.Get("/auth/session", handler.GetSession)
//
// // Client tarafından kullanım (React örneği)
//
//	const checkSession = async () => {
//	  const response = await fetch('/auth/session', {
//	    credentials: 'include' // Cookie gönderimi için gerekli
//	  });
//
//	  const data = await response.json();
//
//	  if (data.session) {
//	    // Kullanıcı giriş yapmış
//	    console.log('Hoş geldin:', data.user.name);
//	    setUser(data.user);
//	    setIsAuthenticated(true);
//	  } else {
//	    // Kullanıcı giriş yapmamış
//	    console.log('Oturum bulunamadı');
//	    setIsAuthenticated(false);
//	  }
//	};
//
// // Uygulama başlangıcında
//
//	useEffect(() => {
//	  checkSession();
//	}, []);
//
// // Periyodik kontrol (opsiyonel)
//
//	useEffect(() => {
//	  const interval = setInterval(checkSession, 5 * 60 * 1000); // 5 dakikada bir
//	  return () => clearInterval(interval);
//	}, []);
//
// ```
//
// # Güvenlik Özellikleri
//
// - Session token cookie'den okunur (güvenli)
// - Token veritabanında doğrulanır (server-side validation)
// - Süresi dolmuş session'lar otomatik reddedilir
// - Geçersiz token için cookie temizlenir (cleanup)
// - User bilgileri preload edilir (N+1 query önlenir)
//
// # Önemli Notlar
//
// - Bu endpoint authentication gerektirmez (public)
// - Session yoksa veya geçersizse null döner (hata değil)
// - Cookie otomatik olarak gönderilir (credentials: 'include' ile)
// - Response her zaman 200 OK döner (bilgi sızıntısı önlenir)
// - Geçersiz session için cookie client tarafında temizlenir
//
// # Best Practices
//
// - Uygulama başlangıcında session kontrolü yapın
// - Session bilgilerini global state'te saklayın (Redux, Context)
// - Periyodik session refresh düşünün (long-running apps)
// - Session null ise kullanıcıyı login sayfasına yönlendirin
// - Loading state kullanın (session kontrolü sırasında)
//
// # Performans Notları
//
// - User bilgileri eager loading ile yüklenir (tek query)
// - Session validation veritabanı sorgusu gerektirir
// - Cookie okuma işlemi minimal overhead
// - Response caching düşünülebilir (kısa süreli)
//
// # Hata Durumları
//
// - Token yoksa: `{"session": null}` döner
// - Token geçersizse: Cookie temizlenir, `{"session": null}` döner
// - Session süresi dolmuşsa: Cookie temizlenir, `{"session": null}` döner
// - Veritabanı hatası: Cookie temizlenir, `{"session": null}` döner
//
// # Gelecek Geliştirmeler
//
// - Session refresh mekanizması (sliding expiration)
// - Multiple device tracking (session listesi)
// - Session revocation (uzaktan çıkış)
// - Activity tracking (son aktivite zamanı)
func (h *Handler) GetSession(c *context.Context) error {
	cookieName := "__Host-session_token"
	if h.environment == "test" {
		cookieName = "session_token"
	}

	token := c.Cookies(cookieName)
	if token == "" {
		return c.JSON(fiber.Map{"session": nil})
	}

	session, err := h.service.ValidateSession(c.Context(), token)
	if err != nil {
		c.ClearCookie(cookieName)
		return c.JSON(fiber.Map{"session": nil})
	}

	return c.JSON(fiber.Map{
		"session": fiber.Map{
			"token":   session.Token,
			"expires": session.ExpiresAt,
		},
		"user": session.User, // Preloaded? Service FindByToken preloads User.
	})
}

// SessionMiddleware, korumalı endpoint'ler için kimlik doğrulama middleware'idir.
//
// Bu middleware, gelen isteklerin geçerli bir session token'ı içerip içermediğini kontrol eder.
// Geçerli session varsa, kullanıcı ve session bilgilerini context'e ekler ve isteğin devam etmesine
// izin verir. Geçersiz veya eksik session durumunda 401 Unauthorized hatası döndürür.
//
// # Kullanım Senaryoları
//
// 1. **Protected Routes**: Sadece giriş yapmış kullanıcıların erişebileceği endpoint'ler
// 2. **User Context**: Handler'larda kullanıcı bilgisine erişim gerektiğinde
// 3. **Authorization**: Rol tabanlı yetkilendirme için ön koşul
// 4. **Audit Logging**: Kullanıcı aktivitelerini kaydetme
//
// # Middleware Kullanımı
//
// ```go
// // Tek bir route için
// app.Get("/api/profile", handler.SessionMiddleware, profileHandler)
//
// // Route grubu için
// api := app.Group("/api")
// api.Use(handler.SessionMiddleware)
// api.Get("/profile", profileHandler)
// api.Get("/settings", settingsHandler)
// api.Post("/posts", createPostHandler)
//
// // Nested groups
// admin := api.Group("/admin")
// admin.Use(handler.SessionMiddleware)
// admin.Use(adminRoleMiddleware) // Ek middleware'ler
// admin.Get("/users", listUsersHandler)
// ```
//
// # Context'ten Veri Okuma
//
// ```go
// // Handler içinde kullanıcı bilgisine erişim
//
//	func profileHandler(c *context.Context) error {
//	  // Session bilgisi
//	  session := c.Locals("session").(*auth.Session)
//	  fmt.Println("Session Token:", session.Token)
//	  fmt.Println("Expires At:", session.ExpiresAt)
//
//	  // Kullanıcı bilgisi
//	  user := c.Locals("user").(*models.User)
//	  fmt.Println("User ID:", user.ID)
//	  fmt.Println("User Email:", user.Email)
//	  fmt.Println("User Name:", user.Name)
//
//	  return c.JSON(fiber.Map{
//	    "user": user,
//	  })
//	}
//
// ```
//
// # İş Akışı
//
// 1. **Cookie Reading**: Session token cookie'den okunur
// 2. **Token Validation**: Token yoksa 401 Unauthorized döner
// 3. **Session Validation**: Token veritabanında doğrulanır
// 4. **Expiry Check**: Session süresi kontrol edilir
// 5. **Context Population**: Session ve user bilgileri context'e eklenir
// 6. **Next Handler**: Sonraki middleware/handler çağrılır
//
// # Response
//
// **Başarılı**: Sonraki handler'a geçer, context'te şunlar bulunur:
// - `c.Locals("session")`: Session objesi
// - `c.Locals("user")`: User objesi
//
// **Başarısız (401 Unauthorized):**
// ```json
//
//	{
//	  "error": "Unauthorized"
//	}
//
// ```
//
// # Güvenlik Özellikleri
//
// - Session token cookie'den güvenli şekilde okunur
// - Token veritabanında doğrulanır (server-side validation)
// - Süresi dolmuş session'lar otomatik reddedilir
// - Geçersiz token için cookie temizlenir
// - User bilgileri preload edilir (N+1 query önlenir)
// - 401 status code ile unauthorized erişim engellenir
//
// # Önemli Notlar
//
// - Bu middleware korumalı route'larda kullanılmalıdır
// - Public endpoint'lerde kullanılmamalıdır (login, register, vb.)
// - Context'e eklenen veriler type assertion ile okunmalıdır
// - Session ve user bilgileri her istekte fresh olarak yüklenir
// - Cookie otomatik olarak gönderilir (credentials: 'include' ile)
//
// # Client Tarafı Kullanımı
//
// ```javascript
// // Korumalı endpoint'e istek
//
//	const response = await fetch('/api/profile', {
//	  method: 'GET',
//	  credentials: 'include' // Cookie gönderimi için ZORUNLU
//	});
//
//	if (response.status === 401) {
//	  // Session geçersiz veya yok
//	  console.log('Oturum süresi dolmuş, lütfen tekrar giriş yapın');
//	  window.location.href = '/login';
//	} else if (response.ok) {
//
//	  const data = await response.json();
//	  console.log('Kullanıcı:', data.user);
//	}
//
// ```
//
// # CORS Yapılandırması
//
// ```go
// // Fiber CORS middleware yapılandırması
//
//	app.Use(cors.New(cors.Config{
//	  AllowOrigins:     "https://yourdomain.com",
//	  AllowCredentials: true, // Cookie için ZORUNLU
//	  AllowHeaders:     "Origin, Content-Type, Accept",
//	}))
//
// ```
//
// # Best Practices
//
// 1. **Middleware Sırası**: SessionMiddleware'i diğer auth middleware'lerden önce kullanın
// 2. **Error Handling**: 401 hatalarını client'ta yakalayın ve login'e yönlendirin
// 3. **Token Refresh**: Long-running apps için session refresh mekanizması düşünün
// 4. **Logging**: Unauthorized erişim denemelerini loglayın
// 5. **Rate Limiting**: Brute-force saldırılarına karşı rate limiting ekleyin
//
// # Performans Notları
//
// - Her istekte veritabanı sorgusu yapılır (session validation)
// - User bilgileri eager loading ile yüklenir (tek query)
// - Cookie okuma işlemi minimal overhead
// - Caching stratejisi düşünülebilir (dikkatli kullanın)
//
// # Hata Durumları
//
// - **Token yok**: 401 Unauthorized, cookie temizlenmez
// - **Token geçersiz**: 401 Unauthorized, cookie temizlenir
// - **Session süresi dolmuş**: 401 Unauthorized, cookie temizlenir
// - **Veritabanı hatası**: 401 Unauthorized, cookie temizlenir
//
// # Güvenlik Uyarıları
//
// - ⚠️ Public endpoint'lerde bu middleware'i kullanmayın
// - ⚠️ CORS AllowCredentials: true ayarını dikkatli yapılandırın
// - ⚠️ Context'ten okunan verileri type assertion ile kontrol edin
// - ⚠️ Session validation her istekte yapılır (performans vs güvenlik)
//
// # Gelecek Geliştirmeler
//
// - Session caching (Redis ile)
// - Role-based access control (RBAC)
// - Permission-based authorization
// - Multi-tenancy support
// - API key authentication (alternatif)
// - JWT token support (stateless alternative)
func (h *Handler) SessionMiddleware(c *context.Context) error {
	if apiKeyAuth, ok := c.Locals(middleware.APIKeyAuthenticatedLocalKey).(bool); ok && apiKeyAuth {
		apiUser := &user.User{
			Name:  "API Key",
			Email: "api-key@panel.local",
			Role:  "admin",
		}
		c.Locals("user", apiUser)
		c.Locals("session", &session.Session{
			Token: "api-key",
			User:  apiUser,
		})
		return c.Next()
	}

	cookieName := "__Host-session_token"
	if h.environment == "test" {
		cookieName = "session_token"
	}

	token := c.Cookies(cookieName)
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	session, err := h.service.ValidateSession(c.Context(), token)
	if err != nil {
		c.ClearCookie(cookieName)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	c.Locals("session", session)
	c.Locals("user", session.User)

	return c.Next()
}

// ForgotPasswordRequest, şifre sıfırlama isteği için gerekli bilgileri içeren yapıdır.
//
// Bu yapı, kullanıcının unuttuğu şifresini sıfırlamak için email adresini alır.
// HTTP request body'sinden JSON formatında parse edilir.
//
// # Alanlar
//
// - `Email`: Şifre sıfırlama linki gönderilecek email adresi (zorunlu)
//
// # JSON Örneği
//
// ```json
//
//	{
//	  "email": "kullanici@example.com"
//	}
//
// ```
//
// # Güvenlik Özellikleri
//
// - Email varlığı kontrol edilir ancak sonuç kullanıcıya bildirilmez (bilgi sızıntısı önlenir)
// - Şifre sıfırlama token'ı güvenli şekilde oluşturulur (cryptographically secure)
// - Token sınırlı süre için geçerlidir (genellikle 1 saat)
// - Token tek kullanımlıktır (kullanıldıktan sonra geçersiz olur)
//
// # Kullanım Örneği
//
// ```go
// // Client tarafından gönderilen request
// POST /auth/forgot-password
// Content-Type: application/json
//
//	{
//	  "email": "ahmet@example.com"
//	}
//
// ```
//
// # Önemli Notlar
//
// - Email adresi sistemde kayıtlı olmasa bile başarılı response döner
// - Bu davranış, email enumeration saldırılarını önler
// - Şifre sıfırlama linki email ile gönderilir
// - Link belirli bir süre sonra otomatik olarak geçersiz olur
type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

// ForgotPassword, kullanıcının unuttuğu şifresini sıfırlaması için email gönderir.
//
// Bu fonksiyon, şifre sıfırlama sürecini başlatır. Kullanıcının email adresine güvenli bir
// şifre sıfırlama linki gönderir. Güvenlik nedeniyle, email adresinin sistemde kayıtlı olup
// olmadığı hakkında bilgi vermez (email enumeration saldırılarını önler).
//
// # HTTP Endpoint
//
// ```
// POST /auth/forgot-password
// Content-Type: application/json
// ```
//
// # Request Body
//
// ```json
//
//	{
//	  "email": "kullanici@example.com"
//	}
//
// ```
//
// # Response
//
// **Her Zaman Başarılı (200 OK):**
// ```json
//
//	{
//	  "message": "If an account exists with this email, a password reset link has been sent."
//	}
//
// ```
//
// # İş Akışı
//
// 1. **Request Parsing**: JSON body'den ForgotPasswordRequest parse edilir
// 2. **Email Lookup**: Email adresi veritabanında aranır
// 3. **Token Generation**: Güvenli reset token oluşturulur (cryptographically secure)
// 4. **Token Storage**: Token veritabanında saklanır (hash'lenmiş, süreli)
// 5. **Email Sending**: Şifre sıfırlama linki email ile gönderilir
// 6. **Generic Response**: Email varlığından bağımsız genel mesaj döndürülür
//
// # Güvenlik Özellikleri
//
// ## 1. Email Enumeration Prevention
// - Email sistemde kayıtlı olmasa bile aynı başarı mesajı döner
// - Response timing attack'lere karşı korunur
// - Kullanıcı varlığı hakkında bilgi sızdırılmaz
//
// ## 2. Secure Token Generation
// - Cryptographically secure random token (32+ bytes)
// - Token hash'lenerek veritabanında saklanır
// - Token tek kullanımlıktır (kullanıldıktan sonra silinir)
// - Token sınırlı süre için geçerlidir (genellikle 1 saat)
//
// ## 3. Rate Limiting
// - Aynı email için sık istek engellenir
// - IP-based rate limiting uygulanabilir
// - Spam ve abuse koruması
//
// ## 4. Email Security
// - Reset linki HTTPS üzerinden gönderilir
// - Link tek kullanımlıktır
// - Link süresi dolunca otomatik geçersiz olur
//
// # Kullanım Örneği
//
// ```go
// // Route tanımı
// app.Post("/auth/forgot-password", handler.ForgotPassword)
//
// // Client tarafından kullanım
//
//	const forgotPassword = async (email) => {
//	  const response = await fetch('/auth/forgot-password', {
//	    method: 'POST',
//	    headers: { 'Content-Type': 'application/json' },
//	    body: JSON.stringify({ email })
//	  });
//
//	  if (response.ok) {
//	    const data = await response.json();
//	    // Her zaman aynı mesaj
//	    alert(data.message);
//	    // Kullanıcıyı bilgilendir
//	    showNotification(
//	      'Şifre sıfırlama linki gönderildi',
//	      'Email adresinizi kontrol edin. Link 1 saat geçerlidir.'
//	    );
//	  }
//	};
//
// // Form submit handler
//
//	const handleSubmit = (e) => {
//	  e.preventDefault();
//	  const email = e.target.email.value;
//	  forgotPassword(email);
//	};
//
// ```
//
// # Email İçeriği Örneği
//
// ```
// Konu: Şifre Sıfırlama Talebi
//
// Merhaba,
//
// Hesabınız için şifre sıfırlama talebinde bulunuldu.
// Şifrenizi sıfırlamak için aşağıdaki linke tıklayın:
//
// https://yourdomain.com/reset-password?token=<secure-token>
//
// Bu link 1 saat geçerlidir ve tek kullanımlıktır.
//
// Eğer bu talebi siz yapmadıysanız, bu emaili görmezden gelebilirsiniz.
// Şifreniz değiştirilmeyecektir.
//
// Saygılarımızla,
// Panel.go Ekibi
// ```
//
// # Reset Password Flow
//
// 1. **Forgot Password**: Kullanıcı email adresini girer
// 2. **Email Sent**: Şifre sıfırlama linki gönderilir
// 3. **Click Link**: Kullanıcı email'deki linke tıklar
// 4. **Validate Token**: Token doğrulanır (GET /reset-password?token=xxx)
// 5. **New Password Form**: Yeni şifre formu gösterilir
// 6. **Submit New Password**: Yeni şifre gönderilir (POST /reset-password)
// 7. **Password Updated**: Şifre güncellenir, token silinir
// 8. **Auto Login**: Kullanıcı otomatik giriş yapabilir (opsiyonel)
//
// # Hata Durumları
//
// - `400 Bad Request`: Geçersiz request body formatı
// - `500 Internal Server Error`: Email gönderimi başarısız (nadiren döner)
//
// **Not**: Email bulunamadığında bile 200 OK döner (güvenlik)
//
// # Önemli Notlar
//
// - Response her zaman aynı mesajı döndürür (email enumeration önlenir)
// - Email gönderimi asenkron olabilir (performans için)
// - Token veritabanında hash'lenmiş olarak saklanır
// - Eski/kullanılmış token'lar periyodik olarak temizlenmelidir
// - Rate limiting mutlaka uygulanmalıdır (spam önleme)
//
// # Best Practices
//
// 1. **User Experience**:
//   - Kullanıcıya net talimatlar verin
//   - Email gelmezse spam klasörünü kontrol etmesini söyleyin
//   - Token süresi hakkında bilgilendirin
//   - Yeni istek gönderebileceğini belirtin
//
// 2. **Security**:
//   - Token'ları güvenli şekilde oluşturun (crypto/rand)
//   - Token'ları hash'leyerek saklayın
//   - Kısa süre için geçerli tutun (1 saat)
//   - Rate limiting uygulayın
//   - Email enumeration'a izin vermeyin
//
// 3. **Email Delivery**:
//   - Güvenilir email servisi kullanın (SendGrid, AWS SES)
//   - Email template'leri profesyonel tutun
//   - SPF, DKIM, DMARC ayarlarını yapın
//   - Bounce ve complaint'leri takip edin
//
// 4. **Monitoring**:
//   - Başarısız email gönderimlerini loglayın
//   - Şüpheli aktiviteleri tespit edin
//   - Token kullanım oranlarını izleyin
//   - Abuse pattern'lerini takip edin
//
// # Güvenlik Uyarıları
//
// - ⚠️ Email enumeration'a asla izin vermeyin (timing attack dahil)
// - ⚠️ Token'ları düz metin olarak saklamayın
// - ⚠️ Token süresini çok uzun tutmayın (max 1-2 saat)
// - ⚠️ Rate limiting olmadan production'a çıkmayın
// - ⚠️ HTTPS olmadan reset link'i göndermeyin
// - ⚠️ Kullanılmış token'ları hemen silin
//
// # Performans Notları
//
// - Email gönderimi asenkron yapılabilir (queue sistemi)
// - Token generation CPU-intensive değildir
// - Veritabanı yazma işlemi gerektirir
// - Email servisi timeout'u dikkate alın
//
// # Gelecek Geliştirmeler
//
// - SMS ile şifre sıfırlama (alternatif kanal)
// - Security questions (ek doğrulama)
// - Account recovery codes (backup method)
// - Multi-factor authentication (2FA)
// - Biometric recovery (mobile apps)
// - Self-service account recovery
//
// # İlgili Endpoint'ler
//
// - `POST /auth/reset-password`: Token ile şifre sıfırlama
// - `GET /auth/validate-reset-token`: Token geçerliliği kontrolü
// - `POST /auth/resend-reset-email`: Reset email'i tekrar gönderme
//
// # Örnek Kullanım Senaryoları
//
// **Senaryo 1: Başarılı Şifre Sıfırlama**
// 1. Kullanıcı forgot-password sayfasına gider
// 2. Email adresini girer ve gönderir
// 3. Başarı mesajı görür
// 4. Email'ini kontrol eder
// 5. Reset linkine tıklar
// 6. Yeni şifresini girer
// 7. Şifre güncellenir
// 8. Login sayfasına yönlendirilir
//
// **Senaryo 2: Yanlış Email**
// 1. Kullanıcı yanlış email adresi girer
// 2. Aynı başarı mesajını görür (güvenlik)
// 3. Email gelmez (çünkü hesap yok)
// 4. Kullanıcı doğru email'i hatırlar
// 5. Tekrar dener
//
// **Senaryo 3: Token Süresi Dolmuş**
// 1. Kullanıcı email'deki linke geç tıklar
// 2. "Token süresi dolmuş" hatası alır
// 3. Yeni reset talebi oluşturur
// 4. Yeni email gelir
// 5. Hızlıca işlemi tamamlar
func (h *Handler) ForgotPassword(c *context.Context) error {
	var req ForgotPasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := h.service.ForgotPassword(c.Context(), req.Email); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Always return success for security (don't reveal if email exists)
	return c.JSON(fiber.Map{
		"message": "If an account exists with this email, a password reset link has been sent.",
	})
}
