package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sherwin-77/go-echo-template/internal/http/handler"
	"github.com/sherwin-77/go-echo-template/internal/http/middlewares"
	"github.com/sherwin-77/go-echo-template/pkg/route"
)

func UserRoutes(userHandler handler.UserHandler, middlware middlewares.Middleware, authMiddleware middlewares.AuthMiddleware) ([]route.Route, []echo.MiddlewareFunc) {
	routes := []route.Route{
		{
			Method:      http.MethodPost,
			Path:        "/register",
			Handler:     userHandler.Register,
			Middlewares: []echo.MiddlewareFunc{},
		},
		{
			Method:      http.MethodPost,
			Path:        "/login",
			Handler:     userHandler.Login,
			Middlewares: []echo.MiddlewareFunc{},
		},
		{
			Method:  http.MethodPut,
			Path:    "/profile",
			Handler: userHandler.EditProfile,
			Middlewares: []echo.MiddlewareFunc{
				authMiddleware.Authenticated,
			},
		},
	}

	var middlewareFuncs []echo.MiddlewareFunc

	return routes, middlewareFuncs
}

func AdminUserRoutes(userHandler handler.UserHandler, middleware middlewares.Middleware, authMiddleware middlewares.AuthMiddleware) ([]route.Route, []echo.MiddlewareFunc) {
	routes := []route.Route{
		{
			Method:      http.MethodGet,
			Path:        "/users",
			Handler:     userHandler.GetUsers,
			Middlewares: []echo.MiddlewareFunc{},
		},
		{
			Method:      http.MethodPost,
			Path:        "/users",
			Handler:     userHandler.CreateUser,
			Middlewares: []echo.MiddlewareFunc{},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/users/:id",
			Handler: userHandler.UpdateUser,
			Middlewares: []echo.MiddlewareFunc{
				middleware.ValidateUUID([]string{"id"}),
			},
		},
		{
			Method:  http.MethodPatch,
			Path:    "/users/:id/role",
			Handler: userHandler.ChangeRole,
			Middlewares: []echo.MiddlewareFunc{
				middleware.ValidateUUID([]string{"id"}),
			},
		},
		{
			Method:  http.MethodGet,
			Path:    "/users/:id",
			Handler: userHandler.GetUserByID,
			Middlewares: []echo.MiddlewareFunc{
				middleware.ValidateUUID([]string{"id"}),
			},
		},
	}

	middlewareFuncs := []echo.MiddlewareFunc{
		authMiddleware.Authenticated,
		authMiddleware.AuthLevel(2),
	}

	return routes, middlewareFuncs

}
