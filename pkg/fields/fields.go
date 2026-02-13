package fields

import "strings"

// Bu fonksiyon, yeni bir alan şeması oluşturur ve temel konfigürasyonunu başlatır.
//
// Kullanım Senaryosu:
// - Tüm alan türlerinin temelini oluşturur
// - Alan adını ve veritabanı anahtarını ayarlar
// - Varsayılan özellikleri başlatır
//
// Parametreler:
//   - name: Alanın görüntü adı (örn: "Kullanıcı Adı", "E-posta")
//   - attribute: İsteğe bağlı veritabanı sütun adı (varsayılan: name'i küçük harfe çevirip boşlukları alt çizgi ile değiştirir)
//
// Dönüş Değeri:
//   - *Schema: Yapılandırılmış alan şeması pointer'ı
//
// Örnek Kullanım:
//
//	field := NewField("Kullanıcı Adı")                    // Key: "kullanıcı_adı"
//	field := NewField("Kullanıcı Adı", "username")        // Key: "username"
//
// Önemli Notlar:
//   - Döndürülen Schema pointer'ı method chaining için kullanılabilir
//   - Props haritası boş olarak başlatılır, sonradan özellikler eklenebilir
//   - TextAlign varsayılan olarak "left" (sol hizalı) olarak ayarlanır
func NewField(name string, attribute ...string) *Schema {
	attr := strings.ToLower(strings.ReplaceAll(name, " ", "_"))
	if len(attribute) > 0 {
		attr = attribute[0]
	}

	return &Schema{
		Name:      name,
		Key:       attr,
		Props:     make(map[string]interface{}),
		TextAlign: "left",
		LabelText: name,
	}
}

// ============================================================================
// TEMEL ALAN TÜRLERİ (Basic Field Types)
// ============================================================================

// Bu fonksiyon, benzersiz kimlik alanını oluşturur ve yapılandırır.
//
// Kullanım Senaryosu:
// - Veritabanı kayıtlarının birincil anahtarını temsil eder
// - Genellikle otomatik olarak oluşturulan ve değiştirilemeyen bir alandır
// - Yalnızca liste görünümünde gösterilir, form düzenlemesinde gizlenir
//
// Parametreler:
//   - name: İsteğe bağlı alan adı (varsayılan: "ID")
//
// Dönüş Değeri:
//   - *Schema: Yapılandırılmış ID alan şeması pointer'ı
//
// Örnek Kullanım:
//
//	idField := ID()                    // Varsayılan "ID" adı ile
//	idField := ID("Kayıt Numarası")    // Özel ad ile
//
// Önemli Notlar:
//   - Veritabanı anahtarı her zaman "id" olarak ayarlanır
//   - OnlyOnList() metodu çağrılarak sadece liste görünümünde gösterilir
//   - Frontend bileşeni "id-field" olarak ayarlanır
func ID(name ...string) *Schema {
	label := "ID"
	if len(name) > 0 {
		label = name[0]
	}
	f := NewField(label, "id")
	f.View = "id-field" // Frontend bileşen adı
	f.OnlyOnList()
	return f
}

// Bu fonksiyon, standart metin giriş alanı (input type="text") oluşturur ve yapılandırır.
//
// Kullanım Senaryosu:
// - Kullanıcı adı, başlık, kısa açıklama gibi tek satırlı metin girişi için
// - Basit string veri türü için en yaygın kullanılan alan
// - Maksimum uzunluk kısıtlaması eklenebilir
//
// Parametreler:
//   - name: Alanın görüntü adı (örn: "Başlık", "Kullanıcı Adı")
//   - attribute: İsteğe bağlı veritabanı sütun adı
//
// Dönüş Değeri:
//   - *Schema: Yapılandırılmış metin alan şeması pointer'ı
//
// Örnek Kullanım:
//
//	titleField := Text("Başlık")
//	usernameField := Text("Kullanıcı Adı", "username")
//	emailField := Text("E-posta", "email_address")
//
// Önemli Notlar:
//   - Frontend bileşeni "text-field" olarak ayarlanır
//   - TYPE_TEXT sabiti ile alan türü belirtilir
//   - Method chaining ile ek özellikler eklenebilir (örn: .Required(), .MaxLength(100))
func Text(name string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "text-field"
	f.Type = TYPE_TEXT
	return f
}

// Bu fonksiyon, çok satırlı metin giriş alanı (textarea) oluşturur ve yapılandırır.
//
// Kullanım Senaryosu:
// - Uzun açıklamalar, notlar, biyografi gibi çok satırlı metin girişi için
// - Sabit yükseklikte veya dinamik olarak genişleyebilen textarea
// - Markdown veya HTML desteği eklenebilir
//
// Parametreler:
//   - name: Alanın görüntü adı (örn: "Açıklama", "Notlar")
//   - attribute: İsteğe bağlı veritabanı sütun adı
//
// Dönüş Değeri:
//   - *Schema: Yapılandırılmış textarea alan şeması pointer'ı
//
// Örnek Kullanım:
//
//	descField := Textarea("Açıklama")
//	bioField := Textarea("Biyografi", "biography")
//
// Önemli Notlar:
//   - Frontend bileşeni "textarea-field" olarak ayarlanır
//   - TYPE_TEXTAREA sabiti ile alan türü belirtilir
//   - Satır sayısı (rows) özelliği Props haritasına eklenebilir
func Textarea(name string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "textarea-field"
	f.Type = TYPE_TEXTAREA
	return f
}

