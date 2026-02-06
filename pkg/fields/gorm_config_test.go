package fields

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGormConfigBuilder(t *testing.T) {
	config := NewGormConfig().
		WithPrimaryKey().
		WithColumn("user_id").
		WithType("bigint").
		WithSize(100).
		WithIndex("idx_users").
		WithNotNull().
		WithDefault("0").
		WithComment("User identifier")

	assert.True(t, config.PrimaryKey)
	assert.True(t, config.AutoIncrement)
	assert.Equal(t, "user_id", config.Column)
	assert.Equal(t, "bigint", config.Type)
	assert.Equal(t, 100, config.Size)
	assert.True(t, config.Index)
	assert.Equal(t, "idx_users", config.IndexName)
	assert.True(t, config.NotNull)
	assert.Equal(t, "0", config.Default)
	assert.Equal(t, "User identifier", config.Comment)
}

func TestGormConfigToTag(t *testing.T) {
	tests := []struct {
		name     string
		config   *GormConfig
		wantPart string
	}{
		{
			"primary key",
			NewGormConfig().WithPrimaryKey(),
			"primaryKey",
		},
		{
			"unique index",
			NewGormConfig().WithUniqueIndex("uniq_email"),
			"uniqueIndex:uniq_email",
		},
		{
			"custom type",
			NewGormConfig().WithType("varchar(100)"),
			"type:varchar(100)",
		},
		{
			"not null",
			NewGormConfig().WithNotNull(),
			"not null",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tag := tt.config.ToGormTag()
			assert.Contains(t, tag, tt.wantPart)
		})
	}
}

func TestSchemaGormFluentAPI(t *testing.T) {
	// Test basic Gorm config assignment
	field := Text("Email").
		Gorm(NewGormConfig().WithUniqueIndex())

	schema := field.(*Schema)
	assert.True(t, schema.HasGormConfig())
	assert.True(t, schema.GetGormConfig().UniqueIndex)

	// Test convenience methods
	field2 := Text("Name").GormIndex("idx_name")
	schema2 := field2.(*Schema)
	schema2.GormSize(100)

	assert.True(t, schema2.HasGormConfig())
	assert.True(t, schema2.GetGormConfig().Index)
	assert.Equal(t, "idx_name", schema2.GetGormConfig().IndexName)
	assert.Equal(t, 100, schema2.GetGormConfig().Size)
}
