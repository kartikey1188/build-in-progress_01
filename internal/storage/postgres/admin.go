package postgres

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/kartikey1188/build-in-progress_01/internal/models"
	"github.com/kartikey1188/build-in-progress_01/internal/types"
	"gorm.io/gorm"
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

func (p *Postgres) GetAllPickupRequests() ([]types.PickupRequest, error) {
	var models []models.PickupRequest
	err := p.GormDB.Find(&models).Error
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	var requests []types.PickupRequest
	for _, model := range models {
		requests = append(requests, convertPickupRequestModelToType(model))
	}

	return requests, nil
}

// GetFacilities retrieves facilities with optional filtering
func (p *Postgres) GetFacilities(status string, location string, wasteType string) ([]types.Facility, error) {
	var facilityModels []models.Facility
	query := p.GormDB.Where("is_active = ?", true)

	// Apply filters if provided
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if location != "" {
		query = query.Where("location ILIKE ?", "%"+location+"%")
	}
	if wasteType != "" {
		query = query.Where("waste_types LIKE ?", "%"+wasteType+"%")
	}

	if err := query.Find(&facilityModels).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch facilities: %w", err)
	}

	facilities := make([]types.Facility, 0, len(facilityModels))
	for _, model := range facilityModels {
		facilities = append(facilities, convertFacilityModelToType(model))
	}

	return facilities, nil
}

// GetFacilityByID retrieves a specific facility by its ID
func (p *Postgres) GetFacilityByID(facilityID int64) (types.Facility, error) {
	var facilityModel models.Facility
	if err := p.GormDB.Where("facility_id = ? AND is_active = ?", facilityID, true).First(&facilityModel).Error; err != nil {
		return types.Facility{}, fmt.Errorf("facility with ID %d not found: %w", facilityID, err)
	}

	return convertFacilityModelToType(facilityModel), nil
}

// CreateFacility creates a new facility and returns its ID
func (p *Postgres) CreateFacility(facility types.Facility) (int64, error) {
	// Convert waste types slice to JSON string
	wasteTypesJSON, err := json.Marshal(facility.WasteTypes)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal waste types: %w", err)
	}

	facilityModel := models.Facility{
		Name:            facility.Name,
		Location:        facility.Location,
		Status:          facility.Status,
		Capacity:        facility.Capacity,
		ComplianceScore: facility.ComplianceScore,
		DailyVolume:     facility.DailyVolume,
		Permit:          facility.Permit,
		WasteTypes:      string(wasteTypesJSON),
		Collectors:      facility.Collectors,
		CreatedAt:       facility.CreatedAt.Time,
		UpdatedAt:       facility.UpdatedAt.Time,
		IsActive:        facility.IsActive,
	}

	if err := p.GormDB.Create(&facilityModel).Error; err != nil {
		return 0, fmt.Errorf("failed to create facility: %w", err)
	}

	return facilityModel.FacilityID, nil
}

// UpdateFacility updates an existing facility
func (p *Postgres) UpdateFacility(facilityID int64, facility types.UpdateFacilityRequest) error {
	updates := make(map[string]interface{})
	updates["updated_at"] = time.Now()

	// Add non-empty fields to updates
	if facility.Name != "" {
		updates["name"] = facility.Name
	}
	if facility.Location != "" {
		updates["location"] = facility.Location
	}
	if facility.Status != "" {
		updates["status"] = facility.Status
	}
	if facility.Capacity != 0 {
		updates["capacity"] = facility.Capacity
	}
	if facility.ComplianceScore != 0 {
		updates["compliance_score"] = facility.ComplianceScore
	}
	if facility.DailyVolume != 0 {
		updates["daily_volume"] = facility.DailyVolume
	}
	if facility.Permit != "" {
		updates["permit"] = facility.Permit
	}
	if len(facility.WasteTypes) > 0 {
		wasteTypesJSON, err := json.Marshal(facility.WasteTypes)
		if err != nil {
			return fmt.Errorf("failed to marshal waste types: %w", err)
		}
		updates["waste_types"] = string(wasteTypesJSON)
	}
	if facility.Collectors != 0 {
		updates["collectors"] = facility.Collectors
	}
	if facility.IsActive != nil {
		updates["is_active"] = *facility.IsActive
	}

	result := p.GormDB.Model(&models.Facility{}).Where("facility_id = ?", facilityID).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update facility: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("facility with ID %d not found", facilityID)
	}

	return nil
}

