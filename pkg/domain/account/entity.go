// Bu paket, kullanıcı hesaplarının (account) domain modelini ve repository arayüzünü tanımlar.
// Hesaplar, OAuth2 sağlayıcıları (Google, GitHub vb.) veya kimlik bilgileri (email/şifre) aracılığıyla
// kullanıcı kimlik doğrulamasını yönetir.
package account

import (
	"context"
	"time"

	"github.com/ferdiunal/panel.go/pkg/domain/user"
)

// Bu yapı, bir kullanıcının kimlik doğrulama hesabını temsil eder.
// Bir kullanıcı birden fazla hesaba sahip olabilir (örneğin: Google, GitHub, email/şifre).
//
// # Kullanım Senaryoları
// - OAuth2 sağlayıcılarıyla entegrasyon (Google, GitHub, Microsoft vb.)
// - Email/şifre tabanlı kimlik doğrulama
// - Birden fazla kimlik doğrulama yöntemi desteği
// - Token yönetimi ve yenileme
//
// # Önemli Notlar
// - Password alanı asla JSON yanıtında döndürülmez (json:"-" etiketi)
// - AccountID, kimlik bilgileri sağlayıcıları için null olabilir
// - AccessToken ve RefreshToken hassas verilerdir ve güvenli şekilde saklanmalıdır
// - Tüm token alanları veritabanında şifreli olarak saklanmalıdır
//
// # Örnek Kullanım
// ```go
//
//	account := &Account{
//	    UserID:     1,
//	    ProviderID: "google",
//	    AccountID:  &googleUserID,
//	    AccessToken: "encrypted_token_here",
//	    Scope: "email profile",
//	}
//
// ```
type Account struct {
	// ID: Hesabın benzersiz tanımlayıcısı (birincil anahtar)
	// Veritabanında otomatik olarak artan bir değerdir.
	ID uint `json:"id" gorm:"primaryKey"`

	// AccountID: Sağlayıcının (provider) tarafından verilen kullanıcı kimliği
	// Örneğin: Google User ID, GitHub User ID vb.
	// Kimlik bilgileri sağlayıcıları için null olabilir (nullable).
	// Veritabanında indekslenmiştir (hızlı arama için).
	//
	// Önemli: Bu alan sağlayıcı tarafından benzersiz olmalıdır.
	AccountID *string `json:"accountId" gorm:"index"`

	// ProviderID: Kimlik doğrulama sağlayıcısının tanımlayıcısı
	// Örnek değerler: "credential" (email/şifre), "google", "github", "microsoft"
	// Veritabanında indekslenmiştir (hızlı arama için).
	//
	// Kullanım: Hangi sağlayıcıyla kimlik doğrulandığını belirlemek için kullanılır.
	ProviderID string `json:"providerId" gorm:"index"`

	// UserID: Bu hesabın ait olduğu kullanıcının kimliği
	// Veritabanında indekslenmiştir (bir kullanıcının tüm hesaplarını bulmak için).
	// Foreign key olarak User tablosuna referans verir.
	UserID uint `json:"userId" gorm:"index"`

	// AccessToken: OAuth2 sağlayıcısından alınan erişim tokeni
	// Bu token, sağlayıcının API'sine erişim için kullanılır.
	// Hassas veri: Veritabanında şifreli olarak saklanmalıdır.
	// JSON yanıtında omitempty etiketi ile isteğe bağlı olarak döndürülür.
	//
	// Uyarı: Asla log dosyalarına veya hata mesajlarına yazılmamalıdır.
	AccessToken string `json:"accessToken,omitempty"`

	// RefreshToken: OAuth2 sağlayıcısından alınan yenileme tokeni
	// AccessToken süresi dolduğunda, bu token kullanılarak yeni bir AccessToken alınır.
	// Hassas veri: Veritabanında şifreli olarak saklanmalıdır.
	// JSON yanıtında omitempty etiketi ile isteğe bağlı olarak döndürülür.
	//
	// Uyarı: Asla istemciye gönderilmemelidir, sadece sunucu tarafında kullanılmalıdır.
	RefreshToken string `json:"refreshToken,omitempty"`

	// IDToken: OpenID Connect sağlayıcılarından alınan ID tokeni
	// Kullanıcı bilgilerini içeren JWT formatında bir tokendır.
	// Hassas veri: Veritabanında şifreli olarak saklanmalıdır.
	// JSON yanıtında omitempty etiketi ile isteğe bağlı olarak döndürülür.
	//
	// Kullanım: Kullanıcı profil bilgilerini doğrulamak için kullanılır.
	IDToken string `json:"idToken,omitempty"`

	// AccessTokenExpiresAt: AccessToken'ın süresi dolma zamanı
	// Null ise, token süresi dolmaz (nadiren kullanılır).
	// Bu zaman geçtiyse, RefreshToken kullanılarak yeni bir token alınmalıdır.
	//
	// Örnek: 2024-12-31T23:59:59Z
	AccessTokenExpiresAt *time.Time `json:"accessTokenExpiresAt,omitempty"`

	// RefreshTokenExpiresAt: RefreshToken'ın süresi dolma zamanı
	// Null ise, refresh token süresi dolmaz.
	// Bu zaman geçtiyse, kullanıcı yeniden kimlik doğrulaması yapmalıdır.
	//
	// Örnek: 2025-12-31T23:59:59Z
	RefreshTokenExpiresAt *time.Time `json:"refreshTokenExpiresAt,omitempty"`

	// Password: Şifrelenmiş kullanıcı şifresi (sadece email/şifre sağlayıcısı için)
	// Bcrypt veya benzer güvenli hash algoritması kullanılarak şifrelenmiş olmalıdır.
	// JSON yanıtında asla döndürülmez (json:"-" etiketi).
	// Veritabanında saklanır ancak istemciye gönderilmez.
	//
	// Uyarı: Asla düz metin olarak saklanmamalıdır.
	// Uyarı: Asla log dosyalarına yazılmamalıdır.
	Password string `json:"-"`

	// Scope: OAuth2 sağlayıcısından istenen izinler (scope)
	// Boşlukla ayrılmış izin listesi.
	// Örnek: "email profile openid"
	//
	// Kullanım: Hangi izinlerle kimlik doğrulandığını belirlemek için kullanılır.
	Scope string `json:"scope,omitempty"`

	// CreatedAt: Hesabın oluşturulma zamanı
	// Veritabanında otomatik olarak ayarlanır.
	// Veritabanında indekslenmiştir (zaman aralığı sorgularında hızlı arama için).
	//
	// Örnek: 2024-01-15T10:30:00Z
	CreatedAt time.Time `json:"createdAt" gorm:"index"`

	// UpdatedAt: Hesabın son güncellenme zamanı
	// Veritabanında otomatik olarak güncellenir.
	// Veritabanında indekslenmiştir (zaman aralığı sorgularında hızlı arama için).
	//
	// Örnek: 2024-01-20T15:45:30Z
	UpdatedAt time.Time `json:"updatedAt" gorm:"index"`

	// User: Bu hesabın ait olduğu kullanıcı nesnesi
	// GORM tarafından otomatik olarak yüklenir (eager loading).
	// Foreign key: UserID alanı aracılığıyla User tablosuna referans verir.
	// JSON yanıtında omitempty etiketi ile isteğe bağlı olarak döndürülür.
	//
	// Kullanım: Hesapla birlikte kullanıcı bilgilerini almak için.
	// Örnek: account.User.Email
	User *user.User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// GetID, hesabın ID'sini döndürür.
// Bu metod, notification handler gibi interface{ GetID() uint } bekleyen yerler için gereklidir.
func (a *Account) GetID() uint {
	return a.ID
}

// Bu interface, hesap (account) veri erişim katmanı (repository) için sözleşmeyi tanımlar.
// Repository pattern kullanılarak, veritabanı işlemleri soyutlanır ve test edilebilir hale getirilir.
//
// # Kullanım Senaryoları
// - Hesap oluşturma, okuma, güncelleme, silme işlemleri
// - OAuth2 sağlayıcılarıyla entegrasyon
// - Kullanıcı kimlik doğrulaması
// - Hesap yönetimi
//
// # Implementasyon Notları
// - Tüm metodlar context.Context parametresi alır (iptal ve timeout desteği için)
// - Tüm metodlar hata döndürür (error handling)
// - Veritabanı işlemleri asenkron olarak yapılabilir
//
// # Örnek Implementasyon
// ```go
//
//	type AccountRepository struct {
//	    db *gorm.DB
//	}
//
//	func (r *AccountRepository) Create(ctx context.Context, account *Account) error {
//	    return r.db.WithContext(ctx).Create(account).Error
//	}
//
// ```
type Repository interface {
	// Bu metod, yeni bir hesap oluşturur ve veritabanına kaydeder.
	//
	// # Parametreler
	// - ctx: İşlem bağlamı (context), iptal ve timeout desteği için
	// - account: Oluşturulacak hesap nesnesi
	//
	// # Dönüş Değeri
	// - error: İşlem başarılı ise nil, aksi takdirde hata
	//
	// # Kullanım Senaryoları
	// - Yeni OAuth2 hesabı oluşturma
	// - Email/şifre hesabı oluşturma
	// - Sosyal medya hesabı bağlama
	//
	// # Örnek Kullanım
	// ```go
	// account := &Account{
	//     UserID: 1,
	//     ProviderID: "google",
	//     AccountID: &googleID,
	//     AccessToken: "token",
	// }
	// err := repo.Create(ctx, account)
	// if err != nil {
	//     log.Printf("Hesap oluşturma hatası: %v", err)
	// }
	// ```
	//
	// # Önemli Notlar
	// - Hesap zaten varsa, hata döndürmelidir
	// - UserID geçerli bir kullanıcıya referans etmelidir
	// - ProviderID ve AccountID kombinasyonu benzersiz olmalıdır
	Create(ctx context.Context, account *Account) error

	// Bu metod, verilen ID'ye sahip hesabı bulur ve döndürür.
	//
	// # Parametreler
	// - ctx: İşlem bağlamı (context), iptal ve timeout desteği için
	// - id: Aranacak hesabın ID'si
	//
	// # Dönüş Değeri
	// - *Account: Bulunan hesap nesnesi (nil ise bulunamadı)
	// - error: İşlem başarılı ise nil, aksi takdirde hata
	//
	// # Kullanım Senaryoları
	// - Hesap detaylarını görüntüleme
	// - Hesap bilgilerini güncelleme öncesi alma
	// - Hesap silme öncesi doğrulama
	//
	// # Örnek Kullanım
	// ```go
	// account, err := repo.FindByID(ctx, 123)
	// if err != nil {
	//     log.Printf("Hesap bulunamadı: %v", err)
	// }
	// if account != nil {
	//     fmt.Printf("Hesap: %+v\n", account)
	// }
	// ```
	//
	// # Önemli Notlar
	// - Hesap bulunamadıysa, nil ve hata döndürmelidir
	// - Veritabanında indekslenmiş alan olduğu için hızlı olmalıdır
	FindByID(ctx context.Context, id uint) (*Account, error)

	// Bu metod, sağlayıcı ve hesap ID'sine göre hesabı bulur.
	// OAuth2 sağlayıcılarıyla entegrasyon için kullanılır.
	//
	// # Parametreler
	// - ctx: İşlem bağlamı (context), iptal ve timeout desteği için
	// - providerID: Sağlayıcı tanımlayıcısı (örn: "google", "github")
	// - accountID: Sağlayıcının verdiği kullanıcı ID'si
	//
	// # Dönüş Değeri
	// - *Account: Bulunan hesap nesnesi (nil ise bulunamadı)
	// - error: İşlem başarılı ise nil, aksi takdirde hata
	//
	// # Kullanım Senaryoları
	// - OAuth2 callback işleminde hesap bulma
	// - Sosyal medya hesabıyla giriş yapma
	// - Mevcut hesabı kontrol etme
	//
	// # Örnek Kullanım
	// ```go
	// account, err := repo.FindByProvider(ctx, "google", "118234567890")
	// if err != nil {
	//     log.Printf("Hesap bulunamadı: %v", err)
	// }
	// if account != nil {
	//     fmt.Printf("Mevcut hesap bulundu: %d\n", account.UserID)
	// } else {
	//     // Yeni hesap oluştur
	// }
	// ```
	//
	// # Önemli Notlar
	// - providerID ve accountID kombinasyonu benzersiz olmalıdır
	// - Veritabanında indekslenmiş alanlar olduğu için hızlı olmalıdır
	// - Hesap bulunamadıysa, nil ve hata döndürmelidir
	FindByProvider(ctx context.Context, providerID, accountID string) (*Account, error)

	// Bu metod, verilen kullanıcı ID'sine ait tüm hesapları bulur.
	// Bir kullanıcının birden fazla kimlik doğrulama yöntemi olabilir.
	//
	// # Parametreler
	// - ctx: İşlem bağlamı (context), iptal ve timeout desteği için
	// - userID: Hesapları aranacak kullanıcının ID'si
	//
	// # Dönüş Değeri
	// - []Account: Bulunan hesapların listesi (boş ise hiç hesap yok)
	// - error: İşlem başarılı ise nil, aksi takdirde hata
	//
	// # Kullanım Senaryoları
	// - Kullanıcının tüm bağlı hesaplarını listeleme
	// - Hesap yönetimi sayfasında gösterme
	// - Hesap silme işleminde kontrol
	//
	// # Örnek Kullanım
	// ```go
	// accounts, err := repo.FindByUserID(ctx, 1)
	// if err != nil {
	//     log.Printf("Hesaplar alınamadı: %v", err)
	// }
	// for _, account := range accounts {
	//     fmt.Printf("Sağlayıcı: %s\n", account.ProviderID)
	// }
	// ```
	//
	// # Önemli Notlar
	// - UserID veritabanında indekslenmiş olduğu için hızlı olmalıdır
	// - Hiç hesap yoksa, boş slice ve nil hata döndürmelidir
	// - Sonuçlar ProviderID'ye göre sıralanabilir
	FindByUserID(ctx context.Context, userID uint) ([]Account, error)

	// Bu metod, mevcut bir hesabı günceller.
	// Token yenileme, scope güncelleme vb. işlemler için kullanılır.
	//
	// # Parametreler
	// - ctx: İşlem bağlamı (context), iptal ve timeout desteği için
	// - account: Güncellenecek hesap nesnesi (ID alanı zorunlu)
	//
	// # Dönüş Değeri
	// - error: İşlem başarılı ise nil, aksi takdirde hata
	//
	// # Kullanım Senaryoları
	// - AccessToken yenileme
	// - RefreshToken güncelleme
	// - Scope güncelleme
	// - Hesap bilgilerini değiştirme
	//
	// # Örnek Kullanım
	// ```go
	// account.AccessToken = "new_token"
	// account.AccessTokenExpiresAt = &newExpireTime
	// err := repo.Update(ctx, account)
	// if err != nil {
	//     log.Printf("Hesap güncellenemedi: %v", err)
	// }
	// ```
	//
	// # Önemli Notlar
	// - Hesap ID'si zorunlu olmalıdır
	// - Hesap bulunamadıysa, hata döndürmelidir
	// - UpdatedAt alanı otomatik olarak güncellenmeli
	// - Hassas veriler (token) şifreli olarak saklanmalıdır
	Update(ctx context.Context, account *Account) error

	// Bu metod, verilen ID'ye sahip hesabı siler.
	// Hesap silme işleminde kullanılır.
	//
	// # Parametreler
	// - ctx: İşlem bağlamı (context), iptal ve timeout desteği için
	// - id: Silinecek hesabın ID'si
	//
	// # Dönüş Değeri
	// - error: İşlem başarılı ise nil, aksi takdirde hata
	//
	// # Kullanım Senaryoları
	// - Kullanıcı hesabı silme
	// - Sosyal medya hesabı bağlantısını kaldırma
	// - Yanlış hesabı silme
	//
	// # Örnek Kullanım
	// ```go
	// err := repo.Delete(ctx, 123)
	// if err != nil {
	//     log.Printf("Hesap silinemedi: %v", err)
	// }
	// ```
	//
	// # Önemli Notlar
	// - Hesap bulunamadıysa, hata döndürmelidir
	// - Silme işlemi geri alınamaz (soft delete değil, hard delete)
	// - Veritabanı kısıtlamaları kontrol edilmelidir (foreign key vb.)
	// - Audit log kaydı tutulmalıdır (silinme işlemi)
	Delete(ctx context.Context, id uint) error
}
