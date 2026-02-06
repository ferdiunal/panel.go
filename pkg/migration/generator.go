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
func (mg *MigrationGenerator) AutoMigrate() error {
	for _, r := range mg.resources {
		model := r.Model()
		if model == nil {
			continue
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
	tableName := mg.getTableName(r)

	for _, field := range r.Fields() {
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
	// GORM naming convention: slug'dan tablo adı türet
	slug := r.Slug()
	return strcase.ToSnake(slug)
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
}

// GetFieldInfos, resource'un tüm alanlarının bilgilerini döner.
func (mg *MigrationGenerator) GetFieldInfos(r resource.Resource) []FieldInfo {
	var infos []FieldInfo

	for _, field := range r.Fields() {
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

// GenerateModelStub, resource'dan Go model stub'ı oluşturur.
// Bu stub, manuel model oluşturmak için referans olarak kullanılabilir.
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
