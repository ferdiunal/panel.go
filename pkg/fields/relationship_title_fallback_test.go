package fields

import "testing"

type relationshipTitleTestResource struct {
	title string
}

func (r relationshipTitleTestResource) RecordTitle(any) string {
	return r.title
}

func TestResolveRelationshipRecordTitle(t *testing.T) {
	type testRecord struct {
		ID    int
		Name  string
		Title string
	}

	record := testRecord{
		ID:    10,
		Name:  "Product Name",
		Title: "Product Title",
	}

	tests := []struct {
		name       string
		resource   relationshipTitleTestResource
		idValue    any
		expectText string
	}{
		{
			name:       "uses record title when meaningful",
			resource:   relationshipTitleTestResource{title: "Custom Title"},
			idValue:    10,
			expectText: "Custom Title",
		},
		{
			name:       "falls back when record title is empty",
			resource:   relationshipTitleTestResource{title: ""},
			idValue:    10,
			expectText: "Product Name",
		},
		{
			name:       "falls back when record title equals id",
			resource:   relationshipTitleTestResource{title: "10"},
			idValue:    10,
			expectText: "Product Name",
		},
		{
			name:       "falls back when record title equals hash id",
			resource:   relationshipTitleTestResource{title: "#10"},
			idValue:    10,
			expectText: "Product Name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolveRelationshipRecordTitle(tt.resource, record, tt.idValue)
			if got != tt.expectText {
				t.Fatalf("expected %q, got %q", tt.expectText, got)
			}
		})
	}
}

func TestResolveRelationshipRecordTitleWithoutFallbackField(t *testing.T) {
	type testRecord struct {
		ID int
	}

	record := testRecord{ID: 22}
	res := relationshipTitleTestResource{title: "22"}

	got := resolveRelationshipRecordTitle(res, record, 22)
	if got != "22" {
		t.Fatalf("expected id title to be preserved when no fallback exists, got %q", got)
	}
}
