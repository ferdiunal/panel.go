package widget

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"gorm.io/gorm"
)

/// # CardType - Kart Türü Tanımı
///
/// `CardType`, panel arayüzünde gösterilecek kartların türünü temsil eden bir string türüdür.
/// Laravel Nova'nın Card konseptine uygun olarak tasarlanmıştır.
///
/// ## Desteklenen Türler
///
/// - `CardTypeValue`: Değer kartları (örn: toplam satış, aktif kullanıcı sayısı)
/// - `CardTypeTrend`: Trend kartları (örn: haftalık büyüme, aylık değişim)
/// - `CardTypeTable`: Tablo kartları (örn: son işlemler, ürün listesi)
/// - `CardTypePartition`: Bölüm kartları (örn: kategori dağılımı, durum dağılımı)
/// - `CardTypeProgress`: İlerleme kartları (örn: proje tamamlanma yüzdesi, görev ilerlemesi)
///
/// ## Kullanım Senaryosu
///
/// ```go
/// // Değer kartı oluşturma
/// card := NewCard("Toplam Satış", "value-metric").
///     SetWidth("1/3")
/// // card.GetType() => CardTypeValue
/// ```
type CardType string

const (
	/// Değer kartı türü - Tek bir metrik değerini gösterir
	CardTypeValue CardType = "value"
	/// Trend kartı türü - Zaman içindeki değişimi gösterir
	CardTypeTrend CardType = "trend"
	/// Tablo kartı türü - Tabular veri gösterir
	CardTypeTable CardType = "table"
	/// Bölüm kartı türü - Veri dağılımını gösterir
	CardTypePartition CardType = "partition"
	/// İlerleme kartı türü - Tamamlanma yüzdesini gösterir
	CardTypeProgress CardType = "progress"
)

