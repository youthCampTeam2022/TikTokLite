package model

import (
	"log"
	"strconv"
)

func GetCommentNumRedis(videoID int64) (count int64) {
	conn := RedisCache.Conn()
	defer func() {
		_ = conn.Close()
	}()
	commentKey := CommentKey(videoID)
	//count
	num, err := conn.Do("HGET","comment",commentKey,"EX",5)
	if err != nil {
		count = GetCommentNum(videoID)
		SetCommentNumRedis(videoID,count)
		return
	}
	return num.(int64)
}

func CommentKey(videoID int64)string  {
	return "comment:" + strconv.FormatInt(videoID,10)
}

func SetCommentNumRedis(videoID int64,num int64)  {
	conn := RedisCache.Conn()
	defer func() {
		_ = conn.Close()
	}()
	commentKey := CommentKey(videoID)
	_, err := conn.Do("HSET","comment",commentKey,num)
	if err != nil {
		log.Print("err in SetCommentNumRedis:",err)
		return
	}
}