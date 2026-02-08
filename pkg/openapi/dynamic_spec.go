// Package openapi, resource'lar için dinamik OpenAPI spesifikasyonu oluşturur.
//
// Bu paket, Panel.go resource'larını otomatik olarak OpenAPI path ve schema'larına çevirir:
// - Resource field'larını OpenAPI schema'ya map eder
// - CRUD endpoint'leri oluşturur (List, Get, Create, Update, Delete)
// - Action endpoint'leri oluşturur
// - Filter ve sort parametrelerini ekler
package openapi

import (
	"fmt"
	"reflect"

	"github.com/ferdiunal/panel.go/pkg/core"
	"github.com/ferdiunal/panel.go/pkg/resource"
)

// DynamicSpecGenerator, resource'lar için dinamik OpenAPI spesifikasyonu oluşturur.
//
// ## Özellikler
// - Resource field'larını otomatik schema'ya çevirir
// - CRUD endpoint'leri oluşturur
// - Action endpoint'leri ekler
// - Filter ve sort parametrelerini tanımlar
//
// ## Kullanım Örneği
//
//	generator := NewDynamicSpecGenerator()
//	paths := generator.GenerateResourcePaths(resources)
//	schemas := generator.GenerateResourceSchemas(resources)
type DynamicSpecGenerator struct {
	mapper *FieldTypeMapper
}

// NewDynamicSpecGenerator, yeni bir DynamicSpecGenerator oluşturur.
func NewDynamicSpecGenerator() *DynamicSpecGenerator {
	return &DynamicSpecGenerator{
		mapper: NewFieldTypeMapper(),
	}
}

// GenerateResourcePaths, tüm resource'lar için OpenAPI path'leri oluşturur.
//
// ## Parametreler
// - resources: Resource map'i (slug -> Resource)
//
// ## Döndürür
// - OpenAPI Paths objesi
func (g *DynamicSpecGenerator) GenerateResourcePaths(resources map[string]resource.Resource) Paths {
	paths := make(Paths)

	for slug, res := range resources {
		// OpenAPI'de görünür değilse skip et
		if !res.OpenAPIEnabled() {
			continue
		}

		// Collection endpoint: GET /api/resources/{slug}
		collectionPath := fmt.Sprintf("/api/resources/%s", slug)
		paths[collectionPath] = *g.generateCollectionPathItem(res)

		// Item endpoint: GET /api/resources/{slug}/{id}
		itemPath := fmt.Sprintf("/api/resources/%s/{id}", slug)
		paths[itemPath] = *g.generateItemPathItem(res)

		// Action endpoints
		for _, action := range res.GetActions() {
			actionPath := fmt.Sprintf("/api/resources/%s/actions/%s", slug, action.GetSlug())
			paths[actionPath] = *g.generateActionPathItem(res, action)
		}
	}

	return paths
}

// generateCollectionPathItem, collection endpoint için PathItem oluşturur.
func (g *DynamicSpecGenerator) generateCollectionPathItem(res resource.Resource) *PathItem {
	schemaName := g.getSchemaName(res)

	return &PathItem{
		Get: &Operation{
			Summary:     fmt.Sprintf("List %s", res.Title()),
			Description: fmt.Sprintf("Retrieve a paginated list of %s", res.Title()),
			Tags:        []string{res.Title()},
			Parameters:  g.generateListParameters(res),
			Responses: Responses{
				"200": Response{
					Description: "Successful response",
					Content: map[string]MediaType{
						"application/json": MediaType{
							Schema: &Schema{
								Type: "object",
								Properties: map[string]Schema{
									"data": Schema{
										Type: "array",
										Items: &Schema{
											Ref: fmt.Sprintf("#/components/schemas/%s", schemaName),
										},
									},
									"meta": Schema{
										Ref: "#/components/schemas/PaginationMeta",
									},
								},
							},
						},
					},
				},
			},
		},
		Post: &Operation{
			Summary:     fmt.Sprintf("Create %s", res.Title()),
			Description: fmt.Sprintf("Create a new %s", res.Title()),
			Tags:        []string{res.Title()},
			RequestBody: &RequestBody{
				Required: true,
				Content: map[string]MediaType{
					"application/json": MediaType{
						Schema: &Schema{
							Ref: fmt.Sprintf("#/components/schemas/%sInput", schemaName),
						},
					},
				},
			},
			Responses: Responses{
				"201": Response{
					Description: "Resource created successfully",
					Content: map[string]MediaType{
						"application/json": MediaType{
							Schema: &Schema{
								Ref: fmt.Sprintf("#/components/schemas/%s", schemaName),
							},
						},
					},
				},
			},
		},
	}
}

