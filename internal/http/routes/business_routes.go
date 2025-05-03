package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kartikey1188/build-in-progress_01/internal/http/handlers/business"
	"github.com/kartikey1188/build-in-progress_01/internal/http/middleware"
	"github.com/kartikey1188/build-in-progress_01/internal/storage"
)

func BusinessRoutes(router *gin.Engine, storage storage.Storage) {
	business_routes := router.Group("/business")
	business_routes.Use(middleware.BusinessOnly())

	business_routes.GET("/:id", business.GetBusinessByID(storage))
	business_routes.GET("", business.GetBusinessByEmail(storage))
	business_routes.PATCH("/profile/:id", business.UpdateBusinessProfile(storage))

}
