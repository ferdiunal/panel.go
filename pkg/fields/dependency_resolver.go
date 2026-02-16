// # Alan Bağımlılık Çözücü Paketi
//
// Bu paket, form alanları arasındaki bağımlılıkları yönetir ve çözer.
// Bir alanın değeri değiştiğinde, ona bağımlı diğer alanların otomatik olarak
// güncellenmesini sağlar.
//
// ## Temel Özellikler
//
// - **Bağımlılık Grafiği**: Alanlar arası bağımlılık ilişkilerini graf yapısında tutar
// - **Otomatik Güncelleme**: Değişen alanlara bağlı tüm alanları otomatik günceller
// - **Döngüsel Bağımlılık Tespiti**: Sonsuz döngülere neden olabilecek bağımlılıkları tespit eder
// - **BFS/DFS Algoritmaları**: Etkin graf traversal algoritmaları kullanır
//
// ## Kullanım Senaryoları
//
// 1. **Cascade Seçimler**: Ülke seçildiğinde şehir listesinin güncellenmesi
// 2. **Koşullu Alanlar**: Bir checkbox işaretlendiğinde ilgili alanların gösterilmesi
// 3. **Dinamik Validasyon**: Bir alanın değerine göre diğer alanların validasyon kurallarının değişmesi
// 4. **Hesaplanan Alanlar**: Birden fazla alanın değerine göre otomatik hesaplama yapılması
//
// ## Örnek Kullanım
//
// ```go
// // Bağımlılık çözücü oluştur
// resolver := NewDependencyResolver(fields, "form")
//
// // Döngüsel bağımlılık kontrolü
//
//	if err := resolver.DetectCircularDependencies(); err != nil {
//	    log.Fatal(err)
//	}
//
// // Değişen alanlar için bağımlılıkları çöz
// updates, err := resolver.ResolveDependencies(formData, []string{"country"}, ctx)
//
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// ```
package fields

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

// # DependencyResolver
//
// Alan bağımlılıklarını çözen ve yöneten ana yapı.
//
// ## Amaç
//
// Form alanları arasındaki bağımlılık ilişkilerini yönetir ve bir alan değiştiğinde
// ona bağımlı diğer alanların otomatik olarak güncellenmesini sağlar.
//
// ## Yapı Alanları
//
// - `fields`: Tüm form alanlarının listesi
// - `context`: Bağımlılık çözme bağlamı (örn: "form", "filter", "detail")
//
// ## Çalışma Prensibi
//
// 1. Alanlar arası bağımlılık grafiği oluşturulur
// 2. Değişen alanlar tespit edilir
// 3. BFS algoritması ile etkilenen tüm alanlar bulunur
// 4. Her etkilenen alan için ilgili callback fonksiyonları çalıştırılır
// 5. Güncellenmiş alan bilgileri döndürülür
//
// ## Avantajlar
//
// - **Performans**: Graf yapısı sayesinde O(V+E) karmaşıklığında çalışır
// - **Esneklik**: Context bazlı farklı bağımlılık kuralları tanımlanabilir
// - **Güvenlik**: Döngüsel bağımlılık tespiti ile sonsuz döngüler önlenir
// - **Modülerlik**: Her alan kendi bağımlılık callback'ini tanımlayabilir
//
// ## Dezavantajlar
//
// - Çok karmaşık bağımlılık ağlarında performans düşebilir
// - Callback fonksiyonlarının doğru yazılması gerekir
//
// ## Önemli Notlar
//
// ⚠️ **Uyarı**: Döngüsel bağımlılıklar sonsuz döngüye neden olabilir.
// Mutlaka `DetectCircularDependencies()` ile kontrol edin.
//
// 💡 **İpucu**: Context parametresi ile aynı alanlar farklı bağlamlarda
// farklı davranabilir (örn: form vs filter).
type DependencyResolver struct {
	fields  []*Schema
	context string
}

