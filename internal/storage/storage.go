package storage

import "github.com/kartikey1188/build-in-progress_01/internal/types"

type Storage interface {
	LoginAndRegister
	AdminControls
	CollectorOperations
	ServiceCategoryOperations
	VehicleOperations
	DriverOperations
	// PickupRequestOperations
	CollectorDetailOperations
}

type LoginAndRegister interface {
	CreateCollectorUser(user types.Collector) (int64, error)
	CreateBusinessUser(user types.Business) (int64, error)
	GetUserByEmail(email string) (types.User, error)
	UpdateLastLogin(userID int64, lastLogin types.DateTime) error
	GetCollectorByEmail(email string) (types.Collector, error)
	GetBusinessByEmail(email string) (types.Business, error)
	GetUserById(userID int64) (types.User, error)
}

type AdminControls interface {
	FlagUser(userID string) error
	UnflagUser(userID string) error
	VerifyUser(userID string) error
	UnverifyUser(userID string) error
}

type CollectorOperations interface {
	GetCollector(id string) (types.Collector, error)
	GetCollectors() ([]types.Collector, error)
	UpdateProfile(userId int64, collector types.CollectorUpdate) (int64, error)
}

type ServiceCategoryOperations interface {
	CreateServiceCategory(csc types.CollectorServiceCategory) (int64, error)
	UpdateServiceCategory(id string, csc types.CollectorServiceCategory) error
	DeleteServiceCategory(id string) error
}

type VehicleOperations interface {
	AddVehicle(v types.CollectorVehicle) (int64, error)
	UpdateVehicle(id string, v types.CollectorVehicle) error
	ActivateVehicle(id string) error
	DeactivateVehicle(id string) error
}

type DriverOperations interface {
	RegisterDriver(d types.CollectorDriver) (int64, error)
	UpdateDriver(id string, d types.CollectorDriver) error
	AssignVehicleToDriver(driverID string, vehicleID int64) error
	UpdateDriverLocation(input types.CollectorDriverLocation) error
	GetDriverTrips(driverID int64) ([]types.PickupRequest, error)
}

// type PickupRequestOperations interface {
// 	GetPickupRequests() ([]types.PickupRequest, error)
// 	GetPickupRequestDetail(id string) (types.PickupRequest, error)
// 	AssignPickup(requestID string, driverID int64, vehicleID int64) error
// 	UpdateTripStatus(tripID string, status string) error
// }

type CollectorDetailOperations interface {
	GetCollectorServiceCategories(id string) ([]types.CollectorServiceCategory, error)
	GetCollectorVehicles(id string) ([]types.CollectorVehicle, error)
	GetCollectorDrivers(id string) ([]types.CollectorDriver, error)
	GetCollectorDetails(id string) (types.Collector, error)
}
