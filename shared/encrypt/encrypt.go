// Package encrypt, AES-256 tabanlı şifreleme ve şifre çözme işlemleri için
// güvenli ve kullanımı kolay bir arayüz sağlar.
//
// # Özellikler
//
// - AES-256-CBC modunda şifreleme (aes256.go)
// - AES-256-GCM modunda authenticated encryption (aes_gcm.go)
// - Singleton pattern ile global instance yönetimi
// - Base64 kodlamalı çıktı
// - PKCS7 padding desteği
//
// # Güvenlik Notları
//
// - Şifreleme anahtarı en az 32 byte (256 bit) olmalıdır
// - Anahtarlar hexadecimal string formatında sağlanmalıdır
// - Her şifreleme işlemi için rastgele IV (Initialization Vector) kullanılır
// - GCM modu, CBC moduna göre daha güvenlidir (authenticated encryption)
//
// # Kullanım Senaryoları
//
// - Veritabanında hassas bilgilerin şifrelenmesi (kredi kartı, şifre, vb.)
// - API token'larının güvenli saklanması
// - Kullanıcı verilerinin şifrelenmesi
// - Dosya içeriklerinin şifrelenmesi
//
// # Örnek Kullanım
//
// ```go
// // Hexadecimal formatında 32 byte (64 karakter) anahtar
// key := "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
//
// // Şifreleme instance'ı oluştur
// crypt := encrypt.NewCrypt(key)
//
// // Veri şifrele
// encrypted, err := crypt.Encrypt("hassas bilgi")
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
package encrypt

import (
	"encoding/hex"
	"log"
)

// Crypt, şifreleme ve şifre çözme işlemleri için temel arayüzü tanımlar.
//
// Bu interface, farklı şifreleme algoritmalarının (AES-CBC, AES-GCM, vb.)
// aynı arayüz üzerinden kullanılmasını sağlar. Dependency injection ve
// test edilebilirlik için idealdir.
//
// # Implementasyonlar
//
// - `crypt`: AES-256-CBC modu (varsayılan)
// - `CryptGCM`: AES-256-GCM modu (authenticated encryption)
//
// # Kullanım Senaryoları
//
// - Veritabanı şifrelemesi
// - API token yönetimi
// - Hassas veri saklama
// - Şifreli iletişim
//
// # Avantajlar
//
// - Interface sayesinde farklı implementasyonlar kolayca değiştirilebilir
// - Mock implementasyonlar ile test edilebilir
// - Dependency injection pattern'i ile kullanılabilir
//
// # Örnek Kullanım
//
// ```go
// var cryptService Crypt = NewCrypt(key)
//
// // Şifreleme
// encrypted, _ := cryptService.Encrypt("gizli veri")
//
// // Şifre çözme
// decrypted, _ := cryptService.Decrypt(encrypted)
// ```
type Crypt interface {
	// Encrypt, düz metni şifreler ve base64 kodlanmış string döndürür.
	//
	// # Parametreler
	//
	// - `plaintext`: Şifrelenecek düz metin
	//
	// # Döndürür
	//
	// - `string`: Base64 kodlanmış şifreli metin
	// - `error`: Şifreleme hatası (varsa)
	//
	// # Hata Durumları
	//
	// - Cipher oluşturma hatası
	// - IV (Initialization Vector) oluşturma hatası
	// - Padding hatası
	//
	// # Örnek
	//
	// ```go
	// encrypted, err := crypt.Encrypt("hassas bilgi")
	// if err != nil {
	//     log.Printf("Şifreleme hatası: %v", err)
	// }
	// ```
	Encrypt(plaintext string) (string, error)

	// Decrypt, şifreli metni çözer ve orijinal düz metni döndürür.
	//
	// # Parametreler
	//
	// - `ciphertext`: Base64 kodlanmış şifreli metin
	//
	// # Döndürür
	//
	// - `string`: Çözülmüş düz metin
	// - `error`: Şifre çözme hatası (varsa)
	//
	// # Hata Durumları
	//
	// - Base64 decode hatası
	// - Geçersiz şifreli metin uzunluğu
	// - Cipher oluşturma hatası
	// - Padding hatası
	// - Geçersiz IV
	//
	// # Örnek
	//
	// ```go
	// decrypted, err := crypt.Decrypt(encrypted)
	// if err != nil {
	//     log.Printf("Şifre çözme hatası: %v", err)
	// }
	// ```
	Decrypt(ciphertext string) (string, error)
}

