package data

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/query"
)

// Sort, veritabanı sorgularında sıralama işlemleri için kullanılan yapıdır.
//
// Bu yapı, bir sorgunun hangi kolona göre ve hangi yönde (artan/azalan) sıralanacağını belirtir.
// Genellikle QueryRequest içinde bir dizi olarak kullanılır ve birden fazla sıralama kriteri
// tanımlanabilir.
//
// # Alanlar
//
// - `Column`: Sıralama yapılacak kolon adı (örn: "created_at", "name", "price")
// - `Direction`: Sıralama yönü - "asc" (artan) veya "desc" (azalan)
//
// # Kullanım Senaryoları
//
// - Tablo verilerini kullanıcı tercihine göre sıralama
// - Çoklu kolon sıralaması (önce ada göre, sonra tarihe göre)
// - API isteklerinde dinamik sıralama parametreleri
// - Admin panellerinde tablo başlıklarına tıklayarak sıralama
//
// # Örnek Kullanım
//
// ```go
// // Tek sıralama
// sort := Sort{
//     Column:    "created_at",
//     Direction: "desc",
// }
//
// // Çoklu sıralama
// sorts := []Sort{
//     {Column: "status", Direction: "asc"},
//     {Column: "created_at", Direction: "desc"},
// }
// ```
//
// # Önemli Notlar
//
// - Direction değeri genellikle "asc" veya "desc" olmalıdır
// - Column değeri veritabanı şemasında mevcut bir kolon olmalıdır
// - SQL injection saldırılarına karşı Column değeri doğrulanmalıdır
// - Boş Direction değeri varsayılan olarak "asc" kabul edilebilir
type Sort struct {
	Column    string `json:"column"`
	Direction string `json:"direction"`
}

// QueryRequest, veri sorgulama işlemleri için kullanılan istek yapısıdır.
//
// Bu yapı, sayfalama, sıralama, filtreleme ve arama gibi tüm sorgu parametrelerini
// tek bir yapıda toplar. RESTful API'lerde liste endpoint'lerinde kullanılır ve
// kullanıcının veri üzerinde esnek sorgulama yapmasını sağlar.
//
// # Alanlar
//
// - `Page`: Sayfa numarası (1'den başlar)
// - `PerPage`: Sayfa başına gösterilecek kayıt sayısı
// - `Sorts`: Sıralama kriterleri dizisi (birden fazla sıralama desteklenir)
// - `Filters`: Filtreleme kriterleri dizisi (karmaşık filtreler desteklenir)
// - `Search`: Genel arama terimi (birden fazla kolonda arama yapar)
//
// # Kullanım Senaryoları
//
// - Admin panellerinde tablo verilerini listeleme
// - API'de sayfalanmış veri döndürme
// - Gelişmiş filtreleme ve arama özellikleri
// - Kullanıcı tercihlerini kaydetme ve geri yükleme
// - Export işlemlerinde veri seçimi
//
// # Örnek Kullanım
//
// ```go
// // Basit sayfalama
// req := QueryRequest{
//     Page:    1,
//     PerPage: 20,
// }
//
// // Arama ile birlikte
// req := QueryRequest{
//     Page:    1,
//     PerPage: 20,
//     Search:  "john",
// }
//
// // Tam özellikli sorgu
// req := QueryRequest{
//     Page:    2,
//     PerPage: 50,
//     Search:  "active",
//     Sorts: []Sort{
//         {Column: "created_at", Direction: "desc"},
//     },
//     Filters: []query.Filter{
//         {Field: "status", Operator: "=", Value: "active"},
//         {Field: "price", Operator: ">", Value: 100},
//     },
// }
// ```
//
// # Avantajlar
//
// - Tüm sorgu parametreleri tek bir yapıda toplanır
// - JSON serileştirme/deserileştirme desteği
// - Tip güvenli parametre yönetimi
// - Genişletilebilir yapı
//
// # Önemli Notlar
//
// - Page değeri 0 veya negatif olmamalıdır (genellikle 1'den başlar)
// - PerPage değeri makul bir üst sınıra sahip olmalıdır (örn: 100)
// - Search terimi SQL injection'a karşı temizlenmelidir
// - Filters dizisi boş olabilir (filtresiz sorgu)
// - Sorts dizisi boş olabilir (varsayılan sıralama kullanılır)
type QueryRequest struct {
	Page    int            `json:"page"`
	PerPage int            `json:"per_page"`
	Sorts   []Sort         `json:"sorts"`
	Filters []query.Filter `json:"filters"`
	Search  string         `json:"search"`
}

