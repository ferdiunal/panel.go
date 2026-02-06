// Package core provides fundamental interfaces and types for the panel system.
// This package has no dependencies on internal packages and serves as the foundation
// for the entire architecture. All other packages should depend on core, not vice versa.
//
// The core package defines:
//   - Element interface: Common interface for form and list fields
//   - ResourceContext: Context during resource operations
//   - ElementType and ElementContext: Type definitions and constants
//   - Callback function types: For visibility and storage operations
package core

import (
	"github.com/gofiber/fiber/v2"
)

// AutoOptionsConfig holds configuration for automatic options generation
type AutoOptionsConfig struct {
	Enabled      bool
	DisplayField string
}

// Element is the common interface for fields used in form and list views.
// It provides methods for data extraction, serialization, visibility control,
// and fluent configuration of field properties.
type Element interface {
	// Basic Accessors

	// GetKey returns the unique identifier for this element.
	// The key is used to map the element to a field in the resource model.
	GetKey() string

	// GetView returns the view type of this element.
	// The view type determines how the element is rendered in the UI.
	GetView() string

	// GetContext returns the context in which this element is displayed.
	// The context indicates where the element should be shown (form, list, detail).
	GetContext() ElementContext

	// Data Processing

	// Extract extracts data from the given resource and populates the element.
	// The resource parameter is typically a struct or map containing the data to extract.
	Extract(resource any)

	// JsonSerialize serializes the element to a JSON-compatible map.
	// Returns a map containing all element properties ready for JSON encoding.
	JsonSerialize() map[string]any

	// Visibility Control

	// IsVisible determines if the element should be visible in the given context.
	// Returns true if the element should be displayed, false otherwise.
	IsVisible(ctx *ResourceContext) bool

	// IsSearchable returns whether this element can be searched.
	// Searchable elements are included in global search operations.
	IsSearchable() bool

	// Fluent Setters - View Control

	// SetName sets the display name of the element.
	// The name is shown as the field label in the UI.
	SetName(name string) Element

	// SetKey sets the unique identifier of the element.
	// The key is used to map the element to a field in the resource.
	SetKey(key string) Element

	// OnList makes the element visible on list views.
	// This is additive - it doesn't hide the element from other views.
	OnList() Element

	// OnDetail makes the element visible on detail views.
	// This is additive - it doesn't hide the element from other views.
	OnDetail() Element

	// OnForm makes the element visible on form views.
	// This is additive - it doesn't hide the element from other views.
	OnForm() Element

	// HideOnList hides the element on list views.
	// The element will still be visible on other views unless explicitly hidden.
	HideOnList() Element

	// HideOnDetail hides the element on detail views.
	// The element will still be visible on other views unless explicitly hidden.
	HideOnDetail() Element

	// HideOnCreate hides the element on create forms.
	// The element will still be visible on update forms and other views.
	HideOnCreate() Element

	// HideOnUpdate hides the element on update forms.
	// The element will still be visible on create forms and other views.
	HideOnUpdate() Element

	// OnlyOnList shows the element only on list views.
	// This hides the element from all other views.
	OnlyOnList() Element

	// OnlyOnDetail shows the element only on detail views.
	// This hides the element from all other views.
	OnlyOnDetail() Element

	// OnlyOnCreate shows the element only on create forms.
	// This hides the element from all other views.
	OnlyOnCreate() Element

	// OnlyOnUpdate shows the element only on update forms.
	// This hides the element from all other views.
	OnlyOnUpdate() Element

	// OnlyOnForm shows the element only on form views (create and update).
	// This hides the element from list and detail views.
	OnlyOnForm() Element

	// Fluent Setters - Properties

	// ReadOnly marks the element as read-only.
	// Read-only elements are displayed but cannot be modified by users.
	ReadOnly() Element

	// WithProps adds custom properties to the element.
	// Custom properties can be used to pass additional data to the frontend.
	WithProps(key string, value any) Element

	// Disabled marks the element as disabled.
	// Disabled elements are visible but not interactive.
	Disabled() Element

	// Immutable marks the element as immutable (cannot be changed after creation).
	// Immutable elements can be set during creation but not during updates.
	Immutable() Element

	// Required marks the element as required.
	// Required elements must have a value before the form can be submitted.
	Required() Element

	// Nullable marks the element as nullable.
	// Nullable elements can have null/nil values.
	Nullable() Element

	// Placeholder sets the placeholder text for the element.
	// Placeholder text is shown when the element is empty.
	Placeholder(placeholder string) Element

	// Label sets the label text for the element.
	// The label is displayed above or beside the element.
	Label(label string) Element

	// HelpText sets the help text for the element.
	// Help text provides additional guidance to users.
	HelpText(helpText string) Element

	// Filterable marks the element as filterable.
	// Filterable elements can be used to filter list views.
	Filterable() Element

	// Sortable marks the element as sortable.
	// Sortable elements can be used to sort list views.
	Sortable() Element

	// Searchable marks the element as searchable.
	// Searchable elements are included in global search operations.
	Searchable() Element

	// Stacked marks the element as stacked (full width).
	// Stacked elements take up the full width of their container.
	Stacked() Element

	// SetTextAlign sets the text alignment for the element.
	// Valid values are "left", "center", "right".
	SetTextAlign(align string) Element

	// Callbacks

	// CanSee sets a visibility callback function.
	// The callback determines whether the element should be visible based on the resource context.
	CanSee(fn VisibilityFunc) Element

	// StoreAs sets a storage callback function for file uploads.
	// The callback handles custom file storage logic and returns the stored file path.
	StoreAs(fn StorageCallbackFunc) Element

	// GetStorageCallback returns the storage callback function.
	// Returns nil if no storage callback has been set.
	GetStorageCallback() StorageCallbackFunc

	// Resolve sets a callback to resolve the value before display.
	// The callback can transform the value before it is shown to the user.
	Resolve(fn func(value any, item any, c *fiber.Ctx) any) Element

	// GetResolveCallback returns the resolve callback function.
	// Returns nil if no resolve callback has been set.
	GetResolveCallback() func(value any, item any, c *fiber.Ctx) any

	// Modify sets a callback to modify the value before storage.
	// The callback can transform the value before it is saved to the database.
	Modify(fn func(value any, c *fiber.Ctx) any) Element

	// GetModifyCallback returns the modify callback function.
	// Returns nil if no modify callback has been set.
	GetModifyCallback() func(value any, c *fiber.Ctx) any

	// Other

	// Options sets the available options for select-type elements.
	// The options parameter can be a slice of values or a map of key-value pairs.
	Options(options any) Element

	// GetAutoOptionsConfig returns the AutoOptions configuration.
	GetAutoOptionsConfig() AutoOptionsConfig

	// Default sets the default value for the element.
	// The default value is used when creating new resources.
	Default(value any) Element

	// Extended Field System Methods

	// IsHidden determines if the element should be hidden in the given visibility context.
	// Returns true if the element should be hidden, false otherwise.
	IsHidden(ctx VisibilityContext) bool

	// ResolveForDisplay resolves the element's value for display purposes.
	// This method transforms the value before it is shown to the user.
	// It differs from Resolve in that it's specifically for display formatting.
	ResolveForDisplay(item any) (any, error)

	// GetDependencies returns a list of field names that this element depends on.
	// Dependencies are used to determine field visibility and validation order.
	GetDependencies() []string

	// IsConditionallyVisible determines if the element should be visible based on the item's values.
	// This allows for dynamic visibility based on other field values.
	IsConditionallyVisible(item any) bool

	// GetMetadata returns metadata about the element.
	// Metadata can include information about field relationships, constraints, and custom properties.
	GetMetadata() map[string]any

	// Validation Methods

	// GetValidationRules returns the validation rules for this element.
	GetValidationRules() []interface{}

	// AddValidationRule adds a validation rule to this element.
	AddValidationRule(rule interface{}) Element

	// ValidateValue validates a value against the element's validation rules.
	ValidateValue(value interface{}) error

	// GetCustomValidators returns custom validator functions for this element.
	GetCustomValidators() []interface{}

	// Display Methods

	// GetDisplayCallback returns the display callback function.
	GetDisplayCallback() func(interface{}) string

	// GetDisplayedAs returns the display format string.
	GetDisplayedAs() string

	// ShouldDisplayUsingLabels returns whether to display using labels.
	ShouldDisplayUsingLabels() bool

	// GetResolveHandle returns the resolve handle for client-side component interaction.
	GetResolveHandle() string

	// Dependency Methods

	// SetDependencies sets the field dependencies for this element.
	SetDependencies(deps []string) Element

	// GetDependencyRules returns the dependency rules for this element.
	GetDependencyRules() map[string]interface{}

	// ResolveDependencies evaluates dependency rules based on context.
	ResolveDependencies(context interface{}) bool

	// Suggestion Methods

	// GetSuggestionsCallback returns the suggestions callback function.
	GetSuggestionsCallback() func(string) []interface{}

	// GetAutoCompleteURL returns the autocomplete URL.
	GetAutoCompleteURL() string

	// GetMinCharsForSuggestions returns the minimum characters for suggestions.
	GetMinCharsForSuggestions() int

	// GetSuggestions returns suggestions for a query.
	GetSuggestions(query string) []interface{}

	// Attachment Methods

	// GetAcceptedMimeTypes returns the accepted MIME types.
	GetAcceptedMimeTypes() []string

	// GetMaxFileSize returns the maximum file size.
	GetMaxFileSize() int64

	// GetStorageDisk returns the storage disk.
	GetStorageDisk() string

	// GetStoragePath returns the storage path.
	GetStoragePath() string

	// ValidateAttachment validates an attachment.
	ValidateAttachment(filename string, size int64) error

	// GetUploadCallback returns the upload callback function.
	GetUploadCallback() func(interface{}, interface{}) error

	// ShouldRemoveEXIFData returns whether to remove EXIF data.
	ShouldRemoveEXIFData() bool

	// RemoveEXIFData removes EXIF data from a file.
	RemoveEXIFData(ctx interface{}, file interface{}) error

	// Repeater Methods

	// IsRepeaterField returns whether this is a repeater field.
	IsRepeaterField() bool

	// GetRepeaterFields returns the repeater fields.
	GetRepeaterFields() []Element

	// GetMinRepeats returns the minimum number of repeats.
	GetMinRepeats() int

	// GetMaxRepeats returns the maximum number of repeats.
	GetMaxRepeats() int

	// ValidateRepeats validates the number of repeats.
	ValidateRepeats(count int) error

	// Rich Text Methods

	// GetEditorType returns the editor type.
	GetEditorType() string

	// GetEditorLanguage returns the editor language.
	GetEditorLanguage() string

	// GetEditorTheme returns the editor theme.
	GetEditorTheme() string

	// Status Methods

	// GetStatusColors returns the status colors mapping.
	GetStatusColors() map[string]string

	// GetBadgeVariant returns the badge variant.
	GetBadgeVariant() string

	// Pivot Methods

	// IsPivot returns whether this is a pivot field.
	IsPivot() bool

	// GetPivotResourceName returns the pivot resource name.
	GetPivotResourceName() string
}
