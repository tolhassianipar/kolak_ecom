package models

import (
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name       string     `gorm:"type:varchar(255);not null"`
	Permission string     `gorm:"type:varchar(255);not null"`
	Users      []User     `gorm:"many2many:user_role;"`
	UserRoles  []UserRole `gorm:"foreignkey:RoleId"`
}
