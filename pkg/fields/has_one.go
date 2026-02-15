package fields

import (
	"reflect"
)

// HasOneField, one-to-one ilişkiyi temsil eder (örn. User -> Profile).
//
// HasOne ilişkisi, bir kaydın tek bir ilişkili kayda sahip olduğunu belirtir.
// Bu, veritabanında ilişkili tabloda foreign key ile temsil edilir.
//
// # Kullanım Senaryoları
//
// - **User -> Profile**: Bir kullanıcının bir profili vardır
// - **Country -> Capital**: Bir ülkenin bir başkenti vardır
// - **Order -> Invoice**: Bir siparişin bir faturası vardır
//
// # Özellikler
//
// - **Tip Güvenliği**: Resource instance veya string slug kullanılabilir
// - **Foreign Key Özelleştirme**: İlişkili tablodaki foreign key sütunu özelleştirilebilir
// - **Owner Key Özelleştirme**: Ana tablodaki referans sütunu özelleştirilebilir
// - **Eager/Lazy Loading**: Yükleme stratejisi seçimi
// - **GORM Yapılandırması**: Foreign key ve references özelleştirme
// - **Hover Card**: Index ve detail sayfalarında hover card desteği
//
// # Kullanım Örneği
//
//	// String slug ile
//	field := fields.HasOne("Profile", "profile", "profiles").
//	    ForeignKey("user_id").
//	    WithEagerLoad()
//
//	// Resource instance ile (tip güvenli)
//	field := fields.HasOne("Profile", "profile", user.NewProfileResource()).
//	    ForeignKey("user_id").
//	    WithEagerLoad()
//
//	// Hover card ile
//	field := fields.HasOne("Profile", "profile", "profiles").
//	    ForeignKey("user_id").
//	    WithHoverCard(fields.NewHoverCardConfig().
//	        WithAvatar("avatar", "").
//	        WithGrid([]fields.HoverCardGridField{
//	            {Key: "bio", Label: "Bio", Type: "text"},
//	            {Key: "location", Label: "Konum", Type: "text", Icon: "map-pin"},
//	        }, "2-column"))
//
// Daha fazla bilgi için docs/Relationships.md dosyasına bakın.
type HasOneField struct {
	Schema
	RelatedResourceSlug string
	RelatedResource     interface{} // resource.Resource interface (interface{} to avoid circular import)
	ForeignKeyColumn    string
	OwnerKeyColumn      string
	QueryCallback       func(query interface{}) interface{}
	LoadingStrategy     LoadingStrategy
	GormRelationConfig  *RelationshipGormConfig
	hoverCardConfig     *HoverCardConfig
}

// HasOne, yeni bir HasOne ilişki alanı oluşturur.
//
// Bu fonksiyon, hem string slug hem de resource instance kabul eder.
// Resource instance kullanımı tip güvenliği sağlar ve refactoring'i kolaylaştırır.
//
// # Parametreler
//
// - **name**: Alanın görünen adı (örn. "Profile", "Profil")
// - **key**: İlişki key'i (örn. "profile")
// - **relatedResource**: İlgili resource (string slug veya resource instance)
//
// # String Slug Kullanımı
//
//	field := fields.HasOne("Profile", "profile", "profiles")
//
// # Resource Instance Kullanımı (Önerilen)
//
//	field := fields.HasOne("Profile", "profile", user.NewProfileResource())
//
// **Avantajlar:**
// - ✅ Tip güvenliği (derleme zamanı kontrolü)
// - ✅ Refactoring desteği
// - ✅ IDE desteği (autocomplete, go-to-definition)
//
// # Varsayılan Değerler
//
// - **ForeignKeyColumn**: slug + "_id" (örn. "profiles_id")
// - **OwnerKeyColumn**: "id" (ana tablonun primary key'i)
// - **LoadingStrategy**: EAGER_LOADING (N+1 sorgu problemini önler)
//
// Döndürür:
//   - Yapılandırılmış HasOneField pointer'ı
func HasOne(name, key string, relatedResource interface{}) *HasOneField {
	// Resource interface'inden slug'ı al
	type resourceSlugger interface {
		Slug() string
	}

	var slug string
	var resourceInstance interface{}

	// Check if relatedResource is a string or a resource instance
	if slugStr, ok := relatedResource.(string); ok {
		// String slug provided
		slug = slugStr
	} else if res, ok := relatedResource.(resourceSlugger); ok {
		// Resource instance provided
		slug = res.Slug()
		resourceInstance = relatedResource
	} else {
		// Fallback: empty slug
		slug = ""
	}

	h := &HasOneField{
		Schema: Schema{
			LabelText: name,
			Name:      name,
			Key:       key,
			View:      "has-one-field",
			Type:      TYPE_RELATIONSHIP,
			Props:     make(map[string]interface{}),
		},
		RelatedResourceSlug: slug,
		RelatedResource:     resourceInstance,
		ForeignKeyColumn:    slug + "_id",
		OwnerKeyColumn:      "id",
		LoadingStrategy:     EAGER_LOADING,
		GormRelationConfig: NewRelationshipGormConfig().
			WithForeignKey(slug + "_id").
			WithReferences("id"),
	}
	// Store relationship details in props for generic access (when Schema interface is used)
	h.WithProps("related_resource", slug)
	if resourceInstance != nil {
		h.WithProps("related_resource_instance", resourceInstance)
	}
	h.WithProps("foreign_key", h.ForeignKeyColumn)
	return h
}

// AutoOptions, ilişkili tablodan otomatik seçenek üretimini etkinleştirir.
//
// Bu method, HasOne ilişkisinde ilişkili kayıtların otomatik olarak seçenek listesi
// olarak sunulmasını sağlar. Frontend'de dropdown veya select bileşenlerinde kullanılır.
//
// # Parametreler
//
// - **displayField**: Seçenek etiketinde gösterilecek sütun adı (örn. "name", "title")
//
// # Kullanım Senaryoları
//
// - **Profil Seçimi**: Kullanıcı profili seçerken profil adlarını gösterme
// - **Başkent Seçimi**: Ülke başkentlerini dropdown'da listeleme
// - **Fatura Seçimi**: Sipariş için fatura seçerken fatura numaralarını gösterme
//
// # Kullanım Örneği
//
//	field := fields.HasOne("Profile", "profile", "profiles").
//	    AutoOptions("full_name").
//	    ForeignKey("user_id")
//
// # Önemli Notlar
//
// - displayField, ilişkili tabloda mevcut bir sütun olmalıdır
// - Büyük veri setlerinde performans sorunlarına neden olabilir
// - Sayfalama ve arama özellikleri ile birlikte kullanılması önerilir
//
// Döndürür:
//   - Yapılandırılmış HasOneField pointer'ı (method chaining için)
func (h *HasOneField) AutoOptions(displayField string) *HasOneField {
	h.Schema.AutoOptions(displayField)
	return h
}

// ForeignKey, foreign key sütun adını ayarlar.
//
// Bu method, ilişkili tablodaki foreign key sütununun adını özelleştirmenizi sağlar.
// Varsayılan olarak, foreign key "{slug}_id" formatında oluşturulur.
//
// # Parametreler
//
// - **key**: Foreign key sütun adı (örn. "user_id", "owner_id", "parent_id")
//
// # Kullanım Senaryoları
//
// - **Özel İsimlendirme**: Standart dışı foreign key isimleri kullanma
// - **Legacy Veritabanları**: Mevcut veritabanı şemalarına uyum sağlama
// - **Çoklu İlişkiler**: Aynı tabloya birden fazla ilişki tanımlama
//
// # Kullanım Örneği
//
//	// Standart kullanım
//	field := fields.HasOne("Profile", "profile", "profiles").
//	    ForeignKey("user_id")
//
//	// Özel isimlendirme
//	field := fields.HasOne("Manager", "manager", "users").
//	    ForeignKey("manager_user_id")
//
//	// Legacy veritabanı
//	field := fields.HasOne("Invoice", "invoice", "invoices").
//	    ForeignKey("order_fk")
//
// # GORM Entegrasyonu
//
// Bu ayar, GORM'un `foreignKey` tag'i ile senkronize edilir:
//
//	type User struct {
//	    ID      uint
//	    Profile Profile `gorm:"foreignKey:user_id"`
//	}
//
// # Önemli Notlar
//
// - Foreign key sütunu ilişkili tabloda mevcut olmalıdır
// - Sütun tipi, owner key ile uyumlu olmalıdır (genellikle uint veya int)
// - Index eklenmesi performans için önerilir
// - Foreign key constraint'leri veritabanı seviyesinde tanımlanmalıdır
//
// Döndürür:
//   - Yapılandırılmış HasOneField pointer'ı (method chaining için)
func (h *HasOneField) ForeignKey(key string) *HasOneField {
	h.ForeignKeyColumn = key
	h.WithProps("foreign_key", key)
	return h
}

// OwnerKey, owner key sütun adını ayarlar.
//
// Bu method, ana tablodaki (parent table) referans sütununun adını özelleştirmenizi sağlar.
// Varsayılan olarak, owner key "id" (primary key) olarak ayarlanır.
//
// # Parametreler
//
// - **key**: Owner key sütun adı (örn. "id", "uuid", "code")
//
// # Kullanım Senaryoları
//
// - **UUID Primary Key**: UUID kullanan tablolarda
// - **Composite Keys**: Birleşik anahtar kullanımında
// - **Custom Primary Keys**: Özel primary key isimlendirmelerinde
// - **Legacy Sistemler**: Standart dışı primary key yapılarında
//
// # Kullanım Örneği
//
//	// UUID primary key
//	field := fields.HasOne("Profile", "profile", "profiles").
//	    ForeignKey("user_uuid").
//	    OwnerKey("uuid")
//
//	// Özel primary key
//	field := fields.HasOne("Invoice", "invoice", "invoices").
//	    ForeignKey("order_code").
//	    OwnerKey("order_code")
//
//	// Composite key senaryosu
//	field := fields.HasOne("Detail", "detail", "details").
//	    ForeignKey("parent_id").
//	    OwnerKey("custom_id")
//
// # GORM Entegrasyonu
//
// Bu ayar, GORM'un `references` tag'i ile senkronize edilir:
//
//	type User struct {
//	    UUID    string  `gorm:"primaryKey"`
//	    Profile Profile `gorm:"foreignKey:user_uuid;references:uuid"`
//	}
//
// # Önemli Notlar
//
// - Owner key sütunu ana tabloda mevcut olmalıdır
// - Genellikle primary key veya unique key olmalıdır
// - Foreign key ile aynı veri tipinde olmalıdır
// - Index'lenmiş olması performans için kritiktir
//
// # Uyarılar
//
// - Owner key değiştirilirse, foreign key ile uyumlu olduğundan emin olun
// - Non-unique owner key kullanımı veri tutarsızlığına yol açabilir
// - Migration'larda bu ilişkiyi doğru tanımladığınızdan emin olun
//
// Döndürür:
//   - Yapılandırılmış HasOneField pointer'ı (method chaining için)
func (h *HasOneField) OwnerKey(key string) *HasOneField {
	h.OwnerKeyColumn = key
	return h
}

// Query, ilişki sorgusunu özelleştirmek için callback fonksiyonu ayarlar.
//
// Bu method, ilişkili kayıtlar yüklenirken uygulanacak özel sorgu koşullarını
// tanımlamanızı sağlar. Filtreleme, sıralama, eager loading gibi işlemler için kullanılır.
//
// # Parametreler
//
// - **fn**: Sorgu callback fonksiyonu - GORM query nesnesini alır ve değiştirilmiş query döndürür
//
// # Kullanım Senaryoları
//
// - **Filtreleme**: Sadece aktif kayıtları yükleme
// - **Sıralama**: İlişkili kayıtları belirli bir sıraya göre getirme
// - **Eager Loading**: İlişkili kayıtların alt ilişkilerini de yükleme
// - **Soft Delete**: Silinmiş kayıtları dahil etme veya hariç tutma
// - **Performans**: Select ile sadece gerekli sütunları çekme
//
// # Kullanım Örnekleri
//
//	// Sadece aktif profilleri yükle
//	field := fields.HasOne("Profile", "profile", "profiles").
//	    Query(func(q interface{}) interface{} {
//	        return q.(*gorm.DB).Where("status = ?", "active")
//	    })
//
//	// Profili ülke bilgisi ile birlikte yükle
//	field := fields.HasOne("Profile", "profile", "profiles").
//	    Query(func(q interface{}) interface{} {
//	        return q.(*gorm.DB).Preload("Country")
//	    })
//
//	// Sadece belirli sütunları çek (performans optimizasyonu)
//	field := fields.HasOne("Profile", "profile", "profiles").
//	    Query(func(q interface{}) interface{} {
//	        return q.(*gorm.DB).Select("id", "full_name", "avatar")
//	    })
//
//	// Sıralama ekle
//	field := fields.HasOne("Invoice", "invoice", "invoices").
//	    Query(func(q interface{}) interface{} {
//	        return q.(*gorm.DB).Order("created_at DESC")
//	    })
//
//	// Birden fazla koşul
//	field := fields.HasOne("Profile", "profile", "profiles").
//	    Query(func(q interface{}) interface{} {
//	        db := q.(*gorm.DB)
//	        return db.Where("verified = ?", true).
//	            Where("status = ?", "active").
//	            Order("updated_at DESC")
//	    })
//
// # GORM Entegrasyonu
//
// Callback fonksiyonu, GORM'un *gorm.DB tipini alır ve değiştirir:
//
//	func(q interface{}) interface{} {
//	    db := q.(*gorm.DB)
//	    // GORM query chain'i burada
//	    return db.Where(...).Order(...)
//	}
//
// # Önemli Notlar
//
// - Callback fonksiyonu her ilişki yüklemesinde çalışır
// - Type assertion kullanırken dikkatli olun (panic riski)
// - Performans için sadece gerekli sütunları Select ile çekin
// - N+1 sorgu problemini önlemek için Preload kullanın
//
// # Performans İpuçları
//
// - **Select**: Sadece gerekli sütunları çekin
// - **Index**: Where koşullarında kullanılan sütunları index'leyin
// - **Preload**: Alt ilişkileri tek sorguda yükleyin
// - **Omit**: Gereksiz büyük sütunları (BLOB, TEXT) hariç tutun
//
// # Uyarılar
//
// - Callback içinde panic oluşursa uygulama çökebilir
// - Type assertion başarısız olursa runtime error oluşur
// - Karmaşık sorgular performansı etkileyebilir
// - Callback'te yapılan değişiklikler tüm yüklemeleri etkiler
//
// Döndürür:
//   - Yapılandırılmış HasOneField pointer'ı (method chaining için)
func (h *HasOneField) Query(fn func(interface{}) interface{}) *HasOneField {
	h.QueryCallback = fn
	return h
}

