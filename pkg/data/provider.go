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
// Bu interface, farklı veri kaynaklarıyla (GORM, MongoDB, REST API vb.) çalışabilen
// soyut bir katman sağlar. Repository pattern'in bir uygulamasıdır ve veri erişim
// mantığını iş mantığından ayırır.
//
// # Metodlar
//
// ## Index
// Sayfalanmış, filtrelenmiş ve sıralanmış veri listesi döndürür.
//
// Parametreler:
// - `ctx`: İstek bağlamı (kullanıcı, yetkilendirme, transaction vb.)
// - `req`: Sorgu parametreleri (sayfalama, filtreleme, sıralama, arama)
//
// Döndürür:
// - `*QueryResponse`: Sayfalanmış veri ve meta bilgiler
// - `error`: Hata durumunda hata mesajı
//
// ## Show
// Belirli bir kaydı ID'sine göre getirir.
//
// Parametreler:
// - `ctx`: İstek bağlamı
// - `id`: Kaydın benzersiz kimliği (string formatında)
//
// Döndürür:
// - `interface{}`: Bulunan kayıt (tip dönüşümü gerekebilir)
// - `error`: Kayıt bulunamazsa veya hata durumunda hata mesajı
//
// ## Create
// Yeni bir kayıt oluşturur.
//
// Parametreler:
// - `ctx`: İstek bağlamı
// - `data`: Oluşturulacak kaydın verileri (key-value map)
//
// Döndürür:
// - `interface{}`: Oluşturulan kayıt (genellikle ID ile birlikte)
// - `error`: Validasyon hatası veya veritabanı hatası
//
// ## Update
// Mevcut bir kaydı günceller.
//
// Parametreler:
// - `ctx`: İstek bağlamı
// - `id`: Güncellenecek kaydın ID'si
// - `data`: Güncellenecek alanlar (key-value map, sadece değişen alanlar)
//
// Döndürür:
// - `interface{}`: Güncellenmiş kayıt
// - `error`: Kayıt bulunamazsa veya güncelleme başarısız olursa hata
//
// ## Delete
// Bir kaydı siler (soft delete veya hard delete olabilir).
//
// Parametreler:
// - `ctx`: İstek bağlamı
// - `id`: Silinecek kaydın ID'si
//
// Döndürür:
// - `error`: Kayıt bulunamazsa veya silme başarısız olursa hata
//
// ## SetSearchColumns
// Arama işleminde kullanılacak kolonları ayarlar.
//
// Parametreler:
// - `cols`: Arama yapılacak kolon isimleri dizisi
//
// Örnek: `[]string{"name", "email", "description"}`
//
// ## SetWith
// Eager loading için ilişkileri ayarlar (JOIN veya preload).
//
// Parametreler:
// - `rels`: Yüklenecek ilişki isimleri dizisi
//
// Örnek: `[]string{"User", "Category", "Tags"}`
//
// # Kullanım Senaryoları
//
// - RESTful API endpoint'lerinde veri yönetimi
// - Admin panellerinde CRUD operasyonları
// - Farklı veritabanı sistemleri arasında geçiş yapabilme
// - Test edilebilir kod yazma (mock implementasyonlar)
// - Mikroservis mimarisinde veri katmanı soyutlama
//
// # Örnek Implementasyon
//
// ```go
// type UserProvider struct {
//     db *gorm.DB
//     searchColumns []string
//     withRelations []string
// }
//
// func (p *UserProvider) Index(ctx *context.Context, req QueryRequest) (*QueryResponse, error) {
//     query := p.db.Model(&User{})
//
//     // Arama uygula
//     if req.Search != "" {
//         // searchColumns kullanarak arama yap
//     }
//
//     // Filtreleri uygula
//     for _, filter := range req.Filters {
//         // Filter uygula
//     }
//
//     // Sıralama uygula
//     for _, sort := range req.Sorts {
//         query = query.Order(sort.Column + " " + sort.Direction)
//     }
//
//     // Sayfalama uygula
//     offset := (req.Page - 1) * req.PerPage
//     query = query.Offset(offset).Limit(req.PerPage)
//
//     var users []User
//     var total int64
//
//     p.db.Model(&User{}).Count(&total)
//     query.Find(&users)
//
//     return &QueryResponse{
//         Items:   convertToInterface(users),
//         Total:   total,
//         Page:    req.Page,
//         PerPage: req.PerPage,
//     }, nil
// }
// ```
//
// # Örnek Kullanım
//
// ```go
// // Provider oluşturma
// provider := NewGormProvider(db, &User{})
// provider.SetSearchColumns([]string{"name", "email"})
// provider.SetWith([]string{"Profile", "Posts"})
//
// // Liste getirme
// response, err := provider.Index(ctx, QueryRequest{
//     Page:    1,
//     PerPage: 20,
//     Search:  "john",
//     Sorts: []Sort{{Column: "created_at", Direction: "desc"}},
// })
//
// // Tekil kayıt getirme
// user, err := provider.Show(ctx, "123")
//
// // Yeni kayıt oluşturma
// newUser, err := provider.Create(ctx, map[string]interface{}{
//     "name":  "John Doe",
//     "email": "john@example.com",
// })
//
// // Güncelleme
// updated, err := provider.Update(ctx, "123", map[string]interface{}{
//     "name": "Jane Doe",
// })
//
// // Silme
// err := provider.Delete(ctx, "123")
// ```
//
// # Avantajlar
//
// - **Soyutlama**: Veri kaynağından bağımsız kod yazma
// - **Test Edilebilirlik**: Mock implementasyonlar ile kolay test
// - **Esneklik**: Farklı veri kaynakları için aynı interface
// - **Bakım Kolaylığı**: Veri erişim mantığı tek yerde
// - **Yeniden Kullanılabilirlik**: Aynı interface farklı modeller için kullanılabilir
//
// # Dezavantajlar
//
// - **Performans**: Soyutlama katmanı minimal overhead ekler
// - **Karmaşıklık**: Basit uygulamalar için fazla soyutlama olabilir
// - **Tip Güvenliği**: interface{} kullanımı tip dönüşümü gerektirir
//
// # Önemli Notlar
//
// - **Thread Safety**: Implementasyonlar thread-safe olmalıdır
// - **Transaction Desteği**: Context üzerinden transaction yönetimi yapılabilir
// - **Error Handling**: Standart error tipleri kullanılmalıdır (NotFound, ValidationError vb.)
// - **Validation**: Create ve Update metodları veri validasyonu yapmalıdır
// - **Authorization**: Context üzerinden yetkilendirme kontrolleri yapılabilir
// - **Soft Delete**: Delete metodu soft delete destekleyebilir
// - **Audit Log**: Tüm operasyonlar audit log'a kaydedilebilir
// - **Caching**: Implementasyonlar cache katmanı ekleyebilir
// - **Rate Limiting**: Context üzerinden rate limiting uygulanabilir
//
// # İlgili Tipler
//
// - `QueryRequest`: Sorgu parametreleri
// - `QueryResponse`: Sorgu yanıtı
// - `Sort`: Sıralama bilgisi
// - `query.Filter`: Filtreleme bilgisi
// - `context.Context`: İstek bağlamı
type DataProvider interface {
	Index(ctx *context.Context, req QueryRequest) (*QueryResponse, error)
	Show(ctx *context.Context, id string) (interface{}, error)
	Create(ctx *context.Context, data map[string]interface{}) (interface{}, error)
	Update(ctx *context.Context, id string, data map[string]interface{}) (interface{}, error)
	Delete(ctx *context.Context, id string) error
	SetSearchColumns(cols []string)
	SetWith(rels []string)
	SetRelationshipFields(fields []fields.RelationshipField)
}
