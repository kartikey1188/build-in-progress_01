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

func (p *Postgres) UpdateCollectorProfile(userID int64, update types.CollectorUpdate) (int64, error) {

	_, err1 := p.GetCollectorByID(userID)
	if err1 != nil {
		return 0, fmt.Errorf("collector ID not found")
	}

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
	_, err1 := p.GetCollectorByID(int64(userID))
	if err1 != nil {
		return 0, fmt.Errorf("collector ID not found")
	}

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
	_, err1 := p.GetCollectorByID(int64(userID))
	if err1 != nil {
		return fmt.Errorf("collector ID not found")
	}
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
	_, err1 := p.GetCollectorByID(categoryID)
	if err1 != nil {
		return fmt.Errorf("collector ID not found")
	}
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
	_, err1 := p.GetCollectorByID(int64(userID))
	if err1 != nil {
		return 0, fmt.Errorf("collector ID not found")
	}
	_, err := p.GetVehicle(uint64(input.VehicleID))
	if err != nil {
		return 0, fmt.Errorf("vehicle not found")
	}
	model := models.CollectorVehicle{
		VehicleID:            input.VehicleID,
		CollectorID:          int64(userID),
		VehicleNumber:        input.VehicleNumber,
		MaintenanceDate:      input.MaintenanceDate.Time,
		IsActive:             input.IsActive,
		GPSTrackingID:        input.GPSTrackingID,
		RegistrationDocument: input.RegistrationDocument,
		RegistrationExpiry:   input.RegistrationExpiry.Time,
	}

	if err := p.GormDB.Create(&model).Error; err != nil {
		return 0, fmt.Errorf("create failed: %w", err)
	}

	return model.VehicleID, nil
}