// Bu fonksiyon, zengin metin editörü (WYSIWYG - What You See Is What You Get) alanı oluşturur.
//
// Kullanım Senaryosu:
// - Blog yazıları, ürün açıklamaları, HTML içeriği gibi zengin metin girişi için
// - Metin biçimlendirmesi (bold, italic, underline), listeler, linkler vb. destekler
// - Genellikle TinyMCE, Quill veya benzer editörler ile entegre edilir
//
// Parametreler:
//   - name: Alanın görüntü adı (örn: "İçerik", "Açıklama")
//   - attribute: İsteğe bağlı veritabanı sütun adı
//
// Dönüş Değeri:
//   - *Schema: Yapılandırılmış zengin metin alan şeması pointer'ı
//
// Örnek Kullanım:
//
//	contentField := RichText("İçerik")
//	descField := RichText("Ürün Açıklaması", "product_description")
//
// Önemli Notlar:
//   - Frontend bileşeni "richtext-field" olarak ayarlanır
//   - TYPE_RICHTEXT sabiti ile alan türü belirtilir
//   - HTML sanitization gerekli olabilir (XSS koruması için)
func RichText(name string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "richtext-field"
	f.Type = TYPE_RICHTEXT
	return f
}

// Bu fonksiyon, şifre giriş alanı (input type="password") oluşturur ve yapılandırır.
//
// Kullanım Senaryosu:
// - Kullanıcı şifresi, API anahtarı, gizli token gibi hassas bilgilerin girişi için
// - Girilen karakterler maskelenir ve görüntülenmez
// - Genellikle şifre güç göstergesi ve göster/gizle butonu ile birlikte kullanılır
//
// Parametreler:
//   - name: Alanın görüntü adı (örn: "Şifre", "Yeni Şifre")
//   - attribute: İsteğe bağlı veritabanı sütun adı
//
// Dönüş Değeri:
//   - *Schema: Yapılandırılmış şifre alan şeması pointer'ı
//
// Örnek Kullanım:
//
//	passField := Password("Şifre")
//	newPassField := Password("Yeni Şifre", "new_password")
//
// Önemli Notlar:
//   - Frontend bileşeni "password-field" olarak ayarlanır
//   - TYPE_PASSWORD sabiti ile alan türü belirtilir
//   - Backend'de şifre hash'lenmeli (bcrypt, argon2 vb.)
//   - HTTPS üzerinden iletilmeli
func Password(name string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "password-field"
	f.Type = TYPE_PASSWORD
	return f
}

// Bu fonksiyon, sayı giriş alanı (input type="number") oluşturur ve yapılandırır.
//
// Kullanım Senaryosu:
// - Yaş, fiyat, miktar, telefon numarası gibi sayısal değerlerin girişi için
// - Min/max değer kısıtlaması ve adım (step) ayarı yapılabilir
// - Otomatik olarak sayı dışı karakterleri filtreler
//
// Parametreler:
//   - name: Alanın görüntü adı (örn: "Fiyat", "Miktar")
//   - attribute: İsteğe bağlı veritabanı sütun adı
//
// Dönüş Değeri:
//   - *Schema: Yapılandırılmış sayı alan şeması pointer'ı
//
// Örnek Kullanım:
//
//	priceField := Number("Fiyat")
//	quantityField := Number("Miktar", "quantity")
//	ageField := Number("Yaş", "age")
//
// Önemli Notlar:
//   - Frontend bileşeni "number-field" olarak ayarlanır
//   - TYPE_NUMBER sabiti ile alan türü belirtilir
//   - Min/Max değerleri Props haritasına eklenebilir
//   - Step değeri (örn: 0.01 para birimi için) ayarlanabilir
func Number(name string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "number-field"
	f.Type = TYPE_NUMBER
	return f
}

// Bu fonksiyon, e-posta giriş alanı (input type="email") oluşturur ve yapılandırır.
//
// Kullanım Senaryosu:
// - Kullanıcı e-posta adresi, iletişim e-postası gibi e-posta girişi için
// - Otomatik olarak e-posta formatı doğrulaması yapılır
// - Tarayıcı tarafından native e-posta doğrulaması desteklenir
//
// Parametreler:
//   - name: Alanın görüntü adı (örn: "E-posta", "İletişim E-postası")
//   - attribute: İsteğe bağlı veritabanı sütun adı
//
// Dönüş Değeri:
//   - *Schema: Yapılandırılmış e-posta alan şeması pointer'ı
//
// Örnek Kullanım:
//
//	emailField := Email("E-posta")
//	contactField := Email("İletişim E-postası", "contact_email")
//
// Önemli Notlar:
//   - Frontend bileşeni "email-field" olarak ayarlanır
//   - TYPE_EMAIL sabiti ile alan türü belirtilir
//   - Backend'de regex veya email validator kütüphanesi ile doğrulama yapılmalı
//   - Benzersizlik kısıtlaması (unique constraint) eklenebilir
func Email(name string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "email-field"
	f.Type = TYPE_EMAIL
	return f
}

