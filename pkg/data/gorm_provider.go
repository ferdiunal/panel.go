package data

import (
	stdcontext "context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/query"
	"github.com/iancoleman/strcase"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type GormDataProvider struct {
	DB                *gorm.DB
	Model             interface{}
	SearchColumns     []string
	WithRelationships []string
	columnValidator   *ColumnValidator
}

func NewGormDataProvider(db *gorm.DB, model interface{}) *GormDataProvider {
	// Initialize column validator for SQL injection protection
	validator, err := NewColumnValidator(db, model)
	if err != nil {
		// Log error but don't fail - fall back to basic validation
		fmt.Printf("[SECURITY WARNING] Failed to create column validator: %v\n", err)
	}

	return &GormDataProvider{
		DB:              db,
		Model:           model,
		columnValidator: validator,
	}
}

// getContext safely extracts the standard context from our custom Context
// Returns context.Background() if ctx is nil or has no underlying context
func (p *GormDataProvider) getContext(ctx *context.Context) stdcontext.Context {
	if ctx == nil {
		return stdcontext.Background()
	}
	stdCtx := ctx.Context()
	if stdCtx == nil {
		return stdcontext.Background()
	}
	return stdCtx
}

func (p *GormDataProvider) SetSearchColumns(cols []string) {
	p.SearchColumns = cols
}

func (p *GormDataProvider) SetWith(rels []string) {
	p.WithRelationships = rels
}

// applyFilters applies advanced filter conditions to the GORM query
// Supports operators: eq, neq, gt, gte, lt, lte, like, nlike, in, nin, null, nnull, between
// SECURITY: Validates column names to prevent SQL injection
func (p *GormDataProvider) applyFilters(db *gorm.DB, filters []query.Filter) *gorm.DB {
	for _, f := range filters {
		if f.Field == "" {
			continue
		}

		// SECURITY: Validate column name to prevent SQL injection
		safeColumn := f.Field
		if p.columnValidator != nil {
			validatedCol, err := p.columnValidator.ValidateColumn(f.Field)
			if err != nil {
				// Skip invalid columns - don't expose error to user
				fmt.Printf("[SECURITY] Rejected invalid column in filter: %s\n", f.Field)
				continue
			}
			safeColumn = validatedCol
		} else {
			// Fallback: sanitize column name if validator not available
			safeColumn = SanitizeColumnName(f.Field)
		}

		switch f.Operator {
		case query.OpEqual:
			db = db.Where(fmt.Sprintf("%s = ?", safeColumn), f.Value)

		case query.OpNotEqual:
			db = db.Where(fmt.Sprintf("%s != ?", safeColumn), f.Value)

		case query.OpGreaterThan:
			db = db.Where(fmt.Sprintf("%s > ?", safeColumn), f.Value)

		case query.OpGreaterEq:
			db = db.Where(fmt.Sprintf("%s >= ?", safeColumn), f.Value)

		case query.OpLessThan:
			db = db.Where(fmt.Sprintf("%s < ?", safeColumn), f.Value)

		case query.OpLessEq:
			db = db.Where(fmt.Sprintf("%s <= ?", safeColumn), f.Value)

		case query.OpLike:
			if strVal, ok := f.Value.(string); ok {
				db = db.Where(fmt.Sprintf("%s LIKE ?", safeColumn), "%"+strVal+"%")
			}

		case query.OpNotLike:
			if strVal, ok := f.Value.(string); ok {
				db = db.Where(fmt.Sprintf("%s NOT LIKE ?", safeColumn), "%"+strVal+"%")
			}

		case query.OpIn:
			if vals, ok := f.Value.([]string); ok && len(vals) > 0 {
				db = db.Where(fmt.Sprintf("%s IN ?", safeColumn), vals)
			}

		case query.OpNotIn:
			if vals, ok := f.Value.([]string); ok && len(vals) > 0 {
				db = db.Where(fmt.Sprintf("%s NOT IN ?", safeColumn), vals)
			}

		case query.OpIsNull:
			if boolVal, ok := f.Value.(bool); ok && boolVal {
				db = db.Where(fmt.Sprintf("%s IS NULL", safeColumn))
			}

		case query.OpIsNotNull:
			if boolVal, ok := f.Value.(bool); ok && boolVal {
				db = db.Where(fmt.Sprintf("%s IS NOT NULL", safeColumn))
			}

		case query.OpBetween:
			if vals, ok := f.Value.([]string); ok && len(vals) == 2 {
				db = db.Where(fmt.Sprintf("%s BETWEEN ? AND ?", safeColumn), vals[0], vals[1])
			}

		default:
			// Default to equality
			db = db.Where(fmt.Sprintf("%s = ?", safeColumn), f.Value)
		}
	}
	return db
}

func (p *GormDataProvider) Index(ctx *context.Context, req QueryRequest) (*QueryResponse, error) {
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

	stdCtx := p.getContext(ctx)
	db := p.DB.WithContext(stdCtx).Model(p.Model)

	// Apply Eager Loading
	for _, rel := range p.WithRelationships {
		db = db.Preload(rel)
	}

	// Apply Advanced Filters
	if len(req.Filters) > 0 {
		db = p.applyFilters(db, req.Filters)
	}

	// Apply Search with column validation
	fmt.Printf("[GORM] Search: %q, SearchColumns: %v\n", req.Search, p.SearchColumns)
	if req.Search != "" && len(p.SearchColumns) > 0 {
		searchQuery := p.DB.WithContext(stdCtx).Session(&gorm.Session{NewDB: true})
		for _, col := range p.SearchColumns {
			// SECURITY: Validate search column names
			safeColumn := col
			if p.columnValidator != nil {
				validatedCol, err := p.columnValidator.ValidateColumn(col)
				if err != nil {
					// Skip invalid columns - don't expose error to user
					fmt.Printf("[SECURITY] Rejected invalid search column: %s\n", col)
					continue
				}
				safeColumn = validatedCol
			} else {
				// Fallback: sanitize column name if validator not available
				safeColumn = SanitizeColumnName(col)
			}
			searchQuery = searchQuery.Or(fmt.Sprintf("%s LIKE ?", safeColumn), "%"+req.Search+"%")
		}
		db = db.Where(searchQuery)
		fmt.Printf("[GORM] Search applied for columns: %v\n", p.SearchColumns)
	} else {
		fmt.Printf("[GORM] Search NOT applied - Search empty: %v, SearchColumns empty: %v\n", req.Search == "", len(p.SearchColumns) == 0)
	}

	// Count Total
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	// Sorting with column validation
	if len(req.Sorts) > 0 {
		for _, sort := range req.Sorts {
			if sort.Column != "" {
				// SECURITY: Validate sort column names
				safeColumn := sort.Column
				if p.columnValidator != nil {
					validatedCol, err := p.columnValidator.ValidateColumn(sort.Column)
					if err != nil {
						// Skip invalid columns - don't expose error to user
						fmt.Printf("[SECURITY] Rejected invalid sort column: %s\n", sort.Column)
						continue
					}
					safeColumn = validatedCol
				} else {
					// Fallback: sanitize column name if validator not available
					safeColumn = SanitizeColumnName(sort.Column)
				}

				direction := "ASC"
				if strings.ToUpper(sort.Direction) == "DESC" {
					direction = "DESC"
				}
				db = db.Order(fmt.Sprintf("%s %s", safeColumn, direction))
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

func (p *GormDataProvider) Show(ctx *context.Context, id string) (interface{}, error) {
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

	stdCtx := p.getContext(ctx)
	db := p.DB.WithContext(stdCtx).Model(p.Model)
	for _, rel := range p.WithRelationships {
		db = db.Preload(rel)
	}

	if err := db.Where("id = ?", id).First(result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func (p *GormDataProvider) Create(ctx *context.Context, data map[string]interface{}) (interface{}, error) {
	stmt := &gorm.Statement{DB: p.DB}
	if err := stmt.Parse(p.Model); err != nil {
		return nil, err
	}
	modelSchema := stmt.Schema

	stdCtx := p.getContext(ctx)
	modelType := reflect.TypeOf(p.Model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	newItem := reflect.New(modelType).Interface()

	validData := make(map[string]interface{})
	for k, v := range data {
		field := modelSchema.LookUpField(k)
		if field != nil && field.DBName != "" {
			validData[k] = v

			// Set field value on newItem to ensure it's populated for Create
			modelVal := reflect.ValueOf(newItem)
			if modelVal.Kind() == reflect.Ptr {
				modelVal = modelVal.Elem()
			}
			if err := field.Set(stdCtx, modelVal, v); err != nil {
				fmt.Printf("[GORM] Error setting field %s: %v\n", field.Name, err)
			}
		}
	}

	// Set timestamps
	now := time.Now()
	if createdAtField := modelSchema.LookUpField("CreatedAt"); createdAtField != nil {
		modelVal := reflect.ValueOf(newItem)
		if modelVal.Kind() == reflect.Ptr {
			modelVal = modelVal.Elem()
		}
		createdAtField.Set(stdCtx, modelVal, now)
	}
	if updatedAtField := modelSchema.LookUpField("UpdatedAt"); updatedAtField != nil {
		modelVal := reflect.ValueOf(newItem)
		if modelVal.Kind() == reflect.Ptr {
			modelVal = modelVal.Elem()
		}
		updatedAtField.Set(stdCtx, modelVal, now)
	}

	// Use Create with struct to ensure ID backfilling and hooks execution
	if err := p.DB.WithContext(stdCtx).Create(newItem).Error; err != nil {
		return nil, err
	}

	// Handle Associations
	for k, v := range data {
		field := modelSchema.LookUpField(k)
		if field == nil {
			field = modelSchema.LookUpField(strcase.ToCamel(k))
		}
		if field != nil {
			if field.DBName == "" {
				if rel, ok := modelSchema.Relationships.Relations[field.Name]; ok {
					relName := field.Name
					switch rel.Type {
					case schema.HasOne:
						relType := rel.FieldSchema.ModelType
						if relType.Kind() == reflect.Ptr {
							relType = relType.Elem()
						}
						relInstance := reflect.New(relType).Interface()

						if v != nil {
							if err := p.DB.WithContext(stdCtx).First(relInstance, v).Error; err == nil {
								p.DB.WithContext(stdCtx).Model(newItem).Association(relName).Replace(relInstance)
							}
						}
					case schema.Many2Many:
						// Handle BelongsToMany (Many2Many)
						// v is likely []interface{} or []string of IDs
						var ids []interface{}
						val := reflect.ValueOf(v)
						if val.Kind() == reflect.Slice {
							for i := 0; i < val.Len(); i++ {
								ids = append(ids, val.Index(i).Interface())
							}
						} else {
							ids = append(ids, v)
						}

						if len(ids) > 0 {
							relType := rel.FieldSchema.ModelType
							if relType.Kind() == reflect.Slice {
								relType = relType.Elem()
							}
							if relType.Kind() == reflect.Ptr {
								relType = relType.Elem()
							}

							// Create a slice of related structs with just IDs
							// GORM Association().Replace() works best with struct instances or slice of structs
							// It can also take slice of primary keys but let's try to be safe

							// Ideally we should find them to ensure they exist, but for performance just binding IDs might work
							// if we use Omit("Example.*") to avoid updating them, but Association Replace handles linking.

							// Let's query them to be safe and GORM-compliant
							sliceType := reflect.SliceOf(reflect.PtrTo(relType))
							relatedItems := reflect.New(sliceType).Interface()

							// Assuming ID is integer or string, GORM Find with slice of IDs works
							if err := p.DB.WithContext(stdCtx).Where("id IN ?", ids).Find(relatedItems).Error; err == nil {
								if err := p.DB.WithContext(stdCtx).Model(newItem).Association(relName).Replace(relatedItems); err != nil {
									// Log error?
								}
							}
						}
					}
				}
			}
		}
	}

	// Return fresh item using ID
	if modelSchema.PrioritizedPrimaryField != nil {
		idVal := reflect.ValueOf(newItem).Elem().FieldByName(modelSchema.PrioritizedPrimaryField.Name).Interface()
		id := fmt.Sprint(idVal)
		if id != "" && id != "0" {
			return p.Show(ctx, id)
		}
	}

	return newItem, nil
}

func (p *GormDataProvider) Update(ctx *context.Context, id string, data map[string]interface{}) (interface{}, error) {
	stmt := &gorm.Statement{DB: p.DB}
	if err := stmt.Parse(p.Model); err != nil {
		return nil, err
	}
	modelSchema := stmt.Schema

	stdCtx := p.getContext(ctx)
	modelType := reflect.TypeOf(p.Model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	item := reflect.New(modelType).Interface()

	if err := p.DB.WithContext(stdCtx).First(item, "id = ?", id).Error; err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})
	for k, v := range data {
		field := modelSchema.LookUpField(k)
		if field == nil {
			field = modelSchema.LookUpField(strcase.ToCamel(k))
		}
		if field != nil {
			if field.DBName != "" {
				updates[k] = v
			} else if rel, ok := modelSchema.Relationships.Relations[field.Name]; ok {
				relName := field.Name
				switch rel.Type {
				case schema.HasOne, schema.BelongsTo:
					relType := rel.FieldSchema.ModelType
					if relType.Kind() == reflect.Ptr {
						relType = relType.Elem()
					}
					relInstance := reflect.New(relType).Interface()

					if v != nil {
						if err := p.DB.WithContext(stdCtx).First(relInstance, v).Error; err == nil {
							p.DB.WithContext(stdCtx).Model(item).Association(relName).Replace(relInstance)
						}
					} else {
						p.DB.WithContext(stdCtx).Model(item).Association(relName).Clear()
					}
				case schema.Many2Many:
					// Handle BelongsToMany (Many2Many) update
					// v is likely []interface{} or []string of IDs
					var ids []interface{}
					if v != nil {
						val := reflect.ValueOf(v)
						if val.Kind() == reflect.Slice {
							for i := 0; i < val.Len(); i++ {
								ids = append(ids, val.Index(i).Interface())
							}
						} else {
							ids = append(ids, v)
						}
					}

					if len(ids) > 0 {
						relType := rel.FieldSchema.ModelType
						if relType.Kind() == reflect.Slice {
							relType = relType.Elem()
						}
						if relType.Kind() == reflect.Ptr {
							relType = relType.Elem()
						}

						sliceType := reflect.SliceOf(reflect.PtrTo(relType))
						relatedItems := reflect.New(sliceType).Interface()

						if err := p.DB.WithContext(stdCtx).Where("id IN ?", ids).Find(relatedItems).Error; err == nil {
							p.DB.WithContext(stdCtx).Model(item).Association(relName).Replace(relatedItems)
						}
					} else {
						// If empty list sent, clear associations
						p.DB.WithContext(stdCtx).Model(item).Association(relName).Clear()
					}
				}

			}
		}
	}

	updates["updated_at"] = time.Now()
	if len(updates) > 0 {
		if err := p.DB.WithContext(stdCtx).Model(item).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	return p.Show(ctx, id)
}

func (p *GormDataProvider) Delete(ctx *context.Context, id string) error {
	stdCtx := p.getContext(ctx)
	return p.DB.WithContext(stdCtx).Model(p.Model).Where("id = ?", id).Delete(nil).Error
}
