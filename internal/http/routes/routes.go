package routes

import (
	"cloud.google.com/go/pubsub"
	"github.com/gin-gonic/gin"
	"github.com/kartikey1188/build-in-progress_01/internal/storage"
)

func SetupRoutes(router *gin.Engine, storage storage.Storage, pubsubClient *pubsub.Client) {
	SetupAuth(router, storage)
	Admin(router, storage, pubsubClient)
	CollectorRoutes(router, storage, pubsubClient)
	General(router, storage, pubsubClient)
	BusinessRoutes(router, storage, pubsubClient)
	DriverRoutes(router, storage, pubsubClient)
}