// Bu fonksiyon, görsel yükleme alanı oluşturur ve yapılandırır.
//
// Kullanım Senaryosu:
// - Ürün resimleri, profil fotoğrafları, banner görselleri gibi görsel dosyaları yüklemek için
// - Genellikle JPG, PNG, WebP gibi görsel formatlarını destekler
// - Ön izleme (preview) ve crop işlemleri yapılabilir
//
// Parametreler:
//   - name: Alanın görüntü adı (örn: "Ürün Resmi", "Profil Fotoğrafı")
//   - attribute: İsteğe bağlı veritabanı sütun adı
//
// Dönüş Değeri:
//   - *Schema: Yapılandırılmış görsel alan şeması pointer'ı
//
// Örnek Kullanım:
//
//	imageField := Image("Ürün Resmi")
//	avatarField := Image("Profil Fotoğrafı", "avatar")
//
// Önemli Notlar:
//   - Frontend bileşeni "image-field" olarak ayarlanır
//   - TYPE_FILE sabiti ile alan türü belirtilir
//   - Dosya boyutu ve format kısıtlaması eklenebilir
//   - CDN veya cloud storage entegrasyonu yapılabilir
func Image(name string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "image-field"
	f.Type = TYPE_FILE
	return f
}

// Bu fonksiyon, telefon numarası giriş alanı oluşturur ve yapılandırır.
//
// Kullanım Senaryosu:
// - Müşteri telefon numarası, iletişim numarası gibi telefon girişi için
// - Genellikle metin input olarak render edilir, ancak telefon formatı doğrulaması yapılabilir
// - Ülkeye göre farklı telefon formatları desteklenebilir
//
// Parametreler:
//   - name: Alanın görüntü adı (örn: "Telefon", "Cep Telefonu")
//   - attribute: İsteğe bağlı veritabanı sütun adı
//
// Dönüş Değeri:
//   - *Schema: Yapılandırılmış telefon alan şeması pointer'ı
//
// Örnek Kullanım:
//
//	phoneField := Tel("Telefon")
//	mobileField := Tel("Cep Telefonu", "mobile_phone")
//
// Önemli Notlar:
//   - Frontend bileşeni "text-field" olarak ayarlanır (genellikle text input kullanılır)
//   - TYPE_TEL sabiti ile alan türü belirtilir
//   - Backend'de regex veya libphonenumber gibi kütüphaneler ile doğrulama yapılmalı
//   - Uluslararası format desteği eklenebilir
func Tel(name string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "tel-field" // Genellikle text input kullanılır
	f.Type = TYPE_TEL
	return f
}

// Bu fonksiyon, video yükleme alanı oluşturur ve yapılandırır.
//
// Kullanım Senaryosu:
// - Ürün tanıtım videoları, eğitim videoları, demo videoları gibi video dosyaları yüklemek için
// - Genellikle MP4, WebM, OGG gibi video formatlarını destekler
// - Video ön izlemesi ve oynatıcı entegrasyonu yapılabilir
//
// Parametreler:
//   - name: Alanın görüntü adı (örn: "Tanıtım Videosu", "Demo")
//   - attribute: İsteğe bağlı veritabanı sütun adı
//
// Dönüş Değeri:
//   - *Schema: Yapılandırılmış video alan şeması pointer'ı
//
// Örnek Kullanım:
//
//	videoField := Video("Tanıtım Videosu")
//	demoField := Video("Demo Videosu", "demo_video")
//
// Önemli Notlar:
//   - Frontend bileşeni "file-field" olarak ayarlanır
//   - TYPE_VIDEO sabiti ile alan türü belirtilir
//   - Dosya boyutu kısıtlaması eklenebilir (videolar genellikle büyük dosyalardır)
//   - Video codec ve format doğrulaması yapılabilir
func Video(name string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "file-field"
	f.Type = TYPE_VIDEO
	return f
}

// Bu fonksiyon, ses dosyası yükleme alanı oluşturur ve yapılandırır.
//
// Kullanım Senaryosu:
// - Podcast bölümleri, müzik dosyaları, ses notları gibi ses dosyaları yüklemek için
// - Genellikle MP3, WAV, OGG, M4A gibi ses formatlarını destekler
// - Ses oynatıcı entegrasyonu yapılabilir
//
// Parametreler:
//   - name: Alanın görüntü adı (örn: "Ses Dosyası", "Podcast")
//   - attribute: İsteğe bağlı veritabanı sütun adı
//
// Dönüş Değeri:
//   - *Schema: Yapılandırılmış ses alan şeması pointer'ı
//
// Örnek Kullanım:
//
//	audioField := Audio("Ses Dosyası")
//	podcastField := Audio("Podcast Bölümü", "podcast_episode")
//
// Önemli Notlar:
//   - Frontend bileşeni "file-field" olarak ayarlanır
//   - TYPE_AUDIO sabiti ile alan türü belirtilir
//   - Ses codec ve format doğrulaması yapılabilir
//   - Ses süresi metadata'sı çıkarılabilir
func Audio(name string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "file-field"
	f.Type = TYPE_AUDIO
	return f
}

