package models

import (
	"time"
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title     string    `gorm:"uniqueIndex;not null" json:"title,omitempty"`
	Content   string    `gorm:"not null" json:"content,omitempty"`
	Image     string    `gorm:"not null" json:"image,omitempty"`
	UserID   uint      `gorm:"column:user_id" json:"user_id,omitempty"`
	User	 User		`json:"user"`
}

type PostQueryResponse struct {
	Title     string    `json:"title,omitempty"`
	Content   string    `json:"content,omitempty"`
	Image     string    `json:"image,omitempty"`
	UserID    uint      `json:"user_id,omitempty"`
}

type CreatePostRequest struct {
	Title     string    `json:"title"  binding:"required"`
	Content   string    `json:"content" binding:"required"`
	Image     string    `json:"image" binding:"required"`
	UserID   uint      `json:"user_id,omitempty"`
}

type UpdatePost struct {
	Title     string    `json:"title,omitempty"`
	Content   string    `json:"content,omitempty"`
	Image     string    `json:"image,omitempty"`
	UserID   uint      `json:"user_id,omitempty"`
	CreateAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}
