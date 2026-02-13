package fields

// MorphToMany, polimorfik many-to-many ilişkiyi temsil eder (örn. Tag -> Taggable: posts, videos).
//
// MorphToMany ilişkisi, bir kaydın farklı tiplerdeki birden fazla kayda ait olabileceğini belirtir.
// Bu, veritabanında polimorfik pivot tablo ile temsil edilir.
//
// # Kullanım Senaryoları
//
// - **Tag -> Taggable**: Bir etiket hem Post'lara hem de Video'lara ait olabilir
// - **Image -> Imageable**: Bir resim hem User'lara hem de Product'lara ait olabilir
//
// # Veritabanı Yapısı
//
// Polimorfik pivot tablo genellikle şu yapıya sahiptir:
//
//	CREATE TABLE taggables (
//	    tag_id INT,              -- Ana tablo FK
//	    taggable_id INT,         -- Polimorfik ID
//	    taggable_type VARCHAR,   -- Polimorfik tip ("post", "video")
//	    PRIMARY KEY (tag_id, taggable_id, taggable_type)
//	);
//
// # Kullanım Örneği
//
//	field := fields.NewMorphToMany("Taggable", "taggable").
//	    Types(map[string]string{
//	        "post":  "posts",
//	        "video": "videos",
//	    }).
//	    Displays(map[string]string{
//	        "post":  "title",
//	        "video": "name",
//	    }).
//	    PivotTable("taggables").
//	    AutoOptions("name").
//	    WithEagerLoad()
type MorphToMany struct {
	Schema
	TypeMappings       map[string]string // Type => Resource slug mapping
	DisplayMappings    map[string]string // Type => Display field name
	PivotTableName     string            // Pivot tablo adı
	ForeignKeyColumn   string            // Ana tablonun FK (örn. "tag_id")
	RelatedKeyColumn   string            // Polimorfik ID sütunu (örn. "taggable_id")
	MorphTypeColumn    string            // Polimorfik tip sütunu (örn. "taggable_type")
	QueryCallback      func(query interface{}) interface{}
	LoadingStrategy    LoadingStrategy
	GormRelationConfig *RelationshipGormConfig
}

// NewMorphToMany, yeni bir MorphToMany polimorfik many-to-many ilişki alanı oluşturur.
//
// # Parametreler
//
// - **name**: Alanın görünen adı (örn. "Taggable", "Imageable")
// - **key**: İlişki key'i (örn. "taggable", "imageable")
//
// # Varsayılan Değerler
//
// - **PivotTableName**: key + "s" (örn. "taggables")
// - **ForeignKeyColumn**: "tag_id" (ana tablonun foreign key'i)
// - **RelatedKeyColumn**: key + "_id" (örn. "taggable_id")
// - **MorphTypeColumn**: key + "_type" (örn. "taggable_type")
// - **LoadingStrategy**: EAGER_LOADING
//
// Döndürür:
//   - Yapılandırılmış MorphToMany pointer'ı
func NewMorphToMany(name, key string) *MorphToMany {
	pivotTable := key + "s"

	m := &MorphToMany{
		Schema: Schema{
			Name:      name,
			LabelText: name,
			Key:       key,
			View:      "morph-to-many-field",
			Type:      TYPE_RELATIONSHIP,
			Props: map[string]interface{}{
				"types":    []map[string]string{},
				"displays": map[string]string{},
			},
		},
		TypeMappings:     make(map[string]string),
		DisplayMappings:  make(map[string]string),
		PivotTableName:   pivotTable,
		ForeignKeyColumn: "tag_id",
		RelatedKeyColumn: key + "_id",
		MorphTypeColumn:  key + "_type",
		LoadingStrategy:  EAGER_LOADING,
		GormRelationConfig: NewRelationshipGormConfig().
			WithPivotTable(pivotTable, "tag_id", key+"_id").
			WithPolymorphic(key+"_type", key+"_id"),
	}
	return m
}

