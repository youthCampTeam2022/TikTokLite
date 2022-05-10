package model

import "gorm.io/gorm"

type Follow struct {
	gorm.Model
	UserID int64
	FollowerID int64
}
