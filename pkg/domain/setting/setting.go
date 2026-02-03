package setting

import (
	"time"
)

// Setting represents a key-value pair setting.
type Setting struct {
	Key       string                 `gorm:"primaryKey;type:varchar(255)" json:"key"`
	Value     map[string]interface{} `gorm:"type:jsonb;serializer:json" json:"value"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

// TableName overrides the table name used by User to `settings`.
func (Setting) TableName() string {
	return "settings"
}
