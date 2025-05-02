package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/kartikey1188/build-in-progress_01/internal/models"
	"github.com/kartikey1188/build-in-progress_01/internal/types"
)

func (p *Postgres) GetAllServiceCategories() ([]types.ServiceCategory, error) {
	var categories []models.ServiceCategory
	err := p.GormDB.Find(&categories).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch service categories: %w", err)
	}

	result := make([]types.ServiceCategory, len(categories))
	for i, category := range categories {
		result[i] = types.ServiceCategory{
			CategoryID: category.CategoryID,
			WasteType:  category.WasteType,
		}
	}
	return result, nil
}

func (p *Postgres) GetServiceCategory(categoryID uint64) (types.ServiceCategory, error) {
	var category models.ServiceCategory
	err := p.GormDB.First(&category, "category_id = ?", categoryID).Error
	if err != nil {
		return types.ServiceCategory{}, fmt.Errorf("failed to fetch service category: %w", err)
	}

	return types.ServiceCategory{
		CategoryID: category.CategoryID,
		WasteType:  category.WasteType,
	}, nil
}

func (p *Postgres) GetAllVehicles() ([]types.Vehicle, error) {
	var vehicles []models.Vehicle
	err := p.GormDB.Find(&vehicles).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch vehicles: %w", err)
	}

	result := make([]types.Vehicle, len(vehicles))
	for i, vehicle := range vehicles {
		result[i] = types.Vehicle{
			VehicleID:   vehicle.VehicleID,
			VehicleType: vehicle.VehicleType,
			Capacity:    vehicle.Capacity,
		}
	}
	return result, nil
}

func (p *Postgres) GetVehicle(vehicleID uint64) (types.Vehicle, error) {
	var vehicle models.Vehicle
	err := p.GormDB.First(&vehicle, "vehicle_id = ?", vehicleID).Error
	if err != nil {
		return types.Vehicle{}, fmt.Errorf("failed to fetch vehicle: %w", err)
	}

	return types.Vehicle{
		VehicleID:   vehicle.VehicleID,
		VehicleType: vehicle.VehicleType,
		Capacity:    vehicle.Capacity,
	}, nil
}

func (p *Postgres) GetUserByID(userID uint64) (types.User, error) {
	var user models.User
	err := p.GormDB.First(&user, "user_id = ?", userID).Error
	if err != nil {
		return types.User{}, fmt.Errorf("failed to fetch user: %w", err)
	}

	return types.User{
		UserID:       user.UserID,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		FullName:     user.FullName,
		PhoneNumber:  user.PhoneNumber,
		Address:      user.Address,
		Registration: types.Date{Time: user.Registration},
		Role:         user.Role,
		IsActive:     user.IsActive,
		ProfileImage: user.ProfileImage,
		LastLogin:    types.DateTime{Time: user.LastLogin},
		IsVerified:   user.IsVerified,
		IsFlagged:    user.IsFlagged,
	}, nil
}

func (p *Postgres) GetUserByEmail(email string) (types.User, error) {
	var user types.User
	var registration, lastLogin time.Time

	err := p.SqlDB.QueryRow(`
		SELECT user_id, email, password_hash, full_name, phone_number,
			address, registration_date, role, is_active, profile_image,
			last_login, is_verified, is_flagged
		FROM users
		WHERE email = $1
		LIMIT 1`, email,
	).Scan(
		&user.UserID, &user.Email, &user.PasswordHash, &user.FullName, &user.PhoneNumber,
		&user.Address, &registration, &user.Role, &user.IsActive, &user.ProfileImage,
		&lastLogin, &user.IsVerified, &user.IsFlagged,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.User{}, fmt.Errorf("user not found")
		}
		return types.User{}, fmt.Errorf("query error: %w", err)
	}

	user.Registration = types.Date{Time: registration}
	user.LastLogin = types.DateTime{Time: lastLogin}
	return user, nil
}
