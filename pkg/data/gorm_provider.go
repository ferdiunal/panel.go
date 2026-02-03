package data

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"gorm.io/gorm"
)

type GormDataProvider struct {
	DB                *gorm.DB
	Model             interface{}
	SearchColumns     []string
	WithRelationships []string
}

func NewGormDataProvider(db *gorm.DB, model interface{}) *GormDataProvider {
	return &GormDataProvider{
		DB:    db,
		Model: model,
	}
}

func (p *GormDataProvider) SetSearchColumns(cols []string) {
	p.SearchColumns = cols
}

func (p *GormDataProvider) SetWith(rels []string) {
	p.WithRelationships = rels
}

func (p *GormDataProvider) Index(ctx context.Context, req QueryRequest) (*QueryResponse, error) {
	var total int64
	// We need a slice of the model type to hold results.
	// Since Model is interface{}, we might need reflection to create a slice of that type,
	// or we can just hope GORM handles Find(&[]Interface{}) correctly if we pass a pointer to a slice of models.
	// Actually, usually users pass a struct instance as Model.
	// Gorm's db.Model(model) works for setting the table.
	// But Find needs a destination.
	// Let's assume we return []map[string]interface{} for generic usage if we don't know the slice type,
	// OR we assume the user might want typed results.
	// But FieldHandler expects []interface{} in Items.

	// Simplest approach for Generic provider: Use map[string]interface{} for dynamic results
	// OR use reflection to make a slice of the Model's type.

	// Let's try map[string]interface{} for maximum flexibility in this generic provider,
	// unless we strictly want the structs.
	// If we use structs, we need to use reflect.New(reflect.SliceOf(reflect.TypeOf(p.Model))).Interface()

	// Let's start with just using db.Model(p.Model).Find(&results) where results is []map[string]interface{}
	// GORM supports finding into a map.

	// Log the incoming sorts to Provider
	fmt.Printf("DEBUG: Provider Index Req Sorts: %+v\n", req.Sorts)

	db := p.DB.WithContext(ctx).Debug().Model(p.Model)

	// Apply Eager Loading
	for _, rel := range p.WithRelationships {
		db = db.Preload(rel)
	}

	// Apply Filters (Basic equality)
	for k, v := range req.Filters {
		db = db.Where(fmt.Sprintf("%s = ?", k), v)
	}

	// Apply Search
	if req.Search != "" && len(p.SearchColumns) > 0 {
		searchQuery := p.DB.WithContext(ctx).Session(&gorm.Session{NewDB: true})
		for _, col := range p.SearchColumns {
			searchQuery = searchQuery.Or(fmt.Sprintf("%s LIKE ?", col), "%"+req.Search+"%")
		}
		db = db.Where(searchQuery)
	}

	// Count Total
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	// Sorting
	if len(req.Sorts) > 0 {
		for _, sort := range req.Sorts {
			if sort.Column != "" {
				direction := "ASC"
				if strings.ToUpper(sort.Direction) == "DESC" {
					direction = "DESC"
				}
				db = db.Order(fmt.Sprintf("%s %s", sort.Column, direction))
			}
		}
	}

	// Pagination
	offset := (req.Page - 1) * req.PerPage
	db = db.Offset(offset).Limit(req.PerPage)

	// Execute Query
	// Use reflection to create a slice of the model type
	modelType := reflect.TypeOf(p.Model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	sliceType := reflect.SliceOf(modelType)
	resultsPtr := reflect.New(sliceType)

	if err := db.Find(resultsPtr.Interface()).Error; err != nil {
		return nil, err
	}

	// Convert to []interface{}
	// Since FieldHandler (and Tests) mostly expect maps for dynamic access, lets convert strict structs to maps
	// This also ensures we respect JSON tags.
	resultsVal := resultsPtr.Elem()
	items := make([]interface{}, resultsVal.Len())
	for i := 0; i < resultsVal.Len(); i++ {
		items[i] = resultsVal.Index(i).Addr().Interface()
	}

	return &QueryResponse{
		Items:   items,
		Total:   total,
		Page:    req.Page,
		PerPage: req.PerPage,
	}, nil
}

func (p *GormDataProvider) Show(ctx context.Context, id string) (interface{}, error) {
	// Create a new instance of the model to hold the result
	// We use p.Model's type
	// But simpler: just use map[string]interface{} for dynamic nature or try to use the model type via reflection if needed.
	// For GORM, if we pass p.Model (which is a pointer to a struct), it works but we might overwrite the original p.Model if we are not careful or if we reuse it?
	// Actually p.Model is just a template.
	// Let's use map[string]interface{} to return standard format for the handler.

	// Create a new instance of the model to hold the result
	modelType := reflect.TypeOf(p.Model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	result := reflect.New(modelType).Interface()

	db := p.DB.WithContext(ctx).Model(p.Model)
	for _, rel := range p.WithRelationships {
		db = db.Preload(rel)
	}

	if err := db.Where("id = ?", id).First(result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (p *GormDataProvider) Create(ctx context.Context, data map[string]interface{}) (interface{}, error) {
	// 1. Create a new instance of the model
	modelType := reflect.TypeOf(p.Model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	newItem := reflect.New(modelType).Interface()

	// 2. Convert map to struct (using json roundtrip for simplicity and tag support)
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, newItem); err != nil {
		return nil, err
	}

	// 3. Create in DB
	if err := p.DB.WithContext(ctx).Model(p.Model).Create(newItem).Error; err != nil {
		return nil, err
	}

	return newItem, nil
}

func (p *GormDataProvider) Update(ctx context.Context, id string, data map[string]interface{}) (interface{}, error) {
	data["updated_at"] = time.Now()
	if err := p.DB.WithContext(ctx).Model(p.Model).Where("id = ?", id).Updates(data).Error; err != nil {
		return nil, err
	}
	// Return updated struct or just the data?
	// Let's return the fresh data
	return p.Show(ctx, id)
}

func (p *GormDataProvider) Delete(ctx context.Context, id string) error {
	return p.DB.WithContext(ctx).Model(p.Model).Where("id = ?", id).Delete(nil).Error
}
