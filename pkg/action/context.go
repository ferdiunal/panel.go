package action

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// ActionContext holds the context during action execution.
// It contains the selected models, field values, user information, and database connection.
type ActionContext struct {
	// Models are the selected resource instances to operate on
	Models []interface{}

	// Fields contains the action field values submitted by the user
	Fields map[string]interface{}

	// User is the authenticated user performing the action
	User interface{}

	// Resource is the resource slug (e.g., "posts", "users")
	Resource string

	// DB is the database connection
	DB *gorm.DB

	// Ctx is the Fiber HTTP context
	Ctx *fiber.Ctx
}