// WithEagerLoad, yükleme stratejisini eager loading olarak ayarlar.
//
// Eager loading, ilişkili kayıtların ana kayıtlarla birlikte tek sorguda yüklenmesini sağlar.
// Bu, N+1 sorgu problemini önler ve performansı önemli ölçüde artırır.
//
// # Eager Loading Nedir?
//
// Eager loading, ilişkili verilerin önceden (eager) yüklenmesi anlamına gelir.
// Ana kayıtlar sorgulanırken, ilişkili kayıtlar da aynı anda JOIN veya ayrı
// bir sorgu ile yüklenir ve cache'lenir.
//
// # Kullanım Senaryoları
//
// - **Liste Görünümleri**: Birden fazla kaydın ilişkileriyle birlikte gösterilmesi
// - **API Responses**: İlişkili verilerin tek response'da dönülmesi
// - **Raporlama**: Toplu veri çekimlerinde performans optimizasyonu
// - **Dashboard**: Özet bilgilerin hızlı yüklenmesi
// - **Export İşlemleri**: Büyük veri setlerinin verimli işlenmesi
//
// # Kullanım Örneği
//
//	// Eager loading ile (önerilen)
//	field := fields.HasOne("Profile", "profile", "profiles").
//	    WithEagerLoad()
//
//	// Kullanıcıları profilleriyle birlikte yükle
//	// SQL: SELECT * FROM users; SELECT * FROM profiles WHERE user_id IN (1,2,3,...)
//	users := []User{}
//	db.Preload("Profile").Find(&users)
//
// # Performans Karşılaştırması
//
// **Lazy Loading (N+1 Problem):**
//
//	// 1 sorgu: Kullanıcıları getir
//	users := []User{} // 100 kullanıcı
//	db.Find(&users)
//
//	// 100 sorgu: Her kullanıcı için profil getir
//	for _, user := range users {
//	    db.Model(&user).Association("Profile").Find(&user.Profile)
//	}
//	// Toplam: 101 sorgu
//
// **Eager Loading:**
//
//	// 2 sorgu: Kullanıcıları ve tüm profilleri getir
//	users := []User{}
//	db.Preload("Profile").Find(&users)
//	// Toplam: 2 sorgu (veya JOIN ile 1 sorgu)
//
// # GORM Entegrasyonu
//
// Eager loading, GORM'un Preload fonksiyonu ile çalışır:
//
//	db.Preload("Profile").Find(&users)
//	db.Preload("Profile.Country").Find(&users) // İç içe ilişkiler
//
// # Avantajlar
//
// - **Performans**: N+1 sorgu problemini önler
// - **Hız**: Tek seferde tüm veriler yüklenir
// - **Verimlilik**: Veritabanı bağlantı sayısını azaltır
// - **Öngörülebilirlik**: Sorgu sayısı sabittir
// - **Cache Dostu**: Tüm veriler bellekte hazır
//
// # Dezavantajlar
//
// - **Bellek Kullanımı**: Tüm ilişkili veriler belleğe yüklenir
// - **Gereksiz Veri**: Kullanılmayacak veriler de yüklenebilir
// - **İlk Yükleme**: İlk sorgu daha yavaş olabilir
// - **Karmaşık Sorgular**: Çok sayıda ilişki JOIN karmaşıklığı artırır
//
// # Ne Zaman Kullanılmalı?
//
// **Eager Loading Kullan:**
// - Liste görünümlerinde (index, table)
// - İlişkili veri kesinlikle gerekli olduğunda
// - Birden fazla kayıt işlenirken
// - API response'larında
// - Export/raporlama işlemlerinde
//
// **Lazy Loading Kullan:**
// - Tek kayıt görünümlerinde (show, detail)
// - İlişkili veri nadiren gerekli olduğunda
// - Bellek kısıtlı ortamlarda
// - Çok büyük ilişkili veri setlerinde
//
// # Önemli Notlar
//
// - Varsayılan strateji EAGER_LOADING'dir (N+1 problemini önlemek için)
// - Çok sayıda ilişki için selective preloading kullanın
// - İç içe ilişkiler için nokta notasyonu kullanın: "Profile.Country"
// - Conditional preloading için Query callback kullanın
//
// # İleri Seviye Kullanım
//
//	// Koşullu eager loading
//	field := fields.HasOne("Profile", "profile", "profiles").
//	    WithEagerLoad().
//	    Query(func(q interface{}) interface{} {
//	        return q.(*gorm.DB).Where("verified = ?", true)
//	    })
//
//	// İç içe ilişkiler
//	db.Preload("Profile.Country").
//	   Preload("Profile.Avatar").
//	   Find(&users)
//
// Döndürür:
//   - Yapılandırılmış HasOneField pointer'ı (method chaining için)
func (h *HasOneField) WithEagerLoad() *HasOneField {
	h.LoadingStrategy = EAGER_LOADING
	return h
}

// WithLazyLoad, yükleme stratejisini lazy loading olarak ayarlar.
//
// Lazy loading, ilişkili kayıtların sadece gerektiğinde (on-demand) yüklenmesini sağlar.
// Bu, ilk sorgu performansını artırır ancak N+1 sorgu problemine yol açabilir.
//
// # Lazy Loading Nedir?
//
// Lazy loading, ilişkili verilerin geç (lazy) yüklenmesi anlamına gelir.
// Ana kayıtlar sorgulanırken ilişkili kayıtlar yüklenmez, sadece ilişkiye
// erişildiğinde ayrı bir sorgu ile yüklenir.
//
// # Kullanım Senaryoları
//
// - **Detay Görünümleri**: Tek kaydın detaylı gösteriminde
// - **Koşullu Yükleme**: İlişkili veri sadece belirli durumlarda gerekli
// - **Bellek Optimizasyonu**: Büyük ilişkili veri setlerinde
// - **API Pagination**: Sayfalanmış sonuçlarda gereksiz veri yüklemesini önleme
// - **Selective Loading**: Kullanıcı tercihine göre veri yükleme
//
// # Kullanım Örneği
//
//	// Lazy loading ile
//	field := fields.HasOne("Profile", "profile", "profiles").
//	    WithLazyLoad()
//
//	// Kullanıcıyı yükle (profil yüklenmez)
//	user := User{}
//	db.First(&user, 1)
//
//	// Profil sadece erişildiğinde yüklenir
//	db.Model(&user).Association("Profile").Find(&user.Profile)
//
// # N+1 Sorgu Problemi
//
// Lazy loading'in en büyük dezavantajı N+1 sorgu problemidir:
//
//	// 1 sorgu: 100 kullanıcı getir
//	users := []User{}
//	db.Find(&users)
//
//	// 100 sorgu: Her kullanıcı için ayrı profil sorgusu
//	for _, user := range users {
//	    db.Model(&user).Association("Profile").Find(&user.Profile)
//	    // Her iterasyonda yeni bir sorgu!
//	}
//	// Toplam: 101 sorgu (çok yavaş!)
//
// # Performans Karşılaştırması
//
// **Lazy Loading:**
// - İlk yükleme: Hızlı (sadece ana kayıtlar)
// - İlişki erişimi: Yavaş (her erişimde yeni sorgu)
// - Toplam sorgu: N+1 (N = kayıt sayısı)
// - Bellek: Az (sadece gerekli veriler)
//
// **Eager Loading:**
// - İlk yükleme: Yavaş (tüm veriler)
// - İlişki erişimi: Hızlı (zaten yüklü)
// - Toplam sorgu: 1-2 (sabit)
// - Bellek: Fazla (tüm veriler)
//
// # Avantajlar
//
// - **İlk Yükleme Hızı**: Ana kayıtlar hızlı yüklenir
// - **Bellek Verimliliği**: Sadece gerekli veriler yüklenir
// - **Esneklik**: İlişki isteğe bağlı yüklenir
// - **Basitlik**: Otomatik çalışır, ekstra kod gerekmez
// - **Koşullu Yükleme**: Sadece gerektiğinde yükleme
//
// # Dezavantajlar
//
// - **N+1 Problem**: Çok sayıda ek sorgu oluşur
// - **Performans**: Liste görünümlerinde çok yavaş
// - **Öngörülemezlik**: Sorgu sayısı değişkendir
// - **Veritabanı Yükü**: Çok sayıda bağlantı açılır
// - **Debugging**: Sorgu sayısını takip etmek zor
//
// # Ne Zaman Kullanılmalı?
//
// **Lazy Loading Kullan:**
// - Tek kayıt görünümlerinde (show, detail, edit)
// - İlişkili veri nadiren gerekli olduğunda
// - Çok büyük ilişkili veri setlerinde
// - Bellek kısıtlı ortamlarda
// - Koşullu veri yükleme senaryolarında
//
// **Eager Loading Kullan:**
// - Liste görünümlerinde (index, table, grid)
// - İlişkili veri kesinlikle gerekli olduğunda
// - Birden fazla kayıt işlenirken
// - API response'larında
// - Export/raporlama işlemlerinde
//
// # N+1 Problemini Önleme
//
// Lazy loading kullanırken N+1 problemini önlemek için:
//
//	// 1. Manuel Preload kullan
//	db.Preload("Profile").Find(&users)
//
//	// 2. Batch loading kullan
//	db.Model(&users).Association("Profile").Find(&profiles)
//
//	// 3. Eager loading'e geç
//	field.WithEagerLoad()
//
// # Önemli Notlar
//
// - Liste görünümlerinde lazy loading kullanmayın (N+1 riski)
// - Tek kayıt görünümlerinde lazy loading tercih edilebilir
// - Production'da sorgu sayısını mutlaka monitor edin
// - GORM'un Debug() modu ile sorguları kontrol edin
//
// # İleri Seviye Kullanım
//
//	// Koşullu lazy loading
//	field := fields.HasOne("Profile", "profile", "profiles").
//	    WithLazyLoad()
//
//	// Runtime'da eager loading'e geç
//	if needsProfile {
//	    db.Preload("Profile").Find(&users)
//	} else {
//	    db.Find(&users)
//	}
//
// # Debugging İpuçları
//
//	// Sorgu sayısını kontrol et
//	db.Debug().Find(&users) // SQL loglarını gösterir
//
//	// N+1 tespiti için
//	// Her iterasyonda "SELECT * FROM profiles" görüyorsanız N+1 var!
//
// # Uyarılar
//
// - Production'da lazy loading kullanırken dikkatli olun
// - N+1 problemi ciddi performans sorunlarına yol açar
// - Büyük veri setlerinde veritabanını aşırı yükleyebilir
// - Liste görünümlerinde kesinlikle eager loading kullanın
//
// Döndürür:
//   - Yapılandırılmış HasOneField pointer'ı (method chaining için)
func (h *HasOneField) WithLazyLoad() *HasOneField {
	h.LoadingStrategy = LAZY_LOADING
	return h
}

