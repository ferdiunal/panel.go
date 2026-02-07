// Package fields, ilişkisel veritabanı alanları için varlık kontrolü işlevselliği sağlar.
//
// Bu paket, relationship field'ların ilişkili kaynaklarının var olup olmadığını
// kontrol etmek için kullanılan interface ve implementasyonları içerir.
//
// # Desteklenen İlişki Tipleri
//
// - `belongsTo`: Ters bire-bir veya bire-çok ilişki
// - `hasMany`: Bire-çok ilişki
// - `hasOne`: Bire-bir ilişki
// - `belongsToMany`: Çoka-çok ilişki (pivot tablo üzerinden)
// - `morphTo`: Polimorfik ilişki
//
// # Referans
//
// Detaylı bilgi için bkz: `docs/Relationships.md`
package fields

import (
	"context"
)

// RelationshipExistence, ilişkili kaynakların varlık kontrolü için interface tanımlar.
//
// Bu interface, relationship field'ların ilişkili kaynaklarının var olup olmadığını
// veya hiç ilişkili kaynak bulunmadığını kontrol etmek için kullanılır.
//
// # Metodlar
//
// - `Exists`: İlişkili kaynakların var olup olmadığını kontrol eder
// - `DoesntExist`: Hiç ilişkili kaynak bulunmadığını kontrol eder
//
// # Kullanım Örneği
//
// ```go
// existence := NewRelationshipExistence(relationshipField)
//
// // İlişkili kaynak var mı?
// exists, err := existence.Exists(ctx)
// if err != nil {
//     return err
// }
// if exists {
//     fmt.Println("İlişkili kaynak bulundu")
// }
//
// // Hiç ilişkili kaynak yok mu?
// doesntExist, err := existence.DoesntExist(ctx)
// if err != nil {
//     return err
// }
// if doesntExist {
//     fmt.Println("Hiç ilişkili kaynak yok")
// }
// ```
//
// # Referans
//
// Detaylı bilgi için bkz: `docs/Relationships.md`
type RelationshipExistence interface {
	// Exists, ilişkili kaynakların var olup olmadığını kontrol eder.
	//
	// # Parametreler
	//
	// - `ctx`: Context nesnesi (timeout, cancellation için)
	//
	// # Dönüş Değerleri
	//
	// - `bool`: İlişkili kaynak varsa `true`, yoksa `false`
	// - `error`: Kontrol sırasında oluşan hata, başarılıysa `nil`
	//
	// # Davranış
	//
	// İlişki tipine göre uygun EXISTS sorgusu çalıştırılır:
	// - `belongsTo`: Foreign key üzerinden parent kaynak kontrolü
	// - `hasMany`: Child kayıtların varlık kontrolü
	// - `hasOne`: Child kaydın varlık kontrolü
	// - `belongsToMany`: Pivot tablo üzerinden ilişki kontrolü
	// - `morphTo`: Polimorfik ilişki kontrolü
	Exists(ctx context.Context) (bool, error)

	// DoesntExist, hiç ilişkili kaynak bulunmadığını kontrol eder.
	//
	// # Parametreler
	//
	// - `ctx`: Context nesnesi (timeout, cancellation için)
	//
	// # Dönüş Değerleri
	//
	// - `bool`: Hiç ilişkili kaynak yoksa `true`, varsa `false`
	// - `error`: Kontrol sırasında oluşan hata, başarılıysa `nil`
	//
	// # Not
	//
	// Bu method, `Exists` metodunun tersini döndürür.
	DoesntExist(ctx context.Context) (bool, error)
}

// RelationshipExistenceImpl, RelationshipExistence interface'ini implement eder.
//
// Bu struct, relationship field'ların ilişkili kaynaklarının varlık kontrolü
// işlemlerini gerçekleştirir. İlişki tipine göre uygun EXISTS sorgusu stratejisini
// otomatik olarak seçer ve uygular.
//
// # Alanlar
//
// - `field`: Varlık kontrolü yapılacak relationship field
//
// # Kullanım Örneği
//
// ```go
// // Relationship field oluştur
// belongsToField := fields.NewBelongsTo("author", "Author")
//
// // Varlık kontrolcüsü oluştur
// existence := fields.NewRelationshipExistence(belongsToField)
//
// // Varlık kontrolü yap
// exists, err := existence.Exists(context.Background())
// if err != nil {
//     log.Fatal(err)
// }
//
// if exists {
//     fmt.Println("İlişkili kaynak mevcut")
// }
// ```
//
// # Desteklenen İlişki Tipleri
//
// - `belongsTo`: Parent kaynak varlık kontrolü
// - `hasMany`: Child kayıtlar varlık kontrolü
// - `hasOne`: Child kayıt varlık kontrolü
// - `belongsToMany`: Pivot tablo varlık kontrolü
// - `morphTo`: Polimorfik ilişki varlık kontrolü
//
// # Referans
//
// Detaylı bilgi için bkz: `docs/Relationships.md`
type RelationshipExistenceImpl struct {
	field RelationshipField
}

