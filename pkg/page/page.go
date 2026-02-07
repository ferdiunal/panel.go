// Package page, panel uygulamasında özel sayfaları (Dashboard, Ayarlar vb.) tanımlamak ve yönetmek için kullanılır.
//
// # Genel Bakış
//
// Bu paket, panel uygulamasında dinamik olarak özel sayfalar oluşturmak için bir arayüz ve temel yapı sağlar.
// Sayfalar, menüde görüntülenmek, form alanları içermek, veri kaydetmek ve erişim kontrolü yapmak gibi
// özelliklere sahip olabilir.
//
// # Temel Kavramlar
//
// - **Page Interface**: Tüm özel sayfaların uyması gereken sözleşmeyi tanımlar
// - **Base Struct**: Page arayüzünü implemente eden ve embedding için kullanılabilen temel yapı
// - **Slug**: Sayfanın URL'deki benzersiz tanımlayıcısı
// - **Navigation**: Menüde görüntülenme, sıralama ve görünürlük ayarları
// - **Access Control**: Kullanıcı izinlerine göre sayfa erişim kontrolü
//
// # Kullanım Senaryoları
//
// 1. **Dashboard Sayfası**: Sistem istatistiklerini ve widget'ları gösteren ana sayfa
// 2. **Ayarlar Sayfası**: Uygulama konfigürasyonlarını yönetmek için form alanları içeren sayfa
// 3. **Raporlama Sayfası**: Özel raporlar ve grafikler gösteren sayfa
// 4. **Yönetim Sayfası**: Sistem yönetimi için özel işlevler içeren sayfa
//
// # Örnek Kullanım
//
//	type DashboardPage struct {
//		page.Base
//	}
//
//	func (d *DashboardPage) Slug() string {
//		return "dashboard"
//	}
//
//	func (d *DashboardPage) Title() string {
//		return "Kontrol Paneli"
//	}
//
//	func (d *DashboardPage) Cards() []widget.Card {
//		return []widget.Card{
//			// Widget'ları buraya ekleyin
//		}
//	}
package page

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/widget"
	"gorm.io/gorm"
)

