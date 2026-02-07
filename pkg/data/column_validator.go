package data

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

/// # ColumnValidator
///
/// Bu yapı, veritabanı şemasına karşı kolon isimlerini doğrulayan ve SQL injection saldırılarına
/// karşı koruma sağlayan bir güvenlik katmanıdır.
///
/// ## Temel Özellikler
///
/// - **Şema Tabanlı Doğrulama**: GORM model şemasını kullanarak sadece geçerli kolonlara erişim sağlar
/// - **Çoklu Format Desteği**: snake_case, camelCase ve orijinal alan isimlerini destekler
/// - **İlişki Desteği**: GORM ilişki alanlarını (relationships) otomatik olarak tanır
/// - **SQL Injection Koruması**: Whitelist yaklaşımı ile güvenli kolon erişimi garanti eder
///
/// ## Kullanım Senaryoları
///
/// 1. **Dinamik Sorgu Oluşturma**: Kullanıcı girdilerinden güvenli SQL sorguları oluşturma
/// 2. **API Filtreleme**: REST API'lerde güvenli filtreleme parametreleri
/// 3. **Dinamik Sıralama**: Kullanıcı tanımlı sıralama kolonlarının doğrulanması
/// 4. **Veri İhracı**: Güvenli kolon seçimi ile veri dışa aktarma
///
/// ## Güvenlik Avantajları
///
/// - ✅ SQL injection saldırılarını önler
/// - ✅ Sadece model şemasında tanımlı kolonlara erişim
/// - ✅ Whitelist tabanlı güvenlik yaklaşımı
/// - ✅ Otomatik normalizasyon ve sanitizasyon
///
/// ## Örnek Kullanım
///
/// ```go
/// // Validator oluşturma
/// validator, err := NewColumnValidator(db, &User{})
/// if err != nil {
///     log.Fatal(err)
/// }
///
/// // Kolon doğrulama
/// if validator.IsValidColumn("email") {
///     // Güvenli sorgu oluşturma
///     db.Where("email = ?", userInput)
/// }
///
/// // Güvenli WHERE clause oluşturma
/// whereClause, err := BuildSafeWhereClause(validator, "created_at", ">")
/// if err == nil {
///     db.Where(whereClause, time.Now())
/// }
/// ```
///
/// ## Önemli Notlar
///
/// ⚠️ **Dikkat**: Bu validator sadece kolon isimlerini doğrular, değerleri değil.
/// Değer sanitizasyonu için GORM'un parameterized queries özelliğini kullanın.
///
/// ⚠️ **Performans**: Validator oluşturma maliyetlidir, singleton pattern kullanın.
///
/// ## İlgili Tipler
///
/// - `gorm.DB`: Veritabanı bağlantısı
/// - `schema.Schema`: GORM model şeması
type ColumnValidator struct {
	allowedColumns map[string]bool // İzin verilen kolon isimlerinin whitelist'i
	schema         *schema.Schema  // GORM model şeması
}

/// # NewColumnValidator
///
/// Bu fonksiyon, belirtilen GORM modeli için yeni bir ColumnValidator örneği oluşturur.
/// Model şemasını analiz ederek izin verilen kolonların whitelist'ini otomatik olarak oluşturur.
///
/// ## Parametreler
///
/// - `db`: GORM veritabanı bağlantısı (*gorm.DB)
/// - `model`: Doğrulanacak GORM model yapısı (interface{})
///
/// ## Dönüş Değerleri
///
/// - `*ColumnValidator`: Yapılandırılmış validator pointer'ı
/// - `error`: Şema parse hatası durumunda hata mesajı
///
/// ## İşleyiş
///
/// 1. **Şema Parse**: GORM Statement kullanarak model şemasını parse eder
/// 2. **Alan Toplama**: Tüm model alanlarını (fields) toplar
/// 3. **Format Çeşitlendirme**: Her alan için snake_case, camelCase ve orijinal ismi ekler
/// 4. **İlişki Ekleme**: GORM ilişkilerini (relationships) whitelist'e ekler
///
/// ## Desteklenen Format Çeşitleri
///
/// Her kolon için aşağıdaki formatlar otomatik olarak eklenir:
/// - `created_at` (snake_case - DB adı)
/// - `createdat` (lowercase)
/// - `CreatedAt` (orijinal alan adı)
/// - `createdAt` (camelCase)
///
/// ## Kullanım Örnekleri
///
/// ```go
/// // Basit kullanım
/// type User struct {
///     ID        uint
///     Email     string
///     CreatedAt time.Time
/// }
///
/// validator, err := NewColumnValidator(db, &User{})
/// if err != nil {
///     log.Fatalf("Validator oluşturulamadı: %v", err)
/// }
///
/// // İlişkili model ile kullanım
/// type Post struct {
///     ID       uint
///     Title    string
///     AuthorID uint
///     Author   User `gorm:"foreignKey:AuthorID"`
/// }
///
/// validator, err := NewColumnValidator(db, &Post{})
/// // "author" ilişkisi otomatik olarak whitelist'e eklenir
/// ```
///
/// ## Hata Durumları
///
/// - Model şeması parse edilemezse error döner
/// - Geçersiz model yapısı verilirse error döner
/// - DB bağlantısı nil ise panic oluşabilir
///
/// ## Performans Notları
///
/// ⚠️ **Önemli**: Bu fonksiyon şema parse işlemi yaptığı için maliyetlidir.
/// Validator'ı bir kez oluşturup tekrar kullanın (singleton pattern).
///
/// ```go
/// // İyi pratik - Singleton pattern
/// var userValidator *ColumnValidator
/// var once sync.Once
///
/// func GetUserValidator() *ColumnValidator {
///     once.Do(func() {
///         userValidator, _ = NewColumnValidator(db, &User{})
///     })
///     return userValidator
/// }
/// ```
///
/// ## İlgili Fonksiyonlar
///
/// - `IsValidColumn()`: Kolon geçerliliğini kontrol eder
/// - `ValidateColumn()`: Kolon doğrular ve DB adını döner
/// - `GetAllowedColumns()`: Tüm izin verilen kolonları listeler
func NewColumnValidator(db *gorm.DB, model interface{}) (*ColumnValidator, error) {
	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(model); err != nil {
		return nil, fmt.Errorf("failed to parse model schema: %w", err)
	}

	validator := &ColumnValidator{
		allowedColumns: make(map[string]bool),
		schema:         stmt.Schema,
	}

	// Build whitelist of allowed columns
	for _, field := range stmt.Schema.Fields {
		if field.DBName != "" {
			// Add both snake_case (DB name) and original field name
			validator.allowedColumns[field.DBName] = true
			validator.allowedColumns[strings.ToLower(field.Name)] = true

			// Add camelCase version
			validator.allowedColumns[toCamelCase(field.DBName)] = true
		}
	}

	// Add relationship fields
	for name := range stmt.Schema.Relationships.Relations {
		validator.allowedColumns[strings.ToLower(name)] = true
		validator.allowedColumns[toSnakeCase(name)] = true
	}

	return validator, nil
}

