package fields

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// TestDB provides a test database connection
type TestDB struct {
	DB *sql.DB
	t  *testing.T
}

// NewTestDB creates a new test database
func NewTestDB(t *testing.T) *TestDB {
	// Use in-memory SQLite database for tests
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("Failed to open test database: %v", err)
	}

	// Create tables
	if err := createTestTables(db); err != nil {
		t.Fatalf("Failed to create test tables: %v", err)
	}

	// Seed test data
	if err := seedTestData(db); err != nil {
		t.Fatalf("Failed to seed test data: %v", err)
	}

	return &TestDB{
		DB: db,
		t:  t,
	}
}

// Close closes the test database
func (tdb *TestDB) Close() error {
	return tdb.DB.Close()
}

// createTestTables creates all necessary test tables
func createTestTables(db *sql.DB) error {
	tables := []string{
		// Users table
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		// Posts table
		`CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			content TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,

		// Comments table
		`CREATE TABLE IF NOT EXISTS comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (post_id) REFERENCES posts(id),
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,

		// Tags table
		`CREATE TABLE IF NOT EXISTS tags (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,

		// Post-Tag pivot table (BelongsToMany)
		`CREATE TABLE IF NOT EXISTS post_tag (
			post_id INTEGER NOT NULL,
			tag_id INTEGER NOT NULL,
			PRIMARY KEY (post_id, tag_id),
			FOREIGN KEY (post_id) REFERENCES posts(id),
			FOREIGN KEY (tag_id) REFERENCES tags(id)
		)`,

		// Taggable polymorphic table
		`CREATE TABLE IF NOT EXISTS taggables (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			taggable_id INTEGER NOT NULL,
			taggable_type TEXT NOT NULL,
			tag_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (tag_id) REFERENCES tags(id)
		)`,

		// Profiles table (HasOne)
		`CREATE TABLE IF NOT EXISTS profiles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL UNIQUE,
			bio TEXT,
			avatar_url TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
	}

	for _, table := range tables {
		if _, err := db.Exec(table); err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	return nil
}

// seedTestData seeds test data into the database
func seedTestData(db *sql.DB) error {
	// Insert users
	users := []struct {
		name  string
		email string
	}{
		{"John Doe", "john@example.com"},
		{"Jane Smith", "jane@example.com"},
		{"Bob Johnson", "bob@example.com"},
	}

	for _, user := range users {
		_, err := db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", user.name, user.email)
		if err != nil {
			return fmt.Errorf("failed to insert user: %w", err)
		}
	}

	// Insert posts
	posts := []struct {
		userID  int
		title   string
		content string
	}{
		{1, "First Post", "This is the first post"},
		{1, "Second Post", "This is the second post"},
		{2, "Jane's Post", "Jane's first post"},
	}

	for _, post := range posts {
		_, err := db.Exec("INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)", post.userID, post.title, post.content)
		if err != nil {
			return fmt.Errorf("failed to insert post: %w", err)
		}
	}

	// Insert comments
	comments := []struct {
		postID  int
		userID  int
		content string
	}{
		{1, 2, "Great post!"},
		{1, 3, "Thanks for sharing"},
		{2, 2, "Interesting"},
	}

	for _, comment := range comments {
		_, err := db.Exec("INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)", comment.postID, comment.userID, comment.content)
		if err != nil {
			return fmt.Errorf("failed to insert comment: %w", err)
		}
	}

	// Insert tags
	tags := []string{"golang", "database", "testing", "api"}
	for _, tag := range tags {
		_, err := db.Exec("INSERT INTO tags (name) VALUES (?)", tag)
		if err != nil {
			return fmt.Errorf("failed to insert tag: %w", err)
		}
	}

	// Insert post-tag relationships
	postTags := []struct {
		postID int
		tagID  int
	}{
		{1, 1}, // First post - golang
		{1, 2}, // First post - database
		{2, 1}, // Second post - golang
		{3, 3}, // Jane's post - testing
	}

	for _, pt := range postTags {
		_, err := db.Exec("INSERT INTO post_tag (post_id, tag_id) VALUES (?, ?)", pt.postID, pt.tagID)
		if err != nil {
			return fmt.Errorf("failed to insert post-tag: %w", err)
		}
	}

	// Insert profiles
	profiles := []struct {
		userID    int
		bio       string
		avatarURL string
	}{
		{1, "Software developer", "https://example.com/avatar1.jpg"},
		{2, "Designer", "https://example.com/avatar2.jpg"},
	}

	for _, profile := range profiles {
		_, err := db.Exec("INSERT INTO profiles (user_id, bio, avatar_url) VALUES (?, ?, ?)", profile.userID, profile.bio, profile.avatarURL)
		if err != nil {
			return fmt.Errorf("failed to insert profile: %w", err)
		}
	}

	// Insert taggable relationships (polymorphic)
	taggables := []struct {
		taggableID   int
		taggableType string
		tagID        int
	}{
		{1, "post", 1},    // Post 1 tagged with golang
		{1, "post", 2},    // Post 1 tagged with database
		{1, "comment", 3}, // Comment 1 tagged with testing
	}

	for _, taggable := range taggables {
		_, err := db.Exec("INSERT INTO taggables (taggable_id, taggable_type, tag_id) VALUES (?, ?, ?)", taggable.taggableID, taggable.taggableType, taggable.tagID)
		if err != nil {
			return fmt.Errorf("failed to insert taggable: %w", err)
		}
	}

	return nil
}

// Dummy resource structs for testing

// User represents a user resource
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Post represents a post resource
type Post struct {
	ID      int    `json:"id"`
	UserID  int    `json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  *User  `json:"author,omitempty"`
}

// Comment represents a comment resource
type Comment struct {
	ID      int    `json:"id"`
	PostID  int    `json:"post_id"`
	UserID  int    `json:"user_id"`
	Content string `json:"content"`
	Post    *Post  `json:"post,omitempty"`
	Author  *User  `json:"author,omitempty"`
}

// Tag represents a tag resource
type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Profile represents a user profile resource
type Profile struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Bio       string `json:"bio"`
	AvatarURL string `json:"avatar_url"`
	User      *User  `json:"user,omitempty"`
}

// PostTag represents a post-tag relationship
type PostTag struct {
	PostID int `json:"post_id"`
	TagID  int `json:"tag_id"`
}

// Taggable represents a polymorphic taggable relationship
type Taggable struct {
	ID           int    `json:"id"`
	TaggableID   int    `json:"taggable_id"`
	TaggableType string `json:"taggable_type"`
	TagID        int    `json:"tag_id"`
}

// GetTestDataDir returns the test data directory
func GetTestDataDir() string {
	return os.TempDir()
}
