// Package orm, veritabanı işlemleri için repository pattern'ı uygulayan veri erişim katmanını içerir.
// Bu paket, GORM ORM kütüphanesi kullanarak veritabanı operasyonlarını soyutlar ve
// domain katmanı ile veritabanı arasında bir köprü görevi görür.
package orm

import (
	"context"

	"github.com/ferdiunal/panel.go/pkg/domain/session"
	"gorm.io/gorm"
)

// Bu yapı, oturum (session) verilerinin veritabanında yönetilmesini sağlayan repository'dir.
// SessionRepository, GORM veritabanı bağlantısını kullanarak session CRUD operasyonlarını gerçekleştirir.
//
// Kullanım Senaryoları:
// - Kullanıcı oturum bilgilerini oluşturma ve depolama
// - Oturum token'ı ile kullanıcı doğrulama
// - Oturum ID'si ile oturum bilgilerini sorgulama
// - Oturum sonlandırma ve silme işlemleri
// - Kullanıcı çıkış yaptığında tüm oturumlarını temizleme
//
// Önemli Notlar:
// - Tüm metotlar context parametresi alır, bu sayede işlemleri iptal edebilirsiniz
// - Preload("User") kullanılarak ilişkili User verisi otomatik olarak yüklenir
// - Veritabanı hataları error olarak döndürülür, nil ise işlem başarılıdır
type SessionRepository struct {
	// db, GORM veritabanı bağlantı nesnesidir.
	// Bu bağlantı tüm veritabanı operasyonları için kullanılır.
	db *gorm.DB
}

// Bu fonksiyon, SessionRepository'nin yeni bir örneğini oluşturur ve başlatır.
// Dependency injection pattern'ı kullanarak GORM bağlantısını alır.
//
// Parametreler:
// - db (*gorm.DB): GORM veritabanı bağlantı nesnesi
//
// Dönüş Değeri:
// - *SessionRepository: Yapılandırılmış SessionRepository pointer'ı
//
// Kullanım Örneği:
//   db := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
//   repo := NewSessionRepository(db)
//
// Önemli Notlar:
// - db parametresi nil olmamalıdır, aksi takdirde runtime hatası oluşur
// - Genellikle uygulama başlangıcında bir kez çağrılır
// - Döndürülen repository tüm oturum işlemleri için kullanılabilir
func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

// Bu metod, yeni bir oturum kaydını veritabanına oluşturur.
// Verilen session nesnesini veritabanına kaydeder ve otomatik ID ataması yapılır.
//
// Parametreler:
// - ctx (context.Context): İşlem bağlamı, timeout ve iptal sinyalleri için kullanılır
// - s (*session.Session): Oluşturulacak oturum nesnesi
//
// Dönüş Değeri:
// - error: İşlem başarılı ise nil, aksi takdirde hata mesajı
//
// Kullanım Örneği:
//   newSession := &session.Session{
//       UserID: 1,
//       Token: "abc123xyz",
//       ExpiresAt: time.Now().Add(24 * time.Hour),
//   }
//   err := repo.Create(ctx, newSession)
//   if err != nil {
//       log.Printf("Oturum oluşturulamadı: %v", err)
//   }
//
// Önemli Notlar:
// - Session nesnesi nil olmamalıdır
// - Token alanı benzersiz (unique) olmalıdır
// - Context timeout'u aşılırsa işlem iptal edilir
// - Başarılı oluşturmada session nesnesine otomatik ID atanır
func (r *SessionRepository) Create(ctx context.Context, s *session.Session) error {
	return r.db.WithContext(ctx).Create(s).Error
}

