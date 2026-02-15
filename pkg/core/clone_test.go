package core_test

import (
	"testing"

	"github.com/ferdiunal/panel.go/pkg/core"
	"github.com/ferdiunal/panel.go/pkg/fields"
)

func TestCloneElement_IsolatesNestedMapState(t *testing.T) {
	original := fields.Text("Name", "name")
	original.Props["meta"] = map[string]interface{}{
		"nested": map[string]interface{}{
			"value": "original",
		},
	}

	clonedElement := core.CloneElement(original)
	cloned, ok := clonedElement.(*fields.Schema)
	if !ok {
		t.Fatalf("expected *fields.Schema clone, got %T", clonedElement)
	}
	if cloned == original {
		t.Fatalf("expected cloned pointer to differ from original")
	}

	nested := cloned.Props["meta"].(map[string]interface{})["nested"].(map[string]interface{})
	nested["value"] = "changed"

	originalNested := original.Props["meta"].(map[string]interface{})["nested"].(map[string]interface{})
	if originalNested["value"] != "original" {
		t.Fatalf("expected original nested map to stay isolated, got %v", originalNested["value"])
	}
}

func TestResourceContextGetOrCloneField_CachesIsolatedClone(t *testing.T) {
	rc := core.NewResourceContext(nil, nil, nil)
	original := fields.Text("Name", "name")

	first := rc.GetOrCloneField("name", original)
	second := rc.GetOrCloneField("name", original)

	if first == nil || second == nil {
		t.Fatalf("expected cloned field instances")
	}
	if first != second {
		t.Fatalf("expected cache hit to return same cloned field instance")
	}
	if first == original {
		t.Fatalf("expected cached field to be a clone, not the original")
	}
}