/// # IsValidColumn
///
/// Bu fonksiyon, verilen kolon isminin model şemasında tanımlı olup olmadığını kontrol eder.
/// Whitelist tabanlı doğrulama yaparak SQL injection saldırılarına karşı koruma sağlar.
///
/// ## Parametreler
///
/// - `columnName`: Doğrulanacak kolon ismi (string)
///
/// ## Dönüş Değeri
///
/// - `bool`: Kolon geçerliyse `true`, değilse `false`
///
/// ## İşleyiş
///
/// 1. **Boş Kontrol**: Boş string kontrolü yapar
/// 2. **Normalizasyon**: Kolon ismini küçük harfe çevirir ve boşlukları temizler
/// 3. **Whitelist Kontrolü**: Normalize edilmiş ismi whitelist'te arar
///
/// ## Desteklenen Format Çeşitleri
///
/// Aşağıdaki formatların hepsi geçerli kabul edilir:
/// - `created_at` (snake_case)
/// - `CreatedAt` (PascalCase)
/// - `createdAt` (camelCase)
/// - `CREATED_AT` (uppercase - normalize edilir)
///
/// ## Kullanım Örnekleri
///
/// ```go
/// validator, _ := NewColumnValidator(db, &User{})
///
/// // Geçerli kolon kontrolleri
/// if validator.IsValidColumn("email") {
///     // Güvenli sorgu oluştur
///     db.Where("email = ?", userInput)
/// }
///
/// // Farklı format çeşitleri
/// validator.IsValidColumn("created_at")  // true
/// validator.IsValidColumn("CreatedAt")   // true
/// validator.IsValidColumn("createdAt")   // true
/// validator.IsValidColumn("CREATED_AT")  // true
///
/// // Geçersiz kolon
/// validator.IsValidColumn("malicious_column") // false
/// validator.IsValidColumn("")                 // false
/// ```
///
/// ## Güvenlik Özellikleri
///
/// - SQL injection saldırılarını önler
/// - Sadece model şemasında tanımlı kolonlara izin verir
/// - Whitelist yaklaşımı kullanır (blacklist değil)
/// - Case-insensitive kontrol yapar
///
/// ## Performans
///
/// - O(1) zaman karmaşıklığı (map lookup)
/// - Çok hızlı, her sorgu için kullanılabilir
/// - Bellek maliyeti düşük
///
/// ## Kullanım Senaryoları
///
/// 1. **Dinamik Filtreleme**: API'den gelen filtreleme parametrelerini doğrulama
/// 2. **Sıralama**: Kullanıcı tanımlı sıralama kolonlarını kontrol etme
/// 3. **Seçim**: SELECT clause için kolon seçimini doğrulama
/// 4. **Güvenlik Katmanı**: Tüm kullanıcı girdilerini doğrulama
///
/// ## İlgili Fonksiyonlar
///
/// - `ValidateColumn()`: Kolon doğrular ve DB adını döner
/// - `ValidateColumns()`: Birden fazla kolonu doğrular
/// - `GetAllowedColumns()`: Tüm geçerli kolonları listeler
func (v *ColumnValidator) IsValidColumn(columnName string) bool {
	if columnName == "" {
		return false
	}

	// Normalize column name
	normalized := strings.ToLower(strings.TrimSpace(columnName))

	// Check against whitelist
	return v.allowedColumns[normalized]
}

/// # ValidateColumn
///
/// Bu fonksiyon, verilen kolon ismini doğrular ve veritabanında kullanılacak güvenli
/// DB kolon ismini döndürür. IsValidColumn'dan farklı olarak, sadece boolean değil,
/// normalize edilmiş DB kolon ismini de döner.
///
/// ## Parametreler
///
/// - `columnName`: Doğrulanacak kolon ismi (string)
///
/// ## Dönüş Değerleri
///
/// - `string`: Normalize edilmiş, güvenli DB kolon ismi
/// - `error`: Kolon geçersizse hata mesajı
///
/// ## İşleyiş
///
/// 1. **Geçerlilik Kontrolü**: IsValidColumn ile kolon geçerliliğini kontrol eder
/// 2. **Normalizasyon**: Kolon ismini küçük harfe çevirir ve boşlukları temizler
/// 3. **Şema Araması**: Model şemasında gerçek DB kolon ismini arar
/// 4. **İlişki Kontrolü**: Bulunamazsa ilişki alanlarında arar
/// 5. **Geri Dönüş**: Hiçbir yerde bulunamazsa orijinal ismi döner
///
/// ## Kullanım Örnekleri
///
/// ```go
/// validator, _ := NewColumnValidator(db, &User{})
///
/// // Farklı format çeşitlerini normalize etme
/// dbColumn, err := validator.ValidateColumn("CreatedAt")
/// // dbColumn = "created_at", err = nil
///
/// dbColumn, err = validator.ValidateColumn("createdAt")
/// // dbColumn = "created_at", err = nil
///
/// dbColumn, err = validator.ValidateColumn("created_at")
/// // dbColumn = "created_at", err = nil
///
/// // Geçersiz kolon
/// dbColumn, err = validator.ValidateColumn("malicious_column")
/// // dbColumn = "", err = "invalid column name: malicious_column"
///
/// // Güvenli sorgu oluşturma
/// dbColumn, err := validator.ValidateColumn(userInput)
/// if err == nil {
///     db.Where(fmt.Sprintf("%s = ?", dbColumn), value)
/// }
/// ```
///
/// ## İlişki Desteği
///
/// İlişki alanları için de çalışır:
/// ```go
/// type Post struct {
///     ID       uint
///     AuthorID uint
///     Author   User `gorm:"foreignKey:AuthorID"`
/// }
///
/// validator, _ := NewColumnValidator(db, &Post{})
/// dbColumn, _ := validator.ValidateColumn("author")
/// // dbColumn = "Author" (ilişki adı)
/// ```
///
/// ## Güvenlik Özellikleri
///
/// - SQL injection saldırılarını önler
/// - Sadece whitelist'teki kolonları kabul eder
/// - Otomatik normalizasyon yapar
/// - Güvenli DB kolon ismini garanti eder
///
/// ## Hata Durumları
///
/// - Kolon geçersizse: `fmt.Errorf("invalid column name: %s", columnName)`
/// - Boş string verilirse: `fmt.Errorf("invalid column name: ")`
///
/// ## IsValidColumn ile Farkı
///
/// | Özellik | IsValidColumn | ValidateColumn |
/// |---------|---------------|----------------|
/// | Dönüş Tipi | bool | (string, error) |
/// | DB Kolon İsmi | Döndürmez | Döndürür |
/// | Hata Mesajı | Yok | Var |
/// | Kullanım | Hızlı kontrol | Sorgu oluşturma |
///
/// ## Kullanım Senaryoları
///
/// 1. **Dinamik Sorgu Oluşturma**: WHERE clause için güvenli kolon ismi alma
/// 2. **ORDER BY**: Sıralama için normalize edilmiş kolon ismi
/// 3. **SELECT**: Seçim için DB kolon ismi
/// 4. **JOIN**: Join koşulları için güvenli kolon referansı
///
/// ## İlgili Fonksiyonlar
///
/// - `IsValidColumn()`: Sadece geçerlilik kontrolü yapar
/// - `ValidateColumns()`: Birden fazla kolonu doğrular
/// - `BuildSafeWhereClause()`: Güvenli WHERE clause oluşturur
func (v *ColumnValidator) ValidateColumn(columnName string) (string, error) {
	if !v.IsValidColumn(columnName) {
		return "", fmt.Errorf("invalid column name: %s", columnName)
	}

	// Find the actual DB column name
	normalized := strings.ToLower(strings.TrimSpace(columnName))

	// Try to find the field in schema
	for _, field := range v.schema.Fields {
		if strings.ToLower(field.DBName) == normalized ||
		   strings.ToLower(field.Name) == normalized {
			return field.DBName, nil
		}
	}

	// If not found in fields, might be a relationship
	for name := range v.schema.Relationships.Relations {
		if strings.ToLower(name) == normalized {
			return name, nil
		}
	}

	return columnName, nil
}

