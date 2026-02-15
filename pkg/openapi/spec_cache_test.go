package openapi

import (
	"testing"
	"time"
)

func TestSpecCacheGetReturnsImmutableClone(t *testing.T) {
	cache := &specCache{ttl: time.Minute}
	cache.set(&OpenAPISpec{
		OpenAPI: "3.0.3",
		Paths: map[string]PathItem{
			"/api/test": {},
		},
	})

	first := cache.get()
	if first == nil {
		t.Fatalf("expected cached spec")
	}

	first.OpenAPI = "modified"
	first.Paths["/api/new"] = PathItem{}

	second := cache.get()
	if second == nil {
		t.Fatalf("expected cached spec clone")
	}

	if second.OpenAPI != "3.0.3" {
		t.Fatalf("expected immutable cached OpenAPI version, got %s", second.OpenAPI)
	}
	if _, exists := second.Paths["/api/new"]; exists {
		t.Fatalf("expected cached paths to remain immutable")
	}
}
