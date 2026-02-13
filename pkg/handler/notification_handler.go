package handler

import (
	"strconv"

	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/notification"
	"github.com/gofiber/fiber/v2"
)

// NotificationHandler, bildirim (notification) ile ilgili HTTP isteklerini yöneten handler yapısıdır.
//
// # Genel Bakış
//
// Bu yapı, kullanıcı bildirimlerinin yönetimi için gerekli tüm HTTP endpoint'lerini sağlar.
// Bildirim servisi ile etkileşime geçerek okunmamış bildirimleri getirme, bildirimleri
// okundu olarak işaretleme ve toplu işlemler gibi operasyonları gerçekleştirir.
//
// # Özellikler
//
// - Okunmamış bildirimleri listeleme
// - Tekil bildirim okundu işaretleme
// - Toplu bildirim okundu işaretleme
// - Kullanıcı bazlı yetkilendirme kontrolü
// - RESTful API standartlarına uygun yanıt formatları
//
// # Kullanım Senaryoları
//
// 1. **Bildirim Merkezi**: Kullanıcının okunmamış bildirimlerini görüntüleme
// 2. **Bildirim Yönetimi**: Bildirimleri tek tek veya toplu olarak okundu işaretleme
// 3. **Gerçek Zamanlı Güncellemeler**: WebSocket veya polling ile bildirim sayısı takibi
//
// # Güvenlik
//
// - Tüm endpoint'ler kullanıcı kimlik doğrulaması gerektirir
// - Her kullanıcı sadece kendi bildirimlerine erişebilir
// - Context'ten alınan user bilgisi ile yetkilendirme yapılır
//
// # Örnek Kullanım
//
//	```go
//	// Handler oluşturma
//	notificationService := notification.NewService(db)
//	handler := NewNotificationHandler(notificationService)
//
//	// Route tanımlama
//	app.Get("/api/notifications/unread", handler.HandleGetUnreadNotifications)
//	app.Put("/api/notifications/:id/read", handler.HandleMarkAsRead)
//	app.Put("/api/notifications/read-all", handler.HandleMarkAllAsRead)
//	```
//
// # Bağımlılıklar
//
// - `notification.Service`: Bildirim iş mantığını yöneten servis
// - `context.Context`: Fiber context wrapper'ı
// - `fiber.v2`: HTTP framework
//
// # Notlar
//
// - Handler'lar middleware chain'inden geçtikten sonra çalışır
// - User bilgisi context'e middleware tarafından enjekte edilmelidir
// - Tüm yanıtlar JSON formatındadır
type NotificationHandler struct {
	// Service, bildirim operasyonlarını gerçekleştiren servis instance'ı
	Service *notification.Service
}

// NewNotificationHandler, yeni bir NotificationHandler instance'ı oluşturur.
//
// # Genel Bakış
//
// Bu fonksiyon, bildirim servisi ile yapılandırılmış bir NotificationHandler oluşturur.
// Factory pattern kullanarak handler'ın doğru şekilde başlatılmasını sağlar.
//
// # Parametreler
//
// - `service`: Bildirim operasyonlarını yöneten notification.Service pointer'ı
//   - nil olmamalıdır
//   - Veritabanı bağlantısı ile yapılandırılmış olmalıdır
//
// # Döndürür
//
// - Yapılandırılmış NotificationHandler pointer'ı
//
// # Kullanım Senaryoları
//
// 1. **Uygulama Başlatma**: Ana uygulama başlatılırken handler'ları oluşturma
// 2. **Dependency Injection**: IoC container'dan servis alıp handler oluşturma
// 3. **Test Ortamı**: Mock servis ile test handler'ı oluşturma
//
// # Örnek Kullanım
//
//	```go
//	// Üretim ortamı
//	db := setupDatabase()
//	notificationService := notification.NewService(db)
//	handler := NewNotificationHandler(notificationService)
//
//	// Test ortamı
//	mockService := &notification.MockService{}
//	testHandler := NewNotificationHandler(mockService)
//	```
//
// # Notlar
//
// - Handler oluşturulduktan sonra route'lara bağlanmalıdır
// - Servis nil kontrolü yapılmaz, çağıran tarafın sorumluluğundadır
// - Handler thread-safe değildir, her request için ayrı context kullanılmalıdır
//
// # Bağımlılıklar
//
// - `notification.Service`: Bildirim iş mantığı servisi
func NewNotificationHandler(service *notification.Service) *NotificationHandler {
	return &NotificationHandler{
		Service: service,
	}
}

