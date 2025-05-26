package driver

import (
	"net/http"
	"strconv"

	"cloud.google.com/go/pubsub"
	"github.com/gin-gonic/gin"
	"github.com/kartikey1188/build-in-progress_01/internal/pub_sub"
	"github.com/kartikey1188/build-in-progress_01/internal/storage"
	"github.com/kartikey1188/build-in-progress_01/internal/utils/response"
)

func StartDelivery(storage storage.Storage, pubsubClient *pubsub.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		pickupRequestID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pickup-request ID"})
			return
		}

		err1 := pub_sub.StartDelivery(storage, pubsubClient, pickupRequestID)
		if err1 != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err1))
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "OK", "Start Delivery for Request ID": pickupRequestID})
	}
}
func EndDelivery(storage storage.Storage, pubsubClient *pubsub.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		pickupRequestID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pickup-request ID"})
			return
		}

		err1 := pub_sub.EndDelivery(storage, pubsubClient, pickupRequestID)
		if err1 != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err1))
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "OK", "Completed Delivery for Request ID": pickupRequestID})
	}
}
