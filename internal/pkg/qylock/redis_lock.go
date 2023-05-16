package qylock

import (
	"context"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
)

const (
	DEFAULT_EXP = 30
)

type RedisLockIface interface {
	Set(key string, expire time.Duration) (isLock bool, identify string, err error)
	Release(key, identify string) bool
	SetWithContext(ctx context.Context, key string, expire time.Duration) (bool, string, error)
	ReleaseWithContext(ctx context.Context, key string, randVal string) error
}

type RedisLock struct {
	redisPool *redis.Pool
}

func NewRedis(redisPool *redis.Pool) RedisLockIface {
	return &RedisLock{redisPool: redisPool}
}

func (r *RedisLock) Set(key string, expire time.Duration) (isLock bool, identify string, err error) {
	conn := r.redisPool.Get()
	defer conn.Close()

	identify = uuid.New().String()
	if expire < 0 {
		expire = DEFAULT_EXP * time.Second
	}
	reply, err := conn.Do("SET", key, identify, "NX", "EX", int(expire/time.Millisecond))
	result, err := redis.String(reply, err)
	return (err == nil && result == "OK"), identify, err
}

func (r *RedisLock) Release(key, identify string) bool {
	conn := r.redisPool.Get()
	defer conn.Close()
	var deleteScript = redis.NewScript(1, `
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("DEL", KEYS[1])
	else
		return 0
	end
`)
	status, err := deleteScript.Do(conn, key, identify)
	return err == nil && status != int32(0)
}

func (r *RedisLock) SetWithContext(ctx context.Context, key string, expire time.Duration) (bool, string, error) {
	conn := r.redisPool.Get()
	defer conn.Close()
	randVal := time.Now().Format("2006-01-02 15:04:05.000")
	var (
		reply any
		err   error
	)
	if expire < 0 {
		reply, err = conn.Do("SET", key, randVal, "NX", "PX", int(DEFAULT_EXP*1000))
	} else {
		reply, err = conn.Do("SET", key, randVal, "NX", "PX", int(expire/time.Millisecond))
	}
	if err != nil {
		return false, "", err
	}
	if reply == nil {
		return false, "", nil
	}

	return true, randVal, nil
}

func (r *RedisLock) ReleaseWithContext(ctx context.Context, key string, randVal string) error {
	conn := r.redisPool.Get()
	defer conn.Close()

	luaScript := `
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end;
	`
	script := redis.NewScript(1, luaScript)
	_, err := script.Do(conn, key, randVal)

	return err
}