// Types, polimorfik ilişki için tip eşlemelerini ayarlar.
//
// # Parametreler
//
// - **types**: Tip değeri -> resource slug eşlemesi (örn. {"post": "posts", "video": "videos"})
//
// # Kullanım Örneği
//
//	field := fields.NewMorphToMany("Taggable", "taggable").
//	    Types(map[string]string{
//	        "post":  "posts",
//	        "video": "videos",
//	    })
//
// Döndürür:
//   - MorphToMany pointer'ı (method chaining için)
func (m *MorphToMany) Types(types map[string]string) *MorphToMany {
	m.TypeMappings = types
	m.Props["types"] = m.formatTypesForFrontend(types)
	return m
}

// Displays, her tip için görüntüleme alanını ayarlar.
//
// # Parametreler
//
// - **displays**: Tip değeri -> görüntüleme alanı eşlemesi (örn. {"post": "title", "video": "name"})
//
// # Kullanım Örneği
//
//	field := fields.NewMorphToMany("Taggable", "taggable").
//	    Displays(map[string]string{
//	        "post":  "title",
//	        "video": "name",
//	    })
//
// Döndürür:
//   - MorphToMany pointer'ı (method chaining için)
func (m *MorphToMany) Displays(displays map[string]string) *MorphToMany {
	m.DisplayMappings = displays
	m.Props["displays"] = displays
	return m
}

// formatTypesForFrontend, tip eşlemelerini frontend select seçeneklerine dönüştürür.
func (m *MorphToMany) formatTypesForFrontend(types map[string]string) []map[string]string {
	var options []map[string]string
	for dbType, resourceSlug := range types {
		label := resourceSlug
		if len(resourceSlug) > 0 {
			label = string(resourceSlug[0]-32) + resourceSlug[1:]
		}

		options = append(options, map[string]string{
			"label": label,
			"value": dbType,
			"slug":  resourceSlug,
		})
	}
	return options
}

// PivotTable, pivot tablo adını özelleştirir.
//
// # Parametreler
//
// - **tableName**: Pivot tablo adı (örn. "taggables", "imageables")
//
// Döndürür:
//   - MorphToMany pointer'ı (method chaining için)
func (m *MorphToMany) PivotTable(tableName string) *MorphToMany {
	m.PivotTableName = tableName
	if m.GormRelationConfig != nil {
		m.GormRelationConfig.PivotTable = tableName
	}
	return m
}

// ForeignKey, ana tablonun foreign key sütun adını ayarlar.
//
// # Parametreler
//
// - **key**: Foreign key sütun adı (örn. "tag_id")
//
// Döndürür:
//   - MorphToMany pointer'ı (method chaining için)
func (m *MorphToMany) ForeignKey(key string) *MorphToMany {
	m.ForeignKeyColumn = key
	if m.GormRelationConfig != nil {
		m.GormRelationConfig.JoinForeignKey = key
	}
	return m
}

// RelatedKey, polimorfik ID sütun adını ayarlar.
//
// # Parametreler
//
// - **key**: Polimorfik ID sütun adı (örn. "taggable_id")
//
// Döndürür:
//   - MorphToMany pointer'ı (method chaining için)
func (m *MorphToMany) RelatedKey(key string) *MorphToMany {
	m.RelatedKeyColumn = key
	if m.GormRelationConfig != nil {
		m.GormRelationConfig.JoinReferences = key
	}
	return m
}

// MorphType, polimorfik tip sütun adını ayarlar.
//
// # Parametreler
//
// - **column**: Polimorfik tip sütun adı (örn. "taggable_type")
//
// Döndürür:
//   - MorphToMany pointer'ı (method chaining için)
func (m *MorphToMany) MorphType(column string) *MorphToMany {
	m.MorphTypeColumn = column
	if m.GormRelationConfig != nil {
		m.GormRelationConfig.PolymorphicType = column
	}
	return m
}