// # NewDependencyResolver
//
// Yeni bir bağımlılık çözücü oluşturur.
//
// ## Parametreler
//
// - `fields`: Bağımlılık çözümlemesi yapılacak alan listesi
// - `context`: Bağımlılık çözme bağlamı (örn: "form", "filter", "detail")
//
// ## Dönüş Değeri
//
// Yapılandırılmış `*DependencyResolver` örneği döner.
//
// ## Kullanım Örneği
//
// ```go
//
//	fields := []*Schema{
//	    {Key: "country", DependsOnFields: []string{}},
//	    {Key: "city", DependsOnFields: []string{"country"}},
//	    {Key: "district", DependsOnFields: []string{"city"}},
//	}
//
// resolver := NewDependencyResolver(fields, "form")
// ```
//
// ## Önemli Notlar
//
// - Context parametresi, aynı alanların farklı bağlamlarda farklı davranmasını sağlar
// - Oluşturulduktan sonra `DetectCircularDependencies()` ile kontrol yapılması önerilir
func NewDependencyResolver(fields []*Schema, context string) *DependencyResolver {
	return &DependencyResolver{
		fields:  fields,
		context: context,
	}
}

// # ResolveDependencies
//
// Değişen alanlara bağlı tüm alanları tespit eder ve güncelleme bilgilerini döner.
//
// ## Amaç
//
// Form verilerinde değişiklik olan alanları tespit edip, bu alanlara bağımlı olan
// diğer tüm alanları bulur ve her biri için ilgili callback fonksiyonlarını çalıştırarak
// güncellenmiş alan bilgilerini döner.
//
// ## Parametreler
//
// - `formData`: Güncel form verilerini içeren map (alan adı -> değer)
// - `changedFields`: Değişiklik yapılan alan anahtarlarının listesi
// - `ctx`: Fiber context nesnesi (HTTP request/response bilgileri için)
//
// ## Dönüş Değerleri
//
// - `map[string]*FieldUpdate`: Alan anahtarı -> güncelleme bilgisi map'i
// - `error`: Hata durumunda hata mesajı
//
// ## Çalışma Algoritması
//
// 1. **Graf Oluşturma**: Tüm alanlar arası bağımlılık grafiği oluşturulur
// 2. **BFS Traversal**: Değişen alanlardan başlayarak BFS ile etkilenen alanlar bulunur
// 3. **Callback Çalıştırma**: Her etkilenen alan için context'e uygun callback çalıştırılır
// 4. **Sonuç Toplama**: Tüm güncellemeler bir map'te toplanır
//
// ## Kullanım Örneği
//
// ```go
// // Form verisi
//
//	formData := map[string]interface{}{
//	    "country": "TR",
//	    "city": "Istanbul",
//	}
//
// // Değişen alanlar
// changedFields := []string{"country"}
//
// // Bağımlılıkları çöz
// updates, err := resolver.ResolveDependencies(formData, changedFields, ctx)
//
//	if err != nil {
//	    return err
//	}
//
// // Güncellemeleri uygula
//
//	for fieldKey, update := range updates {
//	    fmt.Printf("Alan %s güncellendi: %+v\n", fieldKey, update)
//	}
//
// ```
//
// ## Performans
//
// - **Zaman Karmaşıklığı**: O(V + E) - V: alan sayısı, E: bağımlılık sayısı
// - **Alan Karmaşıklığı**: O(V) - Visited ve affected map'leri için
//
// ## Önemli Notlar
//
// ⚠️ **Uyarı**: Callback fonksiyonları içinde hata oluşursa, bu alan için
// güncelleme döndürülmez ancak diğer alanlar işlenmeye devam eder.
//
// 💡 **İpucu**: Context parametresi sayesinde aynı alan farklı bağlamlarda
// (form, filter, detail) farklı callback'ler kullanabilir.
//
// 📌 **Not**: Döngüsel bağımlılıklar varsa sonsuz döngüye girmez, visited
// map'i sayesinde her alan sadece bir kez işlenir.
func (r *DependencyResolver) ResolveDependencies(
	formData map[string]interface{},
	changedFields []string,
	ctx *fiber.Ctx,
) (map[string]*FieldUpdate, error) {
	updates := make(map[string]*FieldUpdate)

	log.Printf(
		"[depends][resolver] start context=%s changedFields=%v formData=%s fieldCount=%d",
		r.context,
		changedFields,
		toDependencyJSON(formData),
		len(r.fields),
	)

	// Build dependency graph
	dependencyGraph := r.buildDependencyGraph()
	log.Printf("[depends][resolver] dependency-graph context=%s graph=%s", r.context, toDependencyJSON(dependencyGraph))

	// Find affected fields
	affectedFields := r.findAffectedFields(dependencyGraph, changedFields)
	log.Printf("[depends][resolver] affected-fields context=%s changed=%v affected=%v", r.context, changedFields, affectedFields)

	// Execute callbacks for affected fields
	for _, fieldKey := range affectedFields {
		field := r.findFieldByKey(fieldKey)
		if field == nil {
			log.Printf("[depends][resolver] skip-missing-field key=%s", fieldKey)
			continue
		}

		// Get the appropriate callback based on context
		callback := field.GetDependencyCallback(r.context)
		if callback == nil {
			log.Printf(
				"[depends][resolver] skip-no-callback key=%s context=%s dependsOn=%v",
				fieldKey,
				r.context,
				field.DependsOnFields,
			)
			continue
		}

		log.Printf(
			"[depends][resolver] callback-exec key=%s context=%s dependsOn=%v",
			fieldKey,
			r.context,
			field.DependsOnFields,
		)

		// Execute callback
		update := callback(field, formData, ctx)
		if update != nil {
			updates[fieldKey] = update
			log.Printf("[depends][resolver] callback-update key=%s update=%s", fieldKey, toDependencyJSON(update))
			continue
		}

		log.Printf("[depends][resolver] callback-nil key=%s", fieldKey)
	}

	log.Printf("[depends][resolver] done context=%s updates=%s", r.context, toDependencyJSON(updates))

	return updates, nil
}

