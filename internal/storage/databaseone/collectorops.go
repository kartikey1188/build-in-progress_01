package databaseone

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/kartikey1188/build-in-progress_01/internal/types"
)

func (p *Postgres) GetCollectorByID(id int64) (types.Collector, error) {
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

	var userSet []string
	var userParams []interface{}
	userParamIndex := 1

	if collector.FullName != nil {
		userSet = append(userSet, fmt.Sprintf("full_name = $%d", userParamIndex))
		userParams = append(userParams, *collector.FullName)
		userParamIndex++
	}

	if collector.PhoneNumber != nil {
		userSet = append(userSet, fmt.Sprintf("phone_number = $%d", userParamIndex))
		userParams = append(userParams, *collector.PhoneNumber)
		userParamIndex++
	}

	if collector.Address != nil {
		userSet = append(userSet, fmt.Sprintf("address = $%d", userParamIndex))
		userParams = append(userParams, *collector.Address)
		userParamIndex++
	}

	if collector.ProfileImage != nil {
		userSet = append(userSet, fmt.Sprintf("profile_image = $%d", userParamIndex))
		userParams = append(userParams, *collector.ProfileImage)
		userParamIndex++
	}

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

	var collectorSet []string
	var collectorParams []interface{}
	collectorParamIndex := 1

	if collector.CompanyName != nil {
		collectorSet = append(collectorSet, fmt.Sprintf("company_name = $%d", collectorParamIndex))
		collectorParams = append(collectorParams, *collector.CompanyName)
		collectorParamIndex++
	}

	if collector.LicenseNumber != nil {
		collectorSet = append(collectorSet, fmt.Sprintf("license_number = $%d", collectorParamIndex))
		collectorParams = append(collectorParams, *collector.LicenseNumber)
		collectorParamIndex++
	}

	if collector.AuthorizedCategories != nil {
		collectorSet = append(collectorSet, fmt.Sprintf("authorized_categories = $%d", collectorParamIndex))
		collectorParams = append(collectorParams, *collector.AuthorizedCategories)
		collectorParamIndex++
	}

	if collector.Capacity != nil {
		collectorSet = append(collectorSet, fmt.Sprintf("capacity = $%d", collectorParamIndex))
		collectorParams = append(collectorParams, *collector.Capacity)
		collectorParamIndex++
	}

	if collector.LicenseExpiry != nil {
		collectorSet = append(collectorSet, fmt.Sprintf("license_expiry = $%d", collectorParamIndex))
		collectorParams = append(collectorParams, collector.LicenseExpiry.Time)
		collectorParamIndex++
	}

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

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return userID, nil
}

func (p *Postgres) AddCollectorServiceCategory(input types.CollectorServiceCategory) (int64, error) {
	query := `
        INSERT INTO collector_service_categories (category_id, collector_id, waste_type, price_per_kg, maximum_capacity, handling_requirements)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id
    `
	var id int64
	err := p.Db.QueryRow(query,
		input.CategoryID,
		input.CollectorID,
		input.WasteType,
		input.PricePerKg,
		input.MaximumCapacity,
		input.HandlingRequirements,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to add collector service category: %w", err)
	}
	return id, nil
}

func (p *Postgres) UpdateCollectorServiceCategory(id int64, collectorID int64, input types.CollectorServiceCategory) error {
	query := `
        UPDATE collector_service_categories
        SET category_id = $1, waste_type = $2, price_per_kg = $3, maximum_capacity = $4, handling_requirements = $5
        WHERE id = $6 AND collector_id = $7
    `
	_, err := p.Db.Exec(query,
		input.CategoryID,
		input.WasteType,
		input.PricePerKg,
		input.MaximumCapacity,
		input.HandlingRequirements,
		id,
		collectorID,
	)
	if err != nil {
		return fmt.Errorf("failed to update collector service category: %w", err)
	}
	return nil
}

func (p *Postgres) DeleteCollectorServiceCategory(id int64, collectorID int64) error {
	query := `
        DELETE FROM collector_service_categories
        WHERE id = $1 AND collector_id = $2
    `
	_, err := p.Db.Exec(query, id, collectorID)
	if err != nil {
		return fmt.Errorf("failed to delete collector service category: %w", err)
	}
	return nil
}

