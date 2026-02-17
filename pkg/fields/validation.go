// Package fields, admin panel için alan (field) tanımlamalarını sağlar.
//
// Bu dosya, alan doğrulama (validation) sistemini içerir.
// Alanlar için çeşitli doğrulama kuralları tanımlanabilir ve bu kurallar
// form gönderildiğinde otomatik olarak çalıştırılır.
package fields

import (
	"fmt"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
)

// ValidationRule, bir alan için doğrulama kuralını temsil eder.
//
// Doğrulama kuralları, form verilerinin belirli kriterlere uygun olup olmadığını kontrol eder.
// Her kural bir ad, parametreler ve hata mesajı içerir.
//
// # Kullanım Örneği
//
//	rule := fields.Required()
//	rule := fields.MinLength(5)
//	rule := fields.EmailRule()
//
// Daha fazla bilgi için docs/Fields.md dosyasına bakın.
type ValidationRule struct {
	Name       string        // Kural adı (required, email, min, vb.)
	Parameters []interface{} // Kural parametreleri
	Message    string        // Hata mesajı
}

// ValidatorFunc, bir değeri doğrulayan fonksiyondur.
//
// Bu tip, özel doğrulama mantığı tanımlamak için kullanılır.
// Fonksiyon, değer geçerliyse nil, geçersizse hata döndürür.
//
// # Kullanım Örneği
//
//	validator := func(value interface{}, context interface{}) error {
//	    if value == nil {
//	        return fmt.Errorf("değer boş olamaz")
//	    }
//	    return nil
//	}
type ValidatorFunc func(value interface{}, context interface{}) error

// ValidationResult, doğrulama sonucunu temsil eder.
//
// Bu yapı, doğrulama işleminin başarılı olup olmadığını ve
// varsa hata mesajlarını içerir.
//
// # Kullanım Örneği
//
//	result := ValidationResult{
//	    IsValid: false,
//	    Errors: map[string][]string{
//	        "email": {"Geçerli bir e-posta adresi giriniz"},
//	    },
//	}
type ValidationResult struct {
	IsValid bool                // Doğrulama başarılı mı?
	Errors  map[string][]string // Alan adı -> hata mesajları
}

// Built-in validators (Yerleşik doğrulayıcılar)

// Required, zorunlu alanlar için doğrulama kuralı döndürür.
//
// Bu kural, alanın boş olmamasını zorunlu kılar.
// Boş string, nil, boş slice gibi değerler geçersiz sayılır.
//
// # Kullanım Örneği
//
//	field.Rules(fields.Required())
//	// Alan boş bırakılamaz
//
// Döndürür:
//   - ValidationRule: "This field is required" hata mesajı ile
func Required() ValidationRule {
	return ValidationRule{
		Name:    "required",
		Message: "This field is required",
	}
}

// EmailRule, e-posta formatı için doğrulama kuralı döndürür.
//
// Bu kural, alanın geçerli bir e-posta adresi formatında olmasını zorunlu kılar.
// RFC 5322 standardına uygun e-posta adresleri kabul edilir.
//
// # Kullanım Örneği
//
//	field.Rules(fields.EmailRule())
//	// Geçerli e-posta formatı: user@example.com
//
// Döndürür:
//   - ValidationRule: "This field must be a valid email address" hata mesajı ile
func EmailRule() ValidationRule {
	return ValidationRule{
		Name:    "email",
		Message: "This field must be a valid email address",
	}
}

// URL, URL formatı için doğrulama kuralı döndürür.
//
// Bu kural, alanın geçerli bir URL formatında olmasını zorunlu kılar.
// HTTP, HTTPS gibi protokoller desteklenir.
//
// # Kullanım Örneği
//
//	field.Rules(fields.URL())
//	// Geçerli URL formatı: https://example.com
//
// Döndürür:
//   - ValidationRule: "This field must be a valid URL" hata mesajı ile
func URL() ValidationRule {
	return ValidationRule{
		Name:    "url",
		Message: "This field must be a valid URL",
	}
}

// Min, minimum değer için doğrulama kuralı döndürür.
//
// Bu kural, sayısal alanların belirtilen minimum değerden büyük veya eşit olmasını zorunlu kılar.
//
// # Parametreler
//
// - **min**: Minimum değer (int, float, vb.)
//
// # Kullanım Örneği
//
//	field.Rules(fields.Min(18))
//	// Değer en az 18 olmalıdır
//
// Döndürür:
//   - ValidationRule: "This field must be at least {min}" hata mesajı ile
func Min(min interface{}) ValidationRule {
	return ValidationRule{
		Name:       "min",
		Parameters: []interface{}{min},
		Message:    fmt.Sprintf("This field must be at least %v", min),
	}
}

