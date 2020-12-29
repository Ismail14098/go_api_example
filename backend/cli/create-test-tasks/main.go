package main

import (
	"fmt"
	"github.com/Ismail14098/agyn_test_rest/api/v1.0/auth"
	"github.com/Ismail14098/agyn_test_rest/database"
	"github.com/Ismail14098/agyn_test_rest/database/models"
	"github.com/Ismail14098/agyn_test_rest/lib/status"
	"github.com/joho/godotenv"
	"strconv"
	"time"
)

func main(){
	err := godotenv.Load()
	db, _ := database.Initialize()

	var user models.User
	password := "test1234"

	passwordhash, err := auth.Hash(password)
	if err != nil {
		panic(err)
	}

	user.Username = "test1234"
	user.PasswordHash = passwordhash
	user.Firstname = "test"
	user.Lastname = "test"
	user.Email = "test@test.com"
	db.Create(&user)

	var role models.Role
	db.Where("name = ?", "user").Find(&role)
	var userRole models.UserRole
	userRole.UserId = user.ID
	userRole.RoleId = role.ID
	db.Create(&userRole)

	var category models.Category
	category.Name = "Test"
	db.Create(&category)

	var tasks []models.Task
	for i:=1;i<11;i++ {
		task := models.Task{
			Title:      "Test task-"+strconv.Itoa(i),
			Text:       "Test task",
			CategoryId: category.ID,
			AuthorId:   user.ID,
			Status:     status.InProgress,
			ExpTime:    time.Now().Local().AddDate(0,0, i),
		}
		tasks = append(tasks,task)
	}
	db.Create(&tasks)

	fmt.Println("Initial population with test user and tasks has done successfully")
}


