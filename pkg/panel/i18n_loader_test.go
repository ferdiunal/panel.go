package panel

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestMergedLocaleLoader_MergesEmbeddedAndUser(t *testing.T) {
	tempDir := t.TempDir()
	userLocalePath := filepath.Join(tempDir, "en.yaml")
	userLocale := `button:
  cancel: "Abort"
custom:
  greeting: "hello"`

	if err := os.WriteFile(userLocalePath, []byte(userLocale), 0644); err != nil {
		t.Fatalf("failed to write user locale file: %v", err)
	}

	loader := newMergedLocaleLoader(tempDir)
	data, err := loader.LoadMessage(filepath.Join(tempDir, "en.yaml"))
	if err != nil {
		t.Fatalf("failed to load merged locale: %v", err)
	}

	var merged map[string]interface{}
	if err := yaml.Unmarshal(data, &merged); err != nil {
		t.Fatalf("failed to parse merged locale: %v", err)
	}

	button, ok := merged["button"].(map[string]interface{})
	if !ok {
		t.Fatalf("button section not found in merged locale")
	}

	if got := button["cancel"]; got != "Abort" {
		t.Fatalf("expected user override for button.cancel, got %v", got)
	}

	if got := button["create"]; got != "Create" {
		t.Fatalf("expected embedded fallback for button.create, got %v", got)
	}
}

func TestMergedLocaleLoader_RegionalLangFallsBackToBase(t *testing.T) {
	tempDir := t.TempDir()
	userLocalePath := filepath.Join(tempDir, "en.yaml")
	userLocale := `button:
  cancel: "Stop"`

	if err := os.WriteFile(userLocalePath, []byte(userLocale), 0644); err != nil {
		t.Fatalf("failed to write user locale file: %v", err)
	}

	loader := newMergedLocaleLoader(tempDir)
	data, err := loader.LoadMessage(filepath.Join(tempDir, "en-US.yaml"))
	if err != nil {
		t.Fatalf("failed to load regional merged locale: %v", err)
	}

	var merged map[string]interface{}
	if err := yaml.Unmarshal(data, &merged); err != nil {
		t.Fatalf("failed to parse merged locale: %v", err)
	}

	button, ok := merged["button"].(map[string]interface{})
	if !ok {
		t.Fatalf("button section not found in merged locale")
	}

	if got := button["cancel"]; got != "Stop" {
		t.Fatalf("expected base locale user override for button.cancel, got %v", got)
	}

	if got := button["create"]; got != "Create" {
		t.Fatalf("expected base locale embedded fallback for button.create, got %v", got)
	}
}
