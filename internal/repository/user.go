package repository

import (
	"context"

	"github.com/sherwin-77/go-echo-template/internal/entity"
	"gorm.io/gorm"
)

type UserRepository interface {
	BaseRepository
	GetUsers(ctx context.Context, tx *gorm.DB) ([]entity.User, error)
	GetUserByID(ctx context.Context, tx *gorm.DB, id string) (*entity.User, error)
	GetUserByEmail(ctx context.Context, tx *gorm.DB, email string) (*entity.User, error)
	CreateUser(ctx context.Context, tx *gorm.DB, user *entity.User) error
	UpdateUser(ctx context.Context, tx *gorm.DB, user *entity.User) error
	DeleteUser(ctx context.Context, tx *gorm.DB, user *entity.User) error
	AddRoles(ctx context.Context, tx *gorm.DB, user *entity.User, roles []*entity.Role) error
	RemoveRoles(ctx context.Context, tx *gorm.DB, user *entity.User, roles []*entity.Role) error
}
type userRepository struct {
	baseRepository
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{baseRepository{db}}
}

func (r *userRepository) GetUsers(ctx context.Context, tx *gorm.DB) ([]entity.User, error) {
	var users []entity.User

	if err := tx.WithContext(ctx).Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, tx *gorm.DB, id string) (*entity.User, error) {
	var user entity.User

	if err := tx.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, tx *gorm.DB, email string) (*entity.User, error) {
	var user entity.User

	if err := tx.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) CreateUser(ctx context.Context, tx *gorm.DB, user *entity.User) error {
	if err := tx.WithContext(ctx).Create(user).Error; err != nil {
		return err
	}

	return nil
}

func (r *userRepository) UpdateUser(ctx context.Context, tx *gorm.DB, user *entity.User) error {
	if err := tx.WithContext(ctx).Save(user).Error; err != nil {
		return err
	}

	return nil
}

func (r *userRepository) DeleteUser(ctx context.Context, tx *gorm.DB, user *entity.User) error {
	if err := tx.WithContext(ctx).Delete(user).Error; err != nil {
		return err
	}

	return nil
}

func (r *userRepository) AddRoles(ctx context.Context, tx *gorm.DB, user *entity.User, roles []*entity.Role) error {
	if err := tx.WithContext(ctx).Model(user).Association("Roles").Append(roles); err != nil {
		return err
	}

	return nil
}

func (r *userRepository) RemoveRoles(ctx context.Context, tx *gorm.DB, user *entity.User, roles []*entity.Role) error {
	if err := tx.WithContext(ctx).Model(user).Association("Roles").Delete(roles); err != nil {
		return err
	}

	return nil
}
