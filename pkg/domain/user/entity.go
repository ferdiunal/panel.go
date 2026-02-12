// Bu paket, kullanıcı domain modelini ve repository arayüzünü tanımlar.
// Kullanıcı yönetimi, kimlik doğrulama ve yetkilendirme işlemleri için temel yapıları içerir.
package user

import (
	"context"
	"time"
)

// Bu yapı, sistem içindeki bir kullanıcıyı temsil eder.
//
// User yapısı, kullanıcı bilgilerini depolamak ve yönetmek için kullanılır.
// Veritabanında users tablosuna karşılık gelir ve GORM ORM tarafından yönetilir.
//
// Alanlar:
//   - ID: Kullanıcının benzersiz tanımlayıcısı (birincil anahtar)
//   - Name: Kullanıcının tam adı (indekslenmiş, arama performansı için)
//   - Email: Kullanıcının e-posta adresi (benzersiz indeks, çift kayıt önleme)
//   - EmailVerified: E-posta adresinin doğrulanıp doğrulanmadığını belirten bayrak
//   - Image: Kullanıcının profil resmi URL'si
//   - Role: Kullanıcının sistem içindeki rolü (admin, user, moderator vb.)
//   - CreatedAt: Kullanıcı hesabının oluşturulma tarihi ve saati
//   - UpdatedAt: Kullanıcı bilgilerinin son güncellenme tarihi ve saati
//
// Kullanım Senaryoları:
//   - Kullanıcı kaydı ve profil yönetimi
//   - Kimlik doğrulama ve oturum yönetimi
//   - Rol tabanlı erişim kontrolü (RBAC)
//   - Kullanıcı arama ve filtreleme işlemleri
//
// Önemli Notlar:
//   - Email alanı benzersiz olmalıdır, aynı e-posta ile iki kullanıcı kaydı yapılamaz
//   - CreatedAt ve UpdatedAt alanları GORM tarafından otomatik olarak yönetilir
//   - Role alanı yetkilendirme kararları için kritik öneme sahiptir
//   - EmailVerified bayrağı, e-posta doğrulama akışında kullanılır
//
// Örnek Kullanım:
//
//	user := &User{
//	    Name:          "Ahmet Yılmaz",
//	    Email:         "ahmet@example.com",
//	    EmailVerified: true,
//	    Role:          "user",
//	}
//	err := repo.CreateUser(ctx, user)
type User struct {
	// ID: Kullanıcının benzersiz tanımlayıcısı
	// Veritabanında birincil anahtar olarak kullanılır
	// GORM tarafından otomatik olarak artırılır
	ID uint `json:"id" gorm:"primaryKey"`

	// Name: Kullanıcının tam adı
	// İndekslenmiş alan, ad bazında arama performansını iyileştirir
	// Boş bırakılamaz, en az 1 karakter içermelidir
	Name string `json:"name" gorm:"index"`

	// Email: Kullanıcının e-posta adresi
	// Benzersiz indeks ile korunur, sistem içinde her e-posta benzersiz olmalıdır
	// Kimlik doğrulama ve iletişim için kullanılır
	// RFC 5322 standardına uygun olmalıdır
	Email string `json:"email" gorm:"uniqueIndex"`

	// EmailVerified: E-posta doğrulama durumu
	// true: E-posta adresi doğrulanmış ve onaylanmış
	// false: E-posta adresi henüz doğrulanmamış
	// Hassas işlemler için doğrulama gerekli olabilir
	EmailVerified bool `json:"emailVerified"`

	// Image: Kullanıcının profil resmi URL'si
	// Kullanıcı arayüzünde profil fotoğrafı göstermek için kullanılır
	// Boş bırakılabilir, varsayılan profil resmi kullanılır
	// Tam URL formatında olmalıdır (https://...)
	Image string `json:"image"`

	// Role: Kullanıcının sistem içindeki rolü
	// İndekslenmiş alan, rol bazında filtreleme performansını iyileştirir
	// Olası değerler: "admin", "user", "moderator", "guest" vb.
	// Yetkilendirme kararları bu alan temel alınarak verilir
	// Rol tabanlı erişim kontrolü (RBAC) için kritik
	Role string `json:"role" gorm:"index"`

	// CreatedAt: Kullanıcı hesabının oluşturulma tarihi ve saati
	// İndekslenmiş alan, tarih bazında sorgular için performans sağlar
	// GORM tarafından otomatik olarak ayarlanır
	// Hesap yaşını ve kullanıcı istatistiklerini belirlemek için kullanılır
	CreatedAt time.Time `json:"createdAt" gorm:"index"`

	// UpdatedAt: Kullanıcı bilgilerinin son güncellenme tarihi ve saati
	// İndekslenmiş alan, son değişiklikleri takip etmek için kullanılır
	// GORM tarafından otomatik olarak güncellenir
	// Denetim izleri ve değişiklik geçmişi için önemlidir
	UpdatedAt time.Time `json:"updatedAt" gorm:"index"`
}

