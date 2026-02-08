package handler

import (
	"github.com/ferdiunal/panel.go/pkg/openapi"
	"github.com/ferdiunal/panel.go/pkg/resource"
	"github.com/gofiber/fiber/v2"
)

// OpenAPIHandler, OpenAPI spesifikasyonu ve dokümantasyon endpoint'lerini yönetir.
//
// ## Özellikler
// - OpenAPI 3.0.3 spec JSON endpoint'i
// - Swagger UI arayüzü
// - ReDoc arayüzü
// - RapiDoc arayüzü
//
// ## Kullanım Örneği
//
//	handler := NewOpenAPIHandler(resources, config)
//	app.Get("/api/openapi.json", handler.GetSpec)
//	app.Get("/api/docs", handler.SwaggerUI)
type OpenAPIHandler struct {
	resources map[string]resource.Resource
	config    openapi.SpecGeneratorConfig
	generator *openapi.SpecGenerator
}

// NewOpenAPIHandler, yeni bir OpenAPIHandler oluşturur.
//
// ## Parametreler
// - resources: Resource map'i (slug -> Resource)
// - config: OpenAPI yapılandırması
//
// ## Döndürür
// - OpenAPIHandler pointer'ı
func NewOpenAPIHandler(resources map[string]resource.Resource, config openapi.SpecGeneratorConfig) *OpenAPIHandler {
	generator := openapi.NewSpecGenerator(resources, config)
	return &OpenAPIHandler{
		resources: resources,
		config:    config,
		generator: generator,
	}
}

// GetSpec, OpenAPI spesifikasyonunu JSON formatında döner.
//
// ## Endpoint
// GET /api/openapi.json
//
// ## Response
// - 200: OpenAPI spec JSON
// - 500: Internal server error
//
// ## Kullanım Örneği
//
//	app.Get("/api/openapi.json", handler.GetSpec)
func (h *OpenAPIHandler) GetSpec(c *fiber.Ctx) error {
	// Spec'i generator'dan al
	spec, err := h.generator.GetSpec()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to generate OpenAPI spec",
		})
	}

	return c.JSON(spec)
}

// SwaggerUI, Swagger UI arayüzünü gösterir.
//
// ## Endpoint
// GET /api/docs
//
// ## Response
// - 200: Swagger UI HTML
//
// ## Kullanım Örneği
//
//	app.Get("/api/docs", handler.SwaggerUI)
func (h *OpenAPIHandler) SwaggerUI(c *fiber.Ctx) error {
	// Spec URL'ini oluştur
	specURL := "/api/openapi.json"

	html := openapi.SwaggerUIHTML(specURL, h.config.Title)
	c.Type("html")
	return c.SendString(html)
}

// ReDocUI, ReDoc arayüzünü gösterir.
//
// ## Endpoint
// GET /api/docs/redoc
//
// ## Response
// - 200: ReDoc HTML
//
// ## Kullanım Örneği
//
//	app.Get("/api/docs/redoc", handler.ReDocUI)
func (h *OpenAPIHandler) ReDocUI(c *fiber.Ctx) error {
	// Spec URL'ini oluştur
	specURL := "/api/openapi.json"

	html := openapi.ReDocHTML(specURL, h.config.Title)
	c.Type("html")
	return c.SendString(html)
}

// RapidocUI, RapiDoc arayüzünü gösterir.
//
// ## Endpoint
// GET /api/docs/rapidoc
//
// ## Response
// - 200: RapiDoc HTML
//
// ## Kullanım Örneği
//
//	app.Get("/api/docs/rapidoc", handler.RapidocUI)
func (h *OpenAPIHandler) RapidocUI(c *fiber.Ctx) error {
	// Spec URL'ini oluştur
	specURL := "/api/openapi.json"

	html := openapi.RapidocHTML(specURL, h.config.Title)
	c.Type("html")
	return c.SendString(html)
}

// RefreshSpec, OpenAPI spec cache'ini temizler ve yeniden oluşturulmasını sağlar.
//
// ## Kullanım Örneği
//
//	handler.RefreshSpec()
func (h *OpenAPIHandler) RefreshSpec() {
	h.generator.InvalidateCache()
}
