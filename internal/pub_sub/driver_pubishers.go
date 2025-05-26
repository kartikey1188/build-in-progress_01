package pub_sub

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/kartikey1188/build-in-progress_01/internal/storage"
)

const DeliveryTopic = "DELIVERY"

func StartDelivery(storage storage.Storage, pubsubClient *pubsub.Client, pickupRequestID int64) error {

	pickupRequest, err := storage.GetPickupRequestByID(pickupRequestID)
	if err != nil {
		fmt.Printf("Error retrieving pickup request ID: %v", err)
		return err
	}

	// Converting pickupRequest to JSON
	messageData, err := json.Marshal(pickupRequest)
	if err != nil {
		fmt.Printf("Error marshaling pickup request: %v", err)
		return err
	}

	// Publishing to the topic
	ctx := context.Background()
	topic := pubsubClient.Topic("DELIVERY")

	result := topic.Publish(ctx, &pubsub.Message{
		Data: messageData,
	})

	// Confirming whether the message was published
	_, err = result.Get(ctx)
	if err != nil {
		fmt.Printf("Error publishing to topic %s: %v", DeliveryTopic, err)
		return err
	}

	fmt.Println("Published pickup request to Pub/Sub topic:", DeliveryTopic)

	return nil
}

func EndDelivery(storage storage.Storage, pubsubClient *pubsub.Client, pickupRequestID int64) error {
	err := storage.EndDelivery(pickupRequestID)
	if err != nil {
		fmt.Printf("Error completing delivery in the database: %v", err)
		return err
	}

	fmt.Println("Delivery completed successfully in the database")

	pickupRequest, err := storage.GetPickupRequestByID(pickupRequestID)
	if err != nil {
		fmt.Printf("Error retrieving pickup request ID: %v", err)
		return err
	}

	// Converting pickupRequest to JSON
	messageData, err := json.Marshal(pickupRequest)
	if err != nil {
		fmt.Printf("Error marshaling pickup request: %v", err)
		return err
	}

	// Publishing to the topic
	ctx := context.Background()
	topic := pubsubClient.Topic("DELIVERY")

	result := topic.Publish(ctx, &pubsub.Message{
		Data: messageData,
	})

	// Confirming whether the message was published
	_, err = result.Get(ctx)
	if err != nil {
		fmt.Printf("Error publishing to topic %s: %v", DeliveryTopic, err)
		return err
	}

	fmt.Println("Published delivery completion to Pub/Sub topic:", DeliveryTopic)

	return nil
}
