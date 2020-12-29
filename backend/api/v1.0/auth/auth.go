package auth

import (
	"github.com/Ismail14098/agyn_test_rest/middlewars"
	"github.com/gin-gonic/gin"
)

func ApplyRoutes(group *gin.RouterGroup){
	auth := group.Group("/auth")
	{
		auth.POST("/login", login)
		auth.POST("/register", middlewars.Authorized, middlewars.IsAdmin, register)
		auth.GET("/check", check)
		auth.GET("/logout", middlewars.Authorized, logout)
	}
}
