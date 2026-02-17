package core_test

import (
	"testing"

	"github.com/ferdiunal/panel.go/pkg/core"
)

// TestElementInterfaceExists verifies that the Element interface is properly defined
func TestElementInterfaceExists(t *testing.T) {
	// This test ensures the Element interface exists and can be referenced
	var _ core.Element = nil // This will compile if the interface exists
}

// TestElementInterfaceLocation verifies that Element interface is in pkg/core
func TestElementInterfaceLocation(t *testing.T) {
	// Verify we can import and reference the Element interface from core package
	var element core.Element
	if element != nil {
		t.Error("Element should be nil by default")
	}
}

// TestElementInterfaceHasRequiredMethods verifies the interface has all required methods
// This is a compile-time check - if any method is missing, this won't compile
func TestElementInterfaceHasRequiredMethods(t *testing.T) {
	// Create a mock implementation to verify all methods exist
	var element core.Element

	// If this compiles, all methods are present in the interface
	// We're not testing implementation, just interface completeness
	if element != nil {
		// Basic Accessors
		_ = element.GetKey()
		_ = element.GetView()
		_ = element.GetContext()

		// Data Processing
		element.Extract(nil)
		_ = element.JsonSerialize()

		// Visibility Control
		_ = element.IsVisible(nil)
		_ = element.IsSearchable()

		// Fluent Setters - View Control
		_ = element.SetName("")
		_ = element.SetKey("")
		_ = element.OnList()
		_ = element.OnDetail()
		_ = element.OnForm()
		_ = element.HideOnList()
		_ = element.HideOnDetail()
		_ = element.HideOnCreate()
		_ = element.HideOnUpdate()
		_ = element.HideOnApi()
		_ = element.OnlyOnList()
		_ = element.OnlyOnDetail()
		_ = element.OnlyOnCreate()
		_ = element.OnlyOnUpdate()
		_ = element.OnlyOnForm()

		// Fluent Setters - Properties
		_ = element.ReadOnly()
		_ = element.WithProps("", nil)
		_ = element.Disabled()
		_ = element.Immutable()
		_ = element.Required()
		_ = element.Nullable()
		_ = element.Placeholder("")
		_ = element.Label("")
		_ = element.HelpText("")
		_ = element.Filterable()
		_ = element.Sortable()
		_ = element.Searchable()
		_ = element.Stacked()
		_ = element.SetTextAlign("")

		// Callbacks
		_ = element.CanSee(nil)
		_ = element.StoreAs(nil)
		_ = element.GetStorageCallback()
		_ = element.Resolve(nil)
		_ = element.GetResolveCallback()
		_ = element.Modify(nil)
		_ = element.GetModifyCallback()

		// Display
		_ = element.Display(nil)
		_ = element.DisplayAs("")
		_ = element.DisplayUsingLabels()

		// Other
		_ = element.Options(nil)
		_ = element.Default(nil)
	}
}
