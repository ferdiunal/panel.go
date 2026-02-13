package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type TestModel struct {
	ID        uint   `gorm:"primarykey"`
	Name      string `gorm:"column:name"`
	Email     string `gorm:"column:email"`
	Age       int    `gorm:"column:age"`
	CreatedAt int64  `gorm:"column:created_at"`
}

func TestColumnValidator(t *testing.T) {
	// Setup test database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Create validator
	validator, err := NewColumnValidator(db, &TestModel{})
	assert.NoError(t, err)
	assert.NotNil(t, validator)

	// Test valid columns
	t.Run("Valid columns", func(t *testing.T) {
		validColumns := []string{"id", "name", "email", "age", "created_at", "ID", "Name", "Email"}
		for _, col := range validColumns {
			assert.True(t, validator.IsValidColumn(col), "Column %s should be valid", col)
		}
	})

	// Test invalid columns
	t.Run("Invalid columns", func(t *testing.T) {
		invalidColumns := []string{"invalid", "password", "secret", "admin", "1=1", "id OR 1=1"}
		for _, col := range invalidColumns {
			assert.False(t, validator.IsValidColumn(col), "Column %s should be invalid", col)
		}
	})

	// Test ValidateColumn
	t.Run("ValidateColumn returns DB column name", func(t *testing.T) {
		dbCol, err := validator.ValidateColumn("name")
		assert.NoError(t, err)
		assert.Equal(t, "name", dbCol)

		dbCol, err = validator.ValidateColumn("Name")
		assert.NoError(t, err)
		assert.Equal(t, "name", dbCol)
	})

	// Test ValidateColumn with invalid column
	t.Run("ValidateColumn rejects invalid column", func(t *testing.T) {
		_, err := validator.ValidateColumn("invalid_column")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid column name")
	})
}

func TestSanitizeColumnName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"name", "name"},
		{"user_id", "user_id"},
		{"table.column", "table.column"},
		{"id OR 1=1", "idOR11"},
		{"name; DROP TABLE users", "nameDROPTABLEusers"},
		{"id' OR '1'='1", "idOR11"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := SanitizeColumnName(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsValidOperator(t *testing.T) {
	validOps := []string{"=", "!=", ">", ">=", "<", "<=", "LIKE", "NOT LIKE", "IN", "NOT IN", "IS NULL", "IS NOT NULL", "BETWEEN"}
	for _, op := range validOps {
		assert.True(t, IsValidOperator(op), "Operator %s should be valid", op)
	}

	invalidOps := []string{"DROP", "DELETE", "UPDATE", "INSERT", "OR", "AND", "UNION"}
	for _, op := range invalidOps {
		assert.False(t, IsValidOperator(op), "Operator %s should be invalid", op)
	}
}

func TestBuildSafeWhereClause(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	validator, err := NewColumnValidator(db, &TestModel{})
	assert.NoError(t, err)

	// Test valid clause
	clause, err := BuildSafeWhereClause(validator, "name", "=")
	assert.NoError(t, err)
	assert.Equal(t, "name = ?", clause)

	// Test invalid column
	_, err = BuildSafeWhereClause(validator, "invalid_column", "=")
	assert.Error(t, err)

	// Test invalid operator
	_, err = BuildSafeWhereClause(validator, "name", "DROP")
	assert.Error(t, err)
}
