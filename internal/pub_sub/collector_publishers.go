package pub_sub

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/pubsub"
	"github.com/kartikey1188/build-in-progress_01/internal/storage"
	"github.com/kartikey1188/build-in-progress_01/internal/types"
)

func AcceptPickupRequest(storage storage.Storage, pubsubClient *pubsub.Client, pickupRequestID int64) error {
	// Saving to the database
	err := storage.AcceptPickupRequest(pickupRequestID)
	if err != nil {
		fmt.Printf("Error accepting pickup request in the database: %v", err)
		return err
	}

	fmt.Println("Pickup request accepted successfully in the database")

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
	topic := pubsubClient.Topic("PICKUP-REQUESTS")

	result := topic.Publish(ctx, &pubsub.Message{
		Data: messageData,
	})

	// Confirming whether the message was published
	_, err = result.Get(ctx)
	if err != nil {
		fmt.Printf("Error publishing to topic %s: %v", PickupRequestsTopic, err)
		return err
	}

	fmt.Println("Published pickup request to Pub/Sub topic:", PickupRequestsTopic)

	return nil
}

func RejectPickupRequest(storage storage.Storage, pubsubClient *pubsub.Client, pickupRequestID int64) error {
	// Saving to the database
	err := storage.RejectPickupRequest(pickupRequestID)
	if err != nil {
		fmt.Printf("Error rejecting pickup request in the database: %v", err)
		return err
	}

	fmt.Println("Pickup request rejected successfully in the database")

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
	topic := pubsubClient.Topic("PICKUP-REQUESTS")

	result := topic.Publish(ctx, &pubsub.Message{
		Data: messageData,
	})

	// Confirming whether the message was published
	_, err = result.Get(ctx)
	if err != nil {
		fmt.Printf("Error publishing to topic %s: %v", PickupRequestsTopic, err)
		return err
	}

	fmt.Println("Published pickup request to Pub/Sub topic:", PickupRequestsTopic)

	return nil
}

const AssignmentsTopic = "ASSIGNMENTS"

func AssignTripToDriver(storage storage.Storage, pubsubClient *pubsub.Client, pickupRequestID int64, driver types.CollectorDriver) error {
	// Saving to the database
	err := storage.AssignTripToDriver(pickupRequestID, driver.UserID)
	if err != nil {
		fmt.Printf("Error assigning trip to driver in the database: %v", err)
		return err
	}

	fmt.Println("Trip assigned to driver successfully in the database")

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
	topic := pubsubClient.Topic("ASSIGNMENTS")

	result := topic.Publish(ctx, &pubsub.Message{
		Data: messageData,
	})

	// Confirming whether the message was published
	_, err = result.Get(ctx)
	if err != nil {
		fmt.Printf("Error publishing to topic %s: %v", AssignmentsTopic, err)
		return err
	}

	fmt.Println("Published assignment to Pub/Sub topic:", AssignmentsTopic)

	return nil
}

func UnassignTripFromDriver(storage storage.Storage, pubsubClient *pubsub.Client, pickupRequestID int64) error {
	// Saving to the database
	err := storage.UnassignTripFromDriver(pickupRequestID)
	if err != nil {
		fmt.Printf("Error unassigning trip from driver in the database: %v", err)
		return err
	}

	fmt.Println("Trip unassigned from driver successfully in the database")

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
	topic := pubsubClient.Topic("ASSIGNMENTS")

	result := topic.Publish(ctx, &pubsub.Message{
		Data: messageData,
	})

	// Confirming whether the message was published
	_, err = result.Get(ctx)
	if err != nil {
		fmt.Printf("Error publishing to topic %s: %v", AssignmentsTopic, err)
		return err
	}

	fmt.Println("Published unassignment to Pub/Sub topic:", AssignmentsTopic)

	return nil
}
