package panel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	appContext "github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/resource"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type internalRESTAPIUser struct {
	ID     uint   `json:"id" gorm:"primaryKey"`
	Name   string `json:"name"`
	Secret string `json:"secret"`
}

func (internalRESTAPIUser) TableName() string {
	return "internal_rest_api_users"
}

type internalRESTAPIUserFieldResolver struct{}

func (r *internalRESTAPIUserFieldResolver) ResolveFields(_ *appContext.Context) []fields.Element {
	return []fields.Element{
		fields.ID(),
		fields.Text("Name", "name").Required(),
		fields.Text("Secret", "secret").HideOnApi(),
	}
}

type internalRESTAPIUserResource struct {
	resource.OptimizedBase
}

func newInternalRESTAPIUserResource() *internalRESTAPIUserResource {
	res := &internalRESTAPIUserResource{}
	res.SetModel(&internalRESTAPIUser{})
	res.SetSlug("internal-rest-users")
	res.SetTitle("Internal REST Users")
	res.SetGroup("System")
	res.SetIcon("users")
	res.SetVisible(true)
	res.SetFieldResolver(&internalRESTAPIUserFieldResolver{})
	return res
}

func setupInternalRESTAPIPanel(t *testing.T, cfg Config) *Panel {
	t.Helper()

	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", strings.ReplaceAll(t.Name(), "/", "_"))
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect db: %v", err)
	}
	if err := db.AutoMigrate(&internalRESTAPIUser{}); err != nil {
		t.Fatalf("failed to auto migrate test model: %v", err)
	}
	if err := db.Create(&internalRESTAPIUser{Name: "First User", Secret: "top-secret"}).Error; err != nil {
		t.Fatalf("failed to seed test model: %v", err)
	}

	cfg.Database = DatabaseConfig{Instance: db}
	cfg.Environment = "test"

	p := New(cfg)
	p.RegisterResource(newInternalRESTAPIUserResource())
	return p
}

func TestInternalRESTAPI_FeatureDisabled(t *testing.T) {
	p := setupInternalRESTAPIPanel(t, Config{
		Features: FeatureConfig{
			RestAPI: false,
		},
		RESTAPI: RESTAPIConfig{
			Keys: []string{"internal-secret"},
		},
	})

	req := httptest.NewRequest("GET", "/internal-api/internal-rest-users", nil)
	req.Header.Set("X-Internal-API-Key", "internal-secret")

	resp, err := testFiberRequest(p.Fiber, req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != 404 {
		t.Fatalf("expected status 404 when feature disabled, got %d", resp.StatusCode)
	}
}

func TestInternalRESTAPI_RequiresValidKey(t *testing.T) {
	p := setupInternalRESTAPIPanel(t, Config{
		Features: FeatureConfig{
			RestAPI: true,
		},
		RESTAPI: RESTAPIConfig{
			Keys: []string{"internal-secret"},
		},
	})

	missingKeyReq := httptest.NewRequest("GET", "/internal-api/internal-rest-users", nil)
	missingKeyResp, err := testFiberRequest(p.Fiber, missingKeyReq)
	if err != nil {
		t.Fatalf("missing key request failed: %v", err)
	}
	if missingKeyResp.StatusCode != 401 {
		t.Fatalf("expected status 401 without key, got %d", missingKeyResp.StatusCode)
	}

	invalidKeyReq := httptest.NewRequest("GET", "/internal-api/internal-rest-users", nil)
	invalidKeyReq.Header.Set("X-Internal-API-Key", "wrong-key")
	invalidKeyResp, err := testFiberRequest(p.Fiber, invalidKeyReq)
	if err != nil {
		t.Fatalf("invalid key request failed: %v", err)
	}
	if invalidKeyResp.StatusCode != 401 {
		t.Fatalf("expected status 401 with invalid key, got %d", invalidKeyResp.StatusCode)
	}

	validKeyReq := httptest.NewRequest("GET", "/internal-api/internal-rest-users", nil)
	validKeyReq.Header.Set("X-Internal-API-Key", "internal-secret")
	validKeyResp, err := testFiberRequest(p.Fiber, validKeyReq)
	if err != nil {
		t.Fatalf("valid key request failed: %v", err)
	}
	if validKeyResp.StatusCode != 200 {
		t.Fatalf("expected status 200 with valid key, got %d", validKeyResp.StatusCode)
	}
}

func TestInternalRESTAPI_DetailUpdateDeleteAndValidate(t *testing.T) {
	p := setupInternalRESTAPIPanel(t, Config{
		Features: FeatureConfig{
			RestAPI: true,
		},
		RESTAPI: RESTAPIConfig{
			Keys: []string{"internal-secret"},
		},
	})

	detailReq := httptest.NewRequest("GET", "/internal-api/internal-rest-users/1", nil)
	detailReq.Header.Set("X-Internal-API-Key", "internal-secret")
	detailResp, err := testFiberRequest(p.Fiber, detailReq)
	if err != nil {
		t.Fatalf("detail request failed: %v", err)
	}
	if detailResp.StatusCode != 200 {
		t.Fatalf("expected status 200 from detail endpoint, got %d", detailResp.StatusCode)
	}

	invalidUpdateBody, _ := json.Marshal(map[string]any{
		"name": "",
	})
	updateReq := httptest.NewRequest("PUT", "/internal-api/internal-rest-users/1", bytes.NewReader(invalidUpdateBody))
	updateReq.Header.Set("Content-Type", "application/json")
	updateReq.Header.Set("X-Internal-API-Key", "internal-secret")
	updateResp, err := testFiberRequest(p.Fiber, updateReq)
	if err != nil {
		t.Fatalf("update request failed: %v", err)
	}
	if updateResp.StatusCode != 422 {
		t.Fatalf("expected status 422 for validation error, got %d", updateResp.StatusCode)
	}

	deleteReq := httptest.NewRequest("DELETE", "/internal-api/internal-rest-users/1", nil)
	deleteReq.Header.Set("X-Internal-API-Key", "internal-secret")
	deleteResp, err := testFiberRequest(p.Fiber, deleteReq)
	if err != nil {
		t.Fatalf("delete request failed: %v", err)
	}
	if deleteResp.StatusCode != 200 {
		t.Fatalf("expected status 200 from delete endpoint, got %d", deleteResp.StatusCode)
	}
}