// generateItemPathItem, item endpoint için PathItem oluşturur.
func (g *DynamicSpecGenerator) generateItemPathItem(res resource.Resource) *PathItem {
	schemaName := g.getSchemaName(res)

	return &PathItem{
		Get: &Operation{
			Summary:     fmt.Sprintf("Get %s", res.Title()),
			Description: fmt.Sprintf("Retrieve a single %s by ID", res.Title()),
			Tags:        []string{res.Title()},
			Parameters: []Parameter{
				{
					Name:        "id",
					In:          "path",
					Required:    true,
					Description: "Resource ID",
					Schema: &Schema{
						Type: "integer",
					},
				},
			},
			Responses: Responses{
				"200": Response{
					Description: "Successful response",
					Content: map[string]MediaType{
						"application/json": MediaType{
							Schema: &Schema{
								Ref: fmt.Sprintf("#/components/schemas/%s", schemaName),
							},
						},
					},
				},
			},
		},
		Put: &Operation{
			Summary:     fmt.Sprintf("Update %s", res.Title()),
			Description: fmt.Sprintf("Update an existing %s", res.Title()),
			Tags:        []string{res.Title()},
			Parameters: []Parameter{
				{
					Name:        "id",
					In:          "path",
					Required:    true,
					Description: "Resource ID",
					Schema: &Schema{
						Type: "integer",
					},
				},
			},
			RequestBody: &RequestBody{
				Required: true,
				Content: map[string]MediaType{
					"application/json": MediaType{
						Schema: &Schema{
							Ref: fmt.Sprintf("#/components/schemas/%sInput", schemaName),
						},
					},
				},
			},
			Responses: Responses{
				"200": Response{
					Description: "Resource updated successfully",
					Content: map[string]MediaType{
						"application/json": MediaType{
							Schema: &Schema{
								Ref: fmt.Sprintf("#/components/schemas/%s", schemaName),
							},
						},
					},
				},
			},
		},
		Delete: &Operation{
			Summary:     fmt.Sprintf("Delete %s", res.Title()),
			Description: fmt.Sprintf("Delete a %s", res.Title()),
			Tags:        []string{res.Title()},
			Parameters: []Parameter{
				{
					Name:        "id",
					In:          "path",
					Required:    true,
					Description: "Resource ID",
					Schema: &Schema{
						Type: "integer",
					},
				},
			},
			Responses: Responses{
				"204": Response{
					Description: "Resource deleted successfully",
				},
			},
		},
	}
}

