package models

import (
	"github.com/Ismail14098/agyn_test_rest/lib/common"
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name string `gorm:"uniqueIndex"`
}

func (r *Role) Serialize() common.JSON{
	return common.JSON{
		"id": r.ID,
		"name": r.Name,
	}
}

func (r *Role) Read(m common.JSON){
	r.ID = uint(m["id"].(float64))
	r.Name = m["name"].(string)
}