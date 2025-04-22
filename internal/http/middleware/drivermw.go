package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/kartikey1188/build-in-progress_01/internal/utils/response"
)

func DriverOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.GeneralError(fmt.Errorf("authorization header required")))
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.GeneralError(fmt.Errorf("invalid authorization header format")))
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			jwtSecret := os.Getenv("JWT_SECRET")
			if jwtSecret == "" {
				return nil, fmt.Errorf("jwt secret not configured")
			}

			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				response.GeneralError(fmt.Errorf("invalid or expired token")),
			)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(
				http.StatusUnauthorized,
				response.GeneralError(fmt.Errorf("invalid token claims")),
			)
			return
		}

		role, ok := claims["role"].(string)
		if !ok || role != "Driver" {
			c.AbortWithStatusJSON(
				http.StatusForbidden,
				response.GeneralError(fmt.Errorf("insufficient permissions")),
			)
			return
		}

		c.Next()
	}
}
