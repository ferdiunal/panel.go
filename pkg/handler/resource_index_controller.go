package handler

import (
	"fmt"
	"strconv"

	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/data"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/resource"
	"github.com/gofiber/fiber/v2"
)

// HandleResourceIndex lists resources with pagination, sorting, and filtering.
// It handles GET requests to /api/resource/:resource endpoint.
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

	// Parse Query Request
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "10"))

	// Parse Sort
	var sorts []data.Sort

	// 1. Try parsing sort_column as a map (e.g. sort_column[id]=desc)
	// Fiber's QueryParser is needed for maps
	type SortQuery struct {
		SortColumn map[string]string `query:"sort_column"`
	}

	sq := new(SortQuery)
	// We only care if parsing succeeds and has data
	if err := c.QueryParser(sq); err == nil && len(sq.SortColumn) > 0 {
		for col, dir := range sq.SortColumn {
			sorts = append(sorts, data.Sort{
				Column:    col,
				Direction: dir,
			})
		}
	} else {
		// 2. Fallback to simple string check if map failed or empty
		// (Example: ?sort_column=id&sort_direction=desc - legacy support or simple single sort)
		sCol := c.Query("sort_column")
		if sCol != "" {
			sDir := c.Query("sort_direction")
			if sDir == "" {
				sDir = "asc"
			}
			sorts = append(sorts, data.Sort{
				Column:    sCol,
				Direction: sDir,
			})
		}
	}

	// 3. Defaults from Resource if no sorts provided
	if len(sorts) == 0 {
		fmt.Printf("DEBUG: No API sorts provided. Checking defaults.\n")
		var sortables []resource.Sortable
		if h.Resource != nil {
			sortables = h.Resource.GetSortable()
		}
		if len(sortables) > 0 {
			fmt.Printf("DEBUG: Found %d defaults from Resource: %+v\n", len(sortables), sortables)
			for _, s := range sortables {
				sorts = append(sorts, data.Sort{
					Column:    s.Column,
					Direction: s.Direction,
				})
			}
		} else {
			fmt.Printf("DEBUG: No defaults from Resource. Using absolute backup.\n")
			// Absolute backup
			sorts = append(sorts, data.Sort{
				Column:    "created_at",
				Direction: "desc",
			})
		}
	} else {
		fmt.Printf("DEBUG: API sorts provided: %+v\n", sorts)
	}

	fmt.Printf("DEBUG: Final Sorts: %+v\n", sorts)

	req := data.QueryRequest{
		Page:    page,
		PerPage: perPage,
		Sorts:   sorts,
		Search:  c.Query("search"),
		// Filters would be parsed here from query params or body
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
