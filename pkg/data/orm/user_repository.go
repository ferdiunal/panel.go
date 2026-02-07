// Bu paket, ORM (Object-Relational Mapping) katmanında kullanıcı veri erişim işlemlerini yönetir.
// GORM kütüphanesi üzerinden veritabanı operasyonlarını gerçekleştirir ve domain katmanı ile
// veri sağlayıcı arasında köprü görevi yapar.
package orm

import (
	"context"

	pkgContext "github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/data"
	"github.com/ferdiunal/panel.go/pkg/domain/user"
	"gorm.io/gorm"
)

// Bu yapı, kullanıcı veri erişim işlemlerini yönetmek için repository pattern'ini implement eder.
// GormDataProvider'ı embed ederek generic CRUD operasyonlarını miras alır ve GORM veritabanı
// bağlantısını doğrudan tutar. Hem domain-specific (CreateUser, FindByID) hem de generic
// (Create, Update, Delete) metodları destekler.
//
// Kullanım Senaryoları:
// - Kullanıcı oluşturma, güncelleme, silme işlemleri
// - Kullanıcı arama (ID veya Email ile)
// - Kullanıcı sayısını alma
// - Sayfalanmış kullanıcı listesi alma
//
// Önemli Notlar:
// - Context desteği ile timeout ve cancellation işlemleri yapılabilir
// - Hem typed (CreateUser) hem generic (Create) metodlar mevcuttur
// - GormDataProvider embed edildiği için generic sorgu yetenekleri otomatik olarak kullanılabilir
type UserRepository struct {
	// GormDataProvider, generic CRUD operasyonları için temel işlevselliği sağlar
	*data.GormDataProvider
	// db, GORM veritabanı bağlantısı pointer'ı. Domain-specific operasyonlar için doğrudan kullanılır
	db *gorm.DB
}

// Bu fonksiyon, UserRepository'nin yeni bir örneğini oluşturur ve başlatır.
// Verilen GORM veritabanı bağlantısını kullanarak repository'yi yapılandırır.
//
// Parametreler:
// - db: GORM veritabanı bağlantısı pointer'ı
//
// Dönüş Değeri:
// - *UserRepository: Yapılandırılmış UserRepository pointer'ı
//
// Kullanım Örneği:
//   db := gorm.Open(sqlite.Open("test.db"))
//   userRepo := NewUserRepository(db)
//
// Önemli Notlar:
// - Verilen db pointer'ı nil olmamalıdır
// - GormDataProvider otomatik olarak user.User modeli ile başlatılır
// - Repository'nin tüm metodları bu fonksiyon tarafından oluşturulan örnek üzerinde çalışır
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		GormDataProvider: data.NewGormDataProvider(db, &user.User{}),
		db:               db,
	}
}

// Bu metod, domain katmanı tarafından kullanılan typed interface'i implement eder.
// Verilen User nesnesini veritabanına yeni bir kayıt olarak ekler.
// Context desteği ile timeout ve cancellation işlemleri yapılabilir.
//
// Parametreler:
// - ctx: Context nesnesi (timeout, cancellation ve deadline bilgisi içerir)
// - u: Oluşturulacak User nesnesinin pointer'ı
//
// Dönüş Değeri:
// - error: İşlem başarılı ise nil, aksi takdirde hata nesnesi
//
// Kullanım Örneği:
//   ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//   defer cancel()
//   newUser := &user.User{Email: "test@example.com", Name: "Test User"}
//   err := userRepo.CreateUser(ctx, newUser)
//   if err != nil {
//       log.Printf("Kullanıcı oluşturma hatası: %v", err)
//   }
//
// Önemli Notlar:
// - User nesnesinin gerekli alanları (Email, Name vb.) doldurulmuş olmalıdır
// - Veritabanı kısıtlamaları (unique constraints) ihlal edilirse hata döner
// - Context timeout'u aşılırsa işlem iptal edilir
// - Başarılı oluşturma sonrası User nesnesine ID atanır
func (r *UserRepository) CreateUser(ctx context.Context, u *user.User) error {
	return r.db.WithContext(ctx).Create(u).Error
}