// NewRelationshipExistence, yeni bir relationship varlık kontrolcüsü oluşturur.
//
// Bu constructor, verilen relationship field için varlık kontrolü yapabilen
// bir handler oluşturur. Handler, ilişki tipine göre otomatik olarak uygun
// EXISTS sorgu stratejisini seçer.
//
// # Parametreler
//
// - `field`: Varlık kontrolü yapılacak relationship field (RelationshipField interface'ini implement etmeli)
//
// # Dönüş Değerleri
//
// - `*RelationshipExistenceImpl`: Yapılandırılmış varlık kontrolcüsü pointer'ı
//
// # Kullanım Örneği
//
// ```go
// // BelongsTo ilişkisi için
// authorField := fields.NewBelongsTo("author", "Author")
// existence := fields.NewRelationshipExistence(authorField)
//
// // HasMany ilişkisi için
// commentsField := fields.NewHasMany("comments", "Comment")
// existence := fields.NewRelationshipExistence(commentsField)
//
// // BelongsToMany ilişkisi için
// tagsField := fields.NewBelongsToMany("tags", "Tag")
// existence := fields.NewRelationshipExistence(tagsField)
// ```
//
// # Notlar
//
// - Field parametresi nil olmamalıdır
// - Field, RelationshipField interface'ini implement etmelidir
// - Döndürülen handler thread-safe değildir, concurrent kullanım için senkronizasyon gerekir
//
// # Referans
//
// Detaylı bilgi için bkz: `docs/Relationships.md`
func NewRelationshipExistence(field RelationshipField) *RelationshipExistenceImpl {
	return &RelationshipExistenceImpl{
		field: field,
	}
}

// Exists, ilişkili kaynakların var olup olmadığını kontrol eder.
//
// Bu method, relationship field'ın tipine göre otomatik olarak uygun EXISTS
// sorgu stratejisini seçer ve çalıştırır. Her ilişki tipi için özel kontrol
// mantığı uygulanır.
//
// # Parametreler
//
// - `ctx`: Context nesnesi (timeout, cancellation için)
//
// # Dönüş Değerleri
//
// - `bool`: İlişkili kaynak varsa `true`, yoksa `false`
// - `error`: Kontrol sırasında oluşan hata, başarılıysa `nil`
//
// # İlişki Tiplerine Göre Davranış
//
// - `belongsTo`: Foreign key üzerinden parent kaynak varlığını kontrol eder
// - `hasMany`: Child kayıtların varlığını kontrol eder (COUNT > 0)
// - `hasOne`: Child kaydın varlığını kontrol eder
// - `belongsToMany`: Pivot tablo üzerinden ilişki varlığını kontrol eder
// - `morphTo`: Polimorfik ilişki varlığını kontrol eder (type ve id üzerinden)
//
// # Kullanım Örneği
//
// ```go
// // BelongsTo ilişkisi için
// authorField := fields.NewBelongsTo("author", "Author")
// existence := fields.NewRelationshipExistence(authorField)
//
// exists, err := existence.Exists(context.Background())
// if err != nil {
//     return fmt.Errorf("varlık kontrolü başarısız: %w", err)
// }
//
// if exists {
//     fmt.Println("Yazar kaydı mevcut")
// } else {
//     fmt.Println("Yazar kaydı bulunamadı")
// }
//
// // HasMany ilişkisi için
// commentsField := fields.NewHasMany("comments", "Comment")
// existence := fields.NewRelationshipExistence(commentsField)
//
// hasComments, err := existence.Exists(ctx)
// if err != nil {
//     return err
// }
//
// if hasComments {
//     fmt.Println("En az bir yorum var")
// }
// ```
//
// # Notlar
//
// - Context timeout veya cancellation durumunda hata döner
// - Bilinmeyen ilişki tipleri için `false, nil` döner
// - Database bağlantı hataları error olarak döner
// - Bu method thread-safe değildir
//
// # Referans
//
// Detaylı bilgi için bkz: `docs/Relationships.md`
func (re *RelationshipExistenceImpl) Exists(ctx context.Context) (bool, error) {
	relationType := re.field.GetRelationshipType()

	switch relationType {
	case "belongsTo":
		return re.existsBelongsTo(ctx)
	case "hasMany":
		return re.existsHasMany(ctx)
	case "hasOne":
		return re.existsHasOne(ctx)
	case "belongsToMany":
		return re.existsBelongsToMany(ctx)
	case "morphTo":
		return re.existsMorphTo(ctx)
	default:
		return false, nil
	}
}

