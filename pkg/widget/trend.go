package widget

import (
	"fmt"
	"time"

	"github.com/ferdiunal/panel.go/pkg/context"
	"gorm.io/gorm"
)

// Bu yapı, zaman serisi verilerini görüntülemek için kullanılan trend widget'ını temsil eder.
// Trend widget'ı, belirli bir tarih aralığında veri noktalarının sayısını gösteren bir grafik bileşenidir.
//
// Kullanım Senaryoları:
// - Son 30/60/90 gün içindeki kullanıcı kayıtlarının trendini göstermek
// - Satış verilerinin zaman içindeki değişimini görselleştirmek
// - Sistem etkinliklerinin günlük dağılımını analiz etmek
//
// Alanlar:
// - Title: Widget'ın başlığı (örn: "Yeni Kullanıcılar")
// - Ranges: Kullanıcının seçebileceği gün aralıkları (örn: [30, 60, 90])
// - QueryFunc: Veritabanından veri çeken özel sorgu fonksiyonu
//
// Örnek Kullanım:
//
//	widget := NewTrendWidget("Yeni Kullanıcılar", &User{}, "created_at")
//	data, err := widget.Resolve(ctx, db)
type Trend struct {
	// Title: Widget'ın başlığı, UI'da gösterilecek metin
	Title string
	// Ranges: Kullanıcının seçebileceği gün aralıkları (örn: 30, 60, 90 gün)
	Ranges []int
	// QueryFunc: Veritabanından trend verilerini çeken özel sorgu fonksiyonu.
	// Parametreler: context.Context ve *gorm.DB
	// Dönüş: Veri noktaları ve hata
	QueryFunc func(ctx *context.Context, db *gorm.DB) ([]interface{}, error)
}

// Bu metod, widget'ın adını döndürür.
// Dönüş: Widget'ın başlığı (Title alanı)
//
// Örnek Kullanım:
//
//	widget := NewTrendWidget("Yeni Kullanıcılar", &User{}, "created_at")
//	name := widget.Name() // "Yeni Kullanıcılar"
func (w *Trend) Name() string {
	return w.Title
}

// Bu metod, widget'ın Vue/React bileşen adını döndürür.
// Dönüş: "trend-metric" - Frontend'de kullanılacak bileşen adı
//
// Kullanım Senaryosu:
// - Frontend, bu değeri kullanarak doğru bileşeni render eder
// - Dinamik bileşen yükleme için kullanılır
func (w *Trend) Component() string {
	return "trend-metric"
}

// Bu metod, widget'ın UI'daki genişliğini döndürür.
// Dönüş: "1/3" - Grid sisteminde 1/3 genişlik (3 sütundan 1'i)
//
// Önemli Not:
// - Şu anda sabit değer döndürür, gelecekte yapılandırılabilir hale getirilebilir
// - Grid sistemi 12 sütunlu olduğunda 1/3 = 4 sütun anlamına gelir
func (w *Trend) Width() string {
	return "1/3" // Trend defaults to 1/3, could be configurable
}

// Bu metod, widget'ın türünü döndürür.
// Dönüş: CardTypeTrend - Widget'ın tip tanımlayıcısı
//
// Kullanım Senaryosu:
// - Widget'ları türlerine göre filtrelemek
// - Farklı widget türlerine göre farklı işlemler yapmak
func (w *Trend) GetType() CardType {
	return CardTypeTrend
}

// Bu metod, widget'ın verilerini çözer ve frontend'e göndermek için hazırlar.
// Parametreler:
// - ctx: İstek bağlamı (query parametreleri, kullanıcı bilgisi vb.)
// - db: Veritabanı bağlantısı
// Dönüş: Veri ve hata
//
// Akış:
// 1. QueryFunc'u çağırarak veritabanından veri alır
// 2. Veriyi başlık ile birlikte map'e dönüştürür
// 3. Frontend'e göndermek için hazır hale getirir
//
// Örnek Dönüş Değeri:
//
//	{
//	  "data": [
//	    {"date": "2026-02-01", "value": 10},
//	    {"date": "2026-02-02", "value": 15}
//	  ],
//	  "title": "Yeni Kullanıcılar"
//	}
//
// Hata Durumu:
// - QueryFunc hata döndürürse, nil ve hata döndürülür
func (w *Trend) Resolve(ctx *context.Context, db *gorm.DB) (interface{}, error) {
	data, err := w.QueryFunc(ctx, db)
	if err != nil {
		return nil, err
	}

	chartData := normalizeAreaChartData(data)
	return map[string]interface{}{
		"data":      data,
		"chartData": chartData,
		"title":     w.Title,
	}, nil
}

