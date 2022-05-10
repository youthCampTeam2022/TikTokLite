package model

import "gorm.io/gorm"

type Video struct {
	gorm.Model
	Author int64
	Name string
	CoverName string
}
