// Package openapi, OpenAPI 3.0.3 spesifikasyonu için tip tanımlamalarını sağlar.
//
// Bu paket, OpenAPI spesifikasyonunun temel yapı taşlarını içerir:
// - OpenAPI spec yapısı
// - Path ve operation tanımlamaları
// - Schema ve component tanımlamaları
// - Parameter ve response tanımlamaları
//
// ## Kullanım Örneği
//
//	spec := &OpenAPISpec{
//	    OpenAPI: "3.0.3",
//	    Info: Info{
//	        Title:   "Panel.go API",
//	        Version: "1.0.0",
//	    },
//	}
package openapi

// OpenAPISpec, OpenAPI 3.0.3 spesifikasyonunun ana yapısıdır.
//
// Bu yapı, API'nin tüm bilgilerini içerir: endpoint'ler, şemalar, güvenlik ayarları vb.
//
// ## Alanlar
//   - OpenAPI: OpenAPI spesifikasyon versiyonu (örn: "3.0.3")
//   - Info: API hakkında genel bilgiler (başlık, versiyon, açıklama)
//   - Servers: API sunucularının listesi
//   - Paths: API endpoint'lerinin tanımları
//   - Components: Yeniden kullanılabilir bileşenler (şemalar, parametreler, vb.)
//   - Security: Global güvenlik gereksinimleri
//   - Tags: Endpoint'leri gruplamak için etiketler
//
// ## Kullanım Örneği
//
//	spec := &OpenAPISpec{
//	    OpenAPI: "3.0.3",
//	    Info: Info{
//	        Title:       "Panel.go API",
//	        Description: "RESTful API for Panel.go admin panel",
//	        Version:     "1.0.0",
//	    },
//	    Servers: []Server{
//	        {URL: "http://localhost:8080", Description: "Development server"},
//	    },
//	    Paths: make(map[string]PathItem),
//	    Components: Components{
//	        Schemas:         make(map[string]Schema),
//	        SecuritySchemes: make(map[string]SecurityScheme),
//	    },
//	}
type OpenAPISpec struct {
	OpenAPI    string                `json:"openapi"`              // OpenAPI spesifikasyon versiyonu (3.0.3)
	Info       Info                  `json:"info"`                 // API hakkında genel bilgiler
	Servers    []Server              `json:"servers,omitempty"`    // API sunucularının listesi
	Paths      map[string]PathItem   `json:"paths"`                // API endpoint'lerinin tanımları
	Components Components            `json:"components,omitempty"` // Yeniden kullanılabilir bileşenler
	Security   []SecurityRequirement `json:"security,omitempty"`   // Global güvenlik gereksinimleri
	Tags       []Tag                 `json:"tags,omitempty"`       // Endpoint grupları için etiketler
}

// Info, API hakkında genel bilgileri içerir.
//
// ## Alanlar
//   - Title: API'nin başlığı (örn: "Panel.go API")
//   - Description: API'nin açıklaması
//   - Version: API versiyonu (örn: "1.0.0")
//   - Contact: İletişim bilgileri (opsiyonel)
//   - License: Lisans bilgileri (opsiyonel)
type Info struct {
	Title       string   `json:"title"`                 // API başlığı
	Description string   `json:"description,omitempty"` // API açıklaması
	Version     string   `json:"version"`               // API versiyonu
	Contact     *Contact `json:"contact,omitempty"`     // İletişim bilgileri
	License     *License `json:"license,omitempty"`     // Lisans bilgileri
}

// Contact, API ile ilgili iletişim bilgilerini içerir.
type Contact struct {
	Name  string `json:"name,omitempty"`  // İletişim kişisi adı
	URL   string `json:"url,omitempty"`   // İletişim URL'i
	Email string `json:"email,omitempty"` // İletişim e-posta adresi
}

// License, API'nin lisans bilgilerini içerir.
type License struct {
	Name string `json:"name"`           // Lisans adı (örn: "MIT")
	URL  string `json:"url,omitempty"`  // Lisans URL'i
}

// Server, API sunucusunun bilgilerini içerir.
//
// ## Alanlar
//   - URL: Sunucu URL'i (örn: "http://localhost:8080")
//   - Description: Sunucu açıklaması (örn: "Development server")
type Server struct {
	URL         string `json:"url"`                   // Sunucu URL'i
	Description string `json:"description,omitempty"` // Sunucu açıklaması
}