// HandleGetUnreadNotifications, mevcut kullanıcının okunmamış bildirimlerini döndürür.
//
// # Genel Bakış
//
// Bu fonksiyon, kimliği doğrulanmış kullanıcının okunmamış tüm bildirimlerini getirir.
// Context'ten kullanıcı bilgisini alır, yetkilendirme kontrolü yapar ve kullanıcıya
// ait okunmamış bildirimleri JSON formatında döndürür.
//
// # HTTP Endpoint
//
// - **Method**: GET
// - **Path**: `/api/notifications/unread` (örnek)
// - **Auth**: Gerekli (middleware ile sağlanmalı)
//
// # İstek
//
// - **Headers**:
//   - `Authorization`: Bearer token veya session cookie
// - **Body**: Yok
// - **Query Params**: Yok
//
// # Yanıt
//
// **Başarılı (200 OK)**:
//
//	```json
//	{
//	  "data": [
//	    {
//	      "id": 1,
//	      "user_id": 123,
//	      "title": "Yeni mesaj",
//	      "message": "Profilinize yeni bir yorum yapıldı",
//	      "type": "comment",
//	      "read_at": null,
//	      "created_at": "2026-02-07T10:30:00Z"
//	    }
//	  ]
//	}
//	```
//
// **Hata Yanıtları**:
//
// - **401 Unauthorized**: Kullanıcı kimliği doğrulanmamış
//
//	```json
//	{"error": "Unauthorized"}
//	```
//
// - **500 Internal Server Error**: Veritabanı hatası veya geçersiz kullanıcı
//
//	```json
//	{"error": "Invalid user"}
//	{"error": "database error message"}
//	```
//
// # Parametreler
//
// - `c`: Fiber context wrapper'ı, kullanıcı bilgisi ve HTTP işlemleri için kullanılır
//
// # Döndürür
//
// - `error`: HTTP yanıt hatası veya nil
//
// # Kullanım Senaryoları
//
// 1. **Bildirim Merkezi**: Kullanıcı bildirim merkezini açtığında okunmamış bildirimleri gösterme
// 2. **Badge Sayısı**: Navbar'da okunmamış bildirim sayısını gösterme
// 3. **Polling**: Periyodik olarak yeni bildirimleri kontrol etme
// 4. **WebSocket Alternatifi**: Gerçek zamanlı bildirim yerine polling ile güncelleme
//
// # İş Akışı
//
// 1. Context'ten user bilgisini al (`c.Locals("user")`)
// 2. User nil kontrolü yap, nil ise 401 döndür
// 3. User'dan ID'yi çıkar (GetID() interface'i kullanarak)
// 4. ID çıkarılamazsa 500 döndür
// 5. Service üzerinden okunmamış bildirimleri getir
// 6. Hata varsa 500 döndür
// 7. Bildirimleri JSON formatında döndür
//
// # Güvenlik
//
// - **Yetkilendirme**: Middleware tarafından user context'e enjekte edilmelidir
// - **Veri İzolasyonu**: Her kullanıcı sadece kendi bildirimlerini görebilir
// - **SQL Injection**: Service katmanında parametreli sorgular kullanılmalıdır
//
// # Performans Notları
//
// - Okunmamış bildirimlerde index kullanılmalıdır (user_id, read_at)
// - Çok sayıda bildirim varsa pagination eklenebilir
// - Cache mekanizması ile performans artırılabilir
//
// # Örnek Kullanım
//
//	```go
//	// Route tanımlama
//	app.Get("/api/notifications/unread",
//	    authMiddleware,
//	    handler.HandleGetUnreadNotifications,
//	)
//
//	// Frontend'den çağrı (JavaScript)
//	fetch('/api/notifications/unread', {
//	    headers: {
//	        'Authorization': 'Bearer ' + token
//	    }
//	})
//	.then(res => res.json())
//	.then(data => {
//	    console.log('Okunmamış bildirimler:', data.data);
//	});
//	```
//
// # Avantajlar
//
// - Basit ve anlaşılır API
// - RESTful standartlara uygun
// - Kullanıcı bazlı veri izolasyonu
// - Kolay test edilebilir
//
// # Dezavantajlar
//
// - Gerçek zamanlı değil (polling gerektirir)
// - Çok sayıda bildirimde performans sorunu olabilir
// - Pagination olmadan tüm bildirimleri döndürür
//
// # Önemli Notlar
//
// - User bilgisi middleware tarafından context'e eklenmelidir
// - GetID() interface'i user model'inde implement edilmelidir
// - Service katmanında hata yönetimi yapılmalıdır
// - Production'da rate limiting eklenmelidir
//
// # İyileştirme Önerileri
//
// - Pagination desteği eklenebilir (limit, offset)
// - Filtreleme seçenekleri eklenebilir (type, date range)
// - Cache mekanizması eklenebilir (Redis)
// - WebSocket ile gerçek zamanlı bildirim desteği
//
// # Bağımlılıklar
//
// - `context.Context`: Fiber context wrapper
// - `notification.Service`: Bildirim servisi
// - Auth Middleware: User bilgisini context'e ekleyen middleware
func (h *NotificationHandler) HandleGetUnreadNotifications(c *context.Context) error {
	// Get user from context
	user := c.Locals("user")
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Extract user ID
	var userID uint
	if u, ok := user.(interface{ GetID() uint }); ok {
		userID = u.GetID()
	} else {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid user"})
	}

	// Get unread notifications
	notifications, err := h.Service.GetUnreadNotifications(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data": notifications,
	})
}

