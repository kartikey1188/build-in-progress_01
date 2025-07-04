package postgres

import (
	"encoding/json"

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

// convertBusinessModelToType converts a Business model and its User model to types.Business
func convertBusinessModelToType(b models.Business, u models.User) types.Business {
	return types.Business{
		User: types.User{
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
		},
		Business_name:       b.BusinessName,
		Business_type:       b.BusinessType,
		Registration_number: b.RegistrationNumber,
		Gst_id:              b.GstID,
		Business_address:    b.BusinessAddress,
	}
}

func convertPickupRequestModelToType(model models.PickupRequest) types.PickupRequest {
	return types.PickupRequest{
		RequestID:            model.RequestID,
		BusinessID:           model.BusinessID,
		CollectorID:          model.CollectorID,
		WasteType:            model.WasteType,
		Quantity:             model.Quantity,
		PickupDate:           types.DateTime{Time: model.PickupDate},
		Status:               model.Status,
		HandlingRequirements: model.HandlingRequirements,
		AssignedDriver:       model.AssignedDriver,
		AssignedVehicle:      model.AssignedVehicle,
		CreatedAt:            types.DateTime{Time: model.CreatedAt},
	}
}

func convertCollectorDriverModelToType(driver models.CollectorDriver, user models.User) types.CollectorDriver {
	return types.CollectorDriver{
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
		CollectorID:   driver.CollectorID,
		LicenseNumber: driver.LicenseNumber,
		DriverName:    driver.DriverName,
		LicenseExpiry: types.Date{Time: driver.LicenseExpiry},
		IsEmployed:    driver.IsEmployed,
		IsActive:      driver.IsActive,
		Rating:        driver.Rating,
		JoiningDate:   types.Date{Time: driver.JoiningDate},
	}
}

func convertFacilityModelToType(facility models.Facility) types.Facility {
	var wasteTypes []string
	if facility.WasteTypes != "" {
		// Parse JSON string to slice
		json.Unmarshal([]byte(facility.WasteTypes), &wasteTypes)
	}

	return types.Facility{
		FacilityID:      facility.FacilityID,
		Name:            facility.Name,
		Location:        facility.Location,
		Status:          facility.Status,
		Capacity:        facility.Capacity,
		ComplianceScore: facility.ComplianceScore,
		DailyVolume:     facility.DailyVolume,
		Permit:          facility.Permit,
		WasteTypes:      wasteTypes,
		Collectors:      facility.Collectors,
		CreatedAt:       types.DateTime{Time: facility.CreatedAt},
		UpdatedAt:       types.DateTime{Time: facility.UpdatedAt},
		IsActive:        facility.IsActive,
	}
}
