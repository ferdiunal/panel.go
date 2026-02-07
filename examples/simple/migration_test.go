package main

import (
	"strings"
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

// Test: AutoMigrate with real blog resources
func TestBlogResources_AutoMigrate(t *testing.T) {
	db := setupTestDB(t)
	mg := migration.NewMigrationGenerator(db)

	// Register all blog resources
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
		var count int64
		err := db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&count).Error
		require.NoError(t, err)
		assert.Equal(t, int64(1), count, "Table %s should exist", table)
	}

	// Verify pivot table exists
	var count int64
	err = db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='post_tags'").Scan(&count).Error
	require.NoError(t, err)
	assert.Equal(t, int64(1), count, "Pivot table post_tags should exist")
}

// Test: CRUD operations with Author resource
func TestAuthor_CRUD(t *testing.T) {
	db := setupTestDB(t)
	mg := migration.NewMigrationGenerator(db)

	mg.RegisterResources(
		blog.NewAuthorResource(),
		blog.NewProfileResource(),
	)

	err := mg.AutoMigrate()
	require.NoError(t, err)

	t.Run("Create Author", func(t *testing.T) {
		author := &blog.Author{
			Name:  "John Doe",
			Email: "john@example.com",
		}

		err := db.Create(author).Error
		require.NoError(t, err)
		assert.NotZero(t, author.ID, "Author ID should be set")
	})

	t.Run("Read Author", func(t *testing.T) {
		var author blog.Author
		err := db.Where("email = ?", "john@example.com").First(&author).Error
		require.NoError(t, err)
		assert.Equal(t, "John Doe", author.Name)
	})

	t.Run("Update Author", func(t *testing.T) {
		var author blog.Author
		err := db.Where("email = ?", "john@example.com").First(&author).Error
		require.NoError(t, err)

		author.Name = "Jane Doe"
		err = db.Save(&author).Error
		require.NoError(t, err)

		var updated blog.Author
		err = db.First(&updated, author.ID).Error
		require.NoError(t, err)
		assert.Equal(t, "Jane Doe", updated.Name)
	})

	t.Run("Delete Author", func(t *testing.T) {
		var author blog.Author
		err := db.Where("email = ?", "john@example.com").First(&author).Error
		require.NoError(t, err)

		err = db.Delete(&author).Error
		require.NoError(t, err)

		var count int64
		err = db.Model(&blog.Author{}).Where("email = ?", "john@example.com").Count(&count).Error
		require.NoError(t, err)
		assert.Equal(t, int64(0), count, "Author should be deleted")
	})
}

// Test: Post with Author relationship
func TestPost_WithAuthor(t *testing.T) {
	db := setupTestDB(t)
	mg := migration.NewMigrationGenerator(db)

	mg.RegisterResources(
		blog.NewAuthorResource(),
		blog.NewPostResource(),
	)

	err := mg.AutoMigrate()
	require.NoError(t, err)

	// Create author
	author := &blog.Author{
		Name:  "Author",
		Email: "author@example.com",
	}
	err = db.Create(author).Error
	require.NoError(t, err)

	// Create post
	post := &blog.Post{
		Title:    "Test Post",
		Content:  "Test Content",
		AuthorID: author.ID,
	}
	err = db.Create(post).Error
	require.NoError(t, err)

	// Load post with author
	var loadedPost blog.Post
	err = db.Preload("Author").First(&loadedPost, post.ID).Error
	require.NoError(t, err)
	assert.NotNil(t, loadedPost.Author)
	assert.Equal(t, "Author", loadedPost.Author.Name)
}

// Test: Author with Profile (HasOne)
func TestAuthor_WithProfile(t *testing.T) {
	db := setupTestDB(t)
	mg := migration.NewMigrationGenerator(db)

	mg.RegisterResources(
		blog.NewAuthorResource(),
		blog.NewProfileResource(),
	)

	err := mg.AutoMigrate()
	require.NoError(t, err)

	// Create author
	author := &blog.Author{
		Name:  "Author",
		Email: "author@example.com",
	}
	err = db.Create(author).Error
	require.NoError(t, err)

	// Create profile
	profile := &blog.Profile{
		AuthorID: &author.ID,
		Bio:      "Test Bio",
		Website:  "https://example.com",
	}
	err = db.Create(profile).Error
	require.NoError(t, err)

	// Load author with profile
	var loadedAuthor blog.Author
	err = db.Preload("Profile").First(&loadedAuthor, author.ID).Error
	require.NoError(t, err)
	assert.NotNil(t, loadedAuthor.Profile)
	assert.Equal(t, "Test Bio", loadedAuthor.Profile.Bio)
}

