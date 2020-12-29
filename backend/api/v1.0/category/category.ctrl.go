package category

import (
	"github.com/Ismail14098/agyn_test_rest/database/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func create(ctx *gin.Context){
	db := ctx.MustGet("db").(*gorm.DB)

	type RequestBody struct {
		Name string `json:"name" binding:"required,correctCategoryName"`
	}

	var body RequestBody
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(400, gin.H{
			"error": "check input data",
		})
		return
	}

	category := models.Category{
		Name: body.Name,
	}

	error := db.Create(&category).Error

	if error == nil {
		ctx.JSON(200, gin.H{
			"category": category,
		})
	} else {
		ctx.JSON(200, gin.H{
			"status": "error",
		})
	}

}

func view(ctx *gin.Context){
	db := ctx.MustGet("db").(*gorm.DB)
	id := ctx.Param("id")
	var category models.Category
	db.First(&category, id)
	ctx.JSON(200, gin.H{
		"category": category,
	})
}

func all(ctx *gin.Context){
	db := ctx.MustGet("db").(*gorm.DB)
	var categories []models.Category
	db.First(&categories)
	ctx.JSON(200, gin.H{
		"categories": categories,
	})
}
