package postgres

import (
	"errors"
	"fmt"

	"github.com/kartikey1188/build-in-progress_01/internal/models"
	"github.com/kartikey1188/build-in-progress_01/internal/types"
	"gorm.io/gorm"
)

// GetBusinessByID retrieves a business by its ID.
func (p *Postgres) GetBusinessByID(id int64) (types.Business, error) {
	var businessModel models.Business
	var userModel models.User

	err := p.GormDB.Where("user_id = ?", id).First(&businessModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return types.Business{}, fmt.Errorf("business not found")
		}
		return types.Business{}, fmt.Errorf("database error: %w", err)
	}

	err = p.GormDB.First(&userModel, "user_id = ?", businessModel.UserID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return types.Business{}, fmt.Errorf("associated user not found")
		}
		return types.Business{}, fmt.Errorf("database error: %w", err)
	}

	return convertBusinessModelToType(businessModel, userModel), nil
}

// UpdateBusinessProfile updates a business's profile.
func (p *Postgres) UpdateBusinessProfile(userID int64, update types.BusinessUpdate) (int64, error) {
	// Verify business exists
	_, err := p.GetBusinessByID(userID)
	if err != nil {
		return 0, fmt.Errorf("business ID not found")
	}

	err = p.GormDB.Transaction(func(tx *gorm.DB) error {
		// Update User fields
		userUpdates := models.User{
			Email:        update.Email,
			PasswordHash: update.PasswordHash,
			FullName:     update.FullName,
			PhoneNumber:  update.PhoneNumber,
			Address:      update.Address,
			ProfileImage: update.ProfileImage,
		}
		if err := tx.Model(&models.User{}).Where("user_id = ?", userID).Updates(userUpdates).Error; err != nil {
			return err
		}

		// Update Business fields
		businessUpdates := models.Business{
			BusinessName:       update.Business_name,
			BusinessType:       update.Business_type,
			RegistrationNumber: update.Registration_number,
			GstID:              update.Gst_id,
			BusinessAddress:    update.Business_address,
		}
		if err := tx.Model(&models.Business{}).Where("user_id = ?", userID).Updates(businessUpdates).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("update failed: %w", err)
	}

	return userID, nil
}

func (p *Postgres) CreatePickupRequest(request types.PickupRequest) (int64, error) {
	_, err := p.GetBusinessByID(request.BusinessID)
	if err != nil {
		return 0, fmt.Errorf("business ID not found: %w", err)
	}

	_, err = p.GetCollectorByID(request.CollectorID)
	if err != nil {
		return 0, fmt.Errorf("collector ID not found: %w", err)
	}

	model := models.PickupRequest{
		BusinessID:           request.BusinessID,
		CollectorID:          request.CollectorID,
		WasteType:            request.WasteType,
		Quantity:             request.Quantity,
		PickupDate:           request.PickupDate.Time,
		Status:               request.Status,
		HandlingRequirements: request.HandlingRequirements,
		AssignedDriver:       request.AssignedDriver,
		AssignedVehicle:      request.AssignedVehicle,
		CreatedAt:            request.CreatedAt.Time,
	}

	if err := p.GormDB.Create(&model).Error; err != nil {
		return 0, fmt.Errorf("failed to create pickup request: %w", err)
	}

	return model.RequestID, nil
}

func (p *Postgres) GetPickupRequestByID(id int64) (types.PickupRequest, error) {
	var model models.PickupRequest
	err := p.GormDB.First(&model, "request_id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return types.PickupRequest{}, fmt.Errorf("pickup request not found")
		}
		return types.PickupRequest{}, fmt.Errorf("database error: %w", err)
	}

	return convertPickupRequestModelToType(model), nil
}

func (p *Postgres) GetAllPickupRequestsForBusiness(businessID int64) ([]types.PickupRequest, error) {
	var models []models.PickupRequest
	err := p.GormDB.Where("business_id = ?", businessID).Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	var requests []types.PickupRequest
	for _, model := range models {
		requests = append(requests, convertPickupRequestModelToType(model))
	}

	return requests, nil
}

func (p *Postgres) UpdatePickupRequest(requestID int64, input types.UpdatePickupRequest) error {
	// Check if the pickup request exists
	var existing models.PickupRequest
	err := p.GormDB.First(&existing, "request_id = ?", requestID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("pickup request not found")
		}
		return fmt.Errorf("database error: %w", err)
	}

	updates := models.PickupRequest{
		WasteType:            input.WasteType,
		Quantity:             input.Quantity,
		PickupDate:           input.PickupDate.Time,
		Status:               input.Status,
		HandlingRequirements: input.HandlingRequirements,
		AssignedDriver:       input.AssignedDriver,
		AssignedVehicle:      input.AssignedVehicle,
		CreatedAt:            input.CreatedAt.Time,
	}

	result := p.GormDB.Model(&models.PickupRequest{}).
		Where("request_id = ?", requestID).
		Updates(updates)

	if result.Error != nil {
		return fmt.Errorf("update failed: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}
	return nil
}
