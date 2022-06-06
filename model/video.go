package model

import (
	"fmt"
	//"github.com/gomodule/redigo/redis"
	"github.com/gistao/RedisGo-Async/redis"
	"gorm.io/gorm"
	"log"
	"strconv"
	"time"
)

const (
	feedSize  = 4
	aliveTime = time.Hour * 24
)

type Video struct {
	gorm.Model
	AuthorId int64
	Title    string `gorm:"type:varchar(255)" ,json:"title"`
	PlayUrl  string `gorm:"type:varchar(255)" ,json:"play_url"`
	CoverUrl string `gorm:"type:varchar(255)" ,json:"cover_url"`
}

func (v *Video) Create() error {

	err := DB.Create(&v).Error
	now := time.Now().UnixMilli()
	//更新authorfeed
	_ = insertAuthorFeed(v.AuthorId, int64(v.ID), now)
	//更新userfeed
	_ = pushNewVideoToActiveUsersFeed(v.AuthorId, int64(v.ID), now)
	return err
	//if err != nil {
	//	return errors.New("insert AuthorFeed in Redis failed")
	//}
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

func GetVideoCreateTime(videoID int64) int64 {
	var t time.Time
	DB.Model(&Video{}).Where("id = ?", videoID).Select("created_at").Scan(&t)
	return t.UnixMilli()
}

// Authorfeed增加新的视频
func insertAuthorFeed(userID, videoID, now int64) (err error) {
	authorFeedKey := fmt.Sprintf("%s:%s", strconv.FormatInt(userID, 10), "authorfeed")
	conn := RedisCache.Conn()
	defer conn.Close()
	_, err = conn.Do("ZADD", authorFeedKey, now, videoID)
	return
}
func pushNewVideoToActiveUsersFeed(userID, videoID, now int64) (err error) {
	var followers []int64
	followers, err = NewFollowManagerRepository().RedisGetFollowerList(userID)
	if err != nil {
		log.Println("fail to get follower list by", userID)
		return err
	}
	var loginTimeKey, userFeedKey, ids string
	var loginTime int64
	//conn := RedisCache.Conn()
	conn := RedisCache.AsynConn()
	defer conn.Close()
	for i := 0; i < len(followers); i++ {
		ids = strconv.FormatInt(followers[i], 10)
		loginTimeKey = fmt.Sprintf("%s", ids)
		//loginTime, err = redis.Int64(conn.Do("GET", loginTimeKey))
		loginTime, err = redis.Int64(conn.Do("HGET", "aliveUser", loginTimeKey))
		if err != nil {
			if err == redis.ErrNil {
				continue
			}
			return err
		}
		//检查登录是否超时
		if time.UnixMilli(loginTime).Add(aliveTime).After(time.Now()) {
			userFeedKey = fmt.Sprintf("%s:%s", ids, "userfeed")
			//_, err = conn.Do("ZADD", userFeedKey, now, videoID)
			_, err = conn.AsyncDo("ZADD", userFeedKey, now, videoID)
			if err != nil {
				log.Println("userFeed Push failed:", err.Error())
			}
		} else {
			//_, _ = conn.Do("HDEL", "aliveUser", loginTimeKey)
			_, _ = conn.AsyncDo("HDEL", "aliveUser", loginTimeKey)
		}
	}
	userFeedKey = fmt.Sprintf("%s:%s", strconv.FormatInt(userID, 10), "userfeed")
	//_, _ = conn.Do("ZADD", userFeedKey, now, videoID)
	_, _ = conn.AsyncDo("ZADD", userFeedKey, now, videoID)
	return nil
}
func GetUserFeedRedis(latestTime time.Time, userId int64) ([]int64, error) {
	id := strconv.FormatInt(userId, 10)
	key := fmt.Sprintf("%s:%s", id, "userfeed")
	//conn := RedisCache.Conn()
	conn := RedisCache.AsynConn()
	defer conn.Close()
	var err error
	var offset, timeStamp int64
	//存储对应用户feed流的偏移量
	userOffset := fmt.Sprintf("%s:%s", id, "offset")
	offset, err = redis.Int64(conn.Do("get", userOffset))
	if err == redis.ErrNil {
		offset = 0
	}
	//存储对应用户feed流的起始时间
	userFeedTimeStamp := fmt.Sprintf("%s:%s", id, "feedtimestamp")
	timeStamp, err = redis.Int64(conn.Do("GET", userFeedTimeStamp))
	if err == redis.ErrNil {
		timeStamp = latestTime.UnixMilli()
		//_,_=conn.Do("SET", userFeedTimeStamp, timeStamp)
		_, _ = conn.AsyncDo("SET", userFeedTimeStamp, timeStamp)
	}
	var vals []int64
	vals, err = redis.Int64s(conn.Do("ZREVRANGEBYSCORE", key, timeStamp, 0, "withscores", "limit", offset, feedSize))
	offset += feedSize
	//意味着feed流已查询到底
	if len(vals) < feedSize*2 {
		offset = 0
		timeStamp = time.Now().UnixMilli()
		//_,_=conn.Do("SET", userFeedTimeStamp, timeStamp)
		_, _ = conn.AsyncDo("SET", userFeedTimeStamp, timeStamp)
	}
	//_, _ = conn.Do("SET", userOffset, offset)
	_, _ = conn.AsyncDo("SET", userOffset, offset)
	return vals, nil
}

func GetVideoByID(videoID int64) (Video, error) {
	var video Video
	err := DB.Where("id=?", videoID).First(&video).Error
	return video, err
}
