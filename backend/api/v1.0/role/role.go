package role

import (
	"github.com/Ismail14098/agyn_test_rest/middlewars"
	"github.com/gin-gonic/gin"
)

func ApplyRoutes(group *gin.RouterGroup){
	profile := group.Group("/role")
	{
		profile.GET("/all", middlewars.Authorized, all)
		//profile.DELETE("/delete/:id", middlewars.Authorized, middlewars.IsAdmin, drop)
	}
}

