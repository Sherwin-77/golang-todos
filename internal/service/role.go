package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/sherwin-77/go-echo-template/internal/entity"
	"github.com/sherwin-77/go-echo-template/internal/http/dto"
	"github.com/sherwin-77/go-echo-template/internal/repository"
	"github.com/sherwin-77/go-echo-template/pkg/caches"
)

type RoleService interface {
	GetRoles(ctx context.Context) ([]entity.Role, error)
	GetRoleByID(ctx context.Context, id string) (*entity.Role, error)
	CreateRole(ctx context.Context, request dto.RoleRequest) (*entity.Role, error)
	UpdateRole(ctx context.Context, request dto.UpdateRoleRequest) (*entity.Role, error)
	DeleteRole(ctx context.Context, id string) error
}

type roleService struct {
	roleRepository repository.RoleRepository
	cache          caches.Cache
}

func NewRoleService(roleRepository repository.RoleRepository, cache caches.Cache) RoleService {
	return &roleService{roleRepository, cache}
}

func (s *roleService) GetRoles(ctx context.Context) ([]entity.Role, error) {
	roleKey := "roles:all"
	var roles []entity.Role
	cachedData := s.cache.Get(roleKey)
	if cachedData != "" {
		if err := json.Unmarshal([]byte(cachedData), &roles); err != nil {
			return nil, err
		}
	} else {
		var err error
		db := s.roleRepository.SingleTransaction()
		roles, err = s.roleRepository.GetRoles(ctx, db)
		if err != nil {
			return nil, err
		}

		data, _ := json.Marshal(roles)

		if err := s.cache.Set(roleKey, string(data), 5*time.Minute); err != nil {
			return nil, err
		}
	}

	return roles, nil
}

func (s *roleService) GetRoleByID(ctx context.Context, id string) (*entity.Role, error) {
	roleKey := "roles:" + id
	role := &entity.Role{}
	cachedData := s.cache.Get(roleKey)
	if cachedData != "" {
		if err := json.Unmarshal([]byte(cachedData), role); err != nil {
			return nil, err
		}
	} else {
		var err error
		db := s.roleRepository.SingleTransaction()
		role, err = s.roleRepository.GetRoleByID(ctx, db, id)
		if err != nil {
			return nil, err
		}

		data, _ := json.Marshal(role)

		if err := s.cache.Set(roleKey, string(data), 5*time.Minute); err != nil {
			return nil, err
		}
	}

	return role, nil
}

func (s *roleService) CreateRole(ctx context.Context, request dto.RoleRequest) (*entity.Role, error) {
	db := s.roleRepository.SingleTransaction()
	newRole := entity.Role{
		Name:      request.Name,
		AuthLevel: request.AuthLevel,
	}

	if err := s.roleRepository.CreateRole(ctx, db, &newRole); err != nil {
		return nil, err
	}

	if err := s.cache.Del("roles:all"); err != nil {
		return nil, err
	}

	return &newRole, nil
}

func (s *roleService) UpdateRole(ctx context.Context, request dto.UpdateRoleRequest) (*entity.Role, error) {
	db := s.roleRepository.SingleTransaction()
	role, err := s.roleRepository.GetRoleByID(ctx, db, request.ID)
	if err != nil {
		return nil, err
	}

	role.Name = request.Name
	role.AuthLevel = request.AuthLevel

	if err := s.roleRepository.UpdateRole(ctx, db, role); err != nil {
		return nil, err
	}

	if err := s.cache.Del("roles:" + role.ID.String()); err != nil {
		return nil, err
	}

	if err := s.cache.Del("roles:all"); err != nil {
		return nil, err
	}

	return role, err
}

func (s *roleService) DeleteRole(ctx context.Context, id string) error {
	db := s.roleRepository.SingleTransaction()
	role, err := s.roleRepository.GetRoleByID(ctx, db, id)
	if err != nil {
		return err
	}

	if err := s.roleRepository.DeleteRole(ctx, db, role); err != nil {
		return err
	}

	if err := s.cache.Del("roles:" + id); err != nil {
		return err
	}

	if err := s.cache.Del("roles:all"); err != nil {
		return err
	}

	return nil
}
