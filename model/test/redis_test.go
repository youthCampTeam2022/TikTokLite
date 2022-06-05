package test

import (
	"TikTokLite/model"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"sort"
	"strconv"
	"testing"
)

func TestName(t *testing.T) {
	model.RedisInit()
	conn := model.RedisCache.Conn()
	defer func() {
		_ = conn.Close()
	}()
	fk1 := model.ID2FavoriteKey(4)
	fk2 := model.ID2FavoriteKey(5)
	fk3 := model.ID2FavoriteKey(6)
	_, _ = conn.Do("HSET","test1",fk1,0)
	_, _ = conn.Do("HSET","test1",fk2,56)
	_, _ = conn.Do("HSET","test1",fk3,2)

	values, err := redis.Values(conn.Do("HGETALL", "test1"))
	if err != nil {
		log.Print("err in GetTopFavorite:",err)
		return
	}
	var favoTop model.FavoriteTops
	for i := 0; i < len(values); i+=2 {
		b := values[i+1].([]uint8)
		num, _ := strconv.Atoi(string(b))
		favoTop = append(favoTop, model.FavoriteTop{
			Id:  model.FavoriteKey2ID(string(values[i].([]uint8))),
			Num: num,
		})
	}
	sort.Sort(favoTop)
	fmt.Println(favoTop)
}

func TestNameGetTopComment(t *testing.T) {
	//RedisInit()
	fk1 := model.ID2CommentKey(4)
	fmt.Println(model.CommentKey2ID(fk1))
}


