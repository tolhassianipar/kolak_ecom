package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/tolhassianipar/kolak_ecom/controllers"
	"github.com/tolhassianipar/kolak_ecom/middleware"
)

type FileRouteController struct {
	fileController controllers.FileController
}

func NewRouteFileController(fileController controllers.FileController) FileRouteController {
	return FileRouteController{fileController}
}

func (pc *FileRouteController) FileRoute(rg *gin.RouterGroup) {

	router := rg.Group("files")
	router.Use(middleware.IsAuthenticated())
	router.POST("/", pc.fileController.CreateFile)
}
