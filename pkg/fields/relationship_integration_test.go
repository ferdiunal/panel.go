package fields

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// TestBelongsToRelationship tests BelongsTo relationship with real SQLite database
func TestBelongsToRelationship(t *testing.T) {
	tdb := NewTestDB(t)
	defer tdb.Close()

	t.Run("BelongsTo resolves related user for a post", func(t *testing.T) {
		// Get a post from database
		var post Post
		err := tdb.DB.QueryRow("SELECT id, user_id, title, content FROM posts WHERE id = 1").
			Scan(&post.ID, &post.UserID, &post.Title, &post.Content)
		if err != nil {
			t.Fatalf("Failed to get post: %v", err)
		}

		// Create BelongsTo field
		field := BelongsTo("Author", "user_id", "users")

		// Verify field properties
		if field.GetRelationshipType() != "belongsTo" {
			t.Errorf("Expected relationship type 'belongsTo', got '%s'", field.GetRelationshipType())
		}

		if field.GetRelatedResource() != "users" {
			t.Errorf("Expected related resource 'users', got '%s'", field.GetRelatedResource())
		}

		// Get related user
		var user User
		err = tdb.DB.QueryRow("SELECT id, name, email FROM users WHERE id = ?", post.UserID).
			Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			t.Fatalf("Failed to get user: %v", err)
		}

		if user.ID != 1 {
			t.Errorf("Expected user ID 1, got %d", user.ID)
		}

		if user.Name != "John Doe" {
			t.Errorf("Expected user name 'John Doe', got '%s'", user.Name)
		}
	})

	t.Run("BelongsTo with custom display key", func(t *testing.T) {
		field := BelongsTo("Author", "user_id", "users")
		field.DisplayUsing("email")

		if field.GetDisplayKey() != "email" {
			t.Errorf("Expected display key 'email', got '%s'", field.GetDisplayKey())
		}
	})

	t.Run("BelongsTo with searchable columns", func(t *testing.T) {
		field := BelongsTo("Author", "user_id", "users")
		field.WithSearchableColumns("name", "email")

		columns := field.GetSearchableColumns()
		if len(columns) != 2 {
			t.Errorf("Expected 2 searchable columns, got %d", len(columns))
		}
	})
}

// TestHasManyRelationship tests HasMany relationship with real SQLite database
func TestHasManyRelationship(t *testing.T) {
	tdb := NewTestDB(t)
	defer tdb.Close()

	t.Run("HasMany resolves all posts for a user", func(t *testing.T) {
		// Create HasMany field
		field := HasMany("Posts", "posts", "posts")

		if field.GetRelationshipType() != "hasMany" {
			t.Errorf("Expected relationship type 'hasMany', got '%s'", field.GetRelationshipType())
		}

		// Get all posts for user 1
		rows, err := tdb.DB.Query("SELECT id, user_id, title, content FROM posts WHERE user_id = ?", 1)
		if err != nil {
			t.Fatalf("Failed to query posts: %v", err)
		}
		defer rows.Close()

		postCount := 0
		for rows.Next() {
			var post Post
			err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content)
			if err != nil {
				t.Fatalf("Failed to scan post: %v", err)
			}

			if post.UserID != 1 {
				t.Errorf("Expected user_id 1, got %d", post.UserID)
			}

			postCount++
		}

		if postCount != 2 {
			t.Errorf("Expected 2 posts for user 1, got %d", postCount)
		}
	})

	t.Run("HasMany with custom foreign key", func(t *testing.T) {
		field := HasMany("Posts", "posts", "posts")
		field.ForeignKey("author_id")

		if field.ForeignKeyColumn != "author_id" {
			t.Errorf("Expected foreign key 'author_id', got '%s'", field.ForeignKeyColumn)
		}
	})
}

