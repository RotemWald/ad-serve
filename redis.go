package main

import "github.com/gomodule/redigo/redis"

func NewRedisPool(idleConnections int, activeConnections int, address string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:   idleConnections,
		MaxActive: activeConnections,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", address)
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}