/// # Card - Kart Arayüzü
///
/// `Card`, panel sistemindeki tüm kart türlerinin uyması gereken temel arayüzü tanımlar.
/// Laravel Nova'nın Card konseptine uygun olarak tasarlanmıştır.
/// Metrikler (Value, Trend) sadece Card'ın özel türleridir.
///
/// ## Arayüz Metodları
///
/// ### Name() string
/// Kartın adını/başlığını döndürür. Arayüzde gösterilecek başlık metnidir.
///
/// ### Component() string
/// Frontend bileşeninin adını döndürür (örn: "value-metric", "custom-card", "trend-card").
/// Bu değer frontend tarafında hangi React/Vue bileşeninin render edileceğini belirler.
///
/// ### Width() string
/// Kartın genişliğini döndürür. Desteklenen değerler:
/// - "1/3": Sayfanın 1/3'ü (varsayılan)
/// - "1/2": Sayfanın yarısı
/// - "full": Sayfanın tamamı
/// - "2/3": Sayfanın 2/3'ü
///
/// ### GetType() CardType
/// Kartın türünü döndürür (value, trend, table, partition, progress).
///
/// ### Resolve(ctx *context.Context, db *gorm.DB) (interface{}, error)
/// Kartın verilerini çözer ve döndürür. Veritabanı sorgularını burada yapılır.
/// Parametreler:
/// - ctx: İstek bağlamı (kullanıcı bilgisi, izinler vb.)
/// - db: GORM veritabanı bağlantısı
/// Dönüş: Kartın göstereceği veri ve hata (varsa)
///
/// ### HandleError(err error) map[string]interface{}
/// Hata durumunda hata yanıtını hazırlar. Hata mesajı ve kart başlığını içerir.
///
/// ### GetMetadata() map[string]interface{}
/// Kartın meta verilerini döndürür (ad, bileşen, genişlik, tür).
/// Frontend tarafında kartı yapılandırmak için kullanılır.
///
/// ### JsonSerialize() map[string]interface{}
/// Kartı JSON formatında serileştirir. API yanıtında gönderilecek veridir.
///
/// ## Kullanım Senaryosu
///
/// ```go
/// // Özel kart oluşturma
/// card := NewCard("Satış Özeti", "sales-summary").
///     SetWidth("1/2").
///     SetContent(map[string]interface{}{
///         "total": 15000,
///         "currency": "USD",
///     })
///
/// // Kartı çözme
/// data, err := card.Resolve(ctx, db)
/// if err != nil {
///     errorResponse := card.HandleError(err)
///     // Hata yanıtını gönder
/// }
///
/// // Kartı JSON olarak serileştirme
/// jsonData := card.JsonSerialize()
/// // API yanıtında gönder
/// ```
///
/// ## Avantajlar
///
/// - Tutarlı arayüz: Tüm kartlar aynı metodları uygular
/// - Esneklik: Özel kartlar kolayca oluşturulabilir
/// - Tip güvenliği: Go'nun tip sistemi ile compile-time kontrol
/// - Hata yönetimi: Standart hata işleme mekanizması
///
/// ## Önemli Notlar
///
/// - Resolve() metodu veritabanı işlemleri yapabilir, bu nedenle async olarak çalışabilir
/// - GetMetadata() ve JsonSerialize() frontend tarafında kartı yapılandırmak için kullanılır
/// - Component() değeri frontend tarafında tanımlı bir bileşen adı olmalıdır
type Card interface {
	/// Kartın adını/başlığını döndürür
	Name() string
	/// Frontend bileşeninin adını döndürür (örn: "value-metric", "custom-card")
	Component() string
	/// Kartın genişliğini döndürür ("1/3", "1/2", "full", vb.)
	Width() string
	/// Kartın türünü döndürür (value, trend, table, partition, progress)
	GetType() CardType
	/// Kartın verilerini çözer ve döndürür
	Resolve(ctx *context.Context, db *gorm.DB) (interface{}, error)
	/// Hata durumunda hata yanıtını hazırlar
	HandleError(err error) map[string]interface{}
	/// Kartın meta verilerini döndürür
	GetMetadata() map[string]interface{}
	/// Kartı JSON formatında serileştirir
	JsonSerialize() map[string]interface{}
}

/// # CardResolver - Kart Veri Çözücü Arayüzü
///
/// `CardResolver`, kart verilerini çözmek için kullanılan basit bir arayüzü tanımlar.
/// Sadece Resolve metodunu içerir ve daha spesifik veri çözme işlemleri için kullanılabilir.
///
/// ## Arayüz Metodları
///
/// ### Resolve(ctx *context.Context, db *gorm.DB) (interface{}, error)
/// Kartın verilerini çözer ve döndürür.
/// Parametreler:
/// - ctx: İstek bağlamı
/// - db: GORM veritabanı bağlantısı
/// Dönüş: Çözülen veri ve hata (varsa)
///
/// ## Kullanım Senaryosu
///
/// ```go
/// // Özel resolver oluşturma
/// type SalesResolver struct {
///     Period string
/// }
///
/// func (r *SalesResolver) Resolve(ctx *context.Context, db *gorm.DB) (interface{}, error) {
///     var total int64
///     db.Model(&Sale{}).Where("period = ?", r.Period).Count(&total)
///     return map[string]interface{}{"total": total}, nil
/// }
/// ```
///
/// ## Avantajlar
///
/// - Basit ve odaklanmış: Sadece veri çözme işlemini yapar
/// - Yeniden kullanılabilir: Farklı kartlar tarafından kullanılabilir
/// - Test edilebilir: Bağımsız olarak test edilebilir
type CardResolver interface {
	/// Kartın verilerini çözer ve döndürür
	Resolve(ctx *context.Context, db *gorm.DB) (interface{}, error)
}

