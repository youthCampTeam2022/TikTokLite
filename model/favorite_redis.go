package model

import (
	"github.com/gomodule/redigo/redis"
	"log"
	"sort"
	"strconv"
)

const FavoriteHash = "favorite"

func GetFavoriteNumRedis(videoID int64) (count int64) {
	conn := RedisCache.Conn()
	defer func() {
		_ = conn.Close()
	}()
	favoriteKey := ID2FavoriteKey(videoID)
	//count
	num, err := redis.Int64(conn.Do("HGET", FavoriteHash, favoriteKey))
	if err != nil  {
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
	conn := RedisCache.Conn()
	defer func() {
		_ = conn.Close()
	}()
	favoriteKey := ID2FavoriteKey(videoID)
	//todo: 过期时间没设置
	_, err := conn.Do("HSET", FavoriteHash, favoriteKey, num)
	if err != nil {
		log.Print("err in SetFavoriteNumRedis:", err)
		return
	}
}

func IncrFavoriteRedis(videoID int64) {
	conn := RedisCache.Conn()
	defer func() {
		_ = conn.Close()
	}()
	favoriteKey := ID2FavoriteKey(videoID)
	_, err := conn.Do("HINCRBY", FavoriteHash, favoriteKey, 1)
	if err != nil {
		log.Print("err in IncrFavoriteRedis:", err)
		return
	}
}

func DecrFavoriteRedis(videoID int64) {
	conn := RedisCache.Conn()
	defer func() {
		_ = conn.Close()
	}()
	favoriteKey := ID2FavoriteKey(videoID)
	_, err := conn.Do("HINCRBY", FavoriteHash, favoriteKey, -1)
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
	values, err := redis.Values(conn.Do("HGETALL", FavoriteHash))
	if err != nil {
		log.Print("err in GetTopFavorite:", err)
		return nil
	}
	var favoTop FavoriteTops
	for i := 0; i < len(values); i += 2 {
		b := values[i+1].([]uint8)
		num, _ := strconv.Atoi(string(b))
		favoTop = append(favoTop, FavoriteTop{
			Id:  FavoriteKey2ID(string(values[i].([]uint8))),
			Num: num,
		})
	}
	sort.Sort(favoTop)
	//现有的不够,从其他补（待定）
	//if n > favoTop.Len(){
	//	return nil
	//}
	top = make(map[int64]int)
	for i := 0; i < min(n, len(favoTop)); i++ {
		top[favoTop[i].Id] = favoTop[i].Num
	}
	return top
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type FavoriteTops []FavoriteTop

type FavoriteTop struct {
	Id  int64
	Num int
}

func (f FavoriteTops) Len() int {
	return len(f)
}

func (f FavoriteTops) Less(i, j int) bool {
	return f[i].Num > f[j].Num
}

func (f FavoriteTops) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}
