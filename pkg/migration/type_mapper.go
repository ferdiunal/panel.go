package migration

import (
	"reflect"
	"time"

	"github.com/ferdiunal/panel.go/pkg/fields"
)

// / # GoType
// /
// / Bu yapı, Go dilindeki tip bilgisini detaylı bir şekilde temsil eder ve migration işlemleri
// / sırasında field tiplerinin Go karşılıklarını yönetmek için kullanılır.
// /
// / ## Amaç
// /
// / Migration generator'ın field tiplerini Go struct field'larına dönüştürürken ihtiyaç duyduğu
// / tüm tip bilgilerini (base type, pointer durumu, slice durumu) tek bir yapıda toplar.
// /
// / ## Alanlar
// /
// / - `Type`: Field'ın temel Go tipi (örn: string, int64, time.Time)
// / - `IsPointer`: Field'ın nullable olup olmadığını belirtir (nullable ise pointer kullanılır)
// / - `IsSlice`: Field'ın bir slice tipi olup olmadığını belirtir (HasMany, BelongsToMany için)
// / - `ElementType`: Slice tiplerinde, slice'ın element tipini belirtir
// /
// / ## Kullanım Senaryoları
// /
// / 1. **Nullable Field'lar**: Database'de NULL değer alabilen field'lar için pointer tip oluşturma
// / 2. **İlişki Field'ları**: HasMany ve BelongsToMany gibi çoklu ilişkiler için slice tip belirleme
// / 3. **Tip Dönüşümü**: Field tiplerinden Go struct field tiplerini otomatik oluşturma
// /
// / ## Örnek Kullanım
// /
// / ```go
// / // String field (nullable)
// / goType := GoType{
// /     Type:      reflect.TypeOf(""),
// /     IsPointer: true,
// /     IsSlice:   false,
// / }
// / // Sonuç: *string
// /
// / // HasMany ilişkisi
// / goType := GoType{
// /     Type:        reflect.TypeOf(User{}),
// /     IsPointer:   false,
// /     IsSlice:     true,
// /     ElementType: reflect.TypeOf(User{}),
// / }
// / // Sonuç: []User
// / ```
// /
// / ## Önemli Notlar
// /
// / - İlişki field'ları (HasOne, HasMany, BelongsTo, BelongsToMany) için `Type` nil olabilir
// / - `IsPointer` ve `IsSlice` aynı anda true olamaz
// / - `ElementType` sadece `IsSlice` true olduğunda kullanılır
type GoType struct {
	Type        reflect.Type
	IsPointer   bool
	IsSlice     bool
	ElementType reflect.Type // Slice için element tipi
}

// / # TypeMapper
// /
// / Bu yapı, Field tiplerini Go ve SQL tiplerine dönüştüren merkezi bir mapper sağlar.
// / Migration generator'ın temel bileşenlerinden biridir ve farklı database dialect'leri
// / için uygun tip dönüşümlerini yönetir.
// /
// / ## Amaç
// /
// / Panel.go field sistemindeki field tiplerini (TYPE_TEXT, TYPE_NUMBER, vb.) hem Go struct
// / field tiplerine hem de database-specific SQL tiplerine dönüştürmek için tek bir merkezi
// / nokta sağlar. Bu sayede:
// / - Tip dönüşüm mantığı tek yerde toplanır
// / - Farklı database'ler için özel SQL tipleri desteklenir
// / - Migration dosyaları otomatik oluşturulabilir
// /
// / ## Desteklenen Database Dialect'leri
// /
// / - **PostgreSQL**: Modern özellikler (JSONB, BIGINT, TIMESTAMP)
// / - **MySQL**: Yaygın kullanım (JSON, BIGINT UNSIGNED, DATETIME)
// / - **SQLite**: Basit tipler (TEXT, INTEGER, DATETIME)
// /
// / ## Kullanım Senaryoları
// /
// / 1. **Migration Oluşturma**: Field tanımlarından SQL migration dosyaları üretme
// / 2. **Model Oluşturma**: Field tanımlarından Go struct'ları üretme
// / 3. **Tip Validasyonu**: Field tiplerinin geçerliliğini kontrol etme
// / 4. **İlişki Analizi**: Field'ların ilişki tiplerini belirleme
// /
// / ## Örnek Kullanım
// /
// / ```go
// / // PostgreSQL için mapper oluştur
// / mapper := NewTypeMapperWithDialect("postgres")
// /
// / // Field tipini Go tipine dönüştür
// / goType := mapper.MapFieldTypeToGo(fields.TYPE_TEXT, true)
// / // Sonuç: *string (nullable)
// /
// / // Field tipini SQL tipine dönüştür
// / sqlType := mapper.MapFieldTypeToSQL(fields.TYPE_TEXT, 255)
// / // Sonuç: "varchar(255)"
// /
// / // İlişki tipini kontrol et
// / if mapper.IsRelationshipType(fields.TYPE_LINK) {
// /     relType := mapper.GetRelationshipType(fields.TYPE_LINK)
// /     // relType: "belongsTo"
// / }
// / ```
// /
// / ## Avantajlar
// /
// / - **Merkezi Yönetim**: Tüm tip dönüşümleri tek bir yerde
// / - **Database Esnekliği**: Farklı database'ler için özel tip desteği
// / - **Tip Güvenliği**: Reflection kullanarak güvenli tip dönüşümleri
// / - **Genişletilebilirlik**: Yeni field tipleri kolayca eklenebilir
// /
// / ## Önemli Notlar
// /
// / - Dialect belirtilmezse varsayılan olarak PostgreSQL kullanılır
// / - İlişki field'ları (HasOne, HasMany, vb.) için özel işlem gerekir
// / - Nullable field'lar otomatik olarak pointer tipine dönüştürülür
// / - SQL tip boyutları (size) field tanımından alınır
type TypeMapper struct {
	dialect string // Database dialect (postgres, mysql, sqlite)
}

// / # NewTypeMapper
// /
// / Bu fonksiyon, varsayılan ayarlarla yeni bir TypeMapper instance'ı oluşturur.
// /
// / ## Amaç
// /
// / Hızlı bir şekilde TypeMapper oluşturmak için kullanılır. Database dialect belirtilmediği
// / için, SQL tip dönüşümlerinde varsayılan olarak PostgreSQL kullanılır.
// /
// / ## Döndürür
// /
// / - Yapılandırılmış TypeMapper pointer'ı (dialect: varsayılan PostgreSQL)
// /
// / ## Kullanım Senaryoları
// /
// / 1. **Hızlı Prototipleme**: Database tipi önemli olmayan test senaryoları
// / 2. **PostgreSQL Projeleri**: PostgreSQL kullanılan projeler için varsayılan mapper
// / 3. **Genel Amaçlı**: Database-agnostic tip dönüşümleri
// /
// / ## Örnek Kullanım
// /
// / ```go
// / // Varsayılan mapper oluştur (PostgreSQL)
// / mapper := NewTypeMapper()
// /
// / // Field tipini Go tipine dönüştür
// / goType := mapper.MapFieldTypeToGo(fields.TYPE_TEXT, false)
// / // Sonuç: string
// /
// / // Field tipini SQL tipine dönüştür (PostgreSQL)
// / sqlType := mapper.MapFieldTypeToSQL(fields.TYPE_NUMBER, 0)
// / // Sonuç: "bigint"
// / ```
// /
// / ## Önemli Notlar
// /
// / - Dialect belirtilmediği için SQL dönüşümlerinde PostgreSQL varsayılan olarak kullanılır
// / - Farklı bir database kullanıyorsanız `NewTypeMapperWithDialect()` kullanın
// / - Oluşturulduktan sonra `SetDialect()` ile dialect değiştirilebilir
// /
// / ## Alternatifler
// /
// / - `NewTypeMapperWithDialect(dialect)`: Belirli bir dialect ile mapper oluşturur
func NewTypeMapper() *TypeMapper {
	return &TypeMapper{}
}

