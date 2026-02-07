package main

import (
	"fmt"
	"log"

	"github.com/ferdiunal/panel.go/pkg/action"
	"github.com/ferdiunal/panel.go/pkg/fields"
	"github.com/ferdiunal/panel.go/pkg/panel"
	"github.com/ferdiunal/panel.go/pkg/resource"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Post model
type Post struct {
	ID      uint   `gorm:"primaryKey"`
	Title   string `gorm:"size:255"`
	Content string `gorm:"type:text"`
	Status  string `gorm:"size:50;default:'draft'"`
}

// PostResource demonstrates the Actions System
type PostResource struct {
	resource.Base
}

func NewPostResource() *PostResource {
	r := &PostResource{}
	r.DataModel = &Post{}
	r.Identifier = "posts"
	r.Label = "Posts"
	r.IconName = "file-text"
	r.GroupName = "Content"

	// Define fields
	r.FieldsVal = []fields.Element{
		fields.ID(),
		fields.Text("Title", "title").Required().Sortable().Searchable(),
		fields.Textarea("Content", "content").Required(),
		fields.Select("Status", "status").
			Options(map[string]string{
				"draft":     "Draft",
				"published": "Published",
				"archived":  "Archived",
			}).
			Default("draft"),
	}

	// Define actions
	r.ActionsVal = []resource.Action{
		// 1. Publish Action - Simple action without fields
		action.New("Publish Posts").
			SetIcon("check-circle").
			Confirm("Are you sure you want to publish these posts?").
			Handle(func(ctx *action.ActionContext) error {
				for _, model := range ctx.Models {
					post := model.(*Post)
					post.Status = "published"
					if err := ctx.DB.Save(post).Error; err != nil {
						return err
					}
				}
				return nil
			}),

		// 2. Archive Action - Destructive action
		action.New("Archive Posts").
			SetIcon("archive").
			Destructive().
			Confirm("Are you sure you want to archive these posts? They will no longer be visible.").
			ConfirmButton("Archive").
			Handle(func(ctx *action.ActionContext) error {
				for _, model := range ctx.Models {
					post := model.(*Post)
					post.Status = "archived"
					if err := ctx.DB.Save(post).Error; err != nil {
						return err
					}
				}
				return nil
			}),

		// 3. Update Status Action - Action with fields
		action.New("Update Status").
			SetIcon("edit").
			WithFields(
				fields.Select("New Status", "status").
					Options(map[string]string{
						"draft":     "Draft",
						"published": "Published",
						"archived":  "Archived",
					}).
					Required(),
			).
			Handle(func(ctx *action.ActionContext) error {
				newStatus := ctx.Fields["status"].(string)
				for _, model := range ctx.Models {
					post := model.(*Post)
					post.Status = newStatus
					if err := ctx.DB.Save(post).Error; err != nil {
						return err
					}
				}
				return nil
			}),

		// 4. Send Notification Action - Action with multiple fields
		action.New("Send Notification").
			SetIcon("mail").
			WithFields(
				fields.Text("Subject", "subject").
					Required().
					Placeholder("Enter notification subject"),
				fields.Textarea("Message", "message").
					Required().
					Placeholder("Enter notification message"),
				fields.Switch("Send Email", "send_email").
					Default(true),
			).
			Handle(func(ctx *action.ActionContext) error {
				subject := ctx.Fields["subject"].(string)
				message := ctx.Fields["message"].(string)
				sendEmail := ctx.Fields["send_email"].(bool)

				fmt.Printf("Sending notification to %d posts:\n", len(ctx.Models))
				fmt.Printf("Subject: %s\n", subject)
				fmt.Printf("Message: %s\n", message)
				fmt.Printf("Send Email: %v\n", sendEmail)

				// Implement notification logic here
				return nil
			}),

		// 5. Export CSV Action - Built-in action
		action.ExportCSV("posts.csv"),

		// 6. Delete Action - Built-in destructive action
		action.Delete().OnlyOnIndex(),

		// 7. Conditional Action - Only runs if conditions are met
		action.New("Feature Post").
			SetIcon("star").
			CanRun(func(ctx *action.ActionContext) bool {
				// Only allow featuring published posts
				for _, model := range ctx.Models {
					post := model.(*Post)
					if post.Status != "published" {
						return false
					}
				}
				return true
			}).
			Handle(func(ctx *action.ActionContext) error {
				fmt.Printf("Featuring %d posts\n", len(ctx.Models))
				// Implement feature logic here
				return nil
			}),
	}

	return r
}

func main() {
	// Setup database
	db, err := gorm.Open(sqlite.Open("actions_example.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Auto migrate
	db.AutoMigrate(&Post{})

	// Seed some data
	posts := []Post{
		{Title: "First Post", Content: "This is the first post", Status: "draft"},
		{Title: "Second Post", Content: "This is the second post", Status: "published"},
		{Title: "Third Post", Content: "This is the third post", Status: "draft"},
	}
	for _, post := range posts {
		db.FirstOrCreate(&post, Post{Title: post.Title})
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
			NewPostResource(),
		},
	})

	fmt.Println("Actions System Example")
	fmt.Println("======================")
	fmt.Println("Server running on http://localhost:3000")
	fmt.Println("")
	fmt.Println("Available Actions:")
	fmt.Println("1. Publish Posts - Simple action to publish draft posts")
	fmt.Println("2. Archive Posts - Destructive action to archive posts")
	fmt.Println("3. Update Status - Action with field to change post status")
	fmt.Println("4. Send Notification - Action with multiple fields")
	fmt.Println("5. Export CSV - Built-in action to export posts to CSV")
	fmt.Println("6. Delete - Built-in destructive action to delete posts")
	fmt.Println("7. Feature Post - Conditional action (only for published posts)")
	fmt.Println("")
	fmt.Println("Try it out:")
	fmt.Println("1. Go to http://localhost:3000")
	fmt.Println("2. Navigate to Posts")
	fmt.Println("3. Select one or more posts using checkboxes")
	fmt.Println("4. Click the 'Actions' dropdown")
	fmt.Println("5. Select an action and fill in any required fields")
	fmt.Println("6. Confirm the action")

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
