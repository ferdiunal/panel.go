// Bu paket, kimlik doğrulama (authentication) işlemlerini yönetir.
// Kullanıcı kaydı, giriş, oturum yönetimi ve şifre sıfırlama gibi temel auth işlevlerini sağlar.
package auth

import (
	"context"
	"errors"
	"time"

	"github.com/ferdiunal/panel.go/pkg/domain/account"
	"github.com/ferdiunal/panel.go/pkg/domain/session"
	"github.com/ferdiunal/panel.go/pkg/domain/user"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Kimlik doğrulama işlemleri sırasında oluşabilecek hata değişkenleri.
// Bu hatalar, kullanıcı kaydı, giriş ve oturum doğrulama sırasında döndürülür.
var (
	// ErrUserNotFound: Belirtilen e-posta adresiyle kullanıcı bulunamadığında döndürülür.
	// Kullanım senaryosu: Giriş sırasında e-posta veritabanında yoksa bu hata oluşur.
	ErrUserNotFound = errors.New("user not found")

	// ErrInvalidCredentials: E-posta veya şifre yanlış olduğunda döndürülür.
	// Kullanım senaryosu: Giriş sırasında şifre eşleşmediğinde veya hesap bulunamadığında.
	// Güvenlik notu: Hem e-posta hem şifre hataları için aynı mesaj kullanılır (timing attack önleme).
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrEmailAlreadyExists: Kayıt sırasında e-posta adresi zaten kullanımda olduğunda döndürülür.
	// Kullanım senaryosu: Yeni kullanıcı kaydı sırasında e-posta benzersizlik kontrolü başarısız olduğunda.
	ErrEmailAlreadyExists = errors.New("email already exists")
)

// Bu yapı, kimlik doğrulama hizmetinin ana bileşenidir.
// Kullanıcı kaydı, giriş, oturum yönetimi ve şifre sıfırlama işlemlerini koordine eder.
//
// Yapı Alanları:
//   - userRepo: Kullanıcı verilerini veritabanından yönetmek için repository
//   - sessionRepo: Oturum verilerini veritabanından yönetmek için repository
//   - accountRepo: Hesap (account) verilerini veritabanından yönetmek için repository
//
// Kullanım Senaryosu:
//   Service, dependency injection pattern kullanılarak oluşturulur.
//   Tüm repository'ler constructor aracılığıyla enjekte edilir.
//
// Örnek:
//   service := auth.NewService(userRepo, sessionRepo, accountRepo)
//   user, err := service.RegisterEmail(ctx, "John", "john@example.com", "password123")
type Service struct {
	// userRepo: Kullanıcı bilgilerini sorgulamak ve oluşturmak için kullanılan repository
	userRepo user.Repository

	// sessionRepo: Oturum bilgilerini yönetmek için kullanılan repository
	sessionRepo session.Repository

	// accountRepo: Hesap bilgilerini (şifre, provider vb.) yönetmek için kullanılan repository
	accountRepo account.Repository
}

// Bu fonksiyon, kimlik doğrulama hizmetinin yeni bir örneğini oluşturur.
// Dependency injection pattern kullanarak tüm repository'leri alır ve Service yapısını başlatır.
//
// Parametreler:
//   - u (user.Repository): Kullanıcı verilerini yönetmek için repository
//   - s (session.Repository): Oturum verilerini yönetmek için repository
//   - a (account.Repository): Hesap verilerini yönetmek için repository
//
// Dönüş Değeri:
//   - *Service: Yapılandırılmış Service pointer'ı
//
// Kullanım Senaryosu:
//   Uygulama başlangıcında, tüm repository'ler oluşturulduktan sonra
//   NewService çağrılarak Service örneği oluşturulur.
//
// Örnek:
//   userRepo := user.NewRepository(db)
//   sessionRepo := session.NewRepository(db)
//   accountRepo := account.NewRepository(db)
//   authService := NewService(userRepo, sessionRepo, accountRepo)
//
// Önemli Notlar:
//   - Tüm repository parametreleri zorunludur (nil olamaz)
//   - Repository'ler dış tarafından yönetilir, Service tarafından kapatılmaz
func NewService(u user.Repository, s session.Repository, a account.Repository) *Service {
	return &Service{
		userRepo:    u,
		sessionRepo: s,
		accountRepo: a,
	}
}

// Bu metod, e-posta ve şifre kullanarak yeni bir kullanıcı kaydı oluşturur.
// Kullanıcı kaydı sırasında şifre bcrypt ile hash'lenir ve veritabanına kaydedilir.
// İlk kayıt yapan kullanıcı otomatik olarak admin rolü alır, sonraki kullanıcılar "user" rolü alır.
//
// Parametreler:
//   - ctx (context.Context): İşlem için context (timeout, cancellation vb.)
//   - name (string): Kullanıcının adı (örn: "John Doe")
//   - email (string): Kullanıcının e-posta adresi (benzersiz olmalı)
//   - password (string): Kullanıcının şifresi (düz metin, bcrypt ile hash'lenecek)
//
// Dönüş Değeri:
//   - *user.User: Başarılı kayıt durumunda oluşturulan kullanıcı nesnesi
//   - error: Hata durumunda (e-posta zaten var, şifre hash'leme hatası vb.)
//
// Olası Hatalar:
//   - ErrEmailAlreadyExists: E-posta adresi zaten kullanımda
//   - bcrypt hata: Şifre hash'leme başarısız
//   - Repository hata: Veritabanı işlemi başarısız
//
// Kullanım Senaryosu:
//   Yeni kullanıcı kaydı sırasında çağrılır. Örneğin, web uygulamasında
//   kayıt formundan gelen veriler bu metoda iletilir.
//
// Örnek:
//   user, err := authService.RegisterEmail(ctx, "John Doe", "john@example.com", "securePassword123")
//   if err != nil {
//       if err == auth.ErrEmailAlreadyExists {
//           // E-posta zaten kayıtlı
//       }
//       return err
//   }
//   // Kullanıcı başarıyla kaydedildi
//   fmt.Printf("Yeni kullanıcı: %s (ID: %s, Rol: %s)\n", user.Name, user.ID, user.Role)
//
// Önemli Notlar:
//   - İlk kayıt yapan kullanıcı otomatik olarak admin rolü alır
//   - Şifre bcrypt.DefaultCost (12) ile hash'lenir
//   - E-posta adresi benzersizlik kontrolü yapılır
//   - Kullanıcı ve hesap (account) iki ayrı işlemde oluşturulur
//   - Hesap oluşturma başarısız olursa, kullanıcı veritabanında kalabilir (transaction önerilir)
//   - EmailVerified başlangıçta false olarak ayarlanır
func (s *Service) RegisterEmail(ctx context.Context, name, email, password string) (*user.User, error) {
	// E-posta adresinin zaten kullanımda olup olmadığını kontrol et
	existing, _ := s.userRepo.FindByEmail(ctx, email)
	if existing != nil {
		return nil, ErrEmailAlreadyExists
	}

	// Şifreyi bcrypt ile hash'le
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Rol belirle: İlk kullanıcı admin, diğerleri user
	role := "user"
	userCount, err := s.userRepo.Count(ctx)
	if err == nil && userCount == 0 {
		// Bu ilk kullanıcı, admin rolü ver
		role = "admin"
	}

	// Kullanıcı nesnesi oluştur
	u := &user.User{
		Name:          name,
		Email:         email,
		EmailVerified: false,
		Role:          role,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Kullanıcıyı veritabanına kaydet
	if err := s.userRepo.CreateUser(ctx, u); err != nil {
		return nil, err
	}

	// Hesap (account) nesnesi oluştur
	acc := &account.Account{
		UserID:     u.ID,
		ProviderID: "credential",
		AccountID:  nil, // Credential provider'ın harici hesap ID'si yoktur
		Password:   string(hashed),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Hesabı veritabanına kaydet
	if err := s.accountRepo.Create(ctx, acc); err != nil {
		// Not: Hesap oluşturma başarısız olursa, kullanıcı veritabanında kalabilir
		// Üretim ortamında transaction kullanılması önerilir
		return nil, err
	}

	return u, nil
}

// Bu metod, e-posta ve şifre kullanarak kullanıcı girişi gerçekleştirir.
// Başarılı giriş durumunda yeni bir oturum (session) oluşturur ve döndürür.
// İstemci IP adresi ve User-Agent bilgileri oturum kaydında saklanır.
//
// Parametreler:
//   - ctx (context.Context): İşlem için context (timeout, cancellation vb.)
//   - email (string): Kullanıcının e-posta adresi
//   - password (string): Kullanıcının şifresi (düz metin)
//   - ip (string): İstemcinin IP adresi (örn: "192.168.1.1")
//   - userAgent (string): İstemcinin User-Agent bilgisi (örn: "Mozilla/5.0...")
//
// Dönüş Değeri:
//   - *session.Session: Başarılı giriş durumunda oluşturulan oturum nesnesi
//   - error: Hata durumunda (kullanıcı bulunamadı, şifre yanlış vb.)
//
// Olası Hatalar:
//   - ErrInvalidCredentials: E-posta bulunamadı veya şifre yanlış
//   - Repository hata: Veritabanı işlemi başarısız
//
// Kullanım Senaryosu:
//   Web uygulamasında giriş formundan gelen e-posta ve şifre bilgileri
//   bu metoda iletilir. Başarılı giriş durumunda oturum token'ı
//   istemciye gönderilir ve cookie'de saklanır.
//
// Örnek:
//   session, err := authService.LoginEmail(
//       ctx,
//       "john@example.com",
//       "securePassword123",
//       "192.168.1.100",
//       "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
//   )
//   if err != nil {
//       if err == auth.ErrInvalidCredentials {
//           // E-posta veya şifre yanlış
//       }
//       return err
//   }
//   // Giriş başarılı, oturum token'ı: session.Token
//   fmt.Printf("Oturum oluşturuldu: %s (Süresi: %v)\n", session.Token, session.ExpiresAt)
//
// Önemli Notlar:
//   - Oturum 7 gün (168 saat) geçerlidir
//   - Şifre bcrypt.CompareHashAndPassword ile doğrulanır
//   - Hem e-posta hem şifre hataları için aynı hata döndürülür (güvenlik)
//   - IP adresi ve User-Agent oturum kaydında saklanır (güvenlik denetimi için)
//   - Oturum token'ı UUID v7 formatında oluşturulur
//   - Credential provider'ı için hesap aranır
func (s *Service) LoginEmail(ctx context.Context, email, password string, ip, userAgent string) (*session.Session, error) {
	// Kullanıcıyı e-posta adresine göre bul
	u, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Credential provider'ı için hesabı bul
	// Hesap repository'den doğrudan sorgula veya kullanıcı hesaplarından ara
	acc, err := s.accountRepo.FindByProvider(ctx, "credential", email)
	if err != nil {
		// Fallback: Kullanıcının tüm hesaplarını al ve credential provider'ını ara
		accounts, err := s.accountRepo.FindByUserID(ctx, u.ID)
		if err != nil {
			return nil, ErrInvalidCredentials
		}
		found := false
		for _, a := range accounts {
			if a.ProviderID == "credential" {
				acc = &a
				found = true
				break
			}
		}
		if !found {
			return nil, ErrInvalidCredentials
		}
	}

	// Şifreyi doğrula
	if err := bcrypt.CompareHashAndPassword([]byte(acc.Password), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Yeni oturum oluştur
	sessionToken, _ := uuid.NewV7()
	sess := &session.Session{
		UserID:    u.ID,
		Token:     sessionToken.String(),
		ExpiresAt: time.Now().Add(24 * 7 * time.Hour), // 7 gün
		IPAddress: ip,
		UserAgent: userAgent,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Oturumu veritabanına kaydet
	if err := s.sessionRepo.Create(ctx, sess); err != nil {
		return nil, err
	}

	// Oturumu kullanıcı bilgileriyle birlikte getir
	return s.sessionRepo.FindByToken(ctx, sess.Token)
}

// Bu metod, verilen oturum token'ını doğrular ve geçerliliğini kontrol eder.
// Token'ın veritabanında var olup olmadığını ve süresi dolup dolmadığını kontrol eder.
//
// Parametreler:
//   - ctx (context.Context): İşlem için context (timeout, cancellation vb.)
//   - token (string): Doğrulanacak oturum token'ı (UUID v7 formatında)
//
// Dönüş Değeri:
//   - *session.Session: Token geçerli ise oturum nesnesi
//   - error: Token bulunamadı veya süresi dolmuş ise hata
//
// Olası Hatalar:
//   - "session expired": Oturum süresi dolmuş
//   - Repository hata: Veritabanı işlemi başarısız
//
// Kullanım Senaryosu:
//   Her HTTP isteğinde, istemciden gelen oturum token'ı bu metoda iletilir.
//   Token geçerli ise, kullanıcı kimliği doğrulanmış kabul edilir.
//   Token geçersiz ise, kullanıcı yeniden giriş yapmaya yönlendirilir.
//
// Örnek:
//   session, err := authService.ValidateSession(ctx, "550e8400-e29b-41d4-a716-446655440000")
//   if err != nil {
//       if err.Error() == "session expired" {
//           // Oturum süresi dolmuş, yeniden giriş gerekli
//       }
//       return err
//   }
//   // Token geçerli, session.UserID ile kullanıcı bilgisine erişebilir
//   fmt.Printf("Oturum geçerli, Kullanıcı ID: %s\n", session.UserID)
//
// Önemli Notlar:
//   - Oturum süresi dolmuş ise, oturum veritabanında kalabilir (temizleme yapılmaz)
//   - Token'ın tam eşleşmesi gerekir (case-sensitive)
//   - Oturum uzatma (extension) şu anda yapılmamaktadır (TODO)
//   - Middleware'de kullanılması önerilir
func (s *Service) ValidateSession(ctx context.Context, token string) (*session.Session, error) {
	// Token'a göre oturumu bul
	sess, err := s.sessionRepo.FindByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	// Oturum süresinin dolup dolmadığını kontrol et
	if sess.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("session expired")
	}

	// Opsiyonel: Oturum süresini uzat?
	// TODO: Oturum uzatma özelliği eklenebilir
	return sess, nil
}

// Bu metod, verilen oturum token'ını iptal eder ve kullanıcıyı çıkış yapmış duruma getirir.
// Token'a karşılık gelen oturumu veritabanından siler.
//
// Parametreler:
//   - ctx (context.Context): İşlem için context (timeout, cancellation vb.)
//   - token (string): İptal edilecek oturum token'ı (UUID v7 formatında)
//
// Dönüş Değeri:
//   - error: Silme işlemi başarısız ise hata, başarılı ise nil
//
// Olası Hatalar:
//   - Repository hata: Veritabanı işlemi başarısız
//
// Kullanım Senaryosu:
//   Kullanıcı çıkış butonuna tıkladığında veya oturum süresi dolduğunda
//   bu metod çağrılarak oturum iptal edilir.
//
// Örnek:
//   err := authService.Logout(ctx, "550e8400-e29b-41d4-a716-446655440000")
//   if err != nil {
//       // Çıkış işlemi başarısız
//       return err
//   }
//   // Çıkış başarılı, kullanıcı yeniden giriş yapmaya yönlendir
//   fmt.Println("Çıkış başarılı")
//
// Önemli Notlar:
//   - Token'ın veritabanında var olması gerekmez (idempotent)
//   - Çıkış sonrası istemci tarafında cookie silinmelidir
//   - Middleware'de kullanılması önerilir
func (s *Service) Logout(ctx context.Context, token string) error {
	return s.sessionRepo.DeleteByToken(ctx, token)
}

// Bu metod, şifresi unutulan kullanıcılar için şifre sıfırlama işlemini başlatır.
// Güvenlik nedeniyle, e-posta adresinin veritabanında var olup olmadığı açıklanmaz.
// Şu anda sadece başarı döndürür, gerçek şifre sıfırlama e-posta gönderimi yapılmamaktadır.
//
// Parametreler:
//   - ctx (context.Context): İşlem için context (timeout, cancellation vb.)
//   - email (string): Şifresi sıfırlanacak kullanıcının e-posta adresi
//
// Dönüş Değeri:
//   - error: Hata durumunda (şu anda her zaman nil döndürür)
//
// Kullanım Senaryosu:
//   Kullanıcı "Şifremi Unuttum" bağlantısına tıkladığında, e-posta adresi
//   bu metoda iletilir. Başarılı ise, kullanıcıya e-posta gönderilir.
//
// Örnek:
//   err := authService.ForgotPassword(ctx, "john@example.com")
//   if err != nil {
//       return err
//   }
//   // Şifre sıfırlama e-postası gönderildi (veya gönderilecek)
//   fmt.Println("Şifre sıfırlama bağlantısı e-postanıza gönderildi")
//
// Önemli Notlar:
//   - Güvenlik: E-posta adresinin var olup olmadığı açıklanmaz (timing attack önleme)
//   - TODO: Şu anda sadece başarı döndürür, gerçek implementasyon yapılmamıştır
//   - Üretim ortamında aşağıdaki adımlar uygulanmalıdır:
//     1. Benzersiz şifre sıfırlama token'ı oluştur
//     2. Token'ı veritabanında sakla (süresi sınırlı, örn: 1 saat)
//     3. Şifre sıfırlama bağlantısı içeren e-posta gönder
//     4. Başarı döndür
//   - E-posta gönderimi için SMTP veya e-posta servisi (SendGrid, AWS SES vb.) gerekli
func (s *Service) ForgotPassword(ctx context.Context, email string) error {
	// Kullanıcının e-posta adresine göre var olup olmadığını kontrol et
	u, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		// Güvenlik: E-posta adresinin var olup olmadığını açıklama
		return nil
	}

	if u == nil {
		// Güvenlik: E-posta adresinin var olup olmadığını açıklama
		return nil
	}

	// TODO: Şifre sıfırlama token'ı oluştur ve e-posta gönder
	// Üretim ortamında aşağıdaki adımlar uygulanmalıdır:
	// 1. Benzersiz şifre sıfırlama token'ı oluştur (UUID v7 veya random string)
	// 2. Token'ı veritabanında sakla (süresi sınırlı, örn: 1 saat)
	// 3. Şifre sıfırlama bağlantısı içeren e-posta gönder
	//    Örn: https://example.com/reset-password?token=<token>
	// 4. Başarı döndür
	//
	// Örnek implementasyon:
	//   resetToken := uuid.NewV7().String()
	//   resetTokenHash := bcrypt.GenerateFromPassword([]byte(resetToken), bcrypt.DefaultCost)
	//   resetRecord := &PasswordReset{
	//       UserID:    u.ID,
	//       Token:     string(resetTokenHash),
	//       ExpiresAt: time.Now().Add(1 * time.Hour),
	//   }
	//   s.passwordResetRepo.Create(ctx, resetRecord)
	//   s.emailService.SendPasswordResetEmail(u.Email, resetToken)

	return nil
}
