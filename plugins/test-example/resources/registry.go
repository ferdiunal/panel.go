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
func Register(slug string, factory func() resource.Resource) {
	mu.Lock()
	defer mu.Unlock()
	registry[slug] = factory
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
