package models

import (
	"sync/atomic"
	"time"

	"github.com/garyburd/redigo/redis"
)

var clientSenders map[int]*MemberModel

func ClientList() *map[int]*MemberModel {
	if clientSenders == nil {
		clientSenders = make(map[int]*MemberModel)
	}
	return &clientSenders
}

var hash int32 = 1

//拿每个新连接的hash值
func GetNewHash() int32 {
	result := atomic.AddInt32(&hash, 1)
	return result
}

var chatSenders map[string]*ChatChannelModel

func ChatList() *map[string]*ChatChannelModel {
	if chatSenders == nil {
		chatSenders = make(map[string]*ChatChannelModel)
	}
	return &chatSenders
}

var pool *redis.Pool

func GetPool() *redis.Pool {
	return pool
}

//初始化一个pool
func NewPool(server string) *redis.Pool {

	pool = &redis.Pool{
		MaxIdle:     3,
		MaxActive:   100,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			// if _, err := c.Do("AUTH", password); err != nil {
			// 	c.Close()
			// 	return nil, err
			// }
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
	return pool
}
