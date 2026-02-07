package fields

import (
	"github.com/gofiber/fiber/v2"
)

// DependencyCallbackFunc is the callback function type for dependency changes
type DependencyCallbackFunc func(
	field *Schema,
	formData map[string]interface{},
	ctx *fiber.Ctx,
) *FieldUpdate

// FieldUpdate represents updates to be applied to a field based on dependency changes
type FieldUpdate struct {
	Visible     *bool                  `json:"visible,omitempty"`
	ReadOnly    *bool                  `json:"readonly,omitempty"`
	Required    *bool                  `json:"required,omitempty"`
	Disabled    *bool                  `json:"disabled,omitempty"`
	HelpText    *string                `json:"helpText,omitempty"`
	Placeholder *string                `json:"placeholder,omitempty"`
	Options     map[string]interface{} `json:"options,omitempty"`
	Value       interface{}            `json:"value,omitempty"`
	Rules       []ValidationRule       `json:"rules,omitempty"`
}

// NewFieldUpdate creates a new FieldUpdate instance
func NewFieldUpdate() *FieldUpdate {
	return &FieldUpdate{}
}

// Show makes the field visible
func (u *FieldUpdate) Show() *FieldUpdate {
	visible := true
	u.Visible = &visible
	return u
}

// Hide makes the field hidden
func (u *FieldUpdate) Hide() *FieldUpdate {
	visible := false
	u.Visible = &visible
	return u
}

// MakeReadOnly makes the field read-only
func (u *FieldUpdate) MakeReadOnly() *FieldUpdate {
	readOnly := true
	u.ReadOnly = &readOnly
	return u
}

// MakeEditable makes the field editable
func (u *FieldUpdate) MakeEditable() *FieldUpdate {
	readOnly := false
	u.ReadOnly = &readOnly
	return u
}

// MakeRequired makes the field required
func (u *FieldUpdate) MakeRequired() *FieldUpdate {
	required := true
	u.Required = &required
	return u
}

// MakeOptional makes the field optional
func (u *FieldUpdate) MakeOptional() *FieldUpdate {
	required := false
	u.Required = &required
	return u
}

// Enable enables the field
func (u *FieldUpdate) Enable() *FieldUpdate {
	disabled := false
	u.Disabled = &disabled
	return u
}

// Disable disables the field
func (u *FieldUpdate) Disable() *FieldUpdate {
	disabled := true
	u.Disabled = &disabled
	return u
}

// SetHelpText sets the help text
func (u *FieldUpdate) SetHelpText(text string) *FieldUpdate {
	u.HelpText = &text
	return u
}

// SetPlaceholder sets the placeholder text
func (u *FieldUpdate) SetPlaceholder(text string) *FieldUpdate {
	u.Placeholder = &text
	return u
}

// SetOptions sets the field options (for select fields)
func (u *FieldUpdate) SetOptions(options map[string]interface{}) *FieldUpdate {
	u.Options = options
	return u
}

// SetValue sets the field value
func (u *FieldUpdate) SetValue(value interface{}) *FieldUpdate {
	u.Value = value
	return u
}

// SetRules sets validation rules
func (u *FieldUpdate) SetRules(rules []ValidationRule) *FieldUpdate {
	u.Rules = rules
	return u
}

// AddRule adds a validation rule
func (u *FieldUpdate) AddRule(rule ValidationRule) *FieldUpdate {
	u.Rules = append(u.Rules, rule)
	return u
}
