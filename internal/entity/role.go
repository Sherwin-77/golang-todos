package entity

type Role struct {
	BaseEntity
	Name      string `json:"name" gorm:"type:varchar(255);not null"`
	AuthLevel int    `json:"auth_level" gorm:"type:integer;not null"`
}
