package admin

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kartikey1188/build-in-progress_01/internal/storage"
	"github.com/kartikey1188/build-in-progress_01/internal/types"
	"github.com/kartikey1188/build-in-progress_01/internal/utils/response"
)

func VerifyUser(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if err := storage.VerifyUser(id); err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "OK", "Verified User ID": id})
	}
}

func UnverifyUser(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if err := storage.UnverifyUser(id); err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "OK", "Unverified User ID": id})
	}
}

func FlagUser(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if err := storage.FlagUser(id); err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "OK", "Flagged User ID": id})
	}
}

func UnflagUser(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		if err := storage.UnflagUser(id); err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "OK", "Unflagged User ID": id})
	}
}

func AddServiceCategory(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var serviceCategory types.ServiceCategory

		if err := c.ShouldBindJSON(&serviceCategory); err != nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		id, err := storage.AddServiceCategory(serviceCategory)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "OK", "Added Service Category ID": id})
	}
}

func AddVehicle(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var vehicle types.Vehicle

		if err := c.ShouldBindJSON(&vehicle); err != nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		id, err := storage.AddVehicle(vehicle)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "OK", "Added Vehicle ID": id})
	}
}

func DeleteServiceCategory(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid collector ID"})
			return
		}

		if err := storage.DeleteServiceCategory(id); err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "OK", "Deleted Service Category ID": id})
	}
}

func DeleteVehicle(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid collector ID"})
			return
		}

		if err := storage.DeleteVehicle(id); err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "OK", "Deleted Vehicle ID": id})
	}
}

func GetAllCollectors(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		collectors, err := storage.GetAllCollectors()
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		c.JSON(http.StatusOK, collectors)
	}
}

func GetAllBusinesses(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		businesses, err := storage.GetAllBusinesses()
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		c.JSON(http.StatusOK, businesses)
	}
}

func GetAllUsers(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := storage.GetAllUsers()
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		c.JSON(http.StatusOK, users)
	}
}

func GetAllPickupRequests(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		pickupRequests, err := storage.GetAllPickupRequests()
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		c.JSON(http.StatusOK, pickupRequests)
	}
}

// GetFacilities retrieves all facilities with optional filtering
func GetFacilities(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get optional query parameters
		status := c.Query("status")
		location := c.Query("location")
		wasteType := c.Query("type")

		facilities, err := storage.GetFacilities(status, location, wasteType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		c.JSON(http.StatusOK, facilities)
	}
}

// GetFacilityByID retrieves a specific facility by ID
func GetFacilityByID(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		facilityID, err := strconv.ParseInt(c.Param("facility_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid facility ID"})
			return
		}

		facility, err := storage.GetFacilityByID(facilityID)
		if err != nil {
			c.JSON(http.StatusNotFound, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, facility)
	}
}

// CreateFacility creates a new facility
func CreateFacility(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var facility types.Facility
		if err := c.ShouldBindJSON(&facility); err != nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// Set defaults for new facility
		facility.CreatedAt = types.DateTime{Time: time.Now()}
		facility.UpdatedAt = types.DateTime{Time: time.Now()}
		facility.IsActive = true

		facilityID, err := storage.CreateFacility(facility)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"status":      "OK",
			"facility_id": facilityID,
			"message":     "Facility created successfully",
		})
	}
}

// UpdateFacility updates an existing facility
func UpdateFacility(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		facilityID, err := strconv.ParseInt(c.Param("facility_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid facility ID"})
			return
		}

		var updateRequest types.UpdateFacilityRequest
		if err := c.ShouldBindJSON(&updateRequest); err != nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		if err := storage.UpdateFacility(facilityID, updateRequest); err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":      "OK",
			"facility_id": facilityID,
			"message":     "Facility updated successfully",
		})
	}
}

// DeleteFacility deletes (soft delete) a facility
func DeleteFacility(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		facilityID, err := strconv.ParseInt(c.Param("facility_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid facility ID"})
			return
		}

		if err := storage.DeleteFacility(facilityID); err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":      "OK",
			"facility_id": facilityID,
			"message":     "Facility deleted successfully",
		})
	}
}