/// # BaseCard - Temel Kart Yapısı
///
/// `BaseCard`, tüm kart türleri için ortak alanları ve varsayılan davranışları sağlayan
/// temel bir yapıdır. Diğer kart yapıları bu yapıyı gömüp (embed) kendi özel alanlarını eklerler.
///
/// ## Alanlar
///
/// - `TitleStr`: Kartın başlığı/adı
/// - `ComponentStr`: Frontend bileşeninin adı
/// - `WidthStr`: Kartın genişliği
/// - `CardTypeVal`: Kartın türü
///
/// ## Varsayılan Davranışlar
///
/// - Component: Boş ise "card" döndürür
/// - Width: Boş ise "1/3" döndürür
/// - GetType: Boş ise CardTypeValue döndürür
///
/// ## Kullanım Senaryosu
///
/// ```go
/// // BaseCard'ı gömüp özel kart oluşturma
/// type SalesCard struct {
///     BaseCard
///     SalesData map[string]interface{}
/// }
///
/// func (c *SalesCard) Resolve(ctx *context.Context, db *gorm.DB) (interface{}, error) {
///     // Özel çözme mantığı
///     return c.SalesData, nil
/// }
/// ```
///
/// ## Avantajlar
///
/// - Kod tekrarını azaltır: Ortak alanlar bir kez tanımlanır
/// - Tutarlılık: Tüm kartlar aynı varsayılan davranışa sahip
/// - Genişletilebilirlik: Yeni kart türleri kolayca oluşturulabilir
type BaseCard struct {
	/// Kartın başlığı/adı
	TitleStr string
	/// Frontend bileşeninin adı
	ComponentStr string
	/// Kartın genişliği ("1/3", "1/2", "full", vb.)
	WidthStr string
	/// Kartın türü (value, trend, table, partition, progress)
	CardTypeVal CardType
}

/// ## BaseCard Metodları
///
/// ### Name() string
/// Kartın adını/başlığını döndürür. TitleStr alanının değerini direkt olarak döndürür.
///
/// **Dönüş**: Kartın başlığı
///
/// **Örnek**:
/// ```go
/// card := &BaseCard{TitleStr: "Satış Özeti"}
/// fmt.Println(card.Name()) // Output: "Satış Özeti"
/// ```
func (c *BaseCard) Name() string {
	return c.TitleStr
}

/// ### Component() string
/// Frontend bileşeninin adını döndürür. ComponentStr boş ise varsayılan "card" değerini döndürür.
///
/// **Varsayılan Değer**: "card"
///
/// **Dönüş**: Frontend bileşeninin adı
///
/// **Örnek**:
/// ```go
/// card1 := &BaseCard{ComponentStr: "value-metric"}
/// fmt.Println(card1.Component()) // Output: "value-metric"
///
/// card2 := &BaseCard{ComponentStr: ""}
/// fmt.Println(card2.Component()) // Output: "card" (varsayılan)
/// ```
///
/// **Önemli Not**: Frontend tarafında bu bileşen adı tanımlı olmalıdır.
/// Tanımlanmamış bir bileşen adı kullanılırsa frontend hata verebilir.
func (c *BaseCard) Component() string {
	if c.ComponentStr == "" {
		return "card" // Varsayılan bileşen
	}
	return c.ComponentStr
}

/// ### Width() string
/// Kartın genişliğini döndürür. WidthStr boş ise varsayılan "1/3" değerini döndürür.
///
/// **Desteklenen Değerler**:
/// - "1/3": Sayfanın 1/3'ü (varsayılan)
/// - "1/2": Sayfanın yarısı
/// - "2/3": Sayfanın 2/3'ü
/// - "full": Sayfanın tamamı
///
/// **Varsayılan Değer**: "1/3"
///
/// **Dönüş**: Kartın genişliği
///
/// **Örnek**:
/// ```go
/// card1 := &BaseCard{WidthStr: "1/2"}
/// fmt.Println(card1.Width()) // Output: "1/2"
///
/// card2 := &BaseCard{WidthStr: ""}
/// fmt.Println(card2.Width()) // Output: "1/3" (varsayılan)
/// ```
///
/// **Avantajlar**:
/// - Responsive tasarım: Farklı genişliklerde kartlar oluşturabilirsiniz
/// - Esneklik: Dinamik olarak genişlik değiştirebilirsiniz
///
/// **Uyarı**: Geçersiz genişlik değerleri frontend tarafında hata verebilir.
func (c *BaseCard) Width() string {
	if c.WidthStr == "" {
		return "1/3"
	}
	return c.WidthStr
}

