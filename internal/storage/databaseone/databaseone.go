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
