package query

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Bu yapı, bir kaynağın sıralama (sort) konfigürasyonunu temsil eder.
//
// Kullanım Senaryosu:
// - Veritabanı sorgularında sıralama düzeni belirtmek için kullanılır
// - Birden fazla sütuna göre sıralama yapılabilir
//
// Örnek:
//
//	Sort{Column: "created_at", Direction: "desc"}
//	Sort{Column: "name", Direction: "asc"}
//
// Önemli Notlar:
// - Direction değeri "asc" (artan) veya "desc" (azalan) olmalıdır
// - Geçersiz direction değerleri otomatik olarak "asc" olarak ayarlanır
type Sort struct {
	Column    string `json:"column"`    // Sıralanacak sütun adı (örn: "id", "created_at", "name")
	Direction string `json:"direction"` // Sıralama yönü: "asc" (artan) veya "desc" (azalan)
}

// Bu yapı, bir kaynağın (resource) sorgu parametrelerini temsil eder.
//
// Kullanım Senaryosu:
// - HTTP GET isteklerinden gelen sorgu parametrelerini parse etmek için kullanılır
// - Arama, filtreleme, sıralama ve sayfalama işlemlerini yönetir
// - Hem yeni nested format hem de eski legacy format'ı destekler
//
// Desteklenen Format Örnekleri:
//
//	Nested Format:
//	  users[search]=john
//	  users[page]=2
//	  users[per_page]=20
//	  users[sort][created_at]=desc
//	  users[filters][status][eq]=active
//
//	Legacy Format:
//	  search=john
//	  page=2
//	  per_page=20
//	  sort_column=created_at
//	  sort_direction=desc
//	  filters[status]=active
//
// Önemli Notlar:
// - Page değeri 1'den başlar (0 geçersizdir)
// - PerPage maksimum 100 olabilir
// - Filters ve Sorts dinamik olarak eklenir
type ResourceQueryParams struct {
	Search  string   // Arama sorgusu (örn: "john" -> tüm aranabilir alanlarda arama yapar)
	Sorts   []Sort   // Sıralama konfigürasyonları (birden fazla sütuna göre sıralama desteklenir)
	Filters []Filter // Filtreleme koşulları (alan, operatör ve değer kombinasyonları)
	Page    int      // Sayfa numarası (1'den başlar, varsayılan: 1)
	PerPage int      // Sayfa başına kayıt sayısı (varsayılan: 10, maksimum: 100)
	View    string   // Index görünümü: "table" (varsayılan) veya "grid"

	// Relationship parametreleri
	ViaResource     string // İlişkili olduğu ana kaynak (örn: "organizations")
	ViaResourceId   string // Ana kaynağın ID'si (örn: "16")
	ViaRelationship string // İlişki adı (örn: "addresses")
}

// Bu fonksiyon, varsayılan sorgu parametrelerini oluşturur ve döndürür.
//
// Kullanım Senaryosu:
// - Yeni bir ResourceQueryParams nesnesi oluşturmak için kullanılır
// - Tüm alanlar varsayılan değerlerle başlatılır
//
// Parametreler: Yok
//
// Dönüş Değeri:
// - *ResourceQueryParams: Varsayılan değerlerle yapılandırılmış pointer'ı
//   - Page: 1 (ilk sayfa)
//   - PerPage: 10 (sayfa başına 10 kayıt)
//   - Filters: Boş slice
//   - Sorts: Boş slice
//   - Search: Boş string
//
// Örnek Kullanım:
//
//	params := DefaultParams()
//	params.Search = "john"
//	params.Page = 2
//
// Önemli Notlar:
// - Her çağrıda yeni bir pointer döndürülür
// - Slice'lar 0 kapasitesiyle başlatılır (dinamik büyüme için)
func DefaultParams() *ResourceQueryParams {
	return &ResourceQueryParams{
		Page:    1,
		PerPage: 10,
		View:    "table",
		Filters: make([]Filter, 0),
		Sorts:   make([]Sort, 0),
	}
}