/// ### GetType() CardType
/// Kartın türünü döndürür. CardTypeVal boş ise varsayılan CardTypeValue değerini döndürür.
///
/// **Varsayılan Değer**: CardTypeValue
///
/// **Dönüş**: Kartın türü (value, trend, table, partition, progress)
///
/// **Örnek**:
/// ```go
/// card1 := &BaseCard{CardTypeVal: CardTypeTrend}
/// fmt.Println(card1.GetType()) // Output: "trend"
///
/// card2 := &BaseCard{CardTypeVal: ""}
/// fmt.Println(card2.GetType()) // Output: "value" (varsayılan)
/// ```
///
/// **Kullanım Senaryosu**:
/// ```go
/// card := NewCard("Satış", "sales-card")
/// if card.GetType() == CardTypeValue {
///     // Değer kartı için özel işlem
/// }
/// ```
func (c *BaseCard) GetType() CardType {
	if c.CardTypeVal == "" {
		return CardTypeValue // Varsayılan tür
	}
	return c.CardTypeVal
}

/// ### HandleError(err error) map[string]interface{}
/// Hata durumunda hata yanıtını hazırlar. Hata mesajı ve kart başlığını içeren
/// bir harita döndürür. Frontend tarafında hata göstermek için kullanılır.
///
/// **Parametreler**:
/// - err: Oluşan hata
///
/// **Dönüş**: Hata bilgisini içeren harita
/// - "error": Hata mesajı
/// - "title": Kartın başlığı
///
/// **Örnek**:
/// ```go
/// card := &BaseCard{TitleStr: "Satış Özeti"}
/// err := errors.New("veritabanı bağlantısı başarısız")
/// errorResponse := card.HandleError(err)
/// // Output: map[string]interface{}{
/// //     "error": "veritabanı bağlantısı başarısız",
/// //     "title": "Satış Özeti",
/// // }
/// ```
///
/// **Kullanım Senaryosu**:
/// ```go
/// data, err := card.Resolve(ctx, db)
/// if err != nil {
///     errorResponse := card.HandleError(err)
///     // API yanıtında gönder
///     return c.JSON(http.StatusInternalServerError, errorResponse)
/// }
/// ```
///
/// **Avantajlar**:
/// - Tutarlı hata formatı: Tüm kartlar aynı hata formatını kullanır
/// - Kolay hata takibi: Hata mesajı ve kart başlığı birlikte gönderilir
func (c *BaseCard) HandleError(err error) map[string]interface{} {
	return map[string]interface{}{
		"error": err.Error(),
		"title": c.TitleStr,
	}
}

