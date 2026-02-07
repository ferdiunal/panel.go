// Bu paket, doğrulama (verification) varlıklarının veritabanı işlemlerini yönetmek için
// GORM ORM kütüphanesi kullanarak repository pattern'ini uygular.
// Doğrulama kayıtlarının oluşturulması, sorgulanması ve silinmesi gibi temel CRUD
// operasyonlarını sağlar.
package orm

import (
	"context"

	"github.com/ferdiunal/panel.go/pkg/domain/verification"
	"gorm.io/gorm"
)

// Bu yapı, doğrulama (verification) varlıklarının veritabanı işlemlerini gerçekleştiren
// repository'dir. GORM veritabanı bağlantısını içerir ve tüm doğrulama ile ilgili
// veritabanı operasyonlarını yönetir.
//
// Kullanım Senaryoları:
// - Kullanıcı e-posta doğrulaması için doğrulama kayıtları oluşturma
// - Doğrulama token'ları aracılığıyla doğrulama bilgilerini sorgulama
// - Tamamlanan veya süresi dolan doğrulama kayıtlarını silme
// - Belirli bir tanımlayıcıya (identifier) göre doğrulama kayıtlarını temizleme
//
// Örnek Kullanım:
//
//	repo := NewVerificationRepository(db)
//	verification := &verification.Verification{
//		Token: "abc123xyz",
//		Identifier: "user@example.com",
//	}
//	err := repo.Create(context.Background(), verification)
//	if err != nil {
//		log.Fatal(err)
//	}
type VerificationRepository struct {
	// db, GORM veritabanı bağlantısını temsil eder.
	// Tüm veritabanı operasyonları bu bağlantı üzerinden gerçekleştirilir.
	db *gorm.DB
}

// Bu fonksiyon, yeni bir VerificationRepository örneği oluşturur ve döndürür.
// Verilen GORM veritabanı bağlantısını kullanarak repository'yi başlatır.
//
// Parametreler:
// - db (*gorm.DB): GORM veritabanı bağlantısı. Bu bağlantı tüm veritabanı
//   operasyonları için kullanılacaktır.
//
// Döndürür:
// - *VerificationRepository: Başlatılmış VerificationRepository pointer'ı
//
// Önemli Notlar:
// - Verilen db parametresi nil olmamalıdır, aksi takdirde runtime hatası oluşur
// - Repository, verilen db bağlantısını doğrudan kullanır, kopyalamaz
// - Aynı db bağlantısı birden fazla repository tarafından paylaşılabilir
//
// Kullanım Örneği:
//
//	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
//	if err != nil {
//		panic("veritabanı bağlantısı başarısız")
//	}
//	repo := NewVerificationRepository(db)
func NewVerificationRepository(db *gorm.DB) *VerificationRepository {
	return &VerificationRepository{db: db}
}

// Bu metod, yeni bir doğrulama kaydını veritabanına oluşturur.
// Verilen verification nesnesini veritabanına kaydeder ve otomatik olarak
// ID ve timestamp alanlarını doldurur.
//
// Parametreler:
// - ctx (context.Context): İşlem bağlamı. Zaman aşımı ve iptal sinyallerini
//   destekler. Veritabanı operasyonunun ne kadar süre çalışabileceğini kontrol eder.
// - v (*verification.Verification): Oluşturulacak doğrulama nesnesi. Bu nesne
//   token, identifier ve diğer gerekli alanları içermelidir.
//
// Döndürür:
// - error: İşlem başarılı ise nil, aksi takdirde hata mesajı
//
// Olası Hatalar:
// - Veritabanı bağlantı hatası
// - Veri doğrulama hatası (validation error)
// - Benzersizlik kısıtlaması ihlali (unique constraint violation)
// - Bağlam zaman aşımı (context deadline exceeded)
//
// Önemli Notlar:
// - Verilen verification nesnesi nil olmamalıdır
// - Operasyon başarılı olursa, verification nesnesine veritabanı tarafından
//   atanan ID değeri otomatik olarak doldurulur
// - Context iptal edilirse, operasyon durdurulur ve hata döndürülür
//
// Kullanım Örneği:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	v := &verification.Verification{
//		Token: "secure_token_123",
//		Identifier: "user@example.com",
//		Type: "email",
//	}
//	err := repo.Create(ctx, v)
//	if err != nil {
//		log.Printf("Doğrulama oluşturma hatası: %v", err)
//	}
func (r *VerificationRepository) Create(ctx context.Context, v *verification.Verification) error {
	return r.db.WithContext(ctx).Create(v).Error
}

