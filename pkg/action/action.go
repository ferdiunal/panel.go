package action

import (
	"fmt"
	"strings"

	"github.com/ferdiunal/panel.go/pkg/core"
)

// Action defines the interface for resource actions.
// Actions allow users to perform bulk operations on selected resources.
type Action interface {
	// GetName returns the display name of the action
	GetName() string

	// GetSlug returns the URL-safe identifier for the action
	GetSlug() string

	// GetIcon returns the icon name for the action
	GetIcon() string

	// GetConfirmText returns the confirmation message text
	GetConfirmText() string

	// GetConfirmButtonText returns the confirm button text
	GetConfirmButtonText() string

	// GetCancelButtonText returns the cancel button text
	GetCancelButtonText() string

	// IsDestructive returns whether this action is destructive
	IsDestructive() bool

	// OnlyOnIndex returns whether this action is only available on index view
	OnlyOnIndex() bool

	// OnlyOnDetail returns whether this action is only available on detail view
	OnlyOnDetail() bool

	// ShowInline returns whether this action should be shown inline
	ShowInline() bool

	// GetFields returns the fields required for this action
	GetFields() []core.Element

	// Execute performs the action on the given context
	Execute(ctx *ActionContext) error

	// CanRun determines if the action can be executed in the given context
	CanRun(ctx *ActionContext) bool
}

// BaseAction provides a base implementation of the Action interface.
// It can be embedded in custom actions to inherit default behavior.
type BaseAction struct {
	Name              string
	Slug              string
	Icon              string
	ConfirmText       string
	ConfirmButtonText string
	CancelButtonText  string
	DestructiveType   bool
	OnlyOnIndexFlag   bool
	OnlyOnDetailFlag  bool
	ShowInlineFlag    bool
	Fields            []core.Element
	HandleFunc        func(ctx *ActionContext) error
	CanRunFunc        func(ctx *ActionContext) bool
}

// New creates a new BaseAction with the given name.
// The slug is automatically generated from the name.
func New(name string) *BaseAction {
	slug := strings.ToLower(strings.ReplaceAll(name, " ", "-"))
	return &BaseAction{
		Name:              name,
		Slug:              slug,
		ConfirmButtonText: "Confirm",
		CancelButtonText:  "Cancel",
	}
}

// Fluent API methods for configuring the action

// SetName sets the display name of the action
func (a *BaseAction) SetName(name string) *BaseAction {
	a.Name = name
	return a
}

// SetSlug sets the URL-safe identifier for the action
func (a *BaseAction) SetSlug(slug string) *BaseAction {
	a.Slug = slug
	return a
}

// SetIcon sets the icon name for the action
func (a *BaseAction) SetIcon(icon string) *BaseAction {
	a.Icon = icon
	return a
}

// Confirm sets the confirmation message text
func (a *BaseAction) Confirm(text string) *BaseAction {
	a.ConfirmText = text
	return a
}

// ConfirmButton sets the confirm button text
func (a *BaseAction) ConfirmButton(text string) *BaseAction {
	a.ConfirmButtonText = text
	return a
}

// CancelButton sets the cancel button text
func (a *BaseAction) CancelButton(text string) *BaseAction {
	a.CancelButtonText = text
	return a
}

// Destructive marks the action as destructive
func (a *BaseAction) Destructive() *BaseAction {
	a.DestructiveType = true
	return a
}

// OnlyOnIndex marks the action as only available on index view
func (a *BaseAction) OnlyOnIndex() *BaseAction {
	a.OnlyOnIndexFlag = true
	return a
}

// OnlyOnDetail marks the action as only available on detail view
func (a *BaseAction) OnlyOnDetail() *BaseAction {
	a.OnlyOnDetailFlag = true
	return a
}

// ShowInline marks the action to be shown inline
func (a *BaseAction) ShowInline() *BaseAction {
	a.ShowInlineFlag = true
	return a
}

// WithFields sets the fields required for this action
func (a *BaseAction) WithFields(fields ...core.Element) *BaseAction {
	a.Fields = fields
	return a
}

// Handle sets the action handler function
func (a *BaseAction) Handle(fn func(ctx *ActionContext) error) *BaseAction {
	a.HandleFunc = fn
	return a
}

// CanRun sets the function to determine if the action can be executed
func (a *BaseAction) CanRun(fn func(ctx *ActionContext) bool) *BaseAction {
	a.CanRunFunc = fn
	return a
}

// Interface implementation

// GetName returns the display name of the action
func (a *BaseAction) GetName() string {
	return a.Name
}

// GetSlug returns the URL-safe identifier for the action
func (a *BaseAction) GetSlug() string {
	return a.Slug
}

// GetIcon returns the icon name for the action
func (a *BaseAction) GetIcon() string {
	return a.Icon
}

// GetConfirmText returns the confirmation message text
func (a *BaseAction) GetConfirmText() string {
	return a.ConfirmText
}

// GetConfirmButtonText returns the confirm button text
func (a *BaseAction) GetConfirmButtonText() string {
	return a.ConfirmButtonText
}

// GetCancelButtonText returns the cancel button text
func (a *BaseAction) GetCancelButtonText() string {
	return a.CancelButtonText
}

// IsDestructive returns whether this action is destructive
func (a *BaseAction) IsDestructive() bool {
	return a.DestructiveType
}

// OnlyOnIndex returns whether this action is only available on index view
func (a *BaseAction) OnlyOnIndex() bool {
	return a.OnlyOnIndexFlag
}

// OnlyOnDetail returns whether this action is only available on detail view
func (a *BaseAction) OnlyOnDetail() bool {
	return a.OnlyOnDetailFlag
}

// ShowInline returns whether this action should be shown inline
func (a *BaseAction) ShowInline() bool {
	return a.ShowInlineFlag
}

// GetFields returns the fields required for this action
func (a *BaseAction) GetFields() []core.Element {
	return a.Fields
}

// Execute performs the action on the given context
func (a *BaseAction) Execute(ctx *ActionContext) error {
	if a.HandleFunc == nil {
		return fmt.Errorf("action handler not defined")
	}
	return a.HandleFunc(ctx)
}

// CanRun determines if the action can be executed in the given context
func (a *BaseAction) CanRun(ctx *ActionContext) bool {
	if a.CanRunFunc == nil {
		return true
	}
	return a.CanRunFunc(ctx)
}
