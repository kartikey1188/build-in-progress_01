package postgres

import (
	"errors"
	"fmt"

	"github.com/kartikey1188/build-in-progress_01/internal/models"
	"github.com/kartikey1188/build-in-progress_01/internal/types"
	"gorm.io/gorm"
)

func (p *Postgres) GetCollectorByID(id int64) (types.Collector, error) {
	var collectorModel models.Collector
	var userModel models.User

	// Fetching the collector
	err := p.GormDB.Where("user_id = ?", id).First(&collectorModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return types.Collector{}, fmt.Errorf("collector not found")
		}
		return types.Collector{}, fmt.Errorf("database error: %w", err)
	}

	// Fetching the corresponding user
	err = p.GormDB.First(&userModel, "user_id = ?", collectorModel.UserID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return types.Collector{}, fmt.Errorf("associated user not found")
		}
		return types.Collector{}, fmt.Errorf("database error: %w", err)
	}

	return convertCollectorModelToType(collectorModel, userModel), nil
}

func (p *Postgres) GetCollectors() ([]types.Collector, error) {
	var collectorModels []models.Collector

	// Fetching all collectors
	err := p.GormDB.Find(&collectorModels).Error
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	collectors := make([]types.Collector, 0, len(collectorModels))

	for _, model := range collectorModels {
		var userModel models.User

		// Fetching the user associated with each collector
		err := p.GormDB.First(&userModel, "user_id = ?", model.UserID).Error
		if err != nil {
			return nil, fmt.Errorf("failed to fetch user for collector %d: %w", model.UserID, err)
		}

		collectors = append(collectors, convertCollectorModelToType(model, userModel))
	}

	return collectors, nil
}