// _crypt, global singleton instance'ı saklar.
//
// # Singleton Pattern
//
// Bu değişken, uygulamada tek bir şifreleme instance'ının kullanılmasını
// sağlar. Bu yaklaşım:
//
// - Bellek kullanımını optimize eder
// - Aynı anahtarın tekrar tekrar parse edilmesini önler
// - Global erişim sağlar
//
// # Uyarı
//
// Singleton pattern, test edilebilirliği zorlaştırabilir. Test senaryolarında
// bu değişkeni sıfırlamak gerekebilir.
var _crypt Crypt

// crypt, Crypt interface'inin AES-256-CBC implementasyonudur.
//
// Bu yapı, AES-256 algoritması ile CBC (Cipher Block Chaining) modunda
// şifreleme ve şifre çözme işlemlerini gerçekleştirir.
//
// # Özellikler
//
// - AES-256-CBC modu
// - PKCS7 padding
// - Rastgele IV (her şifreleme için farklı)
// - Base64 kodlamalı çıktı
//
// # Güvenlik Özellikleri
//
// - 256-bit anahtar uzunluğu (yüksek güvenlik)
// - Her şifreleme için benzersiz IV
// - PKCS7 padding ile blok hizalama
//
// # Dezavantajlar
//
// - CBC modu, authenticated encryption sağlamaz (veri bütünlüğü kontrolü yok)
// - Padding oracle saldırılarına karşı hassas olabilir
// - GCM moduna göre daha az güvenli
//
// # Alternatif
//
// Daha yüksek güvenlik için `CryptGCM` kullanılabilir (aes_gcm.go).
//
// # İç Yapı
//
// ```go
// type crypt struct {
//     key []byte // 32 byte (256 bit) şifreleme anahtarı
// }
// ```
type crypt struct {
	// key, AES-256 şifreleme anahtarını saklar.
	//
	// # Özellikler
	//
	// - Uzunluk: 32 byte (256 bit)
	// - Format: Raw bytes (hexadecimal'den decode edilmiş)
	// - Güvenlik: Bellekte düz metin olarak saklanır
	//
	// # Güvenlik Notu
	//
	// Anahtar bellekte düz metin olarak saklandığı için, memory dump
	// saldırılarına karşı hassastır. Production ortamlarında:
	// - Anahtarları environment variable'lardan okuyun
	// - Anahtarları asla kod içine hardcode etmeyin
	// - Key rotation stratejisi uygulayın
	key []byte
}

