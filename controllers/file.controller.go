package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tolhassianipar/kolak_ecom/models"
	"gorm.io/gorm"
)

type FileController struct {
	DB *gorm.DB
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func generateRandomString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func NewFileController(DB *gorm.DB) FileController {
	return FileController{DB}
}

func (pc *FileController) CreateFile(ctx *gin.Context) {
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
		return
	}
	files := form.File["images"]
	print(files)
	var images = make([]models.File, len(files))
	for index, file := range files {
		fileName := generateRandomString(10) + ".png"

		dirPath := filepath.Join(".", "static", "images")
		filePath := filepath.Join(dirPath, fileName)
		// Create directory if does not exist
		// can be replaced with clod storage
		if _, err = os.Stat(dirPath); os.IsNotExist(err) {
			err = os.MkdirAll(dirPath, os.ModeDir)
			if err != nil {
				ctx.JSON(http.StatusBadRequest, err.Error())
				return
			}
		}
		// Create file that will hold the image
		outputFile, err := os.Create(filePath)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}
		defer outputFile.Close()

		// Open the temporary file that contains the uploaded image
		inputFile, err := file.Open()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
		}
		defer inputFile.Close()

		// Copy the temporary image to the permanent location outputFile
		_, err = io.Copy(outputFile, inputFile)
		if err != nil {
			ctx.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}

		fileSize := (uint)(file.Size)
		images[index] = models.File{Filename: fileName, FilePath: string(filepath.Separator) + "api/" + filePath, FileSize: fileSize, OriginalName: file.Filename}
	}

	// database := infrastructure.GetDb()
	// category := models.Category{Name: name, Description: description, Images: categoryImages}
	result := pc.DB.Create(&images)

	if result.Error != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": images, "length": len(images)})
	// TODO: Why it is performing a SELECT SQL Query per image?
	// Even worse, it is selecting category_id, why??
	// SELECT "tag_id", "product_id" FROM "file_uploads"  WHERE (id = insertedFileUploadId)
	// err = database.Create(&category).Error

	// var payload *models.CreateRoleRequest

	// if err := ctx.ShouldBindJSON(&payload); err != nil {
	// 	ctx.JSON(http.StatusBadRequest, err.Error())
	// 	return
	// }

	// newRole := models.Role{
	// 	Name:       payload.Name,
	// 	Permission: payload.Permission,
	// }

	// result := pc.DB.Create(&newRole)
	// if result.Error != nil {
	// 	if strings.Contains(result.Error.Error(), "duplicate key") {
	// 		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "Role with that title already exists"})
	// 		return
	// 	}
	// 	ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
	// 	return
	// }

	// ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newRole})
}

func (pc *FileController) UpdateRole(ctx *gin.Context) {
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

func (pc *FileController) FindRoles(ctx *gin.Context) {
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
