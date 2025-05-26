package pub_sub

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/kartikey1188/build-in-progress_01/internal/storage"
	"github.com/kartikey1188/build-in-progress_01/internal/types"
)

func StartAcceptPickupRequestSubscriber(ctx context.Context, storage storage.Storage, client *pubsub.Client, subscriptionID string) error {
	sub := client.Subscription(subscriptionID)

	err := sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		fmt.Printf("Received message: %s\n", string(msg.Data))

		var pr types.PickupRequest
		if err := json.Unmarshal(msg.Data, &pr); err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			log.Println("Nacking message due to unmarshalling failure")
			msg.Nack()
			return
		}

		collector, err := storage.GetCollectorByID(pr.CollectorID)
		if err != nil {
			log.Printf("Error fetching collector: %v", err)
			log.Println("Nacking message due to GetCollectorByID failure")
			msg.Nack()
			return
		}

		business, err := storage.GetBusinessByID(pr.BusinessID)
		if err != nil {
			log.Printf("Error fetching collector: %v", err)
			log.Println("Nacking message due to GetCollectorByID failure")
			msg.Nack()
			return
		}

		heading := fmt.Sprintf("Pickup Request Accepted")

		message := fmt.Sprintf("Dear Collector,\n\nYou have accepted a pickup request (ID: %d).\nBusiness ID: %d\nWaste Type: %s\nQuantity: %.2f\n\nRegards,\nXphora AI", pr.RequestID, pr.BusinessID, pr.WasteType, pr.Quantity)
		if err := sendNotification(collector, heading, message); err != nil {
			log.Printf("Error sending email: %v", err)
			log.Println("Nacking message due to sendNotification failure")
			msg.Nack()
			return
		}

		fmt.Printf("Notification sent to collector: %s\n", collector.Email)

		message = fmt.Sprintf("Dear Business,\n\nYour pickup request (ID: %d) has been accepted by collector %s.\n\nRegards,\nXphora AI", pr.RequestID, collector.FullName)
		if err := sendNotification(business, heading, message); err != nil {
			log.Printf("Error sending email: %v", err)
			log.Println("Nacking message due to sendNotification failure")
			msg.Nack()
			return
		}

		fmt.Printf("Notification sent to business: %s\n", business.Email)

		msg.Ack() // Marking message as successfully handled
	})

	if err != nil {
		log.Fatalf("Failed to receive messages: %v", err)
		return err
	}
	return nil
}

func StartRejectPickupRequestSubscriber(ctx context.Context, storage storage.Storage, client *pubsub.Client, subscriptionID string) error {
	sub := client.Subscription(subscriptionID)

	err := sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		fmt.Printf("Received message: %s\n", string(msg.Data))

		var pr types.PickupRequest
		if err := json.Unmarshal(msg.Data, &pr); err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			log.Println("Nacking message due to unmarshalling failure")
			msg.Nack()
			return
		}

		collector, err := storage.GetCollectorByID(pr.CollectorID)
		if err != nil {
			log.Printf("Error fetching collector: %v", err)
			log.Println("Nacking message due to GetCollectorByID failure")
			msg.Nack()
			return
		}

		business, err := storage.GetBusinessByID(pr.BusinessID)
		if err != nil {
			log.Printf("Error fetching collector: %v", err)
			log.Println("Nacking message due to GetCollectorByID failure")
			msg.Nack()
			return
		}

		heading := fmt.Sprintf("Pickup Request Rejected")

		message := fmt.Sprintf("Pickup Request Rejected\n\nDear Business,\n\nYour pickup request (ID: %d) has been rejected by the collector.\nBusiness ID: %d\nWaste Type: %s\nQuantity: %.2f\n\nRegards,\nXphora AI", pr.RequestID, pr.BusinessID, pr.WasteType, pr.Quantity)

		if err := sendNotification(collector, heading, message); err != nil {
			log.Printf("Error sending email: %v", err)
			log.Println("Nacking message due to sendNotification failure")
			msg.Nack()
			return
		}
		fmt.Printf("Notification sent to collector: %s\n", collector.Email)

		message = fmt.Sprintf("Pickup Request Rejected\n\nDear Business,\n\nYour pickup request (ID: %d) has been rejected by the collector.\nBusiness ID: %d\nWaste Type: %s\nQuantity: %.2f\n\nRegards,\nXphora AI", pr.RequestID, pr.BusinessID, pr.WasteType, pr.Quantity)

		if err := sendNotification(business, heading, message); err != nil {
			log.Printf("Error sending email: %v", err)
			log.Println("Nacking message due to sendNotification failure")
			msg.Nack()
			return
		}

		fmt.Printf("Notification sent to business: %s\n", business.Email)

		msg.Ack() // Marking message as successfully handled
	})

	if err != nil {
		log.Fatalf("Failed to receive messages: %v", err)
		return err
	}
	return nil
}