// DeleteFacility marks a facility as inactive (soft delete)
func (p *Postgres) DeleteFacility(facilityID int64) error {
	result := p.GormDB.Model(&models.Facility{}).Where("facility_id = ?", facilityID).Updates(map[string]interface{}{
		"is_active":  false,
		"updated_at": time.Now(),
	})

	if result.Error != nil {
		return fmt.Errorf("failed to delete facility: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("facility with ID %d not found", facilityID)
	}

	return nil
}

// AssignCollectorToFacility assigns a collector to a facility
func (p *Postgres) AssignCollectorToFacility(facilityID int64, collectorFacility types.CollectorFacility) error {
	// Check if facility exists and is active
	var facility models.Facility
	if err := p.GormDB.Where("facility_id = ? AND is_active = ?", facilityID, true).First(&facility).Error; err != nil {
		return fmt.Errorf("facility not found or inactive: %w", err)
	}

	// Check if collector exists and is active
	var collector models.Collector
	if err := p.GormDB.Where("user_id = ?", collectorFacility.CollectorID).First(&collector).Error; err != nil {
		return fmt.Errorf("collector not found: %w", err)
	}

	// Create the relationship
	collectorFacilityModel := models.CollectorFacility{
		CollectorID:          collectorFacility.CollectorID,
		FacilityID:           facilityID,
		AssignmentDate:       collectorFacility.AssignmentDate.Time,
		ProcessingVolume:     collectorFacility.ProcessingVolume,
		HandlingRequirements: collectorFacility.HandlingRequirements,
		IsActive:             collectorFacility.IsActive,
	}

	if err := p.GormDB.Create(&collectorFacilityModel).Error; err != nil {
		return fmt.Errorf("failed to assign collector to facility: %w", err)
	}

	// Update facility collectors count
	if err := p.GormDB.Model(&facility).Update("collectors", gorm.Expr("collectors + ?", 1)).Error; err != nil {
		return fmt.Errorf("failed to update facility collectors count: %w", err)
	}

	return nil
}

// UpdateCollectorFacility updates a collector-facility relationship
func (p *Postgres) UpdateCollectorFacility(facilityID int64, collectorID int64, request types.UpdateCollectorFacilityRequest) error {
	updates := make(map[string]interface{})

	if request.ProcessingVolume != 0 {
		updates["processing_volume"] = request.ProcessingVolume
	}
	if request.HandlingRequirements != "" {
		updates["handling_requirements"] = request.HandlingRequirements
	}
	if request.IsActive != nil {
		updates["is_active"] = *request.IsActive
	}
	if request.LastProcessingDate != nil {
		updates["last_processing_date"] = request.LastProcessingDate.Time
	}

	result := p.GormDB.Model(&models.CollectorFacility{}).
		Where("facility_id = ? AND collector_id = ?", facilityID, collectorID).
		Updates(updates)

	if result.Error != nil {
		return fmt.Errorf("failed to update collector-facility relationship: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("collector-facility relationship not found")
	}

	return nil
}

// RemoveCollectorFromFacility removes a collector from a facility (soft delete)
func (p *Postgres) RemoveCollectorFromFacility(facilityID int64, collectorID int64) error {
	tx := p.GormDB.Begin()

	// Soft delete the relationship
	if err := tx.Model(&models.CollectorFacility{}).
		Where("facility_id = ? AND collector_id = ?", facilityID, collectorID).
		Update("is_active", false).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to deactivate collector-facility relationship: %w", err)
	}

	// Update facility collectors count
	if err := tx.Model(&models.Facility{}).
		Where("facility_id = ?", facilityID).
		Update("collectors", gorm.Expr("collectors - ?", 1)).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update facility collectors count: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetFacilityCollectors gets all collectors assigned to a facility
func (p *Postgres) GetFacilityCollectors(facilityID int64) ([]types.CollectorFacility, error) {
	var collectorFacilities []models.CollectorFacility
	if err := p.GormDB.Where("facility_id = ?", facilityID).Find(&collectorFacilities).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch facility collectors: %w", err)
	}

	result := make([]types.CollectorFacility, len(collectorFacilities))
	for i, cf := range collectorFacilities {
		result[i] = types.CollectorFacility{
			CollectorID:          cf.CollectorID,
			FacilityID:           cf.FacilityID,
			AssignmentDate:       types.DateTime{Time: cf.AssignmentDate},
			LastProcessingDate:   types.DateTime{Time: cf.LastProcessingDate},
			ProcessingVolume:     cf.ProcessingVolume,
			HandlingRequirements: cf.HandlingRequirements,
			IsActive:             cf.IsActive,
		}
	}

	return result, nil
}

// GetCollectorFacilities gets all facilities a collector is assigned to
func (p *Postgres) GetCollectorFacilities(collectorID int64) ([]types.CollectorFacility, error) {
	var collectorFacilities []models.CollectorFacility
	if err := p.GormDB.Where("collector_id = ?", collectorID).Find(&collectorFacilities).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch collector facilities: %w", err)
	}

	result := make([]types.CollectorFacility, len(collectorFacilities))
	for i, cf := range collectorFacilities {
		result[i] = types.CollectorFacility{
			CollectorID:          cf.CollectorID,
			FacilityID:           cf.FacilityID,
			AssignmentDate:       types.DateTime{Time: cf.AssignmentDate},
			LastProcessingDate:   types.DateTime{Time: cf.LastProcessingDate},
			ProcessingVolume:     cf.ProcessingVolume,
			HandlingRequirements: cf.HandlingRequirements,
			IsActive:             cf.IsActive,
		}
	}

	return result, nil
}

