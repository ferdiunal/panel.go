package page

import (
	"github.com/ferdiunal/panel.go/pkg/domain/user"
	"github.com/ferdiunal/panel.go/pkg/widget"
)

// ============================================================================
// DASHBOARD SAYFASI - Ana Kontrol Paneli
// ============================================================================
//
// Dashboard struct'ı, panel uygulamasının ana kontrol panelini temsil eder.
// Sistem özeti, istatistikler ve önemli metrikleri görüntülemek için kullanılır.
//
// # Amaç ve Kullanım
//
// Dashboard, yönetici panelinin giriş noktası olarak hizmet eder. Sistem
// genelinde hızlı bir özet sunarak, kullanıcılara sistem durumu hakkında
// anlık bilgi sağlar.
//
// # Temel Özellikler
//
// - Özel slug ve başlık tanımlaması
// - Navigasyon menüsünde öncelikli konumlandırma (-1 sırası)
// - Dinamik widget kartları desteği
// - Sistem ikonları entegrasyonu
// - Genişletilebilir yapı
//
// # Kullanım Senaryoları
//
// 1. **Yönetici Panelinin Ana Sayfası**: Sistem açıldığında ilk görüntülenen sayfa
// 2. **Sistem İstatistiklerinin Gösterilmesi**: Toplam kullanıcı, aktivite vb.
// 3. **Hızlı Erişim Noktası**: Diğer sayfalara hızlı navigasyon
// 4. **Sistem Sağlığı Kontrolü**: Önemli metriklerin izlenmesi
// 5. **Veri Özeti**: Veritabanı istatistiklerinin merkezi gösterimi
//
// # Avantajlar
//
// - Basit ve anlaşılır yapı
// - Genişletilebilir widget sistemi
// - Merkezi yönetim noktası
// - Hızlı veri görüntüleme
// - Base struct'tan inherit edilerek kod tekrarı azaltılır
// - Diğer sayfalarla tutarlı arayüz
//
// # Dezavantajlar ve Sınırlamalar
//
// - Sabit widget yapısı (şu anda sadece kullanıcı sayısı)
// - Dinamik widget ekleme/kaldırma mekanizması eksik
// - Veri yenileme stratejisi tanımlanmamış
// - Caching mekanizması bulunmamaktadır
// - N+1 sorgu problemi oluşabilir
//
// # Önemli Notlar
//
// - Dashboard her zaman navigasyon menüsünün başında görünür (-1 sırası)
// - Base struct'tan inherit edilir, bu nedenle temel sayfa işlevselliğini sağlar
// - Widget kartları dinamik olarak oluşturulabilir
// - Performans için kart sayısı optimize edilmeli
//
// # Gelecek Geliştirmeler
//
// - Dinamik kart ekleme/kaldırma
// - Kart sıralaması özelleştirmesi
// - Kart görünürlüğü kontrolleri
// - Veri caching stratejisi
// - Gerçek zamanlı veri güncellemeleri
// - Kullanıcı tercihlerine göre özelleştirme
//
// # Örnek Kullanım
//
// ```go
// dashboard := &Dashboard{}
// slug := dashboard.Slug()           // "dashboard"
// title := dashboard.Title()         // "Dashboard"
// desc := dashboard.Description()    // "Sistem özeti ve istatistikleri"
// icon := dashboard.Icon()           // "layout-dashboard"
// order := dashboard.NavigationOrder() // -1
// cards := dashboard.Cards()         // Widget kartları
// ```
//
type Dashboard struct {
	Base
}

// ============================================================================
// Slug() - Sayfa Tanımlayıcısı
// ============================================================================
//
// Dashboard sayfasının URL'de kullanılan benzersiz tanımlayıcısını döndürür.
// Bu değer, sayfa yönlendirmesinde ve URL oluşturmada kullanılır.
//
// # Dönüş Değeri
//
// - `string`: "dashboard" - Sayfanın URL slug'ı
//
// # Teknik Detaylar
//
// - Slug, URL'de sayfa tanımlaması için kullanılır
// - Sistem genelinde benzersiz olmalıdır
// - Küçük harfler ve tire (-) içerebilir
// - Boşluk ve özel karakterler içermemelidir
//
// # Kullanım Örneği
//
// ```go
// dashboard := &Dashboard{}
// slug := dashboard.Slug()
// // Çıktı: "dashboard"
//
// // URL oluşturmada kullanılır
// // URL: /admin/dashboard
// // URL: /panel/dashboard
// ```
//
// # Önemli Notlar
//
// - Bu değer değiştirilmesi mevcut URL'leri kıracaktır
// - Eski URL'lere yönlendirme (redirect) eklenmelidir
// - SEO açısından önemlidir
// - Diğer sayfalarla çakışmamalıdır
//
// # İlişkili Metodlar
//
// - Title(): Sayfa başlığı
// - Description(): Sayfa açıklaması
// - Icon(): Sayfa ikonu
//
func (d *Dashboard) Slug() string {
	return "dashboard"
}