// Test: Post with Tags (BelongsToMany)
func TestPost_WithTags(t *testing.T) {
	db := setupTestDB(t)
	mg := migration.NewMigrationGenerator(db)

	mg.RegisterResources(
		blog.NewPostResource(),
		blog.NewTagResource(),
	)

	err := mg.AutoMigrate()
	require.NoError(t, err)

	// Create post
	post := &blog.Post{
		Title:   "Test Post",
		Content: "Test Content",
	}
	err = db.Create(post).Error
	require.NoError(t, err)

	// Create tags
	tags := []blog.Tag{
		{Name: "Go"},
		{Name: "GORM"},
	}
	err = db.Create(&tags).Error
	require.NoError(t, err)

	// Associate tags with post
	err = db.Model(&post).Association("Tags").Append(&tags)
	require.NoError(t, err)

	// Load post with tags
	var loadedPost blog.Post
	err = db.Preload("Tags").First(&loadedPost, post.ID).Error
	require.NoError(t, err)
	assert.Len(t, loadedPost.Tags, 2)
}

// Test: Comment with MorphTo relationship
func TestComment_WithMorphTo(t *testing.T) {
	db := setupTestDB(t)
	mg := migration.NewMigrationGenerator(db)

	mg.RegisterResources(
		blog.NewPostResource(),
		blog.NewCommentResource(),
	)

	err := mg.AutoMigrate()
	require.NoError(t, err)

	// Create post
	post := &blog.Post{
		Title:   "Test Post",
		Content: "Test Content",
	}
	err = db.Create(post).Error
	require.NoError(t, err)

	// Create comment
	comment := &blog.Comment{
		Content:         "Test Comment",
		CommentableID:   post.ID,
		CommentableType: "posts",
	}
	err = db.Create(comment).Error
	require.NoError(t, err)

	// Verify comment
	var loadedComment blog.Comment
	err = db.First(&loadedComment, comment.ID).Error
	require.NoError(t, err)
	assert.Equal(t, post.ID, loadedComment.CommentableID)
	assert.Equal(t, "posts", loadedComment.CommentableType)
}

// Test: Indexes are created for searchable fields
func TestBlogResources_Indexes(t *testing.T) {
	db := setupTestDB(t)
	mg := migration.NewMigrationGenerator(db)

	mg.RegisterResources(
		blog.NewAuthorResource(),
		blog.NewPostResource(),
	)

	err := mg.AutoMigrate()
	require.NoError(t, err)

	// Check if indexes exist
	var count int64
	err = db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='index' AND name LIKE 'idx_%'").Scan(&count).Error
	require.NoError(t, err)
	assert.Greater(t, count, int64(0), "Indexes should be created for searchable fields")
}

// Test: Soft delete works correctly
func TestBlogResources_SoftDelete(t *testing.T) {
	db := setupTestDB(t)
	mg := migration.NewMigrationGenerator(db)

	mg.RegisterResources(blog.NewAuthorResource())

	err := mg.AutoMigrate()
	require.NoError(t, err)

	// Create author
	author := &blog.Author{
		Name:  "Author",
		Email: "author@example.com",
	}
	err = db.Create(author).Error
	require.NoError(t, err)

	// Soft delete
	err = db.Delete(&author).Error
	require.NoError(t, err)

	// Verify author is soft deleted
	var count int64
	err = db.Model(&blog.Author{}).Where("id = ?", author.ID).Count(&count).Error
	require.NoError(t, err)
	assert.Equal(t, int64(0), count, "Author should not be found in normal query")

	// Verify author exists with Unscoped
	err = db.Unscoped().Model(&blog.Author{}).Where("id = ?", author.ID).Count(&count).Error
	require.NoError(t, err)
	assert.Equal(t, int64(1), count, "Author should be found with Unscoped")
}

// Test: Migration is idempotent (doesn't break existing data)
func TestBlogResources_IdempotentMigration(t *testing.T) {
	db := setupTestDB(t)
	mg := migration.NewMigrationGenerator(db)

	mg.RegisterResources(blog.NewAuthorResource())

	// First migration
	err := mg.AutoMigrate()
	require.NoError(t, err)

	// Create test data
	author := &blog.Author{
		Name:  "Author",
		Email: "author@example.com",
	}
	err = db.Create(author).Error
	require.NoError(t, err)

	// Second migration (should be idempotent)
	err = mg.AutoMigrate()
	require.NoError(t, err)

	// Verify data still exists
	var loadedAuthor blog.Author
	err = db.First(&loadedAuthor, author.ID).Error
	require.NoError(t, err)
	assert.Equal(t, "Author", loadedAuthor.Name)
}

