// Package fields, ilişkisel alanlar için kısıtlama (constraint) yönetimi sağlar.
//
// Bu paket, relationship field'ları üzerinde sorgu kısıtlamaları uygulamak için
// gerekli interface ve implementasyonları içerir. LIMIT, OFFSET, WHERE ve WHERE IN
// gibi SQL kısıtlamalarını yönetir.
//
// # Referans
//
// Detaylı kullanım ve örnekler için: [docs/Relationships.md](../../docs/Relationships.md)
package fields

import (
	"context"
)

// RelationshipConstraints, ilişkisel alanlar için sorgu kısıtlamaları uygulayan interface'dir.
//
// Bu interface, relationship field'ları üzerinde LIMIT, OFFSET, WHERE ve WHERE IN
// gibi SQL kısıtlamalarını uygulamak için metodlar sağlar. Her metod context alır
// ve sonuçları döndürür.
//
// # Metodlar
//
// - `ApplyLimit`: Sorgu sonuç sayısını sınırlar
// - `ApplyOffset`: Sorgu başlangıç noktasını kaydırır
// - `ApplyWhere`: WHERE koşulu ekler
// - `ApplyWhereIn`: WHERE IN koşulu ekler
//
// # Örnek Kullanım
//
// ```go
// constraints := NewRelationshipConstraints(relationshipField)
//
// // Limit uygula
// results, err := constraints.ApplyLimit(ctx, 10)
//
// // Offset uygula
// results, err = constraints.ApplyOffset(ctx, 5)
//
// // WHERE koşulu ekle
// results, err = constraints.ApplyWhere(ctx, "status", "=", "active")
//
// // WHERE IN koşulu ekle
// results, err = constraints.ApplyWhereIn(ctx, "id", []interface{}{1, 2, 3})
// ```
//
// # Referans
//
// Detaylı bilgi için: [docs/Relationships.md](../../docs/Relationships.md)
type RelationshipConstraints interface {
	// ApplyLimit, sorgu sonuç sayısını sınırlar.
	//
	// # Parametreler
	//
	// - `ctx`: Context nesnesi
	// - `limit`: Maksimum sonuç sayısı (negatif değerler 0'a dönüştürülür)
	//
	// # Döndürür
	//
	// - `[]interface{}`: Kısıtlanmış sonuç listesi
	// - `error`: Hata durumunda hata nesnesi
	//
	// # Örnek
	//
	// ```go
	// results, err := constraints.ApplyLimit(ctx, 10)
	// if err != nil {
	//     log.Fatal(err)
	// }
	// ```
	ApplyLimit(ctx context.Context, limit int) ([]interface{}, error)

	// ApplyOffset, sorgu başlangıç noktasını kaydırır.
	//
	// # Parametreler
	//
	// - `ctx`: Context nesnesi
	// - `offset`: Kaç kayıt atlanacağı (negatif değerler 0'a dönüştürülür)
	//
	// # Döndürür
	//
	// - `[]interface{}`: Kaydırılmış sonuç listesi
	// - `error`: Hata durumunda hata nesnesi
	//
	// # Örnek
	//
	// ```go
	// // İlk 5 kaydı atla
	// results, err := constraints.ApplyOffset(ctx, 5)
	// if err != nil {
	//     log.Fatal(err)
	// }
	// ```
	ApplyOffset(ctx context.Context, offset int) ([]interface{}, error)

	// ApplyWhere, sorguya WHERE koşulu ekler.
	//
	// # Parametreler
	//
	// - `ctx`: Context nesnesi
	// - `column`: Koşul uygulanacak sütun adı
	// - `operator`: Karşılaştırma operatörü (=, !=, >, <, >=, <=, LIKE, vb.)
	// - `value`: Karşılaştırılacak değer
	//
	// # Döndürür
	//
	// - `[]interface{}`: Filtrelenmiş sonuç listesi
	// - `error`: Hata durumunda hata nesnesi
	//
	// # Örnek
	//
	// ```go
	// // Aktif kayıtları getir
	// results, err := constraints.ApplyWhere(ctx, "status", "=", "active")
	//
	// // Fiyatı 100'den büyük olanları getir
	// results, err = constraints.ApplyWhere(ctx, "price", ">", 100)
	// ```
	//
	// # Not
	//
	// Boş column değeri durumunda işlem yapılmaz ve boş liste döner.
	ApplyWhere(ctx context.Context, column string, operator string, value interface{}) ([]interface{}, error)

	// ApplyWhereIn, sorguya WHERE IN koşulu ekler.
	//
	// # Parametreler
	//
	// - `ctx`: Context nesnesi
	// - `column`: Koşul uygulanacak sütun adı
	// - `values`: Kontrol edilecek değerler listesi
	//
	// # Döndürür
	//
	// - `[]interface{}`: Filtrelenmiş sonuç listesi
	// - `error`: Hata durumunda hata nesnesi
	//
	// # Örnek
	//
	// ```go
	// // Belirli ID'lere sahip kayıtları getir
	// results, err := constraints.ApplyWhereIn(ctx, "id", []interface{}{1, 2, 3, 5, 8})
	//
	// // Belirli kategorilerdeki kayıtları getir
	// results, err = constraints.ApplyWhereIn(ctx, "category", []interface{}{"tech", "science", "art"})
	// ```
	//
	// # Not
	//
	// Boş column veya values durumunda işlem yapılmaz ve boş liste döner.
	ApplyWhereIn(ctx context.Context, column string, values []interface{}) ([]interface{}, error)
}

