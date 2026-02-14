package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	appContext "github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/gofiber/fiber/v2"
)

type relationTestResource struct {
	*MockResource
	elements []fields.Element
}

func (r *relationTestResource) Fields() []fields.Element {
	return r.elements
}

func newRelationTestHandler(elements []fields.Element) *FieldHandler {
	h := NewFieldHandler(&MockDataProvider{})
	h.Resource = &relationTestResource{
		MockResource: &MockResource{},
		elements:     elements,
	}
	return h
}

func runResolveDependenciesRequest(
	t *testing.T,
	h *FieldHandler,
	payload map[string]interface{},
) (int, map[string]interface{}) {
	t.Helper()

	app := fiber.New()
	app.Post("/resolve", appContext.Wrap(func(c *appContext.Context) error {
		return HandleResolveDependencies(h, c)
	}))

	raw, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("failed to marshal payload: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/resolve", bytes.NewReader(raw))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	decoded := map[string]interface{}{}
	if len(bytes.TrimSpace(respBody)) > 0 {
		if err := json.Unmarshal(respBody, &decoded); err != nil {
			t.Fatalf("failed to unmarshal response body: %v", err)
		}
	}

	return resp.StatusCode, decoded
}

func runParseBodyRequest(
	t *testing.T,
	h *FieldHandler,
	contentType string,
	body io.Reader,
) map[string]interface{} {
	t.Helper()

	app := fiber.New()

	var parsed map[string]interface{}
	var parseErr error

	app.Post("/parse", appContext.Wrap(func(c *appContext.Context) error {
		parsed, parseErr = h.parseBody(c)
		if parseErr != nil {
			return parseErr
		}
		return c.SendStatus(fiber.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodPost, "/parse", body)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if parseErr != nil {
		t.Fatalf("parseBody returned error: %v", parseErr)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		respBody, _ := io.ReadAll(resp.Body)
		t.Fatalf("unexpected status code %d, body: %s", resp.StatusCode, string(respBody))
	}

	return parsed
}

func TestHandleResolveDependencies_AcceptsEditAlias(t *testing.T) {
	cityField := fields.Select("City", "city_id")
	cityField.DependsOn("country_id")
	cityField.OnDependencyChangeUpdating(func(
		field *fields.Schema,
		formData map[string]interface{},
		ctx *fiber.Ctx,
	) *fields.FieldUpdate {
		return fields.NewFieldUpdate().SetValue("from-update-context")
	})

	h := newRelationTestHandler([]fields.Element{
		fields.Select("Country", "country_id"),
		cityField,
	})

	statusCode, response := runResolveDependenciesRequest(t, h, map[string]interface{}{
		"formData": map[string]interface{}{
			"country_id": "1",
		},
		"context":       "edit",
		"changedFields": []string{"country_id"},
	})

	if statusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, statusCode)
	}

	updates, ok := response["fields"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected response.fields to be map, got %T", response["fields"])
	}

	cityUpdate, ok := updates["city_id"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected city_id update, got: %#v", updates["city_id"])
	}

	if cityUpdate["value"] != "from-update-context" {
		t.Fatalf("expected city_id.value to be from-update-context, got %#v", cityUpdate["value"])
	}
}

func TestHandleResolveDependencies_RelationFieldCascade(t *testing.T) {
	assigneeField := fields.BelongsTo("Assignee", "assignee_id", "users")
	assigneeField.DependsOn("team_id")
	assigneeField.OnDependencyChange(func(
		field *fields.Schema,
		formData map[string]interface{},
		ctx *fiber.Ctx,
	) *fields.FieldUpdate {
		return fields.NewFieldUpdate().SetValue("assignee-reset")
	})

	roleField := fields.Text("Role", "role_id")
	roleField.DependsOn("assignee_id")
	roleField.OnDependencyChange(func(
		field *fields.Schema,
		formData map[string]interface{},
		ctx *fiber.Ctx,
	) *fields.FieldUpdate {
		return fields.NewFieldUpdate().SetValue("role-reset")
	})

	h := newRelationTestHandler([]fields.Element{
		fields.Text("Team", "team_id"),
		assigneeField,
		roleField,
	})

	statusCode, response := runResolveDependenciesRequest(t, h, map[string]interface{}{
		"formData": map[string]interface{}{
			"team_id": "99",
		},
		"context":       "update",
		"changedFields": []string{"team_id"},
	})

	if statusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, statusCode)
	}

	updates, ok := response["fields"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected response.fields to be map, got %T", response["fields"])
	}

	assigneeUpdate, ok := updates["assignee_id"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected assignee_id update, got: %#v", updates["assignee_id"])
	}
	if assigneeUpdate["value"] != "assignee-reset" {
		t.Fatalf("expected assignee_id.value to be assignee-reset, got %#v", assigneeUpdate["value"])
	}

	roleUpdate, ok := updates["role_id"].(map[string]interface{})
	if !ok {
		t.Fatalf("expected role_id update, got: %#v", updates["role_id"])
	}
	if roleUpdate["value"] != "role-reset" {
		t.Fatalf("expected role_id.value to be role-reset, got %#v", roleUpdate["value"])
	}
}