// Test: Model-less resource migration (Products)
func TestProducts_ModellessMigration(t *testing.T) {
	db := setupTestDB(t)
	mg := migration.NewMigrationGenerator(db)

	// Register product resource (no model, only fields)
	mg.RegisterResource(products.NewProductResource())

	// Run AutoMigrate - should use createTableFromFields
	err := mg.AutoMigrate()
	require.NoError(t, err, "AutoMigrate should work without model")

	// Verify table exists
	var count int64
	err = db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='products'").Scan(&count).Error
	require.NoError(t, err)
	assert.Equal(t, int64(1), count, "Products table should exist")

	// Verify columns exist
	type ColumnInfo struct {
		Name string
		Type string
	}
	var columns []ColumnInfo
	err = db.Raw("PRAGMA table_info(products)").Scan(&columns).Error
	require.NoError(t, err)

	columnNames := make(map[string]bool)
	for _, col := range columns {
		columnNames[col.Name] = true
	}

	// Check expected columns
	expectedColumns := []string{"id", "name", "description", "details", "price", "stock", "created_at", "updated_at", "deleted_at"}
	for _, expected := range expectedColumns {
		assert.True(t, columnNames[expected], "Column %s should exist", expected)
	}

	// Verify SQL types
	columnTypes := make(map[string]string)
	for _, col := range columns {
		columnTypes[col.Name] = strings.ToLower(col.Type)
	}

	// Text field → VARCHAR
	assert.Contains(t, columnTypes["name"], "varchar", "name should be VARCHAR")
	// Textarea field → TEXT
	assert.Equal(t, "text", columnTypes["description"], "description should be TEXT")
	// RichText field → TEXT
	assert.Equal(t, "text", columnTypes["details"], "details should be TEXT")

	// Test CRUD operations
	t.Run("Create Product", func(t *testing.T) {
		result := db.Exec("INSERT INTO products (name, description, price, stock, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
			"Test Product", "Test Description", 99.99, 10, "2024-01-01 00:00:00", "2024-01-01 00:00:00")
		require.NoError(t, result.Error)
		assert.Equal(t, int64(1), result.RowsAffected)
	})

	t.Run("Read Product", func(t *testing.T) {
		var product struct {
			ID          uint
			Name        string
			Description string
			Price       float64
			Stock       int
		}
		err := db.Raw("SELECT id, name, description, price, stock FROM products WHERE name = ?", "Test Product").Scan(&product).Error
		require.NoError(t, err)
		assert.Equal(t, "Test Product", product.Name)
		assert.Equal(t, 99.99, product.Price)
	})

	t.Run("Update Product", func(t *testing.T) {
		result := db.Exec("UPDATE products SET price = ?, updated_at = ? WHERE name = ?", 89.99, "2024-01-02 00:00:00", "Test Product")
		require.NoError(t, result.Error)
		assert.Equal(t, int64(1), result.RowsAffected)

		var price float64
		err := db.Raw("SELECT price FROM products WHERE name = ?", "Test Product").Scan(&price).Error
		require.NoError(t, err)
		assert.Equal(t, 89.99, price)
	})

	t.Run("Delete Product", func(t *testing.T) {
		result := db.Exec("DELETE FROM products WHERE name = ?", "Test Product")
		require.NoError(t, result.Error)
		assert.Equal(t, int64(1), result.RowsAffected)

		var count int64
		err := db.Raw("SELECT COUNT(*) FROM products WHERE name = ?", "Test Product").Scan(&count).Error
		require.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})
}

// Test: Idempotent migration for model-less resources
func TestProducts_IdempotentModellessMigration(t *testing.T) {
	db := setupTestDB(t)
	mg := migration.NewMigrationGenerator(db)

	mg.RegisterResource(products.NewProductResource())

	// First migration
	err := mg.AutoMigrate()
	require.NoError(t, err)

	// Insert test data
	result := db.Exec("INSERT INTO products (name, price, created_at, updated_at) VALUES (?, ?, ?, ?)",
		"Test Product", 99.99, "2024-01-01 00:00:00", "2024-01-01 00:00:00")
	require.NoError(t, result.Error)

	// Second migration (should be idempotent)
	err = mg.AutoMigrate()
	require.NoError(t, err)

	// Verify data still exists
	var count int64
	err = db.Raw("SELECT COUNT(*) FROM products WHERE name = ?", "Test Product").Scan(&count).Error
	require.NoError(t, err)
	assert.Equal(t, int64(1), count, "Data should still exist after second migration")
}

