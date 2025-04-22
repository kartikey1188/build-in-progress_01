package collectorcontrols

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kartikey1188/build-in-progress_01/internal/storage"
	"github.com/kartikey1188/build-in-progress_01/internal/types"
	"github.com/kartikey1188/build-in-progress_01/internal/utils/response"
)

func UpdateDriverLocation(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input types.CollectorDriverLocation
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		err := storage.UpdateDriverLocation(input)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "OK"})
	}
}

func GetDriverTrips(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		driverID, _ := c.Get("driver_id")
		trips, err := storage.GetDriverTrips(driverID.(int64))
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, trips)
	}
}

// func UpdateTripStatus(storage storage.Storage) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		tripID := c.Param("id")
// 		var input struct {
// 			Status string `json:"status"`
// 		}
// 		if err := c.BindJSON(&input); err != nil {
// 			c.JSON(http.StatusBadRequest, response.GeneralError(err))
// 			return
// 		}

// 		err := storage.UpdateTripStatus(tripID, input.Status)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
// 			return
// 		}

// 		c.JSON(http.StatusOK, gin.H{"status": "OK", "trip_id": tripID})
// 	}
// }
