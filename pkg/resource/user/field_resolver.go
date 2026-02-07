package user

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/permission"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Bu yapı, kullanıcı kaynağı (resource) için alan çözümleyicisini temsil eder.
//
// UserFieldResolver, panel.go uygulamasında kullanıcı yönetimi arayüzünde gösterilecek
// form alanlarını dinamik olarak tanımlamaktan sorumludur. Bu yapı, fields.Resolver
// arayüzünü uygulayarak, kullanıcı oluşturma ve düzenleme formlarında hangi alanların
// gösterileceğini, nasıl doğrulanacağını ve nasıl saklanacağını belirler.
//
// Kullanım Senaryosu:
// - Yönetici panelinde kullanıcı oluşturma/düzenleme formlarının dinamik oluşturulması
// - Farklı kullanıcı rollerine göre alan görünürlüğünün kontrol edilmesi
// - Dosya yükleme (profil resmi) işlemlerinin yönetilmesi
// - Rol seçimi için dinamik seçeneklerin sağlanması
//
// Önemli Notlar:
// - Bu yapı durum taşımaz (stateless), sadece alan tanımlarını sağlar
// - Tüm alan konfigürasyonları method chaining kullanılarak yapılır
// - Dosya depolama işlemleri "./storage/public" dizinine yapılır
// - Roller dinamik olarak permission yöneticisinden alınır
type UserFieldResolver struct{}

