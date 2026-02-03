package fields

import "strings"

// NewField, yeni bir alan şeması oluşturur.
func NewField(name string, attribute ...string) *Schema {
	attr := strings.ToLower(strings.ReplaceAll(name, " ", "_"))
	if len(attribute) > 0 {
		attr = attribute[0]
	}

	return &Schema{
		Name:      name,
		Key:       attr,
		Props:     make(map[string]interface{}),
		TextAlign: "left",
	}
}

// Specific Field Implementations

// ID, benzersiz kimlik alanını oluşturur. Varsayılan olarak "id" anahtarını kullanır ve sadece listede görünür.
func ID(name ...string) *Schema {
	label := "ID"
	if len(name) > 0 {
		label = name[0]
	}
	f := NewField(label, "id")
	f.View = "id-field" // Frontend bileşen adı
	f.OnlyOnList()
	return f
}

// Text, standart metin giriş alanı (input type="text") oluşturur.
func Text(name string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "text-field"
	f.Type = TYPE_TEXT
	return f
}

// Password, şifre giriş alanı (input type="password") oluşturur.
func Password(name string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "password-field"
	f.Type = TYPE_PASSWORD
	return f
}

// Number, sayı giriş alanı (input type="number") oluşturur.
func Number(name string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "number-field"
	f.Type = TYPE_NUMBER
	return f
}

// Email, e-posta giriş alanı (input type="email") oluşturur.
func Email(name string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "email-field"
	f.Type = TYPE_EMAIL
	return f
}

// Image, görsel yükleme alanı oluşturur.
func Image(name string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "image-field"
	f.Type = TYPE_FILE
	return f
}

// Tel, telefon numarası giriş alanı oluşturur.
func Tel(name string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "text-field" // Genellikle text input kullanılır
	f.Type = TYPE_TEL
	return f
}

// Video, video yükleme alanı oluşturur.
func Video(name string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "file-field"
	f.Type = TYPE_VIDEO
	return f
}

// Audio, ses dosyası yükleme alanı oluşturur.
func Audio(name string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "file-field"
	f.Type = TYPE_AUDIO
	return f
}

// Date, tarih seçim alanı (datepicker) oluşturur.
func Date(name string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "date-field"
	f.Type = TYPE_DATE
	return f
}

// DateTime, tarih ve saat seçim alanı oluşturur.
func DateTime(name string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "datetime-field"
	f.Type = TYPE_DATETIME
	return f
}

// File, dosya yükleme alanı oluşturur.
func File(name string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "file-field"
	f.Type = TYPE_FILE
	return f
}

// KeyValue, anahtar-değer ikilisi girişi sağlayan alan oluşturur.
func KeyValue(name string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "key-value-field"
	f.Type = TYPE_KEY_VALUE
	return f
}

// Relationships

// Link, başka bir kaynağa bağlantı (BelongsTo) oluşturur.
func Link(name string, resource string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "link-field"
	f.Type = TYPE_LINK
	f.Props["resource"] = resource
	return f
}

// Detail, bir kaynağın detayını gösteren (HasOne) alan oluşturur. Genellikle listede gizlenir.
func Detail(name string, resource string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "detail-field"
	f.Type = TYPE_DETAIL
	f.Props["resource"] = resource
	f.Context = HIDE_ON_LIST // Generally hidden on list
	return f
}

// Collection, ilişkili kayıtların listesini gösteren (HasMany) alan oluşturur.
func Collection(name string, resource string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "collection-field"
	f.Type = TYPE_COLLECTION
	f.Props["resource"] = resource
	f.Context = HIDE_ON_LIST
	return f
}

// Connect, çoktan çoka (BelongsToMany) ilişki kurmak için kullanılır.
func Connect(name string, resource string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "connect-field"
	f.Type = TYPE_CONNECT
	f.Props["resource"] = resource
	f.Context = HIDE_ON_LIST
	return f
}

// PolyLink, polimorfik ilişki bağlantısı (MorphTo) oluşturur.
func PolyLink(name string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "poly-link-field"
	f.Type = TYPE_POLY_LINK
	return f
}

// PolyDetail, polimorfik detay (MorphOne) oluşturur.
func PolyDetail(name string, resource string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "poly-detail-field"
	f.Type = TYPE_POLY_DETAIL
	f.Props["resource"] = resource
	f.Context = HIDE_ON_LIST
	return f
}

// PolyCollection, polimorfik koleksiyon (MorphMany) oluşturur.
func PolyCollection(name string, resource string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "poly-collection-field"
	f.Type = TYPE_POLY_COLLECTION
	f.Props["resource"] = resource
	f.Context = HIDE_ON_LIST
	return f
}

func PolyConnect(name string, resource string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "poly-connect-field"
	f.Type = TYPE_POLY_CONNECT
	f.Props["resource"] = resource
	f.Context = HIDE_ON_LIST
	return f
}

// Switch, boolean değerler için switch/toggle bileşeni oluşturur.
func Switch(name string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "switch-field"
	f.Type = TYPE_BOOLEAN
	return f
}

// Combobox, çoktan seçmeli veya arama yapılabilir seçim alanı oluşturur.
func Combobox(name string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "combobox-field"
	f.Type = TYPE_SELECT // veya TYPE_COMBOBOX eğer yeni tip gerekiyorsa, şimdilik Select mantığında.
	return f
}

// Select, standart seçim listesi oluşturur.
func Select(name string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "select-field"
	f.Type = TYPE_SELECT
	return f
}
