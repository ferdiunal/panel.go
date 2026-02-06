package query

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Sort represents a sort configuration
type Sort struct {
	Column    string `json:"column"`
	Direction string `json:"direction"`
}

// ResourceQueryParams represents parsed query params for a resource
type ResourceQueryParams struct {
	Search  string   // resource[search]=...
	Sorts   []Sort   // resource[sort][column]=direction
	Filters []Filter // resource[filters][field][op]=value
	Page    int      // resource[page]=...
	PerPage int      // resource[per_page]=...
}

// DefaultParams returns default query params
func DefaultParams() *ResourceQueryParams {
	return &ResourceQueryParams{
		Page:    1,
		PerPage: 15,
		Filters: make([]Filter, 0),
		Sorts:   make([]Sort, 0),
	}
}

// ParseResourceQuery parses nested query params like resource[search]=...
// Falls back to legacy format if nested format not found
//
// Supported formats:
//   - Nested: users[search]=query, users[sort][id]=asc, users[filters][status][eq]=active
//   - Legacy: search=query, sort_column=id, sort_direction=asc
func ParseResourceQuery(c *fiber.Ctx, resourceName string) *ResourceQueryParams {
	params := DefaultParams()

	// Get raw query string and decode it
	rawQuery := string(c.Request().URI().QueryString())

	fmt.Printf("[PARSER] Resource: %s, RawQuery: %s\n", resourceName, rawQuery)

	// Try nested format first using decoded query string
	if parseNestedFormat(rawQuery, resourceName, params) {
		fmt.Printf("[PARSER] Nested format parsed: Search=%q, Page=%d, PerPage=%d\n", params.Search, params.Page, params.PerPage)
		return params
	}

	fmt.Printf("[PARSER] Nested format not found, trying legacy\n")

	// Fallback to legacy format
	parseLegacyFormat(c, params)
	return params
}

// parseNestedFormat parses the new nested format: resource[key]=value
func parseNestedFormat(rawQuery string, resource string, params *ResourceQueryParams) bool {
	if rawQuery == "" {
		fmt.Printf("[NESTED] Empty rawQuery\n")
		return false
	}

	// URL decode the query string first
	decodedQuery, err := url.QueryUnescape(rawQuery)
	if err != nil {
		decodedQuery = rawQuery
	}

	fmt.Printf("[NESTED] Decoded query: %s\n", decodedQuery)

	// Parse the decoded query string
	values, err := url.ParseQuery(decodedQuery)
	if err != nil {
		fmt.Printf("[NESTED] Parse error: %v\n", err)
		return false
	}

	fmt.Printf("[NESTED] Parsed values: %+v\n", values)

	found := false
	prefix := resource + "["
	fmt.Printf("[NESTED] Looking for prefix: %s\n", prefix)

	for key, vals := range values {
		fmt.Printf("[NESTED] Key: %s, Vals: %v\n", key, vals)
		if len(vals) == 0 {
			continue
		}
		value := vals[0] // Take first value

		if !strings.HasPrefix(key, prefix) {
			continue
		}
		found = true

		// Remove resource prefix and trailing bracket
		// users[search] -> search
		// users[sort][id] -> sort][id
		// users[filters][status][eq] -> filters][status][eq
		inner := strings.TrimPrefix(key, prefix)
		inner = strings.TrimSuffix(inner, "]")

		switch {
		case inner == "search":
			params.Search = value

		case inner == "page":
			if p, err := strconv.Atoi(value); err == nil && p > 0 {
				params.Page = p
			}

		case inner == "per_page":
			if pp, err := strconv.Atoi(value); err == nil && pp > 0 && pp <= 100 {
				params.PerPage = pp
			}

		case strings.HasPrefix(inner, "sort]["):
			// sort][name -> name
			column := strings.TrimPrefix(inner, "sort][")
			if column != "" {
				direction := strings.ToLower(value)
				if direction != "asc" && direction != "desc" {
					direction = "asc"
				}
				params.Sorts = append(params.Sorts, Sort{
					Column:    column,
					Direction: direction,
				})
			}

		case strings.HasPrefix(inner, "filters]["):
			parseFilterParam(inner, value, params)
		}
	}

	return found
}

