package fields

import (
	"context"
	"testing"
)

// TestNewRelationshipPagination tests creating a new relationship pagination handler
func TestNewRelationshipPagination(t *testing.T) {
	field := NewBelongsTo("Author", "author_id", "authors")
	pagination := NewRelationshipPagination(field)

	if pagination == nil {
		t.Error("Expected non-nil pagination handler")
	}

	if pagination.field != field {
		t.Error("Expected field to be set")
	}

	if pagination.page != 1 {
		t.Errorf("Expected page 1, got %d", pagination.page)
	}

	if pagination.perPage != 15 {
		t.Errorf("Expected perPage 15, got %d", pagination.perPage)
	}
}

// TestApplyPagination tests applying pagination
func TestApplyPagination(t *testing.T) {
	field := NewBelongsTo("Author", "author_id", "authors")
	pagination := NewRelationshipPagination(field)

	ctx := context.Background()
	results, err := pagination.ApplyPagination(ctx, 2, 10)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if results == nil {
		t.Error("Expected non-nil results")
	}

	if pagination.page != 2 {
		t.Errorf("Expected page 2, got %d", pagination.page)
	}

	if pagination.perPage != 10 {
		t.Errorf("Expected perPage 10, got %d", pagination.perPage)
	}
}

// TestApplyPaginationInvalidPage tests applying pagination with invalid page
func TestApplyPaginationInvalidPage(t *testing.T) {
	field := NewBelongsTo("Author", "author_id", "authors")
	pagination := NewRelationshipPagination(field)

	ctx := context.Background()
	_, err := pagination.ApplyPagination(ctx, 0, 10)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if pagination.page != 1 {
		t.Errorf("Expected page 1 (default), got %d", pagination.page)
	}
}

// TestApplyPaginationInvalidPerPage tests applying pagination with invalid per-page
func TestApplyPaginationInvalidPerPage(t *testing.T) {
	field := NewBelongsTo("Author", "author_id", "authors")
	pagination := NewRelationshipPagination(field)

	ctx := context.Background()
	_, err := pagination.ApplyPagination(ctx, 1, 0)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if pagination.perPage != 15 {
		t.Errorf("Expected perPage 15 (default), got %d", pagination.perPage)
	}
}

// TestApplyPaginationMaxPerPage tests applying pagination with per-page exceeding maximum
func TestApplyPaginationMaxPerPage(t *testing.T) {
	field := NewBelongsTo("Author", "author_id", "authors")
	pagination := NewRelationshipPagination(field)

	ctx := context.Background()
	_, err := pagination.ApplyPagination(ctx, 1, 200)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if pagination.perPage != 100 {
		t.Errorf("Expected perPage 100 (max), got %d", pagination.perPage)
	}
}

// TestGetPageInfo tests getting page info
func TestGetPageInfo(t *testing.T) {
	field := NewBelongsTo("Author", "author_id", "authors")
	pagination := NewRelationshipPagination(field)

	ctx := context.Background()
	pagination.ApplyPagination(ctx, 2, 10)
	pagination.SetTotal(100)

	pageInfo := pagination.GetPageInfo()

	if pageInfo["current_page"] != 2 {
		t.Errorf("Expected current_page 2, got %v", pageInfo["current_page"])
	}

	if pageInfo["per_page"] != 10 {
		t.Errorf("Expected per_page 10, got %v", pageInfo["per_page"])
	}

	if pageInfo["total"] != int64(100) {
		t.Errorf("Expected total 100, got %v", pageInfo["total"])
	}

	if pageInfo["total_pages"] != int64(10) {
		t.Errorf("Expected total_pages 10, got %v", pageInfo["total_pages"])
	}

	if pageInfo["from"] != 10 {
		t.Errorf("Expected from 10, got %v", pageInfo["from"])
	}

	if pageInfo["to"] != 20 {
		t.Errorf("Expected to 20, got %v", pageInfo["to"])
	}
}

// TestGetPageInfoFirstPage tests getting page info for first page
func TestGetPageInfoFirstPage(t *testing.T) {
	field := NewBelongsTo("Author", "author_id", "authors")
	pagination := NewRelationshipPagination(field)

	ctx := context.Background()
	pagination.ApplyPagination(ctx, 1, 10)
	pagination.SetTotal(100)

	pageInfo := pagination.GetPageInfo()

	if pageInfo["from"] != 0 {
		t.Errorf("Expected from 0, got %v", pageInfo["from"])
	}

	if pageInfo["to"] != 10 {
		t.Errorf("Expected to 10, got %v", pageInfo["to"])
	}
}

// TestSetTotal tests setting total count
func TestSetTotal(t *testing.T) {
	field := NewBelongsTo("Author", "author_id", "authors")
	pagination := NewRelationshipPagination(field)

	pagination.SetTotal(50)

	if pagination.GetTotal() != 50 {
		t.Errorf("Expected total 50, got %d", pagination.GetTotal())
	}
}

// TestGetPage tests getting page
func TestGetPage(t *testing.T) {
	field := NewBelongsTo("Author", "author_id", "authors")
	pagination := NewRelationshipPagination(field)

	ctx := context.Background()
	pagination.ApplyPagination(ctx, 3, 10)

	if pagination.GetPage() != 3 {
		t.Errorf("Expected page 3, got %d", pagination.GetPage())
	}
}

// TestGetPerPage tests getting per-page count
func TestGetPerPage(t *testing.T) {
	field := NewBelongsTo("Author", "author_id", "authors")
	pagination := NewRelationshipPagination(field)

	ctx := context.Background()
	pagination.ApplyPagination(ctx, 1, 25)

	if pagination.GetPerPage() != 25 {
		t.Errorf("Expected perPage 25, got %d", pagination.GetPerPage())
	}
}
