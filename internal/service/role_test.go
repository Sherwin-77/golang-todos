package service_test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/sherwin-77/go-echo-template/internal/entity"
	"github.com/sherwin-77/go-echo-template/internal/http/dto"
	"testing"

	"github.com/sherwin-77/go-echo-template/internal/service"
	mock_caches "github.com/sherwin-77/go-echo-template/test/mock/pkg/caches"
	mock_repository "github.com/sherwin-77/go-echo-template/test/mock/repository"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type RoleTestSuite struct {
	suite.Suite
	ctrl        *gomock.Controller
	repo        *mock_repository.MockRoleRepository
	userRepo    *mock_repository.MockUserRepository
	cache       *mock_caches.MockCache
	roleService service.RoleService
}

func (s *RoleTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.repo = mock_repository.NewMockRoleRepository(s.ctrl)
	s.userRepo = mock_repository.NewMockUserRepository(s.ctrl)
	s.cache = mock_caches.NewMockCache(s.ctrl)
	s.roleService = service.NewRoleService(s.repo, s.cache)
}

func TestRoleService(t *testing.T) {
	suite.Run(t, new(RoleTestSuite))
}
func (s *RoleTestSuite) TestGetRoles() {
	keyFindAll := "roles:all"
	roles := make([]entity.Role, 0)
	marshalledData, _ := json.Marshal(roles)

	s.Run("Failed unmarshal", func() {
		s.cache.EXPECT().Get(keyFindAll).Return("invalid")
		result, err := s.roleService.GetRoles(context.Background())

		s.Error(err)
		s.Nil(result)
	})

	s.Run("Failed to get roles", func() {
		s.cache.EXPECT().Get(keyFindAll).Return("")
		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetRoles(gomock.Any(), gomock.Any()).Return(nil, errors.New("get roles error"))
		result, err := s.roleService.GetRoles(context.Background())

		s.Error(err)
		s.Nil(result)
	})

	s.Run("Failed to set cache", func() {
		s.cache.EXPECT().Get(keyFindAll).Return("")
		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetRoles(gomock.Any(), gomock.Any()).Return(roles, nil)
		s.cache.EXPECT().Set(keyFindAll, string(marshalledData), gomock.Any()).Return(errors.New("set cache error"))
		result, err := s.roleService.GetRoles(context.Background())

		s.Error(err)
		s.Nil(result)
	})

	s.Run("Get roles successfully", func() {
		s.cache.EXPECT().Get(keyFindAll).Return("")
		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetRoles(gomock.Any(), gomock.Any()).Return(roles, nil)
		s.cache.EXPECT().Set(keyFindAll, string(marshalledData), gomock.Any()).Return(nil)
		result, err := s.roleService.GetRoles(context.Background())

		s.Nil(err)
		s.Equal(roles, result)
	})

	s.Run("Get roles from cache", func() {
		s.cache.EXPECT().Get(keyFindAll).Return(string(marshalledData))
		result, err := s.roleService.GetRoles(context.Background())

		s.Nil(err)
		s.Equal(roles, result)
	})
}

func (s *RoleTestSuite) TestGetRoleByID() {
	roleId := uuid.NewString()
	keyFindRole := "roles:" + roleId
	role := &entity.Role{}
	role.ID = uuid.MustParse(roleId)
	marshalledData, _ := json.Marshal(role)

	s.Run("Failed unmarshal", func() {
		s.cache.EXPECT().Get(keyFindRole).Return("invalid")
		result, err := s.roleService.GetRoleByID(context.Background(), roleId)

		s.Error(err)
		s.Nil(result)
	})

	s.Run("Failed to get role", func() {
		errorTest := errors.New("get role error")

		s.cache.EXPECT().Get(keyFindRole).Return("")
		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetRoleByID(gomock.Any(), gomock.Any(), roleId).Return(nil, errorTest)
		result, err := s.roleService.GetRoleByID(context.Background(), roleId)

		s.ErrorIs(err, errorTest)
		s.Nil(result)
	})

	s.Run("Failed to set cache", func() {
		errorTest := errors.New("set cache error")

		s.cache.EXPECT().Get(keyFindRole).Return("")
		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetRoleByID(gomock.Any(), gomock.Any(), roleId).Return(&*role, nil)
		s.cache.EXPECT().Set(keyFindRole, string(marshalledData), gomock.Any()).Return(errorTest)
		result, err := s.roleService.GetRoleByID(context.Background(), roleId)

		s.ErrorIs(err, errorTest)
		s.Nil(result)
	})

	s.Run("Get role successfully", func() {
		s.cache.EXPECT().Get(keyFindRole).Return("")
		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetRoleByID(gomock.Any(), gomock.Any(), roleId).Return(&*role, nil)
		s.cache.EXPECT().Set(keyFindRole, string(marshalledData), gomock.Any()).Return(nil)
		result, err := s.roleService.GetRoleByID(context.Background(), roleId)

		s.Nil(err)
		s.Equal(role, result)
	})

	s.Run("Get role from cache", func() {
		s.cache.EXPECT().Get(keyFindRole).Return(string(marshalledData))
		result, err := s.roleService.GetRoleByID(context.Background(), roleId)

		s.Nil(err)
		s.Equal(role, result)
	})
}