// Page, panel uygulamasında özel sayfaları tanımlamak için kullanılan arayüzdür.
//
// # Açıklama
//
// Page arayüzü, panel uygulamasında Dashboard, Ayarlar, Raporlar gibi özel sayfaları
// oluşturmak için gerekli tüm metodları tanımlar. Her sayfa, bu arayüzü implemente
// etmek zorundadır.
//
// # Metodlar
//
// Arayüz, sayfanın temel özelliklerini (başlık, açıklama, ikon vb.), görüntülenme
// ayarlarını (menüde görünürlük, sıralama, grup) ve işlevselliğini (form alanları,
// veri kaydetme, erişim kontrolü) tanımlar.
//
// # Avantajlar
//
// - Tutarlı sayfa yapısı sağlar
// - Menü yönetimini otomatikleştirir
// - Erişim kontrolü mekanizması sunar
// - Form alanlarını dinamik olarak yönetir
// - Widget ve kartları destekler
//
// # Dezavantajlar
//
// - Tüm metodları implemente etmek zorunludur (Base struct kullanılmadığı takdirde)
// - Karmaşık sayfalar için çok sayıda metod gerekebilir
//
// # Önemli Notlar
//
// - Slug değeri benzersiz olmalıdır (URL'de kullanılır)
// - CanAccess metodu her istek için çağrılır, performans önemlidir
// - Save metodu veri doğrulaması yapmalıdır
// - NavigationOrder değeri küçük olan sayfalar menüde önce görünür
//
// # Uyarılar
//
// - Slug değeri URL-safe karakterler içermelidir
// - CanAccess metodu nil context ile çağrılabilir, kontrol edin
// - Save metodu transaction içinde çalışabilir, dikkatli olun
type Page interface {
	// Slug, sayfanın URL'deki benzersiz tanımlayıcısıdır.
	//
	// # Açıklama
	//
	// Slug, sayfaya erişmek için URL'de kullanılan benzersiz bir tanımlayıcıdır.
	// Örneğin, "dashboard" slug'ı "/panel/dashboard" URL'sine karşılık gelir.
	//
	// # Dönüş Değeri
	//
	// Sayfanın URL-safe slug değeri. Boş string dönerse sayfa erişilemez.
	//
	// # Örnek
	//
	//	func (d *DashboardPage) Slug() string {
	//		return "dashboard"
	//	}
	//
	// # Önemli Notlar
	//
	// - Slug değeri benzersiz olmalıdır
	// - Slug değeri URL-safe karakterler içermelidir (a-z, 0-9, -)
	// - Slug değeri küçük harfle yazılmalıdır
	// - Boş string dönerse sayfa erişilemez
	Slug() string

	// Title, menüde ve sayfada görünecek başlıktır.
	//
	// # Açıklama
	//
	// Title, sayfanın insan tarafından okunabilir adıdır. Menüde ve sayfa başlığında
	// görüntülenir. Kullanıcı arayüzünde görünen ana başlıktır.
	//
	// # Dönüş Değeri
	//
	// Sayfanın başlık metni. Boş string dönerse "Başlıksız Sayfa" gibi varsayılan
	// bir başlık kullanılabilir.
	//
	// # Örnek
	//
	//	func (d *DashboardPage) Title() string {
	//		return "Kontrol Paneli"
	//	}
	//
	// # Önemli Notlar
	//
	// - Başlık kısa ve açıklayıcı olmalıdır
	// - Başlık çok uzun olmamalıdır (50 karakterden az önerilir)
	// - Başlık Türkçe karakterler içerebilir
	Title() string

	// Description, sayfanın açıklamasıdır.
	//
	// # Açıklama
	//
	// Description, sayfanın ne yaptığını ve amacını açıklayan bir metindir.
	// Menüde tooltip olarak veya sayfa hakkında bilgi olarak görüntülenebilir.
	//
	// # Dönüş Değeri
	//
	// Sayfanın açıklama metni. Boş string dönerse açıklama gösterilmez.
	//
	// # Örnek
	//
	//	func (d *DashboardPage) Description() string {
	//		return "Sistem istatistiklerini ve önemli metrikleri görüntüleyin"
	//	}
	//
	// # Önemli Notlar
	//
	// - Açıklama kısa ve öz olmalıdır
	// - Açıklama 200 karakterden az olmalıdır
	// - Açıklama sayfanın amacını açıkça belirtmelidir
	Description() string

	// Cards, sayfada gösterilecek widget kartlarını döner.
	//
	// # Açıklama
	//
	// Cards, sayfada görüntülenecek widget kartlarının listesini döner.
	// Kartlar, istatistikler, grafikler, tablolar vb. içerebilir.
	//
	// # Dönüş Değeri
	//
	// widget.Card arayüzünü implemente eden kartların dilimi. Boş dilim dönerse
	// sayfada kart gösterilmez.
	//
	// # Örnek
	//
	//	func (d *DashboardPage) Cards() []widget.Card {
	//		return []widget.Card{
	//			&StatisticCard{Title: "Toplam Kullanıcı", Value: 1234},
	//			&ChartCard{Title: "Aylık Satışlar", Data: chartData},
	//		}
	//	}
	//
	// # Önemli Notlar
	//
	// - Kartlar sayfada sırayla gösterilir
	// - Çok fazla kart performansı etkileyebilir
	// - Kartlar nil olmamalıdır
	// - Kartlar asenkron olarak yüklenebilir
	Cards() []widget.Card

	// Fields, sayfada gösterilecek form alanlarını döner.
	//
	// # Açıklama
	//
	// Fields, sayfada görüntülenecek form alanlarının listesini döner.
	// Alanlar, metin girişi, seçim kutusu, tarih seçici vb. olabilir.
	//
	// # Dönüş Değeri
	//
	// fields.Element arayüzünü implemente eden alanların dilimi. Boş dilim dönerse
	// sayfada form alanı gösterilmez.
	//
	// # Örnek
	//
	//	func (s *SettingsPage) Fields() []fields.Element {
	//		return []fields.Element{
	//			&fields.Text{Name: "app_name", Label: "Uygulama Adı"},
	//			&fields.Email{Name: "admin_email", Label: "Yönetici E-postası"},
	//		}
	//	}
	//
	// # Önemli Notlar
	//
	// - Alanlar sayfada sırayla gösterilir
	// - Her alan benzersiz bir Name değerine sahip olmalıdır
	// - Alanlar nil olmamalıdır
	// - Alanlar doğrulama kuralları içerebilir
	Fields() []fields.Element

	// Save, sayfa formundan gelen verileri işler ve kaydeder.
	//
	// # Açıklama
	//
	// Save metodu, kullanıcı tarafından form aracılığıyla gönderilen verileri
	// işler ve kaydeder. Veri doğrulaması, dönüştürme ve veritabanı işlemleri
	// burada yapılır.
	//
	// # Parametreler
	//
	// - c: İstek bağlamı (Context), kullanıcı bilgisi ve diğer bağlam verilerini içerir
	// - db: GORM veritabanı bağlantısı, veri tabanı işlemleri için kullanılır
	// - data: Form verilerinin harita gösterimi, alanların Name değerleri anahtar olarak kullanılır
	//
	// # Dönüş Değeri
	//
	// Hata durumunda error döner, başarılı olursa nil döner.
	//
	// # Örnek
	//
	//	func (s *SettingsPage) Save(c *context.Context, db *gorm.DB, data map[string]any) error {
	//		appName, ok := data["app_name"].(string)
	//		if !ok {
	//			return fmt.Errorf("geçersiz uygulama adı")
	//		}
	//
	//		// Veri doğrulaması
	//		if appName == "" {
	//			return fmt.Errorf("uygulama adı boş olamaz")
	//		}
	//
	//		// Veritabanına kaydet
	//		return db.Model(&Setting{}).Where("key = ?", "app_name").Update("value", appName).Error
	//	}
	//
	// # Önemli Notlar
	//
	// - Veri doğrulaması yapmalıdır
	// - Hata mesajları açıklayıcı olmalıdır
	// - Transaction içinde çalışabilir
	// - Veritabanı hataları döndürülmelidir
	// - Güvenlik kontrolleri yapılmalıdır
	//
	// # Uyarılar
	//
	// - Context nil olabilir, kontrol edin
	// - data haritası nil olabilir
	// - Veritabanı bağlantısı kapalı olabilir
	// - Kullanıcı izinleri kontrol edilmelidir
	Save(c *context.Context, db *gorm.DB, data map[string]any) error

	// Icon, menüde görünecek ikon adıdır.
	//
	// # Açıklama
	//
	// Icon, menüde sayfanın yanında görüntülenecek ikon adıdır. Genellikle
	// Font Awesome veya Material Icons gibi ikon kütüphanelerinin ikon adları kullanılır.
	//
	// # Dönüş Değeri
	//
	// Ikon adı (örn: "dashboard", "settings", "chart-bar"). Boş string dönerse
	// varsayılan ikon kullanılır.
	//
	// # Örnek
	//
	//	func (d *DashboardPage) Icon() string {
	//		return "dashboard"
	//	}
	//
	//	func (s *SettingsPage) Icon() string {
	//		return "cog"
	//	}
	//
	// # Önemli Notlar
	//
	// - Ikon adı ikon kütüphanesi tarafından desteklenmelidir
	// - Ikon adı küçük harfle yazılmalıdır
	// - Ikon adı kısa olmalıdır
	Icon() string

	// Group, menüde hangi grup altında görüneceğini belirler.
	//
	// # Açıklama
	//
	// Group, menüde sayfanın hangi kategori altında görüneceğini belirler.
	// Örneğin, "Yönetim", "Raporlar", "Ayarlar" gibi gruplar oluşturabilirsiniz.
	//
	// # Dönüş Değeri
	//
	// Grup adı. Boş string dönerse varsayılan grup ("Genel") kullanılır.
	//
	// # Örnek
	//
	//	func (u *UserManagementPage) Group() string {
	//		return "Yönetim"
	//	}
	//
	//	func (r *ReportPage) Group() string {
	//		return "Raporlar"
	//	}
	//
	// # Önemli Notlar
	//
	// - Grup adı kısa ve açıklayıcı olmalıdır
	// - Aynı grup adına sahip sayfalar birlikte görüntülenir
	// - Grup adı Türkçe karakterler içerebilir
	Group() string

	// NavigationOrder, menüdeki sıralama önceliğini döner.
	//
	// # Açıklama
	//
	// NavigationOrder, menüde sayfaların görüntülenme sırasını belirler.
	// Küçük değerler menüde daha önce görünür.
	//
	// # Dönüş Değeri
	//
	// Sıralama önceliği (0-100 arası önerilir). Küçük değer = daha yüksek öncelik.
	//
	// # Örnek
	//
	//	func (d *DashboardPage) NavigationOrder() int {
	//		return 1  // Menüde ilk sırada görünür
	//	}
	//
	//	func (s *SettingsPage) NavigationOrder() int {
	//		return 99  // Menüde son sırada görünür
	//	}
	//
	// # Önemli Notlar
	//
	// - Değer ne kadar küçükse, menüde o kadar önce görünür
	// - Aynı değere sahip sayfalar alfabetik sırayla görüntülenir
	// - Negatif değerler kullanılabilir
	// - Varsayılan değer 99'dur
	NavigationOrder() int

	// Visible, sayfanın menüde görünüp görünmeyeceğini belirler.
	//
	// # Açıklama
	//
	// Visible, sayfanın menüde görüntülenip görüntülenmeyeceğini kontrol eder.
	// Gizli sayfalar menüde görünmez ancak doğrudan URL aracılığıyla erişilebilir.
	//
	// # Dönüş Değeri
	//
	// true dönerse sayfa menüde görünür, false dönerse gizlenir.
	//
	// # Örnek
	//
	//	func (d *DashboardPage) Visible() bool {
	//		return true  // Menüde görünür
	//	}
	//
	//	func (h *HiddenPage) Visible() bool {
	//		return false  // Menüde gizlenir
	//	}
	//
	// # Önemli Notlar
	//
	// - Gizli sayfalar menüde görünmez ancak erişilebilir
	// - Dinamik olarak görünürlük değiştirebilirsiniz
	// - Varsayılan değer true'dur
	Visible() bool

	// CanAccess, kullanıcının sayfaya erişip erişemeyeceğini belirler.
	//
	// # Açıklama
	//
	// CanAccess, kullanıcının sayfaya erişme izni olup olmadığını kontrol eder.
	// Rol tabanlı erişim kontrolü (RBAC) veya diğer yetkilendirme mekanizmaları
	// burada uygulanır.
	//
	// # Parametreler
	//
	// - c: İstek bağlamı (Context), kullanıcı bilgisi ve izinleri içerir
	//
	// # Dönüş Değeri
	//
	// true dönerse kullanıcı sayfaya erişebilir, false dönerse erişim reddedilir.
	//
	// # Örnek
	//
	//	func (a *AdminPage) CanAccess(c *context.Context) bool {
	//		// Sadece yöneticiler erişebilir
	//		return c.User != nil && c.User.IsAdmin
	//	}
	//
	//	func (d *DashboardPage) CanAccess(c *context.Context) bool {
	//		// Tüm oturum açmış kullanıcılar erişebilir
	//		return c.User != nil
	//	}
	//
	// # Önemli Notlar
	//
	// - Context nil olabilir, kontrol edin
	// - Performans önemlidir (her istek için çağrılır)
	// - Veritabanı sorgusu yapmaktan kaçının
	// - Varsayılan değer true'dur
	//
	// # Uyarılar
	//
	// - Context nil olabilir
	// - Kullanıcı bilgisi eksik olabilir
	// - Performans kritiktir
	CanAccess(c *context.Context) bool
}

