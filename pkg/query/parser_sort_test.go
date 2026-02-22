package query

import "testing"

func TestParseNestedFormat_HyphenatedResourceSortAsc(t *testing.T) {
	params := DefaultParams()

	found := parseNestedFormat(
		"page-sections[sort][id]=asc&page-sections[page]=1&page-sections[per_page]=10",
		"page-sections",
		params,
	)

	if !found {
		t.Fatalf("expected nested format to be parsed")
	}

	if len(params.Sorts) != 1 {
		t.Fatalf("expected 1 sort, got %d", len(params.Sorts))
	}

	if params.Sorts[0].Column != "id" {
		t.Fatalf("expected sort column id, got %q", params.Sorts[0].Column)
	}

	if params.Sorts[0].Direction != "asc" {
		t.Fatalf("expected sort direction asc, got %q", params.Sorts[0].Direction)
	}
}

func TestParseNestedFormat_DuplicateSortDirectionUsesLastValue(t *testing.T) {
	params := DefaultParams()

	found := parseNestedFormat(
		"page-sections[sort][id]=desc&page-sections[sort][id]=asc",
		"page-sections",
		params,
	)

	if !found {
		t.Fatalf("expected nested format to be parsed")
	}

	if len(params.Sorts) != 1 {
		t.Fatalf("expected 1 sort, got %d", len(params.Sorts))
	}

	if params.Sorts[0].Direction != "asc" {
		t.Fatalf("expected last duplicate sort direction to win (asc), got %q", params.Sorts[0].Direction)
	}
}
