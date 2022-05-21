package model

import (
	"github.com/gomodule/redigo/redis"
	"time"
)

const (
	FOREVER = time.Duration(-1)
)

type Cache struct {
	pool              *redis.Pool
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
	return &Cache{pool: pool, defaultExpiration: defaultExpiration}
}
func (c *Cache) Conn() redis.Conn {
	return c.pool.Get()
}
