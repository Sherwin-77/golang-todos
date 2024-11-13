package repository_test

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/sherwin-77/golang-todos/internal/entity"
	"github.com/sherwin-77/golang-todos/internal/repository"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"regexp"
	"testing"
)

type TodoTestSuite struct {
	suite.Suite
	db   *gorm.DB
	mock sqlmock.Sqlmock
	repo repository.TodoRepository
}

func TestTodoTestSuite(t *testing.T) {
	suite.Run(t, new(TodoTestSuite))
}

func (s *TodoTestSuite) SetupSuite() {
	db, mock, err := sqlmock.New()
	if err != nil {
		s.FailNow("Failed to create mock db", err.Error())
	}
	s.db, err = gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		s.FailNow("Failed to open mock db", err)
	}

	s.mock = mock
	s.repo = repository.NewTodoRepository(s.db)
}

func (s *TodoTestSuite) AfterTest(string, string) {
	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.FailNow("Failed to meet expectations", err)
	}
}

func (s *TodoTestSuite) TestGetTodosByUserID() {
	userID := uuid.NewString()

	s.Run("Failed to get todos", func() {
		s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "todos" WHERE user_id = $1`)).
			WithArgs(userID).
			WillReturnError(gorm.ErrRecordNotFound)

		result, err := s.repo.GetTodosByUserID(context.Background(), s.db, userID)
		s.ErrorAs(err, &gorm.ErrRecordNotFound)
		s.Nil(result)
	})

	s.Run("Get todos successfully", func() {
		s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "todos" WHERE user_id = $1`)).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
				AddRow(uuid.NewString(), "Todo 1").
				AddRow(uuid.NewString(), "Todo 2"))

		result, err := s.repo.GetTodosByUserID(context.Background(), s.db, userID)
		s.Nil(err)
		s.Len(result, 2)
	})
}

func (s *TodoTestSuite) TestGetTodosFiltered() {
	todoID := uuid.NewString()

	s.Run("Failed to get todos", func() {
		s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "todos" WHERE id != $1 ORDER BY id LIMIT $2 OFFSET $3`)).
			WithArgs(todoID, 1, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		result, err := s.repo.GetTodosFiltered(context.Background(), s.db, 1, 1, "id", "id != ?", todoID)
		s.ErrorAs(err, &gorm.ErrRecordNotFound)
		s.Nil(result)
	})

	s.Run("Get todos successfully", func() {
		s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "todos" WHERE id != $1 ORDER BY id LIMIT $2 OFFSET $3`)).
			WithArgs(todoID, 1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
				AddRow(uuid.NewString(), "Todo 1").
				AddRow(uuid.NewString(), "Todo 2"))

		result, err := s.repo.GetTodosFiltered(context.Background(), s.db, 1, 1, "id", "id != ?", todoID)
		s.Nil(err)
		s.Len(result, 2)
	})
}

func (s *TodoTestSuite) TestGetTodoByID() {
	todoID := uuid.NewString()

	s.Run("Todo not found", func() {
		s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "todos" WHERE id = $1 ORDER BY "todos"."id" LIMIT $2`)).
			WithArgs(todoID, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		result, err := s.repo.GetTodoByID(context.Background(), s.db, todoID)
		s.ErrorAs(err, &gorm.ErrRecordNotFound)
		s.Nil(result)
	})

	s.Run("Get todo successfully", func() {
		s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "todos" WHERE id = $1 ORDER BY "todos"."id" LIMIT $2`)).
			WithArgs(todoID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(todoID, "Todo 1"))

		result, err := s.repo.GetTodoByID(context.Background(), s.db, todoID)
		s.Nil(err)
		s.NotNil(result)
		s.Equal(todoID, result.ID.String())
	})
}

func (s *TodoTestSuite) TestCreateTodo() {
	s.Run("Failed to create todo", func() {
		s.mock.ExpectBegin()
		s.mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "todos"`)).WillReturnError(gorm.ErrInvalidData)
		s.mock.ExpectRollback()

		err := s.repo.CreateTodo(context.Background(), s.db, &entity.Todo{})
		s.ErrorAs(err, &gorm.ErrInvalidData)
	})

	s.Run("Create todo successfully", func() {
		s.mock.ExpectBegin()
		s.mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "todos"`)).WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectCommit()

		err := s.repo.CreateTodo(context.Background(), s.db, &entity.Todo{})
		s.Nil(err)
	})
}

func (s *TodoTestSuite) TestUpdateTodo() {
	s.Run("Failed to update todo", func() {
		todo := &entity.Todo{}
		todo.ID = uuid.Must(uuid.NewV7())
		s.mock.ExpectBegin()
		s.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "todos"`)).WillReturnError(gorm.ErrInvalidData)
		s.mock.ExpectRollback()

		err := s.repo.UpdateTodo(context.Background(), s.db, todo)
		s.ErrorAs(err, &gorm.ErrInvalidData)
	})

	s.Run("Update todo successfully", func() {
		todo := &entity.Todo{}
		todo.ID = uuid.Must(uuid.NewV7())
		s.mock.ExpectBegin()
		s.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "todos"`)).WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectCommit()

		err := s.repo.UpdateTodo(context.Background(), s.db, todo)
		s.Nil(err)
	})
}

func (s *TodoTestSuite) TestDeleteTodo() {
	s.Run("Failed to delete todo", func() {
		todo := &entity.Todo{}
		todo.ID = uuid.Must(uuid.NewV7())
		s.mock.ExpectBegin()
		s.mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "todos" WHERE "todos"."id" = $1`)).
			WithArgs(todo.ID).
			WillReturnError(gorm.ErrInvalidData)
		s.mock.ExpectRollback()

		err := s.repo.DeleteTodo(context.Background(), s.db, todo)
		s.ErrorAs(err, &gorm.ErrInvalidData)
	})

	s.Run("Delete todo successfully", func() {
		todo := &entity.Todo{}
		todo.ID = uuid.Must(uuid.NewV7())
		s.mock.ExpectBegin()
		s.mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "todos" WHERE "todos"."id" = $1`)).
			WithArgs(todo.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectCommit()

		err := s.repo.DeleteTodo(context.Background(), s.db, todo)
		s.Nil(err)
	})
}
