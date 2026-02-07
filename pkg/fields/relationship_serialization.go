// Package fields, ilişki serileştirme işlemlerini yönetir.
//
// Bu paket, relationship field'larının JSON formatına dönüştürülmesi için
// gerekli interface ve implementasyonları sağlar.
//
// # Referans
//
// Detaylı bilgi için: [docs/Relationships.md](../../docs/Relationships.md)
package fields

import (
	"encoding/json"
)

// RelationshipSerialization, ilişkilerin JSON serileştirme işlemlerini yönetir.
//
// Bu interface, relationship field'larının JSON formatına dönüştürülmesi için
// gerekli metodları tanımlar. Hem tekil hem de çoklu ilişkilerin serileştirmesini destekler.
//
// # Kullanım
//
// ```go
// serializer := NewRelationshipSerialization(relationshipField)
//
// // Tekil ilişki serileştirme
// data, err := serializer.SerializeRelationship(user)
// if err != nil {
//     log.Fatal(err)
// }
//
// // Çoklu ilişki serileştirme
// items, err := serializer.SerializeRelationships(users)
// if err != nil {
//     log.Fatal(err)
// }
// ```
//
// # Referans
//
// Detaylı bilgi için: [docs/Relationships.md](../../docs/Relationships.md)
type RelationshipSerialization interface {
	// SerializeRelationship, tekil bir ilişkiyi JSON formatına dönüştürür.
	//
	// # Parametreler
	//
	// - `item`: Serileştirilecek ilişki verisi (nil olabilir)
	//
	// # Döndürür
	//
	// - `map[string]interface{}`: Serileştirilmiş veri (type, name, resource, value içerir)
	// - `error`: Hata durumunda hata mesajı
	//
	// # Örnek
	//
	// ```go
	// data, err := serializer.SerializeRelationship(user)
	// // Sonuç: {"type": "belongsTo", "name": "author", "resource": "users", "value": {...}}
	// ```
	SerializeRelationship(item interface{}) (map[string]interface{}, error)

	// SerializeRelationships, çoklu ilişkileri JSON formatına dönüştürür.
	//
	// # Parametreler
	//
	// - `items`: Serileştirilecek ilişki verileri dizisi (nil veya boş olabilir)
	//
	// # Döndürür
	//
	// - `[]map[string]interface{}`: Serileştirilmiş veri dizisi
	// - `error`: Hata durumunda hata mesajı
	//
	// # Örnek
	//
	// ```go
	// items, err := serializer.SerializeRelationships(users)
	// // Sonuç: [{"type": "hasMany", "name": "posts", "resource": "posts", "value": {...}}, ...]
	// ```
	SerializeRelationships(items []interface{}) ([]map[string]interface{}, error)
}

// RelationshipSerializationImpl, RelationshipSerialization interface'ini implement eder.
//
// Bu struct, relationship field'larının JSON serileştirme işlemlerini gerçekleştirir.
// İlişki tipini, adını ve ilgili resource bilgilerini içeren JSON çıktısı üretir.
//
// # Alanlar
//
// - `field`: Serileştirilecek relationship field
//
// # Örnek
//
// ```go
// impl := &RelationshipSerializationImpl{
//     field: belongsToField,
// }
// ```
type RelationshipSerializationImpl struct {
	field RelationshipField
}

// NewRelationshipSerialization, yeni bir relationship serialization handler oluşturur.
//
// Bu constructor, verilen relationship field için bir serileştirme handler'ı başlatır.
//
// # Parametreler
//
// - `field`: Serileştirilecek relationship field
//
// # Döndürür
//
// - Yapılandırılmış RelationshipSerializationImpl pointer'ı
//
// # Örnek
//
// ```go
// belongsTo := fields.NewBelongsTo("author").
//     SetRelatedResource("users")
//
// serializer := fields.NewRelationshipSerialization(belongsTo)
//
// // Kullanım
// data, err := serializer.SerializeRelationship(user)
// ```
//
// # Notlar
//
// - Field parametresi nil olmamalıdır
// - Döndürülen handler tüm serileştirme metodlarını destekler
func NewRelationshipSerialization(field RelationshipField) *RelationshipSerializationImpl {
	return &RelationshipSerializationImpl{
		field: field,
	}
}

// SerializeRelationship, tekil bir ilişkiyi JSON formatına dönüştürür.
//
// Bu method, verilen ilişki verisini JSON-uyumlu bir map yapısına dönüştürür.
// İlişki tipi, adı, resource bilgisi ve değeri içeren bir yapı oluşturur.
//
// # Parametreler
//
// - `item`: Serileştirilecek ilişki verisi (nil olabilir)
//
// # Döndürür
//
// - `map[string]interface{}`: Serileştirilmiş veri yapısı
//   - `type`: İlişki tipi (belongsTo, hasMany, hasOne, belongsToMany, morphTo)
//   - `name`: İlişki adı
//   - `resource`: İlgili resource adı
//   - `value`: İlişki verisi (nil olabilir)
// - `error`: Hata durumunda hata mesajı (şu an için her zaman nil)
//
// # Örnek
//
// ```go
// serializer := NewRelationshipSerialization(belongsToField)
// data, err := serializer.SerializeRelationship(user)
// // Sonuç: {
// //   "type": "belongsTo",
// //   "name": "author",
// //   "resource": "users",
// //   "value": {...}
// // }
// ```
//
// # Notlar
//
// - Nil değerler için sadece `value: nil` içeren bir map döner
// - İlişki metadata'sı field'dan otomatik olarak alınır
func (rs *RelationshipSerializationImpl) SerializeRelationship(item interface{}) (map[string]interface{}, error) {
	if item == nil {
		return map[string]interface{}{
			"value": nil,
		}, nil
	}

	// Convert item to JSON-compatible format
	jsonData := map[string]interface{}{
		"type":     rs.field.GetRelationshipType(),
		"name":     rs.field.GetRelationshipName(),
		"resource": rs.field.GetRelatedResource(),
		"value":    item,
	}

	return jsonData, nil
}