func (p *Postgres) AddCollectorVehicle(input types.CollectorVehicle) (int64, error) {
	query := `
        INSERT INTO collector_vehicles (vehicle_id, collector_id, vehicle_type, capacity, vehicle_number, maintenance_date, is_active, gps_tracking_id, assigned_driver_id, registration_document, registration_expiry)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
        RETURNING id
    `
	var id int64
	err := p.Db.QueryRow(query,
		input.VehicleID,
		input.CollectorID,
		input.VehicleType,
		input.Capacity,
		input.VehicleNumber,
		input.MaintenanceDate.Time,
		input.IsActive,
		input.GPSTrackingID,
		input.AssignedDriverID,
		input.RegistrationDocument,
		input.RegistrationExpiry.Time,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to add collector vehicle: %w", err)
	}
	return id, nil
}

func (p *Postgres) UpdateCollectorVehicle(id int64, collectorID int64, input types.CollectorVehicle) error {
	query := `
        UPDATE collector_vehicles
        SET vehicle_id = $1, vehicle_type = $2, capacity = $3, vehicle_number = $4, maintenance_date = $5, is_active = $6, gps_tracking_id = $7, assigned_driver_id = $8, registration_document = $9, registration_expiry = $10
        WHERE id = $11 AND collector_id = $12
    `
	_, err := p.Db.Exec(query,
		input.VehicleID,
		input.VehicleType,
		input.Capacity,
		input.VehicleNumber,
		input.MaintenanceDate.Time,
		input.IsActive,
		input.GPSTrackingID,
		input.AssignedDriverID,
		input.RegistrationDocument,
		input.RegistrationExpiry.Time,
		id,
		collectorID,
	)
	if err != nil {
		return fmt.Errorf("failed to update collector vehicle: %w", err)
	}
	return nil
}

func (p *Postgres) ActivateCollectorVehicle(id int64, collectorID int64) error {
	query := `
        UPDATE collector_vehicles
        SET is_active = true
        WHERE id = $1 AND collector_id = $2
    `
	_, err := p.Db.Exec(query, id, collectorID)
	if err != nil {
		return fmt.Errorf("failed to activate collector vehicle: %w", err)
	}
	return nil
}

func (p *Postgres) DeactivateCollectorVehicle(id int64, collectorID int64) error {
	query := `
        UPDATE collector_vehicles
        SET is_active = false
        WHERE id = $1 AND collector_id = $2
    `
	_, err := p.Db.Exec(query, id, collectorID)
	if err != nil {
		return fmt.Errorf("failed to deactivate collector vehicle: %w", err)
	}
	return nil
}

func (p *Postgres) AddCollectorDriver(input types.CollectorDriver) (int64, error) {
	query := `
        INSERT INTO collector_drivers (collector_id, license_number, license_expiry, assigned_vehicle_id, is_employed, is_active, rating, joining_date)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING driver_id
    `
	var driverID int64
	err := p.Db.QueryRow(query,
		input.CollectorID,
		input.LicenseNumber,
		input.LicenseExpiry.Time,
		input.AssignedVehicleID,
		input.IsEmployed,
		input.IsActive,
		input.Rating,
		input.JoiningDate.Time,
	).Scan(&driverID)
	if err != nil {
		return 0, fmt.Errorf("failed to add collector driver: %w", err)
	}
	return driverID, nil
}

func (p *Postgres) UpdateCollectorDriver(id int64, collectorID int64, input types.CollectorDriver) error {
	query := `
        UPDATE collector_drivers
        SET license_number = $1, license_expiry = $2, assigned_vehicle_id = $3, is_employed = $4, is_active = $5, rating = $6, joining_date = $7
        WHERE driver_id = $8 AND collector_id = $9
    `
	_, err := p.Db.Exec(query,
		input.LicenseNumber,
		input.LicenseExpiry.Time,
		input.AssignedVehicleID,
		input.IsEmployed,
		input.IsActive,
		input.Rating,
		input.JoiningDate.Time,
		id,
		collectorID,
	)
	if err != nil {
		return fmt.Errorf("failed to update collector driver: %w", err)
	}
	return nil
}

