package model

import (
	//"github.com/gomodule/redigo/redis"
	"github.com/gistao/RedisGo-Async/redis"
	"time"
)

const (
	FOREVER = time.Duration(-1)
)

type Cache struct {
	pool              *redis.Pool
	asyncPool         *redis.AsyncPool
	defaultExpiration time.Duration
}

func NewRedisCache(db int, host string, defaultExpiration time.Duration) *Cache {
	pool := &redis.Pool{
		MaxIdle:     100,
		MaxActive:   1000,
		IdleTimeout: time.Duration(100) * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", host, redis.DialDatabase(db))
			if err != nil {
				return nil, err
			}
			//if _, err = conn.Do("AUTH", "XXXXXX"); err != nil {
			//	conn.Close()
			//	return nil, err
			//}
			return conn, nil
		},
	}
	asyncPool := &redis.AsyncPool{
		Dial: func() (redis.AsynConn, error) {
			conn, err := redis.AsyncDial("tcp", host, redis.DialDatabase(db))
			if err != nil {
				return nil, err
			}

			return conn, nil
		},
		MaxGetCount: 1000,
	}
	return &Cache{pool: pool, asyncPool: asyncPool, defaultExpiration: defaultExpiration}
}

func (c *Cache) Conn() redis.Conn {
	return c.pool.Get()
}
func (c *Cache) AsynConn() redis.AsynConn {
	return c.asyncPool.Get()
}