func (s *RoleTestSuite) TestCreateRole() {
	emptyRole := &entity.Role{}

	s.Run("Failed to create role", func() {
		errorTest := errors.New("create role error")

		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().CreateRole(gomock.Any(), gomock.Any(), gomock.Any()).Return(errorTest)
		result, err := s.roleService.CreateRole(context.Background(), dto.RoleRequest{})

		s.ErrorIs(err, errorTest)
		s.Nil(result)
	})

	s.Run("Failed to delete cache", func() {
		errorTest := errors.New("delete cache error")

		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().CreateRole(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		s.cache.EXPECT().Del("roles:all").Return(errorTest)
		result, err := s.roleService.CreateRole(context.Background(), dto.RoleRequest{})

		s.ErrorIs(err, errorTest)
		s.Nil(result)
	})

	s.Run("Create role successfully", func() {
		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().CreateRole(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		s.cache.EXPECT().Del("roles:all").Return(nil)
		result, err := s.roleService.CreateRole(context.Background(), dto.RoleRequest{
			Name:      "Admin",
			AuthLevel: 3,
		})

		s.Nil(err)
		s.NotEqual(emptyRole, result)
	})
}

func (s *RoleTestSuite) TestUpdateRole() {
	roleID := uuid.NewString()
	emptyRole := &entity.Role{}
	emptyRole.ID = uuid.MustParse(roleID)

	s.Run("Failed to get role", func() {
		errorTest := errors.New("get role error")

		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetRoleByID(gomock.Any(), gomock.Any(), roleID).Return(nil, errorTest)
		result, err := s.roleService.UpdateRole(context.Background(), dto.UpdateRoleRequest{
			ID: roleID,
		})

		s.ErrorIs(err, errorTest)
		s.Nil(result)
	})

	s.Run("Failed to update role", func() {
		roleRet := *emptyRole
		errorTest := errors.New("update role error")

		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetRoleByID(gomock.Any(), gomock.Any(), roleID).Return(&roleRet, nil)
		s.repo.EXPECT().UpdateRole(gomock.Any(), gomock.Any(), gomock.Any()).Return(errorTest)
		result, err := s.roleService.UpdateRole(context.Background(), dto.UpdateRoleRequest{
			ID: roleID,
			RoleRequest: dto.RoleRequest{
				Name:      "Admin",
				AuthLevel: 3,
			},
		})

		s.ErrorIs(err, errorTest)
		s.Nil(result)
	})

	s.Run("Failed to delete role cache", func() {
		roleRet := *emptyRole
		errorTest := errors.New("delete cache error")

		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetRoleByID(gomock.Any(), gomock.Any(), roleID).Return(&roleRet, nil)
		s.repo.EXPECT().UpdateRole(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		s.cache.EXPECT().Del("roles:" + roleID).Return(errorTest)
		result, err := s.roleService.UpdateRole(context.Background(), dto.UpdateRoleRequest{
			ID: roleID,
			RoleRequest: dto.RoleRequest{
				Name:      "Admin",
				AuthLevel: 3,
			},
		})

		s.ErrorIs(err, errorTest)
		s.Nil(result)
	})

	s.Run("Failed to delete roles cache", func() {
		roleRet := *emptyRole
		errorTest := errors.New("delete cache error")

		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetRoleByID(gomock.Any(), gomock.Any(), roleID).Return(&roleRet, nil)
		s.repo.EXPECT().UpdateRole(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		s.cache.EXPECT().Del("roles:" + roleID).Return(nil)
		s.cache.EXPECT().Del("roles:all").Return(errorTest)
		result, err := s.roleService.UpdateRole(context.Background(), dto.UpdateRoleRequest{
			ID: roleID,
			RoleRequest: dto.RoleRequest{
				Name:      "Admin",
				AuthLevel: 3,
			},
		})

		s.ErrorIs(err, errorTest)
		s.Nil(result)
	})

	s.Run("Update role successfully", func() {
		roleRet := *emptyRole

		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetRoleByID(gomock.Any(), gomock.Any(), roleID).Return(&roleRet, nil)
		s.repo.EXPECT().UpdateRole(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
		s.cache.EXPECT().Del("roles:" + roleID).Return(nil)
		s.cache.EXPECT().Del("roles:all").Return(nil)
		result, err := s.roleService.UpdateRole(context.Background(), dto.UpdateRoleRequest{
			ID: roleID,
			RoleRequest: dto.RoleRequest{
				Name:      "Admin",
				AuthLevel: 3,
			}})

		s.Nil(err)
		s.NotEqual(emptyRole, result)
	})
}

func (s *RoleTestSuite) TestDeleteRole() {
	roleID := uuid.NewString()
	emptyRole := &entity.Role{}
	emptyRole.ID = uuid.MustParse(roleID)

	s.Run("Failed to get role", func() {
		errorTest := errors.New("get role error")
		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetRoleByID(gomock.Any(), gomock.Any(), roleID).Return(nil, errorTest)
		err := s.roleService.DeleteRole(context.Background(), roleID)

		s.ErrorIs(err, errorTest)
	})

	s.Run("Failed to delete role", func() {
		errorTest := errors.New("delete role error")
		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetRoleByID(gomock.Any(), gomock.Any(), roleID).Return(emptyRole, nil)
		s.repo.EXPECT().DeleteRole(gomock.Any(), gomock.Any(), emptyRole).Return(errorTest)
		err := s.roleService.DeleteRole(context.Background(), roleID)

		s.ErrorIs(err, errorTest)
	})

	s.Run("Failed to delete role cache", func() {
		errorTest := errors.New("delete cache error")
		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetRoleByID(gomock.Any(), gomock.Any(), roleID).Return(emptyRole, nil)
		s.repo.EXPECT().DeleteRole(gomock.Any(), gomock.Any(), emptyRole).Return(nil)
		s.cache.EXPECT().Del("roles:" + roleID).Return(errorTest)
		err := s.roleService.DeleteRole(context.Background(), roleID)

		s.ErrorIs(err, errorTest)
	})

	s.Run("Failed to delete roles cache", func() {
		errorTest := errors.New("delete cache error")
		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetRoleByID(gomock.Any(), gomock.Any(), roleID).Return(emptyRole, nil)
		s.repo.EXPECT().DeleteRole(gomock.Any(), gomock.Any(), emptyRole).Return(nil)
		s.cache.EXPECT().Del("roles:" + roleID).Return(nil)
		s.cache.EXPECT().Del("roles:all").Return(errorTest)
		err := s.roleService.DeleteRole(context.Background(), roleID)

		s.ErrorIs(err, errorTest)
	})

	s.Run("Delete role successfully", func() {
		s.repo.EXPECT().SingleTransaction().Return(nil)
		s.repo.EXPECT().GetRoleByID(gomock.Any(), gomock.Any(), roleID).Return(emptyRole, nil)
		s.repo.EXPECT().DeleteRole(gomock.Any(), gomock.Any(), emptyRole).Return(nil)
		s.cache.EXPECT().Del("roles:" + roleID).Return(nil)
		s.cache.EXPECT().Del("roles:all").Return(nil)
		err := s.roleService.DeleteRole(context.Background(), roleID)

		s.Nil(err)
	})
}
