package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tolhassianipar/kolak_ecom/initializers"
	"github.com/tolhassianipar/kolak_ecom/models"
	"github.com/tolhassianipar/kolak_ecom/utils"
	"golang.org/x/exp/slices"
)

func DeserializeUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var access_token string
		cookie, err := ctx.Cookie("access_token")

		authorizationHeader := ctx.Request.Header.Get("Authorization")
		fields := strings.Fields(authorizationHeader)

		if len(fields) != 0 && fields[0] == "Bearer" {
			access_token = fields[1]
		} else if err == nil {
			access_token = cookie
		}

		if access_token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You are not logged in"})
			return
		}

		config, _ := initializers.LoadConfig(".")
		sub, err := utils.ValidateToken(access_token, config.AccessTokenPublicKey)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": err.Error()})
			return
		}

		var user models.User
		result := initializers.DB.Preload("Cart").First(&user, "id = ?", fmt.Sprint(sub))
		if result.Error != nil {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user belonging to this token no logger exists"})
			return
		}

		ctx.Set("currentUser", user)
		ctx.Next()
	}
}

func GetUser(ctx *gin.Context, userId string) models.User {
	var user models.User
	result := initializers.DB.Preload("Cart").Preload("Roles").First(&user, "id = ?", userId)
	// println(user.Roles[0].Permission)
	if result.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "no user is found"})
		return models.User{}
	}
	return user
}

func SetUserCtx(ctx *gin.Context, userId string) {

	var user models.User
	result := initializers.DB.Preload("Cart").First(&user, "id = ?", userId)
	if result.Error != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user belonging to this token no logger exists"})
		return
	}

	ctx.Set("currentUser", user)
}

func IsAuthorized(permission string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var access_token string
		cookie, err := ctx.Cookie("access_token")

		authorizationHeader := ctx.Request.Header.Get("Authorization")
		fields := strings.Fields(authorizationHeader)

		if len(fields) != 0 && fields[0] == "Bearer" {
			access_token = fields[1]
		} else if err == nil {
			access_token = cookie
		}

		if access_token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You are not logged in"})
			return
		}

		config, _ := initializers.LoadConfig(".")
		sub, err := utils.ValidateToken(access_token, config.AccessTokenPublicKey)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": err.Error()})
			return
		}

		SetUserCtx(ctx, fmt.Sprint(sub))
		currentUser := GetUser(ctx, fmt.Sprint(sub))

		userPermissions := []string{}
		for i := 0; i < len(currentUser.Roles); i++ {
			userPermissions = append(userPermissions, currentUser.Roles[i].Permission)
		}
		println("contains result")
		println(!slices.Contains(userPermissions, permission))

		if !slices.Contains(userPermissions, permission) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "error", "message": fmt.Sprintf("You dont have %s permission", permission)})
			return
		}

		ctx.Next()
	}
}

func IsAuthenticated() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var access_token string
		cookie, err := ctx.Cookie("access_token")

		authorizationHeader := ctx.Request.Header.Get("Authorization")
		fields := strings.Fields(authorizationHeader)

		if len(fields) != 0 && fields[0] == "Bearer" {
			access_token = fields[1]
		} else if err == nil {
			access_token = cookie
		}

		if access_token == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You are not logged in"})
			return
		}

		config, _ := initializers.LoadConfig(".")
		sub, err := utils.ValidateToken(access_token, config.AccessTokenPublicKey)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		SetUserCtx(ctx, fmt.Sprint(sub))
		ctx.Next()
	}

}
