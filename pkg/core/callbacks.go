package core

import (
	"mime/multipart"

	"github.com/gofiber/fiber/v2"
)

// VisibilityFunc is a function type for controlling element visibility.
// It receives a ResourceContext and returns true if the element should be visible,
// false otherwise.
//
// This allows for dynamic visibility control based on the current context,
// such as user permissions, resource state, or other conditions.
type VisibilityFunc func(ctx *ResourceContext) bool

// StorageCallbackFunc is a function type for handling file storage.
// It receives the Fiber context and the uploaded file, and returns the
// storage path or URL where the file was saved, or an error if the
// storage operation failed.
//
// This allows for custom file storage implementations, such as saving to
// local disk, cloud storage (S3, GCS), or other storage backends.
type StorageCallbackFunc func(c *fiber.Ctx, file *multipart.FileHeader) (string, error)