/// ### GetMetadata() map[string]interface{}
/// Kartın meta verilerini döndürür. Frontend tarafında kartı yapılandırmak için kullanılır.
/// Kartın adı, bileşeni, genişliği ve türünü içerir.
///
/// **Dönüş**: Meta veri haritası
/// - "name": Kartın adı
/// - "component": Frontend bileşeninin adı
/// - "width": Kartın genişliği
/// - "type": Kartın türü
///
/// **Örnek**:
/// ```go
/// card := &BaseCard{
///     TitleStr:     "Satış Özeti",
///     ComponentStr: "sales-card",
///     WidthStr:     "1/2",
///     CardTypeVal:  CardTypeValue,
/// }
/// metadata := card.GetMetadata()
/// // Output: map[string]interface{}{
/// //     "name":      "Satış Özeti",
/// //     "component": "sales-card",
/// //     "width":     "1/2",
/// //     "type":      "value",
/// // }
/// ```
///
/// **Kullanım Senaryosu**:
/// ```go
/// // Frontend tarafında kartı yapılandırma
/// metadata := card.GetMetadata()
/// component := metadata["component"].(string)
/// width := metadata["width"].(string)
/// // React/Vue bileşenini render et
/// ```
///
/// **Avantajlar**:
/// - Merkezi yapılandırma: Tüm meta veriler bir yerden alınır
/// - Tip güvenliği: Go'nun tip sistemi ile kontrol edilir
/// - Genişletilebilirlik: Yeni meta veriler kolayca eklenebilir
func (c *BaseCard) GetMetadata() map[string]interface{} {
	return map[string]interface{}{
		"name":      c.TitleStr,
		"component": c.Component(),
		"width":     c.Width(),
		"type":      c.GetType(),
	}
}


/// # CustomCard - Özel Kart Yapısı
///
/// `CustomCard`, dışarıdan keyfi kartlar oluşturmaya izin veren esnek bir yapıdır.
/// BaseCard'ı gömüp (embed) kendi özel alanlarını ekler. Statik içerik veya dinamik
/// veri çözme işlemleri için kullanılabilir.
///
/// ## Alanlar
///
/// - `BaseCard`: Gömülü temel kart yapısı (TitleStr, ComponentStr, WidthStr, CardTypeVal)
/// - `Content`: Statik içerik veya başlangıç verisi
/// - `ResolveFunc`: Kartın verilerini çözmek için özel fonksiyon (opsiyonel)
///
/// ## Kullanım Senaryoları
///
/// ### 1. Statik İçerik ile Kart
/// ```go
/// card := NewCard("Hoş Geldiniz", "welcome-card").
///     SetContent(map[string]interface{}{
///         "message": "Panel'e hoş geldiniz!",
///         "icon": "welcome",
///     })
/// ```
///
/// ### 2. Dinamik Veri ile Kart
/// ```go
/// card := NewCard("Satış Özeti", "sales-card").
///     SetWidth("1/2")
/// card.ResolveFunc = func(ctx *context.Context, db *gorm.DB) (interface{}, error) {
///     var total int64
///     db.Model(&Sale{}).Count(&total)
///     return map[string]interface{}{"total": total}, nil
/// }
/// ```
///
/// ### 3. Fluent API ile Kart
/// ```go
/// card := NewCard("Kullanıcılar", "users-card").
///     SetWidth("full").
///     SetContent([]map[string]interface{}{
///         {"id": 1, "name": "Ali"},
///         {"id": 2, "name": "Ayşe"},
///     })
/// ```
///
/// ## Avantajlar
///
/// - Esneklik: Statik ve dinamik içerik destekler
/// - Basitlik: Hızlı kart oluşturma
/// - Genişletilebilirlik: ResolveFunc ile özel mantık eklenebilir
/// - Fluent API: Zincir yöntemi ile kolay yapılandırma
///
/// ## Dezavantajlar
///
/// - Tür güvenliği: Content alanı interface{} olduğu için tür dönüşümü gerekebilir
/// - Sınırlı doğrulama: Özel doğrulama mantığı eklemek zor olabilir
type CustomCard struct {
	/// Gömülü temel kart yapısı
	BaseCard
	/// Statik içerik veya başlangıç verisi
	Content interface{}
	/// Kartın verilerini çözmek için özel fonksiyon (opsiyonel)
	ResolveFunc func(ctx *context.Context, db *gorm.DB) (interface{}, error)
}

