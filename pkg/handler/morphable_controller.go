// Package handler, panel API için HTTP istek işleyicilerini sağlar.
//
// Bu dosya, MorphTo ilişki seçeneklerini yönetmek için MorphableController'ı içerir.
// Polimorfik ilişkilerde (polymorphic relationships) dinamik kaynak listelerini
// sağlayarak, kullanıcıların farklı model türlerinden kayıt seçmesine olanak tanır.
//
// # Kullanım Senaryoları
//
// - Bir Comment modelinin hem Post hem de Video'ya ait olabilmesi
// - Bir Image modelinin User, Product veya Article ile ilişkilendirilebilmesi
// - Bir Tag sisteminin farklı içerik türlerine uygulanabilmesi
//
// # Önemli Notlar
//
// - Bu controller, GORM veritabanı bağlantısı gerektirir
// - Kaynak kayıt defteri (resource registry) yapılandırılmış olmalıdır
// - MorphTo alanları önceden tanımlanmış olmalıdır
package handler

import (
	"fmt"
	"strings"

	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// MorphableController, MorphTo ilişki kaynaklarının listelenmesini yönetir.
//
// Bu yapı, polimorfik ilişki alanları için seçilen türe göre seçenekler sağlar.
// Veritabanı bağlantısı ve kaynak kayıt defteri ile çalışarak, dinamik olarak
// farklı model türlerinden kayıtları getirir ve kullanıcıya sunar.
//
// # Kullanım Senaryoları
//
// - Bir yorum sisteminde, yorumun hem blog yazısına hem de videoya ait olabilmesi
// - Bir etiket sisteminde, etiketin farklı içerik türlerine uygulanabilmesi
// - Bir medya yöneticisinde, görselin farklı modellere bağlanabilmesi
//
// # Yapı Alanları
//
// - `DB`: GORM veritabanı bağlantısı - Kaynak sorgularını çalıştırmak için kullanılır
// - `Resources`: Kaynak kayıt defteri - Tüm kayıtlı kaynakların haritası
//
// # Önemli Notlar
//
// - DB bağlantısı nil olmamalıdır
// - Resources haritası başlatılmış olmalıdır
// - Her kaynak türü için uygun tablo adı yapılandırılmış olmalıdır
//
// # Örnek Kullanım
//
//	controller := &MorphableController{
//	    DB:        db,
//	    Resources: resourceRegistry,
//	}
type MorphableController struct {
	DB        *gorm.DB
	Resources map[string]interface{} // Resource registry
}

// NewMorphableController, yeni bir MorphableController örneği oluşturur.
//
// Bu fonksiyon, MorphableController yapısını başlatmak için fabrika fonksiyonu görevi görür.
// Veritabanı bağlantısı ve kaynak kayıt defterini alarak, polimorfik ilişkileri yönetebilecek
// bir controller örneği döndürür.
//
// # Parametreler
//
// - `db`: GORM veritabanı bağlantısı - Kaynak sorgularını çalıştırmak için kullanılır
// - `resources`: Kaynak kayıt defteri - Tüm kayıtlı kaynakların haritası (key: slug, value: resource)
//
// # Dönüş Değeri
//
// Başlatılmış bir `*MorphableController` örneği döndürür.
//
// # Kullanım Senaryoları
//
// - Uygulama başlangıcında controller'ı başlatmak
// - Router'a MorphTo endpoint'lerini eklemek için controller oluşturmak
// - Test senaryolarında mock controller oluşturmak
//
// # Örnek Kullanım
//
//	db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
//	resources := map[string]interface{}{
//	    "posts": &PostResource{},
//	    "videos": &VideoResource{},
//	}
//	controller := NewMorphableController(db, resources)
//
// # Önemli Notlar
//
// - `db` parametresi nil olmamalıdır, aksi takdirde sorgular başarısız olur
// - `resources` haritası boş olabilir ancak başlatılmış olmalıdır
// - Controller oluşturulduktan sonra DB ve Resources alanları değiştirilebilir
//
// # Avantajlar
//
// - Bağımlılıkların açık bir şekilde enjekte edilmesi (Dependency Injection)
// - Test edilebilirlik için mock bağımlılıklar kullanılabilir
// - Tek sorumluluk prensibi (Single Responsibility Principle) uygulanır
func NewMorphableController(db *gorm.DB, resources map[string]interface{}) *MorphableController {
	return &MorphableController{
		DB:        db,
		Resources: resources,
	}
}

// MorphableOption, MorphTo alanı için tek bir seçeneği temsil eder.
//
// Bu yapı, polimorfik ilişki alanlarında kullanıcıya sunulacak seçenekleri tanımlar.
// Her seçenek, bir kaynak kaydını temsil eder ve kullanıcı arayüzünde görüntülenmek
// üzere gerekli bilgileri içerir.
//
// # Kullanım Senaryoları
//
// - Dropdown/select listelerinde seçenekleri göstermek
// - Autocomplete arama sonuçlarını sunmak
// - Mevcut seçili değeri görüntülemek
// - Kullanıcı dostu kaynak seçimi sağlamak
//
// # Yapı Alanları
//
// - `Value`: Kaydın benzersiz tanımlayıcısı (genellikle ID) - Veritabanında saklanacak değer
// - `Display`: Kullanıcıya gösterilecek metin - Kaydın okunabilir adı
// - `Avatar`: Opsiyonel avatar/resim URL'i - Görsel zenginlik için kullanılır
// - `Subtitle`: Opsiyonel alt başlık - Ek bilgi göstermek için kullanılır
//
// # JSON Çıktı Örneği
//
//	{
//	    "value": 123,
//	    "display": "John Doe",
//	    "avatar": "https://example.com/avatar.jpg",
//	    "subtitle": "john@example.com"
//	}
//
// # Önemli Notlar
//
// - `Value` alanı interface{} türündedir, farklı ID türlerini destekler (int, string, uuid)
// - `Avatar` ve `Subtitle` alanları opsiyoneldir, JSON'da boşsa gösterilmez
// - `Display` alanı her zaman doldurulmalıdır, kullanıcı deneyimi için kritiktir
//
// # Avantajlar
//
// - Esnek veri yapısı: Farklı kaynak türlerini destekler
// - Zengin kullanıcı arayüzü: Avatar ve subtitle ile görsel zenginlik
// - JSON serileştirme: API yanıtları için optimize edilmiş
// - Opsiyonel alanlar: Gereksiz veri transferini önler
type MorphableOption struct {
	Value    interface{} `json:"value"`
	Display  string      `json:"display"`
	Avatar   string      `json:"avatar,omitempty"`
	Subtitle string      `json:"subtitle,omitempty"`
}

// HandleMorphable, GET /api/resource/:resource/morphable/:field endpoint'ini işler.
//
// Bu fonksiyon, MorphTo alanı ile ilişkilendirilebilecek kaynakların listesini döndürür.
// Polimorfik ilişkilerde, kullanıcının seçtiği türe göre dinamik olarak kayıtları getirir
// ve kullanıcı arayüzünde gösterilmek üzere formatlar.
//
// # Kullanım Senaryoları
//
// - Dropdown/select listelerinde seçenekleri yüklemek
// - Autocomplete arama yaparken sonuçları getirmek
// - Mevcut seçili değeri yüklemek (edit formlarında)
// - Farklı kaynak türlerinden kayıtları dinamik olarak listelemek
//
// # İstek Parametreleri
//
// ## URL Parametreleri
//
// - `:resource` - MorphTo alanını içeren kaynak slug'ı (örn: "comments")
// - `:field` - MorphTo alan anahtarı (örn: "commentable")
//
// ## Query Parametreleri
//
// - `type` (zorunlu) - Seçenekleri getirmek için kaynak türü/slug'ı (örn: "posts", "videos")
// - `search` (opsiyonel) - Sonuçları filtrelemek için arama sorgusu
// - `per_page` (opsiyonel) - Sayfa başına sonuç sayısı (varsayılan: 10)
// - `current` (opsiyonel) - Mevcut seçili değer ID'si (başlangıç değerini yüklemek için)
//
// # Yanıt Formatı
//
//	{
//	    "resources": [
//	        {
//	            "value": 123,
//	            "display": "John Doe",
//	            "avatar": "https://example.com/avatar.jpg",
//	            "subtitle": "john@example.com"
//	        }
//	    ],
//	    "softDeletes": false
//	}
//
// # Hata Durumları
//
// - `404 Not Found`: MorphTo alanı bulunamadı
// - `400 Bad Request`: type parametresi eksik veya geçersiz
// - `500 Internal Server Error`: Veritabanı sorgusu başarısız oldu
//
// # İşlem Akışı
//
// 1. URL ve query parametrelerini çıkar
// 2. MorphTo alanını bul ve doğrula
// 3. Kaynak türünü doğrula (type parametresi)
// 4. Veritabanından kayıtları sorgula
// 5. Sonuçları formatla ve döndür
//
// # Örnek İstekler
//
// ## Temel Kullanım
//
//	GET /api/resource/comments/morphable/commentable?type=posts&per_page=10
//
// ## Arama ile Kullanım
//
//	GET /api/resource/comments/morphable/commentable?type=users&search=john&per_page=5
//
// ## Mevcut Değeri Yükleme
//
//	GET /api/resource/comments/morphable/commentable?type=posts&current=123
//
// # Önemli Notlar
//
// - `type` parametresi zorunludur, eksikse 400 hatası döner
// - MorphTo alanı önceden tanımlanmış olmalıdır
// - Kaynak türü, MorphTo alanının types listesinde olmalıdır
// - Arama, name, title, email ve username alanlarında yapılır
// - `current` parametresi verilirse, o kayıt sonuçlarda yoksa ayrıca getirilir
// - Soft delete desteği henüz implement edilmemiş (TODO)
//
// # Avantajlar
//
// - Dinamik kaynak yükleme: Farklı türlerden kayıtları destekler
// - Arama desteği: Kullanıcı dostu filtreleme
// - Sayfalama: Performans optimizasyonu
// - Mevcut değer yükleme: Edit formlarında kullanışlı
// - Esnek yapı: Farklı kaynak türlerine uyum sağlar
//
// # Dezavantajlar
//
// - Sabit alan listesi: Sadece name, title, email, username alanlarında arama yapar
// - Soft delete desteği yok: Silinmiş kayıtlar filtrelenmez
// - Özelleştirme sınırlı: Display değeri otomatik belirlenir
//
// # Gereksinimler
//
// - MorphTo alanları için ilgili kaynakları listeleme
// - Veritabanı bağlantısı aktif olmalıdır
// - FieldHandler başlatılmış olmalıdır
func HandleMorphable(h *FieldHandler, c *context.Context) error {
	fieldKey := c.Params("field")
	resourceType := c.Query("type")
	search := c.Query("search", "")
	perPage := c.QueryInt("per_page", 10)
	current := c.Query("current", "")

	// Find the MorphTo field
	var morphToField *fields.MorphTo
	for _, element := range h.getElements(c) {
		if element.GetKey() == fieldKey {
			if mt, ok := element.(*fields.MorphTo); ok {
				morphToField = mt
				break
			}
		}
	}

	if morphToField == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "MorphTo field not found",
		})
	}

	// Validate type is registered in MorphTo types
	if resourceType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "type parameter is required",
		})
	}

	// Check if type is valid
	resourceSlug, err := morphToField.GetResourceForType(resourceType)
	if err != nil {
		// Type might be the slug directly, not the database type
		// Check if it exists in the type mappings as a value (slug)
		found := false
		for _, slug := range morphToField.GetTypes() {
			if slug == resourceType {
				resourceSlug = resourceType
				found = true
				break
			}
		}
		if !found {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Invalid type: %s", resourceType),
			})
		}
	}

	// Query the related resource table - Get GORM DB from provider
	db, ok := h.Provider.GetClient().(*gorm.DB)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database client not available",
		})
	}

	options, err := queryMorphableResources(db, resourceSlug, search, perPage, current)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch resources",
		})
	}

	return c.JSON(fiber.Map{
		"resources":   options,
		"softDeletes": false, // TODO: Implement soft delete detection
	})
}

