package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kartikey1188/build-in-progress_01/internal/http/handlers/handleuser"
	"github.com/kartikey1188/build-in-progress_01/internal/http/handlers/home"
	"github.com/kartikey1188/build-in-progress_01/internal/storage"
)

func SetupRoutes(router *gin.Engine, storage storage.Storage) {
	router.GET("/", home.Home())
	router.POST("/auth/register/business", handleuser.CreateBusinessUser(storage))   // Changed
	router.POST("/auth/register/collector", handleuser.CreateCollectorUser(storage)) // Changed
	router.POST("/auth/login", handleuser.Login(storage))
	router.StaticFile("/docs/openapi.yaml", "./docs/openapi.yaml")
}