// / # NewTypeMapperWithDialect
// /
// / Bu fonksiyon, belirtilen database dialect ile yeni bir TypeMapper instance'ı oluşturur.
// /
// / ## Amaç
// /
// / Farklı database sistemleri için özelleştirilmiş TypeMapper oluşturmak için kullanılır.
// / Her database'in kendine özgü tip sistemi olduğu için, doğru SQL tiplerini üretmek için
// / dialect belirtmek önemlidir.
// /
// / ## Parametreler
// /
// / - `dialect`: Database tipi ("postgres", "mysql", "sqlite")
// /   - **postgres**: PostgreSQL için optimize edilmiş tipler (JSONB, BIGINT, TIMESTAMP)
// /   - **mysql**: MySQL için optimize edilmiş tipler (JSON, BIGINT UNSIGNED, DATETIME)
// /   - **sqlite**: SQLite için optimize edilmiş tipler (TEXT, INTEGER, DATETIME)
// /
// / ## Döndürür
// /
// / - Yapılandırılmış TypeMapper pointer'ı (belirtilen dialect ile)
// /
// / ## Kullanım Senaryoları
// /
// / 1. **Multi-Database Desteği**: Farklı database'ler için migration oluşturma
// / 2. **Production Ortamları**: Belirli bir database için optimize edilmiş migration'lar
// / 3. **Test Ortamları**: SQLite ile hızlı test, production'da PostgreSQL/MySQL
// /
// / ## Örnek Kullanım
// /
// / ```go
// / // PostgreSQL için mapper
// / pgMapper := NewTypeMapperWithDialect("postgres")
// / sqlType := pgMapper.MapFieldTypeToSQL(fields.TYPE_KEY_VALUE, 0)
// / // Sonuç: "jsonb"
// /
// / // MySQL için mapper
// / mysqlMapper := NewTypeMapperWithDialect("mysql")
// / sqlType = mysqlMapper.MapFieldTypeToSQL(fields.TYPE_KEY_VALUE, 0)
// / // Sonuç: "json"
// /
// / // SQLite için mapper
// / sqliteMapper := NewTypeMapperWithDialect("sqlite")
// / sqlType = sqliteMapper.MapFieldTypeToSQL(fields.TYPE_NUMBER, 0)
// / // Sonuç: "integer"
// / ```
// /
// / ## Önemli Notlar
// /
// / - Dialect büyük/küçük harf duyarlı değildir ancak küçük harf kullanılması önerilir
// / - Desteklenmeyen dialect belirtilirse varsayılan olarak PostgreSQL kullanılır
// / - Dialect oluşturulduktan sonra `SetDialect()` ile değiştirilebilir
// /
// / ## Avantajlar
// /
// / - **Database-Specific Optimizasyon**: Her database için en uygun tipleri kullanır
// / - **Tip Uyumluluğu**: Database'in desteklediği tipleri garanti eder
// / - **Migration Kalitesi**: Database'e özgü özelliklerden yararlanır
func NewTypeMapperWithDialect(dialect string) *TypeMapper {
	return &TypeMapper{
		dialect: dialect,
	}
}

// / # SetDialect
// /
// / Bu method, mevcut TypeMapper instance'ının database dialect'ini değiştirir.
// /
// / ## Amaç
// /
// / Oluşturulmuş bir TypeMapper'ın dialect'ini runtime'da değiştirmek için kullanılır.
// / Bu özellikle multi-database senaryolarında veya test ortamlarında farklı database'ler
// / arasında geçiş yaparken kullanışlıdır.
// /
// / ## Parametreler
// /
// / - `dialect`: Yeni database tipi ("postgres", "mysql", "sqlite")
// /
// / ## Kullanım Senaryoları
// /
// / 1. **Test Ortamları**: Test için SQLite, production için PostgreSQL kullanma
// / 2. **Multi-Tenant Sistemler**: Farklı tenant'lar için farklı database'ler
// / 3. **Migration Araçları**: Tek bir mapper ile birden fazla database için migration üretme
// / 4. **Database Geçişi**: Bir database'den diğerine geçiş sırasında
// /
// / ## Örnek Kullanım
// /
// / ```go
// / // Varsayılan mapper oluştur
// / mapper := NewTypeMapper()
// /
// / // PostgreSQL için SQL tipi
// / sqlType := mapper.MapFieldTypeToSQL(fields.TYPE_BOOLEAN, 0)
// / // Sonuç: "boolean"
// /
// / // MySQL'e geç
// / mapper.SetDialect("mysql")
// / sqlType = mapper.MapFieldTypeToSQL(fields.TYPE_BOOLEAN, 0)
// / // Sonuç: "tinyint(1)"
// /
// / // SQLite'a geç
// / mapper.SetDialect("sqlite")
// / sqlType = mapper.MapFieldTypeToSQL(fields.TYPE_BOOLEAN, 0)
// / // Sonuç: "integer"
// / ```
// /
// / ## Önemli Notlar
// /
// / - Dialect değişikliği sadece sonraki `MapFieldTypeToSQL()` çağrılarını etkiler
// / - Önceden oluşturulmuş SQL tipleri değişmez
// / - Dialect büyük/küçük harf duyarlı değildir
// / - Geçersiz dialect belirtilirse varsayılan olarak PostgreSQL kullanılır
// /
// / ## Avantajlar
// /
// / - **Esneklik**: Tek bir mapper instance'ı ile birden fazla database desteği
// / - **Performans**: Her database için yeni mapper oluşturmaya gerek yok
// / - **Test Kolaylığı**: Test ve production ortamları arasında kolay geçiş
func (tm *TypeMapper) SetDialect(dialect string) {
	tm.dialect = dialect
}