func normalizeAreaChartData(data []interface{}) []map[string]interface{} {
	chartData := make([]map[string]interface{}, 0, len(data))

	for i, item := range data {
		row, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		desktop, hasDesktop := toInt64(row["desktop"])
		if !hasDesktop {
			desktop, _ = toInt64(row["value"])
		}

		mobile, hasMobile := toInt64(row["mobile"])
		if !hasMobile {
			mobile = 0
		}

		month := ""
		if rawMonth, ok := row["month"].(string); ok && rawMonth != "" {
			month = rawMonth
		}

		if month == "" {
			if rawDate, ok := row["date"].(string); ok && rawDate != "" {
				month = rawDate
			}
		}

		if month == "" {
			month = fmt.Sprintf("item-%d", i+1)
		}

		normalized := map[string]interface{}{
			"month":   month,
			"desktop": desktop,
			"mobile":  mobile,
		}

		if rawDate, ok := row["date"].(string); ok && rawDate != "" {
			normalized["date"] = rawDate
		}

		chartData = append(chartData, normalized)
	}

	return chartData
}

func toInt64(value interface{}) (int64, bool) {
	switch v := value.(type) {
	case int:
		return int64(v), true
	case int8:
		return int64(v), true
	case int16:
		return int64(v), true
	case int32:
		return int64(v), true
	case int64:
		return v, true
	case uint:
		return int64(v), true
	case uint8:
		return int64(v), true
	case uint16:
		return int64(v), true
	case uint32:
		return int64(v), true
	case uint64:
		return int64(v), true
	case float32:
		return int64(v), true
	case float64:
		return int64(v), true
	default:
		return 0, false
	}
}

// Bu metod, hata durumunda frontend'e göndermek için hata bilgisini hazırlar.
// Parametreler:
// - err: Oluşan hata
// Dönüş: Hata bilgisini içeren map
//
// Dönüş Değeri Örneği:
//
//	{
//	  "error": "database connection failed",
//	  "title": "Yeni Kullanıcılar",
//	  "type": "trend"
//	}
//
// Kullanım Senaryosu:
// - Resolve() sırasında hata oluştuğunda çağrılır
// - Frontend'de hata mesajı gösterilir
func (w *Trend) HandleError(err error) map[string]interface{} {
	return map[string]interface{}{
		"error": err.Error(),
		"title": w.Title,
		"type":  CardTypeTrend,
	}
}

// Bu metod, widget'ın meta bilgilerini döndürür.
// Dönüş: Widget'ın yapılandırma bilgilerini içeren map
//
// Dönüş Değeri Örneği:
//
//	{
//	  "name": "Yeni Kullanıcılar",
//	  "component": "trend-metric",
//	  "width": "1/3",
//	  "type": "trend",
//	  "ranges": [30, 60, 90]
//	}
//
// Kullanım Senaryosu:
// - Dashboard'un widget'ları hakkında bilgi alması gerektiğinde
// - Widget'ın özelliklerini dinamik olarak okumak için
// - Admin panelinde widget ayarlarını göstermek için
func (w *Trend) GetMetadata() map[string]interface{} {
	return map[string]interface{}{
		"name":      w.Title,
		"component": "trend-metric",
		"width":     "1/3",
		"type":      CardTypeTrend,
		"ranges":    w.Ranges,
	}
}

// Bu metod, widget'ı JSON formatında serileştirir.
// Dönüş: Widget'ın JSON serileştirilebilir map gösterimi
//
// Dönüş Değeri Örneği:
//
//	{
//	  "component": "trend-metric",
//	  "title": "Yeni Kullanıcılar",
//	  "width": "1/3",
//	  "type": "trend",
//	  "ranges": [30, 60, 90]
//	}
//
// Kullanım Senaryosu:
// - Widget'ı JSON olarak kaydetmek
// - Widget'ı API üzerinden göndermek
// - Widget'ın yapılandırmasını dışa aktarmak
//
// Önemli Not:
// - GetMetadata() ile benzer ancak farklı alanlar içerebilir
// - JSON serileştirme için optimize edilmiştir
func (w *Trend) JsonSerialize() map[string]interface{} {
	return map[string]interface{}{
		"component": "trend-metric",
		"title":     w.Title,
		"width":     "1/3",
		"type":      CardTypeTrend,
		"ranges":    w.Ranges,
	}
}