// Extract, resource'dan değeri çıkarır ve işler.
//
// Bu method, HasOne ilişkisinde ilişkili resource'un ID değerini çıkarır.
// Eğer ilişkili veri bir struct ise, struct'ın ID alanını bulup değer olarak ayarlar.
// Bu, frontend'e gönderilecek veriyi hazırlar ve sadece ID'yi döndürür.
//
// # Çalışma Mantığı
//
// 1. Schema.Extract() ile temel extraction yapılır
// 2. Eğer veri nil ise, işlem sonlandırılır
// 3. Eğer veri pointer ise, dereference edilir
// 4. Eğer veri struct ise, ID veya Id alanı aranır
// 5. Bulunan ID değeri Data alanına atanır
//
// # Kullanım Senaryoları
//
// - **API Response**: Frontend'e sadece ilişkili kaydın ID'sini gönderme
// - **Form Data**: Edit formlarında mevcut ilişkili kaydın ID'sini gösterme
// - **Serialization**: JSON response'da nested struct yerine ID kullanma
// - **Data Transformation**: Struct'tan primitive değere dönüşüm
//
// # Örnek Veri Dönüşümü
//
// **Giriş (Struct):**
//
//	user := User{
//	    ID: 1,
//	    Profile: Profile{
//	        ID: 42,
//	        FullName: "John Doe",
//	        Avatar: "avatar.jpg",
//	    },
//	}
//
// **Çıkış (ID):**
//
//	field.Data = 42  // Sadece Profile ID'si
//
// # Desteklenen Veri Tipleri
//
// **Struct:**
//
//	Profile{ID: 42, FullName: "John"} → 42
//
// **Pointer to Struct:**
//
//	&Profile{ID: 42, FullName: "John"} → 42
//
// **Nil Pointer:**
//
//	var profile *Profile = nil → nil
//
// **Primitive (ID zaten var):**
//
//	42 → 42 (değişmez)
//
// # ID Alan İsimleri
//
// Method, aşağıdaki alan isimlerini sırayla arar:
// 1. "ID" (büyük harf, Go convention)
// 2. "Id" (camelCase, alternatif)
//
// # Kullanım Örneği
//
//	// Resource tanımı
//	type User struct {
//	    ID      uint
//	    Profile Profile
//	}
//
//	type Profile struct {
//	    ID       uint
//	    FullName string
//	}
//
//	// Field tanımı
//	field := fields.HasOne("Profile", "profile", "profiles")
//
//	// Extract çağrısı
//	user := User{
//	    ID: 1,
//	    Profile: Profile{ID: 42, FullName: "John"},
//	}
//	field.Extract(user)
//
//	// Sonuç
//	fmt.Println(field.Data) // Output: 42
//
// # JSON Serialization
//
// Extract sonrası JSON response:
//
//	{
//	    "id": 1,
//	    "profile": 42  // Struct yerine sadece ID
//	}
//
// # Reflection Kullanımı
//
// Bu method, Go'nun reflection paketini kullanır:
// - reflect.ValueOf(): Değeri reflection value'ya çevirir
// - v.Kind(): Değerin tipini kontrol eder (Ptr, Struct)
// - v.Elem(): Pointer'ı dereference eder
// - v.FieldByName(): Struct alanını isimle bulur
// - idField.Interface(): Reflection value'yu interface{}'e çevirir
//
// # Önemli Notlar
//
// - ID alanı bulunamazsa, Data değişmeden kalır
// - ID alanı private ise (küçük harfle başlıyorsa) erişilemez
// - ID alanı exported olmalıdır (büyük harfle başlamalı)
// - Nil pointer'lar güvenli şekilde işlenir
//
// # Performans
//
// - Reflection kullanımı nedeniyle normal field access'ten yavaştır
// - Ancak, serialization için gereklidir
// - Büyük veri setlerinde dikkat edilmeli
// - Cache mekanizması ile optimize edilebilir
//
// # Hata Durumları
//
// Method, aşağıdaki durumlarda sessizce başarısız olur:
// - ID alanı bulunamazsa (Data değişmez)
// - ID alanı private ise (erişilemez)
// - Veri tipi desteklenmiyorsa (Data değişmez)
//
// # Alternatif Yaklaşımlar
//
// **Manuel ID Extraction:**
//
//	if user.Profile.ID != 0 {
//	    field.Data = user.Profile.ID
//	}
//
// **Interface Method:**
//
//	type IDGetter interface {
//	    GetID() uint
//	}
//
// **Struct Tag:**
//
//	Profile Profile `json:"profile_id,omitempty"`
//
// # Uyarılar
//
// - Reflection kullanımı runtime overhead ekler
// - Type assertion başarısız olabilir
// - ID alanı yoksa sessizce başarısız olur
// - Panic riski düşük ama mevcut (invalid reflection operations)
//
// Parametreler:
//   - resource: Extract edilecek resource (struct veya pointer)
func (h *HasOneField) Extract(resource interface{}) {
	// Schema.Extract ile ilişki verilerini al
	h.Schema.Extract(resource)

	// Data nil ise çık
	if h.Schema.Data == nil {
		return
	}

	// RelatedResource yoksa mevcut veriyi kullan
	if h.RelatedResource == nil {
		return
	}

	// Data'yı reflection ile işle
	v := reflect.ValueOf(h.Schema.Data)

	// Pointer ise dereference et
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			h.Schema.Data = nil
			return
		}
		v = v.Elem()
	}

	// Struct değilse çık
	if v.Kind() != reflect.Struct {
		return
	}

	record := v.Interface()

	// ID field'ını bul
	var idValue interface{}
	idField := v.FieldByName("ID")
	if !idField.IsValid() {
		idField = v.FieldByName("Id")
	}

	if idField.IsValid() && idField.CanInterface() {
		idValue = idField.Interface()
	}

	// ID bulunamadıysa çık
	if idValue == nil {
		return
	}

	// RelatedResource'dan RecordTitle metodunu çağır (type assertion ile)
	// RelatedResource interface{} tipinde olduğu için type assertion gerekli
	type ResourceWithTitle interface {
		RecordTitle(any) string
	}

	res, ok := h.RelatedResource.(ResourceWithTitle)
	if !ok {
		// RelatedResource RecordTitle metoduna sahip değilse çık
		return
	}

	// RecordTitle ile başlığı al (gerekirse struct field fallback kullan)
	recordTitle := resolveRelationshipRecordTitle(res, record, idValue)

	// Minimal format: {"id": ..., "title": ...}
	h.Schema.Data = map[string]interface{}{
		"id":    idValue,
		"title": recordTitle,
	}
}

// GetRelationshipType, ilişki tipini döndürür.
//
// Bu method, HasOne ilişkisinin tipini string olarak döndürür.
// İlişki tipi, sistem içinde ilişki türünü tanımlamak ve işlemek için kullanılır.
//
// # Dönüş Değeri
//
// - "hasOne": HasOne ilişki tipini belirten sabit string
//
// # Kullanım Senaryoları
//
// - **Tip Kontrolü**: İlişki tipini runtime'da kontrol etme
// - **Routing**: İlişki tipine göre farklı işlemler yapma
// - **Serialization**: JSON/XML çıktısında ilişki tipini belirtme
// - **Validation**: İlişki tipine özel validasyon kuralları uygulama
// - **Logging**: Log mesajlarında ilişki tipini gösterme
//
// # Kullanım Örneği
//
//	field := fields.HasOne("Profile", "profile", "profiles")
//	relType := field.GetRelationshipType()
//	fmt.Println(relType) // Output: "hasOne"
//
//	// Tip kontrolü
//	if field.GetRelationshipType() == "hasOne" {
//	    // HasOne'a özel işlemler
//	}
//
//	// Switch case ile routing
//	switch field.GetRelationshipType() {
//	case "hasOne":
//	    // HasOne işlemleri
//	case "hasMany":
//	    // HasMany işlemleri
//	case "belongsTo":
//	    // BelongsTo işlemleri
//	}
//
// # İlişki Tipleri
//
// Sistemdeki diğer ilişki tipleri:
// - "hasOne": Bir-bir ilişki (bu method)
// - "hasMany": Bir-çok ilişki
// - "belongsTo": Ters bir-bir veya ters bir-çok ilişki
// - "belongsToMany": Çok-çok ilişki
// - "morphTo": Polimorfik ilişki
// - "morphOne": Polimorfik bir-bir ilişki
// - "morphMany": Polimorfik bir-çok ilişki
//
// # JSON Serialization
//
// API response'da kullanımı:
//
//	{
//	    "name": "Profile",
//	    "key": "profile",
//	    "type": "relationship",
//	    "relationship_type": "hasOne",
//	    "related_resource": "profiles"
//	}
//
// # Interface Implementation
//
// Bu method, Relationship interface'ini implement eder:
//
//	type Relationship interface {
//	    GetRelationshipType() string
//	    GetRelatedResource() string
//	    GetRelationshipName() string
//	}
//
// # Önemli Notlar
//
// - Dönüş değeri her zaman "hasOne" string'idir
// - Değer değiştirilemez (immutable)
// - Case-sensitive: "hasOne" (camelCase)
// - Boş string döndürmez
//
// # Performans
//
// - O(1) kompleksitesi (sabit string dönüşü)
// - Bellek allocation yok
// - Çok hızlı çalışır
//
// Döndürür:
//   - "hasOne" string sabiti
func (h *HasOneField) GetRelationshipType() string {
	return "hasOne"
}

// GetRelatedResource, ilişkili resource'un slug'ını döndürür.
//
// Bu method, HasOne ilişkisinde ilişkili resource'un slug değerini döndürür.
// Slug, resource'u sistem içinde tanımlayan benzersiz string identifier'dır.
//
// # Dönüş Değeri
//
// - İlişkili resource'un slug'ı (örn. "profiles", "invoices", "capitals")
//
// # Kullanım Senaryoları
//
// - **Resource Lookup**: İlişkili resource'u slug ile bulma
// - **API Routing**: İlişkili resource'un endpoint'ini oluşturma
// - **Dynamic Loading**: Runtime'da ilişkili resource'u yükleme
// - **Validation**: İlişki tanımının doğruluğunu kontrol etme
// - **Metadata**: İlişki bilgilerini metadata olarak saklama
//
// # Kullanım Örneği
//
//	field := fields.HasOne("Profile", "profile", "profiles")
//	slug := field.GetRelatedResource()
//	fmt.Println(slug) // Output: "profiles"
//
//	// Resource lookup
//	relatedResource := resourceRegistry.Get(field.GetRelatedResource())
//
//	// API endpoint oluşturma
//	endpoint := fmt.Sprintf("/api/%s", field.GetRelatedResource())
//	// Output: "/api/profiles"
//
//	// Dynamic query
//	tableName := field.GetRelatedResource()
//	db.Table(tableName).Where("user_id = ?", userID).First(&profile)
//
// # Slug Formatı
//
// Slug genellikle aşağıdaki formatlarda olur:
// - **Plural Form**: "profiles", "invoices", "users"
// - **Kebab Case**: "user-profiles", "order-items"
// - **Snake Case**: "user_profiles", "order_items"
//
// # Resource Instance vs String Slug
//
// HasOne oluşturulurken iki yöntem kullanılabilir:
//
// **String Slug:**
//
//	field := fields.HasOne("Profile", "profile", "profiles")
//	field.GetRelatedResource() // "profiles"
//
// **Resource Instance:**
//
//	field := fields.HasOne("Profile", "profile", user.NewProfileResource())
//	field.GetRelatedResource() // "profiles" (resource'un Slug() method'undan)
//
// # JSON Serialization
//
// API response'da kullanımı:
//
//	{
//	    "name": "Profile",
//	    "key": "profile",
//	    "type": "relationship",
//	    "relationship_type": "hasOne",
//	    "related_resource": "profiles"  // Bu method'un dönüş değeri
//	}
//
// # Interface Implementation
//
// Bu method, Relationship interface'ini implement eder:
//
//	type Relationship interface {
//	    GetRelationshipType() string
//	    GetRelatedResource() string  // Bu method
//	    GetRelationshipName() string
//	}
//
// # Resource Registry Entegrasyonu
//
//	// Resource registry'den ilişkili resource'u al
//	type ResourceRegistry struct {
//	    resources map[string]Resource
//	}
//
//	func (r *ResourceRegistry) Get(slug string) Resource {
//	    return r.resources[slug]
//	}
//
//	// Kullanım
//	relatedResource := registry.Get(field.GetRelatedResource())
//
// # Önemli Notlar
//
// - Slug değeri HasOne oluşturulurken ayarlanır
// - Değer değiştirilemez (immutable)
// - Boş string olabilir (fallback durumunda)
// - Case-sensitive: "profiles" ≠ "Profiles"
//
// # Hata Durumları
//
// Eğer HasOne oluşturulurken geçersiz bir değer verilirse:
// - String değilse ve Slug() method'u yoksa: boş string döner
// - Nil değer: boş string döner
//
// # Performans
//
// - O(1) kompleksitesi (field access)
// - Bellek allocation yok
// - Çok hızlı çalışır
//
// Döndürür:
//   - İlişkili resource'un slug'ı (string)
func (h *HasOneField) GetRelatedResourceSlug() string {
	return h.RelatedResourceSlug
}

// GetRelationshipName, ilişkinin adını döndürür.
//
// Bu method, HasOne ilişkisinin görünen adını (display name) döndürür.
// Bu ad, kullanıcı arayüzünde gösterilir ve ilişkiyi tanımlar.
//
// # Dönüş Değeri
//
// - İlişkinin görünen adı (örn. "Profile", "Invoice", "Capital")
//
// # Kullanım Senaryoları
//
// - **UI Rendering**: Form etiketlerinde ve tablo başlıklarında gösterme
// - **Validation Messages**: Hata mesajlarında ilişki adını kullanma
// - **Logging**: Log mesajlarında ilişkiyi tanımlama
// - **API Documentation**: API dokümantasyonunda ilişki açıklaması
// - **Debugging**: Debug çıktılarında ilişkiyi belirleme
//
// # Kullanım Örneği
//
//	field := fields.HasOne("Profile", "profile", "profiles")
//	name := field.GetRelationshipName()
//	fmt.Println(name) // Output: "Profile"
//
//	// UI rendering
//	label := fmt.Sprintf("%s:", field.GetRelationshipName())
//	// Output: "Profile:"
//
//	// Validation message
//	errMsg := fmt.Sprintf("%s is required", field.GetRelationshipName())
//	// Output: "Profile is required"
//
//	// Form field
//	<label>{field.GetRelationshipName()}</label>
//	<select name="profile_id">...</select>
//
// # Name vs Key vs Slug
//
// HasOne ilişkisinde üç farklı identifier vardır:
//
// **Name (GetRelationshipName):**
// - Görünen ad, kullanıcı dostu
// - Örnek: "Profile", "User Profile", "Kullanıcı Profili"
// - Kullanım: UI, mesajlar, dokümantasyon
//
// **Key:**
// - Programatik identifier, struct field adı
// - Örnek: "profile", "userProfile", "user_profile"
// - Kullanım: Kod içinde, JSON keys, form names
//
// **Slug (GetRelatedResource):**
// - Resource identifier, tablo/endpoint adı
// - Örnek: "profiles", "user-profiles", "user_profiles"
// - Kullanım: Routing, database, API endpoints
//
// # Çoklu Dil Desteği (i18n)
//
// Name değeri, çoklu dil desteği için kullanılabilir:
//
//	// İngilizce
//	field := fields.HasOne("Profile", "profile", "profiles")
//
//	// Türkçe
//	field := fields.HasOne("Profil", "profile", "profiles")
//
//	// Runtime çeviri
//	name := i18n.Translate(field.GetRelationshipName())
//
// # JSON Serialization
//
// API response'da kullanımı:
//
//	{
//	    "name": "Profile",  // Bu method'un dönüş değeri
//	    "key": "profile",
//	    "type": "relationship",
//	    "relationship_type": "hasOne",
//	    "related_resource": "profiles"
//	}
//
// # Interface Implementation
//
// Bu method, Relationship interface'ini implement eder:
//
//	type Relationship interface {
//	    GetRelationshipType() string
//	    GetRelatedResource() string
//	    GetRelationshipName() string  // Bu method
//	}
//
// # UI Kullanım Örnekleri
//
// **Form Label:**
//
//	<div class="form-group">
//	    <label>{field.GetRelationshipName()}</label>
//	    <select name="profile_id">
//	        <option value="">Select {field.GetRelationshipName()}</option>
//	    </select>
//	</div>
//
// **Table Header:**
//
//	<th>{field.GetRelationshipName()}</th>
//
// **Validation Error:**
//
//	errors := map[string]string{
//	    "profile": fmt.Sprintf("%s is required", field.GetRelationshipName()),
//	}
//
// # Önemli Notlar
//
// - Name değeri HasOne oluşturulurken ayarlanır
// - Değer değiştirilemez (immutable)
// - Boş string olabilir (ancak önerilmez)
// - Kullanıcı dostu olmalıdır (başlık formatında)
//
// # Best Practices
//
// **İyi Örnekler:**
// - "Profile" (kısa, açık)
// - "User Profile" (açıklayıcı)
// - "Invoice Details" (detaylı)
// - "Kullanıcı Profili" (yerelleştirilmiş)
//
// **Kötü Örnekler:**
// - "profile" (küçük harf, programatik)
// - "PROFILE" (tamamı büyük harf)
// - "prof" (kısaltma, belirsiz)
// - "user_profile" (snake_case, programatik)
//
// # Performans
//
// - O(1) kompleksitesi (field access)
// - Bellek allocation yok
// - Çok hızlı çalışır
//
// Döndürür:
//   - İlişkinin görünen adı (string)
func (h *HasOneField) GetRelationshipName() string {
	return h.Name
}