// AutoOptions, ilişkili tablodan otomatik options oluşturmayı etkinleştirir.
//
// Bu metod, MorphToMany ilişkisinde ilişkili kayıtların otomatik olarak yüklenmesini
// ve frontend'de multi-select combobox'ta gösterilmesini sağlar.
//
// # Parametreler
//
// - **displayField**: Option label'ı için kullanılacak sütun adı (örn. "name", "title")
//
// # Kullanım Örneği
//
//	field := fields.NewMorphToMany("Taggable", "taggable").
//	    Types(map[string]string{
//	        "post":  "posts",
//	        "video": "videos",
//	    }).
//	    AutoOptions("name")  // Ana tablodaki "name" sütunu label olarak kullanılır
//
// Döndürür:
//   - MorphToMany pointer'ı (method chaining için)
func (m *MorphToMany) AutoOptions(displayField string) *MorphToMany {
	m.Schema.AutoOptions(displayField)
	return m
}

// Query, ilişki sorgusunu özelleştirmek için callback ayarlar.
func (m *MorphToMany) Query(fn func(interface{}) interface{}) *MorphToMany {
	m.QueryCallback = fn
	return m
}

// WithEagerLoad, yükleme stratejisini eager loading olarak ayarlar.
func (m *MorphToMany) WithEagerLoad() *MorphToMany {
	m.LoadingStrategy = EAGER_LOADING
	return m
}

// WithLazyLoad, yükleme stratejisini lazy loading olarak ayarlar.
func (m *MorphToMany) WithLazyLoad() *MorphToMany {
	m.LoadingStrategy = LAZY_LOADING
	return m
}

// GetRelationshipType, ilişki türünü döndürür.
func (m *MorphToMany) GetRelationshipType() string {
	return "morphToMany"
}

// GetRelatedResource, ilgili resource slug'ını döndürür.
// MorphToMany için uygulanamaz çünkü birden fazla resource'a ait olabilir.
func (m *MorphToMany) GetRelatedResource() string {
	return ""
}

// GetRelationshipName, ilişkinin adını döndürür.
func (m *MorphToMany) GetRelationshipName() string {
	return m.Name
}

// ResolveRelationship, polimorfik many-to-many ilişkiyi çözümler.
func (m *MorphToMany) ResolveRelationship(item interface{}) (interface{}, error) {
	if item == nil {
		return []interface{}{}, nil
	}
	// Gerçek implementasyonda veritabanından ilişkili kayıtlar yüklenir
	return []interface{}{}, nil
}

// ValidateRelationship, ilişkiyi doğrular.
func (m *MorphToMany) ValidateRelationship(value interface{}) error {
	return nil
}

// GetDisplayKey, görüntüleme key'ini döndürür.
func (m *MorphToMany) GetDisplayKey() string {
	return ""
}

// GetSearchableColumns, aranabilir sütunları döndürür.
func (m *MorphToMany) GetSearchableColumns() []string {
	return []string{}
}

// GetQueryCallback, sorgu callback'ini döndürür.
func (m *MorphToMany) GetQueryCallback() func(interface{}) interface{} {
	if m.QueryCallback == nil {
		return func(q interface{}) interface{} { return q }
	}
	return m.QueryCallback
}

// GetLoadingStrategy, yükleme stratejisini döndürür.
func (m *MorphToMany) GetLoadingStrategy() LoadingStrategy {
	if m.LoadingStrategy == "" {
		return EAGER_LOADING
	}
	return m.LoadingStrategy
}

// Searchable, alanı aranabilir olarak işaretler.
func (m *MorphToMany) Searchable() Element {
	m.GlobalSearch = true
	return m
}

// IsRequired, alanın zorunlu olup olmadığını döndürür.
func (m *MorphToMany) IsRequired() bool {
	return m.Schema.IsRequired
}

// GetTypes, tip eşlemelerini döndürür.
func (m *MorphToMany) GetTypes() map[string]string {
	if m.TypeMappings == nil {
		return make(map[string]string)
	}
	return m.TypeMappings
}
