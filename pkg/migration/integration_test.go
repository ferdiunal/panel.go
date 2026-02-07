package migration_test

import (
	"testing"
	"time"

	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/core"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/migration"
	"github.com/ferdiunal/panel.go/pkg/resource"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Test Models
type TestUser struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"type:varchar(255);not null;index"`
	Email     string         `json:"email" gorm:"type:varchar(255);not null;uniqueIndex"`
	Age       int            `json:"age" gorm:"type:int"`
	IsActive  bool           `json:"isActive" gorm:"type:boolean;default:true"`
	Profile   *TestProfile   `json:"profile" gorm:"foreignKey:UserID"`
	Posts     []TestPost     `json:"posts" gorm:"foreignKey:AuthorID"`
	Roles     []TestRole     `json:"roles" gorm:"many2many:user_roles;"`
	CreatedAt time.Time      `json:"createdAt" gorm:"index"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

type TestProfile struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	UserID    *uint          `json:"userId" gorm:"index"`
	Bio       string         `json:"bio" gorm:"type:text"`
	Website   string         `json:"website" gorm:"type:varchar(255)"`
	User      *TestUser      `json:"user"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

type TestPost struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Title     string         `json:"title" gorm:"type:varchar(255);not null;index"`
	Content   string         `json:"content" gorm:"type:text"`
	AuthorID  uint           `json:"authorId" gorm:"index"`
	Author    *TestUser      `json:"author"`
	Tags      []TestTag      `json:"tags" gorm:"many2many:post_tags;"`
	Comments  []TestComment  `json:"comments" gorm:"polymorphic:Commentable;"`
	CreatedAt time.Time      `json:"createdAt" gorm:"index"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

type TestTag struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"type:varchar(100);not null;uniqueIndex"`
	Posts     []TestPost     `json:"posts" gorm:"many2many:post_tags;"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

type TestRole struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"type:varchar(100);not null;uniqueIndex"`
	Users     []TestUser     `json:"users" gorm:"many2many:user_roles;"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

type TestComment struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	Content         string         `json:"content" gorm:"type:text;not null"`
	CommentableID   uint           `json:"commentableId" gorm:"index"`
	CommentableType string         `json:"commentableType" gorm:"type:varchar(100);index"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

// Test Resources
type TestUserResource struct {
	resource.OptimizedBase
}

func NewTestUserResource() *TestUserResource {
	r := &TestUserResource{}
	r.SetModel(&TestUser{})
	r.SetSlug("users")
	r.SetTitle("Users")
	r.SetFieldResolver(&TestUserFieldResolver{})
	return r
}

type TestUserFieldResolver struct{}

func (r *TestUserFieldResolver) ResolveFields(ctx *context.Context) []core.Element {
	return []core.Element{
		fields.ID("ID").Sortable(),
		fields.Text("Name", "name").Required().Sortable().Searchable(),
		fields.Email("Email", "email").Required().Searchable(),
		fields.Number("Age", "age"),
		fields.NewHasOne("Profile", "profile", "profiles").ForeignKey("user_id"),
		fields.NewHasMany("Posts", "posts", "posts").ForeignKey("author_id"),
		fields.NewBelongsToMany("Roles", "user_roles", "roles").
			PivotTable("user_roles").
			ForeignKey("user_id").
			RelatedKey("role_id"),
		fields.DateTime("Created At", "createdAt").ReadOnly(),
	}
}

type TestProfileResource struct {
	resource.OptimizedBase
}

func NewTestProfileResource() *TestProfileResource {
	r := &TestProfileResource{}
	r.SetModel(&TestProfile{})
	r.SetSlug("profiles")
	r.SetTitle("Profiles")
	r.SetFieldResolver(&TestProfileFieldResolver{})
	return r
}

type TestProfileFieldResolver struct{}

func (r *TestProfileFieldResolver) ResolveFields(ctx *context.Context) []core.Element {
	return []core.Element{
		fields.ID("ID").Sortable(),
		fields.Text("Bio", "bio"),
		fields.Text("Website", "website"),
		fields.NewBelongsTo("User", "user_id", "users").DisplayUsing("name"),
		fields.DateTime("Created At", "createdAt").ReadOnly(),
	}
}

type TestTagResource struct {
	resource.OptimizedBase
}

func NewTestTagResource() *TestTagResource {
	r := &TestTagResource{}
	r.SetModel(&TestTag{})
	r.SetSlug("tags")
	r.SetTitle("Tags")
	r.SetFieldResolver(&TestTagFieldResolver{})
	return r
}

type TestTagFieldResolver struct{}

func (r *TestTagFieldResolver) ResolveFields(ctx *context.Context) []core.Element {
	return []core.Element{
		fields.ID("ID").Sortable(),
		fields.Text("Name", "name").Required().Searchable(),
		fields.NewBelongsToMany("Posts", "post_tags", "posts").
			PivotTable("post_tags").
			ForeignKey("tag_id").
			RelatedKey("post_id"),
		fields.DateTime("Created At", "createdAt").ReadOnly(),
	}
}

type TestRoleResource struct {
	resource.OptimizedBase
}

func NewTestRoleResource() *TestRoleResource {
	r := &TestRoleResource{}
	r.SetModel(&TestRole{})
	r.SetSlug("roles")
	r.SetTitle("Roles")
	r.SetFieldResolver(&TestRoleFieldResolver{})
	return r
}

type TestRoleFieldResolver struct{}

func (r *TestRoleFieldResolver) ResolveFields(ctx *context.Context) []core.Element {
	return []core.Element{
		fields.ID("ID").Sortable(),
		fields.Text("Name", "name").Required().Searchable(),
		fields.NewBelongsToMany("Users", "user_roles", "users").
			PivotTable("user_roles").
			ForeignKey("role_id").
			RelatedKey("user_id"),
		fields.DateTime("Created At", "createdAt").ReadOnly(),
	}
}

type TestPostResource struct {
	resource.OptimizedBase
}

func NewTestPostResource() *TestPostResource {
	r := &TestPostResource{}
	r.SetModel(&TestPost{})
	r.SetSlug("posts")
	r.SetTitle("Posts")
	r.SetFieldResolver(&TestPostFieldResolver{})
	return r
}

type TestPostFieldResolver struct{}

func (r *TestPostFieldResolver) ResolveFields(ctx *context.Context) []core.Element {
	return []core.Element{
		fields.ID("ID").Sortable(),
		fields.Text("Title", "title").Required().Sortable().Searchable(),
		fields.Text("Content", "content").Required(),
		fields.NewBelongsTo("Author", "author_id", "users").
			DisplayUsing("name").
			WithSearchableColumns("name", "email"),
		fields.NewBelongsToMany("Tags", "post_tags", "tags").
			PivotTable("post_tags").
			ForeignKey("post_id").
			RelatedKey("tag_id"),
		fields.NewMorphTo("Commentable", "commentable").
			Types(map[string]string{
				"posts": "posts",
			}).
			Displays(map[string]string{
				"posts": "title",
			}),
		fields.DateTime("Created At", "createdAt").ReadOnly(),
	}
}

type TestCommentResource struct {
	resource.OptimizedBase
}

func NewTestCommentResource() *TestCommentResource {
	r := &TestCommentResource{}
	r.SetModel(&TestComment{})
	r.SetSlug("comments")
	r.SetTitle("Comments")
	r.SetFieldResolver(&TestCommentFieldResolver{})
	return r
}

type TestCommentFieldResolver struct{}

func (r *TestCommentFieldResolver) ResolveFields(ctx *context.Context) []core.Element {
	return []core.Element{
		fields.ID("ID").Sortable(),
		fields.Text("Content", "content").Required(),
		fields.NewMorphTo("Commentable", "commentable").
			Types(map[string]string{
				"posts": "posts",
			}).
			Displays(map[string]string{
				"posts": "title",
			}),
		fields.DateTime("Created At", "createdAt").ReadOnly(),
	}
}

// Helper function to setup test database
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err, "Failed to connect to test database")

	return db
}

// Helper function to setup migration generator
func setupMigrationGenerator(t *testing.T, db *gorm.DB) *migration.MigrationGenerator {
	mg := migration.NewMigrationGenerator(db)

	// Register all test resources
	mg.RegisterResources(
		NewTestUserResource(),
		NewTestProfileResource(),
		NewTestPostResource(),
		NewTestTagResource(),
		NewTestRoleResource(),
		NewTestCommentResource(),
	)

	return mg
}

// Test: AutoMigrate creates tables correctly
func TestAutoMigrate_CreatesTablesCorrectly(t *testing.T) {
	db := setupTestDB(t)
	mg := setupMigrationGenerator(t, db)

	// Run AutoMigrate
	err := mg.AutoMigrate()
	require.NoError(t, err, "AutoMigrate should not return error")

	// Verify tables exist
	tables := []string{"test_users", "test_profiles", "test_posts", "test_tags", "test_roles", "test_comments"}
	for _, table := range tables {
		var count int64
		err := db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&count).Error
		require.NoError(t, err)
		assert.Equal(t, int64(1), count, "Table %s should exist", table)
	}

	// Verify pivot tables exist
	pivotTables := []string{"user_roles", "post_tags"}
	for _, table := range pivotTables {
		var count int64
		err := db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?", table).Scan(&count).Error
		require.NoError(t, err)
		assert.Equal(t, int64(1), count, "Pivot table %s should exist", table)
	}
}

// Test: CRUD operations work correctly
func TestCRUD_Operations(t *testing.T) {
	db := setupTestDB(t)
	mg := setupMigrationGenerator(t, db)

	// Run AutoMigrate
	err := mg.AutoMigrate()
	require.NoError(t, err)

	t.Run("Create", func(t *testing.T) {
		user := &TestUser{
			Name:     "John Doe",
			Email:    "john@example.com",
			Age:      30,
			IsActive: true,
		}

		err := db.Create(user).Error
		require.NoError(t, err)
		assert.NotZero(t, user.ID, "User ID should be set after create")
	})

	t.Run("Read", func(t *testing.T) {
		var user TestUser
		err := db.Where("email = ?", "john@example.com").First(&user).Error
		require.NoError(t, err)
		assert.Equal(t, "John Doe", user.Name)
		assert.Equal(t, "john@example.com", user.Email)
		assert.Equal(t, 30, user.Age)
		assert.True(t, user.IsActive)
	})

	t.Run("Update", func(t *testing.T) {
		var user TestUser
		err := db.Where("email = ?", "john@example.com").First(&user).Error
		require.NoError(t, err)

		user.Name = "Jane Doe"
		user.Age = 25

		err = db.Save(&user).Error
		require.NoError(t, err)

		var updated TestUser
		err = db.First(&updated, user.ID).Error
		require.NoError(t, err)
		assert.Equal(t, "Jane Doe", updated.Name)
		assert.Equal(t, 25, updated.Age)
	})

	t.Run("Delete", func(t *testing.T) {
		var user TestUser
		err := db.Where("email = ?", "john@example.com").First(&user).Error
		require.NoError(t, err)

		err = db.Delete(&user).Error
		require.NoError(t, err)

		var count int64
		err = db.Model(&TestUser{}).Where("email = ?", "john@example.com").Count(&count).Error
		require.NoError(t, err)
		assert.Equal(t, int64(0), count, "User should be deleted")
	})
}

// Test: BelongsTo relationship works correctly
func TestBelongsTo_Relationship(t *testing.T) {
	db := setupTestDB(t)
	mg := setupMigrationGenerator(t, db)

	err := mg.AutoMigrate()
	require.NoError(t, err)

	// Create user
	user := &TestUser{
		Name:  "Author",
		Email: "author@example.com",
	}
	err = db.Create(user).Error
	require.NoError(t, err)

	// Create post with author
	post := &TestPost{
		Title:    "Test Post",
		Content:  "Test Content",
		AuthorID: user.ID,
	}
	err = db.Create(post).Error
	require.NoError(t, err)

	// Load post with author
	var loadedPost TestPost
	err = db.Preload("Author").First(&loadedPost, post.ID).Error
	require.NoError(t, err)
	assert.NotNil(t, loadedPost.Author)
	assert.Equal(t, "Author", loadedPost.Author.Name)
}

// Test: HasOne relationship works correctly
func TestHasOne_Relationship(t *testing.T) {
	db := setupTestDB(t)
	mg := setupMigrationGenerator(t, db)

	err := mg.AutoMigrate()
	require.NoError(t, err)

	// Create user
	user := &TestUser{
		Name:  "User",
		Email: "user@example.com",
	}
	err = db.Create(user).Error
	require.NoError(t, err)

	// Create profile
	profile := &TestProfile{
		UserID:  &user.ID,
		Bio:     "Test Bio",
		Website: "https://example.com",
	}
	err = db.Create(profile).Error
	require.NoError(t, err)

	// Load user with profile
	var loadedUser TestUser
	err = db.Preload("Profile").First(&loadedUser, user.ID).Error
	require.NoError(t, err)
	assert.NotNil(t, loadedUser.Profile)
	assert.Equal(t, "Test Bio", loadedUser.Profile.Bio)
}

// Test: HasMany relationship works correctly
func TestHasMany_Relationship(t *testing.T) {
	db := setupTestDB(t)
	mg := setupMigrationGenerator(t, db)

	err := mg.AutoMigrate()
	require.NoError(t, err)

	// Create user
	user := &TestUser{
		Name:  "Author",
		Email: "author@example.com",
	}
	err = db.Create(user).Error
	require.NoError(t, err)

	// Create posts
	posts := []TestPost{
		{Title: "Post 1", Content: "Content 1", AuthorID: user.ID},
		{Title: "Post 2", Content: "Content 2", AuthorID: user.ID},
	}
	err = db.Create(&posts).Error
	require.NoError(t, err)

	// Load user with posts
	var loadedUser TestUser
	err = db.Preload("Posts").First(&loadedUser, user.ID).Error
	require.NoError(t, err)
	assert.Len(t, loadedUser.Posts, 2)
}

// Test: BelongsToMany relationship works correctly
func TestBelongsToMany_Relationship(t *testing.T) {
	db := setupTestDB(t)
	mg := setupMigrationGenerator(t, db)

	err := mg.AutoMigrate()
	require.NoError(t, err)

	// Create user
	user := &TestUser{
		Name:  "User",
		Email: "user@example.com",
	}
	err = db.Create(user).Error
	require.NoError(t, err)

	// Create roles
	roles := []TestRole{
		{Name: "Admin"},
		{Name: "Editor"},
	}
	err = db.Create(&roles).Error
	require.NoError(t, err)

	// Associate roles with user
	err = db.Model(&user).Association("Roles").Append(&roles)
	require.NoError(t, err)

	// Load user with roles
	var loadedUser TestUser
	err = db.Preload("Roles").First(&loadedUser, user.ID).Error
	require.NoError(t, err)
	assert.Len(t, loadedUser.Roles, 2)
}

// Test: MorphTo relationship works correctly
func TestMorphTo_Relationship(t *testing.T) {
	db := setupTestDB(t)
	mg := setupMigrationGenerator(t, db)

	err := mg.AutoMigrate()
	require.NoError(t, err)

	// Create post
	post := &TestPost{
		Title:   "Test Post",
		Content: "Test Content",
	}
	err = db.Create(post).Error
	require.NoError(t, err)

	// Create comment
	comment := &TestComment{
		Content:         "Test Comment",
		CommentableID:   post.ID,
		CommentableType: "posts",
	}
	err = db.Create(comment).Error
	require.NoError(t, err)

	// Load comment
	var loadedComment TestComment
	err = db.First(&loadedComment, comment.ID).Error
	require.NoError(t, err)
	assert.Equal(t, post.ID, loadedComment.CommentableID)
	assert.Equal(t, "posts", loadedComment.CommentableType)
}

// Test: Indexes are created correctly
func TestIndexes_CreatedCorrectly(t *testing.T) {
	db := setupTestDB(t)
	mg := setupMigrationGenerator(t, db)

	err := mg.AutoMigrate()
	require.NoError(t, err)

	// Check if indexes exist
	var count int64
	err = db.Raw("SELECT COUNT(*) FROM sqlite_master WHERE type='index' AND name LIKE 'idx_%'").Scan(&count).Error
	require.NoError(t, err)
	assert.Greater(t, count, int64(0), "Indexes should be created")
}

// Test: Unique constraints work correctly
func TestUniqueConstraints_WorkCorrectly(t *testing.T) {
	db := setupTestDB(t)
	mg := setupMigrationGenerator(t, db)

	err := mg.AutoMigrate()
	require.NoError(t, err)

	// Create first user
	user1 := &TestUser{
		Name:  "User 1",
		Email: "test@example.com",
	}
	err = db.Create(user1).Error
	require.NoError(t, err)

	// Try to create second user with same email
	user2 := &TestUser{
		Name:  "User 2",
		Email: "test@example.com",
	}
	err = db.Create(user2).Error
	assert.Error(t, err, "Should fail due to unique constraint")
}

// Test: Soft delete works correctly
func TestSoftDelete_WorksCorrectly(t *testing.T) {
	db := setupTestDB(t)
	mg := setupMigrationGenerator(t, db)

	err := mg.AutoMigrate()
	require.NoError(t, err)

	// Create user
	user := &TestUser{
		Name:  "User",
		Email: "user@example.com",
	}
	err = db.Create(user).Error
	require.NoError(t, err)

	// Soft delete
	err = db.Delete(&user).Error
	require.NoError(t, err)

	// Verify user is soft deleted
	var count int64
	err = db.Model(&TestUser{}).Where("id = ?", user.ID).Count(&count).Error
	require.NoError(t, err)
	assert.Equal(t, int64(0), count, "User should not be found in normal query")

	// Verify user exists with Unscoped
	err = db.Unscoped().Model(&TestUser{}).Where("id = ?", user.ID).Count(&count).Error
	require.NoError(t, err)
	assert.Equal(t, int64(1), count, "User should be found with Unscoped")
}

// Test: Migration doesn't break existing database
func TestMigration_DoesntBreakExistingDatabase(t *testing.T) {
	db := setupTestDB(t)
	mg := setupMigrationGenerator(t, db)

	// First migration
	err := mg.AutoMigrate()
	require.NoError(t, err)

	// Create test data
	user := &TestUser{
		Name:  "User",
		Email: "user@example.com",
	}
	err = db.Create(user).Error
	require.NoError(t, err)

	// Second migration (should be idempotent)
	err = mg.AutoMigrate()
	require.NoError(t, err)

	// Verify data still exists
	var loadedUser TestUser
	err = db.First(&loadedUser, user.ID).Error
	require.NoError(t, err)
	assert.Equal(t, "User", loadedUser.Name)
}
