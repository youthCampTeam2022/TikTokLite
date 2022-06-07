package model

import (
	//"github.com/gomodule/redigo/redis"
	"github.com/gistao/RedisGo-Async/redis"
	"gorm.io/gorm"
	"strconv"
)

type Follow struct {
	gorm.Model
	UserID     int64 `gorm:"unique_index:idx_follow"`
	FollowerID int64 `gorm:"unique_index:idx_follow"`
}

type IFollowRepository interface {
	//Insert(follow *Follow) error
	//Delete(uid, fid int64) error
	//GetFollowerList(*[]Follow, int64) error
	//GetFollowList(*[]Follow, int64) error
	//上面是操作mysql接口，下面是操作redis接口，暂时先保留两种方式吧，以防万一。
	RedisInsert(uid, fid int64) error
	RedisDelete(uid, fid int64) error
	RedisGetFollowerList(uid int64) ([]int64, error)
	RedisGetFollowList(uid int64) ([]int64, error)
	RedisFollowCount(uid int64) int64
	RedisFollowerCount(uid int64) int64
	RedisIsFollow(userID, followID int64) bool
	GetName(userID int64) string
}
type FollowManagerRepository struct {
	db         *gorm.DB
	redisCache *Cache
}

func NewFollowManagerRepository() *FollowManagerRepository {
	return &FollowManagerRepository{DB, RedisCache}
}

//func (r *FollowManagerRepository) Insert(f *Follow) error {
//	tmp := &Follow{}
//	var err error
//	res := r.db.Where("user_id=? AND follower_id=?", f.UserID, f.FollowerID).First(tmp)
//	if res.Error != nil {
//		if res.Error == gorm.ErrRecordNotFound {
//			err = r.db.Create(f).Error
//
//		} else {
//			err = res.Error
//		}
//	}
//	return err
//	//return r.db.Create(f).Error
//}
//
//func (r *FollowManagerRepository) Delete(uid, fid int64) error {
//	return r.db.Unscoped().Where("follower_id=? AND user_id=?", fid, uid).Delete(&Follow{}).Error
//}
//
////GetFollowerList 获取粉丝用户列表
//func (r *FollowManagerRepository) GetFollowerList(followers *[]Follow, userID int64) error {
//	return r.db.Where("user_id=?", userID).Scan(followers).Error
//}
//
////GetFollowList 获取关注用户列表
//func (r *FollowManagerRepository) GetFollowList(follows *[]Follow, followerID int64) error {
//	return r.db.Where("follows.follower_id=?", followerID).Scan(follows).Error
//}

//id:follow表示id用户关注列表，id:fans表示id用户粉丝列表
func (r *FollowManagerRepository) RedisInsert(uid, fid int64) error {
	//conn := r.redisCache.Conn()
	conn := r.redisCache.Conn()
	defer func() {
		_ = conn.Close()
	}()
	followKey := strconv.Itoa(int(fid)) + ":follow"
	fansKey := strconv.Itoa(int(uid)) + ":fans"
	//开启事务
	conn.Send("MULTI")
	conn.Send("SADD", followKey, uid)
	conn.Send("SADD", fansKey, fid)
	//_, err := conn.Do("EXEC")
	_, err := conn.Do("EXEC")
	return err
}

func (r *FollowManagerRepository) RedisDelete(uid, fid int64) error {
	//conn := r.redisCache.Conn()
	conn := r.redisCache.Conn()
	defer func() {
		_ = conn.Close()
	}()
	followKey := strconv.Itoa(int(fid)) + ":follow"
	fansKey := strconv.Itoa(int(uid)) + ":fans"
	conn.Send("MULTI")
	conn.Send("SREM", followKey, uid)
	conn.Send("SREM", fansKey, fid)
	//_, err := conn.Do("EXEC")
	_, err := conn.Do("EXEC")
	return err
}

func (r *FollowManagerRepository) RedisGetFollowerList(uid int64) ([]int64, error) {
	conn := r.redisCache.Conn()
	defer func() {
		_ = conn.Close()
	}()
	fansKey := strconv.Itoa(int(uid)) + ":fans"
	//查询uid的粉丝列表获得粉丝id集合
	resp, err := redis.Int64s(conn.Do("SMEMBERS", fansKey))
	if err != nil {
		return nil, err
	}
	return resp, err
}
func (r *FollowManagerRepository) RedisGetFollowList(uid int64) ([]int64, error) {
	conn := r.redisCache.Conn()
	defer func() {
		_ = conn.Close()
	}()
	followKey := strconv.Itoa(int(uid)) + ":follow"
	//查询uid的关注列表获得关注id集合
	resp, err := redis.Int64s(conn.Do("SMEMBERS", followKey))
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (r *FollowManagerRepository) RedisFollowCount(uid int64) int64 {
	conn := r.redisCache.Conn()
	defer func() {
		_ = conn.Close()
	}()
	followKey := strconv.Itoa(int(uid)) + ":follow"
	//count
	resp, _ := conn.Do("SCARD", followKey)
	return resp.(int64)
}

func (r *FollowManagerRepository) RedisFollowerCount(uid int64) int64 {
	conn := r.redisCache.Conn()
	defer func() {
		_ = conn.Close()
	}()
	fansKey := strconv.Itoa(int(uid)) + ":fans"
	//count
	resp, _ := conn.Do("SCARD", fansKey)
	return resp.(int64)
}

//RedisIsFollow 判断userID是否有关注followID
func (r *FollowManagerRepository) RedisIsFollow(userID, followID int64) bool {
	conn := r.redisCache.Conn()
	defer func() {
		_ = conn.Close()
	}()
	followKey := strconv.Itoa(int(userID)) + ":follow"
	isFollow, _ := redis.Bool(conn.Do("SISMEMBER", followKey, followID))
	return isFollow
}

//GetName 先从redis查询是否有对应userID的name记录，如果没有就从数据库查找，然后写入redis
func (r *FollowManagerRepository) GetName(userID int64) string {
	conn := r.redisCache.Conn()
	defer func() {
		_ = conn.Close()
	}()
	key := strconv.Itoa(int(userID)) + ":name"
	name, err := redis.String(conn.Do("GET", key))
	//redis不存在记录
	if err == redis.ErrNil {
		//查找数据库
		r.db.Model(&User{}).Select("name").Where("id=?", userID).First(&name)
		//将数据库查找的name写入redis
		if name != "" {
			exp := 60 * 60 * 24 * 3
			_, _ = conn.Do("SET", key, name, "EX", exp)
		}
	}
	return name
}
