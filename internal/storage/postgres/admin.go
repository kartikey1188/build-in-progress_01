package postgres

import (
	"fmt"
	"strconv"

	"github.com/kartikey1188/build-in-progress_01/internal/models"
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

func (p *Postgres) DeleteServiceCategory(categoryID uint64) error {

	result, err := p.SqlDB.Exec("DELETE FROM service_categories WHERE category_id = $1", categoryID)
	if err != nil {
		return fmt.Errorf("failed to delete service category: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("service category with ID %d not found", categoryID)
	}

	return nil
}

func (p *Postgres) DeleteVehicle(vehicleID uint64) error {

	result, err := p.SqlDB.Exec("DELETE FROM vehicles WHERE vehicle_id = $1", vehicleID)
	if err != nil {
		return fmt.Errorf("failed to delete vehicle: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("vehicle with ID %d not found", vehicleID)
	}

	return nil
}

func (p *Postgres) GetAllCollectors() ([]types.Collector, error) {
	type collectorJoin struct {
		models.Collector
		models.User
	}
	var joins []collectorJoin
	if err := p.GormDB.Table("collectors").
		Joins("INNER JOIN users ON users.user_id = collectors.user_id").
		Scan(&joins).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch collectors: %w", err)
	}
	collectors := make([]types.Collector, 0, len(joins))
	for _, join := range joins {
		collectors = append(collectors, convertCollectorModelToType(join.Collector, join.User))
	}
	return collectors, nil
}

func (p *Postgres) GetAllBusinesses() ([]types.Business, error) {
	type businessJoin struct {
		models.Business
		models.User
	}
	var joins []businessJoin
	if err := p.GormDB.Table("businesses").
		Joins("INNER JOIN users ON users.user_id = businesses.user_id").
		Scan(&joins).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch businesses: %w", err)
	}
	businesses := make([]types.Business, 0, len(joins))
	for _, join := range joins {
		businesses = append(businesses, convertBusinessModelToType(join.Business, join.User))
	}
	return businesses, nil
}

func (p *Postgres) GetAllUsers() ([]types.User, error) {
	var userModels []models.User
	if err := p.GormDB.Find(&userModels).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}

	users := make([]types.User, 0, len(userModels))
	for _, u := range userModels {
		users = append(users, types.User{
			UserID:       u.UserID,
			Email:        u.Email,
			PasswordHash: u.PasswordHash,
			FullName:     u.FullName,
			PhoneNumber:  u.PhoneNumber,
			Address:      u.Address,
			Registration: types.Date{Time: u.Registration},
			Role:         u.Role,
			IsActive:     u.IsActive,
			ProfileImage: u.ProfileImage,
			LastLogin:    types.DateTime{Time: u.LastLogin},
			IsVerified:   u.IsVerified,
			IsFlagged:    u.IsFlagged,
		})
	}
	return users, nil
}
