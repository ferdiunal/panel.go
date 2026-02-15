package apikey

import "time"

// APIKey stores a managed API key record.
// Raw token is never persisted; only sha256 hash is stored.
type APIKey struct {
	ID              uint       `json:"id" gorm:"primaryKey"`
	Name            string     `json:"name" gorm:"index"`
	Prefix          string     `json:"prefix" gorm:"index"`
	KeyHash         string     `json:"-" gorm:"uniqueIndex;size:64"`
	CreatedByUserID *uint      `json:"created_by_user_id,omitempty" gorm:"index"`
	LastUsedAt      *time.Time `json:"last_used_at,omitempty" gorm:"index"`
	ExpiresAt       *time.Time `json:"expires_at,omitempty" gorm:"index"`
	RevokedAt       *time.Time `json:"revoked_at,omitempty" gorm:"index"`
	CreatedAt       time.Time  `json:"created_at" gorm:"index"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"index"`
}

func (k *APIKey) IsActive(now time.Time) bool {
	if k == nil {
		return false
	}
	if k.RevokedAt != nil {
		return false
	}
	if k.ExpiresAt != nil && k.ExpiresAt.Before(now) {
		return false
	}
	return true
}

