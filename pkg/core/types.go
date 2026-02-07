package core

// ElementType represents the type of an element (field).
// It determines how the element is rendered and how it handles data.
type ElementType string

// ElementContext represents the context in which an element is displayed.
// It controls the visibility and behavior of elements in different views.
type ElementContext string

// VisibilityContext represents the context in which a field should be visible.
// It defines specific UI contexts like index, detail, create, update, and preview.
type VisibilityContext string

// Element type constants define the available field types.
const (
	// TYPE_TEXT represents a text input field.
	TYPE_TEXT ElementType = "text"

	// TYPE_TEXTAREA represents a multi-line text input field.
	TYPE_TEXTAREA ElementType = "textarea"

	// TYPE_RICHTEXT represents a rich text editor field (WYSIWYG).
	TYPE_RICHTEXT ElementType = "richtext"

	// TYPE_PASSWORD represents a password input field.
	TYPE_PASSWORD ElementType = "password"

	// TYPE_NUMBER represents a numeric input field.
	TYPE_NUMBER ElementType = "number"

	// TYPE_TEL represents a telephone number input field.
	TYPE_TEL ElementType = "tel"

	// TYPE_EMAIL represents an email input field.
	TYPE_EMAIL ElementType = "email"

	// TYPE_AUDIO represents an audio file upload field.
	TYPE_AUDIO ElementType = "audio"

	// TYPE_VIDEO represents a video file upload field.
	TYPE_VIDEO ElementType = "video"

	// TYPE_DATE represents a date picker field.
	TYPE_DATE ElementType = "date"

	// TYPE_DATETIME represents a date and time picker field.
	TYPE_DATETIME ElementType = "datetime"

	// TYPE_FILE represents a file upload field.
	TYPE_FILE ElementType = "file"

	// TYPE_KEY_VALUE represents a key-value pair field.
	TYPE_KEY_VALUE ElementType = "key_value"

	// TYPE_LINK represents a link to another resource.
	TYPE_LINK ElementType = "link"

	// TYPE_COLLECTION represents a collection of related resources.
	TYPE_COLLECTION ElementType = "collection"

	// TYPE_DETAIL represents a detailed view of a related resource.
	TYPE_DETAIL ElementType = "detail"

	// TYPE_CONNECT represents a connection to another resource.
	TYPE_CONNECT ElementType = "connect"

	// TYPE_POLY_LINK represents a polymorphic link to another resource.
	TYPE_POLY_LINK ElementType = "poly_link"

	// TYPE_POLY_DETAIL represents a polymorphic detailed view of a related resource.
	TYPE_POLY_DETAIL ElementType = "poly_detail"

	// TYPE_POLY_COLLECTION represents a polymorphic collection of related resources.
	TYPE_POLY_COLLECTION ElementType = "poly_collection"

	// TYPE_POLY_CONNECT represents a polymorphic connection to another resource.
	TYPE_POLY_CONNECT ElementType = "poly_connect"

	// TYPE_BOOLEAN represents a boolean checkbox field.
	TYPE_BOOLEAN ElementType = "boolean"

	// TYPE_SELECT represents a select dropdown field.
	TYPE_SELECT ElementType = "select"

	// TYPE_PANEL represents a panel/section container for grouping fields.
	TYPE_PANEL ElementType = "panel"

	// TYPE_RELATIONSHIP represents a relationship field.
	TYPE_RELATIONSHIP ElementType = "relationship"
)

// Element context constants define where elements are displayed.
const (
	// CONTEXT_FORM indicates the element is in a form context.
	CONTEXT_FORM ElementContext = "form"

	// CONTEXT_DETAIL indicates the element is in a detail view context.
	CONTEXT_DETAIL ElementContext = "detail"

	// CONTEXT_LIST indicates the element is in a list view context.
	CONTEXT_LIST ElementContext = "list"

	// SHOW_ON_FORM indicates the element should be shown on forms.
	SHOW_ON_FORM ElementContext = "show_on_form"

	// SHOW_ON_DETAIL indicates the element should be shown on detail views.
	SHOW_ON_DETAIL ElementContext = "show_on_detail"

	// SHOW_ON_LIST indicates the element should be shown on list views.
	SHOW_ON_LIST ElementContext = "show_on_list"

	// HIDE_ON_LIST indicates the element should be hidden on list views.
	HIDE_ON_LIST ElementContext = "hide_on_list"

	// HIDE_ON_DETAIL indicates the element should be hidden on detail views.
	HIDE_ON_DETAIL ElementContext = "hide_on_detail"

	// HIDE_ON_CREATE indicates the element should be hidden on create forms.
	HIDE_ON_CREATE ElementContext = "hide_on_create"

	// HIDE_ON_UPDATE indicates the element should be hidden on update forms.
	HIDE_ON_UPDATE ElementContext = "hide_on_update"

	// ONLY_ON_LIST indicates the element should only be shown on list views.
	ONLY_ON_LIST ElementContext = "only_on_list"

	// ONLY_ON_DETAIL indicates the element should only be shown on detail views.
	ONLY_ON_DETAIL ElementContext = "only_on_detail"

	// ONLY_ON_CREATE indicates the element should only be shown on create forms.
	ONLY_ON_CREATE ElementContext = "only_on_create"

	// ONLY_ON_UPDATE indicates the element should only be shown on update forms.
	ONLY_ON_UPDATE ElementContext = "only_on_update"

	// ONLY_ON_FORM indicates the element should only be shown on forms (create and update).
	ONLY_ON_FORM ElementContext = "only_on_form"
)

// Visibility context constants define specific UI contexts for field visibility.
const (
	// ContextIndex indicates the field is in an index/list view context.
	ContextIndex VisibilityContext = "index"

	// ContextDetail indicates the field is in a detail view context.
	ContextDetail VisibilityContext = "detail"

	// ContextCreate indicates the field is in a create form context.
	ContextCreate VisibilityContext = "create"

	// ContextUpdate indicates the field is in an update form context.
	ContextUpdate VisibilityContext = "update"

	// ContextPreview indicates the field is in a preview context.
	ContextPreview VisibilityContext = "preview"
)
