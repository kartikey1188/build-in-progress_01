package routes

import (
	"cloud.google.com/go/pubsub"
	"github.com/gin-gonic/gin"
	"github.com/kartikey1188/build-in-progress_01/internal/http/handlers/admin"
	"github.com/kartikey1188/build-in-progress_01/internal/http/handlers/business"
	"github.com/kartikey1188/build-in-progress_01/internal/http/handlers/collector"
	"github.com/kartikey1188/build-in-progress_01/internal/http/middleware"
	"github.com/kartikey1188/build-in-progress_01/internal/storage"
)

func Admin(router *gin.Engine, storage storage.Storage, pubsubClient *pubsub.Client) {
	admin_routes := router.Group("/admin")
	admin_routes.Use(middleware.AdminOnly())

	admin_routes.PUT("/verify/:id", admin.VerifyUser(storage))
	admin_routes.PUT("/unverify/:id", admin.UnverifyUser(storage))
	admin_routes.PUT("/flag/:id", admin.FlagUser(storage))
	admin_routes.PUT("/unflag/:id", admin.UnflagUser(storage))

	admin_routes.POST("/add/service-category", admin.AddServiceCategory(storage))
	admin_routes.POST("/add/vehicle", admin.AddVehicle(storage))

	admin_routes.DELETE("/delete/service-category/:id", admin.DeleteServiceCategory(storage))
	admin_routes.DELETE("/delete/vehicle/:id", admin.DeleteVehicle(storage))

	admin_routes.GET("/all/collectors", admin.GetAllCollectors(storage))
	admin_routes.GET("/all/businesses", admin.GetAllBusinesses(storage))
	admin_routes.GET("/all/users", admin.GetAllUsers(storage))
	admin_routes.GET("/collector/:id", collector.GetCollectorByID(storage))
	admin_routes.GET("/business/:id", business.GetBusinessByID(storage))

	admin_routes.GET("/all/pickup-requests", admin.GetAllPickupRequests(storage))

	// Facility management routes
	admin_routes.GET("/facilities", admin.GetFacilities(storage))
	admin_routes.GET("/facilities/:facility_id", admin.GetFacilityByID(storage))
	admin_routes.POST("/facilities", admin.CreateFacility(storage))
	admin_routes.PUT("/facilities/:facility_id", admin.UpdateFacility(storage))
	admin_routes.DELETE("/facilities/:facility_id", admin.DeleteFacility(storage))

	// Collector-Facility relationship management
	admin_routes.POST("/facilities/:facility_id/collectors", admin.AssignCollectorToFacility(storage))
	admin_routes.PUT("/facilities/:facility_id/collectors/:collector_id", admin.UpdateCollectorFacility(storage))
	admin_routes.DELETE("/facilities/:facility_id/collectors/:collector_id", admin.RemoveCollectorFromFacility(storage))
	admin_routes.GET("/facilities/:facility_id/collectors", admin.GetFacilityCollectors(storage))
	admin_routes.GET("/collectors/:collector_id/facilities", admin.GetCollectorFacilities(storage))

	// Zone management routes
	admin_routes.GET("/zones", admin.GetZones(storage))
	admin_routes.GET("/zones/:zone_id", admin.GetZoneByID(storage))
	admin_routes.POST("/zones", admin.CreateZone(storage))
	admin_routes.PUT("/zones/:zone_id", admin.UpdateZone(storage))
	admin_routes.DELETE("/zones/:zone_id", admin.DeleteZone(storage))
}
