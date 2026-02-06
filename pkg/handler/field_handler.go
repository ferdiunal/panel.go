// Package handler provides HTTP request handlers for the panel API.
// This package contains the core FieldHandler struct and shared helper methods
// used by all controller functions to handle resource operations.
package handler

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/ferdiunal/panel.go/pkg/auth"
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/core"
	"github.com/ferdiunal/panel.go/pkg/data"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/resource"
	"github.com/ferdiunal/panel.go/pkg/widget"
	"github.com/iancoleman/strcase"

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
	var provider data.DataProvider
	if repo := res.Repository(db); repo != nil {
		provider = repo
	} else {
		provider = data.NewGormDataProvider(db, res.Model())
	}

	cards := res.Cards()

	var searchCols []string
	for _, field := range res.Fields() {
		if field.IsSearchable() {
			searchCols = append(searchCols, field.GetKey())
		}
	}
	provider.SetSearchColumns(searchCols)

	var withRels []string
	withRels = append(withRels, res.With()...)

	for _, element := range res.Fields() {
		if relField, ok := fields.IsRelationshipField(element); ok {
			if relField.GetLoadingStrategy() == fields.EAGER_LOADING {
				// Use the field key as the relationship name for GORM Preload
				// Ideally this should match the struct field name
				// We convert to CamelCase because GORM usually expects struct field names
				key := strcase.ToCamel(relField.GetKey())

				// Check if already exists to avoid duplicates
				exists := false
				for _, existing := range withRels {
					if existing == key {
						exists = true
						break
					}
				}
				if !exists {
					withRels = append(withRels, key)
				}
			}
		}
	}
	provider.SetWith(withRels)

	return &FieldHandler{
		DB:          db,
		Provider:    provider,
		Elements:    res.Fields(),
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
	// Use Lens.Query to get the base query
	lensQuery := lens.Query(db)

	// Initialize Provider with Lens Query
	provider := data.NewGormDataProvider(lensQuery, res.Model())

	var searchCols []string
	for _, field := range lens.Fields() {
		if field.IsSearchable() {
			searchCols = append(searchCols, field.GetKey())
		}
	}
	provider.SetSearchColumns(searchCols)
	// Lens query likely encapsulates necessary preloads, but we could also add them if Lens interface had With()

	return &FieldHandler{
		Provider: provider,
		Elements: lens.Fields(),
		Title:    lens.Name(),
		Resource: res,
	}
}

// Index method is now in resource_index_controller.go as HandleResourceIndex
func (h *FieldHandler) Index(c *context.Context) error {
	return HandleResourceIndex(h, c)
}

// Show method is now in resource_show_controller.go as HandleResourceShow
func (h *FieldHandler) Show(c *context.Context) error {
	return HandleResourceShow(h, c)
}

// Edit method is now in resource_edit_controller.go as HandleResourceEdit
func (h *FieldHandler) Edit(c *context.Context) error {
	return HandleResourceEdit(h, c)
}

// Detail method is now in resource_detail_controller.go as HandleResourceDetail
func (h *FieldHandler) Detail(c *context.Context) error {
	return HandleResourceDetail(h, c)
}

// resolveResourceFields extracts field data from an item using the provided elements
func (h *FieldHandler) resolveResourceFields(c *fiber.Ctx, ctx *core.ResourceContext, item interface{}, elements []fields.Element) map[string]interface{} {
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

		// Resolve options
		h.ResolveFieldOptions(element, serialized, item)

		// Apply callback if exists
		if callback := element.GetResolveCallback(); callback != nil {
			if val, ok := serialized["data"]; ok {
				serialized["data"] = callback(val, item, c)
			}
		}

		resourceData[serialized["key"].(string)] = serialized
	}
	return resourceData
}