// GetZones retrieves zones with optional filtering by type and status
func (p *Postgres) GetZones(zoneType string, status string) ([]types.Zone, error) {
	var zones []models.Zone
	query := p.GormDB.Where("is_active = ?", true)

	if zoneType != "" {
		query = query.Where("type = ?", zoneType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&zones).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch zones: %w", err)
	}

	result := make([]types.Zone, len(zones))
	for i, zone := range zones {
		// Parse restrictions from JSON string to slice
		var restrictions []string
		if zone.Restrictions != "" {
			if err := json.Unmarshal([]byte(zone.Restrictions), &restrictions); err != nil {
				restrictions = []string{} // Default to empty slice if parsing fails
			}
		}

		result[i] = types.Zone{
			ZoneID:          zone.ZoneID,
			Name:            zone.Name,
			Type:            zone.Type,
			Status:          zone.Status,
			Description:     zone.Description,
			Area:            zone.Area,
			ViolationsCount: zone.ViolationsCount,
			Authority:       zone.Authority,
			Restrictions:    restrictions,
			CreatedAt:       types.DateTime{Time: zone.CreatedAt},
			UpdatedAt:       types.DateTime{Time: zone.UpdatedAt},
			IsActive:        zone.IsActive,
		}
	}

	return result, nil
}

