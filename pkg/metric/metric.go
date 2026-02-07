package metric

import (
	"fmt"
	"time"

	"github.com/ferdiunal/panel.go/pkg/context"
	"github.com/ferdiunal/panel.go/pkg/widget"
	"gorm.io/gorm"
)

// Format types for metrics
type Format string

const (
	FormatNumber     Format = "number"
	FormatCurrency   Format = "currency"
	FormatPercentage Format = "percentage"
)

// PartitionMetric represents a pie/donut chart metric
type PartitionMetric struct {
	widget.BaseCard
	QueryFunc  func(db *gorm.DB) (map[string]int64, error)
	Colors     map[string]string
	FormatType Format
}

// NewPartition creates a new partition metric
func NewPartition(title string) *PartitionMetric {
	return &PartitionMetric{
		BaseCard: widget.BaseCard{
			TitleStr:     title,
			ComponentStr: "partition-metric",
			WidthStr:     "1/3",
			CardTypeVal:  "partition",
		},
		Colors:     make(map[string]string),
		FormatType: FormatNumber,
	}
}

// Query sets the query function
func (m *PartitionMetric) Query(fn func(db *gorm.DB) (map[string]int64, error)) *PartitionMetric {
	m.QueryFunc = fn
	return m
}

// SetColors sets custom colors for segments
func (m *PartitionMetric) SetColors(colors map[string]string) *PartitionMetric {
	m.Colors = colors
	return m
}

// SetFormat sets the display format
func (m *PartitionMetric) SetFormat(format Format) *PartitionMetric {
	m.FormatType = format
	return m
}

// SetWidth sets the card width
func (m *PartitionMetric) SetWidth(width string) *PartitionMetric {
	m.WidthStr = width
	return m
}

// Resolve executes the query and returns the data
func (m *PartitionMetric) Resolve(ctx *context.Context, db *gorm.DB) (interface{}, error) {
	if m.QueryFunc == nil {
		return nil, fmt.Errorf("query function not defined")
	}

	data, err := m.QueryFunc(db)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"data":   data,
		"colors": m.Colors,
		"format": m.FormatType,
	}, nil
}

// JsonSerialize serializes the metric for JSON response
func (m *PartitionMetric) JsonSerialize() map[string]interface{} {
	return map[string]interface{}{
		"component": m.Component(),
		"title":     m.Name(),
		"width":     m.Width(),
		"type":      m.GetType(),
		"format":    m.FormatType,
		"colors":    m.Colors,
	}
}

// ProgressMetric represents a progress bar metric
type ProgressMetric struct {
	widget.BaseCard
	CurrentFunc func(db *gorm.DB) (int64, error)
	Target      int64
	FormatType  Format
}

// NewProgress creates a new progress metric
func NewProgress(title string, target int64) *ProgressMetric {
	return &ProgressMetric{
		BaseCard: widget.BaseCard{
			TitleStr:     title,
			ComponentStr: "progress-metric",
			WidthStr:     "1/3",
			CardTypeVal:  "progress",
		},
		Target:     target,
		FormatType: FormatNumber,
	}
}

// Current sets the function to get current value
func (m *ProgressMetric) Current(fn func(db *gorm.DB) (int64, error)) *ProgressMetric {
	m.CurrentFunc = fn
	return m
}

// SetFormat sets the display format
func (m *ProgressMetric) SetFormat(format Format) *ProgressMetric {
	m.FormatType = format
	return m
}

// SetWidth sets the card width
func (m *ProgressMetric) SetWidth(width string) *ProgressMetric {
	m.WidthStr = width
	return m
}

// Resolve executes the query and returns the data
func (m *ProgressMetric) Resolve(ctx *context.Context, db *gorm.DB) (interface{}, error) {
	if m.CurrentFunc == nil {
		return nil, fmt.Errorf("current function not defined")
	}

	current, err := m.CurrentFunc(db)
	if err != nil {
		return nil, err
	}

	percentage := float64(0)
	if m.Target > 0 {
		percentage = (float64(current) / float64(m.Target)) * 100
	}

	return map[string]interface{}{
		"current":    current,
		"target":     m.Target,
		"percentage": percentage,
		"format":     m.FormatType,
	}, nil
}

// JsonSerialize serializes the metric for JSON response
func (m *ProgressMetric) JsonSerialize() map[string]interface{} {
	return map[string]interface{}{
		"component": m.Component(),
		"title":     m.Name(),
		"width":     m.Width(),
		"type":      m.GetType(),
		"target":    m.Target,
		"format":    m.FormatType,
	}
}

// TableMetric represents a table metric
type TableMetric struct {
	widget.BaseCard
	QueryFunc func(db *gorm.DB) ([]map[string]interface{}, error)
	Columns   []TableColumn
}

// TableColumn defines a column in the table metric
type TableColumn struct {
	Key   string
	Label string
	Width string
}

// NewTable creates a new table metric
func NewTable(title string) *TableMetric {
	return &TableMetric{
		BaseCard: widget.BaseCard{
			TitleStr:     title,
			ComponentStr: "table-metric",
			WidthStr:     "full",
			CardTypeVal:  widget.CardTypeTable,
		},
		Columns: []TableColumn{},
	}
}

// Query sets the query function
func (m *TableMetric) Query(fn func(db *gorm.DB) ([]map[string]interface{}, error)) *TableMetric {
	m.QueryFunc = fn
	return m
}

// SetColumns sets the table columns
func (m *TableMetric) SetColumns(columns []TableColumn) *TableMetric {
	m.Columns = columns
	return m
}

// AddColumn adds a column to the table
func (m *TableMetric) AddColumn(key, label, width string) *TableMetric {
	m.Columns = append(m.Columns, TableColumn{
		Key:   key,
		Label: label,
		Width: width,
	})
	return m
}

// SetWidth sets the card width
func (m *TableMetric) SetWidth(width string) *TableMetric {
	m.WidthStr = width
	return m
}

// Resolve executes the query and returns the data
func (m *TableMetric) Resolve(ctx *context.Context, db *gorm.DB) (interface{}, error) {
	if m.QueryFunc == nil {
		return nil, fmt.Errorf("query function not defined")
	}

	data, err := m.QueryFunc(db)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"data":    data,
		"columns": m.Columns,
	}, nil
}

// JsonSerialize serializes the metric for JSON response
func (m *TableMetric) JsonSerialize() map[string]interface{} {
	return map[string]interface{}{
		"component": m.Component(),
		"title":     m.Name(),
		"width":     m.Width(),
		"type":      m.GetType(),
		"columns":   m.Columns,
	}
}

// TrendPoint represents a point in a trend chart
type TrendPoint struct {
	Date  time.Time `json:"date"`
	Value int64     `json:"value"`
}
