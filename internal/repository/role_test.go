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

type RoleTestSuite struct {
	suite.Suite
	db   *gorm.DB
	mock sqlmock.Sqlmock
	repo repository.RoleRepository
}

func TestRoleRepository(t *testing.T) {
	suite.Run(t, new(RoleTestSuite))
}

func (s *RoleTestSuite) SetupSuite() {
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
	s.repo = repository.NewRoleRepository(s.db)
}

func (s *RoleTestSuite) AfterTest(string, string) {
	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.FailNow("Failed to meet expectations", err)
	}
}

func (s *RoleTestSuite) TestGetRoles() {
	s.Run("Failed to get roles", func() {
		s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "roles"`)).
			WillReturnError(gorm.ErrRecordNotFound)

		result, err := s.repo.GetRoles(context.Background(), s.db)
		s.ErrorAs(err, &gorm.ErrRecordNotFound)
		s.Nil(result)
	})

	s.Run("Get roles successfully", func() {
		s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "roles"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "auth_level"}).
				AddRow(uuid.NewString(), "Admin", 3).
				AddRow(uuid.NewString(), "Editor", 2).
				AddRow(uuid.NewString(), "User", 1))

		result, err := s.repo.GetRoles(context.Background(), s.db)
		s.Nil(err)
		s.Len(result, 3)
	})
}

func (s *RoleTestSuite) TestGetRoleByID() {
	s.Run("Role not found", func() {
		id := uuid.NewString()
		s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "roles" WHERE id = $1 ORDER BY "roles"."id" LIMIT $2`)).
			WithArgs(id, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		result, err := s.repo.GetRoleByID(context.Background(), s.db, id)
		s.ErrorAs(err, &gorm.ErrRecordNotFound)
		s.Nil(result)
	})

	s.Run("Get role successfully", func() {
		id := uuid.NewString()
		s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "roles" WHERE id = $1 ORDER BY "roles"."id" LIMIT $2`)).
			WithArgs(id, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "auth_level"}).
				AddRow(id, "Admin", 3))

		result, err := s.repo.GetRoleByID(context.Background(), s.db, id)
		s.Nil(err)
		s.NotNil(result)
		s.Equal(id, result.ID.String())
	})
}

func (s *RoleTestSuite) TestCreateRole() {
	s.Run("Failed to create role", func() {
		s.mock.ExpectBegin()
		s.mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "roles"`)).
			WillReturnError(gorm.ErrInvalidData)
		s.mock.ExpectRollback()

		err := s.repo.CreateRole(context.Background(), s.db, &entity.Role{})
		s.ErrorAs(err, &gorm.ErrInvalidData)
	})

	s.Run("Create role successfully", func() {
		s.mock.ExpectBegin()
		s.mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "roles"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectCommit()

		err := s.repo.CreateRole(context.Background(), s.db, &entity.Role{})
		s.Nil(err)
	})
}

func (s *RoleTestSuite) TestUpdateRole() {
	s.Run("Failed to update role", func() {
		role := &entity.Role{}
		role.ID = uuid.Must(uuid.NewV7())
		s.mock.ExpectBegin()
		s.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "roles"`)).
			WillReturnError(gorm.ErrInvalidData)
		s.mock.ExpectRollback()

		err := s.repo.UpdateRole(context.Background(), s.db, role)
		s.ErrorAs(err, &gorm.ErrInvalidData)
	})

	s.Run("Update role successfully", func() {
		role := &entity.Role{}
		role.ID = uuid.Must(uuid.NewV7())
		s.mock.ExpectBegin()
		s.mock.ExpectExec(regexp.QuoteMeta(`UPDATE "roles"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectCommit()

		err := s.repo.UpdateRole(context.Background(), s.db, role)
		s.Nil(err)
	})
}

func (s *RoleTestSuite) TestDeleteRole() {
	s.Run("Failed to delete role", func() {
		role := &entity.Role{}
		role.ID = uuid.Must(uuid.NewV7())
		s.mock.ExpectBegin()
		s.mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "roles" WHERE "roles"."id" = $1`)).
			WithArgs(role.ID).
			WillReturnError(gorm.ErrInvalidData)
		s.mock.ExpectRollback()

		err := s.repo.DeleteRole(context.Background(), s.db, role)
		s.ErrorAs(err, &gorm.ErrInvalidData)
	})

	s.Run("Delete role successfully", func() {
		role := &entity.Role{}
		role.ID = uuid.Must(uuid.NewV7())
		s.mock.ExpectBegin()
		s.mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "roles" WHERE "roles"."id" = $1`)).
			WithArgs(role.ID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		s.mock.ExpectCommit()

		err := s.repo.DeleteRole(context.Background(), s.db, role)
		s.Nil(err)
	})
}
