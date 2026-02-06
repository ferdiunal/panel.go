package handler

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/data"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/query"
	"github.com/gofiber/fiber/v2"
)

// HandleResourceIndex lists resources with pagination, sorting, and filtering.
// It handles GET requests to /api/resource/:resource endpoint.
//
// Supports two query formats:
//   - Nested: users[search]=query, users[sort][id]=asc, users[filters][status][eq]=active
//   - Legacy: search=query, sort_column=id, sort_direction=asc
func HandleResourceIndex(h *FieldHandler, c *context.Context) error {
	ctx := c.Resource()

	if h.Policy != nil && !h.Policy.ViewAny(c) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Determine elements to use: Context > Handler Defaults
	var elements []fields.Element
	if ctx != nil && len(ctx.Elements) > 0 {
		elements = ctx.Elements
	} else {
		elements = h.Elements
	}

	if len(elements) == 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "No fields defined for this resource",
		})
	}

	// Parse Query Request using new parser
	resourceName := c.Params("resource")
	queryParams := query.ParseResourceQuery(c.Ctx, resourceName)

	// Convert query.Sort to data.Sort
	var sorts []data.Sort
	for _, s := range queryParams.Sorts {
		sorts = append(sorts, data.Sort{
			Column:    s.Column,
			Direction: s.Direction,
		})
	}

	// Apply defaults from Resource if no sorts provided
	if len(sorts) == 0 {
		if h.Resource != nil {
			for _, s := range h.Resource.GetSortable() {
				sorts = append(sorts, data.Sort{
					Column:    s.Column,
					Direction: s.Direction,
				})
			}
		}
		// Absolute fallback
		if len(sorts) == 0 {
			sorts = append(sorts, data.Sort{
				Column:    "created_at",
				Direction: "desc",
			})
		}
	}

	// Build QueryRequest
	req := data.QueryRequest{
		Page:    queryParams.Page,
		PerPage: queryParams.PerPage,
		Sorts:   sorts,
		Search:  queryParams.Search,
		Filters: queryParams.Filters,
	}

	// Fetch Data
	result, err := h.Provider.Index(c, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Map items to resources with fields extracted
	resources := make([]map[string]interface{}, 0)

	for _, item := range result.Items {
		res := h.resolveResourceFields(c.Ctx, c.Resource(), item, elements)
		// Inject per-item policy
		policy := map[string]bool{
			"view":   h.Policy == nil || h.Policy.View(c, item),
			"update": h.Policy == nil || h.Policy.Update(c, item),
			"delete": h.Policy == nil || h.Policy.Delete(c, item),
		}
		res["policy"] = policy
		resources = append(resources, res)
	}

	// Generate headers for frontend table order
	headers := make([]map[string]interface{}, 0)

	for _, element := range elements {
		if !element.IsVisible(c.Resource()) {
			continue
		}

		ctxStr := element.GetContext()
		serialized := element.JsonSerialize()

		// logic for headers (Index/List)
		if ctxStr != fields.HIDE_ON_LIST &&
			ctxStr != fields.ONLY_ON_CREATE &&
			ctxStr != fields.ONLY_ON_UPDATE &&
			ctxStr != fields.ONLY_ON_FORM &&
			ctxStr != fields.ONLY_ON_DETAIL {
			headers = append(headers, serialized)
		}
	}

	return c.JSON(fiber.Map{
		"data": resources,
		"meta": fiber.Map{
			"current_page": result.Page,
			"per_page":     result.PerPage,
			"total":        result.Total,
			"dialog_type":  h.DialogType,
			"title":        h.Resource.Title(),
			"headers":      headers,
			"policy": fiber.Map{
				"create":   h.Policy == nil || h.Policy.Create(c),
				"view_any": h.Policy == nil || h.Policy.ViewAny(c),
				"update":   h.Policy == nil || h.Policy.Update(c, nil),
				"delete":   h.Policy == nil || h.Policy.Delete(c, nil),
			},
		},
	})
}
