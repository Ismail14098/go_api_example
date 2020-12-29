package role

import (
"github.com/Ismail14098/agyn_test_rest/database/models"
"github.com/gin-gonic/gin"
"gorm.io/gorm"
)

func all(ctx *gin.Context){
	db := ctx.MustGet("db").(*gorm.DB)
	var roles []models.Role
	db.Find(&roles)
	ctx.JSON(200, gin.H{
		"roles": roles,
	})
}

