package service

import (
	"TikTokLite/model"
	"fmt"
	"time"
)

func GetFeed(latestTime time.Time, userId int64) ([]Video, int64, error) {
	followService := NewFollowService()
	latestVideo, err := model.GetLatestVideo()
	if err != nil {
		return nil, time.Now().UnixMilli(), err
	}
	var videos []model.Video
	//如果当前时间后的投稿数为0就直接返回数据库中最新的N的视频
	if latestVideo.CreatedAt.UnixMilli() <= latestTime.UnixMilli() {
		fmt.Println("get the latest N")
		videos, err = model.GetTheLatestNVideos()
	} else {
		videos, err = model.GetVideosByLatestTime(latestTime)
	}
	// 根据查询到的videos去查询相关信息，未登录默认userID为-1，is_follow，is_favorite都为false
	videoList := make([]Video, len(videos))
	for i := range videoList {
		videoId := int64(videos[i].ID)
		author := BuildUser(userId, videos[i].AuthorId, followService.FollowRepository)
		isFavorite, err := model.IsFavorite(userId, videoId)
		if err != nil {
			return nil, time.Now().UnixMilli(), err
		}
		videoList[i].Id = videoId
		videoList[i].Author = author
		videoList[i].Title = videos[i].Title
		videoList[i].PlayUrl = videos[i].PlayUrl
		videoList[i].CoverUrl = videos[i].CoverUrl
		videoList[i].FavoriteCount = model.GetFavoriteNum(videoId)
		videoList[i].CommentCount = model.GetCommentNum(videoId)
		videoList[i].IsFavorite = isFavorite
	}
	return videoList, videos[0].CreatedAt.UnixMilli(), nil
}