/// ### Resolve(ctx *context.Context, db *gorm.DB) (interface{}, error)
/// Kartın verilerini çözer ve döndürür. ResolveFunc tanımlanmışsa onu çalıştırır,
/// aksi takdirde Content alanını döndürür.
///
/// **Parametreler**:
/// - ctx: İstek bağlamı (kullanıcı bilgisi, izinler vb.)
/// - db: GORM veritabanı bağlantısı
///
/// **Dönüş**: Kartın göstereceği veri ve hata (varsa)
///
/// **Örnek 1 - Statik İçerik**:
/// ```go
/// card := NewCard("Hoş Geldiniz", "welcome-card").
///     SetContent(map[string]interface{}{"message": "Merhaba!"})
/// data, err := card.Resolve(ctx, db)
/// // Output: map[string]interface{}{"message": "Merhaba!"}, nil
/// ```
///
/// **Örnek 2 - Dinamik Veri**:
/// ```go
/// card := NewCard("Satış", "sales-card")
/// card.ResolveFunc = func(ctx *context.Context, db *gorm.DB) (interface{}, error) {
///     var total int64
///     if err := db.Model(&Sale{}).Count(&total).Error; err != nil {
///         return nil, err
///     }
///     return map[string]interface{}{"total": total}, nil
/// }
/// data, err := card.Resolve(ctx, db)
/// // Output: map[string]interface{}{"total": 15000}, nil
/// ```
///
/// **Kullanım Senaryosu**:
/// ```go
/// // API endpoint'te
/// card := getCard()
/// data, err := card.Resolve(ctx, db)
/// if err != nil {
///     return c.JSON(http.StatusInternalServerError, card.HandleError(err))
/// }
/// return c.JSON(http.StatusOK, data)
/// ```
///
/// **Avantajlar**:
/// - Esnek: Statik ve dinamik içerik destekler
/// - Hata yönetimi: Hata durumunda kontrol edilebilir
/// - Veritabanı erişimi: Gerekli verileri çekmek için db kullanılabilir
///
/// **Uyarılar**:
/// - ResolveFunc uzun süren işlemler yapıyorsa timeout ayarlanmalıdır
/// - Veritabanı hataları düzgün şekilde işlenmelidir
func (c *CustomCard) Resolve(ctx *context.Context, db *gorm.DB) (interface{}, error) {
	if c.ResolveFunc != nil {
		return c.ResolveFunc(ctx, db)
	}
	return c.Content, nil
}

/// ### JsonSerialize() map[string]interface{}
/// Kartı JSON formatında serileştirir. API yanıtında gönderilecek veridir.
/// Kartın tüm önemli bilgilerini içerir: bileşen, başlık, genişlik, tür ve içerik.
///
/// **Dönüş**: Serileştirilmiş kart verisi
/// - "component": Frontend bileşeninin adı
/// - "title": Kartın başlığı
/// - "width": Kartın genişliği
/// - "type": Kartın türü
/// - "content": Kartın içeriği
///
/// **Örnek**:
/// ```go
/// card := NewCard("Satış", "sales-card").
///     SetWidth("1/2").
///     SetContent(map[string]interface{}{"total": 15000})
/// json := card.JsonSerialize()
/// // Output: map[string]interface{}{
/// //     "component": "sales-card",
/// //     "title": "Satış",
/// //     "width": "1/2",
/// //     "type": "value",
/// //     "content": map[string]interface{}{"total": 15000},
/// // }
/// ```
///
/// **Kullanım Senaryosu**:
/// ```go
/// // API endpoint'te
/// card := getCard()
/// jsonData := card.JsonSerialize()
/// return c.JSON(http.StatusOK, jsonData)
/// ```
///
/// **Avantajlar**:
/// - Merkezi serileştirme: Tüm kart verisi bir yerden alınır
/// - Tutarlılık: Tüm kartlar aynı formatı kullanır
/// - Frontend uyumluluğu: Frontend tarafında kolayca parse edilebilir
///
/// **Önemli Not**:
/// Content alanı interface{} olduğu için JSON serileştirmesi sırasında
/// Go'nun standart JSON encoder'ı kullanılır. Özel türler için
/// json.Marshaler arayüzünü uygulamanız gerekebilir.
func (c *CustomCard) JsonSerialize() map[string]interface{} {
	return map[string]interface{}{
		"component": c.Component(),
		"title":     c.Name(),
		"width":     c.Width(),
		"type":      c.GetType(),
		"content":   c.Content,
	}
}

