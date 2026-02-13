// Bu paket, kullanıcı kaynağı (user resource) ile ilgili veri erişim katmanını sağlar.
// Kullanıcı oluşturma, güncelleme, silme ve sorgulama işlemlerini yönetir.
// Özellikle şifre hashleme ve hesap (account) oluşturma gibi özel işlemleri içerir.
package user

import (
	"fmt"
	"time"

	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/data"
	"github.com/ferdiunal/panel.go/pkg/data/orm"
	"github.com/ferdiunal/panel.go/pkg/domain/account"
	"github.com/ferdiunal/panel.go/pkg/domain/user"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Bu yapı, GORM veri sağlayıcısını genişleterek özel kullanıcı oluşturma mantığı ekler.
//
// UserDataProvider, standart GORM işlemlerine ek olarak:
// - Şifre hashleme (bcrypt kullanarak)
// - Otomatik hesap (Account) oluşturma
// - İşlem başarısızlığında geri alma (rollback)
// gibi özel işlemleri gerçekleştirir.
//
// Kullanım Senaryosu:
// Yeni bir kullanıcı kaydı yapılırken, kullanıcı bilgileri ve şifresi alınır.
// Şifre bcrypt ile hashlenir, kullanıcı veritabanına kaydedilir ve
// ardından ilgili hesap kaydı oluşturulur.
//
// Alanlar:
// - *data.GormDataProvider: Temel GORM veri sağlayıcısı (gömülü)
// - client: GORM DB instance'ı (account oluşturma için)
type UserDataProvider struct {
	*data.GormDataProvider
	client            *gorm.DB
	accountRepository *orm.AccountRepository
}

// Bu fonksiyon, yeni bir UserDataProvider örneği oluşturur ve başlatır.
//
// Parametreler:
// - client (*gorm.DB): Ent client instance'ı
//
// Dönüş Değeri:
// - *UserDataProvider: Yapılandırılmış UserDataProvider pointer'ı
//
// Kullanım Örneği:
//
//	db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
//	provider := NewUserDataProvider(db)
//	user, err := provider.Create(ctx, map[string]interface{}{
//	    "email": "user@example.com",
//	    "password": "securepassword123",
//	})
//
// Önemli Notlar:
// - GORM veri sağlayıcısı user.User modeli ile yapılandırılır
// - Client nil ise panic oluşabilir
func NewUserDataProvider(client *gorm.DB) *UserDataProvider {
	return &UserDataProvider{
		GormDataProvider:  data.NewGormDataProvider(client, &user.User{}),
		client:            client,
		accountRepository: orm.NewAccountRepository(client),
	}
}

// Bu metod, varsayılan oluşturma işlemini geçersiz kılarak özel kullanıcı oluşturma mantığı ekler.
//
// Metod, aşağıdaki adımları sırasıyla gerçekleştirir:
// 1. Gelen verilerden şifreyi alır ve doğrular
// 2. Şifreyi bcrypt algoritması ile hashler
// 3. Kullanıcı kaydını GORM ile veritabanına kaydeder
// 4. Kaydedilen kullanıcı için hesap (Account) kaydı oluşturur
// 5. Hesap oluşturma başarısız olursa, kullanıcı kaydını geri alır (rollback)
//
// Parametreler:
// - ctx (*context.Context): İstek bağlamı (context), işlem izleme ve iptal için kullanılır
// - data (map[string]interface{}): Kullanıcı bilgileri
//   - "password" (string): Kullanıcının şifresi (zorunlu)
//   - Diğer alanlar: email, name, vb. kullanıcı özellikleri
//
// Dönüş Değerleri:
// - interface{}: Başarılı olursa oluşturulan User nesnesi
// - error: Hata durumunda hata mesajı
//
// Hata Durumları:
// - "password is required": Şifre alanı boş veya eksik
// - bcrypt.GenerateFromPassword hatası: Şifre hashleme başarısız
// - GORM Create hatası: Kullanıcı kaydı başarısız
// - Account Create hatası: Hesap kaydı başarısız (kullanıcı kaydı geri alınır)
//
// Kullanım Örneği:
//
//	provider := NewUserDataProvider(client)
//	ctx := context.New()
//	user, err := provider.Create(ctx, map[string]interface{}{
//	    "email": "john@example.com",
//	    "name": "John Doe",
//	    "password": "MySecurePassword123!",
//	})
//	if err != nil {
//	    log.Printf("Kullanıcı oluşturma hatası: %v", err)
//	    return
//	}
//	log.Printf("Kullanıcı başarıyla oluşturuldu: %d", user.ID)
//
// Önemli Notlar:
// - Şifre hiçbir zaman veritabanında düz metin olarak kaydedilmez
// - Şifre hashleme için bcrypt.DefaultCost (12) kullanılır
// - Hesap oluşturma başarısız olursa, kullanıcı kaydı otomatik olarak silinir
// - ProviderID "credential" olarak sabitlenmiştir (yerel kimlik doğrulama)
// - AccountID boş string olarak ayarlanır (harici sağlayıcı kimliği yoktur)
// - Başarılı olursa, oluşturulan User nesnesi döndürülür
// - Başarısız olursa, nil ve hata mesajı döndürülür
func (p *UserDataProvider) Create(ctx *context.Context, data map[string]interface{}) (interface{}, error) {
	// Adım 1: Gelen verilerden şifreyi alır ve doğrular
	// Şifre string türünde ve boş olmayan bir değer olmalıdır
	password, ok := data["password"].(string)
	if !ok || password == "" {
		return nil, fmt.Errorf("password is required")
	}

	// Adım 2: Şifreyi bcrypt algoritması ile hashler
	// bcrypt.DefaultCost (12) kullanılarak güvenli bir hash oluşturulur
	// Bu işlem CPU yoğun olabilir ve birkaç saniye sürebilir
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Şifreyi verilerden kaldırır ve hashlenmişini ekler
	// Böylece düz metin şifre veritabanına kaydedilmez
	delete(data, "password")
	data["password"] = string(hashed)

	// Adım 3: Kullanıcı kaydını GORM veri sağlayıcısı ile veritabanına kaydeder
	// Bu işlem, User modelinin tüm alanlarını veritabanına ekler
	result, err := p.GormDataProvider.Create(ctx, data)
	if err != nil {
		return nil, err
	}

	// Döndürülen sonucu User türüne dönüştürür
	// Eğer dönüşüm başarısız olursa, sonuç olduğu gibi döndürülür
	user, ok := result.(*user.User)
	if !ok {
		return result, nil
	}

	// Adım 4: Oluşturulan kullanıcı için hesap (Account) kaydı oluşturur
	// GORM kullanarak Account oluşturulur
	stdCtx := ctx.Context()
	err = p.accountRepository.Create(stdCtx, &account.Account{
		ProviderID: "credential",
		Password:   string(hashed),
		UserID:     user.ID,
		AccountID:  nil,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	})

	if err != nil {
		// Adım 5: Hesap oluşturma başarısız olursa, kullanıcı kaydını geri alır (rollback)
		// Bu, veri tutarlılığını sağlar ve yetim kayıtların oluşmasını önler
		_ = p.GormDataProvider.Delete(ctx, fmt.Sprint(user.ID))
		return nil, err
	}

	// Başarılı olursa, oluşturulan User nesnesi döndürülür
	return user, nil
}
