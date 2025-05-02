package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kartikey1188/build-in-progress_01/internal/http/handlers/collector"
	"github.com/kartikey1188/build-in-progress_01/internal/http/middleware"
	"github.com/kartikey1188/build-in-progress_01/internal/storage"
)

func CollectorRoutes(router *gin.Engine, storage storage.Storage) {
	collector_routes := router.Group("/collector")
	collector_routes.Use(middleware.CollectorOnly())

	collector_routes.PATCH("/profile/:id", collector.UpdateProfile(storage))

	// Service Categories
	collector_routes.POST("/service-categories/:id", collector.OfferServiceCategory(storage))
	collector_routes.PATCH("/service-categories/:id", collector.UpdateOfferedServiceCategory(storage))
	collector_routes.DELETE("/service-categories/:id", collector.DeleteOfferedServiceCategory(storage))

	// Vehicles
	collector_routes.POST("/vehicles/:id", collector.AppendCollectorVehicle(storage))
	collector_routes.PATCH("/vehicles/:id", collector.UpdateCollectorVehicle(storage))
	collector_routes.DELETE("/vehicles/:id", collector.RemoveCollectorVehicle(storage))
	// --> Activating/Deactivating a vehicle can also be done through UpdateVehicle only

	// Drivers
	collector_routes.POST("/drivers/:id", collector.RegisterDriver(storage))
	collector_routes.PATCH("/drivers/:id", collector.UpdateDriver(storage))
	collector_routes.PUT("/drivers/:id/assign-vehicle", collector.AssignVehicleToDriver(storage))

	// Open-access by ID
	router.GET("/collector", collector.ListCollectors(storage))
	router.GET("/collector/:id", collector.GetCollectorDetails(storage))
	router.GET("/collector/:id/service-categories", collector.GetCollectorServiceCategories(storage))
	router.GET("/collector/:id/vehicles", collector.GetCollectorVehicles(storage))
	router.GET("/collector/:id/drivers", collector.GetCollectorDrivers(storage))
}