// ResolveRelationship, ilişkiyi çözerek ilişkili kaydı yükler.
//
// Bu method, HasOne ilişkisinde tek bir ilişkili kaydı veritabanından yükler.
// Gerçek implementasyonda, bu method veritabanı sorgusunu çalıştırır ve
// ilişkili resource'u döndürür.
//
// # Çalışma Mantığı
//
// 1. Item parametresini kontrol eder (nil kontrolü)
// 2. Item'dan owner key değerini çıkarır (örn. user.ID)
// 3. İlişkili tabloda foreign key ile sorgu yapar
// 4. Tek bir ilişkili kaydı bulur ve döndürür
// 5. Hata durumunda error döndürür
//
// # Parametreler
//
// - **item**: Ana kayıt (owner record) - ilişkinin sahibi olan kayıt
//
// # Dönüş Değerleri
//
// - **interface{}**: İlişkili kayıt (tek bir record) veya nil
// - **error**: Hata durumunda error, başarılı ise nil
//
// # Kullanım Senaryoları
//
// - **Lazy Loading**: İlişkili kaydı isteğe bağlı yükleme
// - **Dynamic Resolution**: Runtime'da ilişki çözümleme
// - **API Endpoints**: İlişkili veriyi ayrı endpoint'te sunma
// - **Conditional Loading**: Belirli koşullarda ilişki yükleme
// - **Manual Preloading**: Özel preload stratejileri
//
// # Kullanım Örneği (Teorik)
//
//	// Field tanımı
//	field := fields.HasOne("Profile", "profile", "profiles").
//	    ForeignKey("user_id").
//	    OwnerKey("id")
//
//	// Ana kayıt
//	user := User{ID: 1, Name: "John"}
//
//	// İlişkiyi çözümle
//	profile, err := field.ResolveRelationship(user)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Sonuç
//	fmt.Printf("Profile: %+v\n", profile)
//	// Output: Profile: &Profile{ID: 42, UserID: 1, FullName: "John Doe"}
//
// # Gerçek Implementasyon (Örnek)
//
//	func (h *HasOneField) ResolveRelationship(item interface{}) (interface{}, error) {
//	    if item == nil {
//	        return nil, nil
//	    }
//
//	    // Owner key değerini çıkar
//	    ownerValue := extractFieldValue(item, h.OwnerKeyColumn)
//	    if ownerValue == nil {
//	        return nil, nil
//	    }
//
//	    // İlişkili resource'u bul
//	    relatedResource := resourceRegistry.Get(h.RelatedResourceSlug)
//	    if relatedResource == nil {
//	        return nil, fmt.Errorf("related resource not found: %s", h.RelatedResourceSlug)
//	    }
//
//	    // Veritabanı sorgusu
//	    var result interface{}
//	    query := db.Model(relatedResource.Model()).
//	        Where(fmt.Sprintf("%s = ?", h.ForeignKeyColumn), ownerValue)
//
//	    // Query callback uygula
//	    if h.QueryCallback != nil {
//	        query = h.QueryCallback(query).(*gorm.DB)
//	    }
//
//	    // Tek kayıt bul
//	    err := query.First(&result).Error
//	    if err != nil {
//	        if errors.Is(err, gorm.ErrRecordNotFound) {
//	            return nil, nil // İlişkili kayıt yok
//	        }
//	        return nil, err
//	    }
//
//	    return result, nil
//	}
//
// # SQL Sorgusu
//
// Method, aşağıdaki gibi bir SQL sorgusu oluşturur:
//
//	SELECT * FROM profiles
//	WHERE user_id = 1
//	LIMIT 1
//
// # Hata Durumları
//
// **Nil Item:**
// - Input: nil
// - Output: (nil, nil)
//
// **Kayıt Bulunamadı:**
// - Input: User{ID: 999}
// - Output: (nil, nil) // Error değil, sadece nil
//
// **Veritabanı Hatası:**
// - Input: User{ID: 1}
// - Output: (nil, error) // Bağlantı hatası, syntax error vb.
//
// **Birden Fazla Kayıt:**
// - Input: User{ID: 1}
// - Output: İlk kaydı döndürür (LIMIT 1)
//
// # Performans Optimizasyonu
//
// **Select Specific Columns:**
//
//	field.Query(func(q interface{}) interface{} {
//	    return q.(*gorm.DB).Select("id", "full_name", "avatar")
//	})
//
// **Index Kullanımı:**
//
//	CREATE INDEX idx_profiles_user_id ON profiles(user_id);
//
// **Cache Stratejisi:**
//
//	// Cache'den kontrol et
//	if cached := cache.Get(cacheKey); cached != nil {
//	    return cached, nil
//	}
//
//	// Veritabanından yükle
//	result, err := field.ResolveRelationship(item)
//	if err == nil && result != nil {
//	    cache.Set(cacheKey, result, 5*time.Minute)
//	}
//
// # Eager Loading vs Lazy Loading
//
// **Lazy Loading (Bu Method):**
// - Her kayıt için ayrı sorgu
// - N+1 problem riski
// - İsteğe bağlı yükleme
//
// **Eager Loading (Preload):**
// - Tek sorguda tüm ilişkiler
// - Performanslı
// - Toplu yükleme
//
// # Önemli Notlar
//
// - Şu anki implementasyon placeholder'dır (nil, nil döndürür)
// - Gerçek implementasyon veritabanı erişimi gerektirir
// - Query callback uygulanmalıdır
// - Error handling kritiktir
// - Nil item güvenli şekilde işlenir
//
// # Best Practices
//
// - Eager loading tercih edin (liste görünümlerinde)
// - Lazy loading sadece gerektiğinde kullanın
// - Query callback ile filtreleme yapın
// - Index'leri doğru tanımlayın
// - Cache mekanizması kullanın
// - Error'ları loglamayı unutmayın
//
// # Uyarılar
//
// - N+1 sorgu problemine dikkat edin
// - Nil pointer dereference riskine karşı kontrol yapın
// - Veritabanı bağlantı havuzunu yönetin
// - Timeout ayarlarını yapılandırın
// - Transaction içinde çalışırken dikkatli olun
//
// Döndürür:
//   - İlişkili kayıt (interface{}) veya nil
//   - Hata durumunda error, başarılı ise nil
func (h *HasOneField) ResolveRelationship(item interface{}) (interface{}, error) {
	if item == nil {
		return nil, nil
	}

	// In a real implementation, this would query the database
	// For now, return nil
	return nil, nil
}

// ValidateRelationship, ilişkinin geçerliliğini doğrular.
//
// Bu method, HasOne ilişkisinde en fazla bir ilişkili kaydın var olduğunu doğrular.
// Gerçek implementasyonda, veritabanı constraint'lerini ve iş kurallarını kontrol eder.
//
// # Çalışma Mantığı
//
// 1. Value parametresini kontrol eder (nil, tip, format)
// 2. İlişkili kayıt sayısını kontrol eder (en fazla 1 olmalı)
// 3. Foreign key constraint'lerini doğrular
// 4. İş kurallarını uygular (required, unique, vb.)
// 5. Hata durumunda error döndürür
//
// # Parametreler
//
// - **value**: Doğrulanacak değer (genellikle ilişkili kaydın ID'si)
//
// # Dönüş Değeri
//
// - **error**: Validasyon hatası varsa error, geçerli ise nil
//
// # Kullanım Senaryoları
//
// - **Form Validation**: Kullanıcı input'unu doğrulama
// - **API Validation**: Request payload'unu kontrol etme
// - **Data Integrity**: Veri tutarlılığını sağlama
// - **Business Rules**: İş kurallarını uygulama
// - **Constraint Checking**: Veritabanı constraint'lerini kontrol etme
//
// # Kullanım Örneği (Teorik)
//
//	field := fields.HasOne("Profile", "profile", "profiles").
//	    ForeignKey("user_id").
//	    Required()
//
//	// Validasyon
//	err := field.ValidateRelationship(profileID)
//	if err != nil {
//	    return fmt.Errorf("validation failed: %w", err)
//	}
//
// # Gerçek Implementasyon (Örnek)
//
//	func (h *HasOneField) ValidateRelationship(value interface{}) error {
//	    // Nil kontrolü
//	    if value == nil {
//	        if h.IsRequired() {
//	            return fmt.Errorf("%s is required", h.Name)
//	        }
//	        return nil
//	    }
//
//	    // Tip kontrolü
//	    id, ok := value.(uint)
//	    if !ok {
//	        return fmt.Errorf("%s must be a valid ID", h.Name)
//	    }
//
//	    // ID geçerliliği
//	    if id == 0 {
//	        return fmt.Errorf("%s ID cannot be zero", h.Name)
//	    }
//
//	    // İlişkili kayıt var mı?
//	    var count int64
//	    err := db.Model(relatedModel).
//	        Where("id = ?", id).
//	        Count(&count).Error
//	    if err != nil {
//	        return fmt.Errorf("failed to validate %s: %w", h.Name, err)
//	    }
//
//	    if count == 0 {
//	        return fmt.Errorf("%s with ID %d does not exist", h.Name, id)
//	    }
//
//	    // Birden fazla kayıt kontrolü (HasOne için)
//	    if count > 1 {
//	        return fmt.Errorf("%s has multiple records (expected one)", h.Name)
//	    }
//
//	    // Foreign key constraint kontrolü
//	    if h.ForeignKeyColumn != "" {
//	        // Foreign key'in başka bir kayıt tarafından kullanılıp kullanılmadığını kontrol et
//	        var existingCount int64
//	        err := db.Model(relatedModel).
//	            Where(fmt.Sprintf("%s = ?", h.ForeignKeyColumn), ownerID).
//	            Where("id != ?", id).
//	            Count(&existingCount).Error
//	        if err != nil {
//	            return err
//	        }
//
//	        if existingCount > 0 {
//	            return fmt.Errorf("%s already has a related record", h.Name)
//	        }
//	    }
//
//	    return nil
//	}
//
// # Validasyon Kuralları
//
// **Required Kontrolü:**
//
//	if value == nil && field.IsRequired() {
//	    return errors.New("field is required")
//	}
//
// **Tip Kontrolü:**
//
//	if _, ok := value.(uint); !ok {
//	    return errors.New("invalid type")
//	}
//
// **Existence Kontrolü:**
//
//	var count int64
//	db.Model(&Profile{}).Where("id = ?", value).Count(&count)
//	if count == 0 {
//	    return errors.New("record not found")
//	}
//
// **Uniqueness Kontrolü:**
//
//	var count int64
//	db.Model(&Profile{}).Where("user_id = ?", userID).Count(&count)
//	if count > 1 {
//	    return errors.New("multiple records found")
//	}
//
// # Hata Mesajları
//
// **Required Error:**
// - "Profile is required"
// - "Profile cannot be empty"
//
// **Type Error:**
// - "Profile must be a valid ID"
// - "Profile has invalid type"
//
// **Not Found Error:**
// - "Profile with ID 42 does not exist"
// - "Profile not found"
//
// **Duplicate Error:**
// - "Profile already exists for this user"
// - "User already has a profile"
//
// # Custom Validation
//
// Özel validasyon kuralları eklenebilir:
//
//	field := fields.HasOne("Profile", "profile", "profiles").
//	    Validate(func(value interface{}) error {
//	        id := value.(uint)
//	        // Özel validasyon mantığı
//	        if id < 100 {
//	            return errors.New("profile ID must be >= 100")
//	        }
//	        return nil
//	    })
//
// # Veritabanı Constraint'leri
//
// **Foreign Key Constraint:**
//
//	ALTER TABLE profiles
//	ADD CONSTRAINT fk_profiles_user_id
//	FOREIGN KEY (user_id) REFERENCES users(id)
//	ON DELETE CASCADE;
//
// **Unique Constraint:**
//
//	ALTER TABLE profiles
//	ADD CONSTRAINT uk_profiles_user_id
//	UNIQUE (user_id);
//
// # Performans Optimizasyonu
//
// **Cache Kullanımı:**
//
//	// Validation sonuçlarını cache'le
//	cacheKey := fmt.Sprintf("validation:%s:%v", h.Key, value)
//	if cached := cache.Get(cacheKey); cached != nil {
//	    return cached.(error)
//	}
//
// **Batch Validation:**
//
//	// Birden fazla değeri tek sorguda doğrula
//	ids := []uint{1, 2, 3}
//	var count int64
//	db.Model(&Profile{}).Where("id IN ?", ids).Count(&count)
//
// # Önemli Notlar
//
// - Şu anki implementasyon placeholder'dır (nil döndürür)
// - Gerçek implementasyon veritabanı erişimi gerektirir
// - Validasyon kuralları field konfigürasyonuna bağlıdır
// - Error mesajları kullanıcı dostu olmalıdır
// - Performans için cache kullanılabilir
//
// # Best Practices
//
// - Validasyon kurallarını açık ve anlaşılır yapın
// - Error mesajlarını kullanıcı dostu yazın
// - Veritabanı constraint'lerini kullanın
// - Validasyonu hem client hem server tarafında yapın
// - Batch validation ile performansı artırın
//
// # Uyarılar
//
// - Validasyon atlanmamalıdır (güvenlik riski)
// - Error mesajları hassas bilgi içermemelidir
// - Veritabanı hatalarını yakalayın ve loglamayı unutmayın
// - Race condition'lara dikkat edin (concurrent updates)
// - Transaction içinde validasyon yapın
//
// Döndürür:
//   - Validasyon hatası varsa error, geçerli ise nil
func (h *HasOneField) ValidateRelationship(value interface{}) error {
	// Validate that at most one related resource exists
	// In a real implementation, this would check database constraints
	return nil
}

