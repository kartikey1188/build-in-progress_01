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
