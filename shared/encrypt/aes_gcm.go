// Package encrypt, AES-GCM (Galois/Counter Mode) tabanlı güvenli şifreleme ve şifre çözme
// işlemleri için araçlar sağlar.
//
// # Genel Bakış
//
// Bu paket, modern kriptografi standartlarına uygun, hem gizlilik hem de bütünlük koruması
// sağlayan AES-GCM şifreleme yöntemini kullanır. GCM modu, şifrelenmiş verinin hem
// yetkisiz erişime karşı korunmasını hem de değiştirilmediğinin doğrulanmasını sağlar.
//
// # Güvenlik Özellikleri
//
// - **Authenticated Encryption**: Hem şifreleme hem de kimlik doğrulama sağlar
// - **Tamper Protection**: Veri değişikliklerini otomatik olarak tespit eder
// - **Nonce-based**: Her şifreleme işlemi için benzersiz nonce kullanır
// - **AEAD (Authenticated Encryption with Associated Data)**: Endüstri standardı
//
// # Kullanım Senaryoları
//
// - Hassas kullanıcı verilerinin veritabanında saklanması
// - API token'larının güvenli şekilde saklanması
// - Şifreli mesajlaşma sistemleri
// - Dosya şifreleme işlemleri
// - Session verilerinin güvenli saklanması
//
// # Örnek Kullanım
//
// ```go
// // 32 byte'lık bir anahtar oluştur (AES-256 için)
// key := make([]byte, 32)
// if _, err := rand.Read(key); err != nil {
//     log.Fatal(err)
// }
//
// // Şifreleme nesnesi oluştur
// crypt := encrypt.NewCryptGCM(key)
//
// // Veri şifrele
// encrypted, err := crypt.Encrypt("hassas veri")
// if err != nil {
//     log.Fatal(err)
// }
//
// // Veriyi çöz
// decrypted, err := crypt.Decrypt(encrypted)
// if err != nil {
//     log.Fatal(err)
// }
// ```
//
// # Önemli Notlar
//
// - Anahtar boyutu 16, 24 veya 32 byte olmalıdır (AES-128, AES-192, AES-256)
// - Aynı anahtarı güvenli bir şekilde saklamalısınız
// - Nonce'lar otomatik olarak oluşturulur ve şifreli metne eklenir
// - Base64 encoding kullanılarak metin formatında saklanabilir
package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// encryptGCM, düz metni AES-GCM (Authenticated Encryption) kullanarak şifreler.
//
// Bu fonksiyon, verilen düz metni AES-GCM algoritması ile şifreler ve hem gizlilik
// hem de bütünlük koruması sağlar. GCM modu, şifrelenmiş verinin değiştirilmediğini
// doğrulayan bir kimlik doğrulama etiketi (authentication tag) ekler.
//
// # Parametreler
//
// - `plaintext`: Şifrelenecek düz metin (string formatında)
// - `key`: Şifreleme anahtarı (16, 24 veya 32 byte olmalı)
//   - 16 byte: AES-128
//   - 24 byte: AES-192
//   - 32 byte: AES-256 (önerilen)
//
// # Dönüş Değerleri
//
// - `string`: Base64 kodlanmış şifreli metin (nonce + ciphertext + auth tag)
// - `error`: İşlem başarısız olursa hata mesajı
//
// # Çalışma Prensibi
//
// 1. Verilen anahtar ile AES cipher bloğu oluşturulur
// 2. GCM modu aktif edilir
// 3. Kriptografik olarak güvenli rastgele bir nonce üretilir
// 4. Düz metin şifrelenir ve kimlik doğrulama etiketi eklenir
// 5. Nonce + şifreli metin birleştirilerek Base64 ile kodlanır
//
// # Güvenlik Notları
//
// - Her şifreleme işlemi için benzersiz bir nonce kullanılır
// - Nonce, şifreli metnin başına eklenir (ilk 12 byte)
// - Kimlik doğrulama etiketi, şifreli metnin sonuna eklenir
// - Aynı anahtar ile aynı nonce asla tekrar kullanılmamalıdır
//
// # Hata Durumları
//
// - Geçersiz anahtar boyutu
// - Nonce üretimi başarısız olursa
// - GCM modu oluşturulamazsa
//
// # Örnek
//
// ```go
// key := []byte("32-byte-long-secret-key-here!!!!")
// encrypted, err := encryptGCM("gizli mesaj", key)
// if err != nil {
//     log.Fatal(err)
// }
// fmt.Println(encrypted) // Base64 kodlu şifreli metin
// ```
func encryptGCM(plaintext string, key []byte) (string, error) {
	// AES cipher bloğunu oluştur
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("cipher oluşturulamadı: %w", err)
	}

	// GCM modu, kimlik doğrulamalı şifreleme sağlar
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("GCM modu oluşturulamadı: %w", err)
	}

	// Nonce (number used once - bir kez kullanılan sayı) oluştur
	// Nonce, her şifreleme işlemi için benzersiz olmalıdır
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("nonce üretilemedi: %w", err)
	}

	// Şifrele ve kimlik doğrulama etiketi ekle
	// GCM'in Seal metodu, kimlik doğrulama etiketini şifreli metnin sonuna ekler
	// İlk parametre (nonce) şifreli metnin başına eklenir
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	// Base64 ile kodla ve döndür
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// decryptGCM, şifreli metni AES-GCM kullanarak çözer ve kimlik doğrulaması yapar.
//
// Bu fonksiyon, Base64 kodlanmış şifreli metni çözer ve hem şifre çözme hem de
// veri bütünlüğü doğrulaması yapar. Eğer veri değiştirilmişse veya kimlik doğrulama
// etiketi eşleşmezse, fonksiyon hata döndürür.
//
// # Parametreler
//
// - `b64cipher`: Base64 kodlanmış şifreli metin (nonce + ciphertext + auth tag içerir)
// - `key`: Şifre çözme anahtarı (şifrelemede kullanılan anahtar ile aynı olmalı)
//   - 16 byte: AES-128
//   - 24 byte: AES-192
//   - 32 byte: AES-256
//
// # Dönüş Değerleri
//
// - `string`: Çözülmüş düz metin
// - `error`: İşlem başarısız olursa hata mesajı
//
// # Çalışma Prensibi
//
// 1. Base64 kodlanmış veri çözülür
// 2. Verilen anahtar ile AES cipher bloğu oluşturulur
// 3. GCM modu aktif edilir
// 4. Nonce ve şifreli metin ayrıştırılır
// 5. Şifreli metin çözülür ve kimlik doğrulama etiketi kontrol edilir
// 6. Doğrulama başarılıysa düz metin döndürülür
//
// # Güvenlik Özellikleri
//
// - **Otomatik Doğrulama**: Veri bütünlüğü otomatik olarak kontrol edilir
// - **Tamper Detection**: Veri değiştirilmişse hata döndürülür
// - **Authentication Tag**: GCM'in kimlik doğrulama etiketi doğrulanır
// - **Nonce Extraction**: Nonce şifreli metnin başından otomatik çıkarılır
//
// # Hata Durumları
//
// - Base64 çözme hatası (geçersiz format)
// - Geçersiz anahtar boyutu
// - Şifreli metin çok kısa (nonce boyutundan küçük)
// - Kimlik doğrulama başarısız (veri değiştirilmiş)
// - Şifre çözme hatası
//
// # Önemli Notlar
//
// - Şifreleme ve şifre çözme için aynı anahtar kullanılmalıdır
// - Veri değiştirilmişse fonksiyon hata döndürür (güvenlik özelliği)
// - Nonce, şifreli metnin ilk 12 byte'ında saklanır
// - Kimlik doğrulama etiketi, şifreli metnin sonunda saklanır
//
// # Örnek
//
// ```go
// key := []byte("32-byte-long-secret-key-here!!!!")
// encrypted := "base64-encoded-ciphertext..."
// decrypted, err := decryptGCM(encrypted, key)
// if err != nil {
//     log.Fatal(err) // Veri değiştirilmiş veya geçersiz anahtar
// }
// fmt.Println(decrypted) // Orijinal düz metin
// ```
//
// # Güvenlik Uyarıları
//
// - Yanlış anahtar kullanılırsa şifre çözme başarısız olur
// - Veri manipülasyonu otomatik olarak tespit edilir
// - Hata mesajları hassas bilgi içermez (timing attack koruması)
func decryptGCM(b64cipher string, key []byte) (string, error) {
	// Base64 kodlanmış veriyi çöz
	data, err := base64.StdEncoding.DecodeString(b64cipher)
	if err != nil {
		return "", fmt.Errorf("base64 çözülemedi: %w", err)
	}

	// AES cipher bloğunu oluştur
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("cipher oluşturulamadı: %w", err)
	}

	// GCM modu oluştur
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("GCM modu oluşturulamadı: %w", err)
	}

	// Nonce boyutunu al ve veri uzunluğunu kontrol et
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("şifreli metin çok kısa")
	}

	// Nonce ve şifreli metni ayır
	// İlk nonceSize byte nonce, geri kalanı şifreli metin + auth tag
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	// Şifreyi çöz ve kimlik doğrulaması yap
	// Open metodu, kimlik doğrulama etiketi eşleşmezse hata döndürür
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("şifre çözme veya doğrulama başarısız: %w", err)
	}

	return string(plaintext), nil
}