/// # NewCard - Basit Kart Oluşturma Fonksiyonu
///
/// `NewCard`, başlık ve bileşen adı ile basit bir özel kart oluşturur.
/// Varsayılan değerler: genişlik "1/3", tür "value".
///
/// ## Parametreler
///
/// - `title`: Kartın başlığı/adı
/// - `component`: Frontend bileşeninin adı
///
/// ## Dönüş
///
/// Yeni oluşturulan CustomCard işaretçisi
///
/// ## Örnek
///
/// ```go
/// // Basit kart oluşturma
/// card := NewCard("Satış Özeti", "sales-card")
/// // card.Name() => "Satış Özeti"
/// // card.Component() => "sales-card"
/// // card.Width() => "1/3" (varsayılan)
/// // card.GetType() => "value" (varsayılan)
/// ```
///
/// ## Fluent API ile Kullanım
///
/// ```go
/// card := NewCard("Satış", "sales-card").
///     SetWidth("1/2").
///     SetContent(map[string]interface{}{
///         "total": 15000,
///         "currency": "USD",
///     })
/// ```
///
/// ## Kullanım Senaryoları
///
/// ### 1. Dashboard Kartları
/// ```go
/// cards := []Card{
///     NewCard("Toplam Satış", "value-metric").SetWidth("1/3"),
///     NewCard("Aktif Kullanıcılar", "value-metric").SetWidth("1/3"),
///     NewCard("Aylık Büyüme", "trend-metric").SetWidth("1/3"),
/// }
/// ```
///
/// ### 2. Dinamik Kartlar
/// ```go
/// for _, metric := range metrics {
///     card := NewCard(metric.Name, metric.Component).
///         SetWidth(metric.Width).
///         SetContent(metric.Data)
///     cards = append(cards, card)
/// }
/// ```
///
/// ## Avantajlar
///
/// - Basitlik: Hızlı kart oluşturma
/// - Varsayılan değerler: Yaygın kullanım durumları için hazır
/// - Fluent API: Zincir yöntemi ile kolay yapılandırma
///
/// ## Önemli Notlar
///
/// - Component adı frontend tarafında tanımlı olmalıdır
/// - Title boş bırakılmamalıdır (arayüzde gösterilir)
/// - Varsayılan genişlik "1/3" olduğu için çoğu durumda SetWidth() çağrılmalıdır
func NewCard(title, component string) *CustomCard {
	return &CustomCard{
		BaseCard: BaseCard{
			TitleStr:     title,
			ComponentStr: component,
			WidthStr:     "1/3",
			CardTypeVal:  CardTypeValue,
		},
	}
}

/// # CustomCard Setter Metodları
///
/// CustomCard için fluent API setter metodları. Zincir yöntemi ile
/// kolay yapılandırma sağlarlar.

/// ### SetWidth(w string) *CustomCard
/// Kartın genişliğini ayarlar ve kartı döndürür (fluent API).
///
/// **Parametreler**:
/// - w: Kartın genişliği ("1/3", "1/2", "2/3", "full", vb.)
///
/// **Dönüş**: Kartın kendisi (fluent API için)
///
/// **Örnek**:
/// ```go
/// card := NewCard("Satış", "sales-card").
///     SetWidth("1/2")
/// // card.Width() => "1/2"
/// ```
///
/// **Fluent API Zinciri**:
/// ```go
/// card := NewCard("Satış", "sales-card").
///     SetWidth("full").
///     SetContent(data)
/// ```
///
/// **Desteklenen Değerler**:
/// - "1/3": Sayfanın 1/3'ü
/// - "1/2": Sayfanın yarısı
/// - "2/3": Sayfanın 2/3'ü
/// - "full": Sayfanın tamamı
///
/// **Avantajlar**:
/// - Fluent API: Zincir yöntemi ile kolay kullanım
/// - Okunabilirlik: Kod daha okunabilir hale gelir
/// - Esneklik: Dinamik olarak genişlik değiştirebilirsiniz
///
/// **Uyarı**: Geçersiz genişlik değerleri frontend tarafında hata verebilir.
func (c *CustomCard) SetWidth(w string) *CustomCard {
	c.WidthStr = w
	return c
}

