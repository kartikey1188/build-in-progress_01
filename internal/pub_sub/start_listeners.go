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

	err = StartAcceptPickupRequestSubscriber(ctx, storage, client, cfg.AcceptPickupRequestSubscriptionID)
	if err != nil {
		log.Fatalf("Failed to start AcceptPickupRequestSubscriber: %v", err)
		return err
	}
	err = StartRejectPickupRequestSubscriber(ctx, storage, client, cfg.RejectPickupRequestSubscriptionID)
	if err != nil {
		log.Fatalf("Failed to start RejectPickupRequestSubscriber: %v", err)
		return err
	}

	err = StartDeliverySubscriber(ctx, storage, client, cfg.StartDeliverySubscriptionID)
	if err != nil {
		log.Fatalf("Failed to start DeliverySubscriber: %v", err)
		return err
	}

	err = EndDeliverySubscriber(ctx, storage, client, cfg.EndDeliverySubscriptionID)
	if err != nil {
		log.Fatalf("Failed to start EndDeliverySubscriber: %v", err)
		return err
	}

	err = StartAssignDriverSubscriber(ctx, storage, client, cfg.AssignDriverSubscriptionID)
	if err != nil {
		log.Fatalf("Failed to start AssignDriverSubscriber: %v", err)
		return err
	}

	err = StartUnassignDriverSubscriber(ctx, storage, client, cfg.UnassignDriverSubscriptionID)
	if err != nil {
		log.Fatalf("Failed to start UnassignDriverSubscriber: %v", err)
		return err
	}

	return nil
}
