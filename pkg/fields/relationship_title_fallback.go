package fields

import (
	"fmt"
	"reflect"
	"strings"
)

type resourceWithRecordTitle interface {
	RecordTitle(any) string
}

// resolveRelationshipRecordTitle returns RecordTitle and falls back to common
// human-readable struct fields when RecordTitle is empty or equals raw id.
func resolveRelationshipRecordTitle(res resourceWithRecordTitle, record any, idValue any) string {
	if res == nil {
		return ""
	}

	title := strings.TrimSpace(res.RecordTitle(record))
	fallback := strings.TrimSpace(extractFallbackTitle(record))
	if fallback == "" {
		return title
	}

	idText := strings.TrimSpace(fmt.Sprint(idValue))
	if title == "" || (idText != "" && title == idText) || title == "#"+idText {
		return fallback
	}

	return title
}

func extractFallbackTitle(record any) string {
	if record == nil {
		return ""
	}

	v := reflect.ValueOf(record)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return ""
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return ""
	}

	// Prefer the most common title fields.
	preferred := []string{"Name", "Title", "Label", "FullName", "DisplayName", "Slug"}
	for _, fieldName := range preferred {
		field := v.FieldByName(fieldName)
		if !field.IsValid() || !field.CanInterface() {
			continue
		}
		value := strings.TrimSpace(fmt.Sprint(field.Interface()))
		if value != "" && value != "<nil>" {
			return value
		}
	}

	return ""
}
