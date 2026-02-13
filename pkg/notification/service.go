// Bu paket, bildirim (notification) işlemlerini yönetmek için gerekli servisleri içerir.
// Kullanıcılara gönderilen bildirimlerin veritabanında saklanması, okunması ve yönetilmesi
// gibi işlemleri gerçekleştirir.
package notification

import (
	"fmt"

	"github.com/ferdiunal/panel.go/pkg/core"
	"github.com/ferdiunal/panel.go/pkg/data"
	notificationDomain "github.com/ferdiunal/panel.go/pkg/domain/notification"
)

// Bu yapı, bildirim servisi için gerekli bağımlılıkları içerir.
// Service, Provider kullanarak bildirim işlemlerini gerçekleştirir.
//
// Alanlar:
//   - provider: DataProvider instance'ı, tüm bildirim işlemleri için kullanılır
//
// Kullanım Senaryoları:
//   - Kullanıcı işlemleri sonrasında bildirimleri veritabanına kaydetme
//   - Kullanıcının okunmamış bildirimlerini alma
//   - Bildirimleri okundu olarak işaretleme
//   - Tüm bildirimleri toplu olarak okundu olarak işaretleme
//
// Örnek:
//
//	provider := data.NewGormDataProvider(db, &Notification{})
//	service := NewService(provider)
//	err := service.SaveNotifications(ctx)
//	if err != nil {
//	    log.Fatal(err)
//	}
type Service struct {
	provider data.DataProvider
}

// Bu fonksiyon, yeni bir bildirim servisi örneği oluşturur.
// Verilen Provider'ı kullanarak Service yapısını başlatır.
//
// Parametreler:
//   - provider: DataProvider instance'ı, bildirim işlemleri için kullanılacak
//
// Dönüş Değeri:
//   - *Service: Yapılandırılmış Service pointer'ı
//
// Önemli Notlar:
//   - provider parametresi nil olmamalıdır, aksi takdirde runtime hatası oluşur
//   - Döndürülen Service örneği hemen kullanıma hazırdır
//   - Provider pattern kullanılır, ORM'den bağımsız
//
// Örnek:
//
//	provider := data.NewGormDataProvider(db, &Notification{})
//	service := NewService(provider)
func NewService(provider data.DataProvider) *Service {
	return &Service{provider: provider}
}

// Bu metod, ResourceContext içindeki bildirimleri veritabanına kaydeder.
// Context'ten alınan bildirimleri domain modeline dönüştürerek GORM aracılığıyla
// veritabanına kaydeder. Her bildirim için ayrı bir INSERT işlemi gerçekleştirilir.
//
// Parametreler:
//   - ctx: ResourceContext, bildirimleri ve kullanıcı bilgisini içerir
//
// Dönüş Değeri:
//   - error: İşlem başarılı ise nil, aksi takdirde hata mesajı
//
// Kullanım Senaryoları:
//   - Kullanıcı kaydı sonrasında hoşgeldiniz bildirimi kaydetme
//   - İşlem başarısı/başarısızlığı bildirimi kaydetme
//   - Sistem bildirimleri kaydetme
//
// Önemli Notlar:
//   - Context'te bildirim yoksa fonksiyon nil hata döndürür (başarılı)
//   - Kullanıcı bilgisi Context'te yoksa UserID nil olarak kaydedilir
//   - Tüm bildirimlerin Read alanı false olarak başlatılır
//   - Veritabanı hatası durumunda işlem durdurulur ve hata döndürülür
//
// Örnek:
//
//	ctx := &core.ResourceContext{
//	    User: user,
//	    Notifications: []core.Notification{
//	        {Message: "Hoşgeldiniz", Type: "success", Duration: 5000},
//	    },
//	}
//	err := service.SaveNotifications(ctx)
//	if err != nil {
//	    log.Printf("Bildirim kaydedilemedi: %v", err)
//	}
func (s *Service) SaveNotifications(ctx *core.ResourceContext) error {
	notifications := ctx.GetNotifications()
	if len(notifications) == 0 {
		return nil
	}

	// Context'ten kullanıcı ID'sini al (varsa)
	var userID *uint
	if ctx.User != nil {
		if user, ok := ctx.User.(interface{ GetID() uint }); ok {
			id := user.GetID()
			userID = &id
		}
	}

	// Bulk insert için SQL query oluştur
	if len(notifications) == 1 {
		// Tek notification için basit INSERT
		notif := notifications[0]
		err := s.provider.Exec(nil,
			"INSERT INTO notifications (user_id, message, type, duration, read, created_at, updated_at) VALUES (?, ?, ?, ?, ?, NOW(), NOW())",
			userID, notif.Message, notif.Type, notif.Duration, false)
		return err
	}

	// Çoklu notification için bulk INSERT
	values := make([]string, len(notifications))
	args := make([]interface{}, 0, len(notifications)*5)

	for i, notif := range notifications {
		values[i] = "(?, ?, ?, ?, ?, NOW(), NOW())"
		args = append(args, userID, notif.Message, notif.Type, notif.Duration, false)
	}

	sql := "INSERT INTO notifications (user_id, message, type, duration, read, created_at, updated_at) VALUES " +
		fmt.Sprintf("%s", values[0])
	for i := 1; i < len(values); i++ {
		sql += ", " + values[i]
	}

	return s.provider.Exec(nil, sql, args...)
}