// generateActionPathItem, action endpoint için PathItem oluşturur.
func (g *DynamicSpecGenerator) generateActionPathItem(res resource.Resource, action resource.Action) *PathItem {
	return &PathItem{
		Post: &Operation{
			Summary:     action.GetName(),
			Description: fmt.Sprintf("Execute %s action on selected items", action.GetName()),
			Tags:        []string{res.Title()},
			RequestBody: &RequestBody{
				Required: true,
				Content: map[string]MediaType{
					"application/json": MediaType{
						Schema: &Schema{
							Type: "object",
							Properties: map[string]Schema{
								"ids": {
									Type: "array",
									Items: &Schema{
										Type: "integer",
									},
									Description: "IDs of items to perform action on",
								},
							},
							Required: []string{"ids"},
						},
					},
				},
			},
			Responses: Responses{
				"200": Response{
					Description: "Action executed successfully",
					Content: map[string]MediaType{
						"application/json": MediaType{
							Schema: &Schema{
								Type: "object",
								Properties: map[string]Schema{
									"message": {
										Type: "string",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// generateListParameters, liste endpoint'i için parametreleri oluşturur.
func (g *DynamicSpecGenerator) generateListParameters(res resource.Resource) []Parameter {
	params := []Parameter{
		{
			Name:        "page",
			In:          "query",
			Description: "Page number",
			Schema: &Schema{
				Type:    "integer",
				Default: 1,
			},
		},
		{
			Name:        "per_page",
			In:          "query",
			Description: "Items per page",
			Schema: &Schema{
				Type:    "integer",
				Default: 15,
			},
		},
		{
			Name:        "search",
			In:          "query",
			Description: "Search query",
			Schema: &Schema{
				Type: "string",
			},
		},
		{
			Name:        "sort",
			In:          "query",
			Description: "Sort field",
			Schema: &Schema{
				Type: "string",
			},
		},
		{
			Name:        "order",
			In:          "query",
			Description: "Sort order (asc/desc)",
			Schema: &Schema{
				Type: "string",
				Enum: []interface{}{"asc", "desc"},
			},
		},
	}

	// Filter parametrelerini ekle
	for _, filter := range res.GetFilters() {
		params = append(params, Parameter{
			Name:        filter.GetSlug(),
			In:          "query",
			Description: filter.GetName(),
			Schema: &Schema{
				Type: "string",
			},
		})
	}

	return params
}

// GenerateResourceSchemas, tüm resource'lar için OpenAPI schema'ları oluşturur.
//
// ## Parametreler
// - resources: Resource map'i (slug -> Resource)
//
// ## Döndürür
// - OpenAPI Schemas objesi
func (g *DynamicSpecGenerator) GenerateResourceSchemas(resources map[string]resource.Resource) map[string]*Schema {
	schemas := make(map[string]*Schema)

	// Pagination meta schema
	schemas["PaginationMeta"] = &Schema{
		Type: "object",
		Properties: map[string]Schema{
			"current_page": {Type: "integer"},
			"from":         {Type: "integer"},
			"last_page":    {Type: "integer"},
			"per_page":     {Type: "integer"},
			"to":           {Type: "integer"},
			"total":        {Type: "integer"},
		},
	}

	// Her resource için schema oluştur
	for _, res := range resources {
		// OpenAPI'de görünür değilse skip et
		if !res.OpenAPIEnabled() {
			continue
		}

		schemaName := g.getSchemaName(res)

		// Output schema (tüm field'lar)
		schemas[schemaName] = g.generateResourceSchema(res, false)

		// Input schema (sadece form field'ları)
		schemas[schemaName+"Input"] = g.generateResourceSchema(res, true)
	}

	return schemas
}

// generateResourceSchema, bir resource için OpenAPI schema oluşturur.
func (g *DynamicSpecGenerator) generateResourceSchema(res resource.Resource, inputOnly bool) *Schema {
	schema := &Schema{
		Type:       "object",
		Properties: make(map[string]Schema),
		Required:   []string{},
	}

	fields := res.Fields()
	for _, field := range fields {
		element, ok := field.(core.Element)
		if !ok {
			continue
		}

		// Input schema için sadece form field'larını al
		if inputOnly {
			// Context kontrolü - form'da görünmüyorsa skip et
			// Bu basitleştirilmiş bir kontrol, gerçek implementasyonda
			// element.GetContext() ile kontrol edilmeli
			if element.GetKey() == "id" || element.GetKey() == "created_at" || element.GetKey() == "updated_at" {
				continue
			}
		}

		// Field'ı schema'ya çevir
		fieldSchema := g.mapper.MapFieldToSchema(element)
		if fieldSchema != nil {
			schema.Properties[element.GetKey()] = fieldSchema

			// Required kontrolü
			// Bu basitleştirilmiş bir kontrol, gerçek implementasyonda
			// validation rules kontrol edilmeli
			if inputOnly {
				// Required field'ları ekle (basitleştirilmiş)
				// Gerçek implementasyonda validation rules'dan alınmalı
			}
		}
	}

	return schema
}

// getSchemaName, resource için schema adı oluşturur.
func (g *DynamicSpecGenerator) getSchemaName(res resource.Resource) string {
	// Model tipinden schema adı oluştur
	model := res.Model()
	modelType := reflect.TypeOf(model)

	// Pointer ise underlying type'ı al
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	return modelType.Name()
}
