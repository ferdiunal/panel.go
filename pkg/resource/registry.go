package resource

import (
	"sync"
)

// registry, tüm resource'ların global registry'si.
//
// Bu registry, resource'ların slug'larına göre saklanmasını ve erişilmesini sağlar.
// Circular dependency sorununu önlemek için string slug kullanılır ve
// runtime'da resource instance'ları registry'den alınır.
//
// # Thread Safety
//
// Registry, concurrent erişim için RWMutex kullanır:
// - Register: Write lock (resource ekleme)
// - Get: Read lock (resource okuma)
//
// # Kullanım Senaryoları
//
// 1. **Resource Registration**: Her resource kendi init() fonksiyonunda kendini register eder
// 2. **Relationship Resolution**: HasMany, BelongsTo gibi field'lar RelatedResourceSlug'dan resource instance'ını alır
// 3. **Dynamic Resource Loading**: Runtime'da resource'lara erişim
//
// # Örnek Kullanım
//
//	// Resource registration (init fonksiyonunda)
//	func init() {
//	    resource.Register("users", NewUserResource())
//	}
//
//	// Resource retrieval
//	userResource := resource.Get("users")
//	if userResource != nil {
//	    // Resource bulundu, kullan
//	}
var registry = struct {
	sync.RWMutex
	resources map[string]Resource
}{
	resources: make(map[string]Resource),
}

// Register, bir resource'u registry'ye ekler.
//
// Bu fonksiyon, resource'ların slug'larına göre saklanmasını sağlar.
// Aynı slug ile birden fazla resource register edilirse, son register edilen geçerli olur.
//
// # Thread Safety
//
// Bu fonksiyon thread-safe'tir. Concurrent çağrılar güvenlidir.
//
// # Parametreler
//
// - slug: Resource'un benzersiz slug'ı (örn. "users", "posts", "organizations")
// - res: Resource instance'ı
//
// # Kullanım Örneği
//
//	// Resource registration (init fonksiyonunda)
//	func init() {
//	    resource.Register("users", NewUserResource())
//	    resource.Register("posts", NewPostResource())
//	    resource.Register("organizations", NewOrganizationResource())
//	}
//
// # Önemli Notlar
//
// - Register fonksiyonu genellikle init() fonksiyonunda çağrılır
// - Aynı slug ile birden fazla register edilirse, son register edilen geçerli olur
// - nil resource register edilebilir (ama önerilmez)
func Register(slug string, res Resource) {
	registry.Lock()
	defer registry.Unlock()
	registry.resources[slug] = res
}

// Get, slug'a göre bir resource'u registry'den alır.
//
// Bu fonksiyon, relationship field'larının RelatedResource'unu çözmek için kullanılır.
// Eğer slug bulunamazsa nil döner.
//
// # Thread Safety
//
// Bu fonksiyon thread-safe'tir. Concurrent çağrılar güvenlidir.
//
// # Parametreler
//
// - slug: Resource'un benzersiz slug'ı (örn. "users", "posts", "organizations")
//
// # Dönüş Değeri
//
// - Resource instance'ı (bulunursa)
// - nil (bulunamazsa)
//
// # Kullanım Örneği
//
//	// Relationship field'ında RelatedResource'u çözme
//	relatedResource := resource.Get("users")
//	if relatedResource != nil {
//	    // Resource bulundu, kullan
//	    title := relatedResource.RecordTitle(record)
//	} else {
//	    // Resource bulunamadı
//	}
//
// # Önemli Notlar
//
// - Eğer slug bulunamazsa nil döner (panic atmaz)
// - nil kontrolü yapılmalıdır
// - Registry'de olmayan bir slug için Get çağrılması hata değildir
func Get(slug string) Resource {
	registry.RLock()
	defer registry.RUnlock()
	return registry.resources[slug]
}

// List, registry'deki tüm resource'ları döndürür.
//
// Bu fonksiyon, debug ve test amaçlı kullanılabilir.
//
// # Thread Safety
//
// Bu fonksiyon thread-safe'tir. Concurrent çağrılar güvenlidir.
//
// # Dönüş Değeri
//
// - map[string]Resource: Slug -> Resource mapping
//
// # Kullanım Örneği
//
//	// Tüm resource'ları listele
//	resources := resource.List()
//	for slug, res := range resources {
//	    fmt.Printf("Resource: %s -> %s\n", slug, res.Title())
//	}
//
// # Önemli Notlar
//
// - Döndürülen map bir kopyadır, değiştirmek registry'yi etkilemez
// - Bu fonksiyon genellikle debug amaçlı kullanılır
func List() map[string]Resource {
	registry.RLock()
	defer registry.RUnlock()

	// Kopya oluştur (registry'yi korumak için)
	result := make(map[string]Resource, len(registry.resources))
	for slug, res := range registry.resources {
		result[slug] = res
	}
	return result
}

// Clear, registry'deki tüm resource'ları temizler.
//
// Bu fonksiyon, test amaçlı kullanılır. Production'da kullanılmamalıdır.
//
// # Thread Safety
//
// Bu fonksiyon thread-safe'tir. Concurrent çağrılar güvenlidir.
//
// # Kullanım Örneği
//
//	// Test setup
//	func TestSetup(t *testing.T) {
//	    resource.Clear()
//	    resource.Register("users", NewMockUserResource())
//	}
//
// # Önemli Notlar
//
// - Bu fonksiyon sadece test amaçlı kullanılmalıdır
// - Production'da kullanılması önerilmez
func Clear() {
	registry.Lock()
	defer registry.Unlock()
	registry.resources = make(map[string]Resource)
}
