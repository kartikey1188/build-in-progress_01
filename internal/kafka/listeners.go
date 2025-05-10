package kafka

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"sync"

	"github.com/kartikey1188/build-in-progress_01/internal/config"
	"github.com/kartikey1188/build-in-progress_01/internal/storage"
	"github.com/segmentio/kafka-go"
)

func StartKafkaListeners(ctx context.Context, wg *sync.WaitGroup, cfg *config.Config, storage storage.Storage) {
	slog.Info("Starting Kafka listeners...")

	wg.Add(3)

	go PickupRequestListener(ctx, wg, cfg, storage)
	// go startPickupRequestAssignedListener(ctx, wg, cfg, storage)
	// go startTripCompletedListener(ctx, wg, cfg, storage)
}

func PickupRequestListener(ctx context.Context, wg *sync.WaitGroup, cfg *config.Config, storage storage.Storage) {
	defer wg.Done()

	brokers := []string{"localhost:29092"}
	topic := "PICKUP-REQUESTS"
	groupID := "collectors-create-pickup-request"

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		GroupID:  groupID,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	defer reader.Close()

	slog.Info("PickupRequestListener starting to listen...")

	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := reader.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				log.Fatalf("Failed to read message: %s", err)
				continue // Trying again
			}
			fmt.Printf("Received message at offset %d: %s = %s\n", msg.Offset, string(msg.Key), string(msg.Value))
		}
	}
}

// func startPickupRequestAssignedListener(ctx context.Context, wg *sync.WaitGroup, cfg *config.Config, storage storage.Storage) {
// 	defer wg.Done()

// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return
// 			// Add logic to handle messages here
// 		}
// 	}
// }

// func startTripCompletedListener(ctx context.Context, wg *sync.WaitGroup, cfg *config.Config, storage storage.Storage) {
// 	defer wg.Done()

// 	for {
// 		select {
// 		case <-ctx.Done():
// 			return
// 			// Add logic to handle messages here
// 		}
// 	}
// }
