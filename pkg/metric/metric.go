package metric

import (
	"fmt"
	"time"

	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/widget"
	"gorm.io/gorm"
)

// Bu tür, metrik verilerinin görüntülenme formatını tanımlar.
// Desteklenen formatlar: sayı, para birimi ve yüzde.
// Kullanım senaryoları: Dashboard'larda farklı veri türlerinin uygun şekilde gösterilmesi.
type Format string

const (
	// FormatNumber, metrik değerini düz sayı olarak gösterir.
	// Örnek: 1000 -> "1000"
	FormatNumber Format = "number"

	// FormatCurrency, metrik değerini para birimi formatında gösterir.
	// Örnek: 1000 -> "$1,000.00"
	FormatCurrency Format = "currency"

	// FormatPercentage, metrik değerini yüzde formatında gösterir.
	// Örnek: 75 -> "75%"
	FormatPercentage Format = "percentage"
)

// Bu yapı, pasta/halka grafik şeklinde metrik verilerini temsil eder.
// Kategorilere ayrılmış verileri görselleştirmek için kullanılır.
//
// Kullanım senaryoları:
// - Satış dağılımını kategorilere göre göstermek
// - Kullanıcı segmentasyonunu görselleştirmek
// - Pazar payını göstermek
//
// Alanlar:
// - QueryFunc: Veritabanından veri çeken fonksiyon (map[string]int64 döndürür)
// - Colors: Her segment için özel renkler (key: segment adı, value: hex renk kodu)
// - FormatType: Verilerin görüntülenme formatı (sayı, para birimi, yüzde)
type PartitionMetric struct {
	widget.BaseCard
	QueryFunc  func(db *gorm.DB) (map[string]int64, error)
	Colors     map[string]string
	FormatType Format
}

// Bu fonksiyon, yeni bir bölüm metriği oluşturur.
//
// Parametreler:
// - title: Metriğin başlığı (dashboard'da gösterilecek)
//
// Döndürür: Yapılandırılmış PartitionMetric pointer'ı
//
// Kullanım örneği:
//   metric := NewPartition("Satış Dağılımı")
//   metric.Query(func(db *gorm.DB) (map[string]int64, error) {
//       // Veritabanından veri çek
//   })
func NewPartition(title string) *PartitionMetric {
	return &PartitionMetric{
		BaseCard: widget.BaseCard{
			TitleStr:     title,
			ComponentStr: "partition-metric",
			WidthStr:     "1/3",
			CardTypeVal:  "partition",
		},
		Colors:     make(map[string]string),
		FormatType: FormatNumber,
	}
}

// Bu metod, veri çekme fonksiyonunu ayarlar.
//
// Parametreler:
// - fn: Veritabanından veri çeken fonksiyon
//       Giriş: *gorm.DB (veritabanı bağlantısı)
//       Çıkış: map[string]int64 (segment adı -> değer), error
//
// Döndürür: Yapılandırılmış PartitionMetric pointer'ı (method chaining için)
//
// Önemli notlar:
// - Bu metod çağrılmadığı takdirde Resolve() hata döndürür
// - Fonksiyon nil olmamalıdır
//
// Kullanım örneği:
//   metric.Query(func(db *gorm.DB) (map[string]int64, error) {
//       var result map[string]int64
//       db.Table("sales").Select("category, COUNT(*) as count").
//           Group("category").Scan(&result)
//       return result, nil
//   })
func (m *PartitionMetric) Query(fn func(db *gorm.DB) (map[string]int64, error)) *PartitionMetric {
	m.QueryFunc = fn
	return m
}

