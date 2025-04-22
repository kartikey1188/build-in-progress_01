package collectorcontrols

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kartikey1188/build-in-progress_01/internal/storage"
	"github.com/kartikey1188/build-in-progress_01/internal/types"
	"github.com/kartikey1188/build-in-progress_01/internal/utils/response"
)

func ListCollectors(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		collectors, err := storage.GetCollectors()
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, collectors)
	}
}

func UpdateProfile(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get ID from URL parameter
		userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
			return
		}

		var input types.CollectorUpdate
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// Pass the URL parameter ID to the storage layer
		id, err := storage.UpdateProfile(userID, input)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "OK", "Updated Collector ID": id})
	}
}

// Service Category Handlers
func AddServiceCategory(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input types.CollectorServiceCategory
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		id, err := storage.CreateServiceCategory(input)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "OK", "category_id": id})
	}
}

func UpdateServiceCategory(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var input types.CollectorServiceCategory
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		err := storage.UpdateServiceCategory(id, input)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "OK", "updated_category_id": id})
	}
}

func DeleteServiceCategory(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		err := storage.DeleteServiceCategory(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "OK", "deleted_category_id": id})
	}
}

// Vehicle Handlers
func AddVehicle(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input types.CollectorVehicle
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		id, err := storage.AddVehicle(input)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "OK", "vehicle_id": id})
	}
}

func UpdateVehicle(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var input types.CollectorVehicle
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		err := storage.UpdateVehicle(id, input)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "OK", "vehicle_id": id})
	}
}

func ActivateVehicle(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		err := storage.ActivateVehicle(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "OK", "vehicle_id": id})
	}
}

func DeactivateVehicle(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		err := storage.DeactivateVehicle(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "OK", "vehicle_id": id})
	}
}

// Driver Handlers
func RegisterDriver(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input types.CollectorDriver
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		id, err := storage.RegisterDriver(input)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "OK", "driver_id": id})
	}
}

func UpdateDriver(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var input types.CollectorDriver
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		err := storage.UpdateDriver(id, input)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "OK", "driver_id": id})
	}
}

func AssignVehicleToDriver(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		driverID := c.Param("id")
		var input struct {
			VehicleID int64 `json:"vehicle_id"`
		}
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		err := storage.AssignVehicleToDriver(driverID, input.VehicleID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "OK", "driver_id": driverID, "vehicle_id": input.VehicleID})
	}
}

// // Pickup Request Handlers
// func GetPickupRequests(storage storage.Storage) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		requests, err := storage.GetPickupRequests()
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
// 			return
// 		}

// 		c.JSON(http.StatusOK, requests)
// 	}
// }

// func GetPickupRequestDetail(storage storage.Storage) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		id := c.Param("id")
// 		request, err := storage.GetPickupRequestDetail(id)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
// 			return
// 		}

// 		c.JSON(http.StatusOK, request)
// 	}
// }

// func AssignPickup(storage storage.Storage) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		requestID := c.Param("id")
// 		var input struct {
// 			DriverID  int64 `json:"driver_id"`
// 			VehicleID int64 `json:"vehicle_id"`
// 		}
// 		if err := c.BindJSON(&input); err != nil {
// 			c.JSON(http.StatusBadRequest, response.GeneralError(err))
// 			return
// 		}

// 		err := storage.AssignPickup(requestID, input.DriverID, input.VehicleID)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
// 			return
// 		}

// 		c.JSON(http.StatusOK, gin.H{"status": "OK", "request_id": requestID})
// 	}
// }

// Collector Detail Handlers
func GetServiceCategories(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		categories, err := storage.GetCollectorServiceCategories(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, categories)
	}
}

func GetCollectorVehicles(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		vehicles, err := storage.GetCollectorVehicles(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, vehicles)
	}
}

func GetCollectorDrivers(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		drivers, err := storage.GetCollectorDrivers(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, drivers)
	}
}

func GetCollectorDetails(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		collector, err := storage.GetCollectorDetails(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, collector)
	}
}
