package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tolhassianipar/kolak_ecom/controllers"
	"github.com/tolhassianipar/kolak_ecom/middleware"
)

type RoleRouteController struct {
	roleController controllers.RoleController
}

func NewRouteRoleController(roleController controllers.RoleController) RoleRouteController {
	return RoleRouteController{roleController}
}

func (pc *RoleRouteController) RoleRoute(rg *gin.RouterGroup) {

	router := rg.Group("roles")
	router.Use(middleware.IsAuthenticated())
	router.POST("/", pc.roleController.CreateRole)
	router.GET("/", pc.roleController.FindRoles)
	router.PUT("/:roleId", pc.roleController.UpdateRole)
}
