package handler

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type actionHandlerTestModel struct {
	ID   uint `gorm:"primaryKey"`
	Name string
}

func newActionHandlerTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", strings.ReplaceAll(t.Name(), "/", "_"))
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite db: %v", err)
	}

	if err := db.AutoMigrate(&actionHandlerTestModel{}); err != nil {
		t.Fatalf("failed to migrate schema: %v", err)
	}

	seed := []actionHandlerTestModel{
		{ID: 1, Name: "one"},
		{ID: 2, Name: "two"},
		{ID: 3, Name: "three"},
	}
	if err := db.Create(&seed).Error; err != nil {
		t.Fatalf("failed to seed models: %v", err)
	}

	return db
}

func TestLoadActionModelsByIDs_PreservesRequestOrder(t *testing.T) {
	db := newActionHandlerTestDB(t)
	ids := []string{"3", "1", "2"}

	models, err := loadActionModelsByIDs(db, reflect.TypeOf(actionHandlerTestModel{}), ids)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(models) != len(ids) {
		t.Fatalf("expected %d models, got %d", len(ids), len(models))
	}

	for i, model := range models {
		id, ok := extractModelIDString(model)
		if !ok {
			t.Fatalf("failed to extract id at index %d", i)
		}
		if id != ids[i] {
			t.Fatalf("order mismatch at index %d: expected id=%s got id=%s", i, ids[i], id)
		}
	}
}

func TestLoadActionModelsByIDs_MissingIDReturnsNotFound(t *testing.T) {
	db := newActionHandlerTestDB(t)

	_, err := loadActionModelsByIDs(db, reflect.TypeOf(actionHandlerTestModel{}), []string{"1", "999"})
	if err == nil {
		t.Fatalf("expected missing id error")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected gorm.ErrRecordNotFound, got %v", err)
	}
}

func TestLoadActionModelsByIDs_SupportsDuplicateRequestIDs(t *testing.T) {
	db := newActionHandlerTestDB(t)
	ids := []string{"2", "2", "1"}

	models, err := loadActionModelsByIDs(db, reflect.TypeOf(actionHandlerTestModel{}), ids)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(models) != len(ids) {
		t.Fatalf("expected %d models, got %d", len(ids), len(models))
	}

	for i, model := range models {
		id, ok := extractModelIDString(model)
		if !ok {
			t.Fatalf("failed to extract id at index %d", i)
		}
		if id != ids[i] {
			t.Fatalf("duplicate order mismatch at index %d: expected id=%s got id=%s", i, ids[i], id)
		}
	}
}
