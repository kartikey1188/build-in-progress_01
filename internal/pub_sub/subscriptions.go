package pub_sub

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
)

type SubscriptionConfig struct {
	SubscriptionID string
	TopicID        string
}

var requiredSubscriptions = []SubscriptionConfig{
	{
		SubscriptionID: "pickup-request-subscription-id",
		TopicID:        "PICKUP-REQUESTS",
	},
}

func InitSubscriptions(client *pubsub.Client) error {
	ctx := context.Background()

	for _, sub := range requiredSubscriptions {
		subscription := client.Subscription(sub.SubscriptionID)

		exists, err := subscription.Exists(ctx)
		if err != nil {
			return fmt.Errorf("error checking subscription %s: %w", sub.SubscriptionID, err)
		}

		if !exists {
			topic := client.Topic(sub.TopicID)

			_, err := client.CreateSubscription(ctx, sub.SubscriptionID, pubsub.SubscriptionConfig{
				Topic: topic,
			})
			if err != nil {
				return fmt.Errorf("failed to create subscription %s: %w", sub.SubscriptionID, err)
			}
		}
	}

	return nil
}
