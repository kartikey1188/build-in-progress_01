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
	collector_routes.POST("/:id/service-categories", collector.OfferServiceCategory(storage))
	collector_routes.PATCH("/:id/service-categories", collector.UpdateOfferedServiceCategory(storage))
	collector_routes.DELETE("/:id/service-categories", collector.DeleteOfferedServiceCategory(storage))

	// Vehicles
	collector_routes.POST("/:id/vehicles", collector.AppendCollectorVehicle(storage))
	collector_routes.PATCH("/:id/vehicles", collector.UpdateCollectorVehicle(storage))
	collector_routes.DELETE("/:id/vehicles", collector.RemoveCollectorVehicle(storage))
	collector_routes.GET("/:id/vehicles/:vid", collector.GetCollectorVehicle(storage))
	// --> Activating/Deactivating a vehicle can also be done through UpdateVehicle only

	// Drivers
	collector_routes.GET("/:id/drivers", collector.GetCollectorDrivers(storage))
	collector_routes.GET("/:id/drivers/:did", collector.GetCollectorDriver(storage))
	collector_routes.POST("/:id/drivers", collector.AddCollectorDriver(storage))
	collector_routes.PATCH("/:id/drivers", collector.UpdateCollectorDriver(storage))
	collector_routes.DELETE("/:id/drivers", collector.DeleteCollectorDriver(storage))
	collector_routes.PUT("/:id/drivers/assign-vehicle", collector.AssignVehicleToDriver(storage))
	collector_routes.DELETE("/:id/drivers/unassign-vehicle", collector.UnassignVehicleFromDriver(storage))

	// Open-access by ID
	router.GET("/collector", collector.ListCollectors(storage))
	router.GET("/collector/:id", collector.GetCollectorDetails(storage))
	router.GET("/collector/:id/service-categories", collector.GetCollectorServiceCategories(storage))
	router.GET("/collector/:id/vehicles", collector.GetCollectorVehicles(storage))
}
