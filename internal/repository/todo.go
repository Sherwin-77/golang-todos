package repository

import (
	"context"
	"github.com/sherwin-77/golang-todos/internal/entity"
	"gorm.io/gorm"
)

type TodoRepository interface {
	BaseRepository
	GetTodosByUserID(ctx context.Context, tx *gorm.DB, userID string) ([]entity.Todo, error)
	GetTodosFiltered(ctx context.Context, tx *gorm.DB, limit int, offset int, order interface{}, query interface{}, args ...interface{}) ([]entity.Todo, error)
	GetTodoByID(ctx context.Context, tx *gorm.DB, id string) (*entity.Todo, error)
	CreateTodo(ctx context.Context, tx *gorm.DB, todo *entity.Todo) error
	UpdateTodo(ctx context.Context, tx *gorm.DB, todo *entity.Todo) error
	DeleteTodo(ctx context.Context, tx *gorm.DB, todo *entity.Todo) error
}

type todoRepository struct {
	baseRepository
}

func NewTodoRepository(db *gorm.DB) TodoRepository {
	return &todoRepository{baseRepository{db}}
}

func (r *todoRepository) GetTodosByUserID(ctx context.Context, tx *gorm.DB, userID string) ([]entity.Todo, error) {
	var todos []entity.Todo
	if err := tx.WithContext(ctx).Find(&todos, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return todos, nil
}

func (r *todoRepository) GetTodosFiltered(ctx context.Context, tx *gorm.DB, limit int, offset int, order interface{}, query interface{}, args ...interface{}) ([]entity.Todo, error) {
	var todos []entity.Todo

	if err := tx.WithContext(ctx).Where(query, args...).Limit(limit).Offset(offset).Order(order).Find(&todos).Error; err != nil {
		return nil, err
	}
	return todos, nil
}

func (r *todoRepository) GetTodoByID(ctx context.Context, tx *gorm.DB, id string) (*entity.Todo, error) {
	var todo entity.Todo
	if err := tx.WithContext(ctx).First(&todo, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &todo, nil
}

func (r *todoRepository) CreateTodo(ctx context.Context, tx *gorm.DB, todo *entity.Todo) error {
	if err := tx.WithContext(ctx).Create(todo).Error; err != nil {
		return err
	}
	return nil
}

func (r *todoRepository) UpdateTodo(ctx context.Context, tx *gorm.DB, todo *entity.Todo) error {
	if err := tx.WithContext(ctx).Save(todo).Error; err != nil {
		return err
	}
	return nil
}

func (r *todoRepository) DeleteTodo(ctx context.Context, tx *gorm.DB, todo *entity.Todo) error {
	if err := tx.WithContext(ctx).Delete(todo).Error; err != nil {
		return err
	}
	return nil
}
