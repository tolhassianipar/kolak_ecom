package models

import (
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	OrderItems []OrderItem `gorm:"foreignKey:OrderID;association_foreignkey:ID;"`
	Price      int         `gorm:"not null" json:"price,omitempty"`
	UserID     uint        `gorm:"column:user_id" json:"user_id,omitempty"`
	User       User        `json:"user,omitempty"`
}

type OrderResponse struct {
	ID         uint                `json:"id"`
	OrderItems []OrderItemResponse `json:"orderItems"`
	TotalItems int                 `json:"total_item"`
	Price      int                 `json:"price"`
}

type CreateOrderPayload struct {
	CartItemIDs []uint `json:"cartItemIds,omitempty"`
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
