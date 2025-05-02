package general

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kartikey1188/build-in-progress_01/internal/storage"
	"github.com/kartikey1188/build-in-progress_01/internal/utils/response"
)

func GetAllServiceCategories(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		categories, err := storage.GetAllServiceCategories()
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		c.JSON(http.StatusOK, categories)
	}
}

func GetServiceCategory(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
			return
		}

		category, err := storage.GetServiceCategory(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		c.JSON(http.StatusOK, category)
	}
}

func GetAllVehicles(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		vehicles, err := storage.GetAllVehicles()
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		c.JSON(http.StatusOK, vehicles)
	}
}

func GetVehicle(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid vehicle ID"})
			return
		}

		vehicle, err := storage.GetVehicle(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		c.JSON(http.StatusOK, vehicle)
	}
}
