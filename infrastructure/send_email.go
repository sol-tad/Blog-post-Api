package infrastructure

import (
	"fmt"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendOTP sends a One-Time Password (OTP) email to the specified recipient using SendGrid
func SendOTP(toEmail, otp string) error {
	// Sender email address fetched from environment variables
	from := mail.NewEmail("blogapp", os.Getenv("SENDER_EMAIL")) 
	
	// Subject of the email
	subject := "Verify your BlogApp account"
	
	// Recipient email address
	to := mail.NewEmail("", toEmail)

	// Plain text content of the email
	plainText := fmt.Sprintf("Your OTP is: %s", otp)
	
	// HTML content of the email with the OTP in bold
	htmlText := fmt.Sprintf("<strong>Your OTP is: %s</strong>", otp)

	// Create a new single email message
	message := mail.NewSingleEmail(from, subject, to, plainText, htmlText)

	// Create a new SendGrid client using the API key from environment variables
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	
	// Send the email
	response, err := client.Send(message)

	// Log the SendGrid response for debugging purposes
	fmt.Println("SendGrid Status Code:", response.StatusCode)
	fmt.Println("SendGrid Response Body:", response.Body)
	fmt.Println("SendGrid Headers:", response.Headers)

	// Return error if sending the email failed
	if err != nil {
		return err
	}

	// If the response status code indicates failure (400 or above), return an error
	if response.StatusCode >= 400 {
		return fmt.Errorf("sendgrid failed: %v", response.Body)
	}
	
	// Return nil if email was sent successfully
	return nil
}
