package request

type RegisterRequest struct {
	Name                 string `json:"name" validate:"required,min=3,max=100"`
	Email                string `json:"email" validate:"required,email"`
	Password             string `json:"password" validate:"required,min=8"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,min=8,eqfield=Password"`
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
