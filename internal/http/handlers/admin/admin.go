package admin

import (
	"net/http"

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
