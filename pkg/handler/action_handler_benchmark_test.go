package handler

import (
	"reflect"
	"strconv"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func BenchmarkLoadActionModelsByIDs_1000(b *testing.B) {
	db, err := gorm.Open(sqlite.Open("file:benchmark_action_loader?mode=memory&cache=shared"), &gorm.Config{})
	if err != nil {
		b.Fatalf("failed to open sqlite db: %v", err)
	}

	if err := db.AutoMigrate(&actionHandlerTestModel{}); err != nil {
		b.Fatalf("failed to migrate schema: %v", err)
	}

	seed := make([]actionHandlerTestModel, 1000)
	ids := make([]string, 1000)
	for i := range seed {
		id := i + 1
		seed[i] = actionHandlerTestModel{
			ID:   uint(id),
			Name: "user-" + strconv.Itoa(id),
		}
		ids[i] = strconv.Itoa(1000 - i)
	}

	if err := db.Create(&seed).Error; err != nil {
		b.Fatalf("failed to seed models: %v", err)
	}

	modelType := reflect.TypeOf(actionHandlerTestModel{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		models, err := loadActionModelsByIDs(db, modelType, ids)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
		if len(models) != len(ids) {
			b.Fatalf("expected %d models, got %d", len(ids), len(models))
		}
	}
}
