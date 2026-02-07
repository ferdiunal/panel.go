package metric

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Bu yapı, veritabanı sorgularından dönen tarih-değer çiftini temsil eder.
//
// Kullanım Senaryosu:
// - Zaman serisi verilerini (time-series data) sorgulamak için kullanılır
// - Tarih bazlı metrikler (günlük sayılar, toplamlar, ortalamalar) için veri taşıyıcı görevi görür
//
// Alanlar:
// - Date: Sorgu sonucunun tarihi (YYYY-MM-DD formatında string)
// - Value: İlgili tarih için hesaplanan metrik değeri (int64)
//
// Örnek Kullanım:
//   var results []Result
//   db.Model(&User{}).
//       Select("strftime('%Y-%m-%d', created_at) as date, count(*) as value").
//       Group("date").
//       Scan(&results)
//   // results[0] = {Date: "2024-01-15", Value: 42}
//
// Not: JSON tag'ları API yanıtlarında kullanılır
type Result struct {
	Date  string `json:"date"`
	Value int64  `json:"value"`
}

// Bu fonksiyon, belirtilen tarih aralığında kayıtları tarihe göre gruplandırarak sayar.
//
// Kullanım Senaryosu:
// - Günlük kullanıcı kaydı sayısını takip etmek
// - Belirli bir zaman periyodunda oluşturulan nesnelerin sayısını görmek
// - Trend analizi ve zaman serisi verilerini görselleştirmek
//
// Parametreler:
// - db: GORM veritabanı bağlantısı
// - model: Sorgulanacak model (örn: &User{}, &Order{})
// - dateColumn: Tarih sütununun adı (örn: "created_at", "updated_at")
// - days: Geçmiş kaç günün verisi alınacağı (örn: 30 = son 30 gün)
//
// Dönüş Değerleri:
// - []TrendPoint: Tarih ve sayı değerlerini içeren trend noktaları
// - error: Sorgu sırasında oluşan hata (nil ise başarılı)
//
// Örnek Kullanım:
//   trends, err := CountByDateRange(db, &User{}, "created_at", 30)
//   if err != nil {
//       log.Fatal(err)
//   }
//   for _, point := range trends {
//       fmt.Printf("%s: %d yeni kullanıcı\n", point.Date.Format("2006-01-02"), point.Value)
//   }
//
// Önemli Notlar:
// - Eksik tarihleri otomatik olarak sıfır değeriyle doldurur
// - SQLite date formatting kullanır (strftime)
// - Sonuçlar tarih sırasına göre sıralanır (ASC)
// - Veritabanı sorgusu BETWEEN operatörü kullanır
func CountByDateRange(db *gorm.DB, model interface{}, dateColumn string, days int) ([]TrendPoint, error) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	var results []Result

	// SQLite date formatting
	dateExpr := fmt.Sprintf("strftime('%%Y-%%m-%%d', %s)", dateColumn)

	err := db.Model(model).
		Select(fmt.Sprintf("%s as date, count(*) as value", dateExpr)).
		Where(fmt.Sprintf("%s BETWEEN ? AND ?", dateColumn), startDate, endDate).
		Group("date").
		Order("date ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// Convert to TrendPoint and fill gaps
	return fillDateGaps(results, days), nil
}

// Bu fonksiyon, belirtilen tarih aralığında bir sütunun değerlerini tarihe göre gruplandırarak toplar.
//
// Kullanım Senaryosu:
// - Günlük satış tutarını hesaplamak
// - Belirli bir zaman periyodunda oluşturulan siparişlerin toplam değerini görmek
// - Gelir takibi ve finansal raporlama
// - Zaman serisi verilerinde kümülatif değerleri analiz etmek
//
// Parametreler:
// - db: GORM veritabanı bağlantısı
// - model: Sorgulanacak model (örn: &Order{}, &Transaction{})
// - dateColumn: Tarih sütununun adı (örn: "created_at", "transaction_date")
// - sumColumn: Toplanacak sütunun adı (örn: "amount", "total_price")
// - days: Geçmiş kaç günün verisi alınacağı (örn: 30 = son 30 gün)
//
// Dönüş Değerleri:
// - []TrendPoint: Tarih ve toplam değerlerini içeren trend noktaları
// - error: Sorgu sırasında oluşan hata (nil ise başarılı)
//
// Örnek Kullanım:
//   trends, err := SumByDateRange(db, &Order{}, "created_at", "total_amount", 30)
//   if err != nil {
//       log.Fatal(err)
//   }
//   for _, point := range trends {
//       fmt.Printf("%s: %.2f TL toplam satış\n", point.Date.Format("2006-01-02"), float64(point.Value)/100)
//   }
//
// Önemli Notlar:
// - NULL değerleri otomatik olarak 0 olarak işlenir (COALESCE kullanır)
// - Eksik tarihleri otomatik olarak sıfır değeriyle doldurur
// - SQLite date formatting kullanır (strftime)
// - Sonuçlar tarih sırasına göre sıralanır (ASC)
// - Veritabanı sorgusu BETWEEN operatörü kullanır
// - Sayısal sütunlar için kullanılmalıdır (int, float, decimal)
func SumByDateRange(db *gorm.DB, model interface{}, dateColumn, sumColumn string, days int) ([]TrendPoint, error) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	var results []Result

	// SQLite date formatting
	dateExpr := fmt.Sprintf("strftime('%%Y-%%m-%%d', %s)", dateColumn)

	err := db.Model(model).
		Select(fmt.Sprintf("%s as date, COALESCE(SUM(%s), 0) as value", dateExpr, sumColumn)).
		Where(fmt.Sprintf("%s BETWEEN ? AND ?", dateColumn), startDate, endDate).
		Group("date").
		Order("date ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// Convert to TrendPoint and fill gaps
	return fillDateGaps(results, days), nil
}