// ============================================================================
// Title() - Sayfa Başlığı
// ============================================================================
//
// Dashboard sayfasının kullanıcı arayüzünde gösterilecek başlığını döndürür.
// Bu başlık, tarayıcı sekmesinde, sayfa başlığında ve navigasyon menüsünde
// görüntülenir.
//
// # Dönüş Değeri
//
// - `string`: "Dashboard" - Sayfanın görüntü başlığı
//
// # Teknik Detaylar
//
// - Başlık, HTML <title> etiketinde kullanılır
// - Navigasyon menüsünde gösterilir
// - Sayfa başlığı olarak render edilir
// - Tarayıcı sekmesinde görüntülenir
//
// # Kullanım Örneği
//
// ```go
// dashboard := &Dashboard{}
// title := dashboard.Title()
// // Çıktı: "Dashboard"
//
// // HTML'de kullanılır
// // <title>Dashboard - Panel</title>
// // <h1>Dashboard</h1>
// ```
//
// # Yerelleştirme (i18n)
//
// Gelecekte yerelleştirme için uygun bir yerdir:
//
// ```go
// func (d *Dashboard) Title() string {
//     return i18n.Translate("dashboard.title") // "Dashboard"
// }
// ```
//
// # Önemli Notlar
//
// - Kullanıcı arayüzünde görüntülenen metindir
// - Sayfa başlığı ve navigasyon menüsünde kullanılır
// - Kısa ve açıklayıcı olmalıdır
// - Maksimum 50-60 karakter önerilir
//
// # İlişkili Metodlar
//
// - Slug(): Sayfa tanımlayıcısı
// - Description(): Sayfa açıklaması
// - Icon(): Sayfa ikonu
//
func (d *Dashboard) Title() string {
	return "Dashboard"
}

// ============================================================================
// Description() - Sayfa Açıklaması
// ============================================================================
//
// Dashboard sayfasının kısa açıklamasını döndürür. Bu açıklama, navigasyon
// menüsünde tooltip olarak veya sayfa bilgisinde gösterilir. Kullanıcılara
// sayfanın amacını ve içeriğini açıklar.
//
// # Dönüş Değeri
//
// - `string`: "Sistem özeti ve istatistikleri" - Sayfanın açıklaması
//
// # Teknik Detaylar
//
// - Açıklama, navigasyon menüsünde hover tooltip olarak gösterilir
// - Sayfa bilgisinde kullanılabilir
// - Erişilebilirlik (accessibility) için önemlidir
// - Screen reader'lar tarafından okunabilir
//
// # Kullanım Örneği
//
// ```go
// dashboard := &Dashboard{}
// desc := dashboard.Description()
// // Çıktı: "Sistem özeti ve istatistikleri"
//
// // HTML'de kullanılır
// // <div title="Sistem özeti ve istatistikleri">Dashboard</div>
// // <p class="description">Sistem özeti ve istatistikleri</p>
// ```
//
// # Erişilebilirlik (Accessibility)
//
// - Screen reader'lar tarafından okunabilir
// - ARIA açıklamaları için kullanılabilir
// - Kullanıcıların sayfayı anlamasına yardımcı olur
//
// # Önemli Notlar
//
// - Türkçe olarak yazılmıştır
// - Kısa ve açıklayıcı olmalıdır
// - Maksimum 100-150 karakter önerilir
// - Sayfanın gerçek amacını yansıtmalıdır
//
// # İlişkili Metodlar
//
// - Title(): Sayfa başlığı
// - Slug(): Sayfa tanımlayıcısı
// - Icon(): Sayfa ikonu
//
func (d *Dashboard) Description() string {
	return "Sistem özeti ve istatistikleri"
}

