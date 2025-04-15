package databaseone

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/kartikey1188/build-in-progress_01/internal/config"
	"github.com/kartikey1188/build-in-progress_01/internal/types"
)

type Postgres struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Postgres, error) {
	db, err := sql.Open("pgx", cfg.StoragePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := createTables(db); err != nil {
		return nil, err
	}

	return &Postgres{Db: db}, nil
}

func createTables(db *sql.DB) error {
	tables := []string{
		`CREATE TABLE IF NOT EXISTS users (
			user_id SERIAL PRIMARY KEY,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			full_name TEXT NOT NULL,
			phone_number TEXT,
			address TEXT,
			registration_date DATE,
			role TEXT NOT NULL,
			is_active BOOLEAN,
			profile_image TEXT,
			last_login TIMESTAMP,
			is_verified BOOLEAN,
			is_flagged BOOLEAN
		)`,
		`INSERT INTO users (
			email, password_hash, full_name, role, is_active, is_verified, is_flagged, registration_date, last_login
		) VALUES (
		 	'admincontrols@gmail.com', 
			'admin123',
			'Application Admin',
			'Admin',
			TRUE,
			TRUE,
			FALSE,
			CURRENT_DATE,
			CURRENT_TIMESTAMP
		) ON CONFLICT (email) DO NOTHING`,
		`CREATE TABLE IF NOT EXISTS businesses (
			user_id INTEGER PRIMARY KEY REFERENCES users(user_id),
			business_name TEXT NOT NULL,
			business_type TEXT NOT NULL,
			registration_number TEXT NOT NULL,
			gst_id TEXT NOT NULL,
			business_address TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS collectors (
			user_id INTEGER PRIMARY KEY REFERENCES users(user_id),
			company_name TEXT NOT NULL,
			license_number TEXT NOT NULL,
			authorized_categories TEXT NOT NULL,
			capacity BIGINT NOT NULL,
			license_expiry DATE NOT NULL
		)`,
	}

	for _, query := range tables {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to execute table creation: %w", err)
		}
	}

	return nil
}

func (p *Postgres) CreateUser(user types.User) (int64, error) {
	var lastID int64
	err := p.Db.QueryRow(`
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

func (p *Postgres) GetUserByEmail(email string) (types.User, error) {
	var user types.User
	var registration, lastLogin time.Time

	err := p.Db.QueryRow(`
		SELECT user_id, email, password_hash, full_name, phone_number,
			address, registration_date, role, is_active, profile_image,
			last_login, is_verified, is_flagged
		FROM users
		WHERE email = $1
		LIMIT 1`, email,
	).Scan(
		&user.UserID, &user.Email, &user.PasswordHash, &user.FullName, &user.PhoneNumber,
		&user.Address, &registration, &user.Role, &user.IsActive, &user.ProfileImage,
		&lastLogin, &user.IsVerified, &user.IsFlagged,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.User{}, fmt.Errorf("user not found")
		}
		return types.User{}, fmt.Errorf("query error: %w", err)
	}

	user.Registration = types.Date{Time: registration}
	user.LastLogin = types.DateTime{Time: lastLogin}
	return user, nil
}

func (p *Postgres) UpdateLastLogin(userID int64, lastLogin types.DateTime) error {
	_, err := p.Db.Exec(`UPDATE users SET last_login = $1 WHERE user_id = $2`, lastLogin.Time, userID)
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
	_, err = p.Db.Exec(`
		INSERT INTO collectors (
			user_id, company_name, license_number, authorized_categories, capacity, license_expiry
		) VALUES ($1, $2, $3, $4, $5, $6)`,
		id, user.Company_name, user.License_number, user.Authorized_categories, user.Capacity, user.License_expiry.Time)
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
	_, err = p.Db.Exec(`
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
            c.company_name, c.license_number, c.authorized_categories, c.capacity, c.license_expiry
        FROM users u
        JOIN collectors c ON u.user_id = c.user_id
        WHERE u.email = $1
        LIMIT 1
    `
	err := p.Db.QueryRow(query, email).Scan(
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
	err := p.Db.QueryRow(query, email).Scan(
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

func (p *Postgres) GetUserById(userID int64) (types.User, error) {
	var user types.User
	var registration, lastLogin time.Time

	err := p.Db.QueryRow(`
		SELECT user_id, email, password_hash, full_name, phone_number,
			address, registration_date, role, is_active, profile_image,
			last_login, is_verified, is_flagged
		FROM users
		WHERE user_id = $1
		LIMIT 1`, userID,
	).Scan(
		&user.UserID, &user.Email, &user.PasswordHash, &user.FullName, &user.PhoneNumber,
		&user.Address, &registration, &user.Role, &user.IsActive, &user.ProfileImage,
		&lastLogin, &user.IsVerified, &user.IsFlagged,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.User{}, fmt.Errorf("user not found")
		}
		return types.User{}, fmt.Errorf("query error: %w", err)
	}

	user.Registration = types.Date{Time: registration}
	user.LastLogin = types.DateTime{Time: lastLogin}
	return user, nil
}
