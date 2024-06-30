package utils

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/smtp"
	"os"
	"text/template"

	"github.com/joho/godotenv"
)

func GenerateOTP(length int) string {
	var otp string

	for i := 0; i < length; i++ {
		digit := rand.Intn(10)
		otp += fmt.Sprintf("%d", digit)
	}

	return otp
}

func SendEmail(to string, subject string, templateFile string, data interface{}) error {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error loading .env file: %w", err)
	}

	from := os.Getenv("MAIL_FROM")
	password := os.Getenv("MAIL_PASSWORD")
	smtpHost := os.Getenv("MAIL_HOST")
	smtpPort := os.Getenv("MAIL_PORT")

	// Authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Parse HTML template
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		return fmt.Errorf("error parsing template %s: %w", templateFile, err)
	}

	// Prepare email body
	var body bytes.Buffer
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: %s\n%s\n\n", subject, mimeHeaders)))

	// Execute template with dynamic data
	err = tmpl.Execute(&body, data)
	if err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	// Sending email
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, body.Bytes())
	if err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}

	return nil
}