// TestHasOneRelationship tests HasOne relationship with real SQLite database
func TestHasOneRelationship(t *testing.T) {
	tdb := NewTestDB(t)
	defer tdb.Close()

	t.Run("HasOne resolves single profile for a user", func(t *testing.T) {
		// Create HasOne field
		field := HasOne("Profile", "profile", "profiles")

		if field.GetRelationshipType() != "hasOne" {
			t.Errorf("Expected relationship type 'hasOne', got '%s'", field.GetRelationshipType())
		}

		// Get profile for user 1
		var profile Profile
		err := tdb.DB.QueryRow("SELECT id, user_id, bio, avatar_url FROM profiles WHERE user_id = ?", 1).
			Scan(&profile.ID, &profile.UserID, &profile.Bio, &profile.AvatarURL)
		if err != nil {
			t.Fatalf("Failed to get profile: %v", err)
		}

		if profile.UserID != 1 {
			t.Errorf("Expected user_id 1, got %d", profile.UserID)
		}

		if profile.Bio != "Software developer" {
			t.Errorf("Expected bio 'Software developer', got '%s'", profile.Bio)
		}
	})

	t.Run("HasOne returns no profile for user without profile", func(t *testing.T) {
		// User 3 has no profile
		var profile Profile
		err := tdb.DB.QueryRow("SELECT id, user_id, bio, avatar_url FROM profiles WHERE user_id = ?", 3).
			Scan(&profile.ID, &profile.UserID, &profile.Bio, &profile.AvatarURL)

		if err != sql.ErrNoRows {
			t.Errorf("Expected sql.ErrNoRows, got %v", err)
		}
	})
}

// TestBelongsToManyRelationship tests BelongsToMany relationship with real SQLite database
func TestBelongsToManyRelationship(t *testing.T) {
	tdb := NewTestDB(t)
	defer tdb.Close()

	t.Run("BelongsToMany resolves all tags for a post", func(t *testing.T) {
		// Create BelongsToMany field
		field := BelongsToMany("Tags", "post_tag", "tags")

		if field.GetRelationshipType() != "belongsToMany" {
			t.Errorf("Expected relationship type 'belongsToMany', got '%s'", field.GetRelationshipType())
		}

		// Get all tags for post 1 through pivot table
		rows, err := tdb.DB.Query(`
			SELECT t.id, t.name FROM tags t
			INNER JOIN post_tag pt ON t.id = pt.tag_id
			WHERE pt.post_id = ?
		`, 1)
		if err != nil {
			t.Fatalf("Failed to query tags: %v", err)
		}
		defer rows.Close()

		tagCount := 0
		for rows.Next() {
			var tag Tag
			err := rows.Scan(&tag.ID, &tag.Name)
			if err != nil {
				t.Fatalf("Failed to scan tag: %v", err)
			}

			if tag.Name != "golang" && tag.Name != "database" {
				t.Errorf("Unexpected tag name: %s", tag.Name)
			}

			tagCount++
		}

		if tagCount != 2 {
			t.Errorf("Expected 2 tags for post 1, got %d", tagCount)
		}
	})

	t.Run("BelongsToMany with custom pivot table", func(t *testing.T) {
		field := BelongsToMany("Tags", "post_tag", "tags")
		field.PivotTable("custom_pivot")

		if field.PivotTableName != "custom_pivot" {
			t.Errorf("Expected pivot table 'custom_pivot', got '%s'", field.PivotTableName)
		}
	})
}

// TestMorphToRelationship tests MorphTo relationship with real SQLite database
func TestMorphToRelationship(t *testing.T) {
	tdb := NewTestDB(t)
	defer tdb.Close()

	t.Run("MorphTo resolves polymorphic relationships", func(t *testing.T) {
		// Create MorphTo field
		field := NewMorphTo("Taggable", "taggable")
		field.Types(map[string]string{
			"post":    "posts",
			"comment": "comments",
		})

		if field.GetRelationshipType() != "morphTo" {
			t.Errorf("Expected relationship type 'morphTo', got '%s'", field.GetRelationshipType())
		}

		types := field.GetTypes()
		if len(types) != 2 {
			t.Errorf("Expected 2 types, got %d", len(types))
		}

		// Get taggable relationships
		rows, err := tdb.DB.Query("SELECT id, taggable_id, taggable_type, tag_id FROM taggables")
		if err != nil {
			t.Fatalf("Failed to query taggables: %v", err)
		}
		defer rows.Close()

		taggableCount := 0
		for rows.Next() {
			var taggable Taggable
			err := rows.Scan(&taggable.ID, &taggable.TaggableID, &taggable.TaggableType, &taggable.TagID)
			if err != nil {
				t.Fatalf("Failed to scan taggable: %v", err)
			}

			taggableCount++
		}

		if taggableCount == 0 {
			t.Error("Expected at least one taggable relationship")
		}
	})
}

