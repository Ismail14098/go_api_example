package apiv1

import (
	"github.com/Ismail14098/agyn_test_rest/api/v1.0/auth"
	"github.com/Ismail14098/agyn_test_rest/api/v1.0/category"
	"github.com/Ismail14098/agyn_test_rest/api/v1.0/profile"
	"github.com/Ismail14098/agyn_test_rest/api/v1.0/role"
	"github.com/Ismail14098/agyn_test_rest/api/v1.0/tasks"
	"github.com/gin-gonic/gin"
)

func ping(ctx *gin.Context){
	ctx.JSON(200, gin.H{
		"message": "pong",
	})
}

func ApplyRoutes(group *gin.RouterGroup) {
	v1 := group.Group("/v1")
	{
		auth.ApplyRoutes(v1)
		tasks.ApplyRoutes(v1)
		category.ApplyRoutes(v1)
		profile.ApplyRoutes(v1)
		role.ApplyRoutes(v1)
		v1.GET("/ping", ping)
	}
}