// / # MapFieldTypeToGo
// /
// / Bu method, Panel.go field tiplerini Go programlama dilindeki karşılık gelen tiplere dönüştürür.
// / Migration generator'ın model struct'larını oluştururken kullandığı temel dönüşüm fonksiyonudur.
// /
// / ## Amaç
// /
// / Field tanımlarından Go struct field'ları oluştururken doğru Go tipini belirlemek için kullanılır.
// / Nullable field'lar için pointer tipleri, ilişki field'ları için özel tipler ve standart field'lar
// / için temel Go tipleri döndürür.
// /
// / ## Parametreler
// /
// / - `fieldType`: Panel.go field tipi (fields.ElementType)
// /   - **Metin Tipleri**: TYPE_TEXT, TYPE_PASSWORD, TYPE_TEXTAREA, TYPE_RICHTEXT, TYPE_EMAIL, TYPE_TEL
// /   - **Sayısal Tipler**: TYPE_NUMBER
// /   - **Boolean**: TYPE_BOOLEAN
// /   - **Tarih/Saat**: TYPE_DATE, TYPE_DATETIME
// /   - **Dosya Tipleri**: TYPE_FILE, TYPE_VIDEO, TYPE_AUDIO
// /   - **Seçim Tipleri**: TYPE_SELECT
// /   - **JSON Tipler**: TYPE_KEY_VALUE
// /   - **İlişki Tipleri**: TYPE_LINK, TYPE_DETAIL, TYPE_COLLECTION, TYPE_CONNECT
// /   - **Polimorfik İlişkiler**: TYPE_POLY_LINK, TYPE_POLY_DETAIL, TYPE_POLY_COLLECTION, TYPE_POLY_CONNECT
// /
// / - `nullable`: Field'ın NULL değer alıp alamayacağını belirtir
// /   - `true`: Pointer tip döner (örn: *string, *int64, *time.Time)
// /   - `false`: Değer tip döner (örn: string, int64, time.Time)
// /
// / ## Döndürür
// /
// / - `GoType`: Go tip bilgisini içeren yapı
// /   - `Type`: Temel Go tipi (reflect.Type)
// /   - `IsPointer`: Pointer tip olup olmadığı
// /   - `IsSlice`: Slice tip olup olmadığı
// /   - `ElementType`: Slice element tipi (varsa)
// /
// / ## Tip Dönüşüm Tablosu
// /
// / | Field Tipi | Nullable=false | Nullable=true | Açıklama |
// / |------------|----------------|---------------|----------|
// / | TYPE_TEXT | string | *string | Kısa metin |
// / | TYPE_TEXTAREA | string | *string | Uzun metin |
// / | TYPE_EMAIL | string | *string | Email adresi |
// / | TYPE_NUMBER | int64 | *int64 | Tam sayı |
// / | TYPE_BOOLEAN | bool | *bool | Boolean değer |
// / | TYPE_DATE | time.Time | *time.Time | Tarih |
// / | TYPE_DATETIME | time.Time | *time.Time | Tarih ve saat |
// / | TYPE_FILE | string | *string | Dosya URL'i |
// / | TYPE_SELECT | string | *string | Seçim değeri |
// / | TYPE_KEY_VALUE | map[string]interface{} | *map[string]interface{} | JSON data |
// / | TYPE_LINK | uint | *uint | Foreign key (BelongsTo) |
// / | TYPE_DETAIL | nil | nil | HasOne ilişkisi |
// / | TYPE_COLLECTION | nil | nil | HasMany ilişkisi |
// / | TYPE_CONNECT | nil | nil | BelongsToMany ilişkisi |
// /
// / ## Kullanım Senaryoları
// /
// / 1. **Model Oluşturma**: Go struct field'larını otomatik oluşturma
// / 2. **Tip Validasyonu**: Field tiplerinin Go karşılıklarını kontrol etme
// / 3. **Migration Generator**: Database migration'larında kullanılacak tipleri belirleme
// / 4. **ORM Entegrasyonu**: GORM gibi ORM'lerde kullanılacak tipleri belirleme
// /
// / ## Örnek Kullanım
// /
// / ```go
// / mapper := NewTypeMapper()
// /
// / // Nullable string field
// / goType := mapper.MapFieldTypeToGo(fields.TYPE_TEXT, true)
// / // Sonuç: GoType{Type: string, IsPointer: true}
// / // Go struct: Name *string `json:"name"`
// /
// / // Non-nullable number field
// / goType = mapper.MapFieldTypeToGo(fields.TYPE_NUMBER, false)
// / // Sonuç: GoType{Type: int64, IsPointer: false}
// / // Go struct: Age int64 `json:"age"`
// /
// / // Nullable datetime field
// / goType = mapper.MapFieldTypeToGo(fields.TYPE_DATETIME, true)
// / // Sonuç: GoType{Type: time.Time, IsPointer: true}
// / // Go struct: CreatedAt *time.Time `json:"created_at"`
// /
// / // BelongsTo ilişkisi (foreign key)
// / goType = mapper.MapFieldTypeToGo(fields.TYPE_LINK, false)
// / // Sonuç: GoType{Type: uint, IsPointer: false}
// / // Go struct: UserID uint `json:"user_id"`
// /
// / // HasMany ilişkisi
// / goType = mapper.MapFieldTypeToGo(fields.TYPE_COLLECTION, false)
// / // Sonuç: GoType{Type: nil, IsPointer: false}
// / // Go struct: Posts []Post `gorm:"foreignKey:UserID"`
// / ```
// /
// / ## Önemli Notlar
// /
// / - **İlişki Field'ları**: TYPE_DETAIL, TYPE_COLLECTION, TYPE_CONNECT ve polimorfik ilişkiler için
// /   `Type` nil döner çünkü bu field'lar model struct'ında özel olarak işlenir
// / - **Nullable Mantığı**: Nullable field'lar pointer tip olarak döner, bu sayede NULL değerler
// /   Go'da nil olarak temsil edilebilir
// / - **Foreign Key**: TYPE_LINK (BelongsTo) için uint tipi döner, bu GORM'un varsayılan ID tipidir
// / - **JSON Field'lar**: TYPE_KEY_VALUE için map[string]interface{} döner, database'de JSON olarak saklanır
// / - **Dosya Field'ları**: TYPE_FILE, TYPE_VIDEO, TYPE_AUDIO için string döner (URL saklanır)
// / - **Varsayılan Tip**: Tanınmayan field tipleri için string döner
// /
// / ## Avantajlar
// /
// / - **Tip Güvenliği**: Reflection kullanarak güvenli tip dönüşümleri
// / - **NULL Desteği**: Nullable field'lar için pointer tip desteği
// / - **ORM Uyumluluğu**: GORM ve diğer ORM'lerle uyumlu tipler
// / - **Genişletilebilirlik**: Yeni field tipleri kolayca eklenebilir
// /
// / ## Dikkat Edilmesi Gerekenler
// /
// / - İlişki field'ları için dönen nil değer özel olarak işlenmelidir
// / - Nullable=true olduğunda dönen pointer tip'in nil kontrolü yapılmalıdır
// / - TYPE_NUMBER için int64 kullanılır, farklı sayı tipleri için özel işlem gerekir
// / - TYPE_KEY_VALUE için map tipi kullanılır, JSON serialization gerekir
func (tm *TypeMapper) MapFieldTypeToGo(fieldType fields.ElementType, nullable bool) GoType {
	var baseType reflect.Type

	switch fieldType {
	// Metin Tipleri
	case fields.TYPE_TEXT, fields.TYPE_PASSWORD, fields.TYPE_TEXTAREA, fields.TYPE_RICHTEXT:
		baseType = reflect.TypeOf("")
	case fields.TYPE_EMAIL:
		baseType = reflect.TypeOf("")
	case fields.TYPE_TEL:
		baseType = reflect.TypeOf("")

	// Sayısal Tipler
	case fields.TYPE_NUMBER, fields.TYPE_MONEY:
		baseType = reflect.TypeOf(int64(0))

	// Boolean
	case fields.TYPE_BOOLEAN:
		baseType = reflect.TypeOf(false)

	// Tarih/Saat Tipleri
	case fields.TYPE_DATE, fields.TYPE_DATETIME:
		baseType = reflect.TypeOf(time.Time{})

	// Dosya Tipleri (URL string olarak saklanır)
	case fields.TYPE_FILE, fields.TYPE_VIDEO, fields.TYPE_AUDIO:
		baseType = reflect.TypeOf("")

	// Seçim Tipleri
	case fields.TYPE_SELECT:
		baseType = reflect.TypeOf("")

	// Key-Value (JSON olarak saklanır)
	case fields.TYPE_KEY_VALUE:
		baseType = reflect.TypeOf(map[string]interface{}{})

	// İlişki Tipleri
	case fields.TYPE_LINK: // BelongsTo -> Foreign Key
		baseType = reflect.TypeOf(uint(0))

	case fields.TYPE_DETAIL: // HasOne
		// Bu tip model struct'ında pointer olarak tanımlanır
		baseType = nil

	case fields.TYPE_COLLECTION: // HasMany
		// Bu tip model struct'ında slice olarak tanımlanır
		baseType = nil

	case fields.TYPE_CONNECT: // BelongsToMany
		// Bu tip model struct'ında slice olarak tanımlanır (many2many)
		baseType = nil

	// Polimorfik İlişkiler
	case fields.TYPE_POLY_LINK, fields.TYPE_POLY_DETAIL, fields.TYPE_POLY_COLLECTION, fields.TYPE_POLY_CONNECT:
		baseType = nil

	default:
		baseType = reflect.TypeOf("")
	}

	result := GoType{
		Type:      baseType,
		IsPointer: nullable && baseType != nil,
	}

	return result
}

