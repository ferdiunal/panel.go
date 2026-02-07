package data

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// ColumnValidator validates column names against database schema
type ColumnValidator struct {
	allowedColumns map[string]bool
	schema         *schema.Schema
}

// NewColumnValidator creates a validator for a GORM model
func NewColumnValidator(db *gorm.DB, model interface{}) (*ColumnValidator, error) {
	stmt := &gorm.Statement{DB: db}
	if err := stmt.Parse(model); err != nil {
		return nil, fmt.Errorf("failed to parse model schema: %w", err)
	}

	validator := &ColumnValidator{
		allowedColumns: make(map[string]bool),
		schema:         stmt.Schema,
	}

	// Build whitelist of allowed columns
	for _, field := range stmt.Schema.Fields {
		if field.DBName != "" {
			// Add both snake_case (DB name) and original field name
			validator.allowedColumns[field.DBName] = true
			validator.allowedColumns[strings.ToLower(field.Name)] = true

			// Add camelCase version
			validator.allowedColumns[toCamelCase(field.DBName)] = true
		}
	}

	// Add relationship fields
	for name := range stmt.Schema.Relationships.Relations {
		validator.allowedColumns[strings.ToLower(name)] = true
		validator.allowedColumns[toSnakeCase(name)] = true
	}

	return validator, nil
}

// IsValidColumn checks if a column name is valid for the model
func (v *ColumnValidator) IsValidColumn(columnName string) bool {
	if columnName == "" {
		return false
	}

	// Normalize column name
	normalized := strings.ToLower(strings.TrimSpace(columnName))

	// Check against whitelist
	return v.allowedColumns[normalized]
}

// ValidateColumn validates a column name and returns the safe DB column name
func (v *ColumnValidator) ValidateColumn(columnName string) (string, error) {
	if !v.IsValidColumn(columnName) {
		return "", fmt.Errorf("invalid column name: %s", columnName)
	}

	// Find the actual DB column name
	normalized := strings.ToLower(strings.TrimSpace(columnName))

	// Try to find the field in schema
	for _, field := range v.schema.Fields {
		if strings.ToLower(field.DBName) == normalized ||
		   strings.ToLower(field.Name) == normalized {
			return field.DBName, nil
		}
	}

	// If not found in fields, might be a relationship
	for name := range v.schema.Relationships.Relations {
		if strings.ToLower(name) == normalized {
			return name, nil
		}
	}

	return columnName, nil
}

// GetAllowedColumns returns list of allowed column names
func (v *ColumnValidator) GetAllowedColumns() []string {
	columns := make([]string, 0, len(v.allowedColumns))
	for col := range v.allowedColumns {
		columns = append(columns, col)
	}
	return columns
}

// ValidateColumns validates multiple column names
func (v *ColumnValidator) ValidateColumns(columns []string) error {
	for _, col := range columns {
		if !v.IsValidColumn(col) {
			return fmt.Errorf("invalid column name: %s", col)
		}
	}
	return nil
}

// Helper functions

func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

func toCamelCase(s string) string {
	parts := strings.Split(s, "_")
	for i := 1; i < len(parts); i++ {
		if len(parts[i]) > 0 {
			parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
		}
	}
	return strings.Join(parts, "")
}

// SanitizeColumnName removes potentially dangerous characters
func SanitizeColumnName(columnName string) string {
	// Remove any characters that aren't alphanumeric, underscore, or dot
	var result strings.Builder
	for _, r := range columnName {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
		   (r >= '0' && r <= '9') || r == '_' || r == '.' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// IsValidOperator checks if a SQL operator is safe
func IsValidOperator(operator string) bool {
	validOperators := map[string]bool{
		"=":           true,
		"!=":          true,
		">":           true,
		">=":          true,
		"<":           true,
		"<=":          true,
		"LIKE":        true,
		"NOT LIKE":    true,
		"IN":          true,
		"NOT IN":      true,
		"IS NULL":     true,
		"IS NOT NULL": true,
		"BETWEEN":     true,
	}

	return validOperators[strings.ToUpper(strings.TrimSpace(operator))]
}

// BuildSafeWhereClause builds a safe WHERE clause with validated column
func BuildSafeWhereClause(validator *ColumnValidator, column, operator string) (string, error) {
	// Validate column
	safeColumn, err := validator.ValidateColumn(column)
	if err != nil {
		return "", err
	}

	// Validate operator
	if !IsValidOperator(operator) {
		return "", fmt.Errorf("invalid operator: %s", operator)
	}

	// Build safe clause
	return fmt.Sprintf("%s %s ?", safeColumn, operator), nil
}
