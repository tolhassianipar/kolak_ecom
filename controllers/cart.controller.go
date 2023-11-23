package controllers

import (
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/tolhassianipar/kolak_ecom/models"
	"gorm.io/gorm"
)

type CartController struct {
	DB *gorm.DB
}

func NewCartController(DB *gorm.DB) CartController {
	return CartController{DB}
}

func (pd *CartController) CreateCart(ctx *gin.Context) {

	newCart := models.Cart{
		Price: 0,
	}

	result := pd.DB.Create(&newCart)
	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newCart})
}

func (pd *CartController) CreateCartItem(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	productId := ctx.Param("productId")

	var product models.Product
	result := pd.DB.First(&product, "id = ?", productId)

	var payload *models.AddCartItem
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	//validate
	// db.Where("name1 = @name OR name2 = @name", map[string]interface{}{"name": "jinzhu"}).First(&user)
	// SELECT * FROM `users` WHERE name1 = "jinzhu" OR name2 = "jinzhu" ORDER BY `users`.`id` LIMIT 1

	var cartItemToUpdate models.CartItem
	result = pd.DB.Model(models.CartItem{}).Where("cart_id = @cart_id AND product_id = @productId", map[string]interface{}{"cart_id": currentUser.Cart.ID, "productId": productId}).Find(&cartItemToUpdate)

	var cartToUpdate models.Cart
	result = pd.DB.Model(&models.Cart{}).Preload("CartItems").First(&cartToUpdate, "ID = ?", currentUser.Cart.ID)

	var newDataCart models.Cart

	if !(reflect.DeepEqual(cartItemToUpdate, models.CartItem{})) {
		newDataCartItem := models.CartItem{
			Qty:   cartItemToUpdate.Qty + payload.Qty,
			Price: cartItemToUpdate.Price + payload.Qty*product.Price,
		}
		newDataCart = models.Cart{
			Price: cartToUpdate.Price + (payload.Qty * product.Price),
		}
		pd.DB.Model(&cartToUpdate).Updates(newDataCart)
		pd.DB.Model(&cartItemToUpdate).Updates(newDataCartItem)
	} else {
		// Else
		newCartItem := models.CartItem{
			Name:        product.Name,
			Description: product.Description,
			Price:       payload.Qty * product.Price,
			Image:       product.Image,
			Qty:         payload.Qty,
			UserID:      currentUser.ID,
			CartID:      currentUser.CartID, //todo add this to user
			ProductID:   product.ID,
		}

		result = pd.DB.Create(&newCartItem)
		if result.Error != nil {
			ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
			return
		}
		newDataCart = models.Cart{
			Price: cartToUpdate.Price + (payload.Qty * product.Price),
		}
		pd.DB.Model(&cartToUpdate).Updates(newDataCart)
	}

	// get cartItem
	var cartItems []models.CartItem

	result = pd.DB.Model(models.CartItem{}).Find(&cartItems, "cart_id = ?", currentUser.Cart.ID)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error})
		return
	}

	var cartItemResponses []models.CartItemResponse

	for i := 0; i < len(cartItems); i++ {
		cartItemResponses = append(cartItemResponses, models.CartItemResponse{
			Name:        cartItems[i].Name,
			Description: cartItems[i].Description,
			Qty:         cartItems[i].Qty,
			CartID:      cartItems[i].ID,
			Price:       cartItems[i].Price,
		})
	}

	var response models.CartResponse
	response.CartItems = cartItemResponses
	response.Price = cartToUpdate.Price

	result = pd.DB.Model(&models.Cart{}).Preload("CartItems").First(&cartToUpdate, "ID = ?", currentUser.Cart.ID)
	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": response})

}

