package panel

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadTranslations_MergesEmbeddedAndUser(t *testing.T) {
	tempDir := t.TempDir()
	userLocalePath := filepath.Join(tempDir, "en.yaml")

	userLocale := `button:
  cancel: "Abort"
custom:
  greeting: "Hello from user"`

	if err := os.WriteFile(userLocalePath, []byte(userLocale), 0644); err != nil {
		t.Fatalf("failed to write user locale file: %v", err)
	}

	config := Config{}
	config.I18n.RootPath = tempDir

	translations := loadTranslations(config, "en")

	if got := translations["button.cancel"]; got != "Abort" {
		t.Fatalf("expected user override for button.cancel, got %v", got)
	}

	if got := translations["button.create"]; got != "Create" {
		t.Fatalf("expected embedded fallback for button.create, got %v", got)
	}

	if got := translations["custom.greeting"]; got != "Hello from user" {
		t.Fatalf("expected custom user key, got %v", got)
	}
}

func TestLoadTranslations_FallsBackToEmbeddedWhenUserMissing(t *testing.T) {
	config := Config{}
	config.I18n.RootPath = t.TempDir()

	translations := loadTranslations(config, "tr")

	if got := translations["button.create"]; got != "Oluştur" {
		t.Fatalf("expected embedded translation for button.create, got %v", got)
	}

	if got := translations["button.cancel"]; got != "İptal" {
		t.Fatalf("expected embedded translation for button.cancel, got %v", got)
	}
}

func TestLoadTranslations_UsesBaseLocaleFallbackForRegionalLang(t *testing.T) {
	tempDir := t.TempDir()
	userLocalePath := filepath.Join(tempDir, "en.yaml")

	userLocale := `button:
  cancel: "Stop"`

	if err := os.WriteFile(userLocalePath, []byte(userLocale), 0644); err != nil {
		t.Fatalf("failed to write user locale file: %v", err)
	}

	config := Config{}
	config.I18n.RootPath = tempDir

	translations := loadTranslations(config, "en-US")

	if got := translations["button.cancel"]; got != "Stop" {
		t.Fatalf("expected base locale fallback override for button.cancel, got %v", got)
	}

	if got := translations["button.create"]; got != "Create" {
		t.Fatalf("expected embedded base locale fallback for button.create, got %v", got)
	}
}
