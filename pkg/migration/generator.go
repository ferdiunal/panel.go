package migration

import (
	"fmt"
	"strings"

	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/resource"
	"github.com/iancoleman/strcase"
	"gorm.io/gorm"
)

/// # MigrationGenerator
///
/// Bu yapı, Resource tanımlarından veritabanı migration işlemlerini otomatik olarak yönetir.
/// GORM ORM kütüphanesi ile entegre çalışarak, resource field tanımlarını veritabanı
/// şemasına dönüştürür ve gerekli constraint'leri uygular.
///
/// ## Özellikler
///
/// - **Otomatik Migration**: Resource tanımlarından otomatik tablo oluşturma
/// - **Constraint Yönetimi**: Index, unique constraint ve foreign key yönetimi
/// - **İlişki Desteği**: BelongsTo, HasOne, HasMany, BelongsToMany ilişkileri
/// - **Pivot Tablo**: Many-to-Many ilişkiler için otomatik pivot tablo oluşturma
/// - **Multi-Dialect**: PostgreSQL, MySQL, SQLite desteği
/// - **Field Optimizasyonu**: Searchable, sortable, filterable alanlar için otomatik index
///
/// ## Kullanım Senaryoları
///
/// 1. **Yeni Proje Başlangıcı**: Tüm resource'ları kaydet ve AutoMigrate çağır
/// 2. **Geliştirme Ortamı**: Schema değişikliklerini otomatik uygula
/// 3. **Test Ortamı**: Her test öncesi temiz schema oluştur
/// 4. **Model Stub Üretimi**: Mevcut resource'lardan Go struct'ları oluştur
///
/// ## Avantajlar
///
/// - Kod tekrarını azaltır (DRY prensibi)
/// - Resource tanımları ile veritabanı şeması senkronize kalır
/// - Manuel migration yazma ihtiyacını ortadan kaldırır
/// - Performans optimizasyonları otomatik uygulanır
///
/// ## Dezavantajlar
///
/// - Production ortamında dikkatli kullanılmalı (veri kaybı riski)
/// - Karmaşık migration senaryoları için manuel müdahale gerekebilir
/// - Mevcut verileri migrate etmez, sadece schema değişikliği yapar
///
/// ## Önemli Notlar
///
/// ⚠️ **Production Uyarısı**: AutoMigrate production'da veri kaybına neden olabilir.
/// Production ortamında manuel migration'lar kullanın.
///
/// ⚠️ **Sıralama**: Resource'ları ilişki sırasına göre kaydedin (önce parent, sonra child).
///
/// ## Örnek Kullanım
///
/// ```go
/// // Migration generator oluştur
/// mg := migration.NewMigrationGenerator(db)
///
/// // Resource'ları kaydet
/// mg.RegisterResource(userResource).
///    RegisterResource(postResource).
///    RegisterResource(commentResource)
///
/// // Otomatik migration çalıştır
/// if err := mg.AutoMigrate(); err != nil {
///     log.Fatal(err)
/// }
///
/// // Model stub oluştur
/// stub := mg.GenerateModelStub(userResource)
/// fmt.Println(stub)
/// ```
type MigrationGenerator struct {
	db         *gorm.DB           // GORM veritabanı bağlantısı
	resources  []resource.Resource // Kayıtlı resource listesi
	typeMapper *TypeMapper         // Field type dönüştürücü
	dialect    string              // Veritabanı dialect'i (postgres, mysql, sqlite)
}

/// # NewMigrationGenerator
///
/// Bu fonksiyon, yeni bir MigrationGenerator instance'ı oluşturur ve veritabanı
/// dialect'ini otomatik olarak algılar.
///
/// ## Parametreler
///
/// - `db`: GORM veritabanı bağlantısı (*gorm.DB)
///
/// ## Döndürür
///
/// - Yapılandırılmış MigrationGenerator pointer'ı
///
/// ## Özellikler
///
/// - Veritabanı dialect'ini otomatik algılar (PostgreSQL, MySQL, SQLite)
/// - Dialect'e özel TypeMapper oluşturur
/// - Boş resource listesi ile başlatır
///
/// ## Kullanım Örneği
///
/// ```go
/// // PostgreSQL bağlantısı
/// db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
/// if err != nil {
///     log.Fatal(err)
/// }
///
/// // Migration generator oluştur
/// mg := migration.NewMigrationGenerator(db)
///
/// // MySQL için
/// mysqlDB, _ := gorm.Open(mysql.Open(dsn), &gorm.Config{})
/// mgMySQL := migration.NewMigrationGenerator(mysqlDB)
///
/// // SQLite için
/// sqliteDB, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
/// mgSQLite := migration.NewMigrationGenerator(sqliteDB)
/// ```
///
/// ## Önemli Notlar
///
/// ⚠️ Veritabanı bağlantısının açık ve geçerli olduğundan emin olun.
func NewMigrationGenerator(db *gorm.DB) *MigrationGenerator {
	dialect := db.Dialector.Name()
	return &MigrationGenerator{
		db:         db,
		resources:  []resource.Resource{},
		typeMapper: NewTypeMapperWithDialect(dialect),
		dialect:    dialect,
	}
}

/// # RegisterResource
///
/// Bu fonksiyon, migration işlemi için tek bir resource kaydeder ve method chaining
/// desteği sağlar.
///
/// ## Parametreler
///
/// - `r`: Kaydedilecek resource (resource.Resource interface'i)
///
/// ## Döndürür
///
/// - Yapılandırılmış MigrationGenerator pointer'ı (method chaining için)
///
/// ## Kullanım Senaryoları
///
/// 1. **Tekli Kayıt**: Bir resource'u kaydet
/// 2. **Method Chaining**: Birden fazla RegisterResource çağrısını zincirle
/// 3. **Dinamik Kayıt**: Runtime'da resource'ları kaydet
///
/// ## Kullanım Örneği
///
/// ```go
/// mg := migration.NewMigrationGenerator(db)
///
/// // Tekli kayıt
/// mg.RegisterResource(userResource)
///
/// // Method chaining ile çoklu kayıt
/// mg.RegisterResource(userResource).
///    RegisterResource(postResource).
///    RegisterResource(commentResource)
///
/// // Dinamik kayıt
/// for _, res := range dynamicResources {
///     mg.RegisterResource(res)
/// }
/// ```
///
/// ## Önemli Notlar
///
/// - Resource'lar kayıt sırasına göre migrate edilir
/// - İlişkili resource'ları doğru sırada kaydedin (önce parent, sonra child)
/// - Aynı resource birden fazla kez kaydedilebilir (dikkatli olun)
func (mg *MigrationGenerator) RegisterResource(r resource.Resource) *MigrationGenerator {
	mg.resources = append(mg.resources, r)
	return mg
}

/// # RegisterResources
///
/// Bu fonksiyon, birden fazla resource'u tek seferde kaydeder ve method chaining
/// desteği sağlar.
///
/// ## Parametreler
///
/// - `resources`: Kaydedilecek resource'lar (variadic parameter)
///
/// ## Döndürür
///
/// - Yapılandırılmış MigrationGenerator pointer'ı (method chaining için)
///
/// ## Kullanım Senaryoları
///
/// 1. **Toplu Kayıt**: Birden fazla resource'u tek çağrıda kaydet
/// 2. **Slice Expansion**: Resource slice'ını expand ederek kaydet
/// 3. **Modüler Kayıt**: Farklı modüllerden resource'ları toplu kaydet
///
/// ## Kullanım Örneği
///
/// ```go
/// mg := migration.NewMigrationGenerator(db)
///
/// // Toplu kayıt
/// mg.RegisterResources(
///     userResource,
///     postResource,
///     commentResource,
/// )
///
/// // Slice expansion
/// authResources := []resource.Resource{userResource, roleResource}
/// mg.RegisterResources(authResources...)
///
/// // Modüler kayıt
/// mg.RegisterResources(authModule.Resources()...).
///    RegisterResources(blogModule.Resources()...)
/// ```
///
/// ## Avantajlar
///
/// - Daha temiz ve okunabilir kod
/// - Tek satırda çoklu kayıt
/// - Slice'larla kolay entegrasyon
///
/// ## Önemli Notlar
///
/// - Resource'lar verilen sırada kaydedilir
/// - İlişki bağımlılıklarına dikkat edin
func (mg *MigrationGenerator) RegisterResources(resources ...resource.Resource) *MigrationGenerator {
	mg.resources = append(mg.resources, resources...)
	return mg
}

