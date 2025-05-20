package postgres

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/kartikey1188/build-in-progress_01/internal/types"
)

func (p *Postgres) CreateUser(user types.User) (int64, error) {
	var lastID int64
	err := p.SqlDB.QueryRow(`
		INSERT INTO users (
			email, password_hash, full_name, phone_number, address,
			registration_date, role, is_active, profile_image,
			last_login, is_verified, is_flagged
		)
		VALUES ($1, $2, $3, $4, $5,
				$6, $7, $8, $9,
				$10, $11, $12)
		RETURNING user_id`,
		user.Email, user.PasswordHash, user.FullName, user.PhoneNumber, user.Address,
		user.Registration.Time, user.Role, user.IsActive, user.ProfileImage,
		user.LastLogin.Time, user.IsVerified, user.IsFlagged,
	).Scan(&lastID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert user: %w", err)
	}
	return lastID, nil
}

func (p *Postgres) UpdateLastLogin(userID int64, lastLogin types.DateTime) error {
	_, err := p.SqlDB.Exec(`UPDATE users SET last_login = $1 WHERE user_id = $2`, lastLogin.Time, userID)
	if err != nil {
		return fmt.Errorf("failed to update last login: %w", err)
	}
	return nil
}

func (p *Postgres) CreateCollectorUser(user types.Collector) (int64, error) {
	id, err := p.CreateUser(user.User)
	if err != nil {
		return 0, err
	}
	_, err = p.SqlDB.Exec(`
		INSERT INTO collectors (
			user_id, company_name, license_number, capacity, license_expiry
		) VALUES ($1, $2, $3, $4, $5)`,
		id, user.Company_name, user.License_number, user.Capacity, user.License_expiry.Time)
	if err != nil {
		return 0, fmt.Errorf("failed to insert collector: %w", err)
	}
	return id, nil
}

func (p *Postgres) CreateBusinessUser(user types.Business) (int64, error) {
	id, err := p.CreateUser(user.User)
	if err != nil {
		return 0, err
	}
	_, err = p.SqlDB.Exec(`
		INSERT INTO businesses (
			user_id, business_name, business_type, registration_number, gst_id, business_address
		) VALUES ($1, $2, $3, $4, $5, $6)`,
		id, user.Business_name, user.Business_type, user.Registration_number, user.Gst_id, user.Business_address)
	if err != nil {
		return 0, fmt.Errorf("failed to insert business: %w", err)
	}
	return id, nil
}

func (p *Postgres) GetCollectorByEmail(email string) (types.Collector, error) {
	var collector types.Collector
	var user types.User
	var registration, lastLogin, licenseExpiry time.Time

	query := `
        SELECT 
            u.user_id, u.email, u.password_hash, u.full_name, u.phone_number,
            u.address, u.registration_date, u.role, u.is_active, u.profile_image,
            u.last_login, u.is_verified, u.is_flagged,
            c.company_name, c.license_number, c.capacity, c.license_expiry
        FROM users u
        JOIN collectors c ON u.user_id = c.user_id
        WHERE u.email = $1
        LIMIT 1
    `
	err := p.SqlDB.QueryRow(query, email).Scan(
		&user.UserID, &user.Email, &user.PasswordHash, &user.FullName, &user.PhoneNumber,
		&user.Address, &registration, &user.Role, &user.IsActive, &user.ProfileImage,
		&lastLogin, &user.IsVerified, &user.IsFlagged,
		&collector.Company_name, &collector.License_number,
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

func (p *Postgres) GetBusinessByEmail(email string) (types.Business, error) {
	var business types.Business
	var user types.User
	var registration, lastLogin time.Time

	query := `
        SELECT 
            u.user_id, u.email, u.password_hash, u.full_name, u.phone_number,
            u.address, u.registration_date, u.role, u.is_active, u.profile_image,
            u.last_login, u.is_verified, u.is_flagged,
            b.business_name, b.business_type, b.registration_number, b.gst_id, b.business_address
        FROM users u
        JOIN businesses b ON u.user_id = b.user_id
        WHERE u.email = $1
        LIMIT 1
    `
	err := p.SqlDB.QueryRow(query, email).Scan(
		&user.UserID, &user.Email, &user.PasswordHash, &user.FullName, &user.PhoneNumber,
		&user.Address, &registration, &user.Role, &user.IsActive, &user.ProfileImage,
		&lastLogin, &user.IsVerified, &user.IsFlagged,
		&business.Business_name, &business.Business_type, &business.Registration_number,
		&business.Gst_id, &business.Business_address,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Business{}, fmt.Errorf("business not found")
		}
		return types.Business{}, fmt.Errorf("query error: %w", err)
	}

	user.Registration = types.Date{Time: registration}
	user.LastLogin = types.DateTime{Time: lastLogin}
	business.User = user

	return business, nil
}

func (p *Postgres) CreateCollectorDriver(driver types.CollectorDriver, collectorID int64) (int64, error) {
	userID, err := p.CreateUser(driver.User)
	if err != nil {
		return 0, fmt.Errorf("failed to create driver user: %w", err)
	}

	_, err = p.SqlDB.Exec(`
        INSERT INTO collector_drivers (
            user_id, collector_id, license_number, driver_name, 
            license_expiry, is_employed, is_active, rating, joining_date
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		userID,
		collectorID,
		driver.LicenseNumber,
		driver.DriverName,
		driver.LicenseExpiry.Time,
		driver.IsEmployed,
		driver.IsActive,
		driver.Rating,
		driver.JoiningDate.Time,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to create collector driver: %w", err)
	}

	return userID, nil
}

func (p *Postgres) GetCollectorDriverByEmail(email string) (types.CollectorDriver, error) {
	var driver types.CollectorDriver
	var user types.User
	var registration, lastLogin, licenseExpiry, joiningDate time.Time

	query := `
        SELECT 
            u.user_id, u.email, u.password_hash, u.full_name, u.phone_number,
            u.address, u.registration_date, u.role, u.is_active, u.profile_image,
            u.last_login, u.is_verified, u.is_flagged,
            cd.collector_id, cd.license_number, cd.driver_name,
            cd.license_expiry, cd.is_employed, cd.is_active,
            cd.rating, cd.joining_date
        FROM users u
        JOIN collector_drivers cd ON u.user_id = cd.user_id
        WHERE u.email = $1
        LIMIT 1`

	err := p.SqlDB.QueryRow(query, email).Scan(
		&user.UserID, &user.Email, &user.PasswordHash, &user.FullName, &user.PhoneNumber,
		&user.Address, &registration, &user.Role, &user.IsActive, &user.ProfileImage,
		&lastLogin, &user.IsVerified, &user.IsFlagged,
		&driver.CollectorID, &driver.LicenseNumber, &driver.DriverName,
		&licenseExpiry, &driver.IsEmployed, &driver.IsActive,
		&driver.Rating, &joiningDate,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return types.CollectorDriver{}, fmt.Errorf("driver not found")
		}
		return types.CollectorDriver{}, fmt.Errorf("failed to get driver: %w", err)
	}

	user.Registration = types.Date{Time: registration}
	user.LastLogin = types.DateTime{Time: lastLogin}
	driver.LicenseExpiry = types.Date{Time: licenseExpiry}
	driver.JoiningDate = types.Date{Time: joiningDate}
	driver.User = user

	return driver, nil
}
