package handler

import (
	"fmt"
	"sync"

	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/widget"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// / # cardResult
// /
// / Bu yapı, bir kart widget'ının asenkron çözümleme işleminin sonucunu temsil eder.
// / Paralel işleme sırasında her bir kartın durumunu ve verisini takip etmek için kullanılır.
// /
// / ## Alanlar
// /
// / - `card`: Çözümlenen widget.Card arayüzü implementasyonu
// / - `data`: Kartın çözümlenmiş veri içeriği (herhangi bir tip olabilir)
// / - `err`: Çözümleme sırasında oluşan hata (varsa)
// / - `index`: Kartın orijinal sıradaki konumu (sıralama için kritik)
// / - `serialized`: JSON serileştirilmiş kart özellikleri
// /
// / ## Kullanım Senaryosu
// /
// / Bu yapı, fan-out/fan-in concurrency pattern'inde kullanılır:
// / 1. Her kart için bir goroutine başlatılır
// / 2. Her goroutine kendi cardResult'ını oluşturur
// / 3. Sonuçlar channel üzerinden toplanır
// / 4. index alanı sayesinde orijinal sıralama korunur
// /
// / ## Önemli Notlar
// /
// / - Bu yapı sadece internal kullanım içindir (küçük harfle başlar)
// / - Thread-safe değildir, channel üzerinden iletilmek için tasarlanmıştır
// / - index alanı, paralel işleme sonrası sıralamayı korumak için kritiktir
// /
// / ## Örnek Kullanım
// /
// / ```go
// / result := cardResult{
// /     card:       myCard,
// /     data:       resolvedData,
// /     err:        nil,
// /     index:      0,
// /     serialized: map[string]interface{}{"name": "MyCard"},
// / }
// / results <- result // Channel'a gönder
// / ```
type cardResult struct {
	card       widget.Card
	data       interface{}
	err        error
	index      int
	serialized map[string]interface{}
}

// / # HandleCardList
// /
// / Bu fonksiyon, bir kaynak için tüm kart widget'larını listeler ve her kartın verisini
// / paralel olarak çözümler. Asenkron fan-out/fan-in pattern kullanarak yüksek performans sağlar.
// /
// / ## Parametreler
// /
// / - `h *FieldHandler`: Kart listesini ve veritabanı bağlantısını içeren handler
// /   - `h.Cards`: Çözümlenecek widget.Card slice'ı
// /   - `h.DB`: Veritabanı bağlantısı (kartların veri çözümlemesi için)
// / - `c *context.Context`: Fiber context wrapper'ı, HTTP isteği ve yanıtı için
// /
// / ## Dönüş Değeri
// /
// / - `error`: JSON yanıtı gönderme hatası veya nil
// /
// / ## Çalışma Prensibi
// /
// / ### 1. Boş Kontrol
// / Eğer kart listesi boşsa, boş bir array döner.
// /
// / ### 2. Fan-Out (Dağıtım)
// / - Her kart için ayrı bir goroutine başlatılır
// / - Buffered channel kullanılarak non-blocking send sağlanır
// / - WaitGroup ile goroutine'lerin tamamlanması takip edilir
// /
// / ### 3. Paralel Çözümleme
// / Her goroutine:
// / - Kartın temel özelliklerini serileştirir (name, component, width)
// / - `w.Resolve(c, h.DB)` ile kartın verisini çözümler
// / - Sonucu channel'a gönderir
// /
// / ### 4. Fan-In (Toplama)
// / - Ayrı bir goroutine channel'ı kapatır (tüm işler bitince)
// / - Ana goroutine channel'dan sonuçları toplar
// / - Sonuçlar orijinal index'lerine göre sıralanır
// /
// / ### 5. Hata Yönetimi
// / - Her kartın hatası bağımsız olarak ele alınır
// / - Hata durumunda kart atlanmaz, hata mesajı eklenir
// / - Diğer kartların çözümlenmesi devam eder
// /
// / ## Performans Özellikleri
// /
// / ### Avantajlar
// / - **Paralel İşleme**: N kart için O(1) zaman (en yavaş kartın süresi)
// / - **Non-Blocking**: Buffered channel sayesinde goroutine'ler beklemez
// / - **Ölçeklenebilir**: Kart sayısı arttıkça performans avantajı artar
// / - **Hata İzolasyonu**: Bir kartın hatası diğerlerini etkilemez
// /
// / ### Dezavantajlar
// / - **Bellek Kullanımı**: Her kart için bir goroutine = N goroutine
// / - **Goroutine Overhead**: Az sayıda kart için sıralı işlem daha hızlı olabilir
// / - **Veritabanı Yükü**: Paralel sorgular DB'ye aynı anda yük bindirir
// /
// / ## Kullanım Senaryoları
// /
// / ### Senaryo 1: Dashboard Kartları
// / ```go
// / // 10 farklı metrik kartı paralel olarak yükle
// / handler := &FieldHandler{
// /     Cards: []widget.Card{
// /         userCountCard,
// /         revenueCard,
// /         activeSessionsCard,
// /         // ... 7 kart daha
// /     },
// /     DB: db,
// / }
// / err := HandleCardList(handler, ctx)
// / ```
// /
// / ### Senaryo 2: Analitik Paneli
// / ```go
// / // Ağır hesaplama gerektiren kartlar
// / handler := &FieldHandler{
// /     Cards: []widget.Card{
// /         complexQueryCard,    // 2 saniye
// /         aggregationCard,     // 3 saniye
// /         reportCard,          // 1 saniye
// /     },
// /     DB: db,
// / }
// / // Toplam süre: ~3 saniye (en yavaş kart)
// / // Sıralı işlem: 6 saniye olurdu
// / ```
// /
// / ## Önemli Notlar ve Uyarılar
// /
// / ### ⚠️ Kritik Uyarılar
// /
// / 1. **Veritabanı Bağlantı Havuzu**: Paralel sorgular için yeterli DB connection olmalı
// / 2. **Goroutine Sınırı**: Çok fazla kart (>1000) için worker pool pattern düşünülmeli
// / 3. **Context İptali**: Context cancel durumu kontrol edilmiyor, eklenebilir
// / 4. **Bellek Sızıntısı**: Channel kapatılmazsa goroutine leak olabilir (şu an güvenli)
// /
// / ### 💡 İyileştirme Önerileri
// /
// / 1. **Worker Pool**: Sabit sayıda goroutine ile işlem yapılabilir
// / 2. **Timeout**: Her kart için maksimum çözümleme süresi eklenebilir
// / 3. **Circuit Breaker**: Sürekli hata veren kartlar devre dışı bırakılabilir
// / 4. **Caching**: Sık kullanılan kart verileri cache'lenebilir
// /
// / ## JSON Yanıt Formatı
// /
// / ```json
// / {
// /   "data": [
// /     {
// /       "index": 0,
// /       "name": "User Count",
// /       "component": "MetricCard",
// /       "width": "1/3",
// /       "data": {
// /         "value": 1234,
// /         "trend": "+12%"
// /       }
// /     },
// /     {
// /       "index": 1,
// /       "name": "Revenue",
// /       "component": "MetricCard",
// /       "width": "1/3",
// /       "error": "database connection failed"
// /     }
// /   ]
// / }
// / ```
// /
// / ## Hata Durumları
// /
// / - Kart çözümleme hatası: Kart atlanmaz, "error" alanı eklenir
// / - JSON serialization hatası: Fiber error döner
// / - Boş kart listesi: Boş array döner (hata değil)
// /
// / ## Thread Safety
// /
// / - ✅ Goroutine-safe: Her goroutine kendi verisiyle çalışır
// / - ✅ Channel-safe: Buffered channel ve proper close pattern
// / - ✅ WaitGroup-safe: Doğru Add/Done/Wait kullanımı
// / - ⚠️ Context-safe: Context cancel kontrolü yok (eklenebilir)
// /
// / ## Performans Metrikleri
// /
// / Örnek senaryolar:
// / - 5 kart, her biri 100ms: ~100ms (paralel) vs 500ms (sıralı)
// / - 10 kart, her biri 200ms: ~200ms (paralel) vs 2000ms (sıralı)
// / - 100 kart, her biri 50ms: ~50ms (paralel) vs 5000ms (sıralı)
// /
// / ## İlgili Tipler
// /
// / - `widget.Card`: Kart arayüzü
// / - `FieldHandler`: Handler yapısı
// / - `context.Context`: İstek context'i
// / - `cardResult`: Sonuç yapısı
func HandleCardList(h *FieldHandler, c *context.Context) error {
	if len(h.Cards) == 0 {
		return c.JSON(fiber.Map{
			"data": []map[string]interface{}{},
		})
	}

	// Create buffered channel for results (non-blocking sends)
	results := make(chan cardResult, len(h.Cards))

	// WaitGroup to track goroutine completion
	var wg sync.WaitGroup
	wg.Add(len(h.Cards))

	// Fan-out: Launch goroutines asynchronously for each card
	for i, card := range h.Cards {
		go func(idx int, w widget.Card) {
			defer wg.Done() // Mark goroutine as done when finished

			// Serialize base properties
			serialized := w.JsonSerialize()
			serialized["index"] = idx
			serialized["name"] = w.Name()
			serialized["component"] = w.Component()
			serialized["width"] = w.Width()

			// TODO: Card.Resolve() şimdilik kullanılmayacak (aşağıdaki satırlar comment out edildi)

			// Resolve data. If provider/db is unavailable, card still gets a nil db.
			var db *gorm.DB
			if h.Provider != nil {
				if client, ok := h.Provider.GetClient().(*gorm.DB); ok {
					db = client
				}
			}
			data, err := w.Resolve(c, db)

			// Send result to channel
			results <- cardResult{
				card:       w,
				data:       data,
				err:        err,
				index:      idx,
				serialized: serialized,
			}
		}(i, card)
	}

	// Close channel when all goroutines complete (async closer)
	go func() {
		wg.Wait()      // Wait for all goroutines to finish
		close(results) // Close channel to signal completion
	}()

	// Fan-in: Collect results from channel
	resp := make([]map[string]interface{}, len(h.Cards))

	for result := range results {
		if result.err != nil {
			fmt.Printf("Error resolving card %s: %v\n", result.card.Name(), result.err)
			result.serialized["error"] = result.err.Error()
		} else {
			// Assign resolved data to "data" key
			result.serialized["data"] = result.data
		}

		// Store result at original index to maintain order
		resp[result.index] = result.serialized
	}

	return c.JSON(fiber.Map{
		"data": resp,
	})
}
