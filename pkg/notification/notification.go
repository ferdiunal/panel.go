package notification

// Bu type alias, kullanıcı arayüzünde gösterilecek bildirim türlerini temsil eder.
// Desteklenen değerler: "success", "error", "warning", "info"
//
// Kullanım Senaryoları:
// - Form gönderimi sonuçlarını göstermek
// - Sistem hataları ve uyarılarını bildirmek
// - İşlem başarısını onaylamak
// - Bilgilendirici mesajlar göstermek
//
// Örnek:
//   var notifType notification.Type = notification.TypeSuccess
type Type string

// Bu sabitler, desteklenen tüm bildirim türlerini tanımlar.
// Her bir tür, kullanıcı arayüzünde farklı stil ve ikon ile gösterilir.
const (
	// TypeSuccess, başarılı işlemleri göstermek için kullanılır.
	// Genellikle yeşil renk ve başarı ikonu ile gösterilir.
	// Örnek: "Kullanıcı başarıyla oluşturuldu"
	TypeSuccess Type = "success"

	// TypeError, hata durumlarını göstermek için kullanılır.
	// Genellikle kırmızı renk ve hata ikonu ile gösterilir.
	// Örnek: "Veritabanı bağlantısı başarısız"
	TypeError Type = "error"

	// TypeWarning, uyarı mesajlarını göstermek için kullanılır.
	// Genellikle sarı/turuncu renk ve uyarı ikonu ile gösterilir.
	// Örnek: "Bu işlem geri alınamaz"
	TypeWarning Type = "warning"

	// TypeInfo, bilgilendirici mesajları göstermek için kullanılır.
	// Genellikle mavi renk ve bilgi ikonu ile gösterilir.
	// Örnek: "Veriler yükleniyor..."
	TypeInfo Type = "info"
)

// Bu yapı, kullanıcıya gösterilecek bir bildirimi temsil eder.
// Bildirim, bir mesaj, tür ve görüntülenme süresinden oluşur.
//
// Kullanım Senaryoları:
// - API yanıtlarında kullanıcı geri bildirimi sağlamak
// - Frontend'de toast/alert bileşenleri için veri taşımak
// - Kullanıcı işlemlerinin sonuçlarını iletmek
//
// JSON Serileştirmesi:
// Yapı, JSON formatında serileştirilebilir ve API yanıtlarında gönderilebilir.
//
// Örnek JSON:
//   {
//     "message": "İşlem başarıyla tamamlandı",
//     "type": "success",
//     "duration": 3000
//   }
type Notification struct {
	// Message, kullanıcıya gösterilecek bildirim metnidir.
	// Açık, anlaşılır ve kısa olmalıdır (idealde 100 karakterden az).
	// Örnek: "Profil güncellendi"
	Message string `json:"message"`

	// Type, bildirimin türünü belirtir (success, error, warning, info).
	// Bu değer, frontend'de uygun stil ve ikon seçmek için kullanılır.
	Type Type `json:"type"`

	// Duration, bildirimin ekranda gösterilme süresini milisaniye cinsinden belirtir.
	// Varsayılan değer 3000ms (3 saniye) olarak ayarlanır.
	// Önemli bildirimler için daha uzun, basit bilgiler için daha kısa olabilir.
	// Örnek: 5000 = 5 saniye, 0 = manuel kapatma gerekli
	Duration int `json:"duration"`
}

// Bu fonksiyon, belirtilen mesaj ve tür ile yeni bir Notification nesnesi oluşturur.
// Varsayılan olarak 3000ms (3 saniye) görüntülenme süresi ayarlanır.
//
// Parametreler:
//   - message: Kullanıcıya gösterilecek bildirim metni
//   - notifType: Bildirimin türü (TypeSuccess, TypeError, TypeWarning, TypeInfo)
//
// Dönüş Değeri:
//   - *Notification: Yeni oluşturulan Notification yapısının pointer'ı
//
// Kullanım Senaryoları:
//   - Özel bildirim türleri oluşturmak
//   - Dinamik mesajlarla bildirim oluşturmak
//
// Örnek:
//   notif := notification.New("Veriler kaydedildi", notification.TypeSuccess)
//   notif.SetDuration(5000) // 5 saniye göster
//
// Önemli Notlar:
//   - Döndürülen pointer'ı method chaining ile kullanabilirsiniz
//   - SetDuration() metodu ile görüntülenme süresini değiştirebilirsiniz
func New(message string, notifType Type) *Notification {
	return &Notification{
		Message:  message,
		Type:     notifType,
		Duration: 3000, // Varsayılan 3 saniye
	}
}

// Bu fonksiyon, başarı türünde bir Notification nesnesi oluşturur.
// Başarılı işlemleri kullanıcıya bildirmek için kullanılır.
//
// Parametreler:
//   - message: Başarı mesajı (örn: "Kullanıcı oluşturuldu")
//
// Dönüş Değeri:
//   - *Notification: TypeSuccess türünde yeni Notification pointer'ı
//
// Kullanım Senaryoları:
//   - Form gönderimi başarılı olduğunda
//   - Veri kaydı tamamlandığında
//   - İşlem başarıyla gerçekleştiğinde
//
// Örnek:
//   notif := notification.Success("Profil başarıyla güncellendi")
//   // Döndürür: TypeSuccess türünde Notification
//
// Önemli Notlar:
//   - Varsayılan görüntülenme süresi 3000ms'dir
//   - Method chaining ile SetDuration() kullanabilirsiniz
func Success(message string) *Notification {
	return New(message, TypeSuccess)
}

