package model

import (
	"errors"
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	VideoID	int64 `gorm:"index"`
	UserID int64
	Content string `gorm:"type:varchar(255);"`
	User User
}


func (c *Comment) Create()error  {
	return DB.Create(&c).Error
}

func (c *Comment) Delete() error {
	return DB.Delete(&c).Error
}

func (c *Comment) DeleteByUser() error {
	var uid Comment
	DB.Model(&c).First(&uid)
	if uid.User.ID == c.User.ID{
		return DB.Delete(&c).Error
	}
	return errors.New("invalid delete")
}

func GetCommentNum(videoID int64)int  {
	return 0
}

func GetCommentsByVideo(videoID int64)(comments []Comment,err error)  {
	err = DB.Where("video_id = ?",videoID).Find(&comments).Order("created_at DESC").Error
	return
}

func GetCommentsByVideoJoin(videoID int64)(comments []Comment,err error)  {
	err = DB.Preload("User").Where("video_id = ?",videoID).Find(&comments).Order("created_at DESC").Error
	return
}