package models

import "gorm.io/gorm"

type Notification struct {
	gorm.Model
	Text string
	UserId uint
	User User `gorm:"foreignkey:UserId"`
}
