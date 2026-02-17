package core

// ElementType, bir alanın (field) tipini temsil eden string tabanlı bir tür tanımıdır.
//
// Bu tür, panel sisteminde bir alanın nasıl render edileceğini ve verileri nasıl işleyeceğini belirler.
// Her ElementType, frontend'de farklı bir UI bileşenine karşılık gelir ve farklı veri tipleriyle çalışır.
//
// # Kullanım Senaryoları
//
// - Form alanlarının tipini belirlemek için kullanılır
// - Liste görünümlerinde sütun tiplerini tanımlamak için kullanılır
// - Detay sayfalarında alan render tipini belirlemek için kullanılır
// - İlişkisel alanların (relationship) görünüm tipini ayarlamak için kullanılır
//
// # Örnek Kullanım
//
// ```go
// field := fields.Text("title").
//
//	SetType(core.TYPE_TEXT)
//
// // Zengin metin editörü için
// field := fields.Text("content").
//
//	SetType(core.TYPE_RICHTEXT)
//
// // İlişki alanı için
// field := fields.BelongsTo("user").
//
//	SetType(core.TYPE_LINK)
//
// ```
//
// # Önemli Notlar
//
// - ElementType değerleri sabit olarak tanımlanmıştır (TYPE_TEXT, TYPE_TEXTAREA, vb.)
// - Her alan tipi, belirli veri tipleriyle uyumludur
// - İlişkisel alanlar için özel tipler mevcuttur (TYPE_LINK, TYPE_COLLECTION, vb.)
// - Polymorphic ilişkiler için TYPE_POLY_* önekli tipler kullanılır
//
// # İlgili Dokümantasyon
//
// Daha fazla bilgi için bakınız:
// - docs/Fields.md - Alan tipleri ve kullanımları
// - docs/Relationships.md - İlişkisel alan tipleri
type ElementType string

// ElementContext, bir alanın hangi bağlamda (context) görüntüleneceğini belirleyen string tabanlı bir tür tanımıdır.
//
// Bu tür, alanların farklı görünümlerde (liste, detay, form) görünürlüğünü ve davranışını kontrol eder.
// Context değerleri, alanların hangi sayfalarda gösterileceğini veya gizleneceğini belirler.
//
// # Kullanım Senaryoları
//
// - Bir alanı sadece form sayfalarında göstermek
// - Bir alanı liste görünümünden gizlemek
// - Bir alanı sadece oluşturma (create) formunda göstermek
// - Bir alanı güncelleme (update) formunda gizlemek
// - Bir alanı detay sayfasında göstermek
//
// # Context Kategorileri
//
// 1. **Temel Context'ler**: CONTEXT_FORM, CONTEXT_DETAIL, CONTEXT_LIST
// 2. **Show Context'leri**: SHOW_ON_FORM, SHOW_ON_DETAIL, SHOW_ON_LIST
// 3. **Hide Context'leri**: HIDE_ON_LIST, HIDE_ON_DETAIL, HIDE_ON_CREATE, HIDE_ON_UPDATE
// 4. **Only Context'leri**: ONLY_ON_LIST, ONLY_ON_DETAIL, ONLY_ON_CREATE, ONLY_ON_UPDATE, ONLY_ON_FORM
//
// # Örnek Kullanım
//
// ```go
// // Sadece liste görünümünde göster
// field := fields.Text("title").
//
//	ShowOn(core.ONLY_ON_LIST)
//
// // Oluşturma formunda gizle (genellikle ID gibi alanlar için)
// field := fields.ID("id").
//
//	HideOn(core.HIDE_ON_CREATE)
//
// // Sadece detay sayfasında göster
// field := fields.Text("created_at").
//
//	ShowOn(core.ONLY_ON_DETAIL)
//
// // Birden fazla context kullanımı
// field := fields.Text("internal_notes").
//
//	HideOn(core.HIDE_ON_LIST).
//	ShowOn(core.SHOW_ON_DETAIL)
//
// ```
//
// # Avantajlar
//
// - Alanların görünürlüğünü granüler şekilde kontrol edebilme
// - Farklı kullanıcı deneyimleri için aynı alanı farklı şekillerde gösterebilme
// - Performans optimizasyonu (gereksiz alanları gizleyerek)
// - Güvenlik (hassas alanları belirli görünümlerden gizleyerek)
//
// # Önemli Notlar
//
// - ONLY_* context'leri, alanı sadece belirtilen görünümde gösterir
// - HIDE_* context'leri, alanı belirtilen görünümden gizler
// - SHOW_* context'leri, alanı belirtilen görünümde gösterir
// - Birden fazla context aynı anda kullanılabilir
// - Context'ler birbirleriyle çakışabilir, bu durumda son eklenen geçerli olur
//
// # İlgili Dokümantasyon
//
// Daha fazla bilgi için bakınız:
// - docs/Fields.md - Alan görünürlük kontrolleri
type ElementContext string

// VisibilityContext, bir alanın hangi UI bağlamında görünür olması gerektiğini belirleyen string tabanlı bir tür tanımıdır.
//
// Bu tür, ElementContext'ten daha spesifik bir görünürlük kontrolü sağlar ve
// belirli UI durumlarını (index, detail, create, update, preview) temsil eder.
//
// # Kullanım Senaryoları
//
// - Index/liste sayfasında alan görünürlüğünü kontrol etmek
// - Detay sayfasında alan görünürlüğünü kontrol etmek
// - Oluşturma formunda alan görünürlüğünü kontrol etmek
// - Güncelleme formunda alan görünürlüğünü kontrol etmek
// - Önizleme modunda alan görünürlüğünü kontrol etmek
//
// # Context Tipleri
//
// - **ContextIndex**: Liste/index görünümü (tablo görünümü)
// - **ContextDetail**: Detay görünümü (tek kayıt görünümü)
// - **ContextCreate**: Oluşturma formu
// - **ContextUpdate**: Güncelleme formu
// - **ContextPreview**: Önizleme modu
//
// # Örnek Kullanım
//
// ```go
// // Visibility kontrolü ile alan tanımlama
// field := fields.Text("title").
//
//	SetVisibleOn(core.ContextIndex, core.ContextDetail)
//
// // Sadece oluşturma formunda görünür
// field := fields.Text("initial_value").
//
//	SetVisibleOn(core.ContextCreate)
//
// // Güncelleme ve detay sayfasında görünür
// field := fields.Text("updated_info").
//
//	SetVisibleOn(core.ContextUpdate, core.ContextDetail)
//
// // Önizleme modunda gizle
// field := fields.Text("internal_data").
//
//	SetHiddenOn(core.ContextPreview)
//
// ```
//
// # ElementContext ile Farkları
//
// | Özellik | VisibilityContext | ElementContext |
// |---------|-------------------|----------------|
// | Granülerlik | Daha spesifik | Daha genel |
// | Kullanım | Programatik kontrol | Deklaratif kontrol |
// | Esneklik | Yüksek | Orta |
// | Karmaşıklık | Düşük | Orta |
//
// # Avantajlar
//
// - Daha temiz ve anlaşılır kod
// - Spesifik UI durumları için optimize edilmiş
// - Tip güvenliği (string sabitleri)
// - Kolay test edilebilirlik
//
// # Önemli Notlar
//
// - VisibilityContext değerleri sabit olarak tanımlanmıştır
// - Birden fazla context aynı anda kullanılabilir
// - ElementContext ile birlikte kullanılabilir
// - Frontend tarafında bu context'lere göre render yapılır
//
// # İlgili Dokümantasyon
//
// Daha fazla bilgi için bakınız:
// - docs/Fields.md - Alan görünürlük yönetimi
type VisibilityContext string

