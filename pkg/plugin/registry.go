/// # Plugin Registry
///
/// Thread-safe plugin registry. Plugin'leri kaydetmek ve yönetmek için kullanılır.
///
/// ## Özellikler
/// - Thread-safe (RWMutex ile korunur)
/// - Global registry (singleton pattern)
/// - Manuel registration (init() fonksiyonunda)
/// - Plugin listesi ve arama
///
/// ## Kullanım Örneği
/// ```go
/// // Plugin kaydı (init fonksiyonunda)
/// func init() {
///     plugin.Register(&MyPlugin{})
/// }
///
/// // Plugin listesi
/// plugins := plugin.All()
///
/// // Plugin arama
/// p := plugin.Get("my-plugin")
/// ```

package plugin

import (
	"fmt"
	"sync"
)

/// # Registry Struct
///
/// Thread-safe plugin registry. Tüm kayıtlı plugin'leri tutar.
///
/// ## Alanlar
/// - `plugins`: Plugin listesi (slice)
/// - `pluginMap`: Plugin haritası (name -> Plugin)
/// - `mu`: Read-Write mutex (thread-safety için)
///
/// ## Thread Safety
/// Tüm public metodlar RWMutex ile korunur.
/// Read işlemleri RLock, write işlemleri Lock kullanır.
type Registry struct {
	plugins   []Plugin
	pluginMap map[string]Plugin
	mu        sync.RWMutex
}

// Global registry instance (singleton)
var globalRegistry = &Registry{
	plugins:   make([]Plugin, 0),
	pluginMap: make(map[string]Plugin),
}

/// # Register Fonksiyonu
///
/// Plugin'i global registry'ye kaydeder.
/// Genellikle plugin'in init() fonksiyonunda çağrılır.
///
/// ## Parametreler
/// - `p`: Kaydedilecek plugin
///
/// ## Hata Durumları
/// - Plugin nil ise panic oluşturur
/// - Aynı isimde plugin zaten kayıtlıysa panic oluşturur
///
/// ## Kullanım Örneği
/// ```go
/// func init() {
///     plugin.Register(&MyPlugin{})
/// }
/// ```
///
/// ## Önemli Notlar
/// - Bu fonksiyon thread-safe'dir
/// - init() fonksiyonunda çağrılmalıdır
/// - Aynı plugin birden fazla kez kaydedilemez
func Register(p Plugin) {
	if p == nil {
		panic("plugin: cannot register nil plugin")
	}

	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()

	name := p.Name()
	if _, exists := globalRegistry.pluginMap[name]; exists {
		panic(fmt.Sprintf("plugin: plugin '%s' is already registered", name))
	}

	globalRegistry.plugins = append(globalRegistry.plugins, p)
	globalRegistry.pluginMap[name] = p
}

/// # All Fonksiyonu
///
/// Tüm kayıtlı plugin'leri döndürür.
///
/// ## Dönüş Değeri
/// - `[]Plugin`: Kayıtlı plugin'lerin kopyası
///
/// ## Kullanım Örneği
/// ```go
/// plugins := plugin.All()
/// for _, p := range plugins {
///     fmt.Println(p.Name(), p.Version())
/// }
/// ```
///
/// ## Önemli Notlar
/// - Thread-safe'dir
/// - Slice'ın kopyasını döndürür (orijinal değiştirilmez)
func All() []Plugin {
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()

	// Return a copy to prevent external modification
	result := make([]Plugin, len(globalRegistry.plugins))
	copy(result, globalRegistry.plugins)
	return result
}

/// # Get Fonksiyonu
///
/// İsme göre plugin arar ve döndürür.
///
/// ## Parametreler
/// - `name`: Plugin adı
///
/// ## Dönüş Değeri
/// - `Plugin`: Bulunan plugin (nil ise bulunamadı)
///
/// ## Kullanım Örneği
/// ```go
/// p := plugin.Get("analytics-plugin")
/// if p != nil {
///     fmt.Println("Found:", p.Version())
/// }
/// ```
///
/// ## Önemli Notlar
/// - Thread-safe'dir
/// - Plugin bulunamazsa nil döner
func Get(name string) Plugin {
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()

	return globalRegistry.pluginMap[name]
}

/// # Count Fonksiyonu
///
/// Kayıtlı plugin sayısını döndürür.
///
/// ## Dönüş Değeri
/// - `int`: Kayıtlı plugin sayısı
///
/// ## Kullanım Örneği
/// ```go
/// count := plugin.Count()
/// fmt.Printf("Registered plugins: %d\n", count)
/// ```
func Count() int {
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()

	return len(globalRegistry.plugins)
}

/// # Clear Fonksiyonu
///
/// Tüm kayıtlı plugin'leri temizler.
/// Genellikle test senaryolarında kullanılır.
///
/// ## Uyarı
/// Bu fonksiyon production ortamında kullanılmamalıdır!
/// Sadece test amaçlı kullanılmalıdır.
///
/// ## Kullanım Örneği
/// ```go
/// // Test setup
/// func TestMyPlugin(t *testing.T) {
///     plugin.Clear() // Temiz başla
///     plugin.Register(&MyPlugin{})
///     // ... test
/// }
/// ```
func Clear() {
	globalRegistry.mu.Lock()
	defer globalRegistry.mu.Unlock()

	globalRegistry.plugins = make([]Plugin, 0)
	globalRegistry.pluginMap = make(map[string]Plugin)
}

/// # Exists Fonksiyonu
///
/// Plugin'in kayıtlı olup olmadığını kontrol eder.
///
/// ## Parametreler
/// - `name`: Plugin adı
///
/// ## Dönüş Değeri
/// - `bool`: Plugin kayıtlıysa true, değilse false
///
/// ## Kullanım Örneği
/// ```go
/// if plugin.Exists("analytics-plugin") {
///     fmt.Println("Plugin is registered")
/// }
/// ```
func Exists(name string) bool {
	globalRegistry.mu.RLock()
	defer globalRegistry.mu.RUnlock()

	_, exists := globalRegistry.pluginMap[name]
	return exists
}
