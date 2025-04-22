package databaseone

import (
	"fmt"
	"time"

	"github.com/kartikey1188/build-in-progress_01/internal/types"
)

// UpdateDriverLocation updates or inserts the driver's current location
func (p *Postgres) UpdateDriverLocation(input types.CollectorDriverLocation) error {
	_, err := p.Db.Exec(`
		INSERT INTO driver_locations 
		(driver_id, latitude, longitude, timestamp, is_active, trip_id, vehicle_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (driver_id) DO UPDATE SET
			latitude = EXCLUDED.latitude,
			longitude = EXCLUDED.longitude,
			timestamp = EXCLUDED.timestamp,
			is_active = EXCLUDED.is_active,
			trip_id = EXCLUDED.trip_id,
			vehicle_id = EXCLUDED.vehicle_id`,
		input.DriverID, input.Latitude, input.Longitude, input.Timestamp.Time,
		input.IsActive, input.TripID, input.VehicleID,
	)
	if err != nil {
		return fmt.Errorf("failed to update driver location: %w", err)
	}
	return nil
}

// GetDriverTrips retrieves all pickup requests assigned to a driver
func (p *Postgres) GetDriverTrips(driverID int64) ([]types.PickupRequest, error) {
	query := `
		SELECT request_id, business_id, collector_id, waste_type, quantity,
			pickup_date, status, assigned_driver, assigned_vehicle, created_at
		FROM pickup_requests
		WHERE assigned_driver = $1
	`
	rows, err := p.Db.Query(query, driverID)
	if err != nil {
		return nil, fmt.Errorf("failed to query driver trips: %w", err)
	}
	defer rows.Close()

	var trips []types.PickupRequest
	for rows.Next() {
		var pr types.PickupRequest
		var pickupDate, createdAt time.Time
		err := rows.Scan(
			&pr.RequestID, &pr.BusinessID, &pr.CollectorID, &pr.WasteType, &pr.Quantity,
			&pickupDate, &pr.Status, &pr.AssignedDriver, &pr.AssignedVehicle, &createdAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan pickup request: %w", err)
		}
		pr.PickupDate = types.DateTime{Time: pickupDate}
		pr.CreatedAt = types.DateTime{Time: createdAt}
		trips = append(trips, pr)
	}
	return trips, nil
}

// UpdateTripStatus updates the status of a pickup request
func (p *Postgres) UpdateTripStatus(tripID string, status string) error {
	_, err := p.Db.Exec(`
		UPDATE pickup_requests 
		SET status = $1 
		WHERE request_id = $2`,
		status, tripID,
	)
	if err != nil {
		return fmt.Errorf("failed to update trip status: %w", err)
	}
	return nil
}
