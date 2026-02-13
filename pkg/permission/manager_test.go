package permission_test

import (
	"os"
	"testing"

	"github.com/ferdiunal/panel.go/pkg/permission"
	"github.com/stretchr/testify/assert"
)

func TestManager_HasPermission(t *testing.T) {
	// Create a temporary permissions.toml file
	content := `
[admin]
label = "Admin"
permissions = ["*"]

[editor]
label = "Editor"
permissions = ["posts.create", "posts.edit"]

[user]
label = "User"
permissions = ["comments.create"]
`
	tmpfile, err := os.CreateTemp("", "permissions-*.toml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Load the permissions
	mgr, err := permission.Load(tmpfile.Name())
	assert.NoError(t, err)
	assert.NotNil(t, mgr)

	tests := []struct {
		role       string
		permission string
		want       bool
	}{
		{"admin", "anything", true},
		{"admin", "posts.create", true},
		{"editor", "posts.create", true},
		{"editor", "posts.delete", false},
		{"user", "comments.create", true},
		{"user", "posts.create", false},
		{"guest", "anything", false},
	}

	for _, tt := range tests {
		t.Run(tt.role+":"+tt.permission, func(t *testing.T) {
			got := mgr.HasPermission(tt.role, tt.permission)
			assert.Equal(t, tt.want, got)
		})
	}
}