// Bu metod, belirtilen kullanıcının okunmamış bildirimlerini veritabanından alır.
// Kullanıcı ID'sine göre filtrelenmiş, okunmamış (read = false) bildirimleri
// en yeni tarihten başlayarak sıralanmış şekilde döndürür.
//
// Parametreler:
//   - userID: Bildirimleri alınacak kullanıcının ID'si
//
// Dönüş Değeri:
//   - []notificationDomain.Notification: Okunmamış bildirimlerin listesi
//   - error: İşlem başarılı ise nil, aksi takdirde hata mesajı
//
// Kullanım Senaryoları:
//   - Kullanıcı paneline giriş yaptığında okunmamış bildirimleri gösterme
//   - Bildirim çanını güncellemek için okunmamış bildirimleri alma
//   - Kullanıcı bildirimleri kontrol ettiğinde listeyi yenileme
//
// Önemli Notlar:
//   - Sonuçlar created_at alanına göre DESC (azalan) sırada döndürülür
//   - Kullanıcının hiç bildirimi yoksa boş slice döndürülür (hata değil)
//   - Veritabanı bağlantı hatası durumunda error döndürülür
//   - Sadece read = false olan bildirimler döndürülür
//
// Örnek:
//
//	notifications, err := service.GetUnreadNotifications(userID)
//	if err != nil {
//	    log.Printf("Bildirimler alınamadı: %v", err)
//	    return
//	}
//	for _, notif := range notifications {
//	    fmt.Printf("Bildirim: %s\n", notif.Message)
//	}
func (s *Service) GetUnreadNotifications(userID uint) ([]notificationDomain.Notification, error) {
	// TODO: Context parametresi eklenebilir
	// Şimdilik nil context ile çalışıyoruz
	results, err := s.provider.Raw(nil,
		"SELECT * FROM notifications WHERE user_id = ? AND read = ? ORDER BY created_at DESC",
		userID, false)
	if err != nil {
		return nil, err
	}

	// Map'leri struct'a dönüştür
	notifications := make([]notificationDomain.Notification, 0, len(results))
	for _, result := range results {
		notif := notificationDomain.Notification{}

		// ID
		if id, ok := result["id"].(uint); ok {
			notif.ID = id
		} else if id, ok := result["id"].(int64); ok {
			notif.ID = uint(id)
		}

		// UserID
		if userID, ok := result["user_id"].(uint); ok {
			notif.UserID = &userID
		} else if userID, ok := result["user_id"].(int64); ok {
			uid := uint(userID)
			notif.UserID = &uid
		}

		// Message
		if msg, ok := result["message"].(string); ok {
			notif.Message = msg
		}

		// Type
		if typ, ok := result["type"].(string); ok {
			notif.Type = notificationDomain.NotificationType(typ)
		}

		// Duration
		if dur, ok := result["duration"].(int); ok {
			notif.Duration = dur
		} else if dur, ok := result["duration"].(int64); ok {
			notif.Duration = int(dur)
		}

		// Read
		if read, ok := result["read"].(bool); ok {
			notif.Read = read
		}

		notifications = append(notifications, notif)
	}

	return notifications, nil
}

