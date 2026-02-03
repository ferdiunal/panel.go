package user

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	appContext "github.com/ferdiunal/panel.go/internal/context"
	"github.com/ferdiunal/panel.go/internal/data"
	"github.com/ferdiunal/panel.go/internal/data/orm"
	domainUser "github.com/ferdiunal/panel.go/internal/domain/user"
	"github.com/ferdiunal/panel.go/internal/fields"
	"github.com/ferdiunal/panel.go/internal/resource"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserPolicy, Kullanıcı yönetimi için yetkilendirme kurallarını belirler.
// auth.Policy arayüzünü implemente eder.
type UserPolicy struct{}

// ViewAny, kullanıcının listeyi görüntüleyip görüntüleyemeyeceğini belirler.
func (p UserPolicy) ViewAny(ctx *appContext.Context) bool {
	return true
}

// View, kullanıcının belirli bir kaydı görüntüleyip görüntüleyemeyeceğini belirler.
func (p UserPolicy) View(ctx *appContext.Context, model interface{}) bool {
	return true
}

// Create, yeni bir kullanıcı oluşturma yetkisini kontrol eder.
func (p UserPolicy) Create(ctx *appContext.Context) bool {
	return true
}

// Update, mevcut bir kullanıcıyı güncelleme yetkisini kontrol eder.
func (p UserPolicy) Update(ctx *appContext.Context, model interface{}) bool {
	return true
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
						storageUrl := "storage/demo"
						storagePath := "./storage/public/demo"
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
						return fmt.Sprintf("%s/%s", "http://localhost:4555", value)
					}),
				fields.Text("Name").Placeholder("Enter your name"),
				fields.Email("Email").Placeholder("Enter your email"),
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