// GetDisplayKey, display key'i döndürür (HasOne için kullanılmaz).
//
// Bu method, HasOne ilişkisinde kullanılmaz ve boş string döndürür.
// Display key, genellikle HasMany ve BelongsToMany gibi çoklu kayıt
// ilişkilerinde hangi alanın gösterileceğini belirtmek için kullanılır.
//
// # Dönüş Değeri
//
// - Boş string ("") - HasOne için display key kavramı geçerli değildir
//
// # Neden Kullanılmaz?
//
// HasOne ilişkisi tek bir kayıt döndürdüğü için, display key'e ihtiyaç yoktur.
// Display key, çoklu kayıtların listesinde hangi alanın gösterileceğini
// belirlemek için kullanılır.
//
// # Display Key Kullanılan İlişkiler
//
// **HasMany:**
// - Birden fazla kayıt döndürür
// - Her kayıt için display key ile gösterim yapılır
// - Örnek: User -> Posts (post.title gösterilir)
//
// **BelongsToMany:**
// - Çok-çok ilişkide birden fazla kayıt
// - Pivot tablo üzerinden display key kullanılır
// - Örnek: User -> Roles (role.name gösterilir)
//
// **MorphMany:**
// - Polimorfik bir-çok ilişki
// - Display key ile kayıtlar listelenir
// - Örnek: Post -> Comments (comment.body gösterilir)
//
// # HasOne vs HasMany Karşılaştırması
//
// **HasOne (Bu Method):**
//
//	field := fields.HasOne("Profile", "profile", "profiles")
//	displayKey := field.GetDisplayKey() // "" (boş)
//	// Tek kayıt döndürür, display key gerekmez
//
// **HasMany:**
//
//	field := fields.HasMany("Posts", "posts", "posts").
//	    DisplayKey("title")
//	displayKey := field.GetDisplayKey() // "title"
//	// Birden fazla kayıt, her biri title ile gösterilir
//
// # Interface Implementation
//
// Bu method, Relationship interface'ini implement eder:
//
//	type Relationship interface {
//	    GetDisplayKey() string  // Bu method
//	    GetSearchableColumns() []string
//	    // ... diğer method'lar
//	}
//
// # Kullanım Örneği
//
//	field := fields.HasOne("Profile", "profile", "profiles")
//	displayKey := field.GetDisplayKey()
//	fmt.Println(displayKey) // Output: "" (boş string)
//
//	// Tip kontrolü
//	if displayKey == "" {
//	    fmt.Println("Display key not used for HasOne")
//	}
//
// # JSON Serialization
//
// API response'da bu alan genellikle dahil edilmez:
//
//	{
//	    "name": "Profile",
//	    "key": "profile",
//	    "type": "relationship",
//	    "relationship_type": "hasOne",
//	    "related_resource": "profiles"
//	    // display_key yok (boş olduğu için)
//	}
//
// # Alternatif Yaklaşımlar
//
// HasOne için display key yerine, ilişkili kaydın tamamı döndürülür:
//
//	// API Response
//	{
//	    "id": 1,
//	    "name": "John",
//	    "profile": {
//	        "id": 42,
//	        "full_name": "John Doe",
//	        "avatar": "avatar.jpg"
//	    }
//	}
//
// # Önemli Notlar
//
// - Her zaman boş string döndürür
// - HasOne için display key kavramı geçerli değildir
// - Çoklu kayıt ilişkilerinde (HasMany, BelongsToMany) kullanılır
// - Interface uyumluluğu için implement edilmiştir
//
// # Best Practices
//
// HasOne ilişkisinde display key kullanmayın:
//
//	// Yanlış (etkisiz)
//	field := fields.HasOne("Profile", "profile", "profiles").
//	    DisplayKey("full_name") // HasOne'da DisplayKey method'u yok
//
//	// Doğru
//	field := fields.HasOne("Profile", "profile", "profiles")
//	// Display key gerekmez, tüm profile objesi döndürülür
//
// # Performans
//
// - O(1) kompleksitesi (sabit string dönüşü)
// - Bellek allocation yok
// - Çok hızlı çalışır
//
// Döndürür:
//   - Boş string ("") - HasOne için kullanılmaz
func (h *HasOneField) GetDisplayKey() string {
	return ""
}

// GetSearchableColumns, aranabilir sütunları döndürür (HasOne için kullanılmaz).
//
// Bu method, HasOne ilişkisinde kullanılmaz ve boş slice döndürür.
// Searchable columns, genellikle HasMany ve BelongsToMany gibi çoklu kayıt
// ilişkilerinde hangi sütunlarda arama yapılacağını belirtmek için kullanılır.
//
// # Dönüş Değeri
//
// - Boş string slice ([]string{}) - HasOne için searchable columns kavramı geçerli değildir
//
// # Neden Kullanılmaz?
//
// HasOne ilişkisi tek bir kayıt döndürdüğü için, searchable columns'a ihtiyaç yoktur.
// Searchable columns, çoklu kayıtların listesinde arama yapmak için kullanılır.
//
// # Searchable Columns Kullanılan İlişkiler
//
// **HasMany:**
// - Birden fazla kayıt döndürür
// - Kayıtlar arasında arama yapılabilir
// - Örnek: User -> Posts (title, body sütunlarında ara)
//
// **BelongsToMany:**
// - Çok-çok ilişkide birden fazla kayıt
// - İlişkili kayıtlarda arama yapılabilir
// - Örnek: User -> Roles (name, description sütunlarında ara)
//
// **MorphMany:**
// - Polimorfik bir-çok ilişki
// - Searchable columns ile kayıtlar filtrelenir
// - Örnek: Post -> Comments (body, author sütunlarında ara)
//
// # HasOne vs HasMany Karşılaştırması
//
// **HasOne (Bu Method):**
//
//	field := fields.HasOne("Profile", "profile", "profiles")
//	columns := field.GetSearchableColumns() // [] (boş)
//	// Tek kayıt döndürür, arama gerekmez
//
// **HasMany:**
//
//	field := fields.HasMany("Posts", "posts", "posts").
//	    SearchableColumns("title", "body", "excerpt")
//	columns := field.GetSearchableColumns() // ["title", "body", "excerpt"]
//	// Birden fazla kayıt, arama yapılabilir
//
// # Interface Implementation
//
// Bu method, Relationship interface'ini implement eder:
//
//	type Relationship interface {
//	    GetDisplayKey() string
//	    GetSearchableColumns() []string  // Bu method
//	    // ... diğer method'lar
//	}
//
// # Kullanım Örneği
//
//	field := fields.HasOne("Profile", "profile", "profiles")
//	columns := field.GetSearchableColumns()
//	fmt.Println(len(columns)) // Output: 0 (boş slice)
//
//	// Tip kontrolü
//	if len(columns) == 0 {
//	    fmt.Println("Searchable columns not used for HasOne")
//	}
//
// # JSON Serialization
//
// API response'da bu alan genellikle dahil edilmez:
//
//	{
//	    "name": "Profile",
//	    "key": "profile",
//	    "type": "relationship",
//	    "relationship_type": "hasOne",
//	    "related_resource": "profiles"
//	    // searchable_columns yok (boş olduğu için)
//	}
//
// # Alternatif Yaklaşımlar
//
// HasOne için arama yerine, ilişkili kaydın tamamı döndürülür ve
// frontend'de filtreleme yapılabilir:
//
//	// API Response
//	{
//	    "id": 1,
//	    "name": "John",
//	    "profile": {
//	        "id": 42,
//	        "full_name": "John Doe",
//	        "bio": "Software Developer"
//	    }
//	}
//
//	// Frontend'de filtreleme
//	if (user.profile.full_name.includes(searchTerm)) {
//	    // Göster
//	}
//
// # Global Search
//
// HasOne ilişkisi global search'e dahil edilebilir:
//
//	field := fields.HasOne("Profile", "profile", "profiles").
//	    Searchable() // Global search'e dahil et
//
//	// Ancak bu, ilişkili kaydın kendisini değil,
//	// ilişkinin varlığını aranabilir yapar
//
// # Önemli Notlar
//
// - Her zaman boş slice döndürür
// - HasOne için searchable columns kavramı geçerli değildir
// - Çoklu kayıt ilişkilerinde (HasMany, BelongsToMany) kullanılır
// - Interface uyumluluğu için implement edilmiştir
// - Nil değil, boş slice döndürür ([]string{})
//
// # Best Practices
//
// HasOne ilişkisinde searchable columns kullanmayın:
//
//	// Yanlış (etkisiz)
//	field := fields.HasOne("Profile", "profile", "profiles").
//	    SearchableColumns("full_name", "bio") // HasOne'da bu method yok
//
//	// Doğru
//	field := fields.HasOne("Profile", "profile", "profiles")
//	// Searchable columns gerekmez, tek kayıt döndürülür
//
// # Performans
//
// - O(1) kompleksitesi (sabit slice dönüşü)
// - Minimal bellek allocation (boş slice)
// - Çok hızlı çalışır
//
// Döndürür:
//   - Boş string slice ([]string{}) - HasOne için kullanılmaz
func (h *HasOneField) GetSearchableColumns() []string {
	return []string{}
}

// GetQueryCallback, query callback fonksiyonunu döndürür.
//
// Bu method, ilişki sorgusunu özelleştirmek için tanımlanmış callback fonksiyonunu döndürür.
// Eğer callback tanımlanmamışsa, varsayılan olarak sorguyu değiştirmeyen bir fonksiyon döndürür.
//
// # Dönüş Değeri
//
// - Query callback fonksiyonu: func(interface{}) interface{}
// - Callback yoksa: Sorguyu olduğu gibi döndüren varsayılan fonksiyon
//
// # Çalışma Mantığı
//
// 1. QueryCallback field'ını kontrol eder
// 2. Eğer nil ise, varsayılan pass-through fonksiyon döndürür
// 3. Eğer tanımlıysa, kullanıcı tarafından ayarlanan callback'i döndürür
//
// # Kullanım Senaryoları
//
// - **Query Execution**: İlişki yüklenirken callback'i uygulama
// - **Dynamic Filtering**: Runtime'da sorgu özelleştirme
// - **Middleware**: Sorgu pipeline'ına ekleme
// - **Logging**: Sorgu loglaması için callback kullanma
// - **Testing**: Mock callback'ler ile test etme
//
// # Kullanım Örneği
//
//	// Callback tanımlı field
//	field := fields.HasOne("Profile", "profile", "profiles").
//	    Query(func(q interface{}) interface{} {
//	        return q.(*gorm.DB).Where("status = ?", "active")
//	    })
//
//	// Callback'i al ve kullan
//	callback := field.GetQueryCallback()
//	query := db.Model(&Profile{})
//	query = callback(query).(*gorm.DB)
//	query.Find(&profiles)
//
// # Varsayılan Callback
//
// Callback tanımlanmamışsa, varsayılan fonksiyon döndürülür:
//
//	func(q interface{}) interface{} {
//	    return q  // Sorguyu olduğu gibi döndür
//	}
//
// Bu, callback kontrolü yapmadan güvenle kullanılabilir:
//
//	// Her zaman çalışır (nil check gerekmez)
//	callback := field.GetQueryCallback()
//	query = callback(query)
//
// # Callback Tanımlama
//
// Query method'u ile callback tanımlanır:
//
//	field.Query(func(q interface{}) interface{} {
//	    db := q.(*gorm.DB)
//	    return db.Where("verified = ?", true).
//	        Order("created_at DESC")
//	})
//
// # Callback Uygulama
//
// **Eager Loading'de:**
//
//	func loadRelationship(field *HasOneField, ownerID uint) interface{} {
//	    query := db.Model(&Profile{}).
//	        Where("user_id = ?", ownerID)
//
//	    // Callback uygula
//	    callback := field.GetQueryCallback()
//	    query = callback(query).(*gorm.DB)
//
//	    var profile Profile
//	    query.First(&profile)
//	    return profile
//	}
//
// **Lazy Loading'de:**
//
//	func (h *HasOneField) ResolveRelationship(item interface{}) (interface{}, error) {
//	    ownerValue := extractFieldValue(item, h.OwnerKeyColumn)
//
//	    query := db.Model(relatedModel).
//	        Where(fmt.Sprintf("%s = ?", h.ForeignKeyColumn), ownerValue)
//
//	    // Callback uygula
//	    callback := h.GetQueryCallback()
//	    query = callback(query).(*gorm.DB)
//
//	    var result interface{}
//	    err := query.First(&result).Error
//	    return result, err
//	}
//
// # Callback Chaining
//
// Birden fazla callback zincirleme yapılabilir:
//
//	baseCallback := field.GetQueryCallback()
//	enhancedCallback := func(q interface{}) interface{} {
//	    // Önce base callback uygula
//	    q = baseCallback(q)
//	    // Sonra ek işlemler
//	    return q.(*gorm.DB).Limit(10)
//	}
//
// # Type Safety
//
// Callback içinde type assertion kullanılır:
//
//	callback := func(q interface{}) interface{} {
//	    db, ok := q.(*gorm.DB)
//	    if !ok {
//	        return q  // Type assertion başarısız, olduğu gibi döndür
//	    }
//	    return db.Where("active = ?", true)
//	}
//
// # Testing
//
// Test'lerde mock callback kullanılabilir:
//
//	// Test setup
//	mockCallback := func(q interface{}) interface{} {
//	    // Mock davranış
//	    return q
//	}
//
//	field := fields.HasOne("Profile", "profile", "profiles").
//	    Query(mockCallback)
//
//	// Test
//	callback := field.GetQueryCallback()
//	assert.NotNil(t, callback)
//
// # Önemli Notlar
//
// - Method her zaman nil olmayan bir fonksiyon döndürür
// - Varsayılan callback sorguyu değiştirmez (pass-through)
// - Callback nil kontrolü yapmanıza gerek yoktur
// - Type assertion başarısız olursa panic riski vardır
// - Callback içinde error handling yapılmalıdır
//
// # Best Practices
//
// **Güvenli Type Assertion:**
//
//	callback := func(q interface{}) interface{} {
//	    db, ok := q.(*gorm.DB)
//	    if !ok {
//	        log.Warn("Type assertion failed")
//	        return q
//	    }
//	    return db.Where("active = ?", true)
//	}
//
// **Error Handling:**
//
//	callback := func(q interface{}) interface{} {
//	    defer func() {
//	        if r := recover(); r != nil {
//	            log.Error("Callback panic:", r)
//	        }
//	    }()
//	    return q.(*gorm.DB).Where("active = ?", true)
//	}
//
// **Logging:**
//
//	callback := func(q interface{}) interface{} {
//	    log.Debug("Applying query callback")
//	    db := q.(*gorm.DB)
//	    return db.Where("active = ?", true)
//	}
//
// # Performans
//
// - O(1) kompleksitesi (field access veya varsayılan fonksiyon dönüşü)
// - Minimal bellek allocation
// - Callback execution overhead callback içeriğine bağlıdır
//
// # Uyarılar
//
// - Callback içinde panic oluşursa uygulama çökebilir
// - Type assertion başarısız olursa runtime error oluşur
// - Callback'te yapılan değişiklikler tüm sorguları etkiler
// - Thread-safety callback implementasyonuna bağlıdır
//
// Döndürür:
//   - Query callback fonksiyonu (her zaman nil olmayan)
func (h *HasOneField) GetQueryCallback() func(interface{}) interface{} {
	if h.QueryCallback == nil {
		return func(q interface{}) interface{} { return q }
	}
	return h.QueryCallback
}

