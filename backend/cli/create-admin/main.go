package main

import (
	"bufio"
	"fmt"
	"github.com/Ismail14098/agyn_test_rest/api/v1.0/auth"
	"github.com/Ismail14098/agyn_test_rest/database"
	"github.com/Ismail14098/agyn_test_rest/database/models"
	"github.com/joho/godotenv"
	"os"
)

func main(){
	err := godotenv.Load()
	db, _ := database.Initialize()

	var user models.User
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Enter username:")
	var username string
	for scanner.Scan() {
		username = scanner.Text()
		if len(username) > 5 && len(username) < 31 {
			break
		}
		fmt.Println("Enter again ! min len 6, max 30")
	}

	fmt.Println("Enter password:")
	var password string
	for scanner.Scan() {
		password = scanner.Text()
		if len(password) > 7 && len(password) < 37 {
			break
		}
		fmt.Println("Enter again ! min len 8, max 36")
	}

	passwordhash, err := auth.Hash(password)
	if err != nil {
		panic(err)
	}

	user.Username = username
	user.PasswordHash = passwordhash
	user.Firstname = "admin"
	user.Lastname = "admin"
	user.Email = "admin@admin.com"
	db.Create(&user)

	var role models.Role
	db.Where("name = ?", "admin").Find(&role)
	var userRole models.UserRole
	userRole.UserId = user.ID
	userRole.RoleId = role.ID
	db.Create(&userRole)

	fmt.Println("Initial population with admin user has done successfully")
}