package setting

import (
	"context"
	"time"
)

// Setting, uygulama ayarlarını temsil eder
type Setting struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Key       string    `json:"key" gorm:"uniqueIndex"`
	Value     string    `json:"value" gorm:"type:longtext"`
	Type      string    `json:"type"` // string, integer, boolean, json
	Group     string    `json:"group" gorm:"index"`
	Label     string    `json:"label"`
	Help      string    `json:"help"`
	CreatedAt time.Time `json:"createdAt" gorm:"index"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"index"`
}

// Repository, Setting repository interface'i
type Repository interface {
	Create(ctx context.Context, setting *Setting) error
	FindByID(ctx context.Context, id string) (*Setting, error)
	FindByKey(ctx context.Context, key string) (*Setting, error)
	FindByGroup(ctx context.Context, group string) ([]Setting, error)
	Update(ctx context.Context, setting *Setting) error
	Delete(ctx context.Context, id string) error
}
