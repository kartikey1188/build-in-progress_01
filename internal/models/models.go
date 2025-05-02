package models

import (
	"time"
)

type User struct {
	UserID       int64     `gorm:"primaryKey;autoIncrement;column:user_id"`
	Email        string    `gorm:"column:email;unique;not null"`
	PasswordHash string    `gorm:"column:password_hash;not null"`
	FullName     string    `gorm:"column:full_name;not null"`
	PhoneNumber  string    `gorm:"column:phone_number;size:20"`
	Address      string    `gorm:"column:address;type:text"`
	Registration time.Time `gorm:"column:registration_date;not null"`
	Role         string    `gorm:"column:role;not null;check:role IN ('Business','Collector','Admin','Government','Driver')"`
	IsActive     bool      `gorm:"column:is_active;not null;default:true"`
	ProfileImage string    `gorm:"column:profile_image;type:text"`
	LastLogin    time.Time `gorm:"column:last_login"`
	IsVerified   bool      `gorm:"column:is_verified;not null;default:false"`
	IsFlagged    bool      `gorm:"column:is_flagged;not null;default:false"`
}

type Business struct {
	UserID             int64  `gorm:"column:user_id;not null;uniqueIndex;foreignKey:user_id;references:User;onDelete:CASCADE"`
	BusinessName       string `gorm:"column:business_name;not null;size:255"`
	BusinessType       string `gorm:"column:business_type;not null;size:255"`
	RegistrationNumber string `gorm:"column:registration_number;not null;unique;size:100"`
	GstID              string `gorm:"column:gst_id;not null;unique;size:50"`
	BusinessAddress    string `gorm:"column:business_address;not null;type:text"`
}

type Collector struct {
	UserID        int64     `gorm:"column:user_id;not null;uniqueIndex;foreignKey:user_id;references:User;onDelete:CASCADE"`
	CompanyName   string    `gorm:"column:company_name;not null;size:255"`
	LicenseNumber string    `gorm:"column:license_number;not null;unique;size:100"`
	Capacity      int64     `gorm:"column:capacity;not null"`
	LicenseExpiry time.Time `gorm:"column:license_expiry;not null"`
}

type ServiceCategory struct {
	CategoryID int64  `gorm:"primaryKey;autoIncrement;column:category_id"`
	WasteType  string `gorm:"column:waste_type;not null;unique;size:100"`
}

type CollectorServiceCategory struct {
	CollectorID          int64   `gorm:"primaryKey;column:collector_id;not null;index;foreignKey:collector_id;references:Collector;onDelete:CASCADE"`
	CategoryID           int64   `gorm:"primaryKey;column:category_id;not null;index;foreignKey:category_id;references:ServiceCategory;onDelete:CASCADE"`
	PricePerKg           float64 `gorm:"column:price_per_kg;not null;type:decimal(10,2)"`
	MaximumCapacity      float64 `gorm:"column:maximum_capacity;not null;type:decimal(10,2)"`
	HandlingRequirements string  `gorm:"column:handling_requirements;type:text"`
}

type Vehicle struct {
	VehicleID   int64   `gorm:"primaryKey;autoIncrement;column:vehicle_id"`
	VehicleType string  `gorm:"column:vehicle_type;not null;size:50;uniqueIndex:idx_type_capacity"`
	Capacity    float64 `gorm:"column:capacity;not null;type:decimal(10,2);uniqueIndex:idx_type_capacity"`
}

type CollectorVehicle struct {
	CollectorID          int64     `gorm:"primaryKey;column:collector_id;not null;index;foreignKey:collector_id;references:Collector;onDelete:CASCADE"`
	VehicleID            int64     `gorm:"primaryKey;column:vehicle_id;not null;index;foreignKey:vehicle_id;references:Vehicle;onDelete:CASCADE"`
	VehicleNumber        string    `gorm:"column:vehicle_number;not null;unique;size:50"`
	MaintenanceDate      time.Time `gorm:"column:maintenance_date"`
	IsActive             bool      `gorm:"column:is_active;not null;default:true"`
	GPSTrackingID        string    `gorm:"column:gps_tracking_id;size:100"`
	RegistrationDocument string    `gorm:"column:registration_document;type:text"`
	RegistrationExpiry   time.Time `gorm:"column:registration_expiry"`
}

type CollectorDriver struct {
	DriverID      int64     `gorm:"primaryKey;autoIncrement;column:driver_id"`
	CollectorID   int64     `gorm:"primaryKey;column:collector_id;not null;index;foreignKey:collector_id;references:Collector;onDelete:CASCADE"`
	LicenseNumber string    `gorm:"column:license_number;not null;unique;size:100"`
	LicenseExpiry time.Time `gorm:"column:license_expiry;not null"`
	IsEmployed    bool      `gorm:"column:is_employed;not null;default:true"`
	IsActive      bool      `gorm:"column:is_active;not null;default:true"`
	Rating        float64   `gorm:"column:rating;type:decimal(3,2)"`
	JoiningDate   time.Time `gorm:"column:joining_date;not null"`
}

type VehicleDriver struct {
	DriverID  int64 `gorm:"column:driver_id;primaryKey;foreignKey:driver_id;references:CollectorDriver;onDelete:CASCADE"`
	VehicleID int64 `gorm:"column:vehicle_id;primaryKey;foreignKey:vehicle_id;references:CollectorVehicle;onDelete:CASCADE"`
}

type CollectorDriverLocation struct {
	ID        int64     `gorm:"primaryKey;autoIncrement;column:id"`
	DriverID  int64     `gorm:"column:driver_id;not null;index;foreignKey:driver_id;references:CollectorDriver;onDelete:CASCADE"`
	Latitude  float64   `gorm:"column:latitude;not null;type:decimal(9,6)"`
	Longitude float64   `gorm:"column:longitude;not null;type:decimal(9,6)"`
	Timestamp time.Time `gorm:"column:timestamp;not null"`
	IsActive  bool      `gorm:"column:is_active;not null;default:true"`
	TripID    int64     `gorm:"column:trip_id;index"`
	VehicleID int64     `gorm:"column:vehicle_id;index"`
}

type PickupRequest struct {
	RequestID       int64     `gorm:"primaryKey;autoIncrement;column:request_id"`
	BusinessID      int64     `gorm:"column:business_id;not null;index;foreignKey:business_id;references:Business;onDelete:CASCADE"`
	CollectorID     int64     `gorm:"column:collector_id;not null;index;foreignKey:collector_id;references:Collector;onDelete:CASCADE"`
	WasteType       string    `gorm:"column:waste_type;not null;size:100"`
	Quantity        float64   `gorm:"column:quantity;not null;type:decimal(10,2)"`
	PickupDate      time.Time `gorm:"column:pickup_date;not null"`
	Status          string    `gorm:"column:status;not null;size:50;check:status IN ('pending','assigned','completed')"`
	AssignedDriver  int64     `gorm:"column:assigned_driver;index"`
	AssignedVehicle int64     `gorm:"column:assigned_vehicle;index"`
	CreatedAt       time.Time `gorm:"column:created_at;not null;default:CURRENT_TIMESTAMP"`
}
