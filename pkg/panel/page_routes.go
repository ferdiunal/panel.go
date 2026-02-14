/// # Panel Sayfa Rotaları Modülü
///
/// Bu modül, panel uygulamasında sayfa yönetimi ve rota işlemlerini gerçekleştirir.
/// Sayfaların listelenmesi, detaylarının alınması ve verilerinin kaydedilmesi gibi
/// temel işlemleri HTTP endpoint'leri aracılığıyla sağlar.
///
/// ## Genel Özellikler
/// - Sayfa listesi API'si (GET /api/pages)
/// - Sayfa detayları API'si (GET /api/pages/:slug)
/// - Sayfa kaydetme API'si (POST /api/pages/:slug)
/// - Erişim kontrolü ve yetkilendirme
/// - Dinamik ayarlar yönetimi

package panel

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/i18n"
	"github.com/gofiber/fiber/v2"
)

// / ## handlePages - Tüm Sayfaları Listele
// /
// / Panelde kayıtlı olan ve kullanıcının erişim yetkisine sahip olduğu tüm sayfaları
// / JSON formatında döner. Gizli sayfalar ve erişim izni olmayan sayfalar filtrelenir.
// /
// / ### HTTP Bilgisi
// / - **Metod**: GET
// / - **Rota**: /api/pages
// / - **Yetkilendirme**: Gerekli (CanAccess kontrolü yapılır)
// /
// / ### Parametreler
// / - `c *context.Context`: İstek bağlamı, kullanıcı bilgilerini ve yetkilendirmeyi içerir
// /
// / ### Dönüş Değeri
// / - `error`: İşlem sırasında oluşan hata (başarılı ise nil)
// / - JSON Yanıt Yapısı:
// /   ```json
// /   {
// /     "data": [
// /       {
// /         "slug": "dashboard",
// /         "title": "Kontrol Paneli",
// /         "description": "Ana kontrol paneli",
// /         "icon": "dashboard",
// /         "group": "Ana",
// /         "order": 1,
// /         "visible": true
// /       }
// /     ]
// /   }
// /   ```
// /
// / ### Kullanım Senaryoları
// / 1. **Navigasyon Menüsü Oluşturma**: Frontend, kullanıcının erişebileceği sayfaları
// /    bu endpoint'ten alarak dinamik menü oluşturur.
// / 2. **Sayfa Keşfi**: Yeni kullanıcılar, panelde hangi sayfaların mevcut olduğunu
// /    bu endpoint'ten öğrenebilir.
// / 3. **Erişim Kontrolü**: Gizli veya erişim izni olmayan sayfalar otomatik olarak
// /    filtrelenir.
// /
// / ### Kullanım Örneği
// / ```javascript
// / // Frontend tarafında
// / fetch('/api/pages')
// /   .then(res => res.json())
// /   .then(data => {
// /     // data.data içinde sayfaların listesi
// /     const pages = data.data;
// /     pages.forEach(page => {
// /       console.log(`${page.title} (${page.slug})`);
// /     });
// /   });
// / ```
// /
// / ### Avantajlar
// / - Dinamik sayfa keşfi sağlar
// / - Erişim kontrolü otomatik olarak uygulanır
// / - Sayfa sıralaması (order) ile özel düzenleme yapılabilir
// / - Sayfa gruplandırması ile organize menü oluşturulabilir
// /
// / ### Önemli Notlar
// / - Yalnızca `Visible()` true olan sayfalar döndürülür
// / - `CanAccess(c)` false olan sayfalar filtrelenir
// / - Boş liste döndürülebilir (hiç erişilebilir sayfa yoksa)
// / - Sayfa sırası `Order` alanına göre belirlenir
// /
// / ### Uyarılar
// / - Performans: Çok sayıda sayfa varsa, filtreleme işlemi zaman alabilir
// / - Caching: Sık çağrılan endpoint için caching mekanizması düşünülebilir
func (p *Panel) handlePages(c *context.Context) error {
	/// ### PageItem Struct'ı
	/// Sayfa bilgilerini JSON formatında döndürmek için kullanılan yapı.
	///
	/// #### Alanlar
	/// - `Slug string`: Sayfanın benzersiz tanımlayıcısı (URL'de kullanılır)
	/// - `Title string`: Sayfanın başlığı (UI'da gösterilir)
	/// - `Description string`: Sayfanın açıklaması
	/// - `Icon string`: Sayfanın ikonu (CSS class veya icon adı)
	/// - `Group string`: Sayfanın ait olduğu grup (menü organizasyonu için)
	/// - `Order int`: Sayfanın sıra numarası (küçük = önce gösterilir)
	/// - `Visible bool`: Sayfanın görünür olup olmadığı
	type PageItem struct {
		Slug        string `json:"slug"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
		Group       string `json:"group"`
		Order       int    `json:"order"`
		Visible     bool   `json:"visible"`
	}

	items := []PageItem{}

	/// Tüm kayıtlı sayfaları döngüyle kontrol et
	for slug, pg := range p.pages {
		/// Sayfanın görünür olup olmadığını ve kullanıcının erişim yetkisini kontrol et
		/// - Gizli sayfalar atlanır
		/// - Erişim izni olmayan sayfalar atlanır
		if !pg.Visible() || !pg.CanAccess(c) {
			continue
		}

		/// Sayfayı PageItem yapısına dönüştür ve listeye ekle
		items = append(items, PageItem{
			Slug:        slug,
			Title:       i18n.Trans(c.Ctx, pg.Title()),
			Description: i18n.Trans(c.Ctx, pg.Description()),
			Icon:        pg.Icon(),
			Group:       pg.Group(),
			Order:       pg.NavigationOrder(),
			Visible:     pg.Visible(),
		})
	}

	/// JSON formatında yanıt döndür
	return c.JSON(fiber.Map{
		"data": items,
	})
}

// / ## handlePageDetail - Sayfa Detaylarını Getir
// /
// / Belirli bir sayfanın tüm detaylarını, kartlarını (widgets) ve alanlarını (fields)
// / JSON formatında döner. Sayfanın tam yapısını frontend'e sağlar.
// /
// / ### HTTP Bilgisi
// / - **Metod**: GET
// / - **Rota**: /api/pages/:slug
// / - **Parametreler**: slug (URL parametresi)
// / - **Yetkilendirme**: Gerekli (CanAccess kontrolü yapılır)
// /
// / ### Parametreler
// / - `c *context.Context`: İstek bağlamı
// /   - `c.Params("slug")`: Sayfanın benzersiz tanımlayıcısı
// /
// / ### Dönüş Değeri
// / - `error`: İşlem sırasında oluşan hata (başarılı ise nil)
// / - JSON Yanıt Yapısı (Başarılı):
// /   ```json
// /   {
// /     "slug": "dashboard",
// /     "title": "Kontrol Paneli",
// /     "description": "Ana kontrol paneli",
// /     "meta": {
// /       "cards": [
// /         {
// /           "id": "card-1",
// /           "title": "İstatistikler",
// /           "data": { "users": 100, "posts": 50 }
// /         }
// /       ],
// /       "fields": [
// /         {
// /           "key": "site_name",
// /           "label": "Site Adı",
// /           "type": "text",
// /           "data": "Benim Sitesi"
// /         }
// /       ]
// /     }
// /   }
// /   ```
// /
// / ### Hata Yanıtları
// / - **404 Not Found**: Sayfa bulunamadı
// /   ```json
// /   { "error": "Page not found" }
// /   ```
// / - **403 Forbidden**: Erişim izni yok
// /   ```json
// /   { "error": "Access denied" }
// /   ```
// /
// / ### Kullanım Senaryoları
// / 1. **Sayfa Yükleme**: Kullanıcı bir sayfaya gittiğinde, bu endpoint'ten
// /    sayfanın tüm içeriği yüklenir.
// / 2. **Dinamik İçerik**: Kartlar ve alanlar dinamik olarak yüklenir ve
// /    veritabanından veri çekilir.
// / 3. **Ayarlar Sayfası**: Özel olarak "settings" sayfası için ayar değerleri
// /    enjekte edilir.
// / 4. **Form Oluşturma**: Frontend, alanlar listesinden dinamik form oluşturur.
// /
// / ### Kullanım Örneği
// / ```javascript
// / // Frontend tarafında
// / const slug = 'dashboard';
// / fetch(`/api/pages/${slug}`)
// /   .then(res => {
// /     if (res.status === 404) throw new Error('Sayfa bulunamadı');
// /     if (res.status === 403) throw new Error('Erişim reddedildi');
// /     return res.json();
// /   })
// /   .then(data => {
// /     console.log(`Sayfa: ${data.title}`);
// /     console.log(`Kartlar: ${data.meta.cards.length}`);
// /     console.log(`Alanlar: ${data.meta.fields.length}`);
// /   })
// /   .catch(err => console.error(err));
// / ```
// /
// / ### Avantajlar
// / - Sayfanın tüm yapısını tek bir istek ile alır
// / - Kartlar ve alanlar otomatik olarak veri ile doldurulur
// / - Ayarlar sayfası için özel işleme yapılır
// / - Erişim kontrolü otomatik olarak uygulanır
// /
// / ### Dezavantajlar
// / - Kart verisi çekilirken hata oluşursa, null değer döndürülür
// / - Çok sayıda kart varsa, yanıt boyutu büyük olabilir
// / - Veri çekme işlemi zaman alabilir (N+1 sorunu riski)
// /
// / ### Önemli Notlar
// / - Kart verisi çekilirken hata oluşursa, `data` alanı `null` olur
// / - "settings" sayfası için özel işleme yapılır
// / - Dinamik ayar değerleri önce kontrol edilir, sonra varsayılan değerler
// / - Tüm kartlar ve alanlar JSON serileştirilir
// /
// / ### Uyarılar
// / - Performans: Çok sayıda kart varsa, veri çekme işlemi yavaş olabilir
// / - Caching: Sık çağrılan endpoint için caching mekanizması düşünülebilir
// / - Hata Yönetimi: Kart verisi çekilirken hata oluşursa, frontend'de sorun yaşanabilir
// / - N+1 Sorunu: Her kart için ayrı sorgu yapılabilir, bu da performansı etkileyebilir
func (p *Panel) handlePageDetail(c *context.Context) error {
	/// Sayfanın slug'ını URL parametrelerinden al
	slug := c.Params("slug")

	/// Slug'a göre sayfayı panelden bul
	pg, ok := p.pages[slug]
	if !ok {
		/// Sayfa bulunamadı, 404 hatası döndür
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Page not found",
		})
	}

	/// Kullanıcının bu sayfaya erişim yetkisini kontrol et
	if !pg.CanAccess(c) {
		/// Erişim izni yok, 403 hatası döndür
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied",
		})
	}

	/// ### Kartları Hazırla
	/// Sayfanın tüm kartlarını (widgets) işle ve veri ile doldur
	cards := []map[string]interface{}{}
	for _, card := range pg.Cards() {
		/// Kartı JSON formatına dönüştür
		serialized := card.JsonSerialize()

		/// Kartın verilerini veritabanından çek
		/// Eğer veri çekme başarılı ise, serialized yapıya ekle
		if data, err := card.Resolve(c, p.Db); err == nil {
			serialized["data"] = data
		} else {
			/// Veri çekme başarısız ise, null değer ata
			/// Bu, frontend'in boş durum göstermesini sağlar
			serialized["data"] = nil
		}
		cards = append(cards, serialized)
	}

	/// ### Alanları Hazırla
	/// Sayfanın tüm alanlarını (fields) işle ve gerekli verileri enjekte et

	var fieldsList []map[string]interface{}

	var pageFields []fields.Element
	if ctxAware, ok := pg.(interface {
		GetFieldsWithContext(ctx *context.Context) []fields.Element
	}); ok {
		pageFields = ctxAware.GetFieldsWithContext(c)
	} else {
		pageFields = pg.Fields()
	}

	// Inject context to all fields for i18n
	for _, f := range pageFields {
		if setter, ok := f.(interface{ SetContextForI18n(*fiber.Ctx) }); ok {
			setter.SetContextForI18n(c.Ctx)
		}
		fieldsList = append(fieldsList, f.JsonSerialize())
	}

	/// Sayfanın tüm detaylarını JSON formatında döndür
	return c.JSON(fiber.Map{
		"slug":        pg.Slug(),
		"title":       i18n.Trans(c.Ctx, pg.Title()),
		"description": i18n.Trans(c.Ctx, pg.Description()),
		"meta": fiber.Map{
			"cards":  cards,
			"fields": fieldsList,
		},
	})
}