func (p *Postgres) UpdateCollectorVehicle(input types.UpdateCollectorVehicle, userID uint64) error {
	_, err1 := p.GetCollectorByID(int64(userID))
	if err1 != nil {
		return fmt.Errorf("collector ID not found")
	}
	_, err := p.GetVehicle(uint64(input.VehicleID))
	if err != nil {
		return fmt.Errorf("vehicle not found")
	}

	// Using the following method instead of a struct based partial update because of the is_active field (passing "is_active" as false will not update the field, since it is the zero value for bool)

	// Building a map for partial update
	updates := make(map[string]interface{})

	if input.VehicleNumber != "" {
		updates["vehicle_number"] = input.VehicleNumber
	}
	if !input.MaintenanceDate.Time.IsZero() {
		updates["maintenance_date"] = input.MaintenanceDate.Time
	}
	// For booleans, updating if provided (even if false)
	if input.IsActive != nil {
		updates["is_active"] = *input.IsActive
	}
	if input.GPSTrackingID != "" {
		updates["gps_tracking_id"] = input.GPSTrackingID
	}
	if input.RegistrationDocument != "" {
		updates["registration_document"] = input.RegistrationDocument
	}
	if !input.RegistrationExpiry.Time.IsZero() {
		updates["registration_expiry"] = input.RegistrationExpiry.Time
	}

	result := p.GormDB.Model(&models.CollectorVehicle{}).
		Where("vehicle_id = ? AND collector_id = ?", input.VehicleID, userID).
		Updates(updates)

	if result.Error != nil {
		return fmt.Errorf("update failed: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}
	return nil
}

func (p *Postgres) DeleteCollectorVehicle(vehicleID int64, collectorID uint64) error {
	_, err1 := p.GetCollectorByID(int64(collectorID))
	if err1 != nil {
		return fmt.Errorf("collector ID not found")
	}
	_, err := p.GetVehicle(uint64(vehicleID))
	if err != nil {
		return fmt.Errorf("vehicle not found")
	}
	result := p.GormDB.Where("vehicle_id = ? AND collector_id = ?", vehicleID, collectorID).Delete(&models.CollectorVehicle{})

	if result.Error != nil {
		return fmt.Errorf("delete failed: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}

func (p *Postgres) GetCollectorVehicle(collectorID int64, vehicleID int64) (types.CollectorVehicle, error) {
	_, err1 := p.GetCollectorByID(int64(collectorID))
	if err1 != nil {
		return types.CollectorVehicle{}, fmt.Errorf("collector ID not found")
	}
	_, err2 := p.GetVehicle(uint64(vehicleID))
	if err2 != nil {
		return types.CollectorVehicle{}, fmt.Errorf("vehicle not found")
	}

	var vehicle models.CollectorVehicle
	err := p.GormDB.Where("vehicle_id = ? AND collector_id = ?", vehicleID, collectorID).First(&vehicle).Error
	if err != nil {
		return types.CollectorVehicle{}, fmt.Errorf("failed to fetch vehicle: %w", err)
	}
	return types.CollectorVehicle{
		VehicleID:            vehicle.VehicleID,
		CollectorID:          vehicle.CollectorID,
		VehicleNumber:        vehicle.VehicleNumber,
		MaintenanceDate:      types.Date{Time: vehicle.MaintenanceDate},
		IsActive:             vehicle.IsActive,
		GPSTrackingID:        vehicle.GPSTrackingID,
		RegistrationDocument: vehicle.RegistrationDocument,
		RegistrationExpiry:   types.Date{Time: vehicle.RegistrationExpiry},
	}, nil
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
			RegistrationDocument: m.RegistrationDocument,
			RegistrationExpiry:   types.Date{Time: m.RegistrationExpiry},
		}
	}

	return vehicles, nil
}

func (p *Postgres) GetCollectorDriver(collectorID int64, driverID int64) (types.CollectorDriver, error) {
	var driverModel models.CollectorDriver
	err := p.GormDB.Where("driver_id = ? AND collector_id = ?", driverID, collectorID).First(&driverModel).Error
	if err != nil {
		return types.CollectorDriver{}, fmt.Errorf("driver not found: %w", err)
	}

	var userModel models.User
	err = p.GormDB.Where("user_id = ?", driverModel.UserID).First(&userModel).Error
	if err != nil {
		return types.CollectorDriver{}, fmt.Errorf("associated user not found: %w", err)
	}

	return convertCollectorDriverModelToType(driverModel, userModel), nil
}

func (p *Postgres) GetCollectorDrivers(collectorID int64) ([]types.CollectorDriver, error) {
	var driverModels []models.CollectorDriver
	err := p.GormDB.Where("collector_id = ?", collectorID).Find(&driverModels).Error
	if err != nil {
		return nil, fmt.Errorf("error fetching drivers: %w", err)
	}

	drivers := make([]types.CollectorDriver, 0, len(driverModels))
	for _, dm := range driverModels {
		var userModel models.User
		if err := p.GormDB.Where("user_id = ?", dm.UserID).First(&userModel).Error; err != nil {
			return nil, fmt.Errorf("error fetching user for driver %d: %w", dm.UserID, err)
		}
		drivers = append(drivers, convertCollectorDriverModelToType(dm, userModel))
	}
	return drivers, nil
}

// UpdateCollectorDriver updates existing driver details.
func (p *Postgres) UpdateCollectorDriver(input types.UpdateCollectorDriver, collectorID uint64) error {
	_, err := p.GetCollectorByID(int64(collectorID))
	if err != nil {
		return err
	}
	_, err = p.GetCollectorDriver(int64(collectorID), input.DriverID)
	if err != nil {
		return fmt.Errorf("driver ID not found")
	}
	updates := map[string]interface{}{}
	if input.LicenseNumber != "" {
		updates["license_number"] = input.LicenseNumber
	}
	if !input.LicenseExpiry.Time.IsZero() {
		updates["license_expiry"] = input.LicenseExpiry.Time
	}
	if input.IsEmployed != nil {
		updates["is_employed"] = *input.IsEmployed
	}
	if input.IsActive != nil {
		updates["is_active"] = *input.IsActive
	}
	if input.DriverName != "" {
		updates["driver_name"] = input.DriverName
	}
	if input.Rating != 0 {
		updates["rating"] = input.Rating
	}
	if !input.JoiningDate.Time.IsZero() {
		updates["joining_date"] = input.JoiningDate.Time
	}
	result := p.GormDB.Model(&models.CollectorDriver{}).
		Where("driver_id = ? AND collector_id = ?", input.DriverID, collectorID).
		Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows updated")
	}
	return nil
}

// DeleteCollectorDriver removes a driver from a collector.
func (p *Postgres) DeleteCollectorDriver(driverID int64, collectorID uint64) error {
	_, err := p.GetCollectorByID(int64(collectorID))
	if err != nil {
		return err
	}
	_, err = p.GetCollectorDriver(int64(collectorID), driverID)
	if err != nil {
		return fmt.Errorf("driver ID not found")
	}
	result := p.GormDB.Where("driver_id = ? AND collector_id = ?", driverID, collectorID).Delete(&models.CollectorDriver{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows deleted")
	}
	return nil
}

// AssignVehicleToDriver assigns a vehicle to a driver.
func (p *Postgres) AssignVehicleToDriver(driverID int64, vehicleID int64, collectorID uint64) error {
	// Ensuring the collector exists.
	_, err := p.GetCollectorByID(int64(collectorID))
	if err != nil {
		return fmt.Errorf("collector ID not found")
	}
	// Ensuring the driver exists and belongs to the collector.
	_, err = p.GetCollectorDriver(int64(collectorID), driverID)
	if err != nil {
		return fmt.Errorf("driver ID not found")
	}

	// Ensuring the vehicle exists.
	_, err = p.GetVehicle(uint64(vehicleID))
	if err != nil {
		return fmt.Errorf("vehicle ID not found")
	}

	assignment := models.VehicleDriver{
		DriverID:  driverID,
		VehicleID: vehicleID,
	}
	if err := p.GormDB.Create(&assignment).Error; err != nil {
		return err
	}
	return nil
}

func (p *Postgres) UnassignVehicleFromDriver(driverID int64, vehicleID int64, collectorID uint64) error {
	// Ensuring the collector exists.
	_, err := p.GetCollectorByID(int64(collectorID))
	if err != nil {
		return fmt.Errorf("collector ID not found")
	}
	// Ensuring the driver exists and belongs to the collector.
	_, err = p.GetCollectorDriver(int64(collectorID), driverID)
	if err != nil {
		return fmt.Errorf("driver ID not found")
	}

	// Ensuring the vehicle exists.
	_, err = p.GetVehicle(uint64(vehicleID))
	if err != nil {
		return fmt.Errorf("vehicle ID not found")
	}

	result := p.GormDB.Where("driver_id = ? AND vehicle_id = ?", driverID, vehicleID).Delete(&models.VehicleDriver{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows deleted")
	}
	return nil
}

func (p *Postgres) AcceptPickupRequest(requestID int64) error {
	// Fetching the pickup request
	var request models.PickupRequest
	if err := p.GormDB.First(&request, "request_id = ?", requestID).Error; err != nil {
		return fmt.Errorf("pickup request not found: %w", err)
	}

	// Updating the status to accepted
	request.Status = "accepted"
	if err := p.GormDB.Save(&request).Error; err != nil {
		return fmt.Errorf("failed to update pickup request status: %w", err)
	}

	return nil
}

func (p *Postgres) RejectPickupRequest(requestID int64) error {
	// Fetching the pickup request
	var request models.PickupRequest
	if err := p.GormDB.First(&request, "request_id = ?", requestID).Error; err != nil {
		return fmt.Errorf("pickup request not found: %w", err)
	}

	// Updating the status to rejected
	request.Status = "rejected"
	if err := p.GormDB.Save(&request).Error; err != nil {
		return fmt.Errorf("failed to update pickup request status: %w", err)
	}

	return nil
}

func (p *Postgres) AssignTripToDriver(requestID int64, driverID int64) error {
	// Fetching the pickup request
	var request models.PickupRequest
	if err := p.GormDB.First(&request, "request_id = ?", requestID).Error; err != nil {
		return fmt.Errorf("pickup request not found: %w", err)
	}
	// Setting assigned driver
	request.AssignedDriver = driverID
	if err := p.GormDB.Save(&request).Error; err != nil {
		return fmt.Errorf("failed to assign trip to driver: %w", err)
	}
	return nil
}

func (p *Postgres) UnassignTripFromDriver(requestID int64) error {
	// Fetching the pickup request
	var request models.PickupRequest
	if err := p.GormDB.First(&request, "request_id = ?", requestID).Error; err != nil {
		return fmt.Errorf("pickup request not found: %w", err)
	}
	// Unassign the driver by setting to -1
	request.AssignedDriver = -1
	if err := p.GormDB.Save(&request).Error; err != nil {
		return fmt.Errorf("failed to unassign trip from driver: %w", err)
	}
	return nil
}
