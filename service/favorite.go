package service

import (
	"TikTokLite/model"
)

func SetFavorite(videoID, userID int64) error {
	f := model.Favorite{
		VideoID: videoID,
		UserID:  userID,
	}
	return f.UniqueInsert()
}

func CancelFavorite(videoID, userID int64) error {
	f := model.Favorite{
		VideoID: videoID,
		UserID:  userID,
	}
	return f.Delete()
}

func GetFavoriteList(userID int64) ([]model.VideoRes, error) {
	return model.GetFavoriteRes(userID)
}
