package model

import (
	"fmt"
	"log"
	"strconv"
)

const HotFeedKey = "hot"

// HotCounter 暂时设定top20
type HotCounter struct {
	vid int64
	favorite int
	comment int
	time int64
}

// BuildHotFeed 每ns触发一次
func BuildHotFeed()  {
	_, err := RedisCache.Conn().Do("ZREMRANGEBYSCORE", HotFeedKey, "-inf", "+inf")
	if err != nil {
		fmt.Println(err)
		return
	}
	var set map[int64]struct{}
	//去重
	set = make(map[int64]struct{})
	topf := GetTopFavorite(20)
	for key, _ := range topf {
		set[key] = struct{}{}
	}
	topc := GetTopComment(20)
	for key, _ := range topc {
		set[key] = struct{}{}
	}
	for i, _ := range set {
		fmt.Println(i)
		favoriteNum,ok := topf[i]
		if !ok||favoriteNum==0{
			favoriteNum = int(GetFavoriteNumRedis(i))
		}
		commentNum,ok := topc[i]
		if !ok||commentNum==0{
			commentNum = int(GetCommentNumRedis(i))
		}
		h := HotCounter{
			vid:      i,
			favorite: favoriteNum,
			comment:  commentNum,
			time:     GetVideoCreateTime(i),
		}
		InsertHotFeed(i,h.ToScore())
	}
}

func InsertHotFeed(vid int64,score int)  {
	conn := RedisCache.Conn()
	defer func() {
		_ = conn.Close()
	}()
	_, err := conn.Do("ZADD", HotFeedKey,score,vid)
	if err != nil {
		log.Println("err in InsertHotFeed:",err)
	}
}

func PullHotFeed(n int)(res []int)  {
	conn := RedisCache.Conn()
	defer func() {
		_ = conn.Close()
	}()
	do, err := conn.Do("ZRANGEBYSCORE", HotFeedKey,"-inf","+inf")
	if err != nil {
		fmt.Println("err in PullHotFeed:",err)
		return nil
	}
	for _, id := range do.([]interface{}){
		resID, _ := strconv.Atoi(string(id.([]uint8)))
		res = append(res,resID)
	}
	//逆序
	l := len(res)
	for i := 0; i < len(res)/2; i++ {
		res[i],res[l-i-1] = res[l-i-1],res[i]
	}
	if l<n{
		n=l
	}
	return res[:n]
}

// ToScore todo:给时间加权没写好
func (h *HotCounter)ToScore()int  {
	return h.favorite+(h.comment*2)
}