// Max, maksimum değer için doğrulama kuralı döndürür.
//
// Bu kural, sayısal alanların belirtilen maksimum değerden küçük veya eşit olmasını zorunlu kılar.
//
// # Parametreler
//
// - **max**: Maksimum değer (int, float, vb.)
//
// # Kullanım Örneği
//
//	field.Rules(fields.Max(100))
//	// Değer en fazla 100 olmalıdır
//
// Döndürür:
//   - ValidationRule: "This field must be at most {max}" hata mesajı ile
func Max(max interface{}) ValidationRule {
	return ValidationRule{
		Name:       "max",
		Parameters: []interface{}{max},
		Message:    fmt.Sprintf("This field must be at most %v", max),
	}
}

// MinLength, minimum string uzunluğu için doğrulama kuralı döndürür.
//
// Bu kural, string alanların belirtilen minimum karakter sayısından uzun veya eşit olmasını zorunlu kılar.
//
// # Parametreler
//
// - **length**: Minimum karakter sayısı
//
// # Kullanım Örneği
//
//	field.Rules(fields.MinLength(5))
//	// String en az 5 karakter olmalıdır
//
// Döndürür:
//   - ValidationRule: "This field must be at least {length} characters" hata mesajı ile
func MinLength(length int) ValidationRule {
	return ValidationRule{
		Name:       "minLength",
		Parameters: []interface{}{length},
		Message:    fmt.Sprintf("This field must be at least %d characters", length),
	}
}

// MaxLength, maksimum string uzunluğu için doğrulama kuralı döndürür.
//
// Bu kural, string alanların belirtilen maksimum karakter sayısından kısa veya eşit olmasını zorunlu kılar.
//
// # Parametreler
//
// - **length**: Maksimum karakter sayısı
//
// # Kullanım Örneği
//
//	field.Rules(fields.MaxLength(255))
//	// String en fazla 255 karakter olmalıdır
//
// Döndürür:
//   - ValidationRule: "This field must be at most {length} characters" hata mesajı ile
func MaxLength(length int) ValidationRule {
	return ValidationRule{
		Name:       "maxLength",
		Parameters: []interface{}{length},
		Message:    fmt.Sprintf("This field must be at most %d characters", length),
	}
}

// Pattern, regex pattern eşleştirme için doğrulama kuralı döndürür.
//
// Bu kural, alanın belirtilen regex pattern'ine uygun olmasını zorunlu kılar.
//
// # Parametreler
//
// - **pattern**: Regex pattern (örn. "^[A-Z][a-z]+$")
//
// # Kullanım Örneği
//
//	field.Rules(fields.Pattern("^[0-9]{5}$"))
//	// 5 haneli sayı formatı (örn. 12345)
//
// Döndürür:
//   - ValidationRule: "This field format is invalid" hata mesajı ile
func Pattern(pattern string) ValidationRule {
	return ValidationRule{
		Name:       "pattern",
		Parameters: []interface{}{pattern},
		Message:    "This field format is invalid",
	}
}

// Unique, veritabanı benzersizliği için doğrulama kuralı döndürür.
//
// Bu kural, alanın veritabanında benzersiz olmasını zorunlu kılar.
// Aynı değere sahip başka bir kayıt varsa doğrulama başarısız olur.
//
// # Parametreler
//
// - **table**: Tablo adı
// - **column**: Sütun adı
//
// # Kullanım Örneği
//
//	field.Rules(fields.Unique("users", "email"))
//	// E-posta adresi users tablosunda benzersiz olmalıdır
//
// Döndürür:
//   - ValidationRule: "This value already exists" hata mesajı ile
func Unique(table, column string) ValidationRule {
	return ValidationRule{
		Name:       "unique",
		Parameters: []interface{}{table, column},
		Message:    "This value already exists",
	}
}

// Exists, veritabanı varlığı için doğrulama kuralı döndürür.
//
// Bu kural, alanın veritabanında mevcut olmasını zorunlu kılar.
// Belirtilen değere sahip bir kayıt yoksa doğrulama başarısız olur.
//
// # Parametreler
//
// - **table**: Tablo adı
// - **column**: Sütun adı
//
// # Kullanım Örneği
//
//	field.Rules(fields.Exists("categories", "id"))
//	// Kategori ID'si categories tablosunda mevcut olmalıdır
//
// Döndürür:
//   - ValidationRule: "This value does not exist" hata mesajı ile
func Exists(table, column string) ValidationRule {
	return ValidationRule{
		Name:       "exists",
		Parameters: []interface{}{table, column},
		Message:    "This value does not exist",
	}
}