// Bu fonksiyon, tarih seçim alanı (datepicker) oluşturur ve yapılandırır.
//
// Kullanım Senaryosu:
// - Doğum tarihi, etkinlik tarihi, başlangıç tarihi gibi tarih girişi için
// - Takvim arayüzü ile tarih seçimi yapılabilir
// - Tarih formatı ve aralığı kısıtlaması eklenebilir
//
// Parametreler:
//   - name: Alanın görüntü adı (örn: "Doğum Tarihi", "Etkinlik Tarihi")
//   - attribute: İsteğe bağlı veritabanı sütun adı
//
// Dönüş Değeri:
//   - *Schema: Yapılandırılmış tarih alan şeması pointer'ı
//
// Örnek Kullanım:
//
//	dateField := Date("Doğum Tarihi")
//	eventField := Date("Etkinlik Tarihi", "event_date")
//
// Önemli Notlar:
//   - Frontend bileşeni "date-field" olarak ayarlanır
//   - TYPE_DATE sabiti ile alan türü belirtilir
//   - Tarih formatı (YYYY-MM-DD) standart olarak kullanılır
//   - Min/Max tarih kısıtlaması eklenebilir
func Date(name string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "date-field"
	f.Type = TYPE_DATE
	return f
}

// Bu fonksiyon, tarih ve saat seçim alanı oluşturur ve yapılandırır.
//
// Kullanım Senaryosu:
// - Etkinlik başlangıç zamanı, randevu saati, yayın zamanı gibi tarih ve saat girişi için
// - Takvim ve saat seçici arayüzü ile tarih ve saat seçimi yapılabilir
// - Zaman dilimi (timezone) desteği eklenebilir
//
// Parametreler:
//   - name: Alanın görüntü adı (örn: "Randevu Saati", "Yayın Zamanı")
//   - attribute: İsteğe bağlı veritabanı sütun adı
//
// Dönüş Değeri:
//   - *Schema: Yapılandırılmış tarih-saat alan şeması pointer'ı
//
// Örnek Kullanım:
//
//	datetimeField := DateTime("Randevu Saati")
//	eventField := DateTime("Etkinlik Zamanı", "event_datetime")
//
// Önemli Notlar:
//   - Frontend bileşeni "datetime-field" olarak ayarlanır
//   - TYPE_DATETIME sabiti ile alan türü belirtilir
//   - ISO 8601 formatı (YYYY-MM-DDTHH:mm:ss) standart olarak kullanılır
//   - Zaman dilimi bilgisi saklanabilir
func DateTime(name string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "datetime-field"
	f.Type = TYPE_DATETIME
	return f
}

// Bu fonksiyon, genel dosya yükleme alanı oluşturur ve yapılandırır.
//
// Kullanım Senaryosu:
// - PDF, Word, Excel, ZIP gibi çeşitli dosya türlerini yüklemek için
// - Belge yönetimi, raporlar, ek dosyalar gibi genel dosya depolaması için
// - Dosya türü ve boyutu kısıtlaması eklenebilir
//
// Parametreler:
//   - name: Alanın görüntü adı (örn: "Dosya", "Ek Dosya")
//   - attribute: İsteğe bağlı veritabanı sütun adı
//
// Dönüş Değeri:
//   - *Schema: Yapılandırılmış dosya alan şeması pointer'ı
//
// Örnek Kullanım:
//
//	fileField := File("Dosya")
//	attachmentField := File("Ek Dosya", "attachment")
//
// Önemli Notlar:
//   - Frontend bileşeni "file-field" olarak ayarlanır
//   - TYPE_FILE sabiti ile alan türü belirtilir
//   - Dosya türü whitelist'i eklenebilir (örn: .pdf, .doc, .xls)
//   - Maksimum dosya boyutu kısıtlaması yapılabilir
func File(name string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "file-field"
	f.Type = TYPE_FILE
	return f
}

// Bu fonksiyon, anahtar-değer ikilisi girişi sağlayan alan oluşturur ve yapılandırır.
//
// Kullanım Senaryosu:
// - Dinamik ayarlar, meta veriler, özel özellikler gibi anahtar-değer çiftlerini depolamak için
// - JSON formatında saklanabilen esnek veri yapısı
// - Örn: {"color": "red", "size": "large", "material": "cotton"}
//
// Parametreler:
//   - name: Alanın görüntü adı (örn: "Özellikler", "Meta Veriler")
//   - attribute: İsteğe bağlı veritabanı sütun adı
//
// Dönüş Değeri:
//   - *Schema: Yapılandırılmış anahtar-değer alan şeması pointer'ı
//
// Örnek Kullanım:
//
//	kvField := KeyValue("Özellikler")
//	metaField := KeyValue("Meta Veriler", "metadata")
//
// Önemli Notlar:
//   - Frontend bileşeni "key-value-field" olarak ayarlanır
//   - TYPE_KEY_VALUE sabiti ile alan türü belirtilir
//   - Genellikle JSON formatında veritabanında saklanır
//   - Dinamik form alanları için idealdir
func KeyValue(name string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "key-value-field"
	f.Type = TYPE_KEY_VALUE
	return f
}

// ============================================================================
// İLİŞKİ ALANLARI (Relationship Fields)
// ============================================================================