// SerializeRelationships, çoklu ilişkileri JSON formatına dönüştürür.
//
// Bu method, verilen ilişki verileri dizisini JSON-uyumlu map dizisine dönüştürür.
// Her bir ilişki için SerializeRelationship metodunu çağırır.
//
// # Parametreler
//
// - `items`: Serileştirilecek ilişki verileri dizisi (nil veya boş olabilir)
//
// # Döndürür
//
// - `[]map[string]interface{}`: Serileştirilmiş veri dizisi (boş dizi olabilir)
// - `error`: Herhangi bir item'ın serileştirmesinde hata oluşursa hata mesajı
//
// # Örnek
//
// ```go
// serializer := NewRelationshipSerialization(hasManyField)
// items, err := serializer.SerializeRelationships(posts)
// if err != nil {
//     log.Fatal(err)
// }
// // Sonuç: [
// //   {"type": "hasMany", "name": "posts", "resource": "posts", "value": {...}},
// //   {"type": "hasMany", "name": "posts", "resource": "posts", "value": {...}}
// // ]
// ```
//
// # Notlar
//
// - Nil veya boş dizi için boş slice döner (nil değil)
// - Her item için ayrı ayrı SerializeRelationship çağrılır
// - Herhangi bir item'da hata oluşursa işlem durur ve hata döner
func (rs *RelationshipSerializationImpl) SerializeRelationships(items []interface{}) ([]map[string]interface{}, error) {
	if items == nil || len(items) == 0 {
		return []map[string]interface{}{}, nil
	}

	serialized := make([]map[string]interface{}, 0, len(items))

	for _, item := range items {
		jsonData, err := rs.SerializeRelationship(item)
		if err != nil {
			return nil, err
		}
		serialized = append(serialized, jsonData)
	}

	return serialized, nil
}

// ToJSON, ilişkiyi JSON string formatına dönüştürür.
//
// Bu method, SerializeRelationship metodunu kullanarak ilişkiyi önce map yapısına,
// ardından JSON string formatına dönüştürür.
//
// # Parametreler
//
// - `item`: JSON'a dönüştürülecek ilişki verisi (nil olabilir)
//
// # Döndürür
//
// - `string`: JSON formatında string (örn: `{"type":"belongsTo","name":"author",...}`)
// - `error`: Serileştirme veya JSON encoding hatası
//
// # Örnek
//
// ```go
// serializer := NewRelationshipSerialization(belongsToField)
// jsonStr, err := serializer.ToJSON(user)
// if err != nil {
//     log.Fatal(err)
// }
// fmt.Println(jsonStr)
// // Çıktı: {"type":"belongsTo","name":"author","resource":"users","value":{...}}
// ```
//
// # Notlar
//
// - İç içe SerializeRelationship ve json.Marshal çağrıları yapar
// - Her iki aşamada da hata oluşabilir
// - Nil değerler için `{"value":null}` formatında JSON döner
func (rs *RelationshipSerializationImpl) ToJSON(item interface{}) (string, error) {
	jsonData, err := rs.SerializeRelationship(item)
	if err != nil {
		return "", err
	}

	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

// ToJSONArray, çoklu ilişkileri JSON array string formatına dönüştürür.
//
// Bu method, SerializeRelationships metodunu kullanarak ilişkileri önce map dizisine,
// ardından JSON array string formatına dönüştürür.
//
// # Parametreler
//
// - `items`: JSON array'e dönüştürülecek ilişki verileri dizisi (nil veya boş olabilir)
//
// # Döndürür
//
// - `string`: JSON array formatında string (örn: `[{"type":"hasMany",...},{...}]`)
// - `error`: Serileştirme veya JSON encoding hatası
//
// # Örnek
//
// ```go
// serializer := NewRelationshipSerialization(hasManyField)
// jsonStr, err := serializer.ToJSONArray(posts)
// if err != nil {
//     log.Fatal(err)
// }
// fmt.Println(jsonStr)
// // Çıktı: [{"type":"hasMany","name":"posts","resource":"posts","value":{...}},...]
// ```
//
// # Notlar
//
// - İç içe SerializeRelationships ve json.Marshal çağrıları yapar
// - Her iki aşamada da hata oluşabilir
// - Nil veya boş dizi için `[]` formatında JSON döner
// - HasMany ve BelongsToMany ilişkileri için idealdir
func (rs *RelationshipSerializationImpl) ToJSONArray(items []interface{}) (string, error) {
	jsonData, err := rs.SerializeRelationships(items)
	if err != nil {
		return "", err
	}

	jsonBytes, err := json.Marshal(jsonData)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}
