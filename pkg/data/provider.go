package data

import (
	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/query"
)

type Sort struct {
	Column    string `json:"column"`
	Direction string `json:"direction"`
}

type QueryRequest struct {
	Page    int            `json:"page"`
	PerPage int            `json:"per_page"`
	Sorts   []Sort         `json:"sorts"`
	Filters []query.Filter `json:"filters"`
	Search  string         `json:"search"`
}

type QueryResponse struct {
	Items   []interface{} `json:"items"`
	Total   int64         `json:"total"`
	Page    int           `json:"page"`
	PerPage int           `json:"per_page"`
}

type DataProvider interface {
	Index(ctx *context.Context, req QueryRequest) (*QueryResponse, error)
	Show(ctx *context.Context, id string) (interface{}, error)
	Create(ctx *context.Context, data map[string]interface{}) (interface{}, error)
	Update(ctx *context.Context, id string, data map[string]interface{}) (interface{}, error)
	Delete(ctx *context.Context, id string) error
	SetSearchColumns(cols []string)
	SetWith(rels []string)
}
