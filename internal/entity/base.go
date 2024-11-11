package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseEntity struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (b *BaseEntity) BeforeCreate(tx *gorm.DB) error {
	var err error
	if b.ID == uuid.Nil {
		b.ID, err = uuid.NewV7()
	}

	return err
}