// Bu metod, GormDataProvider'ın generic Create metodunu override eder.
// Verilen veri haritasını kullanarak yeni bir User kaydı oluşturur.
// Panel UI ve API tarafından kullanılan generic interface'i implement eder.
//
// Parametreler:
// - ctx: Panel context nesnesi (kullanıcı, izinler ve diğer bilgiler içerir)
// - data: Oluşturulacak kullanıcı verilerini içeren harita (örn: {"email": "test@example.com", "name": "Test"})
//
// Dönüş Değeri:
// - interface{}: Oluşturulan User nesnesini içeren interface{} (type assertion gerekli)
// - error: İşlem başarılı ise nil, aksi takdirde hata nesnesi
//
// Kullanım Örneği:
//   ctx := &pkgContext.Context{User: currentUser}
//   userData := map[string]interface{}{
//       "email": "newuser@example.com",
//       "name": "New User",
//   }
//   result, err := userRepo.Create(ctx, userData)
//   if err != nil {
//       return err
//   }
//   createdUser := result.(*user.User)
//
// Önemli Notlar:
// - GormDataProvider'ın Create metodunu çağırarak generic işlevselliği kullanır
// - Dönüş değeri interface{} olduğu için type assertion yapılması gerekir
// - Panel UI tarafından kullanılan generic CRUD operasyonları için tasarlanmıştır
// - Döndürür: - Oluşturulan User nesnesini içeren interface{} pointer'ı
func (r *UserRepository) Create(ctx *pkgContext.Context, data map[string]interface{}) (interface{}, error) {
	return r.GormDataProvider.Create(ctx, data)
}