// Bu fonksiyon, başka bir kaynağa bağlantı (BelongsTo ilişkisi) oluşturur ve yapılandırır.
//
// Kullanım Senaryosu:
// - Bir kaynağın başka bir kaynağa ait olduğu ilişkiyi temsil eder
// - Örn: Bir ürün bir kategoriye ait, bir yorum bir yazıya ait
// - Seçim alanı ile ilişkili kaynağı seçmek için kullanılır
//
// Parametreler:
//   - name: Alanın görüntü adı (örn: "Kategori", "Yazar")
//   - resource: İlişkili kaynağın adı (örn: "categories", "authors")
//   - attribute: İsteğe bağlı veritabanı sütun adı
//
// Dönüş Değeri:
//   - *Schema: Yapılandırılmış bağlantı alan şeması pointer'ı
//
// Örnek Kullanım:
//
//	categoryLink := Link("Kategori", "categories")
//	authorLink := Link("Yazar", "authors", "author_id")
//
// Önemli Notlar:
//   - Frontend bileşeni "link-field" olarak ayarlanır
//   - TYPE_LINK sabiti ile alan türü belirtilir
//   - Props haritasına "resource" anahtarı ile ilişkili kaynak adı eklenir
//   - Method chaining ile ek özellikler eklenebilir
func Link(name string, resource string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "link-field"
	f.Type = TYPE_LINK
	f.Props["resource"] = resource
	return f
}

// Bu fonksiyon, bir kaynağın detayını gösteren (HasOne ilişkisi) alan oluşturur ve yapılandırır.
//
// Kullanım Senaryosu:
// - Bir kaynağın başka bir kaynağa sahip olduğu bire-bir ilişkiyi temsil eder
// - Örn: Bir kullanıcının bir profili, bir ürünün bir açıklaması
// - Genellikle listede gizlenir, detay sayfasında gösterilir
//
// Parametreler:
//   - name: Alanın görüntü adı (örn: "Profil", "Açıklama")
//   - resource: İlişkili kaynağın adı (örn: "profiles", "descriptions")
//   - attribute: İsteğe bağlı veritabanı sütun adı
//
// Dönüş Değeri:
//   - *Schema: Yapılandırılmış detay alan şeması pointer'ı
//
// Örnek Kullanım:
//
//	profileDetail := Detail("Profil", "profiles")
//	descDetail := Detail("Açıklama", "descriptions", "description_id")
//
// Önemli Notlar:
//   - Frontend bileşeni "detail-field" olarak ayarlanır
//   - TYPE_DETAIL sabiti ile alan türü belirtilir
//   - Context HIDE_ON_LIST ile liste görünümünde gizlenir
//   - Detay sayfasında ilişkili kaynağın bilgileri gösterilir
func Detail(name string, resource string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "detail-field"
	f.Type = TYPE_DETAIL
	f.Props["resource"] = resource
	f.Context = HIDE_ON_LIST // Generally hidden on list
	return f
}

// Bu fonksiyon, ilişkili kayıtların listesini gösteren (HasMany ilişkisi) alan oluşturur ve yapılandırır.
//
// Kullanım Senaryosu:
// - Bir kaynağın birden fazla ilişkili kaynağa sahip olduğu bire-çok ilişkiyi temsil eder
// - Örn: Bir yazının birden fazla yorumu, bir kategorinin birden fazla ürünü
// - Genellikle listede gizlenir, detay sayfasında ilişkili kayıtların tablosu gösterilir
//
// Parametreler:
//   - name: Alanın görüntü adı (örn: "Yorumlar", "Ürünler")
//   - resource: İlişkili kaynağın adı (örn: "comments", "products")
//   - attribute: İsteğe bağlı veritabanı sütun adı
//
// Dönüş Değeri:
//   - *Schema: Yapılandırılmış koleksiyon alan şeması pointer'ı
//
// Örnek Kullanım:
//
//	commentsCollection := Collection("Yorumlar", "comments")
//	productsCollection := Collection("Ürünler", "products", "product_ids")
//
// Önemli Notlar:
//   - Frontend bileşeni "collection-field" olarak ayarlanır
//   - TYPE_COLLECTION sabiti ile alan türü belirtilir
//   - Context HIDE_ON_LIST ile liste görünümünde gizlenir
//   - Detay sayfasında ilişkili kayıtların listesi gösterilir
func Collection(name string, resource string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "collection-field"
	f.Type = TYPE_COLLECTION
	f.Props["resource"] = resource
	f.Context = HIDE_ON_LIST
	return f
}

// Bu fonksiyon, çoktan çoka (BelongsToMany ilişkisi) ilişki kurmak için kullanılır ve yapılandırır.
//
// Kullanım Senaryosu:
// - İki kaynak arasında çoktan çoka ilişkiyi temsil eder
// - Örn: Bir ürünün birden fazla etiketi, bir yazının birden fazla kategorisi
// - Genellikle ara tablo (pivot table) ile yönetilir
// - Listede gizlenir, detay sayfasında ilişkili kayıtları seçmek için kullanılır
//
// Parametreler:
//   - name: Alanın görüntü adı (örn: "Etiketler", "Kategoriler")
//   - resource: İlişkili kaynağın adı (örn: "tags", "categories")
//   - attribute: İsteğe bağlı veritabanı sütun adı
//
// Dönüş Değeri:
//   - *Schema: Yapılandırılmış bağlantı alan şeması pointer'ı
//
// Örnek Kullanım:
//
//	tagsConnect := Connect("Etiketler", "tags")
//	categoriesConnect := Connect("Kategoriler", "categories", "category_ids")
//
// Önemli Notlar:
//   - Frontend bileşeni "connect-field" olarak ayarlanır
//   - TYPE_CONNECT sabiti ile alan türü belirtilir
//   - Context HIDE_ON_LIST ile liste görünümünde gizlenir
//   - Detay sayfasında çoklu seçim yapılabilir
func Connect(name string, resource string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "connect-field"
	f.Type = TYPE_CONNECT
	f.Props["resource"] = resource
	f.Context = HIDE_ON_LIST
	return f
}

