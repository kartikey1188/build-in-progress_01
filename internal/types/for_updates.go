package types

type UserUpdate struct {
	UserID       int64    `json:"user_id"`
	Email        string   `json:"email" binding:"required,email"`
	PasswordHash string   `json:"password_hash"`
	FullName     string   `json:"full_name"`
	PhoneNumber  string   `json:"phone_number"`
	Address      string   `json:"address"`
	Registration Date     `json:"registration_date"`
	Role         string   `json:"role" binding:"oneof=Business Collector Admin Government Driver"`
	IsActive     bool     `json:"is_active"`
	ProfileImage string   `json:"profile_image"`
	LastLogin    DateTime `json:"last_login"`
	IsVerified   bool     `json:"is_verified"`
	IsFlagged    bool     `json:"is_flagged"`
}

type BusinessUpdate struct {
	UserUpdate
	Business_name       string `json:"business_name"`
	Business_type       string `json:"business_type"`
	Registration_number string `json:"registration_number"`
	Gst_id              string `json:"gst_id"`
	Business_address    string `json:"business_address"`
}

type CollectorUpdate struct {
	UserUpdate
	Company_name   string `json:"company_name"`
	License_number string `json:"license_number"`
	Capacity       int64  `json:"capacity"`
	License_expiry Date   `json:"license_expiry"`
}