func StartDeliverySubscriber(ctx context.Context, storage storage.Storage, client *pubsub.Client, subscriptionID string) error {
	sub := client.Subscription(subscriptionID)

	err := sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		fmt.Printf("Received message: %s\n", string(msg.Data))

		var pr types.PickupRequest
		if err := json.Unmarshal(msg.Data, &pr); err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			log.Println("Nacking message due to unmarshalling failure")
			msg.Nack()
			return
		}

		collector, err := storage.GetCollectorByID(pr.CollectorID)
		if err != nil {
			log.Printf("Error fetching collector: %v", err)
			log.Println("Nacking message due to GetCollectorByID failure")
			msg.Nack()
			return
		}

		business, err := storage.GetBusinessByID(pr.BusinessID)
		if err != nil {
			log.Printf("Error fetching collector: %v", err)
			log.Println("Nacking message due to GetCollectorByID failure")
			msg.Nack()
			return
		}

		heading := fmt.Sprintf("Delivery Started for Pickup Request ID: %d", pr.RequestID)
		message := fmt.Sprintf("Delivery Started for Pickup Request ID: %d, by the Driver %d - with the vehicle %d", pr.RequestID, pr.AssignedDriver, pr.AssignedVehicle)

		if err := sendNotification(collector, heading, message); err != nil {
			log.Printf("Error sending email: %v", err)
			log.Println("Nacking message due to sendNotification failure")
			msg.Nack()
			return
		}

		fmt.Printf("Notification sent to collector: %s\n", collector.Email)

		if err := sendNotification(business, heading, message); err != nil {
			log.Printf("Error sending email: %v", err)
			log.Println("Nacking message due to sendNotification failure")
			msg.Nack()
			return
		}

		fmt.Printf("Notification sent to business: %s\n", business.Email)

		msg.Ack() // Marking message as successfully handled
	})

	if err != nil {
		log.Fatalf("Failed to receive messages: %v", err)
		return err
	}
	return nil
}

func EndDeliverySubscriber(ctx context.Context, storage storage.Storage, client *pubsub.Client, subscriptionID string) error {
	sub := client.Subscription(subscriptionID)

	err := sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		fmt.Printf("Received message: %s\n", string(msg.Data))

		var pr types.PickupRequest
		if err := json.Unmarshal(msg.Data, &pr); err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			log.Println("Nacking message due to unmarshalling failure")
			msg.Nack()
			return
		}

		collector, err := storage.GetCollectorByID(pr.CollectorID)
		if err != nil {
			log.Printf("Error fetching collector: %v", err)
			log.Println("Nacking message due to GetCollectorByID failure")
			msg.Nack()
			return
		}

		business, err := storage.GetBusinessByID(pr.BusinessID)
		if err != nil {
			log.Printf("Error fetching collector: %v", err)
			log.Println("Nacking message due to GetCollectorByID failure")
			msg.Nack()
			return
		}
		heading := fmt.Sprintf("Delivery Completed for Pickup Request ID: %d", pr.RequestID)
		message := fmt.Sprintf("Delivery Completed for Pickup Request ID: %d, by the Driver %d - with the vehicle %d", pr.RequestID, pr.AssignedDriver, pr.AssignedVehicle)

		if err := sendNotification(collector, heading, message); err != nil {
			log.Printf("Error sending email: %v", err)
			log.Println("Nacking message due to sendNotification failure")
			msg.Nack()
			return
		}

		fmt.Printf("Notification sent to collector: %s\n", collector.Email)

		if err := sendNotification(business, heading, message); err != nil {
			log.Printf("Error sending email: %v", err)
			log.Println("Nacking message due to sendNotification failure")
			msg.Nack()
			return
		}
		fmt.Printf("Notification sent to business: %s\n", business.Email)

		msg.Ack() // Marking message as successfully handled
	})

	if err != nil {
		log.Fatalf("Failed to receive messages: %v", err)
		return err
	}
	return nil
}

