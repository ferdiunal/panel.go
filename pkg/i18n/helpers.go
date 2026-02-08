// Package i18n provides internationalization helpers for Panel.go
// Laravel'deki __() helper'ına benzer şekilde çalışır
package i18n

import (
	"github.com/gofiber/contrib/fiberi18n/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// Trans, Laravel'deki __() helper'ına benzer şekilde çeviri yapar
// Basit kullanım: Trans(c, "welcome")
// Template ile: Trans(c, "welcomeWithName", map[string]string{"Name": "Ahmet"})
//
// # Kullanım Örnekleri
//
// Basit çeviri:
//
//	message := i18n.Trans(c, "welcome")
//	// Çıktı: "Hoş geldiniz"
//
// Template değişkenleri ile:
//
//	message := i18n.Trans(c, "welcomeWithName", map[string]string{
//	    "Name": "Ahmet",
//	})
//	// Çıktı: "Hoş geldiniz, Ahmet"
//
// Çoklu parametre:
//
//	message := i18n.Trans(c, "orderSummary", map[string]interface{}{
//	    "Count": 5,
//	    "Total": 150.50,
//	})
//	// Çıktı: "5 ürün, toplam: 150.50 TL"
func Trans(c *fiber.Ctx, messageID string, templateData ...map[string]interface{}) string {
	// Nil context kontrolü - resource initialization sırasında context olmayabilir
	if c == nil {
		return messageID // Fallback: messageID'yi döndür
	}

	if len(templateData) > 0 && templateData[0] != nil {
		return fiberi18n.MustLocalize(c, &i18n.LocalizeConfig{
			MessageID:    messageID,
			TemplateData: templateData[0],
		})
	}
	return fiberi18n.MustLocalize(c, messageID)
}

// TransChoice, Laravel'deki trans_choice() helper'ına benzer şekilde çoğul çeviri yapar
// Sayıya göre doğru çoğul formunu seçer
//
// # Kullanım Örnekleri
//
// Tekil/Çoğul:
//
//	message := i18n.TransChoice(c, "items", 1)
//	// Çıktı: "1 öğe"
//
//	message := i18n.TransChoice(c, "items", 5)
//	// Çıktı: "5 öğe"
//
// Template değişkenleri ile:
//
//	message := i18n.TransChoice(c, "itemsWithName", 3, map[string]interface{}{
//	    "Name": "Ürün",
//	})
//	// Çıktı: "3 Ürün"
func TransChoice(c *fiber.Ctx, messageID string, count int, templateData ...map[string]interface{}) string {
	// Nil context kontrolü
	if c == nil {
		return messageID // Fallback: messageID'yi döndür
	}

	data := map[string]interface{}{
		"Count": count,
	}

	// Eğer ek template data varsa birleştir
	if len(templateData) > 0 && templateData[0] != nil {
		for k, v := range templateData[0] {
			data[k] = v
		}
	}

	return fiberi18n.MustLocalize(c, &i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: data,
		PluralCount:  count,
	})
}

// GetLocale, mevcut dili döndürür
// Laravel'deki app()->getLocale() metoduna benzer
//
// # Kullanım Örneği
//
//	lang := i18n.GetLocale(c)
//	// Çıktı: "tr" veya "en"
func GetLocale(c *fiber.Ctx) string {
	// Nil context kontrolü
	if c == nil {
		return "en" // Varsayılan dil
	}

	// fiberi18n middleware'i locale bilgisini c.Locals() ile kaydeder
	if locale, ok := c.Locals("lang").(string); ok {
		return locale
	}
	// Varsayılan dil
	return "en"
}

// HasTranslation, çevirinin var olup olmadığını kontrol eder
// Laravel'deki Lang::has() metoduna benzer
//
// # Kullanım Örneği
//
//	if i18n.HasTranslation(c, "welcome") {
//	    message := i18n.Trans(c, "welcome")
//	}
func HasTranslation(c *fiber.Ctx, messageID string) bool {
	// Nil context kontrolü
	if c == nil {
		return false
	}

	// Çeviriyi dene, hata varsa false döndür
	_, err := fiberi18n.Localize(c, messageID)
	return err == nil
}

// TransWithFallback, çeviri yoksa fallback değeri döndürür
// Laravel'deki __() helper'ının fallback özelliğine benzer
//
// # Kullanım Örneği
//
//	message := i18n.TransWithFallback(c, "unknown.key", "Varsayılan Mesaj")
//	// Çıktı: "Varsayılan Mesaj" (çeviri yoksa)
func TransWithFallback(c *fiber.Ctx, messageID string, fallback string, templateData ...map[string]interface{}) string {
	// Nil context kontrolü
	if c == nil {
		return fallback
	}

	if !HasTranslation(c, messageID) {
		return fallback
	}
	return Trans(c, messageID, templateData...)
}
