package models

import (
	"gorm.io/gorm"
)

type UserRole struct {
	gorm.Model
	User   User `gorm:"association_foreignkey:UserId"`
	UserId uint
	Role   User `gorm:"association_foreignkey:RoleId"`
	RoleId uint
}

type UserRoleRequest struct {
	UserId uint `json:"userId"`
	RoleId uint `json:"roleId"`
}

func (UserRole) TableName() string {
	return "user_role"
}

func Any(roles []Role, f func(Role) bool) bool { // f func as a validator
	for _, role := range roles {
		if f(role) {
			return true
		}
	}
	return false
}