// GetLoadingStrategy, yükleme stratejisini döndürür.
//
// Bu method, HasOne ilişkisinin yükleme stratejisini (eager veya lazy loading) döndürür.
// Eğer strateji ayarlanmamışsa, varsayılan olarak EAGER_LOADING döndürür.
//
// # Dönüş Değeri
//
// - LoadingStrategy: EAGER_LOADING veya LAZY_LOADING
// - Varsayılan: EAGER_LOADING (N+1 sorgu problemini önlemek için)
//
// # Loading Strategy Tipleri
//
// **EAGER_LOADING:**
// - İlişkili kayıtlar ana kayıtlarla birlikte yüklenir
// - N+1 sorgu problemini önler
// - Liste görünümlerinde önerilir
// - Performanslı toplu yükleme
//
// **LAZY_LOADING:**
// - İlişkili kayıtlar sadece gerektiğinde yüklenir
// - İsteğe bağlı yükleme
// - Tek kayıt görünümlerinde kullanılabilir
// - N+1 problem riski vardır
//
// # Kullanım Senaryoları
//
// - **Query Optimization**: Sorgu stratejisini belirleme
// - **Preload Decision**: Preload yapılıp yapılmayacağına karar verme
// - **Performance Tuning**: Performans optimizasyonu için strateji seçimi
// - **Dynamic Loading**: Runtime'da yükleme stratejisini kontrol etme
// - **Debugging**: Yükleme stratejisini loglama ve debug etme
//
// # Kullanım Örneği
//
//	// Eager loading ile field
//	field := fields.HasOne("Profile", "profile", "profiles").
//	    WithEagerLoad()
//	strategy := field.GetLoadingStrategy()
//	fmt.Println(strategy) // Output: EAGER_LOADING
//
//	// Lazy loading ile field
//	field := fields.HasOne("Profile", "profile", "profiles").
//	    WithLazyLoad()
//	strategy := field.GetLoadingStrategy()
//	fmt.Println(strategy) // Output: LAZY_LOADING
//
//	// Varsayılan (strateji ayarlanmamış)
//	field := fields.HasOne("Profile", "profile", "profiles")
//	strategy := field.GetLoadingStrategy()
//	fmt.Println(strategy) // Output: EAGER_LOADING (varsayılan)
//
// # Strateji Bazlı Query Execution
//
//	func loadRelationships(fields []*HasOneField, items []interface{}) {
//	    for _, field := range fields {
//	        strategy := field.GetLoadingStrategy()
//
//	        switch strategy {
//	        case EAGER_LOADING:
//	            // Tüm kayıtları tek sorguda yükle
//	            preloadRelationships(field, items)
//
//	        case LAZY_LOADING:
//	            // Her kayıt için ayrı ayrı yükle (gerektiğinde)
//	            // Hiçbir şey yapma, isteğe bağlı yükleme
//	        }
//	    }
//	}
//
// # Eager Loading Implementation
//
//	if field.GetLoadingStrategy() == EAGER_LOADING {
//	    // GORM Preload kullan
//	    db.Preload(field.Key).Find(&users)
//
//	    // Veya manuel batch loading
//	    ownerIDs := extractOwnerIDs(users)
//	    profiles := loadProfilesByOwnerIDs(ownerIDs)
//	    mapProfilesToUsers(users, profiles)
//	}
//
// # Lazy Loading Implementation
//
//	if field.GetLoadingStrategy() == LAZY_LOADING {
//	    // İlişki sadece erişildiğinde yüklenir
//	    // Hiçbir şey yapma, manuel yükleme gerekir
//	    for _, user := range users {
//	        if needsProfile(user) {
//	            profile, _ := field.ResolveRelationship(user)
//	            user.Profile = profile
//	        }
//	    }
//	}
//
// # Varsayılan Davranış
//
// Strateji ayarlanmamışsa (boş string), varsayılan olarak EAGER_LOADING döndürülür:
//
//	field := fields.HasOne("Profile", "profile", "profiles")
//	// LoadingStrategy field'ı boş string
//
//	strategy := field.GetLoadingStrategy()
//	// EAGER_LOADING döndürür (varsayılan)
//
// Bu, N+1 sorgu problemini önlemek için güvenli bir varsayılandır.
//
// # Performans Karşılaştırması
//
// **100 Kullanıcı + Profilleri:**
//
// **Eager Loading:**
// - Sorgu sayısı: 2 (users + profiles)
// - Süre: ~50ms
// - Bellek: Yüksek (tüm veriler)
//
// **Lazy Loading:**
// - Sorgu sayısı: 101 (users + 100 profile sorgusu)
// - Süre: ~500ms (10x yavaş)
// - Bellek: Düşük (sadece gerekli veriler)
//
// # Strateji Seçim Rehberi
//
// **Eager Loading Kullan:**
// - Liste görünümlerinde (index, table)
// - İlişkili veri kesinlikle gerekli
// - Birden fazla kayıt işlenirken
// - API response'larında
// - Export/raporlama işlemlerinde
//
// **Lazy Loading Kullan:**
// - Tek kayıt görünümlerinde (show, detail)
// - İlişkili veri nadiren gerekli
// - Bellek kısıtlı ortamlarda
// - Çok büyük ilişkili veri setlerinde
// - Koşullu veri yükleme senaryolarında
//
// # Dynamic Strategy Selection
//
// Runtime'da strateji değiştirilebilir:
//
//	func getLoadingStrategy(context Context) LoadingStrategy {
//	    if context.IsBulkOperation() {
//	        return EAGER_LOADING
//	    }
//	    if context.IsSingleRecord() {
//	        return LAZY_LOADING
//	    }
//	    return EAGER_LOADING // Varsayılan
//	}
//
// # Monitoring ve Logging
//
//	strategy := field.GetLoadingStrategy()
//	log.Info("Loading relationship",
//	    "field", field.Name,
//	    "strategy", strategy,
//	    "record_count", len(items))
//
//	if strategy == LAZY_LOADING {
//	    log.Warn("Lazy loading may cause N+1 queries",
//	        "field", field.Name)
//	}
//
// # Testing
//
//	// Test eager loading
//	field := fields.HasOne("Profile", "profile", "profiles").
//	    WithEagerLoad()
//	assert.Equal(t, EAGER_LOADING, field.GetLoadingStrategy())
//
//	// Test lazy loading
//	field = fields.HasOne("Profile", "profile", "profiles").
//	    WithLazyLoad()
//	assert.Equal(t, LAZY_LOADING, field.GetLoadingStrategy())
//
//	// Test default
//	field = fields.HasOne("Profile", "profile", "profiles")
//	assert.Equal(t, EAGER_LOADING, field.GetLoadingStrategy())
//
// # Önemli Notlar
//
// - Varsayılan strateji EAGER_LOADING'dir (güvenli seçim)
// - Boş string kontrolü yapılır ve varsayılan döndürülür
// - Strateji WithEagerLoad() veya WithLazyLoad() ile ayarlanır
// - Liste görünümlerinde mutlaka eager loading kullanın
// - N+1 problemini önlemek için varsayılan eager'dır
//
// # Best Practices
//
// **Liste Görünümlerinde:**
//
//	field := fields.HasOne("Profile", "profile", "profiles").
//	    WithEagerLoad() // Açıkça belirt
//
// **Tek Kayıt Görünümlerinde:**
//
//	field := fields.HasOne("Profile", "profile", "profiles").
//	    WithLazyLoad() // Gerektiğinde yükle
//
// **Varsayılana Güvenme:**
//
//	// Kötü: Strateji belirsiz
//	field := fields.HasOne("Profile", "profile", "profiles")
//
//	// İyi: Strateji açıkça belirtilmiş
//	field := fields.HasOne("Profile", "profile", "profiles").
//	    WithEagerLoad()
//
// # Performans İpuçları
//
// - Production'da eager loading tercih edin
// - Lazy loading sadece gerektiğinde kullanın
// - Sorgu sayısını monitor edin
// - N+1 problemini tespit edin ve düzeltin
// - Cache mekanizması ile optimize edin
//
// # Uyarılar
//
// - Lazy loading N+1 sorgu problemine yol açar
// - Liste görünümlerinde lazy loading kullanmayın
// - Varsayılan eager olsa da açıkça belirtmek daha iyidir
// - Strateji değişikliği tüm sorguları etkiler
//
// Döndürür:
//   - LoadingStrategy: EAGER_LOADING veya LAZY_LOADING (varsayılan: EAGER_LOADING)
func (h *HasOneField) GetLoadingStrategy() LoadingStrategy {
	if h.LoadingStrategy == "" {
		return EAGER_LOADING
	}
	return h.LoadingStrategy
}

