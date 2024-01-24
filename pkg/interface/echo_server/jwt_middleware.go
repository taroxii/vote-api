package echoserver

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
	"github.com/taroxii/vote-api/pkg/config"
	"github.com/taroxii/vote-api/pkg/entity"
)

type GoMiddleware struct {
	// another stuff , may be needed by middleware
}

// InitMiddleware initialize the middleware
func InitMiddleware() *GoMiddleware {
	return &GoMiddleware{}
}

func (m *GoMiddleware) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Get the JWT token from the Authorization header
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing Authorization header")
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid Authorization header")
		}

		// Parse and validate the JWT token

		token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Config.JWTSecret), nil // Replace with your secret key
		})
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
		}
		claims, ok := token.Claims.(*jwt.MapClaims)

		if !ok || !token.Valid {
			return echo.NewHTTPError(http.StatusInternalServerError, "Invalid token")
		}

		data := entity.JWTClaims{}
		b, err := json.Marshal(claims)
		if err != nil {
			// echo.Logger.Info("Error to parse token encoding %s", err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		// Set the user ID in the context
		err = json.Unmarshal(b, &data)
		if err != nil {
			// echo.Logger.Info("Error to parse token encoding %s", err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		c.Set("user", &data)

		return next(c)
	}
}