// Bu metod, her segment için özel renkler ayarlar.
//
// Parametreler:
// - colors: Segment adı -> hex renk kodu eşlemesi
//           Örnek: map[string]string{"Elektronik": "#FF5733", "Giyim": "#33FF57"}
//
// Döndürür: Yapılandırılmış PartitionMetric pointer'ı (method chaining için)
//
// Kullanım örneği:
//   metric.SetColors(map[string]string{
//       "Elektronik": "#FF5733",
//       "Giyim": "#33FF57",
//       "Yiyecek": "#3357FF",
//   })
func (m *PartitionMetric) SetColors(colors map[string]string) *PartitionMetric {
	m.Colors = colors
	return m
}

// Bu metod, metrik verilerinin görüntülenme formatını ayarlar.
//
// Parametreler:
// - format: Görüntülenme formatı (FormatNumber, FormatCurrency, FormatPercentage)
//
// Döndürür: Yapılandırılmış PartitionMetric pointer'ı (method chaining için)
//
// Kullanım örneği:
//   metric.SetFormat(FormatCurrency)  // Para birimi olarak göster
//   metric.SetFormat(FormatPercentage) // Yüzde olarak göster
func (m *PartitionMetric) SetFormat(format Format) *PartitionMetric {
	m.FormatType = format
	return m
}

// Bu metod, kartın genişliğini ayarlar.
//
// Parametreler:
// - width: Genişlik değeri (örn: "1/3", "1/2", "full")
//
// Döndürür: Yapılandırılmış PartitionMetric pointer'ı (method chaining için)
//
// Kullanım örneği:
//   metric.SetWidth("1/2")  // Sayfanın yarısını kapla
//   metric.SetWidth("full") // Sayfanın tamamını kapla
func (m *PartitionMetric) SetWidth(width string) *PartitionMetric {
	m.WidthStr = width
	return m
}

// Bu metod, sorguyu çalıştırır ve metrik verilerini döndürür.
//
// Parametreler:
// - ctx: İstek bağlamı (context.Context)
// - db: Veritabanı bağlantısı (*gorm.DB)
//
// Döndürür:
// - interface{}: Veri, renkler ve format bilgisini içeren map
// - error: Sorgu sırasında oluşan hata (QueryFunc tanımlanmamışsa hata döndürür)
//
// Hata durumları:
// - QueryFunc nil ise: "query function not defined" hatası
// - Veritabanı sorgusu başarısız ise: Sorgu hatası
//
// Döndürülen veri yapısı:
//   {
//       "data": map[string]int64,    // Segment adı -> değer
//       "colors": map[string]string, // Segment adı -> renk
//       "format": Format,            // Görüntülenme formatı
//   }
func (m *PartitionMetric) Resolve(ctx *context.Context, db *gorm.DB) (interface{}, error) {
	if m.QueryFunc == nil {
		return nil, fmt.Errorf("query function not defined")
	}

	data, err := m.QueryFunc(db)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"data":   data,
		"colors": m.Colors,
		"format": m.FormatType,
	}, nil
}

// Bu metod, metriği JSON formatında serileştirir.
//
// Döndürür: Metrik bilgilerini içeren map[string]interface{}
//
// Döndürülen veri yapısı:
//   {
//       "component": string,         // Bileşen adı ("partition-metric")
//       "title": string,             // Metrik başlığı
//       "width": string,             // Kartın genişliği
//       "type": string,              // Kart türü ("partition")
//       "format": Format,            // Görüntülenme formatı
//       "colors": map[string]string, // Segment renkleri
//   }
//
// Kullanım senaryosu:
// - Frontend'e gönderilecek JSON yanıtı oluşturmak
// - Metrik konfigürasyonunu API üzerinden iletmek
func (m *PartitionMetric) JsonSerialize() map[string]interface{} {
	return map[string]interface{}{
		"component": m.Component(),
		"title":     m.Name(),
		"width":     m.Width(),
		"type":      m.GetType(),
		"format":    m.FormatType,
		"colors":    m.Colors,
	}
}

