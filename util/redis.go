package util

import (
    "github.com/garyburd/redigo/redis"
)

func newRedisPool() *redis.Pool {
    return &redis.Pool{
        MaxIdle: 80,
        MaxActive: 12000, // max number of connections
        Dial: func() (redis.Conn, error) {
            c, err := redis.Dial("tcp", "127.0.0.1:6379")
            if err != nil {
                panic(err.Error())
            }
            return c, err
        },
    } 
}

var RedisPool = newRedisPool()
