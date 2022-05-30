package model

import (
	"log"
	"strconv"
)

func GetFavoriteNumRedis(videoID int64) (count int64) {
	conn := RedisCache.Conn()
	defer func() {
		_ = conn.Close()
	}()
	favoriteKey := FavoriteKey(videoID)
	//count
	num, err := conn.Do("HGET","favorite",favoriteKey,"EX",5)
	if err != nil {
		count = GetFavoriteNum(videoID)
		SetFavoriteNumRedis(videoID,count)
		return
	}
	return num.(int64)
}

func FavoriteKey(videoID int64)string  {
	return "favorite:" + strconv.FormatInt(videoID,10)
}

func SetFavoriteNumRedis(videoID int64,num int64)  {
	conn := RedisCache.Conn()
	defer func() {
		_ = conn.Close()
	}()
	favoriteKey := FavoriteKey(videoID)
	_, err := conn.Do("HSET","favorite",favoriteKey,num)
	if err != nil {
		log.Print("err in SetFavoriteNumRedis:",err)
		return
	}
}

func IncrFavoriteRedis(videoID int64)  {
	conn := RedisCache.Conn()
	defer func() {
		_ = conn.Close()
	}()
	favoriteKey := FavoriteKey(videoID)
	_, err := conn.Do("HINCRBY","favorite",favoriteKey,1)
	if err != nil {
		log.Print("err in IncrFavoriteRedis:",err)
		return
	}
}

func DecrFavoriteRedis(videoID int64)  {
	conn := RedisCache.Conn()
	defer func() {
		_ = conn.Close()
	}()
	favoriteKey := FavoriteKey(videoID)
	_, err := conn.Do("HINCRBY","favorite",favoriteKey,-1)
	if err != nil {
		log.Print("err in DecrFavoriteRedis:",err)
		return
	}
}