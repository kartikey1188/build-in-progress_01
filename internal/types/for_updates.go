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
	CategoryID           int64   `json:"category_id"`
	CollectorID          int64   `json:"collector_id"`
	PricePerKg           float64 `json:"price_per_kg"`
	MaximumCapacity      float64 `json:"maximum_capacity"`
	HandlingRequirements string  `json:"handling_requirements"`
}

type UpdateCollectorVehicle struct {
	VehicleID            int64  `json:"vehicle_id"`
	VehicleNumber        string `json:"vehicle_number"`
	MaintenanceDate      Date   `json:"maintenance_date"`
	IsActive             *bool  `json:"is_active"` // changed from bool to *bool
	GPSTrackingID        string `json:"gps_tracking_id"`
	AssignedDriverID     int64  `json:"assigned_driver_id"`
	RegistrationDocument string `json:"registration_document"`
	RegistrationExpiry   Date   `json:"registration_expiry"`
}

type UpdateCollectorDriver struct {
	DriverID      int64   `json:"driver_id"`
	CollectorID   int64   `json:"collector_id"`
	LicenseNumber string  `json:"license_number"`
	DriverName    string  `json:"driver_name"`
	LicenseExpiry Date    `json:"license_expiry"`
	IsEmployed    *bool   `json:"is_employed"`
	IsActive      *bool   `json:"is_active"`
	Rating        float64 `json:"rating"`
	JoiningDate   Date    `json:"joining_date"`
}
