package observability

import (
	"strings"
	"time"

	appContext "github.com/ferdiunal/panel.go/internal/context"
	"github.com/ferdiunal/panel.go/internal/domain/audit"
	"github.com/ferdiunal/panel.go/shared/uuid"
	"gorm.io/gorm"
)

func AuditMiddleware(db *gorm.DB) appContext.Handler {
	return func(c *appContext.Context) error {
		err := c.Next()

		if db == nil || !shouldAudit(c.Method(), c.Path()) {
			return err
		}

		resource, resourceID := extractAuditTarget(c.Path())

		entry := &audit.Log{
			ID:         uuid.NewUUID().String(),
			Action:     actionFromMethod(c.Method()),
			Resource:   resource,
			ResourceID: resourceID,
			Method:     c.Method(),
			Path:       c.Path(),
			StatusCode: c.Response().StatusCode(),
			IPAddress:  c.IP(),
			UserAgent:  c.Get("User-Agent"),
			RequestID:  c.GetRespHeader("X-Request-ID"),
			CreatedAt:  time.Now(),
		}

		if user := c.User(); user != nil {
			entry.UserID = user.ID
		}
		if session := c.Session(); session != nil {
			entry.SessionID = session.ID
		}

		entry.Metadata = map[string]interface{}{
			"query": string(c.Request().URI().QueryString()),
		}

		_ = db.WithContext(c.Context()).Create(entry).Error
		return err
	}
}

func shouldAudit(method, path string) bool {
	switch method {
	case "GET", "HEAD", "OPTIONS":
		return false
	}
	return strings.HasPrefix(path, "/api/")
}

func actionFromMethod(method string) string {
	switch method {
	case "POST":
		return "create"
	case "PUT", "PATCH":
		return "update"
	case "DELETE":
		return "delete"
	default:
		return "unknown"
	}
}

func extractAuditTarget(path string) (string, string) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 2 || parts[0] != "api" {
		return "api", ""
	}

	switch parts[1] {
	case "resource":
		if len(parts) >= 3 {
			resource := parts[2]
			resourceID := ""
			if len(parts) >= 4 {
				resourceID = parts[3]
			}
			return "resource:" + resource, resourceID
		}
		return "resource", ""
	case "pages":
		if len(parts) >= 3 {
			return "page:" + parts[2], parts[2]
		}
		return "page", ""
	case "auth":
		if len(parts) >= 3 {
			return "auth:" + parts[2], ""
		}
		return "auth", ""
	default:
		return parts[1], ""
	}
}