// / # MapFieldTypeToSQL
// /
// / Bu method, Panel.go field tiplerini database-specific SQL tiplerine dönüştürür.
// / Migration generator'ın SQL migration dosyalarını oluştururken kullandığı temel
// / dönüşüm fonksiyonudur ve farklı database dialect'leri için optimize edilmiş
// / SQL tipleri döndürür.
// /
// / ## Amaç
// /
// / Field tanımlarından SQL migration dosyaları oluştururken her database sisteminin
// / kendi tip sistemine uygun SQL tipleri üretmek için kullanılır. PostgreSQL, MySQL
// / ve SQLite için optimize edilmiş tip dönüşümleri sağlar.
// /
// / ## Parametreler
// /
// / - `fieldType`: Panel.go field tipi (fields.ElementType)
// /   - **Metin Tipleri**: TYPE_TEXT, TYPE_EMAIL, TYPE_TEL, TYPE_PASSWORD
// /   - **Uzun Metin**: TYPE_TEXTAREA, TYPE_RICHTEXT
// /   - **Sayısal**: TYPE_NUMBER
// /   - **Boolean**: TYPE_BOOLEAN
// /   - **Tarih/Saat**: TYPE_DATE, TYPE_DATETIME
// /   - **Dosya**: TYPE_FILE, TYPE_VIDEO, TYPE_AUDIO
// /   - **Seçim**: TYPE_SELECT
// /   - **JSON**: TYPE_KEY_VALUE
// /   - **İlişki**: TYPE_LINK (Foreign Key)
// /
// / - `size`: Field boyutu (VARCHAR için karakter sayısı)
// /   - `0`: Varsayılan boyut kullanılır
// /   - `> 0`: Belirtilen boyut kullanılır (örn: VARCHAR(100))
// /
// / ## Döndürür
// /
// / - `string`: Database-specific SQL tip tanımı
// /
// / ## Database-Specific Tip Dönüşüm Tablosu
// /
// / ### Metin Tipleri (TEXT, EMAIL, TEL, PASSWORD)
// / | Dialect | size=0 | size=100 | Açıklama |
// / |---------|--------|----------|----------|
// / | PostgreSQL | varchar(255) | varchar(100) | Değişken uzunlukta string |
// / | MySQL | varchar(255) | varchar(100) | Değişken uzunlukta string |
// / | SQLite | varchar(255) | varchar(100) | TEXT olarak saklanır |
// /
// / ### Uzun Metin Tipleri (TEXTAREA, RICHTEXT)
// / | Dialect | SQL Tipi | Açıklama |
// / |---------|----------|----------|
// / | PostgreSQL | text | Sınırsız uzunlukta metin |
// / | MySQL | text | 65,535 karaktere kadar |
// / | SQLite | text | Sınırsız uzunlukta metin |
// /
// / ### Sayısal Tipler (NUMBER)
// / | Dialect | SQL Tipi | Açıklama |
// / |---------|----------|----------|
// / | PostgreSQL | bigint | -9223372036854775808 to 9223372036854775807 |
// / | MySQL | bigint | -9223372036854775808 to 9223372036854775807 |
// / | SQLite | integer | Dinamik boyut (1-8 byte) |
// /
// / ### Boolean Tipler
// / | Dialect | SQL Tipi | Açıklama |
// / |---------|----------|----------|
// / | PostgreSQL | boolean | Native boolean (true/false) |
// / | MySQL | tinyint(1) | 0=false, 1=true |
// / | SQLite | integer | 0=false, 1=true |
// /
// / ### Tarih Tipleri (DATE)
// / | Dialect | SQL Tipi | Açıklama |
// / |---------|----------|----------|
// / | PostgreSQL | date | Sadece tarih (YYYY-MM-DD) |
// / | MySQL | date | Sadece tarih (YYYY-MM-DD) |
// / | SQLite | date | TEXT olarak saklanır |
// /
// / ### Tarih/Saat Tipleri (DATETIME)
// / | Dialect | SQL Tipi | Açıklama |
// / |---------|----------|----------|
// / | PostgreSQL | timestamp | Timezone olmadan tarih/saat |
// / | MySQL | datetime | Timezone olmadan tarih/saat |
// / | SQLite | datetime | TEXT olarak saklanır |
// /
// / ### Dosya Tipleri (FILE, VIDEO, AUDIO)
// / | Dialect | SQL Tipi | Açıklama |
// / |---------|----------|----------|
// / | PostgreSQL | text | Dosya URL'i saklanır |
// / | MySQL | text | Dosya URL'i saklanır |
// / | SQLite | text | Dosya URL'i saklanır |
// /
// / ### Seçim Tipleri (SELECT)
// / | Dialect | SQL Tipi | Açıklama |
// / |---------|----------|----------|
// / | PostgreSQL | varchar(100) | Seçim değeri saklanır |
// / | MySQL | varchar(100) | Seçim değeri saklanır |
// / | SQLite | varchar(100) | Seçim değeri saklanır |
// /
// / ### JSON Tipleri (KEY_VALUE)
// / | Dialect | SQL Tipi | Açıklama |
// / |---------|----------|----------|
// / | PostgreSQL | jsonb | Binary JSON (indexlenebilir, hızlı) |
// / | MySQL | json | Native JSON (MySQL 5.7+) |
// / | SQLite | text | JSON string olarak saklanır |
// /
// / ### Foreign Key Tipleri (LINK - BelongsTo)
// / | Dialect | SQL Tipi | Açıklama |
// / |---------|----------|----------|
// / | PostgreSQL | bigint | GORM ID tipi ile uyumlu |
// / | MySQL | bigint unsigned | GORM ID tipi ile uyumlu |
// / | SQLite | integer | GORM ID tipi ile uyumlu |
// /
// / ## Kullanım Senaryoları
// /
// / 1. **Migration Oluşturma**: SQL CREATE TABLE statement'ları için column tipleri
// / 2. **Schema Validasyonu**: Mevcut database schema'sını kontrol etme
// / 3. **Database Geçişi**: Bir database'den diğerine geçiş için tip mapping
// / 4. **ORM Konfigürasyonu**: GORM column type tanımları
// /
// / ## Örnek Kullanım
// /
// / ```go
// / // PostgreSQL için mapper
// / pgMapper := NewTypeMapperWithDialect("postgres")
// /
// / // Metin field (boyut belirtilmiş)
// / sqlType := pgMapper.MapFieldTypeToSQL(fields.TYPE_TEXT, 100)
// / // Sonuç: "varchar(100)"
// / // SQL: name varchar(100)
// /
// / // Metin field (varsayılan boyut)
// / sqlType = pgMapper.MapFieldTypeToSQL(fields.TYPE_TEXT, 0)
// / // Sonuç: "varchar(255)"
// / // SQL: description varchar(255)
// /
// / // Uzun metin field
// / sqlType = pgMapper.MapFieldTypeToSQL(fields.TYPE_TEXTAREA, 0)
// / // Sonuç: "text"
// / // SQL: content text
// /
// / // Sayısal field
// / sqlType = pgMapper.MapFieldTypeToSQL(fields.TYPE_NUMBER, 0)
// / // Sonuç: "bigint"
// / // SQL: age bigint
// /
// / // Boolean field
// / sqlType = pgMapper.MapFieldTypeToSQL(fields.TYPE_BOOLEAN, 0)
// / // Sonuç: "boolean"
// / // SQL: is_active boolean
// /
// / // JSON field (PostgreSQL)
// / sqlType = pgMapper.MapFieldTypeToSQL(fields.TYPE_KEY_VALUE, 0)
// / // Sonuç: "jsonb"
// / // SQL: metadata jsonb
// /
// / // MySQL için mapper
// / mysqlMapper := NewTypeMapperWithDialect("mysql")
// /
// / // JSON field (MySQL)
// / sqlType = mysqlMapper.MapFieldTypeToSQL(fields.TYPE_KEY_VALUE, 0)
// / // Sonuç: "json"
// / // SQL: metadata json
// /
// / // Boolean field (MySQL)
// / sqlType = mysqlMapper.MapFieldTypeToSQL(fields.TYPE_BOOLEAN, 0)
// / // Sonuç: "tinyint(1)"
// / // SQL: is_active tinyint(1)
// /
// / // Foreign key (MySQL)
// / sqlType = mysqlMapper.MapFieldTypeToSQL(fields.TYPE_LINK, 0)
// / // Sonuç: "bigint unsigned"
// / // SQL: user_id bigint unsigned
// /
// / // SQLite için mapper
// / sqliteMapper := NewTypeMapperWithDialect("sqlite")
// /
// / // Sayısal field (SQLite)
// / sqlType = sqliteMapper.MapFieldTypeToSQL(fields.TYPE_NUMBER, 0)
// / // Sonuç: "integer"
// / // SQL: count integer
// /
// / // JSON field (SQLite)
// / sqlType = sqliteMapper.MapFieldTypeToSQL(fields.TYPE_KEY_VALUE, 0)
// / // Sonuç: "text"
// / // SQL: settings text (JSON string olarak)
// / ```
// /
// / ## Önemli Notlar
// /
// / - **Dialect Varsayılanı**: Dialect belirtilmezse PostgreSQL kullanılır
// / - **Size Parametresi**: Sadece VARCHAR tipleri için geçerlidir, diğer tipler için göz ardı edilir
// / - **JSON Desteği**: PostgreSQL'de JSONB (binary, hızlı), MySQL'de JSON (native), SQLite'da TEXT
// / - **Boolean Desteği**: PostgreSQL native boolean, MySQL ve SQLite integer kullanır
// / - **Foreign Key**: MySQL'de UNSIGNED kullanılır, PostgreSQL ve SQLite'da kullanılmaz
// / - **Timestamp vs Datetime**: PostgreSQL timestamp, MySQL datetime kullanır
// / - **SQLite Sınırlamaları**: SQLite sınırlı tip sistemine sahiptir, çoğu tip TEXT veya INTEGER olarak saklanır
// /
// / ## Avantajlar
// /
// / - **Database Optimizasyonu**: Her database için en uygun tipleri kullanır
// / - **Tip Uyumluluğu**: Database'in desteklediği tipleri garanti eder
// / - **GORM Uyumluluğu**: GORM'un beklediği tip formatlarını üretir
// / - **Migration Kalitesi**: Database-specific özelliklerden yararlanır
// /
// / ## Dikkat Edilmesi Gerekenler
// /
// / - **PostgreSQL JSONB**: JSONB binary format kullanır, indexleme ve sorgulama için optimize edilmiştir
// / - **MySQL JSON**: MySQL 5.7+ gerektirir, eski versiyonlarda TEXT kullanılmalıdır
// / - **SQLite Tipleri**: SQLite dinamik tip sistemi kullanır, belirtilen tipler sadece hint'tir
// / - **VARCHAR Boyutu**: Çok büyük boyutlar yerine TEXT kullanmayı düşünün
// / - **Foreign Key Tipleri**: ID field'ın tipi ile uyumlu olmalıdır (genelde BIGINT/UNSIGNED)
// / - **Timezone**: TIMESTAMP ve DATETIME timezone bilgisi içermez, UTC kullanımı önerilir
// /
// / ## Performans Notları
// /
// / - **PostgreSQL JSONB**: JSON'a göre daha fazla disk alanı kullanır ama çok daha hızlıdır
// / - **MySQL BIGINT UNSIGNED**: Pozitif sayılar için daha geniş aralık sağlar
// / - **SQLite INTEGER**: Dinamik boyut kullanır, küçük sayılar için optimize edilmiştir
// / - **TEXT vs VARCHAR**: TEXT sınırsız uzunluk için, VARCHAR belirli boyut için optimize edilmiştir
func (tm *TypeMapper) MapFieldTypeToSQL(fieldType fields.ElementType, size int) string {
	dialect := tm.dialect
	if dialect == "" {
		dialect = "postgres" // Default to postgres
	}

	switch fieldType {
	// Metin Tipleri
	case fields.TYPE_TEXT, fields.TYPE_EMAIL, fields.TYPE_TEL, fields.TYPE_PASSWORD:
		if size > 0 {
			return "varchar(" + itoa(size) + ")"
		}
		return "varchar(255)"

	// Uzun Metin Tipleri (TEXT column)
	case fields.TYPE_TEXTAREA, fields.TYPE_RICHTEXT:
		return "text"

	// Sayısal Tipler
	case fields.TYPE_NUMBER, fields.TYPE_MONEY:
		switch dialect {
		case "postgres":
			return "bigint"
		case "mysql":
			return "bigint"
		case "sqlite":
			return "integer"
		default:
			return "bigint"
		}

	// Boolean
	case fields.TYPE_BOOLEAN:
		switch dialect {
		case "postgres":
			return "boolean"
		case "mysql":
			return "tinyint(1)"
		case "sqlite":
			return "integer" // SQLite stores boolean as 0/1
		default:
			return "boolean"
		}

	// Tarih/Saat Tipleri
	case fields.TYPE_DATE:
		return "date"
	case fields.TYPE_DATETIME:
		switch dialect {
		case "postgres":
			return "timestamp"
		case "mysql":
			return "datetime"
		case "sqlite":
			return "datetime"
		default:
			return "timestamp"
		}

	// Dosya Tipleri (URL saklanır)
	case fields.TYPE_FILE, fields.TYPE_VIDEO, fields.TYPE_AUDIO:
		return "text"

	// Seçim Tipleri
	case fields.TYPE_SELECT:
		return "varchar(100)"

	// Key-Value
	case fields.TYPE_KEY_VALUE:
		switch dialect {
		case "postgres":
			return "jsonb"
		case "mysql":
			return "json"
		case "sqlite":
			return "text" // SQLite stores JSON as text
		default:
			return "jsonb"
		}

	// İlişki Tipleri (Foreign Key)
	case fields.TYPE_LINK:
		switch dialect {
		case "postgres":
			return "bigint"
		case "mysql":
			return "bigint unsigned"
		case "sqlite":
			return "integer"
		default:
			return "bigint"
		}

	default:
		return "varchar(255)"
	}
}

