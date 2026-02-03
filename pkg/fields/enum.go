package fields

type ElementType string

type ElementContext string

const (
	TYPE_TEXT            ElementType = "text"
	TYPE_PASSWORD        ElementType = "password"
	TYPE_NUMBER          ElementType = "number"
	TYPE_TEL             ElementType = "tel"
	TYPE_EMAIL           ElementType = "email"
	TYPE_AUDIO           ElementType = "audio"
	TYPE_VIDEO           ElementType = "video"
	TYPE_DATE            ElementType = "date"
	TYPE_DATETIME        ElementType = "datetime"
	TYPE_FILE            ElementType = "file"
	TYPE_KEY_VALUE       ElementType = "key_value"
	TYPE_LINK            ElementType = "link"
	TYPE_COLLECTION      ElementType = "collection"
	TYPE_DETAIL          ElementType = "detail"
	TYPE_CONNECT         ElementType = "connect"
	TYPE_POLY_LINK       ElementType = "poly_link"
	TYPE_POLY_DETAIL     ElementType = "poly_detail"
	TYPE_POLY_COLLECTION ElementType = "poly_collection"
	TYPE_POLY_CONNECT    ElementType = "poly_connect"
	TYPE_BOOLEAN         ElementType = "boolean"
)

const (
	CONTEXT_FORM   ElementContext = "form"
	CONTEXT_DETAIL ElementContext = "detail"
	CONTEXT_LIST   ElementContext = "list"

	// Visibility Flags (can be combined or specific logic used in Schema)
	// For simplicity, we'll map standard "show/hide" concepts to contexts

	SHOW_ON_FORM   ElementContext = "show_on_form"
	SHOW_ON_DETAIL ElementContext = "show_on_detail"
	SHOW_ON_LIST   ElementContext = "show_on_list"

	HIDE_ON_LIST   ElementContext = "hide_on_list"
	HIDE_ON_DETAIL ElementContext = "hide_on_detail"
	HIDE_ON_CREATE ElementContext = "hide_on_create"
	HIDE_ON_UPDATE ElementContext = "hide_on_update"

	ONLY_ON_LIST   ElementContext = "only_on_list"
	ONLY_ON_DETAIL ElementContext = "only_on_detail"
	ONLY_ON_CREATE ElementContext = "only_on_create"
	ONLY_ON_UPDATE ElementContext = "only_on_update"
	ONLY_ON_FORM   ElementContext = "only_on_form"
)
