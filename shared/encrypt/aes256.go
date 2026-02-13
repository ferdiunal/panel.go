// Package encrypt, AES-256 şifreleme ve şifre çözme işlemleri için güvenli kriptografik
// fonksiyonlar sağlar. Bu paket, CBC (Cipher Block Chaining) modu ve PKCS7 padding
// kullanarak veri güvenliğini sağlar.
//
// # Güvenlik Özellikleri
//
// - AES-256 şifreleme algoritması (endüstri standardı)
// - CBC (Cipher Block Chaining) modu ile blok şifreleme
// - Rastgele IV (Initialization Vector) üretimi
// - PKCS7 padding standardı
// - Base64 kodlama ile güvenli veri aktarımı
//
// # Kullanım Senaryoları
//
// - Hassas kullanıcı verilerinin şifrelenmesi
// - Veritabanında güvenli veri saklama
// - API token'larının korunması
// - Konfigürasyon dosyalarında şifre saklama
// - Güvenli veri iletimi
//
// # Önemli Notlar
//
// ⚠️ **Anahtar Yönetimi**: Şifreleme anahtarları güvenli bir şekilde saklanmalı ve
// asla kaynak kodda hardcode edilmemelidir.
//
// ⚠️ **Anahtar Uzunluğu**: AES-256 için 32 byte (256 bit) anahtar kullanılmalıdır.
//
// ⚠️ **IV Güvenliği**: Her şifreleme işlemi için yeni bir rastgele IV üretilir ve
// şifreli verinin başına eklenir.
//
// # Avantajlar
//
// - Endüstri standardı güvenlik seviyesi
// - Hızlı şifreleme/şifre çözme performansı
// - Yaygın kabul görmüş algoritma
// - Go standart kütüphanesi desteği
//
// # Dezavantajlar
//
// - CBC modu padding oracle saldırılarına karşı hassas olabilir
// - Paralel şifreleme desteği yok (GCM modu tercih edilebilir)
// - Her şifreleme işlemi için IV saklanması gerekir
package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// encrypt, verilen düz metni AES-256-CBC algoritması ile şifreler.
//
// Bu fonksiyon, güvenli bir şekilde veri şifrelemek için AES-256 algoritmasını
// CBC (Cipher Block Chaining) modu ile kullanır. Her şifreleme işlemi için
// kriptografik olarak güvenli rastgele bir IV (Initialization Vector) üretir.
//
// # Şifreleme Süreci
//
// 1. AES cipher bloğu oluşturulur (256-bit anahtar ile)
// 2. Rastgele 16-byte IV üretilir
// 3. Düz metin PKCS7 padding ile doldurulur
// 4. CBC modu ile şifreleme yapılır
// 5. IV ve şifreli veri birleştirilir
// 6. Sonuç Base64 formatında kodlanır
//
// # Parametreler
//
// - `plaintext`: Şifrelenecek düz metin (string)
// - `key`: AES-256 şifreleme anahtarı (32 byte/256 bit olmalı)
//
// # Dönüş Değerleri
//
// - `string`: Base64 kodlanmış şifreli veri (IV + ciphertext)
// - `error`: Şifreleme sırasında oluşan hata (varsa)
//
// # Kullanım Örneği
//
// ```go
// key := []byte("12345678901234567890123456789012") // 32 byte
// plaintext := "Gizli mesaj"
//
// encrypted, err := encrypt(plaintext, key)
// if err != nil {
//     log.Fatal(err)
// }
// fmt.Println("Şifreli:", encrypted)
// ```
//
// # Hata Durumları
//
// - Geçersiz anahtar uzunluğu (16, 24 veya 32 byte olmalı)
// - IV üretimi başarısız (rastgele sayı üreteci hatası)
//
// # Güvenlik Notları
//
// ⚠️ **Anahtar Güvenliği**: Anahtar güvenli bir şekilde saklanmalı ve paylaşılmamalıdır.
//
// ⚠️ **IV Benzersizliği**: Her şifreleme işlemi için otomatik olarak yeni bir IV üretilir.
//
// ⚠️ **Veri Bütünlüğü**: Bu fonksiyon veri bütünlüğü kontrolü sağlamaz. HMAC veya
// GCM modu kullanımı önerilir.
//
// # Performans
//
// - Küçük veriler için hızlı (< 1ms)
// - Büyük veriler için lineer zaman karmaşıklığı O(n)
// - Bellek kullanımı: ~2x veri boyutu (padding + IV)
func encrypt(plaintext string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// CBC modu için IV (Initialization Vector) gereklidir
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// PKCS7 padding ile veriyi blok boyutuna tamamla
	padded := pkcs7Pad([]byte(plaintext), aes.BlockSize)

	// Şifreli veri için bellek ayır
	ciphertext := make([]byte, len(padded))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, padded)

	// IV'yi şifreli verinin başına ekle (şifre çözme için gerekli)
	final := append(iv, ciphertext...)
	return base64.StdEncoding.EncodeToString(final), nil
}

