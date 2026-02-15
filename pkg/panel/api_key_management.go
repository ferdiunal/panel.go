package panel

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/domain/apikey"
	"github.com/ferdiunal/panel.go/pkg/middleware"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type createManagedAPIKeyRequest struct {
	Name         string  `json:"name"`
	ExpiresAt    *string `json:"expires_at"`
	ExpiresAtAlt *string `json:"expiresAt"`
}

type managedAPIKeyResponse struct {
	ID              uint       `json:"id"`
	Name            string     `json:"name"`
	Prefix          string     `json:"prefix"`
	CreatedByUserID *uint      `json:"created_by_user_id,omitempty"`
	LastUsedAt      *time.Time `json:"last_used_at,omitempty"`
	ExpiresAt       *time.Time `json:"expires_at,omitempty"`
	RevokedAt       *time.Time `json:"revoked_at,omitempty"`
	Status          string     `json:"status"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

func (p *Panel) handleAPIKeyList(c *context.Context) error {
	if err := requireAPIKeyAdminSession(c); err != nil {
		return err
	}

	var keys []apikey.APIKey
	if err := p.Db.Order("id DESC").Find(&keys).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to list API keys",
		})
	}

	now := time.Now().UTC()
	items := make([]managedAPIKeyResponse, 0, len(keys))
	for _, key := range keys {
		items = append(items, toManagedAPIKeyResponse(key, now))
	}

	return c.JSON(fiber.Map{
		"data": items,
	})
}

func (p *Panel) handleAPIKeyCreate(c *context.Context) error {
	if err := requireAPIKeyAdminSession(c); err != nil {
		return err
	}

	var req createManagedAPIKeyRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON",
		})
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "name is required",
		})
	}

	expiresAt, err := parseManagedAPIKeyExpiry(req.ExpiresAt, req.ExpiresAtAlt)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	rawKey, prefix, keyHash, err := generateManagedAPIKey()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate API key",
		})
	}

	var createdByUserID *uint
	if currentUser := c.User(); currentUser != nil {
		id := currentUser.ID
		createdByUserID = &id
	}

	record := apikey.APIKey{
		Name:            name,
		Prefix:          prefix,
		KeyHash:         keyHash,
		CreatedByUserID: createdByUserID,
		ExpiresAt:       expiresAt,
	}

	if err := p.Db.Create(&record).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create API key",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": toManagedAPIKeyResponse(record, time.Now().UTC()),
		"key":  rawKey,
	})
}

func (p *Panel) handleAPIKeyRevoke(c *context.Context) error {
	if err := requireAPIKeyAdminSession(c); err != nil {
		return err
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil || id == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid api key id",
		})
	}

	var record apikey.APIKey
	if err := p.Db.First(&record, uint(id)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "API key not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to revoke API key",
		})
	}

	now := time.Now().UTC()
	if record.RevokedAt == nil {
		if err := p.Db.Model(&record).Update("revoked_at", now).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to revoke API key",
			})
		}
		record.RevokedAt = &now
	}

	return c.JSON(fiber.Map{
		"data": toManagedAPIKeyResponse(record, now),
	})
}

func (p *Panel) validateManagedAPIKey(c *fiber.Ctx, incoming string) bool {
	if p == nil || p.Db == nil {
		return false
	}

	incoming = strings.TrimSpace(incoming)
	if incoming == "" {
		return false
	}

	hash := sha256.Sum256([]byte(incoming))
	keyHash := hex.EncodeToString(hash[:])
	now := time.Now().UTC()

	var record apikey.APIKey
	err := p.Db.
		Where("key_hash = ?", keyHash).
		Where("revoked_at IS NULL").
		Where("(expires_at IS NULL OR expires_at > ?)", now).
		First(&record).Error
	if err != nil {
		return false
	}

	_ = p.Db.Model(&apikey.APIKey{}).
		Where("id = ?", record.ID).
		Update("last_used_at", now).Error

	c.Locals("api_key_id", record.ID)
	return true
}

func requireAPIKeyAdminSession(c *context.Context) error {
	if c == nil {
		return fiber.ErrUnauthorized
	}

	if apiKeyAuth, ok := c.Locals(middleware.APIKeyAuthenticatedLocalKey).(bool); ok && apiKeyAuth {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "API key management requires session authentication",
		})
	}

	currentUser := c.User()
	if currentUser == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	if !strings.EqualFold(currentUser.Role, "admin") {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied",
		})
	}

	return nil
}

func toManagedAPIKeyResponse(record apikey.APIKey, now time.Time) managedAPIKeyResponse {
	return managedAPIKeyResponse{
		ID:              record.ID,
		Name:            record.Name,
		Prefix:          record.Prefix,
		CreatedByUserID: record.CreatedByUserID,
		LastUsedAt:      record.LastUsedAt,
		ExpiresAt:       record.ExpiresAt,
		RevokedAt:       record.RevokedAt,
		Status:          managedAPIKeyStatus(record, now),
		CreatedAt:       record.CreatedAt,
		UpdatedAt:       record.UpdatedAt,
	}
}

func managedAPIKeyStatus(record apikey.APIKey, now time.Time) string {
	if record.RevokedAt != nil {
		return "revoked"
	}
	if record.ExpiresAt != nil && !record.ExpiresAt.After(now) {
		return "expired"
	}
	return "active"
}

func parseManagedAPIKeyExpiry(primary, fallback *string) (*time.Time, error) {
	raw := ""
	if primary != nil {
		raw = strings.TrimSpace(*primary)
	}
	if raw == "" && fallback != nil {
		raw = strings.TrimSpace(*fallback)
	}
	if raw == "" {
		return nil, nil
	}

	expiresAt, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return nil, fmt.Errorf("expires_at must be RFC3339")
	}

	expiresAt = expiresAt.UTC()
	return &expiresAt, nil
}

func generateManagedAPIKey() (rawKey, prefix, keyHash string, err error) {
	buf := make([]byte, 32)
	if _, err = rand.Read(buf); err != nil {
		return "", "", "", err
	}

	rawKey = "pnl_" + base64.RawURLEncoding.EncodeToString(buf)
	prefix = rawKey
	if len(prefix) > 12 {
		prefix = prefix[:12]
	}

	sum := sha256.Sum256([]byte(rawKey))
	keyHash = hex.EncodeToString(sum[:])
	return rawKey, prefix, keyHash, nil
}
