package fields

// RelationshipGormConfig, ilişki field'ları için GORM yapılandırmasını tanımlar.
// Bu yapı, foreign key, pivot tablo ve polymorphic ilişkiler için kullanılır.
type RelationshipGormConfig struct {
	// Foreign Key Yapılandırması (BelongsTo, HasOne, HasMany)
	ForeignKey     string `json:"-"` // Foreign key sütun adı (örn: "author_id")
	References     string `json:"-"` // Referans sütun (örn: "id")
	OnDelete       string `json:"-"` // ON DELETE davranışı: CASCADE, SET NULL, RESTRICT, NO ACTION
	OnUpdate       string `json:"-"` // ON UPDATE davranışı
	ConstraintName string `json:"-"` // Özel constraint adı

	// Many-to-Many Yapılandırması (BelongsToMany)
	PivotTable     string `json:"-"` // Ara tablo adı (örn: "user_roles")
	JoinForeignKey string `json:"-"` // Join tablosundaki ana model FK (örn: "user_id")
	JoinReferences string `json:"-"` // Join tablosundaki referans FK (örn: "role_id")
	JoinTable      string `json:"-"` // GORM joinTable alias

	// Polymorphic Yapılandırması (MorphTo, MorphOne, MorphMany, MorphToMany)
	Polymorphic      bool   `json:"-"` // Polymorphic ilişki mi?
	PolymorphicType  string `json:"-"` // Type sütun adı (örn: "commentable_type")
	PolymorphicID    string `json:"-"` // ID sütun adı (örn: "commentable_id")
	PolymorphicValue string `json:"-"` // Type değeri (örn: "posts", "users")

	// Eager Loading
	Preload bool `json:"-"` // Otomatik preload yapılsın mı?
}

// NewRelationshipGormConfig, varsayılan değerlerle yeni bir RelationshipGormConfig oluşturur.
func NewRelationshipGormConfig() *RelationshipGormConfig {
	return &RelationshipGormConfig{
		OnDelete: "CASCADE",
		OnUpdate: "CASCADE",
		Preload:  true,
	}
}

// Builder metodları

// WithForeignKey, foreign key sütununu belirler.
func (r *RelationshipGormConfig) WithForeignKey(fk string) *RelationshipGormConfig {
	r.ForeignKey = fk
	return r
}

// WithReferences, referans sütununu belirler.
func (r *RelationshipGormConfig) WithReferences(ref string) *RelationshipGormConfig {
	r.References = ref
	return r
}

// WithOnDelete, ON DELETE davranışını belirler.
func (r *RelationshipGormConfig) WithOnDelete(action string) *RelationshipGormConfig {
	r.OnDelete = action
	return r
}

// WithOnUpdate, ON UPDATE davranışını belirler.
func (r *RelationshipGormConfig) WithOnUpdate(action string) *RelationshipGormConfig {
	r.OnUpdate = action
	return r
}

// WithConstraintName, özel constraint adı belirler.
func (r *RelationshipGormConfig) WithConstraintName(name string) *RelationshipGormConfig {
	r.ConstraintName = name
	return r
}

// WithPivotTable, many-to-many için ara tablo belirler.
func (r *RelationshipGormConfig) WithPivotTable(table, joinFK, joinRef string) *RelationshipGormConfig {
	r.PivotTable = table
	r.JoinForeignKey = joinFK
	r.JoinReferences = joinRef
	return r
}

// WithPolymorphic, polymorphic ilişki yapılandırır.
func (r *RelationshipGormConfig) WithPolymorphic(typeColumn, idColumn string) *RelationshipGormConfig {
	r.Polymorphic = true
	r.PolymorphicType = typeColumn
	r.PolymorphicID = idColumn
	return r
}

// WithPolymorphicValue, polymorphic type değerini belirler.
func (r *RelationshipGormConfig) WithPolymorphicValue(value string) *RelationshipGormConfig {
	r.PolymorphicValue = value
	return r
}

// WithPreload, eager loading'i etkinleştirir/devre dışı bırakır.
func (r *RelationshipGormConfig) WithPreload(enabled bool) *RelationshipGormConfig {
	r.Preload = enabled
	return r
}

// ToGormTag, ilişki için GORM struct tag'ini oluşturur.
// Bu metod, model alanları için gorm tag'leri üretir.
func (r *RelationshipGormConfig) ToGormTag() string {
	var parts []string

	if r.ForeignKey != "" {
		parts = append(parts, "foreignKey:"+r.ForeignKey)
	}
	if r.References != "" {
		parts = append(parts, "references:"+r.References)
	}
	if r.PivotTable != "" {
		parts = append(parts, "many2many:"+r.PivotTable)
	}
	if r.JoinForeignKey != "" {
		parts = append(parts, "joinForeignKey:"+r.JoinForeignKey)
	}
	if r.JoinReferences != "" {
		parts = append(parts, "joinReferences:"+r.JoinReferences)
	}
	if r.Polymorphic {
		parts = append(parts, "polymorphic:"+r.PolymorphicID)
		if r.PolymorphicType != "" {
			parts = append(parts, "polymorphicType:"+r.PolymorphicType)
		}
		if r.PolymorphicValue != "" {
			parts = append(parts, "polymorphicValue:"+r.PolymorphicValue)
		}
	}
	if r.ConstraintName != "" {
		parts = append(parts, "constraint:"+r.ConstraintName)
	}

	return joinWithSemicolon(parts)
}

// ToForeignKeyTag, foreign key sütunu için GORM tag oluşturur.
// BelongsTo ilişkilerinde kullanılır (örn: author_id uint `gorm:"index"`).
func (r *RelationshipGormConfig) ToForeignKeyTag() string {
	parts := []string{"index"}

	if r.OnDelete != "" || r.OnUpdate != "" {
		constraint := "constraint:"
		if r.ConstraintName != "" {
			constraint += r.ConstraintName
		} else {
			constraint += "OnDelete:" + r.OnDelete
			if r.OnUpdate != "" {
				constraint += ",OnUpdate:" + r.OnUpdate
			}
		}
		parts = append(parts, constraint)
	}

	return joinWithSemicolon(parts)
}
