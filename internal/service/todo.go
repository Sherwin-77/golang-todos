package service

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sherwin-77/golang-todos/internal/entity"
	"github.com/sherwin-77/golang-todos/internal/http/dto"
	"github.com/sherwin-77/golang-todos/internal/repository"
	"github.com/sherwin-77/golang-todos/pkg/caches"
	"net/http"
	"time"
)

type TodoService interface {
	GetTodosByUserID(ctx context.Context, userID string) ([]entity.Todo, error)
	GetTodoByID(ctx context.Context, id string, userID string) (*entity.Todo, error)
	CreateTodo(ctx context.Context, request dto.TodoRequest, userID string) (*entity.Todo, error)
	UpdateTodo(ctx context.Context, request dto.UpdateTodoRequest, userID string) (*entity.Todo, error)
	DeleteTodo(ctx context.Context, id string, userID string) error
}

type todoService struct {
	todoRepository repository.TodoRepository
	userRepository repository.UserRepository
	cache          caches.Cache
}

func NewTodoService(todoRepository repository.TodoRepository, userRepository repository.UserRepository, cache caches.Cache) TodoService {
	return &todoService{todoRepository, userRepository, cache}
}

func (s *todoService) GetTodosByUserID(ctx context.Context, userID string) ([]entity.Todo, error) {
	todoKey := "todos:all:" + userID
	var todos []entity.Todo
	cachedData := s.cache.Get(todoKey)
	if cachedData != "" {
		if err := json.Unmarshal([]byte(cachedData), &todos); err != nil {
			return nil, err
		}
	} else {
		var err error
		db := s.todoRepository.SingleTransaction()
		todos, err = s.todoRepository.GetTodosByUserID(ctx, db, userID)
		if err != nil {
			return nil, err
		}

		data, _ := json.Marshal(todos)

		if err := s.cache.Set(todoKey, string(data), 5*time.Minute); err != nil {
			return nil, err
		}
	}

	return todos, nil
}

func (s *todoService) GetTodoByID(ctx context.Context, id string, userID string) (*entity.Todo, error) {
	todoKey := "todos:" + id
	todo := &entity.Todo{}
	cachedData := s.cache.Get(todoKey)
	if cachedData != "" {
		if err := json.Unmarshal([]byte(cachedData), todo); err != nil {
			return nil, err
		}
	} else {
		var err error
		db := s.todoRepository.SingleTransaction()
		todo, err = s.todoRepository.GetTodoByID(ctx, db, id)
		if err != nil {
			return nil, err
		}

		data, _ := json.Marshal(todo)

		if err := s.cache.Set(todoKey, string(data), 5*time.Minute); err != nil {
			return nil, err
		}
	}

	if todo.UserID.String() != userID {
		return nil, echo.NewHTTPError(http.StatusNotFound, http.StatusText(http.StatusNotFound))
	}

	return todo, nil
}

func (s *todoService) CreateTodo(ctx context.Context, request dto.TodoRequest, userID string) (*entity.Todo, error) {
	db := s.todoRepository.SingleTransaction()

	todo := &entity.Todo{
		Title:       request.Title,
		Description: request.Description,
		IsCompleted: request.IsCompleted,
		UserID:      uuid.MustParse(userID),
	}

	if err := s.todoRepository.CreateTodo(ctx, db, todo); err != nil {
		return nil, err
	}

	if err := s.cache.Del("todos:all:" + userID); err != nil {
		return nil, err
	}

	return todo, nil
}

func (s *todoService) UpdateTodo(ctx context.Context, request dto.UpdateTodoRequest, userID string) (*entity.Todo, error) {
	db := s.todoRepository.SingleTransaction()

	todo, err := s.todoRepository.GetTodoByID(ctx, db, request.ID)
	if err != nil {
		return nil, err
	}

	if todo.UserID.String() != userID {
		return nil, echo.NewHTTPError(http.StatusNotFound, "Todo not found")
	}

	todo.Title = request.Title
	todo.Description = request.Description
	todo.IsCompleted = request.IsCompleted

	if err := s.todoRepository.UpdateTodo(ctx, db, todo); err != nil {
		return nil, err
	}

	if err := s.cache.Del("todos:" + todo.ID.String()); err != nil {
		return nil, err
	}

	if err := s.cache.Del("todos:all:" + userID); err != nil {
		return nil, err
	}

	return todo, nil
}

func (s *todoService) DeleteTodo(ctx context.Context, id string, userID string) error {
	db := s.todoRepository.SingleTransaction()

	todo, err := s.todoRepository.GetTodoByID(ctx, db, id)
	if err != nil {
		return err
	}

	if todo.UserID.String() != userID {
		return echo.NewHTTPError(http.StatusNotFound, "Todo not found")
	}

	if err := s.todoRepository.DeleteTodo(ctx, db, todo); err != nil {
		return err
	}

	if err := s.cache.Del("todos:" + todo.ID.String()); err != nil {
		return err
	}

	if err := s.cache.Del("todos:all:" + userID); err != nil {
		return err
	}

	return nil
}