// GetZoneByID retrieves a zone by its ID
func (p *Postgres) GetZoneByID(zoneID int64) (types.Zone, error) {
	var zone models.Zone
	if err := p.GormDB.Where("zone_id = ? AND is_active = ?", zoneID, true).First(&zone).Error; err != nil {
		return types.Zone{}, fmt.Errorf("failed to fetch zone: %w", err)
	}

	// Parse restrictions from JSON string to slice
	var restrictions []string
	if zone.Restrictions != "" {
		if err := json.Unmarshal([]byte(zone.Restrictions), &restrictions); err != nil {
			restrictions = []string{} // Default to empty slice if parsing fails
		}
	}

	result := types.Zone{
		ZoneID:          zone.ZoneID,
		Name:            zone.Name,
		Type:            zone.Type,
		Status:          zone.Status,
		Description:     zone.Description,
		Area:            zone.Area,
		ViolationsCount: zone.ViolationsCount,
		Authority:       zone.Authority,
		Restrictions:    restrictions,
		CreatedAt:       types.DateTime{Time: zone.CreatedAt},
		UpdatedAt:       types.DateTime{Time: zone.UpdatedAt},
		IsActive:        zone.IsActive,
	}

	return result, nil
}

// CreateZone creates a new zone and returns its ID
func (p *Postgres) CreateZone(zone types.Zone) (int64, error) {
	// Convert restrictions slice to JSON string
	restrictionsJSON, err := json.Marshal(zone.Restrictions)
	if err != nil {
		return 0, fmt.Errorf("failed to marshal restrictions: %w", err)
	}

	zoneModel := models.Zone{
		Name:            zone.Name,
		Type:            zone.Type,
		Status:          zone.Status,
		Description:     zone.Description,
		Area:            zone.Area,
		ViolationsCount: zone.ViolationsCount,
		Authority:       zone.Authority,
		Restrictions:    string(restrictionsJSON),
		CreatedAt:       zone.CreatedAt.Time,
		UpdatedAt:       zone.UpdatedAt.Time,
		IsActive:        zone.IsActive,
	}

	if err := p.GormDB.Create(&zoneModel).Error; err != nil {
		return 0, fmt.Errorf("failed to create zone: %w", err)
	}

	return zoneModel.ZoneID, nil
}

// UpdateZone updates an existing zone
func (p *Postgres) UpdateZone(zoneID int64, zone types.UpdateZoneRequest) error {
	updates := make(map[string]interface{})

	if zone.Name != "" {
		updates["name"] = zone.Name
	}
	if zone.Type != "" {
		updates["type"] = zone.Type
	}
	if zone.Status != "" {
		updates["status"] = zone.Status
	}
	if zone.Description != "" {
		updates["description"] = zone.Description
	}
	if zone.Area != "" {
		updates["area"] = zone.Area
	}
	if zone.ViolationsCount != 0 {
		updates["violations_count"] = zone.ViolationsCount
	}
	if zone.Authority != "" {
		updates["authority"] = zone.Authority
	}
	if zone.Restrictions != nil {
		restrictionsJSON, err := json.Marshal(zone.Restrictions)
		if err != nil {
			return fmt.Errorf("failed to marshal restrictions: %w", err)
		}
		updates["restrictions"] = string(restrictionsJSON)
	}
	if zone.IsActive != nil {
		updates["is_active"] = *zone.IsActive
	}

	updates["updated_at"] = time.Now()

	result := p.GormDB.Model(&models.Zone{}).Where("zone_id = ?", zoneID).Updates(updates)

	if result.Error != nil {
		return fmt.Errorf("failed to update zone: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("zone with ID %d not found", zoneID)
	}

	return nil
}

// DeleteZone soft deletes a zone (sets is_active to false)
func (p *Postgres) DeleteZone(zoneID int64) error {
	result := p.GormDB.Model(&models.Zone{}).Where("zone_id = ?", zoneID).Updates(map[string]interface{}{
		"is_active":  false,
		"updated_at": time.Now(),
	})

	if result.Error != nil {
		return fmt.Errorf("failed to delete zone: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("zone with ID %d not found", zoneID)
	}

	return nil
}