// HandleMarkAsRead, belirtilen bildirimi okundu olarak işaretler.
//
// # Genel Bakış
//
// Bu fonksiyon, URL parametresinden alınan bildirim ID'sine sahip bildirimi okundu
// olarak işaretler. Bildirim servisi üzerinden güncelleme işlemini gerçekleştirir
// ve başarılı sonucu JSON formatında döndürür.
//
// # HTTP Endpoint
//
// - **Method**: PUT veya PATCH
// - **Path**: `/api/notifications/:id/read` (örnek)
// - **Auth**: Gerekli (middleware ile sağlanmalı)
//
// # İstek
//
// - **Headers**:
//   - `Authorization`: Bearer token veya session cookie
// - **URL Params**:
//   - `id`: Bildirim ID'si (uint, pozitif tam sayı)
// - **Body**: Yok
//
// # Yanıt
//
// **Başarılı (200 OK)**:
//
//	```json
//	{
//	  "message": "Notification marked as read"
//	}
//	```
//
// **Hata Yanıtları**:
//
// - **400 Bad Request**: Geçersiz bildirim ID formatı
//
//	```json
//	{"error": "Invalid notification ID"}
//	```
//
// - **500 Internal Server Error**: Veritabanı hatası veya bildirim bulunamadı
//
//	```json
//	{"error": "notification not found"}
//	{"error": "database error message"}
//	```
//
// # Parametreler
//
// - `c`: Fiber context wrapper'ı, URL parametreleri ve HTTP işlemleri için kullanılır
//
// # Döndürür
//
// - `error`: HTTP yanıt hatası veya nil
//
// # Kullanım Senaryoları
//
// 1. **Bildirim Tıklama**: Kullanıcı bir bildirime tıkladığında otomatik okundu işaretleme
// 2. **Manuel İşaretleme**: Kullanıcı bildirim üzerindeki "okundu işaretle" butonuna tıklama
// 3. **Bildirim Detayı**: Bildirim detay sayfası açıldığında otomatik işaretleme
// 4. **Toplu İşlem**: Seçili bildirimleri tek tek okundu işaretleme
//
// # İş Akışı
//
// 1. URL parametresinden ID'yi al (`c.Params("id")`)
// 2. ID'yi string'den uint'e dönüştür
// 3. Dönüştürme hatası varsa 400 döndür
// 4. Service üzerinden bildirimi okundu işaretle
// 5. Hata varsa 500 döndür
// 6. Başarı mesajını JSON formatında döndür
//
// # Güvenlik
//
// - **Yetkilendirme**: Kullanıcı sadece kendi bildirimlerini işaretleyebilmeli
// - **ID Validasyonu**: Pozitif tam sayı kontrolü yapılır
// - **SQL Injection**: Service katmanında parametreli sorgular kullanılmalıdır
// - **IDOR Koruması**: Service katmanında user_id kontrolü yapılmalıdır
//
// # Performans Notları
//
// - Tek satır güncelleme işlemi, hızlı çalışır
// - Index kullanımı önemlidir (id, user_id)
// - Transaction kullanımı gerekmez
//
// # Örnek Kullanım
//
//	```go
//	// Route tanımlama
//	app.Put("/api/notifications/:id/read",
//	    authMiddleware,
//	    handler.HandleMarkAsRead,
//	)
//
//	// Frontend'den çağrı (JavaScript)
//	fetch('/api/notifications/123/read', {
//	    method: 'PUT',
//	    headers: {
//	        'Authorization': 'Bearer ' + token
//	    }
//	})
//	.then(res => res.json())
//	.then(data => {
//	    console.log(data.message);
//	});
//	```
//
// # Avantajlar
//
// - Basit ve hızlı işlem
// - RESTful standartlara uygun
// - Idempotent (tekrar çağrılabilir)
// - Minimal veri transferi
//
// # Dezavantajlar
//
// - Kullanıcı yetkilendirmesi handler'da yapılmıyor (service'e bırakılmış)
// - Bildirim bulunamadığında 404 yerine 500 dönüyor
// - Başarılı güncelleme sonrası bildirim verisi dönmüyor
//
// # Önemli Notlar
//
// - Service katmanında user_id kontrolü yapılmalıdır (IDOR koruması)
// - Bildirim zaten okunmuşsa tekrar işaretleme hata vermez
// - read_at alanı current timestamp ile güncellenir
// - Soft delete kullanılıyorsa deleted_at kontrolü yapılmalıdır
//
// # Güvenlik Uyarıları
//
// ⚠️ **IDOR (Insecure Direct Object Reference) Riski**:
// Handler seviyesinde user_id kontrolü yapılmıyor. Service katmanında
// mutlaka şu kontrol yapılmalıdır:
//
//	```go
//	// Service katmanında
//	func (s *Service) MarkAsRead(notificationID, userID uint) error {
//	    result := s.db.Model(&Notification{}).
//	        Where("id = ? AND user_id = ?", notificationID, userID).
//	        Update("read_at", time.Now())
//	    if result.RowsAffected == 0 {
//	        return errors.New("notification not found or unauthorized")
//	    }
//	    return result.Error
//	}
//	```
//
// # İyileştirme Önerileri
//
// 1. **User ID Kontrolü**: Handler'da user context'i al ve service'e gönder
// 2. **404 Yanıtı**: Bildirim bulunamadığında 404 döndür
// 3. **Yanıt Verisi**: Güncellenmiş bildirim verisini döndür
// 4. **Optimistic Locking**: Version kontrolü ile concurrent update koruması
// 5. **Rate Limiting**: Aynı bildirimi tekrar tekrar işaretlemeyi engelle
//
// # İyileştirilmiş Versiyon Örneği
//
//	```go
//	func (h *NotificationHandler) HandleMarkAsRead(c *context.Context) error {
//	    // Get user from context
//	    user := c.Locals("user")
//	    if user == nil {
//	        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
//	    }
//
//	    var userID uint
//	    if u, ok := user.(interface{ GetID() uint }); ok {
//	        userID = u.GetID()
//	    } else {
//	        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid user"})
//	    }
//
//	    // Get notification ID
//	    idStr := c.Params("id")
//	    id, err := strconv.ParseUint(idStr, 10, 32)
//	    if err != nil {
//	        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid notification ID"})
//	    }
//
//	    // Mark as read with user ID check
//	    notification, err := h.Service.MarkAsRead(uint(id), userID)
//	    if err != nil {
//	        if errors.Is(err, ErrNotificationNotFound) {
//	            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Notification not found"})
//	        }
//	        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
//	    }
//
//	    return c.JSON(fiber.Map{
//	        "message": "Notification marked as read",
//	        "data": notification,
//	    })
//	}
//	```
//
// # Bağımlılıklar
//
// - `context.Context`: Fiber context wrapper
// - `notification.Service`: Bildirim servisi
// - `strconv`: String-uint dönüşümü
func (h *NotificationHandler) HandleMarkAsRead(c *context.Context) error {
	// Get notification ID from params
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid notification ID"})
	}

	// Mark as read
	if err := h.Service.MarkAsRead(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Notification marked as read",
	})
}

