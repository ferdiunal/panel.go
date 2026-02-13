package audit

import "time"

type Log struct {
	ID         string                 `json:"id" gorm:"primaryKey;type:uuid"`
	UserID     string                 `json:"user_id,omitempty" gorm:"index;type:uuid"`
	SessionID  string                 `json:"session_id,omitempty" gorm:"index;type:uuid"`
	Action     string                 `json:"action" gorm:"index"`
	Resource   string                 `json:"resource" gorm:"index"`
	ResourceID string                 `json:"resource_id,omitempty" gorm:"index"`
	Method     string                 `json:"method" gorm:"index"`
	Path       string                 `json:"path" gorm:"index"`
	StatusCode int                    `json:"status_code" gorm:"index"`
	IPAddress  string                 `json:"ip_address"`
	UserAgent  string                 `json:"user_agent"`
	RequestID  string                 `json:"request_id" gorm:"index"`
	Metadata   map[string]interface{} `json:"metadata" gorm:"serializer:json"`
	CreatedAt  time.Time              `json:"created_at" gorm:"index"`
}

func (Log) TableName() string {
	return "audit_logs"
}
