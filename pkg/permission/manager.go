package permission

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

// Role defines the structure for a role in permissions.toml
type Role struct {
	Label       string   `toml:"label"`
	Permissions []string `toml:"permissions"`
}

// Config represents the structure of the permissions.toml file.
// Since the file uses top-level keys for roles (e.g. [admin], [user]),
// we parse it as a map.
type Config map[string]Role

// Manager handles permission configurations.
type Manager struct {
	config Config
}

// Global instance
var currentManager *Manager

// Load reads and parses the toml file at the given path.
func Load(path string) (*Manager, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("permissions file could not be read: %w", err)
	}

	var config Config
	if err := toml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("permissions file could not be parsed: %w", err)
	}

	mgr := &Manager{
		config: config,
	}
	currentManager = mgr
	return mgr, nil
}

// GetConfig returns the loaded configuration.
func (m *Manager) GetConfig() Config {
	if m.config == nil {
		return make(Config)
	}
	return m.config
}

// GetRoles returns the list of role keys defined in the configuration.
func (m *Manager) GetRoles() []string {
	roles := make([]string, 0, len(m.config))
	for r := range m.config {
		roles = append(roles, r)
	}
	return roles
}

// GetInstance returns the current manager instance.
func GetInstance() *Manager {
	return currentManager
}

// GetRole returns the Role definition for a specific role key.
func (m *Manager) GetRole(roleKey string) (Role, bool) {
	role, ok := m.config[roleKey]
	return role, ok
}

// HasPermission checks if a role has a specific permission.
func (m *Manager) HasPermission(roleName string, permission string) bool {
	role, ok := m.config[roleName]
	if !ok {
		return false
	}

	for _, p := range role.Permissions {
		if p == "*" || p == permission {
			return true
		}
	}

	return false
}