// GetID, kullanıcının ID'sini döndürür.
// Bu metod, notification handler gibi interface{ GetID() uint } bekleyen yerler için gereklidir.
func (u *User) GetID() uint {
	return u.ID
}

// Bu interface, kullanıcı verilerine erişim ve yönetim işlemlerini tanımlar.
//
// Repository interface, veri erişim katmanı (DAL) için sözleşme sağlar.
// Veritabanı işlemlerini soyutlar ve test edilebilir kod yazılmasını sağlar.
// Dependency Injection deseni ile kullanılır.
//
// Uygulamalar:
//   - GORM tabanlı SQL veritabanı uygulaması
//   - Mock repository (test ve geliştirme için)
//   - Önbellek katmanı ile dekoratör uygulaması
//
// Kullanım Senaryoları:
//   - Kullanıcı CRUD işlemleri
//   - Kimlik doğrulama ve yetkilendirme
//   - Kullanıcı arama ve filtreleme
//   - Toplu işlemler ve raporlama
//
// Önemli Notlar:
//   - Tüm metotlar context parametresi alır, iptal ve zaman aşımı desteği için
//   - Hata yönetimi çağıran tarafından yapılmalıdır
//   - Tüm işlemler atomik olmalıdır (transaction desteği)
//   - Eşzamanlı erişim için thread-safe olmalıdır
//
// Örnek Uygulamalar:
//
//	type UserRepository struct {
//	    db *gorm.DB
//	}
//
//	func (r *UserRepository) CreateUser(ctx context.Context, user *User) error {
//	    return r.db.WithContext(ctx).Create(user).Error
//	}
type Repository interface {
	// Bu metod, yeni bir kullanıcı kaydı oluşturur.
	//
	// CreateUser fonksiyonu, verilen User yapısını veritabanına ekler.
	// Başarılı olursa, kullanıcıya veritabanı tarafından atanan ID değeri atanır.
	//
	// Parametreler:
	//   - ctx: İşlem bağlamı, iptal ve zaman aşımı kontrolü için
	//   - user: Oluşturulacak kullanıcı verisi (pointer)
	//
	// Dönüş Değerleri:
	//   - error: Hata durumunda hata nesnesi, başarılı ise nil
	//
	// Olası Hatalar:
	//   - Benzersiz kısıtlama ihlali (aynı e-posta zaten var)
	//   - Veritabanı bağlantı hatası
	//   - Geçersiz veri (null değerler, format hataları)
	//   - Context iptal edildi
	//
	// Kullanım Örneği:
	//   user := &User{
	//       Name:  "Fatih Kaplan",
	//       Email: "fatih@example.com",
	//       Role:  "user",
	//   }
	//   err := repo.CreateUser(ctx, user)
	//   if err != nil {
	//       log.Printf("Kullanıcı oluşturulamadı: %v", err)
	//   }
	//
	// Önemli Notlar:
	//   - Email alanı benzersiz olmalıdır
	//   - Başarılı oluşturmadan sonra user.ID otomatik olarak atanır
	//   - CreatedAt ve UpdatedAt alanları otomatik olarak ayarlanır
	CreateUser(ctx context.Context, user *User) error

	// Bu metod, ID'ye göre bir kullanıcıyı bulur.
	//
	// FindByID fonksiyonu, verilen ID'ye sahip kullanıcıyı veritabanından alır.
	// Kullanıcı bulunursa, User yapısı döndürülür.
	// Kullanıcı bulunamazsa, nil ve hata döndürülür.
	//
	// Parametreler:
	//   - ctx: İşlem bağlamı, iptal ve zaman aşımı kontrolü için
	//   - id: Aranacak kullanıcının benzersiz tanımlayıcısı
	//
	// Dönüş Değerleri:
	//   - *User: Bulunan kullanıcı verisi (pointer)
	//   - error: Hata durumunda hata nesnesi, başarılı ise nil
	//
	// Olası Hatalar:
	//   - Kullanıcı bulunamadı (record not found)
	//   - Veritabanı bağlantı hatası
	//   - Context iptal edildi
	//
	// Kullanım Örneği:
	//   user, err := repo.FindByID(ctx, 42)
	//   if err != nil {
	//       if errors.Is(err, gorm.ErrRecordNotFound) {
	//           log.Println("Kullanıcı bulunamadı")
	//       }
	//       return
	//   }
	//   log.Printf("Kullanıcı: %s (%s)", user.Name, user.Email)
	//
	// Önemli Notlar:
	//   - ID değeri pozitif bir sayı olmalıdır
	//   - Sık kullanılan işlem, performans kritik olabilir
	//   - Önbellek katmanı ile optimize edilebilir
	FindByID(ctx context.Context, id uint) (*User, error)

	// Bu metod, e-posta adresine göre bir kullanıcıyı bulur.
	//
	// FindByEmail fonksiyonu, verilen e-posta adresine sahip kullanıcıyı bulur.
	// E-posta adresi benzersiz olduğu için, en fazla bir kullanıcı döndürülür.
	// Kimlik doğrulama işlemlerinde sıkça kullanılır.
	//
	// Parametreler:
	//   - ctx: İşlem bağlamı, iptal ve zaman aşımı kontrolü için
	//   - email: Aranacak e-posta adresi
	//
	// Dönüş Değerleri:
	//   - *User: Bulunan kullanıcı verisi (pointer)
	//   - error: Hata durumunda hata nesnesi, başarılı ise nil
	//
	// Olası Hatalar:
	//   - Kullanıcı bulunamadı (record not found)
	//   - Veritabanı bağlantı hatası
	//   - Context iptal edildi
	//
	// Kullanım Örneği:
	//   user, err := repo.FindByEmail(ctx, "kullanici@example.com")
	//   if err != nil {
	//       if errors.Is(err, gorm.ErrRecordNotFound) {
	//           log.Println("Bu e-posta ile kayıtlı kullanıcı yok")
	//       }
	//       return
	//   }
	//   log.Printf("Kullanıcı bulundu: %s (ID: %d)", user.Name, user.ID)
	//
	// Önemli Notlar:
	//   - E-posta adresi benzersiz indeks ile korunur
	//   - Kimlik doğrulama akışında kritik işlem
	//   - E-posta adresi büyük/küçük harf duyarlı olabilir (veritabanına bağlı)
	//   - Giriş doğrulaması yapılmalıdır (SQL injection önleme)
	FindByEmail(ctx context.Context, email string) (*User, error)

	// Bu metod, mevcut bir kullanıcı kaydını günceller.
	//
	// UpdateUser fonksiyonu, verilen User yapısının ID'sine göre kaydı günceller.
	// Güncelleme işleminden sonra UpdatedAt alanı otomatik olarak ayarlanır.
	// Kısmi güncellemeler için GORM'un Select/Omit metotları kullanılabilir.
	//
	// Parametreler:
	//   - ctx: İşlem bağlamı, iptal ve zaman aşımı kontrolü için
	//   - user: Güncellenecek kullanıcı verisi (pointer, ID alanı zorunlu)
	//
	// Dönüş Değerleri:
	//   - error: Hata durumunda hata nesnesi, başarılı ise nil
	//
	// Olası Hatalar:
	//   - Kullanıcı bulunamadı (record not found)
	//   - Benzersiz kısıtlama ihlali (e-posta değiştirilirken)
	//   - Veritabanı bağlantı hatası
	//   - Context iptal edildi
	//   - ID alanı belirtilmedi
	//
	// Kullanım Örneği:
	//   user := &User{
	//       ID:   42,
	//       Name: "Yeni Ad",
	//       Role: "admin",
	//   }
	//   err := repo.UpdateUser(ctx, user)
	//   if err != nil {
	//       log.Printf("Kullanıcı güncellenemedi: %v", err)
	//   }
	//
	// Önemli Notlar:
	//   - ID alanı mutlaka belirtilmelidir
	//   - Tüm alanlar güncellenir, kısmi güncelleme için Select/Omit kullanın
	//   - UpdatedAt alanı otomatik olarak güncellenir
	//   - Şifre gibi hassas alanlar ayrı metotla güncellenmelidir
	UpdateUser(ctx context.Context, user *User) error

	// Bu metod, bir kullanıcı kaydını siler.
	//
	// DeleteUser fonksiyonu, verilen ID'ye sahip kullanıcıyı veritabanından siler.
	// Silme işlemi kalıcıdır ve geri alınamaz.
	// Soft delete gerekiyorsa, ayrı bir metot uygulanmalıdır.
	//
	// Parametreler:
	//   - ctx: İşlem bağlamı, iptal ve zaman aşımı kontrolü için
	//   - id: Silinecek kullanıcının benzersiz tanımlayıcısı
	//
	// Dönüş Değerleri:
	//   - error: Hata durumunda hata nesnesi, başarılı ise nil
	//
	// Olası Hatalar:
	//   - Kullanıcı bulunamadı (record not found)
	//   - Yabancı anahtar kısıtlaması ihlali (ilişkili veriler var)
	//   - Veritabanı bağlantı hatası
	//   - Context iptal edildi
	//
	// Kullanım Örneği:
	//   err := repo.DeleteUser(ctx, 42)
	//   if err != nil {
	//       if errors.Is(err, gorm.ErrRecordNotFound) {
	//           log.Println("Silinecek kullanıcı bulunamadı")
	//       }
	//       return
	//   }
	//   log.Println("Kullanıcı başarıyla silindi")
	//
	// Önemli Notlar:
	//   - Silme işlemi geri alınamaz, dikkatli kullanın
	//   - Yabancı anahtar kısıtlamaları kontrol edilmelidir
	//   - Denetim izleri için silme işlemi kaydedilmelidir
	//   - Soft delete tercih edilebilir (deleted_at alanı)
	DeleteUser(ctx context.Context, id uint) error

	// Bu metod, toplam kullanıcı sayısını döndürür.
	//
	// Count fonksiyonu, veritabanında kayıtlı toplam kullanıcı sayısını hesaplar.
	// Sayfalama, istatistik ve raporlama işlemlerinde kullanılır.
	// Büyük veri setlerinde performans kritik olabilir.
	//
	// Parametreler:
	//   - ctx: İşlem bağlamı, iptal ve zaman aşımı kontrolü için
	//
	// Dönüş Değerleri:
	//   - int64: Toplam kullanıcı sayısı
	//   - error: Hata durumunda hata nesnesi, başarılı ise nil
	//
	// Olası Hatalar:
	//   - Veritabanı bağlantı hatası
	//   - Context iptal edildi
	//   - Sorgu zaman aşımı (çok sayıda kayıt)
	//
	// Kullanım Örneği:
	//   count, err := repo.Count(ctx)
	//   if err != nil {
	//       log.Printf("Kullanıcı sayısı alınamadı: %v", err)
	//       return
	//   }
	//   log.Printf("Toplam kullanıcı sayısı: %d", count)
	//
	// Önemli Notlar:
	//   - Büyük veri setlerinde COUNT sorgusu yavaş olabilir
	//   - Önbellek katmanı ile optimize edilebilir
	//   - Sayfalama işlemlerinde limit/offset ile birlikte kullanılır
	//   - Raporlama ve istatistik amaçlı kullanılır
	Count(ctx context.Context) (int64, error)
}
