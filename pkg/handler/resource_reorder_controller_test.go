package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	appContext "github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/data"
	"github.com/ferdiunal/panel.go/pkg/resource"
	"github.com/gofiber/fiber/v2"
)

type reorderUpdateCall struct {
	id   string
	data map[string]interface{}
}

type reorderTrackingProvider struct {
	MockDataProvider
	updates       []reorderUpdateCall
	beginCalled   int
	commitCalled  int
	rollbackCalls int
}

func (p *reorderTrackingProvider) Update(ctx *appContext.Context, id string, values map[string]interface{}) (interface{}, error) {
	p.updates = append(p.updates, reorderUpdateCall{
		id:   id,
		data: values,
	})
	return values, nil
}

func (p *reorderTrackingProvider) BeginTx(ctx *appContext.Context) (data.DataProvider, error) {
	p.beginCalled++
	return p, nil
}

func (p *reorderTrackingProvider) Commit() error {
	p.commitCalled++
	return nil
}

func (p *reorderTrackingProvider) Rollback() error {
	p.rollbackCalls++
	return nil
}

func TestHandleResourceReorder_Success(t *testing.T) {
	app := fiber.New()
	provider := &reorderTrackingProvider{}

	h := NewFieldHandler(provider)
	h.IndexReorderConfig = resource.IndexReorderConfig{
		Enabled: true,
		Column:  "order_column",
	}

	app.Post("/users/reorder", appContext.Wrap(func(c *appContext.Context) error {
		return HandleResourceReorder(h, c)
	}))

	payload := map[string]interface{}{
		"ids": []int{7, 3, 9},
	}
	rawBody, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/users/reorder", bytes.NewReader(rawBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("Expected status 200, got %d (%s)", resp.StatusCode, string(body))
	}

	if provider.beginCalled != 1 {
		t.Fatalf("Expected BeginTx to be called once, got %d", provider.beginCalled)
	}

	if provider.commitCalled != 1 {
		t.Fatalf("Expected Commit to be called once, got %d", provider.commitCalled)
	}

	if provider.rollbackCalls != 0 {
		t.Fatalf("Expected Rollback not to be called, got %d", provider.rollbackCalls)
	}

	if len(provider.updates) != 3 {
		t.Fatalf("Expected 3 updates, got %d", len(provider.updates))
	}

	expectedIDs := []string{"7", "3", "9"}
	for i, call := range provider.updates {
		if call.id != expectedIDs[i] {
			t.Fatalf("Expected update id %s, got %s", expectedIDs[i], call.id)
		}

		expectedOrder := i + 1
		if got := call.data["order_column"]; got != expectedOrder {
			t.Fatalf("Expected order_column=%d, got %v", expectedOrder, got)
		}
	}
}

func TestHandleResourceReorder_Disabled(t *testing.T) {
	app := fiber.New()
	provider := &reorderTrackingProvider{}
	h := NewFieldHandler(provider)

	app.Post("/users/reorder", appContext.Wrap(func(c *appContext.Context) error {
		return HandleResourceReorder(h, c)
	}))

	req := httptest.NewRequest("POST", "/users/reorder", bytes.NewReader([]byte(`{"ids":[1]}`)))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	if resp.StatusCode != fiber.StatusBadRequest {
		t.Fatalf("Expected status 400, got %d", resp.StatusCode)
	}
}
