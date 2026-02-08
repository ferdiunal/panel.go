// Package openapi, özel field type mapping'lerini yönetir.
//
// Bu paket, kullanıcıların field type'larını OpenAPI schema'ya nasıl map edileceğini
// özelleştirmelerine olanak tanır:
// - Global field type mapping override
// - Specific field mapping override
// - Resource-level mapping override
package openapi

import (
	"sync"

	"github.com/ferdiunal/panel.go/pkg/core"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/resource"
)

// CustomMappingRegistry, özel field type mapping'lerini yönetir.
//
// ## Özellikler
//   - Global field type mapping: Tüm field'lar için varsayılan mapping override
//   - Specific field mapping: Belirli bir field için mapping override
//   - Resource-level mapping: Belirli bir resource için mapping override
//   - Thread-safe: Concurrent erişim için mutex koruması
//
// ## Kullanım Örneği
//
//	registry := NewCustomMappingRegistry()
//
//	// Global field type mapping
//	registry.RegisterFieldTypeMapping(core.TYPE_TEXT, func(field fields.Element) Schema {
//	    return Schema{Type: "string", MaxLength: ptr(500)}
//	})
//
//	// Specific field mapping
//	registry.RegisterFieldMapping("users", "email", func(field fields.Element) Schema {
//	    return Schema{
//	        Type:    "string",
//	        Format:  "email",
//	        Pattern: "^[a-zA-Z0-9._%+-]+@company\\.com$",
//	    }
//	})
//
//	// Resource-level mapping
//	registry.RegisterResourceMapping("products", func(res resource.Resource) map[string]Schema {
//	    return map[string]Schema{
//	        "price": {Type: "number", Format: "double", Minimum: ptr(0.0)},
//	    }
//	})
type CustomMappingRegistry struct {
	// fieldTypeMappings: Global field type mapping'leri
	// Key: ElementType (örn: core.TYPE_TEXT)
	// Value: Mapping fonksiyonu
	fieldTypeMappings map[core.ElementType]FieldMappingFunc

	// fieldMappings: Specific field mapping'leri
	// Key: "resource:field" (örn: "users:email")
	// Value: Mapping fonksiyonu
	fieldMappings map[string]FieldMappingFunc

	// resourceMappings: Resource-level mapping'ler
	// Key: Resource slug (örn: "users")
	// Value: Resource mapping fonksiyonu
	resourceMappings map[string]ResourceMappingFunc

	mu sync.RWMutex
}

// FieldMappingFunc, bir field'ı OpenAPI schema'ya map eden fonksiyon tipidir.
//
// ## Parametreler
//   - field: Map edilecek field
//
// ## Dönüş Değeri
//   - Schema: OpenAPI schema
//
// ## Kullanım Örneği
//
//	mappingFunc := func(field fields.Element) Schema {
//	    return Schema{
//	        Type:        "string",
//	        MaxLength:   ptr(500),
//	        Description: "Custom description",
//	    }
//	}
type FieldMappingFunc func(field fields.Element) Schema

// ResourceMappingFunc, bir resource için field mapping'lerini döndüren fonksiyon tipidir.
//
// ## Parametreler
//   - res: Resource
//
// ## Dönüş Değeri
//   - map[string]Schema: Field key -> Schema mapping'i
//
// ## Kullanım Örneği
//
//	resourceMappingFunc := func(res resource.Resource) map[string]Schema {
//	    return map[string]Schema{
//	        "price": {
//	            Type:        "number",
//	            Format:      "double",
//	            Minimum:     ptr(0.0),
//	            Description: "Product price in TL",
//	        },
//	        "stock": {
//	            Type:    "integer",
//	            Minimum: ptr(0.0),
//	        },
//	    }
//	}
type ResourceMappingFunc func(res resource.Resource) map[string]Schema

// NewCustomMappingRegistry, yeni bir CustomMappingRegistry oluşturur.
//
// ## Dönüş Değeri
//   - *CustomMappingRegistry: Yapılandırılmış registry
//
// ## Kullanım Örneği
//
//	registry := NewCustomMappingRegistry()
func NewCustomMappingRegistry() *CustomMappingRegistry {
	return &CustomMappingRegistry{
		fieldTypeMappings: make(map[core.ElementType]FieldMappingFunc),
		fieldMappings:     make(map[string]FieldMappingFunc),
		resourceMappings:  make(map[string]ResourceMappingFunc),
	}
}

// RegisterFieldTypeMapping, global field type mapping kaydeder.
//
// ## Parametreler
//   - fieldType: Field type (örn: core.TYPE_TEXT)
//   - mappingFunc: Mapping fonksiyonu
//
// ## Kullanım Örneği
//
//	// Text field'lar için max length'i 500 yap
//	registry.RegisterFieldTypeMapping(core.TYPE_TEXT, func(field fields.Element) Schema {
//	    return Schema{
//	        Type:      "string",
//	        MaxLength: ptr(500),
//	    }
//	})
//
//	// Email field'lar için custom pattern ekle
//	registry.RegisterFieldTypeMapping(core.TYPE_EMAIL, func(field fields.Element) Schema {
//	    return Schema{
//	        Type:    "string",
//	        Format:  "email",
//	        Pattern: "^[a-zA-Z0-9._%+-]+@company\\.com$",
//	    }
//	})
func (r *CustomMappingRegistry) RegisterFieldTypeMapping(fieldType core.ElementType, mappingFunc FieldMappingFunc) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.fieldTypeMappings[fieldType] = mappingFunc
}