// PathItem, bir API endpoint'inin tüm HTTP metodlarını içerir.
//
// ## Alanlar
//   - Get: GET metodu için operation
//   - Post: POST metodu için operation
//   - Put: PUT metodu için operation
//   - Delete: DELETE metodu için operation
//   - Patch: PATCH metodu için operation
//   - Parameters: Tüm metodlar için ortak parametreler
//
// ## Kullanım Örneği
//
//	pathItem := PathItem{
//	    Get: &Operation{
//	        Summary:     "List users",
//	        Description: "Returns a list of users",
//	        Tags:        []string{"users"},
//	        Responses: map[string]Response{
//	            "200": {Description: "Successful response"},
//	        },
//	    },
//	}
type PathItem struct {
	Get        *Operation  `json:"get,omitempty"`        // GET metodu
	Post       *Operation  `json:"post,omitempty"`       // POST metodu
	Put        *Operation  `json:"put,omitempty"`        // PUT metodu
	Delete     *Operation  `json:"delete,omitempty"`     // DELETE metodu
	Patch      *Operation  `json:"patch,omitempty"`      // PATCH metodu
	Parameters []Parameter `json:"parameters,omitempty"` // Ortak parametreler
}

// Operation, bir HTTP metodunun detaylarını içerir.
//
// ## Alanlar
//   - Summary: Kısa özet (örn: "List users")
//   - Description: Detaylı açıklama
//   - OperationID: Benzersiz operation ID (örn: "listUsers")
//   - Tags: Endpoint'i gruplamak için etiketler
//   - Parameters: İstek parametreleri (query, path, header)
//   - RequestBody: İstek gövdesi (POST, PUT için)
//   - Responses: Olası yanıtlar (200, 404, 500 vb.)
//   - Security: Bu operation için güvenlik gereksinimleri
//
// ## Kullanım Örneği
//
//	operation := &Operation{
//	    Summary:     "Create user",
//	    Description: "Creates a new user",
//	    OperationID: "createUser",
//	    Tags:        []string{"users"},
//	    RequestBody: &RequestBody{
//	        Required: true,
//	        Content: map[string]MediaType{
//	            "application/json": {
//	                Schema: &Schema{Ref: "#/components/schemas/User"},
//	            },
//	        },
//	    },
//	    Responses: map[string]Response{
//	        "201": {Description: "User created"},
//	        "400": {Description: "Invalid input"},
//	    },
//	}
type Operation struct {
	Summary     string              `json:"summary,omitempty"`     // Kısa özet
	Description string              `json:"description,omitempty"` // Detaylı açıklama
	OperationID string              `json:"operationId,omitempty"` // Benzersiz operation ID
	Tags        []string            `json:"tags,omitempty"`        // Etiketler
	Parameters  []Parameter         `json:"parameters,omitempty"`  // Parametreler
	RequestBody *RequestBody        `json:"requestBody,omitempty"` // İstek gövdesi
	Responses   map[string]Response `json:"responses"`             // Yanıtlar
	Security    []SecurityRequirement `json:"security,omitempty"`  // Güvenlik gereksinimleri
}

// Parameter, bir API parametresini tanımlar.
//
// ## Alanlar
//   - Name: Parametre adı (örn: "id", "page", "search")
//   - In: Parametre konumu ("query", "path", "header", "cookie")
//   - Description: Parametre açıklaması
//   - Required: Zorunlu mu?
//   - Schema: Parametre şeması (tip, format, enum vb.)
//   - Example: Örnek değer
//
// ## Kullanım Örneği
//
//	// Path parameter
//	idParam := Parameter{
//	    Name:        "id",
//	    In:          "path",
//	    Description: "User ID",
//	    Required:    true,
//	    Schema:      &Schema{Type: "integer", Format: "int64"},
//	}
//
//	// Query parameter
//	pageParam := Parameter{
//	    Name:        "page",
//	    In:          "query",
//	    Description: "Page number",
//	    Required:    false,
//	    Schema:      &Schema{Type: "integer", Default: 1},
//	}
type Parameter struct {
	Name        string      `json:"name"`                  // Parametre adı
	In          string      `json:"in"`                    // Konum: query, path, header, cookie
	Description string      `json:"description,omitempty"` // Açıklama
	Required    bool        `json:"required,omitempty"`    // Zorunlu mu?
	Schema      *Schema     `json:"schema,omitempty"`      // Parametre şeması
	Example     interface{} `json:"example,omitempty"`     // Örnek değer
}

// RequestBody, bir HTTP isteğinin gövdesini tanımlar.
//
// ## Alanlar
//   - Description: İstek gövdesi açıklaması
//   - Required: Zorunlu mu?
//   - Content: İçerik tipleri ve şemaları (application/json, multipart/form-data vb.)
//
// ## Kullanım Örneği
//
//	requestBody := &RequestBody{
//	    Description: "User object to create",
//	    Required:    true,
//	    Content: map[string]MediaType{
//	        "application/json": {
//	            Schema: &Schema{
//	                Type: "object",
//	                Properties: map[string]Schema{
//	                    "name":  {Type: "string"},
//	                    "email": {Type: "string", Format: "email"},
//	                },
//	                Required: []string{"name", "email"},
//	            },
//	        },
//	    },
//	}
type RequestBody struct {
	Description string               `json:"description,omitempty"` // Açıklama
	Required    bool                 `json:"required,omitempty"`    // Zorunlu mu?
	Content     map[string]MediaType `json:"content"`               // İçerik tipleri
}

