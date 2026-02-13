/// # Plugin Discovery
///
/// Auto-discovery sistemi. Plugin'leri otomatik olarak keşfeder ve yükler.
/// Bu özellik opsiyoneldir, manuel import tercih edilir.
///
/// ## Özellikler
/// - plugin.yaml dosyası ile plugin tanımlama
/// - Klasör bazlı discovery
/// - Lazy loading desteği
///
/// ## Kullanım Örneği
/// ```go
/// // plugin.yaml dosyası:
/// // name: my-plugin
/// // version: 1.0.0
/// // author: Author Name
/// // description: Plugin description
/// // entry: plugin.so
///
/// // Discovery
/// if err := plugin.Discover("./plugins"); err != nil {
///     log.Fatal(err)
/// }
/// ```
///
/// ## Önemli Notlar
/// - Auto-discovery opsiyoneldir
/// - Manuel import (compile-time) tercih edilir
/// - Production ortamında dikkatli kullanılmalıdır

package plugin

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

/// # PluginDescriptor Struct
///
/// plugin.yaml dosyasının yapısı.
///
/// ## Alanlar
/// - `Metadata`: Plugin meta bilgileri
/// - `Entry`: Plugin entry point (örn: "plugin.so", "main.go")
/// - `Enabled`: Plugin etkin mi?
///
/// ## Kullanım Örneği
/// ```yaml
/// name: my-plugin
/// version: 1.0.0
/// author: Author Name
/// description: Plugin description
/// entry: plugin.so
/// enabled: true
/// ```
type PluginDescriptor struct {
	Metadata `yaml:",inline"`
	Entry    string `json:"entry" yaml:"entry"`
	Enabled  bool   `json:"enabled" yaml:"enabled"`
}

/// # Discover Fonksiyonu
///
/// Verilen klasörde plugin.yaml dosyalarını arar ve plugin'leri keşfeder.
///
/// ## Parametreler
/// - `path`: Plugin klasörü yolu
///
/// ## Dönüş Değeri
/// - `[]PluginDescriptor`: Keşfedilen plugin descriptor'ları
/// - `error`: Hata varsa hata, yoksa nil
///
/// ## Davranış
/// 1. Verilen klasörü tarar
/// 2. Her alt klasörde plugin.yaml arar
/// 3. plugin.yaml dosyalarını parse eder
/// 4. Descriptor listesi döndürür
///
/// ## Kullanım Örneği
/// ```go
/// descriptors, err := plugin.Discover("./plugins")
/// if err != nil {
///     log.Fatal(err)
/// }
///
/// for _, desc := range descriptors {
///     fmt.Println(desc.Name, desc.Version)
/// }
/// ```
///
/// ## Önemli Notlar
/// - Bu fonksiyon plugin'leri yüklemez, sadece keşfeder
/// - Plugin yükleme manuel olarak yapılmalıdır
/// - Production ortamında dikkatli kullanılmalıdır
func Discover(path string) ([]PluginDescriptor, error) {
	var descriptors []PluginDescriptor

	// Klasör var mı kontrol et
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return descriptors, fmt.Errorf("plugin discovery: path does not exist: %s", path)
	}

	// Alt klasörleri tara
	entries, err := os.ReadDir(path)
	if err != nil {
		return descriptors, fmt.Errorf("plugin discovery: failed to read directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// plugin.yaml dosyasını ara
		pluginYamlPath := filepath.Join(path, entry.Name(), "plugin.yaml")
		if _, err := os.Stat(pluginYamlPath); os.IsNotExist(err) {
			continue
		}

		// plugin.yaml dosyasını parse et
		descriptor, err := loadPluginDescriptor(pluginYamlPath)
		if err != nil {
			// Hatalı plugin.yaml dosyalarını atla
			fmt.Printf("Warning: failed to load plugin descriptor from %s: %v\n", pluginYamlPath, err)
			continue
		}

		descriptors = append(descriptors, descriptor)
	}

	return descriptors, nil
}

/// # loadPluginDescriptor Fonksiyonu
///
/// plugin.yaml dosyasını okur ve parse eder.
///
/// ## Parametreler
/// - `path`: plugin.yaml dosyasının yolu
///
/// ## Dönüş Değeri
/// - `PluginDescriptor`: Parse edilmiş descriptor
/// - `error`: Hata varsa hata, yoksa nil
///
/// ## Kullanım Örneği
/// ```go
/// descriptor, err := loadPluginDescriptor("./plugins/my-plugin/plugin.yaml")
/// if err != nil {
///     log.Fatal(err)
/// }
/// ```
func loadPluginDescriptor(path string) (PluginDescriptor, error) {
	var descriptor PluginDescriptor

	// Dosyayı oku
	data, err := os.ReadFile(path)
	if err != nil {
		return descriptor, fmt.Errorf("failed to read plugin descriptor: %w", err)
	}

	// YAML parse et
	if err := yaml.Unmarshal(data, &descriptor); err != nil {
		return descriptor, fmt.Errorf("failed to parse plugin descriptor: %w", err)
	}

	// Validate
	if err := descriptor.Validate(); err != nil {
		return descriptor, fmt.Errorf("invalid plugin descriptor: %w", err)
	}

	// Enabled varsayılan değeri true
	if descriptor.Enabled == false && descriptor.Entry == "" {
		descriptor.Enabled = true
	}

	return descriptor, nil
}

/// # DiscoverAndList Fonksiyonu
///
/// Plugin'leri keşfeder ve metadata listesi döndürür.
/// Sadece etkin plugin'leri listeler.
///
/// ## Parametreler
/// - `path`: Plugin klasörü yolu
///
/// ## Dönüş Değeri
/// - `[]Metadata`: Etkin plugin'lerin metadata listesi
/// - `error`: Hata varsa hata, yoksa nil
///
/// ## Kullanım Örneği
/// ```go
/// metadataList, err := plugin.DiscoverAndList("./plugins")
/// if err != nil {
///     log.Fatal(err)
/// }
///
/// for _, m := range metadataList {
///     fmt.Println(m.String())
/// }
/// ```
func DiscoverAndList(path string) ([]Metadata, error) {
	descriptors, err := Discover(path)
	if err != nil {
		return nil, err
	}

	var metadataList []Metadata
	for _, desc := range descriptors {
		if desc.Enabled {
			metadataList = append(metadataList, desc.Metadata)
		}
	}

	return metadataList, nil
}