// Bu yapı, trend grafiğindeki tek bir veri noktasını temsil eder.
// Veritabanından sorgu sonuçlarını bu yapıya dönüştürerek işlenir.
//
// Alanlar:
// - Date: Veri noktasının tarihi (format: "YYYY-MM-DD")
// - Value: Tarih için toplam sayı (örn: o gün kaç kullanıcı kaydoldu)
//
// Kullanım Senaryosu:
// - Veritabanı sorgusundan gelen sonuçları yapılandırılmış hale getirmek
// - Tarih boşluklarını doldurmadan önce veriyi organize etmek
//
// Örnek:
//
//	TrendValue{
//	  Date: "2026-02-01",
//	  Value: 42,
//	}
type TrendValue struct {
	// Date: Veri noktasının tarihi (RFC3339 formatında: "YYYY-MM-DD")
	Date string `json:"date"`
	// Value: Tarih için toplam sayı (int64 kullanılarak büyük sayıları destekler)
	Value int64 `json:"value"`
}

// Bu fonksiyon, veritabanı sorgusundan gelen veri noktaları arasındaki tarih boşluklarını doldurur.
// Eksik tarihlere 0 değeri atayarak, grafik için tam bir tarih aralığı oluşturur.
//
// Parametreler:
// - results: Veritabanından gelen TrendValue dizisi (sırasız olabilir)
// - days: Doldurulacak gün sayısı (örn: 30, 60, 90)
//
// Dönüş: Tarih boşlukları doldurulmuş, kronolojik sırada veri noktaları
//
// Akış:
// 1. Veritabanı sonuçlarını tarih -> değer map'ine dönüştürür
// 2. Bugünden geriye doğru 'days' gün için tarih oluşturur
// 3. Her tarih için map'te değer varsa kullanır, yoksa 0 atar
// 4. Sonuçları kronolojik sırada (eski -> yeni) döndürür
//
// Örnek Kullanım:
//
//	results := []TrendValue{
//	  {Date: "2026-02-01", Value: 10},
//	  {Date: "2026-02-03", Value: 15}, // 02-02 eksik
//	}
//	filled := fillGaps(results, 3)
//	// Dönüş:
//	// [
//	//   {date: "2026-02-01", value: 10},
//	//   {date: "2026-02-02", value: 0},   // Boşluk dolduruldu
//	//   {date: "2026-02-03", value: 15},
//	// ]
//
// Önemli Notlar:
// - Grafik kütüphaneleri genellikle kronolojik sırada veri bekler
// - Boşluk doldurma, grafik görünümünü daha profesyonel hale getirir
// - Tarih formatı her zaman "YYYY-MM-DD" olmalıdır
func fillGaps(results []TrendValue, days int) []map[string]interface{} {
	now := time.Now()
	dateMap := make(map[string]int64)

	// Veritabanı sonuçlarını tarih -> değer map'ine dönüştür
	// Bu, hızlı arama için O(1) zaman karmaşıklığı sağlar
	for _, res := range results {
		dateMap[res.Date] = res.Value
	}

	finalData := make([]map[string]interface{}, 0)

	// Bugünden geriye doğru 'days' gün için tarih oluştur
	// Grafik kütüphaneleri genellikle kronolojik sırada (eski -> yeni) veri bekler
	// Bu nedenle (bugün - days + 1) tarihinden başlıyoruz
	start := now.AddDate(0, 0, -days+1)

	for i := 0; i < days; i++ {
		d := start.AddDate(0, 0, i)
		dateStr := d.Format("2006-01-02")

		// Map'te tarih varsa değeri kullan, yoksa 0 ata
		val := int64(0)
		if v, ok := dateMap[dateStr]; ok {
			val = v
		}

		finalData = append(finalData, map[string]interface{}{
			"date":  dateStr, // Tooltip ve diğer UI öğeleri için tarih bilgisi
			"value": val,
		})
	}

	return finalData
}