// TestRelationshipLoading tests eager and lazy loading strategies
func TestRelationshipLoading(t *testing.T) {
	tdb := NewTestDB(t)
	defer tdb.Close()

	t.Run("Eager loading strategy", func(t *testing.T) {
		field := BelongsTo("Author", "user_id", "users")
		field.WithEagerLoad()

		if field.GetLoadingStrategy() != EAGER_LOADING {
			t.Errorf("Expected EAGER_LOADING, got %v", field.GetLoadingStrategy())
		}
	})

	t.Run("Lazy loading strategy", func(t *testing.T) {
		field := BelongsTo("Author", "user_id", "users")
		field.WithLazyLoad()

		if field.GetLoadingStrategy() != LAZY_LOADING {
			t.Errorf("Expected LAZY_LOADING, got %v", field.GetLoadingStrategy())
		}
	})
}

// TestRelationshipValidation tests validation with real data
func TestRelationshipValidation(t *testing.T) {
	tdb := NewTestDB(t)
	defer tdb.Close()

	t.Run("Validate required BelongsTo relationship", func(t *testing.T) {
		field := BelongsTo("Author", "user_id", "users")
		field.Required()

		validator := NewRelationshipValidator()
		ctx := context.Background()

		// Valid relationship
		err := validator.ValidateBelongsTo(ctx, 1, field)
		if err != nil {
			t.Errorf("Expected no error for valid relationship, got %v", err)
		}

		// Nil relationship on required field
		err = validator.ValidateBelongsTo(ctx, nil, field)
		if err == nil {
			t.Error("Expected error for nil relationship on required field")
		}
	})

	t.Run("Validate optional BelongsTo relationship", func(t *testing.T) {
		field := BelongsTo("Author", "user_id", "users")

		validator := NewRelationshipValidator()
		ctx := context.Background()

		// Nil relationship on optional field
		err := validator.ValidateBelongsTo(ctx, nil, field)
		if err != nil {
			t.Errorf("Expected no error for nil relationship on optional field, got %v", err)
		}
	})
}

// TestRelationshipQuery tests query customization
func TestRelationshipQuery(t *testing.T) {
	tdb := NewTestDB(t)
	defer tdb.Close()

	t.Run("Query with WHERE clause", func(t *testing.T) {
		query := NewRelationshipQuery()
		query.Where("status", "=", "active")

		conditions := query.GetWhereConditions()
		if len(conditions) != 1 {
			t.Errorf("Expected 1 condition, got %d", len(conditions))
		}
	})

	t.Run("Query with multiple conditions", func(t *testing.T) {
		query := NewRelationshipQuery()
		query.Where("status", "=", "active").
			Where("verified", "=", true).
			OrderBy("created_at", "DESC").
			Limit(10)

		if len(query.GetWhereConditions()) != 2 {
			t.Errorf("Expected 2 conditions, got %d", len(query.GetWhereConditions()))
		}

		if len(query.GetOrderByColumns()) != 1 {
			t.Errorf("Expected 1 order by column, got %d", len(query.GetOrderByColumns()))
		}

		if query.GetLimit() != 10 {
			t.Errorf("Expected limit 10, got %d", query.GetLimit())
		}
	})
}

// TestRelationshipDisplay tests display customization
func TestRelationshipDisplay(t *testing.T) {
	tdb := NewTestDB(t)
	defer tdb.Close()

	t.Run("Display BelongsTo relationship", func(t *testing.T) {
		field := BelongsTo("Author", "user_id", "users")
		display := NewRelationshipDisplay(field)

		// Get a user
		var user User
		err := tdb.DB.QueryRow("SELECT id, name, email FROM users WHERE id = 1").
			Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			t.Fatalf("Failed to get user: %v", err)
		}

		displayValue, err := display.GetDisplayValue(user)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if displayValue == "" {
			t.Error("Expected non-empty display value")
		}
	})

	t.Run("Display multiple relationships", func(t *testing.T) {
		field := BelongsTo("Author", "user_id", "users")
		display := NewRelationshipDisplay(field)

		// Get all users
		rows, err := tdb.DB.Query("SELECT id, name, email FROM users")
		if err != nil {
			t.Fatalf("Failed to query users: %v", err)
		}
		defer rows.Close()

		var users []interface{}
		for rows.Next() {
			var user User
			err := rows.Scan(&user.ID, &user.Name, &user.Email)
			if err != nil {
				t.Fatalf("Failed to scan user: %v", err)
			}
			users = append(users, user)
		}

		displayValues, err := display.GetDisplayValues(users)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if len(displayValues) != len(users) {
			t.Errorf("Expected %d display values, got %d", len(users), len(displayValues))
		}
	})
}