// # buildDependencyGraph
//
// Alanlar arası bağımlılık ilişkilerini graf yapısında oluşturur.
//
// ## Amaç
//
// Tüm alanları tarayarak her alanın hangi alanlara bağımlı olduğunu tespit eder
// ve bu bilgiyi ters yönde (bağımlı olunan alan -> bağımlı olan alanlar) bir
// map yapısında saklar.
//
// ## Dönüş Değeri
//
// `map[string][]string`: Bağımlılık grafiği
// - **Key**: Bağımlı olunan alan anahtarı
// - **Value**: Bu alana bağımlı olan alanların anahtarları listesi
//
// ## Graf Yapısı
//
// Graf, "ters bağımlılık" mantığıyla çalışır:
// - Eğer Alan B, Alan A'ya bağımlıysa
// - Graf'ta: graph["A"] = ["B"]
// - Bu sayede A değiştiğinde B'nin etkilendiği kolayca bulunur
//
// ## Kullanım Örneği
//
// ```go
// // Alanlar:
// // - country (bağımlılık yok)
// // - city (country'ye bağımlı)
// // - district (city'ye bağımlı)
//
// graph := resolver.buildDependencyGraph()
// // Sonuç:
// // {
// //   "country": ["city"],
// //   "city": ["district"]
// // }
//
// // country değiştiğinde city'nin etkilendiğini bul
// affectedByCoutry := graph["country"] // ["city"]
// ```
//
// ## Algoritma
//
// 1. Boş bir graf map'i oluştur
// 2. Her alan için:
//   - Eğer bağımlılığı yoksa atla
//   - Her bağımlılık için:
//   - Graf'ta bağımlı olunan alanı key olarak ekle
//   - Bu key'in value listesine mevcut alanı ekle
//
// ## Performans
//
// - **Zaman Karmaşıklığı**: O(F × D) - F: alan sayısı, D: ortalama bağımlılık sayısı
// - **Alan Karmaşıklığı**: O(E) - E: toplam bağımlılık sayısı
//
// ## Önemli Notlar
//
// 📌 **Not**: Bu fonksiyon private'dır ve sadece ResolveDependencies içinde kullanılır.
//
// 💡 **İpucu**: Graf yapısı sayesinde BFS/DFS algoritmaları ile etkin traversal yapılabilir.
//
// ⚠️ **Uyarı**: Döngüsel bağımlılıklar varsa graf sonsuz döngü içerebilir.
// DetectCircularDependencies() ile kontrol yapılmalıdır.
func (r *DependencyResolver) buildDependencyGraph() map[string][]string {
	graph := make(map[string][]string)

	for _, field := range r.fields {
		if len(field.DependsOnFields) == 0 {
			continue
		}

		for _, dependsOn := range field.DependsOnFields {
			if graph[dependsOn] == nil {
				graph[dependsOn] = []string{}
			}
			graph[dependsOn] = append(graph[dependsOn], field.Key)
		}
	}

	log.Printf("[depends][resolver] graph-built context=%s graph=%s", r.context, toDependencyJSON(graph))

	return graph
}

