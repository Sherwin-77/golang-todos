package builder

import (
	"github.com/labstack/echo/v4"
	"github.com/sherwin-77/go-echo-template/configs"
	"github.com/sherwin-77/go-echo-template/internal/http/handler"
	"github.com/sherwin-77/go-echo-template/internal/http/middlewares"
	"github.com/sherwin-77/go-echo-template/internal/http/router"
	"github.com/sherwin-77/go-echo-template/internal/repository"
	"github.com/sherwin-77/go-echo-template/internal/service"
	"github.com/sherwin-77/go-echo-template/pkg/caches"
	"github.com/sherwin-77/go-echo-template/pkg/tokens"
	"gorm.io/gorm"
)

func BuildV1Routes(config *configs.Config, db *gorm.DB, cache caches.Cache, group *echo.Group) {
	g := group.Group("/v1")

	// Initialize middlewares
	middleware := middlewares.NewMiddleware()
	authMiddleware := middlewares.NewAuthMiddleware(config, db)

	// Initialize repositories
	userRepository := repository.NewUserRepository(db)
	roleRepository := repository.NewRoleRepository(db)

	// Initialize services
	tokenService := tokens.NewTokenService(config.JWTSecret)
	userService := service.NewUserService(tokenService, userRepository, roleRepository, cache)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)

	// Register routes
	routes, middlewares := router.UserRoutes(*userHandler, *middleware, *authMiddleware)
	for _, route := range routes {
		m := append(middlewares, route.Middlewares...)
		g.Add(route.Method, route.Path, route.Handler, m...)
	}

	adminGroup := g.Group("/admin")

	adminUserRoutes, adminMiddlewares := router.AdminUserRoutes(*userHandler, *middleware, *authMiddleware)
	for _, route := range adminUserRoutes {
		m := append(adminMiddlewares, route.Middlewares...)
		adminGroup.Add(route.Method, route.Path, route.Handler, m...)
	}
}
