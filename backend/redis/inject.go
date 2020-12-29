package redis

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func Inject(rdb *redis.Client) gin.HandlerFunc{
	return func(ctx *gin.Context) {
		ctx.Set("rdb",rdb)
		ctx.Next()
	}
}

func InjectContext(c *context.Context) gin.HandlerFunc{
	return func(ctx *gin.Context) {
		ctx.Set("defCtx",c)
		ctx.Next()
	}
}