// Bu yapı, ilerleme çubuğu şeklinde metrik verilerini temsil eder.
// Hedef değere karşı mevcut ilerlemeyi göstermek için kullanılır.
//
// Kullanım senaryoları:
// - Proje tamamlanma yüzdesini göstermek
// - Satış hedefine ulaşılan ilerlemeyi göstermek
// - Depo kapasitesi kullanımını göstermek
// - Müşteri kazanım hedefini takip etmek
//
// Alanlar:
// - CurrentFunc: Mevcut değeri veritabanından çeken fonksiyon
// - Target: Hedef değer (ilerleme çubuğunun %100'ü)
// - FormatType: Verilerin görüntülenme formatı (sayı, para birimi, yüzde)
type ProgressMetric struct {
	widget.BaseCard
	CurrentFunc func(db *gorm.DB) (int64, error)
	Target      int64
	FormatType  Format
}

// Bu fonksiyon, yeni bir ilerleme metriği oluşturur.
//
// Parametreler:
// - title: Metriğin başlığı (dashboard'da gösterilecek)
// - target: Hedef değer (ilerleme çubuğunun %100'ü temsil eder)
//
// Döndürür: Yapılandırılmış ProgressMetric pointer'ı
//
// Kullanım örneği:
//   metric := NewProgress("Satış Hedefi", 100000)
//   metric.Current(func(db *gorm.DB) (int64, error) {
//       var count int64
//       db.Model(&Sale{}).Count(&count)
//       return count, nil
//   })
func NewProgress(title string, target int64) *ProgressMetric {
	return &ProgressMetric{
		BaseCard: widget.BaseCard{
			TitleStr:     title,
			ComponentStr: "progress-metric",
			WidthStr:     "1/3",
			CardTypeVal:  "progress",
		},
		Target:     target,
		FormatType: FormatNumber,
	}
}

// Bu metod, mevcut değeri çeken fonksiyonu ayarlar.
//
// Parametreler:
// - fn: Mevcut değeri veritabanından çeken fonksiyon
//       Giriş: *gorm.DB (veritabanı bağlantısı)
//       Çıkış: int64 (mevcut değer), error
//
// Döndürür: Yapılandırılmış ProgressMetric pointer'ı (method chaining için)
//
// Önemli notlar:
// - Bu metod çağrılmadığı takdirde Resolve() hata döndürür
// - Fonksiyon nil olmamalıdır
// - Döndürülen değer Target ile karşılaştırılarak yüzde hesaplanır
//
// Kullanım örneği:
//   metric.Current(func(db *gorm.DB) (int64, error) {
//       var current int64
//       db.Table("sales").Where("status = ?", "completed").Count(&current)
//       return current, nil
//   })
func (m *ProgressMetric) Current(fn func(db *gorm.DB) (int64, error)) *ProgressMetric {
	m.CurrentFunc = fn
	return m
}

// Bu metod, metrik verilerinin görüntülenme formatını ayarlar.
//
// Parametreler:
// - format: Görüntülenme formatı (FormatNumber, FormatCurrency, FormatPercentage)
//
// Döndürür: Yapılandırılmış ProgressMetric pointer'ı (method chaining için)
//
// Kullanım örneği:
//   metric.SetFormat(FormatCurrency)  // Para birimi olarak göster
//   metric.SetFormat(FormatPercentage) // Yüzde olarak göster
func (m *ProgressMetric) SetFormat(format Format) *ProgressMetric {
	m.FormatType = format
	return m
}

// Bu metod, kartın genişliğini ayarlar.
//
// Parametreler:
// - width: Genişlik değeri (örn: "1/3", "1/2", "full")
//
// Döndürür: Yapılandırılmış ProgressMetric pointer'ı (method chaining için)
//
// Kullanım örneği:
//   metric.SetWidth("1/2")  // Sayfanın yarısını kapla
//   metric.SetWidth("full") // Sayfanın tamamını kapla
func (m *ProgressMetric) SetWidth(width string) *ProgressMetric {
	m.WidthStr = width
	return m
}