// Searchable, alanı aranabilir olarak işaretler (Element interface'ini implement eder).
//
// Bu method, HasOne ilişkisini global search (genel arama) özelliğine dahil eder.
// Global search, kullanıcıların tüm kayıtlar arasında arama yapmasını sağlar.
//
// # Çalışma Mantığı
//
// Method, GlobalSearch field'ını true olarak ayarlar ve Element interface'ini döndürür.
// Bu, field'ın global search sorgularına dahil edilmesini sağlar.
//
// # Dönüş Değeri
//
// - Element interface'i (method chaining için)
//
// # Kullanım Senaryoları
//
// - **Global Search**: Tüm kayıtlarda arama yapılabilir hale getirme
// - **Quick Search**: Hızlı arama özelliğine dahil etme
// - **Search Bar**: Arama çubuğunda ilişki araması
// - **Filter Integration**: Filtreleme sistemine entegrasyon
// - **Admin Panel**: Yönetim panelinde arama özelliği
//
// # Kullanım Örneği
//
//	// Searchable HasOne field
//	field := fields.HasOne("Profile", "profile", "profiles").
//	    Searchable()
//
//	// Method chaining ile diğer özellikler
//	field := fields.HasOne("Profile", "profile", "profiles").
//	    Searchable().
//	    WithEagerLoad().
//	    ForeignKey("user_id")
//
// # Global Search Davranışı
//
// HasOne ilişkisi searchable olarak işaretlendiğinde:
//
// **İlişki Varlığı Aranır:**
// - İlişkili kayıt var mı yok mu kontrol edilir
// - İlişkili kaydın ID'si aranabilir
//
// **İlişkili Kayıt İçeriği Aranmaz:**
// - HasOne için ilişkili kaydın içeriği (örn. profile.full_name) aranmaz
// - Sadece ilişkinin varlığı aranabilir
//
// # HasOne vs HasMany Searchable Karşılaştırması
//
// **HasOne (Bu Method):**
//
//	field := fields.HasOne("Profile", "profile", "profiles").
//	    Searchable()
//	// İlişki varlığı aranır (profile var mı?)
//
// **HasMany:**
//
//	field := fields.HasMany("Posts", "posts", "posts").
//	    Searchable().
//	    SearchableColumns("title", "body")
//	// İlişkili kayıtların içeriği aranır (post title, body)
//
// # Search Query Implementation (Örnek)
//
//	func buildSearchQuery(fields []Field, searchTerm string) *gorm.DB {
//	    query := db.Model(&User{})
//
//	    for _, field := range fields {
//	        if !field.IsSearchable() {
//	            continue
//	        }
//
//	        if hasOneField, ok := field.(*HasOneField); ok {
//	            // HasOne için ilişki varlığı kontrolü
//	            query = query.Or(
//	                db.Joins(hasOneField.Key).
//	                    Where(fmt.Sprintf("%s.id IS NOT NULL", hasOneField.Key)),
//	            )
//	        }
//	    }
//
//	    return query
//	}
//
// # Frontend Integration
//
// Frontend'de searchable field'lar otomatik olarak arama formuna eklenir:
//
//	// API Response
//	{
//	    "fields": [
//	        {
//	            "name": "Profile",
//	            "key": "profile",
//	            "type": "relationship",
//	            "searchable": true  // Bu method ile ayarlanır
//	        }
//	    ]
//	}
//
//	// Frontend Search Form
//	<SearchBar>
//	    <input placeholder="Search users with profile..." />
//	</SearchBar>
//
// # Element Interface Implementation
//
// Bu method, Element interface'ini implement eder:
//
//	type Element interface {
//	    Searchable() Element
//	    IsSearchable() bool
//	    // ... diğer method'lar
//	}
//
// # IsSearchable Kontrolü
//
// Searchable olup olmadığını kontrol etmek için:
//
//	if field.GlobalSearch {
//	    // Field searchable
//	    includeInSearch(field)
//	}
//
// # Kullanım Örnekleri
//
// **Basit Kullanım:**
//
//	field := fields.HasOne("Profile", "profile", "profiles").
//	    Searchable()
//
// **Diğer Özelliklerle Birlikte:**
//
//	field := fields.HasOne("Profile", "profile", "profiles").
//	    Searchable().
//	    WithEagerLoad().
//	    ForeignKey("user_id").
//	    Required()
//
// **Koşullu Searchable:**
//
//	field := fields.HasOne("Profile", "profile", "profiles")
//	if config.EnableProfileSearch {
//	    field.Searchable()
//	}
//
// # Search Query Örnekleri
//
// **İlişki Varlığı Araması:**
//
//	// Profili olan kullanıcıları bul
//	db.Joins("Profile").
//	    Where("profiles.id IS NOT NULL").
//	    Find(&users)
//
// **İlişki Yokluğu Araması:**
//
//	// Profili olmayan kullanıcıları bul
//	db.Joins("LEFT JOIN profiles ON profiles.user_id = users.id").
//	    Where("profiles.id IS NULL").
//	    Find(&users)
//
// # Performans Considerations
//
// **Index Kullanımı:**
//
//	CREATE INDEX idx_profiles_user_id ON profiles(user_id);
//
// **Join Optimizasyonu:**
//
//	// Efficient join
//	db.Joins("Profile").Find(&users)
//
//	// Inefficient (N+1)
//	for _, user := range users {
//	    db.Model(&user).Association("Profile").Find(&user.Profile)
//	}
//
// # Önemli Notlar
//
// - GlobalSearch field'ı true olarak ayarlanır
// - Element interface'i döndürülür (method chaining için)
// - HasOne için sadece ilişki varlığı aranır, içerik aranmaz
// - Frontend'de otomatik olarak arama formuna eklenir
// - Search query'lerde JOIN kullanılır
//
// # Best Practices
//
// **Searchable Kullanımı:**
//
//	// İyi: İlişki varlığı önemli
//	field := fields.HasOne("Profile", "profile", "profiles").
//	    Searchable() // Profili olan/olmayan kullanıcıları bul
//
//	// Kötü: İlişki içeriği aranmak isteniyorsa
//	// HasOne searchable ile içerik aranamaz
//	// İlişkili resource'da ayrı field'lar tanımlayın
//
// **Index Tanımlama:**
//
//	// Migration'da index ekle
//	CREATE INDEX idx_profiles_user_id ON profiles(user_id);
//
// **Query Optimization:**
//
//	// Eager loading ile birlikte kullan
//	field := fields.HasOne("Profile", "profile", "profiles").
//	    Searchable().
//	    WithEagerLoad()
//
// # Alternatif Yaklaşımlar
//
// İlişkili kayıt içeriğinde arama yapmak için:
//
//	// Profile resource'da searchable field'lar tanımla
//	profileResource := resource.New("profiles").
//	    Fields(
//	        fields.Text("Full Name", "full_name").Searchable(),
//	        fields.Text("Bio", "bio").Searchable(),
//	    )
//
//	// User resource'da HasOne tanımla
//	userResource := resource.New("users").
//	    Fields(
//	        fields.HasOne("Profile", "profile", profileResource).
//	            WithEagerLoad(),
//	    )
//
// # Uyarılar
//
// - HasOne searchable ile sadece ilişki varlığı aranır
// - İlişkili kayıt içeriği aranmaz (full_name, bio vb.)
// - Search query'lerde JOIN kullanılır (performans etkisi)
// - Index'lerin doğru tanımlandığından emin olun
// - Büyük veri setlerinde performans test edin
//
// Döndürür:
//   - Element interface'i (method chaining için)
func (h *HasOneField) Searchable() Element {
	h.GlobalSearch = true
	return h
}

// IsRequired, alanın zorunlu olup olmadığını döndürür.
//
// Bu method, HasOne ilişkisinin zorunlu (required) olarak işaretlenip işaretlenmediğini kontrol eder.
// Zorunlu alan, form validasyonunda ve veri kaydında mutlaka doldurulması gereken alandır.
//
// # Dönüş Değeri
//
// - true: Alan zorunludur, mutlaka doldurulmalıdır
// - false: Alan opsiyoneldir, boş bırakılabilir
//
// # Kullanım Senaryoları
//
// - **Form Validation**: Form gönderilmeden önce zorunlu alan kontrolü
// - **API Validation**: Request payload validasyonu
// - **Database Constraints**: NOT NULL constraint kontrolü
// - **UI Rendering**: Zorunlu alanları görsel olarak işaretleme (*)
// - **Error Messages**: Zorunlu alan hata mesajları
//
// # Kullanım Örneği
//
//	// Zorunlu HasOne field
//	field := fields.HasOne("Profile", "profile", "profiles").
//	    Required()
//
//	// Zorunluluk kontrolü
//	if field.IsRequired() {
//	    fmt.Println("Profile is required")
//	}
//
//	// Form validation
//	if field.IsRequired() && profileID == 0 {
//	    return errors.New("Profile is required")
//	}
//
// # Required Method ile Kullanım
//
// Alanı zorunlu yapmak için Required() method'u kullanılır:
//
//	field := fields.HasOne("Profile", "profile", "profiles").
//	    Required() // IsRequired field'ını true yapar
//
//	fmt.Println(field.IsRequired()) // Output: true
//
// # Validation Implementation
//
//	func validateField(field *HasOneField, value interface{}) error {
//	    if field.IsRequired() {
//	        if value == nil {
//	            return fmt.Errorf("%s is required", field.Name)
//	        }
//
//	        // ID kontrolü
//	        if id, ok := value.(uint); ok && id == 0 {
//	            return fmt.Errorf("%s is required", field.Name)
//	        }
//	    }
//	    return nil
//	}
//
// # Frontend Integration
//
// Frontend'de zorunlu alanlar otomatik olarak işaretlenir:
//
//	// API Response
//	{
//	    "name": "Profile",
//	    "key": "profile",
//	    "type": "relationship",
//	    "required": true  // Bu method'un dönüş değeri
//	}
//
//	// Frontend Rendering
//	<label>
//	    Profile {field.required && <span className="text-red-500">*</span>}
//	</label>
//	<select name="profile_id" required={field.required}>
//	    <option value="">Select Profile</option>
//	</select>
//
// # Form Validation Örneği
//
//	// Backend validation
//	func validateUserForm(data map[string]interface{}, fields []*HasOneField) error {
//	    for _, field := range fields {
//	        if field.IsRequired() {
//	            value, exists := data[field.Key]
//	            if !exists || value == nil {
//	                return fmt.Errorf("%s is required", field.Name)
//	            }
//	        }
//	    }
//	    return nil
//	}
//
//	// Kullanım
//	err := validateUserForm(formData, userFields)
//	if err != nil {
//	    return err
//	}
//
// # Database Constraints
//
// Zorunlu alan, veritabanında NOT NULL constraint ile eşleşir:
//
//	// Migration
//	if field.IsRequired() {
//	    // Foreign key NOT NULL olmalı
//	    ALTER TABLE profiles
//	    ADD CONSTRAINT fk_profiles_user_id
//	    FOREIGN KEY (user_id) REFERENCES users(id)
//	    NOT NULL;
//	}
//
// # Error Messages
//
// Zorunlu alan için hata mesajları:
//
//	if field.IsRequired() && value == nil {
//	    return map[string]string{
//	        field.Key: fmt.Sprintf("%s is required", field.Name),
//	    }
//	}
//
//	// Örnek çıktı
//	{
//	    "profile": "Profile is required"
//	}
//
// # UI Rendering
//
// **Form Label:**
//
//	<label htmlFor="profile_id">
//	    {field.Name}
//	    {field.IsRequired() && <span className="text-red-500">*</span>}
//	</label>
//
// **Input Field:**
//
//	<select
//	    id="profile_id"
//	    name="profile_id"
//	    required={field.IsRequired()}
//	    className={field.IsRequired() ? "required" : ""}
//	>
//	    <option value="">Select {field.Name}</option>
//	</select>
//
// **Validation Message:**
//
//	{errors.profile && (
//	    <span className="error">
//	        {field.Name} is required
//	    </span>
//	)}
//
// # API Validation
//
//	// Request handler
//	func createUser(c *gin.Context) {
//	    var data map[string]interface{}
//	    c.BindJSON(&data)
//
//	    // Validate required fields
//	    for _, field := range userResource.Fields() {
//	        if hasOneField, ok := field.(*HasOneField); ok {
//	            if hasOneField.IsRequired() {
//	                value, exists := data[hasOneField.Key]
//	                if !exists || value == nil {
//	                    c.JSON(400, gin.H{
//	                        "error": fmt.Sprintf("%s is required", hasOneField.Name),
//	                    })
//	                    return
//	                }
//	            }
//	        }
//	    }
//
//	    // Create user
//	    // ...
//	}
//
// # Testing
//
//	// Test required field
//	func TestHasOneRequired(t *testing.T) {
//	    field := fields.HasOne("Profile", "profile", "profiles").
//	        Required()
//
//	    assert.True(t, field.IsRequired())
//	}
//
//	// Test optional field
//	func TestHasOneOptional(t *testing.T) {
//	    field := fields.HasOne("Profile", "profile", "profiles")
//
//	    assert.False(t, field.IsRequired())
//	}
//
// # Önemli Notlar
//
// - Schema.IsRequired field'ından değer okunur
// - Varsayılan olarak false'dur (opsiyonel)
// - Required() method'u ile true yapılır
// - Frontend ve backend validasyonunda kullanılır
// - Database constraint'leri ile senkronize olmalıdır
//
// # Best Practices
//
// **Zorunlu İlişkiler:**
//
//	// İyi: Açıkça required olarak işaretle
//	field := fields.HasOne("Profile", "profile", "profiles").
//	    Required()
//
//	// Kötü: Varsayılana güvenme
//	field := fields.HasOne("Profile", "profile", "profiles")
//	// IsRequired() false döner
//
// **Validation:**
//
//	// İyi: Hem frontend hem backend validation
//	// Frontend: HTML5 required attribute
//	// Backend: IsRequired() kontrolü
//
//	// Kötü: Sadece frontend validation
//	// Güvenlik riski, bypass edilebilir
//
// **Error Messages:**
//
//	// İyi: Kullanıcı dostu mesaj
//	if field.IsRequired() && value == nil {
//	    return fmt.Errorf("%s is required", field.Name)
//	}
//
//	// Kötü: Teknik mesaj
//	if field.IsRequired() && value == nil {
//	    return errors.New("validation failed: nil value")
//	}
//
// # Database Migration
//
//	// Required field için migration
//	if field.IsRequired() {
//	    sql := fmt.Sprintf(`
//	        ALTER TABLE %s
//	        MODIFY COLUMN %s INT NOT NULL
//	    `, tableName, field.ForeignKeyColumn)
//	    db.Exec(sql)
//	}
//
// # Conditional Required
//
// Bazı durumlarda koşullu zorunluluk gerekebilir:
//
//	// Örnek: Profil sadece verified kullanıcılar için zorunlu
//	func validateProfile(user User, profileID uint) error {
//	    field := fields.HasOne("Profile", "profile", "profiles")
//
//	    if user.Verified && field.IsRequired() {
//	        if profileID == 0 {
//	            return errors.New("Profile is required for verified users")
//	        }
//	    }
//	    return nil
//	}
//
// # Performans
//
// - O(1) kompleksitesi (field access)
// - Bellek allocation yok
// - Çok hızlı çalışır
//
// # Uyarılar
//
// - Required field'lar mutlaka doldurulmalıdır
// - Frontend validation bypass edilebilir, backend validation şarttır
// - Database constraint'leri ile senkronize olmalıdır
// - Migration'larda NOT NULL constraint eklenmelidir
// - Mevcut verilerde NULL değer varsa migration başarısız olur
//
// Döndürür:
//   - true: Alan zorunludur
//   - false: Alan opsiyoneldir
func (h *HasOneField) IsRequired() bool {
	return h.Schema.IsRequired
}