func StartAssignDriverSubscriber(ctx context.Context, storage storage.Storage, client *pubsub.Client, subscriptionID string) error {
	sub := client.Subscription(subscriptionID)

	err := sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		fmt.Printf("Received message: %s\n", string(msg.Data))

		var pr types.PickupRequest
		if err := json.Unmarshal(msg.Data, &pr); err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			msg.Nack()
			return
		}

		collector, err := storage.GetCollectorByID(pr.CollectorID)
		if err != nil {
			log.Printf("Error fetching collector: %v", err)
			msg.Nack()
			return
		}

		business, err := storage.GetBusinessByID(pr.BusinessID)
		if err != nil {
			log.Printf("Error fetching business: %v", err)
			msg.Nack()
			return
		}

		heading := "Driver Assigned to Pickup Request"

		message := fmt.Sprintf("Dear Collector,\n\nDriver (ID: %d) has been assigned to your pickup request (ID: %d).\n\nRegards,\nXphora AI", pr.AssignedDriver, pr.RequestID)
		if err := sendNotification(collector, heading, message); err != nil {
			log.Printf("Error sending email to collector: %v", err)
			log.Println("Nacking message due to sendNotification failure")
			msg.Nack()
			return
		}
		fmt.Printf("Notification sent to collector: %s\n", collector.Email)

		message = fmt.Sprintf("Dear Business,\n\nDriver (ID: %d) has been assigned to your pickup request (ID: %d).\n\nRegards,\nXphora AI", pr.AssignedDriver, pr.RequestID)
		if err := sendNotification(business, heading, message); err != nil {
			log.Printf("Error sending email to business: %v", err)
			log.Println("Nacking message due to sendNotification failure")
			msg.Nack()
			return
		}
		fmt.Printf("Notification sent to business: %s\n", business.Email)

		msg.Ack() // Marking message as successfully handled
	})

	if err != nil {
		log.Fatalf("Failed to receive messages: %v", err)
		return err
	}
	return nil
}

func StartUnassignDriverSubscriber(ctx context.Context, storage storage.Storage, client *pubsub.Client, subscriptionID string) error {
	sub := client.Subscription(subscriptionID)

	err := sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		fmt.Printf("Received message: %s\n", string(msg.Data))

		var pr types.PickupRequest
		if err := json.Unmarshal(msg.Data, &pr); err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			msg.Nack()
			return
		}

		collector, err := storage.GetCollectorByID(pr.CollectorID)
		if err != nil {
			log.Printf("Error fetching collector: %v", err)
			msg.Nack()
			return
		}

		business, err := storage.GetBusinessByID(pr.BusinessID)
		if err != nil {
			log.Printf("Error fetching business: %v", err)
			msg.Nack()
			return
		}

		heading := "Driver Unassigned from Pickup Request"

		message := fmt.Sprintf("Dear Collector,\n\nDriver has been unassigned from your pickup request (ID: %d).\n\nRegards,\nXphora AI", pr.RequestID)
		if err := sendNotification(collector, heading, message); err != nil {
			log.Printf("Error sending email to collector: %v", err)
			log.Println("Nacking message due to sendNotification failure")
			msg.Nack()
			return
		}
		fmt.Printf("Notification sent to collector: %s\n", collector.Email)

		message = fmt.Sprintf("Dear Business,\n\nDriver has been unassigned from your pickup request (ID: %d).\n\nRegards,\nXphora AI", pr.RequestID)
		if err := sendNotification(business, heading, message); err != nil {
			log.Printf("Error sending email to business: %v", err)
			log.Println("Nacking message due to sendNotification failure")
			msg.Nack()
			return
		}
		fmt.Printf("Notification sent to business: %s\n", business.Email)

		msg.Ack() // Marking message as successfully handled
	})

	if err != nil {
		log.Fatalf("Failed to receive messages: %v", err)
		return err
	}
	return nil
}

// var upgrader = websocket.Upgrader{
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true // Allowing all connections
// 	},
// }

// // Tracking connected websocket clients
// var (
// 	clients   = make(map[*websocket.Conn]bool)
// 	clientsMu sync.Mutex
// )

// // Tracking if we've seen the first location for each driver
// var (
// 	firstLocationSeen   = make(map[int64]bool)
// 	firstLocationSeenMu sync.Mutex
// )

// // WebsocketHandler handles websocket connections
// func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Printf("Error upgrading to websocket: %v", err)
// 		return
// 	}
// 	defer conn.Close()

// 	// Registering new client
// 	clientsMu.Lock()
// 	clients[conn] = true
// 	clientsMu.Unlock()

// 	// Removing client when connection closes
// 	defer func() {
// 		clientsMu.Lock()
// 		delete(clients, conn)
// 		clientsMu.Unlock()
// 	}()

// 	// Keeping the connection alive
// 	for {
// 		_, _, err := conn.ReadMessage()
// 		if err != nil {
// 			break
// 		}
// 	}
// }

// // BroadcastDriverLocation sends driver location to all connected clients
// func BroadcastDriverLocation(driverLocation types.DriverLocation) {
// 	data, err := json.Marshal(driverLocation)
// 	if err != nil {
// 		log.Printf("Error marshaling driver location for broadcast: %v", err)
// 		return
// 	}

