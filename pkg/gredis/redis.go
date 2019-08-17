package gredis

import (
	"Gin-blog-example/pkg/setting"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"time"
)

// Redis 链接池
var RedisConn *redis.Pool

func Setup() error {
	RedisConn = &redis.Pool{
		MaxIdle:     setting.RedisSetting.MaxIdle,     //最大空闲连接数
		MaxActive:   setting.RedisSetting.MaxActive,   //在给定时间内，允许分配的最大连接数（当为零时，没有限制）
		IdleTimeout: setting.RedisSetting.IdleTimeout, //在给定时间内将会保持空闲状态，若到达时间限制则关闭连接（当为零时，没有限制）
		//提供创建和配置应用程序连接的一个函数
		Dial: func() (conn redis.Conn, e error) {
			//1.创建连接
			c, e := redis.Dial("tcp", setting.RedisSetting.Host)
			if e != nil {
				return nil, e
			}

			//2.访问认证
			if setting.RedisSetting.Password != "" {
				// 如果没有设置密码不用执行这句话
				if _, e = c.Do("AUTH", setting.RedisSetting.Password); e != nil {
					c.Close()
					return nil, e
				}
			}
			return c, e
		},
		// 每分钟去检测一下 redis 链接状况.如果出错了 redis 会自动关闭 shutdown
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}

	return nil
}

func Set(key string, data interface{}, time int) error {
	conn := RedisConn.Get()
	defer conn.Close()

	value, e := json.Marshal(data)

	if e != nil {
		return e
	}

	_, e = conn.Do("SET", key, value)
	if e != nil {
		return e
	}

	_, e = conn.Do("EXPIRE", key, time)
	if e != nil {
		return e
	}
	return nil
}

func Exists(key string) bool {
	conn := RedisConn.Get()
	defer conn.Close()

	exists, e := redis.Bool(conn.Do("EXISTS", key))
	if e != nil {
		return false
	}
	return exists
}

func Get(key string) ([]byte, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	reply, e := redis.Bytes(conn.Do("GET", key))
	if e != nil {
		return nil, e
	}
	return reply, nil
}

func Delete(key string) (bool, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("DEL", key))
}

func LikeDeletes(key string) error {
	conn := RedisConn.Get()
	defer conn.Close()
	keys, e := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if e != nil {
		return e
	}

	for _, key := range keys {
		_, e = Delete(key)
		if e != nil {
			return e
		}
	}
	return nil
}