// Bu fonksiyon, HTTP GET isteklerinden gelen sorgu parametrelerini parse eder.
//
// Kullanım Senaryosu:
// - REST API endpoint'lerinde sorgu parametrelerini işlemek için kullanılır
// - Arama, filtreleme, sıralama ve sayfalama parametrelerini otomatik olarak ayıklar
// - Hem yeni nested format hem de eski legacy format'ı destekler (geriye uyumlu)
//
// Parametreler:
// - c *fiber.Ctx: Fiber HTTP context nesnesi (istek bilgilerini içerir)
// - resourceName string: Kaynağın adı (örn: "users", "products", "orders")
//
// Dönüş Değeri:
// - *ResourceQueryParams: Parse edilmiş sorgu parametrelerini içeren yapılandırılmış pointer'ı
//
// Desteklenen Format Örnekleri:
//
// Nested Format (Yeni - Önerilen):
//
//	GET /api/users?users[search]=john&users[page]=2&users[per_page]=20
//	GET /api/users?users[sort][created_at]=desc&users[sort][name]=asc
//	GET /api/users?users[filters][status][eq]=active&users[filters][age][gt]=18
//
// Legacy Format (Eski - Geriye Uyumlu):
//
//	GET /api/users?search=john&page=2&per_page=20
//	GET /api/users?sort_column=created_at&sort_direction=desc
//	GET /api/users?filters[status]=active
//
// Örnek Kullanım:
//
//	func GetUsers(c *fiber.Ctx) error {
//	    params := ParseResourceQuery(c, "users")
//	    if params.HasSearch() {
//	        // Arama yapılacak
//	    }
//	    if params.HasFilters() {
//	        // Filtreleme yapılacak
//	    }
//	    // Veritabanı sorgusunu oluştur
//	    return c.JSON(users)
//	}
//
// Önemli Notlar:
// - Nested format bulunursa legacy format kontrol edilmez (performans için)
// - Geçersiz sayfa numaraları varsayılan değer (1) olarak ayarlanır
// - PerPage maksimum 100 olabilir (güvenlik için)
// - Direction değerleri otomatik olarak "asc" veya "desc" olarak normalize edilir
func ParseResourceQuery(c *fiber.Ctx, resourceName string) *ResourceQueryParams {
	params := DefaultParams()

	// Ham sorgu dizesini al ve decode et
	rawQuery := string(c.Request().URI().QueryString())

	fmt.Printf("[PARSER] Resource: %s, RawQuery: %s\n", resourceName, rawQuery)

	// Önce nested format'ı dene (decode edilmiş sorgu dizesi kullanarak)
	if parseNestedFormat(rawQuery, resourceName, params) {
		fmt.Printf("[PARSER] Nested format parsed: Search=%q, Page=%d, PerPage=%d\n", params.Search, params.Page, params.PerPage)
		return params
	}

	fmt.Printf("[PARSER] Nested format not found, trying legacy\n")

	// Legacy format'a geri dön (geriye uyumluluk için)
	parseLegacyFormat(c, params)
	return params
}