// # findAffectedFields
//
// Değişen alanlara bağımlı olan tüm alanları BFS (Breadth-First Search) algoritması ile bulur.
//
// ## Amaç
//
// Bir veya birden fazla alan değiştiğinde, bu değişiklikten doğrudan veya dolaylı olarak
// etkilenen tüm alanları tespit eder. Bağımlılık zincirini takip ederek cascade etkiyi
// hesaplar.
//
// ## Parametreler
//
// - `graph`: Bağımlılık grafiği (buildDependencyGraph tarafından oluşturulur)
// - `changedFields`: Değişiklik yapılan alan anahtarlarının listesi
//
// ## Dönüş Değeri
//
// `[]string`: Etkilenen alan anahtarlarının listesi (değişen alanlar hariç)
//
// ## Algoritma: BFS (Breadth-First Search)
//
// 1. **Başlangıç**: Değişen alanları kuyruğa ekle
// 2. **Traversal**: Kuyruktan alan çıkar, ziyaret edildi olarak işaretle
// 3. **Bağımlıları Bul**: Graf'tan bu alana bağımlı alanları bul
// 4. **Etkilenenleri Kaydet**: Bağımlı alanları "etkilenen" olarak işaretle
// 5. **Kuyruğa Ekle**: Henüz ziyaret edilmemiş bağımlıları kuyruğa ekle
// 6. **Tekrarla**: Kuyruk boşalana kadar devam et
//
// ## Kullanım Örneği
//
// ```go
// // Graf yapısı:
// // country -> city -> district
// //         -> state
//
//	graph := map[string][]string{
//	    "country": {"city", "state"},
//	    "city": {"district"},
//	}
//
// // country değiştiğinde etkilenen alanlar
// affected := resolver.findAffectedFields(graph, []string{"country"})
// // Sonuç: ["city", "state", "district"]
//
// // city değiştiğinde etkilenen alanlar
// affected = resolver.findAffectedFields(graph, []string{"city"})
// // Sonuç: ["district"]
// ```
//
// ## Performans
//
// - **Zaman Karmaşıklığı**: O(V + E)
//   - V: Graf'taki toplam alan sayısı
//   - E: Graf'taki toplam bağımlılık sayısı
//
// - **Alan Karmaşıklığı**: O(V)
//   - affected, visited ve queue için
//
// ## Cascade Etki Örneği
//
// ```
// Değişiklik: country = "TR"
//
//	↓
//
// Etkilenen: city (İstanbul, Ankara, İzmir seçenekleri yüklenir)
//
//	↓
//
// Etkilenen: district (city'ye göre ilçeler yüklenir)
//
//	↓
//
// Etkilenen: neighborhood (district'e göre mahalleler yüklenir)
// ```
//
// ## Önemli Notlar
//
// 📌 **Not**: Bu fonksiyon private'dır ve sadece ResolveDependencies içinde kullanılır.
//
// 💡 **İpucu**: BFS algoritması sayesinde aynı seviyedeki tüm alanlar önce işlenir,
// sonra bir sonraki seviyeye geçilir (level-order traversal).
//
// ⚠️ **Uyarı**: Döngüsel bağımlılıklar varsa visited map'i sayesinde sonsuz döngüye
// girmez, ancak bağımlılık sırası beklenmedik olabilir.
//
// 🔍 **Detay**: Değişen alanların kendileri sonuç listesine dahil edilmez, sadece
// onlara bağımlı olan alanlar döndürülür.
func (r *DependencyResolver) findAffectedFields(
	graph map[string][]string,
	changedFields []string,
) []string {
	affected := make(map[string]bool)
	visited := make(map[string]bool)

	// BFS to find all affected fields
	queue := make([]string, len(changedFields))
	copy(queue, changedFields)
	log.Printf("[depends][resolver] bfs-start changed=%v", changedFields)

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		log.Printf("[depends][resolver] bfs-pop current=%s pending=%v", current, queue)

		if visited[current] {
			log.Printf("[depends][resolver] bfs-skip-visited current=%s", current)
			continue
		}
		visited[current] = true

		// Get fields that depend on current field
		dependents := graph[current]
		log.Printf("[depends][resolver] bfs-dependents current=%s dependents=%v", current, dependents)
		for _, dependent := range dependents {
			affected[dependent] = true
			log.Printf("[depends][resolver] bfs-affected dependent=%s by=%s", dependent, current)

			// Check for circular dependencies
			if !visited[dependent] {
				queue = append(queue, dependent)
				log.Printf("[depends][resolver] bfs-enqueue dependent=%s queue=%v", dependent, queue)
			}
		}
	}

	// Convert map to slice
	result := make([]string, 0, len(affected))
	for field := range affected {
		result = append(result, field)
	}

	log.Printf("[depends][resolver] bfs-done changed=%v affected=%v", changedFields, result)

	return result
}