// queryMorphableResources, veritabanından morphable kaynak seçeneklerini sorgular.
//
// Bu fonksiyon, belirtilen tablo adından kayıtları getirir ve MorphableOption formatına
// dönüştürür. Arama, filtreleme ve sayfalama desteği sağlar. Ayrıca, mevcut seçili
// değerin sonuçlarda olmasını garanti eder.
//
// # Kullanım Senaryoları
//
// - Dropdown/select listelerinde seçenekleri yüklemek
// - Autocomplete arama sonuçlarını getirmek
// - Mevcut seçili değeri yüklemek (edit formlarında)
// - Farklı tablolardan dinamik kayıt listesi oluşturmak
//
// # Parametreler
//
// - `db`: GORM veritabanı bağlantısı - Sorguları çalıştırmak için kullanılır
// - `tableName`: Sorgulanacak tablo adı (örn: "users", "posts", "products")
// - `search`: Arama sorgusu - name, title, email, username alanlarında arama yapar
// - `limit`: Maksimum sonuç sayısı - Sayfalama için kullanılır (0 = sınırsız)
// - `currentID`: Mevcut seçili değer ID'si - Sonuçlarda yoksa ayrıca getirilir
//
// # Dönüş Değeri
//
// - `[]MorphableOption`: Formatlanmış seçenek listesi
// - `error`: Veritabanı hatası veya nil
//
// # İşlem Akışı
//
// 1. DB bağlantısını kontrol et (nil ise boş liste döndür)
// 2. Tablo adından kayıtları sorgula (id, name, title, email, username alanları)
// 3. Arama filtresi uygula (varsa)
// 4. Limit uygula (varsa)
// 5. Sonuçları MorphableOption formatına dönüştür
// 6. Mevcut ID sonuçlarda yoksa ayrıca getir ve başa ekle
//
// # Arama Mantığı
//
// Arama, aşağıdaki alanlarda LIKE operatörü ile yapılır:
// - name: İsim alanı
// - title: Başlık alanı
// - email: E-posta alanı
// - username: Kullanıcı adı alanı
//
// Örnek: search="john" → "name LIKE '%john%' OR email LIKE '%john%'"
//
// # Display Değeri Önceliği
//
// Display değeri şu öncelik sırasına göre belirlenir:
// 1. name (varsa ve boş değilse)
// 2. title (varsa ve boş değilse)
// 3. email (varsa ve boş değilse)
// 4. username (varsa ve boş değilse)
// 5. #ID (hiçbiri yoksa)
//
// # Mevcut Değer Yükleme
//
// `currentID` parametresi verilirse:
// - Önce normal sorgu çalıştırılır
// - Sonuçlarda currentID yoksa, ayrı bir sorgu ile getirilir
// - Mevcut değer, listenin başına eklenir
// - Bu, edit formlarında seçili değerin görünmesini garanti eder
//
// # Örnek Kullanım
//
//	// Temel kullanım
//	options, err := queryMorphableResources(db, "users", "", 10, "")
//
//	// Arama ile kullanım
//	options, err := queryMorphableResources(db, "posts", "golang", 20, "")
//
//	// Mevcut değer ile kullanım
//	options, err := queryMorphableResources(db, "products", "", 10, "123")
//
// # Önemli Notlar
//
// - DB nil ise hata döndürmez, boş liste döndürür
// - Tablo adı doğrudan SQL'e eklenir, SQL injection riski var (güvenilir kaynaklardan gelmeli)
// - Sadece belirli alanlar sorgulanır (id, name, title, email, username)
// - Arama case-insensitive değildir (veritabanı ayarlarına bağlı)
// - Limit 0 ise tüm sonuçlar getirilir (performans riski)
// - currentID bulunamazsa hata döndürmez, sadece eklenmez
//
// # Avantajlar
//
// - Esnek arama: Birden fazla alanda arama yapar
// - Mevcut değer garantisi: Edit formlarında kullanışlı
// - Sayfalama desteği: Performans optimizasyonu
// - Hata toleransı: DB nil ise boş liste döndürür
// - Otomatik display değeri: En uygun alanı seçer
//
// # Dezavantajlar
//
// - Sabit alan listesi: Sadece belirli alanlar sorgulanır
// - SQL injection riski: Tablo adı doğrudan kullanılır
// - Case-sensitive arama: Veritabanı ayarlarına bağlı
// - Özelleştirme sınırlı: Display değeri otomatik belirlenir
// - Performans: Limit olmadan tüm kayıtları getirebilir
//
// # Güvenlik Uyarıları
//
// - `tableName` parametresi güvenilir kaynaklardan gelmeli
// - Kullanıcı girdisi doğrudan tableName'e verilmemeli
// - Arama parametresi GORM tarafından escape edilir (güvenli)
// - currentID parametresi GORM tarafından escape edilir (güvenli)
func queryMorphableResources(db *gorm.DB, tableName, search string, limit int, currentID string) ([]MorphableOption, error) {
	if db == nil {
		return []MorphableOption{}, nil
	}

	selectedColumns, searchableColumns, err := resolveMorphableColumns(db, tableName)
	if err != nil {
		return nil, err
	}

	selectClause := strings.Join(selectedColumns, ", ")

	var results []map[string]interface{}

	// Build query
	query := db.Table(tableName).Select(selectClause)

	// Apply search filter
	if search != "" && len(searchableColumns) > 0 {
		var (
			conditions []string
			args       []interface{}
		)
		likeValue := "%" + search + "%"
		for _, column := range searchableColumns {
			conditions = append(conditions, fmt.Sprintf("%s LIKE ?", column))
			args = append(args, likeValue)
		}
		query = query.Where(strings.Join(conditions, " OR "), args...)
	}

	// Apply limit
	if limit > 0 {
		query = query.Limit(limit)
	}

	// Execute query
	if err := query.Find(&results).Error; err != nil {
		return nil, err
	}

	// Format results
	options := make([]MorphableOption, 0, len(results))
	for _, r := range results {
		display := getDisplayValue(r)
		options = append(options, MorphableOption{
			Value:   r["id"],
			Display: display,
		})
	}

	// If current ID is provided and not in results, fetch it separately
	if currentID != "" {
		found := false
		for _, opt := range options {
			if fmt.Sprint(opt.Value) == currentID {
				found = true
				break
			}
		}

		if !found {
			var currentResult map[string]interface{}
			if err := db.Table(tableName).
				Select(selectClause).
				Where("id = ?", currentID).
				Take(&currentResult).Error; err == nil {
				display := getDisplayValue(currentResult)
				options = append([]MorphableOption{{
					Value:   currentResult["id"],
					Display: display,
				}}, options...)
			}
		}
	}

	return options, nil
}

