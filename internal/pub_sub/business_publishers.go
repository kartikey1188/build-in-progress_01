package pub_sub

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/kartikey1188/build-in-progress_01/internal/storage"
	"github.com/kartikey1188/build-in-progress_01/internal/types"
)

const PickupRequestsTopic = "PICKUP-REQUESTS"

func CreatePickupRequest(storage storage.Storage, pubsubClient *pubsub.Client, pickupRequest types.PickupRequest) (int64, error) {
	// Saving to the database
	id, err := storage.CreatePickupRequest(pickupRequest)
	if err != nil {
		fmt.Printf("Error creating pickup request in the database: %v", err)
		return 0, err
	}

	fmt.Println("Pickup request created successfully in the database")

	pickupRequest.RequestID = id

	// Converting pickupRequest to JSON
	messageData, err := json.Marshal(pickupRequest)
	if err != nil {
		fmt.Printf("Error marshaling pickup request: %v", err)
		return id, err
	}

	// Publishing to the topic
	ctx := context.Background()
	topic := pubsubClient.Topic("PICKUP-REQUESTS")

	result := topic.Publish(ctx, &pubsub.Message{
		Data: messageData,
	})

	// Confirming whether the message was published
	_, err = result.Get(ctx)
	if err != nil {
		fmt.Printf("Error publishing to topic %s: %v", PickupRequestsTopic, err)
		return id, err
	}

	fmt.Println("Published pickup request to Pub/Sub topic:", PickupRequestsTopic)

	return id, nil
}
