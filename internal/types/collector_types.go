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
	RegistrationDocument string `json:"registration_document"`             // Path to registration document
	RegistrationExpiry   Date   `json:"registration_expiry"`               // Registration expiration date
}

type CollectorDriver struct {
	User
	CollectorID   int64   `json:"collector_id" binding:"required"` // Foreign key to Collector
	LicenseNumber string  `json:"license_number"`                  // Driver's license number
	DriverName    string  `json:"driver_name"`                     // Driver's name
	LicenseExpiry Date    `json:"license_expiry"`                  // License expiration date
	IsEmployed    bool    `json:"is_employed"`                     // Whether driver is currently employed
	IsActive      bool    `json:"is_active"`                       // Whether driver is available for trips
	Rating        float64 `json:"rating"`                          // Driver's performance rating
	JoiningDate   Date    `json:"joining_date"`                    // Date when driver joined
}

type VehicleDriver struct {
	DriverID  int64 `json:"driver_id"`
	VehicleID int64 `json:"vehicle_id"`
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
	RequestID            int64    `json:"request_id"`
	BusinessID           int64    `json:"business_id" binding:"required"`
	CollectorID          int64    `json:"collector_id" binding:"required"`
	WasteType            string   `json:"waste_type"`
	Quantity             float64  `json:"quantity" binding:"required"`
	PickupDate           DateTime `json:"pickup_date"`
	Status               string   `json:"status" binding:"required,oneof=Pending Assigned Completed Cancelled"`
	HandlingRequirements string   `json:"handling_requirements"`
	AssignedDriver       int64    `json:"assigned_driver,omitempty"`
	AssignedVehicle      int64    `json:"assigned_vehicle,omitempty"`
	CreatedAt            DateTime `json:"created_at"`
}

type DriverLocation struct {
	LocationID  int64    `json:"location_id,omitempty"` // Omitempty for creation requests
	DriverID    int64    `json:"driver_id" binding:"required"`
	CollectorID int64    `json:"collector_id,omitempty"` // Often known from context
	VehicleID   int64    `json:"vehicle_id,omitempty"`   // Optional vehicle context
	Latitude    float64  `json:"latitude" binding:"required"`
	Longitude   float64  `json:"longitude" binding:"required"`
	Timestamp   DateTime `json:"timestamp"`          // Matches your existing DateTime type
	Accuracy    float32  `json:"accuracy,omitempty"` // Google Maps accuracy data
	Speed       float32  `json:"speed,omitempty"`    // Optional movement metrics
	Bearing     float32  `json:"bearing,omitempty"`  // Optional direction
	IsActive    bool     `json:"is_active"`          // Whether the driver is active
	TripID      int64    `json:"trip_id,omitempty"`  // Optional trip association
}

type Facility struct {
	FacilityID      int64    `json:"facility_id"`
	Name            string   `json:"name" binding:"required"`
	Location        string   `json:"location" binding:"required"`
	Status          string   `json:"status" binding:"required,oneof=operational near-capacity maintenance closed"`
	Capacity        float64  `json:"capacity" binding:"required,min=0,max=100"`
	ComplianceScore int      `json:"compliance_score" binding:"required,min=0,max=100"`
	DailyVolume     float64  `json:"daily_volume" binding:"required,min=0"`
	Permit          string   `json:"permit" binding:"required"`
	WasteTypes      []string `json:"waste_types" binding:"required"`
	Collectors      int      `json:"collectors"`
	CreatedAt       DateTime `json:"created_at"`
	UpdatedAt       DateTime `json:"updated_at"`
	IsActive        bool     `json:"is_active"`
}

type CollectorFacility struct {
	CollectorID          int64    `json:"collector_id"`
	FacilityID           int64    `json:"facility_id"`
	AssignmentDate       DateTime `json:"assignment_date"`
	LastProcessingDate   DateTime `json:"last_processing_date,omitempty"`
	ProcessingVolume     float64  `json:"processing_volume"`
	HandlingRequirements string   `json:"handling_requirements,omitempty"`
	IsActive             bool     `json:"is_active"`
}

type Zone struct {
	ZoneID          int64    `json:"zone_id"`
	Name            string   `json:"name" binding:"required"`
	Type            string   `json:"type" binding:"required,oneof=environmental cultural commercial security"`
	Status          string   `json:"status" binding:"required,oneof=permanent time-based"`
	Description     string   `json:"description,omitempty"`
	Area            string   `json:"area" binding:"required"` // JSON string containing boundaries
	ViolationsCount int      `json:"violations_count"`
	Authority       string   `json:"authority" binding:"required"` // JSON string containing authority info
	Restrictions    []string `json:"restrictions,omitempty"`
	CreatedAt       DateTime `json:"created_at"`
	UpdatedAt       DateTime `json:"updated_at"`
	IsActive        bool     `json:"is_active"`
}