// ============================================================================
// Icon() - Sayfa İkonu
// ============================================================================
//
// Dashboard sayfasının navigasyon menüsünde ve arayüzde gösterilecek ikon
// adını döndürür. İkon, sayfanın görsel tanımlaması için kullanılır ve
// kullanıcı deneyimini iyileştirir.
//
// # Dönüş Değeri
//
// - `string`: "layout-dashboard" - İkon tanımlayıcısı
//
// # Teknik Detaylar
//
// - İkon adı, kullanılan ikon kütüphanesinde aranır
// - Genellikle Feather Icons veya Material Icons gibi kütüphaneler kullanılır
// - İkon, SVG veya font-based olabilir
// - Responsive tasarımda ölçeklenebilir
//
// # Desteklenen İkonlar
//
// - "layout-dashboard" - Kontrol paneli ikonu (mevcut)
// - Diğer ikonlar ikon kütüphanesine bağlıdır
//
// # Kullanım Örneği
//
// ```go
// dashboard := &Dashboard{}
// icon := dashboard.Icon()
// // Çıktı: "layout-dashboard"
//
// // HTML'de kullanılır
// // <i class="icon-layout-dashboard"></i>
// // <svg class="icon"><!-- layout-dashboard SVG --></svg>
// ```
//
// # İkon Kütüphaneleri
//
// Feather Icons örneği:
// ```html
// <svg class="icon">
//   <use xlink:href="#icon-layout-dashboard"></use>
// </svg>
// ```
//
// Material Icons örneği:
// ```html
// <i class="material-icons">dashboard</i>
// ```
//
// # Önemli Notlar
//
// - İkon adı, kullanılan ikon kütüphanesinde mevcut olmalıdır
// - Navigasyon menüsünde görsel tanımlama sağlar
// - Kullanıcı deneyimini iyileştirir
// - Erişilebilirlik için alt text eklenmelidir
// - İkon boyutu responsive olmalıdır
//
// # İlişkili Metodlar
//
// - Title(): Sayfa başlığı
// - Description(): Sayfa açıklaması
// - Slug(): Sayfa tanımlayıcısı
//
func (d *Dashboard) Icon() string {
	return "layout-dashboard"
}

// ============================================================================
// NavigationOrder() - Navigasyon Sırası
// ============================================================================
//
// Dashboard sayfasının navigasyon menüsünde gösterilme sırasını belirler.
// Daha düşük değerler, menünün başında görünür. Negatif değerler, sayfayı
// menünün en başına yerleştirir.
//
// # Dönüş Değeri
//
// - `int`: -1 - Navigasyon sırası (negatif değer = başa yerleştir)
//
// # Teknik Detaylar
//
// - Sıralama, tüm sayfalar arasında yapılır
// - Aynı değere sahip sayfalar alfabetik sıralanır
// - Negatif değerler pozitif değerlerden önce gelir
// - Dinamik olarak değiştirilebilir
//
// # Sıralama Kuralları
//
// - -1: Menünün en başında (Dashboard)
// - 0: Varsayılan sıra
// - 1, 2, 3...: Sonraya yerleştir
// - Aynı değer: Alfabetik sıralama
//
// # Kullanım Örneği
//
// ```go
// dashboard := &Dashboard{}
// order := dashboard.NavigationOrder()
// // Çıktı: -1
//
// // Menü sıralaması örneği:
// // Dashboard (-1)
// // Users (0)
// // Settings (1)
// // Reports (2)
// ```
//
// # Menü Sıralama Algoritması
//
// ```
// Sayfalar sıralanırken:
// 1. NavigationOrder() değerine göre sırala (küçükten büyüğe)
// 2. Aynı değere sahip sayfaları Title() alfabetik sırala
// 3. Sonuç menüde göster
// ```
//
// # Önemli Notlar
//
// - Dashboard her zaman menünün başında görünür
// - Bu, kullanıcıların hızlı erişim sağlamasını amaçlar
// - Diğer sayfalar bu değere göre sıralanır
// - Dinamik olarak değiştirilebilir (gerekirse)
// - Menü yeniden yüklendiğinde sıralama güncellenir
//
// # İlişkili Metodlar
//
// - Title(): Sayfa başlığı
// - Slug(): Sayfa tanımlayıcısı
// - Icon(): Sayfa ikonu
//
func (d *Dashboard) NavigationOrder() int {
	return -1
}

