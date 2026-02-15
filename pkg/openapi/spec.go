// Package openapi, OpenAPI spesifikasyonu oluşturma ve yönetme işlevlerini sağlar.
//
// Bu paket, Panel.go için OpenAPI 3.0.3 spesifikasyonu oluşturur:
// - Dinamik resource-based endpoint'ler
// - Statik endpoint'ler (auth, init, navigation)
// - Otomatik schema generation
// - Cache mekanizması
package openapi

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/ferdiunal/panel.go/pkg/resource"
	"golang.org/x/sync/singleflight"
)

// SpecGenerator, OpenAPI spesifikasyonu oluşturan ve yöneten yapıdır.
//
// ## Özellikler
//   - Runtime spec generation: Kayıtlı resource'lardan otomatik spec oluşturma
//   - Cache mekanizması: Performans için spec cache'leme
//   - TTL desteği: Cache süresi ayarlanabilir
//   - Thread-safe: Concurrent erişim için mutex koruması
//
// ## Kullanım Örneği
//
//	generator := NewSpecGenerator(panel, SpecGeneratorConfig{
//	    Title:       "Panel.go API",
//	    Version:     "1.0.0",
//	    Description: "RESTful API for Panel.go admin panel",
//	    ServerURL:   "http://localhost:8080",
//	    CacheTTL:    5 * time.Minute,
//	})
//
//	spec, err := generator.GetSpec()
//	if err != nil {
//	    log.Fatal(err)
//	}
type SpecGenerator struct {
	config    SpecGeneratorConfig
	resources map[string]resource.Resource
	mapper    *FieldTypeMapper
	cache     *specCache
	mu        sync.RWMutex
	builds    singleflight.Group
}

// SpecGeneratorConfig, SpecGenerator yapılandırmasını içerir.
//
// ## Alanlar
//   - Title: API başlığı (örn: "Panel.go API")
//   - Version: API versiyonu (örn: "1.0.0")
//   - Description: API açıklaması
//   - ServerURL: API sunucu URL'i (örn: "http://localhost:8080")
//   - ServerDescription: Sunucu açıklaması (örn: "Development server")
//   - CacheTTL: Cache süresi (0 = cache yok, development için önerilir)
//   - CustomMappings: Özel field type mapping'leri
type SpecGeneratorConfig struct {
	Title               string                 // API başlığı
	Version             string                 // API versiyonu
	Description         string                 // API açıklaması
	ServerURL           string                 // API sunucu URL'i
	ServerDescription   string                 // Sunucu açıklaması
	APIKeyHeader        string                 // API key header adı (Swagger authorize için)
	CacheTTL            time.Duration          // Cache süresi (0 = cache yok)
	CustomMappings      *CustomMappingRegistry // Özel field type mapping'leri
	EnableParallelBuild bool                   // Dinamik path/schema üretimi bounded parallel çalıştır
	ParallelWorkers     int                    // Parallel build worker limiti (0 = auto-tune)
}

// specCache, OpenAPI spec cache'ini yönetir.
//
// ## Alanlar
//   - spec: Cache'lenmiş OpenAPI spec
//   - timestamp: Cache oluşturulma zamanı
//   - ttl: Cache süresi
type specCache struct {
	spec      *OpenAPISpec
	timestamp time.Time
	ttl       time.Duration
	mu        sync.RWMutex
}