// ValidateRequired, bir değerin boş olmadığını doğrular.
//
// Bu fonksiyon, değerin nil, boş string, boş slice olmadığını kontrol eder.
// Farklı veri tipleri için farklı boşluk kontrolleri yapar.
//
// # Parametreler
//
// - **value**: Doğrulanacak değer (interface{})
//
// # Desteklenen Tipler
//
// - **nil**: Hata döndürür
// - **string**: Boşluk karakterleri temizlendikten sonra boş string kontrolü
// - **[]interface{}**: Boş slice kontrolü
// - **[]string**: Boş string slice kontrolü
//
// # Kullanım Örneği
//
//	err := fields.ValidateRequired("test")  // nil
//	err := fields.ValidateRequired("")      // error: "value is required"
//	err := fields.ValidateRequired(nil)     // error: "value is required"
//
// Döndürür:
//   - nil: Değer geçerli
//   - error: Değer boş
func ValidateRequired(value interface{}) error {
	if value == nil {
		return fmt.Errorf("value is required")
	}

	switch v := value.(type) {
	case string:
		if strings.TrimSpace(v) == "" {
			return fmt.Errorf("value is required")
		}
	case []interface{}:
		if len(v) == 0 {
			return fmt.Errorf("value is required")
		}
	case []string:
		if len(v) == 0 {
			return fmt.Errorf("value is required")
		}
	}

	return nil
}

// ValidateEmail, e-posta formatını doğrular.
//
// Bu fonksiyon, RFC 5322 standardına uygun e-posta adresi formatını kontrol eder.
// Boş string geçerli sayılır (zorunluluk için Required() kullanın).
//
// # Parametreler
//
// - **value**: Doğrulanacak değer (string olmalı)
//
// # Kullanım Örneği
//
//	err := fields.ValidateEmail("user@example.com")  // nil
//	err := fields.ValidateEmail("invalid-email")     // error: "invalid email format"
//	err := fields.ValidateEmail("")                  // nil (boş geçerli)
//
// Döndürür:
//   - nil: E-posta formatı geçerli veya boş
//   - error: E-posta formatı geçersiz veya tip uyumsuz
func ValidateEmail(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("email must be a string")
	}

	if str == "" {
		return nil // Empty is valid, use Required() for mandatory
	}

	_, err := mail.ParseAddress(str)
	if err != nil {
		return fmt.Errorf("invalid email format")
	}

	return nil
}

// ValidateURL validates URL format
func ValidateURL(value interface{}) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("url must be a string")
	}

	if str == "" {
		return nil // Empty is valid, use Required() for mandatory
	}

	_, err := url.ParseRequestURI(str)
	if err != nil {
		return fmt.Errorf("invalid URL format")
	}

	return nil
}

// ValidateMin validates minimum value
func ValidateMin(value interface{}, min interface{}) error {
	switch v := value.(type) {
	case int:
		minVal, ok := min.(int)
		if !ok {
			return fmt.Errorf("min parameter must be int")
		}
		if v < minVal {
			return fmt.Errorf("value must be at least %d", minVal)
		}
	case float64:
		minVal, ok := min.(float64)
		if !ok {
			return fmt.Errorf("min parameter must be float64")
		}
		if v < minVal {
			return fmt.Errorf("value must be at least %f", minVal)
		}
	case string:
		minVal, ok := min.(int)
		if !ok {
			return fmt.Errorf("min parameter must be int for string length")
		}
		if len(v) < minVal {
			return fmt.Errorf("value must be at least %d characters", minVal)
		}
	}

	return nil
}

// ValidateMax validates maximum value
func ValidateMax(value interface{}, max interface{}) error {
	switch v := value.(type) {
	case int:
		maxVal, ok := max.(int)
		if !ok {
			return fmt.Errorf("max parameter must be int")
		}
		if v > maxVal {
			return fmt.Errorf("value must be at most %d", maxVal)
		}
	case float64:
		maxVal, ok := max.(float64)
		if !ok {
			return fmt.Errorf("max parameter must be float64")
		}
		if v > maxVal {
			return fmt.Errorf("value must be at most %f", maxVal)
		}
	case string:
		maxVal, ok := max.(int)
		if !ok {
			return fmt.Errorf("max parameter must be int for string length")
		}
		if len(v) > maxVal {
			return fmt.Errorf("value must be at most %d characters", maxVal)
		}
	}

	return nil
}

