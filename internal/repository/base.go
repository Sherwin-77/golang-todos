package repository

import "gorm.io/gorm"

type BaseRepository interface {
	WithTransaction(fn func(tx *gorm.DB) error) error
	SingleTransaction() *gorm.DB
	BeginTransaction() *gorm.DB
	Commit(tx *gorm.DB) error
	Rollback(tx *gorm.DB)
}

type baseRepository struct {
	db *gorm.DB
}

func (r *baseRepository) WithTransaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}

func (r *baseRepository) SingleTransaction() *gorm.DB {
	return r.db
}

func (r *baseRepository) BeginTransaction() *gorm.DB {
	return r.db.Begin()
}

func (r *baseRepository) Commit(tx *gorm.DB) error {
	return tx.Commit().Error
}

func (r *baseRepository) Rollback(tx *gorm.DB) {
	tx.Rollback()
}
