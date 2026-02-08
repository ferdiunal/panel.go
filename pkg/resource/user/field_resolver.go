package user

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/i18n"
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
	// Nil context kontrolü - resource initialization sırasında context olmayabilir
	if ctx == nil {
		ctx = &context.Context{} // Boş context oluştur
	}

	return []fields.Element{
		// ID alanı - i18n destekli
		fields.ID().Searchable(),

		// Profil resmi alanı - i18n destekli label
		fields.Image("Image").
			Label(i18n.Trans(ctx.Ctx, "resources.users.fields.image")).
			StoreAs(func(c *fiber.Ctx, file *multipart.FileHeader) (string, error) {
				storageUrl := "/storage/"
				storagePath := "./storage/public"
				ext := filepath.Ext(file.Filename)
				filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
				localPath := filepath.Join(storagePath, filename)

				_ = os.MkdirAll(storagePath, 0755)

				if err := c.SaveFile(file, localPath); err != nil {
					return "", err
				}
				return fmt.Sprintf("%s/%s", storageUrl, filename), nil
			}).
			Resolve(func(value any, item any, c *fiber.Ctx) any {
				if value == "" {
					return nil
				}
				return value
			}),

		// İsim alanı - i18n destekli label ve placeholder
		fields.Text("Name").
			Label(i18n.Trans(ctx.Ctx, "resources.users.fields.name")).
			Placeholder(i18n.Trans(ctx.Ctx, "resources.users.fields.name_placeholder")).
			Searchable(),

		// E-posta alanı - i18n destekli label ve placeholder
		fields.Email("Email").
			Label(i18n.Trans(ctx.Ctx, "resources.users.fields.email")).
			Placeholder(i18n.Trans(ctx.Ctx, "resources.users.fields.email_placeholder")).
			Searchable(),

		// Rol alanı - i18n destekli label ve placeholder
		fields.Select("Role").
			Label(i18n.Trans(ctx.Ctx, "resources.users.fields.role")).
			Placeholder(i18n.Trans(ctx.Ctx, "resources.users.fields.role_placeholder")).
			Resolve(func(value any, item any, c *fiber.Ctx) any {
				if value == "" {
					return nil
				}
				return value
			}).
			Options(func() map[string]string {
				mgr := permission.GetInstance()
				if mgr == nil {
					return map[string]string{}
				}
				roles := mgr.GetRoles()
				options := make(map[string]string)
				for _, role := range roles {
					options[role] = cases.Title(language.English).String(role)
				}
				return options
			}()),

		// Şifre alanı - i18n destekli label ve placeholder
		fields.Password("Password").
			Label(i18n.Trans(ctx.Ctx, "resources.users.fields.password")).
			Placeholder(i18n.Trans(ctx.Ctx, "resources.users.fields.password_placeholder")).
			Nullable().
			OnlyOnForm(),
	}
}