// Bu fonksiyon, yeni nested format'ı parse eder: resource[key]=value
//
// Kullanım Senaryosu:
// - Modern API'lerde nested sorgu parametrelerini işlemek için kullanılır
// - URL-encoded sorgu dizesini decode ederek parse eder
// - Arama, sayfalama, sıralama ve filtreleme parametrelerini ayıklar
//
// Parametreler:
// - rawQuery string: Ham sorgu dizesi (örn: "users[search]=john&users[page]=2")
// - resource string: Kaynağın adı (örn: "users", "products")
// - params *ResourceQueryParams: Parse edilen parametreleri depolamak için pointer
//
// Dönüş Değeri:
// - bool: Nested format bulunup parse edilirse true, aksi takdirde false
//
// Desteklenen Parametreler:
// - resource[search]=value -> Arama sorgusu
// - resource[page]=number -> Sayfa numarası
// - resource[per_page]=number -> Sayfa başına kayıt sayısı
// - resource[sort][column]=direction -> Sıralama (asc/desc)
// - resource[filters][field][operator]=value -> Filtreleme
//
// Örnek Kullanım:
//
//	rawQuery := "users[search]=john&users[page]=2&users[sort][created_at]=desc"
//	params := DefaultParams()
//	found := parseNestedFormat(rawQuery, "users", params)
//	if found {
//	    fmt.Println("Nested format bulundu:", params.Search)
//	}
//
// Önemli Notlar:
// - URL decode işlemi otomatik olarak yapılır
// - Geçersiz sayfa numaraları göz ardı edilir
// - PerPage maksimum 100 olabilir
// - Direction değerleri otomatik olarak normalize edilir
// - Hata durumunda false döndürülür
func parseNestedFormat(rawQuery string, resource string, params *ResourceQueryParams) bool {
	if rawQuery == "" {
		fmt.Printf("[NESTED] Empty rawQuery\n")
		return false
	}

	// Sorgu dizesini önce URL decode et
	decodedQuery, err := url.QueryUnescape(rawQuery)
	if err != nil {
		decodedQuery = rawQuery
	}

	fmt.Printf("[NESTED] Decoded query: %s\n", decodedQuery)

	// Decode edilmiş sorgu dizesini parse et
	values, err := url.ParseQuery(decodedQuery)
	if err != nil {
		fmt.Printf("[NESTED] Parse error: %v\n", err)
		return false
	}

	fmt.Printf("[NESTED] Parsed values: %+v\n", values)

	found := false
	if resource != "" {
		fmt.Printf("[NESTED] Looking for prefix: %s\n", resource+"[")
	} else {
		fmt.Printf("[NESTED] Resource prefix missing, accepting any nested resource key\n")
	}

	for key, vals := range values {
		fmt.Printf("[NESTED] Key: %s, Vals: %v\n", key, vals)
		if len(vals) == 0 {
			continue
		}
		// Duplicate key durumunda en son gönderilen değeri kullan.
		// Frontend (qs) "last wins" semantiği ile çalışır.
		value := vals[len(vals)-1]

		inner, ok := nestedInnerKeyForResource(key, resource)
		if !ok {
			continue
		}
		found = true

		switch {
		case inner == "search":
			params.Search = value

		case inner == "page":
			if p, err := strconv.Atoi(value); err == nil && p > 0 {
				params.Page = p
			}

		case inner == "per_page":
			if pp, err := strconv.Atoi(value); err == nil && pp > 0 && pp <= 100 {
				params.PerPage = pp
			}

		case inner == "view":
			params.View = normalizeIndexView(value)

		case strings.HasPrefix(inner, "sort]["):
			// sort][name -> name
			column := strings.TrimPrefix(inner, "sort][")
			if column != "" {
				direction := strings.ToLower(value)
				if direction != "asc" && direction != "desc" {
					direction = "asc"
				}
				params.Sorts = append(params.Sorts, Sort{
					Column:    column,
					Direction: direction,
				})
			}

		case strings.HasPrefix(inner, "filters]["):
			parseFilterParam(inner, value, params)

		case inner == "viaResource":
			params.ViaResource = value
		case inner == "viaResourceId":
			params.ViaResourceId = value
		case inner == "viaRelationship":
			params.ViaRelationship = value
		}
	}

	// Root seviyesindeki via parametrelerini de kontrol et (Nested format içinde root parametreler de olabilir)
	// Not: parseNestedFormat sadece prefix ile başlayanları döngüye alıyor, bu yüzden
	// root parametreleri burada ayrıca kontrol etmeliyiz.
	// Ancak values map'inde tüm parametreler var.
	if val := values.Get("viaResource"); val != "" {
		params.ViaResource = val
	}
	if val := values.Get("viaResourceId"); val != "" {
		params.ViaResourceId = val
	}
	if val := values.Get("viaRelationship"); val != "" {
		params.ViaRelationship = val
	}
	if val := values.Get("view"); val != "" {
		params.View = normalizeIndexView(val)
	}

	return found
}

