package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name string `gorm:"type:char(10);"`
	Password string `gorm:"type:char(20);"`
}
