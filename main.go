package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tolhassianipar/kolak_ecom/controllers"
	"github.com/tolhassianipar/kolak_ecom/initializers"
	"github.com/tolhassianipar/kolak_ecom/routes"
)

var (
	server              *gin.Engine
	AuthController      controllers.AuthController
	AuthRouteController routes.AuthRouteController

	UserController      controllers.UserController
	UserRouteController routes.UserRouteController

	PostController      controllers.PostController
	PostRouteController routes.PostRouteController

	ProductController      controllers.ProductController
	ProductRouteController routes.ProductRouteController

	CartController      controllers.CartController
	CartRouteController routes.CartRouteController

	OrderController      controllers.OrderController
	OrderRouteController routes.OrderRouteController

	RoleController      controllers.RoleController
	RoleRouteController routes.RoleRouteController

	UserRoleController      controllers.UserRoleController
	UserRoleRouteController routes.UserRoleRouteController

	FileController      controllers.FileController
	FileRouteController routes.FileRouteController
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)

	AuthController = controllers.NewAuthController(initializers.DB)
	AuthRouteController = routes.NewAuthRouteController(AuthController)

	UserController = controllers.NewUserController(initializers.DB)
	UserRouteController = routes.NewRouteUserController(UserController)

	PostController = controllers.NewPostController(initializers.DB)
	PostRouteController = routes.NewRoutePostController(PostController)

	ProductController = controllers.NewProductController(initializers.DB)
	ProductRouteController = routes.NewRouteProductController(ProductController)

	CartController = controllers.NewCartController(initializers.DB)
	CartRouteController = routes.NewRouteCartController(CartController)

	OrderController = controllers.NewOrderController(initializers.DB)
	OrderRouteController = routes.NewRouteOrderController(OrderController)

	RoleController = controllers.NewRoleController(initializers.DB)
	RoleRouteController = routes.NewRouteRoleController(RoleController)

	UserRoleController = controllers.NewUserRoleController(initializers.DB)
	UserRoleRouteController = routes.NewRouteUserRoleController(UserRoleController)

	FileController = controllers.NewFileController(initializers.DB)
	FileRouteController = routes.NewRouteFileController(FileController)

	server = gin.Default()
}

func main() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:8000", config.ClientOrigin}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	router := server.Group("/api")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		message := "Welcome to Golang with Gorm and Postgres"
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
	})
	router.Static("/static", "./static")

	AuthRouteController.AuthRoute(router)
	UserRouteController.UserRoute(router)
	PostRouteController.PostRoute(router)
	ProductRouteController.ProductRoute(router)
	CartRouteController.CartRoute(router)
	OrderRouteController.OrderRoute(router)
	RoleRouteController.RoleRoute(router)
	UserRoleRouteController.UserRoleRoute(router)
	FileRouteController.FileRoute(router)

	log.Fatal(server.Run(":" + config.ServerPort))
}