func nestedInnerKeyForResource(key, resource string) (string, bool) {
	if resource != "" {
		prefix := resource + "["
		if !strings.HasPrefix(key, prefix) {
			return "", false
		}

		inner := strings.TrimPrefix(key, prefix)
		inner = strings.TrimSuffix(inner, "]")
		if inner == "" {
			return "", false
		}

		return inner, true
	}

	openBracket := strings.IndexByte(key, '[')
	if openBracket <= 0 || !strings.HasSuffix(key, "]") {
		return "", false
	}

	inner := key[openBracket+1:]
	inner = strings.TrimSuffix(inner, "]")
	if inner == "" {
		return "", false
	}

	return inner, true
}

// Bu fonksiyon, basit ve gelişmiş filtreleme formatlarını işler.
//
// Kullanım Senaryosu:
// - Nested format'tan çıkarılan filtreleme parametrelerini parse etmek için kullanılır
// - Hem basit (varsayılan eq operatörü) hem de gelişmiş (özel operatör) formatları destekler
// - Operatöre göre değeri uygun şekilde parse eder (string, array, boolean vb.)
//
// Parametreler:
// - inner string: Filtreleme parametresinin iç kısmı (örn: "status" veya "status][eq")
// - value string: Filtreleme değeri (örn: "active" veya "active,pending")
// - params *ResourceQueryParams: Parse edilen filtreleri depolamak için pointer
//
// Desteklenen Format Örnekleri:
//
// Basit Format (varsayılan eq operatörü):
//
//	filters][status = "active" -> {field: status, op: eq, value: active}
//	filters][name = "john" -> {field: name, op: eq, value: john}
//
// Gelişmiş Format (özel operatör):
//
//	filters][status][eq = "active" -> {field: status, op: eq, value: active}
//	filters][age][gt = "18" -> {field: age, op: gt, value: 18}
//	filters][status][in = "active,pending" -> {field: status, op: in, value: [active, pending]}
//	filters][created_at][between = "2024-01-01,2024-12-31" -> {field: created_at, op: between, value: [2024-01-01, 2024-12-31]}
//	filters][deleted_at][is_null = "true" -> {field: deleted_at, op: is_null, value: true}
//
// Desteklenen Operatörler:
// - eq (eşit), neq (eşit değil), gt (büyük), gte (büyük eşit)
// - lt (küçük), lte (küçük eşit), like (benzer), nlike (benzer değil)
// - in (içinde), not_in (içinde değil), between (arasında)
// - is_null (boş), is_not_null (boş değil)
//
// Örnek Kullanım:
//
//	parseFilterParam("status", "active", params)
//	parseFilterParam("status][eq", "active", params)
//	parseFilterParam("age][gt", "18", params)
//	parseFilterParam("status][in", "active,pending", params)
//
// Önemli Notlar:
// - Operatör geçersizse varsayılan olarak "eq" kullanılır
// - in/not_in operatörleri virgülle ayrılmış değerleri array'e dönüştürür
// - between operatörü tam olarak 2 değer gerektirir (virgülle ayrılmış)
// - is_null/is_not_null operatörleri boolean değerlere dönüştürülür
// - Geçersiz format'lar göz ardı edilir (return yapılır)
func parseFilterParam(inner, value string, params *ResourceQueryParams) {
	// "filters][" prefix'ini kaldır
	rest := strings.TrimPrefix(inner, "filters][")

	// "][" ile bölerek parçaları al
	parts := strings.Split(rest, "][")

	var field string
	var operator FilterOperator = OpEqual // varsayılan operatör

	if len(parts) == 1 {
		// Basit format: filters][status = "active"
		field = parts[0]
	} else if len(parts) >= 2 {
		// Gelişmiş format: filters][status][eq = "active"
		field = parts[0]
		if IsValidOperator(parts[1]) {
			operator = FilterOperator(parts[1])
		}
	}

	if field == "" {
		return
	}

	// Operatöre göre değeri parse et
	var parsedValue interface{}

	switch operator {
	case OpIn, OpNotIn:
		// Virgülle ayrılmış değerler -> []string
		parsedValue = strings.Split(value, ",")

	case OpBetween:
		// İki virgülle ayrılmış değer -> []string
		betweenParts := strings.Split(value, ",")
		if len(betweenParts) == 2 {
			parsedValue = betweenParts
		} else {
			// Geçersiz between format, atla
			return
		}

	case OpIsNull, OpIsNotNull:
		// Boolean değer
		parsedValue = value == "true" || value == "1"

	default:
		// String değer (eq, neq, gt, gte, lt, lte, like, nlike için)
		parsedValue = value
	}

	params.Filters = append(params.Filters, Filter{
		Field:    field,
		Operator: operator,
		Value:    parsedValue,
	})
}

