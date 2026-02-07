package query

// Bu tür, veritabanı sorgularında kullanılan filtre operatörlerini temsil eder.
// FilterOperator, SQL ve NoSQL sorgularında koşul belirtmek için kullanılan
// operatör türlerini tanımlar. String tabanlı bir tür olarak tanımlanmıştır
// ve JSON serileştirmesi için uyumludur.
//
// Kullanım Senaryoları:
// - REST API'lerde filtreleme parametreleri
// - Dinamik sorgu oluşturma
// - İstemci-sunucu iletişiminde operatör aktarımı
// - Veritabanı sorgularının inşası
//
// Örnek:
//   filter := Filter{
//       Field:    "age",
//       Operator: OpGreaterThan,
//       Value:    18,
//   }
type FilterOperator string

const (
	// === EŞİTLİK OPERATÖRLERİ ===
	// OpEqual, alan değerinin belirtilen değere eşit olup olmadığını kontrol eder.
	// SQL: field = value
	// Örnek: Filter{Field: "status", Operator: OpEqual, Value: "active"}
	OpEqual FilterOperator = "eq"

	// OpNotEqual, alan değerinin belirtilen değere eşit olmadığını kontrol eder.
	// SQL: field != value
	// Örnek: Filter{Field: "status", Operator: OpNotEqual, Value: "deleted"}
	OpNotEqual FilterOperator = "neq"

	// === KARŞILAŞTIRMA OPERATÖRLERİ ===
	// OpGreaterThan, alan değerinin belirtilen değerden büyük olup olmadığını kontrol eder.
	// SQL: field > value
	// Örnek: Filter{Field: "price", Operator: OpGreaterThan, Value: 100}
	OpGreaterThan FilterOperator = "gt"

	// OpGreaterEq, alan değerinin belirtilen değerden büyük veya eşit olup olmadığını kontrol eder.
	// SQL: field >= value
	// Örnek: Filter{Field: "score", Operator: OpGreaterEq, Value: 80}
	OpGreaterEq FilterOperator = "gte"

	// OpLessThan, alan değerinin belirtilen değerden küçük olup olmadığını kontrol eder.
	// SQL: field < value
	// Örnek: Filter{Field: "age", Operator: OpLessThan, Value: 65}
	OpLessThan FilterOperator = "lt"

	// OpLessEq, alan değerinin belirtilen değerden küçük veya eşit olup olmadığını kontrol eder.
	// SQL: field <= value
	// Örnek: Filter{Field: "quantity", Operator: OpLessEq, Value: 10}
	OpLessEq FilterOperator = "lte"

	// === METİN OPERATÖRLERİ ===
	// OpLike, alan değerinin belirtilen deseni içerip içermediğini kontrol eder.
	// SQL: field LIKE %value%
	// Büyük/küçük harf duyarlılığı veritabanı ayarlarına bağlıdır.
	// Örnek: Filter{Field: "name", Operator: OpLike, Value: "john"}
	OpLike FilterOperator = "like"

	// OpNotLike, alan değerinin belirtilen deseni içermediğini kontrol eder.
	// SQL: field NOT LIKE %value%
	// Örnek: Filter{Field: "email", Operator: OpNotLike, Value: "spam"}
	OpNotLike FilterOperator = "nlike"

	// === LİSTE OPERATÖRLERİ ===
	// OpIn, alan değerinin belirtilen değerler listesinde olup olmadığını kontrol eder.
	// SQL: field IN (values...)
	// Value parametresi slice veya array olmalıdır.
	// Örnek: Filter{Field: "status", Operator: OpIn, Value: []string{"active", "pending"}}
	OpIn FilterOperator = "in"

	// OpNotIn, alan değerinin belirtilen değerler listesinde olmadığını kontrol eder.
	// SQL: field NOT IN (values...)
	// Value parametresi slice veya array olmalıdır.
	// Örnek: Filter{Field: "role", Operator: OpNotIn, Value: []string{"admin", "superuser"}}
	OpNotIn FilterOperator = "nin"

	// === NULL OPERATÖRLERİ ===
	// OpIsNull, alan değerinin NULL olup olmadığını kontrol eder.
	// SQL: field IS NULL
	// Value parametresi bu operatör için göz ardı edilir.
	// Örnek: Filter{Field: "deleted_at", Operator: OpIsNull}
	OpIsNull FilterOperator = "null"

	// OpIsNotNull, alan değerinin NULL olmadığını kontrol eder.
	// SQL: field IS NOT NULL
	// Value parametresi bu operatör için göz ardı edilir.
	// Örnek: Filter{Field: "email", Operator: OpIsNotNull}
	OpIsNotNull FilterOperator = "nnull"

	// === ARALIK OPERATÖRLERİ ===
	// OpBetween, alan değerinin belirtilen iki değer arasında olup olmadığını kontrol eder.
	// SQL: field BETWEEN value1 AND value2
	// Value parametresi [2]interface{} veya []interface{} (2 eleman) olmalıdır.
	// Örnek: Filter{Field: "age", Operator: OpBetween, Value: []int{18, 65}}
	OpBetween FilterOperator = "between"
)