// / ## handlePageSave - Sayfa Verilerini Kaydet
// /
// / Belirli bir sayfanın verilerini alır, işler ve veritabanına kaydeder.
// / Ayarlar sayfası için özel olarak ayarları yeniden yükler (hot reload).
// /
// / ### HTTP Bilgisi
// / - **Metod**: POST
// / - **Rota**: /api/pages/:slug
// / - **Content-Type**: application/json
// / - **Yetkilendirme**: Gerekli (CanAccess kontrolü yapılır)
// /
// / ### Parametreler
// / - `c *context.Context`: İstek bağlamı
// /   - `c.Params("slug")`: Sayfanın benzersiz tanımlayıcısı
// /   - `c.Body()`: JSON formatında sayfa verileri
// /
// / ### İstek Gövdesi Örneği
// / ```json
// / {
// /   "site_name": "Yeni Site Adı",
// /   "register": true,
// /   "forgot_password": true,
// /   "custom_field": "custom_value"
// / }
// / ```
// /
// / ### Dönüş Değeri
// / - `error`: İşlem sırasında oluşan hata (başarılı ise nil)
// / - JSON Yanıt Yapısı (Başarılı):
// /   ```json
// /   {
// /     "message": "Settings saved"
// /   }
// /   ```
// /
// / ### Hata Yanıtları
// / - **404 Not Found**: Sayfa bulunamadı
// /   ```json
// /   { "error": "Page not found" }
// /   ```
// / - **403 Forbidden**: Erişim izni yok
// /   ```json
// /   { "error": "Access denied" }
// /   ```
// / - **400 Bad Request**: Geçersiz JSON
// /   ```json
// /   { "error": "Invalid JSON" }
// /   ```
// / - **500 Internal Server Error**: Kaydetme sırasında hata
// /   ```json
// /   { "error": "Hata mesajı" }
// /   ```
// /
// / ### Kullanım Senaryoları
// / 1. **Ayarlar Güncelleme**: Kullanıcı ayarlar sayfasında değişiklik yaptığında,
// /    bu endpoint'e POST isteği gönderilerek veriler kaydedilir.
// / 2. **Dinamik Ayarlar**: Ayarlar sayfası için dinamik değerler kaydedilir ve
// /    anında yüklenir (hot reload).
// / 3. **Form Gönderimi**: Frontend'deki formlar bu endpoint'e veri göndererek
// /    backend'de işlenir.
// / 4. **Veri Doğrulama**: Gönderilen veriler sayfanın Save metodu tarafından
// /    doğrulanır ve işlenir.
// /
// / ### Kullanım Örneği
// / ```javascript
// / // Frontend tarafında
// / const slug = 'settings';
// / const data = {
// /   site_name: 'Yeni Site Adı',
// /   register: true,
// /   forgot_password: false
// / };
// /
// / fetch(`/api/pages/${slug}`, {
// /   method: 'POST',
// /   headers: {
// /     'Content-Type': 'application/json'
// /   },
// /   body: JSON.stringify(data)
// / })
// /   .then(res => {
// /     if (res.status === 404) throw new Error('Sayfa bulunamadı');
// /     if (res.status === 403) throw new Error('Erişim reddedildi');
// /     if (res.status === 400) throw new Error('Geçersiz JSON');
// /     if (res.status === 500) throw new Error('Sunucu hatası');
// /     return res.json();
// /   })
// /   .then(data => {
// /     console.log(data.message); // "Settings saved"
// /     // Sayfayı yenile veya state'i güncelle
// /   })
// /   .catch(err => console.error(err));
// / ```
// /
// / ### Avantajlar
// / - Sayfanın Save metodu tarafından özel işleme yapılabilir
// / - Ayarlar sayfası için otomatik hot reload
// / - Erişim kontrolü otomatik olarak uygulanır
// / - Hata yönetimi kapsamlı
// / - Dinamik ayarlar sistemi ile esnek yapı
// /
// / ### Dezavantajlar
// / - Kart verisi kaydedilmez (sadece alanlar kaydedilir)
// / - Ayarlar sayfası için özel işleme yapılması gerekir
// / - Veri doğrulama sayfanın Save metoduna bağlıdır
// /
// / ### Önemli Notlar
// / - Yalnızca alanlar (fields) kaydedilir, kartlar (cards) kaydedilmez
// / - Ayarlar sayfası için özel olarak LoadSettings() çağrılır
// / - Gönderilen veriler map[string]interface{} olarak işlenir
// / - Sayfanın Save metodu tarafından veri doğrulanır ve işlenir
// / - Hot reload, ayarlar sayfası için otomatik olarak yapılır
// /
// / ### Uyarılar
// / - Performans: Çok sayıda alan varsa, kaydetme işlemi zaman alabilir
// / - Veri Doğrulama: Sayfanın Save metodunda veri doğrulanmalıdır
// / - Hot Reload: Ayarlar sayfası için LoadSettings() çağrılır, bu işlem zaman alabilir
// / - Hata Yönetimi: Save metodu hata döndürürse, kullanıcıya hata mesajı gösterilir
// / - Atomiklik: Birden fazla alan kaydedilirken, hata oluşursa kısmi kaydetme olabilir
// /
// / ### İlişkili Fonksiyonlar
// / - `Page.Save()`: Sayfanın verilerini kaydeden metod
// / - `Panel.LoadSettings()`: Ayarları yeniden yükleyen metod
// / - `Page.CanAccess()`: Erişim kontrolü yapan metod
func (p *Panel) handlePageSave(c *context.Context) error {
	/// Sayfanın slug'ını URL parametrelerinden al
	slug := c.Params("slug")

	/// Slug'a göre sayfayı panelden bul
	pg, ok := p.pages[slug]
	if !ok {
		/// Sayfa bulunamadı, 404 hatası döndür
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Page not found",
		})
	}

	/// Kullanıcının bu sayfaya erişim yetkisini kontrol et
	if !pg.CanAccess(c) {
		/// Erişim izni yok, 403 hatası döndür
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied",
		})
	}

	/// İstek gövdesini parse et (JSON veya Multipart Form)
	var data map[string]interface{}
	contentType := c.Get("Content-Type")

	if len(contentType) >= 19 && contentType[:19] == "multipart/form-data" {
		/// Multipart Form Data işleme
		form, err := c.MultipartForm()
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid form data",
			})
		}

		data = make(map[string]interface{})

		/// Form değerlerini map'e ekle
		for key, values := range form.Value {
			if len(values) > 0 {
				/// Tekil değer olarak ekle (çoğul alanlar için mantık genişletilebilir)
				data[key] = values[0]
			}
		}

		/// Dosyaları map'e ekle
		for key, files := range form.File {
			if len(files) > 0 {
				/// İlk dosyayı ekle (çoklu dosya yükleme için mantık genişletilebilir)
				data[key] = files[0]
			}
		}
	} else {
		/// JSON işleme
		if err := c.BodyParser(&data); err != nil {
			/// JSON parse hatası, 400 hatası döndür
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid JSON",
			})
		}
	}

	/// Sayfanın Save metodunu çağırarak verileri kaydet
	/// Bu metod, sayfanın özel işlemesini ve veri doğrulamasını yapar
	if err := pg.Save(c, p.Db, data); err != nil {
		/// Kaydetme sırasında hata oluştu, 500 hatası döndür
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	/// ### Hot Reload - Ayarlar Sayfası İçin Özel İşleme
	/// Eğer kaydedilen sayfa "settings" ise, ayarları anında yeniden yükle
	/// Bu, ayarların uygulamada hemen etkili olmasını sağlar
	if slug == "settings" {
		/// LoadSettings() çağrısı başarısız olsa bile, kaydetme işlemi başarılı sayılır
		/// Hata görmezden gelinir, çünkü veriler zaten kaydedilmiştir
		_ = p.LoadSettings()
	}

	/// Başarılı kaydetme mesajı döndür
	return c.JSON(fiber.Map{
		"message": "Settings saved",
	})
}

