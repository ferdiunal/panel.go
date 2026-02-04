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
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// UserFieldResolver, kullanıcı alanlarını çözer
type UserFieldResolver struct{}

// ResolveFields, kullanıcı alanlarını döner
func (r *UserFieldResolver) ResolveFields(ctx *context.Context) []fields.Element {
	return []fields.Element{
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
			Resolve(func(value any, c *fiber.Ctx) any {
				if value == "" {
					return nil
				}
				return value
			}),
		fields.Text("Name").Placeholder("Enter your name"),
		fields.Email("Email").Placeholder("Enter your email"),
		fields.Select("Role").
			Placeholder("Select a role").
			Resolve(func(value any, c *fiber.Ctx) any {
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
		fields.Password("Password").
			Nullable().
			OnlyOnForm().
			Placeholder("Enter your password").
			Modify(func(value any, c *fiber.Ctx) any {
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
	}
}
