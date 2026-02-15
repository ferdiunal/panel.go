package resource

import "strings"

// IndexRowClickAction defines which modal should open when a row is clicked.
type IndexRowClickAction string

const (
	IndexRowClickActionEdit   IndexRowClickAction = "edit"
	IndexRowClickActionDetail IndexRowClickAction = "detail"
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
