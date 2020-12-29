package tasks

import (
	"github.com/Ismail14098/agyn_test_rest/middlewars"
	"github.com/gin-gonic/gin"
)

func ApplyRoutes(group *gin.RouterGroup){
	tasks := group.Group("/task")
	{
		tasks.POST("/create", middlewars.Authorized, create)
		tasks.GET("/view/:id", middlewars.Authorized, view)
		tasks.GET("/all", middlewars.Authorized, all)
		tasks.GET("/show", middlewars.Authorized, middlewars.IsAdmin, show)
		tasks.POST("/edit", middlewars.Authorized, edit)
		tasks.POST("/change", middlewars.Authorized, middlewars.IsAdmin, editSudo)
		tasks.DELETE("/delete/:id", middlewars.Authorized, drop)
		tasks.DELETE("/drop/:id", middlewars.Authorized, middlewars.IsAdmin, dropSudo)
	}
}

