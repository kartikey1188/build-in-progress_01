package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kartikey1188/build-in-progress_01/internal/http/handlers/collectorcontrols"
	"github.com/kartikey1188/build-in-progress_01/internal/http/middleware"
	"github.com/kartikey1188/build-in-progress_01/internal/storage"
)

func CollectorRoutes(router *gin.Engine, storage storage.Storage) {
	collector := router.Group("/collector")
	collector.Use(middleware.CollectorOnly())

	collector.PUT("/profile", collectorcontrols.UpdateProfile(storage))

	// Service Categories
	collector.POST("/service-categories", collectorcontrols.AddServiceCategory(storage))
	collector.PUT("/service-categories/:id", collectorcontrols.UpdateServiceCategory(storage))
	collector.DELETE("/service-categories/:id", collectorcontrols.DeleteServiceCategory(storage))

	// Vehicles
	collector.POST("/vehicles", collectorcontrols.AddVehicle(storage))
	collector.PUT("/vehicles/:id", collectorcontrols.UpdateVehicle(storage))
	collector.PUT("/vehicles/:id/activate", collectorcontrols.ActivateVehicle(storage))
	collector.PUT("/vehicles/:id/deactivate", collectorcontrols.DeactivateVehicle(storage))

	// Drivers
	collector.POST("/drivers", collectorcontrols.RegisterDriver(storage))
	collector.PUT("/drivers/:id", collectorcontrols.UpdateDriver(storage))
	collector.PUT("/drivers/:id/assign-vehicle", collectorcontrols.AssignVehicleToDriver(storage))

	// Pickup Requests
	// collector.GET("/pickup-requests", collectorcontrols.GetPickupRequests(storage))
	// collector.GET("/pickup-requests/:id", collectorcontrols.GetPickupRequestDetail(storage))
	// collector.PUT("/pickup-requests/:id/assign", collectorcontrols.AssignPickup(storage))

	// Open-access by ID
	router.GET("/collectors", collectorcontrols.ListCollectors(storage))
	router.GET("/collectors/:id", collectorcontrols.GetCollectorDetails(storage))
	router.GET("/collectors/:id/service-categories", collectorcontrols.GetServiceCategories(storage))
	router.GET("/collectors/:id/vehicles", collectorcontrols.GetCollectorVehicles(storage))
	router.GET("/collectors/:id/drivers", collectorcontrols.GetCollectorDrivers(storage))

	// Driver-only group
	DriverRoutes(router, storage)
}

func DriverRoutes(router *gin.Engine, storage storage.Storage) {
	driver := router.Group("/driver")
	driver.Use(middleware.DriverOnly())

	driver.GET("/trips", collectorcontrols.GetDriverTrips(storage))
	// driver.PUT("/trips/:id/status", collectorcontrols.UpdateTripStatus(storage))
	driver.POST("/location", collectorcontrols.UpdateDriverLocation(storage))
}
