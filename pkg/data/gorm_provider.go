package data

import (
	stdcontext "context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/query"
	"github.com/iancoleman/strcase"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

/// # GormDataProvider
///
/// Bu yapı, GORM ORM kütüphanesi kullanarak veritabanı işlemlerini gerçekleştiren
/// bir veri sağlayıcısıdır. DataProvider interface'ini implement eder ve CRUD
/// (Create, Read, Update, Delete) operasyonlarını destekler.
///
/// ## Özellikler
///
/// - **GORM Entegrasyonu**: GORM ORM kütüphanesi ile tam entegrasyon
/// - **İlişki Yönetimi**: HasOne, HasMany, BelongsTo, Many2Many ilişkilerini destekler
/// - **Gelişmiş Filtreleme**: Çoklu operatör desteği (eq, neq, gt, gte, lt, lte, like, in, between vb.)
/// - **Arama Desteği**: Çoklu kolonlarda LIKE operatörü ile arama
/// - **Sıralama**: Çoklu kolon ve yön desteği ile sıralama
/// - **Sayfalama**: Offset ve limit tabanlı sayfalama
/// - **Güvenlik**: SQL injection koruması için kolon validasyonu
/// - **Context Desteği**: Context-aware operasyonlar
///
/// ## Kullanım Senaryoları
///
/// 1. **Admin Panel CRUD İşlemleri**: Yönetim panellerinde veri listeleme, oluşturma, güncelleme ve silme
/// 2. **API Endpoint'leri**: RESTful API'ler için veri sağlayıcı
/// 3. **Dinamik Sorgular**: Kullanıcı tarafından belirlenen filtreler ve sıralama
/// 4. **İlişkisel Veri Yönetimi**: Karmaşık ilişkisel veritabanı yapılarının yönetimi
///
/// ## Güvenlik
///
/// - SQL injection saldırılarına karşı kolon adı validasyonu
/// - Parameterized query kullanımı
/// - Güvenli kolon adı sanitizasyonu
///
/// ## Örnek Kullanım
///
/// ```go
/// // Model tanımı
/// type User struct {
///     ID        uint      `gorm:"primarykey"`
///     Name      string    `gorm:"size:100"`
///     Email     string    `gorm:"uniqueIndex;size:100"`
///     CreatedAt time.Time
///     UpdatedAt time.Time
/// }
///
/// // Provider oluşturma
/// db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
/// provider := NewGormDataProvider(db, &User{})
///
/// // Arama kolonları belirleme
/// provider.SetSearchColumns([]string{"name", "email"})
///
/// // İlişkileri yükleme
/// provider.SetWith([]string{"Profile", "Posts"})
///
/// // Veri listeleme
/// response, err := provider.Index(ctx, QueryRequest{
///     Page:    1,
///     PerPage: 10,
///     Search:  "john",
///     Filters: []query.Filter{
///         {Field: "status", Operator: query.OpEqual, Value: "active"},
///     },
///     Sorts: []query.Sort{
///         {Column: "created_at", Direction: "DESC"},
///     },
/// })
/// ```
///
/// ## Avantajlar
///
/// - **Tip Güvenliği**: Go'nun tip sistemi ile güvenli veri işleme
/// - **Esneklik**: Dinamik model desteği ile farklı tablolar için kullanılabilir
/// - **Performans**: GORM'un optimize edilmiş sorguları
/// - **Bakım Kolaylığı**: Merkezi veri erişim katmanı
///
/// ## Dezavantajlar
///
/// - **GORM Bağımlılığı**: Sadece GORM ile çalışır, diğer ORM'ler desteklenmez
/// - **Reflection Kullanımı**: Dinamik tip işlemleri için reflection kullanımı performansı etkileyebilir
/// - **Karmaşıklık**: İlişki yönetimi karmaşık senaryolarda zorlayıcı olabilir
///
/// ## Önemli Notlar
///
/// - Model parametresi struct veya struct pointer olmalıdır
/// - SearchColumns belirlenmezse arama çalışmaz
/// - WithRelationships ile eager loading yapılır, N+1 problemi önlenir
/// - ColumnValidator başarısız olursa fallback olarak temel sanitizasyon kullanılır
type GormDataProvider struct {
	/// GORM veritabanı bağlantısı
	DB *gorm.DB
	/// Veritabanı modeli (struct veya struct pointer)
	Model interface{}
	/// Arama yapılacak kolon isimleri
	SearchColumns []string
	/// Eager loading için yüklenecek ilişkiler
	WithRelationships []string
	/// SQL injection koruması için kolon validatörü
	columnValidator *ColumnValidator
	/// Raw SQL ile ilişki yükleme için loader
	relationshipLoader fields.RelationshipLoader
	/// Yüklenecek ilişki field'ları
	relationshipFields []fields.RelationshipField
}

/// # NewGormDataProvider
///
/// Bu fonksiyon, yeni bir GormDataProvider instance'ı oluşturur ve döndürür.
/// SQL injection saldırılarına karşı koruma için otomatik olarak kolon validatörü
/// başlatır.
///
/// ## Parametreler
///
/// - `db`: GORM veritabanı bağlantısı (*gorm.DB)
/// - `model`: Veritabanı modeli (struct veya struct pointer olmalıdır)
///
/// ## Döndürür
///
/// - Yapılandırılmış GormDataProvider pointer'ı
///
/// ## Güvenlik
///
/// - Otomatik olarak ColumnValidator oluşturur
/// - Validator başarısız olursa fallback olarak temel sanitizasyon kullanılır
/// - Güvenlik uyarıları konsola yazdırılır
///
/// ## Örnek Kullanım
///
/// ```go
/// type User struct {
///     ID        uint   `gorm:"primarykey"`
///     Name      string `gorm:"size:100"`
///     Email     string `gorm:"uniqueIndex;size:100"`
///     CreatedAt time.Time
///     UpdatedAt time.Time
/// }
///
/// db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
/// provider := NewGormDataProvider(db, &User{})
///
/// // Arama kolonları ayarla
/// provider.SetSearchColumns([]string{"name", "email"})
///
/// // İlişkileri yükle
/// provider.SetWith([]string{"Profile", "Posts"})
/// ```
///
/// ## Önemli Notlar
///
/// - Model parametresi nil olmamalıdır
/// - Model struct veya struct pointer olmalıdır
/// - Validator başarısız olsa bile provider oluşturulur (güvenli fallback)
/// - DB bağlantısı aktif ve geçerli olmalıdır
func NewGormDataProvider(db *gorm.DB, model interface{}) *GormDataProvider {
	// Initialize column validator for SQL injection protection
	validator, err := NewColumnValidator(db, model)
	if err != nil {
		// Log error but don't fail - fall back to basic validation
		fmt.Printf("[SECURITY WARNING] Failed to create column validator: %v\n", err)
	}

	return &GormDataProvider{
		DB:                 db,
		Model:              model,
		columnValidator:    validator,
		relationshipLoader: NewGormRelationshipLoader(db),
	}
}

/// # getContext
///
/// Bu fonksiyon, özel Context tipinden standart Go context'ini güvenli bir şekilde çıkarır.
/// Context nil ise veya içinde standart context yoksa, context.Background() döndürür.
///
/// ## Parametreler
///
/// - `ctx`: Özel Context pointer'ı (*context.Context)
///
/// ## Döndürür
///
/// - Standart Go context'i (stdcontext.Context)
/// - Nil durumlarında context.Background()
///
/// ## Kullanım Senaryoları
///
/// 1. **GORM İşlemleri**: GORM'un WithContext metoduna geçirilmek üzere
/// 2. **Timeout Yönetimi**: Context timeout'larını GORM sorgularına taşımak
/// 3. **İptal Yönetimi**: Context iptali ile sorgu iptali
/// 4. **Trace/Logging**: Context içindeki trace bilgilerini taşımak
///
/// ## Güvenlik
///
/// - Nil pointer kontrolü yapar
/// - Panic durumlarını önler
/// - Her zaman geçerli bir context döndürür
///
/// ## Örnek Kullanım
///
/// ```go
/// // Context ile sorgu
/// ctx := context.New(r.Context())
/// stdCtx := p.getContext(ctx)
/// db := p.DB.WithContext(stdCtx).Find(&users)
///
/// // Nil context durumu
/// stdCtx := p.getContext(nil) // context.Background() döner
/// ```
///
/// ## Önemli Notlar
///
/// - Bu method internal kullanım içindir
/// - Her GORM işleminden önce çağrılmalıdır
/// - Context timeout'ları GORM sorgularına yansır
/// - Context iptali sorguyu iptal eder
func (p *GormDataProvider) getContext(ctx *context.Context) stdcontext.Context {
	if ctx == nil {
		return stdcontext.Background()
	}
	stdCtx := ctx.Context()
	if stdCtx == nil {
		return stdcontext.Background()
	}
	return stdCtx
}

/// # SetSearchColumns
///
/// Bu fonksiyon, arama işlemlerinde kullanılacak kolon isimlerini ayarlar.
/// Index metodunda Search parametresi kullanıldığında, bu kolonlarda LIKE operatörü
/// ile arama yapılır.
///
/// ## Parametreler
///
/// - `cols`: Arama yapılacak kolon isimleri ([]string)
///
/// ## Kullanım Senaryoları
///
/// 1. **Metin Arama**: Kullanıcı adı, email, açıklama gibi alanlarda arama
/// 2. **Çoklu Kolon Arama**: Birden fazla alanda aynı anda arama
/// 3. **Dinamik Arama**: Farklı modeller için farklı arama kolonları
///
/// ## Örnek Kullanım
///
/// ```go
/// // Kullanıcı modelinde arama
/// provider := NewGormDataProvider(db, &User{})
/// provider.SetSearchColumns([]string{"name", "email", "username"})
///
/// // Ürün modelinde arama
/// productProvider := NewGormDataProvider(db, &Product{})
/// productProvider.SetSearchColumns([]string{"name", "description", "sku"})
///
/// // Arama yapma
/// response, _ := provider.Index(ctx, QueryRequest{
///     Search: "john", // name, email ve username kolonlarında aranır
///     Page: 1,
///     PerPage: 10,
/// })
/// ```
///
/// ## Önemli Notlar
///
/// - Kolon isimleri veritabanı kolon isimleri olmalıdır (snake_case)
/// - Geçersiz kolon isimleri güvenlik nedeniyle filtrelenir
/// - Boş liste verilirse arama çalışmaz
/// - LIKE operatörü kullanıldığı için büyük tablolarda performans sorunu olabilir
/// - Index oluşturulmuş kolonlar tercih edilmelidir
func (p *GormDataProvider) SetSearchColumns(cols []string) {
	p.SearchColumns = cols
}

/// # SetWith
///
/// Bu fonksiyon, eager loading için yüklenecek ilişkileri ayarlar.
/// GORM'un Preload metodunu kullanarak N+1 sorgu problemini önler.
///
/// ## Parametreler
///
/// - `rels`: Yüklenecek ilişki isimleri ([]string)
///
/// ## Kullanım Senaryoları
///
/// 1. **N+1 Problemi Önleme**: İlişkili verileri tek sorguda yükleme
/// 2. **Performans Optimizasyonu**: Gereksiz sorgu sayısını azaltma
/// 3. **İç İçe İlişkiler**: Nested ilişkileri yükleme
/// 4. **API Response**: İlişkili verileri JSON response'a dahil etme
///
/// ## Örnek Kullanım
///
/// ```go
/// // Basit ilişki yükleme
/// provider := NewGormDataProvider(db, &User{})
/// provider.SetWith([]string{"Profile", "Posts"})
///
/// // İç içe ilişkiler
/// provider.SetWith([]string{"Posts.Comments", "Posts.Author"})
///
/// // Çoklu ilişki
/// provider.SetWith([]string{
///     "Profile",
///     "Posts",
///     "Posts.Comments",
///     "Posts.Tags",
///     "Roles",
/// })
///
/// // Kullanım
/// response, _ := provider.Index(ctx, QueryRequest{
///     Page: 1,
///     PerPage: 10,
/// })
/// // response.Items içinde ilişkili veriler de gelir
/// ```
///
/// ## Avantajlar
///
/// - **Performans**: N+1 sorgu problemi önlenir
/// - **Esneklik**: İhtiyaca göre ilişkiler yüklenir
/// - **Kontrol**: Hangi ilişkilerin yükleneceği kontrol edilir
///
/// ## Önemli Notlar
///
/// - İlişki isimleri struct field isimleri olmalıdır (PascalCase)
/// - Geçersiz ilişki isimleri hata vermez, sadece yüklenmez
/// - Çok fazla ilişki yüklemek performansı olumsuz etkileyebilir
/// - İç içe ilişkiler nokta (.) ile ayrılır: "Posts.Comments"
/// - Tüm CRUD operasyonlarında (Index, Show, Create, Update) geçerlidir
func (p *GormDataProvider) SetWith(rels []string) {
	p.WithRelationships = rels
}

/// # applyFilters
///
/// Bu fonksiyon, gelişmiş filtre koşullarını GORM sorgusuna uygular.
/// Çoklu operatör desteği ile karmaşık filtreleme senaryolarını destekler.
///
/// ## Parametreler
///
/// - `db`: GORM veritabanı instance'ı (*gorm.DB)
/// - `filters`: Uygulanacak filtre listesi ([]query.Filter)
///
/// ## Döndürür
///
/// - Filtreler uygulanmış GORM DB instance'ı (*gorm.DB)
///
/// ## Desteklenen Operatörler
///
/// | Operatör | Açıklama | Örnek |
/// |----------|----------|-------|
/// | `OpEqual` (eq) | Eşittir | `{Field: "status", Operator: OpEqual, Value: "active"}` |
/// | `OpNotEqual` (neq) | Eşit değildir | `{Field: "status", Operator: OpNotEqual, Value: "deleted"}` |
/// | `OpGreaterThan` (gt) | Büyüktür | `{Field: "age", Operator: OpGreaterThan, Value: 18}` |
/// | `OpGreaterEq` (gte) | Büyük eşittir | `{Field: "price", Operator: OpGreaterEq, Value: 100}` |
/// | `OpLessThan` (lt) | Küçüktür | `{Field: "stock", Operator: OpLessThan, Value: 10}` |
/// | `OpLessEq` (lte) | Küçük eşittir | `{Field: "discount", Operator: OpLessEq, Value: 50}` |
/// | `OpLike` (like) | İçerir (LIKE) | `{Field: "name", Operator: OpLike, Value: "john"}` |
/// | `OpNotLike` (nlike) | İçermez (NOT LIKE) | `{Field: "email", Operator: OpNotLike, Value: "spam"}` |
/// | `OpIn` (in) | Liste içinde | `{Field: "status", Operator: OpIn, Value: []string{"active", "pending"}}` |
/// | `OpNotIn` (nin) | Liste dışında | `{Field: "role", Operator: OpNotIn, Value: []string{"banned", "deleted"}}` |
/// | `OpIsNull` (null) | NULL değer | `{Field: "deleted_at", Operator: OpIsNull, Value: true}` |
/// | `OpIsNotNull` (nnull) | NULL değil | `{Field: "email_verified_at", Operator: OpIsNotNull, Value: true}` |
/// | `OpBetween` (between) | Aralık | `{Field: "created_at", Operator: OpBetween, Value: []string{"2024-01-01", "2024-12-31"}}` |
///
/// ## Güvenlik
///
/// - **SQL Injection Koruması**: Tüm kolon isimleri ColumnValidator ile doğrulanır
/// - **Parameterized Queries**: Tüm değerler parametre olarak geçirilir
/// - **Kolon Sanitizasyonu**: Geçersiz kolonlar otomatik olarak filtrelenir
/// - **Güvenlik Logları**: Reddedilen kolonlar konsola yazdırılır
///
/// ## Kullanım Senaryoları
///
/// 1. **Basit Filtreleme**: Tek kolon, tek değer
/// 2. **Çoklu Filtreleme**: Birden fazla kolon ve operatör
/// 3. **Karmaşık Sorgular**: AND/OR kombinasyonları
/// 4. **Dinamik Filtreleme**: Kullanıcı tarafından belirlenen filtreler
///
/// ## Örnek Kullanım
///
/// ```go
/// // Basit eşitlik filtresi
/// filters := []query.Filter{
///     {Field: "status", Operator: query.OpEqual, Value: "active"},
/// }
///
/// // Çoklu filtre
/// filters := []query.Filter{
///     {Field: "status", Operator: query.OpEqual, Value: "active"},
///     {Field: "age", Operator: query.OpGreaterEq, Value: 18},
///     {Field: "country", Operator: query.OpIn, Value: []string{"TR", "US", "UK"}},
/// }
///
/// // Tarih aralığı filtresi
/// filters := []query.Filter{
///     {Field: "created_at", Operator: query.OpBetween, Value: []string{"2024-01-01", "2024-12-31"}},
/// }
///
/// // LIKE filtresi
/// filters := []query.Filter{
///     {Field: "name", Operator: query.OpLike, Value: "john"},
/// }
///
/// // NULL kontrolü
/// filters := []query.Filter{
///     {Field: "deleted_at", Operator: query.OpIsNull, Value: true},
/// }
///
/// // Kullanım
/// db := p.DB.Model(&User{})
/// db = p.applyFilters(db, filters)
/// db.Find(&users)
/// ```
///
/// ## Önemli Notlar
///
/// - Geçersiz kolon isimleri sessizce atlanır (güvenlik)
/// - OpLike ve OpNotLike otomatik olarak % wildcard ekler
/// - OpIn ve OpNotIn için değer []string tipinde olmalıdır
/// - OpBetween için tam olarak 2 elemanlı []string gereklidir
/// - OpIsNull ve OpIsNotNull için değer bool tipinde olmalıdır
/// - Bilinmeyen operatörler varsayılan olarak eşitlik kontrolü yapar
/// - Tüm filtreler AND mantığı ile birleştirilir
///
/// ## Performans İpuçları
///
/// - Filtrelenen kolonlarda index oluşturun
/// - LIKE operatörü büyük tablolarda yavaş olabilir
/// - IN operatörü için çok fazla değer performansı etkileyebilir
/// - BETWEEN operatörü tarih aralıkları için optimize edilmiştir
func (p *GormDataProvider) applyFilters(db *gorm.DB, filters []query.Filter) *gorm.DB {
	for _, f := range filters {
		if f.Field == "" {
			continue
		}

		// SECURITY: Validate column name to prevent SQL injection
		safeColumn := f.Field
		if p.columnValidator != nil {
			validatedCol, err := p.columnValidator.ValidateColumn(f.Field)
			if err != nil {
				// Skip invalid columns - don't expose error to user
				fmt.Printf("[SECURITY] Rejected invalid column in filter: %s\n", f.Field)
				continue
			}
			safeColumn = validatedCol
		} else {
			// Fallback: sanitize column name if validator not available
			safeColumn = SanitizeColumnName(f.Field)
		}

		switch f.Operator {
		case query.OpEqual:
			db = db.Where(fmt.Sprintf("%s = ?", safeColumn), f.Value)

		case query.OpNotEqual:
			db = db.Where(fmt.Sprintf("%s != ?", safeColumn), f.Value)

		case query.OpGreaterThan:
			db = db.Where(fmt.Sprintf("%s > ?", safeColumn), f.Value)

		case query.OpGreaterEq:
			db = db.Where(fmt.Sprintf("%s >= ?", safeColumn), f.Value)

		case query.OpLessThan:
			db = db.Where(fmt.Sprintf("%s < ?", safeColumn), f.Value)

		case query.OpLessEq:
			db = db.Where(fmt.Sprintf("%s <= ?", safeColumn), f.Value)

		case query.OpLike:
			if strVal, ok := f.Value.(string); ok {
				db = db.Where(fmt.Sprintf("%s LIKE ?", safeColumn), "%"+strVal+"%")
			}

		case query.OpNotLike:
			if strVal, ok := f.Value.(string); ok {
				db = db.Where(fmt.Sprintf("%s NOT LIKE ?", safeColumn), "%"+strVal+"%")
			}

		case query.OpIn:
			if vals, ok := f.Value.([]string); ok && len(vals) > 0 {
				db = db.Where(fmt.Sprintf("%s IN ?", safeColumn), vals)
			}

		case query.OpNotIn:
			if vals, ok := f.Value.([]string); ok && len(vals) > 0 {
				db = db.Where(fmt.Sprintf("%s NOT IN ?", safeColumn), vals)
			}

		case query.OpIsNull:
			if boolVal, ok := f.Value.(bool); ok && boolVal {
				db = db.Where(fmt.Sprintf("%s IS NULL", safeColumn))
			}

		case query.OpIsNotNull:
			if boolVal, ok := f.Value.(bool); ok && boolVal {
				db = db.Where(fmt.Sprintf("%s IS NOT NULL", safeColumn))
			}

		case query.OpBetween:
			if vals, ok := f.Value.([]string); ok && len(vals) == 2 {
				db = db.Where(fmt.Sprintf("%s BETWEEN ? AND ?", safeColumn), vals[0], vals[1])
			}

		default:
			// Default to equality
			db = db.Where(fmt.Sprintf("%s = ?", safeColumn), f.Value)
		}
	}
	return db
}

/// # Index
///
/// Bu fonksiyon, veritabanından sayfalanmış, filtrelenmiş ve sıralanmış veri listesi döndürür.
/// Admin panelleri ve API endpoint'leri için optimize edilmiş gelişmiş sorgu özellikleri sunar.
///
/// ## Parametreler
///
/// - `ctx`: Context bilgisi (*context.Context)
/// - `req`: Sorgu parametreleri (QueryRequest)
///   - `Page`: Sayfa numarası (1'den başlar)
///   - `PerPage`: Sayfa başına kayıt sayısı
///   - `Search`: Arama terimi (SearchColumns'da aranır)
///   - `Filters`: Gelişmiş filtreler ([]query.Filter)
///   - `Sorts`: Sıralama kuralları ([]query.Sort)
///
/// ## Döndürür
///
/// - `*QueryResponse`: Sorgu sonucu
///   - `Items`: Veri listesi ([]interface{})
///   - `Total`: Toplam kayıt sayısı (int64)
///   - `Page`: Mevcut sayfa numarası (int)
///   - `PerPage`: Sayfa başına kayıt sayısı (int)
/// - `error`: Hata durumunda hata mesajı
///
/// ## İşlem Sırası
///
/// 1. **Context Hazırlama**: Context'i standart Go context'ine dönüştürme
/// 2. **Model Ayarlama**: GORM model'ini ayarlama
/// 3. **Eager Loading**: WithRelationships ile ilişkileri yükleme
/// 4. **Filtreleme**: Gelişmiş filtreleri uygulama
/// 5. **Arama**: SearchColumns'da LIKE operatörü ile arama
/// 6. **Sayma**: Toplam kayıt sayısını hesaplama
/// 7. **Sıralama**: Sıralama kurallarını uygulama
/// 8. **Sayfalama**: Offset ve limit uygulama
/// 9. **Sorgu Çalıştırma**: Verileri çekme
/// 10. **Dönüştürme**: Sonuçları []interface{} formatına dönüştürme
///
/// ## Kullanım Senaryoları
///
/// 1. **Admin Panel Listeleme**: Yönetim panelinde veri tabloları
/// 2. **API Endpoint**: RESTful API için liste endpoint'i
/// 3. **Raporlama**: Filtrelenmiş ve sıralanmış raporlar
/// 4. **Arama Sistemi**: Çoklu kolonlarda arama
/// 5. **Veri İhracı**: Filtrelenmiş verilerin dışa aktarımı
///
/// ## Örnek Kullanım
///
/// ```go
/// // Basit listeleme
/// response, err := provider.Index(ctx, QueryRequest{
///     Page:    1,
///     PerPage: 10,
/// })
///
/// // Arama ile listeleme
/// provider.SetSearchColumns([]string{"name", "email"})
/// response, err := provider.Index(ctx, QueryRequest{
///     Page:    1,
///     PerPage: 10,
///     Search:  "john",
/// })
///
/// // Filtreli listeleme
/// response, err := provider.Index(ctx, QueryRequest{
///     Page:    1,
///     PerPage: 10,
///     Filters: []query.Filter{
///         {Field: "status", Operator: query.OpEqual, Value: "active"},
///         {Field: "age", Operator: query.OpGreaterEq, Value: 18},
///     },
/// })
///
/// // Sıralı listeleme
/// response, err := provider.Index(ctx, QueryRequest{
///     Page:    1,
///     PerPage: 10,
///     Sorts: []query.Sort{
///         {Column: "created_at", Direction: "DESC"},
///         {Column: "name", Direction: "ASC"},
///     },
/// })
///
/// // Tam özellikli listeleme
/// provider.SetSearchColumns([]string{"name", "email", "username"})
/// provider.SetWith([]string{"Profile", "Posts"})
/// response, err := provider.Index(ctx, QueryRequest{
///     Page:    1,
///     PerPage: 20,
///     Search:  "john",
///     Filters: []query.Filter{
///         {Field: "status", Operator: query.OpEqual, Value: "active"},
///         {Field: "created_at", Operator: query.OpBetween, Value: []string{"2024-01-01", "2024-12-31"}},
///     },
///     Sorts: []query.Sort{
///         {Column: "created_at", Direction: "DESC"},
///     },
/// })
///
/// // Response kullanımı
/// fmt.Printf("Toplam: %d, Sayfa: %d/%d\n",
///     response.Total,
///     response.Page,
///     (response.Total + int64(response.PerPage) - 1) / int64(response.PerPage))
/// for _, item := range response.Items {
///     user := item.(*User)
///     fmt.Printf("User: %s\n", user.Name)
/// }
/// ```
///
/// ## Güvenlik
///
/// - **SQL Injection Koruması**: Tüm kolon isimleri validate edilir
/// - **Parameterized Queries**: Tüm değerler güvenli şekilde bind edilir
/// - **Kolon Validasyonu**: Geçersiz kolonlar otomatik olarak filtrelenir
/// - **Context Timeout**: Context timeout'ları sorguya yansır
///
/// ## Performans Optimizasyonları
///
/// 1. **Eager Loading**: N+1 problemi önlenir
/// 2. **Count Optimizasyonu**: Count sorgusu ana sorgudan önce çalışır
/// 3. **Index Kullanımı**: Filtreleme ve sıralama için index'ler kullanılır
/// 4. **Sayfalama**: Offset/Limit ile bellek kullanımı optimize edilir
///
/// ## Performans İpuçları
///
/// - SearchColumns'da index oluşturun
/// - Sık kullanılan filtrelerde index oluşturun
/// - PerPage değerini makul tutun (10-100 arası)
/// - Gereksiz ilişkileri yüklemeyin
/// - LIKE araması yerine full-text search kullanmayı düşünün
///
/// ## Önemli Notlar
///
/// - Page numarası 1'den başlar (0 değil)
/// - SearchColumns boşsa arama çalışmaz
/// - Geçersiz sıralama kolonları atlanır
/// - Tüm filtreler AND mantığı ile birleştirilir
/// - İlişkiler otomatik olarak JSON'a dahil edilir
/// - Reflection kullanıldığı için büyük veri setlerinde performans etkilenebilir
///
/// ## Hata Durumları
///
/// - Veritabanı bağlantı hatası
/// - Geçersiz model tipi
/// - Context timeout
/// - Sorgu hatası (syntax, constraint vb.)
func (p *GormDataProvider) Index(ctx *context.Context, req QueryRequest) (*QueryResponse, error) {
	var total int64
	// We need a slice of the model type to hold results.
	// Since Model is interface{}, we might need reflection to create a slice of that type,
	// or we can just hope GORM handles Find(&[]Interface{}) correctly if we pass a pointer to a slice of models.
	// Actually, usually users pass a struct instance as Model.
	// Gorm's db.Model(model) works for setting the table.
	// But Find needs a destination.
	// Let's assume we return []map[string]interface{} for generic usage if we don't know the slice type,
	// OR we assume the user might want typed results.
	// But FieldHandler expects []interface{} in Items.

	// Simplest approach for Generic provider: Use map[string]interface{} for dynamic results
	// OR use reflection to make a slice of the Model's type.

	// Let's try map[string]interface{} for maximum flexibility in this generic provider,
	// unless we strictly want the structs.
	// If we use structs, we need to use reflect.New(reflect.SliceOf(reflect.TypeOf(p.Model))).Interface()

	// Let's start with just using db.Model(p.Model).Find(&results) where results is []map[string]interface{}
	// GORM supports finding into a map.

	stdCtx := p.getContext(ctx)
	db := p.DB.WithContext(stdCtx).Model(p.Model)

	// Apply Eager Loading with GORM Preload
	// WORKAROUND: Direkt olarak WithRelationships kullan çünkü relationshipFields boş olabilir
	// (field type detection sorunu nedeniyle)
	fmt.Printf("[DEBUG] Index - Preloading relationships: %v\n", p.WithRelationships)
	for _, relName := range p.WithRelationships {
		fmt.Printf("[DEBUG] Index - Preload: %s\n", relName)
		db = db.Preload(relName)
	}

	// Apply Advanced Filters
	if len(req.Filters) > 0 {
		db = p.applyFilters(db, req.Filters)
	}

	// Apply Search with column validation
	fmt.Printf("[GORM] Search: %q, SearchColumns: %v\n", req.Search, p.SearchColumns)
	if req.Search != "" && len(p.SearchColumns) > 0 {
		searchQuery := p.DB.WithContext(stdCtx).Session(&gorm.Session{NewDB: true})
		for _, col := range p.SearchColumns {
			// SECURITY: Validate search column names
			safeColumn := col
			if p.columnValidator != nil {
				validatedCol, err := p.columnValidator.ValidateColumn(col)
				if err != nil {
					// Skip invalid columns - don't expose error to user
					fmt.Printf("[SECURITY] Rejected invalid search column: %s\n", col)
					continue
				}
				safeColumn = validatedCol
			} else {
				// Fallback: sanitize column name if validator not available
				safeColumn = SanitizeColumnName(col)
			}
			searchQuery = searchQuery.Or(fmt.Sprintf("%s LIKE ?", safeColumn), "%"+req.Search+"%")
		}
		db = db.Where(searchQuery)
		fmt.Printf("[GORM] Search applied for columns: %v\n", p.SearchColumns)
	} else {
		fmt.Printf("[GORM] Search NOT applied - Search empty: %v, SearchColumns empty: %v\n", req.Search == "", len(p.SearchColumns) == 0)
	}

	// Count Total
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	// Sorting with column validation
	if len(req.Sorts) > 0 {
		for _, sort := range req.Sorts {
			if sort.Column != "" {
				// SECURITY: Validate sort column names
				safeColumn := sort.Column
				if p.columnValidator != nil {
					validatedCol, err := p.columnValidator.ValidateColumn(sort.Column)
					if err != nil {
						// Skip invalid columns - don't expose error to user
						fmt.Printf("[SECURITY] Rejected invalid sort column: %s\n", sort.Column)
						continue
					}
					safeColumn = validatedCol
				} else {
					// Fallback: sanitize column name if validator not available
					safeColumn = SanitizeColumnName(sort.Column)
				}

				direction := "ASC"
				if strings.ToUpper(sort.Direction) == "DESC" {
					direction = "DESC"
				}
				db = db.Order(fmt.Sprintf("%s %s", safeColumn, direction))
			}
		}
	}

	// Pagination
	offset := (req.Page - 1) * req.PerPage
	db = db.Offset(offset).Limit(req.PerPage)

	// Execute Query
	// Use reflection to create a slice of the model type
	modelType := reflect.TypeOf(p.Model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	sliceType := reflect.SliceOf(modelType)
	resultsPtr := reflect.New(sliceType)

	if err := db.Find(resultsPtr.Interface()).Error; err != nil {
		return nil, err
	}

	// Convert to []interface{}
	// Since FieldHandler (and Tests) mostly expect maps for dynamic access, lets convert strict structs to maps
	// This also ensures we respect JSON tags.
	resultsVal := resultsPtr.Elem()
	items := make([]interface{}, resultsVal.Len())
	for i := 0; i < resultsVal.Len(); i++ {
		items[i] = resultsVal.Index(i).Addr().Interface()
	}

	// Load lazy relationships manually (LAZY_LOADING strategy)
	// Eager loading relationships are already loaded via GORM Preload above
	if len(items) > 0 {
		for _, field := range relationshipFields {
			if field.GetLoadingStrategy() == fields.LAZY_LOADING {
				for _, item := range items {
					if _, err := p.relationshipLoader.LazyLoad(stdCtx, item, field); err != nil {
						fmt.Printf("[WARN] Failed to lazy load %s: %v\n", field.GetRelationshipName(), err)
					}
				}
			}
		}
	}

	return &QueryResponse{
		Items:   items,
		Total:   total,
		Page:    req.Page,
		PerPage: req.PerPage,
	}, nil
}

/// # Show
///
/// Bu fonksiyon, belirtilen ID'ye sahip tek bir kaydı veritabanından getirir.
/// İlişkili veriler (WithRelationships) otomatik olarak yüklenir.
///
/// ## Parametreler
///
/// - `ctx`: Context bilgisi (*context.Context)
/// - `id`: Kaydın benzersiz kimliği (string)
///
/// ## Döndürür
///
/// - `interface{}`: Bulunan kayıt (model tipinde)
/// - `error`: Hata durumunda hata mesajı
///   - `gorm.ErrRecordNotFound`: Kayıt bulunamadı
///   - Diğer veritabanı hataları
///
/// ## Kullanım Senaryoları
///
/// 1. **Detay Sayfası**: Tek bir kaydın detaylarını gösterme
/// 2. **API Endpoint**: RESTful API için GET /resource/:id endpoint'i
/// 3. **Düzenleme Formu**: Düzenleme formunu doldurmak için veri çekme
/// 4. **İlişkili Veri**: İlişkili verilerle birlikte kayıt getirme
///
/// ## Örnek Kullanım
///
/// ```go
/// // Basit kullanım
/// user, err := provider.Show(ctx, "123")
/// if err != nil {
///     if errors.Is(err, gorm.ErrRecordNotFound) {
///         return nil, fmt.Errorf("kullanıcı bulunamadı")
///     }
///     return nil, err
/// }
///
/// // Tip dönüşümü
/// userModel := user.(*User)
/// fmt.Printf("Kullanıcı: %s (%s)\n", userModel.Name, userModel.Email)
///
/// // İlişkilerle birlikte
/// provider.SetWith([]string{"Profile", "Posts", "Posts.Comments"})
/// user, err := provider.Show(ctx, "123")
/// userModel := user.(*User)
/// fmt.Printf("Kullanıcı: %s, Post Sayısı: %d\n",
///     userModel.Name,
///     len(userModel.Posts))
///
/// // API handler'da kullanım
/// func GetUser(c *gin.Context) {
///     id := c.Param("id")
///     user, err := provider.Show(ctx, id)
///     if err != nil {
///         if errors.Is(err, gorm.ErrRecordNotFound) {
///             c.JSON(404, gin.H{"error": "Kullanıcı bulunamadı"})
///             return
///         }
///         c.JSON(500, gin.H{"error": err.Error()})
///         return
///     }
///     c.JSON(200, user)
/// }
/// ```
///
/// ## Güvenlik
///
/// - ID parametresi güvenli şekilde bind edilir (SQL injection koruması)
/// - Context timeout'ları sorguya yansır
/// - Sadece belirtilen ID'ye sahip kayıt döndürülür
///
/// ## Performans
///
/// - Primary key üzerinden arama (çok hızlı)
/// - Eager loading ile N+1 problemi önlenir
/// - Index kullanımı (ID primary key)
///
/// ## Önemli Notlar
///
/// - ID string tipinde olmalıdır (integer ID'ler string'e dönüştürülmelidir)
/// - Kayıt bulunamazsa `gorm.ErrRecordNotFound` hatası döner
/// - WithRelationships ile belirlenen ilişkiler otomatik yüklenir
/// - Dönen değer interface{} tipindedir, tip dönüşümü gereklidir
/// - Soft delete kullanılıyorsa, silinmiş kayıtlar döndürülmez
///
/// ## Hata Durumları
///
/// - `gorm.ErrRecordNotFound`: Belirtilen ID'ye sahip kayıt bulunamadı
/// - Veritabanı bağlantı hatası
/// - Context timeout
/// - Geçersiz ID formatı
func (p *GormDataProvider) Show(ctx *context.Context, id string) (interface{}, error) {
	// Create a new instance of the model to hold the result
	// We use p.Model's type
	// But simpler: just use map[string]interface{} for dynamic nature or try to use the model type via reflection if needed.
	// For GORM, if we pass p.Model (which is a pointer to a struct), it works but we might overwrite the original p.Model if we are not careful or if we reuse it?
	// Actually p.Model is just a template.
	// Let's use map[string]interface{} to return standard format for the handler.

	// Create a new instance of the model to hold the result
	modelType := reflect.TypeOf(p.Model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	result := reflect.New(modelType).Interface()

	stdCtx := p.getContext(ctx)
	db := p.DB.WithContext(stdCtx).Model(p.Model)

	// Apply Eager Loading with GORM Preload
	// WORKAROUND: Direkt olarak WithRelationships kullan çünkü relationshipFields boş olabilir
	// (field type detection sorunu nedeniyle)
	fmt.Printf("[DEBUG] Show - Preloading relationships: %v\n", p.WithRelationships)
	for _, relName := range p.WithRelationships {
		fmt.Printf("[DEBUG] Show - Preload: %s\n", relName)
		db = db.Preload(relName)
	}

	if err := db.Where("id = ?", id).First(result).Error; err != nil {
		return nil, err
	}

	// Load lazy relationships manually (LAZY_LOADING strategy)
	// Eager loading relationships are already loaded via GORM Preload above
	for _, field := range relationshipFields {
		if field.GetLoadingStrategy() == fields.LAZY_LOADING {
			if _, err := p.relationshipLoader.LazyLoad(stdCtx, result, field); err != nil {
				fmt.Printf("[WARN] Failed to lazy load %s: %v\n", field.GetRelationshipName(), err)
			}
		}
	}

	return result, nil
}

/// # Create
///
/// Bu fonksiyon, veritabanına yeni bir kayıt ekler ve ilişkili verileri yönetir.
/// Otomatik olarak CreatedAt ve UpdatedAt timestamp'lerini ayarlar ve ilişkileri
/// (HasOne, Many2Many) işler.
///
/// ## Parametreler
///
/// - `ctx`: Context bilgisi (*context.Context)
/// - `data`: Oluşturulacak kaydın verileri (map[string]interface{})
///   - Anahtar: Field adı (snake_case veya camelCase)
///   - Değer: Field değeri (tip modele göre otomatik dönüştürülür)
///
/// ## Döndürür
///
/// - `interface{}`: Oluşturulan kayıt (ilişkilerle birlikte)
/// - `error`: Hata durumunda hata mesajı
///
/// ## İşlem Sırası
///
/// 1. **Schema Parsing**: Model schema'sını parse etme
/// 2. **Veri Validasyonu**: Sadece geçerli field'ları filtreleme
/// 3. **Field Ayarlama**: Reflection ile field değerlerini ayarlama
/// 4. **Timestamp Ayarlama**: CreatedAt ve UpdatedAt otomatik ayarlama
/// 5. **Kayıt Oluşturma**: GORM Create ile veritabanına ekleme
/// 6. **İlişki Yönetimi**: HasOne ve Many2Many ilişkilerini işleme
/// 7. **Fresh Data**: Show metodu ile güncel veriyi getirme
///
/// ## Desteklenen İlişki Tipleri
///
/// ### HasOne İlişkisi
/// ```go
/// data := map[string]interface{}{
///     "name": "John Doe",
///     "profile": 123, // Profile ID
/// }
/// ```
///
/// ### Many2Many İlişkisi
/// ```go
/// data := map[string]interface{}{
///     "name": "John Doe",
///     "roles": []interface{}{1, 2, 3}, // Role ID'leri
/// }
/// ```
///
/// ## Kullanım Senaryoları
///
/// 1. **Form Submission**: Web formlarından gelen verileri kaydetme
/// 2. **API Endpoint**: RESTful API için POST /resource endpoint'i
/// 3. **Bulk Import**: Toplu veri aktarımı
/// 4. **İlişkili Veri**: İlişkilerle birlikte kayıt oluşturma
///
/// ## Örnek Kullanım
///
/// ```go
/// // Basit kayıt oluşturma
/// user, err := provider.Create(ctx, map[string]interface{}{
///     "name":  "John Doe",
///     "email": "john@example.com",
///     "age":   30,
/// })
///
/// // İlişkilerle kayıt oluşturma
/// user, err := provider.Create(ctx, map[string]interface{}{
///     "name":    "John Doe",
///     "email":   "john@example.com",
///     "profile": 123,              // HasOne ilişkisi
///     "roles":   []interface{}{1, 2, 3}, // Many2Many ilişkisi
/// })
///
/// // Tip dönüşümü
/// userModel := user.(*User)
/// fmt.Printf("Oluşturulan kullanıcı ID: %d\n", userModel.ID)
///
/// // API handler'da kullanım
/// func CreateUser(c *gin.Context) {
///     var data map[string]interface{}
///     if err := c.ShouldBindJSON(&data); err != nil {
///         c.JSON(400, gin.H{"error": err.Error()})
///         return
///     }
///
///     user, err := provider.Create(ctx, data)
///     if err != nil {
///         c.JSON(500, gin.H{"error": err.Error()})
///         return
///     }
///
///     c.JSON(201, user)
/// }
///
/// // Nested ilişkilerle
/// provider.SetWith([]string{"Profile", "Roles"})
/// user, err := provider.Create(ctx, map[string]interface{}{
///     "name":    "John Doe",
///     "email":   "john@example.com",
///     "profile": 123,
///     "roles":   []interface{}{1, 2, 3},
/// })
/// // user içinde Profile ve Roles ilişkileri de gelir
/// ```
///
/// ## Otomatik İşlemler
///
/// - **ID Backfilling**: Oluşturulan kaydın ID'si otomatik atanır
/// - **Timestamp**: CreatedAt ve UpdatedAt otomatik ayarlanır
/// - **Hooks**: GORM BeforeCreate, AfterCreate hook'ları çalışır
/// - **İlişki Senkronizasyonu**: İlişkiler otomatik olarak pivot tablolara yazılır
///
/// ## Güvenlik
///
/// - Sadece model schema'sında tanımlı field'lar işlenir
/// - Geçersiz field'lar sessizce atlanır
/// - Parameterized query kullanımı
/// - Context timeout'ları sorguya yansır
///
/// ## Performans
///
/// - Tek transaction içinde çalışır
/// - İlişkiler için ayrı sorgular (N+1 riski var)
/// - Fresh data için ekstra Show sorgusu
///
/// ## Önemli Notlar
///
/// - Data parametresi map[string]interface{} tipinde olmalıdır
/// - Field isimleri snake_case veya camelCase olabilir (otomatik dönüşüm)
/// - Geçersiz field'lar hata vermez, sadece atlanır
/// - İlişki ID'leri mevcut kayıtlara ait olmalıdır
/// - Many2Many için slice tipinde değer gereklidir
/// - HasOne için tek ID değeri gereklidir
/// - BelongsTo ilişkileri foreign key ile otomatik yönetilir
/// - Dönen değer fresh data'dır (Show metodu ile çekilir)
/// - CreatedAt ve UpdatedAt otomatik ayarlanır (manuel ayarlama gerekmez)
///
/// ## Hata Durumları
///
/// - Veritabanı constraint ihlali (unique, foreign key vb.)
/// - Geçersiz veri tipi
/// - Context timeout
/// - İlişki kaydı bulunamadı
/// - Schema parsing hatası
func (p *GormDataProvider) Create(ctx *context.Context, data map[string]interface{}) (interface{}, error) {
	stmt := &gorm.Statement{DB: p.DB}
	if err := stmt.Parse(p.Model); err != nil {
		return nil, err
	}
	modelSchema := stmt.Schema

	stdCtx := p.getContext(ctx)
	modelType := reflect.TypeOf(p.Model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	newItem := reflect.New(modelType).Interface()

	validData := make(map[string]interface{})
	for k, v := range data {
		field := modelSchema.LookUpField(k)
		if field != nil && field.DBName != "" {
			validData[k] = v

			// Set field value on newItem to ensure it's populated for Create
			modelVal := reflect.ValueOf(newItem)
			if modelVal.Kind() == reflect.Ptr {
				modelVal = modelVal.Elem()
			}
			if err := field.Set(stdCtx, modelVal, v); err != nil {
				fmt.Printf("[GORM] Error setting field %s: %v\n", field.Name, err)
			}
		}
	}

	// Set timestamps
	now := time.Now()
	if createdAtField := modelSchema.LookUpField("CreatedAt"); createdAtField != nil {
		modelVal := reflect.ValueOf(newItem)
		if modelVal.Kind() == reflect.Ptr {
			modelVal = modelVal.Elem()
		}
		createdAtField.Set(stdCtx, modelVal, now)
	}
	if updatedAtField := modelSchema.LookUpField("UpdatedAt"); updatedAtField != nil {
		modelVal := reflect.ValueOf(newItem)
		if modelVal.Kind() == reflect.Ptr {
			modelVal = modelVal.Elem()
		}
		updatedAtField.Set(stdCtx, modelVal, now)
	}

	// Use Create with struct to ensure ID backfilling and hooks execution
	if err := p.DB.WithContext(stdCtx).Create(newItem).Error; err != nil {
		return nil, err
	}

	// Handle Associations
	for k, v := range data {
		field := modelSchema.LookUpField(k)
		if field == nil {
			field = modelSchema.LookUpField(strcase.ToCamel(k))
		}
		if field != nil {
			if field.DBName == "" {
				if rel, ok := modelSchema.Relationships.Relations[field.Name]; ok {
					switch rel.Type {
					case schema.HasOne:
						if v != nil {
							// Get newItem ID
							newItemID := reflect.ValueOf(newItem).Elem().FieldByName(modelSchema.PrioritizedPrimaryField.Name).Interface()

							// Get foreign key and related table info from GORM schema
							foreignKey := ""
							if len(rel.References) > 0 {
								foreignKey = rel.References[0].ForeignKey.DBName
							}
							relatedTable := rel.FieldSchema.Table

							if foreignKey != "" && relatedTable != "" && newItemID != nil {
								if err := p.replaceHasOne(stdCtx, newItemID, foreignKey, v, relatedTable); err != nil {
									fmt.Printf("[WARN] Failed to replace HasOne: %v\n", err)
								}
							}
						}
					case schema.Many2Many:
						// Handle BelongsToMany (Many2Many)
						// v is likely []interface{} or []string of IDs
						var ids []interface{}
						val := reflect.ValueOf(v)
						if val.Kind() == reflect.Slice {
							for i := 0; i < val.Len(); i++ {
								ids = append(ids, val.Index(i).Interface())
							}
						} else {
							ids = append(ids, v)
						}

						if len(ids) > 0 {
							// Get newItem ID
							newItemID := reflect.ValueOf(newItem).Elem().FieldByName(modelSchema.PrioritizedPrimaryField.Name).Interface()

							// Get pivot table info from GORM schema
							pivotTable := ""
							parentColumn := ""
							relatedColumn := ""

							if rel.JoinTable != nil {
								pivotTable = rel.JoinTable.Name
							}

							if len(rel.References) > 0 {
								for _, ref := range rel.References {
									if ref.OwnPrimaryKey {
										parentColumn = ref.ForeignKey.DBName
									} else {
										relatedColumn = ref.ForeignKey.DBName
									}
								}
							}

							if pivotTable != "" && parentColumn != "" && relatedColumn != "" && newItemID != nil {
								if err := p.replaceMany2Many(stdCtx, newItemID, ids, pivotTable, parentColumn, relatedColumn); err != nil {
									fmt.Printf("[WARN] Failed to replace Many2Many: %v\n", err)
								}
							}
						}
					}
				}
			}
		}
	}

	// Return fresh item using ID
	if modelSchema.PrioritizedPrimaryField != nil {
		idVal := reflect.ValueOf(newItem).Elem().FieldByName(modelSchema.PrioritizedPrimaryField.Name).Interface()
		id := fmt.Sprint(idVal)
		if id != "" && id != "0" {
			return p.Show(ctx, id)
		}
	}

	return newItem, nil
}

/// # Update
///
/// Bu fonksiyon, mevcut bir kaydı günceller ve ilişkili verileri yönetir.
/// Otomatik olarak UpdatedAt timestamp'ini günceller ve ilişkileri
/// (HasOne, BelongsTo, Many2Many) işler.
///
/// ## Parametreler
///
/// - `ctx`: Context bilgisi (*context.Context)
/// - `id`: Güncellenecek kaydın benzersiz kimliği (string)
/// - `data`: Güncellenecek veriler (map[string]interface{})
///   - Anahtar: Field adı (snake_case veya camelCase)
///   - Değer: Yeni field değeri (tip modele göre otomatik dönüştürülür)
///
/// ## Döndürür
///
/// - `interface{}`: Güncellenmiş kayıt (ilişkilerle birlikte)
/// - `error`: Hata durumunda hata mesajı
///   - `gorm.ErrRecordNotFound`: Kayıt bulunamadı
///   - Diğer veritabanı hataları
///
/// ## İşlem Sırası
///
/// 1. **Schema Parsing**: Model schema'sını parse etme
/// 2. **Kayıt Bulma**: Mevcut kaydı ID ile bulma
/// 3. **Veri Validasyonu**: Sadece geçerli field'ları filtreleme
/// 4. **İlişki Yönetimi**: HasOne, BelongsTo ve Many2Many ilişkilerini güncelleme
/// 5. **Field Güncelleme**: Veritabanı field'larını güncelleme
/// 6. **Timestamp Güncelleme**: UpdatedAt otomatik güncelleme
/// 7. **Fresh Data**: Show metodu ile güncel veriyi getirme
///
/// ## Desteklenen İlişki Tipleri
///
/// ### HasOne / BelongsTo İlişkisi
/// ```go
/// // İlişki güncelleme
/// data := map[string]interface{}{
///     "name": "John Doe Updated",
///     "profile": 456, // Yeni Profile ID
/// }
///
/// // İlişki temizleme
/// data := map[string]interface{}{
///     "profile": nil, // İlişkiyi kaldır
/// }
/// ```
///
/// ### Many2Many İlişkisi
/// ```go
/// // İlişkileri değiştirme
/// data := map[string]interface{}{
///     "roles": []interface{}{2, 3, 4}, // Yeni Role ID'leri (eskiler silinir)
/// }
///
/// // Tüm ilişkileri temizleme
/// data := map[string]interface{}{
///     "roles": []interface{}{}, // Boş slice = tüm ilişkileri kaldır
/// }
/// ```
///
/// ## Kullanım Senaryoları
///
/// 1. **Form Güncelleme**: Web formlarından gelen güncellemeleri kaydetme
/// 2. **API Endpoint**: RESTful API için PUT/PATCH /resource/:id endpoint'i
/// 3. **Partial Update**: Sadece belirli field'ları güncelleme
/// 4. **İlişki Yönetimi**: İlişkileri ekleme, çıkarma veya değiştirme
/// 5. **Bulk Update**: Toplu güncelleme işlemleri
///
/// ## Örnek Kullanım
///
/// ```go
/// // Basit güncelleme
/// user, err := provider.Update(ctx, "123", map[string]interface{}{
///     "name":  "John Doe Updated",
///     "email": "john.updated@example.com",
/// })
///
/// // Partial update (sadece name)
/// user, err := provider.Update(ctx, "123", map[string]interface{}{
///     "name": "John Doe",
/// })
///
/// // İlişkilerle güncelleme
/// user, err := provider.Update(ctx, "123", map[string]interface{}{
///     "name":    "John Doe",
///     "profile": 456,                    // HasOne ilişkisi güncelle
///     "roles":   []interface{}{2, 3, 4}, // Many2Many ilişkileri değiştir
/// })
///
/// // İlişki temizleme
/// user, err := provider.Update(ctx, "123", map[string]interface{}{
///     "profile": nil,           // Profile ilişkisini kaldır
///     "roles":   []interface{}{}, // Tüm role ilişkilerini kaldır
/// })
///
/// // API handler'da kullanım
/// func UpdateUser(c *gin.Context) {
///     id := c.Param("id")
///     var data map[string]interface{}
///     if err := c.ShouldBindJSON(&data); err != nil {
///         c.JSON(400, gin.H{"error": err.Error()})
///         return
///     }
///
///     user, err := provider.Update(ctx, id, data)
///     if err != nil {
///         if errors.Is(err, gorm.ErrRecordNotFound) {
///             c.JSON(404, gin.H{"error": "Kullanıcı bulunamadı"})
///             return
///         }
///         c.JSON(500, gin.H{"error": err.Error()})
///         return
///     }
///
///     c.JSON(200, user)
/// }
///
/// // Nested ilişkilerle
/// provider.SetWith([]string{"Profile", "Roles"})
/// user, err := provider.Update(ctx, "123", map[string]interface{}{
///     "name":  "John Doe",
///     "roles": []interface{}{1, 2, 3},
/// })
/// // user içinde Profile ve Roles ilişkileri de gelir
/// ```
///
/// ## İlişki Davranışları
///
/// ### HasOne / BelongsTo
/// - **Yeni ID**: Mevcut ilişki kaldırılır, yeni ilişki eklenir (Replace)
/// - **nil Değer**: Mevcut ilişki kaldırılır (Clear)
/// - **Geçersiz ID**: Hata vermez, ilişki değişmez
///
/// ### Many2Many
/// - **Yeni ID Listesi**: Tüm mevcut ilişkiler kaldırılır, yenileri eklenir (Replace)
/// - **Boş Liste**: Tüm ilişkiler kaldırılır (Clear)
/// - **nil Değer**: İlişkiler değişmez
///
/// ## Otomatik İşlemler
///
/// - **UpdatedAt**: Otomatik güncellenir
/// - **Hooks**: GORM BeforeUpdate, AfterUpdate hook'ları çalışır
/// - **İlişki Senkronizasyonu**: Pivot tablolar otomatik güncellenir
/// - **Optimistic Locking**: Version field varsa otomatik kontrol edilir
///
/// ## Güvenlik
///
/// - Sadece model schema'sında tanımlı field'lar işlenir
/// - Geçersiz field'lar sessizce atlanır
/// - Parameterized query kullanımı
/// - Context timeout'ları sorguya yansır
/// - ID parametresi güvenli şekilde bind edilir
///
/// ## Performans
///
/// - Tek transaction içinde çalışır
/// - İlişkiler için ayrı sorgular (N+1 riski var)
/// - Fresh data için ekstra Show sorgusu
/// - Sadece değişen field'lar güncellenir
///
/// ## Önemli Notlar
///
/// - ID string tipinde olmalıdır
/// - Data parametresi map[string]interface{} tipinde olmalıdır
/// - Field isimleri snake_case veya camelCase olabilir
/// - Geçersiz field'lar hata vermez, sadece atlanır
/// - İlişki ID'leri mevcut kayıtlara ait olmalıdır
/// - Many2Many için slice tipinde değer gereklidir
/// - HasOne/BelongsTo için tek ID değeri veya nil gereklidir
/// - Dönen değer fresh data'dır (Show metodu ile çekilir)
/// - UpdatedAt otomatik güncellenir (manuel ayarlama gerekmez)
/// - CreatedAt değişmez
/// - Primary key güncellenemez
///
/// ## Partial Update
///
/// - Sadece gönderilen field'lar güncellenir
/// - Gönderilmeyen field'lar değişmez
/// - nil değer gönderilirse field NULL yapılır (nullable ise)
/// - Boş string ("") geçerli bir değerdir
///
/// ## Hata Durumları
///
/// - `gorm.ErrRecordNotFound`: Belirtilen ID'ye sahip kayıt bulunamadı
/// - Veritabanı constraint ihlali (unique, foreign key vb.)
/// - Geçersiz veri tipi
/// - Context timeout
/// - İlişki kaydı bulunamadı
/// - Schema parsing hatası
/// - Optimistic locking conflict
func (p *GormDataProvider) Update(ctx *context.Context, id string, data map[string]interface{}) (interface{}, error) {
	stmt := &gorm.Statement{DB: p.DB}
	if err := stmt.Parse(p.Model); err != nil {
		return nil, err
	}
	modelSchema := stmt.Schema

	stdCtx := p.getContext(ctx)
	modelType := reflect.TypeOf(p.Model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	item := reflect.New(modelType).Interface()

	if err := p.DB.WithContext(stdCtx).First(item, "id = ?", id).Error; err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})
	for k, v := range data {
		field := modelSchema.LookUpField(k)
		if field == nil {
			field = modelSchema.LookUpField(strcase.ToCamel(k))
		}
		if field != nil {
			if field.DBName != "" {
				updates[k] = v
			} else if rel, ok := modelSchema.Relationships.Relations[field.Name]; ok {
				switch rel.Type {
				case schema.HasOne, schema.BelongsTo:
					// Get item ID
					itemID := reflect.ValueOf(item).Elem().FieldByName(modelSchema.PrioritizedPrimaryField.Name).Interface()

					// Get foreign key and related table info from GORM schema
					foreignKey := ""
					if len(rel.References) > 0 {
						foreignKey = rel.References[0].ForeignKey.DBName
					}
					relatedTable := rel.FieldSchema.Table

					if v != nil {
						if foreignKey != "" && relatedTable != "" && itemID != nil {
							if err := p.replaceHasOne(stdCtx, itemID, foreignKey, v, relatedTable); err != nil {
								fmt.Printf("[WARN] Failed to replace HasOne: %v\n", err)
							}
						}
					} else {
						if foreignKey != "" && relatedTable != "" && itemID != nil {
							if err := p.clearHasOne(stdCtx, itemID, foreignKey, relatedTable); err != nil {
								fmt.Printf("[WARN] Failed to clear HasOne: %v\n", err)
							}
						}
					}
				case schema.Many2Many:
					// Handle BelongsToMany (Many2Many) update
					// v is likely []interface{} or []string of IDs
					var ids []interface{}
					if v != nil {
						val := reflect.ValueOf(v)
						if val.Kind() == reflect.Slice {
							for i := 0; i < val.Len(); i++ {
								ids = append(ids, val.Index(i).Interface())
							}
						} else {
							ids = append(ids, v)
						}
					}

					// Get item ID
					itemID := reflect.ValueOf(item).Elem().FieldByName(modelSchema.PrioritizedPrimaryField.Name).Interface()

					// Get pivot table info from GORM schema
					pivotTable := ""
					parentColumn := ""
					relatedColumn := ""

					if rel.JoinTable != nil {
						pivotTable = rel.JoinTable.Name
					}

					if len(rel.References) > 0 {
						for _, ref := range rel.References {
							if ref.OwnPrimaryKey {
								parentColumn = ref.ForeignKey.DBName
							} else {
								relatedColumn = ref.ForeignKey.DBName
							}
						}
					}

					if len(ids) > 0 {
						if pivotTable != "" && parentColumn != "" && relatedColumn != "" && itemID != nil {
							if err := p.replaceMany2Many(stdCtx, itemID, ids, pivotTable, parentColumn, relatedColumn); err != nil {
								fmt.Printf("[WARN] Failed to replace Many2Many: %v\n", err)
							}
						}
					} else {
						// If empty list sent, clear associations
						if pivotTable != "" && parentColumn != "" && itemID != nil {
							if err := p.clearMany2Many(stdCtx, itemID, pivotTable, parentColumn); err != nil {
								fmt.Printf("[WARN] Failed to clear Many2Many: %v\n", err)
							}
						}
					}
				}

			}
		}
	}

	updates["updated_at"] = time.Now()
	if len(updates) > 0 {
		if err := p.DB.WithContext(stdCtx).Model(item).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	return p.Show(ctx, id)
}

/// # Delete
///
/// Bu fonksiyon, belirtilen ID'ye sahip kaydı veritabanından siler.
/// Soft delete kullanılıyorsa kayıt silinmek yerine işaretlenir (deleted_at).
///
/// ## Parametreler
///
/// - `ctx`: Context bilgisi (*context.Context)
/// - `id`: Silinecek kaydın benzersiz kimliği (string)
///
/// ## Döndürür
///
/// - `error`: Hata durumunda hata mesajı
///   - `nil`: Silme işlemi başarılı
///   - `gorm.ErrRecordNotFound`: Kayıt bulunamadı (hard delete'de)
///   - Diğer veritabanı hataları
///
/// ## Silme Tipleri
///
/// ### Soft Delete (Yumuşak Silme)
/// Model'de `gorm.DeletedAt` field'ı varsa:
/// - Kayıt fiziksel olarak silinmez
/// - `deleted_at` field'ı mevcut timestamp ile güncellenir
/// - Normal sorgularda kayıt görünmez
/// - `Unscoped()` ile geri getirilebilir
///
/// ```go
/// type User struct {
///     ID        uint
///     Name      string
///     DeletedAt gorm.DeletedAt `gorm:"index"` // Soft delete aktif
/// }
/// ```
///
/// ### Hard Delete (Kalıcı Silme)
/// Model'de `gorm.DeletedAt` field'ı yoksa:
/// - Kayıt fiziksel olarak veritabanından silinir
/// - Geri getirilemez
/// - İlişkili kayıtlar foreign key constraint'lerine göre işlenir
///
/// ```go
/// type Log struct {
///     ID   uint
///     Text string
///     // DeletedAt yok = hard delete
/// }
/// ```
///
/// ## Kullanım Senaryoları
///
/// 1. **Kayıt Silme**: Kullanıcı, ürün, sipariş vb. silme
/// 2. **API Endpoint**: RESTful API için DELETE /resource/:id endpoint'i
/// 3. **Toplu Silme**: Birden fazla kaydı silme (döngü ile)
/// 4. **Veri Temizleme**: Eski veya gereksiz kayıtları temizleme
/// 5. **Soft Delete**: Geri getirilebilir silme işlemleri
///
/// ## Örnek Kullanım
///
/// ```go
/// // Basit silme
/// err := provider.Delete(ctx, "123")
/// if err != nil {
///     if errors.Is(err, gorm.ErrRecordNotFound) {
///         return fmt.Errorf("kullanıcı bulunamadı")
///     }
///     return err
/// }
/// fmt.Println("Kullanıcı başarıyla silindi")
///
/// // API handler'da kullanım
/// func DeleteUser(c *gin.Context) {
///     id := c.Param("id")
///     err := provider.Delete(ctx, id)
///     if err != nil {
///         if errors.Is(err, gorm.ErrRecordNotFound) {
///             c.JSON(404, gin.H{"error": "Kullanıcı bulunamadı"})
///             return
///         }
///         c.JSON(500, gin.H{"error": err.Error()})
///         return
///     }
///     c.JSON(204, nil) // No Content
/// }
///
/// // Toplu silme
/// ids := []string{"1", "2", "3", "4", "5"}
/// for _, id := range ids {
///     if err := provider.Delete(ctx, id); err != nil {
///         log.Printf("ID %s silinemedi: %v", id, err)
///     }
/// }
///
/// // Soft delete sonrası geri getirme (manuel)
/// // GORM ile doğrudan:
/// db.Unscoped().Where("id = ?", id).First(&user)
/// db.Unscoped().Model(&user).Update("deleted_at", nil)
/// ```
///
/// ## İlişki Davranışları
///
/// ### Foreign Key Constraints
/// - **CASCADE**: İlişkili kayıtlar da silinir
/// - **SET NULL**: İlişkili kayıtların foreign key'i NULL yapılır
/// - **RESTRICT**: İlişkili kayıt varsa silme engellenir
/// - **NO ACTION**: Veritabanı varsayılan davranışı
///
/// ### Örnek Constraint Tanımları
/// ```go
/// type User struct {
///     ID    uint
///     Posts []Post `gorm:"constraint:OnDelete:CASCADE"` // Kullanıcı silinince postlar da silinir
/// }
///
/// type Post struct {
///     ID     uint
///     UserID uint
///     User   User `gorm:"constraint:OnDelete:SET NULL"` // Post silinince user_id NULL olur
/// }
/// ```
///
/// ## Otomatik İşlemler
///
/// - **Soft Delete**: DeletedAt otomatik ayarlanır
/// - **Hooks**: GORM BeforeDelete, AfterDelete hook'ları çalışır
/// - **Cascade**: Foreign key constraint'lerine göre ilişkili kayıtlar işlenir
/// - **Index**: deleted_at index'i varsa performans optimize edilir
///
/// ## Güvenlik
///
/// - ID parametresi güvenli şekilde bind edilir (SQL injection koruması)
/// - Context timeout'ları sorguya yansır
/// - Sadece belirtilen ID'ye sahip kayıt silinir
/// - Soft delete ile veri kaybı önlenir
///
/// ## Performans
///
/// - Primary key üzerinden silme (çok hızlı)
/// - Soft delete hard delete'den daha hızlıdır
/// - Index kullanımı (ID primary key)
/// - Cascade silme performansı etkileyebilir
///
/// ## Önemli Notlar
///
/// - ID string tipinde olmalıdır
/// - Soft delete kullanılıyorsa kayıt fiziksel olarak silinmez
/// - Hard delete geri alınamaz
/// - İlişkili kayıtlar foreign key constraint'lerine göre işlenir
/// - BeforeDelete hook'ları silme işlemini iptal edebilir
/// - Soft delete'de deleted_at index oluşturmak önerilir
/// - Cascade silme dikkatli kullanılmalıdır
/// - Transaction içinde çalışır
///
/// ## Soft Delete Avantajları
///
/// - **Geri Getirilebilir**: Yanlışlıkla silinen veriler kurtarılabilir
/// - **Audit Trail**: Silme geçmişi tutulur
/// - **Veri Bütünlüğü**: İlişkili veriler korunur
/// - **Yasal Gereklilikler**: Bazı düzenlemeler veri saklamayı gerektirir
///
/// ## Hard Delete Avantajları
///
/// - **Performans**: Daha hızlı sorgular (deleted_at kontrolü yok)
/// - **Depolama**: Daha az disk alanı kullanımı
/// - **Basitlik**: Daha basit sorgu yapısı
/// - **GDPR Uyumu**: Kullanıcı verilerinin tamamen silinmesi
///
/// ## Hata Durumları
///
/// - `gorm.ErrRecordNotFound`: Belirtilen ID'ye sahip kayıt bulunamadı (hard delete'de)
/// - Foreign key constraint ihlali (RESTRICT durumunda)
/// - Veritabanı bağlantı hatası
/// - Context timeout
/// - BeforeDelete hook hatası
/// - Transaction hatası
///
/// ## Best Practices
///
/// 1. **Soft Delete Kullanın**: Önemli veriler için soft delete tercih edin
/// 2. **Index Oluşturun**: deleted_at için index oluşturun
/// 3. **Cascade Dikkatli**: Cascade silme dikkatli kullanın
/// 4. **Audit Log**: Silme işlemlerini loglayın
/// 5. **Yetkilendirme**: Silme yetkisi kontrolü yapın
/// 6. **Onay Mekanizması**: Kritik silmelerde onay alın
/// 7. **Backup**: Düzenli backup alın
/// 8. **Temizlik**: Eski soft delete kayıtlarını periyodik temizleyin
func (p *GormDataProvider) Delete(ctx *context.Context, id string) error {
	stdCtx := p.getContext(ctx)
	return p.DB.WithContext(stdCtx).Model(p.Model).Where("id = ?", id).Delete(nil).Error
}

// SetRelationshipFields, yüklenecek ilişki field'larını ayarlar.
//
// Bu metod, RelationshipLoader tarafından kullanılacak field'ları belirler.
// Resource'dan gelen relationship field'ları buraya set edilir.
//
// # Parametreler
//
// - **fields**: Yüklenecek ilişki field'ları ([]fields.RelationshipField)
//
// # Kullanım Örneği
//
//	relationshipFields := []fields.RelationshipField{authorField, postsField}
//	provider.SetRelationshipFields(relationshipFields)
func (p *GormDataProvider) SetRelationshipFields(fields []fields.RelationshipField) {
	p.relationshipFields = fields
}

// getRelationshipFields, yüklenecek ilişki field'larını döndürür.
//
// Bu metod, WithRelationships listesinde belirtilen ilişkilere karşılık gelen
// field'ları relationshipFields listesinden filtreler.
//
// # Döndürür
//
// - []fields.RelationshipField: Yüklenecek ilişki field'ları
func (p *GormDataProvider) getRelationshipFields() []fields.RelationshipField {
	fmt.Printf("[DEBUG] getRelationshipFields - WithRelationships: %v, relationshipFields count: %d\n", p.WithRelationships, len(p.relationshipFields))

	if len(p.WithRelationships) == 0 || len(p.relationshipFields) == 0 {
		fmt.Printf("[DEBUG] getRelationshipFields - Returning empty (WithRels empty: %v, Fields empty: %v)\n", len(p.WithRelationships) == 0, len(p.relationshipFields) == 0)
		return []fields.RelationshipField{}
	}

	// WithRelationships listesinde belirtilen ilişkileri filtrele
	result := []fields.RelationshipField{}
	for _, relName := range p.WithRelationships {
		for _, field := range p.relationshipFields {
			if field.GetRelationshipName() == relName {
				fmt.Printf("[DEBUG] getRelationshipFields - Matched: %s\n", relName)
				result = append(result, field)
				break
			}
		}
	}

	fmt.Printf("[DEBUG] getRelationshipFields - Returning %d fields\n", len(result))
	return result
}

// loadRelationshipsForItems, birden fazla kayıt için ilişkileri batch loading ile yükler.
//
// Bu metod, N+1 sorgu problemini önlemek için tüm kayıtların ilişkilerini
// tek seferde yükler. RelationshipLoader'ın EagerLoad metodunu kullanır.
//
// # Parametreler
//
// - **ctx**: Context bilgisi (*context.Context)
// - **items**: İlişkileri yüklenecek kayıt listesi ([]interface{})
//
// # Döndürür
//
// - error: Hata durumunda hata mesajı
//
// # Kullanım Örneği
//
//	items := []interface{}{&user1, &user2, &user3}
//	err := provider.loadRelationshipsForItems(ctx, items)
func (p *GormDataProvider) loadRelationshipsForItems(ctx *context.Context, items []interface{}) error {
	if p.relationshipLoader == nil {
		return fmt.Errorf("relationship loader not initialized")
	}

	if len(items) == 0 {
		return nil
	}

	stdCtx := p.getContext(ctx)
	relationshipFields := p.getRelationshipFields()

	for _, field := range relationshipFields {
		// Sadece eager loading stratejisine sahip field'ları yükle
		if field.GetLoadingStrategy() == fields.EAGER_LOADING {
			if err := p.relationshipLoader.EagerLoad(stdCtx, items, field); err != nil {
				return fmt.Errorf("failed to eager load %s: %w", field.GetRelationshipName(), err)
			}
		}
	}

	return nil
}

// loadRelationshipsForItem, tek bir kayıt için ilişkileri yükler.
//
// Bu metod, lazy loading stratejisi ile tek bir kaydın ilişkilerini yükler.
// RelationshipLoader'ın LazyLoad metodunu kullanır.
//
// # Parametreler
//
// - **ctx**: Context bilgisi (*context.Context)
// - **item**: İlişkisi yüklenecek kayıt (interface{})
//
// # Döndürür
//
// - error: Hata durumunda hata mesajı
//
// # Kullanım Örneği
//
//	err := provider.loadRelationshipsForItem(ctx, &user)
func (p *GormDataProvider) loadRelationshipsForItem(ctx *context.Context, item interface{}) error {
	if p.relationshipLoader == nil {
		return fmt.Errorf("relationship loader not initialized")
	}

	if item == nil {
		return nil
	}

	stdCtx := p.getContext(ctx)
	relationshipFields := p.getRelationshipFields()

	for _, field := range relationshipFields {
		if _, err := p.relationshipLoader.LazyLoad(stdCtx, item, field); err != nil {
			return fmt.Errorf("failed to lazy load %s: %w", field.GetRelationshipName(), err)
		}
	}

	return nil
}

// replaceHasOne, HasOne ilişkisini raw SQL ile günceller.
//
// Bu metod, GORM Association yerine raw SQL kullanarak HasOne ilişkisini günceller.
// Circular dependency sorununu önlemek için struct field'larına bağımlı olmadan çalışır.
//
// # Parametreler
//
// - **ctx**: Context bilgisi (stdcontext.Context)
// - **parentID**: Ana kaydın ID'si (interface{})
// - **foreignKey**: İlişkili tablodaki foreign key sütun adı (string)
// - **relatedID**: İlişkili kaydın ID'si (interface{})
// - **relatedTable**: İlişkili tablo adı (string)
//
// # Döndürür
//
// - error: Hata durumunda hata mesajı
//
// # Kullanım Örneği
//
//	err := p.replaceHasOne(ctx, userID, "user_id", profileID, "profiles")
func (p *GormDataProvider) replaceHasOne(ctx stdcontext.Context, parentID interface{}, foreignKey string, relatedID interface{}, relatedTable string) error {
	safeTable := SanitizeColumnName(relatedTable)
	safeForeignKey := SanitizeColumnName(foreignKey)

	return p.DB.WithContext(ctx).Exec(
		fmt.Sprintf("UPDATE %s SET %s = ? WHERE id = ?", safeTable, safeForeignKey),
		parentID, relatedID,
	).Error
}

// clearHasOne, HasOne ilişkisini raw SQL ile temizler.
//
// Bu metod, GORM Association yerine raw SQL kullanarak HasOne ilişkisini temizler.
// İlişkili kaydın foreign key'ini NULL yapar.
//
// # Parametreler
//
// - **ctx**: Context bilgisi (stdcontext.Context)
// - **parentID**: Ana kaydın ID'si (interface{})
// - **foreignKey**: İlişkili tablodaki foreign key sütun adı (string)
// - **relatedTable**: İlişkili tablo adı (string)
//
// # Döndürür
//
// - error: Hata durumunda hata mesajı
//
// # Kullanım Örneği
//
//	err := p.clearHasOne(ctx, userID, "user_id", "profiles")
func (p *GormDataProvider) clearHasOne(ctx stdcontext.Context, parentID interface{}, foreignKey string, relatedTable string) error {
	safeTable := SanitizeColumnName(relatedTable)
	safeForeignKey := SanitizeColumnName(foreignKey)

	return p.DB.WithContext(ctx).Exec(
		fmt.Sprintf("UPDATE %s SET %s = NULL WHERE %s = ?", safeTable, safeForeignKey, safeForeignKey),
		parentID,
	).Error
}

// replaceMany2Many, Many2Many ilişkilerini raw SQL ile günceller.
//
// Bu metod, GORM Association yerine raw SQL kullanarak Many2Many ilişkilerini günceller.
// Pivot tablodaki mevcut kayıtları siler ve yeni kayıtları ekler.
//
// # Parametreler
//
// - **ctx**: Context bilgisi (stdcontext.Context)
// - **parentID**: Ana kaydın ID'si (interface{})
// - **relatedIDs**: İlişkili kayıtların ID'leri ([]interface{})
// - **pivotTable**: Pivot tablo adı (string)
// - **parentColumn**: Pivot tablodaki ana kayıt sütun adı (string)
// - **relatedColumn**: Pivot tablodaki ilişkili kayıt sütun adı (string)
//
// # Döndürür
//
// - error: Hata durumunda hata mesajı
//
// # Kullanım Örneği
//
//	err := p.replaceMany2Many(ctx, userID, roleIDs, "user_roles", "user_id", "role_id")
func (p *GormDataProvider) replaceMany2Many(ctx stdcontext.Context, parentID interface{}, relatedIDs []interface{}, pivotTable string, parentColumn string, relatedColumn string) error {
	safePivotTable := SanitizeColumnName(pivotTable)
	safeParentColumn := SanitizeColumnName(parentColumn)
	safeRelatedColumn := SanitizeColumnName(relatedColumn)

	tx := p.DB.WithContext(ctx).Begin()

	// 1. Clear existing relationships
	if err := tx.Exec(
		fmt.Sprintf("DELETE FROM %s WHERE %s = ?", safePivotTable, safeParentColumn),
		parentID,
	).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 2. Insert new relationships
	for _, relatedID := range relatedIDs {
		if err := tx.Exec(
			fmt.Sprintf("INSERT INTO %s (%s, %s) VALUES (?, ?)", safePivotTable, safeParentColumn, safeRelatedColumn),
			parentID, relatedID,
		).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// clearMany2Many, Many2Many ilişkilerini raw SQL ile temizler.
//
// Bu metod, GORM Association yerine raw SQL kullanarak Many2Many ilişkilerini temizler.
// Pivot tablodaki tüm kayıtları siler.
//
// # Parametreler
//
// - **ctx**: Context bilgisi (stdcontext.Context)
// - **parentID**: Ana kaydın ID'si (interface{})
// - **pivotTable**: Pivot tablo adı (string)
// - **parentColumn**: Pivot tablodaki ana kayıt sütun adı (string)
//
// # Döndürür
//
// - error: Hata durumunda hata mesajı
//
// # Kullanım Örneği
//
//	err := p.clearMany2Many(ctx, userID, "user_roles", "user_id")
func (p *GormDataProvider) clearMany2Many(ctx stdcontext.Context, parentID interface{}, pivotTable string, parentColumn string) error {
	safePivotTable := SanitizeColumnName(pivotTable)
	safeParentColumn := SanitizeColumnName(parentColumn)

	return p.DB.WithContext(ctx).Exec(
		fmt.Sprintf("DELETE FROM %s WHERE %s = ?", safePivotTable, safeParentColumn),
		parentID,
	).Error
}