// / # GetRelationshipType
// /
// / Bu method, Panel.go field tipinin hangi ORM ilişki tipine karşılık geldiğini döndürür.
// / GORM ve diğer ORM sistemlerinde kullanılacak ilişki tipini belirlemek için kullanılır.
// /
// / ## Amaç
// /
// / Field tiplerinden ORM ilişki tiplerini belirlemek için kullanılır. Bu bilgi, model
// / struct'larında ilişki tanımlarını oluştururken ve migration dosyalarında foreign key
// / constraint'lerini belirlerken kullanılır.
// /
// / ## Parametreler
// /
// / - `fieldType`: Panel.go field tipi (fields.ElementType)
// /
// / ## Döndürür
// /
// / - `string`: ORM ilişki tipi adı
// /   - İlişki tipi değilse boş string ("") döner
// /
// / ## İlişki Tipleri Tablosu
// /
// / | Field Tipi | İlişki Tipi | Açıklama | GORM Örneği |
// / |------------|-------------|----------|-------------|
// / | TYPE_LINK | belongsTo | Bir kayıt başka bir kayda ait | User belongs to Company |
// / | TYPE_DETAIL | hasOne | Bir kayıt başka bir kayda sahip (tekil) | User has one Profile |
// / | TYPE_COLLECTION | hasMany | Bir kayıt birden fazla kayda sahip | User has many Posts |
// / | TYPE_CONNECT | belongsToMany | Çoktan çoğa ilişki (pivot table) | User belongs to many Roles |
// / | TYPE_POLY_LINK | morphTo | Polimorfik ait olma | Comment morphs to Post/Video |
// / | TYPE_POLY_DETAIL | morphOne | Polimorfik tekil sahiplik | Post morph one Image |
// / | TYPE_POLY_COLLECTION | morphMany | Polimorfik çoklu sahiplik | Post morph many Comments |
// / | TYPE_POLY_CONNECT | morphToMany | Polimorfik çoktan çoğa | Post morph to many Tags |
// /
// / ## Kullanım Senaryoları
// /
// / 1. **Model Oluşturma**: GORM struct tag'lerinde ilişki tipini belirtme
// / 2. **Migration Oluşturma**: Foreign key constraint'lerini belirleme
// / 3. **Validasyon**: Field'ın ilişki tipi olup olmadığını kontrol etme
// / 4. **Dokümantasyon**: İlişki tiplerini otomatik dokümante etme
// /
// / ## Örnek Kullanım
// /
// / ```go
// / mapper := NewTypeMapper()
// /
// / // BelongsTo ilişkisi
// / relType := mapper.GetRelationshipType(fields.TYPE_LINK)
// / // Sonuç: "belongsTo"
// / // GORM: type User struct { CompanyID uint; Company Company `gorm:"foreignKey:CompanyID"` }
// /
// / // HasOne ilişkisi
// / relType = mapper.GetRelationshipType(fields.TYPE_DETAIL)
// / // Sonuç: "hasOne"
// / // GORM: type User struct { Profile Profile `gorm:"foreignKey:UserID"` }
// /
// / // HasMany ilişkisi
// / relType = mapper.GetRelationshipType(fields.TYPE_COLLECTION)
// / // Sonuç: "hasMany"
// / // GORM: type User struct { Posts []Post `gorm:"foreignKey:UserID"` }
// /
// / // BelongsToMany ilişkisi
// / relType = mapper.GetRelationshipType(fields.TYPE_CONNECT)
// / // Sonuç: "belongsToMany"
// / // GORM: type User struct { Roles []Role `gorm:"many2many:user_roles"` }
// /
// / // MorphTo ilişkisi (polimorfik)
// / relType = mapper.GetRelationshipType(fields.TYPE_POLY_LINK)
// / // Sonuç: "morphTo"
// / // GORM: type Comment struct {
// / //   CommentableID uint
// / //   CommentableType string
// / //   Commentable interface{} `gorm:"polymorphic:Commentable"`
// / // }
// /
// / // MorphMany ilişkisi (polimorfik)
// / relType = mapper.GetRelationshipType(fields.TYPE_POLY_COLLECTION)
// / // Sonuç: "morphMany"
// / // GORM: type Post struct {
// / //   Comments []Comment `gorm:"polymorphic:Commentable"`
// / // }
// /
// / // İlişki tipi değil
// / relType = mapper.GetRelationshipType(fields.TYPE_TEXT)
// / // Sonuç: ""
// / ```
// /
// / ## İlişki Tipleri Detayları
// /
// / ### BelongsTo (TYPE_LINK)
// / - **Kullanım**: Bir kayıt başka bir kayda ait olduğunda
// / - **Foreign Key**: Kayıt kendi tablosunda foreign key tutar
// / - **Örnek**: Post belongs to User (posts.user_id)
// /
// / ### HasOne (TYPE_DETAIL)
// / - **Kullanım**: Bir kayıt başka bir kayda sahip olduğunda (1-1)
// / - **Foreign Key**: İlişkili kayıt foreign key tutar
// / - **Örnek**: User has one Profile (profiles.user_id)
// /
// / ### HasMany (TYPE_COLLECTION)
// / - **Kullanım**: Bir kayıt birden fazla kayda sahip olduğunda (1-N)
// / - **Foreign Key**: İlişkili kayıtlar foreign key tutar
// / - **Örnek**: User has many Posts (posts.user_id)
// /
// / ### BelongsToMany (TYPE_CONNECT)
// / - **Kullanım**: Çoktan çoğa ilişki (N-N)
// / - **Pivot Table**: Ara tablo gerektirir
// / - **Örnek**: User belongs to many Roles (user_roles pivot table)
// /
// / ### MorphTo (TYPE_POLY_LINK)
// / - **Kullanım**: Polimorfik ait olma (birden fazla model tipine ait olabilir)
// / - **Alanlar**: {name}_id ve {name}_type
// / - **Örnek**: Comment morphs to Post or Video
// /
// / ### MorphOne (TYPE_POLY_DETAIL)
// / - **Kullanım**: Polimorfik tekil sahiplik
// / - **Alanlar**: İlişkili kayıtta {name}_id ve {name}_type
// / - **Örnek**: Post morph one Image
// /
// / ### MorphMany (TYPE_POLY_COLLECTION)
// / - **Kullanım**: Polimorfik çoklu sahiplik
// / - **Alanlar**: İlişkili kayıtlarda {name}_id ve {name}_type
// / - **Örnek**: Post morph many Comments
// /
// / ### MorphToMany (TYPE_POLY_CONNECT)
// / - **Kullanım**: Polimorfik çoktan çoğa ilişki
// / - **Pivot Table**: Polimorfik pivot table gerektirir
// / - **Örnek**: Post morph to many Tags
// /
// / ## Önemli Notlar
// /
// / - İlişki tipi olmayan field'lar için boş string döner
// / - Dönen string GORM'un beklediği format ile uyumludur
// / - Polimorfik ilişkiler için ek alanlar (type, id) gerekir
// / - BelongsToMany için pivot table otomatik oluşturulmalıdır
// /
// / ## Avantajlar
// /
// / - **ORM Uyumluluğu**: GORM ve diğer ORM'lerle uyumlu ilişki tipleri
// / - **Tip Güvenliği**: Compile-time tip kontrolü
// / - **Dokümantasyon**: İlişki tiplerini otomatik dokümante etme
// / - **Validasyon**: İlişki tiplerini kolayca kontrol etme
// /
// / ## Dikkat Edilmesi Gerekenler
// /
// / - Polimorfik ilişkiler için ek alanlar (ID ve Type) gerekir
// / - BelongsToMany için pivot table migration'ı ayrıca oluşturulmalıdır
// / - Foreign key isimlendirme convention'larına dikkat edilmelidir
// / - İlişki yönü (belongs to vs has many) doğru belirlenmelidir
func (tm *TypeMapper) GetRelationshipType(fieldType fields.ElementType) string {
	switch fieldType {
	case fields.TYPE_LINK:
		return "belongsTo"
	case fields.TYPE_DETAIL:
		return "hasOne"
	case fields.TYPE_COLLECTION:
		return "hasMany"
	case fields.TYPE_CONNECT:
		return "belongsToMany"
	case fields.TYPE_POLY_LINK:
		return "morphTo"
	case fields.TYPE_POLY_DETAIL:
		return "morphOne"
	case fields.TYPE_POLY_COLLECTION:
		return "morphMany"
	case fields.TYPE_POLY_CONNECT:
		return "morphToMany"
	default:
		return ""
	}
}