// GetTypes, tip eşlemelerini döndürür (HasOne için kullanılmaz).
//
// Bu method, HasOne ilişkisinde kullanılmaz ve boş map döndürür.
// Type mappings, genellikle form field'larında input tiplerini ve
// validasyon kurallarını tanımlamak için kullanılır.
//
// # Dönüş Değeri
//
// - Boş map (map[string]string{}) - HasOne için type mappings kavramı geçerli değildir
//
// # Neden Kullanılmaz?
//
// HasOne ilişkisi bir relationship field'dır ve primitive tip mapping'e ihtiyaç duymaz.
// Type mappings, genellikle Text, Number, Email gibi primitive field'larda
// input tipini ve validasyon kurallarını belirtmek için kullanılır.
//
// # Type Mappings Kullanılan Field'lar
//
// **Text Field:**
//
//	field := fields.Text("Name", "name")
//	types := field.GetTypes()
//	// {"input": "text", "validation": "string"}
//
// **Number Field:**
//
//	field := fields.Number("Age", "age")
//	types := field.GetTypes()
//	// {"input": "number", "validation": "integer"}
//
// **Email Field:**
//
//	field := fields.Email("Email", "email")
//	types := field.GetTypes()
//	// {"input": "email", "validation": "email"}
//
// **Date Field:**
//
//	field := fields.Date("Birth Date", "birth_date")
//	types := field.GetTypes()
//	// {"input": "date", "validation": "date"}
//
// # HasOne vs Primitive Fields Karşılaştırması
//
// **HasOne (Bu Method):**
//
//	field := fields.HasOne("Profile", "profile", "profiles")
//	types := field.GetTypes() // {} (boş map)
//	// İlişki field'ı, tip mapping gerekmez
//
// **Text Field:**
//
//	field := fields.Text("Name", "name")
//	types := field.GetTypes() // {"input": "text", "validation": "string"}
//	// Primitive field, tip mapping gerekir
//
// # Interface Implementation
//
// Bu method, Field interface'ini implement eder:
//
//	type Field interface {
//	    GetTypes() map[string]string
//	    // ... diğer method'lar
//	}
//
// # Kullanım Örneği
//
//	field := fields.HasOne("Profile", "profile", "profiles")
//	types := field.GetTypes()
//	fmt.Println(len(types)) // Output: 0 (boş map)
//
//	// Tip kontrolü
//	if len(types) == 0 {
//	    fmt.Println("Type mappings not used for HasOne")
//	}
//
// # JSON Serialization
//
// API response'da bu alan genellikle dahil edilmez:
//
//	{
//	    "name": "Profile",
//	    "key": "profile",
//	    "type": "relationship",
//	    "relationship_type": "hasOne",
//	    "related_resource": "profiles"
//	    // types yok (boş olduğu için)
//	}
//
// # Type Mappings Kullanım Alanları
//
// Type mappings, aşağıdaki alanlarda kullanılır:
//
// **Frontend Rendering:**
// - HTML input type belirleme
// - Form validation kuralları
// - Input maskeleme (telefon, kredi kartı vb.)
//
// **Backend Validation:**
// - Veri tipi kontrolü
// - Format validasyonu
// - Type casting
//
// **Database Schema:**
// - Sütun tipi belirleme
// - Migration oluşturma
// - Index stratejisi
//
// # HasOne İçin Alternatif
//
// HasOne ilişkisinde tip bilgisi farklı şekilde saklanır:
//
//	// İlişki tipi
//	field.GetRelationshipType() // "hasOne"
//
//	// Field tipi
//	field.Type // TYPE_RELATIONSHIP
//
//	// View component
//	field.View // "has-one-field"
//
// # Frontend Component Selection
//
// Frontend'de component seçimi View field'ına göre yapılır:
//
//	// HasOne için
//	switch field.View {
//	case "has-one-field":
//	    return <HasOneSelect field={field} />
//	case "text-field":
//	    return <TextInput field={field} />
//	case "number-field":
//	    return <NumberInput field={field} />
//	}
//
// # Type Mappings Örneği (Diğer Field'lar)
//
// **Text Field:**
//
//	{
//	    "input": "text",
//	    "validation": "string",
//	    "html_type": "text",
//	    "db_type": "VARCHAR"
//	}
//
// **Number Field:**
//
//	{
//	    "input": "number",
//	    "validation": "integer",
//	    "html_type": "number",
//	    "db_type": "INT"
//	}
//
// **Email Field:**
//
//	{
//	    "input": "email",
//	    "validation": "email",
//	    "html_type": "email",
//	    "db_type": "VARCHAR"
//	}
//
// # Önemli Notlar
//
// - Her zaman boş map döndürür
// - HasOne için type mappings kavramı geçerli değildir
// - Primitive field'larda (Text, Number, Email) kullanılır
// - Interface uyumluluğu için implement edilmiştir
// - Nil değil, boş map döndürür (map[string]string{})
//
// # Best Practices
//
// HasOne ilişkisinde type mappings kullanmayın:
//
//	// Yanlış (etkisiz)
//	types := field.GetTypes()
//	if types["input"] == "text" {
//	    // HasOne için bu kontrol anlamsız
//	}
//
//	// Doğru
//	if field.GetRelationshipType() == "hasOne" {
//	    // HasOne'a özel işlemler
//	}
//
// # Field Type Kontrolü
//
// Field tipini kontrol etmek için doğru method'ları kullanın:
//
//	// HasOne için
//	if field.Type == TYPE_RELATIONSHIP {
//	    relType := field.GetRelationshipType() // "hasOne"
//	}
//
//	// Primitive field için
//	if field.Type == TYPE_TEXT {
//	    types := field.GetTypes() // {"input": "text", ...}
//	}
//
// # Migration ve Schema
//
// HasOne için database schema bilgisi farklı şekilde alınır:
//
//	// Foreign key bilgisi
//	foreignKey := field.ForeignKeyColumn // "user_id"
//	ownerKey := field.OwnerKeyColumn     // "id"
//
//	// Migration
//	ALTER TABLE profiles
//	ADD CONSTRAINT fk_profiles_user_id
//	FOREIGN KEY (user_id) REFERENCES users(id);
//
// # Performans
//
// - O(1) kompleksitesi (boş map dönüşü)
// - Minimal bellek allocation (boş map)
// - Çok hızlı çalışır
//
// # Interface Uyumluluğu
//
// Bu method, Field interface'ini implement etmek için gereklidir:
//
//	type Field interface {
//	    GetTypes() map[string]string  // Tüm field'lar implement etmeli
//	    GetName() string
//	    GetKey() string
//	    // ... diğer method'lar
//	}
//
// HasOne bir Field olduğu için bu method'u implement etmelidir,
// ancak relationship field'lar için type mappings kullanılmaz.
//
// # Debugging
//
// Field tipini debug ederken:
//
//	fmt.Printf("Field: %s\n", field.Name)
//	fmt.Printf("Type: %s\n", field.Type)
//	fmt.Printf("View: %s\n", field.View)
//
//	if field.Type == TYPE_RELATIONSHIP {
//	    fmt.Printf("Relationship Type: %s\n", field.GetRelationshipType())
//	    fmt.Printf("Related Resource: %s\n", field.GetRelatedResource())
//	} else {
//	    fmt.Printf("Types: %+v\n", field.GetTypes())
//	}
//
// # Uyarılar
//
// - HasOne için GetTypes() kullanmayın (her zaman boş)
// - Tip bilgisi için GetRelationshipType() kullanın
// - Frontend component seçimi için View field'ını kullanın
// - Database schema için ForeignKeyColumn ve OwnerKeyColumn kullanın
//
// Döndürür:
//   - Boş map (map[string]string{}) - HasOne için kullanılmaz
func (h *HasOneField) GetTypes() map[string]string {
	return make(map[string]string)
}

// WithHoverCard, hover card konfigürasyonunu ayarlar.
//
// Bu metod, index ve detail sayfalarında ilişkili kaydın hover card ile
// nasıl görüntüleneceğini belirler.
//
// # Parametreler
//
// - **config**: Hover card konfigürasyonu
//
// # Kullanım Örneği (Deprecated - Yeni API kullanın)
//
//	field := fields.HasOne("Profile", "profile", "profiles").
//	    WithHoverCard(*fields.NewHoverCardConfig())
//
// # Yeni API (Önerilen)
//
//	field := fields.HasOne("Profile", "profile", "profiles").
//	    HoverCard(&ProfileHoverCard{}).
//	    ResolveHoverCard(func(ctx context.Context, record interface{}, relatedID interface{}, field fields.Field) (interface{}, error) {
//	        // Custom logic
//	        return &ProfileHoverCard{...}, nil
//	    })
//
// Döndürür:
//   - HasOneField pointer'ı (method chaining için)
func (h *HasOneField) WithHoverCard(config HoverCardConfig) *HasOneField {
	h.hoverCardConfig = &config
	h.WithProps("hover_card", config)
	return h
}

// HoverCard, hover card struct'ını ayarlar ve hover card'ı etkinleştirir.
//
// Bu metod, hover card için kullanılacak struct'ı belirler ve
// hover card özelliğini aktif eder.
//
// # Parametreler
//
// - **hoverStruct**: Hover card verisi için kullanılacak struct (örn. &ProfileHoverCard{})
//
// # Kullanım Örneği
//
//	type ProfileHoverCard struct {
//	    Avatar string `json:"avatar"`
//	    Bio    string `json:"bio"`
//	    Location string `json:"location"`
//	}
//
//	field := fields.HasOne("Profile", "profile", "profiles").
//	    ForeignKey("user_id").
//	    HoverCard(&ProfileHoverCard{})
//
// Döndürür:
//   - HasOneField pointer'ı (method chaining için)
func (h *HasOneField) HoverCard(hoverStruct interface{}) *HasOneField {
	if h.hoverCardConfig == nil {
		h.hoverCardConfig = NewHoverCardConfig()
	}
	h.hoverCardConfig.SetStruct(hoverStruct)
	h.WithProps("hover_card_enabled", true)
	return h
}

// ResolveHoverCard, hover card verilerini çözmek için callback fonksiyonunu ayarlar.
//
// Bu metod, hover card açıldığında çağrılacak resolver fonksiyonunu belirler.
// Resolver, ilişkili kaydın hover card verilerini döndürür.
//
// # Parametreler
//
// - **resolver**: Hover card resolver callback fonksiyonu
//
// # Kullanım Örneği
//
//	field := fields.HasOne("Profile", "profile", "profiles").
//	    ForeignKey("user_id").
//	    HoverCard(&ProfileHoverCard{}).
//	    ResolveHoverCard(func(ctx context.Context, record interface{}, relatedID interface{}, field fields.Field) (interface{}, error) {
//	        // İlişkili kaydı veritabanından al
//	        profile := &Profile{}
//	        if err := db.First(profile, relatedID).Error; err != nil {
//	            return nil, err
//	        }
//
//	        // Hover card verisini döndür
//	        return &ProfileHoverCard{
//	            Avatar: profile.Avatar,
//	            Bio: profile.Bio,
//	            Location: profile.Location,
//	        }, nil
//	    })
//
// # API Endpoint
//
// Frontend, hover card açıldığında şu endpoint'e istek atar:
//
//	GET /api/resource/{resource}/resolver/{field_name}?id={related_id}
//	POST /api/resource/{resource}/resolver/{field_name} (body: {id: related_id})
//
// Döndürür:
//   - HasOneField pointer'ı (method chaining için)
func (h *HasOneField) ResolveHoverCard(resolver HoverCardResolver) *HasOneField {
	if h.hoverCardConfig == nil {
		h.hoverCardConfig = NewHoverCardConfig()
	}
	h.hoverCardConfig.SetResolver(resolver)
	return h
}

// GetHoverCard, hover card konfigürasyonunu döndürür.
//
// Bu metod, hover card konfigürasyonunu alır.
//
// Döndürür:
//   - HoverCardConfig pointer'ı (nil olabilir)
func (h *HasOneField) GetHoverCard() *HoverCardConfig {
	return h.hoverCardConfig
}

// GetRelatedTableName, ilişkili tablo adını döndürür.
//
// Bu metod, HasOne ilişkisinde kullanılan ilişkili tablonun adını döndürür.
// Raw SQL sorguları için kullanılır.
//
// # Dönüş Değeri
//
// - İlişkili tablo adı (örn. "profiles", "invoices", "capitals")
//
// # Kullanım Örneği
//
//	field := fields.HasOne("Profile", "profile", "profiles")
//	tableName := field.GetRelatedTableName() // "profiles"
//
// Döndürür:
//   - İlişkili tablo adı
func (h *HasOneField) GetRelatedTableName() string {
	return h.RelatedResourceSlug
}

// GetForeignKeyColumn, foreign key sütun adını döndürür.
//
// Bu metod, HasOne ilişkisinde kullanılan foreign key sütununun adını döndürür.
// Foreign key, ilişkili tablodaki referans sütunudur.
//
// # Dönüş Değeri
//
// - Foreign key sütun adı (örn. "user_id", "order_id", "country_id")
//
// # Kullanım Örneği
//
//	field := fields.HasOne("Profile", "profile", "profiles").ForeignKey("user_id")
//	foreignKey := field.GetForeignKeyColumn() // "user_id"
//
// Döndürür:
//   - Foreign key sütun adı
func (h *HasOneField) GetForeignKeyColumn() string {
	return h.ForeignKeyColumn
}

// GetOwnerKeyColumn, owner key sütun adını döndürür.
//
// Bu metod, HasOne ilişkisinde kullanılan owner key sütununun adını döndürür.
// Owner key, ana tablodaki referans sütunudur (genellikle primary key).
//
// # Dönüş Değeri
//
// - Owner key sütun adı (örn. "id", "uuid")
//
// # Kullanım Örneği
//
//	field := fields.HasOne("Profile", "profile", "profiles").OwnerKey("id")
//	ownerKey := field.GetOwnerKeyColumn() // "id"
//
// Döndürür:
//   - Owner key sütun adı
func (h *HasOneField) GetOwnerKeyColumn() string {
	return h.OwnerKeyColumn
}