// MediaType, bir içerik tipinin şemasını tanımlar.
//
// ## Alanlar
//   - Schema: İçerik şeması
//   - Example: Örnek değer
//   - Examples: Birden fazla örnek
type MediaType struct {
	Schema   *Schema                `json:"schema,omitempty"`   // İçerik şeması
	Example  interface{}            `json:"example,omitempty"`  // Tek örnek
	Examples map[string]interface{} `json:"examples,omitempty"` // Birden fazla örnek
}

// Response, bir HTTP yanıtını tanımlar.
//
// ## Alanlar
//   - Description: Yanıt açıklaması (zorunlu)
//   - Content: Yanıt içeriği (application/json vb.)
//   - Headers: Yanıt başlıkları
//
// ## Kullanım Örneği
//
//	response := Response{
//	    Description: "Successful response",
//	    Content: map[string]MediaType{
//	        "application/json": {
//	            Schema: &Schema{
//	                Type: "object",
//	                Properties: map[string]Schema{
//	                    "data": {Type: "array", Items: &Schema{Ref: "#/components/schemas/User"}},
//	                    "meta": {Ref: "#/components/schemas/PaginationMeta"},
//	                },
//	            },
//	        },
//	    },
//	}
type Response struct {
	Description string               `json:"description"`           // Yanıt açıklaması (zorunlu)
	Content     map[string]MediaType `json:"content,omitempty"`     // Yanıt içeriği
	Headers     map[string]Header    `json:"headers,omitempty"`     // Yanıt başlıkları
}

// Header, bir HTTP başlığını tanımlar.
type Header struct {
	Description string  `json:"description,omitempty"` // Başlık açıklaması
	Schema      *Schema `json:"schema,omitempty"`      // Başlık şeması
}

// Schema, bir veri yapısını tanımlar (JSON Schema benzeri).
//
// ## Alanlar
//   - Ref: Başka bir şemaya referans (#/components/schemas/User)
//   - Type: Veri tipi (string, number, integer, boolean, array, object)
//   - Format: Veri formatı (date, date-time, email, uri, binary vb.)
//   - Properties: Object tipi için özellikler
//   - Items: Array tipi için öğe şeması
//   - Required: Zorunlu alanlar
//   - Enum: İzin verilen değerler
//   - Default: Varsayılan değer
//   - Example: Örnek değer
//   - Description: Açıklama
//   - Minimum, Maximum: Sayısal değerler için min/max
//   - MinLength, MaxLength: String değerler için min/max uzunluk
//   - Pattern: Regex pattern
//   - Nullable: Null değer alabilir mi?
//
// ## Kullanım Örneği
//
//	// String schema
//	nameSchema := Schema{
//	    Type:        "string",
//	    Description: "User name",
//	    MinLength:   ptr(1),
//	    MaxLength:   ptr(255),
//	    Example:     "John Doe",
//	}
//
//	// Object schema
//	userSchema := Schema{
//	    Type: "object",
//	    Properties: map[string]Schema{
//	        "id":    {Type: "integer", Format: "int64"},
//	        "name":  {Type: "string"},
//	        "email": {Type: "string", Format: "email"},
//	    },
//	    Required: []string{"name", "email"},
//	}
//
//	// Array schema
//	usersSchema := Schema{
//	    Type:  "array",
//	    Items: &Schema{Ref: "#/components/schemas/User"},
//	}
type Schema struct {
	Ref         string            `json:"$ref,omitempty"`        // Referans (#/components/schemas/User)
	Type        string            `json:"type,omitempty"`        // Veri tipi
	Format      string            `json:"format,omitempty"`      // Veri formatı
	Properties  map[string]Schema `json:"properties,omitempty"`  // Object özellikleri
	Items       *Schema           `json:"items,omitempty"`       // Array öğeleri
	Required    []string          `json:"required,omitempty"`    // Zorunlu alanlar
	Enum        []interface{}     `json:"enum,omitempty"`        // İzin verilen değerler
	Default     interface{}       `json:"default,omitempty"`     // Varsayılan değer
	Example     interface{}       `json:"example,omitempty"`     // Örnek değer
	Description string            `json:"description,omitempty"` // Açıklama
	Minimum     *float64          `json:"minimum,omitempty"`     // Minimum değer
	Maximum     *float64          `json:"maximum,omitempty"`     // Maximum değer
	MinLength   *int              `json:"minLength,omitempty"`   // Minimum uzunluk
	MaxLength   *int              `json:"maxLength,omitempty"`   // Maximum uzunluk
	Pattern     string            `json:"pattern,omitempty"`     // Regex pattern
	Nullable    bool              `json:"nullable,omitempty"`    // Null değer alabilir mi?
	ReadOnly    bool              `json:"readOnly,omitempty"`    // Salt okunur mu?
	WriteOnly   bool              `json:"writeOnly,omitempty"`   // Sadece yazılabilir mi?
}

