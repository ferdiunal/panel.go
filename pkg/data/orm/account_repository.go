package orm

import (
	"context"

	"github.com/ferdiunal/panel.go/pkg/domain/account"
	"gorm.io/gorm"
)

// Bu yapı, hesap (account) veri tabanı işlemlerini yönetmek için kullanılan repository'dir.
// AccountRepository, GORM ORM kütüphanesi aracılığıyla hesap verilerinin oluşturulması,
// okunması, güncellenmesi ve silinmesi gibi CRUD operasyonlarını gerçekleştirir.
//
// Kullanım Senaryoları:
// - Kullanıcı hesaplarının veri tabanında yönetilmesi
// - Sosyal medya sağlayıcılarıyla entegrasyon (OAuth, SAML vb.)
// - Hesap bilgilerinin sorgulanması ve filtrelenmesi
// - Çoklu hesap desteği (bir kullanıcının birden fazla sosyal medya hesabı)
//
// Örnek Kullanım:
//   db := gorm.Open(...)
//   repo := NewAccountRepository(db)
//   account, err := repo.FindByID(ctx, 1)
//   if err != nil {
//       log.Fatal(err)
//   }
type AccountRepository struct {
	// db, GORM veri tabanı bağlantısını temsil eden pointer'dır.
	// Bu alan, tüm veri tabanı işlemlerinde kullanılır ve
	// veri tabanı sorgularının yürütülmesini sağlar.
	db *gorm.DB
}

// Bu fonksiyon, yeni bir AccountRepository örneği oluşturur ve döndürür.
// Repository pattern'ı kullanarak veri tabanı işlemlerini kapsüller.
//
// Parametreler:
//   - db (*gorm.DB): GORM veri tabanı bağlantısı, nil olmamalıdır
//
// Dönüş Değeri:
//   - Döndürür: Yapılandırılmış AccountRepository pointer'ı
//
// Kullanım Örneği:
//   db := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
//   repo := NewAccountRepository(db)
//   account, err := repo.FindByID(context.Background(), 1)
//
// Önemli Notlar:
//   - Verilen db parametresi nil olmamalıdır, aksi takdirde runtime hatası oluşur
//   - Repository, db bağlantısının yaşam döngüsünü yönetmez
//   - Aynı db bağlantısı birden fazla repository tarafından paylaşılabilir
func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

// Bu metod, yeni bir hesap kaydını veri tabanına oluşturur.
// Hesap nesnesi veri tabanında INSERT işlemi ile kaydedilir.
//
// Parametreler:
//   - ctx (context.Context): İşlem bağlamı, zaman aşımı ve iptal sinyalleri için
//   - a (*account.Account): Oluşturulacak hesap nesnesi, nil olmamalıdır
//
// Dönüş Değeri:
//   - error: İşlem başarılı ise nil, aksi takdirde hata mesajı
//
// Kullanım Örneği:
//   newAccount := &account.Account{
//       UserID: 1,
//       ProviderID: "google",
//       AccountID: "user@gmail.com",
//   }
//   err := repo.Create(ctx, newAccount)
//   if err != nil {
//       log.Printf("Hesap oluşturulamadı: %v", err)
//   }
//
// Önemli Notlar:
//   - Hesap nesnesi geçerli bir UserID içermelidir
//   - Aynı ProviderID ve AccountID kombinasyonu benzersiz olmalıdır (unique constraint)
//   - Context zaman aşımına uğrarsa işlem iptal edilir
//   - Başarılı oluşturma sonrası hesap nesnesine ID atanır
//   - Yabancı anahtar kısıtlaması: UserID geçerli bir kullanıcıya ait olmalıdır
func (r *AccountRepository) Create(ctx context.Context, a *account.Account) error {
	return r.db.WithContext(ctx).Create(a).Error
}