// # findFieldByKey
//
// Verilen anahtar değerine sahip alanı bulur ve döndürür.
//
// ## Amaç
//
// Alan listesinde doğrusal arama yaparak belirtilen key değerine sahip
// Schema nesnesini bulur. Bu fonksiyon, bağımlılık çözme sürecinde
// alan anahtarından alan nesnesine erişim için kullanılır.
//
// ## Parametreler
//
// - `key`: Aranacak alan anahtarı (string)
//
// ## Dönüş Değeri
//
// - `*Schema`: Bulunan alan nesnesi, bulunamazsa `nil`
//
// ## Kullanım Örneği
//
// ```go
// // Alan listesi
//
//	fields := []*Schema{
//	    {Key: "country", Label: "Ülke"},
//	    {Key: "city", Label: "Şehir"},
//	    {Key: "district", Label: "İlçe"},
//	}
//
// resolver := NewDependencyResolver(fields, "form")
//
// // Alan bul
// cityField := resolver.findFieldByKey("city")
//
//	if cityField != nil {
//	    fmt.Println(cityField.Label) // "Şehir"
//	}
//
// // Olmayan alan
// unknownField := resolver.findFieldByKey("unknown")
//
//	if unknownField == nil {
//	    fmt.Println("Alan bulunamadı")
//	}
//
// ```
//
// ## Performans
//
// - **Zaman Karmaşıklığı**: O(n) - n: toplam alan sayısı
// - **Alan Karmaşıklığı**: O(1) - Sabit bellek kullanımı
//
// ## Optimizasyon Önerileri
//
// Eğer bu fonksiyon sık çağrılıyorsa ve performans kritikse:
// 1. Alan listesini map[string]*Schema yapısında tutmak (O(1) erişim)
// 2. Lazy initialization ile ilk çağrıda map oluşturmak
// 3. Cache mekanizması eklemek
//
// ```go
// // Örnek optimizasyon
//
//	type DependencyResolver struct {
//	    fields    []*Schema
//	    fieldMap  map[string]*Schema // Cache
//	    context   string
//	}
//
//	func (r *DependencyResolver) findFieldByKey(key string) *Schema {
//	    if r.fieldMap == nil {
//	        r.fieldMap = make(map[string]*Schema)
//	        for _, field := range r.fields {
//	            r.fieldMap[field.Key] = field
//	        }
//	    }
//	    return r.fieldMap[key]
//	}
//
// ```
//
// ## Önemli Notlar
//
// Bu fonksiyon private'dır ve sadece DependencyResolver içinde kullanılır.
//
// Alan bulunamazsa nil döner, bu durumun kontrol edilmesi gerekir.
//
// Key değerleri case-sensitive'dir, "City" ve "city" farklı kabul edilir.
func (r *DependencyResolver) findFieldByKey(key string) *Schema {
	for _, field := range r.fields {
		if field.Key == key {
			return field
		}
	}
	return nil
}

