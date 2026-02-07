package panel

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/gofiber/fiber/v2"
)

// handlePages returns a list of all registered pages.
// GET /api/pages
func (p *PanelField) handlePages(c *context.Context) error {
	type PageItem struct {
		Slug        string `json:"slug"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
		Group       string `json:"group"`
		Order       int    `json:"order"`
		Visible     bool   `json:"visible"`
	}

	items := []PageItem{}
	for slug, pg := range p.pages {
		if !pg.Visible() || !pg.CanAccess(c) {
			continue
		}
		items = append(items, PageItem{
			Slug:        slug,
			Title:       pg.Title(),
			Description: pg.Description(),
			Icon:        pg.Icon(),
			Group:       pg.Group(),
			Order:       pg.NavigationOrder(),
			Visible:     pg.Visible(),
		})
	}

	return c.JSON(fiber.Map{
		"data": items,
	})
}

// handlePageDetail returns the details of a specific page, including widgets.
// GET /api/pages/:slug
func (p *PanelField) handlePageDetail(c *context.Context) error {
	slug := c.Params("slug")

	pg, ok := p.pages[slug]
	if !ok {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Page not found",
		})
	}

	// Check access
	if !pg.CanAccess(c) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied",
		})
	}

	// Prepare cards
	cards := []map[string]interface{}{}
	for _, card := range pg.Cards() {
		serialized := card.JsonSerialize()
		if data, err := card.Resolve(c, p.Db); err == nil {
			serialized["data"] = data
		} else {
			// On error, maybe send null or log it?
			// For now, let's just send what we have, frontend might show empty state if data is missing?
			// But frontend crashes if data is missing.
			// Let's ensure data key exists at least.
			serialized["data"] = nil
		}
		cards = append(cards, serialized)
	}

	// Prepare fields
	// Prepare fields
	var fieldsList []map[string]interface{}
	for _, f := range pg.Fields() {
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
		"title":       pg.Title(),
		"description": pg.Description(),
		"meta": fiber.Map{
			"cards":  cards,
			"fields": fieldsList,
		},
	})
}

// POST /api/pages/:slug
func (p *PanelField) handlePageSave(c *context.Context) error {
	slug := c.Params("slug")

	pg, ok := p.pages[slug]
	if !ok {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Page not found",
		})
	}

	// Check access
	if !pg.CanAccess(c) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied",
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
