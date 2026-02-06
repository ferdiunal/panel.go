package blog

import (
	"time"

	"gorm.io/gorm"
)

type Author struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name"`
	Email     string         `json:"email"`
	Profile   *Profile       `json:"profile"`
	Posts     []Post         `json:"posts"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

type Profile struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	AuthorID  *uint          `json:"authorId"`
	Author    *Author        `json:"author"`
	Bio       string         `json:"bio"`
	Website   string         `json:"website"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

type Post struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Title     string         `json:"title"`
	Content   string         `json:"content"`
	AuthorID  uint           `json:"authorId"`
	Author    *Author        `json:"author"`
	Tags      []*Tag         `json:"tags" gorm:"many2many:post_tags;"`
	Comments  []Comment      `json:"comments" gorm:"polymorphic:Commentable;"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

type Tag struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name"`
	Posts     []*Post        `json:"posts" gorm:"many2many:post_tags;"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

type Comment struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	Content         string         `json:"content"`
	CommentableID   uint           `json:"commentableId"`
	CommentableType string         `json:"commentableType"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}