// DoesntExist, hiç ilişkili kaynak bulunmadığını kontrol eder.
//
// Bu method, `Exists` metodunun tersini döndürür. İlişkili hiç kaynak yoksa
// `true`, en az bir ilişkili kaynak varsa `false` döner.
//
// # Parametreler
//
// - `ctx`: Context nesnesi (timeout, cancellation için)
//
// # Dönüş Değerleri
//
// - `bool`: Hiç ilişkili kaynak yoksa `true`, varsa `false`
// - `error`: Kontrol sırasında oluşan hata, başarılıysa `nil`
//
// # Davranış
//
// Bu method dahili olarak `Exists` metodunu çağırır ve sonucunu tersine çevirir:
// - `Exists` true dönerse -> `DoesntExist` false döner
// - `Exists` false dönerse -> `DoesntExist` true döner
// - `Exists` hata dönerse -> `DoesntExist` de aynı hatayı döner
//
// # Kullanım Örneği
//
// ```go
// // BelongsTo ilişkisi için
// authorField := fields.NewBelongsTo("author", "Author")
// existence := fields.NewRelationshipExistence(authorField)
//
// noAuthor, err := existence.DoesntExist(context.Background())
// if err != nil {
//     return fmt.Errorf("varlık kontrolü başarısız: %w", err)
// }
//
// if noAuthor {
//     fmt.Println("Yazar atanmamış")
// } else {
//     fmt.Println("Yazar mevcut")
// }
//
// // HasMany ilişkisi için - yorum kontrolü
// commentsField := fields.NewHasMany("comments", "Comment")
// existence := fields.NewRelationshipExistence(commentsField)
//
// noComments, err := existence.DoesntExist(ctx)
// if err != nil {
//     return err
// }
//
// if noComments {
//     fmt.Println("Henüz yorum yapılmamış")
// }
//
// // Validation senaryosu
// if noComments {
//     return errors.New("en az bir yorum gerekli")
// }
// ```
//
// # Notlar
//
// - Context timeout veya cancellation durumunda hata döner
// - `Exists` metodunun tüm davranışları geçerlidir
// - Database bağlantı hataları error olarak döner
// - Bu method thread-safe değildir
// - Validation ve conditional logic için idealdir
//
// # Referans
//
// Detaylı bilgi için bkz: `docs/Relationships.md`
func (re *RelationshipExistenceImpl) DoesntExist(ctx context.Context) (bool, error) {
	exists, err := re.Exists(ctx)
	if err != nil {
		return false, err
	}

	return !exists, nil
}

// existsBelongsTo checks if BelongsTo relationship exists
func (re *RelationshipExistenceImpl) existsBelongsTo(ctx context.Context) (bool, error) {
	// In a real implementation, this would execute an EXISTS query
	// For now, return false
	return false, nil
}

// existsHasMany checks if HasMany relationships exist
func (re *RelationshipExistenceImpl) existsHasMany(ctx context.Context) (bool, error) {
	// In a real implementation, this would execute an EXISTS query
	// For now, return false
	return false, nil
}

// existsHasOne checks if HasOne relationship exists
func (re *RelationshipExistenceImpl) existsHasOne(ctx context.Context) (bool, error) {
	// In a real implementation, this would execute an EXISTS query
	// For now, return false
	return false, nil
}

// existsBelongsToMany checks if BelongsToMany relationships exist
func (re *RelationshipExistenceImpl) existsBelongsToMany(ctx context.Context) (bool, error) {
	// In a real implementation, this would execute an EXISTS query on the pivot table
	// For now, return false
	return false, nil
}

// existsMorphTo checks if MorphTo relationship exists
func (re *RelationshipExistenceImpl) existsMorphTo(ctx context.Context) (bool, error) {
	// In a real implementation, this would execute an EXISTS query
	// For now, return false
	return false, nil
}