// ValidateMinLength validates minimum string length
func ValidateMinLength(value interface{}, length int) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("minLength validation requires string value")
	}

	if len(str) < length {
		return fmt.Errorf("value must be at least %d characters", length)
	}

	return nil
}

// ValidateMaxLength validates maximum string length
func ValidateMaxLength(value interface{}, length int) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("maxLength validation requires string value")
	}

	if len(str) > length {
		return fmt.Errorf("value must be at most %d characters", length)
	}

	return nil
}

// ValidatePattern validates regex pattern matching
func ValidatePattern(value interface{}, pattern string) error {
	str, ok := value.(string)
	if !ok {
		return fmt.Errorf("pattern validation requires string value")
	}

	if str == "" {
		return nil // Empty is valid, use Required() for mandatory
	}

	matched, err := regexp.MatchString(pattern, str)
	if err != nil {
		return fmt.Errorf("invalid regex pattern: %v", err)
	}

	if !matched {
		return fmt.Errorf("value does not match required pattern")
	}

	return nil
}

// ApplyValidationRule applies a validation rule to a value
func ApplyValidationRule(rule ValidationRule, value interface{}) error {
	switch rule.Name {
	case "required":
		return ValidateRequired(value)
	case "email":
		return ValidateEmail(value)
	case "url":
		return ValidateURL(value)
	case "min":
		if len(rule.Parameters) > 0 {
			return ValidateMin(value, rule.Parameters[0])
		}
	case "max":
		if len(rule.Parameters) > 0 {
			return ValidateMax(value, rule.Parameters[0])
		}
	case "minLength":
		if len(rule.Parameters) > 0 {
			if length, ok := rule.Parameters[0].(int); ok {
				return ValidateMinLength(value, length)
			}
		}
	case "maxLength":
		if len(rule.Parameters) > 0 {
			if length, ok := rule.Parameters[0].(int); ok {
				return ValidateMaxLength(value, length)
			}
		}
	case "pattern":
		if len(rule.Parameters) > 0 {
			if pattern, ok := rule.Parameters[0].(string); ok {
				return ValidatePattern(value, pattern)
			}
		}
	case "unique":
		// Unique validation requires database access, implement in handler
		return nil
	case "exists":
		// Exists validation requires database access, implement in handler
		return nil
	}

	return nil
}

// ConditionalValidator represents a validator that applies conditionally
type ConditionalValidator struct {
	Condition func(context interface{}) bool
	Validator ValidatorFunc
}

// ApplyConditionalValidation applies a validator only if condition is met
func ApplyConditionalValidation(validator ConditionalValidator, value interface{}, context interface{}) error {
	if validator.Condition(context) {
		return validator.Validator(value, context)
	}
	return nil
}

// CustomValidatorRegistry stores custom validators
type CustomValidatorRegistry struct {
	validators map[string]ValidatorFunc
}

// NewCustomValidatorRegistry creates a new custom validator registry
func NewCustomValidatorRegistry() *CustomValidatorRegistry {
	return &CustomValidatorRegistry{
		validators: make(map[string]ValidatorFunc),
	}
}

// Register registers a custom validator
func (r *CustomValidatorRegistry) Register(name string, validator ValidatorFunc) {
	r.validators[name] = validator
}

// Get retrieves a custom validator
func (r *CustomValidatorRegistry) Get(name string) (ValidatorFunc, bool) {
	validator, ok := r.validators[name]
	return validator, ok
}

// Apply applies a custom validator
func (r *CustomValidatorRegistry) Apply(name string, value interface{}, context interface{}) error {
	validator, ok := r.Get(name)
	if !ok {
		return fmt.Errorf("custom validator '%s' not found", name)
	}
	return validator(value, context)
}

// Global custom validator registry
var globalValidatorRegistry = NewCustomValidatorRegistry()

// RegisterCustomValidator registers a global custom validator
func RegisterCustomValidator(name string, validator ValidatorFunc) {
	globalValidatorRegistry.Register(name, validator)
}

// GetCustomValidator retrieves a global custom validator
func GetCustomValidator(name string) (ValidatorFunc, bool) {
	return globalValidatorRegistry.Get(name)
}

// ApplyCustomValidator applies a global custom validator
func ApplyCustomValidator(name string, value interface{}, context interface{}) error {
	return globalValidatorRegistry.Apply(name, value, context)
}

func MergeValidationRules(base, extra []ValidationRule) []ValidationRule {
	if len(extra) == 0 {
		result := make([]ValidationRule, len(base))
		copy(result, base)
		return result
	}
	merged := make([]ValidationRule, 0, len(base)+len(extra))
	merged = append(merged, base...)
	merged = append(merged, extra...)
	return merged
}
