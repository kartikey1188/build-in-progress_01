package storage

import (
	"github.com/kartikey1188/build-in-progress_01/internal/types"
)

type Storage interface {
	LoginAndRegister
	Admin
	Collector
	General
	Business
	Driver
}

type LoginAndRegister interface {
	CreateCollectorUser(user types.Collector) (int64, error)
	CreateBusinessUser(user types.Business) (int64, error)
	CreateCollectorDriver(user types.CollectorDriver, collectorID int64) (int64, error)
	UpdateLastLogin(userID int64, lastLogin types.DateTime) error
	GetCollectorByEmail(email string) (types.Collector, error)
	GetBusinessByEmail(email string) (types.Business, error)
	GetCollectorDriverByEmail(email string) (types.CollectorDriver, error)
}

type Admin interface {
	FlagUser(userID string) error
	UnflagUser(userID string) error
	VerifyUser(userID string) error
	UnverifyUser(userID string) error

	AddVehicle(v types.Vehicle) (int64, error)
	AddServiceCategory(sc types.ServiceCategory) (int64, error)
	DeleteVehicle(vehicleID uint64) error
	DeleteServiceCategory(categoryID uint64) error
	GetAllCollectors() ([]types.Collector, error)
	GetAllBusinesses() ([]types.Business, error)
	GetAllUsers() ([]types.User, error)

	GetAllPickupRequests() ([]types.PickupRequest, error)
}

type General interface {
	GetAllServiceCategories() ([]types.ServiceCategory, error)
	GetServiceCategory(categoryID uint64) (types.ServiceCategory, error)
	GetAllVehicles() ([]types.Vehicle, error)
	GetVehicle(vehicleID uint64) (types.Vehicle, error)
	GetUserByID(userID uint64) (types.User, error)
	GetUserByEmail(email string) (types.User, error)
}

type Collector interface {
	GetCollectors() ([]types.Collector, error)
	GetCollectorByID(id int64) (types.Collector, error)

	UpdateCollectorProfile(userID int64, input types.CollectorUpdate) (int64, error)

	AddCollectorServiceCategory(input types.CollectorServiceCategory, userID uint64) (int64, error)
	UpdateCollectorServiceCategory(input types.UpdateCollectorServiceCategory, userID uint64) error
	DeleteCollectorServiceCategory(category_id int64, collectorID uint64) error

	AddCollectorVehicle(input types.CollectorVehicle, userID uint64) (int64, error)
	UpdateCollectorVehicle(input types.UpdateCollectorVehicle, userID uint64) error
	DeleteCollectorVehicle(vehicleID int64, collectorID uint64) error
	GetCollectorVehicle(collectorID int64, vehicleID int64) (types.CollectorVehicle, error)
	// Activating/Deactivating a vehicle can also be done through UpdateCollectorVehicle only

	UpdateCollectorDriver(input types.UpdateCollectorDriver, collectorID uint64) error
	AssignVehicleToDriver(driverID int64, vehicleID int64, collectorID uint64) error
	UnassignVehicleFromDriver(driverID int64, vehicleID int64, collectorID uint64) error
	GetCollectorDriver(collectorID int64, driverID int64) (types.CollectorDriver, error)
	DeleteCollectorDriver(driverID int64, collectorID uint64) error
	GetCollectorDrivers(collectorID int64) ([]types.CollectorDriver, error)

	GetCollectorServiceCategories(collectorID int64) ([]types.CollectorServiceCategory, error)
	GetCollectorVehicles(collectorID int64) ([]types.CollectorVehicle, error)

	// StoreDriverLocation(location models.DriverLocation) error

	AcceptPickupRequest(requestID int64) error
	RejectPickupRequest(requestID int64) error

	AssignTripToDriver(requestID int64, driverID int64) error
	UnassignTripFromDriver(requestID int64) error
}

type Business interface {
	GetBusinessByID(id int64) (types.Business, error)
	GetBusinessByEmail(email string) (types.Business, error)
	UpdateBusinessProfile(userID int64, input types.BusinessUpdate) (int64, error)

	GetPickupRequestByID(requestID int64) (types.PickupRequest, error)
	CreatePickupRequest(request types.PickupRequest) (int64, error)
	GetAllPickupRequestsForBusiness(businessID int64) ([]types.PickupRequest, error)
	UpdatePickupRequest(requestID int64, input types.UpdatePickupRequest) error
}

type Driver interface {
	EndDelivery(requestID int64) error
}
