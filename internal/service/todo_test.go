package service_test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sherwin-77/golang-todos/internal/entity"
	"github.com/sherwin-77/golang-todos/internal/http/dto"
	"github.com/sherwin-77/golang-todos/internal/service"
	mock_caches "github.com/sherwin-77/golang-todos/test/mock/pkg/caches"
	mock_repository "github.com/sherwin-77/golang-todos/test/mock/repository"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"testing"
)

type TodoTestSuite struct {
	suite.Suite
	ctrl        *gomock.Controller
	repo        *mock_repository.MockTodoRepository
	userRepo    *mock_repository.MockUserRepository
	cache       *mock_caches.MockCache
	todoService service.TodoService
}

func (s *TodoTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.repo = mock_repository.NewMockTodoRepository(s.ctrl)
	s.userRepo = mock_repository.NewMockUserRepository(s.ctrl)
	s.cache = mock_caches.NewMockCache(s.ctrl)
	s.todoService = service.NewTodoService(s.repo, s.userRepo, s.cache)
}

func TestTodoService(t *testing.T) {
	suite.Run(t, new(TodoTestSuite))
}

func (s *TodoTestSuite) TestGetTodosByUserID() {
	userID := uuid.New().String()
	keyFindAll := "todos:all:" + userID
	todos := make([]entity.Todo, 0)
	marshalledData, _ := json.Marshal(todos)

	s.Run("Failed unmarshal", func() {
		s.cache.EXPECT().Get(keyFindAll).Return("invalid")
		result, err := s.todoService.GetTodosByUserID(context.Background(), userID)

		s.Error(err)
		s.Nil(result)
	})

	s.Run("Failed to get todos", func() {
		errorTest := errors.New("get todos error")
		s.cache.EXPECT().Get(keyFindAll).Return("")
		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetTodosByUserID(gomock.Any(), gomock.Any(), userID).Return(nil, errorTest)
		result, err := s.todoService.GetTodosByUserID(context.Background(), userID)

		s.ErrorIs(err, errorTest)
		s.Nil(result)
	})

	s.Run("Failed to set cache", func() {
		errorTest := errors.New("set cache error")
		s.cache.EXPECT().Get(keyFindAll).Return("")
		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetTodosByUserID(gomock.Any(), gomock.Any(), userID).Return(todos, nil)
		s.cache.EXPECT().Set(keyFindAll, string(marshalledData), gomock.Any()).Return(errorTest)
		result, err := s.todoService.GetTodosByUserID(context.Background(), userID)

		s.ErrorIs(err, errorTest)
		s.Nil(result)
	})

	s.Run("Successfully get todos", func() {
		s.cache.EXPECT().Get(keyFindAll).Return("")
		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetTodosByUserID(gomock.Any(), gomock.Any(), userID).Return(todos, nil)
		s.cache.EXPECT().Set(keyFindAll, string(marshalledData), gomock.Any()).Return(nil)
		result, err := s.todoService.GetTodosByUserID(context.Background(), userID)

		s.Nil(err)
		s.Equal(todos, result)
	})

	s.Run("Successfully get todos from cache", func() {
		s.cache.EXPECT().Get(keyFindAll).Return(string(marshalledData))
		result, err := s.todoService.GetTodosByUserID(context.Background(), userID)

		s.Nil(err)
		s.Equal(todos, result)
	})
}

