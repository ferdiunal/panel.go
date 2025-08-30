package handler

import (
	"fmt"

	"github.com/a-h/templ"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"panel.go/internal/ent"
)

type Options struct {
	Ent   *ent.Client
	Store *session.Store
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
