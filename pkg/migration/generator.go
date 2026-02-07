package migration

import (
	"fmt"
	"strings"

	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/resource"
	"github.com/iancoleman/strcase"
	"gorm.io/gorm"
)

// MigrationGenerator, Resource tanımlarından veritabanı migration işlemlerini yönetir.
type MigrationGenerator struct {
	db         *gorm.DB
	resources  []resource.Resource
	typeMapper *TypeMapper
	dialect    string
}

// NewMigrationGenerator, yeni bir MigrationGenerator oluşturur.
func NewMigrationGenerator(db *gorm.DB) *MigrationGenerator {
	dialect := db.Dialector.Name()
	return &MigrationGenerator{
		db:         db,
		resources:  []resource.Resource{},
		typeMapper: NewTypeMapperWithDialect(dialect),
		dialect:    dialect,
	}
}

// RegisterResource, migration için resource kaydeder.
func (mg *MigrationGenerator) RegisterResource(r resource.Resource) *MigrationGenerator {
	mg.resources = append(mg.resources, r)
	return mg
}

// RegisterResources, birden fazla resource'u kaydeder.
func (mg *MigrationGenerator) RegisterResources(resources ...resource.Resource) *MigrationGenerator {
	mg.resources = append(mg.resources, resources...)
	return mg
}

// AutoMigrate, kayıtlı tüm resource'ların modellerini migrate eder.
// Model olmayan resource'lar için hata döner.
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

// applyFieldConstraints, field tanımlarından ek constraint'ler oluşturur.
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

// getTableName, resource'dan tablo adını çıkarır.
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

// hasIndex, tabloda index var mı kontrol eder (GORM Migrator kullanarak).
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

// hasIndexWithModel, tabloda index var mı kontrol eder (GORM Migrator kullanarak - model instance ile).
func (mg *MigrationGenerator) hasIndexWithModel(model interface{}, column string) bool {
	tableName := mg.getTableNameFromModel(model)
	indexName := fmt.Sprintf("idx_%s_%s", tableName, column)

	// GORM Migrator kullanarak - dialect-aware ve güvenli
	return mg.db.Migrator().HasIndex(model, indexName)
}

// hasUniqueIndexWithModel, tabloda unique index var mı kontrol eder (GORM Migrator kullanarak - model instance ile).
func (mg *MigrationGenerator) hasUniqueIndexWithModel(model interface{}, column string) bool {
	tableName := mg.getTableNameFromModel(model)
	indexName := fmt.Sprintf("uniq_%s_%s", tableName, column)

	// GORM Migrator kullanarak - dialect-aware ve güvenli
	return mg.db.Migrator().HasIndex(model, indexName)
}

// createIndexWithModel, index oluşturur (GORM Migrator kullanarak - model instance ile).
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

// getTableNameFromModel, model instance'ından tablo adını çıkarır.
func (mg *MigrationGenerator) getTableNameFromModel(model interface{}) string {
	stmt := &gorm.Statement{DB: mg.db}
	err := stmt.Parse(model)
	if err != nil {
		return ""
	}
	return stmt.Table
}

// hasUniqueIndex, tabloda unique index var mı kontrol eder.
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

// createIndex, index oluşturur.
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

// createPivotTable, BelongsToMany ilişkileri için pivot tablo oluşturur.
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

// FieldInfo, alan bilgilerini içerir.
type FieldInfo struct {
	Name         string
	Key          string
	GoType       string
	SQLType      string
	GormTag      string
	IsRequired   bool
	IsNullable   bool
	IsSearchable bool
	IsSortable   bool
	IsFilterable bool
	IsRelation   bool
	RelationType string

	// İlişki Bilgileri
	RelatedResource  string // İlişkili resource slug'ı
	ForeignKey       string // Foreign key sütunu
	PivotTable       string // Pivot tablo adı (BelongsToMany için)
	RelationGormTag  string // İlişki için GORM tag'i
}

// GetFieldInfos, resource'un tüm alanlarının bilgilerini döner.
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

// buildGormTag, schema'dan GORM tag oluşturur.
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

// buildRelationshipFieldInfo, ilişkisel field'dan FieldInfo oluşturur.
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

// GenerateModelStub, resource'dan Go model stub'ı oluşturur.
// Bu stub, manuel model oluşturmak için referans olarak kullanılabilir.
// İlişkisel field'ları da otomatik olarak ekler.
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