// Bu metod, belirtilen ID'ye sahip bildirimi okundu olarak işaretler.
// UPDATE query ile direkt veritabanında güncelleme yapar.
//
// Parametreler:
//   - notificationID: İşaretlenecek bildirimin ID'si
//
// Dönüş Değeri:
//   - error: İşlem başarılı ise nil, aksi takdirde hata mesajı
//
// Kullanım Senaryoları:
//   - Kullanıcı bir bildirimi tıkladığında okundu olarak işaretleme
//   - Bildirim detaylarını görüntülerken okundu durumunu güncelleme
//   - Bildirim geçmişinde okundu olarak işaretleme
//
// Önemli Notlar:
//   - Bildirim bulunamazsa hata döndürülmez (0 satır etkilenir)
//   - Başarılı işlem sonrasında read_at alanı otomatik olarak güncellenir
//   - Zaten okunmuş bir bildirimi tekrar işaretlemek güvenlidir
//
// Örnek:
//
//	err := service.MarkAsRead(notificationID)
//	if err != nil {
//	    log.Printf("Hata: %v", err)
//	    return
//	}
//	log.Println("Bildirim okundu olarak işaretlendi")
func (s *Service) MarkAsRead(notificationID uint) error {
	// TODO: Context parametresi eklenebilir
	return s.provider.Exec(nil,
		"UPDATE notifications SET read = ?, read_at = NOW() WHERE id = ?",
		true, notificationID)
}

// Bu metod, belirtilen kullanıcının tüm okunmamış bildirimlerini toplu olarak
// okundu olarak işaretler. Tek bir UPDATE sorgusu ile tüm bildirimleri günceller,
// bu da performans açısından daha verimlidir.
//
// Parametreler:
//   - userID: Bildirimlerini işaretlenecek kullanıcının ID'si
//
// Dönüş Değeri:
//   - error: İşlem başarılı ise nil, aksi takdirde hata mesajı
//
// Kullanım Senaryoları:
//   - Kullanıcı "Tümünü Okundu Olarak İşaretle" butonuna tıkladığında
//   - Bildirim panelini kapatırken tüm bildirimleri okundu olarak işaretleme
//   - Toplu bildirim yönetimi işlemlerinde
//
// Önemli Notlar:
//   - Sadece read = false olan bildirimler güncellenir
//   - read_at alanı otomatik olarak NOW() (şu anki zaman) olarak ayarlanır
//   - Kullanıcının okunmamış bildirimi yoksa işlem başarılı olur (0 satır etkilenir)
//   - Performans açısından GetUnreadNotifications + MarkAsRead döngüsünden daha iyidir
//
// Örnek:
//
//	err := service.MarkAllAsRead(userID)
//	if err != nil {
//	    log.Printf("Bildirimler işaretlenemedi: %v", err)
//	    return
//	}
//	log.Println("Tüm bildirimler okundu olarak işaretlendi")
func (s *Service) MarkAllAsRead(userID uint) error {
	// TODO: Context parametresi eklenebilir
	return s.provider.Exec(nil,
		"UPDATE notifications SET read = ?, read_at = NOW() WHERE user_id = ? AND read = ?",
		true, userID, false)
}
