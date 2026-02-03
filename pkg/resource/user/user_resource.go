package user

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	appContext "github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/data"
	"github.com/ferdiunal/panel.go/pkg/data/orm"
	domainUser "github.com/ferdiunal/panel.go/pkg/domain/user"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/permission"
	"github.com/ferdiunal/panel.go/pkg/resource"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

// UserPolicy, Kullanıcı yönetimi için yetkilendirme kurallarını belirler.
// auth.Policy arayüzünü implemente eder.
type UserPolicy struct{}

// ViewAny, kullanıcının listeyi görüntüleyip görüntüleyemeyeceğini belirler.
func (p UserPolicy) ViewAny(ctx *appContext.Context) bool {
	return ctx.HasPermission("users.view_any")
}

// View, kullanıcının belirli bir kaydı görüntüleyip görüntüleyemeyeceğini belirler.
func (p UserPolicy) View(ctx *appContext.Context, model interface{}) bool {
	return ctx.HasPermission("users.view")
}

// Create, yeni bir kullanıcı oluşturma yetkisini kontrol eder.
func (p UserPolicy) Create(ctx *appContext.Context) bool {
	return ctx.HasPermission("users.create")
}

// Update, mevcut bir kullanıcıyı güncelleme yetkisini kontrol eder.
func (p UserPolicy) Update(ctx *appContext.Context, model interface{}) bool {
	return ctx.HasPermission("users.update")
}

// Delete, bir kullanıcıyı silme yetkisini kontrol eder.
// Kendini silmeyi engeller.
func (p UserPolicy) Delete(ctx *appContext.Context, model interface{}) bool {
	// Genel yetki kontrolü (model nil ise)
	if model == nil {
		return true
	}

	userModel, ok := model.(*domainUser.User)
	if !ok {
		return false
	}

	authUser := ctx.User()
	if authUser == nil {
		return false
	}

	// Kendini silmeyi engelle
	if userModel.ID == authUser.ID {
		return false
	}

	return true
}

// UserResourceWrapper embeds GenericResource to override Repository
type UserResourceWrapper struct {
	resource.Base
}

func (r UserResourceWrapper) Repository(db *gorm.DB) data.DataProvider {
	return orm.NewUserRepository(db)
}

// GetUserResource, Kullanıcı kaynağının (Resource) konfigürasyonunu döner.
// Alan tanımları, görünüm ayarları ve diğer meta verileri içerir.
func GetUserResource() resource.Resource {
	return UserResourceWrapper{
		Base: resource.Base{
			DataModel:  &domainUser.User{},
			Label:      "Users",
			Identifier: "users",
			IconName:   "users",
			GroupName:  "System",
			DialogType: resource.DialogTypeSheet,
			FieldsVal: []fields.Element{
				fields.ID(),
				fields.Image("Image").
					StoreAs(func(c *fiber.Ctx, file *multipart.FileHeader) (string, error) {
						storageUrl := "/storage/"
						storagePath := "./storage/public"
						ext := filepath.Ext(file.Filename)
						filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
						localPath := filepath.Join(storagePath, filename)

						// Klasörün varlığından emin ol
						_ = os.MkdirAll(storagePath, 0755)

						if err := c.SaveFile(file, localPath); err != nil {
							return "", err
						}
						// Public URL döndür
						return fmt.Sprintf("%s/%s", storageUrl, filename), nil
					}).
					Resolve(func(value interface{}) interface{} {
						if value == "" {
							return nil
						}
						// Tam URL oluştur
						// Not: Gerçek ortamda APP_URL env değişkeni kullanılmalı
						return value
					}),
				fields.Text("Name").Placeholder("Enter your name"),
				fields.Email("Email").Placeholder("Enter your email"),
				fields.Select("Role").
					Placeholder("Select a role").
					Resolve(func(value interface{}) interface{} {
						if value == "" {
							return nil
						}
						return value
					}).
					Options(func() map[string]string {
						mgr := permission.GetInstance()
						if mgr == nil {
							// Fallback or panic depending on requirement.
							// Since permissions MUST be loaded, empty map or default is safe.
							return map[string]string{}
						}
						roles := mgr.GetRoles()
						options := make(map[string]string)
						for _, role := range roles {
							// Assuming role name is also the label for now, or title case it
							options[role] = cases.Title(language.English).String(role) // You might want to use strcase.ToTitle(role)
						}
						return options
					}()),
				fields.Password("Password").
					Nullable().
					OnlyOnForm().
					Placeholder("Enter your password").
					Modify(func(value interface{}) interface{} {
						// Şifre hashleme işlemi
						str, ok := value.(string)
						if !ok || str == "" {
							return value
						}
						hash, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.DefaultCost)
						if err != nil {
							return value
						}
						return string(hash)
					}),
			},
			Sortable: []resource.Sortable{
				{
					Column:    "updated_at",
					Direction: "desc",
				},
				{
					Column:    "created_at",
					Direction: "desc",
				},
			},
			PolicyVal: UserPolicy{},
		},
	}
}
