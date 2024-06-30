package entity

import "time"

type OtpType string

const (
	EmailVerification OtpType = "email_verification"
	PasswordReset     OtpType = "password_reset"
)

type OtpCode struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Otp        string    `gorm:"type:varchar(4);not null" json:"otp"`
	Type       OtpType   `gorm:"type:enum('email_verification', 'password_reset');not null" json:"type"`
	ExpiredAt  time.Time `json:"expired_at"`
	IsVerified bool      `gorm:"default:false" json:"is_verified"`
	UserID     uint      `gorm:"not null" json:"user_id"`
	User       User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"user"`
}
