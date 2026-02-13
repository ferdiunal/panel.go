package handler

import (
	"fmt"
	"strconv"
	"strings"

	stdContext "context"

	"github.com/ferdiunal/panel.go/internal/auth"
	"github.com/ferdiunal/panel.go/internal/context"
	"github.com/ferdiunal/panel.go/internal/data"
	"github.com/ferdiunal/panel.go/internal/fields"
	"github.com/ferdiunal/panel.go/internal/resource"
	"github.com/ferdiunal/panel.go/internal/widget"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type FieldHandler struct {
	DB          *gorm.DB
	Provider    data.DataProvider
	Elements    []fields.Element
	Policy      auth.Policy
	StoragePath string
	StorageURL  string
	Resource    resource.Resource
	Cards       []widget.Card
	Title       string
	DialogType  resource.DialogType
}

func NewFieldHandler(provider data.DataProvider) *FieldHandler {
	return &FieldHandler{
		Provider: provider,
	}
}

func (h *FieldHandler) SetElements(elements []fields.Element) {
	h.Elements = elements
}

// NewResourceHandler creates a FieldHandler from a Resource definition
func NewResourceHandler(db *gorm.DB, res resource.Resource, storagePath, storageURL string) *FieldHandler {
	elements := fields.CloneElements(res.Fields())

	var provider data.DataProvider
	if repo := res.Repository(db); repo != nil {
		provider = repo
	} else {
		provider = data.NewGormDataProvider(db, res.Model())
	}

	cards := res.Cards()

	var searchCols []string
	for _, field := range elements {
		if field.IsSearchable() {
			searchCols = append(searchCols, field.GetKey())
		}
	}
	provider.SetSearchColumns(searchCols)
	provider.SetWith(res.With())

	return &FieldHandler{
		DB:          db,
		Provider:    provider,
		Elements:    elements,
		Policy:      res.Policy(),
		Resource:    res,
		StoragePath: storagePath,
		StorageURL:  storageURL,
		Cards:       cards,
		Title:       res.Title(),
		DialogType:  res.GetDialogType(),
	}
}

// NewLensHandler creates a FieldHandler from a Resource and Lens definition
func NewLensHandler(db *gorm.DB, res resource.Resource, lens resource.Lens) *FieldHandler {
	lensElements := fields.CloneElements(lens.Fields())

	// Use Lens.Query to get the base query
	lensQuery := lens.Query(db)

	// Initialize Provider with Lens Query
	provider := data.NewGormDataProvider(lensQuery, res.Model())

	var searchCols []string
	for _, field := range lensElements {
		if field.IsSearchable() {
			searchCols = append(searchCols, field.GetKey())
		}
	}
	provider.SetSearchColumns(searchCols)
	// Lens query likely encapsulates necessary preloads, but we could also add them if Lens interface had With()

	return &FieldHandler{
		Provider: provider,
		Elements: lensElements,
		Title:    lens.Name(),
		Resource: res,
	}
}

func (h *FieldHandler) Index(c *context.Context) error {
	ctx := context.FromFiber(c.Ctx)

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

	// Calculate Context if strictly needed for other things, but here we just need elements
	if ctx == nil {
		// Create a temporary context if needed, or just proceed with elements
		// We might need to ensure ctx is not nil for some edge cases,
		// but for Index logic below we use 'elements' variable.
		// However, let's create a minimal context for consistency if we wanted to pass it down.
		// For now, adhering to the logic below.
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
	result, err := h.Provider.Index(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Map items to resources with fields extracted
	resources := make([]map[string]interface{}, 0)

	for _, item := range result.Items {
		res := h.resolveResourceFields(c.Context(), item, elements)
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
		if !element.IsVisible(c.Context()) {
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

func (h *FieldHandler) Show(c *context.Context) error {
	id := c.Params("id")
	item, err := h.Provider.Show(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Not found"})
	}

	if h.Policy != nil && !h.Policy.View(c, item) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Determine elements to use
	// For Show, we generally use the handler's elements, but context overrides could happen if middleware sets them
	// For now, simpler to use h.Elements as middleware might be Index specific or general.
	// But let's check context just in case consistent with Index.
	ctx := context.FromFiber(c.Ctx)
	var elements []fields.Element
	if ctx != nil && len(ctx.Elements) > 0 {
		elements = ctx.Elements
	} else {
		elements = h.Elements
	}

	return c.JSON(fiber.Map{
		"data": h.resolveResourceFields(c.Context(), item, elements),
		"meta": fiber.Map{
			"title": h.Resource.Title(),
			"policy": fiber.Map{
				"view":   h.Policy == nil || h.Policy.View(c, item),
				"update": h.Policy == nil || h.Policy.Update(c, item),
				"delete": h.Policy == nil || h.Policy.Delete(c, item),
			},
		},
	})
}

func (h *FieldHandler) Edit(c *context.Context) error {
	id := c.Params("id")
	item, err := h.Provider.Show(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Not found"})
	}

	if h.Policy != nil && !h.Policy.Update(c, item) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Filter for Update
	var updateElements []fields.Element
	for _, element := range h.Elements {
		if !element.IsVisible(c.Context()) {
			continue
		}

		ctxStr := element.GetContext()
		if ctxStr != fields.HIDE_ON_UPDATE &&
			ctxStr != fields.ONLY_ON_LIST &&
			ctxStr != fields.ONLY_ON_DETAIL &&
			ctxStr != fields.ONLY_ON_CREATE {
			updateElements = append(updateElements, element)
		}
	}

	// Resolve fields with values
	// Note: resolveResourceFields returns a map keyed by field key.
	// But for the Form, we might prefer an ordered list or just the map.
	// The frontend ResourceForm takes "fields" (array of definitions) and "initialData" (map of values).
	// If we return the resolved fields, they contain BOTH definition and value (in 'data' key).
	// So we can return a list of these fields.

	resolvedMap := h.resolveResourceFields(c.Context(), item, updateElements)

	// Convert map to ordered slice based on h.Elements order
	var orderedFields []map[string]interface{}
	for _, element := range updateElements {
		if val, ok := resolvedMap[element.GetKey()]; ok {
			// Cast val to map
			if fieldMap, ok := val.(map[string]interface{}); ok {
				orderedFields = append(orderedFields, fieldMap)
			}
		}
	}

	return c.JSON(fiber.Map{
		"fields": orderedFields,
		"meta": fiber.Map{
			"title": h.Title,
		},
	})
}

func (h *FieldHandler) Detail(c *context.Context) error {
	id := c.Params("id")
	item, err := h.Provider.Show(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Not found"})
	}

	if h.Policy != nil && !h.Policy.View(c, item) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Filter for Detail
	var detailElements []fields.Element
	for _, element := range h.Elements {
		if !element.IsVisible(c.Context()) {
			continue
		}

		ctxStr := element.GetContext()

		// Skip if explicitly hidden on detail or restricted to other contexts
		if ctxStr == fields.HIDE_ON_DETAIL ||
			ctxStr == fields.ONLY_ON_LIST ||
			ctxStr == fields.ONLY_ON_FORM ||
			ctxStr == fields.HIDE_ON_UPDATE {
			continue
		}
		detailElements = append(detailElements, element)
	}

	// Resolve fields with values
	resolvedMap := h.resolveResourceFields(c.Context(), item, detailElements)

	// Convert map to ordered slice based on h.Elements order (preserving filtered list order)
	var orderedFields []map[string]interface{}
	// Iterate over detailElements to preserve order
	for _, element := range detailElements {
		if val, ok := resolvedMap[element.GetKey()]; ok {
			if fieldMap, ok := val.(map[string]interface{}); ok {
				orderedFields = append(orderedFields, fieldMap)
			}
		}
	}

	return c.JSON(fiber.Map{
		"fields": orderedFields,
		"meta": fiber.Map{
			"title": h.Title,
		},
	})
}

// resolveResourceFields extracts field data from an item using the provided elements
func (h *FieldHandler) resolveResourceFields(ctx stdContext.Context, item interface{}, elements []fields.Element) map[string]interface{} {
	resourceData := make(map[string]interface{})
	for _, element := range elements {
		if !element.IsVisible(ctx) {
			continue
		}
		// Clone logic or direct access warning applies here too as noted in previous Index implementation.
		// For now, we proceed with direct extraction which mutates element state.
		// In a real high-concurrency scenario, elements should be cloned or Extract should return value.
		element.Extract(item)
		serialized := element.JsonSerialize()

		// Apply callback if exists
		if callback := element.GetResolveCallback(); callback != nil {
			if val, ok := serialized["data"]; ok {
				serialized["data"] = callback(val)
			}
		}

		resourceData[serialized["key"].(string)] = serialized
	}
	return resourceData
}

// parseBody extracts data from request
func (h *FieldHandler) parseBody(c *context.Context) (map[string]interface{}, error) {
	var body = make(map[string]interface{})

	// Check content type
	ctype := c.Ctx.Get("Content-Type")
	if !strings.Contains(ctype, "multipart/form-data") && !strings.Contains(ctype, "application/json") {
		// Try to parse json body manually if not standard
		if err := c.Ctx.BodyParser(&body); err != nil {
			return nil, err
		}
		return body, nil
	}

	// Handle JSON Body
	if err := c.Ctx.BodyParser(&body); err == nil {
		// If success, we have data. But if it's multipart, we might have partial data in body?
		// Fiber BodyParser handles multipart form fields too.
	}

	// Handle Form Data (Multipart)
	if form, err := c.Ctx.MultipartForm(); err == nil {
		for key, values := range form.Value {
			if len(values) > 0 {
				var isFileType bool
				for _, el := range h.Elements {
					if el.GetKey() == key {
						if el.JsonSerialize()["type"] == fields.TYPE_FILE ||
							el.JsonSerialize()["type"] == fields.TYPE_VIDEO ||
							el.JsonSerialize()["type"] == fields.TYPE_AUDIO {
							isFileType = true
						}
						break
					}
				}

				if !isFileType {
					body[key] = values[0]
				}
			}
		}
		for key, files := range form.File {
			if len(files) > 0 {
				file := files[0]
				var path string

				// Check for matching element and callback
				var callback fields.StorageCallbackFunc
				for _, el := range h.Elements {
					if el.GetKey() == key {
						callback = el.GetStorageCallback()
						break
					}
				}

				if callback != nil {
					// User defined storage logic (Assuming it takes *fiber.Ctx)
					path, err = callback(c.Ctx, file)
					if err != nil {
						return nil, err
					}
				} else {
					// Use Resource StoreHandler (Takes *context.Context)
					path, err = h.Resource.StoreHandler(c, file, h.StoragePath, h.StorageURL)
					if err != nil {
						return nil, err
					}
				}

				body[key] = path
			}
		}
	}

	// Apply ModifyCallback for all fields
	for _, el := range h.Elements {
		if val, ok := body[el.GetKey()]; ok {
			if callback := el.GetModifyCallback(); callback != nil {
				body[el.GetKey()] = callback(val)
			}
		}
	}

	return body, nil
}

func (h *FieldHandler) Store(c *context.Context) error {
	data, err := h.parseBody(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if h.Policy != nil && !h.Policy.Create(c) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Unauthorized"})
	}

	result, err := h.Provider.Create(c.Context(), data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": h.resolveResourceFields(c.Context(), result, h.Elements)})
}

func (h *FieldHandler) Update(c *context.Context) error {
	id := c.Params("id")
	data, err := h.parseBody(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Fetch existing to check policy
	item, err := h.Provider.Show(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Not found"})
	}

	if h.Policy != nil && !h.Policy.Update(c, item) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Unauthorized"})
	}

	result, err := h.Provider.Update(c.Context(), id, data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": h.resolveResourceFields(c.Context(), result, h.Elements)})
}

func (h *FieldHandler) Destroy(c *context.Context) error {
	id := c.Params("id")

	// Fetch for Policy Check
	item, err := h.Provider.Show(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Not found"})
	}

	if h.Policy != nil && !h.Policy.Delete(c, item) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if err := h.Provider.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Deleted successfully"})
}

func (h *FieldHandler) ListCards(c *context.Context) error {
	resp := make([]map[string]interface{}, 0)
	for i, w := range h.Cards {
		// Serialize base properties (component, width, etc.)
		serialized := w.JsonSerialize()
		serialized["index"] = i
		serialized["name"] = w.Name()
		serialized["component"] = w.Component()
		serialized["width"] = w.Width()

		// Resolve data
		data, err := w.Resolve(c, h.DB)
		if err != nil {
			fmt.Printf("Error resolving card %s: %v\n", w.Name(), err)
			serialized["error"] = err.Error()
		} else {
			// Assign resolved data to "data" key
			serialized["data"] = data
		}

		resp = append(resp, serialized)
	}

	return c.JSON(fiber.Map{
		"data": resp,
	})
}

func (h *FieldHandler) GetCard(c *context.Context) error {
	index, err := c.ParamsInt("index")
	if err != nil || index < 0 || index >= len(h.Cards) {
		return c.Status(404).JSON(fiber.Map{"error": "Card not found"})
	}

	w := h.Cards[index]

	// Resolve data
	data, err := w.Resolve(c, h.DB)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"data": data,
	})
}

func (h *FieldHandler) List(c *context.Context) error {
	ctx := context.FromFiber(c.Ctx) // Use c.Ctx to get the *fiber.Ctx for FromFiber
	if ctx == nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Field context not found",
		})
	}

	response := make([]map[string]interface{}, 0)

	for _, element := range fields.CloneElements(ctx.Elements) {
		// Resolve value if resource is present
		if ctx.Resource != nil {
			element.Extract(ctx.Resource)
		}
		response = append(response, element.JsonSerialize())
	}

	return c.JSON(response)
}

// Middleware to inject context (simplified for now)
func FieldContextMiddleware(resource interface{}, elements []fields.Element) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := context.NewResourceContext(c, resource, elements)
		c.Locals(context.ResourceContextKey, ctx)
		return c.Next()
	}
}
