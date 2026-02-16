package metric

import (
	"fmt"
	"sort"
	"strings"
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
//
//	metric := NewPartition("Satış Dağılımı")
//	metric.Query(func(db *gorm.DB) (map[string]int64, error) {
//	    // Veritabanından veri çek
//	})
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
//   - fn: Veritabanından veri çeken fonksiyon
//     Giriş: *gorm.DB (veritabanı bağlantısı)
//     Çıkış: map[string]int64 (segment adı -> değer), error
//
// Döndürür: Yapılandırılmış PartitionMetric pointer'ı (method chaining için)
//
// Önemli notlar:
// - Bu metod çağrılmadığı takdirde Resolve() hata döndürür
// - Fonksiyon nil olmamalıdır
//
// Kullanım örneği:
//
//	metric.Query(func(db *gorm.DB) (map[string]int64, error) {
//	    var result map[string]int64
//	    db.Table("sales").Select("category, COUNT(*) as count").
//	        Group("category").Scan(&result)
//	    return result, nil
//	})
func (m *PartitionMetric) Query(fn func(db *gorm.DB) (map[string]int64, error)) *PartitionMetric {
	m.QueryFunc = fn
	return m
}

// Bu metod, her segment için özel renkler ayarlar.
//
// Parametreler:
//   - colors: Segment adı -> hex renk kodu eşlemesi
//     Örnek: map[string]string{"Elektronik": "#FF5733", "Giyim": "#33FF57"}
//
// Döndürür: Yapılandırılmış PartitionMetric pointer'ı (method chaining için)
//
// Kullanım örneği:
//
//	metric.SetColors(map[string]string{
//	    "Elektronik": "#FF5733",
//	    "Giyim": "#33FF57",
//	    "Yiyecek": "#3357FF",
//	})
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
//
//	metric.SetFormat(FormatCurrency)  // Para birimi olarak göster
//	metric.SetFormat(FormatPercentage) // Yüzde olarak göster
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
//
//	metric.SetWidth("1/2")  // Sayfanın yarısını kapla
//	metric.SetWidth("full") // Sayfanın tamamını kapla
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
//
//	{
//	    "data": map[string]int64,    // Segment adı -> değer
//	    "colors": map[string]string, // Segment adı -> renk
//	    "format": Format,            // Görüntülenme formatı
//	}
func (m *PartitionMetric) Resolve(ctx *context.Context, db *gorm.DB) (interface{}, error) {
	if m.QueryFunc == nil {
		return nil, fmt.Errorf("query function not defined")
	}

	data, err := m.QueryFunc(db)
	if err != nil {
		return nil, err
	}

	chartData, chartColors := buildPartitionChartData(data, m.Colors)

	return map[string]interface{}{
		"data":        data,
		"colors":      m.Colors,
		"format":      m.FormatType,
		"chartData":   chartData,
		"chartColors": chartColors,
	}, nil
}

// Bu metod, metriği JSON formatında serileştirir.
//
// Döndürür: Metrik bilgilerini içeren map[string]interface{}
//
// Döndürülen veri yapısı:
//
//	{
//	    "component": string,         // Bileşen adı ("partition-metric")
//	    "title": string,             // Metrik başlığı
//	    "width": string,             // Kartın genişliği
//	    "type": string,              // Kart türü ("partition")
//	    "format": Format,            // Görüntülenme formatı
//	    "colors": map[string]string, // Segment renkleri
//	}
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
	CurrentFunc  func(db *gorm.DB) (int64, error)
	HistoryFunc  func(db *gorm.DB) ([]map[string]interface{}, error)
	Target       int64
	FormatType   Format
	Subtitle     string
	Series       map[string]ProgressSeriesConfig
	ActiveSeries string
}

type ProgressSeriesConfig struct {
	Key     string `json:"key,omitempty"`
	Label   string `json:"label,omitempty"`
	Color   string `json:"color,omitempty"`
	Enabled bool   `json:"enabled"`
}

