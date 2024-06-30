package entity

import "time"

type User struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Name       string    `gorm:"type:varchar(100);not null" json:"name"`
	Email      string    `gorm:"type:varchar(100);not null;unique" json:"email"`
	Password   string    `gorm:"type:varchar(100);column:password;not null" json:"-"`
	IsVerified bool      `gorm:"default:false" json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