// RegisterFieldMapping, specific field mapping kaydeder.
//
// ## Parametreler
//   - resourceSlug: Resource slug (örn: "users")
//   - fieldKey: Field key (örn: "email")
//   - mappingFunc: Mapping fonksiyonu
//
// ## Kullanım Örneği
//
//	// users resource'undaki email field'ı için özel mapping
//	registry.RegisterFieldMapping("users", "email", func(field fields.Element) Schema {
//	    return Schema{
//	        Type:        "string",
//	        Format:      "email",
//	        Pattern:     "^[a-zA-Z0-9._%+-]+@company\\.com$",
//	        Description: "Company email address (must end with @company.com)",
//	        Example:     "user@company.com",
//	    }
//	})
//
//	// products resource'undaki price field'ı için özel mapping
//	registry.RegisterFieldMapping("products", "price", func(field fields.Element) Schema {
//	    return Schema{
//	        Type:        "number",
//	        Format:      "double",
//	        Minimum:     ptr(0.0),
//	        Description: "Product price in TL",
//	        Example:     99.99,
//	    }
//	})
func (r *CustomMappingRegistry) RegisterFieldMapping(resourceSlug, fieldKey string, mappingFunc FieldMappingFunc) {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := resourceSlug + ":" + fieldKey
	r.fieldMappings[key] = mappingFunc
}

// RegisterResourceMapping, resource-level mapping kaydeder.
//
// ## Parametreler
//   - resourceSlug: Resource slug (örn: "products")
//   - mappingFunc: Resource mapping fonksiyonu
//
// ## Kullanım Örneği
//
//	// products resource'u için tüm field'ların mapping'ini özelleştir
//	registry.RegisterResourceMapping("products", func(res resource.Resource) map[string]Schema {
//	    return map[string]Schema{
//	        "price": {
//	            Type:        "number",
//	            Format:      "double",
//	            Minimum:     ptr(0.0),
//	            Description: "Product price in TL",
//	        },
//	        "stock": {
//	            Type:        "integer",
//	            Minimum:     ptr(0.0),
//	            Description: "Stock quantity",
//	        },
//	        "sku": {
//	            Type:        "string",
//	            Pattern:     "^[A-Z0-9-]+$",
//	            Description: "Stock Keeping Unit",
//	        },
//	    }
//	})
func (r *CustomMappingRegistry) RegisterResourceMapping(resourceSlug string, mappingFunc ResourceMappingFunc) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.resourceMappings[resourceSlug] = mappingFunc
}

// GetFieldTypeMapping, global field type mapping'i döndürür.
//
// ## Parametreler
//   - fieldType: Field type
//
// ## Dönüş Değeri
//   - FieldMappingFunc: Mapping fonksiyonu
//   - bool: Mapping bulundu mu?
//
// ## Kullanım Örneği
//
//	if mappingFunc, ok := registry.GetFieldTypeMapping(core.TYPE_TEXT); ok {
//	    schema := mappingFunc(field)
//	}
func (r *CustomMappingRegistry) GetFieldTypeMapping(fieldType core.ElementType) (FieldMappingFunc, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	mappingFunc, ok := r.fieldTypeMappings[fieldType]
	return mappingFunc, ok
}

// GetFieldMapping, specific field mapping'i döndürür.
//
// ## Parametreler
//   - field: Field
//
// ## Dönüş Değeri
//   - Schema: OpenAPI schema
//   - bool: Mapping bulundu mu?
//
// ## Davranış
//   1. Specific field mapping kontrolü (resource:field)
//   2. Global field type mapping kontrolü
//   3. Hiçbiri yoksa false döner
//
// ## Kullanım Örneği
//
//	if schema, ok := registry.GetFieldMapping(field); ok {
//	    // Custom schema kullan
//	} else {
//	    // Varsayılan schema kullan
//	}
func (r *CustomMappingRegistry) GetFieldMapping(field fields.Element) (Schema, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Specific field mapping kontrolü
	// Not: Resource slug'ını field'dan alamıyoruz, bu yüzden sadece field type mapping'i kontrol ediyoruz
	// Specific field mapping için GetFieldMappingByKey metodunu kullanın

	// Global field type mapping kontrolü
	if mappingFunc, ok := r.fieldTypeMappings[field.GetType()]; ok {
		return mappingFunc(field), true
	}

	return Schema{}, false
}

