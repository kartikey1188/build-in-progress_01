package databaseone

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/kartikey1188/build-in-progress_01/internal/types"
)

// Collector operations
func (p *Postgres) GetCollectorDetails(id string) (types.Collector, error) {
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