// Bu fonksiyon, belirtilen sütuna göre kayıtları gruplandırarak her grup için sayı döndürür.
//
// Kullanım Senaryosu:
// - Kullanıcıları ülkeye göre gruplandırarak her ülkedeki kullanıcı sayısını görmek
// - Siparişleri duruma göre gruplandırarak (pending, completed, cancelled) sayılarını görmek
// - Kategorilere göre ürün sayısını analiz etmek
// - Kategorik veri dağılımını anlamak
//
// Parametreler:
// - db: GORM veritabanı bağlantısı
// - model: Sorgulanacak model (örn: &User{}, &Order{})
// - column: Gruplandırılacak sütunun adı (örn: "country", "status", "category")
//
// Dönüş Değerleri:
// - map[string]int64: Sütun değeri -> sayı eşlemesi
// - error: Sorgu sırasında oluşan hata (nil ise başarılı)
//
// Örnek Kullanım:
//   groupData, err := GroupByColumn(db, &User{}, "country")
//   if err != nil {
//       log.Fatal(err)
//   }
//   for country, count := range groupData {
//       fmt.Printf("%s: %d kullanıcı\n", country, count)
//   }
//   // Çıktı:
//   // Turkey: 1500 kullanıcı
//   // Germany: 800 kullanıcı
//   // France: 600 kullanıcı
//
// Önemli Notlar:
// - Sonuçlar map olarak döndürülür (sırasız)
// - NULL değerleri string olarak "NULL" veya boş string olarak görünebilir
// - Büyük veri setleri için performans etkileyebilir
// - Kategorik sütunlar için idealdir (string, enum, status vb.)
func GroupByColumn(db *gorm.DB, model interface{}, column string) (map[string]int64, error) {
	type Result struct {
		Key   string `json:"key"`
		Value int64  `json:"value"`
	}

	var results []Result

	err := db.Model(model).
		Select(fmt.Sprintf("%s as key, count(*) as value", column)).
		Group(column).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// Convert to map
	data := make(map[string]int64)
	for _, r := range results {
		data[r.Key] = r.Value
	}

	return data, nil
}

// Bu fonksiyon, belirtilen koşulu sağlayan kayıtların sayısını döndürür.
//
// Kullanım Senaryosu:
// - Aktif kullanıcı sayısını bulmak
// - Belirli bir durumda olan siparişleri saymak
// - Belirli bir tarihten sonra oluşturulan kayıtları saymak
// - Koşullu filtreleme ile veri analizi yapmak
//
// Parametreler:
// - db: GORM veritabanı bağlantısı
// - model: Sorgulanacak model (örn: &User{}, &Order{})
// - condition: WHERE koşulu (örn: "status = ?", "age > ?", "email LIKE ?")
// - args: Koşuldaki placeholder'lar için değerler (örn: "active", 18, "%@gmail.com")
//
// Dönüş Değerleri:
// - int64: Koşulu sağlayan kayıt sayısı
// - error: Sorgu sırasında oluşan hata (nil ise başarılı)
//
// Örnek Kullanım:
//   // Aktif kullanıcı sayısı
//   count, err := CountWhere(db, &User{}, "status = ?", "active")
//   if err != nil {
//       log.Fatal(err)
//   }
//   fmt.Printf("Aktif kullanıcı: %d\n", count)
//
//   // Birden fazla koşul
//   count, err := CountWhere(db, &Order{}, "status = ? AND total > ?", "completed", 1000)
//   if err != nil {
//       log.Fatal(err)
//   }
//   fmt.Printf("1000 TL üzeri tamamlanan siparişler: %d\n", count)
//
// Önemli Notlar:
// - Parametreli sorgular kullanır (SQL injection'a karşı güvenli)
// - Koşul parametresi SQL WHERE cümlesi olmalıdır
// - args parametresi koşuldaki ? placeholder'larına karşılık gelir
// - Hata durumunda 0 ve error döndürür
func CountWhere(db *gorm.DB, model interface{}, condition string, args ...interface{}) (int64, error) {
	var count int64
	err := db.Model(model).Where(condition, args...).Count(&count).Error
	return count, err
}