func resolveMorphableColumns(db *gorm.DB, tableName string) ([]string, []string, error) {
	rows, err := db.Table(tableName).Select("*").Limit(1).Rows()
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	rawColumns, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}

	availableColumns := make(map[string]struct{}, len(rawColumns))
	for _, column := range rawColumns {
		availableColumns[strings.ToLower(strings.TrimSpace(column))] = struct{}{}
	}

	candidateColumns := []string{"id", "name", "title", "email", "username"}
	searchCandidates := []string{"name", "title", "email", "username"}

	selectedColumns := make([]string, 0, len(candidateColumns))
	for _, column := range candidateColumns {
		if _, ok := availableColumns[column]; ok {
			selectedColumns = append(selectedColumns, column)
		}
	}

	hasIDColumn := false
	for _, column := range selectedColumns {
		if column == "id" {
			hasIDColumn = true
			break
		}
	}
	if !hasIDColumn {
		return nil, nil, fmt.Errorf("table %s does not contain an id column", tableName)
	}

	searchableColumns := make([]string, 0, len(searchCandidates))
	for _, column := range searchCandidates {
		if _, ok := availableColumns[column]; ok {
			searchableColumns = append(searchableColumns, column)
		}
	}

	return selectedColumns, searchableColumns, nil
}

// getDisplayValue, bir kaynak kaydından en uygun görüntüleme değerini çıkarır.
//
// Bu fonksiyon, veritabanından gelen kayıt verilerini analiz ederek, kullanıcıya
// gösterilmek üzere en anlamlı değeri seçer. Öncelik sırasına göre farklı alanları
// kontrol eder ve ilk uygun değeri döndürür.
//
// # Kullanım Senaryoları
//
// - Dropdown/select listelerinde kayıt adlarını göstermek
// - Autocomplete sonuçlarında okunabilir metinler sunmak
// - Kullanıcı dostu kayıt tanımlamaları oluşturmak
// - Farklı tablo yapılarına uyum sağlamak
//
// # Parametreler
//
// - `r`: Kayıt verilerini içeren map - Veritabanından gelen ham veri
//
// # Dönüş Değeri
//
// En uygun görüntüleme değerini string olarak döndürür.
// Hiçbir uygun alan bulunamazsa "#ID" formatında ID değerini döndürür.
// ID de yoksa "Unknown" döndürür.
//
// # Öncelik Sırası
//
// Fonksiyon, aşağıdaki öncelik sırasına göre alanları kontrol eder:
//
// 1. **name**: İsim alanı (en yüksek öncelik)
// 2. **title**: Başlık alanı
// 3. **email**: E-posta alanı
// 4. **username**: Kullanıcı adı alanı
// 5. **id**: ID alanı (fallback, "#123" formatında)
// 6. **"Unknown"**: Hiçbir alan yoksa (son çare)
//
// # İşlem Mantığı
//
// Her alan için şu kontroller yapılır:
// - Alan map'te var mı?
// - Alan değeri nil değil mi?
// - String'e çevrildiğinde boş değil mi?
// - String'e çevrildiğinde "<nil>" değil mi?
//
// İlk geçerli alan bulunduğunda döndürülür.
//
// # Örnek Kullanım
//
//	// İsim alanı olan kayıt
//	record1 := map[string]interface{}{
//	    "id": 1,
//	    "name": "John Doe",
//	    "email": "john@example.com",
//	}
//	display1 := getDisplayValue(record1) // "John Doe"
//
//	// Sadece email olan kayıt
//	record2 := map[string]interface{}{
//	    "id": 2,
//	    "email": "jane@example.com",
//	}
//	display2 := getDisplayValue(record2) // "jane@example.com"
//
//	// Sadece ID olan kayıt
//	record3 := map[string]interface{}{
//	    "id": 3,
//	}
//	display3 := getDisplayValue(record3) // "#3"
//
//	// Boş kayıt
//	record4 := map[string]interface{}{}
//	display4 := getDisplayValue(record4) // "Unknown"
//
// # Önemli Notlar
//
// - Fonksiyon, nil değerleri güvenli bir şekilde işler
// - Boş string değerleri atlanır
// - "<nil>" string değeri geçersiz sayılır (GORM'un nil değerleri string'e çevirme davranışı)
// - ID değeri "#" öneki ile formatlanır (örn: "#123")
// - Öncelik sırası değiştirilemez (sabit kodlanmış)
//
// # Avantajlar
//
// - Esnek: Farklı tablo yapılarına uyum sağlar
// - Güvenli: Nil değerleri ve boş stringleri işler
// - Kullanıcı dostu: En anlamlı değeri seçer
// - Fallback mekanizması: Her zaman bir değer döndürür
// - Basit: Karmaşık mantık gerektirmez
//
// # Dezavantajlar
//
// - Sabit öncelik sırası: Özelleştirilemez
// - Sınırlı alan listesi: Sadece 4 alan kontrol edilir
// - Birleştirme yok: Birden fazla alanı birleştirmez (örn: "John Doe (john@example.com)")
// - Dil desteği yok: Çok dilli içerik için özel mantık yok
//
// # Geliştirme Önerileri
//
// - Öncelik sırasını yapılandırılabilir hale getirmek
// - Birden fazla alanı birleştirme desteği eklemek
// - Özel formatlama fonksiyonları desteklemek
// - Dil bazlı alan seçimi eklemek
func getDisplayValue(r map[string]interface{}) string {
	// Priority: name > title > email > username > id
	displayFields := []string{"name", "title", "email", "username"}

	for _, field := range displayFields {
		if val, ok := r[field]; ok && val != nil {
			str := fmt.Sprint(val)
			if str != "" && str != "<nil>" {
				return str
			}
		}
	}

	// Fallback to ID
	if id, ok := r["id"]; ok {
		return fmt.Sprintf("#%v", id)
	}

	return "Unknown"
}