// ResolveFieldOptions resolves dynamic options for a field element.
// This handles AutoOptions for relationships and executes manual option callbacks.
func (h *FieldHandler) ResolveFieldOptions(element fields.Element, serialized map[string]interface{}, item interface{}) {
	props, _ := serialized["props"].(map[string]interface{})
	if props == nil {
		props = make(map[string]interface{})
		serialized["props"] = props
	}

	// Handle AutoOptions via Config (works even if element is *Schema due to fluent API)
	config := element.GetAutoOptionsConfig()
	fmt.Printf("[DEBUG] ResolveFieldOptions - Key: %s, View: %s, Enabled: %v\n", element.GetKey(), element.GetView(), config.Enabled)
	if config.Enabled {
		if _, hasOpts := props["options"]; !hasOpts {
			table, _ := props["related_resource"].(string)
			display := config.DisplayField
			fmt.Printf("[DEBUG] AutoOptions Table: %s, Display: %s\n", table, display)

			if table != "" && display != "" {
				var results []map[string]interface{}
				view := element.GetView()

				if view == "has-one-field" {
					fk, _ := props["foreign_key"].(string)
					fmt.Printf("[DEBUG] HasOne Query - Table: %s, FK: %s\n", table, fk)
					if h.DB != nil && fk != "" {
						query := h.DB.Table(table).Select("id, " + display)

						var itemID interface{}
						if item != nil {
							val := reflect.ValueOf(item)
							if val.Kind() == reflect.Ptr {
								val = val.Elem()
							}
							if val.Kind() == reflect.Struct {
								// Try to find ID or Id field
								idField := val.FieldByName("ID")
								if !idField.IsValid() {
									idField = val.FieldByName("Id")
								}
								if idField.IsValid() {
									itemID = idField.Interface()
								}
							}
						}

						if itemID != nil {
							query = query.Where(fk+" IS NULL OR "+fk+" = ?", itemID)
						} else {
							query = query.Where(fk + " IS NULL OR " + fk + " = 0")
						}
						query.Find(&results)
					}
				} else if view == "belongs-to-field" || view == "belongs-to-many-field" {
					fmt.Printf("[DEBUG] BelongsTo Query - Table: %s\n", table)
					if h.DB != nil {
						h.DB.Table(table).Select("id, " + display).Find(&results)
					}
				}

				fmt.Printf("[DEBUG] Query Result Count: %d\n", len(results))
				opts := make(map[string]string)
				for _, r := range results {
					if val, ok := r[display]; ok {
						opts[fmt.Sprint(r["id"])] = fmt.Sprint(val)
					}
				}
				props["options"] = opts
			}
		}
	}

	// Resolve dynamic options if present (callback)
	if optsFunc, ok := props["options"].(func() map[string]string); ok {
		props["options"] = optsFunc()
	} else if optsFunc, ok := props["options"].(func() map[string]interface{}); ok {
		props["options"] = optsFunc()
	}
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
				normalizedKey := key
				if strings.HasSuffix(key, "[]") {
					normalizedKey = strings.TrimSuffix(key, "[]")
				}

				var isFileType bool
				for _, el := range h.Elements {
					if el.GetKey() == normalizedKey {
						if el.JsonSerialize()["type"] == fields.TYPE_FILE ||
							el.JsonSerialize()["type"] == fields.TYPE_VIDEO ||
							el.JsonSerialize()["type"] == fields.TYPE_AUDIO {
							isFileType = true
						}
						break
					}
				}

				if !isFileType {
					if strings.HasSuffix(key, "[]") || len(values) > 1 {
						body[normalizedKey] = values
					} else {
						body[normalizedKey] = values[0]
					}
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

		// Handle missing BelongsToMany fields in multipart/form-data
		// If a BelongsToMany field is missing from the request, it implies the user unchecked all options.
		// We set it to an empty slice so GormDataProvider clears the relationships.
		for _, el := range h.Elements {
			if _, ok := fields.IsRelationshipField(el); ok {
				// We specifically check for BelongsToMany by view or type if interface allows
				// Using reflection or type check if possible, or View/Type string
				if el.JsonSerialize()["view"] == "belongs-to-many-field" {
					key := el.GetKey()
					if _, exists := body[key]; !exists {
						body[key] = []interface{}{}
					}
				}
			}
		}
	}

	// Apply ModifyCallback for all fields
	for _, el := range h.Elements {
		if val, ok := body[el.GetKey()]; ok {
			if callback := el.GetModifyCallback(); callback != nil {
				body[el.GetKey()] = callback(val, c.Ctx)
			}
		}
	}

	return body, nil
}