// Bu fonksiyon, belirtilen koşulu sağlayan kayıtların belirli bir sütunundaki değerleri toplar.
//
// Kullanım Senaryosu:
// - Aktif kullanıcıların toplam harcamasını hesaplamak
// - Belirli bir durumda olan siparişlerin toplam tutarını bulmak
// - Belirli bir kategorideki ürünlerin toplam stok miktarını görmek
// - Koşullu finansal raporlama ve analiz
//
// Parametreler:
// - db: GORM veritabanı bağlantısı
// - model: Sorgulanacak model (örn: &Order{}, &Transaction{})
// - column: Toplanacak sütunun adı (örn: "amount", "total_price", "quantity")
// - condition: WHERE koşulu (örn: "status = ?", "user_id = ?", "created_at > ?")
// - args: Koşuldaki placeholder'lar için değerler
//
// Dönüş Değerleri:
// - int64: Koşulu sağlayan kayıtların sütun değerlerinin toplamı
// - error: Sorgu sırasında oluşan hata (nil ise başarılı)
//
// Örnek Kullanım:
//   // Tamamlanan siparişlerin toplam tutarı
//   total, err := SumWhere(db, &Order{}, "total_amount", "status = ?", "completed")
//   if err != nil {
//       log.Fatal(err)
//   }
//   fmt.Printf("Tamamlanan siparişler toplam: %.2f TL\n", float64(total)/100)
//
//   // Belirli bir kullanıcının harcaması
//   total, err := SumWhere(db, &Transaction{}, "amount", "user_id = ? AND type = ?", 123, "expense")
//   if err != nil {
//       log.Fatal(err)
//   }
//   fmt.Printf("Kullanıcı 123 toplam harcama: %.2f TL\n", float64(total)/100)
//
// Önemli Notlar:
// - NULL değerleri otomatik olarak 0 olarak işlenir (COALESCE kullanır)
// - Parametreli sorgular kullanır (SQL injection'a karşı güvenli)
// - Koşul sağlayan kayıt yoksa 0 döndürür
// - Sayısal sütunlar için kullanılmalıdır (int, float, decimal)
// - Hata durumunda 0 ve error döndürür
func SumWhere(db *gorm.DB, model interface{}, column, condition string, args ...interface{}) (int64, error) {
	type Result struct {
		Total int64 `json:"total"`
	}

	var result Result
	err := db.Model(model).
		Select(fmt.Sprintf("COALESCE(SUM(%s), 0) as total", column)).
		Where(condition, args...).
		Scan(&result).Error

	return result.Total, err
}