// Bu yapı, tek bir filtre koşulunu temsil eder.
// Filter, bir veritabanı alanı, bir operatör ve bir değer kombinasyonundan oluşur.
// JSON serileştirmesi için uyumludur ve REST API'lerde kullanılabilir.
//
// Alanlar:
// - Field: Filtrelenecek veritabanı alanının adı (örn: "email", "age", "created_at")
// - Operator: Uygulanacak karşılaştırma operatörü (OpEqual, OpGreaterThan, vb.)
// - Value: Karşılaştırma için kullanılacak değer (interface{} olarak depolanır)
//
// Kullanım Senaryoları:
// - REST API'lerde filtreleme parametreleri
// - Dinamik sorgu oluşturma
// - İstemci tarafından gönderilen filtreleme istekleri
// - Veritabanı sorgularının inşası
//
// Önemli Notlar:
// - Value alanı interface{} türündedir, bu nedenle herhangi bir tür depolanabilir
// - OpIsNull ve OpIsNotNull operatörleri için Value göz ardı edilir
// - OpBetween operatörü için Value iki elemanlı bir slice olmalıdır
// - OpIn ve OpNotIn operatörleri için Value slice veya array olmalıdır
//
// Örnek:
//   filter1 := Filter{
//       Field:    "age",
//       Operator: OpGreaterThan,
//       Value:    18,
//   }
//   filter2 := Filter{
//       Field:    "status",
//       Operator: OpIn,
//       Value:    []string{"active", "pending"},
//   }
type Filter struct {
	// Field, filtrelenecek veritabanı alanının adı
	Field string `json:"field"`

	// Operator, uygulanacak karşılaştırma operatörü
	Operator FilterOperator `json:"operator"`

	// Value, karşılaştırma için kullanılacak değer
	Value interface{} `json:"value"`
}

// Bu yapı, birden fazla filtreyi mantıksal operatörlerle (AND/OR) birleştiren bir grup temsil eder.
// FilterGroup, karmaşık filtreleme koşulları oluşturmak için kullanılır.
// JSON serileştirmesi için uyumludur ve REST API'lerde kullanılabilir.
//
// Alanlar:
// - Logic: Filtreleri birleştirmek için kullanılacak mantıksal operatör ("and" veya "or")
// - Filters: Grupta yer alan Filter nesnelerinin slice'ı
//
// Kullanım Senaryoları:
// - Karmaşık filtreleme koşulları (örn: (age > 18 AND status = 'active') OR role = 'admin')
// - İstemci tarafından gönderilen karmaşık filtreleme istekleri
// - Dinamik sorgu oluşturma
// - Gelişmiş arama ve filtreleme özellikleri
//
// Önemli Notlar:
// - Logic alanı "and" veya "or" olmalıdır (küçük harf)
// - Filters slice'ı boş olabilir, ancak bu durumda hiçbir filtre uygulanmaz
// - İç içe FilterGroup'lar desteklenmeyebilir (uygulamaya bağlı)
// - Tüm filtreler aynı mantıksal operatörle birleştirilir
//
// Örnek:
//   group := FilterGroup{
//       Logic: "and",
//       Filters: []Filter{
//           {Field: "age", Operator: OpGreaterThan, Value: 18},
//           {Field: "status", Operator: OpEqual, Value: "active"},
//       },
//   }
type FilterGroup struct {
	// Logic, filtreleri birleştirmek için kullanılacak mantıksal operatör ("and" veya "or")
	Logic string `json:"logic"`

	// Filters, grupta yer alan Filter nesnelerinin slice'ı
	Filters []Filter `json:"filters"`
}

