package profile

import (
	"context"
	"github.com/Ismail14098/agyn_test_rest/database/models"
	"github.com/Ismail14098/agyn_test_rest/lib/common"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"strconv"
)

func edit(ctx *gin.Context){
	type RequestBody struct {
		Username string `json:"username" binding:"required,correctLogin"`
		Email string `json:"email" binding:"required"`
		Firstname string `json:"firstname" binding:"required,correctName"`
		Lastname string `json:"lastname" binding:"required,correctName"`
	}
	var body RequestBody
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(400, gin.H{
			"error": "check input data",
		})
		return
	}

	user := ctx.MustGet("user").(models.User)
	db := ctx.MustGet("db").(*gorm.DB)

	user.Username = body.Username
	user.Email = body.Email
	user.Firstname = body.Firstname
	user.Lastname = body.Lastname

	num := db.Save(&user).RowsAffected
	if num > 0 {
		ctx.Set("user",user)

		ctx.JSON(200, common.JSON{
			"user":  user.Serialize(),
		})
	} else {
		ctx.JSON(400, common.JSON{
			"status": "error",
		})
	}
}

func editSudo(ctx *gin.Context){
	type RequestBody struct {
		Username string `json:"username" binding:"required,correctLogin"`
		Email string `json:"email" binding:"required"`
		Firstname string `json:"firstname" binding:"required,correctName"`
		Lastname string `json:"lastname" binding:"required,correctName"`
		RoleId string `json:"role" binding:"required"`
	}

	var body RequestBody
	if err := ctx.BindJSON(&body); err != nil {
		ctx.JSON(400, gin.H{
			"error": "check input data",
		})
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	roleId, err1 := strconv.Atoi(body.RoleId)
	if err != nil || err1 != nil {
		ctx.AbortWithStatus(400)
		return
	}

	db := ctx.MustGet("db").(*gorm.DB)
	var user models.User
	db.First(&user, id)

	var userRole models.UserRole
	db.Where("user_id = ?", user.ID).Find(&userRole)
	userRole.RoleId = uint(roleId)
	db.Save(&userRole)

	user.Username = body.Username
	user.Email = body.Email
	user.Firstname = body.Firstname
	user.Lastname = body.Lastname

	num := db.Save(&user).RowsAffected
	if num > 0 {
		ctx.Set("user",user)

		ctx.JSON(200, common.JSON{
			"user":  user.Serialize(),
		})
	} else {
		ctx.JSON(400, common.JSON{
			"status": "error",
		})
	}
}

func drop(ctx *gin.Context){
	db := ctx.MustGet("db").(*gorm.DB)

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "check id",
		})
		return
	}
	var user models.User
	result := db.Find(&user,id)

	if result.RowsAffected == 0 {
		ctx.AbortWithStatus(400)
		return
	}

	var roles []models.UserRole
	result = db.Where("user_id = ?", id).Delete(&roles)

	var tasks []models.Task
	result = db.Where("author_id = ?", id).Delete(&tasks)

	rdb := ctx.MustGet("rdb").(*redis.Client)
	defCtx := ctx.MustGet("defCtx").(*context.Context)
	key := "tasks_user_"+strconv.Itoa(id)
	rdb.Del(*defCtx, key)

	result = db.Delete(&user)
	num := result.RowsAffected

	if num > 0 {
		ctx.JSON(200, gin.H{
			"status": "success",
		})
	} else {
		ctx.JSON(400, gin.H{
			"status": "error",
		})
	}
}

func show(ctx *gin.Context){
	db := ctx.MustGet("db").(*gorm.DB)

	var userRoles []models.UserRole
	page := ctx.DefaultQuery("page", "1")
	pageInt, err := strconv.Atoi(page)

	if err != nil || pageInt < 1 {
		ctx.JSON(400, gin.H{
			"error": "check id",
		})
		return
	}

	offset := 10 * (pageInt-1)
	db.Preload("User").Preload("Role").Offset(offset).Limit(10).Find(&userRoles)

	if len(userRoles) == 0 {
		ctx.AbortWithStatus(400)
		return
	}

	ctx.JSON(200, gin.H{
		"users": userRoles,
	})
}

func get(ctx *gin.Context){
	db := ctx.MustGet("db").(*gorm.DB)

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, gin.H{
			"error": "check id",
		})
		return
	}
	var user models.User
	db.Find(&user,id)

	ctx.JSON(200, gin.H{
		"user": user,
	})
}
