package resources

import (
	"time"

	"cargo.go/entity"
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/core"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/resource"
)

// TestInvalidModel - İlişki field'ı OLMAYAN model (hata vermeli)
//
// Bu model, BelongsTo field'ı tanımlanmış ama struct'ta Organization ilişki field'ı yok.
// Panel başlatılırken açıklayıcı hata mesajı vermeli.
type TestInvalidModel struct {
	ID             uint64    `gorm:"primaryKey;autoIncrement;column:id;bigint"`
	OrganizationID uint64    `gorm:"not null;column:organization_id;bigint"`
	// ❌ Organization field'ı YOK - bu hata vermeli
	Name      string    `gorm:"not null;column:name;varchar(255)"`
	CreatedAt time.Time `gorm:"autoCreateTime;column:created_at;timestamptz"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;column:updated_at;timestamptz"`
}

type TestInvalidResource struct {
	resource.OptimizedBase
}

func NewTestInvalidResource() *TestInvalidResource {
	r := &TestInvalidResource{}
	r.SetSlug("test-invalid")
	r.SetTitle("Test Invalid")
	r.SetIcon("test")
	r.SetGroup("Testing")
	r.SetModel(&TestInvalidModel{})
	r.SetFieldResolver(&TestInvalidResolveFields{})
	r.SetVisible(false) // Test resource'u gizli
	return r
}

type TestInvalidResolveFields struct{}

func (r *TestInvalidResolveFields) ResolveFields(ctx *context.Context) []core.Element {
	return []core.Element{
		fields.ID("ID", "id"),
		fields.BelongsTo("Organization", "organization_id", "organizations"), // ❌ Bu hata vermeli
		fields.Text("Name", "name"),
		fields.Date("CreatedAt", "created_at").HideOnCreate().HideOnUpdate(),
		fields.Date("UpdatedAt", "updated_at").HideOnCreate().HideOnUpdate(),
	}
}

// TestValidModel - İlişki field'ı OLAN model (hata vermemeli)
//
// Bu model, BelongsTo field'ı tanımlanmış ve struct'ta Organization ilişki field'ı var.
// Panel başlatılırken hata vermemeli.
type TestValidModel struct {
	ID             uint64             `gorm:"primaryKey;autoIncrement;column:id;bigint"`
	OrganizationID uint64             `gorm:"not null;column:organization_id;bigint"`
	Organization   *entity.Organization `gorm:"foreignKey:OrganizationID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // ✅ İlişki field'ı var
	Name           string             `gorm:"not null;column:name;varchar(255)"`
	CreatedAt      time.Time          `gorm:"autoCreateTime;column:created_at;timestamptz"`
	UpdatedAt      time.Time          `gorm:"autoUpdateTime;column:updated_at;timestamptz"`
}

type TestValidResource struct {
	resource.OptimizedBase
}

func NewTestValidResource() *TestValidResource {
	r := &TestValidResource{}
	r.SetSlug("test-valid")
	r.SetTitle("Test Valid")
	r.SetIcon("test")
	r.SetGroup("Testing")
	r.SetModel(&TestValidModel{})
	r.SetFieldResolver(&TestValidResolveFields{})
	r.SetVisible(false) // Test resource'u gizli
	return r
}

type TestValidResolveFields struct{}

func (r *TestValidResolveFields) ResolveFields(ctx *context.Context) []core.Element {
	return []core.Element{
		fields.ID("ID", "id"),
		fields.BelongsTo("Organization", "organization_id", "organizations"), // ✅ Bu hata vermemeli
		fields.Text("Name", "name"),
		fields.Date("CreatedAt", "created_at").HideOnCreate().HideOnUpdate(),
		fields.Date("UpdatedAt", "updated_at").HideOnCreate().HideOnUpdate(),
	}
}
