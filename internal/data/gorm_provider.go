package data

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (p *GormDataProvider) columnAliases() map[string]string {
	aliases := map[string]string{}

	stmt := &gorm.Statement{DB: p.DB}
	if err := stmt.Parse(p.Model); err != nil || stmt.Schema == nil {
		return aliases
	}

	for _, field := range stmt.Schema.Fields {
		if field.DBName == "" {
			continue
		}

		dbName := field.DBName
		aliases[dbName] = dbName
		aliases[field.Name] = dbName
		aliases[strings.ToLower(field.Name)] = dbName

		if tag := field.Tag.Get("json"); tag != "" {
			jsonName := strings.Split(tag, ",")[0]
			if jsonName != "" && jsonName != "-" {
				aliases[jsonName] = dbName
			}
		}
	}

	return aliases
}

func normalizeColumn(column string, aliases map[string]string) (string, bool) {
	col := strings.TrimSpace(column)
	if col == "" {
		return "", false
	}

	dbCol, ok := aliases[col]
	return dbCol, ok
}

func (p *GormDataProvider) Index(ctx context.Context, req QueryRequest) (*QueryResponse, error) {
	var total int64
	aliases := p.columnAliases()

	db := p.DB.WithContext(ctx).Model(p.Model)

	// Apply Eager Loading
	for _, rel := range p.WithRelationships {
		db = db.Preload(rel)
	}

	// Apply Filters (Basic equality)
	for k, v := range req.Filters {
		col, ok := normalizeColumn(k, aliases)
		if !ok {
			continue
		}
		db = db.Where(clause.Eq{
			Column: clause.Column{Name: col},
			Value:  v,
		})
	}

	// Apply Search
	if req.Search != "" && len(p.SearchColumns) > 0 {
		searchQuery := p.DB.WithContext(ctx).Session(&gorm.Session{NewDB: true}).Model(p.Model)
		hasValidSearchColumn := false
		for _, col := range p.SearchColumns {
			dbCol, ok := normalizeColumn(col, aliases)
			if !ok {
				continue
			}

			condition := clause.Like{
				Column: clause.Column{Name: dbCol},
				Value:  "%" + req.Search + "%",
			}

			if !hasValidSearchColumn {
				searchQuery = searchQuery.Where(condition)
				hasValidSearchColumn = true
				continue
			}
			searchQuery = searchQuery.Or(condition)
		}

		if hasValidSearchColumn {
			db = db.Where(searchQuery)
		}
	}

	// Count Total
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	// Sorting
	if len(req.Sorts) > 0 {
		for _, sort := range req.Sorts {
			col, ok := normalizeColumn(sort.Column, aliases)
			if !ok {
				continue
			}

			db = db.Order(clause.OrderByColumn{
				Column: clause.Column{Name: col},
				Desc:   strings.EqualFold(sort.Direction, "desc"),
			})
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