// Base, Page arayüzünü implemente eden temel yapı ve embedding için kullanılabilir.
//
// # Açıklama
//
// Base struct, Page arayüzünün tüm metodlarının varsayılan uygulamalarını sağlar.
// Bu struct, embedding (gömme) yoluyla kullanılarak, özel sayfa türlerinin sadece
// gerekli metodları override etmesine olanak tanır. Böylece kod tekrarı azalır ve
// bakım kolaylaşır.
//
// # Kullanım Senaryoları
//
// 1. **Hızlı Prototip Oluşturma**: Temel bir sayfa hızlıca oluşturmak için
// 2. **Miras Yoluyla Genişletme**: Base'i embed ederek sadece gerekli metodları override etmek
// 3. **Varsayılan Davranış**: Tüm metodlar için makul varsayılan değerler sağlamak
//
// # Varsayılan Değerler
//
// - Slug: "" (boş string)
// - Title: "" (boş string)
// - Description: "" (boş string)
// - Icon: "circle" (varsayılan ikon)
// - Group: "Genel" (varsayılan grup)
// - NavigationOrder: 99 (menüde son sırada)
// - Visible: true (menüde görünür)
// - CanAccess: true (tüm kullanıcılar erişebilir)
// - Cards: [] (boş dilim)
// - Fields: [] (boş dilim)
// - Save: nil (hata döndürmez)
//
// # Örnek Kullanım
//
//	type DashboardPage struct {
//		page.Base  // Base'i embed et
//	}
//
//	// Sadece gerekli metodları override et
//	func (d *DashboardPage) Slug() string {
//		return "dashboard"
//	}
//
//	func (d *DashboardPage) Title() string {
//		return "Kontrol Paneli"
//	}
//
//	func (d *DashboardPage) Cards() []widget.Card {
//		return []widget.Card{
//			// Widget'ları buraya ekle
//		}
//	}
//	// Diğer metodlar Base'den kalıtılır
//
// # Avantajlar
//
// - Kod tekrarını azaltır
// - Hızlı prototip oluşturmaya olanak tanır
// - Tutarlı varsayılan davranış sağlar
// - Embedding yoluyla esnek genişletme imkanı
//
// # Dezavantajlar
//
// - Tüm metodlar override edilmediğinde beklenmedik davranışlar olabilir
// - Varsayılan değerler her zaman uygun olmayabilir
//
// # Önemli Notlar
//
// - Base struct'ı doğrudan kullanmayın, her zaman embed edin
// - Override etmediğiniz metodlar varsayılan değerleri döner
// - Slug değeri boş kalırsa sayfa erişilemez
// - Title değeri boş kalırsa menüde görüntülenme sorunları olabilir
type Base struct {
}

