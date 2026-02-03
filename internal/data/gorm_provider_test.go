package data

import (
	"context"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type TestUser struct {
	ID    uint       `gorm:"primaryKey" json:"id"`
	Name  string     `json:"name"`
	Email string     `json:"email"`
	Posts []TestPost `json:"posts" gorm:"foreignKey:UserID"`
}

type TestPost struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	UserID uint   `json:"user_id"`
	Title  string `json:"title"`
}

func TestGormDataProvider_Index(t *testing.T) {
	// Setup In-Memory DB
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to db: %v", err)
	}

	// Migrate
	db.AutoMigrate(&TestUser{})

	// Seed Data
	users := []TestUser{
		{Name: "Alice", Email: "alice@example.com"},
		{Name: "Bob", Email: "bob@example.com"},
		{Name: "Charlie", Email: "charlie@example.com"},
	}
	db.Create(&users)

	// Test Provider
	provider := NewGormDataProvider(db, &TestUser{})

	ctx := context.Background()
	req := QueryRequest{
		Page:    1,
		PerPage: 2,
		Sorts: []Sort{
			{Column: "name", Direction: "asc"},
		},
	}

	resp, err := provider.Index(ctx, req)
	if err != nil {
		t.Fatalf("Index failed: %v", err)
	}

	if resp.Total != 3 {
		t.Errorf("Expected total 3, got %d", resp.Total)
	}
	if len(resp.Items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(resp.Items))
	}

	// Since we are fetching into interface{}, check type
	item1 := resp.Items[0].(*TestUser)
	if item1.Name != "Alice" {
		t.Errorf("Expected Alice, got %v", item1.Name)
	}

	// Test Sort DESC
	req.Sorts[0].Direction = "desc"
	respDesc, _ := provider.Index(ctx, req)
	itemDesc := respDesc.Items[0].(*TestUser)
	if itemDesc.Name != "Charlie" { // C comes after A and B, so desc should be Charlie
		t.Errorf("Expected Charlie, got %v", itemDesc.Name)
	}
}

func TestGormDataProvider_Search(t *testing.T) {
	// Setup In-Memory DB
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to db: %v", err)
	}

	// Clean slate
	db.Migrator().DropTable(&TestUser{})

	// Migrate
	db.AutoMigrate(&TestUser{})

	// Seed Data
	users := []TestUser{
		{Name: "Alice", Email: "alice@example.com"},     // No match
		{Name: "Bob", Email: "bobby@test.com"},          // Match by Name "Bob"
		{Name: "Charlie", Email: "charlie@example.com"}, // Match by Email or Name
		{Name: "Dave", Email: "dave@bob.com"},           // Match by Email "bob"
	}
	db.Create(&users)

	// Test Provider
	provider := NewGormDataProvider(db, &TestUser{})

	// 1. Configure Search Columns
	provider.SetSearchColumns([]string{"name", "email"})

	// 2. Search for "Bob" (Should match Bob by name and Dave by email)
	ctx := context.Background()
	req := QueryRequest{
		Page:    1,
		PerPage: 10,
		Search:  "Bob",
	}

	resp, err := provider.Index(ctx, req)
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if resp.Total != 2 {
		t.Errorf("Expected 2 matches for 'Bob', got %d", resp.Total)
	}

	// Verify items
	foundNames := make(map[string]bool)
	for _, item := range resp.Items {
		u := item.(*TestUser)
		foundNames[u.Name] = true
	}

	if !foundNames["Bob"] {
		t.Errorf("Expected to find Bob")
	}
	if !foundNames["Dave"] {
		t.Errorf("Expected to find Dave (via email match)")
	}

	// 3. Search for "Ali" (Should match Alice)
	req.Search = "Ali"
	resp, _ = provider.Index(ctx, req)
	if resp.Total != 1 {
		t.Errorf("Expected 1 match for 'Ali', got %d", resp.Total)
	}
}

func TestGormDataProvider_With(t *testing.T) {
	// Setup In-Memory DB
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to db: %v", err)
	}

	// Clean slate
	db.Migrator().DropTable(&TestUser{}, &TestPost{})

	// Migrate
	db.AutoMigrate(&TestUser{}, &TestPost{})

	// Seed Data
	db.Create(&TestUser{
		Name:  "User With Posts",
		Email: "rel@example.com",
		Posts: []TestPost{
			{Title: "Post 1"},
			{Title: "Post 2"},
		},
	})

	// Test Provider
	provider := NewGormDataProvider(db, &TestUser{})

	// 1. Configure Eager Loading
	provider.SetWith([]string{"Posts"})

	// 2. Fetch Index
	ctx := context.Background()
	req := QueryRequest{
		Page:    1,
		PerPage: 10,
	}

	resp, err := provider.Index(ctx, req)
	if err != nil {
		t.Fatalf("Index failed: %v", err)
	}

	if len(resp.Items) != 1 {
		t.Fatalf("Expected 1 user, got %d", len(resp.Items))
	}

	// Verify Posts are loaded
	// Since we return []map[string]interface{}, GORM's map output might not include nested structs automatically unless we use JSON serialization trick or GORM specific map mode.
	// Wait, our GormDataProvider uses `db.Find(&results)` where results is `[]map[string]interface{}`.
	// GORM DOES NOT automatically populate nested structs into map[string]interface{} when using Find(&map). It usually works better with Struct slice.
	// However, our Provider is Generic.
	// If `Model` is `&TestUser{}`, `db.Model(p.Model)` knows the schema.
	// But `db.Find(&results)` into map might lose the relationship if GORM doesn't flatten it.
	// Actually, Gorm Preload works by populating fields on the struct. If we scan into map, Preload might be ignored or not mapped.
	// Let's verify this behavior. If it fails, we might need to change Provider to use Reflect to creating a slice of the Model type.
	// Given the previous task status was "Success", let's assume map scanning works OR we need to accept we might need reflection.
	// Let's check the Item.

	// Actually, to support relations properly in a generic way, we probably SHOULD use reflection to create a slice of the real type.
	// But for now let's write the test and see if it fails.

	item := resp.Items[0].(*TestUser)
	if len(item.Posts) != 2 {
		t.Fatalf("Expected 2 posts, got %d", len(item.Posts))
	}
}