// QueryResponse, veri sorgulama işlemlerinin sonucunu temsil eden yapıdır.
//
// Bu yapı, sayfalanmış veri listesi ile birlikte sayfalama bilgilerini içerir.
// API yanıtlarında standart bir format sağlar ve frontend'in sayfalama kontrollerini
// oluşturması için gerekli tüm bilgileri içerir.
//
// # Alanlar
//
// - `Items`: Sorgu sonucunda dönen kayıtlar dizisi (herhangi bir tip olabilir)
// - `Total`: Toplam kayıt sayısı (filtreleme uygulandıktan sonra)
// - `Page`: Mevcut sayfa numarası
// - `PerPage`: Sayfa başına kayıt sayısı
//
// # Kullanım Senaryoları
//
// - API'den sayfalanmış veri döndürme
// - Frontend tablo bileşenlerine veri sağlama
// - Sayfalama kontrollerini oluşturma (toplam sayfa sayısı hesaplama)
// - Export işlemlerinde veri önizleme
// - Infinite scroll implementasyonu
//
// # Örnek Kullanım
//
// ```go
// // Başarılı sorgu yanıtı
// response := &QueryResponse{
//     Items:   users,
//     Total:   150,
//     Page:    1,
//     PerPage: 20,
// }
//
// // Toplam sayfa sayısını hesaplama
// totalPages := (response.Total + int64(response.PerPage) - 1) / int64(response.PerPage)
//
// // Sonraki sayfa var mı kontrolü
// hasNextPage := response.Page * response.PerPage < int(response.Total)
// ```
//
// # Hesaplanabilir Değerler
//
// Bu yapıdan şu değerler hesaplanabilir:
// - Toplam sayfa sayısı: `ceil(Total / PerPage)`
// - Sonraki sayfa var mı: `Page * PerPage < Total`
// - Önceki sayfa var mı: `Page > 1`
// - Mevcut sayfadaki kayıt aralığı: `[(Page-1)*PerPage + 1, min(Page*PerPage, Total)]`
//
// # Avantajlar
//
// - Standart API yanıt formatı
// - Frontend sayfalama için gerekli tüm bilgi
// - JSON serileştirme desteği
// - Tip güvenli veri yönetimi
//
// # Önemli Notlar
//
// - Items boş dizi olabilir (sonuç bulunamadığında)
// - Total değeri 0 olabilir (hiç kayıt yoksa)
// - Items dizisinin uzunluğu PerPage'den küçük olabilir (son sayfa)
// - Items interface{} tipinde olduğu için tip dönüşümü gerekebilir
type QueryResponse struct {
	Items   []interface{} `json:"items"`
	Total   int64         `json:"total"`
	Page    int           `json:"page"`
	PerPage int           `json:"per_page"`
}

