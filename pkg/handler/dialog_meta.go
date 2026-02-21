package handler

import "github.com/ferdiunal/panel.go/pkg/resource"

func resolveDialogMeta(h *FieldHandler) (resource.DialogType, resource.DialogSize) {
	dialogType := resource.DialogTypeSheet
	dialogSize := resource.DialogSizeMD

	if h == nil {
		return dialogType, dialogSize
	}

	if h.DialogType != "" {
		dialogType = h.DialogType
	}
	if h.DialogSize != "" {
		dialogSize = h.DialogSize
	}

	return dialogType, dialogSize
}