// # DetectCircularDependencies
//
// Alan bağımlılıkları arasında döngüsel (circular) bağımlılık olup olmadığını tespit eder.
//
// ## Amaç
//
// Bağımlılık grafiğinde döngü (cycle) olup olmadığını kontrol eder. Döngüsel bağımlılıklar
// sonsuz döngülere ve stack overflow hatalarına neden olabileceği için, sistem başlatılırken
// veya alan tanımları değiştiğinde bu kontrolün yapılması kritik önem taşır.
//
// ## Dönüş Değeri
//
// - `nil`: Döngüsel bağımlılık yok, sistem güvenli
// - `error`: Döngüsel bağımlılık tespit edildi, hata mesajında ilgili alan belirtilir
//
// ## Döngüsel Bağımlılık Nedir?
//
// Döngüsel bağımlılık, alanların birbirine doğrudan veya dolaylı olarak bağımlı olduğu
// ve bir döngü oluşturduğu durumdur.
//
// ### Örnekler
//
// **Doğrudan Döngü:**
// ```
// Alan A -> Alan B -> Alan A
// ```
//
// **Dolaylı Döngü:**
// ```
// Alan A -> Alan B -> Alan C -> Alan A
// ```
//
// **Karmaşık Döngü:**
// ```
// Alan A -> Alan B -> Alan C
//
//	  ↓         ↓
//	Alan D -> Alan E -> Alan A
//
// ```
//
// ## Kullanım Örneği
//
// ```go
// // Alan tanımları
//
//	fields := []*Schema{
//	    {Key: "country", DependsOnFields: []string{"city"}},  // Hatalı!
//	    {Key: "city", DependsOnFields: []string{"country"}},  // Döngü!
//	}
//
// resolver := NewDependencyResolver(fields, "form")
//
// // Döngüsel bağımlılık kontrolü
//
//	if err := resolver.DetectCircularDependencies(); err != nil {
//	    log.Fatal(err) // "circular dependency detected involving field: country"
//	}
//
// ```
//
// ## Doğru Kullanım
//
// ```go
// // Doğru alan tanımları (tek yönlü bağımlılık)
//
//	fields := []*Schema{
//	    {Key: "country", DependsOnFields: []string{}},
//	    {Key: "city", DependsOnFields: []string{"country"}},
//	    {Key: "district", DependsOnFields: []string{"city"}},
//	}
//
// resolver := NewDependencyResolver(fields, "form")
//
// // Kontrol başarılı
//
//	if err := resolver.DetectCircularDependencies(); err != nil {
//	    log.Fatal(err)
//	}
//
// // Hata yok, sistem güvenli
// ```
//
// ## Algoritma: DFS (Depth-First Search)
//
// 1. Her alan için DFS başlat (henüz ziyaret edilmemişse)
// 2. Alanı ziyaret edildi olarak işaretle
// 3. Alanı recursion stack'e ekle
// 4. Alanın bağımlılarını kontrol et:
//   - Bağımlı henüz ziyaret edilmemişse, recursive DFS çağrısı yap
//   - Bağımlı recursion stack'te varsa, döngü tespit edildi
//
// 5. Alanı recursion stack'ten çıkar
//
// ## Performans
//
// - **Zaman Karmaşıklığı**: O(V + E)
//   - V: Toplam alan sayısı
//   - E: Toplam bağımlılık sayısı
//
// - **Alan Karmaşıklığı**: O(V)
//   - visited ve recStack map'leri için
//
// ## Ne Zaman Çağrılmalı?
//
// 1. **Sistem Başlatma**: Uygulama başlarken tüm alanlar için kontrol
// 2. **Alan Tanımı Değişikliği**: Yeni alan eklendiğinde veya bağımlılık değiştiğinde
// 3. **Geliştirme Aşaması**: Unit testlerde otomatik kontrol
// 4. **Deployment Öncesi**: CI/CD pipeline'da validasyon
//
// ## Önemli Notlar
//
// **UYARI**: Bu fonksiyon mutlaka çağrılmalıdır. Döngüsel bağımlılıklar runtime'da
// sonsuz döngüye ve sistem çökmesine neden olabilir.
//
// **NOT**: Hata mesajı sadece döngüye dahil olan alanlardan birini gösterir.
// Tüm döngüyü görmek için ek analiz gerekebilir.
//
// **İPUCU**: Geliştirme ortamında panic kullanarak erken tespit yapılabilir:
// ```go
//
//	if err := resolver.DetectCircularDependencies(); err != nil {
//	    panic(err) // Geliştirme ortamında hemen fark edilir
//	}
//
// ```
//
// ## Test Örneği
//
// ```go
//
//	func TestCircularDependency(t *testing.T) {
//	    fields := []*Schema{
//	        {Key: "a", DependsOnFields: []string{"b"}},
//	        {Key: "b", DependsOnFields: []string{"c"}},
//	        {Key: "c", DependsOnFields: []string{"a"}}, // Döngü!
//	    }
//
//	    resolver := NewDependencyResolver(fields, "test")
//	    err := resolver.DetectCircularDependencies()
//
//	    assert.Error(t, err)
//	    assert.Contains(t, err.Error(), "circular dependency")
//	}
//
// ```
func (r *DependencyResolver) DetectCircularDependencies() error {
	graph := r.buildDependencyGraph()
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	log.Printf("[depends][resolver] circular-check-start context=%s fieldCount=%d", r.context, len(r.fields))

	for _, field := range r.fields {
		if !visited[field.Key] {
			if r.hasCycle(field.Key, graph, visited, recStack) {
				log.Printf("[depends][resolver] circular-check-failed context=%s field=%s", r.context, field.Key)
				return fmt.Errorf("circular dependency detected involving field: %s", field.Key)
			}
		}
	}

	log.Printf("[depends][resolver] circular-check-ok context=%s", r.context)

	return nil
}