// RelationshipConstraintsImpl, RelationshipConstraints interface'inin implementasyonudur.
//
// Bu struct, relationship field'ları üzerinde sorgu kısıtlamalarını yönetir.
// Limit, offset ve WHERE koşullarını saklar ve uygular.
//
// # Alanlar
//
// - `field`: İlişkili field nesnesi
// - `limit`: Maksimum sonuç sayısı (0 = sınırsız)
// - `offset`: Atlanacak kayıt sayısı
// - `constraints`: WHERE koşulları map'i
//
// # Örnek Kullanım
//
// ```go
// // Yeni constraints handler oluştur
// constraints := NewRelationshipConstraints(relationshipField)
//
// // Kısıtlamaları uygula
// constraints.ApplyLimit(ctx, 10)
// constraints.ApplyOffset(ctx, 5)
// constraints.ApplyWhere(ctx, "status", "=", "active")
//
// // Değerleri oku
// limit := constraints.GetLimit()        // 10
// offset := constraints.GetOffset()      // 5
// where := constraints.GetConstraints()  // map[status:...]
// ```
//
// # Referans
//
// Detaylı bilgi için: [docs/Relationships.md](../../docs/Relationships.md)
type RelationshipConstraintsImpl struct {
	field       RelationshipField
	limit       int
	offset      int
	constraints map[string]interface{}
}

// NewRelationshipConstraints, yeni bir relationship constraints handler oluşturur.
//
// Bu constructor, verilen relationship field için kısıtlama yöneticisi oluşturur.
// Tüm kısıtlamalar varsayılan değerlerle (limit=0, offset=0, boş constraints)
// başlatılır.
//
// # Parametreler
//
// - `field`: Kısıtlamaların uygulanacağı RelationshipField nesnesi
//
// # Döndürür
//
// - Yapılandırılmış RelationshipConstraintsImpl pointer'ı
//
// # Örnek
//
// ```go
// // Relationship field oluştur
// relationshipField := fields.NewBelongsTo("user", "User")
//
// // Constraints handler oluştur
// constraints := NewRelationshipConstraints(relationshipField)
//
// // Kısıtlamaları uygula
// results, err := constraints.ApplyLimit(ctx, 10)
// if err != nil {
//     log.Fatal(err)
// }
// ```
//
// # Not
//
// - Limit ve offset başlangıçta 0 olarak ayarlanır
// - Constraints map'i boş olarak başlatılır
// - Field parametresi nil olmamalıdır
//
// # Referans
//
// Detaylı bilgi için: [docs/Relationships.md](../../docs/Relationships.md)
func NewRelationshipConstraints(field RelationshipField) *RelationshipConstraintsImpl {
	return &RelationshipConstraintsImpl{
		field:       field,
		limit:       0,
		offset:      0,
		constraints: make(map[string]interface{}),
	}
}