// CryptGCM, bu yapı AES-GCM (Galois/Counter Mode) kullanarak kimlik doğrulamalı
// şifreleme işlemleri gerçekleştirir.
//
// # Genel Bakış
//
// CryptGCM, modern kriptografi standartlarına uygun, güvenli şifreleme ve şifre çözme
// işlemleri için kullanılan bir yapıdır. AES-GCM algoritması, hem veri gizliliği
// hem de veri bütünlüğü sağlar, bu sayede şifrelenmiş verinin yetkisiz erişime
// karşı korunması ve değiştirilmediğinin doğrulanması mümkün olur.
//
// # Özellikler
//
// - **Thread-Safe**: Aynı CryptGCM instance'ı birden fazla goroutine'den güvenle kullanılabilir
// - **Stateless**: Her şifreleme işlemi bağımsızdır, durum saklanmaz
// - **AEAD**: Authenticated Encryption with Associated Data standardını uygular
// - **Nonce Management**: Nonce'lar otomatik olarak üretilir ve yönetilir
// - **Base64 Encoding**: Şifreli metinler Base64 ile kodlanır, metin formatında saklanabilir
//
// # Kullanım Senaryoları
//
// - **Veritabanı Şifreleme**: Hassas kullanıcı verilerinin (kredi kartı, SSN, vb.) şifrelenmesi
// - **API Token Yönetimi**: OAuth token'ları, API anahtarlarının güvenli saklanması
// - **Session Yönetimi**: Kullanıcı oturum verilerinin şifrelenmesi
// - **Dosya Şifreleme**: Hassas dosyaların şifrelenmesi
// - **Mesajlaşma**: Uçtan uca şifreli mesajlaşma sistemleri
// - **Konfigürasyon Güvenliği**: Hassas yapılandırma değerlerinin şifrelenmesi
//
// # Güvenlik Özellikleri
//
// - **AES-256 Desteği**: 32 byte anahtar ile en güçlü AES şifreleme
// - **GCM Mode**: Galois/Counter Mode ile hem şifreleme hem kimlik doğrulama
// - **Tamper Detection**: Veri değişikliklerini otomatik tespit eder
// - **Unique Nonce**: Her şifreleme için benzersiz nonce kullanır
// - **Authentication Tag**: 16 byte'lık kimlik doğrulama etiketi
//
// # Performans Karakteristikleri
//
// - **Hızlı**: GCM modu, donanım hızlandırması ile çok hızlıdır
// - **Paralel**: Şifreleme ve şifre çözme işlemleri paralelleştirilebilir
// - **Düşük Overhead**: Minimal bellek kullanımı
// - **Ölçeklenebilir**: Büyük veri setleri için uygundur
//
// # Avantajlar
//
// - Hem gizlilik hem de bütünlük koruması sağlar
// - Endüstri standardı ve yaygın olarak kullanılır
// - Donanım hızlandırması desteği (AES-NI)
// - Basit ve kullanımı kolay API
// - Thread-safe ve stateless tasarım
//
// # Dezavantajlar
//
// - Anahtar yönetimi kullanıcının sorumluluğundadır
// - Nonce tekrarı kritik güvenlik açığına yol açar (ancak bu implementasyonda otomatik yönetilir)
// - Şifreli metin, düz metinden %33 daha büyüktür (Base64 overhead)
//
// # Örnek Kullanım
//
// ```go
// // Güvenli anahtar oluştur (AES-256 için 32 byte)
// key := make([]byte, 32)
// if _, err := rand.Read(key); err != nil {
//     log.Fatal(err)
// }
//
// // CryptGCM instance'ı oluştur
// crypt := encrypt.NewCryptGCM(key)
//
// // Veri şifrele
// plaintext := "Hassas kullanıcı verisi"
// encrypted, err := crypt.Encrypt(plaintext)
// if err != nil {
//     log.Fatal(err)
// }
// fmt.Println("Şifreli:", encrypted)
//
// // Veriyi çöz
// decrypted, err := crypt.Decrypt(encrypted)
// if err != nil {
//     log.Fatal(err)
// }
// fmt.Println("Çözülmüş:", decrypted)
// ```
//
// # Veritabanı Entegrasyonu Örneği
//
// ```go
// type User struct {
//     ID       uint
//     Email    string
//     Password string // Şifrelenmiş olarak saklanacak
// }
//
// func (u *User) SetPassword(password string, crypt *CryptGCM) error {
//     encrypted, err := crypt.Encrypt(password)
//     if err != nil {
//         return err
//     }
//     u.Password = encrypted
//     return nil
// }
//
// func (u *User) GetPassword(crypt *CryptGCM) (string, error) {
//     return crypt.Decrypt(u.Password)
// }
// ```
//
// # Önemli Güvenlik Notları
//
// - **Anahtar Güvenliği**: Anahtarı asla kaynak kodunda saklamayın
// - **Anahtar Rotasyonu**: Düzenli olarak anahtarları değiştirin
// - **Anahtar Boyutu**: 32 byte (AES-256) kullanın, daha güvenlidir
// - **Anahtar Saklama**: Anahtarları güvenli bir key management sisteminde saklayın
// - **Nonce Tekrarı**: Bu implementasyon otomatik nonce yönetir, manuel müdahale gerektirmez
// - **Veri Bütünlüğü**: Şifre çözme başarısız olursa, veri değiştirilmiş demektir
//
// # Best Practices
//
// 1. Anahtarları environment variable'lardan veya secret management sistemlerinden okuyun
// 2. Aynı CryptGCM instance'ını birden fazla işlem için yeniden kullanın
// 3. Hata durumlarını loglamayın (hassas bilgi sızıntısı riski)
// 4. Şifreli verileri veritabanında TEXT veya VARCHAR olarak saklayın
// 5. Anahtar rotasyonu için migration stratejisi planlayın
//
// # Teknik Detaylar
//
// - **Nonce Boyutu**: 12 byte (GCM standardı)
// - **Authentication Tag**: 16 byte (GCM standardı)
// - **Anahtar Boyutları**: 16, 24 veya 32 byte (AES-128, AES-192, AES-256)
// - **Encoding**: Base64 Standard Encoding
// - **Format**: [nonce][ciphertext][auth_tag] -> Base64
type CryptGCM struct {
	key []byte // Şifreleme anahtarı (16, 24 veya 32 byte)
}

