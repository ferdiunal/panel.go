package fields

import (
	"context"
	"fmt"
)

// LoadingStrategy defines how relationships are loaded
type LoadingStrategy string

const (
	EAGER_LOADING LoadingStrategy = "eager"
	LAZY_LOADING  LoadingStrategy = "lazy"
)

// RelationshipField represents a database relationship in the field system
type RelationshipField interface {
	Element

	// Relationship Type Methods
	GetRelationshipType() string // "belongsTo", "hasMany", "hasOne", "belongsToMany", "morphTo"
	GetRelatedResource() string  // Related resource slug
	GetRelationshipName() string // Relationship name

	// Relationship Resolution
	ResolveRelationship(item interface{}) (interface{}, error)

	// Query Customization
	GetQueryCallback() func(interface{}) interface{}

	// Loading Strategy
	GetLoadingStrategy() LoadingStrategy

	// Relationship Validation
	ValidateRelationship(value interface{}) error

	// Relationship Display
	GetDisplayKey() string          // Key to display for BelongsTo
	GetSearchableColumns() []string // Searchable columns for BelongsTo

	// Required check
	IsRequired() bool

	// Get types for MorphTo
	GetTypes() map[string]string
}

// RelationshipError represents an error that occurred during relationship operations
type RelationshipError struct {
	FieldName        string
	RelationshipType string
	Message          string
	Context          map[string]interface{}
}

// Error implements the error interface
func (e *RelationshipError) Error() string {
	return fmt.Sprintf("relationship error in field '%s' (%s): %s", e.FieldName, e.RelationshipType, e.Message)
}

// RelationshipLoader handles loading relationships with different strategies
type RelationshipLoader interface {
	// Load related data using eager loading strategy
	EagerLoad(ctx context.Context, items []interface{}, field RelationshipField) error

	// Load related data using lazy loading strategy
	LazyLoad(ctx context.Context, item interface{}, field RelationshipField) (interface{}, error)

	// Load with constraints applied
	LoadWithConstraints(ctx context.Context, item interface{}, field RelationshipField, constraints map[string]interface{}) (interface{}, error)
}

// RelationshipValidator handles validation of relationships
type RelationshipValidator interface {
	// Validate that related resource exists
	ValidateExists(ctx context.Context, value interface{}, field RelationshipField) error

	// Validate foreign key references
	ValidateForeignKey(ctx context.Context, value interface{}, field RelationshipField) error

	// Validate pivot table entries
	ValidatePivot(ctx context.Context, value interface{}, field RelationshipField) error

	// Validate morph type is registered
	ValidateMorphType(ctx context.Context, value interface{}, field RelationshipField) error
}

// RelationshipQuery represents a query builder for relationships
type RelationshipQuery interface {
	// Apply WHERE clause
	Where(column string, operator string, value interface{}) RelationshipQuery

	// Apply WHERE IN clause
	WhereIn(column string, values []interface{}) RelationshipQuery

	// Apply ORDER BY clause
	OrderBy(column string, direction string) RelationshipQuery

	// Apply LIMIT clause
	Limit(limit int) RelationshipQuery

	// Apply OFFSET clause
	Offset(offset int) RelationshipQuery

	// Get count of results
	Count(ctx context.Context) (int64, error)

	// Check if results exist
	Exists(ctx context.Context) (bool, error)

	// Execute query and get results
	Get(ctx context.Context) ([]interface{}, error)

	// Execute query and get first result
	First(ctx context.Context) (interface{}, error)
}

// IsRelationshipField checks if an element is a relationship field
func IsRelationshipField(e Element) (RelationshipField, bool) {
	rf, ok := e.(RelationshipField)
	return rf, ok
}
