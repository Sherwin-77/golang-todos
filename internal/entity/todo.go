package entity

import "github.com/google/uuid"

type Todo struct {
	BaseEntity
	Title       string    `json:"title" gorm:"type:varchar(255);not null"`
	Description string    `json:"description" gorm:"type:varchar(255);"`
	IsCompleted bool      `json:"is_completed" gorm:"type:bool;not null;default:false"`
	UserID      uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`

	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}
