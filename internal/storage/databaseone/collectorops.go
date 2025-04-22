package databaseone

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/kartikey1188/build-in-progress_01/internal/types"
)

// Collector operations
func (p *Postgres) GetCollector(id string) (types.Collector, error) {
	var collector types.Collector
	var user types.User
	var registration, lastLogin, licenseExpiry time.Time

	query := `
        SELECT 
            u.user_id, u.email, u.password_hash, u.full_name, u.phone_number,
            u.address, u.registration_date, u.role, u.is_active, u.profile_image,
            u.last_login, u.is_verified, u.is_flagged,
            c.company_name, c.license_number, c.authorized_categories, c.capacity, c.license_expiry
        FROM users u
        JOIN collectors c ON u.user_id = c.user_id
        WHERE u.user_id = $1
        LIMIT 1
    `
	err := p.Db.QueryRow(query, id).Scan(
		&user.UserID, &user.Email, &user.PasswordHash, &user.FullName, &user.PhoneNumber,
		&user.Address, &registration, &user.Role, &user.IsActive, &user.ProfileImage,
		&lastLogin, &user.IsVerified, &user.IsFlagged,
		&collector.Company_name, &collector.License_number, &collector.Authorized_categories,
		&collector.Capacity, &licenseExpiry,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Collector{}, fmt.Errorf("collector not found")
		}
		return types.Collector{}, fmt.Errorf("query error: %w", err)
	}

	user.Registration = types.Date{Time: registration}
	user.LastLogin = types.DateTime{Time: lastLogin}
	collector.User = user
	collector.License_expiry = types.Date{Time: licenseExpiry}

	return collector, nil
}

func (p *Postgres) GetCollectors() ([]types.Collector, error) {
	query := `
        SELECT 
            u.user_id, u.email, u.full_name, u.phone_number, u.address,
            u.registration_date, u.profile_image, u.last_login,
            c.company_name, c.license_number, c.authorized_categories, c.capacity, c.license_expiry
        FROM users u
        JOIN collectors c ON u.user_id = c.user_id
        WHERE u.role = 'Collector'
    `
	rows, err := p.Db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query collectors: %w", err)
	}
	defer rows.Close()

	var collectors []types.Collector
	for rows.Next() {
		var c types.Collector
		var user types.User
		var regDate, lastLogin, licenseExp time.Time

		err := rows.Scan(
			&user.UserID, &user.Email, &user.FullName, &user.PhoneNumber, &user.Address,
			&regDate, &user.ProfileImage, &lastLogin,
			&c.Company_name, &c.License_number, &c.Authorized_categories,
			&c.Capacity, &licenseExp,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan collector: %w", err)
		}

		user.Registration = types.Date{Time: regDate}
		user.LastLogin = types.DateTime{Time: lastLogin}
		c.User = user
		c.License_expiry = types.Date{Time: licenseExp}
		collectors = append(collectors, c)
	}

	return collectors, nil
}

