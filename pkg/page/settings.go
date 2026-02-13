// Package page, sistem yönetim panelinin sayfa bileşenlerini içerir.
//
// Bu paket, panel uygulamasında farklı sayfa türlerini (Settings, Dashboard, vb.)
// tanımlamak ve yönetmek için kullanılan temel yapıları ve arayüzleri sağlar.
package page

import (
	"encoding/json"

	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/domain/setting"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/widget"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Settings, sistem ayarlarını yönetmek için kullanılan sayfa bileşenidir.
//
// # Açıklama
// Settings struct'ı, panel uygulamasında sistem genelinde ayarları (konfigürasyonları)
// yönetmek için tasarlanmıştır. Veritabanında saklanan anahtar-değer çiftleri olarak
// çalışır ve dinamik form alanları aracılığıyla kullanıcı arayüzünde sunulur.
//
// # Yapı Alanları
// - Base: Tüm sayfa bileşenlerinin ortak özelliklerini içeren gömülü struct
// - Elements: Ayarlar sayfasında gösterilecek form alanlarının listesi
// - HideInNavigation: Ayarlar sayfasının navigasyon menüsünde gizlenip gizlenmeyeceğini belirler
//
// # Kullanım Senaryoları
// 1. Sistem genelinde konfigürasyonları yönetmek (site adı, logo, tema, vb.)
// 2. Dinamik form alanları ile kullanıcı ayarlarını toplamak
// 3. Ayarları veritabanında kalıcı olarak saklamak
// 4. Navigasyon menüsünde ayarlar sayfasını göstermek/gizlemek
//
// # Örnek Kullanım
// ```go
// settings := &Settings{
//     Elements: []fields.Element{
//         // Form alanları buraya eklenir
//     },
//     HideInNavigation: false,
// }
//
// // Ayarları kaydetmek
// err := settings.Save(ctx, db, map[string]interface{}{
//     "site_name": "Benim Sitesi",
//     "site_url": "https://example.com",
// })
// ```
//
// # Avantajlar
// - Dinamik form alanları ile esnek ayar yönetimi
// - Veritabanında merkezi olarak saklanan ayarlar
// - Navigasyon menüsünde gösterilip gizlenebilir
// - Farklı veri tiplerini (string, JSON) destekler
//
// # Dezavantajlar
// - Karmaşık ayarlar için JSON serileştirmesi gerekebilir
// - Ayarların doğrulanması manuel olarak yapılmalıdır
//
// # Önemli Notlar
// - Ayarlar veritabanında anahtar-değer çiftleri olarak saklanır
// - Aynı anahtarla yeni bir ayar kaydedilirse, eski değer güncellenir
// - Tüm değerler string olarak veya JSON string olarak saklanır
type Settings struct {
	Base
	// Elements, ayarlar sayfasında gösterilecek form alanlarının listesidir.
	// Her Element, bir form alanını temsil eder (TextInput, Select, Checkbox, vb.)
	Elements []fields.Element

	// HideInNavigation, ayarlar sayfasının navigasyon menüsünde gizlenip gizlenmeyeceğini belirler.
	// true ise sayfa menüde görünmez, false ise görünür.
	HideInNavigation bool
}

// Slug, ayarlar sayfasının URL'de kullanılan benzersiz tanımlayıcısını döndürür.
//
// # Dönüş Değeri
// "settings" - Sayfanın URL slug'ı
//
// # Kullanım
// Bu metot, sayfa yönlendirmesi ve URL oluşturma işlemlerinde kullanılır.
// Örneğin: /admin/settings
//
// # Örnek
// ```go
// settings := &Settings{}
// slug := settings.Slug() // "settings"
// ```
func (p *Settings) Slug() string {
	return "settings"
}

// Title, ayarlar sayfasının başlığını döndürür.
//
// # Dönüş Değeri
// "Settings" - Sayfanın görüntülenecek başlığı
//
// # Kullanım
// Bu başlık, sayfa başlığı, tarayıcı sekmesi ve navigasyon menüsünde gösterilir.
//
// # Örnek
// ```go
// settings := &Settings{}
// title := settings.Title() // "Settings"
// ```
func (p *Settings) Title() string {
	return "Settings"
}

// Description, ayarlar sayfasının açıklamasını döndürür.
//
// # Dönüş Değeri
// "Sistem ayarlarını yönetin" - Sayfanın açıklaması
//
// # Kullanım
// Bu açıklama, sayfa hakkında bilgi vermek için kullanıcı arayüzünde gösterilir.
// Genellikle sayfa başlığının altında veya navigasyon menüsünde tooltip olarak görünür.
//
// # Örnek
// ```go
// settings := &Settings{}
// desc := settings.Description() // "Sistem ayarlarını yönetin"
// ```
func (p *Settings) Description() string {
	return "Sistem ayarlarını yönetin"
}

// Group, ayarlar sayfasının ait olduğu grup/kategorisini döndürür.
//
// # Dönüş Değeri
// "System" - Sayfanın ait olduğu grup adı
//
// # Kullanım
// Navigasyon menüsünde sayfaları gruplandırmak için kullanılır.
// Aynı grup adına sahip sayfalar menüde birlikte gösterilir.
//
// # Örnek
// ```go
// settings := &Settings{}
// group := settings.Group() // "System"
// // Menüde "System" başlığı altında gösterilir
// ```
func (p *Settings) Group() string {
	return "System"
}

// NavigationOrder, navigasyon menüsünde sayfanın gösterilme sırasını belirler.
//
// # Dönüş Değeri
// 100 - Sıra numarası (daha yüksek sayı, menünün daha aşağısında gösterilir)
//
// # Kullanım
// Navigasyon menüsünde sayfaları sıralamak için kullanılır.
// Sistem sayfaları genellikle menünün en altında gösterilir.
//
// # Sıra Kuralları
// - Düşük sayılar (0-50): Menünün üst kısmında
// - Orta sayılar (50-100): Menünün orta kısmında
// - Yüksek sayılar (100+): Menünün alt kısmında
//
// # Örnek
// ```go
// settings := &Settings{}
// order := settings.NavigationOrder() // 100
// // Sistem sayfaları menünün en altında gösterilir
// ```
func (p *Settings) NavigationOrder() int {
	return 100 // Sistem öğeleri genellikle menünün en altında gösterilir
}

// Visible, ayarlar sayfasının navigasyon menüsünde görünür olup olmadığını belirler.
//
// # Dönüş Değeri
// bool - true ise sayfa menüde görünür, false ise gizlidir
//
// # Mantık
// HideInNavigation alanının ters değerini döndürür.
// - HideInNavigation = false → Visible = true (görünür)
// - HideInNavigation = true → Visible = false (gizli)
//
// # Kullanım
// Navigasyon menüsü oluşturulurken, hangi sayfaların gösterilip gösterilmeyeceğini
// belirlemek için kullanılır.
//
// # Örnek
// ```go
// settings := &Settings{HideInNavigation: false}
// visible := settings.Visible() // true
//
// settings2 := &Settings{HideInNavigation: true}
// visible2 := settings2.Visible() // false
// ```
func (p *Settings) Visible() bool {
	return !p.HideInNavigation
}

// Cards, ayarlar sayfasında gösterilecek widget kartlarını döndürür.
//
// # Dönüş Değeri
// []widget.Card - Boş bir kart listesi
//
// # Kullanım
// Ayarlar sayfasında istatistik, grafik veya özet bilgiler göstermek için
// kullanılabilecek kartları tanımlar. Şu anda boş bir liste döndürülmektedir.
//
// # Gelecek Geliştirmeler
// Ayarlar sayfasında sistem istatistikleri veya özet bilgiler göstermek için
// bu metot genişletilebilir.
//
// # Örnek
// ```go
// settings := &Settings{}
// cards := settings.Cards() // []widget.Card{} (boş)
// ```
func (p *Settings) Cards() []widget.Card {
	return []widget.Card{}
}

// Fields, ayarlar sayfasında gösterilecek form alanlarını döndürür.
//
// # Dönüş Değeri
// []fields.Element - Sayfada gösterilecek form alanlarının listesi
//
// # Kullanım
// Ayarlar sayfasının form alanlarını dinamik olarak sağlamak için kullanılır.
// Her Element, bir form alanını temsil eder (TextInput, Select, Checkbox, vb.)
//
// # Örnek
// ```go
// settings := &Settings{
//     Elements: []fields.Element{
//         // Form alanları
//     },
// }
// fields := settings.Fields() // Elements döndürülür
// ```
func (p *Settings) Fields() []fields.Element {
	return p.Elements
}

// Save, ayarlar sayfasından gelen verileri veritabanına kaydeder.
//
// # Parametreler
// - c (*context.Context): İstek bağlamı, kullanıcı ve oturum bilgilerini içerir
// - db (*gorm.DB): Veritabanı bağlantısı, GORM ORM aracılığıyla
// - data (map[string]interface{}): Kaydedilecek ayarlar (anahtar-değer çiftleri)
//
// # Dönüş Değeri
// error - İşlem başarılı ise nil, hata ise error nesnesi
//
// # Çalışma Mantığı
// 1. Gelen veri haritasındaki her anahtar-değer çiftini döngüye alır
// 2. Değeri string'e dönüştürür:
//    - Eğer zaten string ise olduğu gibi kullanır
//    - Diğer türler için JSON serileştirmesi yapar
// 3. Veritabanında aynı anahtarla bir ayar varsa günceller, yoksa yeni oluşturur
// 4. GORM'un OnConflict özelliğini kullanarak upsert işlemi yapar
//
// # Kullanım Senaryoları
// 1. Sistem ayarlarını kaydetmek
// 2. Kullanıcı tercihlerini güncellemek
// 3. Konfigürasyon değerlerini değiştirmek
//
// # Örnek Kullanım
// ```go
// settings := &Settings{}
// err := settings.Save(ctx, db, map[string]interface{}{
//     "site_name": "Benim Sitesi",
//     "site_url": "https://example.com",
//     "items_per_page": 20,
//     "theme_config": map[string]string{"color": "blue"},
// })
// if err != nil {
//     log.Printf("Ayarlar kaydedilemedi: %v", err)
// }
// ```
//
// # Veri Dönüştürme Örnekleri
// ```
// String değer:
//   "site_name" → "Benim Sitesi"
//
// Sayı değer:
//   "items_per_page" → "20" (JSON string olarak)
//
// Kompleks değer:
//   "theme_config" → "{\"color\":\"blue\"}" (JSON string olarak)
// ```
//
// # Veritabanı İşlemi
// - Anahtar benzersiz (unique) olmalıdır
// - Aynı anahtarla yeni bir kayıt gelirse, eski değer güncellenir
// - updated_at alanı otomatik olarak güncellenir
//
// # Avantajlar
// - Farklı veri tiplerini destekler
// - Veritabanında merkezi olarak saklar
// - Aynı anahtarla yeni değer gelirse otomatik günceller
// - Hata yönetimi yapılır
//
// # Dezavantajlar
// - Kompleks veri yapıları JSON'a dönüştürülmeli
// - Veri doğrulaması bu metotta yapılmaz
// - JSON serileştirmesi başarısız olursa sessizce yoksayılır
//
// # Önemli Notlar
// - Tüm değerler string olarak veritabanında saklanır
// - JSON serileştirmesi başarısız olursa hata yoksayılır (b, _ := json.Marshal)
// - Veritabanı hataları döndürülür
// - Bağlam (context) şu anda kullanılmamaktadır, gelecekte audit log için kullanılabilir
//
// # Uyarılar
// - Büyük veri yapıları JSON serileştirmesi sırasında bellek tüketebilir
// - Veritabanı bağlantısı açık olmalıdır
// - Hata durumunda işlem kısmen tamamlanmış olabilir (atomik değildir)
func (p *Settings) Save(c *context.Context, db *gorm.DB, data map[string]interface{}) error {
	for key, value := range data {
		// Değeri string'e dönüştür
		var strValue string
		if v, ok := value.(string); ok {
			// Zaten string ise olduğu gibi kullan
			strValue = v
		} else {
			// String olmayan değerler için JSON serileştirmesi yap
			b, _ := json.Marshal(value)
			strValue = string(b)
		}

		// Ayar nesnesi oluştur
		s := setting.Setting{
			Key:   key,
			Value: strValue,
		}

		// Veritabanına kaydet (upsert işlemi)
		// Aynı anahtarla bir kayıt varsa güncelle, yoksa yeni oluştur
		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "key"}},
			DoUpdates: clause.AssignmentColumns([]string{"value", "updated_at"}),
		}).Create(&s).Error; err != nil {
			return err
		}
	}
	return nil
}
