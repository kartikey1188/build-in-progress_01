package routes

import (
	"cloud.google.com/go/pubsub"
	"github.com/gin-gonic/gin"
	"github.com/kartikey1188/build-in-progress_01/internal/http/handlers/driver"
	"github.com/kartikey1188/build-in-progress_01/internal/http/middleware"
	"github.com/kartikey1188/build-in-progress_01/internal/storage"
)

func DriverRoutes(router *gin.Engine, storage storage.Storage, pubsubClient *pubsub.Client) {
	driver_routes := router.Group("/driver")
	driver_routes.Use(middleware.DriverOnly())

	driver_routes.POST("/delivery/:id/start", driver.StartDelivery(storage, pubsubClient))
	driver_routes.POST("/delivery/:id/end", driver.EndDelivery(storage, pubsubClient))
}