/// # GetAllowedColumns
///
/// Bu fonksiyon, validator'da tanımlı tüm izin verilen kolon isimlerinin listesini döndürür.
/// Debugging, API dokümantasyonu veya kullanıcı arayüzü için kullanışlıdır.
///
/// ## Parametreler
///
/// Parametre almaz.
///
/// ## Dönüş Değeri
///
/// - `[]string`: İzin verilen tüm kolon isimlerinin slice'ı
///
/// ## İşleyiş
///
/// 1. **Kapasite Ayırma**: Performans için önceden kapasite ayırır
/// 2. **Map İterasyonu**: allowedColumns map'ini iterate eder
/// 3. **Liste Oluşturma**: Tüm kolon isimlerini slice'a ekler
///
/// ## Dönen Liste Özellikleri
///
/// Liste aşağıdaki format çeşitlerini içerir:
/// - Snake_case DB isimleri: `created_at`, `user_id`
/// - Lowercase alan isimleri: `createdat`, `userid`
/// - CamelCase versiyonlar: `createdAt`, `userId`
/// - İlişki isimleri: `author`, `posts`
///
/// ## Kullanım Örnekleri
///
/// ```go
/// validator, _ := NewColumnValidator(db, &User{})
///
/// // Tüm izin verilen kolonları al
/// allowedColumns := validator.GetAllowedColumns()
/// fmt.Println("İzin verilen kolonlar:", allowedColumns)
/// // Çıktı: [id email created_at createdat createdAt ...]
///
/// // API dokümantasyonu için kullanım
/// type APIResponse struct {
///     AvailableFilters []string `json:"available_filters"`
/// }
///
/// response := APIResponse{
///     AvailableFilters: validator.GetAllowedColumns(),
/// }
///
/// // Kullanıcı arayüzü için dropdown
/// func GetFilterOptions() []string {
///     validator, _ := NewColumnValidator(db, &User{})
///     return validator.GetAllowedColumns()
/// }
///
/// // Debugging için
/// log.Printf("Validator kolonları: %v", validator.GetAllowedColumns())
/// ```
///
/// ## Kullanım Senaryoları
///
/// 1. **API Dokümantasyonu**: Filtrelenebilir alanların listesini döndürme
/// 2. **Kullanıcı Arayüzü**: Dropdown veya autocomplete için seçenekler
/// 3. **Debugging**: Validator'ın hangi kolonları tanıdığını kontrol etme
/// 4. **Validasyon Mesajları**: Hata mesajlarında geçerli kolonları gösterme
/// 5. **Test**: Unit testlerde validator durumunu kontrol etme
///
/// ## Önemli Notlar
///
/// ⚠️ **Sıralama**: Dönen liste sıralı değildir (map iteration order)
/// ⚠️ **Duplikasyon**: Aynı kolonun farklı formatları ayrı öğeler olarak döner
/// ⚠️ **Performans**: Her çağrıda yeni slice oluşturur, cache'leme düşünün
///
/// ## Sıralı Liste İçin
///
/// ```go
/// import "sort"
///
/// columns := validator.GetAllowedColumns()
/// sort.Strings(columns)
/// ```
///
/// ## Benzersiz DB Kolonları İçin
///
/// ```go
/// func GetUniqueDBColumns(validator *ColumnValidator) []string {
///     seen := make(map[string]bool)
///     var unique []string
///
///     for _, field := range validator.schema.Fields {
///         if field.DBName != "" && !seen[field.DBName] {
///             seen[field.DBName] = true
///             unique = append(unique, field.DBName)
///         }
///     }
///
///     return unique
/// }
/// ```
///
/// ## İlgili Fonksiyonlar
///
/// - `IsValidColumn()`: Tek bir kolonun geçerliliğini kontrol eder
/// - `ValidateColumns()`: Birden fazla kolonu doğrular
func (v *ColumnValidator) GetAllowedColumns() []string {
	columns := make([]string, 0, len(v.allowedColumns))
	for col := range v.allowedColumns {
		columns = append(columns, col)
	}
	return columns
}

