package repository_test

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sherwin-77/go-echo-template/internal/repository"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type BaseTestSuite struct {
	suite.Suite
	db   *gorm.DB
	mock sqlmock.Sqlmock
	repo repository.BaseRepository
}

func TestBaseRepository(t *testing.T) {
	suite.Run(t, new(BaseTestSuite))
}

func (s *BaseTestSuite) SetupSuite() {
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

func (s *BaseTestSuite) AfterTest(string, string) {
	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.FailNow("Failed to meet expectations", err)
	}
}

func (s *BaseTestSuite) TestWithTransaction() {
	s.Run("Failed to start transaction", func() {
		s.mock.ExpectBegin().WillReturnError(gorm.ErrInvalidTransaction)
		err := s.repo.WithTransaction(func(tx *gorm.DB) error {
			return nil
		})
		s.ErrorAs(err, &gorm.ErrInvalidTransaction)
	})

	s.Run("Failed to commit transaction", func() {
		s.mock.ExpectBegin()
		s.mock.ExpectCommit().WillReturnError(gorm.ErrInvalidTransaction)
		err := s.repo.WithTransaction(func(tx *gorm.DB) error {
			return nil
		})
		s.ErrorAs(err, &gorm.ErrInvalidTransaction)
	})

	s.Run("Rollback transaction successfully", func() {
		s.mock.ExpectBegin()
		s.mock.ExpectRollback()
		err := s.repo.WithTransaction(func(tx *gorm.DB) error {
			return gorm.ErrInvalidTransaction
		})
		s.ErrorAs(err, &gorm.ErrInvalidTransaction)
	})

	s.Run("Commit transaction successfully", func() {
		s.mock.ExpectBegin()
		s.mock.ExpectCommit()
		err := s.repo.WithTransaction(func(tx *gorm.DB) error {
			return nil
		})
		s.Nil(err)
	})
}

func (s *BaseTestSuite) TestControlledTransaction() {
	s.Run("Failed to start transaction", func() {
		s.mock.ExpectBegin().WillReturnError(gorm.ErrInvalidTransaction)
		db := s.repo.BeginTransaction()
		err := s.repo.Commit(db)
		s.ErrorAs(err, &gorm.ErrInvalidTransaction)
	})

	s.Run("Failed to commit transaction", func() {
		s.mock.ExpectBegin()
		s.mock.ExpectCommit().WillReturnError(gorm.ErrInvalidTransaction)
		db := s.repo.BeginTransaction()
		err := s.repo.Commit(db)
		s.ErrorAs(err, &gorm.ErrInvalidTransaction)
	})

	s.Run("Rollback transaction successfully", func() {
		s.mock.ExpectBegin()
		s.mock.ExpectRollback()
		db := s.repo.BeginTransaction()
		s.repo.Rollback(db)
	})

	s.Run("Commit transaction successfully", func() {
		s.mock.ExpectBegin()
		s.mock.ExpectCommit()
		db := s.repo.BeginTransaction()
		s.repo.Commit(db)
	})
}
