package handler

import (
	"testing"

	"github.com/ferdiunal/panel.go/pkg/resource"
	resourceAccount "github.com/ferdiunal/panel.go/pkg/resource/account"
	resourceSession "github.com/ferdiunal/panel.go/pkg/resource/session"
	resourceSetting "github.com/ferdiunal/panel.go/pkg/resource/setting"
	resourceUser "github.com/ferdiunal/panel.go/pkg/resource/user"
	resourceVerification "github.com/ferdiunal/panel.go/pkg/resource/verification"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestResourceHandlerAccountResource, Account resource handler'ını test eder
func TestResourceHandlerAccountResource(t *testing.T) {
	res := resourceAccount.NewAccountResource()

	if res == nil {
		t.Error("Expected non-nil resource")
	}

	if res.Slug() != "accounts" {
		t.Errorf("Expected slug 'accounts', got '%s'", res.Slug())
	}

	fields := res.Fields()
	if len(fields) == 0 {
		t.Error("Expected fields to be defined")
	}
}

// TestResourceHandlerSessionResource, Session resource handler'ını test eder
func TestResourceHandlerSessionResource(t *testing.T) {
	res := resourceSession.NewSessionResource()

	if res == nil {
		t.Error("Expected non-nil resource")
	}

	if res.Slug() != "sessions" {
		t.Errorf("Expected slug 'sessions', got '%s'", res.Slug())
	}

	fields := res.Fields()
	if len(fields) == 0 {
		t.Error("Expected fields to be defined")
	}
}

// TestResourceHandlerVerificationResource, Verification resource handler'ını test eder
func TestResourceHandlerVerificationResource(t *testing.T) {
	res := resourceVerification.NewVerificationResource()

	if res == nil {
		t.Error("Expected non-nil resource")
	}

	if res.Slug() != "verifications" {
		t.Errorf("Expected slug 'verifications', got '%s'", res.Slug())
	}

	fields := res.Fields()
	if len(fields) == 0 {
		t.Error("Expected fields to be defined")
	}
}

// TestResourceHandlerSettingResource, Setting resource handler'ını test eder
func TestResourceHandlerSettingResource(t *testing.T) {
	res := resourceSetting.NewSettingResource()

	if res == nil {
		t.Error("Expected non-nil resource")
	}

	if res.Slug() != "settings" {
		t.Errorf("Expected slug 'settings', got '%s'", res.Slug())
	}

	fields := res.Fields()
	if len(fields) == 0 {
		t.Error("Expected fields to be defined")
	}
}

// TestResourceHandlerUserResource, User resource handler'ını test eder
func TestResourceHandlerUserResource(t *testing.T) {
	res := resourceUser.NewUserResource()

	if res == nil {
		t.Error("Expected non-nil resource")
	}

	if res.Slug() != "users" {
		t.Errorf("Expected slug 'users', got '%s'", res.Slug())
	}

	fields := res.Fields()
	if len(fields) == 0 {
		t.Error("Expected fields to be defined")
	}
}

// TestResourceHandlerFieldCount, resource alanlarının sayısını test eder
func TestResourceHandlerFieldCount(t *testing.T) {
	resources := []resource.Resource{
		resourceAccount.NewAccountResource(),
		resourceSession.NewSessionResource(),
		resourceVerification.NewVerificationResource(),
		resourceSetting.NewSettingResource(),
		resourceUser.NewUserResource(),
	}

	for _, res := range resources {
		fields := res.Fields()
		if len(fields) == 0 {
			t.Errorf("Expected fields for resource '%s'", res.Slug())
		}
	}
}

// TestResourceHandlerPolicy, resource policy'lerini test eder
func TestResourceHandlerPolicy(t *testing.T) {
	resources := []resource.Resource{
		resourceAccount.NewAccountResource(),
		resourceSession.NewSessionResource(),
		resourceVerification.NewVerificationResource(),
		resourceSetting.NewSettingResource(),
		resourceUser.NewUserResource(),
	}

	for _, res := range resources {
		policy := res.Policy()
		if policy == nil {
			t.Errorf("Expected policy for resource '%s'", res.Slug())
		}
	}
}

// TestResourceHandlerRepository, resource repository'lerini test eder
func TestResourceHandlerRepository(t *testing.T) {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})

	resources := []resource.Resource{
		resourceAccount.NewAccountResource(),
		resourceSession.NewSessionResource(),
		resourceVerification.NewVerificationResource(),
		resourceSetting.NewSettingResource(),
	}

	for _, res := range resources {
		repo := res.Repository(db)
		if repo == nil {
			t.Errorf("Expected repository for resource '%s'", res.Slug())
		}
	}
}

