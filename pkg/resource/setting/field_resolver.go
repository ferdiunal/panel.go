package setting

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/core"
	"github.com/ferdiunal/panel.go/pkg/fields"
)

// SettingFieldResolver, Setting alanlarını çözer
type SettingFieldResolver struct{}

// ResolveFields, Setting alanlarını döner
func (r *SettingFieldResolver) ResolveFields(ctx *context.Context) []core.Element {
	return []core.Element{
		(&fields.Schema{
			Key:   "id",
			Name:  "ID",
			View:  "text",
			Props: make(map[string]interface{}),
		}).ReadOnly().OnlyOnDetail(),

		(&fields.Schema{
			Key:   "key",
			Name:  "Key",
			View:  "text",
			Props: make(map[string]interface{}),
		}).OnList().OnDetail().OnForm().Required(),

		(&fields.Schema{
			Key:   "value",
			Name:  "Value",
			View:  "textarea",
			Props: make(map[string]interface{}),
		}).OnList().OnDetail().OnForm(),

		(&fields.Schema{
			Key:  "type",
			Name: "Type",
			View: "select",
			Props: map[string]interface{}{
				"options": map[string]string{
					"string":  "String",
					"integer": "Integer",
					"boolean": "Boolean",
					"json":    "JSON",
				},
			},
		}).OnList().OnDetail().OnForm(),

		(&fields.Schema{
			Key:   "group",
			Name:  "Group",
			View:  "text",
			Props: make(map[string]interface{}),
		}).OnList().OnDetail().OnForm(),

		(&fields.Schema{
			Key:   "label",
			Name:  "Label",
			View:  "text",
			Props: make(map[string]interface{}),
		}).OnDetail().OnForm(),

		(&fields.Schema{
			Key:   "help",
			Name:  "Help Text",
			View:  "textarea",
			Props: make(map[string]interface{}),
		}).OnDetail().OnForm(),

		(&fields.Schema{
			Key:   "created_at",
			Name:  "Created At",
			View:  "datetime",
			Props: make(map[string]interface{}),
		}).ReadOnly().OnList().OnDetail(),

		(&fields.Schema{
			Key:   "updated_at",
			Name:  "Updated At",
			View:  "datetime",
			Props: make(map[string]interface{}),
		}).ReadOnly().OnList().OnDetail(),
	}
}
