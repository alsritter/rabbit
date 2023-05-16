package cache

import (
	"context"
	"errors"

	"github.com/gomodule/redigo/redis"
)

var (
	ErrExpire = errors.New("expire failed")
)

const (
	// Redis 的时间单位是秒
	TimeUnit_MIN   = 60           // 1分钟
	TimeUnit_10MIN = 10 * 60      // 10分钟
	TimeUnit_HOUR  = 60 * 60      // 1小时
	TimeUnit_DAY   = 24 * 60 * 60 // 1天
)

// base 缓存基础结构
type base struct {
	ctx       context.Context
	redisPool *redis.Pool
	keyPrefix string
}

func (b *base) buildKey(key string) string {
	return b.keyPrefix + key
}

// checkTTL 检查 key 的 ttl(time to live)
func (b *base) checkTTL(conn redis.Conn, key string) bool {
	ttl, err := redis.Int(conn.Do("TTL", key))
	if err != nil {
		return false
	}
	if ttl <= 0 {
		return false
	}
	return true
}

// queryMake 添加查询的 Mark 针对一些已经查询过但是为空的数据
func (b *base) queryMark(conn redis.Conn, key string, expireSeconds int32) error {
	_, err := conn.Do("SET", "query_mark:"+key, "1", "EX", expireSeconds)
	if err != nil {
		return err
	}
	return nil
}

func (b *base) checkQueryMark(conn redis.Conn, key string) bool {
	exist, err := redis.Bool(conn.Do("EXISTS", "query_mark:"+key))
	if err != nil {
		return false
	}
	return exist
}