// AssignCollectorToFacility assigns a collector to a facility
func AssignCollectorToFacility(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		facilityID, err := strconv.ParseInt(c.Param("facility_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid facility ID"})
			return
		}

		var collectorFacility types.CollectorFacility
		if err := c.ShouldBindJSON(&collectorFacility); err != nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// Set the facility ID from the URL parameter
		collectorFacility.FacilityID = facilityID
		collectorFacility.AssignmentDate = types.DateTime{Time: time.Now()}
		collectorFacility.IsActive = true

		if err := storage.AssignCollectorToFacility(facilityID, collectorFacility); err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":       "OK",
			"facility_id":  facilityID,
			"collector_id": collectorFacility.CollectorID,
			"message":      "Collector assigned to facility successfully",
		})
	}
}

// UpdateCollectorFacility updates a collector-facility relationship
func UpdateCollectorFacility(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		facilityID, err := strconv.ParseInt(c.Param("facility_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid facility ID"})
			return
		}

		collectorID, err := strconv.ParseInt(c.Param("collector_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid collector ID"})
			return
		}

		var request types.UpdateCollectorFacilityRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		if err := storage.UpdateCollectorFacility(facilityID, collectorID, request); err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":       "OK",
			"facility_id":  facilityID,
			"collector_id": collectorID,
			"message":      "Collector-facility relationship updated successfully",
		})
	}
}

// RemoveCollectorFromFacility removes a collector from a facility
func RemoveCollectorFromFacility(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		facilityID, err := strconv.ParseInt(c.Param("facility_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid facility ID"})
			return
		}

		collectorID, err := strconv.ParseInt(c.Param("collector_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid collector ID"})
			return
		}

		if err := storage.RemoveCollectorFromFacility(facilityID, collectorID); err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":       "OK",
			"facility_id":  facilityID,
			"collector_id": collectorID,
			"message":      "Collector removed from facility successfully",
		})
	}
}

// GetFacilityCollectors gets all collectors assigned to a facility
func GetFacilityCollectors(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		facilityID, err := strconv.ParseInt(c.Param("facility_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid facility ID"})
			return
		}

		collectors, err := storage.GetFacilityCollectors(facilityID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, collectors)
	}
}

// GetCollectorFacilities gets all facilities a collector is assigned to
func GetCollectorFacilities(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		collectorID, err := strconv.ParseInt(c.Param("collector_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid collector ID"})
			return
		}

		facilities, err := storage.GetCollectorFacilities(collectorID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, facilities)
	}
}

// GetZones gets all zones with optional filtering
func GetZones(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		zoneType := c.Query("type")
		status := c.Query("status")

		zones, err := storage.GetZones(zoneType, status)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, zones)
	}
}

// GetZoneByID gets a zone by its ID
func GetZoneByID(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		zoneID, err := strconv.ParseInt(c.Param("zone_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid zone ID"})
			return
		}

		zone, err := storage.GetZoneByID(zoneID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, zone)
	}
}

// CreateZone creates a new zone
func CreateZone(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var zone types.Zone
		if err := c.ShouldBindJSON(&zone); err != nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// Set defaults for new zone
		zone.CreatedAt = types.DateTime{Time: time.Now()}
		zone.UpdatedAt = types.DateTime{Time: time.Now()}
		zone.IsActive = true
		zone.ViolationsCount = 0

		zoneID, err := storage.CreateZone(zone)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"status":  "OK",
			"zone_id": zoneID,
			"message": "Zone created successfully",
		})
	}
}

// UpdateZone updates an existing zone
func UpdateZone(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		zoneID, err := strconv.ParseInt(c.Param("zone_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid zone ID"})
			return
		}

		var zoneUpdate types.UpdateZoneRequest
		if err := c.ShouldBindJSON(&zoneUpdate); err != nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		if err := storage.UpdateZone(zoneID, zoneUpdate); err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "OK",
			"zone_id": zoneID,
			"message": "Zone updated successfully",
		})
	}
}

// DeleteZone deletes a zone (soft delete)
func DeleteZone(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		zoneID, err := strconv.ParseInt(c.Param("zone_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid zone ID"})
			return
		}

		if err := storage.DeleteZone(zoneID); err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "OK",
			"zone_id": zoneID,
			"message": "Zone deleted successfully",
		})
	}
}