func (s *TodoTestSuite) TestGetTodoByID() {
	todoID := uuid.New().String()
	userID := uuid.New().String()
	keyFindTodo := "todos:" + todoID
	todo := &entity.Todo{}
	todo.ID = uuid.MustParse(todoID)
	todo.UserID = uuid.MustParse(userID)
	marshalledData, _ := json.Marshal(todo)

	s.Run("Failed unmarshal", func() {
		s.cache.EXPECT().Get(keyFindTodo).Return("invalid")
		result, err := s.todoService.GetTodoByID(context.Background(), todoID, userID)

		s.Error(err)
		s.Nil(result)
	})

	s.Run("Failed to get todo", func() {
		errorTest := errors.New("get todo error")
		s.cache.EXPECT().Get(keyFindTodo).Return("")
		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetTodoByID(gomock.Any(), gomock.Any(), todoID).Return(nil, errorTest)
		result, err := s.todoService.GetTodoByID(context.Background(), todoID, userID)

		s.ErrorIs(err, errorTest)
		s.Nil(result)
	})

	s.Run("Failed to set cache", func() {
		errorTest := errors.New("set cache error")
		s.cache.EXPECT().Get(keyFindTodo).Return("")
		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetTodoByID(gomock.Any(), gomock.Any(), todoID).Return(todo, nil)
		s.cache.EXPECT().Set(keyFindTodo, string(marshalledData), gomock.Any()).Return(errorTest)
		result, err := s.todoService.GetTodoByID(context.Background(), todoID, userID)

		s.ErrorIs(err, errorTest)
		s.Nil(result)
	})

	s.Run("User ID mismatch", func() {
		var e *echo.HTTPError
		s.cache.EXPECT().Get(keyFindTodo).Return("")
		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetTodoByID(gomock.Any(), gomock.Any(), todoID).Return(todo, nil)
		s.cache.EXPECT().Set(keyFindTodo, string(marshalledData), gomock.Any()).Return(nil)
		result, err := s.todoService.GetTodoByID(context.Background(), todoID, uuid.NewString())

		s.ErrorAs(err, &e)
		s.Nil(result)
	})

	s.Run("Successfully get todo", func() {
		s.cache.EXPECT().Get(keyFindTodo).Return("")
		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetTodoByID(gomock.Any(), gomock.Any(), todoID).Return(todo, nil)
		s.cache.EXPECT().Set(keyFindTodo, string(marshalledData), gomock.Any()).Return(nil)
		result, err := s.todoService.GetTodoByID(context.Background(), todoID, userID)

		s.Nil(err)
		s.Equal(todo, result)
	})

	s.Run("Successfully get todo from cache", func() {
		s.cache.EXPECT().Get(keyFindTodo).Return(string(marshalledData))
		result, err := s.todoService.GetTodoByID(context.Background(), todoID, userID)

		s.Nil(err)
		s.Equal(todo, result)
	})
}

