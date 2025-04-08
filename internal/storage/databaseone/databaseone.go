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

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
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
	)`)
	if err != nil {
		return nil, fmt.Errorf("failed to create users table: %w", err)
	}

	return &Postgres{
		Db: db,
	}, nil
}

func (p *Postgres) CreateUser(user types.User) (int64, error) {
	var lastId int64
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
	).Scan(&lastId)
	if err != nil {
		return 0, fmt.Errorf("failed to insert user: %w", err)
	}
	return lastId, nil
}

func (p *Postgres) GetUserByEmail(email string) (types.User, error) {
	var user types.User
	var registration time.Time
	var lastLogin time.Time

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