// Slug, Base struct için varsayılan slug değerini döner.
//
// # Açıklama
//
// Base struct'ın Slug metodu boş string döner. Özel sayfa türleri bu metodu
// override ederek kendi slug değerlerini belirtmelidir.
//
// # Dönüş Değeri
//
// Boş string (""). Sayfa erişilemez hale gelir.
//
// # Örnek
//
//	type MyPage struct {
//		page.Base
//	}
//
//	func (m *MyPage) Slug() string {
//		return "my-page"  // Base'in varsayılan değerini override et
//	}
//
// # Önemli Notlar
//
// - Base'in Slug metodu boş string döner
// - Özel sayfa türleri bu metodu override etmelidir
// - Slug değeri benzersiz olmalıdır
func (b Base) Slug() string {
	return ""
}

// Title, Base struct için varsayılan başlık değerini döner.
//
// # Açıklama
//
// Base struct'ın Title metodu boş string döner. Özel sayfa türleri bu metodu
// override ederek kendi başlık değerlerini belirtmelidir.
//
// # Dönüş Değeri
//
// Boş string (""). Menüde başlık görüntülenmez.
//
// # Örnek
//
//	type MyPage struct {
//		page.Base
//	}
//
//	func (m *MyPage) Title() string {
//		return "Benim Sayfam"  // Base'in varsayılan değerini override et
//	}
//
// # Önemli Notlar
//
// - Base'in Title metodu boş string döner
// - Özel sayfa türleri bu metodu override etmelidir
// - Başlık kısa ve açıklayıcı olmalıdır
func (b Base) Title() string {
	return ""
}

