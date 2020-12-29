package models

import (
	"github.com/Ismail14098/agyn_test_rest/lib/common"
	"gorm.io/gorm"
)

type UserRole struct {
	gorm.Model
	UserId uint
	User User `gorm:"foreignkey:UserId"`
	RoleId uint
	Role Role `gorm:"foreignkey:RoleId"`
}

func (ur *UserRole) Serialize() common.JSON {
	return common.JSON{
		"id": ur.ID,
		"user": ur.User.Serialize(),
		"role": ur.Role.Serialize(),
	}
}

func (ur *UserRole) Read(m common.JSON)  {
	ur.ID = uint(m["id"].(float64))
	ur.UserId = m["user_id"].(uint)
	ur.RoleId = m["role_id"].(uint)
}