// Bu metod, sorguyu çalıştırır ve ilerleme verilerini döndürür.
//
// Parametreler:
// - ctx: İstek bağlamı (context.Context)
// - db: Veritabanı bağlantısı (*gorm.DB)
//
// Döndürür:
// - interface{}: Mevcut değer, hedef, yüzde ve format bilgisini içeren map
// - error: Sorgu sırasında oluşan hata (CurrentFunc tanımlanmamışsa hata döndürür)
//
// Hata durumları:
// - CurrentFunc nil ise: "current function not defined" hatası
// - Veritabanı sorgusu başarısız ise: Sorgu hatası
//
// Döndürülen veri yapısı:
//   {
//       "current": int64,      // Mevcut değer
//       "target": int64,       // Hedef değer
//       "percentage": float64, // Yüzde (0-100 arası)
//       "format": Format,      // Görüntülenme formatı
//   }
//
// Önemli notlar:
// - Yüzde otomatik olarak hesaplanır: (current / target) * 100
// - Target 0 ise yüzde 0 olarak ayarlanır (bölme hatası önlemek için)
func (m *ProgressMetric) Resolve(ctx *context.Context, db *gorm.DB) (interface{}, error) {
	if m.CurrentFunc == nil {
		return nil, fmt.Errorf("current function not defined")
	}

	current, err := m.CurrentFunc(db)
	if err != nil {
		return nil, err
	}

	percentage := float64(0)
	if m.Target > 0 {
		percentage = (float64(current) / float64(m.Target)) * 100
	}

	return map[string]interface{}{
		"current":    current,
		"target":     m.Target,
		"percentage": percentage,
		"format":     m.FormatType,
	}, nil
}

// Bu metod, metriği JSON formatında serileştirir.
//
// Döndürür: Metrik bilgilerini içeren map[string]interface{}
//
// Döndürülen veri yapısı:
//   {
//       "component": string, // Bileşen adı ("progress-metric")
//       "title": string,     // Metrik başlığı
//       "width": string,     // Kartın genişliği
//       "type": string,      // Kart türü ("progress")
//       "target": int64,     // Hedef değer
//       "format": Format,    // Görüntülenme formatı
//   }
//
// Kullanım senaryosu:
// - Frontend'e gönderilecek JSON yanıtı oluşturmak
// - Metrik konfigürasyonunu API üzerinden iletmek
func (m *ProgressMetric) JsonSerialize() map[string]interface{} {
	return map[string]interface{}{
		"component": m.Component(),
		"title":     m.Name(),
		"width":     m.Width(),
		"type":      m.GetType(),
		"target":    m.Target,
		"format":    m.FormatType,
	}
}

// Bu yapı, tablo şeklinde metrik verilerini temsil eder.
// Veritabanından çekilen verileri satır ve sütun formatında göstermek için kullanılır.
//
// Kullanım senaryoları:
// - Son işlemleri listeleyen tablo göstermek
// - Ürün envanterini göstermek
// - Müşteri listesini göstermek
// - Raporlama verilerini göstermek
//
// Alanlar:
// - QueryFunc: Veritabanından veri çeken fonksiyon ([]map[string]interface{} döndürür)
// - Columns: Tablonun sütun tanımlamaları (key, label, width)
type TableMetric struct {
	widget.BaseCard
	QueryFunc func(db *gorm.DB) ([]map[string]interface{}, error)
	Columns   []TableColumn
}

// Bu yapı, tablo sütununu tanımlar.
//
// Alanlar:
// - Key: Veri haritasındaki anahtar (örn: "id", "name", "email")
// - Label: Sütun başlığı (kullanıcıya gösterilecek metin)
// - Width: Sütun genişliği (CSS genişlik değeri, örn: "100px", "20%")
//
// Kullanım örneği:
//   TableColumn{
//       Key: "id",
//       Label: "Kimlik",
//       Width: "80px",
//   }
type TableColumn struct {
	Key   string
	Label string
	Width string
}

