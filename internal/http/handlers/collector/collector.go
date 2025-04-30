package collector

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kartikey1188/build-in-progress_01/internal/storage"
	"github.com/kartikey1188/build-in-progress_01/internal/types"
	"github.com/kartikey1188/build-in-progress_01/internal/utils/response"
)

// ListCollectors lists all collectors (unchanged)
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

// UpdateProfile updates a collector's profile (unchanged)
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

// OfferServiceCategory allows a collector to offer a new service category
func OfferServiceCategory(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		uidAny, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusUnauthorized, response.GeneralError(fmt.Errorf("user id missing")))
			return
		}
		userID := uidAny.(uint64)

		var input types.CollectorServiceCategory
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		id, err := storage.AddCollectorServiceCategory(input, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "OK", "Added Collector Service Category ID": id})
	}
}

// UpdateOfferedServiceCategory updates an existing service category offered by a collector
func UpdateOfferedServiceCategory(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		collectorID, exists := c.Get("collectorID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid service category ID"})
			return
		}

		var input types.CollectorServiceCategory
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		err = storage.UpdateCollectorServiceCategory(id, collectorID.(int64), input)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "OK", "Updated Service Category ID": id})
	}
}

// DeleteOfferedServiceCategory deletes a service category offered by a collector
func DeleteOfferedServiceCategory(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		uidAny, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusUnauthorized, response.GeneralError(fmt.Errorf("user id missing")))
			return
		}
		userID := uidAny.(uint64)

		type CategoryIdStruct struct {
			CategoryID int64 `json:"category_id"`
		}

		var req CategoryIdStruct
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		err := storage.DeleteCollectorServiceCategory(req.CategoryID, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":                      "OK",
			"Deleted Service Category ID": req.CategoryID,
		})
	}
}

// AppendVehicle allows a collector to append a new vehicle
func AppendVehicle(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		uidAny, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusUnauthorized, response.GeneralError(fmt.Errorf("user id missing")))
			return
		}
		userID := uidAny.(uint64)

		var input types.CollectorVehicle
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		id, err := storage.AddCollectorVehicle(input, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "OK", "Added Collector Vehicle ID": id})
	}
}

// UpdateVehicle updates an existing vehicle for a collector
func UpdateVehicle(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		collectorID, exists := c.Get("collectorID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid vehicle ID"})
			return
		}

		var input types.CollectorVehicle
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		err = storage.UpdateCollectorVehicle(id, collectorID.(int64), input)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "OK", "Updated Vehicle ID": id})
	}
}

// ActivateVehicle activates a collector's vehicle
func ActivateVehicle(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		collectorID, exists := c.Get("collectorID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid vehicle ID"})
			return
		}

		err = storage.ActivateCollectorVehicle(id, collectorID.(int64))
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "OK", "Activated Vehicle ID": id})
	}
}

// DeactivateVehicle deactivates a collector's vehicle
func DeactivateVehicle(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		collectorID, exists := c.Get("collectorID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid vehicle ID"})
			return
		}

		err = storage.DeactivateCollectorVehicle(id, collectorID.(int64))
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "OK", "Deactivated Vehicle ID": id})
	}
}

// RegisterDriver registers a new driver for a collector
func RegisterDriver(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		collectorID, exists := c.Get("collectorID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		var input types.CollectorDriver
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		input.CollectorID = collectorID.(int64)

		id, err := storage.AddCollectorDriver(input)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "OK", "Added Driver ID": id})
	}
}

// UpdateDriver updates an existing driver for a collector
func UpdateDriver(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		collectorID, exists := c.Get("collectorID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid driver ID"})
			return
		}

		var input types.CollectorDriver
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		err = storage.UpdateCollectorDriver(id, collectorID.(int64), input)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "OK", "Updated Driver ID": id})
	}
}

// AssignVehicleToDriver assigns a vehicle to a driver
func AssignVehicleToDriver(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		collectorID, exists := c.Get("collectorID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		driverID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid driver ID"})
			return
		}

		var input struct {
			VehicleID int64 `json:"vehicle_id"`
		}
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		err = storage.AssignVehicleToDriver(driverID, input.VehicleID, collectorID.(int64))
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "OK", "Assigned Vehicle ID": input.VehicleID, "to Driver ID": driverID})
	}
}

// GetCollectorDetails retrieves details of a specific collector
func GetCollectorDetails(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid collector ID"})
			return
		}

		collector, err := storage.GetCollectorByID(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, collector)
	}
}

// GetServiceCategories retrieves service categories offered by a collector
func GetServiceCategories(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid collector ID"})
			return
		}

		categories, err := storage.GetCollectorServiceCategories(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, categories)
	}
}

// GetCollectorVehicles retrieves vehicles appended by a collector
func GetCollectorVehicles(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid collector ID"})
			return
		}

		vehicles, err := storage.GetCollectorVehicles(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, vehicles)
	}
}

// GetCollectorDrivers retrieves drivers registered by a collector
func GetCollectorDrivers(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid collector ID"})
			return
		}

		drivers, err := storage.GetCollectorDrivers(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, drivers)
	}
}
