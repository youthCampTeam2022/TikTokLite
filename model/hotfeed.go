package model

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"time"
)

const HotFeedKey = "hot"

// HotCounter 暂时设定top20
type HotCounter struct {
	vid      int64
	favorite int
	comment  int
	time     int64
}

// BuildHotFeed 每ns触发一次
func BuildHotFeed() {
	var set map[int64]bool
	//去重
	set = make(map[int64]bool)
	topf := GetTopFavorite(20)
	for key, _ := range topf {
		set[key] = true
	}
	topc := GetTopComment(20)
	for key, _ := range topc {
		set[key] = true
	}
	for i, _ := range set {
		favoriteNum, ok := topf[i]
		if !ok || favoriteNum == 0 {
			favoriteNum = int(GetFavoriteNumRedis(i))
		}
		commentNum, ok := topc[i]
		if !ok || commentNum == 0 {
			commentNum = int(GetCommentNumRedis(i))
		}
		h := HotCounter{
			vid:      i,
			favorite: favoriteNum,
			comment:  commentNum,
			time:     GetVideoCreateTime(i),
		}
		InsertHotFeed(i, h.ToScore())
	}
}

func InsertHotFeed(vid int64, score int) {
	conn := RedisCache.Conn()
	defer func() {
		_ = conn.Close()
	}()
	_, err := conn.Do("ZADD", HotFeedKey, score, vid)
	if err != nil {
		log.Println("err in InsertHotFeed:", err)
	}
}

func PullHotFeed(n int) []int64 {
	conn := RedisCache.Conn()
	defer func() {
		_ = conn.Close()
	}()
	res, err := redis.Int64s(conn.Do("ZREVRANGEBYSCORE", HotFeedKey, "+inf", "-inf"))
	if err != nil {
		fmt.Println("err in PullHotFeed:", err)
		return nil
	}
	if n > len(res) {
		n = len(res)
	}
	return res[:n]
}

// ToScore todo:给时间加权没写好
func (h *HotCounter) ToScore() int {
	return h.favorite + (h.comment * 2)
}

func CheckAliveUserAndPushHotFeed() {
	conn := RedisCache.Conn()
	defer conn.Close()

	vals, err := redis.Int64Map(conn.Do("HGETALL", "aliveUser"))
	if err != nil {
		log.Println("push hotfeed error:", err.Error())
		return
	}
	hots := PullHotFeed(20)
	var userFeedKey string
	for k, v := range vals {
		if time.UnixMilli(v).Add(aliveTime).After(time.Now()) {
			userFeedKey = fmt.Sprintf("%s:%s", k, "userfeed")
			for i := 0; i < len(hots); i++ {
				createTime := GetVideoCreateTime(hots[i])
				conn.Do("ZADD", userFeedKey, createTime, hots[i])
			}
		} else {
			_, _ = conn.Do("HDEL", "aliveUser", k)
		}
	}
}
