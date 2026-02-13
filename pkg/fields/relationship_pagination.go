// Package fields, ilişkisel veritabanı alanları ve sayfalama işlevselliği sağlar.
//
// Bu paket, ilişkili kayıtlar için sayfalama yönetimi ve metadata işlemlerini içerir.
// Sayfa numarası, sayfa boyutu ve toplam kayıt sayısı gibi bilgileri yönetir.
//
// Daha fazla bilgi için docs/Relationships.md dosyasına bakın.
package fields

import (
	"context"
)

// RelationshipPagination, ilişkiler için sayfalama işlevselliğini yönetir.
//
// Bu interface, ilişkili kayıtları sayfalamak için metodlar sağlar.
// Sayfa numarası ve sayfa boyutu belirlenebilir, sayfalama metadata'sı alınabilir.
//
// # Kullanım Örneği
//
//	pagination := fields.NewRelationshipPagination(field)
//	results, err := pagination.ApplyPagination(ctx, 1, 15)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	pageInfo := pagination.GetPageInfo()
//
// Daha fazla bilgi için docs/Relationships.md dosyasına bakın.
type RelationshipPagination interface {
	// ApplyPagination, ilişki sorgusuna sayfalama uygular.
	//
	// # Parametreler
	//
	//   - ctx: Context nesnesi
	//   - page: Sayfa numarası (minimum 1, geçersiz değerler 1'e ayarlanır)
	//   - perPage: Sayfa başına kayıt sayısı (minimum 1, maksimum 100)
	//
	// # Dönüş Değerleri
	//
	//   - []interface{}: Sayfalanmış kayıt listesi
	//   - error: Hata durumunda hata nesnesi
	//
	// # Önemli Notlar
	//
	//   - Sayfa numarası 1'den küçükse otomatik olarak 1'e ayarlanır
	//   - Sayfa boyutu 1'den küçükse otomatik olarak 15'e ayarlanır
	//   - Sayfa boyutu 100'den büyükse otomatik olarak 100'e sınırlanır
	ApplyPagination(ctx context.Context, page int, perPage int) ([]interface{}, error)

	// GetPageInfo, sayfalama metadata'sını döndürür.
	//
	// # Dönüş Değeri
	//
	// Map içeriği:
	//   - current_page: Mevcut sayfa numarası
	//   - per_page: Sayfa başına kayıt sayısı
	//   - total: Toplam kayıt sayısı
	//   - total_pages: Toplam sayfa sayısı
	//   - from: Başlangıç kayıt indeksi
	//   - to: Bitiş kayıt indeksi
	//
	// # Kullanım Örneği
	//
	//	info := pagination.GetPageInfo()
	//	currentPage := info["current_page"].(int)
	//	totalPages := info["total_pages"].(int64)
	GetPageInfo() map[string]interface{}
}

// RelationshipPaginationImpl, RelationshipPagination interface'ini uygular.
//
// Bu struct, ilişki sayfalama işlemlerinin somut implementasyonunu sağlar.
// Sayfa numarası, sayfa boyutu ve toplam kayıt sayısı bilgilerini tutar.
//
// # Alanlar
//
//   - field: İlişki alanı referansı
//   - page: Mevcut sayfa numarası (varsayılan: 1)
//   - perPage: Sayfa başına kayıt sayısı (varsayılan: 15)
//   - total: Toplam kayıt sayısı (varsayılan: 0)
//
// # Kullanım Örneği
//
//	impl := &RelationshipPaginationImpl{
//	    field:   myField,
//	    page:    1,
//	    perPage: 20,
//	    total:   100,
//	}
//
// Daha fazla bilgi için docs/Relationships.md dosyasına bakın.
type RelationshipPaginationImpl struct {
	field   RelationshipField
	page    int
	perPage int
	total   int64
}

// NewRelationshipPagination, yeni bir ilişki sayfalama yöneticisi oluşturur.
//
// Bu constructor, varsayılan değerlerle yapılandırılmış bir sayfalama nesnesi döndürür.
// Sayfa numarası 1, sayfa boyutu 15 ve toplam kayıt sayısı 0 olarak başlatılır.
//
// # Parametreler
//
//   - field: İlişki alanı (RelationshipField interface'ini uygulayan nesne)
//
// # Dönüş Değeri
//
//   - *RelationshipPaginationImpl: Yapılandırılmış sayfalama nesnesi pointer'ı
//
// # Kullanım Örneği
//
//	field := fields.NewBelongsTo("user", "User")
//	pagination := fields.NewRelationshipPagination(field)
//	// pagination.page = 1, pagination.perPage = 15, pagination.total = 0
//
// # Varsayılan Değerler
//
//   - Sayfa numarası: 1
//   - Sayfa boyutu: 15
//   - Toplam kayıt: 0
//
// Daha fazla bilgi için docs/Relationships.md dosyasına bakın.
func NewRelationshipPagination(field RelationshipField) *RelationshipPaginationImpl {
	return &RelationshipPaginationImpl{
		field:   field,
		page:    1,
		perPage: 15,
		total:   0,
	}
}

