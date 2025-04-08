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

func CreateUser(storage storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user types.User

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// Hashing the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}
		user.PasswordHash = string(hashedPassword)

		// Auto-setting registration date
		user.Registration = types.Date{Time: time.Now()}

		// Applying server-side default values
		user.IsActive = true
		user.IsVerified = false
		user.IsFlagged = false

		// Fields the user needs to pass at the frontend: Email, PasswordHash, FullName, PhoneNumber (Optional), Address (Optional), ProfileImage (Optional), Role (should be one of Business, Collector, Admin or Government)

		lastId, err := storage.CreateUser(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, response.GeneralError(err))
			return
		}

		slog.Info("user created successfully", slog.String("User ID", fmt.Sprint(lastId)))

		c.JSON(http.StatusCreated, gin.H{
			"status": "OK",
			"user":   lastId,
		})
	}
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

		c.JSON(http.StatusOK, gin.H{
			"status": "OK",
			"token":  signedToken,
			"user":   user,
		})
	}
}
