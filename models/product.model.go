package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string `gorm:"not null" json:"name,omitempty"`
	Description string `gorm:"not null" json:"description,omitempty"`
	Price       int    `gorm:"not null" json:"price,omitempty"`
	Qty         int    `gorm:"default:10" json:"qty,omitempty"`
	Image       string `gorm:"not null" json:"image,omitempty"`
	UserID      uint   `gorm:"column:user_id" json:"user_id,omitempty"`
	User        User   `json:"user"`
}

type ProductQueryResponse struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Price       int    `json:"price,omitempty"`
	Image       string `json:"image,omitempty"`
}

type CreateProductResponse struct {
	ID          uint   `json:id,omitempty`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Price       int    `json:"price,omitempty"`
	Image       string `json:"image,omitempty"`
}

type CreateProductRequest struct {
	Name        string `json:"name"  binding:"required"`
	Description string `json:"description" binding:"required"`
	Price       int    `json:"price" binding:"required"`
	Image       string `json:"image" binding:"required"`
	UserID      uint   `json:"user_id,omitempty"`
}

type UpdateProduct struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Price       int    `json:"price,omitempty"`
	Image       string `json:"image,omitempty"`
}
