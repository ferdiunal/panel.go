package main

import (
	"testing"

	"github.com/ferdiunal/panel.go/examples/simple/blog"
	"github.com/ferdiunal/panel.go/examples/simple/products"
	"github.com/ferdiunal/panel.go/pkg/migration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err, "Failed to connect to test database")

	return db
}

// Test: AutoMigrate with blog resources (model-based)
func TestBlogResources_AutoMigrate(t *testing.T) {
	db := setupTestDB(t)
	mg := migration.NewMigrationGenerator(db)

	// Register all blog resources (all have models)
	mg.RegisterResources(
		blog.NewAuthorResource(),
		blog.NewProfileResource(),
		blog.NewPostResource(),
		blog.NewTagResource(),
		blog.NewCommentResource(),
	)

	// Run AutoMigrate
	err := mg.AutoMigrate()
	require.NoError(t, err, "AutoMigrate should not return error")

	// Verify tables exist
	tables := []string{"authors", "profiles", "posts", "tags", "comments"}
	for _, table := range tables {
		assert.True(t, db.Migrator().HasTable(table), "Table %s should exist", table)
	}

	// Verify pivot table exists
	assert.True(t, db.Migrator().HasTable("post_tags"), "Pivot table post_tags should exist")
}

// Test: Product resource with model
func TestProducts_ModelBasedMigration(t *testing.T) {
	db := setupTestDB(t)
	mg := migration.NewMigrationGenerator(db)

	// Register product resource (has model)
	mg.RegisterResource(products.NewProductResource())

	// Run AutoMigrate
	err := mg.AutoMigrate()
	require.NoError(t, err, "AutoMigrate should work with model")

	// Verify table exists
	assert.True(t, db.Migrator().HasTable("products"), "Products table should exist")

	// Verify columns exist using GORM Migrator
	assert.True(t, db.Migrator().HasColumn(&products.Product{}, "name"), "Column name should exist")
	assert.True(t, db.Migrator().HasColumn(&products.Product{}, "description"), "Column description should exist")
	assert.True(t, db.Migrator().HasColumn(&products.Product{}, "details"), "Column details should exist")
	assert.True(t, db.Migrator().HasColumn(&products.Product{}, "price"), "Column price should exist")
	assert.True(t, db.Migrator().HasColumn(&products.Product{}, "stock"), "Column stock should exist")

	// Test CRUD operations with GORM
	t.Run("Create Product", func(t *testing.T) {
		product := &products.Product{
			Name:        "Test Product",
			Description: "Test Description",
			Details:     "<p>Rich text content</p>",
			Price:       99.99,
			Stock:       10,
		}
		err := db.Create(product).Error
		require.NoError(t, err)
		assert.NotZero(t, product.ID)
	})

	t.Run("Read Product", func(t *testing.T) {
		var product products.Product
		err := db.Where("name = ?", "Test Product").First(&product).Error
		require.NoError(t, err)
		assert.Equal(t, "Test Product", product.Name)
		assert.Equal(t, 99.99, product.Price)
	})

	t.Run("Update Product", func(t *testing.T) {
		var product products.Product
		err := db.Where("name = ?", "Test Product").First(&product).Error
		require.NoError(t, err)

		product.Price = 89.99
		err = db.Save(&product).Error
		require.NoError(t, err)

		var updated products.Product
		err = db.First(&updated, product.ID).Error
		require.NoError(t, err)
		assert.Equal(t, 89.99, updated.Price)
	})

	t.Run("Soft Delete Product", func(t *testing.T) {
		var product products.Product
		err := db.Where("name = ?", "Test Product").First(&product).Error
		require.NoError(t, err)

		err = db.Delete(&product).Error
		require.NoError(t, err)

		// Should not find with normal query (soft deleted)
		var count int64
		err = db.Model(&products.Product{}).Where("name = ?", "Test Product").Count(&count).Error
		require.NoError(t, err)
		assert.Equal(t, int64(0), count, "Product should be soft deleted")

		// Should find with Unscoped
		err = db.Unscoped().Model(&products.Product{}).Where("name = ?", "Test Product").Count(&count).Error
		require.NoError(t, err)
		assert.Equal(t, int64(1), count, "Product should exist with Unscoped")
	})
}

// Test: Idempotent migration
func TestIdempotentMigration(t *testing.T) {
	db := setupTestDB(t)
	mg := migration.NewMigrationGenerator(db)

	mg.RegisterResource(products.NewProductResource())

	// First migration
	err := mg.AutoMigrate()
	require.NoError(t, err)

	// Create test data
	product := &products.Product{
		Name:  "Test Product",
		Price: 99.99,
	}
	err = db.Create(product).Error
	require.NoError(t, err)

	// Second migration (should be idempotent)
	err = mg.AutoMigrate()
	require.NoError(t, err)

	// Verify data still exists
	var loadedProduct products.Product
	err = db.First(&loadedProduct, product.ID).Error
	require.NoError(t, err)
	assert.Equal(t, "Test Product", loadedProduct.Name)
}

// Test: Resource without model should fail
func TestResourceWithoutModel_ShouldFail(t *testing.T) {
	db := setupTestDB(t)
	mg := migration.NewMigrationGenerator(db)

	// Create a resource without model
	type NoModelResource struct {
		products.ProductResource
	}
	noModelResource := &NoModelResource{}
	noModelResource.SetSlug("no_model")
	noModelResource.SetTitle("No Model")

	mg.RegisterResource(noModelResource)

	// Should fail because no model
	err := mg.AutoMigrate()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "has no model")
}