// Bu fonksiyon, yeni bir tablo metriği oluşturur.
//
// Parametreler:
// - title: Metriğin başlığı (dashboard'da gösterilecek)
//
// Döndürür: Yapılandırılmış TableMetric pointer'ı
//
// Kullanım örneği:
//   metric := NewTable("Son İşlemler")
//   metric.AddColumn("id", "Kimlik", "80px")
//   metric.AddColumn("name", "Ad", "200px")
//   metric.Query(func(db *gorm.DB) ([]map[string]interface{}, error) {
//       // Veritabanından veri çek
//   })
func NewTable(title string) *TableMetric {
	return &TableMetric{
		BaseCard: widget.BaseCard{
			TitleStr:     title,
			ComponentStr: "table-metric",
			WidthStr:     "full",
			CardTypeVal:  widget.CardTypeTable,
		},
		Columns: []TableColumn{},
	}
}

// Bu metod, veri çekme fonksiyonunu ayarlar.
//
// Parametreler:
// - fn: Veritabanından veri çeken fonksiyon
//       Giriş: *gorm.DB (veritabanı bağlantısı)
//       Çıkış: []map[string]interface{} (satır verileri), error
//
// Döndürür: Yapılandırılmış TableMetric pointer'ı (method chaining için)
//
// Önemli notlar:
// - Bu metod çağrılmadığı takdirde Resolve() hata döndürür
// - Fonksiyon nil olmamalıdır
// - Döndürülen harita anahtarları Columns'daki Key değerleriyle eşleşmelidir
//
// Kullanım örneği:
//   metric.Query(func(db *gorm.DB) ([]map[string]interface{}, error) {
//       var transactions []map[string]interface{}
//       db.Table("transactions").
//           Select("id, user_name, amount, created_at").
//           Limit(10).
//           Scan(&transactions)
//       return transactions, nil
//   })
func (m *TableMetric) Query(fn func(db *gorm.DB) ([]map[string]interface{}, error)) *TableMetric {
	m.QueryFunc = fn
	return m
}

// Bu metod, tablo sütunlarını ayarlar.
//
// Parametreler:
// - columns: Sütun tanımlamalarının listesi
//
// Döndürür: Yapılandırılmış TableMetric pointer'ı (method chaining için)
//
// Önemli notlar:
// - Bu metod mevcut sütunları tamamen değiştirir
// - AddColumn() ile tek tek sütun eklemek daha esnek olabilir
//
// Kullanım örneği:
//   columns := []TableColumn{
//       {Key: "id", Label: "Kimlik", Width: "80px"},
//       {Key: "name", Label: "Ad", Width: "200px"},
//       {Key: "email", Label: "E-posta", Width: "250px"},
//   }
//   metric.SetColumns(columns)
func (m *TableMetric) SetColumns(columns []TableColumn) *TableMetric {
	m.Columns = columns
	return m
}

// Bu metod, tabloya yeni bir sütun ekler.
//
// Parametreler:
// - key: Veri haritasındaki anahtar (örn: "id", "name")
// - label: Sütun başlığı (kullanıcıya gösterilecek metin)
// - width: Sütun genişliği (CSS genişlik değeri)
//
// Döndürür: Yapılandırılmış TableMetric pointer'ı (method chaining için)
//
// Kullanım örneği:
//   metric.AddColumn("id", "Kimlik", "80px")
//   metric.AddColumn("name", "Ad", "200px")
//   metric.AddColumn("email", "E-posta", "250px")
func (m *TableMetric) AddColumn(key, label, width string) *TableMetric {
	m.Columns = append(m.Columns, TableColumn{
		Key:   key,
		Label: label,
		Width: width,
	})
	return m
}