// ApplyLimit, sorgu sonuç sayısını sınırlar.
//
// Bu metod, relationship sorgusu için maksimum sonuç sayısını belirler.
// Negatif değerler otomatik olarak 0'a dönüştürülür.
//
// # Parametreler
//
// - `ctx`: Context nesnesi (timeout, cancellation için)
// - `limit`: Maksimum sonuç sayısı (0 = sınırsız, negatif değerler 0'a çevrilir)
//
// # Döndürür
//
// - `[]interface{}`: Kısıtlanmış sonuç listesi (şu anki implementasyonda boş)
// - `error`: Hata durumunda hata nesnesi
//
// # Örnek
//
// ```go
// constraints := NewRelationshipConstraints(relationshipField)
//
// // İlk 10 kaydı getir
// results, err := constraints.ApplyLimit(ctx, 10)
// if err != nil {
//     log.Fatal(err)
// }
//
// // Limit değerini kontrol et
// fmt.Println(constraints.GetLimit()) // 10
// ```
//
// # Not
//
// - Negatif limit değerleri 0'a dönüştürülür
// - Limit değeri struct içinde saklanır
// - Gerçek implementasyonda bu değer sorguya uygulanır
func (rc *RelationshipConstraintsImpl) ApplyLimit(ctx context.Context, limit int) ([]interface{}, error) {
	if limit < 0 {
		limit = 0
	}

	rc.limit = limit

	// In a real implementation, this would apply the limit to the query
	return []interface{}{}, nil
}

// ApplyOffset, sorgu başlangıç noktasını kaydırır.
//
// Bu metod, relationship sorgusu için kaç kaydın atlanacağını belirler.
// Pagination (sayfalama) işlemleri için kullanılır. Negatif değerler
// otomatik olarak 0'a dönüştürülür.
//
// # Parametreler
//
// - `ctx`: Context nesnesi (timeout, cancellation için)
// - `offset`: Atlanacak kayıt sayısı (negatif değerler 0'a çevrilir)
//
// # Döndürür
//
// - `[]interface{}`: Kaydırılmış sonuç listesi (şu anki implementasyonda boş)
// - `error`: Hata durumunda hata nesnesi
//
// # Örnek
//
// ```go
// constraints := NewRelationshipConstraints(relationshipField)
//
// // Sayfa 2 için (her sayfada 10 kayıt)
// constraints.ApplyLimit(ctx, 10)
// results, err := constraints.ApplyOffset(ctx, 10)
// if err != nil {
//     log.Fatal(err)
// }
//
// // Offset değerini kontrol et
// fmt.Println(constraints.GetOffset()) // 10
// ```
//
// # Not
//
// - Negatif offset değerleri 0'a dönüştürülür
// - Offset değeri struct içinde saklanır
// - Genellikle limit ile birlikte kullanılır (pagination)
// - Gerçek implementasyonda bu değer sorguya uygulanır
func (rc *RelationshipConstraintsImpl) ApplyOffset(ctx context.Context, offset int) ([]interface{}, error) {
	if offset < 0 {
		offset = 0
	}

	rc.offset = offset

	// In a real implementation, this would apply the offset to the query
	return []interface{}{}, nil
}

