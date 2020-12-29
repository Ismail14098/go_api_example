package profile

import (
	"github.com/Ismail14098/agyn_test_rest/middlewars"
	"github.com/gin-gonic/gin"
)

func ApplyRoutes(group *gin.RouterGroup){
	profile := group.Group("/profile")
	{
		profile.POST("/edit", middlewars.Authorized, edit)
		profile.POST("/change/:id", middlewars.Authorized, middlewars.IsAdmin, editSudo)
		profile.DELETE("/delete/:id", middlewars.Authorized, middlewars.IsAdmin, drop)
		profile.GET("/show", middlewars.Authorized, middlewars.IsAdmin, show)
		profile.GET("/get/:id", middlewars.Authorized, middlewars.IsAdmin, get)
	}
}
