package data

import "testing"

func TestShouldSkipPreloadForViaResource(t *testing.T) {
	lookup := map[string]string{
		"Category":        "categories",
		"ProductVariants": "product_variants",
	}

	tests := []struct {
		name       string
		preload    string
		via        string
		expectSkip bool
	}{
		{
			name:       "skip direct matching relation",
			preload:    "Category",
			via:        "categories",
			expectSkip: true,
		},
		{
			name:       "skip nested relation by root",
			preload:    "Category.Children",
			via:        "categories",
			expectSkip: true,
		},
		{
			name:       "skip when hyphen and underscore differ",
			preload:    "ProductVariants",
			via:        "product-variants",
			expectSkip: true,
		},
		{
			name:       "do not skip unrelated relation",
			preload:    "ProductVariants",
			via:        "categories",
			expectSkip: false,
		},
		{
			name:       "do not skip when via resource empty",
			preload:    "Category",
			via:        "",
			expectSkip: false,
		},
		{
			name:       "do not skip unknown relation",
			preload:    "UnknownRelation",
			via:        "categories",
			expectSkip: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldSkipPreloadForViaResource(tt.preload, tt.via, lookup)
			if got != tt.expectSkip {
				t.Fatalf("expected skip=%v, got %v", tt.expectSkip, got)
			}
		})
	}
}

func TestNormalizeResourceTableIdentifier(t *testing.T) {
	got := normalizeResourceTableIdentifier(" Product-Variants ")
	if got != "product_variants" {
		t.Fatalf("expected normalized identifier to be product_variants, got %q", got)
	}
}