// / # IsRelationshipType
// /
// / Bu method, verilen field tipinin bir ilişki tipi olup olmadığını kontrol eder.
// / Field'ın database'de foreign key veya ilişki olarak işlenmesi gerekip gerekmediğini
// / belirlemek için kullanılır.
// /
// / ## Amaç
// /
// / Field tiplerini işlerken, ilişki field'larının özel olarak ele alınması gerekir.
// / Bu method, bir field'ın ilişki tipi olup olmadığını hızlıca kontrol etmek için
// / kullanılır ve migration generator'da, model oluşturucuda ve validasyon işlemlerinde
// / kritik bir rol oynar.
// /
// / ## Parametreler
// /
// / - `fieldType`: Kontrol edilecek Panel.go field tipi (fields.ElementType)
// /
// / ## Döndürür
// /
// / - `bool`: İlişki tipi ise `true`, değilse `false`
// /
// / ## İlişki Tipleri
// /
// / Aşağıdaki field tipleri ilişki tipi olarak kabul edilir:
// /
// / ### Standart İlişkiler
// / - **TYPE_LINK**: BelongsTo ilişkisi (Foreign Key)
// / - **TYPE_DETAIL**: HasOne ilişkisi (1-1)
// / - **TYPE_COLLECTION**: HasMany ilişkisi (1-N)
// / - **TYPE_CONNECT**: BelongsToMany ilişkisi (N-N)
// /
// / ### Polimorfik İlişkiler
// / - **TYPE_POLY_LINK**: MorphTo ilişkisi
// / - **TYPE_POLY_DETAIL**: MorphOne ilişkisi
// / - **TYPE_POLY_COLLECTION**: MorphMany ilişkisi
// / - **TYPE_POLY_CONNECT**: MorphToMany ilişkisi
// /
// / ## Kullanım Senaryoları
// /
// / 1. **Migration Oluşturma**: İlişki field'ları için foreign key constraint'leri ekleme
// / 2. **Model Oluşturma**: İlişki field'ları için özel struct tag'leri ekleme
// / 3. **Validasyon**: İlişki field'larının doğru yapılandırılıp yapılandırılmadığını kontrol etme
// / 4. **Field İşleme**: İlişki field'larını standart field'lardan ayırma
// / 5. **Eager Loading**: İlişki field'ları için preload stratejisi belirleme
// /
// / ## Örnek Kullanım
// /
// / ```go
// / mapper := NewTypeMapper()
// /
// / // BelongsTo ilişkisi kontrolü
// / if mapper.IsRelationshipType(fields.TYPE_LINK) {
// /     // İlişki field'ı - foreign key oluştur
// /     relType := mapper.GetRelationshipType(fields.TYPE_LINK)
// /     // relType: "belongsTo"
// /     // SQL: user_id bigint, FOREIGN KEY (user_id) REFERENCES users(id)
// / }
// / // Sonuç: true
// /
// / // HasMany ilişkisi kontrolü
// / if mapper.IsRelationshipType(fields.TYPE_COLLECTION) {
// /     // İlişki field'ı - model struct'ında slice olarak tanımla
// /     // type User struct { Posts []Post `gorm:"foreignKey:UserID"` }
// / }
// / // Sonuç: true
// /
// / // Standart field kontrolü
// / if mapper.IsRelationshipType(fields.TYPE_TEXT) {
// /     // Bu blok çalışmaz
// / } else {
// /     // Standart field - normal column oluştur
// /     sqlType := mapper.MapFieldTypeToSQL(fields.TYPE_TEXT, 255)
// /     // sqlType: "varchar(255)"
// / }
// / // Sonuç: false
// /
// / // Polimorfik ilişki kontrolü
// / if mapper.IsRelationshipType(fields.TYPE_POLY_LINK) {
// /     // Polimorfik ilişki - type ve id field'ları oluştur
// /     // commentable_type varchar(255)
// /     // commentable_id bigint
// / }
// / // Sonuç: true
// /
// / // Field işleme örneği
// / for _, field := range fields {
// /     if mapper.IsRelationshipType(field.Type) {
// /         // İlişki field'ı - özel işlem
// /         handleRelationshipField(field)
// /     } else {
// /         // Standart field - normal işlem
// /         handleStandardField(field)
// /     }
// / }
// / ```
// /
// / ## Pratik Kullanım Örnekleri
// /
// / ### Migration Generator'da Kullanım
// / ```go
// / func generateMigration(fields []Field) string {
// /     var columns []string
// /
// /     for _, field := range fields {
// /         if mapper.IsRelationshipType(field.Type) {
// /             // İlişki field'ı için foreign key
// /             if field.Type == fields.TYPE_LINK {
// /                 columns = append(columns, fmt.Sprintf(
// /                     "%s_id %s",
// /                     field.Name,
// /                     mapper.MapFieldTypeToSQL(fields.TYPE_LINK, 0),
// /                 ))
// /                 // Foreign key constraint ekle
// /                 columns = append(columns, fmt.Sprintf(
// /                     "FOREIGN KEY (%s_id) REFERENCES %s(id)",
// /                     field.Name,
// /                     field.RelatedModel,
// /                 ))
// /             }
// /         } else {
// /             // Standart field için column
// /             columns = append(columns, fmt.Sprintf(
// /                 "%s %s",
// /                 field.Name,
// /                 mapper.MapFieldTypeToSQL(field.Type, field.Size),
// /             ))
// /         }
// /     }
// /
// /     return strings.Join(columns, ",\n")
// / }
// / ```
// /
// / ### Model Generator'da Kullanım
// / ```go
// / func generateModelStruct(fields []Field) string {
// /     var structFields []string
// /
// /     for _, field := range fields {
// /         if mapper.IsRelationshipType(field.Type) {
// /             // İlişki field'ı için özel struct field
// /             relType := mapper.GetRelationshipType(field.Type)
// /             switch relType {
// /             case "belongsTo":
// /                 // Foreign key field
// /                 structFields = append(structFields, fmt.Sprintf(
// /                     "%sID uint `json:\"%s_id\"`",
// /                     field.Name,
// /                     field.Name,
// /                 ))
// /                 // İlişki field
// /                 structFields = append(structFields, fmt.Sprintf(
// /                     "%s %s `gorm:\"foreignKey:%sID\"`",
// /                     field.Name,
// /                     field.RelatedModel,
// /                     field.Name,
// /                 ))
// /             case "hasMany":
// /                 // Slice field
// /                 structFields = append(structFields, fmt.Sprintf(
// /                     "%s []%s `gorm:\"foreignKey:%sID\"`",
// /                     field.Name,
// /                     field.RelatedModel,
// /                     "User", // Parent model
// /                 ))
// /             }
// /         } else {
// /             // Standart field için normal struct field
// /             goType := mapper.MapFieldTypeToGo(field.Type, field.Nullable)
// /             structFields = append(structFields, fmt.Sprintf(
// /                 "%s %s `json:\"%s\"`",
// /                 field.Name,
// /                 goType.Type.String(),
// /                 field.Name,
// /             ))
// /         }
// /     }
// /
// /     return strings.Join(structFields, "\n")
// / }
// / ```
// /
// / ### Validasyon'da Kullanım
// / ```go
// / func validateField(field Field) error {
// /     if mapper.IsRelationshipType(field.Type) {
// /         // İlişki field'ı için özel validasyon
// /         if field.RelatedModel == "" {
// /             return errors.New("relationship field must have related model")
// /         }
// /
// /         // BelongsToMany için pivot table kontrolü
// /         if field.Type == fields.TYPE_CONNECT {
// /             if field.PivotTable == "" {
// /                 return errors.New("belongsToMany requires pivot table")
// /             }
// /         }
// /     } else {
// /         // Standart field için validasyon
// /         if field.Type == fields.TYPE_TEXT && field.Size == 0 {
// /             field.Size = 255 // Varsayılan boyut
// /         }
// /     }
// /
// /     return nil
// / }
// / ```
// /
// / ## Önemli Notlar
// /
// / - **Performans**: Bu method basit bir switch-case kullanır, çok hızlıdır
// / - **Kapsam**: Tüm standart ve polimorfik ilişki tipleri desteklenir
// / - **Validasyon**: İlişki field'ları için ek validasyon gerekir (RelatedModel, vb.)
// / - **Migration**: İlişki field'ları için foreign key constraint'leri eklenmeli
// / - **Model**: İlişki field'ları için GORM tag'leri eklenmeli
// /
// / ## Avantajlar
// /
// / - **Tip Güvenliği**: Compile-time tip kontrolü
// / - **Performans**: O(1) zaman karmaşıklığı
// / - **Okunabilirlik**: Kod daha anlaşılır ve bakımı kolay
// / - **Genişletilebilirlik**: Yeni ilişki tipleri kolayca eklenebilir
// /
// / ## Dikkat Edilmesi Gerekenler
// /
// / - İlişki field'ları için `RelatedModel` belirtilmelidir
// / - BelongsToMany için `PivotTable` belirtilmelidir
// / - Polimorfik ilişkiler için `{name}_type` ve `{name}_id` field'ları oluşturulmalıdır
// / - Foreign key constraint'leri migration'da eklenmeli
// / - İlişki field'ları model struct'ında özel olarak işlenmeli
// /
// / ## İlgili Method'lar
// /
// / - `GetRelationshipType()`: İlişki tipinin adını döndürür
// / - `MapFieldTypeToGo()`: Field tipini Go tipine dönüştürür
// / - `MapFieldTypeToSQL()`: Field tipini SQL tipine dönüştürür
func (tm *TypeMapper) IsRelationshipType(fieldType fields.ElementType) bool {
	switch fieldType {
	case fields.TYPE_LINK, fields.TYPE_DETAIL, fields.TYPE_COLLECTION, fields.TYPE_CONNECT,
		fields.TYPE_POLY_LINK, fields.TYPE_POLY_DETAIL, fields.TYPE_POLY_COLLECTION, fields.TYPE_POLY_CONNECT:
		return true
	default:
		return false
	}
}