func TestParseBody_RelationSentinelConvertedToNil(t *testing.T) {
	h := newRelationTestHandler([]fields.Element{
		fields.BelongsTo("Category", "category_id", "categories"),
	})

	parsed := runParseBodyRequest(
		t,
		h,
		"application/json",
		bytes.NewBufferString(`{"category_id":"__PANEL_NULL__"}`),
	)

	value, exists := parsed["category_id"]
	if !exists {
		t.Fatalf("expected category_id key in parsed body")
	}
	if value != nil {
		t.Fatalf("expected category_id to be nil, got %#v", value)
	}
}

func TestParseBody_MorphToSentinelClearsTypeAndID(t *testing.T) {
	h := newRelationTestHandler([]fields.Element{
		fields.NewMorphTo("Commentable", "commentable"),
	})

	parsed := runParseBodyRequest(
		t,
		h,
		"application/json",
		bytes.NewBufferString(`{"commentable":"__PANEL_NULL__"}`),
	)

	if _, exists := parsed["commentable"]; exists {
		t.Fatalf("expected commentable composite key to be removed")
	}

	typeValue, typeExists := parsed["commentable_type"]
	if !typeExists {
		t.Fatalf("expected commentable_type key in parsed body")
	}
	if typeValue != nil {
		t.Fatalf("expected commentable_type to be nil, got %#v", typeValue)
	}

	idValue, idExists := parsed["commentable_id"]
	if !idExists {
		t.Fatalf("expected commentable_id key in parsed body")
	}
	if idValue != nil {
		t.Fatalf("expected commentable_id to be nil, got %#v", idValue)
	}
}

func TestParseBody_MultipartMissingBelongsToManyBecomesEmptySlice(t *testing.T) {
	h := newRelationTestHandler([]fields.Element{
		fields.Text("Title", "title"),
		fields.BelongsToMany("Tags", "tags", "tags"),
	})

	var formBody bytes.Buffer
	writer := multipart.NewWriter(&formBody)
	if err := writer.WriteField("title", "Smoke Test"); err != nil {
		t.Fatalf("failed to write title field: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("failed to close multipart writer: %v", err)
	}

	parsed := runParseBodyRequest(
		t,
		h,
		writer.FormDataContentType(),
		&formBody,
	)

	if parsed["title"] != "Smoke Test" {
		t.Fatalf("expected title to be Smoke Test, got %#v", parsed["title"])
	}

	tagsRaw, exists := parsed["tags"]
	if !exists {
		t.Fatalf("expected tags key in parsed body")
	}

	tags, ok := tagsRaw.([]interface{})
	if !ok {
		t.Fatalf("expected tags to be []interface{}, got %T", tagsRaw)
	}

	if len(tags) != 0 {
		t.Fatalf("expected tags to be empty slice, got %#v", tags)
	}
}
