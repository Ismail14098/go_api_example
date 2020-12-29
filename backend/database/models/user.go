package models

import (
	"github.com/Ismail14098/agyn_test_rest/lib/common"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username 	 string `gorm:"uniqueIndex"`
	Email 		 string `gorm:"uniqueIndex"`
	Firstname  	 string
	Lastname     string
	PasswordHash string
}

func (u *User) Serialize() common.JSON{
	return common.JSON{
		"id": u.ID,
		"username": u.Username,
		"email": u.Email,
		"firstname": u.Firstname,
		"lastname": u.Lastname,
	}
}

func (u *User) Read(m common.JSON){
	u.ID = uint(m["id"].(float64))
	u.Username = m["username"].(string)
	u.Email = m["email"].(string)
	u.Firstname = m["firstname"].(string)
	u.Lastname = m["lastname"].(string)
}