package models

import (
	"gorm.io/gorm"
)

type Cart struct {
	gorm.Model
	CartItems []CartItem `gorm:"foreignKey:CartID;association_foreignkey:ID;"`
	Price     int        `gorm:"not null" json:"price,omitempty"`
}

type CartResponse struct {
	CartItems []CartItemResponse `json:"cartItems,omitempty"`
	Price     int                `json:"price,omitempty"`
}

// func (order *Order) GetOrderStatusAsString() string {
// 	switch order.OrderStatus {
// 	case 0:
// 		return "processed"
// 	case 1:
// 		return "delivered"
// 	case 2:
// 		return "shipped"
// 	default:
// 		return "unknown"
// 	}
// }
