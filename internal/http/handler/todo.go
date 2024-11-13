package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/sherwin-77/golang-todos/internal/http/dto"
	"github.com/sherwin-77/golang-todos/internal/service"
	"github.com/sherwin-77/golang-todos/pkg/response"
	"net/http"
)

type TodoHandler struct {
	TodoService service.TodoService
}

func NewTodoHandler(todoService service.TodoService) *TodoHandler {
	return &TodoHandler{todoService}
}

func (h *TodoHandler) GetTodosByUserID(ctx echo.Context) error {
	userID := ctx.Get("user_id").(string)

	todos, err := h.TodoService.GetTodosByUserID(ctx.Request().Context(), userID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "Success", todos, nil))
}

func (h *TodoHandler) GetTodoByID(ctx echo.Context) error {
	userID := ctx.Get("user_id").(string)
	todoID := ctx.Param("id")
	if todoID == "" {
		return echo.NewHTTPError(http.StatusNotFound, http.StatusText(http.StatusNotFound))
	}

	todo, err := h.TodoService.GetTodoByID(ctx.Request().Context(), todoID, userID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "Success", todo, nil))
}

func (h *TodoHandler) CreateTodo(ctx echo.Context) error {
	userID := ctx.Get("user_id").(string)
	var req dto.TodoRequest

	if err := ctx.Bind(&req); err != nil {
		return err
	}

	if err := ctx.Validate(req); err != nil {
		return err
	}

	todo, err := h.TodoService.CreateTodo(ctx.Request().Context(), req, userID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusCreated, response.NewResponse(http.StatusCreated, "Todo created successfully", todo, nil))
}

func (h *TodoHandler) UpdateTodo(ctx echo.Context) error {
	userID := ctx.Get("user_id").(string)
	var req dto.UpdateTodoRequest

	if err := ctx.Bind(&req); err != nil {
		return err
	}

	if err := ctx.Validate(req); err != nil {
		return err
	}

	todo, err := h.TodoService.UpdateTodo(ctx.Request().Context(), req, userID)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "Todo updated successfully", todo, nil))
}

func (h *TodoHandler) DeleteTodo(ctx echo.Context) error {
	userID := ctx.Get("user_id").(string)
	todoID := ctx.Param("id")
	if todoID == "" {
		return echo.NewHTTPError(http.StatusNotFound, http.StatusText(http.StatusNotFound))
	}

	if err := h.TodoService.DeleteTodo(ctx.Request().Context(), todoID, userID); err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, response.NewResponse(http.StatusOK, "Todo deleted successfully", nil, nil))
}