func defaultProgressSeries() map[string]ProgressSeriesConfig {
	return map[string]ProgressSeriesConfig{
		"desktop": {
			Key:     "desktop",
			Label:   "Desktop",
			Color:   "var(--chart-1)",
			Enabled: true,
		},
		"mobile": {
			Key:     "mobile",
			Label:   "Mobile",
			Color:   "var(--chart-2)",
			Enabled: true,
		},
	}
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
//
//	metric := NewProgress("Satış Hedefi", 100000)
//	metric.Current(func(db *gorm.DB) (int64, error) {
//	    var count int64
//	    db.Model(&Sale{}).Count(&count)
//	    return count, nil
//	})
func NewProgress(title string, target int64) *ProgressMetric {
	return &ProgressMetric{
		BaseCard: widget.BaseCard{
			TitleStr:     title,
			ComponentStr: "progress-metric",
			WidthStr:     "1/3",
			CardTypeVal:  "progress",
		},
		Target:       target,
		FormatType:   FormatNumber,
		ActiveSeries: "desktop",
		Series:       map[string]ProgressSeriesConfig{},
	}
}

// Bu metod, mevcut değeri çeken fonksiyonu ayarlar.
//
// Parametreler:
//   - fn: Mevcut değeri veritabanından çeken fonksiyon
//     Giriş: *gorm.DB (veritabanı bağlantısı)
//     Çıkış: int64 (mevcut değer), error
//
// Döndürür: Yapılandırılmış ProgressMetric pointer'ı (method chaining için)
//
// Önemli notlar:
// - Bu metod çağrılmadığı takdirde Resolve() hata döndürür
// - Fonksiyon nil olmamalıdır
// - Döndürülen değer Target ile karşılaştırılarak yüzde hesaplanır
//
// Kullanım örneği:
//
//	metric.Current(func(db *gorm.DB) (int64, error) {
//	    var current int64
//	    db.Table("sales").Where("status = ?", "completed").Count(&current)
//	    return current, nil
//	})
func (m *ProgressMetric) Current(fn func(db *gorm.DB) (int64, error)) *ProgressMetric {
	m.CurrentFunc = fn
	return m
}

// History, line chart için zaman serisi verisi döndüren sorgu fonksiyonunu ayarlar.
//
// Dönen veri, aşağıdaki alanları içermelidir:
//   - date (string, YYYY-MM-DD)
//   - series konfigürasyonunda tanımlanan seri key'leri (number)
//     (geri uyumluluk için desktop/mobile da desteklenir)
func (m *ProgressMetric) History(fn func(db *gorm.DB) ([]map[string]interface{}, error)) *ProgressMetric {
	m.HistoryFunc = fn
	return m
}

func (m *ProgressMetric) SetSubtitle(subtitle string) *ProgressMetric {
	m.Subtitle = subtitle
	return m
}

func (m *ProgressMetric) ensureSeriesMap() {
	if m.Series == nil {
		m.Series = make(map[string]ProgressSeriesConfig)
	}
}

func (m *ProgressMetric) resolveSeriesAlias(seriesKey string) string {
	alias := normalizeProgressSeriesAlias(seriesKey)
	if alias == "" {
		return ""
	}
	if _, ok := m.Series[alias]; ok {
		return alias
	}

	// Backward compatibility: if caller references old data keys directly.
	for key, cfg := range m.Series {
		if strings.EqualFold(cfg.Key, alias) {
			return key
		}
	}

	return alias
}

func (m *ProgressMetric) ensureSeriesEntry(seriesKey string) (string, ProgressSeriesConfig) {
	m.ensureSeriesMap()

	alias := m.resolveSeriesAlias(seriesKey)
	if alias == "" {
		return "", ProgressSeriesConfig{}
	}

	cfg, exists := m.Series[alias]
	if !exists {
		cfg = ProgressSeriesConfig{
			Key:     alias,
			Label:   defaultProgressSeriesLabel(alias),
			Color:   defaultChartColor(len(m.Series)),
			Enabled: true,
		}
	}
	if cfg.Key == "" {
		cfg.Key = alias
	}

	return alias, cfg
}

func (m *ProgressMetric) SetSeriesLabel(seriesKey, label string) *ProgressMetric {
	alias, cfg := m.ensureSeriesEntry(seriesKey)
	if alias == "" {
		return m
	}

	if strings.TrimSpace(label) != "" {
		cfg.Label = strings.TrimSpace(label)
	}
	m.Series[alias] = cfg
	return m
}

func (m *ProgressMetric) SetSeriesColor(seriesKey, color string) *ProgressMetric {
	alias, cfg := m.ensureSeriesEntry(seriesKey)
	if alias == "" {
		return m
	}

	if strings.TrimSpace(color) != "" {
		cfg.Color = strings.TrimSpace(color)
	}
	m.Series[alias] = cfg
	return m
}

func (m *ProgressMetric) SetSeriesEnabled(seriesKey string, enabled bool) *ProgressMetric {
	alias, cfg := m.ensureSeriesEntry(seriesKey)
	if alias == "" {
		return m
	}

	cfg.Enabled = enabled
	m.Series[alias] = cfg
	return m
}

func (m *ProgressMetric) SetSeriesKey(seriesAlias, dataKey string) *ProgressMetric {
	alias, cfg := m.ensureSeriesEntry(seriesAlias)
	if alias == "" {
		return m
	}

	normalizedDataKey := normalizeProgressSeriesDataKey(dataKey, cfg.Key)
	if normalizedDataKey == "" {
		return m
	}

	cfg.Key = normalizedDataKey
	m.Series[alias] = cfg

	if strings.EqualFold(m.ActiveSeries, alias) {
		m.ActiveSeries = normalizedDataKey
	}

	return m
}

func (m *ProgressMetric) SetActiveSeries(seriesKey string) *ProgressMetric {
	m.ActiveSeries = strings.TrimSpace(seriesKey)
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
//
//	metric.SetFormat(FormatCurrency)  // Para birimi olarak göster
//	metric.SetFormat(FormatPercentage) // Yüzde olarak göster
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
//
//	metric.SetWidth("1/2")  // Sayfanın yarısını kapla
//	metric.SetWidth("full") // Sayfanın tamamını kapla
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
//
//	{
//	    "current": int64,      // Mevcut değer
//	    "target": int64,       // Hedef değer
//	    "percentage": float64, // Yüzde (0-100 arası)
//	    "format": Format,      // Görüntülenme formatı
//	}
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

	seriesConfig := m.resolveSeriesConfig()
	activeSeries := m.resolveActiveSeriesFromConfig(seriesConfig)
	seriesOrder := orderedProgressSeriesAliases(seriesConfig, activeSeries)
	chartData := make([]map[string]interface{}, 0)
	if m.HistoryFunc != nil {
		historyData, err := m.HistoryFunc(db)
		if err != nil {
			return nil, err
		}
		chartData = normalizeLineChartData(historyData, current, m.Target, seriesConfig, seriesOrder, activeSeries)
	} else {
		chartData = buildProgressFallbackChartData(current, m.Target, seriesConfig, seriesOrder, activeSeries)
	}

	return map[string]interface{}{
		"current":      current,
		"target":       m.Target,
		"percentage":   percentage,
		"format":       m.FormatType,
		"chartData":    chartData,
		"subtitle":     m.Subtitle,
		"series":       seriesConfig,
		"activeSeries": activeSeries,
		"seriesOrder":  seriesOrder,
	}, nil
}

// Bu metod, metriği JSON formatında serileştirir.
//
// Döndürür: Metrik bilgilerini içeren map[string]interface{}
//
// Döndürülen veri yapısı:
//
//	{
//	    "component": string, // Bileşen adı ("progress-metric")
//	    "title": string,     // Metrik başlığı
//	    "width": string,     // Kartın genişliği
//	    "type": string,      // Kart türü ("progress")
//	    "target": int64,     // Hedef değer
//	    "format": Format,    // Görüntülenme formatı
//	}
//
// Kullanım senaryosu:
// - Frontend'e gönderilecek JSON yanıtı oluşturmak
// - Metrik konfigürasyonunu API üzerinden iletmek
func (m *ProgressMetric) JsonSerialize() map[string]interface{} {
	seriesConfig := m.resolveSeriesConfig()
	activeSeries := m.resolveActiveSeriesFromConfig(seriesConfig)
	return map[string]interface{}{
		"component":    m.Component(),
		"title":        m.Name(),
		"width":        m.Width(),
		"type":         m.GetType(),
		"target":       m.Target,
		"format":       m.FormatType,
		"subtitle":     m.Subtitle,
		"series":       seriesConfig,
		"activeSeries": activeSeries,
		"seriesOrder":  orderedProgressSeriesAliases(seriesConfig, activeSeries),
	}
}

func (m *ProgressMetric) resolveSeriesConfig() map[string]ProgressSeriesConfig {
	source := m.Series
	if len(source) == 0 {
		source = defaultProgressSeries()
	}

	series := make(map[string]ProgressSeriesConfig, len(source))
	sourceKeys := make([]string, 0, len(source))
	for key := range source {
		sourceKeys = append(sourceKeys, key)
	}
	sort.Strings(sourceKeys)

	for _, rawKey := range sourceKeys {
		cfg := source[rawKey]
		alias := normalizeProgressSeriesAlias(rawKey)
		if alias == "" {
			continue
		}

		resolved := series[alias]
		dataKey := normalizeProgressSeriesDataKey(cfg.Key, alias)
		if resolved.Key == "" {
			resolved.Key = dataKey
		}
		if cfg.Label != "" {
			resolved.Label = cfg.Label
		}
		if cfg.Color != "" {
			resolved.Color = cfg.Color
		}
		resolved.Enabled = cfg.Enabled
		series[alias] = resolved
	}

	aliases := resolveProgressSeriesAliases(series)
	usedDataKeys := make(map[string]struct{}, len(aliases))

	for i, alias := range aliases {
		cfg := series[alias]
		cfg.Key = normalizeProgressSeriesDataKey(cfg.Key, alias)
		if cfg.Label == "" {
			cfg.Label = defaultProgressSeriesLabel(alias)
		}
		if cfg.Color == "" {
			cfg.Color = defaultChartColor(i)
		}

		baseKey := cfg.Key
		uniqueKey := baseKey
		uniqueCounter := 2
		for {
			if _, exists := usedDataKeys[uniqueKey]; !exists {
				break
			}
			uniqueKey = normalizeProgressSeriesDataKey(fmt.Sprintf("%s-%d", baseKey, uniqueCounter), alias)
			uniqueCounter++
		}
		cfg.Key = uniqueKey
		usedDataKeys[uniqueKey] = struct{}{}
		series[alias] = cfg
	}

	enabledFound := false
	for _, alias := range aliases {
		if series[alias].Enabled {
			enabledFound = true
			break
		}
	}
	if !enabledFound && len(aliases) > 0 {
		firstAlias := aliases[0]
		cfg := series[firstAlias]
		cfg.Enabled = true
		series[firstAlias] = cfg
	}

	return series
}

func (m *ProgressMetric) resolveActiveSeries() string {
	return m.resolveActiveSeriesFromConfig(m.resolveSeriesConfig())
}

func (m *ProgressMetric) resolveActiveSeriesFromConfig(series map[string]ProgressSeriesConfig) string {
	active := strings.TrimSpace(m.ActiveSeries)
	normalizedActive := normalizeProgressSeriesAlias(active)

	if cfg, exists := series[normalizedActive]; exists && cfg.Enabled {
		return cfg.Key
	}

	aliases := resolveProgressSeriesAliases(series)
	for _, alias := range aliases {
		cfg, ok := series[alias]
		if !ok || !cfg.Enabled {
			continue
		}
		if strings.EqualFold(cfg.Key, active) {
			return cfg.Key
		}
	}

	for _, alias := range aliases {
		cfg := series[alias]
		if cfg.Enabled {
			return cfg.Key
		}
	}

	if len(aliases) > 0 {
		return series[aliases[0]].Key
	}

	return "desktop"
}

func resolveProgressSeriesAliases(series map[string]ProgressSeriesConfig) []string {
	aliases := make([]string, 0, len(series))
	for alias := range series {
		aliases = append(aliases, alias)
	}
	sort.Strings(aliases)
	return aliases
}

func orderedProgressSeriesAliases(series map[string]ProgressSeriesConfig, activeSeriesKey string) []string {
	aliases := resolveProgressSeriesAliases(series)
	sort.SliceStable(aliases, func(i, j int) bool {
		iActive := strings.EqualFold(series[aliases[i]].Key, activeSeriesKey)
		jActive := strings.EqualFold(series[aliases[j]].Key, activeSeriesKey)
		if iActive != jActive {
			return iActive
		}
		return aliases[i] < aliases[j]
	})
	return aliases
}

func defaultProgressSeriesLabel(key string) string {
	switch strings.TrimSpace(strings.ToLower(key)) {
	case "desktop":
		return "Desktop"
	case "mobile":
		return "Mobile"
	}

	parts := strings.Split(strings.ReplaceAll(key, "_", "-"), "-")
	words := make([]string, 0, len(parts))
	for _, part := range parts {
		if part == "" {
			continue
		}
		words = append(words, strings.ToUpper(part[:1])+part[1:])
	}
	if len(words) == 0 {
		return "Series"
	}

	return strings.Join(words, " ")
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
//
//	TableColumn{
//	    Key: "id",
//	    Label: "Kimlik",
//	    Width: "80px",
//	}
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
//
//	metric := NewTable("Son İşlemler")
//	metric.AddColumn("id", "Kimlik", "80px")
//	metric.AddColumn("name", "Ad", "200px")
//	metric.Query(func(db *gorm.DB) ([]map[string]interface{}, error) {
//	    // Veritabanından veri çek
//	})
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
//   - fn: Veritabanından veri çeken fonksiyon
//     Giriş: *gorm.DB (veritabanı bağlantısı)
//     Çıkış: []map[string]interface{} (satır verileri), error
//
// Döndürür: Yapılandırılmış TableMetric pointer'ı (method chaining için)
//
// Önemli notlar:
// - Bu metod çağrılmadığı takdirde Resolve() hata döndürür
// - Fonksiyon nil olmamalıdır
// - Döndürülen harita anahtarları Columns'daki Key değerleriyle eşleşmelidir
//
// Kullanım örneği:
//
//	metric.Query(func(db *gorm.DB) ([]map[string]interface{}, error) {
//	    var transactions []map[string]interface{}
//	    db.Table("transactions").
//	        Select("id, user_name, amount, created_at").
//	        Limit(10).
//	        Scan(&transactions)
//	    return transactions, nil
//	})
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
//
//	columns := []TableColumn{
//	    {Key: "id", Label: "Kimlik", Width: "80px"},
//	    {Key: "name", Label: "Ad", Width: "200px"},
//	    {Key: "email", Label: "E-posta", Width: "250px"},
//	}
//	metric.SetColumns(columns)
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
//
//	metric.AddColumn("id", "Kimlik", "80px")
//	metric.AddColumn("name", "Ad", "200px")
//	metric.AddColumn("email", "E-posta", "250px")
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
//
//	metric.SetWidth("full")  // Sayfanın tamamını kapla
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
//
//	{
//	    "data": []map[string]interface{},  // Tablo satırları
//	    "columns": []TableColumn,          // Sütun tanımlamaları
//	}
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
//
//	{
//	    "component": string,      // Bileşen adı ("table-metric")
//	    "title": string,          // Metrik başlığı
//	    "width": string,          // Kartın genişliği
//	    "type": string,           // Kart türü ("table")
//	    "columns": []TableColumn, // Sütun tanımlamaları
//	}
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

func buildPartitionChartData(data map[string]int64, colors map[string]string) ([]map[string]interface{}, map[string]string) {
	keys := make([]string, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	chartData := make([]map[string]interface{}, 0, len(keys))
	chartColors := make(map[string]string, len(keys))

	for i, key := range keys {
		normalizedKey := normalizeChartKey(key)
		color := colors[key]
		if color == "" {
			color = colors[normalizedKey]
		}
		if color == "" {
			color = defaultChartColor(i)
		}

		chartData = append(chartData, map[string]interface{}{
			"month":   normalizedKey,
			"label":   key,
			"desktop": data[key],
			"fill":    fmt.Sprintf("var(--color-%s)", normalizedKey),
		})

		chartColors[normalizedKey] = color
	}

	return chartData, chartColors
}

func normalizeLineChartData(
	rows []map[string]interface{},
	current, target int64,
	series map[string]ProgressSeriesConfig,
	seriesAliases []string,
	activeSeries string,
) []map[string]interface{} {
	if len(seriesAliases) == 0 {
		seriesAliases = orderedProgressSeriesAliases(series, activeSeries)
	}
	if len(seriesAliases) == 0 {
		seriesAliases = []string{"desktop"}
	}

	chartData := make([]map[string]interface{}, 0, len(rows))
	for i, row := range rows {
		date := firstString(row, "date", "month")
		if date == "" {
			date = time.Now().AddDate(0, 0, -(len(rows) - i - 1)).Format("2006-01-02")
		}

		normalized := map[string]interface{}{
			"date": date,
		}

		for seriesIndex, alias := range seriesAliases {
			cfg, ok := series[alias]
			if !ok {
				continue
			}

			fallbackKeys := []string{cfg.Key, alias}
			switch seriesIndex {
			case 0:
				fallbackKeys = append(fallbackKeys, "current", "value", "desktop")
			case 1:
				fallbackKeys = append(fallbackKeys, "target", "mobile")
			}

			value, exists := firstInt64(row, fallbackKeys...)
			if !exists {
				switch seriesIndex {
				case 0:
					value = current
				case 1:
					value = target
				default:
					value = 0
				}
			}

			normalized[cfg.Key] = value
		}

		chartData = append(chartData, normalized)
	}

	if len(chartData) == 0 {
		return buildProgressFallbackChartData(current, target, series, seriesAliases, activeSeries)
	}

	return chartData
}

func buildProgressFallbackChartData(
	current, target int64,
	series map[string]ProgressSeriesConfig,
	seriesAliases []string,
	activeSeries string,
) []map[string]interface{} {
	if len(seriesAliases) == 0 {
		seriesAliases = orderedProgressSeriesAliases(series, activeSeries)
	}
	if len(seriesAliases) == 0 {
		seriesAliases = []string{"desktop"}
		series["desktop"] = ProgressSeriesConfig{
			Key:     "desktop",
			Label:   "Desktop",
			Color:   defaultChartColor(0),
			Enabled: true,
		}
	}

	const days = 30
	chartData := make([]map[string]interface{}, 0, days)
	now := time.Now()

	for i := 0; i < days; i++ {
		date := now.AddDate(0, 0, -(days - i - 1)).Format("2006-01-02")
		progressValue := int64(float64(current) * float64(i+1) / float64(days))

		normalized := map[string]interface{}{
			"date": date,
		}

		for seriesIndex, alias := range seriesAliases {
			cfg, ok := series[alias]
			if !ok {
				continue
			}

			switch seriesIndex {
			case 0:
				normalized[cfg.Key] = progressValue
			case 1:
				normalized[cfg.Key] = target
			default:
				normalized[cfg.Key] = 0
			}
		}

		chartData = append(chartData, normalized)
	}

	return chartData
}

func normalizeProgressSeriesAlias(seriesKey string) string {
	return normalizeProgressSeriesDataKey(seriesKey, "")
}

func normalizeProgressSeriesDataKey(value, fallback string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		if strings.TrimSpace(fallback) == "" {
			return ""
		}
		value = fallback
	}

	normalized := normalizeChartKey(value)
	if normalized == "" {
		if strings.TrimSpace(fallback) == "" {
			return ""
		}
		return normalizeChartKey(fallback)
	}
	return normalized
}

func defaultChartColor(index int) string {
	return fmt.Sprintf("var(--chart-%d)", (index%5)+1)
}

func normalizeChartKey(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "" {
		return "item"
	}

	builder := strings.Builder{}
	lastDash := false
	for _, r := range value {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			builder.WriteRune(r)
			lastDash = false
			continue
		}

		if !lastDash {
			builder.WriteRune('-')
			lastDash = true
		}
	}

	result := strings.Trim(builder.String(), "-")
	if result == "" {
		return "item"
	}
	return result
}

func firstString(row map[string]interface{}, keys ...string) string {
	for _, key := range keys {
		if value, ok := row[key].(string); ok && value != "" {
			return value
		}
	}
	return ""
}

func firstInt64(row map[string]interface{}, keys ...string) (int64, bool) {
	for _, key := range keys {
		if value, ok := toInt64(row[key]); ok {
			return value, true
		}
	}
	return 0, false
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
//
//	points := []TrendPoint{
//	    {Date: time.Now().AddDate(0, 0, -7), Value: 1000},
//	    {Date: time.Now().AddDate(0, 0, -6), Value: 1200},
//	    {Date: time.Now().AddDate(0, 0, -5), Value: 1100},
//	}
type TrendPoint struct {
	// Date, veri noktasının tarihini temsil eder.
	// JSON'da "date" anahtarı ile serileştirilir.
	Date time.Time `json:"date"`

	// Value, veri noktasının sayısal değerini temsil eder.
	// JSON'da "value" anahtarı ile serileştirilir.
	Value int64 `json:"value"`
}
