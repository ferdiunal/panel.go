// Package domain, bildirim (notification) domain modellerini ve işlemlerini içerir.
// Bu paket, kullanıcı bildirimleri ile ilgili tüm veri yapılarını ve iş mantığını tanımlar.
package domain

import (
	"time"

	"gorm.io/gorm"
)

// Bu tür alias, bildirim türlerini temsil eder.
// Bildirim sisteminde farklı kategorilerdeki mesajları ayırt etmek için kullanılır.
// Desteklenen değerler: "success", "error", "warning", "info"
//
// Kullanım Senaryoları:
// - Başarılı işlemler için "success" türü
// - Hata durumları için "error" türü
// - Uyarı mesajları için "warning" türü
// - Bilgilendirme mesajları için "info" türü
//
// Örnek:
//   var notifType NotificationType = NotificationTypeSuccess
//   notification.Type = notifType
type NotificationType string

// Bu sabitler, bildirim türlerinin önceden tanımlanmış değerlerini içerir.
// Veritabanında tutarlılık sağlamak için bu sabitler kullanılmalıdır.
//
// NotificationTypeSuccess: Başarılı işlem tamamlandığında kullanılır (örn: "Kayıt başarılı")
// NotificationTypeError: Hata oluştuğunda kullanılır (örn: "Giriş başarısız")
// NotificationTypeWarning: Uyarı mesajları için kullanılır (örn: "Süresi dolmak üzere")
// NotificationTypeInfo: Genel bilgilendirme için kullanılır (örn: "Yeni güncelleme mevcut")
const (
	NotificationTypeSuccess NotificationType = "success"
	NotificationTypeError   NotificationType = "error"
	NotificationTypeWarning NotificationType = "warning"
	NotificationTypeInfo    NotificationType = "info"
)

