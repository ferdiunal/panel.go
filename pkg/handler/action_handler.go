package handler

import (
	"fmt"
	"reflect"

	"github.com/ferdiunal/panel.go/pkg/action"
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/gofiber/fiber/v2"
)

// HandleActionList returns the list of available actions for a resource.
// It serializes action metadata including name, slug, icon, confirmation settings,
// visibility flags, and field definitions.
func HandleActionList(h *FieldHandler, c *context.Context) error {
	// Policy check - user must have view permission
	if h.Policy != nil && !h.Policy.ViewAny(c) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Get actions from resource
	actions := h.Resource.GetActions()

	// Serialize actions
	serialized := make([]map[string]interface{}, 0, len(actions))
	for _, act := range actions {
		// Check if action implements the new action.Action interface
		if newAction, ok := act.(action.Action); ok {
			fields := make([]map[string]interface{}, 0)
			for _, field := range newAction.GetFields() {
				fields = append(fields, field.JsonSerialize())
			}

			serialized = append(serialized, map[string]interface{}{
				"name":               newAction.GetName(),
				"slug":               newAction.GetSlug(),
				"icon":               newAction.GetIcon(),
				"confirmText":        newAction.GetConfirmText(),
				"confirmButtonText":  newAction.GetConfirmButtonText(),
				"cancelButtonText":   newAction.GetCancelButtonText(),
				"destructive":        newAction.IsDestructive(),
				"onlyOnIndex":        newAction.OnlyOnIndex(),
				"onlyOnDetail":       newAction.OnlyOnDetail(),
				"showInline":         newAction.ShowInline(),
				"fields":             fields,
			})
		}
	}

	return c.JSON(fiber.Map{
		"actions": serialized,
	})
}

// HandleActionExecute executes an action on selected resources.
// It validates permissions, loads models, checks action eligibility,
// and executes the action with proper error handling.
func HandleActionExecute(h *FieldHandler, c *context.Context) error {
	actionSlug := c.Params("action")

	// Policy check - user must have update permission
	if h.Policy != nil && !h.Policy.Update(c, nil) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Find the action
	var targetAction action.Action
	for _, act := range h.Resource.GetActions() {
		if newAction, ok := act.(action.Action); ok {
			if newAction.GetSlug() == actionSlug {
				targetAction = newAction
				break
			}
		}
	}

	if targetAction == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Action not found",
		})
	}

	// Parse request body
	var body struct {
		IDs    []string               `json:"ids"`
		Fields map[string]interface{} `json:"fields"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if len(body.IDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No items selected",
		})
	}

	// Load models
	models := make([]interface{}, 0, len(body.IDs))
	modelType := reflect.TypeOf(h.Resource.Model())

	// Handle pointer types
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	for _, id := range body.IDs {
		model := reflect.New(modelType).Interface()
		if err := h.DB.First(model, "id = ?", id).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": fmt.Sprintf("Model with ID %s not found", id),
			})
		}
		models = append(models, model)
	}

	// Store fields and DB in context locals for action execution
	c.Locals("action_fields", body.Fields)
	c.Locals("db", h.DB)

	// Create action context for CanRun check
	ctx := &action.ActionContext{
		Models:   models,
		Fields:   body.Fields,
		User:     c.Locals("user"),
		Resource: h.Resource.Slug(),
		DB:       h.DB,
		Ctx:      c.Ctx,
	}

	// Check if action can run
	if !targetAction.CanRun(ctx) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Action cannot be executed in this context",
		})
	}

	// Execute action with new signature
	if err := targetAction.Execute(c, models); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Action executed successfully on %d item(s)", len(models)),
		"count":   len(models),
	})
}
