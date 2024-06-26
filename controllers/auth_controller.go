package controllers

import (
	"errors"
	"go-news-api/database"
	"go-news-api/models/entity"
	"go-news-api/models/request"
	"go-news-api/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func Register(ctx *fiber.Ctx) error {
	request := new(request.RegisterRequest)

	// Parse request body
	if err := ctx.BodyParser(request); err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusBadRequest, "Failed to register", err)
	}

	// Validate request
	if err := utils.Validate.Struct(request); err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusBadRequest, "Failed to register", err)
	}

	// Create user
	user := entity.User{
		Name:  request.Name,
		Email: request.Email,
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to register", err)
	}
	user.Password = hashedPassword

	// Save user
	if err := database.DB.Create(&user).Error; err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to register", err)
	}

	return utils.SendSuccessResponse(ctx, fiber.StatusCreated, "Sucessfully registered", fiber.Map{
		"user": user,
	})
}

func Login(ctx *fiber.Ctx) error {
	request := new(request.LoginRequest)

	// Parse request body
	if err := ctx.BodyParser(request); err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusBadRequest, "Failed to login", err)
	}

	// Validate request
	if err := utils.Validate.Struct(request); err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusBadRequest, "Failed to login", err)
	}

	// Find user
	var user entity.User
	if err := database.DB.Where("email = ?", request.Email).First(&user).Error; err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusUnauthorized, "Failed to login", err)
	}

	// Check password
	if err := utils.VerifyPassword(request.Password, user.Password); err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusUnauthorized, "Failed to login", err)
	}

	// Generate JWT token
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	token, err := utils.GenerateToken(&claims)
	if err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusNotFound, "Failed to login", err)
	}

	return utils.SendSuccessResponse(ctx, fiber.StatusOK, "Successfully logged in", fiber.Map{
		"token": token,
		"user":  user,
	})
}

func SendVerificationEmail(ctx *fiber.Ctx) error {
	request := new(request.SendVerificationEmailRequest)

	// Parse request body
	if err := ctx.BodyParser(request); err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusBadRequest, "Failed to send verification email", err)
	}

	// Validate request
	if err := utils.Validate.Struct(request); err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusBadRequest, "Failed to send verification email", err)
	}

	// Find user
	var user entity.User
	if err := database.DB.Where("email = ?", request.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.SendErrorResponse(ctx, fiber.StatusNotFound, "Failed to send verification email", err)
		}
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to send verification email", err)
	}

	// Check if email is already verified
	if user.IsVerified {
		return utils.SendErrorResponse(ctx, fiber.StatusBadRequest, "Failed to send verification email", errors.New("email is already verified"))
	}

	// Generate OTP
	otp := utils.GenerateOTP(4)
	otpExpiredAt := time.Now().Add(time.Minute * 10)

	// Check if OTP already exists for the user
	var existingOtp entity.OtpCode
	err := database.DB.Where("user_id = ? AND type = ?", user.ID, entity.EmailVerification).First(&existingOtp).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to send verification email", err)
	}

	if err == gorm.ErrRecordNotFound {
		// Create new OTP
		otpCode := entity.OtpCode{
			Otp:       otp,
			ExpiredAt: otpExpiredAt,
			Type:      entity.EmailVerification,
			UserID:    user.ID,
		}

		if err := database.DB.Create(&otpCode).Error; err != nil {
			return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to send verification email", err)
		}
	} else {
		// Update existing OTP
		existingOtp.Otp = otp
		existingOtp.ExpiredAt = otpExpiredAt

		if err := database.DB.Save(&existingOtp).Error; err != nil {
			return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to send verification email", err)
		}
	}

	// Send email
	if err := utils.SendEmail(user.Email, "Verify your email", "views/emails/verification.html", fiber.Map{
		"name": user.Name,
		"otp":  otp,
	}); err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to send verification email", err)
	}

	return utils.SendSuccessResponse(ctx, fiber.StatusOK, "Successfully sent verification email", nil)
}

func VerifyEmail(ctx *fiber.Ctx) error {
	request := new(request.VerifyEmailRequest)

	// Parse request body
	if err := ctx.BodyParser(request); err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusBadRequest, "Failed to verify email", err)
	}

	// Check if user exists
	var user entity.User
	err := database.DB.Where("email = ?", request.Email).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.SendErrorResponse(ctx, fiber.StatusNotFound, "Failed to verify email", err)
		}
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to verify email", err)
	}

	// Check if email is already verified
	if user.IsVerified {
		return utils.SendErrorResponse(ctx, fiber.StatusBadRequest, "Failed to verify email", errors.New("email is already verified"))
	}

	// Check if OTP is valid
	var otpCode entity.OtpCode
	err = database.DB.Where("user_id = ? AND type = ?", user.ID, entity.EmailVerification).First(&otpCode).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.SendErrorResponse(ctx, fiber.StatusUnauthorized, "Failed to verify email", errors.New("invalid or expired OTP code"))
		}
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to verify email", err)
	}

	if otpCode.Otp != request.Otp || time.Now().After(otpCode.ExpiredAt) {
		return utils.SendErrorResponse(ctx, fiber.StatusUnauthorized, "Failed to verify email", errors.New("invalid or expired OTP code"))
	}

	// Update user's email verification status
	err = database.DB.Model(&user).Updates(map[string]interface{}{
		"is_verified": true,
	}).Error
	if err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to verify email", err)
	}

	// Remove OTP code after successful verification
	if err := database.DB.Delete(&otpCode).Error; err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to verify email", err)
	}

	return utils.SendSuccessResponse(ctx, fiber.StatusOK, "Email has been verified", fiber.Map{
		"user": user,
	})
}

