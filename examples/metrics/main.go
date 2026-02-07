package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/metric"
	"github.com/ferdiunal/panel.go/pkg/panel"
	"github.com/ferdiunal/panel.go/pkg/resource"
	"github.com/ferdiunal/panel.go/pkg/widget"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Order model
type Order struct {
	ID        uint      `gorm:"primaryKey"`
	Amount    float64   `gorm:"type:decimal(10,2)"`
	Status    string    `gorm:"size:50"`
	CreatedAt time.Time `gorm:"index"`
}

// OrderResource demonstrates the Metrics System
type OrderResource struct {
	resource.Base
}

func NewOrderResource() *OrderResource {
	r := &OrderResource{}
	r.DataModel = &Order{}
	r.Identifier = "orders"
	r.Label = "Orders"
	r.IconName = "shopping-cart"
	r.GroupName = "Sales"

	// Define fields
	r.FieldsVal = []fields.Element{
		fields.ID(),
		fields.Number("Amount", "amount").Required(),
		fields.Select("Status", "status").
			Options(map[string]string{
				"pending":   "Pending",
				"completed": "Completed",
				"cancelled": "Cancelled",
			}).
			Default("pending"),
		fields.DateTime("Created At", "created_at"),
	}

	// Define metrics/cards
	r.WidgetsVal = []widget.Card{
		// 1. Value Metric - Total Orders
		&widget.Value{
			Title: "Total Orders",
			QueryFunc: func(ctx *context.Context, db *gorm.DB) (int64, error) {
				var count int64
				err := db.Model(&Order{}).Count(&count).Error
				return count, err
			},
		},

		// 2. Trend Metric - Orders Over Time
		widget.NewTrendWidget("Orders Trend", &Order{}, "created_at"),

		// 3. Partition Metric - Orders by Status
		metric.NewPartition("Orders by Status").
			SetWidth("1/2").
			Query(func(db *gorm.DB) (map[string]int64, error) {
				return metric.GroupByColumn(db, &Order{}, "status")
			}).
			SetColors(map[string]string{
				"pending":   "#f59e0b",
				"completed": "#10b981",
				"cancelled": "#ef4444",
			}),

		// 4. Progress Metric - Monthly Goal
		metric.NewProgress("Monthly Goal", 1000).
			SetWidth("1/3").
			Current(func(db *gorm.DB) (int64, error) {
				startOfMonth := time.Now().AddDate(0, 0, -time.Now().Day()+1)
				return metric.CountWhere(db, &Order{}, "created_at >= ?", startOfMonth)
			}),

		// 5. Table Metric - Recent Orders
		metric.NewTable("Recent Orders").
			SetWidth("full").
			AddColumn("id", "ID", "80px").
			AddColumn("amount", "Amount", "120px").
			AddColumn("status", "Status", "120px").
			AddColumn("created_at", "Date", "150px").
			Query(func(db *gorm.DB) ([]map[string]interface{}, error) {
				var orders []Order
				err := db.Model(&Order{}).
					Order("created_at DESC").
					Limit(10).
					Find(&orders).Error

				if err != nil {
					return nil, err
				}

				// Convert to map slice
				result := make([]map[string]interface{}, len(orders))
				for i, order := range orders {
					result[i] = map[string]interface{}{
						"id":         order.ID,
						"amount":     fmt.Sprintf("$%.2f", order.Amount),
						"status":     order.Status,
						"created_at": order.CreatedAt.Format("2006-01-02 15:04"),
					}
				}
				return result, nil
			}),

		// 6. Value Metric with Currency Format
		&widget.Value{
			Title: "Total Revenue",
			QueryFunc: func(ctx *context.Context, db *gorm.DB) (int64, error) {
				var total float64
				err := db.Model(&Order{}).
					Where("status = ?", "completed").
					Select("COALESCE(SUM(amount), 0)").
					Scan(&total).Error
				return int64(total), err
			},
		},

		// 7. Progress Metric - Completion Rate
		metric.NewProgress("Completion Rate", 100).
			SetWidth("1/3").
			SetFormat(metric.FormatPercentage).
			Current(func(db *gorm.DB) (int64, error) {
				var total, completed int64
				db.Model(&Order{}).Count(&total)
				db.Model(&Order{}).Where("status = ?", "completed").Count(&completed)

				if total == 0 {
					return 0, nil
				}
				return (completed * 100) / total, nil
			}),
	}

	return r
}

func main() {
	// Setup database
	db, err := gorm.Open(sqlite.Open("metrics_example.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Auto migrate
	db.AutoMigrate(&Order{})

	// Seed some data
	orders := []Order{
		{Amount: 99.99, Status: "completed", CreatedAt: time.Now().AddDate(0, 0, -5)},
		{Amount: 149.99, Status: "completed", CreatedAt: time.Now().AddDate(0, 0, -4)},
		{Amount: 79.99, Status: "pending", CreatedAt: time.Now().AddDate(0, 0, -3)},
		{Amount: 199.99, Status: "completed", CreatedAt: time.Now().AddDate(0, 0, -2)},
		{Amount: 59.99, Status: "cancelled", CreatedAt: time.Now().AddDate(0, 0, -1)},
		{Amount: 129.99, Status: "pending", CreatedAt: time.Now()},
	}
	for _, order := range orders {
		db.FirstOrCreate(&order, Order{Amount: order.Amount, CreatedAt: order.CreatedAt})
	}

	// Create panel
	p := panel.New(panel.Config{
		Database: panel.DatabaseConfig{
			Instance: db,
		},
		Server: panel.ServerConfig{
			Host: "localhost",
			Port: "3000",
		},
		Resources: []resource.Resource{
			NewOrderResource(),
		},
	})

	fmt.Println("Metrics System Example")
	fmt.Println("======================")
	fmt.Println("Server running on http://localhost:3000")
	fmt.Println("")
	fmt.Println("Available Metrics:")
	fmt.Println("1. Value Metric - Total Orders count")
	fmt.Println("2. Trend Metric - Orders over time chart")
	fmt.Println("3. Partition Metric - Orders by status (pie chart)")
	fmt.Println("4. Progress Metric - Monthly goal progress bar")
	fmt.Println("5. Table Metric - Recent orders table")
	fmt.Println("6. Value Metric - Total revenue")
	fmt.Println("7. Progress Metric - Completion rate percentage")
	fmt.Println("")
	fmt.Println("Try it out:")
	fmt.Println("1. Go to http://localhost:3000")
	fmt.Println("2. Navigate to Orders")
	fmt.Println("3. View the dashboard with all metrics")

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
