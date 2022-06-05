package model

import (
	"gorm.io/gorm"
	"time"
)

const (
	feedSize = 3
)

type Video struct {
	gorm.Model
	AuthorId int64
	Title    string `gorm:"type:varchar(255)"`
	PlayUrl  string `gorm:"type:varchar(255)"`
	CoverUrl string `gorm:"type:varchar(255)"`
}

func (v *Video) Create() error {
	return DB.Create(&v).Error
}

func GetVideosByUserId(userId int64) ([]Video, error) {
	var videos []Video
	query := DB.Where("author_id = ?", userId).Find(&videos)
	return videos, query.Error
}

func GetVideosByLatestTime(latestTime time.Time) ([]Video, error) {
	var videos []Video
	query := DB.Order("created_at desc").Where("created_at > ?", latestTime).Limit(feedSize).Find(&videos)
	return videos, query.Error
}

func GetTheLatestNVideos() ([]Video, error) {
	var videos []Video
	query := DB.Order("created_at desc").Limit(feedSize).Find(&videos)
	return videos, query.Error
}

func GetLatestVideo() (Video, error) {
	var video Video
	query := DB.Last(&video)
	return video, query.Error
}

func GetVideoCreateTime(videoID int64)int64  {
	var t time.Time
	DB.Model(&Video{}).Where("id = ?",videoID).Select("created_at").Scan(&t)
	return t.Unix()
}