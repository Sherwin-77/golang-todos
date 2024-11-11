package repository_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/sherwin-77/go-echo-template/internal/entity"
	"github.com/sherwin-77/go-echo-template/internal/repository"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type UserTestSuite struct {
	suite.Suite
	db   *gorm.DB
	mock sqlmock.Sqlmock
	repo repository.UserRepository
}

func TestUserRepository(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}

func (s *UserTestSuite) SetupSuite() {
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
	s.repo = repository.NewUserRepository(s.db)
}

func (s *UserTestSuite) AfterTest(string, string) {
	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.FailNow("Failed to meet expectations", err)
	}
}

func (s *UserTestSuite) TestGetUsers() {
	s.Run("Failed to get users", func() {
		s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
			WillReturnError(gorm.ErrRecordNotFound)

		result, err := s.repo.GetUsers(context.Background(), s.db)
		s.ErrorAs(err, &gorm.ErrRecordNotFound)
		s.Nil(result)
	})

	s.Run("Get users successfully", func() {
		s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password"}).
				AddRow(uuid.NewString(), "admin", "password").
				AddRow(uuid.NewString(), "editor", "password"))

		result, err := s.repo.GetUsers(context.Background(), s.db)
		s.Nil(err)
		s.Len(result, 2)
	})
}

func (s *UserTestSuite) TestGetUserByID() {
	s.Run("User not found", func() {
		id := uuid.NewString()
		s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 ORDER BY "users"."id" LIMIT $2`)).
			WithArgs(id, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		result, err := s.repo.GetUserByID(context.Background(), s.db, id)
		s.ErrorAs(err, &gorm.ErrRecordNotFound)
		s.Nil(result)
	})

	s.Run("Get user successfully", func() {
		id := uuid.NewString()
		s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 ORDER BY "users"."id" LIMIT $2`)).
			WithArgs(id, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).
				AddRow(id))

		result, err := s.repo.GetUserByID(context.Background(), s.db, id)
		s.Nil(err)
		s.NotNil(result)
		s.Equal(id, result.ID.String())
	})
}

func (s *UserTestSuite) TestGetUserByEmail() {
	s.Run("User not found", func() {
		email := "admin@example.com"
		s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 ORDER BY "users"."id" LIMIT $2`)).
			WithArgs(email, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		result, err := s.repo.GetUserByEmail(context.Background(), s.db, email)
		s.ErrorAs(err, &gorm.ErrRecordNotFound)
		s.Nil(result)
	})

	s.Run("Get user successfully", func() {
		email := "admin@example.com"
		s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 ORDER BY "users"."id" LIMIT $2`)).
			WithArgs(email, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).
				AddRow(uuid.NewString(), email))

		result, err := s.repo.GetUserByEmail(context.Background(), s.db, email)
		s.Nil(err)
		s.NotNil(result)
		s.Equal(email, result.Email)
	})

}

func (s *UserTestSuite) TestCreateUser() {
	s.Run("Failed to create user", func() {
		user := &entity.User{}

		s.mock.ExpectBegin()
		s.mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "users"`)).
			WillReturnError(gorm.ErrInvalidData)
		s.mock.ExpectRollback()

		err := s.repo.CreateUser(context.Background(), s.db, user)
		s.ErrorAs(err, &gorm.ErrInvalidData)
	})

	s.Run("Create user successfully", func() {
		user := &entity.User{}

		s.mock.ExpectBegin()
		s.mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "users"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectCommit()

		err := s.repo.CreateUser(context.Background(), s.db, user)
		s.Nil(err)
	})
}

func (s *UserTestSuite) TestUpdateUser() {
	s.Run("Failed to update user", func() {
		user := &entity.User{}
		user.ID = uuid.Must(uuid.NewV7())

		s.mock.ExpectBegin()
		s.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users"`)).
			WillReturnError(gorm.ErrInvalidData)
		s.mock.ExpectRollback()

		err := s.repo.UpdateUser(context.Background(), s.db, user)
		s.ErrorAs(err, &gorm.ErrInvalidData)
	})

	s.Run("Update user successfully", func() {
		user := &entity.User{}
		user.ID = uuid.Must(uuid.NewV7())

		s.mock.ExpectBegin()
		s.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectCommit()

		err := s.repo.UpdateUser(context.Background(), s.db, user)
		s.Nil(err)
	})
}

func (s *UserTestSuite) TestDeleteUser() {
	s.Run("Failed to delete user", func() {
		user := &entity.User{}
		user.ID = uuid.Must(uuid.NewV7())

		s.mock.ExpectBegin()
		s.mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "users" WHERE "users"."id" = $1`)).
			WithArgs(user.ID).
			WillReturnError(gorm.ErrInvalidData)
		s.mock.ExpectRollback()

		err := s.repo.DeleteUser(context.Background(), s.db, user)
		s.ErrorAs(err, &gorm.ErrInvalidData)
	})

	s.Run("Delete user successfully", func() {
		user := &entity.User{}
		user.ID = uuid.Must(uuid.NewV7())

		s.mock.ExpectBegin()
		s.mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "users" WHERE "users"."id" = $1`)).
			WithArgs(user.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectCommit()

		err := s.repo.DeleteUser(context.Background(), s.db, user)
		s.Nil(err)
	})
}

func (s *UserTestSuite) TestAddRoles() {
	s.Run("Failed to add roles", func() {
		user := &entity.User{}
		user.ID = uuid.Must(uuid.NewV7())

		var roles []*entity.Role
		roles = append(roles, &entity.Role{})

		s.mock.ExpectBegin()
		s.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "roles"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "role_users"`)).
			WillReturnError(gorm.ErrInvalidData)
		s.mock.ExpectRollback()

		err := s.repo.AddRoles(context.Background(), s.db, user, roles)
		s.ErrorAs(err, &gorm.ErrInvalidData)
	})

	s.Run("Add roles successfully", func() {
		user := &entity.User{}
		user.ID = uuid.Must(uuid.NewV7())

		var roles []*entity.Role
		roles = append(roles, &entity.Role{})

		s.mock.ExpectBegin()
		s.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "roles"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "role_users"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectCommit()

		err := s.repo.AddRoles(context.Background(), s.db, user, roles)
		s.Nil(err)
	})
}

func (s *UserTestSuite) TestRemoveRoles() {
	s.Run("Failed to remove roles", func() {
		user := &entity.User{}
		user.ID = uuid.Must(uuid.NewV7())

		var roles []*entity.Role
		roles = append(roles, &entity.Role{})

		s.mock.ExpectBegin()
		s.mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "role_users" WHERE "role_users"."user_id" = $1`)).
			WithArgs(user.ID).
			WillReturnError(gorm.ErrInvalidData)
		s.mock.ExpectRollback()

		err := s.repo.RemoveRoles(context.Background(), s.db, user, roles)
		s.ErrorAs(err, &gorm.ErrInvalidData)
	})

	s.Run("Remove roles successfully", func() {
		user := &entity.User{}
		user.ID = uuid.Must(uuid.NewV7())

		var roles []*entity.Role
		roles = append(roles, &entity.Role{})

		s.mock.ExpectBegin()
		s.mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "role_users" WHERE "role_users"."user_id" = $1`)).
			WithArgs(user.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectCommit()

		err := s.repo.RemoveRoles(context.Background(), s.db, user, roles)
		s.Nil(err)
	})
}
