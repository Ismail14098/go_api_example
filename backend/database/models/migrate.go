package models

import (
	"fmt"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB){
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Role{})
	db.AutoMigrate(&UserRole{})
	db.AutoMigrate(&Category{})
	db.AutoMigrate(&Task{})
	db.AutoMigrate(&Notification{})
	//db.Model(&Model{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	fmt.Println("Auto Migration has done successfully")
}
