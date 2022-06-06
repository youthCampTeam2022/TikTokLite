package service

import (
	"TikTokLite/model"
	"fmt"
	//"github.com/gomodule/redigo/redis"
	"github.com/gistao/RedisGo-Async/redis"
	"log"
	"strconv"
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
		videoList[i].FavoriteCount = model.GetFavoriteNumRedis(videoId)
		videoList[i].CommentCount = model.GetCommentNumRedis(videoId)
		videoList[i].IsFavorite = isFavorite
	}
	return videoList, videos[0].CreatedAt.UnixMilli(), nil
}

func BuildVideo(userID int64, _video model.Video) Video {
	var video Video
	videoID := int64(_video.ID)
	isFavorite, err := model.IsFavorite(userID, videoID)
	//log.Println(video)
	if err != nil {
		return video
	}
	video.Id = videoID
	video.Author = BuildUser(userID, _video.AuthorId, NewFollowService().FollowRepository)
	video.Title = _video.Title
	video.PlayUrl = _video.PlayUrl
	video.CoverUrl = _video.CoverUrl
	video.FavoriteCount = model.GetFavoriteNumRedis(videoID)
	video.CommentCount = model.GetCommentNumRedis(videoID)
	video.IsFavorite = isFavorite
	return video
}
func GetUserFeed(latestTime time.Time, userId int64) ([]Video, int64, error) {
	videoIDs, _ := model.GetUserFeedRedis(latestTime, userId)
	//conn := model.RedisCache.Conn()
	//defer conn.Close()
	var video Video
	var videos []Video
	var v model.Video
	//log.Println("videoID:", videoIDs)
	for i := 0; i < len(videoIDs); i += 2 {
		id := videoIDs[i]
		//s, err := redis.Bytes(conn.Do("HGET", "videos", id))
		//if err != nil {
		//	if err == redis.ErrNil {
		//		v, err = model.GetVideoByID(id)
		//		video = BuildVideo(userId, v)
		//		val, err := json.Marshal(video)
		//		if err != nil {
		//			log.Println("marshal video failed:", err.Error())
		//		}
		//		_, err = conn.Do("HMSET", "videos", id, val)
		//		if err != nil {
		//			log.Println("cache video in redis failed:", err.Error())
		//		}
		//	} else {
		//		return nil, -1, err
		//	}
		//} else {
		//	err := json.Unmarshal(s, &video)
		//	//log.Println("!!!!!!!", video)
		//	if err != nil {
		//		log.Println("unmarshal failed:", err.Error())
		//		return nil, -1, err
		//	}
		//}
		//log.Println(video)
		var err error
		v, err = model.GetVideoByID(id)
		if err != nil {
			log.Println("get video by id failed:", err.Error())
			return nil, -1, err
		}
		video = BuildVideo(userId, v)
		videos = append(videos, video)
	}
	//log.Println(videos)
	var nextTime = time.Now().UnixMilli()
	if len(videos) > 0 {
		nextTime = videoIDs[len(videoIDs)-1]
	}
	return videos, nextTime, nil
}

func UserFeedInit(userID int64) {
	follows, err := model.NewFollowManagerRepository().RedisGetFollowList(userID)
	if err != nil {
		log.Println("get followlist failed:", err.Error())
		return
	}
	//conn := model.RedisCache.Conn()
	conn := model.RedisCache.AsynConn()
	defer conn.Close()
	userFeedKey := fmt.Sprintf("%s:%s", strconv.FormatInt(userID, 10), "userfeed")
	for _, id := range follows {
		authorkey := fmt.Sprintf("%s:%s", strconv.FormatInt(id, 10), "authorfeed")
		vals, _ := redis.Values(conn.Do("ZREVRANGEBYSCORE", authorkey, time.Now().UnixMilli(), 0, "withscores", "limit", 0, 10))
		for i := 0; i < len(vals); i += 2 {
			k, _ := redis.Int64(vals[i], nil)
			v, _ := redis.Int64(vals[i+1], nil)

			//_, err = conn.Do("ZADD", userFeedKey, v, k)
			_, err = conn.AsyncDo("ZADD", userFeedKey, v, k)
			if err != nil {
				log.Println("userfeed set failed:", err.Error())
			}
		}
	}
	hots := model.PullHotFeed(20)
	for i := 0; i < len(hots); i++ {
		createTime := model.GetVideoCreateTime(hots[i])
		//_,_=conn.Do("ZADD", userFeedKey, createTime, hots[i])
		_, _ = conn.AsyncDo("ZADD", userFeedKey, createTime, hots[i])
	}
}

func AuthorFeedPushToNewFollower(authorID, followerID int64) {
	//conn := model.RedisCache.Conn()
	conn := model.RedisCache.AsynConn()
	defer conn.Close()
	authorFeedKey := fmt.Sprintf("%s:%s", strconv.FormatInt(authorID, 10), "authorfeed")
	videos, err := redis.Int64s(conn.Do("ZREVRANGEBYSCORE", authorFeedKey, "+inf", "-inf", "withscores", "limit", 0, 10))
	if err != nil {
		log.Println("get authorfeed error:", err.Error())
		return
	}
	userFeedKey := fmt.Sprintf("%s:%s", strconv.FormatInt(followerID, 10), "userfeed")
	for i := 0; i < len(videos); i += 2 {
		//_, err = conn.Do("ZADD", userFeedKey, videos[i+1], videos[i])
		_, err = conn.AsyncDo("ZADD", userFeedKey, videos[i+1], videos[i])
	}
}