// Bu fonksiyon, polimorfik ilişki bağlantısı (MorphTo ilişkisi) oluşturur ve yapılandırır.
//
// Kullanım Senaryosu:
// - Bir kaynağın birden fazla farklı kaynak türüne ait olabileceği polimorfik ilişkiyi temsil eder
// - Örn: Bir yorum bir yazıya veya bir ürüne ait olabilir
// - Genellikle "commentable_type" ve "commentable_id" sütunları ile yönetilir
//
// Parametreler:
//   - name: Alanın görüntü adı (örn: "Bağlantılı Kaynak")
//   - attribute: İsteğe bağlı veritabanı sütun adı
//
// Dönüş Değeri:
//   - *Schema: Yapılandırılmış polimorfik bağlantı alan şeması pointer'ı
//
// Örnek Kullanım:
//
//	polyLink := PolyLink("Bağlantılı Kaynak")
//	commentableLink := PolyLink("Yorum Yapılan", "commentable")
//
// Önemli Notlar:
//   - Frontend bileşeni "poly-link-field" olarak ayarlanır
//   - TYPE_POLY_LINK sabiti ile alan türü belirtilir
//   - Polimorfik ilişkiler karmaşık olabilir, dikkatli tasarım gerekir
func PolyLink(name string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "poly-link-field"
	f.Type = TYPE_POLY_LINK
	return f
}

// Bu fonksiyon, polimorfik detay (MorphOne ilişkisi) oluşturur ve yapılandırır.
//
// Kullanım Senaryosu:
// - Bir kaynağın birden fazla farklı kaynak türüne ait olabileceği polimorfik bire-bir ilişkiyi temsil eder
// - Örn: Bir ürünün bir açıklaması olabilir, bir yazının bir özeti olabilir
// - Genellikle listede gizlenir, detay sayfasında gösterilir
//
// Parametreler:
//   - name: Alanın görüntü adı (örn: "Açıklama", "Özet")
//   - resource: İlişkili kaynağın adı (örn: "descriptions", "summaries")
//   - attribute: İsteğe bağlı veritabanı sütun adı
//
// Dönüş Değeri:
//   - *Schema: Yapılandırılmış polimorfik detay alan şeması pointer'ı
//
// Örnek Kullanım:
//
//	polyDetail := PolyDetail("Açıklama", "descriptions")
//	summaryDetail := PolyDetail("Özet", "summaries", "summary_id")
//
// Önemli Notlar:
//   - Frontend bileşeni "poly-detail-field" olarak ayarlanır
//   - TYPE_POLY_DETAIL sabiti ile alan türü belirtilir
//   - Context HIDE_ON_LIST ile liste görünümünde gizlenir
func PolyDetail(name string, resource string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "poly-detail-field"
	f.Type = TYPE_POLY_DETAIL
	f.Props["resource"] = resource
	f.Context = HIDE_ON_LIST
	return f
}

// Bu fonksiyon, polimorfik koleksiyon (MorphMany ilişkisi) oluşturur ve yapılandırır.
//
// Kullanım Senaryosu:
// - Bir kaynağın birden fazla farklı kaynak türüne ait olabileceği polimorfik bire-çok ilişkiyi temsil eder
// - Örn: Bir ürünün birden fazla açıklaması olabilir, bir yazının birden fazla yorumu olabilir
// - Genellikle listede gizlenir, detay sayfasında ilişkili kayıtların listesi gösterilir
//
// Parametreler:
//   - name: Alanın görüntü adı (örn: "Açıklamalar", "Yorumlar")
//   - resource: İlişkili kaynağın adı (örn: "descriptions", "comments")
//   - attribute: İsteğe bağlı veritabanı sütun adı
//
// Dönüş Değeri:
//   - *Schema: Yapılandırılmış polimorfik koleksiyon alan şeması pointer'ı
//
// Örnek Kullanım:
//
//	polyCollection := PolyCollection("Açıklamalar", "descriptions")
//	commentsCollection := PolyCollection("Yorumlar", "comments", "comment_ids")
//
// Önemli Notlar:
//   - Frontend bileşeni "poly-collection-field" olarak ayarlanır
//   - TYPE_POLY_COLLECTION sabiti ile alan türü belirtilir
//   - Context HIDE_ON_LIST ile liste görünümünde gizlenir
func PolyCollection(name string, resource string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "poly-collection-field"
	f.Type = TYPE_POLY_COLLECTION
	f.Props["resource"] = resource
	f.Context = HIDE_ON_LIST
	return f
}

