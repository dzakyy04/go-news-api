package entity

type Category struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"type:varchar(100);not null" json:"name"`
	Description string `gorm:"type:text;not null" json:"description"`
}
