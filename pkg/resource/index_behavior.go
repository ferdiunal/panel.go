package resource

import "strings"

// IndexRowClickAction defines which modal should open when a row is clicked.
type IndexRowClickAction string

const (
	IndexRowClickActionEdit   IndexRowClickAction = "edit"
	IndexRowClickActionDetail IndexRowClickAction = "detail"
)

// IndexPaginationType defines which pagination UI should be used on index pages.
type IndexPaginationType string

const (
	// IndexPaginationTypeLinks renders classic pagination with page numbers.
	IndexPaginationTypeLinks IndexPaginationType = "links"
	// IndexPaginationTypeSimple renders previous/next controls only.
	IndexPaginationTypeSimple IndexPaginationType = "simple"
	// IndexPaginationTypeLoadMore renders an incremental "load more" action.
	IndexPaginationTypeLoadMore IndexPaginationType = "load_more"
)

// IndexReorderConfig defines drag-drop reorder behavior for index tables.
type IndexReorderConfig struct {
	Enabled bool   `json:"enabled"`
	Column  string `json:"column"`
}

func NormalizeIndexRowClickAction(action IndexRowClickAction) IndexRowClickAction {
	switch strings.ToLower(strings.TrimSpace(string(action))) {
	case string(IndexRowClickActionDetail):
		return IndexRowClickActionDetail
	default:
		return IndexRowClickActionEdit
	}
}

func NormalizeIndexReorderColumn(column string) string {
	return strings.TrimSpace(column)
}

func NormalizeIndexPaginationType(paginationType IndexPaginationType) IndexPaginationType {
	switch strings.ToLower(strings.TrimSpace(string(paginationType))) {
	case string(IndexPaginationTypeSimple):
		return IndexPaginationTypeSimple
	case string(IndexPaginationTypeLoadMore), "load-more", "loadmore":
		return IndexPaginationTypeLoadMore
	default:
		return IndexPaginationTypeLinks
	}
}