// ApplyPagination, ilişki sorgusuna sayfalama uygular.
//
// Bu metod, sayfa numarası ve sayfa boyutunu doğrular, geçersiz değerleri düzeltir
// ve sayfalama parametrelerini ayarlar. Sayfa boyutu maksimum 100 ile sınırlandırılır.
//
// # Parametreler
//
//   - ctx: Context nesnesi (iptal ve timeout yönetimi için)
//   - page: Sayfa numarası (minimum 1, geçersiz değerler 1'e ayarlanır)
//   - perPage: Sayfa başına kayıt sayısı (minimum 1, maksimum 100)
//
// # Dönüş Değerleri
//
//   - []interface{}: Sayfalanmış kayıt listesi
//   - error: Hata durumunda hata nesnesi, başarılı ise nil
//
// # Doğrulama Kuralları
//
//   - page < 1 ise otomatik olarak 1'e ayarlanır
//   - perPage < 1 ise otomatik olarak 15'e ayarlanır
//   - perPage > 100 ise otomatik olarak 100'e sınırlanır
//
// # Kullanım Örneği
//
//	pagination := fields.NewRelationshipPagination(field)
//	results, err := pagination.ApplyPagination(ctx, 2, 20)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	// results içinde 2. sayfanın 20 kaydı bulunur
//
// # Önemli Notlar
//
//   - Şu anki implementasyon boş slice döndürür (placeholder)
//   - Gerçek implementasyonda LIMIT ve OFFSET ile veritabanı sorgusu yapılır
//   - Sayfa boyutu güvenlik nedeniyle 100 ile sınırlandırılmıştır
//
// Daha fazla bilgi için docs/Relationships.md dosyasına bakın.
func (rp *RelationshipPaginationImpl) ApplyPagination(ctx context.Context, page int, perPage int) ([]interface{}, error) {
	if page < 1 {
		page = 1
	}

	if perPage < 1 {
		perPage = 15
	}

	// Limit per-page to maximum allowed
	if perPage > 100 {
		perPage = 100
	}

	rp.page = page
	rp.perPage = perPage

	// In a real implementation, this would query the database with LIMIT and OFFSET
	// For now, return empty slice
	return []interface{}{}, nil
}

// GetPageInfo, sayfalama metadata'sını döndürür.
//
// Bu metod, mevcut sayfalama durumu hakkında detaylı bilgi içeren bir map döndürür.
// Toplam sayfa sayısı, mevcut sayfa, kayıt aralıkları gibi bilgiler içerir.
//
// # Dönüş Değeri
//
// Map içeriği:
//   - current_page (int): Mevcut sayfa numarası
//   - per_page (int): Sayfa başına kayıt sayısı
//   - total (int64): Toplam kayıt sayısı
//   - total_pages (int64): Toplam sayfa sayısı
//   - from (int): Başlangıç kayıt indeksi (0-based)
//   - to (int): Bitiş kayıt indeksi (0-based)
//
// # Kullanım Örneği
//
//	pagination := fields.NewRelationshipPagination(field)
//	pagination.SetTotal(100)
//	pagination.ApplyPagination(ctx, 2, 15)
//
//	info := pagination.GetPageInfo()
//	fmt.Printf("Sayfa %d/%d\n", info["current_page"], info["total_pages"])
//	fmt.Printf("Kayıt %d-%d / %d\n", info["from"], info["to"], info["total"])
//	// Çıktı: Sayfa 2/7
//	// Çıktı: Kayıt 15-30 / 100
//
// # Hesaplama Detayları
//
//   - total_pages = ceil(total / perPage)
//   - from = (current_page - 1) * perPage
//   - to = current_page * perPage
//
// # Önemli Notlar
//
//   - perPage 0 ise total_pages 0 olarak döner (sıfıra bölme hatası önlenir)
//   - from ve to değerleri 0-based indeks kullanır
//   - SetTotal() ile toplam kayıt sayısı ayarlanmalıdır
//
// Daha fazla bilgi için docs/Relationships.md dosyasına bakın.
func (rp *RelationshipPaginationImpl) GetPageInfo() map[string]interface{} {
	totalPages := int64(0)
	if rp.perPage > 0 {
		totalPages = (rp.total + int64(rp.perPage) - 1) / int64(rp.perPage)
	}

	return map[string]interface{}{
		"current_page": rp.page,
		"per_page":     rp.perPage,
		"total":        rp.total,
		"total_pages":  totalPages,
		"from":         (rp.page - 1) * rp.perPage,
		"to":           rp.page * rp.perPage,
	}
}