// Element type constants - Alan tipi sabitleri
//
// Bu sabitler, panel sisteminde kullanılabilecek tüm alan tiplerini tanımlar.
// Her sabit, frontend'de farklı bir UI bileşenine karşılık gelir.
//
// # Alan Tipi Kategorileri
//
// 1. **Metin Alanları**: TYPE_TEXT, TYPE_TEXTAREA, TYPE_RICHTEXT, TYPE_PASSWORD
// 2. **Sayısal Alanlar**: TYPE_NUMBER
// 3. **İletişim Alanları**: TYPE_TEL, TYPE_EMAIL
// 4. **Medya Alanları**: TYPE_AUDIO, TYPE_VIDEO, TYPE_FILE
// 5. **Tarih/Zaman Alanları**: TYPE_DATE, TYPE_DATETIME
// 6. **Yapılandırılmış Alanlar**: TYPE_KEY_VALUE
// 7. **İlişkisel Alanlar**: TYPE_LINK, TYPE_COLLECTION, TYPE_DETAIL, TYPE_CONNECT
// 8. **Polymorphic Alanlar**: TYPE_POLY_LINK, TYPE_POLY_DETAIL, TYPE_POLY_COLLECTION, TYPE_POLY_CONNECT
// 9. **Seçim Alanları**: TYPE_BOOLEAN, TYPE_SELECT, TYPE_BOOLEAN_GROUP
// 10. **Görsel Alanlar**: TYPE_BADGE, TYPE_CODE, TYPE_COLOR
// 11. **Konteyner Alanlar**: TYPE_PANEL
// 12. **Genel İlişki Alanları**: TYPE_RELATIONSHIP
//
// # Önemli Notlar
//
// - İlişkisel alanlar için docs/Relationships.md dosyasına bakınız
// - Her alan tipi, belirli veri tipleriyle uyumludur
// - Polymorphic alanlar, birden fazla model tipine bağlanabilir
const (
	// TYPE_TEXT, tek satırlık metin girişi için kullanılan alan tipidir.
	//
	// Bu alan tipi, kısa metin verilerini (başlık, isim, slug vb.) saklamak için kullanılır.
	// HTML input[type="text"] elementine karşılık gelir.
	//
	// # Kullanım Senaryoları
	//
	// - Başlık (title) alanları
	// - İsim (name) alanları
	// - Slug alanları
	// - Kısa açıklama alanları
	// - URL path alanları
	//
	// # Örnek Kullanım
	//
	// ```go
	// fields.Text("title").
	//     SetLabel("Başlık").
	//     SetPlaceholder("Makale başlığını girin").
	//     SetRequired(true).
	//     SetMaxLength(255)
	// ```
	//
	// # Özellikler
	//
	// - Tek satırlık giriş
	// - Maksimum karakter sınırı ayarlanabilir
	// - Placeholder desteği
	// - Validasyon kuralları eklenebilir
	// - Otomatik trim işlemi
	TYPE_TEXT ElementType = "text"

	// TYPE_TEXTAREA, çok satırlı metin girişi için kullanılan alan tipidir.
	//
	// Bu alan tipi, uzun metin verilerini (açıklama, içerik, not vb.) saklamak için kullanılır.
	// HTML textarea elementine karşılık gelir.
	//
	// # Kullanım Senaryoları
	//
	// - Açıklama (description) alanları
	// - Not (notes) alanları
	// - Adres alanları
	// - Kısa içerik alanları
	// - Meta description alanları
	//
	// # Örnek Kullanım
	//
	// ```go
	// fields.Textarea("description").
	//     SetLabel("Açıklama").
	//     SetPlaceholder("Ürün açıklamasını girin").
	//     SetRows(5).
	//     SetMaxLength(1000)
	// ```
	//
	// # Özellikler
	//
	// - Çok satırlı giriş
	// - Satır sayısı ayarlanabilir
	// - Otomatik yükseklik ayarı
	// - Karakter sayacı gösterilebilir
	TYPE_TEXTAREA ElementType = "textarea"

	// TYPE_RICHTEXT, zengin metin editörü (WYSIWYG) için kullanılan alan tipidir.
	//
	// Bu alan tipi, HTML formatında zengin içerik oluşturmak için kullanılır.
	// Metin biçimlendirme, resim ekleme, link oluşturma gibi özellikler sunar.
	//
	// # Kullanım Senaryoları
	//
	// - Blog yazısı içeriği
	// - Sayfa içeriği (CMS)
	// - Ürün detaylı açıklaması
	// - E-posta şablonları
	// - Duyuru içeriği
	//
	// # Örnek Kullanım
	//
	// ```go
	// fields.Richtext("content").
	//     SetLabel("İçerik").
	//     SetRequired(true).
	//     SetToolbar([]string{"bold", "italic", "link", "image"})
	// ```
	//
	// # Özellikler
	//
	// - WYSIWYG editör (TinyMCE, Quill, vb.)
	// - HTML çıktısı
	// - Resim yükleme desteği
	// - Link ekleme
	// - Metin biçimlendirme (bold, italic, underline)
	// - Liste oluşturma
	// - Tablo ekleme
	//
	// # Önemli Notlar
	//
	// - XSS saldırılarına karşı HTML sanitization yapılmalıdır
	// - Büyük içerikler için performans optimizasyonu gerekebilir
	TYPE_RICHTEXT ElementType = "richtext"

	// TYPE_PASSWORD, şifre girişi için kullanılan alan tipidir.
	//
	// Bu alan tipi, hassas bilgileri (şifre, API key vb.) gizli şekilde almak için kullanılır.
	// HTML input[type="password"] elementine karşılık gelir.
	//
	// # Kullanım Senaryoları
	//
	// - Kullanıcı şifresi
	// - API anahtarları
	// - Token değerleri
	// - Gizli yapılandırma değerleri
	//
	// # Örnek Kullanım
	//
	// ```go
	// fields.Password("password").
	//     SetLabel("Şifre").
	//     SetRequired(true).
	//     SetMinLength(8).
	//     SetHelp("En az 8 karakter olmalıdır")
	// ```
	//
	// # Özellikler
	//
	// - Girilen değer gizlenir (*** şeklinde gösterilir)
	// - Şifre gücü göstergesi eklenebilir
	// - Göster/Gizle toggle'ı eklenebilir
	// - Validasyon kuralları (min length, complexity)
	//
	// # Güvenlik Notları
	//
	// - Şifreler veritabanında hash'lenerek saklanmalıdır
	// - HTTPS kullanımı zorunludur
	// - Şifre politikaları uygulanmalıdır
	TYPE_PASSWORD ElementType = "password"

	// TYPE_NUMBER, sayısal değer girişi için kullanılan alan tipidir.
	//
	// Bu alan tipi, tam sayı veya ondalıklı sayı değerlerini almak için kullanılır.
	// HTML input[type="number"] elementine karşılık gelir.
	//
	// # Kullanım Senaryoları
	//
	// - Fiyat alanları
	// - Miktar alanları
	// - Yaş alanları
	// - Sıralama (order) alanları
	// - Yüzde değerleri
	//
	// # Örnek Kullanım
	//
	// ```go
	// fields.Number("price").
	//     SetLabel("Fiyat").
	//     SetMin(0).
	//     SetMax(999999).
	//     SetStep(0.01).
	//     SetPrefix("₺")
	// ```
	//
	// # Özellikler
	//
	// - Minimum/maksimum değer sınırı
	// - Adım (step) değeri ayarlanabilir
	// - Prefix/suffix eklenebilir (₺, $, %, vb.)
	// - Ondalık basamak sayısı ayarlanabilir
	// - Artır/azalt butonları
	TYPE_NUMBER ElementType = "number"

	// TYPE_MONEY, para tutarı girişi/gösterimi için kullanılan alan tipidir.
	//
	// Bu alan tipi, sayısal bir değeri bir para birimi (currency) ile birlikte
	// formda toplamak ve list/detail ekranlarında locale-aware olarak göstermek
	// için kullanılır.
	TYPE_MONEY ElementType = "money"

	// TYPE_TEL, telefon numarası girişi için kullanılan alan tipidir.
	//
	// Bu alan tipi, telefon numaralarını formatlanmış şekilde almak için kullanılır.
	// HTML input[type="tel"] elementine karşılık gelir.
	//
	// # Kullanım Senaryoları
	//
	// - Kullanıcı telefon numarası
	// - İletişim telefonu
	// - Acil durum telefonu
	// - Faks numarası
	//
	// # Örnek Kullanım
	//
	// ```go
	// fields.Tel("phone").
	//     SetLabel("Telefon").
	//     SetPlaceholder("(5XX) XXX XX XX").
	//     SetMask("(999) 999 99 99")
	// ```
	//
	// # Özellikler
	//
	// - Otomatik formatlama
	// - Maske (mask) desteği
	// - Ülke kodu seçici eklenebilir
	// - Validasyon (telefon numarası formatı)
	TYPE_TEL ElementType = "tel"

	// TYPE_EMAIL, e-posta adresi girişi için kullanılan alan tipidir.
	//
	// Bu alan tipi, e-posta adreslerini doğrulanmış şekilde almak için kullanılır.
	// HTML input[type="email"] elementine karşılık gelir.
	//
	// # Kullanım Senaryoları
	//
	// - Kullanıcı e-posta adresi
	// - İletişim e-postası
	// - Bildirim e-postası
	// - CC/BCC e-postaları
	//
	// # Örnek Kullanım
	//
	// ```go
	// fields.Email("email").
	//     SetLabel("E-posta").
	//     SetRequired(true).
	//     SetUnique(true).
	//     SetPlaceholder("ornek@email.com")
	// ```
	//
	// # Özellikler
	//
	// - Otomatik e-posta validasyonu
	// - Küçük harfe çevirme (lowercase)
	// - Trim işlemi
	// - Unique kontrolü
	TYPE_EMAIL ElementType = "email"

	// TYPE_AUDIO, ses dosyası yükleme için kullanılan alan tipidir.
	//
	// Bu alan tipi, ses dosyalarını (MP3, WAV, OGG vb.) yüklemek için kullanılır.
	// Ses önizleme ve oynatma özellikleri sunar.
	//
	// # Kullanım Senaryoları
	//
	// - Podcast dosyaları
	// - Müzik dosyaları
	// - Ses kayıtları
	// - Sesli mesajlar
	// - Ses efektleri
	//
	// # Örnek Kullanım
	//
	// ```go
	// fields.Audio("podcast").
	//     SetLabel("Podcast Dosyası").
	//     SetAccept([]string{".mp3", ".wav", ".ogg"}).
	//     SetMaxSize(50 * 1024 * 1024) // 50MB
	// ```
	//
	// # Özellikler
	//
	// - Ses önizleme (audio player)
	// - Dosya boyutu kontrolü
	// - Format kontrolü (MP3, WAV, OGG, vb.)
	// - Sürükle-bırak yükleme
	// - İlerleme çubuğu
	TYPE_AUDIO ElementType = "audio"

	// TYPE_VIDEO, video dosyası yükleme için kullanılan alan tipidir.
	//
	// Bu alan tipi, video dosyalarını (MP4, WebM, AVI vb.) yüklemek için kullanılır.
	// Video önizleme ve oynatma özellikleri sunar.
	//
	// # Kullanım Senaryoları
	//
	// - Ürün tanıtım videoları
	// - Eğitim videoları
	// - Reklam videoları
	// - Video içerikler
	// - Canlı yayın kayıtları
	//
	// # Örnek Kullanım
	//
	// ```go
	// fields.Video("promo_video").
	//     SetLabel("Tanıtım Videosu").
	//     SetAccept([]string{".mp4", ".webm"}).
	//     SetMaxSize(100 * 1024 * 1024) // 100MB
	// ```
	//
	// # Özellikler
	//
	// - Video önizleme (video player)
	// - Thumbnail oluşturma
	// - Dosya boyutu kontrolü
	// - Format kontrolü (MP4, WebM, AVI, vb.)
	// - Sürükle-bırak yükleme
	// - İlerleme çubuğu
	//
	// # Önemli Notlar
	//
	// - Büyük dosyalar için chunk upload kullanılmalıdır
	// - Video encoding/transcoding gerekebilir
	TYPE_VIDEO ElementType = "video"

	// TYPE_DATE, tarih seçici için kullanılan alan tipidir.
	//
	// Bu alan tipi, tarih değerlerini (gün/ay/yıl) seçmek için kullanılır.
	// Takvim (datepicker) bileşeni sunar.
	//
	// # Kullanım Senaryoları
	//
	// - Doğum tarihi
	// - Başlangıç/bitiş tarihleri
	// - Son kullanma tarihi
	// - Yayın tarihi
	// - Etkinlik tarihi
	//
	// # Örnek Kullanım
	//
	// ```go
	// fields.Date("birth_date").
	//     SetLabel("Doğum Tarihi").
	//     SetMin("1900-01-01").
	//     SetMax("2024-12-31").
	//     SetFormat("DD/MM/YYYY")
	// ```
	//
	// # Özellikler
	//
	// - Takvim (datepicker) bileşeni
	// - Minimum/maksimum tarih sınırı
	// - Tarih formatı özelleştirme
	// - Bugün butonu
	// - Hızlı tarih seçimi
	TYPE_DATE ElementType = "date"

	// TYPE_DATETIME, tarih ve saat seçici için kullanılan alan tipidir.
	//
	// Bu alan tipi, tarih ve saat değerlerini birlikte seçmek için kullanılır.
	// Takvim ve saat seçici bileşenleri sunar.
	//
	// # Kullanım Senaryoları
	//
	// - Etkinlik başlangıç/bitiş zamanı
	// - Randevu zamanı
	// - Yayın zamanı (publish date)
	// - Son güncelleme zamanı
	// - Zamanlı görevler
	//
	// # Örnek Kullanım
	//
	// ```go
	// fields.DateTime("published_at").
	//     SetLabel("Yayın Zamanı").
	//     SetFormat("DD/MM/YYYY HH:mm").
	//     SetTimezone("Europe/Istanbul")
	// ```
	//
	// # Özellikler
	//
	// - Takvim ve saat seçici
	// - Timezone desteği
	// - Format özelleştirme
	// - Şimdi butonu
	// - 12/24 saat formatı
	TYPE_DATETIME ElementType = "datetime"

	// TYPE_FILE, genel dosya yükleme için kullanılan alan tipidir.
	//
	// Bu alan tipi, herhangi bir dosya tipini yüklemek için kullanılır.
	// Resim, PDF, Word, Excel gibi tüm dosya tiplerini destekler.
	//
	// # Kullanım Senaryoları
	//
	// - Döküman yükleme (PDF, Word, Excel)
	// - Resim yükleme
	// - Arşiv dosyaları (ZIP, RAR)
	// - Sertifika dosyaları
	// - Ek dosyalar
	//
	// # Örnek Kullanım
	//
	// ```go
	// fields.File("attachment").
	//     SetLabel("Ek Dosya").
	//     SetAccept([]string{".pdf", ".doc", ".docx"}).
	//     SetMaxSize(10 * 1024 * 1024). // 10MB
	//     SetMultiple(true)
	// ```
	//
	// # Özellikler
	//
	// - Çoklu dosya yükleme
	// - Dosya tipi kısıtlama
	// - Dosya boyutu kontrolü
	// - Sürükle-bırak yükleme
	// - İlerleme çubuğu
	// - Önizleme (resimler için)
	//
	// # Güvenlik Notları
	//
	// - Dosya tipi validasyonu yapılmalıdır
	// - Dosya boyutu sınırı konulmalıdır
	// - Zararlı dosyalar taranmalıdır
	// - Dosya isimleri sanitize edilmelidir
	TYPE_FILE ElementType = "file"

	// TYPE_KEY_VALUE, anahtar-değer çifti girişi için kullanılan alan tipidir.
	//
	// Bu alan tipi, dinamik anahtar-değer çiftlerini saklamak için kullanılır.
	// JSON formatında veri saklar.
	//
	// # Kullanım Senaryoları
	//
	// - Meta veriler (metadata)
	// - Özel alanlar (custom fields)
	// - Yapılandırma ayarları
	// - HTTP headers
	// - Çeviri anahtarları
	//
	// # Örnek Kullanım
	//
	// ```go
	// fields.KeyValue("metadata").
	//     SetLabel("Meta Veriler").
	//     SetKeyLabel("Anahtar").
	//     SetValueLabel("Değer").
	//     SetAddButtonText("Yeni Ekle")
	// ```
	//
	// # Özellikler
	//
	// - Dinamik satır ekleme/çıkarma
	// - Anahtar-değer çifti girişi
	// - JSON çıktısı
	// - Sıralama (drag & drop)
	// - Validasyon (unique keys)
	//
	// # Veri Formatı
	//
	// ```json
	// {
	//   "color": "red",
	//   "size": "large",
	//   "weight": "500g"
	// }
	// ```
	TYPE_KEY_VALUE ElementType = "key_value"

	// TYPE_LINK, başka bir kaynağa (resource) bağlantı göstermek için kullanılan alan tipidir.
	//
	// Bu alan tipi, BelongsTo ilişkilerinde ilişkili kaydın linkini göstermek için kullanılır.
	// Tıklanabilir bir link olarak render edilir ve ilişkili kaydın detay sayfasına yönlendirir.
	//
	// # Kullanım Senaryoları
	//
	// - BelongsTo ilişkilerinde ilişkili kayıt linki
	// - HasOne ilişkilerinde ilişkili kayıt linki
	// - Liste görünümünde ilişkili kayıt gösterimi
	// - Detay sayfasında ilişkili kayıt referansı
	//
	// # Örnek Kullanım
	//
	// ```go
	// fields.BelongsTo("user").
	//     SetLabel("Kullanıcı").
	//     SetType(core.TYPE_LINK).
	//     SetDisplayField("name")
	// ```
	//
	// # Özellikler
	//
	// - Tıklanabilir link
	// - İlişkili kaydın detay sayfasına yönlendirme
	// - Özel display field belirleme
	// - Hover efekti
	//
	// # İlgili Dokümantasyon
	//
	// Daha fazla bilgi için bakınız:
	// - docs/Relationships.md - İlişki tipleri ve kullanımları
	TYPE_LINK ElementType = "link"

	// TYPE_COLLECTION, ilişkili kayıtların koleksiyonunu göstermek için kullanılan alan tipidir.
	//
	// Bu alan tipi, HasMany ve BelongsToMany ilişkilerinde ilişkili kayıtların listesini göstermek için kullanılır.
	// Genellikle detay sayfasında tablo formatında gösterilir.
	//
	// # Kullanım Senaryoları
	//
	// - HasMany ilişkilerinde alt kayıtların listesi
	// - BelongsToMany ilişkilerinde ilişkili kayıtların listesi
	// - Detay sayfasında ilişkili kayıtlar tablosu
	// - İç içe (nested) kayıt yönetimi
	//
	// # Örnek Kullanım
	//
	// ```go
	// fields.HasMany("posts").
	//     SetLabel("Yazılar").
	//     SetType(core.TYPE_COLLECTION).
	//     SetDisplayFields([]string{"title", "status", "created_at"})
	// ```
	//
	// # Özellikler
	//
	// - Tablo formatında gösterim
	// - Sayfalama (pagination) desteği
	// - Sıralama (sorting) desteği
	// - Filtreleme desteği
	// - Satır içi (inline) düzenleme
	// - Yeni kayıt ekleme
	//
	// # İlgili Dokümantasyon
	//
	// Daha fazla bilgi için bakınız:
	// - docs/Relationships.md - HasMany ve BelongsToMany kullanımları
	TYPE_COLLECTION ElementType = "collection"

	// TYPE_DETAIL, ilişkili kaydın detaylı görünümünü göstermek için kullanılan alan tipidir.
	//
	// Bu alan tipi, ilişkili kaydın tüm alanlarını detaylı şekilde göstermek için kullanılır.
	// Genellikle BelongsTo ve HasOne ilişkilerinde kullanılır.
	//
	// # Kullanım Senaryoları
	//
	// - BelongsTo ilişkisinde ilişkili kaydın detayları
	// - HasOne ilişkisinde ilişkili kaydın detayları
	// - İç içe (nested) kayıt görünümü
	// - Genişletilmiş bilgi gösterimi
	//
	// # Örnek Kullanım
	//
	// ```go
	// fields.BelongsTo("user").
	//     SetLabel("Kullanıcı Detayları").
	//     SetType(core.TYPE_DETAIL).
	//     SetDisplayFields([]string{"name", "email", "phone", "address"})
	// ```
	//
	// # Özellikler
	//
	// - Detaylı alan gösterimi
	// - Özelleştirilebilir alan listesi
	// - Genişletilmiş/daraltılmış görünüm
	// - İlişkili kaydın tüm bilgileri
	//
	// # İlgili Dokümantasyon
	//
	// Daha fazla bilgi için bakınız:
	// - docs/Relationships.md - İlişki detay görünümleri
	TYPE_DETAIL ElementType = "detail"

	// TYPE_CONNECT, başka bir kaynağa bağlantı oluşturmak için kullanılan alan tipidir.
	//
	// Bu alan tipi, ilişki kurmak için kayıt seçimi yapmayı sağlar.
	// Genellikle form sayfalarında (create/update) kullanılır.
	//
	// # Kullanım Senaryoları
	//
	// - BelongsTo ilişkisinde kayıt seçimi
	// - HasOne ilişkisinde kayıt seçimi
	// - BelongsToMany ilişkisinde çoklu kayıt seçimi
	// - İlişki oluşturma/güncelleme
	//
	// # Örnek Kullanım
	//
	// ```go
	// fields.BelongsTo("category").
	//     SetLabel("Kategori").
	//     SetType(core.TYPE_CONNECT).
	//     SetSearchable(true).
	//     SetDisplayField("name")
	// ```
	//
	// # Özellikler
	//
	// - Arama (search) desteği
	// - Dropdown/select bileşeni
	// - Autocomplete desteği
	// - Çoklu seçim (BelongsToMany için)
	// - Yeni kayıt oluşturma (inline create)
	//
	// # İlgili Dokümantasyon
	//
	// Daha fazla bilgi için bakınız:
	// - docs/Relationships.md - İlişki bağlantı yönetimi
	TYPE_CONNECT ElementType = "connect"

	// TYPE_POLY_LINK, polymorphic ilişkide başka bir kaynağa bağlantı göstermek için kullanılan alan tipidir.
	//
	// Bu alan tipi, MorphTo ilişkilerinde ilişkili kaydın linkini göstermek için kullanılır.
	// İlişkili kayıt farklı model tiplerinden olabilir.
	//
	// # Kullanım Senaryoları
	//
	// - MorphTo ilişkilerinde ilişkili kayıt linki
	// - Polymorphic ilişkilerde kayıt referansı
	// - Çok tipli (multi-type) ilişki gösterimi
	// - Yorumlar, beğeniler gibi polymorphic yapılar
	//
	// # Örnek Kullanım
	//
	// ```go
	// fields.MorphTo("commentable").
	//     SetLabel("Yorumlanan").
	//     SetType(core.TYPE_POLY_LINK).
	//     SetTypes([]string{"Post", "Video"})
	// ```
	//
	// # Özellikler
	//
	// - Çoklu model tipi desteği
	// - Dinamik link oluşturma
	// - Tip göstergesi (badge)
	// - Özel display field
	//
	// # İlgili Dokümantasyon
	//
	// Daha fazla bilgi için bakınız:
	// - docs/Relationships.md - Polymorphic ilişkiler
	TYPE_POLY_LINK ElementType = "poly_link"

	// TYPE_POLY_DETAIL, polymorphic ilişkide ilişkili kaydın detaylı görünümünü göstermek için kullanılan alan tipidir.
	//
	// Bu alan tipi, MorphTo ilişkilerinde ilişkili kaydın tüm alanlarını detaylı şekilde göstermek için kullanılır.
	// İlişkili kayıt farklı model tiplerinden olabilir ve her tip için farklı alanlar gösterilebilir.
	//
	// # Kullanım Senaryoları
	//
	// - MorphTo ilişkisinde ilişkili kaydın detayları
	// - Polymorphic ilişkilerde genişletilmiş bilgi
	// - Tip bazlı alan gösterimi
	// - Dinamik içerik gösterimi
	//
	// # Örnek Kullanım
	//
	// ```go
	// fields.MorphTo("taggable").
	//     SetLabel("Etiketlenen Detayları").
	//     SetType(core.TYPE_POLY_DETAIL).
	//     SetTypes([]string{"Post", "Video", "Product"})
	// ```
	//
	// # Özellikler
	//
	// - Tip bazlı alan gösterimi
	// - Dinamik detay render
	// - Genişletilmiş/daraltılmış görünüm
	// - Çoklu model tipi desteği
	//
	// # İlgili Dokümantasyon
	//
	// Daha fazla bilgi için bakınız:
	// - docs/Relationships.md - Polymorphic ilişki detayları
	TYPE_POLY_DETAIL ElementType = "poly_detail"

	// TYPE_POLY_COLLECTION, polymorphic ilişkide ilişkili kayıtların koleksiyonunu göstermek için kullanılan alan tipidir.
	//
	// Bu alan tipi, MorphMany ve MorphToMany ilişkilerinde ilişkili kayıtların listesini göstermek için kullanılır.
	// İlişkili kayıtlar farklı model tiplerinden olabilir.
	//
	// # Kullanım Senaryoları
	//
	// - MorphMany ilişkilerinde alt kayıtların listesi
	// - MorphToMany ilişkilerinde ilişkili kayıtların listesi
	// - Polymorphic koleksiyon yönetimi
	// - Yorumlar, etiketler gibi polymorphic listeler
	//
	// # Örnek Kullanım
	//
	// ```go
	// fields.MorphMany("comments").
	//     SetLabel("Yorumlar").
	//     SetType(core.TYPE_POLY_COLLECTION).
	//     SetDisplayFields([]string{"content", "user", "created_at"})
	// ```
	//
	// # Özellikler
	//
	// - Çoklu model tipi desteği
	// - Tablo formatında gösterim
	// - Tip filtreleme
	// - Sayfalama ve sıralama
	// - Satır içi düzenleme
	//
	// # İlgili Dokümantasyon
	//
	// Daha fazla bilgi için bakınız:
	// - docs/Relationships.md - Polymorphic koleksiyonlar
	TYPE_POLY_COLLECTION ElementType = "poly_collection"

	// TYPE_POLY_CONNECT, polymorphic ilişkide başka bir kaynağa bağlantı oluşturmak için kullanılan alan tipidir.
	//
	// Bu alan tipi, polymorphic ilişki kurmak için kayıt seçimi yapmayı sağlar.
	// Önce model tipi seçilir, sonra o tipten kayıt seçilir.
	//
	// # Kullanım Senaryoları
	//
	// - MorphTo ilişkisinde kayıt seçimi
	// - MorphToMany ilişkisinde çoklu kayıt seçimi
	// - Tip bazlı kayıt seçimi
	// - Polymorphic ilişki oluşturma
	//
	// # Örnek Kullanım
	//
	// ```go
	// fields.MorphTo("imageable").
	//     SetLabel("Resim Sahibi").
	//     SetType(core.TYPE_POLY_CONNECT).
	//     SetTypes([]string{"User", "Product", "Post"}).
	//     SetSearchable(true)
	// ```
	//
	// # Özellikler
	//
	// - İki aşamalı seçim (tip + kayıt)
	// - Tip bazlı arama
	// - Autocomplete desteği
	// - Çoklu seçim desteği
	// - Dinamik form alanları
	//
	// # İlgili Dokümantasyon
	//
	// Daha fazla bilgi için bakınız:
	// - docs/Relationships.md - Polymorphic ilişki bağlantıları
	TYPE_POLY_CONNECT ElementType = "poly_connect"

	// TYPE_BOOLEAN, boolean (true/false) değer girişi için kullanılan alan tipidir.
	//
	// Bu alan tipi, iki durumlu (açık/kapalı, evet/hayır) değerleri saklamak için kullanılır.
	// Checkbox veya toggle switch olarak render edilir.
	//
	// # Kullanım Senaryoları
	//
	// - Aktif/pasif durumu
	// - Yayında/taslak durumu
	// - Özellik açma/kapama (feature flags)
	// - Onay/red durumu
	// - Görünürlük kontrolü
	//
	// # Örnek Kullanım
	//
	// ```go
	// fields.Boolean("is_active").
	//     SetLabel("Aktif").
	//     SetDefault(true).
	//     SetHelp("Kaydın aktif olup olmadığını belirler")
	// ```
	//
	// # Özellikler
	//
	// - Checkbox veya toggle switch
	// - Varsayılan değer ayarlama
	// - True/false label özelleştirme
	// - Hızlı değiştirme (quick toggle)
	//
	// # Önemli Notlar
	//
	// - Veritabanında genellikle tinyint(1) veya boolean olarak saklanır
	// - Null değer alabilir (nullable)
	// - Varsayılan değer false'tur
	TYPE_BOOLEAN ElementType = "boolean"

	// TYPE_SELECT, açılır liste (dropdown) seçimi için kullanılan alan tipidir.
	//
	// Bu alan tipi, önceden tanımlanmış seçenekler arasından seçim yapmak için kullanılır.
	// Tek veya çoklu seçim yapılabilir.
	//
	// # Kullanım Senaryoları
	//
	// - Durum (status) seçimi
	// - Kategori seçimi
	// - Rol seçimi
	// - Öncelik (priority) seçimi
	// - Enum değer seçimi
	//
	// # Örnek Kullanım
	//
	// ```go
	// fields.Select("status").
	//     SetLabel("Durum").
	//     SetOptions([]map[string]interface{}{
	//         {"value": "draft", "label": "Taslak"},
	//         {"value": "published", "label": "Yayında"},
	//         {"value": "archived", "label": "Arşivlendi"},
	//     }).
	//     SetDefault("draft")
	// ```
	//
	// # Özellikler
	//
	// - Tek/çoklu seçim
	// - Arama (searchable) desteği
	// - Gruplandırma (grouping)
	// - Özel option render
	// - Placeholder desteği
	// - Temizle (clear) butonu
	//
	// # Veri Formatı
	//
	// Seçenekler şu formatta tanımlanır:
	// ```go
	// []map[string]interface{}{
	//     {"value": "key1", "label": "Görünen İsim 1"},
	//     {"value": "key2", "label": "Görünen İsim 2"},
	// }
	// ```
	TYPE_SELECT ElementType = "select"

	// TYPE_PANEL, alanları gruplamak için kullanılan konteyner alan tipidir.
	//
	// Bu alan tipi, ilgili alanları görsel olarak gruplamak ve organize etmek için kullanılır.
	// Başlık, açıklama ve daraltılabilir (collapsible) özellikler sunar.
	//
	// # Kullanım Senaryoları
	//
	// - İlgili alanları gruplama
	// - Form bölümleri oluşturma
	// - Görsel organizasyon
	// - Karmaşık formları basitleştirme
	// - Koşullu alan grupları
	//
	// # Örnek Kullanım
	//
	// ```go
	// fields.Panel("Kişisel Bilgiler").
	//     SetFields([]core.Field{
	//         fields.Text("first_name"),
	//         fields.Text("last_name"),
	//         fields.Email("email"),
	//         fields.Tel("phone"),
	//     }).
	//     SetCollapsible(true).
	//     SetCollapsed(false)
	// ```
	//
	// # Özellikler
	//
	// - Başlık ve açıklama
	// - Daraltılabilir (collapsible)
	// - İç içe (nested) panel desteği
	// - Özel stil ve sınıf
	// - Koşullu görünürlük
	//
	// # Avantajlar
	//
	// - Daha iyi kullanıcı deneyimi
	// - Görsel hiyerarşi
	// - Form organizasyonu
	// - Karmaşıklığı azaltma
	TYPE_PANEL ElementType = "panel"

	// TYPE_TABS, alanları tab'lara ayırmak için kullanılan konteyner alan tipidir.
	//
	// Bu alan tipi, ilgili alanları tab'lar halinde organize etmek için kullanılır.
	// Her tab, kendi başlığı ve içeriği ile ayrı bir bölüm oluşturur.
	//
	// # Kullanım Senaryoları
	//
	// - Çoklu dil desteği (Türkçe, İngilizce, vb. tab'ları)
	// - Kategorize edilmiş form alanları (Genel Bilgiler, Adres, İletişim tab'ları)
	// - Karmaşık formları organize etme
	// - İlgili alanları gruplandırma
	// - Uzun formları bölümlere ayırma
	//
	// # Örnek Kullanım
	//
	// ```go
	// fields.Tabs("Ürün Bilgileri").
	//     AddTab("tr", "Türkçe", []core.Element{
	//         fields.Text("Başlık", "title_tr"),
	//         fields.Textarea("Açıklama", "description_tr"),
	//     }).
	//     AddTab("en", "English", []core.Element{
	//         fields.Text("Title", "title_en"),
	//         fields.Textarea("Description", "description_en"),
	//     }).
	//     WithSide("top").
	//     WithVariant("line")
	// ```
	//
	// # Özellikler
	//
	// - Çoklu tab desteği
	// - Tab pozisyonu (top, bottom, left, right)
	// - Tab variant'ı (default, line)
	// - Varsayılan aktif tab
	// - Her tab'da farklı field'lar
	//
	// # Avantajlar
	//
	// - Daha iyi kullanıcı deneyimi
	// - Form organizasyonu
	// - Alan tasarrufu
	// - Görsel hiyerarşi
	TYPE_TABS ElementType = "tabs"

	// TYPE_STACK, birden fazla alanı tek hücrede/bölgede birlikte göstermek için
	// kullanılan konteyner alan tipidir.
	//
	// Özellikle Display callback içinde birden fazla küçük bileşeni (örn. badge)
	// yan yana veya alt alta döndürmek için kullanılır.
	TYPE_STACK ElementType = "stack"

	// TYPE_RELATIONSHIP, genel ilişki alanı için kullanılan alan tipidir.
	//
	// Bu alan tipi, tüm ilişki tiplerini (BelongsTo, HasMany, BelongsToMany, vb.) destekleyen
	// genel bir ilişki alanıdır. İlişki tipine göre otomatik olarak uygun bileşeni render eder.
	//
	// # Kullanım Senaryoları
	//
	// - Dinamik ilişki alanları
	// - Genel ilişki yönetimi
	// - Tip bağımsız ilişki gösterimi
	// - Otomatik ilişki render
	//
	// # Örnek Kullanım
	//
	// ```go
	// fields.Relationship("user").
	//     SetLabel("Kullanıcı").
	//     SetRelationType("BelongsTo").
	//     SetDisplayField("name")
	// ```
	//
	// # Özellikler
	//
	// - Tüm ilişki tiplerini destekler
	// - Otomatik bileşen seçimi
	// - Esnek yapılandırma
	// - Tip bazlı render
	//
	// # İlgili Dokümantasyon
	//
	// Daha fazla bilgi için bakınız:
	// - docs/Relationships.md - Tüm ilişki tipleri
	// - docs/Fields.md - Alan tipleri
	TYPE_RELATIONSHIP ElementType = "relationship"

	// TYPE_BADGE, durum veya etiket göstermek için kullanılan alan tipidir.
	//
	// Bu alan tipi, renkli badge/etiket olarak değer göstermek için kullanılır.
	// Genellikle durum (status), öncelik (priority) gibi alanlar için kullanılır.
	//
	// # Kullanım Senaryoları
	//
	// - Durum gösterimi (aktif, pasif, beklemede)
	// - Öncelik gösterimi (düşük, orta, yüksek)
	// - Kategori etiketleri
	// - Rol gösterimi
	// - Tip gösterimi
	//
	// # Örnek Kullanım
	//
	// ```go
	// fields.Badge("status").
	//     SetLabel("Durum").
	//     SetColorMap(map[string]string{
	//         "draft": "gray",
	//         "published": "green",
	//         "archived": "red",
	//     }).
	//     SetLabelMap(map[string]string{
	//         "draft": "Taslak",
	//         "published": "Yayında",
	//         "archived": "Arşivlendi",
	//     })
	// ```
	//
	// # Özellikler
	//
	// - Renkli badge gösterimi
	// - Değer bazlı renk eşleme
	// - Özel label eşleme
	// - Icon desteği
	// - Özel stil
	//
	// # Renk Seçenekleri
	//
	// - gray, red, orange, yellow, green, blue, indigo, purple, pink
	TYPE_BADGE ElementType = "badge"

	// TYPE_CODE, kod editörü için kullanılan alan tipidir.
	//
	// Bu alan tipi, kod yazmak ve düzenlemek için syntax highlighting destekli
	// bir editör sunar. Çeşitli programlama dillerini destekler.
	//
	// # Kullanım Senaryoları
	//
	// - Kod snippet'leri
	// - JSON/XML yapılandırma
	// - SQL sorguları
	// - HTML/CSS/JavaScript kodu
	// - API yanıtları
	//
	// # Örnek Kullanım
	//
	// ```go
	// fields.Code("config").
	//     SetLabel("Yapılandırma").
	//     SetLanguage("json").
	//     SetTheme("monokai").
	//     SetHeight(300)
	// ```
	//
	// # Özellikler
	//
	// - Syntax highlighting
	// - Çoklu dil desteği (JavaScript, Python, Go, JSON, vb.)
	// - Tema desteği
	// - Satır numaraları
	// - Otomatik girinti
	// - Kod tamamlama
	// - Arama ve değiştirme
	//
	// # Desteklenen Diller
	//
	// - javascript, typescript, python, go, php, ruby, java, c, cpp, csharp
	// - html, css, scss, json, xml, yaml, sql, markdown, bash
	TYPE_CODE ElementType = "code"

	// TYPE_COLOR, renk seçici için kullanılan alan tipidir.
	//
	// Bu alan tipi, renk değerlerini seçmek ve saklamak için kullanılır.
	// HEX, RGB, RGBA formatlarını destekler.
	//
	// # Kullanım Senaryoları
	//
	// - Tema renkleri
	// - Marka renkleri
	// - UI renk ayarları
	// - Kategori renkleri
	// - Etiket renkleri
	//
	// # Örnek Kullanım
	//
	// ```go
	// fields.Color("brand_color").
	//     SetLabel("Marka Rengi").
	//     SetDefault("#3B82F6").
	//     SetFormat("hex").
	//     SetShowAlpha(true)
	// ```
	//
	// # Özellikler
	//
	// - Renk seçici (color picker)
	// - HEX, RGB, RGBA formatları
	// - Alpha (şeffaflık) desteği
	// - Renk paletleri
	// - Son kullanılan renkler
	// - Manuel giriş desteği
	//
	// # Veri Formatları
	//
	// - HEX: #3B82F6
	// - RGB: rgb(59, 130, 246)
	// - RGBA: rgba(59, 130, 246, 0.8)
	TYPE_COLOR ElementType = "color"

	// TYPE_BOOLEAN_GROUP, birden fazla boolean seçeneği için kullanılan alan tipidir.
	//
	// Bu alan tipi, ilgili boolean seçeneklerini gruplamak için kullanılır.
	// Checkbox grubu olarak render edilir.
	//
	// # Kullanım Senaryoları
	//
	// - İzin (permission) seçimleri
	// - Özellik (feature) seçimleri
	// - Bildirim tercihleri
	// - Görünürlük ayarları
	// - Çoklu onay alanları
	//
	// # Örnek Kullanım
	//
	// ```go
	// fields.BooleanGroup("permissions").
	//     SetLabel("İzinler").
	//     SetOptions([]map[string]interface{}{
	//         {"value": "read", "label": "Okuma"},
	//         {"value": "write", "label": "Yazma"},
	//         {"value": "delete", "label": "Silme"},
	//     }).
	//     SetDefault([]string{"read"})
	// ```
	//
	// # Özellikler
	//
	// - Çoklu checkbox
	// - Tümünü seç/temizle
	// - Gruplandırma
	// - Varsayılan seçimler
	// - Açıklama metinleri
	//
	// # Veri Formatı
	//
	// Seçili değerler array olarak saklanır:
	// ```json
	// ["read", "write"]
	// ```
	TYPE_BOOLEAN_GROUP ElementType = "boolean_group"
)