// parseFilterParam handles both simple and advanced filter formats
//
// Simple format (defaults to eq operator):
//
//	filters][status = "active" -> {field: status, op: eq, value: active}
//
// Advanced format:
//
//	filters][status][eq = "active" -> {field: status, op: eq, value: active}
//	filters][age][gt = "18" -> {field: age, op: gt, value: 18}
//	filters][status][in = "active,pending" -> {field: status, op: in, value: [active, pending]}
func parseFilterParam(inner, value string, params *ResourceQueryParams) {
	// Remove "filters][" prefix
	rest := strings.TrimPrefix(inner, "filters][")

	// Split by "][" to get parts
	parts := strings.Split(rest, "][")

	var field string
	var operator FilterOperator = OpEqual // default operator

	if len(parts) == 1 {
		// Simple format: filters][status = "active"
		field = parts[0]
	} else if len(parts) >= 2 {
		// Advanced format: filters][status][eq = "active"
		field = parts[0]
		if IsValidOperator(parts[1]) {
			operator = FilterOperator(parts[1])
		}
	}

	if field == "" {
		return
	}

	// Parse value based on operator type
	var parsedValue interface{}

	switch operator {
	case OpIn, OpNotIn:
		// Comma-separated values -> []string
		parsedValue = strings.Split(value, ",")

	case OpBetween:
		// Two comma-separated values -> []string
		betweenParts := strings.Split(value, ",")
		if len(betweenParts) == 2 {
			parsedValue = betweenParts
		} else {
			// Invalid between format, skip
			return
		}

	case OpIsNull, OpIsNotNull:
		// Boolean value
		parsedValue = value == "true" || value == "1"

	default:
		// String value for eq, neq, gt, gte, lt, lte, like, nlike
		parsedValue = value
	}

	params.Filters = append(params.Filters, Filter{
		Field:    field,
		Operator: operator,
		Value:    parsedValue,
	})
}

// parseLegacyFormat parses the old flat format for backward compatibility
func parseLegacyFormat(c *fiber.Ctx, params *ResourceQueryParams) {
	// Page
	if p, err := strconv.Atoi(c.Query("page", "1")); err == nil && p > 0 {
		params.Page = p
	}

	// Per page
	if pp, err := strconv.Atoi(c.Query("per_page", "15")); err == nil && pp > 0 && pp <= 100 {
		params.PerPage = pp
	}

	// Search
	if search := c.Query("search"); search != "" {
		params.Search = search
	}

	// Sort (legacy format: sort_column + sort_direction)
	if col := c.Query("sort_column"); col != "" {
		dir := strings.ToLower(c.Query("sort_direction", "asc"))
		if dir != "asc" && dir != "desc" {
			dir = "asc"
		}
		params.Sorts = append(params.Sorts, Sort{
			Column:    col,
			Direction: dir,
		})
	}

	// Legacy filter format using QueryParser for map
	type LegacyFilters struct {
		Filters map[string]string `query:"filters"`
	}
	lf := new(LegacyFilters)
	if err := c.QueryParser(lf); err == nil && len(lf.Filters) > 0 {
		for field, value := range lf.Filters {
			params.Filters = append(params.Filters, Filter{
				Field:    field,
				Operator: OpEqual,
				Value:    value,
			})
		}
	}
}

// HasSearch returns true if search query is set
func (p *ResourceQueryParams) HasSearch() bool {
	return p.Search != ""
}

// HasSorts returns true if any sorts are defined
func (p *ResourceQueryParams) HasSorts() bool {
	return len(p.Sorts) > 0
}

// HasFilters returns true if any filters are defined
func (p *ResourceQueryParams) HasFilters() bool {
	return len(p.Filters) > 0
}
