package pub_sub

import (
	"fmt"
	"os"

	"github.com/go-mail/mail"
	"github.com/kartikey1188/build-in-progress_01/internal/types"
)

func sendNotification(recipient types.Collector, pr types.PickupRequest) error {
	err := gmail(recipient, pr)
	if err != nil {
		return fmt.Errorf("failed to send gmail: %w", err)
	}

	return nil
}

func gmail(recipient types.Collector, pr types.PickupRequest) error {
	m := mail.NewMessage()
	m.SetHeader("From", os.Getenv("SMTP_USERNAME"))
	m.SetHeader("To", recipient.Email)
	m.SetHeader("Subject", "New Pickup Request Notification")

	body := fmt.Sprintf("Dear Collector,\n\nYou have a new pickup request (ID: %d).\nBusiness ID: %d\nWaste Type: %s\nQuantity: %.2f\n\nRegards,\nXphora AI ",
		pr.RequestID, pr.BusinessID, pr.WasteType, pr.Quantity)
	m.SetBody("text/plain", body)

	// Configuring SMTP server credentials
	d := mail.NewDialer("smtp.gmail.com", 587, os.Getenv("SMTP_USERNAME"), os.Getenv("SMTP_PASSWORD"))

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