// ApplyWhere, sorguya WHERE koşulu ekler.
//
// Bu metod, relationship sorgusu için WHERE koşulu ekler. Belirtilen sütun,
// operatör ve değer kullanılarak filtreleme yapılır. Koşul constraints map'inde
// saklanır ve gerçek implementasyonda sorguya uygulanır.
//
// # Parametreler
//
// - `ctx`: Context nesnesi (timeout, cancellation için)
// - `column`: Koşul uygulanacak sütun adı (boş olamaz)
// - `operator`: Karşılaştırma operatörü (=, !=, >, <, >=, <=, LIKE, vb.)
// - `value`: Karşılaştırılacak değer (herhangi bir tip olabilir)
//
// # Döndürür
//
// - `[]interface{}`: Filtrelenmiş sonuç listesi (şu anki implementasyonda boş)
// - `error`: Hata durumunda hata nesnesi
//
// # Örnek
//
// ```go
// constraints := NewRelationshipConstraints(relationshipField)
//
// // Aktif kayıtları getir
// results, err := constraints.ApplyWhere(ctx, "status", "=", "active")
// if err != nil {
//     log.Fatal(err)
// }
//
// // Fiyatı 100'den büyük olanları getir
// results, err = constraints.ApplyWhere(ctx, "price", ">", 100)
//
// // LIKE operatörü ile arama
// results, err = constraints.ApplyWhere(ctx, "name", "LIKE", "%test%")
//
// // Koşulları kontrol et
// constraints := constraints.GetConstraints()
// fmt.Println(constraints["status"]) // map[operator:= value:active]
// ```
//
// # Not
//
// - Boş column değeri durumunda işlem yapılmaz ve boş liste döner
// - Aynı sütun için birden fazla koşul eklenirse son eklenen geçerli olur
// - Koşul map formatında saklanır: map[operator:... value:...]
// - Gerçek implementasyonda bu koşul sorguya uygulanır
//
// # Desteklenen Operatörler
//
// - `=`: Eşittir
// - `!=`: Eşit değildir
// - `>`: Büyüktür
// - `<`: Küçüktür
// - `>=`: Büyük eşittir
// - `<=`: Küçük eşittir
// - `LIKE`: Benzer (pattern matching)
// - `NOT LIKE`: Benzer değil
func (rc *RelationshipConstraintsImpl) ApplyWhere(ctx context.Context, column string, operator string, value interface{}) ([]interface{}, error) {
	if column == "" {
		return []interface{}{}, nil
	}

	rc.constraints[column] = map[string]interface{}{
		"operator": operator,
		"value":    value,
	}

	// In a real implementation, this would apply the WHERE constraint to the query
	return []interface{}{}, nil
}

// ApplyWhereIn, sorguya WHERE IN koşulu ekler.
//
// Bu metod, relationship sorgusu için WHERE IN koşulu ekler. Belirtilen sütunun
// değerinin verilen değerler listesinde olup olmadığını kontrol eder. Koşul
// constraints map'inde saklanır ve gerçek implementasyonda sorguya uygulanır.
//
// # Parametreler
//
// - `ctx`: Context nesnesi (timeout, cancellation için)
// - `column`: Koşul uygulanacak sütun adı (boş olamaz)
// - `values`: Kontrol edilecek değerler listesi (boş olamaz)
//
// # Döndürür
//
// - `[]interface{}`: Filtrelenmiş sonuç listesi (şu anki implementasyonda boş)
// - `error`: Hata durumunda hata nesnesi
//
// # Örnek
//
// ```go
// constraints := NewRelationshipConstraints(relationshipField)
//
// // Belirli ID'lere sahip kayıtları getir
// results, err := constraints.ApplyWhereIn(ctx, "id", []interface{}{1, 2, 3, 5, 8})
// if err != nil {
//     log.Fatal(err)
// }
//
// // Belirli kategorilerdeki kayıtları getir
// results, err = constraints.ApplyWhereIn(ctx, "category", []interface{}{"tech", "science", "art"})
//
// // Belirli durumlardaki kayıtları getir
// results, err = constraints.ApplyWhereIn(ctx, "status", []interface{}{"active", "pending", "approved"})
//
// // Koşulları kontrol et
// constraints := constraints.GetConstraints()
// fmt.Println(constraints["id"]) // map[in:[1 2 3 5 8]]
// ```
//
// # Not
//
// - Boş column veya values durumunda işlem yapılmaz ve boş liste döner
// - Aynı sütun için birden fazla koşul eklenirse son eklenen geçerli olur
// - Koşul map formatında saklanır: map[in:[...]]
// - Gerçek implementasyonda bu koşul sorguya uygulanır
// - Values listesi herhangi bir tip içerebilir (int, string, vb.)
//
// # Performans Notu
//
// WHERE IN sorguları büyük değer listeleri ile yavaş olabilir. Çok sayıda
// değer için alternatif yöntemler düşünülmelidir.
func (rc *RelationshipConstraintsImpl) ApplyWhereIn(ctx context.Context, column string, values []interface{}) ([]interface{}, error) {
	if column == "" || len(values) == 0 {
		return []interface{}{}, nil
	}

	rc.constraints[column] = map[string]interface{}{
		"in": values,
	}

	// In a real implementation, this would apply the WHERE IN constraint to the query
	return []interface{}{}, nil
}

