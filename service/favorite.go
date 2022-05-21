package service

import (
	"TikTokLite/model"
)

func SetFavorite(videoID, userID int64) error {
	f := model.Favorite{
		VideoID: videoID,
		UserID:  userID,
	}
	return f.Create()
}

func CancelFavorite(videoID, userID int64) error {
	f := model.Favorite{
		VideoID: videoID,
		UserID:  userID,
	}
	return f.Delete()
}

func GetFavoriteList(userID int64) ([]Video, error) {
	f := &model.Favorite{
		UserID: userID,
	}
	list, err := f.GetFavoriteList()
	if err != nil {
		return nil, err
	}
	videoResult := make([]Video, len(list))
	for i, vid := range list {
		videoResult[i] = Video{
			Id: vid,
		}
		//todo: video缺信息
	}
	return videoResult, nil

}