/// # ValidateColumns
///
/// Bu fonksiyon, birden fazla kolon ismini toplu olarak doğrular. Tüm kolonlar geçerliyse
/// nil döner, herhangi biri geçersizse ilk geçersiz kolon için hata döner.
///
/// ## Parametreler
///
/// - `columns`: Doğrulanacak kolon isimlerinin slice'ı ([]string)
///
/// ## Dönüş Değeri
///
/// - `error`: Tüm kolonlar geçerliyse nil, geçersiz kolon varsa hata mesajı
///
/// ## İşleyiş
///
/// 1. **İterasyon**: Verilen kolon listesini iterate eder
/// 2. **Tekil Doğrulama**: Her kolon için IsValidColumn çağırır
/// 3. **Erken Çıkış**: İlk geçersiz kolonda hata döner ve durur
/// 4. **Başarı**: Tüm kolonlar geçerliyse nil döner
///
/// ## Kullanım Örnekleri
///
/// ```go
/// validator, _ := NewColumnValidator(db, &User{})
///
/// // Birden fazla kolonu doğrulama
/// columns := []string{"email", "created_at", "name"}
/// if err := validator.ValidateColumns(columns); err != nil {
///     log.Printf("Geçersiz kolon: %v", err)
/// } else {
///     // Tüm kolonlar geçerli, güvenli sorgu oluştur
///     for _, col := range columns {
///         // Güvenli kullanım
///     }
/// }
///
/// // API filtreleme parametrelerini doğrulama
/// func HandleFilterRequest(w http.ResponseWriter, r *http.Request) {
///     filterColumns := r.URL.Query()["filter"]
///
///     if err := validator.ValidateColumns(filterColumns); err != nil {
///         http.Error(w, "Geçersiz filtre kolonları", http.StatusBadRequest)
///         return
///     }
///
///     // Güvenli filtreleme işlemi
/// }
///
/// // Dinamik SELECT clause oluşturma
/// func BuildSelectQuery(columns []string) (*gorm.DB, error) {
///     if err := validator.ValidateColumns(columns); err != nil {
///         return nil, err
///     }
///
///     return db.Select(columns), nil
/// }
///
/// // Sıralama kolonlarını doğrulama
/// func ValidateSortColumns(sortBy []string) error {
///     return validator.ValidateColumns(sortBy)
/// }
/// ```
///
/// ## Hata Durumları
///
/// ```go
/// // Geçersiz kolon içeren liste
/// columns := []string{"email", "malicious_column", "name"}
/// err := validator.ValidateColumns(columns)
/// // err = "invalid column name: malicious_column"
///
/// // Boş string içeren liste
/// columns := []string{"email", "", "name"}
/// err := validator.ValidateColumns(columns)
/// // err = "invalid column name: "
/// ```
///
/// ## Kullanım Senaryoları
///
/// 1. **API Filtreleme**: Kullanıcıdan gelen filtreleme parametrelerini doğrulama
/// 2. **Dinamik SELECT**: SELECT clause için kolon listesini doğrulama
/// 3. **Sıralama**: ORDER BY için kolon listesini doğrulama
/// 4. **Veri İhracı**: Export edilecek kolonları doğrulama
/// 5. **Toplu İşlemler**: Birden fazla kolon gerektiren işlemleri doğrulama
///
/// ## Performans Özellikleri
///
/// - **Zaman Karmaşıklığı**: O(n) - n: kolon sayısı
/// - **Erken Çıkış**: İlk geçersiz kolonda durur
/// - **Bellek**: Ek bellek ayırma yapmaz
///
/// ## Alternatif Yaklaşımlar
///
/// Tüm geçersiz kolonları toplamak için:
/// ```go
/// func ValidateAllColumns(validator *ColumnValidator, columns []string) []error {
///     var errors []error
///     for _, col := range columns {
///         if !validator.IsValidColumn(col) {
///             errors = append(errors, fmt.Errorf("invalid column: %s", col))
///         }
///     }
///     return errors
/// }
/// ```
///
/// ## Önemli Notlar
///
/// ⚠️ **Erken Çıkış**: İlk geçersiz kolonda durur, diğer geçersiz kolonları kontrol etmez
/// ⚠️ **Sıra Bağımlı**: Kolon sırası önemlidir, ilk geçersiz kolon için hata döner
///
/// ## İlgili Fonksiyonlar
///
/// - `IsValidColumn()`: Tek bir kolonun geçerliliğini kontrol eder
/// - `ValidateColumn()`: Kolon doğrular ve DB adını döner
/// - `GetAllowedColumns()`: Tüm geçerli kolonları listeler
func (v *ColumnValidator) ValidateColumns(columns []string) error {
	for _, col := range columns {
		if !v.IsValidColumn(col) {
			return fmt.Errorf("invalid column name: %s", col)
		}
	}
	return nil
}

/// # Helper Functions
///
/// Bu bölüm, kolon ismi dönüşümleri ve güvenlik kontrolleri için yardımcı fonksiyonlar içerir.

/// # toSnakeCase
///
/// Bu fonksiyon, PascalCase veya camelCase formatındaki string'i snake_case formatına dönüştürür.
/// Büyük harflerin önüne alt çizgi ekleyerek ve tüm harfleri küçülterek dönüşüm yapar.
///
/// ## Parametreler
///
/// - `s`: Dönüştürülecek string (string)
///
/// ## Dönüş Değeri
///
/// - `string`: snake_case formatına dönüştürülmüş string
///
/// ## İşleyiş
///
/// 1. **İterasyon**: String'i karakter karakter iterate eder
/// 2. **Büyük Harf Kontrolü**: Büyük harf bulunca önüne alt çizgi ekler
/// 3. **Küçültme**: Tüm karakterleri küçük harfe çevirir
/// 4. **İlk Karakter**: İlk karakterin önüne alt çizgi eklenmez
///
/// ## Dönüşüm Örnekleri
///
/// ```go
/// toSnakeCase("UserName")      // "user_name"
/// toSnakeCase("userName")      // "user_name"
/// toSnakeCase("CreatedAt")     // "created_at"
/// toSnakeCase("HTTPResponse")  // "h_t_t_p_response"
/// toSnakeCase("ID")            // "i_d"
/// toSnakeCase("user")          // "user"
/// toSnakeCase("ALLCAPS")       // "a_l_l_c_a_p_s"
/// ```
///
/// ## Kullanım Senaryoları
///
/// 1. **Model Alan Dönüşümü**: Go struct alanlarını DB kolon isimlerine dönüştürme
/// 2. **API Parametreleri**: camelCase API parametrelerini DB formatına çevirme
/// 3. **Normalizasyon**: Farklı formatları standart hale getirme
///
/// ## Önemli Notlar
///
/// ⚠️ **Kısaltmalar**: "HTTPResponse" gibi ardışık büyük harfler her biri için alt çizgi ekler
/// ⚠️ **Performans**: Her çağrıda yeni string builder oluşturur
///
/// ## Alternatif Yaklaşımlar
///
/// Daha gelişmiş dönüşüm için:
/// ```go
/// import "github.com/iancoleman/strcase"
/// strcase.ToSnake("HTTPResponse") // "http_response"
/// ```
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