func (s *TodoTestSuite) TestCreateTodo() {
	userID := uuid.New().String()
	keyFindAll := "todos:all:" + userID
	s.Run("Failed to create todo", func() {
		errorTest := errors.New("create todo error")
		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().CreateTodo(gomock.Any(), gomock.Any(), gomock.Any()).Return(errorTest)
		result, err := s.todoService.CreateTodo(context.Background(), dto.TodoRequest{}, userID)

		s.ErrorIs(err, errorTest)
		s.Nil(result)
	})

	s.Run("Failed to delete cache", func() {
		errorTest := errors.New("delete cache error")
		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().CreateTodo(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		s.cache.EXPECT().Del(keyFindAll).Return(errorTest)
		result, err := s.todoService.CreateTodo(context.Background(), dto.TodoRequest{}, userID)

		s.ErrorIs(err, errorTest)
		s.Nil(result)
	})

	s.Run("Successfully create todo", func() {
		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().CreateTodo(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		s.cache.EXPECT().Del(keyFindAll).Return(nil)
		result, err := s.todoService.CreateTodo(context.Background(), dto.TodoRequest{}, userID)

		s.Nil(err)
		s.NotNil(result)
	})
}

func (s *TodoTestSuite) TestUpdateTodo() {
	userID := uuid.New().String()
	todoID := uuid.New().String()
	keyFindAll := "todos:all:" + userID
	keyFindTodo := "todos:" + todoID
	emptyTodo := &entity.Todo{}
	emptyTodo.ID = uuid.MustParse(todoID)
	emptyTodo.UserID = uuid.MustParse(userID)

	s.Run("Failed to get todo", func() {
		errorTest := errors.New("get todo error")
		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetTodoByID(gomock.Any(), gomock.Any(), todoID).Return(nil, errorTest)
		result, err := s.todoService.UpdateTodo(context.Background(), dto.UpdateTodoRequest{
			ID: todoID,
		}, userID)

		s.ErrorIs(err, errorTest)
		s.Nil(result)
	})

	s.Run("User ID mismatch", func() {
		var e *echo.HTTPError
		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetTodoByID(gomock.Any(), gomock.Any(), todoID).Return(emptyTodo, nil)
		result, err := s.todoService.UpdateTodo(context.Background(), dto.UpdateTodoRequest{
			ID: todoID,
		}, uuid.NewString())

		s.ErrorAs(err, &e)
		s.Nil(result)
	})

	s.Run("Failed to update todo", func() {
		errorTest := errors.New("update todo error")
		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetTodoByID(gomock.Any(), gomock.Any(), todoID).Return(emptyTodo, nil)
		s.repo.EXPECT().UpdateTodo(gomock.Any(), gomock.Any(), gomock.Any()).Return(errorTest)
		result, err := s.todoService.UpdateTodo(context.Background(), dto.UpdateTodoRequest{
			ID: todoID,
		}, userID)

		s.ErrorIs(err, errorTest)
		s.Nil(result)
	})

	s.Run("Failed to delete todo cache", func() {
		errorTest := errors.New("delete todo cache error")
		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetTodoByID(gomock.Any(), gomock.Any(), todoID).Return(emptyTodo, nil)
		s.repo.EXPECT().UpdateTodo(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		s.cache.EXPECT().Del("todos:" + todoID).Return(errorTest)
		result, err := s.todoService.UpdateTodo(context.Background(), dto.UpdateTodoRequest{
			ID: todoID,
		}, userID)

		s.ErrorIs(err, errorTest)
		s.Nil(result)
	})

	s.Run("Failed to delete todos cache", func() {
		errorTest := errors.New("delete todos cache error")
		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetTodoByID(gomock.Any(), gomock.Any(), todoID).Return(emptyTodo, nil)
		s.repo.EXPECT().UpdateTodo(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		s.cache.EXPECT().Del(keyFindTodo).Return(errorTest)
		result, err := s.todoService.UpdateTodo(context.Background(), dto.UpdateTodoRequest{
			ID: todoID,
		}, userID)

		s.ErrorIs(err, errorTest)
		s.Nil(result)
	})

	s.Run("Successfully update todo", func() {
		todoRet := *emptyTodo
		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetTodoByID(gomock.Any(), gomock.Any(), todoID).Return(&todoRet, nil)
		s.repo.EXPECT().UpdateTodo(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		s.cache.EXPECT().Del(keyFindTodo).Return(nil)
		s.cache.EXPECT().Del(keyFindAll).Return(nil)
		result, err := s.todoService.UpdateTodo(context.Background(), dto.UpdateTodoRequest{
			ID:    todoID,
			Title: "Todo",
		}, userID)

		s.Nil(err)
		s.NotEqual(emptyTodo, result)
	})
}

func (s *TodoTestSuite) TestDeleteTodo() {
	todoID := uuid.NewString()
	userID := uuid.NewString()
	emptyTodo := &entity.Todo{}
	emptyTodo.ID = uuid.MustParse(todoID)
	emptyTodo.UserID = uuid.MustParse(userID)
	keyFindAll := "todos:all:" + userID
	keyFindTodo := "todos:" + todoID

	s.Run("Failed to get todo", func() {
		errorTest := errors.New("get todo error")
		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetTodoByID(gomock.Any(), gomock.Any(), todoID).Return(nil, errorTest)
		err := s.todoService.DeleteTodo(context.Background(), todoID, userID)

		s.ErrorIs(err, errorTest)
	})

	s.Run("User ID mismatch", func() {
		var e *echo.HTTPError
		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetTodoByID(gomock.Any(), gomock.Any(), todoID).Return(emptyTodo, nil)
		err := s.todoService.DeleteTodo(context.Background(), todoID, uuid.NewString())

		s.ErrorAs(err, &e)
	})

	s.Run("Failed to delete todo", func() {
		errorTest := errors.New("delete todo error")
		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetTodoByID(gomock.Any(), gomock.Any(), todoID).Return(emptyTodo, nil)
		s.repo.EXPECT().DeleteTodo(gomock.Any(), gomock.Any(), emptyTodo).Return(errorTest)
		err := s.todoService.DeleteTodo(context.Background(), todoID, userID)

		s.ErrorIs(err, errorTest)
	})

	s.Run("Failed to delete todo cache", func() {
		errorTest := errors.New("delete todo cache error")
		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetTodoByID(gomock.Any(), gomock.Any(), todoID).Return(emptyTodo, nil)
		s.repo.EXPECT().DeleteTodo(gomock.Any(), gomock.Any(), emptyTodo).Return(nil)
		s.cache.EXPECT().Del(keyFindTodo).Return(errorTest)
		err := s.todoService.DeleteTodo(context.Background(), todoID, userID)

		s.ErrorIs(err, errorTest)
	})

	s.Run("Failed to delete todos cache", func() {
		errorTest := errors.New("delete todos cache error")
		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetTodoByID(gomock.Any(), gomock.Any(), todoID).Return(emptyTodo, nil)
		s.repo.EXPECT().DeleteTodo(gomock.Any(), gomock.Any(), emptyTodo).Return(nil)
		s.cache.EXPECT().Del(keyFindTodo).Return(nil)
		s.cache.EXPECT().Del(keyFindAll).Return(errorTest)
		err := s.todoService.DeleteTodo(context.Background(), todoID, userID)

		s.ErrorIs(err, errorTest)
	})

	s.Run("Successfully delete todo", func() {
		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetTodoByID(gomock.Any(), gomock.Any(), todoID).Return(emptyTodo, nil)
		s.repo.EXPECT().DeleteTodo(gomock.Any(), gomock.Any(), emptyTodo).Return(nil)
		s.cache.EXPECT().Del(keyFindTodo).Return(nil)
		s.cache.EXPECT().Del(keyFindAll).Return(nil)
		err := s.todoService.DeleteTodo(context.Background(), todoID, userID)

		s.Nil(err)
	})
}
