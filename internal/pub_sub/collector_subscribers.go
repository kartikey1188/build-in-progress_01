package pub_sub

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/kartikey1188/build-in-progress_01/internal/storage"
	"github.com/kartikey1188/build-in-progress_01/internal/types"
)

func StartPickupRequestSubscriber(ctx context.Context, storage storage.Storage, client *pubsub.Client, subscriptionID string) error {
	sub := client.Subscription(subscriptionID)

	err := sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		fmt.Printf("Received message: %s\n", string(msg.Data))

		var pr types.PickupRequest
		if err := json.Unmarshal(msg.Data, &pr); err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			log.Println("Nacking message due to unmarshalling failure")
			msg.Nack()
			return
		}

		collector, err := storage.GetCollectorByID(pr.CollectorID)
		if err != nil {
			log.Printf("Error fetching collector: %v", err)
			log.Println("Nacking message due to GetCollectorByID failure")
			msg.Nack()
			return
		}

		heading := fmt.Sprintf("New Request Received")
		message1 := fmt.Sprintf("Dear Collector,\n\nYou have a new pickup request (ID: %d).\nBusiness ID: %d\nWaste Type: %s\nQuantity: %.2f\n\nRegards,\nXphora AI", pr.RequestID, pr.BusinessID, pr.WasteType, pr.Quantity)

		if err := sendNotification(collector, heading, message1); err != nil {
			log.Printf("Error sending email: %v", err)
			log.Println("Nacking message due to sendNotification failure")
			msg.Nack()
			return
		}

		fmt.Printf("Notification sent to collector: %s\n", collector.Email)

		msg.Ack() // Marking message as successfully handled
	})

	if err != nil {
		log.Fatalf("Failed to receive messages: %v", err)
		return err
	}
	return nil
}
