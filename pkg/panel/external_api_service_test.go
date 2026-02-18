package panel

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"
)

func TestExternalAPI_FeatureDisabled(t *testing.T) {
	p := setupInternalRESTAPIPanel(t, Config{
		Features: FeatureConfig{
			ExternalAPI: false,
		},
		ExternalAPI: ExternalAPIConfig{
			Keys: []string{"external-secret"},
		},
	})

	req := httptest.NewRequest("GET", "/api/internal-rest-users", nil)
	req.Header.Set("X-External-API-Key", "external-secret")

	resp, err := testFiberRequest(p.Fiber, req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != 404 {
		t.Fatalf("expected status 404 when feature disabled, got %d", resp.StatusCode)
	}
}

func TestExternalAPI_RequiresValidKey(t *testing.T) {
	p := setupInternalRESTAPIPanel(t, Config{
		Features: FeatureConfig{
			ExternalAPI: true,
		},
		ExternalAPI: ExternalAPIConfig{
			Keys: []string{"external-secret"},
		},
	})

	missingKeyReq := httptest.NewRequest("GET", "/api/internal-rest-users", nil)
	missingKeyResp, err := testFiberRequest(p.Fiber, missingKeyReq)
	if err != nil {
		t.Fatalf("missing key request failed: %v", err)
	}
	if missingKeyResp.StatusCode != 401 {
		t.Fatalf("expected status 401 without key, got %d", missingKeyResp.StatusCode)
	}

	invalidKeyReq := httptest.NewRequest("GET", "/api/internal-rest-users", nil)
	invalidKeyReq.Header.Set("X-External-API-Key", "wrong-key")
	invalidKeyResp, err := testFiberRequest(p.Fiber, invalidKeyReq)
	if err != nil {
		t.Fatalf("invalid key request failed: %v", err)
	}
	if invalidKeyResp.StatusCode != 401 {
		t.Fatalf("expected status 401 with invalid key, got %d", invalidKeyResp.StatusCode)
	}
}

func TestExternalAPI_ReturnsPlainFieldValues(t *testing.T) {
	p := setupInternalRESTAPIPanel(t, Config{
		Features: FeatureConfig{
			ExternalAPI: true,
		},
		ExternalAPI: ExternalAPIConfig{
			Keys: []string{"external-secret"},
		},
	})

	indexReq := httptest.NewRequest("GET", "/api/internal-rest-users", nil)
	indexReq.Header.Set("X-External-API-Key", "external-secret")
	indexResp, err := testFiberRequest(p.Fiber, indexReq)
	if err != nil {
		t.Fatalf("index request failed: %v", err)
	}
	if indexResp.StatusCode != 200 {
		t.Fatalf("expected status 200 for index, got %d", indexResp.StatusCode)
	}

	var indexPayload map[string]interface{}
	if err := json.NewDecoder(indexResp.Body).Decode(&indexPayload); err != nil {
		t.Fatalf("failed to decode index response: %v", err)
	}

	indexData, ok := indexPayload["data"].([]interface{})
	if !ok || len(indexData) == 0 {
		t.Fatalf("expected non-empty data array in index response")
	}

	firstRecord, ok := indexData[0].(map[string]interface{})
	if !ok {
		t.Fatalf("expected record map in index data")
	}

	nameValue, exists := firstRecord["name"]
	if !exists {
		t.Fatalf("expected name field in flattened index record")
	}
	if _, isWrapped := nameValue.(map[string]interface{}); isWrapped {
		t.Fatalf("expected plain name value, got wrapped field payload")
	}
	if _, exists := firstRecord["secret"]; exists {
		t.Fatalf("expected secret field to be hidden by HideOnApi on index response")
	}

	showReq := httptest.NewRequest("GET", "/api/internal-rest-users/1", nil)
	showReq.Header.Set("X-External-API-Key", "external-secret")
	showResp, err := testFiberRequest(p.Fiber, showReq)
	if err != nil {
		t.Fatalf("show request failed: %v", err)
	}
	if showResp.StatusCode != 200 {
		t.Fatalf("expected status 200 for show, got %d", showResp.StatusCode)
	}

	var showPayload map[string]interface{}
	if err := json.NewDecoder(showResp.Body).Decode(&showPayload); err != nil {
		t.Fatalf("failed to decode show response: %v", err)
	}

	showData, ok := showPayload["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected data object in show response")
	}

	if _, isWrapped := showData["name"].(map[string]interface{}); isWrapped {
		t.Fatalf("expected plain name value in show response")
	}
	if _, exists := showData["secret"]; exists {
		t.Fatalf("expected secret field to be hidden by HideOnApi on show response")
	}
}

func TestExternalAPI_AllowsPanelAPIKey(t *testing.T) {
	p := setupInternalRESTAPIPanel(t, Config{
		Features: FeatureConfig{
			ExternalAPI: true,
		},
		APIKey: APIKeyConfig{
			Enabled: true,
			Header:  "X-API-Key",
			Keys:    []string{"panel-shared-key"},
		},
	})

	req := httptest.NewRequest("GET", "/api/internal-rest-users", nil)
	req.Header.Set("X-API-Key", "panel-shared-key")

	resp, err := testFiberRequest(p.Fiber, req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Fatalf("expected status 200 with panel api key, got %d", resp.StatusCode)
	}
}

func TestExternalAPI_ValidationMatchesResourceRules(t *testing.T) {
	p := setupInternalRESTAPIPanel(t, Config{
		Features: FeatureConfig{
			ExternalAPI: true,
		},
		ExternalAPI: ExternalAPIConfig{
			Keys: []string{"external-secret"},
		},
	})

	body, _ := json.Marshal(map[string]any{
		"name": "",
	})
	updateReq := httptest.NewRequest("PUT", "/api/internal-rest-users/1", bytes.NewReader(body))
	updateReq.Header.Set("Content-Type", "application/json")
	updateReq.Header.Set("X-External-API-Key", "external-secret")
	updateResp, err := testFiberRequest(p.Fiber, updateReq)
	if err != nil {
		t.Fatalf("update request failed: %v", err)
	}

	if updateResp.StatusCode != 422 {
		t.Fatalf("expected status 422 for validation error, got %d", updateResp.StatusCode)
	}

	var payload map[string]interface{}
	if err := json.NewDecoder(updateResp.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode validation response: %v", err)
	}

	errorsMap, ok := payload["errors"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected errors object in validation response")
	}

	if _, exists := errorsMap["name"]; !exists {
		t.Fatalf("expected name validation error")
	}
}
