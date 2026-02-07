package core

import "github.com/gofiber/fiber/v2"

// Notification represents a user notification (lightweight version for context)
type Notification struct {
	Message  string `json:"message"`
	Type     string `json:"type"`
	Duration int    `json:"duration"`
	UserID   *uint  `json:"user_id,omitempty"`
}

// ResourceContext holds the context during a resource operation.
// It contains the resource being operated on, the elements (fields) associated with it,
// visibility context, the item being operated on, the user performing the operation,
// and the HTTP request context.
//
// ResourceContext is typically stored in fiber.Locals and can be retrieved using
// the ResourceContextKey constant.
//
// Requirement 15.1: THE Sistem SHALL context'in kaynak metadata'sını taşımasını sağlamalıdır
// Requirement 15.2: THE Sistem SHALL bileşenlerin context'ten alan resolver'larına erişmesine izin vermelidir
type ResourceContext struct {
	// Resource is the domain entity being operated on.
	Resource any

	// Lens is the optional lens (filtered view) being applied to the resource.
	// This is nil if no lens is active.
	Lens any

	// VisibilityCtx is the context in which fields should be visible.
	// It determines which fields are shown (index, detail, create, update, preview).
	VisibilityCtx VisibilityContext

	// Item is the specific resource instance being operated on.
	// This is the actual data object (e.g., a User struct instance).
	Item any

	// User is the user performing the operation.
	// This can be used for authorization and audit logging.
	User any

	// Elements are the fields associated with this resource.
	Elements []Element

	// fieldResolvers maps field names to their resolver functions.
	// These resolvers can be used to dynamically compute or transform field values.
	fieldResolvers map[string]Resolver

	// Request is the Fiber HTTP context.
	Request *fiber.Ctx

	// Notifications are the notifications to be sent to the user
	Notifications []Notification
}

// ResourceContextKey is the key used to store ResourceContext in fiber.Locals.
const ResourceContextKey = "resource_context"

// NewResourceContext creates a new ResourceContext with the given parameters.
//
// Parameters:
//   - c: The Fiber HTTP context
//   - resource: The domain entity being operated on
//   - elements: The fields associated with this resource
//
// Returns:
//   - A pointer to the newly created ResourceContext
func NewResourceContext(c *fiber.Ctx, resource any, elements []Element) *ResourceContext {
	return &ResourceContext{
		Resource:       resource,
		Elements:       elements,
		Request:        c,
		fieldResolvers: make(map[string]Resolver),
		Notifications:  []Notification{},
	}
}

// NewResourceContextWithVisibility creates a new ResourceContext with visibility context.
//
// Parameters:
//   - c: The Fiber HTTP context
//   - resource: The domain entity being operated on
//   - lens: The optional lens (filtered view) being applied
//   - visibilityCtx: The context in which fields should be visible
//   - item: The specific resource instance being operated on
//   - user: The user performing the operation
//   - elements: The fields associated with this resource
//
// Returns:
//   - A pointer to the newly created ResourceContext
//
// Requirement 15.1: THE Sistem SHALL context oluşturulduğunda tüm gerekli kaynak bilgisini başlatmalıdır
func NewResourceContextWithVisibility(
	c *fiber.Ctx,
	resource any,
	lens any,
	visibilityCtx VisibilityContext,
	item any,
	user any,
	elements []Element,
) *ResourceContext {
	return &ResourceContext{
		Resource:       resource,
		Lens:           lens,
		VisibilityCtx:  visibilityCtx,
		Item:           item,
		User:           user,
		Elements:       elements,
		Request:        c,
		fieldResolvers: make(map[string]Resolver),
		Notifications:  []Notification{},
	}
}

// GetFieldResolver returns the resolver for a specific field.
// If no resolver is registered for the field, it returns nil and an error.
//
// Parameters:
//   - fieldName: The name of the field to get the resolver for
//
// Returns:
//   - The Resolver for the field, or nil if not found
//   - An error if the resolver is not found
//
// Requirement 15.2: THE Sistem SHALL bileşenlerin context'ten alan resolver'larına erişmesine izin vermelidir
// Requirement 15.3: THE Sistem SHALL bileşen-alan iletişimini desteklemelidir
func (rc *ResourceContext) GetFieldResolver(fieldName string) (Resolver, error) {
	resolver, ok := rc.fieldResolvers[fieldName]
	if !ok {
		return nil, fiber.NewError(fiber.StatusNotFound, "field resolver not found: "+fieldName)
	}
	return resolver, nil
}

// SetFieldResolver registers a resolver for a specific field.
//
// Parameters:
//   - fieldName: The name of the field to register the resolver for
//   - resolver: The Resolver to register
func (rc *ResourceContext) SetFieldResolver(fieldName string, resolver Resolver) {
	rc.fieldResolvers[fieldName] = resolver
}

// ResolveField resolves a field value using the registered resolver.
// If no resolver is registered, it returns an error.
//
// Parameters:
//   - fieldName: The name of the field to resolve
//   - item: The item to resolve the field for
//   - params: Additional parameters for the resolver
//
// Returns:
//   - The resolved value
//   - An error if resolution fails or no resolver is found
//
// Requirement 15.2: THE Sistem SHALL bileşenlerin context'ten alan resolver'larına erişmesine izin vermelidir
// Requirement 15.3: THE Sistem SHALL bileşen-alan iletişimini desteklemelidir
func (rc *ResourceContext) ResolveField(fieldName string, item interface{}, params map[string]interface{}) (interface{}, error) {
	resolver, err := rc.GetFieldResolver(fieldName)
	if err != nil {
		return nil, err
	}
	return resolver.Resolve(item, params)
}

// GetResourceMetadata returns metadata about the resource and current context.
// This includes information about the resource, lens, visibility context, and available fields.
//
// Returns:
//   - A map containing resource metadata
//
// Requirement 15.4: WHEN context oluşturulduğunda, THE Sistem SHALL tüm gerekli kaynak bilgisini başlatmalıdır
func (rc *ResourceContext) GetResourceMetadata() map[string]interface{} {
	fieldNames := make([]string, 0, len(rc.Elements))
	for _, element := range rc.Elements {
		if element.IsVisible(rc) {
			fieldNames = append(fieldNames, element.GetKey())
		}
	}

	resolverNames := make([]string, 0, len(rc.fieldResolvers))
	for name := range rc.fieldResolvers {
		resolverNames = append(resolverNames, name)
	}

	return map[string]interface{}{
		"visibility_context": string(rc.VisibilityCtx),
		"fields":             fieldNames,
		"resolvers":          resolverNames,
		"has_lens":           rc.Lens != nil,
		"has_item":           rc.Item != nil,
		"has_user":           rc.User != nil,
	}
}

// Notify adds a notification to the context
func (rc *ResourceContext) Notify(message string, notifType string) {
	rc.Notifications = append(rc.Notifications, Notification{
		Message:  message,
		Type:     notifType,
		Duration: 3000,
	})
}

// NotifySuccess adds a success notification
func (rc *ResourceContext) NotifySuccess(message string) {
	rc.Notify(message, "success")
}

// NotifyError adds an error notification
func (rc *ResourceContext) NotifyError(message string) {
	rc.Notify(message, "error")
}

// NotifyWarning adds a warning notification
func (rc *ResourceContext) NotifyWarning(message string) {
	rc.Notify(message, "warning")
}

// NotifyInfo adds an info notification
func (rc *ResourceContext) NotifyInfo(message string) {
	rc.Notify(message, "info")
}

// GetNotifications returns all notifications
func (rc *ResourceContext) GetNotifications() []Notification {
	return rc.Notifications
}
