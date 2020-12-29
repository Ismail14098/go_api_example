package middlewars

import (
	"github.com/Ismail14098/agyn_test_rest/database/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func IsAdmin(ctx *gin.Context){
	user := ctx.MustGet("user").(models.User)
	db := ctx.MustGet("db").(*gorm.DB)
	
	var adminRole models.Role
	db.Where("name = ?", "admin").Find(&adminRole)
	
	var userRole models.UserRole
	var count int64
	db.Where("user_id = ? AND role_id = ?", user.ID, adminRole.ID).Find(&userRole).Count(&count)

	if count == 0 {
		ctx.AbortWithStatus(401)
		return
	}

}
