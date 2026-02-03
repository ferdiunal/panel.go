package fields

import (
	"context"
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
)

// Schema, bir alanın temel yapılandırmasını ve durumunu tutan yapıdır.
// JSON serileştirme ve veri taşıma için kullanılır.
type Schema struct {
	Name               string                              `json:"name"`      // Görünen Ad
	Key                string                              `json:"key"`       // Veri Anahtarı
	View               string                              `json:"view"`      // Frontend Bileşeni
	Data               interface{}                         `json:"data"`      // Alan Değeri
	Type               ElementType                         `json:"type"`      // Veri Tipi
	Context            ElementContext                      `json:"context"`   // Görünüm Bağlamı (List, Detail, Form)
	IsReadOnly         bool                                `json:"read_only"` // Salt okunur mu?
	IsDisabled         bool                                `json:"disabled"`  // Devre dışı mı?
	IsImmutable        bool                                `json:"immutable"` // Değiştirilemez mi?
	Props              map[string]interface{}              `json:"props"`     // Ekstra özellikler
	IsRequired         bool                                `json:"required"`  // Zorunlu mu?
	IsNullable         bool                                `json:"nullable"`  // Boş bırakılabilir mi?
	PlaceholderText    string                              `json:"placeholder"`
	LabelText          string                              `json:"label"`
	HelpTextContent    string                              `json:"help_text"`
	IsFilterable       bool                                `json:"filterable"`
	IsSortable         bool                                `json:"sortable"`
	GlobalSearch       bool                                `json:"searchable"`
	IsStacked          bool                                `json:"stacked"`
	TextAlign          string                              `json:"text_align"`
	Suggestions        []interface{}                       `json:"suggestions"`
	ExtractCallback    func(value interface{}) interface{} `json:"-"`
	VisibilityCallback VisibilityFunc                      `json:"-"`
	StorageCallback    StorageCallbackFunc                 `json:"-"`
	ModifyCallback     func(value interface{}) interface{} `json:"-"`
}

func (s *Schema) GetKey() string {
	return s.Key
}

func (s *Schema) GetView() string {
	return s.View
}

func (s *Schema) GetContext() ElementContext {
	return s.Context
}

func (s *Schema) Extract(resource interface{}) {
	if resource == nil {
		return
	}

	// Use reflection to get the value
	v := reflect.ValueOf(resource)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	var value interface{}

	switch v.Kind() {
	case reflect.Struct:
		// Try to find the field by name or json tag
		fieldVal := v.FieldByName(strcase.ToCamel(s.Key))
		if !fieldVal.IsValid() {
			// Iterate over fields to check json tags if name doesn't match directly
			for i := 0; i < v.NumField(); i++ {
				typeField := v.Type().Field(i)
				tag := typeField.Tag.Get("json")
				if tag == s.Key || strings.Split(tag, ",")[0] == s.Key {
					fieldVal = v.Field(i)
					break
				}
			}
		}

		if fieldVal.IsValid() && fieldVal.CanInterface() {
			value = fieldVal.Interface()
		}
	case reflect.Map:
		// Check if it's a map[string]interface{} or similar
		val := v.MapIndex(reflect.ValueOf(s.Key))
		if val.IsValid() && val.CanInterface() {
			value = val.Interface()
		}
	}

	s.Data = value
}

func (s *Schema) JsonSerialize() map[string]interface{} {
	return map[string]interface{}{
		"view":        s.View,
		"type":        s.Type,
		"key":         s.Key,
		"name":        s.Name,
		"data":        s.Data,
		"props":       s.Props,
		"context":     s.Context,
		"placeholder": s.PlaceholderText,
		"label":       s.LabelText,
		"help_text":   s.HelpTextContent,
		"read_only":   s.IsReadOnly,
		"disabled":    s.IsDisabled,
		"required":    s.IsRequired,
		"nullable":    s.IsNullable,
		"sortable":    s.IsSortable,
		"filterable":  s.IsFilterable,
		"stacked":     s.IsStacked,
		"text_align":  s.TextAlign,
	}
}

// Fluent Setters

func (s *Schema) SetName(name string) Element {
	s.Name = name
	return s
}