// Bu fonksiyon, sorgu sonuçlarındaki eksik tarihleri sıfır değeriyle doldurur.
//
// Kullanım Senaryosu:
// - Zaman serisi grafiklerinde boş günleri göstermek
// - Trend analizi için tutarlı veri sağlamak
// - Tarih aralığında hiç veri olmayan günleri 0 değeriyle temsil etmek
// - Grafiklerde kesintisiz bir çizgi oluşturmak
//
// Parametreler:
// - results: Veritabanı sorgusundan dönen Result slice'ı
// - days: Doldurulacak toplam gün sayısı
//
// Dönüş Değerleri:
// - []TrendPoint: Tüm tarihleri içeren ve eksik günleri 0 ile doldurmuş TrendPoint slice'ı
//
// Örnek Kullanım:
//   // Veritabanından sadece 3 gün için veri geldi
//   results := []Result{
//       {Date: "2024-01-15", Value: 10},
//       {Date: "2024-01-17", Value: 20},
//       {Date: "2024-01-19", Value: 15},
//   }
//   // 5 günlük veri istiyoruz
//   points := fillDateGaps(results, 5)
//   // Çıktı:
//   // {Date: 2024-01-15, Value: 10}
//   // {Date: 2024-01-16, Value: 0}   <- Dolduruldu
//   // {Date: 2024-01-17, Value: 20}
//   // {Date: 2024-01-18, Value: 0}   <- Dolduruldu
//   // {Date: 2024-01-19, Value: 15}
//
// Önemli Notlar:
// - İç fonksiyon olarak kullanılır (private)
// - Tarih formatı "2006-01-02" (YYYY-MM-DD) olmalıdır
// - Bugünden geriye doğru tarihler oluşturur
// - Eksik günleri otomatik olarak 0 değeriyle doldurur
// - Grafik ve trend analizi için gerekli olan tutarlı veri sağlar
func fillDateGaps(results []Result, days int) []TrendPoint {
	now := time.Now()
	dateMap := make(map[string]int64)

	// Populate map with query results
	for _, res := range results {
		dateMap[res.Date] = res.Value
	}

	points := make([]TrendPoint, 0, days)

	// Generate dates from (now - days) to now
	start := now.AddDate(0, 0, -days+1)

	for i := 0; i < days; i++ {
		d := start.AddDate(0, 0, i)
		dateStr := d.Format("2006-01-02")

		val := int64(0)
		if v, ok := dateMap[dateStr]; ok {
			val = v
		}

		points = append(points, TrendPoint{
			Date:  d,
			Value: val,
		})
	}

	return points
}

// Bu fonksiyon, belirtilen tarih aralığında bir sütunun ortalama değerini tarihe göre gruplandırarak hesaplar.
//
// Kullanım Senaryosu:
// - Günlük ortalama sipariş tutarını hesaplamak
// - Belirli bir zaman periyodunda ortalama kullanıcı puanını görmek
// - Ürün fiyatlarının ortalama değişimini takip etmek
// - Performans metriklerinin ortalama değerini analiz etmek
//
// Parametreler:
// - db: GORM veritabanı bağlantısı
// - model: Sorgulanacak model (örn: &Order{}, &Review{})
// - dateColumn: Tarih sütununun adı (örn: "created_at", "review_date")
// - avgColumn: Ortalaması alınacak sütunun adı (örn: "amount", "rating", "price")
// - days: Geçmiş kaç günün verisi alınacağı (örn: 30 = son 30 gün)
//
// Dönüş Değerleri:
// - []TrendPoint: Tarih ve ortalama değerlerini içeren trend noktaları
// - error: Sorgu sırasında oluşan hata (nil ise başarılı)
//
// Örnek Kullanım:
//   // Günlük ortalama sipariş tutarı
//   trends, err := AverageByDateRange(db, &Order{}, "created_at", "total_amount", 30)
//   if err != nil {
//       log.Fatal(err)
//   }
//   for _, point := range trends {
//       fmt.Printf("%s: %.2f TL ortalama\n", point.Date.Format("2006-01-02"), float64(point.Value)/100)
//   }
//
//   // Günlük ortalama ürün puanı
//   trends, err := AverageByDateRange(db, &Review{}, "created_at", "rating", 7)
//   if err != nil {
//       log.Fatal(err)
//   }
//   for _, point := range trends {
//       fmt.Printf("%s: %.1f/5 ortalama puan\n", point.Date.Format("2006-01-02"), float64(point.Value)/10)
//   }
//
// Önemli Notlar:
// - NULL değerleri otomatik olarak 0 olarak işlenir (COALESCE kullanır)
// - Eksik tarihleri otomatik olarak sıfır değeriyle doldurur
// - SQLite date formatting kullanır (strftime)
// - Sonuçlar tarih sırasına göre sıralanır (ASC)
// - Veritabanı sorgusu BETWEEN operatörü kullanır
// - Sayısal sütunlar için kullanılmalıdır (int, float, decimal)
// - Ortalama değerleri int64 olarak döndürülür (ondalık kısım kaybolabilir)
func AverageByDateRange(db *gorm.DB, model interface{}, dateColumn, avgColumn string, days int) ([]TrendPoint, error) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	var results []Result

	// SQLite date formatting
	dateExpr := fmt.Sprintf("strftime('%%Y-%%m-%%d', %s)", dateColumn)

	err := db.Model(model).
		Select(fmt.Sprintf("%s as date, COALESCE(AVG(%s), 0) as value", dateExpr, avgColumn)).
		Where(fmt.Sprintf("%s BETWEEN ? AND ?", dateColumn), startDate, endDate).
		Group("date").
		Order("date ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// Convert to TrendPoint and fill gaps
	return fillDateGaps(results, days), nil
}