func (p *Postgres) UpdateProfile(userID int64, update types.CollectorUpdate) (int64, error) {
	err := p.GormDB.Transaction(func(tx *gorm.DB) error {
		// Updating User fields
		userUpdates := models.User{
			Email:        update.Email,
			PasswordHash: update.PasswordHash,
			FullName:     update.FullName,
			PhoneNumber:  update.PhoneNumber,
			Address:      update.Address,
			IsActive:     update.IsActive,
			ProfileImage: update.ProfileImage,
			LastLogin:    update.LastLogin.Time,
			IsVerified:   update.IsVerified,
			IsFlagged:    update.IsFlagged,
		}

		if err := tx.Model(&models.User{}).Where("user_id = ?", userID).Updates(userUpdates).Error; err != nil {
			return err
		}

		// Updating Collector fields
		collectorUpdates := models.Collector{
			CompanyName:   update.Company_name,
			LicenseNumber: update.License_number,
			Capacity:      update.Capacity,
			LicenseExpiry: update.License_expiry.Time,
		}

		if err := tx.Model(&models.Collector{}).Where("user_id = ?", userID).Updates(collectorUpdates).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return 0, fmt.Errorf("update failed: %w", err)
	}

	return userID, nil
}

func (p *Postgres) AddCollectorServiceCategory(input types.CollectorServiceCategory, userID uint64) (int64, error) {
	_, err := p.GetServiceCategory(uint64(input.CategoryID))
	if err != nil {
		return 0, fmt.Errorf("service category not found")
	}
	model := models.CollectorServiceCategory{
		CategoryID:           input.CategoryID,
		CollectorID:          int64(userID),
		PricePerKg:           input.PricePerKg,
		MaximumCapacity:      input.MaximumCapacity,
		HandlingRequirements: input.HandlingRequirements,
	}

	if err := p.GormDB.Create(&model).Error; err != nil {
		return 0, fmt.Errorf("create failed: %w", err)
	}

	return model.CategoryID, nil
}

func (p *Postgres) UpdateCollectorServiceCategory(input types.UpdateCollectorServiceCategory, userID uint64) error {
	_, err := p.GetServiceCategory(uint64(input.CategoryID))
	if err != nil {
		return fmt.Errorf("service category not found")
	}
	result := p.GormDB.Model(&models.CollectorServiceCategory{}).
		Where("category_id = ? AND collector_id = ?", input.CategoryID, userID).
		Updates(models.CollectorServiceCategory{
			PricePerKg:           input.PricePerKg,
			MaximumCapacity:      input.MaximumCapacity,
			HandlingRequirements: input.HandlingRequirements,
		})

	if result.Error != nil {
		return fmt.Errorf("update failed: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}

func (p *Postgres) DeleteCollectorServiceCategory(categoryID int64, collectorID uint64) error {
	_, err := p.GetServiceCategory(uint64(categoryID))
	if err != nil {
		return fmt.Errorf("service category not found")
	}
	result := p.GormDB.Where("category_id = ? AND collector_id = ?", categoryID, collectorID).Delete(&models.CollectorServiceCategory{})

	if result.Error != nil {
		return fmt.Errorf("delete failed: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}

func (p *Postgres) AddCollectorVehicle(input types.CollectorVehicle, userID uint64) (int64, error) {
	model := models.CollectorVehicle{
		VehicleID:            input.VehicleID,
		CollectorID:          int64(userID),
		VehicleNumber:        input.VehicleNumber,
		MaintenanceDate:      input.MaintenanceDate.Time,
		IsActive:             input.IsActive,
		GPSTrackingID:        input.GPSTrackingID,
		AssignedDriverID:     input.AssignedDriverID,
		RegistrationDocument: input.RegistrationDocument,
		RegistrationExpiry:   input.RegistrationExpiry.Time,
	}

	if err := p.GormDB.Create(&model).Error; err != nil {
		return 0, fmt.Errorf("create failed: %w", err)
	}

	return model.VehicleID, nil
}

func (p *Postgres) UpdateCollectorVehicle(id int64, collectorID int64, input types.CollectorVehicle) error {
	updates := models.CollectorVehicle{
		VehicleID:            input.VehicleID,
		VehicleNumber:        input.VehicleNumber,
		MaintenanceDate:      input.MaintenanceDate.Time,
		IsActive:             input.IsActive,
		GPSTrackingID:        input.GPSTrackingID,
		AssignedDriverID:     input.AssignedDriverID,
		RegistrationDocument: input.RegistrationDocument,
		RegistrationExpiry:   input.RegistrationExpiry.Time,
	}

	result := p.GormDB.Model(&models.CollectorVehicle{}).
		Where("collector_vehicle_id = ? AND collector_id = ?", id, collectorID).
		Updates(updates)

	if result.Error != nil {
		return fmt.Errorf("update failed: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no vehicle found with id %d for collector %d", id, collectorID)
	}
	return nil
}

func (p *Postgres) ActivateCollectorVehicle(id int64, collectorID int64) error {
	result := p.GormDB.Model(&models.CollectorVehicle{}).
		Where("collector_vehicle_id = ? AND collector_id = ?", id, collectorID).
		Update("is_active", true)

	if result.Error != nil {
		return fmt.Errorf("activation failed: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no vehicle found with id %d for collector %d", id, collectorID)
	}
	return nil
}

func (p *Postgres) DeactivateCollectorVehicle(id int64, collectorID int64) error {
	result := p.GormDB.Model(&models.CollectorVehicle{}).
		Where("collector_vehicle_id = ? AND collector_id = ?", id, collectorID).
		Update("is_active", false)

	if result.Error != nil {
		return fmt.Errorf("deactivation failed: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no vehicle found with id %d for collector %d", id, collectorID)
	}
	return nil
}

func (p *Postgres) AddCollectorDriver(input types.CollectorDriver) (int64, error) {
	model := models.CollectorDriver{
		CollectorID:       input.CollectorID,
		LicenseNumber:     input.LicenseNumber,
		LicenseExpiry:     input.LicenseExpiry.Time,
		IsEmployed:        input.IsEmployed,
		IsActive:          input.IsActive,
		Rating:            input.Rating,
		JoiningDate:       input.JoiningDate.Time,
		AssignedVehicleID: input.AssignedVehicleID,
	}

	if err := p.GormDB.Create(&model).Error; err != nil {
		return 0, fmt.Errorf("create failed: %w", err)
	}

	return model.DriverID, nil
}

func (p *Postgres) UpdateCollectorDriver(id int64, collectorID int64, input types.CollectorDriver) error {
	updates := models.CollectorDriver{
		LicenseNumber:     input.LicenseNumber,
		LicenseExpiry:     input.LicenseExpiry.Time,
		IsEmployed:        input.IsEmployed,
		IsActive:          input.IsActive,
		Rating:            input.Rating,
		JoiningDate:       input.JoiningDate.Time,
		AssignedVehicleID: input.AssignedVehicleID,
	}

	result := p.GormDB.Model(&models.CollectorDriver{}).
		Where("driver_id = ? AND collector_id = ?", id, collectorID).
		Updates(updates)

	if result.Error != nil {
		return fmt.Errorf("update failed: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no driver found with id %d for collector %d", id, collectorID)
	}
	return nil
}

func (p *Postgres) AssignVehicleToDriver(driverID int64, vehicleID int64, collectorID int64) error {
	// Verifying whether vehicle belongs to collector
	var vehicle models.CollectorVehicle
	if err := p.GormDB.Where("collector_vehicle_id = ? AND collector_id = ?", vehicleID, collectorID).
		First(&vehicle).Error; err != nil {
		return fmt.Errorf("vehicle not found or access denied: %w", err)
	}

	result := p.GormDB.Model(&models.CollectorDriver{}).
		Where("driver_id = ? AND collector_id = ?", driverID, collectorID).
		Update("assigned_vehicle_id", vehicleID)

	if result.Error != nil {
		return fmt.Errorf("assignment failed: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no driver found with id %d for collector %d", driverID, collectorID)
	}
	return nil
}

func (p *Postgres) GetCollectorServiceCategories(collectorID int64) ([]types.CollectorServiceCategory, error) {
	var models []models.CollectorServiceCategory

	if err := p.GormDB.Where("collector_id = ?", collectorID).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	categories := make([]types.CollectorServiceCategory, len(models))
	for i, m := range models {
		categories[i] = types.CollectorServiceCategory{
			CategoryID:           m.CategoryID,
			CollectorID:          m.CollectorID,
			PricePerKg:           m.PricePerKg,
			MaximumCapacity:      m.MaximumCapacity,
			HandlingRequirements: m.HandlingRequirements,
		}
	}

	return categories, nil
}

func (p *Postgres) GetCollectorVehicles(collectorID int64) ([]types.CollectorVehicle, error) {
	var models []models.CollectorVehicle

	if err := p.GormDB.Where("collector_id = ?", collectorID).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	vehicles := make([]types.CollectorVehicle, len(models))
	for i, m := range models {
		vehicles[i] = types.CollectorVehicle{
			VehicleID:            m.VehicleID,
			CollectorID:          m.CollectorID,
			VehicleNumber:        m.VehicleNumber,
			MaintenanceDate:      types.Date{Time: m.MaintenanceDate},
			IsActive:             m.IsActive,
			GPSTrackingID:        m.GPSTrackingID,
			AssignedDriverID:     m.AssignedDriverID,
			RegistrationDocument: m.RegistrationDocument,
			RegistrationExpiry:   types.Date{Time: m.RegistrationExpiry},
		}
	}

	return vehicles, nil
}

func (p *Postgres) GetCollectorDrivers(collectorID int64) ([]types.CollectorDriver, error) {
	var models []models.CollectorDriver

	if err := p.GormDB.Where("collector_id = ?", collectorID).Find(&models).Error; err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	drivers := make([]types.CollectorDriver, len(models))
	for i, m := range models {
		drivers[i] = types.CollectorDriver{
			DriverID:          m.DriverID,
			CollectorID:       m.CollectorID,
			LicenseNumber:     m.LicenseNumber,
			LicenseExpiry:     types.Date{Time: m.LicenseExpiry},
			AssignedVehicleID: m.AssignedVehicleID,
			IsActive:          m.IsActive,
			Rating:            m.Rating,
			JoiningDate:       types.Date{Time: m.JoiningDate},
			IsEmployed:        m.IsEmployed,
		}
	}

	return drivers, nil
}
