package business

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
		email := c.Param("email")
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
