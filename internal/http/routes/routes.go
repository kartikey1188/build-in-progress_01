package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kartikey1188/build-in-progress_01/internal/storage"
)

func SetupRoutes(router *gin.Engine, storage storage.Storage) {
	SetupAuth(router, storage)
	Admin(router, storage)
	CollectorRoutes(router, storage)
}
