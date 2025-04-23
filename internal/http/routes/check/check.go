package check

import (
	"github.com/gin-gonic/gin"
	"github.com/kartikey1188/build-in-progress_01/internal/http/handlers/collectorcontrols"
	"github.com/kartikey1188/build-in-progress_01/internal/storage"
)

func CollectorRoutes(router *gin.Engine, storage storage.Storage) {
	collector := router.Group("/collector")
	// collector.Use(middleware.CollectorOnly())

	collector.PUT("/profile/:id", collectorcontrols.UpdateProfile(storage))

	// Service Categories
	collector.POST("/service-categories", collectorcontrols.OfferServiceCategory(storage))
	collector.PUT("/service-categories/:id", collectorcontrols.UpdateOfferedServiceCategory(storage))
	collector.DELETE("/service-categories/:id", collectorcontrols.DeleteOfferedServiceCategory(storage))

	// Vehicles
	collector.POST("/vehicles", collectorcontrols.AppendVehicle(storage))
	collector.PUT("/vehicles/:id", collectorcontrols.UpdateVehicle(storage))
	collector.PUT("/vehicles/:id/activate", collectorcontrols.ActivateVehicle(storage))
	collector.PUT("/vehicles/:id/deactivate", collectorcontrols.DeactivateVehicle(storage))

	// Drivers
	collector.POST("/drivers", collectorcontrols.RegisterDriver(storage))
	collector.PUT("/drivers/:id", collectorcontrols.UpdateDriver(storage))
	collector.PUT("/drivers/:id/assign-vehicle", collectorcontrols.AssignVehicleToDriver(storage))

	// Open-access by ID
	router.GET("/collectors", collectorcontrols.ListCollectors(storage))
	router.GET("/collectors/:id", collectorcontrols.GetCollectorDetails(storage))
	router.GET("/collectors/:id/service-categories", collectorcontrols.GetServiceCategories(storage))
	router.GET("/collectors/:id/vehicles", collectorcontrols.GetCollectorVehicles(storage))
	router.GET("/collectors/:id/drivers", collectorcontrols.GetCollectorDrivers(storage))
}
