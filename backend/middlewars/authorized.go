package middlewars

import (
	"github.com/gin-gonic/gin"
)

// Authorized blocks unauthorized requests
func Authorized(ctx *gin.Context) {
	_, exists := ctx.Get("user")
	if !exists {
		ctx.AbortWithStatus(401)
		return
	}
}
