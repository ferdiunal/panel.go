package core_test

import (
	"testing"

	"github.com/ferdiunal/panel.go/pkg/core"
)

// TestElementTypeConstants verifies that all ElementType constants are defined
func TestElementTypeConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant core.ElementType
		expected string
	}{
		{"TYPE_TEXT", core.TYPE_TEXT, "text"},
		{"TYPE_PASSWORD", core.TYPE_PASSWORD, "password"},
		{"TYPE_NUMBER", core.TYPE_NUMBER, "number"},
		{"TYPE_TEL", core.TYPE_TEL, "tel"},
		{"TYPE_EMAIL", core.TYPE_EMAIL, "email"},
		{"TYPE_AUDIO", core.TYPE_AUDIO, "audio"},
		{"TYPE_VIDEO", core.TYPE_VIDEO, "video"},
		{"TYPE_DATE", core.TYPE_DATE, "date"},
		{"TYPE_DATETIME", core.TYPE_DATETIME, "datetime"},
		{"TYPE_FILE", core.TYPE_FILE, "file"},
		{"TYPE_KEY_VALUE", core.TYPE_KEY_VALUE, "key_value"},
		{"TYPE_LINK", core.TYPE_LINK, "link"},
		{"TYPE_COLLECTION", core.TYPE_COLLECTION, "collection"},
		{"TYPE_DETAIL", core.TYPE_DETAIL, "detail"},
		{"TYPE_CONNECT", core.TYPE_CONNECT, "connect"},
		{"TYPE_POLY_LINK", core.TYPE_POLY_LINK, "poly_link"},
		{"TYPE_POLY_DETAIL", core.TYPE_POLY_DETAIL, "poly_detail"},
		{"TYPE_POLY_COLLECTION", core.TYPE_POLY_COLLECTION, "poly_collection"},
		{"TYPE_POLY_CONNECT", core.TYPE_POLY_CONNECT, "poly_connect"},
		{"TYPE_BOOLEAN", core.TYPE_BOOLEAN, "boolean"},
		{"TYPE_SELECT", core.TYPE_SELECT, "select"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.constant) != tt.expected {
				t.Errorf("%s = %q, want %q", tt.name, tt.constant, tt.expected)
			}
		})
	}
}

// TestElementContextConstants verifies that all ElementContext constants are defined
func TestElementContextConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant core.ElementContext
		expected string
	}{
		{"CONTEXT_FORM", core.CONTEXT_FORM, "form"},
		{"CONTEXT_DETAIL", core.CONTEXT_DETAIL, "detail"},
		{"CONTEXT_LIST", core.CONTEXT_LIST, "list"},
		{"SHOW_ON_FORM", core.SHOW_ON_FORM, "show_on_form"},
		{"SHOW_ON_DETAIL", core.SHOW_ON_DETAIL, "show_on_detail"},
		{"SHOW_ON_LIST", core.SHOW_ON_LIST, "show_on_list"},
		{"HIDE_ON_LIST", core.HIDE_ON_LIST, "hide_on_list"},
		{"HIDE_ON_DETAIL", core.HIDE_ON_DETAIL, "hide_on_detail"},
		{"HIDE_ON_CREATE", core.HIDE_ON_CREATE, "hide_on_create"},
		{"HIDE_ON_UPDATE", core.HIDE_ON_UPDATE, "hide_on_update"},
		{"HIDE_ON_API", core.HIDE_ON_API, "hide_on_api"},
		{"ONLY_ON_LIST", core.ONLY_ON_LIST, "only_on_list"},
		{"ONLY_ON_DETAIL", core.ONLY_ON_DETAIL, "only_on_detail"},
		{"ONLY_ON_CREATE", core.ONLY_ON_CREATE, "only_on_create"},
		{"ONLY_ON_UPDATE", core.ONLY_ON_UPDATE, "only_on_update"},
		{"ONLY_ON_FORM", core.ONLY_ON_FORM, "only_on_form"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.constant) != tt.expected {
				t.Errorf("%s = %q, want %q", tt.name, tt.constant, tt.expected)
			}
		})
	}
}

// TestElementTypeIsString verifies that ElementType is a string type
func TestElementTypeIsString(t *testing.T) {
	var et core.ElementType = "custom"
	if string(et) != "custom" {
		t.Errorf("ElementType should be convertible to string")
	}
}

// TestElementContextIsString verifies that ElementContext is a string type
func TestElementContextIsString(t *testing.T) {
	var ec core.ElementContext = "custom"
	if string(ec) != "custom" {
		t.Errorf("ElementContext should be convertible to string")
	}
}
