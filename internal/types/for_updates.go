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

type UpdatePickupRequest struct {
	WasteType            string   `json:"waste_type"`
	Quantity             float64  `json:"quantity"`
	PickupDate           DateTime `json:"pickup_date"`
	Status               string   `json:"status" binding:"required,oneof=Pending Assigned Completed Cancelled"`
	HandlingRequirements string   `json:"handling_requirements"`
	AssignedDriver       int64    `json:"assigned_driver,omitempty"`
	AssignedVehicle      int64    `json:"assigned_vehicle,omitempty"`
	CreatedAt            DateTime `json:"created_at"`
}

type UpdateFacilityRequest struct {
	FacilityID      int64    `json:"facility_id"`
	Name            string   `json:"name,omitempty"`
	Location        string   `json:"location,omitempty"`
	Status          string   `json:"status,omitempty" binding:"omitempty,oneof=operational near-capacity maintenance closed"`
	Capacity        float64  `json:"capacity,omitempty" binding:"omitempty,min=0,max=100"`
	ComplianceScore int      `json:"compliance_score,omitempty" binding:"omitempty,min=0,max=100"`
	DailyVolume     float64  `json:"daily_volume,omitempty" binding:"omitempty,min=0"`
	Permit          string   `json:"permit,omitempty"`
	WasteTypes      []string `json:"waste_types,omitempty"`
	Collectors      int      `json:"collectors,omitempty"`
	IsActive        *bool    `json:"is_active,omitempty"`
}

type UpdateCollectorFacilityRequest struct {
	ProcessingVolume     float64   `json:"processing_volume,omitempty" binding:"omitempty,min=0"`
	HandlingRequirements string    `json:"handling_requirements,omitempty"`
	IsActive             *bool     `json:"is_active,omitempty"`
	LastProcessingDate   *DateTime `json:"last_processing_date,omitempty"`
}

type UpdateZoneRequest struct {
	ZoneID          int64    `json:"zone_id"`
	Name            string   `json:"name,omitempty"`
	Type            string   `json:"type,omitempty" binding:"omitempty,oneof=environmental cultural commercial security"`
	Status          string   `json:"status,omitempty" binding:"omitempty,oneof=permanent time-based"`
	Description     string   `json:"description,omitempty"`
	Area            string   `json:"area,omitempty"`
	ViolationsCount int      `json:"violations_count,omitempty"`
	Authority       string   `json:"authority,omitempty"`
	Restrictions    []string `json:"restrictions,omitempty"`
	IsActive        *bool    `json:"is_active,omitempty"`
}
