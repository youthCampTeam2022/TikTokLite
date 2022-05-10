package model

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	VideoID	int64 `gorm:"index"`
	UserID int64
	Content string `gorm:"type:varchar(255);"`
}

func (c *Comment) Create()  {

}

func (c *Comment) Delete()  {

}

func GetCommentNum(videoID int64)int  {
	return 0
}