// Bu metod, kullanıcı kaynağı için tüm form alanlarını tanımlar ve döner.
//
// Parametreler:
//   - ctx (*context.Context): İstek bağlamı, kullanıcı bilgileri ve izinleri içerir
//
// Dönüş Değeri:
//   - []fields.Element: Tanımlanan form alanlarının dilimi
//
// Alanlar:
//   1. ID: Arama yapılabilir, otomatik olarak oluşturulan benzersiz tanımlayıcı
//   2. Image: Profil resmi, dosya yükleme desteği ile
//   3. Name: Kullanıcı adı, arama yapılabilir
//   4. Email: E-posta adresi, arama yapılabilir
//   5. Role: Kullanıcı rolü, dinamik seçenekler ile
//   6. Password: Şifre, sadece form görünümünde, boş bırakılabilir
//
// Kullanım Örneği:
//   resolver := &UserFieldResolver{}
//   fields := resolver.ResolveFields(ctx)
//   // fields dilimi form oluşturmak için kullanılır
//
// Önemli Notlar:
//   - Resim dosyaları UnixNano timestamp ile adlandırılır (çakışma riski düşük)
//   - Boş resim değerleri nil olarak döndürülür (UI'de gösterilmez)
//   - Şifre alanı sadece form görünümünde gösterilir, liste görünümünde gizlenir
//   - Roller permission yöneticisinden dinamik olarak alınır
//   - Depolama başarısız olursa hata döndürülür
func (r *UserFieldResolver) ResolveFields(ctx *context.Context) []fields.Element {
	return []fields.Element{
		// Bu alan, kullanıcının benzersiz tanımlayıcısını temsil eder.
		// Özellikler:
		//   - Otomatik olarak veritabanı tarafından oluşturulur
		//   - Arama yapılabilir (Searchable), kullanıcıları ID ile bulabilirsiniz
		//   - Liste görünümünde gösterilir
		//   - Düzenleme formunda salt okunur olarak gösterilir
		fields.ID().Searchable(),

		// Bu alan, kullanıcının profil resmini yönetir.
		// Özellikler:
		//   - Dosya yükleme desteği (multipart/form-data)
		//   - Yüklenen dosyalar "./storage/public" dizinine kaydedilir
		//   - Dosya adı UnixNano timestamp + orijinal uzantı ile oluşturulur
		//   - Boş değerler UI'de gösterilmez (nil döndürülür)
		//
		// Depolama Stratejisi:
		//   - Dosya adı: [UnixNano timestamp][orijinal uzantı]
		//   - Örnek: 1707329574123456789.jpg
		//   - Public URL: /storage/[dosya adı]
		//   - Çakışma riski: Çok düşük (nanosecond hassasiyeti)
		//
		// Hata Yönetimi:
		//   - Klasör oluştulamıyorsa sessizce devam eder
		//   - Dosya kaydedilemiyorsa hata döndürülür
		//   - Boş dosya başlığı işlenmez
		fields.Image("Image").
			StoreAs(func(c *fiber.Ctx, file *multipart.FileHeader) (string, error) {
				storageUrl := "/storage/"
				storagePath := "./storage/public"
				ext := filepath.Ext(file.Filename)
				filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
				localPath := filepath.Join(storagePath, filename)

				// Klasörün varlığından emin ol, yoksa oluştur (0755 izinleri ile)
				_ = os.MkdirAll(storagePath, 0755)

				if err := c.SaveFile(file, localPath); err != nil {
					return "", err
				}
				// Public URL döndür (frontend'de erişilebilir)
				return fmt.Sprintf("%s/%s", storageUrl, filename), nil
			}).
			Resolve(func(value any, item any, c *fiber.Ctx) any {
				// Boş string değerleri nil olarak döndür (UI'de gösterilmez)
				if value == "" {
					return nil
				}
				return value
			}),

		// Bu alan, kullanıcının adını temsil eder.
		// Özellikler:
		//   - Metin giriş alanı
		//   - Arama yapılabilir (Searchable), kullanıcıları ad ile bulabilirsiniz
		//   - Placeholder: "Enter your name"
		//   - Zorunlu alan (nullable değildir)
		//   - Form ve liste görünümünde gösterilir
		fields.Text("Name").Searchable().Placeholder("Enter your name"),

		// Bu alan, kullanıcının e-posta adresini temsil eder.
		// Özellikler:
		//   - E-posta doğrulaması yapılır (email format)
		//   - Arama yapılabilir (Searchable), kullanıcıları e-posta ile bulabilirsiniz
		//   - Placeholder: "Enter your email"
		//   - Zorunlu alan (nullable değildir)
		//   - Benzersiz olması önerilir (veritabanı seviyesinde kontrol edilmelidir)
		fields.Email("Email").Searchable().Placeholder("Enter your email"),

		// Bu alan, kullanıcının rolünü temsil eder.
		// Özellikler:
		//   - Seçim alanı (dropdown/select)
		//   - Seçenekler dinamik olarak permission yöneticisinden alınır
		//   - Boş değerler nil olarak döndürülür (UI'de gösterilmez)
		//   - Placeholder: "Select a role"
		//   - Rol adları başlık durumuna dönüştürülür (Title Case)
		//
		// Dinamik Seçenekler:
		//   - Permission yöneticisinden tüm mevcut roller alınır
		//   - Her rol için key=role_name, value=Title(role_name) şeklinde seçenek oluşturulur
		//   - Örnek: "admin" -> "Admin", "editor" -> "Editor"
		//   - Permission yöneticisi null ise boş seçenekler döndürülür
		//
		// Hata Yönetimi:
		//   - Permission yöneticisi başlatılmamışsa boş seçenekler döndürülür
		//   - Rol listesi alınamıyorsa boş seçenekler döndürülür
		fields.Select("Role").
			Placeholder("Select a role").
			Resolve(func(value any, item any, c *fiber.Ctx) any {
				// Boş string değerleri nil olarak döndür (UI'de gösterilmez)
				if value == "" {
					return nil
				}
				return value
			}).
			Options(func() map[string]string {
				// Permission yöneticisinin singleton instance'ını al
				mgr := permission.GetInstance()
				if mgr == nil {
					return map[string]string{}
				}
				// Tüm mevcut rolleri al
				roles := mgr.GetRoles()
				options := make(map[string]string)
				// Her rol için seçenek oluştur (key=role, value=Title(role))
				for _, role := range roles {
					options[role] = cases.Title(language.English).String(role)
				}
				return options
			}()),

		// Bu alan, kullanıcının şifresini temsil eder.
		// Özellikler:
		//   - Şifre giriş alanı (maskelenmiş görünüm)
		//   - Nullable: Boş bırakılabilir (güncelleme sırasında şifre değiştirilmeyebilir)
		//   - OnlyOnForm: Sadece form görünümünde gösterilir, liste görünümünde gizlenir
		//   - Placeholder: "Enter your password"
		//   - Güvenlik nedeniyle veritabanında hash'lenmiş olarak saklanmalıdır
		//
		// Kullanım Senaryoları:
		//   - Yeni kullanıcı oluşturma: Şifre zorunludur
		//   - Kullanıcı düzenleme: Şifre boş bırakılabilir (değiştirilmez)
		//   - Şifre değiştirilmek isteniyorsa yeni şifre girilir
		//
		// Önemli Notlar:
		//   - Şifre asla liste görünümünde gösterilmez
		//   - Şifre asla veritabanında düz metin olarak saklanmamalıdır
		//   - Şifre hash'leme işlemi backend'de yapılmalıdır
		//   - Boş şifre değeri güncelleme sırasında mevcut şifreyi korur
		fields.Password("Password").
			Nullable().
			OnlyOnForm().
			Placeholder("Enter your password"),
	}
}
