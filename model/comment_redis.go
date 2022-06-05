package model

import (
	"github.com/gomodule/redigo/redis"
	"log"
	"sort"
	"strconv"
)

const CommentHash = "comment"

func GetCommentNumRedis(videoID int64) (count int64) {
	conn := RedisCache.Conn()
	defer func() {
		_ = conn.Close()
	}()
	commentKey := ID2CommentKey(videoID)
	//count
	num, err := conn.Do("HGET", CommentHash, commentKey)
	if err != nil {
		count = GetCommentNum(videoID)
		SetCommentNumRedis(videoID, count)
		return
	}
	return num.(int64)
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
	_, err := conn.Do("HSET", CommentHash, commentKey, num)
	if err != nil {
		log.Print("err in SetCommentNumRedis:", err)
		return
	}
}

func GetTopComment(n int) (top map[int64]int) {
	conn := RedisCache.Conn()
	defer func() {
		_ = conn.Close()
	}()
	values, err := redis.Values(conn.Do("HGETALL", CommentHash))
	if err != nil {
		log.Print("err in GetTopComment:", err)
		return nil
	}
	var commTop CommentTops
	for i := 0; i < len(values); i += 2 {
		b := values[i+1].([]uint8)
		num, _ := strconv.Atoi(string(b))
		commTop = append(commTop, CommentTop{
			id:  FavoriteKey2ID(string(values[i].([]uint8))),
			num: num,
		})
	}
	sort.Sort(commTop)
	//现有的不够,从其他补（待定）
	//if n > commTop.Len(){
	//	return nil
	//}
	top = make(map[int64]int)
	for i := 0; i < min(n, len(commTop)); i++ {
		top[commTop[i].id] = commTop[i].num
	}
	return top
}

type CommentTops []CommentTop

type CommentTop struct {
	id  int64
	num int
}

func (c CommentTops) Len() int {
	return len(c)
}

func (c CommentTops) Less(i, j int) bool {
	return c[i].num > c[j].num
}

func (c CommentTops) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
