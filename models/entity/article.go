package entity

import "time"

type ArticleStatus string

const (
	Published ArticleStatus = "published"
	Draft     ArticleStatus = "draft"
	Review    ArticleStatus = "review"
)

type Article struct {
	ID         uint          `gorm:"primaryKey" json:"id"`
	Title      string        `gorm:"type:varchar(100);not null" json:"title"`
	Slug       string        `gorm:"type:varchar(100);not null;unique" json:"slug"`
	Thumbnail  string        `gorm:"type:varchar(100);not null" json:"thumbnail"`
	Content    string        `gorm:"type:text;not null" json:"content"`
	CategoryID uint          `gorm:"not null" json:"category_id"`
	Category   Category      `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"category"`
	Tags       []Tag         `gorm:"many2many:article_tags;" json:"tags"`
	Author     User          `gorm:"foreignKey:AuthorID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"author"`
	AuthorID   uint          `gorm:"not null" json:"author_id"`
	Status     ArticleStatus `gorm:"type:enum('draft', 'published', 'archived');default:draft" json:"status"`
	Comments   []Comment     `gorm:"foreignKey:ArticleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"comments"`
	CreatedAt  time.Time     `json:"created_at"`
	UpdatedAt  time.Time     `json:"updated_at"`
}
