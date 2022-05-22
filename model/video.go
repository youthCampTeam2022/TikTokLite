package model

import "gorm.io/gorm"

type Video struct {
	gorm.Model
	AuthorId int64
	Title    string
	PlayUrl  string
	CoverUrl string
}