func (pd *CartController) DeleteCartItem(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	cartItemId := ctx.Param("cartItemId")
	var newDataCart models.Cart
	var cartItemToDelete models.CartItem
	result := pd.DB.Model(models.CartItem{}).First(&cartItemToDelete, "id", cartItemId)

	var cartToUpdate models.Cart
	result = pd.DB.Model(&models.Cart{}).Preload("CartItems").First(&cartToUpdate, "ID = ?", currentUser.Cart.ID)

	newDataCart = models.Cart{
		Price: cartToUpdate.Price - cartItemToDelete.Price,
	}

	pd.DB.Model(&cartToUpdate).Updates(newDataCart)

	result = pd.DB.Delete(&models.CartItem{}, "id = ?", cartItemId)

	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No cartitem with that title exists"})
		return
	}

	// get latest cartItem
	var cartItems []models.CartItem

	result = pd.DB.Model(models.CartItem{}).Find(&cartItems, "cart_id = ?", currentUser.Cart.ID)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error})
		return
	}

	// get cart
	var cartItemResponses []models.CartItemResponse

	for i := 0; i < len(cartItems); i++ {
		cartItemResponses = append(cartItemResponses, models.CartItemResponse{
			ID:          cartItems[i].ID,
			Name:        cartItems[i].Name,
			Description: cartItems[i].Description,
			Qty:         cartItems[i].Qty,
			CartID:      cartItems[i].ID,
			Price:       cartItems[i].Price,
		})
	}

	var response models.CartResponse
	response.CartItems = cartItemResponses
	response.Price = cartToUpdate.Price

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": response})
}

func (pd *CartController) UpdateCartItem(ctx *gin.Context) {
	cartItemId := ctx.Param("cartItemId")
	currentUser := ctx.MustGet("currentUser").(models.User)
	var payload *models.UpdateCartItem
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	var updatedCartItem models.CartItem
	result := pd.DB.Preload("Product").First(&updatedCartItem, "id = ?", cartItemId)
	if result.Error != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "No cartItem with that title exists"})
		return
	}

	priceBefore := updatedCartItem.Price
	productPrice := updatedCartItem.Product.Price

	cartItemToUpdate := models.CartItem{
		Qty:   payload.Qty,
		Price: payload.Qty * productPrice,
	}

	pd.DB.Model(&updatedCartItem).Updates(cartItemToUpdate)

	var cartToUpdate models.Cart
	result = pd.DB.Model(&models.Cart{}).First(&cartToUpdate, "ID = ?", currentUser.Cart.ID)

	newDataCart := models.Cart{
		Price: cartToUpdate.Price - priceBefore + (payload.Qty * productPrice),
	}
	pd.DB.Model(&cartToUpdate).Updates(newDataCart)

	// get latest cartItem
	var cartItems []models.CartItem

	result = pd.DB.Model(models.CartItem{}).Find(&cartItems, "cart_id = ?", currentUser.Cart.ID)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error})
		return
	}

	//get cart
	var cartItemResponses []models.CartItemResponse

	for i := 0; i < len(cartItems); i++ {
		cartItemResponses = append(cartItemResponses, models.CartItemResponse{
			ID:          cartItems[i].ID,
			Name:        cartItems[i].Name,
			Description: cartItems[i].Description,
			Qty:         cartItems[i].Qty,
			CartID:      cartItems[i].ID,
			Price:       cartItems[i].Price,
		})
	}

	var response models.CartResponse
	response.CartItems = cartItemResponses
	response.Price = cartToUpdate.Price

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": response})
}

func (pd *CartController) FindCarts(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)
	// userId := ctx.Param("userId")
	// var page = ctx.DefaultQuery("page", "1")
	// var limit = ctx.DefaultQuery("limit", "10")

	// intPage, _ := strconv.Atoi(page)
	// intLimit, _ := strconv.Atoi(limit)
	// offset := (intPage - 1) * intLimit
	var cartResponse models.Cart
	result := pd.DB.Model(&models.Cart{}).First(&cartResponse, "ID = ?", currentUser.Cart.ID)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error})
		return
	}

	var cartItems []models.CartItem
	result = pd.DB.Model(models.CartItem{}).Find(&cartItems, "cart_id = ?", currentUser.Cart.ID)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error})
		return
	}

	var cartItemResponses []models.CartItemResponse

	for i := 0; i < len(cartItems); i++ {
		cartItemResponses = append(cartItemResponses, models.CartItemResponse{
			Name:        cartItems[i].Name,
			Description: cartItems[i].Description,
			Qty:         cartItems[i].Qty,
			CartID:      cartItems[i].ID,
		})
	}

	var response models.CartResponse
	response.CartItems = cartItemResponses
	response.Price = cartResponse.Price

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": response})
}