// DataProvider, veri sağlayıcıları için standart CRUD operasyonlarını tanımlayan interface'dir.
//
// Bu interface, farklı veri kaynaklarıyla (GORM, Ent, MongoDB, REST API vb.) çalışabilen
// soyut bir katman sağlar. Repository pattern'in bir uygulamasıdır ve veri erişim
// mantığını iş mantığından ayırır.
//
// # Metodlar
//
// ## CRUD Operasyonları
//
// ### Index
// Sayfalanmış, filtrelenmiş ve sıralanmış veri listesi döndürür.
//
// ### Show
// Belirli bir kaydı ID'sine göre getirir.
//
// ### Create
// Yeni bir kayıt oluşturur.
//
// ### Update
// Mevcut bir kaydı günceller.
//
// ### Delete
// Bir kaydı siler (soft delete veya hard delete olabilir).
//
// ## Configuration Metodları
//
// ### SetSearchColumns
// Arama işleminde kullanılacak kolonları ayarlar.
//
// ### SetWith
// Eager loading için ilişkileri ayarlar (JOIN veya preload).
//
// ### SetRelationshipFields
// Relationship field'larını ayarlar.
//
// ## Dinamik Query Metodları
//
// ### QueryTable
// Dinamik tablo query'leri için kullanılır (MorphTo, HasOne, BelongsTo için).
//
// Parametreler:
// - `ctx`: İstek bağlamı
// - `table`: Tablo adı
// - `conditions`: WHERE koşulları (key-value map)
//
// Döndürür:
// - `[]map[string]interface{}`: Sorgu sonuçları
// - `error`: Hata durumunda hata mesajı
//
// Örnek:
// ```go
// results, err := provider.QueryTable(ctx, "users", map[string]interface{}{
//     "id": "123",
//     "status": "active",
// })
// ```
//
// ### QueryRelationship
// Relationship query'leri için kullanılır (display field çekmek için).
//
// Parametreler:
// - `ctx`: İstek bağlamı
// - `relationshipType`: İlişki tipi (tablo adı)
// - `foreignKey`: Foreign key kolon adı
// - `foreignValue`: Foreign key değeri
// - `displayField`: Gösterilecek field adı
//
// Döndürür:
// - `interface{}`: Display field değeri
// - `error`: Hata durumunda hata mesajı
//
// Örnek:
// ```go
// displayValue, err := provider.QueryRelationship(ctx, "users", "id", "123", "name")
// ```
//
// ## Transaction Metodları
//
// ### BeginTx
// Yeni bir transaction başlatır ve transaction içinde çalışan yeni bir Provider döndürür.
//
// ### Commit
// Aktif transaction'ı commit eder.
//
// ### Rollback
// Aktif transaction'ı rollback eder.
//
// Örnek:
// ```go
// txProvider, err := provider.BeginTx(ctx)
// if err != nil {
//     return err
// }
//
// _, err = txProvider.Create(ctx, data1)
// if err != nil {
//     txProvider.Rollback()
//     return err
// }
//
// _, err = txProvider.Create(ctx, data2)
// if err != nil {
//     txProvider.Rollback()
//     return err
// }
//
// return txProvider.Commit()
// ```
//
// ## Raw Query Metodları
//
// ### Raw
// Raw SQL query çalıştırır ve sonuçları döndürür (SELECT için).
//
// ### Exec
// Raw SQL query çalıştırır (INSERT, UPDATE, DELETE için).
//
// Örnek:
// ```go
// results, err := provider.Raw(ctx, "SELECT * FROM users WHERE status = ?", "active")
// err := provider.Exec(ctx, "UPDATE users SET status = ? WHERE id = ?", "inactive", "123")
// ```
//
// # Kullanım Senaryoları
//
// - RESTful API endpoint'lerinde veri yönetimi
// - Admin panellerinde CRUD operasyonları
// - Farklı veritabanı sistemleri arasında geçiş yapabilme
// - Test edilebilir kod yazma (mock implementasyonlar)
// - Mikroservis mimarisinde veri katmanı soyutlama
// - Dinamik query'ler (MorphTo, HasOne, BelongsTo)
// - Transaction yönetimi
// - Raw SQL query'ler
//
// # Avantajlar
//
// - **Tam Soyutlama**: ORM'den tamamen bağımsız kod yazma
// - **Test Edilebilirlik**: Mock implementasyonlar ile kolay test
// - **Esneklik**: Farklı veri kaynakları için aynı interface
// - **Bakım Kolaylığı**: Veri erişim mantığı tek yerde
// - **Yeniden Kullanılabilirlik**: Aynı interface farklı modeller için kullanılabilir
// - **Güvenlik**: Handler'lar ORM'e direkt erişemez
//
// # Önemli Notlar
//
// - **GetClient() Yok**: Interface'de GetClient() metodu YOK! Handler'lar ORM'e direkt erişemez.
// - **Private getClient()**: Implementasyonlar internal olarak private getClient() kullanabilir.
// - **Thread Safety**: Implementasyonlar thread-safe olmalıdır
// - **Transaction Desteği**: BeginTx, Commit, Rollback metodları ile transaction yönetimi
// - **Error Handling**: Standart error tipleri kullanılmalıdır (NotFound, ValidationError vb.)
//
// # İlgili Tipler
//
// - `QueryRequest`: Sorgu parametreleri
// - `QueryResponse`: Sorgu yanıtı
// - `Sort`: Sıralama bilgisi
// - `query.Filter`: Filtreleme bilgisi
// - `context.Context`: İstek bağlamı
type DataProvider interface {
	// CRUD Operasyonları
	Index(ctx *context.Context, req QueryRequest) (*QueryResponse, error)
	Show(ctx *context.Context, id string) (interface{}, error)
	Create(ctx *context.Context, data map[string]interface{}) (interface{}, error)
	Update(ctx *context.Context, id string, data map[string]interface{}) (interface{}, error)
	Delete(ctx *context.Context, id string) error

	// Configuration
	SetSearchColumns(cols []string)
	SetWith(rels []string)
	SetRelationshipFields(fields []fields.RelationshipField)

	// Dinamik Query'ler (Card, MorphTo, HasOne için)
	QueryTable(ctx *context.Context, table string, conditions map[string]interface{}) ([]map[string]interface{}, error)
	QueryRelationship(ctx *context.Context, relationshipType string, foreignKey string, foreignValue interface{}, displayField string) (interface{}, error)

	// Transaction (Action'lar için)
	BeginTx(ctx *context.Context) (DataProvider, error)
	Commit() error
	Rollback() error

	// Raw Query (özel durumlar için)
	Raw(ctx *context.Context, sql string, args ...interface{}) ([]map[string]interface{}, error)
	Exec(ctx *context.Context, sql string, args ...interface{}) error

	// ORM Client (Card.Resolve ve Action.Execute için)
	// UYARI: Bu metod sadece Card.Resolve() ve Action.Execute() gibi özel durumlar için kullanılmalıdır.
	// Normal işlemler için Provider metodlarını (Raw, Exec, QueryTable, vb.) kullanın.
	GetClient() interface{}
}