/// ### SetContent(content interface{}) *CustomCard
/// Kartın içeriğini ayarlar ve kartı döndürür (fluent API).
/// İçerik statik veri veya dinamik olarak Resolve() tarafından çözülen veri olabilir.
///
/// **Parametreler**:
/// - content: Kartın içeriği (herhangi bir Go türü olabilir)
///
/// **Dönüş**: Kartın kendisi (fluent API için)
///
/// **Örnek 1 - Harita İçeriği**:
/// ```go
/// card := NewCard("Satış", "sales-card").
///     SetContent(map[string]interface{}{
///         "total": 15000,
///         "currency": "USD",
///     })
/// ```
///
/// **Örnek 2 - Dizi İçeriği**:
/// ```go
/// card := NewCard("Kullanıcılar", "users-card").
///     SetContent([]map[string]interface{}{
///         {"id": 1, "name": "Ali"},
///         {"id": 2, "name": "Ayşe"},
///     })
/// ```
///
/// **Örnek 3 - Struct İçeriği**:
/// ```go
/// type SalesData struct {
///     Total    int64
///     Currency string
/// }
/// card := NewCard("Satış", "sales-card").
///     SetContent(SalesData{Total: 15000, Currency: "USD"})
/// ```
///
/// **Fluent API Zinciri**:
/// ```go
/// card := NewCard("Satış", "sales-card").
///     SetWidth("1/2").
///     SetContent(map[string]interface{}{"total": 15000})
/// ```
///
/// **Kullanım Senaryoları**:
/// ```go
/// // 1. Statik içerik
/// card := NewCard("Hoş Geldiniz", "welcome-card").
///     SetContent("Panel'e hoş geldiniz!")
///
/// // 2. Dinamik içerik (Resolve() ile birlikte)
/// card := NewCard("Satış", "sales-card")
/// card.ResolveFunc = func(ctx *context.Context, db *gorm.DB) (interface{}, error) {
///     // Veritabanından veri çek
///     return data, nil
/// }
/// card.SetContent(nil) // Başlangıç değeri
///
/// // 3. Karmaşık içerik
/// card := NewCard("Dashboard", "dashboard-card").
///     SetContent(map[string]interface{}{
///         "widgets": []map[string]interface{}{
///             {"type": "metric", "value": 100},
///             {"type": "chart", "data": []int{1, 2, 3}},
///         },
///     })
/// ```
///
/// **Avantajlar**:
/// - Esneklik: Herhangi bir Go türü destekler
/// - Fluent API: Zincir yöntemi ile kolay kullanım
/// - Tip güvenliği: Go'nun tip sistemi ile kontrol edilir
///
/// **Dezavantajlar**:
/// - Tür dönüşümü: interface{} olduğu için tür dönüşümü gerekebilir
/// - Doğrulama: İçerik doğrulaması manuel olarak yapılmalıdır
///
/// **Önemli Notlar**:
/// - Content alanı nil olabilir (Resolve() ile dinamik olarak ayarlanabilir)
/// - JSON serileştirmesi sırasında Go'nun standart JSON encoder'ı kullanılır
/// - Özel türler için json.Marshaler arayüzünü uygulamanız gerekebilir
func (c *CustomCard) SetContent(content interface{}) *CustomCard {
	c.Content = content
	return c
}
