package widget

import (
	"fmt"
	"time"

	"github.com/ferdiunal/panel.go/internal/context"
	"gorm.io/gorm"
)

type Trend struct {
	Title     string
	Ranges    []int
	QueryFunc func(ctx *context.Context, db *gorm.DB) ([]interface{}, error)
}

func (w *Trend) Name() string {
	return w.Title
}

func (w *Trend) Component() string {
	return "trend-metric"
}

func (w *Trend) Width() string {
	return "1/3" // Trend defaults to 1/3, could be configurable
}

func (w *Trend) Resolve(ctx *context.Context, db *gorm.DB) (interface{}, error) {
	data, err := w.QueryFunc(ctx, db)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"data":  data,
		"title": w.Title,
	}, nil
}

func (w *Trend) JsonSerialize() map[string]interface{} {
	return map[string]interface{}{
		"component": "trend-metric",
		"title":     w.Title,
		"width":     "1/3",
		"ranges":    w.Ranges,
	}
}

// TrendValue represents a single data point in the trend chart
type TrendValue struct {
	Date  string `json:"date"`
	Value int64  `json:"value"`
}

// Helper to fill date gaps and return formatted data
func fillGaps(results []TrendValue, days int) []map[string]interface{} {
	now := time.Now()
	dateMap := make(map[string]int64)

	// Populate map with query results
	for _, res := range results {
		dateMap[res.Date] = res.Value
	}

	finalData := make([]map[string]interface{}, 0)

	// Iterate backwards from today for 'days' count
	// Chart usually expects chronological order (oldest to newest)
	// So let's generate dates from (now - days) to now.
	start := now.AddDate(0, 0, -days+1)

	for i := 0; i < days; i++ {
		d := start.AddDate(0, 0, i)
		dateStr := d.Format("2006-01-02")

		val := int64(0)
		if v, ok := dateMap[dateStr]; ok {
			val = v
		}

		finalData = append(finalData, map[string]interface{}{
			"date":  dateStr, // Optional: might be useful for tooltip
			"value": val,
		})
	}

	return finalData
}

// NewTrendWidget creates a trend widget that groups count by day for the last 30 days.
// Note: This currently assumes SQLite syntax for date grouping.
func NewTrendWidget(title string, model interface{}, dateColumn string) *Trend {
	return &Trend{
		Title:  title,
		Ranges: []int{30, 60, 90},
		QueryFunc: func(ctx *context.Context, db *gorm.DB) ([]interface{}, error) {
			// Get range from query param, default to 30
			days := ctx.QueryInt("range", 30)

			// Validate range against allowed ranges
			valid := false
			allowedRanges := []int{30, 60, 90}
			for _, r := range allowedRanges {
				if r == days {
					valid = true
					break
				}
			}
			if !valid {
				days = 30
			}

			endDate := time.Now()
			startDate := endDate.AddDate(0, 0, -days)

			var results []TrendValue

			// SQLite strftime('%Y-%m-%d', column)
			// Ensure dateColumn is safe or validated if user input (here it's code definition)
			dateExpr := fmt.Sprintf("strftime('%%Y-%%m-%%d', %s)", dateColumn)

			err := db.Model(model).
				Select(fmt.Sprintf("%s as date, count(*) as value", dateExpr)).
				Where(fmt.Sprintf("%s BETWEEN ? AND ?", dateColumn), startDate, endDate).
				Group("date").
				Order("date ASC").
				Scan(&results).Error

			if err != nil {
				return nil, err
			}

			// Fill gaps
			filled := fillGaps(results, days)

			// Convert to []interface{} required by signature
			interfaceSlice := make([]interface{}, len(filled))
			for i, v := range filled {
				interfaceSlice[i] = v
			}

			return interfaceSlice, nil
		},
	}
}
