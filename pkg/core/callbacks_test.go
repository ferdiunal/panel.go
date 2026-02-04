package core_test

import (
	"mime/multipart"
	"testing"

	"github.com/ferdiunal/panel.go/pkg/core"
	"github.com/gofiber/fiber/v2"
)

// TestVisibilityFuncType verifies VisibilityFunc type signature
func TestVisibilityFuncType(t *testing.T) {
	// Create a sample visibility function
	var visibilityFunc core.VisibilityFunc = func(ctx *core.ResourceContext) bool {
		return ctx != nil && ctx.Resource != nil
	}

	// Test with nil context
	if visibilityFunc(nil) {
		t.Error("VisibilityFunc should return false for nil context")
	}

	// Test with valid context
	ctx := &core.ResourceContext{
		Resource: "test",
		Elements: []core.Element{},
		Request:  nil,
	}
	if !visibilityFunc(ctx) {
		t.Error("VisibilityFunc should return true for valid context with resource")
	}
}

// TestStorageCallbackFuncType verifies StorageCallbackFunc type signature
func TestStorageCallbackFuncType(t *testing.T) {
	// Create a sample storage callback function
	var storageFunc core.StorageCallbackFunc = func(c *fiber.Ctx, file *multipart.FileHeader) (string, error) {
		if file == nil {
			return "", fiber.NewError(fiber.StatusBadRequest, "no file provided")
		}
		// Simulate successful storage
		return "/uploads/" + file.Filename, nil
	}

	// Test with nil file
	path, err := storageFunc(nil, nil)
	if err == nil {
		t.Error("StorageCallbackFunc should return error for nil file")
	}
	if path != "" {
		t.Error("StorageCallbackFunc should return empty path on error")
	}

	// Test with valid file
	mockFile := &multipart.FileHeader{
		Filename: "test.txt",
		Size:     100,
	}
	path, err = storageFunc(nil, mockFile)
	if err != nil {
		t.Errorf("StorageCallbackFunc should not return error for valid file: %v", err)
	}
	expectedPath := "/uploads/test.txt"
	if path != expectedPath {
		t.Errorf("StorageCallbackFunc returned wrong path: got %s, want %s", path, expectedPath)
	}
}

// TestCallbackTypesAreExported verifies callback types are exported
func TestCallbackTypesAreExported(t *testing.T) {
	// This test verifies that the callback types can be used outside the core package
	// If this compiles, the types are properly exported

	var _ core.VisibilityFunc
	var _ core.StorageCallbackFunc
}