// Description, Base struct için varsayılan açıklama değerini döner.
//
// # Açıklama
//
// Base struct'ın Description metodu boş string döner. Özel sayfa türleri bu metodu
// override ederek kendi açıklama değerlerini belirtmelidir.
//
// # Dönüş Değeri
//
// Boş string (""). Açıklama görüntülenmez.
//
// # Örnek
//
//	type MyPage struct {
//		page.Base
//	}
//
//	func (m *MyPage) Description() string {
//		return "Bu sayfa benim özel sayfamdır"  // Base'in varsayılan değerini override et
//	}
//
// # Önemli Notlar
//
// - Base'in Description metodu boş string döner
// - Özel sayfa türleri bu metodu override etmelidir
// - Açıklama kısa ve öz olmalıdır
func (b Base) Description() string {
	return ""
}

// Icon, Base struct için varsayılan ikon adını döner.
//
// # Açıklama
//
// Base struct'ın Icon metodu "circle" ikon adını döner. Bu, tüm sayfalar için
// tutarlı bir varsayılan ikon sağlar. Özel sayfa türleri bu metodu override
// ederek kendi ikon adlarını belirtebilir.
//
// # Dönüş Değeri
//
// "circle" ikon adı. Menüde bu ikon görüntülenir.
//
// # Örnek
//
//	type DashboardPage struct {
//		page.Base
//	}
//
//	func (d *DashboardPage) Icon() string {
//		return "dashboard"  // Base'in varsayılan değerini override et
//	}
//
// # Önemli Notlar
//
// - Base'in Icon metodu "circle" döner
// - Ikon adı ikon kütüphanesi tarafından desteklenmelidir
// - Varsayılan ikon tüm sayfalar için uygun olabilir
func (b Base) Icon() string {
	return "circle" // Varsayılan ikon
}

