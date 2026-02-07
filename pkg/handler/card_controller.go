package handler

import (
	"fmt"
	"sync"

	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/widget"
	"github.com/gofiber/fiber/v2"
)

// cardResult holds the result of a card resolution
type cardResult struct {
	card       widget.Card
	data       interface{}
	err        error
	index      int
	serialized map[string]interface{}
}

// HandleCardList handles listing all cards for a resource.
// It resolves each card's data in parallel using async fan-out/fan-in pattern.
func HandleCardList(h *FieldHandler, c *context.Context) error {
	if len(h.Cards) == 0 {
		return c.JSON(fiber.Map{
			"data": []map[string]interface{}{},
		})
	}

	// Create buffered channel for results (non-blocking sends)
	results := make(chan cardResult, len(h.Cards))

	// WaitGroup to track goroutine completion
	var wg sync.WaitGroup
	wg.Add(len(h.Cards))

	// Fan-out: Launch goroutines asynchronously for each card
	for i, card := range h.Cards {
		go func(idx int, w widget.Card) {
			defer wg.Done() // Mark goroutine as done when finished

			// Serialize base properties
			serialized := w.JsonSerialize()
			serialized["index"] = idx
			serialized["name"] = w.Name()
			serialized["component"] = w.Component()
			serialized["width"] = w.Width()

			// Resolve data
			data, err := w.Resolve(c, h.DB)

			// Send result to channel
			results <- cardResult{
				card:       w,
				data:       data,
				err:        err,
				index:      idx,
				serialized: serialized,
			}
		}(i, card)
	}

	// Close channel when all goroutines complete (async closer)
	go func() {
		wg.Wait()      // Wait for all goroutines to finish
		close(results) // Close channel to signal completion
	}()

	// Fan-in: Collect results from channel
	resp := make([]map[string]interface{}, len(h.Cards))

	for result := range results {
		if result.err != nil {
			fmt.Printf("Error resolving card %s: %v\n", result.card.Name(), result.err)
			result.serialized["error"] = result.err.Error()
		} else {
			// Assign resolved data to "data" key
			result.serialized["data"] = result.data
		}

		// Store result at original index to maintain order
		resp[result.index] = result.serialized
	}

	return c.JSON(fiber.Map{
		"data": resp,
	})
}
