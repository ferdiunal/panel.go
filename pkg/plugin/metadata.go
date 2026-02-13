/// # Plugin Metadata
///
/// Plugin metadata yapıları ve yardımcı fonksiyonlar.
/// Plugin'lerin meta bilgilerini (ad, versiyon, yazar, bağımlılıklar) yönetir.
///
/// ## Kullanım Örneği
/// ```go
/// metadata := plugin.Metadata{
///     Name:        "analytics-plugin",
///     Version:     "1.0.0",
///     Author:      "John Doe",
///     Description: "Analytics and reporting plugin",
///     Homepage:    "https://github.com/user/analytics-plugin",
///     License:     "MIT",
///     Dependencies: []string{"database-plugin@^1.0.0"},
/// }
/// ```

package plugin

import (
	"fmt"
	"strings"
)

/// # Metadata Struct
///
/// Plugin'in meta bilgilerini tutar.
///
/// ## Alanlar
/// - `Name`: Plugin adı (kebab-case önerilir)
/// - `Version`: Semantic versioning (örn: "1.0.0")
/// - `Author`: Plugin yazarı
/// - `Description`: Kısa açıklama
/// - `Homepage`: Plugin web sitesi veya repository URL'i
/// - `License`: Lisans türü (örn: "MIT", "Apache-2.0")
/// - `Dependencies`: Bağımlı plugin'ler (örn: ["plugin-name@^1.0.0"])
/// - `Tags`: Arama için etiketler
///
/// ## Kullanım Örneği
/// ```go
/// metadata := plugin.Metadata{
///     Name:        "my-plugin",
///     Version:     "1.0.0",
///     Author:      "Author Name",
///     Description: "Plugin description",
///     License:     "MIT",
///     Tags:        []string{"analytics", "reporting"},
/// }
/// ```
type Metadata struct {
	Name         string   `json:"name" yaml:"name"`
	Version      string   `json:"version" yaml:"version"`
	Author       string   `json:"author" yaml:"author"`
	Description  string   `json:"description" yaml:"description"`
	Homepage     string   `json:"homepage,omitempty" yaml:"homepage,omitempty"`
	License      string   `json:"license,omitempty" yaml:"license,omitempty"`
	Dependencies []string `json:"dependencies,omitempty" yaml:"dependencies,omitempty"`
	Tags         []string `json:"tags,omitempty" yaml:"tags,omitempty"`
}

/// # Validate Metodu
///
/// Metadata'nın geçerli olup olmadığını kontrol eder.
///
/// ## Dönüş Değeri
/// - `error`: Geçersizse hata, geçerliyse nil
///
/// ## Kontroller
/// - Name boş olmamalı
/// - Version boş olmamalı
/// - Author boş olmamalı
///
/// ## Kullanım Örneği
/// ```go
/// if err := metadata.Validate(); err != nil {
///     log.Fatal(err)
/// }
/// ```
func (m *Metadata) Validate() error {
	if strings.TrimSpace(m.Name) == "" {
		return fmt.Errorf("plugin metadata: name is required")
	}
	if strings.TrimSpace(m.Version) == "" {
		return fmt.Errorf("plugin metadata: version is required")
	}
	if strings.TrimSpace(m.Author) == "" {
		return fmt.Errorf("plugin metadata: author is required")
	}
	return nil
}

/// # String Metodu
///
/// Metadata'yı okunabilir string formatına çevirir.
///
/// ## Dönüş Değeri
/// - `string`: "name@version by author" formatında string
///
/// ## Kullanım Örneği
/// ```go
/// fmt.Println(metadata.String())
/// // Output: my-plugin@1.0.0 by Author Name
/// ```
func (m *Metadata) String() string {
	return fmt.Sprintf("%s@%s by %s", m.Name, m.Version, m.Author)
}

/// # GetMetadata Fonksiyonu
///
/// Plugin'den metadata bilgilerini çıkarır.
///
/// ## Parametreler
/// - `p`: Plugin instance
///
/// ## Dönüş Değeri
/// - `Metadata`: Plugin'in metadata'sı
///
/// ## Kullanım Örneği
/// ```go
/// metadata := plugin.GetMetadata(myPlugin)
/// fmt.Println(metadata.String())
/// ```
func GetMetadata(p Plugin) Metadata {
	return Metadata{
		Name:        p.Name(),
		Version:     p.Version(),
		Author:      p.Author(),
		Description: p.Description(),
	}
}

/// # ListMetadata Fonksiyonu
///
/// Tüm kayıtlı plugin'lerin metadata'larını döndürür.
///
/// ## Dönüş Değeri
/// - `[]Metadata`: Tüm plugin'lerin metadata listesi
///
/// ## Kullanım Örneği
/// ```go
/// metadataList := plugin.ListMetadata()
/// for _, m := range metadataList {
///     fmt.Println(m.String())
/// }
/// ```
func ListMetadata() []Metadata {
	plugins := All()
	result := make([]Metadata, len(plugins))
	for i, p := range plugins {
		result[i] = GetMetadata(p)
	}
	return result
}
