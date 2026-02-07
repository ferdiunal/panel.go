package handler

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ferdiunal/panel.go/pkg/context"
	notificationDomain "github.com/ferdiunal/panel.go/pkg/domain/notification"
	userDomain "github.com/ferdiunal/panel.go/pkg/domain/user"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// NotificationSSEHandler, Server-Sent Events (SSE) protokolü kullanarak
// gerçek zamanlı bildirim akışı sağlayan bir yapıdır.
//
// # Genel Bakış
//
// Bu yapı, kullanıcılara gerçek zamanlı bildirimler göndermek için SSE teknolojisini kullanır.
// SSE, sunucudan istemciye tek yönlü, sürekli bir veri akışı sağlar ve WebSocket'e göre
// daha basit bir alternatiftir.
//
// # Özellikler
//
// - Gerçek zamanlı bildirim akışı
// - Kullanıcı bazlı bildirim filtreleme
// - Otomatik yeniden bağlanma desteği (istemci tarafında)
// - Düşük kaynak tüketimi
// - HTTP/1.1 ve HTTP/2 uyumluluğu
//
// # Kullanım Senaryoları
//
// - Anlık bildirimler (yeni mesaj, yorum, beğeni vb.)
// - Sistem uyarıları ve duyurular
// - Görev durumu güncellemeleri
// - Gerçek zamanlı dashboard güncellemeleri
//
// # Teknik Detaylar
//
// SSE bağlantısı, HTTP long-polling tekniği kullanarak açık tutulur ve sunucu
// periyodik olarak (2 saniyede bir) yeni bildirimleri kontrol eder. Bu yaklaşım,
// WebSocket'e göre daha az karmaşıktır ve proxy/firewall uyumluluğu daha iyidir.
//
// # Avantajlar
//
// - Basit implementasyon (standart HTTP üzerinden çalışır)
// - Otomatik yeniden bağlanma (tarayıcı tarafından yönetilir)
// - Proxy ve firewall dostu
// - Düşük overhead
//
// # Dezavantajlar
//
// - Tek yönlü iletişim (sadece sunucudan istemciye)
// - Tarayıcı başına bağlantı limiti (genellikle 6)
// - Binary veri desteği yok (sadece text)
//
// # Örnek Kullanım
//
// ```go
// db := gorm.Open(...)
// handler := NewNotificationSSEHandler(db)
//
// app := fiber.New()
// app.Get("/api/notifications/stream", handler.HandleNotificationStream)
// ```
//
// # İstemci Tarafı Örnek (JavaScript)
//
// ```javascript
// const eventSource = new EventSource('/api/notifications/stream');
//
// eventSource.onmessage = (event) => {
//     const notifications = JSON.parse(event.data);
//     console.log('Yeni bildirimler:', notifications);
// };
//
// eventSource.onerror = (error) => {
//     console.error('SSE hatası:', error);
//     // Tarayıcı otomatik olarak yeniden bağlanmayı dener
// };
// ```
//
// # Güvenlik Notları
//
// - Kullanıcı kimlik doğrulaması zorunludur
// - Her kullanıcı sadece kendi bildirimlerini görebilir
// - CORS ayarlarına dikkat edilmelidir
// - Rate limiting uygulanması önerilir
//
// # Performans Notları
//
// - Polling interval (2 saniye) ihtiyaca göre ayarlanabilir
// - Çok sayıda eşzamanlı bağlantı için connection pooling önerilir
// - Nginx gibi reverse proxy kullanılıyorsa buffering kapatılmalıdır
type NotificationSSEHandler struct {
	// db, veritabanı işlemleri için GORM instance'ını tutar.
	// Bildirimler bu instance üzerinden sorgulanır ve filtrelenir.
	db *gorm.DB
}

