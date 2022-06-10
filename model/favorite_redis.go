package model

import (
	//"github.com/gomodule/redigo/redis"
	"github.com/gistao/RedisGo-Async/redis"
	"log"
	"strconv"
)

const FavoriteSortedSet = "favoriteSet"

func GetFavoriteNumRedis(videoID int64) (count int64) {
	conn := RedisCache.Conn()
	defer func() {
		_ = conn.Close()
	}()
	favoriteKey := ID2FavoriteKey(videoID)
	num, err := redis.Int64(conn.Do("ZSCORE", FavoriteSortedSet, favoriteKey))
	if err != nil {
		count = GetFavoriteNum(videoID)
		SetFavoriteNumRedis(videoID, count)
		return
	}
	return num
}

func ID2FavoriteKey(videoID int64) string {
	return "favorite:" + strconv.FormatInt(videoID, 10)
}

func FavoriteKey2ID(key string) int64 {
	res, _ := strconv.Atoi(key[9:])
	return int64(res)
}

func SetFavoriteNumRedis(videoID int64, num int64) {
	conn := RedisCache.AsynConn()
	defer func() {
		_ = conn.Close()
	}()
	favoriteKey := ID2FavoriteKey(videoID)
	_, err := conn.AsyncDo("ZADD", FavoriteSortedSet, num, favoriteKey)
	if err != nil {
		log.Print("err in SetFavoriteNumRedis:", err)
		return
	}
}

func IncrFavoriteRedis(videoID int64) {
	conn := RedisCache.AsynConn()
	defer func() {
		_ = conn.Close()
	}()
	favoriteKey := ID2FavoriteKey(videoID)
	_, err := conn.AsyncDo("ZINCRBY", FavoriteSortedSet, 1, favoriteKey)
	if err != nil {
		log.Print("err in IncrFavoriteRedis:", err)
		return
	}
}

func DecrFavoriteRedis(videoID int64) {
	conn := RedisCache.AsynConn()
	defer func() {
		_ = conn.Close()
	}()
	favoriteKey := ID2FavoriteKey(videoID)
	_, err := conn.AsyncDo("ZINCRBY", FavoriteSortedSet, -1, favoriteKey)
	if err != nil {
		log.Print("err in DecrFavoriteRedis:", err)
		return
	}
}

func GetTopFavorite(n int) (top map[int64]int) {
	conn := RedisCache.Conn()
	defer func() {
		_ = conn.Close()
	}()
	values, err := redis.Values(conn.Do("ZREVRANGE", FavoriteSortedSet, 0, n, "WITHSCORES"))
	if err != nil {
		log.Println("err in GetTopFavorite:", err)
		return nil
	}
	top = make(map[int64]int)
	for i := 0; i < len(values); i += 2 {
		key, _ := redis.String(values[i], nil)
		v, _ := redis.Int64(values[i+1], nil)
		if FavoriteKey2ID(key) == 0 {
			continue
		}
		top[FavoriteKey2ID(key)] = int(v)
	}
	return top
}

func GetTotalFavoritedRedis(userID int64)(count int64){
	arr := GetVideoIDsByUser(userID)
	if len(arr) == 0 {
		return 0
	}
	for _, val := range arr {
		count += GetFavoriteNumRedis(val)
	}
	return count
}