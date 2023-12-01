package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `gorm:"type:varchar(255);not null"`
	Email    string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
	Provider string `gorm:"not null"`
	File     File   `json:"photo"`
	FileID   uint   `gorm:"column:file_id" json:"photo_id,omitempty"`
	Verified bool   `gorm:"not null"`
	Post     []Post `gorm:"Foreignkey:UserID;association_foreignkey:ID;"`
	Cart     Cart   `json:"cart"`
	CartID   uint   `gorm:"column:cart_id" json:"cart_id,omitempty"`
	Roles    []Role `gorm:"many2many:user_role;"`
	// UserRoles []UserRole `gorm:"foreignkey:UserId"`
	Dob time.Time `json:"dob"`
}

type UserUpdateRequest struct {
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Verified bool   `json:"verified,omitempty"`
	Roles    []uint `json:"roles,omitempty"`
	Dob      string `json:"dob,omitempty"`
	FileID   uint   `json:"photoId,omitempty"`
}

type SignUpInput struct {
	Name            string `json:"name" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required,min=8"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required"`
	Photo           string `json:"photo" binding:"required"`
	Dob             string `json:"dob"`
	FileID          uint   `json:"photoId,omitempty"`
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
	Dob       string    `json:"dob"`
}
