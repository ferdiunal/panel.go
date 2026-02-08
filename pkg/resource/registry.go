package resource

import (
	"sync"
)

// registry - Resource factory'lerini saklayan merkezi kayıt
// Circular dependency problemini çözmek için kullanılır
var (
	registry = make(map[string]func() Resource)
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
// func init() {
//     resource.Register("organizations", func() resource.Resource {
//         return NewOrganizationResource()
//     })
// }
// ```
//
// # Parametreler
//
// - **slug**: Resource'un benzersiz slug'ı (örn: "organizations", "addresses")
// - **factory**: Resource instance'ı döndüren factory fonksiyonu
func Register(slug string, factory func() Resource) {
	mu.Lock()
	defer mu.Unlock()
	registry[slug] = factory
}

// Get kayıtlı bir resource'u slug'ına göre alır
//
// Eğer resource kayıtlı değilse nil döner.
//
// Kullanım örneği:
// ```go
// orgResource := resource.Get("organizations")
// if orgResource != nil {
//     // Resource kullan
// }
// ```
//
// # Parametreler
//
// - **slug**: Resource'un slug'ı
//
// # Döndürür
//
// - Resource: Resource instance'ı veya nil
func Get(slug string) Resource {
	mu.RLock()
	defer mu.RUnlock()
	if factory, ok := registry[slug]; ok {
		return factory()
	}
	return nil
}

// GetOrPanic kayıtlı bir resource'u alır, bulamazsa panic yapar
//
// Bu fonksiyon, resource'un mutlaka kayıtlı olması gereken durumlarda kullanılır.
// Eğer resource kayıtlı değilse, uygulama başlatma sırasında hata verir.
//
// Kullanım örneği:
// ```go
// // ResolveFields içinde kullanım
// fields.BelongsTo("Organization", "organization_id", resource.GetOrPanic("organizations"))
// ```
//
// # Parametreler
//
// - **slug**: Resource'un slug'ı
//
// # Döndürür
//
// - Resource: Resource instance'ı
//
// # Panic
//
// - Resource kayıtlı değilse panic yapar
func GetOrPanic(slug string) Resource {
	r := Get(slug)
	if r == nil {
		panic("resource not found: " + slug)
	}
	return r
}

// List tüm kayıtlı resource slug'larını döndürür
//
// Debug ve test amaçlı kullanılır.
//
// # Döndürür
//
// - []string: Kayıtlı resource slug'larının listesi
func List() []string {
	mu.RLock()
	defer mu.RUnlock()
	slugs := make([]string, 0, len(registry))
	for slug := range registry {
		slugs = append(slugs, slug)
	}
	return slugs
}

// Clear tüm kayıtlı resource'ları temizler
//
// Test amaçlı kullanılır.
func Clear() {
	mu.Lock()
	defer mu.Unlock()
	registry = make(map[string]func() Resource)
}
