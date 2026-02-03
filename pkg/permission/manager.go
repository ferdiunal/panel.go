package permission

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

// ResourcePermissions, her bir kaynak (users, roles, settings vb.) için izin tanımlarını tutar.
type ResourcePermissions struct {
	Label   string   `toml:"label"`
	Actions []string `toml:"actions"`
}

// Config, permissions.toml dosyasının yapısını temsil eder.
type Config struct {
	SystemRoles []string                       `toml:"system_roles"`
	Resources   map[string]ResourcePermissions `toml:"resources"`
}

// Manager, izinleri yöneten yapı.
type Manager struct {
	config Config
}

// Global instance
var currentManager *Manager

// Load, belirtilen yoldaki toml dosyasını okuyup parse eder.
func Load(path string) (*Manager, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("permissions dosyası okunamadı: %w", err)
	}

	var config Config
	if err := toml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("permissions dosyası parse edilemedi: %w", err)
	}

	mgr := &Manager{
		config: config,
	}
	currentManager = mgr
	return mgr, nil
}

// GetConfig returns the loaded configuration.
func (m *Manager) GetConfig() Config {
	// If resources map is nil, initialize it to avoid nil pointer issues if used elsewhere
	if m.config.Resources == nil {
		m.config.Resources = make(map[string]ResourcePermissions)
	}
	return m.config
}

// GetRoles returns the list of system roles defined in the configuration.
func (m *Manager) GetRoles() []string {
	return m.config.SystemRoles
}

// GetInstance returns the current manager instance.
func GetInstance() *Manager {
	return currentManager
}

// GetAllResources returns all resource keys.
func (m *Manager) GetAllResources() []string {
	keys := make([]string, 0, len(m.config.Resources))
	for k := range m.config.Resources {
		keys = append(keys, k)
	}
	return keys
}

// GetPermissionsForResource returns permissions for a specific resource.
func (m *Manager) GetPermissionsForResource(resource string) ([]string, bool) {
	if res, ok := m.config.Resources[resource]; ok {
		return res.Actions, true
	}
	return nil, false
}