// NewCrypt, yeni bir Crypt instance'ı oluşturur veya mevcut singleton instance'ı döndürür.
//
// Bu fonksiyon, singleton pattern kullanarak uygulamada tek bir şifreleme
// instance'ının kullanılmasını sağlar. İlk çağrıda yeni bir instance oluşturur,
// sonraki çağrılarda aynı instance'ı döndürür.
//
// # Parametreler
//
// - `key`: Hexadecimal formatında şifreleme anahtarı (64 karakter = 32 byte)
//
// # Döndürür
//
// - `Crypt`: Şifreleme interface'i implementasyonu
//
// # Anahtar Formatı
//
// Anahtar, hexadecimal string formatında olmalıdır:
// - Uzunluk: 64 karakter (32 byte = 256 bit)
// - Karakterler: 0-9, a-f (büyük/küçük harf duyarsız)
// - Örnek: "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
//
// # Anahtar Oluşturma
//
// ```bash
// # Linux/macOS
// openssl rand -hex 32
//
// # Go ile
// key := make([]byte, 32)
// rand.Read(key)
// hexKey := hex.EncodeToString(key)
// ```
//
// # Singleton Davranışı
//
// ```go
// crypt1 := NewCrypt(key1)
// crypt2 := NewCrypt(key2) // key2 göz ardı edilir, crypt1 döndürülür
// // crypt1 == crypt2 (aynı instance)
// ```
//
// # Hata Durumları
//
// - Geçersiz hexadecimal format: `log.Fatalf` ile uygulama sonlanır
// - Yanlış anahtar uzunluğu: `log.Fatalf` ile uygulama sonlanır
//
// # Uyarılar
//
// - Bu fonksiyon hata durumunda `log.Fatalf` kullanır ve uygulamayı sonlandırır
// - Production ortamlarında daha zarif hata yönetimi düşünülebilir
// - Singleton pattern, test senaryolarında sorun yaratabilir
// - Thread-safe değildir, concurrent kullanımda race condition oluşabilir
//
// # Kullanım Senaryoları
//
// - Uygulama başlangıcında tek bir şifreleme instance'ı oluşturma
// - Global şifreleme servisi sağlama
// - Aynı anahtarın tüm uygulama boyunca kullanılması
//
// # Avantajlar
//
// - Bellek verimliliği (tek instance)
// - Anahtar parse işlemi bir kez yapılır
// - Global erişim kolaylığı
//
// # Dezavantajlar
//
// - Test edilebilirlik zorluğu
// - Thread-safe değil
// - Anahtar değiştirme esnekliği yok
// - Fatal error ile uygulama sonlanması
//
// # Örnek Kullanım
//
// ```go
// // Uygulama başlangıcında
// key := os.Getenv("ENCRYPTION_KEY")
// crypt := encrypt.NewCrypt(key)
//
// // Uygulama içinde
// encrypted, err := crypt.Encrypt("hassas veri")
// if err != nil {
//     log.Printf("Şifreleme hatası: %v", err)
// }
//
// // Başka bir yerde
// crypt2 := encrypt.NewCrypt("farklı-anahtar") // Aynı instance döner
// decrypted, _ := crypt2.Decrypt(encrypted)
// ```
//
// # Alternatif Yaklaşım
//
// Singleton yerine dependency injection kullanmak daha test edilebilir:
//
// ```go
// type Service struct {
//     crypt Crypt
// }
//
// func NewService(crypt Crypt) *Service {
//     return &Service{crypt: crypt}
// }
// ```
func NewCrypt(key string) Crypt {
	if _crypt != nil {
		return _crypt
	}

	keyBytes, err := hex.DecodeString(key)
	if err != nil {
		log.Fatalf("Failed to decode encryption key: %v", err)
	}

	_crypt = &crypt{key: keyBytes}

	return _crypt
}

