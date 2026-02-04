package fields

import (
	"context"
	"testing"
)

// TestNewRelationshipQuery tests creating a new relationship query
func TestNewRelationshipQuery(t *testing.T) {
	query := NewRelationshipQuery()

	if query == nil {
		t.Error("Expected non-nil query")
	}

	if len(query.GetWhereConditions()) != 0 {
		t.Error("Expected empty where conditions")
	}

	if len(query.GetOrderByColumns()) != 0 {
		t.Error("Expected empty order by columns")
	}

	if query.GetLimit() != 0 {
		t.Error("Expected limit to be 0")
	}

	if query.GetOffset() != 0 {
		t.Error("Expected offset to be 0")
	}
}

// TestRelationshipQueryWhere tests Where method
func TestRelationshipQueryWhere(t *testing.T) {
	query := NewRelationshipQuery()
	result := query.Where("name", "=", "John")

	if result != query {
		t.Error("Where should return the query for chaining")
	}

	conditions := query.GetWhereConditions()
	if len(conditions) != 1 {
		t.Errorf("Expected 1 where condition, got %d", len(conditions))
	}

	if conditions[0]["column"] != "name" {
		t.Errorf("Expected column 'name', got '%v'", conditions[0]["column"])
	}

	if conditions[0]["operator"] != "=" {
		t.Errorf("Expected operator '=', got '%v'", conditions[0]["operator"])
	}

	if conditions[0]["value"] != "John" {
		t.Errorf("Expected value 'John', got '%v'", conditions[0]["value"])
	}
}

// TestRelationshipQueryWhereIn tests WhereIn method
func TestRelationshipQueryWhereIn(t *testing.T) {
	query := NewRelationshipQuery()
	values := []interface{}{1, 2, 3}
	result := query.WhereIn("id", values)

	if result != query {
		t.Error("WhereIn should return the query for chaining")
	}

	conditions := query.GetWhereConditions()
	if len(conditions) != 1 {
		t.Errorf("Expected 1 where condition, got %d", len(conditions))
	}

	if conditions[0]["column"] != "id" {
		t.Errorf("Expected column 'id', got '%v'", conditions[0]["column"])
	}

	if len(conditions[0]["in"].([]interface{})) != 3 {
		t.Errorf("Expected 3 values in IN clause, got %d", len(conditions[0]["in"].([]interface{})))
	}
}

// TestRelationshipQueryOrderBy tests OrderBy method
func TestRelationshipQueryOrderBy(t *testing.T) {
	query := NewRelationshipQuery()
	result := query.OrderBy("name", "ASC")

	if result != query {
		t.Error("OrderBy should return the query for chaining")
	}

	orderBy := query.GetOrderByColumns()
	if len(orderBy) != 1 {
		t.Errorf("Expected 1 order by column, got %d", len(orderBy))
	}

	if orderBy[0]["column"] != "name" {
		t.Errorf("Expected column 'name', got '%v'", orderBy[0]["column"])
	}

	if orderBy[0]["direction"] != "ASC" {
		t.Errorf("Expected direction 'ASC', got '%v'", orderBy[0]["direction"])
	}
}

// TestRelationshipQueryLimit tests Limit method
func TestRelationshipQueryLimit(t *testing.T) {
	query := NewRelationshipQuery()
	result := query.Limit(10)

	if result != query {
		t.Error("Limit should return the query for chaining")
	}

	if query.GetLimit() != 10 {
		t.Errorf("Expected limit 10, got %d", query.GetLimit())
	}
}

// TestRelationshipQueryOffset tests Offset method
func TestRelationshipQueryOffset(t *testing.T) {
	query := NewRelationshipQuery()
	result := query.Offset(20)

	if result != query {
		t.Error("Offset should return the query for chaining")
	}

	if query.GetOffset() != 20 {
		t.Errorf("Expected offset 20, got %d", query.GetOffset())
	}
}

// TestRelationshipQueryCount tests Count method
func TestRelationshipQueryCount(t *testing.T) {
	query := NewRelationshipQuery()
	ctx := context.Background()

	count, err := query.Count(ctx)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if count != 0 {
		t.Errorf("Expected count 0, got %d", count)
	}
}

// TestRelationshipQueryExists tests Exists method
func TestRelationshipQueryExists(t *testing.T) {
	query := NewRelationshipQuery()
	ctx := context.Background()

	exists, err := query.Exists(ctx)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if exists != false {
		t.Errorf("Expected exists false, got %v", exists)
	}
}

// TestRelationshipQueryGet tests Get method
func TestRelationshipQueryGet(t *testing.T) {
	query := NewRelationshipQuery()
	ctx := context.Background()

	results, err := query.Get(ctx)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(results) != 0 {
		t.Errorf("Expected 0 results, got %d", len(results))
	}
}

// TestRelationshipQueryFirst tests First method
func TestRelationshipQueryFirst(t *testing.T) {
	query := NewRelationshipQuery()
	ctx := context.Background()

	result, err := query.First(ctx)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result != nil {
		t.Errorf("Expected nil result, got %v", result)
	}
}

// TestRelationshipQueryChaining tests method chaining
func TestRelationshipQueryChaining(t *testing.T) {
	query := NewRelationshipQuery()
	result := query.
		Where("status", "=", "active").
		OrderBy("name", "ASC").
		Limit(10).
		Offset(5)

	if result != query {
		t.Error("Chaining should return the same query")
	}

	if len(query.GetWhereConditions()) != 1 {
		t.Errorf("Expected 1 where condition, got %d", len(query.GetWhereConditions()))
	}

	if len(query.GetOrderByColumns()) != 1 {
		t.Errorf("Expected 1 order by column, got %d", len(query.GetOrderByColumns()))
	}

	if query.GetLimit() != 10 {
		t.Errorf("Expected limit 10, got %d", query.GetLimit())
	}

	if query.GetOffset() != 5 {
		t.Errorf("Expected offset 5, got %d", query.GetOffset())
	}
}

// TestRelationshipQueryString tests String method
func TestRelationshipQueryString(t *testing.T) {
	query := NewRelationshipQuery().
		Where("status", "=", "active").
		OrderBy("name", "ASC").
		Limit(10).
		Offset(5)

	// Cast to implementation to access String method
	impl := query.(*RelationshipQueryImpl)
	queryStr := impl.String()

	if queryStr == "" {
		t.Error("Expected non-empty query string")
	}

	if !contains(queryStr, "WHERE") {
		t.Error("Expected WHERE clause in query string")
	}

	if !contains(queryStr, "ORDER BY") {
		t.Error("Expected ORDER BY clause in query string")
	}

	if !contains(queryStr, "LIMIT") {
		t.Error("Expected LIMIT clause in query string")
	}

	if !contains(queryStr, "OFFSET") {
		t.Error("Expected OFFSET clause in query string")
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