// Bu metod, kartın genişliğini ayarlar.
//
// Parametreler:
// - width: Genişlik değeri (örn: "1/3", "1/2", "full")
//
// Döndürür: Yapılandırılmış TableMetric pointer'ı (method chaining için)
//
// Kullanım örneği:
//   metric.SetWidth("full")  // Sayfanın tamamını kapla
func (m *TableMetric) SetWidth(width string) *TableMetric {
	m.WidthStr = width
	return m
}

// Bu metod, sorguyu çalıştırır ve tablo verilerini döndürür.
//
// Parametreler:
// - ctx: İstek bağlamı (context.Context)
// - db: Veritabanı bağlantısı (*gorm.DB)
//
// Döndürür:
// - interface{}: Veri satırları ve sütun tanımlamalarını içeren map
// - error: Sorgu sırasında oluşan hata (QueryFunc tanımlanmamışsa hata döndürür)
//
// Hata durumları:
// - QueryFunc nil ise: "query function not defined" hatası
// - Veritabanı sorgusu başarısız ise: Sorgu hatası
//
// Döndürülen veri yapısı:
//   {
//       "data": []map[string]interface{},  // Tablo satırları
//       "columns": []TableColumn,          // Sütun tanımlamaları
//   }
//
// Önemli notlar:
// - Döndürülen veri satırlarının anahtarları Columns'daki Key değerleriyle eşleşmelidir
// - Boş sonuç seti döndürülebilir (hata değildir)
func (m *TableMetric) Resolve(ctx *context.Context, db *gorm.DB) (interface{}, error) {
	if m.QueryFunc == nil {
		return nil, fmt.Errorf("query function not defined")
	}

	data, err := m.QueryFunc(db)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"data":    data,
		"columns": m.Columns,
	}, nil
}

// Bu metod, metriği JSON formatında serileştirir.
//
// Döndürür: Metrik bilgilerini içeren map[string]interface{}
//
// Döndürülen veri yapısı:
//   {
//       "component": string,      // Bileşen adı ("table-metric")
//       "title": string,          // Metrik başlığı
//       "width": string,          // Kartın genişliği
//       "type": string,           // Kart türü ("table")
//       "columns": []TableColumn, // Sütun tanımlamaları
//   }
//
// Kullanım senaryosu:
// - Frontend'e gönderilecek JSON yanıtı oluşturmak
// - Metrik konfigürasyonunu API üzerinden iletmek
func (m *TableMetric) JsonSerialize() map[string]interface{} {
	return map[string]interface{}{
		"component": m.Component(),
		"title":     m.Name(),
		"width":     m.Width(),
		"type":      m.GetType(),
		"columns":   m.Columns,
	}
}

// Bu yapı, trend grafiğinde bir veri noktasını temsil eder.
// Zaman serisi verilerini göstermek için kullanılır.
//
// Kullanım senaryoları:
// - Günlük satış trendini göstermek
// - Aylık kullanıcı büyümesini göstermek
// - Haftalık ziyaretçi sayısını göstermek
// - Tarihsel veri analizi yapmak
//
// Alanlar:
// - Date: Veri noktasının tarihi (time.Time)
// - Value: Veri noktasının değeri (int64)
//
// JSON Serileştirme:
// - date: ISO 8601 formatında tarih
// - value: Sayısal değer
//
// Kullanım örneği:
//   points := []TrendPoint{
//       {Date: time.Now().AddDate(0, 0, -7), Value: 1000},
//       {Date: time.Now().AddDate(0, 0, -6), Value: 1200},
//       {Date: time.Now().AddDate(0, 0, -5), Value: 1100},
//   }
type TrendPoint struct {
	// Date, veri noktasının tarihini temsil eder.
	// JSON'da "date" anahtarı ile serileştirilir.
	Date time.Time `json:"date"`

	// Value, veri noktasının sayısal değerini temsil eder.
	// JSON'da "value" anahtarı ile serileştirilir.
	Value int64 `json:"value"`
}
