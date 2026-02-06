package fields

// GormConfig, alan için GORM veritabanı yapılandırmasını tanımlar.
// Bu yapı, alan tanımlarından otomatik migration ve model oluşturma için kullanılır.
type GormConfig struct {
	// Anahtar Yapılandırması
	PrimaryKey    bool `json:"-"` // Birincil anahtar mı?
	AutoIncrement bool `json:"-"` // Otomatik artış mı?

	// Sütun Yapılandırması
	Column string `json:"-"` // Özel sütun adı (varsayılan: Key'den türetilir)
	Type   string `json:"-"` // SQL tipi (örn: "varchar(255)", "text", "decimal(10,2)")
	Size   int    `json:"-"` // Sütun boyutu (VARCHAR için)

	// Sayısal Hassasiyet
	Precision int `json:"-"` // Ondalık hassasiyet (DECIMAL tipi için)
	Scale     int `json:"-"` // Ondalık ölçek (DECIMAL tipi için)

	// İndeks Yapılandırması
	Index       bool   `json:"-"` // Normal indeks oluştur
	UniqueIndex bool   `json:"-"` // Benzersiz indeks oluştur
	IndexName   string `json:"-"` // Özel indeks adı

	// Varsayılan Değer
	Default interface{} `json:"-"` // Varsayılan değer

	// Kısıtlamalar
	NotNull bool   `json:"-"` // NOT NULL kısıtlaması
	Comment string `json:"-"` // Sütun yorumu

	// İlişki Yapılandırması (Foreign Key)
	ForeignKey       string `json:"-"` // Foreign key alanı (örn: "AuthorID")
	References       string `json:"-"` // Referans tablo.alan (örn: "users.id")
	OnDelete         string `json:"-"` // ON DELETE davranışı (CASCADE, SET NULL, vb.)
	OnUpdate         string `json:"-"` // ON UPDATE davranışı
	ManyToManyTable  string `json:"-"` // Many-to-many ara tablo adı
	JoinForeignKey   string `json:"-"` // Join tablosundaki foreign key
	JoinReferences   string `json:"-"` // Join tablosundaki referans
}

// NewGormConfig, varsayılan değerlerle yeni bir GormConfig oluşturur.
func NewGormConfig() *GormConfig {
	return &GormConfig{}
}

// WithPrimaryKey, birincil anahtar olarak işaretler.
func (g *GormConfig) WithPrimaryKey() *GormConfig {
	g.PrimaryKey = true
	g.AutoIncrement = true
	return g
}

// WithAutoIncrement, otomatik artış olarak işaretler.
func (g *GormConfig) WithAutoIncrement() *GormConfig {
	g.AutoIncrement = true
	return g
}

// WithColumn, özel sütun adı belirler.
func (g *GormConfig) WithColumn(name string) *GormConfig {
	g.Column = name
	return g
}

// WithType, SQL tipini belirler.
func (g *GormConfig) WithType(sqlType string) *GormConfig {
	g.Type = sqlType
	return g
}

// WithSize, sütun boyutunu belirler.
func (g *GormConfig) WithSize(size int) *GormConfig {
	g.Size = size
	return g
}

// WithPrecision, ondalık hassasiyeti belirler.
func (g *GormConfig) WithPrecision(precision, scale int) *GormConfig {
	g.Precision = precision
	g.Scale = scale
	return g
}

// WithIndex, normal indeks oluşturur.
func (g *GormConfig) WithIndex(name ...string) *GormConfig {
	g.Index = true
	if len(name) > 0 {
		g.IndexName = name[0]
	}
	return g
}

// WithUniqueIndex, benzersiz indeks oluşturur.
func (g *GormConfig) WithUniqueIndex(name ...string) *GormConfig {
	g.UniqueIndex = true
	if len(name) > 0 {
		g.IndexName = name[0]
	}
	return g
}

// WithDefault, varsayılan değer belirler.
func (g *GormConfig) WithDefault(value interface{}) *GormConfig {
	g.Default = value
	return g
}

// WithNotNull, NOT NULL kısıtlaması ekler.
func (g *GormConfig) WithNotNull() *GormConfig {
	g.NotNull = true
	return g
}

// WithComment, sütun yorumu ekler.
func (g *GormConfig) WithComment(comment string) *GormConfig {
	g.Comment = comment
	return g
}

// WithForeignKey, foreign key ilişkisi tanımlar.
func (g *GormConfig) WithForeignKey(fk, references string) *GormConfig {
	g.ForeignKey = fk
	g.References = references
	return g
}

// WithOnDelete, ON DELETE davranışını belirler.
func (g *GormConfig) WithOnDelete(action string) *GormConfig {
	g.OnDelete = action
	return g
}

// WithOnUpdate, ON UPDATE davranışını belirler.
func (g *GormConfig) WithOnUpdate(action string) *GormConfig {
	g.OnUpdate = action
	return g
}

// WithManyToMany, many-to-many ilişki için ara tablo tanımlar.
func (g *GormConfig) WithManyToMany(tableName string) *GormConfig {
	g.ManyToManyTable = tableName
	return g
}

// ToGormTag, GormConfig'i GORM struct tag string'ine dönüştürür.
func (g *GormConfig) ToGormTag() string {
	var parts []string

	if g.PrimaryKey {
		parts = append(parts, "primaryKey")
	}
	if g.AutoIncrement {
		parts = append(parts, "autoIncrement")
	}
	if g.Column != "" {
		parts = append(parts, "column:"+g.Column)
	}
	if g.Type != "" {
		parts = append(parts, "type:"+g.Type)
	}
	if g.Size > 0 {
		parts = append(parts, "size:"+itoa(g.Size))
	}
	if g.Precision > 0 {
		parts = append(parts, "precision:"+itoa(g.Precision))
	}
	if g.Scale > 0 {
		parts = append(parts, "scale:"+itoa(g.Scale))
	}
	if g.Index {
		if g.IndexName != "" {
			parts = append(parts, "index:"+g.IndexName)
		} else {
			parts = append(parts, "index")
		}
	}
	if g.UniqueIndex {
		if g.IndexName != "" {
			parts = append(parts, "uniqueIndex:"+g.IndexName)
		} else {
			parts = append(parts, "uniqueIndex")
		}
	}
	if g.NotNull {
		parts = append(parts, "not null")
	}
	if g.Default != nil {
		parts = append(parts, "default:"+formatDefault(g.Default))
	}
	if g.Comment != "" {
		parts = append(parts, "comment:"+g.Comment)
	}
	if g.ManyToManyTable != "" {
		parts = append(parts, "many2many:"+g.ManyToManyTable)
	}

	return joinWithSemicolon(parts)
}

// Yardımcı fonksiyonlar
func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	result := ""
	for i > 0 {
		result = string(rune('0'+i%10)) + result
		i /= 10
	}
	return result
}

func formatDefault(v interface{}) string {
	switch val := v.(type) {
	case string:
		return "'" + val + "'"
	case bool:
		if val {
			return "true"
		}
		return "false"
	case int, int64, float64:
		return itoa(int(val.(int)))
	default:
		return ""
	}
}

func joinWithSemicolon(parts []string) string {
	if len(parts) == 0 {
		return ""
	}
	result := parts[0]
	for i := 1; i < len(parts); i++ {
		result += ";" + parts[i]
	}
	return result
}
