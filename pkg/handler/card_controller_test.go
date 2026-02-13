package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	appContext "github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/widget"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// MockCard implements widget.Card interface for testing
type MockCard struct {
	name      string
	component string
	width     string
	cardType  widget.CardType
	data      interface{}
	err       error
}

func (m *MockCard) Name() string      { return m.name }
func (m *MockCard) Component() string { return m.component }
func (m *MockCard) Width() string     { return m.width }
func (m *MockCard) GetType() widget.CardType {
	if m.cardType == "" {
		return widget.CardTypeValue
	}
	return m.cardType
}
func (m *MockCard) Resolve(c *appContext.Context, db *gorm.DB) (interface{}, error) {
	return m.data, m.err
}
func (m *MockCard) HandleError(err error) map[string]interface{} {
	return map[string]interface{}{
		"error": err.Error(),
		"title": m.name,
	}
}
func (m *MockCard) GetMetadata() map[string]interface{} {
	return map[string]interface{}{
		"name":      m.name,
		"component": m.component,
		"width":     m.width,
		"type":      m.GetType(),
	}
}
func (m *MockCard) JsonSerialize() map[string]interface{} {
	return map[string]interface{}{
		"name":      m.name,
		"component": m.component,
		"width":     m.width,
		"type":      m.GetType(),
	}
}

func TestHandleCardList_Success(t *testing.T) {
	app := fiber.New()

	h := &FieldHandler{
		Cards: []widget.Card{
			&MockCard{name: "Card1", component: "ValueCard", width: "1/3", data: map[string]interface{}{"value": 100}},
			&MockCard{name: "Card2", component: "ChartCard", width: "2/3", data: map[string]interface{}{"values": []int{1, 2, 3}}},
		},
	}

	app.Get("/cards", appContext.Wrap(func(c *appContext.Context) error {
		return HandleCardList(h, c)
	}))

	req := httptest.NewRequest("GET", "/cards", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	dataList := response["data"].([]interface{})
	if len(dataList) != 2 {
		t.Errorf("Expected 2 cards, got %d", len(dataList))
	}
}

func TestHandleCardList_WithError(t *testing.T) {
	app := fiber.New()

	h := &FieldHandler{
		Cards: []widget.Card{
			&MockCard{name: "Card1", component: "ValueCard", width: "1/3", err: errors.New("resolution failed")},
		},
	}

	app.Get("/cards", appContext.Wrap(func(c *appContext.Context) error {
		return HandleCardList(h, c)
	}))

	req := httptest.NewRequest("GET", "/cards", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	dataList := response["data"].([]interface{})
	if len(dataList) != 1 {
		t.Errorf("Expected 1 card, got %d", len(dataList))
	}

	card := dataList[0].(map[string]interface{})
	if card["error"] == nil {
		t.Error("Expected error field in card response")
	}
}

func TestHandleCardDetail_Success(t *testing.T) {
	app := fiber.New()

	h := &FieldHandler{
		Cards: []widget.Card{
			&MockCard{name: "Card1", component: "ValueCard", width: "1/3", data: map[string]interface{}{"value": 100}},
		},
	}

	app.Get("/cards/:index", appContext.Wrap(func(c *appContext.Context) error {
		return HandleCardDetail(h, c)
	}))

	req := httptest.NewRequest("GET", "/cards/0", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["data"] == nil {
		t.Error("Expected data field in response")
	}
}

func TestHandleCardDetail_InvalidIndex(t *testing.T) {
	app := fiber.New()

	h := &FieldHandler{
		Cards: []widget.Card{
			&MockCard{name: "Card1", component: "ValueCard", width: "1/3", data: map[string]interface{}{"value": 100}},
		},
	}

	app.Get("/cards/:index", appContext.Wrap(func(c *appContext.Context) error {
		return HandleCardDetail(h, c)
	}))

	req := httptest.NewRequest("GET", "/cards/999", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	if resp.StatusCode != 404 {
		t.Errorf("Expected status 404, got %d", resp.StatusCode)
	}
}

func TestHandleCardDetail_ResolutionError(t *testing.T) {
	app := fiber.New()

	h := &FieldHandler{
		Cards: []widget.Card{
			&MockCard{name: "Card1", component: "ValueCard", width: "1/3", err: errors.New("resolution failed")},
		},
	}

	app.Get("/cards/:index", appContext.Wrap(func(c *appContext.Context) error {
		return HandleCardDetail(h, c)
	}))

	req := httptest.NewRequest("GET", "/cards/0", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	if resp.StatusCode != 500 {
		t.Errorf("Expected status 500, got %d", resp.StatusCode)
	}
}