// Store method is now in resource_store_controller.go as HandleResourceStore
func (h *FieldHandler) Store(c *context.Context) error {
	return HandleResourceStore(h, c)
}

// Update method is now in resource_update_controller.go as HandleResourceUpdate
func (h *FieldHandler) Update(c *context.Context) error {
	return HandleResourceUpdate(h, c)
}

// Destroy method is now in resource_destroy_controller.go as HandleResourceDestroy
func (h *FieldHandler) Destroy(c *context.Context) error {
	return HandleResourceDestroy(h, c)
}

// Create method is now in resource_create_controller.go as HandleResourceCreate
func (h *FieldHandler) Create(c *context.Context) error {
	return HandleResourceCreate(h, c)
}

// ListCards method is now in card_controller.go as HandleCardList
func (h *FieldHandler) ListCards(c *context.Context) error {
	return HandleCardList(h, c)
}

// GetCard method is now in card_detail_controller.go as HandleCardDetail
func (h *FieldHandler) GetCard(c *context.Context) error {
	return HandleCardDetail(h, c)
}

// List method is now in field_controller.go as HandleFieldList
func (h *FieldHandler) List(c *context.Context) error {
	return HandleFieldList(h, c)
}

// FieldContextMiddleware creates a ResourceContext and injects it into the request.
// This middleware initializes the context with resource metadata, visibility context,
// and field resolvers for use by downstream handlers.
//
// Parameters:
//   - resource: The resource being operated on
//   - lens: The optional lens (filtered view) being applied
//   - visibilityCtx: The context in which fields should be visible
//   - elements: The fields associated with this resource
//
// Returns:
//   - A Fiber handler that injects the ResourceContext
//
// Requirement 15.1: THE Sistem SHALL ResourceContext'i oluşturmak için middleware'i güncelle
// Requirement 15.4: WHEN context oluşturulduğunda, THE Sistem SHALL tüm gerekli kaynak bilgisini başlatmalıdır
func FieldContextMiddleware(resource interface{}, lens interface{}, visibilityCtx core.VisibilityContext, elements []fields.Element) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Create ResourceContext with full visibility and metadata
		ctx := core.NewResourceContextWithVisibility(
			c,
			resource,
			lens,
			visibilityCtx,
			nil, // Item will be set by handlers if needed
			nil, // User will be set by handlers if needed
			elements,
		)
		c.Locals(core.ResourceContextKey, ctx)
		return c.Next()
	}
}

// FieldContextMiddlewareWithItem creates a ResourceContext with an item and injects it into the request.
// This is used when operating on a specific resource instance.
//
// Parameters:
//   - resource: The resource being operated on
//   - lens: The optional lens (filtered view) being applied
//   - visibilityCtx: The context in which fields should be visible
//   - item: The specific resource instance being operated on
//   - user: The user performing the operation
//   - elements: The fields associated with this resource
//
// Returns:
//   - A Fiber handler that injects the ResourceContext
//
// Requirement 15.1: THE Sistem SHALL ResourceContext'i oluşturmak için middleware'i güncelle
// Requirement 15.4: WHEN context oluşturulduğunda, THE Sistem SHALL tüm gerekli kaynak bilgisini başlatmalıdır
func FieldContextMiddlewareWithItem(resource interface{}, lens interface{}, visibilityCtx core.VisibilityContext, item interface{}, user interface{}, elements []fields.Element) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Create ResourceContext with full visibility and metadata
		ctx := core.NewResourceContextWithVisibility(
			c,
			resource,
			lens,
			visibilityCtx,
			item,
			user,
			elements,
		)
		c.Locals(core.ResourceContextKey, ctx)
		return c.Next()
	}
}
