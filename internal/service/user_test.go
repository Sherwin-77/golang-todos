package service_test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sherwin-77/go-echo-template/internal/entity"
	"github.com/sherwin-77/go-echo-template/internal/http/dto"
	"github.com/sherwin-77/go-echo-template/internal/service"
	mock_caches "github.com/sherwin-77/go-echo-template/test/mock/pkg/caches"
	mock_tokens "github.com/sherwin-77/go-echo-template/test/mock/pkg/tokens"
	mock_repository "github.com/sherwin-77/go-echo-template/test/mock/repository"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"testing"
)

type UserTestSuite struct {
	suite.Suite
	ctrl         *gomock.Controller
	repo         *mock_repository.MockUserRepository
	roleRepo     *mock_repository.MockRoleRepository
	tokenService *mock_tokens.MockTokenService
	cache        *mock_caches.MockCache
	userService  service.UserService
}

func (s *UserTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.repo = mock_repository.NewMockUserRepository(s.ctrl)
	s.roleRepo = mock_repository.NewMockRoleRepository(s.ctrl)
	s.tokenService = mock_tokens.NewMockTokenService(s.ctrl)
	s.cache = mock_caches.NewMockCache(s.ctrl)
	s.userService = service.NewUserService(s.tokenService, s.repo, s.roleRepo, s.cache)
}

func TestUserService(t *testing.T) {
	suite.Run(t, new(UserTestSuite))
}

func (s *UserTestSuite) TestGetUsers() {
	keyFindAll := "users:all"
	users := []entity.User{
		{
			Username: "admin",
			Email:    "admin",
		},
	}
	marshalledData, _ := json.Marshal(users)

	s.Run("Failed unmarshal", func() {
		s.cache.EXPECT().Get(keyFindAll).Return("invalid")
		result, err := s.userService.GetUsers(context.Background())

		s.Error(err)
		s.Nil(result)
	})

	s.Run("Failed to get users", func() {
		s.cache.EXPECT().Get(keyFindAll).Return("")
		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetUsers(gomock.Any(), gomock.Any()).Return(nil, errors.New("get users error"))
		result, err := s.userService.GetUsers(context.Background())

		s.Error(err)
		s.Nil(result)
	})

	s.Run("Failed to set cache", func() {
		s.cache.EXPECT().Get(keyFindAll).Return("")
		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetUsers(gomock.Any(), gomock.Any()).Return(users, nil)
		s.cache.EXPECT().Set(keyFindAll, string(marshalledData), gomock.Any()).Return(errors.New("set cache error"))
		result, err := s.userService.GetUsers(context.Background())

		s.Error(err)
		s.Nil(result)
	})

	s.Run("Get users successfully", func() {
		s.cache.EXPECT().Get(keyFindAll).Return("")
		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetUsers(gomock.Any(), gomock.Any()).Return(users, nil)
		s.cache.EXPECT().Set(keyFindAll, string(marshalledData), gomock.Any()).Return(nil)
		result, err := s.userService.GetUsers(context.Background())

		s.Nil(err)
		s.Equal(users, result)
	})

	s.Run("Get users from cache", func() {
		s.cache.EXPECT().Get(keyFindAll).Return(string(marshalledData))
		result, err := s.userService.GetUsers(context.Background())

		s.Nil(err)
		s.Equal(users, result)
	})
}