// decrypt, AES-256-CBC ile şifrelenmiş veriyi çözer ve orijinal düz metni döndürür.
//
// Bu fonksiyon, encrypt fonksiyonu tarafından şifrelenmiş Base64 kodlu veriyi
// çözer. Şifreli verinin başındaki IV'yi otomatik olarak ayıklar ve şifre çözme
// işleminde kullanır.
//
// # Şifre Çözme Süreci
//
// 1. Base64 kodlu veri decode edilir
// 2. İlk 16 byte IV olarak ayrıştırılır
// 3. Kalan kısım şifreli veri olarak işlenir
// 4. AES cipher bloğu oluşturulur
// 5. CBC modu ile şifre çözülür
// 6. PKCS7 padding kaldırılır
// 7. Orijinal düz metin döndürülür
//
// # Parametreler
//
// - `b64cipher`: Base64 kodlanmış şifreli veri (IV + ciphertext)
// - `key`: Şifreleme sırasında kullanılan aynı AES-256 anahtarı (32 byte)
//
// # Dönüş Değerleri
//
// - `string`: Çözülmüş düz metin
// - `error`: Şifre çözme sırasında oluşan hata (varsa)
//
// # Kullanım Örneği
//
// ```go
// key := []byte("12345678901234567890123456789012") // 32 byte
// encrypted := "base64_encoded_encrypted_data..."
//
// decrypted, err := decrypt(encrypted, key)
// if err != nil {
//     log.Fatal(err)
// }
// fmt.Println("Çözülmüş:", decrypted)
// ```
//
// # Hata Durumları
//
// - Geçersiz Base64 formatı
// - Şifreli veri çok kısa (< 16 byte)
// - Geçersiz anahtar uzunluğu
// - Şifreli veri blok boyutunun katı değil
// - Geçersiz padding (veri bozulmuş veya yanlış anahtar)
//
// # Güvenlik Notları
//
// ⚠️ **Anahtar Eşleşmesi**: Şifreleme sırasında kullanılan aynı anahtar kullanılmalıdır.
//
// ⚠️ **Veri Bütünlüğü**: Bu fonksiyon veri bütünlüğü kontrolü yapmaz. Veri manipüle
// edilmişse padding hatası alınabilir.
//
// ⚠️ **Hata Mesajları**: Padding hataları timing attack'lere karşı hassas olabilir.
// Üretim ortamında genel hata mesajları kullanılmalıdır.
//
// # Performans
//
// - Küçük veriler için hızlı (< 1ms)
// - Büyük veriler için lineer zaman karmaşıklığı O(n)
// - Bellek kullanımı: ~2x veri boyutu
//
// # Uyumluluk
//
// Bu fonksiyon, aynı paketteki encrypt fonksiyonu ile tam uyumludur.
// Farklı sistemler arası veri alışverişinde aynı anahtar ve algoritma
// kullanıldığından emin olunmalıdır.
func decrypt(b64cipher string, key []byte) (string, error) {
	// Base64 kodlu veriyi decode et
	data, err := base64.StdEncoding.DecodeString(b64cipher)
	if err != nil {
		return "", err
	}

	// Veri en az bir blok boyutunda olmalı (IV için)
	if len(data) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	// İlk 16 byte IV, geri kalanı şifreli veri
	iv := data[:aes.BlockSize]
	ciphertext := data[aes.BlockSize:]

	// AES cipher bloğu oluştur
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Şifreli veri blok boyutunun katı olmalı
	if len(ciphertext)%aes.BlockSize != 0 {
		return "", fmt.Errorf("ciphertext is not a multiple of the block size")
	}

	// CBC modu ile şifre çöz
	mode := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(ciphertext))
	mode.CryptBlocks(decrypted, ciphertext)

	// PKCS7 padding'i kaldır
	unpadded, err := pkcs7Unpad(decrypted, aes.BlockSize)
	if err != nil {
		return "", err
	}

	return string(unpadded), nil
}