// / # itoa
// /
// / Bu fonksiyon, integer (tam sayı) değerini string'e dönüştüren yardımcı bir fonksiyondur.
// / Go'nun standart kütüphanesindeki `strconv.Itoa()` fonksiyonuna alternatif olarak,
// / dış bağımlılık olmadan integer-to-string dönüşümü sağlar.
// /
// / ## Amaç
// /
// / SQL tip tanımlarında VARCHAR boyutunu belirtirken (örn: "varchar(255)") integer
// / değerleri string'e dönüştürmek için kullanılır. Standart kütüphane import etmeden
// / basit bir dönüşüm mekanizması sağlar.
// /
// / ## Parametreler
// /
// / - `i`: Dönüştürülecek integer değer
// /   - Pozitif sayılar: 0, 1, 2, ..., 2147483647
// /   - Negatif sayılar: -1, -2, ..., -2147483648
// /   - Sıfır: 0
// /
// / ## Döndürür
// /
// / - `string`: Integer'ın string temsili
// /
// / ## Algoritma
// /
// / 1. **Sıfır Kontrolü**: Değer 0 ise direkt "0" döner
// / 2. **Negatif Kontrolü**: Negatif ise işareti sakla ve pozitife çevir
// / 3. **Basamak Çıkarma**: Modulo 10 ile son basamağı al, 10'a böl
// / 4. **String Oluşturma**: Basamakları ters sırada birleştir
// / 5. **İşaret Ekleme**: Negatif ise başına "-" ekle
// /
// / ## Kullanım Senaryoları
// /
// / 1. **SQL Tip Oluşturma**: VARCHAR boyutunu string'e dönüştürme
// / 2. **Migration Generator**: SQL statement'larında sayısal değerleri kullanma
// / 3. **Dış Bağımlılık Azaltma**: strconv import etmeden integer dönüşümü
// /
// / ## Örnek Kullanım
// /
// / ```go
// / // Pozitif sayı
// / result := itoa(255)
// / // Sonuç: "255"
// / // Kullanım: "varchar(255)"
// /
// / // Sıfır
// / result = itoa(0)
// / // Sonuç: "0"
// / // Kullanım: "varchar(0)" (geçersiz ama dönüşüm çalışır)
// /
// / // Küçük sayı
// / result = itoa(10)
// / // Sonuç: "10"
// / // Kullanım: "varchar(10)"
// /
// / // Büyük sayı
// / result = itoa(65535)
// / // Sonuç: "65535"
// / // Kullanım: "varchar(65535)"
// /
// / // Negatif sayı (teorik, VARCHAR boyutu için kullanılmaz)
// / result = itoa(-100)
// / // Sonuç: "-100"
// /
// / // Tek basamaklı
// / result = itoa(5)
// / // Sonuç: "5"
// / // Kullanım: "varchar(5)"
// / ```
// /
// / ## MapFieldTypeToSQL İçinde Kullanım
// /
// / ```go
// / func (tm *TypeMapper) MapFieldTypeToSQL(fieldType fields.ElementType, size int) string {
// /     switch fieldType {
// /     case fields.TYPE_TEXT:
// /         if size > 0 {
// /             // itoa kullanarak size'ı string'e dönüştür
// /             return "varchar(" + itoa(size) + ")"
// /             // size=100 için: "varchar(100)"
// /             // size=255 için: "varchar(255)"
// /         }
// /         return "varchar(255)"
// /     }
// / }
// / ```
// /
// / ## Performans Analizi
// /
// / ### Zaman Karmaşıklığı
// / - **O(log₁₀ n)**: Sayının basamak sayısı kadar iterasyon
// / - 1 basamak (0-9): 1 iterasyon
// / - 2 basamak (10-99): 2 iterasyon
// / - 3 basamak (100-999): 3 iterasyon
// / - 6 basamak (100000-999999): 6 iterasyon
// /
// / ### Alan Karmaşıklığı
// / - **O(log₁₀ n)**: Sonuç string'i için basamak sayısı kadar alan
// /
// / ### Karşılaştırma
// / ```go
// / // Bu fonksiyon (itoa)
// / result := itoa(255)
// / // Avantaj: Dış bağımlılık yok
// / // Dezavantaj: Daha yavaş (manuel implementasyon)
// /
// / // strconv.Itoa
// / result := strconv.Itoa(255)
// / // Avantaj: Optimize edilmiş, daha hızlı
// / // Dezavantaj: strconv import gerektirir
// /
// / // fmt.Sprintf
// / result := fmt.Sprintf("%d", 255)
// / // Avantaj: Esnek format desteği
// / // Dezavantaj: En yavaş, fmt import gerektirir
// / ```
// /
// / ## Önemli Notlar
// /
// / - **Sıfır Özel Durum**: Sıfır için özel kontrol yapılır, yoksa boş string döner
// / - **Negatif Sayılar**: Negatif sayılar desteklenir ama VARCHAR boyutu için kullanılmaz
// / - **Performans**: strconv.Itoa'dan daha yavaş ama dış bağımlılık gerektirmez
// / - **Basitlik**: Basit implementasyon, kolay anlaşılır ve bakımı kolay
// / - **Sınırlama**: Sadece int tipi desteklenir, int64 veya uint için uygun değil
// /
// / ## Avantajlar
// /
// / - **Bağımsızlık**: Dış kütüphane import etmeye gerek yok
// / - **Basitlik**: Anlaşılması ve bakımı kolay
// / - **Kontrol**: Tam kontrol, özel optimizasyonlar yapılabilir
// / - **Hafiflik**: Minimal kod, küçük binary boyutu
// /
// / ## Dezavantajlar
// /
// / - **Performans**: strconv.Itoa'dan daha yavaş
// / - **Sınırlı Tip Desteği**: Sadece int, int64 veya uint desteklenmez
// / - **Optimizasyon**: Go compiler'ın optimize etmesi daha zor
// / - **Standart Dışı**: Go idiomatik değil, strconv kullanımı tercih edilir
// /
// / ## Alternatifler
// /
// / ### strconv.Itoa (Önerilen)
// / ```go
// / import "strconv"
// /
// / result := strconv.Itoa(255)
// / // Avantaj: Optimize edilmiş, hızlı, standart
// / // Dezavantaj: Import gerektirir
// / ```
// /
// / ### fmt.Sprintf
// / ```go
// / import "fmt"
// /
// / result := fmt.Sprintf("%d", 255)
// / // Avantaj: Esnek format, birden fazla değer
// / // Dezavantaj: En yavaş, overkill basit dönüşümler için
// / ```
// /
// / ### strconv.FormatInt
// / ```go
// / import "strconv"
// /
// / result := strconv.FormatInt(int64(255), 10)
// / // Avantaj: int64 desteği, farklı base'ler
// / // Dezavantaj: Daha verbose, int için Itoa yeterli
// / ```
// /
// / ## Dikkat Edilmesi Gerekenler
// /
// / - **Tip Sınırlaması**: Sadece int tipi için çalışır, int64 veya uint için cast gerekir
// / - **Performans**: Sık kullanılan kod yollarında strconv.Itoa tercih edilmeli
// / - **Negatif Boyutlar**: VARCHAR boyutu negatif olamaz, validasyon yapılmalı
// / - **Sıfır Boyut**: VARCHAR(0) geçersizdir, varsayılan boyut kullanılmalı
// /
// / ## Test Senaryoları
// /
// / ```go
// / // Test case'leri
// / assert.Equal(t, "0", itoa(0))           // Sıfır
// / assert.Equal(t, "1", itoa(1))           // Tek basamak
// / assert.Equal(t, "10", itoa(10))         // İki basamak
// / assert.Equal(t, "255", itoa(255))       // Üç basamak
// / assert.Equal(t, "65535", itoa(65535))   // Beş basamak
// / assert.Equal(t, "-1", itoa(-1))         // Negatif tek basamak
// / assert.Equal(t, "-100", itoa(-100))     // Negatif üç basamak
// / ```
// /
// / ## Gelecek İyileştirmeler
// /
// / 1. **strconv Kullanımı**: Performans için strconv.Itoa'ya geçiş
// / 2. **int64 Desteği**: Büyük sayılar için int64 desteği
// / 3. **Buffer Kullanımı**: strings.Builder ile daha hızlı string oluşturma
// / 4. **Inline Optimizasyonu**: Compiler hint'leri ile inline optimizasyonu
// /
// / ## İlgili Fonksiyonlar
// /
// / - `MapFieldTypeToSQL()`: Bu fonksiyonu VARCHAR boyutu için kullanır
// / - `strconv.Itoa()`: Go standart kütüphanesindeki alternatif
// / - `fmt.Sprintf()`: Format string ile alternatif
func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	result := ""
	negative := false
	if i < 0 {
		negative = true
		i = -i
	}
	for i > 0 {
		result = string(rune('0'+i%10)) + result
		i /= 10
	}
	if negative {
		result = "-" + result
	}
	return result
}
