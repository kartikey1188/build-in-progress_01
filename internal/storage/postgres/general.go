package postgres

import (
	"fmt"

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
