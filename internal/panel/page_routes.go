package panel

import (
	"github.com/ferdiunal/panel.go/internal/context"
	"github.com/ferdiunal/panel.go/pkg/i18n"
	"github.com/gofiber/fiber/v2"
)

// handlePages returns a list of all registered pages.
// GET /api/pages
func (p *Panel) handlePages(c *context.Context) error {
	type PageItem struct {
		Slug  string `json:"slug"`
		Title string `json:"title"`
	}

	items := []PageItem{}
	for slug, pg := range p.pages {
		items = append(items, PageItem{
			Slug:  slug,
			Title: i18n.Trans(c.Ctx, pg.Title()),
		})
	}

	return c.JSON(fiber.Map{
		"data": items,
	})
}

// handlePageDetail returns the details of a specific page, including widgets.
// GET /api/pages/:slug
func (p *Panel) handlePageDetail(c *context.Context) error {
	slug := c.Params("slug")

	pg, ok := p.pages[slug]
	if !ok {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Page not found",
		})
	}

	// Prepare cards
	cards := []map[string]interface{}{}
	for _, c := range pg.Cards() {
		cards = append(cards, c.JsonSerialize())
	}

	// Prepare fields
	// Prepare fields
	var fieldsList []map[string]interface{}
	for _, f := range pg.Fields() {
		// Inject Context for i18n
		if setter, ok := f.(interface{ SetContextForI18n(*fiber.Ctx) }); ok {
			setter.SetContextForI18n(c.Ctx)
		}

		serialized := f.JsonSerialize()

		// Inject value for Settings Page
		if pg.Slug() == "settings" {
			if key, ok := serialized["key"].(string); ok {
				// Try Dynamic Values first
				if val, exists := p.Config.SettingsValues.Values[key]; exists {
					serialized["data"] = val
				} else {
					// Fallback to struct fields for defaults
					switch key {
					case "site_name":
						serialized["data"] = p.Config.SettingsValues.SiteName
					case "register":
						serialized["data"] = p.Config.SettingsValues.Register
					case "forgot_password":
						serialized["data"] = p.Config.SettingsValues.ForgotPassword
					}
				}
			}
		}

		fieldsList = append(fieldsList, serialized)
	}

	return c.JSON(fiber.Map{
		"slug":        pg.Slug(),
		"title":       i18n.Trans(c.Ctx, pg.Title()),
		"description": i18n.Trans(c.Ctx, pg.Description()),
		"meta": fiber.Map{
			"cards":  cards,
			"fields": fieldsList,
		},
	})
}

// POST /api/pages/:slug
func (p *Panel) handlePageSave(c *context.Context) error {
	slug := c.Params("slug")

	pg, ok := p.pages[slug]
	if !ok {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Page not found",
		})
	}

	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON",
		})
	}

	if err := pg.Save(c, p.Db, data); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Hot Reload Settings
	if slug == "settings" {
		_ = p.LoadSettings()
	}

	return c.JSON(fiber.Map{
		"message": "Settings saved",
	})
}

// RegisterPage registers a new page to the panel.
// This is already defined in app.go, so we don't need it here.
// But keeping the file focused on page handlers.