// GetLimit, mevcut limit değerini döndürür.
//
// Bu metod, ApplyLimit ile ayarlanan limit değerini okur.
//
// # Döndürür
//
// - Mevcut limit değeri (0 = sınırsız)
//
// # Örnek
//
// ```go
// constraints := NewRelationshipConstraints(relationshipField)
// constraints.ApplyLimit(ctx, 10)
//
// limit := constraints.GetLimit()
// fmt.Println(limit) // 10
// ```
func (rc *RelationshipConstraintsImpl) GetLimit() int {
	return rc.limit
}

// GetOffset, mevcut offset değerini döndürür.
//
// Bu metod, ApplyOffset ile ayarlanan offset değerini okur.
//
// # Döndürür
//
// - Mevcut offset değeri (0 = baştan başla)
//
// # Örnek
//
// ```go
// constraints := NewRelationshipConstraints(relationshipField)
// constraints.ApplyOffset(ctx, 20)
//
// offset := constraints.GetOffset()
// fmt.Println(offset) // 20
// ```
func (rc *RelationshipConstraintsImpl) GetOffset() int {
	return rc.offset
}

// GetConstraints, tüm WHERE koşullarını döndürür.
//
// Bu metod, ApplyWhere ve ApplyWhereIn ile eklenen tüm koşulları
// map formatında döndürür. Her sütun için koşul bilgisi içerir.
//
// # Döndürür
//
// - WHERE koşulları map'i (sütun adı -> koşul detayları)
//
// # Map Formatı
//
// WHERE koşulu için:
// ```go
// map[string]interface{}{
//     "column_name": map[string]interface{}{
//         "operator": "=",
//         "value": "some_value",
//     },
// }
// ```
//
// WHERE IN koşulu için:
// ```go
// map[string]interface{}{
//     "column_name": map[string]interface{}{
//         "in": []interface{}{1, 2, 3},
//     },
// }
// ```
//
// # Örnek
//
// ```go
// constraints := NewRelationshipConstraints(relationshipField)
// constraints.ApplyWhere(ctx, "status", "=", "active")
// constraints.ApplyWhereIn(ctx, "id", []interface{}{1, 2, 3})
//
// allConstraints := constraints.GetConstraints()
// fmt.Println(allConstraints)
// // Output:
// // map[
// //   status:map[operator:= value:active]
// //   id:map[in:[1 2 3]]
// // ]
//
// // Belirli bir koşulu kontrol et
// if statusConstraint, ok := allConstraints["status"]; ok {
//     constraint := statusConstraint.(map[string]interface{})
//     fmt.Println(constraint["operator"]) // =
//     fmt.Println(constraint["value"])    // active
// }
// ```
//
// # Not
//
// - Döndürülen map referans olarak döner, değişiklikler orijinal map'i etkiler
// - Boş map döndürülmesi hiç koşul eklenmediği anlamına gelir
func (rc *RelationshipConstraintsImpl) GetConstraints() map[string]interface{} {
	return rc.constraints
}