// 	clientsMu.Lock()
// 	defer clientsMu.Unlock()

// 	for client := range clients {
// 		err := client.WriteMessage(websocket.TextMessage, data)
// 		if err != nil {
// 			log.Printf("Error writing to websocket: %v", err)
// 			client.Close()
// 			delete(clients, client)
// 		}
// 	}
// }

// // StartDriverLocationSubscriber subscribes to driver location updates
// func StartDriverLocationSubscriber(ctx context.Context, store storage.Storage, client *pubsub.Client, subscriptionID string) error {
// 	sub := client.Subscription(subscriptionID)

// 	// Tracking whether we've processed at least one location for each driver
// 	driverFirstLocations := make(map[int64]bool)
// 	driverLastLocations := make(map[int64]types.DriverLocation)
// 	var locationMutex sync.Mutex

// 	// This will run in the background, receiving messages
// 	err := sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
// 		fmt.Printf("Received driver location update: %s\n", string(msg.Data))

// 		var driverLocation types.DriverLocation
// 		if err := json.Unmarshal(msg.Data, &driverLocation); err != nil {
// 			log.Printf("Error unmarshalling driver location: %v", err)
// 			msg.Nack()
// 			return
// 		}

// 		// Printing to terminal
// 		fmt.Printf("Driver %d location: Lat: %f, Long: %f, Time: %v\n",
// 			driverLocation.DriverID, driverLocation.Latitude, driverLocation.Longitude, driverLocation.Timestamp)

// 		// Checking if this is the first location we've seen for this driver
// 		locationMutex.Lock()
// 		isFirst := !driverFirstLocations[driverLocation.DriverID]
// 		if isFirst {
// 			driverFirstLocations[driverLocation.DriverID] = true

// 			// Saving first location to database with Point="START"
// 			modelLocation := models.DriverLocation{
// 				DriverID:    driverLocation.DriverID,
// 				CollectorID: driverLocation.CollectorID,
// 				VehicleID:   driverLocation.VehicleID,
// 				Latitude:    driverLocation.Latitude,
// 				Longitude:   driverLocation.Longitude,
// 				Timestamp:   driverLocation.Timestamp.Time,
// 				Accuracy:    driverLocation.Accuracy,
// 				Speed:       driverLocation.Speed,
// 				Bearing:     driverLocation.Bearing,
// 				Date:        time.Now(),
// 				Point:       "START",
// 			}

// 			// Saving to database
// 			storeErr := storeDriverLocation(store, modelLocation)
// 			if storeErr != nil {
// 				log.Printf("Error storing START driver location: %v", storeErr)
// 			} else {
// 				fmt.Printf("Stored START location for driver %d\n", driverLocation.DriverID)
// 			}
// 		}

// 		// Always updating the last known location
// 		driverLastLocations[driverLocation.DriverID] = driverLocation
// 		locationMutex.Unlock()

// 		// Broadcasting to websocket clients
// 		BroadcastDriverLocation(driverLocation)

// 		msg.Ack()
// 	})

// 	// Handling termination - storing the last location for each driver
// 	go func() {
// 		<-ctx.Done() // Waiting for context cancellation
// 		locationMutex.Lock()
// 		defer locationMutex.Unlock()

// 		// For each driver, storing their last location as an "END" point
// 		for driverID, lastLocation := range driverLastLocations {
// 			modelLocation := models.DriverLocation{
// 				DriverID:    lastLocation.DriverID,
// 				CollectorID: lastLocation.CollectorID,
// 				VehicleID:   lastLocation.VehicleID,
// 				Latitude:    lastLocation.Latitude,
// 				Longitude:   lastLocation.Longitude,
// 				Timestamp:   lastLocation.Timestamp.Time,
// 				Accuracy:    lastLocation.Accuracy,
// 				Speed:       lastLocation.Speed,
// 				Bearing:     lastLocation.Bearing,
// 				Date:        time.Now(),
// 				Point:       "END",
// 			}

// 			storeErr := storeDriverLocation(store, modelLocation)
// 			if storeErr != nil {
// 				log.Printf("Error storing END driver location for driver %d: %v", driverID, storeErr)
// 			} else {
// 				fmt.Printf("Stored END location for driver %d\n", driverID)
// 			}
// 		}
// 	}()

// 	if err != nil {
// 		log.Printf("Failed to receive driver location messages: %v", err)
// 		return err
// 	}

// 	return nil
// }

// // Helper function to store driver location in the database
// func storeDriverLocation(store storage.Storage, location models.DriverLocation) error {
// 	// Calling the database method to store the location
// 	return store.StoreDriverLocation(location)
// }