// Bu fonksiyon, eski flat format'ı parse eder (geriye uyumluluk için).
//
// Kullanım Senaryosu:
// - Eski API'lerde veya legacy sistemlerde sorgu parametrelerini işlemek için kullanılır
// - Nested format bulunmadığında otomatik olarak çağrılır
// - Basit ve düz sorgu parametrelerini destekler
//
// Parametreler:
// - c *fiber.Ctx: Fiber HTTP context nesnesi
// - params *ResourceQueryParams: Parse edilen parametreleri depolamak için pointer
//
// Desteklenen Legacy Parametreler:
// - page: Sayfa numarası (varsayılan: 1)
// - per_page: Sayfa başına kayıt sayısı (varsayılan: 10, maksimum: 100)
// - search: Arama sorgusu
// - sort_column: Sıralanacak sütun adı
// - sort_direction: Sıralama yönü (asc/desc, varsayılan: asc)
// - filters[field]: Filtreleme değeri (basit format, operatör: eq)
//
// Örnek Kullanım:
//
//	GET /api/users?page=2&per_page=20&search=john&sort_column=created_at&sort_direction=desc
//	GET /api/users?filters[status]=active&filters[role]=admin
//
// Önemli Notlar:
// - Geçersiz sayfa numaraları varsayılan değer (1) olarak ayarlanır
// - PerPage maksimum 100 olabilir (güvenlik için)
// - Direction değerleri otomatik olarak normalize edilir
// - Filters basit format'ta (operatör: eq) parse edilir
// - Bu fonksiyon nested format bulunmadığında fallback olarak kullanılır
func parseLegacyFormat(c *fiber.Ctx, params *ResourceQueryParams) {
	// Sayfa
	if p, err := strconv.Atoi(c.Query("page", "1")); err == nil && p > 0 {
		params.Page = p
	}

	// Sayfa başına kayıt sayısı
	if pp, err := strconv.Atoi(c.Query("per_page", "15")); err == nil && pp > 0 && pp <= 100 {
		params.PerPage = pp
	}

	// Arama
	if search := c.Query("search"); search != "" {
		params.Search = search
	}

	// Sıralama (legacy format: sort_column + sort_direction)
	// Compatibility: sort_by + sort_order (frontend lens view)
	sortColumn := c.Query("sort_column")
	if sortColumn == "" {
		sortColumn = c.Query("sort_by")
	}
	if sortColumn != "" {
		sortDirection := c.Query("sort_direction")
		if sortDirection == "" {
			sortDirection = c.Query("sort_order", "asc")
		}
		dir := strings.ToLower(sortDirection)
		if dir != "asc" && dir != "desc" {
			dir = "asc"
		}
		params.Sorts = append(params.Sorts, Sort{
			Column:    sortColumn,
			Direction: dir,
		})
	}

	// Legacy filtreleme format'ı (QueryParser kullanarak map'e dönüştür)
	type LegacyFilters struct {
		Filters map[string]string `query:"filters"`
	}
	lf := new(LegacyFilters)
	if err := c.QueryParser(lf); err == nil && len(lf.Filters) > 0 {
		for field, value := range lf.Filters {
			params.Filters = append(params.Filters, Filter{
				Field:    field,
				Operator: OpEqual,
				Value:    value,
			})
		}
	}

	// Relationship parametreleri
	if viaResource := c.Query("viaResource"); viaResource != "" {
		params.ViaResource = viaResource
	}
	if viaResourceId := c.Query("viaResourceId"); viaResourceId != "" {
		params.ViaResourceId = viaResourceId
	}
	if viaRelationship := c.Query("viaRelationship"); viaRelationship != "" {
		params.ViaRelationship = viaRelationship
	}

	// Görünüm modu
	if view := c.Query("view"); view != "" {
		params.View = normalizeIndexView(view)
	}
}