// Element context constants - Alan bağlam sabitleri
//
// Bu sabitler, alanların hangi görünümlerde (view) gösterileceğini veya gizleneceğini kontrol eder.
// Her sabit, belirli bir görünürlük kuralını temsil eder.
//
// # Context Kategorileri
//
// 1. **Temel Context'ler**: Genel bağlamları tanımlar (CONTEXT_FORM, CONTEXT_DETAIL, CONTEXT_LIST)
// 2. **Show Context'leri**: Alanın gösterileceği yerleri belirtir (SHOW_ON_*)
// 3. **Hide Context'leri**: Alanın gizleneceği yerleri belirtir (HIDE_ON_*)
// 4. **Only Context'leri**: Alanın SADECE gösterileceği yerleri belirtir (ONLY_ON_*)
//
// # Kullanım Mantığı
//
// - **SHOW_ON_***: Alanı belirtilen görünümde göster (diğer görünümlerde varsayılan davranış)
// - **HIDE_ON_***: Alanı belirtilen görünümde gizle (diğer görünümlerde göster)
// - **ONLY_ON_***: Alanı SADECE belirtilen görünümde göster (diğerlerinde gizle)
//
// # Önemli Notlar
//
// - Birden fazla context aynı anda kullanılabilir
// - Context'ler birbirleriyle çakışabilir, bu durumda son eklenen geçerli olur
// - ONLY_* context'leri en kısıtlayıcı olanıdır
const (
	// CONTEXT_FORM, alanın form bağlamında olduğunu belirtir.
	//
	// Bu context, alanın oluşturma (create) veya güncelleme (update) formlarında
	// kullanıldığını gösterir. Genel bir form context'idir.
	//
	// # Kullanım Senaryoları
	//
	// - Form alanlarını tanımlarken
	// - Form validasyonu için
	// - Form davranışlarını kontrol ederken
	//
	// # Örnek Kullanım
	//
	// ```go
	// field := fields.Text("title").
	//     SetContext(core.CONTEXT_FORM)
	// ```
	//
	// # Önemli Notlar
	//
	// - Hem create hem update formlarını kapsar
	// - Daha spesifik context'ler için ONLY_ON_CREATE veya ONLY_ON_UPDATE kullanın
	CONTEXT_FORM ElementContext = "form"

	// CONTEXT_DETAIL, alanın detay görünümü bağlamında olduğunu belirtir.
	//
	// Bu context, alanın bir kaydın detay sayfasında gösterildiğini belirtir.
	// Genellikle salt okunur (read-only) görünüm için kullanılır.
	//
	// # Kullanım Senaryoları
	//
	// - Detay sayfası alanlarını tanımlarken
	// - Salt okunur alan gösterimi için
	// - Genişletilmiş bilgi gösterimi için
	//
	// # Örnek Kullanım
	//
	// ```go
	// field := fields.Text("created_at").
	//     SetContext(core.CONTEXT_DETAIL)
	// ```
	CONTEXT_DETAIL ElementContext = "detail"

	// CONTEXT_LIST, alanın liste görünümü bağlamında olduğunu belirtir.
	//
	// Bu context, alanın kayıt listesi/tablo görünümünde gösterildiğini belirtir.
	// Genellikle özet bilgiler için kullanılır.
	//
	// # Kullanım Senaryoları
	//
	// - Liste/tablo sütunlarını tanımlarken
	// - Özet bilgi gösterimi için
	// - Sıralanabilir alanlar için
	// - Filtrelenebilir alanlar için
	//
	// # Örnek Kullanım
	//
	// ```go
	// field := fields.Text("title").
	//     SetContext(core.CONTEXT_LIST).
	//     SetSortable(true)
	// ```
	CONTEXT_LIST ElementContext = "list"

	// SHOW_ON_FORM, alanın form sayfalarında gösterilmesi gerektiğini belirtir.
	//
	// Bu context, alanın hem oluşturma (create) hem de güncelleme (update)
	// formlarında gösterilmesini sağlar.
	//
	// # Kullanım Senaryoları
	//
	// - Düzenlenebilir alanları belirtmek için
	// - Form-specific alanlar için
	// - Kullanıcı girişi gereken alanlar için
	//
	// # Örnek Kullanım
	//
	// ```go
	// field := fields.Text("title").
	//     ShowOn(core.SHOW_ON_FORM)
	// ```
	//
	// # Farklar
	//
	// - SHOW_ON_FORM: Diğer görünümlerde de gösterilebilir
	// - ONLY_ON_FORM: SADECE form sayfalarında gösterilir
	SHOW_ON_FORM ElementContext = "show_on_form"

	// SHOW_ON_DETAIL, alanın detay sayfasında gösterilmesi gerektiğini belirtir.
	//
	// Bu context, alanın kayıt detay sayfasında gösterilmesini sağlar.
	// Genellikle salt okunur bilgiler için kullanılır.
	//
	// # Kullanım Senaryoları
	//
	// - Detay sayfasında gösterilecek alanlar için
	// - Salt okunur bilgiler için
	// - Genişletilmiş bilgi gösterimi için
	// - Timestamp alanları için (created_at, updated_at)
	//
	// # Örnek Kullanım
	//
	// ```go
	// field := fields.DateTime("created_at").
	//     ShowOn(core.SHOW_ON_DETAIL)
	// ```
	//
	// # Farklar
	//
	// - SHOW_ON_DETAIL: Diğer görünümlerde de gösterilebilir
	// - ONLY_ON_DETAIL: SADECE detay sayfasında gösterilir
	SHOW_ON_DETAIL ElementContext = "show_on_detail"

	// SHOW_ON_LIST, alanın liste görünümünde gösterilmesi gerektiğini belirtir.
	//
	// Bu context, alanın kayıt listesi/tablo görünümünde gösterilmesini sağlar.
	// Genellikle özet bilgiler ve sıralanabilir alanlar için kullanılır.
	//
	// # Kullanım Senaryoları
	//
	// - Liste/tablo sütunları için
	// - Özet bilgi gösterimi için
	// - Sıralanabilir alanlar için
	// - Aranabilir alanlar için
	//
	// # Örnek Kullanım
	//
	// ```go
	// field := fields.Text("title").
	//     ShowOn(core.SHOW_ON_LIST).
	//     SetSortable(true).
	//     SetSearchable(true)
	// ```
	//
	// # Farklar
	//
	// - SHOW_ON_LIST: Diğer görünümlerde de gösterilebilir
	// - ONLY_ON_LIST: SADECE liste görünümünde gösterilir
	SHOW_ON_LIST ElementContext = "show_on_list"

	// HIDE_ON_LIST, alanın liste görünümünde gizlenmesi gerektiğini belirtir.
	//
	// Bu context, alanın kayıt listesi/tablo görünümünde gizlenmesini sağlar.
	// Diğer görünümlerde (form, detay) gösterilir.
	//
	// # Kullanım Senaryoları
	//
	// - Uzun metin alanlarını liste görünümünden gizlemek
	// - Detaylı bilgileri liste görünümünden gizlemek
	// - Performans optimizasyonu için
	// - Tablo genişliğini kontrol etmek için
	//
	// # Örnek Kullanım
	//
	// ```go
	// // Uzun açıklama alanını listeden gizle
	// field := fields.Textarea("description").
	//     HideOn(core.HIDE_ON_LIST)
	//
	// // Zengin metin içeriğini listeden gizle
	// field := fields.Richtext("content").
	//     HideOn(core.HIDE_ON_LIST)
	// ```
	//
	// # Avantajlar
	//
	// - Liste görünümü performansını artırır
	// - Tablo genişliğini optimize eder
	// - Kullanıcı deneyimini iyileştirir
	HIDE_ON_LIST ElementContext = "hide_on_list"

	// HIDE_ON_DETAIL, alanın detay sayfasında gizlenmesi gerektiğini belirtir.
	//
	// Bu context, alanın kayıt detay sayfasında gizlenmesini sağlar.
	// Diğer görünümlerde (form, liste) gösterilir.
	//
	// # Kullanım Senaryoları
	//
	// - Sadece form ve liste için gerekli alanları gizlemek
	// - Detay sayfasında gereksiz bilgileri gizlemek
	// - Tekrarlayan bilgileri gizlemek
	//
	// # Örnek Kullanım
	//
	// ```go
	// // Sıralama alanını detay sayfasından gizle
	// field := fields.Number("order").
	//     HideOn(core.HIDE_ON_DETAIL)
	// ```
	HIDE_ON_DETAIL ElementContext = "hide_on_detail"

	// HIDE_ON_CREATE, alanın oluşturma formunda gizlenmesi gerektiğini belirtir.
	//
	// Bu context, alanın yeni kayıt oluşturma formunda gizlenmesini sağlar.
	// Güncelleme formunda ve diğer görünümlerde gösterilir.
	//
	// # Kullanım Senaryoları
	//
	// - ID alanlarını oluşturma formundan gizlemek
	// - Otomatik oluşturulan alanları gizlemek (created_at, updated_at)
	// - Sistem alanlarını gizlemek
	// - İlişkisel alanları gizlemek (henüz kayıt yok)
	//
	// # Örnek Kullanım
	//
	// ```go
	// // ID alanını oluşturma formundan gizle
	// field := fields.ID("id").
	//     HideOn(core.HIDE_ON_CREATE)
	//
	// // Oluşturma zamanını oluşturma formundan gizle
	// field := fields.DateTime("created_at").
	//     HideOn(core.HIDE_ON_CREATE)
	//
	// // HasMany ilişkisini oluşturma formundan gizle
	// field := fields.HasMany("posts").
	//     HideOn(core.HIDE_ON_CREATE)
	// ```
	//
	// # Yaygın Kullanım Örnekleri
	//
	// - ID alanları
	// - Timestamp alanları (created_at, updated_at, deleted_at)
	// - HasMany ilişkileri (henüz kayıt yok)
	// - Otomatik hesaplanan alanlar
	HIDE_ON_CREATE ElementContext = "hide_on_create"

	// HIDE_ON_UPDATE, alanın güncelleme formunda gizlenmesi gerektiğini belirtir.
	//
	// Bu context, alanın kayıt güncelleme formunda gizlenmesini sağlar.
	// Oluşturma formunda ve diğer görünümlerde gösterilir.
	//
	// # Kullanım Senaryoları
	//
	// - Sadece oluşturma sırasında ayarlanabilen alanları gizlemek
	// - Değiştirilmemesi gereken alanları korumak
	// - Güncelleme formunu basitleştirmek
	//
	// # Örnek Kullanım
	//
	// ```go
	// // Slug alanını güncelleme formundan gizle (sadece oluşturmada ayarlanır)
	// field := fields.Text("slug").
	//     HideOn(core.HIDE_ON_UPDATE)
	//
	// // Başlangıç değerini güncelleme formundan gizle
	// field := fields.Number("initial_value").
	//     HideOn(core.HIDE_ON_UPDATE)
	// ```
	//
	// # Yaygın Kullanım Örnekleri
	//
	// - Slug alanları (genellikle sadece oluşturmada ayarlanır)
	// - Başlangıç değerleri
	// - Değişmez (immutable) alanlar
	HIDE_ON_UPDATE ElementContext = "hide_on_update"

	// HIDE_ON_API, alanın external API yanıtlarında gizlenmesi gerektiğini belirtir.
	//
	// Bu context özellikle external API kullanımında hassas/sistem alanlarının
	// dış dünyaya açılmasını engellemek için kullanılır.
	HIDE_ON_API ElementContext = "hide_on_api"

	// ONLY_ON_LIST, alanın SADECE liste görünümünde gösterilmesi gerektiğini belirtir.
	//
	// Bu context, alanın SADECE kayıt listesi/tablo görünümünde gösterilmesini sağlar.
	// Form ve detay sayfalarında gizlenir.
	//
	// # Kullanım Senaryoları
	//
	// - Sadece liste için özel hesaplanan alanlar
	// - Özet bilgiler (örn: toplam, sayı)
	// - Liste-specific göstergeler
	// - Hızlı eylem butonları
	//
	// # Örnek Kullanım
	//
	// ```go
	// // Sadece listede gösterilecek özet alan
	// field := fields.Text("summary").
	//     ShowOn(core.ONLY_ON_LIST)
	//
	// // Sadece listede gösterilecek hesaplanan alan
	// field := fields.Number("total_posts").
	//     ShowOn(core.ONLY_ON_LIST)
	// ```
	//
	// # Önemli Notlar
	//
	// - En kısıtlayıcı context'tir
	// - Diğer tüm görünümlerden gizler
	// - Performans optimizasyonu için kullanılabilir
	ONLY_ON_LIST ElementContext = "only_on_list"

	// ONLY_ON_DETAIL, alanın SADECE detay sayfasında gösterilmesi gerektiğini belirtir.
	//
	// Bu context, alanın SADECE kayıt detay sayfasında gösterilmesini sağlar.
	// Form ve liste görünümlerinde gizlenir.
	//
	// # Kullanım Senaryoları
	//
	// - Sadece detay için özel hesaplanan alanlar
	// - Genişletilmiş bilgiler
	// - Salt okunur sistem bilgileri
	// - Timestamp alanları
	// - İlişkisel veri gösterimleri
	//
	// # Örnek Kullanım
	//
	// ```go
	// // Sadece detay sayfasında gösterilecek timestamp
	// field := fields.DateTime("created_at").
	//     ShowOn(core.ONLY_ON_DETAIL)
	//
	// // Sadece detay sayfasında gösterilecek ilişki
	// field := fields.HasMany("posts").
	//     ShowOn(core.ONLY_ON_DETAIL)
	//
	// // Sadece detay sayfasında gösterilecek hesaplanan alan
	// field := fields.Text("full_address").
	//     ShowOn(core.ONLY_ON_DETAIL)
	// ```
	//
	// # Yaygın Kullanım Örnekleri
	//
	// - Timestamp alanları (created_at, updated_at, deleted_at)
	// - HasMany ilişkileri
	// - Hesaplanan alanlar
	// - Genişletilmiş bilgiler
	ONLY_ON_DETAIL ElementContext = "only_on_detail"

	// ONLY_ON_CREATE, alanın SADECE oluşturma formunda gösterilmesi gerektiğini belirtir.
	//
	// Bu context, alanın SADECE yeni kayıt oluşturma formunda gösterilmesini sağlar.
	// Güncelleme formu, liste ve detay sayfalarında gizlenir.
	//
	// # Kullanım Senaryoları
	//
	// - Sadece oluşturma sırasında ayarlanabilen alanlar
	// - Başlangıç değerleri
	// - Değişmez (immutable) alanlar
	// - Şifre alanları (oluşturma için)
	//
	// # Örnek Kullanım
	//
	// ```go
	// // Sadece oluşturma formunda gösterilecek şifre alanı
	// field := fields.Password("password").
	//     ShowOn(core.ONLY_ON_CREATE)
	//
	// // Sadece oluşturma formunda gösterilecek başlangıç değeri
	// field := fields.Number("initial_balance").
	//     ShowOn(core.ONLY_ON_CREATE)
	// ```
	//
	// # Yaygın Kullanım Örnekleri
	//
	// - Şifre alanları (oluşturma için)
	// - Başlangıç değerleri
	// - Değişmez alanlar
	// - Tek seferlik ayarlar
	ONLY_ON_CREATE ElementContext = "only_on_create"

	// ONLY_ON_UPDATE, alanın SADECE güncelleme formunda gösterilmesi gerektiğini belirtir.
	//
	// Bu context, alanın SADECE kayıt güncelleme formunda gösterilmesini sağlar.
	// Oluşturma formu, liste ve detay sayfalarında gizlenir.
	//
	// # Kullanım Senaryoları
	//
	// - Sadece güncelleme sırasında değiştirilebilen alanlar
	// - Şifre değiştirme alanları
	// - Güncelleme-specific alanlar
	//
	// # Örnek Kullanım
	//
	// ```go
	// // Sadece güncelleme formunda gösterilecek şifre değiştirme alanı
	// field := fields.Password("new_password").
	//     ShowOn(core.ONLY_ON_UPDATE).
	//     SetLabel("Yeni Şifre (Değiştirmek için doldurun)")
	//
	// // Sadece güncelleme formunda gösterilecek durum değiştirme alanı
	// field := fields.Select("status_change_reason").
	//     ShowOn(core.ONLY_ON_UPDATE)
	// ```
	//
	// # Yaygın Kullanım Örnekleri
	//
	// - Şifre değiştirme alanları
	// - Durum değiştirme nedeni alanları
	// - Güncelleme notları
	ONLY_ON_UPDATE ElementContext = "only_on_update"

	// ONLY_ON_FORM, alanın SADECE form sayfalarında gösterilmesi gerektiğini belirtir.
	//
	// Bu context, alanın SADECE oluşturma ve güncelleme formlarında gösterilmesini sağlar.
	// Liste ve detay sayfalarında gizlenir.
	//
	// # Kullanım Senaryoları
	//
	// - Sadece form için gerekli alanlar
	// - Düzenlenebilir alanlar
	// - Kullanıcı girişi gereken alanlar
	// - Form-specific yardımcı alanlar
	//
	// # Örnek Kullanım
	//
	// ```go
	// // Sadece formlarda gösterilecek alan
	// field := fields.Text("internal_notes").
	//     ShowOn(core.ONLY_ON_FORM)
	//
	// // Sadece formlarda gösterilecek yardımcı alan
	// field := fields.Boolean("send_notification").
	//     ShowOn(core.ONLY_ON_FORM).
	//     SetDefault(true)
	// ```
	//
	// # Farklar
	//
	// - ONLY_ON_FORM: Hem create hem update formlarında gösterilir
	// - ONLY_ON_CREATE: Sadece create formunda gösterilir
	// - ONLY_ON_UPDATE: Sadece update formunda gösterilir
	//
	// # Yaygın Kullanım Örnekleri
	//
	// - İç notlar
	// - Yardımcı alanlar
	// - Form-specific seçenekler
	// - Geçici alanlar
	ONLY_ON_FORM ElementContext = "only_on_form"
)

