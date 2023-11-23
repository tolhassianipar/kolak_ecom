package controllers

import (
	"net/http"
	"strconv"


	"github.com/gin-gonic/gin"
	"github.com/tolhassianipar/kolak_ecom/models"
	"gorm.io/gorm"
)

type ProductCartController struct {
	DB *gorm.DB
}

func NewProductCartController(DB *gorm.DB) ProductCartController {
	return ProductCartController{DB}
}

func (pdc *ProductCartController) CreateProduct(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	var payload *models.CreateProductRequest

	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	newProduct := models.Product{
		Name:     payload.Name,
		Description:   payload.Description,
		Price:     payload.Price,
		Image:     payload.Image,
		UserID:   currentUser.ID,
	}

	result := pdc.DB.Create(&newProduct)
	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}
	newProdResponse := models.CreateProductResponse{
		Name: newProduct.Name,
		Description: newProduct.Description,
		Price: newProduct.Price,
		Image: newProduct.Image,
	}
	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newProdResponse})
}

func (pdc *ProductCartController) UpdateProduct(ctx *gin.Context) {
	productId := ctx.Param("productId")

	var payload *models.UpdateProduct
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	var updatedProduct models.Product
	result := pdc.DB.First(&updatedProduct, "id = ?", productId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No product with that title exists"})
		return
	}
	
	productToUpdate := models.Product{
		Name:     payload.Name,
		Description:   payload.Description,
		Price:     payload.Price,
		Image:     payload.Image,
	}

	pdc.DB.Model(&updatedProduct).Updates(productToUpdate)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": updatedProduct})
}

func (pdc *ProductCartController) FindProductById(ctx *gin.Context) {
	productId := ctx.Param("productId")

	var product models.Product
	result := pdc.DB.Preload("User").First(&product, "id = ?", productId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No product with that title exists"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": product})
}

func (pdc *ProductCartController) FindProducts(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var products []models.Product
	results := pdc.DB.Preload("User").Limit(intLimit).Offset(offset).Find(&products)
	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(products), "data": products})
}

func (pdc *ProductCartController) FindProductsByUser(ctx *gin.Context) {
	userId := ctx.Param("userId")
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	intPage, _ := strconv.Atoi(page)
	intLimit, _ := strconv.Atoi(limit)
	offset := (intPage - 1) * intLimit

	var productresponses []models.ProductQueryResponse
    results := pdc.DB.Model(&models.Product{}).Limit(intLimit).Offset(offset).Find(&productresponses, "user_id = ?", userId)

	// var product []models.product
	// results := pdc.DB.Preload("User").Limit(intLimit).Offset(offset).Find(&product, "user_id = ?", userId)
	if results.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": results.Error})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(productresponses), "data": productresponses})
}

func (pdc *ProductCartController) DeleteProduct(ctx *gin.Context) {
	productId := ctx.Param("productId")

	result := pdc.DB.Delete(&models.Product{}, "id = ?", productId)

	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No product with that title exists"})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