// Group, Base struct için varsayılan grup adını döner.
//
// # Açıklama
//
// Base struct'ın Group metodu "Genel" grup adını döner. Bu, tüm sayfalar için
// tutarlı bir varsayılan grup sağlar. Özel sayfa türleri bu metodu override
// ederek kendi grup adlarını belirtebilir.
//
// # Dönüş Değeri
//
// "Genel" grup adı. Menüde bu grup altında görüntülenir.
//
// # Örnek
//
//	type AdminPage struct {
//		page.Base
//	}
//
//	func (a *AdminPage) Group() string {
//		return "Yönetim"  // Base'in varsayılan değerini override et
//	}
//
// # Önemli Notlar
//
// - Base'in Group metodu "Genel" döner
// - Aynı grup adına sahip sayfalar birlikte görüntülenir
// - Varsayılan grup tüm sayfalar için uygun olabilir
func (b Base) Group() string {
	return "Genel" // Varsayılan grup
}

// NavigationOrder, Base struct için varsayılan sıralama önceliğini döner.
//
// # Açıklama
//
// Base struct'ın NavigationOrder metodu 99 değerini döner. Bu, sayfaların menüde
// son sırada görüntülenmesini sağlar. Özel sayfa türleri bu metodu override
// ederek kendi sıralama önceliklerini belirtebilir.
//
// # Dönüş Değeri
//
// 99 (sıralama önceliği). Menüde son sırada görüntülenir.
//
// # Örnek
//
//	type DashboardPage struct {
//		page.Base
//	}
//
//	func (d *DashboardPage) NavigationOrder() int {
//		return 1  // Base'in varsayılan değerini override et, menüde ilk sırada görünür
//	}
//
// # Önemli Notlar
//
// - Base'in NavigationOrder metodu 99 döner
// - Küçük değer = daha yüksek öncelik (menüde daha önce görünür)
// - Varsayılan değer 99, menüde son sırada görüntülenir
func (b Base) NavigationOrder() int {
	return 99
}

// Visible, Base struct için varsayılan görünürlük değerini döner.
//
// # Açıklama
//
// Base struct'ın Visible metodu true döner. Bu, tüm sayfaların varsayılan olarak
// menüde görünür olmasını sağlar. Özel sayfa türleri bu metodu override ederek
// sayfaları gizleyebilir.
//
// # Dönüş Değeri
//
// true (görünür). Sayfa menüde görüntülenir.
//
// # Örnek
//
//	type HiddenPage struct {
//		page.Base
//	}
//
//	func (h *HiddenPage) Visible() bool {
//		return false  // Base'in varsayılan değerini override et, menüde gizle
//	}
//
// # Önemli Notlar
//
// - Base'in Visible metodu true döner
// - Gizli sayfalar menüde görünmez ancak erişilebilir
// - Varsayılan davranış sayfaları görünür yapar
func (b Base) Visible() bool {
	return true
}

