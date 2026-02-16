package resource

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	appContext "github.com/ferdiunal/panel.go/pkg/context"
	"github.com/gofiber/fiber/v2"
)

func TestOptimizedBaseStoreHandler_DefaultStorage(t *testing.T) {
	tmpDir := t.TempDir()
	r := &OptimizedBase{}

	app := fiber.New()
	app.Post("/upload", appContext.Wrap(func(c *appContext.Context) error {
		file, err := c.Ctx.FormFile("image")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		storedPath, err := r.StoreHandler(c, file, tmpDir, "/storage")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"path": storedPath,
		})
	}))

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("image", "hero.png")
	if err != nil {
		t.Fatalf("failed to create form file: %v", err)
	}
	if _, err := part.Write([]byte("fake-image-content")); err != nil {
		t.Fatalf("failed to write fake file content: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("failed to close multipart writer: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/upload", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d, body: %s", http.StatusOK, resp.StatusCode, string(respBody))
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(respBody, &payload); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	storedPath, ok := payload["path"].(string)
	if !ok || strings.TrimSpace(storedPath) == "" {
		t.Fatalf("expected non-empty stored path, got %#v", payload["path"])
	}

	if !strings.HasPrefix(storedPath, "/storage/") {
		t.Fatalf("expected stored path to start with /storage/, got %q", storedPath)
	}

	fileName := strings.TrimPrefix(storedPath, "/storage/")
	if fileName == "" || fileName == storedPath {
		t.Fatalf("expected filename in stored path, got %q", storedPath)
	}

	savedFilePath := filepath.Join(tmpDir, fileName)
	if _, err := os.Stat(savedFilePath); err != nil {
		t.Fatalf("expected saved file at %q, stat error: %v", savedFilePath, err)
	}
}
