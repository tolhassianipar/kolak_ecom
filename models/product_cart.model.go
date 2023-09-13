package models

import (
	"gorm.io/gorm"
)

type ProductCart struct {
	gorm.Model
	Qty	 int		`gorm:"not null" json:"qty,omitempty"`
	ProductID   uint      `gorm:"column:product_id" json:"product_id,omitempty"`
	Product	 Producy		`json:"product"`
}

type ProductQueryResponse struct {
	Name     string    `json:"name,omitempty"`
	Description   string    `json:"description,omitempty"`
	Price	 int		`json:"price,omitempty"`
	Image    string    `json:"image,omitempty"`
}

type CreateProductResponse struct {
	Name     string    `json:"name,omitempty"`
	Description   string    `json:"description,omitempty"`
	Price	 int		`json:"price,omitempty"`
	Image    string    `json:"image,omitempty"`
}

type CreateProductRequest struct {
	Name     string    `json:"name"  binding:"required"`
	Description   string    `json:"description" binding:"required"`
	Price	  int		`json:"price" binding:"required"`
	Image     string    `json:"image" binding:"required"`
	UserID   uint      `json:"user_id,omitempty"`
}

type UpdateProduct struct {
	Name     string    `json:"name,omitempty"`
	Description   string    `json:"description,omitempty"`
	Price	 int		`json:"price,omitempty"`
	Image    string    `json:"image,omitempty"`
}
