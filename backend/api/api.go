package api

import (
	apiv1 "github.com/Ismail14098/agyn_test_rest/api/v1.0"
	"github.com/gin-gonic/gin"
)

func ApplyRoutes(app *gin.Engine){
	api := app.Group("/api")
	{
		apiv1.ApplyRoutes(api)
	}
}