func normalizeIndexView(raw string) string {
	view := strings.ToLower(strings.TrimSpace(raw))
	if view == "grid" {
		return "grid"
	}
	return "table"
}

// Bu metod, arama sorgusu ayarlanıp ayarlanmadığını kontrol eder.
//
// Kullanım Senaryosu:
// - Arama işleminin yapılması gerekip gerekmediğini belirlemek için kullanılır
// - Koşullu arama mantığı uygulamak için faydalıdır
//
// Parametreler: Yok
//
// Dönüş Değeri:
// - bool: Arama sorgusu boş değilse true, aksi takdirde false
//
// Örnek Kullanım:
//
//	params := ParseResourceQuery(c, "users")
//	if params.HasSearch() {
//	    // Arama yapılacak
//	    results := db.Where("name LIKE ?", "%"+params.Search+"%").Find(&users)
//	} else {
//	    // Arama yapılmayacak
//	    results := db.Find(&users)
//	}
//
// Önemli Notlar:
// - Boş string ("") false döndürür
// - Sadece whitespace içeren string'ler true döndürür (trim yapılmaz)
func (p *ResourceQueryParams) HasSearch() bool {
	return p.Search != ""
}

// Bu metod, sıralama konfigürasyonunun ayarlanıp ayarlanmadığını kontrol eder.
//
// Kullanım Senaryosu:
// - Sıralama işleminin yapılması gerekip gerekmediğini belirlemek için kullanılır
// - Dinamik sıralama mantığı uygulamak için faydalıdır
//
// Parametreler: Yok
//
// Dönüş Değeri:
// - bool: En az bir sıralama konfigürasyonu varsa true, aksi takdirde false
//
// Örnek Kullanım:
//
//	params := ParseResourceQuery(c, "users")
//	if params.HasSorts() {
//	    // Sıralama yapılacak
//	    query := db
//	    for _, sort := range params.Sorts {
//	        query = query.Order(sort.Column + " " + sort.Direction)
//	    }
//	    results := query.Find(&users)
//	} else {
//	    // Varsayılan sıralama
//	    results := db.Order("id DESC").Find(&users)
//	}
//
// Önemli Notlar:
// - Sorts slice'ı boşsa false döndürür
// - Birden fazla sıralama konfigürasyonu olabilir
func (p *ResourceQueryParams) HasSorts() bool {
	return len(p.Sorts) > 0
}

// Bu metod, filtreleme koşulunun ayarlanıp ayarlanmadığını kontrol eder.
//
// Kullanım Senaryosu:
// - Filtreleme işleminin yapılması gerekip gerekmediğini belirlemek için kullanılır
// - Dinamik filtreleme mantığı uygulamak için faydalıdır
//
// Parametreler: Yok
//
// Dönüş Değeri:
// - bool: En az bir filtreleme koşulu varsa true, aksi takdirde false
//
// Örnek Kullanım:
//
//	params := ParseResourceQuery(c, "users")
//	if params.HasFilters() {
//	    // Filtreleme yapılacak
//	    query := db
//	    for _, filter := range params.Filters {
//	        query = applyFilter(query, filter)
//	    }
//	    results := query.Find(&users)
//	} else {
//	    // Filtreleme yapılmayacak
//	    results := db.Find(&users)
//	}
//
// Önemli Notlar:
// - Filters slice'ı boşsa false döndürür
// - Birden fazla filtreleme koşulu olabilir
// - Her filtreleme koşulu farklı operatörleri destekler (eq, gt, in, vb.)
func (p *ResourceQueryParams) HasFilters() bool {
	return len(p.Filters) > 0
}
