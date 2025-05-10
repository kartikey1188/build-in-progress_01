package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kartikey1188/build-in-progress_01/internal/http/handlers/admin"
	"github.com/kartikey1188/build-in-progress_01/internal/http/handlers/business"
	"github.com/kartikey1188/build-in-progress_01/internal/http/handlers/collector"
	"github.com/kartikey1188/build-in-progress_01/internal/http/middleware"
	"github.com/kartikey1188/build-in-progress_01/internal/storage"
)

func Admin(router *gin.Engine, storage storage.Storage) {
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
}
