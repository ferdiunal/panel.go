package query

// FilterOperator represents a filter operation type
type FilterOperator string

const (
	// Equality operators
	OpEqual    FilterOperator = "eq"  // field = value
	OpNotEqual FilterOperator = "neq" // field != value

	// Comparison operators
	OpGreaterThan FilterOperator = "gt"  // field > value
	OpGreaterEq   FilterOperator = "gte" // field >= value
	OpLessThan    FilterOperator = "lt"  // field < value
	OpLessEq      FilterOperator = "lte" // field <= value

	// String operators
	OpLike    FilterOperator = "like"  // field LIKE %value%
	OpNotLike FilterOperator = "nlike" // field NOT LIKE %value%

	// List operators
	OpIn    FilterOperator = "in"  // field IN (values...)
	OpNotIn FilterOperator = "nin" // field NOT IN (values...)

	// Null operators
	OpIsNull    FilterOperator = "null"  // field IS NULL
	OpIsNotNull FilterOperator = "nnull" // field IS NOT NULL

	// Range operator
	OpBetween FilterOperator = "between" // field BETWEEN value1 AND value2
)

// Filter represents a single filter condition
type Filter struct {
	Field    string         `json:"field"`
	Operator FilterOperator `json:"operator"`
	Value    interface{}    `json:"value"`
}

// FilterGroup represents a group of filters with AND/OR logic
type FilterGroup struct {
	Logic   string   `json:"logic"` // "and" or "or"
	Filters []Filter `json:"filters"`
}

// validOperators is the list of all valid operators
var validOperators = []FilterOperator{
	OpEqual, OpNotEqual,
	OpGreaterThan, OpGreaterEq, OpLessThan, OpLessEq,
	OpLike, OpNotLike,
	OpIn, OpNotIn,
	OpIsNull, OpIsNotNull,
	OpBetween,
}

// ValidOperators returns list of valid operators
func ValidOperators() []FilterOperator {
	return validOperators
}

// IsValidOperator checks if operator is valid
func IsValidOperator(op string) bool {
	for _, valid := range validOperators {
		if string(valid) == op {
			return true
		}
	}
	return false
}

// GetOperator returns FilterOperator from string, defaults to OpEqual
func GetOperator(op string) FilterOperator {
	if IsValidOperator(op) {
		return FilterOperator(op)
	}
	return OpEqual
}
