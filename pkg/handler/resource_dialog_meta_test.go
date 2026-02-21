package handler

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	appContext "github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/core"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/resource"
	"github.com/gofiber/fiber/v2"
)

type dialogMetaDataProvider struct {
	MockDataProvider
	item interface{}
}

func (p *dialogMetaDataProvider) Show(ctx *appContext.Context, id string) (interface{}, error) {
	return p.item, nil
}

func TestHandleResourceDetail_IncludesDialogMeta(t *testing.T) {
	app := fiber.New()

	fieldDefs := []fields.Element{
		fields.ID(),
		fields.Text("Full Name", "full_name"),
	}

	provider := &dialogMetaDataProvider{
		item: User{ID: 1, FullName: "Detail User"},
	}

	h := NewFieldHandler(provider)
	h.Resource = &MockResource{}
	h.DialogType = resource.DialogTypeModal
	h.DialogSize = resource.DialogSizeLG

	app.Get(
		"/users/:id/detail",
		FieldContextMiddleware(nil, h.Resource, core.ContextDetail, fieldDefs),
		appContext.Wrap(h.Detail),
	)

	req := httptest.NewRequest("GET", "/users/1/detail", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	meta, ok := payload["meta"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected meta map, got %T", payload["meta"])
	}
	if meta["dialog_type"] != string(resource.DialogTypeModal) {
		t.Fatalf("expected dialog_type %q, got %v", resource.DialogTypeModal, meta["dialog_type"])
	}
	if meta["dialog_size"] != string(resource.DialogSizeLG) {
		t.Fatalf("expected dialog_size %q, got %v", resource.DialogSizeLG, meta["dialog_size"])
	}
}

func TestHandleResourceEdit_IncludesDialogMeta(t *testing.T) {
	app := fiber.New()

	fieldDefs := []fields.Element{
		fields.ID(),
		fields.Text("Full Name", "full_name"),
	}

	provider := &dialogMetaDataProvider{
		item: User{ID: 1, FullName: "Edit User"},
	}

	h := NewFieldHandler(provider)
	h.Resource = &MockResource{}
	h.DialogType = resource.DialogTypeModal
	h.DialogSize = resource.DialogSizeLG

	app.Get(
		"/users/:id/edit",
		FieldContextMiddleware(nil, h.Resource, core.ContextUpdate, fieldDefs),
		appContext.Wrap(h.Edit),
	)

	req := httptest.NewRequest("GET", "/users/1/edit", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != fiber.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	meta, ok := payload["meta"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected meta map, got %T", payload["meta"])
	}
	if meta["dialog_type"] != string(resource.DialogTypeModal) {
		t.Fatalf("expected dialog_type %q, got %v", resource.DialogTypeModal, meta["dialog_type"])
	}
	if meta["dialog_size"] != string(resource.DialogSizeLG) {
		t.Fatalf("expected dialog_size %q, got %v", resource.DialogSizeLG, meta["dialog_size"])
	}
}
