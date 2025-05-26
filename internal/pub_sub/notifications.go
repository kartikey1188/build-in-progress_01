package pub_sub

import (
	"fmt"
	"os"

	"github.com/go-mail/mail"
	"github.com/kartikey1188/build-in-progress_01/internal/types"
)

func sendNotification(recipient types.Notifiable, heading string, message string) error {
	err := gmail(recipient, heading, message)
	if err != nil {
		return fmt.Errorf("failed to send gmail: %w", err)
	}

	return nil
}

func gmail(recipient types.Notifiable, heading string, message string) error {
	m := mail.NewMessage()
	m.SetHeader("From", os.Getenv("SMTP_USERNAME"))
	m.SetHeader("To", recipient.GetEmail())
	m.SetHeader("Subject", heading)

	body := fmt.Sprint(message)
	m.SetBody("text/plain", body)

	// Configuring SMTP server credentials
	d := mail.NewDialer("smtp.gmail.com", 587, os.Getenv("SMTP_USERNAME"), os.Getenv("SMTP_PASSWORD"))

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
