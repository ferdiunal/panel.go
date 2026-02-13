package resources

import (
	"testing"

	"github.com/ferdiunal/panel.go/pkg/resource"
)

// TestRegister - Registry'ye resource kaydı testi
func TestRegister(t *testing.T) {
	// Test için registry'yi temizle
	Clear()

	// Mock resource factory
	mockFactory := func() resource.Resource {
		return nil
	}

	// Resource'u kaydet
	Register("test-resource", mockFactory)

	// Kayıtlı resource'ları kontrol et
	slugs := List()
	if len(slugs) != 1 {
		t.Errorf("Expected 1 registered resource, got %d", len(slugs))
	}

	if slugs[0] != "test-resource" {
		t.Errorf("Expected slug 'test-resource', got '%s'", slugs[0])
	}
}

// TestGet - Registry'den resource alma testi
func TestGet(t *testing.T) {
	// Test için registry'yi temizle
	Clear()

	// Mock resource factory
	mockFactory := func() resource.Resource {
		return nil
	}

	// Resource'u kaydet
	Register("test-resource", mockFactory)

	// Resource'u al
	res := Get("test-resource")
	if res != nil {
		t.Error("Expected nil resource from mock factory")
	}

	// Olmayan resource'u al
	res = Get("non-existent")
	if res != nil {
		t.Error("Expected nil for non-existent resource")
	}
}

// TestGetOrPanic - GetOrPanic fonksiyonu testi
func TestGetOrPanic(t *testing.T) {
	// Test için registry'yi temizle
	Clear()

	// Mock resource factory (nil döndürür)
	mockFactory := func() resource.Resource {
		return nil
	}

	// Resource'u kaydet
	Register("test-resource", mockFactory)

	// Get() ile resource'u al (nil döndürmeli)
	res := Get("test-resource")
	if res != nil {
		t.Error("Expected nil resource from mock factory")
	}
}

// TestGetOrPanicWithNonExistent - Olmayan resource için panic testi
func TestGetOrPanicWithNonExistent(t *testing.T) {
	// Test için registry'yi temizle
	Clear()

	// Panic bekleniyor
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic for non-existent resource")
		} else {
			expectedMsg := "resource not found: non-existent"
			if r != expectedMsg {
				t.Errorf("Expected panic message '%s', got '%v'", expectedMsg, r)
			}
		}
	}()

	// Olmayan resource'u al (panic olmalı)
	GetOrPanic("non-existent")
}

// TestList - Kayıtlı resource'ları listeleme testi
func TestList(t *testing.T) {
	// Test için registry'yi temizle
	Clear()

	// Mock resource factory
	mockFactory := func() resource.Resource {
		return nil
	}

	// Birden fazla resource kaydet
	Register("resource-1", mockFactory)
	Register("resource-2", mockFactory)
	Register("resource-3", mockFactory)

	// Kayıtlı resource'ları listele
	slugs := List()
	if len(slugs) != 3 {
		t.Errorf("Expected 3 registered resources, got %d", len(slugs))
	}

	// Slug'ların varlığını kontrol et
	slugMap := make(map[string]bool)
	for _, slug := range slugs {
		slugMap[slug] = true
	}

	expectedSlugs := []string{"resource-1", "resource-2", "resource-3"}
	for _, expected := range expectedSlugs {
		if !slugMap[expected] {
			t.Errorf("Expected slug '%s' not found in list", expected)
		}
	}
}

// TestClear - Registry temizleme testi
func TestClear(t *testing.T) {
	// Test için registry'yi temizle
	Clear()

	// Mock resource factory
	mockFactory := func() resource.Resource {
		return nil
	}

	// Resource'ları kaydet
	Register("resource-1", mockFactory)
	Register("resource-2", mockFactory)

	// Kayıtlı resource sayısını kontrol et
	slugs := List()
	if len(slugs) != 2 {
		t.Errorf("Expected 2 registered resources before clear, got %d", len(slugs))
	}

	// Registry'yi temizle
	Clear()

	// Temizlendikten sonra kontrol et
	slugs = List()
	if len(slugs) != 0 {
		t.Errorf("Expected 0 registered resources after clear, got %d", len(slugs))
	}
}

// TestConcurrentAccess - Eşzamanlı erişim testi (thread-safety)
func TestConcurrentAccess(t *testing.T) {
	// Test için registry'yi temizle
	Clear()

	// Mock resource factory
	mockFactory := func() resource.Resource {
		return nil
	}

	// Eşzamanlı kayıt
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(index int) {
			Register("resource-"+string(rune('0'+index)), mockFactory)
			done <- true
		}(i)
	}

	// Tüm goroutine'lerin bitmesini bekle
	for i := 0; i < 10; i++ {
		<-done
	}

	// Kayıtlı resource sayısını kontrol et
	slugs := List()
	if len(slugs) != 10 {
		t.Errorf("Expected 10 registered resources, got %d", len(slugs))
	}

	// Eşzamanlı okuma
	for i := 0; i < 10; i++ {
		go func(index int) {
			Get("resource-" + string(rune('0'+index)))
			done <- true
		}(i)
	}

	// Tüm goroutine'lerin bitmesini bekle
	for i := 0; i < 10; i++ {
		<-done
	}
}
