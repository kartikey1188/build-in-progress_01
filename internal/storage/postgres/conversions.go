package postgres

import (
	"github.com/kartikey1188/build-in-progress_01/internal/models"
	"github.com/kartikey1188/build-in-progress_01/internal/types"
)

// convertCollectorModelToType transforms a models.Collector and its associated User into types.Collector
func convertCollectorModelToType(collector models.Collector, user models.User) types.Collector {
	return types.Collector{
		User: types.User{
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
		},
		Company_name:   collector.CompanyName,
		License_number: collector.LicenseNumber,
		Capacity:       collector.Capacity,
		License_expiry: types.Date{Time: collector.LicenseExpiry},
	}
}

// convertCollectorTypeToModel transforms types.Collector into models.Collector and models.User
func convertCollectorTypeToModel(collector types.Collector) (models.Collector, models.User) {
	return models.Collector{
			UserID:        collector.User.UserID,
			CompanyName:   collector.Company_name,
			LicenseNumber: collector.License_number,
			Capacity:      collector.Capacity,
			LicenseExpiry: collector.License_expiry.Time,
		}, models.User{
			UserID:       collector.User.UserID,
			Email:        collector.User.Email,
			PasswordHash: collector.User.PasswordHash,
			FullName:     collector.User.FullName,
			PhoneNumber:  collector.User.PhoneNumber,
			Address:      collector.User.Address,
			Registration: collector.User.Registration.Time,
			Role:         collector.User.Role,
			IsActive:     collector.User.IsActive,
			ProfileImage: collector.User.ProfileImage,
			LastLogin:    collector.User.LastLogin.Time,
			IsVerified:   collector.User.IsVerified,
			IsFlagged:    collector.User.IsFlagged,
		}
}