// NewCryptGCM, bu fonksiyon yeni bir GCM tabanlı şifreleme instance'ı oluşturur.
//
// Bu constructor fonksiyon, verilen şifreleme anahtarı ile yeni bir CryptGCM
// yapısı oluşturur ve döndürür. Oluşturulan instance, thread-safe olup
// birden fazla goroutine'den güvenle kullanılabilir.
//
// # Parametreler
//
// - `key`: Şifreleme anahtarı (byte slice)
//   - **16 byte**: AES-128 (hızlı, orta güvenlik)
//   - **24 byte**: AES-192 (dengeli)
//   - **32 byte**: AES-256 (en güvenli, önerilen)
//
// # Dönüş Değeri
//
// - `*CryptGCM`: Yapılandırılmış CryptGCM instance'ının pointer'ı
//
// # Anahtar Gereksinimleri
//
// - Anahtar boyutu 16, 24 veya 32 byte olmalıdır
// - Kriptografik olarak güvenli rastgele sayı üreteci ile oluşturulmalıdır
// - Tahmin edilemez olmalıdır (dictionary attack'lara karşı)
// - Güvenli bir şekilde saklanmalıdır (environment variable, secret manager, vb.)
//
// # Güvenlik Önerileri
//
// 1. **Anahtar Üretimi**: crypto/rand kullanarak üretin
// 2. **Anahtar Boyutu**: 32 byte (AES-256) kullanın
// 3. **Anahtar Saklama**: Kaynak kodunda saklamayın
// 4. **Anahtar Rotasyonu**: Düzenli olarak değiştirin
// 5. **Anahtar Paylaşımı**: Güvenli kanallar üzerinden paylaşın
//
// # Kullanım Örnekleri
//
// ## Temel Kullanım
//
// ```go
// // Güvenli anahtar oluştur
// key := make([]byte, 32)
// if _, err := rand.Read(key); err != nil {
//     log.Fatal(err)
// }
//
// // CryptGCM instance'ı oluştur
// crypt := encrypt.NewCryptGCM(key)
//
// // Şimdi şifreleme/şifre çözme yapabilirsiniz
// encrypted, _ := crypt.Encrypt("hassas veri")
// decrypted, _ := crypt.Decrypt(encrypted)
// ```
//
// ## Environment Variable'dan Anahtar Okuma
//
// ```go
// import (
//     "encoding/hex"
//     "os"
// )
//
// // Hex formatında anahtar oku
// keyHex := os.Getenv("ENCRYPTION_KEY")
// key, err := hex.DecodeString(keyHex)
// if err != nil {
//     log.Fatal("Geçersiz anahtar formatı:", err)
// }
//
// crypt := encrypt.NewCryptGCM(key)
// ```
//
// ## Singleton Pattern ile Kullanım
//
// ```go
// var (
//     cryptInstance *encrypt.CryptGCM
//     once          sync.Once
// )
//
// func GetCrypt() *encrypt.CryptGCM {
//     once.Do(func() {
//         key := loadKeyFromSecureStorage()
//         cryptInstance = encrypt.NewCryptGCM(key)
//     })
//     return cryptInstance
// }
// ```
//
// ## Dependency Injection ile Kullanım
//
// ```go
// type UserService struct {
//     crypt *encrypt.CryptGCM
//     db    *gorm.DB
// }
//
// func NewUserService(key []byte, db *gorm.DB) *UserService {
//     return &UserService{
//         crypt: encrypt.NewCryptGCM(key),
//         db:    db,
//     }
// }
//
// func (s *UserService) SaveUser(user *User) error {
//     encrypted, err := s.crypt.Encrypt(user.SensitiveData)
//     if err != nil {
//         return err
//     }
//     user.EncryptedData = encrypted
//     return s.db.Create(user).Error
// }
// ```
//
// # Performans Notları
//
// - Constructor çağrısı hafif bir işlemdir (sadece struct oluşturur)
// - Anahtar doğrulaması yapılmaz (ilk Encrypt/Decrypt çağrısında yapılır)
// - Instance'ı yeniden kullanmak performans açısından önerilir
// - Thread-safe olduğu için global instance kullanılabilir
//
// # Best Practices
//
// 1. **Tek Instance**: Aynı anahtar için tek bir instance oluşturun ve yeniden kullanın
// 2. **Anahtar Yönetimi**: Anahtarı güvenli bir key management sisteminde saklayın
// 3. **Error Handling**: Anahtar yükleme hatalarını uygun şekilde yönetin
// 4. **Logging**: Anahtar değerini asla loglama
// 5. **Testing**: Test ortamında farklı anahtar kullanın
//
// # Önemli Uyarılar
//
// - Anahtar değerini asla kaynak kodunda hardcode etmeyin
// - Anahtar değerini asla version control sistemine commit etmeyin
// - Anahtar değerini asla loglara yazdırmayın
// - Anahtar değerini asla hata mesajlarında göstermeyin
// - Üretim ve test ortamları için farklı anahtarlar kullanın
//
// # Anahtar Üretme Örneği
//
// ```bash
// # Komut satırından güvenli anahtar üretme
// openssl rand -hex 32  # 32 byte hex string
// openssl rand -base64 32  # 32 byte base64 string
// ```
//
// ```go
// // Go ile güvenli anahtar üretme
// func GenerateKey() ([]byte, error) {
//     key := make([]byte, 32)
//     if _, err := rand.Read(key); err != nil {
//         return nil, err
//     }
//     return key, nil
// }
// ```
//
// # Thread Safety
//
// Bu fonksiyon ve döndürdüğü CryptGCM instance'ı tamamen thread-safe'tir.
// Aynı instance'ı birden fazla goroutine'den eşzamanlı olarak kullanabilirsiniz.
//
// # Döndürür
//
// - Yapılandırılmış CryptGCM instance'ının pointer'ı
func NewCryptGCM(key []byte) *CryptGCM {
	return &CryptGCM{key: key}
}

