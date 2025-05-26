package pub_sub

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/kartikey1188/build-in-progress_01/internal/types"
)

const DriverLocationTopic = "DRIVER-LOCATION"

func PublishDriverLocation(pubsubClient *pubsub.Client, driverLocation types.DriverLocation) error {
	// Converting driverLocation to JSON
	messageData, err := json.Marshal(driverLocation)
	if err != nil {
		fmt.Printf("Error marshaling driver location: %v", err)
		return err
	}

	// Publishing to the topic
	ctx := context.Background()
	topic := pubsubClient.Topic(DriverLocationTopic)

	result := topic.Publish(ctx, &pubsub.Message{
		Data: messageData,
	})

	// Confirming whether the message was published
	_, err = result.Get(ctx)
	if err != nil {
		fmt.Printf("Error publishing to topic %s: %v", DriverLocationTopic, err)
		return err
	}

	fmt.Printf("Published driver location for driver %d to Pub/Sub topic: %s\n", driverLocation.DriverID, DriverLocationTopic)

	return nil
}
