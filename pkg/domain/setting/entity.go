// Package setting, uygulama ayarlarının yönetimi için domain katmanını sağlar.
// Bu paket, sistem genelinde kullanılan yapılandırma değerlerinin depolanması,
// alınması ve güncellenmesi işlemlerini yönetir.
package setting

import (
	"context"
	"time"
)

// Bu yapı, uygulama ayarlarını temsil eder ve veritabanında depolanır.
//
// Setting yapısı, sistem genelinde kullanılan yapılandırma değerlerini
// yönetmek için kullanılır. Her ayar, benzersiz bir anahtar (Key) ile
// tanımlanır ve farklı türlerde değerler (string, integer, boolean, json)
// depolanabilir.
//
// Kullanım Senaryoları:
// - Uygulama başlangıç ayarlarının depolanması
// - Dinamik yapılandırma değerlerinin yönetimi
// - Sistem genelinde erişilebilir ayarların merkezi yönetimi
// - Kullanıcı tercihlerinin depolanması
//
// Örnek Kullanım:
//   setting := &Setting{
//       Key:   "app.name",
//       Value: "Panel Application",
//       Type:  "string",
//       Group: "app",
//       Label: "Uygulama Adı",
//       Help:  "Uygulamanın görüntülenecek adı",
//   }
//
// Önemli Notlar:
// - Key alanı benzersiz olmalı ve sistem genelinde tekrar etmemelidir
// - Value alanı metin formatında depolanır, JSON değerler string olarak saklanır
// - Group alanı ayarları kategorize etmek için kullanılır
// - CreatedAt ve UpdatedAt otomatik olarak GORM tarafından yönetilir
type Setting struct {
	// ID, ayarın benzersiz tanımlayıcısıdır ve birincil anahtar olarak kullanılır.
	// Veritabanında otomatik olarak artan bir değerdir.
	// Örnek: 1, 2, 3, ...
	ID uint `json:"id" gorm:"primaryKey"`

	// Key, ayarın benzersiz anahtarıdır ve sistem genelinde tekrar etmemelidir.
	// Ayarı tanımlamak ve erişmek için kullanılır.
	// Örnek: "app.name", "app.version", "email.smtp.host"
	// Uyarı: Bu alan benzersiz indekslenmiştir, aynı Key ile iki ayar oluşturulamaz.
	Key string `json:"key" gorm:"uniqueIndex"`

	// Value, ayarın gerçek değeridir ve metin formatında depolanır.
	// JSON, integer, boolean gibi karmaşık değerler string olarak saklanır.
	// Örnek: "true", "123", "{\"host\":\"localhost\"}"
	// Not: Büyük değerler için text türü kullanılır, performans için optimize edilmiştir.
	Value string `json:"value" gorm:"type:text"`

	// Type, Value alanının veri türünü belirtir.
	// Desteklenen türler: "string", "integer", "boolean", "json"
	// Örnek: "string", "integer", "boolean", "json"
	// Kullanım: Değeri doğru türe dönüştürmek için kullanılır.
	Type string `json:"type"`

	// Group, ayarları kategorize etmek için kullanılır.
	// Aynı grup içindeki ayarlar birlikte yönetilir.
	// Örnek: "app", "email", "database", "security"
	// Not: Bu alan indekslenmiştir, grup bazında hızlı sorgu için optimize edilmiştir.
	Group string `json:"group" gorm:"index"`

	// Label, ayarın kullanıcı dostu adıdır ve UI'da gösterilir.
	// Örnek: "Uygulama Adı", "SMTP Sunucusu", "Maksimum Dosya Boyutu"
	// Kullanım: Admin panelinde ayarları görüntülerken kullanılır.
	Label string `json:"label"`

	// Help, ayarın açıklaması ve kullanım talimatlarıdır.
	// Örnek: "Uygulamanın başlık çubuğunda gösterilecek adı", "E-posta göndermek için SMTP sunucusu"
	// Kullanım: Kullanıcılara ayarın ne için olduğunu açıklamak için kullanılır.
	Help string `json:"help"`

	// CreatedAt, ayarın oluşturulduğu tarih ve saattir.
	// GORM tarafından otomatik olarak ayarlanır.
	// Örnek: 2024-01-15 10:30:45
	// Not: Bu alan indekslenmiştir, tarih bazında sorgu performansı için optimize edilmiştir.
	CreatedAt time.Time `json:"createdAt" gorm:"index"`

	// UpdatedAt, ayarın son güncellendiği tarih ve saattir.
	// GORM tarafından otomatik olarak güncellenir.
	// Örnek: 2024-01-20 14:22:10
	// Not: Bu alan indekslenmiştir, en son güncellemeleri bulmak için kullanılır.
	UpdatedAt time.Time `json:"updatedAt" gorm:"index"`
}