// Components, yeniden kullanılabilir bileşenleri içerir.
//
// ## Alanlar
//   - Schemas: Şema tanımlamaları (User, Post, Category vb.)
//   - Parameters: Parametre tanımlamaları
//   - RequestBodies: İstek gövdesi tanımlamaları
//   - Responses: Yanıt tanımlamaları
//   - SecuritySchemes: Güvenlik şeması tanımlamaları
//
// ## Kullanım Örneği
//
//	components := Components{
//	    Schemas: map[string]Schema{
//	        "User": {
//	            Type: "object",
//	            Properties: map[string]Schema{
//	                "id":    {Type: "integer"},
//	                "name":  {Type: "string"},
//	                "email": {Type: "string", Format: "email"},
//	            },
//	        },
//	    },
//	    SecuritySchemes: map[string]SecurityScheme{
//	        "cookieAuth": {
//	            Type: "apiKey",
//	            In:   "cookie",
//	            Name: "session_token",
//	        },
//	    },
//	}
type Components struct {
	Schemas         map[string]Schema         `json:"schemas,omitempty"`         // Şema tanımlamaları
	Parameters      map[string]Parameter      `json:"parameters,omitempty"`      // Parametre tanımlamaları
	RequestBodies   map[string]RequestBody    `json:"requestBodies,omitempty"`   // İstek gövdesi tanımlamaları
	Responses       map[string]Response       `json:"responses,omitempty"`       // Yanıt tanımlamaları
	SecuritySchemes map[string]SecurityScheme `json:"securitySchemes,omitempty"` // Güvenlik şeması tanımlamaları
}

// SecurityScheme, bir güvenlik şemasını tanımlar.
//
// ## Alanlar
//   - Type: Güvenlik tipi (apiKey, http, oauth2, openIdConnect)
//   - In: API key konumu (query, header, cookie)
//   - Name: API key adı
//   - Scheme: HTTP authentication şeması (basic, bearer)
//   - BearerFormat: Bearer token formatı (JWT)
//   - Description: Açıklama
//
// ## Kullanım Örneği
//
//	// Cookie-based authentication
//	cookieAuth := SecurityScheme{
//	    Type:        "apiKey",
//	    In:          "cookie",
//	    Name:        "session_token",
//	    Description: "Session cookie authentication",
//	}
//
//	// Bearer token authentication
//	bearerAuth := SecurityScheme{
//	    Type:         "http",
//	    Scheme:       "bearer",
//	    BearerFormat: "JWT",
//	    Description:  "JWT bearer token authentication",
//	}
type SecurityScheme struct {
	Type         string `json:"type"`                   // Güvenlik tipi
	In           string `json:"in,omitempty"`           // API key konumu
	Name         string `json:"name,omitempty"`         // API key adı
	Scheme       string `json:"scheme,omitempty"`       // HTTP authentication şeması
	BearerFormat string `json:"bearerFormat,omitempty"` // Bearer token formatı
	Description  string `json:"description,omitempty"`  // Açıklama
}

// SecurityRequirement, bir güvenlik gereksinimini tanımlar.
//
// ## Kullanım Örneği
//
//	// Cookie authentication required
//	security := SecurityRequirement{
//	    "cookieAuth": []string{},
//	}
type SecurityRequirement map[string][]string

// Tag, endpoint'leri gruplamak için kullanılan bir etikettir.
//
// ## Alanlar
//   - Name: Etiket adı (örn: "users", "posts")
//   - Description: Etiket açıklaması
//
// ## Kullanım Örneği
//
//	tag := Tag{
//	    Name:        "users",
//	    Description: "User management endpoints",
//	}
type Tag struct {
	Name        string `json:"name"`                  // Etiket adı
	Description string `json:"description,omitempty"` // Etiket açıklaması
}

// ptr, bir değerin pointer'ını döndüren yardımcı fonksiyon.
//
// ## Kullanım Örneği
//
//	schema := Schema{
//	    Type:      "string",
//	    MinLength: ptr(1),
//	    MaxLength: ptr(255),
//	}
func ptr[T any](v T) *T {
	return &v
}
