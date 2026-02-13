// Package openapi, Swagger UI için HTML handler sağlar.
//
// Bu paket, OpenAPI spesifikasyonunu görselleştirmek için Swagger UI arayüzünü sunar.
package openapi

import (
	"fmt"
)

// SwaggerUIHTML, Swagger UI için HTML içeriği oluşturur.
//
// ## Parametreler
// - specURL: OpenAPI spec JSON endpoint URL'i
// - title: Swagger UI başlığı
//
// ## Döndürür
// - HTML string
//
// ## Kullanım Örneği
//
//	html := SwaggerUIHTML("/api/openapi.json", "Panel.go API")
//	c.Type("html")
//	return c.SendString(html)
func SwaggerUIHTML(specURL, title string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s - Swagger UI</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui.css">
    <style>
        html {
            box-sizing: border-box;
            overflow: -moz-scrollbars-vertical;
            overflow-y: scroll;
        }
        *, *:before, *:after {
            box-sizing: inherit;
        }
        body {
            margin: 0;
            padding: 0;
        }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@5.11.0/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            window.ui = SwaggerUIBundle({
                url: "%s",
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout",
                persistAuthorization: true,
                tryItOutEnabled: true,
                filter: true,
                syntaxHighlight: {
                    activate: true,
                    theme: "monokai"
                }
            });
        };
    </script>
</body>
</html>`, title, specURL)
}

// ReDocHTML, ReDoc UI için HTML içeriği oluşturur.
//
// ReDoc, Swagger UI'ya alternatif, daha modern bir dokümantasyon arayüzüdür.
//
// ## Parametreler
// - specURL: OpenAPI spec JSON endpoint URL'i
// - title: ReDoc başlığı
//
// ## Döndürür
// - HTML string
//
// ## Kullanım Örneği
//
//	html := ReDocHTML("/api/openapi.json", "Panel.go API")
//	c.Type("html")
//	return c.SendString(html)
func ReDocHTML(specURL, title string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s - ReDoc</title>
    <style>
        body {
            margin: 0;
            padding: 0;
        }
    </style>
</head>
<body>
    <redoc spec-url='%s'></redoc>
    <script src="https://cdn.redoc.ly/redoc/latest/bundles/redoc.standalone.js"></script>
</body>
</html>`, title, specURL)
}

// RapidocHTML, RapiDoc UI için HTML içeriği oluşturur.
//
// RapiDoc, modern ve özelleştirilebilir bir OpenAPI dokümantasyon arayüzüdür.
//
// ## Parametreler
// - specURL: OpenAPI spec JSON endpoint URL'i
// - title: RapiDoc başlığı
//
// ## Döndürür
// - HTML string
//
// ## Kullanım Örneği
//
//	html := RapidocHTML("/api/openapi.json", "Panel.go API")
//	c.Type("html")
//	return c.SendString(html)
func RapidocHTML(specURL, title string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s - RapiDoc</title>
    <script type="module" src="https://unpkg.com/rapidoc/dist/rapidoc-min.js"></script>
</head>
<body>
    <rapi-doc
        spec-url="%s"
        theme="dark"
        bg-color="#1a1a1a"
        text-color="#f0f0f0"
        primary-color="#3b82f6"
        render-style="read"
        show-header="true"
        show-info="true"
        allow-authentication="true"
        allow-server-selection="true"
        allow-api-list-style-selection="true"
        schema-style="tree"
        schema-expand-level="2"
        default-schema-tab="model"
    ></rapi-doc>
</body>
</html>`, title, specURL)
}
