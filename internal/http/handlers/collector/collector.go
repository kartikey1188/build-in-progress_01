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
		id, err := storage.UpdateCollectorProfile(userID, input)
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
		uidStr := c.Param("id")

		userID, err := strconv.ParseUint(uidStr, 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid ID format"})
			return
		}

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
		uidStr := c.Param("id")

		userID, err := strconv.ParseUint(uidStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
			return
		}

		var input types.UpdateCollectorServiceCategory
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		if input.CategoryID == 0 {
			c.JSON(http.StatusBadRequest, response.GeneralError(fmt.Errorf("category_id is required")))
			return
		}

		err = storage.UpdateCollectorServiceCategory(input, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "OK", "Updated Service Category ID": userID})
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
		if req.CategoryID == 0 {
			c.JSON(http.StatusBadRequest, response.GeneralError(fmt.Errorf("category_id is required")))
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
func AppendCollectorVehicle(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		uidStr := c.Param("id")

		userID, err := strconv.ParseUint(uidStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
			return
		}

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
func UpdateCollectorVehicle(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
			return
		}

		var input types.UpdateCollectorVehicle
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		if input.VehicleID == 0 {
			c.JSON(http.StatusBadRequest, response.GeneralError(fmt.Errorf("vehicle_id is required")))
			return
		}

		err = storage.UpdateCollectorVehicle(input, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "OK", "Updated Vehicle ID": userID})
	}
}

func RemoveCollectorVehicle(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err1 := strconv.ParseUint(c.Param("id"), 10, 64)
		if err1 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
			return
		}

		type VehicleIdStruct struct {
			VehicleID int64 `json:"vehicle_id"`
		}

		var req VehicleIdStruct
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}
		if req.VehicleID == 0 {
			c.JSON(http.StatusBadRequest, response.GeneralError(fmt.Errorf("vehicle_id is required")))
			return
		}

		err := storage.DeleteCollectorVehicle(req.VehicleID, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":             "OK",
			"Deleted Vehicle ID": req.VehicleID,
		})
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
func GetCollectorServiceCategories(storage storage.Storage) gin.HandlerFunc {
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

func GetCollectorVehicle(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		collectorID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid collector ID"})
			return
		}
		vehicleID, err := strconv.ParseInt(c.Param("vid"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid vehicle ID"})
			return
		}
		vehicle, err := storage.GetCollectorVehicle(collectorID, vehicleID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, vehicle)
	}
}

// GetCollectorDrivers lists all drivers for a collector.
func GetCollectorDrivers(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		collectorID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid collector ID"})
			return
		}
		drivers, err := storage.GetCollectorDrivers(collectorID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, drivers)
	}
}

// GetCollectorDriver retrieves a single driver by driver ID.
func GetCollectorDriver(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		collectorID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid collector ID"})
			return
		}
		driverID, err := strconv.ParseInt(c.Param("did"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid driver ID"})
			return
		}
		driver, err := storage.GetCollectorDriver(collectorID, driverID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, driver)
	}
}

// AddCollectorDriver adds a new driver for a collector.
func AddCollectorDriver(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		collectorID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid collector ID"})
			return
		}
		var input types.CollectorDriver
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		id, err := storage.AddCollectorDriver(input, collectorID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "OK", "Added Collector Driver ID": id})
	}
}

// UpdateCollectorDriver updates an existing driver for a collector.
func UpdateCollectorDriver(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		collectorID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid collector ID"})
			return
		}
		var input types.UpdateCollectorDriver
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if input.DriverID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "driver_id is required"})
			return
		}
		err = storage.UpdateCollectorDriver(input, collectorID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "OK", "Updated Driver ID": input.DriverID})
	}
}

// DeleteCollectorDriver deletes a driver from a collector.
func DeleteCollectorDriver(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		collectorID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid collector ID"})
			return
		}
		type DriverIDStruct struct {
			DriverID int64 `json:"driver_id"`
		}
		var req DriverIDStruct
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if req.DriverID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "driver_id is required"})
			return
		}
		err = storage.DeleteCollectorDriver(req.DriverID, collectorID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "OK", "Deleted Driver ID": req.DriverID})
	}
}

// AssignVehicleToDriver assigns a vehicle to a driver.
func AssignVehicleToDriver(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		collectorID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid collector ID"})
			return
		}
		type Assignment struct {
			DriverID  int64 `json:"driver_id" binding:"required"`
			VehicleID int64 `json:"vehicle_id" binding:"required"`
		}
		var assign Assignment
		if err := c.ShouldBindJSON(&assign); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err = storage.AssignVehicleToDriver(assign.DriverID, assign.VehicleID, collectorID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "OK", "vehicle_id": assign.VehicleID, "driver_id": assign.DriverID, "message": "vehicle assigned successfully"})
	}
}

func UnassignVehicleFromDriver(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		collectorID, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid collector ID"})
			return
		}
		type Assignment struct {
			DriverID  int64 `json:"driver_id" binding:"required"`
			VehicleID int64 `json:"vehicle_id" binding:"required"`
		}
		var assign Assignment
		if err := c.ShouldBindJSON(&assign); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		err = storage.UnassignVehicleFromDriver(assign.DriverID, assign.VehicleID, collectorID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "OK", "vehicle_id": assign.VehicleID, "driver_id": assign.DriverID, "message": "vehicle unassigned successfully"})
	}
}

func GetCollectorByID(storage storage.Storage) gin.HandlerFunc {
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

func GetCollectorByEmail(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		type Email struct {
			Email string `json:"email" binding:"required,email"`
		}
		var emailInput Email
		if err := c.ShouldBindJSON(&emailInput); err != nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}
		email := emailInput.Email
		collector, err := storage.GetCollectorByEmail(email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		c.JSON(http.StatusOK, collector)
	}
}
