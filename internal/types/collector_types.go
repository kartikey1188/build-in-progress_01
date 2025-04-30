package types

type ServiceCategory struct {
	CategoryID int64  `json:"category_id"`                   // Primary key
	WasteType  string `json:"waste_type" binding:"required"` // Type of waste accepted
}

type CollectorServiceCategory struct {
	CategoryID           int64   `json:"category_id"`                         // Foreign key to ServiceCategory
	CollectorID          int64   `json:"collector_id"`                        // Foreign key to Collector
	PricePerKg           float64 `json:"price_per_kg" binding:"required"`     // Cost per kg for this waste type
	MaximumCapacity      float64 `json:"maximum_capacity" binding:"required"` // Max capacity for this waste type
	HandlingRequirements string  `json:"handling_requirements,omitempty"`     // Special handling instructions
}

type Vehicle struct {
	VehicleID   int64   `json:"vehicle_id"`                      // Primary key
	VehicleType string  `json:"vehicle_type" binding:"required"` // Type of vehicle
	Capacity    float64 `json:"capacity" binding:"required"`     // Vehicle capacity in kg
}

type CollectorVehicle struct {
	VehicleID            int64  `json:"vehicle_id"`                        // Foreign key to Vehicle
	CollectorID          int64  `json:"collector_id"`                      // Foreign key to Collector
	VehicleNumber        string `json:"vehicle_number" binding:"required"` // Vehicle registration number
	MaintenanceDate      Date   `json:"maintenance_date"`                  // Last maintenance date
	IsActive             bool   `json:"is_active"`                         // Whether vehicle is in service
	GPSTrackingID        string `json:"gps_tracking_id"`                   // GPS tracking device ID
	AssignedDriverID     int64  `json:"assigned_driver_id"`                // Currently assigned driver
	RegistrationDocument string `json:"registration_document"`             // Path to registration document
	RegistrationExpiry   Date   `json:"registration_expiry"`               // Registration expiration date
}

type CollectorDriver struct {
	DriverID          int64   `json:"driver_id"`           // Primary key
	CollectorID       int64   `json:"collector_id"`        // Foreign key to Collector
	LicenseNumber     string  `json:"license_number"`      // Driver's license number
	LicenseExpiry     Date    `json:"license_expiry"`      // License expiration date
	AssignedVehicleID int64   `json:"assigned_vehicle_id"` // Currently assigned vehicle
	IsEmployed        bool    `json:"is_employed"`         // Whether driver is currently employed
	IsActive          bool    `json:"is_active"`           // Whether driver is available for trips
	Rating            float64 `json:"rating"`              // Driver's performance rating
	JoiningDate       Date    `json:"joining_date"`        // Date when driver joined
}

type CollectorDriverLocation struct {
	DriverID  int64    `json:"driver_id"`  // Foreign key to CollectorDriver
	Latitude  float64  `json:"latitude"`   // Current latitude
	Longitude float64  `json:"longitude"`  // Current longitude
	Timestamp DateTime `json:"timestamp"`  // When location was recorded
	IsActive  bool     `json:"is_active"`  // Whether driver is active
	TripID    int64    `json:"trip_id"`    // Current trip ID (if any)
	VehicleID int64    `json:"vehicle_id"` // Current vehicle ID
}

type PickupRequest struct {
	RequestID       int64    `json:"request_id"`
	BusinessID      int64    `json:"business_id"`
	CollectorID     int64    `json:"collector_id"`
	WasteType       string   `json:"waste_type"`
	Quantity        float64  `json:"quantity"`
	PickupDate      DateTime `json:"pickup_date"`
	Status          string   `json:"status"` // e.g., "pending", "assigned", "completed"
	AssignedDriver  int64    `json:"assigned_driver,omitempty"`
	AssignedVehicle int64    `json:"assigned_vehicle,omitempty"`
	CreatedAt       DateTime `json:"created_at"`
}
