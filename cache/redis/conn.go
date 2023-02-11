package redis

import (
	"fmt"
	"time"

	"github.com/ggvylf/filestore/config"
	"github.com/gomodule/redigo/redis"
)

var (
	pool      *redis.Pool
	redisHost = config.RedisHost
	redisPass = config.RedisPass
	redisDb   = config.RedisDb
)

func init() {

	pool = newPool()
}

func RedisPool() *redis.Pool {
	return pool
}

// 初始化redis连接池
func newPool() *redis.Pool {
	return &redis.Pool{

		MaxIdle:     50,
		MaxActive:   30,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {

			c, err := redis.Dial("tcp", redisHost)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}

			// requirepass
			if _, err := c.Do("AUTH", redisPass); err != nil {
				c.Close()
				return nil, err
			}

			// 默认db[0]
			if _, err := c.Do("SELECT", redisDb); err != nil {
				c.Close()
				return nil, err
			}

			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {

				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}
