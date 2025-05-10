package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/kartikey1188/build-in-progress_01/internal/storage"
	"github.com/kartikey1188/build-in-progress_01/internal/types"
	"github.com/segmentio/kafka-go"
)

func CreatePickupRequest(storage storage.Storage, pickupRequest types.PickupRequest) (int64, error) {

	id, err2 := storage.CreatePickupRequest(pickupRequest)
	if err2 != nil {
		fmt.Printf("Error creating pickup request in the database: %v", err2)
		return 0, err2
	}

	fmt.Println("Pickup request created successfully in the database")

	// Creating a Kafka topic if it doesn't exist
	brokers := []string{"localhost:29092"}
	topic := "PICKUP-REQUESTS"

	err := createTopic(brokers, topic)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	fmt.Println("Topic created")

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
	defer writer.Close()

	// Converting request to JSON for the message value
	requestJSON, err := json.Marshal(pickupRequest)
	if err != nil {
		log.Printf("Error marshaling pickup request: %v", err)
		return 0, err
	}

	ctx := context.Background()
	err = writer.WriteMessages(ctx,
		kafka.Message{
			Key:   []byte(fmt.Sprintf("request-%d", pickupRequest.BusinessID)),
			Value: requestJSON,
		})

	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	fmt.Println("Successfully published message to topic")
	return id, nil
}