func (p *Postgres) AssignVehicleToDriver(driverID int64, vehicleID int64, collectorID int64) error {
	tx, err := p.Db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
        UPDATE collector_drivers
        SET assigned_vehicle_id = $1
        WHERE driver_id = $2 AND collector_id = $3
    `, vehicleID, driverID, collectorID)
	if err != nil {
		return fmt.Errorf("failed to assign vehicle to driver: %w", err)
	}

	_, err = tx.Exec(`
        UPDATE collector_vehicles
        SET assigned_driver_id = $1
        WHERE id = $2 AND collector_id = $3
    `, driverID, vehicleID, collectorID)
	if err != nil {
		return fmt.Errorf("failed to assign driver to vehicle: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}

func (p *Postgres) GetCollectorServiceCategories(collectorID int64) ([]types.CollectorServiceCategory, error) {
	query := `
        SELECT category_id, collector_id, waste_type, price_per_kg, maximum_capacity, handling_requirements
        FROM collector_service_categories
        WHERE collector_id = $1
    `
	rows, err := p.Db.Query(query, collectorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get collector service categories: %w", err)
	}
	defer rows.Close()

	var categories []types.CollectorServiceCategory
	for rows.Next() {
		var c types.CollectorServiceCategory
		err := rows.Scan(
			&c.CategoryID,
			&c.CollectorID,
			&c.WasteType,
			&c.PricePerKg,
			&c.MaximumCapacity,
			&c.HandlingRequirements,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan collector service category: %w", err)
		}
		categories = append(categories, c)
	}
	return categories, nil
}

func (p *Postgres) GetCollectorVehicles(collectorID int64) ([]types.CollectorVehicle, error) {
	query := `
        SELECT vehicle_id, collector_id, vehicle_type, capacity, vehicle_number, maintenance_date, is_active, gps_tracking_id, assigned_driver_id, registration_document, registration_expiry
        FROM collector_vehicles
        WHERE collector_id = $1
    `
	rows, err := p.Db.Query(query, collectorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get collector vehicles: %w", err)
	}
	defer rows.Close()

	var vehicles []types.CollectorVehicle
	for rows.Next() {
		var v types.CollectorVehicle
		var maintenanceDate, registrationExpiry time.Time
		err := rows.Scan(
			&v.VehicleID,
			&v.CollectorID,
			&v.VehicleType,
			&v.Capacity,
			&v.VehicleNumber,
			&maintenanceDate,
			&v.IsActive,
			&v.GPSTrackingID,
			&v.AssignedDriverID,
			&v.RegistrationDocument,
			&registrationExpiry,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan collector vehicle: %w", err)
		}
		v.MaintenanceDate = types.Date{Time: maintenanceDate}
		v.RegistrationExpiry = types.Date{Time: registrationExpiry}
		vehicles = append(vehicles, v)
	}
	return vehicles, nil
}

func (p *Postgres) GetCollectorDrivers(collectorID int64) ([]types.CollectorDriver, error) {
	query := `
        SELECT driver_id, collector_id, license_number, license_expiry, assigned_vehicle_id, is_employed, is_active, rating, joining_date
        FROM collector_drivers
        WHERE collector_id = $1
    `
	rows, err := p.Db.Query(query, collectorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get collector drivers: %w", err)
	}
	defer rows.Close()

	var drivers []types.CollectorDriver
	for rows.Next() {
		var d types.CollectorDriver
		var licenseExpiry, joiningDate time.Time
		err := rows.Scan(
			&d.DriverID,
			&d.CollectorID,
			&d.LicenseNumber,
			&licenseExpiry,
			&d.AssignedVehicleID,
			&d.IsEmployed,
			&d.IsActive,
			&d.Rating,
			&joiningDate,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan collector driver: %w", err)
		}
		d.LicenseExpiry = types.Date{Time: licenseExpiry}
		d.JoiningDate = types.Date{Time: joiningDate}
		drivers = append(drivers, d)
	}
	return drivers, nil
}