/// # AutoMigrate
///
/// Bu fonksiyon, kayıtlı tüm resource'ların modellerini otomatik olarak migrate eder.
/// GORM'un AutoMigrate özelliğini kullanarak tabloları oluşturur veya günceller,
/// ardından field tanımlarından ek constraint'leri uygular.
///
/// ## Döndürür
///
/// - `error`: Migration başarısız olursa hata, başarılıysa nil
///
/// ## İşlem Adımları
///
/// 1. Her resource için model kontrolü yapar
/// 2. GORM AutoMigrate ile tablo oluşturur/günceller
/// 3. Field constraint'lerini uygular (index, unique, foreign key)
/// 4. İlişkisel field'lar için ek işlemler yapar
///
/// ## Uygulanan Constraint'ler
///
/// - **Index**: Searchable, sortable, filterable alanlar için
/// - **Unique Index**: Unique validation rule'u olan alanlar için
/// - **Foreign Key Index**: BelongsTo ilişkileri için
/// - **Pivot Table**: BelongsToMany ilişkileri için
///
/// ## Kullanım Senaryoları
///
/// 1. **İlk Kurulum**: Yeni projede tüm tabloları oluştur
/// 2. **Geliştirme**: Schema değişikliklerini otomatik uygula
/// 3. **Test**: Her test öncesi temiz schema oluştur
/// 4. **CI/CD**: Otomatik deployment'ta schema güncelle
///
/// ## Kullanım Örneği
///
/// ```go
/// mg := migration.NewMigrationGenerator(db)
///
/// // Resource'ları kaydet
/// mg.RegisterResources(
///     userResource,
///     postResource,
///     commentResource,
/// )
///
/// // Migration çalıştır
/// if err := mg.AutoMigrate(); err != nil {
///     log.Fatalf("Migration failed: %v", err)
/// }
///
/// // Transaction içinde
/// err := db.Transaction(func(tx *gorm.DB) error {
///     mg := migration.NewMigrationGenerator(tx)
///     mg.RegisterResources(resources...)
///     return mg.AutoMigrate()
/// })
/// ```
///
/// ## Hata Durumları
///
/// - Resource'un model'i yoksa: "resource X has no model" hatası
/// - GORM migration başarısızsa: "migration failed for X" hatası
/// - Constraint uygulaması başarısızsa: "field constraints failed for X" hatası
///
/// ## Avantajlar
///
/// - Otomatik schema yönetimi
/// - Performans optimizasyonları (otomatik index'ler)
/// - İlişki yönetimi (pivot tablolar, foreign key'ler)
/// - Dialect-aware (PostgreSQL, MySQL, SQLite)
///
/// ## Dezavantajlar
///
/// - Mevcut verileri migrate etmez
/// - Karmaşık migration'lar için yetersiz kalabilir
/// - Production'da veri kaybı riski
///
/// ## Önemli Notlar
///
/// ⚠️ **PRODUCTION UYARISI**: AutoMigrate production ortamında dikkatli kullanılmalıdır.
/// Sütun silme, tip değiştirme gibi işlemler veri kaybına neden olabilir.
///
/// ⚠️ **SIRA ÖNEMLİ**: Resource'ları ilişki bağımlılıklarına göre sıralayın.
/// Önce parent resource'ları, sonra child resource'ları kaydedin.
///
/// ⚠️ **TRANSACTION**: Kritik ortamlarda transaction içinde çalıştırın.
///
/// ## Best Practices
///
/// ```go
/// // Development ortamı
/// if os.Getenv("ENV") == "development" {
///     if err := mg.AutoMigrate(); err != nil {
///         log.Fatal(err)
///     }
/// }
///
/// // Production ortamı - manuel migration kullanın
/// if os.Getenv("ENV") == "production" {
///     log.Println("Use manual migrations in production")
///     // Manuel migration dosyalarını çalıştır
/// }
/// ```
func (mg *MigrationGenerator) AutoMigrate() error {
	for _, r := range mg.resources {
		model := r.Model()
		if model == nil {
			return fmt.Errorf("resource %s has no model - all resources must have a model for migration", r.Slug())
		}

		// GORM AutoMigrate
		if err := mg.db.AutoMigrate(model); err != nil {
			return fmt.Errorf("migration failed for %s: %w", r.Slug(), err)
		}

		// Field constraint'lerini uygula
		if err := mg.applyFieldConstraints(r); err != nil {
			return fmt.Errorf("field constraints failed for %s: %w", r.Slug(), err)
		}
	}
	return nil
}

