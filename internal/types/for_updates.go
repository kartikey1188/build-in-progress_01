package types

type UserUpdate struct {
	UserID       int64    `json:"user_id"`
	Email        string   `json:"email"`
	PasswordHash string   `json:"password_hash"`
	FullName     string   `json:"full_name"`
	PhoneNumber  string   `json:"phone_number"`
	Address      string   `json:"address"`
	Registration Date     `json:"registration_date"`
	Role         string   `json:"role"`
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

type UpdateCollectorServiceCategory struct {
	CategoryID           int64   `json:"category_id"`           // Foreign key to ServiceCategory
	CollectorID          int64   `json:"collector_id"`          // Foreign key to Collector
	PricePerKg           float64 `json:"price_per_kg"`          // Cost per kg for this waste type
	MaximumCapacity      float64 `json:"maximum_capacity"`      // Max capacity for this waste type
	HandlingRequirements string  `json:"handling_requirements"` // Special handling instructions
}
