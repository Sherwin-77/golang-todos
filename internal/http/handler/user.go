package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sherwin-77/go-echo-template/internal/http/dto"
	"github.com/sherwin-77/go-echo-template/internal/service"
	"github.com/sherwin-77/go-echo-template/pkg/response"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService}
}

/**
 * Admin Handlers
**/

func (h *UserHandler) GetUsers(ctx echo.Context) error {
	users, err := h.userService.GetUsers(ctx.Request().Context())

	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "Success", users, nil))
}

func (h *UserHandler) GetUserByID(ctx echo.Context) error {
	userID := ctx.Param("id")
	if userID == "" {
		return echo.NewHTTPError(http.StatusNotFound, http.StatusText(http.StatusNotFound))
	}

	user, err := h.userService.GetUserByID(ctx.Request().Context(), userID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "Success", user, nil))
}

func (h *UserHandler) CreateUser(ctx echo.Context) error {
	var req dto.UserRequest

	if err := ctx.Bind(&req); err != nil {
		return err
	}

	if err := ctx.Validate(req); err != nil {
		return err
	}

	user, err := h.userService.CreateUser(ctx.Request().Context(), req)

	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, response.NewResponse(http.StatusCreated, "User Created", user, nil))
}

func (h *UserHandler) UpdateUser(ctx echo.Context) error {
	var req dto.UpdateUserRequest

	if err := ctx.Bind(&req); err != nil {
		return err
	}

	if err := ctx.Validate(req); err != nil {
		return err
	}

	user, err := h.userService.UpdateUser(ctx.Request().Context(), req)

	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "User Updated", user, nil))
}

func (h *UserHandler) DeleteUser(ctx echo.Context) error {
	userID := ctx.Param("id")
	if userID == "" {
		return echo.NewHTTPError(http.StatusNotFound, http.StatusText(http.StatusNotFound))
	}

	if err := h.userService.DeleteUser(ctx.Request().Context(), userID); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "User Deleted", nil, nil))
}

func (h *UserHandler) ChangeRole(ctx echo.Context) error {
	var req dto.ChangeRoleRequest

	if err := ctx.Bind(&req); err != nil {
		return err
	}

	if err := ctx.Validate(req); err != nil {
		return err
	}

	if err := h.userService.ChangeRole(ctx.Request().Context(), req); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "Role Changed", nil, nil))
}

/**
 * User Handlers
**/

func (h *UserHandler) Register(ctx echo.Context) error {
	return h.CreateUser(ctx)
}

func (h *UserHandler) Login(ctx echo.Context) error {
	var req dto.LoginRequest

	if err := ctx.Bind(&req); err != nil {
		return err
	}

	if err := ctx.Validate(req); err != nil {
		return err
	}

	token, err := h.userService.Login(ctx.Request().Context(), req)

	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "Login Success", token, nil))
}

func (h *UserHandler) EditProfile(ctx echo.Context) error {
	userID := ctx.Get("user_id").(string)
	var req dto.UpdateUserRequest

	if err := ctx.Bind(&req); err != nil {
		return err
	}

	req.ID = userID

	if err := ctx.Validate(req); err != nil {
		return err
	}

	if userID != req.ID {
		return echo.NewHTTPError(http.StatusForbidden, http.StatusText(http.StatusForbidden))
	}

	return h.UpdateUser(ctx)
}