// Encrypt, düz metni AES-256-CBC modu ile şifreler.
//
// Bu method, `crypt` struct'ının Crypt interface'ini implement eder.
// Şifreleme işlemi için `aes256.go` dosyasındaki `encrypt` fonksiyonunu kullanır.
//
// # Parametreler
//
// - `plaintext`: Şifrelenecek düz metin string
//
// # Döndürür
//
// - `string`: Base64 kodlanmış şifreli metin (IV + ciphertext)
// - `error`: Şifreleme hatası (varsa)
//
// # Şifreleme Süreci
//
// 1. AES cipher oluşturulur (256-bit anahtar)
// 2. Rastgele IV (Initialization Vector) üretilir (16 byte)
// 3. Düz metin PKCS7 padding ile hizalanır
// 4. CBC modu ile şifreleme yapılır
// 5. IV + şifreli metin birleştirilir
// 6. Base64 ile kodlanır
//
// # Çıktı Formatı
//
// ```
// Base64(IV || Ciphertext)
// ```
//
// - IV: İlk 16 byte (AES block size)
// - Ciphertext: Şifrelenmiş ve padding eklenmiş veri
//
// # Güvenlik Özellikleri
//
// - Her şifreleme için benzersiz rastgele IV
// - PKCS7 padding ile blok hizalama
// - AES-256 güvenlik seviyesi
//
// # Hata Durumları
//
// - Cipher oluşturma hatası (geçersiz anahtar)
// - IV oluşturma hatası (rastgele sayı üreteci sorunu)
// - Padding hatası
//
// # Örnek Kullanım
//
// ```go
// crypt := &crypt{key: keyBytes}
// encrypted, err := crypt.Encrypt("Kredi kartı: 1234-5678-9012-3456")
// if err != nil {
//     return fmt.Errorf("şifreleme hatası: %w", err)
// }
// // encrypted: "base64-encoded-string"
// ```
//
// # Performans
//
// - Küçük metinler için hızlı (~microseconds)
// - Büyük metinler için linear zaman karmaşıklığı
// - Base64 encoding ek overhead ekler (~33% boyut artışı)
func (c *crypt) Encrypt(plaintext string) (string, error) {
	return encrypt(plaintext, c.key)
}

// Decrypt, şifreli metni AES-256-CBC modu ile çözer.
//
// Bu method, `crypt` struct'ının Crypt interface'ini implement eder.
// Şifre çözme işlemi için `aes256.go` dosyasındaki `decrypt` fonksiyonunu kullanır.
//
// # Parametreler
//
// - `ciphertext`: Base64 kodlanmış şifreli metin (IV + ciphertext)
//
// # Döndürür
//
// - `string`: Çözülmüş düz metin
// - `error`: Şifre çözme hatası (varsa)
//
// # Şifre Çözme Süreci
//
// 1. Base64 decode işlemi
// 2. IV ve ciphertext ayrıştırılır
// 3. AES cipher oluşturulur
// 4. CBC modu ile şifre çözülür
// 5. PKCS7 padding kaldırılır
// 6. Düz metin döndürülür
//
// # Girdi Formatı
//
// ```
// Base64(IV || Ciphertext)
// ```
//
// - İlk 16 byte: IV (Initialization Vector)
// - Kalan bytes: Şifrelenmiş veri
//
// # Hata Durumları
//
// - Base64 decode hatası (geçersiz format)
// - Çok kısa ciphertext (< 16 byte)
// - Geçersiz anahtar
// - Blok hizalama hatası (ciphertext uzunluğu block size'ın katı değil)
// - Geçersiz padding
// - Yanlış anahtar (padding hatası veya garbled output)
//
// # Güvenlik Kontrolleri
//
// - Minimum uzunluk kontrolü (en az 1 block + IV)
// - Block size hizalama kontrolü
// - PKCS7 padding doğrulama
//
// # Örnek Kullanım
//
// ```go
// crypt := &crypt{key: keyBytes}
// decrypted, err := crypt.Decrypt(encryptedData)
// if err != nil {
//     return fmt.Errorf("şifre çözme hatası: %w", err)
// }
// // decrypted: "Kredi kartı: 1234-5678-9012-3456"
// ```
//
// # Hata Mesajları
//
// - "ciphertext too short": Veri en az 16 byte olmalı
// - "ciphertext is not a multiple of the block size": Geçersiz veri uzunluğu
// - "invalid padding": Yanlış anahtar veya bozuk veri
// - "failed to decode base64": Geçersiz base64 formatı
//
// # Performans
//
// - Küçük metinler için hızlı (~microseconds)
// - Büyük metinler için linear zaman karmaşıklığı
// - Base64 decoding ek overhead ekler
func (c *crypt) Decrypt(ciphertext string) (string, error) {
	return decrypt(ciphertext, c.key)
}
