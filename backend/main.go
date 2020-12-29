package main

import (
	"context"
	"github.com/Ismail14098/agyn_test_rest/api"
	"github.com/Ismail14098/agyn_test_rest/database"
	"github.com/Ismail14098/agyn_test_rest/lib/validators"
	"github.com/Ismail14098/agyn_test_rest/middlewars"
	"github.com/Ismail14098/agyn_test_rest/redis"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"os"
)

func main(){
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	port := os.Getenv("PORT")

	app := gin.Default()

	//CORS for jquery
	app.Use(middlewars.CORS())

	//PostgreSQL
	db, _ := database.Initialize()
	app.Use(database.Inject(db))

	//Context
	ctx := context.Background()
	app.Use(redis.InjectContext(&ctx))

	//Redis
	rdb := redis.Initialize()
	app.Use(redis.Inject(rdb))

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validators.Validate(v)
	}

	app.Use(middlewars.JWTMiddleware())
	api.ApplyRoutes(app)
	app.Run(":" + port)
}