// NewSpecGenerator, yeni bir SpecGenerator oluşturur.
//
// ## Parametreler
//   - resources: Kayıtlı resource'ların haritası
//   - config: Generator yapılandırması
//
// ## Dönüş Değeri
//   - *SpecGenerator: Yapılandırılmış generator
//
// ## Kullanım Örneği
//
//	generator := NewSpecGenerator(panel.resources, SpecGeneratorConfig{
//	    Title:       "Panel.go API",
//	    Version:     "1.0.0",
//	    Description: "RESTful API for Panel.go admin panel",
//	    ServerURL:   "http://localhost:8080",
//	    CacheTTL:    5 * time.Minute,
//	})
func NewSpecGenerator(resources map[string]resource.Resource, config SpecGeneratorConfig) *SpecGenerator {
	// Varsayılan değerleri ayarla
	if config.Title == "" {
		config.Title = "Panel.go API"
	}
	if config.Version == "" {
		config.Version = "1.0.0"
	}
	if config.Description == "" {
		config.Description = "RESTful API for Panel.go admin panel"
	}
	if config.ServerURL == "" {
		config.ServerURL = "http://localhost:8080"
	}
	if config.ServerDescription == "" {
		config.ServerDescription = "API Server"
	}
	if config.APIKeyHeader == "" {
		config.APIKeyHeader = "X-API-Key"
	}

	mapper := NewFieldTypeMapper()
	if config.CustomMappings != nil {
		mapper.registry = config.CustomMappings
	}

	return &SpecGenerator{
		config:    config,
		resources: resources,
		mapper:    mapper,
		cache: &specCache{
			ttl: config.CacheTTL,
		},
	}
}

// GetSpec, OpenAPI spesifikasyonunu döndürür.
//
// ## Davranış
//   - Cache varsa ve geçerliyse cache'den döner
//   - Cache yoksa veya süresi dolmuşsa yeniden oluşturur
//   - Thread-safe: Concurrent erişim için mutex koruması
//
// ## Dönüş Değeri
//   - *OpenAPISpec: OpenAPI spesifikasyonu
//   - error: Hata varsa hata, aksi takdirde nil
//
// ## Kullanım Örneği
//
//	spec, err := generator.GetSpec()
//	if err != nil {
//	    return err
//	}
//	// spec'i JSON olarak serialize et
func (g *SpecGenerator) GetSpec() (*OpenAPISpec, error) {
	// Cache kontrolü
	if cached := g.cache.get(); cached != nil {
		return cached, nil
	}

	value, err, _ := g.builds.Do("openapi-spec", func() (interface{}, error) {
		// Double-check: başka bir goroutine cache'i doldurmuş olabilir.
		if cached := g.cache.get(); cached != nil {
			return cached, nil
		}

		// Yeni spec oluştur
		spec, buildErr := g.generateSpec()
		if buildErr != nil {
			return nil, buildErr
		}

		// Cache'e kaydet
		g.cache.set(spec)

		return spec, nil
	})
	if err != nil {
		return nil, err
	}

	spec, _ := value.(*OpenAPISpec)
	if spec == nil {
		return nil, fmt.Errorf("failed to generate OpenAPI spec")
	}
	return spec, nil
}

// InvalidateCache, cache'i geçersiz kılar.
//
// ## Kullanım Örneği
//
//	// Resource değişikliğinden sonra cache'i temizle
//	generator.InvalidateCache()
func (g *SpecGenerator) InvalidateCache() {
	g.cache.invalidate()
}

// generateSpec, OpenAPI spesifikasyonunu oluşturur.
//
// ## Davranış
//  1. Temel spec yapısını oluşturur
//  2. Statik endpoint'leri ekler (auth, init, navigation)
//  3. Dinamik resource endpoint'lerini ekler
//  4. Component'leri (schemas, security schemes) ekler
//  5. Tag'leri ekler
//
// ## Dönüş Değeri
//   - *OpenAPISpec: Oluşturulan OpenAPI spesifikasyonu
//   - error: Hata varsa hata, aksi takdirde nil
func (g *SpecGenerator) generateSpec() (*OpenAPISpec, error) {
	spec := &OpenAPISpec{
		OpenAPI: "3.0.3",
		Info: Info{
			Title:       g.config.Title,
			Version:     g.config.Version,
			Description: g.config.Description,
		},
		Servers: []Server{
			{
				URL:         g.config.ServerURL,
				Description: g.config.ServerDescription,
			},
		},
		Paths: make(map[string]PathItem),
		Components: Components{
			Schemas:         make(map[string]Schema),
			SecuritySchemes: make(map[string]SecurityScheme),
		},
		Tags: []Tag{},
	}

	// Security scheme ekle (cookie-based authentication)
	spec.Components.SecuritySchemes["cookieAuth"] = SecurityScheme{
		Type:        "apiKey",
		In:          "cookie",
		Name:        "session_token",
		Description: "Session cookie authentication",
	}
	spec.Components.SecuritySchemes["apiKeyAuth"] = SecurityScheme{
		Type:        "apiKey",
		In:          "header",
		Name:        g.config.APIKeyHeader,
		Description: "API key authentication",
	}

	// Global security requirement ekle
	spec.Security = []SecurityRequirement{
		{"cookieAuth": []string{}},
		{"apiKeyAuth": []string{}},
	}

	// Statik endpoint'leri ekle
	if err := g.addStaticEndpoints(spec); err != nil {
		return nil, fmt.Errorf("failed to add static endpoints: %w", err)
	}

	// Dinamik resource endpoint'lerini ekle
	if err := g.addDynamicEndpoints(spec); err != nil {
		return nil, fmt.Errorf("failed to add dynamic endpoints: %w", err)
	}

	// Common schemas ekle
	g.addCommonSchemas(spec)

	return spec, nil
}

