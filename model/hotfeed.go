package model

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"strconv"
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

func PullHotFeed(n int) (res []int) {
	conn := RedisCache.Conn()
	defer func() {
		_ = conn.Close()
	}()
	do, err := conn.Do("ZRANGEBYSCORE", HotFeedKey, "-inf", "+inf")
	if err != nil {
		fmt.Println("err in PullHotFeed:", err)
		return nil
	}
	for _, id := range do.([]interface{}) {
		resID, _ := strconv.Atoi(string(id.([]uint8)))
		res = append(res, resID)
	}
	//逆序
	l := len(res)
	for i := 0; i < len(res)/2; i++ {
		res[i], res[l-i-1] = res[l-i-1], res[i]
	}
	return res[:n]
}
func PullHotFeed2(n int) []int64 {
	conn := RedisCache.Conn()
	defer func() {
		_ = conn.Close()
	}()
	res, err := redis.Int64s(conn.Do("ZREVRANGEBYSCORE", HotFeedKey, "+inf", "-inf"))
	if err != nil {
		fmt.Println("err in PullHotFeed:", err)
		return nil
	}
	if 2*n > len(res) {
		n = len(res) / 2
	}
	return res[:2*n]
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
	hots := PullHotFeed2(20)
	var userFeedKey string
	for k, v := range vals {
		if time.UnixMilli(v).Add(aliveTime).After(time.Now()) {
			userFeedKey = fmt.Sprintf("%s:%s", k, "userfeed")
			for i := 0; i < len(hots); i += 2 {
				conn.Do("ZADD", userFeedKey, hots[i+1], hots[i])
			}
		} else {
			_, _ = conn.Do("HDEL", "aliveUser", k)
		}
	}
}