func GetProfile(ctx *fiber.Ctx) error {
	// Get user from context
	user := ctx.Locals("user").(*entity.User)

	if user == nil {
		return utils.SendErrorResponse(ctx, fiber.StatusUnauthorized, "Unauthorized", errors.New("user not found"))
	}

	return utils.SendSuccessResponse(ctx, fiber.StatusOK, "Successfully get profile", fiber.Map{
		"user": user,
	})
}

func SendResetPasswordEmail(ctx *fiber.Ctx) error {
	request := new(request.SendResetPasswordEmailRequest)

	// Parse request body
	if err := ctx.BodyParser(request); err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusBadRequest, "Failed to send reset password email", err)
	}

	// Validate request
	if err := utils.Validate.Struct(request); err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusBadRequest, "Failed to send reset password email", err)
	}

	// Find user
	var user entity.User
	if err := database.DB.Where("email = ?", request.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.SendErrorResponse(ctx, fiber.StatusNotFound, "Failed to send reset password email", err)
		}
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to send reset password email", err)
	}

	// Generate OTP
	otp := utils.GenerateOTP(4)
	otpExpiredAt := time.Now().Add(time.Minute * 10)

	// Check if OTP already exists for the user
	var existingOtp entity.OtpCode
	err := database.DB.Where("user_id = ? AND type = ?", user.ID, entity.PasswordReset).First(&existingOtp).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to send reset password email", err)
	}

	if err == gorm.ErrRecordNotFound {
		// Create new OTP
		otpCode := entity.OtpCode{
			Otp:       otp,
			ExpiredAt: otpExpiredAt,
			Type:      entity.PasswordReset,
			UserID:    user.ID,
		}

		if err := database.DB.Create(&otpCode).Error; err != nil {
			return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to send reset password email", err)
		}
	} else {
		// Update existing OTP
		existingOtp.Otp = otp
		existingOtp.ExpiredAt = otpExpiredAt

		if err := database.DB.Save(&existingOtp).Error; err != nil {
			return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to send reset password email", err)
		}
	}

	// Send email
	if err := utils.SendEmail(user.Email, "Reset your password", "views/emails/reset_password.html", fiber.Map{
		"name": user.Name,
		"otp":  otp,
	}); err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to send reset password email", err)
	}

	return utils.SendSuccessResponse(ctx, fiber.StatusOK, "Successfully sent reset password email", nil)
}

func VerifyOtpReset(ctx *fiber.Ctx) error {
	request := new(request.VerifyOtpResetRequest)

	// Parse request body
	if err := ctx.BodyParser(request); err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusBadRequest, "Failed to verify OTP", err)
	}

	// Validate request
	if err := utils.Validate.Struct(request); err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusBadRequest, "Failed to verify OTP", err)
	}

	// Find user
	var user entity.User
	if err := database.DB.Where("email = ?", request.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.SendErrorResponse(ctx, fiber.StatusNotFound, "Failed to verify OTP", err)
		}
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to verify OTP", err)
	}

	// Find OTP
	var otpCode entity.OtpCode
	if err := database.DB.Where("user_id = ? AND type = ?", user.ID, entity.PasswordReset).First(&otpCode).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.SendErrorResponse(ctx, fiber.StatusNotFound, "Failed to verify OTP", err)
		}
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to verify OTP", err)
	}

	// Check if OTP is already verified
	if otpCode.IsVerified {
		return utils.SendErrorResponse(ctx, fiber.StatusBadRequest, "Failed to verify OTP", errors.New("otp is already verified"))
	}

	// Check if OTP is valid
	if otpCode.Otp != request.Otp || time.Now().After(otpCode.ExpiredAt) {
		return utils.SendErrorResponse(ctx, fiber.StatusUnauthorized, "Failed to verify OTP", errors.New("invalid or expired OTP code"))
	}

	// Set OTP to verified
	err := database.DB.Model(&otpCode).Updates(map[string]interface{}{
		"is_verified": true,
	}).Error
	if err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to verify OTP", err)
	}

	return utils.SendSuccessResponse(ctx, fiber.StatusOK, "Successfully verified OTP", nil)
}

func ResetPassword(ctx *fiber.Ctx) error {
	request := new(request.ResetPasswordRequest)

	// Parse request body
	if err := ctx.BodyParser(request); err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusBadRequest, "Failed to reset password", err)
	}

	// Validate request
	if err := utils.Validate.Struct(request); err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusBadRequest, "Failed to reset password", err)
	}

	// Find user
	var user entity.User
	if err := database.DB.Where("email = ?", request.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.SendErrorResponse(ctx, fiber.StatusNotFound, "Failed to reset password", err)
		}
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to reset password", err)
	}

	// Find OTP
	var otpCode entity.OtpCode
	if err := database.DB.Where("user_id = ? AND type = ?", user.ID, entity.PasswordReset).First(&otpCode).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.SendErrorResponse(ctx, fiber.StatusNotFound, "Failed to reset password", err)
		}
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to reset password", err)
	}

	// Check if OTP is already verified
	if !otpCode.IsVerified {
		return utils.SendErrorResponse(ctx, fiber.StatusBadRequest, "Failed to reset password", errors.New("otp is not verified"))
	}

	// Hash password
	newHashedPassword, err := utils.HashPassword(request.NewPassword)
	if err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to reset password", err)
	}

	// Update password
	err = database.DB.Model(&user).Update("password", newHashedPassword).Error
	if err != nil {
		return utils.SendErrorResponse(ctx, fiber.StatusInternalServerError, "Failed to reset password", err)
	}

	return utils.SendSuccessResponse(ctx, fiber.StatusOK, "Successfully reset password", nil)
}
