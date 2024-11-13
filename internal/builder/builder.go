package builder

import (
	"github.com/labstack/echo/v4"
	"github.com/sherwin-77/golang-todos/configs"
	"github.com/sherwin-77/golang-todos/internal/http/handler"
	"github.com/sherwin-77/golang-todos/internal/http/middlewares"
	"github.com/sherwin-77/golang-todos/internal/http/router"
	"github.com/sherwin-77/golang-todos/internal/repository"
	"github.com/sherwin-77/golang-todos/internal/service"
	"github.com/sherwin-77/golang-todos/pkg/caches"
	"github.com/sherwin-77/golang-todos/pkg/tokens"
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
	todoRepository := repository.NewTodoRepository(db)

	// Initialize services
	tokenService := tokens.NewTokenService(config.JWTSecret)
	userService := service.NewUserService(tokenService, userRepository, roleRepository, cache)
	todoService := service.NewTodoService(todoRepository, userRepository, cache)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)
	todoHandler := handler.NewTodoHandler(todoService)

	// Register routes
	userRoutes, userMiddlewares := router.UserRoutes(*userHandler, *middleware, *authMiddleware)
	for _, route := range userRoutes {
		m := append(userMiddlewares, route.Middlewares...)
		g.Add(route.Method, route.Path, route.Handler, m...)
	}

	todoRoutes, todoMiddlewares := router.TodoRoutes(*todoHandler, *middleware, *authMiddleware)
	for _, route := range todoRoutes {
		m := append(todoMiddlewares, route.Middlewares...)
		g.Add(route.Method, route.Path, route.Handler, m...)
	}

	adminGroup := g.Group("/admin")

	adminUserRoutes, adminMiddlewares := router.AdminUserRoutes(*userHandler, *middleware, *authMiddleware)
	for _, route := range adminUserRoutes {
		m := append(adminMiddlewares, route.Middlewares...)
		adminGroup.Add(route.Method, route.Path, route.Handler, m...)
	}
}
