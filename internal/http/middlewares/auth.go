package middlewares

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/sherwin-77/go-echo-template/configs"
	"gorm.io/gorm"
)

type AuthMiddleware struct {
	config *configs.Config
	db     *gorm.DB
}

func NewAuthMiddleware(config *configs.Config, db *gorm.DB) *AuthMiddleware {
	return &AuthMiddleware{config, db}
}

func (m *AuthMiddleware) Authenticated(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		// Extract the "Authorization" header.
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		}

		// Split the token from the header.
		splitToken := strings.Split(authHeader, "Bearer ")
		if len(splitToken) != 2 {
			return echo.NewHTTPError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		}

		tokenString := strings.TrimSpace(splitToken[1])

		// Parse the JWT token.
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, echo.NewHTTPError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
			}
			return []byte(m.config.JWTSecret), nil
		})

		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			return echo.NewHTTPError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
		}

		c.Set("user_id", claims["id"])

		return next(c)
	}
}

func (m *AuthMiddleware) AuthLevel(level int) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if m.db == nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Database connection not available")
			}

			userIDData := c.Get("user_id")
			if userIDData == nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
			}

			userID := userIDData.(string)
			var userLevel int

			m.db.Table("role_users").
				Where("user_id = ?", userID).
				Joins("JOIN roles ON role_users.role_id = roles.id").
				Select("MAX(roles.auth_level)").
				Scan(&userLevel)

			if userLevel < level {
				return echo.NewHTTPError(http.StatusForbidden, "Insufficient permission")
			}
			return next(c)
		}
	}
}
