package fields

import (
	"fmt"
)

// RelationshipDisplay handles display customization for relationships
type RelationshipDisplay interface {
	// GetDisplayValue returns the display value for a relationship
	GetDisplayValue(item interface{}) (string, error)

	// GetDisplayValues returns display values for multiple items
	GetDisplayValues(items []interface{}) ([]string, error)

	// FormatDisplayValue formats a display value
	FormatDisplayValue(value interface{}) string
}

// RelationshipDisplayImpl implements RelationshipDisplay
type RelationshipDisplayImpl struct {
	field RelationshipField
}

// NewRelationshipDisplay creates a new relationship display handler
func NewRelationshipDisplay(field RelationshipField) *RelationshipDisplayImpl {
	return &RelationshipDisplayImpl{
		field: field,
	}
}

// GetDisplayValue returns the display value for a relationship
func (rd *RelationshipDisplayImpl) GetDisplayValue(item interface{}) (string, error) {
	if item == nil {
		return "", nil
	}

	relationType := rd.field.GetRelationshipType()

	switch relationType {
	case "belongsTo":
		return rd.displayBelongsTo(item)
	case "hasMany":
		return rd.displayHasMany(item)
	case "hasOne":
		return rd.displayHasOne(item)
	case "belongsToMany":
		return rd.displayBelongsToMany(item)
	case "morphTo":
		return rd.displayMorphTo(item)
	default:
		return "", fmt.Errorf("unknown relationship type: %s", relationType)
	}
}

// GetDisplayValues returns display values for multiple items
func (rd *RelationshipDisplayImpl) GetDisplayValues(items []interface{}) ([]string, error) {
	if items == nil || len(items) == 0 {
		return []string{}, nil
	}

	displayValues := make([]string, 0, len(items))

	for _, item := range items {
		displayValue, err := rd.GetDisplayValue(item)
		if err != nil {
			return nil, err
		}
		displayValues = append(displayValues, displayValue)
	}

	return displayValues, nil
}

// FormatDisplayValue formats a display value
func (rd *RelationshipDisplayImpl) FormatDisplayValue(value interface{}) string {
	if value == nil {
		return ""
	}

	return fmt.Sprintf("%v", value)
}

// displayBelongsTo displays a BelongsTo relationship
func (rd *RelationshipDisplayImpl) displayBelongsTo(item interface{}) (string, error) {
	// For BelongsTo, show the value from the DisplayKey column
	// In a real implementation, this would extract the value from the related resource
	displayKey := rd.field.GetDisplayKey()
	if displayKey == "" {
		displayKey = "name"
	}

	return fmt.Sprintf("Related resource (key: %s)", displayKey), nil
}

// displayHasMany displays a HasMany relationship
func (rd *RelationshipDisplayImpl) displayHasMany(item interface{}) (string, error) {
	// For HasMany, show count or list of related resources
	// In a real implementation, this would count or list the related resources
	return "Multiple related resources", nil
}

// displayHasOne displays a HasOne relationship
func (rd *RelationshipDisplayImpl) displayHasOne(item interface{}) (string, error) {
	// For HasOne, show the related resource or empty state
	// In a real implementation, this would show the related resource or "No related resource"
	if item == nil {
		return "No related resource", nil
	}

	return "Related resource", nil
}

// displayBelongsToMany displays a BelongsToMany relationship
func (rd *RelationshipDisplayImpl) displayBelongsToMany(item interface{}) (string, error) {
	// For BelongsToMany, show list of related resources
	// In a real implementation, this would list the related resources
	return "Multiple related resources", nil
}

// displayMorphTo displays a MorphTo relationship
func (rd *RelationshipDisplayImpl) displayMorphTo(item interface{}) (string, error) {
	// For MorphTo, show related resource with type indicator
	// In a real implementation, this would show the related resource with its type
	return "Polymorphic related resource", nil
}