// pkcs7Pad, veriyi PKCS7 standardına göre doldurur (padding).
//
// Bu fonksiyon, blok şifreleme algoritmaları için gerekli olan padding işlemini
// gerçekleştirir. PKCS7 standardı, verinin blok boyutunun tam katı olmasını
// sağlamak için veri sonuna belirli sayıda byte ekler.
//
// # PKCS7 Padding Standardı
//
// PKCS7 padding, eklenen byte sayısını her padding byte'ında saklar:
// - 1 byte eksikse: 0x01 eklenir
// - 2 byte eksikse: 0x02 0x02 eklenir
// - 3 byte eksikse: 0x03 0x03 0x03 eklenir
// - Veri tam blok boyutundaysa: tam bir blok padding eklenir (0x10 x 16)
//
// # Padding Süreci
//
// 1. Veri uzunluğunun blok boyutuna göre modunu al
// 2. Eksik byte sayısını hesapla
// 3. Her padding byte'ına eksik byte sayısını yaz
// 4. Padding'i verinin sonuna ekle
//
// # Parametreler
//
// - `data`: Doldurulacak veri (byte dizisi)
// - `blockSize`: Blok boyutu (AES için 16 byte)
//
// # Dönüş Değeri
//
// - `[]byte`: Padding eklenmiş veri
//
// # Kullanım Örneği
//
// ```go
// data := []byte("Hello")           // 5 byte
// blockSize := 16                   // AES blok boyutu
// padded := pkcs7Pad(data, blockSize)
// // Sonuç: "Hello" + 11 adet 0x0B byte (toplam 16 byte)
// ```
//
// # Özel Durumlar
//
// - Veri zaten blok boyutunun tam katıysa, tam bir blok padding eklenir
// - Bu, padding'in her zaman kaldırılabilir olmasını garanti eder
// - Minimum 1, maksimum blockSize kadar padding eklenir
//
// # Güvenlik Notları
//
// ⚠️ **Padding Oracle**: PKCS7 padding, padding oracle saldırılarına karşı
// hassas olabilir. Hata mesajları dikkatli yönetilmelidir.
//
// ⚠️ **Veri Bütünlüğü**: Padding, veri bütünlüğü kontrolü sağlamaz. HMAC veya
// authenticated encryption kullanımı önerilir.
//
// # Performans
//
// - Sabit zaman karmaşıklığı O(1)
// - Bellek kullanımı: veri boyutu + maksimum blockSize
// - Çok hızlı işlem (< 1μs)
//
// # Standart Uyumluluk
//
// Bu implementasyon RFC 5652 (PKCS #7) standardına tam uyumludur ve
// diğer PKCS7 uyumlu sistemlerle uyumlu çalışır.
func pkcs7Pad(data []byte, blockSize int) []byte {
	// Eksik byte sayısını hesapla
	padding := blockSize - len(data)%blockSize
	// Padding byte'larını oluştur (her byte padding sayısını içerir)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	// Padding'i verinin sonuna ekle
	return append(data, padtext...)
}

