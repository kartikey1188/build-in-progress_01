package storage

import "github.com/kartikey1188/build-in-progress_01/internal/types"

type Storage interface {
	LoginAndRegister
	Admin
	Collectors
	General
}

type LoginAndRegister interface {
	CreateCollectorUser(user types.Collector) (int64, error)
	CreateBusinessUser(user types.Business) (int64, error)
	UpdateLastLogin(userID int64, lastLogin types.DateTime) error
	GetCollectorByEmail(email string) (types.Collector, error)
	GetBusinessByEmail(email string) (types.Business, error)
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
}

type General interface {
	GetAllServiceCategories() ([]types.ServiceCategory, error)
	GetServiceCategory(categoryID uint64) (types.ServiceCategory, error)
	GetAllVehicles() ([]types.Vehicle, error)
	GetVehicle(vehicleID uint64) (types.Vehicle, error)
	GetUserByID(userID uint64) (types.User, error)
	GetUserByEmail(email string) (types.User, error)
}

type Collectors interface {
	GetCollectors() ([]types.Collector, error)
	GetCollectorByID(id int64) (types.Collector, error)

	UpdateProfile(userID int64, input types.CollectorUpdate) (int64, error)

	AddCollectorServiceCategory(input types.CollectorServiceCategory, userID uint64) (int64, error)
	UpdateCollectorServiceCategory(input types.UpdateCollectorServiceCategory, userID uint64) error
	DeleteCollectorServiceCategory(category_id int64, collectorID uint64) error

	AddCollectorVehicle(input types.CollectorVehicle, userID uint64) (int64, error)
	UpdateCollectorVehicle(input types.UpdateCollectorVehicle, userID uint64) error
	DeleteCollectorVehicle(vehicleID int64, collectorID uint64) error
	// Activating/Deactivating a vehicle can also be done through UpdateCollectorVehicle only

	AddCollectorDriver(input types.CollectorDriver) (int64, error)
	UpdateCollectorDriver(id int64, collectorID int64, input types.CollectorDriver) error
	AssignVehicleToDriver(driverID int64, vehicleID int64, collectorID int64) error

	GetCollectorServiceCategories(collectorID int64) ([]types.CollectorServiceCategory, error)
	GetCollectorVehicles(collectorID int64) ([]types.CollectorVehicle, error)
	GetCollectorDrivers(collectorID int64) ([]types.CollectorDriver, error)
}