// Bu fonksiyon, polimorfik çoktan çoka (BelongsToMany ilişkisi) ilişki kurmak için kullanılır ve yapılandırır.
//
// Kullanım Senaryosu:
// - Bir kaynağın birden fazla farklı kaynak türüne çoktan çoka ilişkiyi temsil eder
// - Örn: Bir ürünün birden fazla kategoriye ait olabilir, bir yazının birden fazla etiketi olabilir
// - Genellikle ara tablo (pivot table) ile yönetilir
// - Listede gizlenir, detay sayfasında ilişkili kayıtları seçmek için kullanılır
//
// Parametreler:
//   - name: Alanın görüntü adı (örn: "Kategoriler", "Etiketler")
//   - resource: İlişkili kaynağın adı (örn: "categories", "tags")
//   - attribute: İsteğe bağlı veritabanı sütun adı
//
// Dönüş Değeri:
//   - *Schema: Yapılandırılmış polimorfik bağlantı alan şeması pointer'ı
//
// Örnek Kullanım:
//
//	polyConnect := PolyConnect("Kategoriler", "categories")
//	tagsConnect := PolyConnect("Etiketler", "tags", "tag_ids")
//
// Önemli Notlar:
//   - Frontend bileşeni "poly-connect-field" olarak ayarlanır
//   - TYPE_POLY_CONNECT sabiti ile alan türü belirtilir
//   - Context HIDE_ON_LIST ile liste görünümünde gizlenir
//   - Detay sayfasında çoklu seçim yapılabilir
func PolyConnect(name string, resource string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "poly-connect-field"
	f.Type = TYPE_POLY_CONNECT
	f.Props["resource"] = resource
	f.Context = HIDE_ON_LIST
	return f
}

// ============================================================================
// SEÇIM VE DURUM ALANLARI (Selection and Status Fields)
// ============================================================================

// Bu fonksiyon, boolean değerler için switch/toggle bileşeni oluşturur ve yapılandırır.
//
// Kullanım Senaryosu:
// - Aktif/Pasif, Evet/Hayır, Açık/Kapalı gibi boolean değerleri için
// - Toggle switch bileşeni ile kolay seçim yapılabilir
// - Genellikle veritabanında 0/1 veya true/false olarak saklanır
//
// Parametreler:
//   - name: Alanın görüntü adı (örn: "Aktif", "Yayınla")
//   - attribute: İsteğe bağlı veritabanı sütun adı
//
// Dönüş Değeri:
//   - *Schema: Yapılandırılmış switch alan şeması pointer'ı
//
// Örnek Kullanım:
//
//	activeSwitch := Switch("Aktif")
//	publishSwitch := Switch("Yayınla", "is_published")
//
// Önemli Notlar:
//   - Frontend bileşeni "switch-field" olarak ayarlanır
//   - TYPE_BOOLEAN sabiti ile alan türü belirtilir
//   - Genellikle checkbox yerine toggle switch kullanılır
//   - Varsayılan değer (true/false) ayarlanabilir
func Switch(name string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "switch-field"
	f.LabelText = name
	f.Type = TYPE_BOOLEAN
	return f
}

// Bu fonksiyon, çoktan seçmeli veya arama yapılabilir seçim alanı oluşturur ve yapılandırır.
//
// Kullanım Senaryosu:
// - Uzun seçenekler listesinden arama yaparak seçim yapmak için
// - Otomatik tamamlama (autocomplete) özelliği ile hızlı seçim
// - Genellikle API'den dinamik olarak yüklenen seçenekler için
//
// Parametreler:
//   - name: Alanın görüntü adı (örn: "Ürün", "Müşteri")
//   - attribute: İsteğe bağlı veritabanı sütun adı
//
// Dönüş Değeri:
//   - *Schema: Yapılandırılmış combobox alan şeması pointer'ı
//
// Örnek Kullanım:
//
//	productCombo := Combobox("Ürün")
//	customerCombo := Combobox("Müşteri", "customer_id")
//
// Önemli Notlar:
//   - Frontend bileşeni "combobox-field" olarak ayarlanır
//   - TYPE_SELECT sabiti ile alan türü belirtilir (veya TYPE_COMBOBOX eğer yeni tip gerekiyorsa)
//   - Arama ve filtreleme işlevselliği sağlar
//   - Genellikle Select'ten daha gelişmiş bir bileşendir
func Combobox(name string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "combobox-field"
	f.Type = TYPE_SELECT // veya TYPE_COMBOBOX eğer yeni tip gerekiyorsa, şimdilik Select mantığında.
	return f
}

// Bu fonksiyon, standart seçim listesi oluşturur ve yapılandırır.
//
// Kullanım Senaryosu:
// - Önceden tanımlanmış seçeneklerden birini seçmek için
// - Durum, kategori, tür gibi sınırlı sayıda seçeneği olan alanlar için
// - Dropdown veya select bileşeni ile seçim yapılır
//
// Parametreler:
//   - name: Alanın görüntü adı (örn: "Durum", "Tür")
//   - attribute: İsteğe bağlı veritabanı sütun adı
//
// Dönüş Değeri:
//   - *Schema: Yapılandırılmış seçim alan şeması pointer'ı
//
// Örnek Kullanım:
//
//	statusSelect := Select("Durum")
//	typeSelect := Select("Tür", "type_id")
//
// Önemli Notlar:
//   - Frontend bileşeni "select-field" olarak ayarlanır
//   - TYPE_SELECT sabiti ile alan türü belirtilir
//   - Seçenekler Props haritasına "options" anahtarı ile eklenebilir
//   - Varsayılan değer ayarlanabilir
func Select(name string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "select-field"
	f.Type = TYPE_SELECT
	return f
}

