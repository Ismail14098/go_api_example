package main

import (
	"fmt"
	"github.com/Ismail14098/agyn_test_rest/database"
	"github.com/Ismail14098/agyn_test_rest/database/models"
	"github.com/joho/godotenv"
)
func main(){
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	db, _ := database.Initialize()
	userRole := models.Role{Name:"user"}
	adminRole := models.Role{Name:"admin"}
	db.Create(&userRole)
	db.Create(&adminRole)
	fmt.Println("Initial population with roles has done successfully")
}