// Visibility context constants - Görünürlük bağlam sabitleri
//
// Bu sabitler, alanların hangi spesifik UI bağlamında görünür olması gerektiğini belirler.
// ElementContext'ten daha spesifik ve programatik bir kontrol sağlar.
//
// # Context Tipleri
//
// 1. **ContextIndex**: Liste/index görünümü (tablo görünümü)
// 2. **ContextDetail**: Detay görünümü (tek kayıt görünümü)
// 3. **ContextCreate**: Oluşturma formu
// 4. **ContextUpdate**: Güncelleme formu
// 5. **ContextPreview**: Önizleme modu
//
// # ElementContext ile Farkları
//
// | Özellik | VisibilityContext | ElementContext |
// |---------|-------------------|----------------|
// | Granülerlik | Daha spesifik | Daha genel |
// | Kullanım | Programatik kontrol | Deklaratif kontrol |
// | Esneklik | Yüksek | Orta |
// | Karmaşıklık | Düşük | Orta |
// | API Kullanımı | SetVisibleOn() | ShowOn(), HideOn() |
//
// # Kullanım Senaryoları
//
// - Programatik görünürlük kontrolü
// - Koşullu alan gösterimi
// - Dinamik form oluşturma
// - API tabanlı alan yönetimi
//
// # Örnek Kullanım
//
// ```go
// // Birden fazla context ile kullanım
// field := fields.Text("title").
//
//	SetVisibleOn(core.ContextIndex, core.ContextDetail)
//
// // Tek context ile kullanım
// field := fields.Text("internal_notes").
//
//	SetVisibleOn(core.ContextCreate)
//
// // Gizleme için
// field := fields.Text("sensitive_data").
//
//	SetHiddenOn(core.ContextPreview)
//
// ```
//
// # Önemli Notlar
//
// - VisibilityContext değerleri sabit olarak tanımlanmıştır
// - Birden fazla context aynı anda kullanılabilir
// - ElementContext ile birlikte kullanılabilir
// - Frontend tarafında bu context'lere göre render yapılır
// - Daha temiz ve tip güvenli kod sağlar
const (
	// ContextIndex, alanın index/liste görünümü bağlamında olduğunu belirtir.
	//
	// Bu context, alanın kayıt listesi/tablo görünümünde gösterildiğini belirtir.
	// Genellikle özet bilgiler, sıralanabilir ve filtrelenebilir alanlar için kullanılır.
	//
	// # Kullanım Senaryoları
	//
	// - Liste/tablo sütunlarını tanımlarken
	// - Özet bilgi gösterimi için
	// - Sıralanabilir alanlar için
	// - Aranabilir alanlar için
	// - Filtrelenebilir alanlar için
	//
	// # Örnek Kullanım
	//
	// ```go
	// // Sadece index görünümünde göster
	// field := fields.Text("title").
	//     SetVisibleOn(core.ContextIndex)
	//
	// // Index ve detail görünümlerinde göster
	// field := fields.Text("title").
	//     SetVisibleOn(core.ContextIndex, core.ContextDetail)
	//
	// // Index görünümünde sıralanabilir alan
	// field := fields.Text("title").
	//     SetVisibleOn(core.ContextIndex).
	//     SetSortable(true).
	//     SetSearchable(true)
	// ```
	//
	// # Özellikler
	//
	// - Tablo formatında gösterim
	// - Sıralama (sorting) desteği
	// - Filtreleme desteği
	// - Arama desteği
	// - Sayfalama (pagination)
	//
	// # Yaygın Kullanım Örnekleri
	//
	// - Başlık (title) alanları
	// - Durum (status) alanları
	// - Tarih alanları (created_at, updated_at)
	// - İlişkisel alanlar (user, category)
	// - Badge alanları
	//
	// # Performans Notları
	//
	// - Sadece gerekli alanları index'te gösterin
	// - Uzun metin alanlarını index'ten gizleyin
	// - İlişkisel alanlar için eager loading kullanın
	ContextIndex VisibilityContext = "index"

	// ContextDetail, alanın detay görünümü bağlamında olduğunu belirtir.
	//
	// Bu context, alanın bir kaydın detay sayfasında gösterildiğini belirtir.
	// Genellikle tüm bilgilerin gösterildiği, salt okunur görünüm için kullanılır.
	//
	// # Kullanım Senaryoları
	//
	// - Detay sayfası alanlarını tanımlarken
	// - Salt okunur alan gösterimi için
	// - Genişletilmiş bilgi gösterimi için
	// - İlişkisel veri gösterimi için
	// - Timestamp alanları için
	//
	// # Örnek Kullanım
	//
	// ```go
	// // Sadece detay sayfasında göster
	// field := fields.DateTime("created_at").
	//     SetVisibleOn(core.ContextDetail)
	//
	// // Detay ve index görünümlerinde göster
	// field := fields.Text("title").
	//     SetVisibleOn(core.ContextIndex, core.ContextDetail)
	//
	// // Detay sayfasında ilişkisel veri göster
	// field := fields.HasMany("posts").
	//     SetVisibleOn(core.ContextDetail)
	// ```
	//
	// # Özellikler
	//
	// - Detaylı alan gösterimi
	// - Salt okunur görünüm
	// - İlişkisel veri gösterimi
	// - Genişletilmiş bilgiler
	// - Tüm alan detayları
	//
	// # Yaygın Kullanım Örnekleri
	//
	// - Tüm metin alanları
	// - Timestamp alanları (created_at, updated_at, deleted_at)
	// - HasMany ilişkileri
	// - BelongsToMany ilişkileri
	// - Zengin metin (richtext) alanları
	// - Dosya alanları
	//
	// # Avantajlar
	//
	// - Kullanıcı tüm bilgileri görebilir
	// - İlişkisel veriler detaylı gösterilebilir
	// - Daha iyi kullanıcı deneyimi
	ContextDetail VisibilityContext = "detail"

	// ContextCreate, alanın oluşturma formu bağlamında olduğunu belirtir.
	//
	// Bu context, alanın yeni kayıt oluşturma formunda gösterildiğini belirtir.
	// Genellikle kullanıcı girişi gereken, düzenlenebilir alanlar için kullanılır.
	//
	// # Kullanım Senaryoları
	//
	// - Oluşturma formunda gösterilecek alanlar için
	// - Kullanıcı girişi gereken alanlar için
	// - Zorunlu (required) alanlar için
	// - Başlangıç değerleri için
	//
	// # Örnek Kullanım
	//
	// ```go
	// // Sadece oluşturma formunda göster
	// field := fields.Password("password").
	//     SetVisibleOn(core.ContextCreate).
	//     SetRequired(true)
	//
	// // Oluşturma ve güncelleme formlarında göster
	// field := fields.Text("title").
	//     SetVisibleOn(core.ContextCreate, core.ContextUpdate).
	//     SetRequired(true)
	//
	// // Oluşturma formunda varsayılan değer ile
	// field := fields.Boolean("is_active").
	//     SetVisibleOn(core.ContextCreate, core.ContextUpdate).
	//     SetDefault(true)
	// ```
	//
	// # Özellikler
	//
	// - Düzenlenebilir alanlar
	// - Validasyon kuralları
	// - Varsayılan değerler
	// - Zorunlu alan kontrolü
	// - Form validasyonu
	//
	// # Yaygın Kullanım Örnekleri
	//
	// - Tüm temel alanlar (text, textarea, number, vb.)
	// - Şifre alanları
	// - İlişki seçim alanları (BelongsTo)
	// - Dosya yükleme alanları
	// - Boolean alanları
	//
	// # Gizlenmesi Gereken Alanlar
	//
	// - ID alanları (otomatik oluşturulur)
	// - Timestamp alanları (otomatik oluşturulur)
	// - HasMany ilişkileri (henüz kayıt yok)
	// - Hesaplanan alanlar
	//
	// # Önemli Notlar
	//
	// - Zorunlu alanlar mutlaka ContextCreate'te görünür olmalıdır
	// - Varsayılan değerler ayarlanabilir
	// - Validasyon kuralları eklenmelidir
	ContextCreate VisibilityContext = "create"

	// ContextUpdate, alanın güncelleme formu bağlamında olduğunu belirtir.
	//
	// Bu context, alanın mevcut kayıt güncelleme formunda gösterildiğini belirtir.
	// Genellikle düzenlenebilir alanlar için kullanılır.
	//
	// # Kullanım Senaryoları
	//
	// - Güncelleme formunda gösterilecek alanlar için
	// - Düzenlenebilir alanlar için
	// - Değiştirilebilir bilgiler için
	// - Durum güncellemeleri için
	//
	// # Örnek Kullanım
	//
	// ```go
	// // Sadece güncelleme formunda göster
	// field := fields.Password("new_password").
	//     SetVisibleOn(core.ContextUpdate).
	//     SetLabel("Yeni Şifre (Değiştirmek için doldurun)")
	//
	// // Oluşturma ve güncelleme formlarında göster
	// field := fields.Text("title").
	//     SetVisibleOn(core.ContextCreate, core.ContextUpdate).
	//     SetRequired(true)
	//
	// // Güncelleme formunda salt okunur göster
	// field := fields.Text("slug").
	//     SetVisibleOn(core.ContextUpdate).
	//     SetReadonly(true)
	// ```
	//
	// # Özellikler
	//
	// - Düzenlenebilir alanlar
	// - Mevcut değerleri gösterme
	// - Validasyon kuralları
	// - Salt okunur (readonly) alanlar
	// - Koşullu alanlar
	//
	// # Yaygın Kullanım Örnekleri
	//
	// - Tüm temel alanlar (text, textarea, number, vb.)
	// - Şifre değiştirme alanları
	// - Durum (status) alanları
	// - İlişki güncelleme alanları
	// - Boolean alanları
	//
	// # Gizlenmesi Gereken Alanlar
	//
	// - Değişmez (immutable) alanlar
	// - Otomatik oluşturulan alanlar
	// - Sadece oluşturma için gerekli alanlar
	//
	// # ContextCreate ile Farkları
	//
	// | Özellik | ContextCreate | ContextUpdate |
	// |---------|---------------|---------------|
	// | Kayıt | Yeni kayıt | Mevcut kayıt |
	// | Değerler | Boş/varsayılan | Mevcut değerler |
	// | ID | Yok | Var |
	// | İlişkiler | Sınırlı | Tam |
	//
	// # Önemli Notlar
	//
	// - Mevcut değerler otomatik olarak doldurulur
	// - Validasyon kuralları uygulanır
	// - Salt okunur alanlar gösterilebilir
	ContextUpdate VisibilityContext = "update"

	// ContextPreview, alanın önizleme modu bağlamında olduğunu belirtir.
	//
	// Bu context, alanın önizleme modunda gösterildiğini belirtir.
	// Genellikle yayınlanmadan önce içeriğin nasıl görüneceğini göstermek için kullanılır.
	//
	// # Kullanım Senaryoları
	//
	// - İçerik önizleme için
	// - Yayın öncesi kontrol için
	// - Kullanıcıya gösterilecek alanlar için
	// - Public görünüm simülasyonu için
	//
	// # Örnek Kullanım
	//
	// ```go
	// // Önizlemede göster
	// field := fields.Text("title").
	//     SetVisibleOn(core.ContextPreview)
	//
	// // Önizlemede gizle (hassas bilgi)
	// field := fields.Text("internal_notes").
	//     SetHiddenOn(core.ContextPreview)
	//
	// // Önizlemede sadece public alanları göster
	// field := fields.Richtext("content").
	//     SetVisibleOn(core.ContextPreview, core.ContextDetail)
	// ```
	//
	// # Özellikler
	//
	// - Salt okunur görünüm
	// - Public görünüm simülasyonu
	// - Hassas bilgileri gizleme
	// - Yayın öncesi kontrol
	// - Kullanıcı deneyimi testi
	//
	// # Yaygın Kullanım Örnekleri
	//
	// - Başlık (title) alanları
	// - İçerik (content) alanları
	// - Resim alanları
	// - Video alanları
	// - Public bilgiler
	//
	// # Gizlenmesi Gereken Alanlar
	//
	// - İç notlar (internal notes)
	// - Sistem alanları
	// - Hassas bilgiler
	// - Admin-only alanlar
	// - Metadata alanları
	//
	// # Avantajlar
	//
	// - Yayınlanmadan önce kontrol
	// - Kullanıcı deneyimi testi
	// - Hata önleme
	// - Güvenlik (hassas bilgileri gizleme)
	//
	// # Önemli Notlar
	//
	// - Önizleme modu genellikle salt okunurdur
	// - Hassas bilgiler gizlenmelidir
	// - Public görünümü simüle eder
	// - SEO kontrolü için kullanılabilir
	ContextPreview VisibilityContext = "preview"
)
