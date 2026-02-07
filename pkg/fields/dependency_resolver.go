package fields

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// DependencyResolver handles field dependency resolution
type DependencyResolver struct {
	fields  []*Schema
	context string
}

// NewDependencyResolver creates a new dependency resolver
func NewDependencyResolver(fields []*Schema, context string) *DependencyResolver {
	return &DependencyResolver{
		fields:  fields,
		context: context,
	}
}

// ResolveDependencies resolves field dependencies based on changed fields and form data
func (r *DependencyResolver) ResolveDependencies(
	formData map[string]interface{},
	changedFields []string,
	ctx *fiber.Ctx,
) (map[string]*FieldUpdate, error) {
	updates := make(map[string]*FieldUpdate)

	// Build dependency graph
	dependencyGraph := r.buildDependencyGraph()

	// Find affected fields
	affectedFields := r.findAffectedFields(dependencyGraph, changedFields)

	// Execute callbacks for affected fields
	for _, fieldKey := range affectedFields {
		field := r.findFieldByKey(fieldKey)
		if field == nil {
			continue
		}

		// Get the appropriate callback based on context
		callback := field.GetDependencyCallback(r.context)
		if callback == nil {
			continue
		}

		// Execute callback
		update := callback(field, formData, ctx)
		if update != nil {
			updates[fieldKey] = update
		}
	}

	return updates, nil
}

// buildDependencyGraph builds a map of field dependencies
// Returns a map where key is the field that is depended upon,
// and value is a list of fields that depend on it
func (r *DependencyResolver) buildDependencyGraph() map[string][]string {
	graph := make(map[string][]string)

	for _, field := range r.fields {
		if len(field.DependsOnFields) == 0 {
			continue
		}

		for _, dependsOn := range field.DependsOnFields {
			if graph[dependsOn] == nil {
				graph[dependsOn] = []string{}
			}
			graph[dependsOn] = append(graph[dependsOn], field.Key)
		}
	}

	return graph
}

// findAffectedFields finds all fields affected by the changed fields
func (r *DependencyResolver) findAffectedFields(
	graph map[string][]string,
	changedFields []string,
) []string {
	affected := make(map[string]bool)
	visited := make(map[string]bool)

	// BFS to find all affected fields
	queue := make([]string, len(changedFields))
	copy(queue, changedFields)

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if visited[current] {
			continue
		}
		visited[current] = true

		// Get fields that depend on current field
		dependents := graph[current]
		for _, dependent := range dependents {
			affected[dependent] = true

			// Check for circular dependencies
			if !visited[dependent] {
				queue = append(queue, dependent)
			}
		}
	}

	// Convert map to slice
	result := make([]string, 0, len(affected))
	for field := range affected {
		result = append(result, field)
	}

	return result
}

// findFieldByKey finds a field by its key
func (r *DependencyResolver) findFieldByKey(key string) *Schema {
	for _, field := range r.fields {
		if field.Key == key {
			return field
		}
	}
	return nil
}

// DetectCircularDependencies detects circular dependencies in the field graph
func (r *DependencyResolver) DetectCircularDependencies() error {
	graph := r.buildDependencyGraph()
	visited := make(map[string]bool)
	recStack := make(map[string]bool)

	for _, field := range r.fields {
		if !visited[field.Key] {
			if r.hasCycle(field.Key, graph, visited, recStack) {
				return fmt.Errorf("circular dependency detected involving field: %s", field.Key)
			}
		}
	}

	return nil
}

// hasCycle checks if there's a cycle starting from the given field
func (r *DependencyResolver) hasCycle(
	fieldKey string,
	graph map[string][]string,
	visited map[string]bool,
	recStack map[string]bool,
) bool {
	visited[fieldKey] = true
	recStack[fieldKey] = true

	// Get all fields that depend on this field
	dependents := graph[fieldKey]
	for _, dependent := range dependents {
		if !visited[dependent] {
			if r.hasCycle(dependent, graph, visited, recStack) {
				return true
			}
		} else if recStack[dependent] {
			return true
		}
	}

	recStack[fieldKey] = false
	return false
}
