package models

import (
	"github.com/Ismail14098/agyn_test_rest/lib/common"
	"gorm.io/gorm"
	"time"
)

type Task struct {
	gorm.Model
	Title string
	Text string `gorm:"type:text"`
	CategoryId uint
	Category Category `gorm:"foreignkey:CategoryId"`
	AuthorId uint
	Author User `gorm:"foreignkey:AuthorId"`
	Status string
	ExpTime time.Time
}

func (t *Task) Serialize() common.JSON {
	return common.JSON{
		"id": t.ID,
		"title": t.Title,
		"text": t.Text,
		"category": t.Category.Serialize(),
		"author": t.Author.Serialize(),
		"status": t.Status,
		"expTime": t.ExpTime,
	}
}

func (t *Task) Read(m common.JSON)  {
	t.ID = uint(m["id"].(int))
	t.Title = m["title"].(string)
	t.Text = m["text"].(string)
	t.CategoryId = m["category_id"].(uint)
	t.AuthorId = m["author_id"].(uint)
	t.Status = m["status"].(string)
	t.ExpTime = m["expTime"].(time.Time)
}