// addStaticEndpoints, statik endpoint'leri spec'e ekler.
//
// ## Eklenen Endpoint'ler
//   - POST /api/auth/sign-in/email: Email ile giriş
//   - POST /api/auth/sign-up/email: Email ile kayıt
//   - POST /api/auth/sign-out: Çıkış
//   - POST /api/auth/forgot-password: Şifremi unuttum
//   - GET /api/auth/session: Oturum bilgisi
//   - GET /api/init: Uygulama başlatma bilgileri
//   - GET /api/navigation: Navigasyon menüsü
func (g *SpecGenerator) addStaticEndpoints(spec *OpenAPISpec) error {
	staticGen := NewStaticSpecGenerator()
	staticPaths := staticGen.GenerateStaticPaths()

	// Statik path'leri spec'e ekle
	for path, pathItem := range staticPaths {
		spec.Paths[path] = pathItem
	}

	// Statik tag'leri ekle
	spec.Tags = append(spec.Tags, Tag{
		Name:        "auth",
		Description: "Authentication endpoints",
	})
	spec.Tags = append(spec.Tags, Tag{
		Name:        "system",
		Description: "System endpoints",
	})

	return nil
}

// addDynamicEndpoints, dinamik resource endpoint'lerini spec'e ekler.
//
// ## Davranış
//   - Her resource için CRUD endpoint'leri oluşturur
//   - Resource schema'larını oluşturur
//   - Tag'leri ekler
func (g *SpecGenerator) addDynamicEndpoints(spec *OpenAPISpec) error {
	dynamicGen := NewDynamicSpecGenerator()

	// Tüm resource'lar için path'leri oluştur
	paths := dynamicGen.GenerateResourcePaths(g.resources)
	if g.config.EnableParallelBuild {
		paths = dynamicGen.GenerateResourcePathsParallel(g.resources, g.config.ParallelWorkers)
	}

	// Path'leri spec'e ekle
	for path, pathItem := range paths {
		spec.Paths[path] = pathItem
	}

	// Tüm resource'lar için schema'ları oluştur
	schemas := dynamicGen.GenerateResourceSchemas(g.resources)
	if g.config.EnableParallelBuild {
		schemas = dynamicGen.GenerateResourceSchemasParallel(g.resources, g.config.ParallelWorkers)
	}

	// Schema'ları ve tag'leri ekle
	for slug, res := range g.resources {
		// Schema'yı components'e ekle
		schemaName := toPascalCase(slug)
		if schema, ok := schemas[schemaName]; ok {
			spec.Components.Schemas[schemaName] = *schema
		}

		// Tag ekle
		spec.Tags = append(spec.Tags, Tag{
			Name:        slug,
			Description: fmt.Sprintf("%s resource endpoints", res.Title()),
		})
	}

	return nil
}

