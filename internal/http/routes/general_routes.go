package routes

import (
	"cloud.google.com/go/pubsub"
	"github.com/gin-gonic/gin"
	"github.com/kartikey1188/build-in-progress_01/internal/http/handlers/general"
	"github.com/kartikey1188/build-in-progress_01/internal/storage"
)

func General(router *gin.Engine, storage storage.Storage, pubsubClient *pubsub.Client) {
	general_routes := router.Group("/general")

	general_routes.GET("/all-categories", general.GetAllServiceCategories(storage))
	general_routes.GET("/category/:id", general.GetServiceCategory(storage))
	general_routes.GET("/all-vehicles", general.GetAllVehicles(storage))
	general_routes.GET("/vehicle/:id", general.GetVehicle(storage))
	general_routes.GET("/user/:id", general.GetUserByID(storage))
	general_routes.GET("/user/email", general.GetUserByEmail(storage))
}
