package handler

import (
	"fmt"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"panel.go/internal/resource/user_resource"
	"panel.go/internal/service"
)

type Services struct {
	AuthService *service.AuthService
}

type Options struct {
	Store   *session.Store
	Service *Services
}

type HandlerFunc func(c *fiber.Ctx) error

func View(c *fiber.Ctx, title string, componentFunc func(csrfToken string, title string, locale string) templ.Component) error {
	csrfToken := ""
	locale := "tr"
	if token := c.Locals("csrf"); token != nil {
		csrfToken = token.(string)
	}

	title = fmt.Sprintf("%s | %s", title, "Panel")

	c.Set("Content-Type", "text/html")
	component := componentFunc(csrfToken, title, locale)
	return component.Render(c.Context(), c.Response().BodyWriter())
}

func WithAuthView(c *fiber.Ctx, title string, componentFunc func(csrfToken string, title string, locale string, user *user_resource.UserResource) templ.Component) error {
	csrfToken := ""
	locale := "tr"
	if token := c.Locals("csrf"); token != nil {
		csrfToken = token.(string)
	}

	title = fmt.Sprintf("%s | %s", title, "Panel")

	c.Set("Content-Type", "text/html")
	user := c.Locals("user").(*user_resource.UserResource)
	fmt.Println(user)
	component := componentFunc(csrfToken, title, locale, user)
	return component.Render(c.Context(), c.Response().BodyWriter())
}
