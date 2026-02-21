package resource

import (
	"fmt"
	"strings"

	"github.com/ferdiunal/panel.go/pkg/i18n"
	"github.com/gofiber/fiber/v2"
)

type descriptionWithContextProvider interface {
	DescriptionWithContext(ctx *fiber.Ctx) string
}

type descriptionProvider interface {
	Description() string
}

// DescriptionWithContext resolves a resource description with best-effort fallbacks.
// Priority:
//  1. Optional DescriptionWithContext(ctx) on the resource
//  2. Optional Description() on the resource
//  3. i18n key: resources.<slug>.description (and slug with "-" -> "_")
func DescriptionWithContext(res Resource, ctx *fiber.Ctx, slug string) string {
	if res == nil {
		return ""
	}

	if ctx != nil {
		if describer, ok := any(res).(descriptionWithContextProvider); ok {
			if description := strings.TrimSpace(describer.DescriptionWithContext(ctx)); description != "" {
				return description
			}
		}
	}

	if describer, ok := any(res).(descriptionProvider); ok {
		if description := strings.TrimSpace(describer.Description()); description != "" {
			return description
		}
	}

	if ctx == nil {
		return ""
	}

	resourceSlug := strings.TrimSpace(slug)
	if resourceSlug == "" {
		resourceSlug = strings.TrimSpace(res.Slug())
	}
	if resourceSlug == "" {
		return ""
	}

	keys := []string{
		fmt.Sprintf("resources.%s.description", resourceSlug),
	}
	normalizedSlug := strings.ReplaceAll(resourceSlug, "-", "_")
	if normalizedSlug != resourceSlug {
		keys = append(keys, fmt.Sprintf("resources.%s.description", normalizedSlug))
	}

	for _, key := range keys {
		translated := i18n.Trans(ctx, key)
		if translated == key {
			continue
		}
		if description := strings.TrimSpace(translated); description != "" {
			return description
		}
	}

	return ""
}