// Bu metod, verilen ID'ye göre veritabanından User kaydını bulur.
// Context desteği ile timeout ve cancellation işlemleri yapılabilir.
//
// Parametreler:
// - ctx: Context nesnesi (timeout, cancellation ve deadline bilgisi içerir)
// - id: Aranacak kullanıcının ID'si (uint türü)
//
// Dönüş Değeri:
// - *user.User: Bulunan User nesnesinin pointer'ı
// - error: Kayıt bulunamadı ise gorm.ErrRecordNotFound, diğer hatalar için hata nesnesi
//
// Kullanım Örneği:
//   ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//   defer cancel()
//   user, err := userRepo.FindByID(ctx, 1)
//   if err != nil {
//       if errors.Is(err, gorm.ErrRecordNotFound) {
//           log.Println("Kullanıcı bulunamadı")
//       }
//       return
//   }
//   log.Printf("Kullanıcı: %s (%s)", user.Name, user.Email)
//
// Önemli Notlar:
// - ID'nin geçerli bir uint değeri olması gerekir
// - Kayıt bulunamadı ise gorm.ErrRecordNotFound hatası döner
// - Context timeout'u aşılırsa işlem iptal edilir
// - Veritabanı bağlantı hatası ise ilgili hata döner
func (r *UserRepository) FindByID(ctx context.Context, id uint) (*user.User, error) {
	var u user.User
	if err := r.db.WithContext(ctx).First(&u, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// Bu metod, verilen email adresine göre veritabanından User kaydını bulur.
// Context desteği ile timeout ve cancellation işlemleri yapılabilir.
// Genellikle kullanıcı girişi ve email doğrulama işlemlerinde kullanılır.
//
// Parametreler:
// - ctx: Context nesnesi (timeout, cancellation ve deadline bilgisi içerir)
// - email: Aranacak kullanıcının email adresi (string türü)
//
// Dönüş Değeri:
// - *user.User: Bulunan User nesnesinin pointer'ı
// - error: Kayıt bulunamadı ise gorm.ErrRecordNotFound, diğer hatalar için hata nesnesi
//
// Kullanım Örneği:
//   ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//   defer cancel()
//   user, err := userRepo.FindByEmail(ctx, "user@example.com")
//   if err != nil {
//       if errors.Is(err, gorm.ErrRecordNotFound) {
//           log.Println("Bu email adresine sahip kullanıcı bulunamadı")
//       }
//       return
//   }
//   log.Printf("Kullanıcı bulundu: %s (ID: %d)", user.Name, user.ID)
//
// Önemli Notlar:
// - Email adresi case-sensitive olarak aranır (veritabanı ayarlarına bağlı)
// - Kayıt bulunamadı ise gorm.ErrRecordNotFound hatası döner
// - Context timeout'u aşılırsa işlem iptal edilir
// - Email adresi unique constraint'i varsa en fazla bir kayıt döner
// - Girişi ve email doğrulama işlemlerinde sıkça kullanılır
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	var u user.User
	if err := r.db.WithContext(ctx).First(&u, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// Bu metod, domain katmanı tarafından kullanılan typed interface'i implement eder.
// Verilen User nesnesinin tüm alanlarını veritabanında günceller.
// Context desteği ile timeout ve cancellation işlemleri yapılabilir.
//
// Parametreler:
// - ctx: Context nesnesi (timeout, cancellation ve deadline bilgisi içerir)
// - u: Güncellenecek User nesnesinin pointer'ı (ID alanı set olmalıdır)
//
// Dönüş Değeri:
// - error: İşlem başarılı ise nil, aksi takdirde hata nesnesi
//
// Kullanım Örneği:
//   ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//   defer cancel()
//   user, _ := userRepo.FindByID(ctx, 1)
//   user.Name = "Updated Name"
//   user.Email = "newemail@example.com"
//   err := userRepo.UpdateUser(ctx, user)
//   if err != nil {
//       log.Printf("Kullanıcı güncelleme hatası: %v", err)
//   }
//
// Önemli Notlar:
// - User nesnesinin ID alanı set olmalıdır (hangi kaydın güncellenecek olduğunu belirtir)
// - Tüm alanlar güncellenir (zero değerler de dahil)
// - Veritabanı kısıtlamaları ihlal edilirse hata döner
// - Context timeout'u aşılırsa işlem iptal edilir
// - Başarılı güncelleme sonrası User nesnesinin alanları veritabanındaki değerlerle eşleşir
func (r *UserRepository) UpdateUser(ctx context.Context, u *user.User) error {
	return r.db.WithContext(ctx).Save(u).Error
}

// Bu metod, GormDataProvider'ın generic Update metodunu override eder.
// Verilen veri haritasını kullanarak User kaydını günceller.
// Panel UI ve API tarafından kullanılan generic interface'i implement eder.
//
// Parametreler:
// - ctx: Panel context nesnesi (kullanıcı, izinler ve diğer bilgiler içerir)
// - id: Güncellenecek kullanıcının ID'si (string türü)
// - data: Güncellenecek alanları içeren harita (örn: {"name": "Updated Name", "email": "new@example.com"})
//
// Dönüş Değeri:
// - interface{}: Güncellenen User nesnesini içeren interface{} (type assertion gerekli)
// - error: İşlem başarılı ise nil, aksi takdirde hata nesnesi
//
// Kullanım Örneği:
//   ctx := &pkgContext.Context{User: currentUser}
//   updateData := map[string]interface{}{
//       "name": "Updated Name",
//       "email": "updated@example.com",
//   }
//   result, err := userRepo.Update(ctx, "1", updateData)
//   if err != nil {
//       return err
//   }
//   updatedUser := result.(*user.User)
//
// Önemli Notlar:
// - GormDataProvider'ın Update metodunu çağırarak generic işlevselliği kullanır
// - Dönüş değeri interface{} olduğu için type assertion yapılması gerekir
// - Panel UI tarafından kullanılan generic CRUD operasyonları için tasarlanmıştır
// - Döndürür: - Güncellenen User nesnesini içeren interface{} pointer'ı
func (r *UserRepository) Update(ctx *pkgContext.Context, id string, data map[string]interface{}) (interface{}, error) {
	return r.GormDataProvider.Update(ctx, id, data)
}

// Bu metod, GormDataProvider'ın generic Delete metodunu override eder.
// Verilen ID'ye göre User kaydını veritabanından siler.
// Panel UI ve API tarafından kullanılan generic interface'i implement eder.
//
// Parametreler:
// - ctx: Panel context nesnesi (kullanıcı, izinler ve diğer bilgiler içerir)
// - id: Silinecek kullanıcının ID'si (string türü)
//
// Dönüş Değeri:
// - error: İşlem başarılı ise nil, aksi takdirde hata nesnesi
//
// Kullanım Örneği:
//   ctx := &pkgContext.Context{User: currentUser}
//   err := userRepo.Delete(ctx, "1")
//   if err != nil {
//       log.Printf("Kullanıcı silme hatası: %v", err)
//   }
//
// Önemli Notlar:
// - GormDataProvider'ın Delete metodunu çağırarak generic işlevselliği kullanır
// - Panel UI tarafından kullanılan generic CRUD operasyonları için tasarlanmıştır
// - Silme işlemi geri alınamaz (soft delete değildir)
// - Yabancı anahtar kısıtlamaları varsa hata döner
func (r *UserRepository) Delete(ctx *pkgContext.Context, id string) error {
	return r.GormDataProvider.Delete(ctx, id)
}

// Bu metod, domain katmanı tarafından kullanılan typed interface'i implement eder.
// Verilen ID'ye göre User kaydını veritabanından siler.
// Context desteği ile timeout ve cancellation işlemleri yapılabilir.
//
// Parametreler:
// - ctx: Context nesnesi (timeout, cancellation ve deadline bilgisi içerir)
// - id: Silinecek kullanıcının ID'si (uint türü)
//
// Dönüş Değeri:
// - error: İşlem başarılı ise nil, aksi takdirde hata nesnesi
//
// Kullanım Örneği:
//   ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//   defer cancel()
//   err := userRepo.DeleteUser(ctx, 1)
//   if err != nil {
//       log.Printf("Kullanıcı silme hatası: %v", err)
//   }
//
// Önemli Notlar:
// - ID'nin geçerli bir uint değeri olması gerekir
// - Silme işlemi geri alınamaz (soft delete değildir)
// - Context timeout'u aşılırsa işlem iptal edilir
// - Yabancı anahtar kısıtlamaları varsa hata döner
// - Veritabanı bağlantı hatası ise ilgili hata döner
func (r *UserRepository) DeleteUser(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&user.User{}, "id = ?", id).Error
}

// Bu metod, veritabanında kayıtlı olan toplam User sayısını döndürür.
// Context desteği ile timeout ve cancellation işlemleri yapılabilir.
// Genellikle sayfalama ve istatistik işlemlerinde kullanılır.
//
// Parametreler:
// - ctx: Context nesnesi (timeout, cancellation ve deadline bilgisi içerir)
//
// Dönüş Değeri:
// - int64: Veritabanında kayıtlı toplam User sayısı
// - error: İşlem başarılı ise nil, aksi takdirde hata nesnesi
//
// Kullanım Örneği:
//   ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//   defer cancel()
//   count, err := userRepo.Count(ctx)
//   if err != nil {
//       log.Printf("Kullanıcı sayısı alma hatası: %v", err)
//       return
//   }
//   log.Printf("Toplam kullanıcı sayısı: %d", count)
//
// Önemli Notlar:
// - Soft delete kullanılıyorsa, silinmiş kayıtlar sayılmaz
// - Context timeout'u aşılırsa işlem iptal edilir
// - Veritabanı bağlantı hatası ise ilgili hata döner
// - Büyük veri setlerinde performans etkileyebilir
// - Sayfalama ve istatistik işlemlerinde sıkça kullanılır
func (r *UserRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&user.User{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// Bu metod, GormDataProvider'ın generic Show metodunu override eder.
// Verilen ID'ye göre User kaydını veritabanından bulur ve döndürür.
// Panel UI ve API tarafından kullanılan generic interface'i implement eder.
//
// Parametreler:
// - ctx: Panel context nesnesi (kullanıcı, izinler ve diğer bilgiler içerir)
// - id: Gösterilecek kullanıcının ID'si (string türü)
//
// Dönüş Değeri:
// - interface{}: Bulunan User nesnesini içeren interface{} (type assertion gerekli)
// - error: Kayıt bulunamadı ise gorm.ErrRecordNotFound, diğer hatalar için hata nesnesi
//
// Kullanım Örneği:
//   ctx := &pkgContext.Context{User: currentUser}
//   result, err := userRepo.Show(ctx, "1")
//   if err != nil {
//       if errors.Is(err, gorm.ErrRecordNotFound) {
//           return "Kullanıcı bulunamadı"
//       }
//       return err
//   }
//   user := result.(*user.User)
//   log.Printf("Kullanıcı: %s (%s)", user.Name, user.Email)
//
// Önemli Notlar:
// - GormDataProvider'ın Show metodunu çağırarak generic işlevselliği kullanır
// - Dönüş değeri interface{} olduğu için type assertion yapılması gerekir
// - Panel UI tarafından kullanılan generic CRUD operasyonları için tasarlanmıştır
// - Döndürür: - Bulunan User nesnesini içeren interface{} pointer'ı
func (r *UserRepository) Show(ctx *pkgContext.Context, id string) (interface{}, error) {
	return r.GormDataProvider.Show(ctx, id)
}

// Bu metod, GormDataProvider'ın generic Index metodunu override eder.
// Verilen sorgu parametrelerine göre User kayıtlarının sayfalanmış listesini döndürür.
// Panel UI ve API tarafından kullanılan generic interface'i implement eder.
// Filtreleme, sıralama ve sayfalama işlemlerini destekler.
//
// Parametreler:
// - ctx: Panel context nesnesi (kullanıcı, izinler ve diğer bilgiler içerir)
// - req: Sorgu parametrelerini içeren QueryRequest nesnesi (sayfa, limit, sıralama, filtreler vb.)
//
// Dönüş Değeri:
// - *data.QueryResponse: Sorgu sonuçlarını içeren QueryResponse pointer'ı (User listesi, toplam sayı vb.)
// - error: İşlem başarılı ise nil, aksi takdirde hata nesnesi
//
// Kullanım Örneği:
//   ctx := &pkgContext.Context{User: currentUser}
//   req := data.QueryRequest{
//       Page: 1,
//       Limit: 10,
//       Sort: "name",
//       Order: "asc",
//   }
//   response, err := userRepo.Index(ctx, req)
//   if err != nil {
//       log.Printf("Kullanıcı listesi alma hatası: %v", err)
//       return
//   }
//   log.Printf("Toplam: %d, Sayfa: %d, Kayıtlar: %v", response.Total, response.Page, response.Data)
//
// Önemli Notlar:
// - GormDataProvider'ın Index metodunu çağırarak generic işlevselliği kullanır
// - Sayfalama, sıralama ve filtreleme işlemlerini destekler
// - Panel UI tarafından kullanılan generic CRUD operasyonları için tasarlanmıştır
// - Döndürür: - Sorgu sonuçlarını içeren QueryResponse pointer'ı
func (r *UserRepository) Index(ctx *pkgContext.Context, req data.QueryRequest) (*data.QueryResponse, error) {
	return r.GormDataProvider.Index(ctx, req)
}
