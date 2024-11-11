package repository

import (
	"context"

	"github.com/sherwin-77/go-echo-template/internal/entity"
	"gorm.io/gorm"
)

type RoleRepository interface {
	BaseRepository
	GetRoles(ctx context.Context, tx *gorm.DB) ([]entity.Role, error)
	GetRoleByID(ctx context.Context, tx *gorm.DB, id string) (*entity.Role, error)
	CreateRole(ctx context.Context, tx *gorm.DB, role *entity.Role) error
	UpdateRole(ctx context.Context, tx *gorm.DB, role *entity.Role) error
	DeleteRole(ctx context.Context, tx *gorm.DB, role *entity.Role) error
}

type roleRepository struct {
	baseRepository
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{baseRepository{db}}
}

func (r *roleRepository) GetRoles(ctx context.Context, tx *gorm.DB) ([]entity.Role, error) {
	var roles []entity.Role

	if err := tx.WithContext(ctx).Find(&roles).Error; err != nil {
		return nil, err
	}

	return roles, nil
}

func (r *roleRepository) GetRoleByID(ctx context.Context, tx *gorm.DB, id string) (*entity.Role, error) {
	var role entity.Role

	if err := tx.WithContext(ctx).Where("id = ?", id).First(&role).Error; err != nil {
		return nil, err
	}

	return &role, nil
}

func (r *roleRepository) CreateRole(ctx context.Context, tx *gorm.DB, role *entity.Role) error {
	if err := tx.WithContext(ctx).Create(role).Error; err != nil {
		return err
	}

	return nil
}

func (r *roleRepository) UpdateRole(ctx context.Context, tx *gorm.DB, role *entity.Role) error {
	if err := tx.WithContext(ctx).Save(role).Error; err != nil {
		return err
	}

	return nil
}

func (r *roleRepository) DeleteRole(ctx context.Context, tx *gorm.DB, role *entity.Role) error {
	if err := tx.WithContext(ctx).Delete(role).Error; err != nil {
		return err
	}

	return nil
}