// Bu metod, verilen ID'ye sahip oturum kaydını veritabanından bulur.
// Oturum bulunursa ilişkili User bilgisi de otomatik olarak yüklenir (Preload).
//
// Parametreler:
// - ctx (context.Context): İşlem bağlamı, timeout ve iptal sinyalleri için kullanılır
// - id (uint): Aranacak oturum ID'si
//
// Dönüş Değeri:
// - *session.Session: Bulunan oturum nesnesi (nil ise bulunamadı)
// - error: İşlem başarılı ise nil, aksi takdirde hata mesajı
//
// Kullanım Örneği:
//   session, err := repo.FindByID(ctx, 42)
//   if err != nil {
//       if errors.Is(err, gorm.ErrRecordNotFound) {
//           log.Println("Oturum bulunamadı")
//       }
//       return
//   }
//   log.Printf("Oturum bulundu: %+v", session)
//
// Önemli Notlar:
// - ID sıfır veya negatif olmamalıdır
// - Oturum bulunamadığında gorm.ErrRecordNotFound hatası döndürülür
// - User ilişkisi otomatik olarak yüklenir, bu sayede session.User erişilebilir
// - Context timeout'u aşılırsa işlem iptal edilir
func (r *SessionRepository) FindByID(ctx context.Context, id uint) (*session.Session, error) {
	var s session.Session
	if err := r.db.WithContext(ctx).Preload("User").First(&s, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

// Bu metod, verilen token'a sahip oturum kaydını veritabanından bulur.
// Token tabanlı oturum doğrulaması için kullanılır. İlişkili User bilgisi otomatik yüklenir.
//
// Parametreler:
// - ctx (context.Context): İşlem bağlamı, timeout ve iptal sinyalleri için kullanılır
// - token (string): Aranacak oturum token'ı
//
// Dönüş Değeri:
// - *session.Session: Bulunan oturum nesnesi (nil ise bulunamadı)
// - error: İşlem başarılı ise nil, aksi takdirde hata mesajı
//
// Kullanım Örneği:
//   // HTTP isteğinden token alınır
//   token := r.Header.Get("Authorization")
//   session, err := repo.FindByToken(ctx, token)
//   if err != nil {
//       if errors.Is(err, gorm.ErrRecordNotFound) {
//           http.Error(w, "Geçersiz token", http.StatusUnauthorized)
//       }
//       return
//   }
//   // Oturum geçerli, kullanıcı doğrulandı
//   log.Printf("Kullanıcı doğrulandı: %d", session.UserID)
//
// Önemli Notlar:
// - Token boş string olmamalıdır
// - Token benzersiz (unique) olmalıdır
// - Oturum bulunamadığında gorm.ErrRecordNotFound hatası döndürülür
// - User ilişkisi otomatik olarak yüklenir
// - API kimlik doğrulaması için sıkça kullanılır
// - Context timeout'u aşılırsa işlem iptal edilir
func (r *SessionRepository) FindByToken(ctx context.Context, token string) (*session.Session, error) {
	var s session.Session
	if err := r.db.WithContext(ctx).Preload("User").First(&s, "token = ?", token).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

// Bu metod, verilen ID'ye sahip oturum kaydını veritabanından siler.
// Oturum sonlandırma ve temizleme işlemleri için kullanılır.
//
// Parametreler:
// - ctx (context.Context): İşlem bağlamı, timeout ve iptal sinyalleri için kullanılır
// - id (uint): Silinecek oturum ID'si
//
// Dönüş Değeri:
// - error: İşlem başarılı ise nil, aksi takdirde hata mesajı
//
// Kullanım Örneği:
//   err := repo.Delete(ctx, 42)
//   if err != nil {
//       log.Printf("Oturum silinemedi: %v", err)
//       return
//   }
//   log.Println("Oturum başarıyla silindi")
//
// Önemli Notlar:
// - ID sıfır veya negatif olmamalıdır
// - Oturum bulunamasa bile hata döndürülmez (soft delete değildir)
// - Silme işlemi geri alınamaz
// - Context timeout'u aşılırsa işlem iptal edilir
// - Veritabanı kısıtlamaları (foreign key) varsa hata döndürülebilir
func (r *SessionRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&session.Session{}, "id = ?", id).Error
}

// Bu metod, verilen token'a sahip oturum kaydını veritabanından siler.
// Token tabanlı oturum sonlandırma için kullanılır.
//
// Parametreler:
// - ctx (context.Context): İşlem bağlamı, timeout ve iptal sinyalleri için kullanılır
// - token (string): Silinecek oturum token'ı
//
// Dönüş Değeri:
// - error: İşlem başarılı ise nil, aksi takdirde hata mesajı
//
// Kullanım Örneği:
//   // Kullanıcı çıkış yaptığında
//   token := r.Header.Get("Authorization")
//   err := repo.DeleteByToken(ctx, token)
//   if err != nil {
//       log.Printf("Oturum silinemedi: %v", err)
//       return
//   }
//   log.Println("Oturum başarıyla sonlandırıldı")
//
// Önemli Notlar:
// - Token boş string olmamalıdır
// - Oturum bulunamsa bile hata döndürülmez
// - Silme işlemi geri alınamaz
// - Kullanıcı çıkış (logout) işleminde sıkça kullanılır
// - Context timeout'u aşılırsa işlem iptal edilir
func (r *SessionRepository) DeleteByToken(ctx context.Context, token string) error {
	return r.db.WithContext(ctx).Delete(&session.Session{}, "token = ?", token).Error
}

// Bu metod, verilen kullanıcı ID'sine ait tüm oturum kayıtlarını veritabanından siler.
// Kullanıcı çıkış yaptığında tüm oturumlarını temizlemek için kullanılır.
//
// Parametreler:
// - ctx (context.Context): İşlem bağlamı, timeout ve iptal sinyalleri için kullanılır
// - userID (uint): Oturumları silinecek kullanıcı ID'si
//
// Dönüş Değeri:
// - error: İşlem başarılı ise nil, aksi takdirde hata mesajı
//
// Kullanım Örneği:
//   // Kullanıcı hesabı silindiğinde tüm oturumlarını temizle
//   err := repo.DeleteByUserID(ctx, userID)
//   if err != nil {
//       log.Printf("Kullanıcı oturumları silinemedi: %v", err)
//       return
//   }
//   log.Printf("Kullanıcı %d'nin tüm oturumları silindi", userID)
//
// Önemli Notlar:
// - userID sıfır veya negatif olmamalıdır
// - Bir kullanıcının birden fazla oturumu olabilir (farklı cihazlardan)
// - Bu metod tüm oturumları siler, seçici silme yapmaz
// - Silme işlemi geri alınamaz
// - Güvenlik nedenleriyle hesap kapatma veya şifre değişikliğinde kullanılır
// - Context timeout'u aşılırsa işlem iptal edilir
// - Silinen oturum sayısı RowsAffected ile kontrol edilebilir
func (r *SessionRepository) DeleteByUserID(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).Delete(&session.Session{}, "user_id = ?", userID).Error
}
