package handler

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/gofiber/fiber/v2"
)

// ResolveDependenciesRequest represents the request body for dependency resolution
type ResolveDependenciesRequest struct {
	FormData      map[string]interface{} `json:"formData"`
	Context       string                 `json:"context"`
	ChangedFields []string               `json:"changedFields"`
	ResourceID    interface{}            `json:"resourceId"`
}

// HandleResolveDependencies handles field dependency resolution
// It resolves field dependencies based on changed fields and form data
//
// Route: POST /resources/:resource/fields/resolve-dependencies
// Method: POST
//
// Request body:
//   - formData: Current form data
//   - context: "create" or "update"
//   - changedFields: List of fields that changed
//   - resourceId: Resource ID (for update context)
//
// Response:
//   - fields: Map of field keys to field updates
func HandleResolveDependencies(h *FieldHandler, c *context.Context) error {
	// Parse request body
	var req ResolveDependenciesRequest
	if err := c.Ctx.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate context
	if req.Context != "create" && req.Context != "update" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid context. Must be 'create' or 'update'",
		})
	}

	// Convert elements to Schema fields
	schemaFields := make([]*fields.Schema, 0, len(h.Elements))
	for _, element := range h.Elements {
		if schema, ok := element.(*fields.Schema); ok {
			schemaFields = append(schemaFields, schema)
		}
	}

	// Create dependency resolver
	resolver := fields.NewDependencyResolver(schemaFields, req.Context)

	// Check for circular dependencies
	if err := resolver.DetectCircularDependencies(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Resolve dependencies
	updates, err := resolver.ResolveDependencies(req.FormData, req.ChangedFields, c.Ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to resolve dependencies",
		})
	}

	// Return field updates
	return c.JSON(fiber.Map{
		"fields": updates,
	})
}