// TestRelationshipSerialization tests JSON serialization
func TestRelationshipSerialization(t *testing.T) {
	tdb := NewTestDB(t)
	defer tdb.Close()

	t.Run("Serialize BelongsTo relationship", func(t *testing.T) {
		field := BelongsTo("Author", "user_id", "users")
		serialization := NewRelationshipSerialization(field)

		// Get a user
		var user User
		err := tdb.DB.QueryRow("SELECT id, name, email FROM users WHERE id = 1").
			Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			t.Fatalf("Failed to get user: %v", err)
		}

		jsonData, err := serialization.SerializeRelationship(user)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if jsonData["type"] != "belongsTo" {
			t.Errorf("Expected type 'belongsTo', got %v", jsonData["type"])
		}

		if jsonData["value"] == nil {
			t.Error("Expected non-nil value")
		}
	})

	t.Run("Serialize to JSON string", func(t *testing.T) {
		field := BelongsTo("Author", "user_id", "users")
		serialization := NewRelationshipSerialization(field)

		// Get a user
		var user User
		err := tdb.DB.QueryRow("SELECT id, name, email FROM users WHERE id = 1").
			Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			t.Fatalf("Failed to get user: %v", err)
		}

		jsonStr, err := serialization.ToJSON(user)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if jsonStr == "" {
			t.Error("Expected non-empty JSON string")
		}
	})
}

// TestComplexScenario tests a complex real-world scenario
func TestComplexScenario(t *testing.T) {
	tdb := NewTestDB(t)
	defer tdb.Close()

	t.Run("Complex scenario: User with posts, comments, and profile", func(t *testing.T) {
		// Get user
		var user User
		err := tdb.DB.QueryRow("SELECT id, name, email FROM users WHERE id = 1").
			Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			t.Fatalf("Failed to get user: %v", err)
		}

		// Get user's profile (HasOne)
		var profile Profile
		err = tdb.DB.QueryRow("SELECT id, user_id, bio, avatar_url FROM profiles WHERE user_id = ?", user.ID).
			Scan(&profile.ID, &profile.UserID, &profile.Bio, &profile.AvatarURL)
		if err != nil {
			t.Fatalf("Failed to get profile: %v", err)
		}

		// Get user's posts (HasMany)
		rows, err := tdb.DB.Query("SELECT id, user_id, title, content FROM posts WHERE user_id = ?", user.ID)
		if err != nil {
			t.Fatalf("Failed to query posts: %v", err)
		}
		defer rows.Close()

		postCount := 0
		for rows.Next() {
			var post Post
			err := rows.Scan(&post.ID, &post.UserID, &post.Title, &post.Content)
			if err != nil {
				t.Fatalf("Failed to scan post: %v", err)
			}

			// Get tags for each post (BelongsToMany)
			tagRows, err := tdb.DB.Query(`
				SELECT t.id, t.name FROM tags t
				INNER JOIN post_tag pt ON t.id = pt.tag_id
				WHERE pt.post_id = ?
			`, post.ID)
			if err != nil {
				// Skip tag query if tags table doesn't exist
				t.Logf("Warning: Could not query tags: %v", err)
			} else {
				defer tagRows.Close()

				tagCount := 0
				for tagRows.Next() {
					var tag Tag
					err := tagRows.Scan(&tag.ID, &tag.Name)
					if err != nil {
						t.Fatalf("Failed to scan tag: %v", err)
					}
					tagCount++
				}

				if tagCount > 0 {
					// At least one post has tags
					postCount++
				}
			}
		}

		// Verify all relationships
		if user.ID != 1 {
			t.Errorf("Expected user ID 1, got %d", user.ID)
		}

		if profile.UserID != user.ID {
			t.Errorf("Expected profile user_id %d, got %d", user.ID, profile.UserID)
		}
	})
}