// SetTotal, toplam kayıt sayısını ayarlar.
//
// Bu metod, sayfalama hesaplamaları için gerekli olan toplam kayıt sayısını belirler.
// GetPageInfo() metodunun doğru sonuçlar döndürebilmesi için bu metod çağrılmalıdır.
//
// # Parametreler
//
//   - total: Toplam kayıt sayısı (int64)
//
// # Kullanım Örneği
//
//	pagination := fields.NewRelationshipPagination(field)
//	pagination.SetTotal(150) // 150 kayıt var
//	pagination.ApplyPagination(ctx, 1, 20)
//
//	info := pagination.GetPageInfo()
//	// info["total"] = 150
//	// info["total_pages"] = 8 (ceil(150/20))
//
// # Önemli Notlar
//
//   - Bu metod ApplyPagination() çağrılmadan önce veya sonra çağrılabilir
//   - Toplam kayıt sayısı genellikle veritabanı COUNT sorgusu ile elde edilir
//   - Negatif değerler kabul edilir ancak mantıksal olarak 0 veya pozitif olmalıdır
//
// Daha fazla bilgi için docs/Relationships.md dosyasına bakın.
func (rp *RelationshipPaginationImpl) SetTotal(total int64) {
	rp.total = total
}

// GetPage, mevcut sayfa numarasını döndürür.
//
// Bu metod, ApplyPagination() ile ayarlanan veya varsayılan sayfa numarasını döndürür.
//
// # Dönüş Değeri
//
//   - int: Mevcut sayfa numarası (minimum 1)
//
// # Kullanım Örneği
//
//	pagination := fields.NewRelationshipPagination(field)
//	fmt.Println(pagination.GetPage()) // Çıktı: 1 (varsayılan)
//
//	pagination.ApplyPagination(ctx, 3, 15)
//	fmt.Println(pagination.GetPage()) // Çıktı: 3
//
// # Önemli Notlar
//
//   - Varsayılan değer 1'dir (NewRelationshipPagination ile oluşturulduğunda)
//   - ApplyPagination() çağrıldıktan sonra güncellenir
//
// Daha fazla bilgi için docs/Relationships.md dosyasına bakın.
func (rp *RelationshipPaginationImpl) GetPage() int {
	return rp.page
}

// GetPerPage, sayfa başına kayıt sayısını döndürür.
//
// Bu metod, ApplyPagination() ile ayarlanan veya varsayılan sayfa boyutunu döndürür.
//
// # Dönüş Değeri
//
//   - int: Sayfa başına kayıt sayısı (minimum 1, maksimum 100)
//
// # Kullanım Örneği
//
//	pagination := fields.NewRelationshipPagination(field)
//	fmt.Println(pagination.GetPerPage()) // Çıktı: 15 (varsayılan)
//
//	pagination.ApplyPagination(ctx, 1, 25)
//	fmt.Println(pagination.GetPerPage()) // Çıktı: 25
//
//	pagination.ApplyPagination(ctx, 1, 150)
//	fmt.Println(pagination.GetPerPage()) // Çıktı: 100 (maksimum sınır)
//
// # Önemli Notlar
//
//   - Varsayılan değer 15'tir (NewRelationshipPagination ile oluşturulduğunda)
//   - ApplyPagination() çağrıldıktan sonra güncellenir
//   - Değer otomatik olarak 1-100 aralığında sınırlandırılır
//
// Daha fazla bilgi için docs/Relationships.md dosyasına bakın.
func (rp *RelationshipPaginationImpl) GetPerPage() int {
	return rp.perPage
}

// GetTotal, toplam kayıt sayısını döndürür.
//
// Bu metod, SetTotal() ile ayarlanan toplam kayıt sayısını döndürür.
// Sayfalama metadata hesaplamaları için kullanılır.
//
// # Dönüş Değeri
//
//   - int64: Toplam kayıt sayısı
//
// # Kullanım Örneği
//
//	pagination := fields.NewRelationshipPagination(field)
//	fmt.Println(pagination.GetTotal()) // Çıktı: 0 (varsayılan)
//
//	pagination.SetTotal(250)
//	fmt.Println(pagination.GetTotal()) // Çıktı: 250
//
//	// Toplam sayfa sayısını hesapla
//	totalPages := (pagination.GetTotal() + int64(pagination.GetPerPage()) - 1) / int64(pagination.GetPerPage())
//	fmt.Println(totalPages) // Çıktı: 17 (ceil(250/15))
//
// # Önemli Notlar
//
//   - Varsayılan değer 0'dır (NewRelationshipPagination ile oluşturulduğunda)
//   - SetTotal() ile güncellenir
//   - GetPageInfo() bu değeri kullanarak total_pages hesaplar
//
// Daha fazla bilgi için docs/Relationships.md dosyasına bakın.
func (rp *RelationshipPaginationImpl) GetTotal() int64 {
	return rp.total
}
