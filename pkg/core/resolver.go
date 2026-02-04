package core

// Resolver is an interface for resolving field values dynamically.
// It allows fields to compute or transform their values based on the item and parameters.
type Resolver interface {
	// Resolve computes or transforms a field value based on the item and parameters.
	// The item parameter is typically a struct or map containing the data to resolve.
	// The params parameter contains additional parameters for the resolution.
	// Returns the resolved value or an error if resolution fails.
	Resolve(item interface{}, params map[string]interface{}) (interface{}, error)
}
