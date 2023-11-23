package models

import (
	"gorm.io/gorm"
)

type CartItem struct {
	gorm.Model
	Name        string  `gorm:"not null" json:"name,omitempty"`
	Description string  `gorm:"not null" json:"description,omitempty"`
	Price       int     `gorm:"not null" json:"price,omitempty"`
	Image       string  `gorm:"not null" json:"image,omitempty"`
	Qty         int     `gorm:"not null" json:"qty,omitempty"`
	UserID      uint    `gorm:"column:user_id" json:"user_id,omitempty"`
	CartID      uint    `gorm:"column:cart_id" json:"cart_id,omitempty"`
	ProductID   uint    `gorm:"column:product_id" json:"product_id,omitempty"`
	User        User    `json:"user"`
	Product     Product `json:"product"`
}

type CartItemResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Qty         int    `json:"qty,omitempty"`
	CartID      uint   `json:"cart_id,omitempty"`
	Price       int    `json:"price,omitempty"`
}

type AddCartItem struct {
	Qty int `json:"qty"`
}

type UpdateCartItem struct {
	Qty int `json:"qty,omitempty"`
}

// type ProductQueryResponse struct {
// 	Name     string    `json:"name,omitempty"`
// 	Description   string    `json:"description,omitempty"`
// 	Price	 int		`json:"price,omitempty"`
// 	Image    string    `json:"image,omitempty"`
// }

// type CreateProductResponse struct {
// 	Name     string    `json:"name,omitempty"`
// 	Description   string    `json:"description,omitempty"`
// 	Price	 int		`json:"price,omitempty"`
// 	Image    string    `json:"image,omitempty"`
// }

type CreateProductCartRequest struct {
	Qty       int  `json:"qty,omitempty"`
	CartID    uint `json:"cart_id,omitempty"`
	ProductID uint `json:"product_id,omitempty"`
}

// type UpdateProduct struct {
// 	Qty	 int		`json:"qty,omitempty"`
// }