// Bu fonksiyon, hata türünde bir Notification nesnesi oluşturur.
// Hata durumlarını kullanıcıya bildirmek için kullanılır.
//
// Parametreler:
//   - message: Hata mesajı (örn: "Veritabanı bağlantısı başarısız")
//
// Dönüş Değeri:
//   - *Notification: TypeError türünde yeni Notification pointer'ı
//
// Kullanım Senaryoları:
//   - Veritabanı işlemleri başarısız olduğunda
//   - API çağrıları hata döndürdüğünde
//   - Doğrulama hataları oluştuğunda
//
// Örnek:
//   notif := notification.Error("E-posta zaten kayıtlı")
//   // Döndürür: TypeError türünde Notification
//
// Önemli Notlar:
//   - Hata mesajları kullanıcı dostu olmalıdır
//   - Teknik detaylar yerine çözüm önerileri sunun
//   - Varsayılan görüntülenme süresi 3000ms'dir
func Error(message string) *Notification {
	return New(message, TypeError)
}

// Bu fonksiyon, uyarı türünde bir Notification nesnesi oluşturur.
// Dikkat çekmesi gereken durumları kullanıcıya bildirmek için kullanılır.
//
// Parametreler:
//   - message: Uyarı mesajı (örn: "Bu işlem geri alınamaz")
//
// Dönüş Değeri:
//   - *Notification: TypeWarning türünde yeni Notification pointer'ı
//
// Kullanım Senaryoları:
//   - Tehlikeli işlemler gerçekleştirilmeden önce
//   - Sistem kaynakları sınırlandığında
//   - Kullanıcı dikkatini çekmesi gereken durumlar
//
// Örnek:
//   notif := notification.Warning("Tüm veriler silinecek. Devam etmek istiyor musunuz?")
//   // Döndürür: TypeWarning türünde Notification
//
// Önemli Notlar:
//   - Uyarı mesajları açık ve anlaşılır olmalıdır
//   - Kullanıcıya seçenek sunmayı düşünün
//   - Varsayılan görüntülenme süresi 3000ms'dir
func Warning(message string) *Notification {
	return New(message, TypeWarning)
}

// Bu fonksiyon, bilgi türünde bir Notification nesnesi oluşturur.
// Bilgilendirici mesajları kullanıcıya göstermek için kullanılır.
//
// Parametreler:
//   - message: Bilgi mesajı (örn: "Veriler yükleniyor...")
//
// Dönüş Değeri:
//   - *Notification: TypeInfo türünde yeni Notification pointer'ı
//
// Kullanım Senaryoları:
//   - İşlem devam ederken durum güncellemeleri
//   - Sistem bilgileri göstermek
//   - Kullanıcıya rehberlik mesajları
//
// Örnek:
//   notif := notification.Info("Veriler yükleniyor, lütfen bekleyin...")
//   // Döndürür: TypeInfo türünde Notification
//
// Önemli Notlar:
//   - Bilgi mesajları kısa ve öz olmalıdır
//   - Uzun işlemler için daha uzun Duration ayarlayabilirsiniz
//   - Varsayılan görüntülenme süresi 3000ms'dir
func Info(message string) *Notification {
	return New(message, TypeInfo)
}

// Bu metod, bildirimin ekranda gösterilme süresini ayarlar.
// Method chaining destekler, böylece birden fazla ayarı zincirleme yapabilirsiniz.
//
// Parametreler:
//   - duration: Görüntülenme süresi milisaniye cinsinden (örn: 5000 = 5 saniye)
//
// Dönüş Değeri:
//   - *Notification: Yapılandırılmış Notification pointer'ı (method chaining için)
//
// Kullanım Senaryoları:
//   - Varsayılan 3 saniyeden daha uzun gösterim süresi istediğinde
//   - Önemli bildirimler için daha uzun süre ayarlamak
//   - Hızlı bilgiler için daha kısa süre ayarlamak
//
// Örnek - Basit Kullanım:
//   notif := notification.Success("Kaydedildi")
//   notif.SetDuration(5000) // 5 saniye göster
//
// Örnek - Method Chaining:
//   notif := notification.New("İşlem tamamlandı", notification.TypeSuccess).
//     SetDuration(7000)
//
// Örnek - Farklı Süreler:
//   // Hızlı bilgi
//   notification.Info("Yükleniyor").SetDuration(1000)
//
//   // Önemli uyarı
//   notification.Warning("Dikkat!").SetDuration(10000)
//
//   // Manuel kapatma gerekli
//   notification.Error("Kritik hata").SetDuration(0)
//
// Önemli Notlar:
//   - Duration = 0 ayarlanırsa, bildirim manuel olarak kapatılmalıdır
//   - Çok uzun süreler (>10000ms) kullanıcı deneyimini olumsuz etkileyebilir
//   - Frontend'de bu değer milisaniye cinsinden kullanılır
//   - Method chaining ile diğer ayarlamalar yapılabilir (gelecek genişletmeler için)
func (n *Notification) SetDuration(duration int) *Notification {
	n.Duration = duration
	return n
}