/// # applyFieldConstraints
///
/// Bu fonksiyon, resource'un field tanımlarından ek veritabanı constraint'lerini
/// otomatik olarak oluşturur ve uygular. GORM'un AutoMigrate'inin yapmadığı
/// optimizasyonları ve constraint'leri ekler.
///
/// ## Parametreler
///
/// - `r`: Constraint'leri uygulanacak resource
///
/// ## Döndürür
///
/// - `error`: İşlem başarısız olursa hata, başarılıysa nil
///
/// ## Uygulanan Constraint'ler
///
/// ### İlişkisel Field'lar
/// - **BelongsTo**: Foreign key için index oluşturur
/// - **BelongsToMany**: Pivot tablo oluşturur
/// - **HasOne/HasMany**: İlişki yapılandırması
///
/// ### Normal Field'lar
/// - **GlobalSearch**: Arama yapılabilir alanlar için index
/// - **IsSortable**: Sıralanabilir alanlar için index
/// - **IsFilterable**: Filtrelenebilir alanlar için index
/// - **UniqueIndex**: GormConfig'den unique index
/// - **Index**: GormConfig'den normal index
/// - **Unique Validation**: Validation rules'dan unique constraint
///
/// ## İşlem Akışı
///
/// ```
/// 1. Resource'un tüm field'larını tara
/// 2. İlişkisel field'ları kontrol et
///    - BelongsTo → Foreign key index
///    - BelongsToMany → Pivot tablo
/// 3. Normal field'ları kontrol et
///    - Searchable → Index
///    - Sortable → Index
///    - Filterable → Index
///    - GormConfig → Index/Unique
///    - Validation → Unique constraint
/// 4. Mevcut constraint'leri kontrol et (duplicate önleme)
/// 5. Yeni constraint'leri oluştur
/// ```
///
/// ## Kullanım Örneği
///
/// ```go
/// // Otomatik olarak AutoMigrate içinde çağrılır
/// // Manuel kullanım:
/// if err := mg.applyFieldConstraints(userResource); err != nil {
///     log.Fatal(err)
/// }
/// ```
///
/// ## Performans Optimizasyonları
///
/// - Duplicate index kontrolü yapar
/// - Sadece gerekli index'leri oluşturur
/// - Dialect-aware SQL kullanır
///
/// ## Önemli Notlar
///
/// - Bu fonksiyon internal'dır, doğrudan çağrılmamalıdır
/// - AutoMigrate tarafından otomatik çağrılır
/// - Index oluşturma işlemi zaman alabilir (büyük tablolarda)
func (mg *MigrationGenerator) applyFieldConstraints(r resource.Resource) error {
	model := r.Model()
	if model == nil {
		return fmt.Errorf("resource %s has no model", r.Slug())
	}

	for _, field := range r.Fields() {
		// İlişkisel field'ları kontrol et
		if relField, ok := fields.IsRelationshipField(field); ok {
			// BelongsTo için foreign key index'i
			if relField.GetRelationshipType() == "belongsTo" {
				if bt, ok := relField.(*fields.BelongsToField); ok {
					if bt.GormRelationConfig != nil && bt.GormRelationConfig.ForeignKey != "" {
						fkColumn := bt.GormRelationConfig.ForeignKey
						if !mg.hasIndexWithModel(model, fkColumn) {
							if err := mg.createIndexWithModel(model, fkColumn, false); err != nil {
								return err
							}
						}
					}
				}
			}

			// BelongsToMany için pivot tablo
			if relField.GetRelationshipType() == "belongsToMany" {
				if btm, ok := relField.(*fields.BelongsToManyField); ok {
					if err := mg.createPivotTable(btm); err != nil {
						return err
					}
				}
			}

			continue
		}

		// Normal field'lar için mevcut logic
		schema, ok := field.(*fields.Schema)
		if !ok {
			continue
		}

		// Searchable alanlar için index
		if schema.GlobalSearch && !mg.hasIndexWithModel(model, schema.Key) {
			if err := mg.createIndexWithModel(model, schema.Key, false); err != nil {
				return err
			}
		}

		// Sortable alanlar için index
		if schema.IsSortable && !mg.hasIndexWithModel(model, schema.Key) {
			if err := mg.createIndexWithModel(model, schema.Key, false); err != nil {
				return err
			}
		}

		// Filterable alanlar için index
		if schema.IsFilterable && !mg.hasIndexWithModel(model, schema.Key) {
			if err := mg.createIndexWithModel(model, schema.Key, false); err != nil {
				return err
			}
		}

		// GormConfig'den constraint'ler
		if schema.HasGormConfig() {
			config := schema.GetGormConfig()

			// Unique Index
			if config.UniqueIndex && !mg.hasUniqueIndexWithModel(model, schema.Key) {
				if err := mg.createIndexWithModel(model, schema.Key, true); err != nil {
					return err
				}
			}

			// Normal Index
			if config.Index && !mg.hasIndexWithModel(model, schema.Key) {
				if err := mg.createIndexWithModel(model, schema.Key, false); err != nil {
					return err
				}
			}
		}

		// Validation rules'dan unique constraint
		for _, rule := range schema.ValidationRules {
			if rule.Name == "unique" {
				if !mg.hasUniqueIndexWithModel(model, schema.Key) {
					if err := mg.createIndexWithModel(model, schema.Key, true); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

/// # getTableName
///
/// Bu fonksiyon, resource'dan veritabanı tablo adını çıkarır. GORM'un
/// NamingStrategy'sini kullanarak doğru tablo adını belirler.
///
/// ## Parametreler
///
/// - `r`: Tablo adı çıkarılacak resource
///
/// ## Döndürür
///
/// - `string`: Veritabanı tablo adı
///
/// ## İşlem Akışı
///
/// 1. Resource'un model'ini al
/// 2. Model varsa GORM'un NamingStrategy'sini kullan
/// 3. Model yoksa slug'dan tablo adı türet (snake_case)
///
/// ## Kullanım Örneği
///
/// ```go
/// tableName := mg.getTableName(userResource)
/// // Çıktı: "users"
///
/// tableName := mg.getTableName(blogPostResource)
/// // Çıktı: "blog_posts"
/// ```
///
/// ## Önemli Notlar
///
/// - GORM'un naming convention'ını takip eder
/// - Fallback olarak slug'dan türetme yapar
/// - Snake case formatında döner
func (mg *MigrationGenerator) getTableName(r resource.Resource) string {
	// GORM'dan gerçek tablo adını al
	model := r.Model()
	if model == nil {
		// Fallback: slug'dan tablo adı türet
		slug := r.Slug()
		return strcase.ToSnake(slug)
	}

	// GORM'un NamingStrategy'sini kullanarak tablo adını al
	stmt := &gorm.Statement{DB: mg.db}
	err := stmt.Parse(model)
	if err != nil {
		// Fallback: slug'dan tablo adı türet
		slug := r.Slug()
		return strcase.ToSnake(slug)
	}

	return stmt.Table
}

/// # hasIndex
///
/// Bu fonksiyon, belirtilen tabloda belirtilen sütun için index'in var olup olmadığını
/// kontrol eder. Dialect-aware SQL sorguları kullanarak veritabanına özgü kontrol yapar.
///
/// ## Parametreler
///
/// - `table`: Tablo adı
/// - `column`: Sütun adı
///
/// ## Döndürür
///
/// - `bool`: Index varsa true, yoksa false
///
/// ## Desteklenen Dialect'ler
///
/// - **PostgreSQL**: pg_indexes system catalog kullanır
/// - **MySQL**: information_schema.statistics kullanır
/// - **SQLite**: sqlite_master kullanır
///
/// ## Index Adlandırma Kuralı
///
/// Index adı: `idx_{table}_{column}` formatında oluşturulur
///
/// ## Kullanım Örneği
///
/// ```go
/// // Index kontrolü
/// exists := mg.hasIndex("users", "email")
/// if !exists {
///     mg.createIndex("users", "email", false)
/// }
/// ```
///
/// ## Önemli Notlar
///
/// - Bu fonksiyon internal'dır, doğrudan çağrılmamalıdır
/// - Model instance gerektirmez, sadece tablo adı kullanır
/// - hasIndexWithModel tercih edilmelidir (GORM Migrator kullanır)
func (mg *MigrationGenerator) hasIndex(table, column string) bool {
	indexName := fmt.Sprintf("idx_%s_%s", table, column)

	// GORM Migrator kullanarak - dialect-aware ve güvenli
	// Not: Migrator.HasIndex() model instance gerektirir, ama biz sadece tablo adını biliyoruz
	// Bu yüzden GORM'nin internal query'lerini kullanmak zorundayız
	// Alternatif: Her resource için model instance'ı kullan

	// Daha iyi yaklaşım: GORM'nin kendi index kontrolünü kullan
	// Ama bu durumda model instance gerekiyor
	// Şimdilik EXISTS query kullanıyoruz (daha performanslı)
	var exists bool

	switch mg.dialect {
	case "postgres":
		mg.db.Raw("SELECT EXISTS(SELECT 1 FROM pg_indexes WHERE schemaname = 'public' AND tablename = ? AND indexname = ?)",
			table, indexName).Scan(&exists)
	case "mysql":
		mg.db.Raw("SELECT EXISTS(SELECT 1 FROM information_schema.statistics WHERE table_schema = DATABASE() AND table_name = ? AND index_name = ?)",
			table, indexName).Scan(&exists)
	default:
		// SQLite için
		var count int64
		mg.db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='index' AND name=?", indexName).Scan(&count)
		exists = count > 0
	}

	return exists
}

/// # hasIndexWithModel
///
/// Bu fonksiyon, model instance kullanarak tabloda index'in var olup olmadığını
/// kontrol eder. GORM Migrator API'sini kullanarak dialect-aware ve güvenli
/// kontrol yapar.
///
/// ## Parametreler
///
/// - `model`: GORM model instance'ı (interface{})
/// - `column`: Sütun adı
///
/// ## Döndürür
///
/// - `bool`: Index varsa true, yoksa false
///
/// ## Avantajlar
///
/// - GORM Migrator API kullanır (güvenli ve dialect-aware)
/// - Model'den tablo adını otomatik çıkarır
/// - Tüm veritabanı dialect'lerini destekler
///
/// ## Index Adlandırma Kuralı
///
/// Index adı: `idx_{table}_{column}` formatında oluşturulur
///
/// ## Kullanım Örneği
///
/// ```go
/// // Model ile index kontrolü
/// user := &User{}
/// exists := mg.hasIndexWithModel(user, "email")
/// if !exists {
///     mg.createIndexWithModel(user, "email", false)
/// }
/// ```
///
/// ## Önemli Notlar
///
/// - hasIndex yerine bu fonksiyon tercih edilmelidir
/// - GORM'un internal API'sini kullanır
/// - Model instance gerektirir
func (mg *MigrationGenerator) hasIndexWithModel(model interface{}, column string) bool {
	tableName := mg.getTableNameFromModel(model)
	indexName := fmt.Sprintf("idx_%s_%s", tableName, column)

	// GORM Migrator kullanarak - dialect-aware ve güvenli
	return mg.db.Migrator().HasIndex(model, indexName)
}

/// # hasUniqueIndexWithModel
///
/// Bu fonksiyon, model instance kullanarak tabloda unique index'in var olup
/// olmadığını kontrol eder. GORM Migrator API'sini kullanarak dialect-aware
/// ve güvenli kontrol yapar.
///
/// ## Parametreler
///
/// - `model`: GORM model instance'ı (interface{})
/// - `column`: Sütun adı
///
/// ## Döndürür
///
/// - `bool`: Unique index varsa true, yoksa false
///
/// ## Unique Index Adlandırma Kuralı
///
/// Index adı: `uniq_{table}_{column}` formatında oluşturulur
///
/// ## Kullanım Örneği
///
/// ```go
/// // Unique index kontrolü
/// user := &User{}
/// exists := mg.hasUniqueIndexWithModel(user, "email")
/// if !exists {
///     mg.createIndexWithModel(user, "email", true)
/// }
/// ```
///
/// ## Önemli Notlar
///
/// - Normal index'ten farklı adlandırma kullanır (uniq_ prefix)
/// - GORM Migrator API kullanır
/// - Unique constraint kontrolü yapar
func (mg *MigrationGenerator) hasUniqueIndexWithModel(model interface{}, column string) bool {
	tableName := mg.getTableNameFromModel(model)
	indexName := fmt.Sprintf("uniq_%s_%s", tableName, column)

	// GORM Migrator kullanarak - dialect-aware ve güvenli
	return mg.db.Migrator().HasIndex(model, indexName)
}

/// # createIndexWithModel
///
/// Bu fonksiyon, model instance kullanarak tabloda index veya unique index
/// oluşturur. GORM Migrator API'sini kullanarak dialect-aware ve güvenli
/// index oluşturma işlemi yapar.
///
/// ## Parametreler
///
/// - `model`: GORM model instance'ı (interface{})
/// - `column`: Index oluşturulacak sütun adı
/// - `unique`: true ise unique index, false ise normal index oluşturur
///
/// ## Döndürür
///
/// - `error`: İşlem başarısız olursa hata, başarılıysa nil
///
/// ## Index Türleri
///
/// ### Normal Index
/// - Adlandırma: `idx_{table}_{column}`
/// - Kullanım: Arama, sıralama, filtreleme performansı
/// - GORM Migrator.CreateIndex kullanır
///
/// ### Unique Index
/// - Adlandırma: `uniq_{table}_{column}`
/// - Kullanım: Benzersizlik constraint'i
/// - Manuel SQL kullanır (GORM'da unique index metodu yok)
///
/// ## Kullanım Örneği
///
/// ```go
/// user := &User{}
///
/// // Normal index oluştur
/// if err := mg.createIndexWithModel(user, "email", false); err != nil {
///     log.Fatal(err)
/// }
///
/// // Unique index oluştur
/// if err := mg.createIndexWithModel(user, "username", true); err != nil {
///     log.Fatal(err)
/// }
/// ```
///
/// ## Dialect-Aware SQL
///
/// Unique index için dialect'e özel SQL kullanır:
/// - PostgreSQL: `CREATE UNIQUE INDEX IF NOT EXISTS ...`
/// - MySQL: `CREATE UNIQUE INDEX IF NOT EXISTS ...`
/// - SQLite: `CREATE UNIQUE INDEX IF NOT EXISTS ...`
///
/// ## Önemli Notlar
///
/// - Mevcut index kontrolü yapmaz, çağıran tarafından kontrol edilmelidir
/// - IF NOT EXISTS kullanır (duplicate hata önleme)
/// - Unique index için manuel SQL gerekir (GORM limitasyonu)
func (mg *MigrationGenerator) createIndexWithModel(model interface{}, column string, unique bool) error {
	tableName := mg.getTableNameFromModel(model)

	if unique {
		// Unique index için - GORM Migrator'da unique index oluşturmak için özel bir yöntem yok
		// Bu yüzden manuel SQL kullanıyoruz (dialect-aware)
		indexName := fmt.Sprintf("uniq_%s_%s", tableName, column)
		indexType := "UNIQUE INDEX"
		sql := fmt.Sprintf("CREATE %s IF NOT EXISTS %s ON %s(%s)", indexType, indexName, tableName, column)
		return mg.db.Exec(sql).Error
	}

	// Normal index için GORM Migrator kullan
	indexName := fmt.Sprintf("idx_%s_%s", tableName, column)

	// GORM Migrator'ın CreateIndex metodu field adı veya index adı alabilir
	// Biz index adını kullanıyoruz
	return mg.db.Migrator().CreateIndex(model, indexName)
}

/// # getTableNameFromModel
///
/// Bu fonksiyon, GORM model instance'ından veritabanı tablo adını çıkarır.
/// GORM'un NamingStrategy'sini kullanarak doğru tablo adını belirler.
///
/// ## Parametreler
///
/// - `model`: GORM model instance'ı (interface{})
///
/// ## Döndürür
///
/// - `string`: Veritabanı tablo adı (boş string hata durumunda)
///
/// ## İşlem Akışı
///
/// 1. GORM Statement oluştur
/// 2. Model'i parse et
/// 3. Statement'tan tablo adını al
/// 4. Hata durumunda boş string döner
///
/// ## Kullanım Örneği
///
/// ```go
/// user := &User{}
/// tableName := mg.getTableNameFromModel(user)
/// // Çıktı: "users"
///
/// blogPost := &BlogPost{}
/// tableName := mg.getTableNameFromModel(blogPost)
/// // Çıktı: "blog_posts"
/// ```
///
/// ## GORM NamingStrategy
///
/// GORM'un default naming strategy'si:
/// - Struct adını plural yapar (User → users)
/// - CamelCase'i snake_case'e çevirir (BlogPost → blog_posts)
/// - Custom TableName() metodu varsa onu kullanır
///
/// ## Önemli Notlar
///
/// - GORM'un internal API'sini kullanır
/// - Custom TableName() metodunu destekler
/// - Hata durumunda boş string döner (error dönmez)
func (mg *MigrationGenerator) getTableNameFromModel(model interface{}) string {
	stmt := &gorm.Statement{DB: mg.db}
	err := stmt.Parse(model)
	if err != nil {
		return ""
	}
	return stmt.Table
}

/// # hasUniqueIndex
///
/// Bu fonksiyon, belirtilen tabloda belirtilen sütun için unique index'in
/// var olup olmadığını kontrol eder. Dialect-aware SQL sorguları kullanarak
/// veritabanına özgü kontrol yapar.
///
/// ## Parametreler
///
/// - `table`: Tablo adı
/// - `column`: Sütun adı
///
/// ## Döndürür
///
/// - `bool`: Unique index varsa true, yoksa false
///
/// ## Desteklenen Dialect'ler
///
/// ### PostgreSQL
/// - pg_indexes ve pg_index system catalog kullanır
/// - indisunique flag'ini kontrol eder
///
/// ### MySQL
/// - information_schema.statistics kullanır
/// - non_unique = 0 kontrolü yapar
///
/// ### SQLite
/// - sqlite_master kullanır
/// - Index adından kontrol yapar
///
/// ## Unique Index Adlandırma Kuralı
///
/// Index adı: `uniq_{table}_{column}` formatında oluşturulur
///
/// ## Kullanım Örneği
///
/// ```go
/// // Unique index kontrolü
/// exists := mg.hasUniqueIndex("users", "email")
/// if !exists {
///     mg.createIndex("users", "email", true)
/// }
/// ```
///
/// ## Önemli Notlar
///
/// - Bu fonksiyon internal'dır, doğrudan çağrılmamalıdır
/// - hasUniqueIndexWithModel tercih edilmelidir (GORM Migrator kullanır)
/// - Dialect'e özel SQL sorguları kullanır
func (mg *MigrationGenerator) hasUniqueIndex(table, column string) bool {
	indexName := fmt.Sprintf("uniq_%s_%s", table, column)
	var count int64

	switch mg.dialect {
	case "postgres":
		// PostgreSQL için - unique index kontrolü
		mg.db.Raw(`SELECT COUNT(*) FROM pg_indexes i
			JOIN pg_class c ON i.indexname = c.relname
			JOIN pg_index idx ON c.oid = idx.indexrelid
			WHERE i.schemaname = 'public' AND i.tablename = ? AND i.indexname = ? AND idx.indisunique = true`,
			table, indexName).Scan(&count)
	case "mysql":
		// MySQL için - unique index kontrolü
		mg.db.Raw(`SELECT COUNT(*) FROM information_schema.statistics
			WHERE table_schema = DATABASE() AND table_name = ? AND index_name = ? AND non_unique = 0`,
			table, indexName).Scan(&count)
	default:
		// SQLite için (default)
		mg.db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='index' AND name=?", indexName).Scan(&count)
	}

	return count > 0
}

/// # createIndex
///
/// Bu fonksiyon, belirtilen tabloda belirtilen sütun için index veya unique
/// index oluşturur. Manuel SQL kullanarak dialect-aware index oluşturma
/// işlemi yapar.
///
/// ## Parametreler
///
/// - `table`: Tablo adı
/// - `column`: Index oluşturulacak sütun adı
/// - `unique`: true ise unique index, false ise normal index oluşturur
///
/// ## Döndürür
///
/// - `error`: İşlem başarısız olursa hata, başarılıysa nil
///
/// ## Index Türleri
///
/// ### Normal Index
/// - Adlandırma: `idx_{table}_{column}`
/// - SQL: `CREATE INDEX IF NOT EXISTS idx_table_column ON table(column)`
///
/// ### Unique Index
/// - Adlandırma: `uniq_{table}_{column}`
/// - SQL: `CREATE UNIQUE INDEX IF NOT EXISTS uniq_table_column ON table(column)`
///
/// ## Kullanım Örneği
///
/// ```go
/// // Normal index oluştur
/// if err := mg.createIndex("users", "email", false); err != nil {
///     log.Fatal(err)
/// }
///
/// // Unique index oluştur
/// if err := mg.createIndex("users", "username", true); err != nil {
///     log.Fatal(err)
/// }
/// ```
///
/// ## SQL Format
///
/// ```sql
/// -- Normal index
/// CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)
///
/// -- Unique index
/// CREATE UNIQUE INDEX IF NOT EXISTS uniq_users_username ON users(username)
/// ```
///
/// ## Önemli Notlar
///
/// - IF NOT EXISTS kullanır (duplicate hata önleme)
/// - Tüm dialect'lerde çalışır
/// - createIndexWithModel tercih edilmelidir (GORM Migrator kullanır)
func (mg *MigrationGenerator) createIndex(table, column string, unique bool) error {
	indexType := "INDEX"
	indexPrefix := "idx"
	if unique {
		indexType = "UNIQUE INDEX"
		indexPrefix = "uniq"
	}

	indexName := fmt.Sprintf("%s_%s_%s", indexPrefix, table, column)
	sql := fmt.Sprintf("CREATE %s IF NOT EXISTS %s ON %s(%s)", indexType, indexName, table, column)

	return mg.db.Exec(sql).Error
}

/// # createPivotTable
///
/// Bu fonksiyon, BelongsToMany (Many-to-Many) ilişkileri için pivot tablo
/// oluşturur. Pivot tablo, iki model arasındaki çoktan-çoğa ilişkiyi
/// yönetmek için kullanılır.
///
/// ## Parametreler
///
/// - `btm`: BelongsToManyField instance'ı (pivot tablo bilgilerini içerir)
///
/// ## Döndürür
///
/// - `error`: İşlem başarısız olursa hata, başarılıysa nil
///
/// ## Pivot Tablo Yapısı
///
/// ```sql
/// CREATE TABLE pivot_table (
///     foreign_key_column BIGINT NOT NULL,  -- İlk model'in ID'si
///     related_key_column BIGINT NOT NULL,  -- İkinci model'in ID'si
///     PRIMARY KEY (foreign_key_column, related_key_column)
/// )
/// ```
///
/// ## İşlem Adımları
///
/// 1. Pivot tablo var mı kontrol et (dialect-aware)
/// 2. Tablo yoksa oluştur (composite primary key ile)
/// 3. Her iki sütun için index oluştur (performans için)
///
/// ## Dialect-Aware SQL
///
/// ### PostgreSQL & MySQL
/// - BIGINT veri tipi kullanır
/// - Composite primary key destekler
///
/// ### SQLite
/// - INTEGER veri tipi kullanır
/// - Composite primary key destekler
///
/// ## Kullanım Örneği
///
/// ```go
/// // BelongsToMany field tanımı
/// rolesField := fields.NewBelongsToManyField("roles", "Role").
///     SetPivotTable("user_roles").
///     SetForeignKey("user_id").
///     SetRelatedKey("role_id")
///
/// // Pivot tablo oluştur (otomatik olarak applyFieldConstraints içinde çağrılır)
/// if err := mg.createPivotTable(rolesField); err != nil {
///     log.Fatal(err)
/// }
/// ```
///
/// ## Oluşturulan Index'ler
///
/// - `idx_pivot_table_foreign_key`: Foreign key sütunu için
/// - `idx_pivot_table_related_key`: Related key sütunu için
///
/// ## Kullanım Senaryoları
///
/// 1. **User-Role İlişkisi**: Kullanıcılar birden fazla role sahip olabilir
/// 2. **Post-Tag İlişkisi**: Gönderiler birden fazla etikete sahip olabilir
/// 3. **Student-Course İlişkisi**: Öğrenciler birden fazla kursa kayıtlı olabilir
///
/// ## Önemli Notlar
///
/// ⚠️ Pivot tablo zaten varsa işlem atlanır (duplicate önleme)
/// ⚠️ Composite primary key kullanır (duplicate ilişki önleme)
/// ⚠️ Her iki sütun için index oluşturur (join performansı)
///
/// ## Best Practices
///
/// - Pivot tablo adı: `{model1}_{model2}` formatında olmalı
/// - Foreign key: `{model1}_id` formatında olmalı
/// - Related key: `{model2}_id` formatında olmalı
func (mg *MigrationGenerator) createPivotTable(btm *fields.BelongsToManyField) error {
	var count int64

	// Pivot tablo zaten var mı kontrol et
	switch mg.dialect {
	case "postgres":
		mg.db.Raw("SELECT COUNT(*) FROM pg_tables WHERE schemaname = 'public' AND tablename = ?", btm.PivotTableName).Scan(&count)
	case "mysql":
		mg.db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?", btm.PivotTableName).Scan(&count)
	default:
		// SQLite için
		mg.db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?", btm.PivotTableName).Scan(&count)
	}

	if count > 0 {
		return nil // Tablo zaten var
	}

	// Pivot tablo oluştur - dialect-aware SQL
	var sql string
	switch mg.dialect {
	case "postgres":
		sql = fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s (
				%s BIGINT NOT NULL,
				%s BIGINT NOT NULL,
				PRIMARY KEY (%s, %s)
			)
		`, btm.PivotTableName, btm.ForeignKeyColumn, btm.RelatedKeyColumn, btm.ForeignKeyColumn, btm.RelatedKeyColumn)
	case "mysql":
		sql = fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s (
				%s BIGINT NOT NULL,
				%s BIGINT NOT NULL,
				PRIMARY KEY (%s, %s)
			)
		`, btm.PivotTableName, btm.ForeignKeyColumn, btm.RelatedKeyColumn, btm.ForeignKeyColumn, btm.RelatedKeyColumn)
	default:
		// SQLite için (INTEGER)
		sql = fmt.Sprintf(`
			CREATE TABLE IF NOT EXISTS %s (
				%s INTEGER NOT NULL,
				%s INTEGER NOT NULL,
				PRIMARY KEY (%s, %s)
			)
		`, btm.PivotTableName, btm.ForeignKeyColumn, btm.RelatedKeyColumn, btm.ForeignKeyColumn, btm.RelatedKeyColumn)
	}

	if err := mg.db.Exec(sql).Error; err != nil {
		return fmt.Errorf("failed to create pivot table %s: %w", btm.PivotTableName, err)
	}

	// Index'ler ekle
	if err := mg.createIndex(btm.PivotTableName, btm.ForeignKeyColumn, false); err != nil {
		return err
	}
	if err := mg.createIndex(btm.PivotTableName, btm.RelatedKeyColumn, false); err != nil {
		return err
	}

	return nil
}

/// # FieldInfo
///
/// Bu yapı, resource field'larının detaylı bilgilerini içerir. Migration ve
/// model stub oluşturma işlemlerinde kullanılır.
///
/// ## Özellikler
///
/// ### Temel Bilgiler
/// - **Name**: Field'ın görünen adı (human-readable)
/// - **Key**: Field'ın veritabanı sütun adı (snake_case)
/// - **GoType**: Go dilinde field'ın tipi (örn: "string", "*User", "[]Post")
/// - **SQLType**: Veritabanında field'ın tipi (örn: "VARCHAR(255)", "BIGINT")
/// - **GormTag**: GORM tag string'i (örn: "column:email;type:varchar;not null")
///
/// ### Constraint Bilgileri
/// - **IsRequired**: Field zorunlu mu (NOT NULL)
/// - **IsNullable**: Field nullable mı
/// - **IsSearchable**: Field aranabilir mi (index gerektirir)
/// - **IsSortable**: Field sıralanabilir mi (index gerektirir)
/// - **IsFilterable**: Field filtrelenebilir mi (index gerektirir)
///
/// ### İlişki Bilgileri
/// - **IsRelation**: Field bir ilişki mi
/// - **RelationType**: İlişki tipi ("belongsTo", "hasOne", "hasMany", "belongsToMany")
/// - **RelatedResource**: İlişkili resource'un slug'ı
/// - **ForeignKey**: Foreign key sütun adı (BelongsTo, HasOne, HasMany için)
/// - **PivotTable**: Pivot tablo adı (BelongsToMany için)
/// - **RelationGormTag**: İlişki için GORM tag'i
///
/// ## Kullanım Senaryoları
///
/// 1. **Model Stub Oluşturma**: Go struct'ları otomatik oluştur
/// 2. **Migration Analizi**: Field'ların constraint'lerini analiz et
/// 3. **Dokümantasyon**: Field'ların özelliklerini dokümante et
/// 4. **Validasyon**: Field tanımlarını doğrula
///
/// ## Kullanım Örneği
///
/// ```go
/// // Field bilgilerini al
/// infos := mg.GetFieldInfos(userResource)
///
/// for _, info := range infos {
///     fmt.Printf("Field: %s\n", info.Name)
///     fmt.Printf("  Type: %s (SQL: %s)\n", info.GoType, info.SQLType)
///     fmt.Printf("  Required: %v\n", info.IsRequired)
///
///     if info.IsRelation {
///         fmt.Printf("  Relation: %s to %s\n",
///             info.RelationType, info.RelatedResource)
///     }
/// }
/// ```
///
/// ## İlişki Tipleri ve Örnekler
///
/// ### BelongsTo
/// ```go
/// FieldInfo{
///     Name: "Author",
///     Key: "author",
///     GoType: "*User",
///     IsRelation: true,
///     RelationType: "belongsTo",
///     RelatedResource: "users",
///     ForeignKey: "author_id",
///     RelationGormTag: "foreignKey:AuthorID",
/// }
/// ```
///
/// ### HasMany
/// ```go
/// FieldInfo{
///     Name: "Posts",
///     Key: "posts",
///     GoType: "[]Post",
///     IsRelation: true,
///     RelationType: "hasMany",
///     RelatedResource: "posts",
///     ForeignKey: "user_id",
///     RelationGormTag: "foreignKey:UserID",
/// }
/// ```
///
/// ### BelongsToMany
/// ```go
/// FieldInfo{
///     Name: "Roles",
///     Key: "roles",
///     GoType: "[]*Role",
///     IsRelation: true,
///     RelationType: "belongsToMany",
///     RelatedResource: "roles",
///     PivotTable: "user_roles",
///     RelationGormTag: "many2many:user_roles",
/// }
/// ```
///
/// ## Önemli Notlar
///
/// - Normal field'lar için RelationType boş string'dir
/// - İlişkisel field'lar için GoType otomatik türetilir
/// - GORM tag'leri dialect-aware oluşturulur
type FieldInfo struct {
	Name         string // Field'ın görünen adı
	Key          string // Veritabanı sütun adı
	GoType       string // Go dilinde field tipi
	SQLType      string // Veritabanında field tipi
	GormTag      string // GORM tag string'i
	IsRequired   bool   // Field zorunlu mu
	IsNullable   bool   // Field nullable mı
	IsSearchable bool   // Field aranabilir mi
	IsSortable   bool   // Field sıralanabilir mi
	IsFilterable bool   // Field filtrelenebilir mi
	IsRelation   bool   // Field bir ilişki mi
	RelationType string // İlişki tipi

	// İlişki Bilgileri
	RelatedResource  string // İlişkili resource slug'ı
	ForeignKey       string // Foreign key sütunu
	PivotTable       string // Pivot tablo adı (BelongsToMany için)
	RelationGormTag  string // İlişki için GORM tag'i
}

/// # GetFieldInfos
///
/// Bu fonksiyon, resource'un tüm field'larının detaylı bilgilerini döner.
/// Hem normal field'ları hem de ilişkisel field'ları analiz eder ve
/// FieldInfo slice'ı olarak döner.
///
/// ## Parametreler
///
/// - `r`: Bilgileri çıkarılacak resource
///
/// ## Döndürür
///
/// - `[]FieldInfo`: Field bilgileri slice'ı
///
/// ## İşlem Akışı
///
/// 1. Resource'un tüm field'larını tara
/// 2. Her field için tip kontrolü yap
///    - İlişkisel field mi? → buildRelationshipFieldInfo çağır
///    - Normal field mi? → Field bilgilerini topla
/// 3. Go type, SQL type ve GORM tag oluştur
/// 4. Constraint bilgilerini ekle
/// 5. FieldInfo slice'ına ekle
///
/// ## Toplanan Bilgiler
///
/// ### Normal Field'lar İçin
/// - Field adı ve key
/// - Go ve SQL tipleri
/// - GORM tag (column, type, constraints)
/// - Required, nullable, searchable, sortable, filterable flag'leri
///
/// ### İlişkisel Field'lar İçin
/// - İlişki tipi (belongsTo, hasOne, hasMany, belongsToMany)
/// - İlişkili resource
/// - Foreign key veya pivot tablo bilgisi
/// - İlişki GORM tag'i
/// - Go type (pointer, slice, vb.)
///
/// ## Kullanım Örneği
///
/// ```go
/// // User resource için field bilgilerini al
/// infos := mg.GetFieldInfos(userResource)
///
/// // Field'ları listele
/// for _, info := range infos {
///     if info.IsRelation {
///         fmt.Printf("Relation: %s -> %s (%s)\n",
///             info.Key, info.RelatedResource, info.RelationType)
///     } else {
///         fmt.Printf("Field: %s (%s)\n", info.Key, info.SQLType)
///     }
/// }
///
/// // Model stub oluşturmak için kullan
/// stub := mg.GenerateModelStub(userResource)
///
/// // Searchable field'ları bul
/// searchableFields := []string{}
/// for _, info := range infos {
///     if info.IsSearchable {
///         searchableFields = append(searchableFields, info.Key)
///     }
/// }
/// ```
///
/// ## Kullanım Senaryoları
///
/// 1. **Model Stub Oluşturma**: GenerateModelStub içinde kullanılır
/// 2. **Migration Analizi**: Hangi constraint'lerin gerekli olduğunu belirle
/// 3. **API Dokümantasyonu**: Field'ların özelliklerini dokümante et
/// 4. **Validasyon**: Field tanımlarını doğrula
/// 5. **Code Generation**: Otomatik kod üretimi için
///
/// ## Önemli Notlar
///
/// - Schema olmayan field'lar atlanır
/// - İlişkisel field'lar özel işleme tabi tutulur
/// - Go type'lar nullable duruma göre pointer olabilir
/// - SQL type'lar dialect'e göre değişir
func (mg *MigrationGenerator) GetFieldInfos(r resource.Resource) []FieldInfo {
	var infos []FieldInfo

	for _, field := range r.Fields() {
		// İlişkisel field'ları kontrol et
		if relField, ok := fields.IsRelationshipField(field); ok {
			info := mg.buildRelationshipFieldInfo(relField)
			infos = append(infos, info)
			continue
		}

		// Normal field'lar
		schema, ok := field.(*fields.Schema)
		if !ok {
			continue
		}

		info := FieldInfo{
			Name:         schema.Name,
			Key:          schema.Key,
			SQLType:      mg.typeMapper.MapFieldTypeToSQL(schema.Type, 0),
			IsRequired:   schema.IsRequired,
			IsNullable:   schema.IsNullable,
			IsSearchable: schema.GlobalSearch,
			IsSortable:   schema.IsSortable,
			IsFilterable: schema.IsFilterable,
			IsRelation:   mg.typeMapper.IsRelationshipType(schema.Type),
			RelationType: mg.typeMapper.GetRelationshipType(schema.Type),
		}

		// Go type
		goType := mg.typeMapper.MapFieldTypeToGo(schema.Type, schema.IsNullable)
		if goType.Type != nil {
			info.GoType = goType.Type.String()
			if goType.IsPointer {
				info.GoType = "*" + info.GoType
			}
		}

		// GORM tag
		info.GormTag = mg.buildGormTag(schema)

		infos = append(infos, info)
	}

	return infos
}

/// # buildGormTag
///
/// Bu fonksiyon, field schema'sından GORM tag string'i oluşturur. GORM'un
/// struct tag formatına uygun olarak field özelliklerini tag'e dönüştürür.
///
/// ## Parametreler
///
/// - `schema`: Field schema'sı (*fields.Schema)
///
/// ## Döndürür
///
/// - `string`: GORM tag string'i (örn: "column:email;type:varchar;not null;index")
///
/// ## Oluşturulan Tag Bileşenleri
///
/// ### Temel Bileşenler
/// 1. **column**: Sütun adı (her zaman eklenir)
/// 2. **type**: SQL veri tipi (dialect-aware)
/// 3. **not null**: Field required ise eklenir
/// 4. **index**: Field searchable ise eklenir
///
/// ### GormConfig'den Bileşenler
/// - **primaryKey**: Primary key ise
/// - **unique**: Unique constraint ise
/// - **default**: Default değer varsa
/// - **size**: Boyut belirtilmişse
/// - **precision**: Precision belirtilmişse
/// - **autoIncrement**: Auto increment ise
/// - Ve diğer GORM config özellikleri
///
/// ## Tag Format
///
/// GORM tag formatı: `gorm:"key1:value1;key2:value2;flag1;flag2"`
///
/// ## Kullanım Örneği
///
/// ```go
/// // Email field için tag oluştur
/// emailSchema := &fields.Schema{
///     Key: "email",
///     Type: "string",
///     IsRequired: true,
///     GlobalSearch: true,
/// }
/// tag := mg.buildGormTag(emailSchema)
/// // Çıktı: "column:email;type:varchar(255);not null;index"
///
/// // ID field için tag oluştur
/// idSchema := &fields.Schema{
///     Key: "id",
///     Type: "id",
///     GormConfig: &fields.GormConfig{
///         PrimaryKey: true,
///         AutoIncrement: true,
///     },
/// }
/// tag := mg.buildGormTag(idSchema)
/// // Çıktı: "column:id;primaryKey;autoIncrement;type:bigint"
/// ```
///
/// ## Dialect-Aware SQL Types
///
/// - **PostgreSQL**: VARCHAR, BIGINT, TIMESTAMP, JSONB
/// - **MySQL**: VARCHAR, BIGINT, DATETIME, JSON
/// - **SQLite**: TEXT, INTEGER, DATETIME, TEXT
///
/// ## Önemli Notlar
///
/// - Tag bileşenleri noktalı virgül (;) ile ayrılır
/// - Sıralama önemlidir (column her zaman ilk)
/// - GormConfig varsa önceliklidir
/// - SQL type dialect'e göre değişir
func (mg *MigrationGenerator) buildGormTag(schema *fields.Schema) string {
	var parts []string

	// Column name
	parts = append(parts, "column:"+schema.Key)

	// GormConfig'den tag
	if schema.HasGormConfig() {
		config := schema.GetGormConfig()
		if tag := config.ToGormTag(); tag != "" {
			parts = append(parts, tag)
		}
	}

	// SQL type
	sqlType := mg.typeMapper.MapFieldTypeToSQL(schema.Type, 0)
	parts = append(parts, "type:"+sqlType)

	// Not null
	if schema.IsRequired {
		parts = append(parts, "not null")
	}

	// Index for searchable
	if schema.GlobalSearch {
		parts = append(parts, "index")
	}

	return strings.Join(parts, ";")
}

/// # buildRelationshipFieldInfo
///
/// Bu fonksiyon, ilişkisel field'dan FieldInfo oluşturur. Her ilişki tipine
/// özel olarak gerekli bilgileri çıkarır ve FieldInfo struct'ına dönüştürür.
///
/// ## Parametreler
///
/// - `relField`: İlişkisel field (fields.RelationshipField interface'i)
///
/// ## Döndürür
///
/// - `FieldInfo`: İlişki bilgilerini içeren FieldInfo struct'ı
///
/// ## Desteklenen İlişki Tipleri
///
/// ### BelongsTo (N:1)
/// - Foreign key field'ı gerektirir
/// - Go type: Pointer to related struct (*User)
/// - GORM tag: foreignKey, references
/// - Örnek: Post belongs to User
///
/// ### HasOne (1:1)
/// - Foreign key related model'de
/// - Go type: Pointer to related struct (*Profile)
/// - GORM tag: foreignKey, references
/// - Örnek: User has one Profile
///
/// ### HasMany (1:N)
/// - Foreign key related model'de
/// - Go type: Slice of related struct ([]Post)
/// - GORM tag: foreignKey, references
/// - Örnek: User has many Posts
///
/// ### BelongsToMany (N:M)
/// - Pivot tablo gerektirir
/// - Go type: Slice of pointers to related struct ([]*Role)
/// - GORM tag: many2many, joinForeignKey, joinReferences
/// - Örnek: User belongs to many Roles
///
/// ## Kullanım Örneği
///
/// ```go
/// // BelongsTo field
/// authorField := fields.NewBelongsToField("author", "User").
///     SetForeignKey("author_id")
/// info := mg.buildRelationshipFieldInfo(authorField)
/// // info.GoType = "*User"
/// // info.ForeignKey = "author_id"
/// // info.RelationType = "belongsTo"
///
/// // HasMany field
/// postsField := fields.NewHasManyField("posts", "Post").
///     SetForeignKey("user_id")
/// info := mg.buildRelationshipFieldInfo(postsField)
/// // info.GoType = "[]Post"
/// // info.ForeignKey = "user_id"
/// // info.RelationType = "hasMany"
///
/// // BelongsToMany field
/// rolesField := fields.NewBelongsToManyField("roles", "Role").
///     SetPivotTable("user_roles")
/// info := mg.buildRelationshipFieldInfo(rolesField)
/// // info.GoType = "[]*Role"
/// // info.PivotTable = "user_roles"
/// // info.RelationType = "belongsToMany"
/// ```
///
/// ## Go Type Türetme Kuralları
///
/// 1. Related resource slug'ından struct adı türet
/// 2. Plural ise singular yap (users → User)
/// 3. CamelCase'e çevir (blog_posts → BlogPost)
/// 4. İlişki tipine göre pointer/slice ekle
///
/// ## GORM Tag Oluşturma
///
/// Her ilişki tipi için GormRelationConfig'den tag oluşturulur:
/// - BelongsTo: `foreignKey:AuthorID;references:ID`
/// - HasOne: `foreignKey:UserID;references:ID`
/// - HasMany: `foreignKey:UserID;references:ID`
/// - BelongsToMany: `many2many:user_roles;joinForeignKey:UserID;joinReferences:RoleID`
///
/// ## Önemli Notlar
///
/// - Related resource adı plural olabilir, singular'a çevrilir
/// - Foreign key adlandırması: {model}_id formatında
/// - Pivot tablo adlandırması: {model1}_{model2} formatında
/// - Go type'lar GORM convention'ını takip eder
func (mg *MigrationGenerator) buildRelationshipFieldInfo(relField fields.RelationshipField) FieldInfo {
	info := FieldInfo{
		Name:            relField.GetRelationshipName(),
		Key:             relField.GetKey(),
		IsRelation:      true,
		RelationType:    relField.GetRelationshipType(),
		RelatedResource: relField.GetRelatedResource(),
	}

	// İlişki tipine göre bilgileri ayarla
	switch relField.GetRelationshipType() {
	case "belongsTo":
		// BelongsTo için foreign key field'ı gerekir
		if bt, ok := relField.(*fields.BelongsToField); ok {
			if bt.GormRelationConfig != nil {
				info.ForeignKey = bt.GormRelationConfig.ForeignKey
				info.RelationGormTag = bt.GormRelationConfig.ToGormTag()
				// Go type: pointer to related struct
				relatedType := strcase.ToCamel(info.RelatedResource)
				if strings.HasSuffix(relatedType, "s") {
					relatedType = strings.TrimSuffix(relatedType, "s")
				}
				info.GoType = "*" + relatedType
			}
		}
	case "belongsToMany":
		// BelongsToMany için pivot tablo gerekir
		if btm, ok := relField.(*fields.BelongsToManyField); ok {
			info.PivotTable = btm.PivotTableName
			if btm.GormRelationConfig != nil {
				info.RelationGormTag = btm.GormRelationConfig.ToGormTag()
			}
			// Go type: slice of pointers to related struct
			relatedType := strcase.ToCamel(info.RelatedResource)
			if strings.HasSuffix(relatedType, "s") {
				relatedType = strings.TrimSuffix(relatedType, "s")
			}
			info.GoType = "[]*" + relatedType
		}
	case "hasOne":
		// HasOne için GORM tag gerekir
		if ho, ok := relField.(*fields.HasOneField); ok {
			if ho.GormRelationConfig != nil {
				info.ForeignKey = ho.GormRelationConfig.ForeignKey
				info.RelationGormTag = ho.GormRelationConfig.ToGormTag()
			}
			// Go type: pointer to related struct
			relatedType := strcase.ToCamel(info.RelatedResource)
			if strings.HasSuffix(relatedType, "s") {
				relatedType = strings.TrimSuffix(relatedType, "s")
			}
			info.GoType = "*" + relatedType
		}
	case "hasMany":
		// HasMany için GORM tag gerekir
		if hm, ok := relField.(*fields.HasManyField); ok {
			if hm.GormRelationConfig != nil {
				info.ForeignKey = hm.GormRelationConfig.ForeignKey
				info.RelationGormTag = hm.GormRelationConfig.ToGormTag()
			}
			// Go type: slice of related struct
			relatedType := strcase.ToCamel(info.RelatedResource)
			if strings.HasSuffix(relatedType, "s") {
				relatedType = strings.TrimSuffix(relatedType, "s")
			}
			info.GoType = "[]" + relatedType
		}
	}

	return info
}

/// # GenerateModelStub
///
/// Bu fonksiyon, resource tanımından Go model struct stub'ı oluşturur. Manuel
/// model oluşturmak için referans olarak kullanılabilir. İlişkisel field'ları
/// da otomatik olarak ekler ve GORM tag'lerini oluşturur.
///
/// ## Parametreler
///
/// - `r`: Model stub'ı oluşturulacak resource
///
/// ## Döndürür
///
/// - `string`: Go struct tanımı (string formatında)
///
/// ## Oluşturulan Struct Yapısı
///
/// ```go
/// type ModelName struct {
///     ID        uint      `json:"id" gorm:"primaryKey"`
///     // Field'lardan oluşturulan alanlar
///     Email     string    `json:"email" gorm:"column:email;type:varchar;not null;index"`
///     // İlişkisel field'lar
///     AuthorID  uint      `json:"author_id" gorm:"index"`
///     Author    *User     `json:"author" gorm:"foreignKey:AuthorID"`
///     // Timestamp alanları
///     CreatedAt time.Time `json:"createdAt" gorm:"index"`
///     UpdatedAt time.Time `json:"updatedAt" gorm:"index"`
/// }
/// ```
///
/// ## Otomatik Eklenen Alanlar
///
/// 1. **ID**: Primary key (uint, auto increment)
/// 2. **Field'lar**: Resource field tanımlarından
/// 3. **Foreign Key'ler**: BelongsTo ilişkileri için
/// 4. **İlişki Field'ları**: Tüm ilişki tipleri için
/// 5. **CreatedAt**: Oluşturma zamanı
/// 6. **UpdatedAt**: Güncelleme zamanı
///
/// ## İlişki Field'ları
///
/// ### BelongsTo
/// ```go
/// AuthorID  uint  `json:"author_id" gorm:"index"`
/// Author    *User `json:"author" gorm:"foreignKey:AuthorID"`
/// ```
///
/// ### HasOne
/// ```go
/// Profile *Profile `json:"profile" gorm:"foreignKey:UserID"`
/// ```
///
/// ### HasMany
/// ```go
/// Posts []Post `json:"posts" gorm:"foreignKey:UserID"`
/// ```
///
/// ### BelongsToMany
/// ```go
/// Roles []*Role `json:"roles" gorm:"many2many:user_roles"`
/// ```
///
/// ## Kullanım Senaryoları
///
/// 1. **Hızlı Prototipleme**: Resource'dan model oluştur
/// 2. **Referans**: Manuel model yazarken referans al
/// 3. **Dokümantasyon**: Model yapısını dokümante et
/// 4. **Code Generation**: Otomatik kod üretimi
///
/// ## Kullanım Örneği
///
/// ```go
/// // User resource tanımla
/// userResource := resource.New("users").
///     AddField(fields.NewTextField("name").SetRequired(true)).
///     AddField(fields.NewTextField("email").SetRequired(true).SetSearchable(true)).
///     AddField(fields.NewHasManyField("posts", "Post"))
///
/// // Model stub oluştur
/// stub := mg.GenerateModelStub(userResource)
/// fmt.Println(stub)
///
/// // Çıktı:
/// // type User struct {
/// //     ID        uint      `json:"id" gorm:"primaryKey"`
/// //     Name      string    `json:"name" gorm:"column:name;type:varchar;not null"`
/// //     Email     string    `json:"email" gorm:"column:email;type:varchar;not null;index"`
/// //     Posts     []Post    `json:"posts" gorm:"foreignKey:UserID"`
/// //     CreatedAt time.Time `json:"createdAt" gorm:"index"`
/// //     UpdatedAt time.Time `json:"updatedAt" gorm:"index"`
/// // }
///
/// // Dosyaya yaz
/// err := os.WriteFile("models/user.go", []byte(stub), 0644)
/// ```
///
/// ## Struct Adlandırma Kuralları
///
/// 1. Resource slug'ından türet (users → User)
/// 2. Plural ise singular yap
/// 3. CamelCase'e çevir (blog_posts → BlogPost)
/// 4. İlk harf büyük (exported)
///
/// ## Field Adlandırma Kuralları
///
/// 1. Field key'den türet (email → Email)
/// 2. CamelCase'e çevir (first_name → FirstName)
/// 3. İlk harf büyük (exported)
///
/// ## JSON ve GORM Tag'leri
///
/// - **JSON tag**: Field key (snake_case)
/// - **GORM tag**: Column, type, constraints
///
/// ## Avantajlar
///
/// - Hızlı model oluşturma
/// - İlişkileri otomatik ekler
/// - GORM tag'lerini otomatik oluşturur
/// - Dialect-aware SQL type'lar
/// - Constraint'leri otomatik ekler
///
/// ## Dezavantajlar
///
/// - Basit stub oluşturur (custom logic yok)
/// - Method'lar eklenmez
/// - Validation logic eklenmez
/// - Custom type'lar desteklenmez
///
/// ## Önemli Notlar
///
/// ⚠️ Bu fonksiyon sadece stub oluşturur, production-ready model değildir.
/// ⚠️ Oluşturulan model'i ihtiyaçlarınıza göre özelleştirin.
/// ⚠️ Custom business logic ekleyin.
/// ⚠️ Validation method'ları ekleyin.
///
/// ## Best Practices
///
/// ```go
/// // 1. Stub oluştur
/// stub := mg.GenerateModelStub(userResource)
///
/// // 2. Dosyaya yaz
/// err := os.WriteFile("models/user.go", []byte(stub), 0644)
///
/// // 3. Manuel olarak özelleştir:
/// // - Custom method'lar ekle
/// // - Validation logic ekle
/// // - Business logic ekle
/// // - Hook'lar ekle (BeforeCreate, AfterUpdate, vb.)
/// ```
func (mg *MigrationGenerator) GenerateModelStub(r resource.Resource) string {
	var sb strings.Builder

	structName := strcase.ToCamel(r.Slug())
	// Tekil form için son 's' karakterini kaldır (basit çoğul)
	if strings.HasSuffix(structName, "s") {
		structName = strings.TrimSuffix(structName, "s")
	}

	sb.WriteString(fmt.Sprintf("type %s struct {\n", structName))

	// ID alanı
	sb.WriteString("\tID        uint      `json:\"id\" gorm:\"primaryKey\"`\n")

	// Field'lardan alanlar
	for _, info := range mg.GetFieldInfos(r) {
		if info.Key == "id" {
			continue // ID zaten eklendi
		}

		// İlişkisel field'lar için özel işlem
		if info.IsRelation {
			// BelongsTo için foreign key field'ı ekle
			if info.RelationType == "belongsTo" && info.ForeignKey != "" {
				fkFieldName := strcase.ToCamel(info.ForeignKey)
				// Foreign key için basit GORM tag
				fkGormTag := "index"
				sb.WriteString(fmt.Sprintf("\t%s uint `json:\"%s\" gorm:\"%s\"`\n",
					fkFieldName, info.ForeignKey, fkGormTag))
			}

			// İlişki field'ı ekle
			relationFieldName := strcase.ToCamel(info.Key)
			goType := info.GoType
			if goType == "" {
				// Fallback: related resource'dan tip oluştur
				relatedType := strcase.ToCamel(info.RelatedResource)
				if strings.HasSuffix(relatedType, "s") {
					relatedType = strings.TrimSuffix(relatedType, "s")
				}
				switch info.RelationType {
				case "belongsTo", "hasOne":
					goType = "*" + relatedType
				case "hasMany":
					goType = "[]" + relatedType
				case "belongsToMany":
					goType = "[]*" + relatedType
				}
			}

			gormTag := info.RelationGormTag
			jsonTag := info.Key

			sb.WriteString(fmt.Sprintf("\t%s %s `json:\"%s\" gorm:\"%s\"`\n",
				relationFieldName, goType, jsonTag, gormTag))
			continue
		}

		// Normal field'lar
		fieldName := strcase.ToCamel(info.Key)
		goType := info.GoType
		if goType == "" {
			goType = "string"
		}

		gormTag := info.GormTag
		jsonTag := info.Key

		sb.WriteString(fmt.Sprintf("\t%s %s `json:\"%s\" gorm:\"%s\"`\n",
			fieldName, goType, jsonTag, gormTag))
	}

	// Timestamp alanları
	sb.WriteString("\tCreatedAt time.Time `json:\"createdAt\" gorm:\"index\"`\n")
	sb.WriteString("\tUpdatedAt time.Time `json:\"updatedAt\" gorm:\"index\"`\n")

	sb.WriteString("}\n")

	return sb.String()
}

// createTableFromFields, resource'un field tanımlarından tablo oluşturur.
// Model olmayan resource'lar için kullanılır.
