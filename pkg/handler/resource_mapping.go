package handler

import (
	stdcontext "context"
	"fmt"

	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/fields"
	internalconcurrency "github.com/ferdiunal/panel.go/pkg/internal/concurrency"
)

func resolveResourcesWithPolicy(
	h *FieldHandler,
	c *context.Context,
	items []interface{},
	elements []fields.Element,
) ([]map[string]interface{}, error) {
	if !h.usePipelineV2() {
		resources := make([]map[string]interface{}, 0, len(items))
		for _, item := range items {
			res, err := h.resolveResourceFields(c.Ctx, c.Resource(), item, elements)
			if err != nil {
				return nil, err
			}
			res["policy"] = map[string]bool{
				"view":   h.Policy == nil || h.Policy.View(c, item),
				"update": h.Policy == nil || h.Policy.Update(c, item),
				"delete": h.Policy == nil || h.Policy.Delete(c, item),
			}
			resources = append(resources, res)
		}
		return resources, nil
	}

	indices := make([]int, len(items))
	for i := range items {
		indices[i] = i
	}

	failFast := h.shouldFailFast()
	workers := h.resolveFieldWorkers(len(indices))

	resources, err := internalconcurrency.MapOrdered(c.Context(), indices, workers, failFast, func(_ stdcontext.Context, _ int, idx int) (map[string]interface{}, error) {
		item := items[idx]
		itemCtx := cloneResourceContextForIsolation(c.Resource())
		itemElements := cloneElementsForIsolation(elements)

		res, resolveErr := h.resolveResourceFields(c.Ctx, itemCtx, item, itemElements)
		policy := map[string]bool{
			"view":   h.Policy == nil || h.Policy.View(c, item),
			"update": h.Policy == nil || h.Policy.Update(c, item),
			"delete": h.Policy == nil || h.Policy.Delete(c, item),
		}

		if resolveErr != nil {
			if failFast {
				return nil, fmt.Errorf("resolve resource fields: %w", resolveErr)
			}
			return map[string]interface{}{
				"policy": policy,
				"error":  resolveErr.Error(),
			}, nil
		}

		res["policy"] = policy
		return res, nil
	})
	if err != nil {
		return nil, err
	}

	return resources, nil
}
