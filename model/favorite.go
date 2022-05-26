package model

import (
	"gorm.io/gorm"
)

type Favorite struct {
	gorm.Model
	VideoID int64 `gorm:"index"`
	UserID  int64
}

func (f *Favorite) Create() error {
	return DB.Create(&f).Error
}

func (f *Favorite) Delete() error {
	return DB.Delete(&f).Error
}

func GetFavoriteNum(videoID int64) (count int64) {
	DB.Model(&Favorite{}).Where("video_id = ?", videoID).Count(&count)
	return
}

func (f *Favorite) GetFavoriteList() ([]int64, error) {
	var res []int64
	err := DB.Model(f).Select("video_id").Where("user_id = ?", f.UserID).Scan(&res).Error
	return res, err
}

func IsFavorite(userId, videoId int64) (bool, error) {
	var count int64
	err := DB.Model(&Favorite{}).Where("video_id = ? and user_id = ?", videoId, userId).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, err
}