// Encrypt, bu fonksiyon düz metni AES-GCM kullanarak şifreler.
//
// Bu method, CryptGCM instance'ının sakladığı anahtar ile verilen düz metni
// şifreler. Her şifreleme işlemi için benzersiz bir nonce otomatik olarak
// üretilir ve şifreli metnin başına eklenir. Sonuç Base64 ile kodlanarak
// metin formatında döndürülür.
//
// # Parametreler
//
// - `plaintext`: Şifrelenecek düz metin (string formatında)
//   - Boş string kabul edilir
//   - UTF-8 karakterler desteklenir
//   - Boyut sınırı yoktur (bellek izin verdiği sürece)
//
// # Dönüş Değerleri
//
// - `string`: Base64 kodlanmış şifreli metin
//   - Format: Base64([nonce][ciphertext][auth_tag])
//   - Veritabanında TEXT veya VARCHAR olarak saklanabilir
//   - URL-safe değildir (gerekirse URLEncoding.EncodeToString kullanın)
// - `error`: İşlem başarısız olursa hata mesajı
//   - Cipher oluşturma hatası
//   - GCM modu oluşturma hatası
//   - Nonce üretme hatası
//
// # Çalışma Prensibi
//
// 1. Instance'ın anahtarı ile AES cipher bloğu oluşturulur
// 2. GCM (Galois/Counter Mode) aktif edilir
// 3. Kriptografik olarak güvenli 12 byte'lık nonce üretilir
// 4. Düz metin şifrelenir ve 16 byte'lık auth tag eklenir
// 5. [nonce + ciphertext + auth_tag] Base64 ile kodlanır
//
// # Güvenlik Özellikleri
//
// - **Unique Nonce**: Her çağrıda benzersiz nonce üretilir
// - **Authentication**: Kimlik doğrulama etiketi otomatik eklenir
// - **Tamper-Proof**: Veri değişikliği decrypt sırasında tespit edilir
// - **Thread-Safe**: Birden fazla goroutine'den güvenle çağrılabilir
//
// # Performans Karakteristikleri
//
// - **Hız**: ~1-2 GB/s (AES-NI ile)
// - **Overhead**: Düz metin boyutuna +28 byte (nonce:12 + tag:16) + Base64 overhead (%33)
// - **Bellek**: Minimal allocation, düz metin boyutuna orantılı
// - **Paralellik**: Thread-safe, concurrent kullanım desteklenir
//
// # Kullanım Örnekleri
//
// ## Temel Kullanım
//
// ```go
// crypt := encrypt.NewCryptGCM(key)
// encrypted, err := crypt.Encrypt("hassas veri")
// if err != nil {
//     log.Fatal(err)
// }
// fmt.Println(encrypted) // Base64 string
// ```
//
// ## Veritabanı Entegrasyonu
//
// ```go
// type User struct {
//     ID       uint
//     Email    string
//     CreditCard string `gorm:"type:text"` // Şifrelenmiş olarak saklanacak
// }
//
// func (u *User) BeforeSave(tx *gorm.DB) error {
//     crypt := getCryptInstance() // Singleton pattern
//     encrypted, err := crypt.Encrypt(u.CreditCard)
//     if err != nil {
//         return err
//     }
//     u.CreditCard = encrypted
//     return nil
// }
// ```
//
// ## API Token Şifreleme
//
// ```go
// type APIToken struct {
//     UserID    uint
//     Token     string // Şifrelenmiş
//     ExpiresAt time.Time
// }
//
// func CreateToken(userID uint, crypt *encrypt.CryptGCM) (*APIToken, error) {
//     rawToken := generateRandomToken()
//     encrypted, err := crypt.Encrypt(rawToken)
//     if err != nil {
//         return nil, err
//     }
//     return &APIToken{
//         UserID: userID,
//         Token:  encrypted,
//         ExpiresAt: time.Now().Add(24 * time.Hour),
//     }, nil
// }
// ```
//
// ## Batch Şifreleme
//
// ```go
// func EncryptBatch(items []string, crypt *encrypt.CryptGCM) ([]string, error) {
//     encrypted := make([]string, len(items))
//     for i, item := range items {
//         enc, err := crypt.Encrypt(item)
//         if err != nil {
//             return nil, fmt.Errorf("item %d: %w", i, err)
//         }
//         encrypted[i] = enc
//     }
//     return encrypted, nil
// }
// ```
//
// ## Concurrent Şifreleme
//
// ```go
// func EncryptConcurrent(items []string, crypt *encrypt.CryptGCM) ([]string, error) {
//     var wg sync.WaitGroup
//     encrypted := make([]string, len(items))
//     errChan := make(chan error, len(items))
//
//     for i, item := range items {
//         wg.Add(1)
//         go func(idx int, data string) {
//             defer wg.Done()
//             enc, err := crypt.Encrypt(data)
//             if err != nil {
//                 errChan <- err
//                 return
//             }
//             encrypted[idx] = enc
//         }(i, item)
//     }
//
//     wg.Wait()
//     close(errChan)
//
//     if err := <-errChan; err != nil {
//         return nil, err
//     }
//     return encrypted, nil
// }
// ```
//
// ## JSON Şifreleme
//
// ```go
// func EncryptJSON(data interface{}, crypt *encrypt.CryptGCM) (string, error) {
//     jsonBytes, err := json.Marshal(data)
//     if err != nil {
//         return "", err
//     }
//     return crypt.Encrypt(string(jsonBytes))
// }
// ```
//
// # Hata Yönetimi
//
// ```go
// encrypted, err := crypt.Encrypt(plaintext)
// if err != nil {
//     // Hata türüne göre işlem yap
//     if strings.Contains(err.Error(), "cipher") {
//         log.Error("Geçersiz anahtar boyutu")
//     } else if strings.Contains(err.Error(), "nonce") {
//         log.Error("Rastgele sayı üretimi başarısız")
//     }
//     return err
// }
// ```
//
// # Best Practices
//
// 1. **Aynı Instance**: Aynı CryptGCM instance'ını yeniden kullanın
// 2. **Error Handling**: Hataları uygun şekilde yönetin
// 3. **Logging**: Düz metni asla loglama
// 4. **Validation**: Şifrelemeden önce input validation yapın
// 5. **Storage**: Şifreli metni TEXT/VARCHAR olarak saklayın
//
// # Önemli Notlar
//
// - Her çağrıda farklı şifreli metin üretilir (farklı nonce)
// - Aynı düz metin bile farklı şifreli metinler üretir
// - Şifreli metin boyutu: len(plaintext) + 28 byte + Base64 overhead
// - Thread-safe, concurrent kullanım güvenlidir
// - Nonce collision riski yok (kriptografik RNG kullanılır)
//
// # Güvenlik Uyarıları
//
// - Düz metni asla loglama
// - Şifreleme hatalarını kullanıcıya detaylı göstermeyin
// - Anahtar rotasyonu için plan yapın
// - Şifreli verileri düzenli olarak yedekleyin
//
// # Performans İpuçları
//
// - Büyük veriler için streaming encryption düşünün
// - Batch işlemler için goroutine kullanın
// - Instance'ı cache'leyin, her seferinde yeniden oluşturmayın
// - Veritabanı indexleme için şifreli alanları kullanmayın
func (c *CryptGCM) Encrypt(plaintext string) (string, error) {
	return encryptGCM(plaintext, c.key)
}

