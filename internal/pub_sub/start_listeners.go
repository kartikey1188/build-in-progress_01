package pub_sub

import (
	"context"
	"log"
	"log/slog"

	"cloud.google.com/go/pubsub"
	"github.com/kartikey1188/build-in-progress_01/internal/config"
	"github.com/kartikey1188/build-in-progress_01/internal/storage"
)

func StartListeners(ctx context.Context, client *pubsub.Client, cfg *config.Config, storage storage.Storage) error {
	slog.Info("Starting listeners...")

	err := StartPickupRequestSubscriber(ctx, storage, client, cfg.PickupRequestSubscriptionID)
	if err != nil {
		log.Fatalf("Failed to start PickupRequestSubscriber: %v", err)
		return err
	}

	return nil
}