// Bu metod, verilen token değerine göre doğrulama kaydını veritabanından bulur
// ve döndürür. Token benzersiz bir tanımlayıcı olarak kullanılır.
//
// Parametreler:
// - ctx (context.Context): İşlem bağlamı. Zaman aşımı ve iptal sinyallerini
//   destekler. Sorgu operasyonunun ne kadar süre çalışabileceğini kontrol eder.
// - token (string): Aranacak doğrulama token'ı. Bu değer veritabanında
//   benzersiz olmalıdır.
//
// Döndürür:
// - *verification.Verification: Bulunan doğrulama nesnesi pointer'ı
// - error: İşlem başarılı ise nil, aksi takdirde hata mesajı
//
// Olası Hatalar:
// - gorm.ErrRecordNotFound: Verilen token'a sahip kayıt bulunamadı
// - Veritabanı bağlantı hatası
// - Bağlam zaman aşımı (context deadline exceeded)
//
// Önemli Notlar:
// - Token boş string olmamalıdır
// - Kayıt bulunamadı ise, gorm.ErrRecordNotFound hatası döndürülür
// - Döndürülen pointer'ı nil kontrol etmeden kullanmayın
// - Sorgu sonucu nil ise, hata değerini kontrol edin
//
// Kullanım Örneği:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	v, err := repo.FindByToken(ctx, "secure_token_123")
//	if err != nil {
//		if errors.Is(err, gorm.ErrRecordNotFound) {
//			log.Println("Doğrulama token'ı bulunamadı")
//		} else {
//			log.Printf("Veritabanı hatası: %v", err)
//		}
//		return
//	}
//	log.Printf("Doğrulama bulundu: %+v", v)
func (r *VerificationRepository) FindByToken(ctx context.Context, token string) (*verification.Verification, error) {
	var v verification.Verification
	if err := r.db.WithContext(ctx).First(&v, "token = ?", token).Error; err != nil {
		return nil, err
	}
	return &v, nil
}

// Bu metod, verilen ID değerine göre doğrulama kaydını veritabanından siler.
// Silme işlemi kalıcıdır ve geri alınamaz.
//
// Parametreler:
// - ctx (context.Context): İşlem bağlamı. Zaman aşımı ve iptal sinyallerini
//   destekler. Silme operasyonunun ne kadar süre çalışabileceğini kontrol eder.
// - id (uint): Silinecek doğrulama kaydının birincil anahtarı (ID).
//
// Döndürür:
// - error: İşlem başarılı ise nil, aksi takdirde hata mesajı
//
// Olası Hatalar:
// - Veritabanı bağlantı hatası
// - Bağlam zaman aşımı (context deadline exceeded)
// - Yabancı anahtar kısıtlaması ihlali (foreign key constraint violation)
//
// Önemli Notlar:
// - ID değeri 0 olmamalıdır (geçersiz ID)
// - Silme işlemi başarılı olsa bile, etkilenen satır sayısı kontrol edilmez
// - Kayıt bulunamasa bile, hata döndürülmez (GORM davranışı)
// - Silme işlemi geri alınamaz, bu nedenle dikkatli kullanılmalıdır
//
// Kullanım Örneği:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	err := repo.Delete(ctx, 123)
//	if err != nil {
//		log.Printf("Doğrulama silme hatası: %v", err)
//	} else {
//		log.Println("Doğrulama başarıyla silindi")
//	}
func (r *VerificationRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&verification.Verification{}, "id = ?", id).Error
}

// Bu metod, verilen identifier (tanımlayıcı) değerine göre doğrulama kaydını
// veritabanından siler. Identifier genellikle e-posta adresi veya kullanıcı adı
// gibi bir benzersiz tanımlayıcıdır. Silme işlemi kalıcıdır ve geri alınamaz.
//
// Parametreler:
// - ctx (context.Context): İşlem bağlamı. Zaman aşımı ve iptal sinyallerini
//   destekler. Silme operasyonunun ne kadar süre çalışabileceğini kontrol eder.
// - identifier (string): Silinecek doğrulama kaydının tanımlayıcısı.
//   Genellikle e-posta adresi, telefon numarası veya kullanıcı adı olabilir.
//
// Döndürür:
// - error: İşlem başarılı ise nil, aksi takdirde hata mesajı
//
// Olası Hatalar:
// - Veritabanı bağlantı hatası
// - Bağlam zaman aşımı (context deadline exceeded)
// - Yabancı anahtar kısıtlaması ihlali (foreign key constraint violation)
//
// Önemli Notlar:
// - Identifier boş string olmamalıdır
// - Aynı identifier'a sahip birden fazla kayıt varsa, hepsi silinir
// - Silme işlemi başarılı olsa bile, etkilenen satır sayısı kontrol edilmez
// - Kayıt bulunamasa bile, hata döndürülmez (GORM davranışı)
// - Silme işlemi geri alınamaz, bu nedenle dikkatli kullanılmalıdır
// - Genellikle kullanıcı kaydı silinirken veya e-posta değiştirilirken kullanılır
//
// Kullanım Örneği:
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	err := repo.DeleteByIdentifier(ctx, "user@example.com")
//	if err != nil {
//		log.Printf("Doğrulama silme hatası: %v", err)
//	} else {
//		log.Println("Kullanıcının tüm doğrulama kayıtları silindi")
//	}
//
// Uyarı:
// - Bu metod, aynı identifier'a sahip tüm doğrulama kayıtlarını siler.
//   Eğer sadece belirli bir kaydı silmek istiyorsanız, Delete() metodunu kullanın.
func (r *VerificationRepository) DeleteByIdentifier(ctx context.Context, identifier string) error {
	return r.db.WithContext(ctx).Delete(&verification.Verification{}, "identifier = ?", identifier).Error
}
