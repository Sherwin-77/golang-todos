package entity

type User struct {
	BaseEntity
	Username string `json:"username" gorm:"type:varchar(255);not null"`
	Email    string `json:"email" gorm:"type:varchar(255);not null;uniqueIndex"`
	Password string `json:"-"`

	Roles []*Role `json:"roles,omitempty" gorm:"many2many:role_users;"`
}
