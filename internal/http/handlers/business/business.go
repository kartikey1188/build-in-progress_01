package business

import (
	"net/http"
	"strconv"

	"cloud.google.com/go/pubsub"
	"github.com/gin-gonic/gin"
	"github.com/kartikey1188/build-in-progress_01/internal/pub_sub"
	"github.com/kartikey1188/build-in-progress_01/internal/storage"
	"github.com/kartikey1188/build-in-progress_01/internal/types"
	"github.com/kartikey1188/build-in-progress_01/internal/utils/response"
)

// GetBusinessByID retrieves a business by its ID.
func GetBusinessByID(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid business ID"})
			return
		}
		business, err := storage.GetBusinessByID(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		c.JSON(http.StatusOK, business)
	}
}

// GetBusinessByEmail retrieves a business by its email.
func GetBusinessByEmail(storage storage.Storage) gin.HandlerFunc {
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
		business, err := storage.GetBusinessByEmail(email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		c.JSON(http.StatusOK, business)
	}
}

// UpdateBusinessProfile updates a business's profile.
func UpdateBusinessProfile(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid business ID"})
			return
		}
		var input types.BusinessUpdate
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}
		updatedID, err := storage.UpdateBusinessProfile(userID, input)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "OK", "Updated Business ID": updatedID})
	}
}

func CreatePickupRequest(storage storage.Storage, pubsubClient *pubsub.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input types.PickupRequest
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		_, err := storage.GetBusinessByID(input.BusinessID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "business ID not found"})
			return
		}
		_, err = storage.GetCollectorByID(input.CollectorID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "collector ID not found"})
			return
		}

		id, err1 := pub_sub.CreatePickupRequest(storage, pubsubClient, input)
		if err1 != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err1))
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Pickup request created successfully", "pickup_request_id": id})
	}
}

func GetPickupRequestByID(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pickup request ID"})
			return
		}
		pickupRequest, err := storage.GetPickupRequestByID(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		c.JSON(http.StatusOK, pickupRequest)
	}
}

func GetAllPickupRequestsForBusiness(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		businessID, err := strconv.ParseInt(c.Param("business_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid business ID"})
			return
		}
		pickupRequests, err := storage.GetAllPickupRequestsForBusiness(businessID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		c.JSON(http.StatusOK, pickupRequests)
	}
}

func UpdatePickupRequest(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pickup request ID"})
			return
		}

		_, err = storage.GetPickupRequestByID(id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "pickup request ID not found"})
			return
		}

		var input types.UpdatePickupRequest
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}
		err = storage.UpdatePickupRequest(id, input)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Pickup request updated successfully"})
	}
}

// func CancelPickupRequest(storage storage.Storage) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pickup request ID"})
// 			return
// 		}

// 		_, err = storage.GetPickupRequestByID(id)
// 		if err != nil {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "pickup request ID not found"})
// 			return
// 		}

// 		err = kafka.CancelPickupRequest(id)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
// 			return
// 		}
// 		c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Pickup request cancelled successfully"})
// 	}
// }

// func UpdatePickupRequest(storage storage.Storage) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pickup request ID"})
// 			return
// 		}

// 		_, err = storage.GetPickupRequestByID(id)
// 		if err != nil {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "pickup request ID not found"})
// 			return
// 		}

// 		var input types.UpdatePickupRequest
// 		if err := c.ShouldBindJSON(&input); err != nil {
// 			c.JSON(http.StatusBadRequest, response.GeneralError(err))
// 			return
// 		}
// 		err = kafka.UpdatePickupRequest(id, input)
// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
// 			return
// 		}
// 		c.JSON(http.StatusOK, gin.H{"status": "OK", "message": "Pickup request updated successfully"})
// 	}
// }