// Bu fonksiyon, belirtilen model ve tarih sütunu için trend widget'ı oluşturur.
// Widget, son 30/60/90 gün içindeki kayıtların günlük sayısını gösterir.
//
// Parametreler:
// - title: Widget'ın başlığı (örn: "Yeni Kullanıcılar")
// - model: Veritabanı modeli (örn: &User{}, &Order{})
// - dateColumn: Tarih sütununun adı (örn: "created_at", "updated_at")
//
// Dönüş: Yapılandırılmış Trend widget'ı pointer'ı
//
// Akış:
//  1. Trend yapısını başlık ve aralıklarla oluşturur
//  2. QueryFunc'u tanımlar:
//     a. Query parametresinden aralık alır (varsayılan: 30)
//     b. Aralığı izin verilen değerlere karşı doğrular
//     c. Veritabanından günlük sayıları sorgular
//     d. Tarih boşluklarını doldurur
//     e. Sonuçları döndürür
//
// Kullanım Senaryoları:
// - Yeni kullanıcı kayıtlarının trendini göstermek
// - Satış verilerinin zaman içindeki değişimini analiz etmek
// - Sistem etkinliklerinin günlük dağılımını görselleştirmek
//
// Örnek Kullanım:
//
//	// Yeni kullanıcılar için trend widget'ı oluştur
//	widget := NewTrendWidget("Yeni Kullanıcılar", &User{}, "created_at")
//
//	// Widget'ı dashboard'a ekle
//	dashboard.AddWidget(widget)
//
//	// Veriyi çöz ve frontend'e gönder
//	data, err := widget.Resolve(ctx, db)
//	if err != nil {
//	  errorData := widget.HandleError(err)
//	  // Hata mesajını frontend'e gönder
//	}
//
// Önemli Notlar:
// - Şu anda SQLite strftime() fonksiyonunu kullanır
// - PostgreSQL veya MySQL için dateExpr değiştirilmesi gerekir
// - dateColumn parametresi SQL injection'a karşı doğrulanmalıdır
// - Varsayılan aralıklar: 30, 60, 90 gün
// - Geçersiz aralık istekleri otomatik olarak 30 güne sıfırlanır
//
// SQL Sorgusu Örneği (SQLite):
//
//	SELECT strftime('%Y-%m-%d', created_at) as date, count(*) as value
//	FROM users
//	WHERE created_at BETWEEN ? AND ?
//	GROUP BY date
//	ORDER BY date ASC
//
// Veritabanı Uyumluluğu:
// - SQLite: strftime('%Y-%m-%d', column)
// - PostgreSQL: DATE(column)
// - MySQL: DATE(column)
func NewTrendWidget(title string, model interface{}, dateColumn string) *Trend {
	return &Trend{
		Title:  title,
		Ranges: []int{30, 60, 90},
		QueryFunc: func(ctx *context.Context, db *gorm.DB) ([]interface{}, error) {
			// Query parametresinden aralık değerini al, varsayılan: 30 gün
			// Kullanıcı "?range=60" gibi bir query parametresi gönderebilir
			days := ctx.QueryInt("range", 30)

			// Aralığı izin verilen değerlere karşı doğrula
			// Güvenlik: Sadece 30, 60, 90 gün aralıklarına izin ver
			valid := false
			allowedRanges := []int{30, 60, 90}
			for _, r := range allowedRanges {
				if r == days {
					valid = true
					break
				}
			}
			// Geçersiz aralık istekleri otomatik olarak 30 güne sıfırlanır
			if !valid {
				days = 30
			}

			// Tarih aralığını hesapla
			// endDate: Bugün
			// startDate: Bugünden geriye doğru 'days' gün
			endDate := time.Now()
			startDate := endDate.AddDate(0, 0, -days)

			var results []TrendValue

			// SQLite strftime() fonksiyonunu kullanarak tarih formatı
			// Diğer veritabanları için:
			// - PostgreSQL: DATE(column)
			// - MySQL: DATE(column)
			// dateColumn parametresi SQL injection'a karşı doğrulanmalıdır
			dateExpr := fmt.Sprintf("strftime('%%Y-%%m-%%d', %s)", dateColumn)

			// GORM sorgusu:
			// 1. Model'i belirt
			// 2. Tarih ve sayı sütunlarını seç
			// 3. Tarih aralığına göre filtrele
			// 4. Tarih'e göre grupla
			// 5. Tarih'e göre sırala (eski -> yeni)
			// 6. Sonuçları TrendValue yapısına dönüştür
			err := db.Model(model).
				Select(fmt.Sprintf("%s as date, count(*) as value", dateExpr)).
				Where(fmt.Sprintf("%s BETWEEN ? AND ?", dateColumn), startDate, endDate).
				Group("date").
				Order("date ASC").
				Scan(&results).Error

			// Sorgu hatası kontrol et
			if err != nil {
				return nil, err
			}

			// Tarih boşluklarını doldur
			// Eksik tarihlere 0 değeri atanır
			filled := fillGaps(results, days)

			// []interface{} türüne dönüştür
			// Go'da interface{} slice'ı doğrudan type assertion yapılamaz
			// Bu nedenle manuel olarak dönüştürülmesi gerekir
			interfaceSlice := make([]interface{}, len(filled))
			for i, v := range filled {
				interfaceSlice[i] = v
			}

			return interfaceSlice, nil
		},
	}
}