// pkcs7Unpad, PKCS7 padding'ini kaldırır ve orijinal veriyi döndürür.
//
// Bu fonksiyon, pkcs7Pad tarafından eklenen padding byte'larını güvenli bir
// şekilde kaldırır. Padding'in geçerliliğini kontrol eder ve bozuk veri
// durumunda hata döndürür.
//
// # Padding Kaldırma Süreci
//
// 1. Veri uzunluğunu ve blok boyutunu kontrol et
// 2. Son byte'ı oku (padding sayısını içerir)
// 3. Padding sayısının geçerliliğini doğrula
// 4. Tüm padding byte'larının doğru değere sahip olduğunu kontrol et
// 5. Padding'i kaldır ve orijinal veriyi döndür
//
// # Parametreler
//
// - `data`: Padding içeren veri (byte dizisi)
// - `blockSize`: Blok boyutu (AES için 16 byte)
//
// # Dönüş Değerleri
//
// - `[]byte`: Padding kaldırılmış orijinal veri
// - `error`: Geçersiz padding durumunda hata
//
// # Kullanım Örneği
//
// ```go
// // Padding'li veri (16 byte)
// padded := []byte("Hello\x0b\x0b\x0b\x0b\x0b\x0b\x0b\x0b\x0b\x0b\x0b")
// blockSize := 16
//
// unpadded, err := pkcs7Unpad(padded, blockSize)
// if err != nil {
//     log.Fatal(err)
// }
// fmt.Println(string(unpadded)) // "Hello"
// ```
//
// # Hata Durumları
//
// - Veri uzunluğu 0 veya blok boyutunun katı değil
// - Padding sayısı blok boyutundan büyük veya 0
// - Padding byte'ları tutarsız (farklı değerler içeriyor)
// - Veri bozulmuş veya yanlış anahtar ile şifre çözülmüş
//
// # Güvenlik Notları
//
// ⚠️ **Padding Oracle Saldırıları**: Bu fonksiyon padding hatalarını açıkça
// raporlar. Üretim ortamında, hata mesajları genel tutulmalı ve timing
// saldırılarına karşı önlem alınmalıdır.
//
// ⚠️ **Veri Doğrulama**: Padding hatası, veri bozulması veya yanlış anahtar
// kullanımının göstergesi olabilir. Bu durumlar güvenlik loglarına kaydedilmelidir.
//
// ⚠️ **Constant-Time Karşılaştırma**: Mevcut implementasyon constant-time değildir.
// Yüksek güvenlik gerektiren uygulamalarda constant-time karşılaştırma kullanılmalıdır.
//
// # Performans
//
// - Lineer zaman karmaşıklığı O(n) - padding doğrulama için
// - Bellek kullanımı: veri boyutu - padding
// - Çok hızlı işlem (< 1μs küçük veriler için)
//
// # Standart Uyumluluk
//
// Bu implementasyon RFC 5652 (PKCS #7) standardına tam uyumludur.
// pkcs7Pad fonksiyonu ile tam uyumlu çalışır.
//
// # Özel Durumlar
//
// - Tam blok padding (16 byte padding) doğru şekilde işlenir
// - Minimum 1, maksimum blockSize kadar padding kaldırılabilir
// - Sıfır uzunluklu veri geçersiz kabul edilir
//
// # Hata Yönetimi
//
// Üretim ortamında, padding hatalarını kullanıcıya açıkça bildirmek yerine:
// ```go
// if err != nil {
//     log.Error("Decryption failed", "error", err)
//     return nil, errors.New("decryption failed")
// }
// ```
// şeklinde genel bir hata mesajı kullanılması önerilir.
func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	length := len(data)
	// Veri uzunluğu 0 veya blok boyutunun katı değilse hata
	if length == 0 || length%blockSize != 0 {
		return nil, fmt.Errorf("invalid padding size")
	}

	// Son byte padding sayısını içerir
	padding := int(data[length-1])
	// Padding sayısı geçerli aralıkta olmalı (1 ile blockSize arası)
	if padding > blockSize || padding == 0 {
		return nil, fmt.Errorf("invalid padding")
	}

	// Tüm padding byte'larının aynı değere sahip olduğunu doğrula
	for i := 0; i < padding; i++ {
		if data[length-1-i] != byte(padding) {
			return nil, fmt.Errorf("invalid padding byte")
		}
	}

	// Padding'i kaldır ve orijinal veriyi döndür
	return data[:length-padding], nil
}