// Bu yapı, kullanıcı bildirimleri için veritabanı modelini temsil eder.
// Sistem tarafından oluşturulan ve kullanıcılara gösterilen tüm bildirimleri saklar.
//
// Kullanım Senaryoları:
// - Kullanıcı işlemlerinin sonuçlarını bildirmek
// - Sistem olaylarını kullanıcılara iletmek
// - Okunmuş/okunmamış bildirimleri takip etmek
// - Bildirim geçmişini tutmak
//
// Önemli Notlar:
// - UserID nil olabilir (sistem genelinde bildirimleri temsil eder)
// - Duration, bildirim ekranında gösterilme süresini belirler (milisaniye cinsinden)
// - ReadAt, bildirim okunduğu zaman kaydedilir
// - DeletedAt, soft delete için GORM tarafından kullanılır
//
// Veritabanı Özellikleri:
// - ID: Birincil anahtar, otomatik artan
// - UserID: İndekslenmiş, nullable (sistem bildirimleri için)
// - Message: Metin alanı, zorunlu
// - Type: Varchar(20), varsayılan değer 'info'
// - Duration: Varsayılan 3000ms (3 saniye)
// - Read: Varsayılan false (okunmamış)
type Notification struct {
	// ID, bildirimin benzersiz tanımlayıcısıdır.
	// Veritabanında birincil anahtar olarak kullanılır.
	// Otomatik olarak artan bir değerdir.
	ID uint `gorm:"primaryKey" json:"id"`

	// UserID, bildirimin ait olduğu kullanıcının kimliğidir.
	// Nil olabilir - bu durumda bildirim sistem genelinde bir bildirimdir.
	// Veritabanında indekslenmiştir (hızlı sorgu için).
	// JSON çıktısında "user_id" olarak gösterilir.
	UserID *uint `gorm:"index" json:"user_id"`

	// Message, bildirimin ana içeriğidir.
	// Kullanıcıya gösterilecek metin mesajını içerir.
	// Metin alanı olarak saklanır (uzun metinleri destekler).
	// Zorunlu alan - boş olamaz.
	// JSON çıktısında "message" olarak gösterilir.
	Message string `gorm:"type:text;not null" json:"message"`

	// Type, bildirimin kategorisini belirler.
	// Desteklenen değerler: success, error, warning, info
	// Varchar(20) olarak saklanır.
	// Varsayılan değer: 'info'
	// Zorunlu alan - boş olamaz.
	// JSON çıktısında "type" olarak gösterilir.
	Type NotificationType `gorm:"type:varchar(20);not null;default:'info'" json:"type"`

	// Duration, bildirimin kullanıcı arayüzünde gösterilme süresini belirler.
	// Milisaniye (ms) cinsinden değer alır.
	// Varsayılan değer: 3000ms (3 saniye)
	// Örnek: 5000 = 5 saniye, 0 = manuel kapatma gerekli
	// JSON çıktısında "duration" olarak gösterilir.
	Duration int `gorm:"default:3000" json:"duration"`

	// Read, bildirimin okunup okunmadığını belirten boolean değerdir.
	// true: Bildirim kullanıcı tarafından okunmuş
	// false: Bildirim henüz okunmamış
	// Varsayılan değer: false
	// JSON çıktısında "read" olarak gösterilir.
	Read bool `gorm:"default:false" json:"read"`

	// ReadAt, bildirimin okunduğu zamanı kaydeder.
	// Nil olabilir - bildirim okunmadıysa nil kalır.
	// MarkAsRead() metodu çağrıldığında otomatik olarak ayarlanır.
	// JSON çıktısında "read_at" olarak gösterilir.
	ReadAt *time.Time `json:"read_at"`

	// CreatedAt, bildirimin oluşturulduğu zamanı kaydeder.
	// GORM tarafından otomatik olarak ayarlanır.
	// Bildirim geçmişini takip etmek için kullanılır.
	// JSON çıktısında "created_at" olarak gösterilir.
	CreatedAt time.Time `json:"created_at"`

	// UpdatedAt, bildirimin son güncellendiği zamanı kaydeder.
	// GORM tarafından otomatik olarak ayarlanır.
	// Her değişiklikte güncellenir.
	// JSON çıktısında "updated_at" olarak gösterilir.
	UpdatedAt time.Time `json:"updated_at"`

	// DeletedAt, soft delete işlemi için GORM tarafından kullanılır.
	// Nil olabilir - silinmemiş kayıtlar için nil kalır.
	// Veritabanında indekslenmiştir (soft delete sorgularını hızlandırmak için).
	// JSON çıktısında gösterilmez ("-" ile işaretlenmiş).
	// Önemli Not: Silinen bildirimleri sorgulamak için GORM'un Unscoped() metodu kullanılmalıdır.
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// Bu metod, GORM ORM'ine Notification yapısının hangi veritabanı tablosunda saklanacağını söyler.
// GORM varsayılan olarak struct adının çoğul formunu (Notifications) tablo adı olarak kullanır,
// ancak bu metod ile özel bir tablo adı belirtebiliriz.
//
// Dönüş Değeri:
// - string: Veritabanı tablosunun adı ("notifications")
//
// Kullanım Senaryoları:
// - GORM'un otomatik tablo adı belirleme mekanizmasını geçersiz kılmak
// - Özel tablo adlandırma kuralları uygulamak
// - Veritabanı şemasını kontrol etmek
//
// Önemli Notlar:
// - Bu metod GORM tarafından otomatik olarak çağrılır
// - Tablo adı değiştirilirse, veritabanı migration'ları da güncellenmesi gerekir
// - Receiver'ın pointer olması gerekmez (değer receiver yeterlidir)
//
// Örnek:
//   var notif Notification
//   tableName := notif.TableName() // "notifications" döner
func (Notification) TableName() string {
	return "notifications"
}

// Bu metod, bir bildirimi okunmuş olarak işaretler ve veritabanında günceller.
// Bildirim okunduğu zaman kaydedilir ve Read alanı true olarak ayarlanır.
//
// Parametreler:
// - db (*gorm.DB): GORM veritabanı bağlantısı. Veritabanı işlemlerini gerçekleştirmek için kullanılır.
//
// Dönüş Değeri:
// - error: İşlem sırasında oluşan hata. Başarılı olursa nil döner.
//
// Kullanım Senaryoları:
// - Kullanıcı bir bildirimi tıkladığında
// - Bildirim okunmuş olarak işaretlemek gerektiğinde
// - Okunmamış bildirimleri takip etmek için
//
// Yapılan İşlemler:
// 1. Geçerli zamanı alır (time.Now())
// 2. Read alanını true olarak ayarlar
// 3. ReadAt alanını geçerli zaman ile ayarlar
// 4. Değişiklikleri veritabanında kaydeder (db.Save())
//
// Önemli Notlar:
// - Metod receiver'ı pointer (*Notification) olmalıdır çünkü yapıyı değiştirir
// - Veritabanı bağlantısı geçerli olmalıdır, aksi takdirde hata döner
// - Aynı bildirimi birden fazla kez işaretlemek güvenlidir (idempotent)
// - ReadAt zamanı her çağrıda güncellenir
//
// Hata Durumları:
// - Veritabanı bağlantısı kapalıysa hata döner
// - Bildirim veritabanında bulunamazsa hata döner
// - Yazma izni yoksa hata döner
//
// Örnek Kullanım:
//   var notification Notification
//   // ... bildirimi veritabanından yükle ...
//   err := notification.MarkAsRead(db)
//   if err != nil {
//       log.Printf("Bildirim okunmuş olarak işaretlenemedi: %v", err)
//       return err
//   }
//   // Başarılı - bildirim artık okunmuş olarak işaretlendi
//
// Döndürür:
// - Başarılı olursa: nil
// - Hata olursa: GORM hata nesnesi
func (n *Notification) MarkAsRead(db *gorm.DB) error {
	// Geçerli zamanı alır
	now := time.Now()
	// Read alanını true olarak ayarlar (okunmuş olarak işaretler)
	n.Read = true
	// ReadAt alanını geçerli zaman ile ayarlar (okunma zamanını kaydeder)
	n.ReadAt = &now
	// Değişiklikleri veritabanında kaydeder ve hata varsa döner
	return db.Save(n).Error
}
