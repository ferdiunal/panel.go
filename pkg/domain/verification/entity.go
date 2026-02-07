// Bu paket, e-posta doğrulama ve kullanıcı kimlik doğrulama işlemlerini yönetmek için
// gerekli veri yapılarını ve repository arayüzünü tanımlar.
//
// Kullanım Senaryoları:
// - Kullanıcı kayıt sırasında e-posta doğrulama
// - Şifre sıfırlama işlemleri
// - İki faktörlü kimlik doğrulama (2FA)
// - Hesap aktivasyonu
//
// Örnek Kullanım:
//   verification := &Verification{
//       Identifier: "user@example.com",
//       Token: "abc123xyz",
//       ExpiresAt: time.Now().Add(24 * time.Hour),
//   }
//   err := repo.Create(ctx, verification)
package verification

import (
	"context"
	"time"
)

// Bu yapı, doğrulama token'ı ve ilgili meta verilerini temsil eder.
//
// Amaç:
// Kullanıcı doğrulama işlemleri sırasında geçici token'ları ve bunların
// ilişkili bilgilerini depolamak için kullanılır.
//
// Alan Açıklamaları:
// - ID: Veritabanında benzersiz tanımlayıcı (Primary Key)
// - Identifier: Doğrulanacak kullanıcı tanımlayıcısı (e-posta veya kullanıcı ID'si)
// - Token: Doğrulama için kullanılan benzersiz token (örn: UUID veya random string)
// - ExpiresAt: Token'ın geçerlilik süresi bittiği zaman
// - CreatedAt: Doğrulama kaydının oluşturulduğu zaman
// - UpdatedAt: Doğrulama kaydının son güncellendiği zaman
//
// Veritabanı Özellikleri:
// - ID: Primary Key olarak tanımlanmış
// - Identifier: İndekslenmiş (hızlı arama için)
// - Token: Benzersiz indeks (aynı token iki kez oluşturulamaz)
// - ExpiresAt, CreatedAt, UpdatedAt: İndekslenmiş (zaman bazlı sorgular için)
//
// Önemli Notlar:
// - Token'lar güvenli bir şekilde oluşturulmalı ve şifrelenmelidir
// - ExpiresAt zamanı geçmiş token'lar otomatik olarak temizlenmelidir
// - Identifier alanı e-posta veya kullanıcı ID'si olabilir
//
// Kullanım Örneği:
//   v := &Verification{
//       Identifier: "user@example.com",
//       Token: generateSecureToken(),
//       ExpiresAt: time.Now().Add(24 * time.Hour),
//   }
//   if err := repo.Create(ctx, v); err != nil {
//       log.Printf("Doğrulama kaydı oluşturulamadı: %v", err)
//   }
type Verification struct {
	// ID: Doğrulama kaydının benzersiz tanımlayıcısı
	// Veritabanında otomatik olarak artan bir sayı olarak saklanır
	ID uint `json:"id" gorm:"primaryKey"`

	// Identifier: Doğrulanacak kullanıcının tanımlayıcısı
	// Genellikle e-posta adresi veya kullanıcı ID'si olabilir
	// İndekslenmiş: Hızlı arama ve filtreleme için
	// Örnek: "user@example.com" veya "user_123"
	Identifier string `json:"identifier" gorm:"index"`

	// Token: Doğrulama işlemi için kullanılan benzersiz token
	// Benzersiz indeks: Aynı token iki kez oluşturulamaz
	// Güvenlik: Kriptografik olarak güvenli bir şekilde oluşturulmalı
	// Örnek: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
	Token string `json:"token" gorm:"uniqueIndex"`

	// ExpiresAt: Token'ın geçerlilik süresi bittiği zaman
	// İndekslenmiş: Süresi geçmiş token'ları temizlemek için
	// Örnek: 2024-01-15 14:30:00 UTC
	ExpiresAt time.Time `json:"expiresAt" gorm:"index"`

	// CreatedAt: Doğrulama kaydının oluşturulduğu zaman
	// İndekslenmiş: Zaman bazlı sorgular için
	// GORM tarafından otomatik olarak ayarlanır
	CreatedAt time.Time `json:"createdAt" gorm:"index"`

	// UpdatedAt: Doğrulama kaydının son güncellendiği zaman
	// İndekslenmiş: Zaman bazlı sorgular için
	// GORM tarafından otomatik olarak ayarlanır
	UpdatedAt time.Time `json:"updatedAt" gorm:"index"`
}

