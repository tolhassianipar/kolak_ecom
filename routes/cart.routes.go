package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tolhassianipar/kolak_ecom/controllers"
	"github.com/tolhassianipar/kolak_ecom/middleware"
)

type CartRouteController struct {
	cartController controllers.CartController
}

func NewRouteCartController(cartController controllers.CartController) CartRouteController {
	return CartRouteController{cartController}
}

func (pc *CartRouteController) CartRoute(rg *gin.RouterGroup) {

	router := rg.Group("carts")
	router.Use(middleware.DeserializeUser())
	router.POST("/items/:productId", pc.cartController.CreateCartItem)
	router.PUT("/items/:cartItemId", pc.cartController.UpdateCartItem)
	router.DELETE("/items/:cartItemId", pc.cartController.DeleteCartItem)
	router.GET("/", pc.cartController.FindCarts)
}