func (s *UserTestSuite) TestGetUserByID() {
	userID := uuid.NewString()
	keyFindUser := "users:" + userID
	user := &entity.User{}
	user.ID = uuid.MustParse(userID)
	marshalledData, _ := json.Marshal(user)

	s.Run("Failed unmarshal", func() {
		s.cache.EXPECT().Get(keyFindUser).Return("invalid")
		result, err := s.userService.GetUserByID(context.Background(), userID)

		s.Error(err)
		s.Nil(result)
	})

	s.Run("Failed to get user", func() {
		errorTest := errors.New("get user error")

		s.cache.EXPECT().Get(keyFindUser).Return("")
		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetUserByID(gomock.Any(), gomock.Any(), userID).Return(nil, errorTest)
		result, err := s.userService.GetUserByID(context.Background(), userID)

		s.ErrorIs(err, errorTest)
		s.Nil(result)
	})

	s.Run("Failed to set cache", func() {
		errorTest := errors.New("set cache error")
		s.cache.EXPECT().Get(keyFindUser).Return("")
		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetUserByID(gomock.Any(), gomock.Any(), userID).Return(user, nil)
		s.cache.EXPECT().Set(keyFindUser, string(marshalledData), gomock.Any()).Return(errorTest)
		result, err := s.userService.GetUserByID(context.Background(), userID)

		s.ErrorIs(err, errorTest)
		s.Nil(result)
	})

	s.Run("Get user successfully", func() {
		s.cache.EXPECT().Get(keyFindUser).Return("")
		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetUserByID(gomock.Any(), gomock.Any(), userID).Return(user, nil)
		s.cache.EXPECT().Set(keyFindUser, string(marshalledData), gomock.Any()).Return(nil)
		result, err := s.userService.GetUserByID(context.Background(), userID)

		s.Nil(err)
		s.Equal(user, result)
	})

	s.Run("Get user from cache", func() {
		s.cache.EXPECT().Get(keyFindUser).Return(string(marshalledData))
		result, err := s.userService.GetUserByID(context.Background(), userID)

		s.Nil(err)
		s.Equal(user, result)
	})
}

func (s *UserTestSuite) TestCreateUser() {
	emptyUser := &entity.User{}

	s.Run("Failed to create user", func() {
		errorTest := errors.New("create user error")

		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().CreateUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(errorTest)
		result, err := s.userService.CreateUser(context.Background(), dto.UserRequest{})

		s.ErrorIs(err, errorTest)
		s.Nil(result)
	})

	s.Run("Failed to delete cache", func() {
		errorTest := errors.New("delete cache error")

		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().CreateUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		s.cache.EXPECT().Del("users:all").Return(errorTest)
		result, err := s.userService.CreateUser(context.Background(), dto.UserRequest{})

		s.ErrorIs(err, errorTest)
		s.Nil(result)
	})

	s.Run("Create user successfully", func() {
		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().CreateUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		s.cache.EXPECT().Del("users:all").Return(nil)
		result, err := s.userService.CreateUser(context.Background(), dto.UserRequest{
			Username: "admin",
		})

		s.Nil(err)
		s.NotEqual(emptyUser, result)
	})
}

