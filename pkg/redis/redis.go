package redis

import (
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"
)

// Pool RedisPool连接池实例
var Pool *redis.Pool

// ConnRedis 初始化连接池
func ConnRedis() error {
	Pool = &redis.Pool{
		MaxIdle:     viper.GetInt("redis.maxIdle"),
		IdleTimeout: time.Duration(viper.GetInt("redis.idleTimeout")) * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", viper.GetString("redis.redisInfo"),
				redis.DialConnectTimeout(time.Duration(viper.GetInt("redis.connectTimeout"))*time.Second),
				redis.DialReadTimeout(time.Duration(viper.GetInt("redis.readTimeout"))*time.Second),
				redis.DialWriteTimeout(time.Duration(viper.GetInt("redis.writeTimeout"))*time.Second))
			if err != nil {
				return nil, err
			}
			if viper.GetString("redis.redisAuth") != "" {
				if _, err := c.Do("AUTH", viper.GetString("redis.redisAuth")); err != nil {
					c.Close()
					return nil, err
				}
			}
			if _, err := c.Do("SELECT", viper.GetInt("redis.db")); err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Duration(viper.GetInt("redis.idleTimeout"))*time.Second {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
	return Pool.Get().Err()
}

// ExecRedisCommand 执行redis命令
func ExecRedisCommand(command string, args ...interface{}) (interface{}, error) {
	redis := Pool.Get()
	defer redis.Close()

	return redis.Do(command, args...)
}