func toDependencyJSON(v interface{}) string {
	payload, err := json.Marshal(v)
	if err != nil {
		return "<marshal_error>"
	}
	return string(payload)
}

// # hasCycle
//
// Belirtilen alandan başlayarak DFS (Depth-First Search) algoritması ile döngü olup olmadığını kontrol eder.
//
// ## Amaç
//
// Recursive DFS algoritması kullanarak bağımlılık grafiğinde döngü tespit eder.
// Bu fonksiyon DetectCircularDependencies tarafından her alan için çağrılır.
//
// ## Parametreler
//
// - `fieldKey`: Kontrol edilecek alan anahtarı
// - `graph`: Bağımlılık grafiği (alan -> bağımlılar map'i)
// - `visited`: Ziyaret edilen alanları tutan map
// - `recStack`: Recursion stack'te olan alanları tutan map (döngü tespiti için)
//
// ## Dönüş Değeri
//
// - `true`: Döngü tespit edildi
// - `false`: Döngü yok
//
// ## Algoritma: DFS ile Döngü Tespiti
//
// DFS algoritması iki map kullanır:
//
// 1. **visited**: Bir alanın daha önce ziyaret edilip edilmediğini tutar
//   - Gereksiz tekrar ziyaretleri önler
//   - Performans optimizasyonu sağlar
//
// 2. **recStack** (Recursion Stack): Mevcut DFS yolunda hangi alanların olduğunu tutar
//   - Eğer bir alan hem ziyaret edilmişse hem de recStack'te varsa, döngü var demektir
//   - Her DFS dalı tamamlandığında alan recStack'ten çıkarılır
//
// ## Adım Adım Çalışma
//
// ```
//  1. Alanı visited ve recStack'e ekle
//  2. Alanın tüm bağımlılarını kontrol et:
//     a. Bağımlı henüz ziyaret edilmemişse:
//     - Recursive olarak hasCycle çağır
//     - Eğer döngü bulunursa true döndür
//     b. Bağımlı recStack'te varsa:
//     - Döngü tespit edildi, true döndür
//  3. Alanı recStack'ten çıkar (backtrack)
//  4. Döngü bulunamadı, false döndür
//
// ```
//
// ## Görsel Örnek
//
// **Döngü Var:**
// ```
// A -> B -> C -> A
//
// hasCycle("A"):
//
//	visited: {A}, recStack: {A}
//	hasCycle("B"):
//	  visited: {A,B}, recStack: {A,B}
//	  hasCycle("C"):
//	    visited: {A,B,C}, recStack: {A,B,C}
//	    hasCycle("A"):
//	      A visited=true ve recStack=true
//	      DÖNGÜ TESPİT EDİLDİ! -> return true
//
// ```
//
// **Döngü Yok:**
// ```
// A -> B -> C
// A -> D
//
// hasCycle("A"):
//
//	visited: {A}, recStack: {A}
//	hasCycle("B"):
//	  visited: {A,B}, recStack: {A,B}
//	  hasCycle("C"):
//	    visited: {A,B,C}, recStack: {A,B,C}
//	    C'nin bağımlısı yok
//	    recStack: {A,B} (C çıkarıldı)
//	  recStack: {A} (B çıkarıldı)
//	hasCycle("D"):
//	  visited: {A,B,C,D}, recStack: {A,D}
//	  D'nin bağımlısı yok
//	  recStack: {A} (D çıkarıldı)
//	recStack: {} (A çıkarıldı)
//
// DÖNGÜ YOK -> return false
// ```
//
// ## Performans
//
// - **Zaman Karmaşıklığı**: O(V + E)
//   - V: Toplam alan sayısı
//   - E: Toplam bağımlılık sayısı
//   - Her alan ve her bağımlılık en fazla bir kez ziyaret edilir
//
// - **Alan Karmaşıklığı**: O(V)
//   - Recursion stack derinliği en fazla V olabilir
//   - visited ve recStack map'leri O(V) alan kullanır
//
// ## Neden İki Map?
//
// **visited** olmadan:
// - Aynı alanlar tekrar tekrar ziyaret edilir
// - Performans O(V!) gibi çok kötü olur
//
// **recStack** olmadan:
// - Farklı dallardan gelen ziyaretler döngü olarak algılanır
// - Yanlış pozitif sonuçlar üretilir
//
// Örnek:
// ```
// A -> B -> D
// A -> C -> D
// ```
// D iki farklı yoldan ziyaret edilir ama döngü yoktur.
// recStack sayesinde bu durum doğru tespit edilir.
//
// ## Önemli Notlar
//
// **NOT**: Bu fonksiyon private'dır ve sadece DetectCircularDependencies tarafından çağrılır.
//
// **UYARI**: Recursive fonksiyondur, çok derin bağımlılık ağlarında stack overflow
// riski vardır (ancak normal kullanımda bu durum çok nadirdir).
//
// **İPUCU**: recStack'in backtrack edilmesi (fonksiyon sonunda false yapılması)
// kritik öneme sahiptir. Aksi halde yanlış pozitif sonuçlar üretilir.
//
// ## Test Senaryoları
//
// ```go
// // Test 1: Basit döngü
//
//	graph := map[string][]string{
//	    "a": {"b"},
//	    "b": {"a"},
//	}
//
// // hasCycle("a") -> true
//
// // Test 2: Dolaylı döngü
//
//	graph := map[string][]string{
//	    "a": {"b"},
//	    "b": {"c"},
//	    "c": {"a"},
//	}
//
// // hasCycle("a") -> true
//
// // Test 3: Döngü yok
//
//	graph := map[string][]string{
//	    "a": {"b", "c"},
//	    "b": {"d"},
//	    "c": {"d"},
//	}
//
// // hasCycle("a") -> false
//
// // Test 4: Karmaşık graf, döngü yok
//
//	graph := map[string][]string{
//	    "a": {"b", "c"},
//	    "b": {"d"},
//	    "c": {"d", "e"},
//	    "d": {"f"},
//	    "e": {"f"},
//	}
//
// // hasCycle("a") -> false
// ```
func (r *DependencyResolver) hasCycle(
	fieldKey string,
	graph map[string][]string,
	visited map[string]bool,
	recStack map[string]bool,
) bool {
	visited[fieldKey] = true
	recStack[fieldKey] = true

	// Get all fields that depend on this field
	dependents := graph[fieldKey]
	for _, dependent := range dependents {
		if !visited[dependent] {
			if r.hasCycle(dependent, graph, visited, recStack) {
				return true
			}
		} else if recStack[dependent] {
			return true
		}
	}

	recStack[fieldKey] = false
	return false
}
