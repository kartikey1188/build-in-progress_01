package storage

import "github.com/kartikey1188/build-in-progress_01/internal/types"

type Storage interface {
	LoginAndRegister
	AdminControls
	CollectorControls
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
	AddVehicle(v types.Vehicle) (int64, error)
	AddServiceCategory(sc types.ServiceCategory) (int64, error)
}

type CollectorControls interface {
	GetCollectors() ([]types.Collector, error)
	GetCollectorByID(id int64) (types.Collector, error)

	UpdateProfile(userID int64, input types.CollectorUpdate) (int64, error)

	AddCollectorServiceCategory(input types.CollectorServiceCategory) (int64, error)
	UpdateCollectorServiceCategory(id int64, collectorID int64, input types.CollectorServiceCategory) error
	DeleteCollectorServiceCategory(id int64, collectorID int64) error

	AddCollectorVehicle(input types.CollectorVehicle) (int64, error)
	UpdateCollectorVehicle(id int64, collectorID int64, input types.CollectorVehicle) error
	ActivateCollectorVehicle(id int64, collectorID int64) error
	DeactivateCollectorVehicle(id int64, collectorID int64) error

	AddCollectorDriver(input types.CollectorDriver) (int64, error)
	UpdateCollectorDriver(id int64, collectorID int64, input types.CollectorDriver) error
	AssignVehicleToDriver(driverID int64, vehicleID int64, collectorID int64) error

	GetCollectorServiceCategories(collectorID int64) ([]types.CollectorServiceCategory, error)
	GetCollectorVehicles(collectorID int64) ([]types.CollectorVehicle, error)
	GetCollectorDrivers(collectorID int64) ([]types.CollectorDriver, error)
}
