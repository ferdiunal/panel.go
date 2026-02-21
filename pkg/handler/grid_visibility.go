package handler

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/core"
	"github.com/ferdiunal/panel.go/pkg/fields"
)

// gridHiddenDataElements selects fields that are visible on index/table but hidden on grid cards.
// These fields should remain available in row payload for actions/UX flows.
func gridHiddenDataElements(elements []fields.Element, resourceCtx *core.ResourceContext) []fields.Element {
	if len(elements) == 0 {
		return nil
	}

	indexCtx := cloneResourceContextForIsolation(resourceCtx)
	if indexCtx == nil {
		indexCtx = &core.ResourceContext{}
	}
	indexCtx.VisibilityCtx = core.ContextIndex

	gridCtx := cloneResourceContextForIsolation(resourceCtx)
	if gridCtx == nil {
		gridCtx = &core.ResourceContext{}
	}
	gridCtx.VisibilityCtx = core.ContextGrid

	filtered := make([]fields.Element, 0, len(elements))
	for _, element := range elements {
		if element == nil {
			continue
		}
		if element.IsVisible(indexCtx) && !element.IsVisible(gridCtx) {
			filtered = append(filtered, element)
		}
	}

	return filtered
}

// mergeGridHiddenDataIntoRows appends index-visible/grid-hidden fields into list payload rows.
func mergeGridHiddenDataIntoRows(
	h *FieldHandler,
	c *context.Context,
	items []interface{},
	rows []map[string]interface{},
	candidateElements []fields.Element,
) error {
	if h == nil || c == nil || len(candidateElements) == 0 {
		return nil
	}
	if len(items) == 0 || len(rows) == 0 {
		return nil
	}

	max := len(items)
	if len(rows) < max {
		max = len(rows)
	}

	for i := 0; i < max; i++ {
		row := rows[i]
		if row == nil {
			continue
		}

		itemCtx := cloneResourceContextForIsolation(c.Resource())
		if itemCtx == nil {
			itemCtx = core.NewResourceContextWithVisibility(
				c.Ctx,
				h.Resource,
				h.Lens,
				core.ContextIndex,
				nil,
				c.User(),
				nil,
			)
		} else {
			itemCtx.VisibilityCtx = core.ContextIndex
		}

		extraData, err := h.resolveResourceFields(c.Ctx, itemCtx, items[i], candidateElements)
		if err != nil {
			return err
		}

		for key, value := range extraData {
			if _, exists := row[key]; exists {
				continue
			}
			row[key] = value
		}
	}

	return nil
}
