package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/kartikey1188/build-in-progress_01/internal/http/handlers/admincontrols"
	"github.com/kartikey1188/build-in-progress_01/internal/http/middleware"
	"github.com/kartikey1188/build-in-progress_01/internal/storage"
)

func Admin(router *gin.Engine, storage storage.Storage) {
	admin := router.Group("/admin")
	admin.Use(middleware.AdminOnly())

	admin.PUT("/verify/:id", admincontrols.VerifyUser(storage))
	admin.PUT("/unverify/:id", admincontrols.UnverifyUser(storage))
	admin.PUT("/flag/:id", admincontrols.FlagUser(storage))
	admin.PUT("/unflag/:id", admincontrols.UnflagUser(storage))
}
