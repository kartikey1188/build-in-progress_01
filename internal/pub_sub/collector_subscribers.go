package pub_sub

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"
)

func StartPickupRequestSubscriber(ctx context.Context, client *pubsub.Client, subscriptionID string) error {
	sub := client.Subscription(subscriptionID)

	err := sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		fmt.Printf("Received message: %s", string(msg.Data))

		msg.Ack() // Marking message as successfully handled
	})

	if err != nil {
		log.Fatalf("Failed to receive messages: %v", err)
		return err
	}
	return nil
}
