package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tolhassianipar/kolak_ecom/models"
	"gorm.io/gorm"
)

type RoleController struct {
	DB *gorm.DB
}

func NewRoleController(DB *gorm.DB) RoleController {
	return RoleController{DB}
}

func (pc *RoleController) CreateRole(ctx *gin.Context) {
	var payload *models.CreateRoleRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	newRole := models.Role{
		Name:       payload.Name,
		Permission: payload.Permission,
	}

	result := pc.DB.Create(&newRole)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Role with that title already exists"})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	roleResponse := models.RoleResponse{
		Id:         newRole.ID,
		Name:       newRole.Name,
		Permission: newRole.Permission,
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": roleResponse})
}

func (pc *RoleController) UpdateRole(ctx *gin.Context) {
	roleId := ctx.Param("roleId")

	var payload *models.CreateRoleRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	var roletoUpdate models.Role
	result := pc.DB.First(&roletoUpdate, "id = ?", roleId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No post with that title exists"})
		return
	}

	newRoleData := models.Role{
		Name:       payload.Name,
		Permission: payload.Permission,
	}

	pc.DB.Model(&roletoUpdate).Updates(newRoleData)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": roletoUpdate})
}

func (pc *RoleController) FindRoles(ctx *gin.Context) {
	// var page = ctx.DefaultQuery("page", "1")
	// var limit = ctx.DefaultQuery("limit", "10")

	// intPage, _ := strconv.Atoi(page)
	// intLimit, _ := strconv.Atoi(limit)
	// offset := (intPage - 1) * intLimit

	var roles []models.RoleResponse
	results := pc.DB.Model(models.Role{}).Find(&roles)
	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(roles), "data": roles})
}
