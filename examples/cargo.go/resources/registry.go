package resources

import (
	"sync"

	"github.com/ferdiunal/panel.go/pkg/resource"
)

// registry - Resource factory'lerini saklayan merkezi kayıt
// Circular dependency problemini çözmek için kullanılır
var (
	registry = make(map[string]func() resource.Resource)
	mu       sync.RWMutex
)

// Register bir resource factory'sini kayıt eder
//
// Bu fonksiyon, resource'ların init() fonksiyonlarında çağrılır.
// Circular dependency'yi önlemek için resource'lar birbirini import etmez,
// bunun yerine registry'den alır.
//
// Kullanım örneği:
// ```go
//
//	func init() {
//	    resources.Register("organizations", func() resource.Resource {
//	        return NewOrganizationResource()
//	    })
//	}
//
// ```
func Register(slug string, factory func() resource.Resource) {
	mu.Lock()
	defer mu.Unlock()
	registry[slug] = factory

	// Global registry'ye de kaydet (panel.go kütüphanesinin erişebilmesi için)
	// Bu sayede AutoOptions gibi özellikler resource'u ve modelini bulabilir
	resource.Register(slug, factory())
}

// Get kayıtlı bir resource'u slug'ına göre alır
func Get(slug string) resource.Resource {
	mu.RLock()
	defer mu.RUnlock()
	if factory, ok := registry[slug]; ok {
		return factory()
	}
	return nil
}

// GetOrPanic kayıtlı bir resource'u alır, bulamazsa panic yapar
func GetOrPanic(slug string) resource.Resource {
	r := Get(slug)
	if r == nil {
		panic("resource not found: " + slug)
	}
	return r
}

// List tüm kayıtlı resource slug'larını döndürür
func List() []string {
	mu.RLock()
	defer mu.RUnlock()
	slugs := make([]string, 0, len(registry))
	for slug := range registry {
		slugs = append(slugs, slug)
	}
	return slugs
}

// Clear tüm kayıtlı resource'ları temizler (test amaçlı)
func Clear() {
	mu.Lock()
	defer mu.Unlock()
	registry = make(map[string]func() resource.Resource)
}
