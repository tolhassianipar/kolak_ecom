package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tolhassianipar/kolak_ecom/controllers"
	"github.com/tolhassianipar/kolak_ecom/middleware"
)

type ProductRouteController struct {
	productController controllers.ProductController
}

func NewRouteProductController(productController controllers.ProductController) ProductRouteController {
	return ProductRouteController{productController}
}

func (pc *ProductRouteController) ProductRoute(rg *gin.RouterGroup) {

	router := rg.Group("products")
	router.Use(middleware.DeserializeUser())
	router.POST("/", pc.productController.CreateProduct)
	router.GET("/", pc.productController.FindProducts)
	router.PUT("/:productId", pc.productController.UpdateProduct)
	router.GET("/:productId", pc.productController.FindProductById)
	router.GET("user/:userId", pc.productController.FindProductsByUser)
	router.DELETE("/:productId", pc.productController.DeleteProduct)
}