// Bu interface, doğrulama verilerinin veritabanında yönetilmesi için
// gerekli tüm operasyonları tanımlar.
//
// Amaç:
// Repository pattern kullanarak veri erişim katmanını soyutlamak ve
// farklı veritabanı uygulamalarını desteklemek.
//
// Sorumluluklar:
// - Doğrulama kayıtlarını oluşturma
// - Token'a göre doğrulama kayıtlarını bulma
// - Doğrulama kayıtlarını silme
// - Identifier'a göre doğrulama kayıtlarını silme
//
// Bağlam (Context) Kullanımı:
// Tüm metotlar context.Context parametresi alır. Bu:
// - İşlemleri iptal etme yeteneği sağlar
// - Zaman aşımı (timeout) ayarlamaya izin verir
// - İstek bazlı değerleri taşımaya olanak tanır
//
// Hata Yönetimi:
// Tüm metotlar error döndürür. Olası hatalar:
// - Veritabanı bağlantı hataları
// - Benzersizlik ihlalleri (Token için)
// - Kayıt bulunamadı hataları
// - Bağlam iptal hataları
//
// Kullanım Örneği:
//   repo := NewVerificationRepository(db)
//
//   // Doğrulama kaydı oluşturma
//   v := &Verification{...}
//   if err := repo.Create(ctx, v); err != nil {
//       return err
//   }
//
//   // Token'a göre bulma
//   found, err := repo.FindByToken(ctx, "token123")
//   if err != nil {
//       return err
//   }
//
//   // Doğrulama tamamlandıktan sonra silme
//   if err := repo.Delete(ctx, found.ID); err != nil {
//       return err
//   }
type Repository interface {
	// Bu metod, yeni bir doğrulama kaydını veritabanına ekler.
	//
	// Parametreler:
	// - ctx: İşlem bağlamı (zaman aşımı ve iptal için)
	// - v: Oluşturulacak Verification yapısı pointer'ı
	//
	// Dönüş Değeri:
	// - error: İşlem başarılı ise nil, aksi takdirde hata
	//
	// Olası Hatalar:
	// - Benzersizlik ihlali: Token zaten varsa
	// - Veritabanı hatası: Bağlantı problemi
	// - Geçersiz veri: Zorunlu alanlar eksikse
	//
	// Önemli Notlar:
	// - Token benzersiz olmalıdır (uniqueIndex)
	// - Identifier ve Token alanları doldurulmalıdır
	// - CreatedAt ve UpdatedAt otomatik olarak ayarlanır
	//
	// Kullanım Örneği:
	//   v := &Verification{
	//       Identifier: "user@example.com",
	//       Token: "secure_token_123",
	//       ExpiresAt: time.Now().Add(24 * time.Hour),
	//   }
	//   if err := repo.Create(ctx, v); err != nil {
	//       log.Printf("Hata: %v", err)
	//   }
	Create(ctx context.Context, v *Verification) error

	// Bu metod, verilen token'a göre doğrulama kaydını bulur.
	//
	// Parametreler:
	// - ctx: İşlem bağlamı (zaman aşımı ve iptal için)
	// - token: Aranacak token değeri
	//
	// Dönüş Değeri:
	// - *Verification: Bulunan doğrulama kaydı pointer'ı
	// - error: İşlem başarılı ise nil, aksi takdirde hata
	//
	// Olası Hatalar:
	// - Kayıt bulunamadı: Token veritabanında yoksa
	// - Veritabanı hatası: Bağlantı problemi
	// - Bağlam iptal: Context iptal edilmişse
	//
	// Önemli Notlar:
	// - Token benzersiz indekslenmiş olduğu için hızlı arama yapılır
	// - Süresi geçmiş token'lar da döndürülebilir (kontrol etmek gerekir)
	// - Döndürülen pointer'ı değiştirmek orijinal kaydı etkilemez
	//
	// Kullanım Örneği:
	//   v, err := repo.FindByToken(ctx, "secure_token_123")
	//   if err != nil {
	//       log.Printf("Token bulunamadı: %v", err)
	//       return
	//   }
	//
	//   // Token'ın süresi geçip geçmediğini kontrol et
	//   if time.Now().After(v.ExpiresAt) {
	//       log.Println("Token süresi geçmiş")
	//       return
	//   }
	FindByToken(ctx context.Context, token string) (*Verification, error)

	// Bu metod, verilen ID'ye göre doğrulama kaydını siler.
	//
	// Parametreler:
	// - ctx: İşlem bağlamı (zaman aşımı ve iptal için)
	// - id: Silinecek doğrulama kaydının ID'si
	//
	// Dönüş Değeri:
	// - error: İşlem başarılı ise nil, aksi takdirde hata
	//
	// Olası Hatalar:
	// - Kayıt bulunamadı: ID veritabanında yoksa
	// - Veritabanı hatası: Bağlantı problemi
	// - Bağlam iptal: Context iptal edilmişse
	//
	// Önemli Notlar:
	// - Silme işlemi geri alınamaz
	// - Doğrulama başarılı olduktan sonra kayıt silinmelidir
	// - Süresi geçmiş kayıtları temizlemek için kullanılabilir
	//
	// Kullanım Örneği:
	//   if err := repo.Delete(ctx, 42); err != nil {
	//       log.Printf("Silme hatası: %v", err)
	//   } else {
	//       log.Println("Doğrulama kaydı başarıyla silindi")
	//   }
	Delete(ctx context.Context, id uint) error

	// Bu metod, verilen identifier'a göre tüm doğrulama kayıtlarını siler.
	//
	// Parametreler:
	// - ctx: İşlem bağlamı (zaman aşımı ve iptal için)
	// - identifier: Silinecek kayıtların identifier'ı (e-posta veya kullanıcı ID'si)
	//
	// Dönüş Değeri:
	// - error: İşlem başarılı ise nil, aksi takdirde hata
	//
	// Olası Hatalar:
	// - Veritabanı hatası: Bağlantı problemi
	// - Bağlam iptal: Context iptal edilmişse
	//
	// Önemli Notlar:
	// - Aynı identifier'a ait tüm kayıtlar silinir
	// - Silme işlemi geri alınamaz
	// - Bir kullanıcı için birden fazla doğrulama kaydı olabilir
	// - Yeni doğrulama işlemi başlatılmadan önce eski kayıtları temizlemek için kullanılır
	//
	// Kullanım Senaryoları:
	// - Kullanıcı yeni doğrulama isteği gönderdiğinde eski token'ları temizle
	// - Kullanıcı hesabını silerken ilgili tüm doğrulama kayıtlarını sil
	// - Şifre sıfırlama işleminde önceki token'ları geçersiz kıl
	//
	// Kullanım Örneği:
	//   if err := repo.DeleteByIdentifier(ctx, "user@example.com"); err != nil {
	//       log.Printf("Silme hatası: %v", err)
	//   } else {
	//       log.Println("Tüm doğrulama kayıtları silindi")
	//   }
	DeleteByIdentifier(ctx context.Context, identifier string) error
}
