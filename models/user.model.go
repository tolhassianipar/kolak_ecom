package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name      string     `gorm:"type:varchar(255);not null"`
	Email     string     `gorm:"uniqueIndex;not null"`
	Password  string     `gorm:"not null"`
	Provider  string     `gorm:"not null"`
	Photo     string     `gorm:"not null"`
	Verified  bool       `gorm:"not null"`
	Post      []Post     `gorm:"Foreignkey:UserID;association_foreignkey:ID;"`
	Cart      Cart       `json:"cart"`
	CartID    uint       `gorm:"column:cart_id" json:"cart_id,omitempty"`
	Roles     []Role     `gorm:"many2many:user_role;"`
	UserRoles []UserRole `gorm:"foreignkey:UserId"`
}

type SignUpInput struct {
	Name            string `json:"name" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required,min=8"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required"`
	Photo           string `json:"photo" binding:"required"`
}

type SignInInput struct {
	Email    string `json:"email"  binding:"required"`
	Password string `json:"password"  binding:"required"`
}

type UserResponse struct {
	ID        uint      `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Email     string    `json:"email,omitempty"`
	Role      string    `json:"role,omitempty"`
	Photo     string    `json:"photo,omitempty"`
	Provider  string    `json:"provider"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Cart      Cart      `json:"cart"`
}