// CanAccess, Base struct için varsayılan erişim kontrolü değerini döner.
//
// # Açıklama
//
// Base struct'ın CanAccess metodu true döner. Bu, tüm kullanıcıların varsayılan
// olarak sayfaya erişebilmesini sağlar. Özel sayfa türleri bu metodu override
// ederek erişim kontrolü uygulayabilir.
//
// # Parametreler
//
// - c: İstek bağlamı (Context). Kullanıcı bilgisi ve izinleri içerir.
//
// # Dönüş Değeri
//
// true (erişim izni var). Tüm kullanıcılar sayfaya erişebilir.
//
// # Örnek
//
//	type AdminPage struct {
//		page.Base
//	}
//
//	func (a *AdminPage) CanAccess(c *context.Context) bool {
//		// Base'in varsayılan değerini override et
//		return c.User != nil && c.User.IsAdmin
//	}
//
// # Önemli Notlar
//
// - Base'in CanAccess metodu true döner
// - Context nil olabilir, kontrol edin
// - Varsayılan davranış tüm kullanıcılara erişim izni verir
// - Özel sayfa türleri erişim kontrolü uygulamalıdır
func (b Base) CanAccess(c *context.Context) bool {
	return true // Varsayılan: erişim izni ver
}

// Cards, Base struct için varsayılan kartlar dilimini döner.
//
// # Açıklama
//
// Base struct'ın Cards metodu boş bir dilim döner. Özel sayfa türleri bu metodu
// override ederek kendi kartlarını belirtebilir.
//
// # Dönüş Değeri
//
// Boş widget.Card dilimi. Sayfada kart görüntülenmez.
//
// # Örnek
//
//	type DashboardPage struct {
//		page.Base
//	}
//
//	func (d *DashboardPage) Cards() []widget.Card {
//		return []widget.Card{
//			&StatisticCard{Title: "Toplam Kullanıcı", Value: 1234},
//		}
//	}
//
// # Önemli Notlar
//
// - Base'in Cards metodu boş dilim döner
// - Özel sayfa türleri bu metodu override etmelidir
// - Kartlar nil olmamalıdır
func (b Base) Cards() []widget.Card {
	return []widget.Card{}
}

// Fields, Base struct için varsayılan form alanları dilimini döner.
//
// # Açıklama
//
// Base struct'ın Fields metodu boş bir dilim döner. Özel sayfa türleri bu metodu
// override ederek kendi form alanlarını belirtebilir.
//
// # Dönüş Değeri
//
// Boş fields.Element dilimi. Sayfada form alanı görüntülenmez.
//
// # Örnek
//
//	type SettingsPage struct {
//		page.Base
//	}
//
//	func (s *SettingsPage) Fields() []fields.Element {
//		return []fields.Element{
//			&fields.Text{Name: "app_name", Label: "Uygulama Adı"},
//		}
//	}
//
// # Önemli Notlar
//
// - Base'in Fields metodu boş dilim döner
// - Özel sayfa türleri bu metodu override etmelidir
// - Alanlar nil olmamalıdır
func (b Base) Fields() []fields.Element {
	return []fields.Element{}
}

// Save, Base struct için varsayılan veri kaydetme işlemini gerçekleştirir.
//
// # Açıklama
//
// Base struct'ın Save metodu hiçbir işlem yapmaz ve nil döner. Özel sayfa türleri
// bu metodu override ederek veri kaydetme işlemini uygulayabilir.
//
// # Parametreler
//
// - c: İstek bağlamı (Context)
// - db: GORM veritabanı bağlantısı
// - data: Form verilerinin harita gösterimi
//
// # Dönüş Değeri
//
// nil (hata yok). Hiçbir işlem yapılmaz.
//
// # Örnek
//
//	type SettingsPage struct {
//		page.Base
//	}
//
//	func (s *SettingsPage) Save(c *context.Context, db *gorm.DB, data map[string]any) error {
//		// Base'in varsayılan davranışını override et
//		appName, ok := data["app_name"].(string)
//		if !ok {
//			return fmt.Errorf("geçersiz uygulama adı")
//		}
//
//		return db.Model(&Setting{}).Where("key = ?", "app_name").Update("value", appName).Error
//	}
//
// # Önemli Notlar
//
// - Base'in Save metodu hiçbir işlem yapmaz
// - Özel sayfa türleri bu metodu override etmelidir
// - Veri doğrulaması yapmalıdır
// - Hata mesajları açıklayıcı olmalıdır
func (b Base) Save(c *context.Context, db *gorm.DB, data map[string]any) error {
	return nil
}
