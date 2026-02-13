package audit

import (
	"github.com/ferdiunal/panel.go/internal/context"
	"github.com/ferdiunal/panel.go/internal/domain/audit"
	"github.com/ferdiunal/panel.go/internal/domain/user"
	"github.com/ferdiunal/panel.go/internal/fields"
	"github.com/ferdiunal/panel.go/internal/resource"
	"gorm.io/gorm"
)

// AuditLogPolicy, audit log'lar için yetkilendirme politikasıdır.
// Sadece admin kullanıcılar audit log'ları görüntüleyebilir.
type AuditLogPolicy struct{}

// ViewAny, kullanıcının audit log listesini görüntüleyip görüntüleyemeyeceğini belirler.
func (p AuditLogPolicy) ViewAny(ctx *context.Context) bool {
	authUser := ctx.User()
	return authUser != nil && authUser.Role == user.RoleAdmin
}

// View, kullanıcının belirli bir audit log'u görüntüleyip görüntüleyemeyeceğini belirler.
func (p AuditLogPolicy) View(ctx *context.Context, model interface{}) bool {
	authUser := ctx.User()
	return authUser != nil && authUser.Role == user.RoleAdmin
}

// Create, audit log oluşturma yetkisini kontrol eder.
// Audit log'lar otomatik oluşturulur, manuel oluşturma yasaktır.
func (p AuditLogPolicy) Create(ctx *context.Context) bool {
	return false
}

// Update, audit log güncelleme yetkisini kontrol eder.
// Audit log'lar immutable'dır, güncellenemez.
func (p AuditLogPolicy) Update(ctx *context.Context, model interface{}) bool {
	return false
}

// Delete, audit log silme yetkisini kontrol eder.
// Audit log'lar retention policy ile otomatik silinir, manuel silme yasaktır.
func (p AuditLogPolicy) Delete(ctx *context.Context, model interface{}) bool {
	return false
}

// GetAuditLogResource, audit log resource'unu döndürür.
func GetAuditLogResource() resource.Resource {
	return resource.NewResource(
		"audit-logs",
		&audit.Log{},
		func() []fields.Element {
			return []fields.Element{
				fields.Text("ID", "id").
					HideOnIndex().
					HideOnDetail().
					HideOnForm(),

				fields.BelongsTo("User", "user_id", "users").
					DisplayUsing("name").
					Searchable().
					Sortable(),

				fields.Text("Session ID", "session_id").
					HideOnIndex().
					Copyable(),

				fields.Badge("Action", "action").
					Options(map[string]string{
						"create": "success",
						"update": "warning",
						"delete": "destructive",
					}).
					Searchable().
					Sortable(),

				fields.Text("Resource", "resource").
					Searchable().
					Sortable(),

				fields.Text("Resource ID", "resource_id").
					HideOnIndex().
					Copyable(),

				fields.Badge("Method", "method").
					Options(map[string]string{
						"POST":   "default",
						"PUT":    "secondary",
						"PATCH":  "secondary",
						"DELETE": "destructive",
					}).
					Sortable(),

				fields.Text("Path", "path").
					HideOnIndex(),

				fields.Badge("Status", "status_code").
					Options(map[string]string{
						"200": "success",
						"201": "success",
						"204": "success",
						"400": "warning",
						"401": "warning",
						"403": "warning",
						"404": "warning",
						"500": "destructive",
					}).
					Sortable(),

				fields.Text("IP Address", "ip_address").
					Searchable().
					Copyable(),

				fields.Text("User Agent", "user_agent").
					HideOnIndex(),

				fields.Text("Request ID", "request_id").
					HideOnIndex().
					Copyable(),

				fields.Code("Metadata", "metadata").
					Language("json").
					HideOnIndex().
					OnlyOnDetail(),

				fields.DateTime("Created At", "created_at").
					Format("2006-01-02 15:04:05").
					Sortable().
					SortByDesc(),
			}
		},
	).
		SetTitle("Audit Logs").
		SetIcon("shield-check").
		SetGroup("System").
		SetPolicy(AuditLogPolicy{}).
		SetSearchColumns([]string{"action", "resource", "ip_address", "user_agent"}).
		SetPerPage(50).
		SetNavigationOrder(1000) // En sonda göster
}