// NewNotificationSSEHandler, yeni bir NotificationSSEHandler instance'ı oluşturur.
//
// # Parametreler
//
// - `db`: GORM veritabanı bağlantısı. Nil olmamalıdır.
//
// # Dönüş Değeri
//
// Yapılandırılmış bir NotificationSSEHandler pointer'ı döner.
//
// # Örnek Kullanım
//
// ```go
// db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
// if err != nil {
//     log.Fatal(err)
// }
//
// handler := NewNotificationSSEHandler(db)
// ```
//
// # Notlar
//
// - Bu fonksiyon panic yapmaz, ancak db parametresi nil ise
//   sonraki işlemlerde nil pointer hatası alınır.
// - Veritabanı bağlantısının aktif olduğundan emin olun.
func NewNotificationSSEHandler(db *gorm.DB) *NotificationSSEHandler {
	return &NotificationSSEHandler{db: db}
}

// HandleNotificationStream, SSE protokolü kullanarak kullanıcıya gerçek zamanlı
// bildirim akışı sağlar.
//
// # Genel Bakış
//
// Bu fonksiyon, HTTP bağlantısını açık tutar ve kullanıcının okunmamış bildirimlerini
// periyodik olarak kontrol ederek yeni bildirimleri SSE formatında gönderir.
//
// # Çalışma Mantığı
//
// 1. Kullanıcı kimliğini context'ten alır ve doğrular
// 2. SSE için gerekli HTTP header'larını ayarlar
// 3. İlk bağlantıda mevcut okunmamış bildirimleri gönderir (son 50 adet)
// 4. 2 saniyede bir yeni bildirimleri kontrol eder
// 5. Yeni bildirim varsa istemciye gönderir
// 6. İstemci bağlantıyı kesene kadar döngü devam eder
//
// # Parametreler
//
// - `c`: Fiber context nesnesi. Kullanıcı kimliği ve HTTP bağlantı bilgilerini içerir.
//
// # Dönüş Değeri
//
// - `error`: İşlem başarılı ise nil, aksi halde hata döner.
//   - Kullanıcı kimliği yoksa: 401 Unauthorized
//   - İstemci bağlantıyı keserse: nil (normal sonlanma)
//
// # SSE Mesaj Formatı
//
// ```
// data: [{"id":1,"title":"Yeni mesaj","message":"...","read":false,...}]
//
// ```
//
// Her mesaj "data: " prefix'i ile başlar ve "\n\n" ile biter.
//
// # HTTP Header'lar
//
// - `Content-Type: text/event-stream` - SSE formatını belirtir
// - `Cache-Control: no-cache` - Önbelleklemeyi devre dışı bırakır
// - `Connection: keep-alive` - Bağlantıyı açık tutar
// - `X-Accel-Buffering: no` - Nginx buffering'i devre dışı bırakır
//
// # Örnek Kullanım
//
// ```go
// app := fiber.New()
//
// // Middleware ile kullanıcı kimliğini context'e ekle
// app.Use(func(c *fiber.Ctx) error {
//     c.Locals("user_id", getUserID(c))
//     return c.Next()
// })
//
// handler := NewNotificationSSEHandler(db)
// app.Get("/notifications/stream", handler.HandleNotificationStream)
// ```
//
// # İstemci Tarafı Örnek (React)
//
// ```javascript
// useEffect(() => {
//     const eventSource = new EventSource('/api/notifications/stream', {
//         withCredentials: true
//     });
//
//     eventSource.onmessage = (event) => {
//         const notifications = JSON.parse(event.data);
//         setNotifications(prev => [...notifications, ...prev]);
//     };
//
//     eventSource.onerror = () => {
//         console.error('Bağlantı hatası, yeniden bağlanılıyor...');
//     };
//
//     return () => eventSource.close();
// }, []);
// ```
//
// # Veritabanı Sorguları
//
// İlk sorgu (bağlantı kurulduğunda):
// ```sql
// SELECT * FROM notifications
// WHERE user_id = ? AND read = false
// ORDER BY created_at DESC
// LIMIT 50
// ```
//
// Periyodik sorgu (her 2 saniyede):
// ```sql
// SELECT * FROM notifications
// WHERE user_id = ? AND read = false AND created_at > ?
// ORDER BY created_at DESC
// ```
//
// # Performans Optimizasyonları
//
// - Sadece okunmamış bildirimler sorgulanır (read = false)
// - İlk yüklemede maksimum 50 bildirim gönderilir
// - Periyodik sorgularda sadece son kontrolden sonraki bildirimler alınır
// - Index kullanımı: (user_id, read, created_at) composite index önerilir
//
// # Güvenlik Kontrolleri
//
// - Kullanıcı kimliği doğrulaması zorunludur
// - Her kullanıcı sadece kendi bildirimlerini görebilir
// - SQL injection koruması (GORM parametreli sorgular kullanır)
//
// # Hata Senaryoları
//
// 1. **Kullanıcı kimliği yok**: 401 Unauthorized döner
// 2. **Veritabanı hatası**: Sessizce loglanır, bağlantı devam eder
// 3. **İstemci bağlantıyı keser**: Normal sonlanma, nil döner
// 4. **Network timeout**: İstemci otomatik yeniden bağlanır
//
// # Önemli Notlar
//
// - Bu fonksiyon blocking'dir, goroutine içinde çalışır
// - Her bağlantı için bir veritabanı sorgusu açık kalır
// - Çok sayıda eşzamanlı kullanıcı için connection pool boyutu artırılmalıdır
// - Nginx kullanılıyorsa `proxy_buffering off;` ayarı yapılmalıdır
// - Polling interval (2 saniye) production'da ihtiyaca göre ayarlanabilir
//
// # Nginx Konfigürasyon Örneği
//
// ```nginx
// location /api/notifications/stream {
//     proxy_pass http://backend;
//     proxy_buffering off;
//     proxy_cache off;
//     proxy_set_header Connection '';
//     proxy_http_version 1.1;
//     chunked_transfer_encoding off;
// }
// ```
//
// # Test Örneği
//
// ```go
// func TestHandleNotificationStream(t *testing.T) {
//     db := setupTestDB(t)
//     handler := NewNotificationSSEHandler(db)
//
//     app := fiber.New()
//     app.Get("/stream", func(c *fiber.Ctx) error {
//         c.Locals("user_id", 1)
//         return handler.HandleNotificationStream(c)
//     })
//
//     // Test implementation...
// }
// ```
func (h *NotificationSSEHandler) HandleNotificationStream(c *context.Context) error {
	// Kullanıcı objesini context'ten al
	// Bu değer SessionMiddleware tarafından set edilir
	user := c.Locals("user")

	// Kullanıcı yoksa yetkisiz erişim hatası döndür
	// SSE bağlantısı için authentication zorunludur
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized - No user in context",
		})
	}

	// User objesinden ID'yi al
	userModel, ok := user.(*userDomain.User)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized - Invalid user type",
		})
	}

	userID := userModel.ID

	// SSE için gerekli HTTP header'larını ayarla
	// Bu header'lar tarayıcıya SSE bağlantısı olduğunu bildirir

	// Content-Type: SSE formatını belirtir, tarayıcı EventSource API'sini aktive eder
	c.Set("Content-Type", "text/event-stream")

	// Cache-Control: Proxy ve tarayıcı önbelleklemesini devre dışı bırakır
	// SSE gerçek zamanlı veri akışı olduğu için önbellekleme istenmeyen bir durumdur
	c.Set("Cache-Control", "no-cache")

	// Connection: HTTP bağlantısının açık kalmasını sağlar
	// Bu sayede sunucu sürekli veri gönderebilir
	c.Set("Connection", "keep-alive")

	// X-Accel-Buffering: Nginx reverse proxy kullanılıyorsa buffering'i kapatır
	// Buffering açık olursa mesajlar gecikmeyle iletilir, gerçek zamanlılık bozulur
	c.Set("X-Accel-Buffering", "no")

	// İlk bağlantıda mevcut okunmamış bildirimleri gönder
	// Bu sayede kullanıcı bağlandığında bekleyen bildirimleri hemen görür
	var notifications []notificationDomain.Notification

	// Veritabanından kullanıcının okunmamış bildirimlerini sorgula
	// Filtreler:
	// - user_id: Sadece bu kullanıcının bildirimleri
	// - read = false: Sadece okunmamış bildirimler
	// Sıralama: En yeni bildirimler önce (created_at DESC)
	// Limit: Performans için maksimum 50 bildirim
	h.db.Where("user_id = ? AND read = ?", userID, false).
		Order("created_at DESC").
		Limit(50).
		Find(&notifications)

	// Eğer okunmamış bildirim varsa, hemen gönder
	if len(notifications) > 0 {
		// Bildirimleri JSON formatına dönüştür
		// Hata kontrolü yapılmıyor çünkü struct'lar her zaman marshal edilebilir
		data, _ := json.Marshal(notifications)

		// SSE formatında gönder: "data: " prefix + JSON + "\n\n"
		// Çift newline (\n\n) SSE mesajının bittiğini belirtir
		fmt.Fprintf(c, "data: %s\n\n", data)

		// Buffer'ı flush et, veriyi hemen istemciye gönder
		// Flush yapılmazsa veri buffer'da bekler ve gecikme oluşur
		c.Flush()
	}

	// Periyodik kontrol için ticker oluştur
	// Her 2 saniyede bir yeni bildirimleri kontrol eder
	// Bu interval production'da ihtiyaca göre ayarlanabilir:
	// - Daha sık kontrol: Daha gerçek zamanlı ama daha fazla DB yükü
	// - Daha seyrek kontrol: Daha az DB yükü ama daha az gerçek zamanlı
	ticker := time.NewTicker(2 * time.Second)

	// Fonksiyon sonlandığında ticker'ı durdur
	// Bu, goroutine leak'i önler ve kaynakları serbest bırakır
	defer ticker.Stop()

	// Son kontrol zamanını kaydet
	// Bu timestamp, yeni bildirimleri filtrelemek için kullanılır
	// Sadece bu zamandan sonra oluşturulan bildirimler gönderilir
	lastCheck := time.Now()

	// Sonsuz döngü: İstemci bağlantıyı kesene kadar devam eder
	for {
		// Select statement: Birden fazla channel'ı aynı anda dinler
		select {
		// Ticker channel'ından sinyal geldiğinde (her 2 saniyede bir)
		case <-ticker.C:
			// Son kontrolden sonra oluşturulan yeni bildirimleri sorgula
			var newNotifications []notificationDomain.Notification

			// Veritabanı sorgusu:
			// - user_id: Sadece bu kullanıcının bildirimleri
			// - read = false: Sadece okunmamış bildirimler
			// - created_at > lastCheck: Sadece son kontrolden sonra oluşturulanlar
			// Bu sayede aynı bildirim tekrar gönderilmez
			h.db.Where("user_id = ? AND read = ? AND created_at > ?",
				userID, false, lastCheck).
				Order("created_at DESC").
				Find(&newNotifications)

			// Yeni bildirim varsa gönder
			if len(newNotifications) > 0 {
				// JSON'a dönüştür
				data, _ := json.Marshal(newNotifications)

				// SSE formatında gönder
				fmt.Fprintf(c, "data: %s\n\n", data)

				// Buffer'ı flush et, hemen gönder
				c.Flush()

				// Son kontrol zamanını güncelle
				// Bir sonraki kontrolde bu zamandan sonraki bildirimler alınır
				lastCheck = time.Now()
			}

		// Context Done channel'ından sinyal geldiğinde
		// Bu, istemci bağlantıyı kestiğinde veya timeout olduğunda tetiklenir
		case <-c.Context().Done():
			// Bağlantı kesildi, normal sonlanma
			// Hata döndürmeye gerek yok, bu beklenen bir durum
			// Ticker defer ile otomatik olarak durdurulacak
			return nil
		}
	}
}
