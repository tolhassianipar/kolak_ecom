package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tolhassianipar/kolak_ecom/controllers"
	"github.com/tolhassianipar/kolak_ecom/middleware"
)

type OrderRouteController struct {
	orderController controllers.OrderController
}

func NewRouteOrderController(orderController controllers.OrderController) OrderRouteController {
	return OrderRouteController{orderController}
}

func (pc *OrderRouteController) OrderRoute(rg *gin.RouterGroup) {

	router := rg.Group("orders")
	router.Use(middleware.DeserializeUser())
	router.GET("/", pc.orderController.FindOrders)
	router.POST("/", pc.orderController.CreateOrder)
	router.DELETE("/:orderId", pc.orderController.DeleteOrder)
}
