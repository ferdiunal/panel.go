// Package handler provides HTTP request handlers for the panel API.
// This file contains the MorphableController for handling MorphTo relationship options.
package handler

import (
	"fmt"

	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// MorphableController handles MorphTo relationship resource listing.
// It provides options for polymorphic relationship fields based on the selected type.
type MorphableController struct {
	DB        *gorm.DB
	Resources map[string]interface{} // Resource registry
}

// NewMorphableController creates a new MorphableController instance.
func NewMorphableController(db *gorm.DB, resources map[string]interface{}) *MorphableController {
	return &MorphableController{
		DB:        db,
		Resources: resources,
	}
}

// MorphableOption represents a single option for the MorphTo field.
type MorphableOption struct {
	Value    interface{} `json:"value"`
	Display  string      `json:"display"`
	Avatar   string      `json:"avatar,omitempty"`
	Subtitle string      `json:"subtitle,omitempty"`
}

// HandleMorphable handles GET /api/resource/:resource/morphable/:field
// It returns a list of resources that can be associated with a MorphTo field.
//
// Request Parameters:
//   - :resource - The resource slug containing the MorphTo field
//   - :field - The MorphTo field key
//   - type (query) - The resource type/slug to fetch options from (e.g., "posts")
//   - search (query) - Search query for filtering results
//   - per_page (query) - Number of results per page (default: 10)
//   - current (query) - Current selected value ID (for loading initial value)
//
// Response:
//
//	{
//	  "resources": [...],
//	  "softDeletes": false
//	}
//
// Requirement: MorphTo alanları için ilgili kaynakları listeleme
func HandleMorphable(h *FieldHandler, c *context.Context) error {
	fieldKey := c.Params("field")
	resourceType := c.Query("type")
	search := c.Query("search", "")
	perPage := c.QueryInt("per_page", 10)
	current := c.Query("current", "")

	// Find the MorphTo field
	var morphToField *fields.MorphTo
	for _, element := range h.Elements {
		if element.GetKey() == fieldKey {
			if mt, ok := element.(*fields.MorphTo); ok {
				morphToField = mt
				break
			}
		}
	}

	if morphToField == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "MorphTo field not found",
		})
	}

	// Validate type is registered in MorphTo types
	if resourceType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "type parameter is required",
		})
	}

	// Check if type is valid
	resourceSlug, err := morphToField.GetResourceForType(resourceType)
	if err != nil {
		// Type might be the slug directly, not the database type
		// Check if it exists in the type mappings as a value (slug)
		found := false
		for _, slug := range morphToField.GetTypes() {
			if slug == resourceType {
				resourceSlug = resourceType
				found = true
				break
			}
		}
		if !found {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("Invalid type: %s", resourceType),
			})
		}
	}

	// Query the related resource table
	options, err := queryMorphableResources(h.DB, resourceSlug, search, perPage, current)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch resources",
		})
	}

	return c.JSON(fiber.Map{
		"resources":   options,
		"softDeletes": false, // TODO: Implement soft delete detection
	})
}

// queryMorphableResources queries the database for morphable resource options.
// It returns a list of options with id and display value.
func queryMorphableResources(db *gorm.DB, tableName, search string, limit int, currentID string) ([]MorphableOption, error) {
	if db == nil {
		return []MorphableOption{}, nil
	}

	var results []map[string]interface{}

	// Build query
	query := db.Table(tableName).Select("id, name, title, email, username")

	// Apply search filter
	if search != "" {
		query = query.Where(
			"name LIKE ? OR title LIKE ? OR email LIKE ? OR username LIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%",
		)
	}

	// Apply limit
	if limit > 0 {
		query = query.Limit(limit)
	}

	// Execute query
	if err := query.Find(&results).Error; err != nil {
		return nil, err
	}

	// Format results
	options := make([]MorphableOption, 0, len(results))
	for _, r := range results {
		display := getDisplayValue(r)
		options = append(options, MorphableOption{
			Value:   r["id"],
			Display: display,
		})
	}

	// If current ID is provided and not in results, fetch it separately
	if currentID != "" {
		found := false
		for _, opt := range options {
			if fmt.Sprint(opt.Value) == currentID {
				found = true
				break
			}
		}

		if !found {
			var currentResult map[string]interface{}
			if err := db.Table(tableName).
				Select("id, name, title, email, username").
				Where("id = ?", currentID).
				First(&currentResult).Error; err == nil {
				display := getDisplayValue(currentResult)
				options = append([]MorphableOption{{
					Value:   currentResult["id"],
					Display: display,
				}}, options...)
			}
		}
	}

	return options, nil
}

// getDisplayValue extracts the best display value from a resource record.
func getDisplayValue(r map[string]interface{}) string {
	// Priority: name > title > email > username > id
	displayFields := []string{"name", "title", "email", "username"}

	for _, field := range displayFields {
		if val, ok := r[field]; ok && val != nil {
			str := fmt.Sprint(val)
			if str != "" && str != "<nil>" {
				return str
			}
		}
	}

	// Fallback to ID
	if id, ok := r["id"]; ok {
		return fmt.Sprintf("#%v", id)
	}

	return "Unknown"
}