func (s *Schema) SetKey(key string) Element {
	s.Key = key
	return s
}

func (s *Schema) SetContext(context ElementContext) Element {
	s.Context = context
	return s
}

func (s *Schema) OnList() Element {
	return s.SetContext(SHOW_ON_LIST)
}

func (s *Schema) OnDetail() Element {
	return s.SetContext(SHOW_ON_DETAIL)
}

func (s *Schema) OnForm() Element {
	return s.SetContext(SHOW_ON_FORM)
}

func (s *Schema) HideOnList() Element {
	return s.SetContext(HIDE_ON_LIST)
}

func (s *Schema) HideOnDetail() Element {
	return s.SetContext(HIDE_ON_DETAIL)
}

func (s *Schema) HideOnCreate() Element {
	return s.SetContext(HIDE_ON_CREATE)
}

func (s *Schema) HideOnUpdate() Element {
	return s.SetContext(HIDE_ON_UPDATE)
}

func (s *Schema) OnlyOnList() Element {
	return s.SetContext(ONLY_ON_LIST)
}

func (s *Schema) OnlyOnDetail() Element {
	return s.SetContext(ONLY_ON_DETAIL)
}

func (s *Schema) OnlyOnCreate() Element {
	return s.SetContext(ONLY_ON_CREATE)
}

func (s *Schema) OnlyOnUpdate() Element {
	return s.SetContext(ONLY_ON_UPDATE)
}

func (s *Schema) OnlyOnForm() Element {
	return s.SetContext(ONLY_ON_FORM)
}

func (s *Schema) ReadOnly() Element {
	s.IsReadOnly = true
	return s
}

func (s *Schema) WithProps(key string, value interface{}) Element {
	s.Props[key] = value
	return s
}

func (s *Schema) Disabled() Element {
	s.IsDisabled = true
	return s
}

func (s *Schema) Immutable() Element {
	s.IsImmutable = true
	return s
}

func (s *Schema) Required() Element {
	s.IsRequired = true
	return s
}

func (s *Schema) Nullable() Element {
	s.IsNullable = true
	return s
}

func (s *Schema) Placeholder(placeholder string) Element {
	s.PlaceholderText = placeholder
	return s
}

func (s *Schema) Label(label string) Element {
	s.LabelText = label
	return s
}

func (s *Schema) HelpText(helpText string) Element {
	s.HelpTextContent = helpText
	return s
}

func (s *Schema) Filterable() Element {
	s.IsFilterable = true
	return s
}

func (s *Schema) Sortable() Element {
	s.IsSortable = true
	return s
}

func (s *Schema) Searchable() Element {
	s.GlobalSearch = true
	return s
}

func (s *Schema) IsSearchable() bool {
	return s.GlobalSearch
}

func (s *Schema) Stacked() Element {
	s.IsStacked = true
	return s
}

func (s *Schema) SetTextAlign(align string) Element {
	s.TextAlign = align
	return s
}

func (s *Schema) IsVisible(ctx context.Context) bool {
	if s.VisibilityCallback != nil {
		return s.VisibilityCallback(ctx)
	}
	return true
}

func (s *Schema) CanSee(fn VisibilityFunc) Element {
	s.VisibilityCallback = fn
	return s
}

func (s *Schema) StoreAs(fn StorageCallbackFunc) Element {
	s.StorageCallback = fn
	return s
}

func (s *Schema) GetStorageCallback() StorageCallbackFunc {
	return s.StorageCallback
}

func (s *Schema) Resolve(fn func(value interface{}) interface{}) Element {
	s.ExtractCallback = fn
	return s
}

func (s *Schema) GetResolveCallback() func(value interface{}) interface{} {
	return s.ExtractCallback
}

func (s *Schema) Modify(fn func(value interface{}) interface{}) Element {
	s.ModifyCallback = fn
	return s
}

func (s *Schema) GetModifyCallback() func(value interface{}) interface{} {
	return s.ModifyCallback
}

func (s *Schema) Options(options interface{}) Element {
	s.Props["options"] = options
	return s
}

func (s *Schema) Default(value interface{}) Element {
	s.Data = value
	return s
}
