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
}

// NewMigrationGenerator, yeni bir MigrationGenerator oluşturur.
func NewMigrationGenerator(db *gorm.DB) *MigrationGenerator {
	return &MigrationGenerator{
		db:         db,
		resources:  []resource.Resource{},
		typeMapper: NewTypeMapper(),
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
// Model yoksa field-based migration kullanır.
func (mg *MigrationGenerator) AutoMigrate() error {
	for _, r := range mg.resources {
		model := r.Model()
		if model == nil {
			// Model yoksa field-based migration kullan
			if err := mg.createTableFromFields(r); err != nil {
				return fmt.Errorf("field-based migration failed for %s: %w", r.Slug(), err)
			}
		} else {
			// GORM AutoMigrate
			if err := mg.db.AutoMigrate(model); err != nil {
				return fmt.Errorf("migration failed for %s: %w", r.Slug(), err)
			}
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
	tableName := mg.getTableName(r)

	for _, field := range r.Fields() {
		// İlişkisel field'ları kontrol et
		if relField, ok := fields.IsRelationshipField(field); ok {
			// BelongsTo için foreign key index'i
			if relField.GetRelationshipType() == "belongsTo" {
				if bt, ok := relField.(*fields.BelongsTo); ok {
					if bt.GormRelationConfig != nil && bt.GormRelationConfig.ForeignKey != "" {
						fkColumn := bt.GormRelationConfig.ForeignKey
						if !mg.hasIndex(tableName, fkColumn) {
							if err := mg.createIndex(tableName, fkColumn, false); err != nil {
								return err
							}
						}
					}
				}
			}

			// BelongsToMany için pivot tablo
			if relField.GetRelationshipType() == "belongsToMany" {
				if btm, ok := relField.(*fields.BelongsToMany); ok {
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
		if schema.GlobalSearch && !mg.hasIndex(tableName, schema.Key) {
			if err := mg.createIndex(tableName, schema.Key, false); err != nil {
				return err
			}
		}

		// Sortable alanlar için index
		if schema.IsSortable && !mg.hasIndex(tableName, schema.Key) {
			if err := mg.createIndex(tableName, schema.Key, false); err != nil {
				return err
			}
		}

		// Filterable alanlar için index
		if schema.IsFilterable && !mg.hasIndex(tableName, schema.Key) {
			if err := mg.createIndex(tableName, schema.Key, false); err != nil {
				return err
			}
		}

		// GormConfig'den constraint'ler
		if schema.HasGormConfig() {
			config := schema.GetGormConfig()

			// Unique Index
			if config.UniqueIndex && !mg.hasUniqueIndex(tableName, schema.Key) {
				if err := mg.createIndex(tableName, schema.Key, true); err != nil {
					return err
				}
			}

			// Normal Index
			if config.Index && !mg.hasIndex(tableName, schema.Key) {
				if err := mg.createIndex(tableName, schema.Key, false); err != nil {
					return err
				}
			}
		}

		// Validation rules'dan unique constraint
		for _, rule := range schema.ValidationRules {
			if rule.Name == "unique" {
				if !mg.hasUniqueIndex(tableName, schema.Key) {
					if err := mg.createIndex(tableName, schema.Key, true); err != nil {
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

// hasIndex, tabloda index var mı kontrol eder.
func (mg *MigrationGenerator) hasIndex(table, column string) bool {
	indexName := fmt.Sprintf("idx_%s_%s", table, column)
	var count int64

	// SQLite için
	mg.db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='index' AND name=?", indexName).Scan(&count)

	return count > 0
}

// hasUniqueIndex, tabloda unique index var mı kontrol eder.
func (mg *MigrationGenerator) hasUniqueIndex(table, column string) bool {
	indexName := fmt.Sprintf("uniq_%s_%s", table, column)
	var count int64

	mg.db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='index' AND name=?", indexName).Scan(&count)

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
func (mg *MigrationGenerator) createPivotTable(btm *fields.BelongsToMany) error {
	// Pivot tablo zaten var mı kontrol et
	var count int64
	mg.db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?", btm.PivotTableName).Scan(&count)
	if count > 0 {
		return nil // Tablo zaten var
	}

	// Pivot tablo oluştur
	sql := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			%s INTEGER NOT NULL,
			%s INTEGER NOT NULL,
			PRIMARY KEY (%s, %s)
		)
	`, btm.PivotTableName, btm.ForeignKeyColumn, btm.RelatedKeyColumn, btm.ForeignKeyColumn, btm.RelatedKeyColumn)

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
		if bt, ok := relField.(*fields.BelongsTo); ok {
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
		if btm, ok := relField.(*fields.BelongsToMany); ok {
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
		if ho, ok := relField.(*fields.HasOne); ok {
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
		if hm, ok := relField.(*fields.HasMany); ok {
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
func (mg *MigrationGenerator) createTableFromFields(r resource.Resource) error {
	tableName := mg.getTableName(r)

	// Tablo zaten var mı kontrol et
	if mg.db.Migrator().HasTable(tableName) {
		// Tablo var, mevcut yapıyı güncelle (ALTER TABLE)
		return mg.updateTableFromFields(r, tableName)
	}

	// Field'lardan SQL column tanımları oluştur
	var columns []string

	// ID column (primary key)
	columns = append(columns, "id INTEGER PRIMARY KEY AUTOINCREMENT")

	// Timestamp field'larını takip et
	hasCreatedAt := false
	hasUpdatedAt := false
	hasDeletedAt := false

	// Field'lardan column'lar
	for _, field := range r.Fields() {
		// İlişkisel field'ları kontrol et
		if relField, ok := fields.IsRelationshipField(field); ok {
			// BelongsTo için foreign key column'u ekle
			if relField.GetRelationshipType() == "belongsTo" {
				if bt, ok := relField.(*fields.BelongsTo); ok {
					if bt.GormRelationConfig != nil && bt.GormRelationConfig.ForeignKey != "" {
						fkColumn := bt.GormRelationConfig.ForeignKey
						// Foreign key column
						columns = append(columns, fmt.Sprintf("%s INTEGER", fkColumn))
					}
				}
			}
			// Diğer ilişki tipleri için column ekleme (HasOne, HasMany, BelongsToMany için column yok)
			continue
		}

		// Normal field'lar
		schema, ok := field.(*fields.Schema)
		if !ok {
			continue
		}

		if schema.Key == "id" {
			continue // ID zaten eklendi
		}

		// Timestamp field'larını takip et
		if schema.Key == "created_at" {
			hasCreatedAt = true
		} else if schema.Key == "updated_at" {
			hasUpdatedAt = true
		} else if schema.Key == "deleted_at" {
			hasDeletedAt = true
		}

		// SQL type
		sqlType := mg.typeMapper.MapFieldTypeToSQL(schema.Type, 0)

		// Column definition
		columnDef := fmt.Sprintf("%s %s", schema.Key, sqlType)

		// NOT NULL
		if schema.IsRequired && !schema.IsNullable {
			columnDef += " NOT NULL"
		}

		// Default value
		if schema.HasGormConfig() {
			config := schema.GetGormConfig()
			if config.Default != "" {
				columnDef += fmt.Sprintf(" DEFAULT %s", config.Default)
			}
		}

		columns = append(columns, columnDef)
	}

	// Timestamp columns (sadece field'larda tanımlanmamışsa ekle)
	if !hasCreatedAt {
		columns = append(columns, "created_at DATETIME")
	}
	if !hasUpdatedAt {
		columns = append(columns, "updated_at DATETIME")
	}
	if !hasDeletedAt {
		columns = append(columns, "deleted_at DATETIME") // Soft delete
	}

	// CREATE TABLE SQL
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n\t%s\n)", tableName, strings.Join(columns, ",\n\t"))

	// Execute
	if err := mg.db.Exec(sql).Error; err != nil {
		return fmt.Errorf("failed to create table %s: %w", tableName, err)
	}

	return nil
}

// updateTableFromFields, mevcut tabloyu field tanımlarına göre günceller.
// Eksik column'ları ekler, mevcut column'ları değiştirmez (güvenli).
func (mg *MigrationGenerator) updateTableFromFields(r resource.Resource, tableName string) error {
	// Mevcut column'ları al
	type ColumnInfo struct {
		Name string
	}
	var existingColumns []ColumnInfo

	// SQLite için column listesi
	mg.db.Raw(fmt.Sprintf("PRAGMA table_info(%s)", tableName)).Scan(&existingColumns)

	existingColumnMap := make(map[string]bool)
	for _, col := range existingColumns {
		existingColumnMap[col.Name] = true
	}

	// Eksik column'ları ekle
	for _, field := range r.Fields() {
		// İlişkisel field'ları kontrol et
		if relField, ok := fields.IsRelationshipField(field); ok {
			// BelongsTo için foreign key column'u
			if relField.GetRelationshipType() == "belongsTo" {
				if bt, ok := relField.(*fields.BelongsTo); ok {
					if bt.GormRelationConfig != nil && bt.GormRelationConfig.ForeignKey != "" {
						fkColumn := bt.GormRelationConfig.ForeignKey
						if !existingColumnMap[fkColumn] {
							// Column ekle
							sql := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s INTEGER", tableName, fkColumn)
							if err := mg.db.Exec(sql).Error; err != nil {
								return fmt.Errorf("failed to add column %s: %w", fkColumn, err)
							}
						}
					}
				}
			}
			continue
		}

		// Normal field'lar
		schema, ok := field.(*fields.Schema)
		if !ok {
			continue
		}

		if schema.Key == "id" {
			continue
		}

		// Column zaten var mı?
		if existingColumnMap[schema.Key] {
			continue
		}

		// SQL type
		sqlType := mg.typeMapper.MapFieldTypeToSQL(schema.Type, 0)

		// Column definition
		columnDef := fmt.Sprintf("%s %s", schema.Key, sqlType)

		// NOT NULL (dikkat: mevcut tabloya NOT NULL column eklerken default value gerekir)
		if schema.IsRequired && !schema.IsNullable {
			// Mevcut tabloya NOT NULL column eklerken default value ekle
			if schema.HasGormConfig() {
				config := schema.GetGormConfig()
				if config.Default != "" {
					columnDef += fmt.Sprintf(" DEFAULT %s NOT NULL", config.Default)
				}
				// Default value yoksa, NULL'a izin ver (NOT NULL ekleme)
			}
		}

		// ALTER TABLE ADD COLUMN
		sql := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s", tableName, columnDef)
		if err := mg.db.Exec(sql).Error; err != nil {
			return fmt.Errorf("failed to add column %s: %w", schema.Key, err)
		}
	}

	return nil
}