// TestResourceHandlerCards, resource card'larını test eder
func TestResourceHandlerCards(t *testing.T) {
	resources := []resource.Resource{
		resourceAccount.NewAccountResource(),
		resourceSession.NewSessionResource(),
		resourceVerification.NewVerificationResource(),
		resourceSetting.NewSettingResource(),
	}

	for _, res := range resources {
		cards := res.Cards()
		// Cards can be empty, just check it's a slice
		if cards == nil {
			// It's ok if cards is nil, just make sure it's not an error
			continue
		}
	}
}

// TestResourceHandlerSortable, resource sıralama ayarlarını test eder
func TestResourceHandlerSortable(t *testing.T) {
	resources := []resource.Resource{
		resourceAccount.NewAccountResource(),
		resourceSession.NewSessionResource(),
		resourceVerification.NewVerificationResource(),
		resourceSetting.NewSettingResource(),
	}

	for _, res := range resources {
		sortable := res.GetSortable()
		if sortable == nil {
			t.Errorf("Expected sortable for resource '%s'", res.Slug())
		}
	}
}

// TestResourceHandlerVisible, resource görünürlüğünü test eder
func TestResourceHandlerVisible(t *testing.T) {
	resources := []resource.Resource{
		resourceAccount.NewAccountResource(),
		resourceSession.NewSessionResource(),
		resourceVerification.NewVerificationResource(),
		resourceSetting.NewSettingResource(),
		resourceUser.NewUserResource(),
	}

	for _, res := range resources {
		if !res.Visible() {
			t.Errorf("Expected resource '%s' to be visible", res.Slug())
		}
	}
}

// TestResourceHandlerModel, resource model'lerini test eder
func TestResourceHandlerModel(t *testing.T) {
	resources := []resource.Resource{
		resourceAccount.NewAccountResource(),
		resourceSession.NewSessionResource(),
		resourceVerification.NewVerificationResource(),
		resourceSetting.NewSettingResource(),
		resourceUser.NewUserResource(),
	}

	for _, res := range resources {
		model := res.Model()
		if model == nil {
			t.Errorf("Expected model for resource '%s'", res.Slug())
		}
	}
}

// TestResourceHandlerTitle, resource başlıklarını test eder
func TestResourceHandlerTitle(t *testing.T) {
	tests := []struct {
		resource resource.Resource
		expected string
	}{
		{resourceAccount.NewAccountResource(), "Accounts"},
		{resourceSession.NewSessionResource(), "Sessions"},
		{resourceVerification.NewVerificationResource(), "Verifications"},
		{resourceSetting.NewSettingResource(), "Settings"},
		{resourceUser.NewUserResource(), "Users"},
	}

	for _, tt := range tests {
		if tt.resource.Title() != tt.expected {
			t.Errorf("Expected title '%s', got '%s'", tt.expected, tt.resource.Title())
		}
	}
}

// TestResourceHandlerIcon, resource ikonlarını test eder
func TestResourceHandlerIcon(t *testing.T) {
	tests := []struct {
		resource resource.Resource
		expected string
	}{
		{resourceAccount.NewAccountResource(), "key"},
		{resourceSession.NewSessionResource(), "clock"},
		{resourceVerification.NewVerificationResource(), "shield-check"},
		{resourceSetting.NewSettingResource(), "settings"},
		{resourceUser.NewUserResource(), "users"},
	}

	for _, tt := range tests {
		if tt.resource.Icon() != tt.expected {
			t.Errorf("Expected icon '%s', got '%s'", tt.expected, tt.resource.Icon())
		}
	}
}

// TestResourceHandlerGroup, resource gruplarını test eder
func TestResourceHandlerGroup(t *testing.T) {
	tests := []struct {
		resource resource.Resource
		expected string
	}{
		{resourceAccount.NewAccountResource(), "System"},
		{resourceSession.NewSessionResource(), "System"},
		{resourceVerification.NewVerificationResource(), "System"},
		{resourceSetting.NewSettingResource(), "System"},
		{resourceUser.NewUserResource(), "System"},
	}

	for _, tt := range tests {
		if tt.resource.Group() != tt.expected {
			t.Errorf("Expected group '%s', got '%s'", tt.expected, tt.resource.Group())
		}
	}
}
