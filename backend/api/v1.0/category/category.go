package category
import (
	"github.com/Ismail14098/agyn_test_rest/middlewars"
	"github.com/gin-gonic/gin"
)

func ApplyRoutes(group *gin.RouterGroup){
	category := group.Group("/category")
	{
		category.POST("/create", middlewars.Authorized, middlewars.IsAdmin, create)
		category.GET("/view/:id", middlewars.Authorized, middlewars.IsAdmin, view)
		category.GET("/all", all)
	}
}