func (p *Postgres) UpdateProfile(userID int64, collector types.CollectorUpdate) (int64, error) {
	tx, err := p.Db.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// ==================== Users Table Update ====================
	var userSet []string
	var userParams []interface{}
	userParamIndex := 1

	// FullName
	if collector.FullName != nil {
		userSet = append(userSet, fmt.Sprintf("full_name = $%d", userParamIndex))
		userParams = append(userParams, *collector.FullName)
		userParamIndex++
	}

	// PhoneNumber
	if collector.PhoneNumber != nil {
		userSet = append(userSet, fmt.Sprintf("phone_number = $%d", userParamIndex))
		userParams = append(userParams, *collector.PhoneNumber)
		userParamIndex++
	}

	// Address
	if collector.Address != nil {
		userSet = append(userSet, fmt.Sprintf("address = $%d", userParamIndex))
		userParams = append(userParams, *collector.Address)
		userParamIndex++
	}

	// ProfileImage
	if collector.ProfileImage != nil {
		userSet = append(userSet, fmt.Sprintf("profile_image = $%d", userParamIndex))
		userParams = append(userParams, *collector.ProfileImage)
		userParamIndex++
	}

	// Execute users update if any fields are set
	if len(userSet) > 0 {
		userQuery := fmt.Sprintf(
			"UPDATE users SET %s WHERE user_id = $%d",
			strings.Join(userSet, ", "),
			userParamIndex,
		)
		userParams = append(userParams, userID)

		if _, err := tx.Exec(userQuery, userParams...); err != nil {
			return 0, fmt.Errorf("failed to update user: %w", err)
		}
	}

	// ==================== Collectors Table Update ====================
	var collectorSet []string
	var collectorParams []interface{}
	collectorParamIndex := 1

	// CompanyName
	if collector.CompanyName != nil {
		collectorSet = append(collectorSet, fmt.Sprintf("company_name = $%d", collectorParamIndex))
		collectorParams = append(collectorParams, *collector.CompanyName)
		collectorParamIndex++
	}

	// LicenseNumber
	if collector.LicenseNumber != nil {
		collectorSet = append(collectorSet, fmt.Sprintf("license_number = $%d", collectorParamIndex))
		collectorParams = append(collectorParams, *collector.LicenseNumber)
		collectorParamIndex++
	}

	// AuthorizedCategories
	if collector.AuthorizedCategories != nil {
		collectorSet = append(collectorSet, fmt.Sprintf("authorized_categories = $%d", collectorParamIndex))
		collectorParams = append(collectorParams, *collector.AuthorizedCategories)
		collectorParamIndex++
	}

	// Capacity
	if collector.Capacity != nil {
		collectorSet = append(collectorSet, fmt.Sprintf("capacity = $%d", collectorParamIndex))
		collectorParams = append(collectorParams, *collector.Capacity)
		collectorParamIndex++
	}

	// LicenseExpiry
	if collector.LicenseExpiry != nil {
		collectorSet = append(collectorSet, fmt.Sprintf("license_expiry = $%d", collectorParamIndex))
		collectorParams = append(collectorParams, collector.LicenseExpiry.Time)
		collectorParamIndex++
	}

	// Execute collectors update if any fields are set
	if len(collectorSet) > 0 {
		collectorQuery := fmt.Sprintf(
			"UPDATE collectors SET %s WHERE user_id = $%d",
			strings.Join(collectorSet, ", "),
			collectorParamIndex,
		)
		collectorParams = append(collectorParams, userID)

		if _, err := tx.Exec(collectorQuery, collectorParams...); err != nil {
			return 0, fmt.Errorf("failed to update collector: %w", err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return userID, nil
}

// Service Category operations
func (p *Postgres) CreateServiceCategory(csc types.CollectorServiceCategory) (int64, error) {
	var categoryID int64

	// First check if waste type exists
	err := p.Db.QueryRow(
		"SELECT category_id FROM service_categories WHERE waste_type = $1",
		csc.WasteType,
	).Scan(&categoryID)

	if err == sql.ErrNoRows {
		// Create new service category
		err = p.Db.QueryRow(
			"INSERT INTO service_categories (waste_type) VALUES ($1) RETURNING category_id",
			csc.WasteType,
		).Scan(&categoryID)
		if err != nil {
			return 0, fmt.Errorf("failed to create service category: %w", err)
		}
	} else if err != nil {
		return 0, fmt.Errorf("failed to check service categories: %w", err)
	}

	// Link to collector
	_, err = p.Db.Exec(`
		INSERT INTO collector_service_categories (
			category_id, collector_id, price_per_kg, maximum_capacity, handling_requirements
		) VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (category_id, collector_id) DO UPDATE SET
			price_per_kg = EXCLUDED.price_per_kg,
			maximum_capacity = EXCLUDED.maximum_capacity,
			handling_requirements = EXCLUDED.handling_requirements`,
		categoryID, csc.CollectorID, csc.PricePerKg, csc.MaximumCapacity, csc.HandlingRequirements,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to link service category: %w", err)
	}

	return categoryID, nil
}

func (p *Postgres) UpdateServiceCategory(id string, csc types.CollectorServiceCategory) error {
	_, err := p.Db.Exec(`
		UPDATE collector_service_categories SET
			price_per_kg = $1,
			maximum_capacity = $2,
			handling_requirements = $3
		WHERE category_id = $4 AND collector_id = $5`,
		csc.PricePerKg, csc.MaximumCapacity, csc.HandlingRequirements,
		id, csc.CollectorID,
	)
	if err != nil {
		return fmt.Errorf("failed to update service category: %w", err)
	}
	return nil
}

func (p *Postgres) DeleteServiceCategory(id string) error {
	_, err := p.Db.Exec(`
		DELETE FROM collector_service_categories 
		WHERE category_id = $1`,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to delete service category: %w", err)
	}
	return nil
}

// Vehicle operations
func (p *Postgres) AddVehicle(v types.CollectorVehicle) (int64, error) {
	var vehicleID int64
	err := p.Db.QueryRow(`
		INSERT INTO collector_vehicles (
			collector_id, vehicle_type, capacity, vehicle_number,
			maintenance_date, gps_tracking_id, registration_document, registration_expiry
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING vehicle_id`,
		v.CollectorID, v.VehicleType, v.Capacity, v.VehicleNumber,
		v.MaintenanceDate.Time, v.GPSTrackingID, v.RegistrationDocument, v.RegistrationExpiry.Time,
	).Scan(&vehicleID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert vehicle: %w", err)
	}
	return vehicleID, nil
}

func (p *Postgres) UpdateVehicle(id string, v types.CollectorVehicle) error {
	_, err := p.Db.Exec(`
		UPDATE collector_vehicles SET
			vehicle_type = $1,
			capacity = $2,
			vehicle_number = $3,
			maintenance_date = $4,
			gps_tracking_id = $5,
			registration_document = $6,
			registration_expiry = $7
		WHERE vehicle_id = $8`,
		v.VehicleType, v.Capacity, v.VehicleNumber,
		v.MaintenanceDate.Time, v.GPSTrackingID, v.RegistrationDocument, v.RegistrationExpiry.Time,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to update vehicle: %w", err)
	}
	return nil
}

func (p *Postgres) ActivateVehicle(id string) error {
	return p.updateVehicleStatus(id, true)
}

func (p *Postgres) DeactivateVehicle(id string) error {
	return p.updateVehicleStatus(id, false)
}

func (p *Postgres) updateVehicleStatus(id string, active bool) error {
	_, err := p.Db.Exec(`
		UPDATE collector_vehicles SET
			is_active = $1
		WHERE vehicle_id = $2`,
		active, id,
	)
	if err != nil {
		return fmt.Errorf("failed to update vehicle status: %w", err)
	}
	return nil
}

// Driver operations
func (p *Postgres) RegisterDriver(d types.CollectorDriver) (int64, error) {
	var driverID int64
	err := p.Db.QueryRow(`
		INSERT INTO collector_drivers (
			collector_id, license_number, license_expiry, 
			assigned_vehicle_id, is_active, rating, joining_date
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING driver_id`,
		d.CollectorID, d.LicenseNumber, d.LicenseExpiry.Time,
		d.AssignedVehicleID, d.IsActive, d.Rating, d.JoiningDate.Time,
	).Scan(&driverID)
	if err != nil {
		return 0, fmt.Errorf("failed to register driver: %w", err)
	}
	return driverID, nil
}

func (p *Postgres) UpdateDriver(id string, d types.CollectorDriver) error {
	_, err := p.Db.Exec(`
		UPDATE collector_drivers SET
			license_number = $1,
			license_expiry = $2,
			assigned_vehicle_id = $3,
			is_active = $4,
			rating = $5
		WHERE driver_id = $6`,
		d.LicenseNumber, d.LicenseExpiry.Time, d.AssignedVehicleID,
		d.IsActive, d.Rating, id,
	)
	if err != nil {
		return fmt.Errorf("failed to update driver: %w", err)
	}
	return nil
}

func (p *Postgres) AssignVehicleToDriver(driverID string, vehicleID int64) error {
	_, err := p.Db.Exec(`
		UPDATE collector_drivers SET
			assigned_vehicle_id = $1
		WHERE driver_id = $2`,
		vehicleID, driverID,
	)
	if err != nil {
		return fmt.Errorf("failed to assign vehicle: %w", err)
	}
	return nil
}

// Collector details
func (p *Postgres) GetCollectorServiceCategories(id string) ([]types.CollectorServiceCategory, error) {
	query := `
		SELECT csc.category_id, csc.collector_id, sc.waste_type, 
			csc.price_per_kg, csc.maximum_capacity, csc.handling_requirements
		FROM collector_service_categories csc
		JOIN service_categories sc ON csc.category_id = sc.category_id
		WHERE csc.collector_id = $1
	`
	rows, err := p.Db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to query service categories: %w", err)
	}
	defer rows.Close()

	var categories []types.CollectorServiceCategory
	for rows.Next() {
		var c types.CollectorServiceCategory
		err := rows.Scan(
			&c.CategoryID, &c.CollectorID, &c.WasteType,
			&c.PricePerKg, &c.MaximumCapacity, &c.HandlingRequirements,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan service category: %w", err)
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (p *Postgres) GetCollectorVehicles(id string) ([]types.CollectorVehicle, error) {
	query := `
		SELECT vehicle_id, collector_id, vehicle_type, capacity, vehicle_number,
			maintenance_date, is_active, gps_tracking_id, 
			registration_document, registration_expiry
		FROM collector_vehicles
		WHERE collector_id = $1
	`
	rows, err := p.Db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to query vehicles: %w", err)
	}
	defer rows.Close()

	var vehicles []types.CollectorVehicle
	for rows.Next() {
		var v types.CollectorVehicle
		var maintDate, regExpiry time.Time
		err := rows.Scan(
			&v.VehicleID, &v.CollectorID, &v.VehicleType, &v.Capacity, &v.VehicleNumber,
			&maintDate, &v.IsActive, &v.GPSTrackingID,
			&v.RegistrationDocument, &regExpiry,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan vehicle: %w", err)
		}
		v.MaintenanceDate = types.Date{Time: maintDate}
		v.RegistrationExpiry = types.Date{Time: regExpiry}
		vehicles = append(vehicles, v)
	}
	return vehicles, nil
}

func (p *Postgres) GetCollectorDrivers(id string) ([]types.CollectorDriver, error) {
	query := `
		SELECT driver_id, collector_id, license_number, license_expiry,
			assigned_vehicle_id, is_active, rating, joining_date
		FROM collector_drivers
		WHERE collector_id = $1
	`
	rows, err := p.Db.Query(query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to query drivers: %w", err)
	}
	defer rows.Close()

	var drivers []types.CollectorDriver
	for rows.Next() {
		var d types.CollectorDriver
		var licenseExp, joiningDate time.Time
		err := rows.Scan(
			&d.DriverID, &d.CollectorID, &d.LicenseNumber, &licenseExp,
			&d.AssignedVehicleID, &d.IsActive, &d.Rating, &joiningDate,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan driver: %w", err)
		}
		d.LicenseExpiry = types.Date{Time: licenseExp}
		d.JoiningDate = types.Date{Time: joiningDate}
		drivers = append(drivers, d)
	}
	return drivers, nil
}

func (p *Postgres) GetCollectorDetails(id string) (types.Collector, error) {
	return p.GetCollector(id)
}