// Decrypt, bu fonksiyon şifreli metni AES-GCM kullanarak çözer.
//
// Bu method, CryptGCM instance'ının sakladığı anahtar ile verilen şifreli metni
// çözer ve kimlik doğrulaması yapar. Şifreli metin Base64 formatında olmalı ve
// Encrypt method'u ile oluşturulmuş olmalıdır. Veri değiştirilmişse veya yanlış
// anahtar kullanılmışsa, fonksiyon hata döndürür.
//
// # Parametreler
//
// - `ciphertext`: Base64 kodlanmış şifreli metin
//   - Encrypt method'u ile oluşturulmuş olmalı
//   - Format: Base64([nonce][ciphertext][auth_tag])
//   - Geçerli Base64 string olmalı
//
// # Dönüş Değerleri
//
// - `string`: Çözülmüş düz metin (orijinal veri)
//   - UTF-8 string formatında
//   - Şifrelemeden önceki orijinal veri
// - `error`: İşlem başarısız olursa hata mesajı
//   - Base64 çözme hatası
//   - Geçersiz anahtar
//   - Veri değiştirilmiş (authentication failure)
//   - Şifreli metin çok kısa
//   - Cipher oluşturma hatası
//
// # Çalışma Prensibi
//
// 1. Base64 kodlanmış şifreli metin çözülür
// 2. Instance'ın anahtarı ile AES cipher bloğu oluşturulur
// 3. GCM (Galois/Counter Mode) aktif edilir
// 4. Nonce ve şifreli metin + auth tag ayrıştırılır
// 5. Kimlik doğrulama etiketi kontrol edilir
// 6. Şifreli metin çözülür ve düz metin döndürülür
//
// # Güvenlik Özellikleri
//
// - **Authentication Verification**: Kimlik doğrulama etiketi otomatik kontrol edilir
// - **Tamper Detection**: Veri değiştirilmişse hata döndürülür
// - **Integrity Check**: Veri bütünlüğü garanti edilir
// - **Thread-Safe**: Birden fazla goroutine'den güvenle çağrılabilir
// - **Timing-Safe**: Timing attack'lara karşı korumalıdır
//
// # Hata Durumları ve Anlamları
//
// - **"base64 çözülemedi"**: Geçersiz Base64 formatı
// - **"cipher oluşturulamadı"**: Geçersiz anahtar boyutu (16, 24, 32 byte olmalı)
// - **"GCM modu oluşturulamadı"**: Cipher bloğu hatası
// - **"şifreli metin çok kısa"**: Veri bozuk veya eksik
// - **"şifre çözme veya doğrulama başarısız"**: Yanlış anahtar veya veri değiştirilmiş
//
// # Performans Karakteristikleri
//
// - **Hız**: ~1-2 GB/s (AES-NI ile)
// - **Bellek**: Minimal allocation, şifreli metin boyutuna orantılı
// - **Paralellik**: Thread-safe, concurrent kullanım desteklenir
// - **Overhead**: Kimlik doğrulama kontrolü minimal overhead ekler
//
// # Kullanım Örnekleri
//
// ## Temel Kullanım
//
// ```go
// crypt := encrypt.NewCryptGCM(key)
// encrypted := "base64-encoded-ciphertext..."
// decrypted, err := crypt.Decrypt(encrypted)
// if err != nil {
//     log.Fatal(err)
// }
// fmt.Println(decrypted) // Orijinal düz metin
// ```
//
// ## Veritabanından Okuma
//
// ```go
// type User struct {
//     ID         uint
//     Email      string
//     CreditCard string `gorm:"type:text"` // Şifrelenmiş
// }
//
// func (u *User) AfterFind(tx *gorm.DB) error {
//     crypt := getCryptInstance()
//     decrypted, err := crypt.Decrypt(u.CreditCard)
//     if err != nil {
//         return fmt.Errorf("şifre çözme hatası: %w", err)
//     }
//     u.CreditCard = decrypted
//     return nil
// }
// ```
//
// ## Hata Yönetimi ile Kullanım
//
// ```go
// func GetSensitiveData(encrypted string, crypt *encrypt.CryptGCM) (string, error) {
//     decrypted, err := crypt.Decrypt(encrypted)
//     if err != nil {
//         // Hata türüne göre işlem yap
//         if strings.Contains(err.Error(), "doğrulama başarısız") {
//             log.Warn("Veri değiştirilmiş veya yanlış anahtar")
//             return "", fmt.Errorf("güvenlik hatası: veri bütünlüğü bozulmuş")
//         }
//         if strings.Contains(err.Error(), "base64") {
//             log.Warn("Geçersiz veri formatı")
//             return "", fmt.Errorf("veri formatı hatası")
//         }
//         return "", err
//     }
//     return decrypted, nil
// }
// ```
//
// ## Batch Şifre Çözme
//
// ```go
// func DecryptBatch(items []string, crypt *encrypt.CryptGCM) ([]string, error) {
//     decrypted := make([]string, 0, len(items))
//     for i, item := range items {
//         dec, err := crypt.Decrypt(item)
//         if err != nil {
//             log.Warnf("Item %d şifre çözme hatası: %v", i, err)
//             continue // Veya return nil, err
//         }
//         decrypted = append(decrypted, dec)
//     }
//     return decrypted, nil
// }
// ```
//
// ## Concurrent Şifre Çözme
//
// ```go
// func DecryptConcurrent(items []string, crypt *encrypt.CryptGCM) ([]string, error) {
//     var wg sync.WaitGroup
//     decrypted := make([]string, len(items))
//     errChan := make(chan error, len(items))
//
//     for i, item := range items {
//         wg.Add(1)
//         go func(idx int, data string) {
//             defer wg.Done()
//             dec, err := crypt.Decrypt(data)
//             if err != nil {
//                 errChan <- fmt.Errorf("item %d: %w", idx, err)
//                 return
//             }
//             decrypted[idx] = dec
//         }(i, item)
//     }
//
//     wg.Wait()
//     close(errChan)
//
//     if err := <-errChan; err != nil {
//         return nil, err
//     }
//     return decrypted, nil
// }
// ```
//
// ## JSON Şifre Çözme
//
// ```go
// func DecryptJSON(encrypted string, crypt *encrypt.CryptGCM, v interface{}) error {
//     decrypted, err := crypt.Decrypt(encrypted)
//     if err != nil {
//         return err
//     }
//     return json.Unmarshal([]byte(decrypted), v)
// }
//
// // Kullanım
// var user User
// err := DecryptJSON(encryptedJSON, crypt, &user)
// ```
//
// ## API Response Şifre Çözme
//
// ```go
// type APIResponse struct {
//     Data      string `json:"data"` // Şifrelenmiş
//     Timestamp int64  `json:"timestamp"`
// }
//
// func ProcessAPIResponse(resp *APIResponse, crypt *encrypt.CryptGCM) (string, error) {
//     decrypted, err := crypt.Decrypt(resp.Data)
//     if err != nil {
//         return "", fmt.Errorf("API response şifre çözme hatası: %w", err)
//     }
//     return decrypted, nil
// }
// ```
//
// ## Fallback Mekanizması
//
// ```go
// func DecryptWithFallback(encrypted string, primaryKey, backupKey []byte) (string, error) {
//     // Önce primary key ile dene
//     primaryCrypt := encrypt.NewCryptGCM(primaryKey)
//     decrypted, err := primaryCrypt.Decrypt(encrypted)
//     if err == nil {
//         return decrypted, nil
//     }
//
//     // Başarısız olursa backup key ile dene (key rotation için)
//     backupCrypt := encrypt.NewCryptGCM(backupKey)
//     decrypted, err = backupCrypt.Decrypt(encrypted)
//     if err != nil {
//         return "", fmt.Errorf("her iki anahtar ile de şifre çözme başarısız")
//     }
//
//     // Backup key ile başarılı olduysa, yeni key ile yeniden şifrele
//     reencrypted, _ := primaryCrypt.Encrypt(decrypted)
//     // reencrypted'ı veritabanına kaydet
//
//     return decrypted, nil
// }
// ```
//
// # Best Practices
//
// 1. **Error Handling**: Hataları uygun şekilde yönetin ve loglayın
// 2. **Validation**: Şifre çözmeden önce input validation yapın
// 3. **Logging**: Hassas verileri asla loglama
// 4. **Retry Logic**: Geçici hatalar için retry mekanizması ekleyin
// 5. **Monitoring**: Şifre çözme hatalarını izleyin (olası saldırı göstergesi)
//
// # Önemli Notlar
//
// - Yanlış anahtar kullanılırsa şifre çözme başarısız olur
// - Veri değiştirilmişse kimlik doğrulama başarısız olur
// - Hata mesajları hassas bilgi içermez (güvenlik)
// - Thread-safe, concurrent kullanım güvenlidir
// - Aynı instance'ı yeniden kullanmak performans açısından önerilir
//
// # Güvenlik Uyarıları
//
// - **Hata Mesajları**: Kullanıcıya detaylı hata mesajı göstermeyin
// - **Logging**: Şifre çözülmüş veriyi asla loglama
// - **Timing Attacks**: Bu implementasyon timing-safe'tir
// - **Key Rotation**: Düzenli anahtar rotasyonu yapın
// - **Monitoring**: Başarısız şifre çözme denemelerini izleyin
//
// # Performans İpuçları
//
// - Instance'ı cache'leyin, her seferinde yeniden oluşturmayın
// - Büyük veriler için streaming decryption düşünün
// - Batch işlemler için goroutine kullanın
// - Hata durumlarında early return yapın
//
// # Debugging İpuçları
//
// ```go
// // Şifre çözme hatalarını debug etme
// decrypted, err := crypt.Decrypt(encrypted)
// if err != nil {
//     log.Debugf("Şifreli metin uzunluğu: %d", len(encrypted))
//     log.Debugf("Anahtar uzunluğu: %d", len(crypt.key))
//     log.Debugf("Hata: %v", err)
//     // Hassas bilgileri loglama!
// }
// ```
//
// # Yaygın Hatalar ve Çözümleri
//
// 1. **"base64 çözülemedi"**
//    - Çözüm: Veriyi Encrypt ile şifreleyin, manuel oluşturmayın
//
// 2. **"şifreli metin çok kısa"**
//    - Çözüm: Veri bozulmuş, kaynak veriyi kontrol edin
//
// 3. **"doğrulama başarısız"**
//    - Çözüm: Doğru anahtarı kullandığınızdan emin olun
//    - Çözüm: Veri değiştirilmemiş olduğunu kontrol edin
//
// 4. **"cipher oluşturulamadı"**
//    - Çözüm: Anahtar boyutunu kontrol edin (16, 24 veya 32 byte)
func (c *CryptGCM) Decrypt(ciphertext string) (string, error) {
	return decryptGCM(ciphertext, c.key)
}
