package request

type RegisterRequest struct {
	Name                 string `json:"name" validate:"required,min=3,max=100"`
	Email                string `json:"email" validate:"required,email"`
	Password             string `json:"password" validate:"required,min=8"`
	PasswordConfirmation string `json:"password_confirmation" form:"password_confirmation" validate:"required,min=8,eqfield=Password"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type SendVerificationEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type VerifyEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
	Otp   string `json:"otp" validate:"required"`
}

type SendResetPasswordEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type VerifyOtpResetRequest struct {
	Email string `json:"email" validate:"required,email"`
	Otp   string `json:"otp" validate:"required"`
}

type ResetPasswordRequest struct {
	Email                   string `json:"email" validate:"required,email"`
	NewPassword             string `json:"new_password" form:"new_password" validate:"required,min=8"`
	NewPasswordConfirmation string `json:"new_password_confirmation" form:"new_password_confirmation" validate:"required,min=8,eqfield=NewPassword"`
}
