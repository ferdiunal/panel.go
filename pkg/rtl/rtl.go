// Package rtl provides RTL (Right-to-Left) support utilities for Panel.go
// Arapça, İbranice, Farsça gibi sağdan sola yazılan diller için destek sağlar
package rtl

import (
	"github.com/gofiber/fiber/v2"
	"golang.org/x/text/language"
)

// RTLLanguages, sağdan sola yazılan dillerin listesidir
// Bu diller için otomatik olarak dir="rtl" attribute'u eklenir
var RTLLanguages = map[language.Tag]bool{
	language.Arabic:  true, // Arapça (ar)
	language.Hebrew:  true, // İbranice (he)
	language.Persian: true, // Farsça (fa)
	language.Urdu:    true, // Urduca (ur)
}

// IsRTL, verilen dilin RTL olup olmadığını kontrol eder
//
// # Kullanım Örneği
//
//	if rtl.IsRTL(language.Arabic) {
//	    // RTL layout kullan
//	}
func IsRTL(lang language.Tag) bool {
	return RTLLanguages[lang]
}

// IsRTLString, dil kodunun RTL olup olmadığını kontrol eder
//
// # Kullanım Örneği
//
//	if rtl.IsRTLString("ar") {
//	    // RTL layout kullan
//	}
func IsRTLString(langCode string) bool {
	tag, err := language.Parse(langCode)
	if err != nil {
		return false
	}
	return IsRTL(tag)
}

// GetDirection, dil için text direction döndürür ("ltr" veya "rtl")
//
// # Kullanım Örneği
//
//	dir := rtl.GetDirection(language.Arabic)
//	// dir = "rtl"
func GetDirection(lang language.Tag) string {
	if IsRTL(lang) {
		return "rtl"
	}
	return "ltr"
}

// GetDirectionString, dil kodu için text direction döndürür
//
// # Kullanım Örneği
//
//	dir := rtl.GetDirectionString("ar")
//	// dir = "rtl"
func GetDirectionString(langCode string) string {
	tag, err := language.Parse(langCode)
	if err != nil {
		return "ltr"
	}
	return GetDirection(tag)
}

// GetDirectionFromContext, Fiber context'inden mevcut dil için direction döndürür
//
// # Kullanım Örneği
//
//	func MyHandler(c *fiber.Ctx) error {
//	    dir := rtl.GetDirectionFromContext(c)
//	    return c.JSON(fiber.Map{"direction": dir})
//	}
func GetDirectionFromContext(c *fiber.Ctx) string {
	// Accept-Language header'ından veya query param'dan dil al
	lang := c.Query("lang")
	if lang == "" {
		lang = c.Get("Accept-Language")
	}

	if lang == "" {
		return "ltr"
	}

	return GetDirectionString(lang)
}

// Middleware, RTL bilgisini response header'larına ekleyen middleware
//
// # Kullanım Örneği
//
//	app.Use(rtl.Middleware())
func Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		dir := GetDirectionFromContext(c)
		c.Set("X-Text-Direction", dir)
		return c.Next()
	}
}