func (s *UserTestSuite) TestUpdateUser() {
	userId := uuid.NewString()
	emptyUser := &entity.User{}
	emptyUser.ID = uuid.MustParse(userId)

	s.Run("Failed to get user", func() {
		errorTest := errors.New("get user error")

		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetUserByID(gomock.Any(), gomock.Any(), userId).Return(nil, errorTest)
		result, err := s.userService.UpdateUser(context.Background(), dto.UpdateUserRequest{
			ID: userId,
		})

		s.ErrorIs(err, errorTest)
		s.Nil(result)
	})

	s.Run("Failed to update user", func() {
		userRet := *emptyUser
		errorTest := errors.New("update user error")

		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetUserByID(gomock.Any(), gomock.Any(), userId).Return(&userRet, nil)
		s.repo.EXPECT().UpdateUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(errorTest)
		result, err := s.userService.UpdateUser(context.Background(), dto.UpdateUserRequest{
			ID:       userId,
			Username: "admin",
		})

		s.ErrorIs(err, errorTest)
		s.Nil(result)
	})

	s.Run("Failed to delete user cache", func() {
		userRet := *emptyUser
		errorTest := errors.New("delete user cache error")

		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetUserByID(gomock.Any(), gomock.Any(), userId).Return(&userRet, nil)
		s.repo.EXPECT().UpdateUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		s.cache.EXPECT().Del("users:" + userId).Return(errorTest)
		result, err := s.userService.UpdateUser(context.Background(), dto.UpdateUserRequest{
			ID:       userId,
			Username: "admin",
		})

		s.ErrorIs(err, errorTest)
		s.Nil(result)
	})

	s.Run("Failed to delete users cache", func() {
		userRet := *emptyUser
		errorTest := errors.New("delete users cache error")

		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetUserByID(gomock.Any(), gomock.Any(), userId).Return(&userRet, nil)
		s.repo.EXPECT().UpdateUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		s.cache.EXPECT().Del("users:" + userId).Return(nil)
		s.cache.EXPECT().Del("users:all").Return(errorTest)
		result, err := s.userService.UpdateUser(context.Background(), dto.UpdateUserRequest{
			ID:       userId,
			Username: "admin",
		})

		s.ErrorIs(err, errorTest)
		s.Nil(result)
	})

	s.Run("Update user successfully", func() {
		userRet := *emptyUser
		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetUserByID(gomock.Any(), gomock.Any(), userId).Return(&userRet, nil)
		s.repo.EXPECT().UpdateUser(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		s.cache.EXPECT().Del("users:" + userId).Return(nil)
		s.cache.EXPECT().Del("users:all").Return(nil)
		result, err := s.userService.UpdateUser(context.Background(), dto.UpdateUserRequest{
			ID:       userId,
			Username: "admin",
			Email:    "admin@example.com",
			Password: "admin#1234",
		})

		s.Nil(err)
		s.NotEqual(emptyUser, result)
	})
}

func (s *UserTestSuite) TestDeleteUser() {
	userID := uuid.NewString()
	emptyUser := &entity.User{}
	emptyUser.ID = uuid.MustParse(userID)

	s.Run("Failed to get user", func() {
		errorTest := errors.New("get user error")
		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetUserByID(gomock.Any(), gomock.Any(), userID).Return(nil, errorTest)
		err := s.userService.DeleteUser(context.Background(), userID)

		s.ErrorIs(err, errorTest)
	})

	s.Run("Failed to delete user", func() {
		errorTest := errors.New("delete user error")
		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetUserByID(gomock.Any(), gomock.Any(), userID).Return(emptyUser, nil)
		s.repo.EXPECT().DeleteUser(gomock.Any(), gomock.Any(), emptyUser).Return(errorTest)
		err := s.userService.DeleteUser(context.Background(), userID)

		s.ErrorIs(err, errorTest)
	})

	s.Run("Failed to delete user cache", func() {
		errorTest := errors.New("delete user cache error")
		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetUserByID(gomock.Any(), gomock.Any(), userID).Return(emptyUser, nil)
		s.repo.EXPECT().DeleteUser(gomock.Any(), gomock.Any(), emptyUser).Return(nil)
		s.cache.EXPECT().Del("users:" + userID).Return(errorTest)
		err := s.userService.DeleteUser(context.Background(), userID)

		s.ErrorIs(err, errorTest)
	})

	s.Run("Failed to delete users cache", func() {
		errorTest := errors.New("delete users cache error")
		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetUserByID(gomock.Any(), gomock.Any(), userID).Return(emptyUser, nil)
		s.repo.EXPECT().DeleteUser(gomock.Any(), gomock.Any(), emptyUser).Return(nil)
		s.cache.EXPECT().Del("users:" + userID).Return(nil)
		s.cache.EXPECT().Del("users:all").Return(errorTest)
		err := s.userService.DeleteUser(context.Background(), userID)

		s.ErrorIs(err, errorTest)
	})

	s.Run("Delete user successfully", func() {
		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetUserByID(gomock.Any(), gomock.Any(), userID).Return(emptyUser, nil)
		s.repo.EXPECT().DeleteUser(gomock.Any(), gomock.Any(), emptyUser).Return(nil)
		s.cache.EXPECT().Del("users:" + userID).Return(nil)
		s.cache.EXPECT().Del("users:all").Return(nil)
		err := s.userService.DeleteUser(context.Background(), userID)

		s.Nil(err)
	})
}

func (s *UserTestSuite) TestChangeRole() {
	userID := uuid.NewString()
	user := &entity.User{}
	user.ID = uuid.MustParse(userID)

	roleAdd := &entity.Role{}
	roleAdd.ID = uuid.New()
	roleAdd.Name = "admin"

	roleRemove := &entity.Role{}
	roleRemove.ID = uuid.New()
	roleRemove.Name = "user"

	request := dto.ChangeRoleRequest{
		UserID: userID,
		Items: []dto.ChangeRoleRequestItem{
			{
				ID:     roleAdd.ID.String(),
				Action: "add",
			},
			{
				ID:     roleRemove.ID.String(),
				Action: "remove",
			},
		},
	}

	s.Run("Failed to get user", func() {
		errorTest := errors.New("get user error")
		s.repo.EXPECT().WithTransaction(gomock.Any()).DoAndReturn(func(f func(tx *gorm.DB) error) error {
			s.repo.EXPECT().GetUserByID(gomock.Any(), gomock.Any(), userID).Return(nil, errorTest)

			return f(&gorm.DB{})
		})

		err := s.userService.ChangeRole(context.Background(), request)
		s.ErrorIs(err, errorTest)
	})

	s.Run("Failed to get role", func() {
		errorTest := errors.New("get role error")
		s.repo.EXPECT().WithTransaction(gomock.Any()).DoAndReturn(func(f func(tx *gorm.DB) error) error {
			s.repo.EXPECT().GetUserByID(gomock.Any(), gomock.Any(), userID).Return(user, nil)
			s.roleRepo.EXPECT().GetRoleByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errorTest)

			return f(&gorm.DB{})
		})

		err := s.userService.ChangeRole(context.Background(), request)
		s.ErrorIs(err, errorTest)
	})

	s.Run("Invalid action", func() {
		var e *echo.HTTPError
		invalidRequest := dto.ChangeRoleRequest{
			UserID: userID,
			Items: []dto.ChangeRoleRequestItem{
				{
					ID:     roleAdd.ID.String(),
					Action: "invalid",
				},
			},
		}
		s.repo.EXPECT().WithTransaction(gomock.Any()).DoAndReturn(func(f func(tx *gorm.DB) error) error {
			s.repo.EXPECT().GetUserByID(gomock.Any(), gomock.Any(), userID).Return(user, nil)
			s.roleRepo.EXPECT().GetRoleByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(roleAdd, nil)

			return f(&gorm.DB{})
		})

		err := s.userService.ChangeRole(context.Background(), invalidRequest)
		s.ErrorAs(err, &e)
	})

	s.Run("Failed add role", func() {
		errorTest := errors.New("add role error")
		s.repo.EXPECT().WithTransaction(gomock.Any()).DoAndReturn(func(f func(tx *gorm.DB) error) error {
			s.repo.EXPECT().GetUserByID(gomock.Any(), gomock.Any(), userID).Return(user, nil)
			s.roleRepo.EXPECT().GetRoleByID(gomock.Any(), gomock.Any(), roleAdd.ID.String()).Return(roleAdd, nil)
			s.roleRepo.EXPECT().GetRoleByID(gomock.Any(), gomock.Any(), roleRemove.ID.String()).Return(roleRemove, nil)
			s.repo.EXPECT().AddRoles(gomock.Any(), gomock.Any(), user, gomock.Any()).Return(errorTest)

			return f(&gorm.DB{})
		})

		err := s.userService.ChangeRole(context.Background(), request)
		s.ErrorIs(err, errorTest)
	})

	s.Run("Failed remove role", func() {
		errorTest := errors.New("remove role error")
		s.repo.EXPECT().WithTransaction(gomock.Any()).DoAndReturn(func(f func(tx *gorm.DB) error) error {
			s.repo.EXPECT().GetUserByID(gomock.Any(), gomock.Any(), userID).Return(user, nil)
			s.roleRepo.EXPECT().GetRoleByID(gomock.Any(), gomock.Any(), roleAdd.ID.String()).Return(roleAdd, nil)
			s.roleRepo.EXPECT().GetRoleByID(gomock.Any(), gomock.Any(), roleRemove.ID.String()).Return(roleRemove, nil)
			s.repo.EXPECT().AddRoles(gomock.Any(), gomock.Any(), user, gomock.Any()).Return(nil)
			s.repo.EXPECT().RemoveRoles(gomock.Any(), gomock.Any(), user, gomock.Any()).Return(errorTest)

			return f(&gorm.DB{})
		})

		err := s.userService.ChangeRole(context.Background(), request)
		s.ErrorIs(err, errorTest)
	})
}

func (s *UserTestSuite) TestLogin() {
	s.Run("Invalid email or password", func() {
		var e *echo.HTTPError
		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("invalid email or password"))

		result, err := s.userService.Login(context.Background(), dto.LoginRequest{
			Email:    "admin",
			Password: "admin",
		})

		s.ErrorAs(err, &e)
		s.Empty(result)
	})

	s.Run("Login successfully", func() {
		pass, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetUserByEmail(gomock.Any(), gomock.Any(), gomock.Any()).Return(&entity.User{
			Email:    "admin",
			Password: string(pass),
		}, nil)
		s.tokenService.EXPECT().GenerateAccessToken(gomock.Any()).Return("token", nil)
		result, err := s.userService.Login(context.Background(), dto.LoginRequest{
			Email:    "admin",
			Password: "admin",
		})

		s.Nil(err)
		s.NotEmpty(result)
	})
}
