package service

import (
	"TikTokLite/model"
)

// SetFavorite 点赞
func SetFavorite(videoID, userID int64) error {
	f := model.Favorite{
		VideoID: videoID,
		UserID:  userID,
	}
	return f.UniqueInsert()
}

// CancelFavorite 取消点赞
func CancelFavorite(videoID, userID int64) error {
	f := model.Favorite{
		VideoID: videoID,
		UserID:  userID,
	}
	return f.Delete()
}

// GetFavoriteList 获取喜欢列表
func GetFavoriteList(userID int64) ([]model.VideoRes, error) {
	return model.GetFavoriteRes(userID)
}
