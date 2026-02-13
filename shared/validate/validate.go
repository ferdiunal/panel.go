// Bu paket, yapı (struct) doğrulaması ve hata yönetimi için merkezi bir doğrulama sistemi sağlar.
// Go-playground/validator/v10 kütüphanesini kullanarak struct alanlarını doğrular ve
// doğrulama hatalarını yapılandırılmış bir formatta döndürür.
//
// Kullanım Senaryoları:
// - API isteklerinin gelen verilerini doğrulama
// - Form verilerinin geçerliliğini kontrol etme
// - Veritabanı modelleri için veri bütünlüğü sağlama
// - İş kurallarına uygun veri doğrulaması
package validate

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

// Bu global değişken, tüm doğrulama işlemleri için kullanılan validator örneğidir.
// Validator.New() ile oluşturulur ve uygulamanın tamamında paylaşılır.
//
// Önemli Notlar:
// - Validator örneği thread-safe'dir ve eş zamanlı kullanım için güvenlidir
// - Aynı örneği tekrar tekrar kullanmak performans açısından daha verimlidir
// - Özel doğrulama kuralları eklemek için bu değişkeni kullanabilirsiniz
//
// Örnek:
//   Validate.RegisterValidation("custom_rule", customValidationFunc)
var Validate = validator.New()

// Bu yapı, doğrulama hatalarını yapılandırılmış bir formatta temsil eder.
// Alanın adı (key) ile hata bilgilerini (value) eşleştirir.
//
// Yapı Açıklaması:
// - Dış map: Alan adı (string) -> Hata detayları (map)
// - İç map: Hata türü (string) -> Hata mesajı (string)
//
// Örnek Yapı:
//   {
//     "email": {
//       "message": "email"
//     },
//     "age": {
//       "message": "min"
//     }
//   }
//
// Kullanım Senaryoları:
// - API yanıtlarında doğrulama hatalarını döndürme
// - Frontend'e hata bilgilerini JSON formatında gönderme
// - Hata mesajlarını yerelleştirme (i18n) için temel veri yapısı
type ValidationError map[string]map[string]string

// Bu fonksiyon, verilen struct'ı doğrular ve doğrulama hatalarını yapılandırılmış formatta döndürür.
//
// Parametreler:
// - s (interface{}): Doğrulanacak struct örneği. Herhangi bir struct türü olabilir.
//
// Dönüş Değeri:
// - ValidationError: Doğrulama hatalarını içeren map. Hata yoksa boş map döndürülür.
//
// Fonksiyon Davranışı:
// 1. Validator.Struct() kullanarak struct'ı doğrular
// 2. Doğrulama hatası varsa, her hata için:
//    - Alan adını küçük harfe dönüştürür (strings.ToLower)
//    - Hata türünü (tag) "message" anahtarı altında saklar
// 3. Tüm hataları ValidationError formatında döndürür
//
// Kullanım Örnekleri:
//
// Örnek 1: Basit Doğrulama
//   type User struct {
//     Email string `validate:"required,email"`
//     Age   int    `validate:"required,min=18"`
//   }
//
//   user := User{Email: "invalid", Age: 15}
//   errors := ValidateStruct(user)
//   // Sonuç: {"email": {"message": "email"}, "age": {"message": "min"}}
//
// Örnek 2: Geçerli Veri
//   user := User{Email: "test@example.com", Age: 25}
//   errors := ValidateStruct(user)
//   // Sonuç: {} (boş map, hata yok)
//
// Örnek 3: API Handler'da Kullanım
//   func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
//     var user User
//     json.NewDecoder(r.Body).Decode(&user)
//
//     if errors := ValidateStruct(user); len(errors) > 0 {
//       w.Header().Set("Content-Type", "application/json")
//       w.WriteHeader(http.StatusBadRequest)
//       json.NewEncoder(w).Encode(map[string]interface{}{
//         "success": false,
//         "errors": errors,
//       })
//       return
//     }
//     // Doğrulama başarılı, işleme devam et
//   }
//
// Önemli Notlar:
// - Struct alanları doğrulanabilir tag'ler içermelidir (validate:"...")
// - Alan adları otomatik olarak küçük harfe dönüştürülür
// - Hata mesajları validator tag'inin adıdır (örn: "required", "email", "min")
// - Eğer struct doğrulama tag'i içermiyorsa, hata döndürülmez
// - Pointer alanlar için doğrulama yapılmaz, nil değer kabul edilir
//
// Uyarılar:
// - Fonksiyon interface{} kabul ettiği için, yanlış türde veri geçirilirse panic oluşabilir
// - Doğrulama tag'leri yanlış yazılırsa, doğrulama çalışmayabilir
// - Kompleks doğrulama kuralları için özel validator'lar tanımlanmalıdır
func ValidateStruct(s interface{}) ValidationError {
	// Hata bilgilerini saklamak için boş bir map oluştur
	errors := make(map[string]map[string]string)

	// Validator.Struct() ile struct'ı doğrula
	err := Validate.Struct(s)

	// Eğer doğrulama hatası varsa, hataları işle
	if err != nil {
		// validator.ValidationErrors türüne dönüştür ve her hatayı döngüyle işle
		for _, err := range err.(validator.ValidationErrors) {
			// Alan adını küçük harfe dönüştür ve hata bilgisini sakla
			// Örn: "Email" -> "email", hata türü: "email"
			errors[strings.ToLower(err.Field())] = map[string]string{
				"message": err.Tag(),
			}
		}
	}

	// Doğrulama hatalarını döndür (hata yoksa boş map)
	return errors
}