// addCommonSchemas, yaygın kullanılan schema'ları ekler.
//
// ## Eklenen Schema'lar
//   - PaginationMeta: Sayfalama metadata'sı
//   - ErrorResponse: Hata yanıtı
//   - SuccessResponse: Başarı yanıtı
func (g *SpecGenerator) addCommonSchemas(spec *OpenAPISpec) {
	// PaginationMeta schema
	spec.Components.Schemas["PaginationMeta"] = Schema{
		Type: "object",
		Properties: map[string]Schema{
			"total": {
				Type:        "integer",
				Description: "Toplam kayıt sayısı",
				Example:     100,
			},
			"per_page": {
				Type:        "integer",
				Description: "Sayfa başına kayıt sayısı",
				Example:     15,
			},
			"current_page": {
				Type:        "integer",
				Description: "Mevcut sayfa numarası",
				Example:     1,
			},
			"last_page": {
				Type:        "integer",
				Description: "Son sayfa numarası",
				Example:     7,
			},
			"from": {
				Type:        "integer",
				Description: "İlk kayıt numarası",
				Example:     1,
			},
			"to": {
				Type:        "integer",
				Description: "Son kayıt numarası",
				Example:     15,
			},
		},
	}

	// ErrorResponse schema
	spec.Components.Schemas["ErrorResponse"] = Schema{
		Type: "object",
		Properties: map[string]Schema{
			"error": {
				Type:        "string",
				Description: "Hata mesajı",
				Example:     "Resource not found",
			},
			"code": {
				Type:        "string",
				Description: "Hata kodu",
				Example:     "RESOURCE_NOT_FOUND",
			},
		},
		Required: []string{"error"},
	}

	// SuccessResponse schema
	spec.Components.Schemas["SuccessResponse"] = Schema{
		Type: "object",
		Properties: map[string]Schema{
			"message": {
				Type:        "string",
				Description: "Başarı mesajı",
				Example:     "Operation completed successfully",
			},
		},
	}
}

// get, cache'den spec'i döndürür.
//
// ## Davranış
//   - Cache varsa ve geçerliyse spec'i döner
//   - Cache yoksa veya süresi dolmuşsa nil döner
//
// ## Dönüş Değeri
//   - *OpenAPISpec: Cache'lenmiş spec veya nil
func (c *specCache) get() *OpenAPISpec {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Cache yoksa
	if c.spec == nil {
		return nil
	}

	// TTL 0 ise cache kullanma (development mode)
	if c.ttl == 0 {
		return nil
	}

	// Cache süresi dolmuşsa
	if time.Since(c.timestamp) > c.ttl {
		return nil
	}

	if cloned := cloneSpec(c.spec); cloned != nil {
		return cloned
	}
	return c.spec
}

// set, spec'i cache'e kaydeder.
//
// ## Parametreler
//   - spec: Cache'lenecek OpenAPI spec
func (c *specCache) set(spec *OpenAPISpec) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if cloned := cloneSpec(spec); cloned != nil {
		c.spec = cloned
	} else {
		c.spec = spec
	}
	c.timestamp = time.Now()
}

// invalidate, cache'i geçersiz kılar.
func (c *specCache) invalidate() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.spec = nil
	c.timestamp = time.Time{}
}

// toPascalCase, string'i PascalCase'e çevirir.
//
// ## Kullanım Örneği
//
//	toPascalCase("users")        // "Users"
//	toPascalCase("blog-posts")   // "BlogPosts"
//	toPascalCase("user_profiles") // "UserProfiles"
func toPascalCase(s string) string {
	if s == "" {
		return ""
	}

	// İlk harfi büyük yap
	result := ""
	capitalize := true

	for _, ch := range s {
		if ch == '-' || ch == '_' || ch == ' ' {
			capitalize = true
			continue
		}

		if capitalize {
			result += string(ch - 32) // Büyük harfe çevir
			capitalize = false
		} else {
			result += string(ch)
		}
	}

	return result
}

func cloneSpec(spec *OpenAPISpec) *OpenAPISpec {
	if spec == nil {
		return nil
	}

	raw, err := json.Marshal(spec)
	if err != nil {
		return nil
	}

	var cloned OpenAPISpec
	if err := json.Unmarshal(raw, &cloned); err != nil {
		return nil
	}

	return &cloned
}
