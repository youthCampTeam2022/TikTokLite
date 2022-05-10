package model

import "gorm.io/gorm"

type Favorite struct {
	gorm.Model
	VideoID int64 `gorm:"index"`
	UserID int64
}

func (f *Favorite) Create()  {

}

func (f *Favorite) Delete()  {

}

func GetFavoriteNum(videoID int64)int  {
	return 0
}