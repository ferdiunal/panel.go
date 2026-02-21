package data

import (
	"fmt"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type updateNoTimestampArea struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	SortOrder int
}

func newUpdateTimestampTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect sqlite in-memory db: %v", err)
	}

	return db
}

func TestGormDataProvider_Update_SkipsMissingUpdatedAtColumn(t *testing.T) {
	db := newUpdateTimestampTestDB(t)
	if err := db.AutoMigrate(&updateNoTimestampArea{}); err != nil {
		t.Fatalf("failed to migrate table: %v", err)
	}

	area := updateNoTimestampArea{
		Name:      "Area-1",
		SortOrder: 1,
	}
	if err := db.Create(&area).Error; err != nil {
		t.Fatalf("failed to seed area: %v", err)
	}

	provider := NewGormDataProvider(db, &updateNoTimestampArea{})
	if _, err := provider.Update(nil, fmt.Sprint(area.ID), map[string]interface{}{
		"name":       "Area-1 Updated",
		"sort_order": 3,
	}); err != nil {
		t.Fatalf("update should not fail when updated_at column is missing: %v", err)
	}

	var reloaded updateNoTimestampArea
	if err := db.First(&reloaded, area.ID).Error; err != nil {
		t.Fatalf("failed to reload updated row: %v", err)
	}

	if reloaded.Name != "Area-1 Updated" {
		t.Fatalf("expected name to be updated, got %q", reloaded.Name)
	}

	if reloaded.SortOrder != 3 {
		t.Fatalf("expected sort_order to be 3, got %d", reloaded.SortOrder)
	}
}

