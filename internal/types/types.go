package types

import (
	"fmt"
	"time"
)

type User struct {
	UserID       int64    `json:"user_id"`
	Email        string   `json:"email" binding:"required,email"`
	PasswordHash string   `json:"password_hash" binding:"required"`
	FullName     string   `json:"full_name" binding:"required"`
	PhoneNumber  string   `json:"phone_number,omitempty"`
	Address      string   `json:"address,omitempty"`
	Registration Date     `json:"registration_date"`
	Role         string   `json:"role" binding:"required,oneof=Business Collector Admin Government Driver"`
	IsActive     bool     `json:"is_active"`
	ProfileImage string   `json:"profile_image,omitempty"`
	LastLogin    DateTime `json:"last_login"`
	IsVerified   bool     `json:"is_verified"`
	IsFlagged    bool     `json:"is_flagged"`
}

type Business struct {
	User
	Business_name       string `json:"business_name" binding:"required"`
	Business_type       string `json:"business_type" binding:"required"`
	Registration_number string `json:"registration_number" binding:"required"`
	Gst_id              string `json:"gst_id" binding:"required"`
	Business_address    string `json:"business_address" binding:"required"`
}

type Collector struct {
	User
	Company_name   string `json:"company_name" binding:"required"`
	License_number string `json:"license_number" binding:"required"`
	Capacity       int64  `json:"capacity" binding:"required"`
	License_expiry Date   `json:"license_expiry" binding:"required"`
}

type Date struct {
	time.Time
}

type DateTime struct {
	time.Time
}

const (
	dateLayout     = "2006-01-02"
	dateTimeLayout = "2006-01-02 15:04:05"
)

func (d *Date) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = s[1 : len(s)-1]
	t, err := time.Parse(dateLayout, s)
	if err != nil {
		return fmt.Errorf("invalid date format (expected YYYY-MM-DD): %w", err)
	}
	d.Time = t
	return nil
}

func (d Date) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", d.Time.Format(dateLayout))
	return []byte(formatted), nil
}

func (dt *DateTime) UnmarshalJSON(b []byte) error {
	s := string(b)
	s = s[1 : len(s)-1]
	t, err := time.Parse(dateTimeLayout, s)
	if err != nil {
		return fmt.Errorf("invalid datetime format (expected YYYY-MM-DD HH:MM:SS): %w", err)
	}
	dt.Time = t
	return nil
}

func (dt DateTime) MarshalJSON() ([]byte, error) {
	formatted := fmt.Sprintf("\"%s\"", dt.Time.Format(dateTimeLayout))
	return []byte(formatted), nil
}

type Notifiable interface {
	GetEmail() string
}

func (b Business) GetEmail() string {
	return b.Email
}

func (c Collector) GetEmail() string {
	return c.Email
}
