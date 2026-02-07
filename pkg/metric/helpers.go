package metric

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Result represents a date-value pair from database queries
type Result struct {
	Date  string `json:"date"`
	Value int64  `json:"value"`
}

// CountByDateRange counts records grouped by date within a date range
func CountByDateRange(db *gorm.DB, model interface{}, dateColumn string, days int) ([]TrendPoint, error) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	var results []Result

	// SQLite date formatting
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

	// Convert to TrendPoint and fill gaps
	return fillDateGaps(results, days), nil
}

// SumByDateRange sums a column grouped by date within a date range
func SumByDateRange(db *gorm.DB, model interface{}, dateColumn, sumColumn string, days int) ([]TrendPoint, error) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	var results []Result

	// SQLite date formatting
	dateExpr := fmt.Sprintf("strftime('%%Y-%%m-%%d', %s)", dateColumn)

	err := db.Model(model).
		Select(fmt.Sprintf("%s as date, COALESCE(SUM(%s), 0) as value", dateExpr, sumColumn)).
		Where(fmt.Sprintf("%s BETWEEN ? AND ?", dateColumn), startDate, endDate).
		Group("date").
		Order("date ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// Convert to TrendPoint and fill gaps
	return fillDateGaps(results, days), nil
}

// GroupByColumn groups records by a column and returns counts
func GroupByColumn(db *gorm.DB, model interface{}, column string) (map[string]int64, error) {
	type Result struct {
		Key   string `json:"key"`
		Value int64  `json:"value"`
	}

	var results []Result

	err := db.Model(model).
		Select(fmt.Sprintf("%s as key, count(*) as value", column)).
		Group(column).
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// Convert to map
	data := make(map[string]int64)
	for _, r := range results {
		data[r.Key] = r.Value
	}

	return data, nil
}

// CountWhere counts records matching a condition
func CountWhere(db *gorm.DB, model interface{}, condition string, args ...interface{}) (int64, error) {
	var count int64
	err := db.Model(model).Where(condition, args...).Count(&count).Error
	return count, err
}

// SumWhere sums a column for records matching a condition
func SumWhere(db *gorm.DB, model interface{}, column, condition string, args ...interface{}) (int64, error) {
	type Result struct {
		Total int64 `json:"total"`
	}

	var result Result
	err := db.Model(model).
		Select(fmt.Sprintf("COALESCE(SUM(%s), 0) as total", column)).
		Where(condition, args...).
		Scan(&result).Error

	return result.Total, err
}

// fillDateGaps fills missing dates with zero values
func fillDateGaps(results []Result, days int) []TrendPoint {
	now := time.Now()
	dateMap := make(map[string]int64)

	// Populate map with query results
	for _, res := range results {
		dateMap[res.Date] = res.Value
	}

	points := make([]TrendPoint, 0, days)

	// Generate dates from (now - days) to now
	start := now.AddDate(0, 0, -days+1)

	for i := 0; i < days; i++ {
		d := start.AddDate(0, 0, i)
		dateStr := d.Format("2006-01-02")

		val := int64(0)
		if v, ok := dateMap[dateStr]; ok {
			val = v
		}

		points = append(points, TrendPoint{
			Date:  d,
			Value: val,
		})
	}

	return points
}

// AverageByDateRange calculates average of a column grouped by date
func AverageByDateRange(db *gorm.DB, model interface{}, dateColumn, avgColumn string, days int) ([]TrendPoint, error) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	var results []Result

	// SQLite date formatting
	dateExpr := fmt.Sprintf("strftime('%%Y-%%m-%%d', %s)", dateColumn)

	err := db.Model(model).
		Select(fmt.Sprintf("%s as date, COALESCE(AVG(%s), 0) as value", dateExpr, avgColumn)).
		Where(fmt.Sprintf("%s BETWEEN ? AND ?", dateColumn), startDate, endDate).
		Group("date").
		Order("date ASC").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	// Convert to TrendPoint and fill gaps
	return fillDateGaps(results, days), nil
}
