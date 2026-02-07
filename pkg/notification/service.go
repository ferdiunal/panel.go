// Bu paket, bildirim (notification) işlemlerini yönetmek için gerekli servisleri içerir.
// Kullanıcılara gönderilen bildirimlerin veritabanında saklanması, okunması ve yönetilmesi
// gibi işlemleri gerçekleştirir.
package notification

import (
	"github.com/ferdiunal/panel.go/pkg/core"
	notificationDomain "github.com/ferdiunal/panel.go/pkg/domain/notification"
	"gorm.io/gorm"
)

// Bu yapı, bildirim servisi için gerekli bağımlılıkları içerir.
// Service, GORM veritabanı bağlantısını kullanarak bildirim işlemlerini gerçekleştirir.
//
// Alanlar:
//   - db: GORM veritabanı bağlantısı, tüm bildirim işlemleri için kullanılır
//
// Kullanım Senaryoları:
//   - Kullanıcı işlemleri sonrasında bildirimleri veritabanına kaydetme
//   - Kullanıcının okunmamış bildirimlerini alma
//   - Bildirimleri okundu olarak işaretleme
//   - Tüm bildirimleri toplu olarak okundu olarak işaretleme
//
// Örnek:
//   service := NewService(db)
//   err := service.SaveNotifications(ctx)
//   if err != nil {
//       log.Fatal(err)
//   }
type Service struct {
	db *gorm.DB
}

// Bu fonksiyon, yeni bir bildirim servisi örneği oluşturur.
// Verilen GORM veritabanı bağlantısını kullanarak Service yapısını başlatır.
//
// Parametreler:
//   - db: GORM veritabanı bağlantısı, bildirim işlemleri için kullanılacak
//
// Dönüş Değeri:
//   - *Service: Yapılandırılmış Service pointer'ı
//
// Önemli Notlar:
//   - db parametresi nil olmamalıdır, aksi takdirde runtime hatası oluşur
//   - Döndürülen Service örneği hemen kullanıma hazırdır
//
// Örnek:
//   db := gorm.Open(sqlite.Open("test.db"))
//   service := NewService(db)
func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
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
//   ctx := &core.ResourceContext{
//       User: user,
//       Notifications: []core.Notification{
//           {Message: "Hoşgeldiniz", Type: "success", Duration: 5000},
//       },
//   }
//   err := service.SaveNotifications(ctx)
//   if err != nil {
//       log.Printf("Bildirim kaydedilemedi: %v", err)
//   }
func (s *Service) SaveNotifications(ctx *core.ResourceContext) error {
	notifications := ctx.GetNotifications()
	if len(notifications) == 0 {
		return nil
	}

	// Context'ten kullanıcı ID'sini al (varsa)
	// GetID() metodunu destekleyen interface'i kontrol et
	var userID *uint
	if ctx.User != nil {
		if user, ok := ctx.User.(interface{ GetID() uint }); ok {
			id := user.GetID()
			userID = &id
		}
	}

	// Context'teki bildirimleri domain modeline dönüştür ve veritabanına kaydet
	for _, notif := range notifications {
		dbNotif := &notificationDomain.Notification{
			UserID:   userID,
			Message:  notif.Message,
			Type:     notificationDomain.NotificationType(notif.Type),
			Duration: notif.Duration,
			Read:     false,
		}

		if err := s.db.Create(dbNotif).Error; err != nil {
			return err
		}
	}

	return nil
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
//   notifications, err := service.GetUnreadNotifications(userID)
//   if err != nil {
//       log.Printf("Bildirimler alınamadı: %v", err)
//       return
//   }
//   for _, notif := range notifications {
//       fmt.Printf("Bildirim: %s\n", notif.Message)
//   }
func (s *Service) GetUnreadNotifications(userID uint) ([]notificationDomain.Notification, error) {
	var notifications []notificationDomain.Notification
	err := s.db.Where("user_id = ? AND read = ?", userID, false).
		Order("created_at DESC").
		Find(&notifications).Error
	return notifications, err
}

// Bu metod, belirtilen ID'ye sahip bildirimi okundu olarak işaretler.
// Önce bildirimi veritabanından bulur, ardından domain modelin MarkAsRead metodunu
// çağırarak okundu durumunu günceller.
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
//   - Bildirim bulunamazsa GORM "record not found" hatası döndürülür
//   - Başarılı işlem sonrasında read_at alanı otomatik olarak güncellenir
//   - Zaten okunmuş bir bildirimi tekrar işaretlemek güvenlidir
//   - Veritabanı hatası durumunda işlem geri alınır
//
// Örnek:
//   err := service.MarkAsRead(notificationID)
//   if err != nil {
//       if errors.Is(err, gorm.ErrRecordNotFound) {
//           log.Println("Bildirim bulunamadı")
//       } else {
//           log.Printf("Hata: %v", err)
//       }
//       return
//   }
//   log.Println("Bildirim okundu olarak işaretlendi")
func (s *Service) MarkAsRead(notificationID uint) error {
	var notif notificationDomain.Notification
	if err := s.db.First(&notif, notificationID).Error; err != nil {
		return err
	}
	return notif.MarkAsRead(s.db)
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
//   - Veritabanı hatası durumunda işlem geri alınır
//   - Performans açısından GetUnreadNotifications + MarkAsRead döngüsünden daha iyidir
//
// Örnek:
//   err := service.MarkAllAsRead(userID)
//   if err != nil {
//       log.Printf("Bildirimler işaretlenemedi: %v", err)
//       return
//   }
//   log.Println("Tüm bildirimler okundu olarak işaretlendi")
func (s *Service) MarkAllAsRead(userID uint) error {
	return s.db.Model(&notificationDomain.Notification{}).
		Where("user_id = ? AND read = ?", userID, false).
		Updates(map[string]interface{}{
			"read":    true,
			"read_at": gorm.Expr("NOW()"),
		}).Error
}