/// ## Modül Özeti
///
/// Bu modül, panel uygulamasının sayfa yönetimi için temel HTTP endpoint'lerini sağlar.
/// Üç ana işlevi vardır:
///
/// 1. **Sayfa Listesi** (handlePages):
///    - Tüm erişilebilir sayfaları listeler
///    - Navigasyon menüsü oluşturmak için kullanılır
///
/// 2. **Sayfa Detayları** (handlePageDetail):
///    - Belirli bir sayfanın tüm detaylarını döner
///    - Kartlar ve alanlar ile birlikte veri sağlar
///    - Ayarlar sayfası için özel işleme yapar
///
/// 3. **Sayfa Kaydetme** (handlePageSave):
///    - Sayfa verilerini kaydeder
///    - Ayarlar sayfası için hot reload yapar
///    - Hata yönetimi ve doğrulama sağlar
///
/// ### Güvenlik Özellikleri
/// - Tüm endpoint'lerde erişim kontrolü yapılır
/// - Gizli sayfalar filtrelenir
/// - Yetkilendirme kontrol edilir
/// - Hata mesajları güvenli şekilde döndürülür
///
/// ### Performans Özellikleri
/// - Kartlar ve alanlar dinamik olarak yüklenir
/// - Veri çekme işlemi hata toleranslıdır
/// - Hot reload sadece ayarlar sayfası için yapılır
///
/// ### Genişletilebilirlik
/// - Sayfaların Save metodu tarafından özel işleme yapılabilir
/// - Yeni sayfalar kolayca eklenebilir
/// - Dinamik ayarlar sistemi ile esnek yapı
///
/// ### Kullanım Akışı
/// ```
/// 1. Frontend, /api/pages endpoint'ine GET isteği gönderir
/// 2. Sayfaların listesi döndürülür
/// 3. Kullanıcı bir sayfaya tıklar
/// 4. Frontend, /api/pages/:slug endpoint'ine GET isteği gönderir
/// 5. Sayfanın detayları döndürülür
/// 6. Kullanıcı formu doldurur ve gönderir
/// 7. Frontend, /api/pages/:slug endpoint'ine POST isteği gönderir
/// 8. Veriler kaydedilir ve başarı mesajı döndürülür
/// ```
///
/// ### Önemli Notlar
/// - Tüm endpoint'ler context.Context kullanır
/// - Fiber framework'ü kullanılır
/// - JSON formatında veri alışverişi yapılır
/// - Hata yönetimi kapsamlı
/// - Erişim kontrolü her endpoint'te yapılır
