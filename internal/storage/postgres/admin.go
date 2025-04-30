package postgres

import (
	"fmt"
	"strconv"

	"github.com/kartikey1188/build-in-progress_01/internal/types"
)

// VerifyUser sets the user's is_verified status to true
func (p *Postgres) VerifyUser(userID string) error {
	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	result, err := p.SqlDB.Exec("UPDATE users SET is_verified = true WHERE user_id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to verify user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found", id)
	}

	return nil
}

// UnverifyUser sets the user's is_verified status to false
func (p *Postgres) UnverifyUser(userID string) error {
	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	result, err := p.SqlDB.Exec("UPDATE users SET is_verified = false WHERE user_id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to unverify user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found", id)
	}

	return nil
}

// FlagUser sets is_flagged to true and is_active to false
func (p *Postgres) FlagUser(userID string) error {
	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	result, err := p.SqlDB.Exec(
		"UPDATE users SET is_flagged = true, is_active = false WHERE user_id = $1",
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to flag user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found", id)
	}

	return nil
}

// UnflagUser sets is_flagged to false and is_active to true
func (p *Postgres) UnflagUser(userID string) error {
	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid user ID: %w", err)
	}

	result, err := p.SqlDB.Exec(
		"UPDATE users SET is_flagged = false, is_active = true WHERE user_id = $1",
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to unflag user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found", id)
	}

	return nil
}

// AddServiceCategory inserts a new waste service category and returns its ID.
func (p *Postgres) AddServiceCategory(sc types.ServiceCategory) (int64, error) {
	var id int64
	err := p.SqlDB.QueryRow(
		`INSERT INTO service_categories (waste_type)
		 VALUES ($1)
		 RETURNING category_id`,
		sc.WasteType,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert service category: %w", err)
	}
	return id, nil
}

// AddVehicle inserts a new vehicle and returns its ID.
func (p *Postgres) AddVehicle(v types.Vehicle) (int64, error) {
	var id int64
	err := p.SqlDB.QueryRow(
		`INSERT INTO vehicles (vehicle_type, capacity)
		 VALUES ($1, $2)
		 RETURNING vehicle_id`,
		v.VehicleType, v.Capacity,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert vehicle: %w", err)
	}
	return id, nil
}
