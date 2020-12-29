package models

import (
	"github.com/Ismail14098/agyn_test_rest/lib/common"
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name string
}

func (c *Category) Serialize() common.JSON {
	return common.JSON{
		"id": c.ID,
		"name": c.Name,
	}
}

func (c *Category) Read(m common.JSON)  {
	c.ID = uint(m["id"].(int))
	c.Name = m["name"].(string)
}