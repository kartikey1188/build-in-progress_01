package pub_sub

import (
	"context"
	"fmt"

	"cloud.google.com/go/pubsub"
)

var requiredTopics = []string{
	"PICKUP-REQUESTS",
}

func InitTopics(client *pubsub.Client) error {
	ctx := context.Background()

	for _, topicName := range requiredTopics {
		topic := client.Topic(topicName)
		exists, err := topic.Exists(ctx)
		if err != nil {
			return fmt.Errorf("error checking topic %s: %w", topicName, err)
		}

		if !exists {
			_, err = client.CreateTopic(ctx, topicName)
			if err != nil {
				return fmt.Errorf("failed to create topic %s: %w", topicName, err)
			}
		}
	}

	return nil
}