// Bu metod, verilen ID'ye sahip hesabı veri tabanından bulur ve döndürür.
// Hesap bulunursa ilişkili User nesnesi de otomatik olarak yüklenir.
//
// Parametreler:
//   - ctx (context.Context): İşlem bağlamı, zaman aşımı ve iptal sinyalleri için
//   - id (uint): Aranacak hesabın benzersiz tanımlayıcısı (primary key)
//
// Dönüş Değeri:
//   - *account.Account: Bulunan hesap nesnesi (ilişkili User ile birlikte)
//   - error: Hesap bulunamazsa gorm.ErrRecordNotFound, aksi takdirde diğer hata mesajları
//
// Kullanım Örneği:
//   account, err := repo.FindByID(ctx, 123)
//   if err != nil {
//       if errors.Is(err, gorm.ErrRecordNotFound) {
//           log.Println("Hesap bulunamadı")
//       } else {
//           log.Printf("Veri tabanı hatası: %v", err)
//       }
//   } else {
//       log.Printf("Hesap bulundu: %v", account)
//       log.Printf("Kullanıcı: %v", account.User)
//   }
//
// Önemli Notlar:
//   - Metod, ilişkili User nesnesini otomatik olarak yükler (Preload("User"))
//   - Hesap bulunamazsa nil ve gorm.ErrRecordNotFound hatası döndürülür
//   - Context zaman aşımına uğrarsa işlem iptal edilir
//   - Performans için sadece gerekli alanları seçmek istiyorsanız Select() kullanabilirsiniz
func (r *AccountRepository) FindByID(ctx context.Context, id uint) (*account.Account, error) {
	var a account.Account
	if err := r.db.WithContext(ctx).Preload("User").First(&a, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

// Bu metod, sağlayıcı ID'si ve hesap ID'sine göre hesabı bulur ve döndürür.
// OAuth/SAML entegrasyonlarında kullanıcı kimlik doğrulaması için kullanılır.
//
// Parametreler:
//   - ctx (context.Context): İşlem bağlamı, zaman aşımı ve iptal sinyalleri için
//   - providerID (string): Sosyal medya sağlayıcısının tanımlayıcısı
//     Örnekler: "google", "github", "facebook", "microsoft", "apple"
//   - accountID (string): Sağlayıcıdaki hesabın tanımlayıcısı
//     Örnekler: email, user ID, unique identifier
//
// Dönüş Değeri:
//   - *account.Account: Bulunan hesap nesnesi (ilişkili User ile birlikte)
//   - error: Hesap bulunamazsa gorm.ErrRecordNotFound, aksi takdirde diğer hata mesajları
//
// Kullanım Örneği:
//   // Google OAuth callback'inde
//   account, err := repo.FindByProvider(ctx, "google", "user@gmail.com")
//   if err != nil {
//       if errors.Is(err, gorm.ErrRecordNotFound) {
//           // Yeni hesap oluştur
//           newAccount := &account.Account{...}
//           repo.Create(ctx, newAccount)
//       }
//   } else {
//       log.Printf("Mevcut kullanıcı: %d", account.UserID)
//   }
//
// Önemli Notlar:
//   - Bu metod OAuth/SAML entegrasyonlarında kullanıcı kimlik doğrulaması için kullanılır
//   - Metod, ilişkili User nesnesini otomatik olarak yükler (Preload("User"))
//   - ProviderID ve AccountID kombinasyonu benzersiz olmalıdır (unique constraint)
//   - Context zaman aşımına uğrarsa işlem iptal edilir
//   - Sosyal medya sağlayıcılarıyla senkronizasyon sırasında kullanılır
func (r *AccountRepository) FindByProvider(ctx context.Context, providerID, accountID string) (*account.Account, error) {
	var a account.Account
	if err := r.db.WithContext(ctx).Preload("User").First(&a, "provider_id = ? AND account_id = ?", providerID, accountID).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

// Bu metod, verilen kullanıcı ID'sine ait tüm hesapları veri tabanından bulur ve döndürür.
// Bir kullanıcının birden fazla sosyal medya hesabı olabilir.
//
// Parametreler:
//   - ctx (context.Context): İşlem bağlamı, zaman aşımı ve iptal sinyalleri için
//   - userID (uint): Hesapları aranacak kullanıcının benzersiz tanımlayıcısı
//
// Dönüş Değeri:
//   - []account.Account: Bulunan hesapların dilimi (slice)
//   - error: İşlem başarılı ise nil, aksi takdirde hata mesajı
//
// Kullanım Örneği:
//   accounts, err := repo.FindByUserID(ctx, 42)
//   if err != nil {
//       log.Printf("Hesaplar alınamadı: %v", err)
//   } else {
//       for _, acc := range accounts {
//           log.Printf("Sağlayıcı: %s, Hesap ID: %s", acc.ProviderID, acc.AccountID)
//       }
//       if len(accounts) == 0 {
//           log.Println("Kullanıcının bağlı hesabı yok")
//       }
//   }
//
// Önemli Notlar:
//   - Bir kullanıcının birden fazla sosyal medya hesabı olabilir
//   - Hesap bulunamazsa boş bir dilim döndürülür (nil değil, len = 0)
//   - Context zaman aşımına uğrarsa işlem iptal edilir
//   - Bu metod, ilişkili User nesnesini yüklemez (performans için)
//   - Kullanıcı profili sayfasında bağlı hesapları göstermek için kullanılır
func (r *AccountRepository) FindByUserID(ctx context.Context, userID uint) ([]account.Account, error) {
	var accounts []account.Account
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&accounts).Error; err != nil {
		return nil, err
	}
	return accounts, nil
}

// Bu metod, mevcut bir hesap kaydını veri tabanında günceller.
// Hesap nesnesi veri tabanında UPDATE işlemi ile güncellenir.
//
// Parametreler:
//   - ctx (context.Context): İşlem bağlamı, zaman aşımı ve iptal sinyalleri için
//   - a (*account.Account): Güncellenecek hesap nesnesi (ID alanı zorunludur)
//
// Dönüş Değeri:
//   - error: İşlem başarılı ise nil, aksi takdirde hata mesajı
//
// Kullanım Örneği:
//   account, _ := repo.FindByID(ctx, 1)
//   account.AccountID = "newemail@gmail.com"
//   account.AccessToken = "new_token_value"
//   err := repo.Update(ctx, account)
//   if err != nil {
//       log.Printf("Hesap güncellenemedi: %v", err)
//   } else {
//       log.Println("Hesap başarıyla güncellendi")
//   }
//
// Önemli Notlar:
//   - Hesap nesnesi geçerli bir ID içermelidir
//   - GORM Save metodu, tüm alanları günceller (zero değerleri de dahil)
//   - Güncellenecek hesap veri tabanında mevcut olmalıdır
//   - Context zaman aşımına uğrarsa işlem iptal edilir
//   - Unique constraint ihlali durumunda hata döndürülür
//   - Yabancı anahtar kısıtlaması: UserID geçerli bir kullanıcıya ait olmalıdır
func (r *AccountRepository) Update(ctx context.Context, a *account.Account) error {
	return r.db.WithContext(ctx).Save(a).Error
}

// Bu metod, verilen ID'ye sahip hesap kaydını veri tabanından siler.
// Hesap kaydı veri tabanından tamamen kaldırılır (soft delete değil).
//
// Parametreler:
//   - ctx (context.Context): İşlem bağlamı, zaman aşımı ve iptal sinyalleri için
//   - id (uint): Silinecek hesabın benzersiz tanımlayıcısı (primary key)
//
// Dönüş Değeri:
//   - error: İşlem başarılı ise nil, aksi takdirde hata mesajı
//
// Kullanım Örneği:
//   err := repo.Delete(ctx, 123)
//   if err != nil {
//       log.Printf("Hesap silinemedi: %v", err)
//   } else {
//       log.Println("Hesap başarıyla silindi")
//   }
//
// Önemli Notlar:
//   - Silme işlemi geri alınamaz, dikkatli kullanılmalıdır
//   - Hesap veri tabanında mevcut olmasa bile hata döndürmez
//   - Yabancı anahtar kısıtlamaları varsa silme başarısız olabilir
//   - Context zaman aşımına uğrarsa işlem iptal edilir
//   - Soft delete kullanılmıyorsa kayıt tamamen silinir
//   - Kullanıcı hesabını bağlantısını kesmek istediğinde kullanılır
func (r *AccountRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&account.Account{}, "id = ?", id).Error
}
