// Bu paket, kullanıcı oturumlarının (session) yönetimini sağlayan domain katmanını içerir.
// Oturumlar, kullanıcı kimlik doğrulaması ve yetkilendirme işlemleri için kullanılır.
// Paket, Session veri modeli ve Repository arayüzü tanımlar.
package session

import (
	"context"
	"time"

	"github.com/ferdiunal/panel.go/pkg/domain/user"
)

// Bu yapı, kullanıcı oturumunun tüm bilgilerini temsil eder.
// Oturumlar, kullanıcıların sisteme giriş yaptıktan sonra oluşturulan ve
// belirli bir süre boyunca geçerli olan kimlik doğrulama kayıtlarıdır.
//
// Kullanım Senaryoları:
// - Kullanıcı giriş işleminden sonra oturum oluşturma
// - Token tabanlı kimlik doğrulama ve yetkilendirme
// - Oturum süresinin dolup dolmadığını kontrol etme
// - Kullanıcı çıkış işleminde oturum silme
// - Güvenlik denetimi için IP adresi ve User-Agent bilgilerini kaydetme
//
// Önemli Notlar:
// - Token alanı benzersiz (unique) olmalıdır ve her oturum için farklı olmalıdır
// - ExpiresAt zamanı geçmiş oturumlar geçersiz kabul edilir
// - IPAddress ve UserAgent, güvenlik denetimi ve anomali tespiti için kullanılır
// - User ilişkisi, oturumun hangi kullanıcıya ait olduğunu gösterir
type Session struct {
	// ID: Oturumun benzersiz tanımlayıcısı (Primary Key)
	// Veritabanında otomatik olarak artan bir sayıdır
	ID uint `json:"id" gorm:"primaryKey"`

	// Token: Oturumun benzersiz token değeri
	// Kimlik doğrulama işlemlerinde kullanılan güvenli bir string'dir
	// Örnek: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
	// Uyarı: Bu alan benzersiz olmalıdır, aynı token iki oturumda olamaz
	Token string `json:"token" gorm:"uniqueIndex"`

	// UserID: Oturumun ait olduğu kullanıcının ID'si
	// Foreign Key olarak user tablosuna referans verir
	// Örnek: 1, 2, 3 gibi kullanıcı ID'leri
	UserID uint `json:"userId" gorm:"index"`

	// ExpiresAt: Oturumun geçerlilik süresi bitişi zamanı
	// Bu zaman geçtikten sonra oturum otomatik olarak geçersiz sayılır
	// Örnek: 2026-02-07 18:00:00 UTC
	// Uyarı: Sistem saati ile karşılaştırılarak kontrol edilmelidir
	ExpiresAt time.Time `json:"expiresAt" gorm:"index"`

	// IPAddress: Oturumun oluşturulduğu istemcinin IP adresi
	// Güvenlik denetimi ve anomali tespiti için kaydedilir
	// Örnek: "192.168.1.100" veya "2001:0db8:85a3::8a2e:0370:7334"
	// Kullanım: Oturum sırasında IP adresi değişirse güvenlik uyarısı verilebilir
	IPAddress string `json:"ipAddress"`

	// UserAgent: Oturumun oluşturulduğu istemcinin tarayıcı/cihaz bilgisi
	// HTTP User-Agent header'ından alınan değerdir
	// Örnek: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36..."
	// Kullanım: Cihaz değişikliği tespiti ve güvenlik denetimi için kullanılır
	UserAgent string `json:"userAgent"`

	// CreatedAt: Oturumun oluşturulduğu zaman
	// Veritabanı tarafından otomatik olarak ayarlanır
	// Örnek: 2026-02-07 16:00:00 UTC
	// Kullanım: Oturum yaşını hesaplamak ve denetim izleri için
	CreatedAt time.Time `json:"createdAt" gorm:"index"`

	// UpdatedAt: Oturumun son güncellendiği zaman
	// Veritabanı tarafından otomatik olarak ayarlanır
	// Örnek: 2026-02-07 16:30:00 UTC
	// Kullanım: Oturum aktivitesini izlemek için
	UpdatedAt time.Time `json:"updatedAt" gorm:"index"`

	// User: Oturumun ait olduğu kullanıcı nesnesi
	// Foreign Key ilişkisi ile user tablosundan ilgili kullanıcı yüklenir
	// JSON çıktısında omitempty ile belirtilmiştir (boş ise dahil edilmez)
	// Örnek: &user.User{ID: 1, Email: "user@example.com", ...}
	// Uyarı: Lazy loading yapılabilir, her zaman dolu olmayabilir
	User *user.User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// Bu interface, oturum verilerinin veritabanında yönetilmesini sağlayan
// repository katmanının sözleşmesini tanımlar.
// Tüm oturum işlemleri bu interface üzerinden yapılmalıdır.
//
// Kullanım Senaryoları:
// - Oturum oluşturma ve kaydetme
// - Oturum bilgilerini ID veya Token ile sorgulama
// - Oturum silme ve temizleme
// - Kullanıcının tüm oturumlarını silme (çıkış işlemi)
//
// Implementasyon Notları:
// - Tüm metotlar context parametresi alır (iptal ve timeout desteği için)
// - Tüm metotlar hata döndürebilir (veritabanı hataları, kayıt bulunamadı vb.)
// - Context iptal edilirse işlem durdurulmalıdır
type Repository interface {
	// Bu metod, yeni bir oturum kaydını veritabanına ekler.
	// Parametreler:
	//   - ctx: İşlem bağlamı (context), iptal ve timeout desteği sağlar
	//   - session: Kaydedilecek oturum nesnesi pointer'ı
	// Dönüş Değerleri:
	//   - error: İşlem başarılı ise nil, hata ise error nesnesi
	// Kullanım Örneği:
	//   session := &Session{
	//       Token: "abc123xyz",
	//       UserID: 1,
	//       ExpiresAt: time.Now().Add(24 * time.Hour),
	//       IPAddress: "192.168.1.1",
	//       UserAgent: "Mozilla/5.0...",
	//   }
	//   err := repo.Create(ctx, session)
	//   if err != nil {
	//       log.Printf("Oturum oluşturulamadı: %v", err)
	//   }
	// Uyarılar:
	//   - Token alanı benzersiz olmalıdır, aksi takdirde hata döner
	//   - UserID geçerli bir kullanıcıya ait olmalıdır
	//   - ExpiresAt zamanı gelecekte olmalıdır
	Create(ctx context.Context, session *Session) error

	// Bu metod, verilen ID'ye sahip oturum kaydını veritabanından sorgular.
	// Parametreler:
	//   - ctx: İşlem bağlamı (context), iptal ve timeout desteği sağlar
	//   - id: Aranacak oturumun ID'si
	// Dönüş Değerleri:
	//   - *Session: Bulunan oturum nesnesi pointer'ı (nil ise bulunamadı)
	//   - error: İşlem başarılı ise nil, hata ise error nesnesi
	// Kullanım Örneği:
	//   session, err := repo.FindByID(ctx, 42)
	//   if err != nil {
	//       log.Printf("Oturum sorgulanırken hata: %v", err)
	//       return
	//   }
	//   if session == nil {
	//       log.Println("Oturum bulunamadı")
	//       return
	//   }
	//   log.Printf("Oturum bulundu: %+v", session)
	// Uyarılar:
	//   - Oturum bulunamadı ise nil döner (hata değil)
	//   - Veritabanı hatası ise error döner
	FindByID(ctx context.Context, id uint) (*Session, error)

	// Bu metod, verilen token değerine sahip oturum kaydını veritabanından sorgular.
	// Token tabanlı kimlik doğrulama işlemlerinde kullanılır.
	// Parametreler:
	//   - ctx: İşlem bağlamı (context), iptal ve timeout desteği sağlar
	//   - token: Aranacak oturumun token değeri
	// Dönüş Değerleri:
	//   - *Session: Bulunan oturum nesnesi pointer'ı (nil ise bulunamadı)
	//   - error: İşlem başarılı ise nil, hata ise error nesnesi
	// Kullanım Örneği:
	//   token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
	//   session, err := repo.FindByToken(ctx, token)
	//   if err != nil {
	//       log.Printf("Token sorgulanırken hata: %v", err)
	//       return
	//   }
	//   if session == nil {
	//       log.Println("Geçersiz token")
	//       return
	//   }
	//   // Token geçerli, oturum bilgilerini kullan
	//   log.Printf("Kullanıcı ID: %d", session.UserID)
	// Uyarılar:
	//   - Token bulunamadı ise nil döner (hata değil)
	//   - Token benzersiz olmalıdır, birden fazla sonuç dönmemelidir
	//   - Oturum süresi dolmuş olabilir, ExpiresAt kontrol edilmelidir
	FindByToken(ctx context.Context, token string) (*Session, error)

	// Bu metod, verilen ID'ye sahip oturum kaydını veritabanından siler.
	// Parametreler:
	//   - ctx: İşlem bağlamı (context), iptal ve timeout desteği sağlar
	//   - id: Silinecek oturumun ID'si
	// Dönüş Değerleri:
	//   - error: İşlem başarılı ise nil, hata ise error nesnesi
	// Kullanım Örneği:
	//   err := repo.Delete(ctx, 42)
	//   if err != nil {
	//       log.Printf("Oturum silinirken hata: %v", err)
	//       return
	//   }
	//   log.Println("Oturum başarıyla silindi")
	// Uyarılar:
	//   - Oturum bulunamadı ise hata döner
	//   - Silme işlemi geri alınamaz
	//   - Denetim izleri için silme işlemi kaydedilmelidir
	Delete(ctx context.Context, id uint) error

	// Bu metod, verilen token değerine sahip oturum kaydını veritabanından siler.
	// Kullanıcı çıkış işleminde veya token iptal edilmesi gerektiğinde kullanılır.
	// Parametreler:
	//   - ctx: İşlem bağlamı (context), iptal ve timeout desteği sağlar
	//   - token: Silinecek oturumun token değeri
	// Dönüş Değerleri:
	//   - error: İşlem başarılı ise nil, hata ise error nesnesi
	// Kullanım Örneği:
	//   token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
	//   err := repo.DeleteByToken(ctx, token)
	//   if err != nil {
	//       log.Printf("Oturum silinirken hata: %v", err)
	//       return
	//   }
	//   log.Println("Oturum başarıyla sonlandırıldı")
	// Uyarılar:
	//   - Token bulunamadı ise hata döner
	//   - Silme işlemi geri alınamaz
	//   - Kullanıcı çıkış işleminde bu metod kullanılmalıdır
	DeleteByToken(ctx context.Context, token string) error

	// Bu metod, verilen UserID'ye ait tüm oturum kayıtlarını veritabanından siler.
	// Kullanıcının tüm cihazlardan çıkış yapması gerektiğinde kullanılır.
	// Parametreler:
	//   - ctx: İşlem bağlamı (context), iptal ve timeout desteği sağlar
	//   - userID: Oturumları silinecek kullanıcının ID'si
	// Dönüş Değerleri:
	//   - error: İşlem başarılı ise nil, hata ise error nesnesi
	// Kullanım Örneği:
	//   userID := uint(1)
	//   err := repo.DeleteByUserID(ctx, userID)
	//   if err != nil {
	//       log.Printf("Kullanıcı oturumları silinirken hata: %v", err)
	//       return
	//   }
	//   log.Printf("Kullanıcı %d'nin tüm oturumları silindi", userID)
	// Uyarılar:
	//   - Bu işlem kullanıcının tüm oturumlarını siler
	//   - Silme işlemi geri alınamaz
	//   - Güvenlik olayı (şifre değişikliği, hesap ele geçirilmesi vb.) durumunda kullanılır
	//   - Kullanıcı hesabı silindiğinde de bu metod çağrılmalıdır
	DeleteByUserID(ctx context.Context, userID uint) error
}