// Test: Mixed model and model-less resources
func TestMixed_ModelAndModellessResources(t *testing.T) {
	db := setupTestDB(t)
	mg := migration.NewMigrationGenerator(db)

	// Register both model-based and model-less resources
	mg.RegisterResources(
		blog.NewAuthorResource(),    // Has model
		products.NewProductResource(), // No model
	)

	err := mg.AutoMigrate()
	require.NoError(t, err, "Should handle mixed resources")

	// Verify both tables exist
	tables := []string{"authors", "products"}
	for _, table := range tables {
		var count int64
		err := db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&count).Error
		require.NoError(t, err)
		assert.Equal(t, int64(1), count, "Table %s should exist", table)
	}

	// Test data insertion in both tables
	t.Run("Insert into model-based table", func(t *testing.T) {
		author := &blog.Author{
			Name:  "Test Author",
			Email: "test@example.com",
		}
		err := db.Create(author).Error
		require.NoError(t, err)
		assert.NotZero(t, author.ID)
	})

	t.Run("Insert into model-less table", func(t *testing.T) {
		result := db.Exec("INSERT INTO products (name, price, created_at, updated_at) VALUES (?, ?, ?, ?)",
			"Test Product", 99.99, "2024-01-01 00:00:00", "2024-01-01 00:00:00")
		require.NoError(t, result.Error)
		assert.Equal(t, int64(1), result.RowsAffected)
	})
}

// Test: Complex scenario with all relationships
func TestBlogResources_ComplexScenario(t *testing.T) {
	db := setupTestDB(t)
	mg := migration.NewMigrationGenerator(db)

	// Register all resources
	mg.RegisterResources(
		blog.NewAuthorResource(),
		blog.NewProfileResource(),
		blog.NewPostResource(),
		blog.NewTagResource(),
		blog.NewCommentResource(),
	)

	err := mg.AutoMigrate()
	require.NoError(t, err)

	// Create author with profile
	author := &blog.Author{
		Name:  "John Doe",
		Email: "john@example.com",
	}
	err = db.Create(author).Error
	require.NoError(t, err)

	profile := &blog.Profile{
		AuthorID: &author.ID,
		Bio:      "Software Engineer",
		Website:  "https://johndoe.com",
	}
	err = db.Create(profile).Error
	require.NoError(t, err)

	// Create post with tags
	post := &blog.Post{
		Title:    "Getting Started with GORM",
		Content:  "GORM is a fantastic ORM library for Go...",
		AuthorID: author.ID,
	}
	err = db.Create(post).Error
	require.NoError(t, err)

	tags := []blog.Tag{
		{Name: "Go"},
		{Name: "GORM"},
		{Name: "Database"},
	}
	err = db.Create(&tags).Error
	require.NoError(t, err)

	err = db.Model(&post).Association("Tags").Append(&tags)
	require.NoError(t, err)

	// Create comment on post
	comment := &blog.Comment{
		Content:         "Great article!",
		CommentableID:   post.ID,
		CommentableType: "posts",
	}
	err = db.Create(comment).Error
	require.NoError(t, err)

	// Verify everything is connected
	var loadedAuthor blog.Author
	err = db.Preload("Profile").Preload("Posts.Tags").First(&loadedAuthor, author.ID).Error
	require.NoError(t, err)

	assert.Equal(t, "John Doe", loadedAuthor.Name)
	assert.NotNil(t, loadedAuthor.Profile)
	assert.Equal(t, "Software Engineer", loadedAuthor.Profile.Bio)
	assert.Len(t, loadedAuthor.Posts, 1)
	assert.Equal(t, "Getting Started with GORM", loadedAuthor.Posts[0].Title)
	assert.Len(t, loadedAuthor.Posts[0].Tags, 3)

	// Verify comment
	var loadedComment blog.Comment
	err = db.Where("commentable_id = ? AND commentable_type = ?", post.ID, "posts").First(&loadedComment).Error
	require.NoError(t, err)
	assert.Equal(t, "Great article!", loadedComment.Content)
}
