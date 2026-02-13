package fields

// CloneElement creates a request-local copy of an element.
// Schema elements are copied with their mutable maps/slices detached.
func CloneElement(element Element) Element {
	if element == nil {
		return nil
	}

	schema, ok := element.(*Schema)
	if !ok {
		return element
	}

	cloned := *schema

	if schema.Props != nil {
		cloned.Props = make(map[string]interface{}, len(schema.Props))
		for k, v := range schema.Props {
			cloned.Props[k] = v
		}
	}

	if schema.Suggestions != nil {
		cloned.Suggestions = append([]interface{}(nil), schema.Suggestions...)
	}

	return &cloned
}

// CloneElements clones all elements in order.
func CloneElements(elements []Element) []Element {
	if len(elements) == 0 {
		return []Element{}
	}

	cloned := make([]Element, 0, len(elements))
	for _, element := range elements {
		cloned = append(cloned, CloneElement(element))
	}
	return cloned
}
