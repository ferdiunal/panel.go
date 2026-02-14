package migration

import (
	"testing"

	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/stretchr/testify/assert"
)

func TestTypeMapperFieldToGoType(t *testing.T) {
	tm := NewTypeMapper()

	tests := []struct {
		name      string
		fieldType fields.ElementType
		nullable  bool
		wantType  string
	}{
		{"text", fields.TYPE_TEXT, false, "string"},
		{"text nullable", fields.TYPE_TEXT, true, "string"}, // pointer check
		{"number", fields.TYPE_NUMBER, false, "int64"},
		{"boolean", fields.TYPE_BOOLEAN, false, "bool"},
		{"email", fields.TYPE_EMAIL, false, "string"},
		{"select", fields.TYPE_SELECT, false, "string"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tm.MapFieldTypeToGo(tt.fieldType, tt.nullable)
			if result.Type != nil {
				assert.Equal(t, tt.wantType, result.Type.String())
			}
		})
	}
}

func TestTypeMapperFieldToSQLType(t *testing.T) {
	tm := NewTypeMapper()

	tests := []struct {
		name      string
		fieldType fields.ElementType
		size      int
		wantSQL   string
	}{
		{"text default", fields.TYPE_TEXT, 0, "varchar(255)"},
		{"text custom size", fields.TYPE_TEXT, 100, "varchar(100)"},
		{"number", fields.TYPE_NUMBER, 0, "bigint"},
		{"boolean", fields.TYPE_BOOLEAN, 0, "boolean"},
		{"date", fields.TYPE_DATE, 0, "date"},
		{"datetime", fields.TYPE_DATETIME, 0, "timestamp"},
		{"file", fields.TYPE_FILE, 0, "text"},
		{"select", fields.TYPE_SELECT, 0, "varchar(100)"},
		{"key_value", fields.TYPE_KEY_VALUE, 0, "jsonb"},
		{"link (FK)", fields.TYPE_LINK, 0, "bigint"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tm.MapFieldTypeToSQL(tt.fieldType, tt.size)
			assert.Equal(t, tt.wantSQL, result)
		})
	}
}

func TestTypeMapperRelationshipType(t *testing.T) {
	tm := NewTypeMapper()

	tests := []struct {
		fieldType    fields.ElementType
		isRelation   bool
		relationType string
	}{
		{fields.TYPE_LINK, true, "belongsTo"},
		{fields.TYPE_DETAIL, true, "hasOne"},
		{fields.TYPE_COLLECTION, true, "hasMany"},
		{fields.TYPE_CONNECT, true, "belongsToMany"},
		{fields.TYPE_TEXT, false, ""},
		{fields.TYPE_NUMBER, false, ""},
	}

	for _, tt := range tests {
		t.Run(string(tt.fieldType), func(t *testing.T) {
			assert.Equal(t, tt.isRelation, tm.IsRelationshipType(tt.fieldType))
			assert.Equal(t, tt.relationType, tm.GetRelationshipType(tt.fieldType))
		})
	}
}
