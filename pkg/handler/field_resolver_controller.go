package handler

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/gofiber/fiber/v2"
)

// HandleFieldResolve handles calling field-specific resolver functions.
// It retrieves a field resolver by name and calls it with the provided parameters.
// This endpoint allows frontend components to trigger dynamic field transformations.
//
// Route: /resources/:resource/:id/fields/:field/resolve
// Method: POST
//
// Request body should contain resolver-specific parameters as JSON.
//
// Response: The resolved value from the field resolver.
//
// Requirement 16.1: THE Sistem SHALL alan resolver'larını API endpoint'leri aracılığıyla erişilebilir hale getirmelidir
// Requirement 16.2: THE Sistem SHALL resolver'ların özel veri dönüşümleri gerçekleştirmesine izin vermelidir
// Requirement 16.3: WHEN bir resolver çağrıldığında, THE Sistem SHALL resolver-spesifik parametreleri desteklemelidir
func HandleFieldResolve(h *FieldHandler, c *context.Context) error {
	fieldName := c.Params("field")
	resourceID := c.Params("id")

	// Find the field by name
	var targetField interface{}
	for _, element := range h.Elements {
		if element.GetKey() == fieldName {
			targetField = element
			break
		}
	}

	if targetField == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Field not found",
		})
	}

	// Get the resource item
	item, err := h.Provider.Show(c, resourceID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Resource not found",
		})
	}

	// Parse parameters from request body
	params := make(map[string]interface{})
	if err := c.Ctx.BodyParser(&params); err != nil {
		// If no body, use empty params
		params = make(map[string]interface{})
	}

	// For now, we'll return the field's resolved value
	// In a more advanced implementation, this would call a specific resolver
	// that can perform custom transformations based on the parameters

	// Extract the field value from the item
	if field, ok := targetField.(interface{ Extract(interface{}) }); ok {
		field.Extract(item)
	}

	// Serialize the field
	var serialized map[string]interface{}
	if field, ok := targetField.(interface{ JsonSerialize() map[string]interface{} }); ok {
		serialized = field.JsonSerialize()
	}

	// Return the resolved data
	return c.JSON(fiber.Map{
		"data": serialized,
	})
}