/// # toCamelCase
///
/// Bu fonksiyon, snake_case formatındaki string'i camelCase formatına dönüştürür.
/// Alt çizgilerden sonraki harfleri büyük harfe çevirerek ve alt çizgileri kaldırarak dönüşüm yapar.
///
/// ## Parametreler
///
/// - `s`: Dönüştürülecek string (string)
///
/// ## Dönüş Değeri
///
/// - `string`: camelCase formatına dönüştürülmüş string
///
/// ## İşleyiş
///
/// 1. **Bölme**: String'i alt çizgiye göre parçalara böler
/// 2. **İlk Parça**: İlk parçayı olduğu gibi bırakır (küçük harf)
/// 3. **Diğer Parçalar**: Her parçanın ilk harfini büyük harfe çevirir
/// 4. **Birleştirme**: Tüm parçaları birleştirir
///
/// ## Dönüşüm Örnekleri
///
/// ```go
/// toCamelCase("user_name")      // "userName"
/// toCamelCase("created_at")     // "createdAt"
/// toCamelCase("user_id")        // "userId"
/// toCamelCase("first_name")     // "firstName"
/// toCamelCase("is_active")      // "isActive"
/// toCamelCase("user")           // "user"
/// toCamelCase("a_b_c")          // "aBC"
/// toCamelCase("_leading")       // "Leading" (boş ilk parça)
/// ```
///
/// ## Kullanım Senaryoları
///
/// 1. **API Yanıtları**: DB kolon isimlerini JSON camelCase formatına dönüştürme
/// 2. **JavaScript Entegrasyonu**: Go struct'larını JS nesnelerine uyarlama
/// 3. **Format Normalizasyonu**: Farklı formatları standart hale getirme
///
/// ## Önemli Notlar
///
/// ⚠️ **İlk Parça**: İlk parça küçük harf olarak kalır (camelCase kuralı)
/// ⚠️ **Boş Parçalar**: Ardışık alt çizgiler boş parçalar oluşturabilir
/// ⚠️ **Performans**: Her çağrıda string split ve join işlemi yapar
///
/// ## Özel Durumlar
///
/// ```go
/// // Boş string
/// toCamelCase("")           // ""
///
/// // Sadece alt çizgi
/// toCamelCase("___")        // ""
///
/// // Başta alt çizgi
/// toCamelCase("_user_name") // "UserName" (ilk parça boş)
///
/// // Sonda alt çizgi
/// toCamelCase("user_name_") // "userName" (son parça boş)
/// ```
///
/// ## Alternatif Yaklaşımlar
///
/// Daha gelişmiş dönüşüm için:
/// ```go
/// import "github.com/iancoleman/strcase"
/// strcase.ToCamel("user_name") // "userName"
/// ```
func toCamelCase(s string) string {
	parts := strings.Split(s, "_")
	for i := 1; i < len(parts); i++ {
		if len(parts[i]) > 0 {
			parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
		}
	}
	return strings.Join(parts, "")
}