// Bu fonksiyon, badge/status gösterim alanı oluşturur ve yapılandırır.
//
// Kullanım Senaryosu:
// - Durum göstermek için (örn: "Aktif", "Beklemede", "Tamamlandı")
// - Renk kodlu etiketler ile görsel durum gösterimi
// - Genellikle salt okunur (read-only) alan olarak kullanılır
//
// Parametreler:
//   - name: Alanın görüntü adı (örn: "Durum", "Etiket")
//   - attribute: İsteğe bağlı veritabanı sütun adı
//
// Dönüş Değeri:
//   - *Schema: Yapılandırılmış badge alan şeması pointer'ı
//
// Örnek Kullanım:
//
//	statusBadge := Badge("Durum")
//	tagBadge := Badge("Etiket", "tag")
//
// Önemli Notlar:
//   - Frontend bileşeni "badge-field" olarak ayarlanır
//   - TYPE_BADGE sabiti ile alan türü belirtilir
//   - Renk ve stil Props haritasına eklenebilir
//   - Genellikle liste görünümünde gösterilir
func Badge(name string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "badge-field"
	f.Type = TYPE_BADGE
	return f
}

// Bu fonksiyon, kod editörü alanı oluşturur ve yapılandırır.
//
// Kullanım Senaryosu:
// - JSON, JavaScript, SQL, HTML gibi kod yazma ve düzenleme için
// - Syntax highlighting ve kod biçimlendirmesi destekler
// - Genellikle Monaco Editor veya benzer editörler ile entegre edilir
//
// Parametreler:
//   - name: Alanın görüntü adı (örn: "Kod", "SQL Sorgusu")
//   - attribute: İsteğe bağlı veritabanı sütun adı
//
// Dönüş Değeri:
//   - *Schema: Yapılandırılmış kod alan şeması pointer'ı
//
// Örnek Kullanım:
//
//	codeField := Code("Kod")
//	sqlField := Code("SQL Sorgusu", "sql_query")
//
// Önemli Notlar:
//   - Frontend bileşeni "code-field" olarak ayarlanır
//   - TYPE_CODE sabiti ile alan türü belirtilir
//   - Dil (language) Props haritasına eklenebilir (json, javascript, sql, html vb.)
//   - Tema (theme) ayarlanabilir (light, dark vb.)
func Code(name string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "code-field"
	f.Type = TYPE_CODE
	return f
}

// Bu fonksiyon, renk seçici alanı oluşturur ve yapılandırır.
//
// Kullanım Senaryosu:
// - Tema rengi, arka plan rengi, vurgu rengi gibi renk seçimi için
// - Hex, RGB, HSL format desteği sağlar
// - Renk paleti veya özel renk seçici ile seçim yapılabilir
//
// Parametreler:
//   - name: Alanın görüntü adı (örn: "Renk", "Tema Rengi")
//   - attribute: İsteğe bağlı veritabanı sütun adı
//
// Dönüş Değeri:
//   - *Schema: Yapılandırılmış renk alan şeması pointer'ı
//
// Örnek Kullanım:
//
//	colorField := Color("Renk")
//	themeColorField := Color("Tema Rengi", "theme_color")
//
// Önemli Notlar:
//   - Frontend bileşeni "color-field" olarak ayarlanır
//   - TYPE_COLOR sabiti ile alan türü belirtilir
//   - Hex formatı (#RRGGBB) standart olarak kullanılır
//   - Alfa (opacity) değeri eklenebilir
func Color(name string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "color-field"
	f.Type = TYPE_COLOR
	return f
}

// Bu fonksiyon, birden fazla boolean seçeneği checkbox group olarak gösterir ve yapılandırır.
//
// Kullanım Senaryosu:
// - Birden fazla özelliği aynı anda seçmek için (örn: İzinler, Tercihler)
// - Her seçenek bağımsız olarak seçilebilir veya seçimi kaldırılabilir
// - Genellikle JSON array veya virgülle ayrılmış string olarak saklanır
//
// Parametreler:
//   - name: Alanın görüntü adı (örn: "İzinler", "Tercihler")
//   - attribute: İsteğe bağlı veritabanı sütun adı
//
// Dönüş Değeri:
//   - *Schema: Yapılandırılmış boolean grup alan şeması pointer'ı
//
// Örnek Kullanım:
//
//	permissionsGroup := BooleanGroup("İzinler")
//	preferencesGroup := BooleanGroup("Tercihler", "user_preferences")
//
// Önemli Notlar:
//   - Frontend bileşeni "boolean-group-field" olarak ayarlanır
//   - TYPE_BOOLEAN_GROUP sabiti ile alan türü belirtilir
//   - Seçenekler Props haritasına "options" anahtarı ile eklenebilir
//   - Genellikle JSON array olarak veritabanında saklanır
func BooleanGroup(name string, attribute ...string) *Schema {
	f := NewField(name, attribute...)
	f.View = "boolean-group-field"
	f.Type = TYPE_BOOLEAN_GROUP
	return f
}