// Bu değişken, tüm geçerli operatörlerin bir listesini içerir.
// validOperators, operatör doğrulaması için kullanılır ve IsValidOperator
// fonksiyonu tarafından referans alınır.
//
// Önemli Notlar:
// - Bu değişken private (küçük harfle başlayan) olduğu için paket dışından erişilemez
// - ValidOperators() fonksiyonu aracılığıyla public olarak erişilebilir
// - Operatörlerin sırası önemli değildir
var validOperators = []FilterOperator{
	OpEqual, OpNotEqual,
	OpGreaterThan, OpGreaterEq, OpLessThan, OpLessEq,
	OpLike, OpNotLike,
	OpIn, OpNotIn,
	OpIsNull, OpIsNotNull,
	OpBetween,
}

// Bu fonksiyon, tüm geçerli operatörlerin bir listesini döndürür.
// ValidOperators, operatör doğrulaması, UI'de operatör seçenekleri gösterme
// veya operatör listesini almak için kullanılır.
//
// Parametreler: Yok
//
// Dönüş Değeri:
// - []FilterOperator: Tüm geçerli operatörlerin slice'ı
//
// Kullanım Senaryoları:
// - UI'de operatör seçeneklerini gösterme
// - Operatör listesini almak
// - Operatör doğrulaması için referans
// - API dokumentasyonu oluşturma
//
// Önemli Notlar:
// - Döndürülen slice, validOperators değişkeninin doğrudan referansıdır
// - Döndürülen slice'ı değiştirmek validOperators'ı etkileyebilir
// - Güvenlik için döndürülen slice'ın bir kopyasını almayı düşünün
//
// Örnek:
//   operators := ValidOperators()
//   for _, op := range operators {
//       fmt.Println(op)
//   }
func ValidOperators() []FilterOperator {
	return validOperators
}

// Bu fonksiyon, verilen string operatörün geçerli olup olmadığını kontrol eder.
// IsValidOperator, operatör doğrulaması için kullanılır ve operatör
// değerinin kabul edilebilir olup olmadığını belirler.
//
// Parametreler:
// - op (string): Doğrulanacak operatör string'i (örn: "eq", "gt", "like")
//
// Dönüş Değeri:
// - bool: Operatör geçerliyse true, aksi takdirde false
//
// Kullanım Senaryoları:
// - REST API'lerde gelen operatör parametrelerini doğrulama
// - Kullanıcı girdisini kontrol etme
// - Hata ayıklama ve validasyon
// - Operatör değerinin kabul edilebilir olup olmadığını belirme
//
// Önemli Notlar:
// - Operatör karşılaştırması case-sensitive'dir (küçük harf olmalı)
// - Boş string operatörü geçersiz olarak kabul edilir
// - Operatör listesi validOperators değişkeninden alınır
// - Büyük O(n) zaman karmaşıklığına sahiptir (n = operatör sayısı)
//
// Örnek:
//   if IsValidOperator("eq") {
//       fmt.Println("Operatör geçerli")
//   }
//   if !IsValidOperator("invalid") {
//       fmt.Println("Operatör geçersiz")
//   }
func IsValidOperator(op string) bool {
	for _, valid := range validOperators {
		if string(valid) == op {
			return true
		}
	}
	return false
}

// Bu fonksiyon, verilen string'i FilterOperator türüne dönüştürür.
// GetOperator, operatör string'ini geçerli bir FilterOperator'a dönüştürür.
// Eğer operatör geçersizse, varsayılan olarak OpEqual döndürülür.
//
// Parametreler:
// - op (string): Dönüştürülecek operatör string'i (örn: "eq", "gt", "like")
//
// Dönüş Değeri:
// - FilterOperator: Dönüştürülen operatör veya varsayılan OpEqual
//
// Kullanım Senaryoları:
// - REST API'lerde gelen operatör parametrelerini dönüştürme
// - JSON'dan gelen operatör string'lerini işleme
// - Operatör string'lerini FilterOperator türüne dönüştürme
// - Güvenli operatör dönüştürme (geçersiz değerler için fallback)
//
// Önemli Notlar:
// - Geçersiz operatörler için varsayılan olarak OpEqual döndürülür
// - Bu davranış, geçersiz operatörlerin sessizce OpEqual'e dönüştürülmesine neden olabilir
// - Operatör doğrulaması IsValidOperator() fonksiyonu tarafından yapılır
// - Büyük O(n) zaman karmaşıklığına sahiptir (n = operatör sayısı)
//
// Örnek:
//   op1 := GetOperator("gt")      // OpGreaterThan döndürür
//   op2 := GetOperator("invalid") // OpEqual döndürür (varsayılan)
//   op3 := GetOperator("like")    // OpLike döndürür
func GetOperator(op string) FilterOperator {
	if IsValidOperator(op) {
		return FilterOperator(op)
	}
	return OpEqual
}