// Bu interface, Setting veri modeli için repository pattern'ini tanımlar.
//
// Repository interface'i, Setting nesneleri üzerinde CRUD (Create, Read, Update, Delete)
// işlemlerini gerçekleştirmek için gerekli metotları belirtir. Bu interface, veri
// erişim katmanını soyutlar ve farklı veri kaynakları (veritabanı, cache, vb.)
// ile çalışmayı mümkün kılar.
//
// Kullanım Senaryoları:
// - Ayarları veritabanından almak ve depolamak
// - Belirli bir ayarı ID veya Key ile sorgulamak
// - Bir grup içindeki tüm ayarları listelemek
// - Ayarları güncellemek ve silmek
// - Dependency injection aracılığıyla farklı implementasyonlar kullanmak
//
// Implementasyon Örneği:
//   type SettingRepository struct {
//       db *gorm.DB
//   }
//
//   func (r *SettingRepository) Create(ctx context.Context, setting *Setting) error {
//       return r.db.WithContext(ctx).Create(setting).Error
//   }
//
// Önemli Notlar:
// - Tüm metotlar context parametresi alır, iptal edilebilirlik için
// - Tüm metotlar hata döndürür, hata yönetimi için
// - FindByGroup birden fazla Setting döndürür, diğer metotlar tek bir Setting döndürür
type Repository interface {
	// Bu metod, yeni bir Setting kaydını veritabanına oluşturur.
	//
	// Parametreler:
	// - ctx: İşlemin iptal edilebilmesi için context
	// - setting: Oluşturulacak Setting nesnesi pointer'ı
	//
	// Dönüş Değeri:
	// - error: İşlem başarılı ise nil, hata durumunda error nesnesi
	//
	// Kullanım Örneği:
	//   setting := &Setting{
	//       Key:   "app.name",
	//       Value: "My App",
	//       Type:  "string",
	//       Group: "app",
	//   }
	//   err := repo.Create(ctx, setting)
	//   if err != nil {
	//       log.Fatal(err)
	//   }
	//
	// Uyarılar:
	// - Key alanı benzersiz olmalı, aksi takdirde hata döner
	// - Setting nesnesi nil olmamalı
	// - Context iptal edilirse işlem durdurulur
	Create(ctx context.Context, setting *Setting) error

	// Bu metod, ID'si verilen Setting kaydını veritabanından alır.
	//
	// Parametreler:
	// - ctx: İşlemin iptal edilebilmesi için context
	// - id: Aranacak Setting'in ID'si
	//
	// Dönüş Değeri:
	// - *Setting: Bulunan Setting nesnesi pointer'ı
	// - error: Kayıt bulunamadı ise error, başarılı ise nil
	//
	// Kullanım Örneği:
	//   setting, err := repo.FindByID(ctx, 1)
	//   if err != nil {
	//       log.Fatal(err)
	//   }
	//   fmt.Println(setting.Value)
	//
	// Uyarılar:
	// - Kayıt bulunamadı ise error döner (gorm.ErrRecordNotFound)
	// - ID geçerli bir pozitif sayı olmalı
	FindByID(ctx context.Context, id uint) (*Setting, error)

	// Bu metod, Key'i verilen Setting kaydını veritabanından alır.
	//
	// Parametreler:
	// - ctx: İşlemin iptal edilebilmesi için context
	// - key: Aranacak Setting'in Key'i
	//
	// Dönüş Değeri:
	// - *Setting: Bulunan Setting nesnesi pointer'ı
	// - error: Kayıt bulunamadı ise error, başarılı ise nil
	//
	// Kullanım Örneği:
	//   setting, err := repo.FindByKey(ctx, "app.name")
	//   if err != nil {
	//       log.Fatal(err)
	//   }
	//   fmt.Println(setting.Value) // "My App"
	//
	// Uyarılar:
	// - Kayıt bulunamadı ise error döner (gorm.ErrRecordNotFound)
	// - Key boş string olmamalı
	// - Key büyük/küçük harfe duyarlı olabilir
	FindByKey(ctx context.Context, key string) (*Setting, error)

	// Bu metod, Group'u verilen tüm Setting kayıtlarını veritabanından alır.
	//
	// Parametreler:
	// - ctx: İşlemin iptal edilebilmesi için context
	// - group: Aranacak Setting'lerin Group'u
	//
	// Dönüş Değeri:
	// - []Setting: Bulunan Setting nesnelerinin slice'ı
	// - error: Hata durumunda error nesnesi, başarılı ise nil
	//
	// Kullanım Örneği:
	//   settings, err := repo.FindByGroup(ctx, "email")
	//   if err != nil {
	//       log.Fatal(err)
	//   }
	//   for _, setting := range settings {
	//       fmt.Println(setting.Key, setting.Value)
	//   }
	//
	// Uyarılar:
	// - Grup bulunamadı ise boş slice döner (hata değil)
	// - Group boş string olmamalı
	// - Büyük veri setleri için pagination kullanılması önerilir
	FindByGroup(ctx context.Context, group string) ([]Setting, error)

	// Bu metod, mevcut bir Setting kaydını veritabanında günceller.
	//
	// Parametreler:
	// - ctx: İşlemin iptal edilebilmesi için context
	// - setting: Güncellenecek Setting nesnesi pointer'ı
	//
	// Dönüş Değeri:
	// - error: İşlem başarılı ise nil, hata durumunda error nesnesi
	//
	// Kullanım Örneği:
	//   setting, _ := repo.FindByKey(ctx, "app.name")
	//   setting.Value = "Updated App Name"
	//   err := repo.Update(ctx, setting)
	//   if err != nil {
	//       log.Fatal(err)
	//   }
	//
	// Uyarılar:
	// - Setting nesnesi nil olmamalı
	// - Setting'in ID'si geçerli olmalı
	// - UpdatedAt alanı otomatik olarak güncellenir
	// - Key alanı değiştirilmemelidir (benzersiz kısıt nedeniyle)
	Update(ctx context.Context, setting *Setting) error

	// Bu metod, ID'si verilen Setting kaydını veritabanından siler.
	//
	// Parametreler:
	// - ctx: İşlemin iptal edilebilmesi için context
	// - id: Silinecek Setting'in ID'si
	//
	// Dönüş Değeri:
	// - error: İşlem başarılı ise nil, hata durumunda error nesnesi
	//
	// Kullanım Örneği:
	//   err := repo.Delete(ctx, 1)
	//   if err != nil {
	//       log.Fatal(err)
	//   }
	//
	// Uyarılar:
	// - Silme işlemi geri alınamaz
	// - ID geçerli bir pozitif sayı olmalı
	// - Kayıt bulunamadı ise hata döner
	// - Silme işleminden önce bağımlılıkları kontrol edin
	Delete(ctx context.Context, id uint) error
}
