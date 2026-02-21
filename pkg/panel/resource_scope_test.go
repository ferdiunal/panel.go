package panel

import (
	"testing"

	"github.com/ferdiunal/panel.go/pkg/resource"
	resourceAccount "github.com/ferdiunal/panel.go/pkg/resource/account"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newTestPanel(t *testing.T) *Panel {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}

	cfg := Config{
		Database: DatabaseConfig{
			Instance: db,
		},
		Server: ServerConfig{
			Host: "localhost",
			Port: "8080",
		},
		Environment: "test",
	}

	return New(cfg)
}

func TestRegisterResourcePreservesExplicitDialogType(t *testing.T) {
	p := newTestPanel(t)

	accountRes := resourceAccount.NewAccountResource()
	accountRes.SetDialogType(resource.DialogTypeModal)

	p.RegisterResource(accountRes)

	registered, ok := p.resources["accounts"]
	if !ok {
		t.Fatal("expected accounts resource to be registered")
	}

	if registered.GetDialogType() != resource.DialogTypeModal {
		t.Fatalf("expected dialog type %q, got %q", resource.DialogTypeModal, registered.GetDialogType())
	}
}

func TestRegisterResourceSetsDefaultDialogTypeWhenUnset(t *testing.T) {
	p := newTestPanel(t)

	accountRes := resourceAccount.NewAccountResource()
	p.RegisterResource(accountRes)

	registered, ok := p.resources["accounts"]
	if !ok {
		t.Fatal("expected accounts resource to be registered")
	}

	if registered.GetDialogType() != resource.DialogTypeSheet {
		t.Fatalf("expected dialog type %q, got %q", resource.DialogTypeSheet, registered.GetDialogType())
	}
}
