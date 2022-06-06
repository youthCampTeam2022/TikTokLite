package model

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"strconv"
)

const CommentSet = "commentSet"

func GetCommentNumRedis(videoID int64) (count int64) {
	conn := RedisCache.Conn()
	defer func() {
		_ = conn.Close()
	}()
	commentKey := ID2CommentKey(videoID)
	num, err := redis.Int64(conn.Do("ZSCORE", CommentSet, commentKey))
	if err != nil {
		count = GetCommentNum(videoID)
		fmt.Println(count)
		SetCommentNumRedis(videoID, count)
		return
	}
	return num
}

func ID2CommentKey(videoID int64) string {
	return "comment:" + strconv.FormatInt(videoID, 10)
}

func CommentKey2ID(key string) int64 {
	res, _ := strconv.Atoi(key[8:])
	return int64(res)
}

func SetCommentNumRedis(videoID int64, num int64) {
	conn := RedisCache.Conn()
	defer func() {
		_ = conn.Close()
	}()
	commentKey := ID2CommentKey(videoID)
	_, err := conn.Do("ZADD", CommentSet, num, commentKey)
	if err != nil {
		log.Print("err in SetCommentNumRedis:", err)
		return
	}
}

func IncrCommentRedis(videoID int64) {
	conn := RedisCache.Conn()
	defer func() {
		_ = conn.Close()
	}()
	favoriteKey := ID2CommentKey(videoID)
	_, err := conn.Do("ZINCRBY", CommentSet, 1 ,favoriteKey)
	if err != nil {
		log.Print("err in IncrCommentRedis:", err)
		return
	}
}

func DecrCommentRedis(videoID int64) {
	conn := RedisCache.Conn()
	defer func() {
		_ = conn.Close()
	}()
	favoriteKey := ID2CommentKey(videoID)
	_, err := conn.Do("ZINCRBY", CommentSet, -1, favoriteKey)
	if err != nil {
		log.Print("err in DecrCommentRedis:", err)
		return
	}
}

func GetTopComment(n int) (top map[int64]int) {
	conn := RedisCache.Conn()
	defer func() {
		_ = conn.Close()
	}()
	values, err := redis.Values(conn.Do("ZREVRANGE", CommentSet,0,n,"WITHSCORES"))
	if err != nil {
		log.Println("err in GetTopComment:", err)
		return nil
	}
	top = make(map[int64]int)
	for i := 0; i < len(values); i += 2 {
		key,_ := redis.String(values[i],nil)
		v, _ := redis.Int64(values[i+1],nil)
		if CommentKey2ID(key) == 0||v==0{
			continue
		}
		top[CommentKey2ID(key)] = int(v)
	}
	fmt.Println("topComm",top)
	return top
}