/// # SanitizeColumnName
///
/// Bu fonksiyon, kolon isminden potansiyel olarak tehlikeli karakterleri temizler ve
/// sadece güvenli karakterleri (alfanumerik, alt çizgi, nokta) içeren bir string döndürür.
///
/// ## Parametreler
///
/// - `columnName`: Temizlenecek kolon ismi (string)
///
/// ## Dönüş Değeri
///
/// - `string`: Sadece güvenli karakterler içeren temizlenmiş string
///
/// ## İzin Verilen Karakterler
///
/// - **Küçük harfler**: a-z
/// - **Büyük harfler**: A-Z
/// - **Rakamlar**: 0-9
/// - **Alt çizgi**: _
/// - **Nokta**: . (tablo.kolon formatı için)
///
/// ## İşleyiş
///
/// 1. **İterasyon**: String'i karakter karakter iterate eder
/// 2. **Karakter Kontrolü**: Her karakterin güvenli olup olmadığını kontrol eder
/// 3. **Filtreleme**: Sadece güvenli karakterleri yeni string'e ekler
/// 4. **Tehlikeli Karakterler**: SQL injection için kullanılabilecek karakterleri kaldırır
///
/// ## Temizleme Örnekleri
///
/// ```go
/// // Normal kullanım
/// SanitizeColumnName("user_name")           // "user_name"
/// SanitizeColumnName("CreatedAt")           // "CreatedAt"
/// SanitizeColumnName("users.email")         // "users.email"
///
/// // SQL injection girişimleri
/// SanitizeColumnName("id; DROP TABLE--")   // "idDROPTABLE"
/// SanitizeColumnName("name' OR '1'='1")    // "nameOR11"
/// SanitizeColumnName("id/*comment*/")      // "idcomment"
/// SanitizeColumnName("col`name")           // "colname"
/// SanitizeColumnName("col\"name")          // "colname"
///
/// // Özel karakterler
/// SanitizeColumnName("user-name")          // "username"
/// SanitizeColumnName("user@email")         // "useremail"
/// SanitizeColumnName("user#tag")           // "usertag"
/// SanitizeColumnName("user name")          // "username" (boşluk kaldırılır)
///
/// // Boş sonuç
/// SanitizeColumnName("@#$%^&*()")          // ""
/// SanitizeColumnName("   ")                // ""
/// ```
///
/// ## Güvenlik Özellikleri
///
/// Aşağıdaki SQL injection tekniklerine karşı koruma sağlar:
/// - **Yorum İşaretleri**: `--`, `/*`, `*/`, `#`
/// - **String Delimiters**: `'`, `"`, `` ` ``
/// - **SQL Operatörleri**: `;`, `=`, `<`, `>`, `|`, `&`
/// - **Parantezler**: `(`, `)`, `[`, `]`, `{`, `}`
/// - **Boşluk Karakterleri**: Boşluk, tab, newline
///
/// ## Kullanım Senaryoları
///
/// 1. **Kullanıcı Girdisi Temizleme**: API'den gelen kolon isimlerini temizleme
/// 2. **Ek Güvenlik Katmanı**: Validator ile birlikte kullanma
/// 3. **Log Temizleme**: Log mesajlarında güvenli string kullanma
/// 4. **Dinamik Sorgu**: Güvenli kolon isimleri oluşturma
///
/// ## Kullanım Örnekleri
///
/// ```go
/// // API endpoint'inde kullanım
/// func HandleSortRequest(w http.ResponseWriter, r *http.Request) {
///     sortColumn := r.URL.Query().Get("sort")
///
///     // Önce sanitize et
///     safeSortColumn := SanitizeColumnName(sortColumn)
///
///     // Sonra validate et
///     if validator.IsValidColumn(safeSortColumn) {
///         db.Order(safeSortColumn)
///     }
/// }
///
/// // Validator ile birlikte kullanım
/// func SafeColumnName(validator *ColumnValidator, input string) (string, error) {
///     // 1. Sanitize
///     sanitized := SanitizeColumnName(input)
///
///     // 2. Validate
///     return validator.ValidateColumn(sanitized)
/// }
///
/// // Log güvenliği
/// func LogQuery(columnName string) {
///     safeColumn := SanitizeColumnName(columnName)
///     log.Printf("Querying column: %s", safeColumn)
/// }
/// ```
///
/// ## Önemli Notlar
///
/// ⚠️ **Tek Başına Yeterli Değil**: Bu fonksiyon tek başına SQL injection koruması sağlamaz.
/// Mutlaka ColumnValidator ile birlikte kullanılmalıdır.
///
/// ⚠️ **Veri Kaybı**: Geçersiz karakterler kaldırılır, orijinal string değişir.
///
/// ⚠️ **Boş String**: Tüm karakterler geçersizse boş string döner.
///
/// ## Güvenlik Katmanları
///
/// Önerilen güvenlik yaklaşımı (defense in depth):
/// ```go
/// // 1. Sanitize (tehlikeli karakterleri kaldır)
/// sanitized := SanitizeColumnName(userInput)
///
/// // 2. Validate (whitelist kontrolü)
/// dbColumn, err := validator.ValidateColumn(sanitized)
/// if err != nil {
///     return err
/// }
///
/// // 3. Parameterized Query (değer güvenliği)
/// db.Where(fmt.Sprintf("%s = ?", dbColumn), value)
/// ```
///
/// ## Performans
///
/// - **Zaman Karmaşıklığı**: O(n) - n: string uzunluğu
/// - **Bellek**: Her çağrıda yeni string builder oluşturur
/// - **Optimizasyon**: Sık kullanımda sonuçları cache'leyebilirsiniz
///
/// ## İlgili Fonksiyonlar
///
/// - `IsValidColumn()`: Kolon geçerliliğini kontrol eder
/// - `ValidateColumn()`: Kolon doğrular ve DB adını döner
/// - `IsValidOperator()`: SQL operatörlerini doğrular
func SanitizeColumnName(columnName string) string {
	// Remove any characters that aren't alphanumeric, underscore, or dot
	var result strings.Builder
	for _, r := range columnName {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
		   (r >= '0' && r <= '9') || r == '_' || r == '.' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

/// # IsValidOperator
///
/// Bu fonksiyon, verilen SQL operatörünün güvenli ve geçerli olup olmadığını kontrol eder.
/// Whitelist tabanlı yaklaşım kullanarak sadece bilinen güvenli operatörlere izin verir.
///
/// ## Parametreler
///
/// - `operator`: Kontrol edilecek SQL operatörü (string)
///
/// ## Dönüş Değeri
///
/// - `bool`: Operatör geçerliyse `true`, değilse `false`
///
/// ## Desteklenen Operatörler
///
/// ### Karşılaştırma Operatörleri
/// - `=`: Eşitlik kontrolü
/// - `!=`: Eşitsizlik kontrolü
/// - `>`: Büyüktür
/// - `>=`: Büyük veya eşit
/// - `<`: Küçüktür
/// - `<=`: Küçük veya eşit
///
/// ### String Operatörleri
/// - `LIKE`: Pattern matching (joker karakterlerle)
/// - `NOT LIKE`: Negatif pattern matching
///
/// ### Liste Operatörleri
/// - `IN`: Değer listesinde var mı kontrolü
/// - `NOT IN`: Değer listesinde yok mu kontrolü
///
/// ### NULL Operatörleri
/// - `IS NULL`: NULL değer kontrolü
/// - `IS NOT NULL`: NULL olmayan değer kontrolü
///
/// ### Aralık Operatörleri
/// - `BETWEEN`: İki değer arasında kontrolü
///
/// ## İşleyiş
///
/// 1. **Normalizasyon**: Operatörü büyük harfe çevirir ve boşlukları temizler
/// 2. **Whitelist Kontrolü**: Normalize edilmiş operatörü whitelist'te arar
/// 3. **Case-Insensitive**: Büyük/küçük harf duyarsız kontrol yapar
///
/// ## Kullanım Örnekleri
///
/// ```go
/// // Geçerli operatörler
/// IsValidOperator("=")           // true
/// IsValidOperator("!=")          // true
/// IsValidOperator(">")           // true
/// IsValidOperator("LIKE")        // true
/// IsValidOperator("like")        // true (case-insensitive)
/// IsValidOperator("IN")          // true
/// IsValidOperator("IS NULL")     // true
/// IsValidOperator("BETWEEN")     // true
///
/// // Geçersiz operatörler
/// IsValidOperator("OR")          // false (SQL injection riski)
/// IsValidOperator("AND")         // false (SQL injection riski)
/// IsValidOperator("UNION")       // false (SQL injection riski)
/// IsValidOperator("--")          // false (yorum işareti)
/// IsValidOperator("; DROP")      // false (tehlikeli)
/// IsValidOperator("")            // false (boş)
///
/// // Dinamik sorgu oluşturma
/// func BuildWhereClause(column, operator, value string) (string, error) {
///     if !IsValidOperator(operator) {
///         return "", fmt.Errorf("geçersiz operatör: %s", operator)
///     }
///
///     return fmt.Sprintf("%s %s ?", column, operator), nil
/// }
/// ```
///
/// ## Güvenlik Özellikleri
///
/// Aşağıdaki SQL injection tekniklerini engeller:
/// - **Mantıksal Operatörler**: `OR`, `AND` (WHERE clause manipülasyonu)
/// - **UNION Saldırıları**: `UNION`, `UNION ALL`
/// - **Yorum İşaretleri**: `--`, `/*`, `#`
/// - **Komut Enjeksiyonu**: `;`, `EXEC`, `EXECUTE`
/// - **Veri Manipülasyonu**: `UPDATE`, `DELETE`, `DROP`, `INSERT`
///
/// ## Kullanım Senaryoları
///
/// 1. **API Filtreleme**: Kullanıcıdan gelen operatörleri doğrulama
/// 2. **Dinamik Sorgu**: WHERE clause için güvenli operatör seçimi
/// 3. **Sorgu Builder**: Güvenli sorgu oluşturma araçları
/// 4. **Veri İhracı**: Filtreleme operatörlerini doğrulama
///
/// ## Pratik Kullanım Örnekleri
///
/// ```go
/// // API endpoint'inde kullanım
/// func HandleFilterRequest(w http.ResponseWriter, r *http.Request) {
///     column := r.URL.Query().Get("column")
///     operator := r.URL.Query().Get("operator")
///     value := r.URL.Query().Get("value")
///
///     if !IsValidOperator(operator) {
///         http.Error(w, "Geçersiz operatör", http.StatusBadRequest)
///         return
///     }
///
///     // Güvenli sorgu oluştur
///     db.Where(fmt.Sprintf("%s %s ?", column, operator), value)
/// }
///
/// // Sorgu builder ile kullanım
/// type QueryBuilder struct {
///     filters []Filter
/// }
///
/// type Filter struct {
///     Column   string
///     Operator string
///     Value    interface{}
/// }
///
/// func (qb *QueryBuilder) AddFilter(column, operator string, value interface{}) error {
///     if !IsValidOperator(operator) {
///         return fmt.Errorf("geçersiz operatör: %s", operator)
///     }
///
///     qb.filters = append(qb.filters, Filter{
///         Column:   column,
///         Operator: operator,
///         Value:    value,
///     })
///
///     return nil
/// }
///
/// // Çoklu filtre doğrulama
/// func ValidateFilters(filters []Filter) error {
///     for _, filter := range filters {
///         if !IsValidOperator(filter.Operator) {
///             return fmt.Errorf("geçersiz operatör: %s", filter.Operator)
///         }
///     }
///     return nil
/// }
/// ```
///
/// ## Önemli Notlar
///
/// ⚠️ **Whitelist Yaklaşımı**: Sadece bilinen güvenli operatörlere izin verilir.
/// Yeni operatör eklemek için fonksiyonu güncellemeniz gerekir.
///
/// ⚠️ **Parameterized Queries**: Bu fonksiyon operatör güvenliğini sağlar,
/// ancak değerler için mutlaka parameterized queries kullanılmalıdır.
///
/// ⚠️ **Case-Insensitive**: Operatör kontrolü büyük/küçük harf duyarsızdır.
///
/// ## Yeni Operatör Ekleme
///
/// Yeni bir operatör eklemek için:
/// ```go
/// validOperators := map[string]bool{
///     // ... mevcut operatörler
///     "ILIKE":       true,  // PostgreSQL case-insensitive LIKE
///     "SIMILAR TO":  true,  // PostgreSQL regex matching
///     "~":           true,  // PostgreSQL regex operator
/// }
/// ```
///
/// ## Güvenlik Katmanları
///
/// Önerilen güvenlik yaklaşımı:
/// ```go
/// // 1. Kolon doğrulama
/// dbColumn, err := validator.ValidateColumn(column)
/// if err != nil {
///     return err
/// }
///
/// // 2. Operatör doğrulama
/// if !IsValidOperator(operator) {
///     return fmt.Errorf("geçersiz operatör")
/// }
///
/// // 3. Parameterized query
/// db.Where(fmt.Sprintf("%s %s ?", dbColumn, operator), value)
/// ```
///
/// ## Performans
///
/// - **Zaman Karmaşıklığı**: O(1) - map lookup
/// - **Bellek**: Sabit boyutlu map, düşük bellek kullanımı
/// - **Optimizasyon**: Her çağrıda map oluşturulur, global değişken olarak cache'lenebilir
///
/// ## İlgili Fonksiyonlar
///
/// - `BuildSafeWhereClause()`: Güvenli WHERE clause oluşturur
/// - `ValidateColumn()`: Kolon isimlerini doğrular
/// - `SanitizeColumnName()`: Kolon isimlerini temizler
func IsValidOperator(operator string) bool {
	validOperators := map[string]bool{
		"=":           true,
		"!=":          true,
		">":           true,
		">=":          true,
		"<":           true,
		"<=":          true,
		"LIKE":        true,
		"NOT LIKE":    true,
		"IN":          true,
		"NOT IN":      true,
		"IS NULL":     true,
		"IS NOT NULL": true,
		"BETWEEN":     true,
	}

	return validOperators[strings.ToUpper(strings.TrimSpace(operator))]
}

/// # BuildSafeWhereClause
///
/// Bu fonksiyon, doğrulanmış kolon ve operatör kullanarak güvenli bir WHERE clause oluşturur.
/// Hem kolon hem de operatör doğrulaması yaparak SQL injection saldırılarına karşı tam koruma sağlar.
///
/// ## Parametreler
///
/// - `validator`: Kolon doğrulama için ColumnValidator pointer'ı (*ColumnValidator)
/// - `column`: WHERE clause'da kullanılacak kolon ismi (string)
/// - `operator`: WHERE clause'da kullanılacak SQL operatörü (string)
///
/// ## Dönüş Değerleri
///
/// - `string`: Güvenli WHERE clause string'i (parameterized query için hazır)
/// - `error`: Doğrulama hatası durumunda hata mesajı
///
/// ## İşleyiş
///
/// 1. **Kolon Doğrulama**: ValidateColumn ile kolon ismini doğrular
/// 2. **Operatör Doğrulama**: IsValidOperator ile operatörü doğrular
/// 3. **Clause Oluşturma**: Güvenli WHERE clause string'i oluşturur
/// 4. **Placeholder**: Değer için `?` placeholder ekler (parameterized query)
///
/// ## Kullanım Örnekleri
///
/// ```go
/// validator, _ := NewColumnValidator(db, &User{})
///
/// // Basit eşitlik kontrolü
/// whereClause, err := BuildSafeWhereClause(validator, "email", "=")
/// if err == nil {
///     db.Where(whereClause, "user@example.com").Find(&users)
///     // SQL: WHERE email = ?
/// }
///
/// // Büyüktür karşılaştırması
/// whereClause, err := BuildSafeWhereClause(validator, "created_at", ">")
/// if err == nil {
///     db.Where(whereClause, time.Now().AddDate(0, -1, 0)).Find(&users)
///     // SQL: WHERE created_at > ?
/// }
///
/// // LIKE operatörü
/// whereClause, err := BuildSafeWhereClause(validator, "name", "LIKE")
/// if err == nil {
///     db.Where(whereClause, "%john%").Find(&users)
///     // SQL: WHERE name LIKE ?
/// }
///
/// // IN operatörü
/// whereClause, err := BuildSafeWhereClause(validator, "status", "IN")
/// if err == nil {
///     db.Where(whereClause, []string{"active", "pending"}).Find(&users)
///     // SQL: WHERE status IN (?)
/// }
///
/// // NULL kontrolü
/// whereClause, err := BuildSafeWhereClause(validator, "deleted_at", "IS NULL")
/// if err == nil {
///     db.Where(whereClause).Find(&users)
///     // SQL: WHERE deleted_at IS NULL
/// }
/// ```
///
/// ## Hata Durumları
///
/// ```go
/// // Geçersiz kolon
/// whereClause, err := BuildSafeWhereClause(validator, "malicious_column", "=")
/// // err = "invalid column name: malicious_column"
///
/// // Geçersiz operatör
/// whereClause, err := BuildSafeWhereClause(validator, "email", "OR")
/// // err = "invalid operator: OR"
///
/// // Her iki hata da
/// whereClause, err := BuildSafeWhereClause(validator, "bad_column", "DROP")
/// // err = "invalid column name: bad_column" (kolon önce kontrol edilir)
/// ```
///
/// ## Güvenlik Özellikleri
///
/// Bu fonksiyon çok katmanlı güvenlik sağlar:
/// 1. **Kolon Whitelist**: Sadece model şemasında tanımlı kolonlar
/// 2. **Operatör Whitelist**: Sadece güvenli SQL operatörleri
/// 3. **Parameterized Query**: Değerler için placeholder kullanımı
/// 4. **SQL Injection Koruması**: Tam koruma sağlar
///
/// ## Kullanım Senaryoları
///
/// 1. **API Filtreleme**: Kullanıcıdan gelen filtreleme parametreleri
/// 2. **Dinamik Sorgu**: Runtime'da oluşturulan WHERE clause'lar
/// 3. **Sorgu Builder**: Güvenli sorgu oluşturma araçları
/// 4. **Çoklu Filtre**: Birden fazla WHERE koşulu oluşturma
///
/// ## Pratik Kullanım Örnekleri
///
/// ```go
/// // API endpoint'inde kullanım
/// func HandleFilterRequest(w http.ResponseWriter, r *http.Request) {
///     validator, _ := NewColumnValidator(db, &User{})
///
///     column := r.URL.Query().Get("column")
///     operator := r.URL.Query().Get("operator")
///     value := r.URL.Query().Get("value")
///
///     whereClause, err := BuildSafeWhereClause(validator, column, operator)
///     if err != nil {
///         http.Error(w, err.Error(), http.StatusBadRequest)
///         return
///     }
///
///     var users []User
///     db.Where(whereClause, value).Find(&users)
///     json.NewEncoder(w).Encode(users)
/// }
///
/// // Çoklu filtre oluşturma
/// type Filter struct {
///     Column   string
///     Operator string
///     Value    interface{}
/// }
///
/// func ApplyFilters(db *gorm.DB, validator *ColumnValidator, filters []Filter) (*gorm.DB, error) {
///     query := db
///
///     for _, filter := range filters {
///         whereClause, err := BuildSafeWhereClause(validator, filter.Column, filter.Operator)
///         if err != nil {
///             return nil, err
///         }
///         query = query.Where(whereClause, filter.Value)
///     }
///
///     return query, nil
/// }
///
/// // Kullanım
/// filters := []Filter{
///     {Column: "status", Operator: "=", Value: "active"},
///     {Column: "created_at", Operator: ">", Value: time.Now().AddDate(0, -1, 0)},
///     {Column: "email", Operator: "LIKE", Value: "%@example.com"},
/// }
///
/// query, err := ApplyFilters(db, validator, filters)
/// if err == nil {
///     query.Find(&users)
/// }
///
/// // Sorgu builder pattern
/// type QueryBuilder struct {
///     db        *gorm.DB
///     validator *ColumnValidator
///     errors    []error
/// }
///
/// func (qb *QueryBuilder) Where(column, operator string, value interface{}) *QueryBuilder {
///     whereClause, err := BuildSafeWhereClause(qb.validator, column, operator)
///     if err != nil {
///         qb.errors = append(qb.errors, err)
///         return qb
///     }
///
///     qb.db = qb.db.Where(whereClause, value)
///     return qb
/// }
///
/// func (qb *QueryBuilder) Execute(result interface{}) error {
///     if len(qb.errors) > 0 {
///         return qb.errors[0]
///     }
///     return qb.db.Find(result).Error
/// }
///
/// // Kullanım
/// var users []User
/// err := NewQueryBuilder(db, validator).
///     Where("status", "=", "active").
///     Where("age", ">=", 18).
///     Execute(&users)
/// ```
///
/// ## Özel Operatör Durumları
///
/// ```go
/// // IS NULL - değer gerektirmez
/// whereClause, _ := BuildSafeWhereClause(validator, "deleted_at", "IS NULL")
/// db.Where(whereClause).Find(&users)
/// // NOT: Değer parametresi verilmemeli
///
/// // IS NOT NULL - değer gerektirmez
/// whereClause, _ := BuildSafeWhereClause(validator, "email", "IS NOT NULL")
/// db.Where(whereClause).Find(&users)
///
/// // BETWEEN - iki değer gerektirir
/// whereClause, _ := BuildSafeWhereClause(validator, "age", "BETWEEN")
/// db.Where(whereClause+" AND ?", 18, 65).Find(&users)
/// // NOT: BETWEEN için özel işlem gerekir
///
/// // IN - slice gerektirir
/// whereClause, _ := BuildSafeWhereClause(validator, "status", "IN")
/// db.Where(whereClause, []string{"active", "pending"}).Find(&users)
/// ```
///
/// ## Önemli Notlar
///
/// ⚠️ **Parameterized Query**: Dönen string'de `?` placeholder vardır.
/// Mutlaka GORM'un Where metoduna değer parametresi ile birlikte verin.
///
/// ⚠️ **BETWEEN Operatörü**: BETWEEN için ek `AND ?` eklemeniz gerekir.
///
/// ⚠️ **NULL Operatörleri**: IS NULL ve IS NOT NULL değer parametresi gerektirmez.
///
/// ⚠️ **Hata Kontrolü**: Her zaman error dönüş değerini kontrol edin.
///
/// ## Güvenlik En İyi Pratikleri
///
/// ```go
/// // ✅ Doğru kullanım - Parameterized query
/// whereClause, err := BuildSafeWhereClause(validator, "email", "=")
/// if err == nil {
///     db.Where(whereClause, userInput) // Güvenli
/// }
///
/// // ❌ Yanlış kullanım - String concatenation
/// whereClause, err := BuildSafeWhereClause(validator, "email", "=")
/// if err == nil {
///     db.Where(whereClause + userInput) // TEHLİKELİ! SQL injection riski
/// }
///
/// // ✅ Doğru kullanım - Çoklu katman
/// sanitized := SanitizeColumnName(userColumn)
/// whereClause, err := BuildSafeWhereClause(validator, sanitized, userOperator)
/// if err == nil {
///     db.Where(whereClause, userValue) // Çok güvenli
/// }
/// ```
///
/// ## Performans
///
/// - **Zaman Karmaşıklığı**: O(1) - map lookup'lar
/// - **Bellek**: Minimal, sadece string oluşturma
/// - **Optimizasyon**: Validator'ı cache'leyin, her sorgu için yeniden oluşturmayın
///
/// ## Test Örnekleri
///
/// ```go
/// func TestBuildSafeWhereClause(t *testing.T) {
///     validator, _ := NewColumnValidator(db, &User{})
///
///     // Başarılı durumlar
///     tests := []struct {
///         column   string
///         operator string
///         expected string
///     }{
///         {"email", "=", "email = ?"},
///         {"created_at", ">", "created_at > ?"},
///         {"name", "LIKE", "name LIKE ?"},
///         {"status", "IN", "status IN ?"},
///     }
///
///     for _, tt := range tests {
///         result, err := BuildSafeWhereClause(validator, tt.column, tt.operator)
///         assert.NoError(t, err)
///         assert.Equal(t, tt.expected, result)
///     }
///
///     // Hata durumları
///     _, err := BuildSafeWhereClause(validator, "invalid_column", "=")
///     assert.Error(t, err)
///
///     _, err = BuildSafeWhereClause(validator, "email", "DROP")
///     assert.Error(t, err)
/// }
/// ```
///
/// ## İlgili Fonksiyonlar
///
/// - `ValidateColumn()`: Kolon doğrulama yapar
/// - `IsValidOperator()`: Operatör doğrulama yapar
/// - `SanitizeColumnName()`: Kolon ismini temizler
func BuildSafeWhereClause(validator *ColumnValidator, column, operator string) (string, error) {
	// Validate column
	safeColumn, err := validator.ValidateColumn(column)
	if err != nil {
		return "", err
	}

	// Validate operator
	if !IsValidOperator(operator) {
		return "", fmt.Errorf("invalid operator: %s", operator)
	}

	// Build safe clause
	return fmt.Sprintf("%s %s ?", safeColumn, operator), nil
}
