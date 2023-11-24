package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tolhassianipar/kolak_ecom/models"
	"gorm.io/gorm"
)

type OrderController struct {
	DB *gorm.DB
}

func NewOrderController(DB *gorm.DB) OrderController {
	return OrderController{DB}
}

func (pd *OrderController) CreateOrder(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	var payload *models.CreateOrderPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	// create order
	newOrder := models.Order{
		UserID: currentUser.ID,
		Price:  0,
	}

	result := pd.DB.Create(&newOrder)
	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}
	// get cartitem
	cartItemIds := payload.CartItemIDs
	var cartItems []models.CartItem
	result = pd.DB.Model(models.CartItem{}).Preload("Product").Where("id IN ? ", cartItemIds).Find(&cartItems)

	// create orderitem
	var newOrderItem models.OrderItem
	var newOrderItems []models.OrderItem
	totalPrice := 0
	for i := 0; i < len(cartItems); i++ {
		newOrderItem = models.OrderItem{
			Name:        cartItems[i].Name,
			Description: cartItems[i].Description,
			Price:       cartItems[i].Price,
			Image:       cartItems[i].Image,
			Qty:         cartItems[i].Qty,
			UserID:      currentUser.ID,
			OrderID:     newOrder.ID,
			ProductID:   cartItems[i].ProductID,
			BasedPrice:  cartItems[i].Product.Price,
		}
		result = pd.DB.Create(&newOrderItem)
		if result.Error != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
			return
		}
		totalPrice += newOrderItem.Price
		newOrderItems = append(newOrderItems, newOrderItem)
	}
	newOrderData := models.Order{
		Price:      totalPrice,
		OrderItems: newOrderItems,
	}
	result = pd.DB.Model(&newOrder).Updates(newOrderData)
	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}
	//update cart data
	var cartToUpdate models.Cart
	result = pd.DB.First(&cartToUpdate, "id = ?", cartItems[0].CartID)
	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}
	println(cartToUpdate.Price, totalPrice)
	updatedPrice := cartToUpdate.Price - totalPrice
	println("updatedPrice")
	println(updatedPrice)
	// newCartData := models.Cart{
	// 	Price: updatedPrice,
	// }
	cartToUpdate.Price = updatedPrice
	result = pd.DB.Save(cartToUpdate)
	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}
	result = pd.DB.Delete(&cartItems)
	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	var orderItemResponse []models.OrderItemResponse

	for i := 0; i < len(newOrderItems); i++ {
		orderItemResponse = append(orderItemResponse, models.OrderItemResponse{
			ID:          newOrderItems[i].ID,
			Name:        newOrderItems[i].Name,
			Description: newOrderItems[i].Description,
			Qty:         newOrderItems[i].Qty,
			OrderID:     newOrder.ID,
			Price:       newOrderItems[i].Price,
			BasedPrice:  newOrderItems[i].BasedPrice,
		})
	}

	var response models.OrderResponse
	response.OrderItems = orderItemResponse
	response.Price = newOrder.Price
	response.ID = newOrder.ID
	response.TotalItems = len(orderItemResponse)

	// return orderandorderitem response

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": response})
}

func (pd *OrderController) DeleteOrder(ctx *gin.Context) {
	orderId := ctx.Param("orderId")

	result := pd.DB.Delete(&models.Order{}, "id = ?", orderId)

	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No cart with that title exists"})
		return
	}

	pd.FindOrders(ctx)
}

func (pd *OrderController) FindOrders(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)

	// get all orders
	var orders []models.Order
	result := pd.DB.Model(&models.Order{}).Preload("OrderItems").Find(&orders, "user_id = ?", currentUser.ID)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error})
		return
	}

	var orderResponses []models.OrderResponse
	var orderItemResponses []models.OrderItemResponse
	// for every orders set the orderresponse and orderitemresponse

	for i := 0; i < len(orders); i++ {
		orderResponses = append(orderResponses, models.OrderResponse{
			ID:    orders[i].ID,
			Price: orders[i].Price,
		})
		for j := 0; j < len(orders[i].OrderItems); j++ {
			orderItemResponses = append(orderItemResponses, models.OrderItemResponse{
				ID:          orders[i].OrderItems[j].ID,
				Name:        orders[i].OrderItems[j].Name,
				Description: orders[i].OrderItems[j].Description,
				Qty:         orders[i].OrderItems[j].Qty,
				OrderID:     orders[i].ID,
				Price:       orders[i].OrderItems[j].Price,
				BasedPrice:  orders[i].OrderItems[j].BasedPrice,
			})
		}
		orderResponses[i].OrderItems = orderItemResponses
		orderResponses[i].TotalItems = len(orderItemResponses)
		orderItemResponses = nil
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": orderResponses})
}