// HandleMarkAllAsRead, mevcut kullanıcının tüm bildirimlerini okundu olarak işaretler.
//
// # Genel Bakış
//
// Bu fonksiyon, kimliği doğrulanmış kullanıcının tüm okunmamış bildirimlerini
// toplu olarak okundu durumuna getirir. Kullanıcı deneyimini iyileştirmek için
// tek tıkla tüm bildirimleri temizleme özelliği sağlar.
//
// # HTTP Endpoint
//
// - **Method**: PUT veya PATCH
// - **Path**: `/api/notifications/read-all` (örnek)
// - **Auth**: Gerekli (middleware ile sağlanmalı)
//
// # İstek
//
// - **Headers**:
//   - `Authorization`: Bearer token veya session cookie
// - **Body**: Yok
// - **Query Params**: Yok
//
// # Yanıt
//
// **Başarılı (200 OK)**:
//
//	```json
//	{
//	  "message": "All notifications marked as read"
//	}
//	```
//
// **Hata Yanıtları**:
//
// - **401 Unauthorized**: Kullanıcı kimliği doğrulanmamış
//
//	```json
//	{"error": "Unauthorized"}
//	```
//
// - **500 Internal Server Error**: Veritabanı hatası veya geçersiz kullanıcı
//
//	```json
//	{"error": "Invalid user"}
//	{"error": "database error message"}
//	```
//
// # Parametreler
//
// - `c`: Fiber context wrapper'ı, kullanıcı bilgisi ve HTTP işlemleri için kullanılır
//
// # Döndürür
//
// - `error`: HTTP yanıt hatası veya nil
//
// # Kullanım Senaryoları
//
// 1. **Toplu Temizleme**: Kullanıcı "Tümünü okundu işaretle" butonuna tıklama
// 2. **Bildirim Merkezi**: Bildirim panelini kapatırken otomatik temizleme
// 3. **Periyodik Temizleme**: Belirli aralıklarla eski bildirimleri temizleme
// 4. **Kullanıcı Tercihi**: Kullanıcı ayarlarından toplu işaretleme
// 5. **Mobil Uygulama**: Swipe-to-clear-all özelliği
//
// # İş Akışı
//
// 1. Context'ten user bilgisini al (`c.Locals("user")`)
// 2. User nil kontrolü yap, nil ise 401 döndür
// 3. User'dan ID'yi çıkar (GetID() interface'i kullanarak)
// 4. ID çıkarılamazsa 500 döndür
// 5. Service üzerinden tüm bildirimleri okundu işaretle
// 6. Hata varsa 500 döndür
// 7. Başarı mesajını JSON formatında döndür
//
// # Güvenlik
//
// - **Yetkilendirme**: Middleware tarafından user context'e enjekte edilmelidir
// - **Veri İzolasyonu**: Her kullanıcı sadece kendi bildirimlerini işaretleyebilir
// - **SQL Injection**: Service katmanında parametreli sorgular kullanılmalıdır
// - **Mass Assignment**: Sadece read_at alanı güncellenir
//
// # Performans Notları
//
// - Toplu güncelleme işlemi, çok sayıda bildirimde yavaş olabilir
// - Index kullanımı kritiktir (user_id, read_at)
// - Transaction kullanımı önerilir
// - Batch processing ile performans artırılabilir
// - Async işlem olarak çalıştırılabilir (büyük veri setlerinde)
//
// # Örnek Kullanım
//
//	```go
//	// Route tanımlama
//	app.Put("/api/notifications/read-all",
//	    authMiddleware,
//	    handler.HandleMarkAllAsRead,
//	)
//
//	// Frontend'den çağrı (JavaScript)
//	fetch('/api/notifications/read-all', {
//	    method: 'PUT',
//	    headers: {
//	        'Authorization': 'Bearer ' + token
//	    }
//	})
//	.then(res => res.json())
//	.then(data => {
//	    console.log(data.message);
//	    // UI'da bildirim sayısını sıfırla
//	    updateNotificationBadge(0);
//	});
//	```
//
// # Avantajlar
//
// - Kullanıcı deneyimini iyileştirir (tek tıkla temizleme)
// - Basit ve anlaşılır API
// - RESTful standartlara uygun
// - Idempotent (tekrar çağrılabilir)
// - Minimal veri transferi
//
// # Dezavantajlar
//
// - Çok sayıda bildirimde performans sorunu olabilir
// - Geri alma (undo) özelliği yok
// - İşlem sonrası etkilenen kayıt sayısı dönmüyor
// - Async işlem değil, uzun sürebilir
//
// # Önemli Notlar
//
// - Service katmanında WHERE user_id = ? koşulu mutlaka olmalıdır
// - Zaten okunmuş bildirimleri tekrar işaretleme hata vermez
// - read_at alanı current timestamp ile güncellenir
// - Soft delete kullanılıyorsa deleted_at kontrolü yapılmalıdır
// - Transaction kullanımı önerilir (atomicity için)
//
// # Performans Optimizasyonları
//
// **Veritabanı Seviyesi**:
//
//	```sql
//	-- Index oluşturma
//	CREATE INDEX idx_notifications_user_read ON notifications(user_id, read_at);
//
//	-- Toplu güncelleme sorgusu
//	UPDATE notifications
//	SET read_at = NOW()
//	WHERE user_id = ? AND read_at IS NULL;
//	```
//
// **Async İşlem**:
//
//	```go
//	// Büyük veri setleri için async işlem
//	go func() {
//	    if err := h.Service.MarkAllAsRead(userID); err != nil {
//	        log.Error("Failed to mark all as read", "error", err)
//	    }
//	}()
//	return c.JSON(fiber.Map{
//	    "message": "Marking all notifications as read",
//	    "status": "processing",
//	})
//	```
//
// **Batch Processing**:
//
//	```go
//	// Service katmanında batch processing
//	func (s *Service) MarkAllAsRead(userID uint) error {
//	    batchSize := 1000
//	    for {
//	        result := s.db.Model(&Notification{}).
//	            Where("user_id = ? AND read_at IS NULL", userID).
//	            Limit(batchSize).
//	            Update("read_at", time.Now())
//
//	        if result.Error != nil {
//	            return result.Error
//	        }
//
//	        if result.RowsAffected < int64(batchSize) {
//	            break
//	        }
//	    }
//	    return nil
//	}
//	```
//
// # İyileştirme Önerileri
//
// 1. **Etkilenen Kayıt Sayısı**: Kaç bildirimin işaretlendiğini döndür
// 2. **Async İşlem**: Büyük veri setlerinde background job kullan
// 3. **Geri Alma**: Undo özelliği ekle (son işlem timestamp'i sakla)
// 4. **Filtreleme**: Belirli tarih aralığı veya tip için toplu işaretleme
// 5. **Rate Limiting**: Kötüye kullanımı önlemek için rate limit ekle
// 6. **Progress Tracking**: Uzun işlemlerde progress bar göster
//
// # İyileştirilmiş Versiyon Örneği
//
//	```go
//	func (h *NotificationHandler) HandleMarkAllAsRead(c *context.Context) error {
//	    // Get user from context
//	    user := c.Locals("user")
//	    if user == nil {
//	        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
//	    }
//
//	    var userID uint
//	    if u, ok := user.(interface{ GetID() uint }); ok {
//	        userID = u.GetID()
//	    } else {
//	        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid user"})
//	    }
//
//	    // Mark all as read and get count
//	    count, err := h.Service.MarkAllAsReadWithCount(userID)
//	    if err != nil {
//	        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
//	    }
//
//	    return c.JSON(fiber.Map{
//	        "message": "All notifications marked as read",
//	        "count": count,
//	        "timestamp": time.Now(),
//	    })
//	}
//	```
//
// # Test Senaryoları
//
// 1. **Başarılı İşlem**: Okunmamış bildirimleri işaretleme
// 2. **Boş Liste**: Hiç okunmamış bildirim yokken çağrı
// 3. **Unauthorized**: User context'i olmadan çağrı
// 4. **Invalid User**: Geçersiz user objesi ile çağrı
// 5. **Database Error**: Veritabanı bağlantı hatası
// 6. **Concurrent Requests**: Aynı anda birden fazla istek
// 7. **Large Dataset**: 10000+ bildirim ile performans testi
//
// # Kullanım Örnekleri
//
// **React Component**:
//
//	```jsx
//	function NotificationCenter() {
//	    const markAllAsRead = async () => {
//	        try {
//	            const response = await fetch('/api/notifications/read-all', {
//	                method: 'PUT',
//	                headers: {
//	                    'Authorization': `Bearer ${token}`
//	                }
//	            });
//	            const data = await response.json();
//	            toast.success(data.message);
//	            refreshNotifications();
//	        } catch (error) {
//	            toast.error('Failed to mark all as read');
//	        }
//	    };
//
//	    return (
//	        <button onClick={markAllAsRead}>
//	            Tümünü Okundu İşaretle
//	        </button>
//	    );
//	}
//	```
//
// **Vue Component**:
//
//	```vue
//	<template>
//	    <button @click="markAllAsRead">
//	        Tümünü Okundu İşaretle
//	    </button>
//	</template>
//
//	<script>
//	export default {
//	    methods: {
//	        async markAllAsRead() {
//	            try {
//	                const response = await this.$http.put('/api/notifications/read-all');
//	                this.$toast.success(response.data.message);
//	                this.$emit('refresh');
//	            } catch (error) {
//	                this.$toast.error('İşlem başarısız');
//	            }
//	        }
//	    }
//	}
//	</script>
//	```
//
// # Bağımlılıklar
//
// - `context.Context`: Fiber context wrapper
// - `notification.Service`: Bildirim servisi
// - Auth Middleware: User bilgisini context'e ekleyen middleware
//
// # İlgili Fonksiyonlar
//
// - `HandleGetUnreadNotifications`: Okunmamış bildirimleri getir
// - `HandleMarkAsRead`: Tekil bildirim işaretle
// - `notification.Service.MarkAllAsRead`: Service katmanı implementasyonu
func (h *NotificationHandler) HandleMarkAllAsRead(c *context.Context) error {
	// Get user from context
	user := c.Locals("user")
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Extract user ID
	var userID uint
	if u, ok := user.(interface{ GetID() uint }); ok {
		userID = u.GetID()
	} else {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid user"})
	}

	// Mark all as read
	if err := h.Service.MarkAllAsRead(userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "All notifications marked as read",
	})
}
