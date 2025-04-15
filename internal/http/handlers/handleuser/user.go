package handleuser

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kartikey1188/build-in-progress_01/internal/storage"
	"github.com/kartikey1188/build-in-progress_01/internal/types"
	"github.com/kartikey1188/build-in-progress_01/internal/utils/response"
	"golang.org/x/crypto/bcrypt"
)

func CreateBusinessUser(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var business types.Business
		if err := c.ShouldBindJSON(&business); err != nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		if err := setUserDefaults(&business); err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		lastId, err := storage.CreateBusinessUser(business)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		slog.Info("Business user created successfully", slog.String("User ID", fmt.Sprint(lastId)))
		c.JSON(http.StatusCreated, gin.H{"status": "OK", "user": lastId})
	}
}

func CreateCollectorUser(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var collector types.Collector
		if err := c.ShouldBindJSON(&collector); err != nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		if err := setUserDefaults(&collector); err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		lastId, err := storage.CreateCollectorUser(collector)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		slog.Info("Collector user created successfully", slog.String("User ID", fmt.Sprint(lastId)))
		c.JSON(http.StatusCreated, gin.H{"status": "OK", "user": lastId})
	}
}

// Helper function to hash the password and set default values
func setUserDefaults(user interface{}) error {
	// Assuming user has a PasswordHash and Registration field
	var err error
	switch u := user.(type) {
	case *types.Business:
		u.PasswordHash, err = hashPassword(u.PasswordHash)
		if err != nil {
			return err
		}
		u.Registration = types.Date{Time: time.Now()}
		u.IsActive = true
		u.IsVerified = false
		u.IsFlagged = false
	case *types.Collector:
		u.PasswordHash, err = hashPassword(u.PasswordHash)
		if err != nil {
			return err
		}
		u.Registration = types.Date{Time: time.Now()}
		u.IsActive = true
		u.IsVerified = false
		u.IsFlagged = false
	default:
		return fmt.Errorf("unsupported user type")
	}
	return nil
}

// Helper function to hash the password
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func Login(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginData struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}

		if err := c.ShouldBindJSON(&loginData); err != nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		user, err := storage.GetUserByEmail(loginData.Email)
		if err != nil {
			c.JSON(http.StatusUnauthorized, response.GeneralError(fmt.Errorf("invalid email or password")))
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(loginData.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, response.GeneralError(fmt.Errorf("invalid email or password")))
			return
		}
		var business types.Business
		var collector types.Collector
		var admin types.User

		switch user.Role {

		case "Business":
			business, err = storage.GetBusinessByEmail(loginData.Email)
			if err != nil {
				c.JSON(http.StatusUnauthorized, err)
				return
			}

		case "Collector":
			collector, err = storage.GetCollectorByEmail(loginData.Email)
			if err != nil {
				c.JSON(http.StatusUnauthorized, err)
				return
			}

		case "Admin":
			admin, err = storage.GetUserByEmail(loginData.Email)
			if err != nil {
				c.JSON(http.StatusUnauthorized, err)
			}
		}

		// Updating last login timestamp
		user.LastLogin = types.DateTime{Time: time.Now()}
		storage.UpdateLastLogin(user.UserID, user.LastLogin)

		// Generating JWT token

		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			c.JSON(http.StatusInternalServerError, response.GeneralError(fmt.Errorf("JWT_SECRET not set")))
			return
		}

		claims := jwt.MapClaims{
			"user_id": user.UserID,
			"email":   user.Email,
			"role":    user.Role,
			"exp":     time.Now().Add(24 * time.Hour).Unix(),
			"iat":     time.Now().Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signedToken, err := token.SignedString([]byte(secret))
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(fmt.Errorf("failed to generate token")))
			return
		}

		switch user.Role {
		case "Business":
			c.JSON(http.StatusOK, gin.H{
				"status": "OK",
				"token":  signedToken,
				"user":   business,
			})
		case "Collector":
			c.JSON(http.StatusOK, gin.H{
				"status": "OK",
				"token":  signedToken,
				"user":   collector,
			})
		case "Admin":
			c.JSON(http.StatusOK, gin.H{
				"status": "OK",
				"token":  signedToken,
				"user":   admin,
			})
		}
	}
}
