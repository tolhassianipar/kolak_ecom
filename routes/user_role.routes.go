package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tolhassianipar/kolak_ecom/controllers"
)

type UserRoleRouteController struct {
	userRoleController controllers.UserRoleController
}

func NewRouteUserRoleController(userRoleController controllers.UserRoleController) UserRoleRouteController {
	return UserRoleRouteController{userRoleController}
}

func (pc *UserRoleRouteController) UserRoleRoute(rg *gin.RouterGroup) {

	router := rg.Group("userRoles")
	router.POST("/", pc.userRoleController.CreateUserRole)
}
