package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kartikey1188/build-in-progress_01/internal/http/handlers/general"
	"github.com/kartikey1188/build-in-progress_01/internal/storage"
)

func General(router *gin.Engine, storage storage.Storage) {
	general_routes := router.Group("/general")

	general_routes.GET("/all-categories", general.GetAllServiceCategories(storage))
	general_routes.GET("/category/:id", general.GetServiceCategory(storage))
	general_routes.GET("/all-vehicles", general.GetAllVehicles(storage))
	general_routes.GET("/vehicle/:id", general.GetVehicle(storage))
}