// GetFieldMappingByKey, specific field mapping'i resource ve field key ile döndürür.
//
// ## Parametreler
//   - resourceSlug: Resource slug
//   - fieldKey: Field key
//   - field: Field
//
// ## Dönüş Değeri
//   - Schema: OpenAPI schema
//   - bool: Mapping bulundu mu?
//
// ## Kullanım Örneği
//
//	if schema, ok := registry.GetFieldMappingByKey("users", "email", field); ok {
//	    // Custom schema kullan
//	}
func (r *CustomMappingRegistry) GetFieldMappingByKey(resourceSlug, fieldKey string, field fields.Element) (Schema, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Specific field mapping kontrolü
	key := resourceSlug + ":" + fieldKey
	if mappingFunc, ok := r.fieldMappings[key]; ok {
		return mappingFunc(field), true
	}

	// Global field type mapping kontrolü
	if mappingFunc, ok := r.fieldTypeMappings[field.GetType()]; ok {
		return mappingFunc(field), true
	}

	return Schema{}, false
}

// GetResourceMapping, resource-level mapping'i döndürür.
//
// ## Parametreler
//   - resourceSlug: Resource slug
//   - res: Resource
//
// ## Dönüş Değeri
//   - map[string]Schema: Field key -> Schema mapping'i
//   - bool: Mapping bulundu mu?
//
// ## Kullanım Örneği
//
//	if fieldMappings, ok := registry.GetResourceMapping("products", resource); ok {
//	    for fieldKey, schema := range fieldMappings {
//	        // Custom schema kullan
//	    }
//	}
func (r *CustomMappingRegistry) GetResourceMapping(resourceSlug string, res resource.Resource) (map[string]Schema, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if mappingFunc, ok := r.resourceMappings[resourceSlug]; ok {
		return mappingFunc(res), true
	}

	return nil, false
}

// HasFieldTypeMapping, global field type mapping'in olup olmadığını kontrol eder.
//
// ## Parametreler
//   - fieldType: Field type
//
// ## Dönüş Değeri
//   - bool: Mapping var mı?
func (r *CustomMappingRegistry) HasFieldTypeMapping(fieldType core.ElementType) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, ok := r.fieldTypeMappings[fieldType]
	return ok
}

// HasFieldMapping, specific field mapping'in olup olmadığını kontrol eder.
//
// ## Parametreler
//   - resourceSlug: Resource slug
//   - fieldKey: Field key
//
// ## Dönüş Değeri
//   - bool: Mapping var mı?
func (r *CustomMappingRegistry) HasFieldMapping(resourceSlug, fieldKey string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	key := resourceSlug + ":" + fieldKey
	_, ok := r.fieldMappings[key]
	return ok
}

// HasResourceMapping, resource-level mapping'in olup olmadığını kontrol eder.
//
// ## Parametreler
//   - resourceSlug: Resource slug
//
// ## Dönüş Değeri
//   - bool: Mapping var mı?
func (r *CustomMappingRegistry) HasResourceMapping(resourceSlug string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, ok := r.resourceMappings[resourceSlug]
	return ok
}

// ClearFieldTypeMapping, global field type mapping'i temizler.
//
// ## Parametreler
//   - fieldType: Field type
//
// ## Kullanım Örneği
//
//	registry.ClearFieldTypeMapping(core.TYPE_TEXT)
func (r *CustomMappingRegistry) ClearFieldTypeMapping(fieldType core.ElementType) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.fieldTypeMappings, fieldType)
}

// ClearFieldMapping, specific field mapping'i temizler.
//
// ## Parametreler
//   - resourceSlug: Resource slug
//   - fieldKey: Field key
//
// ## Kullanım Örneği
//
//	registry.ClearFieldMapping("users", "email")
func (r *CustomMappingRegistry) ClearFieldMapping(resourceSlug, fieldKey string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := resourceSlug + ":" + fieldKey
	delete(r.fieldMappings, key)
}

// ClearResourceMapping, resource-level mapping'i temizler.
//
// ## Parametreler
//   - resourceSlug: Resource slug
//
// ## Kullanım Örneği
//
//	registry.ClearResourceMapping("products")
func (r *CustomMappingRegistry) ClearResourceMapping(resourceSlug string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.resourceMappings, resourceSlug)
}

// ClearAll, tüm mapping'leri temizler.
//
// ## Kullanım Örneği
//
//	registry.ClearAll()
func (r *CustomMappingRegistry) ClearAll() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.fieldTypeMappings = make(map[core.ElementType]FieldMappingFunc)
	r.fieldMappings = make(map[string]FieldMappingFunc)
	r.resourceMappings = make(map[string]ResourceMappingFunc)
}

// Count, kayıtlı mapping sayısını döndürür.
//
// ## Dönüş Değeri
//   - int: Field type mapping sayısı
//   - int: Field mapping sayısı
//   - int: Resource mapping sayısı
//
// ## Kullanım Örneği
//
//	fieldTypeCount, fieldCount, resourceCount := registry.Count()
//	fmt.Printf("Mappings: %d field types, %d fields, %d resources\n",
//	    fieldTypeCount, fieldCount, resourceCount)
func (r *CustomMappingRegistry) Count() (int, int, int) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.fieldTypeMappings), len(r.fieldMappings), len(r.resourceMappings)
}