// ============================================================================
// Cards() - Sayfa Kartları
// ============================================================================
//
// Dashboard sayfasında gösterilecek widget kartlarını döndürür. Her kart,
// sistem istatistiklerini veya önemli bilgileri görüntüler. Kartlar dinamik
// olarak oluşturulabilir ve çeşitli widget türlerini destekler.
//
// # Dönüş Değeri
//
// - `[]widget.Card`: Widget kartları dizisi
//   - Boş dizi: Kart yok (sayfa boş görünür)
//   - Dolu dizi: Gösterilecek kartlar
//
// # Mevcut Kartlar
//
// 1. **Total Users (Toplam Kullanıcılar)**
//    - Tür: CountWidget
//    - Model: user.User
//    - Gösterir: Sistemdeki toplam kullanıcı sayısı
//    - Veri Kaynağı: Veritabanı sorgusu
//    - Güncelleme: Her sayfa yüklemesinde
//
// # Desteklenen Widget Türleri
//
// - `CountWidget`: Sayı gösterimi (örn: toplam kullanıcı)
// - `ChartWidget`: Grafik gösterimi (örn: trend analizi)
// - `StatWidget`: İstatistik gösterimi (örn: yüzde)
// - `TableWidget`: Tablo gösterimi (örn: son aktiviteler)
// - Özel widget'lar: Genişletilebilir yapı
//
// # Kullanım Örneği
//
// ```go
// dashboard := &Dashboard{}
// cards := dashboard.Cards()
// // Çıktı: [CountWidget("Total Users", user.User)]
//
// // Kartları render etme
// for _, card := range cards {
//     fmt.Println(card.Title())
//     // Çıktı: "Total Users"
// }
// ```
//
// # Kart Ekleme Örneği
//
// ```go
// func (d *Dashboard) Cards() []widget.Card {
//     return []widget.Card{
//         widget.NewCountWidget("Total Users", &user.User{}),
//         widget.NewCountWidget("Active Sessions", &session.Session{}),
//         widget.NewChartWidget("User Growth", chartData),
//         widget.NewStatWidget("System Health", 95),
//     }
// }
// ```
//
// # Performans Notları
//
// - Kartlar her sayfa yüklemesinde oluşturulur
// - Veri tabanı sorguları kartlar tarafından yapılır
// - Büyük veri setleri için pagination önerilir
// - Caching mekanizması eklenebilir
// - N+1 sorgu problemi oluşabilir
//
// # Optimizasyon Önerileri
//
// 1. **Caching**: Sık değişmeyen veriler cache'lenebilir
// 2. **Pagination**: Büyük veri setleri sayfalanabilir
// 3. **Lazy Loading**: Kartlar gerektiğinde yüklenebilir
// 4. **Batch Queries**: Birden fazla sorgu birleştirilebilir
// 5. **Asynchronous Loading**: Kartlar asenkron yüklenebilir
//
// # Önemli Uyarılar
//
// - Kart sayısı performansı etkileyebilir
// - Her kart ayrı bir veri tabanı sorgusu yapabilir
// - N+1 sorgu problemi oluşabilir
// - Veri yenileme sıklığı optimize edilmeli
// - Veritabanı bağlantı havuzu yeterli olmalıdır
//
// # Gelecek Geliştirmeler
//
// - Dinamik kart ekleme/kaldırma
// - Kart sıralaması özelleştirmesi
// - Kart görünürlüğü kontrolleri
// - Veri caching stratejisi
// - Gerçek zamanlı veri güncellemeleri
// - Kullanıcı tercihlerine göre özelleştirme
// - Kart boyutu ayarlaması
// - Kart tema seçimi
//
// # Veri Tabanı Sorgusu Örneği
//
// ```go
// // CountWidget, arka planda şu sorguyu çalıştırır:
// // SELECT COUNT(*) FROM users
//
// // Sonuç: "Total Users: 1,234"
// ```
//
// # İlişkili Metodlar
//
// - Title(): Sayfa başlığı
// - Description(): Sayfa açıklaması
// - Slug(): Sayfa tanımlayıcısı
// - Icon(): Sayfa ikonu
//
func (d *Dashboard) Cards() []widget.Card {
	return []widget.Card{
		widget.NewCountWidget("Total Users", &user.User{}),
	}
}
